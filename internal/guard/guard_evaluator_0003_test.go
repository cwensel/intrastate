package guard_test

// RDR 0003 — evaluator scope and the kernel split: this RDR's evaluator
// decides value semantics over a PRESENT VALUE only, and never reads the
// tag view.

import (
	"reflect"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-34: "This RDR's guard evaluator MUST decide value semantics over a
// present value only. Key presence, existence atoms, absent-key
// unevaluability, and the combination of per-atom verdicts belong to the
// kernel (JDR 0001 §D4, normative in RDR 0007); the evaluator MUST NOT
// read the tag view."
// BOUNDARY
func TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView(t *testing.T) {
	// The seam signature is the proof it cannot read the view: it takes an
	// atom and a value, and RDR 0007's GuardEvaluator has exactly one
	// method. An evaluator that could read the view would need a TagSet
	// parameter or a stored reference.
	var ev guard.Evaluator
	var _ resolve.GuardEvaluator = ev

	method, ok := reflect.TypeOf(ev).MethodByName("Evaluate")
	if !ok {
		t.Fatal("guard.Evaluator declares no Evaluate method")
	}
	fn := method.Type
	// receiver + atom + value
	if fn.NumIn() != 3 {
		t.Fatalf("Evaluate takes %d parameters (receiver included); want 3 — "+
			"atom and value only, never the tag view", fn.NumIn())
	}
	viewType := reflect.TypeOf(resolve.TagSet{})
	for i := range fn.NumIn() {
		if fn.In(i) == viewType {
			t.Errorf("Evaluate parameter %d is a resolve.TagSet; the evaluator "+
				"MUST NOT read the tag view", i)
		}
	}
	if reflect.TypeOf(ev).Kind() == reflect.Struct && reflect.TypeOf(ev).NumField() != 0 {
		t.Errorf("guard.Evaluator carries %d fields; a view-free evaluator "+
			"holds no state to read one from", reflect.TypeOf(ev).NumField())
	}

	// Existence atoms belong to the kernel, so the evaluator does not decide
	// them: handed one, it answers unevaluable rather than inventing a
	// presence verdict it has no view to read.
	exists := resolve.GuardAtom{
		Key: "k", Operator: resolve.OpExists, Literal: resolve.LiteralTrue, Block: resolve.BlockAll,
	}
	if got := ev.Evaluate(exists, "anything"); got != resolve.GuardUnevaluable {
		t.Errorf("Evaluate over an existence atom = %v; want GuardUnevaluable "+
			"— existence atoms belong to the kernel", got)
	}

	// "MUST decide value semantics over a present value" is the positive
	// half: a seam that answers unevaluable to everything reads no view
	// either, and would pass the clauses above while deciding nothing. The
	// cross-RDR contract test drives the present-value × literal × operator
	// product and is the discriminating oracle.
	resolve.TestGuardEvaluatorContract(t, ev)
}

// REQ-35: "an existence atom over an absent key is decided — not
// unevaluable" and "a value atom over an absent key is unevaluable rather
// than false"
// DOMAIN EDGE
func TestReq35_ExistenceOverAbsentIsDecidedAndValueOverAbsentIsUnevaluable(t *testing.T) {
	tbl := resolve.Table{
		Revision: "rev",
		Outcomes: []string{"go"},
		Rows: []resolve.Row{{
			RuleID: "value-over-absent", SourceLocator: "f:1", Outcome: "go",
			Guard: []resolve.GuardAtom{kernelAtom("missing", "eq", "x")},
		}},
	}
	res, err := resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go", Guards: guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("a value atom over an absent key gave %+v; want a "+
			"guard_unevaluable refusal — unevaluable rather than false", res)
	}

	// An existence atom over the same absent key is DECIDED: `exists = false`
	// over an absent key is true, so the row qualifies and plans.
	tbl.Rows = []resolve.Row{{
		RuleID: "exists-over-absent", SourceLocator: "f:2", Outcome: "go",
		Guard: []resolve.GuardAtom{{
			Key: "missing", Operator: resolve.OpExists,
			Literal: resolve.LiteralFalse, Block: resolve.BlockAll,
		}},
	}}
	res, err = resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go", Guards: guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Fatalf("an existence atom over an absent key refused %v; it is "+
			"DECIDED, never unevaluable", res.Refusal.Kind)
	}

	// Discriminating leg: the same value atom over a PRESENT key must
	// DECIDE, so the refusal above is attributable to absence rather than to
	// a seam that cannot decide anything.
	tbl.Rows = []resolve.Row{{
		RuleID: "value-over-present", SourceLocator: "f:3", Outcome: "go",
		Guard: []resolve.GuardAtom{kernelAtom("present", "eq", "x")},
	}}
	res, err = resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go",
		Observed: []resolve.Tag{{Key: "present", Value: "x"}},
		Guards:   guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Fatalf("a value atom over a PRESENT key refused %v; only an ABSENT "+
			"key makes a value atom unevaluable", res.Refusal.Kind)
	}
}

// REQ-36: "an author expressing \"this row applies when tag X is absent\"
// writes an existence atom, not a value atom that happens to miss."
// DOMAIN EDGE
func TestReq36_AbsenceIsExpressedWithAnExistenceAtomNotAMissingValueAtom(t *testing.T) {
	tbl := resolve.Table{Revision: "rev", Outcomes: []string{"go"}}

	// The correct spelling: `exists = false`. It plans.
	tbl.Rows = []resolve.Row{{
		RuleID: "absent-correct", SourceLocator: "f:1", Outcome: "go",
		Guard: []resolve.GuardAtom{{
			Key: "x", Operator: resolve.OpExists,
			Literal: resolve.LiteralFalse, Block: resolve.BlockAll,
		}},
	}}
	res, err := resolve.Resolve(resolve.Input{Table: tbl, Recognized: "go", Guards: guard.Evaluator{}})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Errorf("the existence spelling of \"tag X is absent\" refused %v; "+
			"it is the spelling the RDR names", res.Refusal.Kind)
	}

	// The incorrect spelling: a value atom that happens to miss. It refuses,
	// so an author cannot reach absence through it.
	tbl.Rows = []resolve.Row{{
		RuleID: "absent-wrong", SourceLocator: "f:2", Outcome: "go",
		Guard: []resolve.GuardAtom{kernelAtom("x", "eq", "")},
	}}
	res, err = resolve.Resolve(resolve.Input{Table: tbl, Recognized: "go", Guards: guard.Evaluator{}})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("a value atom over an absent key gave %+v; a value atom that "+
			"happens to miss is not how absence is expressed", res)
	}

	// Discriminating leg: the same value atom over a PRESENT key that does
	// NOT equal the literal decides FALSE and prunes, so the refusal above
	// is attributable to absence rather than to an undecidable seam.
	tbl.Rows = []resolve.Row{
		{
			RuleID: "absent-wrong", SourceLocator: "f:2", Outcome: "go",
			Guard: []resolve.GuardAtom{kernelAtom("x", "eq", "yes")},
		},
		{RuleID: "fallback", SourceLocator: "f:3", Outcome: "go"},
	}
	res, err = resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go",
		Observed: []resolve.Tag{{Key: "x", Value: "no"}},
		Guards:   guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Fatalf("a value atom over a present, non-matching value refused %v; "+
			"it decides FALSE and prunes its row", res.Refusal.Kind)
	}
	if res.Plan.RuleID != "fallback" {
		t.Errorf("the plan names %q; the non-matching row must be pruned, "+
			"leaving `fallback`", res.Plan.RuleID)
	}
}

// REQ-37: Per the `authority` census, the canonical writer of a value-atom
// verdict over a present value is "This RDR's evaluator"; key presence,
// the existence-atom verdict, and per-atom verdict combination are the
// kernel's, and this RDR's evaluator is "**explicitly not an arm**; it
// never reads the tag view".
// BOUNDARY
func TestReq37_ValueVerdictIsThisRDRsAndPresenceCombinationIsTheKernels(t *testing.T) {
	// The evaluator writes the value-atom verdict: swapping it changes the
	// disposition, which proves the kernel delegates rather than deciding.
	tbl := resolve.Table{
		Revision: "rev", Outcomes: []string{"go"},
		Rows: []resolve.Row{{
			RuleID: "r", SourceLocator: "f:1", Outcome: "go",
			Guard: []resolve.GuardAtom{kernelAtom("k", "eq", "yes")},
		}},
	}
	in := resolve.Input{
		Table: tbl, Recognized: "go",
		Observed: []resolve.Tag{{Key: "k", Value: "yes"}},
		Guards:   guard.Evaluator{},
	}
	res, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Fatalf("this RDR's evaluator did not decide a matching present "+
			"value: %v", res.Refusal.Kind)
	}

	// Presence and combination stay the kernel's: with the SAME evaluator,
	// an absent key refuses without the evaluator ever being consulted, and
	// two atoms combine to a single row verdict.
	in.Observed = nil
	res, err = resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("key presence was not decided by the kernel: %+v", res)
	}
}
