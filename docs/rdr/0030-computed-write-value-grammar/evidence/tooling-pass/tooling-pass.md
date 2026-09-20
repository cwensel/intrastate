Model: claude-opus-5[1m]

# Tooling Pass — cli/0030 computed-write-value-grammar

Run: 2026-09-20 · Stage 7 mechanical pre-sweep · Profile `foundational`
Lint receipt: `lint.txt` (`recs lint --locking 0030`)

Regression sweep run after the pre-lock lenses and the Stage 6 reconcile
rewrote this draft.

## Findings

- **C9 [BLOCKING, FIXED IN-PASS] `0030:A12`** — Evidence field 32 lines
  against TEMPLATE.md's one-sentence-plus-stable-anchor form. A
  `foundational` record does not lock over this (`lint --locking` marks it
  blocking; the tool's word, not an eye judgement). Judged per the check's
  one question: the mass IS real verification content and every
  load-bearing anchor is a `path::Symbol`, so truncation was refused.
  **Relocated, not truncated** — the full narrative is reproduced verbatim
  in `artifacts/evidence-a12-a15.md` §A12 and the field keeps the claim,
  the anchors (`internal/guard/grammar.go`,
  `internal/guard/declaration.go::intWidth`, `::IntDomain`,
  `internal/resolve/resolve.go::GuardEvaluator`,
  `guardcontract.go::TestGuardEvaluatorContract`) and the pointer.
- **C9 [BLOCKING, FIXED IN-PASS] `0030:A15`** — the same, 48 lines. Same
  judgement, same disposition: relocated verbatim to
  `artifacts/evidence-a12-a15.md` §A15; the field keeps the
  eight-producer / two-consumer result, `coverage.go::groupHasOverlap`,
  `internal/graphlint/engine.go::identityKey`, `export.go::compareEdges`
  and `normalize.go::normalizeRules`, plus the pointer. "Flagged,
  accepted" was explicitly NOT taken as a disposition.
- **C1 [FIXED IN-PASS] `0030:§cross-cutting-concerns`** — the gate item
  retained at lock carried only TEMPLATE.md's guidance block. Authored in
  place this pass (canonical-form/determinism, versioning, incremental
  adoption, memory/row growth, character encoding; the six inapplicable
  concerns omitted rather than N/A-bulleted), so the element peers cite as
  `cli/0030:G-cross-cutting` is real content.
- **C1 [NOT A DEFECT] four surviving guidance blocks** at
  `§contradiction-check`, `§assumption-verification`, `§scope-verification`
  and `§proportionality`, plus `gate:inline` over the gate section. These
  four are the sub-sections the lock moves out to `artifacts/gate.md`
  behind the one-line pointer; the blocks vanish with the move. Not routed
  to Refine and not hand-deleted.
- **C1** — template spine complete. Zero `template:missing-section`, zero
  `scaffold:row`, zero `contract:template-example`. No TBD, no "see above",
  no single-sentence stand-in, no seed-skeleton header. No
  `## Refinement Context (cluster re-entry)` block exists.
- **C2** — PASS. `ca_off_vocabulary=0`; every one of the 16 assumptions
  carries a Method drawn from the sanctioned set (`Source Search` ×9,
  `MVV Test` ×5, `Spike` ×1, `Peer RDR` ×1). No Evidence Record lacks a
  Method field.
- **C3** — PASS, nothing to report. Of the nine `Source Search` records,
  none resolves its Evidence to the record itself or to a path under this
  RDR's artifact directory. (The check has not fired since the structured
  Evidence Record landed; this run does not change that.)
- **C4** — PASS, vacuous. No assumption carries `method.members ==
  ["Docs Only"]`, so there is no load-bearing Docs-Only claim to pair with
  a Spike or Source Search plan.
- **C5** — PASS. 94 `source-anchor` edges, all `resolved: true`. No
  ABSENT (the repo root resolved, so the check genuinely ran), no bare
  `file:line` without a symbol.
- **C6** — PASS. Five assumptions are `Pending` (A5, A6, A10, A11, A13),
  all `MVV Test`, all unrunnable before Phase 2 by construction and each
  pinned to a named authored scenario (S1/S5b/S5c, S2/S2b, S8, S10,
  S5). A delegated sweep of C1, C2, C3, `§load-bearing-decisions`,
  `§consequences`, `§failure-modes` and `§decision-rationale` found every
  site stating one of those five properties states it NORMATIVELY — what
  the implementation MUST do — never as an observed settled fact; the
  strongest candidate (`§approach`'s "everything after normalization … is
  unchanged") is explicitly hedged in `§consequences` ("the MVV compares
  plans and findings, not dumps … a TEST oracle, not a product surface").
  No checklist box disagrees with the gate. One `Refuted` record (A8) is
  terminal and its refutation is absorbed into C3 and the Decision
  Rationale per ruling Q1.
- **C10** — PASS. Three labelled contracts (C1, C2, C3), no
  `label:contracts`. No `peer-evidence:no-element`: the one `Peer RDR`
  record (A9) cites an element. No `edge:unresolved` and no
  `edge:unresolved-terminal` — predecessors (0002, 0012), overrides
  (`0002:C4`, `0002:C13`), the three joint-decision homes and the
  surface-of pair all resolve.

## Verdict

**PASS** — 0 blocking findings after the two in-pass C9 relocations;
`recs lint --locking 0030` exits 0 (`blocking=0 resolution=0`). The four
residual `placeholder:survived` and the `gate:inline` advisory are
discharged by the lock's own section move, not by an author edit.

First pass; no `$PRIOR_DIR`, so the loop-breaker diff is skipped by its
own rule.
