package cli

// RDR 0005 — the command group and verb surface (REQ-1..REQ-4), the output
// gateway and envelope (REQ-5..REQ-13), and the ownership fences that keep
// this contract from absorbing its peers (REQ-2, REQ-10, REQ-113).
//
// Every assertion drives ExecuteAndEmit — the production emission path — and
// reads the captured streams or the returned *clierr.CLIError. No test here
// reaches into a renderer, a kernel function, or an unexported field.

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// flowVerbs is the exact verb set REQ-1 and REQ-3 fix.
var flowVerbs = []string{"next", "resolve", "read-state", "set-state"}

// flowSubcommand returns the `flow` group's subcommand named name, or nil.
func flowSubcommand(t *testing.T, name string) bool {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		for _, sub := range c.Commands() {
			if sub.Name() == name {
				return true
			}
		}
	}
	return false
}

// flowGroupNames lists the `flow` group's registered subcommand names.
func flowGroupNames(t *testing.T) []string {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() != "flow" {
			continue
		}
		var names []string
		for _, sub := range c.Commands() {
			names = append(names, sub.Name())
		}
		slices.Sort(names)
		return names
	}
	return nil
}

// REQ-1: "The CLI MUST expose one command group for skill integration with
// these verbs: next, resolve, read-state, and set-state."
// REQ-3: "the user-facing group is `flow`, with verbs `next`, `resolve`,
// `read-state`, and `set-state`" — the exact spellings are normative;
// `state` and `run` are named REJECTED group names in the same decision.
// HAPPY PATH
func TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs(t *testing.T) {
	got := flowGroupNames(t)
	if got == nil {
		var roots []string
		for _, c := range NewRootCmd().Commands() {
			roots = append(roots, c.Name())
		}
		t.Fatalf("no `flow` command group is registered; root commands = %v", roots)
	}

	want := slices.Clone(flowVerbs)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("flow group verbs = %v; want exactly %v — the spellings are "+
			"normative (`0005:D-naming`)", got, want)
	}

	// The rejected group names must NOT be registered: `state` and `run` are
	// named rejected alternatives, not aliases.
	for _, c := range NewRootCmd().Commands() {
		if c.Name() == "state" || c.Name() == "run" {
			t.Errorf("group %q is registered; `state` and `run` are named "+
				"REJECTED group names (`0005:D-naming`)", c.Name())
		}
	}
}

// REQ-2: "Other command groups (lint, dump, parse) are outside this contract
// and are owned by the RDR that names them." — the `flow` group MUST NOT
// absorb the shipped `lint` verb (RDR 0006) or `dump` (RDR 0002).
// ADVERSARIAL
func TestReq2_FlowGroupDoesNotAbsorbLintDumpOrParse(t *testing.T) {
	// The fence only says something once the group EXISTS: "the `flow`
	// group does not contain `lint`" is vacuously true of a group that is
	// not registered at all, and a vacuous pass would let this clause go
	// unimplemented unnoticed.
	if flowGroupNames(t) == nil {
		t.Fatal("no `flow` command group is registered; the ownership fence " +
			"this REQ states is not yet assertable")
	}

	for _, foreign := range []string{"lint", "dump", "parse"} {
		if flowSubcommand(t, foreign) {
			t.Errorf("`flow %s` is registered; that verb is owned by the RDR "+
				"that names it, not by this contract", foreign)
		}
	}

	// `lint` stays reachable at ROOT — absorbing it would be the drift.
	var rootLint bool
	for _, c := range NewRootCmd().Commands() {
		if c.Name() == "lint" {
			rootLint = true
		}
	}
	if !rootLint {
		t.Error("root `lint` disappeared; RDR 0006 owns it at root and this " +
			"contract must not move it")
	}
}

// REQ-3: the flag surface is normative — `--flow`, `--model`, `--tag`,
// `--artifact`, `--outcome`, `--evaluate-gates`, `--write`, `--clear`.
// BOUNDARY
func TestReq3_EachVerbRegistersItsNormativeFlagSpellings(t *testing.T) {
	// Which verb owns which flag, per the Technical Design. A flag absent
	// where the contract places it is a missing surface; the assertion is on
	// PRESENCE, not on any implementation detail of parsing.
	want := map[string][]string{
		"next":       {"flow", "model", "tag", "artifact", "evaluate-gates"},
		"resolve":    {"flow", "model", "tag", "artifact", "outcome"},
		"read-state": {"flow", "model", "artifact"},
		"set-state":  {"flow", "model", "artifact", "write", "clear"},
	}

	root := NewRootCmd()
	var group bool
	for _, c := range root.Commands() {
		if c.Name() != "flow" {
			continue
		}
		group = true
		for _, sub := range c.Commands() {
			flags, ok := want[sub.Name()]
			if !ok {
				continue
			}
			for _, f := range flags {
				if sub.Flags().Lookup(f) == nil {
					t.Errorf("`flow %s` does not register --%s; the flag "+
						"spellings are normative (`0005:D-naming`)",
						sub.Name(), f)
				}
			}
		}
	}
	if !group {
		t.Fatal("no `flow` command group is registered")
	}
}

// REQ-5: "All four verbs MUST start RunE by calling
// respond.ValidateMode(cmd), MUST route success through respond.OK, MUST
// route user-facing failure through respond.Fail(cmd, *clierr.CLIError), and
// MUST set SilenceErrors and SilenceUsage."
// ADVERSARIAL — an unrecognized `--as` must be refused by ValidateMode
// BEFORE any verb-specific work, on every verb.
func TestReq5_EveryVerbValidatesOutputModeFirst(t *testing.T) {
	model := writeFlowModel(t, flowMVVModel)

	for _, verb := range flowVerbs {
		t.Run(verb, func(t *testing.T) {
			// `--as=yaml` is invalid. Because ValidateMode runs FIRST, the
			// refusal must be the mode refusal even though the rest of the
			// invocation is deliberately also wrong (no --outcome, no
			// --artifact): the mode error is what wins.
			_, _, err := runCmd(t, "flow", verb, "--model", model, "--as=yaml")
			if err == nil {
				t.Fatal("`--as=yaml` was accepted; ValidateMode refuses it")
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("refusal is not a structured CLIError: %v", err)
			}
			if ce.Code != "flag-invalid-value" {
				t.Errorf("code = %q; want %q — ValidateMode runs FIRST in "+
					"RunE, before any verb-specific refusal", ce.Code,
					"flag-invalid-value")
			}
		})
	}
}

// REQ-5: SilenceErrors and SilenceUsage are set, so cobra's plain-text
// error and usage block never stack above the structured envelope.
// ADVERSARIAL
func TestReq5_VerbsSilenceCobraErrorAndUsageOutput(t *testing.T) {
	if names := flowGroupNames(t); len(names) != len(flowVerbs) {
		t.Fatalf("the `flow` group registers %v; all four verbs must exist "+
			"before their Silence settings are assertable", names)
	}

	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() != "flow" {
			continue
		}
		if !c.SilenceErrors || !c.SilenceUsage {
			t.Errorf("`flow` group: SilenceErrors=%v SilenceUsage=%v; both "+
				"MUST be set", c.SilenceErrors, c.SilenceUsage)
		}
		for _, sub := range c.Commands() {
			if !sub.SilenceErrors || !sub.SilenceUsage {
				t.Errorf("`flow %s`: SilenceErrors=%v SilenceUsage=%v; both "+
					"MUST be set", sub.Name(), sub.SilenceErrors,
					sub.SilenceUsage)
			}
		}
	}
}

// REQ-6: "Under --as=json, each successful invocation MUST emit exactly one
// stdout JSON terminal envelope with type \"ok\" and verb-specific data."
// HAPPY PATH
func TestReq6_JSONSuccessIsExactlyOneOKEnvelopeWithVerbSpecificData(t *testing.T) {
	for _, tc := range flowSuccessInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, append(tc.args, "--as=json")...)
			if err != nil {
				t.Fatalf("invocation failed: %v\nstdout:\n%s", err, stdout)
			}
			lines := nonEmptyLines(stdout)
			if len(lines) != 1 {
				t.Fatalf("stdout carries %d lines; a successful invocation "+
					"emits EXACTLY one terminal envelope:\n%s",
					len(lines), stdout)
			}
			var env struct {
				Type string          `json:"type"`
				Data json.RawMessage `json:"data"`
			}
			if uerr := json.Unmarshal([]byte(lines[0]), &env); uerr != nil {
				t.Fatalf("stdout is not one JSON object: %v\n%s", uerr, stdout)
			}
			if env.Type != "ok" {
				t.Errorf("envelope type = %q; want %q", env.Type, "ok")
			}
			if len(env.Data) == 0 || string(env.Data) == "null" {
				t.Errorf("envelope carries no `data`; every successful "+
					"invocation carries verb-specific data:\n%s", lines[0])
			}
		})
	}
}

// REQ-7: "Under --as=text, each successful invocation MUST emit human output
// derived from the same verb-specific result."
// REQ-10: "New verbs must not use direct Cobra printing for text payloads."
// HAPPY PATH
func TestReq7And10_TextSuccessEmitsHumanOutputFromTheSameResult(t *testing.T) {
	for _, tc := range flowSuccessInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _, err := runCmd(t, tc.args...)
			if err != nil {
				t.Fatalf("invocation failed: %v", err)
			}
			if strings.TrimSpace(stdout) == "" {
				t.Fatal("text mode emitted NO stdout payload; each successful " +
					"invocation emits human output derived from the same " +
					"verb-specific result")
			}
			// Text mode is not the JSON envelope leaking through.
			if strings.HasPrefix(strings.TrimSpace(stdout), `{"type"`) {
				t.Errorf("text mode emitted the JSON envelope:\n%s", stdout)
			}
		})
	}
}

// REQ-9: "Add text success payload rendering through `respond.OK` or a
// respond-owned helper used by `respond.OK`; do not print directly from
// resolver verbs."
// REQ-10: "New verbs must not use direct Cobra printing for text payloads."
// — `internal/cli::newVersionCmd`'s `cmd.Println` is the named drift.
// ADVERSARIAL — the observable proof that the gateway owns the stream: the
// payload rides the writer the gateway was handed, and swapping the mode
// swaps the rendering wholesale. A verb printing directly with `cmd.Println`
// would still write on `--as=json`, putting a second record on stdout.
func TestReq9And10_TextPayloadsRouteThroughTheGatewayNotDirectPrinting(t *testing.T) {
	for _, tc := range flowSuccessInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			// Under --as=json the gateway emits EXACTLY one terminal record.
			// A direct `cmd.Println` of a text payload is mode-blind and
			// would add a second line here.
			jsonOut := requireSuccess(t, append(tc.args, "--as=json")...)
			if n := len(nonEmptyLines(jsonOut)); n != 1 {
				t.Fatalf("--as=json wrote %d stdout lines; the gateway emits "+
					"one terminal record and the verb prints nothing "+
					"directly — a mode-blind `cmd.Println` is exactly this "+
					"defect:\n%s", n, jsonOut)
			}
			if !strings.HasPrefix(strings.TrimSpace(jsonOut), "{") {
				t.Errorf("--as=json stdout is not the JSON envelope; a text "+
					"payload leaked past the gateway:\n%s", jsonOut)
			}

			// Under --as=text the payload IS rendered — so the routing did
			// not simply drop it.
			textOut := requireSuccess(t, tc.args...)
			if strings.TrimSpace(textOut) == "" {
				t.Error("--as=text rendered no payload; the text branch " +
					"must render the verb's result through respond.OK, not " +
					"drop it")
			}
		})
	}
}

// REQ-11: "Text mode renders the same result content in a human-scannable
// order; it does not invent fields absent from the JSON payload".
// REQ-120: "Derive both renderings from one typed result per verb and test
// both modes."
// DOMAIN EDGE — the discriminating oracle: every SCALAR value the JSON
// payload carries must be findable in the text rendering. A text renderer
// built from a second, drifting source fails this.
func TestReq11And120_TextRenderingCarriesEveryJSONScalarValue(t *testing.T) {
	for _, tc := range flowSuccessInvocations(t) {
		t.Run(tc.name, func(t *testing.T) {
			jsonOut, _, jerr := runCmd(t, append(tc.args, "--as=json")...)
			if jerr != nil {
				t.Fatalf("json invocation failed: %v", jerr)
			}
			textOut, _, terr := runCmd(t, tc.args...)
			if terr != nil {
				t.Fatalf("text invocation failed: %v", terr)
			}

			var env struct {
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal([]byte(strings.TrimSpace(jsonOut)), &env); err != nil {
				t.Fatalf("json stdout not parseable: %v\n%s", err, jsonOut)
			}
			var payload any
			if err := json.Unmarshal(env.Data, &payload); err != nil {
				t.Fatalf("data not parseable: %v", err)
			}
			for _, v := range scalarStrings(payload) {
				if v == "" {
					continue
				}
				if !strings.Contains(textOut, v) {
					t.Errorf("text mode omits the value %q that the JSON "+
						"payload carries; both renderings derive from ONE "+
						"typed result per verb\ntext:\n%s", v, textOut)
				}
			}
		})
	}
}

// REQ-12: "The `--as` flag itself — its persistence at the root, its
// `text|json` domain, and its default — stays the root contract's … this RDR
// inherits it unchanged and does not redefine it for the `flow` group."
// ADVERSARIAL
func TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag(t *testing.T) {
	root := NewRootCmd()
	if root.PersistentFlags().Lookup("as") == nil {
		t.Fatal("the root `--as` persistent flag is gone")
	}
	for _, c := range root.Commands() {
		if c.Name() != "flow" {
			continue
		}
		if c.Flags().Lookup("as") != nil || c.PersistentFlags().Lookup("as") != nil {
			t.Error("`flow` re-registers `--as`; the flag stays the ROOT " +
				"contract's and is inherited unchanged")
		}
		for _, sub := range c.Commands() {
			if sub.Flags().Lookup("as") != nil {
				t.Errorf("`flow %s` re-registers `--as`; it is inherited "+
					"unchanged from the root", sub.Name())
			}
		}
	}

	// The default is still text: an invocation with no `--as` renders text.
	for _, tc := range flowSuccessInvocations(t) {
		stdout, _, err := runCmd(t, tc.args...)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if strings.HasPrefix(strings.TrimSpace(stdout), "{") {
			t.Errorf("%s: the default mode emitted JSON; the root default is "+
				"text and this RDR inherits it", tc.name)
		}
		break
	}
}

// REQ-13: "this RDR introduces verb-specific JSON `data` payloads and
// exactly one new `omitempty` `CLIError` field (`findings`), not a new
// terminal envelope or exit group."
// REQ-19: "No new exit group."
// BOUNDARY
func TestReq13And19_NoNewTerminalEnvelopeTypeAndNoNewExitGroup(t *testing.T) {
	// The success envelope's discriminator stays "ok"; the failure
	// envelope's stays the CLIError shape. No third `type` value appears.
	for _, tc := range flowSuccessInvocations(t) {
		stdout, _, err := runCmd(t, append(tc.args, "--as=json")...)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		var env struct {
			Type string `json:"type"`
		}
		if uerr := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &env); uerr != nil {
			t.Fatalf("%s: %v", tc.name, uerr)
		}
		if env.Type != "ok" {
			t.Errorf("%s: terminal type = %q; this RDR introduces NO new "+
				"terminal envelope", tc.name, env.Type)
		}
	}

	// Every refusal this contract raises exits 2 or 3 — the two shipped
	// codes. A new exit group would show up as a third value.
	for _, tc := range flowRefusalInvocations(t) {
		_, _, err := runCmd(t, append(tc.args, "--as=json")...)
		if err == nil {
			t.Errorf("%s: expected a refusal", tc.name)
			continue
		}
		switch code := clierr.ExitCodeFor(err); code {
		case 2, 3:
		default:
			t.Errorf("%s: exit code = %d; this RDR adds NO new exit group — "+
				"every failure is 2 or 3", tc.name, code)
		}
	}
}

// scalarStrings walks a decoded JSON payload and returns every string and
// number value it carries, as text. Map KEYS are excluded: they are the
// payload's field names, and REQ-11 constrains rendered CONTENT, leaving
// text layout free.
func scalarStrings(v any) []string {
	var out []string
	switch t := v.(type) {
	case string:
		out = append(out, t)
	case bool:
		if t {
			out = append(out, "true")
		} else {
			out = append(out, "false")
		}
	case []any:
		for _, e := range t {
			out = append(out, scalarStrings(e)...)
		}
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			out = append(out, scalarStrings(t[k])...)
		}
	}
	return out
}
