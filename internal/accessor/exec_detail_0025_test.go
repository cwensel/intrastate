package accessor_test

// RDR 0025 — the invocation-error carrier and the seam it must not change
// (C4 `detail`, C7; REQ-50, REQ-55..REQ-58, REQ-62..REQ-69, REQ-98..REQ-108,
// REQ-137).
//
// The coverage floor this file exists to hold, and the mutant it kills:
// threading the error parameter through `refusalOf` ALONE reaches gate and
// write but silently drops the tail on READS — the one capability
// established tools bind directly, and so the likeliest source of a stderr
// tail (A11, refuted the single-hook form of the claim). So the read arm
// and the gate/write arms are asserted SEPARATELY and neither stands in for
// the other.
//
// Nothing here mocks the executor. The bindings are fixtures returning a
// typed *accessor.ExecError through the interfaces' existing bare `error`
// slot, which is exactly what a command binding does; the oracle is the
// Refusal the real executor mints.

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// --- fixtures ------------------------------------------------------------

const (
	execRole   = "state"
	execReader = "cmd.read"
	execGate   = "cmd.gate"
	execWriter = "cmd.write"
	execKey    = "status"

	execTimeout = 50 * time.Millisecond
)

// execFailBinding fails every invocation with the supplied error. It is a
// stand-in for a command binding whose child exited non-zero: the error is
// the ONLY channel the seam gives it, which is the whole point of C4's
// typed carrier.
type execFailBinding struct {
	cap accessor.Capability
	err error

	applies int
}

func (b *execFailBinding) Capability() accessor.Capability { return b.cap }

func (b *execFailBinding) Read(context.Context, accessor.Artifact, []string) (
	[]accessor.KeyValue, []string, error,
) {
	return nil, nil, b.err
}

func (b *execFailBinding) Gate(context.Context, accessor.Artifact) (
	accessor.Verdict, string, error,
) {
	return "", "", b.err
}

func (b *execFailBinding) Apply(context.Context, accessor.Artifact, []resolve.Tag) error {
	b.applies++
	return b.err
}

func (b *execFailBinding) Invocations() int { return b.applies }

// execDef builds one definition over a fixture binding.
func execDef(name string, cap accessor.Capability, binding accessor.Binding) accessor.Definition {
	acc := table.Accessor{
		Role:    execRole,
		Path:    "flow.state",
		Keys:    []string{execKey},
		Timeout: execTimeout.String(),
	}
	if cap == accessor.CapWrite {
		acc.ReadBack = true
	}
	return accessor.Definition{
		Identity: accessor.Identity{
			Flow: "cmdflow", Name: name, Capability: cap,
		},
		Accessor: acc,
		Binding:  binding,
	}
}

func execExecutor(defs ...accessor.Definition) *accessor.Executor {
	reg := accessor.Registry{
		Flow:        "cmdflow",
		Definitions: defs,
		OwnedTags:   []string{execKey},
	}
	return accessor.NewExecutor(reg, accessor.Artifacts{
		execRole: {Role: execRole, Path: "flow.state"},
	})
}

func execCtx(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

// --- C4: the typed carrier ------------------------------------------------

// REQ-58: "`ExecError` is `struct { Detail string; Err error }` with
// `Error() string` and `Unwrap() error`: it WRAPS the offending `os/exec`
// error rather than flattening it, so `errors.Is(err, exec.ErrNotFound)`
// survives the trip to the refusal site and the argv0-not-found case stays
// distinguishable from a non-zero exit. `Detail` is the 4 KiB stderr tail,
// not `Err.Error()`"
// BOUNDARY
func TestReq58_ExecErrorWrapsRatherThanFlattensTheOffendingError(t *testing.T) {
	ee := &accessor.ExecError{
		Detail: "fatal: cannot change to '/nope': No such file or directory\n",
		Err:    exec.ErrNotFound,
	}

	if !errors.Is(ee, exec.ErrNotFound) {
		t.Error("errors.Is(ExecError, exec.ErrNotFound) is false; the wrap must " +
			"survive so a spawn failure stays distinguishable from a non-zero exit")
	}
	if got := ee.Unwrap(); !errors.Is(got, exec.ErrNotFound) {
		t.Errorf("Unwrap() = %v; want the wrapped exec.ErrNotFound", got)
	}
	if ee.Detail == ee.Err.Error() {
		t.Error("Detail equals Err.Error(); Detail is the stderr TAIL, not the " +
			"error's own text")
	}
	// A non-zero exit carries no ErrNotFound — the discriminating sibling.
	other := &accessor.ExecError{Detail: "boom", Err: errors.New("exit status 128")}
	if errors.Is(other, exec.ErrNotFound) {
		t.Error("a non-spawn ExecError matches exec.ErrNotFound; the two cases " +
			"must stay distinguishable")
	}

	// The clause says the wrap survives "the trip to the REFUSAL SITE" —
	// the shape alone proves nothing until the executor carries it. A
	// spawn failure (no exit code at all, fixture E1h) and a non-zero exit
	// must arrive as two distinguishable refusals, both carrying a tail.
	spawn := &execFailBinding{
		cap: accessor.CapRead,
		err: &accessor.ExecError{Detail: "", Err: exec.ErrNotFound},
	}
	nonZero := &execFailBinding{
		cap: accessor.CapRead,
		err: &accessor.ExecError{Detail: "exit tail\n", Err: errors.New("exit status 128")},
	}

	a := execExecutor(execDef(execReader, accessor.CapRead, spawn)).
		Read(execCtx(t), execReader)
	b := execExecutor(execDef(execReader, accessor.CapRead, nonZero)).
		Read(execCtx(t), execReader)

	if a.Refusal == nil || b.Refusal == nil {
		t.Fatal("both invocation failures must refuse")
	}
	if b.Refusal.Detail != "exit tail\n" {
		t.Errorf("the non-zero-exit refusal Detail = %q; want the stderr tail — "+
			"the wrap must survive the trip to the refusal site",
			b.Refusal.Detail)
	}
	if a.Refusal.Detail != "" {
		t.Errorf("the spawn-failure refusal Detail = %q; a spawn failure "+
			"produced no stderr and Detail is the TAIL, not Err.Error()",
			a.Refusal.Detail)
	}
}

// REQ-55: "detail:   the stderr tail reaches the refusal as a typed
// `*accessor.ExecError` (Detail string) returned by the command binding
// through the existing `error` return; `executor.go::refusalOf` gains an
// error parameter and `errors.As`-es it."
// REQ-50: "a NEW `Detail string` field on `accessor.Refusal` carries the
// last 4 KiB of stderr"
// REQ-66: "spawn failure, a malformed or oversized stdout envelope, and —
// by default — a non-zero exit are `ClassExecutionFailure`, whose `Detail`
// carries the bounded stderr tail."
// HAPPY PATH
func TestReq55_TheStderrTailReachesTheGateRefusalDetail(t *testing.T) {
	const tail = "fatal: cannot change to '/nope': No such file or directory\n"
	b := &execFailBinding{
		cap: accessor.CapGate,
		err: &accessor.ExecError{Detail: tail, Err: errors.New("exit status 128")},
	}
	e := execExecutor(execDef(execGate, accessor.CapGate, b))

	got := e.Gate(execCtx(t), execGate)

	if got.Refusal == nil {
		t.Fatal("a gate whose binding returned an ExecError did not refuse")
	}
	if got.Refusal.Class != accessor.ClassExecutionFailure {
		t.Errorf("class = %q; want %q", got.Refusal.Class, accessor.ClassExecutionFailure)
	}
	if got.Refusal.Detail != tail {
		t.Errorf("Detail = %q; want the stderr tail %q — the typed error is the "+
			"carrier and refusalOf must errors.As it", got.Refusal.Detail, tail)
	}
	if got.Refusal.Reason != "" {
		t.Errorf("Reason = %q; 0004:C7 reserves Reason for a gate DENY and the "+
			"tail must not be routed through it (REQ-63)", got.Refusal.Reason)
	}
}

// REQ-56: "`executor.go::refusalWithKeys` wraps it and is what
// `Executor.Read` routes every read refusal through, so the error parameter
// threads through BOTH or the read path — the capability that binds
// directly — silently drops its tail."
// REQ-64: "`readOutcome` drops the error today ... so it gains an `err`
// field for the read path to reach `refusalWithKeys`"
// ADVERSARIAL
//
// This is the arm the critique lens opened. A build that threads the
// parameter through `refusalOf` alone passes the gate and write tests above
// and fails HERE — which is exactly why it is a separate test.
func TestReq56_TheStderrTailReachesTheReadRefusalDetailToo(t *testing.T) {
	const tail = "test: --bogus: unexpected operator\n"
	b := &execFailBinding{
		cap: accessor.CapRead,
		err: &accessor.ExecError{Detail: tail, Err: errors.New("exit status 2")},
	}
	e := execExecutor(execDef(execReader, accessor.CapRead, b))

	got := e.Read(execCtx(t), execReader)

	if got.Refusal == nil {
		t.Fatal("a read whose binding returned an ExecError did not refuse")
	}
	if got.Refusal.Class != accessor.ClassExecutionFailure {
		t.Errorf("class = %q; want %q", got.Refusal.Class, accessor.ClassExecutionFailure)
	}
	if got.Refusal.Detail != tail {
		t.Errorf("Detail = %q; want the stderr tail %q — Executor.Read routes "+
			"every read refusal through refusalWithKeys, so the parameter must "+
			"thread through BOTH constructors", got.Refusal.Detail, tail)
	}
}

// REQ-55 (write arm): the same carrier on the write path.
// HAPPY PATH
func TestReq55_TheStderrTailReachesTheWriteRefusalDetail(t *testing.T) {
	const tail = "wrapper: git config failed\n"
	w := &execFailBinding{
		cap: accessor.CapWrite,
		err: &accessor.ExecError{Detail: tail, Err: errors.New("exit status 1")},
	}
	e := execExecutor(execDef(execWriter, accessor.CapWrite, w))

	got := e.Write(execCtx(t), execWriter, resolve.Plan{
		Writes: []resolve.Tag{{Key: execKey, Value: "final"}},
	})

	if got.Refusal == nil {
		t.Fatal("a write whose binding returned an ExecError did not refuse")
	}
	if got.Refusal.Class != accessor.ClassExecutionFailure {
		t.Errorf("class = %q; want %q", got.Refusal.Class, accessor.ClassExecutionFailure)
	}
	if got.Refusal.Detail != tail {
		t.Errorf("Detail = %q; want the stderr tail %q", got.Refusal.Detail, tail)
	}
}

// REQ-57: "`Detail` is empty on every refusal carrying no error (timeout,
// read-back, gate-off) and set only where a binding returned one."
// ADVERSARIAL
//
// The mutant this kills is a build that stuffs SOMETHING into Detail at
// every refusal site — a rendered class name, the accessor id — which would
// pass every "Detail is non-empty" assertion above while making the field
// meaningless as a diagnosis channel.
func TestReq57_DetailIsEmptyOnEveryRefusalCarryingNoInvocationError(t *testing.T) {
	// Positive precondition. "Detail is empty here" is satisfied vacuously
	// while nothing ever populates it, so the discriminating half runs
	// first: a refusal that DOES carry an invocation error must carry the
	// tail, or every arm below proves nothing.
	t.Run("a refusal carrying an invocation error DOES carry the tail", func(t *testing.T) {
		b := &execFailBinding{
			cap: accessor.CapGate,
			err: &accessor.ExecError{Detail: "tail\n", Err: errors.New("exit status 1")},
		}
		e := execExecutor(execDef(execGate, accessor.CapGate, b))

		got := e.Gate(execCtx(t), execGate)

		if got.Refusal == nil {
			t.Fatal("the gate did not refuse")
		}
		if got.Refusal.Detail != "tail\n" {
			t.Fatalf("Detail = %q; want the tail — without this the "+
				"empty-elsewhere arms are vacuous", got.Refusal.Detail)
		}
	})

	t.Run("an untyped binding error carries no tail", func(t *testing.T) {
		b := &execFailBinding{cap: accessor.CapRead, err: errors.New("plain failure")}
		e := execExecutor(execDef(execReader, accessor.CapRead, b))

		got := e.Read(execCtx(t), execReader)

		if got.Refusal == nil {
			t.Fatal("the read did not refuse")
		}
		if got.Refusal.Detail != "" {
			t.Errorf("Detail = %q; a path-backed binding returns an UNTYPED "+
				"error and must carry no tail (REQ-69)", got.Refusal.Detail)
		}
	})

	t.Run("an unknown accessor carries no tail", func(t *testing.T) {
		// Pre-selection: no binding is chosen, so no invocation error can
		// exist. Structural, not a call site remembering to leave it blank.
		e := execExecutor()

		got := e.Read(execCtx(t), "nosuchreader")

		if got.Refusal == nil {
			t.Fatal("an unbound accessor did not refuse")
		}
		if got.Refusal.Class != accessor.ClassUnknownAccessor {
			t.Fatalf("class = %q; want %q", got.Refusal.Class, accessor.ClassUnknownAccessor)
		}
		if got.Refusal.Detail != "" {
			t.Errorf("Detail = %q; a pre-selection refusal can carry no "+
				"invocation error", got.Refusal.Detail)
		}
	})

	t.Run("a read-back-incomplete refusal carries no tail", func(t *testing.T) {
		// A successful write with no reader for the role: the refusal is
		// read-back, which C4 names among the no-error refusals.
		w := &execFailBinding{cap: accessor.CapWrite}
		e := execExecutor(execDef(execWriter, accessor.CapWrite, w))

		got := e.Write(execCtx(t), execWriter, resolve.Plan{
			Writes: []resolve.Tag{{Key: execKey, Value: "final"}},
		})

		if got.Refusal == nil {
			t.Fatal("a write with no reader for its role did not refuse")
		}
		if got.Refusal.Class != accessor.ClassReadBackIncomplete {
			t.Fatalf("class = %q; want %q", got.Refusal.Class,
				accessor.ClassReadBackIncomplete)
		}
		if got.Refusal.Detail != "" {
			t.Errorf("Detail = %q; read-back refusals carry no invocation error",
				got.Refusal.Detail)
		}
		if !got.Refusal.Applied() {
			t.Error("Applied() is false; the write command already ran (0004:C14)")
		}
	})
}

// REQ-62: "Command bindings run under RDR 0004's executor with its refusal
// *classes* unchanged and one additive type change — a `Detail string`
// field on `accessor.Refusal`"
// REQ-63: "`Reason` is not reused because `0004:C7` reserves it for a gate
// deny."
// BOUNDARY — NEGATIVE REQ.
func TestReq62_TheRefusalClassSetIsUnchangedByThisRDR(t *testing.T) {
	want := []accessor.RefusalClass{
		accessor.ClassUnknownAccessor,
		accessor.ClassCapabilityMismatch,
		accessor.ClassTimeout,
		accessor.ClassExecutionFailure,
		accessor.ClassIncompleteRead,
		accessor.ClassGateIndeterminate,
		accessor.ClassReadBackMismatch,
		accessor.ClassReadBackIncomplete,
	}

	got := accessor.RefusalClasses()
	if len(got) != len(want) {
		t.Fatalf("RefusalClasses() = %#v; this RDR adds NO class and the set "+
			"stays %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("RefusalClasses()[%d] = %q; want %q", i, got[i], want[i])
		}
	}

	// The additive half. A struct literal proves only that the field
	// compiles; the clause is that the ADDITION is what carries the tail
	// while the class set stays put, so the oracle is a real refusal.
	b := &execFailBinding{
		cap: accessor.CapGate,
		err: &accessor.ExecError{Detail: "tail", Err: errors.New("exit status 1")},
	}
	r := execExecutor(execDef(execGate, accessor.CapGate, b)).
		Gate(execCtx(t), execGate).Refusal
	if r == nil {
		t.Fatal("the gate did not refuse")
	}
	if r.Detail != "tail" {
		t.Errorf("Refusal.Detail = %q; want %q — the additive field is what "+
			"carries the diagnosis the class set does not grow to hold",
			r.Detail, "tail")
	}
	if r.Class != accessor.ClassExecutionFailure {
		t.Errorf("class = %q; the addition must not mint a new class",
			r.Class)
	}
}

// REQ-69: "Path-backed bindings return untyped errors and are unaffected,
// so every production `NewExecutor` caller and every existing binding
// behave unchanged"
// DOMAIN EDGE — NEGATIVE REQ.
//
// The oracle is a refusal-by-refusal comparison of the SAME failing read
// under an untyped error against the typed one: the class, the diagnosis
// tuple, and the key set must be identical, and only Detail may differ.
func TestReq69_APathBackedBindingBehavesIdenticallyExceptForDetail(t *testing.T) {
	plain := &execFailBinding{cap: accessor.CapRead, err: errors.New("boom")}
	typed := &execFailBinding{
		cap: accessor.CapRead,
		err: &accessor.ExecError{Detail: "tail", Err: errors.New("boom")},
	}

	a := execExecutor(execDef(execReader, accessor.CapRead, plain)).
		Read(execCtx(t), execReader)
	b := execExecutor(execDef(execReader, accessor.CapRead, typed)).
		Read(execCtx(t), execReader)

	if a.Refusal == nil || b.Refusal == nil {
		t.Fatal("both reads must refuse")
	}
	if a.Refusal.Class != b.Refusal.Class {
		t.Errorf("class differs: untyped %q vs typed %q", a.Refusal.Class, b.Refusal.Class)
	}
	if a.Refusal.Accessor != b.Refusal.Accessor ||
		a.Refusal.Capability != b.Refusal.Capability ||
		a.Refusal.Role != b.Refusal.Role ||
		a.Refusal.Timeout != b.Refusal.Timeout {
		t.Errorf("the diagnosis tuple differs: %+v vs %+v", a.Refusal, b.Refusal)
	}
	if a.Refusal.Detail != "" {
		t.Errorf("the untyped path carries Detail %q; it must stay empty",
			a.Refusal.Detail)
	}
	if b.Refusal.Detail != "tail" {
		t.Errorf("the typed path carries Detail %q; want %q", b.Refusal.Detail, "tail")
	}
}

// REQ-68: "A write entry carries `read_back = true` ..., its success is
// never taken from exit status, and it is verified through the role's
// declared reader — the re-read runs under its own bounded timeout
// (`0004:C15`), never the write's residue"
// REQ-38: "A write command's success is never taken from its exit status
// alone — verification is read-back (`0004:C12`)."
// REQ-59: "write:    read_back required; verified only through the role's
// reader, never by exit status"
// ADVERSARIAL
//
// The discriminating fixture: a write binding that returns NO error (its
// child "exited zero") over an artifact the reader reports unchanged. Exit
// status says success; read-back says mismatch. The refusal is the oracle.
func TestReq38_AWriteThatExitsZeroWithoutApplyingStillRefuses(t *testing.T) {
	w := &execFailBinding{cap: accessor.CapWrite} // no error: "exit 0"
	// The reader reports the PRE-write value, as `git config` does after
	// the R3b spike case where the tool ignored stdin.
	r := &staticReadBinding{values: map[string]string{execKey: "draft"}}

	e := execExecutor(
		execDef(execWriter, accessor.CapWrite, w),
		execDef(execReader, accessor.CapRead, r),
	)

	got := e.Write(execCtx(t), execWriter, resolve.Plan{
		Writes: []resolve.Tag{{Key: execKey, Value: "final"}},
	})

	if got.Refusal == nil {
		t.Fatal("a write whose tool exited zero without applying reported " +
			"success; exit status decides no write (0004:C12)")
	}
	if got.Refusal.Class != accessor.ClassReadBackMismatch {
		t.Errorf("class = %q; want %q — read-back observed the artifact and "+
			"found it wrong", got.Refusal.Class, accessor.ClassReadBackMismatch)
	}
	if len(got.Written) != 0 {
		t.Errorf("Written = %#v on a refusal; the disjunction is exclusive", got.Written)
	}
	if r.reads == 0 {
		t.Error("the read-back reader never ran; a write verified by anything " +
			"other than the role's reader is not verified at all")
	}

	// The command-write sibling: the SAME plan, applied by a binding whose
	// child exited NON-zero, must reach the same seam and carry its tail —
	// which is what makes "exit status decides no write" a statement about
	// commands rather than about this fixture.
	failing := &execFailBinding{
		cap: accessor.CapWrite,
		err: &accessor.ExecError{Detail: "wrapper: refused\n", Err: errors.New("exit status 1")},
	}
	fr := &staticReadBinding{values: map[string]string{execKey: "draft"}}
	fe := execExecutor(
		execDef(execWriter, accessor.CapWrite, failing),
		execDef(execReader, accessor.CapRead, fr),
	)
	fgot := fe.Write(execCtx(t), execWriter, resolve.Plan{
		Writes: []resolve.Tag{{Key: execKey, Value: "final"}},
	})
	if fgot.Refusal == nil {
		t.Fatal("the failing command write did not refuse")
	}
	if fgot.Refusal.Detail != "wrapper: refused\n" {
		t.Errorf("Detail = %q; want the command's stderr tail", fgot.Refusal.Detail)
	}
}

// staticReadBinding answers every requested key from a fixed map; a key the
// map lacks is established-absent, exactly as a command reader's
// `exit_absent` arm reports it.
type staticReadBinding struct {
	values map[string]string
	reads  int
}

func (b *staticReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *staticReadBinding) Read(_ context.Context, _ accessor.Artifact, requested []string) (
	[]accessor.KeyValue, []string, error,
) {
	b.reads++
	out := make([]accessor.KeyValue, 0, len(requested))
	for _, k := range requested {
		v, held := b.values[k]
		out = append(out, accessor.KeyValue{Key: k, Value: v, Absent: !held})
	}
	return out, nil, nil
}

// --- C7: the seam this RDR must not change --------------------------------

// REQ-65: "There is exactly **one** slot to land in — ...
// `internal/cli/clierr/clierr.go::CLIError` has a single `Detail string`
// ... the applied-sense text **first** ..., the stderr tail appended after
// it."
// DOMAIN EDGE
//
// Asserted at THIS layer as the precondition the CLI composes from: a
// post-mutation refusal must carry BOTH the applied sense and (when a
// binding supplied one) the tail, as two separable values. The ordering
// itself is asserted at the CLI, where the single slot lives.
func TestReq65_APostMutationRefusalCarriesTheAppliedSenseAndTheTailSeparably(t *testing.T) {
	// A write that applied, then a read-back whose reader fails with a
	// typed error: the applied sense is on the write refusal and the tail
	// comes from the read-back invocation.
	w := &execFailBinding{cap: accessor.CapWrite}
	r := &execFailBinding{
		cap: accessor.CapRead,
		err: &accessor.ExecError{Detail: "bad config line 1\n", Err: errors.New("exit status 128")},
	}
	e := execExecutor(
		execDef(execWriter, accessor.CapWrite, w),
		execDef(execReader, accessor.CapRead, r),
	)

	got := e.Write(execCtx(t), execWriter, resolve.Plan{
		Writes: []resolve.Tag{{Key: execKey, Value: "final"}},
	})

	if got.Refusal == nil {
		t.Fatal("a write whose read-back failed did not refuse")
	}
	if got.Refusal.Class != accessor.ClassReadBackIncomplete {
		t.Errorf("class = %q; want %q — the corrupted artifact no longer "+
			"parses for the reader (S5)", got.Refusal.Class,
			accessor.ClassReadBackIncomplete)
	}
	if !got.Refusal.Applied() {
		t.Error("Applied() is false; the write command already ran (0004:C14)")
	}
	// The read-back's own tail must reach the refusal, or the CLI has
	// nothing to append after the applied-sense text.
	if !strings.Contains(got.Refusal.Detail, "bad config line 1") {
		t.Errorf("Detail = %q; want the read-back invocation's stderr tail",
			got.Refusal.Detail)
	}
}
