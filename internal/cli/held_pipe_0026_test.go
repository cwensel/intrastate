package cli

// RDR 0026 — the CLI rendering of a held-pipe write refusal (S1's rendered
// half), plus the negative witness deviations D1 requires.
//
// The oracle is the envelope the CLI emits, not any internal call: S1 fixes
// the rendered half on "the envelope's `detail` field carrying
// `detailMayHaveApplied` (`flow_exec.go:423`'s existing constant, per
// docs/cli-output-contract.md's error envelope) — the same observable the
// read-back leg already renders, so the two legs are asserted the same way
// and neither needs a new golden" (REQ-121).
//
// Deviations D2 (decided at JDR 0003 §D4 (b)) refines the assertion: the
// direct write leg renders a DISTINCT exit-3 code keyed on `Applied()`,
// with the `detail` text unchanged, and not-applied refusals keep
// `flow-accessor-failed` (REQ-139/REQ-140). D2 states the spelling is
// non-normative, so no test pins the literal beyond what the envelope
// itself requires.
//
// The `applied` sense is unexported by design ("only the write path — which
// knows whether the command already ran — may set it"), so every refusal
// here is minted by the REAL `accessor.Executor` over a fixture binding.
// Nothing mocks the unit under test: `accessorFailureOf` is the unit, and
// the executor is what feeds it.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- fixtures: bindings that mint each refusal shape ---------------------

const (
	hpRole   = "state"
	hpKey    = "status"
	hpWriter = "cmd.write"
	hpReader = "cmd.read"
	hpGate   = "cmd.gate"
)

// hpBinding fails every invocation with the supplied error. It stands in
// for a command binding whose drain reported a held pipe: the error is the
// only channel the seam gives it, which is exactly what `cmdbind` uses.
type hpBinding struct {
	cap accessor.Capability
	err error

	applies int
}

func (b *hpBinding) Capability() accessor.Capability { return b.cap }

func (b *hpBinding) Read(context.Context, accessor.Artifact, []string) (
	[]accessor.KeyValue, []string, error,
) {
	return nil, nil, b.err
}

func (b *hpBinding) Gate(context.Context, accessor.Artifact) (
	accessor.Verdict, string, error,
) {
	return "", "", b.err
}

func (b *hpBinding) Apply(context.Context, accessor.Artifact, []resolve.Tag) error {
	b.applies++
	return b.err
}

func (b *hpBinding) Invocations() int { return b.applies }

func hpDef(name string, cap accessor.Capability, binding accessor.Binding) accessor.Definition {
	acc := table.Accessor{
		Role: hpRole, Path: "flow.state",
		Keys: []string{hpKey}, Timeout: "1s",
	}
	if cap == accessor.CapWrite {
		acc.ReadBack = true
	}
	return accessor.Definition{
		Identity: accessor.Identity{Flow: "hpflow", Name: name, Capability: cap},
		Accessor: acc,
		Binding:  binding,
	}
}

func hpExecutor(defs ...accessor.Definition) *accessor.Executor {
	return accessor.NewExecutor(
		accessor.Registry{Flow: "hpflow", Definitions: defs, OwnedTags: []string{hpKey}},
		accessor.Artifacts{hpRole: {Role: hpRole, Path: "flow.state"}},
	)
}

// hpHeldErr is the refusal a command binding returns for a held pipe: the
// typed `*accessor.ExecError` carrier whose `Err` wraps the typed held-pipe
// error the executor matches (C1 `refusal:`).
func hpHeldErr() error {
	return &accessor.ExecError{
		Detail: "the child's stdout, stderr stayed held past the 500 ms " +
			"bound (exited 0); close or redirect the helper's inherited stdio",
		Err: &accessor.HeldPipeError{
			Held: []string{"stdout", "stderr"}, ExitStatus: "exited 0",
		},
	}
}

// hpWriteRefusal drives the REAL executor's write path with the supplied
// binding error and returns the refusal it minted.
func hpWriteRefusal(t *testing.T, err error) accessor.Refusal {
	t.Helper()

	b := &hpBinding{cap: accessor.CapWrite, err: err}
	got := hpExecutor(hpDef(hpWriter, accessor.CapWrite, b)).
		Write(context.Background(), hpWriter,
			resolve.Plan{RuleID: "r", Writes: []resolve.Tag{{Key: hpKey, Value: "final"}}})
	if got.Refusal == nil {
		t.Fatal("the write did not refuse")
	}
	return *got.Refusal
}

// hpReadBackRefusal drives a write whose command SUCCEEDS and whose
// read-back is the leg that refuses — the arm C1 `refusal:` says is already
// correct and must not be re-keyed.
func hpReadBackRefusal(t *testing.T) accessor.Refusal {
	t.Helper()

	writer := &hpBinding{cap: accessor.CapWrite}
	reader := &hpBinding{cap: accessor.CapRead, err: hpHeldErr()}
	got := hpExecutor(
		hpDef(hpWriter, accessor.CapWrite, writer),
		hpDef(hpReader, accessor.CapRead, reader),
	).Write(context.Background(), hpWriter,
		resolve.Plan{RuleID: "r", Writes: []resolve.Tag{{Key: hpKey, Value: "final"}}})
	if got.Refusal == nil {
		t.Fatal("the read-back leg did not refuse")
	}
	return *got.Refusal
}

// --- S1's rendered half ---------------------------------------------------

// REQ-120 (S1): "on BOTH write legs `Applied()` is true AND the CLI renders
// the \"may have been applied\" text"
// REQ-121 (S1): "The rendered half asserts on the envelope's `detail` field
// carrying `detailMayHaveApplied` ... so the two legs are asserted the same
// way and neither needs a new golden."
// REQ-139 (D2, decided): "a distinct exit-3 CLI code for the applied write
// refusal ..., keyed on `Applied()` in `flow_exec.go::accessorFailureOf`'s
// write arm; `detail` text unchanged; not-applied refusals keep
// `flow-accessor-failed`."
// REQ-140 (D2 check): "S1's rendered half asserts the new code AND the
// prose on the direct write leg"
// REQ-114 (IP Phase 4): "key `flow_exec.go::accessorFailureOf`'s
// `ClassExecutionFailure` arm (`:410-412`) on `Applied()`."
// REQ-32 (C1 `refusal:`): that arm "carries neither a phase check nor a
// `Detail`, so a held-pipe write slips past the \"may have been applied\"
// rendering"
// REQ-96 (MC `disposition`): "Held pipe on the DIRECT write" →
// "`execution_failure` + `Applied()` true ... LOUD once the re-key lands"
// HAPPY PATH
func TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode(t *testing.T) {
	refusal := hpWriteRefusal(t, hpHeldErr())
	if !refusal.Applied() {
		t.Fatal("the executor did not set the applied sense on a held-pipe " +
			"DIRECT write; C1 `refusal:` fixes it as a post-run refusal")
	}

	ce := accessorFailureOf(refusal, phaseWrite)
	if ce == nil {
		t.Fatal("an applied execution_failure rendered no CLI error")
	}
	if ce.Detail != detailMayHaveApplied {
		t.Errorf("detail = %q; want %q. S1's rendered half asserts on the "+
			"envelope's `detail` field carrying `detailMayHaveApplied` — the "+
			"SAME observable the read-back leg already renders, so the two "+
			"legs are asserted the same way and neither needs a new golden",
			ce.Detail, detailMayHaveApplied)
	}
	if ce.Code == codeAccessorFailed {
		t.Errorf("code = %q; the applied write refusal takes a DISTINCT code "+
			"(JDR 0003 §D4 (b), deviations D2), keyed on `Applied()` in the "+
			"write arm. Not-applied refusals keep %q, so a shared code makes "+
			"the two indistinguishable to an agent branching on it",
			ce.Code, codeAccessorFailed)
	}
	if got := clierr.ExitCodeFor(ce); got != 3 {
		t.Errorf("exit code = %d; want 3 — D2 fixes \"a distinct EXIT-3 CLI "+
			"code for the applied write refusal\"", got)
	}
}

// REQ-138 (D1, negative witness): "A unit test on `accessorFailureOf`:
// (write phase, `ClassExecutionFailure`, `Applied()` false) → no
// `detailMayHaveApplied`"
// REQ-140 (D2 check): "D1's negative witness asserts `flow-accessor-failed`
// with no prose on the `Applied()`-false write refusal."
// REQ-117 (IP Phase 4): "S1 gains a negative row: a binding whose command
// cannot start refuses `execution_failure` with `Applied()` FALSE and no
// `detailMayHaveApplied` in the envelope."
// REQ-116: "Only the held-pipe error may set the applied sense there."
// ADVERSARIAL — the neighbouring `ClassTimeout` arm keys on PHASE; copying
// that shape would label every pre-mutation write refusal applied.
func TestReq138_AnUnappliedWritePhaseExecutionFailureRendersNoAppliedProse(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{
			"a command that cannot start",
			&accessor.ExecError{Err: errors.New(
				`exec: "nope": executable file not found in $PATH`)},
		},
		{
			"a plain non-zero exit",
			&accessor.ExecError{Detail: "helper: bad flag\n", Err: errors.New("exit status 2")},
		},
		{
			"an untyped binding error",
			errors.New("the write binding could not open the artifact"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			refusal := hpWriteRefusal(t, tc.err)
			if refusal.Applied() {
				t.Fatalf("the executor set the applied sense on %s; only the "+
					"held-pipe error may do so on the SHARED `err != nil` "+
					"arm", tc.name)
			}

			ce := accessorFailureOf(refusal, phaseWrite)
			if ce == nil {
				t.Fatal("no CLI error was rendered")
			}
			if strings.Contains(ce.Detail, detailMayHaveApplied) {
				t.Errorf("detail = %q on a WRITE-phase execution_failure whose "+
					"`Applied()` is FALSE. The neighbouring `ClassTimeout` arm "+
					"keys on PHASE; copying that shape would label every "+
					"pre-mutation write refusal applied, and the agent skips a "+
					"safe retry. The key is `Applied()`, not the phase",
					ce.Detail)
			}
			if ce.Code != codeAccessorFailed {
				t.Errorf("code = %q; want %q — D2 fixes that not-applied "+
					"refusals KEEP the existing code and only the applied "+
					"write refusal takes the new one", ce.Code, codeAccessorFailed)
			}
		})
	}
}

// REQ-33 (C1 `refusal:`): the read-back arm "already renders with
// `detailMayHaveApplied`. So one site changes, not two"
// REQ-122 (S1): "The READ-BACK leg refuses `read_back_incomplete` and is a
// REGRESSION assertion, not a proof of the gain: that arm already sets
// `detailMayHaveApplied` ..., so the row fails only if the re-key broke a
// path that was correct"
// REQ-97 (MC `disposition`): "Held pipe on the write's READ-BACK" →
// "`read_back_incomplete` + `applied` already true ... LOUD today"
// PRESERVED INVARIANT — this arm is correct today.
// HAPPY PATH
func TestReq122_TheReadBackArmStillRendersTheAppliedProseAndKeepsItsOwnCode(t *testing.T) {
	refusal := hpReadBackRefusal(t)
	if refusal.Class != accessor.ClassReadBackIncomplete {
		t.Fatalf("Class = %q; want %q — the read-back leg's class is NOT "+
			"re-keyed by this record", refusal.Class,
			accessor.ClassReadBackIncomplete)
	}
	if !refusal.Applied() {
		t.Fatal("Applied() = false on the read-back leg; the executor already " +
			"sets it there today")
	}

	ce := accessorFailureOf(refusal, phaseWrite)
	if ce == nil {
		t.Fatal("no CLI error was rendered")
	}
	if ce.Detail != detailMayHaveApplied {
		t.Errorf("detail = %q; want %q — the read-back arm already sets it "+
			"today and this record re-keys ONE site, not two",
			ce.Detail, detailMayHaveApplied)
	}
	if ce.Code != codeReadBackIncomplete {
		t.Errorf("code = %q; want %q — the read-back leg's own code is not "+
			"re-keyed", ce.Code, codeReadBackIncomplete)
	}
}

// REQ-34 (C1 `refusal:`): "On read and gate invocations `Applied()` is
// false as today."
// REQ-138 (D1): "(read phase, refusal wrapping the typed held-pipe error) →
// `Applied()` false"
// ADVERSARIAL — a read or gate refusal wrapping the SAME typed held-pipe
// error must never render the write path's remediation.
func TestReq34_AReadOrGatePhaseHeldPipeRefusalNeverRendersTheAppliedProse(t *testing.T) {
	read := hpExecutor(hpDef(hpReader, accessor.CapRead,
		&hpBinding{cap: accessor.CapRead, err: hpHeldErr()})).
		Read(context.Background(), hpReader)
	if read.Refusal == nil {
		t.Fatal("the held-pipe read did not refuse")
	}
	if read.Refusal.Applied() {
		t.Error("Applied() = true on a READ refusal wrapping the typed " +
			"held-pipe error; C1 `refusal:` fixes it false")
	}

	for _, at := range []phase{phaseRead, phaseGate} {
		ce := accessorFailureOf(*read.Refusal, at)
		if ce == nil {
			t.Fatalf("phase %v rendered no CLI error", at)
		}
		if strings.Contains(ce.Detail, detailMayHaveApplied) {
			t.Errorf("phase %v: detail = %q; a read or gate refusal is never "+
				"applied, so the CLI must not offer the write path's "+
				"remediation", at, ce.Detail)
		}
		if ce.Code != codeAccessorFailed {
			t.Errorf("phase %v: code = %q; want %q", at, ce.Code, codeAccessorFailed)
		}
	}
}

// REQ-83 (MC `authority`): "`withStderrTail` composes into one slot", "one
// slot, applied-sense first"
// REQ-26 (C1 `refusal:`) + deviations D3: the held-pipe reason LEADS and
// the stderr tail FOLLOWS inside the binding's own Detail; at the CLI the
// applied-sense text comes first and the tail is appended after it — one
// slot, no second carrier.
// REQ-136 (XC memory): the retained-output bounds are the only ones, so no
// second Detail carrier is minted.
// BOUNDARY
func TestReq83_TheAppliedSenseAndTheHeldPipeReasonShareTheOneDetailSlotInOrder(t *testing.T) {
	refusal := hpWriteRefusal(t, hpHeldErr())
	ce := withStderrTail(accessorFailureOf(refusal, phaseWrite), refusal.Detail)
	if ce == nil {
		t.Fatal("no CLI error was rendered")
	}
	if !strings.Contains(ce.Detail, detailMayHaveApplied) {
		t.Fatalf("detail = %q; the applied-sense text was lost to the tail",
			ce.Detail)
	}
	if !strings.Contains(ce.Detail, "close or redirect") {
		t.Fatalf("detail = %q; the held-pipe reason did not survive into the "+
			"one Detail slot — the remediation is what tells the author what "+
			"to fix", ce.Detail)
	}
	if strings.Index(ce.Detail, detailMayHaveApplied) >
		strings.Index(ce.Detail, "close or redirect") {
		t.Errorf("detail = %q; the applied-sense text must come FIRST — "+
			"losing \"the write may have applied\" to a diagnostic tail "+
			"would drop the more consequential fact", ce.Detail)
	}
}

// REQ-139 / REQ-140 (D2): the new code joins the exit-3 accessor family and
// is DISTINCT from every existing member, so an agent branching on it
// cannot confuse an applied write refusal with any other refusal.
// ADVERSARIAL
func TestReq139_TheAppliedWriteCodeIsDistinctFromEveryExistingAccessorCode(t *testing.T) {
	applied := accessorFailureOf(hpWriteRefusal(t, hpHeldErr()), phaseWrite)
	if applied == nil {
		t.Fatal("no CLI error was rendered")
	}

	existing := []string{
		codeAccessorFailed, codeAccessorTimeout,
		codeReadBackTimeout, codeReadBackIncomplete, codeReadBackMismatch,
	}
	for _, code := range existing {
		if applied.Code == code {
			t.Fatalf("the applied write refusal renders %q, which is already "+
				"an existing accessor code. D2 fixes it as a DISTINCT exit-3 "+
				"code keyed on `Applied()`, so JDR 0001 §D10's table gains "+
				"the code by citation repair", code)
		}
	}
	if applied.Code == "" {
		t.Fatal("the applied write refusal renders an EMPTY code")
	}
	if !strings.HasPrefix(applied.Code, "flow-") {
		t.Errorf("the applied write code %q does not join the `flow-` family; "+
			"the spelling is non-normative but the namespace is not",
			applied.Code)
	}
}

// REQ-137 (XC `Incremental adoption`, also CON): "No opt-out flag ships and
// none is proposed — a per-binding \"accept partial output\" knob is ALT1
// by another name, rejected for ALT1's reason; the escape hatch is
// redirection at the binding's command, which costs no contract."
// ADVERSARIAL — asserted structurally over the whole command tree, so a
// flag added anywhere fails here.
func TestReq137_NoOptOutFlagForTheHeldPipeRefusalShips(t *testing.T) {
	root := NewRootCmd()
	forbidden := []string{
		"accept-partial-output", "allow-partial-output",
		"allow-held-pipes", "drain-bound", "drain-grace",
	}
	for _, name := range forbidden {
		if root.PersistentFlags().Lookup(name) != nil {
			t.Errorf("the root command carries a persistent `--%s` flag; XC "+
				"`Incremental adoption` states that no opt-out flag ships "+
				"and none is proposed", name)
		}
		for _, sub := range root.Commands() {
			if sub.Flags().Lookup(name) != nil {
				t.Errorf("`%s` carries a `--%s` flag; a per-binding \"accept "+
					"partial output\" knob is ALT1 by another name",
					sub.Name(), name)
			}
		}
	}
}
