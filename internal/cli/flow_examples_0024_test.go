package cli

// RDR 0024 `0024:S8` / `PH4` — the shipped examples and the authoring docs.
//
// The doc assertions use the `readRepoFile` shape peers 0009 and 0011
// already use to pin `docs/cli-output-contract.md` (`0024:S8`).

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
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

	// The oracle is the LOADED model's declared vocabulary, not the source
	// text. `[emit.plan]` and every member string also appear verbatim in
	// the rules' own `[rule.emit]` assignments, so a `strings.Contains`
	// sweep goes green on a model whose domains declare none of them —
	// right bytes, wrong place.
	m := loadExampleModel0024(t, path)

	for key, want := range map[string][]string{
		"plan": {"basic", "pro"},
		"dpa":  {"required", "none"},
	} {
		decl, declared := m.EmitDecls[key]
		if !declared {
			t.Errorf("the pricing example declares no emit key %q; Phase 4 "+
				"declares the example's emit keys. declared = %v",
				key, declaredEmitKeys0024(m))
			continue
		}
		for _, member := range want {
			if !slices.Contains(decl.Domain, member) {
				t.Errorf("the pricing example's %q domain is %v; it must "+
					"cover %q, which its four rows actually author",
					key, decl.Domain, member)
			}
		}
		// `PH4` fixes this example on the FLAT spelling, so no key carries
		// dispositions. The partitioned form's home is the routing example.
		if decl.Dispositions != nil {
			t.Errorf("the pricing example's %q declaration carries "+
				"dispositions %v; `PH4` fixes this example on the flat "+
				"`domain = [ ... ]` spelling", key, decl.Dispositions)
		}
	}

	// The Illustrative Code block shows `[emit.dpa]` with `domain =
	// ["required", "waived"]` — a SHAPE illustration whose members are not
	// this model's, and copying it refuses two of the four rows. This one
	// stays a source scan on purpose: it is an anti-COPY assertion, and a
	// stray `waived` anywhere in the file is what it is about.
	if strings.Contains(src, "waived") {
		t.Error("the pricing example's `dpa` domain names `waived`; that " +
			"member is the ILLUSTRATIVE CODE's shape, not this model's, " +
			"and copying it refuses two of the four rows")
	}

	// The proof, not the reading: it lints exit 0 with the declarations in.
	requireSuccess(t, "lint", "--model", path, "--as=json")
}

// loadExampleModel0024 loads a checked-in example through the loader the
// CLI itself uses, so an assertion lands on the DECLARED vocabulary rather
// than on bytes that happen to be present somewhere in the document.
func loadExampleModel0024(t *testing.T, path string) *table.Model {
	t.Helper()

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	m, err := table.Load(body, path)
	if err != nil {
		t.Fatalf("load %s: %v", path, err)
	}
	return m
}

// REQ-86 second leg: `PH4` fixes the pricing example on the FLAT
// spelling - "Neither `plan` nor `dpa` splits into routes and stops, so
// declaring this model demonstrates `domain = [ ... ]` and nothing else"
// - so its declarations carry no dispositions and `resolve` over it
// answers with an EMPTY `dispositions` object. The non-empty join is
// proved against the routing example, which REQ-87 adds as the
// partitioned form's home.
// BOUNDARY
func TestReq86_0024_ThePricingExampleResolvesWithEmptyDispositions(t *testing.T) {
	stdout := requireSuccess(t, "flow", "resolve",
		"--model", pricingModelPath(t), "--outcome", "decide",
		"--tag", "tier=paid", "--tag", "region=eu", "--as=json")

	data := flowData(t, stdout)
	got := dispositionsOf(t, data)
	if len(got) != 0 {
		t.Errorf("`flow resolve` over the pricing example joined %d "+
			"disposition(s); `PH4` declares both its emit keys with the flat "+
			"`domain = [ ... ]` spelling, which carries none. The partitioned "+
			"form's home is the routing example. payload = %v",
			len(got), data)
	}
}

// REQ-86 / REQ-87: the declarations are load-bearing, so `resolve` over
// the PARTITIONED example surfaces a joined `dispositions` entry rather
// than `{}`. This is the non-empty leg the flat pricing example cannot
// carry.
// HAPPY PATH
func TestReq86_0024_TheRoutingExampleResolvesWithJoinedDispositions(t *testing.T) {
	path := filepath.Join(repoRootFor(t), "models", "examples",
		"routing-decision-table.toml")

	stdout := requireSuccess(t, "flow", "resolve",
		"--model", path, "--outcome", "triage",
		"--tag", "severity=high", "--tag", "owner=assigned", "--as=json")

	data := flowData(t, stdout)
	got := dispositionsOf(t, data)
	if len(got) == 0 {
		t.Errorf("`flow resolve` over the routing example joined NO "+
			"dispositions; REQ-87 adds it precisely as the partitioned "+
			"declaration's worked example. payload = %v", data)
	}

	// An entry exists only when the authored value sits under a
	// disposition, so no surfaced token may be empty.
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

	// The discriminator is the LOADED declaration, not three independent
	// substring hits: `[emit.`, `route = [` and `stop = [` can all be
	// present in a comment, or under three different keys, in a model that
	// partitions nothing.
	var routing []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		if e.Name() == "pricing-decision-table.toml" {
			continue
		}
		if partitionedEmitKey0024(t, filepath.Join(dir, e.Name())) != "" {
			routing = append(routing, e.Name())
		}
	}

	if len(routing) != 1 {
		t.Fatalf("%d examples under models/examples/ declare ONE emit key "+
			"partitioned into `route` and `stop`; Phase 4 adds exactly one. "+
			"matched = %v", len(routing), routing)
	}

	// "a small routing table whose ONE emit key partitions" — the key that
	// partitions is the model's only one. A model declaring a second, flat
	// key alongside it is a different example, and the count is the only
	// thing that says so.
	m := loadExampleModel0024(t, filepath.Join(dir, routing[0]))
	key := partitionedEmitKey0024(t, filepath.Join(dir, routing[0]))
	if len(m.EmitDecls) != 1 {
		t.Fatalf("the routing example declares %d emit keys (%v); `PH4` "+
			"adds a table whose ONE emit key partitions into `route` and "+
			"`stop`", len(m.EmitDecls), declaredEmitKeys0024(m))
	}

	// Its members are exactly the union of the two dispositions, and every
	// one of them carries a disposition — the property that makes the
	// example the partitioned form's worked home.
	decl := m.EmitDecls[key]
	if got := dispositionTokens0024(decl); !slices.Equal(
		got, []string{"route", "stop"}) {
		t.Errorf("the routing example partitions into %v; `PH4` names the "+
			"two disposition tokens `route` and `stop`", got)
	}
	for _, member := range decl.Domain {
		if _, ok := decl.Dispositions[member]; !ok {
			t.Errorf("the routing example's domain member %q sits under no "+
				"disposition; the key's domain is the UNION of the "+
				"partition's member arrays", member)
		}
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

// partitionedEmitKey0024 returns the single emit key the model at path
// declares with a PARTITIONED domain, or the empty string where the model
// declares none. More than one is reported as a failure: `PH4` adds one
// worked example of the partitioned form, not a family.
func partitionedEmitKey0024(t *testing.T, path string) string {
	t.Helper()

	var found []string
	for key, decl := range loadExampleModel0024(t, path).EmitDecls {
		if decl.Dispositions != nil {
			found = append(found, key)
		}
	}
	switch len(found) {
	case 0:
		return ""
	case 1:
		return found[0]
	default:
		t.Fatalf("%s declares %v with partitioned domains; the routing "+
			"example partitions ONE emit key", path, found)
		return ""
	}
}

// dispositionTokens0024 returns the distinct disposition tokens a declared
// domain's members are listed under.
func dispositionTokens0024(decl table.EmitDecl) []string {
	seen := map[string]bool{}
	tokens := make([]string, 0, 2)
	for _, token := range decl.Dispositions {
		if !seen[token] {
			seen[token] = true
			tokens = append(tokens, token)
		}
	}
	slices.Sort(tokens)
	return tokens
}

// declaredEmitKeys0024 names the emit keys a model declares, sorted, for a
// failure message that says what WAS found rather than only what was not.
func declaredEmitKeys0024(m *table.Model) []string {
	keys := make([]string, 0, len(m.EmitDecls))
	for key := range m.EmitDecls {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
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
