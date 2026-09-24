# Coverage — RDR 0030 Counting and stepping a tag without writing out every value by hand

Phase 1 test-author artifact. One row per REQ in `req-list.md`; column 2 is the
Go test that fails if the clause breaks. Every listed test was run and confirmed
RED against the pre-implementation tree (base `48049f9`) before commit.

| REQ | test |
| --- | --- |
| REQ-1 | TestReq1_0030_OnlyTheOneKeyNonZeroIntegerStepTableIsAdmitted |
| REQ-2 | TestReq2_0030_StepIsAdmittedOnlyOnBoundedIntAndNonEmptyEnum |
| REQ-3 | TestReq3_0030_AnUnrepresentableIntWidthIsRefusedForAStep |
| REQ-4 | TestReq4_0030_OneIntWidthRuleLivesInResolveAndGuardReachesIt |
| REQ-5 | TestReq5_0030_TheEmittedCellDenotesExactlyOneLintCell |
| REQ-6 | TestReq6_0030_AWideRepresentableDomainLoadsWithNoRowCeiling |
| REQ-7 | TestReq7_0030_AStepOverADuplicateOrHashCarryingDomainIsRefused |
| REQ-8 | TestReq8_0030_DuplicateAndHashDomainsStillLoadWhenTheyAreNotStepped |
| REQ-9 | TestReq9_0030_TheTableShapeIsWriteBlockOnlyAndEachPathKeepsItsCategory |
| REQ-10 | TestReq10_0030_EveryRefusalThisRecordMintsIsMalformedTagDeclaration |
| REQ-11 | TestReq11_0030_AdmittedCellsAreTheConjunctionOfEveryPositiveAtomOnTheTag |
| REQ-12 | TestReq12_0030_AStepSubsumesAMatchInOnItsTagAuthoredOrInherited |
| REQ-13 | TestReq13_0030_ASubsumedInIsRewrittenAndEveryOtherAtomIsRetained |
| REQ-14 | TestReq14_0030_TheShimIsAConcreteEvaluatorHomedInResolve |
| REQ-15 | TestReq15_0030_TheLoaderAdmitsACellOnlyOnGuardTrue |
| REQ-16 | TestReq16_0030_EachRowCarriesItsCellGuardItsLiteralWriteAndItsSuffix |
| REQ-17 | TestReq17_0030_TheCellIsAppendedEvenWhenExactlyOneIsAdmitted |
| REQ-18 | TestReq18_0030_TheSteppedValueIsComputedWithoutOverflowInEitherDirection |
| REQ-19 | TestReq19_0030_AStepRuleAdmittingNoCellIsRefused |
| REQ-20 | TestReq20_0030_SteppingAndClearingOneKeyIsRefusedBeforeAnyCellIsWalked |
| REQ-21 | TestReq21_0030_EveryExpandedRowRequiresItsSteppedKeyOwned |
| REQ-22 | TestReq22_0030_ExpandStaysTotalWhileMintingSteppedRows |
| REQ-23 | TestReq23_0030_WritesAndNextTagsBothCarryTheSteppedValueWithoutAliasing |
| REQ-24 | TestReq24_0030_EveryExpandedRowCarriesTheRulesBlockIdentically |
| REQ-25 | TestReq25_0030_AnExpandedRowIsIndistinguishableFromItsLiteralTwin |
| REQ-26 | TestReq26_0030_SuffixElementsOrderByKeyAlone |
| REQ-27 | TestReq27_0030_ASteppedValueLeavingTheDomainEitherWayIsRefused |
| REQ-28 | TestReq28_0030_TheResultSlotIsTheLiteralForIntAndTheOffsetForEnum |
| REQ-29 | TestReq29_0030_TheBoundRefusalRidesDetailAloneAndNamesTheCellFirst |
| REQ-30 | TestReq30_0030_TheFirstFailingCellInDomainOrderIsNamed |
| REQ-31 | TestReq31_0030_TheCheckedAddRefusesAnOverflowingCellBeforeRender |
| REQ-32 | TestReq32_0030_TheBoundRefusalNamesAnUnconsultedUnlessAtom |
| REQ-33 | TestReq33_0030_NoConsumerLearnsAComputedValueShape |
| REQ-34 | TestReq34_0030_NoEmittedVocabularyGainsAMember |
| REQ-35 | TestReq35_0030_TheAuthoringGuideDocumentsTheStepFormAndItsHazards |
| REQ-36 | TestReq36_0030_RowSurfacesNameRowsAndRuleSurfacesNameRules |
| REQ-37 | TestReq37_0030_EveryRowCarryingASuffixPublishesItInBothRowSlots |
| REQ-38 | TestReq38_0030_ReadBackJoinsOnTheRecoveredAuthoredID |
| REQ-39 | TestReq39_0030_TheOutputContractNarrowsRuleToTheRowIdentity |
| REQ-40 | TestReq40_0030_ThreeCallSitesReachOneComparisonInResolve |
| REQ-41 | TestReq41_0030_EachCallerKeepsItsOwnUndecidedArmPolicy |
| REQ-42 | TestReq42_0030_ASetLiteralsTwoSpellingsAdmitTheSameCells |
| REQ-43 | TestReq43_0030_AnUnclaimedCapAndAnUninitialisedSteppedTagSurfaceAtLint |
| REQ-44 | TestReq44_0030_AppendingAMemberNeitherRefusesNorWidensTheLadder |
| REQ-45 | TestReq45_0030_AnUnlessExcludedInteriorCellMintsADeadRow |
| REQ-46 | TestReq46_0030_RepeatedLoadsOfASteppedModelAreByteIdentical |
| REQ-47 | TestReq47_0030_AModelThatStepsNothingIsUnchangedButForRowAttribution |
| REQ-48 | TestReq48_0030_LintFindingSetsAgreeUnderTheProjection |
| REQ-49 | TestReq49_0030_ResolvePlansAgreeOverEveryCell |
| REQ-50 | TestReq50_0030_TheUnguardedLadderRefusesNamingRuleTagCellAndResult |
| REQ-51 | TestReq51_0030_TableShapesOffTheWritePathRefuseUnderTheirOwnCategories |
| REQ-52 | TestReq52_0030_EveryStepGrammarAndKindDefectIsANamedLoadRefusal |
| REQ-53 | TestReq53_0030_TheStepLaddersGraphExportDecodesStrictly |
| REQ-54 | TestReq54_0030_ANonFirstExpandedRowJoinsTheAuthoredEmitBlock |
| REQ-55 | TestReq55_0030_DumpPublishesTheSubsumedInRewrittenPerRow |
| REQ-56 | TestReq56_0030_RowSurfacesPublishTheSuffixedIdentity |
| REQ-57 | TestReq57_0030_BothCarriersPublishEachRowsSteppedSuccessor |
| REQ-58 | TestReq58_0030_PerRowFindingsNameTheCellAndGroupFindingsStayBare |
| REQ-MVV | TestMVV_0030_TheStepLadderIsTheLiteralLadderToEveryConsumer |

## REQ-MVV output

Phase 3c leg 2, tree `b9c340f`, full suite (`make check`) exit 0; the output below is unchanged from Phase 2 leg 1 (`d4128bb`).

```
=== RUN   TestMVV_0030_TheStepLadderIsTheLiteralLadderToEveryConsumer
--- PASS: TestMVV_0030_TheStepLadderIsTheLiteralLadderToEveryConsumer (0.17s)
    --- PASS: .../item_1_—_both_fixtures_load (0.00s)
    --- PASS: .../item_2_—_lint_agrees,_overlaps_and_gaps_are_zero,_the_cell_is_named (0.02s)
    --- PASS: .../item_3_—_flow_resolve_agrees_over_every_cell (0.12s)
    --- PASS: .../item_4_—_the_unguarded_cap_is_a_named_load_refusal (0.00s)
    --- PASS: .../item_5_—_the_dead_row_is_advisory_and_the_positive_exclusion_does_not_overlap (0.02s)
    --- PASS: .../item_6_—_the_graph_export_decodes_strictly (0.00s)
    --- PASS: .../item_7_—_one_rule_per_intent,_fewer_than_the_literal_ladder (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/cli	0.471s
```

`intrastate lint --model internal/table/testdata/ladder-step.toml --as json` (first finding):

```
{"code":"graph-idempotent-write","message":"row \"retry\" matches \"open\" on \"status\" and writes the same value back, so the write moves nothing on that key","model":"ladder","severity":"info","rule":"retry#0","span":"ladder:retry","dimension":"status","key":"status","operator":"eq","literal":"open","block":"match","fingerprint":"attempt|all|eq|0,;attempt|all|gte|0,;attempt|all|lt|5,;status|match|eq|open,;#attempt=1,;status=open,;"}
```

`intrastate lint --model internal/table/testdata/ladder-step-unguarded.toml --as json`:

```
{"code":"model-invalid","message":"model does not conform to the transition-model schema","schema_version":"0.1","param":"model","detail":"malformed_tag_declaration: rule retry write attempt: cell 9 steps to 10, which is above max 9","findings":[{"code":"malformed_tag_declaration","message":"malformed_tag_declaration: rule retry write attempt: cell 9 steps to 10, which is above max 9","locator":"internal/table/testdata/ladder-step-unguarded.toml:1"}]}
```

## REQ-MVV runner

`go test -count=1 -v -run 'TestMVV_0030_' ./internal/cli/` plus the two `intrastate lint --model … --as json` calls above, from the worktree root after `make build`.
