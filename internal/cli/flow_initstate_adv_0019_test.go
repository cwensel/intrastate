package cli

// RDR 0019 — Phase 3b ADVERSARIAL probes for `flow init-state`.
//
// Written from the record's `Failure Modes` (`0019:F1`–`F5`) by a reviewer
// assuming the implementation is wrong, and deliberately NOT from the
// existing 0019 suite. Each probe names the failure mode it hunts and the
// REQ that grounds its assertion; none asserts anything the record does not
// already fix.
//
// Two probes survive here. A third — "an exit-3 seed leaves an APPLIED
// store the re-run will not reseed" (F1's exit-group confusion, REQ-67) —
// was written, PASSED, and was discarded rather than kept: it is
// `TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore`
// under a different name, and a redundant copy adds no regression value.

import (
	"path/filepath"
	"testing"
)

// --- ADV-1 (F2) — the torn state the VERB ITSELF produces ----------------

// initTornSeedModel is `initTwoRoleModel`'s shape with role B's writer
// declaring an UNREACHABLE read-back locator, so a seed through it lands
// the mutation AND seals the artifact, and the invocation refuses at exit 3
// with the other role never reached.
//
// It is the only construction that drives the verb's OWN write loop to a
// partial commit. `runFlowInitState` iterates the writers in sorted NAME
// order (`aux` before `state`), so the failing writer must be the one that
// sorts FIRST for the partial commit to be observable — hence `aux` is the
// sealing one and `state` the one the loop never reaches.
const initTornSeedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "inittornseed"
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
required = true

[tags.note]
provenance = "owned"
kind = "scalar"

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

[read.aux]
role = "aux"
path = "flow.aux"
keys = ["note"]
timeout = "2s"

[write.aux]
role = "aux"
path = "flow.aux.readback-unreachable"
keys = ["note"]
timeout = "2s"
read_back = true

[initial]
stage = "seeded"
note = "hello"

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

// ADV-1 — F2 "Torn (multi-writer)": "no cross-writer atomicity — a failed
// writer can leave the artifact seeded for some keys only. The store is
// then non-empty, so a re-run is a NO-OP whose payload lists the missing
// `[initial]` keys — visible rather than silently skipped — but it does NOT
// repair them."
// REQ-47, REQ-51, REQ-60, REQ-67.
//
// The record's torn state is the one the VERB produces, and the shipped S12
// test constructs its torn state with a prior `set-state`, so no test drives
// the verb's own write loop to a partial commit. This one does: the first
// writer in the loop's order seals its artifact and refuses at exit 3, and
// the assertions are the three the record fixes — (a) the partial commit is
// real and the unreached role is untouched, (b) the re-run is a no-op
// SUCCESS that repairs nothing, (c) the re-run's payload names precisely the
// un-seeded key.
//
// STATUS: PASSES against the implementation — the verb defends F2's torn
// case. Kept for the regression value the shipped suite lacks: it is the
// only test whose torn state the verb itself produced.
func TestAdv1_0019_ATornSeedTheVerbItselfProducesIsReportedNotRepaired(t *testing.T) {
	model := writeFlowModel(t, initTornSeedModel)
	dir := t.TempDir()
	artA := filepath.Join(dir, "a.artifact")
	artB := filepath.Join(dir, "b.artifact")
	bindA := artifactBinding(initRoleA, artA)
	bindB := artifactBinding(initRoleB, artB)

	// The seeding arm runs and TEARS: `aux` sorts first, commits `note`
	// into a sealed artifact, and its read-back cannot complete — exit 3,
	// "applied but unverified", with `state` never reached.
	ce := initRefusal(t, 3, initStateArgs(model, bindA, bindB)...)
	t.Logf("the torn seed refused with %s at exit 3", ce.Code)

	// (a) The partial commit is REAL. This is the state F2 names, produced
	// by the verb rather than by a prior `set-state`. An exit-3 refusal a
	// caller could read as "nothing was written" would make the
	// applied-but-unverified disposition false.
	if _, exists := artifactBytes(t, artB); !exists {
		t.Fatalf("the failing writer's artifact %s does not exist; F2's torn "+
			"state requires the write to have LANDED", artB)
	}
	beforeA, existedA := snapshot(t, artA)
	beforeB, existedB := snapshot(t, artB)
	if existedA {
		t.Errorf("the unreached writer's artifact %s exists; the loop refused "+
			"at the FIRST writer, so nothing may have been committed for the "+
			"`state` role", artA)
	}

	// (b) The re-run over the torn store is a NO-OP SUCCESS at exit 0 that
	// repairs nothing: the store is non-empty (it carries `note` and the
	// seal), so the emptiness predicate declines.
	stdout := requireSuccess(t, initStateArgs(model, bindA, bindB)...)
	requireBytesUnchanged(t, artB, beforeB, existedB,
		"the re-run does NOT repair the torn seed")
	requireBytesUnchanged(t, artA, beforeA, existedA,
		"the re-run does NOT seed the role the torn run never reached")

	// (c) The payload names PRECISELY the un-seeded `[initial]` key — what
	// the record calls the thing that keeps the manual repair route
	// mechanical: "it names precisely the keys to pass".
	data := flowData(t, stdout)
	if field, ok := payloadCarriesKeySet(data, []string{"stage"}); !ok {
		t.Errorf("no payload field carries EXACTLY the un-seeded `[initial]` "+
			"key set {stage}; the re-run over a torn store must make the "+
			"missing keys VISIBLE (F2)\npayload fields: %v", keysOf(data))
	} else {
		t.Logf("the torn re-run's absent-key report rides field %q", field)
	}
}

// --- ADV-2 (F2) — the absent report must be answered PER ROLE ------------

// ADV-2 — F2: "The payload's absent-key list is what keeps the manual route
// mechanical — it names precisely the keys to pass."
// REQ-47: on the no-op arm "the payload reports which `[initial]` keys THE
// STORE does not carry".
//
// The adversarial reading: a report that answers presence from the UNION of
// every bound store calls a key present when the artifact its own declared
// writer serves does not carry it. The record's repair route is "explicit
// `set-state` of the listed keys", and `set-state` routes each key to ITS
// OWN writer's artifact — so the question the report must answer is
// per-role, exactly as the emptiness predicate's own quantifier is.
//
// The confound is built THROUGH THE CLI ONLY, so nothing here depends on a
// hand-staged artifact: both roles are first bound to one shared artifact
// and `note` is established there by `set-state`; role B is then rebound to
// a fresh artifact. Role A's store now carries `note` — a key only role B's
// writer serves — while role B's store carries nothing. Rebinding a role to
// a different artifact between invocations is ordinary caller behaviour: the
// binding is per-invocation and no clause forbids it.
//
// The correct absent set is {note, stage}: neither `[initial]` key is
// carried by the store its own writer serves. A union report answers
// {stage}, dropping `note`.
//
// STATUS: FAILS against the implementation.
// `internal/cli/flow_initstate.go::initAbsentKeys` folds every role's key
// set into ONE `present` map before testing membership, so a key present in
// ANY bound store suppresses its absent entry regardless of which writer
// serves it. The harm is direct and was probed separately: an operator who
// runs `set-state --write stage=seeded` — precisely and only what the
// payload lists — is left with `note` still reading ABSENT.
func TestAdv2_0019_TheAbsentReportIsAnsweredPerRoleNotFromTheUnionOfStores(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	dir := t.TempDir()
	shared := filepath.Join(dir, "shared.artifact")
	fresh := filepath.Join(dir, "b.artifact")

	// Establish `note` in the SHARED artifact through the CLI, with both
	// roles bound to it.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding(initRoleA, shared),
		"--artifact", artifactBinding(initRoleB, shared),
		"--write", "note=stray", "--as=json")

	// Rebind role B to a FRESH artifact. Role A's store carries `note`;
	// role B's — the one `note`'s writer serves — carries nothing.
	bindA := artifactBinding(initRoleA, shared)
	bindB := artifactBinding(initRoleB, fresh)

	stdout := requireSuccess(t, initStateArgs(model, bindA, bindB)...)

	data := flowData(t, stdout)
	if field, ok := payloadCarriesKeySet(data, []string{"note", "stage"}); !ok {
		t.Errorf("no payload field carries EXACTLY {note, stage}. Role A's "+
			"store carries `note`, but role A's writer does not SERVE "+
			"`note` — role B's does, and role B's store is empty, as is "+
			"role A's for `stage`. So both `[initial]` keys are absent from "+
			"the store their own writer serves. An absent report answering "+
			"from the UNION of stores calls `note` present and omits it, so "+
			"it does NOT `name precisely the keys to pass` (F2) and an "+
			"operator following it leaves `note` absent\n"+
			"payload fields: %v", keysOf(data))
	} else {
		t.Logf("the per-role absent report rides field %q", field)
	}
}
