package resolve

import (
	"errors"
	"testing"
)

// RDR 0009 Phase 3b — ADVERSARIAL failure-mode tests, kernel half.
//
// Written against the RDR's `## Trade-offs > ### Failure Modes` section
// without reading the Phase 1 suite.

func adv0009EscapeWrites(ruleID, locator string) Row {
	return Row{
		RuleID:        ruleID,
		SourceLocator: locator,
		Outcome:       "adv",
		Escape:        []RefusalKind{KindNoMatch},
		Writes:        []Tag{{Key: "owned-a", Value: "1"}},
	}
}

// --- ADV-2 (kernel half): multi-breach traversal fidelity ---------------
//
// FAILURE MODE: "A single errors.As/AsType call reports only the FIRST
// breach in the chain, so it is the wrong instrument for reading a
// multi-breach report — a test asserting on all offending rows must
// traverse the aggregate."
//
// This pins the shape that makes traversal possible: the aggregate is
// EXACTLY ONE LEVEL DEEP, every element a *EscapeShapeBreachError, every
// element itself classifiable by errors.Is (so Unwrap on the element
// returns the sentinel, not nil), and errors.As over the WHOLE aggregate
// reaches only the first. Any change that nests the join, or that stops an
// element unwrapping to the sentinel, breaks every documented consumer.
func TestAdv0009_TheAggregateIsFlatAndEveryElementClassifies(t *testing.T) {
	tbl := Table{Outcomes: []string{"adv"}, Rows: []Row{
		adv0009EscapeWrites("rule-c", "flow.toml:3"),
		adv0009EscapeWrites("rule-a", "flow.toml:1"),
		adv0009EscapeWrites("rule-b", "flow.toml:2"),
	}}

	err := tbl.CheckValid()
	if err == nil {
		t.Fatal("CheckValid returned nil for a breaching table")
	}
	if !errors.Is(err, ErrEscapeShapeBreach) {
		t.Fatal("the aggregate does not classify as ErrEscapeShapeBreach")
	}

	var joined interface{ Unwrap() []error }
	if !errors.As(err, &joined) {
		t.Fatal("the breach error does not expose Unwrap() []error")
	}
	elems := joined.Unwrap()
	if len(elems) != 3 {
		t.Fatalf("aggregate has %d elements, want 3", len(elems))
	}

	wantOrder := []RowRef{
		{RuleID: "rule-a", SourceLocator: "flow.toml:1"},
		{RuleID: "rule-b", SourceLocator: "flow.toml:2"},
		{RuleID: "rule-c", SourceLocator: "flow.toml:3"},
	}
	for i, elem := range elems {
		// Every element must classify on its own, not only the aggregate.
		if !errors.Is(elem, ErrEscapeShapeBreach) {
			t.Errorf("element %d does not unwrap to the sentinel", i)
		}
		// Every element must be exactly one *EscapeShapeBreachError, and
		// must NOT itself be a nested aggregate.
		var nested interface{ Unwrap() []error }
		if errors.As(elem, &nested) {
			t.Errorf("element %d is itself an aggregate; the flat shape the "+
				"documented traversal depends on is broken", i)
		}
		breach, ok := elem.(*EscapeShapeBreachError)
		if !ok {
			t.Fatalf("element %d is %T, want *EscapeShapeBreachError", i, elem)
		}
		if breach.Ref != wantOrder[i] {
			t.Errorf("element %d Ref = %v, want %v", i, breach.Ref, wantOrder[i])
		}
		if breach.Count != 1 {
			t.Errorf("element %d Count = %d, want 1", i, breach.Count)
		}
	}

	// The RDR's own warning, pinned: a single errors.As over the aggregate
	// reports the FIRST breach alone, so it is the wrong instrument.
	var first *EscapeShapeBreachError
	if !errors.As(err, &first) {
		t.Fatal("errors.As found no breach at all")
	}
	if first.Ref != wantOrder[0] {
		t.Errorf("a single errors.As returned %v, want the first element %v",
			first.Ref, wantOrder[0])
	}
}

// --- ADV-3 (kernel half): no disposition change for conforming tables ---
//
// FAILURE MODE / Consequences: the whole-table entry precondition is a
// deliberate behavior change for BREACHING tables only. "zero disposition
// change for conforming tables" — a conforming table must reach exactly the
// disposition it reached before the precondition existed, including the
// escape-rescue path, and including the arms where the precondition is the
// SECOND full Table.Rows pass.
func TestAdv0009_AConformingTableKeepsItsDisposition(t *testing.T) {
	ordinary := Row{
		RuleID: "ordinary", SourceLocator: "flow.toml:1", Outcome: "adv",
		Match:    []Tag{{Key: "stage", Value: "propose"}},
		NextTags: []Tag{{Key: "stage", Value: "prelock"}},
		Writes:   []Tag{{Key: "stage", Value: "prelock"}},
	}
	// A CONFORMING escape row: non-empty Escape, empty Writes. It must
	// still be able to rescue.
	escape := Row{
		RuleID: "escape", SourceLocator: "flow.toml:2", Outcome: "adv",
		Escape:   []RefusalKind{KindNoMatch},
		NextTags: []Tag{{Key: "stage", Value: "parked"}},
	}
	tbl := Table{Outcomes: []string{"adv"}, Rows: []Row{ordinary, escape}}

	if err := tbl.CheckValid(); err != nil {
		t.Fatalf("a conforming table failed CheckValid: %v", err)
	}

	cases := []struct {
		name       string
		observed   []Tag
		wantPlan   string
		wantEscape bool
	}{
		{name: "ordinary match", observed: []Tag{{Key: "stage", Value: "propose"}},
			wantPlan: "ordinary", wantEscape: false},
		{name: "no ordinary match rescues via the escape row",
			observed: []Tag{{Key: "stage", Value: "elsewhere"}},
			wantPlan: "escape", wantEscape: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Resolve(Input{
				Table: tbl, Observed: tc.observed, Recognized: "adv",
			})
			if err != nil {
				t.Fatalf("the precondition errored on a conforming table: %v", err)
			}
			if got.Plan == nil {
				t.Fatalf("want a plan for %q, got refusal %+v", tc.wantPlan, got.Refusal)
			}
			if got.Plan.RuleID != tc.wantPlan {
				t.Errorf("plan = %q, want %q", got.Plan.RuleID, tc.wantPlan)
			}
			if got.Plan.Escaped != tc.wantEscape {
				t.Errorf("Escaped = %v, want %v", got.Plan.Escaped, tc.wantEscape)
			}
		})
	}

	// An empty table conforms VACUOUSLY and still reaches its modeled
	// disposition rather than an error.
	if err := (Table{}).CheckValid(); err != nil {
		t.Errorf("an empty table must conform vacuously, got %v", err)
	}
	got, err := Resolve(Input{
		Table: Table{Outcomes: []string{"adv"}}, Recognized: "adv",
	})
	if err != nil {
		t.Fatalf("empty conforming table errored: %v", err)
	}
	if got.Refusal == nil || got.Refusal.Kind != KindNoMatch {
		t.Errorf("want no_match on an empty conforming table, got %+v", got)
	}
}

// A breach must outrank every modeled disposition, unmodeled_outcome
// included, and must return the ZERO Result — neither a Plan nor a
// Refusal ("carries the zero Result — neither disposition").
func TestAdv0009_ABreachOutranksEveryModeledDispositionAndCarriesZeroResult(t *testing.T) {
	tbl := Table{Outcomes: []string{"adv"}, Rows: []Row{
		adv0009EscapeWrites("dormant", "flow.toml:9"),
	}}

	// "unmodeled-outcome" is NOT in the alphabet, so without the
	// precondition this resolves to a modeled unmodeled_outcome refusal.
	got, err := Resolve(Input{Table: tbl, Recognized: "not-in-alphabet"})
	if err == nil {
		t.Fatal("a breaching table resolved cleanly on the unmodeled-outcome arm")
	}
	if !errors.Is(err, ErrEscapeShapeBreach) {
		t.Fatalf("err = %v, want an escape-shape breach", err)
	}
	if got.Plan != nil || got.Refusal != nil {
		t.Errorf("a breach must carry the zero Result, got %+v", got)
	}
}
