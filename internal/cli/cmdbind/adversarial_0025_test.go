package cmdbind_test

// Phase 3b adversarial tests for RDR 0025.
//
// Each test here names a failure mode the RDR's own Failure Modes section
// admits, and asserts the property the section claims holds. They are
// written from the contract, not from the implementation.

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

// advScript writes an executable shell script into a fresh temp dir and
// returns its absolute path.
func advScript(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "tool.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatalf("writing the probe tool: %v", err)
	}
	return path
}

// advArtifact binds the `state` role to an absolute path, which C2
// requires of every command entry.
func advArtifact(t *testing.T) accessor.Artifact {
	t.Helper()

	return accessor.Artifact{Role: "state", Path: filepath.Join(t.TempDir(), "state.json")}
}

// --- ADV-1: the exit maps claim a KILLED process --------------------------
//
// `0025:C3` scopes the two exit maps to a process that RAN AND EXITED: "a
// spawn failure — which has no exit code at all — can never be claimed by
// them however broadly they are written". A child SIGKILLed at the
// deadline, or on parent cancellation, has no exit code either — `Wait`
// synthesizes -1 — yet the binding treats it as an exit and consults the
// maps.
//
// `0025:F1` promises that a runtime refusal carries the 0004 diagnosis
// tuple with class `timeout` or `execution_failure`. A killed child that
// returns a VERDICT, or that establishes every declared key ABSENT, is
// neither: it is a success minted from a corpse.
//
// The exit codes below are authorable and lint-green — nothing in C5's
// `command_output_shape` arms rejects a negative code — so this is a
// declaration a model author can write today.

func TestAdvReadExitAbsentMustNotClaimAKilledChild(t *testing.T) {
	t.Parallel()

	tool := advScript(t, "exec sleep 30\n")
	reader := cmdbind.Reader{
		Accessor: table.Accessor{
			Role:       "state",
			Command:    []string{tool},
			Keys:       []string{"state.phase"},
			ExitAbsent: []int{-1},
		},
		Name:   "killed-read",
		Config: cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	values, unreadable, err := reader.Read(ctx, advArtifact(t), []string{"state.phase"})
	if err != nil {
		return // refusing is the contract
	}
	t.Fatalf("a child KILLED at the deadline was claimed by `exit_absent`: "+
		"values=%+v unreadable=%v err=nil — a killed process has no exit "+
		"code, so `0025:C3`'s exit map may not claim it, and `0025:F1` "+
		"promises a timeout or execution_failure refusal here, never "+
		"established absence", values, unreadable)
}

func TestAdvGateExitVerdictsMustNotClaimAKilledChild(t *testing.T) {
	t.Parallel()

	tool := advScript(t, "exec sleep 30\n")
	gate := cmdbind.Gate{
		Accessor: table.Accessor{
			Role:         "state",
			Command:      []string{tool},
			ExitVerdicts: map[string]string{"-1": "allow"},
		},
		Name:   "killed-gate",
		Config: cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	verdict, reason, err := gate.Gate(ctx, advArtifact(t))
	if err != nil {
		return // refusing is the contract
	}
	t.Fatalf("a gate child KILLED at the deadline was laundered into the "+
		"verdict %q (reason %q): `0025:C3` keeps execution failure out of "+
		"the verdict channel, and a killed process has no exit code for "+
		"`exit_verdicts` to claim", verdict, reason)
}

// TestAdvKilledChildOnParentCancelIsNotAnAnswer is the same defect on the
// path where nothing upstream masks it.
//
// `internal/accessor/executor.go::invokeRead` reclassifies only
// `context.DeadlineExceeded`, so a DEADLINE kill is caught one layer up by
// a second, independent decision. `0025:C6` names exactly that shape as
// the thing a safety property may not rest on — "a coincidence between two
// independent decisions". Parent CANCELLATION (an interrupted CLI) yields
// `context.Canceled`, the upstream check does not fire, and the binding's
// answer stands: every declared key ABSENT, from a tool that answered
// nothing.
func TestAdvKilledChildOnParentCancelIsNotAnAnswer(t *testing.T) {
	t.Parallel()

	tool := advScript(t, "exec sleep 30\n")
	reader := cmdbind.Reader{
		Accessor: table.Accessor{
			Role:       "state",
			Command:    []string{tool},
			Keys:       []string{"state.phase"},
			ExitAbsent: []int{-1},
		},
		Name:   "cancelled-read",
		Config: cmdbind.Config{AllowCommands: true},
	}

	parent, interrupt := context.WithCancel(context.Background())
	// The executor wraps the caller's context in the declared timeout; the
	// timeout is generous here, so the ONLY thing that ends the child is
	// the parent cancellation.
	bounded, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()
	go func() {
		time.Sleep(300 * time.Millisecond)
		interrupt()
	}()

	values, unreadable, err := reader.Read(bounded, advArtifact(t), []string{"state.phase"})
	if err != nil {
		return // refusing is the contract
	}
	t.Fatalf("an INTERRUPTED read returned NO error: values=%+v "+
		"unreadable=%v — a child SIGKILLed on parent cancellation wrote no "+
		"stdout and never exited, so `0025:C3` leaves it no way to carry an "+
		"answer (\"a silent tool is a broken tool, not an answer\", and its "+
		"`exit_absent` map is reserved for a process that ran AND exited), "+
		"and `0025:F1` promises a runtime refusal here. Nothing upstream "+
		"reclassifies `context.Canceled`, so whatever this read returned is "+
		"believed: an empty success is read as a complete answer, and an "+
		"established absence additionally arms `0025:F5`'s read_back seal "+
		"against a tool that was killed before it answered", values, unreadable)
}

// --- ADV-2: the stderr channel is drained only AFTER stdout ---------------
//
// `0025:C4`'s deadline triple is called necessary AND SUFFICIENT, and
// `0025:F4` bounds the residual risk to "a child that IGNORES termination
// at deadline". A child that ignores nothing — that writes a verbose
// diagnostic to stderr and then a well-formed envelope to stdout — is
// outside that bound and must succeed.
//
// It does not. `spawn` drains stdout to completion before it reads stderr
// at all, so a child that fills the 64 KiB stderr pipe buffer before
// closing stdout blocks on write(2) while the parent blocks on read(2).
// Neither moves until the deadline kills the child, and the invocation
// that should have returned in milliseconds burns its whole timeout and
// refuses.
//
// The diagnosis is wrong too, which is what makes it a silent-risk rather
// than a visible one: the refusal says the tool "produced no stdout",
// blaming a tool that produced a correct envelope.
func TestAdvVerboseStderrMustNotDeadlockTheDrain(t *testing.T) {
	t.Parallel()

	// 128 KiB of stderr — twice a pipe buffer — emitted BEFORE stdout is
	// closed. Nothing here ignores termination; the tool simply talks.
	tool := advScript(t, `
i=0
while [ $i -lt 1400 ]; do
  echo "warning: this tool is chatty and says so at length ......................." >&2
  i=$((i+1))
done
printf '{"state.phase":"review"}'
`)
	reader := cmdbind.Reader{
		Accessor: table.Accessor{
			Role:    "state",
			Command: []string{tool},
			Keys:    []string{"state.phase"},
		},
		Name:   "chatty-read",
		Config: cmdbind.Config{AllowCommands: true},
	}

	const deadline = 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	start := time.Now()
	values, unreadable, err := reader.Read(ctx, advArtifact(t), []string{"state.phase"})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("a chatty-but-correct tool refused after %v (deadline %v): "+
			"%v — `0025:C4` calls the deadline triple sufficient and "+
			"`0025:F4` bounds the residue to a child that IGNORES "+
			"termination; this one answered correctly and was starved by "+
			"the parent's sequential pipe drain", elapsed, deadline, err)
	}
	if len(unreadable) != 0 {
		t.Fatalf("a chatty-but-correct tool left %v unreadable after %v",
			unreadable, elapsed)
	}
	if len(values) != 1 || values[0].Value != "review" {
		t.Fatalf("got %+v after %v, want state.phase=review", values, elapsed)
	}
	if elapsed >= deadline {
		t.Fatalf("the read consumed its whole %v deadline (%v) for a tool "+
			"that answers immediately", deadline, elapsed)
	}
}

// --- ADV-3: the process group outlives the invocation --------------------
//
// `0025:C4` makes `Setpgid` half of the deadline mechanism, and its
// stated purpose is that "dropping the group signal ORPHANS A GRANDCHILD
// holding the pipe". The group signal is wired to `cmd.Cancel`, which the
// runtime invokes only when the context ends — so on the SUCCESS path the
// group is never signalled and a grandchild the tool spawned survives the
// CLI, holding the pipe it inherited.
//
// `0025:F4` names the residue it accepts: "a child that ignores
// termination AT DEADLINE can outlive its `timeout` refusal". A grandchild
// outliving a SUCCESSFUL invocation is a different mode and is not
// admitted anywhere in the RDR. Its cost is concrete: the parent's
// `io.Copy(io.Discard, …)` drain does not see EOF until the last holder of
// the pipe exits, so one backgrounded grandchild converts a successful
// read into a full-deadline refusal.
func TestAdvBackgroundGrandchildMustNotStallASuccessfulRead(t *testing.T) {
	t.Parallel()

	// The tool answers correctly and exits 0. It also leaves a background
	// helper running — the shape every `tool &`-style wrapper has.
	tool := advScript(t, `
( sleep 20 ) &
printf '{"state.phase":"review"}'
exit 0
`)
	reader := cmdbind.Reader{
		Accessor: table.Accessor{
			Role:    "state",
			Command: []string{tool},
			Keys:    []string{"state.phase"},
		},
		Name:   "backgrounding-read",
		Config: cmdbind.Config{AllowCommands: true},
	}

	const deadline = 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	start := time.Now()
	values, unreadable, err := reader.Read(ctx, advArtifact(t), []string{"state.phase"})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("a tool that answered correctly and exited 0 refused after "+
			"%v (deadline %v): %v — the invocation is held open by a "+
			"grandchild `0025:C4`'s group mechanism is supposed to bound, "+
			"and `0025:F4` admits no such mode on the SUCCESS path",
			elapsed, deadline, err)
	}
	if len(unreadable) != 0 || len(values) != 1 || values[0].Value != "review" {
		t.Fatalf("values=%+v unreadable=%v after %v, want state.phase=review",
			values, unreadable, elapsed)
	}
	if elapsed >= deadline {
		t.Fatalf("a read whose tool exited 0 immediately took %v, its whole "+
			"%v deadline: the drain waits on a grandchild that inherited "+
			"the pipe", elapsed, deadline)
	}
}

// TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead is the second
// half of ADV-3, which no reported test asserted.
//
// Not stalling the drain and not LEAKING THE PROCESS are different
// properties, and closing the first with a pipe close alone would leave the
// second open. `0025:C4` puts the group signal in the mechanism precisely so
// no grandchild is orphaned, and `0025:F4` admits leakage only for a child
// that outlives its `timeout` refusal — never one that outlives a
// SUCCESSFUL invocation, which is the path this asserts.
func TestAdvBackgroundGrandchildDoesNotSurviveASuccessfulRead(t *testing.T) {
	t.Parallel()

	pidFile := filepath.Join(t.TempDir(), "grandchild.pid")
	// The tool answers correctly and exits 0, leaving a background helper
	// that records its own pid — the shape every `tool &`-style wrapper has.
	tool := advScript(t, `
( sleep 20 ) &
echo $! > `+pidFile+`
printf '{"state.phase":"review"}'
exit 0
`)
	reader := cmdbind.Reader{
		Accessor: table.Accessor{
			Role:    "state",
			Command: []string{tool},
			Keys:    []string{"state.phase"},
		},
		Name:   "leaking-read",
		Config: cmdbind.Config{AllowCommands: true},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, _, err := reader.Read(ctx, advArtifact(t), []string{"state.phase"}); err != nil {
		t.Fatalf("the tool answered correctly and exited 0, yet the read "+
			"refused: %v", err)
	}

	pid := advPIDFrom(t, pidFile)
	// The group signal is sent after the direct child is reaped; give the
	// kernel a moment to tear the grandchild down before asking.
	time.Sleep(200 * time.Millisecond)
	if advAlive(pid) {
		t.Fatalf("the grandchild pid %d survived a SUCCESSFUL invocation: "+
			"`0025:C4` signals the child's process GROUP so no grandchild is "+
			"orphaned, and `0025:F4` admits process leakage only for a child "+
			"that outlives its `timeout` refusal", pid)
	}
}

// advAlive reports whether a pid names a live process.
func advAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// advPIDFrom reads the grandchild pid the probe tool recorded.
func advPIDFrom(t *testing.T, path string) int {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the probe tool recorded no grandchild pid: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatalf("the recorded pid %q is not a number: %v", b, err)
	}
	return pid
}

// TestProcGroupBuildTagsMatchTheRefuseList is the drift guard for the ONE
// duplication the platform split introduces.
//
// `0025:C4` rejects a build constraint on the REFUSAL — "a build constraint
// would make the refusal unbuildable-on-Windows rather than observable, and
// lint must stay platform-neutral in the same binary". The refusal, the
// `goos` var and lint are therefore untagged, and only the process-group
// SYSCALLS are split by `//go:build`, because `syscall.Kill` and
// `SysProcAttr.Setpgid` do not exist on Windows at all.
//
// That split states the same platform set twice: once in a `//go:build`
// list, once in `Unsupported()`'s refuse-list predicate. If they diverge,
// one of two silent breakages follows — a platform that compiles the
// no-op mechanism but is NOT refused would spawn children it can never
// terminate as a group, and a platform refused but compiling the syscall
// half would carry unreachable Unix code. Neither shows up as a build
// failure, so it is asserted here.
func TestProcGroupBuildTagsMatchTheRefuseList(t *testing.T) {
	// Direction 1: every platform the build tags name is refused, and the
	// two lists are the SAME set. `procGroupPlatforms` is defined in both
	// halves of the split, so whichever half compiled is the one checked.
	for _, platform := range cmdbind.ProcGroupPlatforms {
		if !cmdbind.UnsupportedOn(platform) {
			t.Errorf("the build tag names %q but Unsupported() admits it: a "+
				"binding would spawn on a platform with a no-op process "+
				"group and could never terminate the group", platform)
		}
	}

	// Direction 2: nothing is refused that the build tags omit. Iterating
	// the known GOOS values keeps this honest — a refuse-list entry added
	// to Unsupported() without a matching build tag is caught here.
	for _, platform := range []string{
		"aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios",
		"js", "linux", "netbsd", "openbsd", "plan9", "solaris", "wasip1",
		"windows",
	} {
		tagged := slices.Contains(cmdbind.ProcGroupPlatforms, platform)
		if refused := cmdbind.UnsupportedOn(platform); refused != tagged {
			t.Errorf("platform %q: Unsupported() = %v but the build-tag "+
				"refuse-list membership = %v; the //go:build list and the "+
				"runtime predicate have drifted", platform, refused, tagged)
		}
	}

	// Direction 3: the half actually compiled agrees with the running
	// platform's own verdict. This is the check that fires if the two
	// //go:build lines stop partitioning GOOS — if both files or neither
	// matched, the package would not build, but a WRONG partition builds
	// fine and is only observable as this mismatch.
	if cmdbind.ProcGroupSupported == cmdbind.Unsupported() {
		t.Errorf("this binary compiled the procGroupSupported=%v half but "+
			"Unsupported() = %v on %s; the mechanism and the refusal "+
			"disagree about the running platform",
			cmdbind.ProcGroupSupported, cmdbind.Unsupported(), runtime.GOOS)
	}
}
