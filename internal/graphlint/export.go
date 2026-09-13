package graphlint

import (
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/table"
)

// Edge is one reachability edge: the rule that traverses it, and the node
// keys of the owned-states it leaves and arrives at.
//
// The identity is the whole `(From, To, Rule)` triple (`0021:D-identity`).
// From and To are node keys — the same rendering `NodeKey` publishes — so
// an edge is self-contained against the node set travelling beside it.
type Edge struct {
	From string
	To   string
	Rule string
}

// NodeKey renders a node's canonical identity, the injective rendering the
// traversal keys its fixpoint on.
//
// It exists because the key is what an exporter must publish as a node's
// id, and the rendering is `key`'s alone: a second one in another package
// would have to re-derive the escaping, and two renderings that drift make
// the exported edge endpoints stop naming the exported nodes.
func NodeKey(n Node) string { return n.key() }

// ReachGraph runs the owned-state traversal and returns the reachable node
// set, the edges between those nodes, and whether the traversal completed
// under the published node ceiling.
//
// It is the sibling `Reach` cannot be: `Reach` discards both the edges and
// the `complete` bool, and widening it would change a shipped signature
// (`0021:C4`, A8). Both functions reach the ONE private `reach` fixpoint,
// so the exported relation is the same relation lint certifies rather than
// a second traversal that could drift from it.
//
// The edges are recovered over the SETTLED node set rather than recorded
// during the widening. An edge discovered mid-fixpoint names a node that
// may still widen, so its endpoint key would be stale; re-walking the
// settled nodes names every endpoint by its final key. The walk runs the
// same `successorsOf` the fixpoint does, so no second traversal semantics
// is introduced — it enumerates, it does not explore.
//
// An incomplete traversal returns what it found with `complete` false. The
// caller reports the ceiling; this function never returns a partial graph
// as though it were whole.
func ReachGraph(m *table.Model) (nodes []Node, edges []Edge, complete bool) {
	nodes, complete = reach(m)

	// Every endpoint must resolve to a node the caller also receives, so the
	// index is built over the settled set and a successor that subsumes into
	// an existing node is named by THAT node's key.
	edges = []Edge{}
	for _, src := range nodes {
		from := src.key()
		for _, e := range edgesFrom(m, src) {
			at := indexOf(nodes, e.to)
			if at < 0 {
				// The successor settled into no node in the fixpoint's final
				// set. It cannot be published as an endpoint without naming a
				// node the document does not carry, so the edge is dropped
				// rather than dangled.
				continue
			}
			edges = append(edges, Edge{
				From: from,
				To:   nodes[at].key(),
				Rule: e.rule,
			})
		}
	}

	slices.SortFunc(edges, compareEdges)
	edges = slices.Compact(edges)
	return nodes, edges, complete
}

// producedEdge pairs a successor node with the rule that produced it.
type producedEdge struct {
	to   Node
	rule string
}

// edgesFrom enumerates the edges leaving one node, each labelled with the
// rule id that traverses it.
//
// It mirrors `successorsOf`'s admission rules exactly — an ESCAPE row is
// skipped before satisfiability is asked, and a non-escape row whose match
// pattern over owned tags is unsatisfiable in the source is not an edge —
// but keeps the rule id `successorsOf` drops when it joins successors by
// presence footprint. Two rules reaching one successor are two EDGES even
// though they are one node, which is what lets a reader see which rule
// traverses which transition.
//
// The escape skip is what makes the exported relation EQUAL to the one the
// fixpoint took, not merely a plausible superset: `Reach`'s contract fixes
// an edge as a normalized NON-ESCAPE row, so admitting an escape row as a
// self-loop here would publish a transition the traversal never took (A2's
// equality bar, and the verb/lint drift `0021:C4` makes structural).
func edgesFrom(m *table.Model, src Node) []producedEdge {
	var out []producedEdge
	for _, row := range m.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		if !matchSatisfiable(m, src, row) {
			continue
		}
		out = append(out, producedEdge{to: successor(m, src, row), rule: row.RuleID})
	}
	return out
}

// compareEdges is C2's edge order: (from, to, rule), so construction order
// is unobservable on the wire.
func compareEdges(a, b Edge) int {
	if c := strings.Compare(a.From, b.From); c != 0 {
		return c
	}
	if c := strings.Compare(a.To, b.To); c != 0 {
		return c
	}
	return strings.Compare(a.Rule, b.Rule)
}
