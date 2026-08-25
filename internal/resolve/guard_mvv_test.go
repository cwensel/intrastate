package resolve_test

// RDR 0007 — Row.RequiresOwned's narrowed contract, the cross-RDR value-
// seam contract-test export, the remaining phase obligations, and the
// runnable Minimum Viable Validation.

import (
	"go/ast"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-59: "Row.RequiresOwned names the owned tag keys the row's post-guard
// transition depends on — the keys its `Writes` require, which per RDR
// 0002 includes an authored clear (normalization renders it as a
// `<clear>` write)."
// BOUNDARY
//
// The doc contract is what Phase 1 narrows, so the doc comment is the
// surface: it must state the post-guard write dependency and must not
// state the superseded "evaluation needs" reading the shipped comment
// carries.
func TestReq59_RequiresOwnedDocContractIsNarrowedToPostGuardWrites(t *testing.T) {
	doc := fieldDoc(t, "Row", "RequiresOwned")
	if doc == "" {
		t.Fatal("Row.RequiresOwned carries no doc comment; Phase 1 narrows its doc contract")
	}
	if !hasText(doc, "Writes") {
		t.Errorf("Row.RequiresOwned's doc does not name Writes:\n%s\n"+
			"the field names the owned keys its Writes require", doc)
	}
	if hasText(doc, "row's evaluation needs") {
		t.Errorf("Row.RequiresOwned's doc still states the superseded "+
			"evaluation-needs reading:\n%s", doc)
	}
}

// REQ-60: "Guard decidability is not RequiresOwned's job: guard-input
// coverage is enforced by the domain rule above."
// DOMAIN EDGE
//
// The discriminating input: a guard reads a key that is NOT in
// RequiresOwned and is absent from the view. The domain rule must refuse
// anyway — coverage is not RequiresOwned's to enforce, so an empty
// RequiresOwned does not make the guard decidable.
func TestReq60_GuardInputCoverageIsNotRequiresOwnedsJob(t *testing.T) {
	row := guardedRow("rdr.coverage", "flows/rdr.toml:120", allAtom(absentKey, opEq, "3"))
	row.RequiresOwned = nil // names no guard key at all

	r := mustRefuse(t, guardInput("rev-req60", newValueSeam(), row),
		resolve.KindGuardUnevaluable)
	if !hasItem(undecidedRuleIDs(r), row.RuleID) {
		t.Errorf("payload names %v; want %q — the domain rule enforces guard-input "+
			"coverage regardless of RequiresOwned", undecidedRuleIDs(r), row.RuleID)
	}
}

// REQ-61: "Listing a guard-read owned key in RequiresOwned remains legal
// and yields the more precise owned_state_unavailable diagnosis among
// survivors — and is the ONLY way a guard over an owned tag is protected
// from a caller-supplied observed tag satisfying presence (A13): the
// guard's key set is NOT required to be a subset of RequiresOwned, because
// guards legitimately read observed and recognized tags."
// DOMAIN EDGE
//
// Testing Strategy row 16's protection half.
func TestReq61_ListingAGuardKeyInRequiresOwnedProtectsAgainstObservedSubstitution(t *testing.T) {
	const key = "gate"
	atom := allAtom(key, opEq, "open")
	seam := newValueSeam().decide(atom, "open", resolve.GuardTrue)

	t.Run("unprotected: an observed tag satisfies presence and the row plans", func(t *testing.T) {
		row := guardedRow("rdr.a13.unprotected", "flows/rdr.toml:121", atom)
		row.RequiresOwned = []string{"status"} // guard key NOT listed

		in := guardInput("rev-req61-open", seam, row)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "open"})

		mustPlan(t, in)
	})

	t.Run("protected: listing the guard key yields owned_state_unavailable", func(t *testing.T) {
		row := guardedRow("rdr.a13.protected", "flows/rdr.toml:122", atom)
		row.RequiresOwned = []string{"status", key}

		in := guardInput("rev-req61-protected", seam, row)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "open"})

		r := mustRefuse(t, in, resolve.KindOwnedStateUnavailable)
		if !hasItem(r.MissingOwned, key) {
			t.Errorf("MissingOwned = %v; want %q — this is the ONLY way a guard over "+
				"an owned tag is protected from a caller-supplied observed tag",
				r.MissingOwned, key)
		}
	})

	t.Run("the guard key set is not required to be a subset of RequiresOwned", func(t *testing.T) {
		// A guard over the OBSERVED key `reviews` and the RECOGNIZED key,
		// with RequiresOwned naming neither, resolves normally.
		row := guardedRow("rdr.a13.subset", "flows/rdr.toml:123",
			existsAll("reviews", true),
			existsAll("recognized", true),
		)
		row.RequiresOwned = []string{"status"}
		mustPlan(t, guardInput("rev-req61-subset", nil, row))
	})
}

// REQ-62: "On an escape row, whose `Writes` MUST be empty (RDR 0009),
// RequiresOwned is therefore empty (A21)."
// BOUNDARY
//
// Testing Strategy row 15: the test MUST build the escape row explicitly
// with empty Writes and empty RequiresOwned and assert the kernel
// consequence. The shipped fixtures_test.go::escapeRow is non-conforming
// on BOTH fields and is deliberately not reused.
func TestReq62_ConformingEscapeRowRaisesNoOwnedStateOfItsOwn(t *testing.T) {
	escape := conformingEscapeRow("rdr.escape.a21", "flows/rdr.toml:124",
		resolve.KindNoMatch, existsAll("status", true))

	if len(escape.Writes) != 0 {
		t.Fatalf("fixture escape row carries Writes %v; RDR 0009 requires them empty",
			escape.Writes)
	}
	if len(escape.RequiresOwned) != 0 {
		t.Fatalf("fixture escape row carries RequiresOwned %v; A21's composition "+
			"makes it empty", escape.RequiresOwned)
	}

	// The kernel consequence: the escape rescues without raising an owned
	// obligation of its own, even though the owned snapshot carries only
	// `status`.
	in := noMatchInput()
	in.Table.Revision = "rev-req62"
	in.Table.Rows = append(in.Table.Rows, escape)
	in.Guards = newValueSeam()

	plan := mustPlan(t, in)
	if plan.RuleID != escape.RuleID || !plan.Escaped {
		t.Errorf("plan names %q (escaped=%v); want the escape %q",
			plan.RuleID, plan.Escaped, escape.RuleID)
	}
	if len(plan.Writes) != 0 {
		t.Errorf("escape plan carries writes %v; RDR 0009 requires them empty", plan.Writes)
	}
}

// REQ-68: "migrate every fixture so the frozen suite drives verdicts
// THROUGH atoms and a value stub — never by injecting row verdicts, which
// would leave the kernel combinator off the tested path (premortem P-20)"
// BOUNDARY
//
// The checkable form of "never by injecting row verdicts": no test seam
// implements GuardEvaluator by keying on anything but the atom it is
// given. Structurally, the narrowed interface makes injection
// unrepresentable — a seam cannot be handed a whole row — so the assertion
// is on the interface itself plus the absence of any row-verdict field on
// Row.
func TestReq68_VerdictsAreDrivenThroughAtomsNotInjectedAtRowLevel(t *testing.T) {
	iface := reflect.TypeOf((*resolve.GuardEvaluator)(nil)).Elem()
	sig := iface.Method(0).Type
	for i := range sig.NumIn() {
		if sig.In(i) == reflect.TypeOf(resolve.Row{}) {
			t.Error("the value seam is handed a whole Row; verdicts must be driven " +
				"through atoms so the kernel combinator stays on the tested path")
		}
	}
	if _, ok := reflect.TypeOf(resolve.Row{}).FieldByName("Verdict"); ok {
		t.Error("Row carries a Verdict field; a row-level verdict leaves the kernel " +
			"combinator off the tested path")
	}
}

// REQ-69: "Phase 1's exit condition is that count with the field actually
// removed" — the frozen suite is expected to pass re-encoded EXCEPT
// Fixup-1d, which is RE-DECIDED (REQ-17).
// BOUNDARY
//
// The checkable half at kernel scope: the shipped guard-text plumbing is
// gone from the kernel's surface, so nothing can still be keyed on it.
func TestReq69_TheGuardTextFieldIsActuallyRemoved(t *testing.T) {
	for _, tc := range []struct {
		typ   reflect.Type
		field string
	}{
		{reflect.TypeOf(resolve.Refusal{}), "Guard"},
	} {
		if _, ok := tc.typ.FieldByName(tc.field); ok {
			t.Errorf("%s still carries %q; Phase 1's exit condition is that count "+
				"with the field actually removed", tc.typ, tc.field)
		}
	}
	guard, ok := reflect.TypeOf(resolve.Row{}).FieldByName("Guard")
	if !ok {
		t.Fatal("resolve.Row has no Guard field")
	}
	if guard.Type.Kind() == reflect.String {
		t.Error("Row.Guard is still a string; the shipped field is replaced by the " +
			"atom slice, not retained beside it")
	}
}

// REQ-71: "Export a small contract-test function for the value seam
// (present value × literal × operator; unparseable value → unevaluable,
// never false) that RDR 0003's implement stage instantiates against its
// evaluator. It is `resolve.TestGuardEvaluatorContract(t *testing.T, seam
// GuardEvaluator)` — named here because it is a CROSS-RDR surface 0003
// must call by name"
// BOUNDARY
//
// The exact name and signature are normative surface. It must live in
// non-`_test.go` source so RDR 0003 can import it: Go test files are not
// importable. Both halves are asserted — the signature by taking its
// address, and the file location by inspecting the package's non-test
// sources.
func TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface(t *testing.T) {
	var fn func(*testing.T, resolve.GuardEvaluator) = resolve.TestGuardEvaluatorContract
	if fn == nil {
		t.Fatal("resolve.TestGuardEvaluatorContract is nil")
	}

	if !kernelDeclaresFunc(t, "TestGuardEvaluatorContract") {
		t.Error("TestGuardEvaluatorContract is not declared in the kernel's non-test " +
			"sources; RDR 0003 must import it by name, and Go test files are not importable")
	}
}

// REQ-71 (behaviour): the contract test's own obligations — "present
// value × literal × operator; unparseable value → unevaluable, never
// false".
// ADVERSARIAL
//
// Two halves, because either alone is satisfiable by a hollow harness.
// (1) The contract test must actually EXERCISE the named product: it is
// driven here by a recording seam, and what it asked is the assertion.
// (2) A seam implementing RDR 0003's semantics honestly must PASS it, so
// the contract test is not simply impossible to satisfy.
func TestReq71_ContractTestExercisesThePresentValueLiteralOperatorProduct(t *testing.T) {
	// Driven by a CONFORMING seam wrapped in a recorder, so the contract
	// test passes while the questions it asked stay observable. Running it
	// against a non-conforming seam is not an option: Go propagates a
	// subtest's failure, so a deliberately failing run would fail this file.
	//
	// The recorder cannot see the contract test's own `want` column — that
	// lives in its local case table and is consumed by its own t.Errorf. It
	// does not have to: the contract test passed against a CONFORMING seam
	// above, so that seam's answers are the `want` column, and the recorder
	// captures them beside each question. What follows asserts the matrix of
	// verdict CLASSES the contract owes per operator, not the incidental
	// shape of today's case table.
	rec := &recordingSeam{inner: conformingContractSeam{}}
	resolve.TestGuardEvaluatorContract(t, rec)

	if len(rec.calls) == 0 {
		t.Fatal("the contract test asked the seam nothing; it must exercise " +
			"present value × literal × operator")
	}

	// Every call carries a PRESENT value and a literal: the contract is
	// over present values only, so an empty value would be out of contract.
	for _, c := range rec.calls {
		if c.Operator == "" {
			t.Errorf("contract case %+v carries no operator", c)
		}
	}

	// The product's operator axis. RDR 0003's operator/kind matrix is the
	// contract's subject, so each typed operator must be exercised BY NAME:
	// a floor of "at least two operators" is discharged by an eq+in harness
	// that never asks the parsing operators anything.
	asked := map[string]bool{}
	for _, c := range rec.calls {
		asked[c.Operator] = true
	}
	for _, op := range []string{opEq, opGte, opIn, opContains} {
		if !asked[op] {
			t.Errorf("the contract test never exercises operator %q; the contract is "+
				"over the present value × literal × OPERATOR product, and %q is in "+
				"RDR 0003's operator matrix", op, op)
		}
	}

	// The product's value axis, PER OPERATOR: each operator must be driven
	// to both TRUE and FALSE while ONE of its two operands is held fixed.
	// Both halves of that are load-bearing, and each rules out a different
	// constant seam.
	//
	// Verdict classes rather than distinct values, because distinctness is
	// only a proxy for decided-both-ways, and an unsound one: `gte` against
	// literal "3" with values "4" and "3" varies the value while answering
	// TRUE to both, so a table that lost its FALSE leg would still admit a
	// seam constant-TRUE for every parseable `gte`.
	//
	// With an operand pinned rather than merely somewhere-per-operator,
	// because two verdicts that share neither operand are explained by a
	// seam that ignores one of them: `gte` "1" vs "9" both against value
	// "5" decides both ways without ever exercising the value axis.
	//
	// EITHER operand may be the pinned one, because which side carries the
	// variation is the operator's own business. `gte` and `eq` vary the
	// value against a fixed bound; `contains` is the mirror — the value is
	// the container and the literal is the probe, so the shipped table
	// varies the probe against a fixed container. Demanding a fixed literal
	// for every operator would reject that legitimate shape.
	sameLiteral := map[string]map[resolve.GuardResult]bool{}
	sameValue := map[string]map[resolve.GuardResult]bool{}
	for i, c := range rec.calls {
		lk := c.Operator + "\x00" + c.Literal
		vk := c.Operator + "\x00" + c.Value
		if sameLiteral[lk] == nil {
			sameLiteral[lk] = map[resolve.GuardResult]bool{}
		}
		if sameValue[vk] == nil {
			sameValue[vk] = map[resolve.GuardResult]bool{}
		}
		sameLiteral[lk][rec.verdicts[i]] = true
		sameValue[vk][rec.verdicts[i]] = true
	}
	bothWays := map[string]bool{}
	for _, group := range []map[string]map[resolve.GuardResult]bool{sameLiteral, sameValue} {
		for k, verdicts := range group {
			if verdicts[resolve.GuardTrue] && verdicts[resolve.GuardFalse] {
				bothWays[k[:strings.Index(k, "\x00")]] = true
			}
		}
	}
	for _, op := range []string{opEq, opGte, opIn, opContains} {
		if asked[op] && !bothWays[op] {
			t.Errorf("the contract test never drives %q to both GuardTrue and "+
				"GuardFalse with one operand held fixed; an operator decided only "+
				"one way cannot be told from a constant, and two verdicts sharing "+
				"neither operand are explained by a seam ignoring one of them", op)
		}
	}

	// The PRESENT VALUE axis proper, which the pinned-operand rule above
	// does not carry on its own: when the fixed operand is the value, the
	// literal did all the discriminating, so that rule alone would admit an
	// operator whose value never moved. The contract is over the present
	// value × literal × operator product, so every operator must be offered
	// at least two distinct values somewhere.
	//
	// Stated per operator over values rather than as a per-operator
	// orientation table, so it holds for any operator the matrix later
	// grows. The shipped table satisfies it on both shapes: `eq`, `gte` and
	// `in` vary the value against a fixed literal, and `contains` varies it
	// between its set value and its unparseable one.
	valuesSeen := map[string]map[string]bool{}
	for _, c := range rec.calls {
		if valuesSeen[c.Operator] == nil {
			valuesSeen[c.Operator] = map[string]bool{}
		}
		valuesSeen[c.Operator][c.Value] = true
	}
	for _, op := range []string{opEq, opGte, opIn, opContains} {
		if asked[op] && len(valuesSeen[op]) < 2 {
			t.Errorf("the contract test offers %q only the value(s) %v; the contract "+
				"is over the PRESENT VALUE × literal × operator product, and an "+
				"operator asked about a single value never exercises the value axis",
				op, valuesSeen[op])
		}
	}

	// The unparseable-value leg the contract exists to enforce. "Cannot
	// parse" is a property of an operator's OWN parse rule, so it is
	// classified per operator against that rule — `gte` parses an integer,
	// `contains` parses a §D13 set on both sides. A global heuristic over
	// all calls is vacuous: `in`'s shipped non-member value "gamma" is a
	// perfectly VALID `in` value, and would discharge the leg with no
	// genuinely unparseable case present anywhere in the table. `in` is
	// therefore deliberately excluded — a non-member string is a decided
	// FALSE for `in`, not an unevaluable.
	unparseableValue := map[string]bool{}
	parseableValue := map[string]bool{}
	unparseableLiteral := map[string]bool{}
	for i, c := range rec.calls {
		var valueOK, literalOK bool
		switch c.Operator {
		case opGte:
			_, valueOK = parseInt(c.Value)
			_, literalOK = parseInt(c.Literal)
		case opContains:
			_, valueOK = parseD13Set(c.Value)
			_, literalOK = parseD13Set(c.Literal)
		default:
			continue
		}
		switch {
		case !valueOK:
			// The obligation is `unparseable value → unevaluable, never
			// false`, so the leg only counts when the seam actually
			// answered unevaluable to it.
			if rec.verdicts[i] == resolve.GuardUnevaluable {
				unparseableValue[c.Operator] = true
			}
		case literalOK:
			// A case counts as one the operator CAN parse only when both
			// operands parse: `gte`'s unparseable-literal leg carries a
			// parseable value but is not a decided case.
			parseableValue[c.Operator] = true
		}
		if !literalOK {
			unparseableLiteral[c.Operator] = true
		}
	}
	for _, op := range []string{opGte, opContains} {
		if !asked[op] {
			continue
		}
		if !unparseableValue[op] {
			t.Errorf("the contract test never offers %q a value its own parse rule "+
				"rejects; `unparseable value → unevaluable, never false` is the "+
				"obligation the seam exists to carry", op)
		}
		if !parseableValue[op] {
			t.Errorf("the contract test offers %q no value it CAN parse; the "+
				"unparseable leg is only meaningful beside a decided one", op)
		}
	}

	// The mirror obligation on the literal side: `gte` cannot compare
	// against a bound it cannot parse either, and the contract ships that
	// case. Asserted only for `gte`, whose literal is a bare integer; a
	// §D13 set literal has no comparable one-sided form.
	if asked[opGte] && !unparseableLiteral[opGte] {
		t.Errorf("the contract test never offers %q an unparseable LITERAL; a bound "+
			"the operator cannot parse is unevaluable on the same grounds as a "+
			"value it cannot parse", opGte)
	}
}

// REQ-71 (satisfiability): a seam implementing RDR 0003's semantics
// honestly passes the exported contract test.
// HAPPY PATH
func TestReq71_AConformingValueSeamSatisfiesTheContractTest(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, conformingContractSeam{})
}

// recordingSeam delegates to a conforming seam and records what it was
// asked, so the contract test's own coverage is observable while the
// contract test itself still passes.
//
// It records the conforming seam's ANSWER beside each question. The
// verdict class is the property the contract owes, and here it is
// observed rather than re-derived: `seamCall` is a shared fixture key
// type, so the verdicts ride in a parallel slice indexed alike.
type recordingSeam struct {
	inner    resolve.GuardEvaluator
	calls    []seamCall
	verdicts []resolve.GuardResult
}

func (s *recordingSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	got := s.inner.Evaluate(atom, value)
	s.calls = append(s.calls, callOf(atom, value))
	s.verdicts = append(s.verdicts, got)
	return got
}

// conformingContractSeam implements RDR 0003's typed semantics narrowly
// enough to satisfy the contract: it compares parseable values and answers
// unevaluable — never false — for a value it cannot parse.
type conformingContractSeam struct{}

func (conformingContractSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	switch atom.Operator {
	case opEq:
		if value == atom.Literal {
			return resolve.GuardTrue
		}
		return resolve.GuardFalse
	case opGte, "gt", "lt", "lte":
		a, aok := parseInt(value)
		b, bok := parseInt(atom.Literal)
		if !aok || !bok {
			return resolve.GuardUnevaluable
		}
		switch atom.Operator {
		case opGte:
			return boolVerdict(a >= b)
		case "gt":
			return boolVerdict(a > b)
		case "lte":
			return boolVerdict(a <= b)
		default:
			return boolVerdict(a < b)
		}
	case opIn:
		members, ok := parseD13Set(atom.Literal)
		if !ok {
			return resolve.GuardUnevaluable
		}
		return boolVerdict(hasItem(members, value))
	case opContains:
		want, wok := parseD13Set(atom.Literal)
		have, hok := parseD13Set(value)
		if !wok || !hok {
			return resolve.GuardUnevaluable
		}
		for _, w := range want {
			if !hasItem(have, w) {
				return resolve.GuardFalse
			}
		}
		return resolve.GuardTrue
	default:
		return resolve.GuardUnevaluable
	}
}

// REQ-73: "Every scenario is a kernel test in `internal/resolve` (package
// `resolve_test`)"
// BOUNDARY
func TestReq73_ScenariosAreKernelTestsInPackageResolveTest(t *testing.T) {
	// This file's own package declaration is the assertion's subject; a
	// scenario written in package `resolve` would reach unexported state
	// and stop testing the kernel's exported contract.
	if got := reflect.TypeOf(resolve.Row{}).PkgPath(); got != "github.com/newcoinc/intrastate/internal/resolve" {
		t.Errorf("kernel package path = %q; the scenarios test internal/resolve", got)
	}
	// And the kernel is reachable only through its exported surface from
	// here: the assembled view has no exported field a scenario could use
	// to inject a state the input tuple did not produce.
	view := reflect.TypeOf(resolve.TagSet{})
	for i := range view.NumField() {
		if view.Field(i).PkgPath == "" {
			t.Errorf("TagSet exports field %q; scenarios must drive the kernel "+
				"through its exported inputs, not by injecting a view",
				view.Field(i).Name)
		}
	}
}

// REQ-74: "Any future edit to the combination tables MUST be
// re-mutation-tested against" scenarios 3–5, which kill the
// FALSE-vs-UNEVALUABLE dominance mutant that "survives all 154 frozen
// tests".
// ADVERSARIAL
//
// The mutant this test exists to kill, stated directly: swapping
// FALSE-dominance for UNEVALUABLE-dominance. It is observable only across
// TWO atoms, and only where the two verdicts lead to different
// dispositions — which is why the escape row is part of the fixture.
func TestReq74_FalseVsUnevaluableDominanceMutantIsKilled(t *testing.T) {
	seam := newValueSeam()
	escape := conformingEscapeRow("rdr.escape.mutant", "flows/rdr.toml:125", resolve.KindNoMatch)

	t.Run("F ∧ U = F ⇒ pruned (UNEVALUABLE-dominance would refuse)", func(t *testing.T) {
		f := atomWithVerdict(seam, "f", resolve.BlockAll, resolve.GuardFalse)
		u := atomWithVerdict(seam, "u", resolve.BlockAll, resolve.GuardUnevaluable)
		row := guardedRow("rdr.mutant.fu", "flows/rdr.toml:126", f, u)

		plan := mustPlan(t, guardInput("rev-req74-fu", seam, row, escape))
		if plan.RuleID != escape.RuleID {
			t.Errorf("plan names %q; want the escape %q — F dominates U", plan.RuleID, escape.RuleID)
		}
	})

	t.Run("T ∧ U = U ⇒ refuses (FALSE-dominance would still refuse; raw min would plan)", func(t *testing.T) {
		tr := atomWithVerdict(seam, "t", resolve.BlockAll, resolve.GuardTrue)
		u := atomWithVerdict(seam, "u2", resolve.BlockAll, resolve.GuardUnevaluable)
		row := guardedRow("rdr.mutant.tu", "flows/rdr.toml:127", tr, u)

		mustRefuse(t, guardInput("rev-req74-tu", seam, row, escape), resolve.KindGuardUnevaluable)
	})
}

// REQ-75: "Guard-path observed-tag substitution: a caller-supplied
// observed tag turns `guard_unevaluable` into a SPECIFIC verdict — assert
// the exact disposition (`--tag X=v` against `X eq v` ⇒ row selected and a
// plan; against `X eq w` ⇒ row pruned, no plan), never merely \"not
// `guard_unevaluable`\", which the masking path would also satisfy"
// DOMAIN EDGE
func TestReq75_ObservedTagSubstitutionYieldsASpecificVerdict(t *testing.T) {
	const key = "lane"
	atom := allAtom(key, opEq, "fast")
	row := guardedRow("rdr.substitution", "flows/rdr.toml:128", atom)
	escape := conformingEscapeRow("rdr.escape.substitution", "flows/rdr.toml:129", resolve.KindNoMatch)

	t.Run("no observed tag ⇒ guard_unevaluable", func(t *testing.T) {
		mustRefuse(t, guardInput("rev-req75-none", eqSeam{}, row, escape),
			resolve.KindGuardUnevaluable)
	})

	t.Run("X=fast against `X eq fast` ⇒ the row is SELECTED and a plan is produced", func(t *testing.T) {
		in := guardInput("rev-req75-match", eqSeam{}, row, escape)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "fast"})

		plan := mustPlan(t, in)
		if plan.RuleID != row.RuleID || plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the guarded row %q selected",
				plan.RuleID, plan.Escaped, row.RuleID)
		}
	})

	t.Run("X=slow against `X eq fast` ⇒ the row is PRUNED and no plan comes from it", func(t *testing.T) {
		in := guardInput("rev-req75-mismatch", eqSeam{}, row, escape)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "slow"})

		plan := mustPlan(t, in)
		if plan.RuleID == row.RuleID {
			t.Errorf("plan names the guarded row %q; the value is decided FALSE, so "+
				"the row is pruned", plan.RuleID)
		}
		if plan.RuleID != escape.RuleID || !plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the escape %q — asserting the "+
				"EXACT disposition, not merely `not guard_unevaluable`",
				plan.RuleID, plan.Escaped, escape.RuleID)
		}
	})
}

// REQ-76: "Two-row absence pattern (`X exists = false` row ‖ `X exists =
// true` + `X eq v` row), plus a bare-value-row negative control that must
// refuse"
// DOMAIN EDGE
//
// Testing Strategy row 9. The pattern is disjoint by construction: exactly
// one of the two rows survives under every state of X, so the table
// resolves without refusing whether X is present or absent.
func TestReq76_TwoRowAbsencePatternIsTotalAndTheBareValueRowIsNot(t *testing.T) {
	const key = "quorum"

	absentRow := guardedRow("rdr.pattern.absent", "flows/rdr.toml:130",
		existsAll(key, false))
	presentRow := guardedRow("rdr.pattern.present", "flows/rdr.toml:131",
		existsAll(key, true),
		allAtom(key, opEq, "3"),
	)
	bareRow := guardedRow("rdr.pattern.bare", "flows/rdr.toml:132",
		allAtom(key, opEq, "3"))

	t.Run("pattern resolves when X is absent", func(t *testing.T) {
		plan := mustPlan(t, guardInput("rev-req76-absent", eqSeam{}, absentRow, presentRow))
		if plan.RuleID != absentRow.RuleID {
			t.Errorf("plan names %q; want the exists=false row %q",
				plan.RuleID, absentRow.RuleID)
		}
	})

	t.Run("pattern resolves when X is present and matches", func(t *testing.T) {
		in := guardInput("rev-req76-present", eqSeam{}, absentRow, presentRow)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "3"})

		plan := mustPlan(t, in)
		if plan.RuleID != presentRow.RuleID {
			t.Errorf("plan names %q; want the conjoined row %q", plan.RuleID, presentRow.RuleID)
		}
	})

	t.Run("negative control: the bare value row refuses when X is absent", func(t *testing.T) {
		mustRefuse(t, guardInput("rev-req76-bare", eqSeam{}, bareRow),
			resolve.KindGuardUnevaluable)
	})
}

// REQ-MVV: "Against the real kernel with a real atom — no stub — the
// masking probe inverts: one candidate row whose guard is a value atom
// over an absent key, whose `RequiresOwned` is SATISFIED, and a modeled
// `no_match` escape row yield `Refusal.Kind == guard_unevaluable`, a nil
// `Plan`, the absent key in the payload, and the row in `Refusal.Rows`;
// the same table with the key present and the value decided FALSE prunes
// and escapes (D8 preserved). … Third scenario:
// *unevaluable-blocks-true-sibling* — row A unevaluable beside row B
// decided TRUE → refusal naming A, no plan. Fourth: an `exists = false`
// atom over the absent key decides TRUE and the row is selected — the one
// sanctioned route from absence to a verdict."
// HAPPY PATH (end-to-end)
//
// The four scenarios run against the real resolve.Resolve. The only
// collaborator is the RDR 0003 value seam, which this RDR narrows to
// "compare a present value against a literal" — scenario 1's row never
// reaches it at all, which is the point.
func TestMVV_GuardPredicateTotality(t *testing.T) {
	const key = "quorum"
	valueAtom := allAtom(key, opEq, "3")

	t.Run("1_masking_probe_inverted", func(t *testing.T) {
		row := guardedRow("rdr.mvv.candidate", "flows/rdr.toml:200", valueAtom)
		row.RequiresOwned = []string{"status"} // SATISFIED — load-bearing
		escape := conformingEscapeRow("rdr.mvv.escape", "flows/rdr.toml:201", resolve.KindNoMatch)

		seam := newValueSeam()
		got := mustResolve(t, guardInput("rev-mvv-1", seam, row, escape))

		if got.Plan != nil {
			t.Fatalf("Plan = %+v; want nil — an absent guard key must never reach a "+
				"plan, escaped or otherwise", got.Plan)
		}
		if got.Refusal.Kind != resolve.KindGuardUnevaluable {
			t.Fatalf("refusal kind = %q; want %q — the escapable no_match is what "+
				"the masking probe used to reach", got.Refusal.Kind, resolve.KindGuardUnevaluable)
		}
		if calls := seam.seen(); len(calls) != 0 {
			t.Errorf("seam saw %v; the evaluator is never asked about an absent key", calls)
		}

		// The absent key in the payload, named per row and per atom.
		wantAtom := resolve.UndecidedAtom{
			Key:      key,
			Block:    resolve.BlockAll,
			Operator: opEq,
			Literal:  "3",
			Reason:   resolve.ReasonAbsent,
		}
		wantRow := resolve.UndecidedRow{
			RuleID:        row.RuleID,
			SourceLocator: row.SourceLocator,
			Atoms:         []resolve.UndecidedAtom{wantAtom},
		}
		if !reflect.DeepEqual(got.Refusal.Undecided, []resolve.UndecidedRow{wantRow}) {
			t.Errorf("payload = %+v; want %+v", got.Refusal.Undecided, []resolve.UndecidedRow{wantRow})
		}

		// And the row in Refusal.Rows.
		wantRefs := []resolve.RowRef{{RuleID: row.RuleID, SourceLocator: row.SourceLocator}}
		if !reflect.DeepEqual(got.Refusal.Rows, wantRefs) {
			t.Errorf("Refusal.Rows = %+v; want %+v", got.Refusal.Rows, wantRefs)
		}
	})

	t.Run("2_key_present_value_false_prunes_and_escapes", func(t *testing.T) {
		row := guardedRow("rdr.mvv.candidate", "flows/rdr.toml:200", valueAtom)
		row.RequiresOwned = []string{"status"}
		escape := conformingEscapeRow("rdr.mvv.escape", "flows/rdr.toml:201", resolve.KindNoMatch)

		in := guardInput("rev-mvv-2", eqSeam{}, row, escape)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "2"}) // decides FALSE

		plan := mustPlan(t, in)
		if plan.RuleID != escape.RuleID || !plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the escape %q — decided "+
				"falsity prunes and D8 is preserved", plan.RuleID, plan.Escaped, escape.RuleID)
		}
	})

	t.Run("3_unevaluable_blocks_true_sibling", func(t *testing.T) {
		rowA := guardedRow("rdr.mvv.a", "flows/rdr.toml:202", valueAtom)
		rowB := guardedRow("rdr.mvv.b", "flows/rdr.toml:203", existsAll("status", true))

		got := mustResolve(t, guardInput("rev-mvv-3", newValueSeam(), rowA, rowB))

		if got.Plan != nil {
			t.Fatalf("Plan = %+v; want nil — a decided-TRUE sibling must not be "+
				"selected while an unevaluable candidate exists", got.Plan)
		}
		if got.Refusal.Kind != resolve.KindGuardUnevaluable {
			t.Fatalf("refusal kind = %q; want %q", got.Refusal.Kind, resolve.KindGuardUnevaluable)
		}
		if ids := undecidedRuleIDs(got.Refusal); !reflect.DeepEqual(ids, []string{rowA.RuleID}) {
			t.Errorf("payload names %v; want exactly the unevaluable row %q",
				ids, rowA.RuleID)
		}
	})

	t.Run("4_exists_false_over_the_absent_key_decides_true", func(t *testing.T) {
		row := guardedRow("rdr.mvv.exists", "flows/rdr.toml:204", existsAll(key, false))
		escape := conformingEscapeRow("rdr.mvv.escape", "flows/rdr.toml:201", resolve.KindNoMatch)

		seam := newValueSeam()
		plan := mustPlan(t, guardInput("rev-mvv-4", seam, row, escape))

		if plan.RuleID != row.RuleID || plan.Escaped {
			t.Errorf("plan names %q (escaped=%v); want the guarded row %q selected — "+
				"the one sanctioned route from absence to a verdict",
				plan.RuleID, plan.Escaped, row.RuleID)
		}
		if calls := seam.seen(); len(calls) != 0 {
			t.Errorf("seam saw %v; an existence atom is decided by the kernel", calls)
		}
	})

	// The controls the `oracle` mini-check names: each scenario has a
	// sibling that must come out the other way, so no leg is an
	// absence-of-error oracle.
	t.Run("oracle_control_scenario_1_vs_2_and_4", func(t *testing.T) {
		row := guardedRow("rdr.mvv.control", "flows/rdr.toml:205", valueAtom)
		row.RequiresOwned = []string{"status"}

		in := guardInput("rev-mvv-control", eqSeam{}, row)
		in.Observed = append(in.Observed, resolve.Tag{Key: key, Value: "3"}) // decides TRUE

		plan := mustPlan(t, in)
		if plan.RuleID != row.RuleID {
			t.Errorf("plan names %q; want %q — the same table with the key present "+
				"and the value TRUE plans, which is what makes scenario 1's refusal "+
				"discriminating", plan.RuleID, row.RuleID)
		}
	})
}

// --- helpers -------------------------------------------------------------

func boolVerdict(b bool) resolve.GuardResult {
	if b {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}

func parseInt(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n := 0
	neg := false
	for i, c := range s {
		if i == 0 && c == '-' {
			neg = true
			continue
		}
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		n = -n
	}
	return n, true
}

// parseD13Set reads the §D13 carriage form: a canonical JSON array of
// strings, members sorted, duplicate-free, compact encoding.
func parseD13Set(s string) ([]string, bool) {
	if len(s) < 2 || s[0] != '[' || s[len(s)-1] != ']' {
		return nil, false
	}
	body := s[1 : len(s)-1]
	if body == "" {
		return nil, true
	}
	var out []string
	for _, part := range strings.Split(body, ",") {
		if len(part) < 2 || part[0] != '"' || part[len(part)-1] != '"' {
			return nil, false
		}
		out = append(out, part[1:len(part)-1])
	}
	return out, true
}

// fieldDoc returns the doc comment on a named struct field of the kernel
// package.
func fieldDoc(t *testing.T, typeName, fieldName string) string {
	t.Helper()
	for _, f := range parseKernelPackage(t) {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != typeName {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range st.Fields.List {
					for _, name := range field.Names {
						if name.Name == fieldName && field.Doc != nil {
							return field.Doc.Text()
						}
					}
				}
			}
		}
	}
	return ""
}

// kernelPackageSource renders the kernel package's non-test sources,
// comments included, so a doc-comment obligation can be checked.
func kernelPackageSource(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	fset := token.NewFileSet()
	for _, f := range parseKernelPackage(t) {
		if err := printer.Fprint(&b, fset, f); err != nil {
			t.Fatalf("rendering the kernel package source: %v", err)
		}
		for _, group := range f.Comments {
			b.WriteString(group.Text())
		}
	}
	return b.String()
}

// hasText reports whether needle appears in haystack once both are
// whitespace-normalized, so a doc-comment obligation is not evaded by
// re-wrapping the line.
func hasText(haystack, needle string) bool {
	return strings.Contains(squash(haystack), squash(needle))
}

// squash collapses every run of whitespace and comment-marker noise to one
// space.
func squash(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, "//", " ")), " ")
}
