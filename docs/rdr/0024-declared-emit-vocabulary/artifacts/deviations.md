# Deviations — RDR 0024 Phase 2

Every entry carries exactly one Type. Entries are appended in the order
they were reached during implementation.

---

## D1 — `TestReq3_0024`'s `set` fixture carried two defects

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**Site:** `internal/table/emit_grammar_0024_test.go::TestReq3_0024_KindIsOneOfFourTokensAndSetIsExcluded`

Phase 1's fixture for "`set` is not an admitted emit kind" authored
`kind = "set"` AND `elements = ["alpha"]`. `elements` is a
tag-declaration key an emit declaration does not admit (REQ-20), so
strict decoding refuses it as `unknown_schema_field` before the kind arm
is ever reached — the test asserted `malformed_emit_declaration` and got
the decoder's slug.

Both refusals are correct; the fixture simply carried two defects, which
`0024:S1` itself forbids ("One defect per fixture — the pipeline is
fail-fast, so a fixture carrying two defects proves nothing about the
second"). Dropped the `elements` line so the fixture is single-defect.
The assertion is unchanged.

---

## D2 — RDR 0010's `TestReq12` step-census technique cannot survive C2's mandated insertion

**Type:** DEPENDENCY-LIMIT
**Status:** resolved from the record's own evidence
**Site:** `internal/table/class_0010_test.go::TestReq12_TheAgreementCheckSitsInsideTheLicensedWindow`

`0024:C2` normatively fixes "**Two steps, in this order, at this
position** … both are inserted immediately after `loadTags` in
`load.go::run`'s step slice". Doing so turned a shipped RDR 0010 test
red.

**Investigated before editing, per REQ-83** ("Any assertion outside this
list going red is a signal the append did more than C4 licenses —
investigate, do not update it"). The finding is that the red is in the
test's *identification technique*, not in its assertion:

- The REQ it carries is `0010:C1`'s "Any step in that window satisfies
  this clause; the RDR fixes the window, not the step's name."
- To locate the class-agreement step without naming it, the test
  computes it as *whichever step is new relative to RDR 0002's ten-step
  baseline*, then asserts that step sits between `loadTags` and
  `checkAccessorBindings`.
- That technique holds only while no later RDR ever adds a loader step.
  0024 adds two, by a normative fence, so the census returns three
  "added" steps and the test fatals before reaching its own assertion.

The subject of the assertion — the class check's position in the window
— is untouched by this RDR: `checkClassAgreement` still sits after
`normalizeRules` and before `checkAccessorBindings`, and the new steps
land well outside that window. So the resolution is to add the two
0024 steps to the test's baseline list, which restores the technique's
premise and leaves the assertion, its REQ, and its window unchanged.

This is not a weakened test: with the baseline corrected, moving
`checkClassAgreement` outside the window still fails it, and adding an
undeclared third step still fatals it.

**Evidence:** `docs/rdr/0024-declared-emit-vocabulary.md` §Normative
Contracts `0024:C2` (the insertion is normative and names the position);
the same fence's note that the five steps between `loadTags` and
`normalizeRules` are independent of both new steps, so no ordering
constraint of 0010's is disturbed.

---

## D3 — `Categories()` gained three entries

**Type:** IMPL-DECISION
**Status:** mechanical translation
**Site:** `internal/table/category.go`

An append-only enum extension, which the launch prompt classes as
mechanical. Recorded only because REQ-71 makes the registration a
contract obligation rather than an implementation detail: "the wire slug
and the registered category are one decision, not two". The three
constants and the three `Categories()` entries land together, in
declaration order, after `CatReservedTagKey`.

---

## D4 — `headerLine` extracted from `tagHeaderLine`'s technique rather than reusing it

**Type:** IMPL-DECISION
**Status:** mechanical translation
**Site:** `internal/table/load.go`

`0024:C2` fixes that a declaration defect "keys on the top-level
`[emit.<key>]` header, which `tagHeaderLine`'s technique reaches
unchanged", and ASSUMPTION-3 leaves the helper's name and arity to the
implementer ("Two helpers or one is an implementation choice; recovering
BOTH lines is not").

Shipped as three unexported helpers in `load.go`: `emitHeaderLine(src,
key)` for the declaration side, `emitRuleLine(src, ruleID)` for the
rule side, and a shared `headerLine(src, spellings...)` carrying the
scan — the decoration stripping, and the decline-on-anything-ambiguous
rule (zero or two matches yield 0) that `tagHeaderLine` documents at
length.

`tagHeaderLine` itself is left byte-for-byte unedited. Folding it into
`headerLine` would be a refactor outside this RDR's REQ set, and
`internal/table/emit_steps_0024_internal_test.go::TestReq33_0024` asserts
`load.go` still declares it.

No new PUBLIC surface: all three helpers are unexported.

---

## D5 — `dispositions` registered with RDR 0023's shipped §D1 oracle (0023-FIRST order)

**Type:** DEPENDENCY-LIMIT
**Status:** resolved from the record's own evidence
**Site:** `internal/cli/flow_partition.go`,
`internal/cli/flow_partition_0023_test.go`

ASSUMPTION-7 assumes the **0024-first** order ("0023 is Final and
unimplemented at HEAD, per Phase 3's own statement"), under which
`0024:S7` discharges as a checked-in note at the `dispositions` field.
That assumption is FALSE at this HEAD: RDR 0023 has landed —
`internal/cli/flow_partition.go` ships the declared ECHO/PLAN partition
table `resolvePayloadGroups`, and `flow_projection.go` drives the
projection off it.

REQ-84 names this branch explicitly and fixes what it costs: "**The §D1
registration obligation, whichever record lands second.** … if 0023
lands first, this change registers `dispositions` with the oracle it
shipped". So the alternative leg applies, and it is the one taken:

- `resolvePayloadGroups` gains `"Dispositions": groupPlan`. Without it
  `TestAdversarialF2_TheShippedDeclarationTableIsTotalOverThePayloadStruct`
  fails on a field with no entry, which is exactly the failure mode
  `0023:C2` designed that oracle to catch.
- `flow_partition_0023_test.go`'s AUTHORED census map gains
  `"dispositions": "plan"`, for the same reason and in the same
  direction: the oracle asserts the authored census against the shipped
  struct in BOTH directions, so a new field with no census entry fails
  it.

The side is PLAN, not ECHO, and `0023:C2` decides it rather than this
RDR: the PLAN group is "rule identity, gate results, **authored answers
and their interpretations**, planned next/writes/clear, the escape
disposition, and revision". A disposition token IS the interpretation of
an authored answer — it is nothing the caller sent, and it is produced
by joining the selected row against the model's declarations. It rides
beside `emit`, which the same census assigns PLAN.

The `0024:S7` note at the field declaration is shipped as well, because
REQ-106 fixes its content independently of which order obtained and it
costs nothing to carry both.

**Evidence:** `docs/rdr/0024-declared-emit-vocabulary.md` `PH3` §"The §D1
registration obligation, whichever record lands second" (REQ-84);
`internal/cli/flow_partition.go`'s shipped `resolvePayloadGroups` and
its `groupPlan` prose; `0023:C2` as quoted in that file.

---

## D6 — `planGroup0023` / `projectedKeyLiteral0023` fixtures extended

**Type:** DEPENDENCY-LIMIT
**Status:** mechanical translation
**Site:** `internal/cli/flow_fixtures_0023_test.go`

Consequential to D5. RDR 0023's suite carries two authored key lists
that enumerate the PLAN side of the wire — `planGroup0023` and
`projectedKeyLiteral0023` (`0023:S2`'s explicit literal list, in default
order minus the deleted keys). Both are censuses of a payload the
appended field changes, and both go red on the append alone.

`dispositions` is inserted at its wire position — immediately after
`emit` — in each, which is the same single-key insertion `0024:C4`
licenses everywhere else. No 0023 assertion is relaxed and no key order
is re-sorted.

---

## D7 — the licensed golden regeneration ran to 28 sites exactly

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**Site:** `internal/cli/flow_demand_0011_test.go`

`PH3`/A2's census predicted 28 whole-payload assertions pinning the
`"emit":{…},"next"` adjacency: 3 `resolveGolden` entries plus the
25-entry `shippedResolveGoldens` sweep. ASSUMPTION-5 reads that count as
evidence rather than a target and asks for a deviation only if it comes
back materially different.

It did not: the census returned exactly 28, and every one changed in
exactly one way — `"dispositions":{}` inserted between `"emit"` and
`"next"`. Recorded for the evidentiary reason `RISK` gives ("**The
residual risk is not diff width but evidentiary loss**"): the
regeneration was a mechanical single-key insertion applied to a tree
carrying only the `dispositions` append, and
`TestReq80_0024_EveryResolveGoldenCarriesTheInsertedDispositionsKey`
re-asserts the adjacency per site.

---

## D9 — `dtDecl0024` did not declare `code`, which the 0010 fixture emits

**Type:** TEST-FIXTURE
**Status:** mechanical translation
**Site:** `internal/cli/flow_dispositions_0024_test.go`,
`internal/cli/flow_emit_surface_0024_test.go`

Phase 1's `dtDecl0024` declares only `verdict`, but the RDR 0010 fixture
it is inserted into (`dtModel0010`) has `cell-xp` authoring TWO keys —
`verdict = "alpha"` and `code = "1"`. `0024:C2` fixes that strictness is
whole-model, not per-key, so every fixture built on it refused
`unknown_emit_key` before any disposition was ever joined.

Added `[emit.code]` with the FLAT spelling (`domain = ["1"]`) to
`dtDecl0024`, and to the four replacement declarations that substitute it
wholesale and would otherwise lose it. The flat spelling is what the
surrounding assertions already require: `TestReq57_0024`'s expected join
is exactly `{"verdict": "route"}`, so `code` must be declared (or the
model refuses) AND carry no disposition (or the join is two entries).

Three sites were left alone because they already handle `code`: the two
`TestReq60_0024` tests declare it themselves, `TestReq59_0024`
substitutes it out of the rule, and `TestReq74_0024` declares it
`scalar`. The "a bare `[emit]` table" case still replaces the whole
declaration with `[emit]`, so it remains a genuine zero-declaration
model.

No assertion changed. This is the same class of fixture defect as D1 —
Phase 1 could not run these fixtures against a real loader, so a
whole-model strictness consequence was not visible until the loader
existed.

---

## D8 — RDR 0023's pre-change golden is folded, not regenerated

**Type:** DEPENDENCY-LIMIT
**Status:** resolved from the record's own evidence
**Site:** `internal/cli/flow_mvv_0023_test.go`,
`docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`

A2's 28-site census (`PH3`) was taken on a tree where RDR 0023 had NOT
landed, so it never enumerated 0023's own MVV golden. That golden went
red on the `dispositions` append, as a 29th whole-payload assertion the
licensed diff does not name.

Investigated per REQ-83 rather than updated. The finding is that this
one is **not regenerable**, and that regenerating it would silently
destroy evidence rather than refresh a fixture:

- `0023:A-8` fixes its provenance normatively — "the capture commit
  precedes the conversion commit, and a golden captured after any Phase 1
  edit does not satisfy REQ-96/REQ-97". It is a capture of a
  **pre-0023 binary**, taken on a `git stash`-clean tree.
- No post-0023 tree can produce those bytes again. Re-capturing it here
  would write the CURRENT build's output into a file whose whole claim is
  that it predates the current build.
- 0023's own REQ-97 names what that costs: "S1–S5 all compare the new
  build against itself, so none of them can see a default-mode regression
  that moves both sides together". Regenerating converts the battery's
  only pre-change comparison into another self-comparison and closes
  nothing.

So the golden's bytes are left **untouched**, and the COMPARISON folds
out 0024's appended key before comparing — the same move the same test
already makes for the checkout prefix via `foldCheckoutRoot`, and for the
same reason: remove the one difference that is not this assertion's
subject, then assert byte-identity over everything else.

The fold is exact rather than permissive, so it strengthens rather than
weakens the site. It anchors on the PREDECESSOR key
(`"emit":{…},"dispositions":{…}`), so it fails if `dispositions` is
absent, if it appears more than once, or if it drifts off its fixed
position after `emit` — which is precisely what `0024:C4` fixes about it.

**Evidence:** `0023:A-8` and `0023:REQ-96`/`REQ-97` as quoted in
`internal/cli/flow_mvv_0023_test.go`'s own doc comments; the existing
`foldCheckoutRoot` precedent in the same function.
