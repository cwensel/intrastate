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
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print build version, commit, and date",
		Long: `Print this build's version, commit, and date.

Under --as=json the same identity is emitted as the structured Info
value under the terminal "ok" envelope, so a caller can pin the build
that produced any other output.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ce := respond.ValidateMode(cmd); ce != nil {
				return respond.Fail(cmd, ce)
			}
			// Both modes route through the gateway. An earlier cut
			// special-cased text with `cmd.Println`, which cobra sends to
			// OutOrStderr — so `$(intrastate version)` captured nothing,
			// and the one direct print in this package sat in the file
			// CONTRIBUTING names as the copy-me example. respond/text.go
			// calls that print out by name as the drift it exists to
			// avoid (REQ-10); this is that drift removed.
			return respond.OK(cmd, respond.Success{Data: version.Get()})
		},
	}
	withExtendedHelp(cmd, versionExtendedDesc)
	return cmd
}

const versionExtendedDesc = `The three fields identify the build, not the model or the flow:

  version  the release identity, or ` + "`dev`" + ` for an unstamped local build.
  commit   the source revision the binary was built from.
  date     the build timestamp.

They are stamped at link time. A build produced without those stamps
reports the ` + "`dev`" + ` identity rather than guessing or failing, so a
locally built binary is still self-identifying.

Pin the build alongside any captured output — the finding codes, the
refusal vocabulary, and the lint bounds are all properties of a
particular build, and ` + "`--help-all`" + ` on this binary is the
authoritative statement of what THIS build does.

  intrastate version --as=json`
