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
// (`0025:C4`, A1's necessary-and-sufficient triple).
const WaitDelay = 500 // milliseconds

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
	// exited reports that the process RAN and exited. The two exit maps are
	// consulted only for a process that ran, so a spawn failure — which has
	// no exit code at all — can never be claimed by them however broadly
	// they are written (`0025:C3`).
	exited bool
}

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
		return invocation{}, refuse("command execution requires the " +
			"`allow_commands` opt-in (--allow-commands); the accessor `" +
			name + "` declares a command and none was given")
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

	// A1's necessary-and-sufficient triple. Dropping the group signal
	// orphans a grandchild holding the pipe; dropping WaitDelay hangs
	// `Wait` on it, and it is also what bounds the stdin write so a
	// non-reading child cannot block the parent (`0025:C4`).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		// The GROUP, not only the direct child.
		if kerr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); kerr != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
	cmd.WaitDelay = WaitDelay * time.Millisecond

	// The stdout bound is enforced by BOUNDING THE READ, so an unbounded
	// child cannot exhaust memory. One byte past the cap is read on
	// purpose: overflow is DETECTED here rather than silently truncated.
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return invocation{}, wrap("", err)
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return invocation{}, wrap("", err)
	}

	if serr := cmd.Start(); serr != nil {
		return invocation{}, wrap("", serr)
	}

	stdout := readBounded(outPipe, StdoutCap+1)
	// The stderr channel is read under the SAME bound as stdout so a
	// verbose tool cannot exhaust memory, and `tail` then keeps the LAST
	// 4 KiB — which is the diagnosis a caller wants when a tool says a lot
	// before it fails (`0025:C4`).
	stderr := readBounded(errPipe, StdoutCap)

	werr := cmd.Wait()

	inv := invocation{stdout: stdout, stderr: tail(stderr)}
	var ee *exec.ExitError
	switch {
	case werr == nil:
		inv.exited = true
	case errors.As(werr, &ee):
		// `Wait` reports a deadline kill as an ExitError too, so the
		// deadline is classified from ctx.Err() by the caller and never
		// from this error (`0025:C4`).
		inv.exited = true
		inv.exitCode = ee.ExitCode()
	default:
		return inv, wrap(inv.stderr, werr)
	}

	if len(inv.stdout) > StdoutCap {
		return inv, wrap(inv.stderr, errors.New("the child's stdout exceeded "+
			"the "+strconv.Itoa(StdoutCap)+" byte bound"))
	}
	return inv, nil
}

// readBounded drains r up to limit bytes, reporting whether more remained.
func readBounded(r io.Reader, limit int) []byte {
	b, err := io.ReadAll(io.LimitReader(r, int64(limit)))
	if err != nil {
		return b
	}
	// Drain the remainder so the child is never blocked on a full pipe: the
	// bound is on what is KEPT, and a child left writing into a full pipe
	// would hang past its deadline rather than being reaped. Overflow is
	// DETECTED from the kept length — the read is bounded at cap+1, so a
	// stdout of exactly the cap is a value and one byte more is a failure.
	_, _ = io.Copy(io.Discard, r)
	return b
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
	inv, err := spawn(ctx, r.Accessor, r.Name, accessor.CapRead, r.Config, art,
		stdinObject(nil))
	if err != nil {
		return nil, nil, err
	}

	// Ordering is normative: a NON-EMPTY stdout is parsed first, and the
	// exit maps apply only to an empty stdout (`0025:C4`).
	if len(inv.stdout) != 0 {
		return r.parse(inv, requested)
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
	return nil, nil, wrap(inv.stderr, errors.New(
		"the read command produced no stdout and exited "+
			strconv.Itoa(inv.exitCode)+
			", which the entry's `exit_absent` does not list"))
}

// parse maps a non-empty stdout onto the split return, per the entry's
// declared output mode.
func (r Reader) parse(inv invocation, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	if r.Accessor.Output != nil && *r.Accessor.Output == "raw" {
		// The single declared key's value is stdout minus EXACTLY ONE
		// trailing "\n" — no trimming, no case folding, no whitespace
		// normalization (`0025:C3`, FX-raw-read).
		if len(requested) != 1 {
			return nil, nil, wrap(inv.stderr, errors.New(
				"raw mode carries one declared key and this read requested "+
					strconv.Itoa(len(requested))))
		}
		value := string(inv.stdout)
		value = strings.TrimSuffix(value, "\n")
		return []accessor.KeyValue{{Key: requested[0], Value: value}}, nil, nil
	}

	var obj map[string]string
	if jerr := json.Unmarshal(inv.stdout, &obj); jerr != nil {
		return nil, nil, wrap(inv.stderr, jerr)
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
			return "", "", wrap(inv.stderr, jerr)
		}
		if !slices.Contains(accessor.Verdicts(), accessor.Verdict(env.Verdict)) {
			return "", "", wrap(inv.stderr, errors.New(
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
	return "", "", wrap(inv.stderr, errors.New(
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
	if inv.exitCode != 0 {
		return wrap(inv.stderr, errors.New(
			"the write command exited "+strconv.Itoa(inv.exitCode)))
	}
	return nil
}

// Invocations counts Apply ENTRIES — one per call, no internal respawn
// (`0025:C7`). The count is about the accessor LAYER's re-entry, not about
// process success, so a failing spawn still counts one.
func (w *Writer) Invocations() int { return w.invocations }
