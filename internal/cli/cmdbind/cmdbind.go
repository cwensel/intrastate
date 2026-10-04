// Package cmdbind is RDR 0025's command-backed `accessor.Binding` family:
// the declared-argv carrier that lets a model delegate a read, a gate, or
// a write to an established external tool.
//
// It is a sibling of `internal/cli/flowbind`, which binds the same three
// interfaces over intrastate's own flat-JSON artifact. Selection between
// them is by CARRIER, at the registry: an entry declaring `path` builds a
// file binding, one declaring `command` builds one of these, and exactly
// one is declared per entry (`0025:C1`).
//
// The contract this package implements, in one place:
//
//   - the executed argv is the DECLARED vector after whole-element
//     `{artifact}` substitution, with no shell and no other rewriting
//     (`0025:C2`);
//   - per-invocation data crosses on stdin as a JSON object of strings and
//     results return on stdout in the shapes `0025:C3` fixes, plus the raw
//     single-value read mode and the two declared exit maps;
//   - the child runs in its own process group under the executor's
//     deadline, with an allowlisted environment and bounded stdout
//     (`0025:C4`);
//   - execution requires the `--allow-commands` opt-in, checked BEFORE any
//     spawn so no execution path can reach one without it (`0025:C6`).
package cmdbind

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// ArtifactPlaceholder is the closed placeholder vocabulary, complete at one
// member in v1: replaced WHOLE-ELEMENT by the caller-bound artifact path
// for the entry's declared role (`0025:C2`).
const ArtifactPlaceholder = "{artifact}"

// StdoutCap bounds the child's stdout. An overflow is an execution failure,
// never a truncation, and a stdout truncated at the cap is NON-empty and so
// is never exit-mapped (`0025:C4`).
const StdoutCap = 1 << 20 // 1 MiB

// StderrTailCap bounds the stderr tail carried on `accessor.Refusal.Detail`
// (`0025:C4`).
const StderrTailCap = 4 << 10 // 4 KiB

// WaitDelay bounds the stdin write and the post-kill pipe drain, so a
// non-reading or slow-draining child cannot hang the CLI past the deadline
// (`0025:C4`; the drain half is now `0026:C1` `bound:`).
//
// It is the ONE bound for every "the CLI waits past the child's exit" case
// — Cmd's stdin write, Cmd's post-Cancel wait, and this package's own drain
// join — and it is in the contract, because `0026:C1` `bound:` publishes a
// committed per-invocation return of `timeout + 2·WaitDelay`. Moving it is
// therefore a contract change, not a tuning knob.
const WaitDelay = 500 // milliseconds

// DrainGrace is how far AHEAD of now the drain's read deadline is re-armed
// once the join's bound has expired (`0026:C1` `bound:`).
//
// It is NOT a second bound. It answers a different question — how far ahead
// the deadline must sit for the kernel to deliver bytes already buffered,
// since Go checks the deadline BEFORE attempting the syscall, so a deadline
// of NOW short-circuits with n=0 and never looks at the pipe. It is never
// model-declarable and it cannot move the refuse/accept boundary: it is an
// offset on the single existing timer's expiry. 50 ms clears the measured
// scheduling tail (~28 ms under 4x-core saturation) with margin while
// costing 10% of one bound.
const DrainGrace = 50 // milliseconds

// ProtocolVersion rides out-of-band in the child env as
// `INTRASTATE_PROTOCOL`, keeping the stdout map flat (`0025:C3`).
const ProtocolVersion = "1"

// The `INTRASTATE_*` overlay names. An entry `env` key or `env_pass` name
// matching the reserved prefix is a load-time defect, so the overlay is
// never shadowed silently (`0025:C4`, `0025:C5`).
const (
	EnvRole       = "INTRASTATE_ROLE"
	EnvCapability = "INTRASTATE_CAPABILITY"
	EnvAccessor   = "INTRASTATE_ACCESSOR"
	EnvProtocol   = "INTRASTATE_PROTOCOL"

	// EnvReservedPrefix is the literal, case-sensitive reserved prefix.
	EnvReservedPrefix = "INTRASTATE_"
)

// AllowlistedVars is the parent-environment allowlist: nothing else is
// inherited (`0025:C4`, A7).
var AllowlistedVars = []string{"PATH", "HOME", "TMPDIR", "LANG"}

// AllowlistedPrefix is the ONE prefix rule — a literal match on `LC_`.
// `env_pass` names whole variables and admits no pattern (`0025:C4`).
const AllowlistedPrefix = "LC_"

// goos is the platform the binding refuses on. It is a package-level var so
// the refusal is testable on Unix CI; overriding it weakens nothing at
// runtime, unlike C4's deadline triple, which carries no injection seam
// (`0025:C4` platform).
var goos = runtime.GOOS

// drainStartDelay withholds each drain goroutine's first read. It is the
// ONE stall seam the drain scenarios share, and it is nil — the zero
// duration — in production: portable Go rather than a build tag, and never
// a sleep on a shipped path.
var drainStartDelay time.Duration

// nonPollableForTest forces the creation-time pollability probe to see a
// read end that refuses a deadline, so F5's condition is simulated AT THE
// CHECK rather than reproduced. No real `os.Pipe` read end refused a
// deadline on either supported OS, so a fixture that reproduced it would be
// fabricating a host condition.
var nonPollableForTest bool

// Unsupported reports whether the running platform is refuse-listed. The
// predicate is refuse-listed, not allow-listed, so every Unix that supports
// the process-group mechanism — the BSDs included — is admitted without
// enumerating it (`0025:C4`).
func Unsupported() bool {
	return goos == "windows" || goos == "js" || goos == "plan9"
}

// Config is what a command binding needs beyond the model entry: the base
// directory a separator-bearing argv0 resolves against (`0025:C2`), and the
// execution gate (`0025:C6`).
type Config struct {
	// BaseDir is `filepath.Dir` of the ABSOLUTIZED model path, supplied by
	// the CLI caller that opened the file. It is empty for a model with no
	// source file, which makes a separator-bearing argv0 a construction-time
	// refusal — never a fall back to the process cwd.
	BaseDir string
	// AllowCommands is the `--allow-commands` gate. With it unset, every
	// command invocation refuses before spawn (`0025:C6`).
	AllowCommands bool
}

// --- the shared invocation -----------------------------------------------

// invocation is one bounded child run: what the binding sends, what the
// child answered, and how it ended.
type invocation struct {
	stdout   []byte
	stderr   string
	exitCode int
	// exited reports that the process RAN and exited OF ITS OWN ACCORD with
	// a real exit code. The two exit maps are consulted only for such a
	// process, so a spawn failure — which has no exit code at all — can
	// never be claimed by them however broadly they are written
	// (`0025:C3`).
	//
	// A process the parent SIGNALLED is the same case: `Wait` reports a
	// SIGKILL as an `*exec.ExitError` whose `ExitCode()` is the synthetic
	// -1, but a killed process no more exited than a spawn failure did, and
	// a lint-green `exit_absent = [-1]` / `exit_verdicts = { "-1" = ... }`
	// would otherwise convert a child killed BEFORE IT ANSWERED into
	// established absence or a verdict the model never decided
	// (`0025:F1`/`0025:F2`/`0025:F3`). The distinction is made HERE, at the
	// binding, and not by the executor's independent `ctx.Err()` check —
	// `0025:C6` forbids a safety property that rests on "a coincidence
	// between two independent decisions", and the executor's check does not
	// cover `context.Canceled` at all.
	exited bool
	// signaled reports that the child was terminated by a signal, which is
	// what the refusal text says instead of the meaningless "exited -1".
	signaled bool
	// nonPollable records that a read end refused a deadline at pipe
	// creation (`os.ErrNoDeadline`), so the join on THIS invocation could
	// not be bounded — F5's host condition, decided before any child
	// existed. `cmdbind` has no logger and gains none, so the condition
	// rides here and surfaces as a `Detail` line on whatever refusal this
	// invocation produces.
	nonPollable bool
}

// exitStatus renders the DIRECT child's own end for the held-pipe reason,
// composed from the fields `spawn` already populates from `Wait` — so no
// new invocation field carries it (`0026:C1` `refusal:`).
func (inv invocation) exitStatus() string {
	switch {
	case inv.signaled:
		return "killed by signal"
	case inv.exited:
		return "exited " + strconv.Itoa(inv.exitCode)
	default:
		return "did not exit"
	}
}

// drainReport is one drain's terminal condition: the ONE three-valued,
// totally ordered report `0026:C1` `precedence:` requires, and not a pair
// of booleans or a bare error. The three states are mutually exclusive and
// their order is fixed — overflow outranks held, held outranks whole — so
// a caller cannot re-derive the precedence differently from this package.
type drainReport int

const (
	// drainWhole is EOF: every writer is gone and the output is whole, as
	// today. It is the LOWEST rank.
	drainWhole drainReport = iota
	// drainHeld is `os.ErrDeadlineExceeded` on the drain's final read under
	// the grace: a writer the group signal could not reach still holds the
	// pipe. It outranks whole however many bytes the grace delivered — the
	// predicate is the REPORTED condition, never the byte count.
	drainHeld
	// drainOverflow is the read passing `StdoutCap`. It is the HIGHEST
	// rank: the bytes were read, so the pipe's terminal state is not what
	// the invocation turns on.
	drainOverflow
)

// spawn performs the whole pre-spawn refusal ladder and, past it, one
// bounded invocation of the declared command.
//
// Every refusal it returns is a typed `*accessor.ExecError`, which is the
// channel the seam's bare `error` slot gives a command binding for its
// stderr tail (`0025:C4`).
func spawn(
	ctx context.Context, acc table.Accessor, name string,
	capability accessor.Capability, cfg Config, art accessor.Artifact,
	stdin []byte,
) (invocation, error) {
	// --- the pre-spawn refusal ladder -----------------------------------
	// Every arm below refuses BEFORE any child exists, which is what the
	// absence-of-a-spawn oracles assert.

	// C6 — the execution gate. A model-declared command is code that runs
	// when the model is used, so execution requires an opt-in outside the
	// model.
	if !cfg.AllowCommands {
		return invocation{}, gateRefusal(name)
	}

	// C4 platform — refuse-listed, so every Unix that supports the
	// process-group mechanism is admitted without being enumerated.
	if Unsupported() {
		return invocation{}, refuse("command entries are unsupported on " +
			goos + "; the accessor `" + name + "` declares a command")
	}

	// C1's runtime residue arm. A binding built from an in-memory model
	// never passed the loader, so the carrier-less entry is refused here
	// rather than falling through to a `Path: ""` file binding — which
	// would read every declared key as absent and confirm an unapplied
	// write.
	if len(acc.Command) == 0 {
		return invocation{}, refuse("the accessor `" + name +
			"` declares neither `path` nor `command`; a carrier-less entry " +
			"is malformed and cannot be invoked")
	}

	argv, err := substitute(acc.Command, art)
	if err != nil {
		return invocation{}, err
	}

	argv0, err := resolveArgv0(argv[0], cfg.BaseDir)
	if err != nil {
		return invocation{}, err
	}

	// --- the bounded invocation -----------------------------------------
	cmd := exec.CommandContext(ctx, argv0, argv[1:]...)
	cmd.Env = childEnv(acc, name, capability)
	cmd.Stdin = bytes.NewReader(stdin)

	// `0025:C4`'s deadline triple, and it was never "necessary and
	// sufficient": dropping the group signal orphans a grandchild holding
	// the pipe, and dropping WaitDelay hangs `Wait` on it and unbounds the
	// stdin write — but a writer OUTSIDE the group is reachable by none of
	// the three, which is the gap `0026:C1` closes by bounding the drain
	// join rather than by widening the kill.
	setProcGroup(cmd)
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		// The GROUP, not only the direct child.
		if kerr := killGroup(cmd.Process.Pid); kerr != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
	cmd.WaitDelay = WaitDelay * time.Millisecond

	// The two output channels are carried on pipes this function OWNS,
	// rather than on `cmd.StdoutPipe`/`cmd.StderrPipe`, because the close
	// ordering is load-bearing: `Cmd.Wait` closes the pipes it owns as soon
	// as the direct child is reaped, which would truncate a read still in
	// flight. Owning them lets the group be released FIRST and the drains
	// then run to whole output WITHIN THE BOUND (`0026:C1` `precedence:`,
	// the successor to `0025:C4`'s "run to a true EOF" — the drains no
	// longer promise an unbounded wait for EOF).
	outR, outW, err := os.Pipe()
	if err != nil {
		return invocation{}, wrap("", err)
	}
	defer func() { _ = outR.Close() }()
	errR, errW, err := os.Pipe()
	if err != nil {
		_ = outW.Close()
		return invocation{}, wrap("", err)
	}
	defer func() { _ = errR.Close() }()

	// F5's pollability probe, run ONCE at pipe creation and never at the
	// join: a far-future deadline, immediately cleared, whose error is the
	// signal. Deciding it here means the fallback is decided before any
	// child exists, so a non-pollable host takes the old unbounded join by
	// a recorded choice rather than by discovering it mid-drain.
	nonPollable := !pollable(outR) || !pollable(errR)
	// An `*os.File` on these fields is handed to the child directly, with no
	// copying goroutine and no entry in the runtime's parent-pipe set; the
	// parent's copy of each write end is closed by `Start`.
	cmd.Stdout = outW
	cmd.Stderr = errW

	serr := cmd.Start()
	_ = outW.Close()
	_ = errW.Close()
	if serr != nil {
		return invocation{}, wrap("", serr)
	}

	// BOTH channels drain CONCURRENTLY. Draining stdout to completion before
	// touching stderr deadlocks a merely VERBOSE tool: a child that fills the
	// 64 KiB stderr pipe buffer before closing stdout blocks in `write(2)`
	// while the parent blocks in `read(2)`, and neither moves until the
	// deadline — well below the 1 MiB cap, so the cap does not bound it.
	// `0025:F4` bounds the accepted residue to "a child that IGNORES
	// termination at deadline"; this child ignores nothing (`0025:C4`).
	//
	// The stdout bound is enforced by BOUNDING THE READ, so an unbounded
	// child cannot exhaust memory. One byte past the cap is read on
	// purpose: overflow is DETECTED here rather than silently truncated.
	// The stderr channel is read under the SAME bound so a verbose tool
	// cannot exhaust memory, and `tail` then keeps the LAST 4 KiB — which is
	// the diagnosis a caller wants when a tool says a lot before it fails.
	var stdout, stderr []byte
	var outReport, errReport drainReport
	// `stop` is the drains' shared HALT — the terminator for a drain that is
	// making progress into output no caller will ever parse. `0026:C1` does
	// not name it, so it is recorded as deviation D9; this comment and
	// `readBounded`'s are its durable record.
	//
	// It is a STATE PREDICATE, not a second bound. It introduces no constant
	// and reads no clock: it answers "is this invocation's outcome already a
	// refusal", never "how long to wait". That is what keeps `0026:C1`
	// `bound:`'s "`WaitDelay` is the ONE bound" and its "`DrainGrace` is NOT
	// a second bound" literally true with the halt shipped. It is also the
	// test `0026`'s own LBD fixes for a new mechanism: `DrainBound` was
	// refused for answering the SAME question with a second value, and this
	// answers a different question again.
	//
	// It follows from `0026:C1` `precedence:` (b) rather than adding to it.
	// That clause requires the parent "MUST NOT read `inv.stdout`/
	// `inv.stderr` until every drain goroutine has returned" and forbids a
	// design "in which the parent proceeds past a live drain" — so once any
	// drain has established a refusal, the only way to honour the clause is
	// for the sibling to return too. The halt is that consequence made
	// mechanical, not new policy.
	//
	// It is a different signal from the grace. The parent's deadline-set is
	// the MARK: the bound expired and each drain's FINAL read is owed.
	// `stop` says something else — the invocation's outcome is ALREADY a
	// refusal, so no further byte from any drain can change it.
	//
	// It is raised by a drain that has itself REPORTED a refusing condition,
	// and never on elapsed time. That distinction is what keeps this out of
	// the mechanism `0026:C1` `precedence:` rejects by name: "A single
	// absolute deadline would instead bound the whole remaining tail against
	// the clock, which turns a slow-arriving tail into a false held-pipe
	// report on a pipe with no writer at all." Nothing here reads a clock.
	// On a drain with no holder — MVV row 6(a), 6(b), S9 — no sibling ever
	// reports held, so `stop` is never raised and the paced tail runs to EOF
	// however long it takes, which is the guarantee the clause fixes.
	//
	// The residue is the flip side of "raised only by a drain that itself
	// reported": when an escapee trickles BOTH pipes with every gap under
	// `DrainGrace`, neither drain reaches a refusing condition, so no drain
	// raises `stop` and the join stays unbounded. That is an ADMITTED,
	// named liveness residue — deviation D10; see `readBounded`.
	var stop atomic.Bool
	// `reaped` gates the halt's OBSERVATION on the direct child being gone.
	// Both drains start BEFORE the wait below, and a drain raises the halt
	// on any read error once it has overflowed — INCLUDING EOF. So a child
	// that overflows stdout and then closes it raises the halt while it is
	// STILL LIVE, and a sibling honouring it there would abandon a pipe that
	// live child is still writing. The child then blocks in `write(2)`, the
	// wait never completes, and the invocation spends its whole deadline and
	// reports a timeout instead of the overflow refusal it owes. That is
	// exactly the serialization `0026:C1` `whole:` forbids ("the drains run
	// CONCURRENTLY and are never stalled") and REQ-52 states negatively.
	//
	// It costs the halt nothing in the shapes it exists for. Every one of
	// them (S9's companion, deviations D9 and D10, MVV row 6) is a `setsid`
	// ESCAPEE holding an inherited pipe, so the direct child is already
	// gone — the refusal lands at ~`WaitDelay + DrainGrace`, necessarily
	// after the wait completed. The gate therefore narrows the halt to
	// precisely the window in which it is unsafe.
	var reaped atomic.Bool
	var drains sync.WaitGroup
	drains.Add(2)
	go func() {
		defer drains.Done()
		if drainStartDelay > 0 {
			time.Sleep(drainStartDelay)
		}
		stdout, outReport = readBounded(outR, StdoutCap+1, &stop, &reaped)
	}()
	go func() {
		defer drains.Done()
		if drainStartDelay > 0 {
			time.Sleep(drainStartDelay)
		}
		// Stderr's overflow does NOT refuse — it is TAILED rather than
		// capped, so `spawn`'s overflow arm consults only the stdout report.
		// `0026:C1` `precedence:` scopes the overflow rank to a drain that
		// "REFUSES the existing overflow error", so ranking a HELD stderr as
		// overflow would discard the condition rather than outrank it and
		// the invocation would return with no error at all — the withheld
		// refusal `0026:C1` `residue:` does not admit.
		stderr, errReport = readBounded(errR, StdoutCap, &stop, &reaped)
	}()

	werr := cmd.Wait()
	// The direct child is gone, so from here a halt raised by either drain
	// can no longer strand the other on a pipe a live writer is filling.
	// See the declaration of `reaped` above.
	reaped.Store(true)

	// The direct child has been reaped. Signal the GROUP now, on the success
	// path as much as the failure path: `cmd.Cancel` fires only when the
	// context ends, so without this a grandchild the tool backgrounded — the
	// shape every wrapper script has — survives the CLI still holding the
	// stdout/stderr pipes it inherited, and the drain above never sees EOF.
	// `0025:C4` states the group signal's purpose as exactly this ("dropping
	// the group signal ORPHANS A GRANDCHILD holding the pipe"), and
	// `0025:F4` admits process leakage only for a child that outlives its
	// `timeout` refusal — never one that outlives a SUCCESSFUL invocation.
	//
	// It is safe on the success path precisely BECAUSE it runs after the
	// wait: the direct child is already gone, so nothing legitimate is
	// killed — whatever remains in the group is the orphan the clause names.
	reapGroup(cmd.Process)

	// Only now, with every write end of both pipes closed, does a drain with
	// no unreachable writer reach EOF. But EOF is not owed: a writer outside
	// the group survives `reapGroup`, and waiting for it is the hang this
	// record exists to end. So the join is BOUNDED — a bounded wait outranks
	// a whole read (`0026:C1` `precedence:`).
	//
	// ONE timer of `WaitDelay` covers BOTH drains. When it fires the parent
	// sets a read deadline of now + `DrainGrace` on its OWN read ends: that
	// deadline IS the join's cross-goroutine mark, and it needs no lock of
	// its own — an `os.File`'s poller is internally synchronised, so setting
	// a deadline against a `Read` in flight is sanctioned. Each drain's
	// final read then reports the PIPE's state rather than the clock, and
	// re-arms the grace before every subsequent read itself.
	//
	// The mark is the JOIN's mechanism and it is not the DRAIN's terminator.
	// A drain still making progress never trips the deadline at all, so the
	// mark alone cannot end it; what ends it is the sibling halt `stop`
	// above, which is a shared flag and is deliberately so (deviation D9).
	//
	// The join stays a REAL join either way: the bound governs how long the
	// parent waits before setting the deadline, never whether it waits for
	// the goroutines to finish. A drain ended by the deadline still runs its
	// `Done`, so `drains.Wait()` below is the happens-before edge between
	// the drains' writes to the byte slices and the read of them here.
	joined := make(chan struct{})
	go func() {
		drains.Wait()
		close(joined)
	}()
	if nonPollable {
		// F5: a read end that refuses a deadline cannot be unblocked, so
		// there is nothing to bound the join WITH. The old unbounded wait
		// is the honest fallback, taken by a choice recorded at the check.
		<-joined
	} else {
		timer := time.NewTimer(WaitDelay * time.Millisecond)
		select {
		case <-joined:
			timer.Stop()
		case <-timer.C:
			grace := time.Now().Add(DrainGrace * time.Millisecond)
			_ = outR.SetReadDeadline(grace)
			_ = errR.SetReadDeadline(grace)
			<-joined
		}
	}

	// `inv` is built BEFORE the arms below and every error return carries
	// this populated value, so `inv.stdout` and `inv.exitCode` are readable
	// on every error path. That is not a licence: the obligation `0026:C1`
	// `refusal:` states is that EVERY caller checks `err` before touching
	// `inv`, and the two error returns below — the drain/`werr` refusal and
	// the overflow refusal — are where the returned stdout may be a stream
	// whose live writer the CLI could not reach. Reading it ahead of `err`
	// derives a verdict from unproven output, which is the silent
	// wrong-answer class this record exists to prevent. Nothing structural
	// enforces this; a new caller must check `err` first.
	inv := invocation{
		stdout: stdout, stderr: tail(stderr), nonPollable: nonPollable,
	}
	var werrRefusal error
	var ee *exec.ExitError
	switch {
	case werr == nil:
		inv.exited = true
	case errors.As(werr, &ee):
		// `Wait` reports a SIGNALLED child as an ExitError too, with the
		// synthetic exit code -1. A killed process did not exit, so the exit
		// maps must not be able to claim it (`0025:C3`); the deadline itself
		// is classified from ctx.Err() by the caller and never from this
		// error (`0025:C4`).
		if status, held := ee.Sys().(syscall.WaitStatus); held && status.Signaled() {
			inv.signaled = true
			inv.exitCode = ee.ExitCode()
			break
		}
		inv.exited = true
		inv.exitCode = ee.ExitCode()
	default:
		// The non-ExitError arm — `exec.ErrWaitDelay` from a child that
		// never read its stdin, per `0026:C1` `stdin:`. It is HELD, not
		// returned: the drain conditions are selected first below, because a
		// held pipe names the helper that held it and `ErrWaitDelay` does
		// not. Both are `execution_failure`, so this orders the reason and
		// never the class.
		//
		// The direct child EXITED 0, and that is a stdlib guarantee rather
		// than an inference from this fixture: `os/exec` returns
		// `ErrWaitDelay` only from `awaitGoroutines`, which runs after
		// `Process.Wait()` has already returned a nil error, and `Wait` is
		// nil-erroring only when `state.Success()` holds. The doc says it
		// outright — ErrWaitDelay is returned "if the process exits with a
		// successful status code but its output pipes are not closed before
		// ... WaitDelay expires". So `exited`/`exitCode` are set here for
		// the same reason the ExitError arm sets them: `exitStatus()`
		// composes them into the held-pipe `Detail` (`0026:C1` `refusal:`),
		// and this is the one arm where `ErrWaitDelay` and a held pipe
		// coincide. Leaving them zero renders "did not exit" for a child
		// that exited cleanly, pointing a reader at the wrong process.
		inv.exited = true
		inv.exitCode = 0
		werrRefusal = werr
	}

	// The drain conditions, in the total order `0026:C1` `precedence:`
	// fixes: overflow outranks held, held outranks whole. Overflow wins
	// because the bytes were READ, which makes the pipe's terminal state not
	// what the invocation turns on.
	if outReport == drainOverflow || len(inv.stdout) > StdoutCap {
		return inv, wrap(inv.detail(""), errors.New("the child's stdout "+
			"exceeded the "+strconv.Itoa(StdoutCap)+" byte bound"))
	}
	// The fold takes the two drains' REPORTED conditions as its input and
	// derives nothing from byte counts or from closed-over state, and
	// returns the names in the fixed order stdout, stderr.
	if held := heldPipes(outReport, errReport); len(held) != 0 {
		return inv, wrap(inv.detail(heldReason(held, inv.exitStatus())),
			&accessor.HeldPipeError{Held: held, ExitStatus: inv.exitStatus()})
	}
	if werrRefusal != nil {
		return inv, wrap(inv.detail(""), werrRefusal)
	}
	return inv, nil
}

// gateRefusal is the C6 refusal: command execution without the opt-in.
// It is shared by `spawn` and by the `steps` writer, which checks it
// before selecting any step so a write that would run nothing still
// refuses without the gate (kata q14r).
//
// It is a refusal about the REQUEST, so it is built with `refuseRequest`
// and takes the exit-2 group. JDR 0001 §D10 rule 1 reserves exit 3 for an
// environment that could not be consulted; a missing opt-in is a property
// of the invocation, and re-issuing it unchanged can never succeed. This
// is the reasoning `0028:C1.3` applies to the same missing flag on the
// edit writer's command-backed read-back, so both gates now route alike.
// `0025:C6`'s CLASS is unchanged: `ErrDeclaredRequest` rides the existing
// `Err` slot and the refusal stays `execution_failure` (JDR 0003 §D3(b)).
func gateRefusal(name string) error {
	return refuseRequest("command execution requires the " +
		"`allow_commands` opt-in (--allow-commands); the accessor `" +
		name + "` declares a command and none was given")
}

// heldPipes folds the two drains' reported terminal conditions into the
// held-pipe names, in the fixed order stdout, stderr, so `len(held) != 0`
// is the consulted predicate and the single-pipe form is the bare name
// (`0026:C1` `precedence:`).
func heldPipes(out, err drainReport) []string {
	var held []string
	if out == drainHeld {
		held = append(held, "stdout")
	}
	if err == drainHeld {
		held = append(held, "stderr")
	}
	return held
}

// heldReason is the held-pipe reason `0026:C1` `refusal:` fixes: the pipe(s)
// held, the bound, the direct child's exit status, and the remediation. The
// single-pipe form is the BARE name; both is the two names comma-separated
// in the fixed order.
func heldReason(held []string, status string) string {
	return "the child's " + strings.Join(held, ", ") + " stayed held past " +
		"the " + strconv.Itoa(WaitDelay) + " ms drain bound (" + status +
		"); close or redirect the helper's inherited stdio"
}

// detail composes this invocation's refusal `Detail` in the ONE slot: the
// held-pipe reason LEADING, then F5's host condition when the join could
// not be bounded, then the stderr tail collected up to the bound
// (`0026:C1` `refusal:`). The reason leads because it is the fact the
// author acts on; the tail is the diagnosis behind it.
func (inv invocation) detail(reason string) string {
	parts := make([]string, 0, 3)
	if reason != "" {
		parts = append(parts, reason)
	}
	if inv.nonPollable {
		// F5 is carried HERE and not on a log line: `cmdbind` has no logger
		// and gains none, so the condition recorded at pipe creation
		// surfaces as a `Detail` line on whatever refusal this invocation
		// produces.
		parts = append(parts, "this host's pipes refused a read deadline "+
			"(os.ErrNoDeadline), so the drain join for this invocation was "+
			"not bounded")
	}
	if inv.stderr != "" {
		parts = append(parts, inv.stderr)
	}
	return strings.Join(parts, "\n")
}

// reapGroup SIGKILLs the child's process group after the direct child has
// been reaped, releasing any grandchild still holding the inherited pipes
// (`0025:C4`).
//
// The negative-pid form addresses the GROUP, and `Setpgid` with no `Pgid`
// makes the child its own group leader, so the group id is the child's pid.
// The guard is not cosmetic: a pid of 0 or 1 would address the CALLER's
// group or init, so a process record without a usable pid signals nothing.
func reapGroup(p *os.Process) {
	if p == nil || p.Pid <= 1 {
		return
	}
	// ESRCH — the ordinary case, an empty group — is nothing to report.
	_ = killGroup(p.Pid)
}

// pollable reports whether a read end accepts a deadline. The probe is a
// FAR-FUTURE deadline, immediately cleared, whose error is the signal
// (`os.ErrNoDeadline` on a fd adopted without the poller registration a
// real `os.Pipe` end carries). It runs once at pipe creation, so F5's
// fallback is decided before any child exists rather than discovered at the
// join.
func pollable(r *os.File) bool {
	if nonPollableForTest {
		return false
	}
	if err := r.SetReadDeadline(time.Now().Add(time.Hour)); err != nil {
		return false
	}
	return r.SetReadDeadline(time.Time{}) == nil
}

// readBounded drains r up to limit bytes and reports its terminal
// condition. The bound is on what is KEPT: the remainder is discarded
// rather than left in the pipe, so a child is never blocked writing into a
// full one.
//
// The loop is explicit rather than `io.ReadAll` because `0026:C1`
// `precedence:` needs a per-read point to re-arm the deadline, and
// `io.ReadAll` owns its own loop. Three behaviours the previous body had
// are preserved exactly: the read stays bounded at the caller's limit — for
// stdout that is `StdoutCap+1`, so exactly the cap is a value and one byte
// more is a detected overflow; the discard-remainder exit carries the
// deadline itself, since it is the same loop; and the PRE-timer path is
// deadline-free, because nothing here arms a deadline until the parent has
// armed one.
//
// The parent's deadline-set at the bound is the mark. The first
// `os.ErrDeadlineExceeded` this loop sees is therefore not a verdict — it
// is the notification that the grace is on, and the loop re-arms
// `now + DrainGrace` and takes ONE more read. It is that FINAL read, taken
// under the drain's own grace, whose result is the report: bytes still
// buffered are delivered and the loop keeps going, EOF is whole output as
// today, and another deadline is a HELD pipe. Refusing on "no EOF after a
// final bounded read" and never on elapsed time is what keeps the grace
// bounding the IDLE GAP between reads rather than the size or the total
// duration of the tail — a single absolute deadline would truncate a slow
// tail on a pipe with no writer at all.
//
// There is NO clock here beyond the grace itself — no ceiling, no absolute
// deadline, nothing keyed on the tail's total duration. `0026:C1` `bound:`
// names `WaitDelay` as "the ONE bound" and states `DrainGrace` "is NOT a
// second bound", so a third constant is not the implementer's to mint.
//
// A writer that keeps delivering bytes FASTER than the grace but never
// reaches EOF therefore cannot be ended by the idle-gap predicate — it
// re-arms forever. What ends it is `stop`: a sibling drain that has already
// reported a refusing condition raises it, because from that moment the
// invocation's outcome is fixed as a refusal and no byte this drain could
// still deliver would be parsed. `stop` is checked at the TOP of the loop
// so an in-flight read is not credited past it, and the drain reports the
// condition it has: held, because it did not reach EOF.
//
// `stop` is a STATE PREDICATE and NOT a second bound — it mints no constant
// and reads no clock, so `0026:C1` `bound:`'s "`WaitDelay` is the ONE
// bound" stays literally true. It is the mechanical consequence of
// `0026:C1` `precedence:` (b)'s "MUST NOT read ... until every drain
// goroutine has returned": once a refusal is established, the sibling must
// return for the parent to honour that clause at all. `0026:C1` does not
// name this terminator, so it is recorded as deviation D9.
//
// ADMITTED RESIDUE (deviation D10). The halt is raised ONLY by a drain that
// has itself reported a refusing condition, so a shape in which no drain
// reports leaves it unraised. That shape exists: an escapee trickling BOTH
// pipes with every idle gap under `DrainGrace`. Both drains keep their
// deadlines ahead of them, neither ever reports, `stop` is never raised and
// the join above stays UNBOUNDED — the old hang, in a shape no fixture
// exercises. It is admitted rather than closed because both mechanisms that
// would close it are already refused by this record. A duration bound on the
// drain is the total-duration bound `precedence:` rejects by name (and is
// exactly what the removed liveness ceiling was). A hard close of the read
// ends is Go's answer (`os/exec`'s `awaitGoroutines` -> `closeDescriptors`)
// and surrenders the stderr tail naming WHICH helper held the pipe — the
// sole justification for this record's deliberate divergence from that prior
// art. So it takes F5's shape: recorded at the check rather than discovered
// at the join, an admitted liveness residue rather than a silent one.
//
// The overflow rank is keyed on the REFUSAL boundary, `StdoutCap`, and not
// on the caller's read limit. `0026:C1` `precedence:` gives the rank its
// reason — "a drain past `StdoutCap` REFUSES the existing overflow error
// whatever its terminal condition, since the bytes were read and the pipe's
// state is then not what the invocation turns on" — so the rank exists
// BECAUSE overflow refuses, and it is that refusal's own boundary which
// scopes it. `spawn` refuses on `len(inv.stdout) > StdoutCap`, so this reads
// the same predicate: exactly the cap is a value and one byte past it is the
// detected overflow.
//
// Keying it on `len(kept) >= limit` instead silently launders the STDERR
// drain, whose limit is exactly `StdoutCap` rather than `StdoutCap+1`
// because stderr is TAILED rather than capped and has no overflow refusal to
// take. A held stderr that reached that limit would rank overflow; `spawn`'s
// overflow arm consults only the stdout report and `heldPipes` counts only
// held, so neither arm would consult it and the invocation would return with
// NO error at all — the withheld refusal `0026:C1` `residue:` does not
// admit. The limit stays a memory bound on both drains: past it the bytes
// are still read and discarded rather than left in the pipe.
//
// The optional `stop` is the sibling drains' shared halt, and the optional
// second flag is `spawn`'s reaped-child gate for OBSERVING it. Both are
// variadic only so the loop keeps the two-argument shape it has always had
// at a call site that has no sibling to be halted by; a site that passes the
// halt without the gate observes it unconditionally, as before.
func readBounded(
	r *os.File, limit int, stop ...*atomic.Bool,
) ([]byte, drainReport) {
	var halt, reaped *atomic.Bool
	if len(stop) != 0 {
		halt = stop[0]
	}
	if len(stop) > 1 {
		reaped = stop[1]
	}
	// A drain may OBSERVE the halt only once the direct child has been
	// reaped. RAISING it stays ungated — a drain always reports its own
	// condition — but honouring it while the child is still live abandons a
	// pipe that child may still be writing, which blocks it in `write(2)`
	// and hangs `Wait`. See the gate's rationale in `spawn`, `0026:C1`
	// `whole:` and REQ-52.
	halted := func() bool {
		if halt == nil || !halt.Load() {
			return false
		}
		return reaped == nil || reaped.Load()
	}
	var kept []byte
	buf := make([]byte, 32<<10)
	graced := false
	overflowed := func() bool { return len(kept) > StdoutCap }
	for {
		// The sibling halt: the terminator for a drain making progress into
		// output no caller will parse. It is a STATE PREDICATE — "the
		// invocation's outcome is already a refusal" — and NOT a second
		// bound: no constant, no clock. `0026:C1` `bound:`'s "`WaitDelay` is
		// the ONE bound" therefore still holds literally. It follows from
		// `0026:C1` `precedence:` (b)'s "MUST NOT read ... until every drain
		// goroutine has returned", which cannot be honoured past an
		// established refusal unless the sibling returns. `0026:C1` names no
		// such terminator, so it is recorded as deviation D9.
		//
		// It is checked BEFORE the grace re-arm and independently of it. A
		// drain still making progress never sees a deadline error, so
		// `graced` never becomes true on the very drain the halt exists to
		// end — gating this on `graced` would leave the trickler running.
		// Observation is gated on the direct child having been reaped. The
		// halt is raised on ANY read error once a drain has overflowed —
		// INCLUDING EOF — so a child that overflows stdout and then closes
		// it raises it while still live. Honouring it there ends this drain
		// on a pipe that live child is still writing: it blocks in
		// `write(2)`, `Wait` never returns, and the invocation spends its
		// whole deadline and reports a timeout instead of the overflow
		// refusal it owes. That is the serialization `0026:C1` `whole:`
		// forbids ("the drains run CONCURRENTLY and are never stalled") and
		// REQ-52 states as a negative requirement.
		if halted() {
			// A sibling already reported a refusing condition, so this
			// invocation refuses whatever this drain does next. The ranking
			// still holds: overflow outranks held on the same drain.
			if overflowed() {
				return kept, drainOverflow
			}
			return kept, drainHeld
		}
		if graced {
			_ = r.SetReadDeadline(time.Now().Add(DrainGrace * time.Millisecond))
		}
		n, err := r.Read(buf)
		if n > 0 {
			if room := limit - len(kept); room > 0 {
				if n < room {
					room = n
				}
				kept = append(kept, buf[:room]...)
			}
			// Past the limit the bytes are DISCARDED rather than left in
			// the pipe, so the child is never blocked writing into a full
			// one — and the same loop, with the same deadline, is what does
			// the discarding.
		}
		if err == nil {
			continue
		}
		// The report is ONE ordered value, so the ranking is applied here
		// and not re-derived by the caller: a drain that filled its bound
		// reports overflow WHATEVER its terminal condition, because the
		// bytes were read and the pipe's state is then not what the
		// invocation turns on.
		if overflowed() {
			if halt != nil {
				halt.Store(true)
			}
			return kept, drainOverflow
		}
		if errors.Is(err, os.ErrDeadlineExceeded) {
			if !graced {
				// The parent's mark, not a verdict: the bound expired and
				// the FINAL read under this drain's own grace is owed.
				graced = true
				continue
			}
			// A refusing condition. Raise `stop` so a sibling that is still
			// making progress into a value that will never be parsed ends
			// too, and the join the parent must complete returns.
			if halt != nil {
				halt.Store(true)
			}
			return kept, drainHeld
		}
		// EOF, or a read end closed out from under the drain: either way no
		// writer is going to deliver more, which is whole output as today.
		return kept, drainWhole
	}
}

// tail keeps the LAST StderrTailCap bytes, which is the diagnosis a caller
// wants when a tool is verbose before it fails.
func tail(b []byte) string {
	if len(b) > StderrTailCap {
		b = b[len(b)-StderrTailCap:]
	}
	return string(b)
}

// refuse builds a pre-spawn refusal: a typed error carrying its own reason
// as the Detail, since no child produced a stderr tail.
func refuse(detail string) error {
	return &accessor.ExecError{Detail: detail}
}

// refuseRequest builds a pre-spawn refusal that is about the REQUEST
// rather than the environment, so it takes the exit-2 group instead of
// `execution_failure`'s default exit 3 (RDR 0028 `0028:C1.3` EXIT
// GROUP:, which governs the entry-level preconditions `order:` names —
// among them C1.6's unbound `{tag.<key>}` and its `-`-prefixed bound
// value).
//
// It also builds `0025:C6`'s gate refusal (see `gateRefusal`), a missing
// opt-in being a property of the request on the same reasoning.
//
// It is deliberately a SIBLING of `refuse` rather than a widening of it:
// 0025's other pre-spawn refusals — the argv0 deny-list, a non-absolute
// `{artifact}` path — keep their existing routing, which this record
// does not amend.
func refuseRequest(detail string) error {
	return &accessor.ExecError{Detail: detail, Err: accessor.ErrDeclaredRequest}
}

// wrap builds an invocation refusal, carrying the stderr tail as the
// Detail and WRAPPING the offending error rather than flattening it, so
// `errors.Is(err, exec.ErrNotFound)` survives to the refusal site
// (`0025:C4`).
func wrap(detail string, err error) error {
	if detail == "" {
		// A refusal with no tail still needs a diagnosable Detail; the
		// error's own text is the only thing there is. It is not the
		// stderr TAIL, which is what the clause distinguishes — there is
		// none.
		return &accessor.ExecError{Err: err}
	}
	return &accessor.ExecError{Detail: detail, Err: err}
}

// --- C2: the authority bound ---------------------------------------------

// substitute performs the ONE rewriting the executed argv admits:
// whole-element replacement of the closed placeholder vocabulary with the
// caller-bound artifact path. Every other element crosses byte-for-byte,
// with no shell, no word splitting, and no glob expansion (`0025:C2`).
func substitute(declared []string, art accessor.Artifact) ([]string, error) {
	out := make([]string, 0, len(declared))
	substituted := false
	for _, el := range declared {
		if el != ArtifactPlaceholder {
			// RDR 0028 `0028:C1.6` — the `{tag.<key>}` family joins this
			// vocabulary, WHOLE-ELEMENT under the same substitution rule,
			// replaced by the bound value of a tag key the model
			// declares. Its two refusals are rules on THIS FAMILY and are
			// not entries on 0027's interpreter deny-list, which stays
			// static.
			//
			// The scan is per-USE-SITE: only an element that NAMES a tag
			// is looked up, so a flag-shaped value bound on the context
			// but referenced by no argv element of this entry reaches no
			// child and is never scanned.
			if key, isTag := table.CommandTagKey(el); isTag {
				value, bound := art.Context[key]
				if !bound {
					// A placeholder is never passed through literally: an
					// unbound `{tag.x}` forwarded as itself would reach
					// the child as an argument and could be read as a
					// filename (`0025:C2`).
					return nil, refuseRequest("the placeholder " + el +
						" is not bound on this invocation's context; a " +
						"placeholder is never passed through literally")
				}
				if strings.HasPrefix(value, "-") {
					// The class is argument injection (CWE-88), not shell
					// injection — `os/exec` runs no shell — so the hazard
					// is a flag-shaped WORD the child reads as a flag.
					// The rule mirrors `{artifact}`'s existing one rather
					// than inventing a second policy, and no `--`
					// separator is inserted: only a child that honours it
					// would be helped and the model cannot know which do.
					return nil, refuseRequest("the placeholder " + el +
						" is bound to " + strconv.Quote(value) +
						", which a child can parse as a flag")
				}
				out = append(out, value)
				continue
			}
			out = append(out, el)
			continue
		}
		if !substituted {
			// The substituted value must be an absolute path, and the
			// binding REFUSES rather than absolutizes: `--artifact` paths
			// are stored verbatim and intrastate's cwd is not the tool's
			// frame of reference, so silently resolving against the wrong
			// base is the failure this refusal exists to prevent
			// (`0025:C2`, premortem P-3).
			if !filepath.IsAbs(art.Path) || strings.HasPrefix(art.Path, "-") {
				return nil, refuse("artifact path must be absolute for a " +
					"command entry; the role `" + art.Role + "` is bound to " +
					strconv.Quote(art.Path) +
					", which a tool can parse as a flag or resolve against " +
					"its own working directory")
			}
			substituted = true
		}
		out = append(out, art.Path)
	}
	return out, nil
}

// resolveArgv0 fixes what the spawn executes, and never restores implicit
// current-directory lookup (`0025:C2`; Go's own `exec.ErrDot` reversal).
//
// A bare name resolves through the PARENT's `PATH` at spawn; a name
// carrying a path separator resolves against the model file's directory;
// an absolute name is already resolved. A separator-bearing argv0 in a
// model with no source file is refused at INVOCATION, before any spawn —
// what it must never do is fall back to the process cwd.
func resolveArgv0(argv0, baseDir string) (string, error) {
	switch {
	case filepath.IsAbs(argv0):
		return argv0, nil
	case !strings.ContainsRune(argv0, filepath.Separator):
		// A bare name. `exec.LookPath` reads the PARENT's PATH, which the
		// C4 allowlist then passes on to the child unchanged (A2).
		resolved, err := exec.LookPath(argv0)
		if err != nil {
			return "", wrap("", err)
		}
		return resolved, nil
	case baseDir == "":
		return "", refuse("the argv0 " + strconv.Quote(argv0) +
			" carries a path separator and resolves against the model " +
			"file's directory, but this model has no source file; the " +
			"process working directory is never the fallback")
	default:
		return filepath.Join(baseDir, argv0), nil
	}
}

// --- C4: the child environment -------------------------------------------

// childEnv composes the child's environment rather than inheriting it. The
// layers are applied in the order C4 fixes, so on a key collision the LATER
// layer wins: parent allowlist, then `env_pass`, then the entry's literal
// `env`, then the `INTRASTATE_*` overlay (`0025:C4`, A7).
//
// A variable that is unset in the parent is simply not passed: absence is
// not an empty value, and synthesizing one would tell the tool something
// the parent never said.
func childEnv(acc table.Accessor, name string, capability accessor.Capability) []string {
	env := map[string]string{}

	// 1 — the parent allowlist, plus the ONE prefix rule.
	for _, kv := range os.Environ() {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		if slices.Contains(AllowlistedVars, k) || strings.HasPrefix(k, AllowlistedPrefix) {
			env[k] = v
		}
	}

	// 2 — `env_pass`: the named, no-glob escape hatch.
	for _, k := range acc.EnvPass {
		if v, held := os.LookupEnv(k); held {
			env[k] = v
		}
	}

	// 3 — the entry's literal `env`.
	maps.Copy(env, acc.Env)

	// 4 — the overlay, which gives wrappers their context without new
	// placeholders. C5's `command_env_conflict` makes it unshadowable.
	env[EnvRole] = acc.Role
	env[EnvCapability] = string(capability)
	env[EnvAccessor] = name
	env[EnvProtocol] = ProtocolVersion

	out := make([]string, 0, len(env))
	for _, k := range slices.Sorted(mapKeys(env)) {
		out = append(out, k+"="+env[k])
	}
	return out
}

func mapKeys(m map[string]string) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// --- C3: the stdin envelope ----------------------------------------------

// stdinObject encodes the flat JSON object of strings every invocation
// sends. A read or gate sends the EMPTY object — the object is always sent,
// and the requested key set is not on it, being already declared in the
// model (`0025:C3`).
func stdinObject(planned []resolve.Tag) []byte {
	obj := make(map[string]string, len(planned))
	for _, t := range planned {
		// A planned `<clear>` crosses UNCHANGED as the literal reserved
		// value: the tool or wrapper performs the removal, and read-back
		// verifies absence (`0025:C3`, `0004:C11`).
		obj[t.Key] = t.Value
	}
	// HTML escaping is DISABLED, so `<`, `>`, and `&` serialize as
	// themselves: a planned `<clear>` must reach the tool as the literal
	// reserved value, and the same non-escaping encoder is what every
	// other wire site in this repo uses (JDR 0001 §D13).
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(obj); err != nil {
		// A map[string]string cannot fail to marshal.
		return []byte("{}")
	}
	// Encode appends a newline; the envelope carries none.
	return bytes.TrimRight(buf.Bytes(), "\n")
}

// --- the read binding -----------------------------------------------------

// Reader is a command-backed `accessor.ReadBinding`.
type Reader struct {
	Accessor table.Accessor
	Name     string
	Config   Config
}

// Capability reports the capability this binding serves.
func (Reader) Capability() accessor.Capability { return accessor.CapRead }

// Read invokes the declared command and maps its envelope onto the seam's
// split return: a parsed key becomes a `KeyValue` in `values`, an omitted
// declared key a name in `unreadable`, and an `exit_absent` hit a present
// record carrying `Absent: true` — never omission from both slices, which
// the executor reads as unreadable (`0025:C3`, `0025:C7`).
func (r Reader) Read(ctx context.Context, art accessor.Artifact, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	// The selectors are compiled BEFORE the spawn: a selector set that
	// cannot be parsed is a defect of the entry, and running the tool
	// first would spend an invocation on output nothing could read. A
	// loaded model never reaches this refusal — lint proves every pointer
	// and the key bijection — so it guards the in-memory construction
	// only (kata wchf).
	sel, err := compileSelectors(r.Accessor.Select, requested)
	if err != nil {
		return nil, nil, refuse(r.Name + ": " + err.Error())
	}

	inv, err := spawn(ctx, r.Accessor, r.Name, accessor.CapRead, r.Config, art,
		stdinObject(nil))
	if err != nil {
		return nil, nil, err
	}

	// Ordering is normative: a NON-EMPTY stdout is parsed first, and the
	// exit maps apply only to an empty stdout (`0025:C4`).
	if len(inv.stdout) != 0 {
		return r.parse(inv, requested, sel)
	}

	// An empty stdout carries meaning ONLY through `exit_absent`. A silent
	// tool is a broken tool, not an answer — exit 0 included.
	if inv.exited && slices.Contains(r.Accessor.ExitAbsent, inv.exitCode) {
		values := make([]accessor.KeyValue, 0, len(requested))
		for _, k := range requested {
			values = append(values, accessor.KeyValue{Key: k, Absent: true})
		}
		return values, nil, nil
	}
	if inv.signaled {
		return nil, nil, wrap(inv.detail(""), errors.New(
			"the read command produced no stdout and was killed by a signal "+
				"before it exited, so its `exit_absent` map does not apply"))
	}
	return nil, nil, wrap(inv.detail(""), errors.New(
		"the read command produced no stdout and exited "+
			strconv.Itoa(inv.exitCode)+
			", which the entry's `exit_absent` does not list"))
}

// parse maps a non-empty stdout onto the split return, per the entry's
// declared output mode.
func (r Reader) parse(inv invocation, requested []string, sel map[string]selector) (
	[]accessor.KeyValue, []string, error,
) {
	if r.Accessor.Output != nil && *r.Accessor.Output == "raw" {
		// The single declared key's value is stdout minus EXACTLY ONE
		// trailing "\n" — no trimming, no case folding, no whitespace
		// normalization (`0025:C3`, FX-raw-read).
		if len(requested) != 1 {
			return nil, nil, wrap(inv.detail(""), errors.New(
				"raw mode carries one declared key and this read requested "+
					strconv.Itoa(len(requested))))
		}
		value := string(inv.stdout)
		value = strings.TrimSuffix(value, "\n")
		return []accessor.KeyValue{{Key: requested[0], Value: value}}, nil, nil
	}

	// An entry declaring `select` reads a json DOCUMENT through its
	// per-key pointers; one declaring none reads C3's flat object below,
	// byte-for-byte as before (kata wchf).
	if sel != nil {
		return parseSelect(inv, requested, sel)
	}

	var obj map[string]string
	if jerr := json.Unmarshal(inv.stdout, &obj); jerr != nil {
		// Name the SHAPE, not just the decoder's complaint. The remedy is
		// to PROJECT the tool's output onto `0025:C3`'s wire shape, and
		// `encoding/json`'s own text ("cannot unmarshal array into Go
		// value of type string") names a Go type the author never wrote.
		// The clause id stays in this comment: the author's remedy is the
		// shape, and a record id is provenance they cannot act on.
		return nil, nil, wrap(inv.detail(""), errors.New(
			"the read command's stdout is not a flat JSON object of "+
				"string values: "+jerr.Error()))
	}

	var values []accessor.KeyValue
	var unreadable []string
	for _, k := range requested {
		v, held := obj[k]
		if !held {
			// A tool that fails to report a key has not established that
			// the key has no value: omission is UNREADABLE, never
			// established-absent (`0025:C3`).
			unreadable = append(unreadable, k)
			continue
		}
		values = append(values, accessor.KeyValue{Key: k, Value: v})
	}
	return values, unreadable, nil
}

// --- the json selectors ---------------------------------------------------

// selector is one key's compiled `select` table (kata wchf): the parsed
// RFC 6901 tokens and the declared projection and absence shapes.
type selector struct {
	pointer       []string
	element       []string
	prefix        string
	absentNull    bool
	absentMissing bool
}

// compileSelectors parses the entry's selectors for the requested keys. It
// returns nil for an entry declaring none, which is the flat-object read.
func compileSelectors(rules map[string]table.SelectRule, requested []string) (
	map[string]selector, error,
) {
	if len(rules) == 0 {
		return nil, nil
	}
	out := make(map[string]selector, len(requested))
	for _, k := range requested {
		rule, ok := rules[k]
		if !ok {
			return nil, errors.New("the entry declares `select` tables but none " +
				"for the key " + k + "; every declared key needs one")
		}
		pointer, err := table.ParseJSONPointer(rule.Pointer)
		if err != nil {
			return nil, errors.New("`select." + k + "`: " + err.Error())
		}
		if rule.Element != "" && rule.Prefix == "" {
			return nil, errors.New("`select." + k + "` declares `element` " +
				"without `prefix`")
		}
		element, err := table.ParseJSONPointer(rule.Element)
		if err != nil {
			return nil, errors.New("`select." + k + "` element " + err.Error())
		}
		out[k] = selector{
			pointer:       pointer,
			element:       element,
			prefix:        rule.Prefix,
			absentNull:    rule.AbsentNull,
			absentMissing: rule.AbsentMissing,
		}
	}
	return out, nil
}

// parseSelect maps a json stdout onto the split return through the
// entry's selectors (kata wchf).
//
// It NARROWS `0025:C3` rather than lifting it: whatever a selector cannot
// establish is UNREADABLE, never a guess. Absence is established only by
// a shape the author DECLARED to mean it, or by a prefix projection over
// a set the tool reported whole in which no member carries the prefix.
func parseSelect(inv invocation, requested []string, sel map[string]selector) (
	[]accessor.KeyValue, []string, error,
) {
	doc, err := decodeDocument(inv.stdout)
	if err != nil {
		return nil, nil, wrap(inv.detail(""), errors.New(
			"the read command's stdout is not one JSON document: "+err.Error()))
	}

	var values []accessor.KeyValue
	var unreadable []string
	for _, k := range requested {
		value, out := sel[k].read(doc)
		switch out {
		case selValue:
			values = append(values, accessor.KeyValue{Key: k, Value: value})
		case selAbsent:
			values = append(values, accessor.KeyValue{Key: k, Absent: true})
		default:
			unreadable = append(unreadable, k)
		}
	}
	return values, unreadable, nil
}

// decodeDocument decodes exactly ONE json value, with numbers kept as
// their own text so `7` reads as "7" and never as "7e+00".
func decodeDocument(b []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var doc any
	if err := dec.Decode(&doc); err != nil {
		return nil, err
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("data follows the first JSON value")
	}
	return doc, nil
}

// selOutcome is what one selector established.
type selOutcome int

const (
	selUnreadable selOutcome = iota
	selValue
	selAbsent
)

// read applies one selector to the decoded document.
func (s selector) read(doc any) (string, selOutcome) {
	v, found := resolvePointer(doc, s.pointer)
	switch found {
	case lookupDrift:
		return "", selUnreadable
	case lookupMissing:
		if s.absentMissing {
			return "", selAbsent
		}
		return "", selUnreadable
	}
	if v == nil {
		if s.absentNull {
			return "", selAbsent
		}
		return "", selUnreadable
	}

	if s.prefix != "" {
		return s.project(v)
	}
	return scalarText(v)
}

// project reads the ONE array member that starts with the prefix,
// stripped of it. Zero matches is established-absent: the tool reported
// the whole set and no member carries the prefix. Two or more is
// UNREADABLE — "first match" would depend on array order and hide the
// corruption a second member is (`0007` totality).
//
// A member the element pointer cannot reach as a string is shape drift
// and makes the KEY unreadable rather than being skipped: a skipped
// member may be exactly the one that carried the prefix.
func (s selector) project(v any) (string, selOutcome) {
	members, ok := v.([]any)
	if !ok {
		return "", selUnreadable
	}
	var matched []string
	for _, m := range members {
		ev, found := resolvePointer(m, s.element)
		str, isString := ev.(string)
		if found != lookupFound || !isString {
			return "", selUnreadable
		}
		if rest, has := strings.CutPrefix(str, s.prefix); has {
			matched = append(matched, rest)
		}
	}
	switch len(matched) {
	case 0:
		return "", selAbsent
	case 1:
		return matched[0], selValue
	default:
		return "", selUnreadable
	}
}

// scalarText renders a json scalar as the key's value: a string as
// itself, a number and a bool as their json text. An object or an array
// is not one value and is UNREADABLE.
func scalarText(v any) (string, selOutcome) {
	switch t := v.(type) {
	case string:
		return t, selValue
	case json.Number:
		return t.String(), selValue
	case bool:
		return strconv.FormatBool(t), selValue
	default:
		return "", selUnreadable
	}
}

// lookup is what a pointer walk found.
type lookup int

const (
	lookupFound lookup = iota
	// lookupMissing: the FINAL token names nothing under a parent that
	// exists — the one shape `absent = ["missing"]` may establish.
	lookupMissing
	// lookupDrift: an intermediate node is missing or is not a container
	// the next token can index. That is the tool's shape moving, not a
	// value being unset, and it stays unreadable whatever is declared.
	lookupDrift
)

// resolvePointer walks parsed RFC 6901 tokens through a decoded document.
// An object member is matched byte-exactly; an array index is `0` or a
// digit string without a leading zero, and `-` names the member past the
// end, so it is always missing (RFC 6901 §4).
func resolvePointer(doc any, tokens []string) (any, lookup) {
	cur := doc
	for i, tok := range tokens {
		last := i == len(tokens)-1
		switch node := cur.(type) {
		case map[string]any:
			v, ok := node[tok]
			if !ok {
				return nil, missingAt(last)
			}
			cur = v
		case []any:
			idx, ok := arrayIndex(tok, len(node))
			if !ok {
				return nil, lookupDrift
			}
			if idx >= len(node) {
				return nil, missingAt(last)
			}
			cur = node[idx]
		default:
			return nil, lookupDrift
		}
	}
	return cur, lookupFound
}

func missingAt(last bool) lookup {
	if last {
		return lookupMissing
	}
	return lookupDrift
}

// arrayIndex parses an RFC 6901 array index. An index too large to
// represent is past the end of any array, so it is reported as `n`.
func arrayIndex(tok string, n int) (int, bool) {
	if tok == "-" {
		return n, true
	}
	if tok == "" || (len(tok) > 1 && tok[0] == '0') {
		return 0, false
	}
	for _, c := range tok {
		if c < '0' || c > '9' {
			return 0, false
		}
	}
	idx, err := strconv.Atoi(tok)
	if err != nil {
		return n, true
	}
	return idx, true
}

// --- the gate binding -----------------------------------------------------

// Gate is a command-backed `accessor.GateBinding`.
type Gate struct {
	Accessor table.Accessor
	Name     string
	Config   Config
}

// Capability reports the capability this binding serves.
func (Gate) Capability() accessor.Capability { return accessor.CapGate }

// Gate invokes the declared command and reads its verdict envelope or its
// declared exit map. Execution failure is never laundered into a verdict:
// an unlisted exit, a spawn failure, or a malformed stdout refuses
// (`0025:C3`).
func (g Gate) Gate(ctx context.Context, art accessor.Artifact) (
	accessor.Verdict, string, error,
) {
	inv, err := spawn(ctx, g.Accessor, g.Name, accessor.CapGate, g.Config, art,
		stdinObject(nil))
	if err != nil {
		return "", "", err
	}

	// Parse first: a well-formed deny envelope with a non-zero exit is a
	// DENY, not a failure (`0025:C4` ordering, premortem P-12).
	if len(inv.stdout) != 0 {
		var env struct {
			Verdict string `json:"verdict"`
			Reason  string `json:"reason"`
		}
		if jerr := json.Unmarshal(inv.stdout, &env); jerr != nil {
			// Same reason as the reader's: the shape is the remedy, and the
			// decoder's Go-type text does not state it.
			return "", "", wrap(inv.detail(""), errors.New(
				"the gate command's stdout is not a JSON object carrying "+
					"the string fields `verdict` and `reason`: "+jerr.Error()))
		}
		if !slices.Contains(accessor.Verdicts(), accessor.Verdict(env.Verdict)) {
			return "", "", wrap(inv.detail(""), errors.New(
				"the gate envelope names the verdict "+strconv.Quote(env.Verdict)+
					", which is not one of allow, deny, indeterminate"))
		}
		return accessor.Verdict(env.Verdict), env.Reason, nil
	}

	// An empty stdout is that verdict only when the entry's map LISTS the
	// exit. The map is consulted only for a process that ran and exited, so
	// a spawn failure — which has no exit code at all — can never be
	// claimed by it however broadly it is written.
	if inv.exited {
		if v, listed := g.Accessor.ExitVerdicts[strconv.Itoa(inv.exitCode)]; listed {
			return accessor.Verdict(v), "", nil
		}
	}
	if inv.signaled {
		return "", "", wrap(inv.detail(""), errors.New(
			"the gate command produced no stdout and was killed by a signal "+
				"before it exited, so its `exit_verdicts` map does not apply"))
	}
	return "", "", wrap(inv.detail(""), errors.New(
		"the gate command produced no stdout and exited "+
			strconv.Itoa(inv.exitCode)+
			", which the entry's `exit_verdicts` does not list"))
}

// --- the write binding ----------------------------------------------------

// Writer is a command-backed `accessor.WriteBinding`.
type Writer struct {
	Accessor table.Accessor
	Name     string
	Config   Config

	invocations int
}

// Capability reports the capability this binding serves.
func (*Writer) Capability() accessor.Capability { return accessor.CapWrite }

// Apply sends the planned tags to the declared command on stdin. Its
// success is never taken from exit status: verification is read-back
// (`0025:C3`, `0004:C12`). One Apply is one spawn and one increment — the
// binding performs no internal respawn (`0025:C7`).
func (w *Writer) Apply(ctx context.Context, art accessor.Artifact, planned []resolve.Tag) error {
	w.invocations++

	inv, err := spawn(ctx, w.Accessor, w.Name, accessor.CapWrite, w.Config, art,
		stdinObject(planned))
	if err != nil {
		return err
	}
	if inv.signaled {
		return wrap(inv.detail(""), errors.New(
			"the write command was killed by a signal before it exited"))
	}
	if inv.exitCode != 0 {
		return wrap(inv.detail(""), errors.New(
			"the write command exited "+strconv.Itoa(inv.exitCode)))
	}
	return nil
}

// Invocations counts Apply ENTRIES — one per call, no internal respawn
// (`0025:C7`). The count is about the accessor LAYER's re-entry, not about
// process success, so a failing spawn still counts one.
func (w *Writer) Invocations() int { return w.invocations }

// --- the steps write binding ----------------------------------------------

// StepsWriter is the `steps`-carried `accessor.WriteBinding` (kata q14r):
// per key, per planned VALUE, a list of LITERAL argv vectors run in
// order. It is how a model drives an argv-only tool — one that ignores
// stdin and changes state through several calls, such as a label remove
// then a label add — with no wrapper script.
//
// No vector carries a planned value. The argv vocabulary is `command`'s
// own, `{artifact}` and the observed `{tag.<key>}` family under the same
// whole-element substitution (`0025:C2`, `0028:C1.6`), and the planned
// object still crosses on stdin (`0025:C3`). What a value changes is
// WHICH literal vectors run, so every executable argv is readable off the
// model.
//
// Each step is one `spawn`, under the same gate, environment, process
// group and bounds as a `command` write. Success is never taken from exit
// status: verification is the executor's one read-back after the last
// step (`0004:C12`).
type StepsWriter struct {
	Accessor table.Accessor
	Name     string
	Config   Config

	invocations int
}

// The executor reaches the prior pre-read only through this capability.
var _ accessor.PriorBinding = (*StepsWriter)(nil)

// step is one selected argv vector and the arm position a refusal names
// it by, `steps.<key>.<arm>.<value>[<i>]`.
type step struct {
	label string
	argv  []string
}

// Capability reports the capability this binding serves.
func (*StepsWriter) Capability() accessor.Capability { return accessor.CapWrite }

// PriorKeys names the planned keys whose table declares any `clear` arm,
// in the entry's declared `keys` order. Only those select argv from the
// prior value; a tool whose `set` natively replaces declares no clear
// arms, and its write runs no pre-read.
func (w *StepsWriter) PriorKeys(planned []resolve.Tag) []string {
	var out []string
	for _, key := range w.Accessor.Keys {
		if len(w.Accessor.Steps[key].Clear) == 0 || slices.Contains(out, key) {
			continue
		}
		if slices.ContainsFunc(planned, func(t resolve.Tag) bool { return t.Key == key }) {
			out = append(out, key)
		}
	}
	return out
}

// Apply runs the selected steps with no prior value. A write whose arms
// need one refuses before any spawn, so a caller that skipped the
// executor's pre-read can never run `set` without the `clear` it owed.
func (w *StepsWriter) Apply(ctx context.Context, art accessor.Artifact, planned []resolve.Tag) error {
	return w.ApplyPrior(ctx, art, planned, nil)
}

// ApplyPrior selects and runs the steps for the planned tags against the
// pre-read prior values (kata q14r).
//
// Selection, per planned key in the entry's declared `keys` order:
//
//   - a value replacing a different held value runs `clear.<prior>` (when
//     that arm exists) then `set.<planned>`;
//   - a value equal to the held one runs `set.<planned>` alone — it
//     re-asserts, and a clear would remove what it then re-adds;
//   - a planned `<clear>` runs `clear.<prior>` alone, and nothing when
//     the key is already absent, which `0004:C11` requires to succeed.
//
// Everything that can refuse without running a child is decided before
// the FIRST spawn: the gate, the selection, and every selected vector's
// placeholder substitution and argv0 resolution. Then the steps run in
// order and the first non-zero exit or signal stops the write. Fail-fast
// rather than run-all: a lost compare-and-set (a claim exiting non-zero)
// must not be followed by the label changes it guarded.
//
// A first-step failure applied nothing and keeps the safe-to-retry sense.
// A later failure is a `*accessor.PartialApplyError`, which the executor
// reports as applied-but-unverified; nothing is retried or undone
// (`0004:C14`).
func (w *StepsWriter) ApplyPrior(ctx context.Context, art accessor.Artifact,
	planned []resolve.Tag, prior []accessor.KeyValue,
) error {
	w.invocations++

	// C6 — checked here and not only in `spawn`, because a selection can
	// run nothing (a `<clear>` of an absent key), and a `steps` entry
	// without the opt-in must refuse whatever it selects.
	if !w.Config.AllowCommands {
		return gateRefusal(w.Name)
	}

	steps, err := w.selectSteps(planned, prior)
	if err != nil {
		return err
	}
	for _, s := range steps {
		if len(s.argv) == 0 {
			return refuse("the accessor `" + w.Name + "` declares the empty step " +
				s.label + ", which cannot be invoked")
		}
		argv, serr := substitute(s.argv, art)
		if serr != nil {
			return serr
		}
		if _, rerr := resolveArgv0(argv[0], w.Config.BaseDir); rerr != nil {
			return rerr
		}
	}

	stdin := stdinObject(planned)
	for i, s := range steps {
		serr := w.runStep(ctx, art, s, stdin)
		if serr == nil {
			continue
		}
		if i == 0 {
			return serr
		}
		return partialApply(i, s.label, serr)
	}
	return nil
}

// selectSteps maps the plan onto the declared arms. Every refusal here
// precedes any spawn.
//
// The arm-coverage refusals are about the REQUEST and take
// `refuseRequest`'s exit-2 group, as `edit_clear_undeclared` does
// (`0028:C1.3` EXIT GROUP:): lint proves a `set` arm per domain value and
// a `clear` arm per value wherever a RULE can clear, but `flow set-state
// --clear` plans a removal no rule declared, and an artifact can hold a
// value outside the domain. Re-running either unchanged cannot help.
func (w *StepsWriter) selectSteps(planned []resolve.Tag, prior []accessor.KeyValue) ([]step, error) {
	want := make(map[string]string, len(planned))
	for _, t := range planned {
		if _, ok := w.Accessor.Steps[t.Key]; !ok {
			return nil, refuse("the accessor `" + w.Name + "` declares no `steps." +
				t.Key + "` table for the planned key " + t.Key)
		}
		want[t.Key] = t.Value
	}
	held := make(map[string]accessor.KeyValue, len(prior))
	for _, kv := range prior {
		held[kv.Key] = kv
	}

	var out []step
	add := func(key, arm, member string, vectors [][]string) {
		for i, argv := range vectors {
			out = append(out, step{
				label: "`steps." + key + "." + arm + "." + member + "[" + strconv.Itoa(i) + "]`",
				argv:  argv,
			})
		}
	}
	for _, key := range w.Accessor.Keys {
		value, ok := want[key]
		if !ok {
			continue
		}
		rule := w.Accessor.Steps[key]

		var before accessor.KeyValue
		if len(rule.Clear) != 0 {
			kv, read := held[key]
			if !read {
				return nil, refuse("the accessor `" + w.Name + "` selects the " +
					"steps for " + key + " from its prior value, and none was read")
			}
			before = kv
		}

		if accessor.IsClear(value) {
			if len(rule.Clear) == 0 {
				return nil, refuseRequest("the accessor `" + w.Name + "` declares no " +
					"`clear` arms for " + key + ", so it cannot remove it")
			}
			if before.Absent {
				continue
			}
			vectors, ok := rule.Clear[before.Value]
			if !ok {
				return nil, refuseRequest("the accessor `" + w.Name + "` cannot remove " +
					key + ": it holds " + strconv.Quote(before.Value) +
					", which has no `clear` arm")
			}
			add(key, "clear", before.Value, vectors)
			continue
		}

		set, ok := rule.Set[value]
		if !ok {
			return nil, refuseRequest("the accessor `" + w.Name + "` declares no `set." +
				value + "` arm for the planned " + key + " = " + strconv.Quote(value))
		}
		if len(rule.Clear) != 0 && !before.Absent && before.Value != value {
			if vectors, ok := rule.Clear[before.Value]; ok {
				add(key, "clear", before.Value, vectors)
			}
		}
		add(key, "set", value, set)
	}
	return out, nil
}

// runStep is one bounded invocation of one selected vector, refused on a
// signal or a non-zero exit exactly as `Writer.Apply` refuses its one.
func (w *StepsWriter) runStep(ctx context.Context, art accessor.Artifact, s step, stdin []byte) error {
	acc := w.Accessor
	acc.Command = s.argv
	inv, err := spawn(ctx, acc, w.Name, accessor.CapWrite, w.Config, art, stdin)
	if err != nil {
		return err
	}
	if inv.signaled {
		return wrap(inv.detail(""), errors.New(
			"the write step "+s.label+" was killed by a signal before it exited"))
	}
	if inv.exitCode != 0 {
		return wrap(inv.detail(""), errors.New(
			"the write step "+s.label+" exited "+strconv.Itoa(inv.exitCode)))
	}
	return nil
}

// partialApply wraps a later step's failure in the applied sense. The
// reason leads the Detail and the failing step's own diagnosis follows,
// as `invocation.detail` orders a held-pipe reason ahead of its tail.
func partialApply(ran int, label string, err error) error {
	detail := "the step " + label + " failed after " + strconv.Itoa(ran) +
		" earlier step(s) of this write exited 0, so the write is partly " +
		"applied; no step is retried or undone"
	if ee, ok := errors.AsType[*accessor.ExecError](err); ok && ee.Detail != "" {
		detail += "\n" + ee.Detail
	} else {
		detail += ": " + err.Error()
	}
	return &accessor.ExecError{
		Detail: detail,
		Err:    &accessor.PartialApplyError{Ran: ran, Step: label, Err: err},
	}
}

// Invocations counts Apply ENTRIES — one per Apply or ApplyPrior call,
// however many steps it runs. The count is about the accessor LAYER's
// re-entry, not process count (`0025:C7`), so a three-step write is one.
func (w *StepsWriter) Invocations() int { return w.invocations }
