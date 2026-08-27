package guard_test

// RDR 0003 — polarity and `unless`, and the runtime-veto narrowing this
// RDR is the recording document for.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-38: "Positive guard atoms MUST live in `all`; negative guard atoms
// MUST live in `unless`. Successful row matching MUST NOT depend on source
// order or first-match priority."
// HAPPY PATH
func TestReq38_PolarityIsByBlockAndMatchingIsOrderIndependent(t *testing.T) {
	m := mustLoadSource(t, mixedBlockSource())

	row := rowByID(t, m, "mixed")
	var sawAll, sawUnless bool
	for _, a := range row.Atoms {
		switch a.Block {
		case table.BlockAll:
			sawAll = true
		case table.BlockUnless:
			sawUnless = true
		}
	}
	if !sawAll || !sawUnless {
		t.Fatalf("row %q lost a block: all=%v unless=%v — positive atoms live "+
			"in `all` and negative in `unless`", row.RuleID, sawAll, sawUnless)
	}

	// Matching does not depend on source order: reversing the rows of a
	// two-row group leaves the disposition unchanged.
	kt := m.KernelTable()
	forward := resolveWith(t, kt, guard.View{"profile": "large", "iter": "1"})
	slices.Reverse(kt.Rows)
	reversed := resolveWith(t, kt, guard.View{"profile": "large", "iter": "1"})
	if !sameDisposition(forward, reversed) {
		t.Errorf("reversing row order changed the disposition (%v vs %v); "+
			"matching MUST NOT depend on source order or first-match priority",
			describe(forward), describe(reversed))
	}
	// "Successful row matching": the disposition must be a SUCCESS in both
	// orders. Two identical refusals agree without matching anything.
	if forward.Refused() {
		t.Fatalf("the polarity fixture did not match successfully: %v; source "+
			"order independence is asserted over SUCCESSFUL matching",
			describe(forward))
	}

	// And the polarity is load-bearing: the same atoms in the OTHER blocks
	// select differently, so `all` and `unless` are not interchangeable.
	if guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"},
	}).Len() == 0 {
		t.Error("a positive `all` atom denotes nothing; positive guard atoms " +
			"live in `all` and must select there")
	}
}

// REQ-39: "a row's accepted assignments are the intersection of all
// positive `all` atom domains minus the single conjunctive assignment set
// matched by the row's full `unless` block. `unless` is not per-atom
// negation, and it does not create source-order priority."
// DOMAIN EDGE
func TestReq39_AcceptedAssignmentsAreAllIntersectionMinusTheUnlessConjunction(t *testing.T) {
	m := mustLoadSource(t, twoAtomUnlessSource())

	accepted := guard.AcceptedAssignments(m, rowByID(t, m, "two-atom-unless"))

	// `all: profile in {small,large}` intersected, minus the SINGLE
	// conjunctive set `profile == large AND flag == true`. Per-atom negation
	// would instead subtract every assignment matching EITHER unless atom,
	// removing `profile=large, flag=false` too.
	kept := guard.Assignment{"profile": "large", "flag": "false"}
	if !accepted.Contains(kept) {
		t.Errorf("accepted set does not contain %v; `unless` is not per-atom "+
			"negation — only the FULL conjunction is subtracted. accepted=%v",
			kept, accepted)
	}
	removed := guard.Assignment{"profile": "large", "flag": "true"}
	if accepted.Contains(removed) {
		t.Errorf("accepted set contains %v; the full `unless` conjunction is "+
			"subtracted", removed)
	}
}

// REQ-40: "a row qualifies only when every `all` atom is true and the
// `unless` predicate set is not fully true; if multiple rows qualify, RDR
// 0001's exact-one resolver refuses instead of choosing by priority."
// DOMAIN EDGE
func TestReq40_RowQualifiesOnAllTrueAndUnlessNotFullyTrueAndTiesRefuse(t *testing.T) {
	m := mustLoadSource(t, twoAtomUnlessSource())
	kt := m.KernelTable()

	// `all` true, `unless` not FULLY true: qualifies.
	res := resolveWith(t, kt, guard.View{"profile": "large", "flag": "false"})
	if res.Refused() {
		t.Errorf("a row with every `all` atom true and `unless` not fully "+
			"true refused %v; it qualifies", res.Refusal.Kind)
	}

	// `unless` fully true: disabled.
	res = resolveWith(t, kt, guard.View{"profile": "large", "flag": "true"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindNoMatch {
		t.Errorf("a row whose full `unless` block is true gave %+v; it is "+
			"disabled", describe(res))
	}

	// Multiple rows qualify: the exact-one resolver refuses, never chooses.
	dup := mustLoadSource(t, twoRuleIdenticalGuardSource())
	res = resolveWith(t, dup.KernelTable(), guard.View{"profile": "large"})
	if !res.Refused() || res.Refusal.Kind != resolve.KindAmbiguousMatch {
		t.Errorf("two qualifying rows gave %v; the exact-one resolver refuses "+
			"instead of choosing by priority", describe(res))
	}
}

// REQ-41: "**`unless` is subtracted two-valued only when its atoms are
// decided.** … an **unevaluable atom inside `unless` makes the whole row
// unevaluable**, not merely un-excluded."
// ADVERSARIAL
func TestReq41_UnevaluableAtomInsideUnlessMakesTheWholeRowUnevaluable(t *testing.T) {
	tbl := resolve.Table{
		Revision: "rev", Outcomes: []string{"go"},
		Rows: []resolve.Row{{
			RuleID: "unless-unevaluable", SourceLocator: "f:1", Outcome: "go",
			Guard: []resolve.GuardAtom{
				kernelAtom("profile", "eq", "large"),
				{Key: "missing", Operator: "eq", Literal: "x", Block: resolve.BlockUnless},
			},
		}},
	}
	res, err := resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go",
		Observed: []resolve.Tag{{Key: "profile", Value: "large"}},
		Guards:   guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("an unevaluable atom inside `unless` gave %v; ¬U = U makes "+
			"the WHOLE ROW unevaluable, not merely un-excluded", describe(res))
	}

	// Discriminating leg: "`unless` is subtracted two-valued only when its
	// atoms are DECIDED". With the same key present the atom decides and the
	// row is merely un-excluded, so the refusal above is attributable to the
	// unevaluable atom rather than to a seam that decides nothing.
	tbl.Rows[0].Guard[1].Key = "decidable"
	res, err = resolve.Resolve(resolve.Input{
		Table: tbl, Recognized: "go",
		Observed: []resolve.Tag{
			{Key: "profile", Value: "large"},
			{Key: "decidable", Value: "not-x"},
		},
		Guards: guard.Evaluator{},
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if res.Refused() {
		t.Fatalf("a DECIDED `unless` atom that does not hold refused %v; the "+
			"block is subtracted two-valued and the row stays selectable",
			describe(res))
	}
}

// REQ-42: "Lint MUST therefore treat a value atom in an `unless` block over
// a key declared optional exactly as it treats one in `all`: the row **can
// refuse**, and the group's claim is withheld under the narrowing above."
// DOMAIN EDGE
func TestReq42_OptionalKeyValueAtomInUnlessWithholdsExactlyAsInAll(t *testing.T) {
	inAll := guard.Lint(mustLoadSource(t, optionalAtomInBlockSource(table.BlockAll)))
	inUnless := guard.Lint(mustLoadSource(t, optionalAtomInBlockSource(table.BlockUnless)))

	for name, reports := range map[string][]guard.GroupReport{"all": inAll, "unless": inUnless} {
		if hasGreen(reports) {
			t.Errorf("the %s-block group certified green; a value atom over an "+
				"optional key means the row can refuse", name)
		}
		if !anyFinding(reports, func(f guard.Finding) bool {
			return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "gate"
		}) {
			t.Errorf("the %s-block group emitted no withholding finding; "+
				"findings=%v", name, allFindings(reports))
		}
	}

	// "exactly as it treats one in `all`": the two verdicts are identical.
	if verdictOf(inAll) != verdictOf(inUnless) {
		t.Errorf("the `all` group verdict %q differs from the `unless` group "+
			"verdict %q; lint MUST treat them identically",
			verdictOf(inAll), verdictOf(inUnless))
	}
}

// REQ-43: "**A can-refuse row contributes no assignments to the coverage
// union**, and the group it sits in has no provable product. … lint MUST
// NOT credit it with its `all`-intersection unsubtracted"
// ADVERSARIAL
func TestReq43_CanRefuseRowContributesNoAssignmentsAndTheGroupIsUnprovable(t *testing.T) {
	m := mustLoadSource(t, optionalAtomInBlockSource(table.BlockAll))
	reports := guard.Lint(m)

	r := reportFor(t, reports, "gate-open")
	if r.Provable {
		t.Error("a group containing a can-refuse row reported a provable " +
			"product; it has none")
	}
	if r.CoverageUnion.Len() != 0 {
		t.Errorf("the coverage union carries %d assignments; a can-refuse row "+
			"contributes none", r.CoverageUnion.Len())
	}
	if r.Green {
		t.Error("a group with a can-refuse row certified green; lint credited " +
			"it with its `all`-intersection unsubtracted")
	}
}

// REQ-44: "a withheld group emits the withholding finding for each
// refusing row, and MUST NOT additionally emit a `graph-coverage-gap`
// naming a witness assignment that a refusing row would in fact accept
// when its key is present. Overlap findings among the group's decidable
// rows are unaffected and still MUST be emitted."
// ADVERSARIAL
func TestReq44_WithheldGroupEmitsNoCoverageGapButStillEmitsOverlap(t *testing.T) {
	m := mustLoadSource(t, withheldPlusOverlapSource())
	reports := guard.Lint(m)

	if n := countCode(reports, guard.CodeCoverageGap); n != 0 {
		t.Errorf("a withheld group emitted %d graph-coverage-gap findings; a "+
			"withheld group has no provable product to gap over", n)
	}

	// One withholding finding per refusing row.
	withholding := findingsWithCode(reports, guard.CodeUnprovableCoverage)
	named := map[string]bool{}
	for _, f := range withholding {
		for _, id := range f.RuleIDs {
			named[id] = true
		}
	}
	if !named["refuser-a"] || !named["refuser-b"] {
		t.Errorf("withholding findings name %v; want one per refusing row "+
			"(refuser-a, refuser-b)", named)
	}

	// Overlap among the group's DECIDABLE rows is unaffected.
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			slices.Contains(f.RuleIDs, "decidable-a") &&
			slices.Contains(f.RuleIDs, "decidable-b")
	}) {
		t.Errorf("overlap among the group's decidable rows was suppressed by "+
			"the withholding; findings=%v", allFindings(reports))
	}
}

// REQ-67: "An exhaustiveness claim MUST NOT be stronger than the runtime it
// describes: lint MUST NOT certify a row group exhaustive when a
// participating row can refuse `guard_unevaluable` under RDR 0007's
// aggregation veto. Where the two disagree the lint promise narrows; the
// runtime veto MUST NOT be weakened."
// ADVERSARIAL
func TestReq67_LintNarrowsWhereItWouldOutrunTheRuntimeVeto(t *testing.T) {
	m := mustLoadSource(t, optionalKeyGroupSource(false))

	// The group is domain-exhaustive over `gate`'s two values, so a lint
	// reading declarations alone would certify it. It must not.
	reports := guard.Lint(m)
	if hasGreen(reports) {
		t.Fatalf("a domain-exhaustive group whose row can refuse "+
			"guard_unevaluable certified green; the claim MUST narrow. reports=%s",
			renderReports(reports))
	}

	// And the runtime veto is not weakened: the same model still refuses.
	res := resolveWith(t, m.KernelTable(), guard.View{})
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("the runtime gave %v; the veto MUST NOT be weakened to make "+
			"lint's claim true", describe(res))
	}

	// Discriminating leg: the narrowing must be exactly where the two
	// DISAGREE. A lint that certifies nothing narrows nowhere, so the
	// always-present twin — which the runtime never refuses — must go green.
	twin := mustLoadSource(t, optionalKeyGroupSource(true))
	if !hasGreen(guard.Lint(twin)) {
		t.Fatalf("the always-present twin did not certify green; a lint that "+
			"certifies nothing is not a narrowing. reports=%s",
			renderReports(guard.Lint(twin)))
	}
	if twinRes := resolveWith(t, twin.KernelTable(), guard.View{"gate": "true"}); twinRes.Refused() {
		t.Errorf("the runtime refused the twin lint certifies (%v); a green "+
			"result must mean resolution succeeds", describe(twinRes))
	}
}

// REQ-68: "**\"Participating row\" here means every row in the group,
// escape rows included** — the participation clause's population, not the
// ordinary-row population the overlap check uses."
// ADVERSARIAL
func TestReq68_NarrowingPopulationIncludesEscapeRows(t *testing.T) {
	// The only possibly-refusing row in this group is a GUARDED ESCAPE row.
	// Reading the overlap check's ordinary-row population into the narrowing
	// would certify this group green.
	m := mustLoadSource(t, guardedEscapeRefuserSource())
	reports := guard.Lint(m)

	if hasGreen(reports) {
		t.Fatalf("a group whose only refusing row is an escape row certified "+
			"green; the narrowing quantifies over EVERY row in the group. reports=%s",
			renderReports(reports))
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage &&
			slices.Contains(f.RuleIDs, "guarded-escape")
	}) {
		t.Errorf("no withholding finding named the guarded escape row; "+
			"findings=%v", allFindings(reports))
	}
}

// REQ-69: "Withholding is therefore decided over the whole group, while
// overlap remains split into the two populations above … The two-population
// reading applies to overlap only; reading it into the narrowing would
// certify green exactly the group the runtime refuses."
// DOMAIN EDGE
func TestReq69_WithholdingIsWholeGroupWhileOverlapStaysTwoPopulation(t *testing.T) {
	m := mustLoadSource(t, guardedEscapeRefuserSource())
	reports := guard.Lint(m)

	// Whole-group withholding: the escape row withholds the group, and it
	// does so by EMITTING the withholding finding rather than by lint being
	// silent about the group altogether.
	if hasGreen(reports) {
		t.Error("withholding was decided over the ordinary-row population only")
	}
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage &&
			slices.Contains(f.RuleIDs, "guarded-escape")
	}) {
		t.Errorf("no withholding finding named the guarded escape row; "+
			"whole-group withholding is an emitted result, not a silence. "+
			"findings=%v", allFindings(reports))
	}

	// Two-population overlap: the escape row overlapping an ORDINARY row is
	// not reported as an ambiguity.
	for _, f := range findingsWithCode(reports, guard.CodeOverlap) {
		hasEscape := slices.Contains(f.RuleIDs, "guarded-escape")
		hasOrdinary := slices.Contains(f.RuleIDs, "ordinary")
		if hasEscape && hasOrdinary {
			t.Errorf("overlap finding %v pairs an escape row with an ordinary "+
				"row; the runtime never matches the two populations", f.RuleIDs)
		}
	}
}

// REQ-70: "A row \"can refuse\" when the row carries a value atom over a
// key **declared optional** **in either block** — `all` or `unless` … the
// same one declared field the optionality clause defines, not a second
// presence property and not a graph query."
// BOUNDARY
func TestReq70_CanRefuseIsAValueAtomOverAnOptionalKeyInEitherBlock(t *testing.T) {
	cases := []struct {
		name  string
		row   table.Row
		decls map[string]table.TagDecl
		want  bool
	}{
		{
			name: "value atom over an optional key in all",
			row:  guardRow("r", table.Atom{Key: "gate", Block: table.BlockAll, Operator: "eq", Literal: []string{"true"}}),
			want: true,
		},
		{
			name: "value atom over an optional key in unless",
			row:  guardRow("r", table.Atom{Key: "gate", Block: table.BlockUnless, Operator: "eq", Literal: []string{"true"}}),
			want: true,
		},
		{
			name: "existence atom over the same optional key",
			row:  guardRow("r", table.Atom{Key: "gate", Block: table.BlockAll, Operator: "exists", Literal: []string{"true"}}),
			want: false,
		},
		{
			name: "value atom over an always-present key",
			row:  guardRow("r", table.Atom{Key: "fixed", Block: table.BlockAll, Operator: "eq", Literal: []string{"true"}}),
			want: false,
		},
		{
			name: "match atom over an optional key",
			row:  guardRow("r", table.Atom{Key: "gate", Block: table.BlockMatch, Operator: "eq", Literal: []string{"true"}}),
			want: false,
		},
	}

	decls := map[string]table.TagDecl{
		"gate":  {Kind: "bool", SingleValued: true},
		"fixed": {Kind: "bool", SingleValued: true, Required: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := guard.CanRefuse(decls, tc.row); got != tc.want {
				t.Errorf("CanRefuse = %v; want %v", got, tc.want)
			}
		})
	}
}

// REQ-71: "Lint MUST decide this syntactically over declarations so the
// test is total; it MUST NOT withhold a claim merely because some
// assignment in the product is unreached."
// ADVERSARIAL
func TestReq71_CanRefuseIsSyntacticOverDeclarationsAndTotal(t *testing.T) {
	decls := map[string]table.TagDecl{"gate": {Kind: "bool", SingleValued: true}}

	// Total: every row gets an answer, including one with no atoms at all.
	if guard.CanRefuse(decls, table.Row{RuleID: "bare"}) {
		t.Error("a row carrying no atoms can refuse; the syntactic test is " +
			"over the row's own value atoms")
	}
	// Total over an UNDECLARED key too — no panic, no reachability query.
	_ = guard.CanRefuse(decls, guardRow("r", table.Atom{
		Key: "nowhere", Block: table.BlockAll, Operator: "eq", Literal: []string{"x"},
	}))

	// An unreached assignment does not withhold: the complete partition has
	// an unreachable-looking branch but every key is always-present, so it
	// certifies green.
	reports := guard.Lint(mustLoadSource(t, completePartitionSource()))
	if !hasGreen(reports) {
		t.Errorf("a group over always-present keys was withheld; a claim MUST "+
			"NOT be withheld merely because some assignment is unreached. reports=%s",
			renderReports(reports))
	}
}

// REQ-72: "An **owned** tag carries a second, graph-level presence
// condition … but that condition governs whether a row may *match*, not
// whether its guard can refuse … Presence for the withholding decision
// reads exactly one field"
// ADVERSARIAL
func TestReq72_WithholdingReadsTheOptionalityFieldNotOwnedReachability(t *testing.T) {
	// Two declarations identical but for PROVENANCE. Owned-ness must not
	// change the withholding decision — only `required` does.
	owned := map[string]table.TagDecl{
		"gate": {Provenance: table.ProvenanceOwned, Kind: "bool", SingleValued: true, Required: true},
	}
	observed := map[string]table.TagDecl{
		"gate": {Provenance: table.ProvenanceObserved, Kind: "bool", SingleValued: true, Required: true},
	}
	row := guardRow("r", table.Atom{Key: "gate", Block: table.BlockAll, Operator: "eq", Literal: []string{"true"}})

	if guard.CanRefuse(owned, row) {
		t.Error("an always-present OWNED tag was read as can-refuse; presence " +
			"for the withholding decision reads exactly the optionality field")
	}
	if guard.CanRefuse(owned, row) != guard.CanRefuse(observed, row) {
		t.Error("provenance changed the withholding decision; it reads " +
			"exactly one field, and provenance is not it")
	}

	// Discriminating leg: flipping the ONE field it does read must flip the
	// answer, in both provenances. A CanRefuse that always answers false
	// passes the two checks above while reading no field at all.
	optionalOwned := map[string]table.TagDecl{
		"gate": {Provenance: table.ProvenanceOwned, Kind: "bool", SingleValued: true},
	}
	optionalObserved := map[string]table.TagDecl{
		"gate": {Provenance: table.ProvenanceObserved, Kind: "bool", SingleValued: true},
	}
	if !guard.CanRefuse(optionalOwned, row) {
		t.Error("an OPTIONAL owned tag was not read as can-refuse; the one " +
			"field the decision reads is the optionality marker")
	}
	if !guard.CanRefuse(optionalObserved, row) {
		t.Error("an OPTIONAL observed tag was not read as can-refuse")
	}
}

// REQ-73: "A withheld exhaustiveness claim MUST be observable, not silent.
// It takes the same blocking inability-to-prove form an unprovable
// dimension already takes — RDR 0006's `graph-unprovable-coverage` — and
// MUST name the participating row and the atom that can refuse, using the
// source rule/context id every other predicate diagnostic names."
// HAPPY PATH
func TestReq73_WithheldClaimNamesTheRowTheAtomAndTheSourceID(t *testing.T) {
	m := mustLoadSource(t, optionalKeyGroupSource(false))
	reports := guard.Lint(m)

	var found []guard.Finding
	for _, f := range findingsWithCode(reports, guard.CodeUnprovableCoverage) {
		if slices.Contains(f.RuleIDs, "gate-open") {
			found = append(found, f)
		}
	}
	if len(found) == 0 {
		t.Fatalf("no withholding finding named the participating row; "+
			"findings=%v", allFindings(reports))
	}
	for _, f := range found {
		if f.Atom == nil {
			t.Errorf("finding %v names no refusing atom; the clause requires "+
				"the atom that can refuse", f)
			continue
		}
		if f.Atom.Key != "gate" {
			t.Errorf("finding names atom over %q; want the refusing atom over "+
				"`gate`", f.Atom.Key)
		}
		if f.Context == "" {
			t.Error("finding carries no source rule/context id")
		}
		if !f.Blocking {
			t.Error("the withheld claim is not blocking; it takes the same " +
				"blocking inability-to-prove form")
		}
	}
}

// REQ-74: "Emitting nothing MUST NOT satisfy this clause: an exit code
// alone cannot distinguish a withheld claim from a proved one."
// ADVERSARIAL
func TestReq74_EmittingNothingDoesNotSatisfyTheWithholdingClause(t *testing.T) {
	withheld := guard.Lint(mustLoadSource(t, optionalKeyGroupSource(false)))

	if len(allFindings(withheld)) == 0 {
		t.Fatal("the withheld group emitted NO findings; emitting nothing " +
			"MUST NOT satisfy the clause — a run that emits nothing fails")
	}

	// And the distinction is visible: the proved twin emits a green verdict
	// the withheld one does not.
	proved := guard.Lint(mustLoadSource(t, optionalKeyGroupSource(true)))
	if hasGreen(withheld) == hasGreen(proved) {
		t.Errorf("a withheld claim and a proved one are indistinguishable "+
			"(green=%v both); an exit code alone cannot tell them apart",
			hasGreen(proved))
	}
}

// REQ-103: "Value atom **inside `unless`** over an absent key" → runtime
// "`guard_unevaluable` refusal — `¬U = U`, so the row is unevaluable, not
// merely un-excluded"; lint "Exhaustiveness claim **withheld** for that
// group; the block is **not** subtracted as if decided"; Loud on both
// surfaces.
// DOMAIN EDGE
func TestReq103_UnlessOverAbsentKeyIsLoudOnBothSurfaces(t *testing.T) {
	// Runtime surface.
	tbl := resolve.Table{
		Revision: "rev", Outcomes: []string{"go"},
		Rows: []resolve.Row{{
			RuleID: "r", SourceLocator: "f:1", Outcome: "go",
			Guard: []resolve.GuardAtom{
				{Key: "gate", Operator: "eq", Literal: "true", Block: resolve.BlockUnless},
			},
		}},
	}
	res, err := resolve.Resolve(resolve.Input{Table: tbl, Recognized: "go", Guards: guard.Evaluator{}})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("runtime gave %v; want guard_unevaluable — ¬U = U", describe(res))
	}

	// Lint surface: withheld, and the block is not subtracted as if decided.
	reports := guard.Lint(mustLoadSource(t, optionalAtomInBlockSource(table.BlockUnless)))
	if hasGreen(reports) {
		t.Error("lint certified green; the claim is withheld and the block " +
			"is not subtracted as if decided")
	}
	if len(allFindings(reports)) == 0 {
		t.Error("lint emitted nothing; the outcome is Loud on both surfaces")
	}
}

// REQ-104: "Value atom over an absent key" → runtime `guard_unevaluable`;
// lint "Exhaustiveness claim **withheld** for that group"; diagnostic "the
// blocking inability-to-prove finding, naming the row and the refusing
// atom"; Loud.
// DOMAIN EDGE
func TestReq104_ValueAtomOverAbsentKeyIsWithheldWithRowAndAtomNamed(t *testing.T) {
	reports := guard.Lint(mustLoadSource(t, optionalKeyGroupSource(false)))

	if hasGreen(reports) {
		t.Fatal("lint certified green over a possibly-absent guard key")
	}
	blocking := findingsWithCode(reports, guard.CodeUnprovableCoverage)
	if len(blocking) == 0 {
		t.Fatalf("no blocking inability-to-prove finding; findings=%v", allFindings(reports))
	}
	for _, f := range blocking {
		if len(f.RuleIDs) == 0 || f.Atom == nil {
			t.Errorf("finding %v names row=%v atom=%v; the diagnostic names "+
				"the row AND the refusing atom", f, f.RuleIDs, f.Atom)
		}
	}
}

// REQ-105: "Existence atom over an absent key" → "Decided (`presence ==
// literal`) — never unevaluable"; lint "Selects `{absent}` on the presence
// dimension (A7)"; no diagnostic.
// DOMAIN EDGE
func TestReq105_ExistenceOverAbsentKeySelectsAbsentAndDrawsNoDiagnostic(t *testing.T) {
	m := mustLoadSource(t, existsPartitionSource())

	// Lint: the `exists = false` atom selects {absent} on the presence
	// dimension, so the two rows partition presence and the group proves.
	reports := guard.Lint(m)
	if !hasGreen(reports) {
		t.Errorf("an exists-partitioned group did not certify green; the "+
			"`exists = false` atom selects {absent}. reports=%s", renderReports(reports))
	}
	if n := len(allFindings(reports)); n != 0 {
		t.Errorf("lint emitted %d findings for an exists partition; the "+
			"disposition row says no diagnostic. findings=%v", n, allFindings(reports))
	}

	// Runtime: decided, never unevaluable.
	res := resolveWith(t, m.KernelTable(), guard.View{})
	if res.Refused() && res.Refusal.Kind == resolve.KindGuardUnevaluable {
		t.Error("an existence atom over an absent key was unevaluable; it is " +
			"DECIDED from presence == literal")
	}
}

// REQ-110: "Row group is domain-exhaustive but a participating row can
// refuse" → "Claim withheld — the narrowing; one finding per refusing row,
// never just the first"; "RDR 0007 payload at runtime **and** a blocking
// lint finding naming the refusing atom — absence of a green result is not
// the artifact".
// ADVERSARIAL
func TestReq110_OneFindingPerRefusingRowAndAbsenceOfGreenIsNotTheArtifact(t *testing.T) {
	m := mustLoadSource(t, twoRefusingRowsSource())
	reports := guard.Lint(m)

	named := map[string]bool{}
	for _, f := range findingsWithCode(reports, guard.CodeUnprovableCoverage) {
		for _, id := range f.RuleIDs {
			named[id] = true
		}
	}
	if !named["refuser-a"] || !named["refuser-b"] {
		t.Errorf("withholding findings name %v; one finding per refusing row, "+
			"never just the first", named)
	}

	// "absence of a green result is not the artifact": a positive finding
	// must exist.
	if len(allFindings(reports)) == 0 {
		t.Fatal("lint emitted nothing; the artifact is the blocking finding, " +
			"not the absence of a green result")
	}

	// And the runtime carries RDR 0007's payload naming the atom.
	res := resolveWith(t, m.KernelTable(), guard.View{})
	if !res.Refused() || res.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Fatalf("runtime gave %v; want the RDR 0007 payload", describe(res))
	}
	if len(res.Refusal.Undecided) == 0 {
		t.Error("the runtime refusal carries no per-atom payload")
	}
}

// REQ-111: "`all` atom decides false" → "Row pruned … Silent by design — a
// decided false is not a defect"; "Full `unless` block decides true" →
// "Row disabled … Excluded intersection subtracted … Silent by design".
// DOMAIN EDGE
func TestReq111_DecidedFalseAndFullUnlessTrueAreSilentByDesign(t *testing.T) {
	m := mustLoadSource(t, completePartitionSource())
	reports := guard.Lint(m)

	// The partition's rows prune each other at runtime, and lint says
	// nothing about it: the group proves clean.
	if !hasGreen(reports) {
		t.Fatalf("the partition did not prove; reports=%s", renderReports(reports))
	}
	if n := len(allFindings(reports)); n != 0 {
		t.Errorf("lint emitted %d findings for a clean partition; a decided "+
			"false is not a defect. findings=%v", n, allFindings(reports))
	}

	// A row disabled by a fully-true `unless` block is likewise silent: the
	// excluded intersection is subtracted, and no finding is minted for it.
	sub := guard.Lint(mustLoadSource(t, twoAtomUnlessSource()))
	for _, f := range allFindings(sub) {
		if f.Code == guard.CodeOverlap || f.Code == guard.CodeUnprovableCoverage {
			t.Errorf("a subtracted `unless` block drew finding %v; the "+
				"subtraction is silent by design", f)
		}
	}
}
