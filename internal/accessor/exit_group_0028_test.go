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
	"path/filepath"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/cli/cmdbind"
	"github.com/cwensel/intrastate/internal/cli/flowbind"
	"github.com/cwensel/intrastate/internal/resolve"
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

	t.Run("the_0025_gate_refusal_takes_the_request_exit_group", func(t *testing.T) {
		// `0025:C6`'s gate refusal is a missing opt-in: a property of the
		// request, which re-issuing unchanged can never satisfy. JDR 0001
		// §D10 rule 1 reserves exit 3 for an environment that could not be
		// consulted, and `0028:C1.3` already routes the same missing flag
		// on the edit read-back to exit 2, so the gate wraps
		// `ErrDeclaredRequest` too. It is not a line edit, so it must not
		// carry the narrower `ErrDeclaredEdit` bit.
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
		if !errors.Is(err, accessor.ErrDeclaredRequest) {
			t.Errorf("`0025:C6`'s gate refusal does not wrap "+
				"`accessor.ErrDeclaredRequest`, so it routes to exit 3 and a "+
				"caller's retry loop spins on a request it must fix: %v", err)
		}
		if errors.Is(err, accessor.ErrDeclaredEdit) {
			t.Errorf("`0025:C6`'s gate refusal wraps "+
				"`accessor.ErrDeclaredEdit`; no line edit minted it: %v", err)
		}
	})
}

// --- ADV-3's re-anchor widening: the arms its fixture does not reach ------

// REQ-38/REQ-39 (`0028:C1.3` re-anchor:).
//
// ADV-3 pins the poisoning arm: a planned rule's replacement makes an
// UNPLANNED sibling's anchor select a line it did not select before, and
// the write must refuse before any byte lands. These are the arms around
// it, which the widening must not break:
//
//   - an unplanned rule whose anchor is untouched by the rewrite still
//     applies cleanly (the widening is not a blanket refusal);
//   - an unplanned rule whose anchor was ALREADY stale or ALREADY
//     ambiguous before this invocation stays that way and is not this
//     plan's refusal to report — the invariant is STABILITY, not
//     cardinality, because an unplanned rule has no rewritten line of
//     its own to be identical to;
//   - a `<clear>` deletion shifts an unplanned rule's line without
//     changing what it selects, so the shift arithmetic must follow it;
//   - a rewrite that makes an unplanned rule's anchor select NOTHING is
//     poisoning too, in the other direction.
func TestReAnchor0028_UnplannedSiblingRulesAreCovered(t *testing.T) {
	const statusAnchor = `^- \*\*Status\*\*: .+$`

	t.Run("untouched_sibling_applies", func(t *testing.T) {
		before := "- **Status**: Draft\n- **Owner**: alice\n"
		path := adv0028Fixture(t, before)
		acc := adv0028Accessor(map[string]table.EditRule{
			"status": {Anchor: statusAnchor, Replace: "- **Status**: {status}"},
			"owner":  {Anchor: `^- \*\*Owner\*\*: .+$`, Replace: "- **Owner**: {owner}"},
		}, "status", "owner")

		err := flowbind.NewEditWriter(acc, adv0028Entry).Apply(
			t.Context(),
			accessor.Artifact{Role: adv0028Role, Path: path},
			[]resolve.Tag{{Key: "status", Value: "Final"}},
		)
		if err != nil {
			t.Fatalf("a rewrite that leaves the unplanned sibling's anchor "+
				"selecting exactly the same line refused: %v", err)
		}
		if got := string(adv0028Read(t, path)); got != "- **Status**: Final\n- **Owner**: alice\n" {
			t.Fatalf("the planned rule did not apply: %q", got)
		}
	})

	t.Run("already_unmatched_sibling_is_not_this_plans_refusal", func(t *testing.T) {
		// `owner` matched nothing BEFORE the edit and matches nothing
		// after it. Nothing changed, so nothing is reported: the stale
		// sibling is the refusal of the next plan that NAMES it.
		before := "- **Status**: Draft\n"
		path := adv0028Fixture(t, before)
		acc := adv0028Accessor(map[string]table.EditRule{
			"status": {Anchor: statusAnchor, Replace: "- **Status**: {status}"},
			"owner":  {Anchor: `^- \*\*Owner\*\*: .+$`, Replace: "- **Owner**: {owner}"},
		}, "status", "owner")

		err := flowbind.NewEditWriter(acc, adv0028Entry).Apply(
			t.Context(),
			accessor.Artifact{Role: adv0028Role, Path: path},
			[]resolve.Tag{{Key: "status", Value: "Final"}},
		)
		if err != nil {
			t.Fatalf("an unplanned sibling that was ALREADY unmatched before "+
				"this invocation was reported as this plan's defect: %v", err)
		}
	})

	t.Run("clear_deletion_shifts_the_sibling_without_poisoning_it", func(t *testing.T) {
		// The planned `note` rule DELETES its line, which moves the
		// unplanned `owner` rule's line up by one. The set it selects is
		// unchanged, so the shift arithmetic must follow the deletion
		// rather than read the move as a change.
		before := "- **Note**: scratch\n- **Owner**: alice\n"
		path := adv0028Fixture(t, before)
		acc := adv0028Accessor(map[string]table.EditRule{
			"note": {
				Anchor:  `^- \*\*Note\*\*: .+$`,
				Replace: "- **Note**: {note}",
				Clear:   table.EditClearLine,
			},
			"owner": {Anchor: `^- \*\*Owner\*\*: .+$`, Replace: "- **Owner**: {owner}"},
		}, "note", "owner")

		err := flowbind.NewEditWriter(acc, adv0028Entry).Apply(
			t.Context(),
			accessor.Artifact{Role: adv0028Role, Path: path},
			[]resolve.Tag{{Key: "note", Value: accessor.ClearSentinel}},
		)
		if err != nil {
			t.Fatalf("a deletion that merely SHIFTED the unplanned sibling's "+
				"line was read as poisoning it: %v", err)
		}
		if got := string(adv0028Read(t, path)); got != "- **Owner**: alice\n" {
			t.Fatalf("the `<clear>` did not delete its line: %q", got)
		}
	})

	t.Run("rewrite_that_de_anchors_a_sibling_refuses_before_mutation", func(t *testing.T) {
		// Poisoning in the other direction: the replacement REMOVES the
		// text the unplanned sibling anchored on, so its anchor now
		// selects nothing. The next plan naming `owner` would refuse
		// `edit_anchor_unmatched` with no way to trace the run that
		// broke it, which is what this pass exists to prevent.
		before := "- **Status**: Draft alice\n"
		path := adv0028Fixture(t, before)
		acc := adv0028Accessor(map[string]table.EditRule{
			"status": {Anchor: statusAnchor, Replace: "- **Status**: {status}"},
			"owner":  {Anchor: `alice`, Replace: "- **Owner**: {owner}"},
		}, "status", "owner")

		err := flowbind.NewEditWriter(acc, adv0028Entry).Apply(
			t.Context(),
			accessor.Artifact{Role: adv0028Role, Path: path},
			[]resolve.Tag{{Key: "status", Value: "Final"}},
		)
		if err == nil {
			t.Fatalf("the rewrite de-anchored the unplanned sibling `owner` " +
				"and applied anyway; `0028:C1.3` re-anchor: requires " +
				"`edit_anchor_unstable` before any write")
		}
		if got := string(adv0028Read(t, path)); got != before {
			t.Fatalf("the edit refused but the artifact changed: %q", got)
		}
	})
}

// REQ-45: "That arm stays 0004's, unchanged and out of this RDR's scope
// (A11, `Verified`): this clause's guarantee is scoped to refusals the
// `edit` binding itself mints, and does NOT assert the executor-minted
// timeout carries a false applied sense." — (0028:C1.3).
//
// BOUNDARY. This is the NEGATIVE half of the exit-group clause and the
// stopping point of the ADV-1/ADV-2 fix. Those widened
// `ErrDeclaredRequest` to the ENTRY-LEVEL preconditions; REQ-45 says the
// executor's own deadline arm is NOT one of them. Without a test the fix
// has no stated boundary, and a later author "finishing" it by wrapping
// the timeout arm would reclassify a genuine ENVIRONMENT failure as a
// request defect — telling a caller whose command merely ran slow that
// re-running the same request cannot help, which is the opposite of true.
//
// The applied-but-unverified sense this REQ also declines to disturb is
// already pinned by `TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense`
// (0004). What is asserted HERE is only what THIS record could break: the
// sentinel it introduced must not reach 0004's arm.
func TestReq45_TheExecutorTimeoutArmDoesNotWrapTheRequestSentinel(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	// The mutation lands, THEN the invocation runs past its deadline —
	// 0004's post-mutation timeout, reached without any `edit` carrier.
	w := &writeBinding{store: s, delay: 10 * fixtureTimeout}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	ref := mustWriteRefuse(t, got, accessor.ClassTimeout)
	if ref.DeclaredRequest() {
		t.Errorf("the executor-minted timeout reports DeclaredRequest() = " +
			"true, so `accessorFailureOf` routes it to the exit-2 REQUEST " +
			"group. REQ-45 scopes this arm OUT of `0028:C1.3`: a command " +
			"that ran past its deadline is an ENVIRONMENT failure and keeps " +
			"exit 3, where re-running unchanged is the honest advice.")
	}
}

// REQ-51 / `0028:C1.3` precedence: step (1), order:.
//
// TRIAGE REGRESSION (roborev job 6783). C1.3 `precedence:` puts C1.6's two
// ARGV preconditions — an unbound `{tag.<key>}`, then a `-`-prefixed bound
// value — first, "decided together before anything is spawned or read",
// and `order:` requires "every refusal in this clause … decided BEFORE any
// byte is written".
//
// The gate pre-check honoured that only for the GATE. When the gate is ON,
// the read-back reader's own argv placeholders were never examined before
// `Apply`: the edit landed, then `invokeRead` refused during read-back, and
// the caller got `read_back_incomplete` with the applied-but-unverified
// sense for what is purely a defect of the request. The artifact was
// mutated where the contract says it must not be.
//
// Deviation D11 does NOT cover this. D11 reasoned over the gate-OFF path,
// where no child is spawned and C1.6's refusals are unreachable. These
// cases require the gate ON — a passed gate is exactly what spawns the
// child that would raise them.
//
// ADVERSARIAL. Both arms assert the artifact is byte-identical: the point
// is not merely which error is reported, but that no write occurred.
func TestReq51_ReadBackReaderArgvPreconditionsRefuseBeforeMutation(t *testing.T) {
	const before = "- **Status**: Draft\n"

	for _, tc := range []struct {
		name    string
		context map[string]string
		why     string
	}{
		{
			name:    "unbound_tag",
			context: nil,
			why: "an unbound `{tag.<key>}` on the read-back reader's argv is " +
				"a property of the REQUEST; C1.3 precedence: step (1) decides " +
				"it before anything is spawned or read",
		},
		{
			name:    "flag_shaped_value",
			context: map[string]string{"nnnn": "--version"},
			why: "a `-`-prefixed bound value is the second half of step (1) " +
				"and is decided in the same place",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "record.md")
			adv0028Write(t, path, before)

			model := &table.Model{
				ID:   "req51",
				Tags: map[string]table.TagDecl{"status": {Provenance: table.ProvenanceOwned}},
				Readers: map[string]table.Accessor{
					"req51-reader": {
						Role:     adv0028Role,
						Keys:     []string{"status"},
						Command:  []string{"/bin/echo", "{tag.nnnn}"},
						Timeout:  "1s",
						ReadBack: true,
					},
				},
				Writers: map[string]table.Accessor{
					adv0028Entry: {
						Role:     adv0028Role,
						Keys:     []string{"status"},
						Timeout:  "1s",
						ReadBack: true,
						Edit: map[string]table.EditRule{
							// A STATIC anchor: the writer itself has no
							// unbound placeholder, so the only precondition
							// in play is the READER's argv.
							"status": {
								Anchor:  `^- \*\*Status\*\*: .+$`,
								Replace: "- **Status**: {status}",
							},
						},
					},
				},
			}

			// The gate is ON. That is what makes the C1.6 refusal
			// reachable at all, and what distinguishes this from D11.
			reg := flowbind.Registry(model, dir, true)
			exec := accessor.NewExecutor(reg, accessor.Artifacts{
				adv0028Role: {Role: adv0028Role, Path: path, Context: tc.context},
			})

			got := exec.Write(t.Context(), adv0028Entry, resolve.Plan{
				Writes: []resolve.Tag{{Key: "status", Value: "Final"}},
			})

			if got.Refusal == nil {
				t.Fatalf("Write succeeded; %s", tc.why)
			}
			if body := string(adv0028Read(t, path)); body != before {
				t.Fatalf("the artifact was MUTATED before the precondition "+
					"refused: %q. `0028:C1.3` order: requires every refusal "+
					"in this clause to be decided BEFORE any byte is "+
					"written.\n%s", body, tc.why)
			}
			if got.Refusal.Applied() {
				t.Errorf("the refusal carries the applied-but-unverified " +
					"sense; nothing was applied, and `0028:C1.3` order: " +
					"forbids that sense for a pre-write refusal")
			}
			if !got.Refusal.DeclaredRequest() {
				t.Errorf("the refusal is not a declared-request failure, so " +
					"it takes exit 3's \"re-run unchanged\" advice for a " +
					"defect no re-run can repair")
			}
		})
	}
}

// REQ-51 second leg / `0028:C1.3` precedence: step (1).
//
// TRIAGE REGRESSION (roborev job 6807). The clause orders the entry-level
// preconditions "C1.6's unbound `{tag.<key>}`, THEN C1.6's `-`-prefixed
// bound value" — all unbound keys ahead of any prefixed value. A single
// argv walk reports whichever defect sits earlier in the vector, so an
// argv whose prefixed value PRECEDES its unbound key inverts the
// precedence the clause fixes.
//
// Observable only through a multi-defect input, which the clause says
// outright: "Fail-fast here is observable only through a multi-defect
// input … so it earns its own scenario rather than riding the
// single-defect fixtures."
//
// BOUNDARY. Both arms carry the SAME two defects and differ only in argv
// order; both must report the unbound key.
func TestReq51_UnboundTagOutranksAPrefixedValueRegardlessOfArgvOrder(t *testing.T) {
	for _, tc := range []struct {
		name string
		argv []string
	}{
		{
			name: "unbound_first",
			argv: []string{"/bin/echo", "{tag.missing}", "{tag.flagged}"},
		},
		{
			// The discriminating arm: a single walk reports the
			// prefixed value here and inverts the precedence.
			name: "prefixed_first",
			argv: []string{"/bin/echo", "{tag.flagged}", "{tag.missing}"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "record.md")
			adv0028Write(t, path, "- **Status**: Draft\n")

			model := &table.Model{
				ID:   "req51order",
				Tags: map[string]table.TagDecl{"status": {Provenance: table.ProvenanceOwned}},
				Readers: map[string]table.Accessor{
					"req51order-reader": {
						Role:     adv0028Role,
						Keys:     []string{"status"},
						Command:  tc.argv,
						Timeout:  "1s",
						ReadBack: true,
					},
				},
				Writers: map[string]table.Accessor{
					adv0028Entry: {
						Role:     adv0028Role,
						Keys:     []string{"status"},
						Timeout:  "1s",
						ReadBack: true,
						Edit: map[string]table.EditRule{
							"status": {
								Anchor:  `^- \*\*Status\*\*: .+$`,
								Replace: "- **Status**: {status}",
							},
						},
					},
				},
			}

			reg := flowbind.Registry(model, dir, true)
			exec := accessor.NewExecutor(reg, accessor.Artifacts{
				adv0028Role: {
					Role: adv0028Role, Path: path,
					// `flagged` is bound to a `-`-prefixed value;
					// `missing` is not bound at all.
					Context: map[string]string{"flagged": "--version"},
				},
			})

			got := exec.Write(t.Context(), adv0028Entry, resolve.Plan{
				Writes: []resolve.Tag{{Key: "status", Value: "Final"}},
			})

			if got.Refusal == nil {
				t.Fatalf("Write succeeded with two entry-level precondition " +
					"defects on the read-back reader's argv")
			}
			if !strings.Contains(got.Refusal.Detail, "missing") {
				t.Errorf("the refusal names %q; want the UNBOUND key "+
					"`missing`. `0028:C1.3` precedence: step (1) orders "+
					"every unbound `{tag.<key>}` ahead of any `-`-prefixed "+
					"bound value, so argv position must not decide which "+
					"defect is reported.\nDetail = %s",
					"flagged", got.Refusal.Detail)
			}
		})
	}
}
