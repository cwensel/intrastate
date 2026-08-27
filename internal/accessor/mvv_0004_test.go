package accessor_test

// RDR 0004 — the Minimum Viable Validation.

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-MVV (`0004:MVV`): "Build a fixture flow with one read accessor, one
// gate accessor, and one write accessor over caller-supplied artifact
// roles. Prove success, timeout, gate denied, gate-indeterminate,
// execution failure, incomplete read, capability mismatch, unsafe
// definition validation, write read-back-mismatch, and
// `read_back_incomplete` dispositions, plus the two A9 boundary cases: an
// absent required key refused as `owned_state_unavailable` at the
// resolver, and a timed-out partial read classified as `timeout`; plus
// A10's case: a write whose re-read cannot read a compared key surfaces
// the applied-but-unverified sense and the test asserts no second write
// and no compensating action occurred; plus A11's cases: a clearing write
// whose re-read shows the key absent, a clear of a key the artifact does
// not hold succeeding, a clear whose re-read still holds the key
// (including as the literal) refused as `read_back_mismatch`, and a read
// holding the literal `<clear>` refused as `incomplete_read`. The write
// success test must assert the re-read owned-tag value equals the
// transition plan's expected value. The read test must assert that an
// accessor which can resolve only some of the requested keys takes the
// refusal branch and is distinguishable from one whose artifact genuinely
// lacks those keys."
// HAPPY PATH
//
// This is the gating validation for RDR 0004: ONE fixture flow — one
// read accessor, one gate accessor, one write accessor, all over
// caller-supplied artifact roles — carried through the nine numbered
// scenarios (TS 1-9), each with its named negative control (ORA).
//
// RDR 0004 DOES declare a Round-Trip / Inverse Invariant:
//
//	write -> read = expected owned-tag value identity
//	              + protected non-owned tag identity
//
// so scenario 3 reconstructs the written state through the read path and
// compares it VALUE-FOR-VALUE against the transition plan's expected
// values — not "did not error", and not a green exit code. Scenario 9
// extends the same invariant's expected-value predicate to ABSENCE for a
// planned `<clear>`.
//
// The Fidelity Table's deliberate weakening is honoured: the invariant is
// asserted at TAG-VALUE granularity, never over the artifact encoding.
func TestMVV_AccessorExecutionSafetyModel(t *testing.T) {
	// ---- the fixture flow --------------------------------------------
	//
	// One read accessor (`state.read`), one gate accessor
	// (`state.gate`), one write accessor (`state.persist`), all bound to
	// the caller-supplied `state` artifact role. Nothing is discovered
	// from ambient process state.

	t.Run("1_unsafe_definition_validation_eight_arms", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		valid := func() (accessor.Registry, []accessor.Identity) {
			rd := readerDef(&readBinding{store: s}, keyStatus, keyProfile)
			gd := gateDef(&gateBinding{verdict: accessor.VerdictAllow})
			wd := writerDef(&writeBinding{store: s}, keyStatus)
			return registryOf(rd, gd, wd),
				[]accessor.Identity{rd.Identity, gd.Identity, wd.Identity}
		}

		// The negative control FIRST: a validator that rejects everything
		// fails here before any arm can pass.
		reg, ids := valid()
		if fs := accessor.Validate(reg, ids); len(fs) != 0 {
			t.Fatalf("validation-ok = %v; want the EMPTY set for the well-formed "+
				"fixture flow — a validator that rejects everything must fail here",
				findingCodes(fs))
		}

		arms := []struct {
			name   string
			break_ func(*accessor.Registry, *[]accessor.Identity)
			want   accessor.ValidationCode
		}{
			{"missing_accessor", func(r *accessor.Registry, i *[]accessor.Identity) {
				*i = append(*i, accessor.Identity{
					Flow: flowID, Name: "absent.read", Capability: accessor.CapRead})
			}, accessor.CodeMissingAccessor},
			{"multiply_bound_accessor", func(r *accessor.Registry, i *[]accessor.Identity) {
				r.Definitions = append(r.Definitions,
					readerDef(&readBinding{store: s}, keyStatus, keyProfile))
			}, accessor.CodeMultiplyBoundAccessor},
			{"capability_mismatch", func(r *accessor.Registry, i *[]accessor.Identity) {
				*i = append(*i, accessor.Identity{
					Flow: flowID, Name: readerName, Capability: accessor.CapWrite})
			}, accessor.CodeCapabilityMismatch},
			{"missing_or_non_positive_timeout", func(r *accessor.Registry, i *[]accessor.Identity) {
				r.Definitions[0].Accessor.Timeout = "0s"
			}, accessor.CodeMissingOrNonPositiveTimeout},
			{"missing_write_read_back", func(r *accessor.Registry, i *[]accessor.Identity) {
				for k := range r.Definitions {
					if r.Definitions[k].Identity.Capability == accessor.CapWrite {
						r.Definitions[k].Accessor.ReadBack = false
					}
				}
			}, accessor.CodeMissingWriteReadBack},
			{"ambient_artifact_discovery", func(r *accessor.Registry, i *[]accessor.Identity) {
				r.Definitions[0].AmbientDiscovery = true
			}, accessor.CodeAmbientArtifactDiscovery},
			{"write_non_owned_tag", func(r *accessor.Registry, i *[]accessor.Identity) {
				for k := range r.Definitions {
					if r.Definitions[k].Identity.Capability == accessor.CapWrite {
						r.Definitions[k].Accessor.Keys = []string{keyProfile}
					}
				}
			}, accessor.CodeWriteNonOwnedTag},
			// The eighth arm — A9's first rule, unwitnessed by the spike.
			{"missing_requested_key_set", func(r *accessor.Registry, i *[]accessor.Identity) {
				r.Definitions[0].Accessor.Keys = nil
			}, accessor.CodeMissingRequestedKeySet},
		}

		if len(arms) != 8 {
			t.Fatalf("the MVV carries %d validation arms; the record names EIGHT", len(arms))
		}
		for _, arm := range arms {
			t.Run(arm.name, func(t *testing.T) {
				reg, ids := valid()
				arm.break_(&reg, &ids)

				fs := accessor.Validate(reg, ids)

				if !hasCode(fs, arm.want) {
					t.Errorf("validation-%s = %v; want the arm's OWN named code %q",
						arm.name, findingCodes(fs), arm.want)
				}
			})
		}
	})

	t.Run("2_read_and_gate_dispositions", func(t *testing.T) {
		// The requested key set is PINNED in the definition, independently
		// of what any artifact carries. ORA 2's second control: a
		// derived-key-set implementation returns a one-key SUCCESS for the
		// truncation arm instead of refusing.
		requested := []string{keyStatus, keyProfile}

		t.Run("complete_read_returns_exactly_the_requested_keys", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			e, _ := readerExec(t, s, requested...)

			values := mustReadSucceed(t, e.Read(ctxOf(t), readerName))

			want := slices.Clone(requested)
			slices.Sort(want)
			if got := keysOf(values); !slices.Equal(got, want) {
				t.Fatalf("read returned keys %v; want exactly %v", got, want)
			}
			for key, wantValue := range map[string]string{
				keyStatus: "Draft", keyProfile: "large",
			} {
				v, _ := valueOf(values, key)
				if v.Value != wantValue || v.Absent {
					t.Errorf("value for %q = %+v; want %q, present", key, v, wantValue)
				}
			}
		})

		t.Run("incomplete_read_refuses_with_no_values_and_names_the_keys", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			s.unreadable[keyProfile] = true
			e, _ := readerExec(t, s, requested...)

			got := e.Read(ctxOf(t), readerName)

			r := mustReadRefuse(t, got, accessor.ClassIncompleteRead)
			if len(got.Values) != 0 {
				t.Errorf("the refusal carries values %+v; it must carry NONE — an "+
					"implementation deriving the requested set from what it read "+
					"reports a one-key success here", got.Values)
			}
			if !slices.Contains(r.Keys, keyProfile) {
				t.Errorf("refusal names keys %v; want %q named", r.Keys, keyProfile)
			}
		})

		t.Run("genuine_absence_is_a_value_not_a_refusal", func(t *testing.T) {
			// Distinguishable from the truncation above by BRANCH.
			s := newStore(map[string]string{keyStatus: "Draft"})
			e, _ := readerExec(t, s, requested...)

			values := mustReadSucceed(t, e.Read(ctxOf(t), readerName))

			v, ok := valueOf(values, keyProfile)
			if !ok || !v.Absent {
				t.Errorf("value for %q = %+v (found=%v); want an ABSENT value on "+
					"the success branch", keyProfile, v, ok)
			}
		})

		t.Run("execution_failure", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft"})
			e := accessor.NewExecutor(
				registryOf(readerDef(&readBinding{store: s, failWith: errBindingFailed},
					keyStatus)), artifactsOf())

			mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassExecutionFailure)
		})

		t.Run("timeout", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft"})
			e := accessor.NewExecutor(
				registryOf(readerDef(&readBinding{store: s, delay: time.Hour}, keyStatus)),
				artifactsOf())

			mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassTimeout)
		})

		t.Run("capability_mismatch", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft"})
			e := accessor.NewExecutor(
				registryOf(writerDef(&writeBinding{store: s}, keyStatus)), artifactsOf())

			mustReadRefuse(t, e.Read(ctxOf(t), writerName), accessor.ClassCapabilityMismatch)
		})

		t.Run("gate_allow_and_deny_are_typed_results", func(t *testing.T) {
			allow := gateExec(t, &gateBinding{verdict: accessor.VerdictAllow}).
				Gate(ctxOf(t), gateName)
			if allow.Refused() || allow.Verdict != accessor.VerdictAllow {
				t.Errorf("gate allow = %+v; want the typed verdict %q",
					allow, accessor.VerdictAllow)
			}

			const reason = "frozen for release"
			deny := gateExec(t, &gateBinding{verdict: accessor.VerdictDeny, reason: reason}).
				Gate(ctxOf(t), gateName)
			if deny.Refused() {
				t.Errorf("gate deny refused with %q; at this boundary deny is a "+
					"TYPED result carrying a reason (JDR 0001 §D9)", deny.Refusal.Class)
			}
			if deny.Verdict != accessor.VerdictDeny || deny.Reason != reason {
				t.Errorf("gate deny = {verdict %q, reason %q}; want {%q, %q}",
					deny.Verdict, deny.Reason, accessor.VerdictDeny, reason)
			}
		})

		t.Run("gate_indeterminate_is_a_refusal", func(t *testing.T) {
			got := gateExec(t, &gateBinding{verdict: accessor.VerdictIndeterminate}).
				Gate(ctxOf(t), gateName)

			if !got.Refused() {
				t.Fatalf("gate indeterminate produced verdict %q; it is a "+
					"refusal-class result", got.Verdict)
			}
			if got.Refusal.Class != accessor.ClassGateIndeterminate {
				t.Errorf("refusal class = %q; want %q",
					got.Refusal.Class, accessor.ClassGateIndeterminate)
			}
		})
	})

	// Scenario 3 is the Round-Trip / Inverse Invariant. The reconstructed
	// value is compared VALUE-FOR-VALUE against the plan's expected
	// values; a green exit or "did not error" is explicitly not enough.
	t.Run("3_write_read_back_round_trip_invariant", func(t *testing.T) {
		t.Run("owned_tag_value_identity_and_non_owned_identity", func(t *testing.T) {
			s := newStore(map[string]string{
				keyStatus:  "Draft",
				keyProfile: "large", // observed: protected, must not move
				keyLabels:  `["alpha"]`,
			})
			// The pre-write snapshot of the protected non-owned values.
			protectedBefore := map[string]string{keyProfile: "large"}

			w := &writeBinding{store: s}
			r := &readBinding{store: s}
			e := accessor.NewExecutor(
				registryOf(readerDef(r, keyStatus, keyLabels, keyProfile),
					writerDef(w, keyStatus, keyLabels)),
				artifactsOf())

			plan := planWriting(
				resolve.Tag{Key: keyStatus, Value: "Final"},
				resolve.Tag{Key: keyLabels, Value: `["alpha","beta"]`},
			)

			got := e.Write(ctxOf(t), writerName, plan)
			mustWriteSucceed(t, got)

			// --- the inverse leg: reconstruct through the READ path -----
			reconstructed := mustReadSucceed(t, e.Read(ctxOf(t), readerName))

			// (a) owned-tag VALUE IDENTITY, value-for-value against the plan.
			for _, planned := range plan.Writes {
				v, ok := valueOf(reconstructed, planned.Key)
				if !ok {
					t.Errorf("the re-read does not carry the planned key %q at all",
						planned.Key)
					continue
				}
				if v.Absent {
					t.Errorf("the re-read reports %q ABSENT; the plan assigned %q",
						planned.Key, planned.Value)
					continue
				}
				if v.Value != planned.Value {
					t.Errorf("write -> read identity broken for %q: re-read holds "+
						"%q, the transition plan expected %q. The invariant is "+
						"VALUE EQUALITY, not \"the write did not error\"",
						planned.Key, v.Value, planned.Value)
				}
			}

			// (b) protected NON-OWNED tag identity.
			for key, before := range protectedBefore {
				v, ok := valueOf(reconstructed, key)
				if !ok || v.Absent {
					t.Errorf("protected non-owned %q vanished across the write; "+
						"values present before the write must be unchanged", key)
					continue
				}
				if v.Value != before {
					t.Errorf("protected non-owned %q moved from %q to %q across "+
						"the write", key, before, v.Value)
				}
			}

			// (c) the written result echoes the same values, so a caller
			//     reading the disposition sees the same identity.
			for _, planned := range plan.Writes {
				if v, ok := seamValueOf(got.Written, planned.Key); !ok || v != planned.Value {
					t.Errorf("WriteResult.Written for %q = %q (found=%v); want the "+
						"verified %q", planned.Key, v, ok, planned.Value)
				}
			}
		})

		t.Run("owned_mismatch_control", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			w := &writeBinding{store: s, corrupt: func(st *store) {
				st.tags[keyStatus] = "corrupt"
			}}
			e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

			mustWriteRefuse(t,
				e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})),
				accessor.ClassReadBackMismatch)
		})

		t.Run("non_owned_mismatch_control", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			w := &writeBinding{store: s, corrupt: func(st *store) {
				st.tags[keyProfile] = "small"
			}}
			e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

			mustWriteRefuse(t,
				e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})),
				accessor.ClassReadBackMismatch)
		})

		t.Run("read_back_incomplete_neither_success_nor_mismatch", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft"})
			w := &writeBinding{store: s}
			e := accessor.NewExecutor(
				registryOf(
					readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
					writerDef(w, keyStatus)),
				artifactsOf())

			got := e.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

			mustWriteRefuse(t, got, accessor.ClassReadBackIncomplete)
		})
	})

	t.Run("4_replay_stability", func(t *testing.T) {
		run := func() accessor.Disposition {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			def := readerDef(&readBinding{store: s}, keyStatus, keyProfile)
			e := accessor.NewExecutor(registryOf(def), artifactsOf())
			return accessor.ReadDisposition(def, e.Read(ctxOf(t), readerName))
		}

		first, second := run(), run()
		if !reflect.DeepEqual(first, second) {
			t.Errorf("replay-identical = false:\nfirst  = %+v\nsecond = %+v",
				first, second)
		}
		if first.Refusal != "" {
			t.Errorf("the replayed success run refused with %q", first.Refusal)
		}

		// Control: the injected refusal is stable AND is not success.
		injected := func() accessor.Disposition {
			def := gateDef(&gateBinding{verdict: accessor.VerdictIndeterminate})
			e := accessor.NewExecutor(registryOf(def), artifactsOf())
			return accessor.GateDisposition(def, e.Gate(ctxOf(t), gateName))
		}
		a, b := injected(), injected()
		if a.Refusal != accessor.ClassGateIndeterminate {
			t.Errorf("replay-injected refusal = %q; want %q — a replay that "+
				"returns success unconditionally fails here",
				a.Refusal, accessor.ClassGateIndeterminate)
		}
		if !reflect.DeepEqual(a, b) {
			t.Errorf("the injected refusal is not stable:\n%+v\n%+v", a, b)
		}
	})

	t.Run("5_no_package_prints", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true

		var got accessor.ReadResult
		stdout, stderr := captureOutput(t, func() {
			e, _ := readerExec(t, s, keyStatus, keyProfile)
			got = e.Read(ctxOf(t), readerName)
		})

		// Positive half — the scenario cannot pass by absence of output.
		if got.Refusal == nil || got.Refusal.Class != accessor.ClassIncompleteRead {
			t.Errorf("returned value = %+v; want a structured refusal carrying "+
				"class %q", got, accessor.ClassIncompleteRead)
		}
		// Negative half — the normative capture control.
		if stdout != "" || stderr != "" {
			t.Errorf("the accessor package printed: stdout=%q stderr=%q",
				stdout, stderr)
		}
	})

	t.Run("6_absence_reaches_the_resolver", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		read := e.Read(ctxOf(t), readerName)
		mustReadSucceed(t, read)

		result, err := resolve.Resolve(resolve.Input{
			Flow:       flowID,
			Table:      requiringTable(keyProfile),
			Owned:      read.OwnedSnapshot(),
			Recognized: recognizedOutcome,
		})
		if err != nil {
			t.Fatalf("the kernel used the Go error path: %v", err)
		}
		if result.Refusal == nil {
			t.Fatalf("the resolver PLANNED over a snapshot that never carried %q; "+
				"a placeholder makes `TagSet.has` true and the refusal never fires",
				keyProfile)
		}
		if result.Refusal.Kind != resolve.KindOwnedStateUnavailable ||
			!slices.Contains(result.Refusal.MissingOwned, keyProfile) {
			t.Errorf("resolver refusal = {kind %q, missing %v}; want %q naming %q",
				result.Refusal.Kind, result.Refusal.MissingOwned,
				resolve.KindOwnedStateUnavailable, keyProfile)
		}

		// Control: a carried key resolves normally, so the refusal above
		// cannot pass vacuously over an implementation that omits everything.
		carried := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		ec, _ := readerExec(t, carried, keyStatus, keyProfile)
		cr := ec.Read(ctxOf(t), readerName)
		mustReadSucceed(t, cr)
		ctrl, err := resolve.Resolve(resolve.Input{
			Flow: flowID, Table: requiringTable(keyProfile),
			Owned: cr.OwnedSnapshot(), Recognized: recognizedOutcome,
		})
		if err != nil {
			t.Fatalf("the kernel used the Go error path: %v", err)
		}
		if ctrl.Plan == nil {
			t.Errorf("the control refused with %+v over a snapshot that DOES "+
				"carry %q; an implementation omitting every key fails here",
				ctrl.Refusal, keyProfile)
		}
	})

	t.Run("7_timeout_outranks_incomplete_read", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		b := &overlapReadBinding{store: s, delay: time.Hour}
		e := accessor.NewExecutor(
			registryOf(readerDef(b, keyStatus, keyProfile)), artifactsOf())

		got := e.Read(ctxOf(t), readerName)

		if !b.sawUnreadable {
			t.Fatal("the fixture did not make both classes true at once")
		}
		mustReadRefuse(t, got, accessor.ClassTimeout)

		// Control: truncation alone must still be `incomplete_read`, so an
		// implementation that always reports `timeout` fails.
		plain := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		plain.unreadable[keyProfile] = true
		ec, _ := readerExec(t, plain, keyStatus, keyProfile)
		mustReadRefuse(t, ec.Read(ctxOf(t), readerName), accessor.ClassIncompleteRead)
	})

	t.Run("8_post_mutation_reporting_and_no_compensation", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
				writerDef(w, keyStatus)),
			artifactsOf())

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		r := mustWriteRefuse(t, got, accessor.ClassReadBackIncomplete)
		if !r.Applied() {
			t.Error("Applied() = false; the refusal must carry the " +
				"applied-but-unverified sense, never \"the write did not occur\"")
		}
		if w.Invocations() != 1 {
			t.Errorf("the write binding ran %d times; want EXACTLY 1 — no retry, "+
				"no undo, no re-derivation", w.Invocations())
		}
		if s.tags[keyStatus] != "Final" {
			t.Errorf("artifact holds %q = %q; the mutation landed and was NOT "+
				"compensated", keyStatus, s.tags[keyStatus])
		}

		// Control: read_back_mismatch asserts the artifact is WRONG, not
		// unverified — collapsing both into one shape fails here.
		ms := newStore(map[string]string{keyStatus: "Draft"})
		mw := &writeBinding{store: ms, corrupt: func(st *store) {
			st.tags[keyStatus] = "corrupt"
		}}
		me := writeExec(t, ms, mw, &readBinding{store: ms}, keyStatus)
		mr := mustWriteRefuse(t,
			me.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})),
			accessor.ClassReadBackMismatch)
		if mr.Applied() {
			t.Error("the mismatch refusal reports Applied()=true; its " +
				"verification RAN and found the artifact wrong — that is a " +
				"different sense from applied-but-unverified")
		}
	})

	t.Run("9_clearing_write", func(t *testing.T) {
		t.Run("clear_a_held_key_re_read_shows_it_absent", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
			w := &writeBinding{store: s}
			r := &readBinding{store: s}
			e := accessor.NewExecutor(
				registryOf(readerDef(r, keyStatus, keyProfile), writerDef(w, keyStatus)),
				artifactsOf())

			mustWriteSucceed(t, e.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})))

			// The inverse leg for a CLEAR: the invariant's expected-value
			// predicate is ABSENCE, asserted through the read path.
			reconstructed := mustReadSucceed(t, e.Read(ctxOf(t), readerName))
			v, ok := valueOf(reconstructed, keyStatus)
			if !ok || !v.Absent {
				t.Errorf("the re-read holds %q = %+v; a planned `%s` removes the "+
					"key, so the round trip's expected value is ABSENCE",
					keyStatus, v, accessor.ClearSentinel)
			}
			if p, ok := valueOf(reconstructed, keyProfile); !ok || p.Value != "large" {
				t.Errorf("protected non-owned %q = %+v; want %q unchanged",
					keyProfile, p, "large")
			}
		})

		t.Run("clear_a_key_the_artifact_does_not_hold_succeeds", func(t *testing.T) {
			s := newStore(map[string]string{keyProfile: "large"})
			w := &writeBinding{store: s}
			e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

			mustWriteSucceed(t, e.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})))

			if w.Invocations() != 1 {
				t.Errorf("the write binding ran %d times; the idempotent clear "+
					"still runs the write", w.Invocations())
			}
		})

		t.Run("clear_whose_re_read_still_holds_the_key_is_mismatch", func(t *testing.T) {
			// The binding stores the LITERAL instead of removing.
			s := newStore(map[string]string{keyStatus: "Draft"})
			w := &writeBinding{store: s, storeLiteral: true}
			e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

			mustWriteRefuse(t, e.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})),
				accessor.ClassReadBackMismatch)
		})

		t.Run("read_holding_the_literal_clear_refuses_incomplete_read", func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: accessor.ClearSentinel})
			e, _ := readerExec(t, s, keyStatus)

			r := mustReadRefuse(t, e.Read(ctxOf(t), readerName),
				accessor.ClassIncompleteRead)
			if !slices.Contains(r.Keys, keyStatus) {
				t.Errorf("refusal names keys %v; want %q named", r.Keys, keyStatus)
			}
		})

		t.Run("control_the_assignment_fixture_still_asserts_presence_and_equality",
			func(t *testing.T) {
				s := newStore(map[string]string{keyStatus: "Draft"})
				w := &writeBinding{store: s}
				r := &readBinding{store: s}
				e := writeExec(t, s, w, r, keyStatus)

				mustWriteSucceed(t, e.Write(ctxOf(t), writerName,
					planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))

				reconstructed := mustReadSucceed(t, e.Read(ctxOf(t), readerName))
				v, ok := valueOf(reconstructed, keyStatus)
				if !ok || v.Absent || v.Value != "Final" {
					t.Errorf("the re-read holds %q = %+v; an ASSIGNMENT asserts "+
						"presence and equality — a read-back that asserts absence "+
						"for every write fails this control", keyStatus, v)
				}
			})
	})
}
