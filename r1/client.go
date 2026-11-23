// Package r1 is a small client for the RUCKUS One API: the calls DogTag needs
// to name access points, pin their channels, hand per-AP settings back to the
// venue and read a venue's inventory.
//
// The shapes used here are the published ones (docs.ruckus.cloud/api, Wi-Fi
// Services and View Model):
//
//	POST {host}/oauth2/token/{tenantId}                    client credentials -> bearer token
//	POST {host}/venues/query                               venues
//	POST {host}/venues/aps/query                           APs, with status and radio channels
//	GET/PUT {host}/venues/{venueId}/aps/{serial}           name, description, GPS
//	GET/PUT {host}/venues/{venueId}/aps/{serial}/radioSettings
//	PUT {host}/venues/{venueId}/aps/{serial}/networkSettings
//	GET/PUT {host}/venues/{venueId}/aps/{serial}/<x>Settings   each with useVenueSettings
//
// Changes are accepted asynchronously: R1 answers 202 and applies them
// afterwards, so "OK" here means accepted, not yet applied.
package r1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// mediaType is the versioned content type the data calls expect. Version 1 is
// the one whose per-AP settings carry useVenueSettings.
const mediaType = "application/vnd.ruckus.v1+json"

var hosts = map[string]string{
	"us":   "api.ruckus.cloud",
	"eu":   "api.eu.ruckus.cloud",
	"asia": "api.asia.ruckus.cloud",
}

// Region normalises what the operator typed (US, eu, Europe, …) to us, eu or
// asia, or returns "" if it is none of them.
func Region(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "us", "na", "usa", "america", "":
		return "us"
	case "eu", "europe":
		return "eu"
	case "asia", "ap", "apac":
		return "asia"
	}
	return ""
}

// Options describes the tenant to sign in to.
type Options struct {
	Region   string
	Tenant   string
	ClientID string
	Secret   string
	// BaseURL replaces the region's API host. It exists for tests and for
	// trying the tool against a stand-in.
	BaseURL string
	Timeout time.Duration
	Trace   func(string)
}

// Client is an authenticated session with one tenant. It is safe for
// concurrent use.
type Client struct {
	opt   Options
	base  string
	hc    *http.Client
	mu    sync.Mutex
	token string
}

// Venue is a venue ("space" in newer R1 wording).
type Venue struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AP is one access point as the AP query reports it.
type AP struct {
	Serial  string
	MAC     string
	Name    string
	Model   string
	VenueID string
	Status  string
	IP      string
	IPType  string
	Radios  []Radio
}

// Radio is one radio's live state.
type Radio struct {
	Band    string
	Channel int
}

// Detail is the part of an AP's configuration the query does not carry.
type Detail struct {
	Name        string
	Description string
	Latitude    string
	Longitude   string
}

// APIError is an answer from R1 that was not a success.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s (HTTP %d)", e.Message, e.Status)
}

// New signs in. A failure here means the region, tenant or key is wrong.
func New(ctx context.Context, o Options) (*Client, error) {
	if o.Tenant == "" || o.ClientID == "" || o.Secret == "" {
		return nil, errors.New("RUCKUS One needs a tenant ID, client ID and shared secret")
	}
	base := strings.TrimRight(o.BaseURL, "/")
	if base == "" {
		h, ok := hosts[Region(o.Region)]
		if !ok {
			return nil, fmt.Errorf("unknown RUCKUS One region %q: use us, eu or asia", o.Region)
		}
		base = "https://" + h
	}
	if o.Timeout <= 0 {
		o.Timeout = 30 * time.Second
	}
	c := &Client{opt: o, base: base, hc: &http.Client{Timeout: o.Timeout}}
	if err := c.signIn(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// Host is the API host in use.
func (c *Client) Host() string {
	if u, err := url.Parse(c.base); err == nil && u.Host != "" {
		return u.Host
	}
	return c.base
}

func (c *Client) logf(format string, a ...any) {
	if c.opt.Trace != nil {
		c.opt.Trace(fmt.Sprintf(format, a...))
	}
}

func (c *Client) signIn(ctx context.Context) error {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.opt.ClientID},
		"client_secret": {c.opt.Secret},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.base+"/oauth2/token/"+url.PathEscape(c.opt.Tenant), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.logf("POST /oauth2/token/{tenant}")
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("RUCKUS One sign-in: %w", plainNetErr(err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RUCKUS One sign-in failed (HTTP %d): %s", resp.StatusCode, reason(body))
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if json.Unmarshal(body, &tok) != nil || tok.AccessToken == "" {
		return errors.New("RUCKUS One sign-in returned no token")
	}
	c.mu.Lock()
	c.token = tok.AccessToken
	c.mu.Unlock()
	return nil
}

// do sends one request and decodes a 2xx answer into out. It signs in again
// once if the token has expired, and backs off when the tenant's rate limit
// is hit: a bulk run over a large venue is exactly what runs into it.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var sent []byte
	if body != nil {
		var err error
		if sent, err = json.Marshal(body); err != nil {
			return err
		}
	}
	reauthed := false
	for attempt := 0; ; attempt++ {
		var rd io.Reader
		if sent != nil {
			rd = bytes.NewReader(sent)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.base+path, rd)
		if err != nil {
			return err
		}
		c.mu.Lock()
		req.Header.Set("Authorization", "Bearer "+c.token)
		c.mu.Unlock()
		req.Header.Set("Accept", mediaType)
		if sent != nil {
			req.Header.Set("Content-Type", mediaType)
			c.logf("%s %s %s", method, path, clip(string(sent), 400))
		} else {
			c.logf("%s %s", method, path)
		}
		resp, err := c.hc.Do(req)
		if err != nil {
			return plainNetErr(err)
		}
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
		resp.Body.Close()
		c.logf("  -> %d %s", resp.StatusCode, clip(strings.TrimSpace(string(raw)), 300))

		switch {
		case resp.StatusCode == http.StatusUnauthorized && !reauthed:
			reauthed = true
			if err := c.signIn(ctx); err != nil {
				return err
			}
			continue
		case resp.StatusCode == http.StatusTooManyRequests && attempt < 6:
			wait := time.Duration(1<<min(attempt, 4)) * time.Second
			if s, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && s > 0 && s <= 60 {
				wait = time.Duration(s) * time.Second
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
			continue
		case resp.StatusCode < 200 || resp.StatusCode > 299:
			return &APIError{Status: resp.StatusCode, Message: reason(raw)}
		}
		if out != nil && len(bytes.TrimSpace(raw)) > 0 {
			if err := json.Unmarshal(raw, out); err != nil {
				return fmt.Errorf("RUCKUS One sent something unreadable: %s", clip(string(raw), 120))
			}
		}
		return nil
	}
}

// Venues lists the tenant's venues.
func (c *Client) Venues(ctx context.Context) ([]Venue, error) {
	var out []Venue
	err := query(ctx, c, "/venues/query", []string{"id", "name"}, nil, func(raw json.RawMessage) error {
		var v []Venue
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v...)
		return nil
	})
	return out, err
}

// APs lists access points: one venue's, or the whole tenant's when venueID is
// empty.
func (c *Client) APs(ctx context.Context, venueID string) ([]AP, error) {
	var filters map[string][]string
	if venueID != "" {
		filters = map[string][]string{"venueId": {venueID}}
	}
	fields := []string{"serialNumber", "macAddress", "name", "model", "venueId", "status", "networkStatus", "radioStatuses"}
	var out []AP
	err := query(ctx, c, "/venues/aps/query", fields, filters, func(raw json.RawMessage) error {
		var page []struct {
			SerialNumber  string `json:"serialNumber"`
			MACAddress    string `json:"macAddress"`
			Name          string `json:"name"`
			Model         string `json:"model"`
			VenueID       string `json:"venueId"`
			Status        string `json:"status"`
			NetworkStatus struct {
				IPAddress     string `json:"ipAddress"`
				IPAddressType string `json:"ipAddressType"`
			} `json:"networkStatus"`
			RadioStatuses []struct {
				Band    string `json:"band"`
				Channel int    `json:"channel"`
			} `json:"radioStatuses"`
		}
		if err := json.Unmarshal(raw, &page); err != nil {
			return err
		}
		for _, d := range page {
			// The filter is the server's to honour; checking again here means
			// a filter it ignored cannot leak another venue's APs into a run.
			if venueID != "" && d.VenueID != venueID {
				continue
			}
			ap := AP{Serial: d.SerialNumber, MAC: strings.ToUpper(d.MACAddress), Name: d.Name, Model: d.Model,
				VenueID: d.VenueID, Status: d.Status, IP: d.NetworkStatus.IPAddress, IPType: d.NetworkStatus.IPAddressType}
			for _, r := range d.RadioStatuses {
				ap.Radios = append(ap.Radios, Radio{Band: r.Band, Channel: r.Channel})
			}
			out = append(out, ap)
		}
		return nil
	})
	return out, err
}

// query walks a view-model query to the end.
func query(ctx context.Context, c *Client, path string, fields []string, filters map[string][]string, each func(json.RawMessage) error) error {
	const pageSize = 1000
	seen := 0
	for page := 1; ; page++ {
		req := map[string]any{"fields": fields, "page": page, "pageSize": pageSize}
		if filters != nil {
			req["filters"] = filters
		}
		var resp struct {
			TotalCount int               `json:"totalCount"`
			Data       []json.RawMessage `json:"data"`
		}
		if err := c.do(ctx, http.MethodPost, path, req, &resp); err != nil {
			return err
		}
		raw, _ := json.Marshal(resp.Data)
		if err := each(raw); err != nil {
			return fmt.Errorf("RUCKUS One sent something unreadable: %w", err)
		}
		seen += len(resp.Data)
		// Driven by the reported total, not page fullness: a short page is not
		// necessarily the last. An empty page is the backstop for a bad total.
		if len(resp.Data) == 0 || seen >= resp.TotalCount {
			return nil
		}
	}
}

func apPath(venueID, serial string) string {
	return "/venues/" + url.PathEscape(venueID) + "/aps/" + url.PathEscape(serial)
}

// Detail reads an AP's name, description and GPS position.
func (c *Client) Detail(ctx context.Context, venueID, serial string) (Detail, error) {
	raw, err := c.getAP(ctx, venueID, serial)
	if err != nil {
		return Detail{}, err
	}
	d := Detail{Name: str(raw["name"]), Description: str(raw["description"])}
	if gps, ok := raw["deviceGps"].(map[string]any); ok {
		d.Latitude, d.Longitude = str(gps["latitude"]), str(gps["longitude"])
	}
	return d, nil
}

func (c *Client) getAP(ctx context.Context, venueID, serial string) (map[string]any, error) {
	var raw map[string]any
	err := c.do(ctx, http.MethodGet, apPath(venueID, serial), nil, &raw)
	return raw, err
}

// UpdateAP changes an AP's name, description and GPS position. Empty values
// are left as they are. The API replaces the whole record and insists on a
// name, so the current one is read first and only the given fields change.
func (c *Client) UpdateAP(ctx context.Context, venueID, serial string, d Detail) error {
	cur, err := c.getAP(ctx, venueID, serial)
	if err != nil {
		return err
	}
	body := map[string]any{}
	for _, k := range []string{"name", "description", "deviceGps", "tags"} {
		if v, ok := cur[k]; ok && v != nil {
			body[k] = v
		}
	}
	if d.Name != "" {
		body["name"] = d.Name
	}
	if d.Description != "" {
		body["description"] = d.Description
	}
	if d.Latitude != "" && d.Longitude != "" {
		body["deviceGps"] = map[string]string{"latitude": d.Latitude, "longitude": d.Longitude}
	}
	return c.do(ctx, http.MethodPut, apPath(venueID, serial), body, nil)
}

// Band names the per-band blocks of an AP's radio settings.
var bandKey = map[string]string{"24": "apRadioParams24G", "50": "apRadioParams50G", "6": "apRadioParams6G"}
var enableKey = map[string]string{"24": "enable24G", "50": "enable50G", "6": "enable6G"}

// Widths maps a channel width in MHz to the API's name for it.
func widthName(mhz int) (string, error) {
	switch mhz {
	case 0:
		return "AUTO", nil
	case 20, 40, 80, 160:
		return strconv.Itoa(mhz) + "MHz", nil
	}
	return "", fmt.Errorf("channel width %d is not one RUCKUS One offers (20, 40, 80 or 160)", mhz)
}

// RadioChange is what to do to one band. A nil field is left alone.
type RadioChange struct {
	Channel  *int  // pin this channel
	Width    *int  // MHz; 0 means automatic
	Enabled  *bool // switch the radio on or off
	UseVenue bool  // hand the band's settings back to the venue
}

// ChangeRadios edits an AP's radio settings: the current settings are read,
// the named bands ("24", "50", "6") are changed, and the result is written
// back, so everything not mentioned stays as it was.
func (c *Client) ChangeRadios(ctx context.Context, venueID, serial string, bands map[string]RadioChange) error {
	path := apPath(venueID, serial) + "/radioSettings"
	var cur map[string]any
	if err := c.do(ctx, http.MethodGet, path, nil, &cur); err != nil {
		return err
	}
	if cur == nil {
		cur = map[string]any{}
	}
	for band, ch := range bands {
		key, ok := bandKey[band]
		if !ok {
			return fmt.Errorf("unknown band %q", band)
		}
		params, _ := cur[key].(map[string]any)
		if params == nil {
			if ch.UseVenue && ch.Channel == nil && ch.Enabled == nil {
				continue // the AP has no such radio; nothing to hand back
			}
			params = map[string]any{}
		}
		if ch.UseVenue {
			params["useVenueSettings"] = true
		}
		if ch.Channel != nil {
			params["method"] = "MANUAL"
			params["manualChannel"] = *ch.Channel
			params["useVenueSettings"] = false
		}
		if ch.Width != nil {
			w, err := widthName(*ch.Width)
			if err != nil {
				return err
			}
			params["channelBandwidth"] = w
			params["useVenueSettings"] = false
		}
		cur[key] = params
		if ch.Enabled != nil {
			cur[enableKey[band]] = *ch.Enabled
		}
	}
	return c.do(ctx, http.MethodPut, path, cur, nil)
}

// UseVenueSettings hands one per-AP settings resource ("ledSettings",
// "smartMonitorSettings", …) back to the venue.
func (c *Client) UseVenueSettings(ctx context.Context, venueID, serial, resource string) error {
	path := apPath(venueID, serial) + "/" + url.PathEscape(resource)
	var cur map[string]any
	if err := c.do(ctx, http.MethodGet, path, nil, &cur); err != nil {
		return err
	}
	if cur == nil {
		cur = map[string]any{}
	}
	cur["useVenueSettings"] = true
	return c.do(ctx, http.MethodPut, path, cur, nil)
}

// SetDynamicIP switches an AP to DHCP.
func (c *Client) SetDynamicIP(ctx context.Context, venueID, serial string) error {
	return c.do(ctx, http.MethodPut, apPath(venueID, serial)+"/networkSettings", map[string]any{"ipType": "DYNAMIC"}, nil)
}

// reason pulls a short message out of an R1 error body.
func reason(body []byte) string {
	var e struct {
		Errors []struct {
			Message string `json:"message"`
			Reason  string `json:"reason"`
		} `json:"errors"`
		Message string `json:"message"`
		Error   string `json:"error_description"`
	}
	if json.Unmarshal(body, &e) == nil {
		if len(e.Errors) > 0 {
			if m := firstNonEmpty(e.Errors[0].Message, e.Errors[0].Reason); m != "" {
				return m
			}
		}
		if m := firstNonEmpty(e.Message, e.Error); m != "" {
			return m
		}
	}
	if s := clip(strings.TrimSpace(string(body)), 200); s != "" {
		return s
	}
	return "no detail"
}

func plainNetErr(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		if ue.Timeout() {
			return errors.New("RUCKUS One did not answer in time")
		}
		return ue.Err
	}
	return err
}

func str(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	}
	return ""
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
