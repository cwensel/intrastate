package resolve

import (
	"maps"
	"slices"
	"testing"
)

// ContractKinds is the declaration fixture TestGuardEvaluatorContract's
// cases assume: tag key → declared kind token, one key per token of the
// five-kind vocabulary. It is a kernel-owned map (no declaration type
// crosses into the kernel), and every call returns a FRESH map, so no
// caller can mutate the fixture another caller reads.
//
// Only the contract test reads it. The resolution path — Resolve and its
// callees — reads no declarations.
func ContractKinds() map[string]string {
	return map[string]string{
		"phase":  "enum",
		"label":  "scalar",
		"count":  "int",
		"flag":   "bool",
		"labels": "set",
	}
}

// contractVocabulary is the five-kind vocabulary the fixture spans.
var contractVocabulary = []string{"bool", "enum", "int", "scalar", "set"}

// TestGuardEvaluatorContract is the cross-RDR contract test for the value
// seam. RDR 0003's implement stage instantiates it against its evaluator,
// so it lives in non-test source and is called by name.
//
// It takes a seam CONSTRUCTOR and builds the seam itself over
// ContractKinds, so an implementer cannot pass a seam built over some
// other declaration mapping. Each case names its own fixture key, chosen
// so the case's operator is one the operator/kind matrix admits for that
// key's kind.
//
// It drives the seam over the present value × literal × operator product
// and enforces the obligation the seam exists to carry: a present value a
// typed operator cannot parse is unevaluable, never false. Absence is not
// part of the contract — the kernel decides presence before the seam is
// ever asked (`0007:C1`), so every case here carries a present value.
func TestGuardEvaluatorContract(t *testing.T, newSeam func(kinds map[string]string) GuardEvaluator) {
	t.Helper()

	// Meta-check: the fixture spans exactly the five-kind vocabulary.
	kinds := ContractKinds()
	if got := slices.Sorted(maps.Values(kinds)); !slices.Equal(slices.Compact(got), contractVocabulary) {
		t.Errorf("ContractKinds carries kind tokens %v; want exactly %v", got, contractVocabulary)
	}

	type contractCase struct {
		name     string
		key      string
		operator string
		literal  string
		value    string
		want     GuardResult
	}
	cases := []contractCase{
		// Equality over a present value: decided both ways against one
		// literal, so a constant-answering seam cannot pass.
		{"eq/equal", "phase", "eq", "Draft", "Draft", GuardTrue},
		{"eq/unequal", "phase", "eq", "Draft", "Final", GuardFalse},

		// Bounded integer comparison: decided both ways against one
		// literal, then the unparseable leg.
		{"gte/above", "count", "gte", "3", "4", GuardTrue},
		{"gte/equal", "count", "gte", "3", "3", GuardTrue},
		{"gte/below", "count", "gte", "3", "2", GuardFalse},
		{"gte/unparseable value", "count", "gte", "3", "many", GuardUnevaluable},
		{"gte/unparseable literal", "count", "gte", "three", "3", GuardUnevaluable},

		// Membership over a §D13 set literal: a canonical JSON array,
		// members sorted, duplicate-free, compact encoding.
		{"in/member", "label", "in", `["alpha","beta"]`, "alpha", GuardTrue},
		{"in/non-member", "label", "in", `["alpha","beta"]`, "gamma", GuardFalse},

		// Set containment: both sides cross in the §D13 form.
		{"contains/superset", "labels", "contains", `["alpha"]`, `["alpha","beta"]`, GuardTrue},
		{"contains/missing member", "labels", "contains", `["gamma"]`, `["alpha","beta"]`, GuardFalse},
		{"contains/unparseable value", "labels", "contains", `["alpha"]`, "alpha", GuardUnevaluable},
		{"contains/null value", "labels", "contains", `["alpha"]`, "null", GuardUnevaluable},

		// `int`: comparison is over PARSED values; a held value that does
		// not parse is unevaluable, never false.
		{"int eq/parsed equal", "count", "eq", "7", "07", GuardTrue},
		{"int eq/parses but differs", "count", "eq", "7", "8", GuardFalse},
		{"int eq/unparseable value", "count", "eq", "7", "many", GuardUnevaluable},
		{"int in/parsed member", "count", "in", `["7","9"]`, "09", GuardTrue},
		{"int in/parses but differs", "count", "in", `["7","9"]`, "8", GuardFalse},
		{"int in/unparseable value", "count", "in", `["7","9"]`, "many", GuardUnevaluable},
		// Defensive: one non-integer member poisons the whole list, and an
		// out-of-range held value is unevaluable (it overflows 64-bit int).
		{"int in/non-integer member", "count", "in", `["7","x"]`, "7", GuardUnevaluable},
		{"int eq/overflow value", "count", "eq", "7", "99999999999999999999", GuardUnevaluable},

		// `bool`: token equality over `true` | `false`; any other spelling
		// is unevaluable.
		{"bool eq/token", "flag", "eq", "true", "true", GuardTrue},
		{"bool eq/other token", "flag", "eq", "true", "false", GuardFalse},
		{"bool eq/non-token value", "flag", "eq", "true", "yes", GuardUnevaluable},
		{"bool in/other token", "flag", "in", `["true"]`, "false", GuardFalse},
		{"bool in/non-token value", "flag", "in", `["true"]`, "1", GuardUnevaluable},

		// Defensive: the matrix admits no `eq`/`in` over `set`, and a key
		// absent from the fixture cannot be typed.
		{"set eq/value equals literal", "labels", "eq", "alpha", "alpha", GuardUnevaluable},
		{"absent key eq/value equals literal", "undeclared", "eq", "x", "x", GuardUnevaluable},
	}

	check := func(seam GuardEvaluator, tc contractCase) {
		t.Helper()
		atom := GuardAtom{
			Key:      tc.key,
			Operator: tc.operator,
			Literal:  tc.literal,
			Block:    BlockAll,
		}
		got := seam.Evaluate(atom, tc.value)
		if got == tc.want {
			return
		}
		if tc.want == GuardUnevaluable && got == GuardFalse {
			t.Errorf("%s: Evaluate(%+v, %q) = GuardFalse; a present value the "+
				"operator cannot parse is unevaluable, never false", tc.name, atom, tc.value)
			return
		}
		t.Errorf("%s: Evaluate(%+v, %q) = %v; want %v", tc.name, atom, tc.value, got, tc.want)
	}

	seam := newSeam(ContractKinds())
	for _, tc := range cases {
		check(seam, tc)
	}

	// Dispatch over the kind token is exhaustive, and a token outside the
	// vocabulary is unevaluable rather than compared. The seam is built over
	// a caller-local copy carrying one extra key, so the published fixture
	// itself stays exactly the five tokens.
	widened := ContractKinds()
	widened["measure"] = "decimal"
	check(newSeam(widened), contractCase{
		"unknown kind token eq", "measure", "eq", "7", "7", GuardUnevaluable,
	})
}
