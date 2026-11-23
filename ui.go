package main

import (
	"context"
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/andreacoppini/dogtag/sz"
)

//go:embed web
var webAssets embed.FS

// keepAliveEvery is how often an idle session is touched so the controller
// does not time it out under an operator who is still editing a spreadsheet.
const keepAliveEvery = 60 * time.Second

// uiServer is the browser console. It binds loopback only: nothing should be
// able to reconfigure a fleet of access points from off this machine.
type uiServer struct {
	opt  options
	addr string // host:port the console is served on

	mu       sync.Mutex
	subs     map[chan Event]struct{}
	history  []Event // the current job, replayed to a client that connects mid-run
	backend  Backend
	zones    []Zone
	stopKeep context.CancelFunc
	running  bool
	cancel   context.CancelFunc
	export   [][]string
	exportAs string
}

func serveUI(opt options) error {
	u := &uiServer{opt: opt, subs: map[chan Event]struct{}{}}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", opt.uiPort))
	if err != nil {
		return fmt.Errorf("cannot start the console: %w", err)
	}
	u.addr = fmt.Sprintf("127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
	link := "http://" + u.addr

	content, err := fs.Sub(webAssets, "web")
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/api/defaults", u.handleDefaults)
	mux.HandleFunc("/api/state", u.handleState)
	mux.HandleFunc("/api/connect", u.handleConnect)
	mux.HandleFunc("/api/disconnect", u.handleDisconnect)
	mux.HandleFunc("/api/parse", u.handleParse)
	mux.HandleFunc("/api/run", u.handleRun)
	mux.HandleFunc("/api/stop", u.handleStop)
	mux.HandleFunc("/api/events", u.handleEvents)
	mux.HandleFunc("/api/download", u.handleDownload)
	mux.HandleFunc("/api/update", u.handleUpdate)

	fmt.Printf("DogTag Plus %s\nConsole: %s\n(leave this window open; close it or press Ctrl-C to quit)\n", version, link)
	openBrowser(link)

	srv := &http.Server{Handler: u.guard(mux), ReadHeaderTimeout: 10 * time.Second}
	return srv.Serve(ln)
}

// guard refuses requests that did not come from the console's own page.
// Binding to loopback keeps other machines out, but not other web pages in
// the operator's browser: without this, any site could post to this port and
// drive the controller with the operator's session.
func (u *uiServer) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, port, _ := net.SplitHostPort(u.addr)
		okHost := r.Host == u.addr || r.Host == "localhost:"+port
		okOrigin := true
		if o := r.Header.Get("Origin"); o != "" {
			p, err := url.Parse(o)
			okOrigin = err == nil && (p.Host == u.addr || p.Host == "localhost:"+port)
		}
		if !okHost || !okOrigin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---- event fan-out ----

func (u *uiServer) publish(e Event) {
	u.mu.Lock()
	if e.Kind != EvSession {
		u.history = append(u.history, e)
	}
	subs := make([]chan Event, 0, len(u.subs))
	for c := range u.subs {
		subs = append(subs, c)
	}
	u.mu.Unlock()
	for _, c := range subs {
		select {
		case c <- e:
		default: // a slow browser must not stall the run
		}
	}
}

func (u *uiServer) logf(format string, a ...any) {
	u.publish(Event{Kind: EvLog, Message: fmt.Sprintf(format, a...)})
}

func (u *uiServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan Event, 1024)
	u.mu.Lock()
	u.subs[ch] = struct{}{}
	backlog := append([]Event(nil), u.history...)
	u.mu.Unlock()
	defer func() {
		u.mu.Lock()
		delete(u.subs, ch)
		u.mu.Unlock()
	}()

	send := func(e Event) bool {
		b, err := json.Marshal(e)
		if err != nil {
			return true
		}
		if _, err := fmt.Fprintf(w, "data: %s\n\n", b); err != nil {
			return false
		}
		flusher.Flush()
		return true
	}
	// Replay so a reload mid-run does not lose the results so far.
	for _, e := range backlog {
		if !send(e) {
			return
		}
	}
	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case e := <-ch:
			if !send(e) {
				return
			}
		case <-keepalive.C:
			if _, err := io.WriteString(w, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// ---- connection ----

// handleDefaults seeds the form from the flags the process was started with.
// Secrets given on the command line stay on this side: the page is told only
// that there is one.
func (u *uiServer) handleDefaults(w http.ResponseWriter, r *http.Request) {
	kind := "sz"
	if u.opt.r1 {
		kind = "r1"
	}
	writeJSON(w, map[string]any{
		"version":     version,
		"kind":        kind,
		"host":        u.opt.host,
		"port":        u.opt.port,
		"user":        u.opt.user,
		"hasPass":     u.opt.pass != "",
		"verifyTls":   u.opt.verifyTLS,
		"apiVersion":  u.opt.apiVersion,
		"r1Region":    u.opt.r1Region,
		"r1Tenant":    u.opt.r1Tenant,
		"r1ClientId":  u.opt.r1ClientID,
		"hasR1Secret": u.opt.r1Secret != "",
		"concurrency": u.opt.concurrency,
		"debug":       u.opt.debug,
	})
}

func (u *uiServer) stateLocked() map[string]any {
	st := map[string]any{"connected": u.backend != nil, "running": u.running, "hasExport": u.export != nil}
	if u.backend != nil {
		st["kind"] = u.backend.Kind()
		st["describe"] = u.backend.Describe()
		st["zones"] = u.zones
		st["overrides"] = u.backend.OverrideOptions()
	}
	return st
}

func (u *uiServer) handleState(w http.ResponseWriter, r *http.Request) {
	u.mu.Lock()
	defer u.mu.Unlock()
	writeJSON(w, u.stateLocked())
}

func (u *uiServer) handleConnect(w http.ResponseWriter, r *http.Request) {
	var req connectRequest
	if !readJSON(w, r, &req) {
		return
	}
	u.mu.Lock()
	busy := u.running
	u.mu.Unlock()
	if busy {
		httpErr(w, errors.New("a run is in progress; stop it first"))
		return
	}
	// A secret given on the command line is used when the form leaves the
	// field empty; it never travels to the browser.
	if req.Pass == "" {
		req.Pass = u.opt.pass
	}
	if req.Secret == "" {
		req.Secret = u.opt.r1Secret
	}
	u.drop(false)

	ctx, cancel := context.WithTimeout(r.Context(), 2*u.opt.timeout)
	defer cancel()
	b, err := connect(ctx, u.opt, req, func(s string) { u.logf("%s", s) })
	if err != nil {
		u.logf("Could not connect: %s", err)
		httpErr(w, err)
		return
	}
	zones, err := b.Zones(ctx)
	if err != nil {
		_ = b.Close(ctx)
		httpErr(w, fmt.Errorf("signed in, but could not list the zones: %w", err))
		return
	}
	keepCtx, stopKeep := context.WithCancel(context.Background())
	u.mu.Lock()
	u.backend, u.zones, u.stopKeep = b, zones, stopKeep
	st := u.stateLocked()
	u.mu.Unlock()
	go u.keepAlive(keepCtx, b)

	u.logf("Connected: %s. %d zones.", b.Describe(), len(zones))
	writeJSON(w, st)
}

// keepAlive touches the session while it is idle, and notices when the
// controller has ended it.
func (u *uiServer) keepAlive(ctx context.Context, b Backend) {
	t := time.NewTicker(keepAliveEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c, cancel := context.WithTimeout(ctx, u.opt.timeout)
			err := b.KeepAlive(c)
			cancel()
			if errors.Is(err, sz.ErrSessionEnded) {
				u.sessionEnded(b)
				return
			}
		}
	}
}

// sessionEnded forgets a session the controller has already dropped.
func (u *uiServer) sessionEnded(b Backend) {
	u.mu.Lock()
	if u.backend != b {
		u.mu.Unlock()
		return
	}
	u.backend, u.zones = nil, nil
	if u.stopKeep != nil {
		u.stopKeep()
	}
	u.mu.Unlock()
	u.logf("The controller ended the session; disconnected.")
	u.publish(Event{Kind: EvSession, Message: "The controller ended the session. Connect again to carry on."})
}

// drop ends the current session, if there is one.
func (u *uiServer) drop(announce bool) {
	u.mu.Lock()
	b := u.backend
	u.backend, u.zones = nil, nil
	if u.stopKeep != nil {
		u.stopKeep()
		u.stopKeep = nil
	}
	u.mu.Unlock()
	if b == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = b.Close(ctx)
	if announce {
		u.logf("Disconnected.")
	}
}

func (u *uiServer) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	u.mu.Lock()
	cancel := u.cancel
	u.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	u.drop(true)
	writeJSON(w, map[string]any{"connected": false})
}

// ---- files ----

// handleParse reads a CSV with the same parser the command line uses, so a
// file means the same thing on both surfaces.
func (u *uiServer) handleParse(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Kind string `json:"kind"` // "naming" or "macs"
		Text string `json:"text"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if req.Kind == "macs" {
		macs, skipped := parseMACs(req.Text)
		writeJSON(w, map[string]any{"macs": macs, "skipped": skipped})
		return
	}
	rows, skipped := parseNaming(req.Text)
	writeJSON(w, map[string]any{"rows": rows, "skipped": skipped})
}

// ---- run control ----

type runRequest struct {
	Job      string `json:"job"` // naming, overrides, channels, export
	ZoneID   string `json:"zoneId"`
	ZoneName string `json:"zoneName"`

	Rows []NameRow `json:"rows"`
	OverrideJob

	Concurrency int `json:"concurrency"`
}

func (u *uiServer) handleRun(w http.ResponseWriter, r *http.Request) {
	var req runRequest
	if !readJSON(w, r, &req) {
		return
	}
	needZone := req.Job == "channels" || req.Job == "export"
	switch {
	case req.Job == "naming" && len(req.Rows) == 0:
		httpErr(w, errors.New("no APs to update: load a CSV file first"))
		return
	case req.Job == "overrides" && len(req.MACs) == 0:
		httpErr(w, errors.New("no APs to work on: load a CSV file first"))
		return
	case needZone && req.ZoneID == "":
		httpErr(w, errors.New("pick a zone first"))
		return
	case req.Job != "naming" && req.Job != "overrides" && !needZone:
		httpErr(w, fmt.Errorf("unknown job %q", req.Job))
		return
	}
	for i := range req.Rows {
		if mac := normalMAC(req.Rows[i].MAC); mac != "" {
			req.Rows[i].MAC = mac
		}
	}
	conc := req.Concurrency
	if conc < 1 {
		conc = u.opt.concurrency
	}

	u.mu.Lock()
	b := u.backend
	if b == nil || u.running {
		busy := u.running
		u.mu.Unlock()
		if busy {
			httpErr(w, errors.New("a run is already in progress"))
		} else {
			httpErr(w, errors.New("not connected: connect to a controller first"))
		}
		return
	}
	u.running = true
	u.history = nil
	ctx, cancel := context.WithCancel(context.Background())
	u.cancel = cancel
	u.mu.Unlock()

	go func() {
		start := time.Now()
		var err error
		if rb, ok := b.(*r1Backend); ok {
			rb.forget() // each run looks the tenant's APs up afresh
		}
		switch req.Job {
		case "naming":
			u.logf("Updating %d APs, %d at a time", len(req.Rows), conc)
			err = runNaming(ctx, b, req.Rows, conc, u.publish)
		case "overrides":
			u.logf("Removing overrides on %d APs, %d at a time", len(req.MACs), conc)
			err = runOverrides(ctx, b, req.OverrideJob, conc, u.publish)
		case "channels":
			u.logf("Counting channels in %s", req.ZoneName)
			_, err = runChannels(ctx, b, req.ZoneID, req.ZoneName, conc, u.publish)
		case "export":
			u.logf("Exporting the APs in %s", req.ZoneName)
			var records [][]string
			if records, err = runExport(ctx, b, req.ZoneID, conc, u.publish); err == nil {
				u.mu.Lock()
				u.export = records
				u.exportAs = "DogTag-" + safeName(req.ZoneName) + ".csv"
				u.mu.Unlock()
			}
		}
		stopped := ctx.Err() != nil
		cancel()
		u.mu.Lock()
		u.running = false
		u.cancel = nil
		u.mu.Unlock()

		elapsed := time.Since(start).Round(time.Millisecond).String()
		switch {
		case errors.Is(err, sz.ErrSessionEnded):
			u.publish(Event{Kind: EvError, Job: req.Job, Message: "The controller ended the session part-way through."})
			u.sessionEnded(b)
		case stopped:
			u.publish(Event{Kind: EvDone, Job: req.Job, Phase: "stopped", Message: "Stopped.", Elapsed: elapsed})
		case err != nil:
			u.publish(Event{Kind: EvError, Job: req.Job, Message: err.Error()})
		default:
			u.publish(Event{Kind: EvDone, Job: req.Job, Phase: "done", Message: "All done in " + elapsed + ".", Elapsed: elapsed})
		}
	}()
	writeJSON(w, map[string]any{"ok": true})
}

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// safeName builds the export's file name from the zone, as the original did.
func safeName(zone string) string {
	s := unsafeName.ReplaceAllString(zone, "-")
	if s == "" || s == "-" {
		return "export"
	}
	return s
}

func (u *uiServer) handleStop(w http.ResponseWriter, r *http.Request) {
	u.mu.Lock()
	cancel := u.cancel
	u.mu.Unlock()
	if cancel != nil {
		cancel()
		u.logf("Stopping...")
	}
	writeJSON(w, map[string]any{"ok": true})
}

// handleDownload hands over the last zone export.
func (u *uiServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	u.mu.Lock()
	records, name := u.export, u.exportAs
	u.mu.Unlock()
	if records == nil {
		http.Error(w, "nothing exported yet", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	_ = csv.NewWriter(w).WriteAll(records)
}

// handleUpdate answers the console's own question about whether a newer
// release exists. It is separate from /api/defaults so that a slow or blocked
// GitHub cannot hold up the page.
func (u *uiServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	info, ok := checkForUpdate(r.Context(), releasesAPI, version, !u.opt.noUpdate)
	if !ok {
		writeJSON(w, map[string]any{"available": false})
		return
	}
	writeJSON(w, map[string]any{"available": true, "latest": info.Latest, "current": version, "url": info.URL})
}

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return false
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 32<<20)).Decode(v); err != nil {
		httpErr(w, err)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func httpErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// openBrowser is best-effort: the URL is printed either way.
func openBrowser(link string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", link)
	case "darwin":
		cmd = exec.Command("open", link)
	default:
		cmd = exec.Command("xdg-open", link)
	}
	_ = cmd.Start()
}
