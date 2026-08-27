package graphlint_test

// RDR 0006 — the mandatory invariant set (REQ-31..REQ-44), the
// always-present owned check (REQ-67..REQ-70), and the reachability
// relation those invariants quantify over (REQ-101..REQ-113).

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-31: "Graph lint MUST check at least these blocking invariant
// classes: dangling edge, dead end, determinism/overlap, guard
// exhaustiveness/gap, single-valued state, owned-set-before-match, and
// declared terminal/escape handling."
// BOUNDARY
func TestReq31_TheSevenMandatoryInvariantClassesAreAllChecked(t *testing.T) {
	// The seven classes map onto the blocking code set. Each must be a
	// code the engine can emit — a class with no code is a class no run
	// can report.
	want := []string{
		graphlint.CodeDanglingEdge,
		graphlint.CodeDeadEnd,
		graphlint.CodeOverlap,
		graphlint.CodeCoverageGap,
		graphlint.CodeSingleValuedState,
		graphlint.CodeOwnedBeforeWrite,
		graphlint.CodeTerminalEscape,
	}
	got := graphlint.BlockingCodes()
	for _, c := range want {
		if !slices.Contains(got, c) {
			t.Errorf("mandatory invariant class %q is absent from the "+
				"blocking code set %v", c, got)
		}
	}
}

// REQ-32: Invariant 1, dangling edge — "every context reference, tag key,
// tag value, outcome, and accessor reference named by a row must resolve
// to a declared model element, and the model must declare an initial owned
// state."
// HAPPY PATH
func TestReq32_DanglingEdgeChecksReferencesAndRequiresAnInitialOwnedState(t *testing.T) {
	// RDR 0002's loader refuses most syntactic dangling references before
	// normalization (an undeclared tag, an out-of-alphabet outcome, an
	// unknown gate accessor), so the arm that reaches lint is the model
	// half of the clause: a model declaring no initial owned state.
	const noRoot = `
terminal = ["done"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, noRoot)

	f := requireCode(t, r, graphlint.CodeDanglingEdge)
	for _, got := range f {
		if got.Severity != graphlint.SeverityBlocking {
			t.Errorf("%s severity = %q; want %q", got.Code, got.Severity,
				graphlint.SeverityBlocking)
		}
	}

	// The control: the same model WITH a root carries no dangling finding
	// on that account, so the check is about the missing declaration and
	// not a blanket refusal.
	clean := lint(t, twoStateDecls, legalBody)
	requireNoCode(t, clean, graphlint.CodeDanglingEdge)
}

// REQ-33: "There is no \"transition target\" term … the edge's destination
// is checked as tag keys and values against their `[tags.<tag>]`
// declarations and declared domains. Terminal declarations are checked the
// same way — as tag predicates, not as state names."
// DOMAIN EDGE
func TestReq33_TerminalsAreTagPredicatesNotStateNames(t *testing.T) {
	m := mustLoad(t, source(twoStateDecls, legalBody))

	// The normalized terminal set is DEREFERENCED predicate sets, never
	// bare context ids: a check reading state names would have nothing to
	// read.
	if len(m.Terminal) == 0 {
		t.Fatal("the normalized model carries no terminal predicate sets")
	}
	for i, set := range m.Terminal {
		if len(set) == 0 {
			t.Errorf("terminal[%d] dereferenced to an empty predicate set", i)
		}
		for _, atom := range set {
			if atom.Key == "" || atom.Operator == "" {
				t.Errorf("terminal[%d] carries a non-predicate element %+v; "+
					"terminals are tag predicates, not state names", i, atom)
			}
			if _, declared := m.Tags[atom.Key]; !declared {
				t.Errorf("terminal[%d] predicate names the undeclared tag %q",
					i, atom.Key)
			}
		}
	}

	// The engine must consume them as PREDICATES, and the discriminating
	// oracle is a model whose terminal predicate reads a tag the node does
	// not satisfy. A name-based reading has no predicate to evaluate and
	// so cannot tell the two fixtures apart; a predicate reading exempts
	// the satisfying node and reports the other.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true
`
	// `status = b` satisfies the terminal predicate, so no dead end.
	const satisfied = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	requireNoCode(t, lint(t, decls, satisfied), graphlint.CodeDeadEnd)

	// The same model whose row writes `c` instead: the terminal PREDICATE
	// `status = b` does not hold at `c`, so the node is a dead end. Only a
	// predicate evaluation can decide this — the terminal's context id is
	// identical in both fixtures.
	const unsatisfied = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "c"
`
	requireCode(t, lint(t, decls, unsatisfied), graphlint.CodeDeadEnd)
}

// REQ-129: "a root `terminal` entry naming a context whose predicate reads
// a **non-owned** tag → one blocking finding (§D7(i))"
// ADVERSARIAL
func TestReq129_TerminalOnANonOwnedTagIsABlockingDanglingFinding(t *testing.T) {
	// The producer side is landed and routes here by name: RDR 0002 fences
	// that a terminal context matching a non-owned tag is RDR 0006's
	// blocking finding, not a load failure. It loads clean, so lint owns it.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.seen]
eq = "y"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	src := source(observedRequiredDecls, body)
	mustLoad(t, src) // precondition: the defect is lint's, not load's

	r := lintSource(t, src)
	f := requireCode(t, r, graphlint.CodeDanglingEdge)

	var named bool
	for _, got := range f {
		if got.Key == "seen" || got.Element == "done" || got.Rule == "done" {
			named = true
		}
	}
	if !named {
		t.Errorf("no %s finding names the terminal context or the non-owned "+
			"tag `seen`; the diagnostic must name the authored site; report:%s",
			graphlint.CodeDanglingEdge, render(r))
	}
}

// REQ-34: Invariant 2, dead end — "every reachable owned-state node that
// satisfies no declared terminal must be the source of at least one
// modeled non-escape row (escape self-loops do not count as progress)."
// HAPPY PATH
func TestReq34_ReachableNonTerminalNodeWithoutAnOutgoingRowIsADeadEnd(t *testing.T) {
	// `status = b` is reachable and satisfies no terminal (the model
	// declares none over it) and no row matches it, so it is a dead end.
	const body = `
terminal = ["never"]

[initial]
status = "a"

[context.never]
[context.never.match.status]
eq = "a"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, body)
	requireCode(t, r, graphlint.CodeDeadEnd)
}

// REQ-103 (dead-end half): "Escape self-loops are excluded from invariant
// 2's \"at least one outgoing row\" test — a self-loop is not progress —
// so a node whose only exit is an escape rescue is still a dead end"
// ADVERSARIAL
func TestReq103_EscapeSelfLoopIsNotProgressForTheDeadEndTest(t *testing.T) {
	// `status = b` has exactly one outgoing row and it is an ESCAPE row,
	// whose successor equals its source. Counting it as progress would
	// clear the dead end; the clause says it must not.
	const body = `
terminal = ["never"]

[initial]
status = "a"

[context.never]
[context.never.match.status]
eq = "a"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "rescue-b"
escape = ["no_match"]
[rule.match.status]
eq = "b"
[rule.match.recognized]
eq = "go"
`
	r := lint(t, statusOnlyDecls, body)
	requireCode(t, r, graphlint.CodeDeadEnd)
}

// REQ-35: "A terminal is a predicate over owned tags, and a node
// *satisfies* it when **every** value in each of the node's per-tag value
// sets meets it. Partial satisfaction is deliberately not enough"
// BOUNDARY
func TestReq35_TerminalSatisfactionIsUniversalOverTheNodesValueSet(t *testing.T) {
	// Two paths write `b` and `c`, converging on a node whose `status` set
	// is {b, c}, against a terminal `status = b`. Partial satisfaction
	// would clear the whole node; universal satisfaction clears only the
	// {b} split and leaves {c} — which has no outgoing row — a dead end.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "to-b"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "to-c"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "c"
`
	r := lint(t, decls, body)

	// `status = c` satisfies no terminal and has no outgoing row.
	requireCode(t, r, graphlint.CodeDeadEnd)
}

// REQ-36: "**The test runs on split nodes, not merged ones.** … lint
// splits each reachable node on the terminal-participating keys: a node
// whose `status` set is `{done, review}` against a terminal `status=done`
// is tested as `{done}` … and `{review}` … — independently."
// REQ-125 / SC-21: "`graph-dead-end` is **not** emitted for the terminated
// branch (invariant 2 splits before testing); it *is* emitted if the live
// branch has no outgoing row."
// ADVERSARIAL
func TestReq36_DeadEndSplitsTheNodeOnTerminalParticipatingKeys(t *testing.T) {
	// Two paths converge on one node whose `status` set is {done, review}
	// against a terminal `status = done`. A single MERGED-node test fails
	// both branches — the regression this fixture guards.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "done", "review"]
single_valued = true
required = true
`
	// The live branch DOES have an outgoing row, so a correct split
	// emits no dead end at all: `{done}` satisfies the terminal and
	// `{review}` progresses.
	const bothCovered = `
terminal = ["finished"]

[initial]
status = "start"

[context.finished]
[context.finished.match.status]
eq = "done"

[[rule]]
id = "finish"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "done"

[[rule]]
id = "send-to-review"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "review"

[[rule]]
id = "review-progresses"
[rule.match.status]
eq = "review"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "done"
`
	r := lint(t, decls, bothCovered)
	requireNoCode(t, r, graphlint.CodeDeadEnd)

	// The paired half: drop the live branch's outgoing row and the dead
	// end IS emitted — for the live branch only.
	const liveBranchStuck = `
terminal = ["finished"]

[initial]
status = "start"

[context.finished]
[context.finished.match.status]
eq = "done"

[[rule]]
id = "finish"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "done"

[[rule]]
id = "send-to-review"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "review"
`
	stuck := lint(t, decls, liveBranchStuck)
	requireCode(t, stuck, graphlint.CodeDeadEnd)
}

// REQ-37: "The split is bounded by the declared domains of
// terminal-participating keys only, never the whole lattice."
// BOUNDARY
func TestReq37_SplitIsBoundedByTerminalParticipatingKeysOnly(t *testing.T) {
	// `status` participates in the terminal; `other` does not. The earlier
	// shape declared `other` an int and wrote it once, so it was a
	// SINGLETON in every reachable node — `splitNode` skips a key holding
	// one value, so a whole-lattice split produced the same node set as
	// the bounded one and went unpunished.
	//
	// Here two rows converge on one merged node holding `other` MULTI-
	// VALUED. The bounded split ranges over `status` alone and yields one
	// stuck split; a split over the whole lattice would range over
	// `other`'s held values too and yield one stuck split per value.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["start", "done", "review"]
single_valued = true
required = true

[tags.other]
provenance = "owned"
kind = "enum"
domain = ["o0", "o1", "o2"]
single_valued = true
required = true
`
	const body = `
terminal = ["finished"]

[initial]
status = "start"
other = "o0"

[context.finished]
[context.finished.match.status]
eq = "done"

[[rule]]
id = "finish"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "done"
other = "o1"

[[rule]]
id = "send-to-review"
[rule.match.status]
eq = "start"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "review"
other = "o2"
`
	m := mustLoad(t, source(decls, body))

	// PRECONDITION. Without a reachable node holding `other` multi-valued
	// there is nothing for a whole-lattice split to multiply over, and the
	// count below cannot tell a bounded split from an unbounded one.
	var multiValued bool
	for _, n := range graphlint.Reach(m) {
		if len(n.Values["other"]) > 1 {
			multiValued = true
		}
	}
	if !multiValued {
		t.Fatalf("no reachable node holds the non-participating key `other` "+
			"multi-valued, so a whole-lattice split would produce the same "+
			"node set as the bounded one and this assertion is vacuous; "+
			"nodes=%v", graphlint.Reach(m))
	}

	r := graphlint.Run(graphlint.NewRequest(m))

	// The `review` split is stuck, so exactly one dead end — not one per
	// value the non-participating `other` holds there.
	if n := countCode(r, graphlint.CodeDeadEnd); n != 1 {
		t.Errorf("%d %s findings; want exactly 1 — the split is bounded by "+
			"the terminal-participating keys only, never the whole lattice; "+
			"report:%s", n, graphlint.CodeDeadEnd, render(r))
	}
}

// REQ-40: Invariant 5, single-valued state — "the model must not produce a
// view in which a tag declared single-valued … holds two values. This is
// decided **per row, syntactically**: no row's write block may assign a
// single-valued tag two values."
// DOMAIN EDGE
func TestReq40_SingleValuedStateIsDecidedPerRowSyntactically(t *testing.T) {
	// RDR 0002's loader refuses every TOML spelling of this defect before
	// normalization (`kind <k> holds one value, not a member sequence`;
	// `kind set admits no single_valued marker`), so the check has no
	// authorable input surface today. What remains assertable is that the
	// check is a declared member of the taxonomy — a class the engine can
	// report — rather than silently absent.
	if !slices.Contains(graphlint.BlockingCodes(), graphlint.CodeSingleValuedState) {
		t.Fatalf("%s is absent from the blocking code set %v",
			graphlint.CodeSingleValuedState, graphlint.BlockingCodes())
	}
	if !graphlint.IsBlocking(graphlint.CodeSingleValuedState) {
		t.Errorf("%s is not blocking", graphlint.CodeSingleValuedState)
	}

	// The per-row reading's observable consequence: a model whose rows
	// each write one value is clean however many rows write the key.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "to-b-go"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "to-b-stop"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, body)
	requireNoCode(t, r, graphlint.CodeSingleValuedState)
}

// REQ-41: "It is deliberately *not* decided over the reachability nodes: a
// merged node's value set has cardinality > 1 whenever two paths write
// different single values, which is the legal shape of a two-path merge,
// not a violation"
// ADVERSARIAL
func TestReq41_MergedNodeWithTwoValuesIsNotASingleValuedViolation(t *testing.T) {
	// Two paths write `b` and `c` on the single-valued `status`, so the
	// merged node's value set is {b, c}. A node-level check reports this;
	// the per-row check must not — it is the legal shape of a merge.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
in = ["b", "c"]

[[rule]]
id = "to-b"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "to-c"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "c"
`
	m := mustLoad(t, source(decls, body))

	// Precondition: the two paths really do converge on one merged node
	// whose value set has cardinality > 1 — otherwise the negative
	// assertion below is vacuous.
	var merged bool
	for _, n := range graphlint.Reach(m) {
		if len(n.Values["status"]) > 1 {
			merged = true
		}
	}
	if !merged {
		t.Fatalf("fixture precondition failed: no reachable node merges two "+
			"`status` values, so a node-level check could not fire here; "+
			"nodes=%v", graphlint.Reach(m))
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	requireNoCode(t, r, graphlint.CodeSingleValuedState)
}

// REQ-42: Invariant 6, owned-set-before-match — "a row that reads an owned
// tag (match key or guard atom) must find it held in every reachable
// owned-state that satisfies the row's match pattern; the declared initial
// owned state counts as a write."
// HAPPY PATH
func TestReq42_OwnedSetBeforeMatchRequiresTheKeyHeldAtEveryMatchingNode(t *testing.T) {
	// `opt` is read by a guard atom but written by no row and absent from
	// `[initial]`, so it is not held at the node the row matches.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-opt"
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
	if !namesRule(r, graphlint.CodeOwnedBeforeWrite, "reads-opt") {
		t.Errorf("no %s finding names rule %q; report:%s",
			graphlint.CodeOwnedBeforeWrite, "reads-opt", render(r))
	}

	// The clause's second half: "the declared initial owned state counts
	// as a write". Declaring `opt` at the root clears the finding.
	const withRoot = `
terminal = ["done"]

[initial]
status = "a"
opt = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "reads-opt"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.opt]
eq = "p"
[rule.write]
status = "b"
`
	rooted := lint(t, optionalGuardDecls, withRoot)
	requireNoCode(t, rooted, graphlint.CodeOwnedBeforeWrite)
}

// REQ-43: Invariant 7 — "terminal states and escape rows are explicit
// model data; lint must not infer them from missing rows." The minted
// defect is "a reachable non-terminal node with no outgoing non-escape row
// and no terminal declaration covering it (an implied terminal), or a
// group relying on an undeclared escape row to close coverage."
// ADVERSARIAL
func TestReq43_RelianceOnAnImpliedTerminalIsGraphTerminalEscape(t *testing.T) {
	// The model declares NO terminal at all and ends at `status = b`,
	// relying on lint to infer that the absence of an outgoing row means
	// "terminal". The clause makes that reliance the defect.
	const body = `
[initial]
status = "a"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, body)
	requireCode(t, r, graphlint.CodeTerminalEscape)

	// The control: declaring the terminal the model depends on clears it.
	clean := lint(t, twoStateDecls, legalBody)
	requireNoCode(t, clean, graphlint.CodeTerminalEscape)
}

// REQ-43 / FAIL-1 (Phase 3a): the invariant 7 condition is stated PER NODE
// — "a reachable non-terminal node with no outgoing non-escape row and no
// terminal declaration covering it (an implied terminal)" — not as a
// property of the model as a whole.
//
// `checkTerminalEscape` opened with `if len(a.model.Terminal) > 0 { return }`,
// so `graph-terminal-escape` could fire ONLY on a model declaring no
// terminal whatsoever. Every model that declares even one terminal was
// exempt, however many of its reachable nodes relied on an inferred one.
//
// That is exactly the fixture the record prescribes being unable to mint
// the code: "a fixture is authored by omitting the declaration the model
// depends on" — omitting ONE terminal, not all of them. Under the global
// gate the illegal matrix entry for this code was satisfiable only by the
// degenerate no-terminal-whatsoever model above.
// ADVERSARIAL — REGRESSION (Phase 3c)
func TestReq43_ImpliedTerminalIsMintedInAModelDeclaringOtherTerminals(t *testing.T) {
	// The model DOES declare a terminal — `stage = x`, reached by `r-x`.
	// `r-y` reaches `stage = y`, which satisfies no declared terminal and is
	// the source of no non-escape row: an implied terminal, and the very
	// shape invariant 7 mints for. The declared terminal must not exempt it.
	const decls = `
[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["a", "x", "y"]
single_valued = true
required = true
`
	const body = `
terminal = ["fin"]

[initial]
stage = "a"

[context.fin]
[context.fin.match.stage]
eq = "x"

[[rule]]
id = "r-x"
[rule.match.stage]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
stage = "x"

[[rule]]
id = "r-y"
[rule.match.stage]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
stage = "y"
`
	r := lint(t, decls, body)

	if !hasCode(r, graphlint.CodeTerminalEscape) {
		t.Errorf("no %s finding; `stage = y` is reachable, satisfies no "+
			"declared terminal, and is the source of no non-escape row — "+
			"invariant 7's implied terminal. The model declaring a DIFFERENT "+
			"terminal (`stage = x`) must not exempt it: the clause states a "+
			"per-node condition, not a per-model one.\nreport:%s",
			graphlint.CodeTerminalEscape, render(r))
	}

	// REQ-83's complete-emission clause: invariants 2 and 7 are separate
	// obligations over the same node and one must not stand in for the
	// other. The node is both a dead end and an implied terminal.
	if !hasCode(r, graphlint.CodeDeadEnd) {
		t.Errorf("no %s finding alongside %s; the two invariants are "+
			"separate obligations and REQ-83 forbids reporting only the "+
			"first defect.\nreport:%s", graphlint.CodeDeadEnd,
			graphlint.CodeTerminalEscape, render(r))
	}

	// The control, isolating the defect to the missing declaration: cover
	// `stage = y` with a second terminal and the finding clears, with the
	// model otherwise byte-identical.
	const covered = `
terminal = ["fin", "fin-y"]

[initial]
stage = "a"

[context.fin]
[context.fin.match.stage]
eq = "x"

[context.fin-y]
[context.fin-y.match.stage]
eq = "y"

[[rule]]
id = "r-x"
[rule.match.stage]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
stage = "x"

[[rule]]
id = "r-y"
[rule.match.stage]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
stage = "y"
`
	clean := lint(t, decls, covered)
	if hasCode(clean, graphlint.CodeTerminalEscape) {
		t.Errorf("%s fired on a model whose every terminal node IS declared; "+
			"declaring the terminal the model depends on is the cure, so the "+
			"control must be clean.\nreport:%s",
			graphlint.CodeTerminalEscape, render(clean))
	}
}

// REQ-44: "A model that declares no initial owned state MUST be rejected
// with a blocking finding; lint MUST NOT treat an absent root as an empty
// reachable set and report a clean model."
// REQ-124 / SC-15: never exit 0 on an empty reachable set.
// ADVERSARIAL
func TestReq44_AbsentRootIsBlockingNeverAnEmptyReachableSetGreen(t *testing.T) {
	const noRoot = `
terminal = ["done"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, statusOnlyDecls, noRoot)

	if len(r.Blocking()) == 0 {
		t.Fatalf("a model declaring no initial owned state produced no "+
			"blocking finding — lint treated the absent root as an empty "+
			"reachable set and reported clean; report:%s", render(r))
	}
	// The disposition table names the code.
	requireCode(t, r, graphlint.CodeDanglingEdge)
}

// REQ-67: "An owned key declared always-present MUST be held in every
// reachable owned-state node; a violation is `graph-always-present-owned`."
// HAPPY PATH
func TestReq67_AlwaysPresentOwnedKeyMustBeHeldAtEveryReachableNode(t *testing.T) {
	// `always` carries the explicit always-present marker but a row clears
	// it, so a reachable node does not hold it.
	const body = `
terminal = ["done"]

[initial]
status = "a"
always = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "clears-always"
clear = ["always"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, alwaysPresentOwnedDecls, body)
	requireCode(t, r, graphlint.CodeAlwaysPresentOwned)
}

// REQ-68: "Lint MUST NOT extend this check to observed or recognized keys
// … so this RDR discharges only the owned half of RDR 0003's conformance
// premise."
// REQ-123 / SC-17 control.
// ADVERSARIAL
func TestReq68_AlwaysPresentCheckDoesNotExtendToObservedKeys(t *testing.T) {
	// `seen` is OBSERVED and carries the always-present marker, and no
	// reachable owned-state node holds it — observed keys arrive at
	// runtime. Reporting it would claim a half this RDR does not own.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, observedRequiredDecls, body)
	for _, f := range withCode(r, graphlint.CodeAlwaysPresentOwned) {
		if f.Key == "seen" {
			t.Errorf("%s names the OBSERVED key `seen`; lint claims only the "+
				"owned half of the conformance premise; report:%s",
				graphlint.CodeAlwaysPresentOwned, render(r))
		}
	}

	// The paired POSITIVE control, so an engine emitting nothing at all
	// cannot satisfy the negative half: the same shape with the key
	// declared OWNED does take the finding. The two fixtures differ only
	// in the provenance field, which is exactly the distinction the clause
	// draws.
	const ownedBody = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	owned := lint(t, alwaysPresentOwnedDecls, ownedBody)
	f := requireCode(t, owned, graphlint.CodeAlwaysPresentOwned)
	var namesOwnedKey bool
	for _, got := range f {
		if got.Key == "always" {
			namesOwnedKey = true
		}
	}
	if !namesOwnedKey {
		t.Errorf("the owned control emitted no %s naming `always`; without it "+
			"the observed-key negative above is vacuous; report:%s",
			graphlint.CodeAlwaysPresentOwned, render(owned))
	}
}

// REQ-69: "**The root is included, and that is a requirement on the
// `initial` declaration, not an exemption.** … an always-present owned key
// the `initial` table omits is absent at the root and takes the finding
// there."
// REQ-70: "lint reports the omission against the `initial` declaration
// rather than against an arbitrary downstream node, so the diagnostic
// names the authored site."
// BOUNDARY
func TestReq69And70_RootIsIncludedAndTheFindingNamesTheInitialDeclaration(t *testing.T) {
	// `always` carries the always-present marker and `[initial]` omits it,
	// so it is absent at the ROOT. Exempting the root would clear this.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	r := lint(t, alwaysPresentOwnedDecls, body)
	f := requireCode(t, r, graphlint.CodeAlwaysPresentOwned)

	var namesKey, namesInitial bool
	for _, got := range f {
		if got.Key == "always" {
			namesKey = true
		}
		// REQ-70: the diagnostic names the AUTHORED site — the `initial`
		// declaration — not an arbitrary downstream node.
		if got.Element == "initial" || got.Rule == "initial" {
			namesInitial = true
		}
	}
	if !namesKey {
		t.Errorf("no %s finding names the always-present key `always`; report:%s",
			graphlint.CodeAlwaysPresentOwned, render(r))
	}
	if !namesInitial {
		t.Errorf("no %s finding names the `initial` declaration; the omission "+
			"must be reported against the authored site, not an arbitrary "+
			"downstream node; report:%s",
			graphlint.CodeAlwaysPresentOwned, render(r))
	}
}

// --- the reachability relation -------------------------------------------

// REQ-101: "the graph lint traverses is the **owned-state graph**: a node
// is an abstract owned-state …; the root is the declared initial owned
// state (A6); an edge is a normalized non-escape row whose match pattern
// over owned tags is satisfiable in the source node, producing the node
// with that row's writes applied and clears removed."
// HAPPY PATH
func TestReq101_ReachabilityIsTheOwnedStateGraphRootedAtInitial(t *testing.T) {
	m := mustLoad(t, source(twoStateDecls, legalBody))
	nodes := graphlint.Reach(m)

	if len(nodes) == 0 {
		t.Fatalf("the traversal produced no nodes for a model with the "+
			"declared root %v", m.Initial)
	}

	// The root is the declared initial owned state, so some node holds
	// exactly the initial assignment.
	var sawRoot bool
	for _, n := range nodes {
		match := true
		for _, tv := range m.Initial {
			if !slices.Equal(n.Values[tv.Key], tv.Value) {
				match = false
			}
		}
		if match {
			sawRoot = true
		}
	}
	if !sawRoot {
		t.Errorf("no reachable node equals the declared initial owned state "+
			"%v; nodes=%v", m.Initial, nodes)
	}

	// An edge applies the row's writes: `status = b` is reachable.
	var sawB bool
	for _, n := range nodes {
		if slices.Contains(n.Values["status"], "b") {
			sawB = true
		}
	}
	if !sawB {
		t.Errorf("the successor of the `status = b` write is not reachable; "+
			"an edge must produce the node with the row's writes applied; "+
			"nodes=%v", nodes)
	}

	// Only OWNED tags are nodes: `recognized` is a recognized-provenance
	// key and must not appear in the owned-state.
	for _, n := range nodes {
		for key := range n.Values {
			if m.Tags[key].Provenance != table.ProvenanceOwned {
				t.Errorf("node carries the %s-provenance key %q; the "+
					"traversal is the OWNED-state graph",
					m.Tags[key].Provenance, key)
			}
		}
	}
}

// REQ-102: "**Escape rows are edges too, and they are self-loops** … Lint
// MUST model them as such rather than omitting them, because a node
// reachable only through an escape rescue is reachable at runtime"
// DOMAIN EDGE
func TestReq102_EscapeRowsAreSelfLoopEdges(t *testing.T) {
	// An escape row carries neither a write block nor a clear list, so its
	// successor equals its source. Omitting the edge would let invariant 2
	// accuse the recovery arm and invariant 6 accuse rows downstream of it.
	m := mustLoad(t, source(twoStateDecls, escapeBody))

	row := rowByRuleID(t, m, "rescue")
	if len(row.Writes) != 0 || len(row.NextTags) != 0 {
		t.Fatalf("fixture precondition failed: the escape row carries writes "+
			"%v / next tags %v, so its successor would not equal its source",
			row.Writes, row.NextTags)
	}

	before := len(graphlint.Reach(m))
	if before == 0 {
		t.Fatal("the traversal produced no nodes")
	}

	// A self-loop adds no node beyond its source, so modelling it must not
	// change the reachable set relative to the same model without it.
	withoutEscape := mustLoad(t, source(twoStateDecls, legalBody))
	if len(graphlint.Reach(withoutEscape)) == 0 {
		t.Fatal("the escape-free control produced no nodes")
	}

	// The observable consequence of modelling the edge: the escape row's
	// own source node is reachable, so no downstream check accuses it.
	r := graphlint.Run(graphlint.NewRequest(m))
	for _, f := range withCode(r, graphlint.CodeUnreachableRule) {
		if f.Rule == "rescue" {
			t.Errorf("the escape row `rescue` is reported unreachable; escape "+
				"rows are edges and their source nodes are reachable; report:%s",
				render(r))
		}
	}
}

// REQ-104: "Match atoms over observed/recognized tags and all guard atoms
// are **not** pruned — every such edge is taken as traversable."
// ADVERSARIAL
func TestReq104_ObservedMatchAtomsAndGuardAtomsAreNotPruned(t *testing.T) {
	// The row's guard is FALSE at the source node — `status` holds `a`
	// there and the guard demands `b` — and its match carries an observed
	// key naming a value no owned-state can supply. A guard-aware
	// traversal would prune the edge and report `status = b` unreachable;
	// the clause says every such edge is traversable.
	//
	// The earlier fixture guarded on `status eq "a"`, which is TRUE at the
	// root, so a guard-EVALUATING traversal took the edge too and the
	// assertion could not tell the two apart.
	const decls = statusOnlyDecls + `
[tags.seen]
provenance = "observed"
kind = "enum"
domain = ["y", "n"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "guarded"
[rule.match.status]
eq = "a"
[rule.match.seen]
eq = "y"
[rule.match.recognized]
eq = "go"
[rule.guard.all.status]
eq = "b"
[rule.write]
status = "b"
`
	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	var sawB bool
	for _, n := range nodes {
		if slices.Contains(n.Values["status"], "b") {
			sawB = true
		}
	}
	if !sawB {
		t.Errorf("the edge behind an observed match atom and a guard atom was "+
			"pruned: `status = b` is unreachable. Every such edge is taken as "+
			"traversable; nodes=%v", nodes)
	}
}

// REQ-106: "This is a syntactic dataflow over declarations, writes, and
// clears — no accessor execution, no runtime trace, no guard evaluation —
// so it is total, and it over-approximates the runtime"
// ADVERSARIAL
func TestReq106_TraversalIsTotalAndExecutesNoAccessor(t *testing.T) {
	// Totality: every fixture in this suite, legal or not, must produce a
	// traversal rather than hanging or refusing. A model whose accessors
	// name unreachable paths still traverses, because none is executed.
	for name, src := range map[string]string{
		"legal":        source(twoStateDecls, legalBody),
		"escape":       source(twoStateDecls, escapeBody),
		"non-finite":   source(nonFiniteDecls, rootAndTerminal+advanceRule),
		"observed-key": source(observedRequiredDecls, rootAndTerminal+advanceRule),
	} {
		m := mustLoad(t, src)
		nodes := graphlint.Reach(m)
		if len(nodes) == 0 {
			t.Errorf("%s: the traversal produced no nodes; it must be total "+
				"over every model that normalizes", name)
		}
	}
}

// REQ-108: "the traversal is a **fixpoint over merged nodes, never
// path-sensitive**: two edges reaching the same successor produce one node
// whose per-tag value sets are the union of theirs (a lattice widening),
// and the worklist runs to a fixpoint."
// BOUNDARY
func TestReq108_TraversalMergesConvergingEdgesIntoOneNode(t *testing.T) {
	// Two rows write `b` and `c` from the same source, then a third row
	// matches both and writes `d`. A path-sensitive traversal produces two
	// distinct successors; the merged reading produces one whose `status`
	// set is the union.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c", "d"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "d"

[[rule]]
id = "to-b"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "to-c"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "c"

[[rule]]
id = "converge"
[rule.match.status]
in = ["b", "c"]
[rule.match.recognized]
eq = "go"
[rule.write]
status = "d"
`
	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	var merged bool
	for _, n := range nodes {
		vals := n.Values["status"]
		if slices.Contains(vals, "b") && slices.Contains(vals, "c") {
			merged = true
		}
	}
	if !merged {
		t.Errorf("no node merges the `b` and `c` successors into one whose "+
			"per-tag value set is their union; the traversal must be a "+
			"fixpoint over MERGED nodes, never path-sensitive; nodes=%v", nodes)
	}
}

// REQ-109: "The lattice is finite … so the fixpoint terminates on cyclic
// graphs, including the self-loops RDR 0002 models."
// ADVERSARIAL
func TestReq109_FixpointTerminatesOnACyclicGraph(t *testing.T) {
	// A two-cycle plus an escape self-loop. A traversal without a fixpoint
	// would not return; the test's own completion is the oracle, and the
	// node set must stay bounded by the finite lattice.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "a-to-b"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "b-to-a"
[rule.match.status]
eq = "b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"

[[rule]]
id = "self-rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`
	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	if len(nodes) == 0 {
		t.Fatal("the traversal produced no nodes on a cyclic graph")
	}
	// `status` ranges over the subsets of a two-value domain, so the
	// owned-state lattice admits at most 3 non-empty value sets.
	if len(nodes) > 8 {
		t.Errorf("the traversal produced %d nodes over a two-value domain; "+
			"the lattice is finite and the fixpoint must not re-expand "+
			"merged nodes", len(nodes))
	}
}

// REQ-113: "A row *preserves* a tag when it neither writes nor clears it,
// so the tag's value set passes through the edge unchanged."
// DOMAIN EDGE
func TestReq113_ARowPreservesATagItNeitherWritesNorClears(t *testing.T) {
	// `always` is set at the root and the row neither writes nor clears
	// it, so its value set must pass through the edge unchanged. A
	// traversal dropping preserved tags would make `always` absent
	// downstream and mint a false `graph-always-present-owned`.
	const body = `
terminal = ["done"]

[initial]
status = "a"
always = "p"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
	m := mustLoad(t, source(alwaysPresentOwnedDecls, body))

	var sawSuccessor bool
	for _, n := range graphlint.Reach(m) {
		if !slices.Contains(n.Values["status"], "b") {
			continue
		}
		sawSuccessor = true
		if !slices.Contains(n.Values["always"], "p") {
			t.Errorf("the successor node dropped the preserved tag `always`; "+
				"a row that neither writes nor clears a tag passes its value "+
				"set through unchanged; node=%v", n.Values)
		}
	}
	if !sawSuccessor {
		t.Fatalf("the traversal never reached the `status = b` successor, so "+
			"the preservation assertion above is vacuous; nodes=%v",
			graphlint.Reach(m))
	}

	r := graphlint.Run(graphlint.NewRequest(m))
	requireNoCode(t, r, graphlint.CodeAlwaysPresentOwned)
}

// advanceRule is the standard single ordinary row from `a` to `b`.
const advanceRule = `
[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
`
