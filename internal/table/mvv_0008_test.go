package table_test

// RDR 0008's gating Minimum Viable Validation, as one runnable end-to-end
// test.
//
// The RDR declares NO Round-Trip / Inverse Invariant (no encode/decode,
// import/export, or inverse operation is introduced), so the MVV's obligation
// is value-for-value assertion of each half's stated outcome rather than a
// reconstruct-and-compare round trip. A green exit code or "did not error" is
// explicitly not sufficient: every leg below compares values.
//
// It lives in `internal/table` because only this package can reach both ends
// — the loader for the normalizer half, and `internal/resolve` (which it
// already imports) for the kernel half. The two halves are the record's own
// split and are asserted together so neither can be declared done alone.

import (
	"testing"

	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-MVV / REQ-72: `0008:MVV` "**Kernel half — executed during this RDR's
// implementation.** A row matching on the tag `recognized` fires against a
// recognized outcome, asserted through the real kernel's assembled view (the
// behavioral conformance form, premortem P-5), and an `Input` supplying an
// owned or observed tag keyed `recognized` — or a row naming it in
// `RequiresOwned` — is rejected by the exported predicate and by `Resolve` at
// entry. Both run against `internal/resolve` as it stands."
// REQ-MVV / REQ-73: `0008:MVV` "**Normalizer half — executed inside RDR 0002's
// implementation.** A table declaring a recognized-provenance tag named
// `recognized` loads and lints clean; the same table with the declaration
// renamed (and separately, with an owned tag named `recognized`) fails
// load/lint in the `reserved_tag_key` category whose failure data carries the
// direction-appropriate remedy name and rule identifier — not a silent
// no-match at resolve time."
// REQ-74: TS-1 "**Expected**: the row is selected; the recognized outcome is
// readable at key `recognized` with `ProvenanceRecognized`."
// HAPPY PATH
func TestMVV0008_ReservedKeyOwnershipEndToEnd(t *testing.T) {
	const reserved = "recognized"

	t.Run("1_normalizer_half_clean_load_under_the_reserved_name", func(t *testing.T) {
		m := mustLoad(t, rdrFixture)

		decl, ok := m.Tags[reserved]
		if !ok {
			t.Fatalf("%s declares no %q tag", rdrFixture, reserved)
		}
		if decl.Provenance != table.ProvenanceRecognized {
			t.Fatalf("%q declaration provenance = %q; want %q",
				reserved, decl.Provenance, table.ProvenanceRecognized)
		}
		// Value-for-value: no OTHER declaration takes the reserved name, and
		// no recognized-provenance declaration sits under another name.
		for key, d := range m.Tags {
			if key == reserved && d.Provenance != table.ProvenanceRecognized {
				t.Errorf("%q is declared with provenance %q", key, d.Provenance)
			}
			if key != reserved && d.Provenance == table.ProvenanceRecognized {
				t.Errorf("recognized provenance declared under the name %q", key)
			}
		}
	})

	t.Run("2_normalizer_half_both_directions_fail_with_their_payload", func(t *testing.T) {
		directions := map[string]struct {
			fixture   string
			offending string
			remedy    string
			rule      string
		}{
			"declaration renamed away from the reserved key": {
				fixture:   "neg/neg-recognized-misnamed.toml",
				offending: "outcome",
				remedy:    reserved,
				rule:      "reserved-tag-key/kernel-owned",
			},
			"owned declaration taking the reserved key": {
				fixture:   "neg/neg-recognized-owned.toml",
				offending: reserved,
				remedy:    "",
				rule:      "reserved-tag-key/author-must-rename",
			},
		}

		for name, d := range directions {
			t.Run(name, func(t *testing.T) {
				f := reservedTagKeyFailure(t, d.fixture)
				if f.Offending != d.offending {
					t.Errorf("offending = %q; want %q", f.Offending, d.offending)
				}
				if f.Remedy != d.remedy {
					t.Errorf("remedy = %q; want %q", f.Remedy, d.remedy)
				}
				if f.Rule != d.rule {
					t.Errorf("rule = %q; want %q", f.Rule, d.rule)
				}
			})
		}
	})

	t.Run("3_kernel_half_a_row_matching_the_reserved_key_fires", func(t *testing.T) {
		in := mvvKernelInput()

		got, err := resolve.Resolve(in)
		if err != nil {
			t.Fatalf("the conforming tuple traveled the Go error path: %v", err)
		}
		if got.Refusal != nil {
			t.Fatalf("resolve refused %q; the row matching on %q must fire",
				got.Refusal.Kind, reserved)
		}
		if got.Plan == nil {
			t.Fatal("no disposition")
		}
		// Value-for-value on the whole plan, not a bare non-nil check.
		if got.Plan.RuleID != "mvv.fires" {
			t.Errorf("RuleID = %q; want %q", got.Plan.RuleID, "mvv.fires")
		}
		if got.Plan.Revision != "rev-mvv-0008" {
			t.Errorf("Revision = %q; want %q", got.Plan.Revision, "rev-mvv-0008")
		}
		if len(got.Plan.NextTags) != 1 ||
			got.Plan.NextTags[0] != (resolve.Tag{Key: "stage", Value: "prelock"}) {
			t.Errorf("NextTags = %+v; want [{stage prelock}]", got.Plan.NextTags)
		}
		if len(got.Plan.Writes) != 1 ||
			got.Plan.Writes[0] != (resolve.Tag{Key: "status", Value: "Final"}) {
			t.Errorf("Writes = %+v; want [{status Final}]", got.Plan.Writes)
		}
		if got.Plan.Escaped {
			t.Error("the row was selected as an escape rescue, not an ordinary match")
		}
	})

	t.Run("4_kernel_half_the_recognized_outcome_is_readable_at_the_reserved_key",
		func(t *testing.T) {
			// Behavioral, per premortem P-5: the binding is observed through
			// the real kernel's selection, never by comparing two literals.
			// A row demanding the OUTCOME's value at the reserved key fires
			// only if the view binds exactly (Input.Recognized) there; a
			// sibling row demanding any other value at the same key must not.
			in := mvvKernelInput()
			in.Table.Rows = append(in.Table.Rows, resolve.Row{
				RuleID:        "mvv.wrong-value",
				SourceLocator: "flows/rdr.toml:2",
				Outcome:       "round-clean",
				Match: []resolve.Tag{
					{Key: "status", Value: "Draft"},
					{Key: reserved, Value: "some-other-outcome"},
				},
				NextTags: []resolve.Tag{{Key: "stage", Value: "wrong"}},
			})

			got, err := resolve.Resolve(in)
			if err != nil {
				t.Fatalf("Go error path: %v", err)
			}
			if got.Plan == nil {
				t.Fatalf("no plan; disposition %+v", got)
			}
			if got.Plan.RuleID != "mvv.fires" {
				t.Fatalf("selected %q; the view binds %q to something other than "+
					"Input.Recognized = %q", got.Plan.RuleID, reserved, in.Recognized)
			}

			// PRESENCE, observed behaviorally: an existence atom over the
			// reserved key is decided by the kernel without consulting the
			// value seam (`0007:C1`), so a row guarded on `exists = true`
			// survives only if the key is in the view. With a nil seam, a row
			// guarded on `exists = false` over the same key is decided FALSE
			// and pruned — it cannot reach a plan.
			present := mvvKernelInput()
			present.Table.Rows[0].Guard = []resolve.GuardAtom{{
				Key:      reserved,
				Operator: resolve.OpExists,
				Literal:  resolve.LiteralTrue,
				Block:    resolve.BlockAll,
			}}
			pGot, pErr := resolve.Resolve(present)
			if pErr != nil {
				t.Fatalf("Go error path: %v", pErr)
			}
			if pGot.Plan == nil || pGot.Plan.RuleID != "mvv.fires" {
				t.Errorf("a row guarded on `%s exists = true` did not fire; the "+
					"key is absent from the assembled view: %+v", reserved, pGot)
			}

			absent := mvvKernelInput()
			absent.Table.Rows[0].Guard = []resolve.GuardAtom{{
				Key:      reserved,
				Operator: resolve.OpExists,
				Literal:  resolve.LiteralFalse,
				Block:    resolve.BlockAll,
			}}
			aGot, aErr := resolve.Resolve(absent)
			if aErr != nil {
				t.Fatalf("Go error path: %v", aErr)
			}
			if aGot.Plan != nil {
				t.Errorf("a row guarded on `%s exists = false` fired; the key "+
					"must be present in the view", reserved)
			}
		})

	t.Run("5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites",
		func(t *testing.T) {
			channels := map[string]func(*resolve.Input){
				"owned tag keyed recognized": func(in *resolve.Input) {
					in.Owned = append(in.Owned, resolve.Tag{Key: reserved, Value: "x"})
				},
				"observed tag keyed recognized": func(in *resolve.Input) {
					in.Observed = append(in.Observed, resolve.Tag{Key: reserved, Value: "x"})
				},
				"row naming it in RequiresOwned": func(in *resolve.Input) {
					in.Table.Rows[0].RequiresOwned =
						append(in.Table.Rows[0].RequiresOwned, reserved)
				},
			}

			for name, breach := range channels {
				t.Run(name, func(t *testing.T) {
					in := mvvKernelInput()
					breach(&in)

					// Site 1: the exported predicate producers may call.
					if err := resolve.CheckInput(in); err == nil {
						t.Error("the exported predicate accepted the breach")
					}

					// Site 2: Resolve at entry — a non-nil error and NO
					// disposition, asserted value-for-value against the zero
					// Result rather than by a bare "did not error".
					got, err := resolve.Resolve(in)
					if err == nil {
						t.Fatal("Resolve did not apply the predicate at entry")
					}
					if (got != resolve.Result{}) {
						t.Errorf("Result = %+v; want the zero value — no Plan and "+
							"no Refusal accompany a producer breach", got)
					}
				})
			}
		})

	t.Run("6_the_two_halves_agree_on_the_spelling", func(t *testing.T) {
		// The declaration channel's reserved name and the key the kernel binds
		// are the same string. This is what makes the two halves one contract
		// rather than two coincidences.
		m := mustLoad(t, rdrFixture)
		if _, ok := m.Tags[table.RecognizedTagKey]; !ok {
			t.Fatalf("the loader's reserved name %q is not the declared key",
				table.RecognizedTagKey)
		}

		// Behavioral: a row whose Match uses the LOADER's spelling selects the
		// row the kernel bound. If the two spellings diverged, no candidate
		// would match and the resolve would refuse.
		in := mvvKernelInput()
		in.Table.Rows[0].Match = []resolve.Tag{
			{Key: "status", Value: "Draft"},
			{Key: table.RecognizedTagKey, Value: in.Recognized},
		}
		got, err := resolve.Resolve(in)
		if err != nil {
			t.Fatalf("Go error path: %v", err)
		}
		if got.Plan == nil || got.Plan.RuleID != "mvv.fires" {
			t.Errorf("the kernel binds the outcome somewhere other than the "+
				"loader's reserved name %q; disposition %+v",
				table.RecognizedTagKey, got)
		}
	})
}

// mvvKernelInput is the MVV's kernel-half tuple: one ordinary row whose Match
// names the reserved key, over a conforming owned/observed pair.
func mvvKernelInput() resolve.Input {
	return resolve.Input{
		Flow: "rdr",
		Table: resolve.Table{
			Revision: "rev-mvv-0008",
			Outcomes: []string{"round-clean"},
			Rows: []resolve.Row{{
				RuleID:        "mvv.fires",
				SourceLocator: "flows/rdr.toml:1",
				Outcome:       "round-clean",
				Match: []resolve.Tag{
					{Key: "status", Value: "Draft"},
					{Key: "recognized", Value: "round-clean"},
				},
				RequiresOwned: []string{"status"},
				NextTags:      []resolve.Tag{{Key: "stage", Value: "prelock"}},
				Writes:        []resolve.Tag{{Key: "status", Value: "Final"}},
			}},
		},
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Observed:   []resolve.Tag{{Key: "reviews", Value: "2"}},
		Recognized: "round-clean",
	}
}
