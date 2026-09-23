# Coverage — RDR 0012 Declared-kind carrier at the guard value seam

Phase 1 artifact. Each REQ-N in `req-list.md` is mapped to the test that
would fail if a future change broke that clause. An uncovered REQ keeps its
row with column 2 EMPTY. None is uncovered, and no REQ is proposed for
demotion.

Test files (all committed RED on `worktree-rdr-0012`):

- `internal/guard/seam_carrier_0012_test.go`: the `0012:C1`–`C4` guard-side
  tests (source and reflection checks over surfaces that exist today) and
  the drivers for the behavioural probes.
- `internal/guard/rdr0012_probe_test.go`: behavioural probes that call
  `guard.NewEvaluator`, `guard.DeclaredKinds`, `resolve.ContractKinds` and the
  constructor-taking `resolve.TestGuardEvaluatorContract` directly. The file
  sits behind the `rdr0012probe` build tag, so `go vet ./...` and the default
  `go test` never compile it. The drivers run it in a
  `go test -tags rdr0012probe` child. A probe that does not compile shows up
  as an `owed` assertion failure, not a build break. Probes named
  `_Mutant_` are expected to FAIL, and their drivers check that they do.
- `internal/resolve/contract_fixture_0012_test.go`: the kernel half
  (`0012:C3` fixture and suite shape, `0012:C1` resolution-path rule), read
  through `go/parser`.
- `internal/cli/guard_typed_0012_test.go`: end-to-end tests through the
  command tree. They cover the owned door (the artifact is seeded with
  `flow set-state`, then the reader, `OwnedSnapshot` and `assemble`), the
  `--tag` door, `flow next`, the `lint --model` C5 refusal, the corpus diff
  (REQ-50) and the REQ-MVV runner.
- `internal/cli/testdata/rdr0012_corpus_verdicts.json`: the REQ-50 "before"
  golden, 123 models linted on the pre-change tree
  (`RDR0012_WRITE_CORPUS_GOLDEN=1` rewrites it).

REQ-MVV runner:
`internal/cli/guard_typed_0012_test.go::TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor`
covers steps 2, 3 (with its control), 4 and 5. For step 4 it runs a probe
child.

| REQ | Test |
| --- | --- |
| REQ-1 | TestReq1_0012_NewEvaluatorIsDeclaredWithTheCarrierSignature |
| REQ-2 | TestReq2_0012_DeclaredKindsIsHomedInGuardAndMapsKeyToKind |
| REQ-3 | TestReq3_0012_DeclaredKindsIsTotalAndOmitsEmptyKind |
| REQ-4 | TestReq4_0012_DeclaredKindsOverANilModelIsEmptyNonNil |
| REQ-5 | TestReq47_0012_EvaluatorHoldsTheDeclarationMappingAndNothingElse |
| REQ-6 | TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction |
| REQ-7 | TestReq49_0012_AZeroValueEvaluatorNeverRevertsToRawComparison |
| REQ-8 | TestReq8_0012_TheSeamSignatureAndAtomShapeAreUnchanged |
| REQ-9 | TestReq9_0012_TheResolutionPathReadsNoDeclarations |
| REQ-10 | TestReq10_0012_NewEvaluatorReturnsTheConcreteEvaluator |
| REQ-11 | TestReq11_0012_GuardAtomsCompareTypedAndMatchAtomsCompareBytes |
| REQ-11 | TestReq11_0012_GuardEqAndInCompareUnderTheDeclaredKind |
| REQ-12 | TestReq12_0012_IntComparesParsedValues |
| REQ-13 | TestReq13_0012_OneUnparseableInMemberPoisonsTheList |
| REQ-14 | TestReq14_0012_OverflowIsUnevaluable |
| REQ-15 | TestReq15_0012_BoolComparesTokensOnly |
| REQ-16 | TestReq16_0012_EnumAndScalarCompareExactStrings |
| REQ-17 | TestReq17_0012_SetEqAndInAreUnevaluable |
| REQ-18 | TestReq18_0012_AbsentKeyIsUnevaluable |
| REQ-19 | TestReq19_0012_EmptyKindDeclaredKeyIsUnevaluable |
| REQ-20 | TestReq20_0012_UnparseableLiteralIsUnevaluable |
| REQ-21 | TestReq21_0012_UnknownKindTokenIsUnevaluable |
| REQ-22 | TestReq22_0012_OrderingAndContainsDoNotConsultTheMapping |
| REQ-23 | TestReq23_0012_UnknownOperatorIsUnevaluableAndNeverPanics |
| REQ-24 | TestReq24_0012_InMembersAreTheDecodedArray |
| REQ-25 | TestReq25_0012_TheSuiteCarriesTheTypedAndDefensiveLegs |
| REQ-26 | TestReq26_0012_EachKindTokenHasADiscriminatingCase |
| REQ-27 | TestReq27_0012_ContractKindsIsAFunctionNotAPackageVar |
| REQ-27 | TestReq27_0012_ContractKindsHandsEachCallerAFreshCopy |
| REQ-28 | TestReq28_0012_TheSuiteTakesASeamConstructor |
| REQ-29 | TestReq29_0012_CasesAreKeyedOntoAdmittedKinds |
| REQ-30 | TestReq30_0012_TheConformingTestSeamIsBoundThroughAConstructor |
| REQ-31 | TestReq31_0012_TheGuardSuiteCallIsRePointedThroughAOneLineAdapter |
| REQ-32 | TestReq32_0012_TheImportablePinIsReTypedNotDeleted |
| REQ-33 | TestReq33_0012_TheCLISiteConstructsOverTheRequestsModel |
| REQ-33 | TestReq33_0012_TheLintSiteConstructsOverTheModelsDeclarations |
| REQ-34 | TestReq34_0012_TheEvaluatorIsThreadedDownNotReDerived |
| REQ-35 | TestReq35_0012_AtomAdmitsValueIsTheSoleRawByteSiteAndKeepsItsShape |
| REQ-36 | TestReq36_0012_ATypedUnevaluableIsNeitherPrunedNorMatched |
| REQ-37 | TestReq37_0012_TheUserFacingDocsStateTheGuardMatchSplit |
| REQ-38 | TestReq38_0012_ATypedUnevaluableRidesTheExistingRefusal |
| REQ-39 | TestReq39_0012_NonCanonicalPredicateIntLiteralsAreRefused |
| REQ-40 | TestReq40_0012_TheBareFloatNegativeZeroIsRefused |
| REQ-41 | TestReq41_0012_TheCheckReachesPredicateLiteralsOnly |
| REQ-42 | TestReq42_0012_CLIValuesStayAcceptedAndCompareTyped |
| REQ-43 | TestReq43_0012_TheRefusalNamesTheSiteAndTheRewrite |
| REQ-44 | TestReq44_0012_TheRefusalReusesMalformedPredicateAtom |
| REQ-45 | TestReq45_0012_CanonicalIsExactlyTheItoaRoundTrip |
| REQ-46 | TestReq46_0012_PerKindDispatchMatrix |
| REQ-47 | TestReq47_0012_EvaluatorHoldsTheDeclarationMappingAndNothingElse |
| REQ-48 | TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction |
| REQ-49 | TestReq49_0012_AZeroValueEvaluatorNeverRevertsToRawComparison |
| REQ-50 | TestReq50_0012_NoCommittedModelsLintVerdictFlips |
| REQ-51 | TestReq51_0012_Cover07RefusesWithTheCanonicalSpellingDiagnostic |
| REQ-52 | TestReq52_0012_TheTagDoorRefusesUpstreamWhileTheOwnedDoorReachesTheSeam |
| REQ-53 | TestReq53_0012_AnOwnedNonCanonicalValueMatchesUnderParsedComparison |
| REQ-54 | TestReq54_0012_FlowNextReportsTypedVerdicts |
| REQ-MVV | TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor |
| REQ-MVV | TestReqMVV_0012_Step4_TheExtendedSuiteIsGreenOverNewEvaluator |

## Notes for Phase 2

- **Preservation clauses** (REQ-8, REQ-35, REQ-41, REQ-42, REQ-50, REQ-52)
  hold on `main` today. Each test pairs the preserved half with the
  observable that turns it red on this tree: the new surface, the typed
  verdict, or the canonical refusal. So no test is green against a
  missing implementation. The preserved halves (the S8 envelope, the
  admitted `[initial]`/`[rule.write]`/CLI spellings, the corpus golden, and
  the atom shape) pass today.
- **REQ-43** pins the site as the exact prefix that the loader's existing
  `malformed_predicate_atom` error for the same atom carries. This is the
  req-list ASSUMPTION on `<site>`: the test derives the prefix from a
  `"x" is not an int` baseline and never hard-codes it.
- **REQ-53 / MVV step 3**: with the unguarded `fallback` on the same
  outcome, "the guarded row MATCHES" can be observed only as
  `flow-ambiguous-match` naming `guarded` beside `fallback`. The
  record's own spike shows the `iter=7` control behaving the same way.
- **Existing tests Phase 2 must migrate** (the RDR names them; they are not
  edited here): `guard_evaluator_0003_test.go::TestReq34…` (its zero-field
  check is narrowed per A2, and the line-67 suite call is re-pointed),
  `guard_mvv_test.go` (the `var fn` re-type and `conformingContractSeam`'s
  constructor), and the `guardSeam()` call sites in
  `escape_shape*_0009_test.go`.
- **Satisfiability check**: the red suite was also run against a
  throwaway, spec-shaped patch in a scratch copy outside the worktree, and
  none of that patch was committed. Every probe-driven test passed,
  including the REQ-25/26 mutants and REQ-29's re-keying. So did every
  owned-door and C5 CLI test. The tests that failed there were only the
  source-shape ones that the scratch patch did not attempt.

## Red confirmation

`go test -count=1 -run '_0012_' ./internal/guard ./internal/resolve ./internal/cli`
(via `rdr-leg-test`) on the committed tree: all 56 new top-level tests FAIL.
There are 32 in `internal/guard`, 5 in `internal/resolve` and 19 in
`internal/cli`, and each fails with its REQ's own assertion. The pre-commit
gate (`gofmt -l .`, `go vet ./...`) passes, because nothing outside the
`rdr0012probe` tag names an absent symbol.

## REQ-MVV runner

`go test -count=1 -v -run 'TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor' ./internal/cli`
(via `rdr-leg-test`, Phase 2 leg 1)

## REQ-MVV output

```text
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_2:_many_refuses_uncomparable
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_3:_07_matches_parsed
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_3_control:_4_prunes
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_4:_the_extended_suite
=== RUN   TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_5:_lint_refuses_00
--- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor (1.36s)
    --- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_2:_many_refuses_uncomparable (0.04s)
    --- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_3:_07_matches_parsed (0.05s)
    --- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_3_control:_4_prunes (0.07s)
    --- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_4:_the_extended_suite (1.18s)
    --- PASS: TestReqMVV_0012_DeclaredKindReachesTheSeamThroughTheOwnedDoor/step_5:_lint_refuses_00 (0.02s)
PASS
ok  	github.com/cwensel/intrastate/internal/cli	1.641s
```
