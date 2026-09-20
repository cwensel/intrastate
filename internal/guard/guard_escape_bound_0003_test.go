package guard_test

// RDR 0003 — escape rows and the two overlap populations, set literals,
// and the published cardinality bound.

import (
	"maps"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-75: "A declared escape row participates in the coverage identity like
// any other row: its accepted assignments are computed from its guard atoms
// and unioned with its peers'. An escape row carrying no guard atoms
// denotes the whole scoped product and therefore closes coverage by itself.
// Lint MUST NOT treat \"an escape row exists\" as a separate
// coverage-satisfying fact outside the union."
// DOMAIN EDGE
func TestReq75_EscapeRowParticipatesInTheUnionAndABareOneClosesCoverage(t *testing.T) {
	m := mustLoadSource(t, bareEscapeClosesGapSource())
	g := groupOf(t, m, "partial")

	// A bare escape row denotes the WHOLE scoped product.
	bare := guard.AcceptedAssignments(m, rowByID(t, m, "bare-escape"))
	if !bare.Equal(guard.Product(m, g)) {
		t.Errorf("a bare escape row denotes %d assignments; the whole scoped "+
			"product is %d", bare.Len(), guard.Product(m, g).Len())
	}

	// Its assignments come through the UNION, not a separate fact: the
	// ordinary row alone leaves a gap, and the union with the escape row
	// closes it.
	partial := guard.AcceptedAssignments(m, rowByID(t, m, "partial"))
	if partial.Equal(guard.Product(m, g)) {
		t.Fatal("the fixture's ordinary row already closes coverage")
	}
	if !guard.CoverageUnion(m, g).Equal(guard.Product(m, g)) {
		t.Error("the union of the ordinary row and the bare escape row does " +
			"not close coverage; the escape row participates in the union")
	}

	// And a guarded escape row contributes only its OWN assignments, not
	// the whole product — proving it is not an "escape row exists" fact.
	gm := mustLoadSource(t, guardedEscapePartialSource())
	gg := groupOf(t, gm, "ordinary")
	guarded := guard.AcceptedAssignments(gm, rowByID(t, gm, "guarded-escape"))
	if guarded.Equal(guard.Product(gm, gg)) {
		t.Error("a GUARDED escape row denoted the whole product; its accepted " +
			"assignments are computed from its guard atoms")
	}
}

// REQ-76: "**An escape row closes coverage only for the failure classes it
// can actually rescue** — those it declares in its rescue list, and of
// those only `no_match` and `ambiguous_match`. … An escape row therefore
// cannot rescue a class it does not declare, nor either blocking class,
// however bare its guard."
// BOUNDARY
func TestReq76_EscapeClosesOnlyItsDeclaredRescuableClasses(t *testing.T) {
	if got := slices.Sorted(slices.Values(guard.RescuableClasses())); !slices.Equal(
		got, []string{string(resolve.KindAmbiguousMatch), string(resolve.KindNoMatch)}) {
		t.Fatalf("RescuableClasses() = %v; want exactly no_match and "+
			"ambiguous_match", got)
	}

	// Neither blocking class is rescuable, however bare the guard.
	for _, blocking := range []resolve.RefusalKind{
		resolve.KindGuardUnevaluable,
		resolve.KindOwnedStateUnavailable,
		resolve.KindUnmodeledOutcome,
	} {
		if slices.Contains(guard.RescuableClasses(), string(blocking)) {
			t.Errorf("%q is rescuable; an escape row cannot rescue a blocking "+
				"class", blocking)
		}
	}

	m := mustLoadSource(t, bareEscapeClosesGapSource())
	g := groupOf(t, m, "partial")
	if !guard.CoverageUnionFor(m, g, string(resolve.KindNoMatch)).Equal(guard.Product(m, g)) {
		t.Error("the escape row did not close its DECLARED class no_match")
	}
	if guard.CoverageUnionFor(m, g, string(resolve.KindAmbiguousMatch)).Equal(guard.Product(m, g)) {
		t.Error("an escape row declaring only no_match closed the " +
			"ambiguous_match arm; it cannot rescue a class it does not declare")
	}
}

// REQ-77: "The coverage union is consequently computed per (scoped row
// group × declared rescuable class), and a row declaring one class MUST NOT
// close the group's other arm"
// DOMAIN EDGE
func TestReq77_CoverageUnionIsPerGroupTimesDeclaredRescuableClass(t *testing.T) {
	// The union is per (group × declared class) whatever the population:
	// this holds on an overlap-free group too, and is REQ-77's own clause.
	m := mustLoadSource(t, bareEscapeClosesGapSource())
	g := groupOf(t, m, "partial")

	noMatch := guard.CoverageUnionFor(m, g, string(resolve.KindNoMatch))
	ambiguous := guard.CoverageUnionFor(m, g, string(resolve.KindAmbiguousMatch))
	if noMatch.Equal(ambiguous) {
		t.Error("the two arms computed one union; the union is per (group × " +
			"declared rescuable class)")
	}

	// The unclosed arm draws its gap finding — witnessed on a group whose
	// ordinary population OVERLAPS, so the `ambiguous_match` arm is
	// reachable. RDR 0006 states the arm mechanics and REQ-77 cites them
	// (`0003:C15`), so an overlap-free group is not a witness for this half:
	// there the arm is vacuously closed and drawing a gap would demand a
	// rescue row `graph-unreachable-rule` then flags.
	reach := mustLoadSource(t, overlappingOrdinaryBareEscapeSource())
	if !anyFinding(guard.Lint(reach), func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap && f.Class == string(resolve.KindAmbiguousMatch)
	}) {
		t.Errorf("the arm the escape row does not declare produced no gap "+
			"finding; findings=%v", allFindings(guard.Lint(reach)))
	}

	// And the arm is NOT demanded where it cannot be reached: the same bare
	// escape row over an overlap-free population draws no `ambiguous_match`
	// gap at all.
	if anyFinding(guard.Lint(m), func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap && f.Class == string(resolve.KindAmbiguousMatch)
	}) {
		t.Errorf("an overlap-free group drew an `ambiguous_match` gap; that "+
			"arm is vacuously closed where it is unreachable; findings=%v",
			allFindings(guard.Lint(m)))
	}
}

// REQ-78: "a bare escape row MUST NOT be read as discharging the narrowing:
// a group whose ordinary row can refuse `guard_unevaluable` still has its
// claim withheld, because at runtime that refusal is returned before the
// escape row is ever consulted."
// ADVERSARIAL
func TestReq78_BareEscapeRowDoesNotDischargeTheNarrowing(t *testing.T) {
	withoutEscape := guard.Lint(mustLoadSource(t, optionalKeyGroupSource(false)))
	withEscape := guard.Lint(mustLoadSource(t, optionalKeyPlusBareEscapeSource()))

	if hasGreen(withEscape) {
		t.Fatalf("adding a bare escape row turned a withheld group green; the "+
			"runtime returns the refusal before the escape row is consulted. "+
			"reports=%s", renderReports(withEscape))
	}
	if verdictOf(withoutEscape) != verdictOf(withEscape) {
		t.Errorf("the verdict changed when a bare escape row was added (%q → "+
			"%q); the escape row cannot rescue guard_unevaluable",
			verdictOf(withoutEscape), verdictOf(withEscape))
	}
	// Discriminating leg: the claim must actually BE withheld in both runs.
	// Two runs that both emit nothing agree trivially.
	for name, reports := range map[string][]guard.GroupReport{
		"without escape": withoutEscape, "with bare escape": withEscape,
	} {
		if countCode(reports, guard.CodeUnprovableCoverage) == 0 {
			t.Errorf("the %s run emitted no withholding finding; the claim is "+
				"withheld on both, not absent on both. findings=%v",
				name, allFindings(reports))
		}
	}

	// And the runtime agrees: the refusal is guard_unevaluable, not a rescue.
	m := mustLoadSource(t, optionalKeyPlusBareEscapeSource())
	res := resolveWith(t, m.KernelTable(), guard.View{})
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("the runtime gave %v; the gate's refusal is returned before "+
			"escapeOrRefuse is reachable", describe(res))
	}
}

// REQ-79: "Overlap is checked in **two separate populations** … an escape
// row overlapping a *guarded* row is **not** a runtime ambiguity and MUST
// NOT be reported as one. What lint MUST still check is overlap **among
// escape rows for the same failure class** … Escape rows are therefore
// excluded from the ordinary-row overlap check and subjected to their own;
// excluding them from **coverage** is what MUST NOT happen."
// ADVERSARIAL
func TestReq79_OverlapIsTwoPopulationsWhileCoverageIsOne(t *testing.T) {
	m := mustLoadSource(t, twoPopulationSource())
	reports := guard.Lint(m)

	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		hasEscape := slices.Contains(f.RuleIDs, "esc-a") || slices.Contains(f.RuleIDs, "esc-b")
		hasOrdinary := slices.Contains(f.RuleIDs, "ord-a") || slices.Contains(f.RuleIDs, "ord-b")
		if hasEscape && hasOrdinary {
			t.Errorf("overlap finding %v mixes the two populations; an escape "+
				"row overlapping a guarded row is not a runtime ambiguity", f.RuleIDs)
		}
	}

	// Overlap among escape rows for one class IS still checked.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			slices.Contains(f.RuleIDs, "esc-a") && slices.Contains(f.RuleIDs, "esc-b")
	}) {
		t.Errorf("two overlapping escape rows for one class drew no finding; "+
			"findings=%v", allFindings(reports))
	}

	// And they are NOT excluded from coverage.
	g := groupOf(t, m, "ord-a")
	union := guard.CoverageUnionFor(m, g, string(resolve.KindNoMatch))
	ordinaryOnly := guard.AcceptedAssignments(m, rowByID(t, m, "ord-a")).
		Union(guard.AcceptedAssignments(m, rowByID(t, m, "ord-b")))
	if union.Equal(ordinaryOnly) {
		t.Error("the escape rows were excluded from coverage; excluding them " +
			"from coverage is what MUST NOT happen")
	}
}

// REQ-80: "lint MUST partition the escape rows into **one population per
// declared failure class**, place a row declaring several classes in
// **each** of those populations, and run the pairwise overlap check within
// each population independently."
// ADVERSARIAL
func TestReq80_EscapeRowsArePartitionedPerClassWithMultiClassRowsInEach(t *testing.T) {
	m := mustLoadSource(t, multiClassEscapeSource())
	reports := guard.Lint(m)

	// `both` declares both classes; `only-no-match` declares one. They
	// collide on no_match only. A per-row partition would miss it; a flat
	// pairing over all escape rows would also pair `only-no-match` with
	// `only-ambiguous`, which share no class.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			f.Class == string(resolve.KindNoMatch) &&
			slices.Contains(f.RuleIDs, "both") &&
			slices.Contains(f.RuleIDs, "only-no-match")
	}) {
		t.Errorf("the multi-class row was not placed in the no_match "+
			"population; findings=%v", allFindings(reports))
	}
	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		if slices.Contains(f.RuleIDs, "only-no-match") &&
			slices.Contains(f.RuleIDs, "only-ambiguous") {
			t.Errorf("finding %v pairs two escape rows sharing no class; the "+
				"check runs within each population independently", f.RuleIDs)
		}
	}
}

// REQ-81: "A pair overlapping in more than one class is one finding per
// class, naming the class, so the author can see which rescue path is
// disabled."
// BOUNDARY
func TestReq81_APairOverlappingInTwoClassesYieldsOneFindingPerClass(t *testing.T) {
	m := mustLoadSource(t, twoClassOverlapPairSource())
	reports := guard.Lint(m)

	classes := map[string]bool{}
	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		if slices.Contains(f.RuleIDs, "pair-a") && slices.Contains(f.RuleIDs, "pair-b") {
			if f.Class == "" {
				t.Errorf("overlap finding %v names no class; the author must "+
					"see which rescue path is disabled", f.RuleIDs)
			}
			classes[f.Class] = true
		}
	}
	want := []string{string(resolve.KindAmbiguousMatch), string(resolve.KindNoMatch)}
	got := slices.Sorted(maps.Keys(classes))
	if !slices.Equal(got, want) {
		t.Errorf("the pair drew findings for classes %v; want one finding per "+
			"class %v", got, want)
	}
}

// REQ-82: "**A bare escape row is not a silent opt-out from the
// guarantee.** … Lint MUST therefore report a bare escape row — one
// carrying no guard atoms — that closes a group's coverage as an
// **observable** result: the group's verdict names the escape row that
// closed it … Emitting a bare green for such a group MUST NOT satisfy this
// clause."
// ADVERSARIAL
func TestReq82_BareEscapeClosureIsObservableAndNamesTheRow(t *testing.T) {
	m := mustLoadSource(t, bareEscapeClosesGapSource())
	reports := guard.Lint(m)
	r := reportFor(t, reports, "partial")

	if r.ClosedByEscape == "" {
		t.Errorf("the group's verdict does not name the escape row that "+
			"closed it; verdict=%s findings=%v", r.Verdict, r.Findings)
	}
	if r.ClosedByEscape != "bare-escape" {
		t.Errorf("the verdict names %q; want the bare escape row "+
			"`bare-escape`", r.ClosedByEscape)
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageClosedByEscape &&
			slices.Contains(f.RuleIDs, "bare-escape")
	}) {
		t.Errorf("no observable result named the bare escape row; findings=%v",
			allFindings(reports))
	}

	// A group proved over its declared domains is DISTINGUISHABLE from one
	// closed by a catch-all: the proved twin names no escape row.
	proved := reportFor(t, guard.Lint(mustLoadSource(t, completePartitionSource())), "part-small")
	if proved.ClosedByEscape != "" {
		t.Errorf("a group proved over its declared domains named escape row "+
			"%q; the reader must be able to tell the two apart", proved.ClosedByEscape)
	}
}

// REQ-83: "A set literal — the right-hand side of `in` and of `contains` —
// has one canonical spelling: an unordered set of typed elements,
// duplicate-free, and compared as a set. Two authored spellings differing
// only in element order or in repeated elements MUST parse to the same
// literal and therefore to the same atom under the identity tuple; a
// repeated element MUST be rejected at parse rather than silently
// collapsed."
// ADVERSARIAL
func TestReq83_SetLiteralIsCanonicalUnorderedAndDuplicateFree(t *testing.T) {
	forward := mustLoadSource(t, setLiteralOrderSource(`["bug", "chore"]`))
	reversed := mustLoadSource(t, setLiteralOrderSource(`["chore", "bug"]`))

	a := rowByID(t, forward, "set-lit")
	b := rowByID(t, reversed, "set-lit")
	if guard.AtomIdentity(a, a.Atoms[0]) != guard.AtomIdentity(b, b.Atoms[0]) {
		t.Errorf("two element orderings produced different atom identities "+
			"(%v vs %v); they parse to the SAME literal",
			guard.AtomIdentity(a, a.Atoms[0]), guard.AtomIdentity(b, b.Atoms[0]))
	}

	// A repeated element is REJECTED at parse, not silently collapsed.
	decls := `
[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
required = true
`
	if cat := loadCategory(t, decls, "[rule.guard.all.labels]\ncontains = [\"bug\", \"bug\"]\n"); cat != table.CatMalformedPredicateAtom {
		t.Errorf("a repeated `contains` element refused as %q; want %q — "+
			"rejected at parse rather than silently collapsed",
			cat, table.CatMalformedPredicateAtom)
	}
}

// REQ-84: "Implementations MUST canonicalize before the literal enters the
// identity tuple, so reordering a set literal cannot change a diagnostic —
// the source-order independence MVV Scenario 6 requires."
// ADVERSARIAL
func TestReq84_CanonicalizationHappensBeforeTheIdentityTuple(t *testing.T) {
	forward := guard.Lint(mustLoadSource(t, setOverlapSource(`["bug", "chore"]`)))
	reversed := guard.Lint(mustLoadSource(t, setOverlapSource(`["chore", "bug"]`)))

	if verdictOf(forward) != verdictOf(reversed) {
		t.Errorf("reordering a set literal changed the verdict (%q → %q); "+
			"canonicalization happens before the identity tuple",
			verdictOf(forward), verdictOf(reversed))
	}
	if countCode(forward, guard.CodeOverlap) != countCode(reversed, guard.CodeOverlap) {
		t.Errorf("reordering a set literal changed the diagnostic count "+
			"(%d vs %d)", countCode(forward, guard.CodeOverlap),
			countCode(reversed, guard.CodeOverlap))
	}

	// Discriminating leg: there must BE a diagnostic for reordering to be
	// unable to change. Two empty runs agree trivially.
	if countCode(forward, guard.CodeOverlap) == 0 {
		t.Fatalf("the overlapping `contains` literals drew no diagnostic; "+
			"canonicalization is what keeps an EXISTING diagnostic stable. "+
			"findings=%v", allFindings(forward))
	}
}

// REQ-85: "Set-valued guard domains MUST be proved with a deterministic
// symbolic or bitset-equivalent representation. If the finite product is
// too large for that proof, lint MUST refuse or downgrade the exhaustiveness
// claim rather than silently capping enumeration."
// BOUNDARY
func TestReq85_SetDomainsAreProvedDeterministicallyAndOverLargeRefuses(t *testing.T) {
	// Deterministic: the same model gives the same verdict on repeated runs.
	m := mustLoadSource(t, setAndEnumSource())
	first := verdictOf(guard.Lint(m))
	for range 5 {
		if got := verdictOf(guard.Lint(m)); got != first {
			t.Fatalf("repeated lint runs disagreed (%q vs %q); the proof "+
				"representation is deterministic", first, got)
		}
	}

	// Too large: refuse rather than silently capping.
	over := guard.Lint(mustLoadSource(t, overLargeProductSource()))
	if hasGreen(over) {
		t.Fatalf("an over-large finite product certified green; lint silently "+
			"capped enumeration. reports=%s", renderReports(over))
	}
	if countCode(over, guard.CodeProductTooLarge) == 0 {
		t.Errorf("no graph-product-too-large finding; findings=%v", allFindings(over))
	}
}

// REQ-86: "\"Too large to prove\" MUST be a declared, model-independent
// bound, not an implementation's incidental limit: the implementation MUST
// publish the bound it enforces, and the same model MUST receive the same
// verdict on every conforming implementation. A bound discovered by
// exhausting memory or wall-clock is not a conforming bound."
// BOUNDARY
func TestReq86_TheBoundIsPublishedDeclaredAndModelIndependent(t *testing.T) {
	b := guard.Bound()
	if b <= 0 {
		t.Fatalf("Bound() = %d; the implementation MUST publish the bound it "+
			"enforces", b)
	}

	// Model-independent: the same constant whatever model is linted.
	for _, src := range []string{completePartitionSource(), overLargeProductSource(), setAndEnumSource()} {
		_ = guard.Lint(mustLoadSource(t, src))
		if guard.Bound() != b {
			t.Fatalf("Bound() changed to %d while linting a model; the bound "+
				"is not per-model", guard.Bound())
		}
	}
}

// REQ-87: "The refusal diagnostic MUST report the product size it computed
// and the bound it exceeded, so an author can tell an over-large product
// from an undeclared dimension."
// HAPPY PATH
func TestReq87_OverLargeRefusalReportsComputedSizeAndBound(t *testing.T) {
	reports := guard.Lint(mustLoadSource(t, overLargeProductSource()))

	found := findingsWithCode(reports, guard.CodeProductTooLarge)
	if len(found) == 0 {
		t.Fatalf("no over-large refusal; findings=%v", allFindings(reports))
	}
	for _, f := range found {
		if f.ComputedSize <= 0 {
			t.Errorf("finding reports computed size %d; it MUST report the "+
				"product size it computed", f.ComputedSize)
		}
		if f.Bound != guard.Bound() {
			t.Errorf("finding reports bound %d; want the published bound %d",
				f.Bound, guard.Bound())
		}
		if f.ComputedSize <= f.Bound {
			t.Errorf("finding reports size %d not exceeding bound %d",
				f.ComputedSize, f.Bound)
		}
	}

	// An undeclared dimension is tellable apart: it carries NO size.
	undeclared := guard.Lint(mustLoadSource(t, unboundedIntGuardSource()))
	for _, f := range allFindings(undeclared) {
		if f.ComputedSize != 0 {
			t.Errorf("an undeclared-dimension finding reports size %d; the "+
				"two classes must be tellable apart", f.ComputedSize)
		}
	}
}

// REQ-88: "The quantity both figures report is the **cardinality of the
// scoped product** — the number of assignments in it, the product of every
// participating dimension's **assignment count** — not a bitset width, byte
// size, or row count."
// BOUNDARY
func TestReq88_TheReportedQuantityIsTheProductCardinality(t *testing.T) {
	m := mustLoadSource(t, productOnlyGapSource())
	g := groupOf(t, m, "gap-a")

	card, ok := guard.Cardinality(m, g)
	if !ok {
		t.Fatal("the group's product has no cardinality")
	}
	// single-valued always-present enum {small,large} = 2; bool = 2.
	if card != 4 {
		t.Errorf("Cardinality = %d; want 4 — the PRODUCT of each dimension's "+
			"assignment count (2 × 2), not a bitset width, byte size, or row "+
			"count (the group carries 2 rows)", card)
	}
	if card == len(g.Rows) {
		t.Error("Cardinality equals the row count; it is the number of " +
			"ASSIGNMENTS in the scoped product")
	}
	if card != guard.Product(m, g).Len() {
		t.Errorf("Cardinality %d != the enumerated product's size %d",
			card, guard.Product(m, g).Len())
	}
}

// REQ-89: The per-kind assignment-count table: "`enum`, `int`
// (single-valued)" → `|domain|`; "`bool` (single-valued)" → 2; "`enum`,
// `bool`, `int` **without** the single-valued marker" → "`2^|domain|` — one
// independent boolean dimension per value"; "`set`" → "`2^|element
// universe|` … **never** `|universe|`"; "any optional key" → "multiplied by
// 2 for its `{present, absent}` presence dimension".
// BOUNDARY
func TestReq89_PerKindAssignmentCountTableIsEnforcedRowByRow(t *testing.T) {
	min0, max3 := 0, 3
	cases := []struct {
		name string
		decl table.TagDecl
		want int
	}{
		{"single-valued enum is |domain|", table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c"}, SingleValued: true, Required: true}, 3},
		{"single-valued int is |domain|", table.TagDecl{Kind: "int", Min: &min0, Max: &max3, SingleValued: true, Required: true}, 4},
		{"single-valued bool is 2", table.TagDecl{Kind: "bool", SingleValued: true, Required: true}, 2},
		{"unmarked enum is 2^|domain|", table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c"}, Required: true}, 8},
		{"unmarked bool is 2^2", table.TagDecl{Kind: "bool", Required: true}, 4},
		{"unmarked int is 2^|domain|", table.TagDecl{Kind: "int", Min: &min0, Max: &max3, Required: true}, 16},
		{"set is 2^|universe|", table.TagDecl{Kind: "set", Elements: []string{"x", "y", "z"}, Required: true}, 8},
		{"optional key multiplies by 2", table.TagDecl{Kind: "bool", SingleValued: true}, 4},
		{"optional set multiplies by 2", table.TagDecl{Kind: "set", Elements: []string{"x", "y"}}, 8},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := guard.AssignmentCount(tc.decl)
			if !ok {
				t.Fatalf("AssignmentCount(%+v) reported no finite count", tc.decl)
			}
			if got != tc.want {
				t.Errorf("AssignmentCount(%+v) = %d; want %d", tc.decl, got, tc.want)
			}
		})
	}

	// "never |universe|": the set row must not be read as its universe size.
	if n, _ := guard.AssignmentCount(table.TagDecl{
		Kind: "set", Elements: []string{"x", "y", "z"}, Required: true,
	}); n == 3 {
		t.Error("a set dimension was counted as |universe| = 3; a set-valued " +
			"tag holds any SUBSET, so its assignment space is the powerset")
	}
}

// REQ-90: "**The table reads the declaration's defaults, not an author's
// intent.** … An implementation MUST apply these defaults when computing
// the cardinality — reading an unmarked declaration as single-valued or
// always-present is the same exponential understatement the `set` row
// forbids"
// ADVERSARIAL
func TestReq90_DefaultsAreAppliedNotTheAuthorsPresumedIntent(t *testing.T) {
	// An unmarked, unrequired enum {a,b,c,d}: 2^4 × 2 = 32. Reading it as
	// single-valued and always-present would give 4.
	unmarked := table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c", "d"}}
	n, ok := guard.AssignmentCount(unmarked)
	if !ok {
		t.Fatal("AssignmentCount reported no finite count for a declared enum")
	}
	if n == 4 {
		t.Fatalf("AssignmentCount(unmarked, optional enum {a,b,c,d}) = 4; the " +
			"declaration's DEFAULTS give 2^4 × 2 = 32 — reading it as " +
			"single-valued and always-present is the exponential understatement " +
			"the clause forbids")
	}
	if n != 32 {
		t.Errorf("AssignmentCount(unmarked, optional enum {a,b,c,d}) = %d; "+
			"want 32 (2^4 for the unmarked marker, ×2 for the omitted "+
			"optionality marker)", n)
	}
}

// REQ-91: "The bound is a single integer constant published by the
// implementation and reported in the diagnostic beside the computed
// cardinality; it is not per-model, per-group, or configurable per run"
// BOUNDARY
func TestReq91_TheBoundIsOneConstantNotPerModelPerGroupOrConfigurable(t *testing.T) {
	b := guard.Bound()

	// Not per-group: every group of a multi-group model reports the same.
	m := mustLoadSource(t, overLargeProductSource())
	for _, f := range findingsWithCode(guard.Lint(m), guard.CodeProductTooLarge) {
		if f.Bound != b {
			t.Errorf("a group reported bound %d; the bound is a single "+
				"constant, not per-group", f.Bound)
		}
	}

	// Not configurable per run: repeated runs report the same constant.
	for range 3 {
		if guard.Bound() != b {
			t.Fatalf("Bound() = %d then %d across runs; it is not "+
				"configurable per run", b, guard.Bound())
		}
	}

	// Discriminating leg: the constant must be PUBLISHED and REPORTED. A
	// bound of zero is stable across every group and run while enforcing
	// nothing.
	if b <= 0 {
		t.Fatalf("Bound() = %d; the bound is a published integer constant", b)
	}
	if countCode(guard.Lint(m), guard.CodeProductTooLarge) == 0 {
		t.Fatalf("the over-large model drew no refusal to report the bound "+
			"beside; findings=%v", allFindings(guard.Lint(m)))
	}
}

// REQ-92: "Two conforming implementations MAY publish different bounds, but
// each MUST return the same verdict for the same model and MUST report
// which bound it applied."
// BOUNDARY
func TestReq92_TheVerdictIsStableAndTheAppliedBoundIsReported(t *testing.T) {
	m := mustLoadSource(t, overLargeProductSource())

	first := verdictOf(guard.Lint(m))
	for range 3 {
		if got := verdictOf(guard.Lint(m)); got != first {
			t.Fatalf("the same model received verdicts %q then %q", first, got)
		}
	}

	found := findingsWithCode(guard.Lint(m), guard.CodeProductTooLarge)
	if len(found) == 0 {
		t.Fatal("no over-large refusal to carry an applied bound")
	}
	for _, f := range found {
		if f.Bound == 0 {
			t.Error("the refusal reports no applied bound; each implementation " +
				"MUST report which bound it applied")
		}
	}
}

// REQ-93: "**A dimension with no finite declared domain has no assignment
// count**, so a product containing one has no cardinality … The too-large
// comparison is therefore defined **only over a fully-provable product** —
// lint MUST compute the cardinality and test the bound after every
// participating dimension is known finite, and MUST NOT report a computed
// size for a product carrying an unprovable dimension."
// ADVERSARIAL
func TestReq93_NoCardinalityIsComputedForAProductCarryingAnUnprovableDimension(t *testing.T) {
	m := mustLoadSource(t, unprovablePlusLargeSource())
	g := groupOf(t, m, "mixed-a")

	if _, ok := guard.Cardinality(m, g); ok {
		t.Error("a product carrying an unprovable dimension reported a " +
			"cardinality; the table defines no value for such a dimension")
	}

	reports := guard.Lint(m)
	for _, f := range allFindings(reports) {
		if f.ComputedSize != 0 {
			t.Errorf("finding %v reports computed size %d for a group with an "+
				"unprovable dimension; the size-bearing diagnostic is "+
				"unavailable when its quantity is undefined", f, f.ComputedSize)
		}
	}
}

// REQ-94: "each unprovable dimension draws its own blocking finding naming
// that dimension, and the over-large refusal — the one that carries
// `(computed size, bound)` — is simply not among the findings for such a
// group"
// ADVERSARIAL
func TestReq94_UnprovableDimensionsAreReportedAndTheOverLargeRefusalIsNot(t *testing.T) {
	m := mustLoadSource(t, unprovablePlusLargeSource())
	reports := guard.Lint(m)

	if n := countCode(reports, guard.CodeProductTooLarge); n != 0 {
		t.Errorf("a group with an unprovable dimension drew %d "+
			"graph-product-too-large findings; that refusal is simply not "+
			"among the findings for such a group", n)
	}
	named := map[string]bool{}
	for _, f := range findingsWithCode(reports, guard.CodeUnprovableCoverage) {
		named[f.Dimension] = true
	}
	if !named["loose"] {
		t.Errorf("the unprovable dimension `loose` drew no blocking finding "+
			"naming it; named=%v", named)
	}
}

// REQ-94 + RDR 0006 C16: "Withholding a group's exhaustiveness claim MUST
// NOT suppress overlap, coverage, or further withholding findings for that
// group." The clause is unconditional, so the over-bound branch that
// withholds on a STRUCTURALLY unprojectable dimension owes the same
// surviving overlap check as the no-finite-domain withholding path does.
// ADVERSARIAL
func TestReq94AndC16_OverBoundStructuralWithholdingStillEmitsDecidableOverlap(t *testing.T) {
	m := mustLoadSource(t, overBoundStructuralPlusOverlapSource())
	g := groupOf(t, m, "decidable-a")

	// The fixture must actually reach the branch under test: a computable
	// cardinality over the bound, with a structurally unprojectable
	// dimension among the participants. Asserting it here keeps a later
	// declaration edit from silently routing the test down another path.
	card, ok := guard.Cardinality(m, g)
	if !ok || card <= guard.Bound() {
		t.Fatalf("Cardinality = (%d, %v); the fixture must present a "+
			"computable product OVER the bound %d to reach the branch under "+
			"test", card, ok, guard.Bound())
	}

	reports := guard.Lint(m)

	// The withholding itself still happens, and it names the structurally
	// unprovable dimension rather than reporting a size.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "loose"
	}) {
		t.Errorf("no withholding finding named the unprovable dimension "+
			"`loose`; findings=%v", allFindings(reports))
	}

	// And the overlap between the two rows that guard only the DECIDABLE
	// dimension survives it. Their accepted assignments are computed over
	// the decidable sub-product, so the ambiguity is decidable without ever
	// projecting `loose`.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			slices.Contains(f.RuleIDs, "decidable-a") &&
			slices.Contains(f.RuleIDs, "decidable-b")
	}) {
		t.Errorf("overlap among the group's decidable rows was suppressed by "+
			"the over-bound structural withholding; C16 forbids the "+
			"suppression. findings=%v", allFindings(reports))
	}
}

// The surviving overlap check must be a CHECK, not an unconditional
// finding: the same fixture with the two decidable rows made mutually
// exclusive reports no overlap.
// DISCRIMINATING
func TestReq94AndC16_OverBoundStructuralWithholdingReportsNoSpuriousOverlap(t *testing.T) {
	m := mustLoadSource(t, overBoundStructuralDisjointSource())
	reports := guard.Lint(m)

	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "loose"
	}) {
		t.Fatalf("the fixture no longer withholds on `loose`; findings=%v",
			allFindings(reports))
	}
	if n := countCode(reports, guard.CodeOverlap); n != 0 {
		t.Errorf("mutually exclusive decidable rows drew %d graph-overlap "+
			"findings; the surviving check intersects accepted assignments, "+
			"it does not fire on the branch being taken. findings=%v",
			n, allFindings(reports))
	}
}

// Negative control: when EVERY row in the group guards the unprovable
// dimension, no row has an accepted-assignment set, so the surviving
// overlap check correctly reports nothing. This is the case REQ-94's own
// fixture exercises, and running overlap on this path must not change it.
// ADVERSARIAL
func TestReq94AndC16_OverBoundStructuralWithholdingFindsNoOverlapWhenNoRowIsDecidable(t *testing.T) {
	m := mustLoadSource(t, overBoundStructuralAllUnprovableSource())
	reports := guard.Lint(m)

	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "loose"
	}) {
		t.Fatalf("the fixture no longer withholds on `loose`; findings=%v",
			allFindings(reports))
	}
	if n := countCode(reports, guard.CodeOverlap); n != 0 {
		t.Errorf("two rows that both guard the unprovable dimension drew %d "+
			"graph-overlap findings; neither has an accepted-assignment set "+
			"to intersect. findings=%v", n, allFindings(reports))
	}
}

// REQ-77 / RDR 0006: the `ambiguous_match` coverage union is drawn from the
// ESCAPE rows declaring the class, and nothing else. That refusal arises
// where the group's ordinary rows overlap, so the rows that caused the
// ambiguity cannot also rescue it — the kernel reaches the escape arm
// precisely because none of them was the exact-one match. A union counting
// ordinary rows would let a fully covering overlap suppress the gap.
// ADVERSARIAL
func TestReq77_AmbiguousMatchCoverageIsDrawnFromEscapeRowsAlone(t *testing.T) {
	// Overlapping ordinary rows that jointly cover the whole product, with
	// NO row declaring `ambiguous_match`: the arm is reachable and unrescued.
	m := mustLoadSource(t, overlappingOrdinaryCoversProductSource())
	g := groupOf(t, m, "ov-a")

	if !guard.CoverageUnion(m, g).Equal(guard.Product(m, g)) {
		t.Fatal("the fixture's ordinary population no longer covers the " +
			"whole product; it cannot witness the suppression")
	}
	if union := guard.EscapeUnionFor(m, g, string(resolve.KindAmbiguousMatch)); union.Len() != 0 {
		t.Errorf("the escape-only union over a group with no "+
			"`ambiguous_match` row holds %d assignments; nothing declares "+
			"the class", union.Len())
	}

	reports := guard.Lint(m)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap
	}) {
		t.Fatalf("the fixture drew no overlap finding, so the "+
			"`ambiguous_match` arm is not reachable; findings=%v",
			allFindings(reports))
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap && f.Class == string(resolve.KindAmbiguousMatch)
	}) {
		t.Errorf("a fully covering OVERLAP suppressed the `ambiguous_match` "+
			"gap; the rows that mint the ambiguity cannot certify it "+
			"rescued. findings=%v", allFindings(reports))
	}

	// A real `ambiguous_match`-declaring escape row DOES rescue the arm: the
	// same overlapping population draws no gap once one is present.
	rescued := guard.Lint(mustLoadSource(t, overlappingOrdinaryAmbiguousEscapeSource()))
	if anyFinding(rescued, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap && f.Class == string(resolve.KindAmbiguousMatch)
	}) {
		t.Errorf("a bare escape row declaring `ambiguous_match` failed to "+
			"close the arm it declares; findings=%v", allFindings(rescued))
	}

	// And where that escape row is the SOLE closer — the ordinary rows
	// overlap AND leave a gap — the closure is reported, not taken silently.
	closed := guard.Lint(mustLoadSource(t, overlappingOrdinaryGapAmbiguousEscapeSource()))
	if anyFinding(closed, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap
	}) {
		t.Errorf("the escape row declaring both rescuable classes left an "+
			"arm open; findings=%v", allFindings(closed))
	}
	if !anyFinding(closed, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageClosedByEscape
	}) {
		t.Errorf("the escape row closed coverage silently; the closure is an "+
			"OBSERVABLE result. findings=%v", allFindings(closed))
	}
}
