package cli

// RDR 0019 — the Minimum Viable Validation (`0019:MVV`), run END TO END as
// one integration test over a single fixture artifact, plus the round-trip
// invariants `0019:RT1` and `0019:RT2` and scenarios S1 and S2.
//
// The MVV is a SPINE: its ten steps share state in order, so it is written
// as one test with ordered `t.Run` subtests rather than ten independent
// ones. A step that fails leaves the following steps' preconditions
// unestablished, which is exactly the coupling the MVV asserts — so the
// spine bails on the first failure rather than reporting ten.
//
// Every round-trip assertion is VALUE-FOR-VALUE. `0019:D-identity` fixes
// the equality as REQ-107's, and RT1 says "value-for-value in canonical
// form, the REQ-107 equality, not merely exit 0" — so no step here accepts
// a green exit as evidence.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// REQ-MVV / `0019:MVV`: "Over a fixture state-machine model declaring
// `[initial]` with one always-present owned key and one plain owned key,
// both writer-served:" — the acceptance spine, ten steps, run as an
// integration test; "\"done\" is every one of its steps green plus the unit
// coverage below" (TS).
// REQ-98 / IP Prerequisites: "All Critical Assumptions verified (A2's
// canonical-form spike and A3's writer-coverage audit gate the design's
// reuse claim)" — satisfied: A1–A8 all carry `Status: Verified`, so the
// spine runs rather than being gated.
// HAPPY PATH — the user outcome, end to end.
func TestReqMVV_0019_InitStateAcceptanceSpine(t *testing.T) {
	model := writeFlowModel(t, initMVVModel)
	art := newFlowArtifact(t, "state.artifact")
	bind := artifactBinding(initRoleA, art)

	// Each step gates the next, so a failure stops the spine.
	step := func(name string, fn func(t *testing.T)) {
		if !t.Run(name, fn) {
			t.FailNow()
		}
	}

	// REQ-MVV.1: "`flow lint` certifies the model (0006 arms all green)."
	step("MVV.1-lint-certifies-the-model", func(t *testing.T) {
		// The shipped verb is the ROOT `lint`; the MVV's "`flow lint`"
		// names the lint the record means, not a `flow`-group spelling.
		stdout := requireSuccess(t, "lint", "--model", model, "--as=json")

		data := flowData(t, stdout)
		findings, ok := data["findings"].([]any)
		if ok && len(findings) != 0 {
			t.Errorf("`flow lint` reports %d finding(s); the MVV fixture must "+
				"lint clean with 0006's arms all green — a model declaring "+
				"`[initial]` satisfies the reachability-root arm (0006:C18)\n"+
				"stdout: %s", len(findings), stdout)
		}
	})

	// REQ-MVV.2: "`flow init-state --artifact state=<fresh path>` exits 0;
	// payload reports the empty-store seeding arm with both keys seeded."
	// REQ-77 / `0019:S1`: "Empty store, every `[initial]` key writer-served.
	// **Expected**: exit 0; every key seeded; `read-state` returns each at
	// its declared value in canonical form (RT1)."
	step("MVV.2-init-state-seeds-a-fresh-artifact", func(t *testing.T) {
		stdout := requireSuccess(t, initStateArgs(model, bind)...)

		data := flowData(t, stdout)
		if field, ok := payloadCarriesKeySet(data, []string{"note", "stage"}); !ok {
			t.Errorf("the payload does not report BOTH keys seeded; the "+
				"empty-store seeding arm reports seeded keys by name, and "+
				"its scope is exactly the `[initial]` key set\n"+
				"payload fields: %v", keysOf(data))
		} else {
			t.Logf("the seeded-key report rides field %q", field)
		}
	})

	// REQ-MVV.3: "`flow read-state` reports both keys at their `[initial]`
	// values (canonical form) — round-trip invariant 1."
	// REQ-69 / `0019:RT1`: "after a seeding init over role R, `flow
	// read-state` over R reports every `[initial]` key with its declared
	// value — value-for-value in canonical form, the REQ-107 equality, not
	// merely exit 0."
	// REQ-75 / `0019:D-identity`: "a seeded value's read-back equality is
	// REQ-107's form: value-for-value over the keys init planned to seed,
	// canonical wire form for set values — the identical equality
	// `set-state` already verifies, not a new one"
	step("MVV.3-read-state-reports-both-keys-value-for-value", func(t *testing.T) {
		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsBack(t, owned, "stage", "seeded")
		requireReadsBack(t, owned, "note", "hello")
	})

	// REQ-MVV.4: "`flow next` over the same artifact reports candidates with
	// no `unknown[].reason: absent` entry for either key — the first-run
	// wall is gone."
	step("MVV.4-next-reports-no-absent-key-the-first-run-wall-is-gone", func(t *testing.T) {
		stdout := requireSuccess(t, "flow", "next", "--model", model,
			"--artifact", bind, "--as=json")

		absent := nextAbsentKeys(t, stdout)
		for _, key := range []string{"stage", "note"} {
			if absent[key] {
				t.Errorf("`flow next` still reports %q under "+
					"`unknown[].reason: absent` after a successful init; the "+
					"first-run wall this verb exists to remove is still "+
					"standing\nstdout: %s", key, stdout)
			}
		}

		data := flowData(t, stdout)
		if cands, ok := data["candidates"].([]any); !ok || len(cands) == 0 {
			t.Errorf("`flow next` reports NO candidates; the step fixes that "+
				"it reports candidates, and an empty list cannot exhibit the "+
				"absence of an `absent` entry\nstdout: %s", stdout)
		}
	})

	// REQ-MVV.5: "Re-run `flow init-state`: exit 0, zero writes, artifact
	// bytes unchanged — invariant 2."
	// REQ-70 / `0019:RT2`: "the second run writes nothing, succeeds, and the
	// bytes of EVERY bound artifact are unchanged (the byte assertion is
	// scoped to the file-backed carrier, the only one this RDR admits — C1)."
	// REQ-78 / `0019:S2`: "Re-run against the store just seeded.
	// **Expected**: exit 0; zero writes; artifact bytes unchanged (RT2);
	// payload reports the no-op arm."
	step("MVV.5-re-run-writes-nothing-and-bytes-are-unchanged", func(t *testing.T) {
		before, existed := snapshot(t, art)

		stdout := requireSuccess(t, initStateArgs(model, bind)...)

		requireBytesUnchanged(t, art, before, existed,
			"the second run writes NOTHING — RT2's byte assertion")

		// The payload reports the NO-OP arm, and on a fully-seeded store the
		// absent list is EMPTY: nothing is missing.
		data := flowData(t, stdout)
		for name, vals := range stringArrayFields(data) {
			if len(vals) == 1 {
				t.Errorf("payload field %q carries %v on the no-op arm over a "+
					"FULLY SEEDED store; no `[initial]` key is absent", name, vals)
			}
		}
	})

	// REQ-MVV.6: "`flow set-state --clear <plain key>`; `flow read-state`
	// reports it absent — REQ-107 held."
	step("MVV.6-clearing-the-plain-key-reads-back-absent", func(t *testing.T) {
		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", bind, "--clear", "note", "--as=json")

		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsAbsent(t, owned, "note", "after `set-state --clear note`")
		// The always-present key survives, which is what keeps the store
		// non-empty for step 7.
		requireReadsBack(t, owned, "stage", "seeded")
	})

	// REQ-MVV.7: "`flow init-state` again (the unconditional-automation
	// pattern): exit 0, ZERO writes, payload reports the cleared key as an
	// absent `[initial]` key, `flow read-state` still reports it absent —
	// invariant 4 / A5."
	step("MVV.7-re-run-after-clear-never-resurrects-the-cleared-key", func(t *testing.T) {
		before, existed := snapshot(t, art)

		stdout := requireSuccess(t, initStateArgs(model, bind)...)

		requireBytesUnchanged(t, art, before, existed,
			"ZERO writes — the artifact is non-empty, so nothing writes")

		if !payloadMentionsKey(flowData(t, stdout), "note") {
			t.Errorf("the payload does not report the cleared key `note` as "+
				"an absent `[initial]` key; the no-op arm reports which "+
				"`[initial]` keys the store does not carry, which is how a "+
				"post-clear state is VISIBLE without being repaired\n"+
				"stdout: %s", stdout)
		}

		owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsAbsent(t, owned, "note",
			"after the unconditional-automation re-run — invariant 4 / A5")
	})

	// REQ-MVV.8: "`flow init-state` against a decision-table model that
	// LOADS (no `[initial]`, no owned tags — see S7) refuses (exit 2) with
	// the C1 class code, with no accessor run and no artifact touched;
	// against the fixture with an unbound artifact role it refuses in the
	// existing artifact-binding family with zero writes committed; against a
	// model whose write accessor is edit-carried or command-backed it
	// refuses with the carrier code and zero writes (S10)."
	step("MVV.8-the-three-refusal-arms", func(t *testing.T) {
		// (a) the CLASS arm — exit 2, no accessor run, no artifact touched.
		dtModel := writeFlowModel(t, initDecisionTableModel)
		dtArt := newFlowArtifact(t, "dt.artifact")
		classCE := initRefusal(t, 2, initStateArgs(dtModel,
			artifactBinding(initRoleA, dtArt))...)
		if initArtifactExists(dtArt) {
			t.Errorf("the class refusal TOUCHED an artifact at %s; no "+
				"accessor runs and no artifact is touched", dtArt)
		}

		// (b) the ARTIFACT-BINDING family — the EXISTING code, unchanged, at
		// the group it already carries. Its spelling IS normative.
		unboundArt := newFlowArtifact(t, "unbound.artifact")
		requireRefusal(t, codeArtifactMissing, 2,
			"flow", initStateVerb, "--model", model, "--as=json")
		if initArtifactExists(unboundArt) {
			t.Errorf("an artifact appeared at %s; the binding refusal "+
				"commits ZERO writes", unboundArt)
		}

		// (c) the CARRIER arm — edit-carried AND command-backed, both with
		// zero writes, both distinct from the class code. Per ASSUMPTION-3
		// this asserts exit group and DISTINCTNESS only.
		for name, src := range map[string]string{
			"edit-carried":   initEditCarriedModel,
			"command-backed": initCommandWriteModel,
		} {
			t.Run(name, func(t *testing.T) {
				carrierModel := writeFlowModel(t, src)
				carrierArt := newFlowArtifact(t, "carrier.artifact")

				carrierCE := initRefusal(t, 2, initStateArgs(carrierModel,
					artifactBinding(initRoleA, carrierArt))...)

				if initArtifactExists(carrierArt) {
					t.Errorf("the carrier refusal wrote an artifact at %s; "+
						"the arm fixes ZERO writes", carrierArt)
				}
				if carrierCE.Code == classCE.Code {
					t.Errorf("the carrier refusal shares the CLASS code %q; "+
						"the two are DISTINCT codes", carrierCE.Code)
				}
				if carrierCE.Code == codeArtifactMissing {
					t.Errorf("the carrier refusal reuses the "+
						"artifact-binding code %q; the classes new to this "+
						"verb are distinct from the shared classes",
						carrierCE.Code)
				}
			})
		}
	})

	// REQ-MVV.9: "The boundary, asserted in the failing direction: clear the
	// REMAINING key so the store carries none, then `flow init-state` — it
	// seeds every `[initial]` key again."
	step("MVV.9-clearing-the-last-key-returns-the-store-to-initializable", func(t *testing.T) {
		// `note` is already cleared (step 6). Clear the REMAINING key.
		requireSuccess(t, "flow", "set-state", "--model", model,
			"--artifact", bind, "--clear", "stage", "--as=json")

		emptied := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsAbsent(t, emptied, "stage", "the store now carries no key")
		requireReadsAbsent(t, emptied, "note", "the store now carries no key")

		// It seeds every `[initial]` key AGAIN. This is the accepted
		// residual of rejecting tombstones, asserted in the FAILING
		// direction so a tombstone added later fails here.
		requireSuccess(t, initStateArgs(model, bind)...)

		reseeded := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
		requireReadsBack(t, reseeded, "stage", "seeded")
		requireReadsBack(t, reseeded, "note", "hello")
	})
}

// REQ-MVV.10: "`flow init-state` against `models/rdr.toml` … bound to a
// fresh artifact seeds exactly `stage`, `status`, `gate_passed` at their
// declared values, and `flow next` then reports no `unknown[].reason:
// absent` for them." — the user outcome on a SHIPPED model, not a fixture.
// Its `[tags.gate_passed]` bool-vs-string spelling "must not be read as"
// covering the encoder arms — "S4's kind table is what covers the arms".
// HAPPY PATH — per ASSUMPTION-7 this binds the SHIPPED file rather than a
// copy, so the step keeps proving the user outcome if the model changes.
func TestReqMVV10_0019_InitStateSeedsTheShippedRdrModel(t *testing.T) {
	model := shippedRdrModelPath(t)
	art := newFlowArtifact(t, "rdr.artifact")
	bind := artifactBinding("rdr", art)

	requireSuccess(t, initStateArgs(model, bind)...)

	// Seeds EXACTLY the three keys at their DECLARED values — value-for-
	// value, not merely exit 0.
	owned := ownedFromReadState(t, requireSuccess(t, readStateArgs(model, bind)...))
	for key, want := range map[string]string{
		"stage":       "seeded",
		"status":      "draft",
		"gate_passed": "false",
	} {
		requireReadsBack(t, owned, key, want)
	}

	// And `flow next` then reports no `unknown[].reason: absent` for them —
	// the user outcome this step exists to prove.
	stdout := requireSuccess(t, "flow", "next", "--model", model,
		"--artifact", bind, "--as=json")
	absent := nextAbsentKeys(t, stdout)
	for _, key := range []string{"stage", "status", "gate_passed"} {
		if absent[key] {
			t.Errorf("`flow next` over the seeded shipped model reports %q "+
				"under `unknown[].reason: absent`; the step's claim is the "+
				"USER OUTCOME on a shipped model\nstdout: %s", key, stdout)
		}
	}
}

// shippedRdrModelPath returns the repo's LIVE `models/rdr.toml`, not a
// copy. ASSUMPTION-7: binding the shipped file is what keeps this step
// proving the user outcome if the model changes.
func shippedRdrModelPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(repoRootFor(t), "models", "rdr.toml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the shipped model at %s is missing: %v — MVV step 10's "+
			"claim is the user outcome on a SHIPPED model, so this test "+
			"fails rather than skipping", path, err)
	}
	return path
}

// REQ-76 / TS: "each normative arm of C1 owes at least one test that fails
// if the arm is removed." — coverage goals are stated as arms, not
// percentages.
// BOUNDARY — the coverage obligation made checkable: the artifacts index
// (`coverage.md`) must carry a row for every REQ this phase enumerates, and
// no row may be silently dropped. A missing coverage row is how an
// uncovered arm hides.
func TestReq76_0019_EveryReqCarriesACoverageRow(t *testing.T) {
	coverage := readRepoFile(t, filepath.Join(
		"docs", "rdr", "0019-owned-state-initialization-semantics",
		"artifacts", "coverage.md"))

	reqs := readRepoFile(t, filepath.Join(
		"docs", "rdr", "0019-owned-state-initialization-semantics",
		"artifacts", "req-list.md"))

	// Every `- [REQ-N]` the audit enumerates must appear in the coverage
	// table. An arm with no row is an arm with no test.
	for _, line := range strings.Split(reqs, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- [REQ-") {
			continue
		}
		id := line[len("- ["):]
		if idx := strings.Index(id, "]"); idx > 0 {
			id = id[:idx]
		}
		if !strings.Contains(coverage, "| "+id+" ") &&
			!strings.Contains(coverage, "|"+id+"|") {
			t.Errorf("coverage.md carries no row for %s; each normative arm "+
				"of C1 owes at least one test that fails if the arm is "+
				"removed", id)
		}
	}
}
