package cli

// RDR 0021 — the vocabulary tiers and their enumeration seams (`0021:C1`'s
// `--emit` set, `0021:C2`'s STABILITY clause), plus the two standing
// constraints this record places on its own tests and on the verb's help
// text.
//
// `0029:C4` obliges an ENUMERATION SEAM of every tier assignment: an
// exported accessor returning the vocabulary's members, in the package that
// owns them. C1 sites the `--emit` seam in `internal/cli`; C2 sites the
// document's seam in the Go struct's json tags, where the compiler is the
// seam and the document's own decode is the by-value assertion.

import (
	"encoding/json"
	"errors"
	"go/ast"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-10: "Its enumeration seam, which `0029:C4` obliges of every tier
// assignment, is a new exported `internal/cli` accessor over the two
// members — the package owns the set and nothing it imports imports it
// back, so the accessor builds with no new import and no cycle."
// HAPPY PATH — the accessor does not exist yet, so this probe reads the
// package SOURCE rather than naming the identifier (see
// graph_probe_0021_test.go).
func TestReq10_TheEmitVocabularyHasAnExportedEnumerationSeam(t *testing.T) {
	funcs := exportedFuncs(t, graphPkgDir(t, "internal", "cli"))

	// The Go spelling is not fixed by the record, so the probe accepts the
	// names a reasonable implementation would choose and fails when NONE is
	// declared. Pinning one spelling would make this a naming test rather
	// than a seam test.
	candidates := []string{
		"EmitFormats", "EmitVocabulary", "GraphEmitFormats", "Emits",
		"DocumentFormats",
	}
	var found string
	var decl *ast.FuncDecl
	for _, name := range candidates {
		if fn, ok := funcs[name]; ok {
			found, decl = name, fn
			break
		}
	}
	if found == "" {
		t.Fatalf("internal/cli declares no exported `--%s` enumeration "+
			"accessor (looked for %v). `0029:C4` obliges a seam of every "+
			"tier assignment, and C1 sites this one in `internal/cli`: the "+
			"package owns the set and nothing it imports imports it back, "+
			"so the accessor builds with no new import and no cycle. "+
			"Without it the `append-only` tier binds only prose — there is "+
			"nothing to check membership over", emitFlag, candidates)
	}

	if results := resultTypes(t, decl); len(results) != 1 ||
		results[0] != "[]string" {
		t.Errorf("%s returns %v; the seam is an accessor over the format "+
			"set's members", found, results)
	}
}

// REQ-9: "That format set is an `append-only` vocabulary (`0029:C2`): a
// third format MAY be added in a minor, none removed or renamed within a
// major, so a consumer MUST NOT assert its cardinality or a member's
// position."
// BOUNDARY — this binds THIS suite too. The oracle is MEMBERSHIP: both
// declared members are present, and nothing here asserts the set's size or
// any member's index, so a third format added in a minor leaves it green.
func TestReq9_TheEmitVocabularyIsAppendOnlyAndAssertedByMembership(t *testing.T) {
	requireGraphVerb(t)

	// Membership, never cardinality: each declared member is accepted by
	// the verb. A third format may join them without failing this.
	path := writeModel(t, legalModel)
	for _, member := range []string{"json", "dot"} {
		if _, _, err := runGraph(t, "--model", path, "--emit", member); err != nil {
			t.Errorf("--%s %s refused: %v; both declared members of the "+
				"append-only format vocabulary are accepted",
				emitFlag, member, err)
		}
	}

	// And the complementary half: a value OUTSIDE the set is still
	// refused, so "append-only" does not mean "anything goes".
	if _, _, err := runGraph(t, "--model", path, "--emit", "xml"); err == nil {
		t.Errorf("--%s xml was accepted; an append-only vocabulary still "+
			"refuses a non-member — the tier governs how the set GROWS, "+
			"not whether it is checked", emitFlag)
	}
}

// REQ-36: "The enumeration seam `0029:C4` obliges is the Go struct's json
// tags, as it is for the `version` payload field names — the compiler is
// the seam, so the document's own decode is the by-value assertion."
// HAPPY PATH — the decode IS the assertion: a document that decodes
// value-for-value into a tagged view proves the json tags carry C2's
// spellings, without any test restating the field list as data.
func TestReq36_TheDocumentsOwnDecodeIsTheByValueAssertion(t *testing.T) {
	body := exportDocument(t, legalModel)

	// Decoding through the tagged view and finding the named members
	// populated is the by-value assertion the compiler underwrites. If a
	// json tag were misspelled, the corresponding field would decode to
	// its zero value while the raw document still parsed as JSON.
	doc := decodeDocument(t, body)

	for _, tc := range []struct {
		member string
		ok     bool
	}{
		{"schema", doc.Schema != ""},
		{"model", doc.Model != ""},
		{"class", doc.Class != ""},
		{"tags", len(doc.Tags) > 0},
		{"rows", len(doc.Rows) > 0},
		{"groups", len(doc.Groups) > 0},
		{"reach", doc.Reach != nil},
	} {
		if !tc.ok {
			t.Errorf("the %q member decoded to its zero value through the "+
				"tagged view; the json tags ARE the enumeration seam, so a "+
				"misspelled tag surfaces exactly here — the document parses "+
				"but the named member does not arrive", tc.member)
		}
	}

	// The nested views decode too, which is the same assertion one level
	// down.
	if doc.Reach != nil {
		if len(doc.Reach.Nodes) > 0 && doc.Reach.Nodes[0].ID == "" {
			t.Error("a `reach.nodes[].id` decoded empty through the tagged view")
		}
		if len(doc.Reach.Edges) > 0 && doc.Reach.Edges[0].Rule == "" {
			t.Error("a `reach.edges[].rule` decoded empty through the tagged view")
		}
	}
}

// REQ-67: "No `load ∘ export` inverse is claimed: the document is derived
// output, never a model source."
// REQ-84 (the RT3 half of Scenario 9: "no `load ∘ export` inverse is
// exercised").
// BOUNDARY — a declared EXEMPTION. The honest observable is the negative:
// the exported document is not a loadable model, and nothing in this record
// round-trips it back through the loader.
func TestReq67_NoLoadAfterExportInverseIsClaimed(t *testing.T) {
	body := exportDocument(t, legalModel)

	// The document is JSON; the loader takes TOML. Feeding the document
	// back to the loader must not be how anything here works, and the
	// exemption is real rather than incidental: the document carries no
	// authored-model grammar.
	var probe map[string]json.RawMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &probe); err != nil {
		t.Fatalf("the document is not JSON: %v", err)
	}
	for _, authored := range []string{"outcomes", "rule", "read", "write", "context"} {
		if _, has := probe[authored]; has {
			t.Errorf("the exported document carries the authored-model key "+
				"%q at top level; RT3 claims NO `load ∘ export` inverse — "+
				"the document is DERIVED output, never a model source, and "+
				"a document shaped like an authored model invites exactly "+
				"the round trip this record declines", authored)
		}
	}

	// And no test in this record may exercise the inverse.
	for _, path := range graphTestFiles(t) {
		src := readTestSource(t, path)
		if strings.Contains(src, "table.Load([]byte(body") ||
			strings.Contains(src, "table.Load([]byte(stdout") ||
			strings.Contains(src, "table.Load([]byte(doc") {
			t.Errorf("%s feeds an EXPORTED document back to the loader; "+
				"RT3 declares no `load ∘ export` inverse, so no test may "+
				"exercise one", filepath.Base(path))
		}
	}
}

// REQ-85: "The MVV is the acceptance floor; these scenarios are what
// \"done\" adds beyond it. Every oracle is byte- or set-equality — none
// asserts an exit code alone."
// ADVERSARIAL — a standing constraint on every test this record authorizes,
// so the subject under test is this record's OWN suite. A test whose whole
// oracle is an exit code would pass against an implementation emitting the
// wrong bytes.
func TestReq85_NoTestInThisRecordAssertsAnExitCodeAlone(t *testing.T) {
	for _, path := range graphTestFiles(t) {
		name := filepath.Base(path)
		src := readTestSource(t, path)

		// Each test function must carry at least one CONTENT oracle beside
		// any exit-code assertion: a byte comparison, a set comparison, or
		// a value assertion. The scan is per function body.
		for _, fn := range splitTestFuncs(src) {
			if fn.name == "" || !strings.HasPrefix(fn.name, "Test") {
				continue
			}
			hasExit := strings.Contains(fn.body, "assertExitTwo") ||
				strings.Contains(fn.body, "ExitCodeFor")
			if !hasExit {
				continue
			}
			hasContent := strings.Contains(fn.body, "assertRefusalCode") ||
				strings.Contains(fn.body, "assertNoDocumentOnStdout") ||
				strings.Contains(fn.body, "assertSameSet") ||
				strings.Contains(fn.body, "assertGolden") ||
				strings.Contains(fn.body, "decodeDocument") ||
				strings.Contains(fn.body, "decodeGeneric") ||
				strings.Contains(fn.body, "ce.Code") ||
				strings.Contains(fn.body, "ce.Param")
			if !hasContent {
				t.Errorf("%s::%s asserts an exit code with no byte- or "+
					"set-equality oracle beside it. The Testing Strategy's "+
					"preamble is explicit: every oracle is byte- or "+
					"set-equality, and none asserts an exit code alone — an "+
					"exit-only test passes against an implementation "+
					"emitting the wrong bytes", name, fn.name)
			}
		}
	}
}

// REQ-95: "Selection-arm refusal codes | `lint.go::runLint` | `graph` verb
// mirrors the arm set and code spellings verbatim; help text is the verb's
// own (C1) | … | lint — C1 mirrors, never redefines"
// REQ-2 (the mirroring half, from the authority table's side).
// BOUNDARY — mirroring is scoped to the ARM SET and its CODE SPELLINGS, and
// explicitly NOT to lint's help/usage text, which each verb words for
// itself. Both halves are asserted: the codes match, the help text does not.
func TestReq95_TheCodesAreMirroredButTheHelpTextIsTheVerbsOwn(t *testing.T) {
	requireGraphVerb(t)

	graphCmd := lookupGraphCmd(t)
	lintCmd := lookupCmd(t, "lint")

	// The help text is the verb's OWN: a `graph` verb describing itself as
	// lint does would tell a caller the wrong thing about what it produces.
	if graphCmd.Short == lintCmd.Short {
		t.Errorf("the `%s` verb's Short text is identical to lint's (%q); "+
			"C1 scopes mirroring to the ARM SET and its CODE SPELLINGS, "+
			"never to lint's help/usage text, which each verb words for "+
			"itself", graphVerb, lintCmd.Short)
	}
	if strings.Contains(strings.ToLower(graphCmd.Short), "invariant") {
		t.Errorf("the `%s` verb's Short text describes checking invariants "+
			"(%q); this verb exports a document and runs no lint invariant "+
			"(C4)", graphVerb, graphCmd.Short)
	}

	// The CODES are mirrored verbatim. Each arm is asserted against lint's
	// own answer for the same input, so a rename on either side breaks it.
	dir := t.TempDir()
	good := writeModelAt(t, dir, "model.toml", legalModel)
	missing := filepath.Join(dir, "absent.toml")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{"both flags", []string{"--model", good, "--flow", "x"}},
		{"neither flag", nil},
		{"flow alone", []string{"--flow", "x"}},
		{"unreadable", []string{"--model", missing}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, graphErr := runGraph(t, append(append([]string{}, tc.args...), "--as=json")...)
			_, _, lintErr := runCmd(t, append(append([]string{"lint"}, tc.args...), "--as=json")...)

			graphCode := refusalCodeOf(t, graphErr)
			lintCode := refusalCodeOf(t, lintErr)
			if graphCode != lintCode {
				t.Errorf("the %s arm refuses %q on `%s` and %q on `lint`; C1 "+
					"mirrors lint's code spellings VERBATIM and never "+
					"redefines them — lint owns these codes (`0006:C19`/`C20`)",
					tc.name, graphCode, graphVerb, lintCode)
			}
		})
	}
}

// refusalCodeOf returns a refusal's structured code, or "" when the call
// succeeded.
func refusalCodeOf(t *testing.T, err error) string {
	t.Helper()

	if err == nil {
		return ""
	}
	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	return ce.Code
}

// graphTestFiles lists this record's own test files, which are the subject
// of the two self-binding constraints above.
func graphTestFiles(t *testing.T) []string {
	t.Helper()

	dir := graphPkgDir(t, "internal", "cli")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var out []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "graph_") &&
			strings.HasSuffix(e.Name(), "_0021_test.go") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	slices.Sort(out)
	if len(out) == 0 {
		t.Fatal("no RDR 0021 test files found; the self-binding constraints " +
			"have no subject")
	}
	return out
}

// readTestSource reads one test file's source.
func readTestSource(t *testing.T, path string) string {
	t.Helper()

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(body)
}

// testFunc is one test function's name and body, for the source scans.
type testFunc struct {
	name string
	body string
}

// splitTestFuncs splits a source file into its top-level function bodies.
func splitTestFuncs(src string) []testFunc {
	var out []testFunc
	var cur testFunc
	var b strings.Builder
	for line := range strings.Lines(src) {
		if strings.HasPrefix(line, "func ") {
			if cur.name != "" {
				cur.body = b.String()
				out = append(out, cur)
			}
			b.Reset()
			name := strings.TrimPrefix(line, "func ")
			if idx := strings.IndexAny(name, "("); idx > 0 {
				name = name[:idx]
			}
			cur = testFunc{name: strings.TrimSpace(name)}
		}
		b.WriteString(line)
	}
	if cur.name != "" {
		cur.body = b.String()
		out = append(out, cur)
	}
	return out
}
