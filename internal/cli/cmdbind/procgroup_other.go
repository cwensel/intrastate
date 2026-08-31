//go:build windows || js || plan9

package cmdbind

import (
	"errors"
	"os/exec"
)

// procGroupPlatforms is the set this file IS built for. It mirrors
// `Unsupported()`'s refuse-list, and a drift-guard test asserts the two
// agree — the build tag above and the runtime predicate are the same
// platform set expressed twice, and they must not diverge (`0025:C4`).
var procGroupPlatforms = []string{"windows", "js", "plan9"}

// procGroupSupported reports whether the syscall-bearing half is the one
// compiled in. It is false here, and `Unsupported()` is true on exactly
// these platforms, so neither function below is reachable at runtime: the
// pre-spawn platform refusal fires first and no child is ever created
// (`0025:C4`).
const procGroupSupported = false

// errNoProcGroup is why `killGroup` cannot succeed here. Returning it —
// rather than nil — is what makes `cmd.Cancel` take its existing fallback
// to `cmd.Process.Kill()`, so the cancel path degrades to the direct child
// instead of silently reporting a termination that never happened.
var errNoProcGroup = errors.New("process groups are unsupported on this platform")

// setProcGroup is a no-op: there is no process-group mechanism to ask for,
// and this platform is refuse-listed before any spawn.
func setProcGroup(_ *exec.Cmd) {}

// killGroup always fails here, which is the honest answer: no group was
// created, so none can be signalled.
func killGroup(_ int) error {
	return errNoProcGroup
}
