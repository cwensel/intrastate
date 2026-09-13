package cli

// RDR 0021 — the verb surface (`0021:C1`): registration, the mirrored
// selection-arm set and its code spellings, `--emit` and its checking
// order, and the exit mapping.
//
// Every assertion drives ExecuteAndEmit, the production emission path, and
// reads the captured streams and the structured error — never an internal
// field. The verb does not exist yet, so each test fails on its own oracle
// (see graph_probe_0021_test.go for why that is a RUNTIME failure).

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// REQ-1: "A new root verb `graph` is registered beside `lint` — outside
// the `flow` group, under `0005:C1`'s carve-out for command groups \"owned
// by the RDR that names them\"."
// REQ-15: "verb `graph`, flag `--emit`."
// HAPPY PATH
func TestReq1And15_GraphIsRegisteredAtRootBesideLint(t *testing.T) {
	root := NewRootCmd()

	var haveGraph, haveLint bool
	for _, c := range root.Commands() {
		switch c.Name() {
		case graphVerb:
			haveGraph = true
		case "lint":
			haveLint = true
		}
		// The verb is deliberately NOT under the `flow` group: `0005:C1`
		// scopes that group to skill integration and carves out command
		// groups owned by the RDR that names them.
		if c.Name() == "flow" {
			for _, sub := range c.Commands() {
				if sub.Name() == graphVerb {
					t.Errorf("`%s` is registered under the `flow` group; "+
						"C1 registers it at ROOT, outside that group",
						graphVerb)
				}
			}
		}
	}
	if !haveLint {
		t.Fatal("no root `lint` command; the fixture for \"beside lint\" is gone")
	}
	if !haveGraph {
		var names []string
		for _, c := range root.Commands() {
			names = append(names, c.Name())
		}
		t.Fatalf("no root %q command is registered; root commands = %v",
			graphVerb, names)
	}
}

// REQ-15: "verb `graph`, flag `--emit`." — the rejected spellings are
// `dump` and `export`; neither may be registered as this verb's name.
// ADVERSARIAL — the defect is shipping the export under a name the
// decision rejected, which would merge two contracts (0002's `dump`) or
// name the act rather than the artifact (`export`).
func TestReq15_TheRejectedVerbSpellingsAreNotRegistered(t *testing.T) {
	requireGraphVerb(t)

	for _, rejected := range []string{"dump", "export"} {
		for _, c := range NewRootCmd().Commands() {
			if c.Name() == rejected {
				t.Errorf("a root %q command is registered; D-naming rejects "+
					"that spelling for this export (0002's dump is a "+
					"distinct rows-only rendering; `export` names the act, "+
					"not the artifact) — the verb is %q",
					rejected, graphVerb)
			}
		}
	}
}

// REQ-2: "Model selection mirrors `lint`'s arm set … and codes verbatim"
// REQ-3: "exactly one of `--model <path>` / `--flow <id>`; both →
// `flag-mutually-exclusive`"
// REQ-89: "Both/neither `--model`/`--flow` | 2 | `flag-mutually-exclusive`
// / `flag-required` | none | loud (C1)"
// INPUT EDGE
func TestReq2And3And89_BothSelectionFlagsIsMutuallyExclusive(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	stdout, _, err := runGraph(t, "--model", path, "--flow", "someid", "--as=json")

	assertRefusalCode(t, err, "flag-mutually-exclusive")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "both --model and --flow")
}

// REQ-4: "`--flow` alone → `flag-invalid-value` (this build resolves no
// ids)"
// REQ-90: "`--flow` alone | 2 | `flag-invalid-value` (no ids resolve this
// build) | none | loud"
// INPUT EDGE
func TestReq4And90_FlowAloneIsFlagInvalidValue(t *testing.T) {
	requireGraphVerb(t)

	stdout, _, err := runGraph(t, "--flow", "someid", "--as=json")

	assertRefusalCode(t, err, "flag-invalid-value")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "--flow alone")

	// The arm names the flag the caller actually used. lint's shipped arm
	// carries `Param: "flow"`, and C1 mirrors the code spellings verbatim.
	var ce *clierr.CLIError
	if errors.As(err, &ce) && ce.Param != "flow" {
		t.Errorf("the `--flow`-alone refusal names param %q; lint's mirrored "+
			"arm names %q — the honest refusal names the flag the caller "+
			"used, not one they did not", ce.Param, "flow")
	}
}

// REQ-5: "neither → `flag-required`"
// REQ-89 (neither half).
// INPUT EDGE
func TestReq5And89_NeitherSelectionFlagIsFlagRequired(t *testing.T) {
	requireGraphVerb(t)

	stdout, _, err := runGraph(t, "--as=json")

	assertRefusalCode(t, err, "flag-required")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "neither --model nor --flow")
}

// REQ-6: "unreadable file → `model-unreadable`"
// REQ-91: "Unreadable file / load failure | 2 | `model-unreadable` /
// `model-invalid` (one findings[] entry per load category) | none | loud"
// INPUT EDGE
func TestReq6And91_UnreadableModelFileIsModelUnreadable(t *testing.T) {
	requireGraphVerb(t)

	missing := filepath.Join(t.TempDir(), "no-such-model.toml")
	stdout, _, err := runGraph(t, "--model", missing, "--as=json")

	assertRefusalCode(t, err, "model-unreadable")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "unreadable model file")
}

// REQ-7: "load failure → `model-invalid` with one findings[] entry per
// load category."
// REQ-91 (load-failure half).
// INPUT EDGE — the model is READABLE and refuses at load, which is a
// different arm from REQ-6's unreadable file.
func TestReq7And91_LoadFailureIsModelInvalidWithPerCategoryFindings(t *testing.T) {
	requireGraphVerb(t)

	// A tag declaration the loader refuses: `set` admits no domain. The
	// same shape `load_locator_test.go` drives, so the category is the
	// loader's own rather than one this test invents.
	const badModel = `
outcomes = ["go"]

[model]
id = "loadfail"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
domain = ["go"]
single_valued = true
required = true

[tags.flavors]
provenance = "owned"
kind = "set"
domain = ["a", "b"]
`
	path := writeModel(t, badModel)
	stdout, _, err := runGraph(t, "--model", path, "--as=json")

	assertRefusalCode(t, err, "model-invalid")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "load failure")

	// The entries come from the shared loader `lint.go`'s arm also calls,
	// so each carries the load-category slug as its `code` — never the
	// envelope code (the Joint-check's reading, ASSUMPTION A-4).
	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if len(ce.Findings) == 0 {
		t.Fatalf("the `model-invalid` refusal carries no findings[]; C1 "+
			"mirrors lint's arm \"with one findings[] entry per load "+
			"category\"\nstdout: %s", stdout)
	}
	for _, f := range ce.Findings {
		if f.Code == "" {
			t.Errorf("a findings[] entry carries no code; each entry's code " +
				"is the load-category slug the shared loader populates")
		}
		if f.Code == "model-invalid" {
			t.Errorf("a findings[] entry carries the ENVELOPE code %q; the "+
				"entry code is the load-category slug, never the envelope's",
				f.Code)
		}
	}
}

// REQ-8: "`--emit <format>` selects the document: `json` (default) or
// `dot`; any other value → `flag-invalid-value` naming `emit`."
// REQ-92: "Unknown `--emit` value | 2 | `flag-invalid-value` naming `emit`
// | none | loud"
// INPUT EDGE
func TestReq8And92_UnknownEmitValueIsFlagInvalidValueNamingEmit(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	stdout, _, err := runGraph(t, "--model", path, "--emit", "xml", "--as=json")

	assertRefusalCode(t, err, "flag-invalid-value")
	assertExitTwo(t, err)
	assertNoDocumentOnStdout(t, stdout, "unknown --emit value")

	var ce *clierr.CLIError
	if errors.As(err, &ce) && ce.Param != emitFlag {
		t.Errorf("the unknown-format refusal names param %q; C1 requires it "+
			"NAME `emit` — a refusal that does not name the offending flag "+
			"leaves the caller guessing which of two flags was wrong",
			ce.Param)
	}
}

// REQ-8 (the two accepted members).
// HAPPY PATH — both declared formats are accepted, so the refusal above is
// discriminating rather than a blanket rejection.
func TestReq8_BothDeclaredEmitFormatsAreAccepted(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	for _, format := range []string{"json", "dot"} {
		if _, _, err := runGraph(t, "--model", path, "--emit", format); err != nil {
			t.Errorf("--emit %s refused: %v; C1 declares `json` (default) "+
				"and `dot` as the two members of the format set",
				format, err)
		}
	}
}

// REQ-8, ASSUMPTION A-1: "`--emit`'s default is applied, not merely
// documented."
// HAPPY PATH — invoking with no `--emit` produces the JSON document,
// byte-identical to an explicit `--emit=json`.
func TestReq8_TheEmitDefaultIsJSONAndIsApplied(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	implicit, _, err := runGraph(t, "--model", path)
	if err != nil {
		t.Fatalf("default emit: %v", err)
	}
	explicit, _, err := runGraph(t, "--model", path, "--emit", "json")
	if err != nil {
		t.Fatalf("--emit json: %v", err)
	}
	if implicit != explicit {
		t.Errorf("the default document differs from an explicit --emit=json;"+
			" C1 declares `json` the DEFAULT, so the two must be the same "+
			"bytes\n--- default ---\n%s\n--- explicit ---\n%s",
			implicit, explicit)
	}
}

// REQ-11: "`--emit` validity is checked with the argument-shaped arms,
// AFTER the `--model`/`--flow` selection arms and BEFORE any file I/O. So
// `graph --model <unreadable> --emit=xml` refuses `flag-invalid-value`
// naming `emit`, never `model-unreadable`; the request is wrong
// independent of the environment."
// ADVERSARIAL — the defect is an implementation that reads the file first,
// making a wrong REQUEST report as a wrong ENVIRONMENT.
func TestReq11_EmitValidityIsCheckedBeforeAnyFileIO(t *testing.T) {
	requireGraphVerb(t)

	missing := filepath.Join(t.TempDir(), "no-such-model.toml")
	_, _, err := runGraph(t, "--model", missing, "--emit", "xml", "--as=json")

	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if ce.Code == "model-unreadable" {
		t.Fatalf("`graph --model <unreadable> --emit=xml` refused %q; C1 "+
			"fixes the ORDER — `--emit` validity is checked BEFORE any file "+
			"I/O, so the answer is `flag-invalid-value` naming `emit`. The "+
			"request is wrong independent of the environment, and reporting "+
			"the environment sends the caller to fix the wrong thing",
			ce.Code)
	}
	if ce.Code != "flag-invalid-value" || ce.Param != emitFlag {
		t.Errorf("refusal = code %q param %q; want code %q naming %q",
			ce.Code, ce.Param, "flag-invalid-value", emitFlag)
	}
}

// REQ-11 (the other half of the order): `--emit` is checked AFTER the
// selection arms, so a bad `--emit` beside a bad SELECTION reports the
// selection.
// BOUNDARY
func TestReq11_SelectionArmsAreCheckedBeforeEmitValidity(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	_, _, err := runGraph(t, "--model", path, "--flow", "someid",
		"--emit", "xml", "--as=json")

	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != "flag-mutually-exclusive" {
		t.Errorf("both selection flags plus a bad --emit refused %q; C1 "+
			"orders the `--emit` check AFTER the `--model`/`--flow` "+
			"selection arms, so the selection arm answers first", ce.Code)
	}
}

// REQ-12: "This is the one value-checked flag on the verb"
// BOUNDARY — `--emit` is the only flag on this verb that rejects a value
// for being outside a declared set. The verb performs no second enum-value
// check of its own: `--as` is validated by the shared gateway
// (`respond.ValidateMode`), and the `--model`/`--flow` arms are PRESENCE
// arms — `--flow` refuses because no id resolves in this build, not
// because the value is outside an enumeration.
//
// The oracle is behavioural rather than a flag-table walk: it feeds each
// of the verb's OWN registered flags an arbitrary value and asserts only
// `--emit` answers `flag-invalid-value` naming itself. A structural walk
// filtering for `--emit` could only ever return 0 or 1 and would pass
// without discriminating anything.
func TestReq12_EmitIsTheOnlyValueCheckedFlagOnTheVerb(t *testing.T) {
	requireGraphVerb(t)

	// The flags the verb registers itself, excluding inherited/persistent
	// ones the gateway owns.
	var own []string
	lookupGraphCmd(t).Flags().VisitAll(func(f *pflag.Flag) {
		own = append(own, f.Name)
	})
	if !slices.Contains(own, emitFlag) {
		t.Fatalf("the verb registers no --%s flag; its own flags are %v",
			emitFlag, own)
	}

	path := writeModel(t, legalModel)
	for _, name := range own {
		if name == emitFlag {
			continue
		}
		// `--model` and `--flow` are the mirrored selection arms; feeding
		// them an arbitrary value must not produce an ENUM-value refusal
		// naming that flag as outside a declared set.
		args := []string{"--" + name, "zzz-not-a-declared-value", "--as=json"}
		if name != "model" && name != "flow" {
			args = append([]string{"--model", path}, args...)
		}
		_, _, err := runGraph(t, args...)

		var ce *clierr.CLIError
		if errors.As(err, &ce) && ce.Code == "flag-invalid-value" &&
			ce.Param == name && name != "flow" {
			t.Errorf("--%s rejected a value with %q naming itself; C1 fixes "+
				"`--emit` as THE one value-checked flag on this verb. "+
				"(`--flow` is exempt: its refusal is a presence arm — no id "+
				"resolves in this build — not an enum-value check.)",
				name, ce.Code)
		}
	}
}

// REQ-13: "RunE starts with `respond.ValidateMode`, success routes through
// `respond.OK`, failure through `respond.Fail`"
// ADVERSARIAL — an unrecognized `--as` value must refuse through the
// gateway's own validator before the verb does any work of its own.
func TestReq13_RunEStartsWithValidateMode(t *testing.T) {
	requireGraphVerb(t)

	path := writeModel(t, legalModel)
	_, _, err := runGraph(t, "--model", path, "--as=yaml")

	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("`--as=yaml` did not refuse through the gateway: %v", err)
	}
	if ce.Code != "flag-invalid-value" || ce.Param != "as" {
		t.Errorf("`--as=yaml` refused code %q param %q; RunE starts with "+
			"respond.ValidateMode, whose refusal names `as`",
			ce.Code, ce.Param)
	}
}

// REQ-14: "exit codes are the existing 0/2 mapping — no new exit group
// (`0005:C1`: exit 3 is environment-not-consulted only, and an unreadable
// model file is a wrong request, so it stays exit 2)."
// ADVERSARIAL — the tempting defect is reading an unreadable file as an
// ENVIRONMENT failure (exit 3), which would make a caller's retry loop
// spin on an input it must instead fix.
func TestReq14_TheVerbUsesOnlyTheExistingZeroTwoExitMapping(t *testing.T) {
	requireGraphVerb(t)

	missing := filepath.Join(t.TempDir(), "no-such-model.toml")
	okPath := writeModel(t, legalModel)

	for _, tc := range []struct {
		name string
		args []string
		want int
	}{
		{"clean model", []string{"--model", okPath}, 0},
		{"unreadable file", []string{"--model", missing}, 2},
		{"neither flag", nil, 2},
		{"unknown emit", []string{"--model", okPath, "--emit", "xml"}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := runGraph(t, tc.args...)
			if got := clierr.ExitCodeFor(err); got != tc.want {
				t.Errorf("exit = %d; want %d. C1 keeps the existing 0/2 "+
					"mapping and adds no exit group — exit 3 is "+
					"environment-not-consulted only, and a wrong request "+
					"stays exit 2", got, tc.want)
			}
		})
	}
}

// REQ-110: "The export is a new root verb beside `lint`, adding no flag to
// `lint` and altering no existing output."
// ADVERSARIAL — the whole point of the sibling-verb shape is that lint's
// surface is untouched; a `--emit` that appeared on `lint` would be
// alternative O2, which the matrix rejected.
func TestReq110_LintGainsNoFlagFromThisRecord(t *testing.T) {
	var lintCmd = lookupCmd(t, "lint")

	if f := lintCmd.Flags().Lookup(emitFlag); f != nil {
		t.Errorf("`lint` registers a --%s flag; this record adds NO flag to "+
			"lint (XC, REQ-110). Putting the format selector on the "+
			"authoritative acceptance verb is alternative O2, rejected "+
			"because it couples the export to lint's verdict", emitFlag)
	}

	// lint's own output is unchanged: a clean model still emits the same
	// success envelope with an empty findings receipt.
	path := writeModel(t, legalModel)
	stdout, _, err := runCmd(t, "lint", "--model", path, "--as=json")
	if err != nil {
		t.Fatalf("lint over the clean fixture refused: %v", err)
	}
	env := parseSuccess(t, stdout)
	if env.Type != "ok" || env.Data == nil || env.Data.Findings == nil {
		t.Errorf("lint's success envelope changed shape: %s", stdout)
	}
}

// --- shared assertions ---------------------------------------------------

// assertRefusalCode asserts the error is a structured CLIError carrying
// code. C1 mirrors lint's code spellings VERBATIM, so each is compared as
// an exact string.
func assertRefusalCode(t *testing.T, err error, code string) {
	t.Helper()

	if err == nil {
		t.Fatalf("the invocation succeeded; want a %q refusal", code)
	}
	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the refusal is not a structured CLIError: %v", err)
	}
	if ce.Code != code {
		t.Errorf("refusal code = %q; want %q — C1 mirrors lint's arm set "+
			"and its code spellings verbatim (`internal/cli/lint.go::runLint`)",
			ce.Code, code)
	}
}

// assertExitTwo asserts the refusal maps to the existing exit 2.
func assertExitTwo(t *testing.T, err error) {
	t.Helper()

	if got := clierr.ExitCodeFor(err); got != 2 {
		t.Errorf("exit = %d; want 2 (GroupUserEnv). C1 adds no exit group "+
			"and exit 3 is environment-not-consulted only", got)
	}
}

// assertNoDocumentOnStdout asserts a refusing arm put NO document on
// stdout — C4's never-a-partial-document rule, which REQ-83 extends to
// EVERY refusing arm rather than only the ceiling one.
func assertNoDocumentOnStdout(t *testing.T, stdout, arm string) {
	t.Helper()

	if strings.Contains(stdout, schemaMarker) {
		t.Errorf("stdout carries the document marker %q on the %s arm; "+
			"C4 refuses rather than emitting a partial document, and a "+
			"refusing arm emits none at all\nstdout: %s",
			schemaMarker, arm, stdout)
	}
	if strings.Contains(stdout, "digraph") {
		t.Errorf("stdout carries a DOT document on the %s arm; a refusing "+
			"arm emits no document (C4, REQ-83)\nstdout: %s", arm, stdout)
	}
}

// --- command lookup ------------------------------------------------------

// lookupCmd returns a registered root command by name, failing if absent.
func lookupCmd(t *testing.T, name string) *cobra.Command {
	t.Helper()

	for _, c := range NewRootCmd().Commands() {
		if c.Name() == name {
			return c
		}
	}
	t.Fatalf("no root %q command is registered", name)
	return nil
}

// lookupGraphCmd returns the `graph` command, failing if absent.
func lookupGraphCmd(t *testing.T) *cobra.Command {
	t.Helper()
	return lookupCmd(t, graphVerb)
}

// writeModelAt writes src under dir with the given name and returns the
// path. Used where a test needs two models in one directory.
func writeModelAt(t *testing.T, dir, name, src string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
	return path
}
