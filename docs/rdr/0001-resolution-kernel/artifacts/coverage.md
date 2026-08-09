# Coverage — RDR 0001 Resolution Kernel Contract

Phase 1 (test author) artifact. Maps every REQ in `req-list.md` to the tests
that would fail if a future change broke that clause.

- Test package: `internal/resolve` (external test package `resolve_test`).
- Files: `internal/resolve/resolve_test.go`, `internal/resolve/mvv_test.go`,
  `internal/resolve/fixtures_test.go` (fixtures only),
  `internal/resolve/boundary_test.go` (package-boundary helpers only),
  `internal/resolve/adversarial_test.go` (Phase 3b),
  `internal/resolve/fixup_test.go` (Phase 3c regression).
- Red gate: **confirmed**. 38 of 43 top-level tests fail against the
  logic-free skeleton; the 5 that pass are negative-contract guards whose
  bite was verified by temporarily violating each prohibition (see
  "Negative-contract guards" below).
- Green gate (Phase 2): **confirmed**. 43 of 43 top-level tests pass against
  the implemented kernel, race-clean, with `golangci-lint run ./...` reporting
  `0 issues`. See "REQ-MVV end-to-end run" below for the recorded output.
- Green gate (Phase 3c): **confirmed**. After the three-defect fixup, the
  whole package is green — **139 passing assertions across 56 top-level
  tests, 0 failures**, race-clean (`go test -race -count=1
  ./internal/resolve/...`), with `golangci-lint run ./...` reporting
  `0 issues`. REQ-MVV re-run end-to-end and still passing; its recorded
  output below is unchanged apart from Go's map-iteration subtest ordering,
  which is harness scheduling and not kernel output.

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

## Phase 3c regression tests (`fixup_test.go`)

Added when the three Phase 3 defects were fixed. Each pins a sub-case that
Phase 3a's chain-of-verification reached by spec-derived probe but that no
committed test covered, so the uniform viability gate cannot regress on
those paths. **Every case was verified red against the pre-fix kernel and
green after** — none is a tautology.

| Test | Closes | REQs | Label |
| --- | --- | --- | --- |
| `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | FAIL-1 sub-case 1d | REQ-5, REQ-12, REQ-23 | ADVERSARIAL |
| `TestFixup1e_AmbiguousClassEscapeIsGatedLikeAnyOtherCandidate` | FAIL-1 sub-case 1e | REQ-5, REQ-12, REQ-15 | ADVERSARIAL |
| `TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder` | FAIL-3 third surface | REQ-1, REQ-2, REQ-10 | DOMAIN EDGE |
| `TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder` | FAIL-3, escape call site | REQ-1, REQ-10 | DOMAIN EDGE |
| `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates` | FAIL-1 + FAIL-2 jointly | REQ-5, REQ-15, REQ-23 | ADVERSARIAL |

Notes on what each buys beyond the Phase 3b adversarial suite:

- **1d** — the escape bypass was not specific to a *live* seam. A nil
  `Guards` seam was bypassed too, whereas an ordinary guarded row with a nil
  seam already refused correctly. A fix that only routed escape rows through
  an existing evaluator would still have leaked here.
- **1e** — the bypass was not specific to the `no_match` class. The
  `ambiguous_match` rescue class bypassed identically, so both rescuable
  classes are pinned, with a control asserting the gate rejects *unviable*
  escapes rather than disabling the rescue mechanism.
- **3c** — `rowRefs` feeds three refusal surfaces, not the two ADV-3
  recorded. Both the ordinary ambiguity payload and the degraded
  escape-ambiguity payload (a distinct call site building rows from the
  escape candidate set) are covered.
- **Uniformity** — the property test that would have caught FAIL-1 and
  FAIL-2 together: an ordinary edge and an escape edge presented with the
  same blocking condition must receive the same refusal kind. This is the
  assertion a future path-specific shortcut would break first.

## Orphans

- **REQs with no test**: none. All 37 REQs plus REQ-MVV are covered.
- **Tests citing no REQ**: none. Every `Test*` function opens with a
  `// REQ-N: "<quote>"` comment, and the Phase 3b/3c tests open with the
  Failure Modes clause plus the REQs they read it against. Non-test helpers
  in `fixtures_test.go` and `boundary_test.go` declare no REQ because they
  assert nothing.

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

**Phase 2 disposition**: recorded as `D1` in `deviations.md`, Type SPEC-UNDER,
Status *needs author decision*. The implementation continues with the narrow
reading — escapes rescue `no_match` and `ambiguous_match` only — grounded in
RDR 0001's own Risks clause ("any priority or escape must be explicit table data
owned by RDR 0002") plus RDR 0002's normative restriction of `escape` lists to
those two classes. No frozen test was changed, and the broad reading remains a
one-line routing change if the author decides otherwise.

---

## REQ-MVV end-to-end run

Command:

```
go test -v -count=1 -run 'TestMVV_ReplayDeterminismAndFiveValueLevelRefusals' ./internal/resolve/
```

Actual output:

```
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/1_replay_yields_value_identical_dispositions
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/no_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/ambiguous_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/owned_state_unavailable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/guard_unevaluable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/unmodeled_outcome
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/ambiguous_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/owned_state_unavailable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/guard_unevaluable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/unmodeled_outcome
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/no_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/no_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/ambiguous_match
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/owned_state_unavailable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/guard_unevaluable
=== RUN   TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/unmodeled_outcome
--- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals (0.00s)
    --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/1_replay_yields_value_identical_dispositions (0.00s)
    --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/no_match (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/ambiguous_match (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/owned_state_unavailable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/guard_unevaluable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/2_one_value_level_refusal_per_kind/unmodeled_outcome (0.00s)
    --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/ambiguous_match (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/owned_state_unavailable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/guard_unevaluable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/unmodeled_outcome (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/3_refusals_use_neither_the_cli_nor_the_go_error_path/no_match (0.00s)
    --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/no_match (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/ambiguous_match (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/owned_state_unavailable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/guard_unevaluable (0.00s)
        --- PASS: TestMVV_ReplayDeterminismAndFiveValueLevelRefusals/4_replay_holds_for_refusal_dispositions_too/unmodeled_outcome (0.00s)
PASS
ok  	github.com/newcoinc/intrastate/internal/resolve	0.189s
```

All three MVV obligations are satisfied by the run above: leg 1 compares
replayed dispositions value-for-value (next tags, writes, revision `rev-1`);
leg 2 produces one value-level refusal for each of the five kinds and checks
closure against `RefusalKinds()`; leg 3 asserts a nil Go error, a refusal on the
`Result` value, that neither `Refusal` nor `Result` implements `error`, and that
the kernel package imports no CLI package. Leg 4 extends replay determinism to
all five refusal dispositions.

## Full-suite result (Phase 2)

```
$ go test -count=1 ./...
ok  	github.com/newcoinc/intrastate/internal/cli	0.187s
ok  	github.com/newcoinc/intrastate/internal/resolve	0.355s
(all other packages: no test files)

$ go test -race -count=1 ./internal/resolve/
ok  	github.com/newcoinc/intrastate/internal/resolve	1.308s

$ golangci-lint run ./...
0 issues.
```

43 of 43 top-level tests in `internal/resolve` pass; 0 fail. No pre-existing
failure was observed anywhere in the module.
