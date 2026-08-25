package accessor

import (
	"context"

	"github.com/newcoinc/intrastate/internal/resolve"
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
