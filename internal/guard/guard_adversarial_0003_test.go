package guard_test

// RDR 0003 — Phase 3b adversarial tests.
//
// Each test below targets one failure mode named in the record's
// `Trade-offs / Failure Modes` and `Trade-offs / Risks and Mitigations`
// sections, and is written to FAIL against a defect rather than to
// document current behaviour. Nothing here weakens an existing assertion;
// these are additions.
//
// The oracle throughout is the RUNTIME the lint claim describes. A lint
// verdict is only as good as its agreement with `internal/resolve`, so
// every assertion drives the real kernel over the real loaded model rather
// than restating lint's own arithmetic back at it.

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- ADV-1 ---------------------------------------------------------------

// adv1EmptySetSource is a `set` dimension over a one-element universe,
// partitioned by a `contains` atom and its `unless` complement.
//
// The scoped product is the powerset of the universe, so it carries the
// EMPTY set as an assignment. Both rows are ordinary, both dimensions are
// finitely declared, the key is always-present, and nothing here can
// refuse — so lint has every reason to certify the group green.
func adv1EmptySetSource() string {
	return declBlock(`
[tags.caps]
provenance = "owned"
kind = "set"
elements = ["x"]
required = true
`) + `
[[rule]]
id = "adv1-has"
source = "t:has"
[rule.match.recognized]
eq = "go"
[rule.guard.all.caps]
contains = ["x"]
[rule.write]

[[rule]]
id = "adv1-not"
source = "t:not"
[rule.match.recognized]
eq = "go"
[rule.guard.unless.caps]
contains = ["x"]
[rule.write]
`
}

// ADV-1 — Failure Modes: "Silent failure would be a false exhaustiveness
// claim; the recovery path is to keep every exactness claim tied to A2 and
// the MVV fixture."
//
// A green claim asserts coverage over CONFORMING views (`0003:C` —
// conformance clause, and `GroupReport.ConformingViewsOnly`). The empty set
// is a conforming value of a `set` tag: `Conforms` admits it, and lint's own
// scoped product enumerates it as an assignment, since a set-valued tag
// holds any SUBSET of its element universe.
//
// But the value seam reads a held `[]` as UNEVALUABLE, not false —
// `parseSetLiteral` refuses an empty array because the published literal
// shape is a NON-EMPTY typed element set — so at runtime the `unless` row is
// unevaluable and the kernel refuses `guard_unevaluable`. Lint reaches the
// opposite conclusion because `valueSatisfies` collapses every non-`true`
// verdict, unevaluable included, into "does not satisfy", which credits the
// `unless` complement with a point the runtime cannot decide.
//
// The result is exactly the silent failure the record names: a group
// certified exhaustive whose green claim covers a conforming view the
// runtime refuses.
//
// ADVERSARIAL
func TestAdv1_GreenClaimCoversAViewTheRuntimeRefuses(t *testing.T) {
	m := mustLoadSource(t, adv1EmptySetSource())
	reports := guard.Lint(m)

	// The empty set is a conforming view: lint's own conformance predicate
	// admits it, so it is inside the scope of any green claim.
	emptyHeld := guard.View{"caps": `[]`}
	if err := guard.Conforms(m, emptyHeld); err != nil {
		t.Fatalf("the fixture's premise is gone: the empty held set is no "+
			"longer a conforming view (%v); ADV-1 asserts a green claim "+
			"about conforming views only, so a non-conforming empty set "+
			"would put this case outside the claim", err)
	}

	kt := m.KernelTable()

	// For EVERY conforming view a green group claims to cover, the runtime
	// must reach a decision. A green claim standing beside a
	// `guard_unevaluable` refusal is the false exhaustiveness claim.
	for _, r := range reports {
		if !r.Green {
			continue
		}
		for _, held := range []string{`[]`, `["x"]`} {
			view := guard.View{"caps": held}
			if !r.Covers(view) {
				continue
			}
			res := resolveWith(t, kt, view)
			if res.Refused() && res.Refusal.Kind == resolve.KindGuardUnevaluable {
				t.Errorf("group %s is certified EXHAUSTIVE and its coverage "+
					"union contains the conforming view caps=%s, yet the "+
					"kernel refuses %s on that very view. A green claim that "+
					"covers a view the runtime cannot decide is the false "+
					"exhaustiveness claim the record's Failure Modes name: "+
					"the scoped product enumerates the empty subset, but the "+
					"value seam reads a held `[]` as UNEVALUABLE rather than "+
					"false, and lint's projection collapses the two.",
					r.Context, held, res.Refusal.Kind)
			}
		}
	}

	// A lint-side assertion stood here, guarding a third fix that would have
	// weakened the runtime veto: it required that either the empty subset
	// leave the product, or the group not be green. It was retired (D15)
	// because its own premise — "the runtime finds its guard UNEVALUABLE
	// over a held `[]`" — is no longer true, and the fix that made it false
	// is not the one it guarded against.
	//
	// `0007:C1` separates the two cases this assertion had merged: "An
	// absent set-valued tag MUST be treated as unevaluable under set
	// containment, NOT as the empty set." RDR 0007 owns the
	// `guard_unevaluable` payload, and its A1 records that clause as
	// pinning `contains` against this record's silence. The kernel
	// implements exactly that split — `evaluateAtom` reports UNEVALUABLE on
	// absence, a conflicting view, or a nil seam, never on a present key
	// holding `[]` — so a held empty set is a decided FALSE under REQ-57's
	// total containment, and REQ-67's veto ("a participating row CAN refuse
	// `guard_unevaluable`") is satisfied rather than evaded.
	//
	// Neither alternative the assertion named survives the enforcement
	// surface: dropping the empty subset makes a `set` dimension 2^n-1
	// against REQ-89's pinned powerset, and withholding needs REQ-58
	// unprojectability, which an atom that projects cleanly does not have.
	//
	// The obligation itself is not retired — it is pinned directly, on this
	// same fixture, by TestFixup_EmptyHeldSetIsDecidedRatherThanRefused,
	// which was verified to fail against the pre-fix evaluator.
}

// --- ADV-2 ---------------------------------------------------------------

// adv2BigSetSource declares a `set` tag over an n-element universe and one
// `contains` atom. The scoped product is 2^n, so every n above 11 is
// already past the published bound of 2048 and MUST be refused rather than
// enumerated.
func adv2BigSetSource(n int) string {
	elements := make([]string, 0, n)
	for i := range n {
		elements = append(elements, fmt.Sprintf("%q", fmt.Sprintf("e%02d", i)))
	}
	return declBlock(`
[tags.caps]
provenance = "owned"
kind = "set"
elements = [`+strings.Join(elements, ", ")+`]
required = true
`) + `
[[rule]]
id = "adv2-row"
source = "t:row"
[rule.match.recognized]
eq = "go"
[rule.guard.all.caps]
contains = ["e00"]
[rule.write]
`
}

// ADV-2 — Risks and Mitigations: "**Risk**: Set-valued domains are
// implemented by naive powerset enumeration. **Mitigation**: Require
// symbolic or bitset-equivalent proof and explicit refusal/downgrade for
// finite products that are too large to prove."
//
// The refusal is emitted — `Cardinality` and `Product` both decline in
// microseconds, and lint reports `graph-product-too-large`. But the refusal
// buys nothing, because `Denotation` consults no bound at all: it calls
// `valueAssignments`, which enumerates the FULL powerset of the element
// universe, before lint ever compares the product against `Bound()`. The
// unprovable-dimension scan (`unprovableReason`) reaches `Denotation` for
// every guard atom, so the enumeration the mitigation forbids runs on the
// refusal path itself.
//
// Historical motivation, measured against the unfixed implementation:
// n=16 → 33k assignments in ~160ms; n=20 → 524k in ~2.3s; n=22 → 2.1M in
// ~10s. That is a factor of 1024 above `Bound()` already materialized in
// memory, from an authored model of twenty-odd element names. The bound is
// documented as "the largest product this implementation's enumerating
// proof representation completes over within its budget" — a claim that
// path contradicted. Those numbers explain why the defect mattered; they
// are NOT the oracle below.
//
// The oracle is STRUCTURAL, on the deterministic D11 seam rather than on a
// clock. `valueAssignments` decides from declaration arithmetic — `if n, ok
// := domainSize(d); !ok || n > Bound() { return nil, false }` — and
// DECLINES the dimension before any enumeration, so the atom's denotation
// is unprojectable. A wall-clock ratio could not tell a bounded
// implementation from fast hardware, and preemption or a GC pause could
// fail a correct one; an unprojectable denotation can only be produced by
// the early return, and only the enumerating regression makes it project.
//
// ADVERSARIAL
func TestAdv2_OverLargeSetProductIsRefusedWithoutEnumeratingIt(t *testing.T) {
	const small, large = 14, 18

	// Both products are far past the bound, so BOTH must be refused. That
	// is the premise: this test is about HOW the refusal is reached, never
	// about the verdict, which the suite already pins.
	for _, n := range []int{small, large} {
		m := mustLoadSource(t, adv2BigSetSource(n))
		g := groupOf(t, m, "adv2-row")
		card, ok := guard.Cardinality(m, g)
		if !ok || card <= guard.Bound() {
			t.Fatalf("the fixture's premise is gone: a %d-element universe "+
				"yields cardinality %d (finite=%v), which is not past the "+
				"published bound %d", n, card, ok, guard.Bound())
		}
	}

	m := mustLoadSource(t, adv2BigSetSource(large))
	row := rowByID(t, m, "adv2-row")

	// The `contains` atom over `caps` is the dimension the record forbids
	// enumerating. Its denotation must be WITHHELD, not computed: a
	// projectable denotation over a 2^18 powerset means `valueAssignments`
	// ran the enumeration and `Denotation` materialized the result.
	var found bool
	for _, atom := range row.Atoms {
		if atom.Key != "caps" || atom.Operator == resolve.OpExists {
			continue
		}
		found = true
		if d := guard.Denotation(m, atom.Key, atom); d.Projectable() {
			t.Errorf("the `%s` atom over a %d-element set universe has a "+
				"PROJECTABLE denotation carrying %d assignments. Its value "+
				"dimension is 2^%d, which is past the published bound %d, "+
				"so `valueAssignments` MUST decline it from declaration "+
				"arithmetic before enumerating anything. A projectable "+
				"denotation here means the naive powerset enumeration the "+
				"record's mitigation forbids ran on the refusal path "+
				"itself.", atom.Operator, large, d.Len(), large, guard.Bound())
		}
	}
	if !found {
		t.Fatalf("the fixture's premise is gone: `adv2-row` carries no " +
			"value atom over `caps`, so nothing here exercises the " +
			"over-large value dimension")
	}

	// The row-level surface agrees, and it is the one lint reads.
	if accepted := guard.AcceptedAssignments(m, row); accepted.Projectable() {
		t.Errorf("`adv2-row` was credited with %d accepted assignments over "+
			"a 2^%d value dimension; a dimension past the published bound "+
			"%d carries no enumerable assignment set to accept from",
			accepted.Len(), large, guard.Bound())
	}

	// And the verdict stays pinned: the refusal is still emitted, so this
	// cannot be mistaken for a fix that simply drops the group from lint.
	if reports := guard.Lint(m); hasGreen(reports) {
		t.Errorf("a %d-element set universe certified green; its product "+
			"is past the published bound and MUST be refused. reports=%s",
			large, renderReports(reports))
	}
}

// BenchmarkAdv2RefusalCost keeps the cost signal ADV-2's original oracle
// chased, without asserting on it. A bounded refusal decides both sizes
// from declaration arithmetic, so the two costs sit within a small constant
// factor; naive powerset enumeration makes the ratio 2^(18-14) = 16x. Run
// it by hand when changing the proof representation.
func BenchmarkAdv2RefusalCost(b *testing.B) {
	for _, n := range []int{14, 18} {
		b.Run(fmt.Sprintf("n%02d", n), func(b *testing.B) {
			m, err := table.Load([]byte(adv2BigSetSource(n)), "fixture.toml")
			if err != nil {
				b.Fatalf("fixture must load clean; refused: %v", err)
			}
			b.ResetTimer()
			for b.Loop() {
				guard.Lint(m)
			}
		})
	}
}

// --- ADV-3 ---------------------------------------------------------------

// adv3Source is one scoped row group carrying two rows with BYTE-IDENTICAL
// guards over a fully provable dimension, plus an unrelated third row whose
// guard names a `scalar` key — a kind that by construction carries no
// finite domain.
//
// The overlap between the first two rows is decidable without reference to
// the scalar dimension at all: both rows constrain only `p`, which is
// declared `bool`, single-valued, and always-present.
func adv3Source() string {
	return declBlock(`
[tags.p]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[tags.s]
provenance = "owned"
kind = "scalar"
required = true
`) + `
[[rule]]
id = "adv3-dup-a"
source = "t:dup-a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.p]
eq = true
[rule.write]

[[rule]]
id = "adv3-dup-b"
source = "t:dup-b"
[rule.match.recognized]
eq = "go"
[rule.guard.all.p]
eq = true
[rule.write]

[[rule]]
id = "adv3-opaque"
source = "t:opaque"
[rule.match.recognized]
eq = "go"
[rule.guard.all.s]
eq = "z"
[rule.write]
`
}

// ADV-3 — Failure Modes: "Visible failures should be typed load or lint
// failures: … non-exhaustive finite domain, overlapping candidate rows,
// guard dimension not provable because it lacks a finite domain …".
// Overlapping candidate rows are a VISIBLE lint failure, and the withheld
// path states the same obligation in its own words: "The coverage-gap
// finding is scoped to a provable product, so none is emitted — but overlap
// among the group's decidable rows is unaffected and still MUST be."
//
// It is not. `Product` returns the UNPROJECTABLE set as soon as any one
// dimension lacks a finite domain, and `acceptedIn` short-circuits on an
// unprojectable product — so EVERY row in the group returns the
// unprojectable set, and `pairwiseOverlaps` skips all of them. One opaque
// `scalar` key on an unrelated row therefore silences the overlap check for
// the whole group.
//
// The runtime does not share the blind spot: it refuses `ambiguous_match`
// on the two duplicate rows. So the group's `graph-unprovable-coverage`
// finding is correct and insufficient — the author is told a dimension is
// unprovable and is NOT told two of their rows collide, which is the defect
// they can actually fix.
//
// The existing REQ-44 case covers the OTHER withholding path — a can-refuse
// row over an optional key — where the product stays projectable and
// overlap still computes. This is the unprovable-dimension path, where it
// does not.
//
// ADVERSARIAL
func TestAdv3_UnprovableDimensionSilencesAnUnrelatedOverlap(t *testing.T) {
	m := mustLoadSource(t, adv3Source())
	reports := guard.Lint(m)

	// Premise: the group is withheld for the reason we intend — the opaque
	// dimension — and not for some incidental defect.
	if hasGreen(reports) {
		t.Fatalf("the fixture's premise is gone: the group certified green "+
			"despite an opaque `scalar` dimension; reports=%s",
			renderReports(reports))
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "s"
	}) {
		t.Fatalf("the fixture's premise is gone: no blocking finding names "+
			"the opaque dimension `s`; findings=%v", allFindings(reports))
	}

	// The runtime's verdict on the duplicate pair, which is what lint is
	// describing: both rows qualify, so the kernel cannot pick one.
	kt := m.KernelTable()
	res := resolveWith(t, kt, guard.View{"p": "true", "s": "z"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Fatalf("the fixture's premise is gone: the duplicate rows no "+
			"longer collide at runtime (%s)", describe(res))
	}

	// The obligation. Two ordinary rows carrying byte-identical guards over
	// a fully provable dimension overlap, and the overlap is decidable
	// without consulting the opaque dimension at all.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			containsAll(f.RuleIDs, "adv3-dup-a", "adv3-dup-b")
	}) {
		t.Errorf("no graph-overlap finding names the duplicate pair "+
			"(adv3-dup-a, adv3-dup-b). Both rows constrain only `p` — "+
			"declared bool, single-valued, always-present — so their "+
			"collision is decidable over a fully provable dimension, and "+
			"the kernel refuses ambiguous_match on it. An unrelated row's "+
			"opaque `scalar` key made Product unprojectable for the whole "+
			"group, so acceptedIn short-circuited for EVERY row and "+
			"pairwiseOverlaps skipped them all. Overlapping candidate rows "+
			"are a visible lint failure the record names, and the withheld "+
			"path's own contract says overlap among decidable rows is "+
			"unaffected. Emitted findings=%v", allFindings(reports))
	}

	// Same defect stated over the accepted-assignment surface, so a fix
	// that only special-cases the finding emitter does not satisfy this.
	a := guard.AcceptedAssignments(m, rowByID(t, m, "adv3-dup-a"))
	b := guard.AcceptedAssignments(m, rowByID(t, m, "adv3-dup-b"))
	if !a.Projectable() || !b.Projectable() {
		t.Errorf("the duplicate rows' accepted assignments are "+
			"unprojectable (a=%v b=%v) even though every dimension THEY "+
			"constrain is finitely declared and single-valued. An atom lint "+
			"cannot project belongs to the row that carries it, not to "+
			"every row sharing its group.",
			a.Projectable(), b.Projectable())
	}
}

// containsAll reports whether ids carries every want.
func containsAll(ids []string, want ...string) bool {
	for _, w := range want {
		var found bool
		for _, id := range ids {
			if id == w {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// --- recorded, not asserted ----------------------------------------------

// TestAdv_RecordedIntBoundInversion records a defect on this RDR's own
// exported surface that the RDR 0002 loader currently masks: `table.Load`
// rejects `min > max` before a declaration reaches the guard package.
//
// The declaration model's clause is stated as read "on this RDR's own
// surface, so a declaration built in memory is judged by it too"
// (`declaration.go::agrees`), and `AssignmentCount` is exported. Against an
// inverted-bound declaration it returns a NEGATIVE cardinality reported as
// a valid finite count, and its unmarked sibling panics with "negative
// shift amount". Neither is reachable through the loader today, which is
// why this is recorded as a finding rather than asserted as a failing
// obligation: it is latent, and it is Phase 3c's to weigh.
//
// The test asserts only the reachability premise, so it stays green while
// the loader holds the line and turns red the day that changes.
func TestAdv_RecordedIntBoundInversion(t *testing.T) {
	src := declBlock(`
[tags.n]
provenance = "owned"
kind = "int"
min = 3
max = 0
single_valued = true
required = true
`)
	if _, err := table.Load([]byte(src), "fixture.toml"); err == nil {
		t.Error("the loader now admits an inverted int bound (min=3, max=0). " +
			"guard.AssignmentCount reports a NEGATIVE cardinality as a valid " +
			"finite count for such a declaration, and its unmarked sibling " +
			"panics on a negative shift, so the loader was the only thing " +
			"keeping the defect latent. See ADV-4.")
	}
}

// --- ADV-4b: the int domain WIDTH overflows before it saturates ----------

// wideIntSource is a single-valued, always-present int tag carrying the
// authored bound `{min..max}`, guarded by one row.
//
// The bound is authored as TOML, so the whole chain under test is the one
// an author reaches: `table.Load` normalizes it, `guard` reads the width.
func wideIntSource(minV, maxV int) string {
	return declBlock(`
[tags.n]
provenance = "owned"
kind = "int"
min = `+strconv.Itoa(minV)+`
max = `+strconv.Itoa(maxV)+`
single_valued = true
required = true
`) + `
[[rule]]
id = "wide"
source = "t:wide"
[rule.match.recognized]
eq = "go"
[rule.guard.all.n]
lt = 3
[rule.write]
`
}

// REQ-86: "\"Too large to prove\" MUST be a declared, model-independent
// bound ... the implementation MUST publish the bound it enforces".
//
// An int domain as wide as the int range evaluates `*Max - *Min + 1` in
// WRAPPING arithmetic: `{MinInt..MaxInt}` yields 0, `{0..MaxInt}` yields
// MinInt, `{MinInt+1..MaxInt}` yields -1. D12 saturated the exponent and
// floored a negative domain, but it scoped that to the SHIFT and never
// covered the SUBTRACTION feeding it — so the wrapped width read as a
// domain of size 0 or 1: comfortably UNDER the published bound, fully
// "provable", and GREEN over a dimension carrying 2^64 values. The bound
// comparison stopped being a comparison, exactly what D12 forbids.
//
// This asserts the ARITHMETIC only and never enumerates, so it terminates
// whatever the fix does.
// ADVERSARIAL
func TestAdv_WideIntDomainWidthSaturatesRatherThanWrapping(t *testing.T) {
	cases := []struct {
		name     string
		min, max int
	}{
		{"the whole int range", math.MinInt, math.MaxInt},
		{"the non-negative half", 0, math.MaxInt},
		{"the range less its floor", math.MinInt + 1, math.MaxInt},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			minV, maxV := tc.min, tc.max
			for _, singleValued := range []bool{true, false} {
				d := table.TagDecl{
					Kind: "int", Min: &minV, Max: &maxV,
					SingleValued: singleValued, Required: true,
				}
				n, ok := guard.AssignmentCount(d)
				if !ok {
					t.Fatalf("AssignmentCount({%d..%d}, single_valued=%v) "+
						"carries no finite domain; a fully declared bound "+
						"names a finite — if astronomical — domain",
						minV, maxV, singleValued)
				}
				if n <= guard.Bound() {
					t.Errorf("AssignmentCount({%d..%d}, single_valued=%v) = %d, "+
						"which is at or under the published bound %d — the "+
						"width wrapped instead of saturating, so lint would "+
						"certify an unbounded dimension as provable",
						minV, maxV, singleValued, n, guard.Bound())
				}
			}

			// IntDomain shares the same unguarded expression as its
			// allocation hint. It is asserted through `AssignmentCount`'s
			// ceiling above rather than by calling it here: at HEAD it
			// enumerates the range one value at a time, so a direct call
			// would not return. `Bound()` is what keeps every caller in
			// this package away from it; the guarded width is what keeps a
			// direct caller from sizing a slice off a wrapped count.
		})
	}
}

// REQ-86 / REQ-93: a finite product larger than the bound draws the
// blocking `graph-product-too-large` refusal.
//
// Before the width was saturated this group certified GREEN over an int
// dimension spanning the whole int range, and `guard.Lint` then hung
// enumerating it one value at a time.
// ADVERSARIAL
func TestAdv_WideIntDomainRefusesRatherThanCertifying(t *testing.T) {
	cases := []struct {
		name     string
		min, max int
	}{
		{"the whole int range", math.MinInt, math.MaxInt},
		{"the non-negative half", 0, math.MaxInt},
		{"the range less its floor", math.MinInt + 1, math.MaxInt},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := mustLoadSource(t, wideIntSource(tc.min, tc.max))

			card, ok := guard.Cardinality(m, guard.Groups(m)[0])
			if !ok || card <= guard.Bound() {
				t.Fatalf("Cardinality over {%d..%d} = (%d, %v); want a finite "+
					"count above the bound %d", tc.min, tc.max, card, ok, guard.Bound())
			}

			reports := guard.Lint(m)
			if countCode(reports, guard.CodeProductTooLarge) == 0 {
				t.Errorf("an int dimension spanning {%d..%d} drew no %q "+
					"finding; findings=%v", tc.min, tc.max,
					guard.CodeProductTooLarge, allFindings(reports))
			}
			for _, f := range findingsWithCode(reports, guard.CodeProductTooLarge) {
				if !f.Blocking {
					t.Errorf("%q finding is not blocking", f.Code)
				}
			}
			if hasGreen(reports) {
				t.Errorf("a group over an int dimension spanning {%d..%d} "+
					"certified green", tc.min, tc.max)
			}
		})
	}
}

// --- ADV-4c: a NARROW int domain whose top endpoint is MaxInt ------------

// REQ-108: lint MUST "refuse or downgrade — never silently cap
// enumeration", which presumes the enumeration terminates at all.
//
// `intWidth` (ADV-4b) refuses a width too wide to represent, but
// `{MaxInt..MaxInt}` is only ONE value: the width is 1, the dimension is
// far under the published bound, and `valueAssignments` therefore admits
// it and enters the loop. The loop's `n <= *Max` test is unfalsifiable at
// that endpoint — the `n++` past the last value wraps to MinInt instead of
// exceeding MaxInt — so enumeration never returns and `guard.Lint` does
// not terminate. Saturating the WIDTH does not reach this; only counting
// the iterations does.
//
// The enumeration runs under a watchdog so this test fails rather than
// hangs against an implementation that still wraps.
// ADVERSARIAL
func TestAdv_NarrowIntDomainAtMaxIntTerminates(t *testing.T) {
	cases := []struct {
		name     string
		min, max int
	}{
		{"the single top value", math.MaxInt, math.MaxInt},
		{"the top two values", math.MaxInt - 1, math.MaxInt},
		{"the single bottom value", math.MinInt, math.MinInt},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			minV, maxV := tc.min, tc.max
			d := table.TagDecl{
				Kind: "int", Min: &minV, Max: &maxV,
				SingleValued: true, Required: true,
			}

			want := maxV - minV + 1
			n, ok := guard.AssignmentCount(d)
			if !ok || n != want {
				t.Fatalf("AssignmentCount({%d..%d}) = (%d, %v); want (%d, true)",
					minV, maxV, n, ok, want)
			}

			// `IntDomain` is the exported enumerator; `guard.Lint` below
			// drives the SEPARATE `valueAssignments` loop, so a wrap
			// reintroduced in either one is caught.
			assertTerminates(t, fmt.Sprintf("IntDomain({%d..%d})", minV, maxV),
				func() {
					values := guard.IntDomain(d)
					if len(values) != want {
						t.Errorf("IntDomain({%d..%d}) enumerated %d values; want %d",
							minV, maxV, len(values), want)
					}
					if len(values) > 0 && values[len(values)-1] != maxV {
						t.Errorf("IntDomain({%d..%d}) ended at %d; want %d",
							minV, maxV, values[len(values)-1], maxV)
					}
				})

			// The production path: `Lint` reaches `valueAssignments`,
			// whose own loop carried the same unfalsifiable test. A
			// declaration this narrow is far under the bound, so the
			// dimension is admitted and genuinely enumerated rather than
			// refused — which is what makes this cover the loop.
			m := mustLoadSource(t, wideIntSource(minV, maxV))
			assertTerminates(t, fmt.Sprintf("Lint over {%d..%d}", minV, maxV),
				func() {
					card, ok := guard.Cardinality(m, guard.Groups(m)[0])
					if !ok || card != want {
						t.Errorf("Cardinality over {%d..%d} = (%d, %v); want (%d, true)",
							minV, maxV, card, ok, want)
					}
					if reports := guard.Lint(m); len(allFindings(reports)) == 0 &&
						!hasGreen(reports) {
						t.Errorf("Lint over {%d..%d} reported neither a finding "+
							"nor a green certification", minV, maxV)
					}
				})
		})
	}
}

// assertTerminates runs `body` under a watchdog and fails if it has not
// returned within the deadline.
//
// The watchdog runs in a SUBPROCESS rather than a goroutine. A goroutine
// cannot be cancelled from outside, and the enumeration under test takes no
// cancellation token — deliberately, since D11 keeps these enumerators
// decided from declaration arithmetic rather than threaded with a context
// that exists only for a test. Failing the test while abandoning the worker
// was measured against a reverted tree: the runaway `append` reached 17GB
// resident in the seconds before the process exited, so "the failure is
// reported before it consumes the machine" was not true of a goroutine.
//
// `Process.Kill` ends the runaway for real: the allocation is reclaimed
// with the child's address space, and the parent's own test run is never
// the process holding it. A hard address-space cap would bound the
// footprint directly but is not portable — Darwin rejects `RLIMIT_AS`
// outright and ignores `ulimit -v`, and `GOMEMLIMIT` cannot collect a
// slice the enumeration keeps live — so the deadline is what bounds it.
//
// The child therefore signals over a pipe immediately BEFORE entering the
// body, which splits one deadline into two: a generous allowance for
// process startup, where nothing is enumerating yet, and a tight one for
// the body, which is the only window a regressed tree allocates in. A
// domain of at most two values returns in microseconds, so a body that
// has not finished in well under a second is not slow but looping.
//
// The child re-runs this same test under `-test.run` with
// `adversarialSubprocessEnv` naming the ONE call it should execute, so a
// hang is attributed to the call that hung rather than to every call in
// the test. The parent asserts on the child's exit status.
const adversarialSubprocessEnv = "GUARD_ADV_SUBPROCESS"

// adversarialReadyFD is where `cmd.ExtraFiles[0]` lands in the child: the
// three standard descriptors come first.
const adversarialReadyFD = 3

// Startup and teardown are generous because they are bounded by the
// machine — process spawn and the test framework's exit path — not by the
// code under test, and nothing is enumerating during either. The body
// deadline is the one that matters: it is the only window a regressed
// enumeration allocates in, and it covers microseconds of honest work.
const (
	adversarialStartupDeadline  = 30 * time.Second
	adversarialBodyDeadline     = 500 * time.Millisecond
	adversarialTeardownDeadline = 30 * time.Second
)

func assertTerminates(t *testing.T, what string, body func()) {
	t.Helper()

	// Inside the child: bracket the body with a signal on each side, so the
	// parent can time the BODY rather than the whole process. Only the call
	// this child was spawned for runs; skipping its siblings keeps one hang
	// from being blamed on them.
	if want := os.Getenv(adversarialSubprocessEnv); want != "" {
		if want == what {
			signal := os.NewFile(adversarialReadyFD, "signal")
			if signal != nil {
				_, _ = signal.Write([]byte{1}) // entering the body
			}
			body()
			if signal != nil {
				_, _ = signal.Write([]byte{1}) // the body returned
				_ = signal.Close()
			}
		}
		return
	}

	self, err := os.Executable()
	if err != nil {
		t.Fatalf("could not locate the test binary: %v", err)
	}

	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not open the readiness pipe: %v", err)
	}
	defer pr.Close()

	cmd := exec.Command(self, "-test.run=^"+t.Name()+"$", "-test.v")
	cmd.Env = append(os.Environ(), adversarialSubprocessEnv+"="+what)
	cmd.ExtraFiles = []*os.File{pw} // becomes adversarialReadyFD in the child
	var out strings.Builder
	cmd.Stdout, cmd.Stderr = &out, &out

	if err := cmd.Start(); err != nil {
		pw.Close()
		t.Fatalf("could not start the watchdog subprocess: %v", err)
	}
	pw.Close() // the parent's copy, so the read below sees EOF if the child dies

	// Each signal byte the child writes, delivered as its own receive.
	signals := make(chan struct{}, 2)
	go func() {
		defer close(signals)
		buf := make([]byte, 1)
		for {
			if _, err := pr.Read(buf); err != nil {
				return // the child exited or closed its end
			}
			signals <- struct{}{}
		}
	}()

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	// Three phases, each bounded by what it actually contains. Startup and
	// teardown are bounded by the machine — process spawn, and the test
	// framework's own exit path — so they get generous allowances, and
	// nothing is enumerating during either. Between the two signals sits
	// the body alone: a domain of at most two values, microseconds of
	// honest work, and the ONLY window a regressed tree allocates in. That
	// is why it gets a sub-second deadline. Timing the body by waiting for
	// process exit instead would fold teardown into the same budget and
	// force it wide enough to cover `-race`, which is what a runaway would
	// then have to grow in.
	// `exited` holds the child's result once observed, so a phase that sees
	// the process end does not consume it away from the teardown check.
	var exited error
	var haveExited bool

	awaitSignal := func(phase string, deadline time.Duration) bool {
		// The signals are AUTHORITATIVE, so take one whenever it is already
		// buffered. A fast child can write both bytes and exit before the
		// parent's first receive, leaving this select with two ready cases;
		// `select` would then pick between them at random and report that
		// the child exited before signalling — a flake on correct code.
		select {
		case _, ok := <-signals:
			if ok {
				return true
			}
		default:
		}

		select {
		case _, ok := <-signals:
			if !ok {
				break
			}
			return true
		case err := <-done:
			exited, haveExited = err, true
			// The child may have signalled on its way out. Its write end is
			// closed now that it has exited, so the reader drains what is
			// left and closes `signals`; ranging to completion collects any
			// such byte, and one of them decides this phase.
			for range signals {
				return true
			}
			if err != nil {
				t.Errorf("%s failed in the watchdog subprocess: %v\n%s",
					what, err, out.String())
			} else {
				t.Errorf("%s: the watchdog subprocess exited before it "+
					"signalled %s\n%s", what, phase, out.String())
			}
			return false
		case <-time.After(deadline):
			_ = cmd.Process.Kill()
			<-done
			if phase == "the body returning" {
				t.Errorf("%s did not terminate: the loop test `n <= max` "+
					"cannot fail at MaxInt, so enumerating a narrow "+
					"dimension runs forever", what)
			} else {
				t.Errorf("%s: the watchdog subprocess never signalled %s\n%s",
					what, phase, out.String())
			}
			return false
		}
		t.Errorf("%s: the watchdog subprocess closed its signal pipe before "+
			"signalling %s\n%s", what, phase, out.String())
		return false
	}

	if !awaitSignal("entering the body", adversarialStartupDeadline) {
		return
	}
	if !awaitSignal("the body returning", adversarialBodyDeadline) {
		return
	}

	if haveExited {
		if exited != nil {
			t.Errorf("%s failed in the watchdog subprocess: %v\n%s",
				what, exited, out.String())
		}
		return
	}

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("%s failed in the watchdog subprocess: %v\n%s",
				what, err, out.String())
		}
	case <-time.After(adversarialTeardownDeadline):
		_ = cmd.Process.Kill()
		<-done
		t.Errorf("%s: the watchdog subprocess did not exit after its body "+
			"returned\n%s", what, out.String())
	}
}

// --- ADV-4d: the OPTIONAL-presence factor overflows ----------------------

// optionalWideIntSource is `wideIntSource`'s optional twin: the same
// single-valued int bound, but `required = false`, so `AssignmentCount`
// applies the presence factor of two on top of the value domain.
func optionalWideIntSource(minV, maxV int) string {
	return declBlock(`
[tags.n]
provenance = "owned"
kind = "int"
min = `+strconv.Itoa(minV)+`
max = `+strconv.Itoa(maxV)+`
single_valued = true
required = false
`) + `
[[rule]]
id = "wide-opt"
source = "t:wide-opt"
[rule.match.recognized]
eq = "go"
[rule.guard.all.n]
lt = 3
[rule.write]
`
}

// REQ-86 / REQ-93: the published bound must be ENFORCED, which a wrapped
// comparison does not do.
//
// The presence factor is the LAST arithmetic standing between a declared
// domain and the bound comparison, and ADV-4b left it unguarded. A
// single-valued `{MinInt..-2}` has a width of MaxInt — honest, and
// `intWidth` rightly admits it — but `count *= 2` for an OPTIONAL key
// wraps it to -2. That reads as under the bound of 2048: fully provable,
// no `graph-product-too-large`, and GREEN over a dimension spanning
// essentially the whole int range. It is D12's defect ("a wrapped negative
// would read as under the bound and certify the very product the clause
// refuses") displaced one operation later.
//
// This asserts the arithmetic and the refusal without enumerating.
// ADVERSARIAL
func TestAdv_OptionalPresenceFactorSaturatesRatherThanWrapping(t *testing.T) {
	cases := []struct {
		name     string
		min, max int
	}{
		{"the negative half less its top", math.MinInt, -2},
		{"the non-negative half less its top", 0, math.MaxInt - 1},
		{"the whole int range", math.MinInt, math.MaxInt},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			minV, maxV := tc.min, tc.max
			d := table.TagDecl{
				Kind: "int", Min: &minV, Max: &maxV,
				SingleValued: true, Required: false,
			}

			n, ok := guard.AssignmentCount(d)
			if !ok {
				t.Fatalf("AssignmentCount({%d..%d}, optional) carries no "+
					"finite domain; a fully declared bound names a finite — "+
					"if astronomical — domain", minV, maxV)
			}
			if n <= guard.Bound() {
				t.Fatalf("AssignmentCount({%d..%d}, optional) = %d, at or "+
					"under the published bound %d — the presence factor "+
					"wrapped instead of saturating, so lint would certify an "+
					"unbounded dimension as provable", minV, maxV, n, guard.Bound())
			}

			m := mustLoadSource(t, optionalWideIntSource(minV, maxV))
			card, ok := guard.Cardinality(m, guard.Groups(m)[0])
			if !ok || card <= guard.Bound() {
				t.Fatalf("Cardinality over optional {%d..%d} = (%d, %v); want "+
					"a finite count above the bound %d",
					minV, maxV, card, ok, guard.Bound())
			}

			reports := guard.Lint(m)
			if countCode(reports, guard.CodeProductTooLarge) == 0 {
				t.Errorf("an optional int dimension spanning {%d..%d} drew no "+
					"%q finding; findings=%v", minV, maxV,
					guard.CodeProductTooLarge, allFindings(reports))
			}
			for _, f := range findingsWithCode(reports, guard.CodeProductTooLarge) {
				if !f.Blocking {
					t.Errorf("%q finding is not blocking", f.Code)
				}
			}
			if hasGreen(reports) {
				t.Errorf("a group over an optional int dimension spanning "+
					"{%d..%d} certified green", minV, maxV)
			}
		})
	}
}
