// Package accessor is RDR 0004's accessor execution safety boundary: the
// declared-capability executor that invokes read, gate, and write
// bindings over caller-supplied artifact roles.
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
	"slices"
	"time"

	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
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

var capabilities = []Capability{CapRead, CapGate, CapWrite}

// Capabilities returns the closed three-member capability vocabulary.
func Capabilities() []Capability { return slices.Clone(capabilities) }

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

var refusalClasses = []RefusalClass{
	ClassUnknownAccessor,
	ClassCapabilityMismatch,
	ClassTimeout,
	ClassExecutionFailure,
	ClassIncompleteRead,
	ClassGateIndeterminate,
	ClassReadBackMismatch,
	ClassReadBackIncomplete,
}

// RefusalClasses returns the closed accessor refusal-class set.
func RefusalClasses() []RefusalClass { return slices.Clone(refusalClasses) }

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

var validationCodes = []ValidationCode{
	CodeMissingAccessor,
	CodeMultiplyBoundAccessor,
	CodeCapabilityMismatch,
	CodeMissingOrNonPositiveTimeout,
	CodeMissingWriteReadBack,
	CodeAmbientArtifactDiscovery,
	CodeWriteNonOwnedTag,
	CodeMissingRequestedKeySet,
}

// ValidationCodes returns the closed eight-member validation-code set.
func ValidationCodes() []ValidationCode { return slices.Clone(validationCodes) }

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
func (d Definition) RequestedKeys() []string { return slices.Clone(d.Accessor.Keys) }

// OwnedKeys returns the writer's owned key set — the keys it may write or
// clear, taken from `keys` (LBD, Definition shape).
func (d Definition) OwnedKeys() []string { return slices.Clone(d.Accessor.Keys) }

// timeout parses the definition's declared timeout metadata. A missing,
// unparseable, or non-positive value is rejected by Validate before
// execution (`0004:C16`); at runtime it yields ok = false.
func (d Definition) timeout() (time.Duration, bool) {
	if d.Accessor.Timeout == "" {
		return 0, false
	}
	v, err := time.ParseDuration(d.Accessor.Timeout)
	if err != nil || v <= 0 {
		return 0, false
	}
	return v, true
}

// Registry is the validated set of accessor definitions for one flow.
type Registry struct {
	Flow        string
	Definitions []Definition
	// OwnedTags names the model's owned tag keys, so validation can
	// refuse a writer naming a non-owned tag.
	OwnedTags []string

	// AllowCommands is the `--allow-commands` gate (0025:C6), carried
	// here as STATE (`0028:C1.3` SITE:, A10).
	//
	// It is a FIELD and never a read across the `cmdbind` seam: `cmdbind`
	// imports `accessor`, so the reverse is an import cycle, and
	// `AllowCommands` on `cmdbind.Config` stays where it is. It is set at
	// `flowbind.Registry`, the single production construction site, which
	// already takes the gate for 0025:C6.
	//
	// The executor needs it to answer C1.3's read-back pre-check: when a
	// write's role reader is command-backed and the gate is off, the
	// write must refuse BEFORE mutation, because a forgotten flag must
	// not produce `read_back_incomplete` for a write that ran no command.
	AllowCommands bool
}

// Lookup selects a binding by capability table and id: an invocation that
// needs capability X selects only from `[X.<id>]`, never from a same-id
// entry under another table (LBD, Selection / predicate).
func (reg Registry) Lookup(name string, capability Capability) (Definition, bool) {
	for _, d := range reg.Definitions {
		if d.Identity.Name == name && d.Identity.Capability == capability {
			return d, true
		}
	}
	return Definition{}, false
}

// bound reports whether any capability table binds the id at all. It is
// what separates "not bound" (`unknown_accessor`) from "bound under
// another capability" (`capability_mismatch`) — two different defects.
func (reg Registry) bound(name string) bool {
	for _, d := range reg.Definitions {
		if d.Identity.Name == name {
			return true
		}
	}
	return false
}

// readerFor selects the read definition serving one artifact role. The
// read-back re-read goes through it, so the verification is subject to
// read completeness and can fail independently of the write (`0004:C13`,
// CA A10).
func (reg Registry) readerFor(role string) (Definition, bool) {
	for _, d := range reg.Definitions {
		if d.Identity.Capability == CapRead && d.Accessor.Role == role {
			return d, true
		}
	}
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

	// Context carries the invocation's bound tag values (RDR 0028 A1 —
	// the ONE seam extension this record owns).
	//
	// It rides the ARTIFACT rather than a second `Apply` argument because
	// the same `art` value reaches both carriers: an `edit` writer's
	// `{tag.<key>}` anchors and a command reader's `{tag.<key>}` argv are
	// the same channel, and a shared artifact's reader could not say
	// which row it read without it (`0028:C1.6` why here:). A separate
	// write-side argument would serve the anchor half and leave the read
	// half unaddressed.
	//
	// It is nil for an invocation that bound no tags, and a binding that
	// references none never reads it: the scan is per-USE-SITE, so a tag
	// bound here but named by no anchor or argv element of an entry is
	// never inspected on that entry's behalf.
	Context map[string]string
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

var verdicts = []Verdict{VerdictAllow, VerdictDeny, VerdictIndeterminate}

// Verdicts returns the closed three-member gate verdict vocabulary.
func Verdicts() []Verdict { return slices.Clone(verdicts) }

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

	// Detail carries the bounded stderr tail of a command-backed
	// invocation (`0025:C4`). It is ADDITIVE: RDR 0004 pins the refusal
	// CLASS set, not this struct's field set, and `Reason` is not reused
	// because `0004:C7` reserves it for a gate deny. It is EMPTY on every
	// refusal carrying no invocation error — timeout, read-back, gate-off
	// — and set only where a binding returned an *ExecError.
	Detail string

	// applied records the post-mutation sense. It is unexported because
	// only the write path — which knows whether the command already ran —
	// may set it.
	applied bool
}

// Applied reports the post-mutation sense: whether the write command
// already ran when this refusal was minted. It is TRUE for
// read_back_incomplete and for a timeout raised after the write command
// ran, and those MUST be understood as "the mutation may have been
// applied and was not verified" — never as a write that did not occur
// (`0004:C14`).
func (r Refusal) Applied() bool { return r.applied }

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
func (r ReadResult) OwnedSnapshot() []resolve.Tag {
	if r.Refusal != nil {
		return nil
	}
	out := make([]resolve.Tag, 0, len(r.Values))
	for _, v := range r.Values {
		if v.Absent {
			// Omission, never a placeholder: a sentinel would make
			// `TagSet.has` true and silently retire the refusal.
			continue
		}
		out = append(out, resolve.Tag{Key: v.Key, Value: v.Value})
	}
	slices.SortFunc(out, func(a, b resolve.Tag) int {
		switch {
		case a.Key < b.Key:
			return -1
		case a.Key > b.Key:
			return 1
		default:
			return 0
		}
	})
	return out
}

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
	// Written echoes the planned owned tags the read-back verified as
	// HELD. A key the plan cleared is OMITTED — read-back verified it
	// absent, and a verified removal is not a written value (`0004:C11`,
	// REQ-32).
	Written []resolve.Tag
	// Refusal is non-nil exactly when the write refused, INCLUDING when
	// the command succeeded but read-back failed (`0004:C12`).
	Refusal *Refusal
}

// Refused reports whether the write took the refusal branch.
func (r WriteResult) Refused() bool { return r.Refusal != nil }

// --- the reserved clear sentinel -----------------------------------------

// ClearSentinel is the reserved tag value a planned removal takes (JDR
// 0001 §D5). It is RDR 0002's constant; this package names it so an
// executor test never re-spells the literal.
const ClearSentinel = table.ClearSentinel

// IsClear reports whether a planned tag value is the reserved removal
// sentinel.
func IsClear(value string) bool { return value == ClearSentinel }

// --- RDR 0025: the invocation-error carrier -------------------------------

// ExecError is the typed carrier a command-backed binding returns through
// the binding interfaces' existing bare `error` slot, so a stderr tail
// reaches the refusal without an interface change (`0025:C4`, A11).
//
// It WRAPS the offending `os/exec` error rather than flattening it, so
// `errors.Is(err, exec.ErrNotFound)` survives the trip to the refusal
// site and the argv0-not-found case stays distinguishable from a non-zero
// exit. Detail is the bounded stderr tail, NOT `Err.Error()`.
type ExecError struct {
	// Detail is the last 4 KiB of the child's stderr (`0025:C4`).
	Detail string
	// Err is the offending error, wrapped and never flattened.
	Err error
}

// Error renders the typed error.
func (e *ExecError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Detail
}

// Unwrap exposes the wrapped error so `errors.Is`/`errors.As` reach it.
func (e *ExecError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
