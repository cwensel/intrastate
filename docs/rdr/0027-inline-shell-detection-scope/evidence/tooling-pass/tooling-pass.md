Model: claude-opus-5[1m]

# Tooling Pass — cli/0027 (inline-shell detection scope)

Date: 2026-09-03 · Iteration 1 (base) · Run at Stage 7 finalize, pre-Gate.
Lint capture: `lint.txt` (`rdr lint --locking 0027`, exit 0,
`blocking=0 resolution=0 placeholder=5 advisory=6`).

## Findings

- **C1 — Template section coverage.** No `template:missing-section` finding. All
  Required sections present and authored: Metadata, Problem Statement, Critical
  Assumptions, Proposed Solution (Approach / Technical Design / Normative
  Contracts / Illustrative Code / Decision Rationale), Alternatives Considered
  (ALT1–ALT3 + Briefly Rejected), Context, Research Findings, Trade-offs,
  Implementation Plan (Prerequisites / MVV / Phases 1–3), Validation, References.
  `## References` is authored, not hollow. Conditional sections present where
  they apply (Load-Bearing Decisions, Capability Dependencies, Existing
  Infrastructure Audit, Performance Expectations).
- **C1 — hollow (placeholder:survived ×5, 502-505 / 511-523 / 529-531 /
  542-546 / 567-592).** All five are inside `§finalization-gate` (472-592), and
  all five are the Gate's own template guidance blocks for the items this pass
  is about to answer. `gate:inline` (472-592) says the same thing at the section
  level. NOT a regression: the lock replaces those four sub-sections with the
  `gate.md` pointer and the fifth (Cross-Cutting) is authored at the Gate.
  No `scaffold:row`, no `contract:template-example`, no seed-skeleton header.
- **C2 — Method-label vocabulary.** `ca_off_vocabulary=0`,
  `ca_off_vocabulary_ids=[]`. Every one of A1–A7 carries a Method field, all
  drawn from the sanctioned set: A1 `Source Search`; A2–A7 `Spike`. Clean.
- **C3 — Source Search self-reference.** One `Source Search` record: A1. Its
  Evidence anchor is `internal/table/load.go::carrierDefect` with concrete
  line cites into the source tree — resolves into the repo, not into
  {RDR_PATH} nor this RDR's artifact dir. No self-reference. Not fired.
- **C4 — Docs Only on load-bearing claims.** Zero `Docs Only` records. Clean.
- **C5 — Symbol resolution.** `anchors_total=35`, `anchors_unresolved=0`,
  `anchors_unlooked=0`, `peer_evidence_unresolved=0`. Every `source-anchor`
  edge resolves `true` (including the cross-repo prior-art anchor
  `pkg/cmd/alias/set/set.go::NewCmdSet`). No bare `file:line` anchor written
  into the body in place of a `path::Symbol`.
- **C6 — Status consistency.** Metadata Status `Draft` (form `none`, no
  re-entry qualifier). `ca_total=7`, `ca_verified=6`, `ca_pending=1`
  (`ca_pending_ids=["0027:A5"]`), `ca_unverified=0`.
  A5 is `Pending`, downgraded at Stage 6 and explicitly marked survivable.
  Checked for settled-fact reliance on A5 elsewhere in the record: the
  `trace` mini-check row 3c states outright that "A5 is no longer a witness
  for it (Pending, downgraded at Stage 6 …), so this row says the form is
  ADMITTED, not that it is owned"; §consequences says the stdin class "waits
  on the stdin successor, which no kata yet tracks (A5)"; §capability-
  dependencies marks the withholding capability `Deferred` and states "this
  record depends on none of it"; C1's `out of scope` line claims admission
  (delivered by the predicate) and ownership only prospectively. No prose
  treats A5's withholding mechanism as a settled fact. No checklist/Gate
  disagreement — the Gate is unwritten at sweep time. Clean.
- **C9 — Evidence-field budget (ADVISORY).** No `evidence:over-budget`
  finding in lint. The six advisory findings are the gate-section conformance
  set above.
- **C10 — Linking.** `label:contracts` not raised — the single normative block
  is labelled `**C1**`. `peer-evidence:no-element`: none. `edge:unresolved`:
  ONE found and FIXED IN-PASS — the Joint-check line (`0027:JC1`, line 202)
  cited `cli/0028:C4`, an element cli/0028 does not have (0028's contracts are
  a single `C1` with sub-clauses C1.1–C1.6; 0028's own JC1 names the arm as
  "C1/C1.4"). Rewritten to `cli/0028:C1 owns the `edit` arm and the tail
  categories (its C1.4 clause)`. Re-linted: exit 0, zero unresolved edges.
  `edge:unresolved-terminal`: none.

## Other pre-Gate mechanical checks

- **Cluster re-entry note.** No `## Refinement Context (cluster re-entry)`
  block in the record. `status_reentry=false`, `reentry_target=none`. Clean.
- **Repeatability-lite.** `--outcome repeatability` → `resolve:determinacy`;
  chained once → `emit.next: none`, rule `determinacy-na` (the record's
  written `Determinacy: n/a` line in Normative Contracts). Not owed.
- **Lens row.** Profile `large`; `--outcome lens` → `emit.row: none`, rule
  `lens-large-row-complete` (grounding → 3amigo → critique all present;
  `critique_models=differ`, so the cross-model pass is real, not a
  single-model fallback). `--outcome floor` → `emit.floor: none`
  (`seam_lineage_count=0`).

## Joint-decision fence (pre-lock, not a Tooling CHECK)

`--outcome fence` first emitted `stopped:overlap-uncited`.
`index --literal-intersect --record 0027` named the pair: **0022 ↔ 0027**
sharing the contract literal `docs/cli-reference.md`, neither citing the other.
0027's Joint-check line had recorded arm 2 clear at the batch propose of
2026-08-31 — before 0022 came in-flight — so the clear was stale, not wrong.

Fired arm 2 against 0022 rather than syncing any copy. Disposition: **no shared
decision**. `docs/cli-reference.md` is generated wholesale from the live command
tree by `internal/cli/docs.go::runDocs` (`make docs`, staleness-gated in
`make check`), so neither record hand-edits it and there is no format, ordering
or exclusivity claim to contend over. The two extend **disjoint enums under
different command sections**: 0022:C4 appends a code to graphlint's blocking
taxonomy (`internal/graphlint/taxonomy.go::BlockingCodes`), rendered by
`internal/cli/lint.go::lintExtendedDesc` under the `lint` section; 0027 Phase 3
adds description text for `table.Category`
(`internal/table/category.go::Categories`), a load-failure discriminator whose
sole consumer is `internal/cli/flow_input.go`, rendered under a `flow` section.
0022's "no new CLI surface" claim is about the finding taxonomy it extends and
is untouched by 0027 adding a body one level down on a different enum; 0027's S6
holds `Categories()` membership and order unchanged either way. No citation is
owed in either direction beyond recording the fire.

Recorded in the record at `0027:JC1` (line 202) — the fired list now reads
`fired → 0028, 0022` and carries the disposition. Re-run: `--outcome fence` →
`op = none`, rule `fence-clear`; `overlap_uncited=0`;
`joint_check_home=homed`; `open_joint_decisions=[]`.

## Verdict

**PASS** — 0 blocking findings after the one C10 in-pass fix (and the
joint-decision fence cleared by firing arm 2, above — a pre-lock gate, not a
Tooling CHECK). Proceed to the Gate's written responses.
