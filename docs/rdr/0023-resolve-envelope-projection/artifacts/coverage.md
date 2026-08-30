# Coverage — RDR 0023 Resolve-envelope projection opt-out

Phase 1 test-authoring artifact. Maps every REQ in
`req-list.md` (REQ-1…REQ-137) to the test that would fail if the clause were
broken, and records the red-gate tally.

## Files

| File | Test funcs | Concern |
| --- | --- | --- |
| `internal/cli/flow_fixtures_0023_test.go` | 0 (helpers) | authored partition lists, the full-emitted-line measurand, payload key-order/raw-field readers, invocation shapes, reflective struct access |
| `internal/cli/flow_projection_0023_test.go` | 20 | `0023:C1` — flag surface, projected key set, absence, differential, width, text mode, projection site |
| `internal/cli/flow_partition_0023_test.go` | 10 | `0023:C2` — declared partition, reflective completeness, carrier constraint, pointer-conversion scope, always-keep core, observed reader execution, A9's empty-container arm |
| `internal/cli/flow_mvv_0023_test.go` | 5 | MVV steps 1–6 incl. the 40% pass bar, the pre-change golden, docs/help surface, the deliverable boundary |

Suite command: `go test ./...`. No implementation code was written.

## Red gate

Run: `go test ./internal/cli/ -count=1 -v` on `worktree-rdr-0023`
(RDR 0024 has NOT landed; `decision_table_0010_test.go:435` pins
`NumField() != 14`, and REQ-64/REQ-90 are asserted against 14 accordingly).

| Outcome | Count |
| --- | --- |
| New 0023 test funcs | **35** |
| RED (fail against the unimplemented tree) | **31** |
| GREEN — `PRESERVES-EXISTING` (see below) | **3** |
| SKIP — conditional on 0024's landing order | **1** |
| Pre-existing tests in `internal/cli` | 267 passing, **0 moved** |

Package-wide totals: `FAIL: 31  PASS: 267  SKIP: 1`.

**No predecessor oracle moved.** Every failure in `internal/cli` is one of
the 31 new 0023 tests; the 0005 / 0010 / 0011 suites pass unedited, which is
`0023:A3`'s claim discharged by execution rather than by assertion
(REQ-135). The new tests reach the still-unbuilt symbols — the declared
assignment table, the reader-execution seam, the total walker — reflectively
or through the production source, never by a direct Go reference, so the
package compiles and each test fails for its own clause instead of taking
all 128 shipped test files down with a build error.

## PRESERVES-EXISTING

Three tests legitimately pass today because they pin **existing,
already-correct behaviour that RDR 0023 explicitly preserves**. Each is
listed with its justification rather than rewritten into an artificial
failure.

| Test | REQs | Why it passes today |
| --- | --- | --- |
| `TestReq64And90And91And135_TheConversionChangesFieldTypesNeverTheFieldCount` | 64, 90, 91, 135 | Pins `NumField() == 14` plus `Emit`/`Gates` adjacency — the count and position `decision_table_0010_test.go` already pins. The RDR's own claim is that the pointer conversion changes field TYPES only, so this must be green before AND after. It goes red exactly when a field is added, which REQ-91 puts out of scope by construction. **[0010-carried]** |
| `TestReq6And14_TheShippedPartialWalkerCannotReachC1sScope` | 6, 14 | Asserts a property of SHIPPED code: `help_all.go::walkCommandTree` reaches strictly fewer commands than a total walk and drops the `completion` subtree. That is the ground C1's reuse prohibition rests on (`0023:A5`: 15 vs 9), so it is true today and must stay true — it goes red if a future edit makes the shipped walker total, at which point the prohibition needs re-deriving. |
| `TestReq29And30_The0011AllFenceIsUntouchedAndInvisibleToPlanOnly` | 29, 30 | **[0011-carried]** — the negative this RDR promises not to disturb. `--all` rides `next` alone, `--select` is absent from `resolve`, and the exact-name `Lookup("all")` probes cannot see a resolve-local `plan-only`. Green today by construction and green after; it goes red if landing this RDR widens the sibling fence. |

Their partial red-arms are already covered elsewhere: the flag-registration
half of REQ-29/30 is asserted in the red
`TestReq2And3And4And5And7And8And9And13And93_...` walk, so the fence's live
edge is not resting on a green test.

## Conditional SKIP

| Test | REQs | Why |
| --- | --- | --- |
| `TestReq77And87_IfDispositionsExistsItIsCarriedUnderTheFlag` | 77, 87 | `dispositions` is 0024's field and that RDR has not landed on this branch. The test skips with an explicit message and becomes live the moment the field exists, in either landing order — which is exactly the ordering tolerance `0023:A4` / REQ-87 record ("either RDR may land first; the second lands with `dispositions` already/newly in the plan group and no contract in either moves"). Making it unconditionally red would assert a field this tree does not have. |

## REQ × test mapping

Every REQ appears in at least one test's `// REQ-N:` comment and is
individually assertable. Status column: **R** = red, **G** = green
(PRESERVES-EXISTING), **S** = conditional skip.

### Flag surface and registration (`0023:C1`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-1 | `TestReq1And15And16_ResolveRegistersABooleanPlanOnlyFlagDefaultFalseNoShorthand` | R |
| REQ-2 | `TestReq2And3And4And5And7And8And9And13And93_PlanOnlyIsRegisteredOnExactlyFlowResolveAcrossTheWholeTree` | R |
| REQ-3 | same | R |
| REQ-4 | same (vacuity guard 1) | R |
| REQ-5 | same (total walk, no name-skip, no `Hidden` gate) | R |
| REQ-6 | `TestReq6And14_TheShippedPartialWalkerCannotReachC1sScope` | **G** |
| REQ-7 | `TestReq2And3And4...` (cobra's own initializers, precondition asserted) | R |
| REQ-8 | same (`completion` named as the discriminating member) | R |
| REQ-9 | same (`t.Fatalf` on a missing auto-generated command — loud, not vacuous) | R |
| REQ-10 | `TestReq10And11And12And132_PlanOnlyOnASiblingVerbIsTheSharedUsageBucketAndMintsNoCode` | R |
| REQ-11 | same (no minted code, no new exit group) | R |
| REQ-12 | same (`command-error`, exit 2, on all three siblings) | R |
| REQ-13 | `TestReq2And3And4...` (walker yields the unfiltered walked set; filter is the caller's) | R |
| REQ-14 | `TestReq6And14_TheShippedPartialWalkerCannotReachC1sScope` | **G** |
| REQ-15 | `TestReq1And15And16_...` (bool, default false, no shorthand) | R |
| REQ-16 | same (rejected spellings absent as aliases) | R |

### Projected payload: partition, key set, absence (`0023:C1`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-17 | `TestReq17And18And22And23And98And108_TheProjectedKeySetIsExactlyTheAuthoredLiteral` | R |
| REQ-18 | same (echo group named individually) | R |
| REQ-19 | `TestReq19And79And113_AProjectedKeyIsAbsentNotNullNotAnEmptyPlaceholder` | R |
| REQ-20 | `TestReq20And26And124And127_EveryCarriedFieldIsByteIdenticalToItsDefaultRendering` | R |
| REQ-21 | `TestReq21And121And122_TheCarriedKeysKeepTheirDefaultModeRelativeOrder` | R |
| REQ-22 | `TestReq17And18And22...` (authored literal, not a complement) | R |
| REQ-23 | same (8 keys + conditional `escape_class`) | R |
| REQ-24 | `TestReq24And25And123_PresenceRuleFieldsAreAssertedAgainstTheSameRunsDefaultNotALiteral` | R |
| REQ-25 | same (comparison against the SAME run's default, both arms) | R |
| REQ-26 | `TestReq20And26And124And127_...` (raw-byte comparison; no decode-and-compare) | R |
| REQ-27 | `TestReq27And28And95And96And97And129_DefaultModeIsByteIdenticalToTheCapturedPreChangeGolden` | R |
| REQ-28 | same | R |
| REQ-29 | `TestReq29And30_The0011AllFenceIsUntouchedAndInvisibleToPlanOnly` | **G** |
| REQ-30 | same | **G** |

### Report-only differential (`0023:C1`, `0023:S1`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-31 | `TestReq31And32And34And45And130And133_TheDifferentialDiffsToExactlyTheEchoKeys` | R |
| REQ-32 | same (exit codes asserted first) | R |
| REQ-33 | `TestReq33And100And131_RefusalEnvelopesAreByteIdenticalPlusOrMinusTheFlag` | R |
| REQ-34 | `TestReq31And32And34...` (strict key-subset) | R |
| REQ-35 | `TestReq35And37And40And41And42And43And44And46_TheInvokedReaderSetIsObservedExecutionNotARecomputation` | R |
| REQ-36 | `TestReq36And38And111_TheProjectedEncodingIsStrictlyShorterOnTheFullEmittedLineInBothModes` | R |
| REQ-37 | `TestReq35And37And40...` (reader pass is the reachable surface) | R |
| REQ-38 | `TestReq36And38And111_...` (`emittedLine` is the full NDJSON record) | R |
| REQ-39 | `TestReq39And57And82_TheTextSubsetAndTextWidthAssertionsAreSeparablyNecessary` | R |
| REQ-40 | `TestReq35And37And40...` (`invokedReaders` recomputation rejected by construction) | R |
| REQ-41 | same + `TestReq41And42_TheReaderRecordingSeamSitsAtTheExecutionAndIsFlagBlind` | R |
| REQ-42 | `TestReq41And42_...` (seam position = the `exec.Read` call, flag-blind) | R |
| REQ-43 | `TestReq35And37And40...` (default run supplies the expected set; non-empty control) | R |
| REQ-44 | same (asserts `readers` IS projected away, so the payload read is unwritable) | R |
| REQ-45 | `TestReq31And32And34...` | R |
| REQ-46 | `TestReq35And37And40...` (`owned`/`readers` cannot be suppressed independently) | R |

### Projection site and flag-blindness (`0023:C1`, `0023:C2`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-47 | `TestReq47And48And49And50And53_NothingUpstreamOfPayloadAssemblyReadsTheFlag` | R |
| REQ-48 | same | R |
| REQ-49 | same + `TestReq41And42_...` (flag absent from `flow_exec.go`) | R |
| REQ-50 | `TestReq47And48...` (via refusal-path equality across 4 decision shapes) | R |
| REQ-51 | `TestReq51And52And67_ProjectionTakesTheFullyAssembledPayloadAndOnlyDeletes` | R |
| REQ-52 | same (conditional assembly is observable as a plan field that tracks the flag) | R |
| REQ-53 | `TestReq47And48...` | R |

### Text mode (`0023:C1`, `0023:S5`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-54 | `TestReq54And94_BothModesRenderTheSameProjectedResultThroughTheOneGateway` | R |
| REQ-55 | `TestReq55And56And58And59And110_ProjectedTextLinesAreAByteIdenticalStableSubset` | R |
| REQ-56 | same (set membership, not subsequence) | R |
| REQ-57 | `TestReq39And57And82_...` | R |
| REQ-58 | `TestReq55And56...` + `TestReq39And57And82_...` (per-line byte identity) | R |
| REQ-59 | `TestReq55And56...` (3 repeated runs) + MVV step 3 | R |

### Partition declaration and reflective oracle (`0023:C2`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-60 | `TestReq60And61And62And63And68And69_EveryPayloadFieldIsAssignedToExactlyOneC2Group` | R |
| REQ-61 | same | R |
| REQ-62 | same | R |
| REQ-63 | same (the census's per-field assignment, authored) | R |
| REQ-64 | `TestReq64And90And91And135_TheConversionChangesFieldTypesNeverTheFieldCount` | **G** |
| REQ-65 | `TestReq65And66And70And71_TheAssignmentIsDeclaredInAStandaloneSiteNotAsAStructMarker` | R |
| REQ-66 | same | R |
| REQ-67 | same + `TestReq51And52And67_...` | R |
| REQ-68 | `TestReq60And61...` (all three assertions, bidirectional equality) | R |
| REQ-69 | same (a field with no entry fails at the first assertion) | R |
| REQ-70 | `TestReq65And66And70And71_...` (forbidden struct-tag carrier + required standalone site) | R |
| REQ-71 | same (carrier file ≠ projection file) | R |
| REQ-72 | `TestReq72And88And89_ThePointerConversionCoversTheFiveEchoFieldsAndNothingElse` | R |

### Always-keep core (`0023:C2`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-73 | `TestReq73And74And76_TheAlwaysKeepCoreIsProjectionInvariant` | R |
| REQ-74 | same (invariance, not unconditional presence; both arms) | R |
| REQ-75 | same — negative REQ: no separate always-keep oracle is minted; the test asserts core membership in S2's literal rather than quantifying over modes | R |
| REQ-76 | same (`revision` carried though constant-empty) | R |
| REQ-77 | `TestReq77And87_IfDispositionsExistsItIsCarriedUnderTheFlag` | **S** |

### Oracle discriminability — negative controls

| REQ | Test | Status |
| --- | --- | --- |
| REQ-78 | `TestReq35And37And40...` (skip-the-reader-pass control) + `TestReq36And38And111_...` (width vs key-set) + `TestReq20And26...` (mutated carried value) | R |
| REQ-79 | `TestReq19And79And113_...` (the `null` absence control, on raw bytes) | R |
| REQ-80 | `TestReq60And61...` (emitted-keys-equal-`plan`: red where S2 stays green) | R |
| REQ-81 | `TestReq81_TheWholeTreeWalkGoesRedWhenASecondCommandOrCompletionRegistersTheFlag` — both controls RUN, incl. planting on `completion` | R |
| REQ-82 | `TestReq39And57And82_...` (no-op text projection) + `TestReq55And56...` (dropped plan line) | R |
| REQ-83 | `TestReq83And112And126_TheOracleBatteryIsTheFiveScenariosAndOwesNoWallTimeBudget` | R |

### Construction obligations (`PRE`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-84 | `TestReq84And85And86And123_EmptyEchoContainersRenderAsBracesAndBracketsNeverNull` (A9 arm) + `TestReq2And3And4...` (S4 walker) | R |
| REQ-85 | `TestReq84And85And86And123_...` | R (default-mode half PASSES today; the `/projected` sub-test is red — see note below) |
| REQ-86 | same | R |
| REQ-87 | `TestReq77And87_...` | **S** |

### Mechanism (`0023:A2`, `PH1`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-88 | `TestReq72And88And89_...` (five echo fields → `*T`, `T` unchanged) | R |
| REQ-89 | same (`escape_class` stays plain `string`) | R |
| REQ-90 | `TestReq64And90And91And135_...` | **G** |
| REQ-91 | same — negative REQ | **G** |
| REQ-92 | `TestReq92_TheFlagRidesResolveAloneAndHandsRespondOKTheProjectedResult` | R |
| REQ-93 | `TestReq2And3And4...` (flag never in the shared registrar) | R |
| REQ-94 | `TestReq54And94_...` — negative REQ: no gateway edit, asserted as both-modes content agreement | R |

### Minimum Viable Validation (`0023:MVV`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-95 | `TestReq27And28And95And96And97And129_...` + `TestMVV0023_.../step1_...` | R |
| REQ-96 | `TestReq27And28...` (golden read from `<art>/mvv-step1-default-golden.json`) | R |
| REQ-97 | same | R |
| REQ-98 | `TestReq17And18And22And23And98And108_...` + `TestMVV0023_.../step2_...` | R |
| REQ-99 | `TestMVV0023_.../step3_projected_text_is_a_stable_byte_identical_subset` | R |
| REQ-100 | `TestReq33And100And131_...` + `TestMVV0023_.../step4_...` | R |
| REQ-101 | `TestMVV0023_.../step5_sibling_verb_refuses_and_the_structural_walk_passes` | R |
| REQ-102 | same — negative REQ: one sibling verb run, the walk is the assertion | R |
| REQ-103 | `TestMVV0023_.../step6_...` | R |
| REQ-104 | same — the 40% bar on all three checked-in fixtures | R |
| REQ-105 | same (message states the route-back, not a recorded number) | R |
| REQ-106 | same — negative REQ: the 79.4% figure is logged as a baseline, never gated | R |
| REQ-107 | `mvvCall0023` runs the checked-in `pricing-decision-table.toml` | R |

### Normative fixtures

| REQ | Test | Status |
| --- | --- | --- |
| REQ-108 | `TestReq17And18And22And23And98And108_...` (runs `pricingCall`, the normative fixture) | R |
| REQ-109 | `TestMVV0023_.../step6_...` (A1 row-S1 baseline logged, per REQ-114) | R |
| REQ-110 | `TestReq55And56And58And59And110_...` (dropped lines named, not counted) | R |
| REQ-111 | `TestReq36And38And111_...` (width asserted unconditionally on 3 shapes) | R |

### Phase 2 oracle enumeration

| REQ | Test | Status |
| --- | --- | --- |
| REQ-112 | `TestReq83And112And126_...` (five oracles, four needing the flag as subject) | R |
| REQ-113 | `TestReq19And79And113_...` (absent-not-null rides the key-set oracle) | R |
| REQ-114 | `TestReq114And125And128_...` — negative REQ: no oracle asserts A1's percentages | R |

### Docs and help (`PH3`)

| REQ | Test | Status |
| --- | --- | --- |
| REQ-115 | `TestReq115And116And117And118And119And120_TheFlagIsDocumentedAndTheOmittedGroupIsNamed` | R |
| REQ-116 | same (`docs/cli-reference.md`, `llms.txt` carry the flag) | R |
| REQ-117 | same | R |
| REQ-118 | same (non-empty usage string) | R |
| REQ-119 | same (help names every omitted echo field) | R |
| REQ-120 | same (`flowResolveExtendedDesc` gains the flag line) | R |

### Determinism and cross-cutting

| REQ | Test | Status |
| --- | --- | --- |
| REQ-121 | `TestReq21And121And122_...` | R |
| REQ-122 | same + `TestReq55And56...` (text stability across repeated runs) | R |
| REQ-123 | `TestReq24And25And123_...` (`emit` as `{}`) + `TestReq84And85And86And123_...` (`{}`/`[]` when empty) | R |
| REQ-124 | `TestReq20And26And124And127_...` | R |
| REQ-125 | `TestReq114And125And128_...` — negative REQ: no in-payload version marker | R |
| REQ-126 | `TestReq83And112And126_...` — negative REQ: no wall-time bar in any 0023 test | R |
| REQ-127 | `TestReq20And26And124And127_...` (one encoder, raw-byte equality) | R |
| REQ-128 | `TestReq114And125And128_...` — negative REQ: tag parsing refuses, unchanged ± the flag | R |
| REQ-129 | `TestReq27And28And95And96And97And129_...` | R |

### Failure modes and consequences

| REQ | Test | Status |
| --- | --- | --- |
| REQ-130 | `TestReq31And32And34And45And130And133_...` | R |
| REQ-131 | `TestReq33And100And131_...` | R |
| REQ-132 | `TestReq10And11And12And132_...` (the same parse refusal makes skew loud both ways) | R |
| REQ-133 | `TestReq31And32And34...` (the diff IS exactly the echo keys, both directions) | R |
| REQ-134 | `TestReq134And136And137_TheFlagIsOptInAndNoConsumerPassesItOnLanding` | R |
| REQ-135 | `TestReq64And90And91And135_...`; discharged package-wide by 267 unmoved passes | **G** |
| REQ-136 | `TestReq134And136And137_...` — negative REQ | R |
| REQ-137 | same — negative REQ: no consumer call site migrated | R |

## Orphans

**REQs with no test: none.** All 137 are mapped above; a scripted sweep for
`REQ-N` over `internal/cli/*_0023_test.go` reports zero missing.

**Tests with no REQ: none.** All 35 new test functions carry at least one
`// REQ-N:` comment. `flow_fixtures_0023_test.go` declares helpers only and
contains no test function, so it is exempt by construction.

## Notes for Phase 2

1. **The pre-change golden is Phase 1's first commit step.** REQ-96 requires
   it captured from the PRE-CHANGE binary on a clean tree, BEFORE the struct
   is edited. `TestReq27And28And95And96And97And129_...` reads it from
   `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`
   and fails with that instruction until it exists. A golden captured after
   any Phase 1 edit does not satisfy REQ-96/REQ-97 (`0023:A-8` fixes the
   ordering, not the path). The default-mode line this tree currently emits
   — captured for reference only, not checked in by this phase — is the
   pricing 2×2 call with `--tag tier=paid --tag region=eu`.

2. **REQ-85's test is deliberately half-green.** The default-mode arm of
   `TestReq84And85And86And123_...` PASSES today and must: it is A9's
   Phase-1 oracle, written to be meaningful BEFORE the pointer conversion
   lands, and it is what catches a nilled container silently changing
   default-mode bytes once the conversion is in. Its `/projected` sub-test
   is red. The top-level function therefore reports FAIL and is counted red
   above, but Phase 2 should not "fix" the green arm.

3. **The three unbuilt symbols are reached indirectly on purpose.** The
   declared assignment table (REQ-66/70/71) is located by reading
   `internal/cli/*.go` sources; the reader-execution seam (REQ-41/42) is
   located by reading `flow_exec.go`; the total walker (REQ-3/5/13) is
   written inline in the test rather than calling a production symbol.
   Phase 2 may replace the inline walker with a shipped one — C1 only
   forbids reusing `help_all.go::walkCommandTree`. If the seam or the table
   ships under a name the probe does not recognise, extend the recognised
   markers rather than weakening the assertion.

4. **REQ-MVV output slot: RECORDED (Phase 2).** All six MVV steps pass.

   Step 6's measured byte table, as `TestMVV0023_.../step6_...` logs it —
   the full emitted line, trailing newline excluded, which is C1's unit:

   ```
   pricing-decision-table:      360 B → 152 B (57.8% saved)
   release-grammar/begin:       482 B → 238 B (50.6% saved)
   review-state-machine/approve: 368 B → 158 B (57.1% saved)
   ```

   All three clear the 40% bar, so REQ-105's route-back is not triggered.

   **Against `evidence/spikes/a1-byte-width.md`.** The in-test defaults
   above are inflated on the DEFAULT side only, because the harness must
   pass an absolute `--model` path (the package's tests run with
   `cwd = internal/cli`) and `model` echoes that argument verbatim. Re-run
   from the repo root with the spikes' own relative-path invocations, the
   shipped projection reproduces A1's recorded rows EXACTLY:

   | row | invocation | default | projected | saved | A1 records |
   | --- | --- | ---: | ---: | ---: | --- |
   | S1 | pricing 2×2, `--tag tier=paid --tag region=eu` | 290 B | 152 B | 47.6% | 290 → 152, 47.6% |
   | S2b | release `begin`, seeded `phase=idle,build-id=none`, `--tag risk=0 --tag 'checks=[]'` | 412 B | 238 B | 42.2% | 412 → 238, 42.2% |

   Both projected widths match the spike to the byte, which is the stronger
   result: the plan group the flag retains is exactly the one A1 measured,
   so REQ-105's "the plan group is carrying materially more than A1
   measured" condition is not merely un-triggered but positively excluded.

   The review fixture's `approve` resolves through a bare escape row, so it
   is also the MVV's witness for the CONDITIONAL ninth key — `escape_class`
   is present by default and present under the flag, per REQ-25/REQ-74:

   ```
   default:   {"type":"ok","data":{"model":…,"revision":"","observed":{},"owned":{"status":"draft"},
                "readers":["review-state"],"outcome":"approve","rule":"approve-otherwise","gates":[],
                "emit":{},"next":{},"writes":{},"clear":[],"escaped":true,"escape_class":"no_match"}}
   projected: {"type":"ok","data":{"revision":"","rule":"approve-otherwise","gates":[],"emit":{},
                "next":{},"writes":{},"clear":[],"escaped":true,"escape_class":"no_match"}}
   ```

   The A1 row-S1 baseline is LOGGED, never asserted, per REQ-114 ("A1's byte
   table is recorded evidence, not a test assertion"), and the
   not-checked-in 79.4% figure is not gated, per REQ-106.

## Phase 2 result

`go test ./...` — all packages pass. `golangci-lint run` — 0 issues.
`make check` (fmt-check, vet, lint, build, graph-lint, docs-check, and
`go test -race`) — green.

All 35 new 0023 test functions pass, together with the 267 pre-existing
`internal/cli` tests, none of which moved — `0023:A3`'s "no predecessor
oracle moves" claim discharged by execution (REQ-135). The one conditional
skip (`dispositions`, pending 0024) remains a skip, which is the ordering
tolerance REQ-87 records.

Seven deviations are recorded in `deviations.md`; none is
`needs author decision`. Five are TEST-FIXTURE, one is a SPEC-DEFECT in
REQ-116's premise about `llms.txt` resolved from REQ-117's own statement of
the obligation, and one is an IMPL-DECISION recording the reader seam's
form, which C1 explicitly delegates.
