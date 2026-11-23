// Command fakesz serves a pretend SmartZone controller and a pretend RUCKUS
// One tenant, for trying DogTag without real ones:
//
//	go run ./cmd/fakesz
//	dogtag-plus -r1-url http://127.0.0.1:8444
//
// then, in the console, connect to SmartZone at 127.0.0.1 (admin / secret) or
// to RUCKUS One with tenant "tenant-1", client ID "client-1", secret "secret".
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/andreacoppini/dogtag/internal/fakesz"
)

func main() {
	szAddr := flag.String("sz", "127.0.0.1:8443", "address for the fake SmartZone (HTTPS)")
	r1Addr := flag.String("r1", "127.0.0.1:8444", "address for the fake RUCKUS One (HTTP)")
	flag.Parse()

	s := fakesz.New()
	ln, err := net.Listen("tcp", *szAddr)
	if err != nil {
		log.Fatal(err)
	}
	// httptest brings its own self-signed certificate, which is exactly what
	// a real controller presents.
	ts := httptest.NewUnstartedServer(s)
	ts.Listener.Close()
	ts.Listener = ln
	ts.StartTLS()
	fmt.Printf("Fake SmartZone on https://%s  (user %q, password %q)\n", *szAddr, s.User, s.Pass)

	r := fakesz.NewR1()
	fmt.Printf("Fake RUCKUS One on http://%s  (tenant %q, client ID %q, secret %q)\nCtrl-C to stop.\n",
		*r1Addr, r.Tenant, r.ClientID, r.Secret)
	log.Fatal(http.ListenAndServe(*r1Addr, r))
}
