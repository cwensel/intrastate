package graphlint_test

// RDR 0006 — the stable-reason table's `tag-not-single-valued` arm has ONE
// source for what counts as a single-value operator: `guard`'s
// classification. This file pins that as an invariant over the whole closed
// vocabulary rather than one member at a time.
//
// The prior spelling of this guarantee was a pair of behavioural tests, one
// per operator (`eq`, then `in`). That shape cannot see a NEW operator added
// to guard's vocabulary and left unclassified by lint, which is exactly the
// drift it was written to catch.

import (
	"fmt"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
)

// singleValueProbeDecls declares a FINITE dimension that carries NO
// single-valued marker, alongside a set-kinded sibling so every operator in
// the closed vocabulary has a dimension its kind admits. Finiteness matters:
// `unprovableReason` decides the DOMAIN arm first, so an infinite dimension
// would report `dimension-not-finite` for every operator and the
// classification under test would never be reached.
const singleValueProbeDecls = statusOnlyDecls + `
[tags.multi]
provenance = "owned"
kind = "int"
min = 1
max = 2
required = true

[tags.bag]
provenance = "owned"
kind = "set"
elements = ["p", "q"]
required = true
`

// singleValueProbeBody builds a model whose one rule carries a single guard
// atom spelling operator over a non-single-valued dimension.
func singleValueProbeBody(t *testing.T, operator string) string {
	t.Helper()

	key, literal := "multi", ""
	switch operator {
	case "eq":
		literal = "1"
	case "in":
		literal = "[1]"
	case "lt", "lte", "gt", "gte":
		literal = "2"
	case "exists":
		literal = "true"
	case "contains":
		key, literal = "bag", `["p"]`
	default:
		t.Fatalf("operator %q is in the closed vocabulary but this fixture "+
			"spells no literal for it; add one rather than skipping the "+
			"operator, or the invariant stops covering it", operator)
	}

	return fmt.Sprintf(`
terminal = ["done"]

[initial]
status = "a"
multi = 1
bag = ["p"]

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "probe"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.%s]
%s = %s
[rule.write]
status = "b"
`, key, operator, literal)
}

// REQ-58: "An atom over a tag not declared single-valued has no projection
// and MUST take `graph-unprovable-coverage` (`0003::A21`)."
// REQ-80 (tag-not-single-valued arm).
//
// The stable reason `tag-not-single-valued` fires for exactly the operators
// guard classifies as narrowing the tag's single held value. Lint must READ
// that classification, never carry its own copy: a copy can drift away from
// the contract the reason's remedy is stated against, which is how an `in`
// atom over a declared finite domain once reported `dimension-not-finite`
// and told the author to declare a domain they had already declared.
//
// Iterating `guard.Operators()` — the closed vocabulary itself — is what
// makes this an invariant: an operator added to guard and left unclassified
// by lint fails here without anyone remembering to extend the table.
// DOMAIN EDGE
func TestReq58_SingleValueClassificationHasOneSource(t *testing.T) {
	operators := guard.Operators()
	if len(operators) == 0 {
		t.Fatalf("the closed operator vocabulary is empty, so the loop " +
			"below asserts nothing")
	}

	for _, operator := range operators {
		t.Run(operator, func(t *testing.T) {
			r := lint(t, singleValueProbeDecls,
				singleValueProbeBody(t, operator))

			var named bool
			for _, got := range withCode(r, graphlint.CodeUnprovableCoverage) {
				if got.Reason == graphlint.ReasonTagNotSingleValued {
					named = true
				}
			}

			// The two must agree in BOTH directions. Asserting only the
			// positive arm would stay green if lint classified every
			// operator as single-value; asserting only the negative arm
			// would stay green if it classified none.
			if want := guard.SingleValueOperator(operator); named != want {
				t.Errorf("an atom spelling %q over a finite dimension "+
					"lacking the single-valued marker reports reason %q: "+
					"%t; guard.SingleValueOperator(%q) = %t. Lint's "+
					"classification must BE guard's, not a copy of it; "+
					"report:%s",
					operator, graphlint.ReasonTagNotSingleValued, named,
					operator, want, render(r))
			}
		})
	}
}
