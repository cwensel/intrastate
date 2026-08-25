package guard

import (
	"encoding/json"
	"slices"
	"strconv"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// Shape names an operator's literal shape, as the operator/kind matrix
// publishes it.
type Shape string

// The literal shapes the operator/kind matrix publishes.
const (
	ShapeScalar     Shape = "one typed scalar"
	ShapeScalarSet  Shape = "non-empty typed scalar set"
	ShapeInteger    Shape = "one typed integer"
	ShapeBoolean    Shape = "boolean"
	ShapeElementSet Shape = "non-empty typed element set"
)

// operatorRow is one row of the operator/kind matrix: the kinds the
// operator accepts and the literal shape it publishes.
type operatorRow struct {
	kinds []string
	shape Shape
}

// matrix is the operator/kind matrix verbatim. `exists` accepts any kind —
// it reads presence, not value; the matrix cell's "provided the tag is
// declared optional" qualifier is a VACUITY report rather than a rejection
// ("well-formed but vacuous — lint reports it as such rather than rejecting
// it"), so it lives in lint and not here.
var matrix = map[string]operatorRow{
	"eq":       {kinds: []string{"enum", "bool", "int", "scalar"}, shape: ShapeScalar},
	"in":       {kinds: []string{"enum", "bool", "int", "scalar"}, shape: ShapeScalarSet},
	"lt":       {kinds: []string{"int"}, shape: ShapeInteger},
	"lte":      {kinds: []string{"int"}, shape: ShapeInteger},
	"gt":       {kinds: []string{"int"}, shape: ShapeInteger},
	"gte":      {kinds: []string{"int"}, shape: ShapeInteger},
	"exists":   {kinds: []string{"enum", "bool", "int", "set", "scalar"}, shape: ShapeBoolean},
	"contains": {kinds: []string{"set"}, shape: ShapeElementSet},
}

// operators is the closed typed vocabulary in matrix order: equality,
// membership, the four bounded integer comparisons, existence, and set
// containment.
var operators = []string{"eq", "in", "lt", "lte", "gt", "gte", "exists", "contains"}

// kinds is the five value-kind tokens, spelled the one way a kind is
// spelled wherever it is named.
var kinds = []string{"enum", "bool", "int", "set", "scalar"}

// Operators returns this RDR's closed operator vocabulary.
func Operators() []string { return slices.Clone(operators) }

// KnownOperator reports whether token is in the closed vocabulary. The set
// is closed, so a free-form expression string is simply not an operator.
func KnownOperator(token string) bool { return slices.Contains(operators, token) }

// Kinds returns the five value-kind tokens.
func Kinds() []string { return slices.Clone(kinds) }

// Accepts reports whether the operator/kind matrix admits the pair. An
// operator outside the closed vocabulary is accepted by no kind.
func Accepts(operator, kind string) bool {
	return slices.Contains(matrix[operator].kinds, kind)
}

// LiteralShape returns the operator's published literal shape.
func LiteralShape(operator string) Shape { return matrix[operator].shape }

// isSingleValueOperator reports whether the operator narrows the tag's
// SINGLE held value: `eq`, `in`, and the four integer comparisons. Such an
// atom projects only over a key the model declares single-valued.
func isSingleValueOperator(operator string) bool {
	switch operator {
	case "eq", "in", "lt", "lte", "gt", "gte":
		return true
	}
	return false
}

// --- the value seam ------------------------------------------------------

// Evaluator is RDR 0003's value-comparison seam: it decides value
// semantics over a PRESENT value and never reads the tag view.
//
// It carries no state, which is what makes "never reads the tag view" a
// property of the type rather than a discipline: there is nowhere to hold a
// view and the signature admits none. Key presence, existence atoms,
// absent-key unevaluability, and the combination of per-atom verdicts are
// the kernel's (RDR 0007).
type Evaluator struct{}

// Evaluate decides one atom against the present value its key holds.
//
// A present value the operator cannot parse is UNEVALUABLE, never false:
// that is the obligation the seam exists to carry, and folding it into
// false would let a malformed literal prune a row silently.
func (Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	switch atom.Operator {
	case "eq":
		return boolResult(value == atom.Literal)
	case "in":
		members, ok := parseSetLiteral(atom.Literal)
		if !ok {
			return resolve.GuardUnevaluable
		}
		return boolResult(slices.Contains(members, value))
	case "lt", "lte", "gt", "gte":
		bound, err := strconv.Atoi(atom.Literal)
		if err != nil {
			return resolve.GuardUnevaluable
		}
		held, err := strconv.Atoi(value)
		if err != nil {
			return resolve.GuardUnevaluable
		}
		return boolResult(compare(atom.Operator, held, bound))
	case "contains":
		want, ok := parseSetLiteral(atom.Literal)
		if !ok {
			return resolve.GuardUnevaluable
		}
		held, ok := parseSetLiteral(value)
		if !ok {
			// A `contains` value crosses the seam as the §D13 canonical JSON
			// array. A bare string is not a held set, so the comparison has
			// no answer — unevaluable, never false.
			return resolve.GuardUnevaluable
		}
		for _, e := range want {
			if !slices.Contains(held, e) {
				return resolve.GuardFalse
			}
		}
		return resolve.GuardTrue
	}
	// `exists` is the kernel's: it is decided from presence alone and never
	// reaches this seam, so being handed one means answering from a view
	// this evaluator does not have. Every unknown operator lands here too —
	// the vocabulary is closed, so a free-form expression string is not an
	// operator and has no reading.
	return resolve.GuardUnevaluable
}

// parseSetLiteral decodes a §D13 canonical JSON array. A NON-EMPTY set is
// the published literal shape for both set-shaped operators, so an empty
// array is not a well-formed literal.
func parseSetLiteral(s string) ([]string, bool) {
	var members []string
	if err := json.Unmarshal([]byte(s), &members); err != nil {
		return nil, false
	}
	if len(members) == 0 {
		return nil, false
	}
	return members, true
}

func compare(operator string, held, bound int) bool {
	switch operator {
	case "lt":
		return held < bound
	case "lte":
		return held <= bound
	case "gt":
		return held > bound
	default:
		return held >= bound
	}
}

func boolResult(b bool) resolve.GuardResult {
	if b {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}
