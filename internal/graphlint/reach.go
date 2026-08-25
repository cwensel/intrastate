package graphlint

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// Node is one abstract owned-state: per owned tag, either absent, or held
// with the set of possible declared values. A tag with no finite domain
// abstracts to held/absent — it is carried with the single opaque value
// below rather than with an enumerated set.
type Node struct {
	// Values maps each held owned tag to its set of possible values,
	// sorted and duplicate-free. A key absent from the map is absent in
	// this node.
	Values map[string][]string
}

// OpaqueValue is the single abstract value a tag with no finite declared
// domain takes in a node. The lattice over such a tag is {held, absent},
// so its held half needs one representative rather than an enumeration.
const OpaqueValue = "<opaque>"

// key renders a node's canonical identity, so the fixpoint can recognize a
// node it has already expanded. Keys and values are both sorted at
// construction, so the rendering is stable.
func (n Node) key() string {
	keys := slices.Sorted(maps.Keys(n.Values))
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(strings.Join(n.Values[k], ","))
		b.WriteString(";")
	}
	return b.String()
}

// clone deep-copies the node so a successor never aliases its source.
func (n Node) clone() Node {
	out := Node{Values: make(map[string][]string, len(n.Values))}
	for k, v := range n.Values {
		out.Values[k] = slices.Clone(v)
	}
	return out
}

// Reach runs the owned-state traversal to a fixpoint over MERGED nodes and
// returns the reachable node set. The root is the declared initial owned
// state; an edge is a normalized non-escape row whose match pattern over
// owned tags is satisfiable in the source node; escape rows are edges too,
// and they are self-loops, so they add no node.
//
// The traversal is a syntactic dataflow over declarations, writes, and
// clears — no accessor execution, no runtime trace, no guard evaluation —
// so it is total, and it over-approximates the runtime. Match atoms over
// observed/recognized tags and all guard atoms are NOT pruned: every such
// edge is taken as traversable.
func Reach(m *table.Model) []Node {
	nodes, _ := reach(m)
	return nodes
}

// reach returns the reachable node set and whether the traversal completed
// under the published node ceiling. An incomplete traversal returns what
// it found — the caller reports the ceiling rather than running on.
func reach(m *table.Model) (nodes []Node, complete bool) {
	if m == nil || len(m.Initial) == 0 {
		// No declared root: the traversal has nothing to start from. The
		// caller must NOT read the empty set as "nothing reachable,
		// therefore clean" — invariant 1 reports the missing declaration.
		return nil, true
	}

	root := Node{Values: map[string][]string{}}
	for _, tv := range m.Initial {
		root.Values[tv.Key] = canonicalValues(tv.Value)
	}

	// The worklist runs to a fixpoint over merged nodes: two edges reaching
	// the same successor produce one node whose per-tag value sets are the
	// union of theirs. Merging is the lattice widening that makes the
	// fixpoint terminate on cyclic graphs, self-loops included.
	index := map[string]int{root.key(): 0}
	nodes = []Node{root}
	complete = true

	for i := 0; i < len(nodes); i++ {
		if len(nodes) > nodeCeiling {
			complete = false
			break
		}
		src := nodes[i]
		for _, row := range m.Rows {
			if row.Kind() == table.KindEscape {
				// An escape row carries neither a write block nor a clear
				// list, so its successor equals its source: modelling it as
				// a self-loop adds no node, but the edge exists, which is
				// what keeps the rows downstream of a rescue reachable.
				continue
			}
			if !matchSatisfiable(m, src, row) {
				continue
			}
			next := successor(m, src, row)
			k := next.key()
			j, seen := index[k]
			if !seen {
				index[k] = len(nodes)
				nodes = append(nodes, next)
				continue
			}
			// Merge into the existing node. A widened node is re-expanded
			// exactly once more, because its key changes and the merged
			// shape is a new entry in the lattice.
			merged, widened := mergeNodes(nodes[j], next)
			if !widened {
				continue
			}
			nodes[j] = merged
			delete(index, k)
			index[merged.key()] = j
			// Re-expand from the widened node: the worklist index rewinds
			// to it so its successors see the wider value sets.
			if j <= i {
				i = j - 1
			}
		}
	}
	return nodes, complete
}

// mergeNodes unions two nodes' per-tag value sets and reports whether the
// result is strictly wider than the first.
func mergeNodes(a, b Node) (Node, bool) {
	out := a.clone()
	var widened bool
	for k, vals := range b.Values {
		before := len(out.Values[k])
		merged := slices.Compact(slices.Sorted(slices.Values(append(slices.Clone(out.Values[k]), vals...))))
		if len(merged) != before {
			widened = true
		}
		out.Values[k] = merged
	}
	return out, widened
}

// successor produces the node an edge yields: the source with the row's
// writes applied and its clears removed. A row PRESERVES a tag when it
// neither writes nor clears it, so that tag's value set passes through
// unchanged.
func successor(m *table.Model, src Node, row table.Row) Node {
	out := src.clone()
	for _, w := range writesOf(row) {
		if isClearValue(w.Value) {
			delete(out.Values, w.Key)
			continue
		}
		if m.Tags[w.Key].Provenance != table.ProvenanceOwned {
			// The traversal is the OWNED-state graph: a write to a
			// non-owned key names no node dimension.
			continue
		}
		// A write REPLACES: it assigns the tag's whole value and supplants
		// whatever was held, so the successor's set is the written value
		// alone rather than a union with the source's (RDR 0002).
		out.Values[w.Key] = canonicalValues(w.Value)
	}
	return out
}

// writesOf returns the row's write block, falling back to its next-state
// tags. The two are populated independently and one is never an alias of
// the other, so the union is what an edge applies.
func writesOf(row table.Row) []table.TagValue {
	out := slices.Clone(row.Writes)
	for _, n := range row.NextTags {
		if !slices.ContainsFunc(out, func(t table.TagValue) bool { return t.Key == n.Key }) {
			out = append(out, n)
		}
	}
	slices.SortFunc(out, func(a, b table.TagValue) int { return strings.Compare(a.Key, b.Key) })
	return out
}

// isClearValue reports whether a normalized value is the clear sentinel.
func isClearValue(value []string) bool {
	return len(value) == 1 && value[0] == table.ClearSentinel
}

// canonicalValues sorts and de-duplicates a value sequence so node keys are
// stable regardless of authored member order.
func canonicalValues(value []string) []string {
	return slices.Compact(slices.Sorted(slices.Values(value)))
}

// matchSatisfiable reports whether the row's match pattern over OWNED tags
// is satisfiable in the node.
//
// Match atoms over observed and recognized tags are not pruned — they name
// values that arrive at runtime, so every such edge is taken as
// traversable — and no guard atom is consulted at all.
func matchSatisfiable(m *table.Model, n Node, row table.Row) bool {
	for _, a := range row.Atoms {
		if a.Block != table.BlockMatch {
			continue
		}
		if m.Tags[a.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		if !ownedAtomSatisfiable(n, a) {
			return false
		}
	}
	return true
}

// ownedAtomSatisfiable reports whether SOME value the node holds for the
// atom's key satisfies it. This is an existential test over a merged node,
// which is the safe direction: a merged node admits a superset of concrete
// views, so a witness real at a concrete view is still visible here.
func ownedAtomSatisfiable(n Node, a table.Atom) bool {
	held := n.Values[a.Key]
	if len(held) == 0 {
		// The key is absent in this node. An `exists = false` atom is
		// satisfied by exactly that; every other atom needs a held value.
		return a.Operator == "exists" && !existsWantsPresent(a)
	}
	if a.Operator == "exists" {
		return existsWantsPresent(a)
	}
	if slices.Contains(held, OpaqueValue) {
		// A tag with no finite declared domain abstracts to held/absent, so
		// any value atom over it is satisfiable wherever it is held.
		return true
	}
	for _, v := range held {
		if atomAdmitsValue(a, v) {
			return true
		}
	}
	return false
}

// existsWantsPresent reports whether an `exists` atom selects the present
// half of its key's presence dimension.
func existsWantsPresent(a table.Atom) bool {
	return len(a.Literal) == 1 && a.Literal[0] == "true"
}

// atomAdmitsValue decides one match atom against one held value, using RDR
// 0003's evaluator so the traversal never defines a second value
// semantics. An UNEVALUABLE verdict is taken as satisfiable: the relation
// over-approximates, and pruning an edge lint cannot decide is the
// false-green direction.
func atomAdmitsValue(a table.Atom, held string) bool {
	literal := strings.Join(a.Literal, "")
	if a.Operator == "in" || a.Operator == "contains" {
		literal = renderSetLiteral(a.Literal)
	}
	verdict := guard.Evaluator{}.Evaluate(resolve.GuardAtom{
		Key:      a.Key,
		Operator: a.Operator,
		Literal:  literal,
		Block:    a.Block,
	}, held)
	return verdict != resolve.GuardFalse
}

// renderSetLiteral spells a member sequence as the canonical JSON array the
// seam compares over (JDR 0001 §D13), which is the form the evaluator's
// `in` and `contains` arms parse.
func renderSetLiteral(members []string) string {
	canonical := slices.Compact(slices.Sorted(slices.Values(members)))
	if canonical == nil {
		canonical = []string{}
	}
	b, err := json.Marshal(canonical)
	if err != nil {
		// json.Marshal of a []string cannot fail.
		return "[]"
	}
	return string(b)
}
