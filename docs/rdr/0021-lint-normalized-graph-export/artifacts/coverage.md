# Coverage — RDR 0021 Lint's normalized-graph export

Phase 1 artifact: every REQ-N from `req-list.md` against the test that would
fail if a future change broke that clause. An uncovered REQ keeps its row
with column 2 EMPTY.

REQ-18 carries no row: it was DEMOTED to `req-list.md`'s `## EXCLUDED`
section, because the clause has no observable of its own — its positive half
is REQ-19…REQ-25 and its negative half is discharged by no test asserting the
Illustrative Code exhibit, so a test covering it directly would assert the
ABSENCE of an assertion. Deviation D1 carries the rationale.

Test files (all committed RED — see *Red confirmation* below):

- `internal/cli/graph_probe_0021_test.go` — the shared probe layer
- `internal/cli/graph_surface_0021_test.go` — `0021:C1`, the verb surface
- `internal/cli/graph_document_0021_test.go` — `0021:C2`, the JSON document
- `internal/cli/graph_dot_0021_test.go` — `0021:C2`'s DOT arm, `0021:S6`
- `internal/cli/graph_determinism_0021_test.go` — `0021:C3`, `0021:S1`/`S2`
- `internal/cli/graph_neutrality_0021_test.go` — `0021:C4`, the ceiling refusal
- `internal/cli/graph_modes_0021_test.go` — `0021:C5`, the four cells
- `internal/cli/graph_stability_0021_test.go` — the tiers and their seams
- `internal/cli/graph_docs_0021_test.go` — `0021:IP` Phase 4's surfaces
- `internal/cli/graph_mvv_0021_test.go` — `0021:MVV`, the acceptance floor
- `internal/graphlint/graph_export_0021_test.go` — `0021:C4`/A8's traversal surface

| REQ | Test |
| --- | --- |
| REQ-1 | TestReq1And15_GraphIsRegisteredAtRootBesideLint |
| REQ-2 | TestReq95_TheCodesAreMirroredButTheHelpTextIsTheVerbsOwn |
| REQ-3 | TestReq2And3And89_BothSelectionFlagsIsMutuallyExclusive |
| REQ-4 | TestReq4And90_FlowAloneIsFlagInvalidValue |
| REQ-5 | TestReq5And89_NeitherSelectionFlagIsFlagRequired |
| REQ-6 | TestReq6And91_UnreadableModelFileIsModelUnreadable |
| REQ-7 | TestReq7And91_LoadFailureIsModelInvalidWithPerCategoryFindings |
| REQ-8 | TestReq8And92_UnknownEmitValueIsFlagInvalidValueNamingEmit; TestReq8_BothDeclaredEmitFormatsAreAccepted; TestReq8_TheEmitDefaultIsJSONAndIsApplied |
| REQ-9 | TestReq9_TheEmitVocabularyIsAppendOnlyAndAssertedByMembership |
| REQ-10 | TestReq10_TheEmitVocabularyHasAnExportedEnumerationSeam |
| REQ-11 | TestReq11_EmitValidityIsCheckedBeforeAnyFileIO; TestReq11_SelectionArmsAreCheckedBeforeEmitValidity |
| REQ-12 | TestReq12_EmitIsTheOnlyValueCheckedFlagOnTheVerb |
| REQ-13 | TestReq13_RunEStartsWithValidateMode; TestReq64_TheSuccessPayloadRidesTheTextLinerSeam |
| REQ-14 | TestReq14_TheVerbUsesOnlyTheExistingZeroTwoExitMapping |
| REQ-15 | TestReq1And15_GraphIsRegisteredAtRootBesideLint; TestReq15_TheRejectedVerbSpellingsAreNotRegistered |
| REQ-16 | TestReq16_TheDocumentCarriesALeadingSchemaField |
| REQ-17 | TestReq17And34And35_AConsumerToleratesUnrecognizedFields |
| REQ-19 | TestReq16_TheDocumentCarriesALeadingSchemaField; TestReq19And46_ModelIsTheAuthoredIDNeverThePathArgument; TestReq19_ClassCarriesTheDeclaredModelClass |
| REQ-20 | TestReq20And28_TagDomainIsPresentExactlyWhenAuthored |
| REQ-21 | TestReq21And32_InitialAndTerminalTravelAsDeclared |
| REQ-22 | TestReq22And23And108_RowsCarryTheDumpFieldVocabulary; TestReq22_AtomsCarryTheirFieldsInTheCanonicalAtomOrder; TestReq22_RowsAreInTheCanonicalRowOrder |
| REQ-23 | TestReq22And23And108_RowsCarryTheDumpFieldVocabulary |
| REQ-24 | TestReq24_GroupsCarryContextAndRules |
| REQ-25 | TestReq25And107_ReachNodeValuesIsATagKeyedObjectOfSortedArrays; TestReq25_ReachEdgesCarryFromToAndRule |
| REQ-26 | TestReq26And107_ReachNodesAndEdgesAreSorted |
| REQ-27 | TestReq27_ANonFiniteOwnedTagCarriesTheOpaqueValueVerbatim |
| REQ-28 | TestReq20And28_TagDomainIsPresentExactlyWhenAuthored; TestReq28And94_NoDeclaredCollectionEverRendersAsNull |
| REQ-29 | TestReq29_ReachCarriesTheRequiredAbstractionMarker |
| REQ-30 | TestReq30_TheSchemaDocsStateTheSoundnessRule |
| REQ-31 | TestReq31_TheDocumentCarriesNoVerdictFindingOrTerminalMarking |
| REQ-32 | TestReq21And32_InitialAndTerminalTravelAsDeclared |
| REQ-33 | TestReq33And65And84_SetMembersSurviveDecodeAsExactArrays |
| REQ-34 | TestReq17And34And35_AConsumerToleratesUnrecognizedFields |
| REQ-35 | TestReq17And34And35_AConsumerToleratesUnrecognizedFields |
| REQ-36 | TestReq36_TheDocumentsOwnDecodeIsTheByValueAssertion |
| REQ-37 | TestReq37_TheEnvelopeVersionIsNeverProjectedIntoTheDocument |
| REQ-38 | TestReq28And94_NoDeclaredCollectionEverRendersAsNull |
| REQ-39 | TestReq39And40And70_TheDOTNodeAndEdgeSetEqualsTheReachBlock; TestReq39_TheDOTDocumentCarriesTheAbstractionMarker; TestReq39_TheDOTDocumentMarksTheInitialNode |
| REQ-40 | TestReq39And40And70_TheDOTNodeAndEdgeSetEqualsTheReachBlock |
| REQ-41 | TestReq41And79_TheDOTArmSurvivesHostileContentAndStillMatches |
| REQ-42 | TestReq42And66_EmissionIsByteIdenticalAcrossInvocations |
| REQ-43 | TestReq43And73_EmissionIsStableUnderProvokedMapOrder |
| REQ-44 | TestReq44And45_TheBytePromiseIsScopedToOneModelAndOneBuild |
| REQ-45 | TestReq44And45_TheBytePromiseIsScopedToOneModelAndOneBuild |
| REQ-46 | TestReq19And46_ModelIsTheAuthoredIDNeverThePathArgument |
| REQ-47 | TestReq47_ExportedIdentitiesNeverCollide |
| REQ-48 | TestReq48And93_TheExportedRelationIsTheSameTraversalLintReaches; TestReq48And93_ThePackageDeclaresOneTraversalNotTwo |
| REQ-49 | TestReq49And50And76_LintIsByteIdenticalWithTheExportCodePresent; TestReq49And104_TheExportEmitsNoFindingsEvenForARefusedModel |
| REQ-50 | TestReq49And50And76_LintIsByteIdenticalWithTheExportCodePresent |
| REQ-51 | TestReq51And77And87_AModelLintRefusesStillExportsAnAssertedDocument |
| REQ-52 | TestReq52And81And88_AnIncompleteTraversalRefusesTooLarge |
| REQ-53 | TestReq53_AnExportedFunctionCarriesNodesEdgesAndCompleteness |
| REQ-54 | TestReq54And55_TheTwoPublishedBoundsAreDistinct |
| REQ-55 | TestReq55_TheGuardProductBoundDoesNotReachTheExportRefusalPath; TestReq54And55_TheTwoPublishedBoundsAreDistinct |
| REQ-56 | TestReq56And100_ReachKeepsItsPresentSignature |
| REQ-57 | TestReq57And61And64_TextModeStdoutIsTheBareDocument |
| REQ-58 | TestReq58_JSONModeEmbedsTheDocumentUnderData |
| REQ-59 | TestReq59And71And78_TheTwoModesAgreeOnOneExportValue |
| REQ-60 | TestReq60And63And96_AllFourCellsAreDefinedAndIndependent |
| REQ-61 | TestReq57And61And64_TextModeStdoutIsTheBareDocument |
| REQ-62 | TestReq62_TheVerbOffersNoCallerControlledProjection |
| REQ-63 | TestReq60And63And96_AllFourCellsAreDefinedAndIndependent |
| REQ-64 | TestReq57And61And64_TextModeStdoutIsTheBareDocument; TestReq64_TheSuccessPayloadRidesTheTextLinerSeam |
| REQ-65 | TestReq33And65And84_SetMembersSurviveDecodeAsExactArrays |
| REQ-66 | TestReq42And66_EmissionIsByteIdenticalAcrossInvocations |
| REQ-67 | TestReq67_NoLoadAfterExportInverseIsClaimed |
| REQ-68 | TestMVV_LintNormalizedGraphExport/step_1 (mvvStep1FixturesLoad) |
| REQ-69 | TestMVV_LintNormalizedGraphExport/step_2 (mvvStep2DoubleEmit) |
| REQ-70 | TestMVV_LintNormalizedGraphExport/step_3 (mvvStep3DOTMatchesReach); TestReq39And40And70_TheDOTNodeAndEdgeSetEqualsTheReachBlock |
| REQ-71 | TestMVV_LintNormalizedGraphExport/step_4 (mvvStep4ModeAgreement); TestReq59And71And78_TheTwoModesAgreeOnOneExportValue |
| REQ-72 | TestMVV_LintNormalizedGraphExport/step_5 (mvvStep5LintNeutrality) |
| REQ-73 | TestReq43And73_EmissionIsStableUnderProvokedMapOrder |
| REQ-74 | TestReq74And75_GoldensPinTheDocumentNotTheEnvelope |
| REQ-75 | TestReq74And75_GoldensPinTheDocumentNotTheEnvelope |
| REQ-76 | TestReq49And50And76_LintIsByteIdenticalWithTheExportCodePresent |
| REQ-77 | TestReq51And77And87_AModelLintRefusesStillExportsAnAssertedDocument |
| REQ-78 | TestReq59And71And78_TheTwoModesAgreeOnOneExportValue |
| REQ-79 | TestReq41And79_TheDOTArmSurvivesHostileContentAndStillMatches |
| REQ-80 | TestReq80_EachDOTIdentifierUnescapesBackToItsSourceID |
| REQ-81 | TestReq52And81And88_AnIncompleteTraversalRefusesTooLarge |
| REQ-82 | TestReq82And83_NoRefusingArmPutsADocumentOnStdout |
| REQ-83 | TestReq82And83_NoRefusingArmPutsADocumentOnStdout |
| REQ-84 | TestReq33And65And84_SetMembersSurviveDecodeAsExactArrays |
| REQ-85 | TestReq85_NoTestInThisRecordAssertsAnExitCodeAlone |
| REQ-86 | TestReq86_ACleanModelExportsADocumentAtExitZero |
| REQ-87 | TestReq51And77And87_AModelLintRefusesStillExportsAnAssertedDocument |
| REQ-88 | TestReq52And81And88_AnIncompleteTraversalRefusesTooLarge |
| REQ-89 | TestReq2And3And89_BothSelectionFlagsIsMutuallyExclusive; TestReq5And89_NeitherSelectionFlagIsFlagRequired |
| REQ-90 | TestReq4And90_FlowAloneIsFlagInvalidValue |
| REQ-91 | TestReq6And91_UnreadableModelFileIsModelUnreadable; TestReq7And91_LoadFailureIsModelInvalidWithPerCategoryFindings |
| REQ-92 | TestReq8And92_UnknownEmitValueIsFlagInvalidValueNamingEmit |
| REQ-93 | TestReq48And93_TheExportedRelationIsTheSameTraversalLintReaches; TestReq48And93_ThePackageDeclaresOneTraversalNotTwo |
| REQ-94 | TestReq28And94_NoDeclaredCollectionEverRendersAsNull |
| REQ-95 | TestReq95_TheCodesAreMirroredButTheHelpTextIsTheVerbsOwn |
| REQ-96 | TestReq60And63And96_AllFourCellsAreDefinedAndIndependent |
| REQ-97 | TestReq97And98And99_ExactlyOneTerminalRecordUnderJSONMode |
| REQ-98 | TestReq97And98And99_ExactlyOneTerminalRecordUnderJSONMode |
| REQ-99 | TestReq97And98And99_ExactlyOneTerminalRecordUnderJSONMode |
| REQ-100 | TestReq56And100_ReachKeepsItsPresentSignature |
| REQ-101 | TestReq101And102_TheExportSurfacesAreBuiltAndReachable |
| REQ-102 | TestReq101And102_TheExportSurfacesAreBuiltAndReachable |
| REQ-103 | TestReq103_TheOutputContractGainsTheDocumentSection; TestReq103_HelpAllGainsTheVerb; TestReq103_TheGeneratedReferenceCarriesTheVerb |
| REQ-104 | TestReq49And104_TheExportEmitsNoFindingsEvenForARefusedModel |
| REQ-105 | TestReq105_AdvisoriesAreIgnoredAndNeverBecomeDocumentContent |
| REQ-106 | TestReq106_TheDocumentClaimsNoHashOrDigest |
| REQ-107 | TestReq25And107_ReachNodeValuesIsATagKeyedObjectOfSortedArrays; TestReq26And107_ReachNodesAndEdgesAreSorted |
| REQ-108 | TestReq22And23And108_RowsCarryTheDumpFieldVocabulary |
| REQ-109 | TestReq109_AngleBracketsAndAmpersandsReachTheWireUnescaped |
| REQ-110 | TestReq110_LintGainsNoFlagFromThisRecord |

## REQ-MVV runner

```sh
go test ./internal/cli/ -run 'TestMVV_LintNormalizedGraphExport' -v
```

`TestMVV_LintNormalizedGraphExport` (`internal/cli/graph_mvv_0021_test.go`)
is the runnable end-to-end test, one subtest per MVV step. The record
declares Round-Trip / Inverse Invariants, so step 2 asserts `export ∘
export` as BYTE identity and step 4 asserts `jq .data` equals the text
document value-for-value — never a green exit code.

## REQ-MVV output

Run in the launch worktree on the fixup tree (the four triage fixes applied):

```
=== RUN   TestMVV_LintNormalizedGraphExport
=== RUN   TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load
=== RUN   TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load/state_machine
=== RUN   TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load/decision_table
=== RUN   TestMVV_LintNormalizedGraphExport/step_2_—_double_emit_is_byte-identical
=== RUN   TestMVV_LintNormalizedGraphExport/step_3_—_DOT_renders_and_matches
=== RUN   TestMVV_LintNormalizedGraphExport/step_4_—_jq_.data_equals_the_document
=== RUN   TestMVV_LintNormalizedGraphExport/step_5_—_lint_neutrality
--- PASS: TestMVV_LintNormalizedGraphExport (0.12s)
    --- PASS: TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load (0.00s)
        --- PASS: TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load/state_machine (0.00s)
        --- PASS: TestMVV_LintNormalizedGraphExport/step_1_—_fixtures_load/decision_table (0.00s)
    --- PASS: TestMVV_LintNormalizedGraphExport/step_2_—_double_emit_is_byte-identical (0.00s)
    --- PASS: TestMVV_LintNormalizedGraphExport/step_3_—_DOT_renders_and_matches (0.12s)
    --- PASS: TestMVV_LintNormalizedGraphExport/step_4_—_jq_.data_equals_the_document (0.00s)
    --- PASS: TestMVV_LintNormalizedGraphExport/step_5_—_lint_neutrality (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/cli	0.368s
```

All five MVV steps pass end to end. Step 2 asserts `export ∘ export` as BYTE
identity and step 4 asserts `jq .data` equals the text document
value-for-value, so neither step rests on a green exit code.

## Red confirmation

79 tests authored; the suite was run on `worktree-rdr-0021` at
`bd7a026` (baseline green).

- **75 RED** — every one a RUNTIME assertion failure, never a compile
  break. The dominant message is `no root "graph" command is registered;
  root commands = [docs flow help lint version]`, plus `internal/graphlint
  declares no exported edges-and-completeness function` and `no golden at
  …/testdata/…`.
- **4 GREEN, each a legitimate invariant guard** over a surface that ships
  TODAY and whose REQ says it must not change. Each goes RED exactly when a
  wrong Phase 2 lands, so none is tautological:
  - `TestReq110_LintGainsNoFlagFromThisRecord` — red if `--emit` lands on
    `lint` (alternative O2) or lint's envelope shape moves.
  - `TestReq56And100_ReachKeepsItsPresentSignature` — red if `Reach`'s
    signature is widened instead of a sibling added (A8's additive rule).
  - `TestReq48And93_ThePackageDeclaresOneTraversalNotTwo` — red if a second
    traversal is added for the export (C4's structural neutrality).
  - `TestReq54And55_TheTwoPublishedBoundsAreDistinct` — red if
    `NodeCeiling()` and `ProductBound()` ever coincide, which would make
    C4's ceiling refusal unattributable to either bound.

Committing red required the probe layer described in
`internal/cli/graph_probe_0021_test.go`: `.githooks/pre-commit` runs
`gofmt -l .` then `go vet ./...` under `set -e`, so a test naming an
unwritten symbol could not be committed at all. The unwritten Go surfaces
are therefore reached through `go/parser` over package source (an absent
function is a `false` return, not a build break) and the unwritten verb
through `ExecuteAndEmit` with the string `"graph"` (cobra refuses it as an
unknown command at run time). No stub was added to the production tree.

Pre-existing tests are unaffected: with these files removed the suite is
`ok` for both `internal/cli` and `internal/graphlint`.

## Notes

- **REQ-18 is the one uncovered row.** "Exact field spellings are normative
  AS SPELLED HERE — the Illustrative Code is an exhibit that disclaims
  literal assertion and shows only some members, so it binds nothing" has
  no observable of its own: its positive half (the field spellings) is
  REQ-19…REQ-25, and its negative half forbids asserting the exhibit, which
  is discharged by no test asserting it. A vacuous test here would assert
  the absence of an assertion.
- **C2's STABILITY clause binds this suite.** No test asserts the document's
  field-set cardinality, a member's ordinal position, or a tail position;
  every field assertion is by NAME. `TestReq17And34And35` pins the
  complementary property — a consumer reading a subset still decodes, and
  an added member does not break it.
- **REQ-9's oracle is membership, never cardinality**, for the same reason:
  a third `--emit` format added in a minor must leave the suite green.
- **The goldens are minted by the real exporter, not by a spike.** `F2`
  scopes A5's sha256 to the ENCODER property and `F3` scopes A6's arm to the
  node/edge set and marker placement, so neither spike's improvised field
  spellings may be promoted. `assertGolden` therefore FAILS rather than
  auto-minting while `internal/cli/testdata/` is absent; Phase 2 writes each
  file from the exporter's own output.
- **Four tests build and run a real subprocess** (`buildGraphBinary`):
  `GODEBUG=randmapiter=1` is read at process start, so S1's provoked
  map-order oracle cannot be satisfied in-process.
- **`TestReq52And81And88` carries a precondition, not a branch.** A fixture
  that failed to exceed `NodeCeiling()` would leave the refusal unexercised,
  so it `t.Fatalf`s rather than skipping — the shape
  `internal/graphlint`'s own ceiling test uses. The generated fixture
  reaches 4103 nodes against the published ceiling of 4096.
- **`TestReq55`'s fixture discriminates the two bounds**: its guard product
  exceeds `ProductBound()` (2048) while its owned-state lattice is 2 nodes,
  far under `NodeCeiling()`. A verb consulting the wrong bound refuses and
  fails the test.
