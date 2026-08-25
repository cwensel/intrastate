package guard

import (
	"encoding/json"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/newcoinc/intrastate/internal/table"
)

// Assignment is one point of the scoped product: dimension to rendered
// held value, plus a presence marker for a key's presence dimension.
type Assignment map[string]string

// Presence names one value of a key's presence dimension.
type Presence string

// The two values of a presence dimension, plus the marker for an atom that
// selects neither — it reads the value dimension and leaves presence free.
const (
	PresencePresent Presence = "present"
	PresenceAbsent  Presence = "absent"
	PresenceEither  Presence = "either"
)

// presenceSuffix names a key's PRESENCE dimension apart from its VALUE
// dimension. The two sit alongside each other in one assignment, so they
// need distinct dimension names, and the separator is a byte no tag key
// carries — RDR 0002's key grammar admits no NUL.
const presenceSuffix = "\x00presence"

// presenceDim renders a key's presence dimension name.
func presenceDim(key string) string { return key + presenceSuffix }

// valueKeyOf recovers the tag key a dimension name belongs to, and whether
// the dimension is the presence half.
func valueKeyOf(dim string) (string, bool) {
	if key, ok := strings.CutSuffix(dim, presenceSuffix); ok {
		return key, true
	}
	return dim, false
}

// AssignmentSet is a set of assignments over a scoped product.
//
// It carries the dimensions it ranges over explicitly, so an EMPTY set over
// two dimensions stays distinguishable from an empty set over none, and
// `Contains` knows which dimensions of a candidate assignment to read.
//
// The zero value is the unprojectable set: no dimensions, no members, and
// `Projectable` false. That is deliberate — an atom lint cannot project is
// not an atom denoting nothing, and the two must not be confused.
type AssignmentSet struct {
	dims        []string
	members     map[string]Assignment
	projectable bool
}

// newSet builds an empty projectable set over dims.
func newSet(dims []string) AssignmentSet {
	d := slices.Clone(dims)
	slices.Sort(d)
	return AssignmentSet{dims: d, members: map[string]Assignment{}, projectable: true}
}

// add inserts one assignment, keyed by its canonical rendering.
func (s *AssignmentSet) add(a Assignment) {
	s.members[renderAssignment(s.dims, a)] = restrict(s.dims, a)
}

// renderAssignment renders an assignment over dims injectively, so two
// assignments agreeing on every dimension share one key.
func renderAssignment(dims []string, a Assignment) string {
	var b strings.Builder
	for _, d := range dims {
		writeField(&b, d)
		writeField(&b, a[d])
	}
	return b.String()
}

// restrict projects an assignment onto dims.
func restrict(dims []string, a Assignment) Assignment {
	out := make(Assignment, len(dims))
	for _, d := range dims {
		out[d] = a[d]
	}
	return out
}

// Len reports how many assignments the set holds.
func (s AssignmentSet) Len() int { return len(s.members) }

// Contains reports whether a is in the set. The candidate is projected onto
// the set's own dimensions first, so a denotation over one key answers for
// an assignment naming only that key.
func (s AssignmentSet) Contains(a Assignment) bool {
	if !s.projectable {
		return false
	}
	_, ok := s.members[renderAssignment(s.dims, a)]
	return ok
}

// Equal reports set equality. Two unprojectable sets are never equal: there
// is no subset of the product either of them names.
func (s AssignmentSet) Equal(other AssignmentSet) bool {
	if !s.projectable || !other.projectable {
		return false
	}
	if !slices.Equal(s.dims, other.dims) || len(s.members) != len(other.members) {
		return false
	}
	for k := range s.members {
		if _, ok := other.members[k]; !ok {
			return false
		}
	}
	return true
}

// Subset reports whether the receiver is contained in other.
func (s AssignmentSet) Subset(other AssignmentSet) bool {
	if !s.projectable || !other.projectable {
		return false
	}
	for _, a := range s.members {
		if !other.Contains(a) {
			return false
		}
	}
	return true
}

// Intersect returns the intersection.
func (s AssignmentSet) Intersect(other AssignmentSet) AssignmentSet {
	if !s.projectable || !other.projectable {
		return AssignmentSet{}
	}
	out := newSet(unionDims(s.dims, other.dims))
	for _, a := range s.members {
		if other.Contains(a) {
			out.add(a)
		}
	}
	return out
}

// Union returns the union.
func (s AssignmentSet) Union(other AssignmentSet) AssignmentSet {
	if !s.projectable || !other.projectable {
		return AssignmentSet{}
	}
	out := newSet(unionDims(s.dims, other.dims))
	for _, a := range s.members {
		out.add(a)
	}
	for _, a := range other.members {
		out.add(a)
	}
	return out
}

// subtract removes every assignment of other from the receiver. It is how
// the row's `unless` block is applied: the single conjunctive assignment
// set the FULL block matches is subtracted, never per-atom negation.
func (s AssignmentSet) subtract(other AssignmentSet) AssignmentSet {
	if !s.projectable || !other.projectable {
		return AssignmentSet{}
	}
	out := newSet(s.dims)
	for _, a := range s.members {
		if !other.Contains(a) {
			out.add(a)
		}
	}
	return out
}

func unionDims(a, b []string) []string {
	out := slices.Clone(a)
	for _, d := range b {
		if !slices.Contains(out, d) {
			out = append(out, d)
		}
	}
	slices.Sort(out)
	return out
}

// Presence reports which presence value the set selects. A set every member
// of which marks its key present selects `present`; likewise `absent`; a
// set spanning both — or one carrying no presence dimension at all —
// selects `either`.
func (s AssignmentSet) Presence() Presence {
	if !s.projectable {
		return ""
	}
	var dim string
	for _, d := range s.dims {
		if _, isPresence := valueKeyOf(d); isPresence {
			dim = d
			break
		}
	}
	if dim == "" {
		return PresenceEither
	}

	seen := map[string]bool{}
	for _, a := range s.members {
		seen[a[dim]] = true
	}
	switch {
	case len(seen) == 1 && seen[string(PresencePresent)]:
		return PresencePresent
	case len(seen) == 1 && seen[string(PresenceAbsent)]:
		return PresenceAbsent
	default:
		return PresenceEither
	}
}

// Values reports the value subset the set selects on its key's value
// dimension, sorted.
func (s AssignmentSet) Values() []string {
	if !s.projectable {
		return nil
	}
	var dim string
	for _, d := range s.dims {
		if _, isPresence := valueKeyOf(d); !isPresence {
			dim = d
			break
		}
	}
	if dim == "" {
		return nil
	}
	seen := map[string]bool{}
	for _, a := range s.members {
		seen[a[dim]] = true
	}
	return slices.Sorted(maps.Keys(seen))
}

// Projectable reports whether lint can project the atom at all.
//
// An atom over a dimension the model does not declare single-valued, or one
// with no finite declared domain, is not a differently-projecting atom — it
// is one lint CANNOT project. Resolving it as "the literal is among the
// held values", "the held set equals the literal", or "the held set is
// contained in the literal" yields different union cardinalities and
// different overlap verdicts on the same model, so lint picks no reading.
func (s AssignmentSet) Projectable() bool { return s.projectable }

// --- dimension enumeration -----------------------------------------------

// valueAssignments enumerates the values a key's VALUE dimension ranges
// over, rendered the way a held value crosses the seam.
func valueAssignments(d table.TagDecl) ([]string, bool) {
	if !agrees(d) {
		return nil, false
	}
	switch d.Kind {
	case "enum":
		if len(d.Domain) == 0 {
			return nil, false
		}
		return spreadValues(d.Domain, d.SingleValued), true
	case "bool":
		return spreadValues([]string{"false", "true"}, d.SingleValued), true
	case "int":
		if d.Min == nil || d.Max == nil {
			return nil, false
		}
		domain := make([]string, 0, *d.Max-*d.Min+1)
		for n := *d.Min; n <= *d.Max; n++ {
			domain = append(domain, strconv.Itoa(n))
		}
		return spreadValues(domain, d.SingleValued), true
	case "set":
		if len(d.Elements) == 0 {
			return nil, false
		}
		// Values co-occur by construction: the space is the powerset of the
		// element universe, never the universe itself.
		return subsets(d.Elements), true
	default:
		return nil, false
	}
}

// spreadValues renders the assignments a value dimension ranges over. Under
// the single-valued marker exactly one value holds, so each domain value is
// one assignment; absent it the values could co-occur, so each SUBSET is
// one assignment, rendered as the canonical JSON array a co-occurring held
// value crosses as.
func spreadValues(domain []string, singleValued bool) []string {
	if singleValued {
		return slices.Clone(domain)
	}
	return subsets(domain)
}

// subsets renders every subset of universe as a canonical JSON array —
// members sorted, duplicate-free, compact (JDR 0001 §D13).
func subsets(universe []string) []string {
	sorted := slices.Compact(slices.Sorted(slices.Values(universe)))
	out := make([]string, 0, 1<<len(sorted))
	for mask := range 1 << len(sorted) {
		members := []string{}
		for i, e := range sorted {
			if mask&(1<<i) != 0 {
				members = append(members, e)
			}
		}
		out = append(out, renderSet(members))
	}
	return out
}

// renderSet encodes a member sequence as its canonical JSON array.
func renderSet(members []string) string {
	b, err := json.Marshal(members)
	if err != nil {
		// json.Marshal of a []string cannot fail.
		return "[]"
	}
	return string(b)
}
