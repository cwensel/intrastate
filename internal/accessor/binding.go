package accessor

import (
	"context"
	"strconv"

	"github.com/cwensel/intrastate/internal/resolve"
)

// KeyValue is one key the binding resolved, with its readability
// disposition. `Absent` is a VALUE (the binding read the artifact
// successfully and the key was not there); unreadability is signalled by
// the binding refusing the key instead (LBD, Absent vs unreadable).
type KeyValue struct {
	Key   string
	Value string
	// Absent reports that the binding established the key is not carried.
	// This is the accessor layer's INTERNAL representation; what crosses
	// the seam omits the key entirely (`0004:C8`).
	Absent bool
}

// Binding is the typed invocation seam. It is the only way the executor
// reaches the outside world: no shell-out, no host callback over
// transition state (Alternatives 2, 3, 5 rejected).
//
// A read or read-back binding implements ReadBinding; a gate binding
// implements GateBinding; a write binding implements WriteBinding.
type Binding interface {
	// Capability reports the capability this binding serves.
	Capability() Capability
}

// ReadBinding resolves requested keys from a caller-supplied artifact.
type ReadBinding interface {
	Binding
	// Read resolves each requested key against art. It returns one
	// KeyValue per key it could read (present or established-absent) and
	// names in unreadable every requested key it could NOT read. An
	// error is an execution failure.
	Read(ctx context.Context, art Artifact, requested []string) (
		values []KeyValue, unreadable []string, err error)
}

// GateBinding answers allow, deny, or indeterminate.
type GateBinding interface {
	Binding
	// Gate decides. An error is an execution failure; Indeterminate is
	// NOT an error, it is the third verdict (`0004:C9`).
	Gate(ctx context.Context, art Artifact) (Verdict, string, error)
}

// WriteBinding applies planned owned-tag writes.
type WriteBinding interface {
	Binding
	// Apply performs the mutation. A `<clear>` planned value is a
	// REMOVAL, not an assignment of the literal (`0004:C11`).
	Apply(ctx context.Context, art Artifact, planned []resolve.Tag) error
	// Invocations reports how many times Apply ran, so a test can assert
	// the accessor layer performed no retry, undo, or re-derivation
	// (`0004:C14`, ORA 8).
	Invocations() int
}

// PriorBinding is the optional capability of a write binding whose
// effect depends on the value a key holds BEFORE the write (kata q14r).
// The `steps` carrier is the one implementation: a replace runs the
// prior value's `clear` arm before the planned value's `set` arm, so the
// binding cannot choose its argv from the plan alone.
//
// It EXTENDS WriteBinding rather than amending it, so every existing
// binding keeps the interface `0025:C7` quotes unchanged. The executor
// pre-reads PriorKeys through the role's reader BEFORE any mutation —
// an unreadable prior refuses with nothing applied — and then calls
// ApplyPrior in place of Apply. ApplyPrior counts as one Apply entry
// on Invocations.
type PriorBinding interface {
	WriteBinding
	// PriorKeys names the planned keys whose prior value the binding
	// needs. Empty means no pre-read runs.
	PriorKeys(planned []resolve.Tag) []string
	// ApplyPrior is Apply with the pre-read prior values of PriorKeys:
	// a present value, or `Absent` for a key the artifact does not hold.
	ApplyPrior(ctx context.Context, art Artifact, planned []resolve.Tag,
		prior []KeyValue) error
}

// PartialApplyError is the typed error a multi-step write binding wraps
// inside `ExecError.Err` when a step fails AFTER an earlier step of the
// same Apply already exited 0 (kata q14r). The mutation is then partly
// applied, so the executor reports the refusal in `0004:C14`'s
// applied-but-unverified sense, as it does for `HeldPipeError`. A failure
// on the FIRST spawned step is not this error: nothing ran, and the
// refusal keeps the safe-to-retry sense.
//
// It is declared here for `HeldPipeError`'s reason: `cmdbind` imports
// `accessor`, so the executor's `errors.As` site can name only a type
// this package owns.
type PartialApplyError struct {
	// Ran is how many steps exited 0 before the failing one.
	Ran int
	// Step names the failing step as `steps.<key>.<arm>.<value>[<i>]`.
	Step string
	// Err is the failing step's own error, wrapped and never flattened.
	Err error
}

// Error renders the partial-application reason.
func (e *PartialApplyError) Error() string {
	if e == nil {
		return ""
	}
	return "the step " + e.Step + " failed after " + strconv.Itoa(e.Ran) +
		" earlier step(s) of this write exited 0: " + e.Err.Error()
}

// Unwrap exposes the failing step's error so `errors.Is`/`errors.As`
// reach it.
func (e *PartialApplyError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
