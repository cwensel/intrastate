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
	"strings"
	"testing"
	"time"

	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
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

	// And state the same defect from the lint side, so a fix that merely
	// stops the kernel refusing does not satisfy this test: either the
	// empty subset leaves the product, or the group must not be green.
	g := groupOf(t, m, "adv1-has")
	notRow := rowByID(t, m, "adv1-not")
	accepted := guard.AcceptedAssignments(m, notRow)
	if accepted.Projectable() && accepted.Contains(guard.Assignment{"caps": `[]`}) {
		t.Errorf("the `unless contains` row accepts the empty-subset "+
			"assignment, but the runtime finds its guard UNEVALUABLE over a "+
			"held `[]` and refuses. Lint MUST NOT credit a row with an "+
			"assignment the runtime cannot decide; product=%d union=%d",
			guard.Product(m, g).Len(), guard.CoverageUnion(m, g).Len())
	}
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
// Cost measured against this implementation: n=16 → 33k assignments in
// ~160ms; n=20 → 524k in ~2.3s; n=22 → 2.1M in ~10s. That is a factor of
// 1024 above `Bound()` already materialized in memory, from an authored
// model of twenty-odd element names. The bound is documented as "the
// largest product this implementation's enumerating proof representation
// completes over within its budget" — a claim this path contradicts.
//
// The test asserts the mitigation, not a wall-clock number: lint over a
// refused product must not scale with the product's size. Two universes
// four elements apart both refuse, so a bounded implementation spends
// comparable time on each; a naively enumerating one spends 16x.
//
// ADVERSARIAL
func TestAdv2_OverLargeSetProductIsRefusedWithoutEnumeratingIt(t *testing.T) {
	const small, large = 14, 18

	// Both products are far past the bound, so BOTH must be refused. That
	// is the premise: this test is about the cost of the refusal, never
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

	measure := func(n int) time.Duration {
		m := mustLoadSource(t, adv2BigSetSource(n))
		start := time.Now()
		reports := guard.Lint(m)
		elapsed := time.Since(start)
		if hasGreen(reports) {
			t.Fatalf("a %d-element set universe certified green; its product "+
				"is past the published bound and MUST be refused", n)
		}
		return elapsed
	}

	// Warm the loader and any lazily-built state so the comparison measures
	// the proof representation rather than first-call overhead.
	measure(small)

	smallCost := measure(small)
	largeCost := measure(large)

	// A bounded (symbolic or bitset-equivalent) refusal decides both from
	// declaration arithmetic, so the two costs sit within a small constant
	// factor. Naive powerset enumeration makes the ratio 2^(large-small) =
	// 16x. Eight is a deliberately slack threshold: it clears any plausible
	// constant-factor difference while still catching the doubling.
	const tolerance = 8
	if largeCost > smallCost*tolerance && largeCost > 50*time.Millisecond {
		t.Errorf("lint over a REFUSED set product scales with the product's "+
			"size: a %d-element universe cost %v and an %d-element one cost "+
			"%v (%.1fx for a %dx larger powerset). Both products are past "+
			"the published bound %d and both are refused, so a refusal that "+
			"stays bounded would cost about the same for each. The "+
			"enumeration happens in Denotation, which consults no bound "+
			"before calling valueAssignments — so the naive powerset "+
			"enumeration the record's mitigation forbids runs on the "+
			"refusal path itself.",
			small, smallCost, large, largeCost,
			float64(largeCost)/float64(max(smallCost, 1)), 1<<(large-small),
			guard.Bound())
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
