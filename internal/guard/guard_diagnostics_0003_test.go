package guard_test

// RDR 0003 — report-every-defect, diagnostic attribution, provenance and
// owned tags, and the remaining disposition rows.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-95: "Lint MUST report every defect it can decide in one pass over a
// row group, not the first one it encounters. Withholding a group's
// exhaustiveness claim MUST NOT suppress overlap findings, coverage
// findings, or further withholding reasons for that same group: each
// unprovable dimension, each refusing row, each overlapping pair, and any
// coverage gap over a provable product is its own finding."
// ADVERSARIAL
func TestReq95_EveryDecidableDefectIsReportedInOnePass(t *testing.T) {
	// A group with a gap AND an overlap, both over a provable product.
	m := mustLoadSource(t, gapAndOverlapSource())
	reports := guard.Lint(m)

	if countCode(reports, guard.CodeCoverageGap) == 0 {
		t.Errorf("the coverage gap was not reported; findings=%v", allFindings(reports))
	}
	if countCode(reports, guard.CodeOverlap) == 0 {
		t.Errorf("the overlap was not reported; findings=%v", allFindings(reports))
	}

	// Two unprovable dimensions in one group draw one finding EACH.
	two := guard.Lint(mustLoadSource(t, twoUnprovableDimensionsSource()))
	dims := map[string]bool{}
	for _, f := range findingsWithCode(two, guard.CodeUnprovableCoverage) {
		dims[f.Dimension] = true
	}
	if !dims["loose_a"] || !dims["loose_b"] {
		t.Errorf("two unprovable dimensions drew findings for %v; each is its "+
			"own finding, not the first one encountered", dims)
	}
}

// REQ-96: "A group with two refusing rows MUST emit a finding for each, so
// the emitted set does not depend on row or dimension iteration order — the
// same source-order independence matching already requires."
// ADVERSARIAL
func TestReq96_TwoRefusingRowsEmitTwoFindingsInAnyIterationOrder(t *testing.T) {
	forward := guard.Lint(mustLoadSource(t, twoRefusingRowsSource()))
	reversed := guard.Lint(mustLoadSource(t, reversedTwoRefusingRowsSource()))

	namedIn := func(reports []guard.GroupReport) []string {
		seen := map[string]bool{}
		for _, f := range findingsWithCode(reports, guard.CodeUnprovableCoverage) {
			for _, id := range f.RuleIDs {
				seen[id] = true
			}
		}
		out := make([]string, 0, len(seen))
		for id := range seen {
			out = append(out, id)
		}
		slices.Sort(out)
		return out
	}

	a, b := namedIn(forward), namedIn(reversed)
	if !slices.Equal(a, b) {
		t.Errorf("reversing row order changed the emitted set (%v vs %v); it "+
			"does not depend on row iteration order", a, b)
	}
	if len(a) != 2 {
		t.Errorf("the emitted set names %v; a group with two refusing rows "+
			"MUST emit a finding for each", a)
	}
}

// REQ-97: "A withheld claim and a coverage or overlap finding MAY be
// emitted together; what MUST NOT be emitted is a green exhaustiveness
// result alongside any withholding reason."
// ADVERSARIAL
func TestReq97_NoGreenResultStandsAlongsideAWithholdingReason(t *testing.T) {
	m := mustLoadSource(t, withheldPlusOverlapSource())
	reports := guard.Lint(m)

	// Co-emission is permitted.
	if countCode(reports, guard.CodeOverlap) == 0 ||
		countCode(reports, guard.CodeUnprovableCoverage) == 0 {
		t.Errorf("a withheld claim and an overlap finding were not emitted "+
			"together; findings=%v", allFindings(reports))
	}

	// A green result alongside a withholding reason is forbidden.
	for _, r := range reports {
		if !r.Green {
			continue
		}
		for _, f := range r.Findings {
			if f.Code == guard.CodeUnprovableCoverage {
				t.Errorf("group %s certified green alongside withholding "+
					"reason %v", r.Context, f)
			}
		}
	}
}

// REQ-98: "Overlap and coverage diagnostics MUST name the source rule id or
// context id that contributed each predicate involved in the finding."
// HAPPY PATH
func TestReq98_OverlapAndCoverageDiagnosticsNameEveryContributingSourceID(t *testing.T) {
	m := mustLoadSource(t, gapAndOverlapSource())
	reports := guard.Lint(m)

	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		if len(f.RuleIDs) < 2 {
			t.Errorf("overlap finding names %v; it MUST name the source id "+
				"contributing EACH predicate involved", f.RuleIDs)
		}
		for _, id := range f.RuleIDs {
			if id == "" {
				t.Errorf("overlap finding %v carries an empty source id", f.RuleIDs)
			}
		}
	}
	for _, f := range findingsWithCode(reports, guard.CodeCoverageGap) {
		if len(f.RuleIDs) == 0 {
			t.Errorf("coverage finding names no source id; findings=%v", f)
		}
	}

	// Discriminating leg: the diagnostics must EXIST to be attributed. A
	// lint that emits nothing names every contributing id vacuously.
	if countCode(reports, guard.CodeOverlap) == 0 || countCode(reports, guard.CodeCoverageGap) == 0 {
		t.Fatalf("no overlap or coverage diagnostic to attribute; findings=%v",
			allFindings(reports))
	}
}

// REQ-99: "A coverage gap has no contributing predicate — it is an absence
// — so it is attributed differently: a gap finding MUST name the selection
// context that scopes the group, every source rule id in that group, and at
// least one concrete uncovered assignment from the scoped product … Naming
// the group without a witness assignment MUST NOT satisfy this clause."
// ADVERSARIAL
func TestReq99_GapFindingNamesContextEveryRuleIDAndAWitnessAssignment(t *testing.T) {
	m := mustLoadSource(t, productOnlyGapSource())
	reports := guard.Lint(m)

	gaps := findingsWithCode(reports, guard.CodeCoverageGap)
	if len(gaps) == 0 {
		t.Fatalf("no gap finding; findings=%v", allFindings(reports))
	}
	g := groupOf(t, m, "gap-a")
	for _, f := range gaps {
		if f.Context == "" {
			t.Error("gap finding names no selection context")
		}
		for _, id := range ruleIDsOf(g) {
			if !slices.Contains(f.RuleIDs, id) {
				t.Errorf("gap finding names %v; it MUST name EVERY source rule "+
					"id in the group (missing %q)", f.RuleIDs, id)
			}
		}
		if len(f.Witness) == 0 {
			t.Fatal("gap finding carries no witness assignment; naming the " +
				"group without one MUST NOT satisfy this clause")
		}
		// The witness must be CONCRETE and genuinely uncovered.
		if guard.CoverageUnion(m, g).Contains(f.Witness) {
			t.Errorf("witness %v is covered by the union; it must be an "+
				"UNCOVERED assignment from the scoped product", f.Witness)
		}
		if !guard.Product(m, g).Contains(f.Witness) {
			t.Errorf("witness %v is not in the scoped product", f.Witness)
		}
	}
}

// REQ-100: "Predicate lint MUST distinguish owned, observed, and recognized
// tags. A row that matches an owned tag MUST be rejected unless every
// reachable predecessor sets or preserves that tag before the match."
// DOMAIN EDGE
func TestReq100_OwnedTagMatchedWithNoReachablePredecessorWriteIsRejected(t *testing.T) {
	m := mustLoadSource(t, ownedBeforeWriteSource())
	reports := guard.Lint(m)

	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOwnedBeforeWrite &&
			f.Dimension == "never_written" &&
			slices.Contains(f.RuleIDs, "reads-unwritten")
	}) {
		t.Errorf("a row matching an owned tag no reachable predecessor writes "+
			"drew no finding; findings=%v", allFindings(reports))
	}

	// The three provenances are distinguished: an OBSERVED tag matched the
	// same way draws no such finding.
	observed := guard.Lint(mustLoadSource(t, observedMatchSource()))
	if anyFinding(observed, func(f guard.Finding) bool {
		return f.Code == guard.CodeOwnedBeforeWrite
	}) {
		t.Error("matching an OBSERVED tag drew the owned-before-write " +
			"finding; lint MUST distinguish the three provenances")
	}
}

// REQ-101: "The reachable-predecessor relation is RDR 0006's owned-state
// reachability relation … this RDR cites it and defines no second one
// (A12)."
// BOUNDARY
func TestReq101_ThisRDRDefinesNoSecondReachabilityRelation(t *testing.T) {
	// The withholding decision is syntactic over ONE declared field and
	// consults no reachability relation. If this RDR defined a second one,
	// the two decisions would have to disagree somewhere: a key that is
	// written by every predecessor but declared optional must STILL be
	// can-refuse.
	m := mustLoadSource(t, optionalButAlwaysWrittenSource())

	decls := m.Tags
	row := rowByID(t, m, "reader")
	if !guard.CanRefuse(decls, row) {
		t.Error("a key written by every reachable predecessor but DECLARED " +
			"optional was read as non-refusing; the withholding decision " +
			"reads the declaration, not a reachability relation")
	}
	if hasGreen(guard.Lint(m)) {
		t.Error("the group certified green; the withholding decision does not " +
			"consult reachability")
	}
}

// REQ-102: "recognized tags are fresh event inputs, observed tags are
// re-read before matching, and owned tags must have a reachable predecessor
// write before a row may match them."
// BOUNDARY
func TestReq102_TheThreeProvenancesCarryTheirThreeStatedRoles(t *testing.T) {
	m := mustLoadSource(t, ownedBeforeWriteSource())

	for key, want := range map[string]table.Provenance{
		"recognized":    table.ProvenanceRecognized,
		"never_written": table.ProvenanceOwned,
	} {
		if got := guard.DeclarationOf(m, key).Provenance; got != want {
			t.Errorf("tag %q carries provenance %q; want %q", key, got, want)
		}
	}

	// Only the OWNED provenance carries the predecessor-write obligation.
	for _, prov := range []table.Provenance{
		table.ProvenanceObserved, table.ProvenanceRecognized,
	} {
		if guard.RequiresPredecessorWrite(prov) {
			t.Errorf("provenance %q carries the predecessor-write obligation; "+
				"only owned tags do", prov)
		}
	}
	if !guard.RequiresPredecessorWrite(table.ProvenanceOwned) {
		t.Error("owned tags do not carry the predecessor-write obligation")
	}
}

// REQ-106: "Zero rows qualify" → lint "Coverage gap if the product is
// provable", `graph-coverage-gap` "naming the selection context, every rule
// id in the group, and one uncovered assignment". "Two or more rows
// qualify" → `graph-overlap` "naming both source rule ids"; "RDR 0001
// refuses — never first-match".
// DOMAIN EDGE
func TestReq106_ZeroQualifyingIsAGapAndTwoQualifyingIsAnOverlap(t *testing.T) {
	m := mustLoadSource(t, gapAndOverlapSource())
	reports := guard.Lint(m)

	gaps := findingsWithCode(reports, guard.CodeCoverageGap)
	if len(gaps) == 0 {
		t.Fatalf("no coverage gap for the uncovered assignment; findings=%v",
			allFindings(reports))
	}
	for _, f := range gaps {
		if f.Context == "" || len(f.RuleIDs) == 0 || len(f.Witness) == 0 {
			t.Errorf("gap finding %v is missing the context, rule ids, or the "+
				"uncovered assignment", f)
		}
	}

	overlaps := findingsWithCode(reports, guard.CodeOverlap)
	if len(overlaps) == 0 {
		t.Fatalf("no overlap finding; findings=%v", allFindings(reports))
	}
	for _, f := range overlaps {
		if len(f.RuleIDs) < 2 {
			t.Errorf("overlap finding names %v; it names BOTH source rule ids",
				f.RuleIDs)
		}
	}

	// And the runtime refuses rather than first-matching.
	res := resolveWith(t, m.KernelTable(), guard.View{"profile": "small", "flag": "true"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Errorf("two qualifying rows gave %v; RDR 0001 refuses — never "+
			"first-match", describe(res))
	}
}

// REQ-107: "Guard dimension lacks a finite declared domain" → "Evaluable at
// runtime"; lint "**Refuse or downgrade** the group's claim, one finding per
// unprovable dimension — never treat examples as complete, never stop at the
// first"; `graph-unprovable-coverage` naming the dimension.
// DOMAIN EDGE
func TestReq107_UnprovableDimensionStaysRuntimeEvaluableAndDrawsOneFindingEach(t *testing.T) {
	m := mustLoadSource(t, twoUnprovableDimensionsSource())

	// Evaluable at runtime.
	res := resolveWith(t, m.KernelTable(), guard.View{"loose_a": "one", "loose_b": "one"})
	if res.Refused() && res.Refusal.Kind == resolve.KindGuardUnevaluable {
		t.Error("the runtime could not evaluate a guard over an unprovable " +
			"dimension; it stays evaluable at runtime")
	}

	// One finding per unprovable dimension, never stopping at the first.
	dims := map[string]bool{}
	for _, f := range findingsWithCode(guard.Lint(m), guard.CodeUnprovableCoverage) {
		if f.Dimension == "" {
			t.Errorf("finding %v names no dimension", f)
		}
		dims[f.Dimension] = true
	}
	if !dims["loose_a"] || !dims["loose_b"] {
		t.Errorf("findings name dimensions %v; want one per unprovable "+
			"dimension (loose_a, loose_b)", dims)
	}
}

// REQ-108: "Finite product too large to prove deterministically" → lint
// "**Refuse or downgrade** — never silently cap enumeration";
// `graph-product-too-large` "reporting the computed product cardinality and
// the published bound".
// DOMAIN EDGE
func TestReq108_OverLargeFiniteProductRefusesAndReportsSizeAndBound(t *testing.T) {
	reports := guard.Lint(mustLoadSource(t, overLargeProductSource()))

	found := findingsWithCode(reports, guard.CodeProductTooLarge)
	if len(found) == 0 {
		t.Fatalf("no graph-product-too-large finding; findings=%v",
			allFindings(reports))
	}
	if hasGreen(reports) {
		t.Error("the over-large group certified green; enumeration was " +
			"silently capped")
	}
	for _, f := range found {
		if f.ComputedSize == 0 || f.Bound == 0 {
			t.Errorf("finding reports size=%d bound=%d; both figures are "+
				"required", f.ComputedSize, f.Bound)
		}
		if !f.Blocking {
			t.Error("the over-large refusal is not blocking")
		}
	}
}

// REQ-109: "Two escape rows matching one failure class" → "Overlap finding
// within the escape population", `graph-overlap` "naming both escape rule
// ids and the failure class". "Escape row overlapping a guarded row" → "**No
// overlap finding**; the escape row still contributes its assignments to the
// coverage union"; Silent by design.
// DOMAIN EDGE
func TestReq109_EscapePairOverlapsWhileEscapeVsGuardedIsSilentButCounted(t *testing.T) {
	m := mustLoadSource(t, twoPopulationSource())
	reports := guard.Lint(m)

	// Two escape rows for one class: an overlap naming both ids AND the class.
	var found bool
	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		if slices.Contains(f.RuleIDs, "esc-a") && slices.Contains(f.RuleIDs, "esc-b") {
			found = true
			if f.Class != string(resolve.KindNoMatch) {
				t.Errorf("the escape overlap names class %q; want no_match", f.Class)
			}
		}
	}
	if !found {
		t.Errorf("two escape rows matching one class drew no overlap finding; "+
			"findings=%v", allFindings(reports))
	}

	// Escape vs guarded: silent, yet the escape row's assignments still count.
	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		if slices.Contains(f.RuleIDs, "esc-a") && slices.Contains(f.RuleIDs, "ord-a") {
			t.Errorf("finding %v pairs an escape row with a guarded row; it is "+
				"silent by design", f.RuleIDs)
		}
	}
	g := groupOf(t, m, "ord-a")
	if !guard.CoverageUnionFor(m, g, string(resolve.KindNoMatch)).
		Contains(guard.Assignment{"profile": "large"}) {
		t.Error("the escape row's assignments did not reach the coverage " +
			"union; it still contributes them")
	}
}
