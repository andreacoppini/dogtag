package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/andreacoppini/dogtag/sz"
)

// Event is one thing a running job has to say. The console streams these to
// the browser and the command line prints them, so both surfaces show the
// same run.
type Event struct {
	Kind     string        `json:"kind"`
	Job      string        `json:"job,omitempty"`
	Message  string        `json:"message,omitempty"`
	Phase    string        `json:"phase,omitempty"`
	Done     int           `json:"done"`
	Failed   int           `json:"failed"`
	Total    int           `json:"total"`
	Row      *RowResult    `json:"row,omitempty"`
	Channels *ChannelCount `json:"channels,omitempty"`
	Elapsed  string        `json:"elapsed,omitempty"`
}

// Event kinds.
const (
	EvLog      = "log"
	EvPhase    = "phase"
	EvRow      = "row"
	EvChannels = "channels"
	EvDone     = "done"
	EvError    = "error"
	EvSession  = "session"
)

// Emitter receives a job's events.
type Emitter func(Event)

// RowResult is the outcome for one AP in a naming or overrides run. OK and
// Error are both lists because clearing overrides is many small requests per
// AP and some can succeed while others fail.
type RowResult struct {
	Index int    `json:"index"`
	ID    string `json:"id"` // MAC address or serial number, as the file gave it
	OK    string `json:"ok"`
	Error string `json:"error"`
}

// ChannelCount is how many connected APs sit on each channel in a zone.
type ChannelCount struct {
	Zone      string    `json:"zone"`
	Started   string    `json:"started"`
	Ended     string    `json:"ended"`
	Total     int       `json:"total"`
	Connected int       `json:"connected"`
	Band24    []ChanRow `json:"band24"`
	Band50    []ChanRow `json:"band50"`
	Band6     []ChanRow `json:"band6"`
}

// ChanRow is one channel and the number of APs on it.
type ChanRow struct {
	Channel int `json:"channel"`
	APs     int `json:"aps"`
}

// OverrideJob says what to clear or force on each AP.
type OverrideJob struct {
	IDs       []string `json:"ids"`
	Overrides []string `json:"overrides"`
	Radio24   string   `json:"radio24"` // "on", "off" or "skip"
	Radio50   string   `json:"radio50"`
	IPDynamic bool     `json:"ipDynamic"`
}

// tally counts progress across workers.
type tally struct {
	done, failed atomic.Int64
	total        int
}

func (t *tally) event(job string, row *RowResult) Event {
	return Event{Kind: EvRow, Job: job, Row: row, Total: t.total,
		Done: int(t.done.Load()), Failed: int(t.failed.Load())}
}

// forEach runs fn for every index with at most conc in flight, and stops
// handing out work once ctx is cancelled.
func forEach(ctx context.Context, n, conc int, fn func(i int)) {
	if conc < 1 {
		conc = 1
	}
	var wg sync.WaitGroup
	next := make(chan int)
	for w := 0; w < min(conc, n); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				fn(i)
			}
		}()
	}
feed:
	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			break feed
		case next <- i:
		}
	}
	close(next)
	wg.Wait()
}

// sessionLost reports whether err means the whole run is pointless, not just
// this AP.
func sessionLost(err error) bool { return errors.Is(err, sz.ErrSessionEnded) }

// runNaming applies a naming file: each AP on its own, many at once. addTo,
// when not empty, is the venue that APs the controller does not have yet are
// added to; otherwise such rows fail.
func runNaming(ctx context.Context, b Backend, rows []NameRow, addTo string, conc int, emit Emitter) error {
	const job = "naming"
	t := &tally{total: len(rows)}
	emit(Event{Kind: EvPhase, Job: job, Phase: "applying", Total: t.total})
	var lost atomic.Bool
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	forEach(ctx, len(rows), conc, func(i int) {
		row := rows[i]
		res := &RowResult{Index: i, ID: row.ID}
		note, err := b.ApplyName(ctx, row, addTo)
		switch {
		case err == nil:
			res.OK = firstNonEmpty(note, "OK")
			t.done.Add(1)
		case ctx.Err() != nil:
			return // stopped mid-request: the row is simply not done
		default:
			if sessionLost(err) {
				lost.Store(true)
				cancel()
			}
			res.Error = err.Error()
			t.failed.Add(1)
		}
		emit(t.event(job, res))
	})
	if lost.Load() {
		return sz.ErrSessionEnded
	}
	return nil
}

// runOverrides clears the chosen overrides on each AP and applies the radio
// and IP choices. The requests for one AP go in order; the APs go in
// parallel.
func runOverrides(ctx context.Context, b Backend, j OverrideJob, conc int, emit Emitter) error {
	const job = "overrides"
	t := &tally{total: len(j.IDs)}
	emit(Event{Kind: EvPhase, Job: job, Phase: "removing", Total: t.total})
	var lost atomic.Bool
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	forEach(ctx, len(j.IDs), conc, func(i int) {
		id := j.IDs[i]
		var ok, bad []string
		note := func(label string, err error) {
			if err == nil {
				ok = append(ok, label)
				return
			}
			if ctx.Err() != nil && !sessionLost(err) {
				return // cut short by Stop, which is not the AP's failure
			}
			if sessionLost(err) {
				lost.Store(true)
				cancel()
			}
			bad = append(bad, fmt.Sprintf("%s (%s)", label, err))
		}
		if err := b.Check(ctx, id); err != nil {
			// An AP the controller cannot act on fails once, with the reason,
			// rather than once for every setting.
			if sessionLost(err) {
				lost.Store(true)
				cancel()
			} else if ctx.Err() != nil {
				return
			}
			bad = []string{err.Error()}
		} else {
			b.RemoveOverrides(ctx, id, j.Overrides, note)
			if on, set := radioChoice(j.Radio24); set && ctx.Err() == nil {
				note("2G", b.SetRadio(ctx, id, "24", on))
			}
			if on, set := radioChoice(j.Radio50); set && ctx.Err() == nil {
				note("5G", b.SetRadio(ctx, id, "50", on))
			}
			if j.IPDynamic && ctx.Err() == nil {
				note("IP", b.SetDynamicIP(ctx, id))
			}
		}
		if ctx.Err() != nil && !lost.Load() {
			return
		}
		if len(bad) > 0 {
			t.failed.Add(1)
		} else {
			t.done.Add(1)
		}
		emit(t.event(job, &RowResult{Index: i, ID: id, OK: strings.Join(ok, ", "), Error: strings.Join(bad, ", ")}))
	})
	if lost.Load() {
		return sz.ErrSessionEnded
	}
	return nil
}

func radioChoice(s string) (on, set bool) {
	switch strings.ToLower(s) {
	case "on":
		return true, true
	case "off":
		return false, true
	}
	return false, false
}

// inventory reads a zone's APs. With detail, each AP costs a further request
// for the fields the listing does not carry; an AP whose detail cannot be
// read keeps its place with what the listing said about it.
func inventory(ctx context.Context, b Backend, job, zoneID string, detail bool, conc int, emit Emitter) ([]APInfo, error) {
	emit(Event{Kind: EvPhase, Job: job, Phase: "listing APs"})
	aps, err := b.ListAPs(ctx, zoneID)
	if err != nil {
		return nil, err
	}
	if len(aps) == 0 {
		return nil, errors.New("there are no APs here")
	}
	if !detail {
		emit(Event{Kind: EvLog, Message: fmt.Sprintf("%d APs", len(aps))})
		return aps, nil
	}
	emit(Event{Kind: EvLog, Message: fmt.Sprintf("%d APs; reading each one", len(aps))})
	emit(Event{Kind: EvPhase, Job: job, Phase: "reading APs", Total: len(aps)})

	t := &tally{total: len(aps)}
	var lost atomic.Bool
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	forEach(ctx, len(aps), conc, func(i int) {
		if err := b.Detail(ctx, &aps[i]); err != nil {
			if ctx.Err() != nil {
				return
			}
			if sessionLost(err) {
				lost.Store(true)
				cancel()
			}
			t.failed.Add(1)
			emit(Event{Kind: EvLog, Message: fmt.Sprintf("%s: %s", firstNonEmpty(aps[i].MAC, aps[i].Serial), err)})
		} else {
			t.done.Add(1)
		}
		emit(Event{Kind: EvPhase, Job: job, Phase: "reading APs", Total: t.total,
			Done: int(t.done.Load()), Failed: int(t.failed.Load())})
	})
	if lost.Load() {
		return nil, sz.ErrSessionEnded
	}
	if ctx.Err() != nil {
		return nil, errStopped
	}
	return aps, nil
}

var errStopped = errors.New("stopped")

// runExport builds the zone export: a header plus one record per AP.
func runExport(ctx context.Context, b Backend, zoneID string, conc int, emit Emitter) ([][]string, error) {
	aps, err := inventory(ctx, b, "export", zoneID, true, conc, emit)
	if err != nil {
		return nil, err
	}
	// The first column is whatever the controller knows the AP by, so that
	// the file loads straight back in: the MAC on SmartZone, the serial number
	// on RUCKUS One, where an AP that has not come online has no MAC yet.
	bySerial := b.Kind() == "r1"
	header := append([]string(nil), exportHeader...)
	if bySerial {
		header[0], header[10] = "Serial", "MAC"
	}
	records := [][]string{header}
	for _, ap := range aps {
		id, other := ap.MAC, ap.Serial
		if bySerial {
			id, other = ap.Serial, ap.MAC
		}
		records = append(records, []string{
			id, ap.Name, ap.Description, ap.Location, ap.Latitude, ap.Longitude,
			"", ap.Ch24, // channel width is not reported, so it is left for the operator
			"", ap.Ch50,
			other, ap.Status, ap.Model, ap.IPType, ap.IP,
		})
	}
	return records, nil
}

// channelText blanks a channel of 0, which is how an absent or switched-off
// radio reads; left in, it would be loaded back as a channel to pin.
func channelText(s string) string {
	if n, err := strconv.Atoi(s); err != nil || n <= 0 {
		return ""
	}
	return s
}

// runChannels counts the connected APs on each channel of a zone.
func runChannels(ctx context.Context, b Backend, zoneID, zoneName string, conc int, emit Emitter) (*ChannelCount, error) {
	out := &ChannelCount{Zone: zoneName, Started: time.Now().Format("2006-01-02 15:04:05")}
	aps, err := inventory(ctx, b, "channels", zoneID, !b.ListHasChannels(), conc, emit)
	if err != nil {
		return nil, err
	}
	out.Total = len(aps)
	b24, b50, b6 := map[int]int{}, map[int]int{}, map[int]int{}
	for _, ap := range aps {
		if !ap.Connected {
			continue // only connected APs are on a channel at all
		}
		out.Connected++
		count := func(m map[int]int, v string) {
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				m[n]++
			}
		}
		count(b24, ap.Ch24)
		count(b50, ap.Ch50)
		count(b6, ap.Ch6)
	}
	out.Band24, out.Band50, out.Band6 = chanRows(b24), chanRows(b50), chanRows(b6)
	out.Ended = time.Now().Format("2006-01-02 15:04:05")
	emit(Event{Kind: EvChannels, Job: "channels", Channels: out})
	return out, nil
}

func chanRows(m map[int]int) []ChanRow {
	rows := make([]ChanRow, 0, len(m))
	for ch, n := range m {
		rows = append(rows, ChanRow{Channel: ch, APs: n})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Channel < rows[j].Channel })
	return rows
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
