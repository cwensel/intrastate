package cli

// RDR 0021 — neutrality (`0021:C4`) and the ceiling refusal.
//
// C4's oracle is MECHANISM-INDEPENDENT by construction: it compares
// `intrastate lint`'s observable output with and without the export code
// present, so it binds whichever edge mechanism the implementation picks
// (post-hoc recovery or A2's in-traversal observer). Nothing here asserts
// HOW the edges are recovered.

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-49: "It MUST NOT run the lint invariants, MUST NOT emit findings, and
// MUST NOT alter any input it shares with lint: `intrastate lint`'s
// verdict, finding set, and bytes are identical with and without the export
// code present"
// REQ-76: "Lint neutrality — run `intrastate lint` over blocking and clean
// fixtures with the export code present, byte-compared against a pre-change
// capture. Expected: identical verdict, finding set, and bytes (C4)"
// REQ-50: "the oracle is MECHANISM-INDEPENDENT"
// ADVERSARIAL — the pre-change capture is the committed golden below. A
// perturbation of lint by the export code changes those bytes.
//
// ASSUMPTION A-3: "bytes" means the full observable output — stdout AND
// stderr — plus the exit code, since a perturbation that moved an advisory
// between streams would otherwise pass a stdout-only compare.
func TestReq49And50And76_LintIsByteIdenticalWithTheExportCodePresent(t *testing.T) {
	for _, tc := range []struct {
		name   string
		src    string
		golden string
	}{
		{"clean model", legalModel, "lint_0021_neutrality_clean.txt"},
		{"advisory model", advisoryModel, "lint_0021_neutrality_advisory.txt"},
		{"blocking model", illegalModel, "lint_0021_neutrality_blocking.txt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeModel(t, tc.src)

			// Both modes, because the two render through different branches
			// of the gateway and a perturbation could touch either.
			var capture strings.Builder
			for _, mode := range []string{"--as=json", "--as=text"} {
				stdout, stderr, err := runCmd(t, "lint", "--model", path, mode)
				fmt.Fprintf(&capture, "=== lint %s ===\nexit: %d\n"+
					"--- stdout ---\n%s\n--- stderr ---\n%s\n",
					mode, clierr.ExitCodeFor(err), stdout, stderr)
			}
			assertGolden(t, tc.golden, capture.String())
		})
	}
}

// REQ-48: "The export runs load, normalization, grouping, and the
// reachability traversal only, reaching the traversal through the same
// `graphlint` entry surface lint uses"
// REQ-93: "lint's call site stays `analysis.go:46`, unchanged"
// HAPPY PATH — structural, not disciplinary: both arms reach ONE traversal,
// so verb/lint drift is impossible by construction rather than by care.
//
// The observable: the relation the export publishes equals the relation
// `graphlint` computes for the same model. A verb that ran its own second
// traversal could diverge from lint's; this asserts it does not.
func TestReq48And93_TheExportedRelationIsTheSameTraversalLintReaches(t *testing.T) {
	requireGraphVerb(t)

	m := mustLoadModelForGraph(t, legalModel)
	want := graphlint.Reach(m)

	doc := decodeDocument(t, exportDocument(t, legalModel))
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	if len(doc.Reach.Nodes) != len(want) {
		t.Errorf("the export published %d reachable nodes; the `graphlint` "+
			"traversal lint reaches computes %d for the same model. C4 "+
			"requires the export reach the traversal through the SAME entry "+
			"surface lint uses, so verb/lint drift is structural rather "+
			"than disciplinary", len(doc.Reach.Nodes), len(want))
	}
}

// REQ-49 (the no-findings half) / REQ-31 restated behaviourally.
// REQ-104: "The verb runs load + normalize + reachability only: it never
// runs the lint invariants, emits no findings, and exports any model that
// loads"
// ADVERSARIAL — a model lint REFUSES still exports, and the document
// carries no finding.
func TestReq49And104_TheExportEmitsNoFindingsEvenForARefusedModel(t *testing.T) {
	requireGraphVerb(t)

	body := exportDocument(t, illegalModel)
	doc := decodeGeneric(t, body)

	for _, forbidden := range []string{"findings", "finding", "verdict"} {
		if _, has := doc[forbidden]; has {
			t.Errorf("the exported document carries %q; C4 forbids the verb "+
				"to emit findings at all — it runs load, normalization, "+
				"grouping and the traversal ONLY", forbidden)
		}
	}
	for _, code := range graphlint.BlockingCodes() {
		if strings.Contains(body, code) {
			t.Errorf("the exported document mentions the blocking finding "+
				"code %q; the export never runs the lint invariants", code)
		}
	}
}

// REQ-51: "The verb succeeds for ANY model that loads, including a model
// lint refuses; the model-loads-but-lint-refuses case is a dedicated
// fixture with an asserted, defined document — never whatever the traversal
// happens to do"
// REQ-77: "Model that loads but lint refuses. Expected: `graph` succeeds
// with an ASSERTED document (defined content, not incidental), carrying no
// verdict or finding field (C2, C4)."
// REQ-87: "Model loads, lint would refuse | 0 | — | ASSERTED document
// (defined content) | loud — the debugging half of the outcome (C4)"
// HAPPY PATH — the document content is ASSERTED, not merely non-empty.
func TestReq51And77And87_AModelLintRefusesStillExportsAnAssertedDocument(t *testing.T) {
	requireGraphVerb(t)

	// The premise: lint really does refuse this fixture. Without it the
	// test would pass against a model lint happens to accept, asserting
	// nothing about the clause.
	path := writeModel(t, illegalModel)
	_, _, lintErr := runCmd(t, "lint", "--model", path, "--as=json")
	if lintErr == nil {
		t.Fatal("the fixture lints CLEAN; the clause under test — a model " +
			"that loads but lint refuses — is unexercised")
	}
	var ce *clierr.CLIError
	if !errors.As(lintErr, &ce) || ce.Code != graphlint.AggregateCode {
		t.Fatalf("lint refused with %v; the fixture must refuse under the "+
			"aggregate code for this to be the refused-model arm", lintErr)
	}

	// The export succeeds over the same model, at exit 0.
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("`graph` refused a model that LOADS (lint's verdict is not "+
			"the export's): %v. C4 makes the verb succeed for ANY model "+
			"that loads — a failing model must still be inspectable as a "+
			"graph, which is the debugging half of the user outcome", err)
	}

	// ASSERTED content, not "it produced something". Each of these is a
	// defined value for this fixture rather than an incidental one.
	doc := decodeDocument(t, stdout)
	if doc.Schema != schemaMarker {
		t.Errorf("schema = %q; want %q", doc.Schema, schemaMarker)
	}
	if doc.Model != "lintfix" {
		t.Errorf("model = %q; want the authored id %q", doc.Model, "lintfix")
	}
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}
	if doc.Reach.Abstraction != abstractionMarker {
		t.Errorf("reach.abstraction = %q; want %q",
			doc.Reach.Abstraction, abstractionMarker)
	}
	if len(doc.Reach.Nodes) == 0 {
		t.Error("the document publishes no reachable nodes for a model with " +
			"a declared root; the content must be DEFINED, never whatever " +
			"the traversal happens to do")
	}
	if len(doc.Rows) == 0 {
		t.Error("the document publishes no rows for a model declaring rules")
	}
}

// REQ-52: "When the traversal is incomplete under the published node
// ceiling (`graphlint.NodeCeiling()`), the verb refuses with the scalar
// code `graph-export-too-large` (GroupUserEnv, exit 2) naming the ceiling
// and the narrow-a-domain remedy — never a partial document, because a
// partial graph diffs as a graph change."
// REQ-81: "Traversal incomplete under the published node ceiling.
// Expected: refusal with `graph-export-too-large` (GroupUserEnv, exit 2)
// naming the ceiling and remedy; no document on stdout (C4)."
// REQ-88 (the disposition row).
// ADVERSARIAL
func TestReq52And81And88_AnIncompleteTraversalRefusesTooLarge(t *testing.T) {
	requireGraphVerb(t)

	src := overCeilingModel(t)
	path := writeModel(t, src)

	// PRECONDITION, not a branch: a fixture that fails to exceed the
	// ceiling leaves the clause unexercised, and a conditional body would
	// let the whole assertion be deleted with the suite still green.
	m := mustLoadModelForGraph(t, src)
	if got := len(graphlint.Reach(m)); got <= graphlint.NodeCeiling() {
		t.Fatalf("the fixture reaches %d nodes, at or under the published "+
			"ceiling of %d, so the refusal under test cannot fire; widen "+
			"the generated fixture rather than accepting a vacuous pass",
			got, graphlint.NodeCeiling())
	}

	stdout, _, err := runGraph(t, "--model", path, "--as=json")

	assertRefusalCode(t, err, tooLargeCode)
	assertExitTwo(t, err)

	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if ce.Group != clierr.GroupUserEnv {
		t.Errorf("the refusal carries group %v; C4 fixes GroupUserEnv "+
			"(exit 2)", ce.Group)
	}
	// A SCALAR refusal: it names one subject, so it carries no findings[].
	if len(ce.Findings) > 0 {
		t.Errorf("the ceiling refusal carries %d findings[]; C4 mints it as "+
			"a SCALAR code, and C4 forbids the verb to emit findings at all",
			len(ce.Findings))
	}

	// It names the ceiling and the remedy.
	text := ce.Message + " " + ce.Detail + " " + ce.Hint
	if !strings.Contains(text, fmt.Sprint(graphlint.NodeCeiling())) {
		t.Errorf("the refusal does not name the published ceiling %d: %q",
			graphlint.NodeCeiling(), text)
	}
	if !strings.Contains(strings.ToLower(text), "narrow") {
		t.Errorf("the refusal does not name the narrow-a-domain remedy: %q — "+
			"the over-ceiling debugger gets the same remedy lint gives, not "+
			"a partial graph (premortem P-13)", text)
	}

	// NEVER a partial document: a partial graph diffs as a graph change.
	assertNoDocumentOnStdout(t, stdout, "over-ceiling traversal")
}

// REQ-55: "The distinct bound this refusal does NOT involve is the
// guard-product one, `graphlint.ProductBound()` (`guard.Bound()`); no
// guard-product input reaches the export's refusal path."
// BOUNDARY — a model whose GUARD PRODUCT is over its bound, but whose
// reachable node set is small, must still export. Refusing it would mean
// the export's refusal path consulted the wrong bound.
func TestReq55_TheGuardProductBoundDoesNotReachTheExportRefusalPath(t *testing.T) {
	requireGraphVerb(t)

	// `mvvWide` declares an owned int whose finite domain far exceeds the
	// published product bound, while its owned-state lattice stays small.
	src := mvvHeader + mvvModelTable + mvvStatus + mvvWide + wideBody

	m := mustLoadModelForGraph(t, src)
	if got := len(graphlint.Reach(m)); got > graphlint.NodeCeiling() {
		t.Fatalf("the fixture reaches %d nodes, over the ceiling of %d; it "+
			"must exercise the PRODUCT bound with a small node set, or it "+
			"cannot tell the two bounds apart", got, graphlint.NodeCeiling())
	}

	if _, _, err := runGraph(t, "--model", writeModel(t, src)); err != nil {
		var ce *clierr.CLIError
		if errors.As(err, &ce) && ce.Code == tooLargeCode {
			t.Errorf("the export refused %q over a model whose GUARD PRODUCT "+
				"exceeds `ProductBound()` but whose node set is under "+
				"`NodeCeiling()`. C4 scopes this refusal to the node-count "+
				"completeness bound alone: no guard-product input reaches "+
				"the export's refusal path", tooLargeCode)
		} else {
			t.Fatalf("the export refused unexpectedly: %v", err)
		}
	}
}

// REQ-105: "selection flags → load (`table.LoadWithAdvisories`'s underlying
// load path — advisories are ignored here; they are lint's advisory
// channel, not graph data) → assemble → render (`--emit`) → respond gateway."
// BOUNDARY — an advisory-carrying model exports a document that carries no
// advisory content: advisories never become document data.
func TestReq105_AdvisoriesAreIgnoredAndNeverBecomeDocumentContent(t *testing.T) {
	requireGraphVerb(t)

	// The premise: this fixture really does provoke an advisory through
	// lint's channel.
	path := writeModel(t, advisoryModel)
	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the advisory fixture does not lint clean: %v", err)
	}
	findings := findingsFromJSON(t, stdout, false)
	if len(findings) == 0 {
		t.Fatal("the fixture carries no advisory findings; the clause under " +
			"test is unexercised")
	}

	body := exportDocument(t, advisoryModel)
	for _, f := range findings {
		if f.Code != "" && strings.Contains(body, f.Code) {
			t.Errorf("the exported document carries the advisory code %q; "+
				"advisories are lint's advisory channel, NOT graph data, "+
				"and the export ignores them", f.Code)
		}
	}
	doc := decodeGeneric(t, body)
	for _, forbidden := range []string{"advisories", "advisory", "notes", "warnings"} {
		if _, has := doc[forbidden]; has {
			t.Errorf("the exported document carries a %q member; the export "+
				"path ignores advisories", forbidden)
		}
	}
}

// --- fixtures and helpers ------------------------------------------------

// mustLoadModelForGraph loads a fixture through the production loader.
func mustLoadModelForGraph(t *testing.T, src string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), "fixture.toml")
	if err != nil {
		t.Fatalf("fixture must load clean; refused: %v", err)
	}
	return m
}

// overCeilingModel generates a model whose reachable node set exceeds the
// published ceiling: one owned enum plus N clearable owned tags, each with
// its own clearing row, so every subset of the keys is a reachable presence
// footprint (2^N). The shape is the one `internal/graphlint`'s own ceiling
// test uses.
func overCeilingModel(t *testing.T) string {
	t.Helper()

	// 13 keys ⇒ 8192 footprints against the published ceiling of 4096.
	const keys = 13

	var b strings.Builder
	b.WriteString("\noutcomes = [\"go\", \"stop\"]\nterminal = [\"done\"]\n\n" +
		"[model]\nid = \"ceiling\"\nversion = 1\n\n" +
		"[tags.recognized]\nprovenance = \"recognized\"\nkind = \"enum\"\n" +
		"single_valued = true\nrequired = true\n\n" +
		"[tags.status]\nprovenance = \"owned\"\nkind = \"enum\"\n" +
		"domain = [\"a\", \"b\"]\nsingle_valued = true\nrequired = true\n")
	for k := range keys {
		fmt.Fprintf(&b, "\n[tags.k%d]\nprovenance = \"owned\"\nkind = \"enum\"\n"+
			"domain = [\"on\"]\nsingle_valued = true\n", k)
	}

	keyList := "\"status\""
	for k := range keys {
		keyList += fmt.Sprintf(", \"k%d\"", k)
	}
	fmt.Fprintf(&b, "\n[read.own]\nrole = \"t\"\npath = \"t.own\"\n"+
		"keys = [%s]\ntimeout = \"2s\"\n\n[write.own]\nrole = \"t\"\n"+
		"path = \"t.own\"\nkeys = [%s]\ntimeout = \"2s\"\nread_back = true\n",
		keyList, keyList)

	b.WriteString("\n[initial]\nstatus = \"a\"\n")
	for k := range keys {
		fmt.Fprintf(&b, "k%d = \"on\"\n", k)
	}
	b.WriteString("\n[context.done]\n[context.done.match.status]\neq = \"b\"\n")
	for k := range keys {
		fmt.Fprintf(&b, "\n[[rule]]\nid = \"drop%d\"\nclear = [\"k%d\"]\n"+
			"[rule.match.status]\neq = \"a\"\n[rule.match.recognized]\n"+
			"eq = \"go\"\n[rule.write]\nstatus = \"a\"\n", k, k)
	}
	return b.String()
}

// wideBody is the body half of the guard-product fixture: a small
// owned-state lattice over a very wide guard dimension.
const wideBody = `
[read.own]
role = "t"
path = "t.own"
keys = ["status", "n"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status", "n"]
timeout = "2s"
read_back = true

[initial]
status = "a"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.n]
gte = 0
[rule.write]
status = "b"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
