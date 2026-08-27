package accessor_test

// RDR 0004 — the desk trace's simultaneity claim, the phase
// deliverables, and the pending-assumption discharge ledger
// (REQ-92..REQ-94, REQ-101, REQ-105, REQ-106, REQ-111, REQ-114..REQ-116,
// REQ-119, REQ-121..REQ-123).
//
// The desk trace's claim is not that each assertion holds in isolation —
// the per-clause files already hold those — but that the whole set holds
// AT ONCE on one invocation. A design that satisfies them one at a time
// and conflicts across them fails here and nowhere else.

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-101 / DESK step 2: nine assertions hold at once on ONE read
// invocation: "read returns typed values or a typed refusal; completeness
// — exactly the requested keys or refuse; requested set is validated
// definition metadata, never derived from what resolved; one unreadable
// key refuses the whole read; absence crosses the seam as omission, not a
// placeholder value; timeout outranks incomplete read; no mutation of
// authoritative artifacts; bounded timeout; no direct stdout/stderr".
// REQ-114 / TS 2: "Successful reads return the complete typed tag values
// for every requested key … and timeout, execution failure, incomplete
// read, and gate indeterminate remain distinct refusal classes."
// REQ-115 / TS 2: "an artifact that genuinely lacks a requested key
// returns it as an absent value, not a refusal."
// REQ-105 / RM: "Completeness is normative — a read that cannot resolve
// every requested key refuses."
// REQ-106 / RM: "Absence crosses as omission from the owned snapshot."
// BOUNDARY
func TestReq101_DeskTraceStep2AssertionsHoldSimultaneously(t *testing.T) {
	requested := []string{keyStatus, keyProfile}

	// The artifact carries `status`, genuinely lacks `profile`, and holds
	// a third key nobody requested — so completeness, absence-as-value,
	// and no-unrequested-key are all in force on the same invocation.
	s := newStore(map[string]string{keyStatus: "Draft", keyLabels: `["alpha"]`})
	before := map[string]string{}
	for k, v := range s.tags {
		before[k] = v
	}

	var got accessor.ReadResult
	var b *readBinding
	stdout, stderr := captureOutput(t, func() {
		var e *accessor.Executor
		e, b = readerExec(t, s, requested...)
		got = e.Read(ctxOf(t), readerName)
	})

	// 1. typed values or a typed refusal — here, the values branch.
	values := mustReadSucceed(t, got)

	// 2. exactly the requested keys.
	want := slices.Clone(requested)
	slices.Sort(want)
	if k := keysOf(values); !slices.Equal(k, want) {
		t.Errorf("returned keys %v; want exactly the requested %v — the "+
			"unrequested %q must not be added", k, want, keyLabels)
	}

	// 3. the requested set is definition metadata, not what resolved.
	if k := b.lastRequested; !slices.Equal(sorted(k), want) {
		t.Errorf("the binding was asked for %v; want the definition's declared "+
			"%v — never a set derived from the artifact", k, want)
	}

	// 4/5. absence is a VALUE here, and it crosses the seam as OMISSION.
	v, ok := valueOf(values, keyProfile)
	if !ok || !v.Absent {
		t.Errorf("value for %q = %+v (found=%v); want an absent VALUE on the "+
			"success branch", keyProfile, v, ok)
	}
	snapshot := got.OwnedSnapshot()
	if sv, present := seamValueOf(snapshot, keyStatus); !present || sv != "Draft" {
		t.Fatalf("the seam snapshot does not carry %q = %q (%+v)",
			keyStatus, "Draft", snapshot)
	}
	if _, present := seamValueOf(snapshot, keyProfile); present {
		t.Errorf("the seam snapshot carries the absent %q; it must be OMITTED, "+
			"never a placeholder", keyProfile)
	}

	// 6. no mutation of authoritative artifacts.
	for k, wantValue := range before {
		if s.tags[k] != wantValue {
			t.Errorf("artifact %q moved from %q to %q across a READ",
				k, wantValue, s.tags[k])
		}
	}
	if len(s.tags) != len(before) {
		t.Errorf("artifact key count moved from %d to %d across a READ",
			len(before), len(s.tags))
	}

	// 7. no direct stdout/stderr.
	if stdout != "" || stderr != "" {
		t.Errorf("the accessor package printed: stdout=%q stderr=%q", stdout, stderr)
	}

	// 8/9. bounded timeout, and timeout outranks incomplete read — the
	// same definition, with the two overlapping in one invocation.
	t.Run("timeout_outranks_incomplete_read_on_the_same_definition", func(t *testing.T) {
		overlap := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		overlap.unreadable[keyProfile] = true
		ob := &overlapReadBinding{store: overlap, delay: time.Hour}
		e := accessor.NewExecutor(registryOf(readerDef(ob, requested...)), artifactsOf())

		start := time.Now()
		r := e.Read(ctxOf(t), readerName)
		elapsed := time.Since(start)

		// Bounded by the DECLARED timeout, not by the binding's delay:
		// comparing against `time.Hour` would pass a deadline wrong by
		// five orders of magnitude and assert nothing about REQ-68.
		if elapsed >= 20*fixtureTimeout {
			t.Fatalf("the read took %v; the invocation is bounded by the "+
				"definition's declared timeout %v", elapsed, fixtureTimeout)
		}
		mustReadRefuse(t, r, accessor.ClassTimeout)

		// The two classes were true AT ONCE: the binding MET the
		// unreadable key and only then ran past its deadline. Without
		// this, an executor whose read timed out before ever reaching
		// the unreadable key passes — and REQ-101's whole claim is that
		// the assertions hold simultaneously on ONE invocation, so a
		// non-overlapping fixture witnesses nothing here (REQ-27).
		if !ob.sawUnreadable {
			t.Fatal("the fixture did not make both classes true at once; the " +
				"unreadable key must be met BEFORE the deadline expires")
		}
	})
}

// REQ-119 / TS 7: "The refusal is `timeout`, not `incomplete_read` — the
// two input classes overlap and timeout takes precedence."
// REQ-116 / TS 3: "Owned-tag values equal to the plan report success;
// mismatched owned-tag values or mutated non-owned observed/recognized
// tag values report read-back mismatch even when command-level write
// invocation succeeded. A read-back whose re-read cannot read a key it
// must compare reports `read_back_incomplete` — neither success nor
// mismatch (A9)."
// REQ-121 / TS 9: "The first two report success and the re-read does not
// hold the key; the third reports `read_back_mismatch`; the read refuses
// `incomplete_read` naming the key."
// BOUNDARY
//
// The three-way write discrimination, asserted as three DIFFERENT
// outcomes over one plan shape. This is the table an implementation that
// collapses any two of them cannot satisfy.
func TestReq116_TheWriteOutcomeTableDiscriminatesThreeWays(t *testing.T) {
	plan := func() resolve.Plan {
		return planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})
	}

	t.Run("equal_reports_success", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		r := &readBinding{store: s}
		e := writeExec(t, s, w, r, keyStatus)

		got := e.Write(ctxOf(t), writerName, plan())

		mustWriteSucceed(t, got)
		// The success must be EARNED: the write ran, the re-read ran, and
		// the artifact holds the planned value. An executor that reports
		// success without doing either passes vacuously.
		if w.Invocations() != 1 || r.reads == 0 {
			t.Fatalf("write invocations = %d, read-back reads = %d; want 1 and "+
				">0 — success is only meaningful once both actually ran",
				w.Invocations(), r.reads)
		}
		if s.tags[keyStatus] != "Final" {
			t.Errorf("artifact holds %q = %q; want %q",
				keyStatus, s.tags[keyStatus], "Final")
		}
	})

	t.Run("mismatched_reports_read_back_mismatch", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			st.tags[keyStatus] = "corrupt"
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		mustWriteRefuse(t, e.Write(ctxOf(t), writerName, plan()),
			accessor.ClassReadBackMismatch)
	})

	t.Run("unreadable_compared_key_reports_read_back_incomplete", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
				writerDef(&writeBinding{store: s}, keyStatus)),
			artifactsOf())

		mustWriteRefuse(t, e.Write(ctxOf(t), writerName, plan()),
			accessor.ClassReadBackIncomplete)
	})
}

// REQ-92 / IP Phase 1: "Define the accessor definition structs,
// capability enum, refusal classes, and validation rules that connect RDR
// 0002 accessor references to declared read/gate/write bindings."
// REQ-93 / IP Phase 2: "Implement context-bound invocation for typed
// accessor bindings. The executor returns structured success/refusal
// values and performs no direct output."
// REQ-94 / IP Phase 3: "Apply planned owned-tag writes through write
// accessors, then re-read the same artifact role and compare each planned
// owned tag for equality — absence for a `<clear>` — plus the pre-write
// observed and recognized tag values. Treat mismatch as a write failure."
// REQ-95 / EIA: "New internal package should own capability validation
// and invocation."
// REQ-111 / TS: "The MVV should become production tests around the
// accessor validator and executor boundary."
// BOUNDARY
//
// The phase deliverables, asserted as a reachable surface: the four
// vocabularies enumerate, the executor is context-bound, and the write
// path performs the read-back comparison.
func TestReq92_ThePhaseDeliverablesAreOneReachableSurface(t *testing.T) {
	t.Run("phase_1_vocabularies_enumerate", func(t *testing.T) {
		if len(accessor.Capabilities()) != 3 {
			t.Errorf("Capabilities() = %v; the capability enum is Phase 1's",
				accessor.Capabilities())
		}
		if len(accessor.RefusalClasses()) == 0 {
			t.Error("RefusalClasses() is empty; the refusal classes are Phase 1's")
		}
		if len(accessor.ValidationCodes()) == 0 {
			t.Error("ValidationCodes() is empty; the validation rules are Phase 1's")
		}
		if len(accessor.Verdicts()) != 3 {
			t.Errorf("Verdicts() = %v; want the three-valued gate vocabulary",
				accessor.Verdicts())
		}
	})

	t.Run("phase_2_invocation_is_context_bound", func(t *testing.T) {
		// The caller's context reaches the binding, so a caller-side
		// cancellation is observable at the seam alongside the
		// per-accessor timeout.
		s := newStore(map[string]string{keyStatus: "Draft"})
		b := &ctxWitnessReadBinding{store: s}
		e := accessor.NewExecutor(registryOf(readerDef(b, keyStatus)), artifactsOf())

		mustReadSucceed(t, e.Read(ctxOf(t), readerName))

		if !b.sawContext {
			t.Error("the binding received no context; invocation is context-bound")
		}
		if !b.sawDeadline {
			t.Error("the context carried no deadline; every accessor invocation " +
				"is bounded by its declared timeout")
		}
		// The CALLER'S context, not merely SOME context. `sawContext`
		// and `sawDeadline` are both satisfied by an executor that
		// DISCARDS what the caller passed and builds a fresh
		// `context.WithTimeout(context.Background(), timeout)` — under
		// which a caller-side cancellation is not observable at the seam
		// at all, which is the property this subtest names.
		if !b.sawCallerValue {
			t.Error("the binding's context does not descend from the caller's; " +
				"the invocation must be bound to the CALLER'S context, so a " +
				"caller-side cancellation is observable at the seam alongside " +
				"the per-accessor timeout")
		}
	})

	t.Run("phase_3_the_write_path_performs_the_read_back", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		r := &readBinding{store: s}
		e := writeExec(t, s, w, r, keyStatus)

		mustWriteSucceed(t, e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})))

		if r.reads == 0 {
			t.Error("the write path performed no re-read; Phase 3 is the " +
				"read-back verification, not the write alone")
		}
	})
}

// REQ-122: "**A9, A10, and A11 Pending** — MVV-proven properties carried
// to lock with the MVV as the named implementation-time plan"; "Verify
// and retire them per-rule rather than flipping A9 as a unit".
// REQ-123: A9's four rules, each verified by a named scenario: the
// missing/empty requested-key-set validation arm (Scenario 1's eighth
// arm), absence-to-resolver (Scenario 6), timeout-outranks-incomplete
// (Scenario 7), and read-back-incomplete (Scenario 3's unreadable-
// compared-key case).
// BOUNDARY
//
// The discharge ledger, asserted PER RULE so a partial refutation is
// visible: each of the six pending properties is exercised here
// independently, matching the record's instruction to retire them
// per-rule rather than as a unit.
func TestReq123_ThePendingAssumptionsAreDischargedPerRule(t *testing.T) {
	t.Run("A9_rule_1_missing_or_empty_requested_key_set", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		def := readerDef(&readBinding{store: s})
		def.Accessor.Keys = nil

		fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

		if !hasCode(fs, accessor.CodeMissingRequestedKeySet) {
			t.Errorf("codes = %v; want %q", findingCodes(fs),
				accessor.CodeMissingRequestedKeySet)
		}
	})

	t.Run("A9_rule_2_absence_crosses_as_omission_to_the_resolver", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e, _ := readerExec(t, s, keyStatus, keyProfile)

		read := e.Read(ctxOf(t), readerName)
		mustReadSucceed(t, read)
		snapshot := read.OwnedSnapshot()
		// The refusal must fire because `profile` was OMITTED, not because
		// the snapshot is empty: an empty snapshot omits every key and
		// would satisfy the refusal for the wrong reason.
		if v, ok := seamValueOf(snapshot, keyStatus); !ok || v != "Draft" {
			t.Fatalf("the snapshot does not carry %q = %q (%+v); the refusal "+
				"below would then fire over an empty snapshot rather than over "+
				"the omitted key", keyStatus, "Draft", snapshot)
		}
		result, err := resolve.Resolve(resolve.Input{
			Flow: flowID, Table: requiringTable(keyProfile),
			Owned: snapshot, Recognized: recognizedOutcome,
		})
		if err != nil {
			t.Fatalf("kernel error path: %v", err)
		}
		if result.Refusal == nil ||
			result.Refusal.Kind != resolve.KindOwnedStateUnavailable {
			t.Errorf("resolver disposition = %+v; want %q",
				result, resolve.KindOwnedStateUnavailable)
		}
	})

	t.Run("A9_rule_3_timeout_outranks_incomplete_read", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		s.unreadable[keyProfile] = true
		b := &overlapReadBinding{store: s, delay: time.Hour}
		e := accessor.NewExecutor(
			registryOf(readerDef(b, keyStatus, keyProfile)), artifactsOf())

		mustReadRefuse(t, e.Read(ctxOf(t), readerName), accessor.ClassTimeout)
	})

	t.Run("A9_rule_4_read_back_incomplete", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
				writerDef(&writeBinding{store: s}, keyStatus)),
			artifactsOf())

		mustWriteRefuse(t, e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})),
			accessor.ClassReadBackIncomplete)
	})

	t.Run("A10_post_mutation_reporting_without_compensation", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
				writerDef(w, keyStatus)),
			artifactsOf())

		r := mustWriteRefuse(t, e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})),
			accessor.ClassReadBackIncomplete)

		if !r.Applied() {
			t.Error("Applied() = false; want the applied-but-unverified sense")
		}
		if w.Invocations() != 1 {
			t.Errorf("write invocations = %d; want exactly 1", w.Invocations())
		}
	})

	t.Run("A11_clear_is_a_removal_verified_by_absence", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		e := writeExec(t, s, &writeBinding{store: s}, &readBinding{store: s}, keyStatus)

		mustWriteSucceed(t, e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: accessor.ClearSentinel})))

		if _, held := s.tags[keyStatus]; held {
			t.Errorf("the artifact still holds %q; a clear is a removal", keyStatus)
		}
	})
}

// --- helpers -------------------------------------------------------------

func sorted(in []string) []string {
	out := slices.Clone(in)
	slices.Sort(out)
	return out
}

// ctxWitnessReadBinding records whether the executor handed it a context
// carrying a deadline, and whether that context DESCENDS from the
// caller's rather than being a fresh one built over `Background()`.
type ctxWitnessReadBinding struct {
	store          *store
	sawContext     bool
	sawDeadline    bool
	sawCallerValue bool
}

func (b *ctxWitnessReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *ctxWitnessReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	b.sawContext = ctx != nil
	if ctx != nil {
		_, b.sawDeadline = ctx.Deadline()
		b.sawCallerValue = ctx.Value(callerCtxKey{}) == callerCtxValue
	}
	inner := &readBinding{store: b.store}
	return inner.Read(ctx, art, requested)
}
