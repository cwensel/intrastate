# Coverage — RDR 0026 bounded-output-drain-on-command-timeout

Phase 1 (tests first). Column 1 is the REQ id from `req-list.md`; column 2
names every test that would fail if that clause broke. An EMPTY column 2 is
the orphan mark.

Test files:

- `internal/cli/cmdbind/fixtures_0026_test.go` — FX-deadline-escape (a real
  `SysProcAttr{Setsid: true}` Go helper) and the shared oracles
- `internal/cli/cmdbind/export_0026_test.go` — the ONE drain-start stall
  seam and the F5 pollability injection, both test-only
- `internal/cli/cmdbind/drain_0026_test.go` — C1 `precedence:`, `whole:`
- `internal/cli/cmdbind/refusal_0026_test.go` — C1 `refusal:`, `class:`,
  `stdin:`, `bound:`, `residue:`, F5, S10
- `internal/cli/cmdbind/contract_0026_test.go` — the structural obligations
  (order in `spawn`, the `readBounded` rewrite, no third build tag, the
  Phase 2 / Phase 5 documentation obligations)
- `internal/accessor/held_pipe_0026_test.go` — the held-pipe error type and
  the applied sense on the DIRECT write leg
- `internal/cli/held_pipe_0026_test.go` — S1's rendered half, D1's negative
  witness, D2's distinct exit-3 code
- `internal/cli/command_mvv_0026_test.go` — REQ-MVV, end to end at the
  `flow resolve` boundary

| REQ | Test |
| --- | --- |
| REQ-MVV | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-MVV.1 | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes; TestReq110_TheFixtureGrandchildGenuinelyEscapesTheProcessGroup |
| REQ-MVV.2 | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-MVV.3 | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-MVV.3a | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-MVV.4 | TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-MVV.5 | TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing |
| REQ-MVV.6 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-MVV.6a | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue |
| REQ-MVV.6b | TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-MVV.6c | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-1 | TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes |
| REQ-2 | TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin; TestReq2_ReleaseFirstKeepsTheInGroupGrandchildCaseFarBelowTheBound |
| REQ-3 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-4 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-5 | TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-6 | TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-7 | TestReq7_OverflowOutranksHeldOnTheSameDrain |
| REQ-8 | TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation; TestReq132_StderrHeldAloneNamesTheBarePipeAndKeepsTheTail; TestReq28_TheHeldPipeErrorWrapsRatherThanFlattensAndSurvivesToTheRefusalSite |
| REQ-9 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue |
| REQ-10 | TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead; TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-11 | TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-12 | TestReq7_OverflowOutranksHeldOnTheSameDrain |
| REQ-13 | TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-14 | TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure; TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-15 | TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure |
| REQ-16 | TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-17 | TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-18 | TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-19 | TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes; TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue |
| REQ-20 | TestReq20_NoGraceOnFlagAndNoExplicitCloseAreAddedToSpawn |
| REQ-21 | TestReq21_TheJoinRemainsARealJoinWithAPartialWriteInFlight |
| REQ-22 | TestReq21_TheJoinRemainsARealJoinWithAPartialWriteInFlight |
| REQ-23 | TestReq21_TheJoinRemainsARealJoinWithAPartialWriteInFlight |
| REQ-24 | TestReq24_NoReadEndOutlivesSpawnAcrossRepeatedHeldPipeRefusals; TestReq20_NoGraceOnFlagAndNoExplicitCloseAreAddedToSpawn |
| REQ-25 | TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-26 | TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-27 | TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-28 | TestReq28_TheHeldPipeErrorWrapsRatherThanFlattensAndSurvivesToTheRefusalSite; TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-29 | TestReq29_TheHeldPipeErrorIsAnExportedAccessorTypeMatchableFromOutsideCmdbind |
| REQ-30 | TestReq29_TheHeldPipeErrorIsAnExportedAccessorTypeMatchableFromOutsideCmdbind; TestReq28_TheHeldPipeErrorWrapsRatherThanFlattensAndSurvivesToTheRefusalSite |
| REQ-31 | TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense |
| REQ-32 | TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense; TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode |
| REQ-33 | TestReq33_TheReadBackLegStaysReadBackIncompleteWithAppliedAlreadySet; TestReq122_TheReadBackArmStillRendersTheAppliedProseAndKeepsItsOwnCode |
| REQ-34 | TestReq34_AHeldPipeOnAReadOrGateLeavesAppliedFalse; TestReq34_AReadOrGatePhaseHeldPipeRefusalNeverRendersTheAppliedProse |
| REQ-35 | TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath; TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing |
| REQ-36 | TestReq36_EveryCallerChecksErrBeforeTouchingTheInvocation |
| REQ-37 | TestReq36_EveryCallerChecksErrBeforeTouchingTheInvocation |
| REQ-38 | TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail; TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-39 | TestReq39_TheBindingNeverClassifiesTimeoutItself |
| REQ-40 | TestReq40_ANonReadingChildWithAStdinPayloadPastTheBufferRefusesAfterTheJoin; TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin |
| REQ-41 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-42 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-43 | TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal |
| REQ-44 | TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes; TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes; TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal |
| REQ-45 | TestReq45_AWriteJourneyAgainstAnEscapingHelperStaysInsideThreeBounds |
| REQ-46 | TestReq45_AWriteJourneyAgainstAnEscapingHelperStaysInsideThreeBounds |
| REQ-47 | TestReq45_AWriteJourneyAgainstAnEscapingHelperStaysInsideThreeBounds; TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal |
| REQ-48 | TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal |
| REQ-49 | TestReq49_AChildThatClosesItsPipesObservesNoChange |
| REQ-50 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-51 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-52 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-53 | TestReq53_TheEscapeeOutlivesTheRefusalAndIsDiagnosedByTheNamedPipes; TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-54 | TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes; TestReq53_TheEscapeeOutlivesTheRefusalAndIsDiagnosedByTheNamedPipes |
| REQ-55 | TestReq53_TheEscapeeOutlivesTheRefusalAndIsDiagnosedByTheNamedPipes; TestReq119_TheOutputContractDocNamesTheHeldPipeRefusalAndTheResidue |
| REQ-56 | TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin |
| REQ-57 | TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath |
| REQ-58 | TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes; TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-59 | TestReq2_ReleaseFirstKeepsTheInGroupGrandchildCaseFarBelowTheBound |
| REQ-60 | TestReq49_AChildThatClosesItsPipesObservesNoChange |
| REQ-61 | TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-62 | TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath |
| REQ-63 | TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine |
| REQ-64 | TestReq39_TheBindingNeverClassifiesTimeoutItself; TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-65 | TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead |
| REQ-66 | TestReq20_NoGraceOnFlagAndNoExplicitCloseAreAddedToSpawn |
| REQ-67 | TestReq11_ReadBoundedOwnsAnExplicitReadLoopThatRearmsTheDeadlinePerRead; TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-68 | TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin; TestReq36_EveryCallerChecksErrBeforeTouchingTheInvocation |
| REQ-69 | TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine; TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-70 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-71 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-72 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-73 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-74 | TestReq132_StderrHeldAloneNamesTheBarePipeAndKeepsTheTail; TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-75 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-76 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReq1_ABoundedWaitOutranksAWholeReadWhenAnEscapeeHoldsThePipes |
| REQ-77 | TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-78 | TestReq7_OverflowOutranksHeldOnTheSameDrain |
| REQ-79 | TestReq79_ADrainConditionIsSelectedAheadOfTheNonExitErrorArm |
| REQ-80 | TestReq39_TheBindingNeverClassifiesTimeoutItself; TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-81 | TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense; TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode |
| REQ-82 | TestReq7_OverflowOutranksHeldOnTheSameDrain |
| REQ-83 | TestReq83_TheAppliedSenseAndTheHeldPipeReasonShareTheOneDetailSlotInOrder |
| REQ-84 | TestReq29_TheHeldPipeErrorIsAnExportedAccessorTypeMatchableFromOutsideCmdbind |
| REQ-85 | TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-86 | TestReq49_AChildThatClosesItsPipesObservesNoChange |
| REQ-87 | TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure |
| REQ-88 | TestReq88_TheStderrTailKeepsTheLastFourKiBAcrossTheRewrite |
| REQ-89 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-90 | TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-91 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-92 | TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath; TestReq109_LossAtTheBoundIsAlwaysBoundedToARefusedInvocation |
| REQ-93 | TestReq49_AChildThatClosesItsPipesObservesNoChange |
| REQ-94 | TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing; TestReq26_TheHeldPipeReasonNamesPipesBoundExitStatusAndRemediation |
| REQ-95 | TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-96 | TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense; TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode |
| REQ-97 | TestReq33_TheReadBackLegStaysReadBackIncompleteWithAppliedAlreadySet; TestReq122_TheReadBackArmStillRendersTheAppliedProseAndKeepsItsOwnCode |
| REQ-98 | TestReq40_ANonReadingChildWithAStdinPayloadPastTheBufferRefusesAfterTheJoin |
| REQ-99 | TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine |
| REQ-100 | TestReq53_TheEscapeeOutlivesTheRefusalAndIsDiagnosedByTheNamedPipes |
| REQ-101 | TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing; TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath |
| REQ-102 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead |
| REQ-103 | TestReq132_StderrHeldAloneNamesTheBarePipeAndKeepsTheTail |
| REQ-104 | TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure |
| REQ-105 | TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-106 | TestReq123_TheExistingGrandchildLeakAndCancelOraclesStillShip; TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-107 | TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine |
| REQ-108 | TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail; TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes |
| REQ-109 | TestReq109_LossAtTheBoundIsAlwaysBoundedToARefusedInvocation |
| REQ-110 | TestReq110_TheFixtureGrandchildGenuinelyEscapesTheProcessGroup |
| REQ-111 | TestReq2_SpawnKeepsTheFixedOrderWaitThenReapGroupThenTheDrainJoin; TestReq29_TheHeldPipeErrorIsAnExportedAccessorTypeMatchableFromOutsideCmdbind |
| REQ-112 | TestReq112_TheC4CommentsAtTheSeamAreRepointedAt0026C1; TestReq112_TheLocked0025RecordIsNotEdited |
| REQ-113 | TestReq110_TheFixtureGrandchildGenuinelyEscapesTheProcessGroup; TestReq123_TheExistingGrandchildLeakAndCancelOraclesStillShip |
| REQ-114 | TestReq31_AHeldPipeOnTheDirectWriteSetsTheAppliedSense; TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode |
| REQ-115 | TestReq45_AWriteJourneyAgainstAnEscapingHelperStaysInsideThreeBounds |
| REQ-116 | TestReq116_OnlyTheHeldPipeErrorSetsTheAppliedSenseOnTheSharedWriteArm |
| REQ-117 | TestReq116_OnlyTheHeldPipeErrorSetsTheAppliedSenseOnTheSharedWriteArm; TestReq138_AnUnappliedWritePhaseExecutionFailureRendersNoAppliedProse |
| REQ-118 | SURVEY-118 (Phase 5 survey, recorded below) |
| REQ-119 | TestReq119_TheOutputContractDocNamesTheHeldPipeRefusalAndTheResidue |
| REQ-120 | TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode; TestReq122_TheReadBackArmStillRendersTheAppliedProseAndKeepsItsOwnCode; TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath |
| REQ-121 | TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode |
| REQ-122 | TestReq122_TheReadBackArmStillRendersTheAppliedProseAndKeepsItsOwnCode; TestReq33_TheReadBackLegStaysReadBackIncompleteWithAppliedAlreadySet |
| REQ-123 | TestReq123_TheExistingGrandchildLeakAndCancelOraclesStillShip; TestReq2_ReleaseFirstKeepsTheInGroupGrandchildCaseFarBelowTheBound |
| REQ-124 | TestReq49_AChildThatClosesItsPipesObservesNoChange; TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure; TestReq88_TheStderrTailKeepsTheLastFourKiBAcrossTheRewrite |
| REQ-125 | TestReq40_ANonReadingChildWithAStdinPayloadPastTheBufferRefusesAfterTheJoin |
| REQ-126 | TestReq126_ANonPollablePipeIsRecordedAtTheCheckAndSurfacesAsADetailLine |
| REQ-127 | TestReq127_AHeldPipeWithTheDeadlineAlsoElapsedIsTimeoutWithAnEmptyDetail |
| REQ-128 | TestReq128_TheStallSeamIsNilInProduction; TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-129 | TestReq51_PromptDrainsRecoverTheWholePayloadAndAStalledOneExactlyOnePipeBuffer |
| REQ-130 | TestReq48_AChildThatResistsTerminationStillReturnsInsideTheCommittedTotal |
| REQ-131 | TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue; TestReq6_ATricklingEscapeeIsRefusedThoughItsDrainIsNeverEmpty |
| REQ-132 | TestReq132_StderrHeldAloneNamesTheBarePipeAndKeepsTheTail |
| REQ-133 | TestReq110_TheFixtureGrandchildGenuinelyEscapesTheProcessGroup |
| REQ-134 | TestReq21_TheJoinRemainsARealJoinWithAPartialWriteInFlight |
| REQ-135 | TestReq135_TheStallSeamAddsNoThirdBuildTagAndThePlatformSplitIsUnchanged |
| REQ-136 | TestReq14_TheCapBoundaryIsExactlyAtCapAValueAndCapPlusOneAFailure; TestReq35_AHeldStdoutIsNeverParsedOnTheReadGateOrWritePath |
| REQ-137 | TestReq137_NoOptOutFlagForTheHeldPipeRefusalShips; TestReq42_DrainGraceShipsBesideWaitDelayAsAMechanismConstantNotASecondBound |
| REQ-138 | TestReq138_AnUnappliedWritePhaseExecutionFailureRendersNoAppliedProse; TestReq116_OnlyTheHeldPipeErrorSetsTheAppliedSenseOnTheSharedWriteArm; TestReq34_AHeldPipeOnAReadOrGateLeavesAppliedFalse |
| REQ-139 | TestReq139_TheAppliedWriteCodeIsDistinctFromEveryExistingAccessorCode |
| REQ-140 | TestReq139_TheAppliedWriteCodeIsDistinctFromEveryExistingAccessorCode; TestReq121_TheAppliedDirectWriteRefusalRendersTheMayHaveAppliedProseAndItsOwnCode; TestReq138_AnUnappliedWritePhaseExecutionFailureRendersNoAppliedProse |

## REQ-118 survey result (IP Phase 5, "SURVEY, then state")

Acceptance for this half is "the recorded result, empty or not". The survey
ran at Phase 3c against the branch tip, after the held-pipe refusal and the
sibling-halt drain landed.

**Population.** Every pre-existing reader / gate / write fixture in the Go
suites, plus the repo's own declared bindings exercised by
`internal/table` (the `command` source kind, `internal/table/source.go`) and
the end-to-end `flow resolve` paths in `internal/cli`.

**Method.** `go test ./...` over the whole module on the branch, compared
against the same run on `main` (the pre-change baseline), both green.

**Result: EMPTY.** No binding that passed before the change refuses after it.

```
ok  github.com/cwensel/intrastate/internal/accessor        1.133s
ok  github.com/cwensel/intrastate/internal/cli             7.373s
ok  github.com/cwensel/intrastate/internal/cli/clierr      0.711s
ok  github.com/cwensel/intrastate/internal/cli/cmdbind    17.894s
ok  github.com/cwensel/intrastate/internal/cli/flowbind    1.399s
ok  github.com/cwensel/intrastate/internal/graphlint       5.705s
ok  github.com/cwensel/intrastate/internal/guard           2.010s
ok  github.com/cwensel/intrastate/internal/resolve         2.114s
ok  github.com/cwensel/intrastate/internal/table           1.760s
ok  github.com/cwensel/intrastate/internal/version         2.294s
```

An empty result is the accepting outcome; a NON-empty result would have been
a Stage-6 route-back rather than a doc note. REQ-119, the "then state" half,
is covered by
`TestReq119_TheOutputContractDocNamesTheHeldPipeRefusalAndTheResidue`.

## Orphan


**No orphans remain.** REQ-118 (IP Phase 5, "SURVEY, then state") was the
one uncovered row; the survey has now RUN and its result is recorded above,
which is exactly the clause's stated acceptance. Retained below is why it
carries a recorded result rather than an assertion. The clause's acceptance criterion is
"the recorded result, empty or not" — a Phase-5 *procedure* whose output is
a written finding, not a program behaviour. No assertion can make a build
that skipped the survey fail: the survey's own artifact is the evidence, and
a test that checked for the artifact's presence would assert that a file was
written, not that the population was measured. It is an implementer
obligation for Phase 5, and a NON-empty result is a Stage-6 route-back
rather than a test failure. REQ-119, the "then state" half, IS covered
(`TestReq119_TheOutputContractDocNamesTheHeldPipeRefusalAndTheResidue`),
because that half has a checkable artifact.

## Red-gate result

Run against the UNMODIFIED implementation.

New-behaviour tests confirmed FAILING (38):

- Compile-red (the three surfaces the RDR mandates do not exist yet):
  `accessor.HeldPipeError` (REQ-28..REQ-30, REQ-84, REQ-111),
  `cmdbind.DrainGrace` (REQ-42, REQ-72), the drain-start stall seam
  `drainStartDelay` (REQ-128) and the F5 pollability injection
  `nonPollableForTest` (REQ-126). Every test in the three packages is
  therefore red at compile time; the per-test verdicts below were taken with
  those four surfaces stubbed, so the behavioural red is real and not an
  artefact of the missing declarations.
- Hangs indefinitely (the defect itself — `internal/cli/cmdbind/cmdbind.go:277`,
  the unbounded `readBounded` drain): TestReq1, TestReq6, TestReq7,
  TestReq26, TestReq29, TestReq35, TestReq36, TestReq39, TestReq45,
  TestReq48, TestReq53, TestReq79, TestReq109, TestReq110, TestReq132
- Hangs at the `flow resolve` boundary (60 s against a 2 s / 20 s declared
  timeout — the originally reported "still blocked after 15 s"):
  TestReqMVV, TestReqMVV5
- Fails on an assertion: TestReq11 (`readBounded` still calls `io.ReadAll`
  and arms no deadline), TestReqMVV6a / TestReqMVV6b / TestReq21 /
  TestReq51 (stalled arm) (the stall seam is not honoured, so the final
  drain is never exercised), TestReq112 (the C4 comments are not
  re-pointed), TestReq119 (the output-contract doc names neither the
  held-pipe refusal nor the residue), TestReq31 (`Applied()` is false on a
  held-pipe direct write), TestReq121 / TestReq83 / TestReq139 (the applied
  write refusal renders neither the prose nor a distinct exit-3 code)

Preserved-invariant tests, green by design (18). Each pins a clause the
record explicitly PRESERVES — C1's `stdin:` / `whole:` / `bound:` lines and
the arms C1 `refusal:` says are "already correct" — so passing today is the
point; they are regression traps across the `readBounded` rewrite:

- TestReq49 (a child that closes its pipes observes no change)
- TestReq14 (the 1 MiB cap boundary: exactly-at-cap is a value, cap+1 is
  overflow)
- TestReq88 (the stderr tail keeps the LAST 4 KiB)
- TestReq51 (prompt arm: the whole 1 MiB is recovered)
- TestReq40 (the unchanged stdin leg)
- TestReq24 (no read end outlives `spawn`)
- TestReq2 (both: `spawn`'s fixed order, and release-first keeping the
  in-group case far below the bound)
- TestReq20 (no `graceOn` flag, no explicit close)
- TestReq123 (the existing grandchild-leak and cancel oracles still ship)
- TestReq135 (no third build tag; the platform split is unchanged)
- TestReq34 / TestReq33 / TestReq122 / TestReq127 / TestReq116 / TestReq28 /
  TestReq138 / TestReq137 (the read-back arm, the read/gate applied-false
  arms, the deadline-first class flip, the Phase-4 negative control, and the
  no-opt-out-flag prohibition — all correct today)

## REQ-MVV runner

The MVV is asserted in the suite by
`internal/cli/command_mvv_0026_test.go::TestReqMVV_FlowResolveDeliversTheTimeoutWhenAnEscapeeHoldsThePipes`
(rows 1-4) and `::TestReqMVV5_TheSuccessShapeRefusesExecutionFailureAndParsesNothing`
(row 5); rows 6(a) and 6(b) are
`internal/cli/cmdbind/drain_0026_test.go::TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue`
and `::TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead`.

The output below is from the shipped binary rather than the suite, run at the
`flow resolve` boundary REQ-MVV.3a requires. FX-deadline-escape is a real Go
helper spawning a real `SysProcAttr{Setsid: true}` grandchild that inherits
stdout and stderr; the model declares one command-backed reader over it.

```sh
go build -o <tmp>/intrastate ./cmd/intrastate
GO111MODULE=off go build -o <tmp>/helper helper.go   # FX-deadline-escape

# rows 1-4 — deadline shape, declared timeout = "2s"
<tmp>/intrastate flow resolve --model escape.toml \
  --artifact state=<tmp>/state.json --outcome advance --allow-commands --as=json

# row 5 — success shape, declared timeout = "20s" so the deadline never elapses
<tmp>/intrastate flow resolve --model success.toml \
  --artifact state=<tmp>/state.json --outcome advance --allow-commands --as=json

# rows 6(a) / 6(b)
go test ./internal/cli/cmdbind/ -run 'TestReqMVV6' -v
```

## REQ-MVV output

Rows 1-4 — deadline shape. The child sleeps past the 2 s deadline while the
`setsid(2)` grandchild holds both pipes:

```text
{"code":"flow-accessor-timeout","message":"the accessor `state` timed out","param":"state"}
exit: 3
elapsed: 2.884238000s
escapee pid: 74411
escapee alive after the refusal: yes
```

Row 3 assertions met: elapsed 2.884 s ≤ the C1 `bound:` committed total of
`2 s + 2·500 ms` plus its stated 100 ms tolerance (3.1 s); the executor
classified `timeout` (deadline-first); the envelope carries NO `detail`
field, which is 0025:C4's rule for the class (F6 — the held-pipe reason is
computed and then not rendered). Row 3a met: `flow resolve` itself
terminated and emitted the class in its envelope, where the original defect
report was "still blocked after 15 s". Row 4 met: the grandchild outlived
the refusal — the admitted residue — and the run killed it.

Row 5 — success shape. The child writes a whole, well-formed envelope and
exits 0 while the grandchild holds both pipes:

```text
{"code":"flow-accessor-failed","message":"the accessor `state` could not be executed","param":"state","detail":"the child's stdout, stderr stayed held past the 500 ms drain bound (exited 0); close or redirect the helper's inherited stdio"}
exit: 3
elapsed: .623834000s
```

Row 5 assertions met: `execution_failure` within `2·500 ms` of the child's
exit plus the tolerance (0.624 s); the reason names the held pipes in the
fixed order "stdout, stderr" and the direct child's "exited 0"; the
remediation is present; and the well-formed envelope on the held stdout was
NOT parsed — no tag values appear anywhere in the answer, only the refusal.

Rows 6(a) and 6(b) — the late-but-unheld controls, run in the package suite:

```text
=== RUN   TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue
--- PASS: TestReqMVV6a_ALateBurstDrainWithNoHolderReturnsTheWholeValue (2.19s)
=== RUN   TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead
--- PASS: TestReqMVV6b_APacedLateTailIsRecoveredWholeBecauseTheDeadlineIsRearmedPerRead (2.63s)
```

6(a) recovers the whole 1 MiB burst, so the deadline is a FUTURE one and not
`now`. 6(b) recovers the whole 1 MiB paced across 32 chunks with 20 ms idle
gaps — a final drain running ~605 ms, some 12·DrainGrace — so the deadline
is re-armed before each read and bounds only the idle gap. A single absolute
`now + DrainGrace` truncates 6(b) at ~15.6% and fails it.
