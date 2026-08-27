package cli

// RDR 0006 — the CLI surface: command placement (REQ-5, REQ-6), the
// normative input contract (REQ-10..REQ-14), and the output envelope
// (REQ-90..REQ-100, REQ-128).
//
// Every assertion drives ExecuteAndEmit — the production emission path —
// and reads the captured streams, never an internal field.

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/graphlint"
)

// asCLIError unwraps err to the structured CLIError the gateway returns.
func asCLIError(err error, target **clierr.CLIError) bool {
	return errors.As(err, target)
}

// --- fixture models ------------------------------------------------------

// legalModel is a minimal conforming model whose one group's two rows
// partition the guard dimension, so it lints clean.
const legalModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "lintfix"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"
read_back = true

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"
`

// illegalModel carries an overlapping ordinary pair and a coverage gap, so
// the run fails with several blocking findings.
const illegalModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "lintfix"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"
read_back = true

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "over-one"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "over-two"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.write]
status = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"
`

// writeModel writes src to a temp file and returns its path.
func writeModel(t *testing.T, src string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "model.toml")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// successEnvelope is the parsed `--as=json` success line.
type successEnvelope struct {
	Type string `json:"type"`
	Data *struct {
		Findings *[]clierr.Finding `json:"findings"`
	} `json:"data"`
}

// failureEnvelope is the parsed `--as=json` failure line. `findings` is a
// TOP-LEVEL sibling of `code`, never nested under an `error` wrapper.
type failureEnvelope struct {
	Code     string            `json:"code"`
	Message  string            `json:"message"`
	Findings *[]clierr.Finding `json:"findings"`
	Error    json.RawMessage   `json:"error"`
	Type     string            `json:"type"`
}

func parseSuccess(t *testing.T, stdout string) successEnvelope {
	t.Helper()

	var env successEnvelope
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	return env
}

func parseFailure(t *testing.T, stdout string) failureEnvelope {
	t.Helper()

	var env failureEnvelope
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	return env
}

// REQ-6: "Root `lint` is deliberately not under RDR 0005's `flow` group" —
// the command is registered at root, not under `flow`.
// REQ-4: root `intrastate lint` is the authoritative CLI surface.
// HAPPY PATH
func TestReq6_LintIsRegisteredAtRootNotUnderFlow(t *testing.T) {
	root := NewRootCmd()

	var lintCmd bool
	for _, c := range root.Commands() {
		if c.Name() == "lint" {
			lintCmd = true
		}
		if c.Name() == "flow" {
			for _, sub := range c.Commands() {
				if sub.Name() == "lint" {
					t.Errorf("`lint` is registered under the `flow` group; " +
						"it is deliberately a ROOT command")
				}
			}
		}
	}
	if !lintCmd {
		var names []string
		for _, c := range root.Commands() {
			names = append(names, c.Name())
		}
		t.Fatalf("no root `lint` command is registered; root commands = %v",
			names)
	}

	// It runs: the authoritative surface is reachable by that name.
	path := writeModel(t, legalModel)
	if _, _, err := runCmd(t, "lint", "--model", path, "--as=json"); err != nil {
		if clierr.ErrorCode(err) == "command-error" {
			t.Fatalf("root `lint` is not invocable: %v", err)
		}
	}
}

// REQ-5: "Lint command success and failure MUST route through
// `respond.OK`, `respond.Fail`, and `CLIError`; the command MUST NOT write
// directly to stdout or stderr."
// REQ-99 / SC-4: "both modes return the same exit behavior with no direct
// stdout/stderr writes."
// ADVERSARIAL
func TestReq5And99_BothModesRouteThroughTheGatewayWithNoDirectWrites(t *testing.T) {
	legal := writeModel(t, legalModel)
	illegal := writeModel(t, illegalModel)

	// JSON success: exactly ONE terminal record on stdout, and it is the
	// respond envelope — a direct write would add a second line.
	stdout, _, err := runCmd(t, "lint", "--model", legal, "--as=json")
	if err != nil {
		t.Fatalf("legal model failed: %v", err)
	}
	if n := len(nonEmptyLines(stdout)); n != 1 {
		t.Errorf("JSON success wrote %d stdout lines; the gateway emits "+
			"exactly one terminal record and the command writes nothing "+
			"directly:\n%s", n, stdout)
	}
	if env := parseSuccess(t, stdout); env.Type != "ok" {
		t.Errorf("success envelope type = %q; want %q", env.Type, "ok")
	}

	// JSON failure: exactly one terminal record, and it is the CLIError.
	failOut, _, ferr := runCmd(t, "lint", "--model", illegal, "--as=json")
	if ferr == nil {
		t.Fatalf("illegal model succeeded; stdout:\n%s", failOut)
	}
	if n := len(nonEmptyLines(failOut)); n != 1 {
		t.Errorf("JSON failure wrote %d stdout lines; want exactly one:\n%s",
			n, failOut)
	}

	// Text mode returns the SAME exit behaviour.
	_, _, textErr := runCmd(t, "lint", "--model", legal)
	if textErr != nil {
		t.Errorf("text mode failed on the legal model: %v", textErr)
	}
	_, _, textFail := runCmd(t, "lint", "--model", illegal)
	if textFail == nil {
		t.Fatal("text mode succeeded on the illegal model; both modes must " +
			"return the same exit behaviour")
	}
	if got, want := clierr.ExitCodeFor(textFail), clierr.ExitCodeFor(ferr); got != want {
		t.Errorf("text exit = %d, json exit = %d; both modes must return the "+
			"same exit behaviour", got, want)
	}
}

// REQ-10: "The command's **input contract is normative** … only cosmetic
// flag spelling defers to RDR 0005."
// REQ-11: "`--flow <id>` selects the model … `--model <path>` names one
// explicitly. They are mutually exclusive — supplying both is a
// `GroupUserEnv` usage error, not a precedence rule."
// INPUT EDGE
func TestReq10And11_FlowAndModelAreMutuallyExclusiveUsageErrors(t *testing.T) {
	path := writeModel(t, legalModel)

	_, _, err := runCmd(t, "lint", "--model", path, "--flow", "rdr", "--as=json")
	if err == nil {
		t.Fatal("supplying both --model and --flow succeeded; they are " +
			"mutually exclusive, not a precedence rule")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if ce.Group != clierr.GroupUserEnv {
		t.Errorf("group = %v; want GroupUserEnv — it is a usage error", ce.Group)
	}
	// The group and the exit code are shared by every usage error, so
	// neither distinguishes THIS refusal. The stable code does, and a
	// consumer branches on it.
	if ce.Code != "flag-mutually-exclusive" {
		t.Errorf("code = %q; want %q — the refusal is identified by its "+
			"stable code, not merely by its group", ce.Code,
			"flag-mutually-exclusive")
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}

	// Both flags exist on the command, so the contract is instantiable.
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() != "lint" {
			continue
		}
		for _, name := range []string{"model", "flow"} {
			if c.Flags().Lookup(name) == nil {
				t.Errorf("the lint command declares no --%s flag; the input "+
					"contract is normative", name)
			}
		}
	}
}

// REQ-12: "Exactly one model is linted per invocation: a corpus is linted
// by invoking the command once per model, so a run's findings never span
// models."
// REQ-13: "The finding's `model` field is the `[model].id` from the model
// itself, never the path."
// BOUNDARY
func TestReq12And13_OneModelPerInvocationAndFindingsNameTheModelID(t *testing.T) {
	path := writeModel(t, illegalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("illegal model succeeded; stdout:\n%s", stdout)
	}
	env := parseFailure(t, stdout)
	if env.Findings == nil {
		t.Fatalf("the failure envelope carries no findings key:\n%s", stdout)
	}
	if len(*env.Findings) == 0 {
		t.Fatalf("the failure envelope carries an empty findings list "+
			"alongside a blocking failure:\n%s", stdout)
	}

	// A run's findings never span models: every finding names the one
	// `[model].id`, and it is the id, never the path.
	for _, f := range *env.Findings {
		if f.Model != "lintfix" {
			t.Errorf("finding %s carries model %q; want the `[model].id` "+
				"%q", f.Code, f.Model, "lintfix")
		}
		if strings.Contains(f.Model, string(os.PathSeparator)) ||
			strings.HasSuffix(f.Model, ".toml") {
			t.Errorf("finding %s carries the PATH %q in its model field",
				f.Code, f.Model)
		}
	}
}

// REQ-14: "Until RDR 0005 lands discovery, only `--model <path>` is
// instantiable."
// INPUT EDGE
func TestReq14_ModelPathIsTheInstantiableForm(t *testing.T) {
	// `--model <path>` lints.
	path := writeModel(t, legalModel)
	if _, _, err := runCmd(t, "lint", "--model", path, "--as=json"); err != nil {
		t.Fatalf("--model <path> is not instantiable: %v", err)
	}

	// `--flow <id>` alone refuses rather than silently linting nothing:
	// discovery has not landed, so it cannot select a model.
	stdout, _, err := runCmd(t, "lint", "--flow", "rdr", "--as=json")
	if err == nil {
		t.Errorf("--flow <id> succeeded without config discovery; it must "+
			"refuse rather than report a clean model; stdout:\n%s", stdout)
	}
}

// REQ-90: "Graph lint failure MUST return one aggregate `CLIError` with
// code `graph-lint-failed` and `GroupUserEnv`" — exit 2.
// HAPPY PATH
func TestReq90_FailureIsOneAggregateCLIErrorAtGroupUserEnv(t *testing.T) {
	path := writeModel(t, illegalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("illegal model succeeded; stdout:\n%s", stdout)
	}
	if code := clierr.ErrorCode(err); code != graphlint.AggregateCode {
		t.Errorf("aggregate code = %q; want %q", code, graphlint.AggregateCode)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("the failure is not a CLIError: %v", err)
	}
	if ce.Group != clierr.GroupUserEnv {
		t.Errorf("group = %v; want GroupUserEnv", ce.Group)
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}

	// ONE aggregate error, not one per finding: the envelope is a single
	// object carrying the list.
	if n := len(nonEmptyLines(stdout)); n != 1 {
		t.Errorf("the failure wrote %d stdout lines; want one aggregate "+
			"envelope:\n%s", n, stdout)
	}
}

// REQ-91: "the individual blocking findings MUST remain machine-readable
// in JSON mode through an append-only typed `findings` field owned by
// `clierr`, not through a verb-local wrapper or a text-only `Detail`
// string."
// REQ-92: "That field is a **top-level sibling of `code`** … reached as
// `.findings` … lint MUST NOT restructure that envelope to create one."
// BOUNDARY
func TestReq91And92_FindingsAreATopLevelSiblingOfCodeOnTheFailureLine(t *testing.T) {
	path := writeModel(t, illegalModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("illegal model succeeded; stdout:\n%s", stdout)
	}

	// Reached as `.findings`, a sibling of `.code`.
	var raw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, stdout)
	}
	if _, ok := raw["code"]; !ok {
		t.Errorf("the failure line carries no top-level `code` key:\n%s", stdout)
	}
	if _, ok := raw["findings"]; !ok {
		t.Fatalf("the failure line carries no top-level `findings` key; it "+
			"is a SIBLING of `code`, not `.error.findings`:\n%s", stdout)
	}
	// No `{"type":"failed","error":{…}}` wrapper was created.
	if _, ok := raw["error"]; ok {
		t.Errorf("the failure line carries an `error` wrapper; EmitJSON "+
			"marshals the CLIError itself and lint MUST NOT restructure the "+
			"shipped envelope:\n%s", stdout)
	}
	if _, ok := raw["type"]; ok {
		t.Errorf("the failure line carries a `type` discriminator; that is "+
			"the output contract's revisit note, not this RDR's to add:\n%s",
			stdout)
	}

	// The findings are structured data, not a text-only Detail string.
	env := parseFailure(t, stdout)
	if env.Findings == nil || len(*env.Findings) == 0 {
		t.Fatalf("the findings list is absent or empty:\n%s", stdout)
	}
	for _, f := range *env.Findings {
		if f.Code == "" || f.Message == "" {
			t.Errorf("a finding is not machine-readable: %+v", f)
		}
	}
}

// REQ-93: "Lint success MUST carry non-blocking findings under the
// existing `respond.Success.Data` payload as `data.findings`."
// HAPPY PATH
func TestReq93_SuccessCarriesNonBlockingFindingsAtDataFindings(t *testing.T) {
	// A clean model carrying an advisory finding: a group closed by a bare
	// escape row.
	path := writeModel(t, advisoryModel)

	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("the advisory model failed: %v\nstdout:\n%s", err, stdout)
	}
	env := parseSuccess(t, stdout)
	if env.Data == nil {
		t.Fatalf("the success envelope carries no `data` key:\n%s", stdout)
	}
	if env.Data.Findings == nil {
		t.Fatalf("the success envelope carries no `data.findings` key:\n%s",
			stdout)
	}
	if len(*env.Data.Findings) == 0 {
		t.Fatalf("the advisory model reported no non-blocking findings; the "+
			"bare-escape closure must be observable without inspecting the "+
			"model:\n%s", stdout)
	}
	for _, f := range *env.Data.Findings {
		if graphlint.IsBlocking(f.Code) {
			t.Errorf("the success payload carries the blocking code %q", f.Code)
		}
	}
}

// REQ-94: "On both surfaces the key MUST be emitted even when the list is
// empty … the `findings` field MUST NOT be `omitempty`, **and** the
// success path MUST assign `respond.Success.Data` a non-nil struct value"
// REQ-128 / SC-16: "the `findings` key is present with value `[]` on
// success (`data.findings`), asserted against a golden JSON file so an
// omitted key fails, and present on the failure envelope alongside
// blocking findings."
// BOUNDARY
func TestReq94And128_TheFindingsKeyIsAlwaysEmittedIncludingTheEmptyList(t *testing.T) {
	// Success: the key is present with value `[]`. The oracle is the RAW
	// JSON, so an omitted key fails rather than decoding to a nil slice.
	legal := writeModel(t, legalModel)
	stdout, _, err := runCmd(t, "lint", "--model", legal, "--as=json")
	if err != nil {
		t.Fatalf("the legal model failed: %v\nstdout:\n%s", err, stdout)
	}

	var raw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &raw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, stdout)
	}
	dataRaw, ok := raw["data"]
	if !ok {
		t.Fatalf("the success line carries no `data` key — `Data` is itself "+
			"`json:\"data,omitempty\"`, so the success path must assign it a "+
			"NON-NIL struct value:\n%s", stdout)
	}
	var data map[string]json.RawMessage
	if jerr := json.Unmarshal(dataRaw, &data); jerr != nil {
		t.Fatalf("`data` is not an object: %v\n%s", jerr, stdout)
	}
	findingsRaw, ok := data["findings"]
	if !ok {
		t.Fatalf("the success line carries no `data.findings` key; the empty "+
			"list is the proof's receipt:\n%s", stdout)
	}
	if got := strings.TrimSpace(string(findingsRaw)); got != "[]" {
		t.Errorf("data.findings = %s; want the empty list `[]` on a clean "+
			"model (a `null` means the field carried a nil slice)", got)
	}

	// Failure: the key is present alongside the blocking findings.
	illegal := writeModel(t, illegalModel)
	failOut, _, ferr := runCmd(t, "lint", "--model", illegal, "--as=json")
	if ferr == nil {
		t.Fatalf("the illegal model succeeded:\n%s", failOut)
	}
	var failRaw map[string]json.RawMessage
	if jerr := json.Unmarshal([]byte(strings.TrimSpace(failOut)), &failRaw); jerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", jerr, failOut)
	}
	if _, ok := failRaw["findings"]; !ok {
		t.Errorf("the failure line carries no `findings` key:\n%s", failOut)
	}
}

// REQ-95: "The `Finding` type MUST be defined in `clierr` as a
// subsystem-agnostic record so `clierr` gains no dependency on the
// graph-lint package."
// REQ-96: "`clierr.Finding` carries those as **declared string-typed
// fields, not enums**: `Key`, `Operator`, `Literal`, `Block`, and `Class`
// are `string`"
// BOUNDARY
func TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields(t *testing.T) {
	// The type is defined in clierr: this reference compiles only if it is.
	var f clierr.Finding

	// The atom and class fields are declared `string`, so `clierr`
	// ascribes them no meaning and imports neither RDR 0002's Block
	// vocabulary nor RDR 0001's RefusalKind.
	//
	// Assigning a bare `""` would NOT say that: an untyped string constant
	// is assignable to any named string type, so `f.Block = ""` compiles
	// just as well if `Block` were declared `table.Block`. Taking the
	// address and binding it to a `*string` does not convert, so these
	// declarations compile only if the fields are `string` itself.
	var (
		_ *string = &f.Key
		_ *string = &f.Operator
		_ *string = &f.Literal
		_ *string = &f.Block
		_ *string = &f.Class
	)

	// clierr gains no dependency on the graph-lint package.
	root := repoRootFor(t)
	src, err := os.ReadFile(filepath.Join(root, "internal", "cli", "clierr", "clierr.go"))
	if err != nil {
		t.Fatalf("read clierr.go: %v", err)
	}
	for _, banned := range []string{"internal/graphlint", "internal/table", "internal/resolve"} {
		if strings.Contains(string(src), banned) {
			t.Errorf("clierr imports %q; the Finding record is "+
				"subsystem-agnostic precisely so it does not", banned)
		}
	}
}

// REQ-97: "Text mode MUST enumerate every finding's code and message,
// which requires extending `clierr.EmitText` (failure) and the
// `respond.OK` text branch (success…)"
// REQ-98 / SC-4: "the set of finding codes appearing in text output equals
// the set in JSON output for the same fixture. … a text renderer that
// drops any finding fails."
// BOUNDARY
func TestReq97And98_TextModeEnumeratesEveryFindingsCodeAndMessage(t *testing.T) {
	for _, tc := range []struct {
		name    string
		model   string
		wantErr bool
	}{
		{"failure", illegalModel, true},
		{"success", advisoryModel, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeModel(t, tc.model)

			jsonOut, jsonErr, err := runCmd(t, "lint", "--model", path, "--as=json")
			if (err != nil) != tc.wantErr {
				t.Fatalf("json run: err = %v; wantErr = %v\n%s%s",
					err, tc.wantErr, jsonOut, jsonErr)
			}
			wantCodes := codesFromJSON(t, jsonOut, tc.wantErr)
			if len(wantCodes) == 0 {
				t.Fatalf("the %s fixture reported no findings in JSON mode, "+
					"so the text oracle is vacuous:\n%s", tc.name, jsonOut)
			}

			textOut, textErrOut, terr := runCmd(t, "lint", "--model", path)
			if (terr != nil) != tc.wantErr {
				t.Fatalf("text run: err = %v; wantErr = %v", terr, tc.wantErr)
			}
			combined := textOut + textErrOut

			// REQ-98 is set EQUALITY, not containment. A one-way
			// JSON->text containment loop green-passes a renderer that
			// invents a code JSON never reported, so the text codes are
			// parsed back out and the two sets compared both directions.
			gotCodes := codesFromText(combined)
			if !slices.Equal(gotCodes, wantCodes) {
				t.Errorf("the text finding-code set is not EQUAL to the "+
					"JSON set:\n  text: %v\n  json: %v\n"+
					"--- text ---\n%s\n--- json ---\n%s",
					gotCodes, wantCodes, combined, jsonOut)
			}
			// And every finding's MESSAGE is enumerated, not merely an
			// aggregate summary.
			for _, msg := range messagesFromJSON(t, jsonOut, tc.wantErr) {
				if msg != "" && !strings.Contains(combined, msg) {
					t.Errorf("text output drops the finding message %q\n"+
						"--- text ---\n%s", msg, combined)
				}
			}
		})
	}
}

// REQ-100: "Model unreadable / not conforming to RDR 0002's schema | per
// RDR 0002 | refused before normalization; lint never runs | none of this
// RDR's codes"
// INPUT EDGE
func TestReq100_NonConformingModelIsRefusedUpstreamWithNoneOfThisRDRsCodes(t *testing.T) {
	for _, tc := range []struct{ name, src string }{
		{"malformed-toml", "this is not ][ toml"},
		{"no-model-table", "outcomes = [\"go\"]\n"},
		{"unsupported-version", "[model]\nid = \"x\"\nversion = 99\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeModel(t, tc.src)

			stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
			if err == nil {
				t.Fatalf("a non-conforming model lint-succeeded:\n%s", stdout)
			}
			// Lint never ran, so NONE of this RDR's codes appears — neither
			// the aggregate nor any finding code.
			code := clierr.ErrorCode(err)
			if code == graphlint.AggregateCode {
				t.Errorf("a non-conforming model returned %q; it is refused "+
					"before normalization and lint never runs", code)
			}
			for _, c := range append(graphlint.BlockingCodes(),
				graphlint.AdvisoryCodes()...) {
				if strings.Contains(stdout, c) {
					t.Errorf("the refusal carries this RDR's code %q:\n%s",
						c, stdout)
				}
			}
		})
	}

	// The discriminating control, so an engine that reports nothing at all
	// cannot satisfy the negative assertions above: a CONFORMING model
	// with a blocking defect DOES return this RDR's aggregate code. The
	// clause draws a line between upstream refusal and lint's own verdict,
	// and both sides must be observable.
	path := writeModel(t, illegalModel)
	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err == nil {
		t.Fatalf("a conforming model with a blocking defect succeeded, so "+
			"the upstream-refusal assertions above are vacuous:\n%s", stdout)
	}
	if code := clierr.ErrorCode(err); code != graphlint.AggregateCode {
		t.Errorf("a conforming model with a blocking defect returned %q; "+
			"want %q — lint DOES run on a model that normalizes",
			code, graphlint.AggregateCode)
	}
}

// REQ-59 / REQ-60 / REQ-126 / SC-18 / SC-22: "Both bounds are
// implementation constants in the command's help output"
// BOUNDARY
func TestReq59And60And126_BothBoundsAppearInTheCommandsHelpOutput(t *testing.T) {
	stdout, stderr, err := runCmd(t, "lint", "--help")
	if err != nil {
		t.Fatalf("lint --help failed: %v", err)
	}
	help := stdout + stderr

	for _, tc := range []struct {
		name  string
		value int
	}{
		{"product bound", graphlint.ProductBound()},
		{"node ceiling", graphlint.NodeCeiling()},
	} {
		if tc.value <= 0 {
			t.Errorf("the published %s is %d; it must be a positive "+
				"implementation constant", tc.name, tc.value)
			continue
		}
		// The value must appear on the line carrying ITS OWN label.
		// Searching the whole help text for the bare number lets the two
		// values be swapped between their labels and stay green, which
		// publishes the wrong number under each name.
		line, ok := helpLineFor(help, tc.name)
		if !ok {
			t.Errorf("the command's help output carries no line labelled "+
				"%q; both bounds are published as named implementation "+
				"constants:\n%s", tc.name, help)
			continue
		}
		// The number is parsed out and compared for EQUALITY. A substring
		// test green-passes any value the published one is a prefix or
		// infix of — `20480` contains `2048` — so a bound that gained a
		// digit would still read as correct.
		got, ok := numberOn(line)
		if !ok {
			t.Errorf("the help line for the %s reads %q, which publishes "+
				"no number; SC-22 asserts the ceiling against the "+
				"PUBLISHED constant, not a test-local value:\n%s",
				tc.name, strings.TrimSpace(line), help)
			continue
		}
		if got != tc.value {
			// A swap between the two labels lands here too, naming the
			// other bound's value as what this label wrongly publishes.
			t.Errorf("the help line for the %s reads %q, publishing %d; "+
				"want %d", tc.name, strings.TrimSpace(line), got, tc.value)
		}
	}
}

// numberOn returns the single decimal number published on a help line, and
// reports whether exactly one exists — two would make "the value under this
// label" ambiguous and the equality assertion unsound.
func numberOn(line string) (int, bool) {
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if len(fields) != 1 {
		return 0, false
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, false
	}
	return n, true
}

// helpLineFor returns the single help line carrying label, and reports
// whether exactly one such line exists — two would make "the line for this
// label" ambiguous and the per-label assertion unsound.
func helpLineFor(help, label string) (string, bool) {
	var found string
	var n int
	for _, line := range strings.Split(help, "\n") {
		if strings.Contains(line, label) {
			found = line
			n++
		}
	}
	if n != 1 {
		return "", false
	}
	return found, true
}

// --- helpers -------------------------------------------------------------

func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
}

func codesFromJSON(t *testing.T, stdout string, failure bool) []string {
	t.Helper()

	var out []string
	for _, f := range findingsFromJSON(t, stdout, failure) {
		if !slices.Contains(out, f.Code) {
			out = append(out, f.Code)
		}
	}
	slices.Sort(out)
	return out
}

// codesFromText recovers the finding codes text mode rendered, so REQ-98's
// "the set of finding codes appearing in text output equals the set in
// JSON output" can be asserted as an equality rather than a containment.
//
// `EmitFindingsText` renders one finding per line as "  <code>: <message>",
// and quotes identity values and folds line-affecting characters out of the
// message, so one finding is one line and the code is the token before the
// first colon. The aggregate `error: <code>: <message>` line is NOT a
// finding and is skipped: it carries the envelope code, which JSON reports
// in `error.code`, not in `findings`.
func codesFromText(combined string) []string {
	var out []string
	for _, line := range strings.Split(combined, "\n") {
		if !strings.HasPrefix(line, "  ") {
			continue
		}
		body := strings.TrimSpace(line)
		code, _, ok := strings.Cut(body, ": ")
		if !ok || code == "" || strings.ContainsAny(code, " \t") {
			// "  detail: ..." / "  hint: ..." are envelope lines, and
			// anything else indented is not a finding line.
			continue
		}
		if code == "detail" || code == "hint" {
			continue
		}
		if !slices.Contains(out, code) {
			out = append(out, code)
		}
	}
	slices.Sort(out)
	return out
}

func messagesFromJSON(t *testing.T, stdout string, failure bool) []string {
	t.Helper()

	var out []string
	for _, f := range findingsFromJSON(t, stdout, failure) {
		out = append(out, f.Message)
	}
	return out
}

func findingsFromJSON(t *testing.T, stdout string, failure bool) []clierr.Finding {
	t.Helper()

	if failure {
		env := parseFailure(t, stdout)
		if env.Findings == nil {
			return nil
		}
		return *env.Findings
	}
	env := parseSuccess(t, stdout)
	if env.Data == nil || env.Data.Findings == nil {
		return nil
	}
	return *env.Data.Findings
}

func repoRootFor(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, serr := os.Stat(filepath.Join(dir, "go.mod")); serr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("no go.mod above the test working directory")
		}
		dir = parent
	}
}

// advisoryModel is a clean model carrying advisory findings only: one
// group closed by a bare escape row, and a redundant row whose accepted
// assignments are a proper subset of a sibling's.
const advisoryModel = `
outcomes = ["go", "stop"]
terminal = ["done"]

[model]
id = "lintfix"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true

[tags.flag]
provenance = "owned"
kind = "bool"
single_valued = true
required = true

[read.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"

[write.own]
role = "t"
path = "t.own"
keys = ["flag", "status"]
timeout = "2s"
read_back = true

[initial]
status = "a"
flag = "false"

[context.done]
[context.done.match.status]
eq = "b"

[[rule]]
id = "advance-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "advance-off"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "go"
[rule.guard.all.flag]
eq = "false"
[rule.write]
status = "b"

[[rule]]
id = "only-on"
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
[rule.guard.all.flag]
eq = "true"
[rule.write]
status = "b"

[[rule]]
id = "bare-rescue"
escape = ["no_match", "ambiguous_match"]
[rule.match.status]
eq = "a"
[rule.match.recognized]
eq = "stop"
`
