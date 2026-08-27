# Coverage — RDR 0009 Ownership of escape-row shape conformance

Phase 1 (test-author) artifact. Every executable REQ in
`req-list.md` maps to at least one test whose header quotes the clause
verbatim and carries exactly one label.

**Totals**: 109 REQs. 107 executable, 2 DEFERRED (REQ-90, REQ-91) by the
record's own Testing-Strategy scope note. 107 covered, 0 orphaned in either
direction.

## Test files

| Path | Package | Scope |
| --- | --- | --- |
| `internal/resolve/escape_shape_fixtures_0009_test.go` | `resolve_test` | fixture builders + structural error readers (no REQ of its own) |
| `internal/resolve/escape_shape_0009_test.go` | `resolve_test` | sections A, C, D, E, F |
| `internal/resolve/escape_shape_mvv_0009_test.go` | `resolve_test` | REQ-MVV + section I scenarios |
| `internal/resolve/escape_shape_surface_0009_test.go` | `resolve_test` | sections G, J, L |
| `internal/table/escape_shape_0009_test.go` | `table_test` | sections B, K |
| `internal/cli/escape_shape_0009_test.go` | `cli` (internal) | section H (C7) |

## Status legend

- **RED** — fails at HEAD; the behaviour does not exist yet.
- **SAT** — `already-satisfied-at-HEAD`; the test pins shipped behaviour
  green so a future change cannot relax it.
- **DEFERRED** — the record places the scenario outside this RDR's Done
  criteria; no test written by design.

Where a REQ is RED only because the package does not compile at HEAD, the
red gate was re-run against a minimal declaration-only stub
(`Table.CheckValid` returning nil, `EscapeShapeBreachError{Ref,Count}`,
`ErrEscapeShapeBreach`, `kernelResolveFailure` delegating to the shipped
generic branch) so each test was exercised individually rather than lost
behind a build failure. The status column reflects that run.

---

## A. The producer obligation

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-1 | `TestReq1_AnEscapeRowMustHaveAnEmptyWritesSlice` | HAPPY PATH | RED |
| REQ-2 | `TestReq2_ANormalizedClearOnAnEscapeRowIsTheSameBreach` | DOMAIN EDGE | RED |
| REQ-3 | `TestReq3_ThePredicateDoesNotExtendToNextTags` | BOUNDARY | RED |
| REQ-4 | `TestReq4_TheBreachPredicateIsTheLengthConjunction` | BOUNDARY | RED |
| REQ-5 | `TestReq5_NilAndEmptySlicesAreIndistinguishableToThePredicate` | INPUT EDGE | RED |
| REQ-6 | `TestReq6_RowKeepsItsSingleShapeWithNoKindField` | ADVERSARIAL | SAT (shape claim; `Row` already carries no kind field) |
| REQ-7 | `TestReq7_TheEscapeDiscriminatorIsTheReusedNonEmptyEscapeSignal` | BOUNDARY | RED |
| REQ-8 | `TestReq8_PlanOfStaysUnconditionalAndAddsNoStrippingLogic` | DOMAIN EDGE | RED |
| REQ-9 | `TestReq9_DataFlowIsUnchangedForEveryConformingTable` | HAPPY PATH | RED |
| REQ-10 | `TestReq10_TheKernelNeverSanitizesAwayABreach` | ADVERSARIAL | RED |

## B. The authored-path enforcer

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-11 | `TestReq11_AnAuthoredEscapeRuleCarryingWritesOrClearsIsRejectedAtLoad` | HAPPY PATH | SAT |
| REQ-12 | `TestReq12_AnEmptyWriteBlockOnAnEscapeRuleIsRejectedToo` | INPUT EDGE | SAT |
| REQ-13 | `TestReq13_TheAuthoredLayerIsStricterThanTheKernelPredicate` | ADVERSARIAL | SAT |
| REQ-14 | `TestReq14_NormalizedEscapeRowsRenderWriteFree` | HAPPY PATH | SAT |
| REQ-15 | `TestReq15_TheAuthoredLayerNeverAcceptsWhatTheKernelRejects` | ADVERSARIAL | SAT |

Section B is `already-satisfied-at-HEAD` in its entirety.
`internal/table/normalize.go::normalizeRule` already rejects an escape rule
carrying a write block, a clear list, or a gate list, keyed on key PRESENCE
(`rule.Write != nil`), under `CatMalformedEscapeDeclaration`. The tests pin
that against relaxation — in particular REQ-13, which fails the moment the
loader is relaxed to a length check and the deliberate strictness gap closes.
The kernel half of each (`CheckValid`) is RED at HEAD, so these tests still
do not compile until Phase 2 lands.

## C. The kernel entry precondition

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-16 | `TestReq16_ResolveEnforcesTheSameObligationAtEntry` | HAPPY PATH | RED |
| REQ-17 | `TestReq17_ResolveKeepsItsOneParameterSignature` | BOUNDARY | SAT (signature already one-parameter) |
| REQ-18 | `TestReq18_ResolveReturnsATypedErrorNamingTheRowAndNoDisposition` | HAPPY PATH | RED |
| REQ-19 | `TestReq19_ADormantBreachingRowStillErrors` | ADVERSARIAL | RED |
| REQ-20 | `TestReq20_TheBreachIsTheErrorPathAndNotAModeledRefusal` | ADVERSARIAL | RED |
| REQ-21 | `TestReq21_TheBreachPrecedesUnmodeledOutcome` | ADVERSARIAL | RED |
| REQ-22 | `TestReq22_TheBreachErrorsWhateverTheAssembledViewWouldHold` | DOMAIN EDGE | RED |
| REQ-23 | `TestReq23_ABreachReturnsTheZeroResultAndANonNilError` | BOUNDARY | RED |
| REQ-24 | `TestReq24_TheZeroResultReportsRefusedFalseOnABreach` | ADVERSARIAL | RED |

## D. The typed error surface

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-25 | `TestReq25_TheBreachErrorIsTypedAndWrapsTheSentinel` | HAPPY PATH | RED |
| REQ-26 | `TestReq26_TheTypeAndSentinelAreExported` | BOUNDARY | RED (compile-level: names absent at HEAD) |
| REQ-27 | `TestReq27_RowIdentityIsRecoverableWithoutParsingProse` | ADVERSARIAL | RED |
| REQ-28 | `TestReq28_TheSentinelIsSpelledErrEscapeShapeBreach` | BOUNDARY | RED |
| REQ-29 | `TestReq29_TheTypedErrorCarriesExactlyOneRowRefFieldAndACount` | BOUNDARY | RED (compile-level) |
| REQ-30 | `TestReq30_EveryPerRowElementUnwrapsToTheSentinel` | BOUNDARY | RED |
| REQ-31 | `TestReq31_ThePredicateIsSpelledTableCheckValid` | BOUNDARY | RED (compile-level) |
| REQ-32 | `TestReq32_TheFourPinnedSpellingsAllExist` | ADVERSARIAL | RED (compile-level) |
| REQ-33 | `TestReq33_TheSingleBreachReturnIsTheJoinWrapperNotTheBareError` | ADVERSARIAL | RED |
| REQ-34 | `TestReq34_TheAggregateIsTheUniformReturnShapeAtEveryCardinality` | BOUNDARY | RED |
| REQ-35 | `TestReq35_TheAggregateIsExactlyOneLevelDeep` | ADVERSARIAL | RED |
| REQ-36 | `TestReq36_ResolveReturnsCheckValidsErrorVerbatim` | ADVERSARIAL | RED |
| REQ-37 | `TestReq37_ThePinnedSpellingsCarryThePinnedTypes` | BOUNDARY | RED (compile-level) |
| REQ-38 | `TestReq38_ASingleErrorsAsReportsOnlyOneOfSeveralBreaches` | DOMAIN EDGE | RED |

## E. Multi-breach reporting, ordering, and collapse

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-39 | `TestReq39_EveryBreachingRowIsReportedInOnePass` | HAPPY PATH | RED |
| REQ-40 | `TestReq40_TheReportIsOrderedByIdentityNotByRowPosition` | ADVERSARIAL | RED |
| REQ-41 | `TestReq41_RowPositionNeverBreaksAnIdentityTie` | ADVERSARIAL | RED |
| REQ-42 | `TestReq42_EqualIdentitiesCollapseToOneReportedError` | HAPPY PATH | RED |
| REQ-43 | `TestReq43_RowsWithNoSourceIdentityCollapseAndCarryTheCount` | INPUT EDGE | RED |
| REQ-44 | `TestReq44_TheCountIsCarriedOnTheStructNotOnlyInProse` | BOUNDARY | RED |
| REQ-45 | `TestReq45_CountIsPerIdentityAndCountsPreCollapseRows` | BOUNDARY | RED |
| REQ-46 | `TestReq46_TheReportIsAggregateSortedAndCollapsed` | HAPPY PATH | RED |

## F. The exported predicate

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-47 | `TestReq47_OnePredicateTwoCallSitesNeverDrift` | HAPPY PATH | RED |
| REQ-48 | `TestReq48_CheckValidIsAMethodOnTableTakingNoArguments` | BOUNDARY | RED (compile-level) |
| REQ-49 | `TestReq49_ThePredicateReturnsErrorAndNotABoolPredicate` | BOUNDARY | RED |
| REQ-50 | `TestReq50_CheckValidsDocFollowsTheStdlibForm` | BOUNDARY | RED |
| REQ-51 | `TestReq51_CheckValidsDocNamesTheOnePropertyAndItsLimits` | ADVERSARIAL | RED |
| REQ-52 | `TestReq52_ARowlessTableConformsVacuously` | INPUT EDGE | RED |
| REQ-53 | `TestReq53_TheTwoEnforcementPointsCannotDesynchronize` | ADVERSARIAL | RED |

## G. Doc-contract amendments

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-54 | `TestReq54_ResolvesEvaluationOrderListGainsAStepZero` | BOUNDARY | RED |
| REQ-55 | `TestReq55_TheDispositionCardinalityDocContractsAreScopedToNilError` | ADVERSARIAL | RED |
| REQ-56 | `TestReq56_TheProducerObligationIsStatedOnTheRowDocContract` | BOUNDARY | RED |
| REQ-57 | `TestReq57_TheValidationSurfaceRecordsWhatIsNotChecked` | ADVERSARIAL | RED |

These four assert on the SOURCE of `internal/resolve/resolve.go`, parsed via
`go/ast`. The doc comments in question are themselves REQ-1 contract
artifacts of RDR 0001 that this RDR is explicitly authorized to amend, so
their content is the assertion subject. This is not an assertion on a log
string or a rendered message — no test in this suite reads `Error()` output
for anything but a failure report.

## H. CLI surfacing (C7)

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-58 | `TestReq58_ABreachReachingTheCLIIsWrappedWithCodeGroupIdentityAndHint` | HAPPY PATH | RED |
| REQ-59 | `TestReq59_TheRowIdentitiesReachTheSerializedEnvelope` | ADVERSARIAL | RED |
| REQ-60 | `TestReq60_TheCarrierIsFindingsAndNotDetail` | BOUNDARY | RED |
| REQ-61 | `TestReq61_TheCarrierIsAStructuredRepeatedSerializedField` | BOUNDARY | SAT (`CLIError.Findings` ships from RDR 0008) |
| REQ-62 | `TestReq62_ClierrStaysALeafAndTheConversionIsTheVerbLayers` | ADVERSARIAL | RED |
| REQ-63 | `TestReq63_ThePerIdentityCountSerializesAlongsideTheIdentities` | BOUNDARY | RED |
| REQ-64 | `TestReq64_TheStableCodeIsSpelledExactly` | BOUNDARY | RED |
| REQ-65 | `TestReq65_TheStaleCauseCommentIsAmended` | BOUNDARY | RED |
| REQ-66 | `TestReq66_TheOutputContractDocumentsTheIdentityCarrier` | BOUNDARY | RED |
| REQ-67 | `TestReq67_TheExitCodeCannotDiscriminateAMissingWrap` | ADVERSARIAL | RED |
| REQ-68 | `TestReq68_TheRefusalTaxonomyStaysClosed` | ADVERSARIAL | RED |
| REQ-69 | `TestReq69_TheRefusalToExitCodeMappingIsUnchanged` | BOUNDARY | RED (mapping half SAT; breach half RED) |
| REQ-109 | `TestReq109_ABreachSurfacedThroughTheCLICarriesCodeIdentityAndHint` | HAPPY PATH | RED |
| REQ-109 | `TestReq109_TheCLICarriesEveryOffendingIdentityInKernelOrder` | ADVERSARIAL | RED |
| REQ-109 | `TestReq109_ANonBreachKernelErrorStillFallsThroughToTheGenericBranch` | ADVERSARIAL | SAT (guards the classifier against blanket-recoding) |

Every C7 assertion is on the typed `Code` and the structured `Findings`,
never on the exit code alone. `TestReq67` pins the reason directly: an
unwrapped kernel error also reaches exit 2 through `cobraErrorToCLIError`,
so the exit code cannot detect a missing wrap.

## I. Executable test scenarios

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| **REQ-70 (REQ-MVV)** | `TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel` | HAPPY PATH | RED |
| REQ-71 | `TestReq71_TheNormativeFixtureValuesArePinned` | BOUNDARY | SAT (the fixture literals are this suite's own) |
| REQ-72 | `TestReq72_BreachYieldsErrorNotPlan` | HAPPY PATH | RED |
| REQ-73 | `TestReq73_DormantRowStillErrors` | ADVERSARIAL | RED |
| REQ-74 | `TestReq74_EmptyNotNilWritesConforms` | INPUT EDGE | RED |
| REQ-75 | `TestReq75_TheDiscriminatingAssertionIsOnTheInputRow` | BOUNDARY | RED |
| REQ-76 | `TestReq76_TheBreachPrecedesEveryModeledDisposition` | ADVERSARIAL | RED |
| REQ-77 | `TestReq77_TheFrozenSuitesEscapeFixturesConform` | HAPPY PATH | RED |
| REQ-78 | `TestReq78_ThePreExistingFixtureSetProducesAnIdenticalOutcomeSet` | DOMAIN EDGE | RED |
| REQ-79 | `TestReq79_TheOutcomeSetIsRegenerableRatherThanFrozenToAnArtifact` | DOMAIN EDGE | RED |
| REQ-80 | `TestReq80_TheCheckIsNonVacuous` | ADVERSARIAL | RED |
| REQ-81 | `TestReq81_TheFrozenFixturesAreInsensitiveToTheCheck` | ADVERSARIAL | RED |
| REQ-82 | `TestReq82_RestoringTheWritesMakesResolveFail` | ADVERSARIAL | RED |
| REQ-83 | `TestReq83_TheExportedPredicateMatchesTheEntryCheckOnScenarios1To3` | HAPPY PATH | RED |
| REQ-84 | `TestReq84_TheZeroResultCallerTrapIsAssertedDirectly` | ADVERSARIAL | RED |
| REQ-85 | `TestReq85_OneCallReportsEveryBreachingIdentity` | HAPPY PATH | RED |
| REQ-86 | `TestReq86_EveryElementTypeAssertsAndIsClassified` | ADVERSARIAL | RED |
| REQ-87 | `TestReq87_TwoSupplyOrdersReportTheIdenticalSequence` | ADVERSARIAL | RED |
| REQ-88 | `TestReq88_TheOrderingIsPinnedAcrossBothIdentityClasses` | BOUNDARY | RED |
| REQ-89 | `TestReq89_SharedIdentitiesCollapseIdenticallyUnderBothOrders` | BOUNDARY | RED |
| **REQ-90** | — | — | **DEFERRED** (TS scenario 8; the record scopes it to RDR 0002's build. This RDR's side is discharged as REQ-97/REQ-98.) |
| **REQ-91** | — | — | **DEFERRED** (TS scenario 11; the record marks it verbatim "not executable by this RDR". Its substance lands as REQ-58…REQ-69 and REQ-109.) |
| REQ-92 | `TestReq92_NoEscapedPlanEverCarriesAWrite` | HAPPY PATH | RED |

### REQ-70 (REQ-MVV) — the round-trip / inverse question

RDR 0009 declares NO encode/decode, import/export, or inverse operation.
The MVV's own claim is therefore not a round-trip but a **zero disposition
change** claim, and it is asserted as VALUE equality, not as a green exit:
subtest 2 compares the conformed table's escaped `*resolve.Plan` field for
field against an explicit `want` literal (`RuleID`, `SourceLocator`,
`NextTags`, `Writes`, `Revision`, `Escaped`) via `reflect.DeepEqual`.
Subtest 3 additionally replays every shipped refusal fixture and compares
dispositions value-identically. "Did not error" is nowhere sufficient.

## J. Fixture conformance (Phase 2)

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-93 | `TestReq93_TheWriteBearingEscapeFixturesAreConformed` | BOUNDARY | RED |
| REQ-94 | `TestReq94_NextTagsStaysOnTheEscapeRowBuilder` | BOUNDARY | RED |
| REQ-95 | `TestReq95_TheCheckAndTheFixtureConformanceLandTogether` | ADVERSARIAL | RED |
| REQ-96 | `TestReq96_TheErrorsAndFmtCapabilitiesAreBothPresent` | BOUNDARY | RED |

## K. Phase 3 — shared normalizer conformance fixtures

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-97 | `TestReq97_TheSharedEscapeConformanceFixtureSetBindsTheAuthoredPath` | HAPPY PATH | SAT |
| REQ-98 | `TestReq98_TheCanonicalAuthoredClearPinsTheSentinelRepresentation` | BOUNDARY | SAT |

The fixture set ships as `EscapeShapeConformanceFixtures()` in
`internal/table/escape_shape_0009_test.go` — five named cases with a
`WantCategory` expectation each, exercised against the real
`table.Load`. It is a *strengthening* over the record's Phase 3 scoping
(which asks only that the fixtures be shipped, deferring their assertion to
0002's build), permitted by the req-list ASSUMPTION on REQ-90/97/98:
"if a later phase finds executing them against `internal/table/` is free,
running them is a strengthening, not a scope breach."

## L. Performance and non-goals

| REQ | Test | Label | Status |
| --- | --- | --- | --- |
| REQ-99 | `TestReq99_ThePreconditionIsOneLinearPassOverTheRows` | BOUNDARY | RED |
| REQ-100 | `TestReq100_TheConformingScanIsAllocationFree` | BOUNDARY | RED |
| REQ-101 | `TestReq101_NoCachingOrMemoizationIsIntroduced` | ADVERSARIAL | RED |
| REQ-102 | `TestReq102_NoHashOrCanonicalSerializationIsIntroduced` | BOUNDARY | SAT (non-goal; asserts an absence) |
| REQ-103 | `TestReq103_AnEscapePlanCarryingWritesCannotBeEmitted` | ADVERSARIAL | RED |
| REQ-104 | `TestReq104_NoModeledConditionMigratesToTheErrorPath` | ADVERSARIAL | RED |
| REQ-105 | `TestReq105_TheKernelVocabularyStandsUnchanged` | BOUNDARY | SAT (no-change claim) |
| REQ-106 | `TestReq106_TheDormantMalformedRowErrorsOnEveryResolve` | DOMAIN EDGE | RED |
| REQ-107 | `TestReq107_EveryEscapeFixtureSiteSatisfiesThePrecondition` | BOUNDARY | RED |
| REQ-108 | `TestReq108_TheBreachReportUsesTheSameSortedRowRefOrder` | BOUNDARY | RED |

---

## Orphans

**REQs with no test**: REQ-90, REQ-91 — both DEFERRED by the record's own
Testing-Strategy scope note, as directed. No other REQ is uncovered.

**Tests citing no REQ**: none. Every `Test*` function in the five files
opens with a `// REQ-N:` header quoting the clause. The three helper files'
non-test functions (`breachElements`, `breachRefs`, `breachCounts`,
`sameRefs`, `sameCounts`, `nonVacuityGate`, `docCommentOf`, `fieldDocOf`,
`resolveErr`, `wrapped`, `breachedKernelError`, `causeFieldComment`,
`loadEscapeShape`, `EscapeShapeConformanceFixtures`) are fixtures and
readers, not assertions, and cite no REQ by design.

---

## Non-vacuity of the red gate

The three packages carrying new tests do not compile at HEAD, which is a
blunt red. To confirm each test discriminates rather than merely failing to
build, the gate was re-run against a **minimal declaration-only stub**:
`Table.CheckValid` returning `nil` unconditionally (Mutant A), the
`EscapeShapeBreachError{Ref, Count}` / `ErrEscapeShapeBreach` declarations,
and `kernelResolveFailure` delegating to the shipped generic
`codeAccessorFailed` branch.

Under that stub, a first pass showed seventeen tests passing that should
not have — every "both call sites agree" and "conforming fixture still
resolves" assertion holds vacuously when the predicate answers `nil` to
everything. Those tests now install `nonVacuityGate(t)` as their first
statement, which requires `CheckValid` AND `Resolve` to reject the canonical
breaching table before any agreement assertion runs. All seventeen are red
under the stub.

The tests that remain green under the stub are the ones marked **SAT**
above plus the pure spelling/shape assertions (REQ-26, REQ-29, REQ-31,
REQ-32, REQ-37, REQ-48), which the stub satisfies *because it declares the
pinned names* — that is the spelling contract doing its job, not a
tautology, and at HEAD they do not compile at all.

---

## ASSUMPTIONS carried from `req-list.md`

Restated here only where they shaped a test:

- **REQ-22 placement is not observable.** The record itself says placement
  is "a cheapness preference, not an observable contract", so no test
  asserts `assemble` was not called. `TestReq22_…` instead pins the
  observable consequence: the breach errors across every variation of the
  input channels `assemble` reads.
- **REQ-80…REQ-82, the mutation check** is a recorded build procedure
  (revert, run, restore), not a permanent source-mutating test. The
  permanently assertable substance is shipped as `TestReq80_TheCheckIsNonVacuous`
  (the check separates breaching from conforming), `TestReq81_…` (the frozen
  fixtures are insensitive to Mutant A), and `TestReq82_…` (Mutant B applied
  to a VALUE rather than to source: restoring `Writes` on a builder-produced
  row makes `Resolve` fail, proving the check is wired into the real path).
- **The error message wording is the implementer's.** REQ-96's test asserts
  only that the per-row `Error()` renders a non-empty string, since REQ-43
  requires diagnostic prose but REQ-27/REQ-44/REQ-89 forbid any test reading
  structure back out of it.
- **`CheckValid` has a value receiver.** `TestReq48_…` and
  `TestReq31_…` look it up via `reflect.TypeOf(resolve.Table{}).MethodByName`,
  which finds a value-receiver method. A pointer receiver would fail those
  two — matching the req-list's stated assumption rather than silently
  admitting both.
- **REQ-77/REQ-78's "frozen suite"** is read as the shipped disposition
  fixture set (`allDispositionInputs`), and the oracle is regenerated in
  process (REQ-79) rather than diffed against the literal `a3-conformed.out`
  bytes, which are not in this tree.

---

## QUESTIONS

Recorded rather than asked; each carries a defensible reading the tests
proceed under.

### Q-A — How does a breach reach the CLI, given RDR 0002 blocks the authored path?

`internal/table/normalize.go` rejects an escape rule carrying a write block,
a clear list, or an empty one, so **no TOML document can produce a breaching
kernel table**. An end-to-end `flow resolve` over a malformed model therefore
yields `flow-model-invalid` at load, never `escape-row-shape-breach`. That is
not a defect — it is exactly the layering C2/C3 describe — but it means the
CLI obligation cannot be exercised by driving `runCmd` over a fixture model.

**Proceeding**: the C7 tests exercise the verb layer's classifier over a
REAL kernel breach, produced by the real `resolve.Resolve` over a hand-built
table — the non-TOML producer path REQ-47 explicitly exists to serve.
Nothing is mocked: both the error under test and the classifier under test
are production code. The test binds the classifier as
`kernelResolveFailure(error) *clierr.CLIError`, sibling to the shipped
`kernelFailure(resolve.Refusal) *clierr.CLIError` that `flow_resolve.go`
already routes refusals through. Naming it is the test's binding on the
implementation, in the same way `0009:C4` binds the kernel's four spellings.
An implementation that inlines the logic in `runFlowResolve` instead must
extract it to satisfy these tests; that is a deliberate, minimal design
constraint and the smallest surface carrying the whole C7 obligation that a
test can reach.

**Alternative considered and rejected**: adding a production test seam to
`flow_resolve.go` (a package-level model override) so `runCmd` could drive a
breaching table end to end. Rejected because it adds production surface for
testing alone, which no predecessor RDR in this repo does; RDR 0008's
`ResolveBypassingPreconditionForTest` is a test-binary-only file, and its
precedent argues for a test-reachable function rather than a runtime hook.

### Q-B — Does `clierr.Finding` gain a `Count` field?

Yes, per the req-list's own Q1 reading (a): a
`Count int json:"count,omitempty"` field on `clierr.Finding`.
`TestReq63_…` asserts the count on the SERIALIZED bytes, so readings (b)
(repeat the finding N times — excluded by REQ-42's collapse) and (c) (render
into prose — excluded by REQ-44/REQ-27) both fail it. `Count == 1` elides on
every non-degenerate path, so no shipped 0005/0006/0008 assertion changes.
The 0008 test that pins `Finding`
(`reserved_key_0008_test.go::TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint`)
asserts field PRESENCE, not an exact field set, so an addition does not
break it. No exact-set assertion on `Finding` was found elsewhere in the
tree.

### Q-C — Which `Finding` field carries the RuleID?

`0009:C7` binds the carrier class, not the field mapping. `clierr.Finding`
offers both `Rule` ("the source rule/context id") and `Param`. `Rule` is the
semantically exact home. `TestReq58_…` accepts EITHER (`Rule` or `Param`),
so the mapping stays the implementer's; `TestReq109_…` asserts only that the
rule id appears somewhere in the serialized finding, and
`TestReq109_TheCLICarriesEveryOffendingIdentityInKernelOrder` keys the
ordering assertion on `Locator`, which C7 unambiguously binds to
`SourceLocator`.

### Q-D — Is the breach code `flow`-prefixed?

No. Every shipped code in `flow_input.go` is `flow`-prefixed, but REQ-64
fixes this one literally as `"escape-row-shape-breach"`.
`TestReq64_TheStableCodeIsSpelledExactly` asserts BOTH the literal and the
absence of the prefix, so a reviewer reading the inconsistency as a typo
cannot "fix" it silently. Same disposition as the req-list's Q4.

### Q-E — JD-5 (precondition precedence vs RDR 0008) is still open.

RDR 0009's Status line carries `[joint decision → JDR 0001 §JD-5:
precondition precedence]`, and JD-5 says "either order is defensible — pick
one and pin it". No REQ in this list binds the order in which
`CheckInput` (0008) and `CheckValid` (0009) run at `Resolve`'s entry.

**Proceeding**: no test asserts the relative order of the two preconditions.
`TestReq109_ANonBreachKernelErrorStillFallsThroughToTheGenericBranch` uses a
0008 reserved-key breach on a table with NO escape rows, and every 0009
breach fixture is free of reserved keys, so no test in this suite can be
decided by the ordering either way. An implementation is free to pick either
order without failing this suite, which is what an open joint decision
requires.

### Q-F — REQ-65's stale `Cause` comment: whose repair is it?

The comment is *already* false at HEAD — `Findings` shipped from RDR 0008
beside it — so this RDR inherits the staleness rather than creating it.
REQ-66 puts the repair in scope here and `TestReq65_TheStaleCauseCommentIsAmended`
is written against it. If a later phase judges the repair belongs to RDR
0008's landing instead, that is a **recorded deviation**, not a silently
dropped edit, per REQ-65's own instruction.

### Q-G — `docs/cli-output-contract.md` already documents `findings`.

REQ-66's doc half is largely discharged; `TestReq66_…` therefore asserts
`findings` AND the `count` key, the genuinely new half Q-B introduces. It
does not re-assert what RDR 0005's
`TestReq131_TheOutputContractDocumentsFindingsAndTheExit3Rule` already pins,
and no assertion in that 0005 test was weakened or touched.

---

## Deviations recorded by this phase

None. No pre-existing test was edited, weakened, or deleted; no non-test
source file was modified. `git diff --stat` against the launch base touches
only the five new test files and this artifact.

---

## REQ-MVV run record (Phase 2)

`REQ-70` / `0009:MVV`, runner
`TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel`, run at the Phase 2
implementation commit. A green exit is not the evidence; the ACTUAL observed
values are recorded below.

```
$ go test ./internal/resolve/ -run TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel -v
=== RUN   TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel
--- PASS: TestMVV0009_EscapeRowShapeConformanceIsOwnedByTheKernel (0.00s)
    --- PASS: .../1_breaching_table_errors_naming_the_row (0.00s)
    --- PASS: .../2_conformed_table_resolves_value_for_value (0.00s)
    --- PASS: .../3_conformed_tables_still_escape_and_still_refuse (0.00s)
    --- PASS: .../4_validator_and_entry_check_are_one_predicate (0.00s)
ok  	github.com/cwensel/intrastate/internal/resolve	0.243s
```

Observed values, printed from the same fixtures the runner asserts on:

```
LEG1 Result={Plan:<nil> Refusal:<nil>} Refused=false err!=nil=true
LEG1 errors.As Ref={RuleID:rdr.escape.needsowned SourceLocator:flows/rdr.toml:90} Count=1
LEG1 errors.Is(ErrEscapeShapeBreach)=true
LEG1 aggregate Unwrap() []error len=1
LEG1 Error()=resolve: escape row "rdr.escape.needsowned" at "flows/rdr.toml:90" carries writes (1 breaching row); an escape row must carry no writes
LEG2 err=<nil> Plan={RuleID:rdr.escape.needsowned SourceLocator:flows/rdr.toml:90 NextTags:[{Key:status Value:Blocked}] Writes:[] Revision:rev-nomatch Escaped:true} Refusal=<nil>
PRED breaching CheckValid!=nil=true conforming CheckValid=<nil>
```

Reading each of the MVV's four obligations against that output:

1. **The breaching table errors, naming the row, with no disposition.**
   `Result{Plan:<nil> Refusal:<nil>}` is the zero `Result`; `err != nil`;
   `errors.As` recovers `RowRef{"rdr.escape.needsowned",
   "flows/rdr.toml:90"}` — the record's normative fixture identity —
   structurally, with `Count == 1`; `errors.Is` classifies the aggregate as
   `ErrEscapeShapeBreach`; and the `errors.Join` wrapper is present even at
   cardinality 1 (`Unwrap() []error` has length 1).
2. **The same table with the writes removed resolves, value for value.**
   `err` is nil and the plan is
   `{rdr.escape.needsowned, flows/rdr.toml:90, NextTags:[status=Blocked],
   Writes:[], rev-nomatch, Escaped:true}` — the escape row's own plan, with
   no writes and `Escaped` set.
3. **The conformed table still escapes and still refuses.** `Escaped:true`
   above; and every shipped refusal fixture still travels the `Result` value
   with a nil error (subtest 3, PASS).
4. **One predicate, two call sites.** `CheckValid` is non-nil on the
   breaching table and plain `nil` on the conforming one, matching
   `Resolve`'s verdict at both.

### C7 CLI surfacing — actual wire output (REQ-58 … REQ-69, REQ-109)

Emitted through the real `clierr.EmitJSON` over the real
`kernelResolveFailure` classifier, fed a real kernel breach:

```
exit=2
{"code":"escape-row-shape-breach","message":"the transition table carries an escape row that also carries writes; an escape row describes no owned-state mutation","hint":"fix the table producer: an escape row must carry no writes","findings":[{"code":"escape-row-shape-breach","message":"this escape row carries writes","locator":"flows/rdr.toml:90","hint":"fix the table producer: an escape row must carry no writes","rule":"rdr.escape.needsowned","count":1}]}
```

And the degenerate producer — three breaching rows sharing `RowRef{"",""}` —
collapsing to ONE finding whose `count` is the only diagnostic content left:

```
{"code":"escape-row-shape-breach","message":"the transition table carries an escape row that also carries writes; an escape row describes no owned-state mutation","hint":"fix the table producer: an escape row must carry no writes","findings":[{"code":"escape-row-shape-breach","message":"this escape row carries writes","hint":"fix the table producer: an escape row must carry no writes","count":3}]}
```

`cause` is absent from both envelopes (it is `json:"-"`), `detail` carries no
identity, and the identities and count reach the wire only through
`findings[]`.

## Deviations recorded by Phase 2

- **D3** — `clierr.Finding` gains `Count int json:"count,omitempty"`.
  SPEC-UNDER, `needs author decision`, proceeding under the unattended
  override on Q-B's reading (a). See `deviations.md`.

No Phase 1 test was edited, weakened, or deleted. The only pre-existing test
files touched are the three REQ-93 names — `fixtures_test.go`,
`adversarial_test.go`, `fixup_test.go` — where the RDR's own Phase 2 plan
mandates dropping `Writes` from the escape-row builder and its two call-site
overrides; each edit removes a fixture value the RDR forbids, and no
assertion in those tests was changed.

## Mutation check run record (REQ-80 … REQ-82, TS scenario 6)

Executed as the recorded build procedure the ASSUMPTIONS section describes
(revert, run, restore) rather than as a shipped source-mutating test. Both
mutants were run with this RDR's own new tests **excluded from the oracle**:
the five `*_0009_test.go` files were held out of `internal/resolve` for both
runs, so the oracle is exactly the pre-existing frozen suite.

**Mutant A** — `Table.CheckValid` returns `nil` unconditionally
(`return nil` inserted as its first statement), fixtures conformed.

```
$ go test ./internal/resolve/
ok  	github.com/cwensel/intrastate/internal/resolve	0.316s
```

The frozen suite still passes, as REQ-81 expects: every pre-existing fixture
conforms, so replacing the check with `return nil` changes none of their
answers and the two mutants stay isolated.

**Mutant B** — the real `CheckValid` restored, the breach re-introduced into
the fixtures by restoring `Writes` on `fixtures_test.go::escapeRow`.

```
$ go test ./internal/resolve/ 2>&1 | grep -c '^--- FAIL'
9
FAIL	github.com/cwensel/intrastate/internal/resolve	0.235s
```

The frozen suite fails in nine tests, as REQ-82 expects: the entry check is
wired into the real evaluation path and is not dead code — a
callable-but-unwired predicate would have left the suite green here.

Both mutations were reverted; `git status --porcelain` showed only this
artifact and `deviations.md` modified afterwards, and the restored suite is
green.
