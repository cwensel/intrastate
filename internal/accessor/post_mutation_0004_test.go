package accessor_test

// RDR 0004 — read-back incompleteness and post-mutation reporting
// (REQ-61..REQ-67, REQ-102, REQ-107, REQ-120, REQ-125).
//
// CA A10's premise, honoured by the fixtures: "no fixture double can
// exercise it — a re-read that cannot fail independently of the write has
// nothing to report." The write binding and the read-back reader are
// SEPARATE objects over the same store, so `readBackFailure` makes the
// re-read fail while the mutation has already landed.
//
// The three-way discrimination this file holds:
//
//	success                → the artifact is right, verified
//	read_back_mismatch     → the artifact is WRONG (verification RAN)
//	read_back_incomplete   → the verification DID NOT RUN; the mutation
//	                         may have been applied and was not verified
//
// A layer that collapses the last two into one "write failed" shape fails
// ORA 8's negative control (REQ-67).

import (
	"context"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
)

// REQ-61: "The read-back re-read is subject to read completeness." — the
// re-read goes through the READ path, not a map clone; the spike's clone
// "cannot fail independently of the write" (CA A10).
// REQ-125: "the read-back fixture MUST be able to fail independently of
// the write."
// REQ-62: "If it cannot read a key it must compare, the write MUST be
// reported as `read_back_incomplete` — the verification did not run — and
// MUST NOT be reported as `read_back_mismatch`, which asserts the
// artifact is wrong, nor as success."
// DOMAIN EDGE (A10 / A9 rule 4, unwitnessed by the spike)
func TestReq61_ReadBackThatCannotReadAComparedKeyIsReadBackIncomplete(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	w := &writeBinding{store: s}
	// The re-read fails on the key it must compare — AFTER the mutation
	// has landed. This is only expressible because the reader is a
	// separate object from the writer.
	r := &readBackFailure{store: s, failOn: keyStatus}
	e := accessor.NewExecutor(
		registryOf(readerDef(r, keyStatus), writerDef(w, keyStatus)), artifactsOf())

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	if !got.Refused() {
		t.Fatalf("write reported SUCCESS; the re-read could not read %q, so the "+
			"verification did not run and success asserts state that was never "+
			"established", keyStatus)
	}
	if got.Refusal.Class == accessor.ClassReadBackMismatch {
		t.Fatalf("refusal class = %q; the re-read could not READ the key — "+
			"mismatch asserts the artifact is WRONG, which the verification "+
			"never established", accessor.ClassReadBackMismatch)
	}
	mustWriteRefuse(t, got, accessor.ClassReadBackIncomplete)

	// The mutation DID land: the write command ran to completion.
	if s.tags[keyStatus] != "Final" {
		t.Errorf("artifact holds %q = %q; the write command ran before the "+
			"re-read failed", keyStatus, s.tags[keyStatus])
	}
	if r.reads == 0 {
		t.Error("the read-back never invoked the read path; the re-read is " +
			"subject to read completeness, so it goes through `read` rather " +
			"than cloning the tag map")
	}
}

// REQ-63: "`read_back_incomplete` and a post-mutation `timeout` are
// reported after the write command already ran. The refusal MUST be
// understood as \"the mutation may have been applied and was not
// verified\""
// REQ-64: "it MUST NOT be represented to callers as a write that did not
// occur."
// REQ-120: "The refusal is `read_back_incomplete` and carries the
// applied-but-unverified sense rather than reading as a write that did
// not occur."
// BOUNDARY (A10)
//
// The sense is asserted POSITIVELY on the returned value, not inferred
// from a log line: `Applied()` reports whether the write command already
// ran when the refusal was minted.
func TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense(t *testing.T) {
	t.Run("read_back_incomplete_is_applied_but_unverified", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		r := &readBackFailure{store: s, failOn: keyStatus}
		e := accessor.NewExecutor(
			registryOf(readerDef(r, keyStatus), writerDef(w, keyStatus)), artifactsOf())

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		ref := mustWriteRefuse(t, got, accessor.ClassReadBackIncomplete)
		if !ref.Applied() {
			t.Error("Applied() = false; the write command already ran, so the " +
				"refusal means \"the mutation may have been applied and was " +
				"not verified\" — it MUST NOT read as a write that did not occur")
		}
	})

	t.Run("post_mutation_timeout_is_applied_but_unverified", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		// The mutation lands, THEN the invocation runs past its deadline.
		w := &writeBinding{store: s, delay: 10 * fixtureTimeout}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		ref := mustWriteRefuse(t, got, accessor.ClassTimeout)
		if s.tags[keyStatus] != "Final" {
			t.Fatalf("the fixture's mutation did not land before the deadline "+
				"(%q = %q); this case requires the write command to have RUN",
				keyStatus, s.tags[keyStatus])
		}
		if !ref.Applied() {
			t.Error("Applied() = false for a timeout raised AFTER the write " +
				"command mutated the artifact; the caller must be able to tell " +
				"\"may have been applied, unverified\" from \"did not occur\"")
		}
	})

	t.Run("control_a_pre_mutation_failure_did_not_occur", func(t *testing.T) {
		// The write command itself failed: nothing was applied. If every
		// write refusal reported Applied()=true, the distinction would be
		// vacuous and the caller's recovery burden would widen to every
		// failure.
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, failErr: errBindingFailed}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		ref := mustWriteRefuse(t, got, accessor.ClassExecutionFailure)
		if ref.Applied() {
			t.Error("Applied() = true for a write whose command FAILED before " +
				"mutating anything; a layer that reports every write refusal as " +
				"possibly-applied makes the distinction vacuous")
		}
		if s.tags[keyStatus] != "Draft" {
			t.Errorf("artifact holds %q = %q; the failing command mutated nothing",
				keyStatus, s.tags[keyStatus])
		}
	})
}

// REQ-65: "The accessor layer MUST NOT attempt compensation: it does not
// retry the write, does not undo it, and does not re-derive the
// artifact's state." — asserted "on the invocation count — a layer that
// retries or undoes shows a second write" (ORA 8).
// REQ-107 / TS 8: "The test asserts the accessor issued exactly one write
// invocation and performed no undo, retry, or re-derivation."
// ADVERSARIAL (A10)
func TestReq65_NoCompensationAfterAPostMutationRefusal(t *testing.T) {
	for name, mk := range map[string]func(*store) (*writeBinding, accessor.Binding, accessor.RefusalClass){
		"read_back_incomplete": func(s *store) (*writeBinding, accessor.Binding, accessor.RefusalClass) {
			return &writeBinding{store: s},
				&readBackFailure{store: s, failOn: keyStatus},
				accessor.ClassReadBackIncomplete
		},
		"read_back_mismatch": func(s *store) (*writeBinding, accessor.Binding, accessor.RefusalClass) {
			return &writeBinding{store: s, corrupt: func(st *store) {
					st.tags[keyStatus] = "corrupt"
				}},
				&readBinding{store: s},
				accessor.ClassReadBackMismatch
		},
	} {
		t.Run(name, func(t *testing.T) {
			s := newStore(map[string]string{keyStatus: "Draft"})
			w, reader, want := mk(s)
			e := accessor.NewExecutor(
				registryOf(readerDef(reader, keyStatus), writerDef(w, keyStatus)),
				artifactsOf())

			before := map[string]string{}
			for k, v := range s.tags {
				before[k] = v
			}

			got := e.Write(ctxOf(t), writerName,
				planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

			mustWriteRefuse(t, got, want)
			if w.Invocations() != 1 {
				t.Errorf("the write binding ran %d times; want EXACTLY 1 — a "+
					"layer that retries or undoes shows a second write. "+
					"Compensation is forbidden; recovery is the caller's",
					w.Invocations())
			}
			// No undo: whatever the write left behind is still there.
			if s.tags[keyStatus] == before[keyStatus] && name == "read_back_incomplete" {
				t.Errorf("artifact %q reverted to its pre-write value %q; the "+
					"accessor layer does not undo the write", keyStatus, before[keyStatus])
			}
		})
	}
}

// REQ-67: The `read_back_mismatch` fixture is the negative control for
// post-mutation reporting: it "must still assert the artifact is wrong
// rather than unverified — an implementation that collapses both into one
// \"write failed\" shape fails that row" (ORA 8).
// BOUNDARY
//
// The two refusals must be DISTINGUISHABLE, and in opposite senses:
// mismatch RAN the verification and found the artifact wrong; incomplete
// never ran it. Same plan, same store shape, different re-read.
func TestReq67_MismatchAndIncompleteAreDistinctNotOneWriteFailedShape(t *testing.T) {
	mismatch := func() accessor.WriteResult {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			st.tags[keyStatus] = "corrupt"
		}}
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBinding{store: s}, keyStatus),
				writerDef(w, keyStatus)), artifactsOf())
		return e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))
	}
	incomplete := func() accessor.WriteResult {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := accessor.NewExecutor(
			registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
				writerDef(w, keyStatus)), artifactsOf())
		return e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))
	}

	m := mustWriteRefuse(t, mismatch(), accessor.ClassReadBackMismatch)
	i := mustWriteRefuse(t, incomplete(), accessor.ClassReadBackIncomplete)

	if m.Class == i.Class {
		t.Fatalf("both refusals carry class %q; an implementation that collapses "+
			"them into one \"write failed\" shape cannot tell a caller whether "+
			"the artifact is WRONG or merely UNVERIFIED", m.Class)
	}
	if m.Applied() {
		t.Error("the mismatch refusal reports Applied()=true; the verification " +
			"RAN and asserts the artifact is wrong — that is not the " +
			"applied-but-unverified sense")
	}
	if !i.Applied() {
		t.Error("the incomplete refusal reports Applied()=false; the " +
			"verification did not run and the mutation may have been applied")
	}
}

// REQ-66: "Recovery is the caller's, and its safe move is to re-read
// before acting."
// REQ-62, naming half: the refusal names the key the re-read could not
// read, so the caller knows what to re-read.
// BOUNDARY
func TestReq66_ReadBackIncompleteNamesWhatTheCallerMustReRead(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	w := &writeBinding{store: s}
	e := accessor.NewExecutor(
		registryOf(readerDef(&readBackFailure{store: s, failOn: keyStatus}, keyStatus),
			writerDef(w, keyStatus)), artifactsOf())

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	r := mustWriteRefuse(t, got, accessor.ClassReadBackIncomplete)
	if len(r.Keys) == 0 {
		t.Fatal("the refusal names no keys; the caller's safe move is to " +
			"re-read before acting, which needs to know WHAT it could not verify")
	}
	found := false
	for _, k := range r.Keys {
		if k == keyStatus {
			found = true
		}
	}
	if !found {
		t.Errorf("refusal names keys %v; want %q — the key the re-read could "+
			"not compare", r.Keys, keyStatus)
	}
}

// REQ-102 / DESK step 5: the full read-back assertion set holds AT ONCE
// on one invocation — same-role re-read, subject to read completeness,
// applied-but-unverified sense, no compensation, equality for an
// assignment and absence for a clear, pre-write non-owned values
// unchanged, mismatch is a write failure.
// BOUNDARY
//
// The desk trace's claim is that these are jointly satisfiable, not that
// each holds in isolation. This walks one write with every assertion in
// force.
func TestReq102_DeskTraceStep5AssertionsHoldSimultaneously(t *testing.T) {
	s := newStore(map[string]string{
		keyStatus:  "Draft",
		keyLabels:  `["alpha"]`,
		keyProfile: "large",
	})
	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := accessor.NewExecutor(
		registryOf(readerDef(r, keyStatus, keyLabels), writerDef(w, keyStatus, keyLabels)),
		artifactsOf())

	// One plan: an assignment AND a clear on the same invocation.
	plan := planWriting(
		resolve.Tag{Key: keyStatus, Value: "Final"},
		resolve.Tag{Key: keyLabels, Value: accessor.ClearSentinel},
	)
	got := e.Write(ctxOf(t), writerName, plan)

	mustWriteSucceed(t, got)

	if s.tags[keyStatus] != "Final" {
		t.Errorf("assignment: %q = %q; want %q held exactly",
			keyStatus, s.tags[keyStatus], "Final")
	}
	if _, held := s.tags[keyLabels]; held {
		t.Errorf("clear: %q is still held as %q; a clear is a removal and its "+
			"read-back asserts ABSENCE", keyLabels, s.tags[keyLabels])
	}
	if s.tags[keyProfile] != "large" {
		t.Errorf("pre-write non-owned %q moved to %q; it must be unchanged",
			keyProfile, s.tags[keyProfile])
	}
	if w.Invocations() != 1 {
		t.Errorf("the write binding ran %d times; want exactly 1", w.Invocations())
	}
	if r.reads == 0 {
		t.Error("the read-back never went through the read path")
	}
}

// --- the independently-failing read-back binding -------------------------

// readBackFailure reads the store normally EXCEPT for failOn, which it
// reports as unreadable. It is a separate object from the write binding,
// which is what lets the re-read fail after the mutation has already
// landed (CA A10, REQ-125).
type readBackFailure struct {
	store  *store
	failOn string
	reads  int
}

func (b *readBackFailure) Capability() accessor.Capability { return accessor.CapRead }

func (b *readBackFailure) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	b.reads++
	var values []accessor.KeyValue
	var unreadable []string
	for _, key := range requested {
		if key == b.failOn {
			unreadable = append(unreadable, key)
			continue
		}
		v, absent, readable := b.store.get(key)
		if !readable {
			unreadable = append(unreadable, key)
			continue
		}
		values = append(values, accessor.KeyValue{Key: key, Value: v, Absent: absent})
	}
	return values, unreadable, nil
}
