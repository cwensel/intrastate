package resolve

import (
	"encoding/json"
	"math"
	"slices"
	"strconv"
	"strings"
)

// Evaluator is RDR 0003's value-comparison seam: it decides value
// semantics over a PRESENT value and never reads the tag view.
//
// It holds the declaration mapping it was constructed over — tag key to
// declared kind token — and NOTHING else: no view, no runtime tag value.
// That keeps "never reads the tag view" a property of the type rather than
// a discipline, while letting `eq`/`in` compare under the declared kind
// (RDR 0012). Key presence, existence atoms, absent-key unevaluability, and
// the combination of per-atom verdicts are the kernel's (RDR 0007).
//
// It lives here rather than in `internal/guard` so the loader, which
// `internal/guard` imports, can compare through it too (`0030:C1`, JDR 0004
// JD-1); `guard.Evaluator` re-exports it.
//
// Construct it with NewEvaluator. A zero value holds no mapping, so every
// `eq`/`in` atom answers unevaluable rather than reverting to raw-string
// comparison.
type Evaluator struct {
	kinds map[string]string
}

// NewEvaluator returns the seam constructed over kinds, the tag key → kind
// token mapping of the model whose rows it will evaluate.
func NewEvaluator(kinds map[string]string) Evaluator {
	return Evaluator{kinds: kinds}
}

// Evaluate decides one atom against the present value its key holds.
//
// A present value the operator cannot parse is UNEVALUABLE, never false:
// that is the obligation the seam exists to carry, and folding it into
// false would let a malformed literal prune a row silently.
func (e Evaluator) Evaluate(atom GuardAtom, value string) GuardResult {
	switch atom.Operator {
	case "eq", "in":
		return e.typedEquality(atom, value)
	case "lt", "lte", "gt", "gte":
		// The ordering operators keep their operator-inferred integer parse:
		// the matrix admits them over `int` alone, so the kind mapping is not
		// consulted.
		bound, err := strconv.Atoi(atom.Literal)
		if err != nil {
			return GuardUnevaluable
		}
		held, err := strconv.Atoi(value)
		if err != nil {
			return GuardUnevaluable
		}
		return boolResult(compareInts(atom.Operator, held, bound))
	case "contains":
		want, ok := parseSetLiteral(atom.Literal)
		if !ok {
			return GuardUnevaluable
		}
		held, ok := parseHeldSet(value)
		if !ok {
			// A `contains` value crosses the seam as the §D13 canonical JSON
			// array. A bare string is not a held set, so the comparison has
			// no answer — unevaluable, never false.
			return GuardUnevaluable
		}
		for _, e := range want {
			if !slices.Contains(held, e) {
				return GuardFalse
			}
		}
		return GuardTrue
	}
	// `exists` is the kernel's: it is decided from presence alone and never
	// reaches this seam, so being handed one means answering from a view
	// this evaluator does not have. Every unknown operator lands here too —
	// the vocabulary is frozen, so a free-form expression string is not an
	// operator and has no reading.
	return GuardUnevaluable
}

// typedEquality decides `eq`/`in` under the atom key's declared kind.
//
// The literal side is the one member of `eq` or the members decoded from
// `in`'s §D13 array, and every member must parse under the kind before any
// comparison: one unparseable member poisons the whole list. A key the
// mapping does not carry, a `set` kind (the matrix admits no `eq`/`in`
// over it), and a token outside the five-kind vocabulary are all atoms the
// seam cannot type — unevaluable, never a raw-string comparison.
func (e Evaluator) typedEquality(atom GuardAtom, value string) GuardResult {
	kind, ok := e.kinds[atom.Key]
	if !ok {
		return GuardUnevaluable
	}
	literals := []string{atom.Literal}
	if atom.Operator == "in" {
		members, ok := parseSetLiteral(atom.Literal)
		if !ok {
			return GuardUnevaluable
		}
		literals = members
	}

	switch kind {
	case "int":
		held, err := strconv.Atoi(value)
		if err != nil {
			return GuardUnevaluable
		}
		parsed := make([]int, 0, len(literals))
		for _, l := range literals {
			n, err := strconv.Atoi(l)
			if err != nil {
				return GuardUnevaluable
			}
			parsed = append(parsed, n)
		}
		return boolResult(slices.Contains(parsed, held))
	case "bool":
		if !isBoolToken(value) {
			return GuardUnevaluable
		}
		for _, l := range literals {
			if !isBoolToken(l) {
				return GuardUnevaluable
			}
		}
		return boolResult(slices.Contains(literals, value))
	case "enum", "scalar":
		return boolResult(slices.Contains(literals, value))
	case "set":
		return GuardUnevaluable
	default:
		return GuardUnevaluable
	}
}

// CompareValue is the value comparison every static caller shares — the
// loader's admitted-cell filter, lint's product walk and graph-lint's
// reachability: it renders a normalized atom's member sequence into the
// seam's literal and returns the evaluator's THREE-valued verdict.
//
// It is the comparison, never the policy: what an undecided verdict means
// is each caller's own (`0030` IP Phase 2).
func CompareValue(ev Evaluator, key, operator string, block Block, members []string, held string) GuardResult {
	return ev.Evaluate(GuardAtom{
		Key:      key,
		Operator: operator,
		Literal:  renderLiteral(operator, members),
		Block:    block,
	}, held)
}

// renderLiteral renders a member sequence as the seam literal: the one
// member of a single-value operator, or the CANONICAL JSON array of a
// set-shaped one — sorted and duplicate-free, since a set literal's two
// spellings are one literal (JDR 0001 §D13).
func renderLiteral(operator string, members []string) string {
	if operator != "in" && operator != "contains" {
		return strings.Join(members, "")
	}
	canonical := slices.Compact(slices.Sorted(slices.Values(members)))
	if canonical == nil {
		canonical = []string{}
	}
	b, err := json.Marshal(canonical)
	if err != nil {
		// json.Marshal of a []string cannot fail.
		return "[]"
	}
	return string(b)
}

// IntWidth returns how many values the inclusive bound `{minV..maxV}`
// names, and whether that count fits in an int.
//
// `maxV - minV + 1` is evaluated in WRAPPING int arithmetic, so a domain
// as wide as the int range — `{MinInt..MaxInt}`, `{0..MaxInt}`,
// `{MinInt+1..MaxInt}` — yields zero or a negative before any ceiling
// applies. A wrapped width read as a tiny domain would be certified over an
// unbounded dimension and then enumerated one value at a time, which does
// not terminate. It is the ONE width rule: lint sizes a declaration with it
// and the loader refuses to step one it cannot represent (`0030:C1`).
func IntWidth(minV, maxV int) (int, bool) {
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

// isBoolToken reports whether s is one of the two boolean tokens.
func isBoolToken(s string) bool {
	return s == "true" || s == "false"
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
	// ("never to false and never to true"). Lint never renders `null`, so
	// only a runtime caller can supply it.
	if members == nil {
		return nil, false
	}
	return members, true
}

func compareInts(operator string, held, bound int) bool {
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
