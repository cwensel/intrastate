package cli

// RDR 0019 — the NO-OP ARM, the cleared-key guarantee, the PAYLOAD's fixed
// content, and the REFUSAL CODES: `0019:C1`'s SEMANTICS head, its no-op and
// payload paragraphs, its refusal-code table, `0019:RT4`, failure modes
// F1–F4, and scenarios S3, S8, S12, S13.
//
// Two rules bind this file specifically:
//
//  1. NEITHER new refusal code's SPELLING is pinned. REQ-65 fixes the exit
//     GROUP (2) and DISTINCTNESS only. Every assertion here compares an
//     observed spelling against ANOTHER observed spelling, or checks the
//     exit, and never against a literal.
//  2. NO PAYLOAD FIELD NAME is pinned beyond what C1 fixes. The payload
//     oracles locate content by SHAPE, so a rename is not a failure but a
//     missing or wrongly-scoped key set is.

import (
	"strings"
	"testing"
)

// --- the no-synthesis prohibition ---------------------------------------

// REQ-2: "Materializing it is an EXPLICIT, PERSISTING act: no read verb, no
// artifact load, and no accessor read path may synthesize, default, or fall
// back to `[initial]` values when a key is absent"
// REQ-3: "a cleared key reads back absent for every reader (REQ-107
// unchanged), and owned state is assembled only from caller-bound
// artifacts (0004:C3 unchanged)."
// REQ-94 / `0019:S13`: "bind an artifact in which that key is ABSENT —
// never seeded, as distinct from cleared — and read it back on every path
// that loads state: `read-state`, and `next`'s candidate computation.
// **Expected**: the key reports ABSENT on every path, and `next` reports it
// under `unknown[].reason: absent`. No read path returns the `[initial]`
// value." — the only scenario pinning REQ-2. "The discriminating control is
// that the same fixture with the key PRESENT reads back its stored value".
// ADVERSARIAL — the prohibition is an INVARIANT, not a snapshot, so it is
// asserted on EVERY path that loads state. The `[initial]` value `seeded`
// is what a synthesizing reader would return, so the test names it.
func TestReq2And3And94_0019_NoReadPathSynthesizesTheInitialValueForAnAbsentKey(t *testing.T) {
	// REQ-2's prohibition is an INVARIANT the verb's existence makes
	// reachable: `[initial]` becomes a RUNTIME bootstrap source here, and
	// the prohibition is what keeps it from also becoming a read-time
	// fallback. Asserted against a tree with no `init-state`, the scenario
	// would only be restating shipped reader behaviour.
	if !flowSubcommand(t, initStateVerb) {
		t.Fatalf("the `flow` group registers no `%s`; `[initial]` is not yet "+
			"a runtime bootstrap source, so the no-synthesis prohibition has "+
			"nothing to discriminate", initStateVerb)
	}
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// A store carrying `note` but NOT `stage`: `stage` is ABSENT — never
	// seeded, as distinct from cleared — and `[initial]` assigns it
	// `seeded`.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--write", "note=hello", "--as=json")

	t.Run("read-state", func(t *testing.T) {
		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		if got, present := owned["stage"]; present {
			t.Errorf("read-state reports `stage` = %q for a key the artifact "+
				"does NOT carry; `[initial]` assigns it %q, and no read verb, "+
				"artifact load, or accessor read path may synthesize, default, "+
				"or fall back to an `[initial]` value", got, "seeded")
		}
	})

	t.Run("next-candidate-computation", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "next", "--model", model,
			"--artifact", bind, "--as=json")

		if !nextAbsentKeys(t, stdout)["stage"] {
			t.Errorf("`next` does not report `stage` under "+
				"`unknown[].reason: absent`; the key is absent from the "+
				"artifact and no read path may fill it from `[initial]`\n"+
				"stdout: %s", stdout)
		}
		for _, s := range payloadValueStrings(t, stdout) {
			if s == "seeded" {
				t.Errorf("`next`'s payload carries the `[initial]` value "+
					"%q for an ABSENT key — a read path synthesized it", s)
			}
		}
	})

	// THE DISCRIMINATING CONTROL: the same fixture with the key PRESENT
	// reads back its STORED value. Without this, a reader that reported
	// every key absent would pass the two subtests above.
	t.Run("control-present-key-reads-its-stored-value", func(t *testing.T) {
		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", bind, "--write", "stage=final", "--as=json")

		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsBack(t, owned, "stage", "final")
	})
}

// --- the no-op arm -------------------------------------------------------

// REQ-47: "A non-empty artifact — torn, partially seeded, post-clear, or
// fully seeded alike — is a NO-OP SUCCESS: zero writes, and the payload
// reports which `[initial]` keys the store does not carry (informational,
// so a torn or post-clear state is visible without being repaired,
// resurrected, or failed on)."
// BOUNDARY — the clause enumerates FOUR store shapes and fixes ONE outcome
// for all of them, so the table walks all four. A per-shape implementation
// would differ on at least one row.
func TestReq47_0019_EveryNonEmptyStoreShapeIsANoOpSuccess(t *testing.T) {
	shapes := []struct {
		name    string
		arrange func(t *testing.T, model, bind string)
		absent  []string
	}{
		{
			name: "fully-seeded",
			arrange: func(t *testing.T, model, bind string) {
				requireSuccess(t, initStateArgs(model, bind)...)
			},
			absent: nil,
		},
		{
			name: "partially-seeded",
			arrange: func(t *testing.T, model, bind string) {
				requireSuccess(t, "flow", "set-state", "--model", model,
					"--artifact", bind, "--write", "stage=seeded", "--as=json")
			},
			absent: []string{"note"},
		},
		{
			name: "post-clear",
			arrange: func(t *testing.T, model, bind string) {
				requireSuccess(t, initStateArgs(model, bind)...)
				requireSuccess(t, "flow", "set-state", "--model", model,
					"--artifact", bind, "--clear", "note", "--as=json")
			},
			absent: []string{"note"},
		},
		{
			name: "torn-carrying-a-non-initial-owned-key-only",
			arrange: func(t *testing.T, model, bind string) {
				requireSuccess(t, "flow", "set-state", "--model", model,
					"--artifact", bind, "--write", "note=other", "--as=json")
			},
			absent: []string{"stage"},
		},
	}

	for _, shape := range shapes {
		t.Run(shape.name, func(t *testing.T) {
			model := writeFlowModel(t, initMVVModel)
			art := newFlowArtifact(t, "state.artifact")
			bind := artifactBinding(initRoleA, art)

			shape.arrange(t, model, bind)

			before, existed := snapshot(t, art)
			stdout := requireSuccess(t, initStateArgs(model, bind)...)
			requireBytesUnchanged(t, art, before, existed,
				"a non-empty artifact is a NO-OP SUCCESS: ZERO writes")

			// The payload reports which `[initial]` keys the store does not
			// carry — informational, so the state is VISIBLE without being
			// repaired, resurrected, or failed on.
			data := flowData(t, stdout)
			for _, key := range shape.absent {
				if !payloadMentionsKey(data, key) {
					t.Errorf("the no-op payload does not name the absent "+
						"`[initial]` key %q; the arm reports which "+
						"`[initial]` keys the store does not carry\n"+
						"payload fields: %v", key, keysOf(data))
				}
			}
		})
	}
}

// REQ-48: "while the store retains AT LEAST ONE key, a cleared key is never
// re-established by init-state under any invocation pattern, including
// unconditional automated re-runs (the artifact is non-empty, so nothing
// writes)."
// REQ-74 / `0019:RT4`: "`init-state ∘ (set-state --clear k)` on a seeded
// artifact carrying at least one key besides k writes nothing, and
// `read-state` still reports k absent — REQ-107's reader guarantee stays
// untouched"
// REQ-79 / `0019:S3`: "Seeded store, one plain key cleared, then re-run.
// **Expected**: exit 0; zero writes; the cleared key still reads absent and
// is listed as an absent `[initial]` key in the payload (RT4 / A5)."
// ADVERSARIAL — "under ANY invocation pattern, including unconditional
// automated re-runs" is the claim, so the test runs init-state REPEATEDLY.
// A single re-run cannot discriminate an implementation that resurrects on
// the second or third pass.
func TestReq48And74And79_0019_AClearedKeyIsNeverResurrectedByRepeatedInitState(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	requireSuccess(t, initStateArgs(model, bind)...)
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bind, "--clear", "note", "--as=json")

	before, existed := snapshot(t, art)

	// The unconditional-automation pattern: five consecutive re-runs.
	for i := range 5 {
		stdout := requireSuccess(t, initStateArgs(model, bind)...)

		requireBytesUnchanged(t, art, before, existed,
			"init-state re-run — the artifact is non-empty, so nothing writes")

		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsAbsent(t, owned, "note",
			"after init-state re-run — REQ-107's reader guarantee stays untouched")
		// `stage` survives: the store retains AT LEAST ONE key, which is
		// what makes the guarantee hold.
		requireReadsBack(t, owned, "stage", "seeded")

		if !payloadMentionsKey(flowData(t, stdout), "note") {
			t.Errorf("re-run %d: the payload does not list `note` as an "+
				"absent `[initial]` key; the no-op arm reports the keys the "+
				"store does not carry", i+1)
		}
	}
}

// REQ-50: "A key added to `[initial]` after seeding is re-established by
// explicit `set-state`, not by init — the payload's absent-key report names
// it." (F4)
// DOMAIN EDGE — the key-added-after-seeding case. The store is seeded under
// a model whose `[initial]` names one key, then the SAME artifact is used
// with a model whose `[initial]` names two. Init must NOT establish the new
// key; the payload must NAME it; and explicit `set-state` must repair it.
func TestReq50_0019_AKeyAddedToInitialAfterSeedingIsRepairedBySetStateNotInit(t *testing.T) {
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// Seed under the ONE-key `[initial]`.
	before := writeFlowModel(t, initOneKeyInitialModel)
	requireSuccess(t, initStateArgs(before, bind)...)

	// The model's `[initial]` now names a SECOND key.
	after := writeFlowModel(t, initMVVModel)

	snap, existed := snapshot(t, art)
	stdout := requireSuccess(t, initStateArgs(after, bind)...)
	requireBytesUnchanged(t, art, snap, existed,
		"init does not re-establish a key added to `[initial]` after seeding")

	if !payloadMentionsKey(flowData(t, stdout), "note") {
		t.Errorf("the payload does not name `note`; a key added to " +
			"`[initial]` after seeding is named by the absent-key report, " +
			"which is how the operator learns the repair is theirs to make")
	}

	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(after, bind)...))
	requireReadsAbsent(t, owned, "note", "init does not establish it")

	// The repair route is EXPLICIT `set-state`.
	requireSuccess(t, "flow", "set-state", "--model", after,
		"--artifact", bind, "--write", "note=hello", "--as=json")
	repaired := ownedFromReadState(t, requireSuccess(t, readStateArgs(after, bind)...))
	requireReadsBack(t, repaired, "note", "hello")
}

// initOneKeyInitialModel is `initMVVModel` with `[initial]` naming only
// `stage`. Seeding under it and then re-running under `initMVVModel`
// exhibits the key-added-after-seeding case (F4).
const initOneKeyInitialModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initmvv"
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
keys = ["stage", "note"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage", "note"]
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
[rule.guard.all.note]
eq = "hello"
[rule.write]
stage = "final"
`

// REQ-51: "the no-op re-run reports the absent `[initial]` keys and never
// repairs them (F2) — visibility over silent merge is the posture, and the
// repair route is explicit `set-state` of the listed keys."
// REQ-93 / `0019:S12`: "The torn state is built by writing role A's
// artifact directly during SETUP, never through the verb … **Expected**:
// NO-OP SUCCESS whose payload names exactly B's un-seeded `[initial]` keys,
// with A's committed values NOT rewritten (asserted on A's artifact
// bytes/mtime). This is the only test of F2".
// DOMAIN EDGE — the TORN state across two roles: role A committed, role B
// un-seeded. It is the only test of F2, and the payload assertion is
// EXACTLY B's keys — a payload naming A's keys too would describe a repair
// scope the verb does not have.
func TestReq51And93_0019_ATornStateIsReportedNotRepaired(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	artA := newFlowArtifact(t, "a.artifact")
	artB := newFlowArtifact(t, "b.artifact")
	bindA := artifactBinding(initRoleA, artA)
	bindB := artifactBinding(initRoleB, artB)

	// SETUP: role A's artifact is written DURING SETUP, never through the
	// verb — S12 fixes that construction, because the verb's own
	// all-or-nothing plan cannot produce a torn state.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bindA, "--artifact", bindB,
		"--write", "stage=seeded", "--as=json")

	beforeA, existedA := snapshot(t, artA)
	beforeB, existedB := snapshot(t, artB)

	stdout := requireSuccess(t, initStateArgs(model, bindA, bindB)...)

	// A's committed values are NOT rewritten.
	requireBytesUnchanged(t, artA, beforeA, existedA,
		"A's committed values are NOT rewritten — visibility over silent merge")
	requireBytesUnchanged(t, artB, beforeB, existedB,
		"B is not repaired either; the whole invocation is a no-op")

	// The payload names EXACTLY B's un-seeded `[initial]` keys.
	data := flowData(t, stdout)
	if field, ok := payloadCarriesKeySet(data, []string{"note"}); !ok {
		t.Errorf("no payload field carries EXACTLY B's un-seeded `[initial]` "+
			"key set {note}; the absent-key report's scope is the "+
			"`[initial]` keys the store does not carry, and `stage` is "+
			"committed\npayload fields: %v", keysOf(data))
	} else {
		t.Logf("the absent-key report rides field %q", field)
	}

	// The repair route is EXPLICIT `set-state` of the listed keys.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", bindB, "--write", "note=hello", "--as=json")
	owned := ownedFromReadState(t, requireSuccess(t,
		readStateArgs(model, bindA, bindB)...))
	requireReadsBack(t, owned, "note", "hello")
}

// REQ-52: "init's claim covers exactly the `[initial]` key set (C1); an
// owned key `[initial]` does not assign can still surface
// `unknown[].reason: absent` in `flow next` after a successful init." (F3)
// DOMAIN EDGE — the SCOPE of init's claim, asserted where it bites: an
// owned key outside `[initial]` is STILL absent after a successful seed, so
// `flow next` still reports it. An implementation that seeded every owned
// key would make this test's `unknown` entry disappear.
func TestReq52_0019_AnOwnedKeyOutsideInitialStillReportsAbsentAfterInit(t *testing.T) {
	model := writeFlowModel(t, initUnassignedOwnedModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

	requireSuccess(t, initStateArgs(model, bind)...)

	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsBack(t, owned, "stage", "seeded")
	requireReadsAbsent(t, owned, "extra",
		"init's claim covers EXACTLY the `[initial]` key set")

	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", bind, "--as=json")
	if !nextAbsentKeys(t, stdout)["extra"] {
		t.Errorf("`flow next` does not report `extra` under "+
			"`unknown[].reason: absent` after a SUCCESSFUL init; an owned key "+
			"`[initial]` does not assign is outside init's claim and can "+
			"still surface as absent (F3)\nstdout: %s", stdout)
	}
}

// initUnassignedOwnedModel declares an owned key `extra` that `[initial]`
// does NOT assign, and a rule guard that demands it — so `flow next` has a
// candidate whose `unknown` list can carry it.
const initUnassignedOwnedModel = `outcomes = ["advance"]
terminal = ["done"]

[model]
id = "initunassigned"
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

[tags.extra]
provenance = "owned"
kind = "scalar"

[read.state]
role = "state"
path = "flow.state"
keys = ["stage", "extra"]
timeout = "2s"

[write.state]
role = "state"
path = "flow.state"
keys = ["stage", "extra"]
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
[rule.guard.all.extra]
eq = "present"
[rule.write]
stage = "final"
`

// --- the payload ---------------------------------------------------------

// REQ-56: "The payload MUST distinguish the seeded-all case from the no-op
// case, and its scope is exactly the `[initial]` key set — the verb claims
// nothing about owned keys `[initial]` does not assign."
// HAPPY PATH — the DISTINCTION is the claim, so both arms are driven and
// their payloads compared. Two arms that rendered identically would leave a
// caller unable to branch on which one ran.
func TestReq56_0019_ThePayloadDistinguishesTheSeededArmFromTheNoOpArm(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	seededOut := requireSuccess(t, initStateArgs(model, bind)...)
	noopOut := requireSuccess(t, initStateArgs(model, bind)...)

	if strings.TrimSpace(seededOut) == strings.TrimSpace(noopOut) {
		t.Fatalf("the seeded-all payload and the no-op payload are "+
			"IDENTICAL; the payload MUST distinguish the two cases\n%s",
			strings.TrimSpace(seededOut))
	}

	// The SEEDED arm reports the seeded keys BY NAME, and its scope is
	// exactly the `[initial]` key set.
	seeded := flowData(t, seededOut)
	if field, ok := payloadCarriesKeySet(seeded, []string{"note", "stage"}); !ok {
		t.Errorf("no payload field carries exactly the `[initial]` key set "+
			"{note, stage} on the SEEDED arm; the payload reports seeded keys "+
			"by NAME\npayload fields: %v", keysOf(seeded))
	} else {
		t.Logf("the seeded-key report rides field %q", field)
	}
}

// REQ-57: "the payload reports seeded keys by NAME and does not echo their
// values — `read-state` is the surface that reports values (RT1), and
// echoing them here would create a second place a seeded value can be read
// from"
// ADVERSARIAL — the forbidden echo. The fixture's seeded values (`seeded`,
// `hello`) are searched for across EVERY string leaf of the payload, so a
// value echoed under any field name fails.
func TestReq57_0019_TheSeededPayloadNamesKeysAndNeverEchoesTheirValues(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

	stdout := requireSuccess(t, initStateArgs(model, bind)...)

	// `hello` is `note`'s seeded value and appears nowhere else in the
	// model's identifiers, so its presence can only be an echo.
	for _, s := range payloadValueStrings(t, stdout) {
		if s == "hello" {
			t.Errorf("the payload echoes the seeded VALUE %q; it reports "+
				"seeded keys by NAME only. `read-state` is the surface that "+
				"reports values (RT1), and echoing them here would create a "+
				"SECOND place a seeded value can be read from, with no "+
				"invariant tying the two\nstdout: %s", s, stdout)
		}
	}
}

// REQ-58: "The absent-key report belongs to the NO-OP arm, which is the
// only arm where the set can be non-empty; on the seeded arm every
// `[initial]` key was just written, so an absent list there is necessarily
// empty and MUST NOT be reported as if it carried information."
// BOUNDARY — the arm-scoping of the report. On the seeded arm no field may
// carry a NON-EMPTY absent list; the constraint is about carrying
// information, not about the field's presence.
func TestReq58_0019_TheSeededArmReportsNoNonEmptyAbsentKeyList(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	bind := artifactBinding(initRoleA, newFlowArtifact(t, "state.artifact"))

	seeded := flowData(t, requireSuccess(t, initStateArgs(model, bind)...))

	// On the seeded arm every `[initial]` key was just written. Any field
	// carrying a PROPER SUBSET of the `[initial]` keys would be an absent
	// list presented as if it carried information.
	for name, vals := range stringArrayFields(seeded) {
		if len(vals) == 0 || len(vals) == 2 {
			continue // empty, or the full seeded-key set
		}
		t.Errorf("payload field %q carries %v on the SEEDED arm — a partial "+
			"`[initial]` key list. Every `[initial]` key was just written, so "+
			"an absent list here is necessarily EMPTY and MUST NOT be "+
			"reported as if it carried information", name, vals)
	}
}

// REQ-59 / `0019:D-wire-byte-format`: "none new; artifacts keep
// `internal/cli/flowbind` shape, and the success payload rides the 0005:C1
// envelope. Field names deferred to implementation, owned here; the
// payload's CONTENT is not deferred"
// BOUNDARY — the envelope is the shipped one, so `flowData` (which asserts
// `type: "ok"` and a JSON-object `data`) must decode both arms. And the
// artifact keeps `flowbind` shape: a flat JSON object.
func TestReq59_0019_TheSuccessPayloadRidesTheShippedEnvelopeAndArtifactShape(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// Both arms ride the 0005:C1 envelope: `flowData` fatals otherwise.
	flowData(t, requireSuccess(t, initStateArgs(model, bind)...))
	flowData(t, requireSuccess(t, initStateArgs(model, bind)...))

	// The artifact keeps `flowbind` shape — none new. `read-state` over the
	// artifact `init-state` wrote must report the seeded values, which it
	// can only do if the shape is the one `flowbind` reads.
	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	requireReadsBack(t, owned, "stage", "seeded")
	requireReadsBack(t, owned, "note", "hello")
}

// --- plan validation and read-back refusals ------------------------------

// REQ-60: "The ENTIRE plan is validated before any write: every `[initial]`
// key must route to exactly one declared writer and every needed artifact
// role must be bound, or the verb refuses with ZERO writes committed —
// per-key commit has no atomicity across writers, so plan-level validation
// is where the all-or-nothing property lives."
// ADVERSARIAL — the all-or-nothing property is only observable across TWO
// writers, since per-key commit has no atomicity. Role A is bound and role
// B is not; if the plan were validated per-writer, A's key would commit
// before B's binding failed.
func TestReq60_0019_ThePlanIsValidatedWholeBeforeAnyWriteCommits(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	artA := newFlowArtifact(t, "a.artifact")
	bindA := artifactBinding(initRoleA, artA)

	beforeA, existedA := snapshot(t, artA)

	// Role B is NEEDED and UNBOUND. Role A is bound and empty, so a
	// per-writer implementation would commit A's `stage` seed first.
	requireRefusal(t, codeArtifactMissing, 2, initStateArgs(model, bindA)...)

	requireBytesUnchanged(t, artA, beforeA, existedA,
		"ZERO writes committed — plan-level validation is where the "+
			"all-or-nothing property lives, because per-key commit has no "+
			"atomicity across writers")
	if initArtifactExists(artA) && !existedA {
		t.Errorf("the bound role's artifact was CREATED at %s despite the "+
			"plan refusing; the ENTIRE plan is validated before ANY write", artA)
	}
}

// REQ-63: "A read-back mismatch is a distinct terminal refusal naming the
// key as PRESENT-AND-UNVERIFIED; the store is then non-empty, so a re-run
// is a no-op that does NOT repair it and MUST NOT be documented as its
// recovery — recovery is an explicit `set-state` (or discarding the
// artifact and re-running init)."
// REQ-87 / `0019:S8`: "the fixture binds the READ accessor to a different
// artifact path than its writer, pre-seeded with a conflicting value for
// the same key … **Expected**: terminal refusal naming the key
// present-and-unverified, distinct from the read-back-INCOMPLETE class an
// unreachable locator raises; a subsequent re-run is a no-op that does NOT
// repair it — asserted, since C1 forbids documenting the re-run as
// recovery."
// REQ-67: "a read-back that COMPLETED and disagreed is exit 2
// (`GroupUserEnv`), while a read-back that could not complete — the
// unreachable-locator and timeout arms — is exit 3
// (`GroupEnvUnavailable`) … So the PRESENT-AND-UNVERIFIED refusal above is
// exit 2, and a seed whose read-back is incomplete exits 3"
// DOMAIN EDGE — the mismatch arm and its non-recovery. The re-run assertion
// is what C1 explicitly requires be ASSERTED rather than merely documented.
func TestReq63And67And87_0019_AReadBackMismatchRefusesAtExitTwoAndTheReRunDoesNotRepairIt(t *testing.T) {
	model := writeFlowModel(t, initMismatchModel)
	writeArt := newFlowArtifact(t, "write.artifact")
	readArt := newFlowArtifact(t, "read.artifact")
	writeBind := artifactBinding("state", writeArt)
	readBind := artifactBinding("mirror", readArt)

	// The READ accessor's artifact is pre-seeded with a CONFLICTING value
	// for the same key, so the read-back COMPLETES and DISAGREES.
	requireSuccess(t, "flow", "set-state", "--model", model,
		"--artifact", artifactBinding("state", readArt),
		"--artifact", artifactBinding("mirror", readArt),
		"--write", "stage=conflicting", "--as=json")

	// The mismatch: exit 2, because the read-back COMPLETED.
	mismatch := initRefusal(t, 2, initStateArgs(model, writeBind, readBind)...)

	// DISTINCT from the read-back-INCOMPLETE class, which is exit 3.
	sealModel := writeFlowModel(t, initSealModel)
	sealArt := newFlowArtifact(t, "sealed.artifact")
	incomplete := initRefusal(t, 3, initStateArgs(sealModel,
		artifactBinding(initRoleA, sealArt))...)

	if mismatch.Code == incomplete.Code {
		t.Errorf("the read-back MISMATCH and the read-back INCOMPLETE class "+
			"share the code %q; a read-back that COMPLETED and disagreed is "+
			"exit 2 while one that could NOT complete is exit 3, and the two "+
			"are distinct refusals", mismatch.Code)
	}

	// The store is now non-empty, so a RE-RUN is a no-op that does NOT
	// repair it. C1 forbids documenting the re-run as recovery, so this is
	// asserted rather than assumed.
	before, existed := snapshot(t, writeArt)
	requireSuccess(t, initStateArgs(model, writeBind, readBind)...)
	requireBytesUnchanged(t, writeArt, before, existed,
		"the re-run is a NO-OP that does NOT repair a present-and-unverified "+
			"key — recovery is an explicit `set-state`")
}

// --- refusal codes and exit groups ---------------------------------------

// REQ-64: "It carries a dedicated code in the `flow-*` family, sharpened
// here as `flow-init-class-unsupported`; the carrier refusal above is
// `flow-init-carrier-unsupported`. Both follow the shipped
// `flow-<subject>-<condition>` convention"
// REQ-65: "What this clause fixes NORMATIVELY is the exit GROUP (2 on both)
// and that the two are DISTINCT codes and distinct from each other and from
// the shared classes — the literal spellings are non-normative on the
// precedent of `0028:C1.3`'s `codeWriteEditRefused`, so no test may pin the
// string."
// REQ-68: "The classes new to this verb get codes recorded in the 0005:C1
// taxonomy extension this override carries."
// BOUNDARY — per ASSUMPTION-3 this asserts EXIT GROUP and DISTINCTNESS
// ONLY. The `flow-` family prefix is the shipped convention the clause
// cites, not a spelling, so it is checked as a prefix; nothing pins the
// `<condition>` segment either code carries.
func TestReq64And65And68_0019_TheTwoNewCodesShareExitTwoAndAreDistinct(t *testing.T) {
	classModel := writeFlowModel(t, initDecisionTableModel)
	carrierModel := writeFlowModel(t, initCommandWriteModel)
	art := newFlowArtifact(t, "state.artifact")

	// EXIT GROUP: 2 on both. `initRefusal` asserts the exit.
	class := initRefusal(t, 2, initStateArgs(classModel)...)
	carrier := initRefusal(t, 2, initStateArgs(carrierModel,
		artifactBinding(initRoleA, art))...)

	// DISTINCT from each other.
	if class.Code == carrier.Code {
		t.Errorf("the class refusal and the carrier refusal share the code "+
			"%q; C1 fixes that the two are DISTINCT codes", class.Code)
	}

	// DISTINCT from the shared classes this verb reuses unchanged.
	shared := []string{
		codeModelNotFound, codeModelInvalid, codeArtifactInvalid,
		codeArtifactMissing, codeWriteUnbound, codeClearUnbound,
	}
	for _, code := range shared {
		for name, got := range map[string]string{"class": class.Code, "carrier": carrier.Code} {
			if got == code {
				t.Errorf("the %s refusal reuses the SHARED code %q; the two "+
					"classes NEW to this verb carry codes distinct from each "+
					"other AND from the shared classes", name, code)
			}
		}
	}

	// The shipped `flow-<subject>-<condition>` convention. This is the
	// FAMILY, not a spelling: nothing here pins either code's string.
	for name, got := range map[string]string{"class": class.Code, "carrier": carrier.Code} {
		if !strings.HasPrefix(got, "flow-") {
			t.Errorf("the %s refusal code %q is not in the `flow-*` family; "+
				"both follow the shipped `flow-<subject>-<condition>` "+
				"convention", name, got)
		}
	}
}

// REQ-66: "Shared refusal classes reuse the existing `flow-*` codes
// unchanged, at the groups those codes already carry — this verb introduces
// no new group for them and may not re-map one. Model selection
// (`::selectModel`/`::selectModelPath`), artifact binding and writer
// routing (`internal/cli/flow_state.go::writerFor`) are all
// `clierr.GroupUserEnv`, exit 2."
// BOUNDARY — the shared classes' codes and groups are ASSERTED EXACTLY,
// because unlike the two new codes their spellings ARE normative. A verb
// that minted its own model-selection or artifact-binding code fails here.
func TestReq66_0019_SharedRefusalClassesReuseTheirShippedCodesAndGroups(t *testing.T) {
	t.Run("model-not-found", func(t *testing.T) {
		requireRefusal(t, codeModelNotFound, 2, "flow", initStateVerb,
			"--model", "/nonexistent/model.toml", "--as=json")
	})

	t.Run("model-invalid", func(t *testing.T) {
		model := writeFlowModel(t, flowInvalidModel)
		requireRefusal(t, codeModelInvalid, 2, initStateArgs(model)...)
	})

	t.Run("artifact-missing", func(t *testing.T) {
		model := writeFlowModel(t, initMVVModel)
		requireRefusal(t, codeArtifactMissing, 2, initStateArgs(model)...)
	})

	t.Run("artifact-invalid", func(t *testing.T) {
		model := writeFlowModel(t, initMVVModel)
		args := []string{"flow", initStateVerb, "--model", model,
			"--artifact", "not-a-role-binding", "--as=json"}
		requireRefusal(t, codeArtifactInvalid, 2, args...)
	})
}

// REQ-67: "a seed whose read-back is incomplete exits 3 with the same
// partial-write disposition; both leave a non-empty store, so both make a
// re-run a no-op that does not repair." (F1)
// DOMAIN EDGE — the exit-3 arm and its DISPOSITION. The seal the failed
// read-back leaves is what makes the store non-empty, which is what makes
// the re-run a no-op.
func TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore(t *testing.T) {
	model := writeFlowModel(t, initSealModel)
	art := newFlowArtifact(t, "sealed.artifact")
	bind := artifactBinding(initRoleA, art)

	// The seed's read-back cannot COMPLETE — the locator is unreachable —
	// so it exits 3, not 2.
	initRefusal(t, 3, initStateArgs(model, bind)...)

	// The store is left NON-EMPTY (it carries the seal), so the re-run is a
	// no-op that does NOT repair.
	before, existed := snapshot(t, art)
	requireSuccess(t, initStateArgs(model, bind)...)
	requireBytesUnchanged(t, art, before, existed,
		"the incomplete read-back left a NON-EMPTY store, so the re-run is a "+
			"no-op that does not repair")
}

// REQ-102 / F2: "A repair path that preserved the cleared-key guarantee (a
// seed scoped to the failed writer's role, gated on that role's artifact
// being empty) is a successor's, not this record's."
// ADVERSARIAL — the non-goal. The tempting feature is a per-role repair
// flag; this record carries none, so any such flag must be unknown.
func TestReq102_0019_NoPerRoleRepairGrammarIsOffered(t *testing.T) {
	model := writeFlowModel(t, initTwoRoleModel)
	bindA := artifactBinding(initRoleA, newFlowArtifact(t, "a.artifact"))
	bindB := artifactBinding(initRoleB, newFlowArtifact(t, "b.artifact"))

	for _, flag := range []string{"--repair", "--role", "--force", "--seed-role"} {
		t.Run(flag, func(t *testing.T) {
			args := append(initStateArgs(model, bindA, bindB), flag, initRoleB)
			ce := requireInitRefusedOnItsOwnTerms(t,
				"`init-state "+flag+"` (a repair path scoped to one role is "+
					"a SUCCESSOR's, not this record's)", args...)
			if !strings.Contains(ce.Message, "unknown flag") {
				t.Errorf("`init-state %s` refused with %s: %s rather than as "+
					"an UNKNOWN FLAG; this record carries no per-role repair "+
					"grammar at all", flag, ce.Code, ce.Message)
			}
		})
	}
}

// REQ-103 / XC "Versioning": "per-finding code identity against it is
// cli/0017's question, not this record's."
// DOMAIN EDGE — the non-goal. This record fixes the two codes' group and
// distinctness and nothing about a per-finding identity scheme, so the
// refusals carry no version or identity field beyond the shipped finding
// shape.
func TestReq103_0019_TheNewRefusalsCarryNoPerFindingIdentityScheme(t *testing.T) {
	model := writeFlowModel(t, initDecisionTableModel)

	stdout, _, err := runCmd(t, initStateArgs(model)...)
	if err == nil {
		t.Fatalf("a decision-table model succeeded")
	}
	requireVerbRegistered(t, mustCLIError(t, err, stdout))

	for _, f := range flowFailureFindings(t, stdout) {
		if f.Class != "" && strings.Contains(f.Class, "@") {
			t.Errorf("a finding carries a versioned identity %q; per-finding "+
				"code identity is cli/0017's question, not this record's",
				f.Class)
		}
	}
}
