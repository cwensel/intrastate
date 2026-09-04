package table_test

// RDR 0028 Phase 3c — the anchor's brace discriminator, corrected.
//
// `0028` S14 names three anchor forms BY NAME — "In `anchor`: `\d{4}`,
// `a{2,3}`, a literal `{{`" — and its Expected requires all three
// through: "every brace form reaches RE2 untouched … `{{` is not an
// escape — and none refuses `edit_anchor_invalid`". C1.2 states the same
// posture as a rule ("never for a brace RE2 itself accepts") and C1.4
// states the counter-rule ("carries any other `{…}` form").
//
// The two are jointly unsatisfiable if "a brace RE2 accepts" is read as
// "a brace that compiles": `regexp.Compile` returns nil for EVERY shape
// below, `{artifact}` and `\d{4}` and `{{` and `[{}]` alike. Deviation D5
// records the reconciliation — the discriminator is PLACEHOLDER SHAPE, a
// `{word}` an author wrote meaning a substitution — and this file is its
// witness on both sides at once, which is what the shipped Phase 1 pair
// (`TestReq17And26And65_…` and `TestReq25And28And121_…`) does not do: the
// forms it holds apart are named in one table so a future narrowing or
// widening of the discriminator cannot satisfy one test by breaking the
// other silently.
//
// Phase 3b's two extra cases, `[{}]` and `[{]`, are here for the same
// reason: both are ordinary RE2 character classes an author would write
// to anchor a literal brace, and both were refused before the correction.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// REQ-121 (S14): "In `anchor`: `\d{4}`, `a{2,3}`, a literal `{{`" …
// "every brace form reaches RE2 untouched … none refuses
// `edit_anchor_invalid`."
// REQ-26: "never for a brace RE2 itself accepts."
// REQ-25: "every other `{`, `}` and `$` is passed to RE2 untouched …
// `{{`/`}}` is NOT an escape there."
// REQ-65: "carries any other `{…}` form."
// REQ-17: an owned or recognized tag key is `edit_anchor_invalid` at LINT.
// INPUT EDGE / ADVERSARIAL
//
// One table, both verdicts. The `refuse` column is the whole of what C1.4
// reaches: a `{word}` that is not the one admitted `{tag.<key>}`
// placeholder. Everything else is RE2's.
func TestEditBraceShape0028_PlaceholderShapeIsTheDiscriminator(t *testing.T) {
	for _, tc := range []struct {
		name   string
		anchor string
		refuse bool
		why    string
	}{
		// --- through to RE2, per C1.2 and S14 ---
		{
			name: "s14_bare_doubled_brace", anchor: `^{{x$`,
			why: "S14 names a literal `{{` in the anchor and C1.2 pins it as " +
				"NOT an escape here; the inner `{` is not a name byte",
		},
		{
			name: "s14_escaped_doubled_brace", anchor: `^\{\{x$`,
			why: "the escaped sibling S14's ESCAPED fixture pins; a backslash " +
				"consumes the byte after it, so the scan never sees the brace",
		},
		{
			name: "s14_bounded_quantifier", anchor: `^\d{4}$`,
			why: "an all-digit body is a repeat spec, which RE2 reads as a " +
				"quantifier — S14 requires it to match four digits",
		},
		{
			name: "s14_range_quantifier", anchor: `^a{2,3}$`,
			why: "a comma is not a name byte, so a range spec is never a " +
				"placeholder shape",
		},
		{
			name: "open_ended_quantifier", anchor: `^a{2,}$`,
			why: "the `{N,}` form, for the same reason as `{N,M}`",
		},
		{
			name: "brace_character_class", anchor: `^[{}]$`,
			why: "Phase 3b: an ordinary RE2 character class an author writes " +
				"to anchor a literal brace; the `{` is closed immediately, so " +
				"it names nothing",
		},
		{
			name: "open_brace_character_class", anchor: `^[{]$`,
			why: "Phase 3b: the same, unclosed — no `}` follows, so no name " +
				"is being written",
		},
		{
			name: "empty_braces", anchor: `^{}$`,
			why: "`{}` names nothing, so it is not a placeholder shape",
		},
		{
			name: "escaped_literal_tag_prefix", anchor: `^\{tag\.nnnn\}$`,
			why: "REQ-27: the scan is for the UNESCAPED prefix, so a document " +
				"legitimately carrying `{tag.` stays anchorable",
		},
		{
			name: "declared_observed_tag_key", anchor: `^\| \[{tag.nnnn}\] \|$`,
			why: "the ONE placeholder an anchor admits",
		},

		// --- refused, per C1.4 ---
		{
			name: "argv_placeholder", anchor: `^{artifact}$`, refuse: true,
			why: "0025:C2's argv placeholder is a `{word}` the anchor does not " +
				"admit; passing it through would silently anchor on literal " +
				"braces the document does not carry",
		},
		{
			name: "the_entrys_own_planned_key", anchor: `^{status}$`, refuse: true,
			why: "REQ-18: the value being written never decides where it is " +
				"written",
		},
		{
			name: "some_other_key", anchor: `^{owner}$`, refuse: true,
			why: "any `{word}` that is not `{tag.<key>}` is C1.4's other form",
		},
		{
			name: "undeclared_tag_key", anchor: `^{tag.absent}$`, refuse: true,
			why: "the placeholder scan reaches it and the key is not declared " +
				"`observed`",
		},
		{
			name: "unterminated_tag_placeholder", anchor: `^{tag.nnnn$`, refuse: true,
			why: "a malformed `{tag.…}` placeholder, caught by the prefix scan",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
				`anchor  = "`+strings.ReplaceAll(tc.anchor, `\`, `\\`)+`"`, 1)
			// The clean `replace` references a capture group most of
			// these anchors do not define, so drop the reference.
			body = strings.Replace(body,
				`replace = "- **Status**: {status}"`,
				`replace = "{status}"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)

			_, err := table.Load([]byte(src), "edit-anchor-brace-shape.toml")
			if tc.refuse {
				if err == nil {
					t.Fatalf("anchor %q loaded clean; want `edit_anchor_invalid` "+
						"— %s", tc.anchor, tc.why)
				}
				if got := editCategoryOf(t, src, "edit-anchor-brace-shape.toml"); got != table.CatEditAnchorInvalid {
					t.Errorf("anchor %q: category = %q; want %q — %s",
						tc.anchor, got, table.CatEditAnchorInvalid, tc.why)
				}
				return
			}
			if err != nil {
				t.Fatalf("anchor %q refused at load: %v\nit must reach RE2 "+
					"untouched — %s", tc.anchor, err, tc.why)
			}
		})
	}
}

// REQ-121's operational half: S14 does not merely require `\d{4}` to
// LOAD, it requires it to "compile as a quantifier and match four
// digits". A discriminator that admitted the pattern but mangled its
// bytes on the way through would satisfy the load assertion above and
// still break the scenario, so the anchor is read back verbatim.
func TestEditBraceShape0028_AdmittedAnchorsAreCarriedVerbatim(t *testing.T) {
	for _, anchor := range []string{`^\d{4}$`, `^a{2,3}$`, `^{{x$`, `^[{}]$`} {
		t.Run(anchor, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
				`anchor  = "`+strings.ReplaceAll(anchor, `\`, `\\`)+`"`, 1)
			body = strings.Replace(body,
				`replace = "- **Status**: {status}"`,
				`replace = "{status}"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)

			m := mustLoadEdit(t, src, "edit-anchor-verbatim.toml")
			if got := m.Writers["record"].Edit["status"].Anchor; got != anchor {
				t.Errorf("anchor = %q; want the declared pattern %q carried "+
					"verbatim — every brace form reaches RE2 untouched", got, anchor)
			}
		})
	}
}
