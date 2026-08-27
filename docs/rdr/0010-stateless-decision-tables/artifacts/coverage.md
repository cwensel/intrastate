# Coverage — RDR 0010 Phase 1 (tests first)

REQ × test, both directions. 110 REQs from
`artifacts/req-list.md`; 81 test functions across four files in three
packages.

**Red gate.** The three touched packages do NOT COMPILE against `main` —
`Model.Class`, `Row.Emit`, `table.EmitValue`, `table.IsDecisionTable`, and
`resolvePayload.Emit` do not exist. That is the red-before-green gate for a
phase whose Implementation Plan is "make the class a declared, checked fact
and let a row carry an answer". It follows RDR 0008's precedent in this repo
(`internal/table/category.go`: "PHASE 1 DECLARATION ONLY. Nothing populates
these yet; the RDR 0008 conformance suite is red against them by design").

**Verification method.** Because a non-compiling package hides the pass/fail
split, the split below was established by compiling against a temporary
five-line stub of the missing API (deleted before every commit; not in any
commit) and running the whole suite. Ten characterization tests were then
re-verified GREEN AGAINST `main` in isolation, with the API-dependent files
moved aside. `go test ./internal/...` on `main` with all 0010 files removed
is fully green: **there are no pre-existing unrelated failures**.

Legend:

- **RED** — new behaviour; fails until Phase 2/3 lands the implementation.
- **GREEN-BC** — characterization / negative REQ pinning behaviour already
  correct on `main`. Verified green against `main` in isolation. Not a
  tautological stub-green: each drives real production code over a real
  fixture and would fail if the implementation regressed the pinned
  behaviour.
- **RED (stub-visible)** — green only under the temporary stub because the
  stub IS the shape the REQ asserts; red without it.

---

## A. Model class declaration (C1)

| REQ | Test | State |
|---|---|---|
| REQ-1 | `table_test.TestReq1_ClassAdmitsExactlyTwoValues` | RED |
| REQ-2 | `table_test.TestReq2_AbsentClassReadsAsStateMachine` | RED |
| REQ-3 | `table_test.TestReq3_UnknownClassIsMalformedModelDeclaration` | RED |
| REQ-4 | `table_test.TestReq4_DecisionTableWithOwnedTagRefusesNamingClassAndOwnedCount`, `TestReq4_OwnedTokenCarriesTheActualCount` | RED |
| REQ-5 | `table_test.TestReq5_ZeroOwnedStateMachineIsNotRefused` | RED (fixture uses `class`) |
| REQ-6 | `table_test.TestReq6_ClassIsDeclaredNeverInferredFromTheOwnedSet` | RED |
| REQ-7 | `table_test.TestReq7_ClassIsAnExportedStringFieldWhoseZeroValueIsStateMachine` | RED (stub-visible) |
| REQ-8 | `table_test.TestReq8_TheFieldIsTheStorageForTheClass` | RED |
| REQ-9 | `table_test.TestReq9_AnUndeclaredTagRefusalPrecedesTheClassDisagreement` | RED |
| REQ-10 | `table_test.TestReq10_DoublyMalformedRefusesOnTheClassNotOnWriterArity` | RED |
| REQ-11 | `table_test.TestReq11_TheClassRefusalWinsEvenWhenTheArityCheckWouldPass` | RED |
| REQ-12 | `table_test.TestReq12_TheAgreementCheckSitsInsideTheLicensedWindow` | RED |
| REQ-13 | `table_test.TestReq13_TheClassIsPerModelAndDeclaredUnlikeRowKind` | RED |

REQ-12 is the one structural (AST) oracle in the suite. C1 states a
POSITION, not a behaviour, and deliberately leaves the step's NAME free
("Any step in that window satisfies this clause"). The test therefore reads
`run`'s step slice and asserts the window, never the name (ASSUMPTION-1).

## B. Class-conditioned rule shape (C2)

| REQ | Test | State |
|---|---|---|
| REQ-14 | `table_test.TestReq15_TheProhibitionsAreCarriedByExistingRefusalArms` (four subtests) | RED |
| REQ-15 | `table_test.TestReq15_TheProhibitionsAreCarriedByExistingRefusalArms` | RED |
| REQ-16 | `table_test.TestReq16_TheNoWriteBlockArmIsConditionedOnClass` | RED |
| REQ-17 | `table_test.TestReq17_TheMatchBlockObligationHoldsInBothClasses` | RED |
| REQ-18 | `table_test.TestReq18_ATerminalOverANonOwnedTagIsNotALoadRefusal`; enforcement arm at `graphlint_test.TestReq65_TheTerminalPredicateArmStaysLiveOverADecisionTable` | RED |

## C. `[rule.emit]` grammar, normalization, dump (C3)

| REQ | Test | State |
|---|---|---|
| REQ-19 | `table_test.TestReq19_EmitIsAdmittedOnBothRuleKindsInBothClasses` | RED |
| REQ-20 | `table_test.TestReq20_ANonStringEmitValueRefusesAsMalformedTOML` (4 cases incl. `[rule.emit.sub]`) | RED |
| REQ-21 | `table_test.TestReq21_NoHandWrittenTypeCheckDiscriminatesTheEmitRefusal` | RED |
| REQ-22 | `table_test.TestReq22_RowCarriesEmitAsANewKeyValueStringType` | RED (stub-visible) |
| REQ-23 | `table_test.TestReq23_DuplicateEmitKeysAreRefusedByTheDecoder` | RED |
| REQ-24 | `table_test.TestReq24_EveryExpandedRowCarriesTheAuthoredEmitBlock` | RED |
| REQ-25 | `table_test.TestReq25_SharingTheEmitSequenceAcrossExpandedRowsIsLicensed` | RED |
| REQ-26 | `table_test.TestReq26_EmitIsNeverMutatedByTheDumpReader` | RED |
| REQ-27 | `table_test.TestReq27_AnEmitKeyIsNotATagKey`, `TestReq27_EmitKeysAndValuesAreCarriedByteForByte` | RED |
| REQ-28 | `table_test.TestReq28_EmitNormalizesKeySortedAndKeysOnLengthNotPresence` (3 subtests) | RED |
| REQ-29 | `table_test.TestReq29_EmitJoinsTheDumpVocabularyAppendedLast` | RED |
| REQ-30 | `table_test.TestReq30_TheEmitCellRendersKeyValuePairsBracketedLikeWrites`, `TestReq30_TheEmitValueIsTheRawAuthoredStringNotRenderValue` | RED |
| REQ-31 | `table_test.TestReq31_AnEmptyEmitRendersTheSameEmptyBracketWritesUses` | RED |
| REQ-32 | `table_test.TestReq32_ADumpListOmittingEmitRefuses` (3 subtests) | RED |
| REQ-33 | `table_test.TestReq33_EmitNeverCrossesTheKernelBoundary` | RED |
| REQ-34 | `cli.TestReq34_TheRuleIDJoinIsSoundOverAnExpandingRule` | RED |
| REQ-35 | `cli.TestReq35_EmitOnTheWireIsKeyOrderedWithHTMLEscapingDisabled` | RED |
| REQ-36 | `cli.TestReq35_EmitOnTheWireIsKeyOrderedWithHTMLEscapingDisabled` (byte-order + determinism subtests) | RED |
| REQ-37 | `graphlint_test.TestReq37_EditingAnEmitBlockDoesNotChangeAFingerprint` | RED |

REQ-30's separator (`; `) is asserted per ASSUMPTION-4 by comparing against
the `writes` column's own rendering rather than a hand-written literal
(REQ-31's test does the comparison explicitly).

## D. Resolve payload (C4)

| REQ | Test | State |
|---|---|---|
| REQ-38 | `cli.TestReq38_ResolveCarriesEmitAsAnObjectNeverNullNeverOmitted` | RED |
| REQ-39 | `cli.TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire` | RED |
| REQ-40 | `cli.TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire` (the full 13-key pre-edit order is the `want`) | RED |
| REQ-41 | `cli.TestReq41_TheStateFieldsKeepTheirShapesAndAreEmpty` | RED |
| REQ-42 | `cli.TestReq42_ASuccessPayloadsGatesIsAlwaysAnAllAllowList` | RED |
| REQ-43 | `cli.TestReq43_EscapedStaysBetweenClearAndEscapeClass` | RED |
| REQ-44 | `cli.TestReq44_AnEscapedPlanCarriesTheEscapeRowsOwnEmit` | RED |
| REQ-45 | `cli.TestReq45_TextModeRendersEmitThroughTheGenericRenderer` | RED |
| REQ-46 | `cli.TestReq46_NoPerVerbTextSpecialCaseIsAddedForEmit` | **GREEN-BC** |
| REQ-47 | `cli.TestReq47_ArtifactIsUnrequiredAndAStrayBindingIsIgnored` | RED |
| REQ-48 | `cli.TestReq48_OutcomeRemainsRequiredOverADecisionTable` | **GREEN-BC** |
| REQ-49 | `cli.TestReq50_FlowNextCarriesNoEmit`; `read-state`/`set-state` untouched, pinned by the shipped 0005 suite | RED |
| REQ-50 | `cli.TestReq50_FlowNextCarriesNoEmit` | RED |
| REQ-51 | `cli.TestReq51_TheInvokedReaderSetIsEmptyOverADecisionTable` | RED |

REQ-46 is an AST oracle over `respond/text.go`, excluding COMMENTS: the file
already says the renderer "emits one path-qualified leaf per line", and a
prose match is not a special case. What the clause forbids is CODE keying on
the field name.

## E. ∅-rooted graph lint (C5)

| REQ | Test | State |
|---|---|---|
| REQ-52 | `graphlint_test.TestReq52_DecisionTableLintRootsAtTheEmptyOwnedStateNode`, `TestReq52_TheLoadedModelCarriesTheClassForLintToRead` | RED |
| REQ-53 | `graphlint_test.TestReq53_TheSeedingPredicateIsClassOrDeclaredRootNeverClassAlone` (both directions) | RED |
| REQ-54 | `graphlint_test.TestReq54_CoverageRunsOverTheAuthoredGuardDimensions` | RED |
| REQ-55 | `graphlint_test.TestReq55_AZeroDimensionDecisionTableGroupIsUnprovable` | RED |
| REQ-56 | `graphlint_test.TestReq56_TheZeroDimensionArmCarriesTheFourthReasonValue` | RED |
| REQ-57 | `graphlint_test.TestReq56_TheZeroDimensionArmCarriesTheFourthReasonValue` ("the emission carries it" subtest: empty reason and a borrowed member both fail) | RED |
| REQ-58 | `graphlint_test.TestReq58_AZeroDimensionGroupTakesOnlyTheUnprovableFinding` | RED |
| REQ-59 | `graphlint_test.TestReq60_TheZeroDimensionArmReturnsRatherThanFallingThrough` | RED |
| REQ-60 | `graphlint_test.TestReq60_TheZeroDimensionArmReturnsRatherThanFallingThrough` | RED |
| REQ-61 | `graphlint_test.TestReq61_TheZeroDimensionArmSuppliesItsOwnMessage` | RED |
| REQ-62 | `graphlint_test.TestReq62_TheZeroDimensionArmBindsTheDecisionTableClassOnly` | **GREEN-BC** |
| REQ-63 | `graphlint_test.TestReq63_TheMachineOnlyInvariantsStaySilentOverADecisionTable` (exact finding-list equality, ASSUMPTION-7) | RED |
| REQ-64 | `graphlint_test.TestReq64_NoFindingNamesTheClass` (6 fixtures × 4 fields) | RED |
| REQ-65 | `graphlint_test.TestReq65_TheTerminalPredicateArmStaysLiveOverADecisionTable` | RED |
| REQ-66 | `graphlint_test.TestReq66_DeclaredTerminalHandlingStaysSilentOverADecisionTable` | RED |
| REQ-67 | `graphlint_test.TestReq67_EveryInvariantIsUnchangedOverAStateMachine` | **GREEN-BC** |
| REQ-68 | `graphlint_test.TestReq68_TheOtherTwoInitialSitesKeepTheBareTest` (4 fixtures) | RED |
| REQ-69 | `graphlint_test.TestReq69_ClassReadersNeverRederiveFromTheOwnedSet` | **GREEN-BC** |

REQ-62 and REQ-69 are the two directions of the false-green C5 exists to
close, and both are green on `main` by construction — REQ-62 because no
zero-dimension arm exists yet, REQ-69 because nothing infers a class today.
They are the guards that make Phase 2's edit non-regressive: REQ-62 fails if
the arm is made unconditional (RISK: 151 zero-dimension groups across 37 of
37 loadable fixtures), REQ-69 fails if a reader derives the class from
`len(owned) == 0`.

## F. Cross-cutting / compatibility

| REQ | Test | State |
|---|---|---|
| REQ-70 | `cli.TestReq70_EveryCheckedInModelLoadsLintsAndResolvesUnchanged` | **GREEN-BC** |
| REQ-71 | `table_test.TestReq71_EveryCheckedInDumpCarryingFixtureNamesEmit` | RED |
| REQ-72 | `graphlint_test.TestReq72_LintIsSinglePassAndDoesNotMutateTheModel` | RED |
| REQ-73 | `cli.TestReq73_TheClassIsNotAVersionAndNoNegotiationIsAdded` | **GREEN-BC** |

REQ-71 covers only the two `[dump]`-carrying fixtures the `table_test`
helpers already read (`rdr-fixture.toml`, `kata-fixture.toml`) by name. The
census is 103 files (A4, ASSUMPTION-8); the remaining 101 are exercised
indirectly — each is loaded by some shipped test, and a `[dump]` list
omitting `emit` refuses at load (REQ-32), so an un-updated fixture reddens
its own test. **Phase 2 must update all 103.**

## G. Failure modes

| REQ | Test | State |
|---|---|---|
| REQ-74 | `cli.TestReq74_AClassOwnedDisagreementSurfacesAsFlowModelInvalid` | RED |
| REQ-75 | `cli.TestReq75_NeitherGuardedSilentModeExitsZeroWithAnEmptyFindingList` (both arms) | RED |
| REQ-76 | `cli.TestReq76_DroppingTheClassRecoversAndNothingPersistentIsWritten` | RED |

## H. Minimum Viable Validation (`0010:MVV`)

All of REQ-77..REQ-85 are discharged by the single runnable end-to-end test
`cli.TestMVV_StatelessDecisionTables`, which drives `ExecuteAndEmit` — the
CI command shape — for every step.

| REQ | Subtest | State |
|---|---|---|
| REQ-77 | `step 1 — the decision table is authorable` | RED |
| REQ-78 | `step 2 — partial lint reports the coverage gap` | RED |
| REQ-79 | `step 3 — complete lint is exactly []` | RED |
| REQ-80 | `step 3 variant — escape row yields only the advisory` | RED |
| REQ-81 | `step 4 — resolve returns rule + emit with no --artifact` | RED |
| REQ-82 | `step 5 — dump renders emit; next carries empty required` | RED (see deviations D2) |
| REQ-83 | `step 6 — negative controls` (3 subtests) | RED (see deviations D1) |
| REQ-84 | `end state — emit round-trips authored to resolved` | RED |
| REQ-85 | `cli.TestReq85_TheCheckedInNavigatorModelIsUntouched` | **GREEN-BC** |

**REQ-84 is the inverse invariant, not an exit code.** For all four cells it
asserts that the `[rule.emit]` block the fixture AUTHORS equals `data.emit`
as returned — first by `reflect.DeepEqual` on the decoded object, then
BYTE-for-byte against the raw wire `emit` object, so a re-encoded or
re-ordered value fails. Its companion subtest drives the shipped 0005
state-machine fixture and asserts it resolves exactly as it did.

## I. Testing Strategy — Done clause and licensed diffs

| REQ | Test | State |
|---|---|---|
| REQ-86 | no dedicated test — the Done clause is the SUITE plus the MVV; discharged by every entry above | — |
| REQ-87 | no dedicated test — a review obligation on the Phase 2/3 diff. The two diffs it will license in Go are recorded as deviations D3 (`clonedRowSliceFields`) and D8 (the 0011 hard-coded payload). | — |
| REQ-88 | no dedicated test — the scope gate for the implementation, not a runtime property. RDR 0011's own `TestReq120_TheProductionDiffIsFlowNextPlusOneTermInFlowExec` is the mechanized approximation, and it will FIRE on this RDR's edits (deviation D7). | — |

REQ-86/87/88 are the only REQs with no test, and each is a REVIEW obligation
by construction. See §orphans below.

## J. Testing Strategy — scenarios

| REQ | Test | State |
|---|---|---|
| REQ-89 (SC-1) | `table_test.TestReq89_LoaderTableOverTheClassKey` (6 cases) | RED |
| REQ-90 (SC-2) | `table_test.TestReq90_DecisionTableRowsNormalizeWithEmptyWriteSurface` + `TestReq16_…` | RED |
| REQ-91 (SC-3) | `table_test.TestReq91_TheEmitNormalizationScenarioAsOneTable` (8 cases) | RED |
| REQ-92 (SC-3) | `table_test.TestReq24_EveryExpandedRowCarriesTheAuthoredEmitBlock` + `cli.TestReq92_ANonFirstExpandedRowResolvesToTheAuthoredEmitBlock` | RED |
| REQ-93 (SC-3) | `table_test.TestReq23_DuplicateEmitKeysAreRefusedByTheDecoder` — the negative REQ's only assertable form: the decoder refuses, so no dedup arm is reachable to test | RED |
| REQ-94 (SC-4) | `graphlint_test.TestReq94_TheGraphLintFixturePairAndTheEscapeVariant` + `TestReq103_TheDecisionTableFixtureSetIsPresentAndDiscriminating` | RED |
| REQ-95 (SC-4) | `graphlint_test.TestReq55_AZeroDimensionDecisionTableGroupIsUnprovable` (positive assertion: exactly one finding, blocking) | RED |
| REQ-96 (SC-4) | `graphlint_test.TestReq96_TheOtherwiseRowIsScopedPerOutcome` | RED |
| REQ-97 (SC-4) | `graphlint_test.TestReq65_TheTerminalPredicateArmStaysLiveOverADecisionTable` | RED |
| REQ-98 (SC-5) | `cli.TestReq38_…` ("an unauthored block is {} and never null" subtest) | RED |
| REQ-99 (SC-5) | `cli.TestReq41_…`, `TestReq45_…`, `TestReq47_…` | RED |
| REQ-100 (SC-6) | no dedicated test — see deviations D6 | — |
| REQ-101 (SC-6) | discharged by writing NO two-toolchain test (negative REQ) | — |

## K. Implementation plan and docs

| REQ | Test | State |
|---|---|---|
| REQ-102 (Phase 1) | the whole `table_test` 0010 set | RED |
| REQ-103 (Phase 2) | `graphlint_test.TestReq103_TheDecisionTableFixtureSetIsPresentAndDiscriminating` | RED |
| REQ-104 (Phase 3) | `cli.TestReq104_EmitIsJoinedByRuleIDAfterSelection` + `TestReq38_…` | **GREEN-BC** (the `rowByID` half) |
| REQ-105 (Phase 4) | the fixture corpus IS the deliverable (ASSUMPTION-9): `table_test` fixtures + `graphlint_test` fixtures + `cli.dtModel0010` | RED |
| REQ-106 (Phase 4) | no test — a DOC obligation on `docs/`, discharged by Phase 4's edit, not by a Go test | — |

## L. Naming and rejected shapes

| REQ | Test | State |
|---|---|---|
| REQ-107 | `cli.TestReq107_TheRejectedNamingSpellingsAreUnauthorable` (6 rejected spellings) | **GREEN-BC** |
| REQ-108 | `cli.TestReq108_TwoMatchingRulesIsAmbiguousMatchInTheDecisionTableClass` | RED |
| REQ-109 | `cli.TestReq109_TheBrieflyRejectedShapesAreAbsent` (BR1, BR2, BR3, BR5, BR6) | RED |
| REQ-110 | `table_test.TestReq110_NoFixtureIsLiftedFromTheIllustrativeCode` | RED |

`0010:BR4` (`emit` in the `flow next` candidate preview) is covered by
`cli.TestReq50_FlowNextCarriesNoEmit` rather than by the BR5 test, since
REQ-50 states it directly.

---

## Orphans

### REQs with no test (5, each by construction)

| REQ | Why |
|---|---|
| REQ-86 | "Done = the MVV passes end to end under `make check`" — the whole suite is the test. Writing a Go test that re-runs `make check` would recurse. |
| REQ-87 | The four licensed diff shapes are a rule for REVIEWING the Phase 2/3 diff, not a runtime property. The two Go-side diffs it will license are pre-identified as D3 and D8 so the review has a checklist. |
| REQ-88 | "Any other changed line is a regression" — the scope gate for the implementation. RDR 0011's `TestReq120_…` is the mechanized approximation already in the tree; see D7. |
| REQ-100 | SC-6's regression sweep is `make check` over every checked-in model and fixture with a licensed-diff oracle. See D6; REQ-71 and REQ-85 take the assertable slices. |
| REQ-106 | A DOC obligation on `docs/` (the authoring doc must state the guard-atom rule and the escape-row "otherwise" idiom). Phase 4's deliverable; no Go test can assert prose. |

REQ-101 is not an orphan: it is a negative REQ satisfied by writing no
two-toolchain test, which is what this phase did.

### Tests with no REQ (0)

None. Every test function opens with a `// REQ-N: "<quote>"` comment naming
the clause it pins, and every helper is called from at least one such test.

---

## Notes for Phase 2 / Phase 3

1. **Two shipped tests will go red and must be extended, never weakened**:
   `internal/table/helpers_test.go::clonedRowSliceFields` (D3) and
   `internal/graphlint/findings_0006_test.go::TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet`
   (D4). Both stay EXACT-set assertions; only their expected sets grow.
2. **`internal/cli/flow_demand_0011_test.go::TestReq35And118_…`** hard-codes
   four `flow resolve` envelopes and needs `"emit":{}` added after
   `"gates":[]` in each — REQ-87 licensed shape (ii). Verified with a stub
   (D8). It is an independent confirmation of C4's position clause.
3. **`internal/cli/flow_next_0011_test.go::TestReq120_…`** is RDR 0011's
   diff guard and FIRES on this RDR's production edits (D7). Resolve it
   before Phase 2's first green commit.
4. **`table.IsDecisionTable`** is the helper name the tests assume (D5). The
   record does not fix it; if a different spelling is chosen, three tests
   need the new name.
5. The 103 `[dump]`-carrying fixtures under `internal/table/testdata/` must
   all gain `"emit"` (REQ-71, ASSUMPTION-8) or they refuse at load.

---

# Phase 2 — REQ-MVV run end to end, actual output

`go test ./...` is GREEN across every package. `TestMVV_StatelessDecisionTables`
passes all 18 subtests. Below is the MVV driven a second time through the
BUILT BINARY (`go build ./cmd/intrastate`), so the recorded output is what
a caller sees rather than what a harness reconstructs. Only the temp-dir
model path is run-specific.

## REQ-78 / step 2 — partial lint reports the coverage gap

```
$ intrastate lint --model partial.toml --as=json ; echo exit=$?
{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","findings":[{"code":"graph-coverage-gap","message":"group dt/decide over rows [cell-xp cell-xq cell-yp] leaves 1 of 4 assignments in its scoped product uncovered for the no_match arm; the coverage union must equal the scoped product","model":"dt","severity":"blocking","rule":"cell-xp","element":"dt/decide","dimension":"a,b","class":"no_match"}]}
exit=2
```

Exit 2, ONE coverage finding, naming the uncovered cell by its group and
its two guard dimensions (`a,b`) with the count `1 of 4`. This is the
POSITIVE finding proving coverage RAN with no root declared — the whole
point of C5's ∅ seeding.

## REQ-79 / step 3 — complete lint is exactly `[]`

```
$ intrastate lint --model complete.toml --as=json ; echo exit=$?
{"type":"ok","data":{"findings":[]}}
exit=0
```

No code from the 0006 taxonomy is present (A10).

## REQ-80 / step 3 variant — escape row yields only the advisory

```
$ intrastate lint --model escape.toml --as=json ; echo exit=$?
{"type":"ok","data":{"findings":[{"code":"graph-coverage-closed-by-escape","message":"the coverage of group dt/decide is closed by the bare escape row \"otherwise\" rather than proved over its declared domains","model":"dt","severity":"info","rule":"otherwise","element":"dt/decide"}]}}
exit=0
```

Exit 0 with exactly the one info-severity advisory — never a bare green.

## REQ-81 / step 4 — resolve returns rule + emit with no `--artifact`

```
$ intrastate flow resolve --model complete.toml --outcome decide \
    --tag a=y --tag b=q --as=json ; echo exit=$?
{"type":"ok","data":{"model":"<model>","revision":"","observed":{"a":"y","b":"q"},"owned":{},"readers":[],"outcome":"decide","rule":"cell-yq","gates":[],"emit":{"code":"4","verdict":"delta"},"next":{},"writes":{},"clear":[],"escaped":false}}
exit=0
```

`data.rule` is the fourth rule's id; `data.emit` is its AUTHORED block with
keys in byte order (`code` before `verdict`); `next`/`writes`/`owned` are
empty and `readers` is `[]`. `emit` sits immediately after `gates` on the
wire, and `escaped` still sits after `clear` (REQ-39/REQ-40/REQ-43).

### REQ-38 / SC-5 — the `{}`-never-`null` arm

```
$ intrastate flow resolve --model noemit.toml --outcome decide \
    --tag a=y --tag b=q --as=json
… "gates":[],"emit":{},"next":{} …
```

### REQ-45 / REQ-99 — text mode, through the GENERIC renderer

```
$ intrastate flow resolve --model complete.toml --outcome decide --tag a=x --tag b=p
clear: (none)
emit.verdict: alpha
escaped: false
gates: (none)
model: <model>
next: (none)
observed.a: x
observed.b: p
outcome: decide
owned: (none)
readers: (none)
revision:
rule: cell-xp
writes: (none)

$ intrastate flow resolve --model noemit.toml --outcome decide --tag a=y --tag b=q | grep '^emit'
emit: (none)
```

One path-qualified leaf per pair, and `flatten`'s empty-container arm for an
unauthored block. No per-verb special case was added (REQ-46).

## REQ-82 / step 5 — dump renders `emit`; `next` carries empty `required`

```
DumpOrder = [identity source kind outcome atoms next writes requires_owned gate escape emit]

MODEL dt rows=4 outcomes=decide
identity=dt.cell-xp source=dt:cell-xp kind=transition outcome=decide atoms=[a.eq=x@all; b.eq=p@all] next=[] writes=[] requires_owned=[] gate=[] escape=[] emit=[verdict=alpha]
identity=dt.cell-xq source=dt:cell-xq kind=transition outcome=decide atoms=[a.eq=x@all; b.eq=q@all] next=[] writes=[] requires_owned=[] gate=[] escape=[] emit=[verdict=beta]
identity=dt.cell-yp source=dt:cell-yp kind=transition outcome=decide atoms=[a.eq=y@all; b.eq=p@all] next=[] writes=[] requires_owned=[] gate=[] escape=[] emit=[verdict=gamma]
identity=dt.cell-yq source=dt:cell-yq kind=transition outcome=decide atoms=[a.eq=y@all; b.eq=q@all] next=[] writes=[] requires_owned=[] gate=[] escape=[] emit=[code=4; verdict=delta]
```

`emit` renders on EVERY row with no `[dump]` declared, appended last after
`escape`, as `key=value` pairs in key order bracketed like `writes` (see
`cell-yq`: `[code=4; verdict=delta]`). Taken through `table.Dump` per
deviations D2 — the repo ships no root `dump` CLI verb.

```
$ intrastate flow next --model complete.toml --as=json ; echo exit=$?
{"type":"ok","data":{…,"candidates":[{"rule":"cell-xp","outcome":"decide","required":[],"unknown":[{"key":"a","reason":"absent"},{"key":"b","reason":"absent"}],"next":{},"writes":{},"clear":[]}, …×4]}}
exit=0
```

Exit 0, every candidate with empty `required` (A11), and NO `emit` anywhere
in the payload or on any candidate — `0010:BR4` rejected (REQ-50).

## REQ-83 / step 6 — negative controls

```
$ intrastate lint --model omitted.toml --as=json ; echo exit=$?
{"code":"graph-lint-failed",…,"findings":[{"code":"graph-dangling-edge","message":"the model declares no initial owned state; add an `[initial]` table assigning every always-present owned tag","model":"dt","severity":"blocking","element":"model"}]}
exit=2

$ intrastate flow resolve --model owned.toml --outcome decide --tag a=x --tag b=p --as=json ; echo exit=$?
{"code":"flow-model-invalid","message":"the selected model could not be loaded","findings":[{"code":"malformed_model_declaration","message":"malformed_model_declaration: [model] class \"decision-table\" declares owned=1; a decision-table model declares zero tags of provenance owned","locator":"<model>:1"}]}
exit=2

$ intrastate lint --model models/rdr.toml --as=json ; echo exit=$?
{"type":"ok","data":{"findings":[]}}
exit=0
```

The class-omitted control takes `0006:C18`'s missing-root finding at
`element = model`, unchanged. The owned-tag control refuses
`malformed model declaration` carrying the literal token `owned=1` and the
class. The checked-in navigator model lints exactly as before and is
untouched.

## REQ-84 — end state

Discharged as the inverse invariant across all four cells: what the fixture
AUTHORS as `[rule.emit]` is what `flow resolve` RETURNS as `data.emit`,
both by `reflect.DeepEqual` on the decoded object and byte-for-byte against
the raw wire object. `owned` is `{}` and `readers` is `[]` on every cell —
no owned state, no accessor, no artifact — and the shipped 0005
state-machine fixture resolves to `advance-draft` with `owned.status=draft`
exactly as it did, its `emit` the empty object.
