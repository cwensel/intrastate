package cmdbind

import "testing"

// SetGOOSForTest overrides the platform the binding checks, restoring it
// when the test ends.
//
// The seam is TEST-ONLY by construction: `export_test.go` is compiled only
// into this package's test binary, so the injectable `goos` var
// `0025:C4` calls for stays out of the shipped public surface. Unlike C4's
// deadline triple — which carries no injection seam, because a way to
// disable it would be a way to weaken the deadline at runtime — overriding
// this one weakens nothing: the refusal it drives is a platform check, and
// forcing it ON is the only thing a test does with it.
func SetGOOSForTest(t *testing.T, value string) {
	t.Helper()

	prev := goos
	goos = value
	t.Cleanup(func() { goos = prev })
}

// ProcGroupPlatforms is the platform set the process-group mechanism's
// build tags name, exported to this package's test binary so the
// drift-guard can compare it against `Unsupported()`'s refuse-list.
var ProcGroupPlatforms = procGroupPlatforms

// ProcGroupSupported reports whether the syscall-bearing half of the
// mechanism is the one compiled into THIS binary.
const ProcGroupSupported = procGroupSupported

// UnsupportedOn reports `Unsupported()`'s verdict for an arbitrary
// platform string, without disturbing the package's `goos` var. The
// drift-guard needs the PREDICATE, not the running platform's answer.
func UnsupportedOn(value string) bool {
	prev := goos
	goos = value
	defer func() { goos = prev }()

	return Unsupported()
}
