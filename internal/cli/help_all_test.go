package cli

// The `--help` contract: intrastate describes itself.
//
// The oracle here is not "the help text says something". It is that the
// help text is DERIVED from the same constants the wire is derived from,
// so a vocabulary change cannot land with a stale help body beside it.
// Each test below pins one way that derivation could rot:
//
//   - a command ships with no extended body at all;
//   - a finding or refusal code is renamed and the help keeps the old
//     spelling (the lists are generated from the taxonomy accessors, so
//     this test fails by finding the NEW code absent);
//   - --help-all stops exiting 0, or starts emitting an envelope on
//     stdout under --as=json where the never-silent contract reserves it;
//   - the scaffold placeholder comes back.

import (
	"bytes"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/spf13/cobra"
)

// runHelp captures a help invocation the way a user experiences it:
// through the production ExecuteAndEmit path, with both streams caught.
func runHelp(t *testing.T, args ...string) (stdout, stderr string) {
	t.Helper()
	cmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	if err := ExecuteAndEmit(cmd, args); err != nil {
		t.Fatalf("%v: unexpected error: %v", args, err)
	}
	return out.String(), errBuf.String()
}

// helpTargets is every command path a user can ask for help on. It is
// walked from the live tree rather than listed, so a verb added without
// help is caught by the tests below rather than by a reader.
func helpTargets(t *testing.T) [][]string {
	t.Helper()
	var paths [][]string
	walkCommandTree(NewRootCmd(), func(c *cobra.Command) {
		if c.Parent() == nil {
			paths = append(paths, nil)
			return
		}
		full := strings.Fields(c.CommandPath())
		paths = append(paths, full[1:])
	})
	return paths
}

// TestHelpAll_EveryCommandCarriesAnExtendedBody is the self-describing
// invariant itself. A verb whose only documentation is its flag list
// tells a caller what to type and nothing about what it means; this
// fails the moment such a verb is added.
func TestHelpAll_EveryCommandCarriesAnExtendedBody(t *testing.T) {
	walkCommandTree(NewRootCmd(), func(c *cobra.Command) {
		if extendedHelpFor(c) == "" {
			t.Errorf("%s carries no extended help body; register one with "+
				"withExtendedHelp so `--help-all` has something to say",
				c.CommandPath())
		}
		if c.Short == "" {
			t.Errorf("%s carries no Short", c.CommandPath())
		}
		if c.Long == "" {
			t.Errorf("%s carries no Long", c.CommandPath())
		}
	})
}

// TestHelpAll_ExitsZeroAndPrintsExtendedEverywhere pins that --help-all
// is a real flag on every command, not just the ones that happened to
// register it. A missing registration would surface as an unknown-flag
// refusal — a non-nil error out of runHelp.
func TestHelpAll_ExitsZeroAndPrintsExtendedEverywhere(t *testing.T) {
	for _, path := range helpTargets(t) {
		args := append(append([]string{}, path...), "--help-all")
		stdout, _ := runHelp(t, args...)
		if !strings.Contains(stdout, "Extended help:") {
			t.Errorf("`intrastate %s` printed no extended body:\n%s",
				strings.Join(args, " "), stdout)
		}
	}
}

// TestHelp_TerseStaysTerse pins the two-tier split. The default --help
// must NOT carry the extended body: the whole point of the tier is that
// the landing page fits a screen.
func TestHelp_TerseStaysTerse(t *testing.T) {
	for _, path := range helpTargets(t) {
		args := append(append([]string{}, path...), "--help")
		stdout, _ := runHelp(t, args...)
		if strings.Contains(stdout, "Extended help:") {
			t.Errorf("`intrastate %s` leaked the extended body into terse help",
				strings.Join(args, " "))
		}
	}
}

// TestRootHelp_HasNoScaffoldPlaceholder pins the specific regression
// this work fixed: the root shipped `<one-line description of what this
// tool does>` as its Long.
func TestRootHelp_HasNoScaffoldPlaceholder(t *testing.T) {
	stdout, _ := runHelp(t, "--help")
	if strings.Contains(stdout, "<one-line description") {
		t.Errorf("root --help still carries the scaffold placeholder:\n%s", stdout)
	}
	// The two surfaces a caller reaches for must both be named on the
	// landing page, or the CLI does not describe itself at the top.
	for _, want := range []string{"lint", "flow", "--as", "--help-all"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("root --help does not mention %q:\n%s", want, stdout)
		}
	}
}

// TestRootHelpAll_DumpsTheWholeTree pins the one-call self-description:
// `intrastate --help-all` and `intrastate help --all` must each carry
// every command, so an agent can read the surface in a single call.
func TestRootHelpAll_DumpsTheWholeTree(t *testing.T) {
	for _, args := range [][]string{{"--help-all"}, {"help", "--all"}} {
		stdout, _ := runHelp(t, args...)
		walkCommandTree(NewRootCmd(), func(c *cobra.Command) {
			if c.Parent() == nil {
				return
			}
			header := "## " + c.CommandPath()
			if !strings.Contains(stdout, header) {
				t.Errorf("`intrastate %s` omits %q",
					strings.Join(args, " "), header)
			}
		})
	}
}

// TestLintHelpAll_PublishesTheLiveTaxonomy is the anti-drift oracle for
// lint. The body is generated from graphlint's own accessors, so this
// passes by construction today — and fails loudly if someone later
// hand-writes the list back into the help text and a code is renamed.
func TestLintHelpAll_PublishesTheLiveTaxonomy(t *testing.T) {
	stdout, _ := runHelp(t, "lint", "--help-all")
	for _, code := range graphlint.BlockingCodes() {
		if !strings.Contains(stdout, code) {
			t.Errorf("lint --help-all omits blocking code %q", code)
		}
	}
	for _, code := range graphlint.AdvisoryCodes() {
		if !strings.Contains(stdout, code) {
			t.Errorf("lint --help-all omits advisory code %q", code)
		}
	}
	if !strings.Contains(stdout, graphlint.AggregateCode) {
		t.Errorf("lint --help-all omits the aggregate code %q",
			graphlint.AggregateCode)
	}
}

// TestResolveHelpAll_PublishesTheKernelMirroredCodes is the same oracle
// for the refusal vocabulary `flow resolve` can return. The codes are
// spelled through kernelCode — the same derivation the emitter uses — so
// a kernel refusal kind renamed on the wire renames itself here too.
func TestResolveHelpAll_PublishesTheKernelMirroredCodes(t *testing.T) {
	stdout, _ := runHelp(t, "flow", "resolve", "--help-all")
	for _, code := range []string{
		codeUnmodeledOutcome,
		codeNoMatch,
		codeAmbiguousMatch,
		codeOwnedStateUnavailable,
		codeGuardUnevaluable,
		codeGateDenied,
		codeGateIndeterminate,
	} {
		if !strings.Contains(stdout, code) {
			t.Errorf("flow resolve --help-all omits refusal code %q", code)
		}
	}
}

// TestHelpAll_UnderJSONModeEmitsNoEnvelope pins the boundary between
// help and the output gateway. Help is not a terminal result: it carries
// no envelope, and asking for help under --as=json must not fabricate
// one on stdout where the single terminal record is reserved.
func TestHelpAll_UnderJSONModeEmitsNoEnvelope(t *testing.T) {
	stdout, _ := runHelp(t, "flow", "next", "--help-all", "--as=json")
	if strings.Contains(stdout, `"type":"ok"`) ||
		strings.Contains(stdout, `"code":`) {
		t.Errorf("help under --as=json emitted an envelope:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Extended help:") {
		t.Errorf("help under --as=json printed no help:\n%s", stdout)
	}
}

// TestHelpAll_DoesNotSuppressUnknownCommandRefusals pins the Args
// hazard registerHelpAllOn documents: wrapping Args on a parent that had
// none flips cobra's unknown-command detection off, and `intrastate
// bogus` would print help at exit 0 instead of refusing. That would make
// a typo look like a success to a caller checking exit status.
func TestHelpAll_DoesNotSuppressUnknownCommandRefusals(t *testing.T) {
	cmd := NewRootCmd()
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	if err := ExecuteAndEmit(cmd, []string{"bogus"}); err == nil {
		t.Fatalf("`intrastate bogus` succeeded; want a structured refusal\n%s",
			out.String())
	}
	cmd = NewRootCmd()
	out.Reset()
	errBuf.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	if err := ExecuteAndEmit(cmd, []string{"flow", "bogus"}); err == nil {
		t.Fatalf("`intrastate flow bogus` succeeded; want a structured refusal\n%s",
			out.String())
	}
}
