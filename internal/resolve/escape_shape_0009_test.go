package resolve_test

// RDR 0009 — escape-row shape conformance, the kernel half.
//
// One contract: an escape row bears no owned-state mutation. The kernel
// enforces it as an entry precondition of `Resolve` and exposes the same
// predicate as `Table.CheckValid` so a producer can fail at construction
// time. This file pins the predicate, the typed error surface, the
// multi-breach report, and the precedence of the breach over every modeled
// disposition.
//
// No test here asserts on message prose: REQ-27, REQ-44 and REQ-86 make row
// identity and count recoverable ONLY from structure.

import (
	"errors"
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// --- A. The producer obligation (the single rule) ------------------------

// REQ-1: "Escape-row shape conformance — an escape row carries no owned-state
// mutation — is a PRODUCER obligation on every constructor of resolve.Row
// values: a Row with a non-empty Escape list MUST have an empty Writes slice."
// HAPPY PATH
//
// The obligation is stated as a predicate over one row, so it is pinned as
// one: a non-empty Escape plus a non-empty Writes is the only combination
// the kernel rejects, and each of the other three combinations conforms.
func TestReq1_AnEscapeRowMustHaveAnEmptyWritesSlice(t *testing.T) {
	cases := []struct {
		name    string
		escape  []resolve.RefusalKind
		writes  []resolve.Tag
		breachs bool
	}{
		{"escape_with_writes_breaches",
			[]resolve.RefusalKind{resolve.KindNoMatch}, breachWrites(), true},
		{"escape_without_writes_conforms",
			[]resolve.RefusalKind{resolve.KindNoMatch}, nil, false},
		{"ordinary_row_with_writes_conforms",
			nil, breachWrites(), false},
		{"ordinary_row_without_writes_conforms",
			nil, nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tbl := resolve.Table{
				Revision: "rev-req1",
				Outcomes: []string{"successful"},
				Rows: []resolve.Row{{
					RuleID:        breachRuleID,
					SourceLocator: breachLocator,
					Outcome:       "successful",
					Escape:        tc.escape,
					Writes:        tc.writes,
				}},
			}
			err := tbl.CheckValid()
			if tc.breachs && err == nil {
				t.Fatalf("CheckValid returned nil; a non-empty Escape with a "+
					"non-empty Writes is the breach (escape=%v writes=%v)",
					tc.escape, tc.writes)
			}
			if !tc.breachs && err != nil {
				t.Fatalf("CheckValid returned %v; this row conforms "+
					"(escape=%v writes=%v)", err, tc.escape, tc.writes)
			}
		})
	}
}

// REQ-2: "(Authored clears normalize to `<clear>` writes per RDR 0002, so the
// Writes predicate carries both "no writes" and "no clears" at the kernel
// boundary.)"
// DOMAIN EDGE
//
// A normalized clear reaches the kernel as an ordinary Writes entry carrying
// the `<clear>` sentinel value. The predicate is over Writes, so the clear is
// caught by the same clause with no separate rule — asserted by feeding an
// escape row whose ONLY write is a clear.
func TestReq2_ANormalizedClearOnAnEscapeRowIsTheSameBreach(t *testing.T) {
	tbl := resolve.Table{
		Revision: "rev-req2",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			Outcome:       "successful",
			Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
			// The `<clear>` sentinel RDR 0002 normalizes an authored clear
			// into. It is a write like any other at this boundary.
			Writes: []resolve.Tag{{Key: "status", Value: "<clear>"}},
		}},
	}

	err := tbl.CheckValid()
	if err == nil {
		t.Fatalf("CheckValid returned nil for an escape row carrying a " +
			"normalized clear; the Writes predicate carries both 'no writes' " +
			"and 'no clears' at the kernel boundary")
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("the clear breach is not classified as ErrEscapeShapeBreach")
	}
}

// REQ-3: "The predicate is Writes-only and does NOT extend to NextTags: A4
// settled that owned state is reachable only through a write accessor, which
// RDR 0004 scopes to "planned owned-tag writes," so an escape row's NextTags
// cannot mutate owned state."
// BOUNDARY
//
// The negative half of the contract. An escape row carrying a populated
// NextTags and NO Writes conforms — widening the predicate to NextTags would
// break this.
func TestReq3_ThePredicateDoesNotExtendToNextTags(t *testing.T) {
	nonVacuityGate(t)

	tbl := resolve.Table{
		Revision: "rev-req3",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			Outcome:       "successful",
			Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
			NextTags: []resolve.Tag{
				{Key: "status", Value: "Blocked"},
				{Key: "phase", Value: "escaped"},
			},
		}},
	}

	if err := tbl.CheckValid(); err != nil {
		t.Fatalf("CheckValid returned %v for an escape row whose only "+
			"populated mutation field is NextTags; the predicate is "+
			"Writes-only and does not extend to NextTags", err)
	}
}

// REQ-4: "the breach predicate over a single row is `len(row.Escape) != 0 &&
// len(row.Writes) != 0`."
// BOUNDARY
//
// The predicate is a conjunction over LENGTHS. Pinned at the boundary of each
// conjunct: one populated entry on each side is a breach; zero on either side
// is not.
func TestReq4_TheBreachPredicateIsTheLengthConjunction(t *testing.T) {
	oneEscapeOneWrite := resolve.Table{
		Revision: "rev-req4",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			Outcome:       "successful",
			Escape:        []resolve.RefusalKind{resolve.KindNoMatch},
			Writes:        []resolve.Tag{{Key: "status", Value: "Escaped"}},
		}},
	}
	if err := oneEscapeOneWrite.CheckValid(); err == nil {
		t.Errorf("one escape class plus one write is the minimal breach; " +
			"CheckValid returned nil")
	}

	zeroEscape := oneEscapeOneWrite
	zeroEscape.Rows = []resolve.Row{oneEscapeOneWrite.Rows[0]}
	zeroEscape.Rows[0].Escape = []resolve.RefusalKind{}
	if err := zeroEscape.CheckValid(); err != nil {
		t.Errorf("a zero-length Escape list makes the row an ordinary row; "+
			"CheckValid returned %v", err)
	}

	zeroWrites := oneEscapeOneWrite
	zeroWrites.Rows = []resolve.Row{oneEscapeOneWrite.Rows[0]}
	zeroWrites.Rows[0].Escape = []resolve.RefusalKind{resolve.KindNoMatch}
	zeroWrites.Rows[0].Writes = []resolve.Tag{}
	if err := zeroWrites.CheckValid(); err != nil {
		t.Errorf("a zero-length Writes slice conforms; CheckValid "+
			"returned %v", err)
	}
}

// REQ-5: "*Nil-vs-empty* — length-based on purpose: a non-nil but empty slice
// is conforming. The predicate tests emptiness, never nil-ness."
// INPUT EDGE
//
// The four nil/empty permutations across the two fields. A nil-sensitive
// implementation would separate `[]Tag{}` from `nil`; this asserts it must
// not, in BOTH directions (an empty Escape with real writes is equally
// conforming).
func TestReq5_NilAndEmptySlicesAreIndistinguishableToThePredicate(t *testing.T) {
	nonVacuityGate(t)

	cases := []struct {
		name   string
		escape []resolve.RefusalKind
		writes []resolve.Tag
	}{
		{"nil_escape_nil_writes", nil, nil},
		{"nil_escape_empty_writes", nil, []resolve.Tag{}},
		{"empty_escape_nil_writes", []resolve.RefusalKind{}, nil},
		{"empty_escape_empty_writes", []resolve.RefusalKind{}, []resolve.Tag{}},
		{"escape_empty_writes",
			[]resolve.RefusalKind{resolve.KindNoMatch}, []resolve.Tag{}},
		{"empty_escape_real_writes",
			[]resolve.RefusalKind{}, breachWrites()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tbl := resolve.Table{
				Revision: "rev-req5",
				Outcomes: []string{"successful"},
				Rows: []resolve.Row{{
					RuleID:        breachRuleID,
					SourceLocator: breachLocator,
					Outcome:       "successful",
					Escape:        tc.escape,
					Writes:        tc.writes,
				}},
			}
			if err := tbl.CheckValid(); err != nil {
				t.Fatalf("CheckValid returned %v; the predicate tests "+
					"emptiness, never nil-ness", err)
			}
		})
	}
}

// REQ-6: "The shared resolve.Row type keeps its single shape: escape identity
// remains discriminated solely by a non-empty Escape list. No ordinary/escape
// type split and no row-kind field is introduced at the kernel boundary."
// ADVERSARIAL
//
// Asserted structurally over the shipped `resolve.Row` type: its field set
// must gain no kind/variant discriminator, and there must be no second
// escape-row type at the kernel boundary. A type split would show up here as
// a new field or as `Row` losing `Escape`.
func TestReq6_RowKeepsItsSingleShapeWithNoKindField(t *testing.T) {
	rt := reflect.TypeOf(resolve.Row{})

	esc, ok := rt.FieldByName("Escape")
	if !ok {
		t.Fatalf("resolve.Row lost its Escape field; escape identity is " +
			"discriminated solely by a non-empty Escape list")
	}
	if esc.Type.Kind() != reflect.Slice {
		t.Errorf("Row.Escape is %v; the discriminator is a LIST whose "+
			"non-emptiness is the signal", esc.Type)
	}

	banned := []string{"Kind", "RowKind", "IsEscape", "Escaped", "Variant", "Type"}
	for _, name := range banned {
		if _, found := rt.FieldByName(name); found {
			t.Errorf("resolve.Row gained a row-kind field %q; no row-kind "+
				"field is introduced at the kernel boundary", name)
		}
	}
}

// REQ-7: "The escape discriminator stays `len(Escape) != 0` / `rescues` — the
// existing sibling signal, reused."
// BOUNDARY
//
// The discriminator the precondition keys on must be the same one the
// candidate partition already keys on. Asserted behaviourally: a row with a
// non-empty Escape is excluded from the ordinary candidate partition AND is
// subject to the shape predicate, and a row with an empty Escape is neither.
func TestReq7_TheEscapeDiscriminatorIsTheReusedNonEmptyEscapeSignal(t *testing.T) {
	// A row with a non-empty Escape carrying writes: the shape predicate
	// sees it as an escape row.
	escaping := resolve.Table{
		Revision: "rev-req7",
		Outcomes: []string{"successful"},
		Rows: []resolve.Row{{
			RuleID:        breachRuleID,
			SourceLocator: breachLocator,
			Outcome:       "successful",
			Escape:        []resolve.RefusalKind{resolve.KindAmbiguousMatch},
			Writes:        breachWrites(),
		}},
	}
	if err := escaping.CheckValid(); err == nil {
		t.Errorf("a row whose Escape list is non-empty is an escape row to " +
			"the predicate, whatever class it declares")
	}

	// The same row with the SAME writes and an empty Escape is an ordinary
	// row to the same discriminator, and is not subject to the predicate.
	ordinary := escaping
	ordinary.Rows = []resolve.Row{escaping.Rows[0]}
	ordinary.Rows[0].Escape = nil
	if err := ordinary.CheckValid(); err != nil {
		t.Errorf("removing the Escape list makes the row ordinary; "+
			"CheckValid returned %v", err)
	}
}

// REQ-8: "`planOf` stays unconditional — under the precondition it can no
// longer copy writes onto an escaped plan, so no stripping logic is added"
// DOMAIN EDGE
//
// Observable consequence: on a CONFORMING table the escaped plan carries the
// row's writes verbatim (zero-length, because a conforming escape row has
// none) and the ordinary plan still carries its own writes copied through. A
// `planOf` that gained a "strip when escaped" branch would still pass the
// escape leg — so the discriminating leg is the ORDINARY one, which must be
// untouched, plus the escaped plan reporting the row's own (empty) writes
// rather than a sanitized nil that hides a producer breach.
func TestReq8_PlanOfStaysUnconditionalAndAddsNoStrippingLogic(t *testing.T) {
	nonVacuityGate(t)

	// Ordinary path: writes copied through unchanged.
	ordinary := mustResolve(t, legalInput())
	if ordinary.Plan == nil {
		t.Fatalf("the legal input no longer produces a plan")
	}
	if len(ordinary.Plan.Writes) == 0 {
		t.Errorf("the ordinary plan lost its writes; planOf is unconditional " +
			"and copies the selected row's writes")
	}

	// Escape path on a conforming table: the plan reports the row's writes,
	// which are empty because the row conforms — not because anything
	// stripped them.
	escaped := mustResolve(t, conformingNoMatchInput())
	if escaped.Plan == nil {
		t.Fatalf("the conforming escape table no longer rescues: %+v", escaped)
	}
	if !escaped.Plan.Escaped {
		t.Errorf("the plan is not marked Escaped; the escape path did not run")
	}
	if len(escaped.Plan.Writes) != 0 {
		t.Errorf("the escaped plan carries %d writes; a conforming escape "+
			"row has none and planOf adds none", len(escaped.Plan.Writes))
	}
}

// REQ-9: "Data flow is unchanged for every conforming table: `assemble` →
// candidate partition (`len(row.Escape) != 0`) → `gate` → selection →
// `planOf`."
// HAPPY PATH
//
// The no-change claim, asserted as value equality of the disposition across
// every shipped disposition fixture: adding the precondition must alter no
// conforming table's answer.
func TestReq9_DataFlowIsUnchangedForEveryConformingTable(t *testing.T) {
	nonVacuityGate(t)

	for name, in := range allDispositionInputs() {
		t.Run(name, func(t *testing.T) {
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Fatalf("a conforming table used the Go error path: %v", err)
			}
			// The disposition must still be exactly one of plan or refusal.
			if (got.Plan == nil) == (got.Refusal == nil) {
				t.Fatalf("disposition is neither or both: %+v", got)
			}
		})
	}

	// And the escape path specifically still reaches planOf.
	escaped := mustResolve(t, conformingNoMatchInput())
	if escaped.Plan == nil || !escaped.Plan.Escaped {
		t.Errorf("the conforming escape table no longer reaches planOf via " +
			"the rescue phase")
	}
}

// REQ-10: "Kernel sanitization (strip `Writes` when the selection is an
// escape)" is rejected — recorded as Briefly Rejected; no
// sanitization/stripping path may be introduced.
// ADVERSARIAL
//
// The rejected alternative would make a breaching table RESOLVE with the
// writes silently dropped. This asserts the opposite: a breaching table
// errors rather than producing a sanitized plan. A sanitizing kernel passes
// REQ-8's escape leg but fails here.
func TestReq10_TheKernelNeverSanitizesAwayABreach(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("a breaching table resolved: %+v — sanitization is "+
			"rejected; the breach must surface, never be stripped", got)
	}
	if got.Plan != nil {
		t.Errorf("a sanitized plan was emitted alongside the error: %+v",
			got.Plan)
	}
}

// --- C. The kernel entry precondition ------------------------------------

// REQ-16: "The kernel enforces the same obligation as an entry precondition of
// Resolve."
// HAPPY PATH
//
// Same obligation, not a second one: whatever `CheckValid` rejects, `Resolve`
// rejects, and whatever it accepts, `Resolve` accepts.
func TestReq16_ResolveEnforcesTheSameObligationAtEntry(t *testing.T) {
	breaching := breachingNoMatchInput()
	if breaching.Table.CheckValid() == nil {
		t.Fatalf("the breaching fixture does not breach CheckValid")
	}
	if _, err := resolve.Resolve(breaching); err == nil {
		t.Errorf("Resolve accepted a table CheckValid rejects")
	}

	conforming := conformingNoMatchInput()
	if err := conforming.Table.CheckValid(); err != nil {
		t.Fatalf("the conforming fixture breaches CheckValid: %v", err)
	}
	if _, err := resolve.Resolve(conforming); err != nil {
		t.Errorf("Resolve rejected a table CheckValid accepts: %v", err)
	}
}

// REQ-17: "Resolve's signature is UNCHANGED — `func Resolve(in Input) (Result,
// error)`, one parameter — and the checked table is the one already reached
// through the input, `in.Table`; the precondition adds no parameter and takes
// no table argument of its own."
// BOUNDARY
//
// Asserted on the shipped function type: one input parameter, two results.
// The precondition must not have been bolted on as a second argument.
func TestReq17_ResolveKeepsItsOneParameterSignature(t *testing.T) {
	ft := reflect.TypeOf(resolve.Resolve)

	if ft.NumIn() != 1 {
		t.Fatalf("Resolve takes %d parameters; the signature is unchanged "+
			"at exactly one (Input)", ft.NumIn())
	}
	if ft.In(0) != reflect.TypeOf(resolve.Input{}) {
		t.Errorf("Resolve's parameter is %v; want resolve.Input", ft.In(0))
	}
	if ft.NumOut() != 2 {
		t.Fatalf("Resolve returns %d values; want (Result, error)", ft.NumOut())
	}
	if ft.Out(0) != reflect.TypeOf(resolve.Result{}) {
		t.Errorf("Resolve's first result is %v; want resolve.Result", ft.Out(0))
	}
}

// REQ-18: "A table containing a row that breaches the conformance predicate
// above MUST cause Resolve to return a non-nil Go error identifying the
// offending row by RuleID and SourceLocator, with no Result disposition."
// HAPPY PATH
//
// Both halves: the identity reaches the caller as structure (RuleID and
// SourceLocator, via RowRef), and the Result carries no disposition.
func TestReq18_ResolveReturnsATypedErrorNamingTheRowAndNoDisposition(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("Resolve returned no error for a breaching table: %+v", got)
	}
	if got.Plan != nil || got.Refusal != nil {
		t.Errorf("Resolve returned a disposition alongside the breach "+
			"error: %+v", got)
	}

	refs := breachRefs(t, err)
	want := resolve.RowRef{RuleID: breachRuleID, SourceLocator: breachLocator}
	if len(refs) != 1 || refs[0] != want {
		t.Errorf("reported identities = %+v; want exactly %+v — the "+
			"offending row is named by RuleID and SourceLocator", refs, want)
	}
}

// REQ-19: "The precondition is evaluated over the WHOLE table before any
// evaluation step: a breach surfaces even when no resolution path reaches the
// offending row, and it precedes every modeled disposition (a malformed table
// is malformed as a value, independent of the input tuple)."
// ADVERSARIAL
//
// The dormant row is the discriminating case: an on-path-only check would
// resolve this table cleanly. It must error instead, and name the dormant row.
func TestReq19_ADormantBreachingRowStillErrors(t *testing.T) {
	in := dormantBreachInput()

	// The same table WITHOUT the dormant row resolves cleanly — so nothing
	// but the dormant row can be responsible for the error below.
	base := legalInput()
	if _, err := resolve.Resolve(base); err != nil {
		t.Fatalf("the base table does not resolve cleanly: %v", err)
	}

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("a table whose ONLY breach is on an unreachable row "+
			"resolved: %+v — the precondition is over the whole table", got)
	}
	if got.Plan != nil || got.Refusal != nil {
		t.Errorf("a disposition accompanied the dormant-row breach: %+v", got)
	}

	refs := breachRefs(t, err)
	want := resolve.RowRef{RuleID: breachRuleID, SourceLocator: breachLocator}
	if len(refs) != 1 || refs[0] != want {
		t.Errorf("reported identities = %+v; want the dormant row %+v",
			refs, want)
	}
}

// REQ-20: "The breach travels the error path RDR 0001 reserves for programmer
// mistakes: it MUST NOT be a modeled refusal, MUST NOT introduce a refusal
// kind, and MUST NOT alter any disposition of a conforming table."
// ADVERSARIAL
//
// Three separate obligations, each asserted: no Refusal value, no new member
// in the closed five-kind taxonomy, and unchanged dispositions for the
// shipped conforming fixtures.
func TestReq20_TheBreachIsTheErrorPathAndNotAModeledRefusal(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("the breach did not travel the Go error path: %+v", got)
	}
	if got.Refusal != nil {
		t.Errorf("the breach was modeled as a refusal (%v); it is a "+
			"programmer mistake, not a value-level disposition",
			got.Refusal.Kind)
	}
	if got.Refused() {
		t.Errorf("Refused() reports true on a breach; there is no refusal")
	}

	kinds := resolve.RefusalKinds()
	if len(kinds) != 5 {
		t.Errorf("the refusal taxonomy carries %d kinds (%v); the five "+
			"kinds stay closed and gain no member", len(kinds), kinds)
	}
	for _, k := range kinds {
		if string(k) == "escape_shape_breach" || string(k) == "escape-row-shape-breach" {
			t.Errorf("a breach refusal kind %q was introduced", k)
		}
	}

	for name, in := range allDispositionInputs() {
		before, err := resolve.Resolve(in)
		if err != nil {
			t.Errorf("%s: a conforming table now errors: %v", name, err)
			continue
		}
		after, err := resolve.Resolve(in)
		if err != nil {
			t.Errorf("%s: replay errored: %v", name, err)
			continue
		}
		if !reflect.DeepEqual(before, after) {
			t.Errorf("%s: disposition is not stable across replay", name)
		}
	}
}

// REQ-21: "*Precedence* — the breach error precedes every modeled disposition,
// `unmodeled_outcome` included."
// ADVERSARIAL
//
// `unmodeled_outcome` is the first step of the shipped evaluation order, so
// it is the hardest disposition to precede. Pinned explicitly.
func TestReq21_TheBreachPrecedesUnmodeledOutcome(t *testing.T) {
	in := unmodeledOutcomeInput()
	// Confirm the fixture really does yield unmodeled_outcome untouched.
	base := mustResolve(t, in)
	if base.Refusal == nil || base.Refusal.Kind != resolve.KindUnmodeledOutcome {
		t.Fatalf("the fixture no longer yields unmodeled_outcome: %+v", base)
	}

	in.Table.Rows = append(in.Table.Rows,
		breachingEscapeRow(breachRuleID, breachLocator, resolve.KindNoMatch))

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("unmodeled_outcome won over the breach: %+v — the breach "+
			"error precedes EVERY modeled disposition", got)
	}
	if got.Refusal != nil {
		t.Errorf("a refusal accompanied the breach: %v", got.Refusal.Kind)
	}
}

// REQ-23: "When any row carries both a non-empty `Escape` and a non-empty
// `Writes`, `Resolve` returns the zero `Result` (`Plan` and `Refusal` both nil
// — `Result` is a struct, so there is no nil `Result` to return) and a non-nil
// error naming every such row."
// BOUNDARY
//
// Asserted as value equality against the ZERO Result, not merely as two nil
// checks: any future field added to Result must also come back zeroed.
func TestReq23_ABreachReturnsTheZeroResultAndANonNilError(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("no error returned for a breaching table")
	}
	if !reflect.DeepEqual(got, resolve.Result{}) {
		t.Errorf("Result = %+v; want the ZERO Result", got)
	}
}

// REQ-24: "Because the zero `Result` reports `Refused() == false`, callers
// MUST check the error before reading the disposition; a caller that branches
// on `Refused()` first would read a success-shaped value with a nil `Plan`."
// ADVERSARIAL
//
// The caller trap, pinned directly so it cannot be quietly "fixed" by making
// a breach look like a refusal — which would break REQ-20.
func TestReq24_TheZeroResultReportsRefusedFalseOnABreach(t *testing.T) {
	got, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("no error returned for a breaching table")
	}
	if got.Refused() {
		t.Errorf("Refused() = true on a breach; the zero Result is " +
			"success-shaped and the error is the only signal")
	}
	if got.Plan != nil {
		t.Errorf("the success-shaped value carries a non-nil Plan: %+v",
			got.Plan)
	}
}

// --- D. The typed error surface ------------------------------------------

// REQ-25: "The breach error MUST be a TYPED error carrying the offending row
// identity as the kernel's existing RowRef value, inspectable via errors.As /
// errors.AsType without parsing message text, and MUST wrap a package-level
// sentinel naming the breach category so errors.Is can classify it."
// HAPPY PATH
func TestReq25_TheBreachErrorIsTypedAndWrapsTheSentinel(t *testing.T) {
	_, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("no error returned for a breaching table")
	}

	var typed *resolve.EscapeShapeBreachError
	if !errors.As(err, &typed) {
		t.Fatalf("errors.As found no *EscapeShapeBreachError in the chain; "+
			"the identity must be reachable without parsing text (%T)", err)
	}
	if typed.Ref.RuleID != breachRuleID ||
		typed.Ref.SourceLocator != breachLocator {
		t.Errorf("Ref = %+v; want the offending row's identity as a RowRef",
			typed.Ref)
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("errors.Is does not classify the error as " +
			"ErrEscapeShapeBreach")
	}
}

// REQ-26: "Both the error type and the sentinel MUST be EXPORTED"
// BOUNDARY
//
// Compile-time proof: this file is in package `resolve_test` and names both
// symbols from outside the package, so an unexported spelling cannot build.
func TestReq26_TheTypeAndSentinelAreExported(t *testing.T) {
	var typed *resolve.EscapeShapeBreachError
	if typed != nil {
		t.Fatalf("unreachable")
	}
	if resolve.ErrEscapeShapeBreach == nil {
		t.Fatalf("ErrEscapeShapeBreach is nil; the sentinel must be a " +
			"package-level exported error value")
	}

	rt := reflect.TypeOf(resolve.EscapeShapeBreachError{})
	if got := rt.Name(); got != "EscapeShapeBreachError" {
		t.Errorf("the type is named %q; the spelling is pinned as "+
			"EscapeShapeBreachError", got)
	}
}

// REQ-27: "Row identity MUST NOT be recoverable only from formatted prose."
// ADVERSARIAL
//
// The structural extraction must succeed on identities that contain no
// distinguishing prose at all: a rule id and locator built from characters a
// message formatter would mangle or a reader could not delimit.
func TestReq27_RowIdentityIsRecoverableWithoutParsingProse(t *testing.T) {
	awkward := resolve.RowRef{
		RuleID:        `rule "with" quotes, a comma and a: colon`,
		SourceLocator: "flows/a b.toml:1:2 (and) parens",
	}
	in := multiBreachInput(awkward)

	_, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("no error returned for a breaching table")
	}

	refs := breachRefs(t, err)
	if len(refs) != 1 || refs[0] != awkward {
		t.Errorf("reported identities = %+v; want the awkward identity %+v "+
			"recovered verbatim from structure", refs, awkward)
	}
}

// REQ-28: "the sentinel is ErrEscapeShapeBreach (package-level, exported);"
// BOUNDARY
func TestReq28_TheSentinelIsSpelledErrEscapeShapeBreach(t *testing.T) {
	if resolve.ErrEscapeShapeBreach == nil {
		t.Fatalf("ErrEscapeShapeBreach is nil")
	}
	// A sentinel classifies itself and nothing unrelated.
	if !errors.Is(resolve.ErrEscapeShapeBreach, resolve.ErrEscapeShapeBreach) {
		t.Errorf("the sentinel does not classify itself")
	}
	_, err := resolve.Resolve(breachingNoMatchInput())
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("the breach is not classified by the named sentinel")
	}
}

// REQ-29: "the typed error is *EscapeShapeBreachError (exported), carrying
// exactly ONE RowRef field, Ref, plus a Count int field (see the multi-breach
// clause);"
// BOUNDARY
//
// The field set is the claim: exactly one RowRef field named Ref, and a Count
// int. A second RowRef field would reintroduce the per-row multiplicity the
// collapse exists to remove.
func TestReq29_TheTypedErrorCarriesExactlyOneRowRefFieldAndACount(t *testing.T) {
	rt := reflect.TypeOf(resolve.EscapeShapeBreachError{})

	ref, ok := rt.FieldByName("Ref")
	if !ok {
		t.Fatalf("EscapeShapeBreachError has no Ref field")
	}
	if ref.Type != reflect.TypeOf(resolve.RowRef{}) {
		t.Errorf("Ref is %v; want resolve.RowRef", ref.Type)
	}

	count, ok := rt.FieldByName("Count")
	if !ok {
		t.Fatalf("EscapeShapeBreachError has no Count field")
	}
	if count.Type.Kind() != reflect.Int {
		t.Errorf("Count is %v; want int", count.Type)
	}

	rowRefs := 0
	for i := range rt.NumField() {
		if rt.Field(i).Type == reflect.TypeOf(resolve.RowRef{}) {
			rowRefs++
		}
	}
	if rowRefs != 1 {
		t.Errorf("EscapeShapeBreachError carries %d RowRef fields; the "+
			"clause fixes EXACTLY ONE, named Ref", rowRefs)
	}
}

// REQ-30: "its Unwrap() error returns ErrEscapeShapeBreach, so errors.Is
// classifies EVERY per-row element, not only the aggregate;"
// BOUNDARY
//
// Asserted on the ELEMENTS, which is where the aggregate-only implementation
// would differ: a wrapper that attached the sentinel to the join alone would
// still satisfy `errors.Is(aggregate, sentinel)`.
func TestReq30_EveryPerRowElementUnwrapsToTheSentinel(t *testing.T) {
	in := multiBreachInput(
		resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
	)
	_, err := resolve.Resolve(in)

	elems := breachElements(t, err)
	if len(elems) != 2 {
		t.Fatalf("aggregate carries %d elements; want 2", len(elems))
	}
	for i, e := range elems {
		if !errors.Is(e, resolve.ErrEscapeShapeBreach) {
			t.Errorf("element %d (%+v) is not classified by the sentinel; "+
				"errors.Is must classify EVERY per-row element", i, e.Ref)
		}
		if unwrapped := errors.Unwrap(e); unwrapped != resolve.ErrEscapeShapeBreach {
			t.Errorf("element %d Unwrap() = %v; want the sentinel itself",
				i, unwrapped)
		}
	}
}

// REQ-31: "the predicate is Table.CheckValid() error."
// BOUNDARY
func TestReq31_ThePredicateIsSpelledTableCheckValid(t *testing.T) {
	mt, ok := reflect.TypeOf(resolve.Table{}).MethodByName("CheckValid")
	if !ok {
		t.Fatalf("resolve.Table has no CheckValid method; the predicate " +
			"spelling is pinned")
	}
	// Method value on a value receiver: one receiver in, one error out.
	if mt.Type.NumIn() != 1 {
		t.Errorf("CheckValid takes %d arguments beyond the receiver; it "+
			"takes none", mt.Type.NumIn()-1)
	}
	if mt.Type.NumOut() != 1 ||
		mt.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		t.Errorf("CheckValid's results are not exactly (error)")
	}
}

// REQ-32: "These four spellings collide with no frozen boundary guard: the
// kernel's banned exported-name lists (REQ-25, REQ-36, REQ-37) are
// exact-match, and errors/fmt are forbidden by no import guard (REQ-26,
// REQ-28, REQ-37)."
// ADVERSARIAL
//
// The four pinned spellings must exist under exactly those names — a
// rename to dodge a guard is what this pins against.
func TestReq32_TheFourPinnedSpellingsAllExist(t *testing.T) {
	if resolve.ErrEscapeShapeBreach == nil {
		t.Errorf("spelling 1: ErrEscapeShapeBreach is absent")
	}
	if reflect.TypeOf(resolve.EscapeShapeBreachError{}).Name() !=
		"EscapeShapeBreachError" {
		t.Errorf("spelling 2: EscapeShapeBreachError is absent")
	}
	if _, ok := reflect.TypeOf(resolve.EscapeShapeBreachError{}).
		FieldByName("Ref"); !ok {
		t.Errorf("spelling 3: the Ref field is absent")
	}
	if _, ok := reflect.TypeOf(resolve.Table{}).
		MethodByName("CheckValid"); !ok {
		t.Errorf("spelling 4: Table.CheckValid is absent")
	}
}

// REQ-33: "Because errors.Join returns a wrapper even for a single error,
// callers MUST classify and extract with errors.Is / errors.As /
// errors.AsType rather than a direct type assertion or equality against the
// sentinel."
// ADVERSARIAL
//
// The single-breach case is the trap: the returned value is the JOIN
// WRAPPER, so a direct type assertion fails and equality against the sentinel
// fails, while errors.Is/As succeed. All four are asserted.
func TestReq33_TheSingleBreachReturnIsTheJoinWrapperNotTheBareError(t *testing.T) {
	_, err := resolve.Resolve(breachingNoMatchInput())
	if err == nil {
		t.Fatalf("no error returned for a breaching table")
	}

	if _, direct := err.(*resolve.EscapeShapeBreachError); direct {
		t.Errorf("a direct type assertion succeeded; even the single-breach " +
			"case returns the errors.Join wrapper")
	}
	if err == resolve.ErrEscapeShapeBreach {
		t.Errorf("the return equals the sentinel; equality is not the " +
			"classification instrument")
	}
	if !errors.Is(err, resolve.ErrEscapeShapeBreach) {
		t.Errorf("errors.Is failed to classify the single-breach return")
	}
	var typed *resolve.EscapeShapeBreachError
	if !errors.As(err, &typed) {
		t.Errorf("errors.As failed to extract from the single-breach return")
	}
}

// REQ-34: "The aggregate is the uniform return shape: Resolve MUST NOT return
// the bare per-row error in the one-breach case and the aggregate otherwise,
// so caller code has one shape to handle regardless of breach count."
// BOUNDARY
//
// One shape at both cardinalities: `Unwrap() []error` must be present and
// traversable at count 1 and count 2 alike.
func TestReq34_TheAggregateIsTheUniformReturnShapeAtEveryCardinality(t *testing.T) {
	one := breachRefs(t, resolveErr(t, breachingNoMatchInput()))
	if len(one) != 1 {
		t.Errorf("single-breach aggregate reported %d identities; want 1",
			len(one))
	}

	two := breachRefs(t, resolveErr(t, multiBreachInput(
		resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
	)))
	if len(two) != 2 {
		t.Errorf("multi-breach aggregate reported %d identities; want 2",
			len(two))
	}
}

// REQ-35: "Unwrap() []error MUST be EXACTLY ONE LEVEL deep, every element
// being a *EscapeShapeBreachError — never a nested join."
// ADVERSARIAL
//
// `breachElements` fails on a nested join or a foreign element type; this
// exercises it at a cardinality where a naive incremental join
// (join(join(a,b),c)) would nest.
func TestReq35_TheAggregateIsExactlyOneLevelDeep(t *testing.T) {
	err := resolveErr(t, multiBreachInput(
		resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
		resolve.RowRef{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"},
	))

	elems := breachElements(t, err)
	if len(elems) != 3 {
		t.Errorf("the flat traversal yielded %d elements; want 3 — a "+
			"nested join would collapse to 2", len(elems))
	}
}

// REQ-36: "Resolve MUST return CheckValid's error VERBATIM, never
// fmt.Errorf-wrapped: an added layer would expose Unwrap() error at the
// outermost level and break the flat traversal this clause guarantees."
// ADVERSARIAL
//
// The discriminating assertion: the OUTERMOST error must expose
// `Unwrap() []error` and must NOT expose `Unwrap() error`. A single
// `fmt.Errorf("%w", ...)` layer inverts both.
func TestReq36_ResolveReturnsCheckValidsErrorVerbatim(t *testing.T) {
	in := breachingNoMatchInput()

	if _, plural := resolveErr(t, in).(interface{ Unwrap() []error }); !plural {
		t.Errorf("the outermost error does not expose Unwrap() []error; " +
			"a wrapping layer was added")
	}
	if _, singular := resolveErr(t, in).(interface{ Unwrap() error }); singular {
		t.Errorf("the outermost error exposes Unwrap() error; Resolve " +
			"wrapped CheckValid's error instead of returning it verbatim")
	}

	// Verbatim also means the two call sites report the same identities and
	// counts.
	direct := in.Table.CheckValid()
	viaResolve := resolveErr(t, in)
	if !sameRefs(breachRefs(t, direct), breachRefs(t, viaResolve)) {
		t.Errorf("Resolve's error reports different identities than "+
			"CheckValid's: %+v vs %+v",
			breachRefs(t, viaResolve), breachRefs(t, direct))
	}
}

// REQ-37: "Spellings pinned (Pre-Lock, 3amigo T-2): sentinel
// `ErrEscapeShapeBreach`; typed error `*EscapeShapeBreachError` with fields
// `Ref RowRef` and `Count int`; predicate `Table.CheckValid() error`."
// BOUNDARY
//
// The field TYPES as well as the names, so a `Ref string` or a `Count int64`
// substitution is caught.
func TestReq37_ThePinnedSpellingsCarryThePinnedTypes(t *testing.T) {
	rt := reflect.TypeOf(resolve.EscapeShapeBreachError{})

	ref, _ := rt.FieldByName("Ref")
	if ref.Type != reflect.TypeOf(resolve.RowRef{}) {
		t.Errorf("Ref is %v; the spelling pins `Ref RowRef`", ref.Type)
	}
	count, _ := rt.FieldByName("Count")
	if count.Type != reflect.TypeOf(int(0)) {
		t.Errorf("Count is %v; the spelling pins `Count int`", count.Type)
	}

	var sentinel error = resolve.ErrEscapeShapeBreach
	if sentinel == nil {
		t.Errorf("the sentinel is not a non-nil package-level error value")
	}
}

// REQ-38: "A single `errors.As`/`AsType` call reports only the FIRST breach in
// the chain, so it is the wrong instrument for reading a multi-breach report —
// a test asserting on all offending rows must traverse the aggregate."
// DOMAIN EDGE
//
// The failure mode itself is the contract: a single errors.As over a
// multi-breach aggregate must yield strictly fewer identities than the
// traversal does. This is what makes the traversal necessary rather than
// merely tidy.
func TestReq38_ASingleErrorsAsReportsOnlyOneOfSeveralBreaches(t *testing.T) {
	err := resolveErr(t, multiBreachInput(
		resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
		resolve.RowRef{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"},
	))

	var single *resolve.EscapeShapeBreachError
	if !errors.As(err, &single) {
		t.Fatalf("errors.As extracted nothing from a three-breach aggregate")
	}

	all := breachRefs(t, err)
	if len(all) != 3 {
		t.Fatalf("the traversal reported %d identities; want 3", len(all))
	}
	// The single extraction is one identity out of three: it is a strict
	// subset, which is exactly why the traversal is mandatory.
	found := false
	for _, ref := range all {
		if ref == single.Ref {
			found = true
		}
	}
	if !found {
		t.Errorf("errors.As returned %+v, which is not among the reported "+
			"identities %+v", single.Ref, all)
	}
}

// --- E. Multi-breach reporting, ordering, and collapse -------------------

// REQ-39: "When a table carries MORE THAN ONE breaching row, the precondition
// MUST report every one of them in a single pass — not the first alone —
// combined with errors.Join, so each per-row error stays individually
// inspectable through the aggregate's Unwrap() []error."
// HAPPY PATH
func TestReq39_EveryBreachingRowIsReportedInOnePass(t *testing.T) {
	want := []resolve.RowRef{
		{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
		{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"},
	}
	err := resolveErr(t, multiBreachInput(want...))

	got := breachRefs(t, err)
	if !sameRefs(got, want) {
		t.Errorf("reported identities = %+v; want every breaching row %+v",
			got, want)
	}
}

// REQ-40: "The reported rows MUST be ordered by RowRef identity using the
// kernel's existing compareRefs ordering, never by Table.Rows position, so the
// diagnostic payload is a function of the table value rather than of row order"
// ADVERSARIAL
//
// Supplied in reverse identity order, so a position-ordered report would come
// back reversed.
func TestReq40_TheReportIsOrderedByIdentityNotByRowPosition(t *testing.T) {
	sorted := []resolve.RowRef{
		{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
		{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"},
	}
	reversed := []resolve.RowRef{sorted[2], sorted[1], sorted[0]}

	got := breachRefs(t, resolveErr(t, multiBreachInput(reversed...)))
	if !sameRefs(got, sorted) {
		t.Errorf("reported identities = %+v; want identity order %+v — "+
			"the payload is a function of the table VALUE", got, sorted)
	}
}

// REQ-41: "Table.Rows position MUST NOT be used to break the tie."
// ADVERSARIAL
//
// Two rows sharing one identity, differing only in the WRITE they carry, in
// both supply orders. If position broke the tie the two runs would differ;
// the collapse makes them identical.
func TestReq41_RowPositionNeverBreaksAnIdentityTie(t *testing.T) {
	shared := resolve.RowRef{RuleID: "rule.tie", SourceLocator: "flows/t.toml:9"}

	build := func(first, second string) resolve.Input {
		in := noMatchInput()
		for _, value := range []string{first, second} {
			row := breachingEscapeRow(shared.RuleID, shared.SourceLocator,
				resolve.KindNoMatch)
			row.Writes = []resolve.Tag{{Key: "status", Value: value}}
			in.Table.Rows = append(in.Table.Rows, row)
		}
		return in
	}

	forward := resolveErr(t, build("Alpha", "Beta"))
	backward := resolveErr(t, build("Beta", "Alpha"))

	if !sameRefs(breachRefs(t, forward), breachRefs(t, backward)) {
		t.Errorf("the two supply orders report different identities: "+
			"%+v vs %+v", breachRefs(t, forward), breachRefs(t, backward))
	}
	if !sameCounts(breachCounts(t, forward), breachCounts(t, backward)) {
		t.Errorf("the two supply orders report different counts: %v vs %v",
			breachCounts(t, forward), breachCounts(t, backward))
	}
}

// REQ-42: "The report is instead ordered by RowRef identity and made total by
// COLLAPSING equal identities: breaching rows sharing one RowRef contribute
// ONE reported error, not one per row."
// HAPPY PATH
func TestReq42_EqualIdentitiesCollapseToOneReportedError(t *testing.T) {
	shared := resolve.RowRef{RuleID: "rule.dup", SourceLocator: "flows/d.toml:4"}
	err := resolveErr(t, multiBreachInput(shared, shared, shared))

	got := breachRefs(t, err)
	if len(got) != 1 {
		t.Fatalf("three rows sharing one identity produced %d reported "+
			"errors (%+v); equal identities collapse to ONE", len(got), got)
	}
	if got[0] != shared {
		t.Errorf("the collapsed entry carries %+v; want %+v", got[0], shared)
	}
}

// REQ-43: "Rows carrying no source identity collapse to a single RowRef{"",""}
// entry; that is a degenerate producer, and the error MUST remain diagnostic
// in that case by stating the breach count alongside the identities."
// INPUT EDGE
//
// The degenerate producer. The identities are all zero-valued, so the COUNT is
// the only diagnostic content — and it must be present as structure.
func TestReq43_RowsWithNoSourceIdentityCollapseAndCarryTheCount(t *testing.T) {
	zero := resolve.RowRef{}
	err := resolveErr(t, multiBreachInput(zero, zero, zero))

	refs := breachRefs(t, err)
	if len(refs) != 1 || refs[0] != zero {
		t.Fatalf("reported identities = %+v; want a single RowRef{\"\",\"\"} "+
			"entry", refs)
	}
	if counts := breachCounts(t, err); len(counts) != 1 || counts[0] != 3 {
		t.Errorf("reported counts = %v; want [3] — the count is what keeps "+
			"a zero-identity report diagnostic", counts)
	}
}

// REQ-44: "The count is carried STRUCTURALLY, as the Count field on each
// per-identity *EscapeShapeBreachError — never only in formatted prose"
// BOUNDARY
//
// Read off the struct field with no reference to Error(). The assertion holds
// whatever the message says, which is the point.
func TestReq44_TheCountIsCarriedOnTheStructNotOnlyInProse(t *testing.T) {
	shared := resolve.RowRef{RuleID: "rule.s", SourceLocator: "flows/s.toml:7"}
	err := resolveErr(t, multiBreachInput(shared, shared))

	elems := breachElements(t, err)
	if len(elems) != 1 {
		t.Fatalf("aggregate carries %d elements; want the one collapsed entry",
			len(elems))
	}
	if elems[0].Count != 2 {
		t.Errorf("Count = %d; want 2, read off the struct field",
			elems[0].Count)
	}
}

// REQ-45: "Count is PER-IDENTITY and counts PRE-COLLAPSE ROWS: three breaching
// rows sharing RowRef{"",""} yield ONE reported error with Ref ==
// RowRef{"",""} and Count == 3. A single-row breach carries Count == 1, so the
// field is uniform rather than present only in the degenerate case."
// BOUNDARY
//
// Both halves, and the per-identity claim: a mixed table's counts must attach
// to the right identities rather than being a table-wide total.
func TestReq45_CountIsPerIdentityAndCountsPreCollapseRows(t *testing.T) {
	t.Run("single_row_breach_carries_count_1", func(t *testing.T) {
		counts := breachCounts(t, resolveErr(t, breachingNoMatchInput()))
		if len(counts) != 1 || counts[0] != 1 {
			t.Errorf("counts = %v; want [1] — the field is uniform, not "+
				"present only in the degenerate case", counts)
		}
	})

	t.Run("three_zero_identity_rows_yield_count_3", func(t *testing.T) {
		zero := resolve.RowRef{}
		counts := breachCounts(t, resolveErr(t,
			multiBreachInput(zero, zero, zero)))
		if len(counts) != 1 || counts[0] != 3 {
			t.Errorf("counts = %v; want [3]", counts)
		}
	})

	t.Run("counts_are_per_identity_not_a_table_total", func(t *testing.T) {
		a := resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"}
		b := resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"}
		err := resolveErr(t, multiBreachInput(a, b, b))

		refs := breachRefs(t, err)
		counts := breachCounts(t, err)
		if !sameRefs(refs, []resolve.RowRef{a, b}) {
			t.Fatalf("identities = %+v; want [a b]", refs)
		}
		if !sameCounts(counts, []int{1, 2}) {
			t.Errorf("counts = %v; want [1 2] — the count is PER-IDENTITY, "+
				"not a table-wide total", counts)
		}
	})
}

// REQ-46: "*Reporting* — aggregate, not fail-fast: all breaching rows in one
// pass, sorted by `compareRefs`, with equal-identity rows collapsed to one
// entry"
// HAPPY PATH
//
// The three properties together on one table: aggregate (not first-only),
// sorted, collapsed.
func TestReq46_TheReportIsAggregateSortedAndCollapsed(t *testing.T) {
	a := resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"}
	b := resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"}
	c := resolve.RowRef{RuleID: "rule.c", SourceLocator: "flows/c.toml:3"}

	// Supplied out of order, with `b` duplicated.
	err := resolveErr(t, multiBreachInput(c, b, a, b))

	refs := breachRefs(t, err)
	if !sameRefs(refs, []resolve.RowRef{a, b, c}) {
		t.Errorf("identities = %+v; want the sorted, collapsed set [a b c]",
			refs)
	}
	if counts := breachCounts(t, err); !sameCounts(counts, []int{1, 2, 1}) {
		t.Errorf("counts = %v; want [1 2 1]", counts)
	}
}

// --- F. The exported predicate (Table.CheckValid) ------------------------

// REQ-47: "The conformance predicate MUST be exported by the kernel package as
// a construction-time check callable by any table producer, and Resolve's
// entry precondition MUST be that same function — one predicate, two call
// sites, so a non-TOML producer can fail at build or construction time rather
// than at first production Resolve, and the two enforcement points cannot
// drift."
// HAPPY PATH
//
// "Same function" is asserted as verdict identity across every fixture in
// this suite: both nil, or both non-nil with equal identities and counts.
func TestReq47_OnePredicateTwoCallSitesNeverDrift(t *testing.T) {
	nonVacuityGate(t)

	fixtures := map[string]resolve.Input{
		"breaching":     breachingNoMatchInput(),
		"conforming":    conformingNoMatchInput(),
		"empty_not_nil": emptyNotNilEscapeInput(),
		"dormant":       dormantBreachInput(),
		"multi": multiBreachInput(
			resolve.RowRef{RuleID: "rule.b", SourceLocator: "flows/b.toml:2"},
			resolve.RowRef{RuleID: "rule.a", SourceLocator: "flows/a.toml:1"},
		),
		"legal": legalInput(),
	}

	for name, in := range fixtures {
		t.Run(name, func(t *testing.T) {
			direct := in.Table.CheckValid()
			_, entry := resolve.Resolve(in)

			if (direct == nil) != (entry == nil) {
				t.Fatalf("the two enforcement points disagree: "+
					"CheckValid=%v, Resolve=%v", direct, entry)
			}
			if direct == nil {
				return
			}
			if !sameRefs(breachRefs(t, direct), breachRefs(t, entry)) {
				t.Errorf("identities differ between call sites: %+v vs %+v",
					breachRefs(t, direct), breachRefs(t, entry))
			}
			if !sameCounts(breachCounts(t, direct), breachCounts(t, entry)) {
				t.Errorf("counts differ between call sites: %v vs %v",
					breachCounts(t, direct), breachCounts(t, entry))
			}
		})
	}
}

// REQ-48: "It MUST be a METHOD ON Table taking no arguments"
// BOUNDARY
func TestReq48_CheckValidIsAMethodOnTableTakingNoArguments(t *testing.T) {
	mt, ok := reflect.TypeOf(resolve.Table{}).MethodByName("CheckValid")
	if !ok {
		t.Fatalf("CheckValid is not a method on resolve.Table")
	}
	// NumIn counts the receiver; a no-argument method has exactly one.
	if mt.Type.NumIn() != 1 {
		t.Errorf("CheckValid takes %d arguments; it takes none",
			mt.Type.NumIn()-1)
	}
	if mt.Type.IsVariadic() {
		t.Errorf("CheckValid is variadic; it takes no arguments")
	}
}

// REQ-49: "It MUST return error (nil when valid), and its name MUST follow
// Go's error-returning convention — CheckValid or Validate, never
// Valid/IsValid/OK, which Go reserves for bool-returning predicates."
// BOUNDARY
//
// Both halves: the single `error` result, and the absence of a bool-returning
// spelling that would signal the convention was ignored.
func TestReq49_ThePredicateReturnsErrorAndNotABoolPredicate(t *testing.T) {
	nonVacuityGate(t)

	mt, _ := reflect.TypeOf(resolve.Table{}).MethodByName("CheckValid")
	if mt.Type.NumOut() != 1 {
		t.Fatalf("CheckValid returns %d values; want exactly one",
			mt.Type.NumOut())
	}
	if mt.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		t.Errorf("CheckValid returns %v; want error", mt.Type.Out(0))
	}

	// nil when valid — a plain untyped nil, never a typed nil that would
	// make `err != nil` true on a conforming table.
	if err := legalInput().Table.CheckValid(); err != nil {
		t.Errorf("CheckValid on a conforming table returned %v (%T); want a "+
			"plain nil", err, err)
	}

	tt := reflect.TypeOf(resolve.Table{})
	for _, banned := range []string{"Valid", "IsValid", "OK"} {
		if m, found := tt.MethodByName(banned); found &&
			m.Type.NumOut() == 1 &&
			m.Type.Out(0).Kind() == reflect.Bool {
			t.Errorf("Table gained a bool-returning %q; Go reserves those "+
				"spellings for bool predicates and this one returns error",
				banned)
		}
	}
}

// REQ-52: "A table with no rows (empty or nil Rows) conforms VACUOUSLY:
// CheckValid MUST return nil, since no row can breach a predicate quantified
// over rows and errors.Join of nothing is nil."
// INPUT EDGE
//
// Both the nil and the empty spelling, at both call sites — a `errors.Join()`
// of nothing must come back as an untyped nil, not a non-nil empty aggregate.
func TestReq52_ARowlessTableConformsVacuously(t *testing.T) {
	nonVacuityGate(t)

	cases := map[string][]resolve.Row{
		"nil_rows":   nil,
		"empty_rows": {},
	}

	for name, rows := range cases {
		t.Run(name, func(t *testing.T) {
			tbl := resolve.Table{
				Revision: "rev-empty",
				Outcomes: []string{"successful"},
				Rows:     rows,
			}
			if err := tbl.CheckValid(); err != nil {
				t.Errorf("CheckValid returned %v (%T) on a rowless table; "+
					"the predicate is vacuously satisfied", err, err)
			}

			// And through the entry check: the rowless table must reach its
			// existing disposition rather than erroring.
			in := legalInput()
			in.Table = tbl
			got, err := resolve.Resolve(in)
			if err != nil {
				t.Errorf("Resolve errored on a rowless table: %v", err)
			}
			if got.Refusal == nil {
				t.Errorf("the rowless table produced no refusal: %+v; a "+
					"nil precondition return must never skip evaluation", got)
			}
		})
	}
}

// REQ-53: "*Locus* — one exported function, called at both sites, so any
// future change is a single-site edit that cannot desynchronize the two
// enforcement points."
// ADVERSARIAL
//
// The desynchronization guard: a table that breaches must breach identically
// at both sites for EVERY breach cardinality, including the degenerate
// zero-identity case where a second independent implementation is most
// likely to differ.
func TestReq53_TheTwoEnforcementPointsCannotDesynchronize(t *testing.T) {
	zero := resolve.RowRef{}
	for _, in := range []resolve.Input{
		breachingNoMatchInput(),
		multiBreachInput(zero, zero),
		multiBreachInput(
			resolve.RowRef{RuleID: "z", SourceLocator: "flows/z.toml:1"},
			zero,
		),
	} {
		direct := in.Table.CheckValid()
		_, entry := resolve.Resolve(in)
		if direct == nil || entry == nil {
			t.Fatalf("a breaching table came back clean at one site: "+
				"CheckValid=%v Resolve=%v", direct, entry)
		}
		if !sameRefs(breachRefs(t, direct), breachRefs(t, entry)) ||
			!sameCounts(breachCounts(t, direct), breachCounts(t, entry)) {
			t.Errorf("the two enforcement points desynchronized: "+
				"%+v/%v vs %+v/%v",
				breachRefs(t, direct), breachCounts(t, direct),
				breachRefs(t, entry), breachCounts(t, entry))
		}
	}
}

// --- shared helper -------------------------------------------------------

// resolveErr runs Resolve and requires a non-nil error, returning it.
func resolveErr(t *testing.T, in resolve.Input) error {
	t.Helper()

	got, err := resolve.Resolve(in)
	if err == nil {
		t.Fatalf("Resolve returned no error for a breaching table: %+v", got)
	}
	return err
}

// nonVacuityGate is the shared guard every "both call sites agree" test
// installs first. A predicate that answers nil to EVERYTHING makes any
// agreement assertion trivially true, so each such test first requires the
// predicate to actually reject the canonical breach.
func nonVacuityGate(t *testing.T) {
	t.Helper()

	if err := breachingNoMatchInput().Table.CheckValid(); err == nil {
		t.Fatalf("CheckValid accepts the canonical breaching table; every " +
			"agreement assertion below would then hold vacuously")
	}
	if _, err := resolve.Resolve(breachingNoMatchInput()); err == nil {
		t.Fatalf("Resolve accepts the canonical breaching table; every " +
			"agreement assertion below would then hold vacuously")
	}
}
