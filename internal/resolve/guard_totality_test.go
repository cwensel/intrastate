package resolve_test

// RDR 0007 — the domain rule itself: value-operator partiality, existence
// totality, empty-block identities, strong-Kleene combination in the
// kernel, gate ordering, and the per-row/per-atom payload.
//
// The coverage floor this file exists to hold: the spike mutation-tested
// the combinator and found that swapping strong-Kleene FALSE-dominance for
// UNEVALUABLE-dominance SURVIVES all 154 frozen tests, because every
// frozen guard is a single atom. The multi-atom scenarios below are what
// kill that mutant (REQ-74).

import (
	"reflect"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-18: "Value-comparing guard operators (equality, membership, bounded
// integer comparison, set containment) are PARTIAL over the assembled
// evaluation view: an atom whose referenced tag key is absent from the
// view MUST evaluate to unevaluable — never to false and never to true."
// DOMAIN EDGE
//
// The discriminating control is the escape row: absence folded to FALSE
// prunes the row and the escape plans; absence folded to TRUE selects the
// row and it plans. Only "unevaluable" refuses.
func TestReq18_ValueOperatorsArePartialOverTheView(t *testing.T) {
	for _, op := range []string{opEq, opIn, opGte, opContains} {
		t.Run(op, func(t *testing.T) {
			row := guardedRow("rdr.partial", "flows/rdr.toml:40", allAtom(absentKey, op, "3"))
			escape := conformingEscapeRow("rdr.escape.partial", "flows/rdr.toml:41", resolve.KindNoMatch)

			r := mustRefuse(t,
				guardInput("rev-req18", eqSeam{}, row, escape),
				resolve.KindGuardUnevaluable)

			if !hasItem(undecidedRuleIDs(r), row.RuleID) {
				t.Errorf("payload names %v; want the partial row %q — never false "+
					"(which would let the escape plan) and never true (which would "+
					"let the row plan)", undecidedRuleIDs(r), row.RuleID)
			}
		})
	}
}

// REQ-19: "An absent set-valued tag MUST be treated as unevaluable under
// set containment, NOT as the empty set."
// DOMAIN EDGE
//
// Testing Strategy row 8. The set literal is carried in the §D13 form —
// canonical JSON array, members sorted, duplicate-free, compact — which
// JDR 0001 §D13 assigns to RDR 0002 and which BUILD-ORDER directs this RDR
// to write against directly.
func TestReq19_AbsentSetValuedTagIsUnevaluableNotTheEmptySet(t *testing.T) {
	literal := d13Set("beta", "alpha", "alpha")
	if literal != `["alpha","beta"]` {
		t.Fatalf("§D13 carriage form = %q; want a sorted, duplicate-free, compact "+
			"JSON array", literal)
	}

	seam := newValueSeam()
	row := guardedRow("rdr.contains.absent", "flows/rdr.toml:42",
		allAtom("labels", opContains, literal))
	escape := conformingEscapeRow("rdr.escape.contains", "flows/rdr.toml:43", resolve.KindNoMatch)

	r := mustRefuse(t,
		guardInput("rev-req19", seam, row, escape),
		resolve.KindGuardUnevaluable)

	if calls := seam.seen(); len(calls) != 0 {
		t.Errorf("seam saw %v; the absent set-valued tag is decided by the KERNEL "+
			"on presence alone", calls)
	}
	reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), "labels", resolve.BlockAll)
	if !ok || reason != resolve.ReasonAbsent {
		t.Errorf("reason = %q (found=%v); want %q — treating the absent tag as the "+
			"empty set would decide `contains` FALSE and let the escape plan",
			reason, ok, resolve.ReasonAbsent)
	}
}

// REQ-72 / Testing Strategy row 8, present-key half: a PRESENT set-valued
// tag whose value is the §D13 canonical JSON array reaches the seam as
// exactly those bytes, and the kernel neither parses nor re-canonicalizes
// it.
// DOMAIN EDGE
func TestReq72_PresentSetValuedTagCrossesTheSeamAsD13Bytes(t *testing.T) {
	value := d13Set("review", "approve")
	literal := d13Set("approve")

	atom := allAtom("labels", opContains, literal)
	seam := newValueSeam().decide(atom, value, resolve.GuardTrue)

	row := guardedRow("rdr.contains.present", "flows/rdr.toml:44", atom)
	row.RequiresOwned = nil

	in := guardInput("rev-req72", seam, row)
	in.Observed = append(in.Observed, resolve.Tag{Key: "labels", Value: value})

	mustPlan(t, in)

	calls := seam.seen()
	if len(calls) != 1 {
		t.Fatalf("seam saw %v; want exactly one call", calls)
	}
	if calls[0].Value != value {
		t.Errorf("seam received value %q; want the exact §D13 bytes %q — the kernel "+
			"neither parses nor re-canonicalizes the carriage form", calls[0].Value, value)
	}
	if calls[0].Literal != literal {
		t.Errorf("seam received literal %q; want %q", calls[0].Literal, literal)
	}
}

// REQ-20: "The existence operator is the sole TOTAL operator: it MUST
// decide true or false from presence or absence of its referenced key
// alone. Its verdict is `presence == literal`: `exists = true` decides
// TRUE when the key is present, `exists = false` decides TRUE when it is
// absent."
// HAPPY PATH
//
// Testing Strategy row 7 requires BOTH polarity channels.
func TestReq20_ExistenceIsTheSoleTotalOperator(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		literal bool
		want    resolve.GuardResult // the atom's verdict
	}{
		{"present key, exists=true", "reviews", true, resolve.GuardTrue},
		{"present key, exists=false", "reviews", false, resolve.GuardFalse},
		{"absent key, exists=true", absentKey, true, resolve.GuardFalse},
		{"absent key, exists=false", absentKey, false, resolve.GuardTrue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// The escape row makes FALSE and UNEVALUABLE distinguishable:
			// FALSE prunes and the escape plans; UNEVALUABLE refuses.
			row := guardedRow("rdr.total", "flows/rdr.toml:45", existsAll(tc.key, tc.literal))
			escape := conformingEscapeRow("rdr.escape.total", "flows/rdr.toml:46", resolve.KindNoMatch)

			got := mustResolve(t, guardInput("rev-req20", newValueSeam(), row, escape))

			if got.Refused() {
				t.Fatalf("existence atom refused %q; the existence operator is TOTAL "+
					"and never refuses", got.Refusal.Kind)
			}
			switch tc.want {
			case resolve.GuardTrue:
				if got.Plan.RuleID != row.RuleID || got.Plan.Escaped {
					t.Errorf("plan names %q (escaped=%v); want the guarded row %q "+
						"selected", got.Plan.RuleID, got.Plan.Escaped, row.RuleID)
				}
			case resolve.GuardFalse:
				if got.Plan.RuleID != escape.RuleID || !got.Plan.Escaped {
					t.Errorf("plan names %q (escaped=%v); want the row pruned and the "+
						"escape %q reached", got.Plan.RuleID, got.Plan.Escaped, escape.RuleID)
				}
			}
		})
	}
}

// REQ-21: "Polarity lives in the literal; `all`/`unless` placement composes
// with it. Absence tests MUST be expressed as existence atoms; no
// value-comparing operator may act as an implicit existence test."
// DOMAIN EDGE
func TestReq21_PolarityLivesInTheLiteralAndComposesWithBlockPlacement(t *testing.T) {
	// `unless (X exists = true)` over an ABSENT key: unless_conj = F,
	// ¬F = T, so the row is selected — the placement composes with the
	// literal rather than replacing it.
	t.Run("unless(exists=true) over an absent key selects the row", func(t *testing.T) {
		row := guardedRow("rdr.polarity.unless", "flows/rdr.toml:47",
			existsUnless(absentKey, true))
		mustPlan(t, guardInput("rev-req21-unless", nil, row))
	})

	// `unless (X exists = true)` over a PRESENT key: unless_conj = T,
	// ¬T = F, so the row is pruned.
	t.Run("unless(exists=true) over a present key prunes the row", func(t *testing.T) {
		row := guardedRow("rdr.polarity.unless.present", "flows/rdr.toml:48",
			existsUnless("reviews", true))
		escape := conformingEscapeRow("rdr.escape.polarity", "flows/rdr.toml:49", resolve.KindNoMatch)

		plan := mustPlan(t, guardInput("rev-req21-present", nil, row, escape))
		if plan.RuleID != escape.RuleID {
			t.Errorf("plan names %q; want the row pruned and the escape %q reached",
				plan.RuleID, escape.RuleID)
		}
	})

	// No value-comparing operator acts as an implicit existence test: an
	// `eq` over an absent key does NOT decide FALSE the way `exists=true`
	// would.
	t.Run("a value operator is not an implicit existence test", func(t *testing.T) {
		valueRow := guardedRow("rdr.polarity.value", "flows/rdr.toml:50",
			allAtom(absentKey, opEq, "3"))
		existsRow := guardedRow("rdr.polarity.exists", "flows/rdr.toml:51",
			existsAll(absentKey, true))
		escape := conformingEscapeRow("rdr.escape.implicit", "flows/rdr.toml:52", resolve.KindNoMatch)

		// The existence row decides FALSE and prunes: the escape plans.
		plan := mustPlan(t, guardInput("rev-req21-exists", eqSeam{}, existsRow, escape))
		if plan.RuleID != escape.RuleID {
			t.Fatalf("exists=true over an absent key: plan names %q; want the escape %q",
				plan.RuleID, escape.RuleID)
		}
		// The value row does NOT: it refuses.
		mustRefuse(t, guardInput("rev-req21-value", eqSeam{}, valueRow, escape),
			resolve.KindGuardUnevaluable)
	})
}

// REQ-30: "A row with no atoms is decided TRUE without consulting the view
// (shipped: `evaluateGuard`'s empty-guard branch)."
// REQ-31: "an empty `all` block is TRUE (empty conjunction)."
// REQ-32: "An empty or omitted `unless` block is ABSENT: the
// `¬(unless_conj)` TERM DROPS OUT of the row verdict, which reduces to
// `all_result`. It is NOT a vacuously-true conjunction … and it is NOT
// `unless_conj = F` either"
// INPUT EDGE
//
// Testing Strategy row 10. REQ-32's discriminating shape: a row with an
// `all` block only. Under `unless_conj = T` the row verdict would be
// `all ∧ ¬T = F` and EVERY such row would be disabled; under
// `unless_conj = F` it would be `all ∧ ¬F = all`, which coincides with the
// stated rule for a TRUE `all` — so the disabling reading is the one this
// test kills, and the FALSE reading is killed by the F-and-U legs below.
func TestReq30to32_EmptyAtomAndEmptyBlockIdentities(t *testing.T) {
	t.Run("REQ-30 no atoms ⇒ TRUE without consulting the view", func(t *testing.T) {
		row := guardedRow("rdr.empty.row", "flows/rdr.toml:53")
		// A hostile seam that answers FALSE to everything proves the view
		// was never consulted for a verdict.
		mustPlan(t, guardInput("rev-req30", falseSeam{}, row))
	})

	t.Run("REQ-31 empty all block ⇒ TRUE", func(t *testing.T) {
		// The `all` block is empty; the row carries only an `unless` atom
		// that is decided FALSE, so unless_conj = F and ¬F = T.
		atom := unlessAtom("status", opEq, "Final")
		row := guardedRow("rdr.empty.all", "flows/rdr.toml:54", atom)
		mustPlan(t, guardInput("rev-req31", eqSeam{}, row))
	})

	t.Run("REQ-32 omitted unless block contributes no operand", func(t *testing.T) {
		// Only `all` atoms, all decided TRUE. Under unless_conj = T the
		// verdict would be F and the row would be disabled.
		row := guardedRow("rdr.empty.unless", "flows/rdr.toml:55",
			allAtom("status", opEq, "Draft"),
			existsAll("reviews", true),
		)
		plan := mustPlan(t, guardInput("rev-req32", eqSeam{}, row))
		if plan.RuleID != row.RuleID {
			t.Errorf("plan names %q; want %q — an omitted unless block drops the "+
				"¬(unless_conj) term out, it does not disable the row",
				plan.RuleID, row.RuleID)
		}
	})

	t.Run("REQ-32 empty unless is not unless_conj = F over a FALSE all", func(t *testing.T) {
		// all = F, unless omitted ⇒ verdict F ⇒ pruned. Under either
		// mistaken reading of the empty block the `all` result still
		// dominates, so this leg pins the reduction to all_result itself.
		row := guardedRow("rdr.empty.unless.false", "flows/rdr.toml:56",
			allAtom("status", opEq, "Final"))
		escape := conformingEscapeRow("rdr.escape.empty", "flows/rdr.toml:57", resolve.KindNoMatch)

		plan := mustPlan(t, guardInput("rev-req32-false", eqSeam{}, row, escape))
		if plan.RuleID != escape.RuleID {
			t.Errorf("plan names %q; want the row pruned (verdict reduces to "+
				"all_result = F) and the escape %q reached", plan.RuleID, escape.RuleID)
		}
	})
}

// REQ-33: "Atom verdicts combine in the KERNEL under strong-Kleene
// three-valued logic. Conjunction is `min` under the TRUTH order
// `F < U < T`; negation is `¬T = F`, `¬F = T`, `¬U = U`" — with the
// normative table `T∧T=T, T∧F=F, T∧U=U, F∧T=F, F∧F=F, F∧U=F, U∧T=U,
// U∧F=F, U∧U=U`.
// DOMAIN EDGE
//
// Testing Strategy row 5: the full matrix across `all` and `unless`,
// including `¬U = U` on a block-level negation. Every leg drives verdicts
// through atoms — the two operand verdicts are produced by real atoms over
// the real view, never injected.
func TestReq33_StrongKleeneMatrixAcrossAllAndUnless(t *testing.T) {
	const (
		T = resolve.GuardTrue
		F = resolve.GuardFalse
		U = resolve.GuardUnevaluable
	)

	conj := []struct {
		a, b, want resolve.GuardResult
	}{
		{T, T, T}, {T, F, F}, {T, U, U},
		{F, T, F}, {F, F, F}, {F, U, F},
		{U, T, U}, {U, F, F}, {U, U, U},
	}

	t.Run("all block conjunction", func(t *testing.T) {
		for _, c := range conj {
			t.Run(verdictName(c.a)+"∧"+verdictName(c.b), func(t *testing.T) {
				seam := newValueSeam()
				a := atomWithVerdict(seam, "a", resolve.BlockAll, c.a)
				b := atomWithVerdict(seam, "b", resolve.BlockAll, c.b)

				row := guardedRow("rdr.k3.all", "flows/rdr.toml:58", a, b)
				assertRowVerdict(t, "rev-req33-all", seam, row, c.want)
			})
		}
	})

	t.Run("unless block: verdict is ¬(unless_conj) with an empty all block", func(t *testing.T) {
		// ¬T = F, ¬F = T, ¬U = U applied to the conjunction of two
		// `unless` atoms.
		for _, c := range conj {
			t.Run("¬("+verdictName(c.a)+"∧"+verdictName(c.b)+")", func(t *testing.T) {
				seam := newValueSeam()
				a := atomWithVerdict(seam, "a", resolve.BlockUnless, c.a)
				b := atomWithVerdict(seam, "b", resolve.BlockUnless, c.b)

				row := guardedRow("rdr.k3.unless", "flows/rdr.toml:59", a, b)
				assertRowVerdict(t, "rev-req33-unless", seam, row, negate(c.want))
			})
		}
	})

	t.Run("¬U = U on a block-level negation", func(t *testing.T) {
		seam := newValueSeam()
		u := atomWithVerdict(seam, "a", resolve.BlockUnless, U)
		row := guardedRow("rdr.k3.notU", "flows/rdr.toml:60", u)
		// ¬U = U, so the row is undecidable — not selected (¬U = T) and
		// not pruned (¬U = F).
		assertRowVerdict(t, "rev-req33-notU", seam, row, U)
	})

	t.Run("cross-block: all_result ∧ ¬(unless_conj)", func(t *testing.T) {
		for _, c := range conj {
			t.Run(verdictName(c.a)+"∧¬"+verdictName(c.b), func(t *testing.T) {
				seam := newValueSeam()
				a := atomWithVerdict(seam, "a", resolve.BlockAll, c.a)
				b := atomWithVerdict(seam, "b", resolve.BlockUnless, c.b)

				row := guardedRow("rdr.k3.cross", "flows/rdr.toml:61", a, b)
				assertRowVerdict(t, "rev-req33-cross", seam, row, kleeneAnd(c.a, negate(c.b)))
			})
		}
	})
}

// REQ-34: "The kernel MUST implement combination against the tables (or
// against an explicit truth-order rank), and MUST NOT `min` the raw
// constant values; the constant order is a shipped fact this RDR does not
// renumber" — shipped `GuardResult` is `GuardFalse, GuardTrue,
// GuardUnevaluable` (`iota` 0,1,2), the order `F < T < U`, which is NOT
// the truth order.
// ADVERSARIAL
//
// Testing Strategy row 21, mutation-killing. Two legs: (a) the shipped
// constant order is unchanged — this RDR does not renumber; (b) `T ∧ U`
// is U, which integer `min` over (T=1, U=2) computes as T.
func TestReq34_CombinationIsTableDrivenNotIntegerMinOverTheConstants(t *testing.T) {
	t.Run("shipped constant order is not renumbered", func(t *testing.T) {
		if resolve.GuardFalse != 0 || resolve.GuardTrue != 1 || resolve.GuardUnevaluable != 2 {
			t.Errorf("GuardResult constants are F=%d T=%d U=%d; want the shipped iota "+
				"order 0,1,2 — renumbering is a silent behaviour change for any "+
				"existing comparison",
				resolve.GuardFalse, resolve.GuardTrue, resolve.GuardUnevaluable)
		}
	})

	t.Run("T ∧ U = U, which integer min computes as T", func(t *testing.T) {
		seam := newValueSeam()
		tAtom := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardTrue)
		uAtom := atomWithVerdict(seam, "b", resolve.BlockAll, resolve.GuardUnevaluable)

		row := guardedRow("rdr.min.tu", "flows/rdr.toml:62", tAtom, uAtom)
		// min(GuardTrue=1, GuardUnevaluable=2) = GuardTrue ⇒ a plan.
		assertRowVerdict(t, "rev-req34-tu", seam, row, resolve.GuardUnevaluable)
	})

	t.Run("U ∧ F = F, which UNEVALUABLE-dominance computes as U", func(t *testing.T) {
		seam := newValueSeam()
		uAtom := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardUnevaluable)
		fAtom := atomWithVerdict(seam, "b", resolve.BlockAll, resolve.GuardFalse)

		row := guardedRow("rdr.min.uf", "flows/rdr.toml:63", uAtom, fAtom)
		assertRowVerdict(t, "rev-req34-uf", seam, row, resolve.GuardFalse)
	})
}

// REQ-35: "The row verdict is `all_result ∧ ¬(unless_conj)`, where
// `unless_conj` is the conjunction of the `unless` block's atoms (`unless`
// is block-level negation, NOT per-atom negation)."
// DOMAIN EDGE
//
// The discriminating input: two `unless` atoms, one TRUE and one FALSE.
// Block-level negation gives ¬(T ∧ F) = ¬F = T (row selected). Per-atom
// negation would give (¬T) ∧ (¬F) = F ∧ T = F (row pruned).
func TestReq35_UnlessIsBlockLevelNegationNotPerAtom(t *testing.T) {
	seam := newValueSeam()
	tAtom := atomWithVerdict(seam, "a", resolve.BlockUnless, resolve.GuardTrue)
	fAtom := atomWithVerdict(seam, "b", resolve.BlockUnless, resolve.GuardFalse)

	row := guardedRow("rdr.unless.blocklevel", "flows/rdr.toml:64", tAtom, fAtom)
	assertRowVerdict(t, "rev-req35", seam, row, resolve.GuardTrue)
}

// REQ-36: "A guard is GuardTrue or GuardFalse only when decided atoms
// alone determine it; any unresolved dependence on an unevaluable atom
// yields GuardUnevaluable."
// DOMAIN EDGE
//
// Testing Strategy row 4 (`T ∧ U = U`), mutation-killing: an unevaluable
// atom beside all-TRUE siblings refuses and the payload names exactly that
// atom.
func TestReq36_UnresolvedDependenceOnAnUnevaluableAtomYieldsUnevaluable(t *testing.T) {
	seam := newValueSeam()
	t1 := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardTrue)
	t2 := atomWithVerdict(seam, "b", resolve.BlockAll, resolve.GuardTrue)
	u := allAtom(absentKey, opEq, "3")

	row := guardedRow("rdr.tu", "flows/rdr.toml:65", t1, t2, u)
	r := mustRefuse(t, guardInput("rev-req36", seam, row), resolve.KindGuardUnevaluable)

	atoms := undecidedAtomsFor(r, row.RuleID)
	if len(atoms) != 1 {
		t.Fatalf("payload carries %d atoms (%v); want exactly the one unevaluable "+
			"atom — the decided siblings are not undecidable", len(atoms), atoms)
	}
	if atoms[0].Key != absentKey {
		t.Errorf("payload names key %q; want %q", atoms[0].Key, absentKey)
	}
}

// REQ-37: "A present-tag atom decided FALSE still yields GuardFalse beside
// an unevaluable atom (`F ∧ U = F`): the falsity is witnessed by present
// state and holds under every resolution of the unevaluable atom — this is
// what makes D8 sound."
// ADVERSARIAL
//
// Testing Strategy rows 3 and 21: mutation-killing, observable only across
// TWO atoms. The escape row is what makes the outcome discriminating: a
// pruned row lets the escape plan, an undecidable one refuses.
func TestReq37_FalseDominatesUnevaluableAcrossTwoAtoms(t *testing.T) {
	seam := newValueSeam()
	f := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardFalse)
	u := allAtom(absentKey, opEq, "3")

	for name, atoms := range map[string][]resolve.GuardAtom{
		"F before U": {f, u},
		"U before F": {u, f},
	} {
		t.Run(name, func(t *testing.T) {
			row := guardedRow("rdr.fu", "flows/rdr.toml:66", atoms...)
			escape := conformingEscapeRow("rdr.escape.fu", "flows/rdr.toml:67", resolve.KindNoMatch)

			plan := mustPlan(t, guardInput("rev-req37", seam, row, escape))
			if plan.RuleID != escape.RuleID || !plan.Escaped {
				t.Errorf("plan names %q (escaped=%v); want the row PRUNED under "+
					"F ∧ U = F and the escape %q reached — this is what makes D8 sound",
					plan.RuleID, plan.Escaped, escape.RuleID)
			}
		})
	}
}

// REQ-58: "What the payload does NOT promise: a row pruned as GuardFalse
// under `F ∧ U = F` is not a survivor, so its absent keys are not
// reported. The guarantee is \"never a plan from absence,\" not \"every
// absence reported\"; an absence is named exactly when it blocked a
// verdict."
// DOMAIN EDGE
func TestReq58_APrunedRowReportsNothing(t *testing.T) {
	seam := newValueSeam()
	f := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardFalse)
	pruned := guardedRow("rdr.pruned", "flows/rdr.toml:68", f, allAtom(absentKey, opEq, "3"))
	survivor := guardedRow("rdr.survivor", "flows/rdr.toml:69", allAtom("gone", opEq, "1"))

	r := mustRefuse(t,
		guardInput("rev-req58", seam, pruned, survivor),
		resolve.KindGuardUnevaluable)

	if hasItem(undecidedRuleIDs(r), pruned.RuleID) {
		t.Errorf("payload names the pruned row %q; a pruned row is not a survivor "+
			"and its absent keys are not reported", pruned.RuleID)
	}
	if !hasItem(undecidedRuleIDs(r), survivor.RuleID) {
		t.Errorf("payload = %v; want the surviving undecidable row %q — an absence "+
			"is named exactly when it blocked a verdict",
			undecidedRuleIDs(r), survivor.RuleID)
	}
}

// REQ-38: "The tables are normative; the evaluator holds no part of them."
// ADVERSARIAL
//
// The checkable form: the seam is asked per atom and never handed a
// combination question. A seam that returned a whole-row verdict would
// have to be consulted once for a two-atom row; the kernel consults it
// twice, once per present-key value atom, and combines the two itself.
func TestReq38_TheEvaluatorHoldsNoPartOfTheCombinationTables(t *testing.T) {
	seam := newValueSeam()
	a := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardTrue)
	b := atomWithVerdict(seam, "b", resolve.BlockUnless, resolve.GuardFalse)

	row := guardedRow("rdr.tables", "flows/rdr.toml:70", a, b)
	mustPlan(t, guardInput("rev-req38", seam, row))

	calls := seam.seen()
	if len(calls) != 2 {
		t.Fatalf("seam saw %d calls (%v); want one per atom — combination is the "+
			"kernel's and no part of it is delegated", len(calls), calls)
	}
	// The seam is told which block each atom sits in as a fact about the
	// atom, but is never asked to negate or combine: each call carries one
	// atom and one value.
	for _, c := range calls {
		if c.Value == "" {
			t.Errorf("seam call %+v carries no present value; the seam is only ever "+
				"asked to compare a present value against a literal", c)
		}
	}
}

// REQ-39: "GATE, THEN COUNT (JDR 0001 §D2, stated as the kernel's
// ordering). Over the candidate set: guard-FALSE rows prune first (D8);
// absent owned state among survivors is reported next as
// `owned_state_unavailable`; then any surviving row whose guard is
// unevaluable vetoes the resolution as `guard_unevaluable`; only then does
// RDR 0001's exact-one count run over the rows the gate returned, and only
// then is escape reachability consulted. `no_match` and `ambiguous_match`
// sit downstream of all three gate facts."
// DOMAIN EDGE
//
// Testing Strategy row 13.
func TestReq39_GateThenCountOrdering(t *testing.T) {
	t.Run("guard-FALSE prunes before the owned sweep sees the row", func(t *testing.T) {
		seam := newValueSeam()
		f := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardFalse)
		pruned := guardedRow("rdr.order.pruned", "flows/rdr.toml:71", f)
		pruned.RequiresOwned = []string{"status", "absent-owned"}

		survivor := guardedRow("rdr.order.survivor", "flows/rdr.toml:72",
			existsAll("status", true))

		plan := mustPlan(t, guardInput("rev-req39-prune", seam, pruned, survivor))
		if plan.RuleID != survivor.RuleID {
			t.Errorf("plan names %q; want %q — the pruned row contributes no "+
				"owned-state obligation", plan.RuleID, survivor.RuleID)
		}
	})

	t.Run("owned state is reported before an unevaluable guard", func(t *testing.T) {
		row := guardedRow("rdr.order.both", "flows/rdr.toml:73",
			allAtom(absentKey, opEq, "3"))
		row.RequiresOwned = []string{"status", "absent-owned"}

		r := mustRefuse(t, guardInput("rev-req39-owned", newValueSeam(), row),
			resolve.KindOwnedStateUnavailable)
		if !hasItem(r.MissingOwned, "absent-owned") {
			t.Errorf("MissingOwned = %v; want %q", r.MissingOwned, "absent-owned")
		}
	})

	t.Run("guard_unevaluable outranks the downstream count", func(t *testing.T) {
		// Two rows would match; one is undecidable. Without the gate, the
		// count would report ambiguous_match, which RDR 0002 makes escapable.
		seam := newValueSeam()
		decided := guardedRow("rdr.order.decided", "flows/rdr.toml:74",
			existsAll("status", true))
		undecidable := guardedRow("rdr.order.undecidable", "flows/rdr.toml:75",
			allAtom(absentKey, opEq, "3"))
		escape := conformingEscapeRow("rdr.escape.order", "flows/rdr.toml:76",
			resolve.KindAmbiguousMatch)

		mustRefuse(t, guardInput("rev-req39-veto", seam, decided, undecidable, escape),
			resolve.KindGuardUnevaluable)
	})

	t.Run("no_match sits downstream of the gate", func(t *testing.T) {
		seam := newValueSeam()
		f := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardFalse)
		row := guardedRow("rdr.order.nomatch", "flows/rdr.toml:77", f)

		mustRefuse(t, guardInput("rev-req39-nomatch", seam, row), resolve.KindNoMatch)
	})
}

// REQ-40: "The escape set is gated identically, as its own row set (D5) —
// \"identically\" by DELEGATION, not by a parallel implementation:
// shipped `escapeOrRefuse` filters the escape rows and then calls the same
// `gate` function (`gate(escapes, in.Guards, view)`), which is why it needs
// no change here"
// DOMAIN EDGE
func TestReq40_EscapeSetIsGatedIdentically(t *testing.T) {
	cases := map[string]struct {
		atoms []resolve.GuardAtom
		want  resolve.RefusalKind // "" means a plan is expected
		plan  bool
	}{
		"escape guard TRUE rescues":       {[]resolve.GuardAtom{existsAll("status", true)}, "", true},
		"escape guard FALSE does not":     {[]resolve.GuardAtom{existsAll("status", false)}, resolve.KindNoMatch, false},
		"escape guard UNEVALUABLE vetoes": {[]resolve.GuardAtom{allAtom(absentKey, opEq, "3")}, resolve.KindGuardUnevaluable, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			escape := conformingEscapeRow("rdr.escape.gated", "flows/rdr.toml:78",
				resolve.KindNoMatch, tc.atoms...)

			in := noMatchInput()
			in.Table.Revision = "rev-req40"
			in.Table.Rows = append(in.Table.Rows, escape)
			in.Guards = newValueSeam()

			if tc.plan {
				plan := mustPlan(t, in)
				if !plan.Escaped {
					t.Errorf("plan from rule %q is not marked escaped", plan.RuleID)
				}
				return
			}
			mustRefuse(t, in, tc.want)
		})
	}
}

// REQ-41: "An unevaluable CANDIDATE is never masked by a decidable escape
// row (the gate returns before the escape path is reached), and an
// unevaluable ESCAPE row yields `guard_unevaluable` in place of the
// candidate-set refusal"
// ADVERSARIAL
//
// Testing Strategy row 14, both legs.
func TestReq41_EscapeSetScopingBothLegs(t *testing.T) {
	t.Run("an unevaluable candidate is not masked by a decidable escape", func(t *testing.T) {
		candidate := guardedRow("rdr.mask.candidate", "flows/rdr.toml:79",
			allAtom(absentKey, opEq, "3"))
		escape := conformingEscapeRow("rdr.mask.escape", "flows/rdr.toml:80",
			resolve.KindNoMatch, existsAll("status", true))

		r := mustRefuse(t, guardInput("rev-req41-mask", newValueSeam(), candidate, escape),
			resolve.KindGuardUnevaluable)
		if !hasItem(undecidedRuleIDs(r), candidate.RuleID) {
			t.Errorf("payload names %v; want the candidate %q — the gate returns "+
				"before the escape path is reached",
				undecidedRuleIDs(r), candidate.RuleID)
		}
	})

	t.Run("an unevaluable escape row replaces the candidate-set refusal", func(t *testing.T) {
		escape := conformingEscapeRow("rdr.mask.escape.undecidable", "flows/rdr.toml:81",
			resolve.KindNoMatch, allAtom(absentKey, opEq, "3"))

		in := noMatchInput()
		in.Table.Revision = "rev-req41-escape"
		in.Table.Rows = append(in.Table.Rows, escape)
		in.Guards = newValueSeam()

		// Q2's stated reading: the escape set's payload rides along, since
		// the escape set is gated by DELEGATION through the same gate that
		// populates the payload.
		r := mustRefuse(t, in, resolve.KindGuardUnevaluable)
		if !hasItem(undecidedRuleIDs(r), escape.RuleID) {
			t.Errorf("payload names %v; want the escape row %q",
				undecidedRuleIDs(r), escape.RuleID)
		}
	})
}

// REQ-42: "The owned-before-unevaluable precedence is pinned to shipped
// behavior, not to D8's original rationale … Authoring docs MUST NOT repeat
// the superseded rationale."
// BOUNDARY
func TestReq42_GateDocCommentDoesNotRepeatTheSupersededRationale(t *testing.T) {
	const superseded = "is frequently the reason the seam could not decide"
	src := kernelPackageSource(t)
	if hasText(src, superseded) {
		t.Errorf("the kernel source still states D8's superseded rationale (%q); "+
			"once RequiresOwned names post-guard write dependencies the two "+
			"refusals diagnose independent problems", superseded)
	}
}

// REQ-43: "Combined reporting — one refusal carrying BOTH payloads … — is
// REJECTED here" — a row failing both ways yields
// `owned_state_unavailable` only, with the owned payload only.
// DOMAIN EDGE
func TestReq43_ARowFailingBothWaysCarriesTheOwnedPayloadOnly(t *testing.T) {
	row := guardedRow("rdr.both.ways", "flows/rdr.toml:82",
		allAtom(absentKey, opEq, "3"))
	row.RequiresOwned = []string{"status", "absent-owned"}

	r := mustRefuse(t, guardInput("rev-req43", newValueSeam(), row),
		resolve.KindOwnedStateUnavailable)

	if len(r.Undecided) != 0 {
		t.Errorf("refusal carries the guard payload %v beside the owned one; one "+
			"refusal, one kind, one payload — combined reporting is rejected",
			r.Undecided)
	}
	if !hasItem(r.MissingOwned, "absent-owned") {
		t.Errorf("MissingOwned = %v; want %q", r.MissingOwned, "absent-owned")
	}
}

// REQ-44: "The kernel maps GuardUnevaluable to `guard_unevaluable`, which
// RDR 0002 excludes from escape lists. Missing artifact state MUST NOT be
// maskable behind an escapable refusal class."
// ADVERSARIAL
func TestReq44_GuardUnevaluableIsNotMaskableBehindAnEscapableClass(t *testing.T) {
	candidate := guardedRow("rdr.notmaskable", "flows/rdr.toml:83",
		allAtom(absentKey, opEq, "3"))

	// Escapes for BOTH escapable classes are modeled and both are
	// decidable. Neither may rescue the guard_unevaluable.
	noMatchEscape := conformingEscapeRow("rdr.escape.nm", "flows/rdr.toml:84",
		resolve.KindNoMatch, existsAll("status", true))
	ambigEscape := conformingEscapeRow("rdr.escape.am", "flows/rdr.toml:85",
		resolve.KindAmbiguousMatch, existsAll("status", true))

	mustRefuse(t,
		guardInput("rev-req44", newValueSeam(), candidate, noMatchEscape, ambigEscape),
		resolve.KindGuardUnevaluable)
}

// REQ-45: "The refusal MUST name what was missing, per row and per atom:
// for every undecidable row, its `(RuleID, SourceLocator)` and, for each
// of its unevaluable atoms, the referenced key, the block, and a reason
// drawn from a closed set — `absent` (the key was not in the view) or
// `uncomparable` (the key was present and its value was not compared to a
// verdict)."
// HAPPY PATH
func TestReq45_PayloadNamesEveryUndecidableRowAndAtom(t *testing.T) {
	seam := newValueSeam()
	absentAll := allAtom(absentKey, opEq, "3")
	uncomparableUnless := unlessAtom("reviews", opGte, "three")
	seam.decide(uncomparableUnless, "2", resolve.GuardUnevaluable)

	row := guardedRow("rdr.payload", "flows/rdr.toml:86", absentAll, uncomparableUnless)
	r := mustRefuse(t, guardInput("rev-req45", seam, row), resolve.KindGuardUnevaluable)

	if len(r.Undecided) != 1 {
		t.Fatalf("payload carries %d rows (%v); want exactly one", len(r.Undecided), r.Undecided)
	}
	got := r.Undecided[0]
	if got.RuleID != row.RuleID || got.SourceLocator != row.SourceLocator {
		t.Errorf("payload row identity = (%q, %q); want (%q, %q)",
			got.RuleID, got.SourceLocator, row.RuleID, row.SourceLocator)
	}
	want := []resolve.UndecidedAtom{
		{Key: absentKey, Block: resolve.BlockAll, Operator: opEq, Literal: "3", Reason: resolve.ReasonAbsent},
		{Key: "reviews", Block: resolve.BlockUnless, Operator: opGte, Literal: "three", Reason: resolve.ReasonUncomparable},
	}
	slices.SortFunc(want, compareUndecidedAtomsForTest)
	gotAtoms := slices.Clone(got.Atoms)
	slices.SortFunc(gotAtoms, compareUndecidedAtomsForTest)
	if !reflect.DeepEqual(gotAtoms, want) {
		t.Errorf("payload atoms = %+v; want %+v", gotAtoms, want)
	}
}

// REQ-46: "`uncomparable` is defined by the OUTCOME, not by which
// component produced it, and covers all three ways a present key fails to
// decide: the seam answered unevaluable (A18); the atom carries a foreign
// or missing literal on `OpExists`; or the seam is NIL, so no comparison
// could be attempted."
// DOMAIN EDGE
//
// Testing Strategy rows 12 and 22.
func TestReq46_UncomparableCoversAllThreeWaysAPresentKeyFailsToDecide(t *testing.T) {
	present := "reviews"

	t.Run("the seam answered unevaluable", func(t *testing.T) {
		atom := allAtom(present, opGte, "three")
		seam := newValueSeam().decide(atom, "2", resolve.GuardUnevaluable)
		row := guardedRow("rdr.unc.seam", "flows/rdr.toml:87", atom)

		r := mustRefuse(t, guardInput("rev-req46-seam", seam, row), resolve.KindGuardUnevaluable)
		assertReason(t, r, row.RuleID, present, resolve.BlockAll, resolve.ReasonUncomparable)
	})

	t.Run("a foreign literal on an OpExists atom", func(t *testing.T) {
		atom := resolve.GuardAtom{Key: present, Operator: resolve.OpExists, Literal: "yes", Block: resolve.BlockAll}
		row := guardedRow("rdr.unc.literal", "flows/rdr.toml:88", atom)

		r := mustRefuse(t, guardInput("rev-req46-literal", newValueSeam(), row),
			resolve.KindGuardUnevaluable)
		assertReason(t, r, row.RuleID, present, resolve.BlockAll, resolve.ReasonUncomparable)
	})

	t.Run("the seam is nil", func(t *testing.T) {
		row := guardedRow("rdr.unc.nil", "flows/rdr.toml:89", allAtom(present, opEq, "2"))

		r := mustRefuse(t, guardInput("rev-req46-nil", nil, row), resolve.KindGuardUnevaluable)
		assertReason(t, r, row.RuleID, present, resolve.BlockAll, resolve.ReasonUncomparable)
	})

	t.Run("an absent key is `absent`, never `uncomparable`", func(t *testing.T) {
		row := guardedRow("rdr.unc.absent", "flows/rdr.toml:91", allAtom(absentKey, opEq, "3"))

		r := mustRefuse(t, guardInput("rev-req46-absent", newValueSeam(), row),
			resolve.KindGuardUnevaluable)
		assertReason(t, r, row.RuleID, absentKey, resolve.BlockAll, resolve.ReasonAbsent)
	})
}

// REQ-51: "The reason values are PRODUCED BY THE VERDICT PASS, not
// re-derived: the per-atom evaluation that decides an atom unevaluable
// emits that atom's payload entry at the same step … A second walk over
// the atoms after the row is known undecidable is FORBIDDEN — it would
// consult the seam twice for every present-key atom, and the seam is not
// required to be pure or cheap."
// ADVERSARIAL
//
// The spec names the seam call count here, so asserting it is asserting
// spec-described behaviour: an undecidable row's present-key atoms must be
// handed across EXACTLY once.
func TestReq51_ReasonsAreProducedByTheVerdictPassNotASecondWalk(t *testing.T) {
	seam := newValueSeam()
	a := allAtom("status", opEq, "Draft")
	b := allAtom("reviews", opGte, "three")
	seam.decide(a, "Draft", resolve.GuardTrue)
	seam.decide(b, "2", resolve.GuardUnevaluable)

	row := guardedRow("rdr.onepass", "flows/rdr.toml:92", a, b)
	mustRefuse(t, guardInput("rev-req51", seam, row), resolve.KindGuardUnevaluable)

	counts := map[seamCall]int{}
	for _, c := range seam.seen() {
		counts[c]++
	}
	for call, n := range counts {
		if n != 1 {
			t.Errorf("the seam was consulted %d times for atom %+v; a second walk "+
				"over the atoms after the row is known undecidable is forbidden", n, call)
		}
	}
	if len(counts) != 2 {
		t.Errorf("seam saw %d distinct atoms (%v); want 2", len(counts), seam.seen())
	}
}

// REQ-52: "This payload replaces the single-valued `Refusal.Guard` text —
// and with it the shipped selection of ONE representative row (`gate`'s
// `slices.MinFunc` over the undecidable set …). Every undecidable row
// appears in the payload, so there is no representative to choose; the
// `MinFunc` call is retired, not re-typed."
// REQ-53: "`Refusal.Rows` continues to carry every undecidable row."
// DOMAIN EDGE
//
// Testing Strategy row 20.
func TestReq52and53_EveryUndecidableRowAppearsInThePayloadAndInRows(t *testing.T) {
	a := guardedRow("rdr.undecidable.a", "flows/rdr.toml:93", allAtom(absentKey, opEq, "1"))
	b := guardedRow("rdr.undecidable.b", "flows/rdr.toml:94", allAtom("gone", opEq, "2"))
	c := guardedRow("rdr.undecidable.c", "flows/rdr.toml:95", allAtom("vanished", opEq, "3"))

	r := mustRefuse(t, guardInput("rev-req52", newValueSeam(), a, b, c),
		resolve.KindGuardUnevaluable)

	wantIDs := []string{a.RuleID, b.RuleID, c.RuleID}
	if got := undecidedRuleIDs(r); !reflect.DeepEqual(got, wantIDs) {
		t.Errorf("payload names %v; want every undecidable row %v — no representative "+
			"row is chosen", got, wantIDs)
	}

	var rowIDs []string
	for _, ref := range r.Rows {
		rowIDs = append(rowIDs, ref.RuleID)
	}
	if !reflect.DeepEqual(rowIDs, wantIDs) {
		t.Errorf("Refusal.Rows names %v; want every undecidable row %v", rowIDs, wantIDs)
	}
}

// REQ-54: "The kernel MUST evaluate every atom of every survivor — no
// short-circuit on a decided block — and MUST sort the payload so it is a
// function of the input tuple (RDR 0001 REQ-1), never of atom or row
// order."
// ADVERSARIAL
func TestReq54_NoShortCircuitOnADecidedBlock(t *testing.T) {
	// The `all` block's first atom is already unevaluable, so a
	// short-circuiting kernel would stop there. Every remaining atom must
	// still be evaluated and every unevaluable one must still be reported.
	seam := newValueSeam()
	uncomparable := allAtom("reviews", opGte, "three")
	seam.decide(uncomparable, "2", resolve.GuardUnevaluable)

	row := guardedRow("rdr.noshortcircuit", "flows/rdr.toml:96",
		allAtom(absentKey, opEq, "1"),
		uncomparable,
		unlessAtom("gone", opEq, "2"),
	)

	r := mustRefuse(t, guardInput("rev-req54", seam, row), resolve.KindGuardUnevaluable)

	atoms := undecidedAtomsFor(r, row.RuleID)
	if len(atoms) != 3 {
		t.Fatalf("payload carries %d atoms (%+v); want all 3 — no short-circuit on "+
			"a decided block", len(atoms), atoms)
	}
	if calls := seam.seen(); len(calls) != 1 {
		t.Errorf("seam saw %v; the present-key atom after an already-unevaluable "+
			"sibling must still be evaluated", calls)
	}
}

// REQ-55: "The ordering is therefore the tuple `(RuleID, SourceLocator,
// key, block, operator token, literal)`, compared field by field in that
// order. Those are the row's identity plus the atom's four §D1 fields, so
// the tuple is total by construction: two entries equal on all six name
// the same atom, which the kernel MUST NOT report twice."
// ADVERSARIAL
//
// Testing Strategy row 17, including the tie case the total order exists
// for: one row carrying two atoms over ONE key, permuted in the slice.
func TestReq55_PayloadIsSortedOnTheTotalSixFieldTuple(t *testing.T) {
	t.Run("row order does not reach the payload", func(t *testing.T) {
		a := guardedRow("rdr.sort.a", "flows/rdr.toml:97", allAtom("k1", opEq, "1"))
		b := guardedRow("rdr.sort.b", "flows/rdr.toml:98", allAtom("k2", opEq, "2"))

		forward := mustRefuse(t, guardInput("rev-req55", newValueSeam(), a, b),
			resolve.KindGuardUnevaluable)
		reversed := mustRefuse(t, guardInput("rev-req55", newValueSeam(), b, a),
			resolve.KindGuardUnevaluable)

		if !reflect.DeepEqual(forward.Undecided, reversed.Undecided) {
			t.Errorf("payload depends on table row order:\n forward = %+v\nreversed = %+v",
				forward.Undecided, reversed.Undecided)
		}
	})

	t.Run("atom order does not reach the payload", func(t *testing.T) {
		x := allAtom("k", opEq, "1")
		y := unlessAtom("k", opEq, "2")
		z := allAtom("k", opGte, "1")

		forward := mustRefuse(t,
			guardInput("rev-req55-atoms", newValueSeam(),
				guardedRow("rdr.sort.atoms", "flows/rdr.toml:99", x, y, z)),
			resolve.KindGuardUnevaluable)
		permuted := mustRefuse(t,
			guardInput("rev-req55-atoms", newValueSeam(),
				guardedRow("rdr.sort.atoms", "flows/rdr.toml:99", z, y, x)),
			resolve.KindGuardUnevaluable)

		if !reflect.DeepEqual(forward.Undecided, permuted.Undecided) {
			t.Errorf("payload depends on atom order — the tie case the total order "+
				"exists for:\n forward = %+v\npermuted = %+v",
				forward.Undecided, permuted.Undecided)
		}
	})

	t.Run("entries are ordered by the six-field tuple", func(t *testing.T) {
		row := guardedRow("rdr.sort.tuple", "flows/rdr.toml:100",
			allAtom("b", opEq, "1"),
			unlessAtom("a", opEq, "1"),
			allAtom("a", opEq, "2"),
			allAtom("a", opEq, "1"),
		)
		r := mustRefuse(t, guardInput("rev-req55-tuple", newValueSeam(), row),
			resolve.KindGuardUnevaluable)

		got := undecidedAtomsFor(r, row.RuleID)
		want := slices.Clone(got)
		slices.SortFunc(want, compareUndecidedAtomsForTest)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("payload atoms = %+v; want them ordered by (key, block, "+
				"operator, literal) = %+v", got, want)
		}
	})
}

// REQ-56: "the kernel MUST compare the literal as the exact bytes the
// normalizer emitted" — the set-literal byte form is §D13's canonical JSON
// array: members sorted, duplicate-free, compact encoding.
// ADVERSARIAL
//
// Two `in` atoms on one key differing only by literal-set content are
// distinguishable by the six-tuple exactly because the literal is compared
// as bytes and §D13 makes those bytes canonical.
func TestReq56_LiteralIsComparedAsExactBytesInTheD13Form(t *testing.T) {
	first := d13Set("alpha", "beta")
	second := d13Set("beta", "gamma")

	row := guardedRow("rdr.literal.bytes", "flows/rdr.toml:101",
		allAtom("labels", opIn, second),
		allAtom("labels", opIn, first),
	)
	r := mustRefuse(t, guardInput("rev-req56", newValueSeam(), row),
		resolve.KindGuardUnevaluable)

	atoms := undecidedAtomsFor(r, row.RuleID)
	if len(atoms) != 2 {
		t.Fatalf("payload carries %d atoms (%+v); two atoms differing only in the "+
			"literal are distinct atoms", len(atoms), atoms)
	}
	if atoms[0].Literal != first || atoms[1].Literal != second {
		t.Errorf("payload literals = [%q %q]; want [%q %q] — compared as exact "+
			"bytes and ordered by them", atoms[0].Literal, atoms[1].Literal, first, second)
	}

	// §D13's canonical form is what the normalizer emits, and the kernel
	// carries it through unchanged.
	if first != `["alpha","beta"]` {
		t.Errorf("§D13 form = %q; want a sorted, duplicate-free, compact JSON array", first)
	}
}

// REQ-57: "\"MUST NOT report twice\" is a kernel obligation, not a
// property of the sort: the kernel emits one entry per unevaluable atom,
// and an atom is identified by the six-tuple, so duplicate suppression is
// by construction of the emit loop — the sort orders entries, it does not
// dedupe them."
// ADVERSARIAL
//
// The discriminating input: a row whose slice carries the SAME atom twice.
// A kernel that emitted per slice element and relied on the sort to tidy up
// would report it twice.
func TestReq57_TheKernelReportsEachAtomExactlyOnce(t *testing.T) {
	atom := allAtom(absentKey, opEq, "3")
	row := guardedRow("rdr.dup", "flows/rdr.toml:102", atom, atom)

	r := mustRefuse(t, guardInput("rev-req57", newValueSeam(), row),
		resolve.KindGuardUnevaluable)

	atoms := undecidedAtomsFor(r, row.RuleID)
	if len(atoms) != 1 {
		t.Errorf("payload carries %d entries (%+v); two entries equal on all six "+
			"fields name the SAME atom, which the kernel must not report twice",
			len(atoms), atoms)
	}
}

// REQ-64: "SURVIVOR MEMBERSHIP. A row whose guard is GuardTrue or
// GuardUnevaluable is a SURVIVOR; only GuardFalse rows are pruned. The
// owned-state scan runs over survivors, so an unevaluable row's
// RequiresOwned keys DO raise owned_state_unavailable, and the undecidable
// check runs over the same set."
// DOMAIN EDGE
func TestReq64_SurvivorMembershipIncludesUnevaluableRows(t *testing.T) {
	t.Run("an unevaluable row's RequiresOwned raises owned_state_unavailable", func(t *testing.T) {
		row := guardedRow("rdr.survivor.owned", "flows/rdr.toml:103",
			allAtom(absentKey, opEq, "3"))
		row.RequiresOwned = []string{"status", "unread-owned"}

		r := mustRefuse(t, guardInput("rev-req64", newValueSeam(), row),
			resolve.KindOwnedStateUnavailable)
		if !hasItem(r.MissingOwned, "unread-owned") {
			t.Errorf("MissingOwned = %v; want %q — an unevaluable row IS a survivor",
				r.MissingOwned, "unread-owned")
		}
	})

	t.Run("a pruned row's RequiresOwned does not", func(t *testing.T) {
		seam := newValueSeam()
		f := atomWithVerdict(seam, "a", resolve.BlockAll, resolve.GuardFalse)
		pruned := guardedRow("rdr.pruned.owned", "flows/rdr.toml:104", f)
		pruned.RequiresOwned = []string{"status", "unread-owned"}
		survivor := guardedRow("rdr.survivor.plain", "flows/rdr.toml:105",
			existsAll("status", true))

		plan := mustPlan(t, guardInput("rev-req64-pruned", seam, pruned, survivor))
		if plan.RuleID != survivor.RuleID {
			t.Errorf("plan names %q; want %q", plan.RuleID, survivor.RuleID)
		}
	})
}

// REQ-65: "The owned scan runs ONCE over the whole survivor set, and
// `MissingOwned` is the deduplicated, sorted union of the missing keys
// across survivors"
// BOUNDARY
//
// Testing Strategy row 23.
func TestReq65_MissingOwnedIsTheDeduplicatedSortedUnionAcrossSurvivors(t *testing.T) {
	a := guardedRow("rdr.owned.a", "flows/rdr.toml:106", existsAll("status", true))
	a.RequiresOwned = []string{"zulu", "shared"}
	b := guardedRow("rdr.owned.b", "flows/rdr.toml:107", existsAll("status", true))
	b.RequiresOwned = []string{"shared", "alpha"}

	forward := mustRefuse(t, guardInput("rev-req65", newValueSeam(), a, b),
		resolve.KindOwnedStateUnavailable)
	reversed := mustRefuse(t, guardInput("rev-req65", newValueSeam(), b, a),
		resolve.KindOwnedStateUnavailable)

	want := []string{"alpha", "shared", "zulu"}
	if !reflect.DeepEqual(forward.MissingOwned, want) {
		t.Errorf("MissingOwned = %v; want the deduplicated, sorted union %v",
			forward.MissingOwned, want)
	}
	if !reflect.DeepEqual(forward.MissingOwned, reversed.MissingOwned) {
		t.Errorf("MissingOwned depends on row order: %v vs %v",
			forward.MissingOwned, reversed.MissingOwned)
	}
}

// REQ-66: "Aggregation is resolution-level: if any surviving candidate
// row's guard is GuardUnevaluable, the resolution MUST refuse
// `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected
// while an unevaluable candidate exists."
// ADVERSARIAL
//
// Testing Strategy row 6; MVV scenario 3.
func TestReq66_AnUnevaluableCandidateVetoesADecidedTrueSibling(t *testing.T) {
	undecidable := guardedRow("rdr.veto.a", "flows/rdr.toml:108",
		allAtom(absentKey, opEq, "3"))
	decided := guardedRow("rdr.veto.b", "flows/rdr.toml:109",
		existsAll("status", true))

	r := mustRefuse(t, guardInput("rev-req66", newValueSeam(), undecidable, decided),
		resolve.KindGuardUnevaluable)

	if !hasItem(undecidedRuleIDs(r), undecidable.RuleID) {
		t.Errorf("payload names %v; want the unevaluable row %q",
			undecidedRuleIDs(r), undecidable.RuleID)
	}
	if hasItem(undecidedRuleIDs(r), decided.RuleID) {
		t.Errorf("payload names the decided-TRUE row %q; only undecidable rows "+
			"appear", decided.RuleID)
	}
}

// REQ-67: "Deviation D8 (guard-FALSE prunes first; a pruned row
// contributes neither candidacy nor an owned-state obligation) is ratified
// as normative, conditional on the domain rule: pruning is safe exactly
// because GuardFalse can only arise from decided atoms — value comparisons
// over present keys, or existence atoms — never from absence folding into
// a value comparison."
// DOMAIN EDGE
//
// Testing Strategy row 2. The conditionality is the testable half: the
// same row shape prunes when the key is PRESENT and the value decides
// FALSE, and refuses when the key is ABSENT.
func TestReq67_D8IsRatifiedConditionalOnTheDomainRule(t *testing.T) {
	escape := conformingEscapeRow("rdr.escape.d8", "flows/rdr.toml:110", resolve.KindNoMatch)

	t.Run("present key, value decided FALSE ⇒ pruned, escape plans", func(t *testing.T) {
		row := guardedRow("rdr.d8", "flows/rdr.toml:111", allAtom("status", opEq, "Final"))
		plan := mustPlan(t, guardInput("rev-req67-false", eqSeam{}, row, escape))
		if plan.RuleID != escape.RuleID || !plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the escape %q — decided "+
				"falsity prunes", plan.RuleID, plan.Escaped, escape.RuleID)
		}
	})

	t.Run("absent key ⇒ refuses, never pruned", func(t *testing.T) {
		row := guardedRow("rdr.d8", "flows/rdr.toml:111", allAtom(absentKey, opEq, "Final"))
		mustRefuse(t, guardInput("rev-req67-absent", eqSeam{}, row, escape),
			resolve.KindGuardUnevaluable)
	})

	t.Run("existence atom deciding FALSE ⇒ pruned", func(t *testing.T) {
		row := guardedRow("rdr.d8", "flows/rdr.toml:111", existsAll(absentKey, true))
		plan := mustPlan(t, guardInput("rev-req67-exists", eqSeam{}, row, escape))
		if plan.RuleID != escape.RuleID {
			t.Errorf("plan names %q; want the escape %q", plan.RuleID, escape.RuleID)
		}
	})
}

// REQ-70: "encode the domain-rule matrix as kernel tests (operators ×
// present/absent × `all`/`unless` × {T,F,U} combinations, empty blocks,
// two-row pattern, escape-set scoping, foreign existence token)"
// DOMAIN EDGE
//
// The matrix axis this test carries is operators × present/absent ×
// all/unless; the {T,F,U} combinations are TestReq33's, empty blocks are
// TestReq30to32's, the two-row pattern is TestReq76's, escape-set scoping
// is TestReq41's, and the foreign token is TestReq24's.
func TestReq70_DomainRuleMatrixOverOperatorsPresenceAndBlocks(t *testing.T) {
	for _, op := range []string{opEq, opIn, opGte, opContains} {
		for _, block := range []resolve.Block{resolve.BlockAll, resolve.BlockUnless} {
			t.Run(op+"/"+string(block)+"/absent", func(t *testing.T) {
				seam := newValueSeam()
				atom := resolve.GuardAtom{Key: absentKey, Operator: op, Literal: "x", Block: block}
				row := guardedRow("rdr.matrix", "flows/rdr.toml:112", atom)

				r := mustRefuse(t, guardInput("rev-req70", seam, row),
					resolve.KindGuardUnevaluable)
				if calls := seam.seen(); len(calls) != 0 {
					t.Errorf("seam saw %v for an absent key in block %q", calls, block)
				}
				assertReason(t, r, row.RuleID, absentKey, block, resolve.ReasonAbsent)
			})

			t.Run(op+"/"+string(block)+"/present", func(t *testing.T) {
				seam := newValueSeam()
				atom := resolve.GuardAtom{Key: "reviews", Operator: op, Literal: "x", Block: block}
				row := guardedRow("rdr.matrix", "flows/rdr.toml:112", atom)

				mustResolve(t, guardInput("rev-req70-present", seam, row))
				if calls := seam.seen(); len(calls) != 1 || calls[0] != callOf(atom, "2") {
					t.Errorf("seam saw %v; a value atom over a present key must be "+
						"handed across in block %q too", calls, block)
				}
			})
		}
	}
}

// --- helpers -------------------------------------------------------------

// atomWithVerdict builds a value atom over a distinct PRESENT key whose
// verdict the seam will report as want, or an atom the KERNEL decides
// unevaluable on absence. Either way the verdict is produced through the
// atom pipeline, never injected at row level (premortem P-20).
func atomWithVerdict(seam *valueSeam, slot string, block resolve.Block, want resolve.GuardResult) resolve.GuardAtom {
	if want == resolve.GuardUnevaluable {
		// Kernel-decided unevaluable: an absent key. Distinct per slot so
		// two U atoms in one row are distinct atoms.
		return resolve.GuardAtom{
			Key: "absent-" + slot, Operator: opEq, Literal: "x", Block: block,
		}
	}
	// Seam-decided: a value atom over the present key `status`, made
	// distinct per slot by its literal.
	atom := resolve.GuardAtom{
		Key: "status", Operator: opEq, Literal: "lit-" + slot, Block: block,
	}
	seam.decide(atom, "Draft", want)
	return atom
}

// assertRowVerdict resolves a table holding row plus a modeled no_match
// escape, and asserts the row's guard verdict by its OBSERVABLE
// consequence: TRUE selects the row, FALSE prunes it and the escape plans,
// UNEVALUABLE refuses guard_unevaluable.
func assertRowVerdict(t *testing.T, revision string, seam resolve.GuardEvaluator, row resolve.Row, want resolve.GuardResult) {
	t.Helper()
	escape := conformingEscapeRow("rdr.escape.verdict", "flows/rdr.toml:113", resolve.KindNoMatch)
	got := mustResolve(t, guardInput(revision, seam, row, escape))

	switch want {
	case resolve.GuardTrue:
		if got.Refused() {
			t.Fatalf("row verdict refused %q; want TRUE (the row is selected)", got.Refusal.Kind)
		}
		if got.Plan.RuleID != row.RuleID || got.Plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the row %q selected",
				got.Plan.RuleID, got.Plan.Escaped, row.RuleID)
		}
	case resolve.GuardFalse:
		if got.Refused() {
			t.Fatalf("row verdict refused %q; want FALSE (the row is pruned and the "+
				"escape plans)", got.Refusal.Kind)
		}
		if got.Plan.RuleID != escape.RuleID || !got.Plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the escape %q reached",
				got.Plan.RuleID, got.Plan.Escaped, escape.RuleID)
		}
	case resolve.GuardUnevaluable:
		if !got.Refused() {
			t.Fatalf("plan named %q; want UNEVALUABLE (guard_unevaluable)", got.Plan.RuleID)
		}
		if got.Refusal.Kind != resolve.KindGuardUnevaluable {
			t.Errorf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindGuardUnevaluable)
		}
	}
}

func assertReason(t *testing.T, r *resolve.Refusal, ruleID, key string, block resolve.Block, want resolve.Reason) {
	t.Helper()
	got, ok := reasonFor(undecidedAtomsFor(r, ruleID), key, block)
	if !ok {
		t.Errorf("payload for %q has no entry for key %q in block %q; the payload "+
			"never omits an unevaluable atom", ruleID, key, block)
		return
	}
	if got != want {
		t.Errorf("reason for %q/%q in block %q = %q; want %q", ruleID, key, block, got, want)
	}
}

// kleeneAnd and negate are the normative tables, restated in the test so a
// kernel that got them wrong cannot supply its own oracle.
func kleeneAnd(a, b resolve.GuardResult) resolve.GuardResult {
	if a == resolve.GuardFalse || b == resolve.GuardFalse {
		return resolve.GuardFalse
	}
	if a == resolve.GuardUnevaluable || b == resolve.GuardUnevaluable {
		return resolve.GuardUnevaluable
	}
	return resolve.GuardTrue
}

func negate(v resolve.GuardResult) resolve.GuardResult {
	switch v {
	case resolve.GuardTrue:
		return resolve.GuardFalse
	case resolve.GuardFalse:
		return resolve.GuardTrue
	default:
		return resolve.GuardUnevaluable
	}
}

func verdictName(v resolve.GuardResult) string {
	switch v {
	case resolve.GuardTrue:
		return "T"
	case resolve.GuardFalse:
		return "F"
	default:
		return "U"
	}
}

// compareUndecidedAtomsForTest orders payload atoms by the four §D1 fields
// of REQ-55's tuple, restated in the test so the kernel cannot supply its
// own oracle for its own sort.
func compareUndecidedAtomsForTest(a, b resolve.UndecidedAtom) int {
	if c := compareStrings(a.Key, b.Key); c != 0 {
		return c
	}
	if c := compareStrings(string(a.Block), string(b.Block)); c != 0 {
		return c
	}
	if c := compareStrings(a.Operator, b.Operator); c != 0 {
		return c
	}
	return compareStrings(a.Literal, b.Literal)
}

func compareStrings(a, b string) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
