package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func consoleFor(opt options) (*uiServer, http.Handler) {
	u := &uiServer{opt: opt, subs: map[chan Event]struct{}{}, addr: "127.0.0.1:9999"}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/defaults", u.handleDefaults)
	mux.HandleFunc("/api/run", u.handleRun)
	mux.HandleFunc("/api/parse", u.handleParse)
	return u, u.guard(mux)
}

func TestConsoleRefusesOtherPages(t *testing.T) {
	_, h := consoleFor(options{})
	cases := []struct {
		host, origin string
		want         int
	}{
		{"127.0.0.1:9999", "", 200},
		{"localhost:9999", "http://localhost:9999", 200},
		{"127.0.0.1:9999", "http://127.0.0.1:9999", 200},
		{"127.0.0.1:9999", "https://evil.example", 403}, // another site posting to the console
		{"evil.example:9999", "", 403},                  // DNS rebinding
		{"127.0.0.1:1234", "", 403},
	}
	for _, c := range cases {
		r := httptest.NewRequest("GET", "/api/defaults", nil)
		r.Host = c.host
		if c.origin != "" {
			r.Header.Set("Origin", c.origin)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != c.want {
			t.Errorf("host %q origin %q: HTTP %d, want %d", c.host, c.origin, w.Code, c.want)
		}
	}
}

func TestSecretsFromTheCommandLineStayOutOfThePage(t *testing.T) {
	_, h := consoleFor(options{host: "10.0.0.5", user: "admin", pass: "hunter2", r1Secret: "s3cret", timeout: time.Second})
	r := httptest.NewRequest("GET", "/api/defaults", nil)
	r.Host = "127.0.0.1:9999"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	body := w.Body.String()
	if strings.Contains(body, "hunter2") || strings.Contains(body, "s3cret") {
		t.Fatalf("a secret reached the browser: %s", body)
	}
	var d map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &d)
	if d["hasPass"] != true || d["hasR1Secret"] != true || d["host"] != "10.0.0.5" {
		t.Errorf("defaults = %v", d)
	}
}

func TestRunNeedsAConnection(t *testing.T) {
	_, h := consoleFor(options{})
	r := httptest.NewRequest("POST", "/api/run", strings.NewReader(`{"job":"naming","rows":[{"id":"D4:68:4D:10:20:A0","name":"x"}]}`))
	r.Host = "127.0.0.1:9999"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 || !strings.Contains(w.Body.String(), "not connected") {
		t.Errorf("HTTP %d %s", w.Code, w.Body.String())
	}
}
