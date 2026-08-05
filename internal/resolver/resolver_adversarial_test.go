package resolver_test

import (
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolver"
)

// REQ-14: "The kernel MUST refuse instead of guessing when ... a guard cannot
// be evaluated ... ."
// ADVERSARIAL
func TestResolveAdversarialRefusesInsteadOfGuessing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   resolver.Input
		want resolver.RefusalKind
	}{
		{
			name: "one apparent match cannot hide an unevaluable competitor",
			in: resolver.Input{
				Owned:      resolver.OwnedSnapshot{Available: true},
				Recognized: resolver.Tag{Name: "result", Value: "accepted"},
				Table: []resolver.Edge{
					{Outcome: "accepted"},
					{Outcome: "accepted", Guard: resolver.Guard{Unevaluable: true}},
				},
			},
			want: resolver.RefusalGuardUnevaluable,
		},
		{
			name: "multiple ambiguity escapes remain ambiguous",
			in: resolver.Input{
				Owned:      resolver.OwnedSnapshot{Available: true},
				Recognized: resolver.Tag{Name: "result", Value: "accepted"},
				Table: []resolver.Edge{
					{Outcome: "accepted"},
					{Outcome: "accepted"},
					{Outcome: "accepted", EscapeFor: resolver.RefusalAmbiguousMatch},
					{Outcome: "accepted", EscapeFor: resolver.RefusalAmbiguousMatch},
				},
			},
			want: resolver.RefusalAmbiguousMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resolver.Resolve(tt.in)
			if err != nil {
				t.Fatalf("Resolve() error = %v, want value-level refusal", err)
			}
			if got.Plan != nil {
				t.Fatalf("Resolve() plan = %#v, want refusal", got.Plan)
			}
			if got.Refusal == nil || got.Refusal.Kind != tt.want {
				t.Fatalf("Resolve() refusal = %#v, want kind %q", got.Refusal, tt.want)
			}
		})
	}
}

// REQ-25: "Replaying that tuple must replay the disposition."
// ADVERSARIAL
func TestResolveAdversarialReturnedPlanCannotMutateReplay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "caller mutation of a plan does not mutate the supplied table"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			in := resolver.Input{
				FlowID:        "rdr",
				TableRevision: "rev-1",
				Owned:         resolver.OwnedSnapshot{Available: true},
				Recognized:    resolver.Tag{Name: "result", Value: "accepted"},
				Table: []resolver.Edge{
					{
						Outcome:  "accepted",
						NextTags: resolver.TagSet{"stage": "done"},
						Writes: []resolver.OwnedTagWrite{
							{Role: "state", Name: "stage", Value: "done"},
						},
					},
				},
			}

			first, err := resolver.Resolve(in)
			if err != nil || first.Plan == nil {
				t.Fatalf("first Resolve() = (%#v, %v), want transition plan", first, err)
			}
			first.Plan.NextTags["stage"] = "corrupted"
			first.Plan.Writes[0].Value = "corrupted"

			got, err := resolver.Resolve(in)
			if err != nil || got.Plan == nil {
				t.Fatalf("replayed Resolve() = (%#v, %v), want transition plan", got, err)
			}
			want := &resolver.TransitionPlan{
				NextTags: resolver.TagSet{"stage": "done"},
				Writes: []resolver.OwnedTagWrite{
					{Role: "state", Name: "stage", Value: "done"},
				},
			}
			if !reflect.DeepEqual(got.Plan, want) {
				t.Fatalf("replayed plan = %#v, want %#v", got.Plan, want)
			}
		})
	}
}
