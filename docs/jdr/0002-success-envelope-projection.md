---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: open
cluster: 0005, 0011, 0017, 0021, 0023, 0024
labels: intrastate, cli-envelope, success-payload, projection, rdr-cluster
---

# JDR 0002 What the Success Envelope Carries, and How a Caller May Reshape It

*Flag spellings and Go identifiers below are **non-normative**; the RDRs own
exact contracts. Source: the 0023×0024 Stage-2 joint check, 2026-08-28.*

## Problem statement

Two verbs already couple on the resolve success-payload shape: 0023 projects
the echo group off `flow resolve`'s payload under an opt-in flag, and 0024
appends `dispositions` to the same payload at the same anchor
(`internal/cli/flow_resolve.go::resolvePayload`). Their Stage-2 joint check
fired on exactly that seam. Every future envelope change re-encounters the
same rule — which fields a caller may drop, and which no mode may ever omit —
and answered per-RDR the rule gets restated and drifts. This registry is the
single home for what the CLI success envelope carries and how callers may
reshape it: the 0005-envelope seam and its successors. 0005 owns the envelope
and the respond gateway the rule constrains; 0011 owns the flag-surface
discipline its instances register under; 0017 (finding surfacing) and 0021
(export surface) are prospective consumers.

## Principles

1. **One home per contract; cite, never restate** — the engine doctrine
   (`$RDR_HOME/stages/README.md`, *Single-source each contract*); why this
   registry exists.
2. **A flag may narrow the report, never the decision** — 0023:C1:
   report-only is oracle-enforced, not aspirational.
3. **Over-reporting is recoverable; silent omission is not** — the tiebreak
   for any field that fits no group cleanly (§D1).
4. **Partitions are enforced by structure, not prose** — an unassigned field
   is a test failure, never a silent default (§D1).

## D1 — The caller-projection partition doctrine

**Decided** — settled 2026-08-28 at the 0023/0024 Stage-2 joint check;
hoisted here at decision level from where 0023's proposer first recorded it
(`cli/0023:C2`, now the verb's instance record). Instance assignments —
concrete field lists, flag semantics, oracle spellings — live in the RDRs.

- Every current and future success-payload field of a projecting verb is
  assigned, when it is added, to exactly one of two groups: **ECHO** — what
  the caller supplied verbatim or can reconstruct from its own prior calls —
  or **PLAN** — what the run decided, its provenance, and what a caller acts
  on or must not misread.
- The partition is **enforced, not prose**: a reflective oracle makes an
  unassigned new field a test failure, never a silent default — a projection
  is never implemented as a bare omit-list whose complement is "whatever
  else exists" (P4).
- An **always-keep core** survives every current and future projection mode
  of the verb — for `flow resolve`: `rule`, `escaped`, `escape_class`,
  `revision`, as that verb's instance of the rule — so a projected payload
  can never launder a rescued plan into an ordinary one or detach a plan
  from the model revision that produced it.
- A field that fits neither group's description cleanly joins **PLAN** —
  over-reporting is recoverable, silent omission is not (P3).
- A **later verb** adopting a caller-controlled success-payload projection
  follows the same two rules: the projection is applied to the verb-specific
  result before `respond.OK` — never in the respond gateway, never differing
  by output mode (P2, and 0005's two-modes agreement) — and the projectable
  set is that verb's echo group, never its decision. No verb is obliged to
  offer a projection, and no flag spelling is reserved.

**Resolved:** lands in **0023** — `flow resolve`'s concrete ECHO/PLAN
assignments, `--plan-only` semantics, its reflective and always-keep oracle
obligations (`cli/0023:C2` cites this entry); and in **0024** —
`dispositions` is a PLAN-group field (its C4 join reads plan-group inputs
only), consistent by citation. 0005, 0011, 0017, 0021 are bound as owners or
consumers of the surfaces the rule constrains; nothing lands in them today.

## What this does not decide

- Which `flow resolve` fields sit in which group, and the `--plan-only`
  wire/byte contract — 0023's instance (`cli/0023:C1`/`C2`).
- The `dispositions` derivation, presence rule, and key order — `0024:C4`.
- Whether any later verb offers a projection, or what flag spelling it picks.
