package cli

// RDR 0021 — the Minimum Viable Validation (`0021:MVV`).
//
// Five steps, one runnable end-to-end test. The record declares Round-Trip
// / Inverse Invariants, so per the launch contract the MVV asserts the
// reconstructed value equals the original BYTE-for-byte — "a green exit
// code or \"did not error\" is NOT sufficient; fidelity loss hides behind a
// passing run".
//
// End-state the record names: one loadable model, four invocations,
// byte-level oracles for determinism, mode agreement, and lint neutrality.

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// REQ-MVV (`0021:MVV`), steps 1–5:
//  1. "Author a small state-machine fixture (two owned states, one
//     terminal, one escape row) and a decision-table fixture."
//  2. "`intrastate graph --model <fixture>` twice → the two stdouts are
//     BYTE-identical, parse as JSON, and carry every C2 field (schema,
//     model identity and class, tags, initial, terminal, rows, groups,
//     reach.nodes, reach.edges, and the `reach` abstraction marker)."
//  3. "`intrastate graph --model <fixture> --emit dot | dot -Tsvg`
//     renders; the DOT node/edge id set equals step 2's `reach` block."
//  4. "`intrastate graph --model <fixture> --as=json | jq .data` equals
//     step 2's document, value-for-value."
//  5. "`intrastate lint` over a fixture WITH blocking findings still
//     refuses identically (byte-compared against a pre-change capture),
//     while `intrastate graph` over the same model succeeds with an
//     ASSERTED document (defined content, not incidental) — the neutrality
//     oracle, mechanism-independent per C4."
//
// HAPPY PATH
func TestMVV_LintNormalizedGraphExport(t *testing.T) {
	requireGraphVerb(t)

	t.Run("step 1 — fixtures load", mvvStep1FixturesLoad)
	t.Run("step 2 — double emit is byte-identical", mvvStep2DoubleEmit)
	t.Run("step 3 — DOT renders and matches", mvvStep3DOTMatchesReach)
	t.Run("step 4 — jq .data equals the document", mvvStep4ModeAgreement)
	t.Run("step 5 — lint neutrality", mvvStep5LintNeutrality)
}

// mvvStateMachine is step 1's state-machine fixture: two owned states
// (`draft`, `final`), one terminal context, one escape row.
const mvvStateMachine = `
outcomes = ["advance"]
terminal = ["done"]

[model]
id = "mvv-machine"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["stage"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["stage"]
timeout = "2s"
read_back = true

[initial]
stage = "draft"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "lock"
[rule.match.stage]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"

[[rule]]
id = "advance-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "advance"
`

// REQ-68: "Author a small state-machine fixture (two owned states, one
// terminal, one escape row) and a decision-table fixture."
// mvvStep1FixturesLoad: "Author a small state-machine fixture (two owned
// states, one terminal, one escape row) and a decision-table fixture."
func mvvStep1FixturesLoad(t *testing.T) {
	for _, tc := range []struct {
		name string
		src  string
	}{
		{"state machine", mvvStateMachine},
		{"decision table", decisionTableFixture},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := mustLoadModelForGraph(t, tc.src)
			if m.ID == "" {
				t.Fatal("the fixture carries no authored [model] id")
			}
		})
	}

	// The state machine's declared shape is what the step names, asserted
	// so a later edit cannot quietly drop the escape row or the terminal.
	m := mustLoadModelForGraph(t, mvvStateMachine)
	if len(m.Terminal) == 0 {
		t.Error("the state-machine fixture declares no terminal predicate set")
	}
	var escapes int
	for _, r := range m.Rows {
		if len(r.Escape) > 0 {
			escapes++
		}
	}
	if escapes == 0 {
		t.Error("the state-machine fixture declares no escape row")
	}
}

// REQ-69: "`intrastate graph --model <fixture>` twice → the two stdouts are
// BYTE-identical, parse as JSON, and carry every C2 field"
// mvvStep2DoubleEmit: "`intrastate graph --model <fixture>` twice → the two
// stdouts are BYTE-identical, parse as JSON, and carry every C2 field".
func mvvStep2DoubleEmit(t *testing.T) {
	path := writeModel(t, mvvStateMachine)

	first, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("first emission refused: %v", err)
	}
	second, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("second emission refused: %v", err)
	}

	// BYTE identity, not "both succeeded".
	if first != second {
		t.Fatalf("two emissions differ; RT2 is `export ∘ export = byte "+
			"identity on any loadable model`, asserted as byte equality "+
			"rather than exit-code green\n--- first ---\n%s\n--- second ---\n%s",
			first, second)
	}

	doc := decodeDocument(t, first)

	// Every field step 2 enumerates, each named individually. This is not a
	// cardinality assertion: C2's STABILITY clause forbids asserting the
	// field set's size, so each member is checked for itself.
	if doc.Schema != schemaMarker {
		t.Errorf("schema = %q; want %q", doc.Schema, schemaMarker)
	}
	if doc.Model != "mvv-machine" {
		t.Errorf("model identity = %q; want the authored id %q",
			doc.Model, "mvv-machine")
	}
	if doc.Class == "" {
		t.Error("the document carries no model class")
	}
	if len(doc.Tags) == 0 {
		t.Error("the document carries no tags")
	}
	if len(doc.Rows) == 0 {
		t.Error("the document carries no rows")
	}
	if len(doc.Groups) == 0 {
		t.Error("the document carries no groups")
	}

	generic := decodeGeneric(t, first)
	for _, member := range []string{"initial", "terminal"} {
		if _, has := generic[member]; !has {
			t.Errorf("the document carries no %q member", member)
		}
	}

	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}
	if len(doc.Reach.Nodes) == 0 {
		t.Error("`reach.nodes` is empty for a model with a declared root")
	}
	if len(doc.Reach.Edges) == 0 {
		t.Error("`reach.edges` is empty for a model declaring an advancing rule")
	}
	if doc.Reach.Abstraction != abstractionMarker {
		t.Errorf("the `reach` abstraction marker = %q; want %q",
			doc.Reach.Abstraction, abstractionMarker)
	}
}

// REQ-70 (the MVV step 3 form of the DOT set equality).
// mvvStep3DOTMatchesReach: "`intrastate graph --model <fixture> --emit dot
// | dot -Tsvg` renders; the DOT node/edge id set equals step 2's `reach`
// block."
func mvvStep3DOTMatchesReach(t *testing.T) {
	path := writeModel(t, mvvStateMachine)

	jsonOut, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("json emission refused: %v", err)
	}
	doc := decodeDocument(t, jsonOut)
	if doc.Reach == nil {
		t.Fatal("the document carries no `reach` block")
	}

	dot, _, err := runGraph(t, "--model", path, "--emit", "dot")
	if err != nil {
		t.Fatalf("dot emission refused: %v", err)
	}

	wantNodes := make([]string, 0, len(doc.Reach.Nodes))
	for _, n := range doc.Reach.Nodes {
		wantNodes = append(wantNodes, n.ID)
	}
	assertSameSet(t, "node", wantNodes, dotNodeIDs(t, dot))

	wantEdges := make([]string, 0, len(doc.Reach.Edges))
	for _, e := range doc.Reach.Edges {
		wantEdges = append(wantEdges, e.From+" -> "+e.To+" ["+e.Rule+"]")
	}
	assertSameSet(t, "edge", wantEdges, dotEdges(t, dot))

	// It renders.
	requireGraphviz(t)
	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("`dot -Tsvg` refused the emitted document: %v\nstderr: %s",
			err, stderr.String())
	}
}

// REQ-71 (the MVV step 4 form of mode agreement).
// mvvStep4ModeAgreement: "`intrastate graph --model <fixture> --as=json |
// jq .data` equals step 2's document, value-for-value."
func mvvStep4ModeAgreement(t *testing.T) {
	path := writeModel(t, mvvStateMachine)

	text, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("text emission refused: %v", err)
	}
	enveloped, _, err := runGraph(t, "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("json emission refused: %v", err)
	}

	var env struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if uerr := json.Unmarshal([]byte(strings.TrimSpace(enveloped)), &env); uerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", uerr, enveloped)
	}
	if env.Type != "ok" {
		t.Errorf("envelope type = %q; want %q", env.Type, "ok")
	}

	var fromText, fromData any
	if uerr := json.Unmarshal([]byte(strings.TrimSpace(text)), &fromText); uerr != nil {
		t.Fatalf("the text-mode document is not JSON: %v", uerr)
	}
	if uerr := json.Unmarshal(env.Data, &fromData); uerr != nil {
		t.Fatalf("`data` is not JSON: %v", uerr)
	}
	if !jsonDeepEqual(fromText, fromData) {
		t.Errorf("`jq .data` does not equal the text-mode document "+
			"value-for-value\n--- text ---\n%s\n--- data ---\n%s",
			text, env.Data)
	}
}

// REQ-72: "`intrastate lint` over a fixture WITH blocking findings still
// refuses identically … while `intrastate graph` over the same model
// succeeds with an ASSERTED document"
// mvvStep5LintNeutrality: "`intrastate lint` over a fixture WITH blocking
// findings still refuses identically (byte-compared against a pre-change
// capture), while `intrastate graph` over the same model succeeds with an
// ASSERTED document (defined content, not incidental) — the neutrality
// oracle, mechanism-independent per C4."
func mvvStep5LintNeutrality(t *testing.T) {
	path := writeModel(t, illegalModel)

	// lint still refuses, and its bytes are compared against the committed
	// pre-change capture rather than against a re-derived expectation.
	var capture strings.Builder
	for _, mode := range []string{"--as=json", "--as=text"} {
		stdout, stderr, err := runCmd(t, "lint", "--model", path, mode)
		if err == nil {
			t.Fatalf("lint %s did not refuse the blocking fixture; the "+
				"neutrality oracle needs a model lint REFUSES", mode)
		}
		capture.WriteString("=== lint " + mode + " ===\nexit: " +
			itoa(clierr.ExitCodeFor(err)) + "\n--- stdout ---\n" + stdout +
			"\n--- stderr ---\n" + stderr + "\n")
	}
	assertGolden(t, "lint_0021_mvv_blocking.txt", capture.String())

	// The export over the SAME model succeeds, with asserted content.
	stdout, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("`graph` refused a model that loads: %v — C4 makes the "+
			"verb succeed for ANY model that loads, including one lint "+
			"refuses", err)
	}
	doc := decodeDocument(t, stdout)
	if doc.Schema != schemaMarker {
		t.Errorf("schema = %q; want %q", doc.Schema, schemaMarker)
	}
	if doc.Model != "lintfix" {
		t.Errorf("model = %q; want %q", doc.Model, "lintfix")
	}
	if doc.Reach == nil || doc.Reach.Abstraction != abstractionMarker {
		t.Errorf("the document carries no `reach` abstraction marker: %+v",
			doc.Reach)
	}
	if len(doc.Rows) == 0 {
		t.Error("the document publishes no rows; the content must be DEFINED")
	}
	generic := decodeGeneric(t, stdout)
	for _, forbidden := range []string{"verdict", "findings"} {
		if _, has := generic[forbidden]; has {
			t.Errorf("the document carries %q; an export is never a lint pass",
				forbidden)
		}
	}
}
