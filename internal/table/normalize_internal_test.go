package table

import (
	"slices"
	"testing"
)

func TestCompareRowsSortsBySecondSuffixElement(t *testing.T) {
	rows := []Row{
		{ModelID: "model", RuleID: "rule", Suffix: []string{"same", "z"}},
		{ModelID: "model", RuleID: "rule", Suffix: []string{"same", "a"}},
	}

	slices.SortFunc(rows, compareRows)

	if got, want := rows[0].Suffix, []string{"same", "a"}; !slices.Equal(got, want) {
		t.Fatalf("first suffix = %v; want %v", got, want)
	}
}
