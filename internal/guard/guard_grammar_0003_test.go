package guard_test

// RDR 0003 — atom grammar, the closed operator vocabulary, the
// operator/kind matrix, and the two equalities.
//
// Nothing here mocks the unit under test: every assertion drives the real
// guard package. The kernel (`internal/resolve`) and the loader
// (`internal/table`) are real collaborators.

import (
	"slices"
	"strings"
	"testing"

	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-1: "A guard predicate MUST be a symbolic atom over a declared tag,
// not a host language callback and not a free-form expression string."
// HAPPY PATH
//
// The checkable content is that the surface a guard travels on is the
// four-field atom — key, operator token, literal, block — and that the
// evaluator is reached with that atom rather than with host code or an
// expression string. A callback surface would have to carry a func; an
// expression surface would have to carry an unparsed string the evaluator
// splits itself.
func TestReq1_GuardIsASymbolicAtomOverADeclaredTag(t *testing.T) {
	var ev guard.Evaluator

	atom := resolve.GuardAtom{
		Key:      "profile",
		Operator: "eq",
		Literal:  "large",
		Block:    resolve.BlockAll,
	}
	if got := ev.Evaluate(atom, "large"); got != resolve.GuardTrue {
		t.Errorf("Evaluate(%+v, %q) = %v; want GuardTrue — a symbolic atom "+
			"over a declared tag decides from its own four fields", atom, "large", got)
	}
	if got := ev.Evaluate(atom, "small"); got != resolve.GuardFalse {
		t.Errorf("Evaluate(%+v, %q) = %v; want GuardFalse", atom, "small", got)
	}

	// A free-form expression string is not an operator: the vocabulary is
	// closed, so "profile == large" is rejected as an operator token, never
	// parsed as an expression.
	expr := resolve.GuardAtom{
		Key:      "profile",
		Operator: `profile == "large"`,
		Literal:  "",
		Block:    resolve.BlockAll,
	}
	if got := ev.Evaluate(expr, "large"); got != resolve.GuardUnevaluable {
		t.Errorf("Evaluate over a free-form expression operator = %v; want "+
			"GuardUnevaluable — there is no expression grammar", got)
	}
}

// REQ-2: "The initial operator vocabulary MUST be closed and typed:
// equality, membership, bounded integer comparison, existence, and set
// containment. Unknown operators MUST be rejected during parse or lint
// before resolution."
// BOUNDARY
func TestReq2_OperatorVocabularyIsClosedAndTyped(t *testing.T) {
	want := []string{"contains", "eq", "exists", "gt", "gte", "in", "lt", "lte"}

	got := slices.Sorted(slices.Values(guard.Operators()))
	if !slices.Equal(got, want) {
		t.Fatalf("Operators() = %v; want exactly %v — equality, membership, "+
			"the four bounded integer comparisons, existence, and set containment", got, want)
	}

	for _, unknown := range []string{"neq", "matches", "regex", "startswith", "", "EQ"} {
		if guard.KnownOperator(unknown) {
			t.Errorf("KnownOperator(%q) = true; the vocabulary is closed", unknown)
		}
		// An unknown operator must be rejected before resolution, in every
		// declared kind. `Accepts` is the matrix gate lint and parse share.
		for _, kind := range guard.Kinds() {
			if guard.Accepts(unknown, kind) {
				t.Errorf("Accepts(%q, %q) = true; an unknown operator is "+
					"accepted by no kind", unknown, kind)
			}
		}
	}
}

// REQ-3: "Each operator MUST declare which tag value kinds it accepts. A
// predicate whose literal cannot be parsed as the declared tag kind MUST
// be rejected before resolution."
// HAPPY PATH
func TestReq3_EachOperatorDeclaresAcceptedKindsAndLiteralsAreKindChecked(t *testing.T) {
	for _, op := range guard.Operators() {
		accepted := 0
		for _, kind := range guard.Kinds() {
			if guard.Accepts(op, kind) {
				accepted++
			}
		}
		if accepted == 0 {
			t.Errorf("operator %q accepts no declared kind; every operator "+
				"MUST declare which kinds it accepts", op)
		}
	}

	// A literal that cannot be parsed as the declared kind is rejected
	// before resolution, at the load surface RDR 0002 carries for this
	// RDR's rejection rules.
	cat := loadCategory(t, `
[tags.iter]
provenance = "owned"
kind = "int"
min = 0
max = 9
required = true
`, `
[rule.guard.all.iter]
lt = "many"
`)
	if cat != table.CatMalformedPredicateAtom {
		t.Errorf("a non-integer literal on an int tag refused as %q; want %q "+
			"— rejected before resolution", cat, table.CatMalformedPredicateAtom)
	}
}

// REQ-4: The operator/kind matrix fixes acceptance per operator: `eq` over
// `enum`, `bool`, `int`, `scalar` with "one typed scalar"; `in` over the
// same kinds with a "non-empty typed scalar set"; "`lt`, `lte`, `gt`,
// `gte`" over `int` with "one typed integer"; `exists` over "any kind,
// provided the tag is declared optional" with a boolean literal;
// `contains` over "`set` with a declared element universe" with a
// "non-empty typed element set".
// BOUNDARY
func TestReq4_OperatorKindMatrixIsExactlyAsPublished(t *testing.T) {
	want := map[string][]string{
		"eq":       {"enum", "bool", "int", "scalar"},
		"in":       {"enum", "bool", "int", "scalar"},
		"lt":       {"int"},
		"lte":      {"int"},
		"gt":       {"int"},
		"gte":      {"int"},
		"exists":   {"enum", "bool", "int", "set", "scalar"},
		"contains": {"set"},
	}

	for op, kinds := range want {
		for _, kind := range guard.Kinds() {
			accept := slices.Contains(kinds, kind)
			if guard.Accepts(op, kind) != accept {
				t.Errorf("Accepts(%q, %q) = %v; want %v — the published "+
					"operator/kind matrix", op, kind, !accept, accept)
			}
		}
	}

	shapes := map[string]guard.Shape{
		"eq":       guard.ShapeScalar,
		"in":       guard.ShapeScalarSet,
		"lt":       guard.ShapeInteger,
		"lte":      guard.ShapeInteger,
		"gt":       guard.ShapeInteger,
		"gte":      guard.ShapeInteger,
		"exists":   guard.ShapeBoolean,
		"contains": guard.ShapeElementSet,
	}
	for op, shape := range shapes {
		if got := guard.LiteralShape(op); got != shape {
			t.Errorf("LiteralShape(%q) = %q; want %q", op, got, shape)
		}
	}

	// "non-empty typed scalar set" / "non-empty typed element set": the
	// two set-shaped literals reject an empty literal set.
	var ev guard.Evaluator
	for _, op := range []string{"in", "contains"} {
		atom := resolve.GuardAtom{Key: "k", Operator: op, Literal: "[]", Block: resolve.BlockAll}
		if got := ev.Evaluate(atom, "alpha"); got != resolve.GuardUnevaluable {
			t.Errorf("Evaluate(%q with an empty literal set) = %v; want "+
				"GuardUnevaluable — the literal shape is a NON-EMPTY set", op, got)
		}
	}
}

// REQ-5: "a key declared always-present contributes no `{absent}`
// assignment, so an `exists` atom over it is well-formed but vacuous —
// lint reports it as such rather than rejecting it"
// DOMAIN EDGE
func TestReq5_ExistsOverAnAlwaysPresentKeyIsVacuousNotRejected(t *testing.T) {
	src := declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`) + `
[[rule]]
id = "vacuous-exists"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
exists = true
[rule.write]
profile = "large"
`
	m := mustLoadSource(t, src)

	// Well-formed: it loads. Not rejected.
	reports := guard.Lint(m)
	if len(reports) == 0 {
		t.Fatal("Lint returned no group reports for a loadable model")
	}

	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeVacuousAtom &&
			f.Dimension == "profile" &&
			slices.Contains(f.RuleIDs, "vacuous-exists")
	}) {
		t.Errorf("lint emitted no vacuous-exists report for an `exists` atom "+
			"over an always-present key; findings=%v — it is reported as "+
			"vacuous rather than rejected", allFindings(reports))
	}
	if anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "profile"
	}) {
		t.Errorf("lint rejected the vacuous `exists` atom; it is well-formed")
	}
}

// REQ-6: "`Key` must resolve to a declared tag, `Operator` must be allowed
// by the operator/kind matrix above, and `Literal` must parse to the
// operator's literal shape."
// ADVERSARIAL
func TestReq6_AtomWellFormednessIsKeyOperatorAndLiteral(t *testing.T) {
	decls := `
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`
	cases := []struct {
		name  string
		guard string
		want  table.Category
	}{
		{
			name:  "key must resolve to a declared tag",
			guard: "[rule.guard.all.undeclared]\neq = \"x\"\n",
			want:  table.CatUnknownTag,
		},
		{
			name:  "operator must be allowed by the matrix",
			guard: "[rule.guard.all.profile]\ncontains = [\"small\"]\n",
			want:  table.CatMalformedPredicateAtom,
		},
		{
			name:  "literal must parse to the operator's literal shape",
			guard: "[rule.guard.all.profile]\nexists = \"yes\"\n",
			want:  table.CatMalformedPredicateAtom,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := loadCategory(t, decls, tc.guard); got != tc.want {
				t.Errorf("refused as %q; want %q", got, tc.want)
			}
		})
	}
}

// REQ-7: "**Source identity is not an atom field** — it is carried by the
// enclosing normalized row (`RuleID`, `SourceLocator`, per RDR 0002 and
// `internal/resolve/resolve.go::Row`), which is why this RDR's atom
// identity tuple is the row's identity joined to the atom's own four
// fields rather than an identity stored on the atom."
// BOUNDARY
func TestReq7_SourceIdentityIsNotAnAtomFieldButJoinsTheIdentityTuple(t *testing.T) {
	row := table.Row{
		RuleID:        "r1",
		SourceLocator: "flows/rdr.toml:10",
		Atoms: []table.Atom{
			{Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"}},
		},
	}

	id := guard.AtomIdentity(row, row.Atoms[0])
	for _, part := range []string{"r1", "flows/rdr.toml:10", "profile", "all", "eq", "large"} {
		if !strings.Contains(id.String(), part) {
			t.Errorf("AtomIdentity(...) = %q; missing %q — the tuple is the "+
				"row's identity JOINED to the atom's four fields", id, part)
		}
	}

	// The same atom under a different enclosing row is a DIFFERENT identity,
	// which is only possible if source identity comes from the row.
	other := row
	other.RuleID = "r2"
	if guard.AtomIdentity(other, other.Atoms[0]) == id {
		t.Error("two rows with different RuleIDs produced one atom identity; " +
			"source identity is carried by the enclosing row")
	}
}

// REQ-8: "a guard atom is identified by the total tuple `(RuleID,
// SourceLocator, key, block, operator, literal)`, never by its index
// within `all` or `unless`."
// ADVERSARIAL
func TestReq8_AtomIdentityIsTheSixFieldTupleNeverThePositionalIndex(t *testing.T) {
	first := table.Atom{Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"}}
	second := table.Atom{Key: "iter", Block: table.BlockAll, Operator: "lt", Literal: []string{"3"}}

	forward := table.Row{RuleID: "r", SourceLocator: "f:1", Atoms: []table.Atom{first, second}}
	reversed := table.Row{RuleID: "r", SourceLocator: "f:1", Atoms: []table.Atom{second, first}}

	// Same atom, different index within the block: one identity.
	if guard.AtomIdentity(forward, first) != guard.AtomIdentity(reversed, first) {
		t.Error("reordering the atoms changed an atom's identity; identity is " +
			"the six-field tuple, never the index within `all` or `unless`")
	}

	// The tuple is TOTAL: one row carrying two atoms over one key in one
	// block, differing only in literal, yields two identities.
	a := table.Atom{Key: "iter", Block: table.BlockAll, Operator: "gte", Literal: []string{"1"}}
	b := table.Atom{Key: "iter", Block: table.BlockAll, Operator: "lt", Literal: []string{"5"}}
	conjoined := table.Row{RuleID: "r", SourceLocator: "f:1", Atoms: []table.Atom{a, b}}
	if guard.AtomIdentity(conjoined, a) == guard.AtomIdentity(conjoined, b) {
		t.Error("two atoms over one key in one block collided on identity; " +
			"the six-field tuple must be total")
	}
}

// REQ-9: "Semantic equality is `(tag, operator, literal)`."
// HAPPY PATH
func TestReq9_SemanticEqualityIsTagOperatorLiteral(t *testing.T) {
	atom := table.Atom{Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"}}

	fromR1 := table.Row{RuleID: "r1", SourceLocator: "a:1", Atoms: []table.Atom{atom}}
	fromR2 := table.Row{RuleID: "r2", SourceLocator: "b:9", Atoms: []table.Atom{atom}}

	if guard.SemanticKey(atom) != guard.SemanticKey(fromR2.Atoms[0]) {
		t.Error("semantic equality distinguished two atoms with the same " +
			"(tag, operator, literal); it is exactly those three fields")
	}
	// It must NOT include source identity: the identity tuple does that.
	if guard.AtomIdentity(fromR1, atom) == guard.AtomIdentity(fromR2, atom) {
		t.Error("the identity tuple ignored source identity")
	}

	// It must not include the block either — block is an identity field, not
	// a semantic one; two atoms differing only in block denote the same
	// subset of the product.
	unlessCopy := atom
	unlessCopy.Block = table.BlockUnless
	if guard.SemanticKey(atom) != guard.SemanticKey(unlessCopy) {
		t.Error("semantic equality distinguished two atoms differing only in " +
			"block; semantic equality is (tag, operator, literal)")
	}
}

// REQ-10: "**Which equality each operation uses is fixed, not left to the
// implementer**: diagnostics, deduplication, and any \"same atom\" claim
// use the **identity tuple**; only domain computation — deciding what
// subset of the product an atom denotes — uses **semantic equality**"
// DOMAIN EDGE
func TestReq10_DiagnosticsUseIdentityAndDomainComputationUsesSemanticEquality(t *testing.T) {
	// Two rules carrying a byte-identical atom. Domain computation must
	// give them the SAME denotation (semantic equality), while the
	// diagnostic must still name both rule ids (identity tuple).
	m := mustLoadSource(t, twoRuleIdenticalGuardSource())

	dom := guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockAll, Operator: "eq", Literal: []string{"large"},
	})
	same := guard.Denotation(m, "profile", table.Atom{
		Key: "profile", Block: table.BlockUnless, Operator: "eq", Literal: []string{"large"},
	})
	if !dom.Equal(same) {
		t.Error("domain computation gave two semantically equal atoms " +
			"different denotations; domain computation uses semantic equality")
	}

	reports := guard.Lint(m)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeOverlap &&
			slices.Contains(f.RuleIDs, "dup-a") && slices.Contains(f.RuleIDs, "dup-b")
	}) {
		t.Errorf("the overlap diagnostic did not name both rule ids; "+
			"diagnostics use the identity tuple. findings=%v", allFindings(reports))
	}
}

// REQ-11: "Overlap detection is therefore *not* an equality test at all:
// it intersects the rows' accepted assignment sets, so two byte-identical
// guards in two different rule ids correctly produce an overlap finding
// naming both rule ids, rather than being deduplicated into one row."
// ADVERSARIAL
func TestReq11_ByteIdenticalGuardsInTwoRulesOverlapRatherThanDeduplicate(t *testing.T) {
	m := mustLoadSource(t, twoRuleIdenticalGuardSource())

	reports := guard.Lint(m)
	overlaps := findingsWithCode(reports, guard.CodeOverlap)
	if len(overlaps) == 0 {
		t.Fatalf("two byte-identical guards in two rule ids produced no "+
			"overlap finding; they were deduplicated. findings=%v", allFindings(reports))
	}
	for _, f := range overlaps {
		if !slices.Contains(f.RuleIDs, "dup-a") || !slices.Contains(f.RuleIDs, "dup-b") {
			t.Errorf("overlap finding names %v; want both dup-a and dup-b", f.RuleIDs)
		}
	}
}

// REQ-12: "equality compares a tag value to one typed literal; membership
// checks a scalar tag against a typed literal set; bounded integer
// comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks presence
// of an optional tag value; set containment checks declared set-valued
// tags against a typed element set."
// HAPPY PATH
func TestReq12_OperatorSemanticsAreAsStated(t *testing.T) {
	var ev guard.Evaluator

	cases := []struct {
		name     string
		operator string
		literal  string
		value    string
		want     resolve.GuardResult
	}{
		{"eq compares to one typed literal", "eq", "large", "large", resolve.GuardTrue},
		{"eq unequal", "eq", "large", "small", resolve.GuardFalse},
		{"in checks membership in a literal set", "in", `["mid","large"]`, "large", resolve.GuardTrue},
		{"in non-member", "in", `["mid","large"]`, "small", resolve.GuardFalse},
		{"lt below bound", "lt", "3", "2", resolve.GuardTrue},
		{"lt at bound", "lt", "3", "3", resolve.GuardFalse},
		{"lte at bound", "lte", "3", "3", resolve.GuardTrue},
		{"gt above bound", "gt", "3", "4", resolve.GuardTrue},
		{"gt at bound", "gt", "3", "3", resolve.GuardFalse},
		{"gte at bound", "gte", "3", "3", resolve.GuardTrue},
		{"contains every listed element", "contains", `["bug"]`, `["bug","chore"]`, resolve.GuardTrue},
		{"contains missing element", "contains", `["urgent"]`, `["bug","chore"]`, resolve.GuardFalse},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			atom := resolve.GuardAtom{
				Key: "subject", Operator: tc.operator, Literal: tc.literal, Block: resolve.BlockAll,
			}
			if got := ev.Evaluate(atom, tc.value); got != tc.want {
				t.Errorf("Evaluate(%+v, %q) = %v; want %v", atom, tc.value, got, tc.want)
			}
		})
	}
}

// A21: a provable value atom requires the tag's `single_valued` marker,
// which makes "which operators narrow the tag's single held value" a
// classification of the closed vocabulary, not a property of one operator.
// `SingleValueOperator` publishes it so consumers read it rather than copy
// it.
//
// The assertion ranges over the WHOLE vocabulary, so an operator added to
// `Operators()` and left unclassified fails here rather than silently
// joining the not-single-value arm. That arm is not a default: `exists`
// reads presence and `contains` reads set membership, and each is excluded
// for a stated reason, not by omission.
// BOUNDARY
func TestReq2_SingleValueOperatorsAreClassifiedOverTheWholeVocabulary(t *testing.T) {
	singleValue := []string{"eq", "in", "lt", "lte", "gt", "gte"}
	notSingleValue := []string{"exists", "contains"}

	// The two arms must PARTITION the vocabulary. Without this, adding an
	// operator to `Operators()` and to neither list below would leave the
	// loop asserting nothing about it.
	classified := slices.Sorted(slices.Values(
		append(slices.Clone(singleValue), notSingleValue...)))
	if got := slices.Sorted(slices.Values(guard.Operators())); !slices.Equal(got, classified) {
		t.Fatalf("Operators() = %v but this test classifies %v; every "+
			"operator in the closed vocabulary MUST be placed in exactly "+
			"one arm — an unclassified operator is the drift this "+
			"assertion exists to catch", got, classified)
	}

	for _, op := range singleValue {
		if !guard.SingleValueOperator(op) {
			t.Errorf("SingleValueOperator(%q) = false; %q narrows the tag's "+
				"single held value, so an atom spelling it projects only "+
				"over a key the model declares single-valued (A21)", op, op)
		}
	}
	for _, op := range notSingleValue {
		if guard.SingleValueOperator(op) {
			t.Errorf("SingleValueOperator(%q) = true; %q does not narrow "+
				"the tag's single held value, so requiring the "+
				"`single_valued` marker for it would name a remedy that "+
				"does not fix anything", op, op)
		}
	}

	// An operator outside the closed vocabulary narrows nothing, exactly as
	// it is accepted by no kind.
	for _, unknown := range []string{"neq", "matches", "", "EQ"} {
		if guard.SingleValueOperator(unknown) {
			t.Errorf("SingleValueOperator(%q) = true; the vocabulary is "+
				"closed and a non-operator classifies as nothing", unknown)
		}
	}
}
