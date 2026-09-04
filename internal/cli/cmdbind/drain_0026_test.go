package cmdbind_test

// RDR 0026 — C1 `precedence:` and C1 `bound:`: a bounded wait outranks a
// whole read.
//
// Every scenario here drives a REAL child through the real
// `cmdbind.Reader.Read`. Nothing mocks `spawn`, no private field is read,
// and no test names the report type, the fold, or the grace's identifier —
// C1 `precedence:` explicitly leaves "the identifier and the Go spelling
// [to] the implementer" (REQ-13) and the Illustrative Code names
// `heldPipes`/`errHeldPipe` as "names for the shape, not the contract"
// (REQ-65). What is asserted is what the contract fixes: which invocations
// RETURN, WHEN they return, WHAT they refuse, and WHAT SURVIVES into the
// refusal.

import (
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/cli/cmdbind"
)

// --- C1 `precedence:` — the bound is delivered ---------------------------

// REQ-1: "precedence: a bounded wait outranks a whole read."
// REQ-19: "Only a held pipe refuses."
// REQ-58 (AP): "the refusal is never withheld, and the invocation returns
// within `timeout + 2·WaitDelay`."
// REQ-44 (C1 `bound:`): "Per invocation, the COMMITTED return bound is
// timeout + 2·WaitDelay, asserted with a scheduling tolerance of +100 ms ...
// a return past the total but inside the tolerance is a pass, and past the
// tolerance is a contract violation."
// REQ-53 (C1 `residue:`): "a process outside the child's process group may
// outlive the refusal"
// REQ-54: "Leakage of a process, or of THAT process's output, is admitted;
// a withheld refusal is not."
// ADVERSARIAL
func TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	r := escReader(t, escTimeout, "deadline", pidfile)

	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatalf("a read whose pipes are held by a `setsid(2)` escapee returned "+
			"values=%+v unreadable=%v with NO refusal after %v; C1 "+
			"`precedence:` fixes that a bounded wait outranks a whole read, "+
			"and C1 `residue:` admits leakage of a process, never a withheld "+
			"refusal", values, unreadable, elapsed)
	}
	if want := escBoundFor(timeout); elapsed > want {
		t.Fatalf("the invocation returned after %v, past the C1 `bound:` "+
			"committed total of timeout + 2·WaitDelay (%v) plus its stated "+
			"%v scheduling tolerance: a return past the tolerance is a "+
			"contract violation", elapsed, timeout+2*escWaitDelay, escTolerance)
	}
	// C1 `residue:` — the escapee OUTLIVES the refusal. That is the admitted
	// leak, and asserting it is what keeps the refusal honest: a suite in
	// which the escapee died would be testing a reachable writer.
	if !escAlive(pid) {
		t.Fatalf("the `setsid(2)` grandchild pid %d did not outlive the "+
			"refusal; C1 `residue:` names it as the admitted leak and this "+
			"fixture is only exercising the held-pipe path while it lives",
			pid)
	}
}

// REQ-4 (C1 `precedence:`): "each drain's FINAL read then reports the
// pipe's state, not the clock: bytes still buffered are delivered, EOF
// (every writer gone) is whole output as today, and a deadline error is a
// HELD pipe."
// REQ-5: "The predicate is the drain's REPORTED terminal condition, not the
// byte count: a drain that ends on `os.ErrDeadlineExceeded` reports held
// however many bytes it delivered during the grace, and one that ends on
// EOF reports whole."
// REQ-6: "\"Empty, with a writer the group signal could not reach\"
// describes the common case; it is not the test, because a writer trickling
// bytes slower than the grace is exactly the case that must refuse (S9's
// companion) and is never empty."
// REQ-105 (MC `oracle`, S9 negative control): "companion: escaped writer
// emitting a byte [on a cadence SLOWER than `DrainGrace`] is REFUSED"
// REQ-131 (S9 companion): "an escaped writer emitting a byte [on a cadence
// SLOWER than `DrainGrace`] forever is REFUSED (never blocked, never EOF)"
//
// Both requirements are quoted with the cadence REPAIRED, as deviation D8
// resolves. The record's literal figure is "a byte every 10 ms", which is a
// FIFTH of the 50 ms grace: under C1 `precedence:`'s own predicate such a
// writer is MAKING PROGRESS, keeps its deadline ahead of it and is never
// reported held — the fixture and the fence selected opposite outcomes on
// the same input. C1's prose names the discriminating case in the other
// direction ("a writer trickling bytes SLOWER than the grace is exactly the
// case that must refuse"), and the fence outranks the scenario (the idiom
// D3 applied at Phase 0). The fixture's gap is therefore `dripGap` = 100 ms,
// twice the grace; see `dripGap` in fixtures_0026_test.go for the margin
// arithmetic.
//
// The repair also changes WHY the row is green. At 10 ms the stdout drain
// never reported on its own: the fixture writes to stdout only, so its idle
// stderr reported held under the grace and the sibling halt is what ended
// the stdout drain. At 100 ms the stdout drain reaches held on its OWN
// mechanism — an idle gap longer than the grace — which is the mechanism the
// row's own text claims to witness.
// REQ-18: "refuse on \"no EOF after a final bounded read\", not on elapsed
// time"
// ADVERSARIAL
func TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The direct child closes its own stderr and exits 0 immediately; the
	// escapee inherits STDOUT ONLY and drips a byte every `dripGap` (100 ms,
	// twice `DrainGrace`) forever. The stdout drain is NEVER empty and NEVER
	// blocked — only its reported terminal condition can refuse it, and at
	// this cadence that condition is the drain's own idle-gap deadline.
	//
	// The escapee holds ONE pipe, which is what makes the row discriminate.
	// Under the both-pipes `trickle` shape the escapee also inherits an idle
	// stderr: that pipe reports held under the grace alone and raises the
	// sibling halt, so stdout was named held by the HALT rather than by its
	// own mechanism, and this row passed for a reason its text did not name.
	// Mutating `readBounded` so a drain that has READ BYTES never reports
	// held on its own idle gap — the exact regression S9's companion exists
	// to catch — left the row green and the refusal text byte-identical.
	// With stderr reaching EOF, "stdout" can only be named by the stdout
	// drain's OWN idle-gap report (deviation D8's recommended amendment).
	r := escReader(t, "10s", "trickle-stdout", pidfile)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatalf("a trickling escaped writer (a byte every 100 ms, never "+
			"blocked, never EOF) was ACCEPTED after %v: values=%+v "+
			"unreadable=%v. C1 `precedence:` states that \"empty, with a "+
			"writer the group signal could not reach\" is the common case "+
			"and NOT the test — a writer trickling slower than the grace is "+
			"exactly the case that must refuse", elapsed, values, unreadable)
	}
	// And it refuses on the drain's reported condition, not on the declared
	// deadline: the deadline here is 10 s and the invocation must return in
	// roughly one bound past the child's own immediate exit.
	if elapsed >= 10*time.Second {
		t.Fatalf("the trickle case took %v — its whole declared deadline. "+
			"C1 `precedence:` refuses on \"no EOF after a final bounded "+
			"read\", not on elapsed time (P-2), so the refusal is owed one "+
			"bound past the child's exit and not the deadline", elapsed)
	}
	// And it names STDOUT — the trickling pipe — which is the assertion that
	// makes this row discriminate. The child's stderr reaches EOF, so no
	// sibling can report held and no halt can be raised: the only way
	// "stdout" appears in the refusal is the stdout drain reporting its own
	// idle gap, which is the mechanism REQ-6 and REQ-131 name. Asserting
	// only `err != nil` and an elapsed bound does not distinguish that
	// mechanism from a sibling naming it, which is what deviation D8's
	// recommended amendment requires here — refused, NAMED PIPE, inside the
	// committed bound.
	if detail := escHeldDetail(t, err); !escNamesHeld(detail, "stdout") {
		t.Fatalf("the refusal reads %q and does not name STDOUT as held. "+
			"The trickling escapee holds stdout alone and the child's "+
			"stderr reached EOF, so the only mechanism that can name it is "+
			"the stdout drain's own idle-gap report — REQ-6's \"a writer "+
			"trickling bytes slower than the grace is exactly the case that "+
			"must refuse\", and REQ-16's \"only a gap longer than the grace "+
			"ends it as held\"", detail)
	}
}

// REQ-9 (C1 `precedence:`): "The grace is what MAKES that final read
// happen: Go checks the deadline BEFORE attempting the syscall
// (`internal/poll/fd_unix.go::(*FD).Read` calls `prepareRead` ahead of
// `syscall.Read`), so a deadline of NOW short-circuits with n=0 and reports
// a held pipe without ever looking at the pipe."
// REQ-16: "a drain making progress keeps its deadline ahead of it and runs
// to EOF, and only a gap longer than the grace ends it as held."
// REQ-76 (LBD selection): "when EOF and the bound both qualify: EOF wins if
// it arrives first (whole output, as today); the bound wins otherwise and
// the invocation refuses."
// REQ-131 (S9): "a LATE drain with no holder ... Expected: accepted, whole
// output, no refusal — the final read under `DrainGrace` delivers the
// buffered bytes and then sees EOF."
// REQ-MVV.6a (burst): "a child that writes 1 MiB to stdout in one `Write`
// and exits with no grandchild, run with the drain goroutine artificially
// delayed past the timer, returns the whole value ... This shape guards C1
// `precedence:` against a deadline of literally `now`, which short-circuits
// before the syscall and returns n=0"
// BOUNDARY
func TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue(t *testing.T) {
	const payload = 1 << 20 // 1 MiB, MVV row 6(a)

	// The drain goroutine is artificially delayed PAST the join's timer, so
	// the bytes are already resident in the pipe and every writer is gone
	// when the final read runs. Only a FUTURE deadline runs that read; a
	// deadline of `now` short-circuits with n=0 and loses them.
	const stall = 2 * time.Second
	cmdbind.SetDrainStartDelayForTest(t, stall)

	r := escReader(t, "10s", "burst", strconv.Itoa(payload))
	r.Accessor.Output = escRaw()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	// The seam must actually STALL the drain, or this row is tautological:
	// a build that ignores the hook drains promptly, never reaches the
	// timer, and passes without exercising the final read at all.
	if elapsed < stall {
		t.Fatalf("the invocation returned in %v with the drain-start stall "+
			"set to %v: the drain was never delayed past the join's timer, so "+
			"the final read under `DrainGrace` was never exercised. S7 fixes "+
			"the stall as ONE test-only hook this row, S9 and MVV row 6 all "+
			"reuse", elapsed, stall)
	}
	if err != nil {
		t.Fatalf("a late burst drain with NO holder refused: %v. C1 "+
			"`precedence:` and S9 both fix this as accepted, whole output, "+
			"no refusal — the final read under `DrainGrace` delivers the "+
			"buffered bytes and then sees EOF", err)
	}
	if len(unreadable) != 0 || len(values) != 1 {
		t.Fatalf("values=%+v unreadable=%v; want the one declared key with "+
			"the whole payload", values, unreadable)
	}
	if got := len(values[0].Value); got != payload {
		t.Fatalf("the late burst drain recovered %d of %d bytes. A deadline "+
			"of literally `now` short-circuits before the syscall and "+
			"returns n=0 (A2, both OSes); the grace is what MAKES the final "+
			"read happen", got, payload)
	}
}

// REQ-10 (C1 `precedence:`): "The deadline is re-armed at `now +
// DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the
// IDLE GAP between reads and never the size or total duration of the tail."
// REQ-17: "A single absolute deadline would instead bound the whole
// remaining tail against the clock, which turns a slow-arriving tail into a
// false held-pipe report on a pipe with no writer at all"
// REQ-67 (IC): "the re-arm lives INSIDE that read loop, not in the parent.
// A one-shot `SetReadDeadline` from the parent before `drains.Wait()` would
// bound the whole remaining tail against the clock"
// REQ-90 (MC `fidelity`): "Final drain after the timer, writer gone" →
// "whole remaining tail, however many reads it takes — the grace bounds the
// idle gap between reads, not the tail"
// REQ-MVV.6b: "**Paced — the mechanism discriminator** ... Every chunk is
// recovered and the drain ends on EOF, because the deadline is re-armed
// before each read and so bounds only the idle gap. An implementation that
// instead sets ONE absolute `now + DrainGrace` for the whole final drain
// truncates here"
// REQ-MVV.6c: "Shape (a) does NOT distinguish re-armed from
// single-absolute and must never be the only witness"
// REQ-102 (MC `oracle`): "6(b) is what fails the single-absolute variant"
// ADVERSARIAL
func TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead(t *testing.T) {
	const (
		payload = 1 << 20 // 1 MiB
		chunks  = 32
	)
	// 0.4·DrainGrace per gap, so the tail runs ~12·DrainGrace in total —
	// far past any single absolute `now + DrainGrace`, while each idle gap
	// stays comfortably inside the grace.
	gap := 2 * time.Duration(cmdbind.DrainGrace) * time.Millisecond / 5

	const stall = 2 * time.Second
	cmdbind.SetDrainStartDelayForTest(t, stall)

	r := escReader(t, "20s", "paced", strconv.Itoa(payload),
		strconv.Itoa(chunks), gap.String())
	r.Accessor.Output = escRaw()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	start := time.Now()
	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	if elapsed < stall {
		t.Fatalf("the invocation returned in %v with the drain-start stall "+
			"set to %v: the drain was never delayed past the join's timer, so "+
			"the FINAL drain — the only place the per-read re-arm is the "+
			"discriminating mechanism — was never exercised", elapsed, stall)
	}
	if err != nil {
		t.Fatalf("a PACED late tail with no holder refused: %v. The tail ran "+
			"~%v — about 12·DrainGrace — with idle gaps of %v, each well "+
			"inside the grace. C1 `precedence:` re-arms the deadline BEFORE "+
			"EACH read, so the grace bounds the idle GAP and never the "+
			"tail's size or total duration", err, time.Duration(chunks)*gap, gap)
	}
	if len(unreadable) != 0 || len(values) != 1 {
		t.Fatalf("values=%+v unreadable=%v; want the one declared key",
			values, unreadable)
	}
	if got := len(values[0].Value); got != payload {
		t.Fatalf("the paced tail recovered %d of %d bytes (%.1f%%). A single "+
			"absolute `now + DrainGrace` for the whole final drain truncates "+
			"here — measured at 15.6%% (A9) — and FAILS this row; only a "+
			"deadline re-armed before each read recovers it whole",
			got, payload, 100*float64(got)/float64(payload))
	}
}

// --- C1 `precedence:` — the ordering inside `spawn` ----------------------

// REQ-2 (C1 `precedence:`): "Order in `spawn` is fixed: Wait (direct child
// reaped) → reapGroup (one kill(2) to -pgid, non-blocking) → drain join
// under ONE timer of WaitDelay covering BOTH read drains."
// REQ-59 (TD): "`reapGroup` as today — release first, so the common
// backgrounded-grandchild case closes its ends before any timer matters."
// REQ-123 (S2): "The existing in-group grandchild-leak and cancel oracles
// re-run under the bounded join ... Expected: unchanged results, elapsed
// far below the bound"
// REQ-106 (MC `oracle`, S2 control): S2 "passes by absence-of-change, so it
// is paired with S7's stalled arm"
// HAPPY PATH
func TestReq2_ReleaseFirstKeepsTheInGroupGrandchildCaseFarBelowTheBound(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// An IN-GROUP grandchild holding stdout: reachable by reapGroup, so the
	// release-first order closes its ends before any timer matters.
	r := escReader(t, "10s", "ingroup", pidfile, `{"state.phase":"review"}`)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("the in-group grandchild oracle refused under the bounded "+
			"join after %v: %v. S2 fixes UNCHANGED results — the prior "+
			"point-fixes stay proven and no correct tool is newly refused",
			elapsed, err)
	}
	if len(unreadable) != 0 || len(values) != 1 || values[0].Value != "review" {
		t.Fatalf("values=%+v unreadable=%v; want the whole envelope",
			values, unreadable)
	}
	// "elapsed far below the bound" — release-first means the timer never
	// matters here, so a full WaitDelay spent is the regression S2 names.
	if elapsed >= escWaitDelay {
		t.Fatalf("the in-group case took %v, at or past one WaitDelay (%v): "+
			"C1 `precedence:` releases the group FIRST so the common "+
			"backgrounded-grandchild case closes its ends before any timer "+
			"matters (S2, A5's worst case was 27.7 ms)", elapsed, escWaitDelay)
	}
	pid := escReadPID(t, pidfile)
	defer escKill(pid)
	// The in-group grandchild is REACHED — that is what distinguishes it
	// from C1 `residue:`'s escapee.
	time.Sleep(200 * time.Millisecond)
	if escAlive(pid) {
		t.Fatalf("the IN-GROUP grandchild pid %d survived; reapGroup's "+
			"group-scoped kill(2) reaches it, unlike the `setsid(2)` "+
			"escapee C1 `residue:` admits", pid)
	}
}

// REQ-21 (C1 `precedence:` (b)): "the parent still MUST NOT read
// `inv.stdout`/`inv.stderr` until every drain goroutine has returned, so
// the join remains a real join — the bound governs how long the parent
// waits before setting the deadline, never whether it waits for the
// goroutines to finish."
// REQ-22: "A drain that ends on the deadline still runs `drains.Done()`, so
// the join completes; a design in which the parent proceeds past a live
// drain is forbidden here by name."
// REQ-23 / REQ-134 (XC): "S3 and S7 run under `-race` and that is the
// check ... S3 + S7 under `-race -count=25` — is a Phase-1 exit condition"
// ADVERSARIAL
func TestReq21_TheJoinRemainsARealJoinWithAPartialWriteInFlight(t *testing.T) {
	// The paced fixture keeps a drain WRITING to its slice while the parent
	// is past the bound: the MVV's own fixture cannot catch a violation of
	// (b) because its child is silent and the drain is blocked in its FIRST
	// read with no partial write in flight (REQ-23). Under `-race` a parent
	// that proceeds past a live drain is a reported data race; without it,
	// the observable is a short or torn value.
	const (
		payload = 1 << 18
		chunks  = 16
	)
	gap := 2 * time.Duration(cmdbind.DrainGrace) * time.Millisecond / 5
	const stall = 2 * time.Second
	cmdbind.SetDrainStartDelayForTest(t, stall)

	r := escReader(t, "20s", "paced", strconv.Itoa(payload),
		strconv.Itoa(chunks), gap.String())
	r.Accessor.Output = escRaw()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	start := time.Now()
	values, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	if elapsed := time.Since(start); elapsed < stall {
		t.Fatalf("the invocation returned in %v with the stall set to %v: no "+
			"drain was live past the bound, so (b) was never under test",
			elapsed, stall)
	}
	if err != nil {
		t.Fatalf("the in-flight-write shape refused: %v", err)
	}
	if len(values) != 1 || len(values[0].Value) != payload {
		got := 0
		if len(values) == 1 {
			got = len(values[0].Value)
		}
		t.Fatalf("the value is %d of %d bytes: the parent read the slice "+
			"before every drain goroutine returned. C1 `precedence:` (b) "+
			"requires the join stay a REAL join — the bound governs how long "+
			"the parent waits before setting the deadline, never whether it "+
			"waits for the goroutines to finish", got, payload)
	}
}

// --- C1 `precedence:` — overflow outranks held ---------------------------

// REQ-7 (C1 `precedence:`): "Overflow outranks held on the same drain: a
// drain past `StdoutCap` refuses the existing overflow error whatever its
// terminal condition, since the bytes were read and the pipe's state is
// then not what the invocation turns on."
// REQ-78 (LBD selection): "When overflow and held both qualify on one drain
// ... overflow wins, and the existing \"stdout exceeded the 1 MiB bound\"
// refusal after the join (`cmdbind.go:326-329`) stands unchanged."
// REQ-82 (MC `authority`): "overflow outranks held ... which is why the
// report is one ordered value and not two booleans"
// REQ-12: "The terminal condition is carried as ONE three-valued report per
// drain ... the three states are mutually exclusive and totally ordered
// here (overflow > held > whole)"
// ADVERSARIAL
func TestReq7_OverflowOutranksHeldOnTheSameDrain(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The child floods past `StdoutCap` AND leaves a `setsid(2)` escapee
	// holding both pipes, so both conditions qualify on the same drain.
	r := escReader(t, "10s", "overflow-and-hold", pidfile,
		strconv.Itoa(cmdbind.StdoutCap+4096))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("a stdout past the 1 MiB bound with a held pipe was accepted")
	}
	detail := escHeldDetail(t, err)
	text := err.Error() + "\n" + detail
	if !strings.Contains(text, strconv.Itoa(cmdbind.StdoutCap)) {
		t.Fatalf("the refusal is %q / detail %q; C1 `precedence:` fixes that "+
			"overflow OUTRANKS held on the same drain — the existing "+
			"\"stdout exceeded the %d byte bound\" refusal stands unchanged, "+
			"because the bytes were read and the pipe's state is then not "+
			"what the invocation turns on",
			err.Error(), detail, cmdbind.StdoutCap)
	}
}

// REQ-79 (LBD selection): "When the `werr` non-ExitError arm
// (`exec.ErrWaitDelay`, C1 `stdin:`) and a drain condition both qualify,
// the drain conditions are selected FIRST — overflow, then held — and the
// non-ExitError arm refuses only if neither fired."
// ADVERSARIAL
func TestReq79_ADrainConditionIsSelectedAheadOfTheNonExitErrorArm(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The WRITE path is the only one that carries a stdin payload:
	// `Reader.Read` always sends the empty object, which fits the pipe
	// buffer and can never outlive the child. `bigPlan` is 128 KiB, past
	// the buffer, so the parent's stdin write blocks on a read end the
	// escapee holds open. The direct child exits normally and immediately,
	// so nothing signals it and `Wait` ends in `exec.ErrWaitDelay` — the
	// non-ExitError arm. The same escapee holds both output pipes, so the
	// held condition qualifies on the same invocation and the two compete.
	w := &cmdbind.Writer{
		Accessor: escEntry(t, "2s", "stdin-holds-exits", pidfile),
		Name:     escName,
		Config:   cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := w.Apply(ctx, escArtifact(t), bigPlan())

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("a held pipe with a competing ErrWaitDelay was accepted")
	}
	text := err.Error() + "\n" + escHeldDetail(t, err)
	if !escNamesHeld(text, "stdout", "stderr") {
		t.Fatalf("the refusal is %q; it does not name the HELD pipes. LBD "+
			"selection fixes that when the `werr` non-ExitError arm and a "+
			"drain condition both qualify the drain conditions are selected "+
			"FIRST — a held pipe names the helper that held it, which "+
			"`ErrWaitDelay` does not, so the more diagnostic reason wins",
			text)
	}
	if errors.Is(err, exec.ErrWaitDelay) {
		t.Fatalf("the refusal is `exec.ErrWaitDelay`: the drain conditions " +
			"are selected FIRST and the non-ExitError arm refuses only if " +
			"neither fired")
	}
}

// --- C1 `whole:` — the unchanged whole-output path -----------------------

// REQ-49 (C1 `whole:`): "a drain that reaches EOF — before the timer, or on
// its final read after it — is whole, as today; the 1 MiB stdout cap and
// 4 KiB stderr tail are unchanged. A child (and group) that closes its
// pipes observes no change."
// REQ-60 (TD): "On EOF within the bound: unchanged — whole output, the
// 1 MiB cap and 4 KiB tail apply as today."
// REQ-86 (MC `fidelity`): "Whole read, child closes pipes" → "byte-for-byte
// unchanged from today"
// REQ-93 (MC `disposition`): "Child closes pipes, exits 0" → "success",
// "none", "envelope parsed", "loud (the value)"
// REQ-124 (S3): "Ordinary paths — child closes its pipes ... Expected:
// byte-for-byte unchanged from today"
// PRESERVED INVARIANT — this holds today and must keep holding across the
// `readBounded` rewrite C1 `precedence:` mandates.
// HAPPY PATH
func TestReq49_AChildThatClosesItsPipesObservesNoChange(t *testing.T) {
	t.Parallel()

	const envelope = `{"state.phase":"review"}`
	r := escReader(t, "10s", "whole", envelope, "diagnostic tail\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
	if err != nil {
		t.Fatalf("a child that closed its pipes and exited 0 refused: %v", err)
	}
	if len(unreadable) != 0 || len(values) != 1 || values[0].Value != "review" {
		t.Fatalf("values=%+v unreadable=%v; want state.phase=review "+
			"byte-for-byte", values, unreadable)
	}
}

// REQ-14 (C1 `precedence:`): "The rewrite must also preserve three
// behaviors the existing suite pins: the read stays bounded at
// `StdoutCap+1` so exactly-at-cap is a value and cap+1 is a failure ..."
// REQ-15: "A re-implementation that shifts the cap boundary is a contract
// violation, not a refactor"
// REQ-87 (MC `fidelity`): "stdout at the 1 MiB cap" → "exactly the cap is a
// value; one byte more is a detected overflow"
// REQ-104 (MC `oracle`, S3 control): "one byte past the cap must still
// refuse overflow"
// REQ-124 (S3): "Assert both boundary cases explicitly here rather than
// relying on the suite's incidental coverage; relaxing an existing
// expectation to make the rewritten loop pass is the failure this row
// names."
// REQ-136 (XC memory): "the existing 1 MiB stdout cap and 4 KiB stderr tail
// are unchanged and remain the only bound on retained output"
// PRESERVED INVARIANT — the cap boundary holds today; C1 `precedence:`
// rewrites the loop that enforces it, so both sides are pinned explicitly.
// BOUNDARY
func TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		size    int
		refuses bool
	}{
		{"one byte under the cap is a value", cmdbind.StdoutCap - 1, false},
		{"exactly the cap is a value", cmdbind.StdoutCap, false},
		{"one byte past the cap is a detected overflow", cmdbind.StdoutCap + 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := escReader(t, "20s", "emit-n", strconv.Itoa(tc.size))
			r.Accessor.Output = escRaw()

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			values, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
			switch {
			case tc.refuses && err == nil:
				t.Fatalf("a stdout of %d bytes — one past the %d byte cap — "+
					"was accepted. C1 `precedence:` states that a "+
					"re-implementation shifting the cap boundary is a "+
					"contract violation, not a refactor",
					tc.size, cmdbind.StdoutCap)
			case !tc.refuses && err != nil:
				t.Fatalf("a stdout of %d bytes (cap = %d) refused: %v. "+
					"Exactly the cap is a VALUE", tc.size, cmdbind.StdoutCap, err)
			case !tc.refuses:
				if len(values) != 1 || len(values[0].Value) != tc.size {
					got := 0
					if len(values) == 1 {
						got = len(values[0].Value)
					}
					t.Fatalf("recovered %d of %d bytes; the bound is on what "+
						"is KEPT and exactly the cap is kept whole",
						got, tc.size)
				}
			}
		})
	}
}

// REQ-88 (MC `fidelity`): "stderr tail" → "LAST 4 KiB preserved"
// REQ-124 (S3): "output at the ... 4 KiB stderr tail boundaries"
// PRESERVED INVARIANT — holds today; re-pinned across the rewrite.
// BOUNDARY
func TestReq88_TheStderrTailKeepsTheLastFourKiBAcrossTheRewrite(t *testing.T) {
	t.Parallel()

	// The fixture writes a positionally distinguishable byte stream, so a
	// tail taken from the WRONG end is caught rather than merely a short one.
	const over = cmdbind.StderrTailCap + 4096
	r := escReader(t, "20s", "stderr-n", strconv.Itoa(over))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	if err == nil {
		t.Fatal("a child with empty stdout and no exit map match was accepted")
	}
	detail := escHeldDetail(t, err)
	if len(detail) != cmdbind.StderrTailCap {
		t.Fatalf("the stderr tail is %d bytes; the cap is %d and the tail is "+
			"the LAST %d bytes, unchanged by this record",
			len(detail), cmdbind.StderrTailCap, cmdbind.StderrTailCap)
	}
	// The LAST bytes, not the first: byte i of the stream is 'a'+(i%26).
	wantLast := byte('a' + ((over - 1) % 26))
	if detail[len(detail)-1] != wantLast {
		t.Fatalf("the tail's last byte is %q; want %q — the tail keeps the "+
			"LAST 4 KiB, which is the diagnosis a caller wants",
			detail[len(detail)-1], wantLast)
	}
}

// REQ-50 (C1 `whole:`): "The completeness guarantee (A3) carries an
// antecedent, stated here because nothing structural enforces it: the two
// drains run CONCURRENTLY and are never stalled."
// REQ-51: "A child whose payload exceeds the ~64 KiB pipe buffer blocks in
// write(2) until drained, so a serialized or stalled drain converts a
// large-answer success into a deadline kill — measured: exactly one pipe
// buffer (65536 bytes) survives, as a correct prefix, where a prompt drain
// returns the whole 1 MiB."
// REQ-52: "this clause is what keeps a later refactor from silently voiding
// A3, since the prohibition is carried by the contract and S7, not by a
// build error"
// REQ-75 (LBD): "**The drains stay concurrent** — A3's no-loss guarantee
// holds only while they do ... Recorded here, and pinned by S7"
// REQ-89 (MC `fidelity`): "Drain past the ~64 KiB pipe buffer, prompt" →
// "whole 1 MiB recovered (S7)"
// REQ-91: "Drain past the pipe buffer, STALLED" → "exactly one pipe buffer
// (65536 bytes) as a **correct prefix**"
// REQ-129 (S7): "prompt drain returns the whole 1 MiB; a stalled drain
// leaves the child blocked in `write(2)`, killed by the group signal, with
// exactly one pipe buffer (65536 bytes) recovered as a correct prefix."
// REQ-128 (S7 seam): "the stall is ONE seam ... one unexported
// package-level hook set by the test and nil in production (not a build
// tag, not a sleep in production code). S7 is a standing regression trap"
// DOMAIN EDGE
func TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer(t *testing.T) {
	const payload = 1 << 20 // 1 MiB, past the ~64 KiB pipe buffer

	t.Run("prompt: the whole 1 MiB is recovered", func(t *testing.T) {
		t.Parallel()

		r := escReader(t, "20s", "emit-n", strconv.Itoa(payload))
		r.Accessor.Output = escRaw()

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		values, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
		if err != nil {
			t.Fatalf("the prompt drain refused a 1 MiB answer: %v — the two "+
				"drains run CONCURRENTLY and are never stalled (C1 `whole:`, "+
				"A3's antecedent)", err)
		}
		if len(values) != 1 || len(values[0].Value) != payload {
			got := 0
			if len(values) == 1 {
				got = len(values[0].Value)
			}
			t.Fatalf("the prompt drain recovered %d of %d bytes; A3's no-loss "+
				"guarantee holds only while the drains stay concurrent",
				got, payload)
		}
	})

	t.Run("stalled: exactly one pipe buffer as a correct prefix", func(t *testing.T) {
		// The stall seam is the ONE test-only hook S7, S9 and MVV row 6 share.
		// A drain stalled past the child's whole life leaves it blocked in
		// `write(2)` and unable to exit, so the group signal kills it and the
		// pipe holds exactly one buffer.
		const stall = 3 * time.Second
		cmdbind.SetDrainStartDelayForTest(t, stall)

		r := escReader(t, "2s", "emit-n", strconv.Itoa(payload))
		r.Accessor.Output = escRaw()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		start := time.Now()
		values, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
		if elapsed := time.Since(start); elapsed < 2*time.Second {
			t.Fatalf("the invocation returned in %v with the drain-start "+
				"stall set to %v and a 2 s deadline: the drain was never "+
				"stalled, so the child was never blocked in `write(2)` and "+
				"this row is not exercising the stall it exists to trap",
				elapsed, stall)
		}
		got := 0
		if len(values) == 1 {
			got = len(values[0].Value)
		}
		const pipeBuffer = 65536
		if err == nil && got == payload {
			t.Fatalf("a STALLED drain recovered the whole %d bytes; S7 fixes "+
				"the measured consequence at exactly one pipe buffer (%d) as "+
				"a correct prefix, and the row exists so a later refactor "+
				"that serializes or delays the drains fails HERE rather than "+
				"silently voiding A3", payload, pipeBuffer)
		}
		if got > 0 && !strings.HasPrefix(strings.Repeat("x", payload), values[0].Value) {
			t.Fatalf("the stalled drain's %d bytes are not a correct PREFIX "+
				"of the payload; the admitted loss is the remainder, never a "+
				"corrupted head", got)
		}
	})
}

// REQ-24 (C1 `precedence:`): "No read end outlives `spawn` (deadline to
// unblock, close to release the fd). The existing defers at pipe creation
// already satisfy this ... An explicit close before them would double-close,
// so none is added"
// REQ-66 (IC): "The block deliberately draws NO `graceOn` and NO
// `closeReads` call ... the existing defers already close both read ends,
// so an explicit close would double-close"
// REQ-20 (C1 `precedence:` (a)): "no separate `graceOn` flag is introduced"
// ADVERSARIAL — the fd-leak oracle. A build that added an explicit close
// would double-close and a build that leaked would grow the descriptor set;
// the observable is the process's own fd count across many invocations.
func TestReq24_NoReadEndOutlivesSpawnAcrossRepeatedHeldPipeRefusals(t *testing.T) {
	t.Parallel()

	before := openFDCount(t)

	for i := range 8 {
		pidfile := escPIDFile(t)
		r := escReader(t, "300ms", "deadline", pidfile)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
		cancel()
		if err == nil {
			t.Fatalf("iteration %d: a held-pipe read was accepted", i)
		}
		escKill(escReadPID(t, pidfile))
	}

	after := openFDCount(t)
	// A leak of two fds per invocation over eight invocations is sixteen;
	// the slack absorbs the runtime's own churn without absorbing a leak.
	if after-before >= 8 {
		t.Fatalf("open descriptors grew from %d to %d over 8 held-pipe "+
			"refusals: C1 `precedence:` states that no read end outlives "+
			"`spawn` — the existing defers release both fds at function "+
			"return, which is after the join, the `werr` switch and the "+
			"overflow check", before, after)
	}
}

// REQ-50 (C1 `whole:`): "the two drains run CONCURRENTLY and are never
// stalled."
// REQ-51: "A child whose payload exceeds the ~64 KiB pipe buffer blocks in
// write(2) until drained, so a serialized or stalled drain converts a
// large-answer success into a deadline kill"
// REQ-52: "this clause is what keeps a later refactor from silently voiding
// A3, since the prohibition is carried by the contract and S7, not by a
// build error" — the NEGATIVE requirement: the drains may not be serialized.
// REQ-75 (LBD): "**The drains stay concurrent**"
// REQ-129 (S7): a stalled drain leaves the child blocked in `write(2)`.
//
// S7 traps a drain stalled at ENTRY. This row traps the same loss arriving
// through the sibling halt, which is the shape S7's seam cannot reach: the
// halt is raised by a drain's OWN report, and `readBounded` raises it on any
// read error once the drain has overflowed — INCLUDING EOF. The drains start
// before `cmd.Wait()`, so a child that overflows stdout and then closes it
// gives the stdout drain EOF-after-overflow while the child is STILL LIVE.
// The halt then ends the stderr drain at the top of its loop, abandoning a
// pipe that live child is still writing: the child blocks in `write(2)`,
// `cmd.Wait()` never returns, and the invocation runs to the declared
// deadline and is misclassified as a timeout — which HIDES the stdout
// overflow refusal it actually owes. Reproduced against a harness of this
// shape: `cmd.Wait()` blocked past 5 s with the stderr drain returned held
// having read 32768 of 262144 bytes, against 21 ms and whole output with the
// halt absent.
//
// The fixture carries no escapee at all, so nothing here is the admitted
// residue: the only writer is the direct child, and the only reason it could
// fail to exit is a drain that stopped reading it.
// ADVERSARIAL — sibling-halt stall (roborev 6717).
func TestReq52_ASiblingHaltNeverStallsALiveChildsOtherDrain(t *testing.T) {
	t.Parallel()

	// Stdout one byte past the cap, so the stdout drain overflows and then
	// reaches EOF on the child's own close. Stderr well past the ~64 KiB
	// pipe buffer, so the child MUST be read to completion to exit.
	const overflow = cmdbind.StdoutCap + 1
	const stderrPayload = 256 << 10

	// The declared deadline is long relative to the committed bound, so a
	// return here is the drain's own doing and an elapsed time near the
	// deadline is the stall this row traps.
	const deadline = 10 * time.Second
	r := escReader(t, "10s", "overflow-then-stderr",
		strconv.Itoa(overflow), strconv.Itoa(stderrPayload))
	r.Accessor.Output = escRaw()

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	start := time.Now()
	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("a stdout of %d bytes — one past the %d byte cap — was "+
			"accepted after %v", overflow, cmdbind.StdoutCap, elapsed)
	}
	// The refusal owed is the stdout OVERFLOW, and it is owed PROMPTLY: the
	// child exits as soon as its stderr is drained, so the invocation costs
	// the drain and not the deadline.
	if budget := escBoundFor(0); elapsed > budget {
		t.Fatalf("the overflowing child refused only after %v, past the "+
			"committed total of %v.\n\n"+
			"The child writes %d bytes to stdout, CLOSES it, then writes "+
			"%d bytes to stderr — past the ~64 KiB pipe buffer, so it "+
			"blocks in `write(2)` until the stderr drain reads it. The "+
			"stdout drain sees EOF AFTER overflowing while the child is "+
			"still live and raises the sibling halt before `cmd.Wait()` has "+
			"returned; the stderr drain then observes that halt at the top "+
			"of its loop and returns, abandoning a pipe a LIVE writer is "+
			"still filling.\n\n"+
			"C1 `whole:` states A3's antecedent — \"the two drains run "+
			"CONCURRENTLY and are never stalled\" — and REQ-52 makes it a "+
			"negative requirement: the drains may not be serialized. A halt "+
			"observed while the child is unreaped serializes them, and the "+
			"consequence is the one S7 measures: the child never exits, so "+
			"the invocation runs to its deadline and reports a TIMEOUT "+
			"instead of the stdout-overflow refusal it owes.",
			elapsed, budget, overflow, stderrPayload)
	}
	// The overflow refusal carries its reason on the WRAPPED error (the
	// `Detail` slot holds the stderr tail), so the oracle reads the
	// refusal's own text.
	if text := escExecError(t, err).Error(); !strings.Contains(text, "exceeded") {
		t.Fatalf("the refusal reads %q; the owed refusal is the stdout "+
			"overflow, which C1 `precedence:` ranks above every drain "+
			"condition because the bytes were READ", text)
	}
}
