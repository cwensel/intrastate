# Impact — 0030-computed-write-value-grammar

families: 11
rows: 168
records: 0002 0012
literals: `is not a tag value` `guard.Evaluator` `Evaluator{}` `NewEvaluator` `intWidth` `atomAdmitsValue` `renderSet` `groupHasOverlap`
repo: /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0030 @48049f9
convention: go

## seam_carrier_0012_test.go (32 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq10_0012_NewEvaluatorReturnsTheConcreteEvaluator | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq11_0012_GuardEqAndInCompareUnderTheDeclaredKind | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq12_0012_IntComparesParsedValues | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq13_0012_OneUnparseableInMemberPoisonsTheList | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq14_0012_OverflowIsUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq15_0012_BoolComparesTokensOnly | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq16_0012_EnumAndScalarCompareExactStrings | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq17_0012_SetEqAndInAreUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq18_0012_AbsentKeyIsUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq19_0012_EmptyKindDeclaredKeyIsUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq1_0012_NewEvaluatorIsDeclaredWithTheCarrierSignature | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq20_0012_UnparseableLiteralIsUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq21_0012_UnknownKindTokenIsUnevaluable | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq22_0012_OrderingAndContainsDoNotConsultTheMapping | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq23_0012_UnknownOperatorIsUnevaluableAndNeverPanics | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq24_0012_InMembersAreTheDecodedArray | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq25_0012_TheSuiteCarriesTheTypedAndDefensiveLegs | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq26_0012_EachKindTokenHasADiscriminatingCase | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq27_0012_ContractKindsHandsEachCallerAFreshCopy | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq29_0012_CasesAreKeyedOntoAdmittedKinds | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq2_0012_DeclaredKindsIsHomedInGuardAndMapsKeyToKind | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq31_0012_TheGuardSuiteCallIsRePointedThroughAOneLineAdapter | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq33_0012_TheLintSiteConstructsOverTheModelsDeclarations | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq35_0012_AtomAdmitsValueIsTheSoleRawByteSiteAndKeepsItsShape | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq3_0012_DeclaredKindsIsTotalAndOmitsEmptyKind | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq46_0012_PerKindDispatchMatrix | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq47_0012_EvaluatorHoldsTheDeclarationMappingAndNothingElse | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq49_0012_AZeroValueEvaluatorNeverRevertsToRawComparison | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq4_0012_DeclaredKindsOverANilModelIsEmptyNonNil | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReq8_0012_TheSeamSignatureAndAtomShapeAreUnchanged | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |
| TestReqMVV_0012_Step4_TheExtendedSuiteIsGreenOverNewEvaluator | internal/guard/seam_carrier_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator;literal:atomAdmitsValue |

## rdr0012_probe_test.go (31 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestRDR0012Probe_Req12_IntComparesParsedValues | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req13_OneUnparseableInMemberPoisonsTheList | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req14_OverflowIsUnevaluableAndInRangeIsPinned | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req15_BoolComparesTokensOnly | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req16_EnumAndScalarCompareExactStrings | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req17_SetEqAndInAreUnevaluable | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req18_AbsentKeyIsUnevaluable | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req19_EmptyKindDeclaredKeyIsUnevaluable | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req1_NewEvaluatorConstructsTheConcreteSeam | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req20_UnparseableLiteralIsUnevaluable | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req21_TheSuitePinsAnUnknownKindToken | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req21_UnknownKindTokenIsUnevaluable | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req22_OrderingAndContainsDoNotConsultTheMapping | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req23_UnknownOperatorIsUnevaluableAndNeverPanics | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req24_InMembersAreTheDecodedArray | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req25_TheSuiteAsksTheNewLegs | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_FixtureTokensAreExactlyTheVocabulary | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_Mutant_bool | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_Mutant_enum | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_Mutant_int | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_Mutant_scalar | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req26_Mutant_set | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req27_ContractKindsHandsEachCallerItsOwnCopy | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req29_CasesAreKeyedOntoAdmittedKinds | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req2_DeclaredKindsMapsKeyToKindToken | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req3_DeclaredKindsIsTotalAndOmitsEmptyKind | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req46_PerKindDispatchMatrix | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req4_DeclaredKindsOverNilModelIsEmptyNonNil | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_Req8_SeamSignatureAndAtomShapeAreUnchanged | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_ReqMVV_ContractSuiteIsGreenOverNewEvaluator | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestRDR0012Probe_ReqMVV_Mutant_rawStringSeam | internal/guard/rdr0012_probe_test.go | literal:guard.Evaluator;literal:NewEvaluator |

## guard_declaration_0003_test.go (21 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq13_DeclarationCarriesAKindAndTheFourOptionalFields | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq14_ValueKindsAreExactlyFiveTokens | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq15_ScalarIsOpaqueAndNeverBearsAnExhaustivenessClaim | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq16_ThisRDROwnsMeaningNotAuthoringOrCarriage | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq17_FiniteDomainIsDeclarableForEveryClaimableKind | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq18_IntBoundEndpointsAreInclusive | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq19_NoFiniteDomainIsWellFormedButUnprovable | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq20_OmittedOptionalityMarkerDeclaresTheKeyOptional | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq21_ConformanceRequiresEveryAlwaysPresentKey | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq22_PresenceIsReadFromTheDeclarationNotAView | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq23_ConformanceIsBothConjunctsAndScopesEveryGreenClaim | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq24_NoClaimIsMadeAboutANonConformingView | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq25_SingleValuedMarkerIsStatableAndBannedOnSet | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq26_SingleValuedContributesDomainSizeUnmarkedContributesPowerset | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq27_SingleValuedIsNeverInferredFromNameSpellingOrFixture | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq28_DomainMustAgreeWithKindAndDisagreementIsRejected | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq29_KindFieldAgreementTableIsEnforcedPerCell | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq30_LiteralOutsideDeclaredDomainIsRejectedAtLoad | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq31_DeclarationErrorsAndLiteralErrorsAreRejectedAtDifferentStages | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq32_TheTwoRejectionsMapOntoRDR0002sTwoLoadCategories | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |
| TestReq33_TheDeclarationModelIsExactlyFiveFields | internal/guard/guard_declaration_0003_test.go | literal:guard.Evaluator |

## guard_narrowing_0003_test.go (20 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq103_UnlessOverAbsentKeyIsLoudOnBothSurfaces | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq104_ValueAtomOverAbsentKeyIsWithheldWithRowAndAtomNamed | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq105_ExistenceOverAbsentKeySelectsAbsentAndDrawsNoDiagnostic | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq110_OneFindingPerRefusingRowAndAbsenceOfGreenIsNotTheArtifact | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq111_DecidedFalseAndFullUnlessTrueAreSilentByDesign | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq38_PolarityIsByBlockAndMatchingIsOrderIndependent | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq39_AcceptedAssignmentsAreAllIntersectionMinusTheUnlessConjunction | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq40_RowQualifiesOnAllTrueAndUnlessNotFullyTrueAndTiesRefuse | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq41_UnevaluableAtomInsideUnlessMakesTheWholeRowUnevaluable | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq42_OptionalKeyValueAtomInUnlessWithholdsExactlyAsInAll | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq43_CanRefuseRowContributesNoAssignmentsAndTheGroupIsUnprovable | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq44_WithheldGroupEmitsNoCoverageGapButStillEmitsOverlap | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq67_LintNarrowsWhereItWouldOutrunTheRuntimeVeto | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq68_NarrowingPopulationIncludesEscapeRows | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq69_WithholdingIsWholeGroupWhileOverlapStaysTwoPopulation | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq70_CanRefuseIsAValueAtomOverAnOptionalKeyInEitherBlock | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq71_CanRefuseIsSyntacticOverDeclarationsAndTotal | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq72_WithholdingReadsTheOptionalityFieldNotOwnedReachability | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq73_WithheldClaimNamesTheRowTheAtomAndTheSourceID | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |
| TestReq74_EmittingNothingDoesNotSatisfyTheWithholdingClause | internal/guard/guard_narrowing_0003_test.go | literal:NewEvaluator |

## guard_typed_0012_test.go (19 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq11_0012_GuardAtomsCompareTypedAndMatchAtomsCompareBytes | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq33_0012_TheCLISiteConstructsOverTheRequestsModel | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq34_0012_TheEvaluatorIsThreadedDownNotReDerived | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq36_0012_ATypedUnevaluableIsNeitherPrunedNorMatched | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq37_0012_TheUserFacingDocsStateTheGuardMatchSplit | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq38_0012_ATypedUnevaluableRidesTheExistingRefusal | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq39_0012_NonCanonicalPredicateIntLiteralsAreRefused | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq40_0012_TheBareFloatNegativeZeroIsRefused | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq41_0012_TheCheckReachesPredicateLiteralsOnly | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq42_0012_CLIValuesStayAcceptedAndCompareTyped | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq43_0012_TheRefusalNamesTheSiteAndTheRewrite | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq44_0012_TheRefusalReusesMalformedPredicateAtom | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq45_0012_CanonicalIsExactlyTheItoaRoundTrip | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq50_0012_NoCommittedModelsLintVerdictFlips | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq51_0012_Cover07RefusesWithTheCanonicalSpellingDiagnostic | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq52_0012_TheTagDoorRefusesUpstreamWhileTheOwnedDoorReachesTheSeam | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq53_0012_AnOwnedNonCanonicalValueMatchesUnderParsedComparison | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReq54_0012_FlowNextReportsTypedVerdicts | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |
| TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor | internal/cli/guard_typed_0012_test.go | literal:NewEvaluator |

## guard_scope_0003_test.go (14 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq112_EveryVisibleFailureIsATypedLoadOrLintFailure | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq113_UndecidableAtRuntimeSurfacesAsGuardUnevaluableNotAParseKind | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq114_EveryGreenClaimIsTiedToTheCoverageIdentity | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq115_NoEncodeDecodePairIsIntroduced | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq116_NoHostPredicateAndNoFreeFormExpressionGrammar | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq117_NoNewThirdPartyDependencyIsIntroduced | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq118_NoPersistentResourceIsCreated | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq119_TheHotPathCarriesNoCallbackParserOrExternalEngine | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq120_TheCanonicalNameIsGuardPredicate | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq121_TheContainerIsRDR0002sAndTheAtomGrammarIsThisRDRs | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq122_TheMatrixAndAtomShapeAreSharedByResolverAndLint | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq123_AllFourFiniteKindsReachCoverageAndOverlapChecks | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq125_DiagnosticsCarryRDR0002IdentitiesAndRDR0006Codes | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |
| TestReq141_OneAtomSetDrivesSelectionOverlapAndProofWithoutOrderPriority | internal/guard/guard_scope_0003_test.go | literal:NewEvaluator |

## guard_grammar_0003_test.go (13 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq10_DiagnosticsUseIdentityAndDomainComputationUsesSemanticEquality | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq11_ByteIdenticalGuardsInTwoRulesOverlapRatherThanDeduplicate | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq12_OperatorSemanticsAreAsStated | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq1_GuardIsASymbolicAtomOverADeclaredTag | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq2_OperatorVocabularyIsClosedAndTyped | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq2_SingleValueOperatorsAreClassifiedOverTheWholeVocabulary | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq3_EachOperatorDeclaresAcceptedKindsAndLiteralsAreKindChecked | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq4_OperatorKindMatrixIsExactlyAsPublished | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq5_ExistsOverAnAlwaysPresentKeyIsVacuousNotRejected | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq6_AtomWellFormednessIsKeyOperatorAndLiteral | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq7_SourceIdentityIsNotAnAtomFieldButJoinsTheIdentityTuple | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq8_AtomIdentityIsTheSixFieldTupleNeverThePositionalIndex | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq9_SemanticEqualityIsTagOperatorLiteral | internal/guard/guard_grammar_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |

## guard_adversarial_0003_test.go (8 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestAdv1_GreenClaimCoversAViewTheRuntimeRefuses | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv2_OverLargeSetProductIsRefusedWithoutEnumeratingIt | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv3_UnprovableDimensionSilencesAnUnrelatedOverlap | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv_NarrowIntDomainAtMaxIntTerminates | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv_OptionalPresenceFactorSaturatesRatherThanWrapping | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv_RecordedIntBoundInversion | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv_WideIntDomainRefusesRatherThanCertifying | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |
| TestAdv_WideIntDomainWidthSaturatesRatherThanWrapping | internal/guard/guard_adversarial_0003_test.go | literal:intWidth |

## contract_fixture_0012_test.go (5 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq27_0012_ContractKindsIsAFunctionNotAPackageVar | internal/resolve/contract_fixture_0012_test.go | literal:guard.Evaluator |
| TestReq28_0012_TheSuiteTakesASeamConstructor | internal/resolve/contract_fixture_0012_test.go | literal:guard.Evaluator |
| TestReq30_0012_TheConformingTestSeamIsBoundThroughAConstructor | internal/resolve/contract_fixture_0012_test.go | literal:guard.Evaluator |
| TestReq32_0012_TheImportablePinIsReTypedNotDeleted | internal/resolve/contract_fixture_0012_test.go | literal:guard.Evaluator |
| TestReq9_0012_TheResolutionPathReadsNoDeclarations | internal/resolve/contract_fixture_0012_test.go | literal:guard.Evaluator |

## guard_evaluator_0003_test.go (4 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq35_ExistenceOverAbsentIsDecidedAndValueOverAbsentIsUnevaluable | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq36_AbsenceIsExpressedWithAnExistenceAtomNotAMissingValueAtom | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |
| TestReq37_ValueVerdictIsThisRDRsAndPresenceCombinationIsTheKernels | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator;literal:NewEvaluator |

## guard_adv_0012_test.go (1 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestAdv1_0012_TheZeroValueCheckCatchesEveryZeroValueForm | internal/guard/guard_adv_0012_test.go | literal:guard.Evaluator;literal:Evaluator{};literal:NewEvaluator |
