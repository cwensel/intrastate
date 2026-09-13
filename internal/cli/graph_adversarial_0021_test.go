package cli

// RDR 0021 — Phase 3b adversarial probes.
//
// The Failure Modes section claims the export's silent failures are all
// either guarded or accepted-and-marked: "nondeterministic bytes across
// runs — caught by the MVV's double-emit byte compare, never shipped
// silently; a partial graph after an incomplete traversal — structurally
// impossible, C4 refuses instead."
//
// That enumeration has a hole. It reasons about edges the traversal took
// but FAILED to publish (the ceiling arm), and never about edges the
// export publishes that the traversal NEVER TOOK. A2 states the bar in the
// other direction and states it as equality, not plausibility: the
// recovered edge list must yield "EXACTLY the edges the traversal itself
// took — over-approximation by merged re-run is the premortem's central
// defect (P-1), so equality, not plausibility, is the bar."
//
// An edge the traversal never took is silent by construction: it carries a
// real rule id and two real node ids, so every shipped oracle accepts it.
// `TestReq25_ReachEdgesCarryFromToAndRule` checks an edge's SHAPE and that
// its endpoints are nodes the document also carries; `TestReq26And107_...`
// checks the sort ORDER; the DOT arm compares itself against the same
// `reach` block. A phantom edge satisfies all three. It surfaces only
// against the traversal's own admission rule, which is what this file
// asserts.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// escapeRowModel carries one ordinary advancing row and one ESCAPE row
// whose match pattern is satisfiable at the root.
//
// `Reach`'s contract fixes what an edge is: "an edge is a normalized
// NON-ESCAPE row whose match pattern over owned tags is satisfiable in the
// source node, producing the node with that row's writes applied and
// clears removed." `successorsOf` implements exactly that — it skips every
// escape row before asking whether the match is satisfiable — so the
// fixpoint never traverses one and no escape row contributes a node.
const escapeRowModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "escloop"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

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

[[rule]]
id = "rescue"
escape = ["no_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
`

// bareEscapeModel carries a catch-all escape row constrained only on the
// recognized outcome, so its match pattern over OWNED tags is satisfiable
// in every node. It scales the defect: the phantom edge count grows with
// the node count rather than staying a fixed cost.
const bareEscapeModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "bareesc"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b", "c"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "c"

[[rule]]
id = "a-to-b"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "b-to-c"
[rule.match.status]
eq = "b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "c"

[[rule]]
id = "catchall"
escape = ["no_match"]
[rule.match.recognized]
eq = "stop"
`

// escapeRuleIDs returns the rule ids of every ESCAPE row in the model, as
// the production loader classifies them.
//
// The classification is read off `table.Row.Kind()` — the same predicate
// `successorsOf` consults — rather than re-derived from the authored
// `escape` list, so this helper cannot disagree with the traversal about
// which rows are escape rows.
func escapeRuleIDs(t *testing.T, modelSrc string) map[string]bool {
	t.Helper()

	m := mustLoadModelForGraph(t, modelSrc)
	out := map[string]bool{}
	for _, r := range m.Rows {
		if r.Kind() == table.KindEscape {
			out[r.RuleID] = true
		}
	}
	if len(out) == 0 {
		t.Fatalf("fixture bug: no row classifies as %v, so the probe would "+
			"assert nothing", table.KindEscape)
	}
	return out
}

// ADV-1: the exported edge relation carries edges the traversal never took.
//
// A2: the recovered edge list yields "EXACTLY the edges the traversal
// itself took ... equality, not plausibility, is the bar."
// `0021:C4`: the export reaches the traversal "through the same `graphlint`
// entry surface lint uses ... so verb/lint drift is structural, not
// disciplinary."
//
// `successorsOf` skips escape rows outright, so the fixpoint takes no edge
// from one. The export's own edge walk instead admits each satisfiable
// escape row as a SELF-LOOP, so the document publishes a relation strictly
// larger than the one lint certifies — the drift C4 says is structurally
// impossible.
//
// ADVERSARIAL — the phantom edge is well-formed, correctly sorted, and
// names two real nodes, so every shipped oracle accepts it.
func TestAdv1_TheExportedEdgeSetCarriesNoEdgeTheTraversalNeverTook(t *testing.T) {
	requireGraphVerb(t)

	for _, tc := range []struct {
		name  string
		model string
	}{
		{"escape row satisfiable at the root", escapeRowModel},
		{"catch-all escape row satisfiable everywhere", bareEscapeModel},
	} {
		t.Run(tc.name, func(t *testing.T) {
			escapes := escapeRuleIDs(t, tc.model)
			doc := decodeDocument(t, exportDocument(t, tc.model))
			if doc.Reach == nil {
				t.Fatal("the document carries no `reach` block")
			}

			for _, e := range doc.Reach.Edges {
				if !escapes[e.Rule] {
					continue
				}
				t.Errorf("the exported relation carries the edge %q -> %q "+
					"labelled with the ESCAPE rule %q, but the traversal "+
					"never took it: `successorsOf` skips every escape row, "+
					"so no escape row contributes an edge to the fixpoint. "+
					"A2 fixes the bar as EXACTLY the edges the traversal "+
					"itself took — equality, not plausibility — and C4 "+
					"requires the exported relation to be the one lint "+
					"certifies. `Reach` states the same rule: an edge is a "+
					"normalized NON-ESCAPE row. A consumer proving a "+
					"universal claim over this relation reasons about a "+
					"transition the certified graph does not contain",
					e.From, e.To, e.Rule)
			}
		})
	}
}

// ADV-2: the phantom edge count scales with the node set.
//
// A catch-all escape row is satisfiable in every node, so the export
// publishes one phantom self-loop PER NODE. This pins the blast radius
// rather than the single instance: a reviewer diffing the exported graph
// across a commit that merely adds a state sees the phantom set grow with
// it, and `0021:C2`'s DOT arm renders every one of them into the diagram a
// reviewer reads.
//
// ADVERSARIAL — a fixed one-off would be a wart; growth with the model is
// what makes the exported relation misleading at scale.
func TestAdv2_NoPhantomSelfLoopIsMintedPerReachableNode(t *testing.T) {
	requireGraphVerb(t)

	escapes := escapeRuleIDs(t, bareEscapeModel)
	doc := decodeDocument(t, exportDocument(t, bareEscapeModel))
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	phantom := 0
	for _, e := range doc.Reach.Edges {
		if escapes[e.Rule] && e.From == e.To {
			phantom++
		}
	}
	if phantom > 0 {
		t.Errorf("the exported relation carries %d self-loop(s) minted by "+
			"escape rows over %d reachable node(s) — one per node, so the "+
			"phantom set grows with the model rather than staying a fixed "+
			"cost. The traversal took none of them: `successorsOf` skips "+
			"every escape row before it asks whether the match is "+
			"satisfiable, so the fixpoint's edge set contains no escape "+
			"edge at all. C2's DOT arm renders each phantom into the "+
			"diagram, so the misreading reaches a human reviewer too",
			phantom, len(doc.Reach.Nodes))
	}
}
