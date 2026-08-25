# Coverage — RDR 0006 Graph Lint Authority And Guarantees

Phase 1 artifact, updated at Phase 2 close. Every REQ in
[`req-list.md`](req-list.md) maps to at least one test, and every test in
the RDR 0006 suite cites at least one REQ. Both orphan directions are
checked mechanically below and both are empty.

**Status at Phase 1: RED. Status at Phase 2: GREEN.** The `Status` column
in the map below records each test's PHASE 1 state — the red-before-green
evidence — and is deliberately left as authored. At Phase 2 close every one
of the 104 tests passes, and no test was edited to make it do so.

At Phase 1, 94 of the 104 top-level RDR 0006 tests were red. The 10 that
were green are **prohibition and structural clauses whose correct Phase-1
state is green**: a violation of each would be an *added* symbol, an
*added* import, or an *added* command registration, so they are regression
guards from the moment they are written rather than implementation drivers.
They are listed explicitly under *Green at Phase 1* below so Phase 2 does
not mistake them for coverage of behaviour.

Phase 2 replaced `internal/graphlint/skeleton.go` with the working engine
(`taxonomy.go`, `reach.go`, `engine.go`, `analysis.go`, `groups.go`,
`coverage.go`), extended `clierr` and `respond` for the text surface, and
authored the repository gate. `go test ./...` passes and `golangci-lint
run ./...` reports 0 issues.

No pre-existing test was weakened, skipped, or deleted, at either phase.
`internal/table`, `internal/guard`, and `internal/resolve` remain green, as
do the shipped `internal/cli` version tests.

## Totals

| Quantity | Count |
| --- | --- |
| REQs in `req-list.md` (REQ-1..REQ-129 + `REQ-MVV`) | 130 |
| Top-level tests in the RDR 0006 suite | 104 |
| REQs with no test | 0 |
| Tests citing no REQ | 0 |
| Tests carrying no label | 0 |
| Tests currently red | 0 (94 at Phase 1; all green at Phase 2) |
| Tests green at Phase 1 (prohibition/structural) | 10 |
| Pre-existing tests broken | 0 |
| Full suite at Phase 2 close | `go test ./...` PASS, `golangci-lint run` 0 issues |

## REQ-MVV output

Recorded by Phase 2. `TestMVV_GraphLintAuthorityAndGuarantees` PASSES, all
five sub-scenarios and all eleven illegal-matrix arms:

```
--- PASS: TestMVV_GraphLintAuthorityAndGuarantees (0.01s)
    --- PASS: .../legal_model_passes
    --- PASS: .../illegal_matrix
        --- PASS: .../dangling-edge/missing-root
        --- PASS: .../dead-end
        --- PASS: .../overlap/ordinary-pair
        --- PASS: .../overlap/escape-pair-shared-class
        --- PASS: .../coverage-gap
        --- PASS: .../unprovable/non-finite-dimension
        --- PASS: .../unprovable/withheld-claim
        --- PASS: .../always-present-owned
        --- PASS: .../owned-before-write
        --- PASS: .../terminal-escape
        --- PASS: .../product-too-large
    --- PASS: .../multi-defect_group
    --- PASS: .../legal_advisory_matrix
    --- PASS: .../findings_key_always_emitted
ok  github.com/newcoinc/intrastate/internal/cli
```

### The matrix, as the shipped command actually answers it

Every row below is one `intrastate lint --model <fixture> --as=json`
invocation through `ExecuteAndEmit` — the production Cobra path and the
production output gateway, per ASSUMPTION-7. `aggregate` is the
`CLIError.Code`; `findings` counts the emitted list.

| Fixture | exit | aggregate | findings | codes emitted |
| --- | --- | --- | --- | --- |
| `legalModel` | 0 | — | 0 | (empty list) |
| `mvvNoRoot` | 2 | `graph-lint-failed` | 1 | `graph-dangling-edge` |
| `mvvDeadEnd` | 2 | `graph-lint-failed` | 1 | `graph-dead-end` |
| `mvvOrdinaryOverlap` | 2 | `graph-lint-failed` | 2 | `graph-coverage-gap` `graph-overlap` |
| `mvvEscapeOverlap` | 2 | `graph-lint-failed` | 1 | `graph-overlap` |
| `mvvCoverageGap` | 2 | `graph-lint-failed` | 1 | `graph-coverage-gap` |
| `mvvNonFinite` | 2 | `graph-lint-failed` | 1 | `graph-unprovable-coverage` |
| `mvvWithheld` | 2 | `graph-lint-failed` | 1 | `graph-unprovable-coverage` |
| `mvvAlwaysPresentOwned` | 2 | `graph-lint-failed` | 1 | `graph-always-present-owned` |
| `mvvOwnedBeforeWrite` | 2 | `graph-lint-failed` | 3 | `graph-always-present-owned` `graph-coverage-gap` `graph-owned-before-write` |
| `mvvTerminalEscape` | 2 | `graph-lint-failed` | 1 | `graph-terminal-escape` |
| `mvvProductTooLarge` | 2 | `graph-lint-failed` | 1 | `graph-product-too-large` |
| `mvvMultiDefect` | 2 | `graph-lint-failed` | 6 | `graph-always-present-owned` `graph-coverage-gap` `graph-overlap` `graph-owned-before-write` `graph-unprovable-coverage` |
| `mvvLegalAdvisory` | 0 | — | 3 | `graph-coverage-closed-by-escape` `graph-redundant-row` `graph-unreachable-rule` |

`graph-single-valued-state` is the one blocking class the matrix cannot
drive: RDR 0002's loader refuses every authorable spelling before
normalization (deviations `D3`). The check is implemented as a total
function over its input and the class is asserted present in the taxonomy.

The multi-defect run emits **five distinct codes across six findings from
one invocation**, so a first-failure engine cannot pass it, and exactly one
`graph-overlap` finding names both rows of the one overlapping pair.

### The three envelopes, verbatim

**Clean model — the empty list IS the receipt** (`data.findings` present,
`[]` not `null`, and `data` itself present despite `json:"data,omitempty"`):

```json
{"type":"ok","data":{"findings":[]}}
```

**A withheld claim — one aggregate `CLIError`, `findings` a TOP-LEVEL
sibling of `code`, no `error` wrapper and no `type` discriminator**, with
the reason discriminator and all four atom fields:

```json
{
  "code": "graph-lint-failed",
  "message": "the model carries blocking graph-lint findings",
  "findings": [
    {
      "code": "graph-unprovable-coverage",
      "model": "lintfix",
      "severity": "blocking",
      "message": "coverage claim withheld for row \"can-refuse\": the tag \"opt\" is declared optional, so the atom can refuse `guard_unevaluable` at runtime and lint declines to promise exhaustiveness the runtime does not deliver",
      "rule": "can-refuse",
      "span": "lintfix:can-refuse",
      "reason": "row-can-refuse",
      "dimension": "opt",
      "key": "opt",
      "operator": "eq",
      "literal": "p",
      "block": "all",
      "fingerprint": "opt|all|eq|p;status|match|eq|a;#status=b;"
    }
  ]
}
```

The message says the claim is **withheld** and names no authoring fault:
`row-can-refuse` is an honest withholding, not a model defect (REQ-81).

**The legal advisory matrix — all three advisory codes on SUCCESS**, at
`data.findings`, none of which changes the disposition:

```json
{"type":"ok","data":{"findings":[
  {"code":"graph-coverage-closed-by-escape","model":"lintfix","severity":"info",
   "message":"the coverage of group lintfix/stop status.eq=a is closed by the bare escape row \"bare-rescue\" rather than proved over its declared domains",
   "rule":"bare-rescue","element":"lintfix/stop status.eq=a"},
  {"code":"graph-redundant-row","model":"lintfix","severity":"info",
   "message":"row \"narrow\" accepts a proper subset of row \"wide\"'s assignments, so it can never be the exact-one match",
   "rule":"narrow","span":"lintfix:narrow","element":"wide",
   "fingerprint":"flag|all|eq|true;status|match|eq|a;#status=b;"},
  {"code":"graph-unreachable-rule","model":"lintfix","severity":"info",
   "message":"no reachable owned-state satisfies the selection context of row \"orphaned\", so the row can never be a candidate",
   "rule":"orphaned","span":"lintfix:orphaned",
   "fingerprint":"status|match|eq|orphan;#status=b;"}
]}}
```

### The repository gate (REQ-119..REQ-122), also green

`models/rdr.toml` is authored and homed, and lints clean in its accepted
state — the SC-23 false-positive census is **zero blocking findings**:

```
$ make graph-lint
bin/intrastate lint --model models/rdr.toml --as=json
{"type":"ok","data":{"findings":[]}}
```

Dropping its `[initial]` table — the deliberate defect REQ-120 authors —
fails the same command with the aggregate code, so the gate bites:

```
code=[graph-lint-failed]
```

## Green at Phase 1

These assert that something is **absent** or that a declared structure
holds. Each would go red if the prohibition were violated.

| Test | REQ | Why green is correct now |
| --- | --- | --- |
| `TestReq10And11_FlowAndModelAreMutuallyExclusiveUsageErrors` | REQ-10, REQ-11 | REQ-10/11 fix the flag contract; the skeleton's RunE already refuses the pair at `GroupUserEnv`. |
| `TestReq118_NoRoundTripOrInverseSurfaceIsIntroduced` | REQ-118 | REQ-118 forbids an encode/decode or import/export pair on the engine. |
| `TestReq14_ModelPathIsTheInstantiableForm` | REQ-14 | REQ-14 requires `--flow` refuse until discovery lands. |
| `TestReq2_ResolverConsultsNoLintVerdict` | REQ-2 | REQ-2 forbids a resolver-side lint import. Red would mean one was added. |
| `TestReq45_OptionalityAndSingleValuednessComeFromTheDeclarationModel` | REQ-16, REQ-45 | REQ-45's declaration-reading half asserts agreement with the landed RDR 0003 peer, which is a cross-package regression guard. |
| `TestReq4_RootLintIsTheAuthoritativeSurface` | REQ-4 | REQ-4 requires the engine expose `Run`; the skeleton declares it. |
| `TestReq6_LintIsRegisteredAtRootNotUnderFlow` | REQ-4, REQ-6 | REQ-6 requires root registration and forbids a `flow lint` sibling. |
| `TestReq7_ExactlyOneEntryPointAndOneRequestBuilder` | REQ-7 | REQ-7 caps the exported surface at one entry point and one builder. Red would mean a second was added. |
| `TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields` | REQ-95, REQ-96 | REQ-95/96 require the record live in `clierr` with string-typed atom fields and no subsystem import. |
| `TestReq9_CanonicalSubsystemNameIsLint` | REQ-9 | REQ-9 fixes the `graph-` code namespace and rejects `validate`/`pre-commit` spellings. |

## REQ x test

| REQ | Test | Label | File | Phase 1 |
| --- | --- | --- | --- | --- |
| REQ-1 | `TestReq1_AnyBlockingFindingRefusesAcceptance` | HAPPY PATH | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-2 | `TestReq2_ResolverConsultsNoLintVerdict` | ADVERSARIAL | `internal/graphlint/authority_0006_test.go` | green |
| REQ-3 | `TestReq3_EngineDefinesNoSecondSourceParser` | ADVERSARIAL | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-4 | `TestReq6_LintIsRegisteredAtRootNotUnderFlow` | HAPPY PATH | `internal/cli/lint_0006_test.go` | green |
| REQ-4 | `TestReq4_RootLintIsTheAuthoritativeSurface` | HAPPY PATH | `internal/graphlint/authority_0006_test.go` | green |
| REQ-5 | `TestReq5And99_BothModesRouteThroughTheGatewayWithNoDirectWrites` | ADVERSARIAL | `internal/cli/lint_0006_test.go` | RED |
| REQ-6 | `TestReq6_LintIsRegisteredAtRootNotUnderFlow` | HAPPY PATH | `internal/cli/lint_0006_test.go` | green |
| REQ-7 | `TestReq7_ExactlyOneEntryPointAndOneRequestBuilder` | BOUNDARY | `internal/graphlint/authority_0006_test.go` | green |
| REQ-8 | `TestReq8_EngineReceivesANormalizedGraphValueNotCommandState` | BOUNDARY | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-9 | `TestReq9_CanonicalSubsystemNameIsLint` | HAPPY PATH | `internal/graphlint/authority_0006_test.go` | green |
| REQ-10 | `TestReq10And11_FlowAndModelAreMutuallyExclusiveUsageErrors` | INPUT EDGE | `internal/cli/lint_0006_test.go` | green |
| REQ-11 | `TestReq10And11_FlowAndModelAreMutuallyExclusiveUsageErrors` | INPUT EDGE | `internal/cli/lint_0006_test.go` | green |
| REQ-12 | `TestReq12And13_OneModelPerInvocationAndFindingsNameTheModelID` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-13 | `TestReq12And13_OneModelPerInvocationAndFindingsNameTheModelID` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-13 | `TestReq13_FindingModelFieldIsTheModelIDNeverThePath` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-14 | `TestReq14_ModelPathIsTheInstantiableForm` | INPUT EDGE | `internal/cli/lint_0006_test.go` | green |
| REQ-15 | `TestReq15_MinimumInputContractIsAvailableToTheEngine` | HAPPY PATH | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-16 | `TestReq16_DeclaredPropertiesAreReadNeverInferred` | ADVERSARIAL | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-16 | `TestReq45_OptionalityAndSingleValuednessComeFromTheDeclarationModel` | BOUNDARY | `internal/graphlint/soundness_0006_test.go` | green |
| REQ-17 | `TestReq17_PipelineRunsNormalizeDeriveCheckEmit` | HAPPY PATH | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-18 | `TestReq18_OwnedReadSetIsDerivedFromAtomsNotRequiresOwned` | ADVERSARIAL | `internal/graphlint/authority_0006_test.go` | RED |
| REQ-19 | `TestReq19_GroupingIsRDR0003sAndNotASecondPredicate` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-20 | `TestReq20And21And52_GroupMembershipIsTheAuthoredMatchPatternNotANode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-21 | `TestReq20And21And52_GroupMembershipIsTheAuthoredMatchPatternNotANode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-22 | `TestReq22_ReachabilityFiltersContextsAndNeverChangesMembership` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-23 | `TestReq23And24_OverlapAndCoverageRunOncePerReachableGroup` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-24 | `TestReq23And24_OverlapAndCoverageRunOncePerReachableGroup` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-25 | `TestReq25And26And27_TwoOrdinaryRowsEnabledTogetherAreRejected` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-26 | `TestReq25And26And27_TwoOrdinaryRowsEnabledTogetherAreRejected` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-27 | `TestReq25And26And27_TwoOrdinaryRowsEnabledTogetherAreRejected` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-28 | `TestReq28_EscapeOverlappingAnOrdinaryRowIsNotReported` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-29 | `TestReq29And30_EscapePairsOverlapPerSharedClassOneFindingEach` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-30 | `TestReq29And30_EscapePairsOverlapPerSharedClassOneFindingEach` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-30 | `TestReq30_OneOverlappingPairYieldsOneFindingNamingBothRows` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-31 | `TestReq31_TheSevenMandatoryInvariantClassesAreAllChecked` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-32 | `TestReq32_DanglingEdgeChecksReferencesAndRequiresAnInitialOwnedState` | HAPPY PATH | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-33 | `TestReq33_TerminalsAreTagPredicatesNotStateNames` | DOMAIN EDGE | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-34 | `TestReq34_ReachableNonTerminalNodeWithoutAnOutgoingRowIsADeadEnd` | HAPPY PATH | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-35 | `TestReq35_TerminalSatisfactionIsUniversalOverTheNodesValueSet` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-36 | `TestReq36_DeadEndSplitsTheNodeOnTerminalParticipatingKeys` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-37 | `TestReq37_SplitIsBoundedByTerminalParticipatingKeysOnly` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-38 | `TestReq38And47_CoverageClaimIsDefaultOnAndProvedAsUnionEqualsProduct` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-39 | `TestReq39_EscapeRowsContributeToTheUnionLikeAnyOtherRow` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-40 | `TestReq40_SingleValuedStateIsDecidedPerRowSyntactically` | DOMAIN EDGE | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-41 | `TestReq41_MergedNodeWithTwoValuesIsNotASingleValuedViolation` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-42 | `TestReq42_OwnedSetBeforeMatchRequiresTheKeyHeldAtEveryMatchingNode` | HAPPY PATH | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-43 | `TestReq43_RelianceOnAnImpliedTerminalIsGraphTerminalEscape` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-44 | `TestReq44_AbsentRootIsBlockingNeverAnEmptyReachableSetGreen` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-45 | `TestReq45And46_NonFiniteDimensionTakesUnprovableCoverage` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-45 | `TestReq45_OptionalityAndSingleValuednessComeFromTheDeclarationModel` | BOUNDARY | `internal/graphlint/soundness_0006_test.go` | green |
| REQ-46 | `TestReq45And46_NonFiniteDimensionTakesUnprovableCoverage` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-47 | `TestReq38And47_CoverageClaimIsDefaultOnAndProvedAsUnionEqualsProduct` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-48 | `TestReq48And49_PresenceDimensionIsInTheProductForOptionalKeysOnly` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-49 | `TestReq48And49_PresenceDimensionIsInTheProductForOptionalKeysOnly` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-50 | `TestReq50And53_CoverageReadsAuthoredRowsAndNeverAReachabilityNode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-51 | `TestReq51_MatchKeysAreNotProductDimensions` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-52 | `TestReq20And21And52_GroupMembershipIsTheAuthoredMatchPatternNotANode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-53 | `TestReq50And53_CoverageReadsAuthoredRowsAndNeverAReachabilityNode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-54 | `TestReq54And56_WithheldClaimNamesTheRowAndTheRefusingAtom` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-55 | `TestReq55And76_WithheldClaimIsBlockingAndNotInTheAdvisoryTier` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-56 | `TestReq54And56_WithheldClaimNamesTheRowAndTheRefusingAtom` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-57 | `TestReq57And59_DeclaredFiniteProductOverTheBoundIsProductTooLarge` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-58 | `TestReq58_ValueAtomOverANonSingleValuedTagIsUnprovable` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-59 | `TestReq59And60And126_BothBoundsAppearInTheCommandsHelpOutput` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-59 | `TestReq57And59_DeclaredFiniteProductOverTheBoundIsProductTooLarge` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-60 | `TestReq59And60And126_BothBoundsAppearInTheCommandsHelpOutput` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-60 | `TestReq60_ReachableNodeSetOverTheCeilingIsProductTooLargeOnTheTraversal` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-61 | `TestReq61And62_EscapeClosesOnlyTheClassesItDeclares` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-62 | `TestReq61And62_EscapeClosesOnlyTheClassesItDeclares` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-63 | `TestReq63_AmbiguousMatchArmIsVacuouslyClosedForAnOverlapFreeGroup` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-64 | `TestReq64_OwnedStateUnavailableIsNeitherRescuableNorMinted` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-65 | `TestReq65And66_BareEscapeClosureIsReportedAndDoesNotFailTheRun` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-66 | `TestReq65And66_BareEscapeClosureIsReportedAndDoesNotFailTheRun` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-67 | `TestReq67_AlwaysPresentOwnedKeyMustBeHeldAtEveryReachableNode` | HAPPY PATH | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-68 | `TestReq68_AlwaysPresentCheckDoesNotExtendToObservedKeys` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-69 | `TestReq69And70_RootIsIncludedAndTheFindingNamesTheInitialDeclaration` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-70 | `TestReq69And70_RootIsIncludedAndTheFindingNamesTheInitialDeclaration` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-71 | `TestReq71_EveryBlockingFindingCarriesTheRequiredIdentityFields` | HAPPY PATH | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-71 | `TestReqSeverityVocabularyIsExactlyBlockingAndInfo` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-72 | `TestReq72And82_AtomAttributedAndEscapeScopedFindingsCarryTheirFields` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-73 | `TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-74 | `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-75 | `TestReq75_AdvisoryFindingsDoNotChangeTheSuccessDisposition` | HAPPY PATH | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-76 | `TestReq55And76_WithheldClaimIsBlockingAndNotInTheAdvisoryTier` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-77 | `TestReq77_RedundantRowIsAProperSubsetOfASiblingsAssignments` | DOMAIN EDGE | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-78 | `TestReq78_UnreachableRuleCoversNoReachableNodeAndTheDeadRuleCase` | DOMAIN EDGE | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-79 | `TestReq79_ExistsOverAnAlwaysPresentKeyIsAVacuousAtomNotARejection` | DOMAIN EDGE | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-80 | `TestReq45And46_NonFiniteDimensionTakesUnprovableCoverage` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-80 | `TestReq54And56_WithheldClaimNamesTheRowAndTheRefusingAtom` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-80 | `TestReq58_ValueAtomOverANonSingleValuedTagIsUnprovable` | DOMAIN EDGE | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-80 | `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-81 | `TestReq54And56_WithheldClaimNamesTheRowAndTheRefusingAtom` | HAPPY PATH | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-81 | `TestReq81_RowCanRefuseMessageReadsAsAWithholdingNotAnError` | DOMAIN EDGE | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-82 | `TestReq72And82_AtomAttributedAndEscapeScopedFindingsCarryTheirFields` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-83 | `TestReq83And84And85_EveryDecidableDefectIsEmittedInOnePass` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-84 | `TestReq83And84And85_EveryDecidableDefectIsEmittedInOnePass` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-85 | `TestReq83And84And85_EveryDecidableDefectIsEmittedInOnePass` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-85 | `TestReq85_EmittedSetIsIndependentOfAuthoredRowOrder` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-86 | `TestReq86And114_FindingsAreOrderedByTheFindingIdentityTuple` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-87 | `TestReq87_FingerprintIsACanonicalSortableSerializationNeverAHash` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-88 | `TestReq88_RuleIDsSortBeforeGraphElementIDs` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-89 | `TestReq85_EmittedSetIsIndependentOfAuthoredRowOrder` | ADVERSARIAL | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-90 | `TestReq90_FailureIsOneAggregateCLIErrorAtGroupUserEnv` | HAPPY PATH | `internal/cli/lint_0006_test.go` | RED |
| REQ-91 | `TestReq91And92_FindingsAreATopLevelSiblingOfCodeOnTheFailureLine` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-92 | `TestReq91And92_FindingsAreATopLevelSiblingOfCodeOnTheFailureLine` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-93 | `TestReq93_SuccessCarriesNonBlockingFindingsAtDataFindings` | HAPPY PATH | `internal/cli/lint_0006_test.go` | RED |
| REQ-94 | `TestReq94And128_TheFindingsKeyIsAlwaysEmittedIncludingTheEmptyList` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-95 | `TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields` | BOUNDARY | `internal/cli/lint_0006_test.go` | green |
| REQ-96 | `TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields` | BOUNDARY | `internal/cli/lint_0006_test.go` | green |
| REQ-97 | `TestReq97And98_TextModeEnumeratesEveryFindingsCodeAndMessage` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-98 | `TestReq97And98_TextModeEnumeratesEveryFindingsCodeAndMessage` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-99 | `TestReq5And99_BothModesRouteThroughTheGatewayWithNoDirectWrites` | ADVERSARIAL | `internal/cli/lint_0006_test.go` | RED |
| REQ-100 | `TestReq100_NonConformingModelIsRefusedUpstreamWithNoneOfThisRDRsCodes` | INPUT EDGE | `internal/cli/lint_0006_test.go` | RED |
| REQ-101 | `TestReq101_ReachabilityIsTheOwnedStateGraphRootedAtInitial` | HAPPY PATH | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-102 | `TestReq102_EscapeRowsAreSelfLoopEdges` | DOMAIN EDGE | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-103 | `TestReq103_EscapeSelfLoopIsNotProgressForTheDeadEndTest` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-104 | `TestReq104_ObservedMatchAtomsAndGuardAtomsAreNotPruned` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-105 | `TestReq105_ASelectionContextIsReachableWhenSomeNodeSatisfiesIt` | HAPPY PATH | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-106 | `TestReq106_TraversalIsTotalAndExecutesNoAccessor` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-107 | `TestReq107_CanRefuseIsDecidedOverOptionalityAndNeverConsultsReachability` | ADVERSARIAL | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-108 | `TestReq108_TraversalMergesConvergingEdgesIntoOneNode` | BOUNDARY | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-109 | `TestReq109_FixpointTerminatesOnACyclicGraph` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-110 | `TestReq110And112_OwnedSetBeforeMatchReadsMergedFixpointNodes` | DOMAIN EDGE | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-111 | `TestReq111_UniversalChecksDoNotReadWidenedNodesDirectly` | ADVERSARIAL | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-112 | `TestReq110And112_OwnedSetBeforeMatchReadsMergedFixpointNodes` | DOMAIN EDGE | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-113 | `TestReq113_ARowPreservesATagItNeitherWritesNorClears` | DOMAIN EDGE | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-114 | `TestReq86And114_FindingsAreOrderedByTheFindingIdentityTuple` | BOUNDARY | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-115 | `TestReq115_WithholdingDominatesBareEscapeClosure` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-116 | `TestReq116And117_NoSuppressionOrWaiverChannelExists` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-117 | `TestReq116And117_NoSuppressionOrWaiverChannelExists` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-117 | `TestReq117_AcceptedFalsePositivesAreCuredByTheModelNotByAWeakerLint` | DOMAIN EDGE | `internal/graphlint/soundness_0006_test.go` | RED |
| REQ-118 | `TestReq118_NoRoundTripOrInverseSurfaceIsIntroduced` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | green |
| REQ-119 | `TestReq119_CICarriesAGraphLintJobOverTheCheckedInModel` | BOUNDARY | `internal/cli/lint_gate_0006_test.go` | RED |
| REQ-119 | `TestReq119_MakeCheckIsWiredToTheSameCommandWithABuildEdge` | BOUNDARY | `internal/cli/lint_gate_0006_test.go` | RED |
| REQ-120 | `TestReq120_ADeliberatelyIllegalEditToTheCheckedInModelFailsTheGate` | ADVERSARIAL | `internal/cli/lint_gate_0006_test.go` | RED |
| REQ-121 | `TestReq121_ATransitionModelIsCheckedIntoTheRepo` | HAPPY PATH | `internal/cli/lint_gate_0006_test.go` | RED |
| REQ-122 | `TestReq122_FalsePositiveCensusOverTheCheckedInModelIsZero` | HAPPY PATH | `internal/cli/lint_gate_0006_test.go` | RED |
| REQ-123 | `TestReq48And49_PresenceDimensionIsInTheProductForOptionalKeysOnly` | BOUNDARY | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-123 | `TestReq68_AlwaysPresentCheckDoesNotExtendToObservedKeys` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-124 | `TestReq20And21And52_GroupMembershipIsTheAuthoredMatchPatternNotANode` | ADVERSARIAL | `internal/graphlint/coverage_0006_test.go` | RED |
| REQ-124 | `TestReq44_AbsentRootIsBlockingNeverAnEmptyReachableSetGreen` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-125 | `TestReq36_DeadEndSplitsTheNodeOnTerminalParticipatingKeys` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-126 | `TestReq59And60And126_BothBoundsAppearInTheCommandsHelpOutput` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-127 | `TestReq127_FindingWithoutARuleIDCarriesASpanOrElementID` | INPUT EDGE | `internal/graphlint/findings_0006_test.go` | RED |
| REQ-128 | `TestReq94And128_TheFindingsKeyIsAlwaysEmittedIncludingTheEmptyList` | BOUNDARY | `internal/cli/lint_0006_test.go` | RED |
| REQ-129 | `TestReq129_TerminalOnANonOwnedTagIsABlockingDanglingFinding` | ADVERSARIAL | `internal/graphlint/invariants_0006_test.go` | RED |
| REQ-MVV | `TestMVV_GraphLintAuthorityAndGuarantees` | HAPPY PATH | `internal/cli/lint_mvv_0006_test.go` | RED |

## Orphan check

Both directions are computed from the test sources and `req-list.md`, not
maintained by hand.

- **REQs with no test: none.** All 130 REQs (REQ-1..REQ-129 plus
  `REQ-MVV`) are cited by at least one test header.
- **Tests citing no REQ: none.** Every top-level `Test*` function in
  `*_0006_test.go` opens with a `// REQ-N: "<quote>"` header.
- **Tests carrying no label: none.** Every test carries exactly one of
  `HAPPY PATH`, `INPUT EDGE`, `BOUNDARY`, `ADVERSARIAL`, `DOMAIN EDGE`.

Several REQs share a test where the clauses are two halves of one
assertion and splitting them would make each half vacuous — REQ-20/21/52
(group membership and the pooling trap are the same fixture's positive and
negative control), REQ-25/26/27, REQ-29/30, REQ-69/70, REQ-91/92,
REQ-94/128, REQ-97/98, REQ-110/112. Each such test cites every REQ it
discharges in its header, so the mapping above stays total in both
directions.

## Files

| File | Area |
| --- | --- |
| `internal/graphlint/fixtures_0006_test.go` | Shared fixture harness. Builds models through the real `internal/table` loader; nothing mocks the unit under test. |
| `internal/graphlint/authority_0006_test.go` | Authority, placement, engine boundary, minimum input contract, the atom-derived read-set. |
| `internal/graphlint/invariants_0006_test.go` | The seven mandatory invariant classes, the always-present owned check, the reachability relation and its join rule. |
| `internal/graphlint/coverage_0006_test.go` | Row grouping, the two overlap populations, exhaustiveness and withholding, the published bounds, escape class scoping, precedence. |
| `internal/graphlint/findings_0006_test.go` | Finding shape, the code taxonomy and severity tiers, emission completeness, deterministic identity ordering. |
| `internal/graphlint/soundness_0006_test.go` | Existential vs universal checks over merged nodes, and the accepted false positives SC-10 buys. |
| `internal/cli/lint_0006_test.go` | The CLI surface and output envelope, driven through `ExecuteAndEmit`. |
| `internal/cli/lint_fixtures_0006_test.go` | The MVV fixture corpus. All 16 models verified to load clean through the RDR 0002 loader. |
| `models/rdr.toml` | The checked-in transition model the repository gate lints (A8). Authored at Phase 2. |
| `internal/cli/lint_mvv_0006_test.go` | `REQ-MVV`, the runnable end-to-end validation. |
| `internal/cli/lint_gate_0006_test.go` | The repository acceptance gate: the CI job, `make check`, and the false-positive census. |

## Implementation surface (Phase 2)

The Phase 1 skeleton was one file of zero-value signatures. Phase 2
replaced it with the working engine, split by concern; every symbol keeps
the name and shape the tests were written against.

| Symbol | File |
| --- | --- |
| The 10 blocking + 4 advisory codes, 3 reasons, 2 severities, `IsBlocking`, `ProductBound`, `NodeCeiling` | `internal/graphlint/taxonomy.go` |
| `graphlint.Reach`, `graphlint.Node`, the fixpoint successor relation and its join | `internal/graphlint/reach.go` |
| `graphlint.Run`, `NewRequest`, `Request`, `Report`, `Fingerprint`, the finding-identity ordering | `internal/graphlint/engine.go` |
| Invariants 1, 2, 5, 7, the always-present owned check, the node ceiling, the unreachable-rule advisory | `internal/graphlint/analysis.go` |
| Invariants 3 and 6, the redundant-row and vacuous-atom advisories | `internal/graphlint/groups.go` |
| Invariant 4: the exhaustiveness claim, the three withholding triggers, the per-class coverage arms | `internal/graphlint/coverage.go` |
| `clierr.Finding`, `CLIError.Findings` (non-`omitempty`), `EmitText`/`EmitFindingsText` | `internal/cli/clierr/clierr.go` |
| `respond.FindingCarrier` and the `OK` text branch (deviation `D5`) | `internal/cli/respond/respond.go` |
| root `lint` verb, `lintPayload` | `internal/cli/lint.go` |
| The checked-in transition model the gate lints | `models/rdr.toml` |
| The `graph-lint` CI job / `make graph-lint` + `check`'s `build` edge | `.github/workflows/ci.yml`, `Makefile` |

`ProductBound()` returns RDR 0003's published bound (`guard.Bound()`, 2048)
rather than a second constant: the proof representation lint declines over
IS `internal/guard`'s enumerating one, so a bound of its own would name a
threshold nothing applies and the two would drift. `NodeCeiling()` is this
RDR's own — 4096 merged nodes — since no peer bounds the traversal.

## Notes for Phase 2 (dispositions recorded inline)

1. **`graph-single-valued-state` has no authorable input surface today.**
   *Disposition: implemented as a total function
   (`analysis.go::checkSingleValuedState`) and left undriven; `D3` stays
   open on RDR 0002's terms.*
   RDR 0002's loader refuses every TOML spelling of a single-valued tag
   written two values before normalization: `kind <k> holds one value, not
   a member sequence` for `enum`/`int`/`bool`/`scalar`, and `kind set
   admits no single_valued marker` for `set`. `TestReq40_...` therefore
   asserts the class is a declared, blocking member of the taxonomy and
   that the per-row reading does not fire on the legal shapes, rather than
   driving the defect through a fixture. Recorded as a deviation.

2. **Invariant 1's syntactic arms are discharged upstream.** An undeclared
   tag, an out-of-alphabet outcome, and an unknown gate accessor are all
   RDR 0002 load refusals. The arms of `graph-dangling-edge` that reach
   lint are the missing `[initial]` declaration (REQ-44) and the terminal
   context reading a non-owned tag (REQ-129), both of which load clean.
   *Disposition: both implemented; the second takes `graph-dangling-edge`
   per ASSUMPTION-1, naming the terminal context index and the non-owned
   key.*

3. **REQ-119..REQ-122 depend on A8.** No transition model is checked in
   today. Per ASSUMPTION-9 it is authored during this implementation and
   homed at `models/rdr.toml`, which is the path the gate tests assert.
   Phase 2 authors it, adds the `graph-lint` CI job, and wires `make
   check`'s missing `build` edge. *Disposition: done. The model lints
   clean — SC-23's false-positive census is zero — and dropping its
   `[initial]` table fails the same command. `D4` is discharged.*
