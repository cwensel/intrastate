package table_test

// Kata wchf — the LOAD half of a `command` read entry's per-key `select`
// tables: the declaration decodes onto the model, and its defects carry
// the three `command_select_*` categories. The read half — pointer
// resolution, prefix projection, declared absence — lives in
// `internal/cli/cmdbind/select_wchf_test.go`, because lint cannot see a
// tool's output and proves the declaration only.
//
// Each category is asserted twice, by the category a refusal carries AND
// by its membership in the registered set, for the reason RDR 0025's and
// 0028's suites give: an unregistered category refuses correctly while
// staying invisible to every consumer that enumerates the set.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// selectModel is a `kata show`-shaped consumer: an owned enum read out of
// a label array by prefix, and an observed scalar read out of a nested
// object whose null and missing forms both mean "unset".
const selectModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "kataflow"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.id]
provenance = "observed"
kind = "scalar"

[tags.owner]
provenance = "observed"
kind = "scalar"

[tags.areas]
provenance = "observed"
kind = "set"
elements = ["cli", "table"]

[tags.lifecycle]
provenance = "owned"
kind = "enum"
domain = ["queued", "shipping"]
single_valued = true
required = true

[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "owner"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"

[read.kata.select.owner]
pointer = "/issue/owner"
absent = ["null", "missing"]

[write.kata]
role = "kata"
command = ["tools/kata-label", "{tag.id}"]
keys = ["lifecycle"]
timeout = "2s"
read_back = true

[initial]
lifecycle = "queued"

[context.done]
[context.done.match.lifecycle]
eq = "shipping"

[[rule]]
id = "ship"
[rule.match.lifecycle]
eq = "queued"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "shipping"
`

// selectReadBlock is the read entry every mutant below rewrites.
const selectReadBlock = `[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "owner"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"

[read.kata.select.owner]
pointer = "/issue/owner"
absent = ["null", "missing"]
`

func TestSelect_TablesDecodeOntoTheReadEntry(t *testing.T) {
	m, err := table.Load([]byte(selectModel), "select.toml")
	if err != nil {
		t.Fatalf("the select model refused: %v", err)
	}
	got := m.Readers["kata"].Select
	want := map[string]table.SelectRule{
		"lifecycle": {Pointer: "/labels", Element: "/label", Prefix: "lifecycle:"},
		"owner":     {Pointer: "/issue/owner", AbsentNull: true, AbsentMissing: true},
	}
	if len(got) != len(want) {
		t.Fatalf("Select = %#v; want %#v", got, want)
	}
	for k, w := range want {
		if got[k] != w {
			t.Errorf("Select[%q] = %#v; want %#v", k, got[k], w)
		}
	}
	if m.Writers["kata"].Select != nil {
		t.Error("an entry declaring no `select` carries a non-nil Select; the " +
			"flat-object read must stay the zero value")
	}
}

func TestSelect_DefectsCarryTheSelectCategories(t *testing.T) {
	cases := []struct {
		name  string
		block string
		want  table.Category
	}{
		{
			name: "select on a path entry",
			block: `[read.kata]
role = "kata"
path = "kata.json"
keys = ["lifecycle", "owner"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/lifecycle"

[read.kata.select.owner]
pointer = "/owner"
`,
			want: table.CatCommandSelectPlacement,
		},
		{
			name: "select on a raw entry",
			block: `[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}"]
output = "raw"
keys = ["lifecycle"]
timeout = "2s"

[read.owner]
role = "kata"
command = ["kata", "owner", "{tag.id}"]
output = "raw"
keys = ["owner"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = ""
`,
			want: table.CatCommandSelectPlacement,
		},
		{
			name: "a key with no select table",
			block: `[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "owner"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"
`,
			want: table.CatCommandSelectKeyMismatch,
		},
		{
			name: "a select table for a key not in keys",
			block: selectReadBlock + `
[read.kata.select.id]
pointer = "/issue/id"
`,
			want: table.CatCommandSelectKeyMismatch,
		},
		{
			name: "a bare select table",
			block: `[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "owner"]
timeout = "2s"

[read.kata.select]
`,
			want: table.CatCommandSelectKeyMismatch,
		},
		{
			name:  "a pointer without a leading slash",
			block: swapEntry(t, selectReadBlock, `pointer = "/issue/owner"`, `pointer = "issue/owner"`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "a pointer carrying a bad escape",
			block: swapEntry(t, selectReadBlock, `pointer = "/issue/owner"`, `pointer = "/issue/~2"`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "no pointer",
			block: swapEntry(t, selectReadBlock, `pointer = "/issue/owner"`, ``),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "an element that does not parse",
			block: swapEntry(t, selectReadBlock, `element = "/label"`, `element = "label"`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "an element without a prefix",
			block: swapEntry(t, selectReadBlock, `prefix = "lifecycle:"`, ``),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "an empty prefix",
			block: swapEntry(t, selectReadBlock, `prefix = "lifecycle:"`, `prefix = ""`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "an absent member outside the closed set",
			block: swapEntry(t, selectReadBlock, `absent = ["null", "missing"]`, `absent = ["empty"]`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name:  "an absent member declared twice",
			block: swapEntry(t, selectReadBlock, `absent = ["null", "missing"]`, `absent = ["null", "null"]`),
			want:  table.CatCommandSelectInvalid,
		},
		{
			name: "a prefix on a set-kind tag",
			block: selectReadBlock + `
[read.areas]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["areas"]
timeout = "2s"

[read.areas.select.areas]
pointer = "/labels"
element = "/label"
prefix = "area:"
`,
			want: table.CatCommandSelectInvalid,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, selectModel, selectReadBlock, tc.block)
			if got := loadCategoryOf(t, src, "select-defect.toml"); got != tc.want {
				t.Errorf("category = %q; want %q", got, tc.want)
			}
		})
	}

	t.Run("select on a write entry", func(t *testing.T) {
		src := swapEntry(t, selectModel, `read_back = true
`, `read_back = true

[write.kata.select.lifecycle]
pointer = "/labels"
`)
		if got := loadCategoryOf(t, src, "select-write.toml"); got != table.CatCommandSelectPlacement {
			t.Errorf("category = %q; want %q", got, table.CatCommandSelectPlacement)
		}
	})
}

func TestSelect_CategoriesAreRegistered(t *testing.T) {
	for _, c := range []table.Category{
		table.CatCommandSelectPlacement,
		table.CatCommandSelectKeyMismatch,
		table.CatCommandSelectInvalid,
	} {
		if !slices.Contains(table.Categories(), c) {
			t.Errorf("%q refuses at its call site but is not in Categories()", c)
		}
	}
}

func TestSelect_ParseJSONPointerFollowsRFC6901(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"/", []string{""}},
		{"/issue/owner", []string{"issue", "owner"}},
		{"/a~1b", []string{"a/b"}},
		{"/m~0n", []string{"m~n"}},
		// `~01` decodes to `~1`, never to `/`: `~1` is replaced FIRST.
		{"/~01", []string{"~1"}},
		{"/labels/0", []string{"labels", "0"}},
	}
	for _, tc := range cases {
		got, err := table.ParseJSONPointer(tc.in)
		if err != nil {
			t.Errorf("ParseJSONPointer(%q) refused: %v", tc.in, err)
			continue
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("ParseJSONPointer(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
	for _, bad := range []string{"issue", "/a~", "/a~2", "#/a"} {
		if _, err := table.ParseJSONPointer(bad); err == nil {
			t.Errorf("ParseJSONPointer(%q) parsed; want a refusal", bad)
		}
	}
}
