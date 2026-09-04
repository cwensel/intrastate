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
- `internal/cli/edit_envelope_0028_test.go` — the CLI ENVELOPE half of
  C1.3's exit-group clause (added in Phase 2; see the orphan notes).

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
| REQ-43 | TestReq43And113_ADeclaredLineEditRefusalTakesTheExit2Group; TestReq43_TheDiscriminatorLeavesGenuineEnvironmentFailuresAtExit3; TestReq43_TheEditBindingIsReachedAndRewritesTheAnchoredLine |
| REQ-44 | `TestReq43And113_ADeclaredLineEditRefusalTakesTheExit2Group` + `TestReq114_TheRefusalCarriesTheRuleAndTokenInFindings` (negative: property asserted instead) |
| REQ-45 | `TestReq45_TheExecutorTimeoutArmDoesNotWrapTheRequestSentinel` |
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
| REQ-102 | `TestReq102_TheOutputContractNamesBothEditNameSets` (cli) |
| REQ-103 | `TestReq43_TheEditBindingIsReachedAndRewritesTheAnchoredLine` (negative: scope control) |
| REQ-104 | TestReq104_ThePinnedAnchorSelectsExactlyOneLinePerInScopeFile |
| REQ-105 | TestReq105_AStatusValueOnTheContinuationLineIsUnmatchedNotRewritten |
| REQ-106 | TestReq16And106_ABoundTagValueIsRegexpQuotedAndCannotAlterThePattern |
| REQ-107 | TestReq49And107_ANoOpEditIsNotWrittenAtAllWitnessedByTheInode |
| REQ-108 | TestReq108_TheTargetsModeIsPreservedAcrossTheWrite |
| REQ-109 | TestReq33And109_ASymlinkedTargetIsResolvedBeforeStagingAndSurvives |
| REQ-110 | TestReq58And110_AnEditAppliesWithTheGateOffAndInvocationsCountsApplies |
| REQ-111 | TestReq52And111_ABoundButUnreferencedTagValueIsNeverScanned; TestReq111_AFlagShapedTagValueIsAdmittedInAnAnchor |
| REQ-112 | TestReq59And60And112_AGateOffCommandReaderRefusesTheWriteBeforeMutation |
| REQ-113 | TestReq43And113_ADeclaredLineEditRefusalTakesTheExit2Group |
| REQ-114 | TestReq114_TheRefusalCarriesTheRuleAndTokenInFindings |
| REQ-115 | TestReq32And115_AnUnreadableTargetRefusesWithTheOSErrorNotAnAnchorResult |
| REQ-116 | TestReq71And116_TheEarlierCategoryInRegistrationOrderIsTheOneReported; TestReq72And116_SiblingTablesReportOneStableCategoryAndNoNamedWinner |
| REQ-117 | TestReq80And117_ClearIsOneWayAndTheNextWriteRefusesUnmatched |
| REQ-118 | TestReq11And12And118_ACarrierLessEntryBuildsARefusingBindingNeverAFileBinding |
| REQ-119 | TestReq119_AWrongLineRewriteIsCaughtByReadBackMismatchWithAppliedFalse |
| REQ-120 | TestReq51To54And120_ApplyTimeRefusalsFailFastInPrecedenceOrder |
| REQ-121 | TestReq25And28And121_AnchorBracesReachRE2Untouched; TestReq22And121_ReplaceEscapesAndDefinedGroupsAreAdmitted; TestReq22And121_ReplaceEscapesAndNonParticipatingGroupsExpandAsStated |
| REQ-122 | TestMVV_DeclaredLineEditWriter/step3_the_reader_splits_the_qualifier_and_reports_status_final |
| REQ-123 | `TestReq106_EveryNamedCategoryExistsAndIsWitnessed` (leg 2) + this file (legs 1, 3) |
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

- **REQ-43, REQ-113, REQ-114** — CLOSED in Phase 2 by
  `internal/cli/edit_envelope_0028_test.go`. The three were unreachable
  from `internal/table/` or `internal/accessor/`, which were Phase 1's two
  test directories; the accessor-seam half of S27 (class
  `execution_failure`, `Applied()` false, the rule id and reason token in
  the Detail) was already covered at REQ-40/41/42.

  The envelope test pins the exit GROUP, the discriminator and the
  `findings[]` shape, and NEVER the CLI code string, which REQ-44 forbids
  and REQ-113 restates. Its load-bearing arm is the negative control: an
  unreadable target is the SAME `execution_failure` class reached WITHOUT
  the request bit and must still exit 3, so an implementation that moved
  the whole class to exit 2 fails rather than passing every positive
  assertion.
- **REQ-45** — CLOSED at the completion gate, no longer an orphan. The
  original note was half right: asserting the applied sense here WOULD
  restate 0004 (`TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense`
  already pins it). But the REQ also bounds THIS record's own change —
  ADV-1/ADV-2 widened `ErrDeclaredRequest` to the entry-level
  preconditions, and REQ-45 is what says the executor's deadline arm is
  not among them. That is a property only 0028 can break, so it earns a
  0028 test:
  `TestReq45_TheExecutorTimeoutArmDoesNotWrapTheRequestSentinel`.
  Mutation-verified — wrapping the sentinel at `executor.go`'s timeout arm
  makes it fail.

- **REQ-44** — a pure prohibition with no observable behaviour: it forbids
  a test from pinning the CLI code string, which is a constraint ON the
  envelope tests rather than a test of its own. Honoured by
  `TestReq43And113_…`, which pins the exit GROUP and the `findings[]`
  shape and never the spelling.
- **REQ-102** — CLOSED at the completion gate, no longer an orphan. The
  doc named neither name set: Phase 2 documented the exit-group behaviour
  but not the twelve wire strings the REQ requires. Both DISJOINT sets of
  six are now written into `docs/cli-output-contract.md` — C1.4's six
  load-time categories and C1.3/C1.5's six apply-time reason tokens, with
  the disjointness stated — and pinned by
  `TestReq102_TheOutputContractNamesBothEditNameSets`, which also asserts
  the disjointness holds of the shipped `table.Categories()` so the doc
  cannot document a distinction the code stops making. Verified
  non-tautological: it fails against the pre-edit doc.
- **REQ-103** — CLOSED as a negative REQ under the predecessor's own
  convention (`0025` coverage.md, "How the negative REQs are covered"): a
  clause that FORBIDS something is covered by a test of the property it
  protects, plus a positive control. The property is that the carrier
  works against a fixture the TEST builds, so no consumer `rdr-write.toml`
  is needed for it to pass — verified: no such file exists anywhere in the
  tree, and the suite is green. The positive control is that the carrier
  is genuinely live, or the claim would be vacuous.

## How the negative REQs are covered

Two clauses FORBID rather than require, so neither has a direct oracle.
Each is covered by a test of the property it protects, paired with a
positive control so the assertion cannot pass vacuously — the convention
RDR 0025's coverage.md sets under this same heading.

| REQ | Property asserted instead | Positive control |
| --- | --- | --- |
| REQ-44 | the envelope tests pin the exit GROUP and the `findings[]` SHAPE, never the CLI code's spelling — grep the file: no test compares a code string literal | the refusal must actually reach the envelope and be routed, or the claim is vacuous (`TestReq43_TheDiscriminatorLeavesGenuineEnvironmentFailuresAtExit3` is the contrasting arm) |
| REQ-103 | the carrier is exercised entirely from test-built fixtures; no consumer `rdr-write.toml` is added or required, and none exists in the tree | the carrier must be live end to end (`TestReq43_TheEditBindingIsReachedAndRewritesTheAnchoredLine` rewrites a real anchored line) |
- **REQ-123** — the Done statement, now PARTLY pinned rather than wholly
  unwitnessed. It has three legs. Leg 2 — "every refusal category fires on
  a fixture that earns it" — is a real testable property and is carried by
  the repo's existing floor test,
  `internal/table/dump_test.go::TestReq106_EveryNamedCategoryExistsAndIsWitnessed`,
  which all six `edit_*` categories are registered in with a checked-in
  `testdata/neg/*.toml` witness each. That test asserts BOTH directions,
  so a seventh category shipping without a witness fails it. Legs 1
  ("every C1–C1.6 clause has a test") and 3 ("the MVV's steps pass end to
  end") are the phase-exit criterion this file and the suite satisfy, not
  behaviours a test can assert about itself. Its "seven MVV steps" is the
  stale count req-list A-7 records; eight run.

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

## Phase 2 correction to the section above

`TestReq19And101_TheEditCarrierIsCarriedOnTheDumpSurface` was NOT on the
green-by-design list and should have been on a third list: it passed
TAUTOLOGICALLY once the type surface existed. Its oracle is
`strings.Contains(table.Dump(m), "edit")`, and the fixture's model id is
`editflow`, so the assertion was satisfied by the `MODEL editflow …`
header alone while `Dump` emitted no accessor at all. Phase 2 found the
dump of an `edit`-carried model byte-identical to that of a carrier-less
one and implemented the carrier line REQ-101/A6 requires; deviations D4
records the format.

Phase 2 moved the EMITTER but not the ORACLE, so this section's original
claim that "the test now discriminates" was false as written: with the
carrier line shipped, the assertion was still `Contains(dumped, "edit")`
and still passed with every `ACCESSOR` line stripped from the dump. The
oracle was rewritten afterwards to walk the declared writers and require,
per declared rule, a FULL-LINE match of D4's carrier line carrying that
rule's anchor, replace and clear, plus an emitted-line COUNT equal to the
declared rule count. Confirmed by mutation: deleting the
`writeEditCarriers` call, skipping keys inside its inner loop, and
rendering the wrong field into `replace=` each turn the test red. The
fixture id stays `editflow` — an oracle that needed it renamed would be
pinning the fixture, not the surface.

That rewrite also carried a two-dump byte-identity check for `0002:C19`'s
sorted walk, and IT was vacuous in the same way for the same reason: the
fixture declares one writer with one key, and with a single element there
is no order to get wrong, so both `slices.Sorted` calls could be deleted
with the check still green. The ordering guarantee now has its own fixture
(`editOrderModel`) — four writer ids and four keys inside one of them, all
declared in the reverse of sorted order — and the oracle names the
expected canonical SEQUENCE rather than comparing two dumps to each other.
Self-comparison is the weaker property: Go randomizes map iteration per
process, not per range, so an unsorted walk can hand back the same order
twice and satisfy it.

Fixture WIDTH turned out to be load-bearing, not cosmetic. At two elements
per dimension each ordering mutant reddened only ~34 of 40 runs, because an
unsorted range over a two-element map reproduces sorted order often enough
that a real regression would have survived CI about one run in seven. At
four it is 40 of 40 for both — dropping the writer-id sort and dropping the
key sort, measured separately.

The lesson generalises three times over. A `Contains` oracle over a whole
rendered surface can be satisfied by an unrelated substring of that
surface, and a fixture id is the likeliest source. Shipping the surface a
tautological test was meant to guard does not by itself repair the test —
the oracle has to move with it, or the claim of coverage is the only thing
that changed. And an assertion is only as discriminating as the fixture it
runs against: an ordering, uniqueness, or precedence property asserted over
a one-element collection cannot fail, so the fixture has to be sized to the
property before the assertion means anything. Each of the three was found
in the fix for the one before it.

The other sixty-five were confirmed red behaviourally: run against a
compiling no-op `edit` binding (a stub `NewEditWriter` whose `Apply`
returns nil, an `Artifact.Context` field, a `Registry.AllowCommands`
field, a `table.EditRule` type and the six `Cat…` constants), all
sixty-five fail. Against the real tree they fail at COMPILE, which is a
stronger red but a weaker signal — the stub run is the evidence that no
test is tautological.

## REQ-MVV runner

The MVV runs at the ACCESSOR SEAM, not through the CLI pipeline.
REQ-MVV-BLOCK records why: the consumer's row-addressed projector verb the
C1.6 README reader invokes does not ship, so `flow resolve | flow
set-state` cannot run end to end. `TestReqMVVBlock_…` probes for that verb
and stays green while it is absent — it flips RED the day the consumer
ships it, which is the signal that Phase 4's CLI half can run.

All eight steps are unblocked at the seam, which is what "Phases 1–3 do not
depend on it" means. The executor, the two `edit` write bindings and a role
reader per artifact are all real; only the projector is a stub, and it
reads the FILE back rather than mirroring what the writer intended, so each
read-back is a genuine re-read that can fail independently of its write.

```sh
go test ./internal/accessor/ \
  -run 'TestMVV_DeclaredLineEditWriter|TestReqMVVBlock' -v
```

## REQ-MVV output

Actual output, 2026-09-04, against the merged tree.

```
--- PASS: TestMVV_DeclaredLineEditWriter (0.02s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step1_the_model_lints_clean_with_no_wrapper_script (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file (0.01s)
        --- PASS: TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file/record_0027 (0.00s)
        --- PASS: TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file/record_0026 (0.00s)
        --- PASS: TestMVV_DeclaredLineEditWriter/step2and3_three_sequential_invocations_each_flip_one_line_per_file/record_0022 (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step3_the_reader_splits_the_qualifier_and_reports_status_final (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step4_a_none_or_stopped_row_applies_nothing_and_does_not_refuse (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step5_without_allow_commands_the_write_refuses_before_mutation (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step6_a_duplicated_readme_row_refuses_edit_anchor_ambiguous (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step7_a_self_de_anchoring_replace_refuses_edit_anchor_unstable (0.00s)
    --- PASS: TestMVV_DeclaredLineEditWriter/step8_an_unlinked_row_refuses_unmatched_while_siblings_still_flip (0.01s)
        --- PASS: TestMVV_DeclaredLineEditWriter/step8_an_unlinked_row_refuses_unmatched_while_siblings_still_flip/the_unlinked_row_refuses_and_names_the_rule (0.00s)
        --- PASS: TestMVV_DeclaredLineEditWriter/step8_an_unlinked_row_refuses_unmatched_while_siblings_still_flip/the_other_rows_in_the_same_readme_still_flip (0.00s)
--- PASS: TestReqMVVBlock_ThePhase4PipelineAwaitsTheConsumersRowAddressedVerb (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/accessor	0.207s
```

Step by step, and what each proves at this seam:

| step | REQ | asserted |
| --- | --- | --- |
| 1 | REQ-MVV.1 | the model lints clean; every writer carries `edit` and neither a `command` nor a `path` — no wrapper script anywhere in the fixture |
| 2+3 | REQ-MVV.2, .3 | three sequential invocations, one per record identity; each flips exactly ONE line per file, byte-compared against a hand-written buffer. The joint record keeps its bracketed qualifier by backreference; the wrapped record's continuation line comes through byte-identical |
| 3 | REQ-122 | the reader reports `status=Final` with the qualifier off the value |
| 4 | REQ-MVV.4 | an empty plan applies nothing and does not refuse; both inodes unmoved |
| 5 | REQ-MVV.5 | with the gate OFF and both readers command-backed, both writes refuse BEFORE mutation — not `read_back_incomplete`, `Applied()` false, both files byte-identical |
| 6 | REQ-MVV.6 | a duplicated README row refuses `edit_anchor_ambiguous`; inode unmoved |
| 7 | REQ-MVV.7 | a self-de-anchoring `replace` refuses `edit_anchor_unstable` before any write |
| 8 | REQ-MVV.8 | an unlinked `| NNNN |` row refuses `edit_anchor_unmatched` and NAMES the rule; the other rows of the SAME README still flip |

What is deferred to Phase 4 is the ENVELOPE, not the behaviour: the
dispositions, byte invariants and refusal tokens above are the record's own
and are asserted here. The CLI envelope's own half — the exit group, the
discriminator and the `findings[]` carriage — is covered separately at
REQ-43/113/114, which no longer depend on the missing verb.

