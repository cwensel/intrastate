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
	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/respond"
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

Every verb takes exactly one of --model <path> or --flow <id>, and binds
each artifact explicitly as --artifact role=path: nothing about a flow's
location is discovered. This build resolves no --flow ids yet, so pass
--model <path>.

Run any verb with --help-all for its refusal codes and worked calls.`,
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

	// `0025:C6` — ONE registration, on the GROUP's persistent set, whose
	// subtree is exactly the `buildRequest` callers and nothing else. That
	// makes the caller set STRUCTURAL rather than a hand-kept verb list: a
	// fifth verb added inside the group inherits the gate, and one added
	// outside it loses command execution loudly. `lint` sits at ROOT,
	// outside the group, so it does not carry the flag — validation is
	// ungated, and a flag that changes nothing is a false affordance.
	cmd.PersistentFlags().Bool("allow-commands", false,
		"permit model-declared command accessors to execute (off by default)")

	cmd.AddCommand(
		newFlowNextCmd(),
		newFlowResolveCmd(),
		newFlowReadStateCmd(),
		newFlowSetStateCmd(),
	)
	withExtendedHelp(cmd, flowExtendedDesc)
	return cmd
}

// flowExtendedDesc is the group's --help-all body: the shared grammar
// and the division of labour among the four verbs. Per-verb refusals
// live on each verb's own extended body.
const flowExtendedDesc = `The four verbs divide one job, and the division is deliberate:

  next        reports what the state does not exclude. It never selects.
  resolve     selects exactly one row, or refuses. It never writes.
  set-state   writes, and verifies by read-back. It never selects.
  read-state  reports what the readers see. It decides nothing.

Nothing links one call to the next. A resolve plan is data you may act
on; set-state re-derives its own legality from its own request, and the
read-back is the only commit-time check. A caller that drops a plan on
the floor has broken nothing.

Shared input grammar

  --model <path>        the transition model. Exactly one of --model or
                        --flow is required.
  --flow <id>           reserved for config discovery. This build
                        registers no ids and refuses with
                        ` + codeModelNotFound + `; use --model.
  --artifact role=path  bind one declared accessor role to a file.
                        Repeatable. Never discovered — an unbound role a
                        verb needs refuses with ` + codeArtifactMissing + `.
  --tag name=value      observed context. Repeatable. A set value is a
                        JSON array literal. A repeated key is
                        ` + codeTagDuplicate + `; a malformed one is
                        ` + codeTagInvalid + `.

  A model that cannot be loaded refuses under ` + codeModelInvalid + `,
  carrying one finding per load category.

  --tag carries OBSERVED tags only. Naming an owned key refuses with
  ` + codeTagOwned + ` and naming a reserved key refuses with
  ` + codeTagReserved + `: owned state comes from the readers and is
  written only through set-state, so accepting it as context would let a
  caller assert state the model owns. --tag is absent from read-state,
  which reports what readers see rather than evaluating anything.

Ordering

  Every input refusal precedes every accessor invocation. A caller whose
  --tag names an owned key learns that, and not that some artifact role
  went unbound: the input is what they must fix, and consulting the
  environment first would report the wrong subject.

  Gates run AFTER exact-one selection and BEFORE the plan is emitted.
  Earlier would gate rows the model never selected; later would emit a
  plan a gate denies.

Reader narrowing

  next and resolve invoke only the readers the model's relevant rows
  demand, so an artifact no candidate needs may stay unbound. read-state
  is diagnostic and runs EVERY declared reader, so it requires every
  declared role to be bound. The demand set is a property of the model,
  not of the flags.

Environment failures

  ` + codeAccessorTimeout + `, ` + codeAccessorFailed + `, and
  ` + codeReadIncomplete + ` exit 3: the environment could not be
  consulted and the identical request may be re-run once it is repaired.
  Every refusal about the request itself exits 2.`

// registerSelectionFlags adds the model-selection flags every verb carries.
// The pair is mutually exclusive and one is required (REQ-22).
func registerSelectionFlags(cmd *cobra.Command) {
	cmd.Flags().String("model", "", "path to the transition model")
	cmd.Flags().String("flow", "", "flow id (reserved; this build resolves none — use --model)")
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
