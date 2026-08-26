# Triage — RDR 0008 Ownership of the recognized-outcome tag key name

Single-pass unattended triage of the roborev findings left by the Stage 8
implementation launch. Window frozen once: `BASE=1718a70`, `HEAD=2cd1845`
(14 commits). Batch label: `batch:rdr-0008`.

Spine: per-commit auto-reviews (jobs 6255, 6260, 6261) — 7 findings.
Cross-commit net: `roborev review --since 1718a70` (job 6263) — 4 findings,
all duplicates of the spine, charged to their per-commit job.
Fix-commit sweep (bounded, one round): job 6264 — 1 finding.

**8 raw findings → 5 distinct after dedup, plus 1 from the bounded sweep.**

## Verdicts

| # | Job(s) | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| F1 | 6255, 6263 | `internal/table/reserved_key_0008_test.go:56` | KATA-BUG | test-oracle | test-coverage | kata `wtne` |
| F2 | 6255, 6263 | `internal/table/reserved_key_0008_test.go:208` | DROP `rdr-adjudicated` | test-oracle | test-coverage | drop |
| F3 | 6255, 6263 | `internal/table/reserved_key_0008_test.go:532` | KATA-BUG | test-oracle | test-coverage | kata `sf4e` |
| F4 | 6260 (×2) | `internal/cli/reserved_key_adv_0008_test.go:95,:172` | DROP `rdr-adjudicated` + KATA-BUG | interface | design-conformance | drop; kata `wagv` (RDR 0006 carrier) |
| F5 | 6260, 6261, 6263 | `internal/table/reserved_key_adv_0008_test.go:115` | FIX-NOW | test-oracle | test-coverage | commit `0ead6e5` |
| F6 | 6264 | `internal/table/reserved_key_adv_0008_test.go:117` | KATA-BUG | test-oracle | test-coverage | kata `jb9v` |

ODC distribution: `test-oracle` 5, `interface` 1. No `n/a-not-a-defect`
pre-filter cases (the intentionally-RED Phase 1 commit drew findings about
coverage gaps, not about its redness).

## Grounding detail

### F1 — observed-provenance declaration untested → kata `wtne` (medium)

REQ-11 and REQ-30 both bind "owned **or observed**"; the suite pins owned
only. Confirmed reachable and correct-today by probe through the real
`table.Load`: `cat="reserved_tag_key" offending="recognized" remedy=""
rule="reserved-tag-key/author-must-rename"`. `load.go:182` guards
`key == RecognizedTagKey && decl.Provenance != ProvenanceRecognized` —
provenance-agnostic, no earlier category preempts it. Fixture sweep found zero
observed declaration fixtures. Deviation D6 scopes only the `unknown tag`
arm, a different arm — no adjudication covers this. Regression-only exposure.

### F2 — write-position check → DROP `rdr-adjudicated`

REQ-23 is a **negative** requirement. `0008:C2` verbatim: "`[rule.write]` is
covered by a *different* pre-existing rule and must not be folded into the
sentence above … what rejects it is RDR 0002's `write to non-owned tag`
category … **This RDR adds no write-position check and depends on that
category holding.**" Predecessor RDR 0002 already pins it
(`TestAdv2_RuleWriteToNonOwnedTagIsRefusedAsSuch`, plus roundtrip/dump/
accessors coverage). The reserved-key check has exactly two call sites
(`load.go:175-190`, `load.go:321-325`); neither is on the `[rule.write]` path,
which fails independently at `normalize.go:528,578` and `load.go:276`. The
proposed regression has no code path to arrive through. Asking 0008 to
re-test a predecessor's contract it explicitly disclaims is scope creep.

### F3 — combined case-fold + trim → kata `sf4e` (low)

REQ-59 requires the advisory when folding and trimming apply **together**; the
case table covers each alone. No live defect: `advisory.go`'s `isNearMiss`
composes `strings.EqualFold(strings.TrimSpace(key), RecognizedTagKey)` in one
step, and a probe with `[tags." Recognized"]` fires the advisory correctly
(`authored=" Recognized" reserved="recognized" rule="reserved-tag-key/near-miss"`).
D5/D9 govern the advisory's *delivery channel*, not its trigger set — no
reversal. Regression-only: a refactor to a disjunction would silently drop
the combined case with the suite green.

### F4 — payload/advisory not asserted on the CLI wire → DROP + kata `wagv`

**Dropped as stated**, per deviation D9 (`SPEC-UNDER`, Phase 3c): "Adjudicated
as belonging to RDR 0005/0006, not this RDR." `0008:C3` binds the three fields
"at the data level"; REQ-93 defers the exit-code map ("The consumer is a test
stub, **not** RDR 0005's exit-code map"); the RDR's Approach item 3 disclaims
the CLI surface. Routing them through `clierr.Finding` would add CLI public
surface 0008's Normative Contracts never name — the ADDITIVE-IS-NOT-EXEMPT
`SPEC-UNDER` trap.

roborev's premise is also inaccurate at HEAD: both tests **do** assert 0008's
owned half unconditionally (`f.Offending`/`f.Remedy`/`f.Rule`;
`a.Authored`/`a.Reserved`/`a.Rule`) before any `t.Skip`; the skipped half is
the wire-side restatement, and both retained assertions were mutation-checked
in Phase 3c.

**But the delivery gap is real**, and D9 named a successor obligation that had
no home. Filed as kata `wagv` against **RDR 0006's carrier**: `lint.go:95-103`
maps every `table.Load` refusal to a bare `model-invalid` + `Detail:
err.Error()` with no `Findings`, which RDR 0006's own REQ-91 forbids
("not through a verb-local wrapper or a **text-only `Detail` string**"). Whole-
category blast radius — 0002's ~20 categories included, not just 0008's rule.
The block traces to 0006's skeleton commit `d0194d9`; 0008 never touched it.
The acceptance test is already written and waiting: un-skip the two `t.Skip`
blocks in `internal/cli/reserved_key_adv_0008_test.go`.

### F5 — probabilistic determinism oracle → FIX-NOW `0ead6e5`

roborev was right: a 200-sample histogram is a probabilistic oracle for an
exact property. Replaced with exact assertions on a single load —
`Offending == "outcome"`, `Remedy == "recognized"`,
`Rule == table.RuleKernelOwned` (constant, not literal), REQ citations kept.

Does **not** over-specify REQ-82: that clause's non-assertion scope is a
**same-direction** fixture (two recognized-provenance declarations, e.g.
`[tags.outcome]` + `[tags.result]`), pinned category-only at
`reserved_key_0008_test.go:66-101`. This fixture is **cross-direction**
(`[tags.recognized]` owned + `[tags.outcome]` recognized), which the Phase 3c
disposition adjudicated: sorting "narrows no latitude". REQ-27/33 make the
rule identifier a token a golden test asserts byte-for-byte.

**Mutation check:** reverting `load.go` to a map range failed the strengthened
test on all 20 runs — and the map range picked the *same* wrong entry every
time, which is precisely the luck the old histogram would have accepted as
proof of determinism. `load.go` byte-identical to HEAD; only the test file
changed.

### F6 — single load does not prove stability → kata `jb9v` (low)

The fix-commit auto-review for `0ead6e5` raises the converse: one invocation
could pass by luck under a map-order regression. Correct, and complementary —
the strongest oracle is exact-assertion-**under**-repetition. Filed rather
than re-fixed, per the bounded one-round fix-then-review sweep. No live
defect (`load.go` sorts by construction), and the mutation check above shows
the current oracle already catches the targeted regression.

## Filed kata (all `src:roborev` + `batch:rdr-0008` + `lifecycle:queued`)

| short_id | severity | pri | area | title |
|---|---|---|---|---|
| `wtne` | medium | 3 | internal-table | Add observed-provenance fixture for the reserved `recognized` tag key |
| `wagv` | medium | 3 | internal-cli | `lint.go` collapses every `table.Load` refusal into a text-only `Detail` |
| `sf4e` | low | 4 | internal-table | Pin REQ-59's combined case-and-whitespace near-miss |
| `jb9v` | low | 4 | internal-table | Combine the exact payload assertion with repeated loads |

All four related to `z53t` (RDR 0008's `kind:rdr-tracked` tracker). No open
questions were baked into any kata — every finding was settled from evidence.

## Jobs closed

6255, 6260, 6261, 6263, 6264 — each answered via `roborev respond` citing the
verdict, its evidence, and the kata short_id or fix commit.

21 further open jobs are **out of window** (pre-`BASE` main-branch debt,
ancestors of `BASE`), untouched by this run.
