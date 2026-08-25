package guard_test

// RDR 0003 — the tag declaration model this RDR owns: the five value
// kinds, finite domains, optionality, single-valuedness, element
// universes, and the domain/kind agreement rules.

import (
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/guard"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-13: "This RDR owns the **tag declaration model** — the typed
// alphabet every guard atom is written against. A tag declaration MUST
// carry a value kind, and — as its kind admits, per the agreement clause
// below — MAY carry a finite domain, an optionality marker, a
// single-valued marker, and (for set-valued kinds) an element universe."
// HAPPY PATH
func TestReq13_DeclarationCarriesAKindAndTheFourOptionalFields(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true

[tags.labels]
provenance = "owned"
kind = "set"
elements = ["bug", "chore"]
`))

	model := guard.DeclarationOf(m, "profile")
	if model.Kind != "enum" {
		t.Errorf("profile kind = %q; want enum — a declaration MUST carry a value kind", model.Kind)
	}
	if !slices.Equal(model.Domain, []string{"small", "large"}) {
		t.Errorf("profile domain = %v; want the declared finite domain", model.Domain)
	}
	if !model.SingleValued {
		t.Error("profile single-valued marker was not carried")
	}
	if model.Optional {
		t.Error("profile declares required = true; the optionality marker was not carried")
	}

	labels := guard.DeclarationOf(m, "labels")
	if !slices.Equal(labels.Elements, []string{"bug", "chore"}) {
		t.Errorf("labels element universe = %v; want the declared universe", labels.Elements)
	}

	// A declaration with no kind is not a declaration.
	if cat := loadDeclCategory(t, `
[tags.nokind]
provenance = "owned"
`); cat != table.CatMalformedTagDeclaration {
		t.Errorf("a declaration with no kind refused as %q; want %q — a "+
			"declaration MUST carry a value kind", cat, table.CatMalformedTagDeclaration)
	}
}

// REQ-14: "The value kinds are exactly five, spelled with these tokens
// wherever a kind is named: in a declaration, in the operator/kind matrix,
// and in a diagnostic — `enum`, `bool`, `int`, `set`, and `scalar`."
// BOUNDARY
func TestReq14_ValueKindsAreExactlyFiveTokens(t *testing.T) {
	want := []string{"bool", "enum", "int", "scalar", "set"}
	got := slices.Sorted(slices.Values(guard.Kinds()))
	if !slices.Equal(got, want) {
		t.Fatalf("Kinds() = %v; want exactly %v", got, want)
	}

	// The same five tokens spell a declaration.
	for _, kind := range guard.Kinds() {
		if !table.IsDeclaredKind(kind) {
			t.Errorf("kind %q is not spellable in a declaration; the tokens "+
				"are the same wherever a kind is named", kind)
		}
	}
	// And nothing else is.
	for _, foreign := range []string{"string", "integer", "boolean", "enumeration", "list", "Enum", ""} {
		if slices.Contains(guard.Kinds(), foreign) {
			t.Errorf("Kinds() admits %q; the vocabulary is exactly five tokens", foreign)
		}
	}
}

// REQ-15: "The `scalar` kind is the opaque scalar: a typed value compared
// only by equality and membership. It carries **no** finite domain and can
// never bear an exhaustiveness claim, so a guard dimension over a `scalar`
// takes the blocking inability-to-prove outcome, and a `scalar`
// declaration carrying a finite domain, an element universe, or a
// single-valued marker is a declaration error."
// DOMAIN EDGE
func TestReq15_ScalarIsOpaqueAndNeverBearsAnExhaustivenessClaim(t *testing.T) {
	// Compared only by equality and membership.
	for _, op := range guard.Operators() {
		want := op == "eq" || op == "in" || op == "exists"
		if guard.Accepts(op, "scalar") != want {
			t.Errorf("Accepts(%q, \"scalar\") = %v; a scalar is compared only "+
				"by equality and membership (plus presence)", op, !want)
		}
	}

	// No finite domain: no assignment count.
	if n, ok := guard.AssignmentCount(table.TagDecl{Kind: "scalar"}); ok {
		t.Errorf("AssignmentCount(scalar) = (%d, true); a scalar carries no "+
			"finite domain and therefore no assignment count", n)
	}

	// A guard dimension over a scalar takes the blocking outcome.
	m := mustLoadSource(t, declBlock(`
[tags.owner]
provenance = "observed"
kind = "scalar"
`)+`
[[rule]]
id = "scalar-guard"
source = "t:scalar"
[rule.match.recognized]
eq = "go"
[rule.guard.all.owner]
eq = "alice"
[rule.write]
`)
	reports := guard.Lint(m)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "owner"
	}) {
		t.Errorf("a guard dimension over a scalar drew no blocking "+
			"inability-to-prove finding; findings=%v", allFindings(reports))
	}
	if hasGreen(reports) {
		t.Error("a group with a scalar guard dimension certified green")
	}

	// Each of the three forbidden fields on a scalar is a declaration error.
	for _, decl := range []string{
		"[tags.s]\nprovenance = \"owned\"\nkind = \"scalar\"\ndomain = [\"a\"]\n",
		"[tags.s]\nprovenance = \"owned\"\nkind = \"scalar\"\nelements = [\"a\"]\n",
		"[tags.s]\nprovenance = \"owned\"\nkind = \"scalar\"\nsingle_valued = true\n",
	} {
		if cat := loadDeclCategory(t, decl); cat != table.CatMalformedTagDeclaration {
			t.Errorf("scalar declaration %q refused as %q; want %q",
				decl, cat, table.CatMalformedTagDeclaration)
		}
	}
}

// REQ-16: "RDR 0002 owns where a declaration is authored and how it is
// carried through normalization; this RDR owns what a declaration means.
// Neither document restates the other (JDR 0001 P6)."
// BOUNDARY
func TestReq16_ThisRDROwnsMeaningNotAuthoringOrCarriage(t *testing.T) {
	// The meaning surface takes an already-loaded model: this package
	// performs no TOML parsing and declares no wire keys of its own.
	m := mustLoadSource(t, declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`))

	d := guard.DeclarationOf(m, "profile")
	n, ok := guard.AssignmentCount(m.Tags["profile"])
	if !ok || n != 2 {
		t.Errorf("AssignmentCount over the loader's TagDecl = (%d, %v); want "+
			"(2, true) — meaning is computed over RDR 0002's carried declaration", n, ok)
	}
	if d.Kind != m.Tags["profile"].Kind {
		t.Errorf("guard.DeclarationOf disagreed with the loader's carried "+
			"kind (%q vs %q); this RDR reads the carried declaration, it does "+
			"not re-author it", d.Kind, m.Tags["profile"].Kind)
	}
}

// REQ-17: "A finite domain MUST be declarable for any kind an
// exhaustiveness claim can range over: an `enum` declares its value set, a
// `bool` is finite by construction, an `int` declares a `{min..max}` bound
// … and a `set` declares the element universe its members are drawn from."
// HAPPY PATH
func TestReq17_FiniteDomainIsDeclarableForEveryClaimableKind(t *testing.T) {
	min0, max3 := 0, 3
	cases := []struct {
		name string
		decl table.TagDecl
		want int
	}{
		{"enum declares its value set", table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c"}, SingleValued: true, Required: true}, 3},
		{"bool is finite by construction", table.TagDecl{Kind: "bool", SingleValued: true, Required: true}, 2},
		{"int declares a {min..max} bound", table.TagDecl{Kind: "int", Min: &min0, Max: &max3, SingleValued: true, Required: true}, 4},
		{"set declares its element universe", table.TagDecl{Kind: "set", Elements: []string{"x", "y"}, Required: true}, 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n, ok := guard.AssignmentCount(tc.decl)
			if !ok {
				t.Fatalf("AssignmentCount(%+v) reported no finite domain", tc.decl)
			}
			if n != tc.want {
				t.Errorf("AssignmentCount(%+v) = %d; want %d", tc.decl, n, tc.want)
			}
		})
	}
}

// REQ-18: "this RDR fixes their meaning — both endpoints are **inclusive**,
// so `{0..3}` has cardinality 4"
// BOUNDARY
func TestReq18_IntBoundEndpointsAreInclusive(t *testing.T) {
	min0, max3 := 0, 3
	decl := table.TagDecl{Kind: "int", Min: &min0, Max: &max3, SingleValued: true, Required: true}

	n, ok := guard.AssignmentCount(decl)
	if !ok {
		t.Fatal("AssignmentCount({0..3}) reported no finite domain")
	}
	if n != 4 {
		t.Errorf("AssignmentCount({0..3}) = %d; want 4 — both endpoints are inclusive", n)
	}

	// A single-point bound has cardinality 1, not 0.
	one := 1
	single := table.TagDecl{Kind: "int", Min: &one, Max: &one, SingleValued: true, Required: true}
	if n, ok := guard.AssignmentCount(single); !ok || n != 1 {
		t.Errorf("AssignmentCount({1..1}) = (%d, %v); want (1, true) — "+
			"inclusive endpoints", n, ok)
	}

	// The declared domain itself enumerates every inclusive value.
	if got := guard.IntDomain(decl); !slices.Equal(got, []int{0, 1, 2, 3}) {
		t.Errorf("IntDomain({0..3}) = %v; want [0 1 2 3]", got)
	}
}

// REQ-19: "A declaration carrying no finite domain is well-formed — the
// tag remains runtime-evaluable — but a guard dimension over it cannot
// carry an exhaustiveness claim, and lint MUST take the blocking
// inability-to-prove outcome for that dimension."
// DOMAIN EDGE
func TestReq19_NoFiniteDomainIsWellFormedButUnprovable(t *testing.T) {
	src := declBlock(`
[tags.iter]
provenance = "owned"
kind = "int"
required = true
`) + `
[[rule]]
id = "unbounded-int"
source = "t:iter"
[rule.match.recognized]
eq = "go"
[rule.guard.all.iter]
lt = 3
[rule.write]
`
	// Well-formed: it loads.
	m := mustLoadSource(t, src)

	// Runtime-evaluable: the evaluator still decides the atom.
	var ev guard.Evaluator
	atom := kernelAtom("iter", "lt", "3")
	if got := ev.Evaluate(atom, "2"); got != resolveTrue() {
		t.Errorf("Evaluate over an unbounded int = %v; want GuardTrue — the "+
			"tag remains runtime-evaluable", got)
	}

	// Lint takes the blocking inability-to-prove outcome for that dimension.
	reports := guard.Lint(m)
	if !anyFinding(reports, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "iter"
	}) {
		t.Errorf("an int declaring no bound drew no blocking finding naming "+
			"the dimension; findings=%v", allFindings(reports))
	}
	if hasGreen(reports) {
		t.Error("a group over an unbounded int certified green")
	}
}

// REQ-20: "A tag declaration MUST be able to state whether the key may be
// absent. A declaration carrying **no optionality marker declares the key
// optional** — the conservative default"
// INPUT EDGE
func TestReq20_OmittedOptionalityMarkerDeclaresTheKeyOptional(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.marked]
provenance = "owned"
kind = "bool"
required = true

[tags.unmarked]
provenance = "owned"
kind = "bool"
`))

	if guard.DeclarationOf(m, "marked").Optional {
		t.Error("required = true declared the key optional; it is always-present")
	}
	if !guard.DeclarationOf(m, "unmarked").Optional {
		t.Error("a declaration carrying no optionality marker was read as " +
			"always-present; the conservative default is OPTIONAL")
	}
}

// REQ-21: "A key declared always-present MUST NOT be absent from a
// conforming evaluation view; a key declared optional MAY be."
// DOMAIN EDGE
func TestReq21_ConformanceRequiresEveryAlwaysPresentKey(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.always]
provenance = "owned"
kind = "bool"
required = true

[tags.maybe]
provenance = "owned"
kind = "bool"
`))

	full := guard.View{"always": "true", "maybe": "false"}
	if err := guard.Conforms(m, full); err != nil {
		t.Errorf("a view carrying both keys does not conform: %v", err)
	}

	withoutOptional := guard.View{"always": "true"}
	if err := guard.Conforms(m, withoutOptional); err != nil {
		t.Errorf("a view missing only the OPTIONAL key does not conform: %v — "+
			"a key declared optional MAY be absent", err)
	}

	withoutRequired := guard.View{"maybe": "false"}
	if err := guard.Conforms(m, withoutRequired); err == nil {
		t.Error("a view missing an always-present key conformed; it MUST NOT be absent")
	}
}

// REQ-22: "Presence is a declared property of the tag, not an observation
// of one view: lint decides `exists` projection and the \"can refuse\"
// narrowing from this declaration, never from a runtime trace."
// ADVERSARIAL
func TestReq22_PresenceIsReadFromTheDeclarationNotAView(t *testing.T) {
	// Two models identical but for the optionality marker. Lint's verdict
	// must differ, and it must differ WITHOUT any view being supplied —
	// Lint takes a model, not a trace.
	optional := mustLoadSource(t, optionalKeyGroupSource(false))
	alwaysPresent := mustLoadSource(t, optionalKeyGroupSource(true))

	withheld := guard.Lint(optional)
	green := guard.Lint(alwaysPresent)

	if !anyFinding(withheld, func(f guard.Finding) bool {
		return f.Code == guard.CodeUnprovableCoverage && f.Dimension == "gate"
	}) {
		t.Errorf("the optional-key model was not withheld; findings=%v — "+
			"presence is read from the declaration", allFindings(withheld))
	}
	if !hasGreen(green) {
		t.Errorf("the always-present model did not certify green; reports=%s",
			renderReports(green))
	}
}

// REQ-23: "An evaluation view **conforms** to the declared model when
// every always-present key is present in it and every single-valued tag
// holds at most one of its declared domain values. Conformance is the
// premise every lint claim in this RDR is conditional on: a green
// exhaustiveness result asserts coverage over conforming views only."
// DOMAIN EDGE
func TestReq23_ConformanceIsBothConjunctsAndScopesEveryGreenClaim(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`))

	if err := guard.Conforms(m, guard.View{"profile": "large"}); err != nil {
		t.Errorf("a conforming view was rejected: %v", err)
	}
	// Conjunct one: every always-present key is present.
	if err := guard.Conforms(m, guard.View{}); err == nil {
		t.Error("a view missing an always-present key conformed")
	}
	// Conjunct two: a single-valued tag holds AT MOST ONE declared value.
	if err := guard.Conforms(m, guard.View{"profile": `["small","large"]`}); err == nil {
		t.Error("a view holding two values of a single-valued tag conformed")
	}

	// A green claim is scoped to conforming views: the report says so.
	reports := guard.Lint(mustLoadSource(t, completePartitionSource()))
	for _, r := range greenGroups(reports) {
		if !r.ConformingViewsOnly {
			t.Errorf("group %s certified green without scoping the claim to "+
				"conforming views", r.Context)
		}
	}
}

// REQ-24: "this RDR MUST NOT be read as promising anything about a view
// that violates the declarations."
// ADVERSARIAL
func TestReq24_NoClaimIsMadeAboutANonConformingView(t *testing.T) {
	m := mustLoadSource(t, completePartitionSource())

	reports := guard.Lint(m)
	if !hasGreen(reports) {
		t.Fatalf("the complete partition did not certify green; reports=%s",
			renderReports(reports))
	}

	// A view violating the always-present declaration is outside the claim:
	// Conforms rejects it, so the green result says nothing about it.
	bad := guard.View{}
	if err := guard.Conforms(m, bad); err == nil {
		t.Fatal("a view violating the declarations conformed")
	}
	for _, r := range greenGroups(reports) {
		if r.Covers(bad) {
			t.Errorf("group %s claims coverage over a non-conforming view; "+
				"this RDR promises nothing about it", r.Context)
		}
	}
}

// REQ-25: "A tag declaration MUST be able to state that the tag is
// **single-valued**: at most one of its declared domain values holds in
// any conforming evaluation view. The marker is meaningful only for a kind
// carrying a finite domain, and a `set`-valued kind MUST NOT carry it."
// BOUNDARY
func TestReq25_SingleValuedMarkerIsStatableAndBannedOnSet(t *testing.T) {
	m := mustLoadSource(t, declBlock(`
[tags.marked]
provenance = "owned"
kind = "enum"
domain = ["a", "b"]
single_valued = true
required = true
`))
	if !guard.DeclarationOf(m, "marked").SingleValued {
		t.Error("the single-valued marker was not carried")
	}
	// "at most one of its declared domain values holds": two values violate it.
	if err := guard.Conforms(m, guard.View{"marked": `["a","b"]`}); err == nil {
		t.Error("a view holding two values of a single-valued tag conformed")
	}

	// A set kind MUST NOT carry the marker.
	if cat := loadDeclCategory(t, `
[tags.labels]
provenance = "owned"
kind = "set"
elements = ["x", "y"]
single_valued = true
`); cat != table.CatMalformedTagDeclaration {
		t.Errorf("single_valued on a set refused as %q; want %q", cat, table.CatMalformedTagDeclaration)
	}

	// Meaningful only for a kind carrying a finite domain: a scalar has none.
	if cat := loadDeclCategory(t, `
[tags.opaque]
provenance = "owned"
kind = "scalar"
single_valued = true
`); cat != table.CatMalformedTagDeclaration {
		t.Errorf("single_valued on a scalar refused as %q; want %q", cat, table.CatMalformedTagDeclaration)
	}
}

// REQ-26: "a single-valued tag contributes **one dimension of `|domain|`
// assignments** (for enum `{a,b,c}`: three), because exactly one value
// holds per view. Absent the marker, a tag whose values could co-occur
// contributes one **independent boolean dimension per value**
// (`2^|domain|`, minus nothing — the model does not assume at least one
// holds)."
// BOUNDARY
func TestReq26_SingleValuedContributesDomainSizeUnmarkedContributesPowerset(t *testing.T) {
	marked := table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c"}, SingleValued: true, Required: true}
	if n, ok := guard.AssignmentCount(marked); !ok || n != 3 {
		t.Errorf("AssignmentCount(single-valued enum {a,b,c}) = (%d, %v); want (3, true)", n, ok)
	}

	unmarked := table.TagDecl{Kind: "enum", Domain: []string{"a", "b", "c"}, Required: true}
	if n, ok := guard.AssignmentCount(unmarked); !ok || n != 8 {
		t.Errorf("AssignmentCount(unmarked enum {a,b,c}) = (%d, %v); want "+
			"(8, true) — 2^|domain|, one independent boolean dimension per value", n, ok)
	}

	// "minus nothing — the model does not assume at least one holds": the
	// all-false assignment is in the space, so it is 2^3 and not 2^3-1.
	if n, _ := guard.AssignmentCount(unmarked); n == 7 {
		t.Error("AssignmentCount(unmarked enum {a,b,c}) = 7; the model does " +
			"not assume at least one value holds, so it is 2^|domain| minus nothing")
	}
}

// REQ-27: "RDR 0006's grouping-dependent lint findings read this field
// from the declaration and MUST NOT infer it from a tag's name, its value
// spelling, or a fixture (JDR 0001 §JD-13)."
// ADVERSARIAL
func TestReq27_SingleValuedIsNeverInferredFromNameSpellingOrFixture(t *testing.T) {
	// A tag NAMED and SPELLED exactly like a single-valued one, but with
	// the marker omitted, must not be read as single-valued.
	m := mustLoadSource(t, declBlock(`
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
required = true
`))
	if guard.DeclarationOf(m, "status").SingleValued {
		t.Error("a tag named `status` with a mutually-exclusive-looking domain " +
			"was inferred single-valued; the marker is read from the declaration")
	}
	if n, _ := guard.AssignmentCount(m.Tags["status"]); n != 4 {
		t.Errorf("AssignmentCount(unmarked status) = %d; want 4 (2^2) — the "+
			"marker was inferred from the tag's name or value spelling", n)
	}

	// And a single-value FIXTURE does not license the marker either.
	if err := guard.Conforms(m, guard.View{"status": `["Draft","Final"]`}); err != nil {
		t.Errorf("an unmarked tag holding two values did not conform: %v — "+
			"without the marker its values may co-occur", err)
	}
}

// REQ-28: "A declared domain MUST agree with its value kind. Each kind
// admits exactly these fields, and any other combination is a declaration
// error that MUST be rejected before normalization completes"
// ADVERSARIAL
func TestReq28_DomainMustAgreeWithKindAndDisagreementIsRejected(t *testing.T) {
	cases := []struct {
		name string
		decl string
	}{
		{"a {min..max} bound on an enum", `
[tags.t]
provenance = "owned"
kind = "enum"
domain = ["a"]
min = 0
max = 3
`},
		{"an element universe on an enum", `
[tags.t]
provenance = "owned"
kind = "enum"
domain = ["a"]
elements = ["x"]
`},
		{"a value set on an int", `
[tags.t]
provenance = "owned"
kind = "int"
domain = ["1", "2"]
`},
		{"an element universe on a bool", `
[tags.t]
provenance = "owned"
kind = "bool"
elements = ["x"]
`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if cat := loadDeclCategory(t, tc.decl); cat != table.CatMalformedTagDeclaration {
				t.Errorf("refused as %q; want %q — rejected before "+
					"normalization completes", cat, table.CatMalformedTagDeclaration)
			}
		})
	}

	// The agreement rule is this RDR's, so it must be decidable on this
	// RDR's own surface and not only observable through the loader: an
	// AGREEING declaration carries an assignment count, a disagreeing one
	// is not a declaration this model admits at all.
	agreeing := table.TagDecl{Kind: "enum", Domain: []string{"a", "b"}, SingleValued: true, Required: true}
	if _, ok := guard.AssignmentCount(agreeing); !ok {
		t.Errorf("an agreeing declaration %+v carries no assignment count; "+
			"agreement is what makes the domain readable", agreeing)
	}
	min0, max3 := 0, 3
	disagreeing := table.TagDecl{
		Kind: "enum", Domain: []string{"a"}, Min: &min0, Max: &max3, SingleValued: true, Required: true,
	}
	if n, ok := guard.AssignmentCount(disagreeing); ok {
		t.Errorf("a disagreeing declaration %+v reported assignment count %d; "+
			"an {min..max} bound on an enum is a declaration error, not a "+
			"second readable domain", disagreeing, n)
	}
}

// REQ-29: The kind/field agreement table: `enum` → value set, no element
// universe, single-valued allowed; `bool` → finite by construction, no
// element universe, single-valued allowed; `int` → `{min..max}`, no
// element universe, single-valued allowed; `set` → element universe
// required to claim exhaustiveness, "**no** — values co-occur by
// construction" for the marker; `scalar` → "**no**" domain, no element
// universe, "**no**" marker. So "an `{min..max}` bound on an `enum`, an
// element universe on any kind but `set`, or a single-valued marker on a
// `set` or a `scalar`, is each a declaration error."
// BOUNDARY
func TestReq29_KindFieldAgreementTableIsEnforcedPerCell(t *testing.T) {
	// Element universe on any kind but `set`.
	for _, kind := range []string{"enum", "bool", "int", "scalar"} {
		decl := "[tags.t]\nprovenance = \"owned\"\nkind = \"" + kind + "\"\nelements = [\"x\"]\n"
		if cat := loadDeclCategory(t, decl); cat != table.CatMalformedTagDeclaration {
			t.Errorf("element universe on %q refused as %q; want %q", kind, cat, table.CatMalformedTagDeclaration)
		}
	}
	// Single-valued marker on a `set` or a `scalar`.
	for _, kind := range []string{"set", "scalar"} {
		decl := "[tags.t]\nprovenance = \"owned\"\nkind = \"" + kind + "\"\nsingle_valued = true\n"
		if cat := loadDeclCategory(t, decl); cat != table.CatMalformedTagDeclaration {
			t.Errorf("single_valued on %q refused as %q; want %q", kind, cat, table.CatMalformedTagDeclaration)
		}
	}
	// Single-valued IS allowed on enum, bool, int.
	for _, kind := range []string{"enum", "bool", "int"} {
		decl := "[tags.t]\nprovenance = \"owned\"\nkind = \"" + kind + "\"\nsingle_valued = true\n"
		if _, err := table.Load([]byte(declBlock(decl)), "fixture.toml"); err != nil {
			t.Errorf("single_valued on %q refused: %v; the table allows it", kind, err)
		}
	}
	// `set` needs its element universe to claim exhaustiveness.
	if n, ok := guard.AssignmentCount(table.TagDecl{Kind: "set", Required: true}); ok {
		t.Errorf("AssignmentCount(set with no element universe) = (%d, true); "+
			"the universe is required to claim exhaustiveness", n)
	}
	// And WITH the universe it does, so the cell above is an agreement rule
	// rather than a blanket refusal.
	if n, ok := guard.AssignmentCount(table.TagDecl{
		Kind: "set", Elements: []string{"x", "y"}, Required: true,
	}); !ok || n != 4 {
		t.Errorf("AssignmentCount(set with a two-element universe) = (%d, %v); "+
			"want (4, true) — the universe is what makes the claim available", n, ok)
	}

	// Each ADMITTED cell of the table is readable, not merely unrejected: an
	// AssignmentCount that refuses everything passes every rejection check
	// above while enforcing no table at all.
	min0, max2 := 0, 2
	for name, d := range map[string]table.TagDecl{
		"enum with a value set":       {Kind: "enum", Domain: []string{"a", "b"}, SingleValued: true, Required: true},
		"bool finite by construction": {Kind: "bool", SingleValued: true, Required: true},
		"int with {min..max}":         {Kind: "int", Min: &min0, Max: &max2, SingleValued: true, Required: true},
	} {
		if _, ok := guard.AssignmentCount(d); !ok {
			t.Errorf("the admitted cell %q carries no assignment count", name)
		}
	}
}

// REQ-30: "A value appearing in a guard literal that lies outside its
// tag's declared domain MUST likewise be rejected — this is the predicate
// semantic kind **literal-outside-declared-domain**, distinct from a
// literal parse failure (the literal parses fine; it is simply not in the
// domain) — so an unsatisfiable atom is a load-time error rather than a
// silently-never-matching row."
// ADVERSARIAL
func TestReq30_LiteralOutsideDeclaredDomainIsRejectedAtLoad(t *testing.T) {
	decls := `
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small", "large"]
single_valued = true
required = true
`
	// The literal parses fine as an enum token; it is simply not in the
	// domain. It must still refuse.
	if cat := loadCategory(t, decls, "[rule.guard.all.profile]\neq = \"gigantic\"\n"); cat != table.CatMalformedPredicateAtom {
		t.Errorf("an out-of-domain enum literal refused as %q; want %q", cat, table.CatMalformedPredicateAtom)
	}

	// The two are DISTINCT predicate semantic kinds: this RDR names them.
	if guard.SemanticKindLiteralOutsideDomain == guard.SemanticKindLiteralParseFailure {
		t.Error("literal-outside-declared-domain and literal parse failure are " +
			"one kind; the RDR names them as distinct")
	}
	if !slices.Contains(guard.SemanticKinds(), guard.SemanticKindLiteralOutsideDomain) {
		t.Errorf("SemanticKinds() = %v; missing literal-outside-declared-domain",
			guard.SemanticKinds())
	}
}

// REQ-31: "a malformed *declaration* (domain disagreeing with its kind) is
// rejected by the declaration loader before normalization completes … a
// *literal-outside-domain* is rejected by guard parsing after the
// declarations have loaded and before rows are yielded."
// DOMAIN EDGE
func TestReq31_DeclarationErrorsAndLiteralErrorsAreRejectedAtDifferentStages(t *testing.T) {
	// A model whose declaration is malformed AND whose literal is out of
	// domain refuses on the DECLARATION — the earlier stage.
	src := declBlock(`
[tags.profile]
provenance = "owned"
kind = "enum"
domain = ["small"]
elements = ["x"]
required = true
`) + `
[[rule]]
id = "r"
[rule.match.recognized]
eq = "go"
[rule.guard.all.profile]
eq = "gigantic"
`
	_, err := table.Load([]byte(src), "fixture.toml")
	if err == nil {
		t.Fatal("a model with both defects loaded clean")
	}
	cat, _ := table.CategoryOf(err)
	if cat != table.CatMalformedTagDeclaration {
		t.Errorf("refused as %q; want %q — the declaration loader runs before "+
			"guard parsing yields rows", cat, table.CatMalformedTagDeclaration)
	}

	// The two stages are separately nameable on this RDR's own surface, so
	// a consumer can tell which one rejected: they map to different load
	// categories, and neither is the other's.
	declStage := guard.LoadCategoryFor(guard.SemanticKindDeclarationDisagreement)
	literalStage := guard.LoadCategoryFor(guard.SemanticKindLiteralOutsideDomain)
	if declStage != table.CatMalformedTagDeclaration {
		t.Errorf("the declaration stage maps to %q; want %q", declStage,
			table.CatMalformedTagDeclaration)
	}
	if literalStage != table.CatMalformedPredicateAtom {
		t.Errorf("the literal stage maps to %q; want %q — it is rejected by "+
			"guard parsing AFTER the declarations have loaded", literalStage,
			table.CatMalformedPredicateAtom)
	}
	if declStage == literalStage {
		t.Error("both stages map to one category; they are rejected at " +
			"different stages by different rules")
	}
}

// REQ-32: "Both are this RDR's rejection rules carried by RDR 0002's load
// categories (JDR 0001 §D7(iii)): the malformed declaration is a
// `malformed tag declaration`, the out-of-domain literal a `malformed
// predicate atom`. RDR 0006 mints nothing for either; RDR 0005 maps both
// onto the envelope under JDR 0001 §JD-8."
// BOUNDARY
func TestReq32_TheTwoRejectionsMapOntoRDR0002sTwoLoadCategories(t *testing.T) {
	if got := guard.LoadCategoryFor(guard.SemanticKindDeclarationDisagreement); got != table.CatMalformedTagDeclaration {
		t.Errorf("declaration disagreement maps to %q; want %q", got, table.CatMalformedTagDeclaration)
	}
	if got := guard.LoadCategoryFor(guard.SemanticKindLiteralOutsideDomain); got != table.CatMalformedPredicateAtom {
		t.Errorf("literal-outside-domain maps to %q; want %q", got, table.CatMalformedPredicateAtom)
	}

	// RDR 0006 mints nothing for either: neither is a lint finding code.
	for _, code := range guard.Codes() {
		if string(code) == string(table.CatMalformedTagDeclaration) ||
			string(code) == string(table.CatMalformedPredicateAtom) {
			t.Errorf("lint code %q duplicates a load category; RDR 0006 mints "+
				"nothing for either rejection", code)
		}
	}
}

// REQ-33: "a tag declaration carries five fields: value kind, finite
// domain, optionality, single-valuedness, and (for set kinds) an element
// universe. … any enumeration of this model in this document states all
// five."
// BOUNDARY
func TestReq33_TheDeclarationModelIsExactlyFiveFields(t *testing.T) {
	want := []string{"element universe", "finite domain", "optionality", "single-valuedness", "value kind"}
	got := slices.Sorted(slices.Values(guard.DeclarationFields()))
	if !slices.Equal(got, want) {
		t.Errorf("DeclarationFields() = %v; want exactly %v — five fields, "+
			"single-valuedness the fifth", got, want)
	}
}
