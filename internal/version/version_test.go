package version

// The build-identity contract.
//
// `resolve` is exercised over injected build settings because the real
// debug.ReadBuildInfo reports ok==false under `go test` — the fallback
// path is unreachable from a unit test otherwise.
//
// The invariant every case below defends: ldflags win, and the VCS
// stamps only ever fill in a field that carried no information. A
// release build's output must stay byte-identical to what it was before
// the fallback existed.

import (
	"runtime/debug"
	"strings"
	"testing"
)

// withSettings builds a read function returning the given VCS settings,
// plus an optional main module version.
func withSettings(mainVersion string, kv ...string) func() (*debug.BuildInfo, bool) {
	var settings []debug.BuildSetting
	for i := 0; i+1 < len(kv); i += 2 {
		settings = append(settings, debug.BuildSetting{Key: kv[i], Value: kv[i+1]})
	}
	bi := &debug.BuildInfo{Settings: settings}
	bi.Main.Version = mainVersion
	return func() (*debug.BuildInfo, bool) { return bi, true }
}

// stubLdflags sets the package-level ldflags vars for one test and
// restores them after, so cases can model a release build.
func stubLdflags(t *testing.T, v, c, d string) {
	t.Helper()
	ov, oc, od := version, commit, date
	version, commit, date = v, c, d
	t.Cleanup(func() { version, commit, date = ov, oc, od })
}

// TestResolve_LdflagsWinOverVCS is the release-parity guarantee. A
// stamped build must ignore the VCS settings entirely — including the
// dirty marker, which the Makefile already folds into VERSION via
// `git describe --dirty`. Honouring both would print it twice.
func TestResolve_LdflagsWinOverVCS(t *testing.T) {
	stubLdflags(t, "v1.2.3", "abc1234", "2026-01-01T00:00:00Z")
	got := resolve(withSettings("v9.9.9",
		"vcs.revision", "ffffffffffffffffffffffffffffffffffffffff",
		"vcs.time", "2099-12-31T23:59:59Z",
		"vcs.modified", "true",
	))
	want := Info{Version: "v1.2.3", Commit: "abc1234", Date: "2026-01-01T00:00:00Z"}
	if got != want {
		t.Errorf("resolve() = %+v; want %+v — ldflags must win outright", got, want)
	}
	if strings.Contains(got.Commit, "-dirty") {
		t.Errorf("commit %q carries a dirty marker; an ldflags build folds "+
			"it into the version string already", got.Commit)
	}
}

// TestResolve_FallsBackToVCS is the gap this closes: an unstamped build
// reported "dev / none / unknown" while carrying the revision all along.
func TestResolve_FallsBackToVCS(t *testing.T) {
	stubLdflags(t, devVersion, noCommit, unknownDate)
	got := resolve(withSettings("",
		"vcs.revision", "4bfa8682c08482e357b40bffa31ff8097ca05ea9",
		"vcs.time", "2026-08-27T16:15:53Z",
		"vcs.modified", "false",
	))
	if got.Commit != "4bfa8682c084" {
		t.Errorf("commit = %q; want the revision truncated to %d chars",
			got.Commit, vcsRevisionLen)
	}
	if got.Date != "2026-08-27T16:15:53Z" {
		t.Errorf("date = %q; want the vcs.time stamp", got.Date)
	}
}

// TestResolve_MarksADirtyTree pins the one marker the VCS path owns.
func TestResolve_MarksADirtyTree(t *testing.T) {
	stubLdflags(t, devVersion, noCommit, unknownDate)
	got := resolve(withSettings("",
		"vcs.revision", "4bfa8682c08482e357b40bffa31ff8097ca05ea9",
		"vcs.modified", "true",
	))
	if !strings.HasSuffix(got.Commit, "-dirty") {
		t.Errorf("commit = %q; want a -dirty suffix on a modified tree", got.Commit)
	}
	if !strings.HasPrefix(got.Commit, "4bfa8682c084") {
		t.Errorf("commit = %q; want the truncated revision before the marker",
			got.Commit)
	}
}

// TestResolve_UsesModuleVersionWhenThereIsNoVCS covers
// `go install module@version`, which embeds a module version and no VCS
// settings at all.
func TestResolve_UsesModuleVersionWhenThereIsNoVCS(t *testing.T) {
	stubLdflags(t, devVersion, noCommit, unknownDate)
	got := resolve(withSettings("v0.4.1"))
	if got.Version != "v0.4.1" {
		t.Errorf("version = %q; want the module version", got.Version)
	}
}

// TestResolve_IgnoresTheDevelPlaceholder pins that "(devel)" — the
// toolchain's local-build placeholder — is not treated as identity. It
// says no more than "dev" does.
func TestResolve_IgnoresTheDevelPlaceholder(t *testing.T) {
	stubLdflags(t, devVersion, noCommit, unknownDate)
	got := resolve(withSettings("(devel)"))
	if got.Version != devVersion {
		t.Errorf("version = %q; want %q — \"(devel)\" carries no identity",
			got.Version, devVersion)
	}
}

// TestResolve_SurvivesAnAbsentBuildInfo pins the `go test` shape, where
// ReadBuildInfo reports ok==false, and `-buildvcs=false` builds.
func TestResolve_SurvivesAnAbsentBuildInfo(t *testing.T) {
	stubLdflags(t, devVersion, noCommit, unknownDate)
	got := resolve(func() (*debug.BuildInfo, bool) { return nil, false })
	want := Info{Version: devVersion, Commit: noCommit, Date: unknownDate}
	if got != want {
		t.Errorf("resolve() = %+v; want the defaults %+v", got, want)
	}
}

// TestString_RendersOneLine pins the human form the CLI prints, which
// respond.TextLiner routes to stdout.
func TestString_RendersOneLine(t *testing.T) {
	got := Info{Version: "v1.0.0", Commit: "abc1234", Date: "2026-01-01T00:00:00Z"}.String()
	want := "v1.0.0 (commit abc1234, built 2026-01-01T00:00:00Z)"
	if got != want {
		t.Errorf("String() = %q; want %q", got, want)
	}
	if strings.Contains(got, "\n") {
		t.Errorf("String() = %q; want a single line", got)
	}
}
