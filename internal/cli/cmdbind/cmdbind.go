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
//   - execution requires the `--allow-commands` opt-in, checked in the
//     CONSTRUCTOR so no execution path can reach a spawn without it
//     (`0025:C6`).
//
// PHASE 1 DECLARATION ONLY. The types and constructors below carry the
// spec-named surface so the RDR 0025 conformance suite compiles and is red
// on BEHAVIOUR; Phase 2 fills in the bodies.
package cmdbind

import (
	"context"
	"runtime"

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

// Reader is a command-backed `accessor.ReadBinding`.
type Reader struct {
	Accessor table.Accessor
	Name     string
	Config   Config
}

// Capability reports the capability this binding serves.
func (Reader) Capability() accessor.Capability { return accessor.CapRead }

// Read invokes the declared command and maps its envelope onto the seam's
// split return (`0025:C3`, `0025:C7`).
func (r Reader) Read(_ context.Context, _ accessor.Artifact, _ []string) (
	values []accessor.KeyValue, unreadable []string, err error,
) {
	return nil, nil, errNotImplemented
}

// Gate is a command-backed `accessor.GateBinding`.
type Gate struct {
	Accessor table.Accessor
	Name     string
	Config   Config
}

// Capability reports the capability this binding serves.
func (Gate) Capability() accessor.Capability { return accessor.CapGate }

// Gate invokes the declared command and reads its verdict envelope or its
// declared exit map (`0025:C3`).
func (g Gate) Gate(_ context.Context, _ accessor.Artifact) (accessor.Verdict, string, error) {
	return "", "", errNotImplemented
}

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
// (`0025:C3`, `0004:C12`).
func (w *Writer) Apply(_ context.Context, _ accessor.Artifact, _ []resolve.Tag) error {
	w.invocations++
	return errNotImplemented
}

// Invocations counts Apply ENTRIES — one per call, no internal respawn
// (`0025:C7`).
func (w *Writer) Invocations() int { return w.invocations }

// errNotImplemented marks the Phase 1 declaration surface. It is a typed
// *accessor.ExecError so the shape of a real refusal is already correct;
// Phase 2 replaces every return site with the real invocation.
var errNotImplemented = &accessor.ExecError{
	Detail: "cmdbind: PHASE 1 DECLARATION ONLY — no invocation is implemented",
}
