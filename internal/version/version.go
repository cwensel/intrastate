// Package version is the single source of build-identity metadata
// (semantic version, commit, build date). The fields are populated at
// link time via -ldflags -X (see the Makefile's LDFLAGS); a `go run`
// or test build leaves the dev defaults in place.
//
// Lives in its own leaf package so any package — the CLI root, an HTTP
// handler, a `--version` formatter — can read build identity without
// importing internal/cli.
package version

import "fmt"

// Set via ldflags at build time, e.g.:
//
//	go build -ldflags "\
//	  -X github.com/newcoinc/intrastate/internal/version.version=0.1.0 \
//	  -X github.com/newcoinc/intrastate/internal/version.commit=abc1234 \
//	  -X github.com/newcoinc/intrastate/internal/version.date=2026-06-17T00:00:00Z"
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Info is the resolved build identity.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns the build identity baked in at link time.
func Get() Info {
	return Info{Version: version, Commit: commit, Date: date}
}

// String renders the build identity for `--version` output.
func (i Info) String() string {
	return fmt.Sprintf("%s (commit %s, built %s)", i.Version, i.Commit, i.Date)
}
