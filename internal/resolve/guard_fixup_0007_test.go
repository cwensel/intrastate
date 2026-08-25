package resolve_test

// Phase 3c regression tests for RDR 0007. Each test pins one Phase 3a
// (CoVe) finding at its exact reported vector, so a regression reproduces
// the original defect rather than a neighbouring one.
//
//   FAIL-1 (`0007:C6`, `0007:C1`) — `evaluateAtoms`' unfenced `default:`
//     arm folded every block it did not name into `all_result`, so a
//     `BlockMatch` atom (representable on `Row.Guard` since JDR 0001 §D12
//     put it on the same `Block` type) became a guard operand. The two
//     legs the CoVe probe reported are BOTH directions of that: an
//     absent-key match atom turned a PLAN into the non-escapable
//     `guard_unevaluable`, and a seam-decided-FALSE match atom PRUNED its
//     row to `no_match`. §D12 calls a `match` entry in the payload
//     "harmless, since match atoms are never unevaluable" — a claim that
//     only holds while match atoms contribute no operand.
//
//   FAIL-2 (`0007:C8`) — the fenced six-tuple `(RuleID, SourceLocator,
//     key, block, operator token, literal)` was split across two sorts:
//     `compareUndecidedRows` compared row identity ONLY, and
//     `slices.SortFunc` is not stable, so two undecidable rows tying on
//     identity emitted in `Table.Rows` position. `0007:C8` fences the key
//     as TOTAL over payload ENTRIES, which makes row identity a component
//     of the entry key rather than a partition boundary.
//
// ADV-1/2/3 are pinned by `guard_adversarial_0007_test.go`; nothing here
// duplicates them.

import (
	"reflect"
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// --- FAIL-1 --------------------------------------------------------------

// FAIL-1: a `BlockMatch` atom is not an operand of `all_result` or of
// `unless_conj`, so it can neither drive a row unevaluable nor prune it.
// REGRESSION (Phase 3a CoVe).
func TestFix0007Fail1_MatchBlockAtomContributesNoOperand(t *testing.T) {
	matchAtom := func(key, op, literal string) resolve.GuardAtom {
		return resolve.GuardAtom{
			Key: key, Operator: op, Literal: literal, Block: resolve.BlockMatch,
		}
	}

	// The CoVe vector verbatim: a decided-TRUE `all` atom beside a
	// `BlockMatch` atom over an absent key. Observed before the fix:
	// REFUSE guard_unevaluable with a payload entry carrying blk="match".
	t.Run("an absent-key match atom does not make the row unevaluable", func(t *testing.T) {
		row := guardedRow("rdr.fail1.absent", "flows/rdr.toml:300",
			existsAll("status", true),
			matchAtom(absentKey, opEq, "v"),
		)

		got := mustResolve(t, guardInput("rev-fail1-absent", newValueSeam(), row))
		if got.Refused() {
			t.Fatalf("kernel refused %q naming %+v; `0007:C6` fences the row verdict "+
				"over `all` and `unless` only, so a `match` atom contributes no "+
				"operand and the row's decided-TRUE guard atom alone settles it",
				got.Refusal.Kind, got.Refusal.Undecided)
		}
		if got.Plan.Escaped || got.Plan.RuleID != row.RuleID {
			t.Errorf("plan = {RuleID:%q Escaped:%v}; want the ordinary row %q",
				got.Plan.RuleID, got.Plan.Escaped, row.RuleID)
		}
	})

	// The CoVe vector's second case: a lone `BlockMatch` atom the seam
	// decides FALSE. Folded into `all_result` it made the row GuardFalse,
	// which prunes it (D8) and degrades an exact-one match to `no_match`.
	t.Run("a seam-decided-FALSE match atom does not prune the row", func(t *testing.T) {
		row := guardedRow("rdr.fail1.false", "flows/rdr.toml:301",
			matchAtom("status", opEq, "Nope"), // eqSeam decides FALSE: status=Draft
		)

		got := mustResolve(t, guardInput("rev-fail1-false", eqSeam{}, row))
		if got.Refused() {
			t.Fatalf("kernel refused %q; a `match` atom is not a guard operand, so it "+
				"cannot yield GuardFalse and cannot prune its row. `0007:C11` ratifies "+
				"pruning only because \"GuardFalse can only arise from decided atoms\" "+
				"— a condition stated over GUARD atoms", got.Refusal.Kind)
		}
		if got.Plan.RuleID != row.RuleID {
			t.Errorf("plan names %q; want %q", got.Plan.RuleID, row.RuleID)
		}
	})

	// §D12: "0007's per-atom payload may therefore carry `match`, which no
	// refusal names — harmless, since match atoms are never unevaluable."
	// The kernel never evaluates one, so it never emits one.
	t.Run("no match atom reaches the payload", func(t *testing.T) {
		row := guardedRow("rdr.fail1.payload", "flows/rdr.toml:302",
			allAtom("guard-key", opEq, "1"), // absent: the real blocker
			matchAtom("match-key", opEq, "2"),
			matchAtom(absentKey, opGte, "3"),
		)

		r := mustRefuse(t, guardInput("rev-fail1-payload", newValueSeam(), row),
			resolve.KindGuardUnevaluable)

		for _, a := range undecidedAtomsFor(r, row.RuleID) {
			if a.Block != resolve.BlockAll && a.Block != resolve.BlockUnless {
				t.Errorf("payload carries %+v; `0007:C8`'s payload names the atoms "+
					"that blocked the GUARD verdict, and an atom contributing no "+
					"operand cannot have blocked it", a)
			}
		}
	})

	// The seam is the RDR 0003 boundary. An atom the verdict formula does
	// not name is not a guard input, so it is never handed across it.
	t.Run("a match atom is never handed to the seam", func(t *testing.T) {
		seam := newValueSeam()
		row := guardedRow("rdr.fail1.seam", "flows/rdr.toml:303",
			existsAll("status", true),
			matchAtom("status", opEq, "Draft"), // present key: would reach the seam
		)

		mustResolve(t, guardInput("rev-fail1-seam", seam, row))
		if calls := seam.seen(); len(calls) != 0 {
			t.Errorf("seam was consulted %d time(s) (%+v); the kernel evaluates only "+
				"the atoms `0007:C6` names as operands", len(calls), calls)
		}
	})
}

// --- FAIL-2 --------------------------------------------------------------

// FAIL-2: the payload sort key is TOTAL over entries, so two undecidable
// rows tying on `(RuleID, SourceLocator)` are still ordered by their atoms
// and the payload never follows `Table.Rows` position.
// REGRESSION (Phase 3a CoVe).
func TestFix0007Fail2_PayloadSortIsTotalOverRowsTyingOnIdentity(t *testing.T) {
	// The CoVe vector verbatim: N rows sharing one identity, each carrying
	// one distinct absent-key atom, presented forward and reversed.
	dup := func(key string) resolve.Row {
		return guardedRow("SAME", "flows/rdr.toml:SAME", allAtom(key, opEq, "v"))
	}
	keys := []string{"k01", "k02", "k03", "k04", "k05", "k06", "k07",
		"k08", "k09", "k10", "k11", "k12", "k13", "k14"}

	rowsFor := func(order []string, build func(string) resolve.Row) []resolve.Row {
		rows := make([]resolve.Row, 0, len(order))
		for _, k := range order {
			rows = append(rows, build(k))
		}
		return rows
	}

	t.Run("candidate set: permuted table order replays one payload", func(t *testing.T) {
		forward := mustRefuse(t,
			guardInput("rev-fail2", newValueSeam(), rowsFor(keys, dup)...),
			resolve.KindGuardUnevaluable)

		reverse := slices.Clone(keys)
		slices.Reverse(reverse)
		reversed := mustRefuse(t,
			guardInput("rev-fail2", newValueSeam(), rowsFor(reverse, dup)...),
			resolve.KindGuardUnevaluable)

		if !reflect.DeepEqual(forward.Undecided, reversed.Undecided) {
			t.Errorf("payload follows Table.Rows position for rows tying on "+
				"(RuleID, SourceLocator). `0007:C8` fences the sort key as TOTAL over "+
				"payload ENTRIES — row identity is a COMPONENT of the entry key, not a "+
				"partition boundary — and `slices.SortFunc` is not stable, so identity "+
				"alone leaves the tie to row order:\n forward = %+v\nreversed = %+v",
				forward.Undecided, reversed.Undecided)
		}
	})

	// The escape set reaches the same `gate` by delegation (`0007:C7`), so
	// the same tie is reachable there. The CoVe probe reproduced it on both.
	t.Run("escape set: permuted table order replays one payload", func(t *testing.T) {
		dupEscape := func(key string) resolve.Row {
			return conformingEscapeRow("SAME.esc", "flows/rdr.toml:SAME.esc",
				resolve.KindNoMatch, allAtom(key, opEq, "v"))
		}

		forward := mustRefuse(t,
			guardInput("rev-fail2-esc", newValueSeam(), rowsFor(keys, dupEscape)...),
			resolve.KindGuardUnevaluable)

		reverse := slices.Clone(keys)
		slices.Reverse(reverse)
		reversed := mustRefuse(t,
			guardInput("rev-fail2-esc", newValueSeam(), rowsFor(reverse, dupEscape)...),
			resolve.KindGuardUnevaluable)

		if !reflect.DeepEqual(forward.Undecided, reversed.Undecided) {
			t.Errorf("escape-set payload follows Table.Rows position for rows tying "+
				"on identity; the escape set is gated by delegation to the same "+
				"`gate` (`0007:C7`) and inherits the same obligation:\n"+
				" forward = %+v\nreversed = %+v", forward.Undecided, reversed.Undecided)
		}
	})

	// Rows tying on identity AND carrying atoms that tie on key differ only
	// in block/operator/literal — the level below the row that `0007:C8`
	// fences the remaining four fields for.
	t.Run("rows tying on identity and key are ordered by the atom fields", func(t *testing.T) {
		variants := []resolve.GuardAtom{
			allAtom("k", opEq, "1"),
			allAtom("k", opEq, "2"),
			allAtom("k", opGte, "1"),
			unlessAtom("k", opEq, "1"),
		}

		build := func(order []int) resolve.Input {
			rows := make([]resolve.Row, 0, len(order))
			for _, i := range order {
				rows = append(rows, guardedRow("TIE", "flows/rdr.toml:TIE", variants[i]))
			}
			return guardInput("rev-fail2-tie", newValueSeam(), rows...)
		}

		forward := mustRefuse(t, build([]int{0, 1, 2, 3}), resolve.KindGuardUnevaluable)
		permuted := mustRefuse(t, build([]int{3, 1, 0, 2}), resolve.KindGuardUnevaluable)

		if !reflect.DeepEqual(forward.Undecided, permuted.Undecided) {
			t.Errorf("payload for identity-tying rows depends on row order once their "+
				"atoms also tie on key:\n forward = %+v\npermuted = %+v",
				forward.Undecided, permuted.Undecided)
		}
	})
}
