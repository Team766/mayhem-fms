Mechanical M-Ayhem FMS - based on [Team 254's Cheesy Arena](https://github.com/Team254/cheesy-arena).
============
This repo contains [Cheesy Arena](https://github.com/Team254/cheesy-arena), to be used as the Field Management System
for the [Mechanical M-Ayhem](https://www.team766.com/mechanical-m-ayhem-1) rookie competition run by Team 766.
To make keeping up with Cheesy Arena and year-specific game customization easier, this codebase is not a long-lived
fork.  It is upstream Cheesy Arena minus a few features M-Ayhem doesn't use, plus a few it needs: optional 2v2 gameplay,
changes for the Arduino PLC we use, and single-game playoff rounds with a best of 3 final.  The current year's game is
applied on top from a short game spec.

Two playbooks, written so that a person or a coding agent can follow them, keep this up to date: one re-applies our
changes onto each new Cheesy Arena release, and the other applies each year's game from its spec.  See
[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for how this all fits together, and [docs/agents](docs/agents) for the
playbooks.

## Manuals

- [Operator manual](docs/OPERATOR_MANUAL.md): building, setting up and running the FMS for an event.
- [Scorer manual](docs/SCORER_MANUAL.md) and [referee manual](docs/REFEREE_MANUAL.md): for the volunteers on the
  tablets.
- [Foul card](docs/FOUL_CARD.md): a one-page summary of the fouls for referees.

## License

Teams may use M-Ayhem FMS freely for practice, scrimmages, and off-season events. See [LICENSE](LICENSE) for more
details.

## Installing

**From source**

1. Download [Go](https://golang.org/dl/) (the version in `go.mod` or later)
1. Clone this GitHub repository to a location of your choice
1. Navigate to the repository's directory in the terminal
1. Compile the code with `go build`
1. Run the `cheesy-arena` or `cheesy-arena.exe` binary
1. Navigate to http://localhost:8080 in your browser (Google Chrome recommended)

**IP address configuration**

When running Cheesy Arena on a playing field with robots, set the IP address of the computer running Cheesy Arena to
10.0.100.5. By a convention baked into the FRC Driver Station software, driver stations will broadcast their presence on
the network to this hardcoded address so that the FMS does not need to discover them by some other method.

When running Cheesy Arena without robots for testing or development, pass the `-dev` flag to bind driver station
listeners to any local IP address.

## Under the hood

Cheesy Arena is written using [Go](https://golang.org), a language developed by Google and first released in 2009. Go
excels in the areas of concurrency, networking, performance, and portability, which makes it ideal for a field
management system.

Cheesy Arena is implemented as a web server, with all human interaction done via browser. The graphical interfaces are
implemented in HTML, JavaScript, and CSS. There are many advantages to this approach &ndash; development of new
graphical elements is rapid, and no software needs to be installed other than on the server. Client web pages send
commands and receive updates using WebSockets.

[Bolt](https://github.com/etcd-io/bbolt) is used as the datastore, and making backups or transferring data from one
installation to another is as simple as copying the database file.

Schedule generation is fast because pre-generated schedules are included with the code. Each schedule contains a certain
number of matches per team for placeholder teams 1 through N, so generating the actual match schedule becomes a simple
exercise in permuting the mapping of real teams to placeholder teams. The pre-generated schedules are checked into this
repository and can be vetted in advance of any events for deviations from the randomness (and other) requirements.

Cheesy Arena includes support for, but doesn't require, networking hardware similar to that used in official FRC events.
Teams are issued their own SSIDs and WPA keys, and when connected to Cheesy Arena are isolated to a VLAN which prevents
any communication other than between the driver station, robot, and event server. The network hardware is reconfigured
via SSH and Telnet commands for the new set of teams when each mach is loaded.

## PLC integration

Mechanical M-Ayhem uses a custom designed Arduino-based PLC, protocol compatible with a subset of the PLC Cheesy Arena uses.
PLC functionality used by M-Ayhem consists of team e-stops and a-stops and the stack lights showing alliance readiness.

The PLC code can be found [here](https://github.com/Team766/fakeplc-arduino), and how the FMS talks to it is in
[docs/ArduinoPlc.md](docs/ArduinoPlc.md).

## Advanced networking

See the [Advanced Networking wiki page](https://github.com/Team254/cheesy-arena/wiki/Advanced-Networking-Concepts) for
instructions on what equipment to obtain and how to configure it in order to support advanced network security.
