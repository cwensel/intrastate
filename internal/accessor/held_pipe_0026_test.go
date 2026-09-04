package accessor_test

// RDR 0026 — the held-pipe error type, the applied sense on the DIRECT
// write leg, and the classification that stays the executor's.
//
// Nothing here mocks the executor. The bindings are fixtures returning the
// typed held-pipe error through the interfaces' existing bare `error` slot,
// which is exactly what `cmdbind` does; the oracle is the Refusal the real
// executor mints. That is also the only way to assert the type's DECLARED
// LOCATION obligation from a third package: `cmdbind` imports `accessor`
// and `accessor` imports no `cli` package, so a concrete `cmdbind` type
// named at an `errors.As` site inside `executor.go` would not compile
// (C1 `refusal:`).

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
)

// heldErr builds the refusal a command binding returns when a drain
// reported a held pipe: the typed `*accessor.ExecError` carrier, whose
// `Detail` leads with the held-pipe reason and whose `Err` wraps the typed
// held-pipe error the executor matches (C1 `refusal:`; deviations D3 fixes
// the ordering — the reason LEADS, the tail FOLLOWS).
func heldErr(held []string, status, tail string) error {
	reason := "the child's " + strings.Join(held, ", ") + " stayed held past " +
		"the 500 ms bound (" + status + "); close or redirect the helper's " +
		"inherited stdio"
	detail := reason
	if tail != "" {
		detail = reason + "\n" + tail
	}
	return &accessor.ExecError{
		Detail: detail,
		Err:    &accessor.HeldPipeError{Held: held, ExitStatus: status},
	}
}

// --- C1 `refusal:` — the applied sense on the DIRECT write leg -----------

// REQ-31 (C1 `refusal:`): "Applied sense: a held pipe on the write path
// after the direct child was reaped is a post-run refusal — `Applied()` is
// true, carried per JDR 0003 §D1 through the executor's write arm."
// REQ-32: "The DIRECT write invocation refuses `execution_failure`, and it
// is the leg that needs the gain: `executor.go::Write`'s `err != nil` arm
// sets no applied sense today"
// REQ-96 (MC `disposition`): "Held pipe on the DIRECT write" →
// "`execution_failure` + `Applied()` true", "held-pipe reason", "no
// read-back", "LOUD once the re-key lands (A8)"
// REQ-114 (IP Phase 4): "match the held-pipe error in
// `executor.go::Write`'s `err != nil` arm ... and set the applied sense
// there"
// REQ-81 (MC `authority`): "`Refusal.applied`; `Applied()` is the only
// reader key (A8)"
// HAPPY PATH
func TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense(t *testing.T) {
	b := &execFailBinding{
		cap: accessor.CapWrite,
		err: heldErr([]string{"stdout", "stderr"}, "exited 0", ""),
	}
	e := execExecutor(execDef(execWriter, accessor.CapWrite, b))

	got := e.Write(execCtx(t), execWriter,
		resolve.Plan{RuleID: "r", Writes: []resolve.Tag{{Key: execKey, Value: "final"}}})

	if got.Refusal == nil {
		t.Fatal("a held pipe on the direct write did not refuse")
	}
	if got.Refusal.Class != accessor.ClassExecutionFailure {
		t.Fatalf("Class = %q; want %q — the DIRECT write invocation refuses "+
			"`execution_failure`, and C1 `refusal:` states the two write legs "+
			"already differ in class and this clause does not flatten them",
			got.Refusal.Class, accessor.ClassExecutionFailure)
	}
	if !got.Refusal.Applied() {
		t.Fatal("Applied() = false on a held-pipe DIRECT write. C1 " +
			"`refusal:` fixes it as a post-run refusal — the direct child was " +
			"already reaped, so the mutation may have been applied — and the " +
			"gain JDR 0003 §D1 rule 3 lands here is exactly that. An agent's " +
			"documented branch on `execution_failure` is retry, so an " +
			"unsensed write applies the mutation twice")
	}
	// No read-back runs: the write's own refusal is terminal.
	if got.Refusal.Class == accessor.ClassReadBackIncomplete {
		t.Error("the refusal is `read_back_incomplete`; a held pipe on the " +
			"DIRECT write refuses before any read-back (MC `disposition`)")
	}
}

// REQ-116 (IP Phase 4 negative control, required): "The `errors.As` match
// lands on the SHARED `err != nil` arm at `:368-371`, which every other
// pre-mutation failure shape also reaches — `resolveArgv0` failure,
// `os.Pipe` creation failure, a plain exec failure from a mistyped binding.
// Only the held-pipe error may set the applied sense there."
// REQ-117 (IP Phase 4): "S1 gains a negative row: a binding whose command
// cannot start refuses `execution_failure` with `Applied()` FALSE"
// REQ-138 (D1, negative witness): "A unit test on `accessorFailureOf`:
// (write phase, `ClassExecutionFailure`, `Applied()` false) → no
// `detailMayHaveApplied`; (read phase, refusal wrapping the typed held-pipe
// error) → `Applied()` false. The `errors.As` for the held-pipe error lives
// in `Executor.Write`'s `err != nil` arm, never in the shared `refusalOf`"
// ADVERSARIAL — a match broadened to `*accessor.ExecError` generally flips
// a write that STRUCTURALLY COULD NOT HAVE RUN to "may have been applied",
// and the agent skips a safe retry: the exact inverse of the defect Phase 4
// fixes, and invisible to the positive row above.
func TestReq116_OnlyTheHeldPipeErrorSetsTheAppliedSenseOnTheSharedWriteArm(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{
			// A spawn failure: no child ever existed.
			"a command that cannot start",
			&accessor.ExecError{Detail: "", Err: errors.New("exec: \"nope\": executable file not found in $PATH")},
		},
		{
			// A plain non-zero exit with a stderr tail: the command ran and
			// FAILED, but nothing says the mutation landed.
			"a plain non-zero exit",
			&accessor.ExecError{Detail: "helper: bad flag\n", Err: errors.New("exit status 2")},
		},
		{
			// An untyped binding error — the path-backed shape.
			"an untyped binding error",
			errors.New("the write binding could not open the artifact"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &execFailBinding{cap: accessor.CapWrite, err: tc.err}
			e := execExecutor(execDef(execWriter, accessor.CapWrite, b))

			got := e.Write(execCtx(t), execWriter,
				resolve.Plan{RuleID: "r", Writes: []resolve.Tag{{Key: execKey, Value: "final"}}})

			if got.Refusal == nil {
				t.Fatal("the write did not refuse")
			}
			if got.Refusal.Class != accessor.ClassExecutionFailure {
				t.Fatalf("Class = %q; want %q", got.Refusal.Class,
					accessor.ClassExecutionFailure)
			}
			if got.Refusal.Applied() {
				t.Fatalf("Applied() = TRUE on %s. The `errors.As` match lands "+
					"on the SHARED `err != nil` arm that every other "+
					"pre-mutation failure shape also reaches, and ONLY the "+
					"held-pipe error may set the applied sense there. A match "+
					"broadened to `*accessor.ExecError` generally flips a "+
					"write that structurally could not have run to \"may have "+
					"been applied\", and the agent skips a safe retry — the "+
					"exact inverse of the defect this phase fixes", tc.name)
			}
		})
	}
}

// REQ-34 (C1 `refusal:`): "On read and gate invocations `Applied()` is
// false as today."
// REQ-138 (D1): "(read phase, refusal wrapping the typed held-pipe error) →
// `Applied()` false"
// ADVERSARIAL — the held-pipe error reaching a READ must not carry the
// write path's applied sense with it.
func TestReq34_AHeldPipeOnAReadOrGateLeavesAppliedFalse(t *testing.T) {
	held := heldErr([]string{"stdout"}, "exited 0", "")

	t.Run("read", func(t *testing.T) {
		b := &execFailBinding{cap: accessor.CapRead, err: held}
		got := execExecutor(execDef(execReader, accessor.CapRead, b)).
			Read(execCtx(t), execReader)

		if got.Refusal == nil {
			t.Fatal("the held-pipe read did not refuse")
		}
		if got.Refusal.Applied() {
			t.Error("Applied() = true on a READ refusal wrapping the typed " +
				"held-pipe error. C1 `refusal:` states it plainly: \"On read " +
				"and gate invocations `Applied()` is false as today\" — the " +
				"applied sense belongs to the write path, which knows whether " +
				"the command already ran")
		}
	})

	t.Run("gate", func(t *testing.T) {
		b := &execFailBinding{cap: accessor.CapGate, err: held}
		got := execExecutor(execDef(execGate, accessor.CapGate, b)).
			Gate(execCtx(t), execGate)

		if got.Refusal == nil {
			t.Fatal("the held-pipe gate did not refuse")
		}
		if got.Refusal.Applied() {
			t.Error("Applied() = true on a GATE refusal wrapping the typed " +
				"held-pipe error; C1 `refusal:` fixes it false")
		}
	})
}

// REQ-33 (C1 `refusal:`): "The READ-BACK leg after the write ran is already
// whole and gains nothing: a held pipe there returns the write's own
// refusal re-classed `read_back_incomplete` with `applied` already set ...
// So one site changes, not two"
// REQ-97 (MC `disposition`): "Held pipe on the write's READ-BACK" →
// "`read_back_incomplete` + `applied` already true", "read-back's own
// error", "no confirmation", "LOUD today"
// REQ-122 (S1): "The READ-BACK leg refuses `read_back_incomplete` and is a
// REGRESSION assertion, not a proof of the gain"
// PRESERVED INVARIANT — this arm is already correct today; the row fails
// only if the re-key broke a path that was correct.
// HAPPY PATH
func TestReq33_TheReadBackLegStaysReadBackIncompleteWithAppliedAlreadySet(t *testing.T) {
	// The write itself succeeds; the READ-BACK's reader is the one whose
	// pipes are held.
	writer := &execFailBinding{cap: accessor.CapWrite}
	reader := &execFailBinding{
		cap: accessor.CapRead,
		err: heldErr([]string{"stdout"}, "exited 0", ""),
	}
	e := execExecutor(
		execDef(execWriter, accessor.CapWrite, writer),
		execDef(execReader, accessor.CapRead, reader),
	)

	got := e.Write(execCtx(t), execWriter,
		resolve.Plan{RuleID: "r", Writes: []resolve.Tag{{Key: execKey, Value: "final"}}})

	if got.Refusal == nil {
		t.Fatal("a held pipe on the write's read-back did not refuse")
	}
	if got.Refusal.Class != accessor.ClassReadBackIncomplete {
		t.Fatalf("Class = %q; want %q — the read-back leg's class is NOT "+
			"re-keyed by this record. \"So one site changes, not two — the "+
			"read-back leg is named here because a reader who assumes the "+
			"class is uniform across both legs would re-key an arm that is "+
			"correct\"", got.Refusal.Class, accessor.ClassReadBackIncomplete)
	}
	if !got.Refusal.Applied() {
		t.Error("Applied() = false on the read-back leg; the executor already " +
			"sets it there today and this row fails only if the re-key broke " +
			"a path that was correct")
	}
}

// --- C1 `class:` — deadline first ----------------------------------------

// REQ-38 (C1 `class:`): "the class is the executor's, as today: ctx
// `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a
// `timeout` refusal carries no Err/Detail, so the held-pipe reason survives
// only on `execution_failure`), else the ExecError ⇒ `execution_failure`."
// REQ-77 (LBD selection): "When `timeout` and `execution_failure` both
// qualify: `timeout` (the executor's existing deadline-first rule, not a
// new tie-break)."
// REQ-95 (MC `disposition`): "Held pipe, deadline ALSO elapsed" →
// "`timeout`", "none — 0025:C4 gives `timeout` no Err/Detail", "no stdout",
// "loud in class, quiet in reason (F6, S6)"
// REQ-108 (F6): "a held pipe that also outran the deadline reports
// `timeout` with no `Err`/`Detail`"
// REQ-127 (S6): "Held stdout with the deadline also elapsed (F6). Expected:
// classified `timeout` with empty `Detail` per 0025:C4 — asserted so the
// class flip is a tested consequence rather than a surprise."
// REQ-74 (LBD grace): "the tail reaches the user on the
// `execution_failure` leg (MVV row 5, S10) and NOT when the declared
// timeout also elapsed"
// DOMAIN EDGE
func TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail(t *testing.T) {
	// The binding sleeps past the declared timeout AND returns the held-pipe
	// refusal, so `timeout` and `execution_failure` both qualify.
	b := &slowHeldBinding{
		cap:   accessor.CapRead,
		sleep: 2 * execTimeout,
		err:   heldErr([]string{"stdout"}, "killed by signal", "helper tail\n"),
	}
	got := execExecutor(execDef(execReader, accessor.CapRead, b)).
		Read(execCtx(t), execReader)

	if got.Refusal == nil {
		t.Fatal("the held-pipe-plus-deadline case did not refuse")
	}
	if got.Refusal.Class != accessor.ClassTimeout {
		t.Fatalf("Class = %q; want %q — C1 `class:` reads the bounded "+
			"context's DeadlineExceeded FIRST, and LBD selection states this "+
			"is the executor's existing deadline-first rule, not a new "+
			"tie-break", got.Refusal.Class, accessor.ClassTimeout)
	}
	if got.Refusal.Detail != "" {
		t.Fatalf("Detail = %q on a `timeout` refusal; 0025:C4 gives `timeout` "+
			"no Err/Detail, so the held-pipe reason survives only on "+
			"`execution_failure`. F6 records this as an ACCEPTED trade-off — "+
			"amending 0025:C4's `detail:` rule is NOT this record's change",
			got.Refusal.Detail)
	}
}

// --- the typed held-pipe error's own shape -------------------------------

// REQ-28 (C1 `refusal:`): "`Err` wraps a typed held-pipe error carrying the
// pipes held and the exit status, which the executor matches with
// `errors.As` at the refusal site."
// REQ-30: "the type is EXPORTED and its held-pipes and exit-status fields
// with it"
// REQ-84 (MC `authority`): "one declaration, one match site"
// REQ-8: "returns the held pipe names in the fixed order stdout, stderr"
// BOUNDARY
func TestReq28_TheHeldPipeErrorWrapsRatherThanFlattensAndSurvivesToTheRefusalSite(t *testing.T) {
	hp := &accessor.HeldPipeError{
		Held:       []string{"stdout", "stderr"},
		ExitStatus: "exited 0",
	}
	carrier := &accessor.ExecError{Detail: "reason\ntail\n", Err: hp}

	var got *accessor.HeldPipeError
	if !errors.As(carrier, &got) {
		t.Fatal("errors.As through the *ExecError carrier does not reach the " +
			"typed held-pipe error; C1 `refusal:` fixes that `Err` WRAPS it " +
			"rather than flattening, so the executor can match it at the " +
			"refusal site")
	}
	if strings.Join(got.Held, ", ") != "stdout, stderr" {
		t.Errorf("Held = %v; want the fixed order stdout, stderr", got.Held)
	}
	if got.ExitStatus != "exited 0" {
		t.Errorf("ExitStatus = %q; want the direct child's status", got.ExitStatus)
	}
	// A sibling ExecError that is NOT a held pipe must not match — that is
	// what keeps the Phase-4 negative control assertable at all.
	other := &accessor.ExecError{Detail: "boom", Err: errors.New("exit status 2")}
	var none *accessor.HeldPipeError
	if errors.As(other, &none) {
		t.Error("a plain non-zero-exit ExecError matches *HeldPipeError; the " +
			"two must stay distinguishable or the shared write arm cannot " +
			"tell a post-run refusal from a pre-mutation one")
	}
}

// --- fixtures -------------------------------------------------------------

// slowHeldBinding sleeps past the declared timeout and then returns the
// held-pipe refusal, so the executor's deadline-first rule and the
// binding's own error both qualify. It mocks nothing under test: the
// executor is real and this stands in for a command binding whose drain
// reported held after the deadline elapsed.
type slowHeldBinding struct {
	cap   accessor.Capability
	sleep time.Duration
	err   error

	applies int
}

func (b *slowHeldBinding) Capability() accessor.Capability { return b.cap }

func (b *slowHeldBinding) wait(ctx context.Context) {
	select {
	case <-time.After(b.sleep):
	case <-ctx.Done():
	}
}

func (b *slowHeldBinding) Read(ctx context.Context, _ accessor.Artifact, _ []string) (
	[]accessor.KeyValue, []string, error,
) {
	b.wait(ctx)
	return nil, nil, b.err
}

func (b *slowHeldBinding) Gate(ctx context.Context, _ accessor.Artifact) (
	accessor.Verdict, string, error,
) {
	b.wait(ctx)
	return "", "", b.err
}

func (b *slowHeldBinding) Apply(ctx context.Context, _ accessor.Artifact, _ []resolve.Tag) error {
	b.applies++
	b.wait(ctx)
	return b.err
}

func (b *slowHeldBinding) Invocations() int { return b.applies }
