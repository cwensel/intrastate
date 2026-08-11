package resolve_test

import (
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// Probe A: unevaluable ESCAPE row over a no_match candidate set.
// RDR 0007's aggregation clause says this MUST NOT become guard_unevaluable.
func TestProbeA_UnevaluableEscapeOverNoMatch(t *testing.T) {
	esc := escapeRow("rdr.escape.unev", "flows/rdr.toml:90", resolve.KindNoMatch)
	esc.Guard = "unknown-predicate"
	in := noMatchInput()
	in.Table.Revision = "probe-a"
	in.Table.Rows = append(in.Table.Rows, esc)
	in.Guards = fixtureGuards{}
	got, _ := resolve.Resolve(in)
	if got.Refusal == nil {
		t.Fatalf("PROBE A: plan=%+v", got.Plan)
	}
	t.Logf("PROBE A: kind=%q guard=%q  (RDR 0007 clause wants no_match)", got.Refusal.Kind, got.Refusal.Guard)
}

// Probe B: unevaluable CANDIDATE row beside a decided-true sibling.
// RDR 0007 MVV scenario *unevaluable-blocks-true-sibling*.
func TestProbeB_UnevaluableBlocksTrueSibling(t *testing.T) {
	tbl := singleMatchTable()
	tbl.Revision = "probe-b"
	tbl.Rows[0].RuleID = "rdr.true"
	tbl.Rows[0].Guard = "always"
	unev := tbl.Rows[0]
	unev.RuleID = "rdr.unev"
	unev.SourceLocator = "flows/rdr.toml:20"
	unev.Guard = "unknown-predicate"
	tbl.Rows = append(tbl.Rows, unev)
	in := resolve.Input{
		Flow: "rdr", Table: tbl,
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Recognized: "successful",
		Guards:     allGuardsTrue("always"),
	}
	got, _ := resolve.Resolve(in)
	if got.Refusal == nil {
		t.Fatalf("PROBE B: MASKING — plan=%q", got.Plan.RuleID)
	}
	t.Logf("PROBE B: kind=%q guard=%q  (blocks true sibling: correct)", got.Refusal.Kind, got.Refusal.Guard)
}
