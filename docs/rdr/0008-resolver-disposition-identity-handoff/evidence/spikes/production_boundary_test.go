package spikes_test

import (
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolver"
)

func TestA2ProductionBoundary(t *testing.T) {
	t.Run("ordinary result owns copied action", func(t *testing.T) {
		table := []resolver.Edge{{
			Outcome:  "accepted",
			NextTags: resolver.TagSet{"stage": "done"},
			Writes:   []resolver.OwnedTagWrite{{Role: "state", Name: "stage", Value: "done"}},
		}}
		in := resolver.Input{
			Owned:      resolver.OwnedSnapshot{Available: true},
			Recognized: resolver.Tag{Name: "result", Value: "accepted"},
			Outcomes:   []string{"accepted"},
			Table:      table,
		}
		got, err := resolver.Resolve(in)
		if err != nil || got.Plan == nil || got.Refusal != nil {
			t.Fatalf("Resolve() = (%#v, %v), want ordinary plan", got, err)
		}
		table[0].NextTags["stage"] = "mutated-after"
		table[0].Writes[0].Value = "mutated-after"
		if !reflect.DeepEqual(got.Plan.NextTags, resolver.TagSet{"stage": "done"}) ||
			!reflect.DeepEqual(got.Plan.Writes, []resolver.OwnedTagWrite{{Role: "state", Name: "stage", Value: "done"}}) {
			t.Fatalf("ordinary plan aliased table storage: %#v", got.Plan)
		}
		t.Log("PASS ordinary action result remains unchanged after source storage mutation")
	})

	t.Run("caller mutation before resolve changes selected row", func(t *testing.T) {
		table := []resolver.Edge{{Outcome: "accepted", NextTags: resolver.TagSet{"stage": "done"}}}
		in := resolver.Input{
			Owned:      resolver.OwnedSnapshot{Available: true},
			Recognized: resolver.Tag{Name: "result", Value: "accepted"},
			Outcomes:   []string{"accepted"},
			Table:      table,
		}
		table[0].Outcome = "other"
		got, err := resolver.Resolve(in)
		if err != nil || got.Plan != nil || got.Refusal == nil || got.Refusal.Kind != resolver.RefusalNoMatch {
			t.Fatalf("Resolve() = (%#v, %v), want no_match after caller mutation", got, err)
		}
		t.Log("FALSIFIED snapshot-before-resolution: exported Input.Table observes caller mutation")
	})

	t.Run("modeled escape is action-bearing plan", func(t *testing.T) {
		table := []resolver.Edge{
			{Outcome: "accepted", Guard: resolver.Guard{All: []resolver.TagPredicate{{Provenance: resolver.ProvenanceRecognized, Name: "result", Value: "rejected"}}}},
			{Outcome: "accepted", EscapeFor: resolver.RefusalNoMatch, NextTags: resolver.TagSet{"stage": "escaped"}, Writes: []resolver.OwnedTagWrite{{Role: "state", Name: "stage", Value: "escaped"}}},
		}
		in := resolver.Input{
			Owned:      resolver.OwnedSnapshot{Available: true},
			Recognized: resolver.Tag{Name: "result", Value: "accepted"},
			Outcomes:   []string{"accepted"},
			Table:      table,
		}
		got, err := resolver.Resolve(in)
		if err != nil || got.Plan == nil || got.Refusal != nil {
			t.Fatalf("Resolve() = (%#v, %v), want current production escape plan", got, err)
		}
		table[1].NextTags["stage"] = "mutated-after"
		table[1].Writes[0].Value = "mutated-after"
		if got.Plan.NextTags["stage"] != "escaped" || got.Plan.Writes[0].Value != "escaped" {
			t.Fatalf("escape plan aliased table storage: %#v", got.Plan)
		}
		t.Log("BLOCK modeled-escape A2 proof: production returns an action-bearing plan, not an identity-bearing refusal")
	})

	t.Run("external caller constructs and replaces row components", func(t *testing.T) {
		row := resolver.Edge{Outcome: "accepted", EscapeFor: resolver.RefusalNoMatch}
		row.Guard = resolver.Guard{Unevaluable: true}
		row.NextTags = resolver.TagSet{"externally": "replaced"}
		row.Writes = []resolver.OwnedTagWrite{{Role: "state", Name: "externally", Value: "replaced"}}
		in := resolver.Input{Table: []resolver.Edge{row}}
		in.Table[0] = resolver.Edge{Outcome: "different"}
		if in.Table[0].Outcome != "different" {
			t.Fatal("unexpected inability to replace row")
		}
		t.Log("FALSIFIED external construction constraint: Edge fields and Input.Table are exported and replaceable")
	})
}
