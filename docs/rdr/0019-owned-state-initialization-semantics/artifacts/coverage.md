# Coverage — RDR 0019 Owned-state initialization semantics

Phase 1 artifact. Column 1 is the REQ id from `req-list.md`; column 2
names the test that fails if that clause is broken. An EMPTY column 2 is
the orphan mark — the REQ carries no test.

All tests live in `internal/cli/` and are RED against the current tree:
`flow init-state` is not registered, so no behavioural assertion can be
satisfied. Every refusal oracle rejects cobra's unregistered-command
error by name (`requireVerbRegistered`), so no test passes on the verb's
absence — the tautology the red gate must exclude.

Each test opens with the verbatim REQ quote it pins and a label of
`HAPPY PATH` / `INPUT EDGE` / `BOUNDARY` / `ADVERSARIAL` / `DOMAIN EDGE`.

Test files:

- `internal/cli/flow_initstate_carrier_0019_test.go`
- `internal/cli/flow_initstate_encoder_0019_test.go`
- `internal/cli/flow_initstate_fixtures_0019_test.go`
- `internal/cli/flow_initstate_mvv_0019_test.go`
- `internal/cli/flow_initstate_noop_0019_test.go`
- `internal/cli/flow_initstate_predicate_0019_test.go`
- `internal/cli/flow_initstate_surface_0019_test.go`

| REQ | Test |
| --- | --- |
| REQ-1 | `TestReq1_0019_InitialIsBothTheLintRootAndTheRuntimeBootstrapSource` |
| REQ-2 | `TestReq2And3And94_0019_NoReadPathSynthesizesTheInitialValueForAnAbsentKey` |
| REQ-3 | `TestReq2And3And94_0019_NoReadPathSynthesizesTheInitialValueForAnAbsentKey` |
| REQ-4 | `TestReq4And6And86_0019_ADecisionTableModelRefusesRatherThanSucceedingVacuously` |
| REQ-5 | `TestReq5_0019_TheClassArmKeysOnIsDecisionTableNotOnLenOwned` |
| REQ-6 | `TestReq4And6And86_0019_ADecisionTableModelRefusesRatherThanSucceedingVacuously` |
| REQ-7 | `TestReq7_0019_TheClassRefusalIsExitTwoNotANoOpSuccess` |
| REQ-8 | `TestReq8And13And14_0019_FlowGroupGainsExactlyInitState` |
| REQ-9 | `TestReq9_0019_NoShippedVerbIsWithdrawn`<br>`TestReq9_0019_RefusalPrintsNoUsageBlock`<br>`TestReq9_0019_ValidateModeRunsFirstOnInitState` |
| REQ-10 | `TestReq10_0019_InitStateCarriesNoWriteGrammar`<br>`TestReq10_0019_InitStateTakesTheSharedSelectionFlags` |
| REQ-11 | `TestReq11_0019_AllowCommandsIsInheritedNotDeclared`<br>`TestReq11_0019_AllowCommandsParsesOnInitState` |
| REQ-12 | `TestReq12_0019_EverySeedIsRoutedThroughTheWriteAccessorWithCommitTimeReadBack` |
| REQ-13 | `TestReq8And13And14_0019_FlowGroupGainsExactlyInitState`<br>`TestReq99And13_0019_SetStateGainsNoFromInitialFlag` |
| REQ-14 | `TestReq8And13And14_0019_FlowGroupGainsExactlyInitState` |
| REQ-15 | `TestReq15_0019_DocSitesEnumerateInitState`<br>`TestReq15_0019_NoSiteStillClaimsFourVerbs` |
| REQ-16 | `TestReq16_0019_LLMsTxtEnumeratesInitState` |
| REQ-17 | `TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts` |
| REQ-18 | `TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts`<br>`TestReq18And19_0019_SetSeedsTakeTheCanonicalSortedCompactUnescapedForm` |
| REQ-19 | `TestReq18And19_0019_SetSeedsTakeTheCanonicalSortedCompactUnescapedForm`<br>`TestReq19And20_0019_ScalarSeedsArePersistedVerbatimNotWrappedInAnArray` |
| REQ-20 | `TestReq19And20_0019_ScalarSeedsArePersistedVerbatimNotWrappedInAnArray` |
| REQ-21 | `TestReq21And22And23And83_0019_LoaderOnlyKindsSeedThroughTheNormalizedPath` |
| REQ-22 | `TestReq21And22And23And83_0019_LoaderOnlyKindsSeedThroughTheNormalizedPath` |
| REQ-23 | `TestReq21And22And23And83_0019_LoaderOnlyKindsSeedThroughTheNormalizedPath` |
| REQ-24 | `TestReq24And30_0019_APartiallySeededStoreIsNeverPerKeyMerged`<br>`TestReq24_0019_AnEmptyStoreSeedsEveryInitialKey`<br>`TestReq24_0019_BothArtifactsEmptySeedsBothRoles` |
| REQ-25 | `TestReq25And26And92_0019_OneSurvivingKeyAnywhereBlocksTheWholeSeed` |
| REQ-26 | `TestReq25And26And92_0019_OneSurvivingKeyAnywhereBlocksTheWholeSeed` |
| REQ-27 | `TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess` |
| REQ-28 | `TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess` |
| REQ-29 | `TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess`<br>`TestReq29_0019_TheEmptinessReadDoesNotInheritTheUnreadableShortCircuit` |
| REQ-30 | `TestReq24And30_0019_APartiallySeededStoreIsNeverPerKeyMerged` |
| REQ-31 | `TestReq31And32And62And85_0019_AnUnboundNeededRoleRefusesBeforeTheEmptinessRead` |
| REQ-32 | `TestReq31And32And62And85_0019_AnUnboundNeededRoleRefusesBeforeTheEmptinessRead` |
| REQ-33 | `TestReq33_0019_TheCarrierGatePreemptsTheRoleBindingRefusal` |
| REQ-34 | `TestReq34And61And84_0019_TwoWritersForOneInitialKeyRefusesAtLoad` |
| REQ-35 | `TestReq35_0019_TheClassRefusalPrecedesTheCarrierRefusal` |
| REQ-36 | `TestReq36_0019_TheCarrierRefusalPreemptsTheAllowCommandsRefusal` |
| REQ-37 | `TestReq37_0019_TheGatePipelineRunsClassCarrierBindingThenPredicate` |
| REQ-38 | `TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding` |
| REQ-39 | `TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding` |
| REQ-40 | `TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding` |
| REQ-41 | `TestReq41_0019_AFileBackedModelIsAdmittedOnBothSides` |
| REQ-42 | `TestReq42_0019_TheCarrierComesFromTheConstructedBindingNotAReDerivation` |
| REQ-43 | `TestReq43And46_0019_TheReaderIsTheFirstMatchInRegistryOrder` |
| REQ-44 | `TestReq44_0019_ANoMatchReachesTheUnboundRoleArmAndMintsNoMissingReaderCode` |
| REQ-45 | `TestReq45_0019_TheCarrierGateInvokesNoBindingMethod` |
| REQ-46 | `TestReq43And46_0019_TheReaderIsTheFirstMatchInRegistryOrder` |
| REQ-47 | `TestReq47_0019_EveryNonEmptyStoreShapeIsANoOpSuccess` |
| REQ-48 | `TestReq48And74And79_0019_AClearedKeyIsNeverResurrectedByRepeatedInitState` |
| REQ-49 | `TestReq49_0019_AnEmptiedStoreAndANeverWrittenOneSeedIdentically`<br>`TestReq49_0019_ClearingTheLastKeyReturnsTheStoreToTheInitializableClass` |
| REQ-50 | `TestReq50_0019_AKeyAddedToInitialAfterSeedingIsRepairedBySetStateNotInit` |
| REQ-51 | `TestReq51And93_0019_ATornStateIsReportedNotRepaired` |
| REQ-52 | `TestReq52_0019_AnOwnedKeyOutsideInitialStillReportsAbsentAfterInit` |
| REQ-53 | `TestReq53And83_0019_EmptyScalarSeedsPresentAndIsDistinctFromCleared` |
| REQ-54 | `TestReq54_0019_TheEmptyScalarIsNotRefusedAtSeedTimeAndAddsNoThirdEncoding` |
| REQ-55 | `TestReq55_0019_ArgvStillRefusesTheEmptyScalarWrite` |
| REQ-56 | `TestReq56_0019_ThePayloadDistinguishesTheSeededArmFromTheNoOpArm` |
| REQ-57 | `TestReq57_0019_TheSeededPayloadNamesKeysAndNeverEchoesTheirValues` |
| REQ-58 | `TestReq58_0019_TheSeededArmReportsNoNonEmptyAbsentKeyList` |
| REQ-59 | `TestReq59_0019_TheSuccessPayloadRidesTheShippedEnvelopeAndArtifactShape` |
| REQ-60 | `TestReq60_0019_ThePlanIsValidatedWholeBeforeAnyWriteCommits` |
| REQ-61 | `TestReq34And61And84_0019_TwoWritersForOneInitialKeyRefusesAtLoad` |
| REQ-62 | `TestReq31And32And62And85_0019_AnUnboundNeededRoleRefusesBeforeTheEmptinessRead` |
| REQ-63 | `TestReq63And67And87_0019_AReadBackMismatchRefusesAtExitTwoAndTheReRunDoesNotRepairIt` |
| REQ-64 | `TestReq64And65And68_0019_TheTwoNewCodesShareExitTwoAndAreDistinct` |
| REQ-65 | `TestReq64And65And68_0019_TheTwoNewCodesShareExitTwoAndAreDistinct` |
| REQ-66 | `TestReq66_0019_SharedRefusalClassesReuseTheirShippedCodesAndGroups` |
| REQ-67 | `TestReq63And67And87_0019_AReadBackMismatchRefusesAtExitTwoAndTheReRunDoesNotRepairIt`<br>`TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore` |
| REQ-68 | `TestReq64And65And68_0019_TheTwoNewCodesShareExitTwoAndAreDistinct` |
| REQ-69 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-70 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-71 | `TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts` |
| REQ-72 | `TestReq72_0019_LoaderOnlyKindsHaveNoArgvTranscriptionToCompare` |
| REQ-73 | `TestReq73And82_0019_TheFloatSpellingDivergenceIsDecidedBehaviour` |
| REQ-74 | `TestReq48And74And79_0019_AClearedKeyIsNeverResurrectedByRepeatedInitState` |
| REQ-75 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-76 | `TestReq76_0019_EveryReqCarriesACoverageRow`<br>`TestReqMVV10_0019_InitStateSeedsTheShippedRdrModel` |
| REQ-77 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-78 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-79 | `TestReq48And74And79_0019_AClearedKeyIsNeverResurrectedByRepeatedInitState` |
| REQ-80 | `TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts` |
| REQ-81 | `TestReq17And18And71And80And81_0019_InitAndSetStateProduceByteIdenticalArtifacts` |
| REQ-82 | `TestReq73And82_0019_TheFloatSpellingDivergenceIsDecidedBehaviour` |
| REQ-83 | `TestReq21And22And23And83_0019_LoaderOnlyKindsSeedThroughTheNormalizedPath`<br>`TestReq53And83_0019_EmptyScalarSeedsPresentAndIsDistinctFromCleared` |
| REQ-84 | `TestReq34And61And84_0019_TwoWritersForOneInitialKeyRefusesAtLoad` |
| REQ-85 | `TestReq31And32And62And85_0019_AnUnboundNeededRoleRefusesBeforeTheEmptinessRead` |
| REQ-86 | `TestReq4And6And86_0019_ADecisionTableModelRefusesRatherThanSucceedingVacuously` |
| REQ-87 | `TestReq63And67And87_0019_AReadBackMismatchRefusesAtExitTwoAndTheReRunDoesNotRepairIt` |
| REQ-88 | `TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess` |
| REQ-89 | `TestReq27And28And29And88And89_0019_ASealedStoreIsANonEmptyOneKeyStoreAndANoOpSuccess` |
| REQ-90 | `TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding` |
| REQ-91 | `TestReq38And39And40And90And91_0019_TheCarrierGateRefusesEveryNonFileBackedBinding` |
| REQ-92 | `TestReq25And26And92_0019_OneSurvivingKeyAnywhereBlocksTheWholeSeed` |
| REQ-93 | `TestReq51And93_0019_ATornStateIsReportedNotRepaired` |
| REQ-94 | `TestReq2And3And94_0019_NoReadPathSynthesizesTheInitialValueForAnAbsentKey` |
| REQ-95 | `TestReq95_0019_HelpAllNamesInitState` |
| REQ-96 | `TestReq96_0019_ModelAuthoringStatesTheFileBackedScope` |
| REQ-97 | `TestReq97_0019_NextPayloadGainsNoInitStateField` |
| REQ-98 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-99 | `TestReq99And100_0019_LintGainsNoArm`<br>`TestReq99And13_0019_SetStateGainsNoFromInitialFlag` |
| REQ-100 | `TestReq99And100_0019_LintGainsNoArm` |
| REQ-101 | `TestReq101_0019_TheRefusedCarriersStayRefused` |
| REQ-102 | `TestReq102_0019_NoPerRoleRepairGrammarIsOffered` |
| REQ-103 | `TestReq103_0019_TheNewRefusalsCarryNoPerFindingIdentityScheme` |
| REQ-MVV | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.1 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.2 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.3 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.4 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.5 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.6 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.7 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.8 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.9 | `TestReqMVV_0019_InitStateAcceptanceSpine` |
| REQ-MVV.10 | `TestReqMVV10_0019_InitStateSeedsTheShippedRdrModel`<br>`TestReqMVV_0019_InitStateAcceptanceSpine` |

**Orphans:** 0 — none; every REQ carries at least one test.


## REQ-MVV runner

```sh
B=./bin/intrastate; M=models/rdr.toml; A="$(mktemp -d)/rdr.artifact"
"$B" flow init-state  --model "$M" --artifact "rdr=$A" --as=json
"$B" flow read-state  --model "$M" --artifact "rdr=$A" --as=json
"$B" flow next        --model "$M" --artifact "rdr=$A" --as=json \
  | jq -c '[.data.candidates[].unknown[] | select(.reason=="absent")]'
"$B" flow init-state  --model "$M" --artifact "rdr=$A" --as=json          # re-run
"$B" flow set-state   --model "$M" --artifact "rdr=$A" --clear status --as=json
"$B" flow init-state  --model "$M" --artifact "rdr=$A" --as=json          # after clear
```

## REQ-MVV output

Run against the SHIPPED `models/rdr.toml` (REQ-MVV.10, ASSUMPTION-7) — the
user outcome, not a fixture. The artifact path is elided to `<fresh>` /
`<same>`; every other byte is verbatim.

```console
$ intrastate flow init-state --model models/rdr.toml --artifact rdr=<fresh> --as=json
{"type":"ok","data":{"model":"models/rdr.toml","revision":"","artifacts":{"rdr":"<fresh>"},"seeded":["gate_passed","stage","status"],"absent":[]}}

$ intrastate flow read-state --model models/rdr.toml --artifact rdr=<same> --as=json
{"type":"ok","data":{"model":"models/rdr.toml","revision":"","artifacts":{"rdr":"<same>"},"readers":[{"id":"rdr-status","keys":["stage","status","gate_passed"],"tags":{"gate_passed":"false","stage":"seeded","status":"draft"}}]}}

$ intrastate flow next --model models/rdr.toml --artifact rdr=<same> --as=json | jq -c '[.data.candidates[].unknown[] | select(.reason=="absent")]'
[]

$ intrastate flow init-state --model models/rdr.toml --artifact rdr=<same> --as=json   # re-run
{"type":"ok","data":{"model":"models/rdr.toml","revision":"","artifacts":{"rdr":"<same>"},"seeded":[],"absent":[]}}

$ intrastate flow set-state --model models/rdr.toml --artifact rdr=<same> --clear status --as=json
{"type":"ok","data":{"model":"models/rdr.toml","revision":"","artifacts":{"rdr":"<same>"},"writers":["rdr-status"],"writes":{},"clear":["status"],"owned":{}}}

$ intrastate flow init-state --model models/rdr.toml --artifact rdr=<same> --as=json   # after clear
{"type":"ok","data":{"model":"models/rdr.toml","revision":"","artifacts":{"rdr":"<same>"},"seeded":[],"absent":["status"]}}
```

Reading it against the ten steps: step 2's seeding arm names all three
`[initial]` keys and an EMPTY absent list; step 3 reads every one back at
its declared value; step 4's `absent` filter is empty, so the first-run wall
is gone; step 5's re-run reports the no-op arm (`seeded: []`) and wrote
nothing; steps 6–7 clear `status` and the next run reports it as an absent
`[initial]` key WITHOUT resurrecting it. Steps 1, 8 and 9 are asserted by
`TestReqMVV_0019_InitStateAcceptanceSpine` over the fixture model, which
carries the decision-table, unbound-role and carrier arms step 8 needs.

**Phase 2 status.** Every REQ above is green except REQ-63 / REQ-67 /
REQ-87's read-back MISMATCH arm, which is unreachable through the
file-backed carrier — see `deviations.md` DEV-9 (SPEC-DEFECT, needs author
decision). The exit-3 INCOMPLETE half of REQ-67 is green
(`TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore`).
