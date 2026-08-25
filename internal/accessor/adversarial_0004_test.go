package accessor_test

// RDR 0004 — Phase 3b ADVERSARIAL failure-mode probes.
//
// These tests are written independently of the Phase 1 suite. They target
// the three shapes in which a competent-but-wrong executor most plausibly
// returns a SUCCESS-SHAPED result over authoritative artifact state it did
// not establish — the RDR's own definition of a silent failure (`0004:FM`,
// REQ-80). This RDR governs authoritative artifact mutation with no ledger
// and no undo, so a green result over a corrupted artifact is the highest
// cost defect the boundary can produce.

import (
	"testing"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// --- ADV-1 ---------------------------------------------------------------

// ADV-1 — a transition plan carrying a key the write accessor does NOT own
// is applied to the authoritative artifact and reported as SUCCESS.
//
// Failure mode (`0004:FM`, "A silent failure is any outcome where the
// accessor returns a success-shaped result over state it did not
// establish"): `Validate` checks only the writer DEFINITION's declared
// `keys` against the model's owned tags (`validate.go::Validate`,
// CodeWriteNonOwnedTag). Nothing at EXECUTION time constrains the keys in
// `plan.Writes` to the writer's owned key set, so a plan naming an
// observed or recognized tag reaches `WriteBinding.Apply` unfiltered.
//
// The corruption is doubly invisible: because `protectedKeys` excludes
// every planned key from the pre-write protected snapshot
// (`executor.go::protectedKeys`), the very tag that was clobbered is
// removed from the non-owned identity comparison, and `verifyReadBack`
// then confirms the clobbered value equals the plan. Read-back reports
// green over an artifact whose observed tag the accessor had no authority
// to touch.
//
// Defends: `0004:C10` / REQ-40, REQ-41 ("A write accessor MUST apply only
// planned owned-tag writes produced by a successful transition. It MUST
// NOT write observed or recognized tags."), REQ-58 (protected non-owned
// tag identity), REQ-80.
func TestAdv1_WriteRefusesAPlanNamingANonOwnedTag(t *testing.T) {
	// `profile` is an OBSERVED tag: registryOf declares OwnedTags as
	// {status, labels} only. The writer owns `status`.
	s := newStore(map[string]string{keyStatus: "draft", keyProfile: "large"})
	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	// The plan names the non-owned observed tag alongside the owned one.
	plan := planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
		resolve.Tag{Key: keyProfile, Value: "small"},
	)

	got := e.Write(ctxOf(t), writerName, plan)

	if !got.Refused() {
		t.Fatalf("Write applied a plan naming the non-owned tag %q and reported SUCCESS "+
			"(Written=%+v); want a refusal: `0004:C10` forbids writing observed or "+
			"recognized tags, and no runtime arm enforces it", keyProfile, got.Written)
	}

	// The refusal must not be a generic execution failure that happens to
	// fire for another reason; it must attribute the non-owned key.
	if class := got.Refusal.Class; class == accessor.ClassReadBackMismatch {
		t.Errorf("Write refused %q; want a refusal that names the non-owned key rather "+
			"than one asserting the artifact is wrong", class)
	}

	// The authoritative artifact must be untouched for the non-owned tag.
	if v, absent, _ := s.get(keyProfile); absent || v != "large" {
		t.Errorf("observed tag %q = (%q, absent=%v) after the refused write; want %q "+
			"unchanged — the write must not reach the artifact at all",
			keyProfile, v, absent, "large")
	}
}

// --- ADV-2 ---------------------------------------------------------------

// ADV-2 — the pre-write protected snapshot fails, the non-owned identity
// guard therefore never runs, and the write reports SUCCESS over a
// clobbered observed tag.
//
// Failure mode: `executor.go::Write` takes the pre-write snapshot through
// `invokeRead` and then does `if raw.class == ""` — on a timeout, an
// execution failure, or an incomplete pre-write read it SILENTLY leaves
// `before` empty and proceeds. `verifyReadBack` iterates `before`, so an
// empty map means the "observed and recognized tag values present before
// the write are unchanged" clause is vacuous. A write binding that also
// mutates a non-owned tag then passes read-back green.
//
// This is the RDR's stated risk verbatim: "A write command succeeds but
// changes the wrong artifact or wrong tag" (RM), whose mitigation is
// same-role read-back verification. A verification that cannot establish
// the pre-write value must not report success; the contract's own vocabulary
// for "the verification did not run" is `read_back_incomplete` (`0004:C13`).
//
// Defends: `0004:C12` / REQ-53, REQ-58 (RT invariant, "any observed or
// recognized tag values present before the write must remain unchanged"),
// REQ-62 (`read_back_incomplete` is neither mismatch nor success), REQ-80.
func TestAdv2_UnreadablePreWriteSnapshotMustNotYieldSuccess(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "draft", keyProfile: "large"})

	// The write clobbers the protected observed tag alongside its own
	// owned write — precisely the shape read-back exists to catch.
	w := &writeBinding{store: s, corrupt: func(st *store) {
		st.tags[keyProfile] = "small"
	}}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	// The PRE-write read of the protected key cannot resolve it. The
	// post-write re-read can, so the executor sees a readable `profile`
	// with no baseline to compare it against.
	s.unreadable[keyProfile] = true
	restore := func() { delete(s.unreadable, keyProfile) }
	w.corrupt = func(st *store) {
		st.tags[keyProfile] = "small"
		restore()
	}

	got := e.Write(ctxOf(t), writerName, planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
	))

	if !got.Refused() {
		t.Fatalf("Write reported SUCCESS (Written=%+v) after an unreadable pre-write "+
			"snapshot left the protected non-owned comparison vacuous, while the write "+
			"changed %q from %q to %q; want a refusal — the verification did not run",
			got.Written, keyProfile, "large", "small")
	}

	switch got.Refusal.Class {
	case accessor.ClassReadBackIncomplete, accessor.ClassReadBackMismatch:
		// Either is defensible: the baseline was unreadable
		// (`read_back_incomplete`), or the boundary otherwise establishes
		// the observed tag moved (`read_back_mismatch`).
	default:
		t.Errorf("Write refused %q; want read_back_incomplete or read_back_mismatch",
			got.Refusal.Class)
	}

	if !got.Refusal.Applied() {
		t.Errorf("Refusal.Applied() = false; the write command already ran, so the " +
			"refusal MUST carry the applied-but-unverified sense (`0004:C14`)")
	}
	if w.Invocations() != 1 {
		t.Errorf("write binding invoked %d times; want exactly 1 — the accessor layer "+
			"performs no retry, undo, or re-derivation (`0004:C14`)", w.Invocations())
	}
}

// --- ADV-3 ---------------------------------------------------------------

// ADV-3 — a successful CLEARING write records the reserved `<clear>`
// literal as a held tag value in its own replay-stable disposition.
//
// Failure mode: `executor.go::Write` returns `WriteResult{Written: planned}`
// on success, and `planned` still carries `Tag{Key: k, Value: "<clear>"}`
// for a cleared key. `disposition.go::WriteDisposition` renders that set
// verbatim through `observedTags`, which only drops KeyValue.Absent — and a
// tag rebuilt from `planned` by `tagsAsValues` is never marked absent. So
// the disposition for a verified REMOVAL asserts the artifact HOLDS the key
// with the reserved literal as its value.
//
// That is exactly the encoding the RDR names as the one that turns its
// safety rule into a regression: "A sentinel value would make an absent key
// read as *present* and silently retire `owned_state_unavailable` for it"
// (LBD, Absence crosses the seam as omission; REQ-32). The clear-semantics
// LBD makes the same point for writes: the literal-storing clear is "the one
// write whose read-back passes green over a tag that was never removed".
// Here the artifact is correct and the RECORD is the placeholder — a
// success-shaped record over state the artifact does not hold, and the
// disposition is what a replay consumes (CA A4).
//
// Defends: `0004:C11` / REQ-45, REQ-46 (a clear is a removal, read-back
// asserts absence), REQ-32 (no sentinel crosses as a present value),
// REQ-80.
func TestAdv3_ClearingWriteDispositionMustNotRecordTheClearLiteral(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "draft", keyLabels: "a", keyProfile: "large"})
	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus, keyLabels)

	plan := planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
		resolve.Tag{Key: keyLabels, Value: table.ClearSentinel},
	)

	got := e.Write(ctxOf(t), writerName, plan)
	if got.Refused() {
		t.Fatalf("clearing write refused %q; the removal read back as absent, so this "+
			"is the SUCCESS control for this probe", got.Refusal.Class)
	}

	// Control: the artifact really did remove the key. If this fails the
	// probe below is testing the wrong thing.
	if _, absent, readable := s.get(keyLabels); !readable || !absent {
		t.Fatalf("artifact still holds %q after a successful clear; the removal itself "+
			"did not happen", keyLabels)
	}

	d := accessor.WriteDisposition(writerDef(w, keyStatus, keyLabels), got)

	for _, tag := range d.Tags {
		if tag.Value == table.ClearSentinel {
			t.Errorf("WriteDisposition recorded %+v: the replay-stable record asserts "+
				"the artifact HOLDS %q as the reserved literal %q, while read-back "+
				"verified the key ABSENT. A cleared key must be OMITTED from the "+
				"recorded tag set, never carried as a sentinel value (REQ-32, REQ-46)",
				tag, tag.Key, table.ClearSentinel)
		}
		if tag.Key == keyLabels {
			t.Errorf("WriteDisposition recorded the cleared key %q with value %q; a "+
				"cleared key must not appear in the recorded tag set at all",
				tag.Key, tag.Value)
		}
	}

	// The same defect on the result value itself: `Written` is documented
	// as "the planned owned tags the read-back verified", and what
	// read-back verified for a cleared key is ABSENCE.
	for _, tag := range got.Written {
		if tag.Key == keyLabels && tag.Value == table.ClearSentinel {
			t.Errorf("WriteResult.Written carries %+v; read-back verified %q absent, "+
				"so echoing the reserved literal as a verified held value is a "+
				"success-shaped record over state the artifact does not hold",
				tag, keyLabels)
		}
	}
}

// --- ADV-4 ---------------------------------------------------------------

// ADV-4 — a PARTIAL pre-write snapshot silently drops the unread keys'
// protection (Phase 3a FAIL-2).
//
// A narrower and more reachable variant of ADV-2: the pre-write read does
// NOT fail as a whole, so `raw.class` is empty and the wholly-failed-
// snapshot guard never fires. It simply names one protected key
// unreadable. `readOutcome.classify` reports that key in its `unread`
// return, which the pre-write path discarded — unlike the post-write
// read-back, which correctly refuses `read_back_incomplete` on a non-empty
// `unread`. The unread key therefore acquires no baseline, and
// `verifyReadBack`'s `before` loop has nothing to compare it against, so a
// write that clobbers it passes green.
//
// `0004:C12` requires observed and recognized tag values "present before
// the write" to be verified unchanged. A key whose pre-write value was
// never established cannot support that claim, so the verification did not
// run for it: `read_back_incomplete` (`0004:C13`, REQ-62), carrying the
// applied-but-unverified sense because the command already ran
// (`0004:C14`, REQ-63).
//
// Defends: `0004:C12` / REQ-53, REQ-58, REQ-62, REQ-80.
func TestAdv4_PartialPreWriteSnapshotMustNotYieldSuccess(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "draft", keyProfile: "large"})

	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	// The pre-write read RESOLVES `status` and names only `profile`
	// unreadable, so the read as a whole does not refuse — `raw.class`
	// stays empty and only `classify`'s `unread` return carries the gap.
	// The post-write re-read sees a readable, clobbered `profile`.
	s.unreadable[keyProfile] = true
	w.corrupt = func(st *store) {
		st.tags[keyProfile] = "small"
		delete(st.unreadable, keyProfile)
	}

	got := e.Write(ctxOf(t), writerName, planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
	))

	if !got.Refused() {
		t.Fatalf("Write reported SUCCESS (Written=%+v) after a PARTIAL pre-write "+
			"snapshot left %q with no baseline, while the write changed it from %q "+
			"to %q; want a refusal — the verification did not run for that key",
			got.Written, keyProfile, "large", "small")
	}

	switch got.Refusal.Class {
	case accessor.ClassReadBackIncomplete, accessor.ClassReadBackMismatch:
	default:
		t.Errorf("Write refused %q; want read_back_incomplete or read_back_mismatch",
			got.Refusal.Class)
	}

	if !got.Refusal.Applied() {
		t.Errorf("Refusal.Applied() = false; the write command already ran, so the " +
			"refusal MUST carry the applied-but-unverified sense (`0004:C14`)")
	}
	if w.Invocations() != 1 {
		t.Errorf("write binding invoked %d times; want exactly 1 — no retry, undo, "+
			"or re-derivation (`0004:C14`)", w.Invocations())
	}
}

// ADV-4 negative control: with a FULLY readable pre-write snapshot and no
// clobber, the identical fixture must SUCCEED. Without this, ADV-4 would
// pass against an executor that refuses every write.
func TestAdv4Control_CompletePreWriteSnapshotStillSucceeds(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "draft", keyProfile: "large"})
	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	got := e.Write(ctxOf(t), writerName, planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
	))
	mustWriteSucceed(t, got)

	if v, absent, _ := s.get(keyProfile); absent || v != "large" {
		t.Errorf("protected tag %q = (%q, absent=%v); want %q untouched",
			keyProfile, v, absent, "large")
	}
}

// --- ADV-5 ---------------------------------------------------------------

// ADV-5 — the write binding MUTATES the planned slice it was handed, and
// the read-back oracle therefore judges the artifact against the value the
// binding itself chose.
//
// Failure mode: `executor.go::Write` clones `plan.Writes` ONCE into
// `planned`. That single backing array is simultaneously the argument to
// `WriteBinding.Apply`, the expectation `verifyReadBack` compares the
// re-read artifact against, the `Expected` field on every typed refusal,
// and the source of the success record `verifiedWritten(planned)`. A
// binding that assigns to `planned[i].Value` before applying it rewrites
// the expectation it is about to be judged against: read-back compares the
// artifact to what the binding wrote, finds them equal, and `Written`
// reports the binding's substituted value as VERIFIED.
//
// That is a self-referential oracle — the verification is defeatable by the
// very component it verifies — and it produces the RDR's own definition of
// a silent failure: a success-shaped result over authoritative artifact
// state the accessor did not establish (`0004:FM`, REQ-80). Bindings are
// not trusted at this boundary: `0004:465-474` makes the binding's own
// reporting the thing the contract polices and `0004:801` names accessor
// bindings a risk surface. Deviation D16 already adjudicated the sibling
// case one hop earlier — the caller-supplied plan is enforced at RUNTIME
// rather than trusted from validation.
//
// The expectation must be immutable with respect to the binding: the
// artifact holds what the binding actually wrote, which differs from the
// plan, so this is `read_back_mismatch` (`0004:C12`, REQ-53).
//
// Defends: `0004:C12` / REQ-53, REQ-80.
func TestAdv5_BindingMutatingItsPlanMustNotRewriteTheReadBackExpectation(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "draft", keyProfile: "large"})

	// The binding substitutes its own value for the planned one, in place,
	// then applies the rewritten plan. Both the artifact and — under
	// aliasing — the expectation end up holding "hijacked".
	w := &writeBinding{store: s, rewritePlan: func(planned []resolve.Tag) {
		for i := range planned {
			if planned[i].Key == keyStatus {
				planned[i].Value = "hijacked"
			}
		}
	}}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	got := e.Write(ctxOf(t), writerName, planWriting(
		resolve.Tag{Key: keyStatus, Value: "final"},
	))

	// Control: the binding really did write its own value, not the plan's.
	if v, _, _ := s.get(keyStatus); v != "hijacked" {
		t.Fatalf("artifact holds %q for %q; the fixture did not substitute its own "+
			"value, so this probe is testing the wrong thing", v, keyStatus)
	}

	if !got.Refused() {
		t.Fatalf("Write reported SUCCESS (Written=%+v) after the binding rewrote the "+
			"planned value from %q to %q in place; read-back compared the artifact "+
			"against the binding's OWN choice rather than the plan — the oracle is "+
			"self-referential", got.Written, "final", "hijacked")
	}
	if got.Refusal.Class != accessor.ClassReadBackMismatch {
		t.Errorf("Write refused %q; want read_back_mismatch — the artifact demonstrably "+
			"holds %q where the plan said %q", got.Refusal.Class, "hijacked", "final")
	}

	// The refusal's own diagnostic must report the PLAN, not the binding's
	// substitution: `Expected` is what the caller asked for.
	if v, ok := seamValueOf(got.Refusal.Expected, keyStatus); !ok || v != "final" {
		t.Errorf("Refusal.Expected carries %q = (%q, present=%v); want the PLANNED "+
			"value %q — a diagnostic echoing the binding's substitution reports the "+
			"mismatch as no mismatch at all", keyStatus, v, ok, "final")
	}

	// Nothing may be reported as verified.
	if len(got.Written) != 0 {
		t.Errorf("WriteResult.Written = %+v on a refusal; want empty", got.Written)
	}
}
