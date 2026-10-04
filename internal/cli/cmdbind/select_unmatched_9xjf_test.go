package cmdbind_test

// Kata 9xjf — the READ half of a `select` table's `unmatched` key: a
// prefix projection with zero matching members reads the declared value
// instead of establishing the key absent. Only that one outcome moves;
// every shape `select` already reads as unreadable stays unreadable.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// unmatchedRule is the kata-flight lifecycle reader: an issue with no
// `lifecycle:*` label is `filed`.
func unmatchedRule() table.SelectRule {
	return table.SelectRule{Pointer: "/labels", Element: "/label",
		Prefix: "lifecycle:", Unmatched: "filed"}
}

func TestSelectUnmatched_ZeroMatchesReadTheDeclaredValue(t *testing.T) {
	cases := map[string]string{
		"other labels only": `{"labels":[{"label":"area:cli"},{"label":"type:bug"}]}`,
		"no labels at all":  `{"labels":[]}`,
	}
	for name, show := range cases {
		t.Run(name, func(t *testing.T) {
			sel := map[string]table.SelectRule{fxKey: unmatchedRule()}
			got, unreadable := selected(t, selectEntry(t, show, sel, fxKey))
			if len(unreadable) != 0 {
				t.Fatalf("unreadable = %q; want none", unreadable)
			}
			if v := got[fxKey]; v.Absent || v.Value != "filed" {
				t.Errorf("%s = %#v; want the declared `unmatched` value \"filed\"",
					fxKey, v)
			}
		})
	}
}

func TestSelectUnmatched_OneMatchStillReadsTheMember(t *testing.T) {
	const show = `{"labels":[{"label":"area:cli"},{"label":"lifecycle:resolving"}]}`

	sel := map[string]table.SelectRule{fxKey: unmatchedRule()}
	got, unreadable := selected(t, selectEntry(t, show, sel, fxKey))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if v := got[fxKey]; v.Absent || v.Value != "resolving" {
		t.Errorf("%s = %#v; want the matched member \"resolving\"", fxKey, v)
	}
}

func TestSelectUnmatched_WhatCannotBeEstablishedStaysUnreadable(t *testing.T) {
	cases := []struct {
		name   string
		stdout string
		sel    table.SelectRule
	}{
		{
			name:   "two members carry the prefix",
			stdout: `{"labels":[{"label":"lifecycle:queued"},{"label":"lifecycle:resolving"}]}`,
			sel:    unmatchedRule(),
		},
		{
			name:   "a missing intermediate node",
			stdout: `{"labels":[]}`,
			sel: table.SelectRule{Pointer: "/issue/labels", Element: "/label",
				Prefix: "lifecycle:", Unmatched: "filed"},
		},
		{
			name:   "a null array with no absent declared",
			stdout: `{"labels":null}`,
			sel:    unmatchedRule(),
		},
		{
			name:   "a member the element pointer cannot reach",
			stdout: `{"labels":[{"name":"area:cli"}]}`,
			sel:    unmatchedRule(),
		},
		{
			name:   "a non-array at the pointer",
			stdout: `{"labels":"area:cli"}`,
			sel:    unmatchedRule(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := map[string]table.SelectRule{fxKey: tc.sel}
			got, unreadable := selected(t, selectEntry(t, tc.stdout, sel, fxKey))
			if !slices.Equal(unreadable, []string{fxKey}) {
				t.Errorf("unreadable = %q, values = %#v; want %s UNREADABLE — "+
					"`unmatched` speaks only for a whole set with no match",
					unreadable, got, fxKey)
			}
		})
	}
}

func TestSelectUnmatched_UndeclaredZeroMatchesStayEstablishedAbsent(t *testing.T) {
	const show = `{"labels":[{"label":"area:cli"}]}`

	rule := unmatchedRule()
	rule.Unmatched = ""
	sel := map[string]table.SelectRule{fxKey: rule}
	got, unreadable := selected(t, selectEntry(t, show, sel, fxKey))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if !got[fxKey].Absent {
		t.Errorf("%s = %#v; want established-absent when no `unmatched` "+
			"is declared", fxKey, got[fxKey])
	}
}
