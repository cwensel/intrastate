# Coverage — RDR 0001 Resolution Kernel Contract

Phase 1 (test author) artifact. Maps every REQ in `req-list.md` to the tests
that would fail if a future change broke that clause.

- Test package: `internal/resolve` (external test package `resolve_test`).
- Files: `internal/resolve/resolve_test.go`, `internal/resolve/mvv_test.go`,
  `internal/resolve/fixtures_test.go` (fixtures only),
  `internal/resolve/boundary_test.go` (package-boundary helpers only).
- Red gate: **confirmed**. 38 of 43 top-level tests fail against the
  logic-free skeleton; the 5 that pass are negative-contract guards whose
  bite was verified by temporarily violating each prohibition (see
  "Negative-contract guards" below).

---

## REQ × test

| REQ | Test | Label |
| --- | --- | --- |
| REQ-1 | `TestReq1_SameInputTupleReturnsSameDisposition` | HAPPY PATH |
| REQ-1 | `TestReq1_DispositionIsExactlyOneOfPlanOrRefusal` | BOUNDARY |
| REQ-2 | `TestReq2_EveryTupleMemberParticipatesInTheDisposition` | HAPPY PATH |
| REQ-3 | `TestReq3_ValueLevelReplayIsIndependentOfInputSliceOrder` | DOMAIN EDGE |
| REQ-4 | `TestReq4_ReplayIsValueIdenticalAndDoesNotMutateCallerInput` | HAPPY PATH |
| REQ-5 | `TestReq5_EachRefusalConditionRefusesInsteadOfGuessing` | ADVERSARIAL |
| REQ-6 | `TestReq6_ModeledRefusalsDoNotUseTheGoErrorPath` | ADVERSARIAL |
| REQ-7 | `TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds` | BOUNDARY |
| REQ-8 | `TestReq8_RefusalKindAloneDrivesCLIMappingWithoutErrorStrings` | DOMAIN EDGE |
| REQ-9 | `TestReq9_EveryRefusalCarriesAKindFromTheClosedSet` | ADVERSARIAL |
| REQ-10 | `TestReq10_RefusalCarriesInputTupleIdentityAndTableRevision` | DOMAIN EDGE |
| REQ-11 | `TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility` | ADVERSARIAL |
| REQ-12 | `TestReq12_RefusalNeverCarriesAGuessedTransition` | ADVERSARIAL |
| REQ-13 | `TestReq13_KernelDoesNotExecuteAccessorsToFillMissingOwnedState` | ADVERSARIAL |
| REQ-14 | `TestReq14_KernelIsStatelessAcrossInterleavedCalls` | ADVERSARIAL |
| REQ-15 | `TestReq15_ExactlyOneMatchAfterGuardsIsTheOnlySuccess` | BOUNDARY |
| REQ-15 | `TestReq15_AllFourListedConditionsRefuseWithoutAnEscapeEdge` | DOMAIN EDGE |
| REQ-15 | `TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce` | DOMAIN EDGE |
| REQ-16 | `TestReq16_ZeroAndMultipleMatchesRefuseUnlessEscapeIsModeled` | BOUNDARY |
| REQ-17 | `TestReq17_OwnedObservedAndRecognizedTagsAllReachSelection` | DOMAIN EDGE |
| REQ-18 | `TestReq18_SuccessCarriesNextTagsAndOwnedTagWrites` | HAPPY PATH |
| REQ-18 | `TestReq18_RefusalCarriesNoWriteDescription` | ADVERSARIAL |
| REQ-19 | `TestReq19_WritesAreDescribedNotApplied` | DOMAIN EDGE |
| REQ-19 | `TestReq19_OnlyOwnedProvenanceTagsAppearInWrites` | DOMAIN EDGE |
| REQ-20 | `TestReq20_EmptyInputTupleStillYieldsAValueDisposition` | INPUT EDGE |
| REQ-21 | `TestReq21_DispositionsAreStructuredValuesNotProse` | DOMAIN EDGE |
| REQ-22 | `TestReq22_PackageIsInternalAndExposesAPureEntryPoint` | BOUNDARY |
| REQ-23 | `TestReq23_GuardEvaluationIsDelegatedToTheInjectedSeam` | DOMAIN EDGE |
| REQ-24 | `TestReq24_KernelAddsNoStateMutationSemantics` | ADVERSARIAL |
| REQ-25 | `TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence` | DOMAIN EDGE |
| REQ-26 | `TestReq26_KernelUsesNoThirdPartyDependency` | BOUNDARY |
| REQ-27 | `TestReq27_ConcurrentResolutionsAreIndependent` | ADVERSARIAL |
| REQ-28 | `TestReq28_KernelIsTestableWithoutAnyCLICommand` | BOUNDARY |
| REQ-29 | `TestReq29_ReplayReturnsValueIdenticalPlansIncludingWrites` | HAPPY PATH |
| REQ-30 | `TestReq30_NoMatchRefusesWithoutPersistenceSideEffect` | DOMAIN EDGE |
| REQ-31 | `TestReq31_AmbiguousMatchRefusesAndNamesTheConflictingRows` | DOMAIN EDGE |
| REQ-32 | `TestReq32_UnmodeledOutcomeRefusesAndIsNotAMereNoMatch` | DOMAIN EDGE |
| REQ-33 | `TestReq33_UnavailableOwnedStateAndUnevaluableGuardAreValueRefusals` | DOMAIN EDGE |
| REQ-34 | `TestReq34_DispositionIsAFunctionOfTheSuppliedTupleAlone` | BOUNDARY |
| REQ-35 | `TestReq35_RefusalDispositionsReplayValueIdentically` | HAPPY PATH |
| REQ-36 | `TestReq36_KernelExposesNoEncodeDecodeOrInverseOperation` | BOUNDARY |
| REQ-37 | `TestReq37_KernelIntroducesNoHashOrCanonicalSerialization` | BOUNDARY |
| REQ-MVV | `TestMVV_ReplayDeterminismAndFiveValueLevelRefusals` | HAPPY PATH |

## Orphans

- **REQs with no test**: none. All 37 REQs plus REQ-MVV are covered.
- **Tests citing no REQ**: none. Every `Test*` function opens with a
  `// REQ-N: "<quote>"` comment. Non-test helpers in `fixtures_test.go` and
  `boundary_test.go` declare no REQ because they assert nothing.

## REQ-MVV decomposition

`TestMVV_ReplayDeterminismAndFiveValueLevelRefusals` is one runnable
end-to-end test carrying all three MVV obligations plus a replay leg for
refusals:

| Subtest | Obligation |
| --- | --- |
| `1_replay_yields_value_identical_dispositions` | Obligation 1 — value-identical dispositions across two calls on the same tuple, compared field-for-field (next tags, writes, revision), not by exit status. |
| `2_one_value_level_refusal_per_kind` | Obligation 2 — one value-level case per each of the five kinds, plus a closure check against `RefusalKinds()`. |
| `3_refusals_use_neither_the_cli_nor_the_go_error_path` | Obligation 3 — nil Go error, refusal on the `Result` value, the refusal value does not implement `error`, and the kernel package imports no CLI package. |
| `4_replay_holds_for_refusal_dispositions_too` | Replay determinism extended to refusals (REQ-1 covers dispositions, not only plans). |

RDR 0001 declares **no** Round-Trip / Inverse Invariant ("No encode/decode,
import/export, or inverse operation is introduced by this RDR"), so the MVV's
replay obligation is value-for-value equality of the disposition rather than a
reconstruct-and-compare round trip. REQ-36 asserts that non-goal directly.

## Negative-contract guards

Five tests pass against the logic-free skeleton because they assert the
*absence* of a forbidden dependency or symbol, which the skeleton correctly
satisfies: REQ-11, REQ-25, REQ-26, REQ-36, REQ-37.

These are not tautological. Their bite was verified by temporarily adding
`crypto/sha256`, `encoding/json`, `os`, and `github.com/spf13/cobra` imports
plus exported `Apply`, `Encode`, and `Hash` functions to the kernel package:
all five flipped to FAIL, and were restored to green only by removing the
violations.

Six further tests initially passed for a genuinely tautological reason — the
stub returned an identical empty `Result{}` from every call, so "two
dispositions are equal" held vacuously. Each was anchored to require a
substantive disposition (a real plan, or a refusal of the expected kind)
before the equality comparison runs: REQ-3, REQ-14, REQ-22, REQ-24, REQ-27,
REQ-34, REQ-35. All are now red.

## Spec tension noted (not silently narrowed)

REQ-15 as quoted lists **four** conditions as escapable — "Zero, multiple,
unavailable, or unevaluable candidates are refusals unless the table contains
a modeled escape edge". RDR 0002 narrows `escape` lists to `no_match` and
`ambiguous_match` only, and `req-list.md` records that narrowing as an
ASSUMPTION.

RDR 0001 is the governing spec for this run, so
`TestReq15_AllFourListedConditionsRefuseWithoutAnEscapeEdge` tests REQ-15 as
written: all four listed conditions refuse when **no** escape edge is modeled.
That assertion is true under both readings, so the test does not prejudge the
tension. `TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce` and
`TestReq16_*` exercise escape *rescue* only for `no_match` and
`ambiguous_match` — the two classes both RDRs agree are escapable. No test
asserts that an escape edge can or cannot rescue `owned_state_unavailable` or
`guard_unevaluable`; that clause needs a cross-RDR decision before it is
pinned by a test.
