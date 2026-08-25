package guard_test

// RDR 0003 — the Testing Strategy scenarios and the Minimum Viable
// Validation.
//
// TestMVV_GuardPredicateExhaustiveness is the gating runnable end-to-end
// validation: it decomposes into the five obligations the REQ list names,
// including the POSITIVE assertion on the narrowing (a run that emits
// nothing must fail) and the negative control that proves the narrowing is
// tight rather than blanket.

import (
	"maps"
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-124: "Build the MVV fixture against representative RDR and kata
// guards and use the result to confirm or adjust the initial operator set.
// … the implementation MVV must add at least one `contains` predicate over
// a declared set-valued tag before the full operator vocabulary is
// accepted."
// HAPPY PATH
func TestReq124_MVVFixtureCarriesAContainsPredicateOverADeclaredSetTag(t *testing.T) {
	rdr := mustLoadSource(t, mvvRDRSlice())
	kata := mustLoadSource(t, mvvKataSlice())

	var contains []table.Atom
	var setKeys []string
	for _, m := range []*table.Model{rdr, kata} {
		for _, row := range m.Rows {
			for _, a := range row.Atoms {
				if a.Operator != "contains" {
					continue
				}
				contains = append(contains, a)
				if guard.DeclarationOf(m, a.Key).Kind != "set" {
					t.Errorf("`contains` atom over %q whose declared kind is "+
						"%q; it must be over a declared SET-valued tag",
						a.Key, guard.DeclarationOf(m, a.Key).Kind)
				}
				if len(guard.DeclarationOf(m, a.Key).Elements) == 0 {
					t.Errorf("set tag %q declares no element universe", a.Key)
				}
				setKeys = append(setKeys, a.Key)
			}
		}
	}
	if len(contains) == 0 {
		t.Fatal("the MVV fixture carries no `contains` predicate; the full " +
			"operator vocabulary is not accepted without one over a declared " +
			"set-valued tag")
	}

	// The full vocabulary is exercised across the two slices.
	seen := map[string]bool{}
	for _, m := range []*table.Model{rdr, kata} {
		for _, row := range m.Rows {
			for _, a := range row.Atoms {
				if a.Block != table.BlockMatch {
					seen[a.Operator] = true
				}
			}
		}
	}
	for _, op := range []string{"eq", "in", "contains", "exists", "gte"} {
		if !seen[op] {
			t.Errorf("the MVV fixture exercises no %q guard atom; operators "+
				"seen=%v", op, slices.Sorted(maps.Keys(seen)))
		}
	}
	if !seen["lt"] && !seen["lte"] && !seen["gt"] && !seen["gte"] {
		t.Error("the MVV fixture carries no bounded integer comparison")
	}
}

// REQ-126: "Phase 3 must therefore declare the markers **and** close the
// routing, or author a separate group for Scenario 2's complete partition —
// reusing it unchanged makes Scenario 2 fail on first run against a correct
// lint"
// HAPPY PATH
func TestReq126_ScenarioTwosPartitionGroupDeclaresMarkersAndClosesRouting(t *testing.T) {
	m := mustLoadSource(t, mvvRDRSlice())
	g := groupOf(t, m, "partition-small")

	// Every guard dimension of the partition group declares the marker.
	for _, key := range guard.Dimensions(m, g) {
		d := guard.DeclarationOf(m, key)
		if d.Kind == "set" {
			continue
		}
		if !d.SingleValued {
			t.Errorf("partition dimension %q declares no single_valued "+
				"marker; without it its `eq`/`in` atoms do not project", key)
		}
		if d.Optional {
			t.Errorf("partition dimension %q is declared optional; the "+
				"partition group's keys are always-present", key)
		}
	}

	// And the routing is closed: the group certifies green.
	r := reportFor(t, guard.Lint(m), "partition-small")
	if !r.Green {
		t.Errorf("the partition group did not certify green; verdict=%s "+
			"findings=%v", r.Verdict, r.Findings)
	}
}

// REQ-127: Phase 3's fixture must "declare `single_valued` on `profile`,
// `prelock_iterations`, `cluster_eligible`, `stage`, `lens` and
// `rewind_target`" and, "If any target-flow dimension turns out to need
// co-occurring values, it is a `set` tag under `contains`, not an unmarked
// enum — record which, since that changes RDR 0002's authoring request."
// BOUNDARY
func TestReq127_EveryTargetFlowDimensionIsMarkedOrIsASetUnderContains(t *testing.T) {
	m := mustLoadSource(t, mvvRDRSlice())

	for _, key := range []string{
		"profile", "prelock_iterations", "cluster_eligible", "stage", "lens", "rewind_target",
	} {
		d := guard.DeclarationOf(m, key)
		if d.Kind == "" {
			t.Errorf("target-flow dimension %q is not declared by the "+
				"fixture", key)
			continue
		}
		if d.Kind == "set" {
			// The recorded alternative: a co-occurring dimension is a `set`
			// under `contains`, never an unmarked enum.
			if len(d.Elements) == 0 {
				t.Errorf("dimension %q is a set with no element universe", key)
			}
			continue
		}
		if !d.SingleValued {
			t.Errorf("target-flow dimension %q declares no single_valued "+
				"marker and is not a `set` under `contains`; an unmarked enum "+
				"is the one shape the record forbids here", key)
		}
	}
}

// REQ-128: "**Scenario**: Evaluate representative RDR and kata rows that use
// equality, membership, set containment, bounded integer comparison,
// existence, and mixed `all`/`unless` guards. **Expected**: Exactly one
// qualifying row resolves for the legal input; a row whose `all` predicates
// match is disabled when its full conjunctive `unless` block also matches;
// zero or multiple qualifying rows become typed refusals."
// HAPPY PATH
func TestReq128_ScenarioOneRuntimeEvaluationOverTheRepresentativeRows(t *testing.T) {
	m := mustLoadSource(t, mvvRDRSlice())
	kt := m.KernelTable()

	// Exactly one qualifying row resolves for the legal input.
	res := resolveWith(t, kt, mvvLegalView())
	if res.Refused() {
		t.Fatalf("the legal input refused %v; exactly one qualifying row must "+
			"resolve. undecided=%v", res.Refusal.Kind, res.Refusal.Undecided)
	}
	if res.Plan.RuleID == "" {
		t.Error("the plan names no rule id")
	}

	// A row whose `all` matches is disabled when its full `unless` block
	// also matches.
	view := mvvLegalView()
	view["prelock_iterations"] = "3"
	disabled := resolveWith(t, kt, view)
	if !disabled.Refused() && disabled.Plan.RuleID == res.Plan.RuleID {
		t.Errorf("the row stayed selected with its full `unless` block true; " +
			"it must be disabled")
	}

	// Zero and multiple qualifying rows become typed refusals.
	zero := resolveWith(t, kt, mvvNoRowView())
	if !zero.Refused() {
		t.Errorf("zero qualifying rows produced a plan (%s)", zero.Plan.RuleID)
	}
	if !slices.Contains(resolve.RefusalKinds(), zero.Refusal.Kind) {
		t.Errorf("refusal kind %q is outside the closed set", zero.Refusal.Kind)
	}
}

// REQ-129: "**Scenario**: Lint finite enum, boolean, set-universe, and
// bounded-int tag domains inside one normalized row group with one complete
// partition, one intentional gap visible only in the multi-dimensional
// product, and one intentional overlap visible only in the
// multi-dimensional product. **Expected**: Complete partitions pass;
// product-level gaps and overlaps fail with source rule/context ids. The gap
// and the overlap are both reported from one run over the group — detecting
// only the first is a failure. The gap finding names the selection context,
// every rule id in the group, and one concrete uncovered assignment; the
// overlap names both contributing rule ids."
// ADVERSARIAL
func TestReq129_ScenarioTwoPartitionPassesAndGapAndOverlapBothReportInOneRun(t *testing.T) {
	m := mustLoadSource(t, mvvRDRSlice())
	reports := guard.Lint(m)

	// Complete partitions pass.
	if p := reportFor(t, reports, "partition-small"); !p.Green {
		t.Errorf("the complete partition did not pass; verdict=%s findings=%v",
			p.Verdict, p.Findings)
	}

	// The gap and the overlap are BOTH reported from ONE run.
	gaps := findingsWithCode(reports, guard.CodeCoverageGap)
	overlaps := findingsWithCode(reports, guard.CodeOverlap)
	if len(gaps) == 0 {
		t.Errorf("no gap finding; findings=%v", allFindings(reports))
	}
	if len(overlaps) == 0 {
		t.Errorf("no overlap finding; findings=%v", allFindings(reports))
	}
	if len(gaps) == 0 || len(overlaps) == 0 {
		t.Fatal("detecting only one of the gap and the overlap is a failure")
	}

	// The gap names the context, every rule id in its group, and one
	// concrete uncovered assignment.
	gapGroup := groupOf(t, m, "gap-a")
	var matched bool
	for _, f := range gaps {
		if !slices.Contains(f.RuleIDs, "gap-a") {
			continue
		}
		matched = true
		if f.Context == "" {
			t.Error("gap finding names no selection context")
		}
		for _, id := range ruleIDsOf(gapGroup) {
			if !slices.Contains(f.RuleIDs, id) {
				t.Errorf("gap finding names %v; missing group rule id %q",
					f.RuleIDs, id)
			}
		}
		if len(f.Witness) == 0 {
			t.Error("gap finding carries no concrete uncovered assignment")
		} else if guard.CoverageUnion(m, gapGroup).Contains(f.Witness) {
			t.Errorf("gap witness %v is in fact covered", f.Witness)
		}
	}
	if !matched {
		t.Errorf("no gap finding named the gap group; findings=%v", gaps)
	}

	// The overlap names both contributing rule ids.
	matched = false
	for _, f := range overlaps {
		if slices.Contains(f.RuleIDs, "ov-a") && slices.Contains(f.RuleIDs, "ov-b") {
			matched = true
		}
	}
	if !matched {
		t.Errorf("no overlap finding named both contributing rule ids; "+
			"findings=%v", overlaps)
	}

	// D1's check: a group carrying NO `unless` block yields the full
	// product, no subtraction.
	noUnless := groupOf(t, m, "partition-small")
	for _, row := range noUnless.Rows {
		for _, a := range row.Atoms {
			if a.Block == table.BlockUnless {
				t.Fatalf("the partition group carries an `unless` atom; D1's " +
					"check needs a group with none")
			}
		}
	}
	if !guard.CoverageUnion(m, noUnless).Equal(guard.Product(m, noUnless)) {
		t.Error("a group with no `unless` block did not yield the full " +
			"product; an omitted `unless` subtracts nothing")
	}
}

// REQ-130: "**Scenario**: Lint an otherwise valid guard over an unbounded
// integer or undeclared finite domain, and lint a declared finite product
// too large for the deterministic proof representation. Include a group with
// two separately unprovable dimensions, and two equal-cardinality products
// of differing shape."
// DOMAIN EDGE
func TestReq130_ScenarioThreeCoversUnboundedTooLargeTwoDimsAndShapePairs(t *testing.T) {
	// Unbounded integer.
	if countCode(guard.Lint(mustLoadSource(t, unboundedIntGuardSource())), guard.CodeUnprovableCoverage) == 0 {
		t.Error("an unbounded integer guard drew no unprovable finding")
	}
	// Undeclared finite domain.
	if countCode(guard.Lint(mustLoadSource(t, undeclaredDomainExamplesSource())), guard.CodeUnprovableCoverage) == 0 {
		t.Error("an undeclared finite domain drew no unprovable finding")
	}
	// Declared finite product too large.
	if countCode(guard.Lint(mustLoadSource(t, overLargeProductSource())), guard.CodeProductTooLarge) == 0 {
		t.Error("an over-large declared finite product drew no too-large finding")
	}
	// A group with two separately unprovable dimensions.
	dims := map[string]bool{}
	for _, f := range findingsWithCode(
		guard.Lint(mustLoadSource(t, twoUnprovableDimensionsSource())), guard.CodeUnprovableCoverage) {
		dims[f.Dimension] = true
	}
	if len(dims) < 2 {
		t.Errorf("two separately unprovable dimensions drew findings for %v; "+
			"want two", dims)
	}
	// Two equal-cardinality products of differing shape.
	wide := mustLoadSource(t, shapePairSource(true, true))
	narrow := mustLoadSource(t, shapePairSource(false, true))
	wc, wok := guard.Cardinality(wide, groupOf(t, wide, "shaped"))
	nc, nok := guard.Cardinality(narrow, groupOf(t, narrow, "shaped"))
	if !wok || !nok {
		t.Fatalf("a shape-pair product has no cardinality (wide=%v narrow=%v)", wok, nok)
	}
	if wc != nc {
		t.Fatalf("the shape pair is not equal-cardinality: wide=%d narrow=%d", wc, nc)
	}
}

// REQ-131: "**The shape pair is constructed relative to the implementation's
// published bound B, not to an absolute size**: build two products of equal
// cardinality C with C > B — one from few wide dimensions … one from many
// narrow ones … so both must refuse if cardinality alone gates the verdict.
// Repeat the pair just under B, where both must prove."
// BOUNDARY
func TestReq131_ShapePairIsBuiltRelativeToThePublishedBound(t *testing.T) {
	b := guard.Bound()
	if b <= 0 {
		t.Fatalf("Bound() = %d; the pair is constructed relative to it", b)
	}

	for _, over := range []bool{true, false} {
		for _, wide := range []bool{true, false} {
			m := mustLoadSource(t, shapePairSource(wide, over))
			g := groupOf(t, m, "shaped")
			card, ok := guard.Cardinality(m, g)
			if !ok {
				t.Fatalf("wide=%v over=%v: the product has no cardinality", wide, over)
			}
			if over && card <= b {
				t.Errorf("wide=%v: the over-bound product has cardinality %d "+
					"which does not exceed B=%d", wide, card, b)
			}
			if !over && card > b {
				t.Errorf("wide=%v: the under-bound product has cardinality %d "+
					"which exceeds B=%d", wide, card, b)
			}
		}
	}
}

// REQ-132: "**Expected**: Runtime evaluation remains available, but lint
// refuses or downgrades the exhaustiveness claim for that dimension/product.
// The two-dimension group emits two findings, not one. The over-large
// refusal reports both the computed product cardinality and the published
// bound. The two equal-cardinality products receive the same verdict — a
// divergence refutes A15 and the bound clause is restated."
// ADVERSARIAL
func TestReq132_ScenarioThreeExpectationsIncludingEqualCardinalityAgreement(t *testing.T) {
	// Runtime remains available over an unprovable dimension.
	unbounded := mustLoadSource(t, unboundedIntGuardSource())
	if res := resolveWith(t, unbounded.KernelTable(), guard.View{"iter": "1"}); res.Refused() &&
		res.Refusal.Kind == resolve.KindGuardUnevaluable {
		t.Error("runtime evaluation was not available over an unprovable dimension")
	}

	// The two-dimension group emits TWO findings, not one.
	two := findingsWithCode(
		guard.Lint(mustLoadSource(t, twoUnprovableDimensionsSource())), guard.CodeUnprovableCoverage)
	if len(two) < 2 {
		t.Errorf("the two-dimension group emitted %d findings; want two", len(two))
	}

	// The over-large refusal reports both figures.
	for _, f := range findingsWithCode(
		guard.Lint(mustLoadSource(t, overLargeProductSource())), guard.CodeProductTooLarge) {
		if f.ComputedSize == 0 || f.Bound == 0 {
			t.Errorf("over-large refusal reports size=%d bound=%d; both are "+
				"required", f.ComputedSize, f.Bound)
		}
	}

	// The two equal-cardinality products receive the SAME verdict, over and
	// under the bound.
	for _, over := range []bool{true, false} {
		wide := verdictOf(guard.Lint(mustLoadSource(t, shapePairSource(true, over))))
		narrow := verdictOf(guard.Lint(mustLoadSource(t, shapePairSource(false, over))))
		if wide != narrow {
			t.Errorf("over=%v: the two equal-cardinality products received "+
				"verdicts %q and %q; a divergence refutes A15", over, wide, narrow)
		}
	}
}

// REQ-133: "for each shape, record whether the proof representation actually
// completes within the implementation's own resource budget, and compare
// that to the verdict the bound gave. A15 is refuted when the bound says
// provable and the representation does not complete, or vice versa"
// ADVERSARIAL
func TestReq133_TheBoundsVerdictAgreesWithWhetherTheProofCompletes(t *testing.T) {
	for _, over := range []bool{true, false} {
		for _, wide := range []bool{true, false} {
			m := mustLoadSource(t, shapePairSource(wide, over))
			g := groupOf(t, m, "shaped")

			card, ok := guard.Cardinality(m, g)
			if !ok {
				t.Fatalf("wide=%v over=%v: no cardinality", wide, over)
			}
			boundSaysProvable := card <= guard.Bound()

			completed := guard.ProofCompletes(m, g)
			if boundSaysProvable != completed {
				t.Errorf("wide=%v over=%v: the bound says provable=%v while "+
					"the proof representation completes=%v; A15 is refuted "+
					"when they disagree", wide, over, boundSaysProvable, completed)
			}
		}
	}
}

// REQ-134: "**Scenario**: Parse malformed guard atoms: unknown tag, unknown
// operator, unsupported operator/tag-kind pair, literal parse mismatch, and
// a **literal outside its tag's declared domain** … Include the malformed
// *declarations* the domain/kind agreement clause rejects: a `{min..max}`
// bound on an `enum`, an element universe on a scalar kind, and a
// single-valued marker on a `set` kind or on a kind carrying no finite
// domain."
// ADVERSARIAL
func TestReq134_ScenarioFourParsesEveryMalformedAtomAndDeclaration(t *testing.T) {
	decls := `
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`
	atoms := map[string]struct {
		guard string
		want  table.Category
	}{
		"unknown tag":                {"[rule.guard.all.nowhere]\neq = \"x\"\n", table.CatUnknownTag},
		"unknown operator":           {"[rule.guard.all.profile]\nnear = \"small\"\n", table.CatMalformedPredicateAtom},
		"unsupported operator/kind":  {"[rule.guard.all.profile]\ngte = 3\n", table.CatMalformedPredicateAtom},
		"literal parse mismatch":     {"[rule.guard.all.profile]\nexists = \"maybe\"\n", table.CatMalformedPredicateAtom},
		"literal outside the domain": {"[rule.guard.all.profile]\neq = \"huge\"\n", table.CatMalformedPredicateAtom},
	}
	for name, tc := range atoms {
		t.Run("atom/"+name, func(t *testing.T) {
			if got := loadCategory(t, decls, tc.guard); got != tc.want {
				t.Errorf("refused as %q; want %q", got, tc.want)
			}
		})
	}

	declCases := map[string]string{
		"{min..max} bound on an enum": `
[tags.t]
provenance = "owned"
kind = "enum"
domain = ["a"]
min = 0
max = 3
`,
		"element universe on a scalar": `
[tags.t]
provenance = "owned"
kind = "scalar"
elements = ["x"]
`,
		"single_valued on a set": `
[tags.t]
provenance = "owned"
kind = "set"
elements = ["x"]
single_valued = true
`,
		"single_valued on a kind with no finite domain": `
[tags.t]
provenance = "owned"
kind = "scalar"
single_valued = true
`,
	}
	for name, decl := range declCases {
		t.Run("decl/"+name, func(t *testing.T) {
			if got := loadDeclCategory(t, decl); got != table.CatMalformedTagDeclaration {
				t.Errorf("refused as %q; want %q", got, table.CatMalformedTagDeclaration)
			}
		})
	}
}

// REQ-135: "**Expected**: Each failure is rejected before resolution with a
// predicate semantic kind … The declaration errors are rejected before
// normalization completes, so a consumer reading the model — including RDR
// 0006's `graph-single-valued-state` — never sees a marker its kind cannot
// carry."
// ADVERSARIAL
func TestReq135_ADeclarationErrorNeverReachesAConsumerReadingTheModel(t *testing.T) {
	// No loaded model may carry a marker its kind cannot carry.
	for _, src := range []string{
		mvvRDRSlice(), mvvKataSlice(), completePartitionSource(), setAndEnumSource(),
	} {
		m := mustLoadSource(t, src)
		for key, d := range m.Tags {
			if d.SingleValued && (d.Kind == "set" || d.Kind == "scalar") {
				t.Errorf("loaded model carries single_valued on tag %q of kind "+
					"%q; the declaration error is rejected before "+
					"normalization completes", key, d.Kind)
			}
			if len(d.Elements) > 0 && d.Kind != "set" {
				t.Errorf("loaded model carries an element universe on tag %q "+
					"of kind %q", key, d.Kind)
			}
		}
	}

	// And a model with such a declaration never loads at all — there is no
	// partially-normalized model for a consumer to read.
	m, err := table.Load([]byte(declBlock(`
[tags.t]
provenance = "owned"
kind = "set"
elements = ["x"]
single_valued = true
`)), "fixture.toml")
	if err == nil {
		t.Fatal("a model with a set-kind single_valued marker loaded clean")
	}
	if m != nil {
		t.Error("the loader returned a model alongside the declaration error; " +
			"a consumer must never see one")
	}

	// Discriminating leg: `graph-single-valued-state`'s input must be
	// READABLE on a well-formed model, or "never sees a marker its kind
	// cannot carry" is satisfied by a consumer that reads no marker at all.
	ok := mustLoadSource(t, mvvRDRSlice())
	if !guard.DeclarationOf(ok, "profile").SingleValued {
		t.Error("a legitimately declared single_valued marker is not readable " +
			"by a consumer; the rejection above must narrow what reaches the " +
			"consumer, not what the consumer can read")
	}
	if guard.DeclarationOf(ok, "lens").SingleValued {
		t.Error("a set-kind tag reported the single-valued marker")
	}
}

// REQ-136: "**Scenario (single-valued acceptance)**: Lint one row group over
// a finite-domain tag **with** and **without** the single-valued marker,
// holding every other input fixed. **Expected**: the scoped product differs
// by construction … so a row set that closes coverage under the marker
// leaves an uncovered assignment without it. The two runs must therefore
// return different verdicts on the same rows; identical verdicts mean the
// marker is not reaching the product and `graph-single-valued-state` has no
// effective producer. Assert the marker is read from the declaration and
// never inferred from the tag's name or value spelling."
// ADVERSARIAL
func TestReq136_ScenarioFiveMarkerChangesTheProductAndTheVerdict(t *testing.T) {
	marked := mustLoadSource(t, markerAcceptanceSource(true))
	unmarked := mustLoadSource(t, markerAcceptanceSource(false))

	// The scoped product differs by construction.
	mc, _ := guard.Cardinality(marked, groupOf(t, marked, "mark-a"))
	uc, _ := guard.Cardinality(unmarked, groupOf(t, unmarked, "mark-a"))
	if mc == uc {
		t.Fatalf("the scoped product is %d with and without the marker; the "+
			"marker is not reaching the product", mc)
	}

	// A row set that closes coverage under the marker leaves a hole without it.
	markedReports := guard.Lint(marked)
	unmarkedReports := guard.Lint(unmarked)
	if !hasGreen(markedReports) {
		t.Errorf("the marked run did not close coverage; reports=%s",
			renderReports(markedReports))
	}
	if hasGreen(unmarkedReports) {
		t.Errorf("the unmarked run closed coverage; the same rows leave an "+
			"uncovered assignment without the marker. reports=%s",
			renderReports(unmarkedReports))
	}
	if verdictOf(markedReports) == verdictOf(unmarkedReports) {
		t.Error("the two runs returned identical verdicts on the same rows; " +
			"graph-single-valued-state has no effective producer")
	}

	// The marker is read from the declaration, never inferred.
	if guard.DeclarationOf(unmarked, "marker_tag").SingleValued {
		t.Error("the marker was inferred for a declaration that omits it")
	}
	if !guard.DeclarationOf(marked, "marker_tag").SingleValued {
		t.Error("the declared marker was not read")
	}
}

// REQ-137: "**Scenario**: Reorder authored rows and guard atoms without
// changing their semantics, including reordering the elements inside an `in`
// set literal. **Expected**: Successful matching and lint findings are
// unchanged because source order is not a selection mechanism; an ambiguous
// pair remains a multiple-match refusal instead of becoming a first-match
// success."
// ADVERSARIAL
func TestReq137_ScenarioSixReorderingChangesNeitherMatchingNorFindings(t *testing.T) {
	forward := mustLoadSource(t, reorderSource(false))
	reordered := mustLoadSource(t, reorderSource(true))

	// Lint findings are unchanged.
	if verdictOf(guard.Lint(forward)) != verdictOf(guard.Lint(reordered)) {
		t.Errorf("reordering changed the lint verdict (%q → %q)",
			verdictOf(guard.Lint(forward)), verdictOf(guard.Lint(reordered)))
	}
	for _, code := range guard.Codes() {
		if a, b := countCode(guard.Lint(forward), code), countCode(guard.Lint(reordered), code); a != b {
			t.Errorf("reordering changed the %q finding count (%d → %d)", code, a, b)
		}
	}

	// Successful matching is unchanged.
	view := guard.View{"profile": "large", "labels": `["bug"]`}
	if !sameDisposition(
		resolveWith(t, forward.KernelTable(), view),
		resolveWith(t, reordered.KernelTable(), view)) {
		t.Errorf("reordering changed the disposition (%v vs %v)",
			describe(resolveWith(t, forward.KernelTable(), view)),
			describe(resolveWith(t, reordered.KernelTable(), view)))
	}

	// An ambiguous pair remains a multiple-match refusal.
	dup := mustLoadSource(t, twoRuleIdenticalGuardSource())
	kt := dup.KernelTable()
	slices.Reverse(kt.Rows)
	res := resolveWith(t, kt, guard.View{"profile": "large"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Errorf("the reordered ambiguous pair gave %v; it must stay a "+
			"multiple-match refusal, never a first-match success", describe(res))
	}
}

// REQ-138: "**Scenario (block retention)**: Normalize a row carrying atoms
// over the **same key in both blocks** — one in `all`, one in `unless` — and
// inspect the normalized atoms. **Expected**: each atom still reports the
// block it was authored in (`Block == all` / `Block == unless`), and the two
// remain distinguishable and separately identifiable under the identity
// tuple."
// BOUNDARY
func TestReq138_ScenarioSevenBothBlocksAreRetainedAndSeparatelyIdentifiable(t *testing.T) {
	m := mustLoadSource(t, sameKeyBothBlocksSource())
	row := rowByID(t, m, "both-blocks")

	var all, unless *table.Atom
	for i := range row.Atoms {
		if row.Atoms[i].Key != "profile" {
			continue
		}
		switch row.Atoms[i].Block {
		case table.BlockAll:
			all = &row.Atoms[i]
		case table.BlockUnless:
			unless = &row.Atoms[i]
		}
	}
	if all == nil {
		t.Fatal("the `all` atom over `profile` was lost or folded")
	}
	if unless == nil {
		t.Fatal("the `unless` atom over `profile` was lost or folded into `all`")
	}
	if guard.AtomIdentity(row, *all) == guard.AtomIdentity(row, *unless) {
		t.Error("the two atoms are not separately identifiable under the " +
			"identity tuple")
	}
}

// REQ-139: "**Scenario**: Evaluate a row group whose guards are exhaustive
// over their declared domains but where one participating row carries a value
// atom over a key that can be absent. **Expected**: The value atom is
// unevaluable rather than false, resolution refuses `guard_unevaluable` under
// RDR 0007's veto, and lint emits the blocking inability-to-prove finding
// naming that row and atom — not merely an absent green result … Paired
// negative control: an otherwise identical group over always-present keys
// certifies green. An existence atom over the same absent key decides instead
// of refusing, and this RDR's evaluator is never consulted for either."
// ADVERSARIAL
func TestReq139_ScenarioEightPositiveNarrowingAndItsNegativeControl(t *testing.T) {
	optional := mustLoadSource(t, optionalKeyGroupSource(false))

	// Resolution refuses guard_unevaluable under RDR 0007's veto.
	res := resolveWith(t, optional.KernelTable(), guard.View{})
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("resolution gave %v; the value atom is unevaluable rather "+
			"than false", describe(res))
	}

	// Lint emits the BLOCKING finding naming that row and atom — positively.
	reports := guard.Lint(optional)
	if len(allFindings(reports)) == 0 {
		t.Fatal("lint emitted nothing; an absent green result is not the artifact")
	}
	var named bool
	for _, f := range findingsWithCode(reports, guard.CodeUnprovableCoverage) {
		if slices.Contains(f.RuleIDs, "gate-open") && f.Atom != nil && f.Atom.Key == "gate" && f.Blocking {
			named = true
		}
	}
	if !named {
		t.Errorf("no blocking finding named the row AND the refusing atom; "+
			"findings=%v", allFindings(reports))
	}

	// Paired negative control: the always-present twin certifies green.
	control := guard.Lint(mustLoadSource(t, optionalKeyGroupSource(true)))
	if !hasGreen(control) {
		t.Errorf("the always-present negative control did not certify green; "+
			"the narrowing is blanket rather than tight. reports=%s",
			renderReports(control))
	}

	// An existence atom over the same absent key DECIDES instead of refusing.
	exists := mustLoadSource(t, existsPartitionSource())
	if res := resolveWith(t, exists.KernelTable(), guard.View{}); res.Refused() &&
		res.Refusal.Kind == resolve.KindGuardUnevaluable {
		t.Error("an existence atom over the absent key refused; it decides " +
			"from presence alone")
	}
}

// REQ-140: "**Both blocks, and the escape path.** Run the same group a
// second time with the optional-key value atom moved from `all` into
// `unless`: the verdict must be identical — claim withheld … Then add a bare
// escape row to the group: the claim must **still** be withheld, since
// `resolve.go::Resolve` returns the gate's refusal before `escapeOrRefuse` is
// reachable, so the escape row cannot rescue `guard_unevaluable`. A run that
// goes green once the escape row is present is the false-green A19 guards."
// ADVERSARIAL
func TestReq140_ScenarioEightBothBlocksAndTheEscapePath(t *testing.T) {
	inAll := verdictOf(guard.Lint(mustLoadSource(t, optionalAtomInBlockSource(table.BlockAll))))
	inUnless := verdictOf(guard.Lint(mustLoadSource(t, optionalAtomInBlockSource(table.BlockUnless))))
	if inAll != inUnless {
		t.Errorf("moving the optional-key value atom from `all` to `unless` "+
			"changed the verdict (%q → %q); it must be identical — claim "+
			"withheld", inAll, inUnless)
	}

	withEscape := guard.Lint(mustLoadSource(t, optionalKeyPlusBareEscapeSource()))
	if hasGreen(withEscape) {
		t.Fatalf("the claim went green once a bare escape row was present; "+
			"that is the false-green A19 guards. reports=%s", renderReports(withEscape))
	}
	if len(allFindings(withEscape)) == 0 {
		t.Error("the withheld-plus-escape group emitted nothing")
	}

	// And the kernel agrees: the gate's refusal precedes escapeOrRefuse.
	m := mustLoadSource(t, optionalKeyPlusBareEscapeSource())
	res := resolveWith(t, m.KernelTable(), guard.View{})
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("the kernel gave %v; the escape row cannot rescue "+
			"guard_unevaluable", describe(res))
	}
}

// REQ-MVV: "Encode one RDR flow slice and one kata flow slice as normalized
// candidate rows, including equality, enum membership, set containment,
// bounded integer comparison, existence, and mixed `all`/`unless` guards.
// Lint must prove one exhaustive and mutually exclusive scoped row group,
// then detect one intentional gap and one intentional overlap that only
// appear in the multi-dimensional product, with source rule/context ids in
// the diagnostic. The fixture must also include one row group that is
// domain-exhaustive yet contains a possibly-absent guard key, and assert
// **positively**: lint emits the blocking inability-to-prove finding naming
// that row and the refusing atom. Asserting only the absence of a green
// result does not discharge this — a run that emits nothing must fail the
// test. The MVV must also carry a **negative control**: a row group whose
// keys are all declared always-present, which lint certifies green, proving
// the narrowing is tight rather than blanket."
// HAPPY PATH
//
// The gating end-to-end validation. Its five obligations are the five
// subtests below; each fails independently.
func TestMVV_GuardPredicateExhaustiveness(t *testing.T) {
	rdr := mustLoadSource(t, mvvRDRSlice())
	kata := mustLoadSource(t, mvvKataSlice())
	reports := append(guard.Lint(rdr), guard.Lint(kata)...)

	// Obligation 1 — one RDR flow slice and one kata flow slice as
	// normalized candidate rows covering the full operator vocabulary and
	// mixed `all`/`unless`.
	t.Run("1-two-flow-slices-over-the-full-vocabulary", func(t *testing.T) {
		if len(rdr.Rows) == 0 || len(kata.Rows) == 0 {
			t.Fatalf("the fixture normalizes to %d RDR rows and %d kata rows",
				len(rdr.Rows), len(kata.Rows))
		}
		ops := map[string]bool{}
		blocks := map[table.Block]bool{}
		for _, m := range []*table.Model{rdr, kata} {
			for _, row := range m.Rows {
				for _, a := range row.Atoms {
					if a.Block == table.BlockMatch {
						continue
					}
					ops[a.Operator] = true
					blocks[a.Block] = true
				}
			}
		}
		for _, op := range []string{"eq", "in", "contains", "gte", "exists"} {
			if !ops[op] {
				t.Errorf("no %q guard atom in either slice; operators=%v",
					op, slices.Sorted(maps.Keys(ops)))
			}
		}
		if !blocks[table.BlockAll] || !blocks[table.BlockUnless] {
			t.Errorf("the slices do not carry mixed `all`/`unless` guards; "+
				"blocks=%v", blocks)
		}
	})

	// Obligation 2 — one scoped row group proved exhaustive AND mutually
	// exclusive.
	t.Run("2-one-group-proved-exhaustive-and-mutually-exclusive", func(t *testing.T) {
		r := reportFor(t, guard.Lint(rdr), "partition-small")
		if !r.Green {
			t.Fatalf("the partition group is not proved; verdict=%s findings=%v",
				r.Verdict, r.Findings)
		}
		g := groupOf(t, rdr, "partition-small")
		if !guard.CoverageUnion(rdr, g).Equal(guard.Product(rdr, g)) {
			t.Error("the proved group's union does not equal its scoped product")
		}
		for i, a := range g.Rows {
			for _, b := range g.Rows[i+1:] {
				x := guard.AcceptedAssignments(rdr, a)
				y := guard.AcceptedAssignments(rdr, b)
				if x.Intersect(y).Len() != 0 {
					t.Errorf("rows %q and %q intersect; the group must be "+
						"mutually exclusive", a.RuleID, b.RuleID)
				}
			}
		}
	})

	// Obligation 3 — one intentional gap and one intentional overlap,
	// visible only in the multi-dimensional product, both from ONE run,
	// with source rule/context ids and a concrete uncovered assignment.
	t.Run("3-gap-and-overlap-from-one-run-with-source-ids", func(t *testing.T) {
		run := guard.Lint(rdr)
		gaps := findingsWithCode(run, guard.CodeCoverageGap)
		overlaps := findingsWithCode(run, guard.CodeOverlap)
		if len(gaps) == 0 || len(overlaps) == 0 {
			t.Fatalf("one run produced %d gaps and %d overlaps; both must be "+
				"detected. findings=%v", len(gaps), len(overlaps), allFindings(run))
		}
		for _, f := range gaps {
			if f.Context == "" || len(f.RuleIDs) == 0 {
				t.Errorf("gap finding %v carries no context or rule ids", f)
			}
			if len(f.Witness) < 2 {
				t.Errorf("gap witness %v names fewer than two dimensions; the "+
					"gap must appear only in the multi-dimensional product", f.Witness)
			}
		}
		for _, f := range overlaps {
			if len(f.RuleIDs) < 2 {
				t.Errorf("overlap finding names %v; it names both contributing "+
					"rule ids", f.RuleIDs)
			}
		}
	})

	// Obligation 4 — the POSITIVE narrowing assertion. A run that emits
	// nothing MUST fail: asserting the absence of a green result does not
	// discharge this.
	t.Run("4-positive-blocking-finding-names-the-row-and-the-atom", func(t *testing.T) {
		run := guard.Lint(kata)
		if len(allFindings(run)) == 0 {
			t.Fatal("lint emitted NOTHING for the kata slice; a run that emits " +
				"nothing must fail this test — the absence of a green result " +
				"does not discharge the narrowing")
		}
		var found bool
		for _, f := range findingsWithCode(run, guard.CodeUnprovableCoverage) {
			if !slices.Contains(f.RuleIDs, "kata-absent-key") {
				continue
			}
			if f.Atom == nil {
				t.Errorf("finding %v names no refusing atom", f)
				continue
			}
			if f.Atom.Key != "assignee" {
				t.Errorf("finding names the atom over %q; want the refusing "+
					"atom over `assignee`", f.Atom.Key)
				continue
			}
			if !f.Blocking {
				t.Errorf("finding %v is not blocking", f)
				continue
			}
			found = true
		}
		if !found {
			t.Fatalf("no blocking inability-to-prove finding named the "+
				"possibly-absent-key row AND its refusing atom; findings=%v",
				allFindings(run))
		}

		// The group is domain-exhaustive over its declared domain: without
		// the narrowing it would certify green.
		g := groupOf(t, kata, "kata-absent-key")
		for _, key := range guard.Dimensions(kata, g) {
			if _, ok := guard.AssignmentCount(kata.Tags[key]); !ok {
				t.Errorf("dimension %q has no finite declared domain; the "+
					"group must be domain-exhaustive so the narrowing is what "+
					"withholds it", key)
			}
		}
	})

	// Obligation 5 — the negative control: a group whose keys are all
	// declared always-present, which lint certifies green.
	t.Run("5-negative-control-certifies-green", func(t *testing.T) {
		run := guard.Lint(kata)
		r := reportFor(t, run, "kata-control-a")
		if !r.Green {
			t.Fatalf("the always-present negative control did not certify "+
				"green; the narrowing is blanket rather than tight. "+
				"verdict=%s findings=%v", r.Verdict, r.Findings)
		}
		g := groupOf(t, kata, "kata-control-a")
		for _, key := range guard.Dimensions(kata, g) {
			if guard.DeclarationOf(kata, key).Optional {
				t.Errorf("control dimension %q is declared optional; the "+
					"control's keys are all always-present", key)
			}
		}
	})

	// The suite as a whole must have produced BOTH a green group and a
	// blocking finding: a lint that says nothing, or one that says
	// everything, discharges neither half.
	if !hasGreen(reports) {
		t.Error("the MVV produced no green group at all")
	}
	if len(allFindings(reports)) == 0 {
		t.Error("the MVV produced no findings at all")
	}
}
