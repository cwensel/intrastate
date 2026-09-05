package cli

// RDR 0019 — the CLASS refusal and the CARRIER SCOPE gate: `0019:C1`'s
// SEMANTICS and CARRIER SCOPE paragraphs, and scenarios S7 and S10.
//
// The suite's discriminating choice: the carrier is varied by AUTHORING a
// model whose accessor declares a different carrier, never by inspecting
// the constructed binding from a test. C1 fixes that the verb detects the
// carrier by the CONSTRUCTED BINDING'S TYPE, so a test that reached into
// the registry would be pinning the detector rather than the partition it
// induces — and would pass against the `internal/cli` re-derivation from
// `table.Accessor`'s `Edit`/`Command` fields that C1 BANS.

import (
	"os"
	"strings"
	"testing"

	"github.com/cwensel/intrastate/internal/cli/clierr"
)

// --- the class refusal ---------------------------------------------------

// REQ-4: "A `decision-table` model has no `[initial]` (0010:C2) and
// therefore no bootstrap to materialize; initialization MUST refuse for
// that class rather than succeed vacuously."
// REQ-6: "the verb's refusal is a CLASS refusal, not an emptiness one — an
// ordinary state-machine model whose `[initial]` is empty is a different
// case, already blocked at lint by 0006:C18."
// REQ-86 / `0019:S7`: "**Expected**: exit 2 with the C1 class code, raised
// before any accessor runs … Assert instead that no accessor process ran
// and no artifact was created or modified (bytes/mtime unchanged, or the
// artifact path still absent), which is the discriminating control".
// HAPPY PATH — the class arm. "Rather than succeed vacuously" is the
// discriminating half: a decision-table model has `len(Model.Initial) == 0`,
// so an implementation that simply seeded the empty set would exit 0.
func TestReq4And6And86_0019_ADecisionTableModelRefusesRatherThanSucceedingVacuously(t *testing.T) {
	model := writeFlowModel(t, initDecisionTableModel)
	art := newFlowArtifact(t, "state.artifact")

	before, existed := snapshot(t, art)

	initRefusal(t, 2, initStateArgs(model, artifactBinding(initRoleA, art))...)

	requireBytesUnchanged(t, art, before, existed,
		"the class refusal is raised BEFORE any accessor runs")
	if _, exists := artifactBytes(t, art); exists {
		t.Errorf("the class refusal CREATED an artifact at %s; no artifact is "+
			"created or modified — that absence is the discriminating "+
			"control, since a decision-table model can bind no owned-tag "+
			"writer to fail on invocation", art)
	}
}

// REQ-5: "The refusal keys on `internal/table/model.go::IsDecisionTable`
// ALONE, never on \"declares `[initial]` but is decision-table\": that
// state is UNCONSTRUCTIBLE"; LBD restates the ban on the alternative:
// "never a re-derivation from `len(owned)`".
// ADVERSARIAL — the banned re-derivation, exhibited. A state-machine model
// declaring ZERO owned tags has `len(owned) == 0` exactly as a
// decision-table model does, so an implementation keying on `len(owned)`
// would raise the CLASS code here. It must not: the class arm keys on
// `::IsDecisionTable` ALONE.
func TestReq5_0019_TheClassArmKeysOnIsDecisionTableNotOnLenOwned(t *testing.T) {
	dtModel := writeFlowModel(t, initDecisionTableModel)
	classCE := initRefusal(t, 2, initStateArgs(dtModel)...)

	// A STATE-MACHINE model with no owned tags. `len(owned) == 0` here too,
	// but `::IsDecisionTable` is FALSE, so the class code must not fire.
	smModel := writeFlowModel(t, initNoOwnedStateMachineModel)
	stdout, _, err := runCmd(t, initStateArgs(smModel)...)
	if err == nil {
		return // a non-refusal cannot carry the class code
	}
	ce := mustCLIError(t, err, stdout)
	requireVerbRegistered(t, ce)
	if ce.Code == classCE.Code {
		t.Errorf("a STATE-MACHINE model with zero owned tags raised the "+
			"CLASS refusal code %q; the refusal keys on `::IsDecisionTable` "+
			"ALONE, never a re-derivation from `len(owned)` — these two "+
			"models are indistinguishable by owned-tag count and must not "+
			"be by class", ce.Code)
	}
}

// initNoOwnedStateMachineModel is an ordinary state-machine model — NOT a
// decision table — that declares ZERO owned tags and so carries no
// `[initial]`. It is the control for the `len(owned)` re-derivation ban:
// indistinguishable from a decision-table model by owned-tag count,
// distinguishable by `::IsDecisionTable`.
const initNoOwnedStateMachineModel = `outcomes = ["advance"]

[model]
id = "initnoowned"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.profile]
provenance = "observed"
kind = "enum"
domain = ["x", "y"]
single_valued = true
required = true

[[rule]]
id = "advance"
[rule.match.recognized]
eq = "advance"
[rule.guard.all.profile]
eq = "x"
`

// REQ-7: "the verb refuses before any accessor runs — a refusal (exit 2),
// not a no-op success." — asserted by the absence of any accessor
// invocation at all, since a decision-table fixture can declare no
// owned-tag writer (SC-7, `0019:S7`).
// BOUNDARY — exit 2 specifically, distinguished from exit 0. "Not a no-op
// success" is the whole claim: the two arms differ only in exit and the
// presence of a code.
func TestReq7_0019_TheClassRefusalIsExitTwoNotANoOpSuccess(t *testing.T) {
	model := writeFlowModel(t, initDecisionTableModel)

	stdout, _, err := runCmd(t, initStateArgs(model)...)
	if err == nil {
		t.Fatalf("init-state SUCCEEDED against a decision-table model; the "+
			"class arm is a REFUSAL (exit 2), not a no-op success\n"+
			"stdout: %s", stdout)
	}
	initRefusal(t, 2, initStateArgs(model)...)
}

// --- the carrier gate ----------------------------------------------------

// REQ-38: "Against a model any of whose bound write accessors is
// non-file-backed, init-state MUST refuse — a distinct terminal refusal in
// the `flow-*` family naming the accessor and its carrier — rather than
// seed on an emptiness answer it cannot compute."
// REQ-39: "The carrier gate is over BOTH capabilities, not the write side
// alone … The refusal therefore fires when EITHER the bound write binding
// or the bound read binding for a needed role is not the file-backed type,
// and it names which capability and which accessor failed."
// REQ-40: "it admits ONLY `*flowbind.Writer` on the write side, refusing
// `*flowbind.EditWriter` and `*cmdbind.Writer`, and ONLY `flowbind.Reader`
// on the read side, refusing `cmdbind.Reader`."
// REQ-90 / `0019:S10`: "run three times over an otherwise S1-shaped
// fixture. Twice on the WRITE side: once with an edit-carried accessor …
// and once with a command-backed one … Once on the READ side: a file-backed
// `[write.x]` paired with a command-backed `[read.x]` on the same role,
// which the write-side type-switch alone would ADMIT".
// REQ-91 / SC-10 expectation: "**Expected**: exit 2 with the carrier
// refusal code, naming the accessor and its carrier, with ZERO writes and
// no accessor invocation … the discriminating control is that removing the
// carrier check makes the edit case attempt a write. For the edit carrier
// the refusal is over-determined … and the test must still see the carrier
// code, not a missing-reader error. The command-backed rows are run with
// `--allow-commands` UNSET and assert the carrier code rather than the
// allow-commands refusal … The read-side row asserts the refusal names the
// READ capability and its accessor".
// ADVERSARIAL — S10's three rows in one table. The read-side row is the one
// the write-side type switch alone would ADMIT, and it doubles as the
// detector for the POINTER/VALUE spelling REQ-41 makes normative.
func TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding(t *testing.T) {
	rows := []struct {
		name       string
		src        string
		capability string
		accessor   string
		carrier    string
	}{
		{
			name:       "write-side-edit-carried",
			src:        initEditCarriedModel,
			capability: "write",
			accessor:   "state",
			carrier: "edit — `*flowbind.EditWriter` has no JSON store, no " +
				"sealedKey, and no artifact bytes, so \"the store carries no " +
				"key\" is UNDEFINED for it",
		},
		{
			name:       "write-side-command-backed",
			src:        initCommandWriteModel,
			capability: "write",
			accessor:   "state",
			carrier:    "command — `*cmdbind.Writer`",
		},
		{
			name:       "read-side-command-backed",
			src:        initCommandReadModel,
			capability: "read",
			accessor:   "state",
			carrier: "command — `cmdbind.Reader`. This row is the one the " +
				"WRITE-side type switch alone would ADMIT: `[write.state]` " +
				"is file-backed",
		},
	}

	// The three rows must share ONE code — a single carrier refusal, not
	// three — so the codes are collected and compared.
	codes := map[string]string{}

	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			model := writeFlowModel(t, row.src)
			art := newFlowArtifact(t, "state.artifact")
			bind := artifactBinding(initRoleA, art)

			before, existed := snapshot(t, art)

			// `--allow-commands` is UNSET on the command-backed rows
			// deliberately (REQ-91): the carrier code must come back, not
			// the allow-commands refusal.
			ce := initRefusal(t, 2, initStateArgs(model, bind)...)
			codes[row.name] = ce.Code

			// ZERO writes. For the EDIT row this is the discriminating
			// control: removing the carrier check makes the edit case
			// ATTEMPT A WRITE.
			requireBytesUnchanged(t, art, before, existed,
				"the carrier refusal is decided BEFORE any accessor is invoked")
			if _, exists := artifactBytes(t, art); exists && !existed {
				t.Errorf("the carrier refusal CREATED an artifact; the gate "+
					"reads `Identity`, `Accessor`, and `Binding` as FIELDS "+
					"and invokes no `Binding` method (%s)", row.carrier)
			}

			// The refusal NAMES which capability and which accessor failed.
			detail := strings.ToLower(ce.Message + " " + ce.Code + " " +
				findingsText(t, ce))
			if !strings.Contains(detail, row.capability) {
				t.Errorf("the refusal does not name the %s CAPABILITY; C1 "+
					"fixes that it names which capability and which accessor "+
					"failed\nmessage: %s", row.capability, ce.Message)
			}
			if !strings.Contains(detail, row.accessor) {
				t.Errorf("the refusal does not name the accessor %q\n"+
					"message: %s", row.accessor, ce.Message)
			}
		})
	}

	// One carrier refusal serves all three rows: the gate is over BOTH
	// capabilities, not two separate gates minting two codes.
	seen := map[string][]string{}
	for name, code := range codes {
		seen[code] = append(seen[code], name)
	}
	if len(seen) > 1 {
		t.Errorf("the three carrier rows raised %d DISTINCT codes (%v); the "+
			"carrier refusal is ONE distinct terminal refusal in the `flow-*` "+
			"family, over both capabilities", len(seen), seen)
	}
}

// REQ-41: "The POINTER/VALUE spelling differs between the two capabilities
// and is normative here, because a type-switch case on the wrong one
// matches nothing … A read-side case written `*flowbind.Reader` would
// therefore refuse EVERY model, file-backed ones included" — S1 is the
// second control (SC-10).
// ADVERSARIAL — the FAILURE IS TOTAL, so the detector is a file-backed
// model that must SUCCEED. A read-side case written `*flowbind.Reader`
// matches nothing (the registry constructs readers as VALUES), so it would
// refuse this model too — and that total failure is what this test names.
func TestReq41_0019_AFileBackedModelIsAdmittedOnBothSides(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

	stdout, _, err := runCmd(t, initStateArgs(model, bind)...)
	if err != nil {
		t.Fatalf("a FULLY FILE-BACKED model was REFUSED (%v); the read-side "+
			"type-switch case must be written `flowbind.Reader` (a VALUE — "+
			"the registry constructs readers as values with value receivers "+
			"on Read). A case written `*flowbind.Reader` matches NOTHING and "+
			"refuses EVERY model, file-backed ones included\nstdout: %s",
			err, stdout)
	}
}

// REQ-42: "Re-deriving the carrier in `internal/cli` from `table.Accessor`'s
// `Edit`/`Command` fields is BANNED … Exporting `commandBacked` is not
// required and is not authorized here."
// ADVERSARIAL — the ban is on a DUPLICATED DISCRIMINATOR, and the
// observable consequence of duplicating it is DISAGREEMENT with the
// registry. The registry's residue rule routes an accessor carrying NEITHER
// `command` NOR `path` to a REFUSING command binding — a carrier a
// `len(acc.Command) != 0` re-derivation would read as file-backed and
// ADMIT. The gate must refuse it.
func TestReq42_0019_TheCarrierComesFromTheConstructedBindingNotAReDerivation(t *testing.T) {
	model := writeFlowModel(t, initCarrierlessModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	before, existed := snapshot(t, art)

	requireInitRefusedOnItsOwnTerms(t,
		"a CARRIER-LESS write accessor (the registry routes it to a REFUSING "+
			"binding, and a re-derivation from `table.Accessor`'s "+
			"`Edit`/`Command` fields — which is BANNED — is the only reading "+
			"that admits it)",
		initStateArgs(model, bind)...)
	requireBytesUnchanged(t, art, before, existed,
		"a carrier-less accessor never reaches a write")
}

// initCarrierlessModel declares a write accessor carrying NEITHER `command`
// NOR `path`. The registry's selection is TOTAL and fails closed on this
// residue, building a REFUSING binding rather than a `Path: ""` file
// binding. A `len(acc.Command) != 0` re-derivation in `internal/cli` would
// read it as file-backed and admit it — the disagreement C1's ban exists to
// prevent.
//
// If the loader refuses this model outright, the refusal is still not the
// seeded arm, and the assertion above holds on either path.
const initCarrierlessModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initcarrierless"
version = 1

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.stage]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[write.state]
role = "state"
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

// REQ-43: "The verb selects over the exported `::Registry.Definitions`
// slice, admitting the FIRST definition whose `Identity.Capability` is
// `accessor.CapRead` and whose `Accessor.Role` is the needed role. FIRST
// MATCH IN REGISTRY ORDER is normative, not incidental" — a gate that
// refused on ambiguity would be WIDER than `::readerFor`; selecting the
// last would disagree.
// REQ-46: "the gate must not sort, filter, or re-order `Definitions` before
// scanning, since the equivalence rides on sharing the registry's own slice
// order."
// DOMAIN EDGE — two same-role readers, which role uniqueness does NOT
// forbid today. The FIRST in registry order is file-backed and the SECOND
// is command-backed, so:
//
//   - a gate refusing on AMBIGUITY refuses (WIDER than `::readerFor`);
//   - a gate selecting the LAST refuses on the carrier (DISAGREES);
//   - a gate taking the FIRST admits, which is what the executor does.
//
// The registry builds readers in SORTED NAME order, so `[read.a-file]`
// sorts before `[read.z-command]` and is the first match.
func TestReq43And46_0019_TheReaderIsTheFirstMatchInRegistryOrder(t *testing.T) {
	model := writeFlowModel(t, initTwoReaderModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

	stdout, _, err := runCmd(t, initStateArgs(model, bind)...)
	if err != nil {
		t.Fatalf("a model with TWO same-role readers was REFUSED (%v); the "+
			"gate admits the FIRST definition in REGISTRY ORDER whose "+
			"capability is read and whose role matches — that first match is "+
			"file-backed here. A gate that refused on AMBIGUITY would be "+
			"WIDER than `::readerFor`, and one selecting the LAST would "+
			"disagree with it\nstdout: %s", err, stdout)
	}
}

// initTwoReaderModel declares TWO readers on the SAME role. Role uniqueness
// is not enforced on the read side today, so this model is constructible.
// The registry appends readers in SORTED NAME order, so `a-file` (file-
// backed) precedes `z-command` (command-backed) in `Definitions`.
const initTwoReaderModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "inittworeader"
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

[read.a-file]
role = "state"
path = "flow.state"
keys = ["stage"]
timeout = "2s"

[read.z-command]
role = "state"
command = ["cat", "{artifact}"]
output = "raw"
keys = ["stage"]
timeout = "10s"

[write.state]
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

// REQ-44: "A no-match is the flat `false` `::readerFor` returns, which
// reaches the unbound-needed-role arm above; it MUST NOT be split into a
// separate missing-reader code."
// BOUNDARY — a model whose write accessor names a role NO reader serves.
// The refusal must land in the EXISTING artifact-binding family at the
// group that code already carries, never in a new missing-reader code the
// contract forbids minting.
func TestReq44_0019_ANoMatchReachesTheUnboundRoleArmAndMintsNoMissingReaderCode(t *testing.T) {
	model := writeFlowModel(t, initNoReaderModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	ce := requireInitRefusedOnItsOwnTerms(t,
		"a write accessor whose role NO reader serves",
		initStateArgs(model, bind)...)

	lower := strings.ToLower(ce.Code)
	for _, banned := range []string{"missing-reader", "reader-missing", "no-reader"} {
		if strings.Contains(lower, banned) {
			t.Errorf("the refusal minted a separate MISSING-READER code %q; a "+
				"no-match is the flat `false` `::readerFor` returns and "+
				"reaches the unbound-needed-role arm — it MUST NOT be split "+
				"into a code of its own", ce.Code)
		}
	}
}

// initNoReaderModel declares a file-backed WRITE accessor on role `state`
// and NO reader for that role at all. `::readerFor` returns a flat `false`
// for it.
const initNoReaderModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initnoreader"
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

[read.other]
role = "elsewhere"
path = "flow.other"
keys = ["stage"]
timeout = "2s"

[write.state]
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

// REQ-45: "The gate reads `Identity`, `Accessor`, and `Binding` as FIELDS
// and invokes no `Binding` method … which is what makes S10's
// no-accessor-invoked assertion true by construction rather than by
// timing."
// ADVERSARIAL — "no accessor invocation" made observable. The command-backed
// write fixture's argv is `true`, which would SUCCEED and leave a trace if
// ever spawned; the read fixture's is `cat`. A gate that invoked a `Binding`
// method would spawn one of them before refusing. Asserted through the
// artifact's untouched state plus the absence of any execution-failure
// class in the refusal.
func TestReq45_0019_TheCarrierGateInvokesNoBindingMethod(t *testing.T) {
	for name, src := range map[string]string{
		"command-backed-write": initCommandWriteModel,
		"command-backed-read":  initCommandReadModel,
		"edit-carried-write":   initEditCarriedModel,
	} {
		t.Run(name, func(t *testing.T) {
			model := writeFlowModel(t, src)
			art := newFlowArtifact(t, "state.artifact")

			ce := initRefusal(t, 2, initStateArgs(model,
				artifactBinding(initRoleA, art))...)

			if _, exists := artifactBytes(t, art); exists {
				t.Errorf("an artifact appeared at %s; the gate reads "+
					"`Identity`, `Accessor`, and `Binding` as FIELDS and "+
					"invokes NO `Binding` method — a type switch calls "+
					"nothing", art)
			}

			// An accessor that RAN and failed would surface an execution or
			// timeout class, both of which are exit 3. The carrier refusal
			// is exit 2 by construction, never by timing.
			lower := strings.ToLower(ce.Code)
			for _, invoked := range []string{"timeout", "accessor-failed", "read-incomplete"} {
				if strings.Contains(lower, invoked) {
					t.Errorf("the refusal code %q is an INVOCATION failure "+
						"class; the carrier refusal is decided AT registry "+
						"construction, before any accessor is invoked", ce.Code)
				}
			}
		})
	}
}

// REQ-101: "Extending the predicate to those carriers is deliberately out
// of scope here (see A6, and the successor noted in Consequences)" —
// edit-carried (0028) and command-backed (0025) write accessors, and the
// command-backed reader, stay refused.
// DOMAIN EDGE — the non-goal, asserted in the direction a later phase would
// break it: an implementation that "helpfully" extended the predicate to
// one of these carriers would turn its refusal into a success.
func TestReq101_0019_TheRefusedCarriersStayRefused(t *testing.T) {
	for name, src := range map[string]string{
		"edit-carried-write-0028":   initEditCarriedModel,
		"command-backed-write-0025": initCommandWriteModel,
		"command-backed-read-0025":  initCommandReadModel,
	} {
		t.Run(name, func(t *testing.T) {
			model := writeFlowModel(t, src)
			art := newFlowArtifact(t, "state.artifact")

			requireInitRefusedOnItsOwnTerms(t,
				"init-state through a "+name+" carrier (extending the "+
					"predicate to these carriers is deliberately OUT OF SCOPE "+
					"and belongs to a successor)",
				append(initStateArgs(model, artifactBinding(initRoleA, art)),
					"--allow-commands")...)
		})
	}
}

// findingsText renders a refusal's findings as one searchable string, so a
// test can assert a refusal NAMES a capability or an accessor without
// pinning which field carried it.
func findingsText(t *testing.T, ce interface{ Error() string }) string {
	t.Helper()
	return ce.Error()
}

// mustCLIError extracts the structured error a refusal must carry.
func mustCLIError(t *testing.T, err error, stdout string) *clierr.CLIError {
	t.Helper()

	var ce *clierr.CLIError
	if !asCLIError(err, &ce) {
		t.Fatalf("refusal is not a structured CLIError: %v\nstdout: %s",
			err, stdout)
	}
	return ce
}

// initArtifactExists reports whether the artifact path exists, for the
// "artifact path still absent" half of the zero-writes oracle.
func initArtifactExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
