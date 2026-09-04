package table_test

// RDR 0028 — the `edit` write carrier's DECLARATION grammar and its six
// load-time defect categories (C1.1, C1.2, C1.4, C1.6).
//
// This file holds the LOAD half of C1. The apply half — selection,
// refusal ordering, byte preservation, `<clear>` — lives in
// `internal/accessor/edit_apply_0028_test.go`, because C1.4's
// `what lint does NOT prove:` fences the two apart by design: a model is
// linted without its artifacts and without the invocation's bindings.
//
// Two coverage floors this file exists to hold.
//
// FIRST: an implementation that refuses correctly at the call site but
// never REGISTERS its categories is invisible to every consumer that
// enumerates `table.Categories()`. So each of the six is asserted twice —
// by the category a refusal carries AND by that category's membership in
// the registered set (C1.4 `registration:`).
//
// SECOND: the two `edit_*` sets are DISJOINT. Only C1.4's six wire
// strings register in `Categories()`; C1.3/C1.5's six apply-time reason
// tokens (`edit_anchor_unmatched`, `edit_anchor_ambiguous`,
// `edit_anchor_collision`, `edit_anchor_unstable`, `edit_value_multiline`,
// `edit_clear_undeclared`) ride the Detail and register NOWHERE. A suite
// that registered all twelve would pass a naive membership check while
// telling a consumer that a stale anchor is a lint defect (req-list A-12).

import (
	"slices"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// --- the edit-carrier fixture --------------------------------------------

// editCarrierModel is a minimal model whose write entry carries `edit`
// and NOTHING else — no `path`, no `command`. Authored rather than
// derived, because the whole point of S1 is that an `edit`-only entry
// LOADS: substituting into a path-backed fixture would leave the old
// locator behind and the "neither" arm would never be exercised.
//
// The reader is `path`-backed so the model validates as a whole while the
// writer under test is the only `edit` entry in it.
const editCarrierModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "editflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.nnnn]
provenance = "observed"
kind = "string"
single_valued = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"

[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "- **Status**: {status}"

[initial]
status = "draft"
`

// editWriteBlock is the write entry as authored above, so a mutant can
// swap the whole block rather than a line whose text recurs.
const editWriteBlock = `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "- **Status**: {status}"`

// editSwap replaces old with replacement exactly once and fails loudly if
// the anchor text is not in the fixture — a silently non-applied mutant
// would test the clean fixture twice and pass for the wrong reason.
func editSwap(t *testing.T, base, old, replacement string) string {
	t.Helper()

	out := strings.Replace(base, old, replacement, 1)
	if out == base {
		t.Fatalf("mutant substitution did not apply; the block\n%s\nis not in the fixture", old)
	}
	return out
}

// editCategoryOf loads src and returns the category its refusal carries.
// The oracle is the CATEGORY, never "an error occurred" and never the
// message text: C1.4 fixes the wire strings as the contract and states
// outright that the hardcoded message is not (REQ-9).
func editCategoryOf(t *testing.T, src, id string) table.Category {
	t.Helper()

	_, err := table.Load([]byte(src), id)
	if err == nil {
		t.Fatalf("%s loaded clean; want a categorized refusal", id)
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("%s refused with a non-categorized error: %v", id, err)
	}
	return cat
}

// mustLoadEdit loads a fixture C1 says must load clean.
func mustLoadEdit(t *testing.T, src, id string) *table.Model {
	t.Helper()

	m, err := table.Load([]byte(src), id)
	if err != nil {
		t.Fatalf("%s must load clean; refused: %v", id, err)
	}
	return m
}

// --- C1.1: the carrier and its grammar -----------------------------------

// REQ-1: "[write.<id>.edit.<key>]           # exactly one table per member
// of `keys`; a table for a key not in `keys` is a defect (C1.4)"
// HAPPY PATH
//
// The oracle is the DECODED rule on the loaded model, not merely that the
// document loaded: a loader that accepts `edit.<key>` and discards it
// would pass a load-clean assertion while losing the whole carrier.
func TestReq1_EditTablesDecodePerDeclaredKey(t *testing.T) {
	m := mustLoadEdit(t, editCarrierModel, "edit-carrier.toml")

	acc := m.Writers["record"]
	if acc.Edit == nil {
		t.Fatalf("write.record carries no decoded `edit` table set")
	}
	rule, ok := acc.Edit["status"]
	if !ok {
		t.Fatalf("write.record.edit = %#v; want a rule table for the `status` key", acc.Edit)
	}
	if got, want := rule.Anchor, `^- \*\*Status\*\*: (.*)$`; got != want {
		t.Errorf("edit.status anchor = %q; want the declared pattern %q", got, want)
	}
	if got, want := rule.Replace, "- **Status**: {status}"; got != want {
		t.Errorf("edit.status replace = %q; want the declared template %q", got, want)
	}
}

// REQ-2: "anchor  = \"<RE2 pattern>\"         # a line is SELECTED when the
// pattern matches anywhere in it (terminator excluded); authors pin
// `^…$`; must select exactly one line (C1.3)"
// REQ-3: "replace = \"<template>\"            # the WHOLE replacement line,
// terminator excluded (C1.2) — never the matched span"
// REQ-14: "the carrier is `edit`; the rule tables are `edit.<key>`; the
// fields are `anchor` / `replace` / `clear`."
// BOUNDARY
//
// The grammar's field NAMES are the contract an author reproduces. A
// loader that accepted `pattern`/`line`/`delete` would ship a different
// declaration language under the same contract id.
func TestReq2And3And14_TheRuleFieldsAreAnchorReplaceAndClear(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			name: "anchor_renamed_is_not_the_field",
			body: strings.Replace(editWriteBlock, "anchor  =", "pattern =", 1),
		},
		{
			name: "replace_renamed_is_not_the_field",
			body: strings.Replace(editWriteBlock, "replace =", "line    =", 1),
		},
		{
			name: "clear_renamed_is_not_the_field",
			body: editWriteBlock + "\nremove  = \"line\"",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, editWriteBlock, tc.body)
			// An unknown field inside a declared table is 0002's
			// `unknown_schema_field`; what this asserts is that the
			// renamed spelling is NOT silently accepted as the field.
			if _, err := table.Load([]byte(src), "edit-field-rename.toml"); err == nil {
				t.Errorf("a rule table using the renamed field loaded clean; "+
					"C1.1 fixes the field names as `anchor`/`replace`/`clear` (%s)", tc.name)
			}
		})
	}
}

// REQ-4: "clear   = \"line\"                  # optional; the disposition
// of a planned `<clear>` (C1.5); absent ⇒ `<clear>` refuses before
// mutation"
// INPUT EDGE
//
// `clear` is OPTIONAL, so the clean fixture — which omits it — must load,
// and a rule declaring `"line"` must load too and carry the value. Both
// halves matter: a loader requiring `clear` and a loader dropping it are
// different defects and only asserting one catches neither.
func TestReq4_ClearIsOptionalAndCarriesItsDeclaredDisposition(t *testing.T) {
	t.Run("omitted_loads_and_is_empty", func(t *testing.T) {
		m := mustLoadEdit(t, editCarrierModel, "edit-carrier.toml")
		if got := m.Writers["record"].Edit["status"].Clear; got != "" {
			t.Errorf("edit.status clear = %q on a rule that omits it; want empty", got)
		}
	})

	t.Run("line_loads_and_carries", func(t *testing.T) {
		src := editSwap(t, editCarrierModel, editWriteBlock,
			editWriteBlock+"\nclear   = \"line\"")
		m := mustLoadEdit(t, src, "edit-clear-line.toml")
		if got, want := m.Writers["record"].Edit["status"].Clear, "line"; got != want {
			t.Errorf("edit.status clear = %q; want the declared %q", got, want)
		}
	})
}

// REQ-5: "role, keys, timeout               # unchanged (0002, 0004:C15)"
// NEGATIVE REQ — these three fields are not redefined by this record;
// `timeout` stays unconditionally required for every entry regardless of
// carrier (A11).
// DOMAIN EDGE
//
// GREEN-BY-DESIGN against existing behaviour on the `role`/`keys` half.
// The load-bearing half is `timeout`: C1.3's `no subprocess:` clause means
// an `edit` write spawns nothing and takes no wall-clock risk, which is
// exactly the argument an implementer would use to make `timeout`
// conditional on the carrier. A11 forbids it.
func TestReq5_RoleKeysAndTimeoutAreUnchangedByTheEditCarrier(t *testing.T) {
	m := mustLoadEdit(t, editCarrierModel, "edit-carrier.toml")
	acc := m.Writers["record"]
	if got, want := acc.Role, "record"; got != want {
		t.Errorf("write.record role = %q; want %q", got, want)
	}
	if got, want := acc.Keys, []string{"status"}; !slices.Equal(got, want) {
		t.Errorf("write.record keys = %#v; want %#v", got, want)
	}
	if got, want := acc.Timeout, "2s"; got != want {
		t.Errorf("write.record timeout = %q; want %q", got, want)
	}

	t.Run("timeout_stays_required_on_an_edit_entry", func(t *testing.T) {
		src := editSwap(t, editCarrierModel, "timeout = \"2s\"\nread_back = true", "read_back = true")
		if _, err := table.Load([]byte(src), "edit-no-timeout.toml"); err == nil {
			t.Errorf("an `edit` entry without `timeout` loaded clean; " +
				"A11 keeps `timeout` unconditionally required regardless of carrier")
		}
	})
}

// REQ-6: "carrier: exactly one of `path` / `command` / `edit` per entry —
// 0025:C1's exactly-one rule with one member appended, not replaced"
// REQ-8: "The \"neither\" arm keeps its 0025:C5 wire string
// `command_and_path_conflict`, but its PREDICATE necessarily widens to
// \"none of the three\" — `carrierDefect`'s `case !hasPath && !hasCommand:`
// gains `&& !hasEdit`, or an `edit`-only entry would be refused as
// carrier-less (S1 asserts it LOADS)."
// HAPPY PATH
//
// S1's positive arm. This is the assertion that catches the whole
// implementation being unreachable: if the widened predicate is missed,
// every `edit` entry refuses as carrier-less and no other test in this
// file can distinguish "not implemented" from "implemented and rejected".
func TestReq6And8_AnEditOnlyEntryLoadsAndIsNotCarrierLess(t *testing.T) {
	m := mustLoadEdit(t, editCarrierModel, "edit-only.toml")

	acc := m.Writers["record"]
	if acc.Path != "" {
		t.Errorf("write.record path = %q; the edit-only fixture declares none", acc.Path)
	}
	if len(acc.Command) != 0 {
		t.Errorf("write.record command = %#v; the edit-only fixture declares none", acc.Command)
	}
	if len(acc.Edit) == 0 {
		t.Errorf("write.record carries no `edit` table; the exactly-one rule " +
			"admits `edit` as a third carrier, not as an ignored field")
	}
}

// REQ-7: "`edit` is admissible on WRITE entries only (a read or gate entry
// carrying `edit` is `edit_carrier_conflict`, C1.4)"
// REQ-63: "edit_carrier_conflict    # `edit` beside `path` or `command`; or
// `edit` on a read/gate entry"
// REQ-10: "`edit_carrier_conflict` follows the \"both\" arm's discipline —
// keyed on the `edit` table being PRESENT, not on it being non-empty"
// ADVERSARIAL
//
// S1's four negative arms plus REQ-10's presence-not-emptiness arm. The
// empty-table case is the mutant this exists to kill: an implementation
// keyed on `len(Edit) != 0` treats `[write.x.edit]` beside a `path` as a
// silently ignored second carrier, which is precisely the failure
// 0025:C1's "both" arm was written to prevent.
func TestReq7And10And63_EditBesideAnotherCarrierOrOffWriteIsCarrierConflict(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			name: "edit_beside_path",
			body: `[write.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^x$"
replace = "y"`,
		},
		{
			name: "edit_beside_command",
			body: `[write.record]
role = "record"
command = ["tools/write", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^x$"
replace = "y"`,
		},
		{
			name: "empty_edit_beside_path_is_still_a_conflict",
			body: `[write.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, editWriteBlock, tc.body)
			if got := editCategoryOf(t, src, "edit-conflict.toml"); got != table.CatEditCarrierConflict {
				t.Errorf("category = %q; want %q — C1.4 keys the conflict on the "+
					"`edit` table being PRESENT, not on it being non-empty",
					got, table.CatEditCarrierConflict)
			}
		})
	}

	// The read/gate arms need their own entry, since the write entry is
	// where the clean fixture's `edit` lives.
	for _, tc := range []struct {
		name  string
		old   string
		added string
	}{
		{
			name: "edit_on_a_read_entry",
			old: `[read.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"`,
			added: `[read.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"

[read.record.edit.status]
anchor  = "^x$"
replace = "y"`,
		},
		{
			name: "edit_on_a_gate_entry",
			old:  "[initial]",
			added: `[gate.approval]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"

[gate.approval.edit.status]
anchor  = "^x$"
replace = "y"

[initial]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, tc.old, tc.added)
			if got := editCategoryOf(t, src, "edit-off-write.toml"); got != table.CatEditCarrierConflict {
				t.Errorf("category = %q; want %q — `edit` is admissible on WRITE "+
					"entries only", got, table.CatEditCarrierConflict)
			}
		})
	}
}

// REQ-9: "Its hardcoded message text, which names two carriers, is updated
// to name three; the wire string is the contract and the message is not
// (C1.4)."
// DOMAIN EDGE
//
// The wire string of the "neither" arm stays `command_and_path_conflict`
// (REQ-8) while its MESSAGE must name three carriers. Asserting the
// message is normally forbidden — here the clause makes the message text
// itself the obligation, and pins the wire string as the thing that must
// NOT move. Both halves are asserted so a rename of either is caught.
func TestReq9_TheNeitherArmKeepsItsWireStringAndNamesThreeCarriers(t *testing.T) {
	// An entry carrying none of the three carriers.
	src := editSwap(t, editCarrierModel, editWriteBlock, `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true`)

	_, err := table.Load([]byte(src), "edit-no-carrier.toml")
	if err == nil {
		t.Fatalf("an entry carrying none of the three carriers loaded clean")
	}
	cat, ok := table.CategoryOf(err)
	if !ok {
		t.Fatalf("refusal is not categorized: %v", err)
	}
	if want := table.CatCommandAndPathConflict; cat != want {
		t.Errorf("category = %q; want the UNCHANGED 0025:C5 wire string %q — "+
			"the predicate widens to \"none of the three\", the string does not move",
			cat, want)
	}
	if msg := err.Error(); !strings.Contains(msg, "edit") {
		t.Errorf("the \"neither\" arm's message %q does not name `edit`; "+
			"C1.1 requires the text be updated to name three carriers", msg)
	}
}

// REQ-64: "edit_key_mismatch        # a `keys` member with no `edit.<key>`
// table, or an `edit.<key>` table for a key not in `keys`"
// BOUNDARY
//
// S2, both directions. The bijection between `keys` and `edit.<key>` is
// what makes "every planned key has a rule" a lint-time guarantee
// (REQ-74) rather than an apply-time surprise.
func TestReq64_KeysAndEditTablesMustBeInBijection(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{
			name: "a_keys_member_with_no_rule_table",
			body: `[write.record]
role = "record"
keys = ["status", "owner"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "- **Status**: {status}"`,
		},
		{
			name: "a_rule_table_for_a_key_not_in_keys",
			body: `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "- **Status**: {status}"

[write.record.edit.owner]
anchor  = "^- \\*\\*Owner\\*\\*: (.*)$"
replace = "- **Owner**: {owner}"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, editWriteBlock, tc.body)
			if got := editCategoryOf(t, src, "edit-key-mismatch.toml"); got != table.CatEditKeyMismatch {
				t.Errorf("category = %q; want %q", got, table.CatEditKeyMismatch)
			}
		})
	}
}

// REQ-65: "edit_anchor_invalid      # anchor fails to compile as RE2
// (checked with every `{tag.<key>}` replaced by a quoted probe), names an
// undeclared tag key, or carries any other `{…}` form"
// REQ-26: "`edit_anchor_invalid` fires when a `{tag.…}` placeholder is
// malformed or names a key not declared `observed`, or when the pattern
// fails to compile — never for a brace RE2 itself accepts."
// REQ-17: "a declared key of either kind is structurally unbindable at
// invocation, and `{tag.<key>}` naming one is `edit_anchor_invalid` at
// LINT (C1.4) rather than a guaranteed runtime `execution_failure`"
// ADVERSARIAL
//
// S3. The last two rows are the ones the clause argues hardest for: an
// `owned` or `recognized` tag key is structurally UNBINDABLE at
// invocation (`flow_input.go::parseTags` refuses both), so admitting it
// here would ship a model that can only ever fail at runtime.
func TestReq17And26And65_AnchorDefectsAreEditAnchorInvalid(t *testing.T) {
	for _, tc := range []struct {
		name   string
		anchor string
	}{
		{name: "uncompilable_re2", anchor: `^- \*\*Status\*\*: (unclosed$`},
		{name: "undeclared_tag_key", anchor: `^\| {tag.absent} \|$`},
		{name: "owned_tag_key_is_unbindable", anchor: `^{tag.status}$`},
		{name: "recognized_tag_key_is_unbindable", anchor: `^{tag.recognized}$`},
		{name: "unterminated_tag_placeholder", anchor: `^{tag.nnnn$`},
		{name: "a_non_tag_brace_placeholder", anchor: `^{artifact}$`},
		{name: "the_entrys_own_planned_key", anchor: `^{status}$`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
				`anchor  = "`+strings.ReplaceAll(tc.anchor, `\`, `\\`)+`"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)
			if got := editCategoryOf(t, src, "edit-anchor-invalid.toml"); got != table.CatEditAnchorInvalid {
				t.Errorf("category = %q; want %q", got, table.CatEditAnchorInvalid)
			}
		})
	}
}

// REQ-18: "The entry's own planned keys are NOT admissible in `anchor`:
// the value being written never decides where it is written"
// NEGATIVE REQ.
// ADVERSARIAL
//
// Asserted as its own test rather than folded into the table above,
// because the reason is a security property and not a syntax rule: an
// anchor that interpolated the value being written would let a caller
// choose the line the write lands on.
func TestReq18_AnEntrysOwnPlannedKeyIsNotAnAnchorPlaceholder(t *testing.T) {
	body := strings.Replace(editWriteBlock,
		`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
		`anchor  = "^- \\*\\*Status\\*\\*: {status}$"`, 1)
	src := editSwap(t, editCarrierModel, editWriteBlock, body)

	if got := editCategoryOf(t, src, "edit-anchor-own-key.toml"); got != table.CatEditAnchorInvalid {
		t.Errorf("category = %q; want %q — the value being written never "+
			"decides where it is written", got, table.CatEditAnchorInvalid)
	}
}

// REQ-25: "In `anchor` ONLY the fixed prefix `{tag.` opens a placeholder,
// scanned to its closing `}`; every other `{`, `}` and `$` is passed to
// RE2 untouched, so `\\d{4}`, `a{2,3}` and `^…$` are ordinary pattern text
// and `{{`/`}}` is NOT an escape there."
// REQ-28: "the escaping dialect in `anchor` is RE2's, in `replace` it is
// this clause's"
// REQ-121 (S14, anchor half): "In `anchor`: `\\d{4}`, `a{2,3}`, a literal
// `{{`."
// INPUT EDGE
//
// The asymmetry IS the scenario. A shared brace scanner over both
// templates would refuse `\d{4}` as a malformed placeholder — a quantifier
// every non-trivial anchor uses — and the suite would still be green on
// every `replace` case.
func TestReq25And28And121_AnchorBracesReachRE2Untouched(t *testing.T) {
	for _, tc := range []struct {
		name   string
		anchor string
	}{
		{name: "bounded_quantifier", anchor: `^\d{4}$`},
		{name: "range_quantifier", anchor: `^a{2,3}$`},
		{name: "doubled_brace_is_not_an_escape", anchor: `^\{\{x$`},
		{name: "dollar_and_caret_are_pattern_text", anchor: `^- \*\*Status\*\*: (.*)$`},
		{name: "escaped_literal_tag_prefix", anchor: `^\{tag\.nnnn\}$`},
		{name: "a_declared_observed_tag_key", anchor: `^\| \[{tag.nnnn}\]\(.*\) \|`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
				`anchor  = "`+strings.ReplaceAll(tc.anchor, `\`, `\\`)+`"`, 1)
			// The clean replace references a group the mutated anchors
			// mostly do not define, so drop the group reference too.
			body = strings.Replace(body,
				`replace = "- **Status**: {status}"`,
				`replace = "{status}"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)
			if _, err := table.Load([]byte(src), "edit-anchor-braces.toml"); err != nil {
				t.Errorf("anchor %q refused at load: %v — every brace form RE2 "+
					"itself accepts must reach it untouched", tc.anchor, err)
			}
		})
	}
}

// REQ-27: "A line whose CONTENT carries the literal text `{tag.` is
// anchored with RE2's own escape (`\\{tag\\.`): the placeholder scan is for
// the unescaped literal prefix, so an escaped brace is ordinary pattern
// text and reaches RE2 untouched."
// ADVERSARIAL
//
// The escaped form must NOT be read as a placeholder — otherwise a
// document that legitimately contains the six bytes `{tag.` becomes
// unanchorable. Asserted by the DECODED anchor still carrying the escape,
// not merely by the model loading: a scanner that stripped the escape
// would load clean and then match the wrong lines at apply.
func TestReq27_AnEscapedTagPrefixIsPatternTextNotAPlaceholder(t *testing.T) {
	body := strings.Replace(editWriteBlock,
		`anchor  = "^- \\*\\*Status\\*\\*: (.*)$"`,
		`anchor  = "^\\{tag\\.nnnn\\} = (.*)$"`, 1)
	src := editSwap(t, editCarrierModel, editWriteBlock, body)

	m := mustLoadEdit(t, src, "edit-anchor-escaped-tag.toml")
	if got, want := m.Writers["record"].Edit["status"].Anchor, `^\{tag\.nnnn\} = (.*)$`; got != want {
		t.Errorf("anchor = %q; want the declared pattern %q carried verbatim — "+
			"the placeholder scan is for the UNESCAPED prefix", got, want)
	}
}

// REQ-66: "edit_template_invalid    # replace carries an unknown
// placeholder, another key's placeholder, a group reference the anchor
// does not define, or a malformed `${…}`/`{…}` token"
// REQ-15: "replace admits: literal text | ${N} and ${name} — the anchor's
// capture groups | {<key>} — the planned value of THE key this table is
// named for, and no other key"
// REQ-24: "In `replace` the vocabulary is closed ... any other `{…}` or
// `$…` form is `edit_template_invalid`, including a bare `$1`/`$name`
// without braces (`Expand` accepts those; this clause does not)."
// ADVERSARIAL
//
// S4, plus S14's `replace` half. The bare-`$1` row is the one that
// separates this clause's parser from a delegation to `regexp.Expand`:
// `Expand` accepts `$1` and this clause refuses it, so an implementation
// that reached for the stdlib helper fails exactly here (REQ-23).
func TestReq15And24And66_ReplaceVocabularyIsClosed(t *testing.T) {
	for _, tc := range []struct {
		name    string
		replace string
	}{
		{name: "unknown_placeholder", replace: "- **Status**: {nope}"},
		{name: "another_keys_placeholder", replace: "- **Status**: {owner}"},
		{name: "a_tag_placeholder_is_not_admitted_in_replace", replace: "- **Status**: {tag.nnnn}"},
		{name: "group_the_anchor_does_not_define", replace: "- **Status**: ${7}"},
		{name: "named_group_the_anchor_does_not_define", replace: "- **Status**: ${nope}"},
		{name: "bare_dollar_digit", replace: "- **Status**: $1"},
		{name: "bare_dollar_name", replace: "- **Status**: $status"},
		{name: "unterminated_brace_token", replace: "- **Status**: {status"},
		{name: "unterminated_dollar_brace_token", replace: "- **Status**: ${1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`replace = "- **Status**: {status}"`,
				`replace = "`+tc.replace+`"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)
			if got := editCategoryOf(t, src, "edit-template-invalid.toml"); got != table.CatEditTemplateInvalid {
				t.Errorf("category = %q; want %q", got, table.CatEditTemplateInvalid)
			}
		})
	}
}

// REQ-22: "escapes: `$$` emits a literal `$`; `{{` and `}}` emit literal
// braces in `replace`; a group that did not participate in the match
// expands to the empty string"
// REQ-121 (S14, replace half)
// INPUT EDGE
//
// The admitted half of the closed vocabulary. A parser that refused the
// escapes would be "closed" in the trivial sense and unable to write a
// line containing a dollar or a brace.
func TestReq22And121_ReplaceEscapesAndDefinedGroupsAreAdmitted(t *testing.T) {
	for _, tc := range []struct {
		name    string
		replace string
	}{
		{name: "doubled_dollar_escape", replace: "- **Status**: $$ {status}"},
		{name: "doubled_braces_escape", replace: "- **Status**: {{{status}}}"},
		{name: "numbered_group_the_anchor_defines", replace: "- **Status**: {status} ${1}"},
		{name: "the_owned_key_alone", replace: "{status}"},
		{name: "literal_text_only", replace: "- **Status**: Final"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.Replace(editWriteBlock,
				`replace = "- **Status**: {status}"`,
				`replace = "`+tc.replace+`"`, 1)
			src := editSwap(t, editCarrierModel, editWriteBlock, body)
			if _, err := table.Load([]byte(src), "edit-template-escapes.toml"); err != nil {
				t.Errorf("replace %q refused at load: %v — the escapes and a "+
					"defined group are inside the closed vocabulary", tc.replace, err)
			}
		})
	}
}

// REQ-67: "edit_clear_invalid       # `clear` outside the closed set
// {\"line\"}"
// REQ-81: "A blank-the-cell disposition (`clear = { replace = … }`) is
// deferred: it is a document-shape decision, and v1's closed set is
// {\"line\"}"
// BOUNDARY
//
// S5. The `clear = { replace = … }` row is the deferred form the clause
// names by hand: admitting it would ship the successor's feature under
// this record's contract.
func TestReq67And81_ClearIsClosedAtTheSingleLiteralLine(t *testing.T) {
	for _, tc := range []struct {
		name  string
		clear string
	}{
		{name: "another_literal", clear: `clear   = "cell"`},
		{name: "the_empty_string", clear: `clear   = ""`},
		{name: "case_variant", clear: `clear   = "Line"`},
		{name: "the_deferred_table_form", clear: `clear   = { replace = "-" }`},
		{name: "a_boolean", clear: `clear   = true`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, editWriteBlock,
				editWriteBlock+"\n"+tc.clear)
			if got := editCategoryOf(t, src, "edit-clear-invalid.toml"); got != table.CatEditClearInvalid {
				t.Errorf("category = %q; want %q", got, table.CatEditClearInvalid)
			}
		})
	}
}

// --- C1.6: `{tag.<key>}` in a command argv --------------------------------

// cmdTagModel is a model whose read, gate and write entries are all
// `command`-backed, so C1.6's argv rules have all three capabilities to
// fire on. `nnnn` is declared `observed`, which is the only provenance
// `{tag.<key>}` admits.
const cmdTagModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "tagflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.nnnn]
provenance = "observed"
kind = "string"
single_valued = true

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "final"]
single_valued = true
required = true

[read.record]
role = "record"
command = ["rdr", "status", "-json", "{tag.nnnn}", "{artifact}"]
keys = ["status"]
timeout = "2s"

[write.record]
role = "record"
command = ["rdr", "set", "{tag.nnnn}", "{artifact}"]
keys = ["status"]
timeout = "2s"
read_back = true

[gate.approval]
role = "record"
command = ["rdr", "check", "{tag.nnnn}", "{artifact}"]
exit_verdicts = { "0" = "allow", "1" = "deny" }
keys = ["status"]
timeout = "2s"

[initial]
status = "draft"
`

// REQ-83: "[read.<id> | gate.<id> | write.<id>]  command = [...,
// \"{tag.<key>}\", ...]"
// REQ-84: "admission: `{tag.<key>}` joins 0025:C2's vocabulary as a
// family, WHOLE-ELEMENT only under 0025:C2's substitution rule, replaced
// by the bound value of a tag key the model DECLARES"
// HAPPY PATH
//
// S7's admitted arm, on all three capabilities. The declared form at a
// NON-LEADING position lints clean — that is the whole point of the
// family, and the negative tests below are meaningless without it.
func TestReq83And84_ADeclaredTagPlaceholderIsAdmittedAtANonLeadingPosition(t *testing.T) {
	m := mustLoadEdit(t, cmdTagModel, "cmd-tag.toml")

	for _, tc := range []struct {
		name string
		argv []string
	}{
		{name: "read", argv: m.Readers["record"].Command},
		{name: "write", argv: m.Writers["record"].Command},
		{name: "gate", argv: m.Gates["approval"].Command},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !slices.Contains(tc.argv, "{tag.nnnn}") {
				t.Errorf("%s argv = %#v; want the declared `{tag.nnnn}` element "+
					"carried whole", tc.name, tc.argv)
			}
		})
	}
}

// REQ-85: "`<key>` undeclared is `command_unknown_placeholder` (0025:C5,
// unchanged wire string); a `{…}` element that is neither `{artifact}` nor
// a declared `{tag.<key>}` stays `command_unknown_placeholder`."
// ADVERSARIAL
//
// S7's unrecognized arms. The wire string must NOT move to an `edit_*`
// category: a consumer keyed on `command_unknown_placeholder` would stop
// seeing the defect it already handles.
func TestReq85_AnUndeclaredOrUnrecognizedBraceElementStaysCommandUnknownPlaceholder(t *testing.T) {
	for _, tc := range []struct {
		name string
		el   string
	}{
		{name: "undeclared_tag_key", el: "{tag.absent}"},
		{name: "not_artifact_and_not_a_tag", el: "{whatever}"},
		{name: "a_bare_tag_prefix", el: "{tag.}"},
		{name: "not_whole_element", el: "x{tag.nnnn}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, cmdTagModel,
				`command = ["rdr", "status", "-json", "{tag.nnnn}", "{artifact}"]`,
				`command = ["rdr", "status", "-json", "`+tc.el+`", "{artifact}"]`)
			if got := editCategoryOf(t, src, "cmd-tag-unknown.toml"); got != table.CatCommandUnknownPlaceholder {
				t.Errorf("category = %q; want the UNCHANGED 0025:C5 wire string %q",
					got, table.CatCommandUnknownPlaceholder)
			}
		})
	}
}

// REQ-68: "edit_tag_argv0           # a `{tag.<key>}` element at argv0 of a
// read, gate or write entry's `command` (C1.6)"
// REQ-86: "POSITION: a `{tag.<key>}` element at argv0 is `edit_tag_argv0`
// at LINT (C1.4) — the executable is the one word a reviewer must be able
// to read off the model"
// REQ-70: "`edit_tag_argv0` ... fires on `command` entries of every kind,
// including entries carrying no `edit` table"
// ADVERSARIAL
//
// S7's position arm. Every fixture here carries NO `edit` table at all —
// that is REQ-70's claim: the category's `edit_` prefix names its owning
// contract, not the carrier it fires on.
func TestReq68And70And86_ADeclaredTagAtArgv0IsEditTagArgv0(t *testing.T) {
	for _, tc := range []struct {
		name string
		old  string
		new  string
	}{
		{
			name: "read_entry",
			old:  `command = ["rdr", "status", "-json", "{tag.nnnn}", "{artifact}"]`,
			new:  `command = ["{tag.nnnn}", "status", "-json", "{artifact}"]`,
		},
		{
			name: "write_entry",
			old:  `command = ["rdr", "set", "{tag.nnnn}", "{artifact}"]`,
			new:  `command = ["{tag.nnnn}", "set", "{artifact}"]`,
		},
		{
			name: "gate_entry",
			old:  `command = ["rdr", "check", "{tag.nnnn}", "{artifact}"]`,
			new:  `command = ["{tag.nnnn}", "check", "{artifact}"]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, cmdTagModel, tc.old, tc.new)
			if got := editCategoryOf(t, src, "cmd-tag-argv0.toml"); got != table.CatEditTagArgv0 {
				t.Errorf("category = %q; want %q — the executable is the one word "+
					"a reviewer must be able to read off the model", got, table.CatEditTagArgv0)
			}
		})
	}

	t.Run("argv1_lints_clean", func(t *testing.T) {
		// The position, not the key, is what the category decides.
		mustLoadEdit(t, cmdTagModel, "cmd-tag-argv1.toml")
	})
}

// REQ-87: "Statically decidable, so it is refused where it is visible
// rather than at spawn; `{artifact}` is unaffected, having no such rule
// and no such reviewer promise to break"
// NEGATIVE REQ: `{artifact}` at argv0 stays admitted (S7).
// DOMAIN EDGE
//
// The rule is on THIS FAMILY only. An implementation that generalised it
// to "no placeholder at argv0" would break `{artifact}`-at-argv0 models
// 0025 already ships, and would do so while every `{tag.…}` assertion
// above stayed green.
func TestReq87_ArtifactAtArgv0StaysAdmitted(t *testing.T) {
	src := editSwap(t, cmdTagModel,
		`command = ["rdr", "status", "-json", "{tag.nnnn}", "{artifact}"]`,
		`command = ["{artifact}", "status", "-json", "{tag.nnnn}"]`)

	if _, err := table.Load([]byte(src), "cmd-artifact-argv0.toml"); err != nil {
		t.Errorf("`{artifact}` at argv0 refused: %v — the argv0 rule is on the "+
			"`{tag.<key>}` family only", err)
	}
}

// --- C1.4: registration and precedence -----------------------------------

// REQ-69: "registration: appended to `table.Categories()` after 0025:C5's
// six, in the order above; typed `Cat…` constants beside the others; the
// wire strings are the contract, the identifiers are not; the list's size
// is not a contract (0025:C5)."
// REQ-76: "The assertion is RELATIVE order, never a tail position or a
// count — the shipped `TestReq77` pins 0025:C5's six as `Categories()`'s
// tail and this append necessarily rewrites it (deviations D1)."
// BOUNDARY
//
// S6. The assertion is deliberately RELATIVE: `Categories()` is
// append-only and its size is explicitly not a contract, so a golden
// keyed on index or length is a scheduled false positive that fails for
// the correct reason at the next append. What is pinned is that the six
// 0025 categories appear contiguously in clause order and that the six
// 0028 categories follow them, also contiguously and in C1.4's order.
func TestReq69And76_TheSixEditCategoriesFollow0025sSixInClauseOrder(t *testing.T) {
	c25 := []table.Category{
		table.CatCommandAndPathConflict,
		table.CatCommandEmpty,
		table.CatCommandUnknownPlaceholder,
		table.CatCommandShellInterpreter,
		table.CatCommandOutputShape,
		table.CatCommandEnvConflict,
	}
	c28 := []table.Category{
		table.CatEditCarrierConflict,
		table.CatEditKeyMismatch,
		table.CatEditAnchorInvalid,
		table.CatEditTemplateInvalid,
		table.CatEditClearInvalid,
		table.CatEditTagArgv0,
	}

	all := table.Categories()
	start := slices.Index(all, c25[0])
	if start < 0 {
		t.Fatalf("Categories() does not carry %q at all", c25[0])
	}

	want := slices.Concat(c25, c28)
	if start+len(want) > len(all) {
		t.Fatalf("Categories() has %d members from index %d; cannot carry the "+
			"twelve contiguous categories", len(all)-start, start)
	}
	got := all[start : start+len(want)]
	if !slices.Equal(got, want) {
		t.Errorf("Categories()[%d:%d] = %#v; want 0025:C5's six in clause order "+
			"followed by C1.4's six in clause order %#v", start, start+len(want), got, want)
	}
}

// REQ-69 (wire strings half): "the wire strings are the contract, the
// identifiers are not"
// REQ-73: "The first five categories fire on the `edit` table;
// `edit_tag_argv0` fires on an entry's `command` argv, so an entry
// carrying `command` reaches it while an `edit` entry never does — the two
// sets are disjoint by carrier and never race"
// req-list A-12: the six C1.4 wire strings and the six apply-time Detail
// reason tokens are DISJOINT sets.
// BOUNDARY
//
// The oracle reads the REGISTERED set by wire STRING, not by Go constant:
// C1.4 fixes the strings as the contract and leaves the identifiers
// unconstrained, so comparing constants to literals would pin the half the
// clause does not fix and pass while a category was invisible.
//
// The second half is A-12's disjointness: registering an apply-time reason
// token in `Categories()` would tell a consumer that a stale anchor is a
// LINT defect, which C1.4 `what lint does NOT prove:` denies outright.
func TestReq69And73_TheSixWireStringsRegisterAndTheApplyTokensDoNot(t *testing.T) {
	registered := map[string]bool{}
	for _, c := range table.Categories() {
		registered[string(c)] = true
	}

	for _, want := range []string{
		"edit_carrier_conflict",
		"edit_key_mismatch",
		"edit_anchor_invalid",
		"edit_template_invalid",
		"edit_clear_invalid",
		"edit_tag_argv0",
	} {
		if !registered[want] {
			t.Errorf("Categories() does not carry the wire string %q; a category "+
				"that refuses correctly but never registers is invisible to every "+
				"consumer enumerating the list", want)
		}
	}

	for _, token := range []string{
		"edit_anchor_unmatched",
		"edit_anchor_ambiguous",
		"edit_anchor_collision",
		"edit_anchor_unstable",
		"edit_value_multiline",
		"edit_clear_undeclared",
	} {
		if registered[token] {
			t.Errorf("Categories() carries %q; that is a C1.3/C1.5 APPLY-time "+
				"Detail reason token and registering it claims lint decides it — "+
				"the two sets are disjoint", token)
		}
	}
}

// REQ-71: "precedence: within one entry, fail-fast in the order above,
// evaluated after 0025:C5's clauses 1–6"
// REQ-116 (S28): "One entry carrying TWO simultaneous load-time defects ...
// the earlier category in C1.4's registration order is the one reported,
// and only it"
// ADVERSARIAL
//
// Fail-fast is observable ONLY through a multi-defect input: with one
// defect per fixture every evaluation order passes, which is why C1.4
// `precedence:` earns its own scenario. Each row asserts the EARLIER
// category and, by construction, that the later one did not win.
func TestReq71And116_TheEarlierCategoryInRegistrationOrderIsTheOneReported(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want table.Category
	}{
		{
			// `edit` beside `path` AND an uncompilable anchor.
			name: "carrier_conflict_beats_anchor_invalid",
			body: `[write.record]
role = "record"
path = "flows/record.toml"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "(unclosed"
replace = "{status}"`,
			want: table.CatEditCarrierConflict,
		},
		{
			// A `keys` mismatch AND a malformed `replace`.
			name: "key_mismatch_beats_template_invalid",
			body: `[write.record]
role = "record"
keys = ["status", "owner"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "$1"`,
			want: table.CatEditKeyMismatch,
		},
		{
			// An uncompilable anchor AND a malformed `replace` on the
			// same rule: anchor is the earlier category.
			name: "anchor_invalid_beats_template_invalid",
			body: `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "(unclosed"
replace = "{nope}"`,
			want: table.CatEditAnchorInvalid,
		},
		{
			// A malformed `replace` AND an out-of-set `clear`.
			name: "template_invalid_beats_clear_invalid",
			body: `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (.*)$"
replace = "{nope}"
clear   = "cell"`,
			want: table.CatEditClearInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := editSwap(t, editCarrierModel, editWriteBlock, tc.body)
			got := editCategoryOf(t, src, "edit-multi-defect.toml")
			if tc.name == "template_invalid_beats_clear_invalid" {
				// Registration order puts edit_template_invalid before
				// edit_clear_invalid, so the template defect wins.
				if got != table.CatEditTemplateInvalid {
					t.Errorf("category = %q; want %q — C1.4's registration order "+
						"puts the template category before the clear category",
						got, table.CatEditTemplateInvalid)
				}
				return
			}
			if got != tc.want {
				t.Errorf("category = %q; want %q — the earlier category in C1.4's "+
					"registration order is the one reported, and only it", got, tc.want)
			}
		})
	}
}

// REQ-72: "across entries and tables, 0025:C5's rules apply unchanged —
// including across the sibling `edit.<key>` tables of ONE entry, which are
// map-ranged ... which of two equally-defective tables is reported is
// unspecified and no test may assert it."
// REQ-116 (S28, sibling half): "TWO `edit.<key>` tables in one entry each
// carrying a defect of the SAME category. The assertion is that exactly
// ONE category is reported and it is that category — NOT which table is
// named."
// NEGATIVE REQ.
// ADVERSARIAL
//
// Go's map iteration is randomised, so a test pinning the winning table
// would assert a guarantee the contract explicitly declines to make. What
// IS asserted is stability of the CATEGORY across repeated loads: the
// randomised order must not leak into the reported category.
func TestReq72And116_SiblingTablesReportOneStableCategoryAndNoNamedWinner(t *testing.T) {
	src := editSwap(t, editCarrierModel, editWriteBlock, `[write.record]
role = "record"
keys = ["status", "owner"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "(unclosed"
replace = "{status}"

[write.record.edit.owner]
anchor  = "(also-unclosed"
replace = "{owner}"`)

	// Repeat so a category that varied with map iteration order shows up
	// rather than passing on one lucky ordering.
	for i := range 32 {
		got := editCategoryOf(t, src, "edit-sibling-defects.toml")
		if got != table.CatEditAnchorInvalid {
			t.Fatalf("load %d: category = %q; want %q on every iteration — the "+
				"WINNING TABLE is unspecified but the CATEGORY is not",
				i, got, table.CatEditAnchorInvalid)
		}
	}
}

// REQ-74: "what lint proves: the carrier is one, every anchor compiles,
// every template parses, every placeholder is closed over the entry's keys
// and the model's declared tags, every `keys` member has one rule, and no
// `{tag.<key>}` sits at argv0."
// HAPPY PATH
//
// The composite: a model satisfying all six obligations loads clean. This
// is the arm every negative test above is measured against — without it a
// loader that refused EVERYTHING would pass the whole negative suite.
func TestReq74_AModelMeetingEverySixLintObligationLoadsClean(t *testing.T) {
	src := editSwap(t, editCarrierModel, editWriteBlock, `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^- \\*\\*Status\\*\\*: (\\w+)(.*)$"
replace = "- **Status**: {status}${2}"
clear   = "line"`)

	m := mustLoadEdit(t, src, "edit-lint-clean.toml")
	if len(m.Writers["record"].Edit) != 1 {
		t.Errorf("the clean model carries %d edit rules; want exactly one for `status`",
			len(m.Writers["record"].Edit))
	}
}

// REQ-75: "What lint does NOT prove: that an anchor matches exactly one
// line of a particular file, or that a bound tag value is not flag-shaped
// — those are C1.3's and C1.6's apply-time refusals, by design (a model is
// linted without its artifacts and without the invocation's bindings)"
// NEGATIVE REQ.
// DOMAIN EDGE
//
// A model whose anchor CANNOT match anything sensible — and whose declared
// `{tag.nnnn}` will be bound to whatever the caller supplies — must still
// LOAD. A loader that reached for the filesystem, or that rejected a
// pattern for being over-broad, would have made lint depend on artifacts
// the clause says it never sees.
func TestReq75_LintDoesNotConsultArtifactsOrBindings(t *testing.T) {
	src := editSwap(t, editCarrierModel, editWriteBlock, `[write.record]
role = "record"
keys = ["status"]
timeout = "2s"
read_back = true

[write.record.edit.status]
anchor  = "^.*$"
replace = "{status}"`)

	// `^.*$` matches every line of every file — a guaranteed apply-time
	// `edit_anchor_ambiguous`. Lint must not anticipate it.
	mustLoadEdit(t, src, "edit-overbroad-anchor.toml")
}

// REQ-101: "Admit `edit` on write entries — `table.Accessor.Edit`, the
// `edit.<key>` tables, C1.4's categories in `carrierDefect`'s order,
// template and anchor parsed at load (C1, C1.2, C1.4); dump/normalize
// round-trip carries it (A6)."
// REQ-19: "parse once: both templates are parsed at LOAD into segments.
// At apply each segment emits bytes; captured text and substituted values
// are never re-scanned for `${…}` or `{…}`."
// BOUNDARY
//
// Two obligations in one fixture.
//
// A6's dump carriage: an `edit`-carried model whose review dump renders
// identically to a carrier-less one hides the writer from the artifact a
// reviewer reads — the surface would not be the model that runs.
//
// `parse once:` at LOAD: the templates are parsed into segments when the
// model is loaded, so a template defect is a LOAD refusal (asserted
// throughout this file) and the parse is not deferred to first apply. The
// witness available at this seam is that the defect surfaces from
// `table.Load` with no artifact in sight.
func TestReq19And101_TheEditCarrierIsCarriedOnTheDumpSurface(t *testing.T) {
	m := mustLoadEdit(t, editCarrierModel, "edit-carrier.toml")

	dumped := table.Dump(m)
	if !strings.Contains(dumped, "edit") {
		t.Errorf("the review dump of an `edit`-carried model does not name the "+
			"carrier anywhere:\n%s", dumped)
	}

	t.Run("templates_parse_at_load_not_at_apply", func(t *testing.T) {
		// No artifact, no binding, no apply — a defective template is
		// refused by Load itself, which is what `parse once:` requires.
		body := strings.Replace(editWriteBlock,
			`replace = "- **Status**: {status}"`,
			`replace = "- **Status**: ${99}"`, 1)
		src := editSwap(t, editCarrierModel, editWriteBlock, body)
		if got := editCategoryOf(t, src, "edit-parse-at-load.toml"); got != table.CatEditTemplateInvalid {
			t.Errorf("category = %q; want %q from Load itself — the templates are "+
				"parsed at LOAD, never deferred to first apply",
				got, table.CatEditTemplateInvalid)
		}
	})
}
