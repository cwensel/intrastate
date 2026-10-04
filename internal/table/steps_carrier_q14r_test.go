package table_test

// Kata q14r — the LOAD half of the `steps` write carrier: per key, per
// planned value, literal argv vectors. The tables decode onto the model,
// the carrier joins the exactly-one rule, and its defects carry the three
// `steps_*` categories — or, for a step vector's own argv, the existing
// `command` clause categories the vector reuses. The apply half — arm
// selection, fail-fast, the applied sense — is
// `internal/cli/steps_q14r_test.go`.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// stepsModel is the consumer's lifecycle: an owned enum held as one
// `lifecycle:<value>` label, moved by removing the old label and adding
// the new one, with a claim in front of `resolving` and an unassign in
// front of `closed`. `rewind` clears the key, so every value it can hold
// owes a clear arm.
const stepsModel = `outcomes = ["advance", "rewind"]
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

[tags.session]
provenance = "observed"
kind = "scalar"

[tags.lifecycle]
provenance = "owned"
kind = "enum"
domain = ["resolving", "shipping", "closed"]
single_valued = true

[tags.attempts]
provenance = "owned"
kind = "int"
min = 0
max = 3

[tags.memo]
provenance = "owned"
kind = "scalar"

[read.memo]
role = "memo"
path = "memo"
keys = ["memo"]
timeout = "2s"

[write.memo]
role = "memo"
path = "memo"
keys = ["memo"]
timeout = "2s"
read_back = true

[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "attempts"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"

[read.kata.select.attempts]
pointer = "/attempts"

` + stepsWriteBlock + `
[initial]
lifecycle = "resolving"

[context.done]
[context.done.match.lifecycle]
eq = "closed"

[[rule]]
id = "ship"
[rule.match.lifecycle]
eq = "resolving"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "shipping"

[[rule]]
id = "close"
[rule.match.lifecycle]
eq = "shipping"
[rule.match.recognized]
eq = "advance"
[rule.write]
lifecycle = "closed"

[[rule]]
id = "rewind"
clear = ["lifecycle"]
[rule.match.lifecycle]
eq = "shipping"
[rule.match.recognized]
eq = "rewind"
[rule.write]
memo = "rewound"
`

// stepsWriteBlock is the write entry every mutant below rewrites.
const stepsWriteBlock = `[write.kata]
role = "kata"
keys = ["lifecycle"]
timeout = "2s"
read_back = true

[write.kata.steps.lifecycle.set]
resolving = [["kata", "claim", "{tag.id}", "--as", "{tag.session}"], ["kata", "label", "add", "{tag.id}", "lifecycle:resolving"]]
shipping = [["kata", "label", "add", "{tag.id}", "lifecycle:shipping"]]
closed = [["kata", "unassign", "{tag.id}"], ["kata", "label", "add", "{tag.id}", "lifecycle:closed"]]

[write.kata.steps.lifecycle.clear]
resolving = [["kata", "label", "rm", "{tag.id}", "lifecycle:resolving"]]
shipping = [["kata", "label", "rm", "{tag.id}", "lifecycle:shipping"]]
closed = [["kata", "label", "rm", "{tag.id}", "lifecycle:closed"]]
`

func TestSteps_TablesDecodeOntoTheWriteEntry(t *testing.T) {
	m, err := table.Load([]byte(stepsModel), "steps.toml")
	if err != nil {
		t.Fatalf("the steps model refused: %v", err)
	}
	w := m.Writers["kata"]
	if w.Command != nil || w.Path != "" || w.Edit != nil {
		t.Fatalf("the steps entry carries another carrier: %#v", w)
	}
	rule, ok := w.Steps["lifecycle"]
	if !ok {
		t.Fatalf("Steps = %#v; want a `lifecycle` table", w.Steps)
	}
	wantSet := [][]string{
		{"kata", "claim", "{tag.id}", "--as", "{tag.session}"},
		{"kata", "label", "add", "{tag.id}", "lifecycle:resolving"},
	}
	if got := rule.Set["resolving"]; !slices.EqualFunc(got, wantSet, slices.Equal) {
		t.Errorf("set.resolving = %q; want %q in array order", got, wantSet)
	}
	if got := rule.Clear["shipping"]; len(got) != 1 ||
		!slices.Equal(got[0], []string{"kata", "label", "rm", "{tag.id}", "lifecycle:shipping"}) {
		t.Errorf("clear.shipping = %q; want the one label-remove step", got)
	}
	if m.Readers["kata"].Steps != nil {
		t.Error("an entry declaring no `steps` carries a non-nil Steps")
	}
}

func TestSteps_ABoolKeyIsAdmitted(t *testing.T) {
	src := swapEntry(t, stepsModel, `[tags.attempts]`, `[tags.blocked]
provenance = "owned"
kind = "bool"

[read.blocked]
role = "kata"
command = ["kata", "blocked", "{tag.id}"]
keys = ["blocked"]
timeout = "2s"

[write.blocked]
role = "kata"
keys = ["blocked"]
timeout = "2s"
read_back = true

[write.blocked.steps.blocked.set]
true = [["kata", "label", "add", "{tag.id}", "blocked"]]
false = [["kata", "label", "rm", "{tag.id}", "blocked"]]

[tags.attempts]`)
	if _, err := table.Load([]byte(src), "steps-bool.toml"); err != nil {
		t.Fatalf("a bool key with a set arm per literal refused: %v", err)
	}
}

func TestSteps_DefectsCarryTheirCategories(t *testing.T) {
	cases := []struct {
		name string
		old  string
		new  string
		want table.Category
	}{
		{
			name: "steps beside a command",
			old:  "keys = [\"lifecycle\"]\ntimeout = \"2s\"\nread_back = true\n",
			new:  "command = [\"kata\", \"set\"]\nkeys = [\"lifecycle\"]\ntimeout = \"2s\"\nread_back = true\n",
			want: table.CatStepsCarrierConflict,
		},
		{
			name: "steps beside a path",
			old:  "keys = [\"lifecycle\"]\ntimeout = \"2s\"\nread_back = true\n",
			new:  "path = \"kata.json\"\nkeys = [\"lifecycle\"]\ntimeout = \"2s\"\nread_back = true\n",
			want: table.CatStepsCarrierConflict,
		},
		{
			name: "steps beside an edit",
			old:  "[write.kata.steps.lifecycle.set]",
			new:  "[write.kata.edit.lifecycle]\nanchor = \"^x$\"\nreplace = \"x\"\n\n[write.kata.steps.lifecycle.set]",
			want: table.CatStepsCarrierConflict,
		},
		{
			name: "a bare steps table",
			old:  stepsWriteBlock,
			new:  "[write.kata]\nrole = \"kata\"\nkeys = [\"lifecycle\"]\ntimeout = \"2s\"\nread_back = true\n\n[write.kata.steps]\n",
			want: table.CatStepsKeyMismatch,
		},
		{
			name: "a steps table for a key not in keys",
			old:  "[write.kata.steps.lifecycle.clear]",
			new:  "[write.kata.steps.attempts.set]\n0 = [[\"kata\", \"reset\"]]\n\n[write.kata.steps.lifecycle.clear]",
			want: table.CatStepsKeyMismatch,
		},
		{
			name: "no set arm for a domain value",
			old:  "shipping = [[\"kata\", \"label\", \"add\", \"{tag.id}\", \"lifecycle:shipping\"]]\n",
			new:  "",
			want: table.CatStepsTableInvalid,
		},
		{
			name: "a set arm for a value outside the domain",
			old:  "[write.kata.steps.lifecycle.clear]",
			new:  "[write.kata.steps.lifecycle.clear]\nrefining = [[\"kata\", \"label\", \"rm\", \"{tag.id}\", \"lifecycle:refining\"]]",
			want: table.CatStepsTableInvalid,
		},
		{
			name: "an empty arm",
			old:  "shipping = [[\"kata\", \"label\", \"add\", \"{tag.id}\", \"lifecycle:shipping\"]]\n",
			new:  "shipping = []\n",
			want: table.CatStepsTableInvalid,
		},
		{
			name: "a cleared key missing a clear arm",
			old:  "closed = [[\"kata\", \"label\", \"rm\", \"{tag.id}\", \"lifecycle:closed\"]]\n",
			new:  "",
			want: table.CatStepsTableInvalid,
		},
		{
			name: "a planned-value placeholder",
			old:  "\"lifecycle:shipping\"]]\n",
			new:  "\"{lifecycle}\"]]\n",
			want: table.CatCommandUnknownPlaceholder,
		},
		{
			name: "a tag placeholder naming the owned key",
			old:  "\"lifecycle:shipping\"]]\n",
			new:  "\"{tag.lifecycle}\"]]\n",
			want: table.CatCommandUnknownPlaceholder,
		},
		{
			name: "inline shell in a step",
			old:  "[\"kata\", \"unassign\", \"{tag.id}\"]",
			new:  "[\"sh\", \"-c\", \"kata unassign x\"]",
			want: table.CatCommandShellInterpreter,
		},
		{
			name: "a tag placeholder at argv0",
			old:  "[\"kata\", \"unassign\", \"{tag.id}\"]",
			new:  "[\"{tag.id}\", \"unassign\"]",
			want: table.CatEditTagArgv0,
		},
		{
			name: "an empty step vector",
			old:  "[\"kata\", \"unassign\", \"{tag.id}\"]",
			new:  "[]",
			want: table.CatCommandEmpty,
		},
		{
			name: "an empty step element",
			old:  "[\"kata\", \"unassign\", \"{tag.id}\"]",
			new:  "[\"kata\", \"\", \"{tag.id}\"]",
			want: table.CatCommandEmpty,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := swapEntry(t, stepsModel, tc.old, tc.new)
			if got := loadCategoryOf(t, src, "steps-defect.toml"); got != tc.want {
				t.Errorf("category = %q; want %q", got, tc.want)
			}
		})
	}

	t.Run("steps on a read entry", func(t *testing.T) {
		src := swapEntry(t, stepsModel, "[read.kata.select.attempts]", `[read.kata.steps.lifecycle.set]
resolving = [["kata", "label"]]

[read.kata.select.attempts]`)
		// `steps` beside the read entry's own `command` is the conflict
		// either way; the write-only arm is reached through a read entry
		// whose carrier is `steps` alone.
		if got := loadCategoryOf(t, src, "steps-read.toml"); got != table.CatStepsCarrierConflict {
			t.Errorf("category = %q; want %q", got, table.CatStepsCarrierConflict)
		}
		alone := swapEntry(t, stepsModel, `[read.kata]
role = "kata"
command = ["kata", "show", "{tag.id}", "--json"]
keys = ["lifecycle", "attempts"]
timeout = "2s"

[read.kata.select.lifecycle]
pointer = "/labels"
element = "/label"
prefix = "lifecycle:"

[read.kata.select.attempts]
pointer = "/attempts"
`, `[read.kata]
role = "kata"
keys = ["lifecycle", "attempts"]
timeout = "2s"

[read.kata.steps.lifecycle.set]
resolving = [["kata", "label"]]
`)
		if got := loadCategoryOf(t, alone, "steps-read-alone.toml"); got != table.CatStepsCarrierConflict {
			t.Errorf("category = %q; want %q", got, table.CatStepsCarrierConflict)
		}
	})

	t.Run("steps on an int key", func(t *testing.T) {
		src := swapEntry(t, stepsModel, `[initial]`, `[write.count]
role = "kata"
keys = ["attempts"]
timeout = "2s"
read_back = true

[write.count.steps.attempts.set]
0 = [["kata", "reset", "{tag.id}"]]

[initial]`)
		if got := loadCategoryOf(t, src, "steps-int.toml"); got != table.CatStepsTableInvalid {
			t.Errorf("category = %q; want %q — an int carries its value on "+
				"stdin through `command`, never per-value argv", got, table.CatStepsTableInvalid)
		}
	})

	t.Run("steps on a scalar key", func(t *testing.T) {
		src := swapEntry(t, stepsModel, `[tags.attempts]`, `[tags.note]
provenance = "owned"
kind = "scalar"

[read.note]
role = "kata"
command = ["kata", "note", "{tag.id}"]
keys = ["note"]
timeout = "2s"

[write.note]
role = "kata"
keys = ["note"]
timeout = "2s"
read_back = true

[write.note.steps.note.set]
x = [["kata", "note", "{tag.id}"]]

[tags.attempts]`)
		if got := loadCategoryOf(t, src, "steps-scalar.toml"); got != table.CatStepsTableInvalid {
			t.Errorf("category = %q; want %q", got, table.CatStepsTableInvalid)
		}
	})
}

func TestSteps_ClearArmsAreOptionalWhenNoRuleClears(t *testing.T) {
	// Without `rewind`, nothing plans `<clear>`, so a tool whose `set`
	// replaces natively may declare no clear arms at all.
	src := swapEntry(t, stepsModel, `
[[rule]]
id = "rewind"
clear = ["lifecycle"]
[rule.match.lifecycle]
eq = "shipping"
[rule.match.recognized]
eq = "rewind"
[rule.write]
memo = "rewound"
`, "")
	src = swapEntry(t, src, `
[write.kata.steps.lifecycle.clear]
resolving = [["kata", "label", "rm", "{tag.id}", "lifecycle:resolving"]]
shipping = [["kata", "label", "rm", "{tag.id}", "lifecycle:shipping"]]
closed = [["kata", "label", "rm", "{tag.id}", "lifecycle:closed"]]
`, "")
	if _, err := table.Load([]byte(src), "steps-noclear.toml"); err != nil {
		t.Fatalf("a steps table with no clear arms, and no rule clearing the "+
			"key, refused: %v", err)
	}
}

func TestSteps_CategoriesAreRegistered(t *testing.T) {
	for _, c := range []table.Category{
		table.CatStepsCarrierConflict,
		table.CatStepsKeyMismatch,
		table.CatStepsTableInvalid,
	} {
		if !slices.Contains(table.Categories(), c) {
			t.Errorf("%q refuses at its call site but is not in Categories()", c)
		}
	}
}
