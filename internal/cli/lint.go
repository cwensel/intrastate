package cli

// SKELETON (RDR 0006 Stage 8, Phase 1). The root `lint` verb — the
// authoritative CLI surface for graph acceptance (`0006:C19`). It is
// registered at ROOT, deliberately not under RDR 0005's `flow` group.
//
// The RunE body is unimplemented: it routes a placeholder through the
// respond gateway so the command exists, the flags parse, and the RDR
// 0006 tests fail on their own assertions rather than on an unknown
// subcommand. Phase 2 wires it to graphlint.NewRequest / graphlint.Run.

import (
	"os"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/newcoinc/intrastate/internal/graphlint"
	"github.com/newcoinc/intrastate/internal/table"
	"github.com/spf13/cobra"
)

// lintPayload is the success-side `data` value. It is a non-nil STRUCT
// value, never a nil map: `respond.Success.Data` is itself
// `json:"data,omitempty"`, so a nil payload would drop the `data` key and
// take the empty-list receipt with it (`0006:C14`).
type lintPayload struct {
	Findings []clierr.Finding `json:"findings"`
}

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Check a transition model against the graph invariants",
		Long: `Check a transition model against the mandatory graph invariants.

A model with any blocking finding is refused: the command returns the
aggregate error ` + graphlint.AggregateCode + ` at exit 2 and carries every
blocking finding in the machine-readable findings list.

Bounds enforced by this build (model-independent implementation
constants, not per-model inputs):

  product bound   ` + itoa(graphlint.ProductBound()) + `
  node ceiling    ` + itoa(graphlint.NodeCeiling()) + `
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE:          runLint,
	}
	cmd.Flags().String("model", "", "path to the transition model to lint")
	cmd.Flags().String("flow", "", "flow id to lint, resolved through config discovery")
	return cmd
}

func runLint(cmd *cobra.Command, _ []string) error {
	if ce := respond.ValidateMode(cmd); ce != nil {
		return respond.Fail(cmd, ce)
	}

	modelPath, _ := cmd.Flags().GetString("model")
	flowID, _ := cmd.Flags().GetString("flow")
	if modelPath != "" && flowID != "" {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-mutually-exclusive",
			Message: "--flow and --model are mutually exclusive",
			Group:   clierr.GroupUserEnv,
		})
	}
	if modelPath == "" {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "flag-required",
			Param:   "model",
			Message: "--model <path> names the transition model to lint",
			Group:   clierr.GroupUserEnv,
		})
	}

	src, err := os.ReadFile(modelPath)
	if err != nil {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "model-unreadable",
			Param:   "model",
			Message: "cannot read " + modelPath,
			Detail:  err.Error(),
			Group:   clierr.GroupUserEnv,
			Cause:   err,
		})
	}
	m, err := table.Load(src, modelPath)
	if err != nil {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:    "model-invalid",
			Param:   "model",
			Message: "model does not conform to the transition-model schema",
			Detail:  err.Error(),
			Group:   clierr.GroupUserEnv,
			Cause:   err,
		})
	}

	report := graphlint.Run(graphlint.NewRequest(m))
	blocking := report.Blocking()
	if len(blocking) > 0 {
		return respond.Fail(cmd, &clierr.CLIError{
			Code:     graphlint.AggregateCode,
			Message:  "the model carries blocking graph-lint findings",
			Group:    clierr.GroupUserEnv,
			Findings: blocking,
		})
	}

	advisory := report.Advisory()
	if advisory == nil {
		advisory = []clierr.Finding{}
	}
	return respond.OK(cmd, respond.Success{Data: lintPayload{Findings: advisory}})
}

// itoa is a local decimal renderer so the long help string stays a
// compile-time concatenation.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
