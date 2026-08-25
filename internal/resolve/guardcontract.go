package resolve

import "testing"

// TestGuardEvaluatorContract is the cross-RDR contract test for the value
// seam. RDR 0003's implement stage instantiates it against its evaluator,
// so it lives in non-test source and is called by name.
//
// It drives the seam over the present value × literal × operator product
// and enforces the obligation the seam exists to carry: a present value a
// typed operator cannot parse is unevaluable, never false. Absence is not
// part of the contract — the kernel decides presence before the seam is
// ever asked (`0007:C1`), so every case here carries a present value.
func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator) {
	t.Helper()

	cases := []struct {
		name     string
		operator string
		literal  string
		value    string
		want     GuardResult
	}{
		// Equality over a present value: decided both ways against one
		// literal, so a constant-answering seam cannot pass.
		{"eq/equal", "eq", "Draft", "Draft", GuardTrue},
		{"eq/unequal", "eq", "Draft", "Final", GuardFalse},

		// Bounded integer comparison: decided both ways against one
		// literal, then the unparseable leg.
		{"gte/above", "gte", "3", "4", GuardTrue},
		{"gte/equal", "gte", "3", "3", GuardTrue},
		{"gte/below", "gte", "3", "2", GuardFalse},
		{"gte/unparseable value", "gte", "3", "many", GuardUnevaluable},
		{"gte/unparseable literal", "gte", "three", "3", GuardUnevaluable},

		// Membership over a §D13 set literal: a canonical JSON array,
		// members sorted, duplicate-free, compact encoding.
		{"in/member", "in", `["alpha","beta"]`, "alpha", GuardTrue},
		{"in/non-member", "in", `["alpha","beta"]`, "gamma", GuardFalse},

		// Set containment: both sides cross in the §D13 form.
		{"contains/superset", "contains", `["alpha"]`, `["alpha","beta"]`, GuardTrue},
		{"contains/missing member", "contains", `["gamma"]`, `["alpha","beta"]`, GuardFalse},
		{"contains/unparseable value", "contains", `["alpha"]`, "alpha", GuardUnevaluable},
	}

	for _, tc := range cases {
		atom := GuardAtom{
			Key:      "subject",
			Operator: tc.operator,
			Literal:  tc.literal,
			Block:    BlockAll,
		}
		got := seam.Evaluate(atom, tc.value)
		if got == tc.want {
			continue
		}
		if tc.want == GuardUnevaluable && got == GuardFalse {
			t.Errorf("%s: Evaluate(%+v, %q) = GuardFalse; a present value the "+
				"operator cannot parse is unevaluable, never false", tc.name, atom, tc.value)
			continue
		}
		t.Errorf("%s: Evaluate(%+v, %q) = %v; want %v", tc.name, atom, tc.value, got, tc.want)
	}
}
