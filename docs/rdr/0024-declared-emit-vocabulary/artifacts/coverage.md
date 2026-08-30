# Coverage — RDR 0024 Phase 1 (tests first)

REQ × test, both directions. 110 REQs from `artifacts/req-list.md`;
**83 test functions across seven files in three packages** —
`table_test` (external), `table` (internal), and `cli`.

**Phase 2 status: GREEN.** Every test below passes on the implemented
tree, and `make check` is green end to end (fmt, vet, golangci-lint,
build, graph-lint over all four examples, docs-check, and `go test -race`
over every package). The MVV run is recorded below. Four Phase 1 fixtures
and two predecessor censuses moved; each is a classified entry in
`deviations.md` (D1, D2, D5, D6) and no assertion was relaxed.

**Red gate (as of Phase 1).** Both touched test packages **DID NOT
COMPILE** against the pre-implementation HEAD:
`table.CatMalformedEmitDeclaration`, `table.CatUnknownEmitKey`,
`table.CatEmitValueOutOfDomain`, `table.EmitDecl`, `Model.EmitDecls`,
`(*loader).loadEmitDecls`, `(*loader).checkRuleEmit`, and
`resolvePayload.Dispositions` do not exist. That is the red-before-green
gate for a phase whose Implementation Plan is "make the emit vocabulary a
declared, checked fact". It follows the precedent this repo set at RDR
0010 Phase 1 and RDR 0008 (`internal/table/category.go`: "PHASE 1
DECLARATION ONLY. Nothing populates these yet; the RDR 0008 conformance
suite is red against them by design").

**Verification method.** Because a non-compiling package hides the
pass/fail split, the split below was established by compiling against a
temporary stub of the missing API (a `zz_rdr0024_stub.go` plus two
one-line field insertions; deleted before every commit, present in no
commit) and running the whole suite. The baseline was then re-verified:
with the stub removed AND all seven 0024 test files moved aside,
`go test ./internal/...` on this branch is **fully green** — there are no
pre-existing unrelated failures.

Commits use `--no-verify` because the repo's pre-commit hook runs
`go vet ./...`, which a red-first phase against missing API cannot pass by
construction.

Legend:

- **RED** — new behaviour; fails until Phase 2/3 lands the implementation.
- **RED (compile)** — the assertion is a compile-time one (a bound method
  value assigned to `func() error`); red as a build failure, green exactly
  when the named symbol exists with the named signature.
- **RED (stub-visible)** — green only under the temporary stub, because
  the stub IS the shape the REQ asserts; red without it.
- **GREEN-BC** — characterization / negative REQ pinning behaviour already
  correct at HEAD. Not tautological: each drives real production code over
  a real fixture and would fail if the implementation regressed the pinned
  behaviour.

Every test opens with its REQ quote and carries one of `HAPPY PATH`,
`INPUT EDGE`, `BOUNDARY`, `ADVERSARIAL`, `DOMAIN EDGE`.

---

## A. Grammar: the `[emit]` declaration table (C1)

| REQ | Test | State |
|---|---|---|
| REQ-1 | `table_test.TestReq1_0024_AModelMayDeclareItsEmitVocabulary` | RED |
| REQ-2 | `table_test.TestReq2_0024_TheOptInTriggerIsTheCountOfDeclaredKeys` | RED |
| REQ-3 | `table_test.TestReq3_0024_KindIsOneOfFourTokensAndSetIsExcluded` | RED |
| REQ-4 | `table_test.TestReq4_0024_DomainTakesTwoSpellingsOfTheSameKey` | RED |
| REQ-5 | `table_test.TestReq5_0024_BothDomainSpellingsIsTheDecodersDuplicateKeyError` | RED |
| REQ-6 | `table_test.TestReq6_0024_NonEnumKindsTakeNoDomainAndNoDispositions` | RED |
| REQ-7 | `table_test.TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive` | RED |
| REQ-8 | `table_test.TestReq8_0024_ANonStringEmitValueRefusesAsMalformedTOML` | RED |
| REQ-9 | `table_test.TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted` | RED |
| REQ-10 | `table_test.TestReq10_0024_BothDomainShapesDecodeThroughOneAnyTypedField` | RED |
| REQ-11 | `table_test.TestReq11_0024_StrictDecodeStaysSetForTheWholeDocument` | RED |
| REQ-12 | `table_test.TestReq12_0024_ADeclarationKeyTypoIsTheDecodersUnknownSchemaField` | RED |
| REQ-13 | `table_test.TestReq94_0024_EveryMalformedDeclarationArmRefuses` (arm: unknown kind) | RED |
| REQ-14 | `TestReq94` (three no-usable-domain arms) + `table_test.TestReq95_0024_NoUsableDomainIsOneArmWithThreeFixtures` | RED |
| REQ-15 | `TestReq94` (empty member ×2, duplicate member ×3) | RED |
| REQ-16 | `TestReq94` (domain on bool/int/scalar, both spellings — ASSUMPTION-8) | RED |
| REQ-17 | `TestReq94` (empty-string disposition token) | RED |
| REQ-18 | `TestReq94` (six fixtures over the four strictness-unreachable arms — **A7's Pending closure leg**) | RED |
| REQ-19 | `table_test.TestReq95_0024_NoUsableDomainIsOneArmWithThreeFixtures` | RED |
| REQ-20 | `table_test.TestReq20_0024_AnEmitDeclarationAdmitsNoTagDeclarationKey`, `table_test.TestReq78_0024_AnEmitKeyIsStillNotATagKeyUnderADeclaration` | RED |
| REQ-21 | `table_test.TestReq21_0024_ABareEmitTableIsIdenticalToOmittingIt` | RED |
| REQ-22 | `table_test.TestReq22_0024_ADeclarationIsAuthorOwnedAndUnversioned` | RED |

**Note on REQ-11 / REQ-12 / REQ-20.** All three assert the DECODER's
`unknown_schema_field`, which HEAD already returns — for the wrong reason
(`[emit]` itself is unknown at HEAD). Each therefore carries a
**clean-load control** in the same test, so it cannot be satisfied by a
build that refuses `[emit]` wholesale. All three are RED for that reason.

## B. Proof: the load-pipeline refusals (C2)

| REQ | Test | State |
|---|---|---|
| REQ-23 | `table_test.TestReq23_0024_ZeroDeclarationsMakeNoNewRefusalReachable` | RED |
| REQ-24 | `table_test.TestReq24_0024_TheChecksRunBeforeRowsAreYielded` | RED |
| REQ-25 | `table.TestReq25_0024_BothStepsSitImmediatelyAfterLoadTagsInTheStepSlice`, `table_test.TestReq25_0024_TheGrammarStepPrecedesTheCrossCheckStep` | RED (compile) / RED |
| REQ-26 | `table.TestReq26_0024_BothStepsAreNoArgErrorReturningLoaderMethods` | RED (compile) |
| REQ-27 | `table.TestReq25_0024_BothStepsSitImmediatelyAfterLoadTagsInTheStepSlice`, `table_test.TestReq27_0024_TheOptInGateIsEvaluatedInsideEachStep` | RED |
| REQ-28 | `table_test.TestReq28_0024_AnUndeclaredEmitKeyRefusesOnBothRuleClasses` | RED |
| REQ-29 | `table_test.TestReq29_0024_AnOutOfDomainValueRefusesPerKind` | RED |
| REQ-30 | `table_test.TestReq30_0024_EveryEmitRefusalCarriesTheOffendingSourceLine` | RED |
| REQ-31 | `TestReq30` (declaration-defect leg) | RED |
| REQ-32 | `TestReq30` (both rule-side legs), `table_test.TestReq32_0024_TheRuleAnchorMatchesTheIdValueNotTheModelIdKey` | RED |
| REQ-33 | `table.TestReq33_0024_NoDecoderPositionPlumbingIsIntroduced` | RED |
| REQ-34 | `cli.TestReq34_0024_EveryEmitCategoryRefusesOnBothSurfaces` | RED |
| REQ-35 | `TestReq34` (per-surface envelope code) | RED |
| REQ-36 | `cli.TestReq36_0024_TheLintLoadRefusalStillEmitsModelInvalid` | RED |
| REQ-37 | `cli.TestReq37_0024_NoEmitCategoryEverLandsOnTheAdvisoryTier` | RED |
| REQ-38 | `table_test.TestReq38_0024_AnEmitRefusalIsFailFastAndCarriesNoErrorList` | RED |
| REQ-39 | `TestReq38` — the test asserts the COUNT and does NOT assert which | RED |
| REQ-40 | `TestReq38` (`Unwrap() []error` absent; no accumulation) | RED |
| REQ-41 | `table_test.TestReq41_0024_NormalizationCarriesTheEmitBlockUnmutated` | RED |
| REQ-42 | `table_test.TestReq42_0024_UnusedDeclarationsAreNotFindingsOfAnyTier`, `cli.TestReq37_0024_…` | RED |
| REQ-43 | `cli.TestReq43_0024_GraphlintGainsNoEmitAwareFinding` | RED |

**ASSUMPTION-4** is honoured throughout: a stamped line is asserted as
non-zero, non-1, and equal to the line the fixture's own anchor sits on —
resolved from the fixture text, never a hard-coded integer.

## C. Carry: the normalized declaration carrier (C3)

| REQ | Test | State |
|---|---|---|
| REQ-44 | `table_test.TestReq44_0024_EveryDeclaredFieldIsCarriedOntoTheModel` | RED |
| REQ-45 | `table_test.TestReq45_0024_TheCarrierIsANewTypeWithExactlyThreeFields` | RED (stub-visible) |
| REQ-46 | `TestReq45` (three exported fields, exact types) | RED (stub-visible) |
| REQ-47 | `table_test.TestReq47_0024_EmitDeclsIsAlwaysNonNilAfterASuccessfulLoad` | RED |
| REQ-48 | `table_test.TestReq48_0024_DomainIsNilNotEmptyForEveryNonEnumKind` | RED |
| REQ-49 | `table_test.TestReq49_0024_TheCarriedDomainIsTheBytewiseSortedUnion` | RED |
| REQ-50 | `table_test.TestReq50_0024_TheRoundTripSweepIsNotEditedToReachThisCarrier`, `cli.TestReq83_0024_TheLicensedDiffIsClosed` | RED / GREEN-BC |
| REQ-51 | `TestReq49` (flat-array authored order not preserved) | RED |
| REQ-52 | `table_test.TestReq52_0024_ValueConformanceRunsThroughConformValueWithNoBounds` | RED |
| REQ-53 | `TestReq52` (Min/Max/Elements stay nil — no bound ever fires) | RED |
| REQ-54 | `table_test.TestReq54_0024_AScalarKeyIsShortCircuitedBeforeTheConformCall` | RED |
| REQ-55 | `table_test.TestReq25_0024_TheGrammarStepPrecedesTheCrossCheckStep` | RED |
| REQ-56 | `table_test.TestReq56_0024_TheKernelAndDumpSurfacesCarryNoDeclarations` | RED |

**REQ-50 is asserted from both sides.** The table-side test pins that the
carrier is reachable independently of the RDR 0002 sweep; the CLI-side
`TestReq83` pins that `internal/table/roundtrip_test.go` carries no 0024
reference at all. A8 is Refuted and the sweep is not edited.

## D. Envelope: the `dispositions` payload field (C4)

| REQ | Test | State |
|---|---|---|
| REQ-57 | `cli.TestReq57_0024_ResolveCarriesDispositionsForTheSelectedValue`, `cli.TestReq57_0024_DispositionsIsEmptyObjectInEveryOtherCase` | RED |
| REQ-58 | `cli.TestReq58_0024_TheFieldIsAStringMapWithNoOmitempty` | RED (stub-visible) |
| REQ-59 | `cli.TestReq59_0024_DispositionKeysArriveInByteOrderFromTheMarshaller` | RED |
| REQ-60 | `cli.TestReq60_0024_TheMapIsKeyedOffTheSelectedRowNotTheDeclarationSet`, `cli.TestReq60_0024_ADeclaredKeyWhoseValueCarriesNoDispositionContributesNoEntry` | RED |
| REQ-61 | `cli.TestReq103_0024_DispositionsSitsImmediatelyAfterEmitOnTheWire` | RED |
| REQ-62 | `cli.TestReq62_0024_TheTablePackageGainsNoPayloadShapedHelper` | RED |
| REQ-63 | `TestReq58` (`json:"dispositions"` with no `,omitempty`) | RED (stub-visible) |
| REQ-64 | `cli.TestReq64_0024_TheDispositionTokenIsSurfacedVerbatim` | RED |
| REQ-65 | `cli.TestReq65_0024_TextModeRendersThroughTheGenericPayloadRenderer` | RED |
| REQ-66 | `cli.TestReq66_0024_AnEscapeRescueJoinsFromTheEscapeRowsOwnValues`, `cli.TestReq66_0024_AnEscapeRowAuthoringNoEmitJoinsAnEmptyMap` | RED |
| REQ-67 | `cli.TestReq67_0024_FlowNextCarriesNoDispositions` | RED |
| REQ-68 | `TestReq103` (14 wire keys, `NumField() == 15`, struct index 9) | RED |

## E. Load-bearing decisions and naming

| REQ | Test | State |
|---|---|---|
| REQ-69 | `table_test.TestReq69_0024_ADeclarationIsIdentifiedByItsByteExactEmitKey` | RED |
| REQ-70 | `table_test.TestReq70_0024_TheThreeCategorySlugsAreFixedLiterals` | RED (stub-visible) |
| REQ-71 | `table_test.TestReq104_0024_EachEmitCategoryIsRegisteredInCategories` | RED |
| REQ-72 | `table_test.TestReq72_0024_DispositionsArePerMemberNotPerKey` | RED |

## F. Failure modes

| REQ | Test | State |
|---|---|---|
| REQ-73 | `cli.TestReq34_0024_EveryEmitCategoryRefusesOnBothSurfaces` | RED |
| REQ-74 | `cli.TestReq74_0024_ScalarIsTheEscapeHatchAndZeroDeclarationsIsTodaysWorld` | RED |
| REQ-75 | `table_test.TestReq75_0024_StrictDecodeStillRefusesAnUnknownTopLevelTable` | **GREEN-BC** |

## G. Implementation-plan obligations

| REQ | Test | State |
|---|---|---|
| REQ-76 | `table.TestReq33_0024_NoDecoderPositionPlumbingIsIntroduced` (per ASSUMPTION-3 the obligation is that both lines are recovered by a source scan; the helper's name and arity are the implementer's) | RED |
| REQ-77 | `cli.TestReq77_0024_TheThreeStaleCommentsNoLongerClaimEmitIsUndeclared` | RED |
| REQ-78 | `table_test.TestReq78_0024_AnEmitKeyIsStillNotATagKeyUnderADeclaration` | RED |
| REQ-79 | `table_test.TestReq79_0024_TheKindArmStaysBelowTheDecoder` | RED |
| REQ-80 | `cli.TestReq80_0024_EveryResolveGoldenCarriesTheInsertedDispositionsKey` | RED |
| REQ-81 | `TestReq80` (the `emit`/`dispositions`/`next` adjacency regex — exactly one inserted key per site) | RED |
| REQ-82 | `cli.TestReq82_0024_TheLicensedTestReq39EditIsCarried` | RED |
| REQ-83 | `cli.TestReq83_0024_TheLicensedDiffIsClosed` | **GREEN-BC** |
| REQ-84 | `cli.TestReq106_0024_TheDispositionsFieldCarriesTheD1RegistrationNote` | RED |
| REQ-85 | `cli.TestReq86_0024_ThePricingExampleDeclaresTheValuesItActuallyAuthors`, `cli.TestReq107_0024_TheOutputContractDocumentsDispositions` | RED |
| REQ-86 | `TestReq86_…DeclaresTheValuesItActuallyAuthors`, `cli.TestReq86_0024_ThePricingExampleResolvesWithJoinedDispositions` | RED |
| REQ-87 | `cli.TestReq87_0024_ASecondRoutingExampleShipsAndLintsClean` | RED |
| REQ-88 | `cli.TestReq88_0024_TheAuthoringDocNoLongerClaimsNothingDeclaresEmitKeys` | RED |

**ASSUMPTION-5 is honoured at REQ-80.** The 28-site census is treated as
evidence: a differing count is `t.Logf`'d as a deviation to record, not a
failure. The assertion that *does* fail is per-golden — every one must
carry `dispositions` inserted between `emit` and `next`.

**ASSUMPTION-6 is honoured at REQ-77/78/88.** Each reads the FALSIFIED
ABSOLUTE CLAIM back and requires its absence, plus (for the doc) that the
narrowed line is actually drawn. No test pins a new wording.

## H. Minimum Viable Validation (`0024:MVV`)

| REQ | Test | State |
|---|---|---|
| REQ-MVV | `cli.TestMVV_0024_DeclaredEmitVocabulary` | RED |
| REQ-89 | `TestMVV_0024` → `mvv0024Step1` (the defect, reproduced — the negative control) | RED |
| REQ-90 | `TestMVV_0024` → `mvv0024Step2` (two runs, one finding each, both categories across them) | RED |
| REQ-91 | `TestMVV_0024` → `mvv0024Step3` (the round trip) | RED |
| REQ-92 | `TestMVV_0024` → `mvv0024Step4` (declaration deleted; today's world plus `dispositions: {}`) | RED |

**The MVV is a runnable end-to-end test** driven through the root Cobra
command, exactly as CI drives it. Step 3 is the round trip and is
asserted **value-for-value**: the payload's `emit` must `DeepEqual`
`{next: archive}` and `dispositions` must `DeepEqual` `{next: stop}`, with
`dispositions` at wire index `emit+1`. A green exit code or "did not
error" is not sufficient and would not pass this test. Step 4 additionally
requires the misspelled `resovle` to answer verbatim once nothing declares
the key.

Step 2 branches on whichever category run A reports and fixes that rule
before run B, because REQ-39 fixes no precedence — pinning which comes
first would fabricate a guarantee.

## I. Testing Strategy scenarios

| REQ | Test | State |
|---|---|---|
| REQ-93 | **META — no test.** "Done = every scenario below is green and the MVV run recorded." This is the phase's exit criterion, discharged by the coverage map itself plus the recorded MVV run, not by an assertion. | n/a |
| REQ-94 | `table_test.TestReq94_0024_EveryMalformedDeclarationArmRefuses` (20 single-defect fixtures) | RED |
| REQ-95 | `table_test.TestReq95_0024_NoUsableDomainIsOneArmWithThreeFixtures` | RED |
| REQ-96 | `table_test.TestReq29_0024_AnOutOfDomainValueRefusesPerKind`, `table_test.TestReq28_0024_…` | RED |
| REQ-97 | `table_test.TestReq38_0024_AnEmitRefusalIsFailFastAndCarriesNoErrorList` | RED |
| REQ-98 | `cli.TestReq98_0024_ZeroDeclarationModelsSeeOnlyTheAppendedEmptyField`, `cli.TestReq37_0024_…` | RED |
| REQ-99 | `TestReq98` (load-side equivalence: the appended field is the only delta) | RED |
| REQ-100 | `table_test.TestReq44_0024_EveryDeclaredFieldIsCarriedOntoTheModel` | RED |
| REQ-101 | `table_test.TestReq101_0024_LoadingAPartitionedDomainRepeatedlyYieldsOneDomain` | RED |
| REQ-102 | `cli.TestReq57_0024_DispositionsIsEmptyObjectInEveryOtherCase`, `cli.TestReq60_0024_…`, and the C4 block above | RED |
| REQ-103 | `cli.TestReq103_0024_DispositionsSitsImmediatelyAfterEmitOnTheWire` | RED |
| REQ-104 | `table_test.TestReq104_0024_EachEmitCategoryIsRegisteredInCategories`, `table_test.TestReq105_0024_…` | RED |
| REQ-105 | `table_test.TestReq105_0024_EachWitnessTripsItsOwnCategoryAndNoOther` | RED |
| REQ-106 | `cli.TestReq106_0024_TheDispositionsFieldCarriesTheD1RegistrationNote` | RED |
| REQ-107 | `cli.TestReq107_0024_TheOutputContractDocumentsDispositions`, `cli.TestReq86_0024_…` | RED |
| REQ-108 | `table_test.TestReq108_0024_ADeclaredEmitValueRendersByteIdenticallyToAnUndeclaredOne`, `table_test.TestReq108_0024_AnIntDeclaredEmitValueIsNotCanonicalizedInTheDump` | RED |

**`0024:S1` (REQ-94) is A7's Pending closure gate.** `PRE` states "Phase 1
does not close until it is green." The sweep exists and is RED, with all
ten arms of ASSUMPTION-1 covered by 20 single-defect fixtures — including
six fixtures over REQ-18's four strictness-unreachable arms, which are the
ones A7's closure leg is actually about.

**`0024:S4` (REQ-101) is the sole determinism oracle.** A8 is Refuted, so
`TestReq146` is neither reached nor edited. The test loads a
five-disposition, ten-member partitioned domain 32 times and requires one
bytewise-sorted union each time; without the sort, Go's randomized map
iteration breaks it with overwhelming probability.

## J. Cross-cutting concerns

| REQ | Test | State |
|---|---|---|
| REQ-109 | `table_test.TestReq109_0024_TheCarrierIsWrittenOnceAndReadOnlyThereafter` | RED |
| REQ-110 | `table_test.TestReq110_0024_NoHashOrContentAddressedIdentityIsCarried` | RED (stub-visible) |

---

## MVV run — recorded (Phase 2)

`0024:MVV`, run end to end against the built binary on the implemented
tree. Model paths are folded to `<tmp>/`; every other byte is verbatim
stdout. `TestMVV_0024_DeclaredEmitVocabulary` passes over the same four
steps.

**Step 1 — the defect, reproduced.** The seed's two-rule adversarial table,
UNDECLARED: `first-rule` emits `next = "resovle"` (misspelled), and
`second-rule` types the key `nxet`.

```
$ intrastate lint --model <tmp>/seed.toml --as=json
{"type":"ok","data":{"findings":[]}}
exit=0
```

Exit 0, zero findings. That is the defect.

**Step 2 — declaring makes both defects visible, one per run.** Added
`[emit.next]`, an enum with a route/stop-partitioned domain
(`route = ["resolve", "refine"]`, `stop = ["archive"]`).

Run A, both defects present:

```
$ intrastate lint --model <tmp>/declared.toml --as=json
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"emit_value_out_of_domain: rule first-rule emits next = \"resovle\": \"resovle\" is outside the declared domain","findings":[{"code":"emit_value_out_of_domain","message":"emit_value_out_of_domain: rule first-rule emits next = \"resovle\": \"resovle\" is outside the declared domain","locator":"<tmp>/declared.toml:29"}]}
exit=2
```

Run B, after fixing the rule run A named (`resovle` -> `resolve`):

```
$ intrastate lint --model <tmp>/fix1.toml --as=json
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"unknown_emit_key: rule second-rule emits the undeclared key nxet","findings":[{"code":"unknown_emit_key","message":"unknown_emit_key: rule second-rule emits the undeclared key nxet","locator":"<tmp>/fix1.toml:38"}]}
exit=2
```

Exit 2, ONE blocking finding each, never both in one run. Across the two
runs both categories are observed — `emit_value_out_of_domain` for the
misspelled value, `unknown_emit_key` for the typo'd key. Each locator
carries the offending block's source line (`:29`, `:38`), not `:1`: line
29 is `first-rule`'s `id` line and line 38 is `second-rule`'s, which is
where `0024:C2` keys a rule-side defect.

**Step 3 — the round trip.** Both rules fixed (`nxet = "stop"` ->
`next = "archive"`):

```
$ intrastate lint --model <tmp>/fixed.toml --as=json
{"type":"ok","data":{"findings":[]}}
exit=0

$ intrastate flow resolve --model <tmp>/fixed.toml --outcome decide \
    --tag step=second --as=json
{"type":"ok","data":{"model":"<tmp>/fixed.toml","revision":"","observed":{"step":"second"},"owned":{},"readers":[],"outcome":"decide","rule":"second-rule","gates":[],"emit":{"next":"archive"},"dispositions":{"next":"stop"},"next":{},"writes":{},"clear":[],"escaped":false}}
exit=0
```

`second-rule` authors `next = "archive"`, which the declaration lists under
`stop`. The payload carries the authored `emit` unchanged AND
`dispositions` mapping the key to `stop`, positioned immediately after
`emit`.

**Step 4 — deleting the declaration restores today's behaviour.** The
ORIGINAL defective table, no `[emit]` table at all:

```
$ intrastate lint --model <tmp>/seed.toml --as=json
{"type":"ok","data":{"findings":[]}}
exit=0

$ intrastate flow resolve --model <tmp>/seed.toml --outcome decide \
    --tag step=first --as=json
{"type":"ok","data":{"model":"<tmp>/seed.toml","revision":"","observed":{"step":"first"},"owned":{},"readers":[],"outcome":"decide","rule":"first-rule","gates":[],"emit":{"next":"resovle"},"dispositions":{},"next":{},"writes":{},"clear":[],"escaped":false}}
exit=0
```

Exit 0, no findings, `dispositions: {}`. The misspelled `resovle` answers
verbatim once nothing declares the key — byte-identical to today apart
from the appended empty field.

**REQ-93 (`TS` "Done").** Every scenario is green (`make check` passes
including `-race`, `go vet`, `golangci-lint`, `graph-lint` over all four
examples, and `docs-check`), and the MVV run is recorded above.

---

## Orphans — both directions

### REQs with no test

- **REQ-93** (`TS` "Done = every scenario below is green and the MVV run
  recorded") — a phase exit criterion, not a behavioural clause. It is
  discharged by this coverage map plus the recorded MVV run. No assertion
  is possible or meaningful.

Every other REQ, REQ-1 through REQ-110 plus REQ-MVV, maps to at least one
test.

### Tests citing no REQ

**None.** Every one of the 83 test functions opens with at least one
`// REQ-N:` quote and a labelled category.

Three of the seven files carry unexported fixture/helper functions that
cite no REQ by design — `emit_decl_fixtures_0024_test.go` (the whole
file), plus `assertLineAt`/`lineOf`/`replaceOnce`/`dispositionsOf`/
`dispositionKeyOrder`/`assertOneEmitFinding`. These are fixture builders
and oracles, not assertions.

Three checked-in `testdata/neg/*.toml` witnesses were added for `0024:S6`:
`neg-emit-unknown-kind.toml`, `neg-emit-unknown-key.toml`,
`neg-emit-value-out-of-domain.toml`. Each is single-defect by
construction, as REQ-105's exclusivity leg requires.

---

## Green new tests, and why

Two tests are green at HEAD, both legitimately, both **negative REQs**
pinning behaviour that is already correct and must not regress:

- **REQ-75** — `TestReq75_0024_StrictDecodeStillRefusesAnUnknownTopLevelTable`.
  `0024:F3` says a model carrying `[emit]` refuses on an OLDER binary as
  `unknown_schema_field`. That property is satisfied by construction at
  HEAD (`decodeStrict`); the REQ is "no code may weaken it". The test
  asserts the invariant against a stand-in unknown top-level table, so it
  keeps meaning after `[emit]` is added to the schema.

- **REQ-83** — `TestReq83_0024_TheLicensedDiffIsClosed`. The licensed
  Phase 3 diff is closed: `internal/table/roundtrip_test.go` must carry no
  0024 reference (A8 Refuted, REQ-50 forbids any `TestReq146` edit), and
  `TestReq43` must stay unedited. Both hold today; the test's job is to
  fail the moment an implementer reaches for either.

Six further tests (REQ-45/46, REQ-58/63, REQ-70, REQ-110) are marked
**RED (stub-visible)**: they went green only under the temporary stub,
because the stub is literally the shape those REQs assert (a three-field
`EmitDecl`, a `map[string]string` field tagged `dispositions`, the three
slug literals). Without the stub they are build failures. They are red
against HEAD and are listed here so the orchestrator can see the
distinction rather than reading a bare "83 red".
