
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
