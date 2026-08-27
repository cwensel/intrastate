package cli

// RDR 0005 — the `flow` GROUP node itself, as distinct from its four verbs.
//
// REQ-5 binds the verbs; the group node sits outside it. But the group's own
// declared intent — "a bare `flow` with no verb is a usage error, not a
// silent success" — is a claim about observable behaviour, and the
// never-silent contract (root.go's package doc, REQ-6/REQ-113) says every
// graceful exit emits exactly one terminal record on the mode's channel.
//
// The asymmetry these tests close: `flow bogusverb` already refuses through
// `command-error` / exit 2, while a bare `flow` exited 0 with raw help on
// stdout — under `--as=json` too, where stdout is reserved for the single
// terminal envelope. No new code is minted here: `command-error` is the
// group-level usage code root.go already ships for the adjacent case.

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// A bare `flow` carries no behaviour of its own, so it is a usage refusal
// with the same identity and exit code as any other cobra-level usage error.
// ADVERSARIAL — this exited 0 with no envelope at all before the fix.
func TestFlowGroup_BareInvocationIsAUsageRefusalNotASilentSuccess(t *testing.T) {
	_, _, err := runCmd(t, "flow")
	if err == nil {
		t.Fatal("a bare `flow` succeeded; the group carries no behaviour of " +
			"its own, so invoking it without a verb is a usage error — " +
			"exiting 0 tells a script the request was honoured")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != "command-error" {
		t.Errorf("code = %q; want %q — the group-level usage code the "+
			"sibling case (`flow <unknown-verb>`) already carries",
			ce.Code, "command-error")
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2 — a usage error is a user-input refusal, "+
			"the `GroupUserEnv` bucket", got)
	}
	// The message must name the remedy: which verbs exist. A refusal that
	// says only "usage error" leaves the caller where it found them.
	for _, verb := range flowVerbs {
		if !strings.Contains(ce.Message+ce.Hint, verb) {
			t.Errorf("neither message %q nor hint %q names the verb %q; the "+
				"refusal must tell the caller what to run instead",
				ce.Message, ce.Hint, verb)
		}
	}
}

// Under `--as=json` stdout carries exactly ONE terminal record and nothing
// else. Raw help text on stdout is not a record, and a consumer that reads
// stdout as NDJSON cannot parse it.
// ADVERSARIAL — before the fix stdout held the full help block, unparseable.
func TestFlowGroup_BareInvocationUnderJSONEmitsTheEnvelopeNotHelpText(t *testing.T) {
	stdout, _, err := runCmd(t, "flow", "--as=json")
	if err == nil {
		t.Fatal("a bare `flow --as=json` succeeded; want a usage refusal")
	}

	line := strings.TrimSpace(stdout)
	if line == "" {
		t.Fatal("stdout is empty under `--as=json`; the never-silent " +
			"contract requires exactly one terminal record per graceful exit")
	}
	if n := strings.Count(line, "\n"); n != 0 {
		t.Errorf("stdout carries %d newlines, so more than one record; "+
			"`--as=json` stdout is exactly one terminal line.\ngot: %q",
			n+1, stdout)
	}
	// The help block's own words. Their presence on stdout is the defect.
	for _, fragment := range []string{"Usage:", "Available Commands:",
		"Global Flags:"} {
		if strings.Contains(stdout, fragment) {
			t.Errorf("stdout under `--as=json` contains the help-text "+
				"fragment %q; help is human output and must never occupy "+
				"the machine channel.\ngot: %q", fragment, stdout)
		}
	}

	var env struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if uerr := json.Unmarshal([]byte(line), &env); uerr != nil {
		t.Fatalf("stdout is not one JSON object: %v\ngot: %q", uerr, stdout)
	}
	if env.Code != "command-error" {
		t.Errorf("envelope code = %q; want %q", env.Code, "command-error")
	}
}

// ValidateMode runs FIRST, before the group's own usage refusal — matching
// the ordering REQ-5 fixes for the four verbs. An unrecognized `--as` is
// what the caller must fix first: they cannot even read the other refusal.
// ADVERSARIAL — the invocation is deliberately ALSO wrong (no verb), so the
// arm discriminates ordering rather than merely "some error happened".
func TestFlowGroup_InvalidOutputModeIsRefusedBeforeTheUsageError(t *testing.T) {
	_, _, err := runCmd(t, "flow", "--as=bogus")
	if err == nil {
		t.Fatal("`flow --as=bogus` succeeded; ValidateMode refuses an " +
			"unrecognized output mode")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != "flag-invalid-value" {
		t.Errorf("code = %q; want %q — the bare invocation is also a usage "+
			"error, and the MODE refusal must win because it runs first, "+
			"exactly as it does on all four verbs", ce.Code,
			"flag-invalid-value")
	}
	if ce.Param != "as" {
		t.Errorf("param = %q; want %q", ce.Param, "as")
	}
}

// The sibling case the fix must not disturb: an unknown verb under the group
// still refuses with the SAME code, so the two usage errors stay symmetric.
// CONTROL — this passed before the fix and must keep passing.
func TestFlowGroup_UnknownVerbStillRefusesWithTheSameUsageCode(t *testing.T) {
	_, _, err := runCmd(t, "flow", "bogusverb")
	if err == nil {
		t.Fatal("`flow bogusverb` succeeded; an unknown verb is a usage error")
	}
	if code := clierr.ErrorCode(err); code != "command-error" {
		t.Errorf("code = %q; want %q", code, "command-error")
	}
	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2", got)
	}
}

// CONTROL — each verb still runs. A group-node refusal that leaked onto the
// verbs would break every real invocation, and `Args: cobra.NoArgs` plus a
// RunE on the parent is exactly the shape where that mistake hides.
func TestFlowGroup_TheVerbsThemselvesStillRun(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)
	art := seedArtifact(t, model, "status=draft")

	stdout := requireSuccess(t, "flow", "read-state", "--model", model,
		"--artifact", artifactBinding(flowStateRole, art),
		"--artifact", artifactBinding(flowOrphanRole, art), "--as=json")
	if flowData(t, stdout) == nil {
		t.Error("`flow read-state` emitted no data payload; the group-node " +
			"refusal must not reach the verbs")
	}
}
