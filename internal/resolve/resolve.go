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
package resolve

import "slices"

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
func RefusalKinds() []RefusalKind {
	return []RefusalKind{
		KindNoMatch,
		KindAmbiguousMatch,
		KindOwnedStateUnavailable,
		KindGuardUnevaluable,
		KindUnmodeledOutcome,
	}
}

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

// recognizedTagKey is the tag key the freshly recognized outcome takes in
// the assembled evaluation view, so a table row can match on it directly
// (REQ-17).
const recognizedTagKey = "recognized"

// Lookup returns the value and provenance for key.
func (s TagSet) Lookup(key string) (value string, prov Provenance, ok bool) {
	tv, ok := s.tags[key]
	if !ok {
		return "", ProvenanceOwned, false
	}
	return tv.value, tv.provenance, true
}

// Len reports how many distinct tag keys the view carries.
func (s TagSet) Len() int { return len(s.tags) }

// has reports whether key is present with the given provenance.
func (s TagSet) has(key string, prov Provenance) bool {
	tv, ok := s.tags[key]
	return ok && tv.provenance == prov
}

// matches reports whether every tag in want is present in the view with
// the same value.
func (s TagSet) matches(want []Tag) bool {
	for _, w := range want {
		tv, ok := s.tags[w.Key]
		if !ok || tv.value != w.Value {
			return false
		}
	}
	return true
}

// assemble merges owned, observed, and freshly recognized tags into one
// evaluation view (REQ-17). Provenance precedence is owned over observed
// over recognized, so the accessor-produced snapshot is never shadowed by
// caller-supplied context. Within one provenance the last tag wins; the
// merge is order-insensitive across provenances, which is what value-level
// replay determinism requires (REQ-3).
func assemble(in Input) TagSet {
	view := TagSet{tags: make(map[string]taggedValue, len(in.Owned)+len(in.Observed)+1)}

	if in.Recognized != "" {
		view.tags[recognizedTagKey] = taggedValue{
			value:      in.Recognized,
			provenance: ProvenanceRecognized,
		}
	}
	for _, t := range in.Observed {
		view.tags[t.Key] = taggedValue{value: t.Value, provenance: ProvenanceObserved}
	}
	for _, t := range in.Owned {
		view.tags[t.Key] = taggedValue{value: t.Value, provenance: ProvenanceOwned}
	}
	return view
}

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

// rescues reports whether the row is modeled to rescue kind.
func (r Row) rescues(kind RefusalKind) bool {
	return slices.Contains(r.Escape, kind)
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

// models reports whether outcome is in the table's declared alphabet.
func (t Table) models(outcome string) bool {
	return slices.Contains(t.Outcomes, outcome)
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
func (r Result) Refused() bool { return r.Refusal != nil }

// Resolve is the kernel's single pure entry point. It returns exactly one
// disposition for the input tuple. A modeled refusal travels the Result
// value with a nil error; the error return is reserved for programmer
// mistakes, not for modeled refusals.
//
// Evaluation order is fixed so the same tuple always replays the same
// disposition (REQ-1):
//
//  1. the recognized outcome must be in the table's declared alphabet,
//     otherwise unmodeled_outcome (which therefore outranks a mere
//     zero-match);
//  2. candidate rows are the non-escape rows for that outcome whose match
//     pattern holds over the assembled view;
//  3. a candidate requiring an owned tag absent from the accessor
//     snapshot yields owned_state_unavailable;
//  4. a candidate whose guard the seam cannot decide yields
//     guard_unevaluable;
//  5. exactly one surviving candidate is the plan; zero is no_match and
//     more than one is ambiguous_match, each subject to rescue by a
//     modeled escape edge that itself matches exactly once.
func Resolve(in Input) (Result, error) {
	view := assemble(in)

	if !in.Table.models(in.Recognized) {
		return refuse(in, Refusal{Kind: KindUnmodeledOutcome}), nil
	}

	var candidates []Row
	for _, row := range in.Table.Rows {
		if len(row.Escape) != 0 {
			// Escape edges are not ordinary candidates; they participate
			// only in the rescue phase for their declared class.
			continue
		}
		if row.Outcome != in.Recognized {
			continue
		}
		if !view.matches(row.Match) {
			continue
		}
		candidates = append(candidates, row)
	}

	if missing := missingOwned(candidates, view); len(missing) > 0 {
		return refuse(in, Refusal{
			Kind:         KindOwnedStateUnavailable,
			MissingOwned: missing,
			Rows:         rowRefsRequiring(candidates, missing),
		}), nil
	}

	var selected []Row
	for _, row := range candidates {
		switch evaluateGuard(in.Guards, row.Guard, view) {
		case GuardTrue:
			selected = append(selected, row)
		case GuardFalse:
			// Decided and does not hold: the row is not a candidate.
		case GuardUnevaluable:
			return refuse(in, Refusal{
				Kind:  KindGuardUnevaluable,
				Guard: row.Guard,
				Rows:  []RowRef{refOf(row)},
			}), nil
		}
	}

	switch len(selected) {
	case 1:
		return Result{Plan: planOf(in, selected[0], false)}, nil
	case 0:
		return escapeOrRefuse(in, view, Refusal{Kind: KindNoMatch}), nil
	default:
		return escapeOrRefuse(in, view, Refusal{
			Kind: KindAmbiguousMatch,
			Rows: rowRefs(selected),
		}), nil
	}
}

// evaluateGuard delegates the predicate to the injected seam (REQ-23). An
// unguarded row needs no seam; a guarded row with no seam is undecidable,
// never evaluated by the kernel itself.
func evaluateGuard(seam GuardEvaluator, guard string, view TagSet) GuardResult {
	if guard == "" {
		return GuardTrue
	}
	if seam == nil {
		return GuardUnevaluable
	}
	return seam.Evaluate(guard, view)
}

// missingOwned returns, in stable order, the owned tag keys any candidate
// requires that the accessor-produced snapshot does not carry.
func missingOwned(candidates []Row, view TagSet) []string {
	var missing []string
	seen := map[string]bool{}
	for _, row := range candidates {
		for _, key := range row.RequiresOwned {
			if view.has(key, ProvenanceOwned) || seen[key] {
				continue
			}
			seen[key] = true
			missing = append(missing, key)
		}
	}
	return missing
}

// escapeOrRefuse rescues a no_match or ambiguous_match refusal when the
// table models exactly one matching escape edge for that class (REQ-15,
// REQ-16). An escape that matches more than once is not an exact-one
// rescue and degrades to ambiguous_match.
func escapeOrRefuse(in Input, view TagSet, r Refusal) Result {
	var escapes []Row
	for _, row := range in.Table.Rows {
		if !row.rescues(r.Kind) {
			continue
		}
		if row.Outcome != in.Recognized || !view.matches(row.Match) {
			continue
		}
		escapes = append(escapes, row)
	}

	switch len(escapes) {
	case 1:
		return Result{Plan: planOf(in, escapes[0], true)}
	case 0:
		return refuse(in, r)
	default:
		return refuse(in, Refusal{Kind: KindAmbiguousMatch, Rows: rowRefs(escapes)})
	}
}

// planOf builds the success disposition for row. Tag slices are copied so
// the plan never aliases the caller's table (REQ-24).
func planOf(in Input, row Row, escaped bool) *Plan {
	return &Plan{
		RuleID:        row.RuleID,
		SourceLocator: row.SourceLocator,
		NextTags:      copyTags(row.NextTags),
		Writes:        copyTags(row.Writes),
		Revision:      in.Table.Revision,
		Escaped:       escaped,
	}
}

// refuse completes a refusal with the input-tuple identity RDR 0001
// requires for diagnosis (REQ-10).
func refuse(in Input, r Refusal) Result {
	r.Revision = in.Table.Revision
	r.Flow = in.Flow
	r.Recognized = in.Recognized
	return Result{Refusal: &r}
}

func refOf(row Row) RowRef {
	return RowRef{RuleID: row.RuleID, SourceLocator: row.SourceLocator}
}

func rowRefs(rows []Row) []RowRef {
	if len(rows) == 0 {
		return nil
	}
	out := make([]RowRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, refOf(row))
	}
	return out
}

// rowRefsRequiring names the candidate rows that asked for one of the
// missing owned keys, so RDR 0005 can report which edge went unevaluable.
func rowRefsRequiring(rows []Row, missing []string) []RowRef {
	want := make(map[string]bool, len(missing))
	for _, key := range missing {
		want[key] = true
	}
	var out []RowRef
	for _, row := range rows {
		for _, key := range row.RequiresOwned {
			if want[key] {
				out = append(out, refOf(row))
				break
			}
		}
	}
	return out
}

func copyTags(in []Tag) []Tag {
	if in == nil {
		return nil
	}
	out := make([]Tag, len(in))
	copy(out, in)
	return out
}
