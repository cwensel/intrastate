package cli

// RDR 0019 — the EMPTY-STORE PREDICATE and the ORDERING of the gates:
// `0019:C1`'s "Seeding is ALL-OR-NOTHING over an EMPTY store" paragraph,
// its ORDERING clause, and scenarios S5, S6, S9, S11.
//
// The suite's discriminating choice: emptiness is varied by ESTABLISHING
// STORE STATE through the CLI (or, for the sealed store, through the
// production write path whose read-back locator is unreachable), never by
// hand-writing an artifact — except where a scenario names direct setup as
// the construction (S12). C1's predicate is a claim about STORE keys, so a
// fixture that hand-built an artifact would be pinning the format instead
// of the predicate.

import (
	"os"
	"strings"
	"testing"
)

// REQ-24: "Seeding is ALL-OR-NOTHING over an EMPTY store, never a per-key
// merge: init-state seeds if and only if EVERY bound artifact carries NO
// key, and then it seeds every `[initial]` key."
// HAPPY PATH — the seeding arm's positive direction: an empty store seeds
// EVERY `[initial]` key, asserted value-for-value through `read-state`
// rather than on exit 0.
func TestReq24_0019_AnEmptyStoreSeedsEveryInitialKey(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	requireSuccess(t, initStateArgs(model, bind)...)

	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsBack(t, owned, "stage", "seeded")
	requireReadsBack(t, owned, "note", "hello")
}

// REQ-24: "never a per-key merge"
// REQ-30: "no implementation may substitute \"every `[initial]` key reads
// absent\" for it, since that is the per-key variant this contract rejects"
// ADVERSARIAL — the per-key merge, exhibited. One `[initial]` key is
// PRESENT and the other ABSENT; a per-key implementation would seed the
// absent one. Under the contract's whole-store predicate the store is
// non-empty and NOTHING writes.
func TestReq24And30_0019_APartiallySeededStoreIsNeverPerKeyMerged(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// Establish exactly ONE of the two `[initial]` keys, through the CLI.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=final", "--as=json")

	before, existed := snapshot(t, art)
	requireSuccess(t, initStateArgs(model, bind)...)
	requireBytesUnchanged(t, art, before, existed,
		"init-state over a PARTIALLY seeded store")

	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	// The absent `[initial]` key stays ABSENT — a per-key merge would have
	// established it at `hello`.
	requireReadsAbsent(t, owned, "note",
		"the store was non-empty, so NOTHING writes")
	// And the present key keeps the value the store already carried, not
	// the `[initial]` one.
	requireReadsBack(t, owned, "stage", "final")
}

// REQ-25: "The quantifier is ALL, not ANY and not per-artifact, and it is
// load-bearing rather than stylistic … a per-artifact or ANY reading would
// seed into a store that still carries a key whenever a SIBLING artifact
// happened to be empty" — SC-11 (`0019:S11`) is its only test.
// REQ-26: "Under ALL, one surviving key anywhere in the bound set blocks
// the whole seed, so the cleared-key guarantee below holds per model rather
// than merely per artifact."
// REQ-92 / `0019:S11`: "Role A's artifact carries a key; role B's is empty.
// Run `init-state`. **Expected**: NO-OP SUCCESS with ZERO writes to BOTH
// artifacts — asserted on both artifacts' bytes, not exit code alone. This
// is the only test of C1's ALL quantifier".
// ADVERSARIAL — the ONLY test of the ALL quantifier, and the only shape
// that can discriminate it: a per-artifact or ANY reading seeds role B's
// empty artifact while role A's still carries a key.
func TestReq25And26And92_0019_OneSurvivingKeyAnywhereBlocksTheWholeSeed(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	artA := newFlowArtifact(t, "a.artifact")
	artB := newFlowArtifact(t, "b.artifact")
	bindA := artifactBinding(initRoleA, artA)
	bindB := artifactBinding(initRoleB, artB)

	// Role A carries a key; role B's artifact is never written.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bindA, "--artifact", bindB,
		"--write", "stage=final", "--as=json")

	beforeA, existedA := snapshot(t, artA)
	beforeB, existedB := snapshot(t, artB)

	requireSuccess(t, initStateArgs(model, bindA, bindB)...)

	requireBytesUnchanged(t, artA, beforeA, existedA,
		"role A (the artifact carrying the surviving key)")
	// The whole quantifier rides on THIS assertion: role B is EMPTY, and a
	// per-artifact or ANY reading would have seeded it.
	requireBytesUnchanged(t, artB, beforeB, existedB,
		"role B (the EMPTY sibling) — the quantifier is ALL over the whole "+
			"bound set, so one surviving key ANYWHERE blocks the whole seed")

	owned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, bindA, bindB)...))
	requireReadsAbsent(t, owned, "note",
		"role B's `[initial]` key must NOT be seeded while role A carries a key")
}

// REQ-24: "init-state seeds if and only if EVERY bound artifact carries NO
// key" — the positive direction over the SAME two-role shape, so the
// ALL-quantifier test above has a control.
// HAPPY PATH — with BOTH artifacts empty, both seed. Without this control
// the S11 test above would also pass against an implementation that never
// seeds anything.
func TestReq24_0019_BothArtifactsEmptySeedsBothRoles(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	bindA := artifactBinding(initRoleA, newFlowArtifact(t, "a.artifact"))
	bindB := artifactBinding(initRoleB, newFlowArtifact(t, "b.artifact"))

	requireSuccess(t, initStateArgs(model, bindA, bindB)...)

	owned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, bindA, bindB)...))
	requireReadsBack(t, owned, "stage", "seeded")
	requireReadsBack(t, owned, "note", "hello")
}

// REQ-27: "The count is over STORE keys, not owned keys, and the difference
// is reachable: a read-back-SEALED artifact carries
// `internal/cli/flowbind/flowbind.go::sealedKey` … so an artifact whose
// owned keys have all been cleared but whose last write declared an
// unreachable read-back locator is a ONE-key, NON-EMPTY store that
// init-state declines to seed."
// REQ-28: "The sealed store is nonetheless a NO-OP SUCCESS at exit 0 here,
// not an exit-3 refusal … What it needs is a key COUNT, which the predicate
// obtains without reading key values"
// REQ-29: "whichever carrier lands MUST answer cardinality on a sealed
// store rather than inheriting `::Reader.Read`'s unreadable short-circuit,
// or this arm degrades from no-op success to an exit-3 refusal and S9
// fails."
// REQ-88 / `0019:S9`: "**Expected**: NO-OP SUCCESS with zero writes — the
// store carries the seal and is therefore non-empty, even though it holds
// no owned key. Asserted on artifact bytes, not just exit code. This pins
// the predicate's count to STORE keys".
// REQ-89 / SC-9 setup obligation: "The test MUST fail loudly rather than
// skip if its setup does not produce exactly `{sealedKey}` — assert the
// post-setup bytes before the scenario runs, so a representation change
// surfaces as a failure naming this coupling."
// DOMAIN EDGE — the store whose only key is the SEAL. This is the arm
// ASSUMPTION-1 turns on: a carrier that inherited `::Reader.Read`'s
// unreadable short-circuit reports exit 3 here instead of exit 0, so this
// test is the one that fails if A6's carrier is read wrong.
func TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess(t *testing.T) {
	model := writeFlowModel(t, initSealModel)
	art := newFlowArtifact(t, "sealed.artifact")
	bind := artifactBinding(initRoleA, art)

	// SETUP: a write whose read-back locator is UNREACHABLE seals the
	// artifact. The write itself refuses at exit 3 — that is the shipped
	// behaviour and is not what this scenario asserts.
	if _, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=seeded", "--as=json"); err == nil {
		t.Fatalf("the sealing write SUCCEEDED; the fixture's read-back " +
			"locator is declared unreachable, so it must refuse at exit 3 " +
			"and leave the artifact SEALED")
	}

	// SETUP OBLIGATION (REQ-89): assert the post-setup bytes BEFORE the
	// scenario runs, and FAIL LOUDLY rather than skip. A representation
	// change must surface here, naming this coupling.
	sealed := mustReadArtifact(t, art)
	if !strings.Contains(string(sealed), "flow.readback-unreachable") {
		t.Fatalf("setup did not seal the artifact — its bytes carry no "+
			"read-back-unreachable marker, so this scenario cannot pin the "+
			"predicate's count to STORE keys. This test fails LOUDLY rather "+
			"than skipping so a representation change names this coupling.\n"+
			"artifact: %q", string(sealed))
	}
	// The store must hold the seal AND NO OWNED KEY: the whole point is a
	// ONE-key store that carries no owned key.
	if strings.Contains(string(sealed), `"stage"`) {
		t.Fatalf("setup left the owned key `stage` in the store; S9 needs an "+
			"artifact whose owned keys are all gone and whose ONLY key is the "+
			"seal\nartifact: %q", string(sealed))
	}

	// THE SCENARIO: exit 0 no-op success, zero writes, asserted on bytes.
	before, existed := snapshot(t, art)
	requireSuccess(t, initStateArgs(model, bind)...)
	requireBytesUnchanged(t, art, before, existed,
		"init-state over a SEALED store — the seal makes the store non-empty "+
			"at cardinality 1, so the verb declines before any write is planned")
}

// REQ-29: the carrier "MUST answer cardinality on a sealed store rather
// than inheriting `::Reader.Read`'s unreadable short-circuit, or this arm
// degrades from no-op success to an exit-3 refusal and S9 fails."
// ADVERSARIAL — the degradation named by name. Asserted as an EXIT
// assertion separate from the bytes assertion above, so a failure reads as
// "the emptiness read inherited the unreadable short-circuit" rather than
// as a generic byte diff.
func TestReq29_0019_TheEmptinessReadDoesNotInheritTheUnreadableShortCircuit(t *testing.T) {
	model := writeFlowModel(t, initSealModel)
	art := newFlowArtifact(t, "sealed.artifact")
	bind := artifactBinding(initRoleA, art)

	if _, _, err := runCmd(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "stage=seeded", "--as=json"); err == nil {
		t.Fatalf("the sealing write SUCCEEDED; the fixture must seal")
	}

	stdout, _, err := runCmd(t, initStateArgs(model, bind)...)
	if err != nil {
		t.Errorf("init-state over a sealed store REFUSED (%v); the emptiness "+
			"read needs a key COUNT, which it obtains WITHOUT reading key "+
			"values. A carrier that inherited `::Reader.Read`'s unreadable "+
			"short-circuit degrades this arm from no-op success to an exit-3 "+
			"refusal\nstdout: %s", err, stdout)
	}
}

// REQ-31: "\"Bound\" is not a filter that can shrink the quantifier's
// domain … An unbound needed role is therefore a REFUSAL, never an artifact
// treated as empty"
// REQ-32: "the ROLE-BINDING half of plan validation runs BEFORE the
// emptiness read, so an invocation that leaves a needed role unbound
// refuses at exit 2 REGARDLESS of what the bound stores contain — it does
// not reach the predicate and cannot take the no-op arm."
// REQ-62: "The role-binding half IS reachable, since roles are bound per
// invocation … and S6 is its test."
// REQ-85 / `0019:S6`: "A required artifact role is unbound. **Expected**:
// refusal in the existing artifact-binding family … zero writes committed,
// asserted on artifact bytes/mtime unchanged, or the artifact path still
// absent."
// BOUNDARY — the ORDER is the claim, so the table varies what the BOUND
// store contains. A half-bound invocation must refuse IDENTICALLY whether
// the bound store is empty or non-empty; an implementation that read the
// predicate first would take the no-op arm on the non-empty row and refuse
// only on the empty one — "the same defect reported under two outcomes".
func TestReq31And32And62And85_0019_AnUnboundNeededRoleRefusesBeforeTheEmptinessRead(t *testing.T) {
	cases := []struct {
		name      string
		seedRoleA bool
	}{
		{name: "bound-store-EMPTY", seedRoleA: false},
		{name: "bound-store-NON-EMPTY", seedRoleA: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := writeFlowModel(t, initTwoRoleModel)
			artA := newFlowArtifact(t, "a.artifact")
			artB := newFlowArtifact(t, "b.artifact")
			bindA := artifactBinding(initRoleA, artA)

			if tc.seedRoleA {
				requireSuccess(t, "flow", "set-state", "--model", model,
					"--artifact", bindA,
					"--artifact", artifactBinding(initRoleB, artB),
					"--write", "stage=final", "--as=json")
			}

			beforeA, existedA := snapshot(t, artA)
			beforeB, existedB := snapshot(t, artB)

			// Role B is NEEDED — `[write.aux]` serves the `[initial]` key
			// `note` — and is left UNBOUND. The refusal is in the EXISTING
			// artifact-binding family, at the group that code already
			// carries (REQ-66): `flow-artifact-missing`, exit 2.
			requireRefusal(t, codeArtifactMissing, 2, initStateArgs(model, bindA)...)

			requireBytesUnchanged(t, artA, beforeA, existedA,
				"an unbound needed role refuses with ZERO writes committed")
			requireBytesUnchanged(t, artB, beforeB, existedB,
				"the unbound role's artifact is never touched")
		})
	}
}

// REQ-33: "It runs AFTER the carrier gate, which is likewise observable: an
// invocation against a non-file-backed model that ALSO omits a needed
// artifact flag refuses under the carrier code, not the artifact-binding
// family."
// BOUNDARY — the ordering's OTHER edge, and the one construction that
// separates the two gates: BOTH defects present at once. The refusal must
// be the CARRIER one, "because the carrier fact is a property of the model
// alone while role-binding is a property of the invocation, and the refusal
// that names the model's own defect is the more useful one".
func TestReq33_0019_TheCarrierGatePreemptsTheRoleBindingRefusal(t *testing.T) {
	model := writeFlowModel(t, initCommandWriteModel)

	// No `--artifact` at all: the needed role is unbound AND the write
	// carrier is command-backed. The carrier code must win.
	ce := initRefusal(t, 2, initStateArgs(model)...)

	if ce.Code == codeArtifactMissing {
		t.Errorf("the refusal is the ARTIFACT-BINDING family (%s); the "+
			"carrier gate runs BEFORE the role-binding half, so a "+
			"non-file-backed model that ALSO omits a needed artifact flag "+
			"refuses under the CARRIER code", ce.Code)
	}
}

// REQ-34: "The WRITER-ARITY half is not ordered here at all, having already
// been settled at load … Only encoding and the write follow the predicate."
// REQ-61: "the writer-arity half is already enforced at LOAD … Such a model
// refuses as `flow-model-invalid` at load, never in the writer-routing
// family, and S5 asserts that ordering. The verb states the requirement (it
// MUST NOT route a key to two writers) but owes no runtime check of its
// own for it."
// REQ-84 / `0019:S5`: "**Expected**: refusal at LOAD (`flow-model-invalid`),
// BEFORE the verb's own plan validation runs — NOT the writer-routing
// family … and asserts ZERO writes committed: artifact bytes/mtime
// unchanged, or the artifact path still absent."
// ADVERSARIAL — a model routing one `[initial]` key to TWO writers. The
// discriminating assertion is the CODE: `flow-model-invalid` (the load
// refusal), never `flow-write-unbound` or another writer-routing code,
// which is what a verb that grew its own arity check would raise.
func TestReq34And61And84_0019_TwoWritersForOneInitialKeyRefusesAtLoad(t *testing.T) {
	model := writeFlowModel(t, initTwoWriterModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	before, existed := snapshot(t, art)

	requireRefusal(t, codeModelInvalid, 2, initStateArgs(model, bind)...)

	requireBytesUnchanged(t, art, before, existed,
		"a model refused at LOAD commits ZERO writes")
}

// initTwoWriterModel routes the ONE `[initial]` key `stage` to TWO declared
// write accessors. `::checkAccessorBindings` folds every `[initial]` key
// into its `written` set and refuses `writerCount[key] != 1` at LOAD, so
// this model never reaches the verb's own plan validation.
const initTwoWriterModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "inittwowriter"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "enum"
domain = ["seeded", "final"]
single_valued = true

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[write.shadow]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"

[context.done]
[context.done.match.stage]
eq = "final"

[[rule]]
id = "advance"
[rule.match.stage]
eq = "seeded"
[rule.match.recognized]
eq = "advance"
[rule.write]
stage = "final"
`

// REQ-35 / `0019:C1` ORDERING: "the CLASS refusal precedes the CARRIER
// refusal." — "the class arm keys on `::IsDecisionTable`, which is
// answerable from the loaded model alone, while the carrier arm needs the
// constructed registry."
// BOUNDARY — the order is fixed for a model that would trip BOTH arms. No
// admitted decision-table model carries `[initial]` or an owned-tag writer,
// so the combination is not constructible through the loader; what IS
// assertable is the cheaper half of the claim: a decision-table model
// refuses WITHOUT a registry being built, which is what "class-first is
// therefore the cheaper order" means operationally.
func TestReq35_0019_TheClassRefusalPrecedesTheCarrierRefusal(t *testing.T) {
	model := writeFlowModel(t, initDecisionTableModel)

	classCE := initRefusal(t, 2, initStateArgs(model)...)

	// The class refusal fires with NO `--artifact` binding and no accessors
	// declared at all, so no registry construction can have gated it. A
	// carrier-first implementation would have to build a registry to
	// decide, and on a model declaring no accessors would reach some other
	// arm entirely.
	carrierModel := writeFlowModel(t, initCommandWriteModel)
	carrierCE := initRefusal(t, 2, initStateArgs(carrierModel)...)

	if classCE.Code == carrierCE.Code {
		t.Errorf("the class refusal and the carrier refusal carry the SAME "+
			"code %q; they are DISTINCT codes, and the ORDERING clause "+
			"depends on telling them apart", classCE.Code)
	}
}

// REQ-36: "The carrier refusal is decided AFTER registry construction (it
// reads the constructed binding) but BEFORE any accessor is invoked, and it
// PREEMPTS the `--allow-commands` refusal
// (`internal/accessor/model.go::Registry.AllowCommands`, also exit 2). A
// command-backed accessor therefore yields the carrier code whether or not
// the opt-in was passed"
// ADVERSARIAL — the preemption is "a real ordering, not a tautology". The
// table runs the command-backed model BOTH with `--allow-commands` UNSET
// and SET: the SAME code must come back, since otherwise "the same model
// would refuse under two different codes depending on a flag irrelevant to
// the carrier fact".
func TestReq36_0019_TheCarrierRefusalPreemptsTheAllowCommandsRefusal(t *testing.T) {
	model := writeFlowModel(t, initCommandWriteModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	unset := initRefusal(t, 2, initStateArgs(model, bind)...)

	setArgs := append(initStateArgs(model, bind), "--allow-commands")
	set := initRefusal(t, 2, setArgs...)

	if unset.Code != set.Code {
		t.Errorf("the refusal code depends on --allow-commands (unset=%q, "+
			"set=%q); the carrier refusal PREEMPTS the allow-commands "+
			"refusal, so a command-backed accessor yields the CARRIER code "+
			"whether or not the opt-in was passed", unset.Code, set.Code)
	}

	// And the code must not be the allow-commands one: a carrier check
	// moved after registry invocation would reach `cmdbind::spawn`'s
	// pre-spawn rung first and fail this row on the wrong code.
	if strings.Contains(unset.Code, "allow-commands") ||
		strings.Contains(strings.ToLower(unset.Message), "--allow-commands") {
		t.Errorf("the refusal is the ALLOW-COMMANDS one (%s: %s); the "+
			"carrier refusal is decided AT registry construction, and both "+
			"allow-commands refusal sites are PAST construction",
			unset.Code, unset.Message)
	}
}

// REQ-37 / IP Phase 1 Step 2: "Gate the carrier on both capabilities (Step
// 1's class refusal has already fired); non-file-backed → refuse (exit 2).
// Then validate that every role a needed writer names is bound; unbound →
// refuse (exit 2) without reading any store. Then read the bound stores;
// non-empty → the no-op arm (report absent `[initial]` keys, write
// nothing). Empty → plan and encode every `[initial]` assignment
// (`writerFor` per key) before any accessor runs."
// BOUNDARY — the four-step pipeline walked end to end, one row per step's
// OUTCOME, so a reordering surfaces as the wrong row failing.
func TestReq37_0019_TheGatePipelineRunsClassCarrierBindingThenPredicate(t *testing.T) {
	t.Run("step1-class-refuses", func(t *testing.T) {
		model := writeFlowModel(t, initDecisionTableModel)
		initRefusal(t, 2, initStateArgs(model)...)
	})

	t.Run("step2-carrier-refuses-exit-2", func(t *testing.T) {
		model := writeFlowModel(t, initCommandWriteModel)
		art := newFlowArtifact(t, "state.artifact")
		initRefusal(t, 2, initStateArgs(model, artifactBinding(initRoleA, art))...)
	})

	t.Run("step3-unbound-role-refuses-without-reading-any-store", func(t *testing.T) {
		model := writeFlowModel(t, initTwoRoleModel)
		artA := newFlowArtifact(t, "a.artifact")
		bindA := artifactBinding(initRoleA, artA)

		before, existed := snapshot(t, artA)
		requireRefusal(t, codeArtifactMissing, 2, initStateArgs(model, bindA)...)
		requireBytesUnchanged(t, artA, before, existed,
			"the role-binding refusal reads NO store")
	})

	t.Run("step4a-non-empty-takes-the-no-op-arm", func(t *testing.T) {
		model := writeFlowModel(t, initMVVModel)
		art := newFlowArtifact(t, "state.artifact")
		bind := artifactBinding(initRoleA, art)

		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", bind, "--write", "stage=final", "--as=json")

		before, existed := snapshot(t, art)
		requireSuccess(t, initStateArgs(model, bind)...)
		requireBytesUnchanged(t, art, before, existed,
			"the no-op arm writes NOTHING")
	})

	t.Run("step4b-empty-plans-and-encodes-every-initial-assignment", func(t *testing.T) {
		model := writeFlowModel(t, initMVVModel)
		bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

		requireSuccess(t, initStateArgs(model, bind)...)
		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsBack(t, owned, "stage", "seeded")
		requireReadsBack(t, owned, "note", "hello")
	})
}

// REQ-49: "clearing the LAST remaining key empties the store, and an
// emptied store is indistinguishable at the content level from a
// never-written one … A subsequent init-state therefore RESEEDS such a
// store. This is the accepted residual of rejecting tombstones … no payload
// or doc may describe init-state as unable to re-establish a cleared key
// without this qualification."
// BOUNDARY — "asserted in the failing direction" (MVV step 9). This is the
// one place the cleared-key guarantee does NOT hold, and it must be pinned
// so a tombstone added later — which would preserve the guarantee here —
// fails this test rather than passing silently.
func TestReq49_0019_ClearingTheLastKeyReturnsTheStoreToTheInitializableClass(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	requireSuccess(t, initStateArgs(model, bind)...)

	// Clear BOTH keys, so the store carries none.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--clear", "note", "--clear", "stage", "--as=json")

	emptied := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsAbsent(t, emptied, "stage", "after clearing the last key")
	requireReadsAbsent(t, emptied, "note", "after clearing the last key")

	// A subsequent init-state RESEEDS: a clear is a key REMOVAL, not a
	// tombstone, so an emptied store and a never-written one are the same
	// store.
	requireSuccess(t, initStateArgs(model, bind)...)

	reseeded := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsBack(t, reseeded, "stage", "seeded")
	requireReadsBack(t, reseeded, "note", "hello")
}

// REQ-49: an emptied store is "indistinguishable at the content level from
// a never-written one".
// DOMAIN EDGE — the claim is about the STORE's content, and the reseed
// above is its consequence. Asserted directly: the artifact emptied by
// clears and the artifact never written seed to the SAME bytes.
func TestReq49_0019_AnEmptiedStoreAndANeverWrittenOneSeedIdentically(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)

	// Route 1: seed, clear every key, re-seed.
	emptied := newFlowArtifact(t, "emptied.artifact")
	emptiedBind := artifactBinding(initRoleA, emptied)
	requireSuccess(t, initStateArgs(model, emptiedBind)...)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", emptiedBind, "--clear", "note", "--clear", "stage", "--as=json")
	requireSuccess(t, initStateArgs(model, emptiedBind)...)

	// Route 2: seed a fresh artifact once.
	fresh := newFlowArtifact(t, "fresh.artifact")
	requireSuccess(t, initStateArgs(model, artifactBinding(initRoleA, fresh))...)

	a := string(mustReadArtifact(t, emptied))
	b := string(mustReadArtifact(t, fresh))
	if a != b {
		t.Errorf("a store emptied by CLEARS seeds differently from a "+
			"NEVER-WRITTEN one\nemptied-then-reseeded: %s\nfresh:                %s\n"+
			"a clear is a key REMOVAL, not a tombstone, so the two stores are "+
			"the same store", a, b)
	}
}

// initFileMTime is a helper for the mtime half of "zero writes committed:
// artifact bytes/mtime unchanged" (S5, S6, S7).
func initFileMTime(t *testing.T, path string) (int64, bool) {
	t.Helper()

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return 0, false
	}
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return fi.ModTime().UnixNano(), true
}
