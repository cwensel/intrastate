# Impact — 0012-declared-kind-carrier-at-the-guard-seam

families: 10
rows: 128
records: 0003 0007
literals: `guard.Evaluator{}` `var ev guard.Evaluator` `guardSeam()` `func(*testing.T, resolve.GuardEvaluator)` `conformingContractSeam{}` `TestGuardEvaluatorContract(` `probeRow(`
repo: /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0012 @901ae3d
convention: go

## guard_declaration_0003_test.go (21 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq13_DeclarationCarriesAKindAndTheFourOptionalFields | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq14_ValueKindsAreExactlyFiveTokens | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq15_ScalarIsOpaqueAndNeverBearsAnExhaustivenessClaim | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq16_ThisRDROwnsMeaningNotAuthoringOrCarriage | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq17_FiniteDomainIsDeclarableForEveryClaimableKind | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq18_IntBoundEndpointsAreInclusive | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq19_NoFiniteDomainIsWellFormedButUnprovable | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq20_OmittedOptionalityMarkerDeclaresTheKeyOptional | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq21_ConformanceRequiresEveryAlwaysPresentKey | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq22_PresenceIsReadFromTheDeclarationNotAView | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq23_ConformanceIsBothConjunctsAndScopesEveryGreenClaim | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq24_NoClaimIsMadeAboutANonConformingView | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq25_SingleValuedMarkerIsStatableAndBannedOnSet | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq26_SingleValuedContributesDomainSizeUnmarkedContributesPowerset | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq27_SingleValuedIsNeverInferredFromNameSpellingOrFixture | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq28_DomainMustAgreeWithKindAndDisagreementIsRejected | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq29_KindFieldAgreementTableIsEnforcedPerCell | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq30_LiteralOutsideDeclaredDomainIsRejectedAtLoad | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq31_DeclarationErrorsAndLiteralErrorsAreRejectedAtDifferentStages | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq32_TheTwoRejectionsMapOntoRDR0002sTwoLoadCategories | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |
| TestReq33_TheDeclarationModelIsExactlyFiveFields | internal/guard/guard_declaration_0003_test.go | literal:var ev guard.Evaluator |

## guard_narrowing_0003_test.go (20 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq103_UnlessOverAbsentKeyIsLoudOnBothSurfaces | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq104_ValueAtomOverAbsentKeyIsWithheldWithRowAndAtomNamed | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq105_ExistenceOverAbsentKeySelectsAbsentAndDrawsNoDiagnostic | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq110_OneFindingPerRefusingRowAndAbsenceOfGreenIsNotTheArtifact | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq111_DecidedFalseAndFullUnlessTrueAreSilentByDesign | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq38_PolarityIsByBlockAndMatchingIsOrderIndependent | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq39_AcceptedAssignmentsAreAllIntersectionMinusTheUnlessConjunction | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq40_RowQualifiesOnAllTrueAndUnlessNotFullyTrueAndTiesRefuse | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq41_UnevaluableAtomInsideUnlessMakesTheWholeRowUnevaluable | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq42_OptionalKeyValueAtomInUnlessWithholdsExactlyAsInAll | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq43_CanRefuseRowContributesNoAssignmentsAndTheGroupIsUnprovable | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq44_WithheldGroupEmitsNoCoverageGapButStillEmitsOverlap | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq67_LintNarrowsWhereItWouldOutrunTheRuntimeVeto | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq68_NarrowingPopulationIncludesEscapeRows | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq69_WithholdingIsWholeGroupWhileOverlapStaysTwoPopulation | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq70_CanRefuseIsAValueAtomOverAnOptionalKeyInEitherBlock | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq71_CanRefuseIsSyntacticOverDeclarationsAndTotal | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq72_WithholdingReadsTheOptionalityFieldNotOwnedReachability | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq73_WithheldClaimNamesTheRowTheAtomAndTheSourceID | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |
| TestReq74_EmittingNothingDoesNotSatisfyTheWithholdingClause | internal/guard/guard_narrowing_0003_test.go | literal:guard.Evaluator{} |

## flow_next_0011_test.go (16 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq2And67_TheProbeBindsItsOwnOutcomeStripsEscapeAndKeepsTheResult | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq34And130_TheKernelPackageIsUntouchedByThisContract | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq4_AnEmptyButPresentOwnedValueIsPresentNotAbsent | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq77And101_APresentAndEqualMatchKeyLeavesNoUnknownEntryForThatKey | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq7And103_AMixedStateRowIsExcludedByTheUnequalKeyAlone | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq82And104_ARecognizedOnlyRowIsACandidateWithNoMatchFacts | internal/cli/flow_next_0011_test.go | literal:probeRow( |
| TestReq83And96_GuardExcludedRowsStayExcludedInBothModes | internal/cli/flow_next_0011_test.go | literal:probeRow( |

## escape_shape_0009_test.go (15 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq109_ABreachSurfacedThroughTheCLICarriesCodeIdentityAndHint | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq109_ANonBreachKernelErrorStillFallsThroughToTheGenericBranch | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq109_TheCLICarriesEveryOffendingIdentityInKernelOrder | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq58_ABreachReachingTheCLIIsWrappedWithCodeGroupIdentityAndHint | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq59_TheRowIdentitiesReachTheSerializedEnvelope | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq60_TheCarrierIsFindingsAndNotDetail | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq61_TheCarrierIsAStructuredRepeatedSerializedField | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq62_ClierrStaysALeafAndTheConversionIsTheVerbLayers | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq63_ThePerIdentityCountSerializesAlongsideTheIdentities | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq64_TheStableCodeIsSpelledExactly | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq65_TheStaleCauseCommentIsAmended | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq66_TheOutputContractDocumentsTheIdentityCarrier | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq67_TheExitCodeCannotDiscriminateAMissingWrap | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq68_TheRefusalTaxonomyStaysClosed | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |
| TestReq69_TheRefusalToExitCodeMappingIsUnchanged | internal/cli/escape_shape_0009_test.go | literal:guardSeam() |

## guard_mvv_test.go (14 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestMVV_GuardPredicateTotality | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq59_RequiresOwnedDocContractIsNarrowedToPostGuardWrites | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq60_GuardInputCoverageIsNotRequiresOwnedsJob | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq61_ListingAGuardKeyInRequiresOwnedProtectsAgainstObservedSubstitution | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq62_ConformingEscapeRowRaisesNoOwnedStateOfItsOwn | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq68_VerdictsAreDrivenThroughAtomsNotInjectedAtRowLevel | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq69_TheGuardTextFieldIsActuallyRemoved | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq71_AConformingValueSeamSatisfiesTheContractTest | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq71_ContractTestExercisesThePresentValueLiteralOperatorProduct | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq73_ScenariosAreKernelTestsInPackageResolveTest | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq74_FalseVsUnevaluableDominanceMutantIsKilled | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq75_ObservedTagSubstitutionYieldsASpecificVerdict | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |
| TestReq76_TwoRowAbsencePatternIsTotalAndTheBareValueRowIsNot | internal/resolve/guard_mvv_test.go | literal:func(*testing.T, resolve.GuardEvaluator);literal:conformingContractSeam{};literal:TestGuardEvaluatorContract( |

## guard_scope_0003_test.go (14 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq112_EveryVisibleFailureIsATypedLoadOrLintFailure | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq113_UndecidableAtRuntimeSurfacesAsGuardUnevaluableNotAParseKind | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq114_EveryGreenClaimIsTiedToTheCoverageIdentity | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq115_NoEncodeDecodePairIsIntroduced | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq116_NoHostPredicateAndNoFreeFormExpressionGrammar | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq117_NoNewThirdPartyDependencyIsIntroduced | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq118_NoPersistentResourceIsCreated | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq119_TheHotPathCarriesNoCallbackParserOrExternalEngine | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq120_TheCanonicalNameIsGuardPredicate | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq121_TheContainerIsRDR0002sAndTheAtomGrammarIsThisRDRs | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq122_TheMatrixAndAtomShapeAreSharedByResolverAndLint | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq123_AllFourFiniteKindsReachCoverageAndOverlapChecks | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq125_DiagnosticsCarryRDR0002IdentitiesAndRDR0006Codes | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |
| TestReq141_OneAtomSetDrivesSelectionOverlapAndProofWithoutOrderPriority | internal/guard/guard_scope_0003_test.go | literal:var ev guard.Evaluator |

## guard_grammar_0003_test.go (13 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq10_DiagnosticsUseIdentityAndDomainComputationUsesSemanticEquality | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq11_ByteIdenticalGuardsInTwoRulesOverlapRatherThanDeduplicate | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq12_OperatorSemanticsAreAsStated | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq1_GuardIsASymbolicAtomOverADeclaredTag | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq2_OperatorVocabularyIsClosedAndTyped | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq2_SingleValueOperatorsAreClassifiedOverTheWholeVocabulary | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq3_EachOperatorDeclaresAcceptedKindsAndLiteralsAreKindChecked | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq4_OperatorKindMatrixIsExactlyAsPublished | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq5_ExistsOverAnAlwaysPresentKeyIsVacuousNotRejected | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq6_AtomWellFormednessIsKeyOperatorAndLiteral | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq7_SourceIdentityIsNotAnAtomFieldButJoinsTheIdentityTuple | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq8_AtomIdentityIsTheSixFieldTupleNeverThePositionalIndex | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |
| TestReq9_SemanticEqualityIsTagOperatorLiteral | internal/guard/guard_grammar_0003_test.go | literal:var ev guard.Evaluator |

## TestAdv0007 (7 tests, 2 files)
| test | file | arm |
|---|---|---|
| TestAdv0007_1_MatchBlockAtomIsNotAGuardOperand | internal/resolve/guard_adversarial_0007_test.go | record:0007 |
| TestAdv0007_2_UnfencedBlockFailsOpenAndReopensTheMaskingPath | internal/resolve/guard_adversarial_0007_test.go | record:0007 |
| TestAdv0007_3_DuplicateKeyMakesTheVerdictAFunctionOfSlicePosition | internal/resolve/guard_adversarial_0007_test.go | record:0007 |
| TestAdv0007_3_ConflictedKeyIsNotEscapableByNamingIt | internal/resolve/match_conflicted_test.go | record:0007 |
| TestAdv0007_3_DuplicateKeyMakesRowSelectionAFunctionOfSlicePosition | internal/resolve/match_conflicted_test.go | record:0007 |
| TestAdv0007_3_EscapeNotNamingTheConflictedKeyStillRescues | internal/resolve/match_conflicted_test.go | record:0007 |
| TestAdv0007_3_MatchNarrowingsHold | internal/resolve/match_conflicted_test.go | record:0007 |

## TestAdv0009 (4 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestAdv0009_ANonBreachKernelErrorKeepsTheGenericInternalCode | internal/cli/escape_shape_adv_0009_test.go | literal:guardSeam() |
| TestAdv0009_MultiBreachTraversalSurvivesAWrappedKernelError | internal/cli/escape_shape_adv_0009_test.go | literal:guardSeam() |
| TestAdv0009_TheCLIBreachEnvelopeIsAFunctionOfTheTableNotRowOrder | internal/cli/escape_shape_adv_0009_test.go | literal:guardSeam() |
| TestAdv0009_ThePerIdentityCountRoundTripsThroughJSON | internal/cli/escape_shape_adv_0009_test.go | literal:guardSeam() |

## guard_evaluator_0003_test.go (4 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator{};literal:var ev guard.Evaluator;literal:TestGuardEvaluatorContract( |
| TestReq35_ExistenceOverAbsentIsDecidedAndValueOverAbsentIsUnevaluable | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator{};literal:var ev guard.Evaluator;literal:TestGuardEvaluatorContract( |
| TestReq36_AbsenceIsExpressedWithAnExistenceAtomNotAMissingValueAtom | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator{};literal:var ev guard.Evaluator;literal:TestGuardEvaluatorContract( |
| TestReq37_ValueVerdictIsThisRDRsAndPresenceCombinationIsTheKernels | internal/guard/guard_evaluator_0003_test.go | literal:guard.Evaluator{};literal:var ev guard.Evaluator;literal:TestGuardEvaluatorContract( |
