package graphlint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/table"
)

// elementModel, elementInitial, and elementTraversal are the graph element
// ids the findings with no source rule to name carry instead. REQ-127
// requires such a finding stay actionable and stable, and these are the
// authored (or structural) sites it names.
const (
	elementModel     = "model"
	elementInitial   = "initial"
	elementTerminal  = "terminal"
	elementTraversal = "traversal"
)

// analysis carries one run's derived graph view and the findings collected
// over it. Deriving once — the traversal, the group partition, the
// reachable-context filter — is what keeps invariants 3 and 4 running once
// per reachable GROUP rather than once per reachable node.
type analysis struct {
	model *table.Model
	// nodes is the reachable owned-state set, a fixpoint over merged nodes.
	nodes []Node
	// complete records whether the traversal finished under the published
	// node ceiling.
	complete bool
	// groups is RDR 0003's partition; this RDR defines no second one.
	groups []guard.Group
	// reachable indexes the groups whose selection context some reachable
	// owned-state satisfies. Reachability is a FILTER over contexts and
	// never changes which rows are in a group.
	reachable map[string]bool

	findings []clierr.Finding
}

func newAnalysis(m *table.Model) *analysis {
	nodes, complete := reach(m)
	a := &analysis{
		model:     m,
		nodes:     nodes,
		complete:  complete,
		groups:    guard.Groups(m),
		reachable: map[string]bool{},
	}
	for _, g := range a.groups {
		if a.contextReachable(g) {
			a.reachable[g.Context.String()] = true
		}
	}
	return a
}

// contextReachable reports whether some reachable owned-state satisfies the
// group's match pattern. Only the OWNED half of the pattern is consulted:
// observed and recognized keys arrive at runtime and name no node
// dimension.
func (a *analysis) contextReachable(g guard.Group) bool {
	for _, n := range a.nodes {
		if a.nodeSatisfiesMatch(n, g.Context.Match) {
			return true
		}
	}
	return false
}

// nodeSatisfiesMatch reports whether the node satisfies every owned atom in
// the match pattern. It is an EXISTENTIAL test per atom — some held value
// satisfies it — which merged nodes may be read directly for.
func (a *analysis) nodeSatisfiesMatch(n Node, match []table.Atom) bool {
	for _, atom := range match {
		if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			continue
		}
		if !ownedAtomSatisfiable(n, atom) {
			return false
		}
	}
	return true
}

// emit appends one finding, filling in the model identity and the severity
// its code declares. No caller sets either, so neither can drift from the
// taxonomy.
func (a *analysis) emit(f clierr.Finding) {
	f.Model = a.model.ID
	f.Severity = severityFor(f.Code)
	a.findings = append(a.findings, f)
}

// --- invariant 1: dangling edge ------------------------------------------

// checkDanglingEdge checks that every reference a row names resolves to a
// declared model element, and that the model declares an initial owned
// state.
//
// RDR 0002's loader refuses every SYNTACTIC dangling reference before
// normalization — an undeclared tag, an out-of-alphabet outcome, an unknown
// gate accessor, an unknown terminal context — so the arms that reach lint
// are the two the loader deliberately routes here: the missing `[initial]`
// declaration, and a terminal context whose predicate reads a NON-OWNED tag
// (JDR 0001 §D7(i), which RDR 0002 fences as this RDR's blocking finding
// rather than a load failure).
func (a *analysis) checkDanglingEdge() {
	if len(a.model.Initial) == 0 {
		// Never "nothing reachable, therefore clean": the absent root is
		// itself the defect, reported against the declaration that is
		// missing rather than against an arbitrary row.
		a.emit(clierr.Finding{
			Code:    CodeDanglingEdge,
			Element: elementModel,
			Message: "the model declares no initial owned state; add an " +
				"`[initial]` table assigning every always-present owned tag",
		})
	}

	// A terminal is a predicate over OWNED tags. A terminal context whose
	// predicate reads an observed or recognized tag names an element no
	// owned-state can resolve, so the edge's destination dangles — the same
	// unresolvable-reference defect invariant 1 already covers, checked as
	// tag keys rather than as state names.
	for i, set := range a.model.Terminal {
		atoms := slices.Clone(set)
		slices.SortFunc(atoms, compareAtoms)
		for _, atom := range atoms {
			decl, declared := a.model.Tags[atom.Key]
			if declared && decl.Provenance == table.ProvenanceOwned {
				continue
			}
			a.emit(clierr.Finding{
				Code:     CodeDanglingEdge,
				Element:  fmt.Sprintf("%s[%d]", elementTerminal, i),
				Key:      atom.Key,
				Operator: atom.Operator,
				Literal:  strings.Join(atom.Literal, ","),
				Block:    string(atom.Block),
				Message: fmt.Sprintf("terminal predicate reads the %s tag %q; "+
					"a terminal is a predicate over owned tags, so a non-owned "+
					"key resolves to no owned-state element",
					provenanceName(decl, declared), atom.Key),
			})
		}
	}
}

// provenanceName spells a declaration's provenance for a diagnostic,
// including the undeclared case.
func provenanceName(decl table.TagDecl, declared bool) string {
	if !declared {
		return "undeclared"
	}
	return string(decl.Provenance)
}

// --- the published node ceiling ------------------------------------------

// checkNodeCeiling reports a traversal whose reachable node set exceeds the
// published ceiling. The finding names the TRAVERSAL, not a group, so it
// carries no rule id — the run is bounded rather than unbounded, and never
// a partial green.
func (a *analysis) checkNodeCeiling() {
	if a.complete {
		return
	}
	a.emit(clierr.Finding{
		Code:    CodeProductTooLarge,
		Element: elementTraversal,
		Message: fmt.Sprintf("the reachable owned-state set exceeds the "+
			"published node ceiling of %d; narrow a declared domain rather "+
			"than running the traversal unbounded", nodeCeiling),
	})
}

// --- invariant 5: single-valued state ------------------------------------

// checkSingleValuedState checks that no row's write block assigns a tag
// declared single-valued two values.
//
// The decision is PER ROW and SYNTACTIC. It is deliberately not decided
// over the reachability nodes: a merged node's value set has cardinality
// greater than one whenever two paths write different single values, which
// is the legal shape of a two-path merge rather than a violation, and
// reporting it would make the abstraction's join rule a false-positive
// generator.
func (a *analysis) checkSingleValuedState() {
	for _, row := range a.model.Rows {
		for _, w := range canonicalTags(writesOf(row)) {
			decl, ok := a.model.Tags[w.Key]
			if !ok || !decl.SingleValued || isClearValue(w.Value) {
				continue
			}
			if len(canonicalValues(w.Value)) <= 1 {
				continue
			}
			a.emit(clierr.Finding{
				Code:        CodeSingleValuedState,
				Rule:        row.RuleID,
				Span:        row.SourceLocator,
				Key:         w.Key,
				Literal:     strings.Join(w.Value, ","),
				Fingerprint: Fingerprint(row),
				Message: fmt.Sprintf("row %q writes %d values to the "+
					"single-valued tag %q; a single-valued tag holds at most "+
					"one declared value", row.RuleID, len(w.Value), w.Key),
			})
		}
	}
}

// --- the always-present owned sibling code -------------------------------

// checkAlwaysPresentOwned checks that every OWNED key declared
// always-present is held in every reachable owned-state node.
//
// The check does not extend to observed or recognized keys: those arrive at
// runtime and the resolver's view assembly reads no declaration, so this
// RDR discharges only the owned half of the conformance premise.
//
// The root is INCLUDED, and that is a requirement on the `initial`
// declaration rather than an exemption: a key the `initial` table omits is
// absent at the root, and the finding is reported against that declaration
// so the diagnostic names the authored site rather than an arbitrary
// downstream node.
func (a *analysis) checkAlwaysPresentOwned() {
	if len(a.model.Initial) == 0 {
		// The absent root is already the dangling-edge finding. Reporting
		// every always-present key against a root that does not exist would
		// bury that diagnostic under derived noise.
		return
	}

	for _, key := range slices.Sorted(slices.Values(ownedRequiredKeys(a.model))) {
		rootHolds := slices.ContainsFunc(a.model.Initial,
			func(tv table.TagValue) bool { return tv.Key == key })
		if !rootHolds {
			a.emit(clierr.Finding{
				Code:    CodeAlwaysPresentOwned,
				Element: elementInitial,
				Key:     key,
				Message: fmt.Sprintf("the owned tag %q is declared "+
					"always-present but `[initial]` does not establish it, so "+
					"it is absent at the root", key),
			})
			continue
		}
		for _, n := range a.nodes {
			if len(n.Values[key]) > 0 {
				continue
			}
			a.emit(clierr.Finding{
				Code:    CodeAlwaysPresentOwned,
				Element: nodeElement(n),
				Key:     key,
				Message: fmt.Sprintf("the owned tag %q is declared "+
					"always-present but the reachable owned-state %s does not "+
					"hold it", key, nodeElement(n)),
			})
			// One finding per key is enough to name the defect; the cure is
			// the write or the clear that drops it, not a per-node census.
			break
		}
	}
}

// ownedRequiredKeys names every owned tag carrying the explicit
// always-present marker, sorted.
func ownedRequiredKeys(m *table.Model) []string {
	var out []string
	for key, decl := range m.Tags {
		if decl.Provenance == table.ProvenanceOwned && decl.Required {
			out = append(out, key)
		}
	}
	slices.Sort(out)
	return out
}

// nodeElement renders a node as a graph element id, so a node-scoped
// finding stays identifiable without a rule to name.
func nodeElement(n Node) string {
	if k := n.key(); k != "" {
		return "node{" + strings.TrimSuffix(k, ";") + "}"
	}
	return "node{}"
}

// --- invariant 2: dead end -----------------------------------------------

// checkDeadEnd checks that every reachable owned-state node satisfying no
// declared terminal is the source of at least one modeled non-escape row.
// Escape self-loops do not count as progress.
//
// A terminal is a predicate over owned tags, and a node SATISFIES it when
// every value in each of the node's per-tag value sets meets it. Partial
// satisfaction is deliberately not enough — accepting it would let a merged
// node close on a path that has not actually terminated.
//
// The test runs on SPLIT nodes, not merged ones. Universal satisfaction is
// anti-monotone in merging, so testing it against a merged node would
// accuse every flow with two routes to completion. The split is exact: each
// split node is a set of concrete views the merged node already stood for.
// It is bounded by the declared domains of the terminal-participating keys
// ONLY, never the whole lattice.
func (a *analysis) checkDeadEnd() {
	if len(a.model.Terminal) == 0 || len(a.nodes) == 0 {
		// A model declaring no terminal at all is invariant 7's defect —
		// reliance on an inferred terminal — not a dead end at every node.
		return
	}

	keys := a.terminalKeys()
	seen := map[string]bool{}
	for _, n := range a.nodes {
		for _, split := range splitNode(n, keys) {
			id := split.key()
			if seen[id] {
				continue
			}
			seen[id] = true
			if a.satisfiesSomeTerminal(split) {
				continue
			}
			if a.hasOutgoingOrdinaryRow(split) {
				continue
			}
			a.emit(clierr.Finding{
				Code:    CodeDeadEnd,
				Element: nodeElement(split),
				Message: fmt.Sprintf("the reachable owned-state %s satisfies "+
					"no declared terminal and is the source of no modeled "+
					"non-escape row; an escape self-loop is not progress",
					nodeElement(split)),
			})
		}
	}
}

// terminalKeys names the terminal-participating owned keys, sorted. They
// are the only keys the split ranges over.
func (a *analysis) terminalKeys() []string {
	var out []string
	for _, set := range a.model.Terminal {
		for _, atom := range set {
			if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
				continue
			}
			if !slices.Contains(out, atom.Key) {
				out = append(out, atom.Key)
			}
		}
	}
	slices.Sort(out)
	return out
}

// splitNode splits one merged node on the terminal-participating keys: each
// such key's value set is expanded into one node per held value, and every
// other key passes through unchanged. Splitting recovers exactness for the
// universal terminal test without touching the rest of the lattice.
func splitNode(n Node, keys []string) []Node {
	out := []Node{n.clone()}
	for _, key := range keys {
		values := n.Values[key]
		if len(values) <= 1 {
			continue
		}
		next := make([]Node, 0, len(out)*len(values))
		for _, base := range out {
			for _, v := range values {
				split := base.clone()
				split.Values[key] = []string{v}
				next = append(next, split)
			}
		}
		out = next
	}
	return out
}

// satisfiesSomeTerminal reports whether the node satisfies some declared
// terminal predicate. Satisfaction is UNIVERSAL over the node's per-tag
// value sets: every value in each set must meet the predicate.
func (a *analysis) satisfiesSomeTerminal(n Node) bool {
	for _, set := range a.model.Terminal {
		if a.nodeMeetsAll(n, set) {
			return true
		}
	}
	return false
}

// nodeMeetsAll reports whether every value the node holds for each of the
// predicate's owned keys meets that predicate.
func (a *analysis) nodeMeetsAll(n Node, set []table.Atom) bool {
	for _, atom := range set {
		if a.model.Tags[atom.Key].Provenance != table.ProvenanceOwned {
			// A non-owned terminal key is invariant 1's dangling finding.
			// It cannot be met by an owned-state, so the predicate as a
			// whole is not satisfied here.
			return false
		}
		held := n.Values[atom.Key]
		if len(held) == 0 {
			return false
		}
		for _, v := range held {
			if !atomAdmitsValue(atom, v) {
				return false
			}
		}
	}
	return len(set) > 0
}

// hasOutgoingOrdinaryRow reports whether the node is the source of some
// modeled non-escape row. An escape row's successor equals its source, so
// counting it would call a self-loop progress.
//
// The node reaching here is ALREADY SPLIT on the terminal-participating
// keys (REQ-36), so on every key that decides terminal satisfaction it
// holds a single value and the test is exact there. On the remaining keys
// it stays a merged node read existentially, and that is deliberate:
// REQ-37 bounds the split "by the declared domains of terminal-participating
// keys only, never the whole lattice", so widening the quantifier to those
// keys is not available to this invariant.
//
// The residual imprecision is the accepted false-NEGATIVE the record books
// against invariant 2 — a node multi-valued on a key participating in no
// terminal can be rescued by an exit that serves only one of its values.
// Closing it requires ranging over key combinations the merged node never
// correlated, which manufactures concrete views no path produces and turns
// the bound REQ-37 sets into false POSITIVES on a conforming model. See D12.
func (a *analysis) hasOutgoingOrdinaryRow(n Node) bool {
	for _, row := range a.model.Rows {
		if row.Kind() == table.KindEscape {
			continue
		}
		if matchSatisfiable(a.model, n, row) {
			return true
		}
	}
	return false
}

// --- invariant 7: declared terminal / escape handling --------------------

// checkTerminalEscape reports a model that RELIES on lint inferring a
// terminal or an escape row from a missing declaration.
//
// Terminal states and escape rows are explicit model data, and lint must
// not infer them from missing rows. The input-observable defect is a model
// whose flow ends at a reachable non-terminal node with no outgoing
// non-escape row and no terminal declaration covering it — an implied
// terminal. The prohibition on lint's own reasoning is what makes such a
// model a defect rather than a silently-accepted shape.
func (a *analysis) checkTerminalEscape() {
	if len(a.model.Terminal) > 0 || len(a.nodes) == 0 {
		return
	}
	for _, n := range a.nodes {
		if a.hasOutgoingOrdinaryRow(n) {
			continue
		}
		a.emit(clierr.Finding{
			Code:    CodeTerminalEscape,
			Element: nodeElement(n),
			Message: fmt.Sprintf("the reachable owned-state %s has no "+
				"outgoing non-escape row and the model declares no terminal "+
				"covering it, so the model relies on an inferred terminal; "+
				"declare it in the root `terminal` list", nodeElement(n)),
		})
		// The reliance is one defect against the missing declaration, not
		// one per node it happens to end at.
		return
	}
}

// --- the unreachable-rule advisory ---------------------------------------

// checkUnreachableRules reports a row no reachable owned-state node
// satisfies, including RDR 0002's dead-rule case — a predicate set
// requiring `recognized` to be absent, which no view reaching a row can
// satisfy. It is advisory: it never changes the success disposition.
func (a *analysis) checkUnreachableRules() {
	if len(a.model.Initial) == 0 {
		// With no root there is no reachable set to filter against, and
		// reporting every row unreachable would bury the missing-root
		// diagnostic under derived noise.
		return
	}
	for _, g := range a.groups {
		if a.reachable[g.Context.String()] {
			continue
		}
		rows := slices.Clone(g.Rows)
		slices.SortFunc(rows, func(x, y table.Row) int { return strings.Compare(x.RuleID, y.RuleID) })
		for _, row := range rows {
			a.emit(clierr.Finding{
				Code:        CodeUnreachableRule,
				Rule:        row.RuleID,
				Span:        row.SourceLocator,
				Fingerprint: Fingerprint(row),
				Message: fmt.Sprintf("no reachable owned-state satisfies the "+
					"selection context of row %q, so the row can never be a "+
					"candidate", row.RuleID),
			})
		}
	}
}
