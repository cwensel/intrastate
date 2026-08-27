package accessor_test

// RDR 0004 — capability, identity, selection, and the caller-supplied
// artifact rule (REQ-1..REQ-12).
//
// The load-bearing discriminator in this file is REQ-4: capability is
// PART of the accessor identity, so a same-id read/write pair is two
// identities and MUST NOT be reported as multiply-bound. An
// implementation keying bindings on the name alone passes every other
// test here and fails that one.

import (
	"context"
	"slices"
	"testing"

	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/table"
)

// REQ-1: "Every accessor definition MUST declare exactly one capability:
// read, gate, or write."
// HAPPY PATH
//
// The vocabulary is closed at three. A fourth member would let a
// definition declare authority the validator cannot reason about.
func TestReq1_CapabilityVocabularyIsClosedAtThree(t *testing.T) {
	got := accessor.Capabilities()
	want := []accessor.Capability{accessor.CapRead, accessor.CapGate, accessor.CapWrite}

	if len(got) != len(want) {
		t.Fatalf("Capabilities() = %v (%d members); want exactly the three %v",
			got, len(got), want)
	}
	for _, w := range want {
		if !slices.Contains(got, w) {
			t.Errorf("Capabilities() = %v; missing the declared capability %q", got, w)
		}
	}
	for _, g := range got {
		if !slices.Contains(want, g) {
			t.Errorf("Capabilities() carries %q, which is outside the closed set %v", g, want)
		}
	}
}

// REQ-7: "the canonical names are \"read accessor\", \"gate accessor\",
// and \"write accessor\". Rejected names: \"hook\" and \"action\" because
// they imply arbitrary transition callbacks."
// INPUT EDGE
//
// Binding on identifier spelling: the wire tokens are `read`, `gate`, and
// `write`, matching RDR 0002's capability tables. "hook" and "action" are
// not capability spellings.
func TestReq7_CapabilitySpellingsAreTheCanonicalNames(t *testing.T) {
	for _, c := range []struct {
		got  accessor.Capability
		want string
	}{
		{accessor.CapRead, "read"},
		{accessor.CapGate, "gate"},
		{accessor.CapWrite, "write"},
	} {
		if string(c.got) != c.want {
			t.Errorf("capability spelling = %q; want %q — the canonical name, "+
				"never \"hook\" or \"action\"", c.got, c.want)
		}
	}
	// Each canonical name must be a MEMBER of the enumerated vocabulary:
	// an empty set trivially excludes "hook" and "action" too.
	for _, want := range []accessor.Capability{
		accessor.CapRead, accessor.CapGate, accessor.CapWrite,
	} {
		if !slices.Contains(accessor.Capabilities(), want) {
			t.Errorf("Capabilities() = %v; the canonical name %q must be a "+
				"member, or excluding the rejected names is vacuous",
				accessor.Capabilities(), want)
		}
	}
	for _, rejected := range []accessor.Capability{"hook", "action"} {
		if slices.Contains(accessor.Capabilities(), rejected) {
			t.Errorf("Capabilities() admits %q; the name is rejected because it "+
				"implies arbitrary transition callbacks", rejected)
		}
	}
}

// REQ-2: "Runtime execution MUST reject any attempt to use an accessor
// for a different capability than the one declared." — the class is
// `capability_mismatch`.
// ADVERSARIAL
//
// The writer is bound; the invocation asks for it as a READ. The
// discriminator is the class name: an implementation that reports
// `unknown_accessor` here has conflated "not bound" with "bound under
// another capability", and the two are different defects.
func TestReq2_OffCapabilityInvocationIsCapabilityMismatch(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	w := &writeBinding{store: s}
	e := accessor.NewExecutor(registryOf(writerDef(w, keyStatus)), artifactsOf())

	got := e.Read(ctxOf(t), writerName)

	r := mustReadRefuse(t, got, accessor.ClassCapabilityMismatch)
	if r.Accessor != writerName {
		t.Errorf("refusal names accessor %q; want %q", r.Accessor, writerName)
	}
	if w.Invocations() != 0 {
		t.Errorf("the write binding ran %d times; an off-capability invocation "+
			"must be rejected BEFORE the binding is reached", w.Invocations())
	}
}

// REQ-6: "Missing or multiply-bound accessors are validation failures." —
// at runtime an unbound name is the refusal class `unknown_accessor`.
// [0002-delivered in part] — the loader's binding counts cover the
// reader/writer key bindings; the identity-triple check is this RDR's.
// ADVERSARIAL
func TestReq6_UnboundNameIsUnknownAccessorAtRuntime(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus)

	got := e.Read(ctxOf(t), "missing.read")

	r := mustReadRefuse(t, got, accessor.ClassUnknownAccessor)
	if r.Accessor != "missing.read" {
		t.Errorf("refusal names accessor %q; want the unbound name %q",
			r.Accessor, "missing.read")
	}
}

// REQ-3: "Within one flow, each `(flow id, accessor name, capability)`
// identity MUST resolve to exactly one accessor binding. Missing and
// multiply-bound identities MUST fail validation before resolution."
// ADVERSARIAL
func TestReq3_IdentityTripleResolvesToExactlyOneBinding(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})

	t.Run("multiply_bound_same_triple", func(t *testing.T) {
		// Two bindings under the SAME (flow, name, capability) triple.
		first := readerDef(&readBinding{store: s}, keyStatus)
		second := readerDef(&readBinding{store: s}, keyStatus)

		fs := accessor.Validate(registryOf(first, second),
			[]accessor.Identity{first.Identity})

		if !hasCode(fs, accessor.CodeMultiplyBoundAccessor) {
			t.Errorf("validation codes = %v; want %q — two bindings share the "+
				"identity triple %+v", findingCodes(fs),
				accessor.CodeMultiplyBoundAccessor, first.Identity)
		}
	})

	t.Run("missing_identity", func(t *testing.T) {
		reg := registryOf(readerDef(&readBinding{store: s}, keyStatus))
		referenced := accessor.Identity{
			Flow: flowID, Name: "absent.read", Capability: accessor.CapRead,
		}

		fs := accessor.Validate(reg, []accessor.Identity{referenced})

		if !hasCode(fs, accessor.CodeMissingAccessor) {
			t.Errorf("validation codes = %v; want %q — the referenced identity "+
				"%+v has no binding", findingCodes(fs),
				accessor.CodeMissingAccessor, referenced)
		}
	})
}

// REQ-4: "an accessor is identified by `(flow id, accessor name,
// capability)`. Capability is part of the identity, so the same id may
// appear in both `[read.x]` and `[write.x]` — two identities, not a
// rebinding (JDR 0001 §D7(ii)); only a second binding of the *same*
// triple is multiply-bound."
// BOUNDARY
//
// This is the discriminating case for the whole identity rule. An
// implementation keying on the name alone reports `multiply_bound_accessor`
// here — the "cannot be rebound" reading §D7 explicitly retired.
func TestReq4_SameIdUnderTwoCapabilitiesIsTwoIdentitiesNotARebinding(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})

	const sharedID = "state"
	reader := accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: sharedID, Capability: accessor.CapRead},
		Accessor: table.Accessor{
			Role: stateRole, Path: statePath,
			Keys: []string{keyStatus}, Timeout: fixtureTimeout.String(),
		},
		Binding: &readBinding{store: s},
	}
	writer := accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: sharedID, Capability: accessor.CapWrite},
		Accessor: table.Accessor{
			Role: stateRole, Path: statePath,
			Keys: []string{keyStatus}, Timeout: fixtureTimeout.String(), ReadBack: true,
		},
		Binding: &writeBinding{store: s},
	}

	reg := registryOf(reader, writer)

	// Positive precondition: BOTH identities must resolve. Without this,
	// a validator that binds nothing at all would satisfy the
	// not-multiply-bound assertion vacuously.
	if _, ok := reg.Lookup(sharedID, accessor.CapRead); !ok {
		t.Fatalf("Lookup(%q, read) found nothing; [read.%s] is bound",
			sharedID, sharedID)
	}
	if _, ok := reg.Lookup(sharedID, accessor.CapWrite); !ok {
		t.Fatalf("Lookup(%q, write) found nothing; [write.%s] is bound",
			sharedID, sharedID)
	}

	fs := accessor.Validate(reg, []accessor.Identity{reader.Identity, writer.Identity})

	if hasCode(fs, accessor.CodeMultiplyBoundAccessor) {
		t.Errorf("validation reported %q for the same id under two capability "+
			"tables; capability is PART of the identity, so [read.%s] and "+
			"[write.%s] are TWO identities, not a rebinding (codes = %v)",
			accessor.CodeMultiplyBoundAccessor, sharedID, sharedID, findingCodes(fs))
	}
	if len(fs) != 0 {
		t.Errorf("validation codes = %v; want the empty set for a legal "+
			"same-id read/write pair", findingCodes(fs))
	}
}

// REQ-5: "the executor selects a binding by capability table and id: an
// invocation that needs capability X selects only from `[X.<id>]`, never
// from a same-id entry under another table."
// BOUNDARY
//
// The negative control is REQ-4's legal pair: with both tables carrying
// the id, a read must reach the READ binding and not the write one.
func TestReq5_SelectionIsByCapabilityTableAndId(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})

	const sharedID = "state"
	rb := &readBinding{store: s}
	wb := &writeBinding{store: s}
	reader := accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: sharedID, Capability: accessor.CapRead},
		Accessor: table.Accessor{
			Role: stateRole, Path: statePath,
			Keys: []string{keyStatus}, Timeout: fixtureTimeout.String(),
		},
		Binding: rb,
	}
	writer := accessor.Definition{
		Identity: accessor.Identity{Flow: flowID, Name: sharedID, Capability: accessor.CapWrite},
		Accessor: table.Accessor{
			Role: stateRole, Path: statePath,
			Keys: []string{keyStatus}, Timeout: fixtureTimeout.String(), ReadBack: true,
		},
		Binding: wb,
	}
	reg := registryOf(reader, writer)

	t.Run("lookup_selects_per_capability", func(t *testing.T) {
		gotRead, ok := reg.Lookup(sharedID, accessor.CapRead)
		if !ok {
			t.Fatalf("Lookup(%q, read) found nothing; the read table binds it", sharedID)
		}
		if gotRead.Identity.Capability != accessor.CapRead {
			t.Errorf("Lookup(%q, read) returned a %q binding; selection must never "+
				"cross to a same-id entry under another table",
				sharedID, gotRead.Identity.Capability)
		}

		gotWrite, ok := reg.Lookup(sharedID, accessor.CapWrite)
		if !ok {
			t.Fatalf("Lookup(%q, write) found nothing; the write table binds it", sharedID)
		}
		if gotWrite.Identity.Capability != accessor.CapWrite {
			t.Errorf("Lookup(%q, write) returned a %q binding", sharedID,
				gotWrite.Identity.Capability)
		}
	})

	t.Run("read_invocation_reaches_the_read_binding", func(t *testing.T) {
		e := accessor.NewExecutor(reg, artifactsOf())

		got := e.Read(ctxOf(t), sharedID)

		mustReadSucceed(t, got)
		if rb.reads != 1 {
			t.Errorf("the read binding ran %d times; want exactly 1", rb.reads)
		}
		if wb.Invocations() != 0 {
			t.Errorf("the write binding ran %d times on a READ invocation; "+
				"selection crossed capability tables", wb.Invocations())
		}
	})
}

// REQ-10: "Accessors MUST operate on caller-supplied artifact roles."
// HAPPY PATH
//
// The artifact the binding sees is the one the CALLER supplied under that
// role — same role name and same path — not one the executor derived.
func TestReq10_BindingReceivesTheCallerSuppliedArtifactRole(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	seen := &roleWitness{}
	e := accessor.NewExecutor(
		registryOf(readerDef(&witnessReadBinding{store: s, witness: seen}, keyStatus)),
		artifactsOf())

	mustReadSucceed(t, e.Read(ctxOf(t), readerName))

	if seen.role != stateRole || seen.path != statePath {
		t.Errorf("binding saw artifact {role=%q path=%q}; want the caller-supplied "+
			"{role=%q path=%q}", seen.role, seen.path, stateRole, statePath)
	}
}

// REQ-11: "The accessor executor MUST NOT discover authoritative
// artifacts from ambient process state." — a definition attempting
// ambient artifact discovery MUST fail validation before execution.
// ADVERSARIAL
func TestReq11_AmbientArtifactDiscoveryFailsValidation(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := readerDef(&readBinding{store: s}, keyStatus)
	def.AmbientDiscovery = true

	fs := accessor.Validate(registryOf(def), []accessor.Identity{def.Identity})

	if !hasCode(fs, accessor.CodeAmbientArtifactDiscovery) {
		t.Errorf("validation codes = %v; want %q — the definition discovers its "+
			"artifact from ambient process state", findingCodes(fs),
			accessor.CodeAmbientArtifactDiscovery)
	}
}

// REQ-12: "Artifact role not supplied" is a refusal, and the class minted
// is `execution_failure`. Restated normatively: "An unsupplied or
// unreadable artifact role is **not** a separate class — it is
// `execution_failure`, matching the Disposition Table."
// INPUT EDGE
//
// The discriminator is that no NEW class is minted. An implementation
// adding `artifact_unavailable` fails here — and FM says `clierr` has no
// artifact-unavailability code to reuse either.
func TestReq12_UnsuppliedArtifactRoleIsExecutionFailureNotItsOwnClass(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	// The registry names role "state"; the caller supplies NOTHING.
	e := accessor.NewExecutor(
		registryOf(readerDef(&readBinding{store: s}, keyStatus)),
		accessor.Artifacts{})

	got := e.Read(ctxOf(t), readerName)

	r := mustReadRefuse(t, got, accessor.ClassExecutionFailure)
	if r.Role != stateRole {
		t.Errorf("refusal names role %q; want the unsupplied role %q for diagnosis",
			r.Role, stateRole)
	}
	if slices.Contains(accessor.RefusalClasses(), accessor.RefusalClass("artifact_unavailable")) {
		t.Error("the refusal-class set carries a separate artifact-unavailability " +
			"class; an unsupplied or unreadable role is execution_failure")
	}
}

// --- witness binding -----------------------------------------------------

type roleWitness struct {
	role string
	path string
}

type witnessReadBinding struct {
	store   *store
	witness *roleWitness
}

func (b *witnessReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *witnessReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	b.witness.role = art.Role
	b.witness.path = art.Path
	inner := &readBinding{store: b.store}
	return inner.Read(ctx, art, requested)
}
