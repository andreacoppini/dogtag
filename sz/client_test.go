package sz_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/andreacoppini/dogtag/internal/fakesz"
	"github.com/andreacoppini/dogtag/sz"
)

func dial(t *testing.T, f *fakesz.Server, o sz.Options) *sz.Client {
	t.Helper()
	ts := httptest.NewTLSServer(f)
	t.Cleanup(ts.Close)
	u, _ := url.Parse(ts.URL)
	o.Host, o.Port = u.Hostname(), u.Port()
	c, err := sz.New(o)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestLoginPicksTheNewestAPIVersion(t *testing.T) {
	f := fakesz.New()
	f.Versions = []string{"v9_0", "v10_0", "v9_1"} // v10_0 sorts below v9_x as text
	c := dial(t, f, sz.Options{})
	if err := c.Login(context.Background(), "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if got := c.APIVersion(); got != "v10_0" {
		t.Errorf("API version = %s, want v10_0", got)
	}
	if c.ControllerVersion() == "" {
		t.Error("controller version not recorded")
	}
}

func TestLoginHonoursAPinnedVersion(t *testing.T) {
	c := dial(t, fakesz.New(), sz.Options{APIVersion: "v9_1"})
	if err := c.Login(context.Background(), "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if got := c.APIVersion(); got != "v9_1" {
		t.Errorf("API version = %s, want v9_1", got)
	}
}

func TestLoginRefusedSaysWhy(t *testing.T) {
	c := dial(t, fakesz.New(), sz.Options{})
	err := c.Login(context.Background(), "admin", "wrong")
	var api *sz.APIError
	if !errors.As(err, &api) || api.Message != "Invalid username or password" {
		t.Fatalf("err = %v, want the controller's own message", err)
	}
}

func TestCertificateIsCheckedWhenAsked(t *testing.T) {
	c := dial(t, fakesz.New(), sz.Options{VerifyTLS: true, APIVersion: "v9_0"})
	if err := c.Login(context.Background(), "admin", "secret"); err == nil {
		t.Fatal("a self-signed certificate was accepted with verification on")
	}
}

func TestListingsArePagedToTheEnd(t *testing.T) {
	f := fakesz.New()
	f.PageSize = 7
	c := dial(t, f, sz.Options{})
	ctx := context.Background()
	if err := c.Login(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	zones, err := c.Zones(ctx)
	if err != nil || len(zones) != 4 {
		t.Fatalf("zones = %d, %v; want 4", len(zones), err)
	}
	if zones[0].Name != "Default Zone" {
		t.Errorf("zones are not sorted by name: first is %q", zones[0].Name)
	}
	aps, err := c.ZoneAPs(ctx, "z-hq")
	if err != nil || len(aps) != 30 {
		t.Fatalf("APs = %d, %v; want all 30 across pages of 7", len(aps), err)
	}
}

func TestSummaryValuesBecomeText(t *testing.T) {
	c := dial(t, fakesz.New(), sz.Options{})
	ctx := context.Background()
	if err := c.Login(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	s, err := c.APSummary(ctx, "D4:68:4D:10:20:A0")
	if err != nil {
		t.Fatal(err)
	}
	if s["wifi24Channel"] != "1" || s["latitude"] != "" || s["connectionState"] != "Connect" {
		t.Errorf("summary = %v", s)
	}
}

func TestAnEndedSessionIsRecognised(t *testing.T) {
	f := fakesz.New()
	c := dial(t, f, sz.Options{})
	ctx := context.Background()
	if err := c.Login(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := c.KeepAlive(ctx); err != nil {
		t.Fatalf("keep-alive on a live session: %v", err)
	}
	f.EndSessions()
	if err := c.KeepAlive(ctx); !errors.Is(err, sz.ErrSessionEnded) {
		t.Errorf("keep-alive = %v, want ErrSessionEnded", err)
	}
	if err := c.PatchAP(ctx, "D4:68:4D:10:20:A0", map[string]any{"name": "x"}); !errors.Is(err, sz.ErrSessionEnded) {
		t.Errorf("patch = %v, want ErrSessionEnded", err)
	}
}

func TestLogoutEndsTheSession(t *testing.T) {
	c := dial(t, fakesz.New(), sz.Options{})
	ctx := context.Background()
	if err := c.Login(ctx, "admin", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := c.Logout(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Zones(ctx); !errors.Is(err, sz.ErrSessionEnded) {
		t.Errorf("after logout: %v, want ErrSessionEnded", err)
	}
}

func TestBadAddressesAreRefusedBeforeAnythingIsSent(t *testing.T) {
	for _, o := range []sz.Options{{Host: ""}, {Host: "sz", Port: "https"}, {Host: "sz", Port: "70000"}} {
		if _, err := sz.New(o); err == nil {
			t.Errorf("New(%+v) succeeded", o)
		}
	}
	c, err := sz.New(sz.Options{Host: "https://sz.example.com:7443/"})
	if err != nil || c.Host() != "sz.example.com" {
		t.Errorf("a pasted URL was not understood: %v, %v", c, err)
	}
}
