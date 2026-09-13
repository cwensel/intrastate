package cli

// RDR 0021 — the contract surfaces Phase 4 lands (`0021:IP` Phase 4, C2's
// soundness sentence).
//
// These are documentation obligations with real observables: a generated
// reference derived from the live command tree, a hand-written wire
// contract, and the schema documentation that states the soundness rule. A
// reader who believes the diagram certifies the model is the failure this
// documentation exists to prevent, so each is asserted rather than assumed.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// readRepoDoc reads one documentation file from the repository root.
func readRepoDoc(t *testing.T, rel ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{repoRootFor(t)}, rel...)...)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", filepath.Join(rel...), err)
	}
	return string(body)
}

// REQ-103: "`docs/cli-output-contract.md` gains the document section;
// `--help-all` gains the verb; the schema docs carry C2's `abstraction`
// marker and its soundness sentence."
// HAPPY PATH — the first of three separate observables.
func TestReq103_TheOutputContractGainsTheDocumentSection(t *testing.T) {
	body := readRepoDoc(t, "docs", "cli-output-contract.md")

	if !strings.Contains(body, schemaMarker) {
		t.Errorf("docs/cli-output-contract.md does not name %q. IP Phase 4 "+
			"has it gain the document section: the wire contract is where a "+
			"consumer learns what this document is and how it is versioned",
			schemaMarker)
	}
	if !strings.Contains(body, graphVerb) {
		t.Errorf("docs/cli-output-contract.md does not mention the `%s` "+
			"verb; the document section describes the surface that emits it",
			graphVerb)
	}
}

// REQ-103 (the `--help-all` half).
// HAPPY PATH — the verb carries an extended body, so `--help-all` describes
// it rather than listing flags alone. The repo's existing help contract
// already requires every command to carry one; this pins it for this verb.
func TestReq103_HelpAllGainsTheVerb(t *testing.T) {
	requireGraphVerb(t)

	stdout, _ := runHelp(t, graphVerb, "--help-all")
	if strings.TrimSpace(stdout) == "" {
		t.Fatalf("`%s --help-all` printed nothing", graphVerb)
	}

	// The extended body must describe what the verb MEANS, not just its
	// flags: the document it emits and the format selector.
	for _, want := range []string{emitFlag, "dot", "json"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("`%s --help-all` does not mention %q; the extended "+
				"reference is where a caller learns the document formats "+
				"this verb selects between", graphVerb, want)
		}
	}
}

// REQ-103 (the generated-reference half).
// HAPPY PATH — `docs/cli-reference.md` and `llms.txt` are GENERATED from
// the live command tree and gated by `make docs-check`, so a new verb must
// appear in both or the gate fails.
func TestReq103_TheGeneratedReferenceCarriesTheVerb(t *testing.T) {
	requireGraphVerb(t)

	reference, llms := renderBoth(t)

	if !strings.Contains(reference, "## intrastate "+graphVerb) {
		t.Errorf("the generated CLI reference carries no `## intrastate %s` "+
			"section. The reference is derived from the live command tree "+
			"and `make docs-check` compares the committed copy against a "+
			"fresh render, so a verb added without regenerating fails CI",
			graphVerb)
	}
	if !strings.Contains(llms, "`intrastate "+graphVerb+"`") {
		t.Errorf("llms.txt carries no `intrastate %s` entry; it is generated "+
			"from the same tree", graphVerb)
	}
}

// REQ-30: "the schema docs state the soundness rule in one sentence:
// universal claims (\"no path does X\") proved over this relation hold at
// runtime; existence claims (\"some path reaches X\") may be spurious."
// REQ-29 (the marker the docs describe).
// HAPPY PATH — the documented soundness rule is what stops a consumer
// "verifying" a property the runtime lacks over an over-approximate
// relation.
func TestReq30_TheSchemaDocsStateTheSoundnessRule(t *testing.T) {
	body := readRepoDoc(t, "docs", "cli-output-contract.md")

	if !strings.Contains(body, abstractionMarker) {
		t.Errorf("the schema documentation does not carry C2's `abstraction` "+
			"marker %q. The document self-describes as the DECLARED "+
			"over-approximation, and the docs are where that contract is "+
			"stated", abstractionMarker)
	}

	// The soundness rule distinguishes the two claim classes. Its WORDING
	// is free; the distinction it draws is not.
	lowered := strings.ToLower(body)
	hasUniversal := strings.Contains(lowered, "universal")
	hasExistence := strings.Contains(lowered, "existence") ||
		strings.Contains(lowered, "existential")
	hasSpurious := strings.Contains(lowered, "spurious")

	if !hasUniversal || !hasExistence || !hasSpurious {
		t.Errorf("the schema documentation does not state the soundness "+
			"rule (universal=%v existence=%v spurious=%v). C2 requires one "+
			"sentence: universal claims proved over this relation hold at "+
			"runtime; existence claims may be spurious. Without it a "+
			"consumer reads the merged relation as the runtime one and "+
			"over-counts edges", hasUniversal, hasExistence, hasSpurious)
	}
}

// REQ-101: "Build the document-assembly component over `*table.Model` +
// `Reach`-with-edges; goldens pin `intrastate.graph/1`"
// REQ-102: "Register root `graph` with C1's arm set, `--emit`, the
// `TextLiner` ride, and the `graph-export-too-large` refusal."
// HAPPY PATH — the two build phases whose observable is that the surfaces
// they name exist and are reachable end to end. The detailed behaviour of
// each is asserted by the suites those clauses name; this pins that the
// phases landed at all.
func TestReq101And102_TheExportSurfacesAreBuiltAndReachable(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)

	// Phase 1's component: a document assembled over the model and the
	// relation, carrying the schema marker the goldens pin.
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("the export refused a clean model: %v", err)
	}
	doc := decodeDocument(t, stdout)
	if doc.Schema != schemaMarker {
		t.Errorf("schema = %q; want %q", doc.Schema, schemaMarker)
	}
	if doc.Reach == nil || len(doc.Reach.Edges) == 0 {
		t.Error("the document carries no reachability edges; Phase 1 builds " +
			"the component over `Reach`-with-edges")
	}

	// Phase 2's verb: the arm set, the flag, and the gateway ride, each
	// reachable through the registered command.
	cmd := lookupGraphCmd(t)
	for _, flag := range []string{"model", "flow", emitFlag} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("the verb registers no --%s flag; C1 fixes the arm set "+
				"(`--model`/`--flow`) and the document selector (`--emit`)",
				flag)
		}
	}
}
