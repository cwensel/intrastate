package guard_test

// RDR 0003 — Phase 3c regression tests.
//
// One test per defect closed in the fixup phase that did not already carry
// a failing test of its own. The adversarial cases (ADV-1, ADV-2, ADV-3)
// arrived with theirs in `guard_adversarial_0003_test.go`; FAIL-1, found
// by the Phase 3a CoVe pass, did not.

import (
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/guard"
)

// fail1CanRefuseSource is a scoped row group in which every dimension is
// finite and single-valued — so every atom in it PROJECTS — while one row
// still carries a value atom over a key declared optional.
//
// That combination is the narrowing's own interesting case: the row can
// refuse `guard_unevaluable` at runtime, yet nothing about its atoms is
// unprojectable, so a projection-only test for undecidability misses it.
// `dec` and `refuser` denote disjoint value subsets of two independent
// dimensions, so their `all`-intersections necessarily intersect in the
// product — which is what makes the mistaken overlap finding visible.
func fail1CanRefuseSource() string {
	return declBlock(`
[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["x", "y"]
single_valued = true

[tags.req]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`) + `
[[rule]]
id = "fail1-dec"
source = "t:dec"
[rule.match.recognized]
eq = "go"
[rule.guard.all.req]
eq = "a"
[rule.write]

[[rule]]
id = "fail1-refuser"
source = "t:refuser"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "x"
[rule.write]
`
}

// FAIL-1 — REQ-43: "**A can-refuse row contributes no assignments to the
// coverage union** … a row that can refuse has **no decidable
// accepted-assignment set to contribute**, so lint MUST NOT credit it with
// its `all`-intersection unsubtracted"; and REQ-44: "Overlap findings among
// the group's **decidable** rows are unaffected and still MUST be emitted."
//
// The withholding decision (`CanRefuse`) and the projection decision
// (`Denotation`) are independent: a can-refuse row whose every dimension is
// finitely declared and single-valued projects perfectly well. Before the
// fix, `acceptedIn` consulted only projection, so such a row was credited
// with an accepted-assignment set, entered the coverage union, and was
// paired by the overlap check against the group's genuinely decidable rows.
//
// `Lint`'s group verdict never depended on it — it withholds on `CanRefuse`
// before reading the union — so the violation was confined to the exported
// denotation/union surface and to the emitted overlap set, which is where
// this test looks.
func TestFail1_CanRefuseRowContributesNoAcceptedAssignments(t *testing.T) {
	m := mustLoadSource(t, fail1CanRefuseSource())
	g := groupOf(t, m, "fail1-dec")

	// Premise: this is the projectable case. Every dimension is finite and
	// single-valued, so the group's product exists and the refusing row's
	// own atom projects — an undecidability test reading projection alone
	// finds nothing wrong here.
	if !guard.Product(m, g).Projectable() {
		t.Fatalf("the fixture's premise is gone: the group's product is no " +
			"longer projectable, so the can-refuse row's crediting could be " +
			"suppressed by projection rather than by the narrowing")
	}
	refuser := rowByID(t, m, "fail1-refuser")
	if !guard.CanRefuse(m.Tags, refuser) {
		t.Fatalf("the fixture's premise is gone: `fail1-refuser` no longer " +
			"carries a value atom over a key declared optional")
	}

	// REQ-43, first half: the row has no decidable accepted-assignment set.
	if accepted := guard.AcceptedAssignments(m, refuser); accepted.Projectable() {
		t.Errorf("a can-refuse row was credited with %d accepted "+
			"assignments; a row that can refuse has NO decidable "+
			"accepted-assignment set to contribute, so lint MUST NOT credit "+
			"it with its `all`-intersection unsubtracted", accepted.Len())
	}

	// The decidable row keeps its own set: withholding is a property of the
	// refusing row, never of every row that shares its group.
	dec := guard.AcceptedAssignments(m, rowByID(t, m, "fail1-dec"))
	if !dec.Projectable() || dec.Len() == 0 {
		t.Errorf("the group's DECIDABLE row lost its accepted assignments "+
			"(projectable=%v len=%d); the narrowing withholds the group's "+
			"claim, it does not make every row undecidable",
			dec.Projectable(), dec.Len())
	}

	// REQ-43, second half: nothing the refusing row would accept reaches
	// the union. `opt=x` is the refuser's own denotation; no decidable row
	// in the group constrains `opt` at all.
	union := guard.CoverageUnion(m, g)
	if union.Projectable() && !dec.Equal(union) {
		t.Errorf("the coverage union (%d) is not the union of the group's "+
			"DECIDABLE rows alone (%d); a can-refuse row contributed "+
			"assignments to it", union.Len(), dec.Len())
	}

	// REQ-44: the surviving overlap check is scoped to the group's
	// decidable rows, so a pair involving the refusing row is never
	// intersected.
	reports := guard.Lint(m)
	for _, f := range allFindings(reports) {
		if f.Code == guard.CodeOverlap && slices.Contains(f.RuleIDs, "fail1-refuser") {
			t.Errorf("graph-overlap names %v, which includes the can-refuse "+
				"row `fail1-refuser`. REQ-44 scopes the surviving overlap "+
				"check to the group's DECIDABLE rows; a row with no "+
				"decidable accepted-assignment set has nothing to intersect, "+
				"so the pair is not a runtime ambiguity and MUST NOT be "+
				"reported as one", f.RuleIDs)
		}
	}

	// And the withholding itself still fires, so this fix cannot be
	// mistaken for one that simply drops the refusing row from lint.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage &&
			slices.Contains(f.RuleIDs, "fail1-refuser")
	}) {
		t.Errorf("no blocking withholding finding names `fail1-refuser`; a "+
			"withheld claim MUST be observable, not silent. findings=%v",
			allFindings(reports))
	}
}
