package table_test

// Kata 9xjf — the LOAD half of a `select` table's `unmatched` key: the
// value a prefix projection reads when no member carries the prefix. The
// read half is `internal/cli/cmdbind/select_unmatched_9xjf_test.go`.
//
// The key is narrow by construction: only beside a `prefix`, never empty,
// inside the tag's domain, and only for a tag declared `required = true` —
// a reader that can no longer report the key absent must not leave an
// "optional" tag owing an absent arm nothing can reach.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// unmatchedLifecycle is the base model's lifecycle selector with the
// no-match value declared.
const unmatchedLifecycle = `[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"
unmatched = "queued"
`

// unmatchedModel is selectModel with `unmatched` declared on the prefix
// projection of the required owned enum.
func unmatchedModel(t *testing.T) string {
	t.Helper()

	return swapEntry(t, selectModel, `[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"
`, unmatchedLifecycle)
}

func TestSelectUnmatched_DecodesOntoTheRule(t *testing.T) {
	m, err := table.Load([]byte(unmatchedModel(t)), "unmatched.toml")
	if err != nil {
		t.Fatalf("a prefix rule declaring an in-domain `unmatched` for a "+
			"required enum refused: %v", err)
	}
	want := table.SelectRule{Pointer: "/labels", Element: "/label",
		Prefix: "lifecycle:", Unmatched: "queued"}
	if got := m.Readers["kata"].Select["lifecycle"]; got != want {
		t.Errorf("Select[lifecycle] = %#v; want %#v", got, want)
	}
	if got := m.Readers["kata"].Select["owner"].Unmatched; got != "" {
		t.Errorf("owner declares no `unmatched` but carries %q", got)
	}
}

func TestSelectUnmatched_DefectsAreSelectInvalid(t *testing.T) {
	cases := []struct {
		name, old, new string
	}{
		{
			name: "without a prefix",
			old:  "element = \"/label\"\nprefix = \"lifecycle:\"\n",
			new:  "",
		},
		{
			name: "empty",
			old:  `unmatched = "queued"`,
			new:  `unmatched = ""`,
		},
		{
			name: "outside the domain",
			old:  `unmatched = "queued"`,
			new:  `unmatched = "filed"`,
		},
		{
			name: "on a tag not declared required",
			old: `domain = ["queued", "shipping"]
single_valued = true
required = true`,
			new: `domain = ["queued", "shipping"]
single_valued = true`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, unmatchedModel(t), tc.old, tc.new)
			_, err := table.Load([]byte(src), "unmatched-defect.toml")
			if err == nil {
				t.Fatal("loaded clean; want command_select_invalid")
			}
			f := cmdFailureOf(t, err)
			if f.Category != table.CatCommandSelectInvalid {
				t.Errorf("category = %q; want %q", f.Category, table.CatCommandSelectInvalid)
			}
			if !strings.Contains(f.Detail, "`unmatched`") {
				t.Errorf("Detail %q does not name `unmatched`; the refusal "+
					"must come from the key under test", f.Detail)
			}
		})
	}
}

// TestSelectUnmatched_ClearSentinelIsReserved pins that `unmatched` is an
// authored tag value, so the reserved `<clear>` sentinel is refused there as
// it is at every other authoring site. A projection reading the sentinel on
// zero matches would be classified unreadable, defeating the declared
// fallback. The scalar case is the one no domain check would catch.
func TestSelectUnmatched_ClearSentinelIsReserved(t *testing.T) {
	scalar := swapEntry(t, unmatchedModel(t), `[tags.owner]
provenance = "observed"
kind = "scalar"
`, `[tags.owner]
provenance = "observed"
kind = "scalar"
required = true
`)
	scalar = swapEntry(t, scalar, `[read.kata.select.owner]
pointer = "/issue/owner"
absent = ["null", "missing"]
`, `[read.kata.select.owner]
pointer = "/labels"
element = "/label"
prefix = "owner:"
unmatched = "nobody"
`)
	if _, err := table.Load([]byte(scalar), "unmatched-scalar.toml"); err != nil {
		t.Fatalf("the scalar base refused before the sentinel was authored: %v", err)
	}

	cases := []struct {
		name, src, old string
	}{
		{name: "on a required scalar", src: scalar, old: `unmatched = "nobody"`},
		{name: "on an enum", src: unmatchedModel(t), old: `unmatched = "queued"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, tc.src, tc.old, `unmatched = "<clear>"`)
			_, err := table.Load([]byte(src), "unmatched-clear.toml")
			if err == nil {
				t.Fatal("loaded clean; want reserved_tag_value")
			}
			f := cmdFailureOf(t, err)
			if f.Category != table.CatReservedTagValue {
				t.Errorf("category = %q; want %q", f.Category, table.CatReservedTagValue)
			}
			if !strings.Contains(f.Detail, "`select.") ||
				!strings.Contains(f.Detail, table.ClearSentinel) {
				t.Errorf("Detail %q does not name the select table and the "+
					"reserved value", f.Detail)
			}
		})
	}
}
