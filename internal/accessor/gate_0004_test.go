package accessor_test

// RDR 0004 — gate accessor invocation semantics (REQ-35..REQ-39).
//
// SCOPE, per req-list Q2: this RDR ships gate INVOCATION semantics only —
// allow / deny / indeterminate, bounded timeout, and "reported, never
// applied". JDR 0001 §D9 lands the gate SITE (when gates run, only the
// selected row's, deny-overrides aggregation, `set-state` never gates) in
// RDR 0005. Nothing here schedules a gate, and nothing here aggregates
// several gates: the executor's own AP says "The resolver stays
// stateless" and CA A5's If-wrong is "The accessor layer becomes an
// orchestrator".
//
// Q1 is settled the layered way (REQ-37): at THIS boundary a deny is a
// TYPED RESULT carrying a reason; it is a refusal at the CLI. JDR 0001
// §D9: "Deny is a refusal at the CLI and a typed result at the accessor."

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
)

// REQ-35: "A gate accessor MUST return allow, deny, or indeterminate."
// HAPPY PATH
func TestReq35_GateVerdictVocabularyIsClosedAtThree(t *testing.T) {
	got := accessor.Verdicts()
	want := []accessor.Verdict{
		accessor.VerdictAllow, accessor.VerdictDeny, accessor.VerdictIndeterminate,
	}

	if len(got) != len(want) {
		t.Fatalf("Verdicts() = %v (%d members); want exactly the three %v",
			got, len(got), want)
	}
	for _, w := range want {
		if !slices.Contains(got, w) {
			t.Errorf("Verdicts() = %v; missing %q", got, w)
		}
	}
}

// REQ-35: "A gate accessor MUST return allow, deny, or indeterminate."
// HAPPY PATH
func TestReq35_AllowIsATypedGateResult(t *testing.T) {
	e := gateExec(t, &gateBinding{verdict: accessor.VerdictAllow})

	got := e.Gate(ctxOf(t), gateName)

	if got.Refused() {
		t.Fatalf("allow refused with class %q; allow is a typed gate result",
			got.Refusal.Class)
	}
	if got.Verdict != accessor.VerdictAllow {
		t.Errorf("verdict = %q; want %q", got.Verdict, accessor.VerdictAllow)
	}
}

// REQ-37: Deny is a **typed gate result carrying a reason** at this
// accessor boundary, and a refusal at the CLI (JDR 0001 §D9; DISP: "Gate
// returns deny | typed gate result | — (deny, with reason)").
// REQ-38: this RDR MUST NOT apply a deny, only report it.
// BOUNDARY
//
// The discriminating assertions: the reason SURVIVES to the caller, and
// the deny does not become an accessor refusal class here. An
// implementation that mints a `gate_denied` refusal class at this
// boundary fails both legs — and FM's "gate denied" sentence names the
// caller-visible disposition, not this boundary's return shape.
func TestReq37_DenyIsATypedResultCarryingAReasonAtTheAccessorBoundary(t *testing.T) {
	const reason = "the flow is frozen for release"
	e := gateExec(t, &gateBinding{verdict: accessor.VerdictDeny, reason: reason})

	got := e.Gate(ctxOf(t), gateName)

	if got.Refused() {
		t.Fatalf("deny took the refusal branch with class %q; at THIS boundary "+
			"a deny is a typed gate result (JDR 0001 §D9 — \"Deny is a refusal "+
			"at the CLI and a typed result at the accessor\")", got.Refusal.Class)
	}
	if got.Verdict != accessor.VerdictDeny {
		t.Errorf("verdict = %q; want %q", got.Verdict, accessor.VerdictDeny)
	}
	if got.Reason != reason {
		t.Errorf("reason = %q; want %q — a deny carries its reason", got.Reason, reason)
	}
	for _, c := range accessor.RefusalClasses() {
		if c == "gate_denied" {
			t.Error("the accessor refusal-class set mints `gate_denied`; deny is " +
				"reported as a typed result here and mapped to a refusal by the " +
				"CLI (RDR 0005), never applied by this layer")
		}
	}
}

// REQ-36: "Indeterminate MUST be a refusal-class result, not a false
// allow and not a false deny." — the class is `gate_indeterminate`.
// ADVERSARIAL
//
// The two negative halves are asserted explicitly: an implementation that
// folds indeterminate into allow (permissive) or into deny (restrictive)
// fails on the verdict assertion, not merely on the class name.
func TestReq36_IndeterminateIsARefusalNeverAFalseAllowOrDeny(t *testing.T) {
	e := gateExec(t, &gateBinding{verdict: accessor.VerdictIndeterminate})

	got := e.Gate(ctxOf(t), gateName)

	if !got.Refused() {
		t.Fatalf("indeterminate produced the typed verdict %q; it is a "+
			"REFUSAL-class result, never a false allow and never a false deny",
			got.Verdict)
	}
	if got.Refusal.Class != accessor.ClassGateIndeterminate {
		t.Errorf("refusal class = %q; want %q",
			got.Refusal.Class, accessor.ClassGateIndeterminate)
	}
	if got.Verdict == accessor.VerdictAllow || got.Verdict == accessor.VerdictDeny {
		t.Errorf("the refusal also carries verdict %q; indeterminate must not "+
			"collapse into a decided verdict", got.Verdict)
	}
}

// REQ-39: "A gate's timeout or execution failure is an *accessor* refusal
// (§D11, exit 3), not a gate result." — a gate that times out yields
// `timeout`, not `gate_indeterminate`.
// BOUNDARY
//
// The discriminating pair: an undecided gate and a timed-out gate look
// alike to a naive implementation ("we did not get an answer"), but they
// are different classes. Collapsing them retires the timeout diagnosis.
func TestReq39_GateTimeoutAndExecutionFailureAreAccessorRefusalsNotGateResults(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		e := gateExec(t, &gateBinding{
			verdict: accessor.VerdictAllow, delay: 10 * fixtureTimeout,
		})

		got := e.Gate(ctxOf(t), gateName)

		if !got.Refused() {
			t.Fatalf("a timed-out gate produced verdict %q", got.Verdict)
		}
		if got.Refusal.Class != accessor.ClassTimeout {
			t.Errorf("refusal class = %q; want %q — a gate's timeout is an "+
				"ACCESSOR refusal, not a gate result",
				got.Refusal.Class, accessor.ClassTimeout)
		}
	})

	t.Run("execution_failure", func(t *testing.T) {
		e := gateExec(t, &gateBinding{failErr: errBindingFailed})

		got := e.Gate(ctxOf(t), gateName)

		if !got.Refused() {
			t.Fatalf("a failed gate produced verdict %q", got.Verdict)
		}
		if got.Refusal.Class != accessor.ClassExecutionFailure {
			t.Errorf("refusal class = %q; want %q",
				got.Refusal.Class, accessor.ClassExecutionFailure)
		}
	})

	t.Run("negative_control_indeterminate_is_still_its_own_class", func(t *testing.T) {
		e := gateExec(t, &gateBinding{verdict: accessor.VerdictIndeterminate})

		got := e.Gate(ctxOf(t), gateName)

		if got.Refusal == nil || got.Refusal.Class != accessor.ClassGateIndeterminate {
			t.Errorf("an undecided gate produced %+v; want %q — an "+
				"implementation that reports `timeout` for every unanswered "+
				"gate fails here", got.Refusal, accessor.ClassGateIndeterminate)
		}
	})
}
