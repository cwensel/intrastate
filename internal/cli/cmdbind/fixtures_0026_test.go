package cmdbind_test

// RDR 0026 — FX-deadline-escape and the shared timing oracles.
//
// Nothing here mocks the unit under test. FX-deadline-escape is a REAL Go
// helper binary that spawns a REAL `SysProcAttr{Setsid: true}` grandchild
// inheriting the parent's stdout and stderr, because the whole failure
// shape is a writer the group signal cannot reach — a shell fixture using
// `setsid sh -c …` does NOT escape on macOS (Background), so the escape
// must be made with a real `setsid(2)` call.
//
// REQ-110: "A `Setsid: true` Go test helper exists in `internal/cli/cmdbind`
// so the escape is real on both OSes (Background: a shell fixture does not
// escape on macOS)"
// REQ-133 (S-coverage): "the FX-deadline-escape helper is exercised on
// darwin and linux in CI, so the escape is real on both (Background)."
// REQ-113 (IP Phase 3): "Add FX-deadline-escape (both shapes) beside the
// existing grandchild-leak and cancel oracles"
//
// The helper's behaviour is selected by argv, so one binary covers every
// shape the record names: the deadline shape (child silent past the
// deadline, grandchild holds both pipes), the success shape (child writes a
// whole envelope and exits 0 while the grandchild holds), the stderr-only
// shape (S10), the paced and burst late-drain shapes (MVV row 6), the
// trickle shape (S9's companion) and the SIGTERM-resisting shape (S8).

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

const (
	// escTimeout is the MVV's declared timeout: `timeout = "2s"`.
	escTimeout = "2s"
	// escRole / escKey / escName name the fixture accessor.
	escRole = "state"
	escKey  = "state.phase"
	escName = "escaping.state"
)

// Bound / grace figures the contract publishes. They are re-stated here as
// DURATIONS derived from the package's own exported constant so a test
// asserting the C1 `bound:` total moves with the constant rather than
// against a transcribed literal.
//
// REQ-41 (C1 `bound:`): "WaitDelay (500 ms) is the ONE bound"
// REQ-44 (C1 `bound:`): "the COMMITTED return bound is timeout +
// 2·WaitDelay, asserted with a scheduling tolerance of +100 ms"
var (
	escWaitDelay = time.Duration(cmdbind.WaitDelay) * time.Millisecond
	escTolerance = 100 * time.Millisecond
)

// escBoundFor returns the C1 `bound:` committed total for a declared
// timeout, with the record's stated scheduling tolerance already added.
func escBoundFor(timeout time.Duration) time.Duration {
	return timeout + 2*escWaitDelay + escTolerance
}

// --- the FX-deadline-escape helper ---------------------------------------

var (
	escOnce sync.Once
	escPath string
	escErr  error
)

// escapeBin builds (once) FX-deadline-escape and returns its absolute path.
//
// It is a real Go binary rather than a shell script for the reason the
// Background records: macOS ships no `setsid` BINARY, so only a real
// `syscall.Setsid` — which is what `SysProcAttr{Setsid: true}` performs —
// escapes the process group on both supported OSes.
func escapeBin(t *testing.T) string {
	t.Helper()

	escOnce.Do(func() {
		dir, err := os.MkdirTemp("", "cmdbind-escape")
		if err != nil {
			escErr = err
			return
		}
		src := filepath.Join(dir, "main.go")
		if werr := os.WriteFile(src, []byte(escapeSource), 0o600); werr != nil {
			escErr = werr
			return
		}
		out := filepath.Join(dir, "escapehelper")
		cmd := exec.Command("go", "build", "-o", out, src)
		cmd.Env = append(os.Environ(), "GO111MODULE=off")
		if b, berr := cmd.CombinedOutput(); berr != nil {
			escErr = errBuild{berr, string(b)}
			return
		}
		escPath = out
	})
	if escErr != nil {
		t.Fatalf("building FX-deadline-escape failed: %v", escErr)
	}
	return escPath
}

// escapeSource is FX-deadline-escape. Its first argument selects a shape;
// the second is always the pidfile the escaped grandchild's pid is recorded
// in, so a test can assert the admitted residue (C1 `residue:`) and kill it.
const escapeSource = `package main

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// escape spawns a grandchild in its OWN SESSION, so the parent's
// group-scoped kill(2) cannot reach it. Which of the parent's own stdio the
// grandchild inherits is the argument: that is what makes "stdout held",
// "stderr held" and "both held" three different fixtures. It is
// escapeHolding with the stdin read end NOT inherited, which is what every
// output-side fixture wants.
func escape(pidfile string, holdOut, holdErr bool, sleep string) {
	escapeHolding(pidfile, false, holdOut, holdErr, sleep)
}

// escapeHolding is escape with the STDIN read end also selectable. An
// escapee holding stdin is what makes the parent's stdin write outlive the
// direct child: the CLI's Wait then ends in exec.ErrWaitDelay rather
// than in an ExitError, which is the ONLY way to reach spawn's
// non-ExitError arm. A fixture that wants that arm must ALSO exit its direct
// child normally and immediately — a child that sleeps past the deadline is
// SIGKILLed by the group-scoped cancel and reports a signalled ExitError
// instead, which takes a different arm entirely.
func escapeHolding(pidfile string, holdIn, holdOut, holdErr bool, sleep string) {
	self, _ := os.Executable()
	sub := exec.Command(self, "hold", sleep)
	sub.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if holdIn {
		sub.Stdin = os.Stdin
	}
	if holdOut {
		sub.Stdout = os.Stdout
	}
	if holdErr {
		sub.Stderr = os.Stderr
	}
	if err := sub.Start(); err != nil {
		os.Exit(95)
	}
	if pidfile != "" {
		_ = os.WriteFile(pidfile, []byte(strconv.Itoa(sub.Process.Pid)), 0o600)
	}
}

// dripGap is the idle gap between the trickler's bytes, and it is the whole
// point of S9's companion: it must be LONGER than DrainGrace (50 ms), which
// is the gap the drain's re-armed deadline discriminates on. C1 precedence:
// names that case directly — "a writer trickling bytes SLOWER than the grace
// is exactly the case that must refuse" — and the same clause fixes the
// opposite outcome for a faster one: "a drain making progress keeps its
// deadline ahead of it and runs to EOF, and only a gap longer than the grace
// ends it as held."
//
// The value shipped originally was 10 ms — a FIFTH of the grace, so the
// stdout drain re-armed forever and never reported on its own. That row
// passed only because this fixture writes to STDOUT ONLY: its inherited
// stderr sat idle with a live writer, reported held under the grace alone,
// and raised the sibling halt that ended the stdout drain. The row was green
// for a reason its own text did not name. Repaired as deviation D8.
//
// 100 ms is 2x the grace: the deadline expires a full grace before the next
// byte could arrive, so the stdout drain reaches held on its OWN mechanism
// with 50 ms of margin against scheduling jitter (the measured tail under
// 4x-core saturation is ~28 ms). It also stays far inside the committed
// bound: the refusal is owed at most WaitDelay (500, the join timer) + one
// gap (100) + one grace (50) = ~650 ms against the committed 2*WaitDelay +
// 100 ms tolerance = 1100 ms, so the row is neither slow nor fragile.
const dripGap = 100 * time.Millisecond

func main() {
	if len(os.Args) < 2 {
		os.Exit(97)
	}
	mode := os.Args[1]
	rest := os.Args[2:]

	switch mode {
	// hold is the GRANDCHILD: it holds whatever stdio it inherited and
	// sleeps, writing nothing. Nothing in the CLI can reach it.
	case "hold":
		d, _ := time.ParseDuration(rest[0])
		time.Sleep(d)

	// deadline is the MVV's deadline shape: the grandchild inherits BOTH
	// pipes and sleeps 60 s; the child itself writes nothing and sleeps
	// past the declared deadline.
	case "deadline":
		escape(rest[0], true, true, "60s")
		time.Sleep(60 * time.Second)

	// success is the MVV's row-5 shape: the grandchild holds BOTH pipes
	// while the child writes a whole envelope and exits 0.
	case "success":
		escape(rest[0], true, true, "60s")
		_, _ = os.Stdout.WriteString(rest[1])
		os.Exit(0)

	// stderr-held is S10: stdout reaches EOF cleanly (the child closes it
	// and no escapee holds it) while STDERR ALONE is held by the escapee.
	// The child writes its stderr tail first, so the tail is buffered and
	// the grace is what preserves it.
	case "stderr-held":
		escape(rest[0], false, true, "60s")
		_, _ = os.Stderr.WriteString(rest[1])
		_, _ = os.Stdout.WriteString(rest[2])
		_ = os.Stdout.Close()
		os.Exit(0)

	// stdout-held is S10's named negative control: the same refusal with
	// the OTHER pipe named, and no tail.
	case "stdout-held":
		escape(rest[0], true, false, "60s")
		_ = os.Stderr.Close()
		os.Exit(0)

	// trickle is S9's COMPANION: an escaped writer emitting a byte on a
	// cadence SLOWER than DrainGrace, forever. It never blocks and never
	// reaches EOF, so a predicate keyed on emptiness or on byte count
	// accepts it and only the drain's reported terminal condition refuses
	// it. See dripGap for why the cadence is what it is (deviation D8).
	case "trickle":
		self, _ := os.Executable()
		sub := exec.Command(self, "drip")
		sub.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		sub.Stdout = os.Stdout
		sub.Stderr = os.Stderr
		if err := sub.Start(); err != nil {
			os.Exit(95)
		}
		_ = os.WriteFile(rest[0], []byte(strconv.Itoa(sub.Process.Pid)), 0o600)
		os.Exit(0)

	// trickle-stdout is S9's companion made DISCRIMINATING. It is trickle
	// with the escapee inheriting STDOUT ONLY and the direct child closing
	// its own stderr, mirroring stdout-held. Under trickle the escapee also
	// inherits an idle stderr, which reports held under the grace alone and
	// raises the sibling halt — so stdout was named held by the HALT rather
	// than by its own idle-gap deadline, and an oracle asserting only that
	// the invocation refused could not tell the two apart. With stderr
	// reaching EOF, the only mechanism that can name stdout is the stdout
	// drain's OWN report, which is the mechanism S9's companion exists to
	// witness (deviation D8's recommended amendment).
	//
	// trickle itself is unchanged: D10's both-pipes residue row and the
	// committed-total row legitimately need that shape.
	case "trickle-stdout":
		self, _ := os.Executable()
		sub := exec.Command(self, "drip")
		sub.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		sub.Stdout = os.Stdout
		if err := sub.Start(); err != nil {
			os.Exit(95)
		}
		_ = os.WriteFile(rest[0], []byte(strconv.Itoa(sub.Process.Pid)), 0o600)
		_ = os.Stderr.Close()
		os.Exit(0)

	case "drip":
		for {
			_, _ = os.Stdout.WriteString("x")
			time.Sleep(dripGap)
		}

	// burst is MVV row 6(a): N bytes on stdout in ONE Write, then exit,
	// with NO grandchild. Every writer is gone when the drain resumes.
	case "burst":
		n, _ := strconv.Atoi(rest[0])
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = 'x'
		}
		_, _ = os.Stdout.Write(buf)
		os.Exit(0)

	// paced is MVV row 6(b), the MECHANISM DISCRIMINATOR: the same payload
	// written in N chunks separated by idle gaps, then exit with no
	// grandchild. Only a deadline RE-ARMED before each read recovers it
	// whole; a single absolute now+DrainGrace truncates.
	case "paced":
		total, _ := strconv.Atoi(rest[0])
		chunks, _ := strconv.Atoi(rest[1])
		gap, _ := time.ParseDuration(rest[2])
		size := total / chunks
		buf := make([]byte, size)
		for i := range buf {
			buf[i] = 'x'
		}
		for i := 0; i < chunks; i++ {
			if i > 0 {
				time.Sleep(gap)
			}
			_, _ = os.Stdout.Write(buf)
		}
		os.Exit(0)

	// resists is S8: the child IGNORES SIGTERM and survives to its
	// WaitDelay while a grandchild holds both pipes. Leg (a) of the bound
	// spends its allowance instead of donating it to leg (b).
	case "resists":
		signal.Ignore(syscall.SIGTERM, syscall.SIGINT)
		escape(rest[0], true, true, "60s")
		time.Sleep(60 * time.Second)

	// overflow-and-hold makes BOTH conditions qualify on the same drain:
	// a stdout past StdoutCap AND an escapee holding the pipes.
	// C1 precedence: fixes that overflow outranks held.
	case "overflow-and-hold":
		escape(rest[0], true, true, "60s")
		n, _ := strconv.Atoi(rest[1])
		buf := make([]byte, 4096)
		for i := range buf {
			buf[i] = 'x'
		}
		for n > 0 {
			w := len(buf)
			if n < w {
				w = n
			}
			m, err := os.Stdout.Write(buf[:w])
			if err != nil {
				break
			}
			n -= m
		}
		os.Exit(0)

	// no-stdin-holds leaves an escapee holding BOTH pipes and then sleeps
	// past the declared deadline. It does NOT reach the werr non-ExitError
	// arm: the sleeping child is SIGKILLed by the group-scoped cancel, so
	// Wait reports a SIGNALLED ExitError. The dual ErrWaitDelay-vs-held
	// competition is stdin-holds-exits, below.
	case "no-stdin-holds":
		escape(rest[0], true, true, "60s")
		d, _ := time.ParseDuration(rest[1])
		time.Sleep(d)

	// stdin-holds-exits stages the REQ-79 competition: the escapee holds
	// stdin AND both output pipes, and the direct child EXITS NORMALLY AND
	// IMMEDIATELY without reading a byte. The immediate exit is what makes
	// this different from every sleeping mode — nothing signals the child,
	// so Wait cannot report an ExitError, and the still-open stdin read end
	// keeps the parent's copy alive until WaitDelay ends it as
	// exec.ErrWaitDelay. The held output pipes qualify at the same instant,
	// so the drain condition and the non-ExitError arm compete for real.
	case "stdin-holds-exits":
		escapeHolding(rest[0], true, true, true, "60s")
		os.Exit(0)

	// stdin-holds-only is the SINGLE-condition half of the same shape: the
	// escapee holds stdin ALONE, so the output drains reach EOF and no
	// drain condition qualifies, leaving exec.ErrWaitDelay as the only
	// reason to refuse (S4, REQ-40/98/125).
	case "stdin-holds-only":
		escapeHolding(rest[0], true, false, false, "60s")
		os.Exit(0)

	// cap-exact / cap-over are S3's boundary shapes, with the child
	// closing its pipes normally: exactly the cap is a VALUE, one byte
	// more is a detected overflow.
	case "emit-n":
		n, _ := strconv.Atoi(rest[0])
		buf := make([]byte, 4096)
		for i := range buf {
			buf[i] = 'x'
		}
		for n > 0 {
			w := len(buf)
			if n < w {
				w = n
			}
			m, err := os.Stdout.Write(buf[:w])
			if err != nil {
				break
			}
			n -= m
		}
		os.Exit(0)

	// overflow-then-stderr is the SIBLING-HALT shape, and it has NO
	// escapee: the DIRECT child overflows stdout, closes it, and only then
	// writes a stderr payload larger than the pipe buffer. The stdout drain
	// therefore sees EOF-after-overflow while this child is still live, so
	// a halt observed before the child is reaped strands the stderr drain
	// on a pipe a live writer is still filling. See REQ-52 and the
	// stdout-overflow oracle in drain_0026_test.go.
	case "overflow-then-stderr":
		n, _ := strconv.Atoi(rest[0])
		buf := make([]byte, 4096)
		for i := range buf {
			buf[i] = 'x'
		}
		for n > 0 {
			w := len(buf)
			if n < w {
				w = n
			}
			m, err := os.Stdout.Write(buf[:w])
			if err != nil {
				break
			}
			n -= m
		}
		_ = os.Stdout.Close()
		e, _ := strconv.Atoi(rest[1])
		ebuf := make([]byte, 4096)
		for i := range ebuf {
			ebuf[i] = 'e'
		}
		for e > 0 {
			w := len(ebuf)
			if e < w {
				w = e
			}
			m, err := os.Stderr.Write(ebuf[:w])
			if err != nil {
				break
			}
			e -= m
		}
		os.Exit(0)

	// whole writes an exact stdout payload and an exact stderr payload,
	// closes both and exits 0. It is S3's byte-for-byte ordinary path.
	case "whole":
		_, _ = os.Stderr.WriteString(rest[1])
		_, _ = os.Stdout.WriteString(rest[0])
		os.Exit(0)

	// stderr-n writes N bytes on stderr, all distinguishable by position,
	// then a terminal marker, so the LAST 4 KiB is assertable.
	case "stderr-n":
		n, _ := strconv.Atoi(rest[0])
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = byte('a' + (i % 26))
		}
		_, _ = os.Stderr.Write(buf)
		os.Exit(0)

	// no-stdin never reads stdin and sleeps, so a stdin payload past the
	// pipe buffer blocks unless Cmd's own WaitDelay bounds the write (S4).
	case "no-stdin":
		d, _ := time.ParseDuration(rest[0])
		time.Sleep(d)

	// ingroup is the EXISTING grandchild-leak oracle's shape, re-run under
	// the bound (S2): a grandchild INSIDE the group that holds stdout and
	// is therefore reachable by reapGroup.
	case "ingroup":
		self, _ := os.Executable()
		sub := exec.Command(self, "hold", "20s")
		sub.Stdout = os.Stdout
		if err := sub.Start(); err != nil {
			os.Exit(95)
		}
		_ = os.WriteFile(rest[0], []byte(strconv.Itoa(sub.Process.Pid)), 0o600)
		_, _ = os.Stdout.WriteString(rest[1])
		os.Exit(0)

	default:
		os.Exit(98)
	}
}
`

// --- entry construction ---------------------------------------------------

// escEntry builds a command-backed read entry over FX-deadline-escape.
func escEntry(t *testing.T, timeout string, argv ...string) table.Accessor {
	t.Helper()

	return table.Accessor{
		Role:    escRole,
		Keys:    []string{escKey},
		Timeout: timeout,
		Command: append([]string{escapeBin(t)}, argv...),
	}
}

// escReader builds the reader every drain scenario drives. The binding is
// the real one; the child is the real fixture.
func escReader(t *testing.T, timeout string, argv ...string) cmdbind.Reader {
	t.Helper()

	return cmdbind.Reader{
		Accessor: escEntry(t, timeout, argv...),
		Name:     escName,
		Config:   cmdbind.Config{AllowCommands: true},
	}
}

// escArtifact returns an ABSOLUTE artifact handle; C2 refuses a relative
// one, so every positive fixture must supply one.
func escArtifact(t *testing.T) accessor.Artifact {
	t.Helper()

	p := filepath.Join(t.TempDir(), "state.json")
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("absolutizing the artifact path failed: %v", err)
	}
	return accessor.Artifact{Role: escRole, Path: abs}
}

// escPIDFile returns a path the fixture records its escaped grandchild's
// pid in.
func escPIDFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "escapee.pid")
}

// --- oracles --------------------------------------------------------------

// escAlive reports whether a pid names a live process. It is the oracle for
// C1 `residue:`'s admitted leak — the escapee OUTLIVES the refusal.
func escAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// escReadPID reads the escaped grandchild's pid, waiting briefly for the
// fixture to have written it. It never fabricates one: a missing pidfile is
// a fixture failure, not a pass.
func escReadPID(t *testing.T, path string) int {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)
	for {
		b, err := os.ReadFile(path)
		if err == nil {
			n, cerr := strconv.Atoi(strings.TrimSpace(string(b)))
			if cerr == nil {
				return n
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("FX-deadline-escape recorded no grandchild pid at %s", path)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// escKill tears the admitted residue down. Every scenario that asserts the
// escapee survived is obliged to kill it — C1 `residue:` admits the leak,
// the suite does not inherit it.
func escKill(pid int) {
	if pid <= 1 {
		return
	}
	if p, err := os.FindProcess(pid); err == nil {
		_ = p.Signal(syscall.SIGKILL)
	}
}

// escExecError extracts the typed refusal carrier every command-binding
// refusal rides (`0025:C4`; C1 `refusal:` "a held pipe refuses through
// `*accessor.ExecError`"). It is `errors.As` and not a type assertion
// because the contract WRAPS rather than flattens.
func escExecError(t *testing.T, err error) *accessor.ExecError {
	t.Helper()

	var ee *accessor.ExecError
	if !errors.As(err, &ee) {
		t.Fatalf("the refusal is not an *accessor.ExecError: %#v (%v)", err, err)
	}
	return ee
}

// escHeldDetail reports the refusal's Detail. The contract fixes the pipe
// NAMES and their fixed order; it explicitly leaves the identifier and Go
// spelling of the report type and of the fold to the implementer (REQ-13,
// REQ-65), so the oracle is the user-visible text and never a type name.
func escHeldDetail(t *testing.T, err error) string {
	t.Helper()
	return escExecError(t, err).Detail
}

// escNamesHeld reports whether a refusal text names the held pipe(s) in the
// fixed order C1 `precedence:` states — "stdout, stderr" for both, and the
// BARE pipe name for one (S10, MVV row 5).
func escNamesHeld(text string, pipes ...string) bool {
	switch len(pipes) {
	case 1:
		return strings.Contains(text, pipes[0])
	case 2:
		return strings.Contains(text, pipes[0]+", "+pipes[1])
	default:
		return false
	}
}

// escRaw selects the entry's raw output mode, so a payload test asserts on
// the child's stdout BYTES rather than on a JSON envelope.
func escRaw() *string {
	s := "raw"
	return &s
}

// openFDCount counts the process's own open descriptors. It is the oracle
// for C1 `precedence:`'s "no read end outlives `spawn`" — a leak grows it,
// and a double-close would surface as a failure elsewhere rather than here.
func openFDCount(t *testing.T) int {
	t.Helper()

	// /dev/fd is the portable-enough listing on darwin and linux, the two
	// OSes the record's spikes cover.
	entries, err := os.ReadDir("/dev/fd")
	if err != nil {
		t.Skipf("this host exposes no /dev/fd listing: %v", err)
	}
	return len(entries)
}
