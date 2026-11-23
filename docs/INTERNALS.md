# Internals

## Shape

```
main.go      flags, the command line
ui.go        the console: a loopback HTTP server, an event stream, web/ embedded
job.go       the four jobs: naming, overrides, export, channels
backend.go   the Backend interface and its two implementations
csv.go       reading naming files and MAC lists; building a SmartZone patch
update.go    the once-a-day release check
sz/          SmartZone public API client
r1/          RUCKUS One API client
web/         index.html, app.js, style.css
```

The jobs are written against `Backend` and know nothing about which controller
is behind it. The console and the command line both call the same job
functions and differ only in what they do with the events a job emits: the
console streams them to the browser over server-sent events, the command line
prints them.

A job works its APs through `forEach`, a bounded worker pool (`-c`, default 8).
Requests for one AP go in order; APs go in parallel. Stop cancels the context:
requests in flight are abandoned and rows that were never reached are left
unmarked rather than counted as failures.

## The console

`serveUI` binds `127.0.0.1` on a free port and opens the browser. Because a
loopback port is reachable by any page open in that browser, every request
passes `guard`, which refuses a `Host` or `Origin` that is not the console's
own. That closes cross-site posts and DNS rebinding.

Secrets given on the command line stay in the process. `/api/defaults` tells
the page only that one exists, and `/api/connect` substitutes it when the
form's field is empty. The page keeps typed secrets in `sessionStorage` (gone
when the tab closes) and everything else, including the loaded file, in
`localStorage`.

The event history of the current job is replayed to a page that connects
mid-run, so a reload does not lose the results so far. Results are matched to
rows by index and MAC together, so a replay cannot land on a list that has
since been edited.

## SmartZone

All under `https://{host}:{port}/wsg/api/public`:

| Call | Used for |
|---|---|
| `GET /apiInfo` | which API versions are offered; the newest is used |
| `POST /{v}/session` | sign in (a session cookie) |
| `GET /{v}/session` | keep-alive, once a minute |
| `DELETE /{v}/session` | sign out |
| `GET /{v}/rkszones` | zones, paged |
| `GET /{v}/aps?zoneId=` | a zone's APs, paged |
| `GET /{v}/aps/{mac}/operational/summary` | one AP's detail, for export and channels |
| `PATCH /{v}/aps/{mac}` | name, description, location, GPS, `wifi24`/`wifi50` channel and width, `wlanService24Enabled`/`wlanService50Enabled` |
| `PATCH /{v}/aps/{mac}/network` | `ipType: Dynamic` |
| `DELETE /{v}/aps/{mac}/{override}` | remove one override |

If `apiInfo` is not answered the client falls back to `v7_0`, the version the
original used. Certificates are not verified unless `-verify-tls` is given,
since controllers almost always present a self-signed one. A 401 on anything
but sign-in means the session has ended; the job stops and the console
disconnects.

## RUCKUS One

Against `https://api.ruckus.cloud`, `api.eu.ruckus.cloud` or
`api.asia.ruckus.cloud`, with `application/vnd.ruckus.v1+json`:

| Call | Used for |
|---|---|
| `POST /oauth2/token/{tenantId}` | client credentials for a bearer token |
| `POST /venues/query` | venues |
| `POST /venues/aps/query` | APs with status, IP and per-radio channel; filtered by `venueId` |
| `GET /venues/{venueId}/aps/{serial}` | name, description, GPS |
| `PUT /venues/{venueId}/aps/{serial}` | the same, written back |
| `GET`, `PUT /venues/{venueId}/aps/{serial}/radioSettings` | static channel and width, radio on/off, `useVenueSettings` per band |
| `GET`, `PUT /venues/{venueId}/aps/{serial}/{x}Settings` | `useVenueSettings: true` |
| `PUT /venues/{venueId}/aps/{serial}/networkSettings` | `ipType: DYNAMIC` |

R1 addresses an AP by venue and serial number, while the files address it by
MAC. The backend reads the tenant's AP list once per run and keeps it as an
index.

Every write is read-modify-write. `PUT` replaces the resource, and the AP
record requires a name, so a row that sets only a description has to carry the
existing name with it; likewise pinning the 2.4 GHz channel must not disturb
5 GHz.

A 401 triggers one fresh sign-in (tokens expire). A 429 is retried with
back-off, honouring `Retry-After`.

## Trying it without a controller

```
go run ./cmd/fakesz
go run . -r1-url http://127.0.0.1:8444
```

Then connect to SmartZone at `127.0.0.1` (`admin` / `secret`), or to RUCKUS One
with tenant `tenant-1`, client ID `client-1`, secret `secret`. The same fakes
back the tests (`go test ./...`).

## Releasing

Add a `What's new in X.Y.Z (date):` block to the README changelog, update
`docs/RELEASE-NOTES.md`, and push a `vX.Y.Z` tag. The workflow refuses to
publish without both, runs the tests, stamps the version into the binaries and
the Windows resources, and attaches every platform plus `SHA256SUMS.txt`.
