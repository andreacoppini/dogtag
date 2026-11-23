// Package sz is a small client for the RUCKUS SmartZone public API: just the
// calls DogTag needs to name access points, pin their channels, clear
// per-AP overrides and read a zone's inventory.
package sz

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FallbackAPIVersion is used when the controller does not say which API
// versions it offers. It is the version the original DogTag was written
// against (SmartZone 5.0).
const FallbackAPIVersion = "v7_0"

// ErrSessionEnded means the controller no longer recognises the session: it
// timed out, or an administrator ended it.
var ErrSessionEnded = errors.New("the SmartZone session has ended")

// Options describes how to reach a controller.
type Options struct {
	Host string
	Port string
	// APIVersion pins the API version ("v9_1"). Empty means ask the controller
	// and use the newest it offers.
	APIVersion string
	// VerifyTLS turns certificate verification on. It is off by default
	// because a SmartZone almost always presents a self-signed certificate.
	VerifyTLS bool
	Timeout   time.Duration
	// Trace, when set, is called with one line per request and response.
	Trace func(string)
}

// Client is one signed-in session with a controller. It is safe for
// concurrent use.
type Client struct {
	base    string // https://host:port/wsg/api/public
	host    string
	http    *http.Client
	trace   func(string)
	pinned  string
	mu      sync.Mutex
	apiVer  string
	version string
	lastUse time.Time
}

// Zone is one AP zone.
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AP is one access point as the zone listing reports it.
type AP struct {
	MAC  string `json:"mac"`
	Name string `json:"name"`
}

// Summary is an AP's operational summary. Every value is kept as text,
// because that is how it is shown and exported; a field the controller left
// out or sent as null is an empty string.
type Summary map[string]string

// APIError is an answer from the controller that was not a success.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("HTTP %d", e.Status)
}

// New prepares a client. Nothing is sent until Login.
func New(o Options) (*Client, error) {
	host := strings.TrimSpace(o.Host)
	host = strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://")
	host = strings.TrimRight(host, "/")
	if host == "" {
		return nil, errors.New("no SmartZone address given")
	}
	port := strings.TrimSpace(o.Port)
	if port == "" {
		port = "8443"
	}
	if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
		return nil, fmt.Errorf("%q is not a port number", port)
	}
	// A pasted "host:port" wins over the separate port field only when it
	// really is one, so a bare IPv6 address is left alone.
	if h, p, err := net.SplitHostPort(host); err == nil {
		host, port = h, p
	}
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Second
	}
	jar, _ := cookiejar.New(nil)
	c := &Client{
		base:   "https://" + net.JoinHostPort(host, port) + "/wsg/api/public",
		host:   host,
		trace:  o.Trace,
		pinned: strings.TrimSpace(o.APIVersion),
		http: &http.Client{
			Jar:     jar,
			Timeout: o.Timeout,
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: !o.VerifyTLS}, //nolint:gosec // opt-in verification; controllers ship self-signed
				MaxIdleConnsPerHost: 32,
				IdleConnTimeout:     60 * time.Second,
			},
		},
	}
	return c, nil
}

// Host is the controller address as given.
func (c *Client) Host() string { return c.host }

// APIVersion is the API version in use, once signed in.
func (c *Client) APIVersion() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.apiVer
}

// ControllerVersion is the SmartZone software version reported at sign-in.
func (c *Client) ControllerVersion() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.version
}

// IdleFor is how long ago the last request was made, which is what the
// keep-alive goes by.
func (c *Client) IdleFor() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Since(c.lastUse)
}

func (c *Client) logf(format string, a ...any) {
	if c.trace != nil {
		c.trace(fmt.Sprintf(format, a...))
	}
}

// do sends one request. path is relative to the versioned API root unless it
// starts with "//", which addresses the unversioned root. out, when not nil,
// receives the decoded JSON body.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var target string
	if strings.HasPrefix(path, "//") {
		target = c.base + path[1:]
	} else {
		target = c.base + "/" + c.APIVersion() + path
	}

	var rd io.Reader
	var sent []byte
	if body != nil {
		var err error
		sent, err = json.Marshal(body)
		if err != nil {
			return err
		}
		rd = bytes.NewReader(sent)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, rd)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	c.mu.Lock()
	c.lastUse = time.Now()
	c.mu.Unlock()

	if sent != nil && !strings.HasSuffix(path, "/session") {
		c.logf("%s %s %s", method, path, sent)
	} else {
		c.logf("%s %s", method, path) // never the sign-in body: it holds the password
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return plainNetErr(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return err
	}
	c.logf("  -> %d %s", resp.StatusCode, clip(strings.TrimSpace(string(raw)), 300))

	if resp.StatusCode == http.StatusUnauthorized {
		if strings.HasSuffix(path, "/session") && method == http.MethodPost {
			return &APIError{Status: resp.StatusCode, Message: firstNonEmpty(errorMessage(raw), "the controller refused that username and password")}
		}
		return ErrSessionEnded
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &APIError{Status: resp.StatusCode, Message: firstNonEmpty(errorMessage(raw), fmt.Sprintf("HTTP %d", resp.StatusCode))}
	}
	if out != nil && len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("the controller sent something that is not JSON: %s", clip(string(raw), 120))
		}
	}
	return nil
}

// errorMessage pulls the readable part out of a SmartZone error body.
func errorMessage(raw []byte) string {
	var e struct {
		Message   string `json:"message"`
		ErrorType string `json:"errorType"`
	}
	if json.Unmarshal(raw, &e) != nil {
		return clip(strings.TrimSpace(string(raw)), 200)
	}
	return firstNonEmpty(e.Message, e.ErrorType)
}

// plainNetErr strips the URL and method Go wraps around a transport error;
// what is left names the actual problem.
func plainNetErr(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Timeout() {
			return errors.New("the controller did not answer in time")
		}
		return ue.Err
	}
	return err
}

// Login signs in and settles which API version to speak.
func (c *Client) Login(ctx context.Context, user, pass string) error {
	ver := c.pinned
	if ver == "" {
		ver = c.discoverVersion(ctx)
	}
	c.mu.Lock()
	c.apiVer = ver
	c.mu.Unlock()

	var out struct {
		ControllerVersion string `json:"controllerVersion"`
	}
	err := c.do(ctx, http.MethodPost, "/session", map[string]string{"username": user, "password": pass}, &out)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.version = out.ControllerVersion
	c.mu.Unlock()
	return nil
}

// discoverVersion asks the controller which API versions it offers and picks
// the newest. Older SmartZone releases drop old API versions, so a version
// that is hard-coded stops working after an upgrade.
func (c *Client) discoverVersion(ctx context.Context) string {
	var info struct {
		Versions []string `json:"apiSupportVersions"`
	}
	if err := c.do(ctx, http.MethodGet, "//apiInfo", nil, &info); err != nil || len(info.Versions) == 0 {
		return FallbackAPIVersion
	}
	return newestVersion(info.Versions)
}

// newestVersion orders "v9_1"-style names numerically; as text, v10_0 would
// sort below v9_0.
func newestVersion(vs []string) string {
	key := func(v string) [2]int {
		var k [2]int
		parts := strings.SplitN(strings.TrimPrefix(v, "v"), "_", 2)
		for i := 0; i < len(parts) && i < 2; i++ {
			k[i], _ = strconv.Atoi(parts[i])
		}
		return k
	}
	best := vs[0]
	for _, v := range vs[1:] {
		a, b := key(v), key(best)
		if a[0] > b[0] || (a[0] == b[0] && a[1] > b[1]) {
			best = v
		}
	}
	return best
}

// Logout ends the session on the controller.
func (c *Client) Logout(ctx context.Context) error {
	return c.do(ctx, http.MethodDelete, "/session", nil, nil)
}

// KeepAlive touches the session so the controller does not time it out.
func (c *Client) KeepAlive(ctx context.Context) error {
	var out map[string]any
	if err := c.do(ctx, http.MethodGet, "/session", nil, &out); err != nil {
		return err
	}
	// Some releases answer an expired session with 200 and success:false.
	if ok, present := out["success"].(bool); present && !ok {
		return ErrSessionEnded
	}
	return nil
}

type page[T any] struct {
	TotalCount int  `json:"totalCount"`
	HasMore    bool `json:"hasMore"`
	List       []T  `json:"list"`
}

// listAll walks a paged listing to the end. The controller caps a page at
// 1000 entries whatever is asked for, so a single large request silently
// truncates a big zone.
func listAll[T any](ctx context.Context, c *Client, path string) ([]T, error) {
	const size = 1000
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	var all []T
	for {
		// The next page starts where this one ended, not a page-size further
		// on: the controller is free to return fewer than were asked for.
		var p page[T]
		if err := c.do(ctx, http.MethodGet, fmt.Sprintf("%s%sindex=%d&listSize=%d", path, sep, len(all), size), nil, &p); err != nil {
			return nil, err
		}
		all = append(all, p.List...)
		if !p.HasMore || len(p.List) == 0 {
			return all, nil
		}
	}
}

// Zones lists the AP zones, sorted by name.
func (c *Client) Zones(ctx context.Context) ([]Zone, error) {
	zones, err := listAll[Zone](ctx, c, "/rkszones")
	if err != nil {
		return nil, err
	}
	sort.SliceStable(zones, func(i, j int) bool {
		return strings.ToLower(zones[i].Name) < strings.ToLower(zones[j].Name)
	})
	return zones, nil
}

// ZoneAPs lists the access points in a zone.
func (c *Client) ZoneAPs(ctx context.Context, zoneID string) ([]AP, error) {
	return listAll[AP](ctx, c, "/aps?zoneId="+url.QueryEscape(zoneID))
}

// APSummary reads one AP's operational summary.
func (c *Client) APSummary(ctx context.Context, mac string) (Summary, error) {
	var raw map[string]any
	if err := c.do(ctx, http.MethodGet, "/aps/"+url.PathEscape(mac)+"/operational/summary", nil, &raw); err != nil {
		return nil, err
	}
	out := Summary{}
	for k, v := range raw {
		out[k] = text(v)
	}
	return out, nil
}

// PatchAP changes fields of an AP's configuration.
func (c *Client) PatchAP(ctx context.Context, mac string, fields map[string]any) error {
	return c.do(ctx, http.MethodPatch, "/aps/"+url.PathEscape(mac), fields, nil)
}

// PatchAPNetwork changes an AP's network settings.
func (c *Client) PatchAPNetwork(ctx context.Context, mac string, fields map[string]any) error {
	return c.do(ctx, http.MethodPatch, "/aps/"+url.PathEscape(mac)+"/network", fields, nil)
}

// DeleteOverride removes one per-AP override, so that the AP inherits the
// setting from its group or zone again.
func (c *Client) DeleteOverride(ctx context.Context, mac, override string) error {
	return c.do(ctx, http.MethodDelete, "/aps/"+url.PathEscape(mac)+"/"+url.PathEscape(override), nil, nil)
}

// text renders a JSON value the way it should read in a table or a CSV cell.
func text(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
