package table_test

// RDR 0024 `0024:C2` — PROOF: the load-pipeline refusals.
//
// The two cross-checks run as steps of the source-to-candidate-rows
// pipeline, ahead of `normalizeRules`, and read the SOURCE rules. The
// oracle is the stable category plus the stamped source line.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// declVerdictCovering declares `verdict` over exactly the three values the
// 0010 decision-table rules author, so a covering model loads clean.
const declVerdictCovering = `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma"]`

// REQ-23 / `0024:S3`: "PROOF (opt-in, whole-model). With ZERO `[emit.*]`
// declarations the load pipeline is byte-for-byte today's: no new refusal
// is reachable."
// HAPPY PATH
func TestReq23_0024_ZeroDeclarationsMakeNoNewRefusalReachable(t *testing.T) {
	// The 0010 corpus authors emit keys nothing declares and values no
	// domain admits. With zero declarations none of that is reachable.
	for name, src := range map[string]string{
		"dtComplete":         dtComplete,
		"dtCompleteEmitting": dtCompleteEmitting,
		"dtEscapeVariant":    dtEscapeVariant,
	} {
		t.Run(name, func(t *testing.T) {
			m := loadSource(t, src, name+".toml")
			if m.EmitDecls == nil {
				t.Error("EmitDecls is nil after a successful zero-declaration " +
					"load; it is empty, never nil")
			}
			if len(m.EmitDecls) != 0 {
				t.Errorf("EmitDecls = %v; a zero-declaration model declares "+
					"none", m.EmitDecls)
			}
		})
	}
}

// REQ-24: "With ONE OR MORE declarations present, the
// source-to-candidate-rows pipeline (`0002:C24`'s \"load\") MUST refuse,
// before yielding rows. \"Before yielding rows\" fixes the position: rows
// are minted in `internal/table/load.go`'s `normalizeRules` step, so both
// checks run as a step ahead of it, beside `loadTags`, and therefore read
// the SOURCE rules (`sourceRule.Emit`) rather than normalized rows."
// ADVERSARIAL
func TestReq24_0024_TheChecksRunBeforeRowsAreYielded(t *testing.T) {
	// The undeclared key lives on a rule that ALSO fails to normalize
	// cleanly is not usable here (fail-fast, unspecified order). Instead
	// the observable is: the refusal returns no model, so no row survived.
	src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
		"verdict = \"alpha\"\n", "verdict = \"alpha\"\nnxet = \"x\"\n", 1)

	m, err := table.Load([]byte(src), "emit-before-rows.toml")
	if err == nil {
		t.Fatal("an undeclared emit key loaded clean; C2 refuses it")
	}
	// The refusal must be C2's, not an incidental one: without this the
	// `m != nil` leg below is satisfied by any refusal at all.
	cat, ok := table.CategoryOf(err)
	if !ok || cat != table.CatUnknownEmitKey {
		t.Fatalf("category = %q (categorized=%v); want %q — the check runs "+
			"as a step of the load pipeline, ahead of `normalizeRules`",
			cat, ok, table.CatUnknownEmitKey)
	}
	if m != nil {
		t.Errorf("the refusal yielded a model with %d rows; the check runs "+
			"BEFORE rows are minted", len(m.Rows))
	}
}

// REQ-25: "**Two steps, in this order, at this position.** The work is
// `loadEmitDecls` (C1's grammar checks over `[emit]`, building C3's
// carrier) then `checkRuleEmit` (this contract's two cross-checks over
// every `sourceRule.Emit`). `loadEmitDecls` MUST precede `checkRuleEmit` …
// and both are inserted immediately after `loadTags` in `load.go::run`'s
// step slice."
// REQ-55: "**The reuse is safe-by-omission, so C2's value refusal DEPENDS
// on C1's arms having already fired.**"
// ADVERSARIAL
func TestReq25_0024_TheGrammarStepPrecedesTheCrossCheckStep(t *testing.T) {
	// A model carrying BOTH a malformed declaration and a rule-side defect
	// against that same declaration. `loadEmitDecls` runs first, so the
	// declaration defect is the one reported: were the order reversed, the
	// cross-check would read a half-built carrier.
	body := `[emit.verdict]
kind = "enum"
domain = ["alpha", "", "gamma"]`
	src := strings.Replace(dtWithEmitDecl(body),
		"verdict = \"beta\"\n", "verdict = \"not-a-member\"\n", 1)
	if !strings.Contains(src, "not-a-member") {
		t.Fatal("the value substitution did not apply")
	}

	f := refuseSource(t, src, "emit-order.toml")
	if f.Category != table.CatMalformedEmitDeclaration {
		t.Errorf("category = %q; want %q — `loadEmitDecls` MUST precede "+
			"`checkRuleEmit`, so a malformed declaration is reported "+
			"rather than a value checked against it",
			f.Category, table.CatMalformedEmitDeclaration)
	}
}

// REQ-27: "**The zero-declaration opt-in gate is evaluated INSIDE each
// step**, not by the caller: the step slice is uniform and has no room for
// a conditional call, so both steps run unconditionally and return `nil`
// early when no key is declared."
// BOUNDARY
func TestReq27_0024_TheOptInGateIsEvaluatedInsideEachStep(t *testing.T) {
	// The observable of an in-step gate on a zero-declaration model is that
	// the model loads clean AND its carrier is initialised — a step that
	// never ran could not have set it. `Model.EmitDecls` is non-nil after
	// EVERY successful load (REQ-47).
	m := loadSource(t, dtComplete, "emit-optin-gate.toml")
	if m.EmitDecls == nil {
		t.Error("`Model.EmitDecls` is nil on a zero-declaration model; both " +
			"steps run unconditionally and return nil early")
	}
}

// REQ-28 / `0024:S2`: "`unknown_emit_key` — any `[rule.emit]` key of any
// rule, ordinary or escape, not declared under `[emit]`. Strictness is
// whole-model, not per-key"
// ADVERSARIAL
func TestReq28_0024_AnUndeclaredEmitKeyRefusesOnBothRuleClasses(t *testing.T) {
	t.Run("ordinary rule", func(t *testing.T) {
		src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
			"verdict = \"alpha\"\n", "verdict = \"alpha\"\nnxet = \"x\"\n", 1)
		f := refuseSource(t, src, "emit-unknown-ordinary.toml")
		if f.Category != table.CatUnknownEmitKey {
			t.Errorf("category = %q; want %q", f.Category, table.CatUnknownEmitKey)
		}
	})

	t.Run("escape rule", func(t *testing.T) {
		// dtEscapeVariant's `otherwise` escape row emits `verdict =
		// "fallback"`; declare `fallback` into the domain so the ONLY
		// defect is the undeclared key on the escape row.
		src := strings.Replace(dtEscapeVariant, dtHeader,
			dtHeader+emitDeclTable(`[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma", "fallback"]`), 1)
		if src == dtEscapeVariant {
			t.Fatal("the declaration substitution did not apply")
		}
		src = strings.Replace(src,
			"verdict = \"fallback\"\n", "verdict = \"fallback\"\nnxet = \"x\"\n", 1)
		f := refuseSource(t, src, "emit-unknown-escape.toml")
		if f.Category != table.CatUnknownEmitKey {
			t.Errorf("category = %q; want %q — strictness is whole-model, "+
				"and an ESCAPE row's emit block is in scope",
				f.Category, table.CatUnknownEmitKey)
		}
	})
}

// REQ-29 / `0024:S2` / REQ-96: "`emit_value_out_of_domain` — an authored
// value that is not a member of its key's declared enum domain, not a
// `bool` token, or not an `int` literal, per C1's kinds (`scalar` values
// are never refused)."
// DOMAIN EDGE
func TestReq29_0024_AnOutOfDomainValueRefusesPerKind(t *testing.T) {
	cases := []struct {
		name    string
		decl    string
		value   string
		refuses bool
	}{
		{
			name:    "enum member outside the declared domain",
			decl:    `kind = "enum"` + "\n" + `domain = ["route-a", "route-b"]`,
			value:   "route-c",
			refuses: true,
		},
		{
			name:    "enum member inside the declared domain",
			decl:    `kind = "enum"` + "\n" + `domain = ["route-a", "route-b"]`,
			value:   "route-a",
			refuses: false,
		},
		{
			name:    "not a bool token",
			decl:    `kind = "bool"`,
			value:   "yes",
			refuses: true,
		},
		{
			name:    "a bool token",
			decl:    `kind = "bool"`,
			value:   "false",
			refuses: false,
		},
		{
			name:    "not an int literal",
			decl:    `kind = "int"`,
			value:   "twelve",
			refuses: true,
		},
		{
			name:    "an int literal",
			decl:    `kind = "int"`,
			value:   "12",
			refuses: false,
		},
		{
			// "`scalar` values are never refused" — the declared-but-
			// unvalidated escape hatch (`0024:F2`, REQ-74).
			name:    "scalar admits any value",
			decl:    `kind = "scalar"`,
			value:   "anything at all",
			refuses: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := declVerdictCovering + "\n\n[emit.extra]\n" + tc.decl
			src := strings.Replace(dtWithEmitDecl(body),
				"verdict = \"alpha\"\n",
				"verdict = \"alpha\"\nextra = \""+tc.value+"\"\n", 1)

			if !tc.refuses {
				loadSource(t, src, "emit-value-ok.toml")
				return
			}
			f := refuseSource(t, src, "emit-value-bad.toml")
			if f.Category != table.CatEmitValueOutOfDomain {
				t.Errorf("category = %q; want %q",
					f.Category, table.CatEmitValueOutOfDomain)
			}
		})
	}
}

// REQ-30: "**All three categories MUST carry a source line**, stamped
// through `internal/table/category.go::atLine` the way `loadTags` stamps a
// tag refusal."
// REQ-31: "a declaration defect keys on the top-level `[emit.<key>]`
// header, which `tagHeaderLine`'s technique reaches unchanged"
// REQ-32: "a rule-side defect (`unknown_emit_key`,
// `emit_value_out_of_domain`) keys on the offending RULE ID, **not on its
// `[rule.emit]` header**."
// ASSUMPTION-4: asserted as a non-zero, non-1 line pointing at the
// offending block, not as an exact fixture integer.
// DOMAIN EDGE
func TestReq30_0024_EveryEmitRefusalCarriesTheOffendingSourceLine(t *testing.T) {
	t.Run("malformed_emit_declaration keys on `[emit.<key>]`", func(t *testing.T) {
		src := dtWithEmitDecl(`[emit.verdict]
kind = "gauge"`)
		f := refuseSource(t, src, "emit-line-decl.toml")
		assertLineAt(t, f, src, "[emit.verdict]")
	})

	t.Run("unknown_emit_key keys on the offending RULE ID", func(t *testing.T) {
		src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
			"verdict = \"gamma\"\n", "verdict = \"gamma\"\nnxet = \"x\"\n", 1)
		f := refuseSource(t, src, "emit-line-unknown.toml")
		// `cell-yp` is the rule authoring `verdict = "gamma"`.
		assertLineAt(t, f, src, "id = \"cell-yp\"")
	})

	t.Run("emit_value_out_of_domain keys on the offending RULE ID", func(t *testing.T) {
		src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
			"verdict = \"gamma\"\n", "verdict = \"not-a-member\"\n", 1)
		f := refuseSource(t, src, "emit-line-value.toml")
		assertLineAt(t, f, src, "id = \"cell-yp\"")
	})
}

// REQ-32 tail: "The rule-side locator anchors on the rule's `id =
// \"<ruleID>\"` line and scans forward to the emit block … The anchor must
// match the id's VALUE, since `[model]` also carries an `id` key."
// ADVERSARIAL
func TestReq32_0024_TheRuleAnchorMatchesTheIdValueNotTheModelIdKey(t *testing.T) {
	// `[model]` carries `id = "dt"` on an EARLIER line than every rule. An
	// anchor matching the bare `id` KEY rather than the VALUE would stamp
	// the model header's line instead of the offending rule's.
	src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
		"verdict = \"gamma\"\n", "verdict = \"not-a-member\"\n", 1)

	f := refuseSource(t, src, "emit-anchor.toml")
	// The refusal must be the rule-side one, and it must be STAMPED:
	// an unstamped (Line == 0) refusal would satisfy the inequality below
	// while grounding nothing.
	if f.Category != table.CatEmitValueOutOfDomain {
		t.Fatalf("category = %q; want %q", f.Category,
			table.CatEmitValueOutOfDomain)
	}
	if f.Line <= 0 {
		t.Fatalf("the refusal carries no source line (%d); the rule-side "+
			"locator anchors on the rule's `id` VALUE", f.Line)
	}
	modelIDLine := lineOf(t, src, "id = \"dt\"")
	if f.Line == modelIDLine {
		t.Errorf("the refusal is stamped at line %d, which is `[model]`'s "+
			"own `id` key; the anchor must match the RULE id's VALUE",
			f.Line)
	}
}

// REQ-38 / `0024:S2` / REQ-97: "**Reporting is FAIL-FAST, one refusal per
// run — inherited from the load tier, not decided here.**"
// REQ-39: "This RDR fixes no precedence among the three emit categories and
// the existing ones, and takes none: order stays unspecified."
// REQ-40: "`0005:C1`'s \"one findings[] entry per category hit\" is
// satisfied vacuously, because no load path yields more than one hit to
// map — this RDR introduces no accumulation."
// ADVERSARIAL
func TestReq38_0024_AnEmitRefusalIsFailFastAndCarriesNoErrorList(t *testing.T) {
	// An emit defect AND a coexisting structural one. Exactly ONE refusal
	// comes back, and the test does NOT assert which.
	src := strings.Replace(dtWithEmitDecl(declVerdictCovering),
		"verdict = \"gamma\"\n", "verdict = \"not-a-member\"\n", 1)
	src = strings.Replace(src, "[tags.b]", "[tags.b]\nnot_a_field = 1", 1)
	if !strings.Contains(src, "not_a_field") {
		t.Fatal("the structural-defect substitution did not apply")
	}

	_, err := table.Load([]byte(src), "emit-failfast.toml")
	if err == nil {
		t.Fatal("a model carrying two defects loaded clean")
	}
	if _, ok := table.CategoryOf(err); !ok {
		t.Fatalf("the refusal carries no category: %v", err)
	}
	// The load tier returns ONE categorized error, never a list.
	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		t.Errorf("the refusal unwraps to %d errors; load is fail-fast and "+
			"this RDR introduces no accumulation", len(multi.Unwrap()))
	}

	// Both defects, alone, refuse — so the coexisting run above really did
	// suppress one of two reachable refusals.
	single := strings.Replace(dtWithEmitDecl(declVerdictCovering),
		"verdict = \"gamma\"\n", "verdict = \"not-a-member\"\n", 1)
	if f := refuseSource(t, single, "emit-failfast-single.toml"); f.Category !=
		table.CatEmitValueOutOfDomain {
		t.Errorf("the emit defect alone refuses %q; want %q",
			f.Category, table.CatEmitValueOutOfDomain)
	}
}

// REQ-41: "Proving the AUTHORED form is proving the executed form:
// `0010:C3` fixes that normalization carries the emit block through
// `expand` unmutated … Nothing is checked at resolve time, and an emit
// value's evaluation semantics — uninterpreted, byte-compared — are
// unchanged (`0010:C3`)."
// HAPPY PATH
func TestReq41_0024_NormalizationCarriesTheEmitBlockUnmutated(t *testing.T) {
	m := loadEmitDecls0024(t, declVerdictCovering, "emit-unmutated.toml")

	rows := rowsByRuleID(m, "cell-xp")
	if len(rows) != 1 {
		t.Fatalf("rule `cell-xp` expanded to %d rows; want one", len(rows))
	}
	if len(rows[0].Emit) != 1 || rows[0].Emit[0].Key != "verdict" ||
		rows[0].Emit[0].Value != "alpha" {
		t.Errorf("row emit = %v; want the authored block unmutated "+
			"(verdict=alpha)", rows[0].Emit)
	}
}

// REQ-42: "A declared key no rule emits, and a declared member no rule
// authors, are NOT findings of any tier (authoring headroom; the advisory
// tier is closed, `0006:C17`)"
// DOMAIN EDGE
func TestReq42_0024_UnusedDeclarationsAreNotFindingsOfAnyTier(t *testing.T) {
	body := `[emit.verdict]
kind = "enum"
domain = ["alpha", "beta", "gamma", "delta-no-rule-authors-this"]

[emit.never-emitted]
kind = "scalar"`

	m, advisories, err := table.LoadWithAdvisories(
		[]byte(dtWithEmitDecl(body)), "emit-headroom.toml")
	if err != nil {
		t.Fatalf("authoring headroom must load clean; refused: %v", err)
	}
	if m == nil {
		t.Fatal("a clean load yielded a nil model")
	}
	if len(advisories) != 0 {
		t.Errorf("unused declarations produced %d advisories; they are NOT "+
			"findings of any tier — the advisory tier is closed",
			len(advisories))
	}
}

// --- helpers -------------------------------------------------------------

// assertLineAt asserts the failure carries a grounded source line pointing
// at the fixture line containing anchor.
//
// ASSUMPTION-4: the property is "non-zero, non-1, pointing at the offending
// block", not an exact integer fixture — a literal line number would break
// on any whitespace edit to the fixture.
func assertLineAt(t *testing.T, f *table.Failure, src, anchor string) {
	t.Helper()

	if f.Line == 0 {
		t.Fatalf("the %s refusal carries no source line; all three "+
			"categories MUST carry one", f.Category)
	}
	if f.Line == 1 {
		t.Errorf("the %s refusal is stamped at line 1, which `MVV` step 2 "+
			"names as the failure mode ('the offending block's SOURCE "+
			"LINE, not `:1`')", f.Category)
	}
	if want := lineOf(t, src, anchor); f.Line != want {
		t.Errorf("the %s refusal is stamped at line %d; the offending block "+
			"%q is at line %d", f.Category, f.Line, anchor, want)
	}
}

// lineOf returns the 1-based line number of the single fixture line whose
// trimmed text equals anchor.
func lineOf(t *testing.T, src, anchor string) int {
	t.Helper()

	line, matches := 0, 0
	for i, raw := range strings.Split(src, "\n") {
		if strings.TrimSpace(raw) == anchor {
			matches++
			line = i + 1
		}
	}
	if matches != 1 {
		t.Fatalf("the fixture carries %d lines equal to %q; the anchor must "+
			"be unique for the assertion to mean anything", matches, anchor)
	}
	return line
}
