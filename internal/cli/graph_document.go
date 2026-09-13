package cli

// RDR 0021 — the exported normalized-graph document (`0021:C2`).
//
// The struct's json tags ARE the enumeration seam (`0029:C4`, REQ-36): the
// compiler is what keeps the wire spellings in step with this record, so
// nothing here restates the field list as data. `schema` is declared FIRST
// because the encoder emits struct fields in declaration order, and C2
// makes the marker both required and LEADING.
//
// Every declared collection is constructed as an empty slice or map rather
// than left nil (REQ-28/REQ-94): a nil Go slice marshals to `null`, and C2
// admits only `[]`/`{}` for an empty collection and ABSENCE for an
// inapplicable optional member. The three states stay distinguishable on
// the wire.

import (
	"slices"
	"strings"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/table"
)

// schemaMarkerValue is C2's document version marker, initial value. It
// versions the DOCUMENT; `0029:C1`'s `schema_version` versions the
// ENVELOPE and is never projected in here (REQ-37).
const schemaMarkerValue = "intrastate.graph/1"

// abstractionMarkerValue is the REQUIRED token on the `reach` block. It
// states the relation is the DECLARED over-approximation — merged nodes,
// guard and observed atoms unpruned — so a consumer cannot read it as the
// runtime relation (REQ-29).
const abstractionMarkerValue = "declared-over-approximation"

// graphDoc is the exported document. Field ORDER is the wire order; field
// NAMES are C2's normative spellings.
type graphDoc struct {
	Schema   string           `json:"schema"`
	Model    string           `json:"model"`
	Class    string           `json:"class"`
	Tags     []graphTagDoc    `json:"tags"`
	Initial  []graphTagValue  `json:"initial"`
	Terminal [][]graphAtomDoc `json:"terminal"`
	Rows     []graphRowDoc    `json:"rows"`
	Groups   []graphGroupDoc  `json:"groups"`
	Reach    graphReachDoc    `json:"reach"`
}

// graphTagDoc is one declared tag. `domain` is present exactly when
// `guard.AssignmentCount` reports the declaration's domain finite, and
// ABSENT — key omitted — otherwise (REQ-20); `omitempty` is what delivers
// the absence, and the finite arm always carries at least one member.
type graphTagDoc struct {
	Name         string   `json:"name"`
	Provenance   string   `json:"provenance"`
	Kind         string   `json:"kind"`
	Required     bool     `json:"required"`
	SingleValued bool     `json:"single_valued"`
	Domain       []string `json:"domain,omitempty"`
}

// graphTagValue is one key bound to its member sequence — a set-valued
// member travels as a JSON ARRAY, which is what closes 0002's lossy
// set-literal rendering for this document (REQ-33).
type graphTagValue struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// graphAtomDoc is one normalized predicate atom, its fields flat.
type graphAtomDoc struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Literal  []string `json:"literal"`
	Block    string   `json:"block"`
}

// graphRowDoc carries RDR 0002's dump field vocabulary — the columns
// `table.DumpColumns()` fixes, in that order, lowercase snake_case
// (REQ-22/REQ-23/REQ-108).
type graphRowDoc struct {
	Identity      string          `json:"identity"`
	Source        string          `json:"source"`
	Kind          string          `json:"kind"`
	Outcome       string          `json:"outcome"`
	Atoms         []graphAtomDoc  `json:"atoms"`
	Next          []graphTagValue `json:"next"`
	Writes        []graphTagValue `json:"writes"`
	RequiresOwned []string        `json:"requires_owned"`
	Gate          []string        `json:"gate"`
	Escape        []string        `json:"escape"`
	Emit          []graphEmitDoc  `json:"emit"`
}

// graphEmitDoc is one `[rule.emit]` pair. Its value is the RAW authored
// string: an emit key is not a tag key, so no set semantics applies.
type graphEmitDoc struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// graphGroupDoc is one scoped row group: the selection context and the
// rule ids it claims (REQ-24).
type graphGroupDoc struct {
	Context string   `json:"context"`
	Rules   []string `json:"rules"`
}

// graphReachDoc is the merged fixpoint relation plus its required
// abstraction marker (REQ-25/REQ-29).
type graphReachDoc struct {
	Abstraction string         `json:"abstraction"`
	Nodes       []graphNodeDoc `json:"nodes"`
	Edges       []graphEdgeDoc `json:"edges"`
}

// graphNodeDoc is one reachable owned-state. `values` is an OBJECT keyed
// by tag name whose every value is that tag's sorted, deduplicated value
// ARRAY — the invention-free projection of `reach.go::Node`'s
// `Values map[string][]string` (REQ-25, A7). Joined `key=value` strings
// live only inside the node key's own fingerprint and are never a value
// here.
type graphNodeDoc struct {
	ID     string              `json:"id"`
	Values map[string][]string `json:"values"`
}

// graphEdgeDoc is one reachability edge, identified by the whole triple.
type graphEdgeDoc struct {
	From string `json:"from"`
	To   string `json:"to"`
	Rule string `json:"rule"`
}

// buildGraphDocument assembles the document over a loaded model and the
// relation the traversal published.
//
// Every sequence is pre-sorted by C2's declared orders before it reaches
// the encoder, so no Go map iteration is observable on the wire (REQ-43,
// REQ-107). The model's own `Rows` arrive pre-sorted in 0002's canonical
// row order and its atoms in 0002's canonical atom order, so those two are
// carried rather than re-sorted — re-sorting them here would be a second
// ordering authority that could drift from 0002's.
func buildGraphDocument(m *table.Model, nodes []graphlint.Node, edges []graphlint.Edge) graphDoc {
	class := m.Class
	if class == "" {
		// The zero value reads everywhere as a state machine, so the
		// document carries the class rather than an empty string a consumer
		// would have to know to interpret.
		class = table.ClassStateMachine
	}

	return graphDoc{
		Schema:   schemaMarkerValue,
		Model:    m.ID,
		Class:    class,
		Tags:     graphTags(m),
		Initial:  graphTagValues(m.Initial),
		Terminal: graphTerminal(m),
		Rows:     graphRows(m),
		Groups:   graphGroups(m),
		Reach:    graphReachDocOf(nodes, edges),
	}
}

// graphTags renders the declared tags in key order.
func graphTags(m *table.Model) []graphTagDoc {
	out := make([]graphTagDoc, 0, len(m.Tags))
	for _, name := range slices.Sorted(graphTagKeys(m.Tags)) {
		decl := m.Tags[name]
		tag := graphTagDoc{
			Name:         name,
			Provenance:   string(decl.Provenance),
			Kind:         decl.Kind,
			Required:     decl.Required,
			SingleValued: decl.SingleValued,
		}
		// `domain` is present exactly when the declaration carries a finite
		// domain, which `guard.AssignmentCount` is the one authority on
		// (REQ-20). A non-finite declaration omits the key entirely rather
		// than carrying an empty array, because absent and empty mean
		// different things here.
		if _, finite := guard.AssignmentCount(decl); finite && len(decl.Domain) > 0 {
			tag.Domain = slices.Clone(decl.Domain)
		}
		out = append(out, tag)
	}
	return out
}

// graphTerminal carries the declared stop set as the DEREFERENCED
// predicate sets the model holds, never the bare context ids (REQ-21,
// REQ-32). No node is marked terminal: which merged nodes satisfy a
// predicate set is RDR 0015's quantifier, and this record declares no
// evaluator (REQ-31).
func graphTerminal(m *table.Model) [][]graphAtomDoc {
	out := make([][]graphAtomDoc, 0, len(m.Terminal))
	for _, set := range m.Terminal {
		out = append(out, graphAtoms(set))
	}
	return out
}

// graphRows renders the candidate rows in the model's own canonical order.
func graphRows(m *table.Model) []graphRowDoc {
	out := make([]graphRowDoc, 0, len(m.Rows))
	for _, r := range m.Rows {
		out = append(out, graphRowDoc{
			Identity:      r.Identity(),
			Source:        r.SourceLocator,
			Kind:          string(r.Kind()),
			Outcome:       r.Outcome,
			Atoms:         graphAtoms(r.Atoms),
			Next:          graphTagValues(r.NextTags),
			Writes:        graphTagValues(r.Writes),
			RequiresOwned: stringsOrEmpty(r.RequiresOwned),
			Gate:          stringsOrEmpty(r.Gate),
			Escape:        stringsOrEmpty(r.Escape),
			Emit:          graphEmit(r.Emit),
		})
	}
	return out
}

// graphAtoms renders a predicate set, carrying 0002's canonical atom order
// the normalized value already holds.
func graphAtoms(atoms []table.Atom) []graphAtomDoc {
	out := make([]graphAtomDoc, 0, len(atoms))
	for _, a := range atoms {
		out = append(out, graphAtomDoc{
			Key:      a.Key,
			Operator: a.Operator,
			Literal:  stringsOrEmpty(a.Literal),
			Block:    string(a.Block),
		})
	}
	return out
}

// graphTagValues renders a key/member-sequence list in key order.
func graphTagValues(values []table.TagValue) []graphTagValue {
	out := make([]graphTagValue, 0, len(values))
	for _, v := range values {
		out = append(out, graphTagValue{Key: v.Key, Value: stringsOrEmpty(v.Value)})
	}
	slices.SortFunc(out, func(a, b graphTagValue) int {
		return strings.Compare(a.Key, b.Key)
	})
	return out
}

// graphEmit renders the authored emit block, already key-sorted.
func graphEmit(emit []table.EmitValue) []graphEmitDoc {
	out := make([]graphEmitDoc, 0, len(emit))
	for _, e := range emit {
		out = append(out, graphEmitDoc{Key: e.Key, Value: e.Value})
	}
	return out
}

// graphGroups renders RDR 0003's partition — this record defines no second
// one — with each group's rule ids sorted so construction order is
// unobservable.
func graphGroups(m *table.Model) []graphGroupDoc {
	groups := guard.Groups(m)
	out := make([]graphGroupDoc, 0, len(groups))
	for _, g := range groups {
		rules := make([]string, 0, len(g.Rows))
		for _, r := range g.Rows {
			if !slices.Contains(rules, r.RuleID) {
				rules = append(rules, r.RuleID)
			}
		}
		slices.Sort(rules)
		out = append(out, graphGroupDoc{Context: g.Context.String(), Rules: rules})
	}
	return out
}

// graphReachDocOf renders the relation: nodes sorted by node key, edges by
// (from, to, rule), each node's per-tag values sorted and deduplicated
// (REQ-26, REQ-107).
func graphReachDocOf(nodes []graphlint.Node, edges []graphlint.Edge) graphReachDoc {
	outNodes := make([]graphNodeDoc, 0, len(nodes))
	for _, n := range nodes {
		values := make(map[string][]string, len(n.Values))
		for key, held := range n.Values {
			// Sorted and deduplicated per tag. The traversal already
			// canonicalizes, and repeating it here costs nothing and makes
			// the wire order a property of this document rather than an
			// inherited one.
			values[key] = slices.Compact(slices.Sorted(slices.Values(held)))
		}
		outNodes = append(outNodes, graphNodeDoc{
			ID:     graphlint.NodeKey(n),
			Values: values,
		})
	}
	slices.SortFunc(outNodes, func(a, b graphNodeDoc) int {
		return strings.Compare(a.ID, b.ID)
	})

	outEdges := make([]graphEdgeDoc, 0, len(edges))
	for _, e := range edges {
		outEdges = append(outEdges, graphEdgeDoc{From: e.From, To: e.To, Rule: e.Rule})
	}

	return graphReachDoc{
		Abstraction: abstractionMarkerValue,
		Nodes:       outNodes,
		Edges:       outEdges,
	}
}

// stringsOrEmpty returns a non-nil clone, so an empty declared collection
// renders `[]` and never `null` (REQ-28/REQ-94).
func stringsOrEmpty(in []string) []string {
	if len(in) == 0 {
		return []string{}
	}
	return slices.Clone(in)
}

// graphTagKeys yields a tag map's keys for sorting, so the declared tags
// reach the wire in key order rather than in Go map range order.
func graphTagKeys[V any](m map[string]V) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
