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
	"os"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/newcoinc/intrastate/internal/version"
	"github.com/spf13/cobra"
)

const rootLongDesc = `intrastate — <one-line description of what this tool does>.

Run any subcommand with --help for its flags. Global flags:

  --as text|json    output mode (default text)`

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
