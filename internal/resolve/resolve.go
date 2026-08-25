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

import (
	"slices"
	"strings"
)

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

// GuardEvaluator is the delegated value-comparison seam owned by RDR
// 0003. The kernel calls it; it does not implement operator semantics.
//
// It is a single-method interface, not a func type, so a missing seam is
// a nil interface value and RDR 0003's evaluator satisfies it by
// declaring the method (`0007:C1`).
type GuardEvaluator interface {
	// Evaluate decides one atom against the present value its key holds
	// in the assembled view. The seam never sees the view: the kernel has
	// already decided presence, so the seam is only ever asked to compare
	// a present value against a literal, and it cannot fold absence into
	// false. It decides true or false under RDR 0003's typed operator
	// semantics, and may answer unevaluable for a present value it cannot
	// compare.
	Evaluate(atom GuardAtom, value string) GuardResult
}

// TagSet is the assembled evaluation view: owned, observed, and freshly
// recognized tags merged into one lookup with provenance preserved.
type TagSet struct {
	tags map[string]taggedValue
}

type taggedValue struct {
	value      string
	provenance Provenance
	// conflicted marks a key that arrived more than once WITHIN one
	// provenance carrying different values. The key is present, but no
	// single value is a function of the input tuple, so a value-comparing
	// guard atom over it is unevaluable (`0007:C8`, see assemble).
	conflicted bool
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

// conflicting reports whether key was supplied more than once within one
// provenance with differing values, so the view carries no single value
// for it (see assemble). Unexported: `0007:C8`'s obligation is on the
// kernel's verdict, and Lookup's exported shape is fenced by `0007:C1`'s
// presence test.
func (s TagSet) conflicting(key string) bool {
	return s.tags[key].conflicted
}

// has reports whether key is present with the given provenance.
func (s TagSet) has(key string, prov Provenance) bool {
	tv, ok := s.tags[key]
	return ok && tv.provenance == prov
}

// matches reports whether every tag in want is present in the view with
// the same value. A CONFLICTED key never matches: the view carries no
// single value for it (see assemble), so comparing the last-merged one
// would make candidate selection a function of slice position rather than
// of the input tuple, which `0007:C8` determinism and deviations D9
// forbid. This mirrors guard.go::evaluateAtom's `ReasonUncomparable`
// posture at the match seam — the key is present, and its value was not
// compared to a verdict. Non-match is the conservative reading: it can
// only turn a plan into a refusal, never the reverse, and it mints no new
// refusal kind (REQ-7 pins the set at five).
func (s TagSet) matches(want []Tag) bool {
	for _, w := range want {
		tv, ok := s.tags[w.Key]
		if !ok || tv.conflicted || tv.value != w.Value {
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
//
// "Last tag wins" is order-insensitive only while no key REPEATS within one
// provenance, and nothing constrains a caller-supplied Observed slice from
// repeating one (A13). Two orderings of the same tag multiset would then
// assemble different values, and a guard verdict — and under D8 whether the
// row is pruned at all — would be a function of slice position rather than
// of the input tuple, which `0007:C8` and RDR 0001 REQ-1 forbid.
//
// Such a key is therefore marked CONFLICTED rather than silently resolved:
// it is present (presence stays provenance-blind and positional-free,
// `0007:C4`), but the view carries no single value for it, so a
// value-comparing atom over it is unevaluable with reason `uncomparable` —
// "the key was present and its value was not compared to a verdict"
// (`0007:C8`). This is the same posture the RDR takes everywhere else: an
// undecidable input is a refusal-class result, never a false allow and
// never a false deny. Both orderings then agree.
//
// A repeat carrying the SAME value is not a conflict — either resolution is
// the same value — and a key crossing PROVENANCES is resolved by the
// precedence above, which is a property of the tuple, not of slice order.
// The merge's control flow, and RDR 0001's non-duplicate merge semantics,
// are otherwise unchanged.
//
// The mark is consumed on BOTH paths that read a value out of the view: the
// GUARD path, where evaluateAtom answers unevaluable with reason
// `uncomparable` rather than consulting the seam, and the MATCH path, where
// matches treats the key as non-matching rather than comparing the retained
// value. Fixing only the guard path leaves `Row.Match` selection positional,
// which is the same defect one filter earlier.
func assemble(in Input) TagSet {
	view := TagSet{tags: make(map[string]taggedValue, len(in.Owned)+len(in.Observed)+1)}

	if in.Recognized != "" {
		view.tags[recognizedTagKey] = taggedValue{
			value:      in.Recognized,
			provenance: ProvenanceRecognized,
		}
	}
	view.merge(in.Observed, ProvenanceObserved)
	view.merge(in.Owned, ProvenanceOwned)
	return view
}

// merge folds one provenance's tags into the view, marking any key the
// slice repeats with differing values as conflicted (see assemble).
func (s TagSet) merge(tags []Tag, prov Provenance) {
	for _, t := range tags {
		next := taggedValue{value: t.Value, provenance: prov}
		if prior, ok := s.tags[t.Key]; ok && prior.provenance == prov {
			next.conflicted = prior.conflicted || prior.value != t.Value
		}
		s.tags[t.Key] = next
	}
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
	// RequiresOwned names the owned tag keys the row's POST-GUARD
	// transition depends on — the keys its Writes require, which per RDR
	// 0002 includes an authored clear (normalization renders it as a
	// `<clear>` write). A key absent from the owned snapshot yields
	// owned_state_unavailable.
	//
	// Guard decidability is not this field's job: guard-input coverage is
	// enforced by the domain rule the atom pipeline implements
	// (`0007:C9`). Listing a guard-read owned key here remains legal and
	// yields the more precise owned_state_unavailable diagnosis — and is
	// the only way a guard over an owned tag is protected from a
	// caller-supplied observed tag satisfying presence. The guard's key
	// set is not required to be a subset of this one, because guards
	// legitimately read observed and recognized tags.
	RequiresOwned []string
	// Guard is the row's predicate as a slice of parsed atoms, the shape
	// JDR 0001 §D1 fixes. An empty slice is an unguarded row. There is no
	// opaque guard string and no reconstruction step (`0007:C1`).
	Guard []GuardAtom

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
	// Undecided names what was missing per row and per atom, on
	// guard_unevaluable: every undecidable row, and for each of its
	// unevaluable atoms the referenced key, the block, the operator, the
	// literal, and the reason. Every undecidable row appears, so there is
	// no representative row to choose (`0007:C8`).
	Undecided []UndecidedRow
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
//  3. every candidate passes the same viability gate (see gate): a row
//     whose guard is decided FALSE is pruned outright and contributes
//     nothing; a surviving row missing required owned state yields
//     owned_state_unavailable; a surviving row whose guard is undecidable
//     yields guard_unevaluable. The guard verdict itself is combined in
//     the kernel from per-atom verdicts (see evaluateAtoms);
//  4. exactly one viable candidate is the plan; zero is no_match and more
//     than one is ambiguous_match, each subject to rescue by a modeled
//     escape edge that is itself viable and matches exactly once.
//
// The gate in step 3 is applied identically to ordinary candidates and to
// escape candidates, so an escape edge can never reach a plan on terms an
// ordinary edge would be refused on (REQ-5, REQ-15, REQ-23).
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

	selected, blocked := gate(candidates, in.Guards, view)
	if blocked != nil {
		return refuse(in, *blocked), nil
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

// gate applies the uniform viability check to rows and partitions them into
// the survivors and, if some row raises a typed blocking condition, the
// refusal that condition warrants.
//
// A row the seam decides GuardFalse is *pruned*: the predicate is decided
// and does not hold, so the row is not an edge at all and its owned-state
// obligation is not part of this resolution's evaluation. This is what
// REQ-15's "exactly one matching edge after guard evaluation" requires —
// a pruned row must not convert a legal exact-one match into a refusal, nor
// mask an escapable zero-match condition.
//
// Among rows the guard does *not* prune, owned state is reported before an
// undecidable guard. The two refusals diagnose independent problems — the
// owned snapshot is missing artifact state, while an undecidable guard is
// a predicate the assembled view cannot settle — and this precedence is
// pinned to shipped behavior, with combined reporting rejected: a row
// failing both ways yields owned_state_unavailable and the owned payload
// only (`0007:C7`).
//
// Rows are visited in table order, but every payload the refusal carries is
// sorted before it is returned, so the disposition stays a function of the
// tuple rather than of slice position (REQ-1, REQ-10).
func gate(rows []Row, seam GuardEvaluator, view TagSet) (selected []Row, blocked *Refusal) {
	var survivors []Row
	verdicts := make([]GuardResult, 0, len(rows))
	payloads := make([][]UndecidedAtom, 0, len(rows))
	for _, row := range rows {
		verdict, atoms := evaluateAtoms(row.Guard, seam, view)
		if verdict == GuardFalse {
			continue
		}
		survivors = append(survivors, row)
		verdicts = append(verdicts, verdict)
		payloads = append(payloads, atoms)
	}

	if missing := missingOwned(survivors, view); len(missing) > 0 {
		return nil, &Refusal{
			Kind:         KindOwnedStateUnavailable,
			MissingOwned: missing,
			Rows:         rowRefsRequiring(survivors, missing),
		}
	}

	var undecidable []Row
	var undecided []UndecidedRow
	for i, row := range survivors {
		switch verdicts[i] {
		case GuardTrue:
			selected = append(selected, row)
		case GuardUnevaluable:
			undecidable = append(undecidable, row)
			undecided = append(undecided, UndecidedRow{
				RuleID:        row.RuleID,
				SourceLocator: row.SourceLocator,
				Atoms:         payloads[i],
			})
		case GuardFalse:
			// Unreachable: pruned above.
		}
	}

	if len(undecidable) > 0 {
		slices.SortFunc(undecided, compareUndecidedRows)
		return nil, &Refusal{
			Kind:      KindGuardUnevaluable,
			Undecided: undecided,
			Rows:      rowRefs(undecidable),
		}
	}

	return selected, nil
}

// missingOwned returns the owned tag keys any candidate requires that the
// accessor-produced snapshot does not carry, sorted by key.
//
// Sorting by key rather than accumulating in row order is what makes the
// payload a function of the input tuple: REQ-2 identifies a resolution by
// the table *revision*, not by row sequence, so two orderings of the same
// row set at the same revision are the same input and must report the same
// diagnosis (REQ-1, REQ-10). Row order is a normalization detail RDR 0002
// owns and must not reach the reported diagnosis.
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
	slices.Sort(missing)
	return missing
}

// escapeOrRefuse rescues a no_match or ambiguous_match refusal when the
// table models exactly one viable escape edge for that class (REQ-15,
// REQ-16). An escape that matches more than once is not an exact-one
// rescue and degrades to ambiguous_match.
//
// Escape candidates pass the same viability gate as ordinary candidates
// before the exact-one count is taken: an escape edge whose guard the seam
// decides FALSE is pruned and does not rescue, and an escape edge missing
// required owned state or carrying an undecidable guard raises that typed
// refusal rather than emitting a plan whose writes the kernel cannot
// justify (REQ-5, REQ-12, REQ-15, REQ-23). The escape path is the kernel's
// most permissive path, so leaving it ungated is where a guessed transition
// would enter.
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

	viable, blocked := gate(escapes, in.Guards, view)
	if blocked != nil {
		return refuse(in, *blocked)
	}

	switch len(viable) {
	case 1:
		return Result{Plan: planOf(in, viable[0], true)}
	case 0:
		return refuse(in, r)
	default:
		return refuse(in, Refusal{Kind: KindAmbiguousMatch, Rows: rowRefs(viable)})
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

// compareRefs orders row references by source identity, so a collection of
// implicated rows can be reported independently of table row order.
func compareRefs(a, b RowRef) int {
	if c := strings.Compare(a.RuleID, b.RuleID); c != 0 {
		return c
	}
	return strings.Compare(a.SourceLocator, b.SourceLocator)
}

// rowRefs names rows by source identity, sorted by that identity. Like
// missingOwned, the sort is what keeps the diagnosis payload a function of
// the input tuple rather than of Table.Rows order (REQ-1, REQ-10).
func rowRefs(rows []Row) []RowRef {
	if len(rows) == 0 {
		return nil
	}
	out := make([]RowRef, 0, len(rows))
	for _, row := range rows {
		out = append(out, refOf(row))
	}
	slices.SortFunc(out, compareRefs)
	return out
}

// rowRefsRequiring names the candidate rows that asked for one of the
// missing owned keys, so RDR 0005 can report which edge went unevaluable.
func rowRefsRequiring(rows []Row, missing []string) []RowRef {
	want := make(map[string]bool, len(missing))
	for _, key := range missing {
		want[key] = true
	}
	var out []Row
	for _, row := range rows {
		for _, key := range row.RequiresOwned {
			if want[key] {
				out = append(out, row)
				break
			}
		}
	}
	return rowRefs(out)
}

func copyTags(in []Tag) []Tag {
	if in == nil {
		return nil
	}
	out := make([]Tag, len(in))
	copy(out, in)
	return out
}
