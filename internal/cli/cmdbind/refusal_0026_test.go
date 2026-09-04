package cmdbind_test

// RDR 0026 — C1 `refusal:`, `class:`, `stdin:`, `bound:` and the failure
// modes F5/F6 the record accepts by name.
//
// The refusal's SHAPE is the contract; its Go spelling is not. C1
// `precedence:` gives the report type's "identifier and the Go spelling"
// to the implementer (REQ-13) and the Illustrative Code names
// `heldPipes`/`errHeldPipe` as "names for the shape, not the contract"
// (REQ-65), so every assertion here is on the user-visible refusal — the
// pipe names, the exit status, the bound, the remediation, the tail — and
// on the ONE type the contract does fix by location: the held-pipe error
// declared in `internal/accessor`.

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- C1 `refusal:` — the carrier and its content -------------------------

// REQ-25 (C1 `refusal:`): "refusal: a held pipe refuses through
// `*accessor.ExecError`, this record's instance of JDR 0003 §D1"
// REQ-26: "`Detail` carries the held-pipe reason — the pipe(s) held
// (\"stdout\"/\"stderr\"/both), the bound, the direct child's exit status
// (exited N / killed by signal) and the remediation (\"close or redirect
// the helper's inherited stdio\") — composed AHEAD of the stderr tail
// collected up to the bound."
// REQ-27: "That exit status is composed from fields `spawn` already
// populates from `werr` — `inv.exited` and `inv.signaled` beside
// `inv.exitCode` — so no new field is added to carry it"
// REQ-28: "`Err` wraps a typed held-pipe error carrying the pipes held and
// the exit status, which the executor matches with `errors.As` at the
// refusal site."
// REQ-61 (TD): "returns `wrap(detail, errHeldPipe)` — `Detail` leading with
// the held pipe(s) and the child's exit status ahead of the stderr tail"
// REQ-83 (MC `authority`): "`refusal.Detail`, ordered per 0025:C4"
// REQ-8: "returns the held pipe names in the fixed order stdout, stderr, so
// `len(held) != 0` is the consulted predicate"
// REQ-MVV.5: "`Err` naming \"stdout, stderr\" and \"exited 0\""
// Ordering per deviations D3: the held-pipe reason LEADS, the stderr tail
// FOLLOWS.
// HAPPY PATH
func TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	// The child writes a whole envelope and EXITS 0 while the escapee holds
	// both pipes — MVV row 5's success shape, the leg on which the
	// held-pipe reason is reachable (F6: it is not on the `timeout` leg).
	r := escReader(t, "20s", "success", pidfile, `{"state.phase":"review"}`)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("a child that exited 0 with both pipes held by a `setsid(2)` " +
			"escapee was accepted; C1 `precedence:` fixes the held pipe as " +
			"the refusal")
	}
	ee := escExecError(t, err)
	detail := ee.Detail

	// The pipe(s) held, in the FIXED order.
	if !escNamesHeld(detail, "stdout", "stderr") {
		t.Errorf("Detail = %q; it does not name the held pipes as "+
			"\"stdout, stderr\". C1 `precedence:` fixes the ORDER (stdout, "+
			"stderr) and MVV row 5 quotes that exact pair", detail)
	}
	// The direct child's exit status — composed from fields `spawn` already
	// populates, which is why no new invocation field carries it.
	if !strings.Contains(detail, "exited 0") {
		t.Errorf("Detail = %q; it does not carry the direct child's exit "+
			"status (\"exited 0\"), which C1 `refusal:` and MVV row 5 both "+
			"name", detail)
	}
	// The bound.
	if !strings.Contains(detail, strconv.Itoa(cmdbind.WaitDelay)) {
		t.Errorf("Detail = %q; it does not name the bound (%d ms). C1 "+
			"`refusal:` lists the bound among what the held-pipe reason "+
			"carries", detail, cmdbind.WaitDelay)
	}
	// The remediation.
	if !strings.Contains(detail, "close or redirect") {
		t.Errorf("Detail = %q; it does not carry the remediation \"close or "+
			"redirect the helper's inherited stdio\", which is what tells an "+
			"author what to fix", detail)
	}
	// `Err` wraps a TYPED held-pipe error the executor can match, not a
	// flattened string.
	if ee.Err == nil {
		t.Error("ExecError.Err is nil; C1 `refusal:` fixes that `Err` WRAPS " +
			"a typed held-pipe error carrying the pipes held and the exit " +
			"status, which the executor matches with `errors.As`")
	}
}

// REQ-29 (C1 `refusal:`): "The type is declared in `internal/accessor`, NOT
// in `cmdbind`: `cmdbind` imports `accessor` ... so a concrete `cmdbind`
// type named at an `errors.As` site inside `executor.go` is an import cycle
// that does not compile."
// REQ-30: "Being matched from `executor.go` and populated from `cmdbind`,
// both outside its own package, the type is EXPORTED and its held-pipes and
// exit-status fields with it — an unexported type would be unnameable at
// the `errors.As` site"
// REQ-84 (MC `authority`): "`internal/accessor` — `cmdbind` imports
// `accessor`, never the reverse"
// REQ-111 (IP Phase 1): "Declare the held-pipe error type in
// `internal/accessor` beside `ExecError`"
// DOMAIN EDGE — the assertion is that the type is NAMEABLE and MATCHABLE
// from outside `cmdbind` at all, which is exactly what the import direction
// requires and what an unexported or `cmdbind`-local type would fail.
func TestReq29_TheHeldPipeErrorIsAnExportedAccessorTypeMatchableFromOutsideCmdbind(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	r := escReader(t, "20s", "success", pidfile, `{"state.phase":"review"}`)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("the held-pipe case was accepted")
	}

	// The match is `errors.As` against a type declared in `internal/accessor`
	// and named HERE, from a package that is neither `accessor` nor
	// `cmdbind` — the same position `executor.go` occupies.
	var hp *accessor.HeldPipeError
	if !errors.As(err, &hp) {
		t.Fatalf("the refusal does not match `*accessor.HeldPipeError` via "+
			"errors.As: %#v. C1 `refusal:` fixes the type's declaration in "+
			"`internal/accessor` (the import direction forbids `cmdbind`) "+
			"and its EXPORT (an unexported type would be unnameable at the "+
			"`errors.As` site in executor.go)", err)
	}
	// The held-pipes and exit-status fields are exported with it.
	if len(hp.Held) == 0 {
		t.Error("HeldPipeError carries no held pipes; C1 `refusal:` fixes " +
			"that the typed error carries \"the pipes held and the exit " +
			"status\"")
	}
	if got := strings.Join(hp.Held, ", "); got != "stdout, stderr" {
		t.Errorf("the held pipes are %q; want %q — the FIXED order C1 "+
			"`precedence:` states", got, "stdout, stderr")
	}
}

// REQ-35 (C1 `refusal:`): "The returned invocation carries no stdout;
// stdout read on a held drain is never parsed on ANY path — read envelope,
// gate verdict, write, write read-back (JDR 0003 §D1)."
// REQ-57 (AP): "Output read on a bound-expired drain is never parsed on any
// path."
// REQ-62 (TD): "Nothing read from a held stdout is handed to `Reader.parse`,
// `Gate.Gate`'s verdict mapping, or the writer's read-back."
// REQ-92 (MC `fidelity`): "Held pipe" → "stdout is never parsed on ANY
// path"; exemption "the whole stream — refused, not truncated-then-read"
// REQ-101 (MC `oracle`, MVV row 5 control): "a build that parses held
// stdout returns values and fails the not-parsed assertion"
// REQ-136 (XC memory): "a partial buffer is never retained as a value"
// REQ-MVV.5: "the envelope NOT parsed (no values returned)"
// ADVERSARIAL — the child writes a PERFECTLY VALID envelope, so a build
// that parses held stdout returns values and this row catches it.
func TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath(t *testing.T) {
	t.Parallel()

	const envelope = `{"state.phase":"review"}`

	t.Run("read: the envelope is not parsed", func(t *testing.T) {
		t.Parallel()

		pidfile := escPIDFile(t)
		r := escReader(t, "20s", "success", pidfile, envelope)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		values, unreadable, err := r.Read(ctx, escArtifact(t), []string{escKey})
		escKill(escReadPID(t, pidfile))

		if err == nil || len(values) != 0 {
			t.Fatalf("a held stdout carrying a WELL-FORMED envelope returned "+
				"values=%+v unreadable=%v err=%v. C1 `refusal:` fixes that "+
				"the returned invocation carries no stdout and a held stdout "+
				"is never parsed on ANY path — refused, not "+
				"truncated-then-read", values, unreadable, err)
		}
	})

	t.Run("gate: no verdict is derived", func(t *testing.T) {
		t.Parallel()

		pidfile := escPIDFile(t)
		g := cmdbind.Gate{
			Accessor: escEntry(t, "20s", "success", pidfile,
				`{"verdict":"allow","reason":"ok"}`),
			Name:   escName,
			Config: cmdbind.Config{AllowCommands: true},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		verdict, reason, err := g.Gate(ctx, escArtifact(t))
		escKill(escReadPID(t, pidfile))

		if err == nil || verdict != "" {
			t.Fatalf("a held stdout carrying a WELL-FORMED allow verdict "+
				"returned verdict=%q reason=%q err=%v; no gate verdict may be "+
				"derived from held output", verdict, reason, err)
		}
	})
}

// REQ-36 (C1 `refusal:`): "\"Carries no stdout\" is a CALLER obligation,
// not a structural guarantee, and is stated as one so a later caller cannot
// inherit it by accident ... The obligation is that every caller checks
// `err` before touching `inv` — true of today's three callers by inspection
// (A1), and NOT enforced by any type or signature."
// REQ-37: "Implementation carries this as a comment at the `inv`
// construction site naming the two error returns; making it structural (an
// error-only return, or zeroing `inv` on the refusal arms) is a
// `readBounded`-adjacent change this record does not take"
// REQ-68 (EIA): "Its signature is unchanged by this record: it returns the
// `invocation` by VALUE (not a pointer), which is what makes C1
// `refusal:`'s \"the field is readable on every error path\" hold without a
// nil check"
// DOMAIN EDGE — the obligation is on the three callers, so the assertable
// consequence is that EVERY caller returns the refusal rather than a value,
// on every capability. That is what a structural change would also give,
// and what a caller that read `inv` ahead of `err` would break.
func TestReq36_EveryCallerChecksErrBeforeTouchingTheInvocation(t *testing.T) {
	t.Parallel()

	const envelope = `{"state.phase":"review"}`

	// The write path: a held pipe must refuse rather than confirm.
	pidfile := escPIDFile(t)
	w := &cmdbind.Writer{
		Accessor: escEntry(t, "20s", "success", pidfile, envelope),
		Name:     escName,
		Config:   cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := w.Apply(ctx, escArtifact(t), nil)
	escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatal("a write whose pipes were held by a `setsid(2)` escapee " +
			"reported success. `spawn` builds `inv` BEFORE the `werr` switch " +
			"and both error arms return that populated value, so the field " +
			"is readable on every error path — the obligation C1 `refusal:` " +
			"states is that every caller checks `err` FIRST")
	}
	if got := w.Invocations(); got != 1 {
		t.Errorf("Invocations() = %d; want 1 — one Apply is one spawn and one "+
			"increment, unchanged by this record", got)
	}
}

// --- C1 `class:` — classification stays the executor's --------------------

// REQ-38 (C1 `class:`): "class: the class is the executor's, as today: ctx
// `DeadlineExceeded` ⇒ `timeout` (deadline-first rule; per 0025:C4 a
// `timeout` refusal carries no Err/Detail, so the held-pipe reason survives
// only on `execution_failure`), else the ExecError ⇒ `execution_failure`."
// REQ-39: "The binding never classifies `timeout` itself"
// REQ-64 (TD): "`spawn` therefore never has to know whether the bound
// expired because of a deadline or because a success-path child leaked a
// writer."
// REQ-80 (MC `authority`): "the EXECUTOR; the binding never classifies
// `timeout` (C1 `class:`)"
// ADVERSARIAL — the binding is driven under an ALREADY-elapsed deadline; it
// must still return a plain binding error and never a class of its own.
func TestReq39_TheBindingNeverClassifiesTimeoutItself(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	r := escReader(t, "300ms", "deadline", pidfile)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatal("the held-pipe case under an elapsed deadline was accepted")
	}
	// The binding's channel is the typed ExecError. It carries no class, and
	// the executor's own deadline-first rule is what mints `timeout` — so a
	// binding error naming the class is the failure this row catches.
	ee := escExecError(t, err)
	text := strings.ToLower(ee.Detail + " " + ee.Error())
	if strings.Contains(text, "timeout") {
		t.Errorf("the BINDING's refusal names a timeout class: detail=%q "+
			"err=%q. C1 `class:` fixes the class as the EXECUTOR's, read "+
			"from its own bounded context's DeadlineExceeded first; the "+
			"binding never classifies `timeout` itself", ee.Detail, ee.Error())
	}
}

// --- C1 `stdin:` — the unchanged stdin leg -------------------------------

// REQ-40 (C1 `stdin:`): "stdin: unchanged and named: the child's stdin is a
// Cmd-owned pipe (`cmd.Stdin` is a `bytes.Reader`), bounded by Cmd's own
// WaitDelay; a non-reading child or escapee holding its read end ends in
// `exec.ErrWaitDelay` from Wait, which `spawn`'s non-ExitError arm refuses
// as `execution_failure` AFTER reapGroup and the join — no arm of `spawn`
// returns between Wait and the join"
// REQ-98 (MC `disposition`): "Non-reading child, stdin past the buffer" →
// "`execution_failure`", "`exec.ErrWaitDelay` via `spawn`'s non-ExitError
// arm", "none", "loud (S4, A7)"
// REQ-125 (S4): "Non-reading child with a stdin payload past the pipe
// buffer (A7). Expected: `exec.ErrWaitDelay` from `Wait` reaches `spawn`'s
// non-`ExitError` arm and refuses `execution_failure` after `reapGroup` and
// the join."
// PRESERVED INVARIANT — the stdin leg is "unchanged and named"; the row
// exists so the rewritten join does not disturb it.
// DOMAIN EDGE
func TestReq40_ANonReadingChildWithAStdinPayloadPastTheBufferRefusesAfterTheJoin(t *testing.T) {
	t.Parallel()

	// A write carries a real stdin payload; the child never reads it and
	// exits immediately, leaving an escapee holding the stdin READ END, so
	// Cmd's WaitDelay is what ends the stdin write. The escapee holds
	// stdin ALONE: the output drains reach EOF, no drain condition
	// qualifies, and `exec.ErrWaitDelay` is the sole reason to refuse.
	pidfile := escPIDFile(t)
	w := &cmdbind.Writer{
		Accessor: escEntry(t, "1s", "stdin-holds-only", pidfile),
		Name:     escName,
		Config:   cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	err := w.Apply(ctx, escArtifact(t), bigPlan())
	elapsed := time.Since(start)

	defer escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatal("a non-reading child with a stdin payload past the pipe " +
			"buffer reported success")
	}
	if !errors.Is(err, exec.ErrWaitDelay) {
		t.Fatalf("the refusal is %q; REQ-98/REQ-125 fix the refusal for a "+
			"non-reading child with a stdin payload past the buffer as "+
			"`exec.ErrWaitDelay` reaching `spawn`'s non-ExitError arm — an "+
			"error of any other kind means the case never staged, so the "+
			"row proves nothing about that arm", err)
	}
	if want := escBoundFor(1 * time.Second); elapsed > want {
		t.Fatalf("the stdin case returned after %v, past the C1 `bound:` "+
			"total plus tolerance (%v): no arm of `spawn` returns between "+
			"Wait and the join, and the whole journey stays inside the one "+
			"bound", elapsed, want)
	}
}

// --- C1 `bound:` — the one bound, no second knob -------------------------

// REQ-41 (C1 `bound:`): "WaitDelay (500 ms) is the ONE bound — the same
// constant, by design, for every \"the CLI waits past the child's exit\"
// case ... no second knob, not model-declarable."
// REQ-42: "DrainGrace (50 ms) is NOT a second bound: it answers a different
// question ... is never model-declarable, and cannot move the refuse/accept
// boundary — it is a mechanism constant inside the one bound, named here
// because every shipped constant of this family is named in its contract"
// REQ-70 (EIA): "C1 `bound:` — no second bound"
// REQ-71 (LBD naming): "the bound is `cmdbind.go::WaitDelay`, reused;
// rejected: a new `DrainBound` constant ... and a model-declarable drain
// bound"
// REQ-72 (LBD grace): "`DrainGrace` (50 ms) ships beside `WaitDelay` ... it
// is an offset on the single existing timer's expiry."
// REQ-73: "50 ms is chosen against the measured scheduling tail"
// REQ-85 (MC `authority`): "`WaitDelay`; no second knob (D-naming)";
// "`DrainGrace` is an OFFSET, not a rival bound"
// REQ-135 (XC): the seam "needs no third build tag — it is portable Go"
// BOUNDARY
func TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound(t *testing.T) {
	t.Parallel()

	// The two constants ship side by side in this package, with the values
	// the contract names. `DrainGrace` is an OFFSET inside the one bound, so
	// it must be strictly smaller than it — a grace at or past the bound
	// would BE a second bound.
	if cmdbind.WaitDelay != 500 {
		t.Errorf("WaitDelay = %d ms; C1 `bound:` names 500 ms as the ONE "+
			"bound, and the constant is in the contract because C1 publishes "+
			"timeout + 2·WaitDelay", cmdbind.WaitDelay)
	}
	if cmdbind.DrainGrace != 50 {
		t.Errorf("DrainGrace = %d ms; C1 `bound:` and D-the-drain-grace both "+
			"name 50 ms, chosen against the measured scheduling tail (A5 saw "+
			"drain-goroutine latency reach ~28 ms under 4x-core saturation) "+
			"so the grace clears it with margin while costing 10%% of one "+
			"bound", cmdbind.DrainGrace)
	}
	if cmdbind.DrainGrace >= cmdbind.WaitDelay {
		t.Errorf("DrainGrace (%d ms) is not smaller than WaitDelay (%d ms); "+
			"the grace is a mechanism constant INSIDE the one bound, not a "+
			"rival bound", cmdbind.DrainGrace, cmdbind.WaitDelay)
	}

	// Not model-declarable: no entry field carries either constant. The
	// author declares `timeout`; the drain bound is the CLI's own mechanism
	// cost, not policy (D-naming).
	if declaresDrainBound() {
		t.Error("a model entry declares a drain bound or grace; C1 `bound:` " +
			"fixes both as NOT model-declarable — a second knob is ALT1's " +
			"parallel-signal defect by another name")
	}
}

// REQ-128 (S7 seam): "one unexported package-level hook set by the test and
// nil in production (not a build tag, not a sleep in production code)"
// ADVERSARIAL — the seam's production value is the assertion; a shipped
// delay would slow every invocation.
func TestReq128_TheStallSeamIsNilInProduction(t *testing.T) {
	if !cmdbind.DrainStartDelayIsNilInProduction() {
		t.Error("the drain-start stall seam carries a non-zero delay in " +
			"production; S7 fixes it as one unexported package-level hook " +
			"set by the test and NIL in production — not a build tag, and " +
			"not a sleep in production code")
	}
}

// REQ-43 (C1 `bound:`): "Cancel sends one signal and returns, so Cmd's
// timer starts at the deadline."
// REQ-44: the committed total, "asserted with a scheduling tolerance of
// +100 ms"
// REQ-47: "The two legs are NOT separately bounded and the contract does
// not claim they are"
// REQ-48: "a child that RESISTS termination spends leg (a)'s allowance and
// is the case where the total tightens — S8 tests it"
// REQ-130 (S8): "a child that RESISTS termination (ignores SIGTERM,
// survives to its `WaitDelay`) while a grandchild holds the pipes ...
// Expected: the invocation still returns within the C1 `bound:` total plus
// its stated scheduling tolerance"
// BOUNDARY — the case where the two legs genuinely compete.
func TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	r := escReader(t, "1s", "resists", pidfile)

	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	elapsed := time.Since(start)

	escKill(escReadPID(t, pidfile))

	if err == nil {
		t.Fatal("a SIGTERM-resisting child with a held pipe was accepted")
	}
	if want := escBoundFor(timeout); elapsed > want {
		t.Fatalf("a child that RESISTS termination returned after %v, past "+
			"the C1 `bound:` committed total of %v plus its stated %v "+
			"scheduling tolerance. This is the case where leg (a) spends its "+
			"allowance instead of donating it to leg (b), so the TOTAL is "+
			"asserted where the two legs genuinely compete",
			elapsed, timeout+2*escWaitDelay, escTolerance)
	}
}

// --- F5 — the non-pollable pipe ------------------------------------------

// REQ-63 (TD): "Pollability is checked once at pipe creation (a far-future
// deadline; `os.ErrNoDeadline` marks the pipe non-pollable) so the fallback
// is decided before any child exists (F5)."
// REQ-69 (EIA): "F5's pollability result is recorded here as one more
// field, which is how it reaches the `Detail` line S5 asserts on"
// REQ-99 (MC `disposition`): "Non-pollable pipe (`os.ErrNoDeadline`)" →
// "unbounded join — the old hang", "recorded on the invocation at the
// check", "a `Detail` line naming the host condition", "loud at the check,
// by design — as a `Detail` line, not a log line; `cmdbind` has no logger"
// REQ-107 (F5): "It runs in `spawn` at pipe creation: a probe
// `SetReadDeadline` on each read end, immediately cleared, whose error is
// the signal. `cmdbind` has no logger and gains none ... the condition is
// carried on the invocation and surfaces as a `Detail` line on any refusal
// that invocation produces, which is what a test asserts on."
// REQ-126 (S5): "simulated at the check by injecting a read end that
// refuses a deadline ... Expected: the condition is recorded at the check,
// not discovered at the join, and is observable as the F5 `Detail` line on
// that invocation's refusal"
// DOMAIN EDGE
func TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine(t *testing.T) {
	// S5 simulates the condition AT THE CHECK: the injected read end is one
	// `os.NewFile` over a dup'd pipe fd — adopted without the original's
	// poller registration — which is the shape F5 confirms produces
	// `os.ErrNoDeadline`.
	cmdbind.SetNonPollableReadEndsForTest(t)

	// The child fails loudly so the invocation produces a refusal for the
	// F5 line to ride on.
	r := escReader(t, "5s", "stderr-n", "64")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
	if err == nil {
		t.Fatal("the fixture did not refuse, so there is no refusal for the " +
			"F5 Detail line to ride on")
	}
	detail := escHeldDetail(t, err)
	if !strings.Contains(detail, "deadline") {
		t.Fatalf("Detail = %q; it carries no line naming the host condition. "+
			"F5 fixes that `cmdbind` has no logger and gains none, so "+
			"\"logged\" is not a log line: the condition is carried on the "+
			"invocation and surfaces as a `Detail` line on any refusal that "+
			"invocation produces, which is what a test asserts on", detail)
	}
}

// --- F6 — held with the deadline also elapsed ----------------------------

// REQ-95 (MC `disposition`): "Held pipe, deadline ALSO elapsed" →
// "`timeout`", "none — 0025:C4 gives `timeout` no Err/Detail", "no stdout",
// "loud in class, **quiet in reason** (F6, S6)"
// REQ-108 (F6): "a held pipe that also outran the deadline reports
// `timeout` with no `Err`/`Detail` (0025:C4's rule for the class), so the
// held-pipe reason is visible only on `execution_failure`."
// REQ-74 (LBD grace): "the tail reaches the user on the `execution_failure`
// leg (MVV row 5, S10) and NOT when the declared timeout also elapsed"
// REQ-77 (LBD selection): "When `timeout` and `execution_failure` both
// qualify: `timeout` (the executor's existing deadline-first rule, not a
// new tie-break)."
// REQ-127 (S6): "Held stdout with the deadline also elapsed (F6). Expected:
// classified `timeout` with empty `Detail` per 0025:C4"
// REQ-MVV.3: "the executor classifies `timeout`, and the refusal's `Detail`
// is empty as 0025:C4 states for `timeout`"
// DOMAIN EDGE — this lives at the executor, which is where the class is
// minted; see held_pipe_0026_test.go for the executor-side row.

// --- S10 — the single-pipe form ------------------------------------------

// REQ-132 (S10): "stdout reaches EOF cleanly while STDERR alone is held by
// the escapee ... Expected: refused; `Detail` names `stderr` alone — the
// single-pipe form is the bare pipe name, and the both-pipes form is the
// two names comma-separated in the fixed order stdout, stderr (MVV row 5's
// \"stdout, stderr\") — and the stderr tail collected up to the bound
// survives into `Detail`"
// REQ-8: "the single-pipe form is the bare name (S10, MVV row 5)"
// REQ-103 (MC `oracle`, S10 control): "stdout-held-alone: the same refusal
// with the other pipe named, and no tail"
// Per deviations D3 the assertion is the tail's SURVIVAL and the bare
// single-pipe name, NOT the tail's position ahead of the reason: the
// contract's ordering (reason leads, tail follows) is what REQ-26 pins.
// DOMAIN EDGE
func TestReq132_StderrHeldAloneNamesTheBarePipeAndKeepsTheTail(t *testing.T) {
	t.Parallel()

	const tailText = "helper: cannot reach the daemon\n"

	t.Run("stderr held alone: the bare name and the surviving tail", func(t *testing.T) {
		t.Parallel()

		pidfile := escPIDFile(t)
		r := escReader(t, "20s", "stderr-held", pidfile, tailText,
			`{"state.phase":"review"}`)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
		escKill(escReadPID(t, pidfile))

		if err == nil {
			t.Fatal("stdout reached EOF cleanly while STDERR alone was held " +
				"by the escapee, and the read was accepted")
		}
		detail := escHeldDetail(t, err)
		if !strings.Contains(detail, "stderr") {
			t.Fatalf("Detail = %q; it does not name `stderr`", detail)
		}
		// The single-pipe form is the BARE name: not "stdout, stderr", and
		// not a bracketed or quoted list.
		if strings.Contains(detail, "stdout, stderr") {
			t.Errorf("Detail = %q names BOTH pipes; only stderr is held, and "+
				"the single-pipe form is the bare pipe name", detail)
		}
		// The tail collected up to the bound SURVIVES into Detail — this is
		// the grace's whole load-bearing justification.
		if !strings.Contains(detail, strings.TrimSpace(tailText)) {
			t.Errorf("Detail = %q does not carry the stderr tail %q collected "+
				"up to the bound. D-the-drain-grace justifies `DrainGrace` on "+
				"exactly this: \"what the grace preserves is the stderr tail\" "+
				"— without it the mechanism has no test", detail, tailText)
		}
	})

	t.Run("negative control: stdout held alone names the other pipe and no tail", func(t *testing.T) {
		t.Parallel()

		pidfile := escPIDFile(t)
		r := escReader(t, "20s", "stdout-held", pidfile)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})
		escKill(escReadPID(t, pidfile))

		if err == nil {
			t.Fatal("stdout held alone was accepted")
		}
		detail := escHeldDetail(t, err)
		if !strings.Contains(detail, "stdout") {
			t.Fatalf("Detail = %q; it does not name `stdout`", detail)
		}
		if strings.Contains(detail, "stderr") {
			t.Errorf("Detail = %q names `stderr`, which is NOT held here. "+
				"S10's named control is the same refusal with the OTHER pipe "+
				"named, and no tail — a fold that names both regardless is "+
				"exactly what it catches", detail)
		}
	})
}

// --- C1 `residue:` — the admitted leak -----------------------------------

// REQ-53 (C1 `residue:`): "0025:F4 restated — a process outside the child's
// process group may outlive the refusal; it keeps write ends whose reader
// is gone, so its next write fails with EPIPE (SIGPIPE under the default
// disposition), and whatever it wrote is never read."
// REQ-55: "Diagnosis: the refusal's Detail names the held pipe(s); the
// leaked process is visible to `ps` until its next write"
// REQ-100 (MC `disposition`): "Escapee outside the group" → "not reaped",
// "EPIPE on its next write", "none", "SILENT — admitted residue"
// REQ-109 (F4): "if A3 is wrong, bytes the direct child wrote can be lost
// at the bound; the refusal still fires, so the loss is bounded to a
// refused invocation, never a parsed one."
// DOMAIN EDGE
func TestReq53_TheEscapeeOutlivesTheRefusalAndIsDiagnosedByTheNamedPipes(t *testing.T) {
	t.Parallel()

	pidfile := escPIDFile(t)
	r := escReader(t, "20s", "success", pidfile, `{"state.phase":"review"}`)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, _, err := r.Read(ctx, escArtifact(t), []string{escKey})

	pid := escReadPID(t, pidfile)
	defer escKill(pid)

	if err == nil {
		t.Fatal("the escapee case was accepted")
	}
	// Not reaped — the admitted residue, SILENT by design.
	if !escAlive(pid) {
		t.Errorf("the escapee pid %d did not outlive the refusal; C1 "+
			"`residue:` restates 0025:F4 as \"a process outside the child's "+
			"process group MAY outlive the refusal\" and nothing in the CLI "+
			"reaps a `setsid(2)` escapee", pid)
	}
	// And the diagnosis is the refusal's Detail naming the held pipe(s) —
	// the only visible half of an otherwise silent residue.
	if detail := escHeldDetail(t, err); !escNamesHeld(detail, "stdout", "stderr") {
		t.Errorf("Detail = %q; C1 `residue:` fixes the DIAGNOSIS as \"the "+
			"refusal's Detail names the held pipe(s)\"", detail)
	}
}

// --- helpers --------------------------------------------------------------

// bigPlan builds a stdin payload past the ~64 KiB pipe buffer, so a
// non-reading child blocks the write unless Cmd's own WaitDelay bounds it.
func bigPlan() []resolve.Tag {
	// One 128 KiB value comfortably exceeds the pipe buffer once encoded.
	return []resolve.Tag{{Key: escKey, Value: strings.Repeat("x", 128<<10)}}
}

// declaresDrainBound reports whether a model entry carries a field naming a
// drain bound or grace. C1 `bound:` fixes both as NOT model-declarable, so
// the assertion is the ABSENCE of such a field on the entry type.
func declaresDrainBound() bool {
	rt := reflect.TypeOf(table.Accessor{})
	for i := range rt.NumField() {
		name := strings.ToLower(rt.Field(i).Name)
		if strings.Contains(name, "drain") || strings.Contains(name, "grace") ||
			strings.Contains(name, "waitdelay") {
			return true
		}
	}
	return false
}
