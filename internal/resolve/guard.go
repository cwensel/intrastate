package resolve

import (
	"slices"
	"strings"
)

// Block names the guard block an atom sits in. RDR 0007 `0007:C1` fixes
// the type as an exported named string type; JDR 0001 §D12 adds
// BlockMatch so one type serves the guard payload and the match pattern's
// future atom representation.
type Block string

const (
	// BlockAll is the conjunctive block: every atom must hold.
	BlockAll Block = "all"
	// BlockUnless is the block-level negated block: the row verdict is
	// all_result ∧ ¬(unless_conj).
	BlockUnless Block = "unless"
	// BlockMatch is the match pattern's block (JDR 0001 §D12). RDR 0007
	// evaluates no match atoms: the row verdict formula `0007:C6` fences
	// names `all` and `unless` only, and match-tag semantics remain RDR
	// 0001's closed-world TagSet.matches.
	BlockMatch Block = "match"
)

// isGuardBlock reports whether the verdict formula `0007:C6` names the
// block as an operand. It names `all` and `unless` and nothing else, so
// every other block — BlockMatch (JDR 0001 §D12), the zero value, and any
// future token — is not a guard operand.
func isGuardBlock(b Block) bool {
	return b == BlockAll || b == BlockUnless
}

// The existence operator's token and its two boolean literal forms. The
// kernel compares an atom's operator and literal against these values
// verbatim, performing no parsing, case-folding, or coercion of its own
// (`0007:C3`).
const (
	// OpExists is the sole total operator's token.
	OpExists = "exists"
	// LiteralTrue is the existence literal that decides TRUE on presence.
	LiteralTrue = "true"
	// LiteralFalse is the existence literal that decides TRUE on absence.
	LiteralFalse = "false"
)

// GuardAtom is one parsed guard predicate: the shape JDR 0001 §D1 fixes.
// The kernel consumes atoms; it does not parse the authored grammar.
type GuardAtom struct {
	// Key is the tag key the atom references, compared to the assembled
	// view by exact string equality (`0007:C4`).
	Key string
	// Operator is the operator token. OpExists is the kernel's; every
	// other token is opaque and belongs to the evaluator seam.
	Operator string
	// Literal is the authored right-hand side, carried as the exact bytes
	// the normalizer emitted.
	Literal string
	// Block is the block the atom sits in.
	Block Block
}

// Reason is the kernel-owned closed set naming why an atom could not be
// decided (`0007:C8`). It is defined by the OUTCOME, not by which
// component produced it.
type Reason string

const (
	// ReasonAbsent: the referenced key was not in the assembled view.
	ReasonAbsent Reason = "absent"
	// ReasonUncomparable: the key was present and its value was not
	// compared to a verdict — the seam answered unevaluable, the atom
	// carried a foreign or missing literal on OpExists, or there was no
	// seam to ask.
	ReasonUncomparable Reason = "uncomparable"
)

// Reasons returns the closed kernel-owned reason set in declaration
// order, mirroring RefusalKinds().
func Reasons() []Reason {
	return []Reason{ReasonAbsent, ReasonUncomparable}
}

// UndecidedAtom names one atom that blocked a row's verdict. It mirrors
// GuardAtom's four fields by name and type, plus the reason.
type UndecidedAtom struct {
	Key      string
	Block    Block
	Operator string
	Literal  string
	Reason   Reason
}

// UndecidedRow names one undecidable row by its source identity and the
// atoms that blocked it.
type UndecidedRow struct {
	RuleID        string
	SourceLocator string
	Atoms         []UndecidedAtom
}

// evaluateAtoms decides one row's guard atom by atom and combines the
// verdicts in the kernel under strong-Kleene three-valued logic
// (`0007:C6`). The row verdict is all_result ∧ ¬(unless_conj); an omitted
// unless block contributes no operand, so the term drops out rather than
// standing in as a vacuous truth or falsity (`0007:C5`).
//
// Every atom of the row is evaluated — there is no short-circuit on a
// decided block — and each atom the pass decides unevaluable emits its
// payload entry at that same step, so no second walk over the atoms is
// ever needed and the seam is consulted at most once per present-key
// value atom (`0007:C8`).
func evaluateAtoms(atoms []GuardAtom, seam GuardEvaluator, view TagSet) (GuardResult, []UndecidedAtom) {
	allResult := GuardTrue
	unlessConj := GuardTrue
	sawUnless := false

	var undecided []UndecidedAtom
	for _, atom := range atoms {
		if !isGuardBlock(atom.Block) {
			// Fail CLOSED at the block boundary, the way `0007:C3` makes the
			// operator boundary fail closed in both directions. `0007:C6`
			// fences the verdict formula over `all` and `unless` only, so an
			// atom in any other block contributes NO operand to either
			// conjunction — the way an omitted `unless` block contributes
			// none (`0007:C5`) — and is never handed to the seam.
			//
			// Sweeping such an atom into `all_result` instead would let it be
			// DECIDED, and a decided-FALSE one would prune its row along with
			// the row's owned-state obligation (D8) — reopening the masking
			// path `0007:C11` ratifies pruning against, whose condition
			// ("GuardFalse can only arise from decided atoms") is stated over
			// GUARD atoms. Reachable via BlockMatch (JDR 0001 §D12), the zero
			// value Block(""), and any future block token.
			//
			// It emits no payload entry either: `0007:C8`'s payload names the
			// atoms that blocked the GUARD verdict, and an atom that
			// contributes no operand cannot have blocked it.
			continue
		}

		verdict, reason, ok := evaluateAtom(atom, seam, view)
		if !ok {
			undecided = appendAtom(undecided, atom, reason)
		}

		switch atom.Block {
		case BlockUnless:
			sawUnless = true
			unlessConj = kleeneAnd(unlessConj, verdict)
		case BlockAll:
			allResult = kleeneAnd(allResult, verdict)
		}
	}

	verdict := allResult
	if sawUnless {
		verdict = kleeneAnd(allResult, kleeneNot(unlessConj))
	}

	if verdict != GuardUnevaluable {
		// Only an undecidable row reports a payload (`0007:C8`); a decided
		// row's entries are discarded.
		return verdict, nil
	}
	slices.SortFunc(undecided, compareUndecidedAtoms)
	return verdict, undecided
}

// evaluateAtom decides one atom. It returns the atom's verdict and, when
// the atom is unevaluable, the reason the payload carries.
//
// Presence is decided first and provenance-blind — the TagSet.Lookup `ok`
// test, never the owned-provenance predicate missingOwned uses
// (`0007:C4`). An existence atom is then decided by the kernel from
// presence alone; a value-comparing atom over an absent key is marked
// unevaluable by the kernel; only a value-comparing atom over a present
// key is handed to the seam. An evaluator that folds absence into false
// cannot be reached, because it is never asked about an absent key
// (`0007:C1`).
func evaluateAtom(atom GuardAtom, seam GuardEvaluator, view TagSet) (verdict GuardResult, reason Reason, decided bool) {
	value, _, present := view.Lookup(atom.Key)

	if atom.Operator == OpExists {
		switch atom.Literal {
		case LiteralTrue:
			return boolResult(present), "", true
		case LiteralFalse:
			return boolResult(!present), "", true
		}
		// A foreign literal on an OpExists atom is the kernel's
		// fail-closed backstop for a normalizer that did not reject it at
		// load. It is never decided from presence, and a missing literal
		// is not LiteralFalse (`0007:C3`). The reason is `uncomparable`
		// whether or not the key is present: what failed is the atom's
		// own literal, not the view — `absent` would name the wrong fact.
		return GuardUnevaluable, ReasonUncomparable, false
	}

	if !present {
		return GuardUnevaluable, ReasonAbsent, false
	}
	if view.conflicting(atom.Key) {
		// The key was supplied more than once within one provenance with
		// differing values, so the view carries no single value to compare
		// (see assemble). The key IS present, so `absent` would be a lie;
		// the value was not compared to a verdict, which is exactly what
		// `uncomparable` names. The seam is not consulted: handing it one
		// of the colliding values would make the verdict a function of
		// slice position, which `0007:C8` forbids.
		return GuardUnevaluable, ReasonUncomparable, false
	}
	if seam == nil {
		// A nil seam is a wiring fact, not a reason of its own: the key is
		// present, so `absent` would be a lie (`0007:C1`).
		return GuardUnevaluable, ReasonUncomparable, false
	}
	switch v := seam.Evaluate(atom, value); v {
	case GuardTrue, GuardFalse:
		return v, "", true
	default:
		return GuardUnevaluable, ReasonUncomparable, false
	}
}

// appendAtom emits one payload entry, suppressing a duplicate by
// construction: two entries equal on all six fields of REQ-55's tuple name
// the same atom, which the kernel must not report twice. The sort orders
// entries; it does not dedupe them (`0007:C8`).
func appendAtom(entries []UndecidedAtom, atom GuardAtom, reason Reason) []UndecidedAtom {
	entry := UndecidedAtom{
		Key:      atom.Key,
		Block:    atom.Block,
		Operator: atom.Operator,
		Literal:  atom.Literal,
		Reason:   reason,
	}
	if slices.Contains(entries, entry) {
		return entries
	}
	return append(entries, entry)
}

// kleeneAnd is strong-Kleene conjunction: min under the TRUTH order
// F < U < T. It is implemented against the normative table, not by
// minimizing the raw GuardResult constant values, whose shipped iota
// order is F < T < U and is not the truth order (`0007:C6`).
func kleeneAnd(a, b GuardResult) GuardResult {
	if a == GuardFalse || b == GuardFalse {
		return GuardFalse
	}
	if a == GuardUnevaluable || b == GuardUnevaluable {
		return GuardUnevaluable
	}
	return GuardTrue
}

// kleeneNot is strong-Kleene negation: ¬T = F, ¬F = T, ¬U = U.
func kleeneNot(v GuardResult) GuardResult {
	switch v {
	case GuardTrue:
		return GuardFalse
	case GuardFalse:
		return GuardTrue
	default:
		return GuardUnevaluable
	}
}

func boolResult(b bool) GuardResult {
	if b {
		return GuardTrue
	}
	return GuardFalse
}

// compareUndecidedAtoms orders payload atoms by the atom half of REQ-55's
// six-field tuple — key, block, operator token, literal — compared field
// by field in that order. The literal is compared as the exact bytes the
// normalizer emitted, which JDR 0001 §D13 makes canonical for a set
// (sorted, duplicate-free, compact JSON array), so the tuple is total.
func compareUndecidedAtoms(a, b UndecidedAtom) int {
	if c := strings.Compare(a.Key, b.Key); c != 0 {
		return c
	}
	if c := strings.Compare(string(a.Block), string(b.Block)); c != 0 {
		return c
	}
	if c := strings.Compare(a.Operator, b.Operator); c != 0 {
		return c
	}
	return strings.Compare(a.Literal, b.Literal)
}

// compareUndecidedRows orders payload rows by REQ-55's six-field tuple:
// source identity first, then — for rows that TIE on identity — the atom
// half, entry by entry over each row's already-sorted atom list.
//
// Identity alone is not total. `0007:C8` fences the sort key as total over
// payload ENTRIES, and the entry key is `(RuleID, SourceLocator, key,
// block, operator token, literal)` — row identity is a COMPONENT of it,
// not a partition boundary. Nothing in this RDR makes `(RuleID,
// SourceLocator)` unique across rows, and `slices.SortFunc` is not stable,
// so comparing identity alone leaves rows that tie on it in `Table.Rows`
// position — the row-order dependence this clause exists to forbid, one
// level above the atom order the entry sort already closes.
//
// Two rows equal under this comparison carry the same identity and the
// same atom entries, so they are indistinguishable in the payload and
// their relative order cannot be observed.
func compareUndecidedRows(a, b UndecidedRow) int {
	if c := strings.Compare(a.RuleID, b.RuleID); c != 0 {
		return c
	}
	if c := strings.Compare(a.SourceLocator, b.SourceLocator); c != 0 {
		return c
	}
	for i := range min(len(a.Atoms), len(b.Atoms)) {
		if c := compareUndecidedAtoms(a.Atoms[i], b.Atoms[i]); c != 0 {
			return c
		}
		if c := strings.Compare(string(a.Atoms[i].Reason), string(b.Atoms[i].Reason)); c != 0 {
			return c
		}
	}
	return len(a.Atoms) - len(b.Atoms)
}
