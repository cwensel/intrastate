package guard_test

// RDR 0003 — row groups, participation, the scoped product, atom
// projection, and exhaustiveness eligibility.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-45: "Coverage and overlap checks MUST be scoped to a normalized row
// group supplied by the transition/lint model, and MUST evaluate the
// participating guard dimensions as one product rather than as independent
// one-dimensional checks."
// ADVERSARIAL
func TestReq45_DimensionsAreEvaluatedAsOneProductNotIndependently(t *testing.T) {
	// This group covers each dimension completely when the two are checked
	// independently, yet leaves a hole in the 2-D product. Independent
	// one-dimensional checks would certify it green.
	m := mustLoadSource(t, productOnlyGapSource())
	reports := guard.Lint(m)

	if hasGreen(reports) {
		t.Fatalf("a group with a gap visible only in the multi-dimensional "+
			"product certified green; the dimensions were checked "+
			"independently. reports=%s", renderReports(reports))
	}
	gaps := findingsWithCode(reports, guard.CodeCoverageGap)
	if len(gaps) == 0 {
		t.Fatalf("no coverage gap finding; findings=%v", allFindings(reports))
	}
	// The witness must be the 2-D hole, not a 1-D one.
	for _, f := range gaps {
		if len(f.Witness) < 2 {
			t.Errorf("gap witness %v names %d dimensions; the hole is visible "+
				"only in the multi-dimensional product", f.Witness, len(f.Witness))
		}
	}
}

// REQ-46: "**The row group is defined here**, not deferred: a scoped row
// group is the set of normalized candidate rows sharing one selection
// context — the same source state and the same recognized outcome — since
// that is exactly the set RDR 0001 resolves exact-one over."
// HAPPY PATH
func TestReq46_RowGroupIsTheRowsSharingOneSelectionContext(t *testing.T) {
	m := mustLoadSource(t, twoContextSource())
	groups := guard.Groups(m)

	if len(groups) != 2 {
		t.Fatalf("Groups returned %d groups; the fixture carries two selection "+
			"contexts (two source states over one outcome)", len(groups))
	}
	for _, g := range groups {
		ids := ruleIDsOf(g)
		switch {
		case slices.Contains(ids, "draft-a"):
			if !slices.Contains(ids, "draft-b") {
				t.Errorf("group %v splits two rows sharing one selection "+
					"context", ids)
			}
			if slices.Contains(ids, "final-a") {
				t.Errorf("group %v merges two different source states", ids)
			}
		case slices.Contains(ids, "final-a"):
			if slices.Contains(ids, "draft-a") {
				t.Errorf("group %v merges two different source states", ids)
			}
		default:
			t.Errorf("unexpected group %v", ids)
		}
	}

	// A differing recognized outcome also splits the group.
	other := mustLoadSource(t, twoOutcomeSource())
	if n := len(guard.Groups(other)); n != 2 {
		t.Errorf("Groups over two recognized outcomes returned %d groups; the "+
			"outcome is half the selection context", n)
	}
}

// REQ-47: "RDR 0006 supplies the graph traversal that enumerates which
// selection contexts are reachable; it does not define the grouping
// predicate, and this RDR does not read one back from it."
// BOUNDARY
func TestReq47_GroupingPredicateIsSelfContainedAndReadsNoTraversal(t *testing.T) {
	m := mustLoadSource(t, twoContextSource())

	// Groups is total over the model alone: no traversal, no reachability
	// input, no second grouping source. Calling it twice yields the same
	// grouping.
	first := guard.Groups(m)
	second := guard.Groups(m)
	if len(first) != len(second) {
		t.Fatalf("Groups is not a function of the model: %d vs %d groups",
			len(first), len(second))
	}
	for i := range first {
		if ctx := first[i].Context.String(); ctx != second[i].Context.String() {
			t.Errorf("group %d context differs across calls: %q vs %q",
				i, ctx, second[i].Context.String())
		}
	}

	// Every row lands in exactly one group — the grouping is a partition,
	// which an externally-supplied reachability set could not guarantee.
	seen := map[string]int{}
	for _, g := range first {
		for _, id := range ruleIDsOf(g) {
			seen[id]++
		}
	}
	for _, r := range m.Rows {
		if seen[r.RuleID] != 1 {
			t.Errorf("row %q lands in %d groups; the grouping predicate is a "+
				"partition over the model's own rows", r.RuleID, seen[r.RuleID])
		}
	}
}

// REQ-48: "A dimension **participates** in a row group when any row in that
// group carries a **guard** atom over that key, in either `all` or
// `unless` — not only when the rows constrain it differently. The authored
// block is what makes an atom a guard atom (JDR 0001 §D6): every atom
// under `guard.all` or `guard.unless` is a guard atom regardless of
// operator, `eq` included."
// DOMAIN EDGE
func TestReq48_ParticipationIsAnyGuardAtomInEitherBlockRegardlessOfOperator(t *testing.T) {
	m := mustLoadSource(t, participationSource())
	g := groupOf(t, m, "part-a")

	dims := guard.Dimensions(m, g)
	// `only_in_all` is carried by one row only, under `all`, with `eq`.
	// `only_in_unless` is carried by one row only, under `unless`.
	for _, key := range []string{"only_in_all", "only_in_unless"} {
		if !slices.Contains(dims, key) {
			t.Errorf("dimension %q does not participate; participation is ANY "+
				"guard atom in EITHER block, not only a differing constraint. "+
				"dimensions=%v", key, dims)
		}
	}
}

// REQ-49: "A key every row constrains identically still bounds the product
// and still requires a finite declared domain; it MUST NOT be dropped from
// the product because it does not discriminate."
// ADVERSARIAL
func TestReq49_IdenticallyConstrainedKeyStaysInTheProduct(t *testing.T) {
	m := mustLoadSource(t, identicalConstraintSource())
	g := groupOf(t, m, "same-a")

	if !slices.Contains(guard.Dimensions(m, g), "shared") {
		t.Errorf("a key every row constrains identically was dropped from the "+
			"product; dimensions=%v", guard.Dimensions(m, g))
	}
	// It still bounds the product: the cardinality includes its factor.
	card, ok := guard.Cardinality(m, g)
	if !ok {
		t.Fatal("the group's product has no cardinality")
	}
	if card%2 != 0 {
		t.Errorf("cardinality %d does not carry the shared dimension's "+
			"factor; the key still bounds the product", card)
	}

	// And it still requires a finite declared domain: undeclaring it makes
	// the group unprovable rather than dropping it.
	unbounded := mustLoadSource(t, identicalConstraintUnboundedSource())
	reports := guard.Lint(unbounded)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "shared"
	}) {
		t.Errorf("an identically-constrained key with no finite domain drew "+
			"no blocking finding; findings=%v", allFindings(reports))
	}
}

// REQ-50: "**Match keys are not product dimensions.** A row's match
// pattern selects which group the row belongs to — the selection context
// this RDR groups by — and is a separate field from its guard"
// ADVERSARIAL
func TestReq50_MatchKeysAreNotProductDimensions(t *testing.T) {
	m := mustLoadSource(t, twoContextSource())
	g := groupOf(t, m, "draft-a")

	if slices.Contains(guard.Dimensions(m, g), "status") {
		t.Errorf("the match key `status` entered the product as a dimension; "+
			"a match pattern selects the GROUP, it is not a dimension within "+
			"it. dimensions=%v", guard.Dimensions(m, g))
	}
	// It is the grouping context instead.
	if !slices.Contains(g.Context.MatchKeys(), "status") {
		t.Errorf("the match key `status` is not part of the selection "+
			"context; context=%s", g.Context)
	}
}

// REQ-51: "**The match/guard split is per atom per group, not per key per
// model.** The same tag key MAY be a match key in one rule and a guard key
// in another … a key enters that group's product when some row in *that*
// group carries a guard atom over it, regardless of how the key is used in
// any other group."
// DOMAIN EDGE
func TestReq51_MatchGuardSplitIsPerAtomPerGroupNotPerKeyPerModel(t *testing.T) {
	m := mustLoadSource(t, dualUseKeySource())

	guarded := groupOf(t, m, "guarded-use")
	matched := groupOf(t, m, "matched-use")

	if !slices.Contains(guard.Dimensions(m, guarded), "profile") {
		t.Errorf("`profile` did not enter the product of the group that "+
			"GUARDS it; dimensions=%v", guard.Dimensions(m, guarded))
	}
	if slices.Contains(guard.Dimensions(m, matched), "profile") {
		t.Errorf("`profile` entered the product of the group that only "+
			"MATCHES it; the split is per atom per group, not per key per "+
			"model. dimensions=%v", guard.Dimensions(m, matched))
	}
}

// REQ-52: "An `exists` atom projects onto the scoped product as a
// **per-key presence dimension**: a two-valued dimension `{present,
// absent}` for that key, alongside the key's value dimension. An `exists`
// atom denotes `{present}` or `{absent}` on it"
// DOMAIN EDGE
func TestReq52_ExistsProjectsOntoAPerKeyPresenceDimension(t *testing.T) {
	m := mustLoadSource(t, existsPartitionSource())
	g := groupOf(t, m, "gate-present")

	if !slices.Contains(guard.PresenceDimensions(m, g), "gate") {
		t.Errorf("`gate` has no presence dimension; an `exists` atom projects "+
			"onto a two-valued {present, absent} dimension. presence=%v",
			guard.PresenceDimensions(m, g))
	}

	// The two-valued dimension sits ALONGSIDE the key's value dimension.
	if !slices.Contains(guard.Dimensions(m, g), "gate") {
		t.Errorf("`gate` has no value dimension; the presence dimension sits "+
			"alongside it, it does not replace it. dimensions=%v",
			guard.Dimensions(m, g))
	}

	present := guard.Denotation(m, "gate", table.Atom{
		Key: "gate", Block: table.BlockAll, Operator: "exists", Literal: []string{"true"},
	})
	absent := guard.Denotation(m, "gate", table.Atom{
		Key: "gate", Block: table.BlockAll, Operator: "exists", Literal: []string{"false"},
	})
	if present.Equal(absent) {
		t.Error("`exists = true` and `exists = false` denote the same subset; " +
			"they denote {present} and {absent}")
	}
}

// REQ-53: "a value atom denotes a subset of the key's value dimension **as
// the operator/kind agreement clause projects it** — one held value
// narrowed for `eq`/`in`/comparisons, which therefore requires the key be
// declared single-valued, or a containing subset for `contains` — and,
// because a value atom over an absent key is unevaluable rather than
// false, implicitly `{present}`."
// DOMAIN EDGE
func TestReq53_ValueAtomDenotesAValueSubsetAndImplicitlyPresent(t *testing.T) {
	m := mustLoadSource(t, optionalKeyGroupSource(false))

	d := guard.Denotation(m, "gate", table.Atom{
		Key: "gate", Block: table.BlockAll, Operator: "eq", Literal: []string{"true"},
	})
	if d.Presence() != guard.PresencePresent {
		t.Errorf("a value atom denotes presence %q; it is implicitly "+
			"{present}, because a value atom over an absent key is "+
			"unevaluable rather than false", d.Presence())
	}
	if d.Values() == nil {
		t.Error("a value atom denotes no value subset")
	}
}

// REQ-54: "`exists` is unaffected by that requirement: it reads presence
// alone, so it projects over a key whose values co-occur exactly as it does
// over a single-valued one."
// DOMAIN EDGE
func TestReq54_ExistsProjectsOverCoOccurringKeysExactlyAsOverSingleValued(t *testing.T) {
	single := mustLoadSource(t, existsOverKindSource(true))
	cooccur := mustLoadSource(t, existsOverKindSource(false))

	a := guard.Denotation(single, "gate", table.Atom{
		Key: "gate", Block: table.BlockAll, Operator: "exists", Literal: []string{"true"},
	})
	b := guard.Denotation(cooccur, "gate", table.Atom{
		Key: "gate", Block: table.BlockAll, Operator: "exists", Literal: []string{"true"},
	})
	if a.Presence() != b.Presence() {
		t.Errorf("`exists` projected differently over a single-valued key "+
			"(%q) and a co-occurring one (%q); it reads presence alone",
			a.Presence(), b.Presence())
	}

	// And it draws no unprojectable finding over the co-occurring key.
	if anyFinding(guard.Lint(cooccur), func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "gate"
	}) {
		t.Error("`exists` over a co-occurring key drew an unprojectable " +
			"finding; it is unaffected by the single-valued requirement")
	}

	// Discriminating leg: `exists` must PROJECT in both, and its projection
	// must be a real selection on the presence dimension. Two atoms that
	// project to nothing are "unaffected" trivially.
	for name, d := range map[string]guard.AssignmentSet{"single-valued": a, "co-occurring": b} {
		if !d.Projectable() {
			t.Errorf("`exists` over the %s key is not projectable; it reads "+
				"presence alone and always projects", name)
		}
		if d.Presence() != guard.PresencePresent {
			t.Errorf("`exists = true` over the %s key selects presence %q; "+
				"want %q", name, d.Presence(), guard.PresencePresent)
		}
	}
}

// REQ-55: "A key declared always-present contributes no presence dimension
// — its `{absent}` assignment is not in the product"
// BOUNDARY
func TestReq55_AlwaysPresentKeyContributesNoPresenceDimension(t *testing.T) {
	m := mustLoadSource(t, optionalKeyGroupSource(true))
	g := groupOf(t, m, "gate-open")

	if slices.Contains(guard.PresenceDimensions(m, g), "gate") {
		t.Errorf("an always-present key contributed a presence dimension; its "+
			"{absent} assignment is not in the product. presence=%v",
			guard.PresenceDimensions(m, g))
	}

	// And the cardinality is not doubled for it.
	card, ok := guard.Cardinality(m, g)
	if !ok {
		t.Fatal("the group's product has no cardinality")
	}
	if card != 2 {
		t.Errorf("cardinality = %d; want 2 — a single-valued always-present "+
			"bool contributes 2, with no ×2 presence factor", card)
	}
}

// REQ-56: "Lint MUST NOT drop `exists` atoms from the product: a group
// carrying one stays provable, and certifying it exhaustive while ignoring
// the presence dimension is the false-green this RDR's narrowing forbids."
// ADVERSARIAL
func TestReq56_ExistsAtomsAreNotDroppedFromTheProduct(t *testing.T) {
	// This group covers both VALUES of `gate` but only the {present} half
	// of its presence dimension. Dropping the `exists` atom would certify
	// it green.
	m := mustLoadSource(t, existsIgnoredWouldGoGreenSource())
	reports := guard.Lint(m)

	if hasGreen(reports) {
		t.Fatalf("a group leaving the {absent} presence assignment uncovered "+
			"certified green; the `exists` atom was dropped from the product. "+
			"reports=%s", renderReports(reports))
	}
	// It stays PROVABLE — this is a gap, not a withheld claim.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeCoverageGap
	}) {
		t.Errorf("no coverage gap finding; a group carrying an `exists` atom "+
			"stays provable. findings=%v", allFindings(reports))
	}
}

// REQ-57: "**`eq`, `in`, and the integer comparisons are single-value
// operators** … so an atom denotes the assignments in which the tag's
// single held value is the literal (`eq`), is a member of the literal set
// (`in`), or satisfies the comparison. `contains` is the operator for a tag
// whose values co-occur, and it projects the other way: the assignments
// whose held set contains every listed element."
// DOMAIN EDGE
func TestReq57_SingleValueOperatorsNarrowAndContainsProjectsTheOtherWay(t *testing.T) {
	m := mustLoadSource(t, setAndEnumSource())

	// `in` narrows the single held value to the listed set.
	in := guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockAll, Operator: "in", Literal: []string{"small", "mid"},
	})
	eq := guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"small"},
	})
	if !eq.Subset(in) {
		t.Error("`eq small` is not a subset of `in [small, mid]`; both narrow " +
			"the tag's SINGLE HELD VALUE")
	}

	// `contains` projects the other way: the assignments whose held set
	// contains every listed element. The FULL universe satisfies
	// `contains [bug]`, which a narrowing reading would exclude.
	contains := guard.Denotation(m, "labels", table.Atom{
		Key: "labels", Block: table.BlockAll, Operator: "contains", Literal: []string{"bug"},
	})
	full := guard.Assignment{"labels": `["bug","chore"]`}
	if !contains.Contains(full) {
		t.Errorf("`contains [bug]` excludes the held set %v; it denotes every "+
			"assignment whose held set CONTAINS every listed element", full)
	}
	only := guard.Assignment{"labels": `["chore"]`}
	if contains.Contains(only) {
		t.Errorf("`contains [bug]` includes the held set %v, which lacks bug", only)
	}
}

// REQ-58: "a value atom over a dimension the model does **not** declare
// single-valued is not a differently-projecting atom — it is one lint
// **cannot project at all** … Such an atom MUST take the blocking
// inability-to-prove outcome for that dimension, exactly as a dimension
// with no finite declared domain does."
// ADVERSARIAL
func TestReq58_ValueAtomOverAnUnmarkedDimensionIsUnprojectable(t *testing.T) {
	m := mustLoadSource(t, unmarkedEnumGuardSource())
	reports := guard.Lint(m)

	if hasGreen(reports) {
		t.Fatalf("a group whose `eq` atom sits over an unmarked dimension "+
			"certified green; lint cannot project it. reports=%s", renderReports(reports))
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "profile"
	}) {
		t.Errorf("no blocking inability-to-prove finding for the unmarked "+
			"dimension; findings=%v", allFindings(reports))
	}
	// "exactly as a dimension with no finite declared domain does": the SAME
	// code, not a distinct one.
	unbounded := guard.Lint(mustLoadSource(t, unboundedIntGuardSource()))
	if verdictOf(reports) != verdictOf(unbounded) {
		t.Errorf("the unmarked-dimension verdict %q differs from the "+
			"no-finite-domain verdict %q; the outcome is the same",
			verdictOf(reports), verdictOf(unbounded))
	}
}

// REQ-59: "Lint MUST NOT silently pick a reading — resolving it as \"the
// literal is among the held values\", \"the held set equals the literal\",
// or \"the held set is contained in the literal\" yields different union
// cardinalities and different overlap verdicts on the same model, which the
// published-bound clause forbids."
// ADVERSARIAL
func TestReq59_LintPicksNoReadingForAnUnprojectableAtom(t *testing.T) {
	m := mustLoadSource(t, unmarkedEnumGuardSource())
	g := groupOf(t, m, "unmarked-a")

	// No union is computed: picking any of the three readings would produce
	// one, and each produces a different cardinality.
	r := reportFor(t, guard.Lint(m), "unmarked-a")
	if r.CoverageUnion.Len() != 0 {
		t.Errorf("a union of %d assignments was computed over an "+
			"unprojectable atom; lint MUST NOT silently pick a reading",
			r.CoverageUnion.Len())
	}
	// And no overlap verdict is reached either.
	if countCode(guard.Lint(m), guard.CodeOverlap) != 0 {
		t.Error("an overlap verdict was reached over unprojectable atoms; " +
			"whether they overlap is the question the clause refuses to guess")
	}
	// The denotation itself reports that it is unprojectable.
	d := guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"},
	})
	if d.Projectable() {
		t.Errorf("the atom over unmarked dimension %v reports as projectable",
			guard.Dimensions(m, g))
	}
}

// REQ-60: "Coverage is `union(row_i accepted assignments) == scoped
// product` for that row group, and overlap is any non-empty `row_i
// accepted assignments intersect row_j accepted assignments`."
// HAPPY PATH
func TestReq60_CoverageIsUnionEqualsProductAndOverlapIsNonEmptyIntersection(t *testing.T) {
	m := mustLoadSource(t, completePartitionSource())
	g := groupOf(t, m, "part-small")

	product := guard.Product(m, g)
	union := guard.CoverageUnion(m, g)
	if !union.Equal(product) {
		t.Errorf("union (%d assignments) != scoped product (%d); coverage is "+
			"exactly that identity", union.Len(), product.Len())
	}

	a := guard.AcceptedAssignments(m, rowByID(t, m, "part-small"))
	b := guard.AcceptedAssignments(m, rowByID(t, m, "part-large"))
	if a.Intersect(b).Len() != 0 {
		t.Errorf("the partition's two rows intersect in %d assignments; a "+
			"partition is disjoint", a.Intersect(b).Len())
	}

	// The overlap fixture's rows DO intersect, and that is what the finding
	// reports.
	dup := mustLoadSource(t, twoRuleIdenticalGuardSource())
	x := guard.AcceptedAssignments(dup, rowByID(t, dup, "dup-a"))
	y := guard.AcceptedAssignments(dup, rowByID(t, dup, "dup-b"))
	if x.Intersect(y).Len() == 0 {
		t.Error("two byte-identical guards have an empty intersection")
	}
}

// REQ-61: "Lint MAY claim guard exhaustiveness only for finite declared
// domains: enum values, booleans, declared set element universes, or
// bounded integer ranges."
// BOUNDARY
func TestReq61_ExhaustivenessIsClaimableOnlyOverFiniteDeclaredDomains(t *testing.T) {
	min0, max2 := 0, 2
	finite := []table.TagDecl{
		{Kind: "enum", Domain: []string{"a", "b"}, SingleValued: true, Required: true},
		{Kind: "bool", SingleValued: true, Required: true},
		{Kind: "set", Elements: []string{"x"}, Required: true},
		{Kind: "int", Min: &min0, Max: &max2, SingleValued: true, Required: true},
	}
	for _, d := range finite {
		if _, ok := guard.AssignmentCount(d); !ok {
			t.Errorf("%+v carries no finite domain; enum values, booleans, "+
				"set element universes and bounded int ranges are finite", d)
		}
	}

	nonFinite := []table.TagDecl{
		{Kind: "scalar", Required: true},
		{Kind: "int", Required: true},
		{Kind: "enum", Required: true},
		{Kind: "set", Required: true},
	}
	for _, d := range nonFinite {
		if n, ok := guard.AssignmentCount(d); ok {
			t.Errorf("%+v reported a finite count of %d; it declares no "+
				"finite domain", d, n)
		}
	}
}

// REQ-62: "If a guard dimension lacks a finite declared domain, lint MUST
// refuse or downgrade an exhaustiveness claim for that dimension rather
// than treating the covered examples as complete."
// ADVERSARIAL
func TestReq62_NoFiniteDomainRefusesRatherThanTreatingExamplesAsComplete(t *testing.T) {
	// The rows enumerate every value the guards mention, so a lint treating
	// covered examples as complete would certify this green.
	m := mustLoadSource(t, undeclaredDomainExamplesSource())
	reports := guard.Lint(m)

	if hasGreen(reports) {
		t.Fatalf("a group over an undeclared domain certified green from its "+
			"covered examples; examples are not a domain. reports=%s",
			renderReports(reports))
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "loose"
	}) {
		t.Errorf("no refusal for the dimension lacking a finite declared "+
			"domain; findings=%v", allFindings(reports))
	}
}

// REQ-63: "The exhaustiveness claim is **default-on for every scoped row
// group whose participating dimensions are all finitely declared**: an
// opt-in flag would let the guarantee be silently skipped exactly where it
// matters."
// HAPPY PATH
func TestReq63_ClaimIsDefaultOnWithNoOptInGesture(t *testing.T) {
	// The fixture carries no opt-in of any kind, yet the gap is reported.
	m := mustLoadSource(t, productOnlyGapSource())
	if len(allFindings(guard.Lint(m))) == 0 {
		t.Fatal("a group with all dimensions finitely declared produced no " +
			"finding without an opt-in; the claim is default-on")
	}

	// And the clean partition is CLAIMED (green), not merely silent — which
	// is what distinguishes default-on from unchecked.
	clean := guard.Lint(mustLoadSource(t, completePartitionSource()))
	if !hasGreen(clean) {
		t.Errorf("a fully-declared clean group produced no claim at all; the "+
			"claim is default-on. reports=%s", renderReports(clean))
	}
}

// REQ-64: "**Leaving a dimension undeclared is not an opt-out.** An
// undeclared dimension does not remove the row group from proof; it makes
// the proof unavailable, which is the Loud blocking inability-to-prove
// outcome … never a silent downgrade to unchecked. No authoring gesture
// quietly exempts a group: a group is either proved, or it carries a
// blocking finding naming the dimension that cannot be proved."
// ADVERSARIAL
func TestReq64_UndeclaredDimensionIsNotAnOptOutAndEveryGroupIsProvedOrBlocked(t *testing.T) {
	m := mustLoadSource(t, undeclaredDomainExamplesSource())
	reports := guard.Lint(m)

	// The group is not REMOVED from proof: lint reports on it, so a lint
	// that simply skips undeclared groups fails here rather than passing
	// the every-group check vacuously.
	if len(reports) == 0 {
		t.Fatal("lint produced no group report for a model with an undeclared " +
			"dimension; leaving a dimension undeclared removed the group from " +
			"proof, which is exactly the opt-out this clause forbids")
	}
	if len(allFindings(reports)) == 0 {
		t.Fatal("lint produced no finding at all; the outcome is the LOUD " +
			"blocking inability-to-prove one, never a silent downgrade")
	}

	for _, r := range reports {
		if r.Green {
			continue
		}
		var blocking []guard.Finding
		for _, f := range r.Findings {
			if f.Blocking {
				blocking = append(blocking, f)
			}
		}
		if len(blocking) == 0 {
			t.Errorf("group %s is neither proved nor carrying a blocking "+
				"finding; there is no third, silently-unchecked state",
				r.Context)
		}
		for _, f := range blocking {
			if f.Code == guard.CodeUnprovableCoverage && f.Dimension == "" {
				t.Errorf("blocking finding %v names no dimension; it must name "+
					"the dimension that cannot be proved", f)
			}
		}
	}
}

// REQ-65: "\"Refuse\" and \"downgrade\" are one outcome, not an author's
// choice: every case this RDR sends to refuse-or-downgrade — a non-finite
// dimension, a finite product too large to prove deterministically, and a
// withheld claim under the narrowing above — MUST produce a blocking
// inability-to-prove finding. This RDR MUST NOT mint a non-blocking warning
// category for these cases"
// BOUNDARY
func TestReq65_AllThreeRefuseOrDowngradeCasesAreBlocking(t *testing.T) {
	cases := map[string]string{
		"non-finite dimension":  unboundedIntGuardSource(),
		"product too large":     overLargeProductSource(),
		"withheld under narrow": optionalKeyGroupSource(false),
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			reports := guard.Lint(mustLoadSource(t, src))
			var blocking int
			for _, f := range allFindings(reports) {
				if f.Blocking {
					blocking++
				}
			}
			if blocking == 0 {
				t.Errorf("produced no blocking finding; findings=%v",
					allFindings(reports))
			}
			if hasGreen(reports) {
				t.Error("certified green")
			}
		})
	}

	// No non-blocking warning category is minted for these cases.
	for _, code := range guard.Codes() {
		if !guard.IsBlocking(code) && guard.IsRefuseOrDowngrade(code) {
			t.Errorf("code %q is a non-blocking refuse-or-downgrade carrier; "+
				"this RDR mints no non-blocking warning category for these", code)
		}
	}
}

// REQ-66: "\"Downgrade\" therefore names the same blocking outcome as
// \"refuse\" on both sides of the seam, and never a silent or advisory
// one."
// BOUNDARY
func TestReq66_DowngradeAndRefuseNameOneBlockingOutcome(t *testing.T) {
	// The two blocking carriers this RDR's assumptions name.
	for _, code := range []guard.Code{guard.CodeUnprovableCoverage, guard.CodeProductTooLarge} {
		if !guard.IsBlocking(code) {
			t.Errorf("code %q is not blocking; downgrade names the same "+
				"blocking outcome as refuse", code)
		}
	}
	// And there is no advisory tier at all in this RDR's own vocabulary.
	for _, sev := range guard.Severities() {
		if sev == "warning" || sev == "advisory" || sev == "info" {
			t.Errorf("severity %q exists; downgrade is never silent or "+
				"advisory", sev)
		}
	}
}
