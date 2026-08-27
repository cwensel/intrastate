package graphlint_test

// Phase 3b adversarial suite for RDR 0006.
//
// Every test here is anchored in the RDR's own `Trade-offs > Failure Modes`
// section, which enumerates what SILENT failure would look like:
//
//	"Silent failure would mean accepting a model with a blocking invariant
//	 defect, … reporting only the first defect per group so a fixed model
//	 still fails, or treating a model with no declared root as clean because
//	 nothing is reachable — the normative command/CI authority, shared lint
//	 engine, complete-emission clause, missing-root rejection, and
//	 by-citation adoption of RDR 0003's clauses are meant to prevent that."
//
// The three attacks below each drive one of those shapes. They use the
// shared fixture helpers in `fixtures_0006_test.go`: real authored TOML
// through the real `internal/table` loader, and the real `graphlint.Run`.
// Nothing is mocked, so a failure here is a defect in the engine and not in
// a stub.

import (
	"maps"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
)

// --- ADV-1: the successor join collapses divergent edges ------------------

// TestAdvSuccessorJoinCollapsesDivergentEdges is ADVERSARIAL.
//
// REQ-108 (LBD, *Join rule and termination*): "two edges reaching **the same
// successor** produce one node whose per-tag value sets are the union of
// theirs (a lattice widening)". The join is licensed only where two edges
// AGREE on their successor.
//
// `reach.go::successorOf` instead folds EVERY edge leaving a node into one
// successor value, and `joinNodes` drops any key absent on either side
// ("absence dominates"). Two rows out of the root that write DIFFERENT owned
// keys therefore annihilate each other: the reachable set loses both keys
// and stands for a concrete view NO path produces.
//
// The Failure Modes clause this lands on is the first one — "accepting a
// model with a blocking invariant defect". The erased keys make two LIVE
// rows look unreachable, which is the success disposition, and their groups
// are then never coverage-checked at all (`groups.go::checkGroups` skips
// `checkCoverage` for an unreachable group). A whole arm of the model goes
// unproven and lint still returns clean.
func TestAdvSuccessorJoinCollapsesDivergentEdges(t *testing.T) {
	// Root is status=a. `goLeft` writes status=b AND left=l; `goRight`
	// writes status=b AND right=r. Both successors are genuinely reachable
	// and genuinely different: {status=b,left=l} and {status=b,right=r}.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true

[tags.left]
provenance = "owned"
kind = "enum"
domain = ["l"]
single_valued = true

[tags.right]
provenance = "owned"
kind = "enum"
domain = ["r"]
single_valued = true
`
	const body = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "c"

[[rule]]
id = "goLeft"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
left = "l"

[[rule]]
id = "goRight"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "b"
right = "r"

[[rule]]
id = "onlyLeft"
[rule.match.status]
eq = "b"
[rule.match.left]
eq = "l"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "c"

[[rule]]
id = "onlyRight"
[rule.match.status]
eq = "b"
[rule.match.right]
eq = "r"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "c"
`

	m := mustLoad(t, source(decls, body))
	nodes := graphlint.Reach(m)

	// The traversal must retain a node holding `left` and a node holding
	// `right`: each is produced by an edge the runtime genuinely takes, and
	// REQ-108 licenses a join only between edges reaching the SAME
	// successor. Collapsing them erases a reachable owned-state.
	t.Run("reachable set keeps the divergent successors", func(t *testing.T) {
		for _, key := range []string{"left", "right"} {
			if !someNodeHolds(nodes, key) {
				t.Errorf("no reachable node holds the owned tag %q, but the "+
					"row writing it is on a root-to-node path; REQ-108 joins "+
					"only edges reaching the SAME successor, so a divergent "+
					"edge's key must survive the traversal.\nnodes:%s",
					key, renderNodes(nodes))
			}
		}
	})

	// The user-visible consequence: two live rows are declared unreachable.
	// `graph-unreachable-rule` is advisory, so lint returns the success
	// disposition while an entire arm of the model went unproven.
	t.Run("live rows are not declared unreachable", func(t *testing.T) {
		r := graphlint.Run(graphlint.NewRequest(m))
		for _, id := range []string{"onlyLeft", "onlyRight"} {
			if namesRule(r, graphlint.CodeUnreachableRule, id) {
				t.Errorf("row %q took %s, but the runtime reaches its "+
					"selection context; the advisory is the success "+
					"disposition, so the group's coverage is never proved "+
					"and lint reports a clean model anyway.\nreport:%s",
					id, graphlint.CodeUnreachableRule, render(r))
			}
		}
	})
}

// someNodeHolds reports whether any node in the set holds key.
func someNodeHolds(nodes []graphlint.Node, key string) bool {
	for _, n := range nodes {
		if len(n.Values[key]) > 0 {
			return true
		}
	}
	return false
}

// renderNodes formats a node set for a failure message.
func renderNodes(nodes []graphlint.Node) string {
	if len(nodes) == 0 {
		return " (empty)"
	}
	out := ""
	for _, n := range nodes {
		out += "\n  " + nodeString(n)
	}
	return out
}

// nodeString renders one node's held tags, keys sorted so a failure
// message is stable across runs.
func nodeString(n graphlint.Node) string {
	out := "{"
	first := true
	for _, k := range slices.Sorted(maps.Keys(n.Values)) {
		if !first {
			out += " "
		}
		first = false
		out += k + "="
		for i, v := range n.Values[k] {
			if i > 0 {
				out += ","
			}
			out += v
		}
	}
	return out + "}"
}

// --- ADV-2: a withheld claim suppresses the group's overlap finding -------

// TestAdvWithholdingSuppressesOverlap is ADVERSARIAL.
//
// REQ-84 (`0006:C16`, complete-emission clause): "Withholding a group's
// exhaustiveness claim MUST NOT suppress overlap, coverage, or further
// withholding findings for that group." The Failure Modes section names the
// same shape — "reporting only the first defect per group so a fixed model
// still fails" — and the complete-emission clause is listed there as the
// mechanism meant to prevent it.
//
// `groups.go::emitOverlaps` skips a row whose accepted-assignment set is not
// projectable (`if !left.Projectable() { continue }`). A row carrying an
// atom over an OPTIONAL key is non-projectable, so adding one optional-key
// guard atom to each of two genuinely-overlapping rows makes the
// `graph-overlap` finding vanish — while the same two rows overlap loudly
// without it.
//
// The overlap is decidable independently of the withholding: it lives on the
// FINITE dimension `p`, whose domain is fully declared. Declining to decide
// it because a different dimension is unprovable is precisely the
// suppression REQ-84 forbids.
//
// The table pairs a control (no optional key, overlap must fire) with the
// attack (identical overlap plus an optional key, overlap must STILL fire).
// The control is what proves the attack is not asserting a defect that was
// never detectable.
func TestAdvWithholdingSuppressesOverlap(t *testing.T) {
	// `p` is a required, single-valued, finitely-declared enum over three
	// values. r1 accepts {x,y}; r2 accepts {y,z}. They intersect at p=y and
	// neither subsumes the other, so this is overlap and not the
	// redundant-row advisory.
	const baseDecls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.p]
provenance = "owned"
kind = "enum"
domain = ["x", "y", "z"]
single_valued = true
required = true
`
	// `opt` carries no `required` marker, so RDR 0003 defaults it optional.
	// An atom over it makes its row able to refuse `guard_unevaluable`, and
	// makes the row's accepted set non-projectable.
	const optDecl = `
[tags.opt]
provenance = "owned"
kind = "enum"
domain = ["o1", "o2"]
single_valued = true
`
	// The optional key is established at the root so this fixture does not
	// also trip invariant 6 — the attack must isolate the suppression.
	const optInitial = "\nopt = \"o1\"\n"

	body := func(extraInitial, extraGuard string) string {
		return `
terminal = ["done"]

[initial]
status = "a"
p = "x"` + extraInitial + `

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "r1"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.p]
in = ["x", "y"]` + guardOn(extraGuard) + `
[rule.write]
status = "b"

[[rule]]
id = "r2"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.p]
in = ["y", "z"]` + guardOn(extraGuard) + `
[rule.write]
status = "b"
`
	}

	tests := []struct {
		name  string
		decls string
		body  string
	}{
		{
			// CONTROL. No optional key, nothing withheld. The overlap on
			// `p=y` is reported. This arm passing is what makes the
			// adversarial arm meaningful.
			name:  "control: overlap with nothing withheld",
			decls: baseDecls,
			body:  body("", ""),
		},
		{
			// ATTACK. Byte-for-byte the same overlap on `p`, plus one
			// optional-key atom per row. REQ-84 says the withheld claim
			// must not suppress the overlap.
			name:  "attack: same overlap with a withheld claim on the group",
			decls: baseDecls + optDecl,
			body:  body(optInitial, "opt"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := lint(t, tc.decls, tc.body)

			f := withCode(r, graphlint.CodeOverlap)
			if len(f) == 0 {
				t.Fatalf("no %s finding; rows r1 and r2 are both enabled by "+
					"p=y over a fully declared finite domain. REQ-84: "+
					"withholding a group's exhaustiveness claim MUST NOT "+
					"suppress overlap for that group.\nreport:%s",
					graphlint.CodeOverlap, render(r))
			}
			// One pair is one finding naming BOTH rows.
			if len(f) != 1 {
				t.Fatalf("%d %s findings; one overlapping pair is one "+
					"finding.\nreport:%s", len(f), graphlint.CodeOverlap,
					render(r))
			}
			if f[0].Rule != "r1" || f[0].Element != "r2" {
				t.Errorf("overlap finding names rule=%q element=%q; want the "+
					"pair (r1, r2).\nreport:%s", f[0].Rule, f[0].Element,
					render(r))
			}
		})
	}
}

// guardOn renders an optional-key guard atom, or nothing when the arm does
// not carry one.
func guardOn(key string) string {
	if key == "" {
		return ""
	}
	return "\n[rule.guard.all." + key + "]\neq = \"o1\""
}

// --- ADV-3: dead end reads a merged node existentially -------------------

// TestAdvDeadEndExistentialOnMergedNode is ADVERSARIAL.
//
// REQ-36 / REQ-111 (INV 2, LBD *Soundness direction*): dead end is a
// UNIVERSAL check, and "a ∀-claim evaluated against a widened domain gets
// *easier* to satisfy, which is the false-green direction." The RDR's remedy
// is to "**split** the node on terminal-participating keys first, recovering
// exactness", because "partial satisfaction is deliberately not enough —
// accepting it would let a merged node close on a path that has not actually
// terminated."
//
// `analysis.go::checkDeadEnd` splits on `terminalKeys()` and then asks
// `hasOutgoingOrdinaryRow`, which is EXISTENTIAL: `matchSatisfiable` accepts
// a node when SOME held value satisfies each atom. The split therefore
// recovers exactness for the terminal half of the test and leaves the
// outgoing-row half reading the merged node directly.
//
// The gap is a node multi-valued on a key that participates in NO terminal.
// Splitting never touches that key, so a node standing for both
// {phase=p} — which has an exit — and {phase=q} — which does not, and meets
// no terminal — is rescued whole by the single exit.
//
// PHASE 3C DISPOSITION — the miss is REAL and the record ACCEPTS it.
//
// Closing it requires ranging the outgoing-row test over the node's
// non-terminal keys, and REQ-37 forecloses exactly that: "The split is
// bounded by the declared domains of terminal-participating keys only,
// **never the whole lattice**." The wider quantifier was implemented and
// measured against the conforming model, and it manufactures concrete views
// no path produces: in `models/rdr.toml` every `stage = "dropped"` write
// also writes `status = "abandoned"`, but the merge decorrelates the two
// keys, so the cross product invents {stage=dropped, status=draft} and
// mints NINE false `graph-dead-end` findings on a model REQ-122 requires
// lint clean. The record anticipates precisely this: "A merged node holds
// every value some path brings, so a check against it can accuse a path the
// runtime never walks. **A path-sensitive reading would be exponential and
// is rejected**" (LBD, Join rule and termination).
//
// This fixture's node and the conforming model's are structurally identical
// after the terminal split — one singleton terminal key beside one
// multi-valued non-terminal key — so no local rule separates the true
// positive here from the false positives there. The separation needs the
// correlation tracking the record rejects.
//
// The test therefore asserts what the record's bounded split DOES
// guarantee, and pins the accepted miss so it stays visible rather than
// silently rediscovered. See D12 in artifacts/deviations.md, which carries
// `Status: needs author decision` for the REQ-111 / REQ-37 tension.
func TestAdvDeadEndExistentialOnMergedNode(t *testing.T) {
	// `status` is the only terminal-participating key, and it is already a
	// singleton at the node under test, so `splitNode` is a no-op there.
	// `phase` is the multi-valued key, and it participates in no terminal.
	const decls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true

[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true
`
	// Root {status=a, phase=p}. `toP` and `toQ` both reach status=b, with
	// phase=p and phase=q respectively, so the merged node is
	// {status:[b], phase:[p,q]}. Only `fromP` leaves it, and only for
	// phase=p. The terminal is status=c, which neither half meets.
	const body = `
terminal = ["done"]

[initial]
status = "a"
phase = "p"

[context.done]
[context.done.match.status]
eq = "c"

[[rule]]
id = "toP"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"
phase = "p"

[[rule]]
id = "toQ"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "b"
phase = "q"

[[rule]]
id = "fromP"
[rule.match.status]
eq = "b"
[rule.match.phase]
eq = "p"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a"
phase = "p"
`

	m := mustLoad(t, source(decls, body))

	// Establish the premise rather than assuming it: the node under attack
	// really is multi-valued on `phase`, so the merged reading is what the
	// check will see.
	if !someNodeHoldsBoth(graphlint.Reach(m), "phase", "p", "q") {
		t.Fatalf("fixture premise failed: no reachable node holds phase in "+
			"both p and q, so this fixture does not exercise the merged "+
			"reading.\nnodes:%s", renderNodes(graphlint.Reach(m)))
	}

	r := graphlint.Run(graphlint.NewRequest(m))

	// The ACCEPTED MISS, pinned. {status=b, phase=q} satisfies no declared
	// terminal and is the source of no non-escape row, so invariant 2 would
	// accuse it under an unbounded split — but `phase` participates in no
	// terminal, and REQ-37 bounds the split to terminal-participating keys
	// "never the whole lattice". Lint reports nothing here BY THE RECORD'S
	// OWN BOUND, not by an implementation slip.
	//
	// Pinning it is what keeps the miss honest: if a later change widens the
	// split, this assertion fails and forces the REQ-122 census to be re-run
	// against `models/rdr.toml` before the widening is accepted.
	if hasCode(r, graphlint.CodeDeadEnd) {
		t.Errorf("%s was emitted for a node multi-valued on `phase`, a key "+
			"participating in NO terminal. REQ-37 bounds invariant 2's split "+
			"to terminal-participating keys only, never the whole lattice, "+
			"and widening it mints false positives on the conforming model "+
			"REQ-122 requires clean. If this widening is intended, re-run the "+
			"REQ-122 census and update D12.\nreport:%s",
			graphlint.CodeDeadEnd, render(r))
	}

	// What the bounded split DOES guarantee, asserted so this fixture still
	// exercises invariant 2 rather than merely recording an absence: the
	// terminal-participating key IS split exactly, so a node whose `status`
	// half meets no terminal and has no exit is still accused.
	t.Run("the terminal-participating key is still split exactly", func(t *testing.T) {
		// Same shape, but the dead half now lives on `status` — the
		// terminal-participating key — so the bounded split reaches it.
		const deadDecls = `
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c", "d"]
single_valued = true
required = true
`
		const deadBody = `
terminal = ["done"]

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "c"

[[rule]]
id = "toB"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "toD"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.write]
status = "d"

[[rule]]
id = "fromB"
[rule.match.status]
eq = "b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "c"
`
		rr := lint(t, deadDecls, deadBody)
		if !hasCode(rr, graphlint.CodeDeadEnd) {
			t.Errorf("no %s finding; the reachable owned-state {status=d} "+
				"satisfies no declared terminal (terminal is status=c) and is "+
				"the source of no non-escape row. `status` IS "+
				"terminal-participating, so REQ-36's split must separate it "+
				"from the {status=b} half that `fromB` rescues.\nreport:%s",
				graphlint.CodeDeadEnd, render(rr))
		}
	})
}

// someNodeHoldsBoth reports whether some node holds key with both values.
func someNodeHoldsBoth(nodes []graphlint.Node, key, a, b string) bool {
	for _, n := range nodes {
		held := n.Values[key]
		if slices.Contains(held, a) && slices.Contains(held, b) {
			return true
		}
	}
	return false
}
