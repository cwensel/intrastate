package resolve_test

import (
	"reflect"
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// REQ-MVV: "Resolve must name and implementation must add a replay test
// that feeds the same table, owned snapshot, observed tags, and recognized
// outcome to the kernel twice and asserts value-identical dispositions. The
// same validation must include at least one value-level refusal each for
// `no_match`, `ambiguous_match`, `owned_state_unavailable`,
// `guard_unevaluable`, and `unmodeled_outcome`, and must assert those
// modeled refusals do not use the CLI or Go error path."
// HAPPY PATH
//
// This is the gating Minimum Viable Validation for RDR 0001. It is one
// runnable end-to-end test with three obligations in a single suite:
//
//  1. replayed input tuple yields value-identical dispositions;
//  2. one value-level refusal case per each of the five kernel refusal
//     kinds;
//  3. every modeled refusal travels the Result value — not the CLI path,
//     not the Go error path.
//
// RDR 0001 declares no Round-Trip / Inverse Invariant ("No encode/decode,
// import/export, or inverse operation is introduced by this RDR"), so the
// replay obligation is value-for-value equality of the disposition rather
// than a reconstruct-and-compare round trip. A green exit or "did not
// error" is explicitly not sufficient here: every leg below compares
// values.
func TestMVV_ReplayDeterminismAndFiveValueLevelRefusals(t *testing.T) {
	t.Run("1_replay_yields_value_identical_dispositions", func(t *testing.T) {
		first, err1 := resolve.Resolve(legalInput())
		second, err2 := resolve.Resolve(legalInput())

		if err1 != nil || err2 != nil {
			t.Fatalf("legal replay used the Go error path: %v / %v", err1, err2)
		}
		if first.Refusal != nil {
			t.Fatalf("legal input tuple refused with kind %q", first.Refusal.Kind)
		}
		if first.Plan == nil || second.Plan == nil {
			t.Fatal("legal input tuple produced no transition plan")
		}

		// Value-for-value, not "did not error".
		if !reflect.DeepEqual(first, second) {
			t.Fatalf("replayed dispositions are not value-identical:\nfirst  = %+v\nsecond = %+v", first, second)
		}
		if !reflect.DeepEqual(first.Plan.NextTags, second.Plan.NextTags) {
			t.Errorf("next tags differ across replay: %+v vs %+v",
				first.Plan.NextTags, second.Plan.NextTags)
		}
		if !reflect.DeepEqual(first.Plan.Writes, second.Plan.Writes) {
			t.Errorf("accessor write descriptions differ across replay: %+v vs %+v",
				first.Plan.Writes, second.Plan.Writes)
		}

		// The replayed plan must be substantive, or "identical" is vacuous.
		wantNext := []resolve.Tag{{Key: "status", Value: "Final"}}
		if !reflect.DeepEqual(first.Plan.NextTags, wantNext) {
			t.Errorf("next tags = %+v; want %+v", first.Plan.NextTags, wantNext)
		}
		wantWrites := []resolve.Tag{{Key: "status", Value: "Final"}}
		if !reflect.DeepEqual(first.Plan.Writes, wantWrites) {
			t.Errorf("write descriptions = %+v; want %+v", first.Plan.Writes, wantWrites)
		}
		if first.Plan.Revision != "rev-1" {
			t.Errorf("plan revision = %q; want %q", first.Plan.Revision, "rev-1")
		}
	})

	t.Run("2_one_value_level_refusal_per_kind", func(t *testing.T) {
		cases := []struct {
			name string
			in   resolve.Input
			kind resolve.RefusalKind
		}{
			{"no_match", noMatchInput(), resolve.KindNoMatch},
			{"ambiguous_match", ambiguousInput(), resolve.KindAmbiguousMatch},
			{"owned_state_unavailable", missingOwnedInput(), resolve.KindOwnedStateUnavailable},
			{"guard_unevaluable", unevaluableGuardInput(), resolve.KindGuardUnevaluable},
			{"unmodeled_outcome", unmodeledOutcomeInput(), resolve.KindUnmodeledOutcome},
		}

		covered := map[resolve.RefusalKind]bool{}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := resolve.Resolve(tc.in)

				// Obligation 3, per case: not the Go error path.
				if err != nil {
					t.Fatalf("modeled refusal %q traveled the Go error path: %v", tc.kind, err)
				}
				// Obligation 3, per case: a value-level disposition.
				if got.Refusal == nil {
					t.Fatalf("modeled refusal %q did not travel the Result value", tc.kind)
				}
				if got.Plan != nil {
					t.Fatalf("refusal %q also carries a transition plan: %+v", tc.kind, *got.Plan)
				}
				if got.Refusal.Kind != tc.kind {
					t.Fatalf("refusal kind = %q; want %q", got.Refusal.Kind, tc.kind)
				}

				// No write description travels a refusal.
				if got.Plan != nil && len(got.Plan.Writes) > 0 {
					t.Error("refusal carries an accessor write description")
				}

				// The refusal is diagnosable from the input tuple.
				if got.Refusal.Revision != tc.in.Table.Revision {
					t.Errorf("refusal revision = %q; want %q",
						got.Refusal.Revision, tc.in.Table.Revision)
				}
			})
			covered[tc.kind] = true
		}

		for _, k := range resolve.RefusalKinds() {
			if !covered[k] {
				t.Errorf("refusal kind %q has no value-level case in the MVV", k)
			}
		}
		if len(covered) != 5 {
			t.Errorf("MVV covers %d refusal kinds; want 5", len(covered))
		}
	})

	t.Run("3_refusals_use_neither_the_cli_nor_the_go_error_path", func(t *testing.T) {
		for name, build := range refusalInputBuilders() {
			t.Run(name, func(t *testing.T) {
				got, err := resolve.Resolve(build())
				if err != nil {
					t.Fatalf("refusal traveled the Go error path: %v", err)
				}
				if got.Refusal == nil {
					t.Fatal("refusal did not travel the value path")
				}

				// Not the CLI path: the refusal value is not an error at
				// all, so nothing downstream can mistake it for one, and
				// the kind alone — with no error and no string
				// inspection — is what RDR 0005 maps.
				if _, isErr := any(got.Refusal).(error); isErr {
					t.Error("the refusal value implements error; modeled refusals must not be errors")
				}
				if _, isErr := any(*got.Refusal).(error); isErr {
					t.Error("the refusal value implements error; modeled refusals must not be errors")
				}
				if _, isErr := any(got).(error); isErr {
					t.Error("the Result value implements error; dispositions must not be errors")
				}
				if got.Refusal.Kind == "" {
					t.Error("refusal carries no kind for the CLI layer to map")
				}
			})
		}

		// The kernel package itself is not wired to the CLI: refusals
		// cannot reach a CLI error envelope from inside the kernel.
		for imp := range kernelPackageImports(t) {
			if isCLIPackage(imp) {
				t.Errorf("kernel imports CLI package %q; refusals must not travel the CLI path", imp)
			}
		}
	})

	t.Run("4_replay_holds_for_refusal_dispositions_too", func(t *testing.T) {
		for name, build := range refusalInputBuilders() {
			t.Run(name, func(t *testing.T) {
				first, err1 := resolve.Resolve(build())
				second, err2 := resolve.Resolve(build())
				if err1 != nil || err2 != nil {
					t.Fatalf("refusal replay used the Go error path: %v / %v", err1, err2)
				}
				if !reflect.DeepEqual(first, second) {
					t.Errorf("replayed refusal dispositions are not value-identical:\nfirst  = %+v\nsecond = %+v",
						first, second)
				}
			})
		}
	})
}
