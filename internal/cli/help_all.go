package cli

// Two-tier `--help`, ported from retrofit's internal/cli/help_all.go.
//
// The default `--help` stays terse: the command's `Long` (an
// orientation blurb), `Usage:`, `Flags:`, `Global Flags:`. Extended
// reference material — vocabulary, wire shapes, exit-code semantics —
// moves into a separate "extended" body that surfaces only on
// `--help-all` (per-command) or `intrastate help --all [<cmd>]`.
//
// Mechanism:
//   - Each command stores its extended body in cmd.Annotations under
//     extendedHelpKey. registerHelpAllOn adds a per-command --help-all
//     flag and hooks ParseFlags → Args → PreRunE: when the flag is
//     set, PreRunE prints terse + extended itself and returns
//     pflag.ErrHelp so cobra's existing ErrHelp path silences the
//     downstream error message and exits 0.
//   - The root carries a wrapping HelpFunc that suppresses the second
//     print cobra issues after PreRunE returns ErrHelp (the emit
//     already happened) AND handles non-Runnable parents like
//     `intrastate --help-all` (where cobra short-circuits to ErrHelp
//     BEFORE PreRunE runs).
//   - The auto-generated `help` subcommand on the root carries an
//     `--all` flag; when set, the help command prints the targeted
//     command in extended mode (or, with no arg, the entire tree).
//
// Why we don't wrap cmd.HelpFunc on every command: cobra's HelpFunc()
// walks up the tree when a child has no own helpFunc, so capturing it
// during depth-first registration pulls in already-wrapped parent
// funcs and re-emits extended bodies once per tree level.
//
// Why not mutate cmd.Long at help-print time: cobra's help template
// runs off cmd.Long; swapping it under help invocations would race in
// tests that share a cmd tree. The Annotations slot keeps the source
// of truth stable.
//
// Help output is deliberately NOT routed through the respond gateway.
// Help is not a terminal result: it is cobra's own stream, exits 0,
// and carries no envelope. Under `--as=json` a caller asking for help
// gets help text, and the never-silent contract still binds every
// path that produces a result or a refusal.

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// extendedHelpKey stores per-command extended help text on
	// cmd.Annotations. Empty / missing means the command has no
	// extended help; --help-all falls back to the terse default.
	extendedHelpKey = "intrastate:extended-help"
	helpAllFlagName = "help-all"

	// helpAllEmittedKey marks a command whose --help-all body was
	// already printed by emitHelpAll, so the wrapper installed on the
	// root's HelpFunc can suppress cobra's downstream HelpFunc call
	// (which would re-print the terse body).
	helpAllEmittedKey = "intrastate:help-all-emitted"
)

// withExtendedHelp records `extended` as the command's extended-help
// body and returns cmd for fluent chaining at constructor sites.
// Whitespace is trimmed because the constructor sites use raw backtick
// strings that often start/end with newlines.
func withExtendedHelp(cmd *cobra.Command, extended string) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[extendedHelpKey] = strings.TrimSpace(extended)
	return cmd
}

// extendedHelpFor returns the extended-help body for cmd, or "" if
// none was registered.
func extendedHelpFor(cmd *cobra.Command) string {
	if cmd.Annotations == nil {
		return ""
	}
	return cmd.Annotations[extendedHelpKey]
}

// registerHelpAllOnTree walks the command tree rooted at root and adds
// the --help-all flag plus the extended-help PreRunE hook to every
// command, then wires --all onto the auto-generated `help` subcommand.
//
// Call this AFTER all subcommands are registered and after every
// command's Args is set: registerHelpAllOn captures cmd.Args at
// registration time. NewRootCmd calls it just before returning.
func registerHelpAllOnTree(root *cobra.Command) {
	walkCommandTree(root, registerHelpAllOn)

	// Install a root-level HelpFunc that:
	//
	//  1. Suppresses the second help print cobra issues after our
	//     PreRunE returned ErrHelp (we already emitted help inside
	//     emitHelpAll). The marker on Annotations[helpAllEmittedKey]
	//     drives the suppression.
	//  2. Emits extended help for non-Runnable parent commands —
	//     `intrastate --help-all` — where cobra short-circuits to
	//     ErrHelp BEFORE PreRunE runs because `!c.Runnable()`.
	//
	// Cobra's HelpFunc() walks up the tree when a child has no own
	// helpFunc, so installing this guard on the root catches every
	// descendant.
	//
	// We deliberately call our own defaultHelp on the fall-through
	// rather than capturing root.HelpFunc() before SetHelpFunc — the
	// captured func would be cobra's defaultHelpFunc, whose semantics
	// are an internal cobra contract.
	root.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c.Annotations != nil && c.Annotations[helpAllEmittedKey] == "1" {
			// Already printed by emitHelpAll; clear the marker so a
			// later --help on the same tree (tests reuse instances)
			// doesn't silently swallow help output.
			delete(c.Annotations, helpAllEmittedKey)
			return
		}
		if helpAllRequested(c) {
			// Non-Runnable parent: PreRunE never fired (cobra returns
			// ErrHelp before the hook). Drive the emit here.
			emitHelpAll(c)
			delete(c.Annotations, helpAllEmittedKey)
			return
		}
		defaultHelp(c, args)
	})

	// Cobra creates the help command lazily on first Execute or via
	// InitDefaultHelpCmd. Force-init so we can wire --all on it now.
	root.InitDefaultHelpCmd()
	for _, sub := range root.Commands() {
		if sub.Name() == "help" {
			wireHelpSubcommandAll(root, sub)
			break
		}
	}
}

// registerHelpAllOn adds the --help-all flag to a single command and
// wires it to print the terse help plus the extended body. The flag is
// registered on EVERY command, including parent groups with no
// extended body, for UX consistency — a user typing
// `intrastate flow --help-all` shouldn't get an unknown-flag refusal.
func registerHelpAllOn(cmd *cobra.Command) {
	if cmd.Flags().Lookup(helpAllFlagName) != nil {
		return
	}
	cmd.Flags().Bool(helpAllFlagName, false,
		"show extended help (vocabulary, wire shapes, exit codes)")

	// Cobra's order is: ParseFlags → check --help → !Runnable → Args
	// validation → PreRunE → RunE. We need --help-all to short-circuit
	// before Args runs, otherwise a command with a positional-arg
	// validator would fail that check before we ever see the flag.
	//
	// Only wrap Args when the command already declared one. Cobra's
	// Find checks `commandFound.Args == nil` and only then runs
	// legacyArgs — the source of the `unknown command "X" for "Y"`
	// error that root.go's ExecuteAndEmit converts into a structured
	// CLIError. Wrapping Args on a parent that previously had nil
	// flips that detection off, and `intrastate blah` would silently
	// print help at exit 0 instead of refusing. Non-Runnable parents
	// never reach the Args check at execute time anyway, so leaving
	// Args nil on them costs nothing.
	if prevArgs := cmd.Args; prevArgs != nil {
		cmd.Args = func(c *cobra.Command, args []string) error {
			if helpAllRequested(c) {
				return nil
			}
			return prevArgs(c, args)
		}
	}

	// PreRunE prints the terse + extended body itself, then returns
	// pflag.ErrHelp so cobra's existing ErrHelp path silences any
	// downstream error message and exits cleanly. emitHelpAll marks
	// the command so we don't double-print when cobra calls HelpFunc
	// on its way out.
	//
	// For non-Runnable parents (the root, `flow`) cobra short-circuits
	// to ErrHelp BEFORE PreRunE runs, so this hook never fires there —
	// the live path for them is the HelpFunc wrapper above.
	prevPreRunE := cmd.PreRunE
	prevPreRun := cmd.PreRun
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if helpAllRequested(c) {
			emitHelpAll(c)
			return pflag.ErrHelp
		}
		if prevPreRunE != nil {
			return prevPreRunE(c, args)
		}
		if prevPreRun != nil {
			prevPreRun(c, args)
		}
		return nil
	}
}

// emitHelpAll writes the terse help and the command's extended body.
// For the root, whose own extended body is the concepts orientation,
// it ALSO walks the tree and dumps every descendant's extended body:
// `intrastate --help-all` is the single self-describing dump an agent
// can read in one call.
func emitHelpAll(c *cobra.Command) {
	defaultHelp(c, nil)
	out := c.OutOrStdout()
	if ext := extendedHelpFor(c); ext != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Extended help:")
		fmt.Fprintln(out)
		fmt.Fprintln(out, ext)
	}
	if c.Parent() == nil {
		writeAllExtended(out, c, true)
	}
	if c.Annotations == nil {
		c.Annotations = map[string]string{}
	}
	c.Annotations[helpAllEmittedKey] = "1"
}

// wireHelpSubcommandAll adds --all to the auto-generated `help`
// subcommand. `intrastate help --all <cmd>` prints <cmd> in extended
// mode; `intrastate help --all` (no arg) walks the tree.
func wireHelpSubcommandAll(root, helpCmd *cobra.Command) {
	if helpCmd.Flags().Lookup("all") != nil {
		return
	}
	helpCmd.Flags().Bool("all", false,
		"print extended help (vocabulary, wire shapes, exit codes)")
	originalRun := helpCmd.Run
	helpCmd.Run = func(c *cobra.Command, args []string) {
		all, _ := c.Flags().GetBool("all")
		if !all {
			if originalRun != nil {
				originalRun(c, args)
			}
			return
		}
		// With no args, dump the whole tree — through the SAME path
		// `intrastate --help-all` takes, so the two documented spellings
		// of "the whole reference" cannot diverge. Calling HelpFunc plus
		// writeAllExtended separately is what let them: that pair printed
		// the terse root body and every DESCENDANT's extended body, but
		// silently dropped the root's own — the vocabulary, output modes,
		// and exit contract — from the alias.
		if len(args) == 0 {
			emitHelpAll(root)
			if root.Annotations != nil {
				delete(root.Annotations, helpAllEmittedKey)
			}
			return
		}
		target, _, err := root.Find(args)
		if target == nil || err != nil {
			c.Printf("Unknown help topic %#q\n", args)
			_ = root.Usage()
			return
		}
		// Drive the same code path as `<target> --help-all`: setting
		// the flag is what root's HelpFunc checks. Reset it afterwards
		// so a later help invocation on the same tree doesn't carry
		// the bit forward (tests reuse the tree).
		_ = target.Flags().Set(helpAllFlagName, "true")
		defer func() { _ = target.Flags().Set(helpAllFlagName, "false") }()
		target.HelpFunc()(target, args)
	}
}

// helpAllRequested reports whether --help-all (or the help
// subcommand's --all forwarded onto the target) was set on cmd.
//
// The Lookup-then-GetBool dance is deliberate: GetBool only errors
// when the named flag is absent or registered as a non-Bool, both of
// which are statically prevented here. A future change that
// re-registered the flag with a different type would surface as this
// returning false rather than panicking.
func helpAllRequested(cmd *cobra.Command) bool {
	if f := cmd.Flags().Lookup(helpAllFlagName); f != nil {
		v, _ := cmd.Flags().GetBool(helpAllFlagName)
		return v
	}
	return false
}

// defaultHelp mirrors cobra's defaultHelpFunc: print Long (or Short)
// then the UsageString. We re-implement it rather than calling cobra's
// exported HelpFunc() because that walks up the tree to the parent
// when c.helpFunc is nil — during our registration the parent has
// already been wrapped, so each ancestor's wrapper would run too and
// emit the extended body once per tree level.
func defaultHelp(c *cobra.Command, _ []string) {
	c.InitDefaultHelpFlag()
	c.InitDefaultVersionFlag()
	out := c.OutOrStdout()
	body := c.Long
	if body == "" {
		body = c.Short
	}
	body = strings.TrimRight(body, " \t\r\n")
	if body != "" {
		fmt.Fprintln(out, body)
		fmt.Fprintln(out)
	}
	if c.Runnable() || c.HasSubCommands() {
		fmt.Fprint(out, c.UsageString())
	}
}

// writeAllExtended walks the tree under root and writes every
// command's terse Long plus extended body to w, each under a
// `## <command path>` header. skipRoot omits the root itself, whose
// body the caller has already printed.
func writeAllExtended(w io.Writer, root *cobra.Command, skipRoot bool) {
	walkCommandTree(root, func(c *cobra.Command) {
		if skipRoot && c == root {
			return
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "## %s\n", c.CommandPath())
		fmt.Fprintln(w)
		if c.Short != "" {
			fmt.Fprintln(w, c.Short)
			fmt.Fprintln(w)
		}
		fmt.Fprint(w, c.UsageString())
		if ext := extendedHelpFor(c); ext != "" {
			fmt.Fprintln(w)
			fmt.Fprintln(w, ext)
		}
	})
}

// walkCommandTree applies fn to root and every descendant,
// depth-first. The auto-generated `help` and `completion` subcommands
// are skipped: they carry no user-authored extended body.
func walkCommandTree(root *cobra.Command, fn func(*cobra.Command)) {
	fn(root)
	for _, sub := range root.Commands() {
		if sub.Name() == "help" || sub.Name() == "completion" {
			continue
		}
		walkCommandTree(sub, fn)
	}
}
