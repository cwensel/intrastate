package guard_test

// RDR 0003 — Phase 3c regression tests.
//
// One test per defect closed in the fixup phase that did not already carry
// a failing test of its own. The adversarial cases (ADV-1, ADV-2, ADV-3)
// arrived with theirs in `guard_adversarial_0003_test.go`; FAIL-1, found
// by the Phase 3a CoVe pass, did not.

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
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
	//
	// The projectability of the union is part of the claim, not a
	// precondition on it: `coverageUnion` returns a ZERO `AssignmentSet{}`
	// when the group's product is unprojectable, and `Equal` is false
	// whenever either side is unprojectable. Guarding the comparison on
	// `union.Projectable()` would therefore let that early return — the
	// regression this test exists to catch — pass as a non-finding.
	union := guard.CoverageUnion(m, g)
	if !union.Projectable() || !dec.Equal(union) {
		t.Errorf("the coverage union (projectable=%v len=%d) is not the "+
			"union of the group's DECIDABLE rows alone (len=%d). An "+
			"UNPROJECTABLE union is itself the regression: `coverageUnion` "+
			"yields the zero `AssignmentSet{}` when the group's product "+
			"does not project, and this group's product does — so a "+
			"withheld union here means the narrowing withdrew the whole "+
			"group rather than the refusing row's contribution. A "+
			"projectable union unequal to `dec` means a can-refuse row "+
			"contributed assignments to it",
			union.Projectable(), union.Len(), dec.Len())
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

// adv2HugeSetSource declares a `set` universe whose powerset exponent is at
// or past the width of the machine integer the cardinality is computed in.
func adv2HugeSetSource(n int) string {
	elements := make([]string, 0, n)
	for i := range n {
		elements = append(elements, strconv.Quote(fmt.Sprintf("e%03d", i)))
	}
	return declBlock(`
[tags.caps]
provenance = "owned"
kind = "set"
elements = [`+strings.Join(elements, ", ")+`]
required = true
`) + `
[[rule]]
id = "huge-row"
source = "t:huge"
[rule.match.recognized]
eq = "go"
[rule.guard.all.caps]
contains = ["e000"]
[rule.write]
`
}

// ADV-2, second half — the bound comparison must stay a COMPARISON.
//
// REQ-86 requires the published bound to be enforced, and REQ-108 requires
// an over-large finite product to be refused, "never silently cap[ped]".
// `Cardinality` already saturates rather than overflowing, on the stated
// grounds that "a wrapped negative would read as under the bound and
// certify the very product the clause refuses" — but `domainSize`'s own
// `1 << |universe|` did not. In Go a shift at or past the integer width
// yields ZERO, so a 64-element element universe reported a value dimension
// of size 0 and a whole-product cardinality of 0: comfortably under the
// bound, fully "provable", and green.
//
// The exponent saturates at the same ceiling `Cardinality` uses, so the
// comparison holds at every universe size the loader admits.
func TestFixup_HugeSetUniverseIsRefusedRatherThanWrappingUnderTheBound(t *testing.T) {
	for _, n := range []int{40, 63, 64, 65, 96} {
		m := mustLoadSource(t, adv2HugeSetSource(n))
		g := groupOf(t, m, "huge-row")

		card, ok := guard.Cardinality(m, g)
		if ok && card <= guard.Bound() {
			t.Errorf("a %d-element set universe reports cardinality %d, "+
				"which is at or under the published bound %d. The powerset "+
				"of %d elements is 2^%d — the bound comparison wrapped "+
				"instead of comparing, so lint would certify a product it "+
				"cannot enumerate", n, card, guard.Bound(), n, n)
		}
		if guard.Product(m, g).Projectable() {
			t.Errorf("a %d-element set universe yielded a PROJECTABLE "+
				"product; enumerating 2^%d assignments is exactly the "+
				"silent cap the record forbids", n, n)
		}
		if reports := guard.Lint(m); hasGreen(reports) {
			t.Errorf("a %d-element set universe certified GREEN; reports=%s",
				n, renderReports(reports))
		}
	}
}

// ADV-4 — an inverted int bound on this RDR's own exported surface.
//
// Phase 3b recorded this as LATENT: `table.Load` refuses `min > max`, so no
// declaration reaching the guard package THROUGH THE LOADER could trigger
// it. But `AssignmentCount` is exported and takes a `table.TagDecl`
// directly, and `agrees`' own doc comment states the agreement rule is read
// "on this RDR's own surface, so a declaration built in memory is judged by
// it too" — which it was not: `agrees` never compared the endpoints.
//
// Against `{3..0}` the table returned a NEGATIVE cardinality reported as a
// valid finite count, and the unmarked sibling panicked in `spread` on
// `1 << -2`. REQ-18 fixes `{min..max}` as inclusive at both endpoints,
// which presumes an ordered pair; an inverted bound names no value, so it
// is not a readable finite domain and the table defines no count for it.
//
// The tripwire `TestAdv_RecordedIntBoundInversion` stays green either way —
// it asserts only that the loader still holds the line. This asserts the
// surface itself.
func TestFixup_InvertedIntBoundCarriesNoReadableDomain(t *testing.T) {
	inverted := func(singleValued bool) table.TagDecl {
		lo, hi := 3, 0
		return table.TagDecl{
			Kind: "int", Min: &lo, Max: &hi,
			SingleValued: singleValued, Required: true,
		}
	}

	for _, singleValued := range []bool{true, false} {
		// Must not panic, and must not report a count.
		n, ok := guard.AssignmentCount(inverted(singleValued))
		if ok {
			t.Errorf("AssignmentCount({int, min:3, max:0, single_valued:%v}) "+
				"= (%d, true); an inverted bound names no value, so it is "+
				"not a readable finite domain and carries NO assignment "+
				"count. Reporting one lets a dimension that cannot be "+
				"enumerated read as provable", singleValued, n)
		}
		if n < 0 {
			t.Errorf("AssignmentCount reported the NEGATIVE cardinality %d "+
				"as a count; a negative reads as under the published bound "+
				"and certifies a product that cannot be proved", n)
		}
	}

	// The ordered sibling is untouched: `{0..3}` is inclusive at both
	// endpoints, so it counts four.
	lo, hi := 0, 3
	ordered := table.TagDecl{
		Kind: "int", Min: &lo, Max: &hi, SingleValued: true, Required: true,
	}
	if n, ok := guard.AssignmentCount(ordered); !ok || n != 4 {
		t.Errorf("AssignmentCount({int, min:0, max:3, single_valued}) = "+
			"(%d, %v); REQ-18 fixes both endpoints inclusive, so `{0..3}` "+
			"has cardinality 4. The inverted-bound guard must not narrow "+
			"the ordered case", n, ok)
	}
}

// ADV-1 — the lint claim and the runtime must agree about the empty subset.
//
// REQ-67: "An exhaustiveness claim MUST NOT be stronger than the runtime it
// describes." The two surfaces disagreed about whether a held `[]` is a
// value at all:
//
//   - REQ-89 makes a `set` dimension `2^|element universe|` — "a set-valued
//     tag holds any SUBSET" — so the scoped product enumerates the empty
//     subset, and `Conforms` admits `caps=[]` as a conforming view.
//   - `Evaluator.Evaluate` refused a held `[]`, because its `contains` arm
//     parsed the HELD value with `parseSetLiteral`, whose non-empty rule is
//     the operator/kind matrix's published shape for the authored LITERAL
//     ("non-empty typed element set", REQ-4).
//
// REQ-57 settles which surface was wrong: `contains` denotes "the
// assignments whose held set CONTAINS every listed element", and containment
// is total over sets — `[]` does not contain `x`, so the answer is FALSE,
// not undecidable. The literal's shape rule never governed the held value.
//
// Dropping the empty subset from the product instead would contradict
// REQ-89's `2^|universe|`, which its own test pins. So the fix is on the
// runtime side, and the two surfaces now agree: the `contains` / `unless
// contains` pair is a genuine complete partition of the powerset.
func TestFixup_EmptyHeldSetIsDecidedRatherThanRefused(t *testing.T) {
	m := mustLoadSource(t, adv1EmptySetSource())
	kt := m.KernelTable()

	// The runtime decides every conforming view, the empty subset included.
	for _, held := range []string{`[]`, `["x"]`} {
		view := guard.View{"caps": held}
		if err := guard.Conforms(m, view); err != nil {
			t.Fatalf("the fixture's premise is gone: caps=%s is no longer a "+
				"conforming view (%v)", held, err)
		}
		res := resolveWith(t, kt, view)
		if res.Refused() {
			t.Errorf("the kernel refuses %s on the conforming view caps=%s. "+
				"`contains` denotes the assignments whose held set CONTAINS "+
				"every listed element (REQ-57), and containment is total "+
				"over sets — the empty subset simply does not contain `x`, "+
				"which is FALSE, not undecidable. The non-empty rule belongs "+
				"to the authored LITERAL (REQ-4), not to the held value",
				res.Refusal.Kind, held)
		}
	}

	// The empty subset is IN the product: REQ-89 makes a set dimension the
	// full powerset of its universe, never the universe minus a member.
	g := groupOf(t, m, "adv1-has")
	product := guard.Product(m, g)
	if !product.Projectable() || !product.Contains(guard.Assignment{"caps": `[]`}) {
		t.Errorf("the scoped product does not carry the empty-subset "+
			"assignment (projectable=%v len=%d). A set-valued tag holds any "+
			"SUBSET, so its assignment space is 2^|universe| — dropping the "+
			"empty subset would make it 2^n-1 and contradict REQ-89",
			product.Projectable(), product.Len())
	}

	// And every conforming view a GREEN group claims to cover, the runtime
	// decides. This is the obligation ADV-1 names, stated directly.
	for _, r := range guard.Lint(m) {
		if !r.Green {
			continue
		}
		for _, held := range []string{`[]`, `["x"]`} {
			view := guard.View{"caps": held}
			if !r.Covers(view) {
				continue
			}
			if res := resolveWith(t, kt, view); res.Refused() {
				t.Errorf("group %s is certified EXHAUSTIVE and covers the "+
					"conforming view caps=%s, yet the kernel refuses %s on "+
					"it. A green claim MUST NOT be stronger than the runtime "+
					"it describes (REQ-67)", r.Context, held, res.Refusal.Kind)
			}
		}
	}
}

// A dimensionless group — one carrying no guard dimension at all — makes no
// exhaustiveness claim, and lint correctly withholds one. But withholding the
// CLAIM is not withholding every check: the overlap invariant ("no
// finite-domain input assignment may enable two ordinary rows") carries no
// finitely-declared precondition, because it reads declarations rather than a
// product. Two guard-atom-free ordinary rows in one selection context are
// enabled by EVERY assignment, so returning before `overlapFindings` emitted
// no verdict and no finding at all — lint exited clean on a model the kernel
// always refuses `ambiguous_match`. Found by roborev triage (job 6087).
// ADVERSARIAL
func TestFixup_DimensionlessGroupStillReportsOrdinaryOverlap(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`)+`
[[rule]]
id = "a"
source = "t:a"
[rule.match.recognized]
eq = "go"
[rule.write]

[[rule]]
id = "b"
source = "t:b"
[rule.match.recognized]
eq = "go"
[rule.write]
`)

	g := groupOf(t, m, "a")
	if len(guard.Dimensions(m, g)) != 0 {
		t.Fatalf("the fixture's premise is gone: the group carries %d "+
			"dimensions, so it is not the dimensionless case under test",
			len(guard.Dimensions(m, g)))
	}

	reports := guard.Lint(m)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap
	}) {
		t.Errorf("two unguarded ordinary rows in one selection context "+
			"produced no `graph-overlap` finding; findings=%v. Every "+
			"assignment enables both rows, and the kernel refuses "+
			"`ambiguous_match` on this model, so lint exiting clean here is "+
			"a false green", allFindings(reports))
	}
}
