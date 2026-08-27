package resolve_test

// Phase 3b adversarial coverage for RDR 0001. These tests are written from
// the premise that the implementation is wrong, and they target the
// interaction, precedence, and determinism traps that a clause-by-clause
// (REQ-N) test author would plausibly miss: each REQ read alone is
// satisfied, and the defect only appears where two clauses meet.
//
// Every case below is anchored in RDR 0001 Trade-offs / Failure Modes:
//
//	"Visible failures are typed refusal values: no_match, ambiguous_match,
//	owned_state_unavailable, guard_unevaluable, or unmodeled_outcome.
//	Silent failure would mean the kernel guessed a transition or executed
//	persistence directly; both are prohibited by the normative contract.
//	Diagnosis starts with the input tuple, the refusal kind, and the
//	transition table revision used for that resolution."

import (
	"reflect"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// ---------------------------------------------------------------------------
// ADV-1: a guard-FALSE row's RequiresOwned poisons an otherwise legal
// resolution.
//
// Failure mode (FM): "Visible failures are typed refusal values ...
// owned_state_unavailable". Read with REQ-15 ("the only successful
// selection is exactly one matching edge after guard evaluation") and the
// req-list ASSUMPTION that owned_state_unavailable is raised "when a
// candidate row's evaluation requires an owned tag absent from the
// accessor-produced owned snapshot -- not when a candidate merely fails to
// match on a present tag".
//
// The trap: Resolve computes missingOwned() over ALL match-pattern
// candidates BEFORE any guard runs. A row the guard seam decides FALSE is
// not a surviving candidate and its evaluation never needed that owned tag
// -- yet its RequiresOwned entry still converts a clean exact-one match
// into owned_state_unavailable. Because owned_state_unavailable is NOT an
// escapable class (RDR 0002 restricts escapes to no_match and
// ambiguous_match), the table author has no way to model around it.
//
// This is a wrong-refusal defect, not a missing refusal: the kernel refuses
// a transition it should have planned, and the refusal kind it reports
// names a condition that does not hold for the selected edge.
func TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch(t *testing.T) {
	// Row "live" is the legal edge: unguarded-true, all owned state present.
	// Row "dead" matches the same tag-set but its guard is decided FALSE, so
	// it is not a surviving candidate. It requires an owned tag the accessor
	// snapshot does not carry.
	table := resolve.Table{
		Revision: "rev-adv-1",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.live",
				SourceLocator: "flows/rdr.toml:10",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"status"},
				Guard:         namedGuard("always"),
				NextTags:      []resolve.Tag{{Key: "status", Value: "Final"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			},
			{
				RuleID:        "rdr.dead",
				SourceLocator: "flows/rdr.toml:20",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"never-present"},
				Guard:         namedGuard("never"),
				NextTags:      []resolve.Tag{{Key: "status", Value: "Abandoned"}},
			},
		},
	}

	in := resolve.Input{
		Flow:       "rdr",
		Table:      table,
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Recognized: "successful",
		Guards: fixtureGuards{decided: map[string]bool{
			"always": true,
			"never":  false,
		}},
	}

	got := mustResolve(t, in)

	if got.Refused() {
		t.Fatalf("kernel refused with kind %q (missing owned %v); "+
			"want a plan for rule %q — the only row that required %q was decided "+
			"guard-FALSE and is not a surviving candidate, so its owned-state "+
			"requirement must not decide the disposition",
			got.Refusal.Kind, got.Refusal.MissingOwned, "rdr.live", "never-present")
	}
	if got.Plan.RuleID != "rdr.live" {
		t.Errorf("selected rule = %q; want %q", got.Plan.RuleID, "rdr.live")
	}
}

// ADV-1b: the same ordering trap, stated as a precedence question the RDR
// does answer indirectly. owned_state_unavailable must describe the edge
// the kernel actually wanted to take. Here NO row survives its guard, so
// the honest disposition is no_match (which an escape edge could rescue) —
// not owned_state_unavailable (which nothing can rescue).
//
// FM: "Diagnosis starts with the input tuple, the refusal kind, and the
// transition table revision" — a refusal kind that misdescribes the
// condition is a diagnosis defect, and it also silently removes the table
// author's modeled no_match escape.
func TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable(t *testing.T) {
	table := resolve.Table{
		Revision: "rev-adv-1b",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{
			{
				RuleID:        "rdr.dead",
				SourceLocator: "flows/rdr.toml:20",
				Outcome:       "successful",
				Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
				RequiresOwned: []string{"never-present"},
				Guard:         namedGuard("never"),
				NextTags:      []resolve.Tag{{Key: "status", Value: "Abandoned"}},
			},
			escapeRow("rdr.escape.nomatch", "flows/rdr.toml:90", resolve.KindNoMatch),
		},
	}

	in := resolve.Input{
		Flow:       "rdr",
		Table:      table,
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Recognized: "successful",
		Guards:     fixtureGuards{decided: map[string]bool{"never": false}},
	}

	got := mustResolve(t, in)

	if got.Refused() {
		t.Fatalf("kernel refused with kind %q; want the modeled no_match escape "+
			"edge to rescue — every ordinary row was decided guard-FALSE, which is "+
			"a zero-match condition, and no_match is an escapable class",
			got.Refusal.Kind)
	}
	if got.Plan.RuleID != "rdr.escape.nomatch" || !got.Plan.Escaped {
		t.Errorf("selected rule = %q (escaped=%v); want the modeled escape edge",
			got.Plan.RuleID, got.Plan.Escaped)
	}
}

// ---------------------------------------------------------------------------
// ADV-2: an escape edge's Guard is never evaluated, so an escape silently
// bypasses the RDR 0003 guard seam.
//
// Failure mode (FM): "Silent failure would mean the kernel guessed a
// transition". Read with REQ-15 ("the only successful selection is exactly
// one matching edge after guard evaluation") and REQ-23 ("guard evaluation
// delegation").
//
// The trap: escapeOrRefuse() filters escape rows on rescues(kind), Outcome,
// and Match — and then emits a Plan. It never calls the guard seam. A
// table author who writes an escape edge behind a guard gets a transition
// the guard forbids. The guard-FALSE case is the sharp one: the kernel
// emits a plan for an edge whose own predicate says it does not hold. That
// is precisely "guessed a transition".
func TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam(t *testing.T) {
	t.Run("guard FALSE must not rescue", func(t *testing.T) {
		escape := escapeRow("rdr.escape.guarded", "flows/rdr.toml:90", resolve.KindNoMatch)
		escape.Guard = namedGuard("never")

		in := noMatchInput()
		in.Table.Revision = "rev-adv-2-false"
		in.Table.Rows = append(in.Table.Rows, escape)
		in.Guards = fixtureGuards{decided: map[string]bool{"never": false}}

		got := mustResolve(t, in)

		if !got.Refused() {
			t.Fatalf("kernel emitted a plan for rule %q (escaped=%v); want a "+
				"no_match refusal — the escape edge's own guard was decided FALSE, "+
				"so selecting it is a guessed transition",
				got.Plan.RuleID, got.Plan.Escaped)
		}
		if got.Refusal.Kind != resolve.KindNoMatch {
			t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindNoMatch)
		}
	})

	t.Run("guard UNEVALUABLE must not rescue", func(t *testing.T) {
		escape := escapeRow("rdr.escape.guarded", "flows/rdr.toml:90", resolve.KindNoMatch)
		escape.Guard = namedGuard("unknown-predicate")

		in := noMatchInput()
		in.Table.Revision = "rev-adv-2-unevaluable"
		in.Table.Rows = append(in.Table.Rows, escape)
		// fixtureGuards reports GuardUnevaluable for predicates it does not know.
		in.Guards = fixtureGuards{}

		got := mustResolve(t, in)

		if !got.Refused() {
			t.Fatalf("kernel emitted a plan for rule %q (escaped=%v); want a typed "+
				"refusal — the escape edge's guard could not be decided, so the "+
				"kernel cannot know the edge holds",
				got.Plan.RuleID, got.Plan.Escaped)
		}
		if got.Refusal.Kind != resolve.KindGuardUnevaluable {
			t.Errorf("refusal kind = %q; want %q — an undecidable escape predicate "+
				"is exactly the guard_unevaluable condition",
				got.Refusal.Kind, resolve.KindGuardUnevaluable)
		}
	})
}

// ADV-2b: an escape edge's RequiresOwned is likewise never checked, so an
// escape can emit owned-tag Writes derived from owned state the accessor
// snapshot never produced.
//
// FM: "Silent failure would mean the kernel guessed a transition or
// executed persistence directly". The plan's Writes are the persistence
// description handed to the RDR 0004 accessor layer; describing a write
// off absent owned state is the kernel guessing.
func TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement(t *testing.T) {
	escape := escapeRow("rdr.escape.needsowned", "flows/rdr.toml:90", resolve.KindNoMatch)
	escape.RequiresOwned = []string{"never-present"}

	in := noMatchInput()
	in.Table.Revision = "rev-adv-2b"
	in.Table.Rows = append(in.Table.Rows, escape)

	got := mustResolve(t, in)

	if !got.Refused() {
		t.Fatalf("kernel emitted a plan for rule %q describing writes %v; want "+
			"owned_state_unavailable — the escape edge required owned tag %q, "+
			"which the accessor snapshot does not carry",
			got.Plan.RuleID, got.Plan.Writes, "never-present")
	}
	if got.Refusal.Kind != resolve.KindOwnedStateUnavailable {
		t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindOwnedStateUnavailable)
	}
}

// ---------------------------------------------------------------------------
// ADV-3: the guard_unevaluable refusal is decided by table row order, so
// the reported refusal is not a function of the input tuple's value.
//
// Failure mode (FM): "Diagnosis starts with the input tuple, the refusal
// kind, and the transition table revision used for that resolution." Read
// with REQ-1 ("Given the same ... resolve returns the same disposition")
// and REQ-3 (value-level replay determinism).
//
// The trap: the guard loop returns on the FIRST GuardUnevaluable row it
// meets, so the refusal payload and Refusal.Rows name whichever undecidable row
// happens to sit earlier in Table.Rows. Two tables that are the same SET of
// rows — same revision, same tuple, same everything a reviewer would call
// "the same table" — produce different refusal payloads. Diagnosis then
// depends on a normalization detail RDR 0002 owns rather than on the input.
//
// Phase 1's TestReq3 covers Owned/Observed slice order. It does not cover
// Table.Rows order, which is the surface that actually leaks here.
func TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder(t *testing.T) {
	rowA := resolve.Row{
		RuleID:        "rdr.unevaluable.a",
		SourceLocator: "flows/rdr.toml:10",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		Guard:         namedGuard("unknown-a"),
		NextTags:      []resolve.Tag{{Key: "status", Value: "A"}},
	}
	rowB := resolve.Row{
		RuleID:        "rdr.unevaluable.b",
		SourceLocator: "flows/rdr.toml:20",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"status"},
		Guard:         namedGuard("unknown-b"),
		NextTags:      []resolve.Tag{{Key: "status", Value: "B"}},
	}

	newInput := func(rows ...resolve.Row) resolve.Input {
		return resolve.Input{
			Flow: "rdr",
			Table: resolve.Table{
				Revision: "rev-adv-3",
				Outcomes: []string{"successful"},
				Rows:     rows,
			},
			Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
			Recognized: "successful",
			// fixtureGuards decides neither predicate.
			Guards: fixtureGuards{},
		}
	}

	forward := mustResolve(t, newInput(rowA, rowB))
	reversed := mustResolve(t, newInput(rowB, rowA))

	if !forward.Refused() || !reversed.Refused() {
		t.Fatalf("both orderings must refuse; forward refused=%v reversed refused=%v",
			forward.Refused(), reversed.Refused())
	}
	if forward.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("forward refusal kind = %q; want %q",
			forward.Refusal.Kind, resolve.KindGuardUnevaluable)
	}

	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
		t.Errorf("refusal payload depends on Table.Rows order.\n"+
			" rows [a,b] -> undecided=%+v rows=%v\n"+
			" rows [b,a] -> undecided=%+v rows=%v\n"+
			"The same tuple and the same table revision must replay the same "+
			"disposition; row order is a normalization detail RDR 0002 owns and "+
			"must not reach the reported diagnosis.",
			forward.Refusal.Undecided, forward.Refusal.Rows,
			reversed.Refusal.Undecided, reversed.Refusal.Rows)
	}
}

// ADV-3b: the same order-dependence in the owned_state_unavailable payload.
// missingOwned() accumulates keys in row-then-slice order, so a table whose
// rows are permuted reports the missing keys in a different order. REQ-10
// makes this payload the diagnosis surface; REQ-1 makes the disposition a
// function of the tuple. The implementation's own comment claims
// missingOwned returns keys "in stable order" — stable with respect to row
// order, which is not the same as stable with respect to the input tuple.
func TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder(t *testing.T) {
	rowA := resolve.Row{
		RuleID:        "rdr.needs.a",
		SourceLocator: "flows/rdr.toml:10",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"alpha"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "A"}},
	}
	rowB := resolve.Row{
		RuleID:        "rdr.needs.b",
		SourceLocator: "flows/rdr.toml:20",
		Outcome:       "successful",
		Match:         []resolve.Tag{{Key: "status", Value: "Draft"}},
		RequiresOwned: []string{"beta"},
		NextTags:      []resolve.Tag{{Key: "status", Value: "B"}},
	}

	newInput := func(rows ...resolve.Row) resolve.Input {
		return resolve.Input{
			Flow: "rdr",
			Table: resolve.Table{
				Revision: "rev-adv-3b",
				Outcomes: []string{"successful"},
				Rows:     rows,
			},
			Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
			Recognized: "successful",
		}
	}

	forward := mustResolve(t, newInput(rowA, rowB))
	reversed := mustResolve(t, newInput(rowB, rowA))

	if !forward.Refused() || !reversed.Refused() {
		t.Fatalf("both orderings must refuse; forward refused=%v reversed refused=%v",
			forward.Refused(), reversed.Refused())
	}
	if forward.Refusal.Kind != resolve.KindOwnedStateUnavailable {
		t.Fatalf("forward refusal kind = %q; want %q",
			forward.Refusal.Kind, resolve.KindOwnedStateUnavailable)
	}

	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
		t.Errorf("owned_state_unavailable payload depends on Table.Rows order.\n"+
			" rows [a,b] -> missing=%v rows=%v\n"+
			" rows [b,a] -> missing=%v rows=%v\n"+
			"Diagnosis (REQ-10) must be a function of the input tuple, not of "+
			"row order.",
			forward.Refusal.MissingOwned, forward.Refusal.Rows,
			reversed.Refusal.MissingOwned, reversed.Refusal.Rows)
	}
}

// ---------------------------------------------------------------------------
// Regression-guard territory: cases an adversarial reader expects to break
// but which the implementation gets right. Kept so a later refactor cannot
// silently regress them.

// ADV-4 (guard): observed tags must not be able to shadow the
// accessor-produced owned snapshot.
//
// FM: "Silent failure would mean the kernel guessed a transition". If a
// caller-supplied observed tag could overwrite an owned tag of the same
// key, the caller could steer the kernel onto an edge the artifact's real
// owned state forbids — and the resulting plan would describe owned-tag
// Writes computed off caller input. assemble() applies owned last, so owned
// wins; this test pins that precedence.
func TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot(t *testing.T) {
	in := legalInput()
	in.Table.Revision = "rev-adv-4"
	// The caller claims status:Final; the accessor snapshot says Draft. The
	// single legal row matches on status:Draft.
	in.Observed = append(in.Observed, resolve.Tag{Key: "status", Value: "Final"})

	got := mustResolve(t, in)

	if got.Refused() {
		t.Fatalf("kernel refused with kind %q; want the owned snapshot "+
			"(status:Draft) to decide the match, not the caller's observed "+
			"status:Final", got.Refusal.Kind)
	}

	// The disposition must be identical to the run without the shadowing tag.
	clean := legalInput()
	clean.Table.Revision = "rev-adv-4"
	want := mustResolve(t, clean)

	if !reflect.DeepEqual(got.Plan, want.Plan) {
		t.Errorf("plan changed when a caller-supplied observed tag shadowed an "+
			"owned key.\n got  = %+v\n want = %+v", got.Plan, want.Plan)
	}
}

// ADV-5 (guard): a required owned key that is present ONLY as an observed
// tag must still refuse owned_state_unavailable.
//
// FM: "Visible failures are typed refusal values ... owned_state_unavailable".
// The req-list ASSUMPTION ties the condition to "an owned tag absent from
// the accessor-produced owned snapshot". A caller must not be able to
// satisfy an owned-state requirement by supplying the key as observed
// context — that would let caller input stand in for artifact state.
// missingOwned() checks has(key, ProvenanceOwned), so provenance is
// enforced; this test pins it.
func TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement(t *testing.T) {
	in := missingOwnedInput()
	in.Table.Revision = "rev-adv-5"
	// Supply the required owned key as caller-observed context instead.
	in.Observed = append(in.Observed, resolve.Tag{Key: "gate", Value: "open"})

	got := mustResolve(t, in)

	if !got.Refused() {
		t.Fatalf("kernel emitted a plan for rule %q; want owned_state_unavailable "+
			"— the required key was supplied as an observed tag, which is caller "+
			"context and not accessor-produced owned state", got.Plan.RuleID)
	}
	if got.Refusal.Kind != resolve.KindOwnedStateUnavailable {
		t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindOwnedStateUnavailable)
	}
	if !contains(got.Refusal.MissingOwned, "gate") {
		t.Errorf("MissingOwned = %v; want it to name %q",
			got.Refusal.MissingOwned, "gate")
	}
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
