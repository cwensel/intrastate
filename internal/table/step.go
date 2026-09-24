package table

import (
	"fmt"
	"iter"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/cwensel/intrastate/internal/resolve"
)

// stepKey is the one key the computed write value admits: `{ step = <n> }`
// (`0030:C1`, `0030:D-wire-byte-format`).
const stepKey = "step"

// stepSpec is one authored step write, validated for its grammar and kind
// but not yet walked over its cells.
type stepSpec struct {
	key  string
	decl TagDecl
	n    int
}

// stepPoint is a validated step resolved into the per-cell lists `expand`
// consumes: the ADMITTED cells in the tag's own domain order, and at each
// the literal the row writes.
type stepPoint struct {
	key    string
	cells  []string
	values []string
}

// parseStep validates a write-block table value as the one admitted
// computed form over decl, returning the refusal detail when it is not.
//
// Every refusal here is on the write-block path and so takes
// `malformed_tag_declaration` at the caller (`0030:C1`).
func parseStep(decl TagDecl, raw map[string]any) (int, string) {
	const form = "a computed write value is the one-key table { step = <n> }, n a non-zero integer"
	v, ok := raw[stepKey]
	if len(raw) != 1 || !ok {
		return 0, form
	}
	n, ok := v.(int64)
	if !ok || n == 0 {
		return 0, form
	}

	const kinds = "a step is admitted only on an int declaring both min and max, " +
		"or an enum with a non-empty domain"
	switch {
	case decl.Kind == "int" && decl.Min != nil && decl.Max != nil:
		// The SAME width rule lint sizes the declaration with: a width lint
		// cannot compute is not one the loader enumerates. It refuses the
		// same way the kind arm does, naming the admitted kinds (`0030:F1`).
		if _, ok := resolve.IntWidth(*decl.Min, *decl.Max); !ok {
			return 0, fmt.Sprintf("a step cannot enumerate int bound {%d..%d}: its width is not representable; %s",
				*decl.Min, *decl.Max, kinds)
		}
	case decl.Kind == "enum" && len(decl.Domain) > 0:
		// Both guards are step-scoped: the step indexes positions and mints
		// the cell as a suffix element, so a repeated member or one carrying
		// the separator is refused here, never on the declaration.
		if dup, ok := firstDuplicate(decl.Domain); ok {
			return 0, "a step cannot index a domain that repeats the member " + strconv.Quote(dup)
		}
		for _, m := range decl.Domain {
			if strings.Contains(m, suffixSep) {
				return 0, "a step cannot suffix the member " + strconv.Quote(m) +
					", which contains the expansion-suffix separator " + suffixSep
			}
		}
	default:
		return 0, kinds
	}
	return int(n), ""
}

// stepCells yields a stepped declaration's CANDIDATE cells in its own domain
// order, each with its domain index: the authored order for `enum`, and
// ascending for `int`.
//
// An `int` walk is narrowed by the positive atoms before any cell is
// visited, so its cost follows the cells the atoms can admit rather than
// the declared width: a representable domain as wide as `math.MaxInt` is
// admitted (`0030:C1`, "representability, not size"), and sizing the walk
// by it crashed or never returned. The narrowing only drops members some
// positive atom must answer false at — the ordered bounds, and the members
// an `eq`/`in` names — so every candidate is still decided by `admits` and
// the admitted set is unchanged.
func stepCells(decl TagDecl, positive []Atom) iter.Seq2[int, string] {
	if decl.Kind == "enum" {
		return func(yield func(int, string) bool) {
			for i, m := range decl.Domain {
				if !yield(i, m) {
					return
				}
			}
		}
	}

	lo, hi := *decl.Min, *decl.Max
	var members []int
	pinned, empty := false, false
	for _, a := range positive {
		var v int
		if a.Operator == "lt" || a.Operator == "lte" || a.Operator == "gt" || a.Operator == "gte" {
			var err error
			if v, err = strconv.Atoi(a.Literal[0]); err != nil {
				continue
			}
		}
		switch a.Operator {
		case "eq", "in":
			if ns, ok := intMembers(a.Literal); ok && (!pinned || len(ns) < len(members)) {
				members, pinned = ns, true
			}
		case "lt":
			// `lt MinInt` admits nothing; the decrement would wrap.
			empty = empty || v == math.MinInt
			if v != math.MinInt {
				hi = min(hi, v-1)
			}
		case "lte":
			hi = min(hi, v)
		case "gt":
			empty = empty || v == math.MaxInt
			if v != math.MaxInt {
				lo = max(lo, v+1)
			}
		case "gte":
			lo = max(lo, v)
		}
	}

	return func(yield func(int, string) bool) {
		if empty || lo > hi {
			return
		}
		if pinned {
			for _, n := range members {
				if n >= lo && n <= hi && !yield(n-*decl.Min, strconv.Itoa(n)) {
					return
				}
			}
			return
		}
		// Counted, not bounded by `n <= hi`, so a walk ending at MaxInt
		// terminates.
		width, _ := resolve.IntWidth(lo, hi)
		for n, i := lo, 0; i < width; n, i = n+1, i+1 {
			if !yield(n-*decl.Min, strconv.Itoa(n)) {
				return
			}
		}
	}
}

// intMembers parses an `eq`/`in` literal's members as ints, ascending and
// duplicate-free, or reports that one does not parse.
func intMembers(literal []string) ([]int, bool) {
	out := make([]int, 0, len(literal))
	for _, s := range literal {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, false
		}
		out = append(out, n)
	}
	slices.Sort(out)
	return slices.Compact(out), true
}

// resolveStep walks a step's cells in domain order, admits each that every
// POSITIVE atom the rule authors on the tag answers true at, and computes
// the stepped literal at each admitted cell.
//
// A cell answering anything but true — false or undecided — is not
// admitted: the loader's policy excludes the undecided arm (`0030:C1`).
// `unless` atoms are not consulted. The first admitted cell whose stepped
// value leaves the domain refuses the rule (`0030:C2`).
func (l *loader) resolveStep(id string, spec stepSpec, predicates []Atom, ev resolve.Evaluator) (stepPoint, error) {
	bad := func(detail string) error {
		return fail(CatMalformedTagDeclaration, "rule "+id+" write "+spec.key+": "+detail)
	}

	var positive, unless []Atom
	for _, a := range predicates {
		switch {
		case a.Key != spec.key:
		case a.Block == BlockUnless:
			unless = append(unless, a)
		default:
			positive = append(positive, a)
		}
	}

	point := stepPoint{key: spec.key}
	for i, cell := range stepCells(spec.decl, positive) {
		if !admits(positive, cell, ev) {
			continue
		}
		value, failure := stepValue(spec, i, cell)
		if failure != "" {
			return stepPoint{}, bad(failure + unlessNote(unless))
		}
		// The stepped literal conforms exactly as an authored literal write
		// does (0030:C2), and a literal write of the sentinel is refused
		// (0002:C11): a domain may declare `<clear>` as a member, so a step
		// landing on it would otherwise mint a write every consumer reads as
		// removal.
		if value == ClearSentinel {
			return stepPoint{}, fail(CatReservedTagValue, "rule "+id+" write "+spec.key+
				": cell "+cell+" steps to the reserved value "+ClearSentinel)
		}
		point.cells = append(point.cells, cell)
		point.values = append(point.values, value)
	}
	if len(point.cells) == 0 {
		return stepPoint{}, bad("the step admits no cell: no member of " + spec.key +
			"'s domain satisfies every positive atom the rule authors on it")
	}
	return point, nil
}

// admits reports whether every positive atom answers TRUE at cell.
func admits(atoms []Atom, cell string, ev resolve.Evaluator) bool {
	for _, a := range atoms {
		if a.Operator == resolve.OpExists {
			// Existence is decided from presence, and a cell is a present
			// value.
			if a.Literal[0] != resolve.LiteralTrue {
				return false
			}
			continue
		}
		if resolve.CompareValue(ev, a.Key, a.Operator, a.Block, a.Literal, cell) != resolve.GuardTrue {
			return false
		}
	}
	return true
}

// stepValue computes the literal `n` positions from the cell at domain
// index i, or the bound failure naming the cell first and then the result.
func stepValue(spec stepSpec, i int, cell string) (string, string) {
	decl, n := spec.decl, spec.n
	if decl.Kind == "enum" {
		last := len(decl.Domain) - 1
		switch {
		case n > 0 && n > last-i:
			return "", fmt.Sprintf("cell %s steps %d past the last member", cell, n-(last-i))
		case n < 0 && i+n < 0:
			// uint64 of the negation reads MinInt's magnitude correctly.
			return "", fmt.Sprintf("cell %s steps %d before the first member", cell, uint64(-(i + n)))
		}
		return decl.Domain[i+n], ""
	}

	c, _ := strconv.Atoi(cell)
	// A CHECKED add, ordered before the render: a cell the step carries out
	// of int range is a bound failure naming the bound it passed.
	switch {
	case n > 0 && c > math.MaxInt-n:
		return "", fmt.Sprintf("cell %d steps out of int range, which is above max %d", c, *decl.Max)
	case n < 0 && c < math.MinInt-n:
		return "", fmt.Sprintf("cell %d steps out of int range, which is below min %d", c, *decl.Min)
	}
	v := c + n
	switch {
	case v > *decl.Max:
		return "", fmt.Sprintf("cell %d steps to %d, which is above max %d", c, v, *decl.Max)
	case v < *decl.Min:
		return "", fmt.Sprintf("cell %d steps to %d, which is below min %d", c, v, *decl.Min)
	}
	return strconv.Itoa(v), ""
}

// unlessNote names the `unless` atoms a bound refusal's rule authors on the
// stepped tag: they did not exclude the cell, because a step never consults
// them, and the exclusion has to be written as a positive atom.
func unlessNote(unless []Atom) string {
	if len(unless) == 0 {
		return ""
	}
	names := make([]string, 0, len(unless))
	for _, a := range unless {
		names = append(names, a.Key+" "+a.Operator+" "+strings.Join(a.Literal, ","))
	}
	return "; the guard.unless atom " + strings.Join(names, ", ") +
		" did not exclude the cell, since a step never consults unless — write the exclusion as a guard.all atom"
}

// tagKinds is the loader's key → declared kind mapping, a key with an empty
// kind omitted, for constructing the value seam over this model.
func (l *loader) tagKinds() map[string]string {
	out := make(map[string]string, len(l.model.Tags))
	for key, decl := range l.model.Tags {
		if decl.Kind != "" {
			out[key] = decl.Kind
		}
	}
	return out
}
