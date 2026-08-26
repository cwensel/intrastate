# Coverage — RDR 0005 Skill Integration CLI Contract

Phase 1 artifact. Maps every REQ in `req-list.md` (132 REQ + REQ-MVV)
to the test(s) that would FAIL if a future change broke that clause,
and flags orphans in BOTH directions.

**Red gate**: all 93 new tests FAIL against the missing implementation;
none passes tautologically. All 675 pre-existing tests still pass, and
no shipped file was modified.

## Summary

| Metric | Value |
| --- | --- |
| REQs in scope | 133 (REQ-1..REQ-132 + REQ-MVV) |
| REQs with at least one test | 133 |
| REQs with NO test (orphan REQs) | 0 |
| New tests | 93 |
| Tests mapping to no REQ (orphan tests) | 0 |
| New tests RED | 93 / 93 |
| New tests passing prematurely | 0 |
| Pre-existing tests still green | 675 / 675 |

## Test files

| File | Role |
| --- | --- |
| `internal/cli/flow_fixtures_0005_test.go` | The fixture model corpus. Every model loads under RDR 0002 except `flowInvalidModel`, which is deliberately unloadable. |
| `internal/cli/flow_harness_0005_test.go` | Invocation tables and the discriminating oracles (`requireRefusal` asserts the exact `flow-*` code and exit). |
| `internal/cli/flow_surface_0005_test.go` | Command group, verb surface, output gateway, envelope (REQ-1..REQ-13). |
| `internal/cli/flow_input_0005_test.go` | Model selection, state input, write grammar, and the input rows of the stable-code table (REQ-22..REQ-34, REQ-61..REQ-63, REQ-80..REQ-91). |
| `internal/cli/flow_next_0005_test.go` | `flow next`, the narrowed read-accessor set, and the gate-disposition mini-check (REQ-35..REQ-47, REQ-74, REQ-75). |
| `internal/cli/flow_resolve_0005_test.go` | `flow resolve`, kernel refusal mapping, and the post-selection gate site (REQ-48..REQ-57, REQ-76, REQ-92..REQ-98). |
| `internal/cli/flow_setstate_0005_test.go` | `flow set-state` / `flow read-state`, read-back-gated success, round-trip invariants (REQ-58..REQ-68, REQ-77, REQ-78, REQ-103..REQ-109). |
| `internal/cli/flow_encoder_0005_test.go` | Canonical set literal, shared encoder, and the `Finding` record on the wire (REQ-70..REQ-73, REQ-14..REQ-18, REQ-79, REQ-129). |
| `internal/cli/flow_exit_0005_test.go` | Exit-code contract, internal-group codes, ownership fences, docs (REQ-19..REQ-21, REQ-99..REQ-102, REQ-117, REQ-131). |
| `internal/cli/flow_mvv_0005_test.go` | **The REQ-MVV runner.** All four verbs, the three-hop byte-identity assertion, five scenarios x two modes, the unbound-reader pair, and `next --evaluate-gates`. |
| `internal/cli/clierr/finding_0005_test.go` | Type-level obligations on the shared record and encoder (REQ-8, REQ-11, REQ-14, REQ-16..REQ-18, REQ-71, REQ-72). |
| `internal/cli/flow_adversarial_0005_test.go` | **Phase 3b/3c.** The four defects the Phase 1 corpus's own fixtures cannot reach: FAIL-1 (every gate reported on a gate refusal), ADV-1 (guard-excluded row neither reported nor gated), ADV-2 (a reader serving only a guard's owned key), ADV-3 (`revision` renders empty). |

## REQ x test map

Every row is RED. `file` is the basename within the table above.

| REQ | Test(s) | File |
| --- | --- | --- |
| REQ-1 | `TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs` | flow_surface |
| REQ-2 | `TestReq2_FlowGroupDoesNotAbsorbLintDumpOrParse` | flow_surface |
| REQ-3 | `TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs`<br>`TestReq3_EachVerbRegistersItsNormativeFlagSpellings` | flow_surface |
| REQ-4 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-5 | `TestReq5_EveryVerbValidatesOutputModeFirst`<br>`TestReq5_VerbsSilenceCobraErrorAndUsageOutput` | flow_surface |
| REQ-6 | `TestReq6_JSONSuccessIsExactlyOneOKEnvelopeWithVerbSpecificData` | flow_surface |
| REQ-7 | `TestReq7And10_TextSuccessEmitsHumanOutputFromTheSameResult` | flow_surface |
| REQ-8 | `TestReq8And16_CLIErrorFindingsIsOmitEmpty`<br>`TestReq8And16_FindingsIsOmitEmptyOnTheCLIErrorEnvelope` | finding<br>flow_encoder |
| REQ-9 | `TestReq9And10_TextPayloadsRouteThroughTheGatewayNotDirectPrinting` | flow_surface |
| REQ-10 | `TestReq7And10_TextSuccessEmitsHumanOutputFromTheSameResult`<br>`TestReq9And10_TextPayloadsRouteThroughTheGatewayNotDirectPrinting` | flow_surface |
| REQ-11 | `TestReq11And120_TextRenderingCarriesEveryJSONScalarValue`<br>`TestReq11And129_TextModeEnumeratesEveryFindingOnePerLine`<br>`TestReq11And17And18_TextRendersThisRDRsIdentityFieldsOnePerLine` | finding<br>flow_encoder<br>flow_surface |
| REQ-12 | `TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag` | flow_surface |
| REQ-13 | `TestReq13And19_NoNewTerminalEnvelopeTypeAndNoNewExitGroup` | flow_surface |
| REQ-14 | `TestReq14And15And18_FindingsAreFlatAndUnsetOptionalFieldsAreAbsent`<br>`TestReq14And16_FindingCarriesTheThirteenNamedFieldsFlat`<br>`TestReq14_UnsetOptionalFindingFieldsAreAbsentFromTheJSON` | finding<br>flow_encoder |
| REQ-15 | `TestReq14And15And18_FindingsAreFlatAndUnsetOptionalFieldsAreAbsent`<br>`TestReq96_GuardUnevaluableCarriesRowsAndFlatAtomFields` | flow_encoder<br>flow_resolve |
| REQ-16 | `TestReq14And16_FindingCarriesTheThirteenNamedFieldsFlat`<br>`TestReq8And16_CLIErrorFindingsIsOmitEmpty`<br>`TestReq8And16_FindingsIsOmitEmptyOnTheCLIErrorEnvelope` | finding<br>flow_encoder |
| REQ-17 | `TestReq11And17And18_TextRendersThisRDRsIdentityFieldsOnePerLine`<br>`TestReq17_EveryFindingMessageStandsAloneWithoutItsStructuredFields` | finding<br>flow_encoder |
| REQ-18 | `TestReq11And17And18_TextRendersThisRDRsIdentityFieldsOnePerLine`<br>`TestReq14And15And18_FindingsAreFlatAndUnsetOptionalFieldsAreAbsent` | finding<br>flow_encoder |
| REQ-19 | `TestReq13And19_NoNewTerminalEnvelopeTypeAndNoNewExitGroup`<br>`TestReq19And20_TheExit3PopulationIsExactlyTheFiveNamedClasses`<br>`TestReq19_AnExit3RefusalIsReRunnableUnchangedToTheSameDisposition` | flow_exit<br>flow_surface |
| REQ-20 | `TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied`<br>`TestReq19And20_TheExit3PopulationIsExactlyTheFiveNamedClasses`<br>`TestReq20_EveryAccessorRefusalClassHasACLIMapping` | flow_exit<br>flow_setstate |
| REQ-21 | `TestReq21_UnavailableAccessorExits3ButIndeterminateGateExits2` | flow_exit |
| REQ-22 | `TestReq22And90_ModelSelectionArityAndNonResolutionShareOneCode` | flow_input |
| REQ-23 | `TestReq23_ModelFileIOBelongsToTheCLINotTheLoader` | flow_input |
| REQ-24 | `TestReq24And91_LoadCategoriesRideOneCodeWithPerCategoryFindings` | flow_input |
| REQ-25 | `TestReq25_RevisionIsNeverCLIDerivedFromPathOrContent` | flow_input |
| REQ-26 | `TestReq26And114_TagEntersAsObservedAndNeverSatisfiesAnOwnedDependency` | flow_input |
| REQ-27 | `TestReq27And82And83_OwnedAndReservedTagsAreRefusedBeforeAnyAccessorRuns`<br>`TestReq34_SetStateNeverTreatsTagValuesAsWrites` | flow_input |
| REQ-28 | `TestReq28And81_DuplicateTagNameIsItsOwnCode` | flow_input |
| REQ-29 | `TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid` | flow_input |
| REQ-30 | `TestReq30And57_ResolveAssemblesOwnedStateFromTheDeclaredReaders` | flow_resolve |
| REQ-31 | `TestReq31And32And88_ArtifactLocationComesOnlyFromExplicitBindings` | flow_input |
| REQ-32 | `TestReq31And32And88_ArtifactLocationComesOnlyFromExplicitBindings` | flow_input |
| REQ-33 | `TestReq33And89_MissingBindingForAnInvokedRoleCarriesTheRoleAsParam` | flow_input |
| REQ-34 | `TestReq34_SetStateNeverTreatsTagValuesAsWrites` | flow_input |
| REQ-35 | `TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet`<br>`TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds` | flow_mvv<br>flow_next |
| REQ-36 | `TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet`<br>`TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds` | flow_mvv<br>flow_next |
| REQ-37 | `TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet`<br>`TestReq37_ReadStateRunsEveryDeclaredReaderAndRefusesTheMissingBinding` | flow_mvv<br>flow_next |
| REQ-38 | `TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds` | flow_next |
| REQ-39 | `TestReq39And59_NextResolveAndReadStateNeverMutateTheArtifact` | flow_next |
| REQ-40 | `TestReq40And74_NextPayloadCarriesTheAlphabetAndCandidateSummaries` | flow_next |
| REQ-41 | `TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults` | flow_next |
| REQ-42 | `TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny`<br>`TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate` | flow_mvv<br>flow_next |
| REQ-43 | `TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny`<br>`TestReq43And46_ADenyUnderNextIsReportedAndStillExitsZero` | flow_mvv<br>flow_next |
| REQ-44 | `TestReq44_NextInventsNoGuardFactAbsentFromEveryChannel` | flow_next |
| REQ-45 | `TestReq45And110And111_TheSameRequestProducesTheSameResultEveryRun` | flow_next |
| REQ-46 | `TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny`<br>`TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults`<br>`TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate`<br>`TestReq43And46_ADenyUnderNextIsReportedAndStillExitsZero` | flow_mvv<br>flow_next |
| REQ-47 | `TestReq47And55_AGateThatCouldNotBeConsultedIsExit3AndNeverADeny`<br>`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal` | flow_adversarial<br>flow_resolve |
| REQ-48 | `TestReq48And53_ResolveReturnsExactlyOnePlanOrExactlyOneRefusal` | flow_resolve |
| REQ-49 | `TestReq49_TheFiveKernelKindsRaiseFiveDistinctCodes`<br>`TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-50 | `TestReq50And76_GatesRunAfterSelectionAndEveryResultIsReported`<br>`TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans`<br>`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal` | flow_adversarial<br>flow_resolve |
| REQ-51 | `TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans`<br>`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal` | flow_adversarial<br>flow_resolve |
| REQ-52 | `TestReq52_AnEscapedPlanIsASuccessCarryingEscapedTrueAndItsClass`<br>`TestReq52_AnOrdinaryPlanIsNotMarkedEscaped` | flow_resolve |
| REQ-53 | `TestReq48And53_ResolveReturnsExactlyOnePlanOrExactlyOneRefusal`<br>`TestReq53And113_ResolveEmitsAPlanAndInitiatesNoWorkItself` | flow_resolve |
| REQ-54 | `TestReq54_AbsentOutcomeIsCLISideAndOutOfAlphabetIsKernelSide` | flow_resolve |
| REQ-55 | `TestReq47And55_AGateThatCouldNotBeConsultedIsExit3AndNeverADeny` | flow_resolve |
| REQ-56 | `TestReq56And119_EveryRefusalFamilyKeepsItsOwnIdentity` | flow_resolve |
| REQ-57 | `TestReq30And57_ResolveAssemblesOwnedStateFromTheDeclaredReaders`<br>`TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults` | flow_next<br>flow_resolve |
| REQ-58 | `TestReq58And77_ReadStateReportsDeclaredKeysBesideTheTagsReturned` | flow_setstate |
| REQ-59 | `TestReq39And59_NextResolveAndReadStateNeverMutateTheArtifact`<br>`TestReq59And115_ReadStateNeverCoercesAGateVerdictIntoATagValue` | flow_next |
| REQ-60 | `TestReq60And78_SetStatePayloadCarriesItsFullDataMinimum` | flow_setstate |
| REQ-61 | `TestReq102_TheThreeInternalCodesAreUnreachableFromTheUserPath`<br>`TestReq61And86And87_UnboundWriteAndClearKeysAreRefusedBeforeAccessors` | flow_exit<br>flow_input |
| REQ-62 | `TestReq62And84_ClearSentinelIsUnauthorableAndHintsAtTheFlag` | flow_input |
| REQ-63 | `TestReq63And85_DuplicateKeyAcrossEitherFlagIsOneCode` | flow_input |
| REQ-64 | `TestReq64_SetStateNeverRunsGateAccessors` | flow_setstate |
| REQ-65 | `TestReq65_SuccessImpliesTheValueIsActuallyReadableAfterwards` | flow_setstate |
| REQ-66 | `TestReq66And107_EmptySetAndClearedKeyStayDistinctThroughReadBack` | flow_setstate |
| REQ-67 | `TestReq68_SetStateAcceptsAnUnlinkedRequestWithNoPriorResolve` | flow_setstate |
| REQ-68 | `TestReq117_TheResolveToSetStateWindowIsAStatedNonGuarantee`<br>`TestReq68_SetStateAcceptsAnUnlinkedRequestWithNoPriorResolve` | flow_exit<br>flow_setstate |
| REQ-69 | `TestReq69_NoPlanFlagShipsOnAnyVerb` | flow_input |
| REQ-70 | `TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical`<br>`TestReq70And73_SetValuesRenderAngleAndAmpersandAsThemselves`<br>`TestReq70_SetValuesCrossingInTagObeyTheSameCanonicalForm` | flow_encoder<br>flow_mvv |
| REQ-71 | `TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical`<br>`TestReq71And72_BothEmitSitesUseTheSameNonEscapingEncoder`<br>`TestReq71And72_EmitJSONRendersAngleAndAmpersandAsThemselves`<br>`TestReq71_TheCanonicalLiteralSurvivesTheJSONEmitSiteByteIdentically` | finding<br>flow_encoder<br>flow_mvv |
| REQ-72 | `TestReq71And72_BothEmitSitesUseTheSameNonEscapingEncoder`<br>`TestReq71And72_EmitJSONRendersAngleAndAmpersandAsThemselves`<br>`TestReq71_TheCanonicalLiteralSurvivesTheJSONEmitSiteByteIdentically` | finding<br>flow_encoder |
| REQ-73 | `TestReq70And73_SetValuesRenderAngleAndAmpersandAsThemselves` | flow_encoder |
| REQ-74 | `TestReq40And74_NextPayloadCarriesTheAlphabetAndCandidateSummaries` | flow_next |
| REQ-75 | `TestReq75And118_CandidateSummariesAreReadFromNormalizedModelData` | flow_next |
| REQ-76 | `TestReq50And76_GatesRunAfterSelectionAndEveryResultIsReported`<br>`TestReq76_ResolvePayloadCarriesItsFullDataMinimum` | flow_resolve |
| REQ-77 | `TestReq58And77_ReadStateReportsDeclaredKeysBesideTheTagsReturned` | flow_setstate |
| REQ-78 | `TestReq60And78_SetStatePayloadCarriesItsFullDataMinimum` | flow_setstate |
| REQ-79 | `TestReq79_CarrierChoiceFollowsTheFailuresShape` | flow_encoder |
| REQ-80 | `TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid` | flow_input |
| REQ-81 | `TestReq28And81_DuplicateTagNameIsItsOwnCode` | flow_input |
| REQ-82 | `TestReq27And82And83_OwnedAndReservedTagsAreRefusedBeforeAnyAccessorRuns` | flow_input |
| REQ-83 | `TestReq27And82And83_OwnedAndReservedTagsAreRefusedBeforeAnyAccessorRuns` | flow_input |
| REQ-84 | `TestReq62And84_ClearSentinelIsUnauthorableAndHintsAtTheFlag` | flow_input |
| REQ-85 | `TestReq63And85_DuplicateKeyAcrossEitherFlagIsOneCode` | flow_input |
| REQ-86 | `TestReq61And86And87_UnboundWriteAndClearKeysAreRefusedBeforeAccessors` | flow_input |
| REQ-87 | `TestReq61And86And87_UnboundWriteAndClearKeysAreRefusedBeforeAccessors` | flow_input |
| REQ-88 | `TestReq31And32And88_ArtifactLocationComesOnlyFromExplicitBindings` | flow_input |
| REQ-89 | `TestReq33And89_MissingBindingForAnInvokedRoleCarriesTheRoleAsParam` | flow_input |
| REQ-90 | `TestReq22And90_ModelSelectionArityAndNonResolutionShareOneCode` | flow_input |
| REQ-91 | `TestReq24And91_LoadCategoriesRideOneCodeWithPerCategoryFindings` | flow_input |
| REQ-92 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-93 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-94 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-95 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne` | flow_resolve |
| REQ-96 | `TestReq4And49_EveryKernelKindMapsToItsMirroredCodeOneToOne`<br>`TestReq96_GuardUnevaluableCarriesRowsAndFlatAtomFields` | flow_resolve |
| REQ-97 | `TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans` | flow_resolve |
| REQ-98 | `TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans` | flow_resolve |
| REQ-99 | `TestReq99And100And101_AccessorRefusalsCarryTheAccessorIdAsParam` | flow_exit |
| REQ-100 | `TestReq99And100And101_AccessorRefusalsCarryTheAccessorIdAsParam` | flow_exit |
| REQ-101 | `TestReq99And100And101_AccessorRefusalsCarryTheAccessorIdAsParam` | flow_exit |
| REQ-102 | `TestReq102_TheThreeInternalCodesAreUnreachableFromTheUserPath` | flow_exit |
| REQ-103 | `TestReq65_SuccessImpliesTheValueIsActuallyReadableAfterwards` | flow_setstate |
| REQ-104 | `TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied` | flow_setstate |
| REQ-105 | `TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied` | flow_setstate |
| REQ-106 | `TestReq106_NoneOfTheThreeSilentFailureShapesBecomesASuccess` | flow_resolve |
| REQ-107 | `TestReq107_SetStateThenReadStateReturnsThePlannedValues`<br>`TestReq66And107_EmptySetAndClearedKeyStayDistinctThroughReadBack` | flow_setstate |
| REQ-108 | `TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical`<br>`TestReq108_PlanValuesAreAcceptedVerbatimBySetStatesGrammar` | flow_mvv<br>flow_setstate |
| REQ-109 | `TestReq109_AnOwnedTagReadByReadStateCannotBeHandedBackThroughTag` | flow_setstate |
| REQ-110 | `TestReq19_AnExit3RefusalIsReRunnableUnchangedToTheSameDisposition`<br>`TestReq45And110And111_TheSameRequestProducesTheSameResultEveryRun` | flow_exit<br>flow_next |
| REQ-111 | `TestReq45And110And111_TheSameRequestProducesTheSameResultEveryRun` | flow_next |
| REQ-112 | `TestReq112_NoNewThirdPartyDependencyIsIntroduced` | flow_input |
| REQ-113 | `TestReq113And121_TheCLIMapsTheKernelsClosedKindSetAndAddsNoKinds`<br>`TestReq53And113_ResolveEmitsAPlanAndInitiatesNoWorkItself` | flow_next<br>flow_resolve |
| REQ-114 | `TestReq26And114_TagEntersAsObservedAndNeverSatisfiesAnOwnedDependency` | flow_input |
| REQ-115 | `TestReq59And115_ReadStateNeverCoercesAGateVerdictIntoATagValue` | flow_next |
| REQ-116 | `TestReq52_AnEscapedPlanIsASuccessCarryingEscapedTrueAndItsClass`<br>`TestReq52_AnOrdinaryPlanIsNotMarkedEscaped` | flow_resolve |
| REQ-117 | `TestReq117_TheResolveToSetStateWindowIsAStatedNonGuarantee` | flow_exit |
| REQ-118 | `TestReq75And118_CandidateSummariesAreReadFromNormalizedModelData` | flow_next |
| REQ-119 | `TestReq56And119_EveryRefusalFamilyKeepsItsOwnIdentity` | flow_resolve |
| REQ-120 | `TestMVV_FiveScenariosInBothModesAgreeOnExitAndIdentity`<br>`TestReq11And120_TextRenderingCarriesEveryJSONScalarValue` | flow_mvv<br>flow_surface |
| REQ-121 | `TestReq113And121_TheCLIMapsTheKernelsClosedKindSetAndAddsNoKinds` | flow_next |
| REQ-122 | `TestMVV_AllFourVerbsProveOutOverOneFixtureBackedFlow` | flow_mvv |
| REQ-123 | `TestReq41And46And57_WithoutTheFlagGateIdsAreUnresolvedFactsNotResults` | flow_next |
| REQ-124 | `TestReq50And76_GatesRunAfterSelectionAndEveryResultIsReported` | flow_resolve |
| REQ-125 | `TestReq54_AbsentOutcomeIsCLISideAndOutOfAlphabetIsKernelSide` | flow_resolve |
| REQ-126 | `TestReq51And97And98_GateDenyAndIndeterminateAreTypedRefusalsNeverPlans` | flow_resolve |
| REQ-127 | `TestReq66And107_EmptySetAndClearedKeyStayDistinctThroughReadBack`<br>`TestReq70And73_SetValuesRenderAngleAndAmpersandAsThemselves` | flow_encoder<br>flow_setstate |
| REQ-128 | `TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied` | flow_setstate |
| REQ-129 | `TestReq11And129_TextModeEnumeratesEveryFindingOnePerLine`<br>`TestReq129_ThreeProducersShareOneRecordAndPopulateDisjointFields`<br>`TestReq14_UnsetOptionalFindingFieldsAreAbsentFromTheJSON` | finding<br>flow_encoder |
| REQ-130 | `TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds`<br>`TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate` | flow_next |
| REQ-131 | `TestReq131_TheOutputContractDocumentsFindingsAndTheExit3Rule` | flow_exit |
| REQ-132 | `TestReq132_TheIllustrativeInvocationShapesParseUnderTheShippedGrammar` | flow_mvv |
| REQ-MVV | `TestMVV_AllFourVerbsProveOutOverOneFixtureBackedFlow`<br>`TestMVV_FiveScenariosInBothModesAgreeOnExitAndIdentity`<br>`TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny`<br>`TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical`<br>`TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet`<br>`TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds`<br>`TestReq37_ReadStateRunsEveryDeclaredReaderAndRefusesTheMissingBinding`<br>`TestReq42And46_OnlyReportedCandidatesGatesRunAndResultsRideTheCandidate`<br>`TestReq43And46_ADenyUnderNextIsReportedAndStillExitsZero` | flow_mvv<br>flow_next |

## Orphans

**REQs with no test: none.** All 132 REQ plus REQ-MVV map to at least
one test.

**Tests mapping to no REQ: none.** Every test opens with at least one
`// REQ-N: "<quote>"` citation and one of the five labels.

## REQ-MVV — end-to-end run (recorded output)

**Refreshed in Phase 3c** after the seven-defect fixup. Re-run against the
`flowMVVModel` fixture through the built binary
(`go build -o intrastate ./cmd/intrastate`), paths shortened to `<m>` for the
model and `<a>` for the artifact. Every obligation `0005:MVV` names is
discharged below.

The only change from the Phase 2 recording is `revision`, which now renders
EMPTY at every payload site. That is DEV-7's decision implementing REQ-25:
RDR 0002's `[model]` block declares no revision field, so every model
"declares none" and the clause's consequent — empty, never a CLI-invented
value — governs every payload. Section 9 is new and records the FAIL-1 fix.

**1. `set-state` — scalar write, set write carrying `<` and `&`, and a clear**

```
$ intrastate flow set-state --model <m> --artifact state=<a> \
    --write status=final --write 'labels=["a<b","x&y"]' --clear stale --as=json
{"type":"ok","data":{"model":"<m>","revision":"",
 "artifacts":{"state":"<a>"},"writers":["state"],
 "writes":{"labels":"[\"a<b\",\"x&y\"]","status":"final"},
 "clear":["stale"],
 "owned":{"labels":"[\"a<b\",\"x&y\"]","status":"final"}}}
EXIT=0
```

`<` and `&` serialize as THEMSELVES at every carrier, and the cleared
`stale` is ABSENT from the read-back-confirmed `owned` — not reported as
holding the sentinel or an empty string.

**2. `read-state` — per-reader declared keys beside the tags returned**

```
$ intrastate flow read-state --model <m> --artifact state=<a> --artifact orphan=<a> --as=json
{"type":"ok","data":{"model":"<m>","revision":"",
 "artifacts":{"orphan":"<a>","state":"<a>"},
 "readers":[{"id":"orphan","keys":["note"],"tags":{}},
            {"id":"state","keys":["status","labels","stale"],
             "tags":{"labels":"[\"a<b\",\"x&y\"]","status":"final"}}]}}
EXIT=0
```

`read-state` runs EVERY declared reader (`orphan` included). `stale`
appears in `state.keys` while carrying no entry in `tags`, so "absent from
the artifact" is distinguishable from "not requested" — `note` is in
neither for this reader.

**3. `next` — the alphabet, with gate ids as unresolved facts**

```
$ intrastate flow next --model <m> --artifact state=<a> --tag profile=mid --as=json
{"type":"ok","data":{"model":"<m>","revision":"",
 "observed":{"profile":"mid"},
 "owned":{"labels":"[\"a<b\",\"x&y\"]","status":"final"},
 "readers":["state"],
 "outcomes":["advance","hold","bail"],
 "candidates":[
   {"rule":"advance-draft","outcome":"advance","required":["stale","status"],
    "unresolved":["stale","approval"],"next":{"stale":"<clear>","status":"final"},
    "writes":{"status":"final"},"clear":["stale"]},
   {"rule":"hold-draft","outcome":"hold","required":["status"],
    "unresolved":[],"next":{"status":"draft"},
    "writes":{"status":"draft"},"clear":[]}]}}
EXIT=0
```

The gate id `approval` rides `unresolved` and NO gate accessor ran.
`readers[]` omits `orphan` — the narrowed invoked set (REQ-35). `clear[]`
splits `stale` out of the preview `writes`, so the `<clear>` sentinel never
appears as a write target. All three declared outcomes are reported
(DEV-4).

**4. `resolve` — one outcome to one plan, one gate on the selected row**

```
$ intrastate flow resolve --model <m> --artifact state=<b> --outcome advance --as=json
{"type":"ok","data":{"model":"<m>","revision":"","observed":{},
 "owned":{"stale":"x","status":"draft"},"readers":["state"],
 "outcome":"advance","rule":"advance-draft",
 "gates":[{"id":"approval","result":"allow"}],
 "next":{"stale":"<clear>","status":"final"},
 "writes":{"status":"final"},"clear":["stale"],"escaped":false}}
EXIT=0
```

EXACTLY one gate — the selected row's — ran and its result is reported.
`readers[]` again omits `orphan`, whose role is unbound throughout.

**5. The unbound-reader pair — the discriminating assertion**

The SAME `--artifact state=<b>` binding succeeds on `next` / `resolve`
(above) and refuses on `read-state`:

```
$ intrastate flow read-state --model <m> --artifact state=<b> --as=json
{"code":"flow-artifact-missing",
 "message":"the artifact role `orphan` that the read accessor `orphan` needs has no --artifact binding",
 "param":"orphan"}
EXIT=2
```

No "run every declared reader" implementation can produce this pair.

**6. One kernel refusal, in both modes**

```
$ intrastate flow resolve --model <m> --artifact state=<b> --outcome bail --as=json
{"code":"flow-no-match","message":"no rule matches the recognized outcome `bail` over the assembled state",
 "findings":[{"code":"flow-no-match","message":"no rule in the model responds to the recognized outcome `bail`"}]}
EXIT=2

$ intrastate flow resolve --model <m> --artifact state=<b> --outcome bail     # --as=text
error: flow-no-match: no rule matches the recognized outcome `bail` over the assembled state
  flow-no-match: no rule in the model responds to the recognized outcome `bail`
EXIT=2
```

Both modes agree on exit and refusal identity; `findings` render
one-per-line under the message in text mode.

**7. One exit-3 accessor failure** (`flowReadBackFailModel`)

```
$ intrastate flow set-state --model <readbackfail> --artifact state=<c> --write status=final --as=json
{"code":"flow-write-readback-incomplete",
 "message":"the post-mutation read-back for the write accessor `state` did not complete",
 "param":"state",
 "detail":"the write command already ran, so the mutation may have been applied and was not verified; inspect the artifact before retrying"}
EXIT=3

$ cat <c>
{"\u0000flow.readback-unreachable":"1","status":"final"}
```

Exit 3 with the "may have been applied" detail, and the mutation genuinely
LANDED (`status: final`) — the applied-but-unverified sense `0004:C14`
requires, not a mismatch and not a write that did not occur.

**8. `resolve` in `--as=text`** — the same result, human-scannable, with
every JSON scalar present:

```
clear[0]: stale
escaped: false
gates[0].id: approval
gates[0].result: allow
model: <m>
next.stale: <clear>
next.status: final
observed: (none)
outcome: advance
owned.stale: x
owned.status: draft
readers[0]: state
revision:
rule: advance-draft
writes.status: final
```

**9. Every gate on the selected row is reported on a gate refusal**
(Phase 3c — FAIL-1; fixture `advGateReportModel`, whose `deny-row` carries an
ALLOWING gate `permits` beside the denying `refuses`)

```
$ intrastate flow resolve --model <g> --artifact state=<d> --outcome deny-path --as=json
{"code":"flow-gate-denied","message":"a gate on the selected rule denied the transition",
 "findings":[{"code":"flow-gate-denied","message":"the gate `permits` answered allow","param":"permits"},
             {"code":"flow-gate-denied","message":"the gate `refuses` denied the transition: the gate flow.gate.deny denied the transition","param":"refuses"}]}
EXIT=2

$ intrastate flow resolve --model <g> --artifact state=<d> --outcome deny-path   # --as=text
error: flow-gate-denied: a gate on the selected rule denied the transition
  flow-gate-denied: the gate `permits` answered allow (param="permits")
  flow-gate-denied: the gate `refuses` denied the transition: the gate flow.gate.deny denied the transition (param="refuses")
EXIT=2
```

`findings[]` carries ONE ENTRY PER GATE ON THE ROW — the allowing gate
included — so "every gate result MUST be reported" (REQ-50) and "no gate
result is ever dropped silently" (REQ-47) hold on the REFUSAL path exactly as
they do on the success path's `gates[]`. Precedence is unchanged: the envelope
`code` is still `flow-gate-denied`, deny over allow and indeterminate.

`flow next --evaluate-gates` over the SAME row reports the same two gates, so
the two verbs agree about what the gates answered:

```
$ intrastate flow next --model <g> --artifact state=<d> --evaluate-gates --as=json
... "gates":[{"id":"permits","result":"allow"},
             {"id":"refuses","result":"deny","reason":"the gate flow.gate.deny denied the transition"}] ...
EXIT=0
```

## Phase 2 result

| Metric | Value |
| --- | --- |
| Suite subtests passing | 241 |
| Tests failing | 3 (all recorded oracle defects: DEV-2, DEV-3, DEV-5) |
| Pre-existing suites | green (accessor, graphlint, guard, resolve, table) |
| `go vet` / `gofmt` / `golangci-lint` | clean (0 issues) |

## Phase 3c result

| Metric | Value |
| --- | --- |
| `internal/cli` + `internal/cli/clierr` tests passing | 291 |
| Tests failing | **0** — full `go test ./...` GREEN |
| Predecessor suites (0001, 0002, 0003, 0004, 0006) | green; none modified or weakened |
| `go vet` / `gofmt` / `golangci-lint run` | clean (0 issues) |
| Defects closed | FAIL-1, ADV-1, ADV-2, ADV-3, DEV-2, DEV-3, DEV-5 |
| Deviations opened | DEV-8 (SPEC-UNDER, mechanical translation — A-4 superseded) |
| Deviations still `needs author decision` | **0** |

Regression tests added:
`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal`
(`internal/cli/flow_adversarial_0005_test.go`). ADV-1, ADV-2, and ADV-3 already
carried Phase 3b tests, which are now green. Each of the three corrected
oracles (DEV-2, DEV-3, DEV-5) was re-verified to still DISCRIMINATE the
behaviour it names by injecting the defect the REQ forbids.

## Notes for Phase 2

- **DEV-1 is honoured as reading (a).** `TestReq35And36And38_NextAndResolveSkipAReaderNoCandidateRowNeeds` asserts the narrowing is the union over ALL model rows: `read.orphan` serves only `note`, which no rule's `RequiresOwned` names, so it is skipped under either reading and the fixture stays discriminating. No test asserts a reader must NOT run for a reason that depends on `--tag` content, so DEV-1's escalation trigger is not fired.
- **Two `clierr` fields are the hard blocker.** `Finding.Param` and `Finding.Locator` (A-7) do not exist. The tests reach them REFLECTIVELY and through JSON rather than by struct literal, deliberately: a struct literal would be a package-wide COMPILE break that also takes RDR 0006's shipped `clierr` suite down. As written, the absence is a readable failure naming the field.
- **`CLIError.Findings` must move to `omitempty` (A-6).** `TestReq8And16_CLIErrorFindingsIsOmitEmpty` shows the current envelope emits `"findings":null` on a refusal that names none. RDR 0006's shipped assertions all run against non-empty lists, so they stay green; that was verified, not assumed.
- **`Finding.Model` and `Finding.Severity` are currently non-`omitempty`.** REQ-14 requires the optional fields be `omitempty`; the shipped record marks these two required. Phase 2 must reconcile this against `0006:REQ-94` the same way A-6 reconciles `CLIError.Findings`.
- **No production `accessor.Binding` implementation ships.** `internal/accessor` defines the typed binding seam but the only implementations are RDR 0004's test fixtures. Phase 3 must supply the binding the CLI uses to reach a caller-owned artifact. No test in this suite asserts the artifact's on-disk FORMAT — every state oracle goes `set-state` -> `read-state` through the CLI — so the implementation is free to choose it.
- **`revision` has no declared carrier in RDR 0002's schema.** `[model]` accepts only `id`, `version`, `description`, `metadata`; `Model.KernelTable()` sets `Revision: m.ID`. `TestReq25_RevisionIsNeverCLIDerivedFromPathOrContent` therefore asserts the NEGATIVE REQ-25 fixes (never a path fragment, never a digest) rather than presuming a field. **Phase 3c settled the positive leg per DEV-7**: since no model CAN declare a revision, every model "declares none" and `flowRequest.revision()` renders EMPTY. `KernelTable().Revision` is left as-is for the kernel's own table-identity use and is simply no longer forwarded as the model revision; `TestAdv3_AModelDeclaringNoRevisionRendersRevisionEmpty` is the oracle.
- **Gate verdicts are fixture-driven, not path-driven.** `flowGateDenyModel` and `flowGateFailModel` name gate ACCESSOR IDS (`deny`, `indeterminate`, `unreachable`); how a fixture gate is made to return a given verdict is Phase 3's binding concern. The oracles assert the resulting refusal identity and exit, which is what the contract fixes.

