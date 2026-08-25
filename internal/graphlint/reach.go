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
//
// Every field is escaped, so the rendering is also INJECTIVE: two nodes
// standing for DIFFERENT owned-states never render identically. That
// matters twice over — this key is the `seen` de-duplication identity both
// terminal walks in `analysis.go` close on, where a collision suppresses a
// finding outright, and it is the fixpoint equality in `reach`, where a
// collision can terminate the widening early. Both are the false-green
// direction REQ-106 forbids.
func (n Node) key() string {
	keys := slices.Sorted(maps.Keys(n.Values))
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(escapeField(k))
		b.WriteString("=")
		b.WriteString(escapeJoin(n.Values[k], ","))
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
// the node identity is what makes that literal rather than incidental: a
// node is keyed by its owned-state, not by the path that reached it, so
// every path arriving at one owned-state arrives at ONE node. Edges
// reaching the SAME successor join into it; edges reaching different
// successors stay distinct, which is all REQ-108 licenses. A path-sensitive
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
		root.Values[tv.Key] = heldValues(m, tv.Key, tv.Value)
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

		// Each successor folds into an EXISTING node whenever one already
		// stands for every concrete view it stands for, so a cycle closes on
		// the node it came from rather than minting a fresh one each time
		// round. The lattice is finite — a finitely-declared tag ranges over
		// its declared domain's subsets, a tag with no finite domain over
		// {held, absent} — so the widening reaches a fixpoint, self-loops
		// and cycles included.
		for _, next := range successorsOf(m, nodes[i]) {
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
	}
	return nodes, complete
}

// successorsOf returns the distinct successor nodes the edges leaving src
// reach.
//
// Two edges reach THE SAME SUCCESSOR — and are joined, per REQ-108 — when
// the nodes they produce hold the same OWNED TAGS. A node is an abstract
// owned-state, "per owned tag: absent, or held with its set of possible
// declared values" (REQ-101), so the presence footprint IS the successor's
// identity and the per-tag value sets are what the lattice widening unions
// over it. Two edges writing different values to the same tag therefore
// still converge on one node whose value set is their union, which is the
// merged, never-path-sensitive fixpoint REQ-108 requires.
//
// Two edges establishing DIFFERENT owned tags are not joined, because they
// do not reach the same successor. Folding them together anyway — a
// functional successor relation, one successor per source node — is unsound
// in the false-green direction: `joinNodes` drops a key absent on either
// side, so two rows writing different keys annihilate each other and the
// relation stands for a concrete view no path produces. That
// under-approximates presence, against REQ-106's over-approximation
// contract and REQ-110's premise that a merged node admits a SUPERSET of
// concrete views. Grouping by the presence footprint first is what makes
// the join total: within one group every key is held on both sides, so the
// absence-dominates arm never fires and no reachable key is erased.
//
// This stays a fixpoint over MERGED nodes and never path-sensitive: a node
// is keyed by its owned-state, not by the path that reached it, so every
// path arriving at one owned-state arrives at one node. Successor
// multiplicity is a property of the transition relation, not of path
// sensitivity — a DFA state has many successors and enumerates no paths —
// and the lattice stays finite, so the widening still terminates.
//
// Escape rows are edges too, and they are SELF-LOOPS: an escape row carries
// neither a write block nor a clear list, so its successor equals its
// source. Modelling them adds no node, but the edge exists — which is what
// keeps a node reachable only through an escape rescue reachable here, and
// what stops invariant 2 accusing the recovery arm and invariant 6 accusing
// the rows downstream of it.
func successorsOf(m *table.Model, src Node) []Node {
	// index maps a successor's presence footprint to its position in out,
	// so edges reaching the same successor join into one node while edges
	// establishing a different owned tag set stay apart.
	index := map[string]int{}
	var out []Node
	for _, row := range m.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		if !matchSatisfiable(m, src, row) {
			continue
		}
		produced := successor(m, src, row)
		id := presenceKey(produced)
		if at, seen := index[id]; seen {
			out[at] = joinNodes(out[at], produced)
			continue
		}
		index[id] = len(out)
		out = append(out, produced)
	}
	return out
}

// presenceKey renders a node's presence footprint: the sorted set of owned
// tags it holds, with no value sets. It is the successor identity the join
// groups on — two nodes sharing it stand for the same abstract owned-state
// up to the widening REQ-108 licenses.
//
// The tag names are escaped, because nothing upstream forbids a tag name
// carrying the `;` this joins on: `internal/table` validates a tag
// declaration's provenance, kind, per-kind fields, and bounds, and imposes
// no charset. Unescaped, the footprints {"a;b"} and {"a", "b"} both render
// `a;b;`, two DIFFERENT owned-states land in one join group, and
// `joinNodes`'s absence-dominates arm erases every key held on only one
// side — a reachable owned-state lint can no longer see, against REQ-106's
// over-approximation contract.
func presenceKey(n Node) string {
	var b strings.Builder
	for _, k := range slices.Sorted(maps.Keys(n.Values)) {
		b.WriteString(escapeField(k))
		b.WriteString(";")
	}
	return b.String()
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
		out.Values[w.Key] = heldValues(m, w.Key, w.Value)
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

// structuralDelimiters are the runes the composite keys in this package and
// in `engine.go` join on. A field carrying one of them, written raw, would
// let two different composites render identically.
const structuralDelimiters = ";=,|#"

// escapeField makes one authored string safe to write into a
// delimiter-joined composite key. It is the ONE escape in the package, so
// no renderer can diverge from another: a sort key and the composite it
// orders must escape identically or the order stops matching the rendering.
//
// The encoding is backslash escaping — `\` first, then each structural
// delimiter — which is injective because the escape is prefix-free and
// reversible. Two alternatives were rejected on REQ-87, which requires the
// fingerprint be "readable, never a hash, sortable": a length prefix
// destroys lexicographic sortability outright (`10:` sorts before `2:`),
// and JSON re-encodes into a form no longer readable at a glance.
//
// Escaping is NOT order-preserving against an arbitrary neighbour, and no
// caller may assume it is. It inserts `\` (U+005C) at the delimiter's
// position, so `"a,b"` sorts before `"a0"` as authored (`,` U+002C < `0`
// U+0030) but after it once escaped (`\` U+005C > `0`). Where RDR 0002's
// canonical order is what matters — `compareAtoms` — the comparison is
// therefore over the canonical AUTHORED members, never over this encoding.
// This function's contract is injectivity, not order.
//
// Callers escape AFTER canonicalization — `canonicalValues` sorts and
// compacts the AUTHORED values, so the canonical form stays a property of
// what the author wrote rather than of the encoding.
func escapeField(s string) string {
	if !strings.ContainsAny(s, "\\"+structuralDelimiters) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	for _, r := range s {
		if r == '\\' || strings.ContainsRune(structuralDelimiters, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// escapeJoin escapes each member and TERMINATES it with sep, rather than
// joining members on sep. The members are already canonical; escaping here
// is the last step before they become part of a composite key.
//
// The terminator, not an infix join, is what makes the rendering injective
// over CARDINALITY as well as content. An infix join renders the empty
// sequence and the one-element sequence holding the empty string
// identically — `[]` and `[""]` both give `""`, so `k=;` stands for two
// different owned-states. Both are authorable: a `set` tag declaring no
// `elements` constrains no member, `conform` loops zero times over an empty
// member sequence, and `write = { k = [] }` and `write = { k = [""] }` both
// load. They are not the same owned-state — `[]` satisfies no value atom
// while `[""]` satisfies `eq ""` — so collapsing them lets the `seen` maps
// and the fixpoint equality treat a reachable state as already visited,
// the false-green direction REQ-106 forbids.
//
// Terminating preserves what D16 bought. The form stays readable (`a,b,`),
// and it stays sortable for the same reason the escape does: every member
// is followed by its terminator, so no member's rendering is a prefix of a
// different member sequence's.
func escapeJoin(values []string, sep string) string {
	var b strings.Builder
	for _, v := range values {
		b.WriteString(escapeField(v))
		b.WriteString(sep)
	}
	return b.String()
}

// heldValues normalizes a value a node comes to HOLD for a key. A tag whose
// declaration carries no finite domain — a `scalar`, or an `int` lacking a
// bound — ranges over {held, absent} rather than over an enumeration, so it
// is carried as the single opaque value rather than as the concrete string
// the author happened to write. Retaining the concrete string would let
// `ownedAtomSatisfiable` PRUNE a value atom over that key, which marks live
// rows unreachable and skips their blocking checks — the false-green
// direction REQ-106 forbids.
//
// `guard.AssignmentCount` is the one authority on whether a declaration
// carries a finite domain, so the abstraction is gated on it rather than on
// the declared kind: an `int` WITH a bound is finite and keeps its concrete
// value, exactly as an enum does.
//
// Only the two sites that INTRODUCE a held value need this — the root and a
// row's write. `joinNodes` unions sets both of whose sides already came
// through here, and the opaque value is idempotent under that union.
func heldValues(m *table.Model, key string, value []string) []string {
	if _, finite := guard.AssignmentCount(m.Tags[key]); !finite {
		return []string{OpaqueValue}
	}
	return canonicalValues(value)
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
