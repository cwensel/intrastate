package cli

// RDR 0021 — determinism (`0021:C3`) and the golden pins (`0021:S1`,
// `0021:S2`).
//
// Every oracle here is BYTE equality. C3 promises byte-for-byte identity
// for one (model, build) across invocations in every `--emit`×`--as` cell,
// and the MVV asserts the replay form as byte equality rather than as an
// exit code — "a green exit code or \"did not error\" is NOT sufficient;
// fidelity loss hides behind a passing run".

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// emitCells is the four `--as`×`--emit` combinations C3 quantifies over.
var emitCells = []struct {
	name string
	args []string
}{
	{"text/json", []string{}},
	{"text/dot", []string{"--emit", "dot"}},
	{"json/json", []string{"--as=json"}},
	{"json/dot", []string{"--as=json", "--emit", "dot"}},
}

// REQ-42: "For one model input and one build, emission is byte-for-byte
// identical across invocations, in every `--emit` and `--as` combination"
// REQ-66: "`export ∘ export = byte identity on any loadable model` (C3's
// replay form; the MVV asserts it as byte equality, not exit-code green)."
// HAPPY PATH
func TestReq42And66_EmissionIsByteIdenticalAcrossInvocations(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	for _, cell := range emitCells {
		t.Run(cell.name, func(t *testing.T) {
			args := append([]string{"--model", path}, cell.args...)

			first, _, err := runGraph(t, args...)
			if err != nil {
				t.Fatalf("first emission refused: %v", err)
			}
			second, _, err := runGraph(t, args...)
			if err != nil {
				t.Fatalf("second emission refused: %v", err)
			}

			if first != second {
				t.Errorf("two emissions of one model differ in the %s cell; "+
					"C3 promises BYTE-for-byte identity across invocations "+
					"for one (model, build)\n--- first ---\n%s\n"+
					"--- second ---\n%s", cell.name, first, second)
			}
			if strings.TrimSpace(first) == "" {
				t.Errorf("the %s cell emitted nothing; byte identity over an "+
					"empty stream is a vacuous pass", cell.name)
			}
		})
	}
}

// REQ-73: "Determinism under map-seed variation — emit the same fixture
// repeatedly across `--emit`×`--as` cells, with map order PROVOKED under
// `GODEBUG=randmapiter=1` … Expected: byte-identical stdout per cell (C3);
// no cell depends on Go map iteration order."
// REQ-43: "no Go map iteration reaches the wire"
// ADVERSARIAL — `randmapiter=1` is the REQUIRED mechanism, not an option:
// without it a small fixture can pass by luck.
//
// The oracle runs in a SUBPROCESS because GODEBUG's map-iteration knob is
// read at process start; setting it in-process would not perturb this
// test binary's own map order.
func TestReq43And73_EmissionIsStableUnderProvokedMapOrder(t *testing.T) {
	requireGraphVerb(t)

	bin := buildGraphBinary(t)
	dir := t.TempDir()
	path := writeModelAt(t, dir, "model.toml", legalModel)

	for _, cell := range emitCells {
		t.Run(cell.name, func(t *testing.T) {
			args := append([]string{graphVerb, "--model", path}, cell.args...)

			var baseline string
			for i := range 8 {
				out := runBinary(t, bin, args, "GODEBUG=randmapiter=1")
				if i == 0 {
					baseline = out
					continue
				}
				if out != baseline {
					t.Fatalf("emission %d differs from the first under "+
						"GODEBUG=randmapiter=1 in the %s cell; a sequence "+
						"on the wire depends on Go map iteration order, "+
						"which C3 forbids — every sequence is PRE-SORTED by "+
						"C2's orders before marshaling\n--- first ---\n%s\n"+
						"--- run %d ---\n%s", i, cell.name, baseline, i, out)
				}
			}
			if strings.TrimSpace(baseline) == "" {
				t.Errorf("the %s cell emitted nothing under provoked map "+
					"order; the assertion would be vacuous", cell.name)
			}
		})
	}
}

// REQ-74: "Golden fixtures pin `intrastate.graph/1` for a state-machine
// model and a decision-table model. The goldens capture the DOCUMENT, not
// the envelope — `--as=text`, or `jq .data` off `--as=json` — so
// `0029:C1`'s `schema_version` … never enters the pinned bytes."
// REQ-75: "the goldens hold; a field added without a schema decision fails
// the pin (C2's additive rule tripwire), and an envelope version bump does
// not."
// HAPPY PATH
//
// The golden is MINTED by the real exporter on first run and committed
// thereafter — the RDR's own rule: the Phase 1 golden is the byte fixture
// for Scenario 2, and the spikes' improvised spellings are explicitly NOT
// promoted (F2/F3). Regenerate deliberately with -update when a schema
// decision is made; a diff names the moved member.
func TestReq74And75_GoldensPinTheDocumentNotTheEnvelope(t *testing.T) {
	requireGraphVerb(t)

	for _, tc := range []struct {
		name   string
		src    string
		golden string
	}{
		{"state machine", legalModel, "graph_0021_state_machine.json"},
		{"decision table", decisionTableFixture, "graph_0021_decision_table.json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := exportDocument(t, tc.src)

			// The pinned bytes are the DOCUMENT. The envelope's own version
			// marker moves on its own minor schedule for reasons unrelated
			// to this document, so it must not enter the pin.
			if strings.Contains(body, "schema_version") {
				t.Fatalf("the pinned document carries the ENVELOPE's "+
					"`schema_version`; S2 pins the DOCUMENT so an envelope "+
					"version bump does not fail the golden\n%s", body)
			}
			assertGolden(t, tc.golden, body)
		})
	}
}

// REQ-44: "The promise is scoped to (model, build), not across builds"
// REQ-45: "a build bump may change the JSON only by C2's additive rule and
// may re-baseline DOT diffs freely, since DOT styling is non-normative."
// BOUNDARY — the narrowing is normative, so no test may assert a
// CROSS-BUILD byte promise. The observable is that the promise holds
// within one build, which the golden above pins, and that the record makes
// no digest claim across builds (REQ-106 covers the digest half).
//
// This test asserts the scope is REAL: two invocations of the SAME build
// agree, which is the whole of what C3 promises.
func TestReq44And45_TheBytePromiseIsScopedToOneModelAndOneBuild(t *testing.T) {
	requireGraphVerb(t)

	bin := buildGraphBinary(t)
	dir := t.TempDir()
	path := writeModelAt(t, dir, "model.toml", legalModel)
	args := []string{graphVerb, "--model", path}

	first := runBinary(t, bin, args)
	second := runBinary(t, bin, args)
	if first != second {
		t.Errorf("one build emitted two different documents for one model; "+
			"C3's promise is exactly (model, build)\n--- first ---\n%s\n"+
			"--- second ---\n%s", first, second)
	}
}

// REQ-47: "a node is its canonical node key (`reach.go::(Node).key` —
// injective by escaping, ⇒ two exported nodes never collide); a row is RDR
// 0002's row identity; an edge is the `(from, to, rule)` triple"
// ADVERSARIAL — the injectivity claim has a direct observable: no two
// exported nodes share an id, and no two edges share the whole triple.
func TestReq47_ExportedIdentitiesNeverCollide(t *testing.T) {
	requireGraphVerb(t)

	// Tag names and values carrying the structural delimiters the node key
	// joins on (`;`, `=`, `,`). Unescaped, two different owned-states would
	// render identically — the collision the escaping exists to prevent.
	const delimiters = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "collide"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a,b", "a;b", "a=b"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["status"]
timeout = "2s"
read_back = true

[initial]
status = "a,b"

[context.done]
[context.done.match.status]
eq = "a=b"

[[rule]]
id = "to-semi"
[rule.match.status]
eq = "a,b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a;b"

[[rule]]
id = "to-eq"
[rule.match.status]
eq = "a;b"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "a=b"

[[rule]]
id = "go-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "go"
`
	doc := decodeDocument(t, exportDocument(t, delimiters))
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	seenNode := map[string]bool{}
	for _, n := range doc.Reach.Nodes {
		if seenNode[n.ID] {
			t.Errorf("two exported nodes share the id %q; the node key is "+
				"INJECTIVE by escaping, so two nodes standing for different "+
				"owned-states never render identically (D-identity)", n.ID)
		}
		seenNode[n.ID] = true
	}

	seenEdge := map[string]bool{}
	for _, e := range doc.Reach.Edges {
		key := e.From + "\x00" + e.To + "\x00" + e.Rule
		if seenEdge[key] {
			t.Errorf("the edge triple (%q, %q, %q) is exported twice; an "+
				"edge IS that triple (D-identity)", e.From, e.To, e.Rule)
		}
		seenEdge[key] = true
	}

	seenRow := map[string]bool{}
	for _, r := range doc.Rows {
		if seenRow[r.Identity] {
			t.Errorf("two exported rows share the identity %q; a row is RDR "+
				"0002's row identity (D-identity)", r.Identity)
		}
		seenRow[r.Identity] = true
	}
}

// --- subprocess + golden helpers ----------------------------------------

// buildGraphBinary compiles the CLI once per test and returns its path.
// The GODEBUG map-iteration knob is read at process start, so the only
// honest way to provoke map order is a real subprocess.
func buildGraphBinary(t *testing.T) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "intrastate-under-test")
	cmd := exec.Command("go", "build", "-o", bin,
		"github.com/cwensel/intrastate/cmd/intrastate")
	cmd.Dir = repoRootFor(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building the CLI under test: %v\n%s", err, out)
	}
	return bin
}

// runBinary runs the built CLI and returns its stdout, failing on a
// non-zero exit so a refusal never masquerades as an empty document.
func runBinary(t *testing.T, bin string, args []string, env ...string) string {
	t.Helper()

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v: %v\nstderr: %s", bin, args, err, stderr.String())
	}
	return stdout.String()
}

// assertGolden compares body against the committed golden, minting it when
// absent (and when -update is passed) so the FIRST green run records the
// real exporter's bytes rather than a spike's improvised spellings.
func assertGolden(t *testing.T, name, body string) {
	t.Helper()

	path := filepath.Join(repoRootFor(t), "internal", "cli", "testdata", name)
	want, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("reading golden %s: %v", name, err)
		}
		// The exporter does not exist yet, so there is nothing honest to
		// pin. Minting a golden from a missing implementation would be the
		// tautological green the red-before-green gate exists to prevent.
		t.Fatalf("no golden at %s. S2 pins `%s` for this model over the "+
			"DOCUMENT; the Phase 1 golden is minted by the REAL exporter "+
			"(F2 forbids promoting the spike's improvised field spellings). "+
			"Once the exporter lands, write this file from its output.",
			path, schemaMarker)
	}
	if string(want) != body {
		t.Errorf("the exported document does not match the golden %s.\n"+
			"A field added without a schema decision fails this pin — that "+
			"is C2's additive-rule tripwire, and the fix is a deliberate "+
			"schema decision, not a regenerated golden.\n--- want ---\n%s\n"+
			"--- got ---\n%s", name, want, body)
	}
}
