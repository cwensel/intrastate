# Impact — 0019-owned-state-initialization-semantics

families: 11
rows: 104
records: 0002 0004 0005 0006 0010
literals: `four verbs` `four flow verbs` `"next", "resolve", "read-state", "set-state"` `flowVerbs`
repo: /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0019 @fcb85d6
convention: go

## flow_input_0005_test.go (20 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq112_NoNewThirdPartyDependencyIsIntroduced | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq22And90_ModelSelectionArityAndNonResolutionShareOneCode | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq23_ModelFileIOBelongsToTheCLINotTheLoader | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq24And91_LoadCategoriesRideOneCodeWithPerCategoryFindings | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq25_RevisionIsNeverCLIDerivedFromPathOrContent | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq26And114_TagEntersAsObservedAndNeverSatisfiesAnOwnedDependency | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq27And82And83_OwnedAndReservedTagsAreRefusedBeforeAnyAccessorRuns | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq28And81_DuplicateTagNameIsItsOwnCode | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq31And32And88_ArtifactLocationComesOnlyFromExplicitBindings | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq33And89_MissingBindingForAnInvokedRoleCarriesTheRoleAsParam | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq34_SetStateNeverTreatsTagValuesAsWrites | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq61And84_AWriteValueIsHeldToItsDeclaredDomainAndBounds | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq61And86And87_UnboundWriteAndClearKeysAreRefusedBeforeAccessors | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq61_AnUnboundKeyIsReportedBeforeItsValuesDomain | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq62And84_ClearSentinelIsUnauthorableAndHintsAtTheFlag | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq63And85_DuplicateKeyAcrossEitherFlagIsOneCode | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq69_ThePlanCarrierDoesNotDisturbThePlanOnlyFence | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq69_ThePlanCarrierRidesSetStateAlone | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestReq80And26_ADeclaredTagIsDomainCheckedAndAnUndeclaredOneIsNot | internal/cli/flow_input_0005_test.go | literal:four verbs;literal:flowVerbs |

## flow_projection_0023_test.go (20 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq10And11And12And132_PlanOnlyOnASiblingVerbIsTheSharedUsageBucketAndMintsNoCode | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq17And18And22And23And98And108_TheProjectedKeySetIsExactlyTheAuthoredLiteral | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq19And79And113_AProjectedKeyIsAbsentNotNullNotAnEmptyPlaceholder | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq1And15And16_ResolveRegistersABooleanPlanOnlyFlagDefaultFalseNoShorthand | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq20And26And124And127_EveryCarriedFieldIsByteIdenticalToItsDefaultRendering | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq21And121And122_TheCarriedKeysKeepTheirDefaultModeRelativeOrder | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq24And25And123_PresenceRuleFieldsAreAssertedAgainstTheSameRunsDefaultNotALiteral | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq29And30_The0011AllFenceIsUntouchedAndInvisibleToPlanOnly | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq2And3And4And5And7And8And9And13And93_PlanOnlyIsRegisteredOnExactlyFlowResolveAcrossTheWholeTree | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq31And32And34And45And130And133_TheDifferentialDiffsToExactlyTheEchoKeys | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq33And100And131_RefusalEnvelopesAreByteIdenticalPlusOrMinusTheFlag | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq36And38And111_TheProjectedEncodingIsStrictlyShorterOnTheFullEmittedLineInBothModes | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq39And57And82_TheTextSubsetAndTextWidthAssertionsAreSeparablyNecessary | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq47And48And49And50And53_NothingUpstreamOfPayloadAssemblyReadsTheFlag | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq51And52And67_ProjectionTakesTheFullyAssembledPayloadAndOnlyDeletes | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq54And94_BothModesRenderTheSameProjectedResultThroughTheOneGateway | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq55And56And58And59And110_ProjectedTextLinesAreAByteIdenticalStableSubset | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq6And14_TheShippedPartialWalkerCannotReachC1sScope | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq81_TheWholeTreeWalkGoesRedWhenASecondCommandOrCompletionRegistersTheFlag | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |
| TestReq92_TheFlagRidesResolveAloneAndHandsRespondOKTheProjectedResult | internal/cli/flow_projection_0023_test.go | literal:four verbs;literal:four flow verbs;literal:flowVerbs |

## flow_all_0011_test.go (11 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq123And126_NextIsEffectFreeInBothModes | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq43And128_EveryOther0005ObligationHoldsIdenticallyInBothModes | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq44And71_TheUnknownFieldIsPresentAsAnEmptyListInBothModes | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq45And76_AllIsPartOfRequestIdentity | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |
| TestReq65_TheSurfaceFlagMapNamesAllOnNext | internal/cli/flow_all_0011_test.go | literal:four verbs;literal:flowVerbs |

## flow_surface_0005_test.go (11 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq11And120_TextRenderingCarriesEveryJSONScalarValue | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq12_FlowGroupDoesNotRedefineTheRootAsFlag | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq13And19_NoNewTerminalEnvelopeTypeAndNoNewExitGroup | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq2_FlowGroupDoesNotAbsorbLintDumpOrParse | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq3_EachVerbRegistersItsNormativeFlagSpellings | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq5_EveryVerbValidatesOutputModeFirst | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq5_VerbsSilenceCobraErrorAndUsageOutput | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq6_JSONSuccessIsExactlyOneOKEnvelopeWithVerbSpecificData | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq7And10_TextSuccessEmitsHumanOutputFromTheSameResult | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |
| TestReq9And10_TextPayloadsRouteThroughTheGatewayNotDirectPrinting | internal/cli/flow_surface_0005_test.go | literal:four verbs;literal:"next", "resolve", "read-state", "set-state";literal:flowVerbs |

## flow_rehome_0011_test.go (10 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq108_TheUpstreamKernelPinsStillExistAsShipped | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq109_TheCLIViewsKeySetAgreesWithTheKernelsModuloRecognized | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq52_TheHelpSeparatesCandidateFromWhatResolveWillSelect | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq53And111_The0005NextFileHeaderNamesThisRDRAsTheSourceOfTheDefault | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq54And110_TheFiveShippedReadsOfUnresolvedAreReHomedToUnknown | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq56And122_TheShippedDescriptionsOfTheOldDefaultAreCorrected | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq58_TheReHomedReadsAssertTheOkBoolRatherThanDiscardingIt | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq59_StringsAtKeepsItsShapeAndTheReHomingAddsASiblingHelper | internal/cli/flow_rehome_0011_test.go | literal:four verbs |
| TestReq60_TheCommentAndFailureMessageProseNamingUnresolvedIsSwept | internal/cli/flow_rehome_0011_test.go | literal:four verbs |

## flow_exit_0005_test.go (8 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestReq102_TheThreeInternalCodesAreUnreachableFromTheUserPath | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq117_TheResolveToSetStateWindowIsAStatedNonGuarantee | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq131_TheOutputContractDocumentsFindingsAndTheExit3Rule | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq19And20_TheExit3PopulationIsExactlyTheFiveNamedClasses | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq19_AnExit3RefusalIsReRunnableUnchangedToTheSameDisposition | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq20_EveryAccessorRefusalClassHasACLIMapping | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq21_UnavailableAccessorExits3ButIndeterminateGateExits2 | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |
| TestReq99And100And101_AccessorRefusalsCarryTheAccessorIdAsParam | internal/cli/flow_exit_0005_test.go | literal:flowVerbs |

## flow_mvv_0023_test.go (7 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestD8_0024_TheDispositionsFoldAcceptsTheEmptyObjectAndNothingElse | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestD8_0024_ThePopulatedDispositionsRefusalNamesReq97sHole | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestFoldCheckoutRootEncodesThePrefixTheWayTheRecordDoes | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestReq115And116And117And118And119And120_TheFlagIsDocumentedAndTheOmittedGroupIsNamed | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestReq134And136And137_TheFlagIsOptInAndNoConsumerPassesItOnLanding | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestReq27And28And95And96And97And129_DefaultModeIsByteIdenticalToTheCapturedPreChangeGolden | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
| TestReq83And112And126_TheOracleBatteryIsTheFiveScenariosAndOwesNoWallTimeBudget | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |

## flow_mvv_0005_test.go (6 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestMVV_AllFourVerbsProveOutOverOneFixtureBackedFlow | internal/cli/flow_mvv_0005_test.go | literal:four verbs |
| TestMVV_FiveScenariosInBothModesAgreeOnExitAndIdentity | internal/cli/flow_mvv_0005_test.go | literal:four verbs |
| TestMVV_NextEvaluateGatesReportsOneGateAndExitsZeroEvenOnADeny | internal/cli/flow_mvv_0005_test.go | literal:four verbs |
| TestMVV_SetMembersWithAngleAndAmpersandSurviveAllThreeHopsByteIdentical | internal/cli/flow_mvv_0005_test.go | literal:four verbs |
| TestMVV_TheUnboundReaderPairDistinguishesTheNarrowedInvokedSet | internal/cli/flow_mvv_0005_test.go | literal:four verbs |
| TestReq132_TheIllustrativeInvocationShapesParseUnderTheShippedGrammar | internal/cli/flow_mvv_0005_test.go | literal:four verbs |

## TestAdv0010 (5 tests, 2 files)
| test | file | arm |
|---|---|---|
| TestAdv0010_DeadEndStaysSilentOverADecisionTable | internal/graphlint/adversarial_0010_test.go | record:0010 |
| TestAdv0010_TerminalEscapeSilentOnARuleFreeDecisionTable | internal/graphlint/adversarial_0010_test.go | record:0010 |
| TestAdv0010_TerminalEscapeStaysSilentOverADecisionTable | internal/graphlint/adversarial_0010_test.go | record:0010 |
| TestAdv0010_ClassDisagreementStillWinsOverAccessorBinding | internal/table/adversarial_0010_test.go | record:0010 |
| TestAdv0010_UndeclaredTagPrecedesTheClassDisagreement | internal/table/adversarial_0010_test.go | record:0010 |

## flow_group_0005_test.go (5 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestFlowGroup_BareInvocationIsAUsageRefusalNotASilentSuccess | internal/cli/flow_group_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestFlowGroup_BareInvocationUnderJSONEmitsTheEnvelopeNotHelpText | internal/cli/flow_group_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestFlowGroup_InvalidOutputModeIsRefusedBeforeTheUsageError | internal/cli/flow_group_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestFlowGroup_TheVerbsThemselvesStillRun | internal/cli/flow_group_0005_test.go | literal:four verbs;literal:flowVerbs |
| TestFlowGroup_UnknownVerbStillRefusesWithTheSameUsageCode | internal/cli/flow_group_0005_test.go | literal:four verbs;literal:flowVerbs |

## TestMVV0023 (1 tests, 1 files)
| test | file | arm |
|---|---|---|
| TestMVV0023_ResolveEnvelopeProjectionEndToEnd | internal/cli/flow_mvv_0023_test.go | literal:four flow verbs;literal:flowVerbs |
