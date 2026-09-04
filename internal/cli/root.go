// Package cli wires the Cobra command tree and owns the boundary
// between cobra's error handling and the structured-error gateway. The
// binary in cmd/intrastate is a thin shim over Execute.
//
// Conventions for new verbs (so future prompts have fewer decisions):
//
//   - Add the verb with cmd.AddCommand(newXxxCmd()) in NewRootCmd.
//   - In each RunE, call respond.ValidateMode(cmd) first, then route
//     every success through respond.OK and every failure through
//     respond.Fail(cmd, &clierr.CLIError{…}). Never print to stdout or
//     stderr directly — the output gateway owns both streams.
//   - Use SilenceErrors + SilenceUsage on every command so cobra's
//     plain-text errors don't stack above the structured envelope;
//     ExecuteAndEmit converts any cobra-level error into a CLIError.
package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
	"github.com/cwensel/intrastate/internal/version"
	"github.com/spf13/cobra"
)

// rootLongDesc is the terse default body: what the tool is, the two
// surfaces a caller reaches for, and a worked invocation of each. The
// vocabulary (tags, provenance, models, artifacts) and the wire/exit
// contract live in rootExtendedDesc, surfaced on --help-all, so the
// default landing page stays one screen.
const rootLongDesc = `intrastate makes workflow state transitions explicit, reviewable, and
deterministic. A flow is authored once as a transition model — a TOML
document of tags, rules, guards, and writes — and every caller navigates
it by asking this CLI, rather than reimplementing the flow in a skill,
script, or agent.

Two surfaces:

  lint    check a model against the graph invariants, before runtime
  flow    drive a model at runtime: what can happen, what does happen,
          and what state was written

Author a model, then check it:

  intrastate lint --model flow.toml

Ask what the current state can do, then commit one outcome:

  intrastate flow next     --model flow.toml --artifact state=state.json
  intrastate flow resolve  --model flow.toml --artifact state=state.json \
      --outcome approved
  intrastate flow set-state --model flow.toml --artifact state=state.json \
      --write status=approved

Given the same model, state, tags, and recognized outcome, the answer is
the same every time: one legal plan, or one typed refusal. intrastate
never guesses which of two matching rules you meant.

Global flags:

  --as text|json    output mode (default text)
  --help-all        the extended reference for any command

Run any subcommand with --help for its flags, or --help-all for its
extended reference. ` + "`intrastate help --all`" + ` prints the full reference for
every command at once.`

// rootExtendedDesc is the vocabulary and contract a caller needs to
// read this CLI's output — surfaced on --help-all. It describes what
// the code does; docs/ carries the longer-form authoring guidance.
const rootExtendedDesc = `Vocabulary

  model       one TOML document declaring the tags, rules, guards, and
              writes of a flow. Its class is either state-machine (it
              owns state and advances it) or decision-table (it owns no
              state and maps a situation to an answer).
  tag         one named piece of state. Its PROVENANCE says who supplies
              it, and that is what decides which flag carries it:
                owned      the model's own state, read and written
                           through accessors. Never settable via --tag;
                           --write is the only authoring channel.
                observed   caller context, supplied as --tag name=value.
                recognized the outcome being asked about, supplied as
                           --outcome on flow resolve.
  rule        one row of the model: a match, optional guards, an
              optional write block, and (for a decision table) an emit
              block of answer values.
  candidate   a rule the supplied state does not exclude. Not a rule
              that will be selected — that is flow resolve's verdict.
  plan        the writes a resolved rule calls for. It is data the
              caller acts on: resolve applies nothing, and nothing links
              a resolve call to the set-state call that follows it.
  accessor    the declared reader, writer, or gate that touches an
              artifact. Artifact locations are never discovered: every
              one arrives as an explicit --artifact role=path binding.

Model selection

  Every model-taking command requires exactly one of:

    --model <path>   the transition model to load
    --flow <id>      a registered flow id

  --flow is the reserved spelling for config discovery. This build
  registers no ids, so --flow refuses with flow-model-not-found and
  names --model as the remedy. Use --model <path> today.

Output modes

  --as text (default)
    stdout carries the verb's human output; stderr carries advisories
    (note:, warning:) and errors (error: <code>: <message>, with
    optional detail: and hint: lines). Every finding renders on its own
    line, with its identity fields appended.

  --as json
    stdout carries EXACTLY ONE terminal record. On success it is
    {"type":"ok",...}; on failure it is the structured error envelope
    {"code":...,"message":...}. stderr carries advisories only, each
    discriminated by "level". The terminal record is emitted on every
    graceful exit — its absence means the process was killed.

  Aggregate failures (gate results, model-load categories, read-back
  mismatches) report one entry per subject in a top-level "findings"
  array, regardless of how many subjects a given run produces. Scalar
  failures (tag, write, artifact, and selection validation) name their
  one subject in "param" instead. Branch on the code to know which
  carrier to read, never on the number of subjects observed.

Exit codes

  0    the command completed. A gate that denied a candidate under
       flow next is a reported result, not a failure.
  2    the request or the model is wrong, or the model said no. Fix
       the input; re-running it unchanged will refuse again.
  3    the environment could not be consulted — an accessor timed out,
       failed to execute, returned an incomplete key set, or a
       post-mutation read-back did not complete. THE RETRIABLE CLASS:
       repair the environment and re-issue the identical request.
  130  interrupted (SIGINT).

  The 2/3 split is what makes a caller's retry loop safe: a refusal
  about the request never exits 3, so a caller cannot spin on an input
  it must instead fix.

  A read-back that did not complete says the write MAY have been
  applied and was not verified. It is never reported as a write that
  did not occur.

Set values

  A set-valued tag crosses the CLI as a canonical JSON array — members
  sorted, deduplicated, compact, and encoded WITHOUT HTML escaping, so
  <, >, and & serialize as themselves. Every emit and compare site
  builds it the same way, which is what makes copying a plan value back
  into a request, and the read-back check itself, byte equality.

  The reserved value <clear> is unauthorable. Removing an owned tag is
  --clear <key> on flow set-state, never --write key=<clear>.

Further reading

  docs/model-authoring.md      the model class, and authoring a
                               decision table's cells and escape row
  docs/cli-output-contract.md  the authoritative wire contract`

// NewRootCmd builds a fresh command tree. Callers MUST construct a new
// tree per invocation (do not share one across goroutines) — cobra
// commands are not goroutine-safe.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "intrastate",
		Short:   "intrastate CLI",
		Long:    rootLongDesc,
		Version: version.Get().String(),
		// Route every error through the structured gateway instead of
		// cobra's default "Error: …\nRun '… --help'…" text. ExecuteAndEmit
		// converts cobra-level errors (unknown flag/command, missing args)
		// into CLIErrors so the never-silent contract holds at the root.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Persistent root flags shared by every verb.
	cmd.PersistentFlags().String(respond.FlagName, "text", "output mode: text | json")

	// Register verbs here. The version subcommand below is the worked
	// example of the respond/clierr wiring every verb follows.
	cmd.AddCommand(newVersionCmd())
	// Root `lint` is the authoritative graph-acceptance surface
	// (`0006:C19`), deliberately not under RDR 0005's `flow` group.
	cmd.AddCommand(newLintCmd())
	// The `flow` group is RDR 0005's skill-integration surface. `lint`
	// stays at ROOT and is deliberately NOT absorbed into it (`0005:C1`).
	cmd.AddCommand(newFlowCmd())
	// Hidden: it addresses maintainers, not callers. See docs.go.
	cmd.AddCommand(newDocsCmd())

	// MUST run last: registerHelpAllOnTree captures each command's Args
	// at registration time and needs every subcommand already attached.
	withExtendedHelp(cmd, rootExtendedDesc)
	registerHelpAllOnTree(cmd)

	return cmd
}

// Execute is the process entry point: build the tree, run it, and map
// the resulting error to an exit code.
func Execute() {
	if err := ExecuteAndEmit(NewRootCmd(), os.Args[1:]); err != nil {
		os.Exit(clierr.ExitCodeFor(err))
	}
}

// ExecuteAndEmit runs cmd against args and guarantees that any
// cobra/pflag-level error (unknown flag, missing positional, unknown
// subcommand) reaches the caller as a structured CLIError, emitted via
// respond.Fail. Verbs already route their own errors through the
// gateway; this closes the gap for errors cobra raises before RunE.
//
// Tests drive this same path so they exercise the production emission
// flow. Callers pass a freshly built tree (NewRootCmd()) per call.
func ExecuteAndEmit(cmd *cobra.Command, args []string) error {
	cmd.SetArgs(args)
	primeAsFlag(cmd, args)

	// Before cobra dispatch — so the token is never mistaken for a
	// subcommand name, and so this lands ahead of any tag validation.
	if token, ok := argvCarriesUpstreamStop(args); ok {
		return respond.Fail(cmd, upstreamStopError(token))
	}

	err := cmd.Execute()
	if err == nil {
		return nil
	}
	var ce *clierr.CLIError
	if errors.As(err, &ce) {
		return ce
	}
	return respond.Fail(cmd, cobraErrorToCLIError(err))
}

// primeAsFlag commits --as onto the persistent flag before cobra's
// command resolution runs, so respond.ModeOf returns the right mode even
// when cobra errors out (e.g. unknown subcommand) before pflag parses.
func primeAsFlag(cmd *cobra.Command, args []string) {
	flag := cmd.PersistentFlags().Lookup(respond.FlagName)
	if flag == nil {
		return
	}
	for i, a := range args {
		switch {
		case a == "--":
			return
		case len(a) > 5 && a[:5] == "--as=":
			_ = flag.Value.Set(a[5:])
			return
		case a == "--as":
			if i+1 < len(args) {
				_ = flag.Value.Set(args[i+1])
			}
			return
		}
	}
}

// upstreamStopPrefix is the shared stop-packet prefix every cooperating
// tool writes its refusals with. A token carrying it cannot be an
// argument a caller meant to type — it is a refusal line that was
// command-substituted into this argv.
const upstreamStopPrefix = "stopped:"

// argvCarriesUpstreamStop reports the first argv token that is, or
// carries as its value, an upstream tool's stop line. It walks the same
// token forms pflag does — a bare positional, `--flag=value`, and the
// `--flag value` pair — so a stop line is caught whichever slot the
// substitution dropped it into. Tokens after a bare `--` are still
// scanned: `--` ends FLAG parsing, not the substitution hazard.
//
// The prefix match deliberately also refuses a legitimate value that
// happens to begin `stopped:` (e.g. `--tag k=stopped:x`). `stopped:` is
// reserved to the stop-packet convention, so that value is unauthorable
// rather than collateral.
func argvCarriesUpstreamStop(args []string) (string, bool) {
	for _, a := range args {
		// A bare positional, or the value half of a `--flag value`
		// pair — both reach this test as their own token.
		if strings.HasPrefix(a, upstreamStopPrefix) {
			return a, true
		}
		if strings.HasPrefix(a, "-") {
			if _, value, found := strings.Cut(a, "="); found &&
				strings.HasPrefix(value, upstreamStopPrefix) {
				return value, true
			}
		}
	}
	return "", false
}

// upstreamStopError names the token as another command's refusal rather
// than as anything this CLI can parse, so the session debugs the tool
// that actually refused. GroupUserEnv fixes the exit at 2, the same
// shape cobraErrorToCLIError and flow.go's bare-verb refusal use.
func upstreamStopError(token string) *clierr.CLIError {
	return &clierr.CLIError{
		Code: codeArgvUpstreamStop,
		Message: fmt.Sprintf(
			"the argv carries %q, a refusal from the command whose output was substituted; run it alone",
			token),
		Group: clierr.GroupUserEnv,
		Hint:  "re-run the substituted command by itself and act on its refusal",
	}
}

// cobraErrorToCLIError converts a cobra/pflag error into a structured
// CLIError so the never-silent invariant holds at the harness boundary.
func cobraErrorToCLIError(err error) *clierr.CLIError {
	return &clierr.CLIError{
		Code:    "command-error",
		Message: err.Error(),
		Group:   clierr.GroupUserEnv,
		Hint:    "run with `--help` to list supported commands and flags",
	}
}
