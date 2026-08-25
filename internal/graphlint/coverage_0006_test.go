package graphlint_test

// RDR 0006 — row grouping (REQ-19..REQ-24), determinism and overlap
// (REQ-25..REQ-30), guard exhaustiveness and the withholding rules
// (REQ-38, REQ-39, REQ-45..REQ-66), and the precedence clause (REQ-115).

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/guard"
)

// REQ-19: "Coverage, overlap, and withholding MUST be decided per scoped
// row group as RDR 0003 defines it … This RDR MUST NOT define a second
// grouping predicate; it supplies only which selection contexts are
// reachable"
// BOUNDARY
func TestReq19_GroupingIsRDR0003sAndNotASecondPredicate(t *testing.T) {
	m := mustLoad(t, source(twoStateDecls, legalBody))

	// The oracle is agreement with the landed peer: whatever this RDR
	// decides per group must be decided over RDR 0003's groups, so a row
	// sits in exactly the group `guard.Groups` puts it in. Defining a
	// second predicate would put some row somewhere else.
	groups := guard.Groups(m)
	if len(groups) == 0 {
		t.Fatal("RDR 0003's grouping produced no groups for the legal fixture")
	}

	// `advance-on` and `advance-off` share a match pattern and an outcome,
	// so RDR 0003 puts them in ONE group.
	var together bool
	for _, g := range groups {
		var ids []string
		for _, row := range g.Rows {
			ids = append(ids, row.RuleID)
		}
		if slices.Contains(ids, "advance-on") && slices.Contains(ids, "advance-off") {
			together = true
		}
	}
	if !together {
		t.Fatalf("fixture precondition failed: RDR 0003 does not group the " +
			"two rows together, so this test cannot detect a second predicate")
	}

	// This RDR's coverage verdict must therefore be a per-group verdict
	// over that same partition: the two rows partition the `flag`
	// dimension, so the group closes and no gap is reported.
	r := graphlint.Run(graphlint.NewRequest(m))
	requireNoCode(t, r, graphlint.CodeCoverageGap)

	// The discriminating half, so an engine deciding nothing at all cannot
	// satisfy the clause: dropping one row leaves the SAME group short of
	// the product, and the per-group verdict flips. Both fixtures share
	// one grouping; only the group's contents differ.
	const shortBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	short := mustLoad(t, source(twoStateDecls, shortBody))
	if len(guard.Groups(short)) != 1 {
		t.Fatalf("fixture precondition failed: the short model has %d groups, "+
			"want 1", len(guard.Groups(short)))
	}
	requireCode(t, graphlint.Run(graphlint.NewRequest(short)),
		graphlint.CodeCoverageGap)
}

// REQ-20: "**Source state is the authored match pattern.** … rows sharing
// a match pattern and a recognized outcome are one group."
// REQ-21: "Two rows whose patterns differ are two groups even when some
// abstract node satisfies both"
// REQ-52: "Lint MUST NOT widen a group's row set by pooling rows
// satisfiable at a shared abstract node"
// REQ-124 / SC-20: "**no run may exit 0 with an empty blocking list**"
// ADVERSARIAL
func TestReq20And21And52_GroupMembershipIsTheAuthoredMatchPatternNotANode(t *testing.T) {
	// SC-20's fixture: two rows with DIFFERENT match patterns
	// (`status = a`, `status = b`), each individually non-exhaustive over
	// the shared `flag` dimension, whose accepted assignments are
	// complementary so a POOLED union would close. Two paths write `a` and
	// `b`, converging on a node whose `status` set is {a, b}, so both rows
	// are satisfiable there — the pooling trap.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seed", "a", "b", "end"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "seed"
flag = "false"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "seed-to-a"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"

[[rule]]
id = "seed-to-b"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "b"

[[rule]]
id = "a-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "end"

[[rule]]
id = "b-off"
[rule.match.status]
eq = "b"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "end"
`
	r := lint(t, decls, body)

	// The negative control is the point of the fixture.
	if len(r.Blocking()) == 0 {
		t.Fatalf("the two-different-match-pattern fixture exited with an "+
			"empty blocking list — exactly the false green a merged-node "+
			"union produces. Membership is the authored match pattern, so "+
			"the rows never pool; report:%s", render(r))
	}
	// Two groups, and a coverage gap on each.
	if n := countCode(r, graphlint.CodeCoverageGap); n < 2 {
		t.Errorf("%d %s findings; want one per group — the rows are two "+
			"groups, each individually non-exhaustive; report:%s",
			n, graphlint.CodeCoverageGap, render(r))
	}

	// The paired POSITIVE control: give both rows the SAME match pattern
	// and they are one group whose union genuinely closes — lints clean.
	const sameBody = `
terminal = ["done"]

[initial]
status = "seed"
flag = "false"

[context.done]
[context.done.match.status]
eq = "end"

[[rule]]
id = "seed-to-a"
[rule.match.status]
eq = "seed"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"

[[rule]]
id = "a-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "end"

[[rule]]
id = "a-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "end"
`
	same := lint(t, decls, sameBody)
	requireNoCode(t, same, graphlint.CodeCoverageGap)
}

// REQ-22: "The reachability relation below decides **which** groups are
// proven, never **which rows are in one**. Membership is syntactic and
// per-row; reachability is a filter over contexts."
// BOUNDARY
func TestReq22_ReachabilityFiltersContextsAndNeverChangesMembership(t *testing.T) {
	// A group whose selection context NO reachable node satisfies is
	// filtered out of the proof obligation — but its rows stay in it. The
	// observable consequence: the unreachable group takes the advisory
	// `graph-unreachable-rule`, not a coverage gap, and its rows are never
	// folded into a reachable group.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "orphan"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"

[[rule]]
id = "orphaned"
[rule.match.status]
eq = "orphan"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	r := lint(t, decls, body)

	// No reachable node satisfies `status = orphan`, so the group is not
	// proven — the filter's job — and the row is reported unreachable.
	if !namesRule(r, graphlint.CodeUnreachableRule, "orphaned") {
		t.Errorf("no %s finding names the row whose context no reachable "+
			"node satisfies; report:%s",
			graphlint.CodeUnreachableRule, render(r))
	}
	// Membership is unchanged by the filter: the orphan row is NOT pooled
	// into the reachable `status = a` group, so that group still closes.
	requireNoCode(t, r, graphlint.CodeCoverageGap)
}

// REQ-23: "invariants 3 and 4 run once per reachable group."
// REQ-24: "Group *membership* is syntactic (the authored match pattern),
// so the group count is bounded by the authored rule count regardless of
// the lattice; only the reachability filter scales with nodes."
// BOUNDARY
func TestReq23And24_OverlapAndCoverageRunOncePerReachableGroup(t *testing.T) {
	// One reachable group carrying an overlapping ordinary pair. Running
	// the check per NODE rather than per group would multiply the finding
	// by the reachable node count; running it once per group yields one.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c", "d"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true
`
	// Several nodes reach `status = a` so the per-node/per-group
	// distinction is observable — and they reach it holding DIFFERENT
	// `flag` values, so the two arrivals are distinct concrete views
	// rather than the same owned-state twice. Letting both paths preserve
	// `flag = false` (the earlier shape) made a per-node implementation
	// and a per-group one indistinguishable: there was only ever one node.
	const body = `
terminal = ["done"]

[initial]
status = "c"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "c-to-a"
[rule.match.status]
eq = "c"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"

[[rule]]
id = "c-to-d"
[rule.match.status]
eq = "c"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "d"

[[rule]]
id = "d-to-a"
[rule.match.status]
eq = "d"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"
flag = "true"

[[rule]]
id = "overlap-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "overlap-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, decls, body)

	// One overlapping pair in one group is one finding, however many
	// nodes reach the group's context.
	if n := countCode(r, graphlint.CodeOverlap); n != 1 {
		t.Errorf("%d %s findings; want exactly 1 — invariant 3 runs once per "+
			"reachable GROUP, not once per reachable node; report:%s",
			n, graphlint.CodeOverlap, render(r))
	}
}

// REQ-25: "Graph lint MUST reject ambiguity instead of relying on source
// order, rendered-row order, or first-match priority to choose between
// enabled rows."
// REQ-26: "when two ordinary rows qualify for the same selection context,
// lint rejects the model. It never selects by order"
// REQ-27: Invariant 3 — "within a scoped row group, no finite-domain input
// assignment may enable two ordinary (non-escape) rows."
// HAPPY PATH
func TestReq25And26And27_TwoOrdinaryRowsEnabledTogetherAreRejected(t *testing.T) {
	// Two ordinary rows in one group with no guard at all: every
	// assignment enables both.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "first"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "second"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, body)
	f := requireCode(t, r, graphlint.CodeOverlap)

	if !graphlint.IsBlocking(graphlint.CodeOverlap) {
		t.Error("graph-overlap is not blocking; lint must REJECT ambiguity")
	}
	for _, got := range f {
		if got.Severity != graphlint.SeverityBlocking {
			t.Errorf("%s severity = %q; want %q", got.Code, got.Severity,
				graphlint.SeverityBlocking)
		}
	}

	// Never selects by order: reversing the authored order changes neither
	// the verdict nor the finding count, so no first-match priority is in
	// play.
	const reversed = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "second"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "first"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	rev := lint(t, statusOnlyDecls, reversed)
	if countCode(rev, graphlint.CodeOverlap) != countCode(r, graphlint.CodeOverlap) {
		t.Errorf("reversing the authored row order changed the overlap "+
			"finding count (%d -> %d); lint must never rely on source order",
			countCode(r, graphlint.CodeOverlap),
			countCode(rev, graphlint.CodeOverlap))
	}
}

// REQ-28: "an escape row overlapping an ordinary row is not a runtime
// ambiguity and is not reported."
// SC-9a; the disposition table's "Silent by design" row.
// ADVERSARIAL
func TestReq28_EscapeOverlappingAnOrdinaryRowIsNotReported(t *testing.T) {
	// One ordinary row and one escape row whose accepted assignments
	// overlap it completely. The kernel matches an escape row only to
	// rescue a refusal, so this is never a runtime ambiguity — reporting
	// it would fail a model the runtime accepts.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "ordinary"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, statusOnlyDecls, body)
	requireNoCode(t, r, graphlint.CodeOverlap)

	// The paired POSITIVE control, so an engine that never reports overlap
	// cannot satisfy the negative half: replace the ordinary row with a
	// second ESCAPE row sharing the class and the same assignments, and
	// the overlap IS reported. The two fixtures differ only in the first
	// row's escape declaration — exactly the two-population distinction.
	const escapePair = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "ordinary"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"

[[rule]]
id = "rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	requireCode(t, lint(t, statusOnlyDecls, escapePair), graphlint.CodeOverlap)
}

// REQ-29: "Escape rows MUST participate in the coverage union and MUST be
// overlap-checked in one population per declared failure class, never
// against ordinary rows"
// REQ-30: "One overlapping pair MUST yield one finding naming both rows;
// escape rows overlapping across several shared failure classes MUST yield
// one finding per shared class."
// SC-9b, SC-13.
// BOUNDARY
func TestReq29And30_EscapePairsOverlapPerSharedClassOneFindingEach(t *testing.T) {
	// Two escape rows sharing BOTH declared classes. The per-class
	// partition yields one finding per shared class — two here — each
	// naming both rows.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "ordinary"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "esc-one"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"

[[rule]]
id = "esc-two"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, statusOnlyDecls, body)
	f := requireCode(t, r, graphlint.CodeOverlap)

	if len(f) != 2 {
		t.Errorf("%d %s findings; want one per SHARED failure class (2); "+
			"report:%s", len(f), graphlint.CodeOverlap, render(r))
	}
	// Each finding is scoped to an escape population, so it carries the
	// failure class — and REQ-30's "one finding NAMING BOTH ROWS" means
	// each per-class finding names esc-one AND esc-two. A finding that
	// named only one row would satisfy a count-plus-nonempty-class check
	// while telling the author half the story.
	classes := map[string]bool{}
	for _, got := range f {
		if got.Class == "" {
			t.Errorf("%s finding scoped to an escape population carries no "+
				"failure class: %+v", got.Code, got)
		}
		classes[got.Class] = true

		// The pair is named in the IDENTITY fields, and is asserted as an
		// exact unordered pair. Substring-matching a blob of concatenated
		// fields would accept a prefixed id (`esc-one-b`) and would also
		// let unrelated message prose stand in for an identity, so the
		// two identity fields are compared directly.
		//
		// This also carries the ordinary-row exclusion: `ordinary` shares
		// the same match pattern but is overlap-checked in a separate
		// population, and an exact pair equal to {esc-one, esc-two}
		// cannot name it.
		pair := []string{got.Rule, got.Element}
		slices.Sort(pair)
		if wantPair := []string{"esc-one", "esc-two"}; !slices.Equal(
			pair, wantPair) {
			t.Errorf("the %s finding for class %q names the row pair %v; "+
				"want exactly %v — one overlapping pair is one finding "+
				"naming BOTH escape rows, and escape rows are "+
				"overlap-checked in one population per declared class, "+
				"NEVER against ordinary rows: %+v",
				got.Code, got.Class, pair, wantPair, got)
		}
	}
	for _, want := range guard.RescuableClasses() {
		if !classes[want] {
			t.Errorf("no %s finding for the shared class %q; report:%s",
				graphlint.CodeOverlap, want, render(r))
		}
	}
}

// REQ-38: Invariant 4 — "every scoped row group whose participating guard
// dimensions are all finitely declared claims closed coverage (default-on,
// never an opt-in annotation), and lint must prove
// `union(row_i accepted assignments) == scoped product`"
// REQ-47: "The claim is default-on for every scoped row group whose
// participating dimensions are all finitely declared."
// HAPPY PATH
func TestReq38And47_CoverageClaimIsDefaultOnAndProvedAsUnionEqualsProduct(t *testing.T) {
	// The claim is DEFAULT-ON: no annotation opts the group in, and a
	// group leaving an assignment uncovered is reported without one.
	const gapBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	gap := lint(t, twoStateDecls, gapBody)
	requireCode(t, gap, graphlint.CodeCoverageGap)

	// The union closing over the product is the pass condition, and the
	// legal fixture's two rows partition `flag`.
	closed := lint(t, twoStateDecls, legalBody)
	requireNoCode(t, closed, graphlint.CodeCoverageGap)
}

// REQ-39: "Escape rows contribute their accepted assignments to the union
// like any other row; there is no separate \"an escape row exists\"
// disjunct."
// DOMAIN EDGE
func TestReq39_EscapeRowsContributeToTheUnionLikeAnyOtherRow(t *testing.T) {
	// The ordinary row covers `flag = true` only. A GUARDED escape row
	// covering `flag = false` closes the union — by contributing its
	// accepted assignments, not by existing. The discriminating control
	// below shows a guarded escape covering the SAME assignment leaves the
	// gap open, which an "escape row exists" disjunct could not do.
	const closes = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "rescue-off"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
`
	requireNoCode(t, lint(t, twoStateDecls, closes), graphlint.CodeCoverageGap)

	// The control: an escape row covering the assignment the ordinary row
	// ALREADY covers leaves `flag = false` uncovered. An "escape row
	// exists" disjunct would close the group here too.
	const doesNotClose = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "rescue-on"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
`
	requireCode(t, lint(t, twoStateDecls, doesNotClose), graphlint.CodeCoverageGap)
}

// REQ-45: "Graph lint MAY claim exhaustiveness only over finite declared
// domains … and MUST read finite domain, optionality, and
// single-valuedness from that declaration"
// REQ-46: "If a participating dimension is not finite or not projectable,
// lint MUST emit `graph-unprovable-coverage` for that dimension."
// REQ-80 (dimension arm): the finding carries reason `dimension-not-finite`
// naming the dimension.
// DOMAIN EDGE
func TestReq45And46_NonFiniteDimensionTakesUnprovableCoverage(t *testing.T) {
	// `free` is a `scalar`: no finite declared domain, so the dimension is
	// not projectable and the claim cannot be made.
	const body = `
terminal = ["done"]

[initial]
status = "a"
free = "anything"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-free"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.free]
eq = "x"
[rule.write]
status = "b"
`
	r := lint(t, nonFiniteDecls, body)
	f := requireCode(t, r, graphlint.CodeUnprovableCoverage)

	var named bool
	for _, got := range f {
		if got.Reason == graphlint.ReasonDimensionNotFinite && got.Dimension == "free" {
			named = true
		}
	}
	if !named {
		t.Errorf("no %s finding carries reason %q naming the dimension "+
			"`free`; report:%s", graphlint.CodeUnprovableCoverage,
			graphlint.ReasonDimensionNotFinite, render(r))
	}
}

// REQ-48: "The scoped product MUST include the presence dimension an
// `exists` atom over an optional key contributes … lint MUST NOT certify a
// group exhaustive while ignoring that dimension's `{absent}` assignment."
// REQ-49 / REQ-123: "A key declared always-present contributes no presence
// dimension" — and the control's key must carry the EXPLICIT marker.
// SC-11.
// BOUNDARY
func TestReq48And49_PresenceDimensionIsInTheProductForOptionalKeysOnly(t *testing.T) {
	// The optional-key group's rows cover every VALUE assignment but leave
	// the `{absent}` assignment uncovered, so dropping the presence
	// dimension is the false green this asserts against.
	const optionalBody = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "present"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
exists = true
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, optionalBody)
	requireCode(t, r, graphlint.CodeCoverageGap)

	// The CONTROL: the same shape over a key carrying RDR 0003's EXPLICIT
	// always-present marker contributes no presence dimension, so the
	// group closes. An unmarked key defaults to optional and would take
	// the x2 presence row, which is why the marker must be explicit.
	const alwaysBody = `
terminal = ["done"]

[initial]
status = "a"
always = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "present"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.always]
exists = true
[rule.write]
status = "b"
`
	ctrl := lint(t, alwaysPresentOwnedDecls, alwaysBody)
	requireNoCode(t, ctrl, graphlint.CodeCoverageGap)
}

// REQ-50: "Coverage is a universal claim, so it MUST be computed from the
// group's authored rows and their declared guard domains alone — never
// from a reachability node."
// REQ-53: "Reachability decides only whether a group is proven at all; it
// never enters the coverage computation."
// ADVERSARIAL
func TestReq50And53_CoverageReadsAuthoredRowsAndNeverAReachabilityNode(t *testing.T) {
	// The group's own rows leave `flag = false` uncovered, but the ONLY
	// reachable node for the group's context holds `flag = true`. A
	// coverage computation consulting the node would narrow the product
	// and certify the group; the authored-rows-only reading reports the
	// gap.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "true"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	r := lint(t, twoStateDecls, body)
	requireCode(t, r, graphlint.CodeCoverageGap)
}

// REQ-51: "The union … and the scoped product both range over the group's
// *participating guard dimensions* …; match keys are not product
// dimensions … and contribute no assignment to either side."
// BOUNDARY
func TestReq51_MatchKeysAreNotProductDimensions(t *testing.T) {
	// `status` is a MATCH key and `flag` is the only guard dimension. If
	// match keys entered the product, the group's two rows could not close
	// it (they say nothing about the other `status` values) and a spurious
	// gap would be reported.
	r := lint(t, twoStateDecls, legalBody)
	requireNoCode(t, r, graphlint.CodeCoverageGap)

	// The discriminating control: the same group with the `flag` dimension
	// half-covered DOES report the gap, so the assertion above is not
	// simply "this engine never reports gaps".
	const halfBody = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	requireCode(t, lint(t, twoStateDecls, halfBody), graphlint.CodeCoverageGap)
}

// REQ-54: "lint MUST NOT certify a row group exhaustive when a
// participating row can refuse `guard_unevaluable`. A withheld claim MUST
// be emitted as `graph-unprovable-coverage`, naming the participating row
// and the refusing atom."
// REQ-56: the "can refuse" test is "a syntactic test over the declared
// optionality field, not a reachability query".
// REQ-81 / REQ-80: reason `row-can-refuse`, and the message must read as an
// honest withholding, not an authoring error.
// SC-8.
// HAPPY PATH
func TestReq54And56_WithheldClaimNamesTheRowAndTheRefusingAtom(t *testing.T) {
	// `opt` is declared OPTIONAL, so a VALUE atom over it can refuse
	// `guard_unevaluable` — a syntactic property of the declaration.
	const body = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, body)
	f := requireCode(t, r, graphlint.CodeUnprovableCoverage)

	var ok bool
	for _, got := range f {
		if got.Reason != graphlint.ReasonRowCanRefuse {
			continue
		}
		if got.Rule != "can-refuse" {
			continue
		}
		// The atom fields are required on this arm.
		if got.Key != "opt" || got.Operator != "eq" ||
			got.Literal == "" || got.Block == "" {
			t.Errorf("the withheld-claim finding does not carry the refusing "+
				"atom's four fields: key=%q operator=%q literal=%q block=%q",
				got.Key, got.Operator, got.Literal, got.Block)
			continue
		}
		ok = true
	}
	if !ok {
		t.Errorf("no %s finding carries reason %q naming the row "+
			"`can-refuse` and its refusing atom; report:%s",
			graphlint.CodeUnprovableCoverage,
			graphlint.ReasonRowCanRefuse, render(r))
	}

	// SC-8's second half: no `graph-coverage-gap` is emitted for the same
	// group — the claim is withheld, not failed.
	requireNoCode(t, r, graphlint.CodeCoverageGap)

	// REQ-56's "not a reachability query": an UNLESS-block atom over the
	// same optional key withholds identically, and the `unless` variant is
	// the fixture SC-8 names alongside the `all` one.
	const unlessBody = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse-unless"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.unless.opt]
eq = "q"
[rule.write]
status = "b"
`
	un := lint(t, optionalGuardDecls, unlessBody)
	if !namesRule(un, graphlint.CodeUnprovableCoverage, "can-refuse-unless") {
		t.Errorf("an `unless`-block atom over the same optional key did not "+
			"withhold the claim; the can-refuse test is syntactic over the "+
			"declared optionality field, not block-sensitive; report:%s",
			render(un))
	}
}

// REQ-55: "This RDR cites that clause and MUST NOT restate it, mint a
// second code for it, or carry it in a non-blocking tier."
// REQ-76: "the advisory tier … must not absorb a withheld claim."
// ADVERSARIAL
func TestReq55And76_WithheldClaimIsBlockingAndNotInTheAdvisoryTier(t *testing.T) {
	const body = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	r := lint(t, optionalGuardDecls, body)
	f := requireCode(t, r, graphlint.CodeUnprovableCoverage)

	for _, got := range f {
		if got.Severity != graphlint.SeverityBlocking {
			t.Errorf("%s severity = %q; there is no non-blocking tier for "+
				"this class", got.Code, got.Severity)
		}
	}
	// The withheld claim appears among the BLOCKING findings, never the
	// advisory ones.
	for _, got := range r.Advisory() {
		if got.Code == graphlint.CodeUnprovableCoverage {
			t.Errorf("the advisory tier absorbed a withheld claim: %+v", got)
		}
	}
	// No second code was minted for it.
	if !slices.Contains(graphlint.AdvisoryCodes(), graphlint.CodeCoverageClosedByEscape) {
		t.Errorf("the advisory tier %v does not carry the bare-escape "+
			"closure code, so the tier under test is not the one the clause "+
			"names", graphlint.AdvisoryCodes())
	}
	for _, c := range graphlint.AdvisoryCodes() {
		if c == graphlint.CodeUnprovableCoverage {
			t.Errorf("%q appears in the advisory tier", c)
		}
	}
}

// REQ-57: "Lint MUST publish the model-independent product bound above
// which it declines to prove coverage, and MUST emit
// `graph-product-too-large` for a group whose declared finite product
// exceeds it"
// REQ-59: "The bound is an implementation constant … not a per-model
// input."
// SC-18.
// BOUNDARY
func TestReq57And59_DeclaredFiniteProductOverTheBoundIsProductTooLarge(t *testing.T) {
	bound := graphlint.ProductBound()
	if bound <= 0 {
		t.Fatalf("the published product bound is %d; it must be a positive "+
			"implementation constant", bound)
	}

	// `n` is declared over a finite int domain far wider than the bound,
	// so the group's declared finite product exceeds it.
	const body = `
terminal = ["done"]

[initial]
status = "a"
n = 0

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "wide"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.n]
gte = 5
[rule.write]
status = "b"
`
	r := lint(t, wideIntDecls, body)
	f := requireCode(t, r, graphlint.CodeProductTooLarge)

	// The disposition table: the group-scoped finding names group + bound.
	// `Rule != "" || Element != ""` would NOT say that — the traversal-scoped
	// spelling of this same code satisfies the `Element` half (see
	// TestReq60_ReachableNodeSetOverTheCeilingIsProductTooLargeOnTheTraversal),
	// so the two would be indistinguishable. The group-scoped finding is the
	// one carrying a source rule id AND a selection-context element AND the
	// dimension whose declared domain blew the bound.
	var namesGroup bool
	for _, got := range f {
		if got.Rule != "wide" {
			continue
		}
		if got.Element == "" || got.Element == "traversal" {
			t.Errorf("the group-scoped %s finding carries element %q; it "+
				"must name the selection context, never the traversal: %+v",
				graphlint.CodeProductTooLarge, got.Element, got)
			continue
		}
		// Compare parsed tokens, never a substring: `Dimension` is a
		// comma-joined list and a bare Contains(..., "n") also matches
		// `recognized` and `unknown`.
		if !slices.Contains(strings.Split(got.Dimension, ","), "n") {
			t.Errorf("the group-scoped %s finding names dimensions %q; it "+
				"must name `n`, whose declared domain exceeds the bound: %+v",
				graphlint.CodeProductTooLarge, got.Dimension, got)
			continue
		}
		// The diagnostic must quote the PUBLISHED bound, not merely
		// report that some bound was exceeded: a message naming a wrong
		// or stale constant misdirects the narrowing it asks for.
		if want := fmt.Sprintf("%d", bound); !strings.Contains(got.Message, want) {
			t.Errorf("the group-scoped %s finding's message %q does not name "+
				"the published product bound %s: %+v",
				graphlint.CodeProductTooLarge, got.Message, want, got)
			continue
		}
		namesGroup = true
	}
	if !namesGroup {
		t.Errorf("no %s finding names the group by its source rule `wide`, "+
			"its selection context, and the offending dimension; report:%s",
			graphlint.CodeProductTooLarge, render(r))
	}

	// The bound is model-INDEPENDENT: a second model with a different
	// shape reports the same constant.
	if graphlint.ProductBound() != bound {
		t.Error("the published product bound varied between calls; it is an " +
			"implementation constant, not a per-model input")
	}
}

// REQ-58: "An atom over a tag not declared single-valued has no projection
// and MUST take `graph-unprovable-coverage` (`0003::A21`)."
// REQ-80 (tag-not-single-valued arm).
// DOMAIN EDGE
func TestReq58_ValueAtomOverANonSingleValuedTagIsUnprovable(t *testing.T) {
	// `multi` declares a finite domain but carries NO single-valued
	// marker, so an `eq` atom over it has no projection.
	const decls = statusOnlyDecls + `
[tags.multi]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"
multi = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-multi"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.multi]
eq = "p"
[rule.write]
status = "b"
`
	r := lint(t, decls, body)
	f := requireCode(t, r, graphlint.CodeUnprovableCoverage)

	var named bool
	for _, got := range f {
		if got.Reason == graphlint.ReasonTagNotSingleValued {
			named = true
			if got.Key != "multi" && got.Dimension != "multi" {
				t.Errorf("the %q finding names neither the key nor the "+
					"dimension `multi`: %+v", got.Reason, got)
			}
		}
	}
	if !named {
		t.Errorf("no %s finding carries reason %q; report:%s",
			graphlint.CodeUnprovableCoverage,
			graphlint.ReasonTagNotSingleValued, render(r))
	}
}

// REQ-60: "Lint MUST therefore publish a **model-independent node
// ceiling** alongside the product bound, and MUST emit
// `graph-product-too-large` naming the traversal (not a group) when a
// model's reachable node set exceeds it, rather than running unboundedly."
// SC-22.
// ADVERSARIAL
func TestReq60_ReachableNodeSetOverTheCeilingIsProductTooLargeOnTheTraversal(t *testing.T) {
	ceiling := graphlint.NodeCeiling()
	if ceiling <= 0 {
		t.Fatalf("the published node ceiling is %d; it must be a positive "+
			"implementation constant", ceiling)
	}

	decls, body := overCeilingFixture(ceilingKeys)
	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	// PRECONDITION, not a branch. The old shape wrapped this whole
	// assertion in `if len(Reach(...)) > ceiling` and the fixture produced
	// two nodes, so the body never ran and the ceiling check could have
	// been deleted outright with the suite still green. A fixture that
	// fails to exceed the ceiling is a BROKEN FIXTURE and must say so.
	if len(nodes) <= ceiling {
		t.Fatalf("the fixture reaches %d nodes, at or under the published "+
			"ceiling of %d, so the clause under test is unexercised; widen "+
			"the generated fixture rather than accepting a vacuous pass",
			len(nodes), ceiling)
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	f := requireCode(t, r, graphlint.CodeProductTooLarge)

	// The ceiling finding names the TRAVERSAL, not a group: it carries the
	// traversal element id and no rule id. The group-scoped spelling of
	// this same code carries both (see
	// TestReq57And59_DeclaredFiniteProductOverTheBoundIsProductTooLarge),
	// so accepting either would not tell the two apart.
	var namesTraversal bool
	for _, got := range f {
		if got.Rule != "" || got.Element != "traversal" {
			continue
		}
		// The diagnostic must quote the PUBLISHED ceiling. Without this
		// the message could name any number -- or none -- and the test
		// would still see a correctly scoped finding.
		if want := fmt.Sprintf("%d", ceiling); !strings.Contains(got.Message, want) {
			t.Errorf("the traversal-scoped %s finding's message %q does not "+
				"name the published node ceiling %s: %+v",
				graphlint.CodeProductTooLarge, got.Message, want, got)
			continue
		}
		namesTraversal = true
	}
	if !namesTraversal {
		t.Errorf("no %s finding names the traversal (rule=\"\", "+
			"element=%q); the node-ceiling finding is traversal-scoped, "+
			"never group-scoped; report:%s",
			graphlint.CodeProductTooLarge, "traversal", render(r))
	}

	// The ceiling is model-independent.
	if graphlint.NodeCeiling() != ceiling {
		t.Error("the published node ceiling varied between calls")
	}
}

// ceilingKeys is the number of clearable owned tags
// `overCeilingFixture` emits. The reachable owned-state set is the
// subset lattice over their presence footprints, so the node count is
// 2^ceilingKeys — 8192 against the published ceiling of 4096, the
// smallest power of two that clears it with margin.
const ceilingKeys = 13

// overCeilingFixture generates a model whose reachable node set exceeds
// the published node ceiling, and returns its declarations and body.
//
// Widening a declared int domain does NOT do this: `reach` folds
// successors into existing nodes and mints a node only for a written
// value, so a declared domain's width never drives the node count. What
// multiplies nodes is the PRESENCE FOOTPRINT — `successorsOf` groups
// successors by which owned tags they hold — so n independently clearable
// optional owned tags give 2^n distinct reachable owned-states.
func overCeilingFixture(keys int) (decls, body string) {
	var d, b strings.Builder

	d.WriteString(`
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`)
	for k := range keys {
		fmt.Fprintf(&d, `
[tags.k%d]
provenance = "owned"
kind = "enum"
domain = ["on"]
single_valued = true
`, k)
	}

	b.WriteString("\nterminal = [\"done\"]\n\n[initial]\nstatus = \"a\"\n")
	for k := range keys {
		fmt.Fprintf(&b, "k%d = \"on\"\n", k)
	}
	b.WriteString(`
[context.done]
[context.done.match.status]
eq = "b"
`)
	// One row per key, each clearing that key alone. Every subset of the
	// keys is therefore a reachable presence footprint.
	for k := range keys {
		fmt.Fprintf(&b, `
[[rule]]
id = "drop%d"
clear = ["k%d"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"
`, k, k)
	}
	return d.String(), b.String()
}

// REQ-61: "An escape row closes coverage only for the failure classes it
// declares."
// REQ-62: "Lint MUST therefore compute the coverage union per (scoped row
// group × declared rescuable class) — `no_match` and `ambiguous_match`
// only — and MUST NOT let a row declaring one class close the group's
// other arms."
// SC-19a.
// DOMAIN EDGE
func TestReq61And62_EscapeClosesOnlyTheClassesItDeclares(t *testing.T) {
	// The group's only escape row declares `no_match` alone, and the
	// group's ordinary rows OVERLAP — so the `ambiguous_match` arm is
	// reachable and its uncovered state must be reported alongside the
	// overlap. Letting the `no_match` row close it is the defect.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "rescue-no-match-only"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, twoStateDecls, body)

	// The group overlaps, so the ambiguous_match arm is reachable.
	requireCode(t, r, graphlint.CodeOverlap)
	// And its coverage is not closed by the no_match-only escape row.
	// REQ-61 is class-SCOPED: a bare `graph-coverage-gap` proves nothing,
	// because a gap on the `no_match` arm — the arm the escape row DOES
	// close — would satisfy it while demonstrating the opposite defect.
	//
	// The fixture has ONE group and ONE uncovered class, so the gap is
	// exactly one finding: counting `ambiguous_match` findings and only
	// requiring a nonzero count would green-pass an engine reporting the
	// same gap twice.
	gap := requireOneCode(t, r, graphlint.CodeCoverageGap)
	if gap.Class != "ambiguous_match" {
		t.Errorf("the %s finding carries class %q; want %q — the arm the "+
			"escape row does NOT declare is the one left open, and "+
			"letting the no_match row close it is the defect. The union "+
			"is computed per (group x declared rescuable class) over "+
			"`no_match` and `ambiguous_match` only: %+v",
			gap.Code, gap.Class, "ambiguous_match", gap)
	}
	// The rescuable class vocabulary is exactly the two RDR 0003 names.
	if got := guard.RescuableClasses(); len(got) != 2 {
		t.Errorf("rescuable classes = %v; the union is computed per (group x "+
			"declared rescuable class) over `no_match` and `ambiguous_match` "+
			"only", got)
	}
}

// REQ-63: "Lint MUST therefore check the `ambiguous_match` arm's coverage
// only for a group that carries a `graph-overlap` finding, and MUST treat
// the arm as vacuously closed for a group whose ordinary population is
// overlap-free."
// SC-19b; the disposition table's "Silent by design" row.
// ADVERSARIAL
func TestReq63_AmbiguousMatchArmIsVacuouslyClosedForAnOverlapFreeGroup(t *testing.T) {
	// The same no_match-only escape row, but the group's ordinary rows do
	// NOT overlap. An overlap-free group cannot produce `ambiguous_match`,
	// so demanding an escape row for that arm would require a row
	// `graph-unreachable-rule` then flags.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"

[[rule]]
id = "rescue-no-match-only"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, twoStateDecls, body)

	requireNoCode(t, r, graphlint.CodeOverlap)
	// The ambiguous_match arm is vacuously closed: no coverage finding.
	requireNoCode(t, r, graphlint.CodeCoverageGap)
	// SC-19: both variants emit the bare-escape closure for the no_match
	// arm — the escape row carries no guard atoms.
	requireCode(t, r, graphlint.CodeCoverageClosedByEscape)
}

// REQ-64: "`owned_state_unavailable` is a runtime refusal class no escape
// row rescues and no lint finding mints; invariant 6's
// `graph-owned-before-write` is the design-time check whose runtime
// counterpart it is, and the two MUST NOT be conflated."
// BOUNDARY
func TestReq64_OwnedStateUnavailableIsNeitherRescuableNorMinted(t *testing.T) {
	// No escape row rescues it: the rescuable vocabulary excludes it, and
	// RDR 0002 already refuses a rule declaring it.
	if slices.Contains(guard.RescuableClasses(), "owned_state_unavailable") {
		t.Errorf("rescuable classes %v include owned_state_unavailable",
			guard.RescuableClasses())
	}
	// No lint finding mints it: it appears in no code and no class field.
	for _, c := range append(graphlint.BlockingCodes(), graphlint.AdvisoryCodes()...) {
		if c == "graph-owned-state-unavailable" || c == "owned_state_unavailable" {
			t.Errorf("the taxonomy mints a code for owned_state_unavailable: %q", c)
		}
	}
	// The two are not conflated: the design-time check keeps its own code.
	if !slices.Contains(graphlint.BlockingCodes(), graphlint.CodeOwnedBeforeWrite) {
		t.Errorf("%s is absent from the blocking code set %v",
			graphlint.CodeOwnedBeforeWrite, graphlint.BlockingCodes())
	}
}

// REQ-65: "A group whose coverage is closed by a bare escape row MUST emit
// `graph-coverage-closed-by-escape` naming that row; a bare green MUST NOT
// satisfy this clause."
// REQ-66: the closure does not fail the run.
// SC-9c.
// HAPPY PATH
func TestReq65And66_BareEscapeClosureIsReportedAndDoesNotFailTheRun(t *testing.T) {
	// The ordinary row covers `flag = true` only; a BARE escape row (one
	// carrying no guard atoms) closes the rest.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, twoStateDecls, body)

	// A bare green must NOT satisfy the clause: the finding is required.
	f := requireCode(t, r, graphlint.CodeCoverageClosedByEscape)
	var named bool
	for _, got := range f {
		if got.Rule == "bare-rescue" {
			named = true
		}
		if got.Severity != graphlint.SeverityInfo {
			t.Errorf("%s severity = %q; want %q — the closure passes",
				got.Code, got.Severity, graphlint.SeverityInfo)
		}
	}
	if !named {
		t.Errorf("no %s finding names the bare escape row `bare-rescue`; "+
			"report:%s", graphlint.CodeCoverageClosedByEscape, render(r))
	}

	// The closure does not fail the run.
	requireClean(t, r)
}

// REQ-115: Desk-trace step 7 — "A group both closable by a bare escape row
// **and** carrying a row that can refuse resolves to
// `graph-unprovable-coverage` (blocking, exit 2) — **not** a success with
// `graph-coverage-closed-by-escape`. Withholding dominates closure"
// SC-12, whose negative half is the assertion that matters.
// ADVERSARIAL
func TestReq115_WithholdingDominatesBareEscapeClosure(t *testing.T) {
	// The conjunction: a bare escape row that WOULD close the group, and a
	// participating row carrying a value atom over an OPTIONAL key, so it
	// can refuse `guard_unevaluable`.
	const body = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "can-refuse"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, optionalGuardDecls, body)

	// Withholding wins: blocking.
	requireCode(t, r, graphlint.CodeUnprovableCoverage)
	if len(r.Blocking()) == 0 {
		t.Fatalf("the group resolved to a success; withholding dominates "+
			"closure; report:%s", render(r))
	}
	// The negative half is the assertion that matters.
	requireNoCode(t, r, graphlint.CodeCoverageClosedByEscape)
}

// REQ-116: "**There is no suppression or waiver mechanism** — no inline
// ignore comment, no allowlist."
// REQ-117: "the cure is a clearer model, never a weaker lint."
// ADVERSARIAL
func TestReq116And117_NoSuppressionOrWaiverChannelExists(t *testing.T) {
	// A model carrying every shape an author might use to ask for silence:
	// a `[model].metadata` table (the one sanctioned extension namespace,
	// carried through UNINTERPRETED) naming the finding it wants waived,
	// and a `source` annotation on the offending rule.
	const body = `
terminal = ["done"]

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "only-on"
source = "lint:ignore graph-coverage-gap"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`
	const waiverHeader = `
[model.metadata]
lint_ignore = ["graph-coverage-gap"]
allowlist = ["only-on"]
`
	src := source(twoStateDecls, waiverHeader+body)
	r := lintSource(t, src)

	// The finding stands: no metadata key and no source annotation
	// suppresses it.
	requireCode(t, r, graphlint.CodeCoverageGap)
}

// REQ-118: "This RDR introduces no encode/decode, import/export, or
// inverse operation."
// BOUNDARY
func TestReq118_NoRoundTripOrInverseSurfaceIsIntroduced(t *testing.T) {
	// Negative REQ. The engine's exported surface must carry no
	// encode/decode or import/export pair — parse/render fidelity stays
	// RDR 0002's.
	funcs := exportedFuncsIn(t, "../../internal/graphlint")
	for _, name := range funcs {
		for _, banned := range []string{
			"Encode", "Decode", "Marshal", "Unmarshal",
			"Import", "Export", "Parse", "Render", "Dump",
		} {
			if name == banned || (len(name) > len(banned) &&
				name[:len(banned)] == banned) {
				t.Errorf("the engine exports %q; this RDR introduces no "+
					"encode/decode, import/export, or inverse operation", name)
			}
		}
	}
}
