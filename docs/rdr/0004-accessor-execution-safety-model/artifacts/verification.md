
## Phase 3b — Adversarial failure modes (independent)

Written without reading the Phase 1 test files or `coverage.md`. Probes target
`0004:FM`'s definition of a silent failure: "any outcome where the accessor
returns a **success-shaped** result over state it did not establish."

Tests live in `internal/accessor/adversarial_0004_test.go`. All three FAIL
against the implementation at this commit; the rest of the accessor suite stays
green.

### ADV-1 — a plan naming a non-owned tag is applied and reported success

`TestAdv1_WriteRefusesAPlanNamingANonOwnedTag` — **FAILS**.

`Validate` checks only the writer *definition's* declared `keys` against
`Registry.OwnedTags` (`validate.go::Validate`, `CodeWriteNonOwnedTag`). Nothing
in `executor.go::Write` constrains the keys in `plan.Writes` to
`Definition.OwnedKeys()`, so a plan naming an observed or recognized tag reaches
`WriteBinding.Apply` unfiltered and mutates the authoritative artifact.

The corruption is doubly invisible: `executor.go::protectedKeys` excludes every
planned key from the pre-write protected snapshot, so the clobbered observed tag
is dropped from the non-owned identity comparison, and `verifyReadBack` then
confirms the clobbered value equals the plan. Read-back reports green over an
artifact tag the accessor had no authority to touch.

Observed: `Write` returned success with
`Written=[{status final} {profile small}]`, where `profile` is an observed tag
absent from `OwnedTags`.

Defends `0004:C10` / REQ-40, REQ-41; REQ-58; REQ-80.

### ADV-2 — a failed pre-write snapshot makes the protected-tag guard vacuous

`TestAdv2_UnreadablePreWriteSnapshotMustNotYieldSuccess` — **FAILS**.

`executor.go::Write` takes the pre-write protected snapshot through `invokeRead`
and guards it with `if raw.class == ""`. On a timeout, execution failure, or
incomplete pre-write read the failure is discarded silently, `before` stays
empty, and execution proceeds to full success. `verifyReadBack` iterates
`before`, so an empty map makes the `0004:C12` clause "observed and recognized
tag values present before the write are unchanged" vacuously true. The same hole
exists one level down: `raw.classify(protected, false)` discards its `unread`
return, narrowing the protected set key by key.

A write binding that also mutates a protected observed tag therefore passes
read-back green. This is the RDR's own risk row ("A write command succeeds but
changes the wrong artifact or wrong tag") with its stated mitigation disabled.
The contract's vocabulary for "the verification did not run" already exists:
`read_back_incomplete` (`0004:C13`), which must be neither mismatch nor success.

Observed: `Write` returned success with `Written=[{status final}]` while the
artifact's `profile` moved from `large` to `small` unnoticed.

Defends `0004:C12` / REQ-53, REQ-58; REQ-62; REQ-80.

### ADV-3 — a verified removal is recorded as holding the reserved `<clear>` literal

`TestAdv3_ClearingWriteDispositionMustNotRecordTheClearLiteral` — **FAILS**.

On success `executor.go::Write` returns `WriteResult{Written: planned}`, and
`planned` still carries `Tag{Key: k, Value: "<clear>"}` for a cleared key.
`disposition.go::WriteDisposition` renders that set verbatim through
`tagsAsValues` + `observedTags`; `observedTags` drops only `KeyValue.Absent`, and
a tag rebuilt from `planned` is never marked absent. So the replay-stable
disposition for a *verified removal* asserts the artifact HOLDS the key with the
reserved literal as its value — while read-back verified that key ABSENT.

The artifact is correct here; the *record* is the placeholder. That is the exact
encoding the RDR names as the regression: "A sentinel value would make an absent
key read as *present* and silently retire `owned_state_unavailable` for it" (LBD,
Absence crosses the seam as omission). `ReadResult.OwnedSnapshot` gets this right
by omitting absent keys; the write disposition path does not, and the disposition
is what a replay consumes (CA A4). `WriteResult.Written` carries the same defect —
it is documented as "the planned owned tags the read-back verified", and what
read-back verified for a cleared key is absence.

Observed: `WriteDisposition(...).Tags` contained `{Key:labels Value:<clear>}`
after a clear whose read-back correctly asserted absence.

Defends `0004:C11` / REQ-45, REQ-46; REQ-32; REQ-80.

### Verdict

`BLOCK` — three real defects, all on the write path, all producing a
success-shaped result over unestablished or misrecorded authoritative state.
ADV-1 and ADV-2 are artifact-corruption paths; ADV-3 is a record-corruption path
that reintroduces the sentinel encoding this RDR exists to forbid. Phase 3c owns
the fixes.

---

# Phase 3a — CoVe verification (independent probes)

Method: written blind to the Phase 1 suite, `coverage.md`, and
`verification.md`'s prior content. For each REQ an input was designed that
would make a *correct* implementation visibly violate it, then run against
`internal/accessor/` via throwaway probe harnesses (`cove_probe*_test.go`,
deleted before commit). Every FAIL-N below was **executed and observed** —
no suspicion is recorded as a finding. Each carries a *passing negative
control* proving the probe discriminates rather than refusing everything.

## Probes run and found conforming

| Probe | REQs exercised | Result |
| --- | --- | --- |
| A | 27, 119, 68, 69 | `timeout` outranks `incomplete_read` when both are true at once; truncation-only control still `incomplete_read` |
| B | 19–24, 29–32, 48, 115 | unrequested key dropped; refusal carries no values and empty snapshot; a silently-dropped key defaults to *unreadable*, not guessed absence; `<clear>` read value refuses `incomplete_read`; genuine absence takes the success branch and is **omitted** from `OwnedSnapshot`, placeholder included |
| C | 2–6, 85 | same-id read/write pair is two identities (not multiply-bound); a true duplicate triple is; off-capability invoke → `capability_mismatch`; unbound → `unknown_accessor`; read selection never reaches the same-id write binding |
| D | 42, 43, 45–47, 49, 53, 55, 57, 62–65, 67, 69, 120 | success/clear/idempotent-clear/stored-literal/owned-mismatch/non-owned-clobber/read-back-incomplete/post-mutation-timeout/exec-failure all classify correctly; `applied=true` exactly on the post-mutation refusals; exactly one `Apply` invocation, no undo; empty write set runs no write |
| F | 50, 54 | read-back selects the **write binding's** role, not a decoy role already holding the expected value |
| G | 35–39 | allow/deny typed with reason (deny is *not* a refusal here); indeterminate is refusal-class; gate timeout → `timeout`, gate error → `execution_failure`, neither folded into `gate_indeterminate` |
| H | 77, 86–89 | `resolve.RefusalKinds()` unextended and disjoint; dispositions compare equal across runs with tag order normalized; injected-refusal run differs from the success run (ORA 4 control) |
| I | 11, 41, 70, 83–85 | all eight arms mint their own named code; valid registry returns the empty set |
| J | 12, 16, 18, 26 | unsupplied role → `execution_failure`; the binding is handed the **definition's** key set, never one it derived; refusal names the unread keys |
| L | 72–74 | stdout/stderr captured around read/gate/write/validate/disposition batteries — both empty (0 bytes) |
| M | 45–47, 62–65, 68 | clear alongside protected keys, mixed assign+clear, unreadable cleared key → `read_back_incomplete`, post-write re-read timeout → `timeout` with `applied=true` |
| N | 14, 16, 52 | set-valued read-back is **byte** equality — a reordered or whitespace-differing JSON array refuses `read_back_mismatch`; a binding mutating its handed key slice cannot corrupt the definition |
| O | 31, 118 | absent required key omitted from the snapshot; carried-key control preserves it |

## FAIL-1 — a write reports SUCCESS when the pre-write snapshot could not be taken

- **Violates**: REQ-53 (`0004:C12` "and that observed and recognized tag
  values present before the write are unchanged"), REQ-80 (a silent failure
  is "any outcome where the accessor returns a **success-shaped** result over
  state it did not establish"), REQ-104, REQ-112.
- **Failing input** (probe E/K): registry with reader `rd` on role `state`
  declaring `keys = ["status","profile"]`, writer `wr` on role `state`
  declaring `keys = ["status"]`, `OwnedTags = ["status"]`. Artifact holds
  `status=old, profile=keep`. The **pre-write snapshot read** fails — probed
  in three independent modes, each reproducing: (a) it names every requested
  key unreadable, (b) it returns an execution error, (c) it exceeds the
  reader's timeout. The write binding then applies `status=new` **and**
  clobbers `profile` to `CLOBBERED`. Plan: `Writes=[{status,new}]`.
- **Observed**: `WriteResult{Written:[{status new}]}`, `Refused()==false` —
  a success. Artifact ends `map[profile:CLOBBERED status:new]`.
- **Spec-required**: a refusal. The non-owned comparison the contract
  mandates never ran, so the executor cannot assert the protected values are
  unchanged; per REQ-62's own logic ("the verification did not run") this is
  `read_back_incomplete`, and per REQ-80 it may not be success-shaped.
- **Root cause**: `internal/accessor/executor.go:269-281`. The pre-write read's
  refusal class is discarded by the `if raw.class == ""` guard, leaving
  `before` empty; `verifyReadBack`'s protected loop then iterates an empty
  map and finds nothing to compare.
- **Negative control (passes)**: probe Q — with a fully readable pre-write
  snapshot the identical `profile` clobber IS caught as `read_back_mismatch`.
  The check works; it is silently skipped exactly when the snapshot fails.

## FAIL-2 — a *partial* pre-write snapshot silently drops the unread keys' protection

- **Violates**: REQ-53, REQ-80. A narrower, more reachable variant of FAIL-1.
- **Failing input** (probe P): as above plus a third non-owned key `owner`
  (`rd` declares `["status","profile","owner"]`, artifact holds
  `owner=alice`). The pre-write snapshot read **succeeds** for `status` and
  `profile` but names only `owner` unreadable. The write applies
  `status=new` and clobbers `owner` to `MALLORY`.
- **Observed**: `Refused()==false`, `Written=[{status new}]`; artifact ends
  `map[owner:MALLORY profile:keep status:new]`.
- **Spec-required**: a refusal — `owner`'s pre-write value was never
  established, so its "unchanged" claim is unverifiable.
- **Root cause**: `internal/accessor/executor.go:273` —
  `vals, _ := raw.classify(protected, false)` discards the `unread` return.
  Unlike the post-write read-back at line 338, which correctly refuses
  `read_back_incomplete` on a non-empty `unread`, the pre-write path drops it.
- **Negative control (passes)**: probe Q, as for FAIL-1.

## FAIL-3 — a plan naming a non-owned tag is applied at runtime and reports success

- **Violates**: REQ-40 (`0004:C10` "A write accessor MUST apply only planned
  owned-tag writes produced by a successful transition"), REQ-41 ("It MUST
  NOT write observed or recognized tags").
- **Failing input** (probe R1): a registry that is **valid** — `Validate`
  returns the empty set, because the writer's *declared* `keys` are
  `["status"]` and `OwnedTags = ["status"]`. The runtime **plan** then names
  a non-owned tag: `Writes = [{status,final}, {profile,small}]`.
- **Observed**: `Refused()==false`,
  `Written=[{status final} {profile small}]`; the artifact's non-owned
  `profile` is mutated from `large` to `small`.
- **Spec-required**: a refusal. `0004:C10` binds the *write*, not only the
  definition. `Executor.Write` passes `plan.Writes` to `binding.Apply`
  unfiltered (`executor.go:251,285`) with no check against
  `def.OwnedKeys()`/`reg.OwnedTags`; `CodeWriteNonOwnedTag`
  (`validate.go:64-73`) inspects only declared metadata, so nothing enforces
  C10 on the plan actually executed.
- **Reading taken**: REQ-41's annotation names a definition-level validation
  arm, and that arm exists. But REQ-40's obligation is on what the accessor
  *applies*, and REQ-80 forbids a success-shaped result over unestablished
  state. A plan is caller-supplied at this boundary (RDR 0005 is
  unimplemented), so an unchecked plan is exactly the smuggling path
  REQ-110's mitigation exists to close. Recorded as a violation of REQ-40.
- **Negative control (passes)**: probe R2 — the same fixture with an
  owned-only plan succeeds, so this is not an executor that refuses
  everything.

## FAIL-4 — a cleared key is recorded as holding the reserved literal `<clear>`

- **Violates**: REQ-46 (`0004:C11` — the read-back "MUST assert the key is
  absent"), REQ-32 (a reserved sentinel must not stand in for a key the
  artifact does not hold), REQ-86/REQ-89 (the replay-stable disposition
  records "the returned tag values"), REQ-80.
- **Failing input** (probe R3): writer `wr` and reader `rd` on role `state`,
  `OwnedTags=["labels"]`, artifact holds `labels=["a"]`. Plan:
  `Writes=[{labels,<clear>}]`. The binding correctly **removes** the key;
  read-back correctly verifies absence and the write succeeds.
- **Observed**: artifact ends `map[]` (key genuinely gone), yet
  `WriteResult.Written = [{Key:labels Value:<clear>}]` and
  `WriteDisposition(...).Tags = [{Key:labels Value:<clear>}]`.
- **Spec-required**: the cleared key must be **omitted** from the recorded
  tag set. The record asserts the artifact holds `labels` with the reserved
  literal — precisely the encoding LBD Clear semantics calls out as "the one
  write whose read-back passes green over a tag that was never removed", and
  the same placeholder-for-absence shape REQ-32 forbids at the seam. A
  replay consumer reading this disposition reconstructs a state the artifact
  does not have.
- **Root cause**: `executor.go:359` returns `WriteResult{Written: planned}`
  — the raw plan, `<clear>` sentinels included — and
  `disposition.go:68` renders it through `tagsAsValues`, which sets
  `Absent:false` for every tag, so `observedTags`' absence filter never
  fires.
- **Negative control (passes)**: probe D1/M2 — an assignment write records
  its real value, so the recorder is not simply dropping everything.

## Note on independence

FAIL-1 and FAIL-4 correspond to failures also reported by Phase 3b's
`internal/accessor/adversarial_0004_test.go` (Adv2, Adv3), and FAIL-3 to its
Adv1. That suite was committed to this branch before this phase ran and was
**not** consulted while designing these probes; each finding above was
reproduced independently by its own harness with its own negative control.
The agreement is corroboration, not restatement. FAIL-2 is additional: the
*partial* pre-write snapshot is a distinct and more reachable path than the
wholly-unreadable snapshot, and it is not closed by fixing the `raw.class`
guard alone — `classify`'s discarded `unread` return must also be honoured.
