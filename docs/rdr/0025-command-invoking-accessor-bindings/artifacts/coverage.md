# Coverage — RDR 0025 command-invoking-accessor-bindings

Phase 1 (tests first). REQ-N x test name; orphans flagged BOTH ways.

Source of REQs: `req-list.md` — 142 REQs plus REQ-MVV.
Suite: `go test ./...`. 79 new tests across 7 test files plus 2 fixture files.

## Orphans

- **A REQ with no test: NONE.** All 142 REQs and REQ-MVV are cited.
- **A test citing no REQ: NONE.** All 79 new tests open with at least one
  `// REQ-N: "<quote>"` header line.

Derived mechanically: for each `func TestReq…`, the contiguous comment block
above it supplies the header citations and its body supplies in-body ones;
both directions were then checked against the full REQ-1..142 + REQ-MVV set.

## Red gate

79 new tests; 79 `--- FAIL` lines tree-wide on `go test ./...`. Every failing
test in the tree is an RDR 0025 test and no pre-existing test regressed
(verified by diffing the failing-test name list against the suite).

The tests COMPILE and fail on BEHAVIOUR rather than on a missing symbol. The
spec-named surface is declared, with no behaviour, following this repo's own
Phase-1 precedent (`internal/table/category.go`'s "PHASE 1 DECLARATION ONLY"
block for RDR 0008's `Failure` fields):

| Declared | Where | Deliberately NOT done |
| --- | --- | --- |
| `table.Accessor` command fields | `internal/table/model.go` | nothing decodes them |
| the six `CatCommand*` constants | `internal/table/category.go` | NOT appended to `Categories()`; no load arm raises them |
| `accessor.ExecError` | `internal/accessor/model.go` | nothing constructs it in production |
| `accessor.Refusal.Detail` | `internal/accessor/model.go` | nothing populates it |
| `cmdbind.{Reader,Gate,Writer}` + constants | `internal/cli/cmdbind/cmdbind.go` | every method returns a not-implemented `*ExecError` |
| `Registry`'s C6 signature | `internal/cli/flowbind/registry.go` | both new parameters ignored |

## Test files

| File | Scope |
| --- | --- |
| `internal/table/command_carrier_0025_test.go` | C1 carrier, C5's six defect classes, registration, precedence |
| `internal/accessor/exec_detail_0025_test.go` | C4 `ExecError` / `Refusal.Detail` / both refusal constructors |
| `internal/cli/cmdbind/fixtures_0025_test.go` | shared fixtures — a real helper child; argv and env oracles |
| `internal/cli/cmdbind/invocation_0025_test.go` | C2 authority bound, C3 envelope, C4 safety; S2/S3/S4/S4b/S6b |
| `internal/cli/cmdbind/seam_0025_test.go` | C7 seam; S5 write read-back |
| `internal/cli/flowbind/registry_command_0025_test.go` | carrier selection, C1 residue, C6 gate site, S5b |
| `internal/cli/command_fixtures_0025_test.go` | CLI fixtures — declared wrapper, model writers, spawn sentinel |
| `internal/cli/command_gate_0025_test.go` | C6 containment (S7), lint ungated, `CLIError.Detail` order |
| `internal/cli/command_mvv_0025_test.go` | **REQ-MVV runner** (S6), end to end against `git config --file` |

## REQ -> test

| REQ | Test(s) | File(s) |
| --- | --- | --- |
| REQ-1 | TestReq1_CommandDecodesAsTheEntrysArgvVector | table/carrier |
| REQ-2 | TestReq2_OutputCarriesJSONOrRawAndDistinguishesOmission | table/carrier |
| REQ-3 | TestReq3And4_ExitMapsDecodeOnTheirOwnCapability | table/carrier |
| REQ-4 | TestReq3And4_ExitMapsDecodeOnTheirOwnCapability | table/carrier |
| REQ-5 | TestReq5And6_EnvAndEnvPassDecodeOnAnEntry | table/carrier |
| REQ-6 | TestReq5And6_EnvAndEnvPassDecodeOnAnEntry <br> TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited | cmdbind/invocation <br> table/carrier |
| REQ-7 | TestReq7_ExactlyOneOfPathOrCommandPerEntry | table/carrier |
| REQ-8 | TestReq8_ACarrierLessEntryBuildsARefusingBindingNotAnEmptyPathOne | flowbind/registry |
| REQ-9 | TestReq9_CommandIsNonEmptyAndCarriesNoEmptyElement | table/carrier |
| REQ-10 | TestReq10_TheOtherFourRulesAreUnchangedForACommandEntry | table/carrier |
| REQ-11 | TestReq8_ACarrierLessEntryBuildsARefusingBindingNotAnEmptyPathOne | flowbind/registry |
| REQ-12 | TestReq12_AllSixCarrierFieldsSurviveOntoTheNormalizedAccessor | table/carrier |
| REQ-13 | TestReq2_OutputCarriesJSONOrRawAndDistinguishesOmission | table/carrier |
| REQ-14 | TestReq7_ExactlyOneOfPathOrCommandPerEntry | table/carrier |
| REQ-15 | TestReq112_TheWireFieldIsCommandAndTheRejectedSpellingsAreNot <br> TestReq26_TheLoaderPerformsNoPathResolutionAndModelCarriesNoBaseDir | table/carrier |
| REQ-16 | TestReq17_TheExecutedArgvIsTheDeclaredVectorAfterOneSubstitution | cmdbind/invocation |
| REQ-17 | TestReq17_TheExecutedArgvIsTheDeclaredVectorAfterOneSubstitution | cmdbind/invocation |
| REQ-18 | TestReq18_PlaceholdersAreWholeElementAndDrawnFromTheClosedV1Vocabulary | table/carrier |
| REQ-19 | TestReq17_TheExecutedArgvIsTheDeclaredVectorAfterOneSubstitution | cmdbind/invocation |
| REQ-20 | TestReq20_ARelativeOrDashLedArtifactPathRefusesBeforeSpawn | cmdbind/invocation |
| REQ-21 | TestReq20_ARelativeOrDashLedArtifactPathRefusesBeforeSpawn | cmdbind/invocation |
| REQ-22 | TestReq20_ARelativeOrDashLedArtifactPathRefusesBeforeSpawn | cmdbind/invocation |
| REQ-23 | TestReq23_ASeparatorBearingArgv0ResolvesAgainstTheModelDirNotTheCwd | cmdbind/invocation |
| REQ-24 | TestReq9_CommandIsNonEmptyAndCarriesNoEmptyElement | table/carrier |
| REQ-25 | TestReq23_ASeparatorBearingArgv0ResolvesAgainstTheModelDirNotTheCwd | cmdbind/invocation |
| REQ-26 | TestReq26_TheLoaderPerformsNoPathResolutionAndModelCarriesNoBaseDir | table/carrier |
| REQ-27 | TestReq23_ASeparatorBearingArgv0ResolvesAgainstTheModelDirNotTheCwd | cmdbind/invocation |
| REQ-28 | TestReq28_ASeparatorBearingArgv0WithNoBaseDirRefusesBeforeSpawn | cmdbind/invocation |
| REQ-29 | TestReq29_AWriteSendsThePlannedTagsAsAFlatObjectOfStrings | cmdbind/invocation |
| REQ-30 | TestReq30_AReadSendsAnEmptyObjectOnStdinAndNeverTheKeySet | cmdbind/invocation |
| REQ-31 | TestReq31_AJSONEnvelopeOmissionIsUnreadableNeverEstablishedAbsent | cmdbind/invocation |
| REQ-32 | TestReq75_OutputShapeArmsCoverRawArityAndOffCapabilityExitMaps <br> TestReq32_RawModeStripsExactlyOneTrailingNewlineAndNothingElse | cmdbind/invocation <br> table/carrier |
| REQ-33 | TestReq33_TheGateEnvelopeCarriesVerdictAndReasonBackThroughTheSeam | cmdbind/invocation |
| REQ-34 | TestReq34_AListedExitAbsentEstablishesEveryDeclaredKeyAbsent | cmdbind/invocation |
| REQ-35 | TestReq35_AnEmptyStdoutWithNoExitMapMatchIsAlwaysExecutionFailure | cmdbind/invocation |
| REQ-36 | TestReq36_AListedExitWithEmptyStdoutIsThatVerdict | cmdbind/invocation |
| REQ-37 | TestReq30_AReadSendsAnEmptyObjectOnStdinAndNeverTheKeySet | cmdbind/invocation |
| REQ-38 | TestReq38_AWriteThatExitsZeroWithoutApplyingStillRefuses | accessor/detail |
| REQ-39 | TestReq39_ExitVerdictsAdmitExactlyTheThreeAccessorVerdicts | table/carrier |
| REQ-40 | TestReq29_AWriteSendsThePlannedTagsAsAFlatObjectOfStrings <br> TestReq127_WriteVerificationIsReadBackAndNeverExitStatus | cmdbind/invocation <br> cmdbind/seam |
| REQ-41 | TestReq36_AListedExitWithEmptyStdoutIsThatVerdict | cmdbind/invocation |
| REQ-42 | TestReq42_ASpawnFailureIsNeverClaimedByAnOverBroadExitMap | cmdbind/invocation |
| REQ-43 | TestReq48_TheStdoutEnvelopeIsParsedBeforeTheExitCodeIsClassified <br> TestReq50_AStdoutOverTheOneMiBCapIsAnExecutionFailureNotATruncation | cmdbind/invocation |
| REQ-44 | TestReq34_AListedExitAbsentEstablishesEveryDeclaredKeyAbsent <br> TestReq106_EstablishedAbsenceCrossesAsTheFlagAndSurvivesClassify | cmdbind/invocation <br> cmdbind/seam |
| REQ-45 | TestReq35_AnEmptyStdoutWithNoExitMapMatchIsAlwaysExecutionFailure | cmdbind/invocation |
| REQ-46 | TestReq46_TheReservedLiteralIsUnreadableOnTheReadPathAndAMismatchOnReadBack | cmdbind/seam |
| REQ-47 | TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited | cmdbind/invocation |
| REQ-48 | TestReq48_TheStdoutEnvelopeIsParsedBeforeTheExitCodeIsClassified | cmdbind/invocation |
| REQ-49 | TestReq49_TheDeadlineBoundsTheWholeProcessGroupAndLeavesNoSurvivor | cmdbind/invocation |
| REQ-50 | TestReq55_TheStderrTailReachesTheGateRefusalDetail <br> TestReq50_AStdoutOverTheOneMiBCapIsAnExecutionFailureNotATruncation | accessor/detail <br> cmdbind/invocation |
| REQ-51 | TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited <br> TestReq51_AnUnsetAllowlistedOrPassedVariableIsNotSynthesized | cmdbind/invocation |
| REQ-52 | TestReq52_OnACollisionTheLaterLayerWins | cmdbind/invocation |
| REQ-53 | TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited | cmdbind/invocation |
| REQ-54 | TestReq54_TheIntrastatePrefixIsReservedFromEnvAndEnvPass | table/carrier |
| REQ-55 | TestReq55_TheStderrTailReachesTheGateRefusalDetail <br> TestReq55_TheStderrTailReachesTheWriteRefusalDetail | accessor/detail |
| REQ-56 | TestReq56_TheStderrTailReachesTheReadRefusalDetailToo | accessor/detail |
| REQ-57 | TestReq57_DetailIsEmptyOnEveryRefusalCarryingNoInvocationError | accessor/detail |
| REQ-58 | TestReq58_ExecErrorWrapsRatherThanFlattensTheOffendingError <br> TestReq42_ASpawnFailureIsNeverClaimedByAnOverBroadExitMap | accessor/detail <br> cmdbind/invocation |
| REQ-59 | TestReq38_AWriteThatExitsZeroWithoutApplyingStillRefuses <br> TestReq127_WriteVerificationIsReadBackAndNeverExitStatus | accessor/detail <br> cmdbind/seam |
| REQ-60 | TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn | cmdbind/invocation |
| REQ-61 | TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn | cmdbind/invocation |
| REQ-62 | TestReq62_TheRefusalClassSetIsUnchangedByThisRDR | accessor/detail |
| REQ-63 | TestReq62_TheRefusalClassSetIsUnchangedByThisRDR | accessor/detail |
| REQ-64 | TestReq56_TheStderrTailReachesTheReadRefusalDetailToo | accessor/detail |
| REQ-65 | TestReq65_APostMutationRefusalCarriesTheAppliedSenseAndTheTailSeparably <br> TestReq65_TheAppliedSenseLeadsAndTheStderrTailFollowsInOneSlot | accessor/detail <br> cli/gate |
| REQ-66 | TestReq55_TheStderrTailReachesTheGateRefusalDetail | accessor/detail |
| REQ-67 | TestReq48_TheStdoutEnvelopeIsParsedBeforeTheExitCodeIsClassified | cmdbind/invocation |
| REQ-68 | TestReq38_AWriteThatExitsZeroWithoutApplyingStillRefuses <br> TestReq127_WriteVerificationIsReadBackAndNeverExitStatus | accessor/detail <br> cmdbind/seam |
| REQ-69 | TestReq69_APathBackedBindingBehavesIdenticallyExceptForDetail | accessor/detail |
| REQ-70 | TestReq49_TheDeadlineBoundsTheWholeProcessGroupAndLeavesNoSurvivor | cmdbind/invocation |
| REQ-71 | TestReq7_ExactlyOneOfPathOrCommandPerEntry | table/carrier |
| REQ-72 | TestReq9_CommandIsNonEmptyAndCarriesNoEmptyElement | table/carrier |
| REQ-73 | TestReq18_PlaceholdersAreWholeElementAndDrawnFromTheClosedV1Vocabulary | table/carrier |
| REQ-74 | TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect | table/carrier |
| REQ-75 | TestReq75_OutputShapeArmsCoverRawArityAndOffCapabilityExitMaps | table/carrier |
| REQ-76 | TestReq54_TheIntrastatePrefixIsReservedFromEnvAndEnvPass | table/carrier |
| REQ-77 | TestReq77_TheSixCommandCategoriesAreRegisteredAtTheTailInClauseOrder | table/carrier |
| REQ-78 | TestReq78_TheSixCategoriesCarryTheirContractWireStrings | table/carrier |
| REQ-79 | TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree | table/carrier |
| REQ-80 | TestReq80_WithinOneEntryTheEarlierClauseWins | table/carrier |
| REQ-81 | TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree | table/carrier |
| REQ-82 | TestReq82_ADefectiveReadEntryMasksADefectiveGateEntry | table/carrier |
| REQ-83 | TestReq83_TheInterpreterSetIsDenyListedNotClosed | table/carrier |
| REQ-84 | TestReq84_TheAccessorValidationCodeSetIsNotExtendedByThisRDR | table/carrier |
| REQ-85 | TestReq77_TheSixCommandCategoriesAreRegisteredAtTheTailInClauseOrder | table/carrier |
| REQ-86 | TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect | table/carrier |
| REQ-87 | TestReq87_AnUndeclaredEntryFieldIsADocumentLevelUnknownSchemaField | table/carrier |
| REQ-88 | TestReq90_LintSitsAtRootAndDoesNotResolveTheFlag | cli/gate |
| REQ-89 | TestReq89_EveryVerbInTheFlowGroupResolvesTheAllowCommandsFlag | cli/gate |
| REQ-90 | TestReq90_LintSitsAtRootAndDoesNotResolveTheFlag | cli/gate |
| REQ-91 | TestReq91_ALookupMissIsAWiringBugAndSurfacesAsOne | cli/gate |
| REQ-92 | TestReq92_WithTheGateOffEveryCapabilityRefusesBeforeSpawn <br> TestReq92_WithoutTheFlagACommandModelRefusesAndNoChildSpawns <br> TestReq92_LintValidatesACommandModelRegardlessOfTheGate | cli/gate <br> cmdbind/invocation |
| REQ-93 | TestReq93_TheGateIsAppliedAtTheSingleConstructionSite | flowbind/registry |
| REQ-94 | TestReq94_RegistryTakesBaseDirAndAllowCommandsAndReturnsNoError | flowbind/registry |
| REQ-95 | TestReq94_RegistryTakesBaseDirAndAllowCommandsAndReturnsNoError | flowbind/registry |
| REQ-96 | TestReq92_WithTheGateOffEveryCapabilityRefusesBeforeSpawn <br> TestReq93_TheGateIsAppliedAtTheSingleConstructionSite | cmdbind/invocation <br> flowbind/registry |
| REQ-97 | TestReq92_WithTheGateOffEveryCapabilityRefusesBeforeSpawn <br> TestReq92_WithoutTheFlagACommandModelRefusesAndNoChildSpawns | cli/gate <br> cmdbind/invocation |
| REQ-98 | TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam |
| REQ-99 | TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam |
| REQ-100 | TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam |
| REQ-101 | TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam |
| REQ-102 | TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam |
| REQ-103 | TestReq103_OneApplyIsOneSpawnAndOneIncrement | cmdbind/seam |
| REQ-104 | TestReq33_TheGateEnvelopeCarriesVerdictAndReasonBackThroughTheSeam <br> TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/invocation <br> cmdbind/seam |
| REQ-105 | TestReq31_AJSONEnvelopeOmissionIsUnreadableNeverEstablishedAbsent | cmdbind/invocation |
| REQ-106 | TestReq34_AListedExitAbsentEstablishesEveryDeclaredKeyAbsent <br> TestReq106_EstablishedAbsenceCrossesAsTheFlagAndSurvivesClassify | cmdbind/invocation <br> cmdbind/seam |
| REQ-107 | TestReq106_EstablishedAbsenceCrossesAsTheFlagAndSurvivesClassify <br> TestReq107_ACommandBindingSignallingAbsenceByOmissionIsRefused | cmdbind/seam |
| REQ-108 | TestReq106_EstablishedAbsenceCrossesAsTheFlagAndSurvivesClassify | cmdbind/seam |
| REQ-109 | TestReq109_TheCarrierDoesNotEnterAccessorIdentity | flowbind/registry |
| REQ-110 | TestReq141_ReadBackComparesTypedTagValuesAndNotArtifactBytes | cmdbind/seam |
| REQ-111 | TestReq30_AReadSendsAnEmptyObjectOnStdinAndNeverTheKeySet | cmdbind/invocation |
| REQ-112 | TestReq112_TheWireFieldIsCommandAndTheRejectedSpellingsAreNot | table/carrier |
| REQ-113 | TestReq113_TheRegistrySelectsTheBindingByCarrierField | flowbind/registry |
| REQ-114 | TestReq113_TheRegistrySelectsTheBindingByCarrierField | flowbind/registry |
| REQ-115 | TestReq12_AllSixCarrierFieldsSurviveOntoTheNormalizedAccessor | table/carrier |
| REQ-116 | TestReq26_TheLoaderPerformsNoPathResolutionAndModelCarriesNoBaseDir <br> TestReq98_TheCommandBindingsSatisfyTheUnchangedInterfaces | cmdbind/seam <br> table/carrier |
| REQ-117 | TestReq93_TheGateIsAppliedAtTheSingleConstructionSite | flowbind/registry |
| REQ-118 | TestReq92_LintValidatesACommandModelRegardlessOfTheGate | cli/gate |
| REQ-119 | TestReq119_TheContractIsTheCarrierNotTheIllustration | cmdbind/seam |
| REQ-120 | TestReq128_TheRegistryEmitsReadersNameSortedForFirstMatchSelection | flowbind/registry |
| REQ-121 | TestReqMVV_ALintedModelAppliesAndVerifiesARealChangeThroughACommand | cli/mvv |
| REQ-122 | TestReq83_TheInterpreterSetIsDenyListedNotClosed <br> TestReq84_TheAccessorValidationCodeSetIsNotExtendedByThisRDR | table/carrier |
| REQ-123 | TestReq17_TheExecutedArgvIsTheDeclaredVectorAfterOneSubstitution | cmdbind/invocation |
| REQ-124 | TestReq35_AnEmptyStdoutWithNoExitMapMatchIsAlwaysExecutionFailure <br> TestReq65_TheAppliedSenseLeadsAndTheStderrTailFollowsInOneSlot | cli/gate <br> cmdbind/invocation |
| REQ-125 | TestReq49_TheDeadlineBoundsTheWholeProcessGroupAndLeavesNoSurvivor | cmdbind/invocation |
| REQ-126 | TestReq51_TheChildEnvIsComposedFromTheAllowlistNotInherited | cmdbind/invocation |
| REQ-127 | TestReq127_WriteVerificationIsReadBackAndNeverExitStatus <br> TestReq46_TheReservedLiteralIsUnreadableOnTheReadPathAndAMismatchOnReadBack | cmdbind/seam |
| REQ-128 | TestReq128_TheRegistryEmitsReadersNameSortedForFirstMatchSelection <br> TestReq128_ASecondReaderOnAWritesOwnedKeyIsRefusedAtLoad | flowbind/registry |
| REQ-129 | TestReqMVV_ALintedModelAppliesAndVerifiesARealChangeThroughACommand | cli/mvv |
| REQ-130 | TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn | cmdbind/invocation |
| REQ-131 | TestReq8_ACarrierLessEntryBuildsARefusingBindingNotAnEmptyPathOne <br> TestReq89_EveryVerbInTheFlowGroupResolvesTheAllowCommandsFlag <br> TestReq92_WithoutTheFlagACommandModelRefusesAndNoChildSpawns | cli/gate <br> flowbind/registry |
| REQ-132 | TestReq132_ALoadTimeCommandDefectNamesTheEntryAndTheOffendingElement | table/carrier |
| REQ-133 | TestReq48_TheStdoutEnvelopeIsParsedBeforeTheExitCodeIsClassified | cmdbind/invocation |
| REQ-134 | TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn | cmdbind/invocation |
| REQ-135 | TestReq20_ARelativeOrDashLedArtifactPathRefusesBeforeSpawn | cmdbind/invocation |
| REQ-136 | TestReq136_AStdinSinkingWriteToolIsNotStaticallyDetectableAndSurfacesAtReadBack | cmdbind/seam |
| REQ-137 | TestReq103_OneApplyIsOneSpawnAndOneIncrement | cmdbind/seam |
| REQ-138 | TestReq138_NoV1LintAdvisoryRefusesACredentialShapedArgvElement | table/carrier |
| REQ-139 | TestReq139_APathBackedModelIsUnchangedByTheCommandCarrier <br> TestReq93_TheGateIsAppliedAtTheSingleConstructionSite | flowbind/registry <br> table/carrier |
| REQ-140 | TestReq32_RawModeStripsExactlyOneTrailingNewlineAndNothingElse | cmdbind/invocation |
| REQ-141 | TestReq141_ReadBackComparesTypedTagValuesAndNotArtifactBytes | cmdbind/seam |
| REQ-142 | TestReq18_PlaceholdersAreWholeElementAndDrawnFromTheClosedV1Vocabulary | table/carrier |
| REQ-MVV | TestReqMVV_ALintedModelAppliesAndVerifiesARealChangeThroughACommand | cli/mvv |

## How the negative REQs are covered

A clause that FORBIDS something has no direct oracle, so each is covered by a
test of the property it protects, paired with a positive control so the
assertion cannot pass vacuously.

| REQ | Property asserted instead | Positive control |
| --- | --- | --- |
| REQ-15, REQ-112 | the WIRE spelling `command` is the contract; the oracle reads TOML and never names a Go field or struct tag | the three rejected spellings (`exec`/`run`/`argv`) decode as unknown fields |
| REQ-26 | the loader takes bytes plus a display sourceID and resolves no path; the model carries no base dir | the carrier must be live, or the claim is vacuous |
| REQ-46 | this RDR adds no executor arm; the two `clearIsUnreadable` senses are asserted separately | neither subtest stands in for the other |
| REQ-63 | `Reason` stays empty on an execution failure | asserted on a refusal that DOES carry a tail |
| REQ-69, REQ-114, REQ-139 | path-backed behaviour is unchanged | a refusal-by-refusal comparison, and a command-side sibling that differs |
| REQ-79, REQ-81 | append-only, duplicate-free, additions only at the tail — never a count, never a cross-entry order | the six additions must be present first |
| REQ-98 | compile-time interface satisfaction, the only oracle Go offers for "no signature changed" | one real invocation per capability |
| REQ-111 | the requested key set never rides stdin | asserted on the child's OBSERVED stdin bytes |
| REQ-119 | a deliberately unlike entry behaves identically | the illustration's tool is never a dependency |
| REQ-125 | the assembled deadline triple, asserted on a real process group | the ablations are CITED from FX-deadline and NOT built as arms — a seam to disable them would ship a way to weaken the deadline |
| REQ-136 | no `command_*` category refuses a stdin-sinking entry | the corruption is observable and surfaces at read-back |
| REQ-138 | no v1 refusal on credential-shaped argv or env | the carrier must be live |
| REQ-141 | typed tag-value equality across a repeated write | byte identity neither asserted nor required |

## Phase-0 readings encoded

Each is marked in the source with a `PHASE-0 READING` comment.

| Q | Reading | Where |
| --- | --- | --- |
| Q1 | C4's allowlist passes `PATH`; argv0 `LookPath` uses the PARENT's `PATH` (A2/C4) | `cmdbind/invocation` — REQ-51 "the allowlist passes", REQ-28 "a BARE argv0" |
| Q2 | `output = "json"` with a well-formed empty `{}` stdout ⇒ every declared key UNREADABLE ⇒ `incomplete_read`, not `execution_failure` | `cmdbind/invocation` REQ-31 omission arm; `cmdbind/seam` REQ-107 |
| Q3 | `read_back = false` on a command write registers NO new C5 arm; the rule is RDR 0004's, inherited | `table/carrier` REQ-10 — asserted as the INHERITED `malformed_accessor_declaration`, not a `command_*` category |
| Q4 | S7's containment arm asserted AS WORDED — a structural `Flags().Lookup` walk | `cli/gate` REQ-89 |

Also recorded, from `req-list.md`'s 20 ASSUMPTION lines, where the reading is
load-bearing for a test's shape: REQ-2 (two output literals), REQ-3/REQ-4
(capability is the discriminator), REQ-20 (the dash guard is on the
substituted value only), REQ-32 (at most one trailing newline), REQ-34/REQ-106
(every declared key, whole-invocation), REQ-50 (overflow is detected, not
truncated), REQ-51 (an unset variable is not synthesized), REQ-54
(case-sensitive literal prefix), REQ-60/REQ-61 (an unexported `goos` reached
through `export_test.go`), REQ-77 (clause order), REQ-91 (an error return,
not a panic), REQ-103 (Apply ENTRIES, even on a failed spawn).

## Normative fixtures cited

| Fixture | Cited by |
| --- | --- |
| **FX-raw-read** (spike R1: stdout `"draft\n"` -> `draft`) | `cmdbind/invocation` REQ-32 |
| **FX-exit-codes** (E1c, E1d, E1e, E1f, E1g, E1h, R6) | `cmdbind/invocation` REQ-35, REQ-36, REQ-34, REQ-42 |
| **FX-deadline** (S1, S2, S3, S4, S5) | `cmdbind/invocation` REQ-49 |

The non-empty-stdout arms of S3 (cases 1, 2, 8) are asserted as
CONTRACT-DERIVED with no fixture oracle, exactly as S3 directs — their
expected values follow from C3/C4's ordering rule and FX-exit-codes says
nothing about them.
