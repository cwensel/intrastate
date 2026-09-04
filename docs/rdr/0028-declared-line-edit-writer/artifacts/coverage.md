# Coverage — RDR 0028 declared-line-edit-writer

Phase 1 test authoring. Column 1 is the REQ id from `req-list.md`; column 2
is the test that holds it. An EMPTY column-2 cell is the orphan mark.

Test files:

- `internal/table/edit_carrier_0028_test.go` — the LOAD half (C1.1 grammar,
  C1.2 admission as decided at load, C1.4 categories, C1.6 argv rules).
- `internal/accessor/edit_apply_0028_test.go` — the APPLY half (C1.2 value
  shape and parse-once, C1.3 execution, C1.5 clear, C1.6 binding, A1 seam).
- `internal/accessor/edit_mvv_0028_test.go` — REQ-MVV's eight steps and
  REQ-MVV-BLOCK.

| REQ | Test |
|---|---|
| REQ-1 | TestReq1_EditTablesDecodePerDeclaredKey |
| REQ-2 | TestReq2And3And14_TheRuleFieldsAreAnchorReplaceAndClear; TestReq2And3_ReplaceEmitsTheWholeLineNotTheMatchedSpan |
| REQ-3 | TestReq2And3And14_TheRuleFieldsAreAnchorReplaceAndClear; TestReq2And3_ReplaceEmitsTheWholeLineNotTheMatchedSpan |
| REQ-4 | TestReq4_ClearIsOptionalAndCarriesItsDeclaredDisposition |
| REQ-5 | TestReq5_RoleKeysAndTimeoutAreUnchangedByTheEditCarrier |
| REQ-6 | TestReq6And8_AnEditOnlyEntryLoadsAndIsNotCarrierLess |
| REQ-7 | TestReq7And10And63_EditBesideAnotherCarrierOrOffWriteIsCarrierConflict |
| REQ-8 | TestReq6And8_AnEditOnlyEntryLoadsAndIsNotCarrierLess |
| REQ-9 | TestReq9_TheNeitherArmKeepsItsWireStringAndNamesThreeCarriers |
| REQ-10 | TestReq7And10And63_EditBesideAnotherCarrierOrOffWriteIsCarrierConflict |
| REQ-11 | TestReq11And12And118_ACarrierLessEntryBuildsARefusingBindingNeverAFileBinding |
| REQ-12 | TestReq11And12And118_ACarrierLessEntryBuildsARefusingBindingNeverAFileBinding |
| REQ-13 | TestReq13_TheCarrierDoesNotEnterTheAccessorIdentityAndARuleIsEntryPlusKey |
| REQ-14 | TestReq2And3And14_TheRuleFieldsAreAnchorReplaceAndClear |
| REQ-15 | TestReq15And24And66_ReplaceVocabularyIsClosed |
| REQ-16 | TestReq16And106_ABoundTagValueIsRegexpQuotedAndCannotAlterThePattern |
| REQ-17 | TestReq17And26And65_AnchorDefectsAreEditAnchorInvalid |
| REQ-18 | TestReq18_AnEntrysOwnPlannedKeyIsNotAnAnchorPlaceholder |
| REQ-19 | TestReq19And101_TheEditCarrierIsCarriedOnTheDumpSurface; TestReq19And23_CapturedTextAndValuesEmitAsLiteralBytesNeverRescanned |
| REQ-20 | TestReq20_AValueCarryingANewlineOrCarriageReturnRefusesBeforeMutation |
| REQ-21 | TestReq21And29_OnlyNewlinesAreRefusedAndNoShellHazardIsScannedFor |
| REQ-22 | TestReq22And121_ReplaceEscapesAndDefinedGroupsAreAdmitted; TestReq22And121_ReplaceEscapesAndNonParticipatingGroupsExpandAsStated |
| REQ-23 | TestReq19And23_CapturedTextAndValuesEmitAsLiteralBytesNeverRescanned |
| REQ-24 | TestReq15And24And66_ReplaceVocabularyIsClosed |
| REQ-25 | TestReq25And28And121_AnchorBracesReachRE2Untouched |
| REQ-26 | TestReq17And26And65_AnchorDefectsAreEditAnchorInvalid |
| REQ-27 | TestReq27_AnEscapedTagPrefixIsPatternTextNotAPlaceholder |
| REQ-28 | TestReq25And28And121_AnchorBracesReachRE2Untouched |
| REQ-29 | TestReq21And29_OnlyNewlinesAreRefusedAndNoShellHazardIsScannedFor |
| REQ-30 | TestReq30And56And57_TerminatorHandlingPreservesLineShape |
| REQ-31 | TestReq31And47_AnAppliedEditRewritesExactlyTheAnchoredLine; TestReq30And56And57_TerminatorHandlingPreservesLineShape |
| REQ-32 | TestReq32And115_AnUnreadableTargetRefusesWithTheOSErrorNotAnAnchorResult |
| REQ-33 | TestReq33And109_ASymlinkedTargetIsResolvedBeforeStagingAndSurvives |
| REQ-34 | TestReq34And100_TheEditWriteEmitsNoSealAndCreatesNoSidecar |
| REQ-35 | TestReq35And36And37_EachRuleMustSelectExactlyOneLine |
| REQ-36 | TestReq35And36And37_EachRuleMustSelectExactlyOneLine |
| REQ-37 | TestReq35And36And37_EachRuleMustSelectExactlyOneLine |
| REQ-38 | TestReq38_AReplacementThatDeAnchorsItselfRefusesEditAnchorUnstable |
| REQ-39 | TestReq39_ReAnchorAssertsLineIdentityNotMerelyCardinality |
| REQ-40 | TestReq40And41And42_ARefusedEditIsExecutionFailureNotAppliedWithARuleScopedDetail |
| REQ-41 | TestReq40And41And42_ARefusedEditIsExecutionFailureNotAppliedWithARuleScopedDetail |
| REQ-42 | TestReq40And41And42_ARefusedEditIsExecutionFailureNotAppliedWithARuleScopedDetail |
| REQ-43 |  |
| REQ-44 |  |
| REQ-45 |  |
| REQ-46 | TestReq46And48_AtomicityIsOneEntryOverOneFileAndNeverCrossEntry |
| REQ-47 | TestReq31And47_AnAppliedEditRewritesExactlyTheAnchoredLine |
| REQ-48 | TestReq46And48_AtomicityIsOneEntryOverOneFileAndNeverCrossEntry |
| REQ-49 | TestReq49And107_ANoOpEditIsNotWrittenAtAllWitnessedByTheInode |
| REQ-50 | TestReq50_ThereIsNoLockAndNoCompareBeforeRename |
| REQ-51 | TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder |
| REQ-52 | TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder; TestReq52And111_ABoundButUnreferencedTagValueIsNeverScanned |
| REQ-53 | TestReq53And79_AnUndeclaredClearRefusesAheadOfSelection; TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder |
| REQ-54 | TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder |
| REQ-55 | TestReq55_TwoEquallyDefectiveSiblingsReportAStableTokenAndNoPinnedWinner |
| REQ-56 | TestReq30And56And57_TerminatorHandlingPreservesLineShape |
| REQ-57 | TestReq30And56And57_TerminatorHandlingPreservesLineShape |
| REQ-58 | TestReq58And110_AnEditAppliesWithTheGateOffAndInvocationsCountsApplies |
| REQ-59 | TestReq59And60And112_AGateOffCommandReaderRefusesTheWriteBeforeMutation |
| REQ-60 | TestReq59And60And112_AGateOffCommandReaderRefusesTheWriteBeforeMutation |
| REQ-61 | TestReq61And62_TheGateIsCarriedAsStateOnTheAccessorRegistry |
| REQ-62 | TestReq61And62_TheGateIsCarriedAsStateOnTheAccessorRegistry |
| REQ-63 | TestReq7And10And63_EditBesideAnotherCarrierOrOffWriteIsCarrierConflict |
| REQ-64 | TestReq64_KeysAndEditTablesMustBeInBijection |
| REQ-65 | TestReq17And26And65_AnchorDefectsAreEditAnchorInvalid |
| REQ-66 | TestReq15And24And66_ReplaceVocabularyIsClosed |
| REQ-67 | TestReq67And81_ClearIsClosedAtTheSingleLiteralLine |
| REQ-68 | TestReq68And70And86_ADeclaredTagAtArgv0IsEditTagArgv0 |
| REQ-69 | TestReq69And76_TheSixEditCategoriesFollow0025sSixInClauseOrder; TestReq69And73_TheSixWireStringsRegisterAndTheApplyTokensDoNot |
| REQ-70 | TestReq68And70And86_ADeclaredTagAtArgv0IsEditTagArgv0 |
| REQ-71 | TestReq71And116_TheEarlierCategoryInRegistrationOrderIsTheOneReported |
| REQ-72 | TestReq72And116_SiblingTablesReportOneStableCategoryAndNoNamedWinner |
| REQ-73 | TestReq69And73_TheSixWireStringsRegisterAndTheApplyTokensDoNot |
| REQ-74 | TestReq74_AModelMeetingEverySixLintObligationLoadsClean |
| REQ-75 | TestReq75_LintDoesNotConsultArtifactsOrBindings |
| REQ-76 | TestReq69And76_TheSixEditCategoriesFollow0025sSixInClauseOrder; TestReq77_TheSixCommandCategoriesAreRegisteredContiguouslyInClauseOrder (D1 rewrite) |
| REQ-77 | TestReq77_ClearLineDeletesTheAnchoredLineWithItsTerminator |
| REQ-78 | TestReq78_AZeroMatchClearSucceedsWithNoWriteButTwoMatchesStillRefuses |
| REQ-79 | TestReq53And79_AnUndeclaredClearRefusesAheadOfSelection |
| REQ-80 | TestReq80And117_ClearIsOneWayAndTheNextWriteRefusesUnmatched |
| REQ-81 | TestReq67And81_ClearIsClosedAtTheSingleLiteralLine |
| REQ-82 | TestReq82_TheClearSentinelIsReadOnThePlannedValueAndNeverOnATemplate |
| REQ-83 | TestReq83And84_ADeclaredTagPlaceholderIsAdmittedAtANonLeadingPosition |
| REQ-84 | TestReq83And84_ADeclaredTagPlaceholderIsAdmittedAtANonLeadingPosition |
| REQ-85 | TestReq85_AnUndeclaredOrUnrecognizedBraceElementStaysCommandUnknownPlaceholder |
| REQ-86 | TestReq68And70And86_ADeclaredTagAtArgv0IsEditTagArgv0 |
| REQ-87 | TestReq87_ArtifactAtArgv0StaysAdmitted |
| REQ-88 | TestReq88And89And90_UnboundAndFlagShapedTagValuesRefuseBeforeSpawn |
| REQ-89 | TestReq88And89And90_UnboundAndFlagShapedTagValuesRefuseBeforeSpawn |
| REQ-90 | TestReq88And89And90_UnboundAndFlagShapedTagValuesRefuseBeforeSpawn |
| REQ-91 | TestReq91_ABoundButUnreferencedFlagShapedTagIsNeverScannedAtArgv |
| REQ-92 | TestReq92And93_ThereIsNoEditReadCarrierAndTheReadSideIsUnchanged |
| REQ-93 | TestReq92And93_ThereIsNoEditReadCarrierAndTheReadSideIsUnchanged |
| REQ-94 | TestReq94_TheShellInterpreterDenyListIsNeitherWidenedNorNarrowed |
| REQ-95 | TestReq95To98_OneContextMapOnArtifactServesBothTheAnchorAndTheReader |
| REQ-96 | TestReq95To98_OneContextMapOnArtifactServesBothTheAnchorAndTheReader |
| REQ-97 | TestReq95To98_OneContextMapOnArtifactServesBothTheAnchorAndTheReader |
| REQ-98 | TestReq95To98_OneContextMapOnArtifactServesBothTheAnchorAndTheReader |
| REQ-99 | TestReq99_TheNonOwnedCheckRunsBeforeApplyAndTheBindingIsNeverReached |
| REQ-100 | TestReq34And100_TheEditWriteEmitsNoSealAndCreatesNoSidecar |
| REQ-101 | TestReq19And101_TheEditCarrierIsCarriedOnTheDumpSurface |
| REQ-102 |  |
| REQ-103 |  |
| REQ-104 | TestReq104_ThePinnedAnchorSelectsExactlyOneLinePerInScopeFile |
| REQ-105 | TestReq105_AStatusValueOnTheContinuationLineIsUnmatchedNotRewritten |
| REQ-106 | TestReq16And106_ABoundTagValueIsRegexpQuotedAndCannotAlterThePattern |
| REQ-107 | TestReq49And107_ANoOpEditIsNotWrittenAtAllWitnessedByTheInode |
| REQ-108 | TestReq108_TheTargetsModeIsPreservedAcrossTheWrite |
| REQ-109 | TestReq33And109_ASymlinkedTargetIsResolvedBeforeStagingAndSurvives |
| REQ-110 | TestReq58And110_AnEditAppliesWithTheGateOffAndInvocationsCountsApplies |
| REQ-111 | TestReq52And111_ABoundButUnreferencedTagValueIsNeverScanned; TestReq111_AFlagShapedTagValueIsAdmittedInAnAnchor |
| REQ-112 | TestReq59And60And112_AGateOffCommandReaderRefusesTheWriteBeforeMutation |
| REQ-113 |  |
| REQ-114 |  |
| REQ-115 | TestReq32And115_AnUnreadableTargetRefusesWithTheOSErrorNotAnAnchorResult |
| REQ-116 | TestReq71And116_TheEarlierCategoryInRegistrationOrderIsTheOneReported; TestReq72And116_SiblingTablesReportOneStableCategoryAndNoNamedWinner |
| REQ-117 | TestReq80And117_ClearIsOneWayAndTheNextWriteRefusesUnmatched |
| REQ-118 | TestReq11And12And118_ACarrierLessEntryBuildsARefusingBindingNeverAFileBinding |
| REQ-119 | TestReq119_AWrongLineRewriteIsCaughtByReadBackMismatchWithAppliedFalse |
| REQ-120 | TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder |
| REQ-121 | TestReq25And28And121_AnchorBracesReachRE2Untouched; TestReq22And121_ReplaceEscapesAndDefinedGroupsAreAdmitted; TestReq22And121_ReplaceEscapesAndNonParticipatingGroupsExpandAsStated |
| REQ-122 | TestMVV_DeclaredLineEditWriter/step3_the_reader_splits_the_qualifier_and_reports_status_final |
| REQ-123 |  |
| REQ-MVV | TestMVV_DeclaredLineEditWriter |
| REQ-MVV.1 | TestMVV_DeclaredLineEditWriter/step1_the_model_lints_clean_with_no_wrapper_script |
| REQ-MVV.2 | TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file |
| REQ-MVV.3 | TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file; TestMVV_DeclaredLineEditWriter/step3_the_reader_splits_the_qualifier_and_reports_status_final |
| REQ-MVV.4 | TestMVV_DeclaredLineEditWriter/step4_a_none_or_stopped_row_applies_nothing_and_does_not_refuse |
| REQ-MVV.5 | TestMVV_DeclaredLineEditWriter/step5_without_allow_commands_the_write_refuses_before_mutation |
| REQ-MVV.6 | TestMVV_DeclaredLineEditWriter/step6_a_duplicated_readme_row_refuses_edit_anchor_ambiguous |
| REQ-MVV.7 | TestMVV_DeclaredLineEditWriter/step7_a_self_de_anchoring_replace_refuses_edit_anchor_unstable |
| REQ-MVV.8 | TestMVV_DeclaredLineEditWriter/step8_an_unlinked_row_refuses_unmatched_while_siblings_still_flip |
| REQ-MVV-BLOCK | TestReqMVVBlock_ThePhase4PipelineAwaitsTheConsumersRowAddressedVerb |

## Orphan notes

Eight REQs carry an empty cell. Each is unreachable from
`internal/table/` or `internal/accessor/`, which are this phase's two test
directories; none is a spec gap.

- **REQ-43, REQ-113, REQ-114** — the CLI ENVELOPE. C1.3's exit-2 group, the
  `flow_exec.go::accessorFailureOf` discriminator and the `findings[]`
  carriage per `docs/cli-output-contract.md` are all observable only at
  `internal/cli/`. The accessor-seam half of S27 (class
  `execution_failure`, `Applied()` false, the rule id and reason token in
  the Detail) IS covered — see REQ-40/41/42. Phase 2 or 3 owes an
  `internal/cli/*_0028_test.go` for the envelope half. REQ-44 forbids
  pinning the CLI code STRING and REQ-113 restates it; the envelope test
  must pin the exit GROUP and the `findings[]` shape, never the spelling.
- **REQ-44, REQ-45** — pure prohibitions with no observable behaviour in
  these packages. REQ-44 forbids a test from pinning the CLI code string,
  which is a constraint ON the envelope test above rather than a test of
  its own. REQ-45 scopes `internal/accessor/executor.go`'s deadline arm OUT
  of this RDR and asserts nothing about it; writing a test would assert
  0004's behaviour under 0028's id.
- **REQ-102** — `docs/cli-output-contract.md` naming the new categories and
  Detail tokens. A DOCS obligation, Phase 4's. The repo has a docs gate
  (`internal/cli/docs_test.go`) that is the natural home once the
  categories exist.
- **REQ-103** — the consumer's `rdr-write.toml` migration, explicitly
  out of scope ("the consumer's kata rdr#yjye, not this plan's").
- **REQ-123** — the Done statement. It is the phase-exit criterion this
  file and the suite together satisfy, not a testable behaviour. Note its
  "seven MVV steps" is the stale count req-list A-7 records; eight run.

## Green-by-design tests

Five of the seventy tests pass against the pre-implementation tree once
the type surface exists, because they assert EXISTING or NEGATIVE
behaviour rather than new `edit`-carrier behaviour. Called out so a
reviewer does not read them as tautologies.

- `TestReq2And3And14_TheRuleFieldsAreAnchorReplaceAndClear` — the renamed
  field spellings are refused by 0002's existing strict decode. The test
  pins that `anchor`/`replace`/`clear` are the grammar's field names.
- `TestReq61And62_TheGateIsCarriedAsStateOnTheAccessorRegistry` — a
  structural assertion that the gate is a FIELD on `accessor.Registry`.
  It cannot compile until the field exists, and it then passes: the
  contract it holds is architectural (no `cmdbind` import, no second
  policy site), which a behavioural assertion cannot reach.
- `TestReq94_TheShellInterpreterDenyListIsNeitherWidenedNorNarrowed` — a
  pure negative over 0027's registered category. It is green today and
  goes red only if this record mints a `{tag.…}` deny-list variant, which
  is exactly the regression it exists to catch.
- `TestReq99_TheNonOwnedCheckRunsBeforeApplyAndTheBindingIsNeverReached` —
  the phase ORDER at `Executor.Write` is 0004's, unchanged. The test pins
  that the `edit` binding does not perturb it.
- `TestReqMVVBlock_ThePhase4PipelineAwaitsTheConsumersRowAddressedVerb` —
  asserts the blocker REQ-MVV-BLOCK names still holds. It flips red when
  the consumer ships the row-addressed projector verb, which is the signal
  that MVV Phase 4's CLI pipeline can run.

The other sixty-five were confirmed red behaviourally: run against a
compiling no-op `edit` binding (a stub `NewEditWriter` whose `Apply`
returns nil, an `Artifact.Context` field, a `Registry.AllowCommands`
field, a `table.EditRule` type and the six `Cat…` constants), all
sixty-five fail. Against the real tree they fail at COMPILE, which is a
stronger red but a weaker signal — the stub run is the evidence that no
test is tautological.
