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
