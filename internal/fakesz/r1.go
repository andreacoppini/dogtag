package fakesz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// R1 is a stand-in RUCKUS One tenant, shaped after the published API: token
// sign-in, the venue and AP queries, and the per-AP resources DogTag edits.
type R1 struct {
	Tenant, ClientID, Secret string
	// FailSerial makes every change to that AP fail.
	FailSerial string
	// PageSize caps a query page, to exercise paging. Zero means 1000.
	PageSize int
	// Throttle answers this many requests with 429 before behaving.
	Throttle int

	mu       sync.Mutex
	venues   []zone
	aps      []*r1AP
	calls    []Call
	settings map[string]map[string]any // "serial/resource" -> body
}

type r1AP struct {
	Serial, MAC, Name, Description, Model, Venue, Status, IP string
	Lat, Lon                                                 string
	Ch24, Ch50                                               int
}

// NewR1 builds a tenant with a small demo estate.
func NewR1() *R1 {
	s := &R1{Tenant: "tenant-1", ClientID: "client-1", Secret: "secret",
		settings: map[string]map[string]any{},
		venues:   []zone{{"v-hq", "HQ Campus"}, {"v-branch", "Branch Office"}, {"v-empty", "Empty Venue"}},
	}
	ch24 := []int{1, 6, 11}
	ch50 := []int{36, 44, 52, 100, 149, 157}
	models := []string{"R650", "R770", "R350", "T350d"}
	for i := 0; i < 28; i++ {
		v := "v-hq"
		if i >= 20 {
			v = "v-branch"
		}
		status := "2_00_Operational"
		if i%9 == 4 {
			status = "3_04_DisconnectedFromCloud"
		}
		s.aps = append(s.aps, &r1AP{
			Serial: fmt.Sprintf("4321%08d", 1000+i),
			MAC:    fmt.Sprintf("C8:03:F5:%02X:%02X:%02X", 0x30+i/16, 0x40+i, 0xB0+i),
			Name:   fmt.Sprintf("R1-AP-%03d", i+1),
			Model:  models[i%len(models)], Venue: v, Status: status,
			IP:   fmt.Sprintf("10.30.%d.%d", 1+i/20, 20+i),
			Ch24: ch24[i%len(ch24)], Ch50: ch50[(i*5)%len(ch50)],
		})
	}
	// Two APs that have been added to a venue but have never come online: R1
	// knows them by serial number alone, with no MAC, address or radios yet.
	for i := 0; i < 2; i++ {
		s.aps = append(s.aps, &r1AP{
			Serial: fmt.Sprintf("4321%08d", 2000+i),
			Name:   fmt.Sprintf("Pending-%d", i+1),
			Model:  "R650", Venue: "v-hq", Status: "1_01_NeverContactedCloud",
		})
	}
	return s
}

// Calls returns the changes made so far.
func (s *R1) Calls() []Call {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Call(nil), s.calls...)
}

// Name returns an AP's current name and description, by MAC or serial.
func (s *R1) Name(id string) (name, description string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.aps {
		if a.Serial == id || (a.MAC != "" && strings.EqualFold(a.MAC, id)) {
			return a.Name, a.Description
		}
	}
	return "", ""
}

func (s *R1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/oauth2/token/"+s.Tenant && r.Method == http.MethodPost {
		_ = r.ParseForm()
		if r.Form.Get("client_id") != s.ClientID || r.Form.Get("client_secret") != s.Secret {
			reply(w, 401, map[string]any{"error": "invalid_client", "error_description": "Bad client credentials"})
			return
		}
		reply(w, 200, map[string]any{"access_token": "token-" + s.Tenant, "token_type": "Bearer", "expires_in": 7200})
		return
	}
	if r.Header.Get("Authorization") != "Bearer token-"+s.Tenant {
		r1Fail(w, 401, "Unauthorized")
		return
	}
	s.mu.Lock()
	if s.Throttle > 0 {
		s.Throttle--
		s.mu.Unlock()
		w.Header().Set("Retry-After", "1")
		r1Fail(w, 429, "Too Many Requests")
		return
	}
	s.mu.Unlock()

	var body map[string]any
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}
	seg := strings.Split(strings.Trim(path, "/"), "/")
	switch {
	case path == "/venues/query" && r.Method == http.MethodPost:
		var list []map[string]any
		for _, v := range s.venues {
			list = append(list, map[string]any{"id": v.ID, "name": v.Name})
		}
		s.page(w, body, list)
	case path == "/venues/aps/query" && r.Method == http.MethodPost:
		want := ""
		if f, ok := body["filters"].(map[string]any); ok {
			if ids, ok := f["venueId"].([]any); ok && len(ids) > 0 {
				want, _ = ids[0].(string)
			}
		}
		list := []map[string]any{}
		s.mu.Lock()
		for _, a := range s.aps {
			if want != "" && a.Venue != want {
				continue
			}
			item := map[string]any{
				"serialNumber": a.Serial, "name": a.Name, "model": a.Model,
				"venueId": a.Venue, "status": a.Status,
			}
			if a.MAC != "" {
				item["macAddress"] = a.MAC
				item["networkStatus"] = map[string]any{"ipAddress": a.IP, "ipAddressType": "dynamic"}
				item["radioStatuses"] = []map[string]any{
					{"band": "2.4G", "channel": a.Ch24, "channelBandwidth": "20"},
					{"band": "5G", "channel": a.Ch50, "channelBandwidth": "80"},
				}
			}
			list = append(list, item)
		}
		s.mu.Unlock()
		s.page(w, body, list)
	case len(seg) == 3 && seg[0] == "venues" && seg[2] == "aps" && r.Method == http.MethodPost:
		s.add(w, seg[1], body)
	case len(seg) >= 4 && seg[0] == "venues" && seg[2] == "aps":
		s.ap(w, r, seg, body)
	default:
		r1Fail(w, 404, "no such resource: "+path)
	}
}

var r1Serial = regexp.MustCompile(`^[1-9][0-9]{11}$`)

// add creates an AP in a venue, as adding one by serial number does.
func (s *R1) add(w http.ResponseWriter, venue string, body map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	known := false
	for _, v := range s.venues {
		known = known || v.ID == venue
	}
	serial, _ := body["serialNumber"].(string)
	name, _ := body["name"].(string)
	switch {
	case !known:
		r1Fail(w, 404, "Venue not found")
		return
	case !r1Serial.MatchString(serial):
		r1Fail(w, 400, "serialNumber is not a valid serial number")
		return
	case len(name) < 2 || len(name) > 32:
		r1Fail(w, 400, "name must be 2 to 32 characters")
		return
	}
	for _, a := range s.aps {
		if a.Serial == serial {
			r1Fail(w, 422, "AP with this serial number already exists")
			return
		}
	}
	ap := &r1AP{Serial: serial, Name: name, Venue: venue, Status: "1_01_NeverContactedCloud"}
	ap.Description, _ = body["description"].(string)
	s.aps = append(s.aps, ap)
	s.calls = append(s.calls, Call{Method: "POST", Path: "/venues/" + venue + "/aps", Body: body})
	reply(w, 202, map[string]any{"requestId": "req-1"})
}

func (s *R1) ap(w http.ResponseWriter, r *http.Request, seg []string, body map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var ap *r1AP
	for _, a := range s.aps {
		if a.Serial == seg[3] && a.Venue == seg[1] {
			ap = a
		}
	}
	if ap == nil {
		r1Fail(w, 404, "AP not found in this venue")
		return
	}
	change := r.Method == http.MethodPut
	if change {
		if s.FailSerial == ap.Serial {
			r1Fail(w, 422, "Validation failed")
			return
		}
		s.calls = append(s.calls, Call{Method: r.Method, Path: "/" + strings.Join(seg[3:], "/"), Body: body})
	}
	accepted := func() { reply(w, 202, map[string]any{"requestId": "req-1"}) }

	if len(seg) == 4 {
		switch r.Method {
		case http.MethodGet:
			out := map[string]any{"name": ap.Name, "description": ap.Description, "model": ap.Model, "tags": []string{}}
			if ap.Lat != "" {
				out["deviceGps"] = map[string]any{"latitude": ap.Lat, "longitude": ap.Lon}
			}
			reply(w, 200, out)
		case http.MethodPut:
			name, _ := body["name"].(string)
			if len(name) < 2 || len(name) > 32 {
				r1Fail(w, 400, "name must be 2 to 32 characters")
				s.calls = s.calls[:len(s.calls)-1]
				return
			}
			ap.Name = name
			ap.Description, _ = body["description"].(string)
			if g, ok := body["deviceGps"].(map[string]any); ok {
				ap.Lat, _ = g["latitude"].(string)
				ap.Lon, _ = g["longitude"].(string)
			}
			accepted()
		default:
			r1Fail(w, 405, "method not allowed")
		}
		return
	}
	res := seg[4]
	key := ap.Serial + "/" + res
	switch {
	case res == "networkSettings" && change:
		accepted()
	case strings.HasSuffix(res, "Settings") && r.Method == http.MethodGet:
		cur := s.settings[key]
		if cur == nil {
			cur = map[string]any{"useVenueSettings": false}
			if res == "radioSettings" {
				cur = map[string]any{
					"enable24G": true, "enable50G": true,
					"apRadioParams24G": map[string]any{"method": "CHANNELFLY", "manualChannel": 0, "channelBandwidth": "AUTO", "useVenueSettings": false},
					"apRadioParams50G": map[string]any{"method": "CHANNELFLY", "manualChannel": 0, "channelBandwidth": "AUTO", "useVenueSettings": false},
				}
			}
		}
		reply(w, 200, cur)
	case strings.HasSuffix(res, "Settings") && change:
		s.settings[key] = body
		if res == "radioSettings" {
			if p, ok := body["apRadioParams24G"].(map[string]any); ok && p["method"] == "MANUAL" {
				if ch, ok := p["manualChannel"].(float64); ok {
					ap.Ch24 = int(ch)
				}
			}
			if p, ok := body["apRadioParams50G"].(map[string]any); ok && p["method"] == "MANUAL" {
				if ch, ok := p["manualChannel"].(float64); ok {
					ap.Ch50 = int(ch)
				}
			}
		}
		accepted()
	default:
		r1Fail(w, 404, "no such resource")
	}
}

func (s *R1) page(w http.ResponseWriter, body map[string]any, list []map[string]any) {
	size := s.PageSize
	if size <= 0 {
		size = 1000
	}
	if n, ok := body["pageSize"].(float64); ok && int(n) > 0 && int(n) < size {
		size = int(n)
	}
	page := 1
	if n, ok := body["page"].(float64); ok && int(n) > 1 {
		page = int(n)
	}
	start := min((page-1)*size, len(list))
	end := min(start+size, len(list))
	reply(w, 200, map[string]any{"totalCount": len(list), "page": page, "data": list[start:end]})
}

func r1Fail(w http.ResponseWriter, status int, msg string) {
	reply(w, status, map[string]any{"requestId": "req-1", "errors": []map[string]any{{"code": fmt.Sprint(status), "message": msg}}})
}
