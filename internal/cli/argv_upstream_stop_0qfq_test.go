package cli

// Kata 0qfq — an argv token beginning `stopped:` is another tool's stop
// line, command-substituted into this invocation, and must be refused as
// such rather than looked up as a subcommand name.
//
// `flow resolve` and `lint` are both `cobra.NoArgs` leaves, so cobra's
// `Find` resolves a positional token as a subcommand BEFORE `Args`
// validation fires, and the refusal came back as
// `command-error: unknown command "stopped:no-such-record"` — the
// upstream refusal erased, the session debugging the wrong tool. The
// check therefore lives at `root.go::ExecuteAndEmit`, the one gateway
// both commands pass through, ahead of cobra dispatch entirely.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// TestKata0qfq_ArgvCarryingUpstreamStopIsRefused is the RED test: every
// token form a substitution can land in must name the upstream tool.
func TestKata0qfq_ArgvCarryingUpstreamStopIsRefused(t *testing.T) {
	const stopLine = "stopped:no-such-record (0999 in /records)"

	cases := map[string][]string{
		// The consumer scenario verbatim: `$(rdr status --tags 0999)`
		// word-split onto the end of a `flow resolve` call.
		"flow resolve positional": {
			"flow", "resolve", "--model", "x.toml", "--outcome", "y", stopLine,
		},
		// `lint` is at root and equally exposed.
		"lint positional": {"lint", "--model", "x.toml", stopLine},
		// A substitution that landed in a flag's value instead: both
		// the `--flag value` pair and the `--flag=value` spelling.
		"flag value pair":   {"flow", "resolve", "--model", stopLine, "--outcome", "y"},
		"flag value joined": {"flow", "resolve", "--model=" + stopLine, "--outcome", "y"},
	}

	for name, args := range cases {
		t.Run(name, func(t *testing.T) {
			_, stderr, err := runCmd(t, args...)
			if err == nil {
				t.Fatalf("args %q: want a refusal, got success", args)
			}
			var ce *clierr.CLIError
			if !asCLIError(err, &ce) {
				t.Fatalf("args %q: error is not a CLIError: %v", args, err)
			}
			if ce.Code != codeArgvUpstreamStop {
				t.Errorf("code = %q; want %q", ce.Code, codeArgvUpstreamStop)
			}
			if got := clierr.ExitCodeFor(err); got != 2 {
				t.Errorf("exit = %d; want 2", got)
			}
			// The message must carry the stop line back verbatim, so the
			// caller sees WHICH refusal arrived, and must name it as
			// coming from a substituted command.
			if !strings.Contains(ce.Message, stopLine) {
				t.Errorf("message %q does not quote the stop line %q", ce.Message, stopLine)
			}
			if !strings.Contains(ce.Message, "substituted") {
				t.Errorf("message %q does not name the token as upstream", ce.Message)
			}
			// The never-silent contract: the refusal reached stderr.
			if !strings.Contains(stderr, codeArgvUpstreamStop) {
				t.Errorf("stderr = %q; want it to carry %q", stderr, codeArgvUpstreamStop)
			}
		})
	}
}

// TestKata0qfq_UnknownCommandStillReportsUnknownCommand is the CONTROL:
// a positional that is merely wrong, not a stop line, keeps its own
// diagnosis. The new gate must not swallow cobra's dispatch errors.
func TestKata0qfq_UnknownCommandStillReportsUnknownCommand(t *testing.T) {
	_, _, err := runCmd(t, "flow", "resolve", "--model", "x.toml", "--outcome", "y", "bogus")
	if err == nil {
		t.Fatal("want a refusal for an unknown positional, got success")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("error is not a CLIError: %v", err)
	}
	if ce.Code == codeArgvUpstreamStop {
		t.Errorf("code = %q; a plain unknown positional must not be read as a stop line", ce.Code)
	}
}

// TestKata0qfq_NormalArgvStillDispatches is the CONTROL for the happy
// path: with no `stopped:` token the gate is transparent and the verb
// reaches its own validation. `nope.toml` does not exist, so the refusal
// that comes back must be the model-load one, proving dispatch ran.
func TestKata0qfq_NormalArgvStillDispatches(t *testing.T) {
	_, _, err := runCmd(t, "flow", "resolve", "--model", "nope.toml", "--outcome", "y")
	if err == nil {
		t.Fatal("want a refusal for a missing model, got success")
	}
	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("error is not a CLIError: %v", err)
	}
	if ce.Code != codeModelNotFound {
		t.Errorf("code = %q; want %q — the verb's own validation must have run",
			ce.Code, codeModelNotFound)
	}
}
