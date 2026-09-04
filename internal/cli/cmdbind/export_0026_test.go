package cmdbind

import (
	"testing"
	"time"
)

// RDR 0026 — the ONE drain-start stall seam.
//
// REQ-128 (S7): "the stall is ONE seam, shared by this row, S9 and MVV
// row 6 so they cannot build two: a test-only delay hook on the drain
// goroutine's start ... one unexported package-level hook set by the test
// and nil in production (not a build tag, not a sleep in production code).
// S7 is a standing regression trap, so the seam outlives the fixture"
// REQ-135 (XC): "A12 confirms the drain-start stall seam this RDR adds
// needs **no third build tag** — it is portable Go in `cmdbind`."
//
// The seam is TEST-ONLY by construction, exactly like `0025:C4`'s `goos`
// override: `export_test.go`-family files compile only into this package's
// test binary, so the hook stays out of the shipped surface. Exactly ONE
// such hook ships and S7, S9 and MVV row 6 all reuse it.
//
// This file is UNTAGGED and so is the hook it exposes — A12's finding that
// the seam needs no third build tag is asserted by this file compiling on
// every platform the package does.
func SetDrainStartDelayForTest(t *testing.T, d time.Duration) {
	t.Helper()

	prev := drainStartDelay
	drainStartDelay = d
	t.Cleanup(func() { drainStartDelay = prev })
}

// DrainStartDelayIsNilInProduction reports the hook's PRODUCTION value —
// the zero duration, i.e. no delay at all. REQ-128 states the hook is "nil
// in production", which the drift guard reads back here rather than by
// inspecting the shipped source.
func DrainStartDelayIsNilInProduction() bool { return drainStartDelay == 0 }

// SetNonPollableReadEndsForTest makes the creation-time pollability probe
// see a read end that refuses a deadline, so S5's condition is SIMULATED at
// the check rather than reproduced.
//
// REQ-126 (S5): "the creation-time pollability check returns
// `os.ErrNoDeadline` (F5), simulated at the check by injecting a read end
// that refuses a deadline (`os.NewFile` over a dup'd pipe fd, the shape F5
// confirms produces the error)"
// REQ-107 (F5): the trigger is "unexercised for a real `os.Pipe`" on both
// supported OSes, so a fixture that reproduced it would be fabricating a
// host condition; the check's own INPUT is what the test controls.
func SetNonPollableReadEndsForTest(t *testing.T) {
	t.Helper()

	prev := nonPollableForTest
	nonPollableForTest = true
	t.Cleanup(func() { nonPollableForTest = prev })
}
