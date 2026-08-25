# Coverage — RDR 0004 Accessor Execution Safety Model

Stage 8, **Phase 1** (test authoring). Every REQ in
[`req-list.md`](req-list.md) mapped to the test that would FAIL if a
future change broke that clause, plus orphan analysis in both
directions.

## Status

| | |
| --- | --- |
| Test files | 11 under `internal/accessor/` |
| Top-level test functions | 80 (79 REQ tests + `TestMVV_…`) |
| **Red at Phase 1** | **77 of 80** |
| Green at Phase 1 | 3, each justified below |
| REQs covered | 125 of 126, plus REQ-MVV |
| REQs with no test | 1 — REQ-126, see Orphans |
| Tests with no REQ | 0 |

The implementation does not exist. `internal/accessor/skeleton.go` is
signatures only, every body returning a zero value, following RDR 0006's
Phase 1 precedent (`d0194d9`): the suite COMPILES so each test fails on
its own assertion naming the REQ it quotes, rather than the package
failing with one undefined-symbol error that attributes to no clause.

Every predecessor suite stays green (`internal/{cli,clierr,graphlint,
guard,resolve,table}`); only `internal/accessor` is red.

## Test files

| File | Cluster |
| --- | --- |
| `fixtures_0004_test.go` | shared fixture bindings and assertion helpers |
| `capability_0004_test.go` | capability, identity, selection, caller-supplied artifacts (REQ-1..12) |
| `read_completeness_0004_test.go` | read branch discipline and completeness (REQ-13..29, 103) |
| `seam_omission_0004_test.go` | absence crosses the seam as omission (REQ-30..34, 118) |
| `gate_0004_test.go` | gate invocation semantics (REQ-35..39) |
| `write_readback_0004_test.go` | write accessors, read-back, `<clear>` (REQ-40..60, 104) |
| `post_mutation_0004_test.go` | read-back incompleteness, post-mutation sense (REQ-61..67, 102, 107, 120, 125) |
| `boundary_0004_test.go` | timeout, no-prints, refusal-class set, disposition totality (REQ-68..82, 96..100, 109) |
| `validation_0004_test.go` | the eight validation arms, replay, ownership (REQ-83..91, 8, 9, 81, 99) |
| `desk_trace_0004_test.go` | desk-trace simultaneity, phases, discharge ledger (REQ-92..95, 101, 105..106, 111, 114..116, 119, 121..123) |
| `mvv_0004_test.go` | REQ-MVV — the nine scenarios with named controls |

## The three green-at-Phase-1 tests

None is a genuine coverage hole; each asserts against a surface that
already exists, which is what the exception in the Phase 1 brief
anticipates.

| Test | Why green is correct |
| --- | --- |
| `TestReq77_AccessorClassesAreDisjointFromTheKernelRefusalKinds` | **[0001-delivered, green-at-phase-1]** — asserts RDR 0001's SHIPPED `resolve.RefusalKinds()` is closed at five and that this RDR does not extend it. The kernel half is real today; the disjointness half becomes discriminating once `RefusalClasses()` is populated. |
| `TestReq73_TheCaptureControlItselfDiscriminates` | **Negative control, green by design** — ORA 5 makes the capture control normative: "A test that deliberately prints must fail the capture assertion; without that control this scenario passes vacuously." This test proves the harness CAN fail, so it must pass. Were it red, `TestReq72_…` would be worthless. |
| `TestReq96_AuthorityTravelsAsATypedBindingNotACommandString` | **Structural, satisfied by the declared surface** — asserts a `Definition` reaches the world through a typed `Binding` and carries no command-string field for authority to hide in (Alternatives 2/3/5 rejected). This is a property of the type declaration, which the skeleton genuinely establishes; an implementation cannot regress it without a signature change the compiler catches. |

Sub-tests that pass inside otherwise-red parents are the **named
negative controls** ORA requires (the valid validation fixture returning
the empty set, truncation-alone still reporting `incomplete_read`, the
assignment fixture still asserting presence). By construction they must
pass; they are what stop the positive assertions passing vacuously.

## Anti-tautology pass

Sixteen tests initially passed against the zero-value skeleton for the
wrong reason — an empty snapshot satisfies "the absent key is omitted",
a write that never runs satisfies "no write ran". Each was rewritten
with a positive precondition the stub cannot meet (the binding ran
exactly N times; the carried key IS present; the artifact holds the
planned value). Thirteen went red; the three above are the residue.

## REQ × test

| REQ | Test | Label |
| --- | --- | --- |
| REQ-1 | `TestReq1_CapabilityVocabularyIsClosedAtThree` | HAPPY PATH |
| REQ-2 | `TestReq2_OffCapabilityInvocationIsCapabilityMismatch` | ADVERSARIAL |
| REQ-3 | `TestReq3_IdentityTripleResolvesToExactlyOneBinding` | ADVERSARIAL |
| REQ-4 | `TestReq4_SameIdUnderTwoCapabilitiesIsTwoIdentitiesNotARebinding` | BOUNDARY |
|  | `TestReq5_SelectionIsByCapabilityTableAndId` | BOUNDARY |
| REQ-5 | `TestReq5_SelectionIsByCapabilityTableAndId` | BOUNDARY |
| REQ-6 | `TestReq6_UnboundNameIsUnknownAccessorAtRuntime` | ADVERSARIAL |
| REQ-7 | `TestReq7_CapabilitySpellingsAreTheCanonicalNames` | INPUT EDGE |
| REQ-8 | `TestReq8_TheDefinitionConsumesRDR0002sCarrierRatherThanASecondOne` | BOUNDARY |
| REQ-9 | `TestReq8_TheDefinitionConsumesRDR0002sCarrierRatherThanASecondOne` | BOUNDARY |
| REQ-10 | `TestReq10_BindingReceivesTheCallerSuppliedArtifactRole` | HAPPY PATH |
| REQ-11 | `TestReq11_AmbientArtifactDiscoveryFailsValidation` | ADVERSARIAL |
| REQ-12 | `TestReq12_UnsuppliedArtifactRoleIsExecutionFailureNotItsOwnClass` | INPUT EDGE |
| REQ-13 | `TestReq13_ReadReturnsTypedValuesOrATypedRefusal` | HAPPY PATH |
| REQ-14 | `TestReq14_ReadDoesNotMutateTheAuthoritativeArtifact` | ADVERSARIAL |
| REQ-15 | `TestReq15_RequestedKeySetIsTheDefinitionsCarriedKeys` | HAPPY PATH |
| REQ-16 | `TestReq16_RequestedSetIsPinnedNotDerivedFromWhatResolved` | ADVERSARIAL |
| REQ-17 | `TestReq17_MissingOrEmptyRequestedKeySetFailsValidation` | INPUT EDGE |
| REQ-18 | `TestReq16_RequestedSetIsPinnedNotDerivedFromWhatResolved` | ADVERSARIAL |
| REQ-19 | `TestReq19_SuccessBranchCarriesExactlyTheRequestedKeys` | HAPPY PATH |
| REQ-20 | `TestReq20_OneUnreadableKeyRefusesTheWholeRead` | ADVERSARIAL |
| REQ-21 | `TestReq21_AbsenceIsAValueAndUnreadabilityIsARefusal` | DOMAIN EDGE |
| REQ-22 | `TestReq20_OneUnreadableKeyRefusesTheWholeRead` | ADVERSARIAL |
| REQ-23 | `TestReq20_OneUnreadableKeyRefusesTheWholeRead` | ADVERSARIAL |
| REQ-24 | `TestReq19_SuccessBranchCarriesExactlyTheRequestedKeys` | HAPPY PATH |
| REQ-25 | `TestReq25_IncompleteReadIsDistinctFromExecutionFailureAndTimeout` | BOUNDARY |
| REQ-26 | `TestReq26_IncompleteReadNamesTheKeysItCouldNotRead` | BOUNDARY |
| REQ-27 | `TestReq27_TimeoutOutranksIncompleteReadWhenBothAreTrue` | BOUNDARY |
| REQ-28 | `TestReq28_ReservedClearLiteralReadAsAValueIsUnreadable` | DOMAIN EDGE |
| REQ-29 | `TestReq29_UnclassifiedKeyDefaultsToUnreadableNotAbsent` | ADVERSARIAL |
| REQ-30 | `TestReq30_AbsentKeyIsOmittedFromTheOwnedSnapshot` | DOMAIN EDGE |
| REQ-31 | `TestReq30_AbsentKeyIsOmittedFromTheOwnedSnapshot` | DOMAIN EDGE |
| REQ-32 | `TestReq32_NoPlaceholderValueCrossesTheSeamForAnAbsentKey` | ADVERSARIAL |
| REQ-33 | `TestReq33_ThePlaceholderEncodingTypechecksAndSilentlyRetiresTheRefusal` | ADVERSARIAL |
| REQ-34 | `TestReq34_MatchOnlyKeyIsOutsideThisRDRsSurface` | DOMAIN EDGE |
| REQ-35 | `TestReq35_GateVerdictVocabularyIsClosedAtThree` | HAPPY PATH |
|  | `TestReq35_AllowIsATypedGateResult` | HAPPY PATH |
| REQ-36 | `TestReq36_IndeterminateIsARefusalNeverAFalseAllowOrDeny` | ADVERSARIAL |
| REQ-37 | `TestReq76_TheRefusalClassSetIsExactlyTheEightWireSpellings` | BOUNDARY |
|  | `TestReq37_DenyIsATypedResultCarryingAReasonAtTheAccessorBoundary` | BOUNDARY |
| REQ-38 | `TestReq37_DenyIsATypedResultCarryingAReasonAtTheAccessorBoundary` | BOUNDARY |
| REQ-39 | `TestReq39_GateTimeoutAndExecutionFailureAreAccessorRefusalsNotGateResults` | BOUNDARY |
| REQ-40 | `TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly` | HAPPY PATH |
| REQ-41 | `TestReq76_TheRefusalClassSetIsExactlyTheEightWireSpellings` | BOUNDARY |
|  | `TestReq41_WriterNamingANonOwnedTagFailsValidation` | ADVERSARIAL |
| REQ-42 | `TestReq42_NoWriteRunsWithoutASuccessfulPlansWrites` | ADVERSARIAL |
| REQ-43 | `TestReq42_NoWriteRunsWithoutASuccessfulPlansWrites` | ADVERSARIAL |
| REQ-44 | `TestReq44_ClearSentinelIsRDR0002sReservedConstant` | HAPPY PATH |
| REQ-45 | `TestReq45_ClearingWriteRemovesTheKeyAndReadBackAssertsAbsence` | DOMAIN EDGE |
| REQ-46 | `TestReq45_ClearingWriteRemovesTheKeyAndReadBackAssertsAbsence` | DOMAIN EDGE |
|  | `TestReq46_AClearThatStoresTheLiteralIsReadBackMismatch` | ADVERSARIAL |
|  | `TestReq46_AClearWhoseKeySurvivesWithAnyValueIsReadBackMismatch` | ADVERSARIAL |
| REQ-47 | `TestReq47_ClearingAKeyTheArtifactDoesNotHoldSucceeds` | BOUNDARY |
| REQ-48 | `TestReq28_ReservedClearLiteralReadAsAValueIsUnreadable` | DOMAIN EDGE |
| REQ-49 | `TestReq49_AssignmentControlStillAssertsPresenceAndEquality` | BOUNDARY |
| REQ-50 | `TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly` | HAPPY PATH |
| REQ-51 | `TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly` | HAPPY PATH |
|  | `TestReq51_ReadBackIsEqualityNotContainment` | DOMAIN EDGE |
| REQ-52 | `TestReq51_ReadBackIsEqualityNotContainment` | DOMAIN EDGE |
| REQ-53 | `TestReq55_ReadBackMismatchIsAWriteFailureOnBothArms` | ADVERSARIAL |
| REQ-54 | `TestReq54_ReadBackReadsTheSameRoleNamedByTheWriteBinding` | ADVERSARIAL |
| REQ-55 | `TestReq55_ReadBackMismatchIsAWriteFailureOnBothArms` | ADVERSARIAL |
| REQ-56 | `TestReq56_WriterWithoutReadBackFailsValidation` | INPUT EDGE |
| REQ-57 | `TestReq55_ReadBackMismatchIsAWriteFailureOnBothArms` | ADVERSARIAL |
| REQ-58 | `TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly` | HAPPY PATH |
| REQ-59 | `TestReq59_ReadBackDoesNotAssertArtifactLevelFidelity` | BOUNDARY |
| REQ-60 | `TestReq60_NonOwnedComparisonExemptsNewKeysAndThePlannedOwnedKeys` | BOUNDARY |
| REQ-61 | `TestReq61_ReadBackThatCannotReadAComparedKeyIsReadBackIncomplete` | DOMAIN EDGE |
| REQ-62 | `TestReq61_ReadBackThatCannotReadAComparedKeyIsReadBackIncomplete` | DOMAIN EDGE |
|  | `TestReq66_ReadBackIncompleteNamesWhatTheCallerMustReRead` | BOUNDARY |
| REQ-63 | `TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense` | BOUNDARY |
| REQ-64 | `TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense` | BOUNDARY |
| REQ-65 | `TestReq65_NoCompensationAfterAPostMutationRefusal` | ADVERSARIAL |
| REQ-66 | `TestReq66_ReadBackIncompleteNamesWhatTheCallerMustReRead` | BOUNDARY |
| REQ-67 | `TestReq67_MismatchAndIncompleteAreDistinctNotOneWriteFailedShape` | BOUNDARY |
| REQ-68 | `TestReq68_EveryInvocationIsBounded` | BOUNDARY |
| REQ-69 | `TestReq69_TimeoutIsDistinctFromExecutionFailureAndReadBackMismatch` | BOUNDARY |
| REQ-70 | `TestReq70_MissingOrNonPositiveTimeoutFailsValidation` | INPUT EDGE |
| REQ-71 | `TestReq68_EveryInvocationIsBounded` | BOUNDARY |
| REQ-72 | `TestReq72_TheAccessorPackageNeverPrints` | ADVERSARIAL |
| REQ-73 | `TestReq72_TheAccessorPackageNeverPrints` | ADVERSARIAL |
|  | `TestReq73_TheCaptureControlItselfDiscriminates` | ADVERSARIAL |
| REQ-74 | `TestReq75_RefusalClassesAreExposedForTheCLIToMap` | BOUNDARY |
| REQ-75 | `TestReq75_RefusalClassesAreExposedForTheCLIToMap` | BOUNDARY |
| REQ-76 | `TestReq76_TheRefusalClassSetIsExactlyTheEightWireSpellings` | BOUNDARY |
| REQ-77 | `TestReq77_AccessorClassesAreDisjointFromTheKernelRefusalKinds` | ADVERSARIAL |
|  | `TestReq81_TheResidualIsLeftOpenNotClosedHere` | BOUNDARY |
| REQ-78 | `TestReq78_TheRefusalCarriesNoExitCodeOrCLICode` | BOUNDARY |
| REQ-79 | `TestReq79_EveryInputClassLandsOnANamedOutcome` | BOUNDARY |
| REQ-80 | `TestReq79_EveryInputClassLandsOnANamedOutcome` | BOUNDARY |
| REQ-81 | `TestReq81_TheResidualIsLeftOpenNotClosedHere` | BOUNDARY |
| REQ-82 | `TestReq82_RefusalPayloadCarriesTheDiagnosisTuple` | BOUNDARY |
| REQ-83 | `TestReq84_EightValidationArmsEachAssertingItsOwnNamedCode` | BOUNDARY |
| REQ-84 | `TestReq84_EightValidationArmsEachAssertingItsOwnNamedCode` | BOUNDARY |
|  | `TestReq84_TheValidationCodeSetIsClosedAtEight` | BOUNDARY |
| REQ-85 | `TestReq84_EightValidationArmsEachAssertingItsOwnNamedCode` | BOUNDARY |
| REQ-86 | `TestReq86_ReplayProducesTheSameDisposition` | HAPPY PATH |
| REQ-87 | `TestReq86_ReplayProducesTheSameDisposition` | HAPPY PATH |
| REQ-88 | `TestReq86_ReplayProducesTheSameDisposition` | HAPPY PATH |
| REQ-89 | `TestReq89_TheDispositionRecordsTheFiveReplayInputs` | BOUNDARY |
| REQ-90 | `TestReq90_WriteThenReadBackIsTheNormativeOrdering` | BOUNDARY |
| REQ-91 | `TestReq90_WriteThenReadBackIsTheNormativeOrdering` | BOUNDARY |
| REQ-92 | `TestReq92_ThePhaseDeliverablesAreOneReachableSurface` | BOUNDARY |
| REQ-93 | `TestReq92_ThePhaseDeliverablesAreOneReachableSurface` | BOUNDARY |
| REQ-94 | `TestReq92_ThePhaseDeliverablesAreOneReachableSurface` | BOUNDARY |
| REQ-95 | `TestReq92_ThePhaseDeliverablesAreOneReachableSurface` | BOUNDARY |
| REQ-96 | `TestReq96_AuthorityTravelsAsATypedBindingNotACommandString` | BOUNDARY |
| REQ-97 | `TestReq96_AuthorityTravelsAsATypedBindingNotACommandString` | BOUNDARY |
| REQ-98 | `TestReq100_TheExecutorLoadsNoConfigAndDiscoversNoPath` | BOUNDARY |
| REQ-99 | `TestReq81_TheResidualIsLeftOpenNotClosedHere` | BOUNDARY |
| REQ-100 | `TestReq100_TheExecutorLoadsNoConfigAndDiscoversNoPath` | BOUNDARY |
| REQ-101 | `TestReq101_DeskTraceStep2AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-102 | `TestReq102_DeskTraceStep5AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-103 | `TestReq103_ReadHasExactlyTwoBranchesNeverAThird` | BOUNDARY |
| REQ-104 | `TestReq54_ReadBackReadsTheSameRoleNamedByTheWriteBinding` | ADVERSARIAL |
| REQ-105 | `TestReq101_DeskTraceStep2AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-106 | `TestReq101_DeskTraceStep2AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-107 | `TestReq65_NoCompensationAfterAPostMutationRefusal` | ADVERSARIAL |
| REQ-108 | `TestReq46_AClearThatStoresTheLiteralIsReadBackMismatch` | ADVERSARIAL |
| REQ-109 | `TestReq69_TimeoutIsDistinctFromExecutionFailureAndReadBackMismatch` | BOUNDARY |
| REQ-110 | `TestReq82_RefusalPayloadCarriesTheDiagnosisTuple` | BOUNDARY |
| REQ-111 | `TestReq92_ThePhaseDeliverablesAreOneReachableSurface` | BOUNDARY |
| REQ-112 | `TestReq40_WriteSucceedsWhenReadBackHoldsThePlannedValueExactly` | HAPPY PATH |
| REQ-113 | `TestReq84_EightValidationArmsEachAssertingItsOwnNamedCode` | BOUNDARY |
| REQ-114 | `TestReq101_DeskTraceStep2AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-115 | `TestReq101_DeskTraceStep2AssertionsHoldSimultaneously` | BOUNDARY |
| REQ-116 | `TestReq116_TheWriteOutcomeTableDiscriminatesThreeWays` | BOUNDARY |
| REQ-117 | `TestReq75_RefusalClassesAreExposedForTheCLIToMap` | BOUNDARY |
| REQ-118 | `TestReq118_AbsentRequiredKeyReachesTheResolverAsOwnedStateUnavailable` | DOMAIN EDGE |
| REQ-119 | `TestReq116_TheWriteOutcomeTableDiscriminatesThreeWays` | BOUNDARY |
| REQ-120 | `TestReq63_PostMutationRefusalsCarryTheAppliedButUnverifiedSense` | BOUNDARY |
| REQ-121 | `TestReq116_TheWriteOutcomeTableDiscriminatesThreeWays` | BOUNDARY |
| REQ-122 | `TestReq123_ThePendingAssumptionsAreDischargedPerRule` | BOUNDARY |
| REQ-123 | `TestReq123_ThePendingAssumptionsAreDischargedPerRule` | BOUNDARY |
| REQ-124 | `TestReq27_TimeoutOutranksIncompleteReadWhenBothAreTrue` | BOUNDARY |
| REQ-125 | `TestReq61_ReadBackThatCannotReadAComparedKeyIsReadBackIncomplete` | DOMAIN EDGE |
| REQ-MVV | `TestMVV_AccessorExecutionSafetyModel` | HAPPY PATH |

## Orphans

### REQs with no test — 1

| REQ | Disposition |
| --- | --- |
| REQ-126 | "Prerequisites carried unchecked at Gate PASS…" — a **status disclosure about the record's checkboxes at lock**, not a behaviour the accessor boundary can exhibit or violate. No input exists for which an implementation could take a wrong branch, so no test can fail if it "broke". Recorded as deviation **D9**; the three peer capabilities it names ARE exercised where this RDR consumes them (`TestReq118_…`/`TestReq34_…` drive RDR 0001's shipped kernel; `TestReq8_…`/`TestReq15_…`/`TestReq44_…` consume RDR 0002's shipped carrier and sentinel). RDR 0003's clause is a guard-package property, outside this RDR's surface. |

### Tests with no REQ — 0

Every test function's doc comment opens with at least one `REQ-N:` or
`REQ-MVV:` citation and carries exactly one label.

## REQ-MVV decomposition

`TestMVV_AccessorExecutionSafetyModel` is one runnable end-to-end test
over ONE fixture flow (one read, one gate, one write accessor over
caller-supplied artifact roles), decomposed into the nine numbered
scenarios TS 1-9, each with the negative control ORA names.

| Scenario | Sub-test | Negative control asserted |
| --- | --- | --- |
| 1 — definition validation | `1_unsafe_definition_validation_eight_arms` | the valid fixture returns the EMPTY set, checked FIRST, so a validator that rejects everything fails before any arm can pass |
| 2 — read / gate dispositions | `2_read_and_gate_dispositions` | the complete read; and the truncation refusal carries NO values, so a derived-key-set implementation reporting a one-key success fails |
| 3 — write read-back | `3_write_read_back_round_trip_invariant` | the owned-mismatch and non-owned-mismatch arms differ in WHICH tag moved; success is the control against a read-back that always reports mismatch |
| 4 — replay stability | `4_replay_stability` | the injected-refusal run, compared for EQUALITY not absence of error |
| 5 — no package prints | `5_no_package_prints` | the positive half (returned value carries the class) plus the capture control |
| 6 — absence reaches the resolver | `6_absence_reaches_the_resolver` | a carried key resolves normally through the same path, so the refusal cannot pass vacuously |
| 7 — timeout precedence | `7_timeout_outranks_incomplete_read` | truncation alone must still report `incomplete_read` |
| 8 — post-mutation reporting | `8_post_mutation_reporting_and_no_compensation` | the `read_back_mismatch` fixture must still assert the artifact is WRONG rather than unverified |
| 9 — clearing write | `9_clearing_write` | the assignment fixture must still assert presence and equality |

### The Round-Trip / Inverse Invariant

RDR 0004 declares one, so a green exit code or "did not error" does NOT
discharge the MVV. Scenario 3 reconstructs the written state **through
the read path** and compares it value-for-value:

```
write -> read = expected owned-tag value identity
              + protected non-owned tag identity
```

- every planned owned key: the re-read holds exactly the plan's expected
  value (not containment — a write replaces the whole value, JDR 0001
  §D7(iv); byte equality over §D13's canonical JSON array for a set);
- every observed/recognized value present before the write: unchanged;
- scenario 9 extends the expected-value predicate to **absence** for a
  planned `<clear>`.

The Fidelity Table's deliberate weakening is honoured: the invariant is
asserted at TAG-VALUE granularity and never over the artifact encoding
(`TestReq59_…` asserts a re-materialized artifact still passes).

## Pending-assumption discharge ledger

A9, A10, and A11 are Pending at lock, to be discharged by the MVV. CA
A9's note on bundling directs "Verify and retire them per-rule rather
than flipping A9 as a unit", so
`TestReq123_ThePendingAssumptionsAreDischargedPerRule` exercises each
independently and a partial refutation stays visible.

| Property | Named scenario | Per-rule test |
| --- | --- | --- |
| A9 rule 1 — missing/empty requested key set | Scenario 1's eighth arm | `A9_rule_1_missing_or_empty_requested_key_set` |
| A9 rule 2 — absence reaches the resolver | Scenario 6 | `A9_rule_2_absence_crosses_as_omission_to_the_resolver` |
| A9 rule 3 — timeout outranks incomplete | Scenario 7 | `A9_rule_3_timeout_outranks_incomplete_read` |
| A9 rule 4 — read-back incomplete | Scenario 3 | `A9_rule_4_read_back_incomplete` |
| A10 — post-mutation reporting, no compensation | Scenario 8 | `A10_post_mutation_reporting_without_compensation` |
| A11 — clear is a removal, verified by absence | Scenario 9 | `A11_clear_is_a_removal_verified_by_absence` |

## The five forbidden spike shapes

CA A9 names four shapes the implementation MUST NOT reproduce (REQ-124)
and CA A10 a fifth (REQ-125). The fixtures are built to make each
observable rather than merely avoided.

| Spike shape | How the fixtures forbid it |
| --- | --- |
| derives the key set from the artifact (`expectedTagKeys`) | `keys` is pinned in the definition; artifacts deliberately disagree with it. `readBinding.lastRequested` records what the executor asked for (`TestReq101_…` step 3) |
| holds absence as an in-map sentinel `<absent>` | absence is the typed `KeyValue.Absent` flag; the seam OMITS the key (`TestReq32_…` names the sentinel as a defect) |
| sleeps before the key loop so timeout and truncation never overlap | `overlapReadBinding` meets the unreadable key AND runs past the deadline in one invocation, asserting `b.sawUnreadable` first |
| re-reads by cloning the tag map without going through `read` | the read-back goes through a `ReadBinding`; `TestReq61_…` asserts `r.reads != 0` |
| a re-read that cannot fail independently of the write | the write binding and the read-back reader are SEPARATE objects over one store; `readBackFailure` fails the re-read while the mutation has landed |

## Cross-layer grounding

Three assertions reach past this RDR's boundary into RDR 0001's shipped
kernel. Each was probed against the real `resolve.Resolve` before being
relied on:

| Input | Kernel disposition | Used by |
| --- | --- | --- |
| owned snapshot OMITS a `RequiresOwned` key | `owned_state_unavailable` naming the key | `TestReq118_…`, MVV 6 |
| owned snapshot CARRIES the key | plan (control, so the refusal is not vacuous) | `TestReq118_…`, MVV 6 |
| owned snapshot carries a `<absent>` PLACEHOLDER | **plan** — the regression is real, `TagSet.has` reads it as present | `TestReq33_…` |
| key consumed only through `Row.Match` | `no_match` — peer-owned residual, left OPEN | `TestReq34_…`, REQ-81 |

## Deviations recorded at Phase 1

D1-D8 are the pre-seeded cluster-gate entries, untouched. Phase 1
appends four, none needing an author decision:

- **D9** — REQ-126 is a status disclosure, not a testable obligation.
- **D10** — the eighth validation arm's code is spelled
  `missing_requested_key_set`, under the req-list's standing ASSUMPTION.
- **D11** — Q1 / D3's check ran: `gate denied` reconciles by citation to
  JDR 0001 §D9 on the layered reading. No fenced conflict survives, so
  D3 does not escalate.
- **D12** — Q2: the gate SITE is RDR 0005's. Phase 1 ships invocation
  semantics only, no scheduling and no aggregation.
