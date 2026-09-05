package cmdbind_test

// RDR 0025 — shared fixtures for the command-binding suite.
//
// Nothing here mocks the unit under test. Every fixture is a REAL child
// process: the suite builds a tiny helper binary once and declares it as
// `argv0`, so the binding spawns, writes stdin, reads stdout, and reaps a
// genuine process group. A fake `exec` seam would make the C4 deadline
// triple and the C4 env allowlist unassertable — those are properties of
// a real spawn or they are nothing.
//
// The helper's behaviour is selected by argv, so one binary covers every
// scenario the RDR names, and the argv the child OBSERVES is the oracle
// for C2's authority bound: the child echoes its own os.Args and os.Environ
// back, which is the only way to assert "the executed argv is exactly the
// declared vector after substitution" without reading the binding's
// internals.

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

const (
	fxRole    = "state"
	fxKey     = "flow.status"
	fxKeyB    = "flow.stage"
	fxName    = "cmd.state"
	fxTimeout = "2s"
)

// --- the helper child ----------------------------------------------------

var (
	helperOnce sync.Once
	helperPath string
	helperErr  error
)

// helperBin builds (once) the fixture child and returns its absolute path.
// It is a real binary because every clause under test — the process group,
// the stdin pipe, the env the child observes — is a property of a real
// spawn.
func helperBin(t *testing.T) string {
	t.Helper()

	helperOnce.Do(func() {
		dir, err := os.MkdirTemp("", "cmdbind-helper")
		if err != nil {
			helperErr = err
			return
		}
		src := filepath.Join(dir, "main.go")
		if werr := os.WriteFile(src, []byte(helperSource), 0o600); werr != nil {
			helperErr = werr
			return
		}
		out := filepath.Join(dir, "cmdhelper")
		cmd := exec.Command("go", "build", "-o", out, src)
		cmd.Env = append(os.Environ(), "GO111MODULE=off")
		if b, berr := cmd.CombinedOutput(); berr != nil {
			helperErr = errBuild{berr, string(b)}
			return
		}
		helperPath = out
	})
	if helperErr != nil {
		t.Fatalf("building the fixture child failed: %v", helperErr)
	}
	return helperPath
}

type errBuild struct {
	err error
	out string
}

func (e errBuild) Error() string { return e.err.Error() + ": " + e.out }

// helperSource is the fixture child. Its first argument selects a mode; the
// rest are mode arguments. Every mode is one of the shapes the RDR's
// scenarios name.
const helperSource = `package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "cmdhelper: no mode")
		os.Exit(97)
	}
	mode := os.Args[1]
	rest := os.Args[2:]

	switch mode {
	// echo-argv writes the child's OWN os.Args back as a JSON array on
	// stdout. It is the oracle for C2's authority bound.
	case "echo-argv":
		_ = json.NewEncoder(os.Stdout).Encode(os.Args)

	// echo-env writes the child's OWN environment back as a JSON object.
	// It is the oracle for C4's env allowlist.
	case "echo-env":
		env := map[string]string{}
		for _, kv := range os.Environ() {
			if k, v, ok := strings.Cut(kv, "="); ok {
				env[k] = v
			}
		}
		_ = json.NewEncoder(os.Stdout).Encode(env)

	// echo-stdin writes whatever arrived on stdin back on stdout verbatim.
	case "echo-stdin":
		b, _ := io.ReadAll(os.Stdin)
		_, _ = os.Stdout.Write(b)

	// stdin-to-file drains stdin into the named file and exits zero. It is
	// the honest wrapper shape: envelope in, tool-native write out.
	case "stdin-to-file":
		b, _ := io.ReadAll(os.Stdin)
		_ = os.WriteFile(rest[0], b, 0o600)

	// emit writes its literal argument on stdout and exits with the code
	// its second argument names.
	case "emit":
		_, _ = os.Stdout.WriteString(rest[0])
		code, _ := strconv.Atoi(rest[1])
		os.Exit(code)

	// emit-stderr writes its argument on STDERR and exits with the code.
	case "emit-stderr":
		_, _ = os.Stderr.WriteString(rest[0])
		code, _ := strconv.Atoi(rest[1])
		os.Exit(code)

	// silent writes nothing and exits with the code its argument names.
	case "silent":
		code, _ := strconv.Atoi(rest[0])
		os.Exit(code)

	// self-signal writes nothing and SIGKILLs itself, so Wait reports a
	// signal and NO exit code at all. It is the "the exit map cannot
	// apply" fixture: a signalled child never reaches an exit status for
	// exit_absent or exit_verdicts to claim.
	case "self-signal":
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGKILL)
		time.Sleep(10 * time.Second)

	// flood writes N bytes on stdout and exits zero.
	case "flood":
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

	// sleep blocks past its argument's duration without reading stdin.
	case "sleep":
		d, _ := time.ParseDuration(rest[0])
		time.Sleep(d)

	// orphan-holds-pipe spawns a GRANDCHILD that holds stdout open and
	// sleeps, then the direct child exits. Without process-group
	// termination the grandchild outlives the deadline holding the pipe.
	case "orphan-holds-pipe":
		self, _ := os.Executable()
		sub := exec.Command(self, "sleep", rest[0])
		sub.Stdout = os.Stdout
		_ = sub.Start()
		// Record the grandchild's pid so the test can assert it did not
		// survive the refusal.
		if len(rest) > 1 {
			_ = os.WriteFile(rest[1], []byte(strconv.Itoa(sub.Process.Pid)), 0o600)
		}
		time.Sleep(50 * time.Millisecond)

	// never-reads-stdin ignores stdin entirely and sleeps, so a large
	// stdin write blocks unless the write itself is bounded.
	case "never-reads-stdin":
		d, _ := time.ParseDuration(rest[0])
		time.Sleep(d)

	// trace appends a line to the named file, proving the child RAN. It is
	// the "absence of a spawn" oracle's positive control.
	case "trace":
		f, err := os.OpenFile(rest[0], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			os.Exit(96)
		}
		_, _ = f.WriteString("ran\n")
		_ = f.Close()

	default:
		fmt.Fprintln(os.Stderr, "cmdhelper: unknown mode "+mode)
		os.Exit(98)
	}
}
`

// --- entry construction ---------------------------------------------------

// entry builds one command-backed accessor entry.
func entry(role string, argv []string, keys ...string) table.Accessor {
	if len(keys) == 0 {
		keys = []string{fxKey}
	}
	return table.Accessor{
		Role:    role,
		Keys:    keys,
		Timeout: fxTimeout,
		Command: argv,
	}
}

func strptr(s string) *string { return &s }

// allowed is the Config every non-gate test uses: commands permitted, and
// a base dir that is the helper's own directory so a separator-bearing
// argv0 resolves.
func allowed(t *testing.T) cmdbind.Config {
	t.Helper()
	return cmdbind.Config{BaseDir: filepath.Dir(helperBin(t)), AllowCommands: true}
}

// artifactAt returns an ABSOLUTE artifact handle under the test's temp dir.
// C2 refuses a relative one, so every positive fixture must be absolute.
func artifactAt(t *testing.T, name string) accessor.Artifact {
	t.Helper()

	p := filepath.Join(t.TempDir(), name)
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("absolutizing the artifact path failed: %v", err)
	}
	return accessor.Artifact{Role: fxRole, Path: abs}
}

// --- oracles --------------------------------------------------------------

// observedArgv runs a read whose child is `echo-argv` and returns the argv
// the CHILD observed. It is how C2's byte-for-byte claim is asserted
// without reading the binding's internals.
func observedArgv(t *testing.T, raw string) []string {
	t.Helper()

	var argv []string
	if err := json.Unmarshal([]byte(raw), &argv); err != nil {
		t.Fatalf("the child did not echo a JSON argv array: %v\n%s", err, raw)
	}
	return argv
}

// observedEnv decodes the child's echoed environment.
func observedEnv(t *testing.T, raw string) map[string]string {
	t.Helper()

	env := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		t.Fatalf("the child did not echo a JSON env object: %v\n%s", err, raw)
	}
	return env
}

// alive reports whether a pid names a live process. It is the oracle for
// "no process from the child's group survives the refusal".
func alive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

// pidFrom reads the grandchild pid the `orphan-holds-pipe` mode recorded.
func pidFrom(t *testing.T, path string) int {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the fixture recorded no grandchild pid: %v", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatalf("the recorded pid %q is not a number: %v", b, err)
	}
	return n
}
