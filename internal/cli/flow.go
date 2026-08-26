package cli

// RDR 0005 — the `flow` command group: the skill-integration CLI contract.
//
// One group, four verbs (`0005:C1`, `0005:D-naming`):
//
//	flow next        enumerate the legal outcomes and their candidate rows
//	flow resolve     select exactly one plan, or refuse
//	flow read-state  report declared readers' tags, diagnostically
//	flow set-state   apply planned owned-tag writes, read-back-verified
//
// The group owns the CLI contract and nothing beneath it. Transition-model
// representation is RDR 0002's, guard semantics RDR 0003's, the resolution
// kernel RDR 0001's, accessor execution safety RDR 0004's, and graph lint
// RDR 0006's — this file parses input, selects a model, invokes those
// seams in a fixed order, and renders ONE terminal result (REQ-113).
//
// Two orderings in here are load-bearing rather than incidental:
//
//  1. Every input refusal precedes every accessor invocation (REQ-27,
//     REQ-61). A caller whose `--tag` names an owned key learns that, not
//     that some artifact role went unbound — the input is what they must
//     fix, and running accessors first would report the wrong subject.
//  2. Gates run AFTER exact-one selection and BEFORE the plan is emitted
//     (JDR 0001 §D9, REQ-50). Running them earlier would gate rows the
//     model never selected; running them later would emit a plan a gate
//     denies.

import (
	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/spf13/cobra"
)

// newFlowCmd builds the `flow` command group.
func newFlowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flow",
		Short: "Drive a transition model from a skill",
		Long: `Drive a transition model from a skill.

The four verbs are the whole skill-integration surface:

  next        list the legal recognized outcomes and their candidate rows
  resolve     map one recognized outcome to exactly one plan, or refuse
  read-state  report what the declared read accessors see
  set-state   apply planned owned-tag writes and verify them by read-back

Model selection takes exactly one of --model <path> or --flow <id>.
Artifact locations are never discovered: every one arrives as an explicit
--artifact role=path binding.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		// A bare `flow` with no verb is a usage error, not a silent
		// success: the group carries no behaviour of its own. It routes
		// through the SAME gateway and the SAME group-level usage code
		// (`command-error`, GroupUserEnv / exit 2) that root.go already
		// gives the sibling case `flow <unknown-verb>` — so both usage
		// errors under this group are one disposition, and neither
		// leaves stdout carrying help text where `--as=json` reserves it
		// for the single terminal record (REQ-113).
		//
		// ValidateMode runs FIRST, matching REQ-5's ordering on the four
		// verbs: a caller who cannot name the output mode cannot read
		// the usage refusal either.
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ce := respond.ValidateMode(cmd); ce != nil {
				return respond.Fail(cmd, ce)
			}
			return respond.Fail(cmd, &clierr.CLIError{
				Code: "command-error",
				Message: "`flow` requires a verb: next, resolve, " +
					"read-state, set-state",
				Group: clierr.GroupUserEnv,
				Hint:  "run `intrastate flow --help` to list the verbs",
			})
		},
	}

	cmd.AddCommand(
		newFlowNextCmd(),
		newFlowResolveCmd(),
		newFlowReadStateCmd(),
		newFlowSetStateCmd(),
	)
	return cmd
}

// registerSelectionFlags adds the model-selection flags every verb carries.
// The pair is mutually exclusive and one is required (REQ-22).
func registerSelectionFlags(cmd *cobra.Command) {
	cmd.Flags().String("model", "", "path to the transition model")
	cmd.Flags().String("flow", "", "flow id, resolved through config discovery")
	cmd.Flags().StringArray("artifact", nil,
		"artifact role binding, as role=path (repeatable)")
}

// registerTagFlag adds `--tag`, the observed-context channel. It is absent
// from `read-state`, which reports what readers see rather than evaluating
// anything against caller context.
func registerTagFlag(cmd *cobra.Command) {
	cmd.Flags().StringArray("tag", nil,
		"observed tag, as name=value (repeatable); set values are JSON arrays")
}
