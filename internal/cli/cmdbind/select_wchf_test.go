package cmdbind_test

// Kata wchf — the READ half of a `command` read entry's per-key `select`
// tables: RFC 6901 pointer resolution over the tool's own nested json,
// prefix projection out of a label array, and the per-key declared
// absence. The load half is `internal/table/command_select_wchf_test.go`.
//
// Every fixture is the real helper child emitting a `kata show`-shaped or
// `kata list`-shaped document, so the scenario is the consumer's: an
// owned enum read out of `labels[]` by prefix, and an owner that is null
// or missing when nobody holds the issue — read with no wrapper script.

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// kataShowSelect is the selector set of a `kata show --json` reader.
func kataShowSelect() map[string]table.SelectRule {
	return map[string]table.SelectRule{
		"lifecycle": {Pointer: "/labels", Element: "/label", Prefix: "lifecycle:"},
		"owner":     {Pointer: "/issue/owner", AbsentNull: true, AbsentMissing: true},
	}
}

// selectEntry builds a json-mode read entry whose child emits stdout.
func selectEntry(t *testing.T, stdout string, sel map[string]table.SelectRule,
	keys ...string,
) table.Accessor {
	t.Helper()

	acc := entry(fxRole, []string{helperBin(t), "emit", stdout, "0"}, keys...)
	acc.Select = sel
	return acc
}

// selected drives one read and indexes the answer by key.
func selected(t *testing.T, acc table.Accessor) (map[string]accessor.KeyValue, []string) {
	t.Helper()

	values, unreadable, err := readOnce(t, acc, artifactAt(t, "kata.json"))
	if err != nil {
		t.Fatalf("the read failed: %v", err)
	}
	got := map[string]accessor.KeyValue{}
	for _, v := range values {
		got[v.Key] = v
	}
	return got, unreadable
}

func TestSelect_ReadsNestedKataShowOutputWithNoWrapper(t *testing.T) {
	const show = `{"issue":{"short_id":"wchf","owner":"x","labels":null},` +
		`"labels":[{"label":"area:cli"},{"label":"lifecycle:refining"},` +
		`{"label":"type:feature"}]}`

	got, unreadable := selected(t,
		selectEntry(t, show, kataShowSelect(), "lifecycle", "owner"))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if v := got["lifecycle"]; v.Absent || v.Value != "refining" {
		t.Errorf("lifecycle = %#v; want the prefix-stripped value \"refining\"", v)
	}
	if v := got["owner"]; v.Absent || v.Value != "x" {
		t.Errorf("owner = %#v; want \"x\"", v)
	}
}

func TestSelect_DeclaredAbsenceIsEstablishedAbsent(t *testing.T) {
	cases := map[string]string{
		"owner null":    `{"issue":{"owner":null},"labels":[]}`,
		"owner missing": `{"issue":{},"labels":[]}`,
	}
	for name, show := range cases {
		t.Run(name, func(t *testing.T) {
			got, unreadable := selected(t,
				selectEntry(t, show, kataShowSelect(), "lifecycle", "owner"))
			if len(unreadable) != 0 {
				t.Fatalf("unreadable = %q; want none", unreadable)
			}
			if !got["owner"].Absent {
				t.Errorf("owner = %#v; want established-absent through the "+
					"declared `absent`", got["owner"])
			}
			// No member carries the prefix in a set the tool reported
			// whole, so the projected key is absent too.
			if !got["lifecycle"].Absent {
				t.Errorf("lifecycle = %#v; want established-absent: zero "+
					"prefix matches", got["lifecycle"])
			}
		})
	}
}

func TestSelect_StringArrayProjectsWithNoElementPointer(t *testing.T) {
	const list = `{"issues":[{"short_id":"wchf","labels":["area:cli","lifecycle:queued"]}]}`

	sel := map[string]table.SelectRule{
		fxKey: {Pointer: "/issues/0/labels", Prefix: "lifecycle:"},
	}
	got, unreadable := selected(t, selectEntry(t, list, sel, fxKey))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if v := got[fxKey]; v.Value != "queued" {
		t.Errorf("%s = %#v; want \"queued\"", fxKey, v)
	}
}

func TestSelect_ScalarsRenderAsTheirJSONText(t *testing.T) {
	const doc = `{"n":7,"big":12345678901234567890,"f":1.50,"b":true,"s":"x/y"}`

	sel := map[string]table.SelectRule{
		"n": {Pointer: "/n"}, "big": {Pointer: "/big"}, "f": {Pointer: "/f"},
		"b": {Pointer: "/b"}, "s": {Pointer: "/s"},
	}
	got, unreadable := selected(t, selectEntry(t, doc, sel, "n", "big", "f", "b", "s"))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	want := map[string]string{
		"n": "7", "big": "12345678901234567890", "f": "1.50", "b": "true", "s": "x/y",
	}
	for k, w := range want {
		if got[k].Value != w {
			t.Errorf("%s = %q; want %q", k, got[k].Value, w)
		}
	}
}

func TestSelect_WhatCannotBeEstablishedIsUnreadable(t *testing.T) {
	cases := []struct {
		name   string
		stdout string
		sel    table.SelectRule
	}{
		{
			name:   "two members carry the prefix",
			stdout: `{"labels":[{"label":"lifecycle:queued"},{"label":"lifecycle:shipping"}]}`,
			sel:    table.SelectRule{Pointer: "/labels", Element: "/label", Prefix: "lifecycle:"},
		},
		{
			name:   "a missing intermediate node, even with absent = missing",
			stdout: `{"labels":[]}`,
			sel:    table.SelectRule{Pointer: "/issue/owner", AbsentNull: true, AbsentMissing: true},
		},
		{
			name:   "a null with no absent declared",
			stdout: `{"issue":{"owner":null}}`,
			sel:    table.SelectRule{Pointer: "/issue/owner"},
		},
		{
			name:   "a missing final token with only null declared",
			stdout: `{"issue":{}}`,
			sel:    table.SelectRule{Pointer: "/issue/owner", AbsentNull: true},
		},
		{
			name:   "an object target",
			stdout: `{"issue":{"owner":{"name":"x"}}}`,
			sel:    table.SelectRule{Pointer: "/issue/owner"},
		},
		{
			name:   "an array target with no prefix",
			stdout: `{"labels":["a"]}`,
			sel:    table.SelectRule{Pointer: "/labels"},
		},
		{
			name:   "a prefix over a non-array",
			stdout: `{"labels":"lifecycle:queued"}`,
			sel:    table.SelectRule{Pointer: "/labels", Prefix: "lifecycle:"},
		},
		{
			name:   "a member the element pointer cannot reach",
			stdout: `{"labels":[{"name":"lifecycle:queued"},{"label":"area:cli"}]}`,
			sel:    table.SelectRule{Pointer: "/labels", Element: "/label", Prefix: "lifecycle:"},
		},
		{
			name:   "a non-index token into an array",
			stdout: `{"issues":[{"owner":"x"}]}`,
			sel:    table.SelectRule{Pointer: "/issues/first/owner"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := map[string]table.SelectRule{fxKey: tc.sel}
			got, unreadable := selected(t, selectEntry(t, tc.stdout, sel, fxKey))
			if !slices.Equal(unreadable, []string{fxKey}) {
				t.Errorf("unreadable = %q, values = %#v; want %s UNREADABLE — "+
					"a value the tool did not establish is never guessed",
					unreadable, got, fxKey)
			}
		})
	}
}

func TestSelect_StdoutThatIsNotOneDocumentIsAnExecutionFailure(t *testing.T) {
	sel := map[string]table.SelectRule{fxKey: {Pointer: "/a"}}
	for _, stdout := range []string{`not json`, `{"a":"x"} {"a":"y"}`} {
		_, _, err := readOnce(t, selectEntry(t, stdout, sel, fxKey), artifactAt(t, "kata.json"))
		if err == nil {
			t.Errorf("stdout %q read clean; want an execution failure", stdout)
		}
	}
}

func TestSelect_AnUnparseableSelectorRefusesBeforeSpawn(t *testing.T) {
	trace := filepath.Join(t.TempDir(), "trace")
	acc := entry(fxRole, []string{helperBin(t), "trace", trace}, fxKey)
	acc.Select = map[string]table.SelectRule{fxKey: {Pointer: "no-slash"}}

	if _, _, err := readOnce(t, acc, artifactAt(t, "kata.json")); err == nil {
		t.Fatal("an in-memory selector that does not parse read clean")
	}
	if _, err := os.Stat(trace); err == nil {
		t.Error("the child ran; a selector defect must refuse before the spawn")
	}
}

// The consumer's acceptance: a write through a tool-native writer is
// verified by reading the label back through the SAME selector entry.
func TestSelect_WriteReadsBackThroughTheSelector(t *testing.T) {
	const key = "lifecycle"
	bin := helperBin(t)

	writeThrough := func(t *testing.T, writeArgv func(art string) []string) accessor.WriteResult {
		t.Helper()

		art := artifactAt(t, "kata.json")
		seed := `{"issue":{"owner":null},"labels":[{"label":"lifecycle:queued"}]}`
		if err := os.WriteFile(art.Path, []byte(seed), 0o600); err != nil {
			t.Fatal(err)
		}
		wacc := entry(fxRole, writeArgv(art.Path), key)
		wacc.ReadBack = true
		racc := entry(fxRole, []string{"cat", "{artifact}"}, key)
		racc.Select = map[string]table.SelectRule{
			key: {Pointer: "/labels", Element: "/label", Prefix: "lifecycle:"},
		}
		reg := accessor.Registry{
			Flow: "kataflow",
			Definitions: []accessor.Definition{
				{
					Identity: accessor.Identity{
						Flow: "kataflow", Name: "kata.read", Capability: accessor.CapRead,
					},
					Accessor: racc,
					Binding:  cmdbind.Reader{Accessor: racc, Name: "kata.read", Config: allowed(t)},
				},
				{
					Identity: accessor.Identity{
						Flow: "kataflow", Name: "kata.write", Capability: accessor.CapWrite,
					},
					Accessor: wacc,
					Binding:  &cmdbind.Writer{Accessor: wacc, Name: "kata.write", Config: allowed(t)},
				},
			},
			OwnedTags: []string{key},
		}
		return accessor.NewExecutor(reg, accessor.Artifacts{fxRole: art}).
			Write(ctxOf(t), "kata.write", resolve.Plan{
				Writes: []resolve.Tag{{Key: key, Value: "shipping"}},
			})
	}

	t.Run("the label lands: the write verifies", func(t *testing.T) {
		got := writeThrough(t, func(art string) []string {
			return []string{bin, "stdin-to-labels", art}
		})
		if got.Refusal != nil {
			t.Fatalf("the write refused %q (Detail %q)", got.Refusal.Class, got.Refusal.Detail)
		}
		if len(got.Written) != 1 || got.Written[0].Value != "shipping" {
			t.Errorf("Written = %#v; want the verified lifecycle", got.Written)
		}
	})

	t.Run("the tool exits zero without applying: read_back_mismatch", func(t *testing.T) {
		got := writeThrough(t, func(string) []string {
			return []string{bin, "silent", "0"}
		})
		if got.Refusal == nil || got.Refusal.Class != accessor.ClassReadBackMismatch {
			t.Errorf("refusal = %#v; want %q — the selector read the old label",
				got.Refusal, accessor.ClassReadBackMismatch)
		}
	})
}
