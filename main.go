// Command dogtag-plus names, configures and audits RUCKUS access points in
// bulk, through the controller that manages them: a SmartZone, or a RUCKUS
// One tenant.
//
// Run it with no arguments for the browser console, or drive it from the
// command line:
//
//	dogtag-plus -ip 10.0.0.5 -user admin -ask-pass -zone "HQ Campus" -export aps.csv
//	dogtag-plus -ip 10.0.0.5 -user admin -ask-pass -names aps.csv
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

// version is overwritten at release time with -X main.version. The fallback
// is deliberately not a number: a locally built binary claiming to be a
// release it is not is worse than one that admits what it is.
var version = "dev"

type options struct {
	host, port, user, pass string
	askPass                bool
	passEnv                string
	verifyTLS              bool
	apiVersion             string

	r1          bool
	r1Region    string
	r1Tenant    string
	r1ClientID  string
	r1Secret    string
	r1SecretEnv string
	r1AskSecret bool
	r1URL       string

	zone          string
	exportPath    string
	namesPath     string
	overridesPath string
	channels      bool
	overrides     string
	radio24       string
	radio50       string
	ipDynamic     bool
	out           string

	concurrency int
	timeout     time.Duration
	debug       bool
	ui          bool
	uiPort      int
	showVers    bool
	noUpdate    bool
}

// yesNo is a boolean flag that also takes the words the original DogTag
// accepted: -debug=yes.
type yesNo struct{ v *bool }

func (y yesNo) String() string   { return "" }
func (y yesNo) IsBoolFlag() bool { return true }
func (y yesNo) Set(s string) error {
	switch strings.ToLower(s) {
	case "yes", "y", "true", "1", "on":
		*y.v = true
	case "no", "n", "false", "0", "off":
		*y.v = false
	default:
		return fmt.Errorf("use yes or no")
	}
	return nil
}

func parseFlags(args []string) (options, error) {
	var o options
	fs := flag.NewFlagSet("dogtag-plus", flag.ContinueOnError)
	alias := func(p *string, def, usage string, names ...string) {
		for i, n := range names {
			u := usage
			if i > 0 {
				u = "same as -" + names[0]
			}
			fs.StringVar(p, n, def, u)
		}
	}
	// The original's spellings come first so its command lines keep working.
	alias(&o.host, "", "SmartZone address", "ip", "ipaddr", "ipaddress", "host")
	alias(&o.port, "8443", "SmartZone API port", "port", "prt")
	alias(&o.user, "admin", "SmartZone username", "user", "usr")
	alias(&o.pass, "", "SmartZone password (prefer -ask-pass: an argument is visible in the process list)", "pass", "pwd", "password")
	fs.BoolVar(&o.askPass, "ask-pass", false, "prompt for the SmartZone password")
	fs.StringVar(&o.passEnv, "pass-env", "", "read the SmartZone password from this environment variable")
	fs.BoolVar(&o.verifyTLS, "verify-tls", false, "verify the controller's certificate (off by default: most are self-signed)")
	fs.StringVar(&o.apiVersion, "api-version", "", "SmartZone API version, e.g. v9_1 (default: the newest the controller offers)")

	fs.BoolVar(&o.r1, "r1", false, "work on a RUCKUS One tenant instead of a SmartZone")
	fs.StringVar(&o.r1Region, "r1-region", "us", "RUCKUS One region: us, eu or asia")
	fs.StringVar(&o.r1Tenant, "r1-tenant", "", "RUCKUS One tenant ID")
	fs.StringVar(&o.r1ClientID, "r1-client-id", "", "RUCKUS One application token client ID")
	fs.StringVar(&o.r1Secret, "r1-secret", "", "RUCKUS One shared secret (prefer -r1-ask-secret)")
	fs.StringVar(&o.r1SecretEnv, "r1-secret-env", "", "read the RUCKUS One shared secret from this environment variable")
	fs.BoolVar(&o.r1AskSecret, "r1-ask-secret", false, "prompt for the RUCKUS One shared secret")
	fs.StringVar(&o.r1URL, "r1-url", "", "RUCKUS One API base URL, instead of the region's (for testing)")

	fs.StringVar(&o.zone, "zone", "", "zone (SmartZone) or venue (RUCKUS One) to export or count, by name or ID")
	fs.StringVar(&o.zone, "venue", "", "same as -zone")
	fs.StringVar(&o.exportPath, "export", "", "write the zone's APs to this CSV file")
	fs.StringVar(&o.namesPath, "names", "", "apply this naming CSV: names, descriptions, locations, GPS, static channels")
	fs.StringVar(&o.overridesPath, "remove-overrides", "", "remove per-AP overrides from the APs listed in this CSV")
	fs.StringVar(&o.overrides, "overrides", "default", `which overrides to remove: "default", "all", "none", or a comma-separated list`)
	fs.StringVar(&o.radio24, "radio24", "skip", "with -remove-overrides, switch the 2.4 GHz radio: on, off or skip")
	fs.StringVar(&o.radio50, "radio50", "skip", "with -remove-overrides, switch the 5 GHz radio: on, off or skip")
	fs.BoolVar(&o.ipDynamic, "ip-dynamic", false, "with -remove-overrides, force the AP's IP address to dynamic")
	fs.BoolVar(&o.channels, "channels", false, "count the connected APs on each channel of the zone")
	fs.StringVar(&o.out, "out", "", "write the per-AP results of -names or -remove-overrides to this CSV file")

	fs.IntVar(&o.concurrency, "c", 8, "requests in flight at once")
	fs.DurationVar(&o.timeout, "timeout", 30*time.Second, "how long to wait for one answer from the controller")
	fs.Var(yesNo{&o.debug}, "debug", "log every API request and answer")
	fs.BoolVar(&o.ui, "ui", false, "open the browser console (the default when no action is given)")
	fs.IntVar(&o.uiPort, "ui-port", 0, "port for the console (default: any free port)")
	fs.BoolVar(&o.showVers, "version", false, "print the version and exit")
	fs.BoolVar(&o.noUpdate, "no-update-check", false, "do not ask GitHub whether a newer version exists (also DOGTAG_NO_UPDATE_CHECK=1)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "DogTag Plus %s\n\nUsage: %s [flags]\n\nWith no action flag (-export, -names, -remove-overrides, -channels) it opens\nthe console in your browser; connection flags pre-fill it.\n\n", version, filepath.Base(os.Args[0]))
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return o, err
	}
	if o.concurrency < 1 {
		o.concurrency = 1
	}
	return o, nil
}

func (o options) hasAction() bool {
	return o.exportPath != "" || o.namesPath != "" || o.overridesPath != "" || o.channels
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	opt, err := parseFlags(args)
	if err != nil {
		return err
	}
	if opt.showVers {
		fmt.Println("DogTag Plus", version)
		return nil
	}
	if opt.passEnv != "" && opt.pass == "" {
		opt.pass = os.Getenv(opt.passEnv)
	}
	if opt.r1SecretEnv != "" && opt.r1Secret == "" {
		opt.r1Secret = os.Getenv(opt.r1SecretEnv)
	}
	if opt.askPass {
		if opt.pass, err = prompt("SmartZone password: "); err != nil {
			return err
		}
	}
	if opt.r1AskSecret {
		if opt.r1Secret, err = prompt("RUCKUS One shared secret: "); err != nil {
			return err
		}
	}

	// No action means the console: that is what a double-click gets, and what
	// the original did with its connection flags too.
	if opt.ui || !opt.hasAction() {
		return serveUI(opt)
	}
	return runCLI(opt)
}

func prompt(label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("cannot read from the terminal: %w", err)
	}
	return string(b), nil
}

func (o options) connectRequest() connectRequest {
	kind := "sz"
	if o.r1 {
		kind = "r1"
	}
	return connectRequest{Kind: kind, Host: o.host, Port: o.port, User: o.user, Pass: o.pass,
		VerifyTLS: o.verifyTLS, APIVer: o.apiVersion,
		Region: o.r1Region, Tenant: o.r1Tenant, ClientID: o.r1ClientID, Secret: o.r1Secret}
}

func runCLI(opt options) error {
	if opt.out != "" && opt.namesPath != "" && opt.overridesPath != "" {
		return errors.New("-out holds one run's results: use -names and -remove-overrides in separate runs to keep both")
	}
	update := startUpdateCheck(opt)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logLine := func(s string) { fmt.Fprintln(os.Stderr, s) }
	b, err := connect(ctx, opt, opt.connectRequest(), logLine)
	if err != nil {
		return err
	}
	defer func() {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = b.Close(c)
	}()
	logLine("Connected: " + b.Describe())

	var zone Zone
	if opt.exportPath != "" || opt.channels {
		zones, err := b.Zones(ctx)
		if err != nil {
			return err
		}
		if zone, err = pickZone(zones, opt.zone); err != nil {
			return err
		}
	}

	results := map[int]*RowResult{}
	failed := 0
	var mu sync.Mutex
	emit := func(e Event) {
		mu.Lock()
		defer mu.Unlock()
		switch e.Kind {
		case EvLog:
			logLine(e.Message)
		case EvRow:
			results[e.Row.Index] = e.Row
			if e.Row.Error != "" {
				failed++
			}
			logLine(fmt.Sprintf("[%d/%d] %s  %s", e.Done+e.Failed, e.Total, e.Row.ID, firstNonEmpty(prefixed("ERROR: ", e.Row.Error), e.Row.OK)))
		}
	}
	var ids []string
	writeOut := func() error {
		defer func() { results, ids = map[int]*RowResult{}, nil }()
		return writeResults(opt.out, ids, results)
	}

	if opt.namesPath != "" {
		text, err := os.ReadFile(opt.namesPath)
		if err != nil {
			return err
		}
		rows, skipped := parseNaming(string(text))
		if len(rows) == 0 {
			return fmt.Errorf("no APs found in %s: the first column of each row must be a MAC address or a serial number", opt.namesPath)
		}
		logLine(fmt.Sprintf("%d APs to update (%d lines skipped)", len(rows), skipped))
		for _, r := range rows {
			ids = append(ids, r.ID)
		}
		if err := runNaming(ctx, b, rows, opt.concurrency, emit); err != nil {
			return err
		}
		if err := writeOut(); err != nil {
			return err
		}
	}
	if opt.overridesPath != "" && ctx.Err() == nil {
		text, err := os.ReadFile(opt.overridesPath)
		if err != nil {
			return err
		}
		list, skipped := parseIDs(string(text))
		if len(list) == 0 {
			return fmt.Errorf("no APs found in %s: the first column of each row must be a MAC address or a serial number", opt.overridesPath)
		}
		ids, err := chooseOverrides(b.OverrideOptions(), opt.overrides)
		if err != nil {
			return err
		}
		logLine(fmt.Sprintf("%d APs, %d overrides each (%d lines skipped)", len(list), len(ids), skipped))
		ids = list
		job := OverrideJob{IDs: list, Overrides: ids, Radio24: opt.radio24, Radio50: opt.radio50, IPDynamic: opt.ipDynamic}
		if err := runOverrides(ctx, b, job, opt.concurrency, emit); err != nil {
			return err
		}
		if err := writeOut(); err != nil {
			return err
		}
	}
	if opt.exportPath != "" && ctx.Err() == nil {
		records, err := runExport(ctx, b, zone.ID, opt.concurrency, emit)
		if err != nil {
			return err
		}
		if err := writeCSV(opt.exportPath, records); err != nil {
			return err
		}
		logLine(fmt.Sprintf("%d APs written to %s", len(records)-1, opt.exportPath))
	}
	if opt.channels && ctx.Err() == nil {
		cc, err := runChannels(ctx, b, zone.ID, zone.Name, opt.concurrency, emit)
		if err != nil {
			return err
		}
		printChannels(cc)
	}
	if ctx.Err() != nil {
		return errStopped
	}
	reportUpdate(update, os.Stderr)
	if failed > 0 {
		// A script has to be able to tell a clean run from a partial one.
		return fmt.Errorf("%d APs could not be changed", failed)
	}
	return nil
}

func prefixed(p, s string) string {
	if s == "" {
		return ""
	}
	return p + s
}

// pickZone finds the zone the operator named, by name or ID.
func pickZone(zones []Zone, want string) (Zone, error) {
	if want == "" {
		return Zone{}, fmt.Errorf("say which zone or venue with -zone; there are %d: %s", len(zones), zoneNames(zones))
	}
	for _, z := range zones {
		if z.ID == want || strings.EqualFold(z.Name, want) {
			return z, nil
		}
	}
	return Zone{}, fmt.Errorf("no zone or venue called %q; there are: %s", want, zoneNames(zones))
}

func zoneNames(zones []Zone) string {
	names := make([]string, len(zones))
	for i, z := range zones {
		names[i] = z.Name
	}
	return strings.Join(names, ", ")
}

// chooseOverrides turns the -overrides flag into the list to remove.
func chooseOverrides(all []OverrideOption, spec string) ([]string, error) {
	var out []string
	switch strings.ToLower(strings.TrimSpace(spec)) {
	case "none":
		return nil, nil
	case "", "default":
		for _, o := range all {
			if o.Default {
				out = append(out, o.ID)
			}
		}
		return out, nil
	case "all":
		for _, o := range all {
			out = append(out, o.ID)
		}
		return out, nil
	}
	known := map[string]string{}
	for _, o := range all {
		known[strings.ToLower(o.ID)] = o.ID
	}
	for _, f := range strings.Split(spec, ",") {
		id, ok := known[strings.ToLower(strings.TrimSpace(f))]
		if !ok {
			return nil, fmt.Errorf("%q is not an override this controller offers", strings.TrimSpace(f))
		}
		out = append(out, id)
	}
	return out, nil
}

func printChannels(cc *ChannelCount) {
	fmt.Printf("%s: %d of %d APs connected\n", cc.Zone, cc.Connected, cc.Total)
	for _, band := range []struct {
		name string
		rows []ChanRow
	}{{"2.4 GHz", cc.Band24}, {"5 GHz", cc.Band50}, {"6 GHz", cc.Band6}} {
		if len(band.rows) == 0 {
			continue
		}
		fmt.Printf("\n%s\n%8s %6s\n", band.name, "Channel", "APs")
		for _, r := range band.rows {
			fmt.Printf("%8d %6d\n", r.Channel, r.APs)
		}
	}
}

func writeCSV(path string, records [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// writeResults writes one line per AP, including any the run never reached.
func writeResults(path string, ids []string, results map[int]*RowResult) error {
	if path == "" {
		return nil
	}
	records := [][]string{{"AP", "OK", "Error"}}
	for i, mac := range ids {
		rec := []string{mac, "", "not attempted"}
		if r := results[i]; r != nil {
			rec = []string{mac, r.OK, r.Error}
		}
		records = append(records, rec)
	}
	return writeCSV(path, records)
}
