//go:build rdr0012probe

package guard_test

// RDR 0012 — the behavioural probes over the declared-kind carrier.
//
// Every case in this file names a surface RDR 0012 ADDS and the tree these
// tests were written against does not yet declare: `guard.NewEvaluator`,
// `guard.DeclaredKinds`, `resolve.ContractKinds`, and the constructor-taking
// `resolve.TestGuardEvaluatorContract`. Naming them directly in an ordinary
// test file would break the build, and `.githooks/pre-commit` runs `go vet
// ./...` under `set -e`, so a red suite written that way cannot be
// committed.
//
// The file is therefore fenced behind the `rdr0012probe` build tag, which
// neither `go vet` nor the default `go test` selects, and is driven from
// `seam_carrier_0012_test.go` by a `go test -tags rdr0012probe` child. A
// probe that does not compile is reported there as the REQ it owes — a
// runtime assertion failure, never a build break. Once the surface lands
// the probes compile, run, and must pass on their own merits.
//
// Probes whose name carries `_Mutant_` are EXPECTED to fail: each drives the
// published conformance suite against a deliberately wrong seam, and the
// driver asserts that the suite caught it.

import (
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/cwensel/intrastate/internal/guard"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// The five-kind vocabulary, RDR 0003's spelling.
const (
	r12KindEnum   = "enum"
	r12KindBool   = "bool"
	r12KindInt    = "int"
	r12KindSet    = "set"
	r12KindScalar = "scalar"
)

// r12Vocabulary is the five-kind vocabulary as a set.
var r12Vocabulary = map[string]bool{
	r12KindEnum: true, r12KindBool: true, r12KindInt: true,
	r12KindSet: true, r12KindScalar: true,
}

// r12Kinds is this file's OWN declaration mapping: one key per kind token.
// It is deliberately independent of `resolve.ContractKinds`, so the seam's
// dispatch is pinned whatever the published fixture's keys are.
func r12Kinds() map[string]string {
	return map[string]string{
		"n": r12KindInt,
		"b": r12KindBool,
		"e": r12KindEnum,
		"s": r12KindScalar,
		"t": r12KindSet,
	}
}

func r12Atom(key, op, literal string) resolve.GuardAtom {
	return resolve.GuardAtom{Key: key, Operator: op, Literal: literal, Block: resolve.BlockAll}
}

// r12Case is one row of a seam table.
type r12Case struct {
	name  string
	key   string
	op    string
	lit   string
	value string
	want  resolve.GuardResult
}

// r12Eval evaluates one atom and converts a panic into a reported failure:
// the seam NEVER panics (`0012:C2`).
func r12Eval(t *testing.T, seam resolve.GuardEvaluator, atom resolve.GuardAtom, value string) (got resolve.GuardResult) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Evaluate(%+v, %q) panicked: %v; the seam never panics", atom, value, r)
			got = resolve.GuardUnevaluable
		}
	}()
	return seam.Evaluate(atom, value)
}

// r12RunCases drives a table against a seam built by NewEvaluator.
func r12RunCases(t *testing.T, kinds map[string]string, cases []r12Case) {
	t.Helper()

	seam := guard.NewEvaluator(kinds)
	for _, tc := range cases {
		atom := r12Atom(tc.key, tc.op, tc.lit)
		got := r12Eval(t, seam, atom, tc.value)
		if got == tc.want {
			continue
		}
		if tc.want == resolve.GuardUnevaluable && got == resolve.GuardFalse {
			t.Errorf("%s: Evaluate(%+v, %q) = GuardFalse; want GuardUnevaluable — a "+
				"present value the declared kind cannot parse is unevaluable, never false",
				tc.name, atom, tc.value)
			continue
		}
		t.Errorf("%s: Evaluate(%+v, %q) = %v; want %v", tc.name, atom, tc.value, got, tc.want)
	}
}

// --- C1: the carrier ------------------------------------------------------

// REQ-1 / REQ-10: `func NewEvaluator(kinds map[string]string) Evaluator` —
// the concrete type, which satisfies the kernel's seam interface.
func TestRDR0012Probe_Req1_NewEvaluatorConstructsTheConcreteSeam(t *testing.T) {
	var concrete guard.Evaluator = guard.NewEvaluator(r12Kinds())
	var seam resolve.GuardEvaluator = concrete

	// The mapping is what it carries: an `int` key compares parsed values,
	// which a mapping-free seam cannot do.
	if got := r12Eval(t, seam, r12Atom("n", "eq", "7"), "07"); got != resolve.GuardTrue {
		t.Errorf("NewEvaluator(kinds) over `n` (int): eq 7 held \"07\" = %v; want "+
			"GuardTrue — the constructed evaluator must hold the mapping it was built over", got)
	}
}

// REQ-2: `DeclaredKinds(m *table.Model) map[string]string` maps tag key →
// declared kind token.
func TestRDR0012Probe_Req2_DeclaredKindsMapsKeyToKindToken(t *testing.T) {
	m := &table.Model{Tags: map[string]table.TagDecl{
		"n": {Kind: r12KindInt},
		"b": {Kind: r12KindBool},
		"e": {Kind: r12KindEnum},
		"s": {Kind: r12KindScalar},
		"t": {Kind: r12KindSet},
	}}
	got := guard.DeclaredKinds(m)
	if want := r12Kinds(); !reflect.DeepEqual(got, want) {
		t.Errorf("DeclaredKinds = %v; want %v", got, want)
	}
}

// REQ-3: TOTAL over `m.Tags`, and an empty-string `Kind` is OMITTED rather
// than mapped to "".
func TestRDR0012Probe_Req3_DeclaredKindsIsTotalAndOmitsEmptyKind(t *testing.T) {
	m := &table.Model{Tags: map[string]table.TagDecl{
		"n":     {Kind: r12KindInt},
		"b":     {Kind: r12KindBool},
		"blank": {Kind: ""},
	}}
	got := guard.DeclaredKinds(m)
	if _, present := got["blank"]; present {
		t.Errorf("DeclaredKinds carries an entry for the empty-Kind key `blank` (%q); "+
			"it must be OMITTED — an empty-string entry is indistinguishable from an "+
			"absent key at C2's arms", got["blank"])
	}
	if got["n"] != r12KindInt || got["b"] != r12KindBool {
		t.Errorf("DeclaredKinds = %v; every declared key with a kind gets an entry", got)
	}
	if len(got) != 2 {
		t.Errorf("DeclaredKinds = %v (%d entries); want exactly the 2 keys that "+
			"declare a kind", got, len(got))
	}
}

// REQ-4: a nil model yields an empty non-nil map, never a panic.
func TestRDR0012Probe_Req4_DeclaredKindsOverNilModelIsEmptyNonNil(t *testing.T) {
	var got map[string]string
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("DeclaredKinds(nil) panicked: %v", r)
			}
		}()
		got = guard.DeclaredKinds(nil)
	}()
	if got == nil {
		t.Fatal("DeclaredKinds(nil) = nil; want an empty NON-NIL map")
	}
	if len(got) != 0 {
		t.Errorf("DeclaredKinds(nil) = %v; want empty", got)
	}
}

// REQ-8: the seam signature and the atom shape are unchanged, and no kind
// rides on the atom.
func TestRDR0012Probe_Req8_SeamSignatureAndAtomShapeAreUnchanged(t *testing.T) {
	var seam resolve.GuardEvaluator = guard.NewEvaluator(r12Kinds())

	m, ok := reflect.TypeOf(seam).MethodByName("Evaluate")
	if !ok {
		t.Fatal("NewEvaluator's value has no Evaluate method")
	}
	want := reflect.TypeOf(func(resolve.GuardAtom, string) resolve.GuardResult { return 0 })
	// Method type carries the receiver first; compare the tail.
	if m.Type.NumIn() != 3 || m.Type.In(1) != want.In(0) || m.Type.In(2) != want.In(1) ||
		m.Type.NumOut() != 1 || m.Type.Out(0) != want.Out(0) {
		t.Errorf("Evaluate has type %v; want Evaluate(atom GuardAtom, value string) GuardResult", m.Type)
	}

	atomT := reflect.TypeOf(resolve.GuardAtom{})
	var fields []string
	for i := range atomT.NumField() {
		fields = append(fields, atomT.Field(i).Name)
	}
	if !slices.Equal(fields, []string{"Key", "Operator", "Literal", "Block"}) {
		t.Errorf("GuardAtom fields = %v; want exactly [Key Operator Literal Block] — "+
			"no kind is ever stamped into an atom", fields)
	}
}

// --- C2: the typed comparison matrix --------------------------------------

// REQ-12: `int` compares PARSED values; an unparseable held value is
// unevaluable, never false.
func TestRDR0012Probe_Req12_IntComparesParsedValues(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"eq/equal", "n", "eq", "7", "7", resolve.GuardTrue},
		{"eq/leading zero parses equal", "n", "eq", "7", "07", resolve.GuardTrue},
		{"eq/plus sign parses equal", "n", "eq", "7", "+7", resolve.GuardTrue},
		{"eq/parses but differs", "n", "eq", "7", "8", resolve.GuardFalse},
		{"eq/unparseable held", "n", "eq", "7", "many", resolve.GuardUnevaluable},
		{"eq/empty held", "n", "eq", "7", "", resolve.GuardUnevaluable},
		{"in/member parsed", "n", "in", `["7","9"]`, "09", resolve.GuardTrue},
		{"in/non-member", "n", "in", `["7","9"]`, "8", resolve.GuardFalse},
		{"in/unparseable held", "n", "in", `["7","9"]`, "many", resolve.GuardUnevaluable},
	})
}

// REQ-13: ONE unparseable `in` member poisons the whole list.
func TestRDR0012Probe_Req13_OneUnparseableInMemberPoisonsTheList(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"held equals a parsed member", "n", "in", `["7","x"]`, "7", resolve.GuardUnevaluable},
		{"held equals the poisoned member's text", "n", "in", `["7","x"]`, "x", resolve.GuardUnevaluable},
		{"held matches nothing", "n", "in", `["7","x"]`, "3", resolve.GuardUnevaluable},
		{"poison first, match later", "n", "in", `["x","7"]`, "7", resolve.GuardUnevaluable},
	})
}

// REQ-14: the parse is `strconv.Atoi`; an out-of-range spelling is
// unevaluable. The overflow spelling overflows 64-bit int, so it overflows
// on every supported target.
func TestRDR0012Probe_Req14_OverflowIsUnevaluableAndInRangeIsPinned(t *testing.T) {
	const overflow = "99999999999999999999"
	if _, err := strconv.Atoi(overflow); !errors.Is(err, strconv.ErrRange) {
		t.Fatalf("fixture %q does not overflow int on this target (%v)", overflow, err)
	}
	r12RunCases(t, r12Kinds(), []r12Case{
		{"eq/overflow held", "n", "eq", "7", overflow, resolve.GuardUnevaluable},
		{"eq/negative overflow held", "n", "eq", "7", "-" + overflow, resolve.GuardUnevaluable},
		{"in/overflow held", "n", "in", `["7"]`, overflow, resolve.GuardUnevaluable},
		{"eq/in-range negative", "n", "eq", "-3", "-3", resolve.GuardTrue},
		{"eq/in-range max int32", "n", "eq", "2147483647", "2147483647", resolve.GuardTrue},
	})
}

// REQ-15: `bool` admits only the tokens `true` | `false`.
func TestRDR0012Probe_Req15_BoolComparesTokensOnly(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"eq/true", "b", "eq", "true", "true", resolve.GuardTrue},
		{"eq/false vs true", "b", "eq", "true", "false", resolve.GuardFalse},
		{"eq/false", "b", "eq", "false", "false", resolve.GuardTrue},
		{"eq/1", "b", "eq", "true", "1", resolve.GuardUnevaluable},
		{"eq/0", "b", "eq", "false", "0", resolve.GuardUnevaluable},
		{"eq/True", "b", "eq", "true", "True", resolve.GuardUnevaluable},
		{"eq/yes", "b", "eq", "true", "yes", resolve.GuardUnevaluable},
		{"in/member", "b", "in", `["false","true"]`, "true", resolve.GuardTrue},
		{"in/non-member", "b", "in", `["true"]`, "false", resolve.GuardFalse},
		{"in/yes", "b", "in", `["true"]`, "yes", resolve.GuardUnevaluable},
	})
}

// REQ-16: `enum` and `scalar` share one arm — exact string equality, no
// parse.
func TestRDR0012Probe_Req16_EnumAndScalarCompareExactStrings(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"enum eq/equal", "e", "eq", "Draft", "Draft", resolve.GuardTrue},
		{"enum eq/case differs", "e", "eq", "Draft", "draft", resolve.GuardFalse},
		{"enum eq/numeric text is not parsed", "e", "eq", "7", "07", resolve.GuardFalse},
		{"enum in/member", "e", "in", `["alpha","beta"]`, "beta", resolve.GuardTrue},
		{"enum in/non-member", "e", "in", `["alpha","beta"]`, "gamma", resolve.GuardFalse},
		{"scalar eq/equal", "s", "eq", "x y", "x y", resolve.GuardTrue},
		{"scalar eq/numeric text is not parsed", "s", "eq", "7", "07", resolve.GuardFalse},
		{"scalar eq/bool text is not parsed", "s", "eq", "true", "True", resolve.GuardFalse},
		{"scalar in/member", "s", "in", `["a","b"]`, "a", resolve.GuardTrue},
	})
}

// REQ-17: `eq`/`in` over a `set` key is unevaluable.
func TestRDR0012Probe_Req17_SetEqAndInAreUnevaluable(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"eq/value equals literal", "t", "eq", "alpha", "alpha", resolve.GuardUnevaluable},
		{"eq/array value equals literal", "t", "eq", `["alpha"]`, `["alpha"]`, resolve.GuardUnevaluable},
		{"in/member", "t", "in", `["alpha","beta"]`, "alpha", resolve.GuardUnevaluable},
	})
}

// REQ-18: a key absent from the mapping is unevaluable.
func TestRDR0012Probe_Req18_AbsentKeyIsUnevaluable(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"eq/value equals literal", "undeclared", "eq", "x", "x", resolve.GuardUnevaluable},
		{"in/member", "undeclared", "in", `["x"]`, "x", resolve.GuardUnevaluable},
	})
	r12RunCases(t, map[string]string{}, []r12Case{
		{"empty mapping eq", "n", "eq", "7", "7", resolve.GuardUnevaluable},
	})
}

// REQ-19: a key DECLARED with `Kind` "" reaches the seam as an absent key —
// `eq`/`in` over it are unevaluable where the raw-string seam compared.
func TestRDR0012Probe_Req19_EmptyKindDeclaredKeyIsUnevaluable(t *testing.T) {
	m := &table.Model{Tags: map[string]table.TagDecl{
		"n":     {Kind: r12KindInt},
		"blank": {Kind: ""},
	}}
	kinds := guard.DeclaredKinds(m)
	if _, present := kinds["blank"]; present {
		t.Errorf("DeclaredKinds maps the empty-Kind key `blank`; it must be omitted")
	}
	r12RunCases(t, kinds, []r12Case{
		{"eq over empty-Kind key", "blank", "eq", "x", "x", resolve.GuardUnevaluable},
		{"in over empty-Kind key", "blank", "in", `["x"]`, "x", resolve.GuardUnevaluable},
	})
}

// REQ-20: an unparseable LITERAL is unevaluable (defense in depth).
func TestRDR0012Probe_Req20_UnparseableLiteralIsUnevaluable(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"int eq/literal not an int", "n", "eq", "seven", "7", resolve.GuardUnevaluable},
		{"int eq/literal equals held text", "n", "eq", "seven", "seven", resolve.GuardUnevaluable},
		{"bool eq/literal not a token", "b", "eq", "yes", "true", resolve.GuardUnevaluable},
		{"bool eq/literal equals held text", "b", "eq", "yes", "yes", resolve.GuardUnevaluable},
		{"bool in/member not a token", "b", "in", `["true","yes"]`, "true", resolve.GuardUnevaluable},
	})
}

// REQ-21: dispatch is exhaustive over the five tokens and the `default`
// arm is unevaluable.
func TestRDR0012Probe_Req21_UnknownKindTokenIsUnevaluable(t *testing.T) {
	kinds := r12Kinds()
	kinds["d"] = "decimal"
	kinds["cap"] = "Int"
	kinds["empty"] = ""
	r12RunCases(t, kinds, []r12Case{
		{"unknown token eq", "d", "eq", "7", "7", resolve.GuardUnevaluable},
		{"unknown token in", "d", "in", `["7"]`, "7", resolve.GuardUnevaluable},
		{"mis-cased token eq", "cap", "eq", "7", "7", resolve.GuardUnevaluable},
		{"empty-string token eq", "empty", "eq", "x", "x", resolve.GuardUnevaluable},
		// Every vocabulary token that admits eq DECIDES a matching value:
		// the default arm is not reached for them.
		{"int decides", "n", "eq", "7", "7", resolve.GuardTrue},
		{"bool decides", "b", "eq", "true", "true", resolve.GuardTrue},
		{"enum decides", "e", "eq", "a", "a", resolve.GuardTrue},
		{"scalar decides", "s", "eq", "a", "a", resolve.GuardTrue},
	})
}

// REQ-22: `lt/lte/gt/gte` and `contains` do not consult the kind lookup;
// a nil-mapping evaluator still answers them normally.
func TestRDR0012Probe_Req22_OrderingAndContainsDoNotConsultTheMapping(t *testing.T) {
	cases := []r12Case{
		{"gte on int key", "n", "gte", "3", "4", resolve.GuardTrue},
		{"gte below", "n", "gte", "3", "2", resolve.GuardFalse},
		{"lt on absent key", "zz", "lt", "3", "2", resolve.GuardTrue},
		{"gt on absent key", "zz", "gt", "3", "2", resolve.GuardFalse},
		{"lte unparseable held", "zz", "lte", "3", "many", resolve.GuardUnevaluable},
		{"contains on set key", "t", "contains", `["alpha"]`, `["alpha","beta"]`, resolve.GuardTrue},
		{"contains missing member", "t", "contains", `["gamma"]`, `["alpha","beta"]`, resolve.GuardFalse},
		{"contains on absent key", "zz", "contains", `["alpha"]`, `["alpha"]`, resolve.GuardTrue},
		{"contains on enum key", "e", "contains", `["alpha"]`, `["alpha"]`, resolve.GuardTrue},
	}
	t.Run("constructed mapping", func(t *testing.T) { r12RunCases(t, r12Kinds(), cases) })
	t.Run("nil mapping", func(t *testing.T) { r12RunCases(t, nil, cases) })
}

// REQ-23: an unrecognized operator — `exists` included — is unevaluable,
// and the seam never panics.
func TestRDR0012Probe_Req23_UnknownOperatorIsUnevaluableAndNeverPanics(t *testing.T) {
	var cases []r12Case
	for _, key := range []string{"n", "b", "e", "s", "t", "zz"} {
		for _, op := range []string{"exists", "matches", "", "EQ", "ne"} {
			cases = append(cases, r12Case{
				name: key + "/" + op, key: key, op: op, lit: "true", value: "true",
				want: resolve.GuardUnevaluable,
			})
		}
	}
	r12RunCases(t, r12Kinds(), cases)
	r12RunCases(t, nil, cases)
}

// REQ-24: `in`'s members are the members decoded from the array; no
// delimiter convention is introduced.
func TestRDR0012Probe_Req24_InMembersAreTheDecodedArray(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		{"int decoded member parsed", "n", "in", `["1","2"]`, "02", resolve.GuardTrue},
		{"int comma-delimited literal is not a list", "n", "in", "1,2", "1", resolve.GuardUnevaluable},
		{"int pipe-delimited literal is not a list", "n", "in", "1|2", "1", resolve.GuardUnevaluable},
		{"bool decoded member", "b", "in", `["true"]`, "true", resolve.GuardTrue},
		{"bool comma-delimited literal is not a list", "b", "in", "true,false", "true", resolve.GuardUnevaluable},
	})
}

// REQ-46 (`0012:S2`): the per-kind dispatch matrix over `eq` and `in`,
// driven through NewEvaluator — the C2 matrix verbatim.
func TestRDR0012Probe_Req46_PerKindDispatchMatrix(t *testing.T) {
	r12RunCases(t, r12Kinds(), []r12Case{
		// int
		{"int eq parsed-equal", "n", "eq", "7", "07", resolve.GuardTrue},
		{"int eq differs", "n", "eq", "7", "4", resolve.GuardFalse},
		{"int eq unparseable", "n", "eq", "7", "many", resolve.GuardUnevaluable},
		{"int in parsed member", "n", "in", `["4","7"]`, "007", resolve.GuardTrue},
		{"int in unparseable", "n", "in", `["4","7"]`, "many", resolve.GuardUnevaluable},
		// bool
		{"bool eq token", "b", "eq", "false", "false", resolve.GuardTrue},
		{"bool eq other token", "b", "eq", "false", "true", resolve.GuardFalse},
		{"bool eq non-token", "b", "eq", "false", "0", resolve.GuardUnevaluable},
		{"bool in non-token", "b", "in", `["true"]`, "1", resolve.GuardUnevaluable},
		// enum
		{"enum eq", "e", "eq", "a", "a", resolve.GuardTrue},
		{"enum eq differs", "e", "eq", "a", "b", resolve.GuardFalse},
		{"enum in", "e", "in", `["a","b"]`, "b", resolve.GuardTrue},
		// scalar
		{"scalar eq", "s", "eq", "07", "07", resolve.GuardTrue},
		{"scalar eq no parse", "s", "eq", "7", "07", resolve.GuardFalse},
		{"scalar in", "s", "in", `["x"]`, "y", resolve.GuardFalse},
		// set
		{"set eq", "t", "eq", "x", "x", resolve.GuardUnevaluable},
		{"set in", "t", "in", `["x"]`, "x", resolve.GuardUnevaluable},
		// key absent
		{"absent eq", "zz", "eq", "x", "x", resolve.GuardUnevaluable},
		{"absent in", "zz", "in", `["x"]`, "x", resolve.GuardUnevaluable},
	})
}

// --- C3: the conformance suite --------------------------------------------

// r12NewSeam is the one-line adapter C3 names.
func r12NewSeam(k map[string]string) resolve.GuardEvaluator { return guard.NewEvaluator(k) }

// MVV step 4: the extended suite, driven against NewEvaluator over the
// published fixture kinds, is green.
func TestRDR0012Probe_ReqMVV_ContractSuiteIsGreenOverNewEvaluator(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12NewSeam)
}

// REQ-26: the published fixture's kind tokens are exactly the five-kind
// vocabulary.
func TestRDR0012Probe_Req26_FixtureTokensAreExactlyTheVocabulary(t *testing.T) {
	tokens := map[string]bool{}
	for _, kind := range resolve.ContractKinds() {
		tokens[kind] = true
	}
	if !reflect.DeepEqual(tokens, r12Vocabulary) {
		t.Errorf("ContractKinds() carries kind tokens %v; want exactly the five-kind "+
			"vocabulary %v", tokens, r12Vocabulary)
	}
}

// r12RawIn is today's raw-string `in`: decode the array and test string
// membership, no kind consulted.
func r12RawIn(literal, value string) resolve.GuardResult {
	var members []string
	if err := json.Unmarshal([]byte(literal), &members); err != nil || len(members) == 0 {
		return resolve.GuardUnevaluable
	}
	if slices.Contains(members, value) {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}

func r12RawEqIn(atom resolve.GuardAtom, value string) resolve.GuardResult {
	if atom.Operator == "in" {
		return r12RawIn(atom.Literal, value)
	}
	if value == atom.Literal {
		return resolve.GuardTrue
	}
	return resolve.GuardFalse
}

// r12MutantSeam is NewEvaluator with ONE kind's eq/in arm replaced by the
// wrong behaviour: the raw-string comparison for kinds that must parse or
// refuse, a blanket unevaluable for the kinds that must compare.
type r12MutantSeam struct {
	inner resolve.GuardEvaluator
	kinds map[string]string
	kind  string
}

func (s r12MutantSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	if (atom.Operator == "eq" || atom.Operator == "in") && s.kinds[atom.Key] == s.kind {
		switch s.kind {
		case r12KindInt, r12KindBool, r12KindSet:
			return r12RawEqIn(atom, value)
		default:
			return resolve.GuardUnevaluable
		}
	}
	return s.inner.Evaluate(atom, value)
}

func r12MutantOf(kind string) func(map[string]string) resolve.GuardEvaluator {
	return func(k map[string]string) resolve.GuardEvaluator {
		return r12MutantSeam{inner: guard.NewEvaluator(k), kinds: k, kind: kind}
	}
}

// REQ-26: at least one DISCRIMINATING case per kind token. Each mutant
// breaks exactly one kind's arm; the driver asserts the suite FAILS it.
func TestRDR0012Probe_Req26_Mutant_int(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12MutantOf(r12KindInt))
}

func TestRDR0012Probe_Req26_Mutant_bool(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12MutantOf(r12KindBool))
}

func TestRDR0012Probe_Req26_Mutant_enum(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12MutantOf(r12KindEnum))
}

func TestRDR0012Probe_Req26_Mutant_scalar(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12MutantOf(r12KindScalar))
}

func TestRDR0012Probe_Req26_Mutant_set(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, r12MutantOf(r12KindSet))
}

// r12RawSeam is TODAY's raw-string seam: eq/in never consult a kind.
type r12RawSeam struct{ inner resolve.GuardEvaluator }

func (s r12RawSeam) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	if atom.Operator == "eq" || atom.Operator == "in" {
		return r12RawEqIn(atom, value)
	}
	return s.inner.Evaluate(atom, value)
}

// MVV step 4 / REQ-25: the int/bool want-Unevaluable legs FAIL against
// today's raw-string arms. The driver asserts the suite fails this seam.
func TestRDR0012Probe_ReqMVV_Mutant_rawStringSeam(t *testing.T) {
	resolve.TestGuardEvaluatorContract(t, func(k map[string]string) resolve.GuardEvaluator {
		return r12RawSeam{inner: guard.NewEvaluator(k)}
	})
}

// REQ-27: ContractKinds is a function returning a FRESH map per call.
func TestRDR0012Probe_Req27_ContractKindsHandsEachCallerItsOwnCopy(t *testing.T) {
	first := resolve.ContractKinds()
	if len(first) == 0 {
		t.Fatal("ContractKinds() is empty")
	}
	n := len(first)
	for k := range first {
		first[k] = "mutated"
	}
	first["injected-by-caller"] = r12KindInt

	second := resolve.ContractKinds()
	if len(second) != n {
		t.Errorf("second ContractKinds() has %d entries; want %d — a caller's "+
			"mutation leaked into the published fixture", len(second), n)
	}
	if _, leaked := second["injected-by-caller"]; leaked {
		t.Error("a key one caller added appears in the next caller's fixture")
	}
	for k, v := range second {
		if v == "mutated" {
			t.Errorf("fixture key %q carries a value a previous caller wrote", k)
		}
	}
}

// r12Call is one question the suite asked, with the kind the asking seam
// held for the key and the constructed evaluator's answer.
type r12Call struct {
	atom    resolve.GuardAtom
	value   string
	kind    string
	mapped  bool
	fixture bool
	verdict resolve.GuardResult
}

type r12Recorder struct {
	inner   resolve.GuardEvaluator
	kinds   map[string]string
	fixture map[string]string
	log     *[]r12Call
}

func (r r12Recorder) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult {
	got := r.inner.Evaluate(atom, value)
	kind, mapped := r.kinds[atom.Key]
	_, inFixture := r.fixture[atom.Key]
	*r.log = append(*r.log, r12Call{
		atom: atom, value: value, kind: kind, mapped: mapped, fixture: inFixture, verdict: got,
	})
	return got
}

// r12Record runs the suite against NewEvaluator wrapped in a recorder and
// returns what it asked.
func r12Record(t *testing.T) ([]r12Call, map[string]string) {
	t.Helper()

	fixture := resolve.ContractKinds()
	var calls []r12Call
	resolve.TestGuardEvaluatorContract(t, func(k map[string]string) resolve.GuardEvaluator {
		copied := make(map[string]string, len(k))
		for key, v := range k {
			copied[key] = v
		}
		return r12Recorder{inner: guard.NewEvaluator(k), kinds: copied, fixture: fixture, log: &calls}
	})
	if len(calls) == 0 {
		t.Fatal("the conformance suite asked the constructed seam nothing")
	}
	return calls, fixture
}

func r12Parses(kind, value string) bool {
	switch kind {
	case r12KindInt:
		_, err := strconv.Atoi(value)
		return err == nil
	case r12KindBool:
		return value == "true" || value == "false"
	}
	return true
}

func r12Overflows(value string) bool {
	_, err := strconv.Atoi(value)
	return errors.Is(err, strconv.ErrRange)
}

// REQ-25 (`0012:S1`): the new legs are asked, and the constructed seam
// answers each as the contract states.
func TestRDR0012Probe_Req25_TheSuiteAsksTheNewLegs(t *testing.T) {
	calls, _ := r12Record(t)

	has := func(pred func(c r12Call) bool) bool {
		return slices.ContainsFunc(calls, pred)
	}
	for _, kind := range []string{r12KindInt, r12KindBool} {
		for _, op := range []string{"eq", "in"} {
			if !has(func(c r12Call) bool {
				return c.kind == kind && c.atom.Operator == op && !r12Parses(kind, c.value) &&
					c.verdict == resolve.GuardUnevaluable
			}) {
				t.Errorf("the suite never asks %s over a %s key with a held value that does "+
					"not parse (want GuardUnevaluable)", op, kind)
			}
			if !has(func(c r12Call) bool {
				return c.kind == kind && c.atom.Operator == op && r12Parses(kind, c.value) &&
					c.verdict == resolve.GuardFalse
			}) {
				t.Errorf("the suite never asks %s over a %s key with a held value that "+
					"parses but differs (want GuardFalse)", op, kind)
			}
		}
	}
	if !has(func(c r12Call) bool {
		if c.kind != r12KindInt || c.atom.Operator != "in" || c.verdict != resolve.GuardUnevaluable {
			return false
		}
		var members []string
		if json.Unmarshal([]byte(c.atom.Literal), &members) != nil {
			return false
		}
		return slices.ContainsFunc(members, func(m string) bool { return !r12Parses(r12KindInt, m) }) &&
			r12Parses(r12KindInt, c.value)
	}) {
		t.Error("the suite never asks `in` over an int key whose list carries a non-integer " +
			"member while the held value parses (want GuardUnevaluable)")
	}
	if !has(func(c r12Call) bool {
		return c.kind == r12KindInt && r12Overflows(c.value) && c.verdict == resolve.GuardUnevaluable
	}) {
		t.Error("the suite never asks an int key with an integer-overflow held value " +
			"(want GuardUnevaluable)")
	}
	if !has(func(c r12Call) bool {
		return c.kind == r12KindSet && c.atom.Operator == "eq" && c.verdict == resolve.GuardUnevaluable
	}) {
		t.Error("the suite never asks `eq` over a set-kind key (want GuardUnevaluable)")
	}
	if !has(func(c r12Call) bool {
		return !c.mapped && !c.fixture && (c.atom.Operator == "eq" || c.atom.Operator == "in") &&
			c.verdict == resolve.GuardUnevaluable
	}) {
		t.Error("the suite never asks eq/in over a key absent from the fixture " +
			"(want GuardUnevaluable)")
	}
}

// REQ-21 (suite side): the suite pins an unknown kind token as unevaluable.
func TestRDR0012Probe_Req21_TheSuitePinsAnUnknownKindToken(t *testing.T) {
	calls, _ := r12Record(t)
	if !slices.ContainsFunc(calls, func(c r12Call) bool {
		return c.mapped && c.kind != "" && !r12Vocabulary[c.kind] &&
			(c.atom.Operator == "eq" || c.atom.Operator == "in") &&
			c.verdict == resolve.GuardUnevaluable
	}) {
		t.Error("the suite never asks eq/in over a key mapped to a token outside the " +
			"five-kind vocabulary (want GuardUnevaluable)")
	}
}

// REQ-29: each case names its own fixture key, admitted by the matrix for
// its operator; the re-keyed legacy cases keep their verdicts.
func TestRDR0012Probe_Req29_CasesAreKeyedOntoAdmittedKinds(t *testing.T) {
	calls, fixture := r12Record(t)

	for _, c := range calls {
		if c.atom.Key == "subject" {
			if _, declared := fixture["subject"]; !declared {
				t.Errorf("case %+v still uses Key \"subject\", which the fixture does not declare", c.atom)
			}
		}
		if !c.fixture && c.verdict != resolve.GuardUnevaluable {
			t.Errorf("case %+v names key %q, which the fixture does not declare, and "+
				"was answered %v", c.atom, c.atom.Key, c.verdict)
		}
		switch c.atom.Operator {
		case "gte", "gt", "lt", "lte":
			if c.fixture && fixture[c.atom.Key] != r12KindInt {
				t.Errorf("ordering case %+v is keyed onto a %q key; the matrix admits only int",
					c.atom, fixture[c.atom.Key])
			}
		case "contains":
			if c.fixture && fixture[c.atom.Key] != r12KindSet {
				t.Errorf("contains case %+v is keyed onto a %q key; the matrix admits only set",
					c.atom, fixture[c.atom.Key])
			}
		}
	}

	legacy := []struct {
		op, lit, value string
		kinds          []string
		want           resolve.GuardResult
	}{
		{"eq", "Draft", "Draft", []string{r12KindEnum, r12KindScalar}, resolve.GuardTrue},
		{"eq", "Draft", "Final", []string{r12KindEnum, r12KindScalar}, resolve.GuardFalse},
		{"gte", "3", "4", []string{r12KindInt}, resolve.GuardTrue},
		{"gte", "3", "3", []string{r12KindInt}, resolve.GuardTrue},
		{"gte", "3", "2", []string{r12KindInt}, resolve.GuardFalse},
		{"gte", "3", "many", []string{r12KindInt}, resolve.GuardUnevaluable},
		{"gte", "three", "3", []string{r12KindInt}, resolve.GuardUnevaluable},
		{"in", `["alpha","beta"]`, "alpha", []string{r12KindEnum, r12KindScalar}, resolve.GuardTrue},
		{"in", `["alpha","beta"]`, "gamma", []string{r12KindEnum, r12KindScalar}, resolve.GuardFalse},
		{"contains", `["alpha"]`, `["alpha","beta"]`, []string{r12KindSet}, resolve.GuardTrue},
		{"contains", `["gamma"]`, `["alpha","beta"]`, []string{r12KindSet}, resolve.GuardFalse},
		{"contains", `["alpha"]`, "alpha", []string{r12KindSet}, resolve.GuardUnevaluable},
		{"contains", `["alpha"]`, "null", []string{r12KindSet}, resolve.GuardUnevaluable},
	}
	for _, l := range legacy {
		found := slices.ContainsFunc(calls, func(c r12Call) bool {
			return c.atom.Operator == l.op && c.atom.Literal == l.lit && c.value == l.value &&
				c.fixture && slices.Contains(l.kinds, fixture[c.atom.Key]) && c.verdict == l.want
		})
		if !found {
			t.Errorf("legacy case %s %s held %q is not asked over a %v fixture key with "+
				"verdict %v; re-keying must keep every existing case and its verdict",
				l.op, l.lit, l.value, l.kinds, l.want)
		}
	}
}
