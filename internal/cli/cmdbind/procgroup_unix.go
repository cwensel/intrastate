//go:build !windows && !js && !plan9

package cmdbind

import (
	"os/exec"
	"syscall"
)

// procGroupPlatforms is the set this file is NOT built for. It mirrors
// `Unsupported()`'s refuse-list, and a drift-guard test asserts the two
// agree — the build tag above and the runtime predicate are the same
// platform set expressed twice, and they must not diverge (`0025:C4`).
var procGroupPlatforms = []string{"windows", "js", "plan9"}

// procGroupSupported reports whether THIS FILE — the syscall-bearing half —
// is the one compiled in. Only the mechanism is split by build tag; the
// refusal, the `goos` var and all lint stay untagged in one
// platform-neutral binary, which is what `0025:C4` requires when it rejects
// a build constraint on the refusal itself.
const procGroupSupported = true

// setProcGroup puts the child in its OWN process group, so the group id is
// the child's pid and the negative-pid signal in `killGroup` addresses the
// group rather than the caller's (`0025:C4`, A1).
func setProcGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killGroup SIGKILLs the process group led by pid. The caller owns the
// guard on pid, since a pid of 0 or 1 would address the caller's group or
// init.
func killGroup(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}
