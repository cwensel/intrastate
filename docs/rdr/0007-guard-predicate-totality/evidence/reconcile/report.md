Model: claude-opus-5[1m]

# Stage 6 Reconcile — RDR 0007 guard predicate totality

Forced gate after all Stage 5 rounds. Every open spike and every
round-disturbed assumption driven to a terminal state, written into the RDR's
Critical Assumptions section (the RDR is the artifact; this file is the audit
trail).

## Stage 5 completeness preflight — PASS

`Profile: foundational` → required lenses `cove 3amigo critique repeatability`;
all four evidence dirs present and resolved.

- `repeatability/` ran the **correct variant**: `variant: full (profile:
  foundational)` with three runs on three distinct base models
  (`claude-opus-5[1m]`, `claude-fable-5`, `glm-5.2:cloud`) plus `diff.md`. No
  lite/full mismatch.
- `critique/` satisfied its `foundational` dual-model obligation: pass A
  (`claude-opus-5[1m]`) + pass B (`claude-fable-5`), diffed by passage anchor in
  `diff-modelB.md`.
- Determinacy trigger: n/a as a separate obligation — the profile's full
  repeatability lens subsumes it.

## Absorption audit (delegated)

All 55 lens findings across four rounds verified absorbed into the current RDR
text or dispositioned with a cite that holds. **No text residue.**

| Round | Findings | Absorbed | Residue |
|---|---|---|---|
| 3amigo | PM/IMP/QA ×15 | 15/15 | A10 |
| critique pass A | C-1..C-15 | 15/15 | A12, A13; C-14 charted |
| critique pass B | B-1..B-15 | 15/15 | A15 (net-new catch) |
| repeatability | D-1..D-10, G-1..G-8 | 18/18 | A14 |
| cove | F-1..F-9 | 9/9 **verified directly** | F-8 charted → RDR 0008 |

Cove was the risk: it has no `dispositions.md`, so every F-N was checked against
the RDR text rather than a ledger. All nine landed (F-1 → assumption A8; F-4 →
Infrastructure Audit split with two rows marked Build; F-5 → empty-block
normative clause + Scenario 9).

## Open set and terminal dispositions

Sources: (1) Pre-Lock needs-verification lists ×4, (2) still-Pending
assumptions, (3) named-but-unrun spikes, (4) exactness-word delta.

| Item | Src | Disposition | Evidence / plan |
|---|---|---|---|
| A5 two-row absence pattern | 1 | **VERIFIED** (runtime settled; load-time rides on A12's obligation) | Prior "load-time half is *refuted*" was too strong — RDR 0003 is silent, not adverse. Narrowed by A14: two-row pattern needed only for the disjunction |
| A9 empty-`unless` identity | 2 | **VERIFIED** | RDR 0002 silent (governs representation + combination timing only); it fixes the analogous absence-identity for *writes*, showing intent. RDR 0003's subtractive algebra independently forces the same reading — the vacuously-true reading would empty every row's accepted set and break coverage for every row group |
| A10 guard-structure mapping (hinge) | 1+2 | **VERIFIED as NAMED OBLIGATION** | Specified nowhere in RDR 0001/0002/0003 or code (demonstrated negative sweeps). *Not refuted*: RDR 0003's Identity decision is already a total string-expressible key fitting `Row.Guard string` losslessly → RDR 0001's `Row` stays closed. Destination: **RDR 0003 `Phase 1: Predicate Model`**, which already charters "the normalized predicate atom shape used by resolver and lint" |
| A12 load-time overlap | 1 | **VERIFIED as INDETERMINATE-BY-SILENCE** → inherited obligation | RDR 0003 never states how an existence atom projects onto the declared-domain product; "optional scalar" confers no domain standing (RDR 0002's tag schema has no optionality field). Overlap/exhaustiveness asymmetry **confirmed**. Not rejected — the projection rule does not exist. Destination: RDR 0003 implement (projection) + RDR 0006 implement (row-group membership) |
| A13 caller-supplied observed | 1 | **VERIFIED as ACCEPTED EXPOSURE** | No constraint anywhere: RDR 0005 never mentions observed tags and treats `--tag` as "context already known to the caller"; RDR 0004 constrains only writes; RDR 0001 normatively *affirms* "caller-supplied observed tags". Bounded: no production code imports `internal/resolve`; `root.go` registers only `version`. Seeded against RDR 0005 |
| A14 `exists` boolean literal | 1 | **VERIFIED** | RDR 0003 never defines "positive"/"negative" atom — polarity is placement-only and structurally independent of the literal; `false` parses as boolean and is not rejectable. Fixture's `exists = true` is under-coverage (validator matches operator *names* only) |
| A15 conjoined value+existence row | 1 | **VERIFIED as NAMED OBLIGATION** (silence, not prohibition) | RDR 0002 names the table path but bounds no operator keys inside it. Structure leans toward admitting: identity tuple is `(tag, operator, literal)`; selection quantifies over *atoms*; both RDRs' own TOML puts multiple atoms under one guard header; spike validator has no at-most-one constraint. Destination: same RDR 0003 Phase 1 step as A10 |
| A6b read-failed vs absent | 2 | **DOWNGRADED** (by decision) + **ROUTED** | Not MVV-critical (the kernel seam cannot express the distinction either way); contract sound without it. Obligation routed to **RDR 0004's implement stage**. Survivable: exposure is that an unevaluable refusal must not be read as "retryable" — stated as operator guidance in Failure Modes |
| Aggregation spike (Probes A/B) | 3 | **ALREADY CLOSED** at Stage 4 | `evidence/spikes/aggregation-probe.md` has command + verbatim output. Probe A **refuted** the escape half of the aggregation clause; correctly resolved by changing the RDR text to shipped behavior, not the kernel — Phase 1 stays doc-comments-only |
| Exactness-word delta | 4 | **CLEAN** | Sweep of all `normative` blocks: every all/every/never/only/total/sole/exactly claim traces to a named assumption or an evidenced source pointer. No orphaned post-mutation claim |

## Hard rules

**Refutation → BLOCKER**: none live. The one refutation in this RDR's history
(spike Probe A, escape-half of the aggregation clause) was resolved at Stage 4
by correcting the RDR to match frozen kernel behavior — the honest direction,
since three shipped tests and RDR 0001 deviation D5 freeze it. Nothing at Stage 6
refutes a claim the RDR relies on. Critically, A10's "If wrong" branch (every
atom-level clause unimplementable, `Row` reopened) is **not triggered**: the
mapping is unowned, not impossible.

**No MVV-critical deferral**: the MVV proves the masking probe inverts —
scenarios 1–3, all kernel-level, all reachable today against `fixtureGuards`
with no atom structure required. None of the deferred obligations (A10, A12,
A15) gates them; they gate Testing Strategy rows 4–8, which the RDR already
declares as Phase 2 work explicitly blocked on A10. The RDR states plainly that
a green MVV must not be read as validating the domain rule.

## Completeness check

No `_Draft placeholder._`, no seed-skeleton header, no bracketed template
placeholders. `## References` carries 12 real citations. Method vocabulary
clean — no `Docs Only` on any assumption.

## Carried to Stage 7 (not Stage 6's call)

**Proportionality / contract count.** Both critique passes independently
(C-14, B-15) counted five or more independently load-bearing contracts against a
`Profile` of "foundational — one contract", and charted the question to
`/rdr-finalize`'s Proportionality response — which is still an unfilled
template. Finalize must either justify the value in writing or correct the field
and split, with the `RequiresOwned` narrowing named as the likeliest seam to
lift out. Recorded here so finalize does not meet it cold.

## Verdict

**RECONCILED.** All items terminal; no BLOCKER. Three dispositions are
obligations inherited from Final peers (A10, A12, A15) and one is a
by-decision downgrade (A6b) — each with a destination that exists, per this
RDR's own never-amend routing rule. Ready for Finalize.
