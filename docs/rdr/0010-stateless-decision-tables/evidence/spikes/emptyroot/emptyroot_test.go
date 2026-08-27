//go:build rdr_spike

// Package emptyroot is an RDR 0010 A1 SPIKE. It is not production code and
// is not wired into any build target: it exists to execute C5's proposed
// "root the reachability relation at the empty owned-state node" WITHOUT
// editing internal/.
//
// The trick: `reach` bails only on len(m.Initial) == 0, and then SKIPS every
// root assignment whose tag is not owned-provenance. So an [initial] over a
// NON-owned key yields root = Node{Values: {}} -- literally the empty
// owned-state node C5 names -- while every owned tag, write block, accessor
// and terminal is stripped from the in-memory model first.
package emptyroot

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/table"
)

// strip removes every owned tag, the declared root, the terminals, and every
// row's write / next / requires-owned, leaving the pure decision table.
func strip(m *table.Model) {
	for k, d := range m.Tags {
		if d.Provenance == table.ProvenanceOwned {
			delete(m.Tags, k)
		}
	}
	m.Initial = nil
	m.Terminal = nil
	for i := range m.Rows {
		m.Rows[i].Writes = nil
		m.Rows[i].NextTags = nil
		m.Rows[i].RequiresOwned = nil
	}
}

func run(t *testing.T, m *table.Model, label string) {
	t.Helper()
	rep := graphlint.Run(graphlint.NewRequest(m))
	b, _ := json.MarshalIndent(rep.Findings, "", "  ")
	t.Logf("%s: nodes=%d findings=%d\n%s",
		label, len(graphlint.Reach(m)), len(rep.Findings), b)
}

func load(t *testing.T) *table.Model {
	t.Helper()
	src, err := os.ReadFile("../a1-baseline.toml")
	if err != nil {
		t.Fatal(err)
	}
	m, err := table.Load(src, "a1-baseline.toml")
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// TestA_Baseline is the control: the model exactly as it loads.
func TestA_Baseline(t *testing.T) {
	run(t, load(t), "A/baseline (scratch owned tag, real root)")
}

// TestB_StrippedNoRoot is the CURRENT behavior of a zero-owned model:
// m.Initial is empty, reach returns nil, no group is reachable, coverage is
// vacuous. This is the false green C5 exists to close.
func TestB_StrippedNoRoot(t *testing.T) {
	m := load(t)
	strip(m)
	run(t, m, "B/stripped, no root (CURRENT zero-owned behavior)")
}

// TestC_StrippedEmptyRoot is C5's proposal EXECUTED: the same stripped model,
// with the traversal rooted at the empty owned-state node. The root
// assignment names an OBSERVED key, which reach skips -- so the root node it
// builds is Node{Values: {}}, the empty owned-state node, and nothing else
// about the run changes.
func TestC_StrippedEmptyRoot(t *testing.T) {
	m := load(t)
	strip(m)
	m.Initial = []table.TagValue{{Key: "d0", Value: []string{"false"}}}
	nodes := graphlint.Reach(m)
	if len(nodes) != 1 || len(nodes[0].Values) != 0 {
		t.Fatalf("expected exactly the empty owned-state node, got %#v", nodes)
	}
	run(t, m, "C/stripped, EMPTY-ROOT seeded (C5 proposal)")
}
