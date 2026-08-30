package cli

// RDR 0025 — shared fixtures for the CLI-level command-binding suite.
//
// Two rules hold here, both inherited from the RDR 0005 harness and both
// load-bearing for whether this suite discriminates anything:
//
//  1. State is established and observed THROUGH THE CLI. The RDR fixes the
//     wire and the carrier, not the artifact's on-disk schema — the
//     artifact's byte form belongs to the bound TOOL, not to intrastate —
//     so the MVV's oracle is `flow read-state`, never a file parse.
//  2. Every refusal oracle names its exact code. "An error occurred" would
//     pass against an unregistered flag (cobra's own `command-error`),
//     which is exactly the tautology the red gate must exclude.
//
// The bound tool is `git config --file`, which is the tool the RDR's own
// A3/A7/A8 spikes bound and which the MVV names as "an established tool
// present in CI". Where git is absent the MVV skips loudly rather than
// degrading to a weaker oracle.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// --- the declared wrapper -------------------------------------------------

// writeWrapper writes the A3 spike's `wrapper-write.sh` shape: stdin JSON
// envelope in, `git config` per key out. Its argv is in the model and
// linted; its BODY is script the model does not carry — which is exactly
// what the Approach and Consequences record as the accepted cost.
func writeWrapper(t *testing.T, dir string) string {
	t.Helper()

	p := filepath.Join(dir, "flowstate-write")
	const body = `#!/bin/sh
# stdin: {"<key>":"<value>", ...}   argv: <artifact>
set -e
artifact="$1"
payload=$(cat)
printf '%s' "$payload" |
  sed -e 's/^{//' -e 's/}$//' -e 's/","/"\n"/g' |
  while IFS= read -r pair; do
    [ -n "$pair" ] || continue
    key=$(printf '%s' "$pair" | sed -e 's/^"//' -e 's/":.*$//')
    val=$(printf '%s' "$pair" | sed -e 's/^.*":"//' -e 's/"$//')
    if [ "$val" = "<clear>" ]; then
      git config --file "$artifact" --unset "$key" 2>/dev/null || true
    else
      git config --file "$artifact" "$key" "$val"
    fi
  done
`
	if err := os.WriteFile(p, []byte(body), 0o700); err != nil {
		t.Fatalf("writing the declared wrapper: %v", err)
	}
	return p
}

// traceArgv returns an argv for a SENTINEL child: one that leaves an
// observable trace if it is ever executed. It is the "absence of a spawn"
// oracle S7 requires — a non-zero exit alone cannot tell a refusal before
// spawn from a child that ran and failed.
func traceArgv(t *testing.T, trace string) []string {
	t.Helper()

	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skipf("no POSIX sh on PATH for the spawn sentinel: %v", err)
	}
	// NOTE: this is a fixture's OWN argv, not a model-declared command —
	// a declared `sh -c` vector is a C5 defect (REQ-74) and the models
	// below never carry one.
	return []string{sh, "-c", "printf ran >> " + strconv.Quote(trace)}
}

// --- model writers --------------------------------------------------------

func writeModelFile(t *testing.T, dir, name, src string) string {
	t.Helper()

	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}
	return p
}

func tomlArgv(argv []string) string {
	quoted := make([]string, 0, len(argv))
	for _, a := range argv {
		quoted = append(quoted, strconv.Quote(a))
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// writeCommandModel writes a model whose WRITE entry carries the supplied
// argv and whose read entry binds `git config --file` directly in raw mode.
func writeCommandModel(t *testing.T, dir string, writeArgv []string) string {
	t.Helper()

	src := commandModelHeader + `
[read.state]
role = "state"
command = ["git", "config", "--file", "{artifact}", "--get", "flow.status"]
output = "raw"
exit_absent = [1]
keys = ["status"]
timeout = "10s"

[write.state]
role = "state"
command = ` + tomlArgv(append(writeArgv, "{artifact}")) + `
keys = ["status"]
timeout = "10s"
read_back = true
` + commandModelTail
	return writeModelFile(t, dir, "command-model.toml", src)
}

// writeCommandReadModel writes a model whose READ entry carries the
// supplied argv, so a read-side failure can be driven without a write.
func writeCommandReadModel(t *testing.T, dir string, readArgv []string) string {
	t.Helper()

	src := commandModelHeader + `
[read.state]
role = "state"
command = ` + tomlArgv(append(readArgv, "{artifact}")) + `
output = "raw"
keys = ["status"]
timeout = "10s"

[write.state]
role = "state"
command = ["true", "{artifact}"]
keys = ["status"]
timeout = "10s"
read_back = true
` + commandModelTail
	return writeModelFile(t, dir, "command-read-model.toml", src)
}

// writeAppliedThenFailingReadBackModel writes a model whose write APPLIES
// (so the applied-but-unverified sense is in force) and whose reader then
// fails with a stderr tail — the shape where BOTH senses compete for the
// CLI's one Detail slot.
func writeAppliedThenFailingReadBackModel(t *testing.T, dir string) string {
	t.Helper()

	failing := filepath.Join(dir, "failing-reader")
	const body = `#!/bin/sh
printf 'wrapper: cannot read the artifact\n' >&2
exit 3
`
	if err := os.WriteFile(failing, []byte(body), 0o700); err != nil {
		t.Fatalf("writing the failing reader: %v", err)
	}

	src := commandModelHeader + `
[read.state]
role = "state"
command = ["` + failing + `", "{artifact}"]
output = "raw"
keys = ["status"]
timeout = "10s"

[write.state]
role = "state"
command = ["true", "{artifact}"]
keys = ["status"]
timeout = "10s"
read_back = true
` + commandModelTail
	return writeModelFile(t, dir, "applied-unverified.toml", src)
}

const commandModelHeader = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "cmdcli"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true
`

const commandModelTail = `
[initial]
status = "draft"

[context.done]
[context.done.match.status]
eq = "final"

[[rule]]
id = "advance-draft"
[rule.match.status]
eq = "draft"
[rule.match.recognized]
eq = "advance"
[rule.write]
status = "final"
`

// requireGit skips when the bound tool is absent. The MVV binds "an
// established tool present in CI", so a missing git is a SKIP with a
// reason, never a silently weakened oracle.
func requireGit(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("the MVV binds `git config --file`, an established tool the "+
			"RDR's own A3 spike bound; git is not on PATH: %v", err)
	}
}
