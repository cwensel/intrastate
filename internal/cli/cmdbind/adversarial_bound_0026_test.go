package cmdbind_test

// RDR 0026 — Phase 3b adversarial test ADV-3.
//
// This one drives a REAL child through the real `cmdbind.Reader.Read` and
// asserts nothing but the published number: C1 `bound:`'s committed
// per-invocation return of `timeout + 2·WaitDelay`, with the record's own
// stated scheduling tolerance of +100 ms.

import (
	"context"
	"testing"
	"time"
)

// TestAdv3_ATricklingEscapeeStillReturnsInsideTheCommittedTotal is ADV-3.
//
// F1 (Failure Modes) promises the trickling escapee is REFUSED, and the
// existing trickle oracle asserts exactly that. What no oracle asserts is
// WHEN. C1 `bound:` publishes the number as a commitment, not a hope:
// "Per invocation, the COMMITTED return bound is timeout + 2·WaitDelay,
// asserted with a scheduling tolerance of +100 ms — stated here once so
// the scenarios cite it rather than each carrying its own: a return past
// the total but inside the tolerance is a pass, and past the tolerance is a
// contract violation." F2's consequence for a caller is what makes it
// load-bearing: an agent sizing its own watchdog from `timeout` and the
// published bound is killed mid-invocation when the real return runs long.
//
// The escape: the liveness ceiling is `2·WaitDelay` measured from the
// DRAIN's entry into the grace, not carved out of the parent's allowance.
// The parent has already spent its full `WaitDelay` waiting on the join
// before it arms the grace at all, so on a trickler the two costs STACK:
//
//	WaitDelay (join timer) + DrainGrace (first grace window) + 2·WaitDelay
//	(the ceiling) ≈ 1550 ms
//
// against a committed drain allowance of 2·WaitDelay = 1000 ms. That is
// ~460 ms past the stated tolerance. The trickler is the ONE shape the
// ceiling exists for, so this is not an exotic corner: it is the ceiling's
// own motivating case overrunning the record's published number.
//
// The declared timeout is long (10 s) and the child exits immediately, so
// leg (a) contributes nothing and the elapsed time measured here is very
// nearly the drain leg alone. That is what makes the assertion a clean read
// on the drain's cost rather than on the deadline's.
//
// ADVERSARIAL
func TestAdv3_ATricklingEscapeeStillReturnsInsideTheCommittedTotal(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The direct child exits 0 at once; the escapee drips a byte every
	// `dripGap` (100 ms, twice `DrainGrace`) forever, so the drain's own
	// re-armed deadline is what ends it — one join timer plus at most one
	// gap plus one grace. The cadence was 10 ms when this row was written,
	// which is a fifth of the grace and never fires it; repaired as D8.
	r := escReader(t, "10s", "trickle", pidfile)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("the trickling escapee was accepted; F1 owes a refusal")
	}

	// The child exits immediately, so leg (a) costs ~nothing and the
	// committed total for THIS invocation is the drain allowance plus the
	// record's stated tolerance.
	budget := escBoundFor(0)
	if elapsed <= budget {
		return
	}
	t.Fatalf("the trickling escapee refused after %v, past the committed "+
		"total of %v (2*WaitDelay = %v, plus the record's stated %v "+
		"scheduling tolerance).\n\n"+
		"The direct child exits 0 immediately, so leg (a) of the bound "+
		"costs ~nothing and this elapsed time is the DRAIN leg alone.\n\n"+
		"The liveness ceiling is 2*WaitDelay measured from the drain's "+
		"entry into the grace rather than carved out of the parent's "+
		"allowance, so it STACKS on the full WaitDelay the parent already "+
		"spent on the join timer: WaitDelay + DrainGrace + 2*WaitDelay.\n\n"+
		"C1 `bound:` publishes timeout + 2*WaitDelay as a COMMITMENT — \"a "+
		"return past the total but inside the tolerance is a pass, and past "+
		"the tolerance is a contract violation\" — and F2's consequence is "+
		"a caller that sized its watchdog from the published number and is "+
		"killed mid-invocation. The trickler is the ceiling's own motivating "+
		"case, so this is the ceiling overrunning the number it was added "+
		"underneath.",
		elapsed, budget, 2*escWaitDelay, escTolerance)
}
