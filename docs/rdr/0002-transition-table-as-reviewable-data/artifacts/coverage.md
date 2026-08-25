# Coverage — RDR 0002 Transition Table As Reviewable Data

Phase 1 artifact. Every REQ in
[`req-list.md`](req-list.md) maps to at least one named test in
`internal/table/*_test.go`, and every test cites exactly one REQ. The
mapping is 1:1 in both directions: **147 REQs, 147 tests, zero orphans.**

Tests are external (`package table_test`), so they exercise only the
exported surface. Nothing mocks the unit under test: every case reads a
promoted fixture off disk and hands the bytes to the real `table.Load`,
and the resolver legs run the real `resolve.Resolve`. The one stub is
`fixtureGuards` in `mvv_test.go`, which stands in for RDR 0003's guard
evaluator **seam** — a peer package, not the unit under test.

## Red gate

Confirmed by compiling and running the whole suite against a temporary
type-only stub (every function returning its zero value), then discarding
it:

| Result | Count |
| --- | --- |
| Top-level tests FAILING against a do-nothing implementation | **136** |
| Subtests FAILING | **243** |
| Top-level tests passing by SHAPE (see below) | 9 |

Without the stub the package does not compile at all —
`go vet ./internal/table/` reports `undefined: table.Model`. That is
legitimate red: the failure names a **missing production symbol**, not a
typo in a test. The stub pass confirmed the tests themselves are
typo-free, since they compiled cleanly against the intended surface.

**The nine that pass by shape** are pure structural assertions over the
package's own source or exported surface, for which no behavioural form
exists — they assert what the package must *be*, not what it must *do*,
and each would fail against a wrong implementation:

| Test | What a wrong implementation does to fail it |
| --- | --- |
| `TestReq8_LoadTakesBytesAndSourceIDAndDoesNoFileIO` | takes a path, or imports `os`/`io/fs`/`path/filepath` |
| `TestReq84_NormalizedRowCarriesOneUnifiedAtomSet` | pre-splits `Row` into `Match`/`Guard` |
| `TestReq104_ThisPackageRestatesNoSelectionProcedure` | declares a package-level `Resolve`/`Select`/`Gate` |
| `TestReq109_PackageDoesNotImportInternalCLI` | imports `internal/cli` |
| `TestReq110_CategoryIdentifiersAreSnakeCase` | spells a category in any other case or shape |
| `TestReq111_CrossRowFindingsAreNotLoadCategories` | mints an `overlap`/`gap`/`dead_row` load category |
| `TestReq116_NoSourceRewriteCapability` | exports `Rewrite`/`Marshal`/`Encode`/`Save` |
| `TestReq117_OverlapIsDeferredAndNotAssertedHere` | mints an overlap verdict before RDR 0006 lands |
| `TestReq140_PackagePlacement` | lives outside `internal/` |

## REQ x test map

| REQ | Label | Test | File |
| --- | --- | --- | --- |
| REQ-MVV | HAPPY PATH | `TestMVV_ParseNormalizeDumpValidateAndResolve` | `mvv_test.go` |
| REQ-1 | HAPPY PATH | `TestReq1_SparseSourceExpandsToMoreRowsThanAuthoredRules` | `format_test.go` |
| REQ-2 | HAPPY PATH | `TestReq2_ClosedLayoutDecodesEveryElement` | `format_test.go` |
| REQ-3 | HAPPY PATH | `TestReq3_EveryAuthoringSiteReachesTheNormalizedRow` | `format_test.go` |
| REQ-4 | HAPPY PATH | `TestReq4_DescriptionAndSourceAreAdmittedAndCarried` | `format_test.go` |
| REQ-5 | ADVERSARIAL | `TestReq5_NoOtherRootKeyOrTableIsAdmitted` | `format_test.go` |
| REQ-6 | ADVERSARIAL | `TestReq6_UnmappedKeyIsAStableRefusalNotASilentNoOp` | `format_test.go` |
| REQ-7 | BOUNDARY | `TestReq7_TOMLParserIsGoTOMLV2` | `format_test.go` |
| REQ-8 | BOUNDARY | `TestReq8_LoadTakesBytesAndSourceIDAndDoesNoFileIO` | `format_test.go` |
| REQ-9 | HAPPY PATH | `TestReq9_MetadataIsFreeFormAndUninterpreted` | `format_test.go` |
| REQ-10 | BOUNDARY | `TestReq10_MetadataReachesNoCandidateRow` | `format_test.go` |
| REQ-11 | ADVERSARIAL | `TestReq11_MetadataAssertedByValueAndNestingEquality` | `format_test.go` |
| REQ-12 | INPUT EDGE | `TestReq12_MetadataShapeIsUnconstrained` | `format_test.go` |
| REQ-13 | ADVERSARIAL | `TestReq13_ModelHeaderAndVersionOne` | `format_test.go` |
| REQ-14 | ADVERSARIAL | `TestReq14_VersionGateRunsBeforeStrictFieldValidation` | `format_test.go` |
| REQ-15 | BOUNDARY | `TestReq15_FixedPrecedenceOrdering` | `format_test.go` |
| REQ-16 | BOUNDARY | `TestReq16_LoadIsFailFastWithOneCategory` | `format_test.go` |
| REQ-17 | BOUNDARY | `TestReq17_LoadReturnsOneRefusalNotAList` | `format_test.go` |
| REQ-18 | BOUNDARY | `TestReq18_AbsorbedDefectsTripNoCategory` | `format_test.go` |
| REQ-19 | ADVERSARIAL | `TestReq19_AccessorDeclarationRequiresAllFourFields` | `accessors_test.go` |
| REQ-20 | ADVERSARIAL | `TestReq20_WriteEntryRequiresReadBackTrue` | `accessors_test.go` |
| REQ-21 | HAPPY PATH | `TestReq21_SameIDMayAppearUnderTwoCapabilityTables` | `accessors_test.go` |
| REQ-22 | ADVERSARIAL | `TestReq22_KeysMustBeDeclaredAndNeverRecognized` | `accessors_test.go` |
| REQ-23 | BOUNDARY | `TestReq23_ReaderArityIsProvenanceScoped` | `accessors_test.go` |
| REQ-24 | ADVERSARIAL | `TestReq24_EveryWrittenKeyIsServedByExactlyOneWriter` | `accessors_test.go` |
| REQ-25 | ADVERSARIAL | `TestReq25_WriterKeysMustBeOwned` | `accessors_test.go` |
| REQ-26 | ADVERSARIAL | `TestReq26_GateIDsResolveOnlyAgainstGateEntries` | `accessors_test.go` |
| REQ-27 | BOUNDARY | `TestReq27_ReaderArityAndWriterProvenanceCarryDistinctCategories` | `accessors_test.go` |
| REQ-28 | HAPPY PATH | `TestReq28_GateListIsCarriedOnTheNormalizedRow` | `accessors_test.go` |
| REQ-29 | ADVERSARIAL | `TestReq29_EscapeRuleMustNotCarryAGateList` | `accessors_test.go` |
| REQ-30 | ADVERSARIAL | `TestReq30_InitialAssignmentsAreValidatedBySite` | `accessors_test.go` |
| REQ-31 | ADVERSARIAL | `TestReq31_TerminalIDsMustResolve` | `accessors_test.go` |
| REQ-32 | BOUNDARY | `TestReq32_TerminalIsDereferencedToPredicateSets` | `accessors_test.go` |
| REQ-33 | INPUT EDGE | `TestReq33_AbsentInitialOrTerminalIsNotALoadFailure` | `accessors_test.go` |
| REQ-34 | ADVERSARIAL | `TestReq34_RuleIDsAreUniqueByExactByteEquality` | `rules_test.go` |
| REQ-35 | ADVERSARIAL | `TestReq35_OrdinaryRuleShape` | `rules_test.go` |
| REQ-36 | ADVERSARIAL | `TestReq36_EscapeRuleCarriesNoPlanBearingField` | `rules_test.go` |
| REQ-37 | BOUNDARY | `TestReq37_AWriteReplaces` | `rules_test.go` |
| REQ-38 | ADVERSARIAL | `TestReq38_SetWriteValueIsAMemberSequenceNotAJoinedString` | `rules_test.go` |
| REQ-39 | ADVERSARIAL | `TestReq39_EscapeListAdmitsOnlyTwoFailureClasses` | `rules_test.go` |
| REQ-40 | HAPPY PATH | `TestReq40_EscapeRuleNormalizesToACandidateRow` | `rules_test.go` |
| REQ-41 | BOUNDARY | `TestReq41_EscapeRowsBindOneOutcomeEachWithNoCatchAll` | `rules_test.go` |
| REQ-42 | HAPPY PATH | `TestReq42_EscapeRuleExpandsUnderIn` | `rules_test.go` |
| REQ-43 | BOUNDARY | `TestReq43_RowKindIsDerivedNotAField` | `rules_test.go` |
| REQ-44 | ADVERSARIAL | `TestReq44_RenderedKindDoesNotFeedBackIntoTheNormalizedValue` | `rules_test.go` |
| REQ-45 | HAPPY PATH | `TestReq45_InheritanceNormalizesToAnExplicitPredicateSet` | `rules_test.go` |
| REQ-46 | ADVERSARIAL | `TestReq46_MergeKeysOnTheFullAtomIdentity` | `rules_test.go` |
| REQ-47 | BOUNDARY | `TestReq47_MergeIsIdempotentOnIdenticalAtoms` | `rules_test.go` |
| REQ-48 | INPUT EDGE | `TestReq48_UnsatisfiableAtomPairIsNotALoadFailure` | `rules_test.go` |
| REQ-49 | HAPPY PATH | `TestReq49_OneUnifiedPredicateSetWithBlockRetained` | `atoms_test.go` |
| REQ-50 | BOUNDARY | `TestReq50_BlockDomainIsTheKernelsThreeValues` | `atoms_test.go` |
| REQ-51 | ADVERSARIAL | `TestReq51_BlocksAreNotFolded` | `atoms_test.go` |
| REQ-52 | INPUT EDGE | `TestReq52_SameAtomInAllAndUnlessSurvivesAsTwo` | `atoms_test.go` |
| REQ-53 | ADVERSARIAL | `TestReq53_ExistenceAtomsEmitTheKernelConstants` | `atoms_test.go` |
| REQ-54 | ADVERSARIAL | `TestReq54_TagKeyIdentityIsExactByteEqualityAtEveryStage` | `atoms_test.go` |
| REQ-55 | ADVERSARIAL | `TestReq55_ReferencesResolveToTheDeclarationKey` | `atoms_test.go` |
| REQ-56 | BOUNDARY | `TestReq56_SetLiteralIsAnOrderedMemberSequence` | `atoms_test.go` |
| REQ-57 | BOUNDARY | `TestReq57_MembersSortByteLexicographically` | `atoms_test.go` |
| REQ-58 | ADVERSARIAL | `TestReq58_LiteralComparesAsAMemberSequenceNotAJoinedString` | `atoms_test.go` |
| REQ-59 | BOUNDARY | `TestReq59_SetValuedFieldsRenderWithMembersUnambiguouslyDelimited` | `atoms_test.go` |
| REQ-60 | ADVERSARIAL | `TestReq60_ExpansionSeparatorIsBannedAtThreeSitesWithThreeCategories` | `atoms_test.go` |
| REQ-61 | ADVERSARIAL | `TestReq61_ClearSentinelIsReservedAtEveryAuthoringSite` | `atoms_test.go` |
| REQ-62 | HAPPY PATH | `TestReq62_EveryMatchedOrWrittenTagIsDeclaredWithProvenance` | `atoms_test.go` |
| REQ-63 | BOUNDARY | `TestReq63_TypeModelWireKeys` | `atoms_test.go` |
| REQ-64 | ADVERSARIAL | `TestReq64_DeclarationSiteAndTheTwoRejectionCategories` | `atoms_test.go` |
| REQ-65 | BOUNDARY | `TestReq65_EveryDeclaredFieldSurvivesNormalization` | `atoms_test.go` |
| REQ-66 | ADVERSARIAL | `TestReq66_RecognizedIsTheReservedName` | `atoms_test.go` |
| REQ-67 | ADVERSARIAL | `TestReq67_OutcomeAlphabetIsNonEmptyDuplicateFreeAndHasNoEmptyMember` | `atoms_test.go` |
| REQ-68 | HAPPY PATH | `TestReq68_EveryRuleBindsExactlyOneOutcomeFromItsMatchBlocks` | `normalize_test.go` |
| REQ-69 | BOUNDARY | `TestReq69_TheRecognizedAtomIsLiftedOutOfThePredicateSet` | `normalize_test.go` |
| REQ-70 | ADVERSARIAL | `TestReq70_ZeroTwoOrOutOfAlphabetOutcomeBindingsRefuse` | `normalize_test.go` |
| REQ-71 | ADVERSARIAL | `TestReq71_RecognizedAtomInAGuardBlockRefuses` | `normalize_test.go` |
| REQ-72 | HAPPY PATH | `TestReq72_EveryMatchBlockInExpandsAndTheyProduct` | `normalize_test.go` |
| REQ-73 | BOUNDARY | `TestReq73_ExpandedInBecomesEqStillCarryingBlockMatch` | `normalize_test.go` |
| REQ-74 | BOUNDARY | `TestReq74_ExpansionSuffixIsASequenceInAtomSortOrder` | `normalize_test.go` |
| REQ-75 | INPUT EDGE | `TestReq75_UnsatisfiableProductRowIsNotALoadFailure` | `normalize_test.go` |
| REQ-76 | BOUNDARY | `TestReq76_SingleMemberInIsOneSpellingOfOneEdge` | `normalize_test.go` |
| REQ-77 | HAPPY PATH | `TestReq77_RequiresOwnedIsDerivedFromWritesAndClears` | `normalize_test.go` |
| REQ-78 | BOUNDARY | `TestReq78_EscapeRowsCarryAnEmptyRequiresOwnedSet` | `normalize_test.go` |
| REQ-79 | ADVERSARIAL | `TestReq79_GuardReadKeysAreNotAddedToRequiresOwned` | `normalize_test.go` |
| REQ-80 | HAPPY PATH | `TestReq80_NextTagsAndWritesAreBothPopulatedIncludingClears` | `normalize_test.go` |
| REQ-81 | BOUNDARY | `TestReq81_EscapeRowsCarryEmptyNextTagsAndWrites` | `normalize_test.go` |
| REQ-82 | ADVERSARIAL | `TestReq82_NextTagsAndWritesAreNotAliased` | `normalize_test.go` |
| REQ-83 | BOUNDARY | `TestReq83_ClearingIsExplicitAndAbsenceIsNotDeletion` | `normalize_test.go` |
| REQ-84 | BOUNDARY | `TestReq84_NormalizedRowCarriesOneUnifiedAtomSet` | `normalize_test.go` |
| REQ-85 | ADVERSARIAL | `TestReq85_HandoffRoutesByBlockNeverByOperator` | `normalize_test.go` |
| REQ-86 | ADVERSARIAL | `TestReq86_MatchBlocksAdmitOnlyEqAndIn` | `normalize_test.go` |
| REQ-87 | ADVERSARIAL | `TestReq87_AdmittedOperatorSetIsClosedAtEight` | `normalize_test.go` |
| REQ-88 | ADVERSARIAL | `TestReq88_AtomLevelValidationIsBlockAgnostic` | `normalize_test.go` |
| REQ-89 | BOUNDARY | `TestReq89_ConformanceIsPerOperator` | `normalize_test.go` |
| REQ-90 | BOUNDARY | `TestReq90_TwoRulesAreLegitimatelyMatchOnly` | `normalize_test.go` |
| REQ-91 | HAPPY PATH | `TestReq91_EveryCandidateRowRetainsRuleIDAndLocator` | `dump_test.go` |
| REQ-92 | ADVERSARIAL | `TestReq92_LocatorIsDerivedFromModelAndRuleIDNeverFromSource` | `dump_test.go` |
| REQ-93 | BOUNDARY | `TestReq93_SetsCrossTheKernelSeamAsCanonicalJSONArrays` | `dump_test.go` |
| REQ-94 | HAPPY PATH | `TestReq94_DumpCarriesEveryFieldOfTheNormalizedValue` | `dump_test.go` |
| REQ-95 | BOUNDARY | `TestReq95_DumpColumnVocabularyIsClosedAndVerbatim` | `dump_test.go` |
| REQ-96 | ADVERSARIAL | `TestReq96_DumpOrderMayReorderButNeverOmit` | `dump_test.go` |
| REQ-97 | BOUNDARY | `TestReq97_DumpSettingsDoNotReachTheNormalizedValue` | `dump_test.go` |
| REQ-98 | BOUNDARY | `TestReq98_RowsSortByTheIdentityTuple` | `dump_test.go` |
| REQ-99 | BOUNDARY | `TestReq99_WithinRowSequencesAreSorted` | `dump_test.go` |
| REQ-100 | ADVERSARIAL | `TestReq100_LocatorDoesNotParticipateInRowOrdering` | `dump_test.go` |
| REQ-101 | ADVERSARIAL | `TestReq101_DuplicateModelIDIsCrossDocumentOnly` | `dump_test.go` |
| REQ-102 | ADVERSARIAL | `TestReq102_DumpIsEmittedFromAPreSortedSequence` | `dump_test.go` |
| REQ-103 | ADVERSARIAL | `TestReq103_NoPositionalFieldAndRuleOrderIsNeutralized` | `dump_test.go` |
| REQ-104 | BOUNDARY | `TestReq104_ThisPackageRestatesNoSelectionProcedure` | `dump_test.go` |
| REQ-105 | BOUNDARY | `TestReq105_MultiModelDumpsCarryTheModelID` | `dump_test.go` |
| REQ-106 | BOUNDARY | `TestReq106_EveryNamedCategoryExistsAndIsWitnessed` | `dump_test.go` |
| REQ-107 | ADVERSARIAL | `TestReq107_AbsentVersionIsNotUnsupportedVersion` | `dump_test.go` |
| REQ-108 | BOUNDARY | `TestReq108_ARefusedDocumentYieldsNoCandidateRows` | `dump_test.go` |
| REQ-109 | ADVERSARIAL | `TestReq109_PackageDoesNotImportInternalCLI` | `dump_test.go` |
| REQ-110 | BOUNDARY | `TestReq110_CategoryIdentifiersAreSnakeCase` | `dump_test.go` |
| REQ-111 | BOUNDARY | `TestReq111_CrossRowFindingsAreNotLoadCategories` | `dump_test.go` |
| REQ-112 | INPUT EDGE | `TestReq112_AmbiguousOverlapIsNotALoadFailure` | `dump_test.go` |
| REQ-113 | HAPPY PATH | `TestReq113_RoundTripValueIdentityAcrossAllThreePermutations` | `roundtrip_test.go` |
| REQ-114 | ADVERSARIAL | `TestReq114_LocatorPresenceAndRuleIdentifyingPartAreCompared` | `roundtrip_test.go` |
| REQ-115 | BOUNDARY | `TestReq115_SetValuedFieldsRenderVisiblyAsSets` | `roundtrip_test.go` |
| REQ-116 | BOUNDARY | `TestReq116_NoSourceRewriteCapability` | `roundtrip_test.go` |
| REQ-117 | BOUNDARY | `TestReq117_OverlapIsDeferredAndNotAssertedHere` | `roundtrip_test.go` |
| REQ-118 | BOUNDARY | `TestReq118_PromotedFixtureSetIsNotNarrowed` | `roundtrip_test.go` |
| REQ-119 | ADVERSARIAL | `TestReq119_OneFixturePerCategoryAssertedByCategory` | `roundtrip_test.go` |
| REQ-120 | BOUNDARY | `TestReq120_BothExistenceArmsProduceDistinctLiterals` | `roundtrip_test.go` |
| REQ-121 | ADVERSARIAL | `TestReq121_MisCasedReferenceFailsUnknownTag` | `roundtrip_test.go` |
| REQ-122 | BOUNDARY | `TestReq122_KeyOrderRuleOrderAndRepeatedDumpDeterminism` | `roundtrip_test.go` |
| REQ-123 | HAPPY PATH | `TestReq123_MetadataDeepEqualityAndAnnotationKeysDecode` | `roundtrip_test.go` |
| REQ-124 | ADVERSARIAL | `TestReq124_NoGoldenHashOfRenderedText` | `roundtrip_test.go` |
| REQ-125 | ADVERSARIAL | `TestReq125_TheNoAliasControlIsAMutationTest` | `roundtrip_test.go` |
| REQ-126 | ADVERSARIAL | `TestReq126_DistinctLiteralControlParameterizedOverFiveDelimiters` | `roundtrip_test.go` |
| REQ-127 | ADVERSARIAL | `TestReq127_WriteValueControlReadsTheNormalizedValue` | `roundtrip_test.go` |
| REQ-128 | BOUNDARY | `TestReq128_ByteIdenticalAtomsCollapseToExactlyOne` | `roundtrip_test.go` |
| REQ-129 | HAPPY PATH | `TestReq129_EscapeExpansionYieldsTwoRowsWithDistinctSuffixes` | `roundtrip_test.go` |
| REQ-130 | ADVERSARIAL | `TestReq130_HandoffRoutingIsTotalAndDisjoint` | `roundtrip_test.go` |
| REQ-131 | ADVERSARIAL | `TestReq131_RefusalsAreAssertedByCategoryNotMessageText` | `roundtrip_test.go` |
| REQ-132 | BOUNDARY | `TestReq132_TheSevenOwedCategoriesAreNowWitnessed` | `roundtrip_test.go` |
| REQ-133 | ADVERSARIAL | `TestReq133_GuardBlockControlsExistForEveryBlockAgnosticRule` | `roundtrip_test.go` |
| REQ-134 | HAPPY PATH | `TestReq134_BothNormativeFixturesLoadClean` | `roundtrip_test.go` |
| REQ-135 | ADVERSARIAL | `TestReq135_V2ShapedPrecedenceControl` | `roundtrip_test.go` |
| REQ-136 | HAPPY PATH | `TestReq136_PositiveControls` | `roundtrip_test.go` |
| REQ-137 | ADVERSARIAL | `TestReq137_ResolveOverTheNormalizedFixtureRows` | `mvv_test.go` |
| REQ-138 | ADVERSARIAL | `TestReq138_ScenarioSixCarriesTwoDistinctCategories` | `roundtrip_test.go` |
| REQ-139 | ADVERSARIAL | `TestReq139_DuplicateModelIDNeedsThePairedDocumentSurface` | `roundtrip_test.go` |
| REQ-140 | BOUNDARY | `TestReq140_PackagePlacement` | `roundtrip_test.go` |
| REQ-141 | BOUNDARY | `TestReq141_PhaseTwoDeliverables` | `roundtrip_test.go` |
| REQ-142 | BOUNDARY | `TestReq142_PhaseThreeValidationFamilies` | `roundtrip_test.go` |
| REQ-143 | BOUNDARY | `TestReq143_KernelHandoffAddsNoOrderingOrCallback` | `roundtrip_test.go` |
| REQ-144 | BOUNDARY | `TestReq144_ParsedRepresentationIsExposedForLint` | `roundtrip_test.go` |
| REQ-145 | BOUNDARY | `TestReq145_ExpansionCountsAreDerivedAndNonBlocking` | `roundtrip_test.go` |
| REQ-146 | ADVERSARIAL | `TestReq146_EveryEmittedSequenceIsASortedSlice` | `roundtrip_test.go` |

## REQ-MVV output

`TestMVV_ParseNormalizeDumpValidateAndResolve` — nine legs, one runnable
end-to-end test.

Because this RDR declares a Round-Trip / Inverse Invariant
(`parse ∘ normalize = expanded-table value identity`), leg 7 compares the
reconstructed candidate-row set **value-for-value** against the original
across all three semantics-preserving authorings (TOML key order, rule
declaration order, `eq` versus single-member `in`), and leg 3 compares
repeated dumps **byte-for-byte**. A green exit code or "did not error" is
not sufficient at any leg.

| Leg | Obligation | Oracle |
| --- | --- | --- |
| 1 | parse two fixtures into typed source data | field-level assertions on the type model, accessor entries, and metadata |
| 2 | normalize into candidate rows | the nine RDR and two kata row identities, compared as a value |
| 3 | dump the expanded table | every row present; eight repeated dumps byte-identical |
| 4 | validate the eight named families | one category per family, plus multi-tag writes asserted positively |
| 5 | one tag-set resolves to exactly one ordinary row | real `resolve.Resolve`; plan rule id and next-tags/writes compared value-for-value |
| 6 | one unmatched tag-set resolves to one modeled escape row | `Plan.Escaped`, the escape rule id, and empty next/writes |
| 7 | Round-Trip value identity | reconstructed row set `DeepEqual` to the original, all three permutations |
| 8 | unsupported version refused before normalization | `unsupported_version`, zero candidate rows, plus the A14 v2-shaped precedence control |
| 9 | overlapping variant reported by lint | **DEFERRED** — see below |

**Actual output**: _(unfilled — Phase 2 records the run.)_

**Leg 9 is deferred, not satisfied.** REQ-117 states outright that the
ambiguous-overlap check is cross-row and therefore RDR 0006's by the arity
split, that RDR 0006 is unimplemented, and that "marking the overlap
assertion green before RDR 0006 lands would be asserting on a stub." Leg 9
therefore asserts only what this RDR owes *before* that handshake: the
overlapping variant **loads** (overlap is not a load failure, REQ-112), and
its rows are distinguishable by row identity and comparable by predicate
set, which is the structure RDR 0006's check consumes.

## Orphans

**REQs with no test: none.** All 146 numbered REQs plus REQ-MVV are cited.

**Tests citing no REQ: none.** Every one of the 147 `Test*` functions opens
with a `// REQ-N: "<quote>"` comment and exactly one label. The remaining
package-level functions are unexported helpers (`readFixture`, `mustLoad`,
`loadCategory`, `rowByID`, `rowsByRuleID`, `rowIdentities`, `atomsOn`,
`atomsInBlock`, `tagValue`, `containsAtom`, `packageGoFiles`,
`setWriteValue`, `escapeRuleBlock`, `hasGuardAtom`, `compareAtoms`,
`assertSortedByKey`, `dumpLine`, `seamValue`, `ownedTags`) and the
`fixtureGuards` seam stub; helpers carry no REQ by design.

## Fixture promotion

`internal/table/testdata/` holds the Stage-4/6 approved set from
`evidence/spikes/iter-2/`, promoted per REQ-118 — **extended, never
narrowed**. `TestReq118_PromotedFixtureSetIsNotNarrowed` enforces that by
walking both trees and comparing by filename, and additionally asserts the
approved row census (nine RDR rows, two kata rows).

Extensions minted this phase, all single-mutation derivatives of
`rdr-fixture.toml`:

| Fixture | Why |
| --- | --- |
| `neg/neg-escape-with-empty-write.toml` | deviations.md **D2** — `write = []` on an escape rule, refused by key PRESENCE not length |
| `neg/neg-write-value-outside-domain.toml`, `neg/neg-write-value-wrong-kind.toml` | deviations.md **D3** — write-block value kind/domain conformance |
| `neg/neg-malformed-toml.toml` | owed category: `malformed_toml` |
| `neg/neg-no-alphabet.toml` | owed category: `missing_recognized_outcome_alphabet` |
| `neg/neg-cyclic-context.toml` | owed category: `cyclic_context_inheritance` |
| `neg/neg-tagdecl-domain-on-bool.toml` | owed category: `malformed_tag_declaration` |
| `neg/neg-no-version.toml`, `neg/neg-no-model-id.toml`, `neg/neg-no-model-table.toml`, `neg/neg-no-recognized-decl.toml` | owed category: `malformed_model_declaration`, four arms |
| `neg/neg-dump-unknown-column.toml`, `neg/neg-dump-repeated-column.toml`, `neg/neg-dump-omitted-column.toml` | owed category: `malformed_dump_declaration`, three arms |
| `dup/dup-model-a.toml`, `dup/dup-model-b.toml`, `dup/dup-model-distinct.toml` | owed category: `duplicate_model_id` — the paired-document surface no single file can provide |

**deviations.md D1 check 1 is discharged**: `kind = "string"` was rewritten
to RDR 0003's `scalar` token across all 71 sites in the promoted fixtures;
`grep -rn 'kind = "string"' internal/table/testdata/` returns 0.
`TestReq64` additionally asserts the loader refuses `kind = "string"` as
`malformed_tag_declaration` and admits exactly the five tokens
`enum`, `bool`, `int`, `set`, `scalar` (`0003:793`).

## Assumptions carried into the tests

Each is grounded and recorded here rather than left implicit. The first
seven restate `req-list.md`'s ASSUMPTIONS as they bind a specific test; the
rest are Phase 1's own, derived from the spec where the REQ text left a
testable detail open.

1. **The package is `internal/table`.** REQ-140/EIA says "New internal
   package can own sparse source structs and normalization" without naming
   it. `table` is taken from the RDR's own Naming decision, which fixes
   "expanded transition table" as the canonical rendered view and the
   Existing Infrastructure Audit's "Normalizer/table package" row.

2. **`Load(src []byte, sourceID string) (*Model, error)`.** REQ-8 fixes
   the two parameters and REQ-17 forbids returning a list instead of a
   refusal; `(value, error)` is the only Go shape satisfying both.
   `TestReq8` and `TestReq17` assert the signature reflectively so a later
   phase cannot silently widen it.

3. **Refusals carry their category through `table.CategoryOf(err) (Category, bool)`.**
   REQ-110 says the identifiers are "an API surface, not message text" and
   REQ-131's oracle is the category. An accessor over `error` is the
   minimal surface that makes the category assertable without fixing the
   error type, which the record does not fix.

4. **`Atom.Literal` is `[]string`.** REQ-56/REQ-58 forbid a
   delimiter-joined string "in the stored value or in any comparison
   derived from it", and REQ-58 requires element-by-element comparison. A
   member sequence is the only Go shape that satisfies both, and
   `TestReq56` asserts the field's kind reflectively.

5. **REQ-38 and REQ-93 are two layers, not a contradiction.** The
   normalized value holds a `[]string` member sequence; the *kernel seam*
   encoding — what lands in `resolve.Tag.Value`, a `string` — is §D13's
   canonical JSON array. `TestReq93` asserts both halves in one test,
   including that the normalized value is NOT the serialized form.

6. **`duplicate_model_id` is a second, set-taking entry point,
   `CheckModelIDs([]*Model) error`.** REQ-101 says a loader handed one
   document MUST NOT report it and that the check belongs to "a caller
   that loads several documents into one dump or lint invocation".
   `TestReq139` asserts both documents load clean alone.

7. **`[model.metadata]` is `map[string]any`.** A1's evidence names exactly
   this shape, and REQ-12 forbids constraining it.

8. **Row identity renders as `<model id>.<rule id>[#<suffix elements>]`.**
   REQ-98 fixes the *tuple* and REQ-60 reserves `#` as the
   expansion-suffix separator; the spike output (`output.txt`) spells it
   this way and REQ-105 requires model qualification. `Row.Identity()`
   exposes it and `Row.Suffix` exposes the sequence separately, so no test
   parses the rendered string to recover the tuple.

9. **`Row.Kind()` is a method, not a field.** REQ-43 says row kind is "a
   derived view property, never a row field"; `TestReq43` asserts no
   `Kind` *field* exists and that the derived value agrees with the
   escape-list predicate on every row of both fixtures.

10. **Q1 is followed under reading (a)** — `Row.Match []Tag` stays the
    equality pattern and `Row.Guard []GuardAtom` takes only `all`/`unless`
    atoms. `TestReq130` therefore asserts the union across two Go fields of
    different types. This mirrors what `0007/artifacts/req-list.md` Q1
    chose and what `internal/resolve/resolve.go` currently declares. If a
    later phase finds `Match` cannot carry a `BlockMatch`-tagged atom's
    operator token, the routing contract (REQ-85) is unchanged but
    `TestReq130`'s assertion shape is.

11. **Q2 is followed by SITE, not by defect** — an `[initial]` value
    outside its declared domain is `malformed_initial_declaration`
    (`TestReq30`), while the same defect on a predicate literal is
    `malformed_predicate_atom` (`TestReq64`). The TS-3 fixture names
    (`neg-initial-bad-value`, `neg-initial-int-out-of-range`) confirm the
    `[initial]` site, and the promoted transcript records exactly those
    categories.

12. **`Model.Terminal` is `[][]Atom`.** REQ-32 requires normalization to
    dereference each terminal id to "the context's explicit predicate set
    over owned tags, in the same atom shape rules use" and to carry the
    predicate sets "never the bare ids". A slice of atom slices is that
    shape; `TestReq32` compares against the dereferenced value, not a
    non-empty check, so a normalizer carrying bare ids fails.

13. **`Dump` and `DumpAll` render; they claim no inverse.** REQ-94 requires
    a dump carrying every field and REQ-105 a model-qualified multi-model
    dump. `TestReq116` bans only genuinely *source-mutating* names
    (`Rewrite`, `Marshal`, `Encode`, `Save`, …) — the Round-Trip clause is
    explicit that the dump has no specified inverse and that a re-readable
    dump grammar is a follow-up RDR.

14. **REQ-39's non-modelable escape classes refuse as
    `malformed_escape_declaration`.** REQ-106's list names that category
    for "a write block, clear list, or `gate` list on an escape rule" and
    names no other for a bad class member; the `escape` list is the escape
    declaration, so its own malformity files there. No fenced category
    covers it separately and the record mints none.

15. **REQ-19's `keys`-member arm carries `unknown tag`, not the accessor
    category.** REQ-19's own text parenthesizes it: "a `keys` member naming
    an undeclared tag (`unknown tag`)". The promoted transcript's
    `neg-recognized-in-keys` line confirms the sibling reserved-key arm
    carries `malformed accessor binding` instead, so the two arms are
    tested apart.

## Open items

**None that the evidence cannot settle.** The two `req-list.md` QUESTIONS
were already resolved there as assumptions and are carried above as items
10 and 11; both are recorded rather than closed only because the evidence
is one-sided rather than conclusive, and both are single-test-shaped, so a
later phase that finds otherwise changes an assertion's shape, not a
contract.

---

## Phase 2 — REQ-MVV end-to-end result (recorded)

Run after the suite went green.

```
$ go test ./internal/table/ -run TestMVV_ParseNormalizeDumpValidateAndResolve -v
--- PASS: TestMVV_ParseNormalizeDumpValidateAndResolve (0.01s)
    --- PASS: .../1_parse_two_fixtures_into_typed_source_data (0.00s)
    --- PASS: .../2_normalize_into_candidate_rows (0.00s)
    --- PASS: .../3_dump_the_expanded_table (0.00s)
    --- PASS: .../4_validate_the_eight_named_families (0.00s)
    --- PASS: .../5_one_sample_tag_set_resolves_to_exactly_one_ordinary_row (0.00s)
    --- PASS: .../6_one_unmatched_tag_set_resolves_to_one_modeled_escape_row (0.00s)
    --- PASS: .../7_round_trip_value_identity (0.00s)
    --- PASS: .../8_unsupported_version_is_refused_before_normalization (0.00s)
    --- PASS: .../9_overlap_is_deferred_to_the_phase_5_lint_handshake (0.00s)
PASS
ok      github.com/newcoinc/intrastate/internal/table   0.170s
```

Leg 9 is DEFERRED by REQ-117, not counted satisfied: it asserts only what
this RDR owes before RDR 0006's lint handshake — the overlapping rows are
distinguishable by identity, bind one outcome, and are comparable by
predicate set.

**Whole-package state**: `go test ./...` green; `internal/table` runs 147
top-level tests / 467 including subtests, 0 failures, stable under
`-count=3`. `go vet ./...` clean; `gofmt -l` empty.

### The normalized dump, as produced

The approved row census and identities (nine RDR rows, two kata rows) hold
value-for-value against `evidence/spikes/iter-2/output.txt`. Two fields
differ from the spike's bytes, both because the record names the spike's
form a defect rather than a permitted reading:

- **The locator is derived**, `<model id>:<rule id>`, not copied from the
  authored `source`. `0002:TD` says so and calls the spike's
  `Locator = rule.Source` "a spike defect against this clause"; REQ-92's
  control edits `source` and asserts the locator does not move.
- **A set-valued field renders bracketed and quoted** — `labels=["needs
  work"]`, not `labels=needs work`. `0002:C10` requires each member
  delimited unambiguously and `0002:RT` requires a set to render visibly
  as a set. Set-ness is the DECLARED kind, so a one-member set still
  brackets; otherwise read-back could not tell `["a"]` from the scalar `a`.

```
MODEL rdr rows=9 outcomes=round-clean,verdict-flapping,reconcile-block,finalized
identity=rdr.continue-prelock#foundational source=rdr:continue-prelock kind=transition outcome=round-clean atoms=[finalized_at.exists=false@all; iter.lt=3@all; profile.eq=foundational@match; profile.eq=small@unless; stage.eq=prelock@match; status.eq=Draft@match] next=[iter=2; stage=prelock] writes=[iter=2; stage=prelock] requires_owned=[iter,stage] gate=[rdr-lock] escape=[]
identity=rdr.continue-prelock#large source=rdr:continue-prelock kind=transition outcome=round-clean atoms=[finalized_at.exists=false@all; iter.lt=3@all; profile.eq=large@match; profile.eq=small@unless; stage.eq=prelock@match; status.eq=Draft@match] next=[iter=2; stage=prelock] writes=[iter=2; stage=prelock] requires_owned=[iter,stage] gate=[rdr-lock] escape=[]
identity=rdr.continue-prelock-cluster#foundational source=rdr:continue-prelock-cluster kind=transition outcome=round-clean atoms=[cluster_ready.eq=true@all; profile.eq=foundational@match; stage.eq=prelock@match; status.eq=Draft@match] next=[prelock_lens=critique; stage=prelock] writes=[prelock_lens=critique; stage=prelock] requires_owned=[prelock_lens,stage] gate=[] escape=[]
identity=rdr.continue-prelock-cluster#large source=rdr:continue-prelock-cluster kind=transition outcome=round-clean atoms=[cluster_ready.eq=true@all; profile.eq=large@match; stage.eq=prelock@match; status.eq=Draft@match] next=[prelock_lens=critique; stage=prelock] writes=[prelock_lens=critique; stage=prelock] requires_owned=[prelock_lens,stage] gate=[] escape=[]
identity=rdr.draft-no-match-escape#reconcile-block source=rdr:draft-no-match-escape kind=escape outcome=reconcile-block atoms=[status.eq=Draft@match] next=[] writes=[] requires_owned=[] gate=[] escape=[no_match]
identity=rdr.draft-no-match-escape#round-clean source=rdr:draft-no-match-escape kind=escape outcome=round-clean atoms=[status.eq=Draft@match] next=[] writes=[] requires_owned=[] gate=[] escape=[no_match]
identity=rdr.reconcile-rewind source=rdr:reconcile-rewind kind=transition outcome=reconcile-block atoms=[status.eq=Draft@match] next=[prelock_lens=<clear>; rewind_scope=assumptions; stage=resolve; status=Draft] writes=[prelock_lens=<clear>; rewind_scope=assumptions; stage=resolve; status=Draft] requires_owned=[prelock_lens,rewind_scope,stage,status] gate=[] escape=[]
identity=rdr.terminal-archive#finalized source=rdr:terminal-archive kind=transition outcome=finalized atoms=[stage.eq=archive@match] next=[stage=archive] writes=[stage=archive] requires_owned=[stage] gate=[] escape=[]
identity=rdr.terminal-archive#verdict-flapping source=rdr:terminal-archive kind=transition outcome=verdict-flapping atoms=[stage.eq=archive@match] next=[stage=archive] writes=[stage=archive] requires_owned=[stage] gate=[] escape=[]
MODEL kata rows=2 outcomes=accepted,needs-work,closed
identity=kata.review-accepted source=kata:review-accepted kind=transition outcome=accepted atoms=[owner.eq=current-session@all; phase.eq=review@match; status.eq=open@match; status.eq=closed@unless] next=[phase=ship; status=accepted] writes=[phase=ship; status=accepted] requires_owned=[phase,status] gate=[] escape=[]
identity=kata.review-needs-work source=kata:review-needs-work kind=transition outcome=needs-work atoms=[phase.eq=review@match; status.eq=open@match] next=[labels=["needs work"]; phase=resolve; status=open] writes=[labels=["needs work"]; phase=resolve; status=open] requires_owned=[labels,phase,status] gate=[] escape=[]
```

### Pre-seeded deviation checks, discharged

- **D1** — discharged **at the promoted set**, which is what the loader
  reads: `grep -rn 'kind = "string"' internal/table/testdata/` → 0, and
  `TestReq64` asserts an unknown kind token refuses
  `malformed_tag_declaration` while all five of RDR 0003's tokens are
  admitted (`table.IsDeclaredKind`). The *spike* directory still carries
  71 hits, because it is Stage-4/6 evidence of what was reviewed then and
  this build is read-only over it; the rewrite D1 asks for landed on
  promotion (REQ-118: extended, never narrowed), which `TestReq118`
  enforces by name. Re-running `gen-cases.py` to refresh the evidence is
  the residual, and it changes no clause and no assertion.
- **D2** — `neg-escape-with-empty-write.toml` refuses
  `malformed_escape_declaration`: `sourceRule.Write` is a POINTER, so the
  loader keys on key presence rather than `len(...) > 0`.
- **D3** — `neg-write-value-outside-domain` and
  `neg-write-value-wrong-kind` both refuse `malformed_tag_declaration`, the
  category D3's check names. No fenced category had to widen.
- **D5** — the §D13 landing is implemented at the kernel seam
  (`model.go::seamValue`): sorted, duplicate-free, compact JSON array,
  matching the form RDR 0007 already wrote its `contains` leg against.
  `TestReq93` asserts all four legs including read-back byte equality.

### New deviations minted this phase

**D6** and **D7**, both TEST-FIXTURE, both `Status: mechanical
translation` — two Phase-1 assertions that were red against every possible
implementation. Each was repaired at the fixture, never by weakening the
assertion; see `deviations.md`.

**Open author decisions: none.**
