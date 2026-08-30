package cli

// RDR 0024 — the surface obligations: which envelope code carries an emit
// refusal, which tiers stay closed, and which named edits Phase 3 licenses.

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
	"github.com/cwensel/intrastate/internal/table"
)

// undeclaredKeyDT0024 is the declared decision table with one rule emitting
// an UNDECLARED key — the `unknown_emit_key` witness on the CLI surface.
func undeclaredKeyDT0024(t *testing.T) string {
	t.Helper()

	return strings.Replace(declaredDT0024(t),
		"verdict = \"alpha\"\ncode = \"1\"",
		"verdict = \"alpha\"\nnxet = \"typo\"", 1)
}

// outOfDomainDT0024 is the declared decision table with one rule emitting a
// value outside the declared domain.
func outOfDomainDT0024(t *testing.T) string {
	t.Helper()

	src := strings.Replace(declaredDT0024(t),
		"verdict = \"alpha\"\ncode = \"1\"", "verdict = \"not-a-member\"", 1)
	if strings.Contains(src, "code = \"1\"") {
		t.Fatal("the out-of-domain substitution did not apply")
	}
	return src
}

// REQ-34: "These two categories and C1's `malformed_emit_declaration` join
// `0002:C24`'s data-level set … and therefore map to `flow-model-invalid`
// under `0005:C1`, and to a blocking finding under `intrastate lint` via
// its load-refusal arm (`internal/cli/lint.go`)."
// REQ-35: "**The envelope code differs by surface and neither is this
// RDR's to change**: the `flow` verbs emit `flow-model-invalid` … while
// `intrastate lint` emits `model-invalid` … The category slug travels in
// the inner finding's `code` on both"
// REQ-73 / `0024:F1`: "a declared model with an undeclared key,
// out-of-domain value, or malformed declaration refuses at load —
// `intrastate lint` exits nonzero with one blocking finding … and `flow
// resolve`/`flow next` refuse `flow-model-invalid` with the same finding."
// ADVERSARIAL
func TestReq34_0024_EveryEmitCategoryRefusesOnBothSurfaces(t *testing.T) {
	cases := map[string]struct {
		src  func(*testing.T) string
		slug table.Category
	}{
		"malformed_emit_declaration": {
			src: func(t *testing.T) string {
				return strings.Replace(declaredDT0024(t), dtDecl0024,
					"\n[emit.verdict]\nkind = \"gauge\"\n", 1)
			},
			slug: table.CatMalformedEmitDeclaration,
		},
		"unknown_emit_key": {
			src: undeclaredKeyDT0024, slug: table.CatUnknownEmitKey,
		},
		"emit_value_out_of_domain": {
			src: outOfDomainDT0024, slug: table.CatEmitValueOutOfDomain,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			model := writeFlowModel(t, tc.src(t))

			t.Run("intrastate lint", func(t *testing.T) {
				ce := requireRefusal(t, "model-invalid", 2,
					"lint", "--model", model, "--as=json")
				assertOneEmitFinding(t, ce.Findings, string(tc.slug), model)
			})

			for _, verb := range []string{"resolve", "next"} {
				t.Run("flow "+verb, func(t *testing.T) {
					args := []string{"flow", verb, "--model", model,
						"--tag", "a=x", "--tag", "b=p", "--as=json"}
					if verb == "resolve" {
						args = append(args, "--outcome", "decide")
					}
					ce := requireRefusal(t, "flow-model-invalid", 2, args...)
					assertOneEmitFinding(t, ce.Findings, string(tc.slug), model)
				})
			}
		})
	}
}

// REQ-37: "All three categories are refusals — nonzero exit, never
// advisory, under every surface that loads the model."
// REQ-42 tail / REQ-98: nothing new lands on the ADVISORY tier.
// ADVERSARIAL — negative REQ.
func TestReq37_0024_NoEmitCategoryEverLandsOnTheAdvisoryTier(t *testing.T) {
	// A clean declared model: its success payload's advisory `findings`
	// list must gain nothing from this RDR.
	model := writeFlowModel(t, declaredDT0024(t))
	stdout := requireSuccess(t, "lint", "--model", model, "--as=json")

	for _, slug := range []string{
		string(table.CatMalformedEmitDeclaration),
		string(table.CatUnknownEmitKey),
		string(table.CatEmitValueOutOfDomain),
	} {
		if strings.Contains(stdout, slug) {
			t.Errorf("a clean lint run mentions %q; all three categories are "+
				"REFUSALS — nonzero exit, never advisory:\n%s", slug, stdout)
		}
	}
}

// REQ-36: "(`lint.go`'s help text naming `codeModelInvalid` for a branch it
// does not emit is pre-existing drift, noted so an implementer does not
// \"fix\" the emitted value to match the prose and break a consumer.)"
// ADVERSARIAL — negative REQ: the emitted value stays `model-invalid`.
func TestReq36_0024_TheLintLoadRefusalStillEmitsModelInvalid(t *testing.T) {
	model := writeFlowModel(t, outOfDomainDT0024(t))

	// The pre-existing drift is not this RDR's to "fix": the EMITTED code
	// stays what consumers already branch on. The inner finding is asserted
	// too, so the outer code cannot be satisfied by an incidental refusal.
	ce := requireRefusal(t, "model-invalid", 2,
		"lint", "--model", model, "--as=json")
	assertOneEmitFinding(t, ce.Findings,
		string(table.CatEmitValueOutOfDomain), model)

	// And `lint.go`'s help text keeps its pre-existing `codeModelInvalid`
	// prose: fixing it to match the emitted value is what REQ-36 forbids.
	if !strings.Contains(readRepoFile(t, "internal/cli/lint.go"),
		"codeModelInvalid") {
		t.Error("`lint.go` no longer names `codeModelInvalid`; that " +
			"pre-existing drift is noted so an implementer does NOT fix " +
			"the emitted value to match the prose and break a consumer")
	}
}

// REQ-43 / `AP`: "The enforcement point is the load pipeline, not lint …
// No graphlint code changes; the advisory tier stays closed (`0006:C17`)
// and the blocking tier's \"at least\" floor (`0006:C3`) is not touched."
// ADVERSARIAL — negative REQ.
func TestReq43_0024_GraphlintGainsNoEmitAwareFinding(t *testing.T) {
	// A model whose declaration is CLEAN but which leaves authoring
	// headroom (a declared member no rule authors, a declared key no rule
	// emits). Graph-lint reports nothing about it, blocking or advisory.
	src := strings.Replace(declaredDT0024(t), dtDecl0024, `
[emit.verdict]
kind = "enum"
[emit.verdict.domain]
route = ["alpha", "beta"]
stop = ["gamma", "never-authored"]

[emit.code]
kind = "enum"
domain = ["1"]

[emit.never-emitted]
kind = "scalar"
`, 1)
	model := writeFlowModel(t, src)
	stdout := requireSuccess(t, "lint", "--model", model, "--as=json")

	for _, token := range []string{"emit", "disposition", "domain"} {
		if strings.Contains(stdout, token) {
			t.Errorf("a clean lint run mentions %q; the enforcement point is "+
				"the LOAD pipeline and graphlint gains no finding:\n%s",
				token, stdout)
		}
	}
}

// REQ-74 / `0024:F2`: "a `scalar`-declared key admits any value —
// declared-but-unvalidated is the documented escape hatch, visible in the
// model source. A model with zero declarations is today's world: nothing
// new fires."
// DOMAIN EDGE
func TestReq74_0024_ScalarIsTheEscapeHatchAndZeroDeclarationsIsTodaysWorld(t *testing.T) {
	t.Run("scalar admits any value", func(t *testing.T) {
		src := strings.Replace(declaredDT0024(t), dtDecl0024,
			"\n[emit.verdict]\nkind = \"scalar\"\n\n[emit.code]\nkind = \"scalar\"\n", 1)
		model := writeFlowModel(t, src)
		requireSuccess(t, "lint", "--model", model, "--as=json")
	})

	t.Run("zero declarations is today's world", func(t *testing.T) {
		model := writeFlowModel(t, undeclaredKeyDT0024Undeclared(t))
		requireSuccess(t, "lint", "--model", model, "--as=json")
	})
}

// undeclaredKeyDT0024Undeclared is the same defective table with the
// `[emit]` declaration REMOVED — the negative control.
func undeclaredKeyDT0024Undeclared(t *testing.T) string {
	t.Helper()

	src := strings.Replace(undeclaredKeyDT0024(t), dtDecl0024, "", 1)
	if strings.Contains(src, "[emit.verdict]") {
		t.Fatal("the declaration removal did not apply")
	}
	return src
}

// REQ-98 / `0024:S3`: "**Expected**: full suite green with the loader
// change and no `models/*.toml` edited; no new refusal reachable; the only
// observable payload delta is C4's appended `dispositions: {}`."
// REQ-99: "**The pass criterion is LOAD-side equivalence, not payload
// byte-identity**"
// HAPPY PATH
func TestReq98_0024_ZeroDeclarationModelsSeeOnlyTheAppendedEmptyField(t *testing.T) {
	// The checked-in model carries no `[emit]` table. Its resolve payload
	// gains `dispositions: {}` and nothing else.
	model := repoModelPath(t)
	art := newFlowArtifact(t, "rdr.artifact")
	bind := artifactBinding("rdr", art)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=resolved",
		"--write", "status=draft", "--write", "gate_passed=false", "--as=json")

	stdout := requireSuccess(t, "flow", "resolve", "--model", model,
		"--artifact", bind, "--outcome", "advance", "--as=json")

	if !strings.Contains(stdout, `"dispositions":{}`) {
		t.Errorf("the zero-declaration payload does not carry "+
			"`\"dispositions\":{}`; that appended empty field is the ONLY "+
			"observable delta:\n%s", stdout)
	}
	data := flowData(t, stdout)
	if got := dispositionsOf(t, data); len(got) != 0 {
		t.Errorf("dispositions = %v; want `{}` on a zero-declaration model",
			got)
	}
}

// REQ-80 / `PH3`: "(i) the 28 inline whole-payload assertions in
// `internal/cli/flow_demand_0011_test.go` — **all 28** pin the
// `\"emit\":{…},\"next\"` adjacency, so every one of them goes red and that
// is the licensed outcome" and "(ii) those 28 ARE the 3 `resolveGolden`
// byte-identity goldens plus the 25-entry `shippedResolveGoldens` sweep in
// the same file, not a further set."
// REQ-81: "the diff must be inspected to confirm every one of the 28
// changed in exactly one way — the inserted `dispositions` key at index 9"
// ASSUMPTION-5: the count is A2's census evidence, not a target.
// ADVERSARIAL
func TestReq80_0024_EveryResolveGoldenCarriesTheInsertedDispositionsKey(t *testing.T) {
	goldens := map[string]string{}
	for outcome, g := range resolveGolden {
		goldens["resolveGolden/"+outcome] = g
	}
	for _, tc := range shippedResolveGoldens {
		goldens["shippedResolveGoldens/"+tc.name+"/"+tc.outcome] = tc.golden
	}

	// The census: the licensed diff is the 3 + 25 sites in this one file,
	// "not a further set".
	if got := len(resolveGolden) + len(shippedResolveGoldens); got != 28 {
		t.Logf("the golden census returned %d sites, not A2's 28; the "+
			"obligation is unchanged and a materially different count is a "+
			"DEVIATION to record, not a REQ failure (ASSUMPTION-5)", got)
	}

	// Every one changed in EXACTLY one way: `dispositions` inserted between
	// `emit` and `next`.
	adjacency := regexp.MustCompile(`"emit":\{[^{}]*\},"dispositions":\{[^{}]*\},"next"`)
	for name, golden := range goldens {
		t.Run(name, func(t *testing.T) {
			if !adjacency.MatchString(golden) {
				t.Errorf("the golden does not carry `dispositions` inserted "+
					"between `emit` and `next`; all %d whole-payload "+
					"assertions pin that adjacency and the regeneration "+
					"inserts exactly one key at index 9:\n%s",
					len(goldens), golden)
			}
		})
	}
}

// REQ-82 / `PH3`: "(iii)
// `internal/cli/decision_table_0010_test.go::TestReq39_
// EmitSitsImmediatelyAfterGatesOnTheWire`, whose 13-key `slices.Equal` list
// becomes 14 keys and whose `NumField() != 14` guard becomes `!= 15` … the
// edit raises both it and the adjacent failure message … Those prose
// strings are part of this licensed edit"
// ADVERSARIAL
func TestReq82_0024_TheLicensedTestReq39EditIsCarried(t *testing.T) {
	src := readRepoFile(t, "internal/cli/decision_table_0010_test.go")

	if !strings.Contains(src, "rt.NumField() != 15") {
		t.Error("`decision_table_0010_test.go` still guards on " +
			"`NumField() != 14`; the licensed edit raises it to 15")
	}
	if strings.Contains(src, "thirteen to fourteen") {
		t.Error("`decision_table_0010_test.go` still says `thirteen to " +
			"fourteen`; the adjacent failure-message prose is part of the " +
			"licensed edit")
	}
	if strings.Contains(src, "the struct declares fourteen fields") {
		t.Error("`decision_table_0010_test.go` still titles its sub-test " +
			"`the struct declares fourteen fields`; the sub-test title is " +
			"part of the licensed edit")
	}
	if !strings.Contains(src, `"dispositions"`) {
		t.Error("`TestReq39`'s `slices.Equal` key list does not carry " +
			"`dispositions`; the 13-key list becomes 14 keys")
	}
}

// REQ-83 / `PH3`: "`TestReq43` (escaped between clear and escape_class)
// asserts relative indices and needs no edit. Any assertion outside this
// list going red is a signal the append did more than C4 licenses —
// investigate, do not update it."
// REQ-50 / A8 Refuted: "**REQ-50 forbids any `TestReq146` edit**".
// ADVERSARIAL — negative REQ: the licensed diff is CLOSED.
func TestReq83_0024_TheLicensedDiffIsClosed(t *testing.T) {
	t.Run("TestReq146 is not edited", func(t *testing.T) {
		src := readRepoFile(t, "internal/table/roundtrip_test.go")
		for _, token := range []string{
			"EmitDecls", "EmitDecl", "Dispositions", "0024",
		} {
			if strings.Contains(src, token) {
				t.Errorf("`roundtrip_test.go` mentions %q; A8 is REFUTED — "+
					"all three of the sweep's sub-tests are scoped to "+
					"`table.Row`, scenario 4 is the sole oracle, and Phase 2 "+
					"owes NO edit to `TestReq146`", token)
			}
		}
	})

	t.Run("TestReq43 asserts relative indices and needs no edit", func(t *testing.T) {
		src := readRepoFile(t, "internal/cli/decision_table_0010_test.go")
		idx := strings.Index(src, "func TestReq43_")
		if idx < 0 {
			t.Fatal("`decision_table_0010_test.go` declares no `TestReq43_`")
		}
		body := src[idx:]
		if end := strings.Index(body, "\nfunc "); end > 0 {
			body = body[:end]
		}
		if strings.Contains(body, "dispositions") {
			t.Error("`TestReq43` mentions `dispositions`; it asserts " +
				"RELATIVE indices and needs no edit — an edit here is a " +
				"signal the append did more than C4 licenses")
		}
	})
}

// REQ-84 / REQ-106 / `0024:S7` / ASSUMPTION-7: in the 0024-first order the
// §D1 registration obligation "discharges instead as a written, checked-in
// note at the `dispositions` field declaration in
// `internal/cli/flow_resolve.go::resolvePayload` naming the PLAN assignment
// and `JDR 0002 §D1` … That note is the Done criterion in that order"
// DOMAIN EDGE
func TestReq106_0024_TheDispositionsFieldCarriesTheD1RegistrationNote(t *testing.T) {
	src := readRepoFile(t, "internal/cli/flow_resolve.go")

	idx := strings.Index(src, "Dispositions map[string]string")
	if idx < 0 {
		t.Fatal("`flow_resolve.go` declares no `Dispositions " +
			"map[string]string` field")
	}
	// The note sits AT the field declaration — the doc comment above it.
	head := src[:idx]
	start := strings.LastIndex(head, "\n\n")
	if start < 0 {
		start = 0
	}
	note := head[start:]

	for _, want := range []string{"PLAN", "JDR 0002"} {
		if !strings.Contains(note, want) {
			t.Errorf("the `Dispositions` field declaration carries no note "+
				"naming %q; in the 0024-first order that written, "+
				"checked-in note IS the Done criterion for `0024:S7`:\n%s",
				want, note)
		}
	}
}

// --- helpers -------------------------------------------------------------

// assertOneEmitFinding asserts the refusal carries EXACTLY one finding
// whose `code` is the category slug and whose `locator` names the model
// file with a grounded, non-`:1` line.
//
// REQ-38: one refusal per run. REQ-30 / `0024:F1`: category slug in the
// finding's `code`, offending file in `locator`. ASSUMPTION-4: the line is
// asserted as non-zero and non-1, not as an integer fixture.
func assertOneEmitFinding(t *testing.T, findings []clierr.Finding, slug, model string) {
	t.Helper()

	if len(findings) != 1 {
		t.Fatalf("the refusal carries %d findings; load is FAIL-FAST — one "+
			"blocking finding per run: %+v", len(findings), findings)
	}
	f := findings[0]
	if f.Code != slug {
		t.Errorf("finding code = %q; want %q — the category slug travels in "+
			"the inner finding's `code`", f.Code, slug)
	}
	base := filepath.Base(model)
	if !strings.Contains(f.Locator, base) {
		t.Errorf("finding locator = %q; want it to name the offending file "+
			"%q", f.Locator, base)
	}
	line := f.Locator[strings.LastIndex(f.Locator, ":")+1:]
	n, err := strconv.Atoi(line)
	if err != nil {
		t.Errorf("finding locator %q carries no line number", f.Locator)
		return
	}
	if n <= 1 {
		t.Errorf("finding locator = %q; the finding must carry the "+
			"offending block's SOURCE LINE, not `:1`", f.Locator)
	}
}
