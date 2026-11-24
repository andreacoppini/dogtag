**DogTag Plus** is a rebuild of DogTag: bulk naming, configuration and auditing
of RUCKUS access points through the controller that manages them. It now works
with **RUCKUS One** as well as **SmartZone**, works many APs at a time instead
of one after another, and is a single file with nothing to install.

## Getting started

Download the file for your platform, unzip it, and run it.

- **Windows** — `dogtag-plus-windows-amd64.exe` (use `-arm64` for Surface and
  Snapdragon machines). Double-click it, or run it from a terminal. If your
  network or browser refuses to download an `.exe`, take
  `dogtag-plus-windows-amd64.zip` instead: the same binary, zipped.
- **macOS** — `dogtag-plus-macos.zip`. One universal binary for Intel and Apple
  Silicon. It is unsigned, so the first launch needs
  `xattr -d com.apple.quarantine dogtag-plus-macos-universal`, or
  right-click → Open.
- **Linux** — `dogtag-plus-linux-amd64.tar.gz` (or `-arm64`).

Started with no arguments it opens a console in your browser, bound to
localhost. Everything it can do is also available as flags: run it with `-h`.

`SHA256SUMS.txt` covers every file in this release.

## What it does

- **AP naming and static channels.** Load a CSV keyed by MAC address and set
  names, descriptions, locations, GPS positions and static channels in one run.
- **Remove overrides.** Hand per-AP settings back to the zone or venue, switch
  radios on or off, and force APs back to DHCP.
- **Channels.** Count the connected APs on each 2.4, 5 and 6 GHz channel.
- **Export.** Write a zone's or venue's APs to a CSV that loads straight back
  in as a naming file.
- **Scaling.** The SmartZone AP and switch capacity calculator.

## What changed from DogTag 5

- **RUCKUS One** is supported alongside SmartZone, through the R1 API. There
  the AP's ID is its **serial number**: key the file by serial and it covers
  APs that have been added to a venue but are not onboarded yet, which have no
  MAC address to go by.
- APs are worked **in parallel**, eight at a time by default.
- The SmartZone API version is **negotiated** with the controller instead of
  being fixed at the one SmartZone 5.0 shipped with.
- Zones and AP lists of any size are read in full; the old 1000-entry cut-off
  is gone.
- CSV files are read properly: quoted commas, quotes inside names, a byte-order
  mark from Excel.
- Every action is available from the command line, and DogTag 5's own flags
  (`-ip=`, `-usr=`, `-pwd=`, `-prt=`, `-debug=yes`) still work.

The original Xojo application is preserved on the `legacy` branch.
