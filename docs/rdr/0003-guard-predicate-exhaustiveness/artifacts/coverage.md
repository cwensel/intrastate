# Coverage — RDR 0003 Guard Predicate Exhaustiveness

Phase 1 artifact. Every REQ in
[`req-list.md`](req-list.md) maps to exactly one test, and every test in
`internal/guard/` cites exactly one REQ. Both orphan directions are
checked mechanically below and both are empty.

**Status: RED.** No implementation has landed. `internal/guard/skeleton.go`
declares signatures only so the suite compiles and each test fails on its
own assertion — naming the REQ its header quotes — rather than the whole
package failing with one undefined-symbol error that attributes to no
clause. All 142 top-level tests fail. See `deviations.md` for the
cluster-gate entries and Phase 2's disposition of them.

## Totals

| Quantity | Count |
| --- | --- |
| REQs in `req-list.md` (141 + `REQ-MVV`) | 142 |
| Tests in `internal/guard/` | 142 |
| REQs with no test | 0 |
| Tests citing no REQ | 0 |
| Top-level tests currently red | 142 |
| Pre-existing tests broken | 0 |

## REQ × test

| REQ | Test | Label | File |
| --- | --- | --- | --- |
| REQ-1 | `TestReq1_GuardIsASymbolicAtomOverADeclaredTag` | HAPPY PATH | `internal/guard/guard_grammar_0003_test.go` |
| REQ-2 | `TestReq2_OperatorVocabularyIsClosedAndTyped` | BOUNDARY | `internal/guard/guard_grammar_0003_test.go` |
| REQ-3 | `TestReq3_EachOperatorDeclaresAcceptedKindsAndLiteralsAreKindChecked` | HAPPY PATH | `internal/guard/guard_grammar_0003_test.go` |
| REQ-4 | `TestReq4_OperatorKindMatrixIsExactlyAsPublished` | BOUNDARY | `internal/guard/guard_grammar_0003_test.go` |
| REQ-5 | `TestReq5_ExistsOverAnAlwaysPresentKeyIsVacuousNotRejected` | DOMAIN EDGE | `internal/guard/guard_grammar_0003_test.go` |
| REQ-6 | `TestReq6_AtomWellFormednessIsKeyOperatorAndLiteral` | ADVERSARIAL | `internal/guard/guard_grammar_0003_test.go` |
| REQ-7 | `TestReq7_SourceIdentityIsNotAnAtomFieldButJoinsTheIdentityTuple` | BOUNDARY | `internal/guard/guard_grammar_0003_test.go` |
| REQ-8 | `TestReq8_AtomIdentityIsTheSixFieldTupleNeverThePositionalIndex` | ADVERSARIAL | `internal/guard/guard_grammar_0003_test.go` |
| REQ-9 | `TestReq9_SemanticEqualityIsTagOperatorLiteral` | HAPPY PATH | `internal/guard/guard_grammar_0003_test.go` |
| REQ-10 | `TestReq10_DiagnosticsUseIdentityAndDomainComputationUsesSemanticEquality` | DOMAIN EDGE | `internal/guard/guard_grammar_0003_test.go` |
| REQ-11 | `TestReq11_ByteIdenticalGuardsInTwoRulesOverlapRatherThanDeduplicate` | ADVERSARIAL | `internal/guard/guard_grammar_0003_test.go` |
| REQ-12 | `TestReq12_OperatorSemanticsAreAsStated` | HAPPY PATH | `internal/guard/guard_grammar_0003_test.go` |
| REQ-13 | `TestReq13_DeclarationCarriesAKindAndTheFourOptionalFields` | HAPPY PATH | `internal/guard/guard_declaration_0003_test.go` |
| REQ-14 | `TestReq14_ValueKindsAreExactlyFiveTokens` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-15 | `TestReq15_ScalarIsOpaqueAndNeverBearsAnExhaustivenessClaim` | DOMAIN EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-16 | `TestReq16_ThisRDROwnsMeaningNotAuthoringOrCarriage` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-17 | `TestReq17_FiniteDomainIsDeclarableForEveryClaimableKind` | HAPPY PATH | `internal/guard/guard_declaration_0003_test.go` |
| REQ-18 | `TestReq18_IntBoundEndpointsAreInclusive` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-19 | `TestReq19_NoFiniteDomainIsWellFormedButUnprovable` | DOMAIN EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-20 | `TestReq20_OmittedOptionalityMarkerDeclaresTheKeyOptional` | INPUT EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-21 | `TestReq21_ConformanceRequiresEveryAlwaysPresentKey` | DOMAIN EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-22 | `TestReq22_PresenceIsReadFromTheDeclarationNotAView` | ADVERSARIAL | `internal/guard/guard_declaration_0003_test.go` |
| REQ-23 | `TestReq23_ConformanceIsBothConjunctsAndScopesEveryGreenClaim` | DOMAIN EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-24 | `TestReq24_NoClaimIsMadeAboutANonConformingView` | ADVERSARIAL | `internal/guard/guard_declaration_0003_test.go` |
| REQ-25 | `TestReq25_SingleValuedMarkerIsStatableAndBannedOnSet` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-26 | `TestReq26_SingleValuedContributesDomainSizeUnmarkedContributesPowerset` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-27 | `TestReq27_SingleValuedIsNeverInferredFromNameSpellingOrFixture` | ADVERSARIAL | `internal/guard/guard_declaration_0003_test.go` |
| REQ-28 | `TestReq28_DomainMustAgreeWithKindAndDisagreementIsRejected` | ADVERSARIAL | `internal/guard/guard_declaration_0003_test.go` |
| REQ-29 | `TestReq29_KindFieldAgreementTableIsEnforcedPerCell` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-30 | `TestReq30_LiteralOutsideDeclaredDomainIsRejectedAtLoad` | ADVERSARIAL | `internal/guard/guard_declaration_0003_test.go` |
| REQ-31 | `TestReq31_DeclarationErrorsAndLiteralErrorsAreRejectedAtDifferentStages` | DOMAIN EDGE | `internal/guard/guard_declaration_0003_test.go` |
| REQ-32 | `TestReq32_TheTwoRejectionsMapOntoRDR0002sTwoLoadCategories` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-33 | `TestReq33_TheDeclarationModelIsExactlyFiveFields` | BOUNDARY | `internal/guard/guard_declaration_0003_test.go` |
| REQ-34 | `TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView` | BOUNDARY | `internal/guard/guard_evaluator_0003_test.go` |
| REQ-35 | `TestReq35_ExistenceOverAbsentIsDecidedAndValueOverAbsentIsUnevaluable` | DOMAIN EDGE | `internal/guard/guard_evaluator_0003_test.go` |
| REQ-36 | `TestReq36_AbsenceIsExpressedWithAnExistenceAtomNotAMissingValueAtom` | DOMAIN EDGE | `internal/guard/guard_evaluator_0003_test.go` |
| REQ-37 | `TestReq37_ValueVerdictIsThisRDRsAndPresenceCombinationIsTheKernels` | BOUNDARY | `internal/guard/guard_evaluator_0003_test.go` |
| REQ-38 | `TestReq38_PolarityIsByBlockAndMatchingIsOrderIndependent` | HAPPY PATH | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-39 | `TestReq39_AcceptedAssignmentsAreAllIntersectionMinusTheUnlessConjunction` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-40 | `TestReq40_RowQualifiesOnAllTrueAndUnlessNotFullyTrueAndTiesRefuse` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-41 | `TestReq41_UnevaluableAtomInsideUnlessMakesTheWholeRowUnevaluable` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-42 | `TestReq42_OptionalKeyValueAtomInUnlessWithholdsExactlyAsInAll` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-43 | `TestReq43_CanRefuseRowContributesNoAssignmentsAndTheGroupIsUnprovable` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-44 | `TestReq44_WithheldGroupEmitsNoCoverageGapButStillEmitsOverlap` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-45 | `TestReq45_DimensionsAreEvaluatedAsOneProductNotIndependently` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-46 | `TestReq46_RowGroupIsTheRowsSharingOneSelectionContext` | HAPPY PATH | `internal/guard/guard_product_0003_test.go` |
| REQ-47 | `TestReq47_GroupingPredicateIsSelfContainedAndReadsNoTraversal` | BOUNDARY | `internal/guard/guard_product_0003_test.go` |
| REQ-48 | `TestReq48_ParticipationIsAnyGuardAtomInEitherBlockRegardlessOfOperator` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-49 | `TestReq49_IdenticallyConstrainedKeyStaysInTheProduct` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-50 | `TestReq50_MatchKeysAreNotProductDimensions` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-51 | `TestReq51_MatchGuardSplitIsPerAtomPerGroupNotPerKeyPerModel` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-52 | `TestReq52_ExistsProjectsOntoAPerKeyPresenceDimension` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-53 | `TestReq53_ValueAtomDenotesAValueSubsetAndImplicitlyPresent` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-54 | `TestReq54_ExistsProjectsOverCoOccurringKeysExactlyAsOverSingleValued` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-55 | `TestReq55_AlwaysPresentKeyContributesNoPresenceDimension` | BOUNDARY | `internal/guard/guard_product_0003_test.go` |
| REQ-56 | `TestReq56_ExistsAtomsAreNotDroppedFromTheProduct` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-57 | `TestReq57_SingleValueOperatorsNarrowAndContainsProjectsTheOtherWay` | DOMAIN EDGE | `internal/guard/guard_product_0003_test.go` |
| REQ-58 | `TestReq58_ValueAtomOverAnUnmarkedDimensionIsUnprojectable` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-59 | `TestReq59_LintPicksNoReadingForAnUnprojectableAtom` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-60 | `TestReq60_CoverageIsUnionEqualsProductAndOverlapIsNonEmptyIntersection` | HAPPY PATH | `internal/guard/guard_product_0003_test.go` |
| REQ-61 | `TestReq61_ExhaustivenessIsClaimableOnlyOverFiniteDeclaredDomains` | BOUNDARY | `internal/guard/guard_product_0003_test.go` |
| REQ-62 | `TestReq62_NoFiniteDomainRefusesRatherThanTreatingExamplesAsComplete` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-63 | `TestReq63_ClaimIsDefaultOnWithNoOptInGesture` | HAPPY PATH | `internal/guard/guard_product_0003_test.go` |
| REQ-64 | `TestReq64_UndeclaredDimensionIsNotAnOptOutAndEveryGroupIsProvedOrBlocked` | ADVERSARIAL | `internal/guard/guard_product_0003_test.go` |
| REQ-65 | `TestReq65_AllThreeRefuseOrDowngradeCasesAreBlocking` | BOUNDARY | `internal/guard/guard_product_0003_test.go` |
| REQ-66 | `TestReq66_DowngradeAndRefuseNameOneBlockingOutcome` | BOUNDARY | `internal/guard/guard_product_0003_test.go` |
| REQ-67 | `TestReq67_LintNarrowsWhereItWouldOutrunTheRuntimeVeto` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-68 | `TestReq68_NarrowingPopulationIncludesEscapeRows` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-69 | `TestReq69_WithholdingIsWholeGroupWhileOverlapStaysTwoPopulation` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-70 | `TestReq70_CanRefuseIsAValueAtomOverAnOptionalKeyInEitherBlock` | BOUNDARY | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-71 | `TestReq71_CanRefuseIsSyntacticOverDeclarationsAndTotal` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-72 | `TestReq72_WithholdingReadsTheOptionalityFieldNotOwnedReachability` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-73 | `TestReq73_WithheldClaimNamesTheRowTheAtomAndTheSourceID` | HAPPY PATH | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-74 | `TestReq74_EmittingNothingDoesNotSatisfyTheWithholdingClause` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-75 | `TestReq75_EscapeRowParticipatesInTheUnionAndABareOneClosesCoverage` | DOMAIN EDGE | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-76 | `TestReq76_EscapeClosesOnlyItsDeclaredRescuableClasses` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-77 | `TestReq77_CoverageUnionIsPerGroupTimesDeclaredRescuableClass` | DOMAIN EDGE | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-78 | `TestReq78_BareEscapeRowDoesNotDischargeTheNarrowing` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-79 | `TestReq79_OverlapIsTwoPopulationsWhileCoverageIsOne` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-80 | `TestReq80_EscapeRowsArePartitionedPerClassWithMultiClassRowsInEach` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-81 | `TestReq81_APairOverlappingInTwoClassesYieldsOneFindingPerClass` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-82 | `TestReq82_BareEscapeClosureIsObservableAndNamesTheRow` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-83 | `TestReq83_SetLiteralIsCanonicalUnorderedAndDuplicateFree` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-84 | `TestReq84_CanonicalizationHappensBeforeTheIdentityTuple` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-85 | `TestReq85_SetDomainsAreProvedDeterministicallyAndOverLargeRefuses` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-86 | `TestReq86_TheBoundIsPublishedDeclaredAndModelIndependent` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-87 | `TestReq87_OverLargeRefusalReportsComputedSizeAndBound` | HAPPY PATH | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-88 | `TestReq88_TheReportedQuantityIsTheProductCardinality` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-89 | `TestReq89_PerKindAssignmentCountTableIsEnforcedRowByRow` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-90 | `TestReq90_DefaultsAreAppliedNotTheAuthorsPresumedIntent` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-91 | `TestReq91_TheBoundIsOneConstantNotPerModelPerGroupOrConfigurable` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-92 | `TestReq92_TheVerdictIsStableAndTheAppliedBoundIsReported` | BOUNDARY | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-93 | `TestReq93_NoCardinalityIsComputedForAProductCarryingAnUnprovableDimension` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-94 | `TestReq94_UnprovableDimensionsAreReportedAndTheOverLargeRefusalIsNot` | ADVERSARIAL | `internal/guard/guard_escape_bound_0003_test.go` |
| REQ-95 | `TestReq95_EveryDecidableDefectIsReportedInOnePass` | ADVERSARIAL | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-96 | `TestReq96_TwoRefusingRowsEmitTwoFindingsInAnyIterationOrder` | ADVERSARIAL | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-97 | `TestReq97_NoGreenResultStandsAlongsideAWithholdingReason` | ADVERSARIAL | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-98 | `TestReq98_OverlapAndCoverageDiagnosticsNameEveryContributingSourceID` | HAPPY PATH | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-99 | `TestReq99_GapFindingNamesContextEveryRuleIDAndAWitnessAssignment` | ADVERSARIAL | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-100 | `TestReq100_OwnedTagMatchedWithNoReachablePredecessorWriteIsRejected` | DOMAIN EDGE | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-101 | `TestReq101_ThisRDRDefinesNoSecondReachabilityRelation` | BOUNDARY | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-102 | `TestReq102_TheThreeProvenancesCarryTheirThreeStatedRoles` | BOUNDARY | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-103 | `TestReq103_UnlessOverAbsentKeyIsLoudOnBothSurfaces` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-104 | `TestReq104_ValueAtomOverAbsentKeyIsWithheldWithRowAndAtomNamed` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-105 | `TestReq105_ExistenceOverAbsentKeySelectsAbsentAndDrawsNoDiagnostic` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-106 | `TestReq106_ZeroQualifyingIsAGapAndTwoQualifyingIsAnOverlap` | DOMAIN EDGE | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-107 | `TestReq107_UnprovableDimensionStaysRuntimeEvaluableAndDrawsOneFindingEach` | DOMAIN EDGE | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-108 | `TestReq108_OverLargeFiniteProductRefusesAndReportsSizeAndBound` | DOMAIN EDGE | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-109 | `TestReq109_EscapePairOverlapsWhileEscapeVsGuardedIsSilentButCounted` | DOMAIN EDGE | `internal/guard/guard_diagnostics_0003_test.go` |
| REQ-110 | `TestReq110_OneFindingPerRefusingRowAndAbsenceOfGreenIsNotTheArtifact` | ADVERSARIAL | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-111 | `TestReq111_DecidedFalseAndFullUnlessTrueAreSilentByDesign` | DOMAIN EDGE | `internal/guard/guard_narrowing_0003_test.go` |
| REQ-112 | `TestReq112_EveryVisibleFailureIsATypedLoadOrLintFailure` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-113 | `TestReq113_UndecidableAtRuntimeSurfacesAsGuardUnevaluableNotAParseKind` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-114 | `TestReq114_EveryGreenClaimIsTiedToTheCoverageIdentity` | ADVERSARIAL | `internal/guard/guard_scope_0003_test.go` |
| REQ-115 | `TestReq115_NoEncodeDecodePairIsIntroduced` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-116 | `TestReq116_NoHostPredicateAndNoFreeFormExpressionGrammar` | ADVERSARIAL | `internal/guard/guard_scope_0003_test.go` |
| REQ-117 | `TestReq117_NoNewThirdPartyDependencyIsIntroduced` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-118 | `TestReq118_NoPersistentResourceIsCreated` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-119 | `TestReq119_TheHotPathCarriesNoCallbackParserOrExternalEngine` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-120 | `TestReq120_TheCanonicalNameIsGuardPredicate` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-121 | `TestReq121_TheContainerIsRDR0002sAndTheAtomGrammarIsThisRDRs` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-122 | `TestReq122_TheMatrixAndAtomShapeAreSharedByResolverAndLint` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-123 | `TestReq123_AllFourFiniteKindsReachCoverageAndOverlapChecks` | HAPPY PATH | `internal/guard/guard_scope_0003_test.go` |
| REQ-124 | `TestReq124_MVVFixtureCarriesAContainsPredicateOverADeclaredSetTag` | HAPPY PATH | `internal/guard/guard_mvv_0003_test.go` |
| REQ-125 | `TestReq125_DiagnosticsCarryRDR0002IdentitiesAndRDR0006Codes` | BOUNDARY | `internal/guard/guard_scope_0003_test.go` |
| REQ-126 | `TestReq126_ScenarioTwosPartitionGroupDeclaresMarkersAndClosesRouting` | HAPPY PATH | `internal/guard/guard_mvv_0003_test.go` |
| REQ-127 | `TestReq127_EveryTargetFlowDimensionIsMarkedOrIsASetUnderContains` | BOUNDARY | `internal/guard/guard_mvv_0003_test.go` |
| REQ-128 | `TestReq128_ScenarioOneRuntimeEvaluationOverTheRepresentativeRows` | HAPPY PATH | `internal/guard/guard_mvv_0003_test.go` |
| REQ-129 | `TestReq129_ScenarioTwoPartitionPassesAndGapAndOverlapBothReportInOneRun` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-130 | `TestReq130_ScenarioThreeCoversUnboundedTooLargeTwoDimsAndShapePairs` | DOMAIN EDGE | `internal/guard/guard_mvv_0003_test.go` |
| REQ-131 | `TestReq131_ShapePairIsBuiltRelativeToThePublishedBound` | BOUNDARY | `internal/guard/guard_mvv_0003_test.go` |
| REQ-132 | `TestReq132_ScenarioThreeExpectationsIncludingEqualCardinalityAgreement` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-133 | `TestReq133_TheBoundsVerdictAgreesWithWhetherTheProofCompletes` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-134 | `TestReq134_ScenarioFourParsesEveryMalformedAtomAndDeclaration` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-135 | `TestReq135_ADeclarationErrorNeverReachesAConsumerReadingTheModel` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-136 | `TestReq136_ScenarioFiveMarkerChangesTheProductAndTheVerdict` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-137 | `TestReq137_ScenarioSixReorderingChangesNeitherMatchingNorFindings` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-138 | `TestReq138_ScenarioSevenBothBlocksAreRetainedAndSeparatelyIdentifiable` | BOUNDARY | `internal/guard/guard_mvv_0003_test.go` |
| REQ-139 | `TestReq139_ScenarioEightPositiveNarrowingAndItsNegativeControl` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-140 | `TestReq140_ScenarioEightBothBlocksAndTheEscapePath` | ADVERSARIAL | `internal/guard/guard_mvv_0003_test.go` |
| REQ-141 | `TestReq141_OneAtomSetDrivesSelectionOverlapAndProofWithoutOrderPriority` | HAPPY PATH | `internal/guard/guard_scope_0003_test.go` |
| REQ-MVV | `TestMVV_GuardPredicateExhaustiveness` | HAPPY PATH | `internal/guard/guard_mvv_0003_test.go` |

## Orphans

**REQs with no test:** none. Every one of the 141 REQs plus `REQ-MVV`
carries a test whose header quotes the clause verbatim.

**Tests citing no REQ:** none. Every `func Test…` in `internal/guard/`
opens with a `// REQ-N: "<quote>"` header and exactly one of the five
labels. Fixture builders and inspection helpers live in
`guard_fixtures_0003_test.go` and declare no test functions.

## REQ-MVV decomposition

`TestMVV_GuardPredicateExhaustiveness` is the gating runnable end-to-end
validation. Its five obligations are five independently-failing subtests:

| Obligation | Subtest |
| --- | --- |
| 1 — one RDR flow slice and one kata flow slice over the full operator vocabulary and mixed `all`/`unless` | `1-two-flow-slices-over-the-full-vocabulary` |
| 2 — one scoped row group proved exhaustive *and* mutually exclusive | `2-one-group-proved-exhaustive-and-mutually-exclusive` |
| 3 — one intentional gap and one intentional overlap, product-only, both from one run, with source ids and a concrete uncovered assignment | `3-gap-and-overlap-from-one-run-with-source-ids` |
| 4 — the **positive** narrowing assertion: a run that emits nothing fails | `4-positive-blocking-finding-names-the-row-and-the-atom` |
| 5 — the **negative control**: an always-present group lint certifies green | `5-negative-control-certifies-green` |

Obligation 4 asserts positively — it fails on an empty finding set before
it looks at anything else, so asserting mere absence of a green result
cannot discharge it. Obligation 5 is what makes the narrowing tight
rather than blanket: a lint that certifies nothing fails it.

This RDR declares **no** Round-Trip / Inverse invariant (REQ-115), so the
MVV carries no `X∘Y = identity` fidelity obligation.

## Label distribution

| Label | Tests |
| --- | --- |
| HAPPY PATH | 19 |
| INPUT EDGE | 1 |
| BOUNDARY | 42 |
| ADVERSARIAL | 49 |
| DOMAIN EDGE | 31 |

## Files

| File | Scope |
| --- | --- |
| `guard_grammar_0003_test.go` | Atom grammar, closed operator vocabulary, operator/kind matrix, the two equalities (REQ-1 … REQ-12) |
| `guard_declaration_0003_test.go` | The tag declaration model this RDR owns (REQ-13 … REQ-33) |
| `guard_evaluator_0003_test.go` | Evaluator scope and the kernel split (REQ-34 … REQ-37) |
| `guard_narrowing_0003_test.go` | Polarity, `unless`, and the runtime-veto narrowing (REQ-38 … REQ-44, REQ-67 … REQ-74, REQ-103 … REQ-105, REQ-110, REQ-111) |
| `guard_product_0003_test.go` | Row groups, participation, the scoped product, eligibility (REQ-45 … REQ-66) |
| `guard_escape_bound_0003_test.go` | Escape rows, set literals, the published cardinality bound (REQ-75 … REQ-94) |
| `guard_diagnostics_0003_test.go` | Report-every-defect, attribution, provenance, remaining disposition rows (REQ-95 … REQ-102, REQ-106 … REQ-109) |
| `guard_scope_0003_test.go` | Failure modes, non-goals, naming, wire ownership, phases (REQ-112 … REQ-125, REQ-141) |
| `guard_mvv_0003_test.go` | Testing Strategy scenarios and the MVV (REQ-124, REQ-126 … REQ-140, `REQ-MVV`) |
| `guard_fixtures_0003_test.go` | Shared fixture builders and lint-result inspection. Declares no tests. |

