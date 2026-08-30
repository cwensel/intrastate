package table_test

// Phase 3b adversarial tests for RDR 0024, anchored in the record's own
// Failure Modes section.
//
// The Failure Modes' "breaks visibly" leg promises that a declared model
// carrying a defect "refuses at load — `intrastate lint` exits nonzero with
// one blocking finding … (category slug in the finding's `code`, OFFENDING
// FILE IN `locator`)", and that an author can "diagnose from the finding's
// category + locator". `0024:C2` makes that normative rather than
// incidental: "All three categories MUST carry a source line", stamped
// through `atLine`, precisely because "without an analogue every emit
// refusal renders line 0/1 and both promises are false".
//
// These tests attack the LOCATOR, which is the single point on which that
// promise rests. They do not weaken or restate the Phase 1 suite: each one
// authors a model that is legal TOML, legal under `0024:C1`, and carries
// exactly one emit defect — the only variable is HOW the offending rule's
// id or the offending declaration's header is spelled.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// advFixture splices an `[emit]` declaration block and one rule's
// `[rule.emit]` block into the shipped multi-rule RDR fixture.
//
// The fixture is used rather than a hand-rolled minimal model because the
// locator's whole difficulty is a REAL model: `[rule.emit]` is spelled
// identically under every `[[rule]]`, which is why `0024:C2` keys the
// rule-side locator on the rule id in the first place.
func advFixture(t *testing.T, decls, ruleID, emitBlock string) string {
	t.Helper()
	raw, err := os.ReadFile("testdata/rdr-fixture.toml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	src := strings.Replace(string(raw), "[[rule]]", decls+"\n[[rule]]", 1)

	anchor := "id = \"" + ruleID + "\"\n"
	i := strings.Index(src, anchor)
	if i < 0 {
		t.Fatalf("fixture carries no rule %q", ruleID)
	}
	j := strings.Index(src[i:], "\n\n")
	if j < 0 {
		t.Fatalf("no end of rule %q", ruleID)
	}
	return src[:i+j] + "\n[rule.emit]\n" + emitBlock + src[i+j:]
}

// advRefusal loads src and returns the single categorized refusal, failing
// the test when the load succeeds or the error is not a *table.Failure.
func advRefusal(t *testing.T, src string) *table.Failure {
	t.Helper()
	if _, err := table.Load([]byte(src), "adv.toml"); err != nil {
		var f *table.Failure
		if !errors.As(err, &f) {
			t.Fatalf("refusal is not a *table.Failure: %v", err)
		}
		return f
	}
	t.Fatal("the model loaded clean; the fixture carries an emit defect and must refuse")
	return nil
}

// advLineOf reports the 1-based line whose text contains want, or 0.
func advLineOf(src, want string) int {
	for i, l := range strings.Split(src, "\n") {
		if strings.Contains(l, want) {
			return i + 1
		}
	}
	return 0
}

const advEnumDecl = `
[emit.next_command]
kind = "enum"
domain = ["alpha"]
`

// ADV-1. The rule-side locator's premise is that the offending rule's id is
// present and unique by the time the check runs. `0024:C2` states it
// outright: "the id is unique by the time any emit work runs
// (`CatDuplicateRuleID` is enforced in `normalize.go`)".
//
// That premise is FALSE as built. C2 also fixes the position of the two
// emit steps as "immediately after `loadTags`" and "before yielding rows",
// i.e. AHEAD of `normalizeRules` — and `normalize.go` is where
// `CatMalformedRuleShape` ("a [[rule]] carries no id"),
// `CatMalformedRuleID`, and `CatDuplicateRuleID` are all enforced. So
// `checkRuleEmit` reads `l.doc.Rule` before any rule id has been proven
// present, well-formed, or unique, and the locator silently degrades to
// line 0 in exactly the cases the RDR assumed away.
//
// Each sub-test authors ONE emit defect and asserts the refusal carries the
// offending rule's source line, which is what the Failure Modes' "diagnose
// from the finding's category + locator" promises.
func TestAdv1_0024_TheRuleSideLocatorSurvivesAnUnprovenRuleID(t *testing.T) {
	t.Run("a rule carrying no id at all", func(t *testing.T) {
		// The rule's `id` line is removed. `checkRuleEmit` reads
		// `rule.ID == nil` as the EMPTY STRING, so `emitRuleLine` returns 0
		// on sight and the detail renders "rule  emits …" — a refusal that
		// names neither the rule nor a line. The author is handed a defect
		// with no way to locate it, which is the failure the locator exists
		// to prevent.
		src := advFixture(t, advEnumDecl, "continue-prelock", "next_command = \"bad\"\n")
		src = strings.Replace(src, "id = \"continue-prelock\"\n", "", 1)

		f := advRefusal(t, src)
		if f.Category != table.CatEmitValueOutOfDomain {
			t.Fatalf("category = %q, want %q", f.Category, table.CatEmitValueOutOfDomain)
		}
		if strings.Contains(f.Detail, "rule  emits") {
			t.Errorf("the detail names no rule and reads %q; a refusal an author "+
				"cannot attribute to a rule is not diagnosable", f.Detail)
		}
		if f.Line <= 0 {
			t.Errorf("Line = %d; `0024:C2` requires every emit refusal to carry a "+
				"source line, and the Failure Modes promise the offending block is "+
				"locatable from the finding", f.Line)
		}
	})

	t.Run("two rules sharing one id", func(t *testing.T) {
		// `CatDuplicateRuleID` is enforced in normalize.go, which runs AFTER
		// checkRuleEmit -- so the emit refusal fires FIRST and its locator
		// sees two matching `id = "…"` lines. `headerLine` returns 0 on any
		// ambiguity, so the refusal carries no line.
		src := advFixture(t, advEnumDecl, "reconcile-rewind", "next_command = \"bad\"\n")
		src = strings.Replace(src,
			"id = \"continue-prelock-cluster\"", "id = \"reconcile-rewind\"", 1)

		f := advRefusal(t, src)
		if f.Line <= 0 {
			t.Errorf("Line = %d for a duplicated rule id; the locator degrades to "+
				"zero exactly where `0024:C2` assumed uniqueness had already been "+
				"proven, so the refusal carries no position", f.Line)
		}
	})

	t.Run("a rule id equal to the model id", func(t *testing.T) {
		// `0024:C2` requires the anchor to match "the id's VALUE, since
		// `[model]` also carries an `id` key" -- but matching the value only
		// moves the collision: a rule named for its model matches BOTH the
		// `[model]` id line and its own, so `headerLine` again reports 0.
		src := advFixture(t, advEnumDecl, "reconcile-rewind", "next_command = \"bad\"\n")
		src = strings.Replace(src, "id = \"reconcile-rewind\"", "id = \"rdr\"", 1)

		f := advRefusal(t, src)
		if f.Line <= 0 {
			t.Errorf("Line = %d for a rule id equal to the model id; the value-match "+
				"anchor collides with `[model]`'s own id line and the refusal loses "+
				"its position", f.Line)
		}
	})
}

// ADV-2. The rule-side locator matches the id line by EXACT TEXT after
// stripping a trailing comment and trimming space: `headerLine(src, "id = "
// + strconv.Quote(ruleID))`. TOML fixes no such spelling. `id="r"`,
// `id  =  "r"`, and a single-quoted literal string `id = 'r'` are all the
// same document to the decoder, and all three decode to the same rule the
// refusal names -- but only the one canonical spelling recovers a line.
//
// The same brittleness bites the `#` character: `headerLine` truncates at
// the first `#` in the raw line INCLUDING one inside a quoted string, so a
// rule id containing `#` can never match its own anchor.
//
// This is a distinct root cause from ADV-1: the id here is present, unique,
// and unambiguous. The locator loses it purely on formatting, which makes
// `0024:C2`'s "MUST carry a source line" contingent on how the author
// happened to type a space.
func TestAdv2_0024_TheRuleAnchorSurvivesLegalTOMLSpellingsOfTheSameID(t *testing.T) {
	canonical := `id = "reconcile-rewind"`

	for _, tc := range []struct{ name, spelling string }{
		{"no spaces around the equals", `id="reconcile-rewind"`},
		{"extra spaces around the equals", `id  =  "reconcile-rewind"`},
		{"a single-quoted literal string", `id = 'reconcile-rewind'`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := advFixture(t, advEnumDecl, "reconcile-rewind", "next_command = \"bad\"\n")
			src = strings.Replace(src, canonical, tc.spelling, 1)

			f := advRefusal(t, src)
			if f.Category != table.CatEmitValueOutOfDomain {
				t.Fatalf("category = %q, want %q", f.Category, table.CatEmitValueOutOfDomain)
			}
			if want := advLineOf(src, tc.spelling); f.Line != want {
				t.Errorf("Line = %d, want %d for id spelled %s; the anchor is an "+
					"exact text match, so a legal respelling of the SAME id loses "+
					"the line `0024:C2` requires", f.Line, want, tc.spelling)
			}
		})
	}

	t.Run("a rule id containing a hash", func(t *testing.T) {
		// `headerLine` strips at the first '#' before comparing, so the
		// quoted id is truncated mid-string and never matches.
		src := advFixture(t, advEnumDecl, "reconcile-rewind", "next_command = \"bad\"\n")
		src = strings.Replace(src, canonical, `id = "reconcile#rewind"`, 1)

		f := advRefusal(t, src)
		if f.Line <= 0 {
			t.Errorf("Line = %d; the comment strip cuts inside the quoted id, so a "+
				"rule whose id carries '#' can never match its own anchor", f.Line)
		}
	})
}

// ADV-3. The DECLARATION-side locator has the mirror defect.
// `emitHeaderLine` scans for a standalone `[emit.<key>]` header line, but
// `0024:C1` fixes the declaration by its DECODED SHAPE, not by its
// authoring: "an empty table and an absent one both decode to an empty
// map", and the `any`-typed `domain` exists precisely because the decoder,
// not the text, is what defines the declaration.
//
// TOML gives an author a second, equally legal spelling of the identical
// document -- an inline table under a bare `[emit]`. It decodes to the same
// `map[string]sourceEmitDecl`, takes the same C1 arms, and refuses with the
// same category. It simply has no `[emit.<key>]` header line to find, so
// the refusal carries no position and `0024:C2`'s "All three categories
// MUST carry a source line" is false for it.
//
// The rule-side leg is included as the control: it keys on the rule id and
// is unaffected by how the DECLARATION was spelled, which isolates the
// defect to `emitHeaderLine`.
func TestAdv3_0024_TheDeclarationLocatorSurvivesAnInlineTableDeclaration(t *testing.T) {
	t.Run("a malformed inline declaration", func(t *testing.T) {
		decl := "\n[emit]\nnext_command = { kind = \"nope\" }\n"
		src := advFixture(t, decl, "continue-prelock", "next_command = \"alpha\"\n")

		f := advRefusal(t, src)
		if f.Category != table.CatMalformedEmitDeclaration {
			t.Fatalf("category = %q, want %q", f.Category, table.CatMalformedEmitDeclaration)
		}
		if want := advLineOf(src, "next_command = { kind"); f.Line != want {
			t.Errorf("Line = %d, want %d; an inline-table declaration decodes to the "+
				"same shape as the bracketed spelling and takes the same C1 arm, but "+
				"`emitHeaderLine` finds no `[emit.next_command]` header and the "+
				"refusal loses the source line `0024:C2` requires", f.Line, want)
		}
	})

	t.Run("the rule-side leg is unaffected by the declaration spelling", func(t *testing.T) {
		decl := "\n[emit]\nnext_command = { kind = \"enum\", domain = [\"alpha\"] }\n"
		src := advFixture(t, decl, "continue-prelock", "next_command = \"bad\"\n")

		f := advRefusal(t, src)
		if f.Category != table.CatEmitValueOutOfDomain {
			t.Fatalf("category = %q, want %q", f.Category, table.CatEmitValueOutOfDomain)
		}
		if want := advLineOf(src, `id = "continue-prelock"`); f.Line != want {
			t.Errorf("Line = %d, want %d; the rule-side locator keys on the rule id "+
				"and must not depend on how the declaration was authored", f.Line, want)
		}
	})
}
