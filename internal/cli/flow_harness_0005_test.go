package cli

// RDR 0005 — the shared invocation harness for the `flow` suite.
//
// Two design rules hold throughout, and both are load-bearing for whether
// this suite actually discriminates the contract:
//
//  1. NOTHING here asserts on the caller-owned artifact's on-disk FORMAT.
//     The RDR fixes the CLI contract, not a file schema, so every state
//     oracle goes through the CLI: state is established with `flow
//     set-state` and observed with `flow read-state`. A round-trip proved
//     that way holds for whatever format the implementation chooses.
//
//  2. Every refusal oracle names its EXACT `flow-*` code. Asserting only
//     "an error occurred" would pass against an unregistered command
//     (cobra's own `command-error`), which is exactly the tautology the red
//     gate must exclude.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// flowInvocation is one named argv the suite drives through ExecuteAndEmit.
type flowInvocation struct {
	name string
	args []string
	// code is the exact `flow-*` code a refusal invocation must carry.
	code string
	// exit is the exact process exit code the refusal must map to.
	exit int
}

// seedArtifact establishes owned state on a fresh artifact by driving
// `flow set-state` through the production path, and returns the artifact
// path. It is how every read-side test gets an artifact whose content the
// CLI itself wrote — no format assumption anywhere.
//
// A seed that cannot be established is a FATAL setup failure, not a
// silently-skipped test: a suite that quietly degrades to "no artifact"
// would stop discriminating the read path.
func seedArtifact(t *testing.T, model string, writes ...string) string {
	t.Helper()

	art := newFlowArtifact(t, "state.artifact")
	args := []string{
		"flow", "set-state",
		"--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
	}
	for _, w := range writes {
		args = append(args, "--write", w)
	}
	if _, _, err := runCmd(t, append(args, "--as=json")...); err != nil {
		t.Fatalf("seeding the fixture artifact failed: %v", err)
	}
	return art
}

// flowSuccessInvocations returns one invocation per verb that the contract
// requires to SUCCEED. Each is self-contained: the artifact is seeded first
// so the read-side verbs have owned state to report.
//
// `next` and `resolve` deliberately bind ONLY the `state` role and leave
// `orphan` unbound (`0005:MVV`, REQ-36).
func flowSuccessInvocations(t *testing.T) []flowInvocation {
	t.Helper()

	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")
	stateBind := artifactBinding(flowStateRole, art)

	return []flowInvocation{
		{
			name: "next",
			args: []string{"flow", "next", "--model", model,
				"--artifact", stateBind, "--tag", "profile=mid"},
		},
		{
			name: "resolve",
			args: []string{"flow", "resolve", "--model", model,
				"--artifact", stateBind, "--outcome", "hold"},
		},
		{
			name: "read-state",
			args: []string{"flow", "read-state", "--model", model,
				"--artifact", stateBind,
				"--artifact", artifactBinding(flowOrphanRole, art)},
		},
		{
			name: "set-state",
			args: []string{"flow", "set-state", "--model", model,
				"--artifact", stateBind, "--write", "status=draft"},
		},
	}
}

// flowRefusalInvocations returns one invocation per refusal FAMILY the
// stable-code table names, each with the exact code and exit it must carry.
// Every entry is reachable from the CLI surface alone.
func flowRefusalInvocations(t *testing.T) []flowInvocation {
	t.Helper()

	model := writeFlowModel(t, flowMVVModel)
	invalid := writeFlowModel(t, flowInvalidModel)
	ambiguous := writeFlowModel(t, flowAmbiguousModel)
	art := seedArtifact(t, model, "status=draft")
	stateBind := artifactBinding(flowStateRole, art)
	ambigArt := seedArtifact(t, ambiguous, "status=draft")

	return []flowInvocation{
		{
			name: "model-not-found/neither",
			args: []string{"flow", "next"},
			code: "flow-model-not-found", exit: 2,
		},
		{
			name: "model-not-found/both",
			args: []string{"flow", "next", "--model", model, "--flow", "mvvflow"},
			code: "flow-model-not-found", exit: 2,
		},
		{
			name: "model-invalid",
			args: []string{"flow", "next", "--model", invalid},
			code: "flow-model-invalid", exit: 2,
		},
		{
			name: "tag-owned",
			args: []string{"flow", "resolve", "--model", model,
				"--artifact", stateBind, "--outcome", "hold",
				"--tag", "status=draft"},
			code: "flow-tag-owned", exit: 2,
		},
		{
			name: "tag-reserved",
			args: []string{"flow", "next", "--model", model,
				"--artifact", stateBind, "--tag", "recognized=hold"},
			code: "flow-tag-reserved", exit: 2,
		},
		{
			name: "tag-duplicate",
			args: []string{"flow", "next", "--model", model,
				"--artifact", stateBind,
				"--tag", "profile=mid", "--tag", "profile=foundational"},
			code: "flow-tag-duplicate", exit: 2,
		},
		{
			name: "tag-invalid/empty-outcome",
			args: []string{"flow", "resolve", "--model", model,
				"--artifact", stateBind, "--outcome", ""},
			code: "flow-tag-invalid", exit: 2,
		},
		{
			name: "artifact-invalid",
			args: []string{"flow", "next", "--model", model,
				"--artifact", "no-equals-sign"},
			code: "flow-artifact-invalid", exit: 2,
		},
		{
			name: "artifact-missing",
			args: []string{"flow", "resolve", "--model", model,
				"--outcome", "hold"},
			code: "flow-artifact-missing", exit: 2,
		},
		{
			name: "unmodeled-outcome",
			args: []string{"flow", "resolve", "--model", model,
				"--artifact", stateBind, "--outcome", "not-an-outcome"},
			code: "flow-unmodeled-outcome", exit: 2,
		},
		{
			name: "ambiguous-match",
			args: []string{"flow", "resolve", "--model", ambiguous,
				"--artifact", artifactBinding(flowStateRole, ambigArt),
				"--outcome", "advance"},
			code: "flow-ambiguous-match", exit: 2,
		},
		{
			name: "write-invalid/clear-sentinel",
			args: []string{"flow", "set-state", "--model", model,
				"--artifact", stateBind, "--write", "status=<clear>"},
			code: "flow-write-invalid", exit: 2,
		},
		{
			name: "write-duplicate",
			args: []string{"flow", "set-state", "--model", model,
				"--artifact", stateBind,
				"--write", "status=draft", "--write", "status=final"},
			code: "flow-write-duplicate", exit: 2,
		},
		{
			name: "write-unbound",
			args: []string{"flow", "set-state", "--model", model,
				"--artifact", stateBind, "--write", "note=x"},
			code: "flow-write-unbound", exit: 2,
		},
		{
			name: "clear-unbound",
			args: []string{"flow", "set-state", "--model", model,
				"--artifact", stateBind, "--clear", "note"},
			code: "flow-clear-unbound", exit: 2,
		},
	}
}

// --- oracles -------------------------------------------------------------

// requireRefusal drives args and asserts the refusal carries EXACTLY code
// and exit. It is the discriminating oracle the whole suite leans on: an
// unregistered command yields cobra's `command-error`, which fails here.
func requireRefusal(t *testing.T, code string, exit int, args ...string) *clierr.CLIError {
	t.Helper()

	stdout, _, err := runCmd(t, args...)
	if err == nil {
		t.Fatalf("invocation succeeded; want refusal %s\nstdout:\n%s",
			code, stdout)
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != code {
		t.Fatalf("refusal code = %q; want %q (the spellings are normative)",
			ce.Code, code)
	}
	if got := clierr.ExitCodeFor(err); got != exit {
		t.Errorf("exit code = %d; want %d — the group fixes the exit", got, exit)
	}
	return ce
}

// requireSuccess drives args and asserts it succeeded, returning stdout.
func requireSuccess(t *testing.T, args ...string) string {
	t.Helper()

	stdout, stderr, err := runCmd(t, args...)
	if err != nil {
		t.Fatalf("invocation failed: %v\nstdout:\n%s\nstderr:\n%s",
			err, stdout, stderr)
	}
	return stdout
}

// flowData decodes the `data` payload of a `--as=json` success line.
func flowData(t *testing.T, stdout string) map[string]any {
	t.Helper()

	var env struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if env.Type != "ok" {
		t.Fatalf("envelope type = %q; want %q", env.Type, "ok")
	}
	var data map[string]any
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("`data` is not a JSON object: %v\n%s", err, string(env.Data))
	}
	return data
}

// flowFinding is the WIRE shape of one `clierr.Finding` as this contract
// fixes it. It is decoded from JSON rather than reusing the Go type so that
// `param` and `locator` — the two fields REQ-14 / REQ-16 / A-7 ADD — are
// asserted as the envelope carries them, and so a missing Go field is a
// test FAILURE with a readable message rather than a package-wide compile
// break that would take RDR 0006's shipped suite down with it.
//
// REQ-15: "each producer populates only the fields it owns, and no producer
// nests its own fields in a sub-object" — every field here is FLAT.
type flowFinding struct {
	Code    string `json:"code"`
	Message string `json:"message"`

	// Param and Locator are this RDR's additions (REQ-18, A-7).
	Param   string `json:"param"`
	Locator string `json:"locator"`

	Hint     string `json:"hint"`
	Severity string `json:"severity"`
	Model    string `json:"model"`
	Rule     string `json:"rule"`

	// The guard atom's four fields sit FLAT on the record, never nested
	// under an `atom` object (REQ-15).
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Literal  string `json:"literal"`
	Block    string `json:"block"`

	Class string `json:"class"`
}

// flowFailureFindings decodes the `findings` array off a `--as=json`
// failure line. `findings` is a TOP-LEVEL sibling of `code` — EmitJSON
// marshals the *CLIError itself, so there is no `error` wrapper (REQ-8).
func flowFailureFindings(t *testing.T, stdout string) []flowFinding {
	t.Helper()

	var env struct {
		Code     string          `json:"code"`
		Findings []flowFinding   `json:"findings"`
		Error    json.RawMessage `json:"error"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); err != nil {
		t.Fatalf("stdout is not one JSON object: %v\n%s", err, stdout)
	}
	if len(env.Error) > 0 {
		t.Errorf("failure line carries an `error` wrapper; `findings` is a "+
			"TOP-LEVEL sibling of `code`:\n%s", stdout)
	}
	return env.Findings
}

// keysOf returns a payload map's keys, for failure messages that need to
// show what the payload actually carried.
func keysOf(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// stringsAt reads payload[key] as a []string, reporting whether it was
// present and array-shaped.
func stringsAt(m map[string]any, key string) ([]string, bool) {
	raw, ok := m[key]
	if !ok {
		return nil, false
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(arr))
	for _, e := range arr {
		s, ok := e.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// objectsAt reads payload[key] as a []map[string]any.
func objectsAt(m map[string]any, key string) ([]map[string]any, bool) {
	raw, ok := m[key]
	if !ok {
		return nil, false
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]map[string]any, 0, len(arr))
	for _, e := range arr {
		o, ok := e.(map[string]any)
		if !ok {
			return nil, false
		}
		out = append(out, o)
	}
	return out, true
}
