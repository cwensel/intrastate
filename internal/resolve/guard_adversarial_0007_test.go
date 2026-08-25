package resolve_test

// Phase 3b adversarial failure-mode tests for RDR 0007 (guard predicate
// totality over an incomplete evaluation view).
//
// These are NOT per-clause conformance tests — `guard_totality_test.go`,
// `guard_atoms_test.go`, and `guard_mvv_test.go` already carry those, and
// they are green. What is exercised here is COMPOSITION: the places where
// two clauses that each hold in isolation admit a case neither one names.
//
// Three failure modes, each anchored to a Failure Modes element of the
// record:
//
//   ADV-1 (`0007:F4`) — the domain rule is scoped to GUARDS. A `match`
//     atom carried in `Row.Guard` is not a guard atom, and `0007:C6`
//     fences the row verdict over `all` and `unless` only. The kernel's
//     combinator has no `match` case: its `default:` arm sweeps every
//     block it does not name into the conjunctive reading, so a match
//     atom over an absent key refuses `guard_unevaluable` — NON-escapable
//     — where F4 requires the escapable `no_match` the closed-world match
//     pattern produces.
//
//   ADV-2 (`0007:F8` + `0007:C11`) — the OPERATOR boundary fails closed
//     (`0007:C3`: a foreign token is a value atom, a foreign literal is
//     `uncomparable`). The BLOCK boundary has no such rule and fails
//     OPEN, into `all`. A decided-FALSE atom in an unrecognized block
//     therefore PRUNES its row, the row's owned-state obligation goes
//     with it (D8), and a modeled `no_match` escape rescues with
//     `Escaped:true` — the Background masking probe (kata `xg7p`) that
//     this RDR exists to close, reproduced through a boundary the
//     fail-closed clause never reached.
//
//   ADV-3 (`0007:F6` + `0007:C8`) — caller-supplied observed tags are
//     unconstrained (A13). `assemble` resolves a key duplicated WITHIN one
//     provenance by "last tag wins", so the assembled value — and with it
//     a guard verdict, and with it whether the row is pruned — is a
//     function of SLICE POSITION, not of the input tuple. `0007:C8`
//     requires the disposition to be a function of the tuple; the shipped
//     doc comment on `assemble` asserts order-insensitivity that only
//     holds when no key repeats. The two orderings of one duplicated
//     observed key give opposite dispositions, and one of them is the
//     masking plan again.
//
// Every fixture drives verdicts THROUGH atoms and the value stub, never by
// injecting a row verdict, per the Testing Strategy preamble.

import (
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// --- ADV-1 ---------------------------------------------------------------

// ADV-1 anchors `0007:F4`: "The domain rule is scoped to guards, so the
// match pattern still folds absence into non-match … predicate PLACEMENT
// decides whether an absent key refuses non-escapably or escapes."
//
// §D12 put `BlockMatch` on the same `Block` type "so one type serves the
// guard payload and the match pattern's future atom representation" — one
// type, no separate slice on `Row`. That is precisely what makes this
// reachable: a match atom is now REPRESENTABLE in `Row.Guard`, and the
// kernel's combinator gives it conjunctive guard semantics.
//
// `0007:C6` fences the verdict formula: "The row verdict is `all_result ∧
// ¬(unless_conj)`, where `unless_conj` is the conjunction of the `unless`
// block's atoms." A `match` atom is an operand of neither term. It must
// contribute NO operand — the same way an omitted `unless` block
// contributes none (`0007:C5`) — leaving the row's guard verdict decided
// by its guard atoms alone.
//
// Currently it contributes to `all_result`, so a match atom over an absent
// key makes the row `guard_unevaluable`: the non-escapable refusal F4
// explicitly says match-pattern absence must NOT produce.
func TestAdv0007_1_MatchBlockAtomIsNotAGuardOperand(t *testing.T) {
	matchAtom := func(key, op, literal string) resolve.GuardAtom {
		return resolve.GuardAtom{
			Key: key, Operator: op, Literal: literal, Block: resolve.BlockMatch,
		}
	}

	t.Run("an absent-key match atom must not refuse guard_unevaluable", func(t *testing.T) {
		// The row's own guard atoms are decided TRUE, so the row's guard
		// verdict is TRUE and it is the single viable candidate. Only the
		// match atom can change that — and `0007:C6` gives it no operand.
		row := guardedRow("rdr.adv1.absent", "flows/rdr.toml:200",
			existsAll("status", true),
			matchAtom(absentKey, opEq, "3"),
		)

		got := mustResolve(t, guardInput("rev-adv1-absent", newValueSeam(), row))
		if got.Refused() {
			t.Fatalf("kernel refused %q naming %+v; a `match` atom is an operand of "+
				"neither `all_result` nor `unless_conj` (`0007:C6` fences the verdict "+
				"formula over `all` and `unless` only), so it cannot drive the row "+
				"guard unevaluable. `0007:F4` scopes the domain rule to GUARDS: an "+
				"absent key referenced by the MATCH pattern folds into non-match and "+
				"yields the escapable `no_match`, never the non-escapable "+
				"`guard_unevaluable`",
				got.Refusal.Kind, got.Refusal.Undecided)
		}
		if got.Plan.RuleID != row.RuleID {
			t.Errorf("plan names %q; want the row %q, whose guard atoms are all decided TRUE",
				got.Plan.RuleID, row.RuleID)
		}
	})

	t.Run("a match atom never reaches the guard payload", func(t *testing.T) {
		// A row genuinely undecidable on a GUARD atom must not have a match
		// atom's key reported beside it: `0007:C8`'s payload names "what was
		// missing" for the guard, and a match key was never a guard input.
		row := guardedRow("rdr.adv1.payload", "flows/rdr.toml:201",
			allAtom("guard-key", opEq, "1"),   // absent -> the real blocker
			matchAtom("match-key", opEq, "2"), // absent, but a match atom
		)

		r := mustRefuse(t, guardInput("rev-adv1-payload", newValueSeam(), row),
			resolve.KindGuardUnevaluable)

		for _, a := range undecidedAtomsFor(r, row.RuleID) {
			if a.Block == resolve.BlockMatch {
				t.Errorf("guard payload carries the match atom %+v; the payload names "+
					"the atoms that blocked the GUARD verdict, and a `match` atom is "+
					"not a guard operand (`0007:C6`, `0007:F4`)", a)
			}
		}
	})
}

// --- ADV-2 ---------------------------------------------------------------

// ADV-2 anchors `0007:F8` ("Evaluator/grammar version skew") read against
// `0007:C3`'s fail-closed rule, and `0007:C11` (D8 ratification).
//
// `0007:C3` states the drift rule for the OPERATOR boundary and makes it
// fail closed in both directions: a foreign token is demoted to a value
// atom (unevaluable on absence), a foreign literal is `uncomparable`.
// `internal/resolve` honours that exactly — `"Exists"`, `"EXISTS"`, and a
// future operator token all fail closed.
//
// The BLOCK boundary has no such rule, and the kernel fails OPEN: the
// combinator's `default:` arm sweeps every unrecognized block into
// `all_result`. So an atom in a block the kernel does not know can be
// DECIDED there, and a decided-FALSE one PRUNES the row.
//
// `0007:C11` ratifies pruning "conditional on the domain rule: pruning is
// safe exactly because GuardFalse can only arise from decided atoms." That
// condition is stated over GUARD atoms. An atom in an unfenced block is
// not one, so the condition does not cover it — and the consequence is the
// Background probe verbatim: "A FALSE guard plus an absent `RequiresOwned`
// key plus a modeled `no_match` escape yields a plan via the escape,
// `Escaped:true`. This is the path where missing artifact state is masked."
//
// Both legs use blocks the kernel cannot claim to have evaluated: the
// §D12 `match` block, and the zero value of `Block` — what an unset field
// on a producer-built atom carries.
func TestAdv0007_2_UnfencedBlockFailsOpenAndReopensTheMaskingPath(t *testing.T) {
	// `gate` absent from the owned snapshot is the missing artifact state.
	// The row's Writes depend on it, so `RequiresOwned` names it
	// (`0007:C9`), and the honest disposition is owned_state_unavailable.
	maskingTable := func(block resolve.Block) resolve.Input {
		row := guardedRow("rdr.adv2.row", "flows/rdr.toml:210",
			resolve.GuardAtom{
				Key:      "gate",
				Operator: resolve.OpExists,
				Literal:  resolve.LiteralTrue,
				Block:    block,
			},
		)
		row.Writes = []resolve.Tag{{Key: "status", Value: "Final"}, {Key: "gate", Value: "open"}}
		row.RequiresOwned = []string{"status", "gate"}

		escape := conformingEscapeRow("rdr.adv2.escape", "flows/rdr.toml:211",
			resolve.KindNoMatch)

		return guardInput("rev-adv2", newValueSeam(), row, escape)
	}

	for name, block := range map[string]resolve.Block{
		"§D12 match block":   resolve.BlockMatch,
		"zero-value block":   resolve.Block(""),
		"future block token": resolve.Block("any"),
	} {
		t.Run(name, func(t *testing.T) {
			got := mustResolve(t, maskingTable(block))
			if !got.Refused() {
				t.Fatalf("kernel emitted a plan for %q (escaped=%v): the atom's block %q "+
					"is not one the verdict formula names, yet it was DECIDED as an "+
					"`all` conjunct, decided FALSE, and pruned the row — taking its "+
					"owned-state obligation on `gate` with it (D8) and letting the "+
					"modeled `no_match` escape rescue. This is the Background masking "+
					"probe (kata `xg7p`) that `0007:C11` ratifies pruning against: "+
					"\"GuardFalse can only arise from DECIDED atoms\". `0007:C3` makes "+
					"operator-token drift fail CLOSED; the block boundary must too "+
					"(`0007:F8`)",
					got.Plan.RuleID, got.Plan.Escaped, block)
			}
			if got.Refusal.Kind != resolve.KindOwnedStateUnavailable {
				t.Errorf("refusal kind = %q; want %q naming the absent owned key `gate`",
					got.Refusal.Kind, resolve.KindOwnedStateUnavailable)
			}
			if !hasItem(got.Refusal.MissingOwned, "gate") {
				t.Errorf("MissingOwned = %v; want it to name `gate`", got.Refusal.MissingOwned)
			}
		})
	}
}

// --- ADV-3 ---------------------------------------------------------------

// ADV-3 anchors `0007:F6` ("Refusal defeated by caller-supplied state" —
// A13, "caller-supplied observed tags are unconstrained") against
// `0007:C8`'s determinism obligation: the kernel "MUST sort the payload so
// it is a function of the input tuple (RDR 0001 REQ-1), never of atom or
// row order."
//
// `0007:C8` closes atom order and row order. It does not reach the third
// ordering the kernel is exposed to — the order of the caller's TAG
// slices. `resolve.go::assemble` resolves that by "Within one provenance
// the last tag wins", and its own doc comment then claims "the merge is
// order-insensitive across provenances, which is what value-level replay
// determinism requires (REQ-3)". The claim holds only while no key
// repeats within a provenance, and nothing constrains a caller-supplied
// `Observed` slice from repeating one (A13).
//
// The consequence is not cosmetic. The two orderings of ONE duplicated
// observed key produce opposite guard verdicts, which under D8 is the
// difference between a pruned row and a survivor — and therefore between
// the masking plan and the honest refusal. `0007:F6`'s accepted risk is
// that an operator can supply a missing key; it does not extend to an
// operator flipping the disposition by REORDERING tags they already
// supplied.
//
// The fix is in the kernel's assembly, not in the caller: a duplicated key
// must resolve to one value by a rule that is a function of the tuple
// (reject it as a malformed input tuple, or make the collision itself
// unevaluable). Either way the two orderings must agree.
func TestAdv0007_3_DuplicateKeyMakesTheVerdictAFunctionOfSlicePosition(t *testing.T) {
	const key = "gate"

	// The row plans only if the guard holds. `eqSeam` decides `eq` by exact
	// byte equality on the present value, so the assembled value for `gate`
	// decides the row: "open" -> TRUE (survivor), anything else -> FALSE
	// (pruned).
	build := func(values ...string) resolve.Input {
		row := guardedRow("rdr.adv3.row", "flows/rdr.toml:220", allAtom(key, opEq, "open"))
		row.Writes = []resolve.Tag{{Key: "status", Value: "Final"}}
		row.RequiresOwned = []string{"status", "audit"} // `audit` is absent from owned

		escape := conformingEscapeRow("rdr.adv3.escape", "flows/rdr.toml:221",
			resolve.KindNoMatch)

		in := guardInput("rev-adv3", eqSeam{}, row, escape)
		for _, v := range values {
			in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: v})
		}
		return in
	}

	t.Run("the same tag multiset in two orders replays identically", func(t *testing.T) {
		forward := mustResolve(t, build("open", "closed"))
		reversed := mustResolve(t, build("closed", "open"))

		describe := func(r resolve.Result) string {
			if r.Refused() {
				return "refusal " + string(r.Refusal.Kind)
			}
			return "plan " + r.Plan.RuleID
		}

		if forward.Refused() != reversed.Refused() {
			t.Fatalf("[gate=open, gate=closed] yields %s but [gate=closed, gate=open] "+
				"yields %s. The two inputs carry the SAME tag multiset and differ only "+
				"in slice position, so RDR 0001 REQ-1 and `0007:C8` make them one input "+
				"tuple with one disposition. `assemble`'s within-provenance \"last tag "+
				"wins\" makes the assembled value — and with it the guard verdict — "+
				"positional, and `assemble`'s doc comment asserts the order-insensitivity "+
				"this disproves",
				describe(forward), describe(reversed))
		}
		if forward.Refused() && forward.Refusal.Kind != reversed.Refusal.Kind {
			t.Errorf("refusal kinds differ across the permutation: %q vs %q",
				forward.Refusal.Kind, reversed.Refusal.Kind)
		}
	})

	t.Run("neither ordering may mask the absent owned key", func(t *testing.T) {
		// The row's Writes depend on owned `audit`, which the snapshot does
		// not carry. Whichever value the duplicate resolves to, the honest
		// dispositions are owned_state_unavailable (survivor) or a refusal —
		// never an `Escaped:true` plan routing around the missing state.
		for name, in := range map[string]resolve.Input{
			"open first":   build("open", "closed"),
			"closed first": build("closed", "open"),
		} {
			t.Run(name, func(t *testing.T) {
				got := mustResolve(t, in)
				if !got.Refused() && got.Plan.Escaped {
					t.Errorf("kernel emitted the escape plan %q with Escaped=true: a "+
						"caller-supplied duplicate observed key pruned the row under D8 "+
						"and routed around the absent owned key `audit`. `0007:F6` "+
						"accepts that an operator can SUPPLY a missing key; it does not "+
						"accept an operator flipping the disposition by reordering tags "+
						"they already supplied", got.Plan.RuleID)
				}
			})
		}
	})

	t.Run("a duplicated OWNED key is equally positional", func(t *testing.T) {
		// The accessor layer produces the owned snapshot (RDR 0004), so a
		// repeated owned key is a producer defect rather than caller input —
		// but the kernel is the component that must not let a producer's
		// slice order decide a verdict.
		ownedBuild := func(values ...string) resolve.Input {
			row := guardedRow("rdr.adv3.owned", "flows/rdr.toml:222", allAtom(key, opEq, "open"))
			in := guardInput("rev-adv3-owned", eqSeam{}, row)
			for _, v := range values {
				in.Owned = append(in.Owned, resolve.Tag{Key: key, Value: v})
			}
			return in
		}

		forward := mustResolve(t, ownedBuild("open", "closed"))
		reversed := mustResolve(t, ownedBuild("closed", "open"))
		if forward.Refused() != reversed.Refused() {
			t.Errorf("a duplicated OWNED key is resolved by slice position too: "+
				"[open, closed] refused=%v, [closed, open] refused=%v",
				forward.Refused(), reversed.Refused())
		}
	})
}
