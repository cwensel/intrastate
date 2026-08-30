package table

// RDR 0024 — `tableHeader`'s canonicalization contract, asserted from
// inside the package because the helper is unexported.
//
// `emitRuleLine` builds its `[[rule]]` census by comparing
// `tableHeader(line)` against the single canonical string `"[[rule]]"`, and
// uses the census index as the decoder's rule ordinal. So every spelling
// TOML admits for the array table `rule` MUST canonicalize to that one
// string: a spelling the census misses drops a block, shifts every later
// ordinal onto the wrong rule, and produces the "pointing at the wrong text
// costs more than pointing at no text" failure `emitRuleLine`'s own doc
// comment disclaims. That is `0024:MVV` step 2's "each carrying the
// offending block's SOURCE LINE, not `:1`" (REQ-30) and REQ-32's rule-side
// keying, both read through the census the locator is built on.

import "testing"

// REQ-30: "**All three categories MUST carry a source line**, stamped
// through `internal/table/category.go::atLine` the way `loadTags` stamps a
// tag refusal."
// REQ-32: "a rule-side defect (`unknown_emit_key`,
// `emit_value_out_of_domain`) keys on the offending RULE ID".
//
// TOML's key grammar admits a bare key, a basic-string key and a
// literal-string key interchangeably: `[[rule]]`, `[["rule"]]` and
// `[['rule']]` all name the SAME array table and all decode into
// `l.doc.Rule`. Whitespace inside the brackets is equally free. Every one
// of them must report the canonical `[[rule]]`.
// INPUT EDGE
func TestTableHeaderCanonicalizesQuotedAndSpacedRuleHeaders(t *testing.T) {
	for _, tc := range []struct{ raw, want string }{
		{"[[rule]]", "[[rule]]"},
		{"[[ rule ]]", "[[rule]]"},
		{`[["rule"]]`, "[[rule]]"},
		{`[['rule']]`, "[[rule]]"},
		{`[[ "rule" ]]`, "[[rule]]"},
		{`[[ 'rule' ]]`, "[[rule]]"},
		{`[["rule"]]  # trailing comment`, "[[rule]]"},
		{`["rule"]`, "[rule]"},
		{`[ 'model' ]`, "[model]"},
	} {
		if got := tableHeader(tc.raw); got != tc.want {
			t.Errorf("tableHeader(%q) = %q, want %q; the spelling names the "+
				"same TOML table, so it must join the caller's census rather "+
				"than shifting every later ordinal", tc.raw, got, tc.want)
		}
	}
}

// Negative control on the same canonicalization. Stripping quotes must not
// widen what counts as a TOP-LEVEL header: a dotted key still belongs to
// the block it sits in and must not close it, and quoting the key does not
// change that. The dot test therefore has to run AFTER the quotes come off,
// or a quoted dotted key is misjudged.
//
// The premortem this guards: over-canonicalizing would fold a genuinely
// distinct table name into `rule` and add a PHANTOM census entry, shifting
// ordinals the other way. Only a matched outer quote pair on an otherwise
// bare key is stripped.
// ADVERSARIAL
func TestTableHeaderRejectsDottedAndUnbalancedQuotedHeaders(t *testing.T) {
	for _, raw := range []string{
		`[["rule.other"]]`,
		`[['rule.other']]`,
		`["rule.emit"]`,
		`[rule.emit]`,
		`[[ "rule.other" ]]`,
	} {
		if got := tableHeader(raw); got != "" {
			t.Errorf("tableHeader(%q) = %q, want \"\"; a dotted key names a "+
				"sub-table and must not be read as a top-level header, however "+
				"it is quoted", raw, got)
		}
	}

	// A quote that is not a MATCHED outer pair is not a quoted key. It is
	// carried through verbatim so that a table genuinely named `"rule` can
	// never be folded into the `rule` census.
	for _, tc := range []struct{ raw, want string }{
		{`[["rule]]`, `[["rule]]`},
		{`[[rule"]]`, `[[rule"]]`},
		{`[["rule']]`, `[["rule']]`},
		{`[[""]]`, `[[]]`},
	} {
		if got := tableHeader(tc.raw); got != tc.want {
			t.Errorf("tableHeader(%q) = %q, want %q; only a MATCHED outer "+
				"quote pair is a quoted key, and folding anything else would "+
				"add a phantom census entry", tc.raw, got, tc.want)
		}
	}
}
