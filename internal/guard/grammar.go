package guard

import (
	"slices"

	"github.com/cwensel/intrastate/internal/resolve"
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

// operators is the frozen typed vocabulary in matrix order: equality,
// membership, the four bounded integer comparisons, existence, and set
// containment.
var operators = []string{"eq", "in", "lt", "lte", "gt", "gte", "exists", "contains"}

// kinds is the five value-kind tokens, spelled the one way a kind is
// spelled wherever it is named.
var kinds = []string{"enum", "bool", "int", "set", "scalar"}

// Operators returns this RDR's frozen operator vocabulary (`0029:C4`).
func Operators() []string { return slices.Clone(operators) }

// KnownOperator reports whether token is in the frozen vocabulary. The set
// admits no member beyond those below, so a free-form expression string is
// simply not an operator.
func KnownOperator(token string) bool { return slices.Contains(operators, token) }

// Kinds returns the five value-kind tokens.
func Kinds() []string { return slices.Clone(kinds) }

// Accepts reports whether the operator/kind matrix admits the pair. An
// operator outside the frozen vocabulary is accepted by no kind.
func Accepts(operator, kind string) bool {
	return slices.Contains(matrix[operator].kinds, kind)
}

// LiteralShape returns the operator's published literal shape.
func LiteralShape(operator string) Shape { return matrix[operator].shape }

// SingleValueOperator reports whether the operator narrows the tag's
// SINGLE held value: `eq`, `in`, and the four integer comparisons. Such an
// atom projects only over a key the model declares single-valued.
func SingleValueOperator(operator string) bool {
	switch operator {
	case "eq", "in", "lt", "lte", "gt", "gte":
		return true
	}
	return false
}

// --- the value seam ------------------------------------------------------

// Evaluator is RDR 0003's value-comparison seam. It is DECLARED in
// `internal/resolve`, beside the kernel it serves, so the loader can compare
// through the same seam the runtime does (`0030:C1`, JDR 0004 JD-1); this
// package re-exports it under its original name.
type Evaluator = resolve.Evaluator

// NewEvaluator returns the seam constructed over kinds, the tag key → kind
// token mapping DeclaredKinds produces for the model whose rows it will
// evaluate. The return is the concrete Evaluator, never widened to the
// kernel's interface.
func NewEvaluator(kinds map[string]string) Evaluator {
	return resolve.NewEvaluator(kinds)
}
