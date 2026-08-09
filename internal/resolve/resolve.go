// Package resolve is the resolver kernel: the pure decision boundary
// between a freshly recognized outcome and the next legal state.
//
// RDR 0001 locks this package's contract. The kernel is stateless and
// non-orchestrating: it consumes a transition table, an accessor-produced
// owned tag snapshot, caller-supplied observed tags, and one freshly
// recognized outcome tag, and it returns exactly one disposition — a
// transition plan or a typed refusal. It never prints output, inspects CLI
// flags, discovers ambient state, chooses artifacts, initiates work, or
// executes persistence.
//
// This file is the Phase 1 boundary skeleton: types and the entry-point
// signature only. It contains no resolution logic.
package resolve

// Tag is a single state fact in the tag-set model. Tags are compared by
// value; the kernel never parses their internal structure.
type Tag struct {
	Key   string
	Value string
}

// Provenance distinguishes where a tag entered the evaluation view. RDR
// 0001's Background requires the kernel to distinguish owned, observed,
// and freshly recognized tags; RDR 0004 forbids writing anything but
// owned tags.
type Provenance int

const (
	// ProvenanceOwned marks a tag read from a caller-provided artifact by
	// the accessor layer. Only owned tags may appear in plan writes.
	ProvenanceOwned Provenance = iota
	// ProvenanceObserved marks a non-owned context tag supplied by the
	// caller.
	ProvenanceObserved
	// ProvenanceRecognized marks the freshly recognized outcome tag.
	ProvenanceRecognized
)

// RefusalKind is the kernel-owned refusal discriminator. The set is
// closed: RDR 0001's Normative Contracts say the kind set "is exactly"
// the five constants below, and RDR 0005 maps each to a CLI error code
// without inspecting error strings.
type RefusalKind string

const (
	// KindNoMatch: no table edge matched the assembled tag-set.
	KindNoMatch RefusalKind = "no_match"
	// KindAmbiguousMatch: more than one edge matched after guards.
	KindAmbiguousMatch RefusalKind = "ambiguous_match"
	// KindOwnedStateUnavailable: evaluation required an owned tag absent
	// from the accessor-produced owned snapshot.
	KindOwnedStateUnavailable RefusalKind = "owned_state_unavailable"
	// KindGuardUnevaluable: the guard seam reported a predicate it could
	// not decide.
	KindGuardUnevaluable RefusalKind = "guard_unevaluable"
	// KindUnmodeledOutcome: the recognized outcome is outside the table's
	// declared recognized-outcome alphabet.
	KindUnmodeledOutcome RefusalKind = "unmodeled_outcome"
)

// RefusalKinds returns the closed kernel-owned refusal kind set.
func RefusalKinds() []RefusalKind { return nil }

// GuardResult is a guard seam verdict. RDR 0003 owns predicate shape; the
// kernel only consumes decided/undecided verdicts.
type GuardResult int

const (
	// GuardFalse: the predicate is decided and does not hold.
	GuardFalse GuardResult = iota
	// GuardTrue: the predicate is decided and holds.
	GuardTrue
	// GuardUnevaluable: the seam cannot decide the predicate. The kernel
	// answers with KindGuardUnevaluable.
	GuardUnevaluable
)

// GuardEvaluator is the delegated guard-evaluation seam owned by RDR
// 0003. The kernel calls it; it does not implement operator semantics.
type GuardEvaluator interface {
	// Evaluate decides guard over the assembled tag-set view.
	Evaluate(guard string, view TagSet) GuardResult
}

// TagSet is the assembled evaluation view: owned, observed, and freshly
// recognized tags merged into one lookup with provenance preserved.
type TagSet struct {
	tags map[string]taggedValue
}

type taggedValue struct {
	value      string
	provenance Provenance
}

// Lookup returns the value and provenance for key.
func (TagSet) Lookup(key string) (value string, prov Provenance, ok bool) {
	return "", ProvenanceOwned, false
}

// Len reports how many distinct tag keys the view carries.
func (TagSet) Len() int { return 0 }

// Row is one normalized candidate edge from the reviewable transition
// table. RDR 0002 owns normalization; the kernel consumes the normalized
// shape and retains each row's source identity for diagnosis.
type Row struct {
	// RuleID and SourceLocator are the source identity RDR 0002 requires
	// every normalized row to retain. The kernel carries them into an
	// ambiguous-match refusal so RDR 0005 can report conflicting rows.
	RuleID        string
	SourceLocator string

	// Outcome is the recognized outcome this row responds to.
	Outcome string
	// Match is the tag pattern the row requires of the evaluation view.
	Match []Tag
	// RequiresOwned names owned tag keys the row's evaluation needs.
	// A key absent from the owned snapshot yields owned_state_unavailable.
	RequiresOwned []string
	// Guard is the predicate handed to the guard seam. Empty means no
	// guard.
	Guard string

	// NextTags is the next state the row transitions to.
	NextTags []Tag
	// Writes describes the owned-tag writes the accessor layer applies.
	Writes []Tag

	// Escape lists the resolver failure classes this row is modeled to
	// rescue. RDR 0002 restricts the list to no_match and
	// ambiguous_match.
	Escape []RefusalKind
}

// Table is the parsed, reviewable transition table the caller supplies.
type Table struct {
	// Revision is the opaque caller-supplied table revision identity. The
	// kernel carries and compares it; it does not parse it.
	Revision string
	// Outcomes is the declared recognized-outcome alphabet. A recognized
	// outcome outside this alphabet yields unmodeled_outcome.
	Outcomes []string
	// Rows are the normalized candidate edges.
	Rows []Row
}

// Input is the resolution input tuple named by RDR 0001's Identity
// decision: flow identity, transition table revision, accessor-produced
// owned tag snapshot, observed tag-set, and freshly recognized outcome
// tag — plus the reviewable table itself.
type Input struct {
	Flow       string
	Table      Table
	Owned      []Tag
	Observed   []Tag
	Recognized string
	// Guards is the delegated guard-evaluation seam. A nil seam means the
	// table's guarded rows cannot be decided.
	Guards GuardEvaluator
}

// Plan is a successful resolution: the next state tags plus the owned-tag
// writes the accessor layer applies back to the caller-provided artifact
// boundary. The kernel describes writes; it never executes them.
type Plan struct {
	// RuleID and SourceLocator identify the single matching row.
	RuleID        string
	SourceLocator string
	// NextTags is the next state.
	NextTags []Tag
	// Writes are the owned-tag writes for the accessor layer.
	Writes []Tag
	// Revision echoes the transition table revision used.
	Revision string
	// Escaped reports whether the selection came from a modeled escape
	// edge rather than an ordinary exact-one match.
	Escaped bool
}

// Refusal is a modeled, value-level refusal disposition. It is not a CLI
// error and not the Go error path.
type Refusal struct {
	// Kind is the stable discriminator RDR 0005 maps to a CLI code.
	Kind RefusalKind
	// Revision echoes the transition table revision used, for diagnosis.
	Revision string
	// Flow echoes the flow identity, for diagnosis.
	Flow string
	// Recognized echoes the freshly recognized outcome, for diagnosis.
	Recognized string
	// Rows carries source identity for the rows implicated in the
	// refusal — the conflicting rows on ambiguous_match.
	Rows []RowRef
	// MissingOwned names owned tag keys absent from the snapshot, on
	// owned_state_unavailable.
	MissingOwned []string
	// Guard names the predicate the seam could not decide, on
	// guard_unevaluable.
	Guard string
}

// RowRef is the source identity of one implicated table row.
type RowRef struct {
	RuleID        string
	SourceLocator string
}

// Result is the kernel disposition: exactly one transition plan or
// exactly one typed refusal, never both and never neither.
type Result struct {
	// Plan is non-nil exactly when the resolution succeeded.
	Plan *Plan
	// Refusal is non-nil exactly when the kernel refused.
	Refusal *Refusal
}

// Refused reports whether the disposition is a modeled refusal.
func (Result) Refused() bool { return false }

// Resolve is the kernel's single pure entry point. It returns exactly one
// disposition for the input tuple. A modeled refusal travels the Result
// value with a nil error; the error return is reserved for programmer
// mistakes, not for modeled refusals.
func Resolve(in Input) (Result, error) { return Result{}, nil }
