package cli

// RDR 0026 — the Minimum Viable Validation, end to end (`0026:MVV`).
//
// This is the runnable integration proof and it is asserted at the
// `flow resolve` BOUNDARY, not only at the binding: "the defect was
// reported there (\"still blocked after 15 s\"), so the row also asserts
// that `flow resolve` itself terminates and emits the class in its
// envelope. A binding-level pass with a hung CLI would not close the
// original report" (REQ-MVV.3a).
//
// The fixture is FX-deadline-escape, built here as a REAL Go binary
// spawning a REAL `SysProcAttr{Setsid: true}` grandchild, because macOS
// ships no `setsid` binary and a shell fixture does not escape there
// (Background, REQ-110). A green exit code is explicitly NOT the oracle at
// any row: every row asserts the CLASS, the ELAPSED time against the C1
// `bound:` committed total, and — on the rows that have one — the recovered
// value BYTE-FOR-BYTE.

import (
	"encoding/json"
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

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
)

// mvvWaitDelay / mvvTolerance restate the C1 `bound:` figures from the
// package's own constant, so a row moves with the constant rather than
// against a transcribed literal.
var (
	mvvWaitDelay = time.Duration(cmdbind.WaitDelay) * time.Millisecond
	mvvTolerance = 100 * time.Millisecond
)

// REQ-MVV (`0026:MVV`): "Fixture FX-deadline-escape, deadline shape, run
// through the real `cmdbind.Reader.Read` under `--allow-commands` with a
// declared `timeout = \"2s\"`"
// REQ-MVV.1: "The child (a Go helper) spawns a grandchild with
// `SysProcAttr{Setsid: true}` that inherits stdout and stderr and sleeps
// 60 s; the child itself writes nothing and sleeps past the deadline."
// REQ-MVV.2: "`flow resolve` invokes the reader; the context deadline
// expires at 2 s; `Cancel` kills the group; `Wait` returns; `reapGroup`
// signals the group; the drains are joined under the bound."
// REQ-MVV.3: "the read returns within `2 s + 2·500 ms` (Normative: elapsed
// ≤ 3 s, asserted with the C1 `bound:` tolerance), the executor classifies
// `timeout`, and the refusal's `Detail` is empty as 0025:C4 states for
// `timeout`."
// REQ-MVV.3a: "Asserted at the `flow resolve` boundary, not only at the
// binding ... the row also asserts that `flow resolve` itself terminates
// and emits the class in its envelope."
// REQ-MVV.4: "The grandchild is still alive (asserted — the admitted
// residue) and the test kills it."
// REQ-58 (AP): "the refusal is never withheld, and the invocation returns
// within `timeout + 2·WaitDelay`."
// REQ-133 (TS): "the FX-deadline-escape helper is exercised on darwin and
// linux in CI, so the escape is real on both"
// HAPPY PATH — "happy" in the sense that the CONTRACT's promised outcome is
// what is asserted: a bounded, classified refusal where today there is a
// hang.
func TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes(t *testing.T) {
	dir := t.TempDir()
	helper := mvvEscapeBin(t)
	pidfile := filepath.Join(dir, "escapee.pid")
	artifact := filepath.Join(dir, "state.json")
	if err := os.WriteFile(artifact, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("seeding the artifact: %v", err)
	}
	model := mvvWriteModel(t, dir, "escape.toml",
		mvvModelSource(helper, "deadline", pidfile, "2s"))

	// --- MVV rows 1-3 -----------------------------------------------------
	start := time.Now()
	stdout, _, err := runCmd(t, "flow", "resolve",
		"--model", model,
		"--artifact", "state="+artifact,
		"--outcome", "advance",
		"--allow-commands",
		"--as=json")
	elapsed := time.Since(start)

	pid := mvvReadPID(t, pidfile)
	defer mvvKill(pid)

	// Row 3a: `flow resolve` itself TERMINATES. The reported defect was a
	// CLI still blocked after 15 s; a binding-level pass with a hung CLI
	// would not close it.
	if err == nil {
		t.Fatalf("`flow resolve` succeeded against a reader whose pipes are "+
			"held by a `setsid(2)` escapee (stdout: %s). The refusal is never "+
			"withheld — C1 `residue:` admits leakage of a PROCESS, never a "+
			"withheld refusal", stdout)
	}

	// Row 3: the Normative bound. timeout (2 s) + 2·WaitDelay, with the C1
	// `bound:` scheduling tolerance.
	want := 2*time.Second + 2*mvvWaitDelay + mvvTolerance
	if elapsed > want {
		t.Fatalf("`flow resolve` returned after %v, past the C1 `bound:` "+
			"committed total of timeout + 2·WaitDelay (%v) plus its stated "+
			"%v tolerance. The originally reported defect is exactly this: "+
			"\"still blocked after 15 s\", deadline (2 s) and WaitDelay "+
			"(500 ms) long expired",
			elapsed, 2*time.Second+2*mvvWaitDelay, mvvTolerance)
	}

	// Row 3: the executor classifies `timeout` — the deadline-first rule.
	code := clierr.ErrorCode(err)
	if code != codeAccessorTimeout && code != codeReadBackTimeout {
		t.Fatalf("code = %q; want a TIMEOUT code. C1 `class:` fixes the class "+
			"as the executor's: ctx `DeadlineExceeded` ⇒ `timeout`, read "+
			"FIRST, and only then the binding's error ⇒ `execution_failure`",
			code)
	}

	// Row 3: the refusal's `Detail` is EMPTY, as 0025:C4 states for
	// `timeout`. This is F6, an accepted trade-off recorded by name: the
	// held-pipe reason is computed and then not rendered.
	if detail := mvvEnvelopeDetail(t, err); detail != "" {
		t.Errorf("the `timeout` refusal's detail = %q; 0025:C4 gives "+
			"`timeout` no Err/Detail, so the held-pipe reason survives only "+
			"on `execution_failure` (F6, S6). This record does NOT amend "+
			"0025:C4's `detail:` rule", detail)
	}

	// Row 4: the grandchild is STILL ALIVE — the admitted residue — and the
	// test kills it.
	if !mvvAlive(pid) {
		t.Errorf("the `setsid(2)` grandchild pid %d did not outlive the "+
			"refusal. C1 `residue:` restates 0025:F4: a process outside the "+
			"child's process group MAY outlive the refusal, and nothing in "+
			"this CLI reaps it. If it died, this fixture is not exercising "+
			"the held-pipe path at all", pid)
	}
}

// REQ-MVV.5: "Success shape, same fixture: the child writes a whole
// envelope and exits 0 while the grandchild holds the pipes — expected
// `execution_failure` within `2·500 ms` of the exit, `Err` naming
// \"stdout, stderr\" and \"exited 0\", and the envelope NOT parsed (no
// values returned)."
// REQ-101 (MC `oracle`): "a build that parses held stdout returns values
// and fails the not-parsed assertion"
// REQ-94 (MC `disposition`): "Held pipe, deadline NOT elapsed" →
// "`execution_failure`", "the held-pipe error ... in `Err`, reason + tail
// in `Detail`", "no stdout", "LOUD"
// REQ-74 (LBD grace): "the tail reaches the user on the
// `execution_failure` leg (MVV row 5, S10)"
// ADVERSARIAL — the child answers CORRECTLY and exits 0, so a build that
// parses held stdout returns a value and passes everything except this row.
func TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing(t *testing.T) {
	dir := t.TempDir()
	helper := mvvEscapeBin(t)
	pidfile := filepath.Join(dir, "escapee.pid")
	artifact := filepath.Join(dir, "state.json")
	if err := os.WriteFile(artifact, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("seeding the artifact: %v", err)
	}
	// The declared timeout is generous, so the deadline never elapses and
	// the class is `execution_failure` rather than `timeout` — the leg on
	// which the held-pipe reason is REACHABLE.
	model := mvvWriteModel(t, dir, "success.toml",
		mvvModelSource(helper, "success", pidfile, "20s"))

	start := time.Now()
	stdout, _, err := runCmd(t, "flow", "resolve",
		"--model", model,
		"--artifact", "state="+artifact,
		"--outcome", "advance",
		"--allow-commands",
		"--as=json")
	elapsed := time.Since(start)

	pid := mvvReadPID(t, pidfile)
	defer mvvKill(pid)

	if err == nil {
		t.Fatalf("a child that wrote a whole envelope and exited 0 while a "+
			"`setsid(2)` grandchild held both pipes was ACCEPTED (stdout: "+
			"%s). C1 `refusal:` fixes that the returned invocation carries no "+
			"stdout and a held stdout is never parsed on ANY path — refused, "+
			"not truncated-then-read", stdout)
	}
	// "within `2·500 ms` of the exit" — the child exits immediately, so the
	// invocation's own start is the anchor, plus the stated tolerance.
	if want := 2*mvvWaitDelay + mvvTolerance; elapsed > want {
		t.Fatalf("the success shape returned after %v; MVV row 5 fixes it "+
			"within 2·WaitDelay (%v) of the child's exit plus the C1 `bound:` "+
			"tolerance (%v)", elapsed, 2*mvvWaitDelay, mvvTolerance)
	}
	if code := clierr.ErrorCode(err); code != codeAccessorFailed {
		t.Errorf("code = %q; want %q — the deadline never elapsed, so the "+
			"executor's deadline-first rule falls through to the binding's "+
			"error ⇒ `execution_failure`", code, codeAccessorFailed)
	}
	// `Err` names the held pipes in the FIXED order, and the direct child's
	// exit status.
	detail := mvvEnvelopeDetail(t, err)
	if !strings.Contains(detail, "stdout, stderr") {
		t.Errorf("detail = %q; MVV row 5 fixes `Err` naming \"stdout, "+
			"stderr\" — the fixed order C1 `precedence:` states", detail)
	}
	if !strings.Contains(detail, "exited 0") {
		t.Errorf("detail = %q; MVV row 5 fixes `Err` naming \"exited 0\" — "+
			"the direct child's exit status, composed from the fields `spawn` "+
			"already populates from `werr`", detail)
	}
	// The envelope is NOT parsed: no values are reported.
	if mvvReportsAnyTag(t, stdout) {
		t.Errorf("`flow resolve` reported tags from a HELD stdout:\n%s\nA "+
			"build that parses held stdout returns values and fails exactly "+
			"this assertion (MC `oracle`, MVV row 5's named control)", stdout)
	}
}

// --- FX-deadline-escape, built once for the CLI-level rows ---------------

var (
	mvvOnce sync.Once
	mvvPath string
	mvvErr  error
)

// mvvEscapeBin builds FX-deadline-escape for the CLI-boundary rows. It is
// the SAME fixture the binding-level suite drives (REQ-113: "Add
// FX-deadline-escape (both shapes) beside the existing grandchild-leak and
// cancel oracles"), rebuilt here because the two packages cannot share a
// test helper.
func mvvEscapeBin(t *testing.T) string {
	t.Helper()

	mvvOnce.Do(func() {
		dir, err := os.MkdirTemp("", "mvv-escape")
		if err != nil {
			mvvErr = err
			return
		}
		src := filepath.Join(dir, "main.go")
		if werr := os.WriteFile(src, []byte(mvvEscapeSource), 0o600); werr != nil {
			mvvErr = werr
			return
		}
		out := filepath.Join(dir, "escapehelper")
		cmd := exec.Command("go", "build", "-o", out, src)
		cmd.Env = append(os.Environ(), "GO111MODULE=off")
		if b, berr := cmd.CombinedOutput(); berr != nil {
			mvvErr = &mvvBuildErr{berr, string(b)}
			return
		}
		mvvPath = out
	})
	if mvvErr != nil {
		t.Fatalf("building FX-deadline-escape failed: %v", mvvErr)
	}
	return mvvPath
}

type mvvBuildErr struct {
	err error
	out string
}

func (e *mvvBuildErr) Error() string { return e.err.Error() + ": " + e.out }

// mvvEscapeSource is FX-deadline-escape's two MVV shapes. The grandchild is
// spawned with a REAL `syscall.Setsid`, which is the only escape that works
// on darwin as well as linux (Background).
const mvvEscapeSource = `package main

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// escape starts the grandchild and leaves the pid bookkeeping to it. This
// process is in the group the declared deadline SIGKILLs, so a write here
// can lose that race under load even though the escape already happened.
func escape(pidfile string) {
	self, _ := os.Executable()
	sub := exec.Command(self, "hold", pidfile)
	sub.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	sub.Stdout = os.Stdout
	sub.Stderr = os.Stderr
	if err := sub.Start(); err != nil {
		os.Exit(95)
	}
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(97)
	}
	switch os.Args[1] {
	// hold is the GRANDCHILD: it holds the inherited stdout and stderr and
	// sleeps 60 s, writing nothing to them. Nothing in the CLI can reach it.
	// Its first act is to record its OWN pid: it is in a new session from
	// the moment it starts, so the group kill cannot skip the write, and
	// the temp file + rename means a reader sees no file or a whole one.
	case "hold":
		if len(os.Args) > 2 {
			tmp := os.Args[2] + ".tmp"
			if os.WriteFile(tmp, []byte(strconv.Itoa(os.Getpid())), 0o600) == nil {
				_ = os.Rename(tmp, os.Args[2])
			}
		}
		time.Sleep(60 * time.Second)

	// deadline is MVV rows 1-4: the grandchild inherits BOTH pipes and
	// sleeps; the child itself writes nothing and sleeps past the deadline.
	case "deadline":
		escape(os.Args[2])
		time.Sleep(60 * time.Second)

	// success is MVV row 5: the child writes a WHOLE, well-formed envelope
	// and exits 0 while the grandchild holds the pipes.
	case "success":
		escape(os.Args[2])
		_, _ = os.Stdout.WriteString(` + "`" + `{"status":"final"}` + "`" + `)
		os.Exit(0)

	default:
		os.Exit(98)
	}
}
`

// --- model and envelope helpers -------------------------------------------

// mvvModelSource declares one command-backed reader over FX-deadline-escape
// with the given declared `timeout`.
func mvvModelSource(helper, mode, pidfile, timeout string) string {
	return `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "escape0026"
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

[read.state]
role = "state"
command = [` + strconv.Quote(helper) + `, ` + strconv.Quote(mode) + `, ` +
		strconv.Quote(pidfile) + `]
keys = ["status"]
timeout = ` + strconv.Quote(timeout) + `

# The writer is path-backed and is NOT under test: the model needs exactly
# one writer for the owned tag to load, and the reader is the command
# binding whose drain the scenario holds.
[write.state]
role = "state"
path = "status"
keys = ["status"]
timeout = "10s"
read_back = true

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
}

func mvvWriteModel(t *testing.T, dir, name, src string) string {
	t.Helper()

	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(src), 0o600); err != nil {
		t.Fatalf("writing the model: %v", err)
	}
	return p
}

// mvvEnvelopeDetail reads the CLI error envelope's `detail` field — the
// observable S1 and the MVV rows both assert on, per
// docs/cli-output-contract.md's error envelope.
func mvvEnvelopeDetail(t *testing.T, err error) string {
	t.Helper()

	var ce *clierr.CLIError
	if !errors.As(err, &ce) {
		t.Fatalf("the CLI error is not a *clierr.CLIError: %v", err)
	}
	return ce.Detail
}

// mvvReportsAnyTag reports whether `flow resolve`'s envelope carries any
// tag value at all. A held stdout must yield NONE.
func mvvReportsAnyTag(t *testing.T, stdout string) bool {
	t.Helper()

	if strings.TrimSpace(stdout) == "" {
		return false
	}
	var env struct {
		Type string `json:"type"`
		Data struct {
			Tags  map[string]string `json:"tags"`
			Owned map[string]string `json:"owned"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		// A non-JSON stdout carries no parsed tags either.
		return false
	}
	return len(env.Data.Tags) != 0 || len(env.Data.Owned) != 0
}

// --- process oracles ------------------------------------------------------

func mvvAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

func mvvKill(pid int) {
	if pid <= 1 {
		return
	}
	if p, err := os.FindProcess(pid); err == nil {
		_ = p.Signal(syscall.SIGKILL)
	}
}

// mvvReadPID reads the escaped grandchild's pid, waiting briefly for the
// fixture to have written it. A missing pidfile is a fixture failure, never
// a pass.
func mvvReadPID(t *testing.T, path string) int {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)
	for {
		if b, err := os.ReadFile(path); err == nil {
			if n, cerr := strconv.Atoi(strings.TrimSpace(string(b))); cerr == nil {
				return n
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("FX-deadline-escape recorded no grandchild pid at %s — "+
				"the fixture did not escape, so nothing was under test", path)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
