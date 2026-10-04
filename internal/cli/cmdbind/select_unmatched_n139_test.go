package cmdbind_test

// Kata n139 — the READ half of `unmatched` on a pointer select: a field the
// tool reports only in some states reads the declared value wherever the
// select's declared `absent` set would establish it absent. Only the
// established-absent outcome moves; every shape `select` reads as
// unreadable stays unreadable.

import (
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/table"
)

// closedReasonRule is the kata-flight close-reason reader: an open issue
// carries no `closed_reason`, and that reads as "open".
func closedReasonRule() table.SelectRule {
	return table.SelectRule{Pointer: "/issue/closed_reason",
		AbsentNull: true, AbsentMissing: true, Unmatched: "open"}
}

func TestSelectUnmatchedN139_AbsentPointerReadsTheDeclaredValue(t *testing.T) {
	cases := []struct {
		name, stdout, want string
	}{
		{name: "a missing leaf", stdout: `{"issue":{"id":1}}`, want: "open"},
		{name: "a null leaf", stdout: `{"issue":{"id":1,"closed_reason":null}}`, want: "open"},
		{name: "a reported reason", stdout: `{"issue":{"id":1,"closed_reason":"wontfix"}}`, want: "wontfix"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := map[string]table.SelectRule{fxKey: closedReasonRule()}
			got, unreadable := selected(t, selectEntry(t, tc.stdout, sel, fxKey))
			if len(unreadable) != 0 {
				t.Fatalf("unreadable = %q; want none", unreadable)
			}
			if v := got[fxKey]; v.Absent || v.Value != tc.want {
				t.Errorf("%s = %#v; want %q", fxKey, v, tc.want)
			}
		})
	}
}

func TestSelectUnmatchedN139_WhatCannotBeEstablishedStaysUnreadable(t *testing.T) {
	onlyMissing := closedReasonRule()
	onlyMissing.AbsentNull = false
	onlyNull := closedReasonRule()
	onlyNull.AbsentMissing = false

	cases := []struct {
		name   string
		stdout string
		sel    table.SelectRule
	}{
		{name: "a missing intermediate node", stdout: `{}`, sel: closedReasonRule()},
		{name: "an object at the leaf", stdout: `{"issue":{"closed_reason":{"v":"done"}}}`, sel: closedReasonRule()},
		{name: "an array at the leaf", stdout: `{"issue":{"closed_reason":["done"]}}`, sel: closedReasonRule()},
		{name: "null with only missing declared", stdout: `{"issue":{"closed_reason":null}}`, sel: onlyMissing},
		{name: "missing with only null declared", stdout: `{"issue":{"id":1}}`, sel: onlyNull},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sel := map[string]table.SelectRule{fxKey: tc.sel}
			got, unreadable := selected(t, selectEntry(t, tc.stdout, sel, fxKey))
			if !slices.Equal(unreadable, []string{fxKey}) {
				t.Errorf("unreadable = %q, values = %#v; want %s UNREADABLE — "+
					"`unmatched` speaks only for a declared-absent shape",
					unreadable, got, fxKey)
			}
		})
	}
}

// A prefix projection that also declares `absent` reads `unmatched` for a
// null array, not established-absent: every established-absent outcome of
// a select carrying `unmatched` reads the declared value.
func TestSelectUnmatchedN139_PrefixWithAbsentNullReadsUnmatched(t *testing.T) {
	rule := unmatchedRule()
	rule.AbsentNull = true
	sel := map[string]table.SelectRule{fxKey: rule}
	got, unreadable := selected(t, selectEntry(t, `{"labels":null}`, sel, fxKey))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if v := got[fxKey]; v.Absent || v.Value != "filed" {
		t.Errorf("%s = %#v; want the declared `unmatched` value \"filed\"", fxKey, v)
	}
}

func TestSelectUnmatchedN139_UndeclaredAbsentPointerStaysEstablishedAbsent(t *testing.T) {
	rule := closedReasonRule()
	rule.Unmatched = ""
	sel := map[string]table.SelectRule{fxKey: rule}
	got, unreadable := selected(t, selectEntry(t, `{"issue":{"id":1}}`, sel, fxKey))
	if len(unreadable) != 0 {
		t.Fatalf("unreadable = %q; want none", unreadable)
	}
	if !got[fxKey].Absent {
		t.Errorf("%s = %#v; want established-absent when no `unmatched` "+
			"is declared", fxKey, got[fxKey])
	}
}
