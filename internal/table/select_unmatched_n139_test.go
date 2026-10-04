package table_test

// Kata n139 — the LOAD half of `unmatched` on a pointer select: beside a
// declared `absent` set, a select without a `prefix` may name the value
// an established-absent field reads. The read half is
// `internal/cli/cmdbind/select_unmatched_n139_test.go`.
//
// Every other load rule is 9xjf's unchanged: never empty, inside the tag's
// domain, only for a tag declared `required = true`, never `<clear>`.

import (
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// pointerUnmatchedLifecycle reads the lifecycle off a nested field the
// tool omits or nulls when the issue is queued.
const pointerUnmatchedLifecycle = `[read.kata.select.lifecycle]
pointer = "/issue/lifecycle"
absent = ["null", "missing"]
unmatched = "queued"
`

// pointerUnmatchedModel is selectModel with the lifecycle read through a
// pointer select that declares `absent` and `unmatched`.
func pointerUnmatchedModel(t *testing.T) string {
	t.Helper()

	return swapEntry(t, unmatchedModel(t), unmatchedLifecycle, pointerUnmatchedLifecycle)
}

func TestSelectUnmatchedN139_PointerWithAbsentDecodesOntoTheRule(t *testing.T) {
	m, err := table.Load([]byte(pointerUnmatchedModel(t)), "pointer-unmatched.toml")
	if err != nil {
		t.Fatalf("a pointer rule declaring `absent` and an in-domain "+
			"`unmatched` for a required enum refused: %v", err)
	}
	want := table.SelectRule{Pointer: "/issue/lifecycle",
		AbsentNull: true, AbsentMissing: true, Unmatched: "queued"}
	if got := m.Readers["kata"].Select["lifecycle"]; got != want {
		t.Errorf("Select[lifecycle] = %#v; want %#v", got, want)
	}
}

func TestSelectUnmatchedN139_PointerDefectsAreSelectInvalid(t *testing.T) {
	cases := []struct {
		name, old, new string
	}{
		{
			name: "no prefix and no absent",
			old:  "absent = [\"null\", \"missing\"]\nunmatched = \"queued\"\n",
			new:  "unmatched = \"queued\"\n",
		},
		{
			name: "an empty absent",
			old:  "absent = [\"null\", \"missing\"]\nunmatched = \"queued\"\n",
			new:  "absent = []\nunmatched = \"queued\"\n",
		},
		{
			name: "empty",
			old:  `unmatched = "queued"`,
			new:  `unmatched = ""`,
		},
		{
			name: "outside the domain",
			old:  `unmatched = "queued"`,
			new:  `unmatched = "open"`,
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
			src := swapEntry(t, pointerUnmatchedModel(t), tc.old, tc.new)
			_, err := table.Load([]byte(src), "pointer-unmatched-defect.toml")
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

func TestSelectUnmatchedN139_PointerClearSentinelIsReserved(t *testing.T) {
	src := swapEntry(t, pointerUnmatchedModel(t), `unmatched = "queued"`, `unmatched = "<clear>"`)
	_, err := table.Load([]byte(src), "pointer-unmatched-clear.toml")
	if err == nil {
		t.Fatal("loaded clean; want reserved_tag_value")
	}
	f := cmdFailureOf(t, err)
	if f.Category != table.CatReservedTagValue {
		t.Errorf("category = %q; want %q", f.Category, table.CatReservedTagValue)
	}
}
