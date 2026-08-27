package table_test

// RDR 0010 Phase 3b — adversarial coverage for the class/owned-set
// agreement check's POSITION in `run`'s fail-fast order.
//
// FAILURE MODES ANCHOR. The record's `Trade-offs / Failure Modes` section
// opens with the visible arm:
//
//	**Visible**: class/owned-set disagreement → `flow-model-invalid` with a
//	`malformed model declaration` finding naming the class and count …
//
// The category and the `owned=<n>` token are both delivered. What is not is
// the ORDER `0010:C1` fixes that visible failure against. C1 states the
// position at BOTH ends and then draws a consequence from the lower one:
//
//	The class is read in `loadModelHeader` (it is `[model]` data) and the
//	agreement is checked in a step at or after `loadTags`, so an
//	undeclared-tag refusal precedes a class-disagreement refusal under
//	`run`'s fail-fast order.
//
// That consequence is the clause under test. `internal/table/load.go::run`
// fixes the step order `loadModelHeader, loadOutcomes, loadTags,
// loadAccessors, loadDump, loadContexts, loadInitial, loadTerminal,
// normalizeRules, checkAccessorBindings`, and the undeclared-tag refusal a
// RULE mints — `rule <id> matches the undeclared tag <k>` — is raised in
// `normalizeRules`, the second-to-last step. The implementation inserted
// `checkClassAgreement` immediately after `loadTags`: the EARLIEST point in
// C1's window, and therefore SEVEN steps ahead of the refusal C1 says must
// precede it.
//
// The window C1 fixes is "at or after `loadTags` and before
// `checkAccessorBindings`". Any step in that window satisfies the window —
// but not every step in it satisfies the ordering consequence C1 draws in
// the same breath, and the implementation picked one that does not. A
// doubly-malformed model reports the class disagreement and hides the
// undeclared tag, which is the reverse of what an author is told to expect:
// they fix the class, reload, and meet a second refusal C1 promised would
// have come first.
//
// This is NOT the doubly-malformed case C1 separately settles. That one is
// "an owned tag AND `[initial]`", and its ceiling ("before
// `checkAccessorBindings`") is what makes the class refusal win — correctly,
// and the implementation delivers it. This test is the FLOOR's consequence,
// which points the other way, and the two are only compatible if the check
// sits between `normalizeRules` and `checkAccessorBindings`.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// advDoublyMalformed declares an owned tag under `class = "decision-table"`
// (a class disagreement) AND references the undeclared tag `nope` from a
// rule's match block (an undeclared-tag refusal). Both are decidable; C1
// fixes which one is reported.
const advDoublyMalformed = `
outcomes = ["go"]

[model]
id = "adv-order"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.st]
provenance = "owned"
kind = "enum"
domain = ["a"]
single_valued = true
required = true

[[rule]]
id = "r1"
[rule.match.recognized]
eq = "go"
[rule.match.nope]
eq = "x"
`

// TestAdv0010_UndeclaredTagPrecedesTheClassDisagreement pins `0010:C1`'s
// "an undeclared-tag refusal precedes a class-disagreement refusal under
// `run`'s fail-fast order".
//
// FAILURE MODE GUARDED: the agreement check placed at the FLOOR of C1's
// window rather than inside it. `checkClassAgreement` runs immediately after
// `loadTags`, so it pre-empts `normalizeRules` — where a rule's
// undeclared-tag refusal is raised — and the model reports
// `malformed_model_declaration` where C1 fixes `unknown_tag`. Both the
// window's floor and its ceiling are satisfiable at once only by a step
// between `normalizeRules` and `checkAccessorBindings`; the chosen position
// satisfies the ceiling and breaks the floor's stated consequence.
func TestAdv0010_UndeclaredTagPrecedesTheClassDisagreement(t *testing.T) {
	_, err := table.Load([]byte(advDoublyMalformed), "adv-0010.toml")
	if err == nil {
		t.Fatal("a doubly-malformed model loaded; want a refusal")
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("refusal carries no category: %v", err)
	}
	if cat != table.CatUnknownTag {
		t.Fatalf("category = %q; want %q — `0010:C1` fixes that an "+
			"undeclared-tag refusal PRECEDES a class-disagreement refusal "+
			"under `run`'s fail-fast order, and the class disagreement is "+
			"what surfaced instead: %v", cat, table.CatUnknownTag, err)
	}
}

// TestAdv0010_ClassDisagreementStillWinsOverAccessorBinding is the control
// for the clause above: C1's CEILING must keep holding once the floor is
// respected. A decision table declaring both an owned tag and `[initial]`
// must refuse on the class, never with the writer-arity diagnostic C2 calls
// "survivable but not self-explanatory".
//
// It passes today and is kept as the pin that any fix for the ordering
// defect must not trade away — moving the check LATER than
// `checkAccessorBindings` would satisfy the floor and break this.
func TestAdv0010_ClassDisagreementStillWinsOverAccessorBinding(t *testing.T) {
	src := strings.Replace(advDoublyMalformed, `[rule.match.nope]
eq = "x"
`, "", 1) + `
[initial]
st = "a"
`
	_, err := table.Load([]byte(src), "adv-0010.toml")
	if err == nil {
		t.Fatal("a decision table declaring owned state and [initial] loaded")
	}
	cat, _ := table.CategoryOf(err)
	if cat != table.CatMalformedModelDeclaration {
		t.Fatalf("category = %q; want %q — `0010:C1` fixes the window's "+
			"CEILING so the doubly-malformed case is determinate on the "+
			"class, never on the accessor-binding writer arity: %v",
			cat, table.CatMalformedModelDeclaration, err)
	}
	if !strings.Contains(err.Error(), "owned=1") {
		t.Errorf("detail = %q; want it to carry the `owned=<n>` token "+
			"`0010:C1` renders verbatim", err.Error())
	}
}
