# Coverage — RDR 0008 Ownership of the recognized-outcome tag key name

Phase 1 (test authoring) artifact. Maps every REQ in
[`req-list.md`](req-list.md) to the test that would fail if a future change
broke that clause, and flags orphans in both directions.

## Colour legend

- **RED** — net-new 0008 surface. The test fails today against the missing
  implementation. These are the red gate.
- **GREEN (predecessor)** — the clause is already satisfied by an implemented
  predecessor (RDR 0002 run 2, RDR 0007 run 1, or the shipped kernel). The test
  is a **regression pin**, legitimately green, and was NOT rewritten to force
  red.

## Net-new surface this suite pins

The repository's `.githooks/pre-commit` runs `go vet ./...` on every commit, so
Phase 1 lands each symbol below as an **unimplemented declaration** (marked
`PHASE 1 DECLARATION ONLY` with a `TODO(rdr-0008 Phase 2)`) rather than leaving
the test packages uncompilable. The red gate is therefore BEHAVIOURAL, which is
the stronger form: 23 tests fail on what the code does, not on what it lacks.

Four net-new tests are green against the declarations alone, and are **not**
tautological — each fails against a plausible wrong Phase 2:

| Test | Why green now | What breaks it |
| --- | --- | --- |
| `TestReq40And69And70_…` | a scope/naming ceiling, satisfied once one predicate exists | a second exported `Input` predicate, or renaming it `Final` |
| `TestReq6And8_…` | kernel-load-bearing half is predecessor-satisfied; predicate-exists half now holds | removing the predicate, or a conformance assertion that stops reading the view |
| `TestReq88And99_…` | conforming inputs must stay nil-error | a Phase 2 predicate that over-rejects (its complement, breach detection, is RED under REQ-36/38) |
| `TestReq66And68_…` | both sites agree while neither breaches | verified against a probe implementation: a working predicate that `Resolve` does not call makes this FAIL |

Declarations landed unimplemented:

| Symbol | Package | REQs |
| --- | --- | --- |
| `CheckInput(Input) error` | `internal/resolve` | 36, 38–46, 50, 66, 68–71, 85–90, 97, 99 |
| `Resolve` entry call to it | `internal/resolve` | 38, 55, 67, 71, 87 |
| `Failure.Offending` / `.Remedy` / `.Rule` | `internal/table` | 26–31, 33–35, 73, 75, 91–93, 101 |
| `Advisory{Authored, Reserved, Rule}` | `internal/table` | 58–64, 80, 81, 100 |
| `LoadWithAdvisories(...) (*Model, []Advisory, error)` | `internal/table` | 58–64, 80, 81, 100 |
| RDR 0008 citation on `recognizedTagKey` | `internal/resolve` | 9 |

## Test files

| Path | Scope |
| --- | --- |
| `internal/resolve/reserved_key_0008_test.go` | blocks 1, 4, 5, 6; enforcement locus |
| `internal/resolve/reserved_key_sameview_0008_test.go` | scenario 3 (deviation D1) |
| `internal/resolve/reserved_key_fixtures_0008_test.go` | fixtures + structural probes |
| `internal/resolve/export_0008_test.go` | package-internal residual path (`package resolve`) |
| `internal/table/reserved_key_0008_test.go` | blocks 2, 3; near-miss advisory |
| `internal/table/reserved_key_fixtures_0008_test.go` | fixtures + payload/advisory readers |
| `internal/table/mvv_0008_test.go` | **REQ-MVV runner** |
| `internal/cli/reserved_key_0008_test.go` | §D8 / §D10 citations (deviation D3) |

New fixtures: `internal/table/testdata/pos-near-miss-{folded,upper,space}.toml`,
`pos-not-near-miss.toml`, `neg/neg-recognized-{owned,quoted-owned,declared-twice}.toml`,
`neg/neg-two-recognized-decls.toml`. All extend the promoted set (REQ-118
permits extension), and all are single-mutation derivatives of `rdr-fixture.toml`.

## REQ × test

### A. The reserved key and the assembled view (block 1)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 1 | `TestReq1_AssembledViewBindsRecognizedOutcomeAtTheReservedKey` | HAPPY PATH | GREEN (predecessor) |
| 2 | `TestReq2_ObservedTagCannotRebindTheReservedKey` | ADVERSARIAL | GREEN (predecessor) |
| 3 | `TestReq3_AbsentOutcomeYieldsAViewWithNoRecognizedKey` | BOUNDARY | GREEN (predecessor) |
| 4 | `TestReq4And5_EmptyOutcomeIsNotClosedByThisRDR` | DOMAIN EDGE | GREEN (predecessor) |
| 5 | `TestReq4And5_EmptyOutcomeIsNotClosedByThisRDR` | DOMAIN EDGE | GREEN (predecessor) |
| 6 | `TestReq6And8_ConformanceIsBehavioralAndThePredicateIsStillRequired` | ADVERSARIAL | GREEN (declaration; see table above) |
| 7 | `TestReq7_TableSideSpellingAgreesWithTheKernelBinding` | DOMAIN EDGE | GREEN (predecessor) |
| 8 | `TestReq6And8_ConformanceIsBehavioralAndThePredicateIsStillRequired` | ADVERSARIAL | GREEN (declaration; see table above) |
| 9 | `TestReq9_RecognizedTagKeyConstantCitesRDR0008` | DOMAIN EDGE | RED |

### B. Declaration naming (block 2)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 10 | `TestReq10And18_RecognizedProvenanceUnderAnotherNameFailsReservedTagKey` | HAPPY PATH | GREEN (predecessor) |
| 11 | `TestReq11_OwnedDeclarationNamedRecognizedFailsReservedTagKey` | HAPPY PATH | GREEN (predecessor) |
| 12 | `TestReq12And82And83_TwoRecognizedDeclarationsYieldExactlyOneNamingFailure` | DOMAIN EDGE | GREEN (predecessor) |
| 13 | `TestReq13And84_SameNameTwiceIsMalformedTOMLNotReservedTagKey` | BOUNDARY | GREEN (predecessor) |
| 14 | `TestReq14And17_TheScopeUnitIsOneModelWithNoCrossFileRule` | DOMAIN EDGE | GREEN (predecessor) |
| 15 | `TestReq15And16_TheMissingDeclarationLowerBoundStaysOutsideReservedTagKey` | DOMAIN EDGE | GREEN (predecessor) |
| 16 | `TestReq15And16_TheMissingDeclarationLowerBoundStaysOutsideReservedTagKey` | DOMAIN EDGE | GREEN (predecessor) |
| 17 | `TestReq14And17_TheScopeUnitIsOneModelWithNoCrossFileRule` | DOMAIN EDGE | GREEN (predecessor) |
| 18 | `TestReq10And18_RecognizedProvenanceUnderAnotherNameFailsReservedTagKey` | HAPPY PATH | GREEN (predecessor) |
| 19 | `TestReq19And56And57And79_NearSpellingsAreOrdinaryUnreservedNames` | INPUT EDGE | GREEN (predecessor) |
| 20 | `TestReq20And79_QuotedFormParsesToTheSameKeyAndIsReserved` | ADVERSARIAL | GREEN (predecessor) |
| 21 | `TestReq21And22And23_ReservedTagKeyIsScopedToDeclarationKeysOnly` | DOMAIN EDGE | GREEN (predecessor) |
| 22 | `TestReq21And22And23_ReservedTagKeyIsScopedToDeclarationKeysOnly` | DOMAIN EDGE | GREEN (predecessor) |
| 23 | `TestReq21And22And23_ReservedTagKeyIsScopedToDeclarationKeysOnly` | DOMAIN EDGE | GREEN (predecessor) |
| 24 | `TestReq24_TheStrayPredicateCaseIsNotClosedByThisRDR` | DOMAIN EDGE | GREEN (predecessor) |
| 25 | `TestReq10And18_RecognizedProvenanceUnderAnotherNameFailsReservedTagKey` + `TestReq26And29And30And92And101_…` | HAPPY PATH / ADVERSARIAL | RED (payload half) |

### C. Failure payload (block 3) — all NET-NEW

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 26 | `TestReq26And29And30And92And101_BothDirectionsCarryTheThreeFieldPayload` | ADVERSARIAL | RED |
| 27 | `TestReq27And31And33_OneCategoryCarriesTwoDistinctRuleIdentifiers` | BOUNDARY | RED |
| 28 | `TestReq28And91And93_PayloadSurvivesAConsumerThatDoesNotKnowTheCategory` | ADVERSARIAL | RED |
| 29 | `TestReq26And29And30And92And101_BothDirectionsCarryTheThreeFieldPayload` | ADVERSARIAL | RED |
| 30 | `TestReq26And29And30And92And101_BothDirectionsCarryTheThreeFieldPayload` | ADVERSARIAL | RED |
| 31 | `TestReq27And31And33_OneCategoryCarriesTwoDistinctRuleIdentifiers` | BOUNDARY | RED |
| 32 | `TestReq32_TheCategoryTokenIsExactlyReservedTagKey` | BOUNDARY | GREEN (predecessor) |
| 33 | `TestReq27And31And33_OneCategoryCarriesTwoDistinctRuleIdentifiers` | BOUNDARY | RED |
| 34 | `TestReq34And35_TheThreeValuesAreDistinctFieldsNotOneConflatedString` | DOMAIN EDGE | RED |
| 35 | `TestReq34And35_TheThreeValuesAreDistinctFieldsNotOneConflatedString` | DOMAIN EDGE | RED |

**Deviation D6 applied**: REQ-26's `unknown tag` arm is NOT asserted. RDR 0002
as implemented trips `CatMalformedModelDeclaration` (`load.go:169-171`) before
any `unknown tag` failure can carry the reserved key, so the arm is
unreachable. The payload binds `reserved_tag_key` failures only. Verified
empirically: `neg/neg-no-recognized-decl.toml` refuses
`malformed_model_declaration`, pinned by `TestReq15And16_…`.

### D. Input producer obligation and the exported predicate (block 4) — NET-NEW

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 36 | `TestReq36And38_ReservedKeyedProducerTagBreachesAtBothCallSites` | ADVERSARIAL | RED |
| 37 | `TestReq37_BreachIsAGoErrorAndMintsNoRefusalKind` | DOMAIN EDGE | RED |
| 38 | `TestReq36And38_ReservedKeyedProducerTagBreachesAtBothCallSites` | ADVERSARIAL | RED |
| 39 | `TestReq39_OnePredicateCarriesAllThreeReservedKeyChannels` | BOUNDARY | RED |
| 40 | `TestReq40And69And70_ExactlyOneNewExportedInputPredicateNotNamedFinal` | ADVERSARIAL | GREEN (ceiling; see table above) |
| 41 | `TestReq41_PredicateReadDomainIsExactlyTheThreeSequences` | BOUNDARY | RED |
| 42 | `TestReq42_DuplicateReservedKeyInOneSequenceIsAdmissibleAndBreaches` | INPUT EDGE | RED |
| 43 | `TestReq43_MultipleBreachesYieldOneErrorNotAnAggregate` | ADVERSARIAL | RED |
| 44 | `TestReq44And45And90_ReservedKeyBreachIsNotSkippedByACoincidentShapeBreach` | ADVERSARIAL | RED |
| 45 | `TestReq44And45And90_ReservedKeyBreachIsNotSkippedByACoincidentShapeBreach` | ADVERSARIAL | RED |
| 46 | `TestReq46And89_ReservedKeyBreachesEvenWithNoRecognizedOutcome` | BOUNDARY | RED |

REQ-45 and REQ-90 deliberately do **not** assert which of two coincident errors
is reported — JDR 0001 §JD-5 is open and REQ-45 forbids encoding an order.
Per ASSUMPTION-5 the doubly-breaching table is hand-constructed on today's
`resolve.Row` (non-empty `Escape` + non-empty `Writes`); RDR 0009 is
unimplemented, so no second check co-fires and none is asserted.

### E. `RequiresOwned` name reservation (block 5)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 47 | `TestReq47And50And97_RequiresOwnedNamingTheReservedKeyBreachesAtBothSites` | ADVERSARIAL | RED |
| 48 | `TestReq48And51And96_OwnedOnlyResidualRemainsReachableOnTheInternalPath` | DOMAIN EDGE | GREEN (predecessor) |
| 49 | `TestReq49And52And98_RequiresOwnedBreachIsNotADataLevelCategory` | DOMAIN EDGE | RED |
| 50 | `TestReq47And50And97_RequiresOwnedNamingTheReservedKeyBreachesAtBothSites` | ADVERSARIAL | RED |
| 51 | `TestReq48And51And96_OwnedOnlyResidualRemainsReachableOnTheInternalPath` | DOMAIN EDGE | GREEN (predecessor) |
| 52 | `TestReq49And52And98_RequiresOwnedBreachIsNotADataLevelCategory` | DOMAIN EDGE | RED |

### F. Unchanged kernel disposition (block 6)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 53 | `TestReq53_NoNewRefusalKindAndTheClosedSetIsUnchanged` | BOUNDARY | GREEN (predecessor) |
| 54 | `TestReq54And94_D3PrecedenceDeterministicallyResolvesTheCollision` | ADVERSARIAL | GREEN (predecessor) |
| 55 | `TestReq55And67And71_ResolveSignatureIsUnchangedAndBreachYieldsNoDisposition` | BOUNDARY | RED |

### G. Identity, near-spellings, and the near-miss advisory (LBD)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 56 | `TestReq19And56And57And79_NearSpellingsAreOrdinaryUnreservedNames` | INPUT EDGE | GREEN (predecessor) |
| 57 | `TestReq19And56And57And79_NearSpellingsAreOrdinaryUnreservedNames` | INPUT EDGE | GREEN (predecessor) |
| 58 | `TestReq58And59And60And80_DisjunctiveNearMissTriggerAndTheTargetedControl` | ADVERSARIAL | RED |
| 59 | `TestReq58And59And60And80_DisjunctiveNearMissTriggerAndTheTargetedControl` | ADVERSARIAL | RED |
| 60 | `TestReq58And59And60And80_DisjunctiveNearMissTriggerAndTheTargetedControl` | ADVERSARIAL | RED |
| 61 | `TestReq61And62And63And81_AdvisoryIsNonBlockingCategorylessAndByteExact` | BOUNDARY | RED |
| 62 | `TestReq61And62And63And81_AdvisoryIsNonBlockingCategorylessAndByteExact` | BOUNDARY | RED |
| 63 | `TestReq61And62And63And81_AdvisoryIsNonBlockingCategorylessAndByteExact` | BOUNDARY | RED |
| 64 | `TestReq64_TheAdvisoryChannelIsSeparateFromTheFailureChannel` | DOMAIN EDGE | RED |
| 65 | `TestReq65_CanonicalNameIsTheBareLiteralWithNoSigil` | BOUNDARY | GREEN (predecessor) |

**Deviation D5 applied**: the advisory is asserted at the `internal/table`
boundary as a list returned **separately from the error**
(`LoadWithAdvisories`). CLI delivery is out of 0008's scope and is not
asserted; REQ-102 pins only that §D10's `Finding.Hint` carrier exists.

### H. Enforcement locus (LBD)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 66 | `TestReq66And68_TheTwoCallSitesAgreeOnEveryInput` | ADVERSARIAL | GREEN (agreement; see table above) |
| 67 | `TestReq55And67And71_ResolveSignatureIsUnchangedAndBreachYieldsNoDisposition` | BOUNDARY | RED |
| 68 | `TestReq66And68_TheTwoCallSitesAgreeOnEveryInput` | ADVERSARIAL | GREEN (agreement; see table above) |
| 69 | `TestReq40And69And70_ExactlyOneNewExportedInputPredicateNotNamedFinal` | ADVERSARIAL | GREEN (ceiling; see table above) |
| 70 | `TestReq40And69And70_ExactlyOneNewExportedInputPredicateNotNamedFinal` | ADVERSARIAL | GREEN (ceiling; see table above) |
| 71 | `TestReq55And67And71_ResolveSignatureIsUnchangedAndBreachYieldsNoDisposition` | BOUNDARY | RED |

### I. Minimum Viable Validation

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 72 (kernel half) | `TestMVV0008_ReservedKeyOwnershipEndToEnd` legs 3–6 | HAPPY PATH | RED (leg 5) |
| 73 (normalizer half) | `TestMVV0008_ReservedKeyOwnershipEndToEnd` legs 1–2; `TestReq73And75_…` | HAPPY PATH | RED (payload half) |

**REQ-MVV is a runnable end-to-end test**:
`internal/table/mvv_0008_test.go::TestMVV0008_ReservedKeyOwnershipEndToEnd`.
The RDR declares **no** Round-Trip / Inverse Invariant (no encode/decode,
import/export, or inverse operation is introduced), so the MVV's obligation is
value-for-value assertion of each half's stated outcome rather than
reconstruct-and-compare. Every leg compares values — plan `RuleID`, `Revision`,
`NextTags`, `Writes`, `Escaped`; the three payload fields byte-for-byte; and
the breach legs assert `Result{}` equality rather than "did not error". A green
exit code is explicitly not sufficient and is nowhere relied on.

### J. Testing Strategy scenarios

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 74 (TS-1) | `TestMVV0008_…` legs 3–4; `TestReq1_AssembledViewBinds…` | HAPPY PATH | GREEN (predecessor) |
| 75 (TS-2) | `TestReq73And75_NormalizerHalfCleanLoadPlusBothFailuresWithPayload` | HAPPY PATH | RED (payload half) |
| 76 (TS-3) | `TestReq76And77And78_GuardSeamAndMatcherReadTheSameRecognizedBinding` | ADVERSARIAL | GREEN (predecessor) |
| 77 (TS-3) | `TestReq76And77And78_…` | ADVERSARIAL | GREEN (predecessor) |
| 78 (TS-3) | `TestReq76And77And78_…` | ADVERSARIAL | GREEN (predecessor) |
| 79 (TS-4) | `TestReq19And56And57And79_…` + `TestReq20And79_…` | INPUT EDGE / ADVERSARIAL | GREEN (predecessor) |
| 80 (TS-4) | `TestReq58And59And60And80_…` | ADVERSARIAL | RED |
| 81 (TS-4) | `TestReq61And62And63And81_…` | BOUNDARY | RED |
| 82 (TS-5) | `TestReq12And82And83_…` | DOMAIN EDGE | GREEN (predecessor) |
| 83 (TS-5) | `TestReq12And82And83_…` | DOMAIN EDGE | GREEN (predecessor) |
| 84 (TS-5) | `TestReq13And84_…` | BOUNDARY | GREEN (predecessor) |
| 85 (TS-6) | `TestReq85And86And87_AllThreeTS6VariantsBreachAtBothSitesWithZeroResult` | ADVERSARIAL | RED |
| 86 (TS-6) | `TestReq85And86And87_…` | ADVERSARIAL | RED |
| 87 (TS-6) | `TestReq85And86And87_…` | ADVERSARIAL | RED |
| 88 (TS-6) | `TestReq88And99_ConformingInputsIncludingTheEmptyTupleStayNilError` | HAPPY PATH | GREEN (no-regression half; see table above) |
| 89 (TS-6) | `TestReq46And89_…` | BOUNDARY | RED |
| 90 (TS-6) | `TestReq44And45And90_…` | ADVERSARIAL | RED |
| 91 (TS-7) | `TestReq28And91And93_…` | ADVERSARIAL | RED |
| 92 (TS-7) | `TestReq26And29And30And92And101_…` | ADVERSARIAL | RED |
| 93 (TS-7) | `TestReq28And91And93_…` | ADVERSARIAL | RED |
| 94 (TS-8) | `TestReq54And94_…` | ADVERSARIAL | GREEN (predecessor) |
| 95 (TS-8) | `TestReq95_TheCollisionPathStaysConstructibleInPackage` | DOMAIN EDGE | RED |
| 96 (TS-9) | `TestReq48And51And96_…` | DOMAIN EDGE | GREEN (predecessor) |
| 97 (TS-9) | `TestReq47And50And97_…` | ADVERSARIAL | RED |
| 98 (TS-9) | `TestReq49And52And98_…` | DOMAIN EDGE | RED |

**Deviation D1 applied and DISCHARGED — no SPEC-DEFECT.** Scenario 3 is written
against RDR 0007's per-atom guard seam (`valueSeam` in
`internal/resolve/guard_fixtures_test.go`, a NEW seam alongside `fixtureGuards`,
not a change to it — REQ-76's own requirement). The test captures the `value`
handed to the guard atom over `recognized` and asserts it equals
`in.Recognized`, while a `Match` on the same key selects the row in the SAME
resolve; provenance (`ProvenanceRecognized`) is read on the package-internal
path over the identical `Input`, since `0007:C1` fences `Evaluate` from the
view. **The same-view property IS assertable at that seam** — the test passes
against the shipped kernel — so D1's escalation branch does not fire.

### K. Phase 3 conformance set

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 99 | `TestReq88And99_…` | HAPPY PATH | GREEN (no-regression half; see table above) |
| 100 | `TestReq19And56And57And79_…` + `TestReq61And62And63And81_…` | INPUT EDGE / BOUNDARY | RED (advisory half) |
| 101 | `TestReq26And29And30And92And101_…` | ADVERSARIAL | RED |

### L. §JD-8 / §JD-9 citations (deviation D3)

| REQ | Test | Label | Colour |
| --- | --- | --- | --- |
| 102 | `TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint` | BOUNDARY | GREEN (predecessor) |
| 103 | `TestReq103_ReservedTagKeyMapsToTheSharedModelInvalidCodeAtExit2` | DOMAIN EDGE | GREEN (predecessor) |
| 104 | `TestReq104_TagOnTheReservedKeyIsAGroupUserEnvRefusalAtTheCLI` | ADVERSARIAL | GREEN (predecessor) |
| 105 | `TestReq105_TagValuesStayObservedAndNeverReachTheRecognizedChannel` | BOUNDARY | GREEN (predecessor) |
| 106 | `TestReq106_NoAccessorKeysEntryMayNameTheReservedKey` | ADVERSARIAL | GREEN (predecessor) |

D3's check also requires the RECORD cite §D10 and §D8 in a body clause. That is
a documentation edit on `docs/rdr/0008-*.md`, not a test, and is **not** in this
Phase's scope — flagged below as a REQ-side orphan.

## Orphans

### REQs with no test

| REQ | Reason |
| --- | --- |
| — (D3's record-citation half) | `grep -n '§D10\|§D8\b' docs/rdr/0008-*.md` hits only the Status line. Adding the body citation is a record edit, and this Phase authors tests only. Carried forward. |

Every REQ-1 … REQ-106 has at least one test. No REQ is unmapped.

### Tests with no REQ

None. Every test in the eight files above opens with at least one `// REQ-N:`
quote. Fixture builders and structural probes
(`reserved_key_fixtures_0008_test.go`, `export_0008_test.go`) carry no `Test`
functions and are therefore not orphan tests.

## Notes on negative REQs

Nine REQs are NEGATIVE — they forbid something rather than requiring it, and
each is nonetheless given an executable witness so over-implementation fails:

| REQ | Forbidden | Witness |
| --- | --- | --- |
| 5 | closing the empty-string-outcome path | `TestReq4And5_…` asserts no new refusal/error appears |
| 12 | a separate cardinality check | `TestReq12And82And83_…` asserts the category is the naming one and the failure is singular |
| 15 | re-homing the lower bound into `reserved_tag_key` | `TestReq15And16_…` asserts it stays `malformed_model_declaration` |
| 22, 23 | a predicate-position or write-position reserved-key check | `TestReq21And22And23_…` asserts those positions keep 0002's categories |
| 24 | closing the stray-predicate case | `TestReq24_…` asserts near-miss-named declarations still load clean |
| 40, 70 | binding the predicate to RDR 0009's `Final` | `TestReq40And69And70_…` |
| 48 | restating/altering `missingOwned`'s owned-only test | `TestReq48And51And96_…` asserts the residual still fires |
| 49, 52, 98 | a `reserved_tag_key` failure for a `RequiresOwned` entry | `TestReq49And52And98_…` |
| 65 | a sigil-guarded canonical name | `TestReq65_…` |
| 95 | testing A5's unreachability half | `TestReq95_…` asserts only that the path stays constructible |

---

## Stage 8 Phase 2 (implementation) — REQ-MVV actual output, 2026-08-26

All 23 net-new tests are green. `go test ./...`, `go vet ./...`, and
`golangci-lint run` (0 issues) are clean across the repository; no
predecessor test was weakened, and no unrelated test changed.

`go test ./internal/table/ -run TestMVV0008_ReservedKeyOwnershipEndToEnd -v`:

```
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/1_normalizer_half_clean_load_under_the_reserved_name
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload/declaration_renamed_away_from_the_reserved_key
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload/owned_declaration_taking_the_reserved_key
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/3_kernel_half_a_row_matching_the_reserved_key_fires
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/4_kernel_half_the_recognized_outcome_is_readable_at_the_reserved_key
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/owned_tag_keyed_recognized
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/observed_tag_keyed_recognized
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/row_naming_it_in_RequiresOwned
=== RUN   TestMVV0008_ReservedKeyOwnershipEndToEnd/6_the_two_halves_agree_on_the_spelling
--- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/1_normalizer_half_clean_load_under_the_reserved_name (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload (0.00s)
        --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload/declaration_renamed_away_from_the_reserved_key (0.00s)
        --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/2_normalizer_half_both_directions_fail_with_their_payload/owned_declaration_taking_the_reserved_key (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/3_kernel_half_a_row_matching_the_reserved_key_fires (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/4_kernel_half_the_recognized_outcome_is_readable_at_the_reserved_key (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites (0.00s)
        --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/owned_tag_keyed_recognized (0.00s)
        --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/observed_tag_keyed_recognized (0.00s)
        --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/5_kernel_half_all_three_breach_channels_are_rejected_at_both_sites/row_naming_it_in_RequiresOwned (0.00s)
    --- PASS: TestMVV0008_ReservedKeyOwnershipEndToEnd/6_the_two_halves_agree_on_the_spelling (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/table	0.198s
```

Both MVV halves are executed: the normalizer half (clean load plus both
failure directions with their payloads) and the kernel half (the reserved-key
match fires; all three breach channels are rejected at both call sites).

### Implementation sites

| Surface | Site |
| --- | --- |
| `resolve.CheckInput` + `Resolve` entry call | `internal/resolve/precondition.go`, `internal/resolve/resolve.go` |
| `Failure.Offending` / `.Remedy` / `.Rule` population | `internal/table/load.go` (`loadTags`) |
| `table.Advisory` / `LoadWithAdvisories` near-miss scan | `internal/table/advisory.go` |
| REQ-9 doc pointer | `internal/resolve/resolve.go::recognizedTagKey` |

The D3 record-citation orphan above is unchanged by this Phase: it remains a
documentation edit on `docs/rdr/0008-*.md`, still owed.
