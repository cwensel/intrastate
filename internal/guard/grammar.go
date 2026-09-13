package guard

import (
	"encoding/json"
	"slices"
	"strconv"

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
		held, ok := parseHeldSet(value)
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
	// the vocabulary is frozen, so a free-form expression string is not an
	// operator and has no reading.
	return resolve.GuardUnevaluable
}

// parseSetLiteral decodes an authored set LITERAL from its §D13 canonical
// JSON array. A NON-EMPTY set is the published literal shape for both
// set-shaped operators, so an empty array is not a well-formed literal.
func parseSetLiteral(s string) ([]string, bool) {
	members, ok := parseHeldSet(s)
	if !ok || len(members) == 0 {
		return nil, false
	}
	return members, true
}

// parseHeldSet decodes a HELD set-valued tag value from its §D13 canonical
// JSON array.
//
// The non-empty requirement is a rule about the authored LITERAL — the
// operator/kind matrix publishes `contains`' literal shape as a "non-empty
// typed element set" — and it does NOT carry over to the held value. A
// set-valued tag holds any SUBSET of its declared element universe, the
// empty subset included: the assignment-count table makes the dimension
// `2^|universe|` by construction, and `Conforms` admits a held `[]` as a
// conforming view.
//
// Containment is total over sets, so `contains ["x"]` against a held `[]`
// has an answer — FALSE, since the held set does not contain `x`. Applying
// the literal's non-empty rule to the held value refused instead, and the
// kernel then refused `guard_unevaluable` on a conforming view that lint's
// own scoped product enumerates. That disagreement is the false
// exhaustiveness claim REQ-67 forbids, reached from the runtime side.
func parseHeldSet(s string) ([]string, bool) {
	var members []string
	if err := json.Unmarshal([]byte(s), &members); err != nil {
		return nil, false
	}
	// `null` decodes without error into a NIL slice, while `[]` decodes into a
	// non-nil empty one — so the nil check is exactly the line between a value
	// that is not a set at all and the empty subset. Without it a held `null`
	// would read as `[]` and `contains` would DECIDE it false, converting an
	// unevaluable into a decided verdict — the direction `0007:C1` forbids
	// ("never to false and never to true"), and the mirror image of the defect
	// D14 fixed. Lint never renders `null` (`renderSet` marshals a []string),
	// so only a runtime caller can supply it.
	if members == nil {
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
