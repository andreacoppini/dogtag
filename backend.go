package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/andreacoppini/dogtag/r1"
	"github.com/andreacoppini/dogtag/sz"
)

// Backend is whatever manages the access points: a SmartZone controller or a
// RUCKUS One tenant. The jobs are written against this, so both do the same
// things from the same files.
type Backend interface {
	Kind() string     // "sz" or "r1"
	Describe() string // one line for the log and the status box
	// Zones lists the groupings APs live in: zones on SmartZone, venues on R1.
	Zones(ctx context.Context) ([]Zone, error)
	// ListAPs lists a zone's APs with whatever the listing itself carries.
	ListAPs(ctx context.Context, zoneID string) ([]APInfo, error)
	// Detail fills in the rest of one AP, at the cost of a request.
	Detail(ctx context.Context, ap *APInfo) error
	// ListHasChannels reports whether ListAPs already carries status and
	// channels, which saves the channel count a request per AP.
	ListHasChannels() bool
	// Check says whether the controller can act on this AP at all, so a row
	// that cannot work fails once instead of once per request. id is a MAC
	// address or a serial number, as the file gave it.
	Check(ctx context.Context, id string) error
	ApplyName(ctx context.Context, row NameRow) error
	OverrideOptions() []OverrideOption
	// RemoveOverrides hands the named settings back to the zone or venue,
	// reporting each one through note.
	RemoveOverrides(ctx context.Context, id string, ids []string, note func(label string, err error))
	SetRadio(ctx context.Context, id, band string, on bool) error
	SetDynamicIP(ctx context.Context, id string) error
	KeepAlive(ctx context.Context) error
	Close(ctx context.Context) error
}

// Zone is a SmartZone zone or an R1 venue.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// APInfo is one access point, in the terms of the export file.
type APInfo struct {
	MAC, Name, Description, Location string
	Latitude, Longitude              string
	Ch24, Ch50, Ch6                  string
	Serial, Status, Model            string
	IPType, IP                       string
	Connected                        bool
	venue                            string // R1 only
}

// OverrideOption is one per-AP setting that can be handed back.
type OverrideOption struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Default bool   `json:"default"`
}

// ---------------------------------------------------------------- SmartZone

type szBackend struct{ c *sz.Client }

func (b *szBackend) Kind() string { return "sz" }

func (b *szBackend) Describe() string {
	return fmt.Sprintf("SmartZone %s at %s (API %s)", b.c.ControllerVersion(), b.c.Host(), b.c.APIVersion())
}

func (b *szBackend) Zones(ctx context.Context) ([]Zone, error) {
	zs, err := b.c.Zones(ctx)
	out := make([]Zone, len(zs))
	for i, z := range zs {
		out[i] = Zone{ID: z.ID, Name: z.Name}
	}
	return out, err
}

func (b *szBackend) ListAPs(ctx context.Context, zoneID string) ([]APInfo, error) {
	aps, err := b.c.ZoneAPs(ctx, zoneID)
	out := make([]APInfo, len(aps))
	for i, a := range aps {
		out[i] = APInfo{MAC: a.MAC, Name: a.Name}
	}
	return out, err
}

func (b *szBackend) ListHasChannels() bool { return false }

func (b *szBackend) Detail(ctx context.Context, ap *APInfo) error {
	s, err := b.c.APSummary(ctx, ap.MAC)
	if err != nil {
		return err
	}
	ap.Name = firstNonEmpty(ap.Name, s["name"])
	ap.Description, ap.Location = s["description"], s["location"]
	ap.Latitude, ap.Longitude = s["latitude"], s["longitude"]
	ap.Ch24, ap.Ch50 = channelText(s["wifi24Channel"]), channelText(s["wifi50Channel"])
	ap.Ch6 = channelText(firstNonEmpty(s["wifi6gChannel"], s["wifi60Channel"]))
	ap.Serial, ap.Status, ap.Model = s["serial"], s["connectionState"], s["model"]
	ap.IPType, ap.IP = s["ipType"], s["ip"]
	ap.Connected = s["connectionState"] == "Connect"
	return nil
}

// Check refuses a serial number: SmartZone's API addresses an AP by its MAC
// address and nothing else.
func (b *szBackend) Check(_ context.Context, id string) error {
	if isSerial(id) {
		return errors.New("SmartZone identifies APs by MAC address; this row gives a serial number")
	}
	return nil
}

func (b *szBackend) ApplyName(ctx context.Context, row NameRow) error {
	if err := b.Check(ctx, row.ID); err != nil {
		return err
	}
	p, err := row.patch()
	if err != nil {
		return err
	}
	return b.c.PatchAP(ctx, row.ID, p)
}

// szOverrides are the per-AP overrides SmartZone can clear. The names are the
// API's own.
var szOverrides = []string{"login", "syslog", "wifi24", "wifi50", "wlanGroup24", "wlanGroup50",
	"smartMonitor", "bonjourGateway", "clientAdmissionControl24", "clientAdmissionControl50",
	"venueProfile", "specific", "channelEvaluationInterval", "autoChannelSelection24",
	"autoChannelSelection50", "picture", "meshOptions"}

func (b *szBackend) OverrideOptions() []OverrideOption {
	out := make([]OverrideOption, len(szOverrides))
	for i, o := range szOverrides {
		out[i] = OverrideOption{ID: o, Label: o, Default: true}
	}
	return out
}

func (b *szBackend) RemoveOverrides(ctx context.Context, mac string, ids []string, note func(string, error)) {
	for i, id := range ids {
		if ctx.Err() != nil {
			return
		}
		err := b.c.DeleteOverride(ctx, mac, id)
		note(id, err)
		// An AP the controller does not know fails every request the same
		// way; one line says so as well as seventeen.
		var api *sz.APIError
		if i == 0 && errors.As(err, &api) && api.Status == http.StatusNotFound {
			return
		}
	}
}

func (b *szBackend) SetRadio(ctx context.Context, mac, band string, on bool) error {
	return b.c.PatchAP(ctx, mac, map[string]any{"wlanService" + band + "Enabled": on})
}

func (b *szBackend) SetDynamicIP(ctx context.Context, mac string) error {
	return b.c.PatchAPNetwork(ctx, mac, map[string]any{"ipType": "Dynamic"})
}

func (b *szBackend) KeepAlive(ctx context.Context) error { return b.c.KeepAlive(ctx) }
func (b *szBackend) Close(ctx context.Context) error     { return b.c.Logout(ctx) }

// --------------------------------------------------------------- RUCKUS One

// r1Backend finds APs by whichever ID a file gives. R1's own key is the
// serial number, and it is the only one an AP has until it has come online:
// an AP added to a venue but not yet onboarded reports no MAC address. So the
// tenant's AP list is read once and indexed by serial and, where there is
// one, by MAC as well.
type r1Backend struct {
	c      *r1.Client
	tenant string

	mu    sync.Mutex
	index map[string]r1.AP // by serial number and by MAC
}

func (b *r1Backend) Kind() string { return "r1" }

func (b *r1Backend) Describe() string {
	return fmt.Sprintf("RUCKUS One tenant %s at %s", b.tenant, b.c.Host())
}

func (b *r1Backend) Zones(ctx context.Context) ([]Zone, error) {
	vs, err := b.c.Venues(ctx)
	out := make([]Zone, len(vs))
	for i, v := range vs {
		out[i] = Zone{ID: v.ID, Name: v.Name}
	}
	sort.SliceStable(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, err
}

func r1Connected(status string) bool {
	// R1 statuses read "2_00_Operational": the leading digit is the class,
	// and 2 is the one where the AP is up and talking to the cloud.
	return strings.HasPrefix(status, "2_") || strings.Contains(strings.ToLower(status), "operational")
}

func (b *r1Backend) ListAPs(ctx context.Context, zoneID string) ([]APInfo, error) {
	aps, err := b.c.APs(ctx, zoneID)
	if err != nil {
		return nil, err
	}
	out := make([]APInfo, 0, len(aps))
	b.mu.Lock()
	if b.index == nil {
		b.index = map[string]r1.AP{}
	}
	for _, a := range aps {
		b.remember(a)
		info := APInfo{MAC: a.MAC, Name: a.Name, Serial: a.Serial, Status: a.Status, Model: a.Model,
			IP: a.IP, IPType: a.IPType, Connected: r1Connected(a.Status), venue: a.VenueID}
		for _, r := range a.Radios {
			if r.Channel <= 0 {
				continue
			}
			ch := strconv.Itoa(r.Channel)
			switch band := strings.TrimSpace(r.Band); {
			case strings.HasPrefix(band, "2"):
				info.Ch24 = ch
			case strings.HasPrefix(band, "6"):
				info.Ch6 = ch
			default:
				info.Ch50 = firstNonEmpty(info.Ch50, ch)
			}
		}
		out = append(out, info)
	}
	b.mu.Unlock()
	return out, nil
}

func (b *r1Backend) ListHasChannels() bool { return true }

func (b *r1Backend) Detail(ctx context.Context, ap *APInfo) error {
	d, err := b.c.Detail(ctx, ap.venue, ap.Serial)
	if err != nil {
		return err
	}
	ap.Name = firstNonEmpty(d.Name, ap.Name)
	ap.Description, ap.Latitude, ap.Longitude = d.Description, d.Latitude, d.Longitude
	return nil
}

// find resolves a serial number or a MAC to the AP R1 knows it as, reading
// the tenant's AP list the first time one is asked for.
func (b *r1Backend) find(ctx context.Context, id string) (r1.AP, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.index == nil {
		aps, err := b.c.APs(ctx, "")
		if err != nil {
			return r1.AP{}, err
		}
		b.index = make(map[string]r1.AP, 2*len(aps))
		for _, a := range aps {
			b.remember(a)
		}
	}
	ap, ok := b.index[id]
	switch {
	case ok:
		return ap, nil
	case isSerial(id):
		return r1.AP{}, errors.New("no AP with this serial number in the tenant: add it to a venue in RUCKUS One first")
	}
	return r1.AP{}, errors.New("no AP with this MAC address in the tenant; an AP that has not come online yet can only be found by serial number")
}

// remember indexes one AP. The caller holds the lock.
func (b *r1Backend) remember(a r1.AP) {
	if a.Serial != "" {
		b.index[a.Serial] = a
	}
	if a.MAC != "" {
		b.index[a.MAC] = a
	}
}

func (b *r1Backend) Check(ctx context.Context, id string) error {
	_, err := b.find(ctx, id)
	return err
}

// forget drops the index, so the next run sees APs added since the last.
func (b *r1Backend) forget() {
	b.mu.Lock()
	b.index = nil
	b.mu.Unlock()
}

func (b *r1Backend) ApplyName(ctx context.Context, row NameRow) error {
	if (row.Latitude == "") != (row.Longitude == "") {
		return errors.New("latitude and longitude must be given together")
	}
	for _, v := range []string{row.Latitude, row.Longitude} {
		if _, err := strconv.ParseFloat(v, 64); v != "" && err != nil {
			return fmt.Errorf("latitude and longitude must be numbers (got %q)", v)
		}
	}
	bands := map[string]r1.RadioChange{}
	for band, cols := range map[string][2]string{"24": {row.Width24, row.Channel24}, "50": {row.Width50, row.Channel50}} {
		var ch r1.RadioChange
		if cols[0] != "" {
			n, err := strconv.Atoi(cols[0])
			if err != nil {
				return fmt.Errorf("channel width %q is not a number", cols[0])
			}
			ch.Width = &n
		}
		if cols[1] != "" {
			n, err := strconv.Atoi(cols[1])
			if err != nil {
				return fmt.Errorf("channel %q is not a number", cols[1])
			}
			ch.Channel = &n
		}
		if ch.Width != nil || ch.Channel != nil {
			bands[band] = ch
		}
	}
	record := row.Name != "" || row.Description != "" || row.Latitude != ""
	if !record && len(bands) == 0 {
		if row.Location != "" {
			return errors.New("nothing to change: RUCKUS One has no location field, and the row sets nothing else")
		}
		return errors.New("nothing to change: every column after the first is empty")
	}
	ap, err := b.find(ctx, row.ID)
	if err != nil {
		return err
	}
	if record {
		err := b.c.UpdateAP(ctx, ap.VenueID, ap.Serial, r1.Detail{Name: row.Name, Description: row.Description,
			Latitude: row.Latitude, Longitude: row.Longitude})
		if err != nil {
			return err
		}
	}
	if len(bands) > 0 {
		if err := b.c.ChangeRadios(ctx, ap.VenueID, ap.Serial, bands); err != nil {
			return fmt.Errorf("radio settings: %w", err)
		}
	}
	return nil
}

// r1Overrides are the per-AP settings R1 lets an AP take from its venue
// instead. The ones unticked by default exist only on some models, where
// asking for them on the rest is an error per AP.
var r1Overrides = []OverrideOption{
	{"radio24", "2.4 GHz radio settings", true},
	{"radio50", "5 GHz radio settings", true},
	{"radio6", "6 GHz radio settings", true},
	{"ledSettings", "LED", true},
	{"smartMonitorSettings", "Smart monitor", true},
	{"bssColoringSettings", "BSS colouring", true},
	{"clientAdmissionControlSettings", "Client admission control", true},
	{"directedMulticastSettings", "Directed multicast", true},
	{"managementTrafficVlanSettings", "Management traffic VLAN", true},
	{"stickyClientSteeringSettings", "Sticky client steering", true},
	{"bandModeSettings", "Band mode", false},
	{"lanPortSpecificSettings", "LAN port PoE settings", false},
	{"usbPortSettings", "USB port", false},
	{"iotSettings", "IoT", false},
	{"antennaTypeSettings", "Antenna type", false},
	{"externalAntennaSettings", "External antenna", false},
}

func (b *r1Backend) OverrideOptions() []OverrideOption { return r1Overrides }

func (b *r1Backend) RemoveOverrides(ctx context.Context, mac string, ids []string, note func(string, error)) {
	ap, err := b.find(ctx, mac)
	if err != nil {
		note("lookup", err)
		return
	}
	// The three radio entries live in one resource, so they go as one change.
	bands := map[string]r1.RadioChange{}
	var radioIDs []string
	for _, id := range ids {
		if band, ok := strings.CutPrefix(id, "radio"); ok {
			bands[band] = r1.RadioChange{UseVenue: true}
			radioIDs = append(radioIDs, id)
		}
	}
	if len(bands) > 0 {
		err := b.c.ChangeRadios(ctx, ap.VenueID, ap.Serial, bands)
		note(strings.Join(radioIDs, "+"), err)
	}
	for _, id := range ids {
		if strings.HasPrefix(id, "radio") {
			continue
		}
		if ctx.Err() != nil {
			return
		}
		note(id, b.c.UseVenueSettings(ctx, ap.VenueID, ap.Serial, id))
	}
}

func (b *r1Backend) SetRadio(ctx context.Context, mac, band string, on bool) error {
	ap, err := b.find(ctx, mac)
	if err != nil {
		return err
	}
	return b.c.ChangeRadios(ctx, ap.VenueID, ap.Serial, map[string]r1.RadioChange{band: {Enabled: &on}})
}

func (b *r1Backend) SetDynamicIP(ctx context.Context, mac string) error {
	ap, err := b.find(ctx, mac)
	if err != nil {
		return err
	}
	return b.c.SetDynamicIP(ctx, ap.VenueID, ap.Serial)
}

// The R1 token is renewed on demand, so there is no session to keep warm or
// to end.
func (b *r1Backend) KeepAlive(context.Context) error { return nil }
func (b *r1Backend) Close(context.Context) error     { return nil }

// ------------------------------------------------------------------ connect

// connectRequest is the Controller section of the console, and the same
// fields the command line takes.
type connectRequest struct {
	Kind string `json:"kind"` // "sz" or "r1"

	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	VerifyTLS bool   `json:"verifyTls"`
	APIVer    string `json:"apiVersion"`
	Debug     bool   `json:"debug"`

	Region   string `json:"r1Region"`
	Tenant   string `json:"r1Tenant"`
	ClientID string `json:"r1ClientId"`
	Secret   string `json:"r1Secret"`
}

// connect signs in to whichever kind of controller the request names.
func connect(ctx context.Context, opt options, req connectRequest, trace func(string)) (Backend, error) {
	if !opt.debug && !req.Debug {
		trace = nil
	}
	if req.Kind == "r1" {
		c, err := r1.New(ctx, r1.Options{Region: req.Region, Tenant: strings.TrimSpace(req.Tenant),
			ClientID: strings.TrimSpace(req.ClientID), Secret: req.Secret,
			BaseURL: opt.r1URL, Timeout: opt.timeout, Trace: trace})
		if err != nil {
			return nil, err
		}
		return &r1Backend{c: c, tenant: strings.TrimSpace(req.Tenant)}, nil
	}
	if strings.TrimSpace(req.Host) == "" || req.User == "" || req.Pass == "" {
		return nil, errors.New("please enter the SmartZone address, username and password")
	}
	c, err := sz.New(sz.Options{Host: req.Host, Port: req.Port, APIVersion: req.APIVer,
		VerifyTLS: req.VerifyTLS, Timeout: opt.timeout, Trace: trace})
	if err != nil {
		return nil, err
	}
	if err := c.Login(ctx, req.User, req.Pass); err != nil {
		return nil, fmt.Errorf("could not sign in to SmartZone: %w", err)
	}
	return &szBackend{c: c}, nil
}
