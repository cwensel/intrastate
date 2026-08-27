// Package version is the single source of build-identity metadata
// (semantic version, commit, build date).
//
// Identity comes from two sources, ldflags first and Go's embedded VCS
// stamps as the fallback. The fallback matters: a plain `go build` or
// `go install` carries no ldflags, and without it such a binary reports
// only "dev" with no route back to a source revision. retrofit adopted
// the same fallback after five high-severity reports were filed from a
// binary that could say nothing else (its kata 2yvh); roborev and kata
// carry it too. This is that scheme, kept uniform with those projects.
//
// Lives in its own leaf package so any package — the CLI root, an HTTP
// handler, a `--version` formatter — can read build identity without
// importing internal/cli.
package version

import (
	"fmt"
	"runtime/debug"
)

// Set via ldflags at build time (see the Makefile's LDFLAGS):
//
//	go build -ldflags "\
//	  -X github.com/cwensel/intrastate/internal/version.version=0.1.0 \
//	  -X github.com/cwensel/intrastate/internal/version.commit=abc1234 \
//	  -X github.com/cwensel/intrastate/internal/version.date=2026-06-17T00:00:00Z"
//
// The zero-information defaults below are what the VCS fallback
// replaces; they survive only when the build carries no VCS stamps
// either (`-buildvcs=false`, or a build from outside a repository).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// The sentinels an unstamped build carries. They are the "ask the VCS
// fallback" signal, so they are compared against rather than spelled
// twice.
const (
	devVersion  = "dev"
	noCommit    = "none"
	unknownDate = "unknown"
)

// vcsRevisionLen is how much of a VCS revision hash is surfaced. Twelve
// hex chars matches retrofit and stays unambiguous well past this
// repo's object count. An ldflags commit is passed through verbatim —
// the Makefile already injects a short SHA, and truncating a release
// commit a second time would be wrong.
const vcsRevisionLen = 12

// readBuildInfo is indirected so tests can inject build settings.
// debug.ReadBuildInfo reports ok==false under `go test`, so the real
// function can never exercise the fallback from a unit test.
var readBuildInfo = debug.ReadBuildInfo

// Info is the resolved build identity.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Get returns this build's identity: ldflags where present, Go's
// embedded VCS stamps otherwise.
func Get() Info {
	return resolve(readBuildInfo)
}

// resolve is Get's testable core. ldflags values take precedence — a
// release build's output stays byte-identical to what it was before the
// fallback existed — and only the zero-information defaults are filled
// in from the VCS stamps.
func resolve(read func() (*debug.BuildInfo, bool)) Info {
	info := Info{Version: version, Commit: commit, Date: date}

	bi, ok := read()
	if !ok || bi == nil {
		return info
	}

	var revision, vcsTime string
	var dirty bool
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			vcsTime = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}

	if info.Commit == noCommit && revision != "" {
		if len(revision) > vcsRevisionLen {
			revision = revision[:vcsRevisionLen]
		}
		// The dirty marker is VCS-fallback presentation only. An ldflags
		// build folds it into the version string already — the Makefile
		// uses `git describe --dirty` — so adding it there would print it
		// twice.
		if dirty {
			revision += "-dirty"
		}
		info.Commit = revision
	}
	if info.Date == unknownDate && vcsTime != "" {
		info.Date = vcsTime
	}
	// `go install module@version` leaves no VCS settings at all; the
	// module version is the only identity such a build carries.
	// "(devel)" is the toolchain's placeholder for a local build and
	// says no more than "dev" does.
	if info.Version == devVersion {
		if mv := bi.Main.Version; mv != "" && mv != "(devel)" {
			info.Version = mv
		}
	}
	return info
}

// TextLine satisfies respond.TextLiner so `intrastate version` renders
// as the one build-identity line rather than a field-per-line payload,
// while still routing through the output gateway (and so, stdout).
// It is String() under the name the gateway asks for; both stay because
// the identity is also useful in error text and test failures.
func (i Info) TextLine() string { return i.String() }

// String renders the build identity for `--version` output.
func (i Info) String() string {
	return fmt.Sprintf("%s (commit %s, built %s)", i.Version, i.Commit, i.Date)
}
