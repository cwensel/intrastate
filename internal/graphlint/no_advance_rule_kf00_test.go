package graphlint_test

// Kata kf00, graph half — a non-advancing row (`advance = false`) is a
// MODELED SELF-LOOP in the owned-state graph.
//
// It carries neither a write block nor a clear list, so `successor` yields
// the source node unchanged, exactly as an escape row does. The two are
// NOT the same to lint, and the difference is what these tests pin:
//
//   - It is an ORDINARY row, so invariant 2 counts it as an outgoing edge
//     and the node is not `graph-dead-end`. `0006:REQ-103` excludes ESCAPE
//     self-loops from that test because an escape is a resolver RESCUE,
//     not a modeled decision; a non-advancing row is a decision the author
//     wrote down and bound to an outcome.
//   - It must NOT satisfy a declared terminal. Terminal satisfaction is a
//     PREDICATE over the node's owned values, and a row that changes none
//     of them leaves that predicate exactly as it found it.
//   - It must not close coverage through `graph-coverage-closed-by-escape`:
//     that finding keys on `Kind() == KindEscape`, which a non-advancing
//     row is not.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
)

// kf00NoAdvanceBody is the mixed model the kata's consumer needs: one row
// that ADVANCES `a -> b` and one that decides at `a` and stays there.
// `b` is the declared terminal, so `a` is non-terminal and is the node
// invariant 2 tests.
const kf00NoAdvanceBody = `
terminal = ["done"]

[initial]
status = "a"

[emit.op]
kind = "enum"
domain = ["stopped:no-gate"]

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

[[rule]]
id = "refuse"
advance = false
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.emit]
op = "stopped:no-gate"
`

// TestKataKf00_NonAdvancingRowIsAModeledEdgeNotADeadEnd pins the
// invariant-2 half: the self-loop is an ordinary row, so a node whose only
// exit is one is NOT reported `graph-dead-end`.
func TestKataKf00_NonAdvancingRowIsAModeledEdgeNotADeadEnd(t *testing.T) {
	// The discriminating fixture drops the advancing row entirely, so
	// `status = a`'s ONLY outgoing row is the non-advancing one. With the
	// escape reading it would be a dead end.
	const only = `
terminal = ["done"]

[initial]
status = "a"

[emit.op]
kind = "enum"
domain = ["stopped:no-gate"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "refuse"
advance = false
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.emit]
op = "stopped:no-gate"
`
	r := lint(t, statusOnlyDecls, only)
	if hasCode(r, graphlint.CodeDeadEnd) {
		t.Errorf("%s reported: a non-advancing row is a MODELED decision "+
			"bound to an outcome, not a resolver rescue, so REQ-103's "+
			"escape exclusion does not reach it; report:%s",
			graphlint.CodeDeadEnd, render(r))
	}
	// Invariant 7 reads the same "no outgoing non-escape row" test, so it
	// must agree — a model authoring the decision explicitly is not
	// relying on an INFERRED terminal.
	if hasCode(r, graphlint.CodeTerminalEscape) {
		t.Errorf("%s reported over an explicitly authored non-advancing row; "+
			"report:%s", graphlint.CodeTerminalEscape, render(r))
	}
}

// TestKataKf00_NonAdvancingRowSatisfiesNoTerminalAndOpensNoEscapeClosure
// pins the two false-green directions: the self-loop must not stand in for
// terminal satisfaction, and it must not read as an escape row anywhere.
func TestKataKf00_NonAdvancingRowSatisfiesNoTerminalAndOpensNoEscapeClosure(t *testing.T) {
	r := lint(t, statusOnlyDecls, kf00NoAdvanceBody)

	// The terminal predicate is `status = b`. The non-advancing row leaves
	// `status = a`, so it changes nothing the predicate reads; the
	// ADVANCING row is what reaches the terminal. Reporting a dead end
	// here would mean the self-loop had been discounted.
	if hasCode(r, graphlint.CodeDeadEnd) {
		t.Errorf("%s reported on a model whose non-terminal node has both an "+
			"advancing and a non-advancing exit; report:%s",
			graphlint.CodeDeadEnd, render(r))
	}
	// The row declares no failure class, so it is not an escape row and
	// cannot close a coverage arm as one.
	if hasCode(r, graphlint.CodeCoverageClosedByEscape) {
		t.Errorf("%s reported: a non-advancing row declares no resolver "+
			"failure class, so it is not the bare-escape closure this "+
			"finding names; report:%s",
			graphlint.CodeCoverageClosedByEscape, render(r))
	}

	// The self-loop mints no node — its successor equals its source — so
	// the reachable set is exactly the two owned states the model names.
	// A third node would mean the traversal had invented one.
	nodes := graphlint.Reach(mustLoad(t, source(statusOnlyDecls, kf00NoAdvanceBody)))
	if len(nodes) != 2 {
		t.Errorf("reached %d nodes; want 2 (`status = a`, `status = b`) — a "+
			"non-advancing row's successor equals its source, so it adds "+
			"no node: %v", len(nodes), nodes)
	}

	// The self-loop keeps its own row REACHABLE: the node it decides at is
	// in the reachable set, so the row is never reported as dead.
	for _, f := range withCode(r, graphlint.CodeUnreachableRule) {
		if f.Rule == "refuse" {
			t.Errorf("the non-advancing row %q is reported unreachable; "+
				"report:%s", f.Rule, render(r))
		}
	}
}

// TestKataKf00_ANonAdvancingRowAtANonTerminalNodeDoesNotRescueADeadEnd is
// the adversarial control on the terminal half: the self-loop must not let
// a node that satisfies NO terminal read as terminated. Invariant 2 clears
// the node because it has a modeled exit, but the node is still
// non-terminal — so a DIFFERENT node with no exit at all is still
// reported, proving the self-loop bought no blanket exemption.
func TestKataKf00_ANonAdvancingRowAtANonTerminalNodeDoesNotRescueADeadEnd(t *testing.T) {
	// `a` decides without advancing AND advances to `b`; `b` satisfies no
	// terminal (`c` is declared) and has no outgoing row at all.
	const body = `
terminal = ["done"]

[initial]
status = "a"

[emit.op]
kind = "enum"
domain = ["stopped:no-gate"]

[context.done]
[context.done.match.status]
eq = "c"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "refuse"
advance = false
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.emit]
op = "stopped:no-gate"
`
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true
`
	r := lint(t, decls, body)
	// `status = b` IS the dead end — it satisfies no terminal and has no
	// exit — and it must be reported EXACTLY once, at that node.
	f := requireOneCode(t, r, graphlint.CodeDeadEnd)
	if f.Element != "node{status=b,}" {
		t.Errorf("%s reported at %s; want node{status=b,} — `status = a` has "+
			"both an advancing and a non-advancing exit, so it is not a dead "+
			"end, and `status = b` has neither; report:%s",
			graphlint.CodeDeadEnd, f.Element, render(r))
	}
}
