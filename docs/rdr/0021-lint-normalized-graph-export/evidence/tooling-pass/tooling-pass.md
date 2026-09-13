Model: claude-opus-5

# Tooling Pass — RDR 0021 (lint's normalized-graph export)

Date: 2026-09-12 · Iteration 1 (no prior pass; loop-breaker diff skipped)
Lint: `rdr lint --locking 0021` → exit 0, `blocking=0 resolution=0 placeholder=5 advisory=12`
(captured at `lint.txt` beside this file).

## CHECK 1 — Template section coverage

No `template:missing-section` findings; the spine is complete (Metadata, Problem
Statement, Critical Assumptions, Proposed Solution, Alternatives, Context,
Research Findings, Trade-offs, Implementation Plan, Validation, Finalization
Gate, References all present and authored).

Hollow sections — all five are the Finalization Gate's own sub-sections, which
this stage authors:

- C1 `placeholder:survived` 1263-1266 — §contradiction-check is template brackets only.
- C1 `placeholder:survived` 1272-1284 — §assumption-verification is template brackets only.
- C1 `placeholder:survived` 1290-1292 — §scope-verification is template brackets only.
- C1 `placeholder:survived` 1303-1307 — §cross-cutting-concerns is template brackets only.
- C1 `placeholder:survived` 1328-1353 — §proportionality is template brackets only.
- C1 `gate:inline` 1233-1353 — gate responses inlined; they belong in `artifacts/gate.md`.

These are MECHANICAL and are discharged by this stage writing the gate responses
(items 1/2/3/5 to `artifacts/gate.md`, item 4 authored in the record). They are
not routed to Refine. No section is hollow outside the gate; no seed-skeleton
header survives.

## CHECK 2 — Method label vocabulary

PASS. `ca_off_vocabulary=0`, `ca_off_vocabulary_ids=[]`. All nine assumptions
carry a Method field drawn from the sanctioned set (Spike, Source Search,
Peer RDR).

## CHECK 3 — Source Search self-reference

PASS. The `Source Search` rows (A3, A9) anchor into `internal/graphlint/`,
`internal/guard/` and `internal/table/` on `main` — no Evidence path resolves to
the record itself or to its artifact directory.

## CHECK 4 — Docs Only on load-bearing claims

PASS. No assumption carries `Method: Docs Only`.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

PASS. No `edge:unresolved`. Every `source-anchor` edge in the projection
reports `resolved: true` (A1 `respond.go::TextLiner`, A2/A3 `reach.go::reach`/
`::Reach`, A3 `guard/declaration.go::AssignmentCount`, `guard/product.go::Groups`,
A9 `reach.go::OpaqueValue` and the `table/` authored-value sites).

## CHECK 6 — Status consistency

One Pending assumption, A9, and it is consistent: the record states the
downgrade explicitly ("Pending — DOWNGRADED at Stage 6 iteration 2,
survivable"), and C2 was refined so no clause depends on the refuted limb —
the discriminator is conditioned on the owed `0002:C11` reservation and cites
JDR 0001 §JD-23 rather than restating it. No settled-fact prose elsewhere
relies on the refuted limb. Not a C6 finding.

## CHECK 9 — Evidence-field budget (ADVISORY)

- C9 `evidence:over-budget` 330-372 — A9's Evidence field is 43 lines (soft cap 30).
  Profile is `large`, not `foundational`, so this is advisory, not blocking. The
  load-bearing anchors (`reach.go::OpaqueValue`, `reach.go:438`,
  `table/model.go::ClearSentinel`, `analysis.go::nodeMeetsAll`) remain findable
  in the field; the balance is the two-limb REFUTED/VERIFIED disposition the
  Stage-6 reconcile produced and the grounding sweep reads. Answered at the
  Gate: keep as-is, do not truncate.

## CHECK 10 — Linking

- Advisory `prose:exactness` ×5 (506, 509, 521, 546, 597) — "canonical",
  "byte-identical", "deterministic" in normative claims. Each is covered:
  A5 (encoder byte-stability, Spike), C3 (determinism), and Testing Strategy
  S1/S2/S5/S9 pin them with named fixtures and the MVV. No action.
- No `label:contracts` finding — C1–C5 are labelled.
- No `peer-evidence:no-element`; no `edge:unresolved`.

## Cluster re-entry note

None. No `## Refinement Context (cluster re-entry)` block survives in the record.

## Verdict

PASS — 0 blocking findings. The five hollow gate sub-sections are this stage's
own authoring obligation, not a regression to route back. Proceed to the Gate's
written responses.

Note: the sweep PASSES, but the joint-decision fence (run separately, after this
sweep) returns `stopped:overlap-uncited` — see the Gate report. That is a
Stage-7 blocker independent of this mechanical pass.
