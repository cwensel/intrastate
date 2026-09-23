package guard

import (
	"encoding/json"
	"errors"
	"maps"
	"math"
	"slices"

	"github.com/cwensel/intrastate/internal/table"
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

// DeclaredKinds is the one producer of the mapping NewEvaluator is built
// over: tag key → declared kind token, total over m.Tags.
//
// A key whose Kind is the empty string is OMITTED rather than mapped to "":
// an empty-string entry would read at the seam exactly like an absent key,
// collapsing two diagnostics into one. A nil model yields an empty, non-nil
// map.
func DeclaredKinds(m *table.Model) map[string]string {
	out := map[string]string{}
	if m == nil {
		return out
	}
	for key, decl := range m.Tags {
		if decl.Kind == "" {
			continue
		}
		out[key] = decl.Kind
	}
	return out
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
		// Saturate rather than double blindly. An optional key's presence
		// factor is the LAST arithmetic between a declared domain and the
		// bound comparison, and it is as able to wrap as the width and the
		// shift before it: a single-valued `{MinInt..-2}` has a width of
		// MaxInt, which doubles to -2 — read as under the bound, fully
		// provable, and GREEN over a dimension spanning the int range.
		// That is the same defect D12 named for the exponent and this
		// record's `intWidth` closed for the subtraction, so it takes the
		// same answer: report the ceiling, keep the domain finite-but-huge,
		// and let `Cardinality` saturate into `graph-product-too-large`.
		if count > cardinalityCeiling/2 {
			return cardinalityCeiling, true
		}
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
		width, ok := intWidth(*d.Min, *d.Max)
		if !ok {
			// The width itself overflowed, so the domain is FINITE but
			// larger than this implementation counts to. Reporting the
			// ceiling makes `Cardinality` saturate into the too-large
			// refusal, the same answer D12 gives an over-large exponent;
			// reporting no domain at all would instead misreport an
			// author's fully-declared bound as carrying no finite domain.
			return cardinalityCeiling, true
		}
		return spread(width, d.SingleValued), true
	case "set":
		if len(d.Elements) == 0 {
			return 0, false
		}
		// Values co-occur by construction, so the marker is never carried
		// and the space is always the powerset.
		return spread(len(d.Elements), false), true
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
//
// The exponential saturates at the same ceiling `Cardinality` saturates at,
// for the same reason it does: `1 << 64` is 0 in Go and `1 << 63` is
// negative, and either would read as UNDER the published bound and certify
// exactly the product the too-large clause refuses. Saturating keeps the
// bound comparison a comparison. A negative domain likewise floors at zero
// rather than shifting by a negative count, which panics.
func spread(domain int, singleValued bool) int {
	domain = max(domain, 0)
	if singleValued {
		return domain
	}
	if domain >= ceilingExponent {
		return cardinalityCeiling
	}
	return 1 << domain
}

// ceilingExponent is the exponent at which the spread saturates:
// `cardinalityCeiling` is `1 << ceilingExponent`.
const ceilingExponent = 40

// intWidth returns how many values the inclusive bound `{minV..maxV}`
// names, and whether that count fits in an int.
//
// `maxV - minV + 1` is evaluated in WRAPPING int arithmetic, so a domain
// as wide as the int range — `{MinInt..MaxInt}`, `{0..MaxInt}`,
// `{MinInt+1..MaxInt}` — yields zero or a negative before any ceiling
// applies. D12 saturated the SHIFT and never covered this SUBTRACTION, so
// a wrapped width read as a tiny domain: under the published bound, fully
// provable, and certified GREEN over an unbounded dimension — and then
// enumerated one value at a time, which does not terminate. Refusing the
// width here is what keeps REQ-86's bound comparison a comparison.
func intWidth(minV, maxV int) (int, bool) {
	if minV > maxV {
		return 0, false
	}
	span := maxV - minV
	if span < 0 || span == math.MaxInt {
		// The subtraction wrapped, or the inclusive `+1` would.
		return 0, false
	}
	return span + 1, true
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
		// An INVERTED bound is not a readable finite domain: REQ-18 fixes
		// `{min..max}` as inclusive at both endpoints, which presumes an
		// ordered pair — `{3..0}` names no value at all. The loader already
		// refuses it; this is the same rule read on this RDR's own surface,
		// so a declaration built in memory is judged by it too rather than
		// yielding a domain size the table defines no value for.
		return !hasDomain && !hasElements && orderedBound(d)
	case "set":
		return !hasDomain && !hasBound && !d.SingleValued
	case "scalar":
		return !hasDomain && !hasBound && !hasElements && !d.SingleValued
	default:
		return false
	}
}

// orderedBound reports whether a declared int bound is ordered. A partial
// bound carries no domain at all, so it never reaches the comparison.
func orderedBound(d table.TagDecl) bool {
	return d.Min == nil || d.Max == nil || *d.Min <= *d.Max
}

// IntDomain enumerates a bounded int declaration's inclusive values. Both
// endpoints are inclusive, so `{0..3}` enumerates four values.
func IntDomain(d table.TagDecl) []int {
	if d.Kind != "int" || d.Min == nil || d.Max == nil {
		return nil
	}
	width, ok := intWidth(*d.Min, *d.Max)
	if !ok {
		return nil
	}
	// Counted, not bounded by `n <= *d.Max` — see `valueAssignments`: a
	// domain ending at MaxInt wraps its own loop variable rather than
	// passing the test, so the exact width drives the iteration.
	out := make([]int, 0, width)
	for n, i := *d.Min, 0; i < width; n, i = n+1, i+1 {
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
		if d.Provenance == table.ProvenanceRecognized {
			// The freshly recognized outcome is not a key the view supplies:
			// the kernel binds it from the recognized outcome itself
			// (`resolve.go::assemble`). Reading its declaration as an
			// obligation on the view would make every view non-conforming
			// and so vacuous every claim conformance is the premise of.
			continue
		}
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
