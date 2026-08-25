package accessor_test

// RDR 0004 — validation before execution: the eight arms, each asserting
// its OWN named code (REQ-83..REQ-85, REQ-113), plus replay stability
// (REQ-86..REQ-89) and the ownership boundaries this RDR must not cross
// (REQ-8, REQ-9, REQ-81, REQ-90..REQ-95, REQ-99, REQ-126).

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/newcoinc/intrastate/internal/accessor"
	"github.com/newcoinc/intrastate/internal/resolve"
	"github.com/newcoinc/intrastate/internal/table"
)

// REQ-83: "validation rejects unknown accessor names, missing or
// multiply-bound accessor identities, capability mismatches, writes to
// non-owned tags, missing timeout/read-back metadata, non-positive
// timeouts, a missing or empty read requested-key set, and ambient
// artifact discovery before resolution."
// REQ-84: "eight validation arms, each asserting its own named code."
// REQ-85: the oracle asserts each fixture's OWN named code — "not merely
// that validation returned non-empty" — with the valid fixture returning
// the empty set as the negative control, "so a validator that rejects
// everything fails".
// REQ-113 / TS 1.
// BOUNDARY
//
// This is MVV Scenario 1. Eight invalid fixtures plus one valid control.
func TestReq84_EightValidationArmsEachAssertingItsOwnNamedCode(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})

	// The valid fixture: one reader, one gate, one writer, all well-formed.
	validRegistry := func() (accessor.Registry, []accessor.Identity) {
		rd := readerDef(&readBinding{store: s}, keyStatus)
		gd := gateDef(&gateBinding{verdict: accessor.VerdictAllow})
		wd := writerDef(&writeBinding{store: s}, keyStatus)
		return registryOf(rd, gd, wd),
			[]accessor.Identity{rd.Identity, gd.Identity, wd.Identity}
	}

	arms := map[string]struct {
		build func() (accessor.Registry, []accessor.Identity)
		want  accessor.ValidationCode
	}{
		"missing_accessor": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				ids = append(ids, accessor.Identity{
					Flow: flowID, Name: "absent.read", Capability: accessor.CapRead,
				})
				return reg, ids
			},
			want: accessor.CodeMissingAccessor,
		},
		"multiply_bound_accessor": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				dup := readerDef(&readBinding{store: s}, keyStatus)
				reg.Definitions = append(reg.Definitions, dup)
				return reg, ids
			},
			want: accessor.CodeMultiplyBoundAccessor,
		},
		"capability_mismatch": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				// The model references `state.read` as a WRITE, but only a
				// read binding carries that id.
				ids = append(ids, accessor.Identity{
					Flow: flowID, Name: readerName, Capability: accessor.CapWrite,
				})
				return reg, ids
			},
			want: accessor.CodeCapabilityMismatch,
		},
		"missing_or_non_positive_timeout": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				reg.Definitions[0].Accessor.Timeout = "0s"
				return reg, ids
			},
			want: accessor.CodeMissingOrNonPositiveTimeout,
		},
		"missing_write_read_back": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				for i := range reg.Definitions {
					if reg.Definitions[i].Identity.Capability == accessor.CapWrite {
						reg.Definitions[i].Accessor.ReadBack = false
					}
				}
				return reg, ids
			},
			want: accessor.CodeMissingWriteReadBack,
		},
		"ambient_artifact_discovery": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				reg.Definitions[0].AmbientDiscovery = true
				return reg, ids
			},
			want: accessor.CodeAmbientArtifactDiscovery,
		},
		"write_non_owned_tag": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				for i := range reg.Definitions {
					if reg.Definitions[i].Identity.Capability == accessor.CapWrite {
						// `profile` is observed, not owned.
						reg.Definitions[i].Accessor.Keys = []string{keyProfile}
					}
				}
				return reg, ids
			},
			want: accessor.CodeWriteNonOwnedTag,
		},
		// The EIGHTH arm — A9, unwitnessed by the spike.
		"missing_requested_key_set": {
			build: func() (accessor.Registry, []accessor.Identity) {
				reg, ids := validRegistry()
				reg.Definitions[0].Accessor.Keys = nil
				return reg, ids
			},
			want: accessor.CodeMissingRequestedKeySet,
		},
	}

	if len(arms) != 8 {
		t.Fatalf("the fixture carries %d arms; the record names EIGHT", len(arms))
	}

	for name, arm := range arms {
		t.Run(name, func(t *testing.T) {
			reg, ids := arm.build()

			fs := accessor.Validate(reg, ids)

			if len(fs) == 0 {
				t.Fatalf("validation returned no findings; want the named code %q",
					arm.want)
			}
			if !hasCode(fs, arm.want) {
				t.Errorf("validation codes = %v; want the arm's OWN named code %q "+
					"— asserting merely that validation returned non-empty does "+
					"not discriminate a dropped arm", findingCodes(fs), arm.want)
			}
		})
	}

	t.Run("negative_control_the_valid_fixture_returns_the_empty_set", func(t *testing.T) {
		reg, ids := validRegistry()

		fs := accessor.Validate(reg, ids)

		if len(fs) != 0 {
			t.Errorf("validation codes = %v for a well-formed registry; want the "+
				"EMPTY set — a validator that rejects everything passes every "+
				"arm above and fails here", findingCodes(fs))
		}
	})
}

// REQ-84, code-set half: the eight arms' codes are exactly the seven
// witnessed names plus the eighth.
// BOUNDARY
func TestReq84_TheValidationCodeSetIsClosedAtEight(t *testing.T) {
	want := []accessor.ValidationCode{
		"missing_accessor",
		"multiply_bound_accessor",
		"capability_mismatch",
		"missing_or_non_positive_timeout",
		"missing_write_read_back",
		"ambient_artifact_discovery",
		"write_non_owned_tag",
		accessor.CodeMissingRequestedKeySet,
	}

	got := slices.Clone(accessor.ValidationCodes())
	slices.Sort(got)
	sortedWant := slices.Clone(want)
	slices.Sort(sortedWant)

	if !slices.Equal(got, sortedWant) {
		t.Errorf("ValidationCodes() = %v; want exactly %v — seven witnessed "+
			"spellings plus the eighth (missing/empty requested key set)",
			got, sortedWant)
	}
}

// REQ-91: "Execution has three phases" — validation before resolution,
// runtime invocation of read and gate accessors with classification, then
// write accessors from a successful plan.
// REQ-90: the illustrative execution order is normative "only in its
// ordering of write-then-read-back".
// BOUNDARY
//
// Ordering, asserted as an observable sequence: the write's mutation
// precedes the read-back's read.
func TestReq90_WriteThenReadBackIsTheNormativeOrdering(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	seq := &sequence{}
	w := &orderedWriteBinding{store: s, seq: seq}
	r := &orderedReadBinding{store: s, seq: seq}
	e := accessor.NewExecutor(
		registryOf(readerDef(r, keyStatus), writerDef(w, keyStatus)), artifactsOf())

	got := e.Write(ctxOf(t), writerName,
		planWriting(resolve.Tag{Key: keyStatus, Value: "Final"}))

	mustWriteSucceed(t, got)
	if !slices.Equal(seq.steps, []string{"write", "read"}) {
		t.Errorf("execution order = %v; want [write read] — the re-read happens "+
			"AFTER command-level success, which is what makes read-back a "+
			"verification rather than a precondition", seq.steps)
	}
}

// REQ-86: "Two runs over the same model and fixture results produce the
// same disposition" — "disposition equality", "not byte-identical output;
// ordering is normalized by sorted map formatting".
// REQ-88: "the two identical runs are compared for equality
// (`replay-identical=true`) rather than for absence of error", with the
// injected-refusal run as control: "a replay that returns success
// unconditionally fails it".
// REQ-87 / TS 4.
// HAPPY PATH
func TestReq86_ReplayProducesTheSameDisposition(t *testing.T) {
	run := func() accessor.Disposition {
		s := newStore(map[string]string{keyStatus: "Draft", keyProfile: "large"})
		def := readerDef(&readBinding{store: s}, keyStatus, keyProfile)
		e := accessor.NewExecutor(registryOf(def), artifactsOf())
		return accessor.ReadDisposition(def, e.Read(ctxOf(t), readerName))
	}

	first, second := run(), run()

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("replayed dispositions are not value-identical:\nfirst  = %+v\n"+
			"second = %+v — the comparison is EQUALITY, not absence of error",
			first, second)
	}
	if first.Refusal != "" {
		t.Fatalf("the replayed run refused with %q; the fixture is a success case",
			first.Refusal)
	}

	t.Run("control_an_injected_refusal_is_also_stable_and_is_not_success", func(t *testing.T) {
		injected := func() accessor.Disposition {
			s := newStore(map[string]string{keyStatus: "Draft"})
			def := gateDef(&gateBinding{verdict: accessor.VerdictIndeterminate})
			e := accessor.NewExecutor(registryOf(def), artifactsOf())
			_ = s
			return accessor.GateDisposition(def, e.Gate(ctxOf(t), gateName))
		}

		a, b := injected(), injected()

		if a.Refusal != accessor.ClassGateIndeterminate {
			t.Fatalf("injected disposition refusal = %q; want %q — a replay that "+
				"returns success unconditionally fails here",
				a.Refusal, accessor.ClassGateIndeterminate)
		}
		if !reflect.DeepEqual(a, b) {
			t.Errorf("the injected refusal is not stable across replay:\n%+v\n%+v", a, b)
		}
	})
}

// REQ-89: "Accessor execution can be deterministic enough for resolver
// replay when the model records artifact role, accessor name, capability,
// timeout, and returned tag values." — the recorded disposition names
// those five.
// BOUNDARY
func TestReq89_TheDispositionRecordsTheFiveReplayInputs(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := readerDef(&readBinding{store: s}, keyStatus)
	e := accessor.NewExecutor(registryOf(def), artifactsOf())

	d := accessor.ReadDisposition(def, e.Read(ctxOf(t), readerName))

	if d.Accessor != readerName {
		t.Errorf("disposition accessor = %q; want %q", d.Accessor, readerName)
	}
	if d.Capability != accessor.CapRead {
		t.Errorf("disposition capability = %q; want %q", d.Capability, accessor.CapRead)
	}
	if d.Role != stateRole {
		t.Errorf("disposition role = %q; want %q", d.Role, stateRole)
	}
	if d.Timeout != fixtureTimeout {
		t.Errorf("disposition timeout = %v; want %v", d.Timeout, fixtureTimeout)
	}
	if v, ok := seamValueOf(d.Tags, keyStatus); !ok || v != "Draft" {
		t.Errorf("disposition tags = %+v; want the returned value %q = %q",
			d.Tags, keyStatus, "Draft")
	}
}

// REQ-8: "RDR 0002 owns the TOML carrier. This RDR owns the accessor
// execution semantics embedded behind accessor references." — this RDR
// MUST NOT introduce a second carrier.
// REQ-9: "the definition is the capability-table entry … `keys` is the
// binding on readers and writers alike, and `read_back = true` is fixed
// on writers. This RDR does not spell the layout; it consumes it."
// [0002-delivered]
// BOUNDARY
//
// Asserted structurally: the definition's declaration fields ARE RDR
// 0002's `table.Accessor`, not a parallel struct this package spells.
func TestReq8_TheDefinitionConsumesRDR0002sCarrierRatherThanASecondOne(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	def := writerDef(&writeBinding{store: s}, keyStatus)

	var carried table.Accessor = def.Accessor
	if carried.Role != stateRole || carried.Path != statePath {
		t.Errorf("carried entry = %+v; want RDR 0002's role/path", carried)
	}
	if !carried.ReadBack {
		t.Error("a writer's carried entry has read_back = false; `read_back = " +
			"true` is FIXED on writers by RDR 0002's loader")
	}
	if !slices.Equal(carried.Keys, []string{keyStatus}) {
		t.Errorf("carried keys = %v; `keys` is the binding on readers and "+
			"writers alike", carried.Keys)
	}

	// `keys` is the writer's OWNED key set, the same field a reader uses
	// as its requested key set — one binding, not two spellings.
	if got := def.OwnedKeys(); !slices.Equal(got, []string{keyStatus}) {
		t.Errorf("OwnedKeys() = %v; want the carried `keys` %v",
			got, []string{keyStatus})
	}
}

// REQ-99: "Add stable refusal codes during implementation or RDR 0005." —
// `internal/cli/clierr::CLIError` codes are append-only.
// REQ-81: "One shape is guarded but **not** fully closed by this RDR … the
// residual is a property of that peer-owned scope, tracked in Capability
// Dependencies rather than claimed closed here." — this residual MUST be
// left OPEN, not closed here.
// BOUNDARY
//
// The accessor package ships classes; it does not reach into the kernel's
// or the CLI's vocabularies. REQ-77 pins the kernel side; this pins that
// the accessor package does not silently widen `Row.RequiresOwned` to
// close REQ-81's residual either.
func TestReq81_TheResidualIsLeftOpenNotClosedHere(t *testing.T) {
	s := newStore(map[string]string{keyStatus: "Draft"})
	e, _ := readerExec(t, s, keyStatus, keyProfile)

	read := e.Read(ctxOf(t), readerName)
	mustReadSucceed(t, read)

	// The executor emits a snapshot. It does NOT emit a RequiresOwned set,
	// a row set, or any other normalized-row field: closing the residual
	// would mean deriving one here, which RDR 0002's normalizer owns.
	snapshot := read.OwnedSnapshot()
	if v, ok := seamValueOf(snapshot, keyStatus); !ok || v != "Draft" {
		t.Fatalf("the snapshot does not carry %q = %q (snapshot = %+v); an "+
			"empty snapshot satisfies the omission assertion vacuously",
			keyStatus, "Draft", snapshot)
	}
	for _, tag := range snapshot {
		if tag.Key == keyProfile {
			t.Errorf("the snapshot carries the absent %q; the executor's only "+
				"move is OMISSION — it does not compensate for the residual by "+
				"synthesizing a value", keyProfile)
		}
	}
}

// --- ordering fixtures ---------------------------------------------------

type sequence struct{ steps []string }

func (s *sequence) mark(step string) { s.steps = append(s.steps, step) }

type orderedWriteBinding struct {
	store   *store
	seq     *sequence
	applies int
}

func (b *orderedWriteBinding) Capability() accessor.Capability { return accessor.CapWrite }

func (b *orderedWriteBinding) Apply(
	ctx context.Context, art accessor.Artifact, planned []resolve.Tag,
) error {
	b.applies++
	b.seq.mark("write")
	inner := &writeBinding{store: b.store}
	return inner.Apply(ctx, art, planned)
}

func (b *orderedWriteBinding) Invocations() int { return b.applies }

type orderedReadBinding struct {
	store *store
	seq   *sequence
}

func (b *orderedReadBinding) Capability() accessor.Capability { return accessor.CapRead }

func (b *orderedReadBinding) Read(
	ctx context.Context, art accessor.Artifact, requested []string,
) ([]accessor.KeyValue, []string, error) {
	b.seq.mark("read")
	inner := &readBinding{store: b.store}
	return inner.Read(ctx, art, requested)
}
