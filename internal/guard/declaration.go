package guard

import (
	"encoding/json"
	"errors"
	"maps"
	"slices"

	"github.com/newcoinc/intrastate/internal/table"
)

// Declaration is one tag's declaration as this RDR reads it: the five
// fields the declaration model carries, plus the provenance RDR 0002
// carries alongside them.
//
// Optionality is spelled as the POSITIVE `Optional` rather than as RDR
// 0002's `required` wire key, because the conservative default this RDR
// fixes is the optional one: a declaration carrying no optionality marker
// declares the key optional, and a field whose zero value is the default
// cannot misread an omitted marker.
type Declaration struct {
	Provenance   table.Provenance
	Kind         string
	Domain       []string
	Elements     []string
	SingleValued bool
	Optional     bool
}

// declarationFields names the five fields the declaration model carries.
// Any enumeration of the model states all five.
var declarationFields = []string{
	"value kind",
	"finite domain",
	"optionality",
	"single-valuedness",
	"element universe",
}

// DeclarationFields names the five fields the declaration model carries.
func DeclarationFields() []string { return slices.Clone(declarationFields) }

// DeclarationOf reads one tag's declaration out of a loaded model. RDR 0002
// owns where a declaration is authored and how it is carried; this reads
// what it carries and never re-authors it.
func DeclarationOf(m *table.Model, key string) Declaration {
	return declarationOf(m.Tags[key])
}

func declarationOf(d table.TagDecl) Declaration {
	return Declaration{
		Provenance:   d.Provenance,
		Kind:         d.Kind,
		Domain:       slices.Clone(d.Domain),
		Elements:     slices.Clone(d.Elements),
		SingleValued: d.SingleValued,
		Optional:     !d.Required,
	}
}

// AssignmentCount returns a dimension's assignment count under the per-kind
// table, and whether the declaration carries a finite domain.
//
// The table reads the declaration's DEFAULTS, not an author's intent:
//
//   - `enum` and `int` single-valued        → |domain|
//   - `bool` single-valued                  → 2
//   - `enum`, `bool`, `int` unmarked        → 2^|domain|, one independent
//     boolean dimension per value
//   - `set`                                 → 2^|element universe|, never
//     |universe| — a set-valued tag holds any SUBSET
//   - any OPTIONAL key                      → ×2 for its {present, absent}
//     presence dimension
//
// A declaration whose domain disagrees with its kind carries no readable
// domain: agreement is what makes the domain readable, and a disagreeing
// declaration is a declaration error the loader rejects before a consumer
// sees it, so it is not a second readable shape here either.
func AssignmentCount(d table.TagDecl) (int, bool) {
	if !agrees(d) {
		return 0, false
	}

	values, ok := domainSize(d)
	if !ok {
		return 0, false
	}

	count := values
	if !d.Required {
		count *= 2
	}
	return count, true
}

// domainSize returns the number of assignments the key's VALUE dimension
// carries, before the presence factor.
func domainSize(d table.TagDecl) (int, bool) {
	switch d.Kind {
	case "enum":
		if len(d.Domain) == 0 {
			return 0, false
		}
		return spread(len(d.Domain), d.SingleValued), true
	case "bool":
		// A bool is finite by construction: its domain is the two literals,
		// declared nowhere and never absent.
		return spread(2, d.SingleValued), true
	case "int":
		if d.Min == nil || d.Max == nil {
			return 0, false
		}
		return spread(*d.Max-*d.Min+1, d.SingleValued), true
	case "set":
		if len(d.Elements) == 0 {
			return 0, false
		}
		// Values co-occur by construction, so the marker is never carried
		// and the space is always the powerset.
		return 1 << len(d.Elements), true
	default:
		// `scalar` is the opaque scalar: it carries no finite domain and can
		// never bear an exhaustiveness claim.
		return 0, false
	}
}

// spread applies the single-valued marker: exactly one value holds under
// the marker, so the dimension is |domain|; absent the marker the values
// could co-occur, so it is one independent boolean dimension per value —
// 2^|domain|, minus nothing, since the model does not assume at least one
// holds.
func spread(domain int, singleValued bool) int {
	if singleValued {
		return domain
	}
	return 1 << domain
}

// agrees reports whether a declared domain agrees with its value kind, per
// the kind/field agreement table. The loader rejects a disagreement before
// normalization completes; this is the same rule read on this RDR's own
// surface, so a declaration built in memory is judged by it too.
func agrees(d table.TagDecl) bool {
	hasDomain := len(d.Domain) > 0
	hasBound := d.Min != nil || d.Max != nil
	hasElements := len(d.Elements) > 0

	switch d.Kind {
	case "enum":
		return !hasBound && !hasElements
	case "bool":
		return !hasDomain && !hasBound && !hasElements
	case "int":
		return !hasDomain && !hasElements
	case "set":
		return !hasDomain && !hasBound && !d.SingleValued
	case "scalar":
		return !hasDomain && !hasBound && !hasElements && !d.SingleValued
	default:
		return false
	}
}

// IntDomain enumerates a bounded int declaration's inclusive values. Both
// endpoints are inclusive, so `{0..3}` enumerates four values.
func IntDomain(d table.TagDecl) []int {
	if d.Kind != "int" || d.Min == nil || d.Max == nil {
		return nil
	}
	out := make([]int, 0, *d.Max-*d.Min+1)
	for n := *d.Min; n <= *d.Max; n++ {
		out = append(out, n)
	}
	return out
}

// RequiresPredecessorWrite reports whether a provenance carries the
// reachable-predecessor-write obligation. Recognized tags are fresh event
// inputs and observed tags are re-read before matching; only an owned tag
// must have a reachable predecessor write before a row may match it.
func RequiresPredecessorWrite(p table.Provenance) bool {
	return p == table.ProvenanceOwned
}

// --- conformance ---------------------------------------------------------

// View is one evaluation view: tag key to rendered held value.
type View map[string]string

// Conforms reports whether v satisfies the declared model's two conjuncts:
// every always-present key is present, and every single-valued tag holds at
// most one of its declared domain values.
//
// Conformance is the premise every lint claim in this RDR is conditional
// on. A green exhaustiveness result asserts coverage over conforming views
// ONLY, and this RDR promises nothing about a view that violates the
// declarations.
func Conforms(m *table.Model, v View) error {
	for _, key := range slices.Sorted(maps.Keys(m.Tags)) {
		d := m.Tags[key]
		held, present := v[key]

		if d.Required && !present {
			return errors.New("always-present key " + key + " is absent from the view")
		}
		if !present {
			continue
		}
		if d.SingleValued && heldValueCount(held) > 1 {
			return errors.New("single-valued tag " + key + " holds more than one value")
		}
	}
	return nil
}

// heldValueCount reports how many values a rendered held value carries. A
// set crosses as its canonical JSON array (JDR 0001 §D13); anything else is
// one value, whatever bytes it holds.
func heldValueCount(held string) int {
	var members []string
	if err := json.Unmarshal([]byte(held), &members); err != nil {
		return 1
	}
	return len(members)
}
