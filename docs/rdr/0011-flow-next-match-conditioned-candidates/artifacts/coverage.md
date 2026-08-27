# Coverage — RDR 0011 Phase 1 tests

Phase 1 artifact: REQ-N x test name, both orphan directions, and the
red-gate classification. Generated from the test headers themselves, so a
REQ cited in a header but asserted by nothing would show as an orphan here.

## Totals

- REQs in `req-list.md`: **130** (REQ-1..REQ-130)
- REQs cited by at least one asserting test: **130/130**
- Orphan REQs (no test): **0**
- Tests written: **62** across 7 `*_0011_test.go` files
- Orphan tests (citing no REQ): **0**
- `red_confirmed`: **113 REQs** covered by at least one FAILING test
- `green_baseline`: **17 REQs** whose every test is green today

## The red-gate classification

The gate is *every REQ describing a behaviour change is red*. It holds: all
53 change-behaviour tests fail against the shipped build, and the 9 green
tests are **preservation or negative obligations** — clauses that require
something to STAY true rather than to BECOME true. Each is listed below with
why it is green rather than under-asserted, and each carries a vacuity guard
where one is possible, because C3 states outright that "a green 0005 suite is
NOT evidence C1 shipped".

No REQ that this RDR CHANGES is green. There is no `unexpected_green`.

### The green-baseline tests, and why each is green

| test | REQs | why it is green today, and what it would catch |
| --- | --- | --- |
| `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` | 46, 47, 65, 98 | `--all` exists nowhere yet, so its absence from `resolve`/`read-state`/`set-state` and the `flow` group holds. It is vacuity-guarded on the four verbs existing (they do), and it FAILS the moment Phase 2 registers the flag on the shared `registerSelectionFlags` or as a persistent flag. |
| `TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode` | 48, 49, 98 | The corroborating behavioural run: `command-error` + exit 2 is already what pflag's failed parse maps to. It fails if the build mints a bespoke `flow-all*` code, which REQ-49 forbids. |
| `TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged` | 35, 118 | Byte-identity is the obligation. The goldens are captured from THIS pre-change tree (A-11), so the test is green by construction now and fails if the demand-set term changes a payload it must not — the payload-diff form S8 says "the build owes" because the spike could not reach it. |
| `TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest` | 63, 90, 106 | The three load/parse refusals already ship; C3 pins them so a future relaxation "fails a named oracle rather than silently reopening the fold". Green is the correct state — the pins are what C1's exactness RESTS on. |
| `TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus` | 64, 66, 105 | Three NEGATIVE obligations naming fixtures the build must NOT contain. Green because this RDR's own corpus is clean; it fails if a later phase authors an empty `[rule.match]`, a conflicted-key `TagSet`, or one of the three unreachable `uncomparable` producers. |
| `TestReq94And119_MVV9TheKernelIsUntouchedInBothOfTheFormsTheMVVNames` | 94, 119 | `git diff --stat internal/resolve` empty, and no kernel test file changed. Green now; fails if Phase 2 reaches into the kernel. |
| `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` | 34, 78, 79, 107, 130 | Same claim from the wire side: `RefusalKinds()` still closes at five, no sixth minted. |
| `TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec` | 120 | A SCOPE claim about the diff. Green now; fails if the predicate spreads past `flow_next.go` + the one `flow_exec.go` term + the files C3 names. |
| `TestReq108_TheUpstreamKernelPinsStillExistAsShipped` | 108 | The three upstream pins A13's merge relies on still exist. Green is the required state. |

### One vacuous green found and fixed

`TestReq58` (ok-bool asserted) and `TestReq59` (`stringsAt` not repointed)
initially passed for the WRONG reason: they scan the re-homed reads, and
before the rename those reads do not exist, so the scan found nothing and
passed. That is exactly the failure C3 warns about. Both now carry a vacuity
guard that FATALs while the files still read `"unresolved"`, so both are RED
until the re-homing lands. They are counted under `red_confirmed`.


## REQ-N x test

`R` = at least one citing test fails today (red). `G` = every citing test is green (preservation / negative obligation).

| REQ | gate | tests |
| --- | --- | --- |
| REQ-1 | R | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates` |
| REQ-2 | R | `TestReq2And67_TheProbeBindsItsOwnOutcomeStripsEscapeAndKeepsTheResult` |
| REQ-3 | R | `TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent` |
| REQ-4 | R | `TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent` |
| REQ-5 | R | `TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent` |
| REQ-6 | R | `TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm` |
| REQ-7 | R | `TestReq7And103_AMixedStateRowIsExcludedByTheUnequalKeyAlone` |
| REQ-8 | R | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-9 | R | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-10 | R | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-11 | R | `TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded` |
| REQ-12 | R | `TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates` |
| REQ-13 | R | `TestReq83And96_GuardExcludedRowsStayExcludedInBothModes` |
| REQ-14 | R | `TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally` |
| REQ-15 | R | `TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent` |
| REQ-16 | R | `TestReq16And114_AMatchOnlyOwnedKeyInvokesItsReaderAndIsMatchDecided` |
| REQ-17 | R | `TestReq17And117_BothVerbsAgreeOnTheReaderSetOverOneModelAndOutcome` |
| REQ-18 | R | `TestReq18And21And115_TheDemandSetIsModeIndependentAndEmptyWhereNoRowMatchOwns`, `TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent` |
| REQ-19 | R | `TestReq19And93_ResolveSeparatesIntoThreeArmsAndRefusesFlowNoMatchInNone` |
| REQ-20 | R | `TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost` |
| REQ-21 | R | `TestReq18And21And115_TheDemandSetIsModeIndependentAndEmptyWhereNoRowMatchOwns` |
| REQ-22 | R | `TestReq22And71_AllFourSourcesLandOnUnknownAsKeyReasonPairs` |
| REQ-23 | R | `TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded` |
| REQ-24 | R | `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable` |
| REQ-25 | R | `TestReq25_TheGuardPayloadIsReadOnlyOnAGuardUnevaluableDisposition` |
| REQ-26 | R | `TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey` |
| REQ-27 | R | `TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey` |
| REQ-28 | R | `TestReq28_GateIdsAreEmittedEvenOnAnOwnedStateUnavailableCandidate` |
| REQ-29 | R | `TestReq29And113_UnknownIsSortedByKeyThenReasonAndDedupedOnThePair` |
| REQ-30 | R | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` |
| REQ-31 | R | `TestReq31And85And112_TheOwnedStateUnavailablePrecedenceLimitIsPinned` |
| REQ-32 | R | `TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes` |
| REQ-33 | R | `TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes` |
| REQ-34 | G | `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` |
| REQ-35 | G | `TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged` |
| REQ-36 | R | `TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand` |
| REQ-37 | R | `TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict` |
| REQ-38 | R | `TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict` |
| REQ-39 | R | `TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears` |
| REQ-40 | R | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` |
| REQ-41 | R | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` |
| REQ-42 | R | `TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears` |
| REQ-43 | R | `TestReq43And128_EveryOther0005ObligationHoldsIdenticallyInBothModes` |
| REQ-44 | R | `TestReq44And71_TheUnknownFieldIsPresentAsAnEmptyListInBothModes` |
| REQ-45 | R | `TestReq45And76_AllIsPartOfRequestIdentity` |
| REQ-46 | G | `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` |
| REQ-47 | R | `TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand`, `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` |
| REQ-48 | G | `TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode` |
| REQ-49 | G | `TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode` |
| REQ-50 | R | `TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand` |
| REQ-51 | R | `TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag` |
| REQ-52 | R | `TestReq52_TheHelpSeparatesCandidateFromWhatResolveWillSelect` |
| REQ-53 | R | `TestReq53And111_The0005NextFileHeaderNamesThisRDRAsTheSourceOfTheDefault` |
| REQ-54 | R | `TestReq54And110_TheFiveShippedReadsOfUnresolvedAreReHomedToUnknown` |
| REQ-55 | R | `TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag` |
| REQ-56 | R | `TestReq56And122_TheShippedDescriptionsOfTheOldDefaultAreCorrected` |
| REQ-57 | R | `TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey` |
| REQ-58 | R | `TestReq58_TheReHomedReadsAssertTheOkBoolRatherThanDiscardingIt` |
| REQ-59 | R | `TestReq59_StringsAtKeepsItsShapeAndTheReHomingAddsASiblingHelper` |
| REQ-60 | R | `TestReq60_TheCommentAndFailureMessageProseNamingUnresolvedIsSwept` |
| REQ-61 | R | `TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow` |
| REQ-62 | R | `TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears` |
| REQ-63 | G | `TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest` |
| REQ-64 | G | `TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus` |
| REQ-65 | R | `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup`, `TestReq65_TheSurfaceFlagMapNamesAllOnNext` |
| REQ-66 | R | `TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus`, `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable` |
| REQ-67 | R | `TestReq2And67_TheProbeBindsItsOwnOutcomeStripsEscapeAndKeepsTheResult` |
| REQ-68 | R | `TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally` |
| REQ-69 | R | `TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes` |
| REQ-70 | R | `TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates` |
| REQ-71 | R | `TestReq44And71_TheUnknownFieldIsPresentAsAnEmptyListInBothModes`, `TestReq22And71_AllFourSourcesLandOnUnknownAsKeyReasonPairs` |
| REQ-72 | R | `TestReq72And73_TextModeCarriesBothMembersThroughTheSharedFlattener` |
| REQ-73 | R | `TestReq72And73_TextModeCarriesBothMembersThroughTheSharedFlattener` |
| REQ-74 | R | `TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes` |
| REQ-75 | R | `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable` |
| REQ-76 | R | `TestReq45And76_AllIsPartOfRequestIdentity`, `TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent` |
| REQ-77 | R | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates`, `TestReq77And101_APresentAndEqualMatchKeyLeavesNoUnknownEntryForThatKey` |
| REQ-78 | R | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates`, `TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates`, `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` |
| REQ-79 | R | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates`, `TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded`, `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` |
| REQ-80 | R | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-81 | R | `TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey` |
| REQ-82 | R | `TestReq82And104_ARecognizedOnlyRowIsACandidateWithNoMatchFacts` |
| REQ-83 | R | `TestReq83And96_GuardExcludedRowsStayExcludedInBothModes` |
| REQ-84 | R | `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable` |
| REQ-85 | R | `TestReq31And85And112_TheOwnedStateUnavailablePrecedenceLimitIsPinned` |
| REQ-86 | R | `TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict` |
| REQ-87 | R | `TestReq87And95_MVV3NarrowsTheCheckedInModelFrom21RowsToExactlyThree` |
| REQ-88 | R | `TestReq88_MVV4UnderAllTheSameModelReportsThe21RowsTheBaselineDid` |
| REQ-89 | R | `TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow` |
| REQ-90 | R | `TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest`, `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates` |
| REQ-91 | R | `TestReq91_MVV7UnknownIsEmptyOnAFullyResolvedCandidateInBothModes` |
| REQ-92 | R | `TestReq92_MVV8TheMatchOnlyOwnedKeyIsReadAndDecidedInBothModes` |
| REQ-93 | R | `TestReq19And93_ResolveSeparatesIntoThreeArmsAndRefusesFlowNoMatchInNone`, `TestReq92_MVV8TheMatchOnlyOwnedKeyIsReadAndDecidedInBothModes` |
| REQ-94 | G | `TestReq94And119_MVV9TheKernelIsUntouchedInBothOfTheFormsTheMVVNames` |
| REQ-95 | R | `TestReq87And95_MVV3NarrowsTheCheckedInModelFrom21RowsToExactlyThree` |
| REQ-96 | R | `TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow`, `TestReq83And96_GuardExcludedRowsStayExcludedInBothModes` |
| REQ-97 | R | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` |
| REQ-98 | G | `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup`, `TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode` |
| REQ-99 | R | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates`, `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-100 | R | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` |
| REQ-101 | R | `TestReq77And101_APresentAndEqualMatchKeyLeavesNoUnknownEntryForThatKey` |
| REQ-102 | R | `TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm` |
| REQ-103 | R | `TestReq7And103_AMixedStateRowIsExcludedByTheUnequalKeyAlone` |
| REQ-104 | R | `TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally`, `TestReq82And104_ARecognizedOnlyRowIsACandidateWithNoMatchFacts` |
| REQ-105 | G | `TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus` |
| REQ-106 | G | `TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest` |
| REQ-107 | G | `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` |
| REQ-108 | G | `TestReq108_TheUpstreamKernelPinsStillExistAsShipped` |
| REQ-109 | R | `TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent`, `TestReq109_TheCLIViewsKeySetAgreesWithTheKernelsModuloRecognized` |
| REQ-110 | R | `TestReq54And110_TheFiveShippedReadsOfUnresolvedAreReHomedToUnknown` |
| REQ-111 | R | `TestReq53And111_The0005NextFileHeaderNamesThisRDRAsTheSourceOfTheDefault` |
| REQ-112 | R | `TestReq91_MVV7UnknownIsEmptyOnAFullyResolvedCandidateInBothModes`, `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable`, `TestReq31And85And112_TheOwnedStateUnavailablePrecedenceLimitIsPinned` |
| REQ-113 | R | `TestReq29And113_UnknownIsSortedByKeyThenReasonAndDedupedOnThePair` |
| REQ-114 | R | `TestReq16And114_AMatchOnlyOwnedKeyInvokesItsReaderAndIsMatchDecided` |
| REQ-115 | R | `TestReq18And21And115_TheDemandSetIsModeIndependentAndEmptyWhereNoRowMatchOwns` |
| REQ-116 | R | `TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost` |
| REQ-117 | R | `TestReq17And117_BothVerbsAgreeOnTheReaderSetOverOneModelAndOutcome` |
| REQ-118 | G | `TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged` |
| REQ-119 | G | `TestReq94And119_MVV9TheKernelIsUntouchedInBothOfTheFormsTheMVVNames` |
| REQ-120 | G | `TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec` |
| REQ-121 | R | `TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag` |
| REQ-122 | R | `TestReq56And122_TheShippedDescriptionsOfTheOldDefaultAreCorrected` |
| REQ-123 | R | `TestReq123And126_NextIsEffectFreeInBothModes` |
| REQ-124 | R | `TestReq124And129_TheReasonVocabularyIsClosedAndEachTokenIsScoped`, `TestReq124_UncomparableNeverRidesAMatchOwnedKeyOrGateEntry` |
| REQ-125 | R | `TestReq125_NoVersionMarkerIsMintedAndTheOldFieldIsGone` |
| REQ-126 | R | `TestReq123And126_NextIsEffectFreeInBothModes` |
| REQ-127 | R | `TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm` |
| REQ-128 | R | `TestReq43And128_EveryOther0005ObligationHoldsIdenticallyInBothModes` |
| REQ-129 | R | `TestReq124And129_TheReasonVocabularyIsClosedAndEachTokenIsScoped` |
| REQ-130 | G | `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` |

## Test x REQ (the other orphan direction)

| file | test | label | REQs | gate |
| --- | --- | --- | --- | --- |
| `flow_all_0011_test.go` | `TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand` | HAPPY PATH | REQ-36, REQ-47, REQ-50 | R |
| `flow_all_0011_test.go` | `TestReq37And38And86_UnderAllTheMatchPatternTakesNoPartInTheVerdict` | HAPPY PATH | REQ-37, REQ-38, REQ-86 | R |
| `flow_all_0011_test.go` | `TestReq39And42_UnderAllAMatchAtomsUnknownEntryDisappears` | ADVERSARIAL | REQ-39, REQ-42, REQ-62 | R |
| `flow_all_0011_test.go` | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` | ADVERSARIAL | REQ-30, REQ-40, REQ-41, REQ-97 | R |
| `flow_all_0011_test.go` | `TestReq43And128_EveryOther0005ObligationHoldsIdenticallyInBothModes` | DOMAIN EDGE | REQ-43, REQ-128 | R |
| `flow_all_0011_test.go` | `TestReq44And71_TheUnknownFieldIsPresentAsAnEmptyListInBothModes` | BOUNDARY | REQ-44, REQ-71 | R |
| `flow_all_0011_test.go` | `TestReq45And76_AllIsPartOfRequestIdentity` | BOUNDARY | REQ-45, REQ-76 | R |
| `flow_all_0011_test.go` | `TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` | ADVERSARIAL | REQ-46, REQ-47, REQ-65, REQ-98 | G |
| `flow_all_0011_test.go` | `TestReq48And49_TheUnknownFlagRefusalIsTheSharedUsageClassAndMintsNoNewCode` | ADVERSARIAL | REQ-48, REQ-49, REQ-98 | G |
| `flow_all_0011_test.go` | `TestReq65_TheSurfaceFlagMapNamesAllOnNext` | HAPPY PATH | REQ-65 | R |
| `flow_all_0011_test.go` | `TestReq123And126_NextIsEffectFreeInBothModes` | ADVERSARIAL | REQ-123, REQ-126 | R |
| `flow_demand_0011_test.go` | `TestReq16And114_AMatchOnlyOwnedKeyInvokesItsReaderAndIsMatchDecided` | ADVERSARIAL | REQ-16, REQ-114 | R |
| `flow_demand_0011_test.go` | `TestReq18And21And115_TheDemandSetIsModeIndependentAndEmptyWhereNoRowMatchOwns` | BOUNDARY | REQ-18, REQ-21, REQ-115 | R |
| `flow_demand_0011_test.go` | `TestReq17And117_BothVerbsAgreeOnTheReaderSetOverOneModelAndOutcome` | ADVERSARIAL | REQ-17, REQ-117 | R |
| `flow_demand_0011_test.go` | `TestReq19And93_ResolveSeparatesIntoThreeArmsAndRefusesFlowNoMatchInNone` | DOMAIN EDGE | REQ-19, REQ-93 | R |
| `flow_demand_0011_test.go` | `TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost` | ADVERSARIAL | REQ-20, REQ-116 | R |
| `flow_demand_0011_test.go` | `TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged` | ADVERSARIAL | REQ-35, REQ-118 | G |
| `flow_mvv_0011_test.go` | `TestReq87And95_MVV3NarrowsTheCheckedInModelFrom21RowsToExactlyThree` | HAPPY PATH | REQ-87, REQ-95 | R |
| `flow_mvv_0011_test.go` | `TestReq88_MVV4UnderAllTheSameModelReportsThe21RowsTheBaselineDid` | BOUNDARY | REQ-88 | R |
| `flow_mvv_0011_test.go` | `TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow` | DOMAIN EDGE | REQ-61, REQ-89, REQ-96 | R |
| `flow_mvv_0011_test.go` | `TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest` | ADVERSARIAL | REQ-63, REQ-90, REQ-106 | G |
| `flow_mvv_0011_test.go` | `TestReq64And66And105_TheForbiddenFixtureShapesAreAbsentFromTheCorpus` | ADVERSARIAL | REQ-64, REQ-66, REQ-105 | G |
| `flow_mvv_0011_test.go` | `TestReq91_MVV7UnknownIsEmptyOnAFullyResolvedCandidateInBothModes` | BOUNDARY | REQ-91, REQ-112 | R |
| `flow_mvv_0011_test.go` | `TestReq92_MVV8TheMatchOnlyOwnedKeyIsReadAndDecidedInBothModes` | HAPPY PATH | REQ-92, REQ-93 | R |
| `flow_mvv_0011_test.go` | `TestReq94And119_MVV9TheKernelIsUntouchedInBothOfTheFormsTheMVVNames` | ADVERSARIAL | REQ-94, REQ-119 | G |
| `flow_next_0011_test.go` | `TestReq1And77And78And79_TheThreeReachableMatchClassesSortToCandidates` | HAPPY PATH | REQ-1, REQ-77, REQ-78, REQ-79, REQ-90, REQ-99 | R |
| `flow_next_0011_test.go` | `TestReq11And23And79_AnAbsentMatchKeyIsReportedAbsentNotExcluded` | INPUT EDGE | REQ-11, REQ-23, REQ-79 | R |
| `flow_next_0011_test.go` | `TestReq77And101_APresentAndEqualMatchKeyLeavesNoUnknownEntryForThatKey` | HAPPY PATH | REQ-77, REQ-101 | R |
| `flow_next_0011_test.go` | `TestReq12And70And78_AMatchExcludedRowsGatesDoNotRunUnderEvaluateGates` | ADVERSARIAL | REQ-12, REQ-70, REQ-78 | R |
| `flow_next_0011_test.go` | `TestReq83And96_GuardExcludedRowsStayExcludedInBothModes` | DOMAIN EDGE | REQ-13, REQ-83, REQ-96 | R |
| `flow_next_0011_test.go` | `TestReq14And68And104_ARowWhoseMatchAtomsAreAllOmittedMatchesUnconditionally` | BOUNDARY | REQ-14, REQ-68, REQ-104 | R |
| `flow_next_0011_test.go` | `TestReq82And104_ARecognizedOnlyRowIsACandidateWithNoMatchFacts` | DOMAIN EDGE | REQ-82, REQ-104 | R |
| `flow_next_0011_test.go` | `TestReq10And80And99And100_TheDeadRowIsTheKernelsConjunctionNotACLIComparison` | ADVERSARIAL | REQ-8, REQ-9, REQ-10, REQ-80, REQ-99, REQ-100 | R |
| `flow_next_0011_test.go` | `TestReq7And103_AMixedStateRowIsExcludedByTheUnequalKeyAlone` | ADVERSARIAL | REQ-7, REQ-103 | R |
| `flow_next_0011_test.go` | `TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm` | ADVERSARIAL | REQ-6, REQ-102, REQ-127 | R |
| `flow_next_0011_test.go` | `TestReq3And4And5_PresenceIsMapMembershipAndTheEmptyStringIsPresent` | BOUNDARY | REQ-3, REQ-4, REQ-5, REQ-109 | R |
| `flow_next_0011_test.go` | `TestReq15And18And76_TheReaderSetAndViewAreFixedOncePerInvocationAndModeIndependent` | BOUNDARY | REQ-15, REQ-18, REQ-76 | R |
| `flow_next_0011_test.go` | `TestReq32And33And69And74_TheAlphabetIsTheFullDeclaredSetInBothModes` | DOMAIN EDGE | REQ-32, REQ-33, REQ-69, REQ-74 | R |
| `flow_next_0011_test.go` | `TestReq34And107And130_TheKernelPackageIsUntouchedByThisContract` | ADVERSARIAL | REQ-34, REQ-78, REQ-79, REQ-107, REQ-130 | G |
| `flow_next_0011_test.go` | `TestReq2And67_TheProbeBindsItsOwnOutcomeStripsEscapeAndKeepsTheResult` | ADVERSARIAL | REQ-2, REQ-67 | R |
| `flow_next_0011_test.go` | `TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec` | ADVERSARIAL | REQ-120 | G |
| `flow_rehome_0011_test.go` | `TestReq54And110_TheFiveShippedReadsOfUnresolvedAreReHomedToUnknown` | ADVERSARIAL | REQ-54, REQ-110 | R |
| `flow_rehome_0011_test.go` | `TestReq58_TheReHomedReadsAssertTheOkBoolRatherThanDiscardingIt` | ADVERSARIAL | REQ-58 | R |
| `flow_rehome_0011_test.go` | `TestReq59_StringsAtKeepsItsShapeAndTheReHomingAddsASiblingHelper` | ADVERSARIAL | REQ-59 | R |
| `flow_rehome_0011_test.go` | `TestReq60_TheCommentAndFailureMessageProseNamingUnresolvedIsSwept` | ADVERSARIAL | REQ-60 | R |
| `flow_rehome_0011_test.go` | `TestReq53And111_The0005NextFileHeaderNamesThisRDRAsTheSourceOfTheDefault` | DOMAIN EDGE | REQ-53, REQ-111 | R |
| `flow_rehome_0011_test.go` | `TestReq51And55And121_TheHelpStatesTheMatchConditionedDefaultAndTheAllFlag` | DOMAIN EDGE | REQ-51, REQ-55, REQ-121 | R |
| `flow_rehome_0011_test.go` | `TestReq52_TheHelpSeparatesCandidateFromWhatResolveWillSelect` | ADVERSARIAL | REQ-52 | R |
| `flow_rehome_0011_test.go` | `TestReq56And122_TheShippedDescriptionsOfTheOldDefaultAreCorrected` | DOMAIN EDGE | REQ-56, REQ-122 | R |
| `flow_rehome_0011_test.go` | `TestReq109_TheCLIViewsKeySetAgreesWithTheKernelsModuloRecognized` | ADVERSARIAL | REQ-109 | R |
| `flow_rehome_0011_test.go` | `TestReq108_TheUpstreamKernelPinsStillExistAsShipped` | ADVERSARIAL | REQ-108 | G |
| `flow_unknown_0011_test.go` | `TestReq22And71_AllFourSourcesLandOnUnknownAsKeyReasonPairs` | HAPPY PATH | REQ-22, REQ-71 | R |
| `flow_unknown_0011_test.go` | `TestReq26And27And57And81_AnUnrunGateIdCarriesNotEvaluatedAndTheIdIsTheKey` | DOMAIN EDGE | REQ-26, REQ-27, REQ-57, REQ-81 | R |
| `flow_unknown_0011_test.go` | `TestReq28_GateIdsAreEmittedEvenOnAnOwnedStateUnavailableCandidate` | ADVERSARIAL | REQ-28 | R |
| `flow_unknown_0011_test.go` | `TestReq24And66And75And84And112_AnUnevaluableGuardOverAPresentKeyIsUncomparable` | DOMAIN EDGE | REQ-24, REQ-66, REQ-75, REQ-84, REQ-112 | R |
| `flow_unknown_0011_test.go` | `TestReq25_TheGuardPayloadIsReadOnlyOnAGuardUnevaluableDisposition` | ADVERSARIAL | REQ-25 | R |
| `flow_unknown_0011_test.go` | `TestReq31And85And112_TheOwnedStateUnavailablePrecedenceLimitIsPinned` | DOMAIN EDGE | REQ-31, REQ-85, REQ-112 | R |
| `flow_unknown_0011_test.go` | `TestReq29And113_UnknownIsSortedByKeyThenReasonAndDedupedOnThePair` | BOUNDARY | REQ-29, REQ-113 | R |
| `flow_unknown_0011_test.go` | `TestReq124And129_TheReasonVocabularyIsClosedAndEachTokenIsScoped` | ADVERSARIAL | REQ-124, REQ-129 | R |
| `flow_unknown_0011_test.go` | `TestReq124_UncomparableNeverRidesAMatchOwnedKeyOrGateEntry` | ADVERSARIAL | REQ-124 | R |
| `flow_unknown_0011_test.go` | `TestReq72And73_TextModeCarriesBothMembersThroughTheSharedFlattener` | DOMAIN EDGE | REQ-72, REQ-73 | R |
| `flow_unknown_0011_test.go` | `TestReq125_NoVersionMarkerIsMintedAndTheOldFieldIsGone` | ADVERSARIAL | REQ-125 | R |

No test cites zero REQs, and every test carries one of the five labels.

## Files

| file | what it pins |
| --- | --- |
| `internal/cli/flow_fixtures_0011_test.go` | the discriminating fixture corpus, the sibling `unknown` reader C3 mandates, and the diff/golden helpers |
| `internal/cli/flow_next_0011_test.go` | C1 — the match-conditioned predicate, the `disposition` grid, the probe's shape |
| `internal/cli/flow_all_0011_test.go` | C2 — the `--all` flag, its filter scope, and the filter/dedup ORDER |
| `internal/cli/flow_unknown_0011_test.go` | the `unknown` payload: four sources, closed vocabulary, dedup, sort, fidelity limit |
| `internal/cli/flow_demand_0011_test.go` | the demand-set term and the `flow resolve` change class it carries, breaking arm included |
| `internal/cli/flow_mvv_0011_test.go` | the MVV, steps 2-9, runnable end to end |
| `internal/cli/flow_rehome_0011_test.go` | C3 — help text, shipped prose, the five-read census, the key-set agreement pin |

## The discriminating oracles Phase 0 flagged

| flagged | test | status |
| --- | --- | --- |
| **REQ-102** set-kinded match key | `TestReq6And102And127_ASetKindedMatchKeyIsComparedInTheKernelsCanonicalForm` | RED |
| **REQ-97** filter/dedup order | `TestReq30And40And41And97_TheAllFilterIsScopedToTheAtomWalkAndDedupRunsLast` | RED |
| **REQ-116** S8 breaking arm | `TestReq20And116_TheBreakingArmIsPinnedAsAKnownCost` | RED |
| **REQ-77..86** `disposition` grid | `TestReq1And77And78And79`, `TestReq12And70And78`, `TestReq10And80And99And100`, `TestReq26And27And57And81`, `TestReq82And104`, `TestReq83And96`, `TestReq24And66And75And84And112`, `TestReq31And85And112`, `TestReq37And38And86` | all RED |
| **REQ-54** five-read census | `TestReq54And110` — reconciled by SYMBOL; the census COUNT is the contract | RED |
| **REQ-60** prose sweep | `TestReq60` — the WORD is the obligation, not the line numbers | RED |

### REQ-54 / REQ-60 line-number reconciliation

Checked against the tree at Phase 1: the five reads are at
`flow_next_0005_test.go:103,163,404`, `flow_adversarial_0005_test.go:446`,
and `flow_mvv_0005_test.go:66` — **exactly** the lines C3 cites, plus `:166`
as the `%#v` argument it names. The tree has NOT drifted. Both oracles are
nonetheless written against the symbol and the count rather than the line, so
they survive any later edit above them.


---

## Phase 2 — REQ-MVV run end to end (REQ-87..REQ-94)

Run against the implemented build on branch `worktree-rdr-0011`. Every step
below is the ACTUAL output, captured live, not a restatement of the
expectation. `make check` passes (fmt-check, vet, golangci-lint `0 issues`,
build, graph-lint, `go test -race ./...` all green).

### MVV 3 — the default narrows `models/rdr.toml` (REQ-87, REQ-95)

Artifact seeded through the CLI at `stage=resolved`, `status=draft`,
`gate_passed=false`.

```
$ intrastate flow next --model models/rdr.toml --artifact rdr=<art> --as=json
candidates: ['prelock', 'resolve-abandon', 'resolve-route-back']
count:      3
outcomes:   ['advance', 'revise', 'abandon']
unknown[prelock]: []
```

Exactly the three rules MVV 3 names, one per outcome; no rule whose
`match.stage` names another value survives, and `outcomes[]` is still the
model's full declared alphabet. The pre-change build reported 21 here.

### MVV 4 — `--all` restores the 21 (REQ-88)

```
$ intrastate flow next --model models/rdr.toml --artifact rdr=<art> --all --as=json
count:    21
outcomes: ['advance', 'revise', 'abandon']
candidates: ['final-abandon', 'final-route-back', 'finalize-blocked',
 'implement', 'prelock', 'prelock-abandon', 'prelock-route-back', 'propose',
 'propose-abandon', 'propose-again', 'reconcile', 'reconcile-abandon',
 'reconcile-route-back', 'refine', 'refine-abandon', 'refine-again',
 'resolve-abandon', 'resolve-assumptions', 'resolve-route-back',
 'seed-abandon', 'seed-revise']
```

21 rows, the same `outcomes[]`. `finalize-pass` — the 22nd — is absent in
BOTH modes, guard-excluded on `gate_passed`, exactly as MVV 4 predicts.

### MVV 5 — guard exclusions and the added discriminating row (REQ-89, REQ-61, REQ-96)

`TestReq61And89And96_MVV5GuardExclusionsAndTheAddedDiscriminatingRow`

```
--- PASS: .../guard-excluded-over-the-fixture-that-authors-it
--- PASS: .../added-row-absent-by-default-present-under-all
```

### MVV 6 — the three pins that license the presence test (REQ-90, REQ-63, REQ-106)

`TestReq63And90And106_MVV6TheThreePinsThatLicenseThePresenceTest`

```
--- PASS: .../two-readers-on-one-owned-key-is-refused-at-load
--- PASS: .../repeated-tag-key-is-refused-flow-tag-duplicate
--- PASS: .../tag-on-an-owned-key-is-refused-flow-tag-owned
```

The three match classes and the dead row are covered by
`TestReq1And77And78And79` / `TestReq10And80And99And100`, both PASS.

### MVV 7 — `unknown` present as `[]` rather than omitted (REQ-91)

`TestReq91_MVV7UnknownIsEmptyOnAFullyResolvedCandidateInBothModes` — PASS.
Confirmed on the wire: a candidate with nothing undecided emits
`"unknown":[]`, in both modes, never an omitted field.

### MVV 8 — the match-only owned key and the three `flow resolve` arms (REQ-92, REQ-93)

`TestReq92_MVV8TheMatchOnlyOwnedKeyIsReadAndDecidedInBothModes` — PASS: the
reader appears in `readers`, the key is in the view, and the row is
match-decided rather than reported `{key, absent}`; the same `readers` set
under `--all`.

`TestReq19And93_ResolveSeparatesIntoThreeArmsAndRefusesFlowNoMatchInNone`

```
--- PASS: .../bound-and-answering-yields-a-plan
--- PASS: .../role-unbound-yields-artifact-missing-not-no-match     (exit 2)
--- PASS: .../reader-refusing-yields-exit-3-not-no-match            (exit 3)
```

`flow-no-match` in none of the three. `flow resolve` over `models/rdr.toml`
and the shipped fixtures is byte-identical to the pre-change goldens
(`TestReq35And118` — PASS, full-payload form).

### MVV 9 — the kernel is untouched, and the 0005 suite (REQ-94)

```
$ make check
golangci-lint run ./...   →  0 issues.
go test -race ./...       →  every package ok
   internal/cli       89.2% of statements
   internal/resolve   95.5% of statements

$ git diff --stat main -- internal/resolve
(empty)

$ go test ./internal/resolve
ok   github.com/cwensel/intrastate/internal/resolve

$ git diff --stat main -- internal/cli/flow_{next,adversarial,mvv}_0005_test.go
 internal/cli/flow_adversarial_0005_test.go | 12 +++--
 internal/cli/flow_mvv_0005_test.go         | 17 +++++---
 internal/cli/flow_next_0005_test.go        | 70 +++++++++++++++++++--------
 3 files changed, 63 insertions(+), 36 deletions(-)

$ grep -c '"--all"' internal/cli/flow_next_0005_test.go
0
```

The 0005 `next` suite passes with only C3's five mechanical re-homings
(ok-bool asserted at every one), the prose sweep, and the header naming
this RDR. **A4's MOVES-UNDER-`--all` = 0 holds** — no 0005 `next` oracle
needed the flag, which `TestReq53And111` asserts as a negative.

### Text-mode rendering (REQ-72, REQ-73)

The shared generic flattener, not a per-verb template:

```
$ intrastate flow next --model <m> --artifact state=<a> --as=text
candidates[0].unknown[0].key: wanted
candidates[0].unknown[0].reason: absent
candidates[1].unknown: (none)
```

Both members of the pair reach the text caller, path-qualified. No bespoke
`key (reason)` form was minted.

## Phase 2 result

- **130/130 REQ green.** All 62 Phase-1 oracles pass; `make check` passes.
- Two Phase-1 oracles were repointed off a shape their own fixture could not
  reach (DEV-6) and off a vacuous decode (DEV-7). Neither weakens an
  assertion; both are recorded in `deviations.md`.
