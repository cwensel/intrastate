Model: claude-opus-5[1m]

# Reconcile — RDR 0003 Guard Predicate Exhaustiveness (iteration 6)

## Verdict

**RECONCILED.** The iteration-5 BLOCKER is discharged. Every open record holds a
terminal disposition, no record refutes a claim this RDR relies on, and no open
record pins the property the MVV proves.

This pass is the re-entry iteration 5 scheduled. Its route — `/rdr-refine 0002`
to land the single-valued marker travelling with A14 — was executed (commit
`c12cc19`, with RDR 0002 resolved on top at `3119744`). This pass verifies the
landing against source and disposes the open set accordingly.

## Preflight — Stage 5 completeness

`Profile: large` → lens row `grounding → 3amigo → critique`, **plus
`repeatability` (lite) because the Determinacy trigger fires**. No lens has run
since iteration 3; iteration 5 verified the same evidence and nothing on disk
changed.

| Lens | Required | On disk (iteration 3) | State |
| --- | --- | --- | --- |
| grounding | yes (`large`) | `evidence/grounding/iter-3/{findings,dispositions,Charted}.md` | complete |
| 3amigo | yes (`large`) | `evidence/3amigo/iter-3/{persona-1..3,consolidation,dispositions}.md` | complete |
| critique | yes (`large`) | `evidence/critique/iter-3/{critique,critique-modelB,diff,dispositions,Charted}.md` — two distinct stamps, diffed behind a barrier | complete (dual-model) |
| repeatability (lite) | yes — Determinacy fired | `evidence/repeatability/iter-3/{run-1,diff,dispositions}.md` | complete |

**Variant check passes**: `run-1.md` reads `variant: lite (profile: large)` and
no `run-2`/`run-3` sits beside it — lite was owed and lite ran.

**Preflight: PASS.**

## Absorption audit (delegated)

One sub-agent over the four iteration-3 lens output dirs. **All 43 findings are
absorbed**; every claimed fix site was verified present in the body. No
unabsorbed finding exists outside the three known residues, and none implies an
unrun spike or an unbooked assumption. The three residues are unchanged and
still correctly charted: **G3-2** (JDR 0001 §JD-4's line anchors into this RDR —
peer document repair, now further off after iteration-3 edits, which confirms
the repair belongs at the JDR touch), **R-10** (RDR 0006's default-on rider —
that peer's decided text), **R-18** (the declaration model shipped without an
alternatives pass — charted to a successor seed for enum-domain
evolution/versioning).

One bookkeeping note from the audit, not a finding: the critique's `Charted.md`
claims the aggregate assessment "names this as the thing to watch" for R-18, and
the iteration-5 rewrite of that section dropped the mention. R-18's disposition
is unchanged and its own `Charted.md` is the durable record.

## Drift check — RDR 0002's refine against this RDR's claims (delegated)

RDR 0002 changed by +181/-55 lines between `c12cc19~1` and `3119744`. A
sub-agent checked every RDR 0003 citation into RDR 0002 against the new text.

**No RDR 0003 claim is contradicted.** Every disturbance runs in this RDR's
favor. Findings:

- **A16 and A14 both landed** (details in Dispositions below).
- **A17's premise is INTACT.** Its cited clause was rewritten as gate-then-count,
  but preserves the sequencing verbatim in substance and is *more* explicit about
  per-failure-class escape matching than the text A17 quoted.
- **The rehoming is ratified bilaterally.** RDR 0002's refine replaced its
  non-normative prose schema list with a normative clause naming all five fields
  and stating "RDR 0003 is the normative home of that model … MUST NOT be
  restated here". The ownership-arms table's "contested arm" cell for the
  declaration model is now settled rather than contested.
- **RDR 0002's profile change (`large` → `foundational`) touches nothing here** —
  this RDR never asserts RDR 0002's profile or lens obligations. The only status
  it asserts is `Draft`, which still holds.
- **Anchor hygiene:** five distinct stale `0002:NNN-NNN` anchors (one used six
  times). Every cited clause is still present and still says what this RDR says
  it says, so these were anchor problems, not findings — repaired in this pass by
  converting to stable section anchors per rdr-common §delegation anchor
  doctrine. One was wrong in substance as well as stale (the `escape`-field
  citation pointed at the selection clause rather than the escape-list clause)
  and now cites the right clause with its wording.

## The open set (four sources)

- **Source 1 — Pre-Lock needs-verification lists.** Unchanged from iteration 5;
  no lens has run since. grounding: none new. 3amigo: A16, A17 (new), A12
  (narrowed). critique: A18, A19, A20 (new). repeatability: A21 (new), A16
  (consequence raised to blocking), A17 (unchanged).
- **Source 2 — Pending/Unverified Critical Assumptions.** Census recomputed
  independently at entry: **21 records, 11 `Verified`, 10 `Pending`** — A10, A12,
  A14, A15, A16, A17, A18, A19, A20, A21. None `Unverified`. After this pass:
  **21 records, 13 `Verified`, 8 `Pending`**.
- **Source 3 — named-but-unrun spikes.** Unchanged. The only `Method: Spike`
  record is A1, `Verified` with captured output
  (`evidence/spikes/iter-2/a1-eval-harness/output.txt`, 28 assertions pass / 0
  fail). A15 and A21 are `MVV Test` implementation-time plans, not draft-time
  spikes. No spike is named-but-unrun.
- **Source 4 — exactness-word delta.** No lens ran this iteration, so the delta
  is what RDR 0002's refine introduced into text this RDR cites: the
  gate-then-count selection ordering ("exactly one surviving non-escape
  normalized candidate row"; "exactly one escape row for that failure class"),
  the per-atom retention MUST and its "MUST NOT fold `unless` atoms into `all`",
  the "carry every declared field through … without loss" carriage clause, and
  the within-row total order "(key, block, operator token, literal)". Each was
  read against the claim it supports; all four corroborate rather than disturb.

## Dispositions

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| **A16** — RDR 0002's authoring surface carries the single-valued marker | 1, 2, 4 | **VERIFIED** (was BLOCKER) | Landed in both required surfaces, each citing this RDR for meaning: the `[tags.<tag>]` schema enumeration (`0002::Technical Design`, "Tag declarations") and the normative type-model clause (`0002::Normative Contracts`), which also states the authoring location normatively and carries the general carriage obligation. Done-condition met to the letter. Confirmed by a delegated verification pass and independently by the drift check. |
| **A14** — normalization preserves each atom's `block` | 2, 4 | **VERIFIED** (was DOWNGRADED) | `0002::Normative Contracts`: "each atom in that set MUST retain the key, operator token, literal, and the block (`all` or `unless`) it was authored in", plus "normalization MUST NOT fold `unless` atoms into `all`", naming this RDR's identity tuple and `unless` semantics as the reason. `block` also carried into the dump field list and the within-row sort key. Stronger than requested. MVV Scenario 7 keeps it honest. |
| A10 — RDR 0006 accepts the row-group division of labour | 2 | DOWNGRADED (carried) | Cycle broken in this document; peer confirmation on RDR 0006's refine (its Direction item 2). Unchanged. |
| A12 — predecessor reachability is decidable | 1, 2 | DOWNGRADED (carried) | Narrowed to the owned-tag clause alone; RDR 0006's Direction item 3. Unchanged. |
| A15 — a single published cardinality bound is sufficient | 2 | DOWNGRADED (carried) | MVV Scenario 3 — predicted vs. observed provability, refuted by divergence in either direction. Unchanged. |
| **A17** — escape-row overlap is checked in two populations | 1, 2, 4 | DOWNGRADED, re-affirmed | Mechanics settled on shipped code (`internal/resolve/resolve.go::Row.rescues`). RDR 0002's rewritten clause corroborates the two-population reading more explicitly than the text A17 quoted, but corroboration from a second document is not the cluster confirmation A17 books. Venue unchanged: RDR 0006's refine + the §JD-14 correction. |
| A18 — a conforming-view check has a producer | 1, 2, 4 | DOWNGRADED (carried) | Premise stated as an explicit conditional; RDR 0007's next touch and/or RDR 0006's refine. Unchanged. |
| A19 — a bare escape row's coverage closure is observable | 1, 2 | DOWNGRADED (carried) | Observability erosion, not a false green; RDR 0006's refine with A17. Unchanged. |
| A20 — atom-naming gets a carrier on both surfaces | 1, 2 | DOWNGRADED (carried) | Lint half is §JD-4's decided duty; runtime half rides RDR 0007's payload reshape. Unchanged. |
| **A21** — requiring `single_valued` leaves the target flows authorable | 1, 2, 4 | DOWNGRADED, sequencing updated | MVV Scenario 4's negative control + Phase 3 declaring the marker on the six named dimensions. Was gated on A16; **now measurable**, since the marker is writable. Still not lock-blocking. |

## Why lock is no longer blocked

Iteration 5 blocked on the second hard rule — no assumption deferred past lock
when the MVV depends on it. The dependency was real: the projection clause makes
`eq`/`in`/comparisons single-value operators whose atoms over an unmarked
dimension cannot be projected at all, so an unwritable marker left MVV
Scenario 2's complete partition (the MVV's only passing exhaustiveness verdict)
and Scenario 8's green negative control unauthorable, degrading the MVV to
refusals only.

The marker is now writable. Both green cases are authorable, and the MVV again
produces a positive proof outcome rather than refusals only. The hard rule is
satisfied on its own terms — the input that pinned the property the MVV proves
has been produced, not deferred.

The eight remaining `Pending` records were each tested against the same rule and
none pins the MVV's property: A10, A12, A17 and A19 are peer agreements about
this RDR's own definitions whose failure mode is a false positive a fixture
catches; A18 and A20 are producer fields whose absence weakens Scenario 8's
diagnostics without removing a verdict; A15 and A21 are measurements the MVV
performs rather than inputs it consumes. Each leaves a provable case standing.

No **§strong-consult** was run this pass and none is owed: there is no reopened
fork, no tiebreaker the evidence won't collapse, and no verdict flapping. The
iteration-5 consult ran on the fork this pass closes, and it closed in the
direction that consult predicted.

## Refutation check

**No spike or source search refuted an assumption this RDR relies on.** The
delegated verification and drift passes read RDR 0002's refined text against
every claim this RDR makes about it; all disturbances were favorable
(requests landed) or neutral (anchors moved, clause rewritten with meaning
preserved). No route-back is owed, so **no new §punt-ledger row is appended** —
the iteration-5 row recording the escaped consequence sweep stands as the
durable record of that punt.

## Completeness check

- No `_Draft placeholder._` in any body section — grep clean.
- No `this is a seed skeleton` header — grep clean.
- No surviving template brackets (`[Required`, `[Conditional`, `[Resource]`,
  `[Capability]`) — grep clean.
- `## References` carries real citations; its §JD-4 entry was stale ("recording
  document still open") and now records the 2026-08-22 closure, with §JD-13 and
  §JD-14 added.

## Dispositions landed in the RDR

Every row above is written into the RDR's Critical Assumptions section, not only
here. A16 and A14 were rewritten as `Verified` records with concrete Evidence
and a **Stage 6 disposition — VERIFIED (iteration 6)** paragraph; A17 gained a
**re-affirmed** paragraph recording the rewritten peer clause as corroboration
rather than confirmation; A21's sequencing note was updated to record that its
gate has lifted.

**Amendment sweep (rdr-common §amendment-sweep)** for the A16/A14 closures:

- `still enumerates four fields` / `not yet writable` / `occurs zero times` —
  every site amended: A16's own Evidence, the §JD-13 checklist item, the
  Authorability paragraph, and the two "Fields Requested Elsewhere" rows (both
  now read **Landed**).
- **Prerequisites checklist** — A16 and A14 entries checked `[x]` with the
  done-condition recorded as met.
- **Authorability paragraph** — the "one exception" is recorded as closed, and
  the claim that green-verdict cases wait on the field is replaced with the
  restored MVV outcome.
- **Aggregate open-assumption assessment** — recount (ten of twenty-one → eight
  of twenty-one), the producer-fields group reduced to A18/A20, and the section's
  lock test kept while its verdict flips: no open record blocks lock. The
  corrected test (does the missing input pin the property the MVV proves) is
  retained, since it is what made A16 blocking and what clears the other eight.
- **Assumption Verification roster** — Verified/Pending lists and census
  corrected to 13/8; the A14 paragraph rewritten from "Pending … single surviving
  producer request" to `Verified`; the A8 paragraph's parenthetical scoped
  correctly ("no record now gates lock on a peer edit").
- **Phase gating** — Phase 3's blocker cleared; the fixture's unprovable Draft
  group re-described as a fixture gap rather than a schema one.
- **Ownership arms table** — the declaration-model row's contested arm settled,
  recording RDR 0002's normative adoption of the rehoming.
- **Anchors** — six `0002:328-334` uses and one `0002-...:317` converted to
  stable section anchors; the `escape`-field citation repointed to the
  escape-list clause it actually depends on.
- **Metadata `Status`** — records the iteration-6 verdict, the discharge of the
  BLOCKER, and the corrected "Remaining open" list (eight records).

## Next

```
/rdr-finalize 0003
```

RDR 0006's refine (A10, A12, A17, A19, A20's lint half) and RDR 0007's next
touch (A18, A20's runtime half) remain independent of lock and can proceed in
any order.
