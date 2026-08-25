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
// below rather than with an enumeration.
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

// key renders a node's canonical identity. Keys and values are both sorted
// at construction, so the rendering is stable and two nodes standing for
// the same owned-state render identically.
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
// owned tags is satisfiable in the source node, producing the node with that
// row's writes applied and clears removed.
//
// The traversal is a syntactic dataflow over declarations, writes, and
// clears — no accessor execution, no runtime trace, no guard evaluation — so
// it is total, and it over-approximates the runtime. Match atoms over
// observed/recognized tags and all guard atoms are NOT pruned: every such
// edge is taken as traversable.
func Reach(m *table.Model) []Node {
	nodes, _ := reach(m)
	return nodes
}

// reach returns the reachable node set and whether the traversal completed
// under the published node ceiling. An incomplete traversal returns what it
// found — the caller reports the ceiling rather than running on.
//
// The traversal is a fixpoint over MERGED nodes, never path-sensitive, and
// the successor relation is what makes that literal rather than incidental:
// every edge leaving one node reaches ONE successor node, whose per-tag
// value sets are the join of what those edges produce. A path-sensitive
// enumeration is therefore unrepresentable here rather than merely avoided,
// and the exponential reading the record rejects has nowhere to arise.
func reach(m *table.Model) (nodes []Node, complete bool) {
	if m == nil || len(m.Initial) == 0 {
		// No declared root: the traversal has nothing to start from. The
		// caller must NOT read the empty set as "nothing reachable,
		// therefore clean" — invariant 1 reports the missing declaration.
		return nil, true
	}

	root := Node{Values: map[string][]string{}}
	for _, tv := range m.Initial {
		if m.Tags[tv.Key].Provenance != table.ProvenanceOwned {
			// The traversal is the OWNED-state graph: a root assignment to
			// any other provenance names no node dimension.
			continue
		}
		root.Values[tv.Key] = canonicalValues(tv.Value)
	}

	nodes = []Node{root}
	// worklist carries the node indices whose successors are not yet
	// settled. A node re-enters it whenever a join widens it, so every node
	// downstream sees the wider value sets.
	worklist := []int{0}
	complete = true

	for len(worklist) > 0 {
		if len(nodes) > nodeCeiling {
			complete = false
			break
		}
		i := worklist[0]
		worklist = worklist[1:]

		next, any := successorOf(m, nodes[i])
		if !any {
			continue
		}

		// The successor folds into an EXISTING node whenever one already
		// stands for every concrete view it stands for, so a cycle closes on
		// the node it came from rather than minting a fresh one each time
		// round. The lattice is finite — a finitely-declared tag ranges over
		// its declared domain's subsets, a tag with no finite domain over
		// {held, absent} — so the widening reaches a fixpoint, self-loops
		// and cycles included.
		j := indexOf(nodes, next)
		if j < 0 {
			nodes = append(nodes, next)
			worklist = append(worklist, len(nodes)-1)
			continue
		}
		merged := joinNodes(nodes[j], next)
		if merged.key() == nodes[j].key() {
			continue
		}
		nodes[j] = merged
		worklist = append(worklist, j)
	}
	return nodes, complete
}

// successorOf joins every edge leaving src into the ONE successor node they
// share, and reports whether any edge left at all.
//
// Escape rows are edges too, and they are SELF-LOOPS: an escape row carries
// neither a write block nor a clear list, so its successor equals its
// source. Modelling them adds no node, but the edge exists — which is what
// keeps a node reachable only through an escape rescue reachable here, and
// what stops invariant 2 accusing the recovery arm and invariant 6 accusing
// the rows downstream of it.
func successorOf(m *table.Model, src Node) (Node, bool) {
	var out Node
	var any bool
	for _, row := range m.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		if !matchSatisfiable(m, src, row) {
			continue
		}
		produced := successor(m, src, row)
		if !any {
			out, any = produced, true
			continue
		}
		out = joinNodes(out, produced)
	}
	return out, any
}

// indexOf finds the node already standing for every concrete view want
// stands for, or -1 when none does.
func indexOf(nodes []Node, want Node) int {
	for i, n := range nodes {
		if subsumes(n, want) {
			return i
		}
	}
	return -1
}

// subsumes reports whether n stands for every concrete view want stands for.
// The two must hold the same key set — a node holding a key want does not is
// a different owned-state, not a wider one — and n's value sets must cover
// want's.
func subsumes(n, want Node) bool {
	if len(n.Values) != len(want.Values) {
		return false
	}
	for k, vals := range want.Values {
		held, ok := n.Values[k]
		if !ok {
			return false
		}
		for _, v := range vals {
			if !slices.Contains(held, v) {
				return false
			}
		}
	}
	return true
}

// joinNodes is the lattice join over two edges reaching the same successor:
// per-tag value sets union, and a tag ABSENT on either side is absent in the
// join.
//
// Absence dominates deliberately. The node stands for every concrete view
// some path reaches it with, and a key one path never established is not
// held on that view — recording it as held would let invariant 6 certify a
// read the runtime finds unavailable, which is the false-green direction.
func joinNodes(a, b Node) Node {
	out := Node{Values: map[string][]string{}}
	for k, left := range a.Values {
		right, held := b.Values[k]
		if !held {
			continue
		}
		out.Values[k] = canonicalValues(append(slices.Clone(left), right...))
	}
	return out
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
			continue
		}
		// A write REPLACES: it assigns the tag's whole value and supplants
		// whatever was held, so the successor's set is the written value
		// alone rather than a union with the source's (RDR 0002).
		out.Values[w.Key] = canonicalValues(w.Value)
	}
	return out
}

// writesOf returns the row's write block plus any next-state tag it does not
// already name. The two are populated independently and one is never an
// alias of the other, so the union is what an edge applies.
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

// matchSatisfiable reports whether the row's match pattern over OWNED tags is
// satisfiable in the node.
//
// Match atoms over observed and recognized tags are not pruned — they name
// values that arrive at runtime, so every such edge is taken as traversable —
// and no guard atom is consulted at all.
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
		return a.Operator == resolve.OpExists && !existsWantsPresent(a)
	}
	if a.Operator == resolve.OpExists {
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
	return len(a.Literal) == 1 && a.Literal[0] == resolve.LiteralTrue
}

// atomAdmitsValue decides one match atom against one held value, using RDR
// 0003's evaluator so the traversal never defines a second value semantics.
// An UNEVALUABLE verdict is taken as satisfiable: the relation
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
// seam compares over (JDR 0001 §D13), which is the form the evaluator's `in`
// and `contains` arms parse.
func renderSetLiteral(members []string) string {
	canonical := canonicalValues(members)
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
