package r1_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreacoppini/dogtag/internal/fakesz"
	"github.com/andreacoppini/dogtag/r1"
)

func dial(t *testing.T, f *fakesz.R1) *r1.Client {
	t.Helper()
	ts := httptest.NewServer(f)
	t.Cleanup(ts.Close)
	c, err := r1.New(context.Background(), r1.Options{BaseURL: ts.URL, Tenant: f.Tenant, ClientID: f.ClientID, Secret: f.Secret})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSignInFailureNamesTheReason(t *testing.T) {
	f := fakesz.NewR1()
	ts := httptest.NewServer(f)
	defer ts.Close()
	_, err := r1.New(context.Background(), r1.Options{BaseURL: ts.URL, Tenant: f.Tenant, ClientID: f.ClientID, Secret: "nope"})
	if err == nil || !strings.Contains(err.Error(), "Bad client credentials") {
		t.Fatalf("err = %v", err)
	}
	if _, err := r1.New(context.Background(), r1.Options{Region: "mars", Tenant: "t", ClientID: "c", Secret: "s"}); err == nil {
		t.Error("an unknown region was accepted")
	}
}

func TestRegionSpellings(t *testing.T) {
	for in, want := range map[string]string{"": "us", "US": "us", "Europe": "eu", "apac": "asia", "moon": ""} {
		if got := r1.Region(in); got != want {
			t.Errorf("Region(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestQueriesArePagedAndFilteredByVenue(t *testing.T) {
	f := fakesz.NewR1()
	f.PageSize = 6
	c := dial(t, f)
	ctx := context.Background()
	all, err := c.APs(ctx, "")
	if err != nil || len(all) != 30 {
		t.Fatalf("tenant APs = %d, %v; want 30 across pages of 6", len(all), err)
	}
	hq, err := c.APs(ctx, "v-hq")
	if err != nil || len(hq) != 22 {
		t.Fatalf("venue APs = %d, %v; want 22", len(hq), err)
	}
	ap := hq[0]
	if ap.MAC != "C8:03:F5:30:40:B0" || ap.Serial == "" || len(ap.Radios) != 2 || ap.IP == "" {
		t.Errorf("AP = %+v", ap)
	}
	venues, err := c.Venues(ctx)
	if err != nil || len(venues) != 3 {
		t.Fatalf("venues = %d, %v", len(venues), err)
	}
}

func TestUpdateKeepsWhatItWasNotAskedToChange(t *testing.T) {
	f := fakesz.NewR1()
	c := dial(t, f)
	ctx := context.Background()
	aps, _ := c.APs(ctx, "v-hq")
	ap := aps[0]

	// The API replaces the record and requires a name, so a description-only
	// change has to carry the existing name with it.
	if err := c.UpdateAP(ctx, ap.VenueID, ap.Serial, r1.Detail{Description: "Lobby"}); err != nil {
		t.Fatal(err)
	}
	if name, desc := f.Name(ap.MAC); name != "R1-AP-001" || desc != "Lobby" {
		t.Errorf("after description change: name %q, description %q", name, desc)
	}
	if err := c.UpdateAP(ctx, ap.VenueID, ap.Serial, r1.Detail{Name: "Front", Latitude: "35.9", Longitude: "14.5"}); err != nil {
		t.Fatal(err)
	}
	d, err := c.Detail(ctx, ap.VenueID, ap.Serial)
	if err != nil || d.Name != "Front" || d.Description != "Lobby" || d.Latitude != "35.9" {
		t.Errorf("detail = %+v, %v", d, err)
	}
}

func TestRadioChangesLeaveTheOtherBandAlone(t *testing.T) {
	f := fakesz.NewR1()
	c := dial(t, f)
	ctx := context.Background()
	aps, _ := c.APs(ctx, "v-hq")
	ap := aps[0]
	ch, w, off := 11, 20, false
	err := c.ChangeRadios(ctx, ap.VenueID, ap.Serial, map[string]r1.RadioChange{"24": {Channel: &ch, Width: &w, Enabled: &off}})
	if err != nil {
		t.Fatal(err)
	}
	calls := f.Calls()
	body := calls[len(calls)-1].Body
	p24 := body["apRadioParams24G"].(map[string]any)
	if p24["method"] != "MANUAL" || p24["manualChannel"] != float64(11) || p24["channelBandwidth"] != "20MHz" || p24["useVenueSettings"] != false {
		t.Errorf("2.4G = %v", p24)
	}
	if body["enable24G"] != false || body["enable50G"] != true {
		t.Errorf("enables = %v / %v", body["enable24G"], body["enable50G"])
	}
	if p50 := body["apRadioParams50G"].(map[string]any); p50["method"] != "CHANNELFLY" {
		t.Errorf("5G was disturbed: %v", p50)
	}

	// Handing a band back, including one the AP does not have.
	err = c.ChangeRadios(ctx, ap.VenueID, ap.Serial, map[string]r1.RadioChange{"24": {UseVenue: true}, "6": {UseVenue: true}})
	if err != nil {
		t.Fatal(err)
	}
	calls = f.Calls()
	body = calls[len(calls)-1].Body
	if body["apRadioParams24G"].(map[string]any)["useVenueSettings"] != true {
		t.Error("2.4G was not handed back to the venue")
	}
	if _, has := body["apRadioParams6G"]; has {
		t.Error("a 6 GHz block was invented for an AP without that radio")
	}

	bad := 30
	if err := c.ChangeRadios(ctx, ap.VenueID, ap.Serial, map[string]r1.RadioChange{"50": {Width: &bad}}); err == nil {
		t.Error("a 30 MHz channel width was accepted")
	}
}

func TestRateLimitIsWaitedOut(t *testing.T) {
	f := fakesz.NewR1()
	c := dial(t, f)
	f.Throttle = 1
	if _, err := c.Venues(context.Background()); err != nil {
		t.Fatalf("a 429 was not retried: %v", err)
	}
}
