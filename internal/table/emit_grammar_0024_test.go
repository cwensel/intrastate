package table_test

// RDR 0024 `0024:C1` — GRAMMAR: the `[emit]` declaration table.
//
// The oracle everywhere is the stable load CATEGORY (`0002:C24`), never
// message text. Fixtures are complete authored TOML handed to the real
// `table.Load`.

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-1: "GRAMMAR. A model MAY declare its emit vocabulary in a top-level
// `[emit]` TOML table, one sub-table per emit key: `[emit.<key>]`."
// HAPPY PATH
func TestReq1_0024_AModelMayDeclareItsEmitVocabulary(t *testing.T) {
	m := loadEmitDecls0024(t, declVerdictEnum, "emit-decl-happy.toml")

	if _, ok := m.EmitDecls["verdict"]; !ok {
		t.Fatalf("a top-level `[emit.verdict]` sub-table declared no emit "+
			"key; EmitDecls = %v", m.EmitDecls)
	}
}

// REQ-2: "**The opt-in trigger is the COUNT of declared keys, never the
// presence of the `[emit]` table** (a bare table with zero sub-tables is a
// zero-declaration model; the tail of this contract gives the reasoning)."
// BOUNDARY
func TestReq2_0024_TheOptInTriggerIsTheCountOfDeclaredKeys(t *testing.T) {
	// A bare `[emit]` table with zero sub-tables. If presence were the
	// trigger, the model's undeclared `verdict` emit keys would refuse
	// `unknown_emit_key`; the count is zero, so nothing fires.
	m := loadEmitDecls0024(t, "[emit]", "emit-bare-table.toml")

	if len(m.EmitDecls) != 0 {
		t.Errorf("a bare `[emit]` table declared %d keys; a table with zero "+
			"sub-tables is a ZERO-declaration model", len(m.EmitDecls))
	}
}

// REQ-3: "`kind` (required): one of `enum | bool | int | scalar` — RDR
// 0003's token spellings reused verbatim" … "`set` is excluded because an
// emit value is one authored string, `0010:C3`"
// DOMAIN EDGE
func TestReq3_0024_KindIsOneOfFourTokensAndSetIsExcluded(t *testing.T) {
	admitted := []string{"enum", "bool", "int", "scalar"}
	for _, kind := range admitted {
		t.Run("admits "+kind, func(t *testing.T) {
			body := "[emit.free]\nkind = \"" + kind + "\""
			if kind == "enum" {
				body += "\ndomain = [\"x\"]"
			}
			// `free` is a declared key no rule emits, which C2 licenses
			// (REQ-42), so the only thing under test is the kind token.
			m := loadEmitDecls0024(t, body+"\n"+declVerdictEnum,
				"emit-kind-"+kind+".toml")
			if got := m.EmitDecls["free"].Kind; got != kind {
				t.Errorf("Kind = %q; want %q", got, kind)
			}
		})
	}

	t.Run("`set` is not an admitted emit kind", func(t *testing.T) {
		// SINGLE DEFECT: the kind token alone. Carrying `elements` too — a
		// tag-declaration key an emit declaration does not admit (REQ-20) —
		// would trip the DECODER's `unknown_schema_field` first and prove
		// nothing about `set`.
		f := refuseEmitDecls0024(t,
			"[emit.verdict]\nkind = \"set\"",
			"emit-kind-set.toml")
		if f.Category != table.CatMalformedEmitDeclaration {
			t.Errorf("category = %q; want %q — `set` is excluded because an "+
				"emit value is one authored string",
				f.Category, table.CatMalformedEmitDeclaration)
		}
	})
}

// REQ-4: "a domain, `enum` only, authored in either of two TOML spellings
// of the SAME `domain` key" — `domain = [ ... ]` — "a flat member array, no
// dispositions; or" `[emit.<key>.domain]` — "a sub-table whose keys are
// MODEL-AUTHORED disposition tokens (e.g. `route`, `stop`, `terminal` —
// intrastate fixes no vocabulary) and whose values are member arrays. The
// key's domain is the union; each member carries the one disposition it is
// listed under."
// HAPPY PATH
func TestReq4_0024_DomainTakesTwoSpellingsOfTheSameKey(t *testing.T) {
	t.Run("flat member array, no dispositions", func(t *testing.T) {
		m := loadEmitDecls0024(t, declVerdictFlat, "emit-domain-flat.toml")
		d := m.EmitDecls["verdict"]
		want := []string{"alpha", "beta", "gamma"}
		if !slices.Equal(d.Domain, want) {
			t.Errorf("Domain = %v; want %v", d.Domain, want)
		}
		// REQ-46 fixes the carrier form, not merely the count:
		// `Dispositions` is "`nil` when the domain carries none". An empty
		// non-nil map is a different value, and the non-enum leg
		// (`TestReq6_0024`) already asserts `nil`, so the flat spelling is
		// held to the same canonical form.
		if d.Dispositions != nil {
			t.Errorf("Dispositions = %v; want nil — a flat array carries "+
				"no dispositions", d.Dispositions)
		}
	})

	t.Run("disposition sub-table: domain is the union", func(t *testing.T) {
		m := loadEmitDecls0024(t, declVerdictEnum, "emit-domain-table.toml")
		d := m.EmitDecls["verdict"]
		want := []string{"alpha", "beta", "gamma"}
		if !slices.Equal(d.Domain, want) {
			t.Errorf("Domain = %v; want the union %v", d.Domain, want)
		}
		// Each member carries the ONE disposition it is listed under.
		for member, want := range map[string]string{
			"alpha": "route", "beta": "route", "gamma": "stop",
		} {
			if got := d.Dispositions[member]; got != want {
				t.Errorf("Dispositions[%q] = %q; want %q", member, got, want)
			}
		}
	})

	t.Run("intrastate fixes no disposition vocabulary", func(t *testing.T) {
		// Tokens the record never names must be admitted verbatim.
		body := `[emit.verdict]
kind = "enum"
[emit.verdict.domain]
"weird token" = ["alpha", "beta"]
zzz = ["gamma"]`
		m := loadEmitDecls0024(t, body, "emit-domain-vocab.toml")
		if got := m.EmitDecls["verdict"].Dispositions["gamma"]; got != "zzz" {
			t.Errorf("Dispositions[\"gamma\"] = %q; want %q — intrastate "+
				"fixes no disposition vocabulary", got, "zzz")
		}
	})
}

// REQ-5: "**These are two spellings of one key, not two coexisting fields,
// so \"both present\" is not a refusal arm this contract owns** — authoring
// both is a duplicate-key error the TOML decoder raises before any check
// here runs. An implementer must not add a hand-written arm for it"
// ADVERSARIAL
func TestReq5_0024_BothDomainSpellingsIsTheDecodersDuplicateKeyError(t *testing.T) {
	body := `[emit.verdict]
kind = "enum"
domain = ["alpha"]
[emit.verdict.domain]
route = ["beta"]`

	f := refuseEmitDecls0024(t, body, "emit-domain-both.toml")
	if f.Category != table.CatMalformedTOML {
		t.Errorf("category = %q; want %q — authoring both spellings is a "+
			"duplicate-key error the DECODER raises; no hand-written arm "+
			"may re-report it", f.Category, table.CatMalformedTOML)
	}
}

// REQ-6: "`bool` fixes the implicit domain `true | false`; `int` constrains
// the value to a base-10 integer literal; `scalar` is the
// declared-but-unvalidated escape hatch — the key is admitted, the value
// unconstrained. None of the three takes a `domain`, and dispositions
// attach only to declared enum members."
// DOMAIN EDGE
func TestReq6_0024_NonEnumKindsTakeNoDomainAndNoDispositions(t *testing.T) {
	m := loadEmitDecls0024(t, `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]

[emit.flag]
kind = "bool"

[emit.count]
kind = "int"

[emit.free]
kind = "scalar"`, "emit-nonenum-kinds.toml")

	for _, key := range []string{"flag", "count", "free"} {
		d := m.EmitDecls[key]
		if d.Domain != nil {
			t.Errorf("EmitDecls[%q].Domain = %v; none of bool/int/scalar "+
				"takes a domain", key, d.Domain)
		}
		if d.Dispositions != nil {
			t.Errorf("EmitDecls[%q].Dispositions = %v; dispositions attach "+
				"only to declared ENUM members", key, d.Dispositions)
		}
	}
}

// REQ-7: "Every kind check is a LEXICAL check on the authored string: no
// value is parsed into a typed representation, canonicalized, or converted
// anywhere downstream — evaluation and the payload carry the authored bytes
// (`0010:C3`)."
// DOMAIN EDGE
func TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive(t *testing.T) {
	// An `int`-declared key whose authored literal is non-canonical. If the
	// check parsed and re-rendered, `007` would become `7` on the row.
	src := strings.Replace(dtWithEmitDecl(`[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]

[emit.code]
kind = "int"`), "verdict = \"alpha\"\n", "verdict = \"alpha\"\ncode = \"007\"\n", 1)

	m := loadSource(t, src, "emit-lexical.toml")
	row := rowByID(t, m, findRowIdentity(t, m, "cell-xp"))
	for _, e := range row.Emit {
		if e.Key != "code" {
			continue
		}
		if e.Value != "007" {
			t.Errorf("row emit `code` = %q; want %q — no value is parsed "+
				"into a typed representation, canonicalized, or converted",
				e.Value, "007")
		}
		return
	}
	t.Fatalf("the selected row carries no `code` emit value: %v", row.Emit)
}

// REQ-8: "The authored value is always a TOML **string**: `sourceRule.Emit`
// is `map[string]string` (`0010:C3`), so under `kind = \"int\"` an author
// writes `count = \"42\"`, and a bare `count = 42` refuses as
// `malformed_toml` from the decoder — not as `emit_value_out_of_domain` —
// exactly as it does today."
// ADVERSARIAL
func TestReq8_0024_ANonStringEmitValueRefusesAsMalformedTOML(t *testing.T) {
	decl := `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]

[emit.count]
kind = "int"`

	t.Run("the quoted int literal is admitted", func(t *testing.T) {
		src := strings.Replace(dtWithEmitDecl(decl),
			"verdict = \"alpha\"\n", "verdict = \"alpha\"\ncount = \"42\"\n", 1)
		loadSource(t, src, "emit-int-quoted.toml")
	})

	t.Run("the bare int literal refuses as malformed_toml", func(t *testing.T) {
		src := strings.Replace(dtWithEmitDecl(decl),
			"verdict = \"alpha\"\n", "verdict = \"alpha\"\ncount = 42\n", 1)
		f := refuseSource(t, src, "emit-int-bare.toml")
		if f.Category != table.CatMalformedTOML {
			t.Errorf("category = %q; want %q — a bare `count = 42` refuses "+
				"from the DECODER, never as %q",
				f.Category, table.CatMalformedTOML,
				table.CatEmitValueOutOfDomain)
		}
	})
}

// REQ-9: "Reusing `ConformValue` (A4) settles non-canonical literals in the
// permissive direction: `03`, `+5` and `-0` are admitted for `int`, because
// `strconv.Atoi` admits them."
// BOUNDARY
func TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted(t *testing.T) {
	decl := `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]

[emit.count]
kind = "int"`

	for _, literal := range []string{"03", "+5", "-0"} {
		t.Run(literal, func(t *testing.T) {
			src := strings.Replace(dtWithEmitDecl(decl),
				"verdict = \"alpha\"\n",
				"verdict = \"alpha\"\ncount = \""+literal+"\"\n", 1)
			// `strconv.Atoi` admits all three, so reuse admits them too.
			loadSource(t, src, "emit-int-"+literal+".toml")
		})
	}
}

// REQ-10: "Because `domain` is two-shaped, its decoded field is the ONE
// place in the source schema that cannot be concretely typed: it is
// declared **`any`** on the source-schema struct, mirroring
// `internal/table/source.go`'s `sourceModel.Metadata`"
// DOMAIN EDGE
func TestReq10_0024_BothDomainShapesDecodeThroughOneAnyTypedField(t *testing.T) {
	// The observable of an `any`-typed field is that BOTH authorings decode
	// without a decoder type error — a concretely typed field would refuse
	// one of the two as `malformed_toml`.
	for name, body := range map[string]string{
		"flat array":        declVerdictFlat,
		"disposition table": declVerdictEnum,
	} {
		t.Run(name, func(t *testing.T) {
			m := loadEmitDecls0024(t, body, "emit-domain-shape.toml")
			if len(m.EmitDecls["verdict"].Domain) != 3 {
				t.Errorf("Domain = %v; want three members from the %s "+
					"spelling", m.EmitDecls["verdict"].Domain, name)
			}
		})
	}
}

// REQ-11: "**The carve-out is that field type, NOT a decoder option** —
// `pelletier/go-toml/v2` has no per-subtree strictness switch, and
// `DisallowUnknownFields` stays set for the whole document"
// ADVERSARIAL
func TestReq11_0024_StrictDecodeStaysSetForTheWholeDocument(t *testing.T) {
	// An unknown field OUTSIDE the emit subtree must still refuse, proving
	// no global strictness toggle was introduced to admit the two-shaped
	// domain.
	src := strings.Replace(dtWithEmitDecl(declVerdictEnum),
		"[model]\nid = \"dt\"", "[model]\nid = \"dt\"\nnot_a_field = \"x\"", 1)
	if !strings.Contains(src, "not_a_field") {
		t.Fatal("the unknown-field substitution did not apply")
	}

	// The control: the SAME document minus the unknown field must LOAD.
	// Without it this assertion is satisfied by a build that refuses
	// `[emit]` itself, which proves nothing about the carve-out.
	loadEmitDecls0024(t, declVerdictEnum, "emit-strict-control.toml")

	f := refuseSource(t, src, "emit-strict-elsewhere.toml")
	if f.Category != table.CatUnknownSchemaField {
		t.Errorf("category = %q; want %q — DisallowUnknownFields stays set "+
			"for the WHOLE document", f.Category, table.CatUnknownSchemaField)
	}
}

// REQ-12: "**That declaration-level refusal is the decoder's, not this
// grammar's**: it arrives as `unknown_schema_field` (A1's `decodeStrict`
// arm), an existing `0002:C24` category and NOT one of the three this RDR
// registers. A test covering a `domaim` typo asserts the decoder's slug,
// and an implementer must not add a hand-written arm to re-report it as
// `malformed_emit_declaration`."
// ADVERSARIAL
func TestReq12_0024_ADeclarationKeyTypoIsTheDecodersUnknownSchemaField(t *testing.T) {
	// The control: correcting the typo must LOAD. Without it the assertion
	// below is satisfied by a build that refuses `[emit]` wholesale.
	loadEmitDecls0024(t, declVerdictFlat, "emit-domain-control.toml")

	f := refuseEmitDecls0024(t, `[emit.verdict]
kind = "enum"
domaim = ["alpha", "beta", "gamma"]`, "emit-domaim-typo.toml")

	if f.Category != table.CatUnknownSchemaField {
		t.Errorf("category = %q; want %q — the declaration-level refusal is "+
			"the DECODER's, never %q",
			f.Category, table.CatUnknownSchemaField,
			table.CatMalformedEmitDeclaration)
	}
}

// REQ-20: "Emit declarations are NOT tag declarations: no provenance, no
// accessor reference, no `min/max/elements/single_valued/required`, and an
// emit key remains barred from match, guard, write, and accessor use
// (`0010:C3`)."
// ADVERSARIAL
func TestReq20_0024_AnEmitDeclarationAdmitsNoTagDeclarationKey(t *testing.T) {
	// The control: the same declaration WITHOUT the tag-declaration key
	// must load. Without it every arm below is satisfied by a build that
	// refuses `[emit]` itself.
	loadEmitDecls0024(t, declVerdictFlat, "emit-tagkey-control.toml")

	for _, key := range []string{
		"provenance = \"observed\"",
		"min = 0",
		"max = 9",
		"elements = [\"alpha\"]",
		"single_valued = true",
		"required = true",
	} {
		t.Run(strings.SplitN(key, " ", 2)[0], func(t *testing.T) {
			body := "[emit.verdict]\nkind = \"enum\"\n" + key +
				"\ndomain = [\"alpha\", \"beta\", \"gamma\"]"
			f := refuseEmitDecls0024(t, body, "emit-tagkey.toml")
			// The key is not in the emit schema at all, so strict decode
			// refuses it — an emit declaration is NOT a tag declaration.
			if f.Category != table.CatUnknownSchemaField {
				t.Errorf("category = %q; want %q — an emit declaration "+
					"carries no tag-declaration key",
					f.Category, table.CatUnknownSchemaField)
			}
		})
	}
}

// REQ-21: "A **bare `[emit]` table with zero sub-tables is a
// zero-declaration model**, identical in every observable to omitting the
// table: it takes C2's opt-in leg, no refusal becomes reachable, and the
// payload carries `dispositions: {}`."
// BOUNDARY
func TestReq21_0024_ABareEmitTableIsIdenticalToOmittingIt(t *testing.T) {
	bare := loadEmitDecls0024(t, "[emit]", "emit-bare.toml")
	omitted := loadSource(t, dtComplete, "emit-omitted.toml")

	if len(bare.EmitDecls) != 0 || len(omitted.EmitDecls) != 0 {
		t.Errorf("bare declared %d keys, omitted declared %d; both are "+
			"zero-declaration models",
			len(bare.EmitDecls), len(omitted.EmitDecls))
	}
	// No refusal is reachable: the rules emit `verdict`, undeclared under
	// both authorings, and neither refuses.
	if len(bare.Rows) != len(omitted.Rows) {
		t.Errorf("bare yielded %d rows, omitted %d; a bare `[emit]` table "+
			"is identical in every observable",
			len(bare.Rows), len(omitted.Rows))
	}
}

// REQ-22: "A declaration is **author-owned and unversioned**: a domain may
// be widened or narrowed by editing the model, and intrastate holds no
// history to check the edit against — every check in this contract reads
// the model file against itself at one point in time."
// DOMAIN EDGE
func TestReq22_0024_ADeclarationIsAuthorOwnedAndUnversioned(t *testing.T) {
	// The wide baseline declares two members beyond the three the 0010
	// rules author. A member no rule authors is not a finding of any tier,
	// so the widened model loads clean.
	wide := loadEmitDecls0024(t, `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma", "delta", "epsilon"]`, "emit-wide.toml")

	// Narrow the domain back to exactly the values the rules author.
	// Nothing consults a history, so the narrowing edit loads clean too.
	// It cannot narrow BELOW the authored values without tripping
	// `checkRuleEmit`'s out-of-domain refusal, which is C2's contract and
	// not this one.
	narrowed := loadEmitDecls0024(t, declVerdictFlat, "emit-narrowed.toml")

	if len(wide.EmitDecls) != 1 || len(narrowed.EmitDecls) != 1 {
		t.Fatalf("both models declare one key; got %d and %d",
			len(wide.EmitDecls), len(narrowed.EmitDecls))
	}

	// The two legs must be observably DIFFERENT models, or the narrowing
	// never happened and the leg would pass against an implementation that
	// refused the edit by consulting a history.
	if got, want := len(wide.EmitDecls["verdict"].Domain), 5; got != want {
		t.Errorf("wide domain carries %d members; want %d", got, want)
	}
	if got, want := len(narrowed.EmitDecls["verdict"].Domain), 3; got != want {
		t.Errorf("narrowed domain carries %d members; want %d", got, want)
	}
}

// findRowIdentity returns the identity string of the single normalized row
// expanded from ruleID. The 0010 decision-table rules expand one-to-one.
func findRowIdentity(t *testing.T, m *table.Model, ruleID string) string {
	t.Helper()

	rows := rowsByRuleID(m, ruleID)
	if len(rows) != 1 {
		t.Fatalf("rule %q expanded to %d rows; want exactly one", ruleID, len(rows))
	}
	return rows[0].Identity()
}
