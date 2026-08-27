package graphlint_test

// RDR 0010 §E — ∅-rooted graph lint over the decision-table class
// (`0010:C5`), plus the negative REQs pinning what stays unchanged.
//
// Every fixture is authored TOML handed to the real `internal/table`
// loader and driven through the real `graphlint.Run`. Nothing mocks the
// engine, and per ASSUMPTION-7 the "MUST NOT emit" clauses are asserted as
// EXACT finding-list equality against the full taxonomy (SC-4), not as a
// per-code negative unit test at each silent site.

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- the decision-table fixture corpus -----------------------------------

// dtHeader0010 is the decision-table preamble: the class declaration, one
// recognized outcome, and two observed enum dimensions each `single_valued`
// and `required`. No `[initial]`, no accessors, no owned tag (`0010:MVV`
// step 1).
const dtHeader0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.b]
provenance = "observed"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true
`

// dtCell renders one ordinary decision-table rule: the outcome binding and
// the two GUARD atoms that discriminate the cell. Guard atoms, not match
// atoms, are what `internal/guard/product.go::Dimensions` collects
// (`0006:C7`, `0003:C13`).
func dtCell(id, a, b string) string {
	return "\n[[rule]]\nid = \"" + id + "\"\n" +
		"[rule.match.recognized]\neq = \"decide\"\n" +
		"[rule.guard.all.a]\neq = \"" + a + "\"\n" +
		"[rule.guard.all.b]\neq = \"" + b + "\"\n" +
		"[rule.emit]\nverdict = \"" + id + "\"\n"
}

// dtPartial0010 covers three of the four (a,b) cells — the MVV step-2
// table, whose coverage gap is the POSITIVE finding proving coverage ran
// with no root declared.
var dtPartial0010 = dtHeader0010 +
	dtCell("cell-xp", "x", "p") +
	dtCell("cell-xq", "x", "q") +
	dtCell("cell-yp", "y", "p")

// dtComplete0010 closes the product with a fourth ORDINARY rule.
var dtComplete0010 = dtPartial0010 + dtCell("cell-yq", "y", "q")

// dtEscape0010 replaces the fourth rule with an escape row rescuing
// `no_match`, itself carrying `[rule.emit]`.
var dtEscape0010 = dtPartial0010 + `
[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "fallback"
`

// dtMatchOnly0010 is SC-4's match-only negative control: the same table
// discriminating by `[rule.match.<key>]` instead of guard atoms, so the
// scoped product has ZERO participating dimensions.
const dtMatchOnly0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.match.a]
eq = "x"
[rule.emit]
verdict = "x"
`

// dtMatchOnlyBareEscape0010 is dtMatchOnly0010 plus a BARE escape row
// rescuing `no_match`. It is the A13 double-report probe: without the
// early `return`, a zero-dimension group carrying a bare escape row yields
// BOTH `graph-unprovable-coverage` and `graph-coverage-closed-by-escape`.
const dtMatchOnlyBareEscape0010 = dtMatchOnly0010 + `
[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
`

// dtTwoOutcomeEscape0010 is SC-4's two-outcome escape control: one escape
// row over a TWO-outcome alphabet. It rescues one outcome and leaves the
// other's coverage open, per `0002:C5`'s per-outcome rescue scoping (A9).
const dtTwoOutcomeEscape0010 = `outcomes = ["decide", "review"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "decide-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.emit]
verdict = "dx"

[[rule]]
id = "decide-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"

[[rule]]
id = "review-x"
[rule.match.recognized]
eq = "review"
[rule.guard.all.a]
eq = "x"
[rule.emit]
verdict = "rx"
`

// dtStrayTerminal0010 is SC-4's stray-`terminal` control: a decision table
// declaring `terminal` over an OBSERVED tag. Invariant 1's
// terminal-predicate arm stays live and reports it at
// `element = terminal[0]` (C2, C5) — the arm distinguished from the
// silenced root arm by ELEMENT, not by code.
//
// Over a decision table the arm cannot false-positive: with zero owned
// tags EVERY terminal predicate reads a non-owned tag.
var dtStrayTerminal0010 = "terminal = [\"stop\"]\n" + dtHeader0010 + `
[context.stop]
[context.stop.match.a]
eq = "x"
` +
	dtCell("cell-xp", "x", "p") +
	dtCell("cell-xq", "x", "q") +
	dtCell("cell-yp", "y", "p") +
	dtCell("cell-yq", "y", "q")

// smRootless0010 is the class-omitted control: the SAME shape with `class`
// omitted, carrying the owned-tag scaffolding `0002:C4`'s write-block arm
// demands of a state machine, and NO `[initial]`. It must take
// `0006:C18`'s missing-root finding at `element = model`, unchanged.
const smRootless0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seen"]
single_valued = true
required = true

[read.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"
read_back = true

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.write]
status = "seen"
`

// --- REQs ----------------------------------------------------------------

// REQ-52: "Graph lint over a `\"decision-table\"` model MUST root the
// reachability relation at the empty owned-state node; the reachable set is
// exactly that node, and every selection context is reachable."
// HAPPY PATH
func TestReq52_DecisionTableLintRootsAtTheEmptyOwnedStateNode(t *testing.T) {
	m := mustLoad(t, dtPartial0010)

	nodes := graphlint.Reach(m)
	if len(nodes) != 1 {
		t.Fatalf("Reach returned %d nodes; the reachable set over a "+
			"decision table is EXACTLY the empty owned-state node", len(nodes))
	}
	if len(nodes[0].Values) != 0 {
		t.Errorf("the root node holds %v; the ∅ node holds nothing",
			nodes[0].Values)
	}

	// Every selection context is reachable, so coverage RUNS: the partial
	// table's gap is the positive finding that proves it.
	r := graphlint.Run(graphlint.NewRequest(m))
	requireCode(t, r, graphlint.CodeCoverageGap)
}

// REQ-53: "The root is keyed on the declared class read from the loaded
// model, **augmenting** the existing `len(Initial)` test rather than
// replacing it: a state-machine model with an empty `[initial]` MUST still
// traverse nothing and take `0006:C18`'s missing-root finding, so the
// seeding predicate is \"decision-table class OR a non-empty declared
// root\", never \"class alone\"."
// ADVERSARIAL — the augmentation, taken in both directions.
func TestReq53_TheSeedingPredicateIsClassOrDeclaredRootNeverClassAlone(t *testing.T) {
	t.Run("decision table seeds at ∅", func(t *testing.T) {
		if n := len(graphlint.Reach(mustLoad(t, dtComplete0010))); n != 1 {
			t.Errorf("Reach = %d nodes; want 1", n)
		}
	})

	t.Run("rootless state machine still traverses nothing", func(t *testing.T) {
		m := mustLoad(t, smRootless0010)
		if n := len(graphlint.Reach(m)); n != 0 {
			t.Errorf("Reach = %d nodes over a rootless STATE MACHINE; the "+
				"class-keyed seed must AUGMENT the len(Initial) test, not "+
				"replace it", n)
		}

		f := requireOneCode(t, graphlint.Run(graphlint.NewRequest(m)),
			graphlint.CodeDanglingEdge)
		if f.Element != "model" {
			t.Errorf("missing-root finding element = %q; want %q "+
				"(`0006:C18` unchanged)", f.Element, "model")
		}
	})
}

// REQ-54: "Overlap (invariant 3), coverage and withholding (invariant 4),
// the redundant-row, vacuous-atom, and unreachable-rule advisories, the
// node ceiling, and the product bound MUST run unchanged over the authored
// **guard** dimensions"
// HAPPY PATH — coverage runs, and it runs over the GUARD dimensions.
func TestReq54_CoverageRunsOverTheAuthoredGuardDimensions(t *testing.T) {
	f := requireOneCode(t,
		graphlint.Run(graphlint.NewRequest(mustLoad(t, dtPartial0010))),
		graphlint.CodeCoverageGap)

	dims := strings.Split(f.Dimension, ",")
	slices.Sort(dims)
	if !slices.Equal(dims, []string{"a", "b"}) {
		t.Errorf("coverage dimensions = %q; want the two authored GUARD "+
			"dimensions a,b — not the tags' observed provenance", f.Dimension)
	}

	t.Run("overlap runs unchanged", func(t *testing.T) {
		overlapping := dtComplete0010 + dtCell("cell-yq-dup", "y", "q")
		r := graphlint.Run(graphlint.NewRequest(mustLoad(t, overlapping)))
		requireCode(t, r, graphlint.CodeOverlap)
	})
}

// REQ-55 / REQ-95 / SC-4: "A `\"decision-table\"` group whose scoped
// product has **zero participating dimensions** MUST take
// `graph-unprovable-coverage` naming the group, and MUST NOT be reported as
// covered." / "the match-only control reports exactly
// `graph-unprovable-coverage` naming the group at exit 2 **while
// incomplete** — a *positive* assertion"
// ADVERSARIAL — the false green C5 exists to close.
func TestReq55_AZeroDimensionDecisionTableGroupIsUnprovable(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtMatchOnly0010)))

	f := requireOneCode(t, r, graphlint.CodeUnprovableCoverage)
	if f.Element != "dt/decide" {
		t.Errorf("finding element = %q; the arm MUST name the GROUP", f.Element)
	}
	if !graphlint.IsBlocking(f.Code) {
		t.Error("the finding is not blocking; the group MUST NOT be " +
			"reported as covered")
	}
	if len(r.Blocking()) == 0 {
		t.Error("the report carries no blocking finding; a match-discriminated " +
			"table exiting clean is the vacuous coverage C5 closes")
	}
}

// REQ-56 / REQ-57: "This RDR therefore requires a **fourth reason value**,
// `no-participating-dimension`, appended to that set." / "`reason` is
// normative for this fence: an emission carrying `\"\"` fails REQ-80
// outright, and one borrowing a dimension-scoped member states a remedy
// that does not apply."
// BOUNDARY
func TestReq56_TheZeroDimensionArmCarriesTheFourthReasonValue(t *testing.T) {
	t.Run("the closed set gains the fourth member", func(t *testing.T) {
		got := slices.Clone(graphlint.Reasons())
		slices.Sort(got)
		want := []string{
			"dimension-not-finite", "no-participating-dimension",
			"row-can-refuse", "tag-not-single-valued",
		}
		if !slices.Equal(got, want) {
			t.Errorf("reason set = %v; want exactly %v — `0006`'s set is "+
				"closed and APPEND-ONLY, and this is the append", got, want)
		}
	})

	t.Run("the emission carries it", func(t *testing.T) {
		f := requireOneCode(t,
			graphlint.Run(graphlint.NewRequest(mustLoad(t, dtMatchOnly0010))),
			graphlint.CodeUnprovableCoverage)

		if f.Reason == "" {
			t.Fatal("the emission carries an empty reason; that fails REQ-80 " +
				"outright")
		}
		if f.Reason != "no-participating-dimension" {
			t.Errorf("reason = %q; want %q — a dimension-scoped member states "+
				"a remedy that does not apply to a group-level emission",
				f.Reason, "no-participating-dimension")
		}
	})
}

// REQ-59 / REQ-60: "Where both would apply, this finding is reported and
// `graph-coverage-closed-by-escape` MUST NOT be" / "Precedence is delivered
// by **returning** from the class-keyed arm, not by ordering alone:
// emitting and then falling through to `emitCoverageArms` yields *both*
// findings over a zero-dimension group carrying a bare escape row, which is
// the double-report this clause forbids"
// ADVERSARIAL — the A13 double-report probe.
func TestReq60_TheZeroDimensionArmReturnsRatherThanFallingThrough(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtMatchOnlyBareEscape0010)))

	requireOneCode(t, r, graphlint.CodeUnprovableCoverage)
	requireNoCode(t, r, graphlint.CodeCoverageClosedByEscape)

	if hasCode(r, graphlint.CodeCoverageClosedByEscape) {
		t.Error("both findings are reported over a zero-dimension group " +
			"carrying a bare escape row; the class-keyed arm must RETURN, " +
			"not fall through to emitCoverageArms")
	}
}

// REQ-58: "the arm MUST be emitted inside the `len(dims) == 0` branch and
// MUST precede `emitCoverageArms`"
// BOUNDARY
//
// The observable consequence of the position: a zero-dimension group takes
// the unprovable finding and NOTHING from invariant 4's other arms — no
// gap, no withholding, no closure advisory.
func TestReq58_AZeroDimensionGroupTakesOnlyTheUnprovableFinding(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtMatchOnlyBareEscape0010)))

	for _, code := range []string{
		graphlint.CodeCoverageGap,
		graphlint.CodeCoverageClosedByEscape,
		graphlint.CodeProductTooLarge,
	} {
		requireNoCode(t, r, code)
	}
	requireOneCode(t, r, graphlint.CodeUnprovableCoverage)
}

// REQ-61 / ASSUMPTION-6: "The arm MUST supply its **own message** rather
// than routing through `internal/graphlint/coverage.go::unprovableMessage`,
// whose `default` limb would phrase the remedy as \"declare the domain\""
// DOMAIN EDGE
//
// The assertion is on the REMEDY the message must not state, not on exact
// wording (the wording is the implementer's, per ASSUMPTION-6).
func TestReq61_TheZeroDimensionArmSuppliesItsOwnMessage(t *testing.T) {
	f := requireOneCode(t,
		graphlint.Run(graphlint.NewRequest(mustLoad(t, dtMatchOnly0010))),
		graphlint.CodeUnprovableCoverage)

	if strings.Contains(f.Message, "declare the domain") {
		t.Errorf("message = %q; it routes through unprovableMessage's default "+
			"limb, whose remedy is for a dimension that does not exist",
			f.Message)
	}
	if f.Message == "" {
		t.Error("the arm supplies no message")
	}
	if !strings.Contains(f.Message, "dt/decide") {
		t.Errorf("message = %q; the arm names the GROUP", f.Message)
	}
}

// REQ-62: "which is why the clause binds the decision-table class only"
// ADVERSARIAL — negative REQ: an unconditional arm would redden 151
// zero-dimension groups across 37 of 37 loadable fixtures (RISK).
func TestReq62_TheZeroDimensionArmBindsTheDecisionTableClassOnly(t *testing.T) {
	// A STATE MACHINE whose one group ranges over no guard dimension. Over
	// this class the shape is legitimate and MUST stay silent.
	const smZeroDim = `outcomes = ["decide"]

[model]
id = "sm"
version = 1

[initial]
status = "draft"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "done"]
single_valued = true
required = true

[read.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"
read_back = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "decide"
[rule.match.status]
eq = "draft"
[rule.write]
status = "done"
`

	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, smZeroDim)))
	requireNoCode(t, r, graphlint.CodeUnprovableCoverage)
}

// REQ-63 / REQ-79 / ASSUMPTION-7: "The missing-root arm of invariant 1
// (`0006:C18`), dead end (2), always-present-owned (5),
// owned-set-before-match (6), and single-valued state are vacuous by
// construction over this class and MUST NOT emit." / "Add the fourth rule;
// lint → exit 0, findings exactly `[]` — no code from the 0006 taxonomy
// present (A10)."
// DOMAIN EDGE — asserted as EXACT finding-list equality against the full
// taxonomy (SC-4, ASSUMPTION-7).
func TestReq63_TheMachineOnlyInvariantsStaySilentOverADecisionTable(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtComplete0010)))

	if len(r.Findings) != 0 {
		t.Fatalf("the complete decision table reports %d findings; want "+
			"EXACTLY none — no code from the 0006 taxonomy is present:%s",
			len(r.Findings), render(r))
	}
}

// REQ-64: "Lint MUST NOT report the class itself as a finding: no finding's
// `Code`, `Element`, `Message`, or `Detail` may name the class or the
// string `decision-table`"
// ADVERSARIAL — negative REQ, in assertable form (`0010:BR6` rejected).
func TestReq64_NoFindingNamesTheClass(t *testing.T) {
	for name, src := range map[string]string{
		"partial":     dtPartial0010,
		"complete":    dtComplete0010,
		"escape":      dtEscape0010,
		"match-only":  dtMatchOnly0010,
		"two-outcome": dtTwoOutcomeEscape0010,
		"terminal":    dtStrayTerminal0010,
	} {
		t.Run(name, func(t *testing.T) {
			r := graphlint.Run(graphlint.NewRequest(mustLoad(t, src)))
			for _, f := range r.Findings {
				for field, v := range map[string]string{
					"Code": f.Code, "Element": f.Element,
					"Message": f.Message, "Hint": f.Hint,
				} {
					if strings.Contains(v, "decision-table") ||
						strings.Contains(v, "state-machine") {
						t.Errorf("finding %s carries %q = %q, naming the class",
							f.Code, field, v)
					}
				}
			}
		})
	}
}

// REQ-65 / REQ-97 / SC-4: "Invariant 1's **terminal-predicate arm** MUST
// stay live" / "the stray-`terminal` control reports `graph-dangling-edge`
// at `element = terminal[0]` — the arm C5 keeps live, distinguishing it
// from the silenced root arm by `element` rather than by code"
// ADVERSARIAL — this is what enforces C2's `terminal` prohibition.
func TestReq65_TheTerminalPredicateArmStaysLiveOverADecisionTable(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtStrayTerminal0010)))

	found := withCode(r, graphlint.CodeDanglingEdge)
	if len(found) == 0 {
		t.Fatalf("a decision table declaring `terminal` over an observed tag "+
			"reports no %s; silencing the code wholesale would remove C2's "+
			"enforcement of the terminal prohibition:%s",
			graphlint.CodeDanglingEdge, render(r))
	}

	var atTerminal bool
	for _, f := range found {
		if f.Element == "terminal[0]" {
			atTerminal = true
		}
		if f.Element == "model" {
			t.Errorf("the MISSING-ROOT arm fired over a decision table at "+
				"element %q; that arm is silenced by class while the "+
				"terminal arm stays live — they are distinguished by "+
				"ELEMENT, not by code", f.Element)
		}
	}
	if !atTerminal {
		t.Errorf("no %s at element `terminal[0]`; findings:%s",
			graphlint.CodeDanglingEdge, render(r))
	}
}

// REQ-66: "Declared-terminal handling (7) proper, and the escape arm of it,
// stay silent by construction: they key on a declared terminal, which C2
// forbids."
// DOMAIN EDGE
func TestReq66_DeclaredTerminalHandlingStaysSilentOverADecisionTable(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtEscape0010)))
	requireNoCode(t, r, graphlint.CodeTerminalEscape)
}

// REQ-67: "Over a `\"state-machine\"` model every 0006 invariant,
// `0006:C18` included, is unchanged."
// DOMAIN EDGE — negative REQ, the regression guarantee SC-6 sweeps.
//
// GREEN BY CONSTRUCTION against main: this is a characterization test
// pinning today's behaviour for the state-machine class.
func TestReq67_EveryInvariantIsUnchangedOverAStateMachine(t *testing.T) {
	m := mustLoad(t, smRootless0010)
	r := graphlint.Run(graphlint.NewRequest(m))

	f := requireOneCode(t, r, graphlint.CodeDanglingEdge)
	if f.Element != "model" {
		t.Errorf("element = %q; want %q", f.Element, "model")
	}
	if len(graphlint.Reach(m)) != 0 {
		t.Error("a rootless state machine traverses a non-empty node set")
	}
}

// REQ-68: "Of the four shipped `len(Initial) == 0` sites, exactly two
// become class-aware: `internal/graphlint/reach.go::reach` and
// `checkDanglingEdge`'s missing-root arm … The other two,
// `checkAlwaysPresentOwned` and `checkUnreachableRules`
// (`internal/graphlint/analysis.go`), keep the bare `len(Initial)` test"
// ADVERSARIAL — negative REQ: converting either of the latter two to the
// class test would flip it live over every decision table.
func TestReq68_TheOtherTwoInitialSitesKeepTheBareTest(t *testing.T) {
	for name, src := range map[string]string{
		"partial":    dtPartial0010,
		"complete":   dtComplete0010,
		"escape":     dtEscape0010,
		"match-only": dtMatchOnly0010,
	} {
		t.Run(name, func(t *testing.T) {
			r := graphlint.Run(graphlint.NewRequest(mustLoad(t, src)))
			// checkAlwaysPresentOwned's finding.
			requireNoCode(t, r, graphlint.CodeAlwaysPresentOwned)
			// checkUnreachableRules' advisory.
			requireNoCode(t, r, graphlint.CodeUnreachableRule)
		})
	}
}

// REQ-69: "Any site asking \"is this a decision table\" reads the `Class`
// field (empty = `state-machine`, C1) and never re-derives it from
// `len(owned) == 0`."
// ADVERSARIAL — negative REQ.
//
// The discriminating pair: the class-omitted control declares an owned tag,
// so a reader deriving the class from `len(owned) == 0` would agree with
// the correct implementation there. The probe below strips the OWNED TAG
// from a state-machine document — same empty owned set as a decision
// table, different declared class — and the missing-root finding must
// still fire.
func TestReq69_ClassReadersNeverRederiveFromTheOwnedSet(t *testing.T) {
	// A state-machine (class omitted) with ZERO owned tags and no root: an
	// implementation reading `len(owned) == 0` as "decision table" seeds it
	// at ∅ and reports nothing.
	const smZeroOwnedNoRoot = `outcomes = ["decide"]

[model]
id = "sm0"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "only"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
`

	m := mustLoad(t, smZeroOwnedNoRoot)
	r := graphlint.Run(graphlint.NewRequest(m))

	f := requireOneCode(t, r, graphlint.CodeDanglingEdge)
	if f.Element != "model" {
		t.Errorf("element = %q; want %q — a zero-owned STATE MACHINE is "+
			"rootless, and the class is declared, never re-derived from "+
			"`len(owned) == 0`", f.Element, "model")
	}
	if n := len(graphlint.Reach(m)); n != 0 {
		t.Errorf("Reach = %d nodes; a zero-owned state machine traverses "+
			"nothing", n)
	}
}

// REQ-37 / SC-4 / A7: "`emit` does not join
// `internal/graphlint/engine.go::Fingerprint`, so editing an emit block
// never changes a finding's identity"
// ADVERSARIAL — negative REQ, asserted by a fingerprint test.
func TestReq37_EditingAnEmitBlockDoesNotChangeAFingerprint(t *testing.T) {
	before := mustLoad(t, dtPartial0010)
	edited := mustLoad(t, strings.ReplaceAll(dtPartial0010,
		`verdict = "cell-xp"`, `verdict = "TOTALLY-DIFFERENT"`))

	if len(before.Rows) != len(edited.Rows) {
		t.Fatalf("the emit edit changed the row count: %d -> %d",
			len(before.Rows), len(edited.Rows))
	}

	for i := range before.Rows {
		a := graphlint.Fingerprint(before.Rows[i])
		b := graphlint.Fingerprint(edited.Rows[i])
		if a != b {
			t.Errorf("%s: fingerprint changed after editing only an emit "+
				"block: %q -> %q — `emit` does not join the fingerprint, so "+
				"a finding's identity is unaffected",
				before.Rows[i].Identity(), a, b)
		}
	}

	// The emit edit must be REAL, or the assertion is vacuous.
	var differed bool
	for i := range before.Rows {
		if !slices.Equal(before.Rows[i].Emit, edited.Rows[i].Emit) {
			differed = true
		}
	}
	if !differed {
		t.Fatal("no row's Emit differed; the fingerprint assertion would be " +
			"vacuous")
	}
}

// REQ-96 / SC-4: "the two-outcome escape control reports the unrescued
// outcome's coverage gap, pinning that \"the otherwise row\" is per-outcome
// (A9)"
// DOMAIN EDGE
func TestReq96_TheOtherwiseRowIsScopedPerOutcome(t *testing.T) {
	r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtTwoOutcomeEscape0010)))

	gaps := requireCode(t, r, graphlint.CodeCoverageGap)
	var reviewGap bool
	for _, f := range gaps {
		if strings.Contains(f.Element, "review") {
			reviewGap = true
		}
		if strings.Contains(f.Element, "decide") {
			t.Errorf("the RESCUED outcome `decide` reports a coverage gap at "+
				"element %q; the escape row closes its own outcome's arm",
				f.Element)
		}
	}
	if !reviewGap {
		t.Errorf("the UNRESCUED outcome `review` reports no coverage gap; "+
			"`the otherwise row` is per-outcome (`0002:C5`):%s", render(r))
	}
}

// REQ-94 / SC-4 (`0010:S4`): the graph-lint fixture PAIR — partial → one
// coverage finding naming the cell; complete, closed by a fourth ORDINARY
// rule → `[]`; plus the escape-"otherwise" variant.
// HAPPY PATH
func TestReq94_TheGraphLintFixturePairAndTheEscapeVariant(t *testing.T) {
	t.Run("partial reports one coverage finding naming the cell", func(t *testing.T) {
		r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtPartial0010)))
		f := requireOneCode(t, r, graphlint.CodeCoverageGap)
		if f.Element != "dt/decide" {
			t.Errorf("element = %q; want %q", f.Element, "dt/decide")
		}
		// "naming the cell": the uncovered cell is identified by the group's
		// dimensions, and the count of uncovered assignments is what makes
		// this a POSITIVE finding rather than a lint that never reached the
		// group machinery. `1 of 4` is the shipped renderer's own phrasing
		// of the 3-of-4-cells-covered table; asserting the count keeps the
		// oracle on the computed quantity rather than on prose.
		if !strings.Contains(f.Message, "1 of 4") {
			t.Errorf("message = %q; the partial table covers three of the "+
				"four (a,b) cells, so exactly one assignment is uncovered",
				f.Message)
		}
	})

	t.Run("complete reports nothing", func(t *testing.T) {
		r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtComplete0010)))
		if len(r.Findings) != 0 {
			t.Errorf("findings = %d; want exactly []:%s", len(r.Findings), render(r))
		}
	})

	t.Run("the escape variant reports only the closure advisory", func(t *testing.T) {
		r := graphlint.Run(graphlint.NewRequest(mustLoad(t, dtEscape0010)))
		if got := codesIn(r); !slices.Equal(got,
			[]string{graphlint.CodeCoverageClosedByEscape}) {
			t.Errorf("codes = %v; want exactly [%s] — never a bare green",
				got, graphlint.CodeCoverageClosedByEscape)
		}
		if len(r.Blocking()) != 0 {
			t.Errorf("the escape variant carries %d blocking findings; the "+
				"advisory is info-severity and the run exits 0",
				len(r.Blocking()))
		}
	})
}

// REQ-72: "no concurrency surface is added. The class is read from the
// loaded model, immutable after load; lint and resolve are single-pass over
// that value."
// DOMAIN EDGE — negative REQ.
//
// Assertable form: lint is idempotent and does not mutate the model it is
// handed.
func TestReq72_LintIsSinglePassAndDoesNotMutateTheModel(t *testing.T) {
	m := mustLoad(t, dtPartial0010)
	class := m.Class

	first := codesIn(graphlint.Run(graphlint.NewRequest(m)))
	second := codesIn(graphlint.Run(graphlint.NewRequest(m)))

	if !slices.Equal(first, second) {
		t.Errorf("two lint runs over one model disagree: %v vs %v", first, second)
	}
	if m.Class != class {
		t.Errorf("lint mutated Model.Class: %q -> %q", class, m.Class)
	}
}

// REQ-103 / IP Phase 2: "the graph-lint fixtures gain a decision-table pair
// (partial → coverage finding; complete → clean) plus the class-omitted
// control."
// HAPPY PATH — the phase's deliverable, asserted as one table.
func TestReq103_TheDecisionTableFixtureSetIsPresentAndDiscriminating(t *testing.T) {
	cases := []struct {
		name  string
		src   string
		codes []string
	}{
		{"partial", dtPartial0010, []string{graphlint.CodeCoverageGap}},
		{"complete", dtComplete0010, nil},
		{"escape", dtEscape0010, []string{graphlint.CodeCoverageClosedByEscape}},
		{"class-omitted", smRootless0010, []string{graphlint.CodeDanglingEdge}},
		{"match-only", dtMatchOnly0010, []string{graphlint.CodeUnprovableCoverage}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := graphlint.Run(graphlint.NewRequest(mustLoad(t, c.src)))
			got := codesIn(r)
			want := c.codes
			if want == nil {
				want = []string{}
			}
			if len(got) == 0 {
				got = []string{}
			}
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Errorf("codes = %v; want exactly %v:%s", got, want, render(r))
			}
		})
	}
}

// --- shared table-package assertion --------------------------------------

// REQ-52 (companion): the decision-table model's own class survives the
// loader, so lint's class-keyed sites have something to read.
// HAPPY PATH
func TestReq52_TheLoadedModelCarriesTheClassForLintToRead(t *testing.T) {
	m := mustLoad(t, dtComplete0010)
	if !table.IsDecisionTable(m) {
		t.Errorf("Class = %q; lint's class-keyed sites read the declared "+
			"class off the loaded model", m.Class)
	}
}
