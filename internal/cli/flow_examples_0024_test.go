package cli

// RDR 0024 `0024:S8` / `PH4` — the shipped examples and the authoring docs.
//
// The doc assertions use the `readRepoFile` shape peers 0009 and 0011
// already use to pin `docs/cli-output-contract.md` (`0024:S8`).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// REQ-85 / `PH4`: "Declare the pricing example's emit keys
// (`models/examples/pricing-decision-table.toml`), and extend
// `docs/cli-output-contract.md`'s resolve section with `dispositions` and
// the declaration grammar"
// REQ-86: "**Declare the example against the values it actually authors**,
// not against the Illustrative Code's shapes: the model emits `plan` ∈
// {`basic`, `pro`} and `dpa` ∈ {`required`, `none`} across its four rows."
// REQ-107 / `0024:S8`: "**Expected**:
// `models/examples/pricing-decision-table.toml` lints exit 0 with its emit
// keys declared — the declared domains must cover the values the example
// actually authors"
// HAPPY PATH
func TestReq86_0024_ThePricingExampleDeclaresTheValuesItActuallyAuthors(t *testing.T) {
	path := pricingModelPath(t)
	src := readRepoFile(t, "models/examples/pricing-decision-table.toml")

	for _, key := range []string{"[emit.plan]", "[emit.dpa]"} {
		if !strings.Contains(src, key) {
			t.Errorf("the pricing example declares no `%s`; Phase 4 declares "+
				"the example's emit keys", key)
		}
	}

	// The Illustrative Code block shows `[emit.dpa]` with `domain =
	// ["required", "waived"]` — a SHAPE illustration whose members are not
	// this model's, and copying it refuses two of the four rows.
	if strings.Contains(src, "waived") {
		t.Error("the pricing example's `dpa` domain names `waived`; that " +
			"member is the ILLUSTRATIVE CODE's shape, not this model's, " +
			"and copying it refuses two of the four rows")
	}

	// The declared domains cover the values the four rows actually author.
	for _, member := range []string{
		`"basic"`, `"pro"`, `"required"`, `"none"`,
	} {
		if !strings.Contains(src, member) {
			t.Errorf("the pricing example does not carry the member %s it "+
				"authors", member)
		}
	}

	// The proof, not the reading: it lints exit 0 with the declarations in.
	requireSuccess(t, "lint", "--model", path, "--as=json")
}

// REQ-86 second leg: the declarations are load-bearing, so `resolve` over
// the example surfaces a joined `dispositions` entry rather than `{}`.
// HAPPY PATH
func TestReq86_0024_ThePricingExampleResolvesWithJoinedDispositions(t *testing.T) {
	stdout := requireSuccess(t, "flow", "resolve",
		"--model", pricingModelPath(t), "--outcome", "decide",
		"--tag", "tier=paid", "--tag", "region=eu", "--as=json")

	data := flowData(t, stdout)
	got := dispositionsOf(t, data)
	if len(got) == 0 {
		t.Errorf("`flow resolve` over the pricing example joined NO "+
			"dispositions; Phase 4 declares its emit keys against the "+
			"values it authors. payload = %v", data)
	}
	// The selected row `paid-eu` authors `plan = "pro"` and `dpa =
	// "required"`; whichever key the example partitions, the surfaced
	// token must be one the declaration lists.
	for key, token := range got {
		if token == "" {
			t.Errorf("dispositions[%q] is the empty string; an entry exists "+
				"only when the declaration lists the authored value under a "+
				"disposition", key)
		}
	}
}

// REQ-87 / `PH4`: "this phase adds a second example — a small routing
// table whose one emit key partitions into `route` and `stop` members —
// under `models/examples/`, and `docs/cli-output-contract.md` documents
// `dispositions` against it."
// `0024:S8`: "Scenario 8 asserts both examples lint clean"
// HAPPY PATH
func TestReq87_0024_ASecondRoutingExampleShipsAndLintsClean(t *testing.T) {
	dir := filepath.Join(repoRootFor(t), "models", "examples")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read models/examples: %v", err)
	}

	var routing []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		if e.Name() == "pricing-decision-table.toml" {
			continue
		}
		body, rerr := os.ReadFile(filepath.Join(dir, e.Name()))
		if rerr != nil {
			t.Fatalf("read %s: %v", e.Name(), rerr)
		}
		src := string(body)
		// The second example's discriminator: ONE emit key partitioned
		// into `route` and `stop` members.
		if strings.Contains(src, "[emit.") &&
			strings.Contains(src, "route = [") &&
			strings.Contains(src, "stop = [") {
			routing = append(routing, e.Name())
		}
	}

	if len(routing) == 0 {
		t.Fatalf("no example under models/examples/ declares one emit key "+
			"partitioned into `route` and `stop` members; Phase 4 adds a "+
			"second example. present = %v", entries)
	}

	// Both examples lint clean.
	requireSuccess(t, "lint", "--model", pricingModelPath(t), "--as=json")
	for _, name := range routing {
		t.Run(name, func(t *testing.T) {
			requireSuccess(t, "lint",
				"--model", filepath.Join(dir, name), "--as=json")
		})
	}
}

// REQ-85 tail / REQ-107 tail: "`docs/cli-output-contract.md` documents
// `dispositions`, asserted with the `readRepoFile` shape peers 0009 and
// 0011 already use to pin that same doc."
// DOMAIN EDGE
func TestReq107_0024_TheOutputContractDocumentsDispositions(t *testing.T) {
	doc := readRepoFile(t, "docs/cli-output-contract.md")

	if !strings.Contains(doc, "dispositions") {
		t.Fatal("`docs/cli-output-contract.md` does not mention " +
			"`dispositions`; Phase 4 extends its resolve section")
	}

	// The declaration grammar is documented alongside, not just the field
	// name in a key list.
	for _, want := range []string{"[emit.", "kind"} {
		if !strings.Contains(doc, want) {
			t.Errorf("`docs/cli-output-contract.md` does not document %q; "+
				"the resolve section carries `dispositions` AND the "+
				"declaration grammar", want)
		}
	}

	// The plan-half key list is a normative fixture in that doc and must
	// carry the new field.
	if !strings.Contains(doc, "`dispositions`") {
		t.Error("`docs/cli-output-contract.md` never names `dispositions` " +
			"in backticks as a payload field")
	}
}

// REQ-88 / `PH4`: "**Amend the two stale `docs/model-authoring.md`
// sentences** in the same change … The doc says emit keys are keys
// \"nothing declares\" (twice …), and that `plan = \"<clear>\"` in an emit
// block \"loads and answers with that literal text\", which a declared
// domain excluding it now refuses. Narrow both to the same line C1 draws"
// ASSUMPTION-6: asserted by reading the narrowed text back — the absence of
// the falsified absolute claim — not by pinning new wording.
// ADVERSARIAL
func TestReq88_0024_TheAuthoringDocNoLongerClaimsNothingDeclaresEmitKeys(t *testing.T) {
	doc := readRepoFile(t, "docs/model-authoring.md")

	if strings.Contains(doc, "nothing declares them") {
		t.Error("`docs/model-authoring.md` still says emit keys are keys " +
			"`nothing declares them`; C1 makes that absolute claim false " +
			"and Phase 4 narrows it to the line C1 draws")
	}
	if strings.Contains(doc, "loads and answers with that literal text") {
		t.Error("`docs/model-authoring.md` still says an emit block's " +
			"`plan = \"<clear>\"` `loads and answers with that literal " +
			"text`; a declared domain excluding it now REFUSES")
	}

	// The narrowed line must actually be drawn, not merely deleted: the doc
	// says emit keys are not TAG keys, and MAY be declared.
	if !strings.Contains(doc, "[emit.") {
		t.Error("`docs/model-authoring.md` never shows the `[emit.<key>]` " +
			"declaration grammar; narrowing the claim means drawing C1's " +
			"line, not deleting the sentence")
	}
}

// REQ-77 / `PH1`: "**Amend three stale comments in the same change.** …
// `internal/table/model.go::EmitValue` (\"undeclared, uninterpreted, and
// compared by exact byte equality\"),
// `internal/table/normalize.go::emitSequence` (\"they are undeclared and
// uninterpreted\"), and `internal/table/dump.go::renderEmit` (\"an emit key
// has neither: it is undeclared, so there is no kind to consult\")"
// ASSUMPTION-6: the code comments are reviewed, not asserted — but the
// FALSIFIED ABSOLUTE CLAIM is checkable, and leaving it is a documented
// defect this RDR names.
// ADVERSARIAL
func TestReq77_0024_TheThreeStaleCommentsNoLongerClaimEmitIsUndeclared(t *testing.T) {
	for rel, stale := range map[string]string{
		"internal/table/model.go":     "undeclared, uninterpreted, and compared by exact byte equality",
		"internal/table/normalize.go": "they are undeclared and uninterpreted",
		"internal/table/dump.go":      "it is undeclared, so there is no kind to consult",
	} {
		t.Run(rel, func(t *testing.T) {
			if strings.Contains(readRepoFile(t, rel), stale) {
				t.Errorf("%s still carries the stale absolute claim %q; C1 "+
					"narrows it — `no kind to consult` becomes `no TAG "+
					"kind`, and `uninterpreted` becomes `never parsed, "+
					"canonicalized, or converted`", rel, stale)
			}
		})
	}
}
