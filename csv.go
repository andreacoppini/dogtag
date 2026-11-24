package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// NameRow is one line of an AP naming file: the ID that identifies
// the AP, and what to set on it. Blank fields are left alone.
type NameRow struct {
	// ID is the AP's MAC address or its serial number. SmartZone knows APs by
	// MAC; RUCKUS One knows them by serial, and an AP that has been added
	// there but has never come online has no MAC to go by at all.
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Width24     string `json:"width24"`
	Channel24   string `json:"channel24"`
	Width50     string `json:"width50"`
	Channel50   string `json:"channel50"`
}

// exportHeader is the layout of a zone export. Its first ten columns are a
// naming file, so an export can be edited in a spreadsheet and loaded
// straight back in.
var exportHeader = []string{"MAC", "Name", "Description", "Location", "GPSLatitude", "GPSLongitude",
	"24Width", "24Channel", "50Width", "50Channel", "Serial", "Status", "Model", "IPType", "IPAddress"}

var macRe = regexp.MustCompile(`^([0-9A-Fa-f]{2})[:-]([0-9A-Fa-f]{2})[:-]([0-9A-Fa-f]{2})[:-]([0-9A-Fa-f]{2})[:-]([0-9A-Fa-f]{2})[:-]([0-9A-Fa-f]{2})$`)

// serialRe is a RUCKUS serial number: twelve digits.
var serialRe = regexp.MustCompile(`^[0-9]{12}$`)

// normalMAC returns the MAC in the form the controllers use, or "" if the
// text is not one.
func normalMAC(s string) string {
	m := macRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return ""
	}
	return strings.ToUpper(strings.Join(m[1:], ":"))
}

// isSerial reports whether an ID is a serial number rather than a MAC.
func isSerial(id string) bool { return serialRe.MatchString(id) }

// normalID returns the AP's ID — a MAC address or a serial number — in its
// standard form, or "" if the text is neither. Anything else in the first
// column — a header, a note — is how a line gets skipped.
func normalID(s string) string {
	if mac := normalMAC(s); mac != "" {
		return mac
	}
	if s = strings.TrimSpace(s); isSerial(s) {
		return s
	}
	return ""
}

// readRecords parses CSV text leniently: files that have been through a
// spreadsheet carry a BOM, ragged rows and stray quotes, and none of that is
// a reason to refuse the file.
func readRecords(text string) [][]string {
	text = strings.TrimPrefix(text, "\ufeff")
	r := csv.NewReader(strings.NewReader(text))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	var out [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			return out
		}
		if err != nil {
			continue // one unreadable line does not cost the rest
		}
		out = append(out, rec)
	}
}

// parseNaming reads a naming file. Lines whose first column is neither a MAC
// address nor a serial number are skipped and counted.
func parseNaming(text string) (rows []NameRow, skipped int) {
	for _, rec := range readRecords(text) {
		if len(rec) == 0 || (len(rec) == 1 && strings.TrimSpace(rec[0]) == "") {
			continue
		}
		id := normalID(rec[0])
		if id == "" {
			skipped++
			continue
		}
		col := func(i int) string {
			if i < len(rec) {
				return strings.TrimSpace(rec[i])
			}
			return ""
		}
		row := NameRow{ID: id, Name: col(1), Description: col(2), Location: col(3),
			Latitude: col(4), Longitude: col(5),
			Width24: col(6), Channel24: col(7), Width50: col(8), Channel50: col(9)}
		// A channel given without a width means 20 MHz, as it always has here.
		if row.Width24 == "" {
			if ch, _ := strconv.Atoi(row.Channel24); ch > 0 && ch < 15 {
				row.Width24 = "20"
			}
		}
		if row.Width50 == "" {
			if ch, _ := strconv.Atoi(row.Channel50); ch > 30 && ch < 200 {
				row.Width50 = "20"
			}
		}
		rows = append(rows, row)
	}
	return rows, skipped
}

// parseIDs reads the first column of a file as a list of APs, by MAC address
// or serial number.
func parseIDs(text string) (ids []string, skipped int) {
	seen := map[string]bool{}
	for _, rec := range readRecords(text) {
		if len(rec) == 0 || strings.TrimSpace(rec[0]) == "" {
			continue
		}
		id := normalID(rec[0])
		if id == "" {
			skipped++
			continue
		}
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids, skipped
}

// patch builds the change to send for one row. Only what the row actually
// specifies is included, so a file that carries names alone does not disturb
// anything else on the AP.
func (r NameRow) patch() (map[string]any, error) {
	p := map[string]any{}
	if r.Name != "" {
		p["name"] = r.Name
	}
	if r.Description != "" {
		p["description"] = r.Description
	}
	if r.Location != "" {
		p["location"] = r.Location
	}
	if r.Latitude != "" || r.Longitude != "" {
		lat, err1 := strconv.ParseFloat(r.Latitude, 64)
		lon, err2 := strconv.ParseFloat(r.Longitude, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("latitude and longitude must both be numbers (got %q, %q)", r.Latitude, r.Longitude)
		}
		p["latitude"], p["longitude"] = lat, lon
	}
	for _, radio := range []struct{ key, width, channel string }{
		{"wifi24", r.Width24, r.Channel24},
		{"wifi50", r.Width50, r.Channel50},
	} {
		m := map[string]any{}
		if radio.width != "" {
			n, err := strconv.Atoi(radio.width)
			if err != nil {
				return nil, fmt.Errorf("channel width %q is not a number", radio.width)
			}
			m["channelWidth"] = n
		}
		if radio.channel != "" {
			n, err := strconv.Atoi(radio.channel)
			if err != nil {
				return nil, fmt.Errorf("channel %q is not a number", radio.channel)
			}
			m["channel"] = n
		}
		if len(m) > 0 {
			p[radio.key] = m
		}
	}
	if len(p) == 0 {
		return nil, fmt.Errorf("nothing to change: every column after the first is empty")
	}
	return p, nil
}
