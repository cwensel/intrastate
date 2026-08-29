Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0024 declared emit vocabulary

Verdict: **RECONCILED**. All nine Critical Assumptions are terminal, every
named spike has a captured run, and no BLOCKER was raised.

## Stage 5 preflight

`intrastate flow resolve --outcome lens` returned `/rdr-reconcile` under rule
`lens-foundational-row-complete` — Profile `foundational`, the
cove -> 3amigo -> critique -> repeatability row complete. The two completion
outcomes a folder cannot show both answered `none` with no caveat:

- `--outcome critique` -> `critique-foundational-complete` (two base models,
  diffed; dual-model obligation discharged).
- `--outcome repeatability` -> `repeatability-full-complete` (variant `full`;
  run-1/2/3 and the diff written).

The `mid`/`large` Determinacy judgement does not apply at `foundational`.

## Open set and dispositions

| Item | Source | Disposition | Evidence / plan |
| --- | --- | --- | --- |
| A1 top-level `[emit]` key is free | 2 | VERIFIED (Stage 4) | `internal/table/source.go::sourceDoc`, `::decodeStrict` |
| A2 `dispositions` append breaks no consumer | 2 | VERIFIED (Stage 4) | `internal/cli/decision_table_0010_test.go::TestReq39_...`; 28 enumerated sites |
| A3 no `[emit.<key>]` authored repo-wide | 2 | VERIFIED (Stage 4) | `evidence/spikes/a3-corpus-sweep.md` |
| A4 kind conformance without `internal/guard` | 2 | VERIFIED (Stage 4) | `internal/table/load.go::ConformValue` |
| A5 consumer emit answers enumerable | 2 | VERIFIED (Stage 4) | `0010:A6`; 54 + 29 `[rule.emit]` blocks, fixed literals |
| A6 DMN allowed-values on output clauses | 2 | **DOWNGRADED** | Not load-bearing; corroborating citation only. Searched across four corpora and the sibling repo, not found (`evidence/research/resolve-prior-art.md`). Plan deliberately UNSCHEDULED — the alignment argument stands on `0002:C22` + two opened peer citations. |
| A7 two-shaped `domain` closure leg | 2 | **DOWNGRADED** | Decode legs Verified (`evidence/spikes/c1-dual-form-domain.md`); CLOSURE leg carried under Testing Strategy scenario 1, which runs in Phase 1 and gates its close. Not MVV-critical. |
| A8 `TestReq146` reaches a model-level carrier | 1 | REFUTED, costless (absorbed at pre-lock) | Sweep is scoped to `table.Row`; C3 + scenario 4 already rewritten as sole oracle. No design change. |
| A9 emit-refusal source line recoverable | 1,3 | **VERIFIED at reconcile** | `evidence/spikes/a9-rule-side-locator.md` |
| Stale `conform*` reuse pointers (Phase 1, Testing Strategy) | 4 (absorption audit) | FIXED in-pass | Now name `ConformValue`; amendment sweep over the `conform` token clean |
| Exactness delta: "total order" x4, "deterministic" x1 | 4 | No action owed | Each site names its comparator (`sortFindings`, `compareRefs`) or its oracle (scenario 4) |

## A9 — the item this gate existed to catch

A9 was MVV-critical and could not be deferred: C2 makes it normative that all
three refusal categories carry a source line, and MVV step 2 asserts "the
offending block's SOURCE LINE, not `:1`". Its rule-side leg was Pending and
unread.

The spike CONFIRMS the operative claim — writable without new plumbing —
because `internal/table/normalize.go::normalizeRule` is a `*loader` method, so
`l.src` and the rule `id` are both in scope at the sole `sourceRule.Emit` read.
No data has to be threaded; only a helper written.

It also **corrected the technique the RDR had made normative**, which is the
substantive catch. The draft keyed the rule-side locator on the `[rule.emit]`
header. That header text is byte-identical for every rule — 63 occurrences in
`rdr-status.toml`, 4 in the in-repo example — so `tagHeaderLine`'s
exactly-one-match-or-zero contract returns 0 on every real multi-rule model.
The locator anchors on the rule id instead (`id = "<ruleID>"`, matching the
VALUE, since `[model]` also carries an `id` key), scanning forward to the emit
block. Uniqueness is guaranteed by `CatDuplicateRuleID` before any emit work
runs, and the ordering precondition was checked mechanically: every rule's `id`
line precedes its `[rule.emit]` header in both real models (63/63, 4/4). The
decoder was ruled out rather than assumed away — `pelletier/go-toml/v2 v2.2.4`
exposes position only on `DecodeError` with unexported fields, and a rule-side
refusal is decided after a successful strict decode.

Had this shipped as drafted, every rule-side emit refusal would have rendered
line 0 on every multi-rule model — silently, since the naive helper returns 0
rather than failing.

## Absorption audit

Delegated over the four lens dirs. Absorption verified round by round:
cove 9/9, 3amigo 30/30 (26 folded, 4 charted), critique 21/21 merged IDs
(16 folded, 5 charted), repeatability 13 disagreements + 7 GUESS clusters
(13 D-items and G-1..G-6 folded, G-7 charted). Claimed dispositions were
spot-checked against the projected contract text rather than trusted.

One real residue, fixed in this pass: Phase 1 and Testing Strategy still named
`conform`/`conformKind`/`conformDomain` "through a narrow helper" — the exact
misdirection cove F2 was filed against, contradicted by C3 ("that struct
literal IS the adapter — no extraction, no new exported helper"), A4, and the
Infrastructure Audit. Both now name `ConformValue`; Testing Strategy flags the
unexported three-arg `conform` as the wrong function.

Charted-to-successor items were verified genuinely out of scope (consumer seam
test and merge-day coverage; scalar-hollowing observability under the closed
advisory tier `0006:C17`; the `dispositions` plural naming held by
`JDR 0002 SS-D1`; a pre-existing `lint.go` help-text drift; repeatability G-7
bounding what the lens proved, not what the record says).

No unbacked new claims were introduced by the fix passes.

## Completeness

- No `_Draft placeholder._` and no seed-skeleton header in any body section.
- `## References` fully authored — peer elements, source reviewed, prior art,
  Stage 4 evidence, related issues. No template brackets.
- `rdr lint 0024` PASS. The five surviving `placeholder:survived` findings are
  all inside the Finalization Gate, which Stage 7 authors.
- All three previously named spikes have captured runs on disk; a fourth
  (`a9-rule-side-locator.md`) was added by this stage.

## Hard rules

- **Refutation**: none. A8 was refuted at pre-lock and is already absorbed with
  no contract change; nothing the RDR currently relies on is refuted.
- **MVV deferral**: none. A9, the one MVV-critical open item, was run now and
  Verified rather than deferred. A6 and A7 are downgraded and neither pins what
  the MVV proves.
