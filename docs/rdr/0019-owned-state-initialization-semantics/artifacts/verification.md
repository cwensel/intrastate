
## Phase 3b — adversarial review (`0019:F1`–`F5`)

Independent pass over the record's `Failure Modes`, written without reading
Phase 3a's findings. Tests live in `internal/cli/flow_initstate_adv_0019_test.go`.

### ADV-1 — F2, the torn state the VERB ITSELF produces — PASSES

**Failure mode.** `0019:F2` "Torn (multi-writer)": "no cross-writer atomicity
— a failed writer can leave the artifact seeded for some keys only. The store
is then non-empty, so a re-run is a NO-OP whose payload lists the missing
`[initial]` keys … but it does NOT repair them." (REQ-47, REQ-51, REQ-60,
REQ-67.)

**The gap in the shipped suite.** S12
(`TestReq51And93_0019_ATornStateIsReportedNotRepaired`) constructs its torn
state with a prior `flow set-state`. No shipped test drives `init-state`'s own
write loop to a partial commit, so F2's central claim was asserted only over a
state the verb did not create.

**The test.** `TestAdv1_0019_ATornSeedTheVerbItselfProducesIsReportedNotRepaired`.
A two-role model whose `aux` writer — first in the loop's sorted name order —
declares an unreachable read-back locator. The seeding arm commits `note` into
a sealed artifact, refuses at exit 3, and never reaches `state`. Asserts: the
partial commit is real and the unreached role's artifact does not exist; the
re-run is a no-op success leaving both artifacts byte-identical; the re-run's
payload names exactly `{stage}`.

**Result: PASSES.** The implementation defends F2's torn case. Kept for
regression value the shipped suite lacks — it is the only test whose torn
state the verb itself produced.

### ADV-2 — F2, the absent report is answered from the UNION of stores — FAILS

**Failure mode.** `0019:F2`: "The payload's absent-key list is what keeps the
manual route mechanical — it names precisely the keys to pass." REQ-47 fixes
that the no-op payload "reports which `[initial]` keys THE STORE does not
carry", and the record's repair route is "explicit `set-state` of the listed
keys" — which routes each key to its own writer's artifact.

**The defect.** `internal/cli/flow_initstate.go::initAbsentKeys` folds every
role's key set into one `present` map before testing membership:

```go
present := map[string]bool{}
for _, keys := range stores {
	for _, key := range keys {
		present[key] = true
	}
}
```

A key carried by ANY bound store suppresses its absent entry, regardless of
which writer serves it. Every other quantifier in the verb is role-scoped
(`initNeededRoles`, `initStoresEmpty` over per-role key sets); this one is not.

**The test.**
`TestAdv2_0019_TheAbsentReportIsAnsweredPerRoleNotFromTheUnionOfStores`.
Built through the CLI only, with no hand-staged artifact: both roles of
`initTwoRoleModel` are bound to one shared artifact and `note` is established
there by `set-state`; role B is then rebound to a fresh artifact (a
per-invocation binding change no clause forbids). Role A's store now carries
`note` — a key only role B's writer serves — while role B's store is empty.
The correct absent set is `{note, stage}`; the implementation reports
`{stage}`.

**Result: FAILS.** Observed payload: `absent: ["stage"]`.

**Harm, probed separately and confirmed.** An operator who runs exactly and
only what the payload lists — `flow set-state --write stage=seeded` — is left
with `note` still reading ABSENT from `flow read-state`. F2's claim that the
list "names precisely the keys to pass" is false on a multi-role model.

**Suggested repair (a later step's job; no production code touched here).**
Answer `initAbsentKeys` per role: for each `[initial]` assignment, test its key
against the key set of the store its own writer's role names, rather than
against the union.

### Discarded probe — F1 exit-group confusion — PASSED, redundant

"An exit-3 seed leaves an APPLIED store the re-run will not reseed" (REQ-67,
F1's "refused before writing" vs "applied but unverified"). Written, PASSED,
and discarded rather than kept: it is
`TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore` under a
different name. The implementation defends this mode; a redundant copy adds no
regression value. ADV-1 asserts the same disposition incidentally, on the
harder multi-writer construction.
