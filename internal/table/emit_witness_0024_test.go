package table_test

// RDR 0024 `0024:S6` / `0024:S9` — category registration, checked-in
// witnesses, and the rendering negative control.
//
// `0024:D-naming` fixes the three slugs as literals and REQ-71 fixes that
// the wire slug and the registered category are ONE decision: scenario 6
// asserts the registration rather than trusting the append.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// `0024:S6`'s three checked-in witnesses are registered in
// `internal/table/dump_test.go::TestReq106_EveryNamedCategoryExistsAndIs
// Witnessed`'s map, which ASSUMPTION-10 reads as THE convention for a new
// load category. A second, parallel map here would have kept them outside
// the authoritative floor — including outside its exclusivity leg and its
// completeness leg — so the registration lives there and the assertions
// below name the three category constants directly.

// REQ-70 / `0024:D-naming`: "The new load categories are
// `malformed_emit_declaration`, `unknown_emit_key`, and
// `emit_value_out_of_domain`" — the three slugs are fixed literals.
// DOMAIN EDGE
func TestReq70_0024_TheThreeCategorySlugsAreFixedLiterals(t *testing.T) {
	for constant, want := range map[table.Category]string{
		table.CatMalformedEmitDeclaration: "malformed_emit_declaration",
		table.CatUnknownEmitKey:           "unknown_emit_key",
		table.CatEmitValueOutOfDomain:     "emit_value_out_of_domain",
	} {
		if string(constant) != want {
			t.Errorf("category constant renders %q; the slug is the fixed "+
				"literal %q", string(constant), want)
		}
	}
}

// REQ-71 / `0024:D-naming`: "Each slug also gets its `Cat*` constant in
// `internal/table/category.go` and an entry in `Categories()`, which is
// what a consumer branches on — the wire slug and the registered category
// are one decision, not two, and scenario 6 asserts the registration rather
// than trusting the append."
// REQ-104 / `0024:S6`: "**Expected**: `malformed_emit_declaration`,
// `unknown_emit_key` and `emit_value_out_of_domain` each appear in
// `table.Categories()`"
// HAPPY PATH
func TestReq104_0024_EachEmitCategoryIsRegisteredInCategories(t *testing.T) {
	declared := table.Categories()
	for _, cat := range emitCategories0024 {
		if !slices.Contains(declared, cat) {
			t.Errorf("category %q is not in table.Categories(); the wire "+
				"slug and the registered category are ONE decision", cat)
		}
	}
}

// emitCategories0024 names the three constants `0024:D-naming` fixes, the
// same three `TestReq70_0024` pins to their literal slugs.
var emitCategories0024 = []table.Category{
	table.CatMalformedEmitDeclaration,
	table.CatUnknownEmitKey,
	table.CatEmitValueOutOfDomain,
}

// REQ-104 tail: "each carries a checked-in `testdata/neg/*.toml` witness
// asserted by category, per the convention `internal/table/dump_test.go`'s
// REQ-119 witness map holds for every existing load category."
// REQ-105 / `0024:S6`: "**The convention has an exclusivity leg (REQ-131)
// this scenario must also carry**: each witness trips its own category *and
// no other*. Assert it per witness. … each witness must be single-defect by
// construction"
//
// Both legs are carried by
// `TestReq106_EveryNamedCategoryExistsAndIsWitnessed`, which the three
// categories are now registered in: its per-category subtest IS the
// exclusivity assertion (`loadCategory(t, rel) != cat` fails on a witness
// tripping a second category first), and its completeness leg is what
// would notice a fourth emit category shipping without one. Restating
// either here against a private map would only reintroduce the parallel
// registration this kata removed.

// REQ-108 / `0024:S9`: "**Expected**: an enum-declared emit value renders
// byte-identically to its undeclared rendering —
// `internal/table/dump.go::renderEmit` keeps bypassing `renderValue`. This
// is the negative control on C1's \"no value is parsed, canonicalized, or
// converted anywhere downstream\""
// ADVERSARIAL
func TestReq108_0024_ADeclaredEmitValueRendersByteIdenticallyToAnUndeclaredOne(t *testing.T) {
	undeclared := table.Dump(loadSource(t, dtComplete, "emit-undeclared.toml"))
	declared := table.Dump(
		loadEmitDecls0024(t, declVerdictCovering, "emit-declared.toml"))

	if undeclared != declared {
		t.Errorf("the dump of a DECLARED model differs from the undeclared "+
			"one; `renderEmit` keeps bypassing `renderValue`, so declaring "+
			"a domain changes no rendered byte.\n--- undeclared ---\n%s\n"+
			"--- declared ---\n%s", undeclared, declared)
	}
}

// REQ-108 second leg: the negative control is only meaningful if a value
// that `renderValue` WOULD quote or canonicalize stays raw.
// DOMAIN EDGE
func TestReq108_0024_AnIntDeclaredEmitValueIsNotCanonicalizedInTheDump(t *testing.T) {
	body := declVerdictCovering + `

[emit.count]
kind = "int"`
	src := replaceOnce(t, dtWithEmitDecl(body),
		"verdict = \"alpha\"\n", "verdict = \"alpha\"\ncount = \"007\"\n")

	rendered := table.Dump(loadSource(t, src, "emit-int-render.toml"))
	if !strings.Contains(rendered, "count=007") {
		t.Errorf("the dump does not render `count=007`; an `int`-declared "+
			"emit value is never parsed, canonicalized, or converted "+
			"anywhere downstream.\n%s", rendered)
	}
}

// REQ-79 / `PH1`: "`internal/table/emit_shape_0010_test.go::TestReq21_
// NoHandWrittenTypeCheckDiscriminatesTheEmitRefusal` … keeps passing … Do
// not \"fix\" TestReq21 when adding the kind arm; if it goes red, the kind
// check has been placed above the decoder, which C1 does not license."
// ADVERSARIAL
func TestReq79_0024_TheKindArmStaysBelowTheDecoder(t *testing.T) {
	// TestReq21's property, restated against a DECLARED model: adding the
	// kind arm must not move a wrong-typed emit value off the decoder's
	// category. The tripwire only means something once a declaration is
	// present, which is the state TestReq21 itself never reaches.
	body := declVerdictCovering + `

[emit.count]
kind = "int"`
	// The control: the same declaration with a QUOTED int value must load
	// clean. Without it every assertion below is satisfied by a build that
	// refuses the `[emit]` table itself, which proves nothing about where
	// the kind arm sits.
	loadSource(t, replaceOnce(t, dtWithEmitDecl(body),
		"verdict = \"alpha\"\n", "verdict = \"alpha\"\ncount = \"42\"\n"),
		"emit-typed-control.toml")

	emitTyped := refuseSource(t, replaceOnce(t, dtWithEmitDecl(body),
		"verdict = \"alpha\"\n", "verdict = \"alpha\"\ncount = 42\n"),
		"emit-typed-declared.toml")

	otherTyped := refuseSource(t, replaceOnce(t, dtWithEmitDecl(body),
		`id = "cell-xp"`, `id = 42`), "other-typed-declared.toml")

	if emitTyped.Category != otherTyped.Category {
		t.Errorf("a wrong-typed emit value under a DECLARED key refuses %q "+
			"while another wrong-typed scalar refuses %q; the kind check "+
			"has been placed ABOVE the decoder, which C1 does not license",
			emitTyped.Category, otherTyped.Category)
	}
	if emitTyped.Category != table.CatMalformedTOML {
		t.Errorf("category = %q; want %q",
			emitTyped.Category, table.CatMalformedTOML)
	}
}

// REQ-78 / `PH1`: "`TestReq27`'s ASSERTIONS need no edit, but its
// failure-message prose restates \"undeclared\" in the same absolute sense
// the three comments do — narrow those strings to \"not a tag key\" in this
// same change"
// REQ-20 tail: `TestReq27_AnEmitKeyIsNotATagKey` keeps passing unchanged.
// ADVERSARIAL
func TestReq78_0024_AnEmitKeyIsStillNotATagKeyUnderADeclaration(t *testing.T) {
	// The assertion TestReq27 makes, restated against a DECLARED model:
	// declaring an emit key must NOT enter it into the tag table, and must
	// not make it usable in match or guard.
	m := loadEmitDecls0024(t, declVerdictCovering, "emit-not-a-tag.toml")
	if _, declared := m.Tags["verdict"]; declared {
		t.Error("the declared emit key `verdict` appears in the model's TAG " +
			"table; an emit key is not a tag key, declared or not")
	}

	src := replaceOnce(t, dtWithEmitDecl(declVerdictCovering),
		"[rule.guard.all.a]\neq = \"x\"\n",
		"[rule.match.verdict]\neq = \"alpha\"\n[rule.guard.all.a]\neq = \"x\"\n")
	f := refuseSource(t, src, "emit-declared-in-match.toml")
	if f.Category != table.CatUnknownTag {
		t.Errorf("category = %q; want %q — a DECLARED emit key remains "+
			"barred from match", f.Category, table.CatUnknownTag)
	}
}

// REQ-75 / `0024:F3`: "a model carrying `[emit]` refuses on an older binary
// as `unknown_schema_field` (strict decode) — loud, never a silent ignore
// of the declaration."
// ADVERSARIAL — no code may weaken strict decode for the whole document.
func TestReq75_0024_StrictDecodeStillRefusesAnUnknownTopLevelTable(t *testing.T) {
	// The forward-compat property is `decodeStrict`'s: an unknown TOP-LEVEL
	// table refuses loudly. `[emit]` is exactly that on a pre-change build,
	// so the invariant is asserted on a stand-in table this build does not
	// know either.
	f := refuseSource(t, dtComplete+"\n[not_a_declared_table]\nx = \"y\"\n",
		"emit-forward-compat.toml")
	if f.Category != table.CatUnknownSchemaField {
		t.Errorf("category = %q; want %q — an undeclared top-level table "+
			"refuses LOUDLY, never as a silent ignore",
			f.Category, table.CatUnknownSchemaField)
	}
}
