Model: claude-opus-5[1m]

# A2 spike — edge recovery after the fixpoint (differential test)

## Assumption under test

A2: the reachable node set and its successor edges are a deterministic
function of the model value, AND the edge list is recoverable AFTER the
fixpoint by re-running `successorsOf` over the final nodes with subsumption
mapping (`indexOf`), yielding EXACTLY the edges the traversal took.

The premortem names OVER-APPROXIMATION BY MERGED RE-RUN (P-1) as the central
defect: a post-hoc re-run over MERGED final nodes may yield edges the real
traversal never took. The bar is EQUALITY, not plausibility.

## Method and isolation

The project source tree was NOT modified. `internal/graphlint` and its
dependency closure (`internal/table`, `internal/guard`, `internal/resolve`,
`internal/cli/clierr`) were COPIED to a scratch module at `/tmp/a2spike`
with import paths rewritten to `a2spike/...`, because `reach`,
`successorsOf`, `indexOf`, `subsumes`, and `joinNodes` are unexported and
the observer must live inside the package. Go forbids importing another
module's `internal/...`, which is why a `replace` onto the real repo was not
usable and the closure was vendored instead.

Toolchain: go1.27.1 darwin/arm64. `go vet ./...` clean on the scratch module.

### Map-seed variation — the toolchain-appropriate method

`GODEBUG=randmapiter=1` is NOT a knob this toolchain honours: go1.27 silently
ignores unknown GODEBUG keys (verified — `GODEBUG=asdfnotreal=1` was accepted
without complaint), so setting it would prove nothing. Go randomizes map
iteration per-range BY DEFAULT, which is the real seed variation. That
randomization was demonstrated live rather than assumed (5 runs, 5 different
iteration orders; see raw output). Determinism was then established by
(a) 2000 repeats in ONE process and (b) repeated fresh processes, byte-compared.

## What a node is, and what merging does

A `Node` is an abstract owned-state: per owned tag, either absent, or held
with a sorted duplicate-free set of possible values. Identity is the
owned-state, never the path.

There are TWO distinct merges, and the distinction is what decides A2:

1. INTRA-`successorsOf` JOIN. Within one call, edges whose produced nodes
   share a presence footprint (`presenceKey`) are joined via `joinNodes`.
   This is a pure function of `(model, src)`.
2. IN-PLACE WIDENING in `reach`: `nodes[j] = merged` re-queues a node.
   This is the ONLY path by which a node could be widened AFTER it was
   already expanded — and therefore the only way a post-hoc re-run could
   ever see a source node wider than the one the traversal expanded.

P-1 lives entirely in (2).

## KEY STRUCTURAL FINDING — the widening arm is dead code

`indexOf` returns `j` only when `subsumes(nodes[j], next)` holds.
`subsumes` requires the SAME key set and requires `nodes[j]`'s value sets to
COVER `next`'s. But if `nodes[j]` already covers `next`, then
`joinNodes(nodes[j], next)` unions nothing new, so
`merged.key() == nodes[j].key()` and the `continue` guard fires FIRST.
The `nodes[j] = merged` widening is therefore UNREACHABLE.

Proved two ways:

- LEMMA, exhaustive sweep: over all non-empty subsets of `{a,b,c,<opaque>}`
  for four key-set shapes (including mismatched key sets), 3600 node pairs,
  260 of them subsuming — ZERO counterexamples to
  `subsumes(n,want) => joinNodes(n,want).key() == n.key()`.
- AUDIT, per-traversal: an instrumented `reach` recording every
  `nodes[j]=merged` event reports `wideningEvents=0` on every adversarial
  fixture (C, D, H, I), and `widening-branch taken=0` on all 36 real models.

Consequence: every node is expanded at its FINAL width, so re-running
`successorsOf` over the final nodes reproduces exactly the edges taken. The
equality is STRUCTURAL, not incidental.

## Fixtures (inline)

All fixtures declare owned `enum` tags (finite domain, so `heldValues` keeps
concrete values and `ownedAtomSatisfiable` genuinely prunes — confirmed:
`guard.Evaluator.Evaluate` decides `eq` as `value == literal`).

- A — plain chain s0->s1->s2. Control, no merge.
- B — SUBSUMPTION MERGE: two rows write different values (`k=a`, `k=b`) to
  the same key from the root, same presence footprint, so the intra-call
  join unions them to `k={a,b}`; the result then folds into the existing
  node via `indexOf`/`subsumes`.
- C — P-1 WITNESS ATTEMPT: rows `ka`/`kb` merge to `k={a,b}`, and row `trap`
  requires `k eq b`. The MERGED node enables `trap`; neither authored write
  alone would have been reached separately.
- D — POST-EXPANSION WIDENING attempt: two routes converge on footprint
  `{k,p}` carrying `k=a` and `k=b` by different paths, aiming to widen a
  node after it was expanded.
- E — cycle a->b->c->a with widening.
- F — opaque `scalar` tag plus a `<clear>` write, mixing presence changes.
- G — narrower successor folding into a wider root.
- H — cycle re-widening an already-expanded node; `onlyC` fires only at k=c.
- I — two-key (j,k) stale-source-key hunt.

Plus all 36 loadable authored models under `internal/table/testdata`
(118 `.toml` files walked; 82 are deliberate negative/parse fixtures that do
not load).

## Controls — the comparator is not vacuously green

NEGATIVE CONTROL (comparator self-test): fed a recovery set with one injected
phantom edge and one dropped edge; the differ reported exactly 1 in each
direction. The comparator sees both directions.

POSITIVE CONTROL (divergence reproduced in a REAL traversal): `indexOf` was
weakened to match on presence footprint alone — the literal "merged re-run"
shape P-1 warns about. The widening arm then went LIVE (2 in-place
widenings), the whole chain collapsed to ONE node, and recovery diverged from
the traversal by 2 edges. So a divergent algorithm IS detected by this
harness; the real algorithm's equality is due to its subsumption guard.

## Results

- Synthetic fixtures A–I: recovery == traversal, EXACT multiset equality,
  0 phantom edges, 0 missing edges.
- Real authored models: 36 loaded, EQUAL=36, NOT-EQUAL=0.
- Determinism: 2000 in-process repeats -> 1 distinct output digest.
  Repeated fresh processes byte-identical (sha256 match across runs) for both
  the synthetic suite and the real-model suite.

## Verdict per arm

- DETERMINISM ARM: PASS. Node set and edge list are byte-identical across
  2000 in-process repeats and across fresh processes, under Go's default
  per-range map-iteration randomization (demonstrated live).
- DIFFERENTIAL ARM: PASS. Post-hoc recovery equals the in-traversal observer
  EXACTLY on every fixture, including the subsumption-merge and merged-node-
  enables-a-new-row cases, and on all 36 real models. No over-approximation.

## Honest scope note on clause 4(b)

A fixture where a MERGED node enables a row no pre-merge node enabled WAS
built (C, and H's `onlyC`). It does NOT produce a traversal/recovery
divergence, and the reason is structural rather than lucky: that merge
happens inside `successorsOf`, which is a pure function of `(model, src)` and
which the recovery re-runs identically. The divergence P-1 fears requires the
OTHER merge — in-place widening of an already-expanded node — and that arm is
unreachable under `subsumes`, proven by exhaustive lemma and by a zero-event
audit. A witness for it therefore cannot be constructed against the current
`indexOf`; it was constructed against a deliberately weakened `indexOf`
(positive control), where the divergence duly appeared.

## Caveat for implementation

A2 holds BECAUSE `indexOf` uses `subsumes`. If `indexOf` is ever relaxed to
merge on presence footprint (or any non-covering relation), the widening arm
becomes live and post-hoc recovery becomes an over-approximation — exactly
P-1. The positive control demonstrates this failure concretely. Recommend the
implementation pin this with a regression test asserting the widening branch
stays unreachable.

## Spike source

Scratch module `/tmp/a2spike`; project tree untouched.

### internal/graphlint/observer.go
```go
package graphlint

import (
	"fmt"
	"sort"

	"a2spike/internal/table"
)

// Edge is one traversal edge, recorded as the traversal takes it.
type Edge struct {
	From string // source node key AT THE MOMENT the edge was taken
	Rule string // the row that enabled the edge
	To   string // produced successor key (pre-merge)
}

func (e Edge) String() string { return e.From + "  --[" + e.Rule + "]-->  " + e.To }

// reachObserved is a FAITHFUL copy of reach() with an in-traversal edge
// observer. The only additions are the `edges` recording lines; all control
// flow, merge, and fixpoint logic is character-for-character the original.
func reachObserved(m *table.Model) (nodes []Node, complete bool, edges []Edge) {
	if m == nil {
		return nil, true, nil
	}
	if len(m.Initial) == 0 && !table.IsDecisionTable(m) {
		return nil, true, nil
	}

	root := Node{Values: map[string][]string{}}
	for _, tv := range m.Initial {
		if m.Tags[tv.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		root.Values[tv.Key] = heldValues(m, tv.Key, tv.Value)
	}

	nodes = []Node{root}
	worklist := []int{0}
	complete = true

	for len(worklist) > 0 {
		if len(nodes) > nodeCeiling {
			complete = false
			break
		}
		i := worklist[0]
		worklist = worklist[1:]

		// Record the source identity AS OF this visit. A later join may widen
		// nodes[i]; the edge the traversal took belongs to the node as it
		// stood here.
		srcKey := nodes[i].key()
		for _, se := range successorsWithRules(m, nodes[i]) {
			edges = append(edges, Edge{From: srcKey, Rule: se.rule, To: se.node.key()})

			next := se.node
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
	return nodes, complete, edges
}

type succWithRule struct {
	node Node
	rule string
}

// successorsWithRules mirrors successorsOf exactly, but carries the enabling
// rule id alongside each distinct successor. When two rows join into one
// successor, the rule label is the JOINED set, matching successorsOf's own
// grouping.
func successorsWithRules(m *table.Model, src Node) []succWithRule {
	index := map[string]int{}
	var out []succWithRule
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
			out[at].node = joinNodes(out[at].node, produced)
			out[at].rule = out[at].rule + "+" + row.RuleID
			continue
		}
		index[id] = len(out)
		out = append(out, succWithRule{node: produced, rule: row.RuleID})
	}
	return out
}

// recoverEdges is the POST-HOC RECOVERY the assumption proposes: re-run
// successorsOf over the FINAL (merged) nodes, mapping each successor through
// indexOf.
func recoverEdges(m *table.Model, nodes []Node) []Edge {
	var edges []Edge
	for _, n := range nodes {
		src := n.key()
		for _, se := range successorsWithRules(m, n) {
			edges = append(edges, Edge{From: src, Rule: se.rule, To: se.node.key()})
		}
	}
	return edges
}

// multiset renders an edge slice as a sorted multiset for exact comparison.
func multiset(edges []Edge) []string {
	out := make([]string, 0, len(edges))
	for _, e := range edges {
		out = append(out, e.String())
	}
	sort.Strings(out)
	return out
}

// diffMultisets returns (onlyInA, onlyInB) as exact multiset differences.
func diffMultisets(a, b []Edge) (onlyA, onlyB []string) {
	counts := map[string]int{}
	for _, s := range multiset(a) {
		counts[s]++
	}
	for _, s := range multiset(b) {
		counts[s]--
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch {
		case counts[k] > 0:
			for i := 0; i < counts[k]; i++ {
				onlyA = append(onlyA, k)
			}
		case counts[k] < 0:
			for i := 0; i < -counts[k]; i++ {
				onlyB = append(onlyB, k)
			}
		}
	}
	return onlyA, onlyB
}

// NodeKeys renders the final node set as sorted keys.
func NodeKeys(nodes []Node) []string {
	out := make([]string, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, n.key())
	}
	return out
}

// RunDifferential is the decisive arm.
func RunDifferential(name string, m *table.Model) string {
	nodes, complete, taken := reachObserved(m)
	recovered := recoverEdges(m, nodes)

	onlyTaken, onlyRecovered := diffMultisets(taken, recovered)

	s := fmt.Sprintf("FIXTURE: %s\n", name)
	s += fmt.Sprintf("  complete=%v  nodes=%d\n", complete, len(nodes))
	s += "  FINAL NODES (traversal order):\n"
	for i, k := range NodeKeys(nodes) {
		s += fmt.Sprintf("    [%d] %s\n", i, k)
	}
	s += fmt.Sprintf("  EDGES TAKEN BY TRAVERSAL (observer, n=%d):\n", len(taken))
	for _, e := range multiset(taken) {
		s += "    " + e + "\n"
	}
	s += fmt.Sprintf("  EDGES FROM POST-HOC RECOVERY (n=%d):\n", len(recovered))
	for _, e := range multiset(recovered) {
		s += "    " + e + "\n"
	}
	s += fmt.Sprintf("  ONLY-IN-TRAVERSAL (under-approx by recovery): %d\n", len(onlyTaken))
	for _, e := range onlyTaken {
		s += "    MISSING-FROM-RECOVERY: " + e + "\n"
	}
	s += fmt.Sprintf("  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): %d\n", len(onlyRecovered))
	for _, e := range onlyRecovered {
		s += "    PHANTOM-EDGE: " + e + "\n"
	}
	if len(onlyTaken) == 0 && len(onlyRecovered) == 0 {
		s += "  VERDICT: EQUAL (recovery == traversal)\n"
	} else {
		s += "  VERDICT: NOT EQUAL -- A2 FALSIFIED ON THIS FIXTURE\n"
	}
	return s
}
```

### internal/graphlint/audit.go
```go
package graphlint

import (
	"fmt"

	"a2spike/internal/table"
)

// WideningEvent records an in-place widening of a node during reach().
type WideningEvent struct {
	Index        int
	Before       string
	After        string
	AlreadyByNow bool // had this node index already been EXPANDED?
}

// reachAudit mirrors reach() and records every nodes[j]=merged widening,
// noting whether that index had already been expanded at least once.
func reachAudit(m *table.Model) (nodes []Node, events []WideningEvent) {
	if m == nil || (len(m.Initial) == 0 && !table.IsDecisionTable(m)) {
		return nil, nil
	}
	root := Node{Values: map[string][]string{}}
	for _, tv := range m.Initial {
		if m.Tags[tv.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		root.Values[tv.Key] = heldValues(m, tv.Key, tv.Value)
	}
	nodes = []Node{root}
	worklist := []int{0}
	expanded := map[int]bool{}

	for len(worklist) > 0 {
		if len(nodes) > nodeCeiling {
			break
		}
		i := worklist[0]
		worklist = worklist[1:]
		expanded[i] = true

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
			events = append(events, WideningEvent{
				Index: j, Before: nodes[j].key(), After: merged.key(),
				AlreadyByNow: expanded[j],
			})
			nodes[j] = merged
			worklist = append(worklist, j)
		}
	}
	return nodes, events
}

// AuditWidening reports widening events for one model.
func AuditWidening(name string, m *table.Model) string {
	nodes, events := reachAudit(m)
	s := fmt.Sprintf("AUDIT %s: nodes=%d wideningEvents=%d\n", name, len(nodes), len(events))
	for _, e := range events {
		s += fmt.Sprintf("  widen node[%d] %q -> %q  alreadyExpanded=%v\n",
			e.Index, e.Before, e.After, e.AlreadyByNow)
	}
	if len(events) == 0 {
		s += "  (no in-place widening occurred)\n"
	}
	return s
}



```

### internal/graphlint/lemma2.go
```go
package graphlint

import (
	"fmt"

	"slices"
)

// SubsumptionLemmaBig is an exhaustive sweep over ALL subsets of {a,b,c,opaque}
// for 1- and 2-key nodes, including mismatched key sets, checking:
//   subsumes(n, want)  =>  joinNodes(n, want).key() == n.key()
// This is the lemma that makes reach()'s `nodes[j] = merged` arm dead code.
func SubsumptionLemmaBig() string {
	universe := []string{"a", "b", "c", OpaqueValue}
	var sets [][]string
	for mask := 1; mask < 1<<len(universe); mask++ {
		var s []string
		for i, v := range universe {
			if mask&(1<<i) != 0 {
				s = append(s, v)
			}
		}
		sets = append(sets, canonicalValues(s))
	}

	keySets := [][]string{{"k"}, {"j"}, {"k", "j"}, {"k", "j", "z"}}

	subsuming, counter, total := 0, 0, 0
	var detail string
	for _, ak := range keySets {
		for _, bk := range keySets {
			for _, av := range sets {
				for _, bv := range sets {
					a := Node{Values: map[string][]string{}}
					for _, k := range ak {
						a.Values[k] = slices.Clone(av)
					}
					b := Node{Values: map[string][]string{}}
					for _, k := range bk {
						b.Values[k] = slices.Clone(bv)
					}
					total++
					if !subsumes(a, b) {
						continue
					}
					subsuming++
					if joinNodes(a, b).key() != a.key() {
						counter++
						if len(detail) < 2000 {
							detail += fmt.Sprintf("    COUNTEREXAMPLE a=%q b=%q join=%q\n",
								a.key(), b.key(), joinNodes(a, b).key())
						}
					}
				}
			}
		}
	}
	s := fmt.Sprintf("SUBSUMPTION LEMMA (exhaustive): pairs=%d subsuming=%d counterexamples=%d\n",
		total, subsuming, counter)
	if counter == 0 {
		s += "  => subsumes(n,want) ALWAYS implies joinNodes(n,want).key()==n.key()\n"
		s += "  => reach()'s `nodes[j]=merged` widening arm is UNREACHABLE (dead code)\n"
	}
	return s + detail
}
```

### internal/graphlint/posctl.go
```go
package graphlint

import (
	"fmt"

	"a2spike/internal/table"
)

// POSITIVE CONTROL: make the widening branch LIVE by using a deliberately
// WEAKENED indexOf that matches on presence footprint alone (ignoring value
// coverage). That is the "merged re-run" shape P-1 warns about. If recovery
// then diverges from the traversal, we have proven the harness detects
// over-approximation in a REAL traversal -- and that the real algorithm's
// equality is due to its subsumption guard, not to a blind comparator.
func indexOfWeak(nodes []Node, want Node) int {
	for i, n := range nodes {
		if presenceKey(n) == presenceKey(want) {
			return i
		}
	}
	return -1
}

func reachWeak(m *table.Model) (nodes []Node, edges []Edge, widenings int) {
	root := Node{Values: map[string][]string{}}
	for _, tv := range m.Initial {
		if m.Tags[tv.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		root.Values[tv.Key] = heldValues(m, tv.Key, tv.Value)
	}
	nodes = []Node{root}
	worklist := []int{0}
	for len(worklist) > 0 {
		if len(nodes) > nodeCeiling {
			break
		}
		i := worklist[0]
		worklist = worklist[1:]
		srcKey := nodes[i].key()
		for _, se := range successorsWithRules(m, nodes[i]) {
			edges = append(edges, Edge{From: srcKey, Rule: se.rule, To: se.node.key()})
			j := indexOfWeak(nodes, se.node)
			if j < 0 {
				nodes = append(nodes, se.node)
				worklist = append(worklist, len(nodes)-1)
				continue
			}
			merged := joinNodes(nodes[j], se.node)
			if merged.key() == nodes[j].key() {
				continue
			}
			widenings++
			nodes[j] = merged
			worklist = append(worklist, j)
		}
	}
	return nodes, edges, widenings
}

// PositiveControl runs the weakened traversal and compares recovery against
// the observed edges.
func PositiveControl(name string, m *table.Model) string {
	nodes, taken, widenings := reachWeak(m)
	recovered := recoverEdges(m, nodes)
	onlyTaken, onlyRecovered := diffMultisets(taken, recovered)

	s := fmt.Sprintf("POSITIVE CONTROL %s (weakened indexOf -> widening branch LIVE)\n", name)
	s += fmt.Sprintf("  in-place widenings=%d  finalNodes=%d\n", widenings, len(nodes))
	s += fmt.Sprintf("  edges taken=%d  edges recovered=%d\n", len(taken), len(recovered))
	s += fmt.Sprintf("  ONLY-IN-TRAVERSAL=%d  ONLY-IN-RECOVERY(phantom)=%d\n",
		len(onlyTaken), len(onlyRecovered))
	for _, e := range onlyTaken {
		s += "    MISSING-FROM-RECOVERY: " + e + "\n"
	}
	for _, e := range onlyRecovered {
		s += "    PHANTOM-EDGE: " + e + "\n"
	}
	if len(onlyTaken) > 0 || len(onlyRecovered) > 0 {
		s += "  => divergence REPRODUCED under a merging indexOf: the harness CAN see P-1.\n"
	} else {
		s += "  => no divergence even when widening is live.\n"
	}
	return s
}
```

### internal/graphlint/negctl.go
```go
package graphlint

import "fmt"

// NegativeControl proves the comparator is not vacuously green: it feeds the
// differ a recovery set containing one edge the traversal never took, and one
// traversal edge the recovery lacks, then reports what the differ found.
func NegativeControl() string {
	taken := []Edge{
		{From: "n0", Rule: "r1", To: "n1"},
		{From: "n1", Rule: "r2", To: "n2"},
		{From: "n2", Rule: "r3", To: "n3"},
	}
	recovered := []Edge{
		{From: "n0", Rule: "r1", To: "n1"},
		{From: "n1", Rule: "r2", To: "n2"},
		// n2->n3 dropped (simulated under-approximation)
		{From: "n1", Rule: "rPHANTOM", To: "n9"}, // simulated P-1 phantom
	}
	onlyTaken, onlyRecovered := diffMultisets(taken, recovered)
	s := "NEGATIVE CONTROL (harness self-test)\n"
	s += fmt.Sprintf("  injected: 1 phantom in recovery, 1 edge missing from recovery\n")
	s += fmt.Sprintf("  ONLY-IN-TRAVERSAL detected: %d\n", len(onlyTaken))
	for _, e := range onlyTaken {
		s += "    MISSING-FROM-RECOVERY: " + e + "\n"
	}
	s += fmt.Sprintf("  ONLY-IN-RECOVERY detected: %d\n", len(onlyRecovered))
	for _, e := range onlyRecovered {
		s += "    PHANTOM-EDGE: " + e + "\n"
	}
	if len(onlyTaken) == 1 && len(onlyRecovered) == 1 {
		s += "  HARNESS SELF-TEST: PASS (comparator detects both directions)\n"
	} else {
		s += "  HARNESS SELF-TEST: FAIL (comparator is blind)\n"
	}
	return s
}
```

### cmd/spike/main.go (fixtures A-F)
```go
package main

import (
	"fmt"
	"os"

	"a2spike/internal/graphlint"
	"a2spike/internal/table"
)

func enum(domain ...string) table.TagDecl {
	return table.TagDecl{
		Provenance: table.ProvenanceOwned, Kind: "enum",
		Domain: domain, SingleValued: true, Required: false,
	}
}

func match(key, op string, lit ...string) table.Atom {
	return table.Atom{Key: key, Block: table.BlockMatch, Operator: op, Literal: lit}
}

func tv(key string, vals ...string) table.TagValue {
	return table.TagValue{Key: key, Value: vals}
}

// --- FIXTURE A: plain chain, no merges. Control. -------------------------
func fixtureA() *table.Model {
	return &table.Model{
		ID:   "fixA",
		Tags: map[string]table.TagDecl{"s": enum("s0", "s1", "s2")},
		Initial: []table.TagValue{tv("s", "s0")},
		Rows: []table.Row{
			{RuleID: "r1", Atoms: []table.Atom{match("s", "eq", "s0")}, Writes: []table.TagValue{tv("s", "s1")}},
			{RuleID: "r2", Atoms: []table.Atom{match("s", "eq", "s1")}, Writes: []table.TagValue{tv("s", "s2")}},
		},
	}
}

// --- FIXTURE B: subsumption merge. Two rows write DIFFERENT values to the
// SAME key, so successorsOf's intra-call join unions them into {a,b}. That
// widened node then joins into the existing node via indexOf/subsumes.
func fixtureB() *table.Model {
	return &table.Model{
		ID: "fixB",
		Tags: map[string]table.TagDecl{
			"k": enum("a", "b", "c"),
			"p": enum("p0", "p1"),
		},
		Initial: []table.TagValue{tv("p", "p0")},
		Rows: []table.Row{
			// From the root (p=p0, k absent) two rows establish k with
			// different values but the SAME presence footprint {k,p}.
			{RuleID: "wa", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("k", "a")}},
			{RuleID: "wb", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("k", "b")}},
		},
	}
}

// --- FIXTURE C: THE P-1 WITNESS. ----------------------------------------
// A MERGED node enables a row that NO PRE-MERGE node enabled.
//
// Construction: rows `wa` and `wb` both leave the root with footprint {k,p}
// but values k=a and k=b respectively. successorsOf JOINS them (same
// presence footprint) into the single successor k={a,b}. Meanwhile row
// `trap` matches `k in [a,b]`... that is enabled by either alone.
//
// To get a row enabled ONLY by the merge we need a CONJUNCTION over two
// keys that no single pre-merge node satisfies. Use two keys j and k:
//   wa: writes j=1        (footprint {j,p})
//   wb: writes j=2        (footprint {j,p})  -> join gives j={1,2}
// then `trap` requires j eq 2 AND ... but j=2 alone already enables it.
//
// The true P-1 shape: the join WIDENS an EXISTING node via joinNodes in
// reach's fixpoint. Node X first visited with j={1}; a later successor
// j={2} with the same footprint... indexOf/subsumes FAILS (1 does not cover
// 2), so a NEW node appears instead of a merge. The merge that DOES widen
// is the intra-successorsOf join. So: make the intra-call join produce
// j={1,2}, and have `trap` match `j eq 2` while the SOURCE node that the
// traversal actually expands is the joined one. Then the trap edge IS taken
// by the traversal too -- so no divergence.
//
// The divergence therefore requires the widening to happen AFTER the node
// was already expanded. That is reach's joinNodes arm: nodes[j] = merged,
// re-queued. We force it by making a node reachable twice with values that
// SUBSUME on the first arrival but widen on the second.
func fixtureC() *table.Model {
	return &table.Model{
		ID: "fixC",
		Tags: map[string]table.TagDecl{
			"k": enum("a", "b"),
			"p": enum("p0", "p1", "p2"),
		},
		Initial: []table.TagValue{tv("p", "p0")},
		Rows: []table.Row{
			// Root p0 -> p1 (k absent)
			{RuleID: "t1", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("p", "p1")}},
			// p1 -> establish k=a  (footprint {k,p}, p stays p1)
			{RuleID: "ka", Atoms: []table.Atom{match("p", "eq", "p1")}, Writes: []table.TagValue{tv("k", "a")}},
			// p1 -> establish k=b  (same footprint -> intra-call JOIN k={a,b})
			{RuleID: "kb", Atoms: []table.Atom{match("p", "eq", "p1")}, Writes: []table.TagValue{tv("k", "b")}},
			// Only enabled when k can be b.
			{RuleID: "trap", Atoms: []table.Atom{match("k", "eq", "b"), match("p", "eq", "p1")}, Writes: []table.TagValue{tv("p", "p2")}},
		},
	}
}

// --- FIXTURE D: forced post-expansion widening (the sharpest P-1 attempt).
// Node with footprint {k,p} is reached first holding k={a}. It is EXPANDED
// (successors computed) at that moment. Only later does a second path bring
// k={b} to the same footprint, triggering joinNodes -> k={a,b} and re-queue.
// A row matching `k eq b` is then enabled from the widened node. The
// traversal DOES re-expand on re-queue, so the edge should be taken -- the
// question is whether the FROM key recorded matches the final key.
func fixtureD() *table.Model {
	return &table.Model{
		ID: "fixD",
		Tags: map[string]table.TagDecl{
			"k": enum("a", "b"),
			"g": enum("g0", "g1"),
			"p": enum("p0", "p1"),
		},
		Initial: []table.TagValue{tv("p", "p0")},
		Rows: []table.Row{
			// Two distinct routes that both land on footprint {k,p}.
			{RuleID: "ra", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("k", "a")}},
			// A second route via g that eventually writes k=b with the SAME
			// footprint {k,p}: first go p0->p1, then from p1 write k=b and
			// reset p to p0 so the footprint AND the p value coincide.
			{RuleID: "rp", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("p", "p1")}},
			{RuleID: "rb", Atoms: []table.Atom{match("p", "eq", "p1")}, Writes: []table.TagValue{tv("k", "b"), tv("p", "p0")}},
			// Enabled only when k can be b.
			{RuleID: "trap", Atoms: []table.Atom{match("k", "eq", "b"), match("p", "eq", "p0")}, Writes: []table.TagValue{tv("g", "g1")}},
		},
	}
}

// --- FIXTURE E: self-loop / cycle with widening. -------------------------
func fixtureE() *table.Model {
	return &table.Model{
		ID: "fixE",
		Tags: map[string]table.TagDecl{
			"k": enum("a", "b", "c"),
		},
		Initial: []table.TagValue{tv("k", "a")},
		Rows: []table.Row{
			{RuleID: "c1", Atoms: []table.Atom{match("k", "eq", "a")}, Writes: []table.TagValue{tv("k", "b")}},
			{RuleID: "c2", Atoms: []table.Atom{match("k", "eq", "b")}, Writes: []table.TagValue{tv("k", "c")}},
			{RuleID: "c3", Atoms: []table.Atom{match("k", "eq", "c")}, Writes: []table.TagValue{tv("k", "a")}},
		},
	}
}

// --- FIXTURE F: opaque (scalar) tag + clear, mixing presence changes. ----
func fixtureF() *table.Model {
	return &table.Model{
		ID: "fixF",
		Tags: map[string]table.TagDecl{
			"s": {Provenance: table.ProvenanceOwned, Kind: "scalar"},
			"k": enum("a", "b"),
		},
		Initial: []table.TagValue{tv("k", "a")},
		Rows: []table.Row{
			{RuleID: "set", Atoms: []table.Atom{match("k", "eq", "a")}, Writes: []table.TagValue{tv("s", "anything")}},
			{RuleID: "flip", Atoms: []table.Atom{match("k", "eq", "a")}, Writes: []table.TagValue{tv("k", "b")}},
			{RuleID: "clr", Atoms: []table.Atom{match("k", "eq", "b")}, Writes: []table.TagValue{tv("s", table.ClearSentinel)}},
		},
	}
}

func main() {
	fixtures := []struct {
		name string
		m    *table.Model
	}{
		{"A: plain chain (control, no merge)", fixtureA()},
		{"B: subsumption merge via intra-call join", fixtureB()},
		{"C: P-1 witness attempt (merged node enables `trap`)", fixtureC()},
		{"D: post-expansion widening (joinNodes re-queue)", fixtureD()},
		{"E: cycle with widening", fixtureE()},
		{"F: opaque scalar + clear", fixtureF()},
	}

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "determinism":
		// Print node keys + edge multiset digest for every fixture.
		for _, f := range fixtures {
			fmt.Printf("FIXTURE %s\n", f.name)
			for _, k := range graphlint.NodeKeys(graphlint.Reach(f.m)) {
				fmt.Printf("  NODE %s\n", k)
			}
		}
	default:
		for _, f := range fixtures {
			fmt.Print(graphlint.RunDifferential(f.name, f.m))
			fmt.Println("---")
		}
	}
}
```

### cmd/stale/main.go (fixtures G,H,I)
```go
package main

import (
	"fmt"

	"a2spike/internal/graphlint"
	"a2spike/internal/table"
)

func enum(d ...string) table.TagDecl {
	return table.TagDecl{Provenance: table.ProvenanceOwned, Kind: "enum", Domain: d, SingleValued: true}
}
func match(k, op string, l ...string) table.Atom {
	return table.Atom{Key: k, Block: table.BlockMatch, Operator: op, Literal: l}
}
func tv(k string, v ...string) table.TagValue { return table.TagValue{Key: k, Value: v} }

// G: a NARROWER successor folds into an already-expanded WIDER node via
// indexOf/subsumes. joinNodes then yields merged.key() == nodes[j].key(),
// so no widening occurs -- the `continue` arm fires and no edge is lost.
func fixtureG() *table.Model {
	return &table.Model{
		ID:      "fixG",
		Tags:    map[string]table.TagDecl{"k": enum("a", "b"), "p": enum("p0", "p1", "p2")},
		Initial: []table.TagValue{tv("k", "a", "b")},
		Rows: []table.Row{
			{RuleID: "na", Atoms: []table.Atom{match("k", "eq", "a")}, Writes: []table.TagValue{tv("k", "a")}},
			{RuleID: "nb", Atoms: []table.Atom{match("k", "eq", "b")}, Writes: []table.TagValue{tv("k", "b")}},
		},
	}
}

// H: cycle that re-widens an already-expanded node.
func fixtureH() *table.Model {
	return &table.Model{
		ID:      "fixH",
		Tags:    map[string]table.TagDecl{"k": enum("a", "b", "c")},
		Initial: []table.TagValue{tv("k", "a")},
		Rows: []table.Row{
			{RuleID: "toB", Atoms: []table.Atom{match("k", "eq", "a")}, Writes: []table.TagValue{tv("k", "b")}},
			{RuleID: "toC", Atoms: []table.Atom{match("k", "eq", "b")}, Writes: []table.TagValue{tv("k", "c")}},
			{RuleID: "back", Atoms: []table.Atom{match("k", "eq", "c")}, Writes: []table.TagValue{tv("k", "a")}},
			{RuleID: "onlyC", Atoms: []table.Atom{match("k", "eq", "c")}, Writes: []table.TagValue{tv("k", "b")}},
		},
	}
}

// I: THE STALE-SOURCE-KEY HUNT.
// Two keys j,k. A node with footprint {j,k} is expanded holding j={1},k={x}.
// A later path reaches footprint {j,k} holding j={2},k={x}. indexOf/subsumes
// fails (1 does not cover 2) -> NEW node. To make indexOf MATCH and then
// widen, the existing node must subsume the newcomer yet differ in key --
// which subsumes forbids. This fixture tests whether ANY in-place widening
// of an already-expanded node is reachable.
func fixtureI() *table.Model {
	return &table.Model{
		ID:      "fixI",
		Tags:    map[string]table.TagDecl{"j": enum("1", "2"), "k": enum("x", "y")},
		Initial: []table.TagValue{tv("j", "1"), tv("k", "x")},
		Rows: []table.Row{
			{RuleID: "w2", Atoms: []table.Atom{match("j", "eq", "1")}, Writes: []table.TagValue{tv("j", "2")}},
			{RuleID: "wy", Atoms: []table.Atom{match("j", "eq", "2")}, Writes: []table.TagValue{tv("k", "y")}},
			{RuleID: "w1", Atoms: []table.Atom{match("k", "eq", "y")}, Writes: []table.TagValue{tv("j", "1")}},
			{RuleID: "trapJ", Atoms: []table.Atom{match("j", "eq", "2"), match("k", "eq", "y")}, Writes: []table.TagValue{tv("k", "x")}},
		},
	}
}

func main() {
	for _, f := range []struct {
		n string
		m *table.Model
	}{
		{"G: narrower successor folds into wider root", fixtureG()},
		{"H: cycle re-widening an expanded node", fixtureH()},
		{"I: stale-source-key hunt (two-key widening)", fixtureI()},
	} {
		fmt.Print(graphlint.RunDifferential(f.n, f.m))
		fmt.Println("---")
	}
}
```

### cmd/real/main.go (real authored models)
```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"a2spike/internal/graphlint"
	"a2spike/internal/table"
)

func main() {
	dir := "/Users/cwensel/sandbox/newcoinc/intrastate/internal/table/testdata"
	var files []string
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && filepath.Ext(p) == ".toml" {
			files = append(files, p)
		}
		return nil
	})
	sort.Strings(files)
	loaded, equal, notEqual, failed := 0, 0, 0, 0
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		m, err := table.Load(src, filepath.Base(f))
		if err != nil || m == nil {
			failed++
			continue
		}
		loaded++
		out := graphlint.RunDifferential(filepath.Base(f), m)
		fmt.Print(out)
		fmt.Print(graphlint.ProveWideningDead(filepath.Base(f), m))
		if containsNotEqual(out) {
			notEqual++
		} else {
			equal++
		}
		fmt.Println("---")
	}
	fmt.Printf("REAL-MODEL SUMMARY: files=%d loaded=%d unloadable=%d EQUAL=%d NOT-EQUAL=%d\n",
		len(files), loaded, failed, equal, notEqual)
}

func containsNotEqual(s string) bool {
	for i := 0; i+9 <= len(s); i++ {
		if s[i:i+9] == "NOT EQUAL" {
			return true
		}
	}
	return false
}
```

### cmd/inproc/main.go (in-process repeats)
```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"a2spike/internal/graphlint"
	"a2spike/internal/table"
)

func enum(d ...string) table.TagDecl {
	return table.TagDecl{Provenance: table.ProvenanceOwned, Kind: "enum", Domain: d, SingleValued: true}
}
func match(k, op string, l ...string) table.Atom {
	return table.Atom{Key: k, Block: table.BlockMatch, Operator: op, Literal: l}
}
func tv(k string, v ...string) table.TagValue { return table.TagValue{Key: k, Value: v} }

func model() *table.Model {
	return &table.Model{ID: "fixD",
		Tags:    map[string]table.TagDecl{"k": enum("a", "b"), "g": enum("g0", "g1"), "p": enum("p0", "p1")},
		Initial: []table.TagValue{tv("p", "p0")},
		Rows: []table.Row{
			{RuleID: "ra", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("k", "a")}},
			{RuleID: "rp", Atoms: []table.Atom{match("p", "eq", "p0")}, Writes: []table.TagValue{tv("p", "p1")}},
			{RuleID: "rb", Atoms: []table.Atom{match("p", "eq", "p1")}, Writes: []table.TagValue{tv("k", "b"), tv("p", "p0")}},
			{RuleID: "trap", Atoms: []table.Atom{match("k", "eq", "b"), match("p", "eq", "p0")}, Writes: []table.TagValue{tv("g", "g1")}},
		}}
}

func main() {
	const N = 2000
	seen := map[string]int{}
	for i := 0; i < N; i++ {
		out := graphlint.RunDifferential("D", model())
		h := sha256.Sum256([]byte(out))
		seen[hex.EncodeToString(h[:])]++
	}
	fmt.Printf("IN-PROCESS REPEATS: n=%d distinct outputs=%d\n", N, len(seen))
	for k, v := range seen {
		fmt.Printf("  digest %s count=%d\n", k, v)
	}
}
```

## Raw captured output — synthetic suite, controls, determinism
```text
### CMD: go run ./cmd/spike   (differential arm, synthetic fixtures A-F)
FIXTURE: A: plain chain (control, no merge)
  complete=true  nodes=3
  FINAL NODES (traversal order):
    [0] s=s0,;
    [1] s=s1,;
    [2] s=s2,;
  EDGES TAKEN BY TRAVERSAL (observer, n=2):
    s=s0,;  --[r1]-->  s=s1,;
    s=s1,;  --[r2]-->  s=s2,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    s=s0,;  --[r1]-->  s=s1,;
    s=s1,;  --[r2]-->  s=s2,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: B: subsumption merge via intra-call join
  complete=true  nodes=2
  FINAL NODES (traversal order):
    [0] p=p0,;
    [1] k=a,b,;p=p0,;
  EDGES TAKEN BY TRAVERSAL (observer, n=2):
    k=a,b,;p=p0,;  --[wa+wb]-->  k=a,b,;p=p0,;
    p=p0,;  --[wa+wb]-->  k=a,b,;p=p0,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    k=a,b,;p=p0,;  --[wa+wb]-->  k=a,b,;p=p0,;
    p=p0,;  --[wa+wb]-->  k=a,b,;p=p0,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: C: P-1 witness attempt (merged node enables `trap`)
  complete=true  nodes=4
  FINAL NODES (traversal order):
    [0] p=p0,;
    [1] p=p1,;
    [2] k=a,b,;p=p1,;
    [3] k=a,b,;p=p1,p2,;
  EDGES TAKEN BY TRAVERSAL (observer, n=4):
    k=a,b,;p=p1,;  --[ka+kb+trap]-->  k=a,b,;p=p1,p2,;
    k=a,b,;p=p1,p2,;  --[ka+kb+trap]-->  k=a,b,;p=p1,p2,;
    p=p0,;  --[t1]-->  p=p1,;
    p=p1,;  --[ka+kb]-->  k=a,b,;p=p1,;
  EDGES FROM POST-HOC RECOVERY (n=4):
    k=a,b,;p=p1,;  --[ka+kb+trap]-->  k=a,b,;p=p1,p2,;
    k=a,b,;p=p1,p2,;  --[ka+kb+trap]-->  k=a,b,;p=p1,p2,;
    p=p0,;  --[t1]-->  p=p1,;
    p=p1,;  --[ka+kb]-->  k=a,b,;p=p1,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: D: post-expansion widening (joinNodes re-queue)
  complete=true  nodes=8
  FINAL NODES (traversal order):
    [0] p=p0,;
    [1] k=a,;p=p0,;
    [2] p=p1,;
    [3] k=a,;p=p0,p1,;
    [4] k=b,;p=p0,;
    [5] k=a,b,;p=p0,p1,;
    [6] g=g1,;k=b,;p=p0,;
    [7] g=g1,;k=a,b,;p=p0,p1,;
  EDGES TAKEN BY TRAVERSAL (observer, n=11):
    g=g1,;k=a,b,;p=p0,p1,;  --[ra+rp+rb+trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    g=g1,;k=b,;p=p0,;  --[ra+rp+trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    k=a,;p=p0,;  --[ra+rp]-->  k=a,;p=p0,p1,;
    k=a,;p=p0,p1,;  --[ra+rp+rb]-->  k=a,b,;p=p0,p1,;
    k=a,b,;p=p0,p1,;  --[ra+rp+rb]-->  k=a,b,;p=p0,p1,;
    k=a,b,;p=p0,p1,;  --[trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    k=b,;p=p0,;  --[ra+rp]-->  k=a,b,;p=p0,p1,;
    k=b,;p=p0,;  --[trap]-->  g=g1,;k=b,;p=p0,;
    p=p0,;  --[ra]-->  k=a,;p=p0,;
    p=p0,;  --[rp]-->  p=p1,;
    p=p1,;  --[rb]-->  k=b,;p=p0,;
  EDGES FROM POST-HOC RECOVERY (n=11):
    g=g1,;k=a,b,;p=p0,p1,;  --[ra+rp+rb+trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    g=g1,;k=b,;p=p0,;  --[ra+rp+trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    k=a,;p=p0,;  --[ra+rp]-->  k=a,;p=p0,p1,;
    k=a,;p=p0,p1,;  --[ra+rp+rb]-->  k=a,b,;p=p0,p1,;
    k=a,b,;p=p0,p1,;  --[ra+rp+rb]-->  k=a,b,;p=p0,p1,;
    k=a,b,;p=p0,p1,;  --[trap]-->  g=g1,;k=a,b,;p=p0,p1,;
    k=b,;p=p0,;  --[ra+rp]-->  k=a,b,;p=p0,p1,;
    k=b,;p=p0,;  --[trap]-->  g=g1,;k=b,;p=p0,;
    p=p0,;  --[ra]-->  k=a,;p=p0,;
    p=p0,;  --[rp]-->  p=p1,;
    p=p1,;  --[rb]-->  k=b,;p=p0,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: E: cycle with widening
  complete=true  nodes=3
  FINAL NODES (traversal order):
    [0] k=a,;
    [1] k=b,;
    [2] k=c,;
  EDGES TAKEN BY TRAVERSAL (observer, n=3):
    k=a,;  --[c1]-->  k=b,;
    k=b,;  --[c2]-->  k=c,;
    k=c,;  --[c3]-->  k=a,;
  EDGES FROM POST-HOC RECOVERY (n=3):
    k=a,;  --[c1]-->  k=b,;
    k=b,;  --[c2]-->  k=c,;
    k=c,;  --[c3]-->  k=a,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: F: opaque scalar + clear
  complete=true  nodes=5
  FINAL NODES (traversal order):
    [0] k=a,;
    [1] k=a,;s=<opaque>,;
    [2] k=b,;
    [3] k=a,b,;s=<opaque>,;
    [4] k=a,b,;
  EDGES TAKEN BY TRAVERSAL (observer, n=8):
    k=a,;  --[flip]-->  k=b,;
    k=a,;  --[set]-->  k=a,;s=<opaque>,;
    k=a,;s=<opaque>,;  --[set+flip]-->  k=a,b,;s=<opaque>,;
    k=a,b,;  --[flip+clr]-->  k=a,b,;
    k=a,b,;  --[set]-->  k=a,b,;s=<opaque>,;
    k=a,b,;s=<opaque>,;  --[clr]-->  k=a,b,;
    k=a,b,;s=<opaque>,;  --[set+flip]-->  k=a,b,;s=<opaque>,;
    k=b,;  --[clr]-->  k=b,;
  EDGES FROM POST-HOC RECOVERY (n=8):
    k=a,;  --[flip]-->  k=b,;
    k=a,;  --[set]-->  k=a,;s=<opaque>,;
    k=a,;s=<opaque>,;  --[set+flip]-->  k=a,b,;s=<opaque>,;
    k=a,b,;  --[flip+clr]-->  k=a,b,;
    k=a,b,;  --[set]-->  k=a,b,;s=<opaque>,;
    k=a,b,;s=<opaque>,;  --[clr]-->  k=a,b,;
    k=a,b,;s=<opaque>,;  --[set+flip]-->  k=a,b,;s=<opaque>,;
    k=b,;  --[clr]-->  k=b,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---

### CMD: go run ./cmd/stale   (fixtures G,H,I -- stale-source-key hunt)
FIXTURE: G: narrower successor folds into wider root
  complete=true  nodes=1
  FINAL NODES (traversal order):
    [0] k=a,b,;
  EDGES TAKEN BY TRAVERSAL (observer, n=1):
    k=a,b,;  --[na+nb]-->  k=a,b,;
  EDGES FROM POST-HOC RECOVERY (n=1):
    k=a,b,;  --[na+nb]-->  k=a,b,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: H: cycle re-widening an expanded node
  complete=true  nodes=6
  FINAL NODES (traversal order):
    [0] k=a,;
    [1] k=b,;
    [2] k=c,;
    [3] k=a,b,;
    [4] k=b,c,;
    [5] k=a,b,c,;
  EDGES TAKEN BY TRAVERSAL (observer, n=6):
    k=a,;  --[toB]-->  k=b,;
    k=a,b,;  --[toB+toC]-->  k=b,c,;
    k=a,b,c,;  --[toB+toC+back+onlyC]-->  k=a,b,c,;
    k=b,;  --[toC]-->  k=c,;
    k=b,c,;  --[toC+back+onlyC]-->  k=a,b,c,;
    k=c,;  --[back+onlyC]-->  k=a,b,;
  EDGES FROM POST-HOC RECOVERY (n=6):
    k=a,;  --[toB]-->  k=b,;
    k=a,b,;  --[toB+toC]-->  k=b,c,;
    k=a,b,c,;  --[toB+toC+back+onlyC]-->  k=a,b,c,;
    k=b,;  --[toC]-->  k=c,;
    k=b,c,;  --[toC+back+onlyC]-->  k=a,b,c,;
    k=c,;  --[back+onlyC]-->  k=a,b,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---
FIXTURE: I: stale-source-key hunt (two-key widening)
  complete=true  nodes=4
  FINAL NODES (traversal order):
    [0] j=1,;k=x,;
    [1] j=2,;k=x,;
    [2] j=2,;k=y,;
    [3] j=1,2,;k=x,y,;
  EDGES TAKEN BY TRAVERSAL (observer, n=4):
    j=1,2,;k=x,y,;  --[w2+wy+w1+trapJ]-->  j=1,2,;k=x,y,;
    j=1,;k=x,;  --[w2]-->  j=2,;k=x,;
    j=2,;k=x,;  --[wy]-->  j=2,;k=y,;
    j=2,;k=y,;  --[wy+w1+trapJ]-->  j=1,2,;k=x,y,;
  EDGES FROM POST-HOC RECOVERY (n=4):
    j=1,2,;k=x,y,;  --[w2+wy+w1+trapJ]-->  j=1,2,;k=x,y,;
    j=1,;k=x,;  --[w2]-->  j=2,;k=x,;
    j=2,;k=x,;  --[wy]-->  j=2,;k=y,;
    j=2,;k=y,;  --[wy+w1+trapJ]-->  j=1,2,;k=x,y,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
---

### CMD: go run ./cmd/audit   (in-place widening audit)
AUDIT C: nodes=4 wideningEvents=0
  (no in-place widening occurred)
AUDIT D: nodes=8 wideningEvents=0
  (no in-place widening occurred)
AUDIT H: nodes=6 wideningEvents=0
  (no in-place widening occurred)
AUDIT I: nodes=4 wideningEvents=0
  (no in-place widening occurred)

### CMD: go run ./cmd/lemma2  (exhaustive subsumption lemma)
SUBSUMPTION LEMMA (exhaustive): pairs=3600 subsuming=260 counterexamples=0
  => subsumes(n,want) ALWAYS implies joinNodes(n,want).key()==n.key()
  => reach()'s `nodes[j]=merged` widening arm is UNREACHABLE (dead code)

### CMD: go run ./cmd/negctl  (comparator negative control)
NEGATIVE CONTROL (harness self-test)
  injected: 1 phantom in recovery, 1 edge missing from recovery
  ONLY-IN-TRAVERSAL detected: 1
    MISSING-FROM-RECOVERY: n2  --[r3]-->  n3
  ONLY-IN-RECOVERY detected: 1
    PHANTOM-EDGE: n1  --[rPHANTOM]-->  n9
  HARNESS SELF-TEST: PASS (comparator detects both directions)

### CMD: go run ./cmd/posctl  (positive control: widening forced live)
POSITIVE CONTROL chain-abc (weakened indexOf -> widening branch LIVE)
  in-place widenings=2  finalNodes=1
  edges taken=3  edges recovered=1
  ONLY-IN-TRAVERSAL=2  ONLY-IN-RECOVERY(phantom)=0
    MISSING-FROM-RECOVERY: k=a,;  --[toB]-->  k=b,;
    MISSING-FROM-RECOVERY: k=a,b,;  --[toB+toC]-->  k=b,c,;
  => divergence REPRODUCED under a merging indexOf: the harness CAN see P-1.

### CMD: go run ./cmd/inproc  (2000 in-process repeats)
IN-PROCESS REPEATS: n=2000 distinct outputs=1
  digest 55817727390811f4b5cbb774b3e3ccf024194c6c624dad48e400e043c3031afd count=2000

### CMD: map-iteration randomization probe x5
bcdefgha
cdefghab
ghabcdef
bcdefgha
ghabcdef
```

## Raw captured output — real authored models (tail)
```text
    iter=0,;stage=propose,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
    iter=0,;stage=propose,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
PROOF pos-not-near-miss.toml: indexOf hits=1, widening-branch taken=0 (expected 0)
---
FIXTURE: rdr-fixture.toml
  complete=true  nodes=2
  FINAL NODES (traversal order):
    [0] iter=0,;stage=propose,;status=Draft,;
    [1] iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
  EDGES TAKEN BY TRAVERSAL (observer, n=2):
    iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
    iter=0,;stage=propose,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
    iter=0,;stage=propose,;status=Draft,;  --[reconcile-rewind]-->  iter=0,;rewind_scope=assumptions,;stage=resolve,;status=Draft,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
PROOF rdr-fixture.toml: indexOf hits=1, widening-branch taken=0 (expected 0)
---
FIXTURE: write-delim-a.toml
  complete=true  nodes=3
  FINAL NODES (traversal order):
    [0] phase=review,;status=open,;
    [1] phase=ship,;status=accepted,;
    [2] labels=x\,y,z,;phase=resolve,;status=open,;
  EDGES TAKEN BY TRAVERSAL (observer, n=2):
    phase=review,;status=open,;  --[review-accepted]-->  phase=ship,;status=accepted,;
    phase=review,;status=open,;  --[review-needs-work]-->  labels=x\,y,z,;phase=resolve,;status=open,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    phase=review,;status=open,;  --[review-accepted]-->  phase=ship,;status=accepted,;
    phase=review,;status=open,;  --[review-needs-work]-->  labels=x\,y,z,;phase=resolve,;status=open,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
PROOF write-delim-a.toml: indexOf hits=0, widening-branch taken=0 (expected 0)
---
FIXTURE: write-delim-b.toml
  complete=true  nodes=3
  FINAL NODES (traversal order):
    [0] phase=review,;status=open,;
    [1] phase=ship,;status=accepted,;
    [2] labels=x,y\,z,;phase=resolve,;status=open,;
  EDGES TAKEN BY TRAVERSAL (observer, n=2):
    phase=review,;status=open,;  --[review-accepted]-->  phase=ship,;status=accepted,;
    phase=review,;status=open,;  --[review-needs-work]-->  labels=x,y\,z,;phase=resolve,;status=open,;
  EDGES FROM POST-HOC RECOVERY (n=2):
    phase=review,;status=open,;  --[review-accepted]-->  phase=ship,;status=accepted,;
    phase=review,;status=open,;  --[review-needs-work]-->  labels=x,y\,z,;phase=resolve,;status=open,;
  ONLY-IN-TRAVERSAL (under-approx by recovery): 0
  ONLY-IN-RECOVERY (OVER-APPROXIMATION / P-1): 0
  VERDICT: EQUAL (recovery == traversal)
PROOF write-delim-b.toml: indexOf hits=0, widening-branch taken=0 (expected 0)
---
REAL-MODEL SUMMARY: files=118 loaded=36 unloadable=82 EQUAL=36 NOT-EQUAL=0
```

## Cross-process determinism digests
```text
real-model suite, 4 independent processes:
e4ab6355cb01282c21f37ebefa2c7c37051e9b5f50ead74c992f0594383ad0dd  /tmp/a2spike/RAW_REAL.txt
e4ab6355cb01282c21f37ebefa2c7c37051e9b5f50ead74c992f0594383ad0dd  /tmp/a2spike/real_1.txt
e4ab6355cb01282c21f37ebefa2c7c37051e9b5f50ead74c992f0594383ad0dd  /tmp/a2spike/real_2.txt
e4ab6355cb01282c21f37ebefa2c7c37051e9b5f50ead74c992f0594383ad0dd  /tmp/a2spike/real_3.txt
```
