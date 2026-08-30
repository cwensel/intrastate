package table_test

// RDR 0024 `0024:S1` — the table-driven malformed-declaration sweep.
//
// This file is A7's Pending CLOSURE GATE. `PRE` states "Phase 1 does not
// close until it is green", and `0024:S1` fixes the sweep's shape: one
// defect per fixture, each loaded on its own, each refusing
// `malformed_emit_declaration` at load with exactly one blocking finding,
// before any row is yielded.
//
// ASSUMPTION-1 reads C1's refusal list as TEN distinct arms. The four the
// record calls "the arms that exist only because strictness cannot reach
// them" (REQ-18) are the ones A7's closure leg covers.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// emitArm is one malformed-declaration fixture: a single defect and the
// arm it exercises.
type emitArm struct {
	name string
	// req is the REQ whose clause this arm discharges.
	req string
	// body is the `[emit]` declaration text, authored verbatim.
	body string
}

// malformedEmitArms enumerates C1's refusal arms, one single-defect
// fixture each, in the order `0024:C1` and `0024:S1` give them.
//
// `0024:S1`: "One defect per fixture — the pipeline is fail-fast, so a
// fixture carrying two defects proves nothing about the second."
var malformedEmitArms = []emitArm{
	// REQ-13: "Refused as `malformed_emit_declaration`: an unknown `kind`
	// token"
	{
		name: "unknown kind token",
		req:  "REQ-13",
		body: `[emit.verdict]
kind = "gauge"`,
	},

	// REQ-14 / REQ-19 / REQ-95: "an `enum` carrying no usable domain" —
	// THREE fixtures against ONE expectation.
	{
		name: "no usable domain: no `domain` key at all",
		req:  "REQ-14",
		body: `[emit.verdict]
kind = "enum"`,
	},
	{
		name: "no usable domain: an empty `[emit.<key>.domain]` sub-table",
		req:  "REQ-14",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]`,
	},
	{
		name: "no usable domain: an empty flat `domain = []`",
		req:  "REQ-14",
		body: `[emit.verdict]
kind = "enum"
domain = []`,
	},

	// REQ-15: "an empty-string member, or a duplicate member (duplicates
	// across disposition lists included — one member, one disposition)"
	{
		name: "an empty-string member (flat)",
		req:  "REQ-15",
		body: `[emit.verdict]
kind = "enum"
domain = ["alpha", "", "gamma"]`,
	},
	{
		name: "an empty-string member (disposition list)",
		req:  "REQ-15",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", ""]
stop = ["gamma"]`,
	},
	{
		name: "a duplicate member within one flat domain",
		req:  "REQ-15",
		body: `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "alpha"]`,
	},
	{
		name: "a duplicate member within one disposition list",
		req:  "REQ-15",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "alpha"]
stop = ["gamma"]`,
	},
	{
		name: "a duplicate member ACROSS disposition lists",
		req:  "REQ-15",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["alpha", "gamma"]`,
	},

	// REQ-16 / ASSUMPTION-8: "a `domain` on a non-enum kind" — the PRESENCE
	// of the key under bool/int/scalar, in EITHER spelling.
	{
		name: "a flat `domain` on kind = bool",
		req:  "REQ-16",
		body: `[emit.flag]
kind = "bool"
domain = ["true", "false"]`,
	},
	{
		name: "a flat `domain` on kind = int",
		req:  "REQ-16",
		body: `[emit.count]
kind = "int"
domain = ["1", "2"]`,
	},
	{
		name: "a flat `domain` on kind = scalar",
		req:  "REQ-16",
		body: `[emit.free]
kind = "scalar"
domain = ["anything"]`,
	},
	{
		name: "a disposition-table `domain` on kind = bool",
		req:  "REQ-16",
		body: `[emit.flag]
kind = "bool"
[emit.flag.domain]
route = ["true"]`,
	},

	// REQ-17: "a disposition token that is the empty string"
	{
		name: "an empty-string disposition token",
		req:  "REQ-17",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
"" = ["alpha", "beta", "gamma"]`,
	},

	// REQ-18: "the arms that exist only because strictness cannot reach
	// them — a `domain` that is neither a flat array of strings nor a table
	// of disposition keys, a disposition whose value is not an array of
	// strings, any nesting below the disposition level, and a non-string
	// member." These four are A7's Pending closure leg.
	{
		name: "a `domain` that is neither a flat array nor a table (string)",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
domain = "alpha"`,
	},
	{
		name: "a `domain` that is neither a flat array nor a table (int)",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
domain = 42`,
	},
	{
		name: "a disposition whose value is not an array of strings",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = "alpha"`,
	},
	{
		name: "nesting below the disposition level",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain.route]
inner = ["alpha"]`,
	},
	{
		name: "a non-string member (flat)",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
domain = ["alpha", 42, "gamma"]`,
	},
	{
		name: "a non-string member (disposition list)",
		req:  "REQ-18",
		body: `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", 42]
stop = ["gamma"]`,
	},
}

// REQ-94 / `0024:S1`: "**Expected**: each case, loaded on its own, refuses
// `malformed_emit_declaration` at load with exactly one blocking finding,
// before any row is yielded. One defect per fixture — the pipeline is
// fail-fast, so a fixture carrying two defects proves nothing about the
// second."
// REQ-13, REQ-14, REQ-15, REQ-16, REQ-17, REQ-18: C1's refusal arms.
// ADVERSARIAL
func TestReq94_0024_EveryMalformedDeclarationArmRefuses(t *testing.T) {
	for _, arm := range malformedEmitArms {
		t.Run(arm.req+"/"+arm.name, func(t *testing.T) {
			model, err := table.Load(
				[]byte(dtWithEmitDecl(arm.body)), "emit-arm.toml")

			if err == nil {
				t.Fatalf("%s (%s) loaded CLEAN; C1 refuses it as %q\n"+
					"--- declaration ---\n%s",
					arm.name, arm.req, table.CatMalformedEmitDeclaration,
					arm.body)
			}
			// "before any row is yielded" — the loader returns no model at
			// all on a refusal, so nothing downstream ever sees a row.
			if model != nil {
				t.Errorf("%s refused AND yielded a model with %d rows; the "+
					"check runs BEFORE any row is yielded",
					arm.name, len(model.Rows))
			}
			cat, ok := table.CategoryOf(err)
			if !ok {
				t.Fatalf("%s refused with a non-categorized error: %v",
					arm.name, err)
			}
			if cat != table.CatMalformedEmitDeclaration {
				t.Errorf("%s refused %q; want %q",
					arm.name, cat, table.CatMalformedEmitDeclaration)
			}
		})
	}
}

// REQ-95 / `0024:S1`: "The no-usable-domain arm takes **three fixtures
// against one expectation** — no `domain` key, an empty
// `[emit.<key>.domain]` sub-table, and `domain = []` … Assert the shared
// category; do NOT assert distinct messages for the first two, which would
// pin a discrimination the decoder cannot make."
// REQ-19: "**\"No usable domain\" is ONE arm, not two, because the decoder
// cannot tell the two authorings apart.**"
// BOUNDARY
func TestReq95_0024_NoUsableDomainIsOneArmWithThreeFixtures(t *testing.T) {
	fixtures := map[string]string{
		"no `domain` key at all": `[emit.verdict]
kind = "enum"`,
		"an empty `[emit.<key>.domain]` sub-table": `[emit.verdict]
kind = "enum"
[emit.verdict.domain]`,
		"an empty flat `domain = []`": `[emit.verdict]
kind = "enum"
domain = []`,
	}

	for name, body := range fixtures {
		t.Run(name, func(t *testing.T) {
			f := refuseEmitDecls0024(t, body, "emit-no-domain.toml")
			if f.Category != table.CatMalformedEmitDeclaration {
				t.Errorf("category = %q; want %q — all three authorings "+
					"refuse under ONE arm",
					f.Category, table.CatMalformedEmitDeclaration)
			}
		})
	}

	// The two indistinguishable authorings must not be discriminated. Both
	// decode to the same absent domain, so the refusal DETAIL they carry
	// must be identical — a differing detail is a discrimination the
	// decoder cannot support.
	absentKey := refuseEmitDecls0024(t, fixtures["no `domain` key at all"],
		"emit-absent-key.toml")
	emptyTable := refuseEmitDecls0024(t,
		fixtures["an empty `[emit.<key>.domain]` sub-table"],
		"emit-empty-table.toml")
	if absentKey.Detail != emptyTable.Detail {
		t.Errorf("the absent-`domain` refusal detail is %q while the empty "+
			"sub-table's is %q; both decode to the same absent domain and "+
			"pinning a difference would fabricate a discrimination the "+
			"decoder cannot make", absentKey.Detail, emptyTable.Detail)
	}
}
