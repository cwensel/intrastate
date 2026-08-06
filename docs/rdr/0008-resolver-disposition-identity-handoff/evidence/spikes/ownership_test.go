package spikes

import (
	"maps"
	"slices"
	"testing"
)

type selectedRule struct {
	modelID, ruleID, expansionSuffix, sourceLocator string
}

type ownedTagWrite struct {
	name, value string
}

type transitionAction struct {
	nextTags map[string]string
	writes   []ownedTagWrite
}

type edge struct {
	selectedRule selectedRule
	action       transitionAction
}

type transitionPlan struct {
	selectedRule selectedRule
	action       transitionAction
}

func plan(edge edge) transitionPlan {
	return transitionPlan{
		selectedRule: edge.selectedRule,
		action: transitionAction{
			nextTags: maps.Clone(edge.action.nextTags),
			writes:   slices.Clone(edge.action.writes),
		},
	}
}

func TestPlanOwnsSelectedIdentityAndAction(t *testing.T) {
	input := edge{
		selectedRule: selectedRule{"model-a", "rule-a", "#1", "rules.toml:7"},
		action: transitionAction{
			nextTags: map[string]string{"status": "Final"},
			writes:   []ownedTagWrite{{"status", "Final"}},
		},
	}
	got := plan(input)

	// Mutate the caller-owned backing storage, then reuse the edge value.
	input.selectedRule.ruleID = "mutated-rule"
	input.action.nextTags["status"] = "Draft"
	input.action.writes[0].value = "Draft"
	input = edge{
		selectedRule: selectedRule{"model-b", "rule-b", "", "rules.toml:20"},
		action: transitionAction{
			nextTags: map[string]string{"status": "Blocked"},
			writes:   []ownedTagWrite{{"status", "Blocked"}},
		},
	}
	_ = plan(input)

	want := transitionPlan{
		selectedRule: selectedRule{"model-a", "rule-a", "#1", "rules.toml:7"},
		action: transitionAction{
			nextTags: map[string]string{"status": "Final"},
			writes:   []ownedTagWrite{{"status", "Final"}},
		},
	}
	if got.selectedRule != want.selectedRule ||
		!maps.Equal(got.action.nextTags, want.action.nextTags) ||
		!slices.Equal(got.action.writes, want.action.writes) {
		t.Fatalf("plan changed after input mutation/reuse: got %#v, want %#v", got, want)
	}

	t.Logf("preserved selected_rule=%s/%s%s source=%s next_status=%s write=%s:%s",
		got.selectedRule.modelID,
		got.selectedRule.ruleID,
		got.selectedRule.expansionSuffix,
		got.selectedRule.sourceLocator,
		got.action.nextTags["status"],
		got.action.writes[0].name,
		got.action.writes[0].value,
	)
}
