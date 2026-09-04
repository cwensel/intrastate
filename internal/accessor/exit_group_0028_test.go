package accessor_test

// RDR 0028 Phase 3c — the THIRD entry-level precondition site.
//
// `0028:C1.3` `order:` enumerates the entry-level preconditions as a set
// of three and then fixes their exit group in one sentence that speaks of
// all of them at once:
//
//	"…and the entry-level preconditions — the read-back gate below,
//	C1.6's unbound tag and C1.6's `-`-prefixed value — name the gate or
//	the placeholder instead, having no rule to name; … EXIT GROUP: these
//	refusals are about the REQUEST, not the environment, so they take the
//	exit-2 group, not `execution_failure`'s default exit 3 — the
//	executor-facing typed `Err` discriminates at
//	`flow_exec.go::accessorFailureOf`"
//
// Phase 3b pinned two of the three (ADV-1's unbound anchor tag, ADV-2's
// read-back gate). The third — C1.6's argv-side pair, minted in
// `cmdbind.go::substitute` — was owed the same fix and is pinned here.
//
// The Phase 1 suite already asserts these two refusals happen BEFORE
// spawn and that the Detail names the placeholder (REQ-88, REQ-89). What
// it does not assert is the discriminator riding beside them, which is
// exactly the shape ADV-1 and ADV-2 were about: the class, the Detail and
// the untouched artifact are all correct, and only the exit group is
// wrong. A caller told to "repair the environment and re-run the same
// request unchanged" over a missing `--tag` retries forever.

import (
	"context"
	"errors"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-41/REQ-43 (`0028:C1.3` order:/EXIT GROUP:) over REQ-88 and REQ-89
// (`0028:C1.6` binding:/VALUE:).
//
// Both C1.6 argv refusals are entry-level preconditions of the same set
// C1.3's EXIT GROUP: sentence governs, so both must wrap the typed
// `accessor.ErrDeclaredRequest` that `refusalOf` reads.
//
// The admitted arm is asserted too: a bound, non-flag value must NOT be
// refused, so the fix cannot be read as "route everything to exit 2".
func TestExitGroup0028_CommandTagPreconditionsTakeTheRequestExitGroup(t *testing.T) {
	for _, tc := range []struct {
		name    string
		context map[string]string
		why     string
	}{
		{
			name:    "unbound_tag",
			context: nil,
			why: "an unbound `{tag.<key>}` is a property of the REQUEST: no " +
				"repair of the environment can bind a tag the caller did not pass",
		},
		{
			name:    "flag_shaped_value",
			context: map[string]string{"nnnn": "--version"},
			why: "a `-`-prefixed bound value is a property of the REQUEST: the " +
				"caller must pass a different value, not re-run this one",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// `false` as argv0 exits non-zero, so a refusal that did NOT
			// happen before spawn would surface the child's error instead
			// of the placeholder's — the same guard the Phase 1 fixture
			// uses.
			acc := table.Accessor{
				Role:    editRole,
				Command: []string{"false", "{tag.nnnn}"},
				Keys:    []string{editKey},
				Timeout: "2s",
			}
			r := cmdbind.Reader{
				Accessor: acc,
				Name:     editReader,
				Config:   cmdbind.Config{AllowCommands: true},
			}

			art := accessor.Artifact{Role: editRole, Path: "/dev/null", Context: tc.context}
			_, _, err := r.Read(context.Background(), art, []string{editKey})
			if err == nil {
				t.Fatalf("the read succeeded; `0028:C1.6` requires a refusal " +
					"before spawn")
			}

			// The Phase 1 assertion, restated so a regression that drops
			// the Detail is caught here too rather than only there.
			if detail := detailOf(t, err); detail == "" {
				t.Errorf("the refusal carries an empty Detail; it must NAME " +
					"the placeholder, having no rule to name")
			}

			// The finding. The typed `Err` is the only carrier
			// `0028:C1.3` gives `accessorFailureOf` for the exit group.
			if !errors.Is(err, accessor.ErrDeclaredRequest) {
				t.Fatalf("the refusal does not wrap "+
					"`accessor.ErrDeclaredRequest`, so `refusalOf` records "+
					"DeclaredRequest() = false and `accessorFailureOf` routes "+
					"it to the environment exit group (exit 3) rather than the "+
					"exit-2 group `0028:C1.3` EXIT GROUP: fixes for the "+
					"entry-level preconditions.\n%s\ngot err = %v", tc.why, err)
			}
		})
	}

	t.Run("bound_non_flag_value_is_not_refused", func(t *testing.T) {
		// The negative control. The fix narrows to the two C1.6 arms and
		// changes no admission decision.
		acc := table.Accessor{
			Role:    editRole,
			Command: []string{"echo", "{tag.nnnn}"},
			Output:  ptr("raw"),
			Keys:    []string{editKey},
			Timeout: "2s",
		}
		r := cmdbind.Reader{
			Accessor: acc,
			Name:     editReader,
			Config:   cmdbind.Config{AllowCommands: true},
		}
		art := accessor.Artifact{
			Role: editRole, Path: "/dev/null",
			Context: map[string]string{"nnnn": "0028"},
		}
		if _, _, err := r.Read(context.Background(), art, []string{editKey}); err != nil {
			t.Fatalf("a bound non-flag value refused: %v", err)
		}
	})

	t.Run("the_0025_gate_refusal_keeps_its_own_routing", func(t *testing.T) {
		// The scope control. `0025:C6`'s gate refusal is NOT one of the
		// three entry-level preconditions `0028:C1.3` enumerates — it is
		// 0025's own pre-spawn ladder, which this record does not amend.
		// The fix is a sibling of `refuse`, not a widening of it, and
		// this arm is what holds that line.
		acc := table.Accessor{
			Role:    editRole,
			Command: []string{"echo", "{tag.nnnn}"},
			Output:  ptr("raw"),
			Keys:    []string{editKey},
			Timeout: "2s",
		}
		r := cmdbind.Reader{
			Accessor: acc,
			Name:     editReader,
			Config:   cmdbind.Config{AllowCommands: false},
		}
		art := accessor.Artifact{
			Role: editRole, Path: "/dev/null",
			Context: map[string]string{"nnnn": "0028"},
		}
		_, _, err := r.Read(context.Background(), art, []string{editKey})
		if err == nil {
			t.Fatalf("the gate-off read succeeded")
		}
		if errors.Is(err, accessor.ErrDeclaredRequest) {
			t.Errorf("`0025:C6`'s gate refusal now wraps " +
				"`accessor.ErrDeclaredRequest`; RDR 0028 amends only C1.6's " +
				"two argv rules and leaves 0025's pre-spawn ladder alone")
		}
	})
}
