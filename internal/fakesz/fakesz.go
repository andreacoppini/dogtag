// Package fakesz is a stand-in SmartZone controller: enough of the public API
// to exercise DogTag end to end without a real one. The tests use it, and
// cmd/fakesz serves it for trying the console by hand.
package fakesz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// AP is one simulated access point.
type AP struct {
	MAC, Name, Description, Location string
	Latitude, Longitude              string
	Zone                             string
	Model, Serial, IP, IPType        string
	State                            string // "Connect", "Disconnect", ...
	Ch24, Ch50                       int
}

// Call is one change the fake was asked to make, kept so tests can check
// exactly what was sent.
type Call struct {
	Method string
	Path   string // relative to the versioned root, e.g. /aps/AA:BB:.../wifi24
	Body   map[string]any
}

// Server is the fake controller.
type Server struct {
	User, Pass string
	Versions   []string
	// FailMAC makes every change to that AP fail, to exercise error rows.
	FailMAC string
	// PageSize caps a listing page, to exercise paging. Zero means 1000.
	PageSize int

	mu       sync.Mutex
	zones    []zone
	aps      []*AP
	sessions map[string]bool
	calls    []Call
	nextSess int
}

type zone struct{ ID, Name string }

// New builds a controller with a small demo estate.
func New() *Server {
	s := &Server{
		User: "admin", Pass: "secret",
		Versions: []string{"v9_0", "v9_1", "v10_0", "v11_1"},
		sessions: map[string]bool{},
		zones: []zone{
			{"z-default", "Default Zone"},
			{"z-staging", "Staging Zone"},
			{"z-hq", "HQ Campus"},
			{"z-empty", "Empty Zone"},
		},
	}
	ch24 := []int{1, 6, 11}
	ch50 := []int{36, 40, 44, 48, 52, 100, 149, 157}
	models := []string{"R650", "R750", "R550", "T350d", "R350"}
	for i := 0; i < 42; i++ {
		z := "z-hq"
		if i >= 30 {
			z = "z-staging"
		}
		state := "Connect"
		if i%11 == 7 {
			state = "Disconnect"
		}
		s.aps = append(s.aps, &AP{
			MAC:      fmt.Sprintf("D4:68:4D:%02X:%02X:%02X", 0x10+i/16, 0x20+i, 0xA0+i),
			Name:     fmt.Sprintf("AP-%03d", i+1),
			Location: fmt.Sprintf("Floor %d", 1+i/8),
			Zone:     z,
			Model:    models[i%len(models)],
			Serial:   fmt.Sprintf("9819020%05d", 400+i),
			IP:       fmt.Sprintf("10.20.%d.%d", 1+i/20, 10+i),
			IPType:   "Dynamic",
			State:    state,
			Ch24:     ch24[i%len(ch24)],
			Ch50:     ch50[(i*3)%len(ch50)],
		})
	}
	return s
}

// Calls returns the changes made so far.
func (s *Server) Calls() []Call {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Call(nil), s.calls...)
}

// AP looks one up by MAC.
func (s *Server) AP(mac string) (AP, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.aps {
		if strings.EqualFold(a.MAC, mac) {
			return *a, true
		}
	}
	return AP{}, false
}

// EndSessions invalidates every session, as a controller-side timeout would.
func (s *Server) EndSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions = map[string]bool{}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const root = "/wsg/api/public/"
	if !strings.HasPrefix(r.URL.Path, root) {
		http.NotFound(w, r)
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, root)
	if rest == "apiInfo" {
		reply(w, 200, map[string]any{"apiSupportVersions": s.Versions})
		return
	}
	ver, path, _ := strings.Cut(rest, "/")
	if !contains(s.Versions, ver) {
		fail(w, 404, "API version "+ver+" is not supported")
		return
	}
	path = "/" + path

	if path == "/session" && r.Method == http.MethodPost {
		var b struct{ Username, Password string }
		_ = json.NewDecoder(r.Body).Decode(&b)
		if b.Username != s.User || b.Password != s.Pass {
			fail(w, 401, "Invalid username or password")
			return
		}
		s.mu.Lock()
		s.nextSess++
		id := "sess-" + strconv.Itoa(s.nextSess)
		s.sessions[id] = true
		s.mu.Unlock()
		http.SetCookie(w, &http.Cookie{Name: "JSESSIONID", Value: id, Path: "/wsg"})
		reply(w, 200, map[string]any{"controllerVersion": "6.1.2.0.441"})
		return
	}

	ck, err := r.Cookie("JSESSIONID")
	s.mu.Lock()
	live := err == nil && s.sessions[ck.Value]
	s.mu.Unlock()
	if !live {
		fail(w, 401, "No active session")
		return
	}

	seg := strings.Split(strings.Trim(path, "/"), "/")
	switch {
	case path == "/session" && r.Method == http.MethodGet:
		reply(w, 200, map[string]any{"cpId": "fake", "apiVersions": s.Versions})
	case path == "/session" && r.Method == http.MethodDelete:
		s.mu.Lock()
		delete(s.sessions, ck.Value)
		s.mu.Unlock()
		w.WriteHeader(204)
	case path == "/rkszones" && r.Method == http.MethodGet:
		var list []map[string]any
		for _, z := range s.zones {
			list = append(list, map[string]any{"id": z.ID, "name": z.Name})
		}
		s.page(w, r, list)
	case path == "/aps" && r.Method == http.MethodGet:
		zoneID := r.URL.Query().Get("zoneId")
		list := []map[string]any{}
		s.mu.Lock()
		for _, a := range s.aps {
			if zoneID == "" || a.Zone == zoneID {
				list = append(list, map[string]any{"mac": a.MAC, "zoneId": a.Zone, "apGroupId": "g-" + a.Zone, "name": a.Name})
			}
		}
		s.mu.Unlock()
		s.page(w, r, list)
	case len(seg) == 4 && seg[0] == "aps" && seg[2] == "operational" && seg[3] == "summary" && r.Method == http.MethodGet:
		a, ok := s.AP(seg[1])
		if !ok {
			fail(w, 404, "AP not found")
			return
		}
		out := map[string]any{
			"mac": a.MAC, "name": a.Name, "description": a.Description, "location": a.Location,
			"latitude": nil, "longitude": nil,
			"wifi24Channel": a.Ch24, "wifi50Channel": a.Ch50,
			"serial": a.Serial, "connectionState": a.State, "model": a.Model,
			"ipType": a.IPType, "ip": a.IP, "apGroupId": "g-" + a.Zone,
		}
		if a.Latitude != "" {
			out["latitude"], _ = strconv.ParseFloat(a.Latitude, 64)
			out["longitude"], _ = strconv.ParseFloat(a.Longitude, 64)
		}
		reply(w, 200, out)
	case len(seg) >= 2 && seg[0] == "aps" && (r.Method == http.MethodPatch || r.Method == http.MethodDelete):
		s.change(w, r, seg)
	default:
		fail(w, 404, "no such resource: "+path)
	}
}

func (s *Server) change(w http.ResponseWriter, r *http.Request, seg []string) {
	mac := seg[1]
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)

	s.mu.Lock()
	defer s.mu.Unlock()
	var ap *AP
	for _, a := range s.aps {
		if strings.EqualFold(a.MAC, mac) {
			ap = a
		}
	}
	if ap == nil {
		fail(w, 404, "AP "+mac+" not found")
		return
	}
	if s.FailMAC != "" && strings.EqualFold(s.FailMAC, mac) {
		fail(w, 422, "Validation failed for "+mac)
		return
	}
	s.calls = append(s.calls, Call{Method: r.Method, Path: "/" + strings.Join(seg, "/"), Body: body})

	if r.Method == http.MethodPatch && len(seg) == 2 {
		if v, ok := body["name"].(string); ok {
			ap.Name = v
		}
		if v, ok := body["description"].(string); ok {
			ap.Description = v
		}
		if v, ok := body["location"].(string); ok {
			ap.Location = v
		}
		if v, ok := body["latitude"].(float64); ok {
			ap.Latitude = strconv.FormatFloat(v, 'f', -1, 64)
		}
		if v, ok := body["longitude"].(float64); ok {
			ap.Longitude = strconv.FormatFloat(v, 'f', -1, 64)
		}
		if m, ok := body["wifi24"].(map[string]any); ok {
			if ch, ok := m["channel"].(float64); ok {
				ap.Ch24 = int(ch)
			}
		}
		if m, ok := body["wifi50"].(map[string]any); ok {
			if ch, ok := m["channel"].(float64); ok {
				ap.Ch50 = int(ch)
			}
		}
	}
	if r.Method == http.MethodPatch && len(seg) == 3 && seg[2] == "network" {
		if v, ok := body["ipType"].(string); ok {
			ap.IPType = v
		}
	}
	w.WriteHeader(204)
}

func (s *Server) page(w http.ResponseWriter, r *http.Request, list []map[string]any) {
	size := s.PageSize
	if size <= 0 {
		size = 1000
	}
	if n, err := strconv.Atoi(r.URL.Query().Get("listSize")); err == nil && n > 0 && n < size {
		size = n
	}
	index, _ := strconv.Atoi(r.URL.Query().Get("index"))
	if index > len(list) {
		index = len(list)
	}
	end := min(index+size, len(list))
	reply(w, 200, map[string]any{
		"totalCount": len(list), "hasMore": end < len(list), "firstIndex": index,
		"list": list[index:end],
	})
}

func reply(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, status int, msg string) {
	reply(w, status, map[string]any{"message": msg, "errorCode": status, "errorType": "Error"})
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
