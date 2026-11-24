package main

import (
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andreacoppini/dogtag/internal/fakesz"
	"github.com/andreacoppini/dogtag/sz"
)

func szBackendFor(t *testing.T, f *fakesz.Server) Backend {
	t.Helper()
	ts := httptest.NewTLSServer(f)
	t.Cleanup(ts.Close)
	u, _ := url.Parse(ts.URL)
	b, err := connect(context.Background(), options{timeout: 5 * time.Second},
		connectRequest{Kind: "sz", Host: u.Hostname(), Port: u.Port(), User: f.User, Pass: f.Pass}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func r1BackendFor(t *testing.T, f *fakesz.R1) Backend {
	t.Helper()
	ts := httptest.NewServer(f)
	t.Cleanup(ts.Close)
	b, err := connect(context.Background(), options{timeout: 5 * time.Second, r1URL: ts.URL},
		connectRequest{Kind: "r1", Tenant: f.Tenant, ClientID: f.ClientID, Secret: f.Secret}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// collect gathers a job's row results by index.
func collect() (Emitter, func() map[int]RowResult) {
	var mu sync.Mutex
	rows := map[int]RowResult{}
	return func(e Event) {
			if e.Kind == EvRow {
				mu.Lock()
				rows[e.Row.Index] = *e.Row
				mu.Unlock()
			}
		}, func() map[int]RowResult {
			mu.Lock()
			defer mu.Unlock()
			return rows
		}
}

func TestParseNaming(t *testing.T) {
	text := "\ufeffMAC,Name,Description\r\n" +
		`d4-68-4d-10-20-a0,"Lobby, front","says ""hi""",Reception,35.1,14.2,,6,,36` + "\r\n" +
		"\r\n" +
		"D4:68:4D:10:21:A1,Cafe,,,,,40,,80,\n" +
		"not a mac,ignored\n" +
		"D4:68:4D:10:22:A2\n"
	rows, skipped := parseNaming(text)
	if len(rows) != 3 || skipped != 2 {
		t.Fatalf("rows = %d, skipped = %d; want 3 and 2 (the header and the note)", len(rows), skipped)
	}
	want := NameRow{ID: "D4:68:4D:10:20:A0", Name: "Lobby, front", Description: `says "hi"`, Location: "Reception",
		Latitude: "35.1", Longitude: "14.2", Width24: "20", Channel24: "6", Width50: "20", Channel50: "36"}
	if rows[0] != want {
		t.Errorf("row 0 = %+v\n     want %+v", rows[0], want)
	}
	// A width with no channel is kept as given, not defaulted.
	if rows[1].Width24 != "40" || rows[1].Channel24 != "" || rows[1].Width50 != "80" {
		t.Errorf("row 1 = %+v", rows[1])
	}
}

func TestParseMACs(t *testing.T) {
	macs, skipped := parseIDs("MAC\nd4:68:4d:10:20:a0,x\nD4:68:4D:10:20:A0\n\n\"D4:68:4D:10:21:A1\"\nfoo\n")
	if !reflect.DeepEqual(macs, []string{"D4:68:4D:10:20:A0", "D4:68:4D:10:21:A1"}) || skipped != 2 {
		t.Errorf("macs = %v, skipped = %d", macs, skipped)
	}
}

func TestPatchCarriesOnlyWhatTheRowSets(t *testing.T) {
	p, err := NameRow{ID: "x", Name: `He said "hi"`, Latitude: "35.5", Longitude: "-14", Width24: "20", Channel24: "6"}.patch()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"name": `He said "hi"`, "latitude": 35.5, "longitude": -14.0,
		"wifi24": map[string]any{"channelWidth": 20, "channel": 6}}
	if !reflect.DeepEqual(p, want) {
		t.Errorf("patch = %v\n   want %v", p, want)
	}
	for _, bad := range []NameRow{{ID: "x"}, {ID: "x", Latitude: "35"}, {ID: "x", Channel50: "auto"}, {ID: "x", Latitude: "n", Longitude: "e"}} {
		if _, err := bad.patch(); err == nil {
			t.Errorf("%+v was accepted", bad)
		}
	}
}

func TestNamingOnSmartZone(t *testing.T) {
	f := fakesz.New()
	f.FailMAC = "D4:68:4D:10:22:A2"
	b := szBackendFor(t, f)
	rows := []NameRow{
		{ID: "D4:68:4D:10:20:A0", Name: "Lobby", Description: "Front desk", Location: "Ground", Channel24: "11", Width24: "20"},
		{ID: "D4:68:4D:10:21:A1", Description: "Only this"},
		{ID: "D4:68:4D:10:22:A2", Name: "Refused"},
		{ID: "D4:68:4D:10:23:A3"},
	}
	emit, got := collect()
	if err := runNaming(context.Background(), b, rows, 4, emit); err != nil {
		t.Fatal(err)
	}
	res := got()
	if res[0].OK != "OK" || res[1].OK != "OK" {
		t.Errorf("rows 0 and 1 = %+v, %+v", res[0], res[1])
	}
	if res[2].Error != "Validation failed for D4:68:4D:10:22:A2" {
		t.Errorf("row 2 error = %q, want the controller's message", res[2].Error)
	}
	if !strings.Contains(res[3].Error, "nothing to change") {
		t.Errorf("row 3 error = %q", res[3].Error)
	}
	ap, _ := f.AP("D4:68:4D:10:20:A0")
	if ap.Name != "Lobby" || ap.Description != "Front desk" || ap.Location != "Ground" || ap.Ch24 != 11 {
		t.Errorf("AP after the run = %+v", ap)
	}
	// A row that sets only the description must not blank the name.
	if ap, _ := f.AP("D4:68:4D:10:21:A1"); ap.Name != "AP-002" || ap.Description != "Only this" {
		t.Errorf("AP 2 after the run = %+v", ap)
	}
	if n := len(f.Calls()); n != 2 {
		t.Errorf("%d changes reached the controller, want 2", n)
	}
}

func TestOverridesOnSmartZone(t *testing.T) {
	f := fakesz.New()
	b := szBackendFor(t, f)
	job := OverrideJob{IDs: []string{"D4:68:4D:10:20:A0", "AA:BB:CC:DD:EE:FF"},
		Overrides: []string{"wifi24", "syslog"}, Radio24: "off", Radio50: "skip", IPDynamic: true}
	emit, got := collect()
	if err := runOverrides(context.Background(), b, job, 2, emit); err != nil {
		t.Fatal(err)
	}
	res := got()
	if res[0].OK != "wifi24, syslog, 2G, IP" || res[0].Error != "" {
		t.Errorf("row 0 = %+v", res[0])
	}
	if res[1].OK != "" || !strings.HasPrefix(res[1].Error, "wifi24 (AP AA:BB:CC:DD:EE:FF not found)") || strings.Contains(res[1].Error, "syslog") {
		t.Errorf("row 1 = %+v; an unknown AP should fail once, not once per override", res[1])
	}
	var paths []string
	for _, c := range f.Calls() {
		paths = append(paths, c.Method+" "+c.Path)
	}
	want := []string{
		"DELETE /aps/D4:68:4D:10:20:A0/wifi24", "DELETE /aps/D4:68:4D:10:20:A0/syslog",
		"PATCH /aps/D4:68:4D:10:20:A0", "PATCH /aps/D4:68:4D:10:20:A0/network",
	}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("requests = %v\n      want %v", paths, want)
	}
	if body := f.Calls()[2].Body; body["wlanService24Enabled"] != false {
		t.Errorf("radio body = %v", body)
	}
}

func TestExportRoundTripsThroughNaming(t *testing.T) {
	f := fakesz.New()
	b := szBackendFor(t, f)
	records, err := runExport(context.Background(), b, "z-staging", 4, func(Event) {})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 13 || !reflect.DeepEqual(records[0], exportHeader) {
		t.Fatalf("%d records, header %v", len(records), records[0])
	}
	var sb strings.Builder
	for _, r := range records {
		sb.WriteString(strings.Join(r, ",") + "\n")
	}
	rows, skipped := parseNaming(sb.String())
	if len(rows) != 12 || skipped != 1 {
		t.Fatalf("an export read back as %d rows with %d skipped; want 12 and the header", len(rows), skipped)
	}
	if rows[0].ID != records[1][0] || rows[0].Name != records[1][1] || rows[0].Channel50 == "" || rows[0].Width50 != "20" {
		t.Errorf("first row back = %+v", rows[0])
	}
}

func TestChannelsCountOnlyConnectedAPs(t *testing.T) {
	for name, b := range map[string]Backend{"sz": szBackendFor(t, fakesz.New()), "r1": r1BackendFor(t, fakesz.NewR1())} {
		zones, err := b.Zones(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		var zone Zone
		for _, z := range zones {
			if z.Name == "HQ Campus" {
				zone = z
			}
		}
		cc, err := runChannels(context.Background(), b, zone.ID, zone.Name, 4, func(Event) {})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		sum := func(rows []ChanRow) (n int) {
			for _, r := range rows {
				n += r.APs
			}
			return n
		}
		if cc.Connected == 0 || cc.Connected >= cc.Total {
			t.Errorf("%s: %d of %d connected; the estate has both kinds", name, cc.Connected, cc.Total)
		}
		if sum(cc.Band24) != cc.Connected || sum(cc.Band50) != cc.Connected {
			t.Errorf("%s: counted %d and %d APs for %d connected", name, sum(cc.Band24), sum(cc.Band50), cc.Connected)
		}
		for i := 1; i < len(cc.Band50); i++ {
			if cc.Band50[i-1].Channel >= cc.Band50[i].Channel {
				t.Errorf("%s: channels are not in order: %v", name, cc.Band50)
			}
		}
	}
}

func TestAnEmptyZoneIsAnError(t *testing.T) {
	b := szBackendFor(t, fakesz.New())
	if _, err := runExport(context.Background(), b, "z-empty", 4, func(Event) {}); err == nil {
		t.Error("exporting an empty zone succeeded")
	}
}

func TestALostSessionStopsTheRun(t *testing.T) {
	f := fakesz.New()
	b := szBackendFor(t, f)
	f.EndSessions()
	rows := make([]NameRow, 20)
	for i := range rows {
		rows[i] = NameRow{ID: "D4:68:4D:10:20:A0", Name: "x"}
	}
	emit, got := collect()
	err := runNaming(context.Background(), b, rows, 2, emit)
	if !errors.Is(err, sz.ErrSessionEnded) {
		t.Fatalf("err = %v, want ErrSessionEnded", err)
	}
	if n := len(got()); n > 4 {
		t.Errorf("%d APs were still attempted after the session ended", n)
	}
}

func TestStopLeavesUnreachedRowsUnmarked(t *testing.T) {
	b := szBackendFor(t, fakesz.New())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	emit, got := collect()
	if err := runNaming(ctx, b, []NameRow{{ID: "D4:68:4D:10:20:A0", Name: "x"}}, 1, emit); err != nil {
		t.Fatal(err)
	}
	if n := len(got()); n != 0 {
		t.Errorf("%d rows were marked after Stop; a row that was never sent is not a failure", n)
	}
}

func TestNamingOnRuckusOne(t *testing.T) {
	f := fakesz.NewR1()
	f.FailSerial = "432100001003"
	b := r1BackendFor(t, f)
	rows := []NameRow{
		{ID: "C8:03:F5:30:40:B0", Name: "Lobby", Location: "ignored", Channel24: "6", Width24: "20", Channel50: "36", Width50: "80"},
		{ID: "C8:03:F5:30:41:B1", Description: "Only this"},
		{ID: "C8:03:F5:30:42:B2", Location: "Nowhere to put this"},
		{ID: "C8:03:F5:30:43:B3", Name: "Refused"},
		{ID: "AA:BB:CC:DD:EE:FF", Name: "Ghost"},
		{ID: "C8:03:F5:30:44:B4", Name: "Half", Latitude: "35"},
	}
	emit, got := collect()
	if err := runNaming(context.Background(), b, rows, 3, emit); err != nil {
		t.Fatal(err)
	}
	res := got()
	if res[0].OK != "OK" || res[1].OK != "OK" {
		t.Errorf("rows 0 and 1 = %+v, %+v", res[0], res[1])
	}
	for i, want := range map[int]string{2: "no location field", 3: "Validation failed", 4: "no AP with this MAC", 5: "given together"} {
		if !strings.Contains(res[i].Error, want) {
			t.Errorf("row %d error = %q, want it to mention %q", i, res[i].Error, want)
		}
	}
	if name, desc := f.Name("C8:03:F5:30:41:B1"); name != "R1-AP-002" || desc != "Only this" {
		t.Errorf("AP 2 = %q / %q; its name should have survived", name, desc)
	}
	var radio map[string]any
	for _, c := range f.Calls() {
		if c.Path == "/432100001000/radioSettings" {
			radio = c.Body
		}
	}
	if radio == nil {
		t.Fatal("no radio settings were written for the AP given channels")
	}
	if p := radio["apRadioParams50G"].(map[string]any); p["manualChannel"] != float64(36) || p["channelBandwidth"] != "80MHz" {
		t.Errorf("5G = %v", p)
	}
}

func TestOverridesOnRuckusOne(t *testing.T) {
	f := fakesz.NewR1()
	b := r1BackendFor(t, f)
	job := OverrideJob{IDs: []string{"C8:03:F5:30:40:B0"}, Overrides: []string{"radio24", "radio50", "ledSettings"},
		Radio24: "skip", Radio50: "off", IPDynamic: true}
	emit, got := collect()
	if err := runOverrides(context.Background(), b, job, 1, emit); err != nil {
		t.Fatal(err)
	}
	if r := got()[0]; r.OK != "radio24+radio50, ledSettings, 5G, IP" || r.Error != "" {
		t.Errorf("row = %+v", r)
	}
	var paths []string
	for _, c := range f.Calls() {
		paths = append(paths, c.Path)
	}
	want := []string{"/432100001000/radioSettings", "/432100001000/ledSettings", "/432100001000/radioSettings", "/432100001000/networkSettings"}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("writes = %v\n    want %v", paths, want)
	}
	calls := f.Calls()
	if calls[1].Body["useVenueSettings"] != true {
		t.Errorf("LED body = %v", calls[1].Body)
	}
	// Switching 5 GHz off must not undo the hand-back that came before it.
	last := calls[2].Body
	if last["enable50G"] != false || last["apRadioParams24G"].(map[string]any)["useVenueSettings"] != true {
		t.Errorf("radio body after switching 5 GHz off = %v", last)
	}
	if calls[3].Body["ipType"] != "DYNAMIC" {
		t.Errorf("network body = %v", calls[3].Body)
	}
}

func TestChooseOverrides(t *testing.T) {
	all := r1Overrides
	def, _ := chooseOverrides(all, "default")
	every, _ := chooseOverrides(all, "all")
	none, _ := chooseOverrides(all, "none")
	if len(def) == 0 || len(def) >= len(every) || len(every) != len(all) || none != nil {
		t.Errorf("default %d, all %d, none %v", len(def), len(every), none)
	}
	got, err := chooseOverrides(all, "LEDsettings, radio24")
	if err != nil || !reflect.DeepEqual(got, []string{"ledSettings", "radio24"}) {
		t.Errorf("list = %v, %v", got, err)
	}
	if _, err := chooseOverrides(all, "ledSettings,nonsense"); err == nil {
		t.Error("an unknown override was accepted")
	}
}

func TestOriginalCommandLineStillParses(t *testing.T) {
	o, err := parseFlags([]string{"-ip=10.0.0.5", "-usr=bob", "-pwd=s3cret", "-prt=7443", "-debug=yes"})
	if err != nil {
		t.Fatal(err)
	}
	if o.host != "10.0.0.5" || o.user != "bob" || o.pass != "s3cret" || o.port != "7443" || !o.debug {
		t.Errorf("options = %+v", o)
	}
	if o.hasAction() {
		t.Error("connection flags alone should open the console, as they always did")
	}
	for _, alias := range []string{"-ipaddr=h", "-ipaddress=h"} {
		if o, _ := parseFlags([]string{alias}); o.host != "h" {
			t.Errorf("%s was not understood", alias)
		}
	}
}

func TestSerialNumbersAreIDsToo(t *testing.T) {
	rows, skipped := parseNaming("Serial,Name\n432100002000,Pending\n 432100001000 ,Online\n12345,too short\n4321000020001,too long\n")
	if len(rows) != 2 || skipped != 3 || rows[0].ID != "432100002000" || rows[1].ID != "432100001000" {
		t.Fatalf("rows = %+v, skipped = %d", rows, skipped)
	}
	ids, _ := parseIDs("432100002000\nC8:03:F5:30:40:B0\n432100002000\n")
	if !reflect.DeepEqual(ids, []string{"432100002000", "C8:03:F5:30:40:B0"}) {
		t.Errorf("ids = %v", ids)
	}
}

// An AP added to a venue but not yet onboarded has no MAC address in RUCKUS
// One; its serial number is all there is to go by.
func TestRuckusOneWorksAPsBySerialBeforeTheyHaveAMAC(t *testing.T) {
	f := fakesz.NewR1()
	b := r1BackendFor(t, f)
	rows := []NameRow{
		{ID: "432100002000", Name: "Lobby-new", Description: "Not installed yet"},
		{ID: "432100001000", Name: "Online-by-serial"},
		{ID: "C8:03:F5:30:41:B1", Name: "Online-by-MAC"},
		{ID: "432199999999", Name: "Not in the tenant"},
	}
	emit, got := collect()
	if err := runNaming(context.Background(), b, rows, 2, emit); err != nil {
		t.Fatal(err)
	}
	res := got()
	for i := 0; i < 3; i++ {
		if res[i].OK != "OK" {
			t.Errorf("row %d = %+v", i, res[i])
		}
	}
	if !strings.Contains(res[3].Error, "no AP with this serial number") {
		t.Errorf("row 3 error = %q", res[3].Error)
	}
	if name, desc := f.Name("432100002000"); name != "Lobby-new" || desc != "Not installed yet" {
		t.Errorf("pending AP = %q / %q", name, desc)
	}
	if name, _ := f.Name("432100001000"); name != "Online-by-serial" {
		t.Errorf("online AP addressed by serial = %q", name)
	}

	emit, got = collect()
	job := OverrideJob{IDs: []string{"432100002000", "432199999999"}, Overrides: []string{"ledSettings"}, Radio24: "off"}
	if err := runOverrides(context.Background(), b, job, 1, emit); err != nil {
		t.Fatal(err)
	}
	res = got()
	if res[0].OK != "ledSettings, 2G" || res[0].Error != "" {
		t.Errorf("overrides by serial = %+v", res[0])
	}
	if res[1].OK != "" || strings.Count(res[1].Error, "serial number") != 1 {
		t.Errorf("an unknown serial should fail once: %+v", res[1])
	}
}

func TestRuckusOneExportIsKeyedBySerial(t *testing.T) {
	b := r1BackendFor(t, fakesz.NewR1())
	records, err := runExport(context.Background(), b, "v-hq", 4, func(Event) {})
	if err != nil {
		t.Fatal(err)
	}
	if records[0][0] != "Serial" || records[0][10] != "MAC" {
		t.Fatalf("header = %v", records[0])
	}
	if len(records) != 23 {
		t.Fatalf("%d records; want the 20 online APs and the 2 not yet onboarded", len(records)-1)
	}
	pending := records[len(records)-1]
	if pending[0] != "432100002001" || pending[10] != "" || pending[1] != "Pending-2" {
		t.Errorf("pending AP record = %v", pending)
	}
	var sb strings.Builder
	for _, r := range records {
		sb.WriteString(strings.Join(r, ",") + "\n")
	}
	rows, skipped := parseNaming(sb.String())
	if len(rows) != 22 || skipped != 1 {
		t.Errorf("the export read back as %d rows with %d skipped; every AP should survive, pending ones included", len(rows), skipped)
	}
}

func TestSmartZoneRefusesASerialNumber(t *testing.T) {
	f := fakesz.New()
	b := szBackendFor(t, f)
	emit, got := collect()
	if err := runNaming(context.Background(), b, []NameRow{{ID: "981902000400", Name: "x"}}, 1, emit); err != nil {
		t.Fatal(err)
	}
	if e := got()[0].Error; !strings.Contains(e, "identifies APs by MAC address") {
		t.Errorf("error = %q", e)
	}
	emit, got = collect()
	job := OverrideJob{IDs: []string{"981902000400"}, Overrides: []string{"wifi24", "syslog"}, IPDynamic: true}
	if err := runOverrides(context.Background(), b, job, 1, emit); err != nil {
		t.Fatal(err)
	}
	if e := got()[0].Error; strings.Count(e, "MAC address") != 1 {
		t.Errorf("error = %q; want the reason once", e)
	}
	if n := len(f.Calls()); n != 0 {
		t.Errorf("%d requests were sent for a row SmartZone cannot address", n)
	}
}
