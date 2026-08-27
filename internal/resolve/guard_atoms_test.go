package resolve_test

// RDR 0007 — seam shape, atom transport, the per-atom nil-seam rule, and
// the kernel-owned constant sets.
//
// Every scenario here is a kernel test in package resolve_test driven
// through atoms and a value stub (REQ-73). The kernel under test is always
// the real resolve.Resolve.

import (
	"go/ast"
	"reflect"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-1: "SEAM. A candidate row carries its guard as a slice of parsed
// atoms — key, operator token, literal, block ∈ {all, unless} — the shape
// JDR 0001 §D1 fixes; this RDR cites it and does not restate the grammar."
// HAPPY PATH
//
// The checkable content of "carries its guard as a slice of parsed atoms"
// is that a row holding two atoms is evaluated as two atoms: the kernel
// must reach BOTH, which a one-slot opaque channel cannot express.
func TestReq1_RowCarriesItsGuardAsASliceOfParsedAtoms(t *testing.T) {
	first := allAtom("status", opEq, "Draft")
	second := allAtom("reviews", opEq, "2")

	seam := newValueSeam().
		decide(first, "Draft", resolve.GuardTrue).
		decide(second, "2", resolve.GuardTrue)

	row := guardedRow("rdr.two.atoms", "flows/rdr.toml:10", first, second)
	plan := mustPlan(t, guardInput("rev-req1", seam, row))

	if plan.RuleID != row.RuleID {
		t.Errorf("plan names rule %q; want %q", plan.RuleID, row.RuleID)
	}

	want := []seamCall{callOf(first, "Draft"), callOf(second, "2")}
	got := seam.seen()
	if len(got) != len(want) {
		t.Fatalf("seam saw %d atoms (%v); want %d — the row's guard is a "+
			"slice of atoms, each carried across the seam on its own", len(got), got, len(want))
	}
	for _, w := range want {
		if !containsCall(got, w) {
			t.Errorf("seam never saw atom %+v; the kernel must carry every parsed atom", w)
		}
	}
}

// REQ-2: "An empty slice is an unguarded row."
// INPUT EDGE
func TestReq2_EmptyAtomSliceIsAnUnguardedRow(t *testing.T) {
	for name, guard := range map[string][]resolve.GuardAtom{
		"nil slice":   nil,
		"empty slice": {},
	} {
		t.Run(name, func(t *testing.T) {
			row := guardedRow("rdr.unguarded", "flows/rdr.toml:11")
			row.Guard = guard

			// A nil seam proves the kernel never reached for one: an
			// unguarded row has no atom to hand across.
			plan := mustPlan(t, guardInput("rev-req2", nil, row))
			if plan.RuleID != row.RuleID {
				t.Errorf("plan names rule %q; want %q", plan.RuleID, row.RuleID)
			}
		})
	}
}

// REQ-3: "There is no opaque guard string and no reconstruction step, so
// no mapping failure exists and no panic or error channel is needed for
// one."
// BOUNDARY
//
// The checkable form of "no opaque guard string" is that the kernel's
// guard-carrying field is not a string, and that the seam's method takes
// no view. Both are structural claims about the exported surface, which is
// where a reconstruction step would have to reappear.
func TestReq3_NoOpaqueGuardStringAndNoReconstructionStep(t *testing.T) {
	guard, ok := reflect.TypeOf(resolve.Row{}).FieldByName("Guard")
	if !ok {
		t.Fatal("resolve.Row has no Guard field")
	}
	if guard.Type.Kind() != reflect.Slice {
		t.Fatalf("Row.Guard is %s; want a slice of parsed atoms — an opaque "+
			"string reintroduces the reconstruction step §D1 removes", guard.Type)
	}
	if guard.Type.Elem() != reflect.TypeOf(resolve.GuardAtom{}) {
		t.Errorf("Row.Guard elements are %s; want resolve.GuardAtom", guard.Type.Elem())
	}

	// No panic channel: a guard the seam cannot decide is a verdict, not a
	// crash, and the Go error return stays reserved for programmer error.
	row := guardedRow("rdr.nopanic", "flows/rdr.toml:12", allAtom("reviews", opEq, "2"))
	in := guardInput("rev-req3", newValueSeam(), row)
	got, err := resolve.Resolve(in)
	if err != nil {
		t.Fatalf("Resolve used the Go error path for an undecidable guard: %v", err)
	}
	if !got.Refused() || got.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("undecidable guard produced %+v; want a guard_unevaluable refusal", got)
	}
}

// REQ-4: "`Refusal.Guard` has no referent once guards are atoms" — the
// shipped `Row.Guard string` and `Refusal.Guard string` fields are
// DELETED, not re-typed.
// BOUNDARY
func TestReq4_RefusalGuardFieldIsDeleted(t *testing.T) {
	if _, ok := reflect.TypeOf(resolve.Refusal{}).FieldByName("Guard"); ok {
		t.Error("resolve.Refusal still carries a Guard field; it has no referent " +
			"once guards are atoms and must be deleted, not re-typed")
	}
}

// REQ-5: "The atom's four fields are named `Key`, `Operator`, `Literal`,
// `Block` — the payload's `UndecidedAtom` mirrors them by name, so
// divergent spellings would make the mirror a mapping."
// BOUNDARY
func TestReq5_AtomAndPayloadAtomShareTheFourFieldNames(t *testing.T) {
	atom := reflect.TypeOf(resolve.GuardAtom{})
	payload := reflect.TypeOf(resolve.UndecidedAtom{})

	for _, name := range []string{"Key", "Operator", "Literal", "Block"} {
		a, aok := atom.FieldByName(name)
		if !aok {
			t.Errorf("resolve.GuardAtom has no field %q", name)
			continue
		}
		p, pok := payload.FieldByName(name)
		if !pok {
			t.Errorf("resolve.UndecidedAtom has no field %q; the payload mirrors "+
				"the atom by name, and a divergent spelling makes the mirror a mapping", name)
			continue
		}
		if a.Type != p.Type {
			t.Errorf("field %q is %s on GuardAtom and %s on UndecidedAtom; the mirror "+
				"is by name AND type", name, a.Type, p.Type)
		}
	}
	if atom.NumField() != 4 {
		t.Errorf("resolve.GuardAtom has %d fields; the atom's four fields are "+
			"Key, Operator, Literal, Block", atom.NumField())
	}
}

// REQ-6: "`Block` is an exported named STRING type carrying exactly two
// constants, `BlockAll = \"all\"` and `BlockUnless = \"unless\"`."
// BOUNDARY
//
// JDR 0001 §D12 (§JD-21) widens this fence: `Block` gains a THIRD
// constant, `BlockMatch`, added in this RDR's Phase 1 — one type, no
// separate slice on Row. See artifacts/deviations.md D1. The obligation
// asserted here is the widened form: an exported named string type
// carrying all three constants.
func TestReq6_BlockIsAnExportedNamedStringTypeWithThreeConstants(t *testing.T) {
	block := reflect.TypeOf(resolve.BlockAll)
	if block.Kind() != reflect.String {
		t.Fatalf("resolve.Block's underlying kind is %s; want string", block.Kind())
	}
	if block.Name() != "Block" || block.PkgPath() == "" {
		t.Fatalf("Block constants have type %q (pkg %q); want the exported named "+
			"type resolve.Block", block.Name(), block.PkgPath())
	}
	if block == reflect.TypeOf("") {
		t.Error("Block is an alias for string; the payload clause requires a named type")
	}

	for name, got := range map[string]struct{ have, want resolve.Block }{
		"BlockAll":    {resolve.BlockAll, "all"},
		"BlockUnless": {resolve.BlockUnless, "unless"},
		// §D12: the third member. Its byte value is "match" by symmetry
		// (artifacts/req-list.md ASSUMPTION).
		"BlockMatch": {resolve.BlockMatch, "match"},
	} {
		if got.have != got.want {
			t.Errorf("%s = %q; want %q", name, got.have, got.want)
		}
		if reflect.TypeOf(got.have) != block {
			t.Errorf("%s is not of type resolve.Block", name)
		}
	}
}

// REQ-7: "The kernel evaluates a guard ATOM BY ATOM. For each atom it MUST
// first decide presence of the referenced key in the assembled view,
// provenance-blind (the `TagSet.Lookup` `ok` test)."
// HAPPY PATH
func TestReq7_KernelDecidesPresencePerAtomBeforeConsultingTheSeam(t *testing.T) {
	present := allAtom("reviews", opEq, "2")
	absent := allAtom(absentKey, opEq, "3")

	seam := newValueSeam().decide(present, "2", resolve.GuardTrue)
	row := guardedRow("rdr.presence", "flows/rdr.toml:13", present, absent)

	r := mustRefuse(t, guardInput("rev-req7", seam, row), resolve.KindGuardUnevaluable)

	// Presence decided per atom: the present-key atom crossed the seam
	// with its value, the absent-key atom did not cross at all.
	if calls := seam.seen(); len(calls) != 1 || calls[0] != callOf(present, "2") {
		t.Errorf("seam saw %v; want exactly the present-key atom carrying its value", calls)
	}
	atoms := undecidedAtomsFor(r, row.RuleID)
	if reason, ok := reasonFor(atoms, absentKey, resolve.BlockAll); !ok || reason != resolve.ReasonAbsent {
		t.Errorf("payload reports %q for the absent key (found=%v); want %q",
			reason, ok, resolve.ReasonAbsent)
	}
}

// REQ-8: "an existence atom is decided by the KERNEL from presence alone,
// and the evaluator MUST NOT be consulted for it"
// HAPPY PATH
func TestReq8_ExistenceAtomsNeverReachTheEvaluator(t *testing.T) {
	cases := map[string]struct {
		atom resolve.GuardAtom
		want bool // want a plan
	}{
		"exists=true over a present key":  {existsAll("reviews", true), true},
		"exists=false over an absent key": {existsAll(absentKey, false), true},
		"exists=true over an absent key":  {existsAll(absentKey, true), false},
		"exists=false over a present key": {existsAll("reviews", false), false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			seam := newValueSeam()
			row := guardedRow("rdr.exists", "flows/rdr.toml:14", tc.atom)
			got := mustResolve(t, guardInput("rev-req8", seam, row))

			if calls := seam.seen(); len(calls) != 0 {
				t.Errorf("evaluator was consulted for an existence atom: %v; the "+
					"kernel decides it from presence alone", calls)
			}
			if tc.want && got.Refused() {
				t.Fatalf("kernel refused %q; the existence atom decides TRUE", got.Refusal.Kind)
			}
			if !tc.want && !got.Refused() {
				t.Fatalf("kernel planned rule %q; the existence atom decides FALSE, "+
					"so the row is pruned and no candidate remains", got.Plan.RuleID)
			}
		})
	}
}

// REQ-9: "a value-comparing atom whose key is ABSENT is marked unevaluable
// by the KERNEL, and the evaluator MUST NOT be consulted for it"
// DOMAIN EDGE
func TestReq9_AbsentKeyValueAtomIsKernelUnevaluableAndNeverReachesTheSeam(t *testing.T) {
	for _, op := range []string{opEq, opIn, opGte, opContains} {
		t.Run(op, func(t *testing.T) {
			seam := newValueSeam()
			atom := allAtom(absentKey, op, "3")
			row := guardedRow("rdr.absent", "flows/rdr.toml:15", atom)

			r := mustRefuse(t, guardInput("rev-req9", seam, row), resolve.KindGuardUnevaluable)

			if calls := seam.seen(); len(calls) != 0 {
				t.Errorf("evaluator was consulted for an absent key: %v; an evaluator "+
					"that is never asked cannot fold absence into false", calls)
			}
			reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), absentKey, resolve.BlockAll)
			if !ok || reason != resolve.ReasonAbsent {
				t.Errorf("payload reason = %q (found=%v); want %q", reason, ok, resolve.ReasonAbsent)
			}
		})
	}
}

// REQ-10: "a value-comparing atom whose key is PRESENT is handed to the
// evaluator seam with the atom and the present value:
// `Evaluate(atom, value) GuardResult` — a single-method INTERFACE named
// `GuardEvaluator`, not a func type, so the nil-seam rule above is a nil
// interface value and RDR 0003's evaluator satisfies it by declaring the
// method."
// BOUNDARY
func TestReq10_GuardEvaluatorIsASingleMethodInterfaceTakingAtomAndValue(t *testing.T) {
	iface := reflect.TypeOf((*resolve.GuardEvaluator)(nil)).Elem()
	if iface.Kind() != reflect.Interface {
		t.Fatalf("resolve.GuardEvaluator is %s; want an interface, not a func type", iface.Kind())
	}
	if iface.NumMethod() != 1 {
		t.Fatalf("GuardEvaluator declares %d methods; want exactly one", iface.NumMethod())
	}
	m := iface.Method(0)
	if m.Name != "Evaluate" {
		t.Errorf("method is %q; want Evaluate", m.Name)
	}
	sig := m.Type
	if sig.NumIn() != 2 || sig.NumOut() != 1 {
		t.Fatalf("Evaluate has signature %s; want Evaluate(GuardAtom, string) GuardResult", sig)
	}
	if sig.In(0) != reflect.TypeOf(resolve.GuardAtom{}) {
		t.Errorf("Evaluate's first parameter is %s; want resolve.GuardAtom", sig.In(0))
	}
	if sig.In(1).Kind() != reflect.String {
		t.Errorf("Evaluate's second parameter is %s; want the present value as a string", sig.In(1))
	}
	if sig.Out(0) != reflect.TypeOf(resolve.GuardResult(0)) {
		t.Errorf("Evaluate returns %s; want resolve.GuardResult", sig.Out(0))
	}
	// The seam never sees the view: no parameter is a TagSet.
	for i := range sig.NumIn() {
		if sig.In(i) == reflect.TypeOf(resolve.TagSet{}) {
			t.Error("Evaluate takes the assembled view; the seam must never see it")
		}
	}

	// And the value handed across is the PRESENT value from the view.
	atom := allAtom("status", opEq, "Draft")
	seam := newValueSeam().decide(atom, "Draft", resolve.GuardTrue)
	mustPlan(t, guardInput("rev-req10", seam, guardedRow("rdr.seam", "flows/rdr.toml:16", atom)))
	if calls := seam.seen(); len(calls) != 1 || calls[0].Value != "Draft" {
		t.Errorf("seam saw %v; want the atom carrying the present value %q", calls, "Draft")
	}
}

// REQ-11: "The seam decides true or false under RDR 0003's typed operator
// semantics, and MAY answer unevaluable for a present value it cannot
// compare (A18). It never sees the view."
// DOMAIN EDGE
func TestReq11_SeamMayAnswerUnevaluableForAPresentValue(t *testing.T) {
	atom := allAtom("reviews", opGte, "three")
	seam := newValueSeam().decide(atom, "2", resolve.GuardUnevaluable)
	row := guardedRow("rdr.a18", "flows/rdr.toml:17", atom)

	r := mustRefuse(t, guardInput("rev-req11", seam, row), resolve.KindGuardUnevaluable)

	reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), "reviews", resolve.BlockAll)
	if !ok || reason != resolve.ReasonUncomparable {
		t.Errorf("payload reason = %q (found=%v); want %q — the key was present, so "+
			"`absent` would be a lie", reason, ok, resolve.ReasonUncomparable)
	}
	if calls := seam.seen(); len(calls) != 1 {
		t.Errorf("seam saw %v; want exactly one call for the present-key value atom", calls)
	}
}

// REQ-12: "The kernel therefore enforces the domain rule by structure: an
// evaluator cannot fold absence into false, because it is never asked
// about an absent key."
// ADVERSARIAL
//
// The adversarial seam answers GuardFalse to everything — precisely the
// absence-as-false fold this RDR exists to make unreachable. A kernel that
// asks it about an absent key would prune the row and reach an escapable
// no_match; a kernel that never asks refuses guard_unevaluable.
func TestReq12_AnAbsenceFoldingEvaluatorCannotBeReachedForAnAbsentKey(t *testing.T) {
	row := guardedRow("rdr.hostile", "flows/rdr.toml:18", allAtom(absentKey, opEq, "3"))
	escape := conformingEscapeRow("rdr.escape.hostile", "flows/rdr.toml:19", resolve.KindNoMatch)

	in := guardInput("rev-req12", falseSeam{}, row, escape)
	r := mustRefuse(t, in, resolve.KindGuardUnevaluable)

	if len(r.Undecided) == 0 {
		t.Error("refusal carries no payload; the absent key must be named")
	}
}

// falseSeam folds every question into GuardFalse. The kernel must never
// reach it for an absent key.
type falseSeam struct{}

func (falseSeam) Evaluate(resolve.GuardAtom, string) resolve.GuardResult { return resolve.GuardFalse }

// REQ-13: "A NIL seam does not make a row unevaluable by itself. The
// kernel decides existence atoms and absent-key atoms without consulting
// the seam, so a row whose atoms are all kernel-decidable resolves under a
// nil seam exactly as it would under a present one."
// BOUNDARY
//
// This is Testing Strategy row 19's discriminating leg: shipped
// evaluateGuard returns GuardUnevaluable for ANY non-empty guard when
// seam == nil. Here a pure-existence row must produce a PLAN and land its
// Writes.
func TestReq13_NilSeamStillResolvesAKernelDecidableRow(t *testing.T) {
	row := guardedRow("rdr.nilseam.exists", "flows/rdr.toml:20",
		existsAll("status", true),
		existsAll(absentKey, false),
	)

	withNil := mustPlan(t, guardInput("rev-req13-nil", nil, row))
	withSeam := mustPlan(t, guardInput("rev-req13-nil", newValueSeam(), row))

	if !reflect.DeepEqual(withNil, withSeam) {
		t.Errorf("nil seam plan %+v differs from present-seam plan %+v; a row whose "+
			"atoms are all kernel-decidable resolves identically either way", withNil, withSeam)
	}
	if !reflect.DeepEqual(withNil.Writes, row.Writes) {
		t.Errorf("plan writes = %v; want %v — the nil-seam row's writes must land",
			withNil.Writes, row.Writes)
	}
}

// REQ-14: "A nil seam yields unevaluable only for an atom that WOULD have
// been handed to it — a value-comparing atom over a present key."
// BOUNDARY
func TestReq14_NilSeamYieldsUnevaluableOnlyForAPresentKeyValueAtom(t *testing.T) {
	cases := map[string]struct {
		atoms      []resolve.GuardAtom
		wantRefuse bool
	}{
		"value atom over a PRESENT key refuses": {
			[]resolve.GuardAtom{allAtom("reviews", opEq, "2")}, true,
		},
		"value atom over an ABSENT key refuses (on absence, not on the seam)": {
			[]resolve.GuardAtom{allAtom(absentKey, opEq, "3")}, true,
		},
		"existence atom over a PRESENT key plans": {
			[]resolve.GuardAtom{existsAll("reviews", true)}, false,
		},
		"existence atom over an ABSENT key plans": {
			[]resolve.GuardAtom{existsAll(absentKey, false)}, false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			row := guardedRow("rdr.nilseam.perAtom", "flows/rdr.toml:21", tc.atoms...)
			got := mustResolve(t, guardInput("rev-req14", nil, row))
			if got.Refused() != tc.wantRefuse {
				t.Fatalf("refused=%v (%+v); want refused=%v", got.Refused(), got, tc.wantRefuse)
			}
		})
	}
}

// REQ-15: "Such an atom is reported in the payload like any other
// unevaluable atom, with reason `uncomparable` (its key is present, so
// `absent` would be a lie); the payload never omits it and never names the
// nil seam as a reason."
// DOMAIN EDGE
//
// Testing Strategy row 22: assert the payload ENTRY, not just the kind.
func TestReq15_NilSeamPresentKeyAtomIsPayloadUncomparable(t *testing.T) {
	row := guardedRow("rdr.nilseam.payload", "flows/rdr.toml:22",
		allAtom("reviews", opGte, "3"),
	)
	r := mustRefuse(t, guardInput("rev-req15", nil, row), resolve.KindGuardUnevaluable)

	atoms := undecidedAtomsFor(r, row.RuleID)
	if len(atoms) != 1 {
		t.Fatalf("payload carries %d atoms (%v); want exactly the one present-key "+
			"value atom — the payload never omits it", len(atoms), atoms)
	}
	got := atoms[0]
	want := resolve.UndecidedAtom{
		Key:      "reviews",
		Block:    resolve.BlockAll,
		Operator: opGte,
		Literal:  "3",
		Reason:   resolve.ReasonUncomparable,
	}
	if got != want {
		t.Errorf("payload atom = %+v; want %+v — a nil seam is a wiring fact, never "+
			"a reason of its own", got, want)
	}
}

// REQ-16: "This narrows the shipped whole-guard behavior (`evaluateGuard`
// returns GuardUnevaluable for any non-empty guard when `seam == nil`) to
// a per-atom rule."
// ADVERSARIAL
//
// Testing Strategy row 19 requires BOTH legs under the SAME nil seam: the
// pure-exists row plans and its Writes land, and a value atom over a
// present key still refuses. Without this, a wiring bug that leaves
// Input.Guards unset ships as silent plan-production.
func TestReq16_NilSeamIsPerAtomNotWholeGuard(t *testing.T) {
	t.Run("pure-exists row plans and its writes land", func(t *testing.T) {
		row := guardedRow("rdr.narrow.plans", "flows/rdr.toml:23", existsAll("status", true))
		plan := mustPlan(t, guardInput("rev-req16-plan", nil, row))
		if len(plan.Writes) == 0 {
			t.Fatal("plan carries no writes; the widened behaviour is precisely that " +
				"a pure-exists row under a nil seam produces a plan whose writes land")
		}
		if !reflect.DeepEqual(plan.Writes, row.Writes) {
			t.Errorf("plan writes = %v; want %v", plan.Writes, row.Writes)
		}
	})

	t.Run("value atom over a present key still refuses under the same nil seam", func(t *testing.T) {
		row := guardedRow("rdr.narrow.refuses", "flows/rdr.toml:24", allAtom("status", opEq, "Draft"))
		mustRefuse(t, guardInput("rev-req16-refuse", nil, row), resolve.KindGuardUnevaluable)
	})
}

// REQ-17: "Frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
// keeps its verdict but changes its REASON … To keep testing the nil-seam
// rule the fixture needs a value atom over a PRESENT key (e.g.
// `reviews >= 3`)."
// ADVERSARIAL
//
// The re-decided disposition, asserted here as its own obligation: a
// guarded ESCAPE row whose atom is a value atom over the PRESENT key
// `reviews` must not rescue under a nil seam, and its payload must say
// `uncomparable` — not `absent`, and not want-of-a-seam.
func TestReq17_Fixup1dRedecidedWithAValueAtomOverAPresentKey(t *testing.T) {
	escape := conformingEscapeRow("rdr.escape.nilseam", "flows/rdr.toml:90",
		resolve.KindNoMatch, allAtom("reviews", opGte, "3"))

	in := noMatchInput()
	in.Table.Revision = "rev-req17"
	in.Table.Rows = append(in.Table.Rows, escape)
	in.Guards = nil // no seam at all, not merely a seam that cannot decide

	r := mustRefuse(t, in, resolve.KindGuardUnevaluable)

	reason, ok := reasonFor(undecidedAtomsFor(r, escape.RuleID), "reviews", resolve.BlockAll)
	if !ok || reason != resolve.ReasonUncomparable {
		t.Errorf("escape payload reason = %q (found=%v); want %q — `reviews` is "+
			"present in the view, so the nil seam is what could not compare it",
			reason, ok, resolve.ReasonUncomparable)
	}
}

// REQ-22: "The constants are `OpExists` (the operator token) and
// `LiteralTrue` / `LiteralFalse` (the two boolean literal forms); the
// kernel compares an atom's operator and literal against these values
// verbatim, performing no parsing, case-folding, or coercion of its own."
// ADVERSARIAL
func TestReq22_KernelComparesOperatorAndLiteralVerbatim(t *testing.T) {
	cases := map[string]resolve.GuardAtom{
		"case-folded token":   {Key: "reviews", Operator: "Exists", Literal: resolve.LiteralTrue, Block: resolve.BlockAll},
		"case-folded literal": {Key: "reviews", Operator: resolve.OpExists, Literal: "True", Block: resolve.BlockAll},
		"padded token":        {Key: "reviews", Operator: " exists", Literal: resolve.LiteralTrue, Block: resolve.BlockAll},
		"coerced literal":     {Key: "reviews", Operator: resolve.OpExists, Literal: "1", Block: resolve.BlockAll},
	}
	for name, atom := range cases {
		t.Run(name, func(t *testing.T) {
			// The seam is nil so nothing but the kernel can decide. A kernel
			// that case-folds or coerces would decide the atom TRUE from
			// presence and emit a plan.
			row := guardedRow("rdr.verbatim", "flows/rdr.toml:25", atom)
			mustRefuse(t, guardInput("rev-req22", nil, row), resolve.KindGuardUnevaluable)
		})
	}
}

// REQ-23: "`OpExists = \"exists\"`, `LiteralTrue = \"true\"`,
// `LiteralFalse = \"false\"`, all lower-case."
// BOUNDARY
func TestReq23_ExistenceTokenAndLiteralBytesAreNormative(t *testing.T) {
	for name, got := range map[string]struct{ have, want string }{
		"OpExists":     {resolve.OpExists, "exists"},
		"LiteralTrue":  {resolve.LiteralTrue, "true"},
		"LiteralFalse": {resolve.LiteralFalse, "false"},
	} {
		if got.have != got.want {
			t.Errorf("%s = %q; want %q — the comparison IS byte equality, so the "+
				"bytes are normative", name, got.have, got.want)
		}
		if strings.ToLower(got.have) != got.have {
			t.Errorf("%s = %q; want all lower-case", name, got.have)
		}
	}
}

// REQ-24: "an existence atom carrying a foreign TOKEN is a value atom to
// the kernel (unevaluable on absence, handed to the seam on presence)."
// ADVERSARIAL
//
// Testing Strategy row 11: drift at this boundary fails CLOSED.
func TestReq24_ForeignExistenceTokenIsAValueAtom(t *testing.T) {
	t.Run("absent key ⇒ unevaluable, seam never consulted", func(t *testing.T) {
		seam := newValueSeam()
		atom := resolve.GuardAtom{Key: absentKey, Operator: "EXISTS", Literal: resolve.LiteralFalse, Block: resolve.BlockAll}
		row := guardedRow("rdr.foreign.absent", "flows/rdr.toml:26", atom)

		r := mustRefuse(t, guardInput("rev-req24-absent", seam, row), resolve.KindGuardUnevaluable)

		if calls := seam.seen(); len(calls) != 0 {
			t.Errorf("seam saw %v; a value atom over an absent key never crosses", calls)
		}
		reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), absentKey, resolve.BlockAll)
		if !ok || reason != resolve.ReasonAbsent {
			t.Errorf("reason = %q (found=%v); want %q", reason, ok, resolve.ReasonAbsent)
		}
		// The kernel must NOT have decided it from presence: `exists=false`
		// over an absent key would otherwise plan.
		_ = r
	})

	t.Run("present key ⇒ handed to the seam", func(t *testing.T) {
		atom := resolve.GuardAtom{Key: "reviews", Operator: "EXISTS", Literal: resolve.LiteralTrue, Block: resolve.BlockAll}
		seam := newValueSeam().decide(atom, "2", resolve.GuardTrue)
		row := guardedRow("rdr.foreign.present", "flows/rdr.toml:27", atom)

		mustPlan(t, guardInput("rev-req24-present", seam, row))

		if calls := seam.seen(); len(calls) != 1 || calls[0] != callOf(atom, "2") {
			t.Errorf("seam saw %v; a foreign-token atom over a present key is a "+
				"value atom and must be handed across", calls)
		}
	})
}

// REQ-25: "A foreign LITERAL on an `OpExists` atom — any value that is
// neither `LiteralTrue` nor `LiteralFalse`, the empty literal included —
// is UNEVALUABLE at the kernel, reason `uncomparable`; the kernel MUST NOT
// decide such an atom from presence, and MUST NOT treat a missing literal
// as `LiteralFalse`."
// ADVERSARIAL
func TestReq25_ForeignLiteralOnAnExistsAtomIsKernelUncomparable(t *testing.T) {
	for name, literal := range map[string]string{
		"empty literal":     "",
		"case-folded true":  "True",
		"numeric":           "1",
		"unrelated word":    "yes",
		"whitespace-padded": " true",
	} {
		t.Run(name, func(t *testing.T) {
			for _, key := range []string{"reviews", absentKey} {
				t.Run(key, func(t *testing.T) {
					seam := newValueSeam()
					atom := resolve.GuardAtom{
						Key: key, Operator: resolve.OpExists,
						Literal: literal, Block: resolve.BlockAll,
					}
					row := guardedRow("rdr.foreign.literal", "flows/rdr.toml:28", atom)

					r := mustRefuse(t, guardInput("rev-req25", seam, row), resolve.KindGuardUnevaluable)

					if calls := seam.seen(); len(calls) != 0 {
						t.Errorf("seam saw %v; a foreign literal on an OpExists atom is "+
							"the KERNEL's fail-closed backstop, not a seam question", calls)
					}
					reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), key, resolve.BlockAll)
					if !ok || reason != resolve.ReasonUncomparable {
						t.Errorf("reason = %q (found=%v); want %q — never decided from "+
							"presence, and a missing literal is not LiteralFalse",
							reason, ok, resolve.ReasonUncomparable)
					}
				})
			}
		})
	}
}

// REQ-27: "PRESENCE IS PROVENANCE-BLIND. \"Present in the assembled
// view\" means the key is in the view under ANY provenance — owned,
// observed, or recognized. An atom MUST NOT be decided differently
// according to how a tag reached the view."
// DOMAIN EDGE
//
// Testing Strategy row 18: TagSet.Lookup `ok`, not TagSet.has.
func TestReq27_PresenceIsProvenanceBlind(t *testing.T) {
	const key = "lane"
	const value = "fast"

	place := map[string]func(*resolve.Input){
		"owned":    func(in *resolve.Input) { in.Owned = append(in.Owned, resolve.Tag{Key: key, Value: value}) },
		"observed": func(in *resolve.Input) { in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: value}) },
	}

	for name, put := range place {
		t.Run("value atom decides identically when the key arrives "+name, func(t *testing.T) {
			atom := allAtom(key, opEq, value)
			seam := newValueSeam().decide(atom, value, resolve.GuardTrue)
			row := guardedRow("rdr.provenance", "flows/rdr.toml:29", atom)
			row.RequiresOwned = nil // do not entangle the owned sweep

			in := guardInput("rev-req27", seam, row)
			put(&in)

			mustPlan(t, in)
			if calls := seam.seen(); len(calls) != 1 || calls[0].Value != value {
				t.Errorf("seam saw %v; presence must be decided the same way under "+
					"every provenance", calls)
			}
		})

		t.Run("existence atom decides identically when the key arrives "+name, func(t *testing.T) {
			row := guardedRow("rdr.provenance.exists", "flows/rdr.toml:30", existsAll(key, true))
			row.RequiresOwned = nil

			in := guardInput("rev-req27-exists", nil, row)
			put(&in)

			mustPlan(t, in)
		})
	}

	t.Run("existence atom decides on the recognized tag", func(t *testing.T) {
		row := guardedRow("rdr.provenance.recognized", "flows/rdr.toml:31",
			existsAll("recognized", true))
		row.RequiresOwned = nil
		mustPlan(t, guardInput("rev-req27-recognized", nil, row))
	})
}

// REQ-28: "Key identity is exact string equality on the key as assembled;
// canonicalization of authored key spellings is the normalizer's (RDR
// 0002), upstream of the kernel, and the kernel performs none (A22)."
// ADVERSARIAL
func TestReq28_KeyIdentityIsExactStringEqualityWithNoCanonicalization(t *testing.T) {
	for _, spelling := range []string{"Reviews", "REVIEWS", "reviews ", " reviews", "re-views", "re_views"} {
		t.Run(spelling, func(t *testing.T) {
			seam := newValueSeam()
			row := guardedRow("rdr.keyid", "flows/rdr.toml:32", allAtom(spelling, opEq, "2"))

			r := mustRefuse(t, guardInput("rev-req28", seam, row), resolve.KindGuardUnevaluable)

			if calls := seam.seen(); len(calls) != 0 {
				t.Errorf("seam saw %v; %q is not the assembled key %q and the kernel "+
					"canonicalizes nothing", calls, spelling, "reviews")
			}
			reason, ok := reasonFor(undecidedAtomsFor(r, row.RuleID), spelling, resolve.BlockAll)
			if !ok || reason != resolve.ReasonAbsent {
				t.Errorf("reason = %q (found=%v); want %q", reason, ok, resolve.ReasonAbsent)
			}
		})
	}
}

// REQ-29: "This is deliberately NOT the kernel's owned-state predicate:
// `missingOwned` tests owned provenance because `owned_state_unavailable`
// is about the owned snapshot; guard decidability takes the broader test.
// The two MUST NOT be conflated."
// DOMAIN EDGE
//
// The discriminating input: the guard's key arrives OBSERVED. Guard
// decidability must accept it (presence is provenance-blind), while the
// owned sweep over the same key must NOT — the two predicates disagree on
// exactly this input, which is what "MUST NOT be conflated" means.
func TestReq29_GuardDecidabilityAndOwnedStateUseDifferentPredicates(t *testing.T) {
	const key = "quorum"
	atom := allAtom(key, opEq, "3")
	seam := newValueSeam().decide(atom, "3", resolve.GuardTrue)

	t.Run("guard decidability accepts an observed key", func(t *testing.T) {
		row := guardedRow("rdr.split.guard", "flows/rdr.toml:33", atom)
		row.RequiresOwned = []string{"status"} // satisfied

		in := guardInput("rev-req29-guard", seam, row)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "3"})

		mustPlan(t, in)
	})

	t.Run("the owned sweep rejects the same observed key", func(t *testing.T) {
		row := guardedRow("rdr.split.owned", "flows/rdr.toml:34", atom)
		row.RequiresOwned = []string{"status", key}

		in := guardInput("rev-req29-owned", seam, row)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "3"})

		r := mustRefuse(t, in, resolve.KindOwnedStateUnavailable)
		if !hasItem(r.MissingOwned, key) {
			t.Errorf("MissingOwned = %v; want %q — owned_state_unavailable is about "+
				"the owned snapshot, and an observed tag does not satisfy it",
				r.MissingOwned, key)
		}
	})
}

// REQ-49: "`Reason` is a kernel-owned closed constant set — `ReasonAbsent`,
// `ReasonUncomparable` — spelled and enumerated the way
// `resolve.go::RefusalKind` / `RefusalKinds()` already spell the refusal
// taxonomy: a named STRING type with exported constants plus an exported
// enumerator `Reasons() []Reason` returning the set in declaration order,
// mirroring `RefusalKinds()`. The enumerator is required surface, not an
// implementation choice"
// BOUNDARY
func TestReq49_ReasonIsAClosedNamedStringSetWithAnEnumerator(t *testing.T) {
	rt := reflect.TypeOf(resolve.ReasonAbsent)
	if rt.Kind() != reflect.String || rt.Name() != "Reason" || rt.PkgPath() == "" {
		t.Fatalf("Reason constants have type %q (kind %s, pkg %q); want the exported "+
			"named string type resolve.Reason", rt.Name(), rt.Kind(), rt.PkgPath())
	}
	if resolve.ReasonAbsent != "absent" {
		t.Errorf("ReasonAbsent = %q; want %q", resolve.ReasonAbsent, "absent")
	}
	if resolve.ReasonUncomparable != "uncomparable" {
		t.Errorf("ReasonUncomparable = %q; want %q", resolve.ReasonUncomparable, "uncomparable")
	}

	want := []resolve.Reason{resolve.ReasonAbsent, resolve.ReasonUncomparable}
	if got := resolve.Reasons(); !reflect.DeepEqual(got, want) {
		t.Errorf("Reasons() = %v; want %v in declaration order", got, want)
	}

	// Mirroring RefusalKinds(): a fresh slice per call, so a caller cannot
	// mutate the kernel's set.
	first := resolve.Reasons()
	if len(first) == 0 {
		t.Fatal("Reasons() returned an empty set")
	}
	first[0] = "mutated"
	if resolve.Reasons()[0] != resolve.ReasonAbsent {
		t.Error("Reasons() aliases kernel state; RefusalKinds() returns a fresh slice per call")
	}
}

// REQ-47: "The closed set stays at two members and every unevaluable atom
// of an undecidable row therefore carries a reason — the payload's
// per-atom completeness obligation above admits no third state and no
// omitted entry. A nil seam is not itself a reason"
// BOUNDARY
func TestReq47_ReasonSetStaysAtExactlyTwoMembers(t *testing.T) {
	if got := len(resolve.Reasons()); got != 2 {
		t.Errorf("Reasons() has %d members (%v); the closed set stays at two — "+
			"a nil seam is a wiring fact, not a third reason",
			got, resolve.Reasons())
	}
	for _, r := range resolve.Reasons() {
		if r != resolve.ReasonAbsent && r != resolve.ReasonUncomparable {
			t.Errorf("Reasons() carries %q; want only %q and %q",
				r, resolve.ReasonAbsent, resolve.ReasonUncomparable)
		}
	}
}

// REQ-48: "`Refusal` carries the field `Undecided []UndecidedRow`, where an
// `UndecidedRow` names its `RuleID`, `SourceLocator`, and `Atoms
// []UndecidedAtom`, and an `UndecidedAtom` names its `Key`, `Block`,
// `Operator`, `Literal`, and `Reason`."
// BOUNDARY
//
// "The payload is named surface, not shape-by-description."
func TestReq48_PayloadIsNamedSurface(t *testing.T) {
	undecided, ok := reflect.TypeOf(resolve.Refusal{}).FieldByName("Undecided")
	if !ok {
		t.Fatal("resolve.Refusal has no Undecided field")
	}
	if undecided.Type != reflect.TypeOf([]resolve.UndecidedRow{}) {
		t.Fatalf("Refusal.Undecided is %s; want []resolve.UndecidedRow", undecided.Type)
	}

	row := reflect.TypeOf(resolve.UndecidedRow{})
	for name, want := range map[string]reflect.Type{
		"RuleID":        reflect.TypeOf(""),
		"SourceLocator": reflect.TypeOf(""),
		"Atoms":         reflect.TypeOf([]resolve.UndecidedAtom{}),
	} {
		f, ok := row.FieldByName(name)
		if !ok {
			t.Errorf("UndecidedRow has no field %q", name)
			continue
		}
		if f.Type != want {
			t.Errorf("UndecidedRow.%s is %s; want %s", name, f.Type, want)
		}
	}

	atom := reflect.TypeOf(resolve.UndecidedAtom{})
	for name, want := range map[string]reflect.Type{
		"Key":      reflect.TypeOf(""),
		"Block":    reflect.TypeOf(resolve.BlockAll),
		"Operator": reflect.TypeOf(""),
		"Literal":  reflect.TypeOf(""),
		"Reason":   reflect.TypeOf(resolve.ReasonAbsent),
	} {
		f, ok := atom.FieldByName(name)
		if !ok {
			t.Errorf("UndecidedAtom has no field %q", name)
			continue
		}
		if f.Type != want {
			t.Errorf("UndecidedAtom.%s is %s; want %s", name, f.Type, want)
		}
	}
}

// REQ-50: "`Block` is the same exported constant type the atom carries,
// not a separately-spelled payload value."
// BOUNDARY
func TestReq50_PayloadBlockIsTheAtomsBlockType(t *testing.T) {
	atomBlock, ok := reflect.TypeOf(resolve.GuardAtom{}).FieldByName("Block")
	if !ok {
		t.Fatal("resolve.GuardAtom has no Block field")
	}
	payloadBlock, ok := reflect.TypeOf(resolve.UndecidedAtom{}).FieldByName("Block")
	if !ok {
		t.Fatal("resolve.UndecidedAtom has no Block field")
	}
	if atomBlock.Type != payloadBlock.Type {
		t.Errorf("atom Block is %s, payload Block is %s; the payload must carry the "+
			"same exported constant type, not a separately-spelled value",
			atomBlock.Type, payloadBlock.Type)
	}
	if payloadBlock.Type == reflect.TypeOf("") {
		t.Error("payload Block is a bare string; it must be the named Block type")
	}
}

// REQ-79: "verdicts keep `GuardUnevaluable` / `guard_unevaluable`;
// rejected: a sixth refusal kind (reopens RDR 0001's closed taxonomy). The
// absent-key payload is a field on `Refusal`, not a new kind."
// BOUNDARY
func TestReq79_NoSixthRefusalKindIsMinted(t *testing.T) {
	kinds := resolve.RefusalKinds()
	if len(kinds) != 5 {
		t.Errorf("RefusalKinds() has %d members (%v); RDR 0001 REQ-7 pins the set at "+
			"exactly five and this RDR mints no sixth", len(kinds), kinds)
	}
	if !hasItem(kinds, resolve.KindGuardUnevaluable) {
		t.Errorf("RefusalKinds() = %v; guard_unevaluable is the verdict this RDR keeps", kinds)
	}
	if resolve.KindGuardUnevaluable != "guard_unevaluable" {
		t.Errorf("KindGuardUnevaluable = %q; want %q", resolve.KindGuardUnevaluable, "guard_unevaluable")
	}
}

// REQ-77: "`assemble`, `gate`, `missingOwned`, `escapeOrRefuse`, and
// `Resolve` are untouched in control flow (the payload is a data change
// `gate` populates, not a reordering)." The change "is confined to
// `internal/resolve/resolve.go`" plus test fixtures.
// BOUNDARY
func TestReq77_ChangeIsConfinedToTheKernelPackage(t *testing.T) {
	for _, name := range []string{"assemble", "gate", "missingOwned", "escapeOrRefuse", "Resolve"} {
		if !kernelDeclaresFunc(t, name) {
			t.Errorf("kernel no longer declares %q; this RDR leaves the control-flow "+
				"functions in place and changes data, not ordering", name)
		}
	}
}

// REQ-78: "The domain rule is scoped to guards, so the match pattern still
// folds absence into non-match. `resolve.go::TagSet.matches` returns false
// on `!ok` … This is deliberate … and NOT closed by this RDR".
// DOMAIN EDGE
//
// The negative contract: a row whose MATCH pattern names an absent key
// must still drop out of candidacy and reach the ESCAPABLE no_match — the
// opposite of what the guard path now does with the same absent key.
func TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch(t *testing.T) {
	row := guardedRow("rdr.match.absent", "flows/rdr.toml:35")
	row.Match = append(row.Match, resolve.Tag{Key: absentKey, Value: "3"})
	escape := conformingEscapeRow("rdr.escape.match", "flows/rdr.toml:36", resolve.KindNoMatch)

	plan := mustPlan(t, guardInput("rev-req78", newValueSeam(), row, escape))
	if !plan.Escaped {
		t.Errorf("plan came from rule %q unescaped; a match pattern over an absent "+
			"key drops the row out of candidacy and yields the ESCAPABLE no_match",
			plan.RuleID)
	}
	if plan.RuleID != escape.RuleID {
		t.Errorf("plan names %q; want the escape row %q", plan.RuleID, escape.RuleID)
	}
}

// REQ-80: "Evaluation stays linear in the atom count per candidate row …
// and the payload sort is over undecidable rows only, on the refusal path.
// `Resolve` is pure and in-process; no new I/O."
// BOUNDARY
//
// A shape constraint, not a benchmark: linearity is checkable as "one seam
// call per present-key value atom, and no more".
func TestReq80_OneSeamCallPerPresentKeyValueAtom(t *testing.T) {
	a := allAtom("status", opEq, "Draft")
	b := allAtom("reviews", opEq, "2")
	c := allAtom(absentKey, opEq, "3") // absent: never crosses
	d := existsAll("status", true)     // existence: never crosses
	seam := newValueSeam().
		decide(a, "Draft", resolve.GuardTrue).
		decide(b, "2", resolve.GuardTrue)

	row := guardedRow("rdr.linear", "flows/rdr.toml:37", a, b, c, d)
	mustRefuse(t, guardInput("rev-req80", seam, row), resolve.KindGuardUnevaluable)

	if got := len(seam.seen()); got != 2 {
		t.Errorf("seam was consulted %d times (%v); want exactly one call per "+
			"present-key value atom", got, seam.seen())
	}
}

// --- helpers -------------------------------------------------------------

func containsCall(calls []seamCall, want seamCall) bool {
	for _, c := range calls {
		if c == want {
			return true
		}
	}
	return false
}

func hasItem[T comparable](s []T, want T) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

// kernelDeclaresFunc reports whether the kernel package declares a
// function or method with the given name.
func kernelDeclaresFunc(t *testing.T, name string) bool {
	t.Helper()
	for _, f := range parseKernelPackage(t) {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Name.Name == name {
				return true
			}
		}
	}
	return false
}
