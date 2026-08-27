package accessor_test

// RDR 0004 — write accessors, read-back verification, and the
// `<clear>`-as-removal rule (REQ-40..REQ-67).
//
// The coverage floor this file exists to hold, stated as the mutants it
// kills:
//
//   - a read-back scoped to owned tags only survives the owned-mismatch
//     fixture and dies on the non-owned one (ORA 3, REQ-57);
//   - a read-back that always reports mismatch dies on the success
//     control (ORA 3, REQ-57);
//   - a clear implemented as assignment of the literal passes a
//     containment check and dies on the ABSENCE assertion (ORA 9,
//     REQ-45/46);
//   - a read-back that asserts absence for every write dies on the
//     assignment control (ORA 9, REQ-49);
//   - a layer that collapses read_back_incomplete and read_back_mismatch
//     into one "write failed" shape dies on REQ-67's pair.

import (
	"context"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/resolve"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-40: "A write accessor MUST apply only planned owned-tag writes
// produced by a successful transition."
// REQ-50: "After a write accessor reports command-level success, the
// executor MUST re-read the same caller-supplied artifact role named by
// the write binding"
// REQ-51: "verify each planned owned tag for equality against the held
// value — a write replaces the whole value, so containment is not
// equality"
// REQ-58 / REQ-112 / MVV Round-Trip: the re-read holds each planned
// owned-tag value EXACTLY, and pre-write observed/recognized values are
// unchanged.
// HAPPY PATH
//
// This is ORA 3's control row: "a read-back that always reports mismatch
// fails this row".
func TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	w := &writeBinding{store: s}
	r := &readBinding{store: s}
	e := writeExec(t, s, w, r, keyStatus)

	plan := planWriting(resolve.Tag{Key: keyStatus, Value: "Final"})
	got := e.Write(ctxOf(t), writerName, plan)

	mustWriteSucceed(t, got)
	if v, ok := seamValueOf(got.Written, keyStatus); !ok || v != "Final" {
		t.Errorf("written tags = %+v; want %q held exactly as %q — the re-read "+
			"value must EQUAL the transition plan's expected value",
			got.Written, keyStatus, "Final")
	}
	if s.tags[keyStatus] != "Final" {
		t.Errorf("artifact holds %q = %q after the write; want %q",
			keyStatus, s.tags[keyStatus], "Final")
	}
	if s.tags[keyProfile] != "large" {
		t.Errorf("pre-write non-owned %q moved from %q to %q; observed and "+
			"recognized values present before the write must be unchanged",
			keyProfile, "large", s.tags[keyProfile])
	}
	if w.Invocations() != 1 {
		t.Errorf("the write binding ran %d times; want exactly 1", w.Invocations())
	}
}

// REQ-51: "a write replaces the whole value, so containment is not
// equality" (JDR 0001 §D7(iv)).
// REQ-52: read-back equality is BYTE equality over RDR 0002's canonical
// JSON array for a set-valued tag — "no third encoding exists"; the
// accessor MUST NOT compare set-valued tags as unordered collections.
// DOMAIN EDGE
//
// The set arm is the one a "compare as sets" implementation passes and a
// byte-equality one refuses: the artifact holds the same MEMBERS in a
// non-canonical order, which is not the §D13 form.
func TestReq51_ReadBackIsEqualityNotContainment(t *testing.T) {
	t.Run("scalar_containment_is_not_equality", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			// The planned value is CONTAINED in what is held, but the write
			// replaces the whole value, so this is a mismatch.
			st.tags[keyStatus] = "Final-pending"
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
	})

	t.Run("set_valued_read_back_is_byte_equality_over_the_D13_array", func(t *testing.T) {
		const canonical = `["alpha","beta"]`
		// Same members, non-canonical order: NOT the §D13 form, so byte
		// equality refuses. An unordered-collection comparison passes.
		const permuted = `["beta","alpha"]`

		s := newStore(map[string]string{keyLabels: `["alpha"]`})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			st.tags[keyLabels] = permuted
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyLabels)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyLabels, Value: canonical}))

		mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
	})

	t.Run("control_the_canonical_array_matches_byte_for_byte", func(t *testing.T) {
		const canonical = `["alpha","beta"]`
		s := newStore(map[string]string{keyLabels: `["alpha"]`})
		e := writeExec(t, s, &writeBinding{store: s}, &readBinding{store: s}, keyLabels)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyLabels, Value: canonical}))

		mustWriteSucceed(t, got)
	})
}

// REQ-55: "A read-back mismatch MUST be reported as a write failure." —
// "A read-back mismatch is a failure even if the write command exited
// successfully."
// REQ-57: the two mismatch arms are distinguished by WHICH tag moved:
// owned (`status`) vs non-owned (`profile`) (ORA 3).
// REQ-53: "and that observed and recognized tag values present before the
// write are unchanged."
// ADVERSARIAL
func TestReq55_ReadBackMismatchIsAWriteFailureOnBothArms(t *testing.T) {
	t.Run("owned_tag_differs", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			st.tags[keyStatus] = "corrupt"
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		r := mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
		if v, ok := seamValueOf(r.Expected, keyStatus); !ok || v != "Final" {
			t.Errorf("Expected = %+v; want %q = %q for diagnosis",
				r.Expected, keyStatus, "Final")
		}
		if v, ok := seamValueOf(r.Observed, keyStatus); !ok || v != "corrupt" {
			t.Errorf("Observed = %+v; want %q = %q for diagnosis",
				r.Observed, keyStatus, "corrupt")
		}
	})

	t.Run("non_owned_tag_changed_alongside", func(t *testing.T) {
		// The OWNED tag landed exactly as planned. Only `profile`, a
		// non-owned tag present before the write, moved. A read-back
		// scoped to owned tags only reports SUCCESS here.
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			st.tags[keyProfile] = "small"
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		if !got.Refused() {
			t.Fatalf("write reported success; the owned tag landed as planned but "+
				"the non-owned %q moved from %q to %q. A read-back scoped to "+
				"owned tags only cannot see this — read-back verifies the "+
				"protected non-owned identity too", keyProfile, "large", "small")
		}
		mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
	})
}

// REQ-60: "tags absent before the write are unconstrained; the write
// binding's own planned owned tags are excluded from this comparison by
// construction."
// BOUNDARY
//
// Two exemptions, each asserted as a SUCCESS that a too-strict read-back
// would refuse: a key that did not exist pre-write appearing afterwards,
// and the planned owned key itself changing (which it must).
func TestReq60_NonOwnedComparisonExemptsNewKeysAndThePlannedOwnedKeys(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	w := &writeBinding{store: s, corrupt: func(st *store) {
		// A key that did not exist before the write. Unconstrained.
		st.tags["annotation"] = "added"
	}}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	mustWriteSucceed(t, got)
	if s.tags[keyStatus] != "Final" {
		t.Errorf("the planned owned key %q holds %q; it is excluded from the "+
			"NON-OWNED comparison by construction and must change as planned",
			keyStatus, s.tags[keyStatus])
	}
}

// REQ-59: "This is not an undo or byte-for-byte artifact invariant" —
// "byte-for-byte artifact identity, formatting, key order, and comments
// are explicitly out of scope; the accessor observes tag values, not the
// artifact encoding."
// BOUNDARY
//
// The deliberate weakening, asserted as a success: the artifact's key
// ORDER and its unrelated encoding change, and read-back still passes
// because it compares TAG VALUES.
func TestReq59_ReadBackDoesNotAssertArtifactLevelFidelity(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	w := &writeBinding{store: s, corrupt: func(st *store) {
		// Re-materialize the same tag VALUES through a different map
		// identity: same values, different encoding-level identity.
		rebuilt := map[string]string{}
		for k, v := range st.tags {
			rebuilt[k] = v
		}
		st.tags = rebuilt
	}}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	mustWriteSucceed(t, got)
	// Positive precondition: the write ran and the re-read compared tag
	// values across the re-materialization. A write that never ran would
	// satisfy "encoding identity is out of scope" vacuously.
	if w.Invocations() != 1 {
		t.Fatalf("the write binding ran %d times; want exactly 1", w.Invocations())
	}
	if s.tags[keyStatus] != "Final" {
		t.Errorf("artifact holds %q = %q; the read-back compared VALUES across "+
			"the re-materialization and passed", keyStatus, s.tags[keyStatus])
	}
}

// REQ-54: "It MUST NOT satisfy read-back verification by discovering an
// ambient artifact or by reading an unrelated role."
// REQ-104: "Require same-role read-back verification against expected
// owned-tag values."
// ADVERSARIAL
//
// Two roles are supplied. The write binding names `state`; an UNRELATED
// role holds the planned value already. A read-back that reads the wrong
// role reports success over an artifact it never wrote.
func TestReq54_ReadBackReadsTheSameRoleNamedByTheWriteBinding(t *testing.T) {
	const otherRole = "sibling"

	written := newStore(map[string]string{keyStatus: "Draft"})
	// The decoy already holds the planned value.
	decoy := newStore(map[string]string{keyStatus: "Final"})

	// The write lands on `written`, but leaves it WRONG.
	w := &writeBinding{store: written, corrupt: func(st *store) {
		st.tags[keyStatus] = "corrupt"
	}}
	roleReader := &roleRoutedReadBinding{
		stores: map[string]*store{stateRole: written, otherRole: decoy},
	}

	reg := registryOf(readerDef(roleReader, keyStatus), writerDef(w, keyStatus))
	arts := accessor.Artifacts{
		stateRole: {Role: stateRole, Path: statePath},
		otherRole: {Role: otherRole, Path: "flows/sibling.toml"},
	}
	e := accessor.NewExecutor(reg, arts)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	if !got.Refused() {
		t.Fatalf("write reported success; role %q holds %q = %q while the "+
			"unrelated role %q happens to hold the planned value. Read-back "+
			"must re-read the SAME caller-supplied role the write binding names",
			stateRole, keyStatus, written.tags[keyStatus], otherRole)
	}
	mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
	if slices.Contains(roleReader.rolesSeen, otherRole) {
		t.Errorf("read-back read roles %v; it must never read the unrelated "+
			"role %q", roleReader.rolesSeen, otherRole)
	}
}

// REQ-41: "It MUST NOT write observed or recognized tags." — a definition
// declaring a write to a non-owned tag MUST fail validation.
// ADVERSARIAL
func TestReq41_WriterNamingANonOwnedTagFailsValidation(t *testing.T) {
	s := newStore(map[string]string{keyProfile: "large"})
	// `profile` is observed: registryOf declares owned = {status, labels}.
	def := writerDef(&writeBinding{store: s}, keyProfile)

	fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

	if !hasCode(fs, accessor.CodeWriteNonOwnedTag) {
		t.Errorf("validation codes = %v; want %q — the writer names the "+
			"non-owned tag %q", findingCodes(fs),
			accessor.CodeWriteNonOwnedTag, keyProfile)
	}
}

// REQ-56: read-back is MANDATORY, not optional: "writers additionally
// carry `read_back = true`, since read-back is mandatory, not optional".
// [0002-delivered] — the loader already refuses a write entry without
// `read_back = true`.
// INPUT EDGE
//
// The shipped surface is cited: `internal/table/load.go::accessorTable`
// refuses the entry at load. This RDR's boundary asserts the same defect
// carries its own named code here.
func TestReq56_WriterWithoutReadBackFailsValidation(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := writerDef(&writeBinding{store: s}, keyStatus)
	def.Accessor.ReadBack = false

	fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

	if !hasCode(fs, accessor.CodeMissingWriteReadBack) {
		t.Errorf("validation codes = %v; want %q — read-back is mandatory, not "+
			"optional", findingCodes(fs), accessor.CodeMissingWriteReadBack)
	}
}

// REQ-42: "write accessors execute only from a successful transition
// plan." — no write runs on `no_match`, `ambiguous_match`, or any other
// refusal.
// REQ-43: an ESCAPED plan reaches the write accessor with an empty write
// set and performs no write (deviations D1, RDR 0009 `0009:1515-1527`;
// JDR 0001 §D9 "an escaped plan is never gated").
// ADVERSARIAL
func TestReq42_NoWriteRunsWithoutASuccessfulPlansWrites(t *testing.T) {
	t.Run("escaped_plan_carries_an_empty_write_set_and_writes_nothing", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		plan := escapedPlan()
		if len(plan.Writes) != 0 {
			t.Fatalf("fixture escaped plan carries writes %+v; RDR 0009 fixes an "+
				"escaped plan's write set as EMPTY", plan.Writes)
		}
		if len(plan.NextTags) != 0 {
			t.Fatalf("fixture escaped plan carries NextTags %+v; an escaped plan "+
				"advances no state", plan.NextTags)
		}

		got := e.Write(ctxOf(t), writerName, plan)

		if w.Invocations() != 0 {
			t.Errorf("the write binding ran %d times on an ESCAPED plan; an "+
				"escaped plan reaches the write accessor with an empty write "+
				"set and performs no write", w.Invocations())
		}
		if got.Refused() {
			t.Errorf("an escaped plan's empty write set refused with class %q; "+
				"performing no write is not a failure", got.Refusal.Class)
		}
		if s.tags[keyStatus] != "Draft" {
			t.Errorf("artifact %q moved to %q; no write ran",
				keyStatus, s.tags[keyStatus])
		}
	})

	t.Run("control_a_successful_plan_with_writes_does_write", func(t *testing.T) {
		// Without this control, an executor that never writes at all
		// passes both no-write assertions vacuously.
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		mustWriteSucceed(t, got)
		if w.Invocations() != 1 {
			t.Errorf("the write binding ran %d times for a SUCCESSFUL plan "+
				"carrying writes; want exactly 1 — an executor that never "+
				"writes passes the escaped-plan assertions vacuously",
				w.Invocations())
		}
	})

	t.Run("only_the_plans_writes_are_applied_never_its_next_tags", func(t *testing.T) {
		// The plan's WRITES are the only tags applied. `NextTags` names
		// the state the matched row transitions TO — the resolver's
		// record of the advance — and is not an instruction to mutate
		// the artifact (deviations D1 carrying RDR 0009's obligation at
		// `0009:1515-1527`, Type TEST-FIXTURE; REQ-40, REQ-43).
		//
		// Every other plan in this suite comes from `planWriting`, which
		// mirrors the two fields, so an executor applying `NextTags`
		// instead of `Writes` is indistinguishable there. Here the two
		// DIFFER: `labels` is named by `NextTags` alone.
		//
		// `labels` is an owned tag of the model AND among this writer's
		// declared keys, so it is a key the accessor MAY write — the
		// refusal arms cannot account for it being left alone. Only the
		// Writes/NextTags distinction can.
		s := newStore(map[string]string{keyStatus: "Draft", keyLabels: `["alpha"]`})
		w := &writeBinding{store: s}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus, keyLabels)

		plan := planAdvancing(
			[]resolve.Tag{{Key: keyStatus, Value: "Final"}},
			[]resolve.Tag{
				{Key: keyStatus, Value: "Final"},
				{Key: keyLabels, Value: `["beta"]`},
			},
		)

		got := e.Write(ctxOf(t), writerName, plan)

		mustWriteSucceed(t, got)

		// The artifact: the NextTags-only key is untouched.
		if s.tags[keyLabels] != `["alpha"]` {
			t.Errorf("artifact %q = %q after the write; want %q unchanged — %q is "+
				"named by the plan's NextTags and NOT by its Writes, and only "+
				"planned WRITES are applied", keyLabels, s.tags[keyLabels],
				`["alpha"]`, keyLabels)
		}
		if s.tags[keyStatus] != "Final" {
			t.Errorf("artifact %q = %q; want %q — the planned write must still land",
				keyStatus, s.tags[keyStatus], "Final")
		}

		// The record: `Written` names exactly the planned writes. A
		// success record naming a key the plan did not write reports the
		// accessor as having established state it never planned.
		if k := tagKeysOf(got.Written); !slices.Equal(k, []string{keyStatus}) {
			t.Errorf("Written names %v; want exactly %v — the write record is the "+
				"plan's Writes, never its NextTags", k, []string{keyStatus})
		}
	})

	t.Run("empty_write_set_writes_nothing", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		e.Write(ctxOf(t), writerName, planWriting())

		if w.Invocations() != 0 {
			t.Errorf("the write binding ran %d times for a plan with no writes",
				w.Invocations())
		}
	})
}

// REQ-45: "A planned write of `<clear>` MUST remove the key from the
// artifact"
// REQ-46: "its read-back MUST assert the key is absent"
// DOMAIN EDGE (A11, unwitnessed by the spike)
//
// ORA 9's oracle: "the assertion is on the key's *absence* in the re-read
// — a stored `<clear>` string makes the key present and the row fails".
func TestReq45_ClearingWriteRemovesTheKeyAndReadBackAssertsAbsence(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
	w := &writeBinding{store: s}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: table.ClearSentinel}))

	mustWriteSucceed(t, got)
	if v, held := s.tags[keyStatus]; held {
		t.Errorf("the artifact still holds %q = %q after a `%s` write; a clear "+
			"is a REMOVAL, not an assignment", keyStatus, v, table.ClearSentinel)
	}
	if s.tags[keyProfile] != "large" {
		t.Errorf("the clear also moved the non-owned %q to %q",
			keyProfile, s.tags[keyProfile])
	}
}

// REQ-46: "a re-read that still holds the key, including as the literal
// string `<clear>`, is `read_back_mismatch`."
// REQ-108: "a clear is a removal, read-back asserts absence, and MVV
// Scenario 9 asserts a stored literal fails."
// ADVERSARIAL (A11)
//
// This is THE defect the reserved-sentinel clause exists to catch: the
// binding stores the literal, so a containment or equality check against
// the planned value `<clear>` passes GREEN over a tag that was never
// removed. Only an ABSENCE assertion fails it.
func TestReq46_AClearThatStoresTheLiteralIsReadBackMismatch(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	w := &writeBinding{store: s, storeLiteral: true}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: table.ClearSentinel}))

	if !got.Refused() {
		t.Fatalf("write reported success; the artifact holds %q = %q — the "+
			"literal was STORED, not removed. A read-back comparing the held "+
			"value against the planned `%s` passes green over a tag that was "+
			"never removed; only an ABSENCE assertion catches it",
			keyStatus, s.tags[keyStatus], table.ClearSentinel)
	}
	mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
}

// REQ-46, non-literal half: "a re-read that still holds the key,
// including as the literal string `<clear>`, is `read_back_mismatch`."
// ADVERSARIAL (A11)
func TestReq46_AClearWhoseKeySurvivesWithAnyValueIsReadBackMismatch(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	w := &writeBinding{store: s, corrupt: func(st *store) {
		// The removal happened, then something put the key back.
		st.tags[keyStatus] = "Draft"
	}}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: table.ClearSentinel}))

	mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
}

// REQ-47: "Clearing a key the artifact does not hold MUST succeed." — the
// idempotent clear.
// BOUNDARY (A11)
func TestReq47_ClearingAKeyTheArtifactDoesNotHoldSucceeds(t *testing.T) {
	// The artifact never held `status`.
	s := newStore(map[string]string{keyProfile: "large"})
	w := &writeBinding{store: s}
	e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: table.ClearSentinel}))

	mustWriteSucceed(t, got)
	// Positive precondition: the clear was ACTUALLY attempted. An
	// executor that ignores every plan reports success over an artifact
	// it never touched.
	if w.Invocations() != 1 {
		t.Fatalf("the write binding ran %d times; want exactly 1 — an "+
			"idempotent clear still runs the write, it just removes nothing",
			w.Invocations())
	}
	if _, held := s.tags[keyStatus]; held {
		t.Errorf("the artifact now holds %q; an idempotent clear removes nothing "+
			"and adds nothing", keyStatus)
	}
	if s.tags[keyProfile] != "large" {
		t.Errorf("the idempotent clear moved the untouched %q to %q",
			keyProfile, s.tags[keyProfile])
	}
}

// REQ-49: the clear cases' negative control is the ASSIGNMENT fixture:
// "The assignment fixture from scenario 3 … must still assert presence
// and equality — a read-back that asserts absence for every write fails
// it" (ORA 9).
// BOUNDARY
func TestReq49_AssignmentControlStillAssertsPresenceAndEquality(t *testing.T) {
	t.Run("assignment_requires_presence", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		mustWriteSucceed(t, got)
		if s.tags[keyStatus] != "Final" {
			t.Errorf("artifact holds %q = %q; an assignment asserts PRESENCE and "+
				"equality — a read-back that asserts absence for every write "+
				"fails this control", keyStatus, s.tags[keyStatus])
		}
	})

	t.Run("an_assignment_whose_key_vanished_is_a_mismatch", func(t *testing.T) {
		s := newStore(map[string]string{keyStatus: "Draft"})
		w := &writeBinding{store: s, corrupt: func(st *store) {
			delete(st.tags, keyStatus)
		}}
		e := writeExec(t, s, w, &readBinding{store: s}, keyStatus)

		got := e.Write(ctxOf(t), writerName,
			planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

		mustWriteRefuse(t, got, accessor.ClassReadBackMismatch)
	})
}

// REQ-44: "`<clear>` is a reserved tag value (JDR 0001 §D5)."
// [0002-delivered] — 0002 refuses `<clear>` as an authored tag value at
// load; this RDR consumes the same constant rather than re-spelling it.
// HAPPY PATH
func TestReq44_ClearSentinelIsRDR0002sReservedConstant(t *testing.T) {
	if accessor.ClearSentinel != table.ClearSentinel {
		t.Errorf("accessor.ClearSentinel = %q; want RDR 0002's %q — this RDR "+
			"consumes the reserved constant, it does not mint a second one",
			accessor.ClearSentinel, table.ClearSentinel)
	}
	if accessor.ClearSentinel != "<clear>" {
		t.Errorf("the reserved sentinel spells %q; want %q (JDR 0001 §D5)",
			accessor.ClearSentinel, "<clear>")
	}
	if !accessor.IsClear(table.ClearSentinel) {
		t.Errorf("IsClear(%q) = false; the sentinel names a removal",
			table.ClearSentinel)
	}
	for _, notClear := range []string{"", "clear", "Final", "<absent>", " <clear>"} {
		if accessor.IsClear(notClear) {
			t.Errorf("IsClear(%q) = true; only the exact reserved sentinel is a "+
				"removal", notClear)
		}
	}
}

// --- role-routed read binding --------------------------------------------

// roleRoutedReadBinding serves a different store per artifact role and
// records which roles it was asked for, so a read-back that strays to an
// unrelated role is observable.
type roleRoutedReadBinding struct {
	stores    map[string]*store
	rolesSeen []string
}

func (b *roleRoutedReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *roleRoutedReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	b.rolesSeen = append(b.rolesSeen, art.Role)
	s, ok := b.stores[art.Role]
	if !ok {
		return nil, nil, errBindingFailed
	}
	inner := &readBinding{store: s}
	return inner.Read(ctx, art, requested)
}
