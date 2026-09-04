package cmdbind

// RDR 0026 — Phase 3b adversarial tests.
//
// These are IN-PACKAGE on purpose. The two failure modes below live inside
// `readBounded`'s terminal-condition report and inside `spawn`'s fold of
// the two reports; an end-to-end fixture can only observe them through a
// child that must first push a megabyte or pace a tail for seconds, which
// makes the witness slow and its diagnosis indirect. Driving `readBounded`
// against a real `os.Pipe` with a real held write end reproduces the exact
// state `spawn` reaches — a parent that has armed `now + DrainGrace` on its
// own read end while a writer it cannot reach still holds the other — and
// names which of the three ordered conditions the drain reported.
//
// Nothing here mocks the read. The pipes are real, the deadline is the
// real one the parent arms, and the loop under test is the shipped one.
//
// The third row is not a failure mode. It documents deviation D10's
// ADMITTED liveness residue by pinning the precondition of the mechanism
// whose absence produces it, at the same level and against the same real
// pipes.

import (
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// reportName renders a drainReport for a failure message, so the diagnosis
// names the condition rather than an integer.
func reportName(r drainReport) string {
	switch r {
	case drainWhole:
		return "whole (EOF)"
	case drainHeld:
		return "held (os.ErrDeadlineExceeded on the final read)"
	case drainOverflow:
		return "overflow (past the read bound)"
	default:
		return "unknown"
	}
}

// heldWriter opens a pipe and writes `total` bytes into it from a
// goroutine that then HOLDS the write end open forever — the escapee
// `reapGroup` cannot reach. It returns the read end.
//
// The write end is deliberately never closed, which is the whole point, and
// the fd is released when the test process ends. The writer blocks in
// `write(2)` once the ~64 KiB pipe buffer fills, so the caller must already
// be draining: nothing here waits for the writes to land.
func heldWriter(t *testing.T, total int) *os.File {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	t.Cleanup(func() { _ = r.Close() })

	go func() {
		buf := make([]byte, 64<<10)
		for i := range buf {
			buf[i] = 'x'
		}
		for total > 0 {
			n := len(buf)
			if total < n {
				n = total
			}
			m, werr := w.Write(buf[:n])
			if werr != nil {
				return
			}
			total -= m
		}
		// The write end is NEVER closed: this stands in for the `setsid(2)`
		// grandchild that survives the group kill still holding the pipe.
		select {}
	}()
	return r
}

// TestAdv1_AHeldStderrPastTheReadBoundStillReportsHeldAndStillRefuses is
// ADV-1.
//
// F1 (Failure Modes): "an escaped grandchild holds a pipe past the bound —
// the invocation refuses `timeout` (deadline elapsed) or
// `execution_failure` (deadline not elapsed) with `Err` naming the pipe(s)
// and the bound". C1 `residue:` fixes the same guarantee from the other
// side: "Leakage of a process, or of THAT process's output, is admitted; a
// withheld refusal is not."
//
// The escape: `spawn` drains stderr with `readBounded(errR, StdoutCap)` —
// a limit of exactly `StdoutCap`, not `StdoutCap+1`. `readBounded` applies
// the overflow rank whenever `len(kept) >= limit`, WHATEVER the terminal
// condition, so a verbose child whose stderr reaches 1 MiB and is then held
// by an escapee reports `drainOverflow` rather than `drainHeld`. `spawn`'s
// overflow arm consults only the STDOUT report (`outReport ==
// drainOverflow || len(inv.stdout) > StdoutCap`) — there is no stderr
// overflow refusal in the contract, because stderr is tailed rather than
// capped — and `heldPipes` counts only `drainHeld`. The held condition is
// therefore laundered into a rank nobody consults, and the invocation
// returns with NO error at all: the refusal is withheld, which is the one
// thing C1 `residue:` does not admit.
//
// C1 `precedence:` scopes the overflow rank deliberately: "Overflow
// outranks held on the same drain: a drain past `StdoutCap` refuses the
// existing overflow error whatever its terminal condition, since the bytes
// were read and the pipe's state is then not what the invocation turns on."
// The rank exists BECAUSE overflow refuses. On stderr it does not refuse,
// so ranking a held stderr as overflow discards the condition instead of
// outranking it.
//
// ADVERSARIAL
func TestAdv1_AHeldStderrPastTheReadBoundStillReportsHeldAndStillRefuses(t *testing.T) {
	// The stderr drain's own limit, verbatim from `spawn`.
	const stderrLimit = StdoutCap

	r := heldWriter(t, stderrLimit+(256<<10))

	done := make(chan drainReport, 1)
	go func() {
		_, report := readBounded(r, stderrLimit)
		done <- report
	}()

	// This is exactly what `spawn` does when its one timer fires: it arms
	// `now + DrainGrace` on its OWN read end and joins.
	time.Sleep(200 * time.Millisecond)
	if derr := r.SetReadDeadline(time.Now().Add(DrainGrace * time.Millisecond)); derr != nil {
		t.Fatalf("SetReadDeadline on a real pipe read end: %v", derr)
	}

	var report drainReport
	select {
	case report = <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the stderr drain never returned after the grace was armed")
	}

	if report == drainHeld {
		return
	}
	t.Fatalf("a stderr pipe HELD by an unreachable writer past the drain "+
		"bound reported %s, not held.\n\n"+
		"`spawn` drains stderr with a limit of exactly StdoutCap (%d), so a "+
		"verbose child that reaches the limit and is then held takes "+
		"`readBounded`'s `len(kept) >= limit` overflow rank BEFORE the "+
		"deadline test. `spawn` then folds the reports: its overflow arm "+
		"reads only the STDOUT report, and `heldPipes` counts only the held "+
		"rank — so this condition is consulted by neither arm and the "+
		"invocation returns with NO error.\n\n"+
		"F1 promises the invocation REFUSES with the pipe named; C1 "+
		"`residue:` admits leaking the process and its output but states "+
		"\"a withheld refusal is not\" admitted. C1 `precedence:` gives the "+
		"overflow rank its scope — \"a drain past `StdoutCap` REFUSES the "+
		"existing overflow error\" — and stderr has no overflow refusal to "+
		"take, so ranking a held stderr as overflow discards the condition "+
		"rather than outranking it.",
		reportName(report), StdoutCap)
}

// TestAdv2_APacedTailReachingARealEOFIsNeverReportedHeld is ADV-2.
//
// F4 (Failure Modes): "bytes the direct child wrote can be lost at the
// bound; the refusal still fires, so the loss is bounded to a refused
// invocation, never a parsed one." F4 scopes the loss to a pipe an escapee
// still holds. This test's pipe has NO holder: the writer closes, EOF
// arrives, and every byte is deliverable.
//
// The escape: `readBounded`'s liveness ceiling is checked at the TOP of the
// loop, before the read, and expires `2·WaitDelay` after the drain entered
// the grace. It is armed once and never re-armed by progress. A writer
// whose per-read idle gaps are all UNDER `DrainGrace` — so the grace itself
// never fires — but whose tail takes longer than `2·WaitDelay` to finish is
// cut off mid-tail and reported HELD, on a pipe that reaches a real EOF and
// has no writer at all.
//
// That is the failure C1 `precedence:` forbids by name: "the grace bounds
// the IDLE GAP between reads and never the size or total duration of the
// tail" and "A single absolute deadline would instead bound the whole
// remaining tail against the clock, which turns a slow-arriving tail into a
// false held-pipe report on a pipe with no writer at all". The ceiling is
// that single absolute deadline, moved out to `2·WaitDelay`; the mechanism
// the clause rejects is unchanged, only its constant is larger. A9 measured
// the shipped re-arm recovering a 1 MiB paced tail in 612 ms (darwin) /
// 663 ms (linux), so the ceiling sits at roughly 1.5x the ALREADY MEASURED
// case — a loaded host or a slower writer cadence reaches it.
//
// The pacing here is deliberately inside the grace on every gap (40 ms
// against a 50 ms `DrainGrace`), so a failure cannot be blamed on a gap the
// grace legitimately refuses: the ONLY thing that ends this drain early is
// the total-duration ceiling.
//
// ADVERSARIAL
func TestAdv2_APacedTailReachingARealEOFIsNeverReportedHeld(t *testing.T) {
	const (
		chunks    = 40
		chunkSize = 4 << 10
		gap       = 40 * time.Millisecond // strictly under DrainGrace (50 ms)
		want      = chunks * chunkSize
	)

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	t.Cleanup(func() { _ = r.Close() })

	go func() {
		buf := make([]byte, chunkSize)
		for i := range buf {
			buf[i] = 'y'
		}
		for i := range chunks {
			if i > 0 {
				time.Sleep(gap)
			}
			if _, werr := w.Write(buf); werr != nil {
				break
			}
		}
		// EVERY writer is gone: the drain's next read is a true EOF.
		_ = w.Close()
	}()

	type outcome struct {
		kept   int
		report drainReport
	}
	done := make(chan outcome, 1)
	go func() {
		kept, report := readBounded(r, StdoutCap+1)
		done <- outcome{len(kept), report}
	}()

	// The parent's bound fires while the writer is still pacing, which is
	// the ordinary late-drain shape (MVV row 6(b), S9).
	time.Sleep(200 * time.Millisecond)
	if derr := r.SetReadDeadline(time.Now().Add(DrainGrace * time.Millisecond)); derr != nil {
		t.Fatalf("SetReadDeadline on a real pipe read end: %v", derr)
	}

	var got outcome
	select {
	case got = <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("the paced drain never returned")
	}

	if got.report == drainWhole && got.kept == want {
		return
	}
	t.Fatalf("a paced tail with NO holder — every idle gap %v, strictly "+
		"under the %d ms grace, and the writer closing for a real EOF — "+
		"reported %s with %d of %d bytes recovered.\n\n"+
		"The liveness ceiling is checked at the top of the loop and expires "+
		"2*WaitDelay (%v) after the drain entered the grace. It is armed "+
		"once and progress never re-arms it, so a tail whose TOTAL duration "+
		"exceeds it is cut off and reported held even though the next read "+
		"would have returned EOF.\n\n"+
		"C1 `precedence:` states the grace \"bounds the IDLE GAP between "+
		"reads and never the size or TOTAL DURATION of the tail\", and "+
		"rejects the rival mechanism by name: \"A single absolute deadline "+
		"would instead bound the whole remaining tail against the clock, "+
		"which turns a slow-arriving tail into a false held-pipe report on "+
		"a pipe with no writer at all\". The ceiling IS that single "+
		"absolute deadline with a larger constant. A9 measured the shipped "+
		"re-arm recovering a 1 MiB paced tail in 612-663 ms, so the ceiling "+
		"sits at ~1.5x an already-measured legitimate case.",
		gap, DrainGrace, reportName(got.report), got.kept, want,
		2*WaitDelay*time.Millisecond)
}

// TestD10_BothPipesTricklingRaisesNoHalt_DocumentedResidue documents an
// ADMITTED liveness residue. It is not a bug report and it does not assert
// a hang: it pins the PRECONDITION of the mechanism whose absence produces
// the residue, so the residue cannot be forgotten and cannot be silently
// closed or silently widened.
//
// The residue (deviation D10): when an escapee trickles BOTH pipes with
// every idle gap under `DrainGrace`, neither drain ever reaches a refusing
// condition. The sibling halt is raised ONLY by a drain that has itself
// reported one, so nothing raises it, and `spawn`'s join — which `0026:C1`
// `precedence:` (b) requires it to complete before reading either byte
// slice — stays unbounded. That is the old hang, in a shape no fixture
// exercises.
//
// It is ADMITTED rather than closed because both mechanisms that would
// close it are already refused by this record. A duration bound on the
// drain is the total-duration bound `0026:C1` `precedence:` rejects by name
// ("the grace bounds the IDLE GAP between reads and never the size or total
// duration of the tail"), and is exactly what the removed liveness ceiling
// was (ADV-2). A hard close of the read ends is Go's answer
// (`os/exec`'s `awaitGoroutines` -> `closeDescriptors`) and surrenders the
// stderr tail naming WHICH helper held the pipe — the sole justification
// for this record's deliberate divergence from that prior art. So the
// residue takes F5's shape: recorded at the check rather than discovered at
// the join.
//
// WHAT IS ASSERTED, and why it is not a hang. Reproducing the residue
// end-to-end means waiting forever, which is not a test. What IS testable
// in bounded time is the mechanism's precondition, in three parts:
//
//  1. A drain trickling UNDER the grace does not raise the halt. Observed
//     over a window many multiples of the grace wide, `stop` stays false —
//     so a second drain in the same state has nothing to end it either, and
//     the join has nothing to wait on but EOF that never comes.
//  2. That drain has also not returned on its own. Both parts together are
//     the residue: no report, therefore no halt, therefore no terminator.
//  3. The halt is what would have ended it. Raising `stop` by hand — which
//     stands in for the sibling report that in this shape never happens —
//     returns the drain promptly, and it reports held. This is the control:
//     it proves part 2 is the ABSENCE of a raised halt and not some other
//     defect in the loop, and it lets the test finish cleanly instead of
//     leaking a wedged goroutine.
//
// If this test ever fails at part 1 or 2, the residue has been closed by
// some new mechanism and D10 needs revisiting — including whether that
// mechanism is one this record permits.
//
// ADVERSARIAL
func TestD10_BothPipesTricklingRaisesNoHalt_DocumentedResidue(t *testing.T) {
	// Strictly under DrainGrace (50 ms), so the drain keeps its re-armed
	// deadline ahead of it forever — the shape D10 names.
	const gap = 10 * time.Millisecond
	// Many multiples of the grace: long enough that a drain which was ever
	// going to report has done so, short enough to keep the row fast. The
	// D10 shape trickles indefinitely, so no window proves "forever"; this
	// one proves the mechanism does not fire on the timescale the grace and
	// the join's own bound operate on.
	const window = 20 * DrainGrace * time.Millisecond

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	t.Cleanup(func() { _ = r.Close() })

	// The trickler stands in for the escapee `reapGroup` cannot reach: it
	// writes forever and NEVER closes, so no read ever returns EOF. It
	// stops when the test's pipe is closed at cleanup.
	go func() {
		for {
			if _, werr := w.Write([]byte("x")); werr != nil {
				return
			}
			time.Sleep(gap)
		}
	}()

	// One `stop`, exactly as `spawn` shares it between its two drains. Only
	// one drain runs here: the second would be in the identical state, and
	// running it would double the goroutines without changing what either
	// reports. What matters is that NOTHING sets this flag.
	var stop atomic.Bool
	done := make(chan drainReport, 1)
	go func() {
		_, report := readBounded(r, StdoutCap+1, &stop)
		done <- report
	}()

	// The parent's mark: its one timer fired and it armed `now + DrainGrace`
	// on its own read end. From here the drain re-arms the grace itself
	// before every read.
	time.Sleep(100 * time.Millisecond)
	if derr := r.SetReadDeadline(time.Now().Add(DrainGrace * time.Millisecond)); derr != nil {
		t.Fatalf("SetReadDeadline on a real pipe read end: %v", derr)
	}

	// Parts 1 and 2: over the window, no report and no halt.
	select {
	case report := <-done:
		t.Fatalf("a drain trickling every %v — under the %d ms grace — "+
			"reported %s within %v.\n\n"+
			"D10 records the both-pipes trickle as an ADMITTED liveness "+
			"residue on the premise that such a drain reaches NO refusing "+
			"condition of its own. If it now reports, some mechanism ends "+
			"it that did not before, and D10's premise no longer holds. "+
			"Check what that mechanism is: `0026:C1` `precedence:` rejects "+
			"a total-duration bound by name, and `0026:C1` `bound:` states "+
			"`WaitDelay` is the ONE bound.",
			gap, DrainGrace, reportName(report), window)
	case <-time.After(window):
	}
	if stop.Load() {
		t.Fatalf("the sibling halt was raised after %v with no drain having "+
			"reported a refusing condition.\n\n"+
			"The halt's whole discipline (D9) is that it is raised ONLY by a "+
			"drain that has itself reported held or overflow, and never on "+
			"elapsed time — that is what keeps it a state predicate rather "+
			"than the second bound `0026:C1` `bound:` forbids. A halt raised "+
			"here would mean something reads a clock.", window)
	}

	// Part 3, the control: the halt is what would have ended it. Raising
	// `stop` by hand stands in for the sibling report that, in the D10
	// shape, never happens.
	stop.Store(true)
	select {
	case report := <-done:
		if report != drainHeld {
			t.Fatalf("a trickling drain ended by the sibling halt reported "+
				"%s, not held. It did not reach EOF, so held is the "+
				"condition it has (`0026:C1` `precedence:`).",
				reportName(report))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the trickling drain did not return after the sibling halt " +
			"was raised. The halt is the ONLY terminator this shape has, so " +
			"if it does not end the drain then D10 is not a residue but an " +
			"unconditional hang.")
	}
}
