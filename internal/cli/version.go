package cli

import (
	"github.com/newcoinc/intrastate/internal/cli/respond"
	"github.com/newcoinc/intrastate/internal/version"
	"github.com/spf13/cobra"
)

// newVersionCmd is the worked example of a verb: validate the output
// mode, then route the result through the respond gateway. In text mode
// it prints the build-identity string; in json mode it emits the
// structured Info under the terminal "ok" envelope.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print build version, commit, and date",
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ce := respond.ValidateMode(cmd); ce != nil {
				return respond.Fail(cmd, ce)
			}
			info := version.Get()
			if respond.ModeOf(cmd) == respond.ModeText {
				cmd.Println(info.String())
				return nil
			}
			return respond.OK(cmd, respond.Success{Data: info})
		},
	}
}
