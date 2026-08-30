package cli

// RDR 0024 — the Minimum Viable Validation (`0024:MVV`), run end to end.
//
// Four steps over ONE authored table — the seed's two-rule adversarial
// table — driven through the ROOT command the way CI drives it. Nothing
// here calls the loader or the join directly.
//
// Two properties of the battery are load-bearing and asserted rather than
// assumed:
//
//  1. Step 1 is the NEGATIVE CONTROL. It reproduces the defect: the same
//     two-rule table, undeclared, lints exit 0 with zero findings. Without
//     it, step 2's refusal proves only that SOMETHING refuses, not that the
//     declaration is what made the defect visible.
//
//  2. Step 2's bar is TWO RUNS, one category each. `0024:MVV`: "Do NOT
//     expect both in one run" — the load tier is fail-fast, so a battery
//     asserting both categories from one invocation would assert a
//     guarantee this RDR explicitly does not make.
//
// Step 3 is the round trip C4 fixes: the payload carries the AUTHORED emit
// unchanged AND the disposition the declaration assigns that authored
// value. Both halves are compared value-for-value; an exit code proves
// neither.

import (
	"reflect"
	"slices"
	"strings"
	"testing"
)

// mvvSeedTable0024 is "the seed's two-rule adversarial table": the first
// rule's emit value misspells a command (`resovle`), the second rule types
// the key `nxet`.
//
// It is a decision table over one observed dimension, so both rules are
// selectable by a discriminating `--tag`, and both defects are reachable.
const mvvSeedTable0024 = `outcomes = ["decide"]

[model]
id = "mvv-0024"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.step]
provenance = "observed"
kind = "enum"
domain = ["first", "second"]
single_valued = true
required = true

[[rule]]
id = "first-rule"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.step]
eq = "first"
[rule.emit]
next = "resovle"

[[rule]]
id = "second-rule"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.step]
eq = "second"
[rule.emit]
nxet = "stop"
`

// mvvDecl0024 is step 2's declaration: `[emit.next]`, an enum with a
// route/stop-partitioned domain.
const mvvDecl0024 = `
[emit.next]
kind = "enum"
[emit.next.domain]
route = ["resolve", "refine"]
stop = ["archive"]
`

// mvvDeclare0024 inserts the declaration ahead of the rules.
func mvvDeclare0024(t *testing.T, src string) string {
	t.Helper()

	const anchor = "[[rule]]\nid = \"first-rule\""
	if !strings.Contains(src, anchor) {
		t.Fatalf("the MVV table no longer carries %q", anchor)
	}
	return strings.Replace(src, anchor, mvvDecl0024+"\n"+anchor, 1)
}

// REQ-MVV / `0024:MVV`.
// HAPPY PATH
func TestMVV_0024_DeclaredEmitVocabulary(t *testing.T) {
	t.Run("step 1: the defect, reproduced", mvv0024Step1)
	t.Run("step 2: declaring makes both defects visible, one per run",
		mvv0024Step2)
	t.Run("step 3: the payload carries emit and its disposition",
		mvv0024Step3)
	t.Run("step 4: deleting the declaration restores today's behaviour",
		mvv0024Step4)
}

// REQ-89 / `MVV` step 1: "Re-author the seed's two-rule adversarial table —
// first rule's emit value misspells a command, second rule types the key
// `nxet` — and confirm it still lints exit 0 with zero findings (the
// defect, reproduced)."
// INPUT EDGE — the negative control.
func mvv0024Step1(t *testing.T) {
	model := writeFlowModel(t, mvvSeedTable0024)

	stdout := requireSuccess(t, "lint", "--model", model, "--as=json")

	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the lint success envelope carries no `data.findings`:\n%s",
			stdout)
	}
	if n := len(*env.Data.Findings); n != 0 {
		t.Errorf("the undeclared adversarial table lints with %d findings; "+
			"the DEFECT is that it lints exit 0 with ZERO findings: %+v",
			n, *env.Data.Findings)
	}
}

// REQ-90 / `MVV` step 2: "Add an `[emit.next]` enum declaration with a
// route/stop-partitioned domain; run `intrastate lint`. **Expected**: exit
// nonzero with ONE blocking finding … Fix that rule and re-run: the second
// defect now surfaces, likewise as one finding. Across the two runs both
// categories are observed — `emit_value_out_of_domain` for the misspelled
// value and `unknown_emit_key` for the typo'd `nxet` key. Do NOT expect
// both in one run"
// ADVERSARIAL
func mvv0024Step2(t *testing.T) {
	declared := mvvDeclare0024(t, mvvSeedTable0024)

	// Run A. Both defects are present; the pipeline is fail-fast, so ONE
	// finding comes back. The RDR fixes no precedence, so the assertion is
	// on the count and on the finding's shape, not on which category.
	runA := writeFlowModel(t, declared)
	ceA := requireRefusal(t, "model-invalid", 2,
		"lint", "--model", runA, "--as=json")
	if len(ceA.Findings) != 1 {
		t.Fatalf("run A carries %d findings; the load tier is FAIL-FAST — "+
			"ONE blocking finding per run: %+v", len(ceA.Findings), ceA.Findings)
	}

	observed := map[string]bool{ceA.Findings[0].Code: true}

	// Run B. Fix the rule run A reported and re-run: the second defect now
	// surfaces, likewise as one finding.
	fixed := declared
	switch ceA.Findings[0].Code {
	case "emit_value_out_of_domain":
		fixed = strings.Replace(fixed, `next = "resovle"`, `next = "resolve"`, 1)
	case "unknown_emit_key":
		fixed = strings.Replace(fixed, `nxet = "stop"`, `next = "archive"`, 1)
	default:
		t.Fatalf("run A refused %q; the two reachable categories are "+
			"`emit_value_out_of_domain` and `unknown_emit_key`",
			ceA.Findings[0].Code)
	}
	if fixed == declared {
		t.Fatalf("the fix for %q did not apply", ceA.Findings[0].Code)
	}

	runB := writeFlowModel(t, fixed)
	ceB := requireRefusal(t, "model-invalid", 2,
		"lint", "--model", runB, "--as=json")
	if len(ceB.Findings) != 1 {
		t.Fatalf("run B carries %d findings; likewise ONE: %+v",
			len(ceB.Findings), ceB.Findings)
	}
	observed[ceB.Findings[0].Code] = true

	// Across the two runs BOTH categories are observed.
	for _, want := range []string{
		"emit_value_out_of_domain", "unknown_emit_key",
	} {
		if !observed[want] {
			t.Errorf("across the two runs the category %q was never "+
				"observed; observed = %v", want, slices.Sorted(mapKeys(observed)))
		}
	}
	// And never both from one run.
	if len(ceA.Findings) > 1 || len(ceB.Findings) > 1 {
		t.Error("a single run reported more than one finding; `0024:MVV` " +
			"says do NOT expect both in one run")
	}

	// Every finding carries the offending block's SOURCE LINE, not `:1`.
	for name, f := range map[string]string{
		"run A": ceA.Findings[0].Locator,
		"run B": ceB.Findings[0].Locator,
	} {
		if strings.HasSuffix(f, ":1") || !strings.Contains(f, ":") {
			t.Errorf("%s's finding locator is %q; each finding carries the "+
				"offending block's SOURCE LINE, not `:1`", name, f)
		}
	}
}

// REQ-91 / `MVV` step 3: "Fix both rules; lint exits 0. Run `flow resolve`
// selecting a rule whose value is listed under `stop`. **Expected**:
// payload carries the authored `emit` unchanged and `dispositions` mapping
// the key to `stop`, positioned immediately after `emit`."
// HAPPY PATH — the round trip, asserted value-for-value.
func mvv0024Step3(t *testing.T) {
	fixed := mvvDeclare0024(t, mvvSeedTable0024)
	fixed = strings.Replace(fixed, `next = "resovle"`, `next = "resolve"`, 1)
	fixed = strings.Replace(fixed, `nxet = "stop"`, `next = "archive"`, 1)
	if strings.Contains(fixed, "resovle") || strings.Contains(fixed, "nxet") {
		t.Fatal("the step-3 fixes did not apply")
	}

	model := writeFlowModel(t, fixed)

	// "lint exits 0"
	requireSuccess(t, "lint", "--model", model, "--as=json")

	// `second-rule` authors `next = "archive"`, listed under `stop`.
	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "step=second", "--as=json")
	data := flowData(t, stdout)

	// (i) the payload carries the AUTHORED emit unchanged — value for
	// value, not "did not error".
	emit, ok := data["emit"].(map[string]any)
	if !ok {
		t.Fatalf("the payload carries no `emit` object: %v", data)
	}
	wantEmit := map[string]any{"next": "archive"}
	if !reflect.DeepEqual(emit, wantEmit) {
		t.Errorf("emit = %v; want %v — the authored block, unchanged",
			emit, wantEmit)
	}

	// (ii) `dispositions` maps the key to `stop`.
	got := dispositionsOf(t, data)
	want := map[string]string{"next": "stop"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("dispositions = %v; want %v — the declaration lists "+
			"`archive` under `stop`", got, want)
	}

	// (iii) positioned IMMEDIATELY after `emit` on the wire.
	order := wireKeyOrder(t, stdout)
	ei := slices.Index(order, "emit")
	di := slices.Index(order, "dispositions")
	if ei < 0 || di < 0 {
		t.Fatalf("the wire carries emit at %d and dispositions at %d; both "+
			"are required. order = %v", ei, di, order)
	}
	if di != ei+1 {
		t.Errorf("`dispositions` is at wire index %d and `emit` at %d; "+
			"`dispositions` sits IMMEDIATELY after `emit`. order = %v",
			di, ei, order)
	}
}

// REQ-92 / `MVV` step 4: "Delete the `[emit]` table; re-run lint and
// resolve over the original (defective) table. **Expected**: exit 0, no
// findings, payload carries `dispositions: {}` — byte-identical behavior to
// today apart from the appended empty field."
// BOUNDARY
func mvv0024Step4(t *testing.T) {
	// The ORIGINAL defective table, with no `[emit]` table at all.
	model := writeFlowModel(t, mvvSeedTable0024)

	// lint: exit 0, no findings.
	stdout := requireSuccess(t, "lint", "--model", model, "--as=json")
	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the lint success envelope carries no `data.findings`:\n%s",
			stdout)
	}
	if n := len(*env.Data.Findings); n != 0 {
		t.Errorf("deleting the `[emit]` table left %d findings; the "+
			"zero-declaration model is today's world: %+v",
			n, *env.Data.Findings)
	}

	// resolve: the payload carries `dispositions: {}` and nothing else new.
	resolved := requireSuccess(t, "flow", "resolve", "--model", model,
		"--outcome", "decide", "--tag", "step=first", "--as=json")
	if !strings.Contains(resolved, `"dispositions":{}`) {
		t.Errorf("the zero-declaration payload does not carry "+
			"`\"dispositions\":{}`; the appended empty field is the only "+
			"delta from today:\n%s", resolved)
	}

	data := flowData(t, resolved)
	if got := dispositionsOf(t, data); len(got) != 0 {
		t.Errorf("dispositions = %v; want `{}`", got)
	}
	// The defective value still answers verbatim — deleting the
	// declaration restores today's behaviour exactly.
	emit, ok := data["emit"].(map[string]any)
	if !ok {
		t.Fatalf("the payload carries no `emit` object: %v", data)
	}
	if want := map[string]any{"next": "resovle"}; !reflect.DeepEqual(emit, want) {
		t.Errorf("emit = %v; want %v — the misspelled value answers "+
			"verbatim once nothing declares the key", emit, want)
	}
}

// mapKeys yields a map's keys for a sorted failure message.
func mapKeys(m map[string]bool) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
