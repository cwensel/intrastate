package cli

// RDR 0021 — the DOT arm (`0021:C2`, A6).
//
// A PURE renderer over the export VALUE: it reads the assembled document
// and consults no `graphlint` or analysis internal, so the two arms cannot
// disagree about the relation. Its node/edge SET and the abstraction
// marker are normative; its styling and attributes are not (REQ-40), so a
// later restyle is free while the identifiers stay pinned.

import (
	"strings"
)

// renderDOT renders the document as a DOT digraph.
//
// One node per reachability node, one edge per reachability edge labelled
// with its rule id, the initial node marked, and the abstraction marker in
// the header so the DIAGRAM carries it too — a reviewer reading only the
// rendered picture must still learn this is the declared
// over-approximation rather than the runtime relation (REQ-39).
//
// The sequences are already in C2's orders, so the rendering inherits the
// document's determinism rather than sorting a second time.
func renderDOT(doc graphDoc) string {
	var b strings.Builder

	// The marker rides a header comment AND the graph label: the comment
	// survives a text diff, the label survives rendering to an image.
	b.WriteString("// " + abstractionMarkerValue + "\n")
	b.WriteString("digraph " + dotQuote(doc.Model) + " {\n")
	b.WriteString("  label=" + dotQuote(doc.Model+" ("+abstractionMarkerValue+")") + ";\n")

	// The initial node is the one the traversal rooted at — the document's
	// first node in key order is not it, so the root is identified by the
	// declared initial assignment rather than by position. A model whose
	// root is the empty owned-state (a decision table) still has one node
	// standing for it.
	initial := initialNodeID(doc)

	for _, n := range doc.Reach.Nodes {
		b.WriteString("  " + dotQuote(n.ID))
		if n.ID == initial {
			// The MECHANISM is styling and therefore free; that the initial
			// node is distinguished is not (REQ-39).
			b.WriteString(" [shape=doublecircle, xlabel=\"initial\"]")
		}
		b.WriteString(";\n")
	}

	for _, e := range doc.Reach.Edges {
		b.WriteString("  " + dotQuote(e.From) + " -> " + dotQuote(e.To) +
			" [label=" + dotQuote(e.Rule) + "];\n")
	}

	b.WriteString("}")
	return b.String()
}

// initialNodeID names the node standing for the declared root: the node
// whose per-tag values are exactly the declared initial assignment over
// owned tags. It returns "" when no node matches, in which case no node is
// marked rather than an arbitrary one being singled out.
func initialNodeID(doc graphDoc) string {
	for _, n := range doc.Reach.Nodes {
		if nodeMatchesInitial(n, doc.Initial) {
			return n.ID
		}
	}
	return ""
}

// nodeMatchesInitial reports whether a node holds exactly the declared
// initial assignment. A declared key the node does not hold, or a held
// value the declaration did not name, both disqualify it — the root node
// is the one the traversal seeded, not merely one reachable from it.
func nodeMatchesInitial(n graphNodeDoc, initial []graphTagValue) bool {
	declared := 0
	for _, tv := range initial {
		held, ok := n.Values[tv.Key]
		if !ok {
			// The key names no node dimension (a non-owned provenance), so
			// it constrains nothing here.
			continue
		}
		declared++
		if !sameMembers(held, tv.Value) {
			return false
		}
	}
	// A node holding dimensions beyond the declared ones is a successor,
	// not the root.
	return declared == len(n.Values)
}

// sameMembers compares two member sequences for equality. The traversal
// abstracts a tag with no finite declared domain to its opaque value, so
// the held sequence may legitimately differ from the authored one; that
// case is admitted rather than treated as a mismatch.
func sameMembers(held, authored []string) bool {
	if len(held) == 1 && held[0] == opaqueHeldValue {
		return true
	}
	if len(held) != len(authored) {
		return false
	}
	for i := range held {
		if held[i] != authored[i] {
			return false
		}
	}
	return true
}

// opaqueHeldValue mirrors `reach.go::OpaqueValue`. It is compared, never
// minted: this record attaches no meaning to the spelling and passes it
// through verbatim (REQ-27).
const opaqueHeldValue = "<opaque>"

// dotQuote renders one identifier as a quoted DOT string.
//
// The escaping ORDER is `\` first, then `"`, then the line break (A6).
// Escaping the backslash last would re-escape the backslashes the earlier
// steps introduced; doing it first makes the encoding injective, so
// inverting it — undo `\n`, then `\"`, then `\\` — returns the source
// identifier exactly (REQ-80). A well-formed but mis-escaped identifier
// still renders at exit 0, which is why the round trip rather than the
// render is the oracle.
func dotQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return `"` + s + `"`
}
