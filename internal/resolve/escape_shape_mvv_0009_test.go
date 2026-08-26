package resolve_test

// RDR 0009 — the Minimum Viable Validation and the record's executable test
// scenarios.
//
// REQ-70 is the gating MVV: one runnable end-to-end test carrying both legs
// of the record's claim — the breaching table errors and names the row, and
// the SAME table with the writes removed resolves, escapes, and refuses
// exactly as before. The RDR declares no encode/decode or inverse operation,
// so "zero disposition change" is asserted as VALUE equality of the
// disposition against the conformed baseline, never as a green exit.

import (
	"errors"
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-70: "A hand-constructed table containing an escape row with populated
// `Writes`: `Resolve` returns a non-nil error naming that row's `RuleID` and
// `SourceLocator`, with no `Plan` and no `Refusal` — while the same table with
// the writes removed resolves, escapes, and refuses exactly as the frozen ADV
// suite proves today (zero disposition change for conforming tables). The
// exported validator and the entry check are exercised as the same predicate."
// HAPPY PATH — REQ-MVV
//
// Four obligations in one runnable suite:
//
//  1. the breaching table errors, naming RuleID and SourceLocator as
//     structure, with the zero Result;
//  2. the SAME table with the writes removed resolves — and its disposition
//     is VALUE-EQUAL to the pre-change baseline the frozen suite proves, not
//     merely non-erroring;
//  3. the conformed table still ESCAPES and still REFUSES where it did;
//  4. the exported validator and the entry check are the same predicate,
//     verdict for verdict.
func TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel(t *testing.T) {
	t.Run("1_breaching_table_errors_naming_the_row", func(t *testing.T) {
		got, err := resolve.Resolve(breachingNoMatchInput())
		if err == nil {
			t.Fatalf("the breaching table resolved: %+v", got)
		}
		if got.Plan != nil {
			t.Errorf("a Plan accompanied the breach: %+v", got.Plan)
		}
		if got.Refusal != nil {
			t.Errorf("a Refusal accompanied the breach: %+v", got.Refusal)
		}

		var typed *resolve.EscapeShapeBreachError
		if !errors.As(err, &typed) {
			t.Fatalf("no *EscapeShapeBreachError in the chain (%T)", err)
		}
		if typed.Ref.RuleID != breachRuleID {
			t.Errorf("RuleID = %q; want %q", typed.Ref.RuleID, breachRuleID)
		}
		if typed.Ref.SourceLocator != breachLocator {
			t.Errorf("SourceLocator = %q; want %q",
				typed.Ref.SourceLocator, breachLocator)
		}
	})

	t.Run("2_conformed_table_resolves_value_for_value", func(t *testing.T) {
		// The baseline: the conformed table's disposition, computed by the
		// same kernel. "Zero disposition change" means the value the escape
		// path produces is exactly the escape row's own plan — asserted
		// field by field, not as "did not error".
		got, err := resolve.Resolve(conformingNoMatchInput())
		if err != nil {
			t.Fatalf("the conformed table errored: %v", err)
		}
		if got.Plan == nil {
			t.Fatalf("the conformed table produced no plan: %+v", got)
		}

		want := &resolve.Plan{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			NextTags:      []resolve.Tag{{Key: "status", Value: "Blocked"}},
			Writes:        nil,
			Revision:      conformingNoMatchInput().Table.Revision,
			Escaped:       true,
		}
		if !reflect.DeepEqual(got.Plan, want) {
			t.Errorf("escaped plan = %+v; want %+v — the conformed table's "+
				"disposition must be value-identical, not merely non-erroring",
				got.Plan, want)
		}
		if got.Refusal != nil {
			t.Errorf("the conformed escape produced a refusal too: %+v",
				got.Refusal)
		}
	})

	t.Run("3_conformed_tables_still_escape_and_still_refuse", func(t *testing.T) {
		// Escapes: the rescue path still fires.
		escaped := mustResolve(t, conformingNoMatchInput())
		if escaped.Plan == nil || !escaped.Plan.Escaped {
			t.Errorf("the conformed table no longer escapes: %+v", escaped)
		}

		// Refuses: a table with NO escape row still refuses no_match, with a
		// value-identical refusal.
		before := mustResolve(t, noMatchInput())
		if before.Refusal == nil ||
			before.Refusal.Kind != resolve.KindNoMatch {
			t.Fatalf("the no-match fixture no longer refuses no_match: %+v",
				before)
		}
		after := mustResolve(t, noMatchInput())
		if !reflect.DeepEqual(before, after) {
			t.Errorf("the refusal is not value-stable across replay: "+
				"%+v vs %+v", before, after)
		}

		// And every shipped refusal fixture still travels the Result value
		// with a nil error.
		for name, in := range allRefusalInputs() {
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Errorf("%s: a conforming refusal used the Go error "+
					"path: %v", name, err)
				continue
			}
			if got.Refusal == nil {
				t.Errorf("%s: no longer refuses: %+v", name, got)
			}
		}
	})

	t.Run("4_validator_and_entry_check_are_one_predicate", func(t *testing.T) {
		for name, in := range map[string]resolve.Input{
			"breaching":   breachingNoMatchInput(),
			"conformed":   conformingNoMatchInput(),
			"empty_slice": emptyNotNilEscapeInput(),
		} {
			direct := in.Table.CheckValid()
			_, entry := resolve.Resolve(in)
			if (direct == nil) != (entry == nil) {
				t.Errorf("%s: the exported validator and the entry check "+
					"disagree: CheckValid=%v Resolve=%v", name, direct, entry)
				continue
			}
			if direct == nil {
				continue
			}
			if !sameRefs(breachRefs(t, direct), breachRefs(t, entry)) {
				t.Errorf("%s: the two call sites report different "+
					"identities", name)
			}
		}
	})
}

// REQ-71: "The three scenarios below are **normative fixtures** (approved at
// Stage 4). Their values are read from the frozen suite and the A3 spike runs,
// not invented"
// BOUNDARY
//
// The fixture VALUES are themselves normative, so a silent edit to them is a
// contract change. Pinned as literals.
func TestReq71_TheNormativeFixtureValuesArePinned(t *testing.T) {
	if breachRuleID != "rdr.escape.needsowned" {
		t.Errorf("breach rule id = %q; the normative fixture fixes "+
			"rdr.escape.needsowned", breachRuleID)
	}
	if breachLocator != "flows/rdr.toml:90" {
		t.Errorf("breach locator = %q; the normative fixture fixes "+
			"flows/rdr.toml:90", breachLocator)
	}
	want := []resolve.Tag{{Key: "status", Value: "Escaped"}}
	if !reflect.DeepEqual(breachWrites(), want) {
		t.Errorf("breach writes = %+v; want %+v", breachWrites(), want)
	}
}

// REQ-72: "**breach-yields-error-not-plan** — an escape row (`RuleID`
// `rdr.escape.needsowned`, `SourceLocator` `flows/rdr.toml:90`, `Escape`
// `[no_match]`) carrying `Writes: []Tag{{Key: "status", Value: "Escaped"}}`,
// in a table where it is selected via the escape path. **Expected**: `Resolve`
// returns the zero `Result` (`Plan: nil, Refusal: nil`) and a non-nil error;
// the error yields that row's `RowRef{"rdr.escape.needsowned",
// "flows/rdr.toml:90"}` through `errors.As`/`AsType` without parsing text, and
// `errors.Is` classifies it as `ErrEscapeShapeBreach` — both through the
// `errors.Join` aggregate, which wraps even this single-breach case. The
// aggregate's `Unwrap() []error` has length **1**; its single
// `*EscapeShapeBreachError` carries `Count == 1`."
// HAPPY PATH — TS scenario 1
func TestReq72_BreachYieldsErrorNotPlan(t *testing.T) {
	in := breachingNoMatchInput()

	// The row IS selected via the escape path when the precondition is
	// absent: with the writes removed the same table escapes to that row.
	baseline := mustResolve(t, conformingNoMatchInput())
	if baseline.Plan == nil || !baseline.Plan.Escaped ||
		baseline.Plan.RuleID != breachRuleID {
		t.Fatalf("the fixture row is not the one selected via the escape "+
			"path: %+v", baseline)
	}

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("Resolve returned no error: %+v", got)
	}
	if !reflect.DeepEqual(got, resolve.Result{}) {
		t.Errorf("Result = %+v; want the zero Result", got)
	}

	var typed *resolve.EscapeShapeBreachError
	if !errors.As(err, &typed) {
		t.Fatalf("errors.As found no *EscapeShapeBreachError")
	}
	want := resolve.RowRef{RuleID: breachRuleID, SourceLocator: breachLocator}
	if typed.Ref != want {
		t.Errorf("Ref = %+v; want %+v", typed.Ref, want)
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("errors.Is did not classify the aggregate")
	}

	elems := breachElements(t, err)
	if len(elems) != 1 {
		t.Fatalf("Unwrap() []error has length %d; want 1", len(elems))
	}
	if elems[0].Count != 1 {
		t.Errorf("Count = %d; want 1", elems[0].Count)
	}
}

// REQ-73: "**dormant-row-still-errors** — a table that resolves cleanly on its
// own, plus an escape row no resolution path reaches (an outcome outside the
// tuple's reach) carrying the same `Writes`. **Expected**: the same error,
// naming the dormant row — pinning the whole-table decision against an
// on-path-only reading."
// ADVERSARIAL — TS scenario 2
func TestReq73_DormantRowStillErrors(t *testing.T) {
	// The table resolves cleanly WITHOUT the dormant row.
	clean := mustResolve(t, legalInput())
	if clean.Plan == nil {
		t.Fatalf("the base table does not resolve cleanly: %+v", clean)
	}

	got, err := resolve.Resolve(dormantBreachInput())
	if err == nil {
		t.Fatalf("the dormant breach did not error: %+v", got)
	}
	if !reflect.DeepEqual(got, resolve.Result{}) {
		t.Errorf("Result = %+v; want the zero Result", got)
	}

	// The SAME error shape as the on-path case: same sentinel, same typed
	// element carrying the dormant row's identity.
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("the dormant breach is not classified by the sentinel")
	}
	refs := breachRefs(t, err)
	want := resolve.RowRef{RuleID: breachRuleID, SourceLocator: breachLocator}
	if len(refs) != 1 || refs[0] != want {
		t.Errorf("identities = %+v; want the dormant row %+v", refs, want)
	}
}

// REQ-74: "**empty-not-nil-conforms** — an escape row whose `Writes` is a
// non-nil, zero-length slice (`[]Tag{}`). **Expected**: no error; the table
// resolves, escapes, and refuses identically to the A3 baseline, and the
// emitted plan carries `len(Plan.Writes) == 0`. Asserted on length, never on
// nil-ness"
// INPUT EDGE — TS scenario 3
func TestReq74_EmptyNotNilWritesConforms(t *testing.T) {
	nonVacuityGate(t)

	in := emptyNotNilEscapeInput()

	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("a non-nil zero-length Writes errored: %v", err)
	}
	if got.Plan == nil {
		t.Fatalf("the table produced no plan: %+v", got)
	}
	if !got.Plan.Escaped {
		t.Errorf("the plan is not marked Escaped; the escape path did not run")
	}
	// Asserted on LENGTH, never on nil-ness.
	if len(got.Plan.Writes) != 0 {
		t.Errorf("len(Plan.Writes) = %d; want 0", len(got.Plan.Writes))
	}

	// Identical to the conformed baseline: same disposition, value for value.
	// REQ-74 requires the table resolve, escape, and refuse "identically to
	// the A3 baseline", so the comparison is a whole-Plan equality, not a
	// sample of fields — SourceLocator, NextTags, and Revision regressions
	// must not slip through. The ONLY permitted difference is the
	// nil-versus-empty `Writes` distinction the scenario exists to admit, so
	// that one field is normalized away first ("asserted on length, never on
	// nil-ness"); the discriminating length assertion above still stands.
	baseline := mustResolve(t, conformingNoMatchInput())
	if got.Plan == nil || baseline.Plan == nil {
		t.Fatalf("both plans must be non-nil to compare: got=%v baseline=%v",
			got.Plan, baseline.Plan)
	}
	// Copy the pointees before normalizing: Plan is a *Plan, so zeroing
	// Writes through the pointer would mutate the resolver's own value.
	gotPlan, basePlan := *got.Plan, *baseline.Plan
	if len(gotPlan.Writes) == 0 {
		gotPlan.Writes = nil
	}
	if len(basePlan.Writes) == 0 {
		basePlan.Writes = nil
	}
	if !reflect.DeepEqual(gotPlan, basePlan) {
		t.Errorf("the empty-slice table's disposition differs from the "+
			"conformed baseline: %+v vs %+v", gotPlan, basePlan)
	}
}

// REQ-75: "The **discriminating** assertion is on the *input row*
// (`row.Writes != nil && len(row.Writes) == 0` yields no error): asserting only
// on the plan would not separate the two states this scenario exists to
// separate"
// BOUNDARY — TS scenario 3
//
// The input row is asserted to be genuinely in the non-nil-empty state — a
// fixture that had drifted to nil would make REQ-74 vacuous — and THAT state
// is what yields no error.
func TestReq75_TheDiscriminatingAssertionIsOnTheInputRow(t *testing.T) {
	nonVacuityGate(t)

	in := emptyNotNilEscapeInput()

	var escape *resolve.Row
	for i := range in.Table.Rows {
		if in.Table.Rows[i].RuleID == breachRuleID {
			escape = &in.Table.Rows[i]
		}
	}
	if escape == nil {
		t.Fatalf("the fixture no longer carries the escape row")
	}
	if escape.Writes == nil {
		t.Fatalf("the input row's Writes is nil; the scenario requires the " +
			"NON-NIL, zero-length state to discriminate")
	}
	if len(escape.Writes) != 0 {
		t.Fatalf("the input row's Writes has length %d; want 0",
			len(escape.Writes))
	}

	if err := in.Table.CheckValid(); err != nil {
		t.Errorf("CheckValid returned %v for a non-nil, zero-length Writes; "+
			"the predicate tests emptiness, never nil-ness", err)
	}
	if _, err := resolve.Resolve(in); err != nil {
		t.Errorf("Resolve returned %v for a non-nil, zero-length Writes", err)
	}
}

// REQ-76: "Breach precedence over each modeled disposition, run once **per
// refusal kind** — a breaching row coexisting with a table state that would
// otherwise yield `unmodeled_outcome`, `no_match`, `ambiguous_match`,
// `owned_state_unavailable`, and `guard_unevaluable` respectively.
// **Expected**: The breach error wins in every case, with no `Result`
// disposition. Exhaustive over the five kinds"
// ADVERSARIAL — TS scenario 4
//
// Exhaustive by construction: the loop is driven by the shipped refusal
// fixtures and the run fails if the set is not the full five kinds, so a
// future sixth kind cannot slip past unexercised.
func TestReq76_TheBreachPrecedesEveryModeledDisposition(t *testing.T) {
	inputs := allRefusalInputs()
	if len(inputs) != 5 {
		t.Fatalf("the refusal fixture set carries %d kinds; the scenario is "+
			"exhaustive over the five", len(inputs))
	}

	seen := map[resolve.RefusalKind]bool{}
	for name, in := range inputs {
		t.Run(name, func(t *testing.T) {
			// Establish which kind this fixture yields untouched.
			base := mustResolve(t, in)
			if base.Refusal == nil {
				t.Fatalf("the fixture no longer refuses: %+v", base)
			}
			seen[base.Refusal.Kind] = true

			// Add a breaching row that is NOT itself a candidate for this
			// request, so the breach — not a changed selection — is what
			// alters the answer.
			breaching := deepCopyInput(in)
			row := breachingEscapeRow(breachRuleID, breachLocator,
				resolve.KindAmbiguousMatch)
			row.Outcome = "outcome-no-request-carries"
			breaching.Table.Rows = append(breaching.Table.Rows, row)

			got, err := resolve.Resolve(breaching)
			if err == nil {
				t.Fatalf("%v won over the breach: %+v — the breach error "+
					"wins in every case", base.Refusal.Kind, got)
			}
			if got.Plan != nil || got.Refusal != nil {
				t.Errorf("a disposition accompanied the breach: %+v", got)
			}
			if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
				t.Errorf("the returned error is not the breach: %v", err)
			}
		})
	}

	for _, kind := range resolve.RefusalKinds() {
		if !seen[kind] {
			t.Errorf("refusal kind %q was never exercised; the scenario is "+
				"exhaustive over the five kinds", kind)
		}
	}
}

// REQ-77: "The full frozen `internal/resolve` suite (ADV/MVV/boundary) run
// against conformed escape fixtures. **Expected**: Every assertion still
// passes."
// HAPPY PATH — TS scenario 5
//
// Executable half: the shared `escapeRow` builder — inherited by sixteen call
// sites across the frozen suite — must itself CONFORM, so every fixture the
// frozen suite builds from it passes the entry precondition. A builder still
// carrying writes makes the frozen suite fatal at `mustResolve`, which is
// exactly what REQ-95 predicts.
func TestReq77_TheFrozenSuitesEscapeFixturesConform(t *testing.T) {
	row := escapeRow(breachRuleID, breachLocator, resolve.KindNoMatch)
	if len(row.Writes) != 0 {
		t.Errorf("the shared escapeRow builder still carries %d writes; "+
			"every frozen-suite fixture inheriting it would fatal at "+
			"Resolve entry", len(row.Writes))
	}

	tbl := resolve.Table{
		Revision: "rev-frozen",
		Outcomes: []string{"successful"},
		Rows:     []resolve.Row{row},
	}
	if err := tbl.CheckValid(); err != nil {
		t.Errorf("the shared escapeRow builder produces a breaching row: %v",
			err)
	}
}

// REQ-78: "The comparison baseline is the A3 spike's **conformed** run
// (`a3-conformed.out`: 154 PASS, 0 FAIL), restricted to the tests existing at
// that revision — **not** `a3-baseline.out` (the pre-conformance tree) and not
// the raw post-Phase-2 count, which necessarily exceeds 154 because this RDR
// adds tests. The oracle is the sorted outcome set over that pre-existing test
// set, which must be identical."
// DOMAIN EDGE — TS scenario 5
//
// The executable substance of the oracle: over the pre-existing fixture set,
// every disposition must still be produced without the Go error path — the
// outcome set is identical, and no fixture migrated onto the error channel.
func TestReq78_ThePreExistingFixtureSetProducesAnIdenticalOutcomeSet(t *testing.T) {
	nonVacuityGate(t)

	type outcome struct {
		planned bool
		kind    resolve.RefusalKind
		ruleID  string
	}

	first := map[string]outcome{}
	for name, in := range allDispositionInputs() {
		got, err := resolve.Resolve(in)
		if err != nil {
			t.Fatalf("%s: a pre-existing fixture migrated onto the Go error "+
				"path: %v", name, err)
		}
		o := outcome{}
		switch {
		case got.Plan != nil:
			o.planned, o.ruleID = true, got.Plan.RuleID
		case got.Refusal != nil:
			o.kind = got.Refusal.Kind
		default:
			t.Fatalf("%s: neither disposition: %+v", name, got)
		}
		first[name] = o
	}

	for name, in := range allDispositionInputs() {
		got, err := resolve.Resolve(in)
		if err != nil {
			t.Errorf("%s: replay errored: %v", name, err)
			continue
		}
		o := outcome{}
		if got.Plan != nil {
			o.planned, o.ruleID = true, got.Plan.RuleID
		} else if got.Refusal != nil {
			o.kind = got.Refusal.Kind
		}
		if o != first[name] {
			t.Errorf("%s: outcome changed: %+v vs %+v", name, o, first[name])
		}
	}
}

// REQ-79: "The baseline is regenerable, not frozen to the spike artifact:
// re-run the frozen suite at the implementation's base commit immediately
// before Phase 1 lands and diff sorted outcome sets over the tests existing at
// that commit"
// DOMAIN EDGE — TS scenario 5
//
// Regenerability is the property: the outcome set must be a pure function of
// the fixtures, so re-deriving it in a fresh run yields the same set. A
// disposition that depended on run order or on iteration order would break
// this and make the oracle un-regenerable.
func TestReq79_TheOutcomeSetIsRegenerableRatherThanFrozenToAnArtifact(t *testing.T) {
	nonVacuityGate(t)

	derive := func() map[string]string {
		out := map[string]string{}
		for name, in := range allDispositionInputs() {
			got, err := resolve.Resolve(in)
			if err != nil {
				out[name] = "error"
				continue
			}
			if got.Plan != nil {
				out[name] = "plan:" + got.Plan.RuleID
				continue
			}
			out[name] = "refusal:" + string(got.Refusal.Kind)
		}
		return out
	}

	if a, b := derive(), derive(); !reflect.DeepEqual(a, b) {
		t.Errorf("the outcome set is not regenerable: %v vs %v", a, b)
	}
}

// REQ-80: "Mutation check — the check's non-vacuity. Two named mutants, run
// separately, with this RDR's **own** new tests (scenarios 1–4, 7, 9, 10, 10b)
// **excluded** from the oracle"
// ADVERSARIAL — TS scenario 6
//
// The mutation check itself is a recorded build procedure (see coverage.md's
// ASSUMPTIONS), not a permanent source-mutating test. What IS permanently
// assertable is the property the two mutants exist to establish: the check is
// non-vacuous — it separates a breaching table from a conforming one, so
// neither "always nil" nor "always error" satisfies it.
func TestReq80_TheCheckIsNonVacuous(t *testing.T) {
	if err := conformingNoMatchInput().Table.CheckValid(); err != nil {
		t.Errorf("CheckValid rejects a conforming table (%v); an "+
			"always-error check would be as vacuous as an always-nil one", err)
	}
	if err := breachingNoMatchInput().Table.CheckValid(); err == nil {
		t.Errorf("CheckValid accepts a breaching table; the check is " +
			"vacuous (Mutant A: return nil unconditionally)")
	}
}

// REQ-81: "*Mutant A* — `Table.CheckValid` returns `nil` unconditionally.
// **Expected**: the frozen suite still passes."
// ADVERSARIAL — TS scenario 6
//
// Mutant A's expectation is that the FROZEN suite is insensitive to the check
// — which is only meaningful if the frozen suite's own fixtures conform.
// Asserted here: every pre-existing disposition fixture passes CheckValid, so
// replacing the check with `return nil` cannot change any of their answers.
func TestReq81_TheFrozenFixturesAreInsensitiveToTheCheck(t *testing.T) {
	nonVacuityGate(t)

	for name, in := range allDispositionInputs() {
		if err := in.Table.CheckValid(); err != nil {
			t.Errorf("%s: a frozen-suite fixture breaches the check (%v); "+
				"Mutant A (return nil) would then change its answer and the "+
				"mutation check would not isolate the two mutants", name, err)
		}
	}
}

// REQ-82: "*Mutant B* — re-introduce the breach into the fixtures (restore
// `Writes` on `fixtures_test.go::escapeRow`) with the entry check intact.
// **Expected**: the suite fails, proving the check is wired into the real
// evaluation path and not dead code."
// ADVERSARIAL — TS scenario 6
//
// Mutant B's substance, assertable without mutating source: restoring the
// writes onto the shared builder's output produces a table that `Resolve`
// REJECTS. That is the proof the check is wired into the real path rather
// than being callable-but-dead.
func TestReq82_RestoringTheWritesMakesResolveFail(t *testing.T) {
	row := escapeRow(breachRuleID, breachLocator, resolve.KindNoMatch)
	row.Writes = breachWrites() // Mutant B, applied to a value not to source.

	in := noMatchInput()
	in.Table.Rows = append(in.Table.Rows, row)

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("Resolve accepted the Mutant-B table: %+v — the entry "+
			"check is not wired into the real evaluation path", got)
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("Resolve errored for some other reason: %v", err)
	}
}

// REQ-83: "The exported predicate called directly by a producer at
// construction time, on the same tables as scenarios 1–3. **Expected**:
// Verdicts identical to `Resolve`'s entry check — one predicate, two call
// sites, no drift. "Identical" is defined as: both nil, or both non-nil with
// equal extracted `[]RowRef` **and** equal per-identity `Count` values, in the
// same order. Not `reflect.DeepEqual` over the two error values"
// HAPPY PATH — TS scenario 7
//
// The three scenario tables, with the record's own definition of "identical"
// — extracted identities and counts in order, never DeepEqual over errors.
func TestReq83_TheExportedPredicateMatchesTheEntryCheckOnScenarios1To3(t *testing.T) {
	nonVacuityGate(t)

	scenarios := map[string]resolve.Input{
		"scenario_1_breach":        breachingNoMatchInput(),
		"scenario_2_dormant":       dormantBreachInput(),
		"scenario_3_empty_not_nil": emptyNotNilEscapeInput(),
	}

	for name, in := range scenarios {
		t.Run(name, func(t *testing.T) {
			direct := in.Table.CheckValid()
			_, entry := resolve.Resolve(in)

			if (direct == nil) != (entry == nil) {
				t.Fatalf("verdicts differ: CheckValid=%v Resolve=%v",
					direct, entry)
			}
			if direct == nil {
				return
			}
			if !sameRefs(breachRefs(t, direct), breachRefs(t, entry)) {
				t.Errorf("extracted []RowRef differ, in value or in order: "+
					"%+v vs %+v",
					breachRefs(t, direct), breachRefs(t, entry))
			}
			if !sameCounts(breachCounts(t, direct), breachCounts(t, entry)) {
				t.Errorf("per-identity Count values differ: %v vs %v",
					breachCounts(t, direct), breachCounts(t, entry))
			}
		})
	}
}

// REQ-84: "**Scenario**: The zero-`Result` caller trap on a breach.
// **Expected**: `Resolve` returns `Result{Plan: nil, Refusal: nil}` alongside
// the non-nil error, so `Refused() == false` on a call that did not succeed.
// Asserted directly"
// ADVERSARIAL — TS scenario 7b
func TestReq84_TheZeroResultCallerTrapIsAssertedDirectly(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("the breach did not error")
	}
	if got.Plan != nil {
		t.Errorf("Plan = %+v; want nil", got.Plan)
	}
	if got.Refusal != nil {
		t.Errorf("Refusal = %+v; want nil", got.Refusal)
	}
	if got.Refused() {
		t.Errorf("Refused() = true; the trap is that a call that did NOT " +
			"succeed reports false")
	}
}

// REQ-85: "A table with two or more breaching rows. **Expected**: One
// `Resolve` call reports **every** breaching identity, not just the first;
// each per-row error is individually recoverable by traversing the aggregate's
// `Unwrap() []error`"
// HAPPY PATH — TS scenario 9
func TestReq85_OneCallReportsEveryBreachingIdentity(t *testing.T) {
	want := []resolve.RowRef{
		{RuleID: "rule.one", SourceLocator: "flows/one.toml:1"},
		{RuleID: "rule.two", SourceLocator: "flows/two.toml:2"},
		{RuleID: "rule.three", SourceLocator: "flows/three.toml:3"},
	}
	err := resolveErr(t, multiBreachInput(want...))

	elems := breachElements(t, err)
	if len(elems) != len(want) {
		t.Fatalf("the aggregate carries %d per-row errors; want %d",
			len(elems), len(want))
	}
	// Each is INDIVIDUALLY recoverable — carrying its own identity, not a
	// shared or repeated one.
	seen := map[resolve.RowRef]bool{}
	for _, e := range elems {
		if seen[e.Ref] {
			t.Errorf("identity %+v reported twice", e.Ref)
		}
		seen[e.Ref] = true
	}
	for _, ref := range want {
		if !seen[ref] {
			t.Errorf("breaching identity %+v was not reported", ref)
		}
	}
}

// REQ-86: "The traversal is **exactly one level deep** and every element
// type-asserts to `*EscapeShapeBreachError`; `errors.Is(elem,
// ErrEscapeShapeBreach)` holds on **each element**, not only on the aggregate.
// The assertion compares the extracted `[]RowRef` — never `reflect.DeepEqual`
// over two `errors.Join` values"
// ADVERSARIAL — TS scenario 9
func TestReq86_EveryElementTypeAssertsAndIsClassified(t *testing.T) {
	err := resolveErr(t, multiBreachInput(
		resolve.RowRef{RuleID: "rule.one", SourceLocator: "flows/one.toml:1"},
		resolve.RowRef{RuleID: "rule.two", SourceLocator: "flows/two.toml:2"},
	))

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("the aggregate does not expose Unwrap() []error (%T)", err)
	}
	for i, elem := range joined.Unwrap() {
		// A DIRECT type assertion, one level deep — not errors.As, which
		// would walk deeper and hide a nested join.
		typed, ok := elem.(*resolve.EscapeShapeBreachError)
		if !ok {
			t.Fatalf("element %d is %T; every element must type-assert to "+
				"*EscapeShapeBreachError", i, elem)
		}
		if !errors.Is(typed, resolve.ErrEscapeShapeBreach) {
			t.Errorf("element %d is not classified by the sentinel; "+
				"errors.Is must hold on EACH element", i)
		}
	}
}

// REQ-87: "The same multi-breach table, rows carrying *distinct* `RowRef`
// identities, supplied in two different orders. **Expected**: Identical
// reported row sequence — ordered by `RowRef` identity (`compareRefs`), never
// by `Table.Rows` position."
// ADVERSARIAL — TS scenario 10
func TestReq87_TwoSupplyOrdersReportTheIdenticalSequence(t *testing.T) {
	a := resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"}
	b := resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"}
	c := resolve.RowRef{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"}

	forward := breachRefs(t, resolveErr(t, multiBreachInput(a, b, c)))
	reverse := breachRefs(t, resolveErr(t, multiBreachInput(c, b, a)))

	if !sameRefs(forward, reverse) {
		t.Fatalf("the two supply orders report different sequences: "+
			"%+v vs %+v", forward, reverse)
	}
	if !sameRefs(forward, []resolve.RowRef{a, b, c}) {
		t.Errorf("reported sequence = %+v; want compareRefs order [a b c]",
			forward)
	}
}

// REQ-88: "The table MUST be **mixed** — at least one row carrying source
// identity and at least one built without it (`RowRef{"",""}`) — so the
// report's ordering is pinned across the two identity classes rather than only
// within one."
// BOUNDARY — TS scenario 10
//
// The zero identity sorts first under compareRefs (empty strings compare
// least), so the mixed table pins the ordering ACROSS the classes.
func TestReq88_TheOrderingIsPinnedAcrossBothIdentityClasses(t *testing.T) {
	zero := resolve.RowRef{}
	named := resolve.RowRef{RuleID: "rule.named", SourceLocator: "flows/n.toml:5"}

	forward := breachRefs(t, resolveErr(t, multiBreachInput(named, zero)))
	reverse := breachRefs(t, resolveErr(t, multiBreachInput(zero, named)))

	if !sameRefs(forward, reverse) {
		t.Fatalf("the mixed table reports different sequences under the two "+
			"supply orders: %+v vs %+v", forward, reverse)
	}
	if !sameRefs(forward, []resolve.RowRef{zero, named}) {
		t.Errorf("reported sequence = %+v; want [RowRef{} named] — the "+
			"empty identity sorts least under compareRefs", forward)
	}
}

// REQ-89: "A multi-breach table whose breaching rows share one `RowRef`
// identity (including the zero-value `RowRef{"",""}` of rows built without
// source identity), supplied in two different orders. **Expected**: Identical
// report under both permutations — equal identities collapse to one reported
// entry. For three rows sharing `RowRef{"",""}`: `Unwrap() []error` has length
// **1**, its single element has `Ref == RowRef{"",""}` and `Count == 3`
// (pre-collapse rows, read off the struct field — no assertion parses message
// text)."
// BOUNDARY — TS scenario 10b
func TestReq89_SharedIdentitiesCollapseIdenticallyUnderBothOrders(t *testing.T) {
	zero := resolve.RowRef{}

	// Three rows sharing the zero identity, distinguished only by their
	// writes so the two permutations are genuinely different tables.
	build := func(order []string) resolve.Input {
		in := noMatchInput()
		for _, value := range order {
			row := breachingEscapeRow("", "", resolve.KindNoMatch)
			row.Writes = []resolve.Tag{{Key: "status", Value: value}}
			in.Table.Rows = append(in.Table.Rows, row)
		}
		return in
	}

	first := resolveErr(t, build([]string{"A", "B", "C"}))
	second := resolveErr(t, build([]string{"C", "A", "B"}))

	for label, err := range map[string]error{"first": first, "second": second} {
		elems := breachElements(t, err)
		if len(elems) != 1 {
			t.Fatalf("%s: Unwrap() []error has length %d; want 1 — equal "+
				"identities collapse to ONE entry", label, len(elems))
		}
		if elems[0].Ref != zero {
			t.Errorf("%s: Ref = %+v; want RowRef{\"\",\"\"}",
				label, elems[0].Ref)
		}
		if elems[0].Count != 3 {
			t.Errorf("%s: Count = %d; want 3 pre-collapse rows, read off "+
				"the struct field", label, elems[0].Count)
		}
	}

	if !sameRefs(breachRefs(t, first), breachRefs(t, second)) ||
		!sameCounts(breachCounts(t, first), breachCounts(t, second)) {
		t.Errorf("the two permutations report different collapsed entries")
	}
}

// REQ-92: "Done means, in user terms: **an escape row can no longer carry an
// owned-tag write to the accessor layer, so no tag value can appear in owned
// state that no authored rule set**"
// HAPPY PATH
//
// The user-facing outcome, asserted end to end at the kernel boundary: no
// reachable path produces a Plan that is BOTH escaped AND carrying writes. The
// plan is the only thing the accessor layer applies, so a plan with no writes
// is a guarantee no unauthored tag value reaches owned state.
func TestReq92_NoEscapedPlanEverCarriesAWrite(t *testing.T) {
	inputs := map[string]resolve.Input{
		"conformed_escape": conformingNoMatchInput(),
		"empty_not_nil":    emptyNotNilEscapeInput(),
		"breaching":        breachingNoMatchInput(),
		"dormant":          dormantBreachInput(),
	}
	for name, in := range allDispositionInputs() {
		inputs["frozen_"+name] = in
	}

	for name, in := range inputs {
		t.Run(name, func(t *testing.T) {
			got, err := resolve.Resolve(in)
			if err != nil {
				// A breaching table errors; no plan reaches the accessor
				// layer at all, which is the strongest form of the claim.
				if got.Plan != nil {
					t.Errorf("a plan accompanied an error: %+v", got.Plan)
				}
				return
			}
			if got.Plan == nil {
				return
			}
			if got.Plan.Escaped && len(got.Plan.Writes) != 0 {
				t.Errorf("an escaped plan carries %d writes (%+v); an escape "+
					"row can no longer carry an owned-tag write to the "+
					"accessor layer", len(got.Plan.Writes), got.Plan.Writes)
			}
		})
	}
}
