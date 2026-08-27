package cli

// RDR 0010 — the Minimum Viable Validation (`0010:MVV`), run end to end
// through the production Cobra path.
//
// The MVV fixes six steps plus an end-state clause (REQ-77..REQ-85). Every
// step below drives `ExecuteAndEmit` — the same command shape intended for
// CI — and asserts on the JSON envelope the gateway emits. A green exit
// code is never sufficient on its own: each step asserts the VALUE the
// step names.
//
// The end-state (REQ-84) is the runnable inverse invariant: what the model
// AUTHORS as `[rule.emit]` is what `flow resolve` RETURNS as `data.emit`,
// byte-for-byte, over a model with no owned state, no accessor, and no
// artifact — while the state-machine control behaves as it did.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// mvvPartial0010 is MVV step 1: `class = "decision-table"`, two observed
// enum dimensions of two values each — each `single_valued` and `required`,
// and discriminated by `[rule.guard.all.<key>]` atoms — with
// `[rule.match.recognized]` binding the outcome, one recognized outcome,
// THREE ordinary rules each carrying `[rule.emit]`, no `[initial]`, no
// accessors.
const mvvPartial0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.b]
provenance = "observed"
kind = "enum"
domain = ["p", "q"]
single_valued = true
required = true

[[rule]]
id = "cell-xp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "alpha"

[[rule]]
id = "cell-xq"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.guard.all.b]
eq = "q"
[rule.emit]
verdict = "beta"

[[rule]]
id = "cell-yp"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.guard.all.b]
eq = "p"
[rule.emit]
verdict = "gamma"
`

// mvvFourthRule is MVV step 3's fourth ORDINARY rule, closing the product.
// Its authored block is what step 4 asserts `data.emit` equals.
const mvvFourthRule = `
[[rule]]
id = "cell-yq"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "y"
[rule.guard.all.b]
eq = "q"
[rule.emit]
verdict = "delta"
code = "4"
`

// mvvFourthRuleAuthored is the authored block of `cell-yq`, spelled here so
// the step-4 assertion compares against what the FIXTURE says rather than
// against a value re-derived from the implementation.
var mvvFourthRuleAuthored = map[string]any{"verdict": "delta", "code": "4"}

// mvvComplete0010 is MVV step 3: the fourth rule added.
const mvvComplete0010 = mvvPartial0010 + mvvFourthRule

// mvvEscapeVariant0010 is MVV step 3's variant: the fourth rule replaced by
// an escape row rescuing `no_match`, carrying `[rule.emit]`.
const mvvEscapeVariant0010 = mvvPartial0010 + `
[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
verdict = "fallback"
`

// mvvClassOmitted0010 is MVV step 6's first negative control: the same
// table with `class` omitted.
//
// DEVIATION (deviations.md D1): the control carries the owned-tag
// scaffolding `0002:C4`'s write-block arm demands of a state machine and
// declares NO `[initial]`. A byte-for-byte `class`-strip of the decision
// table refuses at LOAD with `malformed rule shape` — the very arm REQ-16
// conditions on class — so it could not exhibit the lint finding MVV step 6
// names. What the control pins is unchanged: a state-machine model with an
// empty `[initial]` still takes `0006:C18`'s missing-root finding.
const mvvClassOmitted0010 = `outcomes = ["decide"]

[model]
id = "dt"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.a]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["seen"]
single_valued = true
required = true

[read.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"

[write.state]
role = "state"
path = "s.own"
keys = ["status"]
timeout = "2s"
read_back = true

[[rule]]
id = "cell-x"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.a]
eq = "x"
[rule.write]
status = "seen"
`

// mvvOwnedTagAdded0010 is MVV step 6's second negative control: the same
// declared decision table with one `provenance = "owned"` tag added.
var mvvOwnedTagAdded0010 = strings.Replace(mvvComplete0010, "[tags.a]", `[tags.own]
provenance = "owned"
kind = "scalar"

[tags.a]`, 1)

// REQ-MVV / `0010:MVV` (REQ-77..REQ-85): the six steps and the end-state,
// run end to end through the production Cobra path.
// HAPPY PATH
func TestMVV_StatelessDecisionTables(t *testing.T) {
	t.Run("step 1 — the decision table is authorable", mvvStep1Authorable)
	t.Run("step 2 — partial lint reports the coverage gap", mvvStep2PartialLint)
	t.Run("step 3 — complete lint is exactly []", mvvStep3CompleteLint)
	t.Run("step 3 variant — escape row yields only the advisory", mvvStep3EscapeVariant)
	t.Run("step 4 — resolve returns rule + emit with no --artifact", mvvStep4Resolve)
	t.Run("step 5 — dump renders emit; next carries empty required", mvvStep5DumpAndNext)
	t.Run("step 6 — negative controls", mvvStep6NegativeControls)
	t.Run("end state — emit round-trips authored to resolved", mvvEndStateRoundTrip)
}

// REQ-77 / MVV step 1: "Author a decision-table model: `class =
// \"decision-table\"`, two observed enum dimensions of two values each —
// each declared `single_valued` and `required`, and **discriminated by
// `[rule.guard.all.<key>]` atoms**, with `[rule.match.recognized]` binding
// the outcome — one recognized outcome, three ordinary rules each carrying
// `[rule.emit]`, no `[initial]`, no accessors."
func mvvStep1Authorable(t *testing.T) {
	m, err := table.Load([]byte(mvvPartial0010), "mvv.toml")
	if err != nil {
		t.Fatalf("the MVV decision table is not authorable; refused: %v", err)
	}

	if !table.IsDecisionTable(m) {
		t.Errorf("Class = %q; want the declared decision-table class", m.Class)
	}
	if len(m.Initial) != 0 {
		t.Errorf("Initial = %v; the MVV table declares no `[initial]`", m.Initial)
	}
	if len(m.Readers)+len(m.Writers)+len(m.Gates) != 0 {
		t.Error("the MVV table declares an accessor; it declares none")
	}
	for key, decl := range m.Tags {
		if decl.Provenance == table.ProvenanceOwned {
			t.Errorf("the MVV table declares the owned tag %q; it declares "+
				"zero owned tags", key)
		}
	}
	if len(m.Rows) != 3 {
		t.Errorf("normalized to %d rows; the step-1 table carries three "+
			"ordinary rules", len(m.Rows))
	}
	for _, r := range m.Rows {
		if len(r.Emit) == 0 {
			t.Errorf("%s carries no emit; each of the three ordinary rules "+
				"carries `[rule.emit]`", r.Identity())
		}
	}
}

// REQ-78 / MVV step 2: "`intrastate lint --model m.toml --as=json` → exit 2
// with one coverage finding naming the uncovered cell (a *positive*
// finding, proving coverage ran with no root declared)."
func mvvStep2PartialLint(t *testing.T) {
	path := writeModel(t, mvvPartial0010)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("the partial table linted clean; want exit 2 with a "+
			"coverage finding:\n%s", stdout)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit code = %d; want 2", got)
	}

	env := parseFailure(t, stdout)
	if env.Findings == nil {
		t.Fatalf("the failure line carries no `findings`:\n%s", stdout)
	}
	var gaps int
	for _, f := range *env.Findings {
		if f.Code == "graph-coverage-gap" {
			gaps++
			if f.Element != "dt/decide" {
				t.Errorf("coverage finding element = %q; want the group "+
					"`dt/decide`", f.Element)
			}
		}
	}
	if gaps != 1 {
		t.Errorf("%d coverage findings; want exactly one — the POSITIVE "+
			"finding proving coverage ran with no root declared:\n%s",
			gaps, stdout)
	}
}

// REQ-79 / MVV step 3: "Add the fourth rule; lint → exit 0, findings
// exactly `[]` — no code from the 0006 taxonomy present (A10)."
func mvvStep3CompleteLint(t *testing.T) {
	path := writeModel(t, mvvComplete0010)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the complete table did not lint clean: %v\n%s", err, stdout)
	}

	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the success envelope carries no `data.findings`:\n%s", stdout)
	}
	if n := len(*env.Data.Findings); n != 0 {
		t.Errorf("findings = %d; want exactly `[]` — no code from the 0006 "+
			"taxonomy is present:\n%s", n, stdout)
	}
}

// REQ-80 / MVV step 3 variant: "replace the fourth rule with an escape row
// rescuing `no_match` carrying `[rule.emit]` → exit 0 with only the
// `graph-coverage-closed-by-escape` advisory (A9)."
func mvvStep3EscapeVariant(t *testing.T) {
	path := writeModel(t, mvvEscapeVariant0010)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the escape variant did not exit 0: %v\n%s", err, stdout)
	}

	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		t.Fatalf("the success envelope carries no `data.findings`:\n%s", stdout)
	}
	var codes []string
	for _, f := range *env.Data.Findings {
		codes = append(codes, f.Code)
	}
	if !slices.Equal(codes, []string{"graph-coverage-closed-by-escape"}) {
		t.Errorf("findings = %v; want exactly "+
			"[graph-coverage-closed-by-escape] — never a bare green", codes)
	}
}

// REQ-81 / MVV step 4: "`intrastate flow resolve --model m.toml --outcome
// <o> --tag a=x --tag b=y --as=json` with no `--artifact` → exit 0;
// `data.rule` is the fourth rule's id, `data.emit` equals its authored
// block, `data.next`, `data.writes`, `data.owned` are empty, `data.readers`
// is `[]`."
func mvvStep4Resolve(t *testing.T) {
	model := writeFlowModel(t, mvvComplete0010)

	data := flowData(t, requireSuccess(t, "flow", "resolve",
		"--model", model, "--outcome", "decide",
		"--tag", "a=y", "--tag", "b=q", "--as=json"))

	if data["rule"] != "cell-yq" {
		t.Fatalf("`data.rule` = %#v; want the fourth rule's id %q",
			data["rule"], "cell-yq")
	}

	emit, ok := data["emit"].(map[string]any)
	if !ok {
		t.Fatalf("`data.emit` is not a JSON object: %#v (payload keys: %v)",
			data["emit"], keysOf(data))
	}
	if !reflect.DeepEqual(emit, mvvFourthRuleAuthored) {
		t.Errorf("`data.emit` = %#v; want the fourth rule's AUTHORED block "+
			"%#v", emit, mvvFourthRuleAuthored)
	}

	for _, key := range []string{"next", "writes", "owned"} {
		obj, ok := data[key].(map[string]any)
		if !ok || len(obj) != 0 {
			t.Errorf("`data.%s` = %#v; want the empty object", key, data[key])
		}
	}
	readers, ok := data["readers"].([]any)
	if !ok || len(readers) != 0 {
		t.Errorf("`data.readers` = %#v; want `[]`", data["readers"])
	}
}

// REQ-82 / MVV step 5: "`intrastate dump --model m.toml` renders the `emit`
// column for every row with no `[dump]` declared; `flow next --model
// m.toml --as=json` with no `--artifact` → exit 0, every candidate with
// empty `required` (A11)."
//
// DEVIATION (deviations.md D2): the repo ships no root `dump` CLI verb —
// `NewRootCmd` registers `version`, `lint`, and `flow` only. The dump is
// the `internal/table::Dump` surface `0002:C19` fixes, so the "renders the
// `emit` column for every row" half is taken there. The `flow next` half is
// taken through the CLI as written.
func mvvStep5DumpAndNext(t *testing.T) {
	t.Run("dump renders the emit column for every row", func(t *testing.T) {
		m, err := table.Load([]byte(mvvComplete0010), "mvv.toml")
		if err != nil {
			t.Fatalf("the MVV table must load clean: %v", err)
		}
		if slices.Contains(m.DumpOrder, "emit") == false {
			t.Errorf("DumpOrder = %v; with no `[dump]` declared the default "+
				"set includes `emit`", m.DumpOrder)
		}

		dump := table.Dump(m)
		var rows int
		for _, line := range strings.Split(dump, "\n") {
			if !strings.HasPrefix(line, "identity=") {
				continue
			}
			rows++
			if !strings.Contains(line, " emit=[") {
				t.Errorf("dump row renders no `emit` cell:\n%s", line)
			}
		}
		if rows != 4 {
			t.Errorf("the dump rendered %d rows; the complete table carries 4",
				rows)
		}
	})

	t.Run("flow next exits 0 with empty required on every candidate", func(t *testing.T) {
		model := writeFlowModel(t, mvvComplete0010)
		data := flowData(t, requireSuccess(t, "flow", "next",
			"--model", model, "--as=json"))

		candidates, ok := objectsAt(data, "candidates")
		if !ok || len(candidates) == 0 {
			t.Fatalf("`candidates` is empty or malformed: %#v",
				data["candidates"])
		}
		for _, c := range candidates {
			req, ok := c["required"].([]any)
			if !ok {
				t.Errorf("candidate %#v carries no `required` array: %#v",
					c["rule"], c["required"])
				continue
			}
			if len(req) != 0 {
				t.Errorf("candidate %#v: `required` = %#v; want empty (A11)",
					c["rule"], req)
			}
		}
	})
}

// REQ-83 / MVV step 6: "Negative controls: the same model with `class`
// omitted → lint `graph-dangling-edge` naming the absent `[initial]`
// (0006:C18 unchanged); with one `provenance = \"owned\"` tag added → load
// refusal `malformed model declaration`; `models/rdr.toml` lints and
// resolves as before."
func mvvStep6NegativeControls(t *testing.T) {
	t.Run("class omitted lints graph-dangling-edge at element=model", func(t *testing.T) {
		path := writeModel(t, mvvClassOmitted0010)

		stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
		if err == nil {
			t.Fatalf("the class-omitted control linted clean; want the "+
				"missing-root finding:\n%s", stdout)
		}
		env := parseFailure(t, stdout)
		if env.Findings == nil {
			t.Fatalf("no `findings`:\n%s", stdout)
		}
		var rooted bool
		for _, f := range *env.Findings {
			if f.Code == "graph-dangling-edge" && f.Element == "model" {
				rooted = true
			}
		}
		if !rooted {
			t.Errorf("no `graph-dangling-edge` at `element = model` naming "+
				"the absent `[initial]` (`0006:C18` unchanged):\n%s", stdout)
		}
	})

	t.Run("an owned tag added refuses at load", func(t *testing.T) {
		model := writeFlowModel(t, mvvOwnedTagAdded0010)

		_, err := table.Load([]byte(mvvOwnedTagAdded0010), "mvv-owned.toml")
		if err == nil {
			t.Fatal("a declared decision table carrying one owned tag loaded " +
				"clean; want `malformed model declaration`")
		}
		if cat, ok := table.CategoryOf(err); !ok ||
			cat != table.CatMalformedModelDeclaration {
			t.Errorf("category = %q; want %q",
				cat, table.CatMalformedModelDeclaration)
		}

		// The refusal surfaces at the CLI as a model-invalid refusal, not a
		// crash and not a plan.
		requireRefusal(t, "flow-model-invalid", 2,
			"flow", "resolve", "--model", model, "--outcome", "decide",
			"--tag", "a=x", "--tag", "b=p", "--as=json")
	})

	t.Run("models/rdr.toml lints as before", func(t *testing.T) {
		path := filepath.Join(repoRootFor(t), checkedInModelPath)
		stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
		if err != nil {
			t.Fatalf("the checked-in navigator model no longer lints "+
				"clean: %v\n%s", err, stdout)
		}
	})
}

// REQ-84 / `0010:MVV` end-state: "a stateless model lints for coverage and
// resolves to `rule` + `emit` with no owned state, no accessor, and no
// artifact, while every state-machine model behaves as it did."
//
// The RUNNABLE INVERSE INVARIANT: the emit block the model AUTHORS is what
// `flow resolve` RETURNS, byte-for-byte, for every selectable cell. A green
// exit code is explicitly not sufficient — the reconstructed value is
// compared against the authored one, key for key and byte for byte.
func mvvEndStateRoundTrip(t *testing.T) {
	model := writeFlowModel(t, mvvComplete0010)

	// The authored blocks, read off the FIXTURE rather than off the loader,
	// so the oracle is what the author wrote.
	authored := map[string]map[string]any{
		"cell-xp": {"verdict": "alpha"},
		"cell-xq": {"verdict": "beta"},
		"cell-yp": {"verdict": "gamma"},
		"cell-yq": {"verdict": "delta", "code": "4"},
	}
	selection := map[string][2]string{
		"cell-xp": {"x", "p"},
		"cell-xq": {"x", "q"},
		"cell-yp": {"y", "p"},
		"cell-yq": {"y", "q"},
	}

	for rule, want := range authored {
		t.Run(rule, func(t *testing.T) {
			cell := selection[rule]
			stdout := requireSuccess(t, "flow", "resolve",
				"--model", model, "--outcome", "decide",
				"--tag", "a="+cell[0], "--tag", "b="+cell[1], "--as=json")
			data := flowData(t, stdout)

			if data["rule"] != rule {
				t.Fatalf("`rule` = %#v; want %q for the cell (a=%s, b=%s)",
					data["rule"], rule, cell[0], cell[1])
			}

			got, ok := data["emit"].(map[string]any)
			if !ok {
				t.Fatalf("`emit` is not a JSON object: %#v", data["emit"])
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("`emit` = %#v; want the AUTHORED block %#v — the "+
					"answer a decision table returns is what the row says",
					got, want)
			}

			// Byte-level: the wire object must be the authored pairs in key
			// order, with no re-encoding of the values.
			var wire struct {
				Data struct {
					Emit json.RawMessage `json:"emit"`
				} `json:"data"`
			}
			if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &wire); err != nil {
				t.Fatalf("stdout is not one JSON object: %v", err)
			}
			for k, v := range want {
				pair := `"` + k + `":"` + v.(string) + `"`
				if !strings.Contains(string(wire.Data.Emit), pair) {
					t.Errorf("the wire's `emit` %s does not carry %s "+
						"byte-for-byte", wire.Data.Emit, pair)
				}
			}

			// No owned state, no accessor, no artifact.
			if owned, _ := data["owned"].(map[string]any); len(owned) != 0 {
				t.Errorf("`owned` = %#v; the end-state carries no owned state",
					data["owned"])
			}
			if readers, _ := data["readers"].([]any); len(readers) != 0 {
				t.Errorf("`readers` = %#v; the end-state invokes no accessor",
					data["readers"])
			}
		})
	}

	t.Run("every state-machine model behaves as it did", func(t *testing.T) {
		smModel := writeFlowModel(t, flowMVVModel)
		art := seedArtifact(t, smModel, "status=draft", "stale=obsolete")

		data := flowData(t, requireSuccess(t, "flow", "resolve",
			"--model", smModel,
			"--artifact", artifactBinding(flowStateRole, art),
			"--outcome", "advance", "--as=json"))

		if data["rule"] != "advance-draft" {
			t.Errorf("`rule` = %#v; the 0005 state-machine fixture resolves "+
				"as it did", data["rule"])
		}
		owned, ok := data["owned"].(map[string]any)
		if !ok || owned["status"] != "draft" {
			t.Errorf("`owned` = %#v; owned state is still read from the "+
				"artifact", data["owned"])
		}
		// The state machine's rows author no emit, so its payload carries the
		// empty object — the field is added for both classes.
		emit, ok := data["emit"].(map[string]any)
		if !ok {
			t.Fatalf("`emit` is not a JSON object: %#v", data["emit"])
		}
		if len(emit) != 0 {
			t.Errorf("`emit` = %#v; the 0005 fixture authors no emit block",
				emit)
		}
	})
}

// REQ-85 / `0010:MVV` consumer-side acceptance: "That rewrite happens in
// the consumer repo, not here — intrastate stays generic and
// `models/rdr.toml` is untouched"
// ADVERSARIAL — negative REQ: `models/rdr.toml` MUST NOT be edited.
func TestReq85_TheCheckedInNavigatorModelIsUntouched(t *testing.T) {
	src, err := os.ReadFile(filepath.Join(repoRootFor(t), checkedInModelPath))
	if err != nil {
		t.Fatalf("read %s: %v", checkedInModelPath, err)
	}

	if strings.Contains(string(src), "class =") {
		t.Error("models/rdr.toml declares a `class`; intrastate stays " +
			"generic and the navigator model is untouched by this RDR — " +
			"the consumer rewrite is tracked on rdr#tmxk")
	}
	if strings.Contains(string(src), "[rule.emit]") {
		t.Error("models/rdr.toml authors `[rule.emit]`; it is untouched by " +
			"this RDR")
	}

	m, lerr := table.Load(src, "models/rdr.toml")
	if lerr != nil {
		t.Fatalf("models/rdr.toml no longer loads: %v", lerr)
	}
	if table.IsDecisionTable(m) {
		t.Error("models/rdr.toml reads as a decision table; it is a state " +
			"machine and this RDR leaves it untouched")
	}
}
