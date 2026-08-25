// Package accessor is RDR 0004's accessor execution safety boundary: the
// declared-capability executor that invokes read, gate, and write
// bindings over caller-supplied artifact roles.
//
// SKELETON ONLY (Stage 8, Phase 1). Every declaration below is a
// signature returning a zero value so the RDR 0004 test suite COMPILES
// and each test fails on its own assertion — naming the REQ its header
// quotes — rather than the whole package failing with one
// undefined-symbol error that attributes to no clause. Phase 2 replaces
// each body. Nothing here implements a check.
//
// The boundary this package fixes (`0004:AP`): the executor is NOT a
// shell runner and NOT a state-machine action callback surface. It
// invokes typed bindings selected by accessor name and capability,
// enforces a per-accessor timeout, converts execution errors into stable
// refusal classes, and never prints directly (`0004:C17`). RDR 0002 owns
// the TOML carrier and its load-time validations; this package consumes
// `table.Accessor` values and owns EXECUTION.
package accessor

import (
	"context"
	"time"

	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- capability ----------------------------------------------------------

// Capability is the declared authority class of one accessor. The set is
// closed at three: read, gate, and write (`0004:C1`).
//
// The canonical names are "read accessor", "gate accessor", and "write
// accessor". "hook" and "action" are rejected names because they imply
// arbitrary transition callbacks (LBD, Naming).
type Capability string

const (
	// CapRead reads typed tag values from a caller-supplied artifact role.
	CapRead Capability = "read"
	// CapGate answers allow, deny, or indeterminate.
	CapGate Capability = "gate"
	// CapWrite applies planned owned-tag writes and verifies by read-back.
	CapWrite Capability = "write"
)

// Capabilities returns the closed three-member capability vocabulary.
func Capabilities() []Capability { return nil }

// --- refusal classes -----------------------------------------------------

// RefusalClass is the accessor-owned refusal discriminator. This set is
// DISJOINT from the kernel's closed five-kind
// `resolve.RefusalKinds` set, which this RDR does not extend (`0004:FM`).
type RefusalClass string

const (
	// ClassUnknownAccessor: the invoked accessor name is not bound.
	ClassUnknownAccessor RefusalClass = "unknown_accessor"
	// ClassCapabilityMismatch: the accessor was used off-capability.
	ClassCapabilityMismatch RefusalClass = "capability_mismatch"
	// ClassTimeout: the invocation exceeded its bounded timeout. Its own
	// class, distinct from execution failure and read-back mismatch
	// (`0004:C15`).
	ClassTimeout RefusalClass = "timeout"
	// ClassExecutionFailure: the binding failed, including an unsupplied
	// or unreadable artifact role — NOT a separate class (`0004:FM`).
	ClassExecutionFailure RefusalClass = "execution_failure"
	// ClassIncompleteRead: a requested key could not be read (`0004:C7`).
	ClassIncompleteRead RefusalClass = "incomplete_read"
	// ClassGateIndeterminate: the gate could not decide (`0004:C9`).
	ClassGateIndeterminate RefusalClass = "gate_indeterminate"
	// ClassReadBackMismatch: the re-read asserts the artifact is WRONG.
	ClassReadBackMismatch RefusalClass = "read_back_mismatch"
	// ClassReadBackIncomplete: the verification DID NOT RUN (`0004:C13`).
	ClassReadBackIncomplete RefusalClass = "read_back_incomplete"
)

// RefusalClasses returns the closed accessor refusal-class set.
func RefusalClasses() []RefusalClass { return nil }

// --- validation codes ----------------------------------------------------

// ValidationCode names one definition-validation defect. Validation runs
// BEFORE resolution and rejects the eight shapes (`0004:TD`, TS 1).
type ValidationCode string

const (
	// CodeMissingAccessor: a referenced identity has no binding.
	CodeMissingAccessor ValidationCode = "missing_accessor"
	// CodeMultiplyBoundAccessor: one identity triple has two bindings.
	CodeMultiplyBoundAccessor ValidationCode = "multiply_bound_accessor"
	// CodeCapabilityMismatch: the reference wants a capability the
	// binding does not declare.
	CodeCapabilityMismatch ValidationCode = "capability_mismatch"
	// CodeMissingOrNonPositiveTimeout: timeout metadata absent or <= 0.
	CodeMissingOrNonPositiveTimeout ValidationCode = "missing_or_non_positive_timeout"
	// CodeMissingWriteReadBack: a write entry without read_back = true.
	CodeMissingWriteReadBack ValidationCode = "missing_write_read_back"
	// CodeAmbientArtifactDiscovery: the definition discovers artifacts
	// from ambient process state (`0004:C3`).
	CodeAmbientArtifactDiscovery ValidationCode = "ambient_artifact_discovery"
	// CodeWriteNonOwnedTag: a writer names an observed or recognized tag.
	CodeWriteNonOwnedTag ValidationCode = "write_non_owned_tag"
	// CodeMissingRequestedKeySet is the EIGHTH arm: a reader with a
	// missing or empty requested key set (`0004:C5`; A9, unwitnessed).
	CodeMissingRequestedKeySet ValidationCode = "missing_requested_key_set"
)

// ValidationCodes returns the closed eight-member validation-code set.
func ValidationCodes() []ValidationCode { return nil }

// Finding is one validation defect: its own named code plus the identity
// it attributes to (ORA 1 — each arm asserts its OWN named code, never
// merely that validation returned non-empty).
type Finding struct {
	Code       ValidationCode
	Accessor   string
	Capability Capability
	// Detail names the offending key or role, for diagnosis.
	Detail string
}

// --- definitions and the registry ----------------------------------------

// Identity is the accessor identity triple `(flow id, accessor name,
// capability)`. Capability is PART of the identity, so the same id may
// appear in both `[read.x]` and `[write.x]` — two identities, not a
// rebinding (LBD, Identity; JDR 0001 §D7(ii)).
type Identity struct {
	Flow       string
	Name       string
	Capability Capability
}

// Definition is one validated accessor definition. It consumes RDR
// 0002's capability-table entry (`table.Accessor`) rather than spelling a
// second carrier (LBD, Wire / byte format; REQ-8).
type Definition struct {
	Identity Identity
	// Accessor is RDR 0002's carried entry: Role, Path, Keys, Timeout,
	// ReadBack.
	Accessor table.Accessor
	// Binding is the typed binding this identity invokes. It is never a
	// shell string and never a host callback over transition state.
	Binding Binding
	// AmbientDiscovery marks a definition that would discover its
	// artifact from ambient process state. It exists so validation can
	// REJECT it (`0004:C3`); no valid definition carries it true.
	AmbientDiscovery bool
}

// RequestedKeys returns the reader's validated requested key set, taken
// from the definition's `keys` metadata and NEVER derived from the keys a
// read actually resolved (`0004:C5`, LBD Requested key set).
func (d Definition) RequestedKeys() []string { return nil }

// OwnedKeys returns the writer's owned key set — the keys it may write or
// clear, taken from `keys` (LBD, Definition shape).
func (d Definition) OwnedKeys() []string { return nil }

// Registry is the validated set of accessor definitions for one flow.
type Registry struct {
	Flow        string
	Definitions []Definition
	// OwnedTags names the model's owned tag keys, so validation can
	// refuse a writer naming a non-owned tag.
	OwnedTags []string
}

// Validate runs the eight validation arms over the registry and the
// identities the model references. It returns every finding it can
// decide, each carrying its own named code. A valid registry returns the
// EMPTY set — the negative control that fails a validator which rejects
// everything (ORA 1).
func Validate(reg Registry, referenced []Identity) []Finding { return nil }

// Lookup selects a binding by capability table and id: an invocation that
// needs capability X selects only from `[X.<id>]`, never from a same-id
// entry under another table (LBD, Selection / predicate).
func (reg Registry) Lookup(name string, capability Capability) (Definition, bool) {
	return Definition{}, false
}

// --- the artifact seam ---------------------------------------------------

// Artifacts is the caller-supplied artifact role map. The executor MUST
// operate on these roles and MUST NOT discover authoritative artifacts
// from ambient process state (`0004:C3`).
type Artifacts map[string]Artifact

// Artifact is one caller-supplied artifact role handle. It is opaque to
// the executor: the binding, not the executor, knows how to read it.
type Artifact struct {
	Role string
	Path string
}

// --- bindings ------------------------------------------------------------

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

// --- verdicts ------------------------------------------------------------

// Verdict is a gate accessor's answer. Three-valued (`0004:C9`).
type Verdict string

const (
	// VerdictAllow: the gate permits.
	VerdictAllow Verdict = "allow"
	// VerdictDeny: the gate refuses, WITH A REASON. At this accessor
	// boundary deny is a TYPED RESULT, not a refusal class; it is a
	// refusal at the CLI (JDR 0001 §D9; req-list REQ-37, Q1).
	VerdictDeny Verdict = "deny"
	// VerdictIndeterminate: the gate could not decide. It is a
	// REFUSAL-class result, never a false allow and never a false deny.
	VerdictIndeterminate Verdict = "indeterminate"
)

// Verdicts returns the closed three-member gate verdict vocabulary.
func Verdicts() []Verdict { return nil }

// --- refusals ------------------------------------------------------------

// Refusal is one typed accessor refusal. It is a structured VALUE: the
// package returns it and never prints it (`0004:C17`).
type Refusal struct {
	// Class is the stable discriminator RDR 0005 later maps to a CLI
	// code. `exit code` is out of scope here (`0004:DISP`).
	Class RefusalClass
	// Accessor, Capability, Role, and Timeout are the diagnosis tuple
	// (`0004:FM`: "Diagnosis starts with the accessor identity,
	// capability, artifact role, timeout, and expected versus observed
	// tag values").
	Accessor   string
	Capability Capability
	Role       string
	Timeout    time.Duration
	// Keys names the requested keys the read could not read, on
	// incomplete_read and read_back_incomplete (`0004:C7`).
	Keys []string
	// Expected and Observed carry the read-back comparison, on
	// read_back_mismatch.
	Expected []resolve.Tag
	Observed []resolve.Tag
	// Reason carries a gate's deny reason where the caller mapped a deny
	// into a refusal. It is never set by this package for a deny.
	Reason string
}

// Applied reports the post-mutation sense: whether the write command
// already ran when this refusal was minted. It is TRUE for
// read_back_incomplete and for a timeout raised after the write command
// ran, and those MUST be understood as "the mutation may have been
// applied and was not verified" — never as a write that did not occur
// (`0004:C14`).
func (r Refusal) Applied() bool { return false }

// --- results -------------------------------------------------------------

// ReadResult is a read accessor's disposition: exactly one of Values or
// Refusal, never both and never neither. Completeness constrains WHICH
// branch the disjunction takes; it does not add a third branch
// (`0004:C4`, `0004:C6`, DESK "No CONTRADICTION row").
type ReadResult struct {
	// Values holds exactly the requested keys — none missing, none added
	// — on the success branch. A key the artifact genuinely lacks is
	// carried here with Absent = true.
	Values []KeyValue
	// Refusal is non-nil exactly when the read refused.
	Refusal *Refusal
}

// Refused reports whether the read took the refusal branch.
func (r ReadResult) Refused() bool { return r.Refusal != nil }

// OwnedSnapshot renders the values for the RESOLVER seam. A key the
// artifact genuinely does not carry is OMITTED — never carried as a
// placeholder value — so a required absent key resolves as
// `owned_state_unavailable` rather than matching a sentinel (`0004:C8`,
// MVV Scenario 6).
func (r ReadResult) OwnedSnapshot() []resolve.Tag { return nil }

// GateResult is a gate accessor's disposition. A deny is a typed result
// carrying a reason (Verdict + Reason); indeterminate, timeout, and
// execution failure are refusals.
type GateResult struct {
	Verdict Verdict
	Reason  string
	Refusal *Refusal
}

// Refused reports whether the gate took the refusal branch.
func (r GateResult) Refused() bool { return r.Refusal != nil }

// WriteResult is a write accessor's disposition after read-back
// verification.
type WriteResult struct {
	// Written echoes the planned owned tags the read-back verified.
	Written []resolve.Tag
	// Refusal is non-nil exactly when the write refused, INCLUDING when
	// the command succeeded but read-back failed (`0004:C12`).
	Refusal *Refusal
}

// Refused reports whether the write took the refusal branch.
func (r WriteResult) Refused() bool { return r.Refusal != nil }

// --- the executor --------------------------------------------------------

// Executor invokes validated accessor definitions against
// caller-supplied artifacts. It holds no ambient state and performs no
// scheduling: the caller decides WHEN a gate runs (JDR 0001 §D9 lands the
// gate SITE in RDR 0005; req-list Q2).
type Executor struct {
	Registry  Registry
	Artifacts Artifacts
}

// NewExecutor is THE executor constructor.
func NewExecutor(reg Registry, arts Artifacts) *Executor { return &Executor{} }

// Read invokes the named read accessor. The requested key set comes from
// the definition's validated metadata, never from the keys the read
// happened to resolve (`0004:C5`).
func (e *Executor) Read(ctx context.Context, name string) ReadResult {
	return ReadResult{}
}

// Gate invokes the named gate accessor. It REPORTS the verdict and never
// APPLIES a deny (JDR 0001 §D9; REQ-38).
func (e *Executor) Gate(ctx context.Context, name string) GateResult {
	return GateResult{}
}

// Write invokes the named write accessor for a SUCCESSFUL transition
// plan, then re-reads the same caller-supplied artifact role named by the
// write binding and verifies read-back (`0004:C12`).
//
// The plan's Writes are the ONLY tags applied; an escaped plan carries an
// empty write set and no write runs (deviations D1, RDR 0009
// `0009:1515-1527`).
func (e *Executor) Write(ctx context.Context, name string, plan resolve.Plan) WriteResult {
	return WriteResult{}
}

// --- replay disposition --------------------------------------------------

// Disposition is the recorded, replay-stable outcome of one invocation.
// It records artifact role, accessor name, capability, timeout, and the
// returned tag values — the five the model must record for accessor
// execution to be deterministic enough for resolver replay (CA A4).
type Disposition struct {
	Accessor   string
	Capability Capability
	Role       string
	Timeout    time.Duration
	// Refusal is the refusal class, empty on success.
	Refusal RefusalClass
	// Tags is the returned tag values, rendered deterministically.
	Tags []resolve.Tag
	// Verdict is the gate verdict, empty for non-gate invocations.
	Verdict Verdict
}

// ReadDisposition renders a read result as a replay-stable disposition.
func ReadDisposition(def Definition, r ReadResult) Disposition { return Disposition{} }

// GateDisposition renders a gate result as a replay-stable disposition.
func GateDisposition(def Definition, r GateResult) Disposition { return Disposition{} }

// WriteDisposition renders a write result as a replay-stable disposition.
func WriteDisposition(def Definition, r WriteResult) Disposition { return Disposition{} }

// --- the reserved clear sentinel -----------------------------------------

// ClearSentinel is the reserved tag value a planned removal takes (JDR
// 0001 §D5). It is RDR 0002's constant; this package names it so an
// executor test never re-spells the literal.
const ClearSentinel = table.ClearSentinel

// IsClear reports whether a planned tag value is the reserved removal
// sentinel.
func IsClear(value string) bool { return false }
