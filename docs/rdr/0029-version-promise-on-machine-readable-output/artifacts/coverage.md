# coverage — 0029 version promise on machine-readable output

Column 1 is the REQ id; column 2 is the test that would fail if a future
change broke that clause. An EMPTY column 2 is the orphan mark.

| REQ | Test |
| --- | --- |
| REQ-1 | `TestReq1And8And42_TheOkEnvelopeCarriesSchemaVersion`, `TestReq1And43_TheRefusalRecordCarriesSchemaVersionAndStaysBare`, `TestReq1And58_SchemaVersionIsAStringOfTheFormMajorMinor` |
| REQ-2 | `TestReq2_TheSchemaVersionIsNotDerivedFromTheBinaryVersion` |
| REQ-3 | `TestReq3And4_TheSchemaVersionConstantHasOneHomeInClierr` |
| REQ-4 | `TestReq3And4_TheSchemaVersionConstantHasOneHomeInClierr` |
| REQ-5 | `TestReq5_SchemaVersionIsNeverProjectedIntoData` |
| REQ-6 | `TestReq6And7_SchemaVersionIsNotOmitemptyAndFindingsStaysTheOneStructuredField` |
| REQ-7 | `TestReq6And7_SchemaVersionIsNotOmitemptyAndFindingsStaysTheOneStructuredField` |
| REQ-8 | `TestReq1And8And42_TheOkEnvelopeCarriesSchemaVersion` |
| REQ-9 | `TestReq9And22_TheIntroductionAndMovementRulesArePublished`, `TestMVV0029_AnAgentPinsAVersionAndSurvivesANewFindingCode` |
| REQ-10 | `TestReq10And44_ThePlanDecoderToleratesTheAddedField` |
| REQ-11 | `TestReq11And12And13And14_TheTierRegisterIsPublishedInTheOutputContract`, `TestReq11_TheRegisterCarriesThe0xIntentQualifierBesideTheTierTable` |
| REQ-12 | `TestReq11And12And13And14_TheTierRegisterIsPublishedInTheOutputContract` |
| REQ-13 | `TestReq11And12And13And14_TheTierRegisterIsPublishedInTheOutputContract`, `TestReq13And25And50_TheCategorySetIsAssertedByMembershipNotCardinality` |
| REQ-14 | `TestReq11And12And13And14_TheTierRegisterIsPublishedInTheOutputContract`, `TestReq14And50_TheAdvisoryTierCarriesItsNamedMembersAndNoDuplicate` |
| REQ-15 | `TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus` |
| REQ-16 | `TestReq16And56_TheEmittedIdentifierCarryingTheWordSurvives` |
| REQ-17 | `TestReq17And63_UntieredVocabulariesAndNonTierSensesKeepTheWord` |
| REQ-18 | `TestReq18And59_ACodeOutsideTheBlockingTierTakesInfoSeverity` |
| REQ-19 | `TestReq19And46_AdvisoryFindingsLeaveTheSuccessDispositionUntouched` |
| REQ-21 | `TestReq21_ACommittedSnapshotOfTheSeamMembersAndSeveritiesIsDiffed`, `TestReq21_TheSnapshotIsComparedAgainstTheCurrentTree` |
| REQ-22 | `TestReq9And22_TheIntroductionAndMovementRulesArePublished` |
| REQ-24 | `TestReq24And31And48_TheEnvelopeTypeSeamIsExactlyOk`, `TestReq24And32And48_TheExitCodeSeamIsExactlyTheFiveEmittedIntegers`, `TestReq24And33And48_TheAdvisoryLevelSeamIsExactlyNoteAndWarning`, `TestReq24And48_TheFrozenAnchorVocabulariesMatchTheirDeclaredMembers` |
| REQ-25 | `TestReq25And34And50_TheUnknownReasonSeamIsAppendOnlyAndDuplicateFree`, `TestReq25And35And50_TheBlockSeamIsAppendOnlyAndDuplicateFree`, `TestReq25And50_TheBlockingTierCarriesItsNamedMembersAndNoDuplicate`, `TestReq25And50_TheReasonSetCarriesItsNamedMembersAndNoDuplicate`, `TestReq13And25And50_TheCategorySetIsAssertedByMembershipNotCardinality` |
| REQ-26 | `TestReq26_The0006AdvisoryTierTextIsAmendedOffItsClosure` |
| REQ-27 | `TestReq27_TheVersionPayloadFieldNamesAreFrozen` |
| REQ-28 | `TestReq28And29_TheUntieredNamespacesGainNoEnumerationSeam` |
| REQ-29 | `TestReq28And29_TheUntieredNamespacesGainNoEnumerationSeam` |
| REQ-30 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist` |
| REQ-31 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq24And31And48_TheEnvelopeTypeSeamIsExactlyOk` |
| REQ-32 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq24And32And48_TheExitCodeSeamIsExactlyTheFiveEmittedIntegers` |
| REQ-33 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq24And33And48_TheAdvisoryLevelSeamIsExactlyNoteAndWarning` |
| REQ-34 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq25And34And50_TheUnknownReasonSeamIsAppendOnlyAndDuplicateFree` |
| REQ-35 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq25And35And50_TheBlockSeamIsAppendOnlyAndDuplicateFree` |
| REQ-36 | `TestReq36_TheCLIErrorCodeVocabularyGainsNoSeam` |
| REQ-37 | `TestReq37And62_TheStaleFailedRecordDocCommentsAreRetired` |
| REQ-38 | `TestReq38And48_BothOperatorEnumerationsArePinnedByValue` |
| REQ-42 | `TestReq1And8And42_TheOkEnvelopeCarriesSchemaVersion` |
| REQ-43 | `TestReq1And43_TheRefusalRecordCarriesSchemaVersionAndStaysBare`, `TestReq43_SchemaVersionRidesEveryProvokedFailureShape` |
| REQ-44 | `TestReq10And44_ThePlanDecoderToleratesTheAddedField` |
| REQ-45 | `TestReq45_TheGoldenIsRecapturedDifferingByExactlyTheOneKey` |
| REQ-46 | `TestReq19And46_AdvisoryFindingsLeaveTheSuccessDispositionUntouched`, `TestMVV0029_AnAgentPinsAVersionAndSurvivesANewFindingCode` |
| REQ-47 | `TestReq47_PromotionToBlockingFlipsTheExitAndReplacesTheRecord` |
| REQ-48 | `TestReq24And31And48_TheEnvelopeTypeSeamIsExactlyOk`, `TestReq24And32And48_TheExitCodeSeamIsExactlyTheFiveEmittedIntegers`, `TestReq24And33And48_TheAdvisoryLevelSeamIsExactlyNoteAndWarning`, `TestReq24And48_TheFrozenAnchorVocabulariesMatchTheirDeclaredMembers`, `TestReq38And48_BothOperatorEnumerationsArePinnedByValue` |
| REQ-49 | `TestReq49_ABlockingRunReturnsTheAggregateCodeAndNoOther` |
| REQ-50 | `TestReq25And34And50_TheUnknownReasonSeamIsAppendOnlyAndDuplicateFree`, `TestReq25And35And50_TheBlockSeamIsAppendOnlyAndDuplicateFree`, `TestReq14And50_TheAdvisoryTierCarriesItsNamedMembersAndNoDuplicate`, `TestReq25And50_TheBlockingTierCarriesItsNamedMembersAndNoDuplicate`, `TestReq25And50_TheReasonSetCarriesItsNamedMembersAndNoDuplicate` |
| REQ-51 | `TestReq51And53_TheCardinalityPinningAssertionsAreReplacedNotUpdated` |
| REQ-52 | `TestReq52_NoGoldenSnapshotOfTheJSONEnvelopeIsAdded` |
| REQ-53 | `TestReq51And53_TheCardinalityPinningAssertionsAreReplacedNotUpdated` |
| REQ-54 | `TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus` |
| REQ-55 | `TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus` |
| REQ-56 | `TestReq16And56_TheEmittedIdentifierCarryingTheWordSurvives` |
| REQ-57 | `TestReq15And54And55And57_TheClosedWordingIsRetiredAcrossTheEnumeratedCensus` |
| REQ-58 | `TestReq1And58_SchemaVersionIsAStringOfTheFormMajorMinor` |
| REQ-59 | `TestReq18And59_ACodeOutsideTheBlockingTierTakesInfoSeverity` |
| REQ-60 | `TestReq60_EveryVerbInheritsTheFieldFromTheGateway` |
| REQ-61 | `TestReq30And31Thru35And61_TheFiveOwedEnumerationSeamsExist`, `TestReq61_TheSeamsAreProductionAccessorsNotTestHelpers` |
| REQ-62 | `TestReq37And62_TheStaleFailedRecordDocCommentsAreRetired` |
| REQ-63 | `TestReq17And63_UntieredVocabulariesAndNonTierSensesKeepTheWord` |
| REQ-MVV | `TestMVV0029_AnAgentPinsAVersionAndSurvivesANewFindingCode` |

## Test files

- `internal/cli/schema_probe_0029_test.go` — the reflective probe layer over
  surface that does not exist yet (the constant, the five owed seams).
- `internal/cli/schema_envelope_0029_test.go` — the wire observations on both
  terminal records.
- `internal/cli/schema_tiers_0029_test.go` — frozen by-value and append-only
  membership over C4's census.
- `internal/cli/schema_closed_wording_0029_test.go` — S9's enumerated
  `closed`-retirement census and the stale doc comments.
- `internal/cli/schema_register_0029_test.go` — C2's published register and
  C3's committed snapshot.
- `internal/cli/schema_mvv_0029_test.go` — the golden re-capture, the
  promotion control, and the REQ-MVV runner.
- `internal/graphlint/tier_assertions_0029_test.go` — S8's three replacements
  and the C3 severity split at the engine boundary.

## Orphans (empty column 2)

REQ-20, REQ-23, REQ-39, REQ-40, REQ-41 — each is named in `next_action` as a
demotion with its reason. None is covered by a vacuous test.

## Note on probe shape

`clierr.SchemaVersion` and the five owed seams (`respond::Types`,
`respond::Levels`, `clierr::ExitCodes`, `cli::UnknownReasons`,
`resolve::Blocks`) do not exist at `2f7d22c`. The pre-commit hook runs
`gofmt -l .` and `go vet ./...` under `set -e`, so a suite naming them
directly could not be committed. Every probe therefore reaches them through
`go/parser` over the package source or through the emitted envelope's own
bytes: each red is a runtime assertion failure, never a compile break.

Eight tests initially passed against the absent surface and were rewritten
before commit — each now fails first on a precondition asserting the subject
exists, so none can pass tautologically.
