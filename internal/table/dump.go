package table

import (
	"fmt"
	"slices"
	"strings"
)

// dumpColumns is the closed column vocabulary. It is fixed VERBATIM as the
// normative fixtures author it — `identity` (not `row_identity`), `source`
// (not `source_locator`), `writes` (not `write`) — because those fixtures
// are what implementation tests promote (`0002:C19`).
// `emit` joins the vocabulary APPENDED LAST, after `escape` (`0010:C3`).
// That is the only insertion that leaves every existing column at its
// existing index, so an explicit `[dump]` list gains one member at the end
// and changes in no other way.
var dumpColumns = []string{
	"identity", "source", "kind", "outcome",
	"atoms", "next", "writes", "requires_owned", "gate", "escape", "emit",
}

// DumpColumns returns the closed dump column vocabulary in its canonical
// order.
func DumpColumns() []string {
	return slices.Clone(dumpColumns)
}

// Dump renders one model's expanded table.
//
// The dump is DERIVED from the normalized candidate-row value and carries
// every field of it, plus the derived row-kind column. It is emitted from
// the pre-sorted row sequence, never by iterating a map and never by
// delegating key order to an encoder, so repeated dumps of one model are
// byte-identical (`0002:C19`).
//
// Rendering never feeds back into the normalized value: the kind column is
// computed here and stored nowhere (`0002:C5`).
func Dump(m *Model) string {
	var b strings.Builder
	writeModel(&b, m)
	return b.String()
}

// DumpAll renders several models into one table. A dump spanning more than
// one model carries a model-unique model id, and every row identity is
// model-qualified, so rows from two models never collide (`0002:C19`).
func DumpAll(models []*Model) string {
	var b strings.Builder
	for _, m := range models {
		writeModel(&b, m)
	}
	return b.String()
}

func writeModel(b *strings.Builder, m *Model) {
	fmt.Fprintf(b, "MODEL %s rows=%d outcomes=%s\n",
		m.ID, len(m.Rows), strings.Join(m.Outcomes, ","))
	for _, r := range m.Rows {
		writeRow(b, m, r)
	}
}

func writeRow(b *strings.Builder, m *Model, r Row) {
	order := m.DumpOrder
	if len(order) == 0 {
		order = dumpColumns
	}
	parts := make([]string, 0, len(order))
	for _, col := range order {
		parts = append(parts, col+"="+column(m, r, col))
	}
	b.WriteString(strings.Join(parts, " "))
	b.WriteString("\n")
}

// column renders one column of one row. `[dump]` settings may reorder these
// columns; they may never omit one, which the loader enforces.
func column(m *Model, r Row, col string) string {
	switch col {
	case "identity":
		return r.Identity()
	case "source":
		return r.SourceLocator
	case "kind":
		return string(r.Kind())
	case "outcome":
		return r.Outcome
	case "atoms":
		return bracket(renderAtoms(m, r), "; ")
	case "next":
		return bracket(renderTagValues(m, r.NextTags), "; ")
	case "writes":
		return bracket(renderTagValues(m, r.Writes), "; ")
	case "requires_owned":
		return bracket(r.RequiresOwned, ",")
	case "gate":
		return bracket(r.Gate, ",")
	case "escape":
		return bracket(r.Escape, ",")
	case "emit":
		return bracket(renderEmit(r.Emit), "; ")
	default:
		return ""
	}
}

func bracket(parts []string, sep string) string {
	return "[" + strings.Join(parts, sep) + "]"
}

func renderAtoms(m *Model, r Row) []string {
	out := make([]string, 0, len(r.Atoms))
	for _, a := range r.Atoms {
		out = append(out, fmt.Sprintf("%s.%s=%s@%s",
			a.Key, a.Operator, renderValue(m, a.Key, a.Literal), a.Block))
	}
	return out
}

// renderEmit spells the emit block as `key=value` pairs in key order — the
// sequence is already key-sorted at normalization — bracketed like `writes`
// (`0010:C3`).
//
// The value is the RAW authored string and is deliberately NOT routed
// through renderValue, whose quoting and bracketing key on a tag's declared
// kind and member count. An emit key has neither: it is undeclared, so
// there is no kind to consult, and its value is one string rather than a
// member sequence. This is a READER — it never writes back into the row.
func renderEmit(emit []EmitValue) []string {
	out := make([]string, 0, len(emit))
	for _, e := range emit {
		out = append(out, e.Key+"="+e.Value)
	}
	return out
}

func renderTagValues(m *Model, tags []TagValue) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.Key+"="+renderValue(m, t.Key, t.Value))
	}
	return out
}

// renderValue spells one value for DISPLAY. It is never an identity: no
// merge key, dedup, or sort is derived from it (`0002:C10`).
//
// A SET-valued field is spelled visibly as a set with each member QUOTED,
// because an unquoted separator is forgeable by a member that contains it:
// ["x|y","z"] and ["x","y|z"] would render alike, and a reviewer comparing
// two rows could not see that a difference may be hiding (`0002:C10`,
// `0002:RT`). Set-ness is the declared kind, so a one-member set still
// renders bracketed.
func renderValue(m *Model, key string, members []string) string {
	if isClear(members) {
		return ClearSentinel
	}
	if m.Tags[key].Kind == "set" || len(members) != 1 {
		parts := make([]string, 0, len(members))
		for _, member := range members {
			parts = append(parts, fmt.Sprintf("%q", member))
		}
		return "[" + strings.Join(parts, " ") + "]"
	}
	return members[0]
}
