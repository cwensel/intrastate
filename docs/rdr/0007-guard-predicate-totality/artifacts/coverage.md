# Coverage — RDR 0007 Phase 1 spec tests (red)

Maps every REQ in `req-list.md` to the test(s) that would FAIL if a future
change broke that clause, and flags orphans in both directions.

- Test files: `internal/resolve/guard_atoms_test.go`,
  `internal/resolve/guard_totality_test.go`,
  `internal/resolve/guard_mvv_test.go`,
  `internal/resolve/guard_fixtures_test.go` (fixture builders only).
- Package: `resolve_test` (REQ-73). Framework: Go stdlib `testing`.
- Every test opens with `// REQ-N: "<quote>"` and carries one of
  `HAPPY PATH | INPUT EDGE | BOUNDARY | ADVERSARIAL | DOMAIN EDGE`.

## REQ × test

| REQ | Test | Label | File |
| --- | --- | --- | --- |
| REQ-1 | `TestReq1_RowCarriesItsGuardAsASliceOfParsedAtoms` | HAPPY PATH | atoms |
| REQ-2 | `TestReq2_EmptyAtomSliceIsAnUnguardedRow` | INPUT EDGE | atoms |
| REQ-3 | `TestReq3_NoOpaqueGuardStringAndNoReconstructionStep` | BOUNDARY | atoms |
| REQ-4 | `TestReq4_RefusalGuardFieldIsDeleted` | BOUNDARY | atoms |
| REQ-5 | `TestReq5_AtomAndPayloadAtomShareTheFourFieldNames` | BOUNDARY | atoms |
| REQ-6 | `TestReq6_BlockIsAnExportedNamedStringTypeWithThreeConstants` | BOUNDARY | atoms |
| REQ-7 | `TestReq7_KernelDecidesPresencePerAtomBeforeConsultingTheSeam` | HAPPY PATH | atoms |
| REQ-8 | `TestReq8_ExistenceAtomsNeverReachTheEvaluator` | HAPPY PATH | atoms |
| REQ-9 | `TestReq9_AbsentKeyValueAtomIsKernelUnevaluableAndNeverReachesTheSeam` | DOMAIN EDGE | atoms |
| REQ-10 | `TestReq10_GuardEvaluatorIsASingleMethodInterfaceTakingAtomAndValue` | BOUNDARY | atoms |
| REQ-11 | `TestReq11_SeamMayAnswerUnevaluableForAPresentValue` | DOMAIN EDGE | atoms |
| REQ-12 | `TestReq12_AnAbsenceFoldingEvaluatorCannotBeReachedForAnAbsentKey` | ADVERSARIAL | atoms |
| REQ-13 | `TestReq13_NilSeamStillResolvesAKernelDecidableRow` | BOUNDARY | atoms |
| REQ-14 | `TestReq14_NilSeamYieldsUnevaluableOnlyForAPresentKeyValueAtom` | BOUNDARY | atoms |
| REQ-15 | `TestReq15_NilSeamPresentKeyAtomIsPayloadUncomparable` | DOMAIN EDGE | atoms |
| REQ-16 | `TestReq16_NilSeamIsPerAtomNotWholeGuard` | ADVERSARIAL | atoms |
| REQ-17 | `TestReq17_Fixup1dRedecidedWithAValueAtomOverAPresentKey` **+** the re-decided frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | ADVERSARIAL | atoms, `fixup_test.go` |
| REQ-18 | `TestReq18_ValueOperatorsArePartialOverTheView` | DOMAIN EDGE | totality |
| REQ-19 | `TestReq19_AbsentSetValuedTagIsUnevaluableNotTheEmptySet` | DOMAIN EDGE | totality |
| REQ-20 | `TestReq20_ExistenceIsTheSoleTotalOperator` | HAPPY PATH | totality |
| REQ-21 | `TestReq21_PolarityLivesInTheLiteralAndComposesWithBlockPlacement` | DOMAIN EDGE | totality |
| REQ-22 | `TestReq22_KernelComparesOperatorAndLiteralVerbatim` | ADVERSARIAL | atoms |
| REQ-23 | `TestReq23_ExistenceTokenAndLiteralBytesAreNormative` | BOUNDARY | atoms |
| REQ-24 | `TestReq24_ForeignExistenceTokenIsAValueAtom` | ADVERSARIAL | atoms |
| REQ-25 | `TestReq25_ForeignLiteralOnAnExistsAtomIsKernelUncomparable` | ADVERSARIAL | atoms |
| **REQ-26** | **— (intentional orphan, see below)** | — | — |
| REQ-27 | `TestReq27_PresenceIsProvenanceBlind` | DOMAIN EDGE | atoms |
| REQ-28 | `TestReq28_KeyIdentityIsExactStringEqualityWithNoCanonicalization` | ADVERSARIAL | atoms |
| REQ-29 | `TestReq29_GuardDecidabilityAndOwnedStateUseDifferentPredicates` | DOMAIN EDGE | atoms |
| REQ-30 | `TestReq30to32_EmptyAtomAndEmptyBlockIdentities` (subtest `REQ-30 …`) | INPUT EDGE | totality |
| REQ-31 | `TestReq30to32_EmptyAtomAndEmptyBlockIdentities` (subtest `REQ-31 …`) | INPUT EDGE | totality |
| REQ-32 | `TestReq30to32_EmptyAtomAndEmptyBlockIdentities` (subtests `REQ-32 …` ×2) | INPUT EDGE | totality |
| REQ-33 | `TestReq33_StrongKleeneMatrixAcrossAllAndUnless` | DOMAIN EDGE | totality |
| REQ-34 | `TestReq34_CombinationIsTableDrivenNotIntegerMinOverTheConstants` | ADVERSARIAL | totality |
| REQ-35 | `TestReq35_UnlessIsBlockLevelNegationNotPerAtom` | DOMAIN EDGE | totality |
| REQ-36 | `TestReq36_UnresolvedDependenceOnAnUnevaluableAtomYieldsUnevaluable` | DOMAIN EDGE | totality |
| REQ-37 | `TestReq37_FalseDominatesUnevaluableAcrossTwoAtoms` | ADVERSARIAL | totality |
| REQ-38 | `TestReq38_TheEvaluatorHoldsNoPartOfTheCombinationTables` | ADVERSARIAL | totality |
| REQ-39 | `TestReq39_GateThenCountOrdering` | DOMAIN EDGE | totality |
| REQ-40 | `TestReq40_EscapeSetIsGatedIdentically` | DOMAIN EDGE | totality |
| REQ-41 | `TestReq41_EscapeSetScopingBothLegs` | ADVERSARIAL | totality |
| REQ-42 | `TestReq42_GateDocCommentDoesNotRepeatTheSupersededRationale` | BOUNDARY | totality |
| REQ-43 | `TestReq43_ARowFailingBothWaysCarriesTheOwnedPayloadOnly` | DOMAIN EDGE | totality |
| REQ-44 | `TestReq44_GuardUnevaluableIsNotMaskableBehindAnEscapableClass` | ADVERSARIAL | totality |
| REQ-45 | `TestReq45_PayloadNamesEveryUndecidableRowAndAtom` | HAPPY PATH | totality |
| REQ-46 | `TestReq46_UncomparableCoversAllThreeWaysAPresentKeyFailsToDecide` | DOMAIN EDGE | totality |
| REQ-47 | `TestReq47_ReasonSetStaysAtExactlyTwoMembers` | BOUNDARY | atoms |
| REQ-48 | `TestReq48_PayloadIsNamedSurface` | BOUNDARY | atoms |
| REQ-49 | `TestReq49_ReasonIsAClosedNamedStringSetWithAnEnumerator` | BOUNDARY | atoms |
| REQ-50 | `TestReq50_PayloadBlockIsTheAtomsBlockType` | BOUNDARY | atoms |
| REQ-51 | `TestReq51_ReasonsAreProducedByTheVerdictPassNotASecondWalk` | ADVERSARIAL | totality |
| REQ-52 | `TestReq52and53_EveryUndecidableRowAppearsInThePayloadAndInRows` | DOMAIN EDGE | totality |
| REQ-53 | `TestReq52and53_EveryUndecidableRowAppearsInThePayloadAndInRows` | DOMAIN EDGE | totality |
| REQ-54 | `TestReq54_NoShortCircuitOnADecidedBlock` | ADVERSARIAL | totality |
| REQ-55 | `TestReq55_PayloadIsSortedOnTheTotalSixFieldTuple` | ADVERSARIAL | totality |
| REQ-56 | `TestReq56_LiteralIsComparedAsExactBytesInTheD13Form` | ADVERSARIAL | totality |
| REQ-57 | `TestReq57_TheKernelReportsEachAtomExactlyOnce` | ADVERSARIAL | totality |
| REQ-58 | `TestReq58_APrunedRowReportsNothing` | DOMAIN EDGE | totality |
| REQ-59 | `TestReq59_RequiresOwnedDocContractIsNarrowedToPostGuardWrites` | BOUNDARY | mvv |
| REQ-60 | `TestReq60_GuardInputCoverageIsNotRequiresOwnedsJob` | DOMAIN EDGE | mvv |
| REQ-61 | `TestReq61_ListingAGuardKeyInRequiresOwnedProtectsAgainstObservedSubstitution` | DOMAIN EDGE | mvv |
| REQ-62 | `TestReq62_ConformingEscapeRowRaisesNoOwnedStateOfItsOwn` | BOUNDARY | mvv |
| **REQ-63** | **— (intentional orphan, see below)** | — | — |
| REQ-64 | `TestReq64_SurvivorMembershipIncludesUnevaluableRows` | DOMAIN EDGE | totality |
| REQ-65 | `TestReq65_MissingOwnedIsTheDeduplicatedSortedUnionAcrossSurvivors` | BOUNDARY | totality |
| REQ-66 | `TestReq66_AnUnevaluableCandidateVetoesADecidedTrueSibling` | ADVERSARIAL | totality |
| REQ-67 | `TestReq67_D8IsRatifiedConditionalOnTheDomainRule` | DOMAIN EDGE | totality |
| REQ-68 | `TestReq68_VerdictsAreDrivenThroughAtomsNotInjectedAtRowLevel` **+** the migrated frozen fixtures (`fixtures_test.go::namedGuard`) | BOUNDARY | mvv, `fixtures_test.go` |
| REQ-69 | `TestReq69_TheGuardTextFieldIsActuallyRemoved` | BOUNDARY | mvv |
| REQ-70 | `TestReq70_DomainRuleMatrixOverOperatorsPresenceAndBlocks` (operators × presence × blocks axis; the other axes are REQ-33, REQ-30to32, REQ-76, REQ-41, REQ-24) | DOMAIN EDGE | totality |
| REQ-71 | `TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface`; `TestReq71_ContractTestExercisesThePresentValueLiteralOperatorProduct`; `TestReq71_AConformingValueSeamSatisfiesTheContractTest` | BOUNDARY, ADVERSARIAL, HAPPY PATH | mvv |
| REQ-72 | `TestReq72_PresentSetValuedTagCrossesTheSeamAsD13Bytes` **+** REQ-19's absent-key half; the §D13 carriage form is built by `guard_fixtures_test.go::d13Set` | DOMAIN EDGE | totality, fixtures |
| REQ-73 | `TestReq73_ScenariosAreKernelTestsInPackageResolveTest` | BOUNDARY | mvv |
| REQ-74 | `TestReq74_FalseVsUnevaluableDominanceMutantIsKilled` | ADVERSARIAL | mvv |
| REQ-75 | `TestReq75_ObservedTagSubstitutionYieldsASpecificVerdict` | DOMAIN EDGE | mvv |
| REQ-76 | `TestReq76_TwoRowAbsencePatternIsTotalAndTheBareValueRowIsNot` | DOMAIN EDGE | mvv |
| REQ-77 | `TestReq77_ChangeIsConfinedToTheKernelPackage` | BOUNDARY | atoms |
| REQ-78 | `TestReq78_MatchPatternStillFoldsAbsenceIntoNonMatch` | DOMAIN EDGE | atoms |
| REQ-79 | `TestReq79_NoSixthRefusalKindIsMinted` | BOUNDARY | atoms |
| REQ-80 | `TestReq80_OneSeamCallPerPresentKeyValueAtom` | BOUNDARY | atoms |
| REQ-MVV | `TestMVV_GuardPredicateTotality` (4 scenarios + 1 oracle control) | HAPPY PATH (end-to-end) | mvv |

## Orphans — REQ with no test

Two, both recorded as out-of-scope boundaries by `req-list.md` itself, not
gaps:

- **REQ-26** — "The normalizer additionally MUST reject a foreign literal
  at load (RDR 0002's typed validation)". `req-list.md` marks it **NOT a
  kernel test**; Testing Strategy row 11 says "the load-rejection leg is
  RDR 0002's and is NOT a kernel test (no load path exists under
  `internal/resolve`)". The kernel-side backstop half IS covered, by
  `TestReq25_ForeignLiteralOnAnExistsAtomIsKernelUncomparable`.
- **REQ-63** — "Who populates the field is RDR 0002's (JDR 0001 §JD-3)."
  `req-list.md`: "Out of scope for this build; recorded as a boundary."
  The kernel-side consequence of the field's *contents* IS covered, by
  REQ-59/60/61/62.

## Orphans — test citing no REQ

None. Every `func Test…` in the four new files opens with a `// REQ-N:`
quote block. The non-test helpers that carry no REQ citation are fixture
builders and oracles, listed here so the absence is deliberate:

- `guard_fixtures_test.go` — atom builders (`allAtom`, `unlessAtom`,
  `existsAll`, `existsUnless`), value-seam stubs (`valueSeam`, `eqSeam`),
  table builders (`guardedRow`, `conformingEscapeRow`, `guardInput`),
  payload readers (`undecidedAtomsFor`, `undecidedRuleIDs`, `reasonFor`),
  assertion helpers (`mustRefuse`, `mustPlan`), and `d13Set` (§D13
  carriage form).
- `guard_atoms_test.go` — `falseSeam`, `containsCall`, `hasItem`,
  `kernelDeclaresFunc`.
- `guard_totality_test.go` — `atomWithVerdict`, `assertRowVerdict`,
  `assertReason`, and the restated normative oracles `kleeneAnd`,
  `negate`, `compareUndecidedAtomsForTest` (restated in the test so the
  kernel cannot supply its own oracle for its own tables and its own sort).
- `guard_mvv_test.go` — `conformingContractSeam`, `recordingSeam`,
  `boolVerdict`, `parseInt`, `parseD13Set`, `fieldDoc`,
  `kernelPackageSource`, `hasText`, `squash`.

## Frozen RDR 0001 suite — mechanical adaptation only

Per §scope-discipline, the reshape breaks the *compile* surface of four
frozen files. Each was adapted mechanically; no behavioural assertion was
weakened.

| File | Adaptation |
| --- | --- |
| `fixtures_test.go` | `fixtureGuards.Evaluate` narrowed to the per-atom seam and keyed on `atom.Literal`; `namedGuard(name)` renders each former guard string as ONE atom over the PRESENT key `status`, so the seam's verdict — not an absence — decides the row (REQ-68's "through atoms and a value stub"); `deepCopyInput` copies `Row.Guard`. |
| `adversarial_test.go` | Guard strings → `namedGuard(...)`; ADV-3's failure message reads `Refusal.Undecided` in place of the deleted `Refusal.Guard`. Its assertion (`reflect.DeepEqual` over the whole refusal) is unchanged. |
| `fixup_test.go` | Guard strings → `namedGuard(...)`. **Fixup-1d is RE-DECIDED, not re-encoded** (REQ-17): its atom moves from the absent key `iterations` to the PRESENT key `reviews`, and the `Refusal.Guard` assertion becomes a payload assertion on reason `uncomparable`. The verdict it asserts is unchanged. |
| `resolve_test.go` | `TestReq33`'s `r.Guard == ""` diagnosis check becomes `len(r.Undecided) == 0` — the same obligation ("the refusal names what could not be decided") against the replacement surface. |

## How red was confirmed

1. **In the worktree**: `go build ./...` succeeds (no non-test source was
   touched); `go test ./internal/resolve/` fails to compile on
   `undefined: resolve.GuardAtom`, `resolve.Block`, `resolve.OpExists`,
   `resolve.UndecidedRow`, `resolve.UndecidedAtom`, `resolve.Reason`,
   `resolve.Reasons`, `resolve.TestGuardEvaluatorContract` — the 0007
   Phase 1 surface that does not exist yet. That is legitimate red for
   every test in the four new files.
2. **Anti-tautology check** (scratchpad copy of the worktree only; nothing
   written back): two throwaway kernels were built to make the suite
   compile — (a) a NAIVE one with the new types but no presence step, no
   K3 and no payload, and (b) a MOSTLY-CORRECT one adding the per-atom
   presence rule, K3 combination and per-atom reason production but not
   the payload plumbing into `gate`. Against (a), 52 of 78 new tests fail;
   against (b), 32 still fail. Every test that passed against a throwaway
   kernel was then re-run against a targeted mutant and confirmed to fail:
   `BlockMatch` dropped, `UndecidedAtom.Block` re-spelled `string`,
   `OpExists = "Exists"`, `Reasons()` returning a shared slice,
   UNEVALUABLE-dominance swapped for FALSE-dominance (the mutant that
   survives all 154 frozen tests — killed by REQ-37 and REQ-74), a
   case-folding/coercing existence rule (REQ-22, REQ-25), an empty guard
   folded to unevaluable (REQ-2, REQ-30to32), the kernel
   re-canonicalizing the §D13 literal (REQ-72), unevaluable rows pruned
   from the survivor set (REQ-64), combined owned+guard reporting
   (REQ-43), `TagSet.matches` folding absence into MATCH (REQ-78), the
   owned sweep conflated with provenance-blind presence (REQ-29, REQ-61),
   `MissingOwned` left unsorted (REQ-65), a sixth `RefusalKind` (REQ-79),
   and both a hollow and a one-case `TestGuardEvaluatorContract` (REQ-71).
   No test passes against the missing implementation for a reason
   unrelated to the spec clause it cites.
3. Against (b), the frozen RDR 0001 suite passes in full except
   `TestFixup1d` (RE-DECIDED, expected) and `TestReq33` (waiting on the
   payload plumbing Phase 2 adds) — confirming the mechanical adaptation
   preserved behaviour rather than weakening it.
