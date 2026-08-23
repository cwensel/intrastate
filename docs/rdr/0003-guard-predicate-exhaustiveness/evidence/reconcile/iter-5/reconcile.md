Model: claude-opus-5[1m]

# Reconcile — RDR 0003 Guard Predicate Exhaustiveness (iteration 5)

## Verdict

**NOT RECONCILED — lock is blocked on A16.**

The route is **not** a stage return inside this RDR. Its approach, contract
wording, and evidence discipline are not what is defective: the grammar, the
identity tuple, the `unless` semantics, the assignment-count table, the coverage
derivation, and the projection clause are stated, internally consistent, and
computable. What blocks lock is one additive field in a `Draft` peer.

## Preflight — Stage 5 completeness

`Profile: large` → lens row `grounding → 3amigo → critique`, **plus
`repeatability` (lite) because the Determinacy trigger fires**.

| Lens | Required | On disk (iteration 3) | State |
| --- | --- | --- | --- |
| grounding | yes (`large`) | `evidence/grounding/iter-3/{findings,dispositions,Charted}.md` | complete |
| 3amigo | yes (`large`) | `evidence/3amigo/iter-3/{persona-1..3,consolidation,dispositions}.md` | complete |
| critique | yes (`large`) | `evidence/critique/iter-3/{critique,critique-modelB,diff,dispositions,Charted}.md` — two distinct stamps, diffed behind a barrier | complete (dual-model) |
| repeatability (lite) | yes — Determinacy fired | `evidence/repeatability/iter-3/{run-1,diff,dispositions}.md` | complete |

**Variant check passes**: `run-1.md` reads `variant: lite (profile: large)` and
no `run-2`/`run-3` sits beside it — lite was owed and lite ran.

**Preflight: PASS.** The open set was built.

## Absorption audit (delegated)

One sub-agent over the four iteration-3 lens output dirs. **Every finding is
absorbed**; three residues survive and all three are deliberately charted, not
dropped:

- **G3-2** — JDR 0001 §JD-4's line anchors into this RDR point at the wrong
  clauses. Peer document; repair at the next JDR 0001 touch. Not lock-blocking:
  the clause §JD-4 assigns here is present and correct, only the pointer is off.
- **R-10** — RDR 0006's default-on rider makes a docs-only domain declaration a
  CI break. The rider is that peer's decided text; raise at its refine.
- **R-18** — the declaration model shipped without an alternatives pass (A11's
  only rejected alternative is a *placement*, not a design). Suggested successor:
  an RDR seed for enum-domain evolution/versioning.

## The open set (four sources)

- **Source 1 — Pre-Lock needs-verification lists.** grounding: none new.
  3amigo: A16, A17 (new), A12 (narrowed). critique: A18, A19, A20 (new).
  repeatability: A21 (new), A16 (consequence raised to blocking), A17 (unchanged).
- **Source 2 — Pending/Unverified Critical Assumptions.** Census recomputed
  independently: **21 records, 11 `Verified`, 10 `Pending`** — A10, A12, A14,
  A15, A16, A17, A18, A19, A20, A21. None `Unverified`.
- **Source 3 — named-but-unrun spikes.** A1's evaluation harness has captured
  output (`evidence/spikes/iter-2/a1-eval-harness/output.txt`, 28 assertions
  pass / 0 fail) and the older `check.sh` transcript is retained with its scope
  correctly demoted. The remaining named runs (A15's Scenario 3 measurement,
  A21's Scenario 4 negative control) are implementation-time MVV work with named
  plans, not draft-time spikes. Desk computations (D-1's 8,192 recomputation,
  the census extraction, the conformance greps) left no artifact but are written
  into the RDR body and were re-verified here against source.
- **Source 4 — exactness-word delta.** The iteration-3 passes introduced or
  amended: the assignment-count table's declaration defaults; the
  single-value-operator projection clause; the `unless` Kleene clause; the
  cardinality-undefined-over-unprovable-dimension clause; the escape-row
  failure-class partition; the narrowing's widened population; the conformance
  definition; the guard-only participation rule; `{min..max}` inclusivity; the
  five pinned kind tokens. Each is carried by a record or grounded at its
  authority. The one that changed a lock verdict is the projection clause.

## Dispositions

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| **A16** — RDR 0002's authoring surface carries the single-valued marker | 1, 2 | **BLOCKER** | Confirmed against source: `single-valued`/`single_valued` occurs **zero** times in `docs/rdr/0002-transition-table-as-reviewable-data.md`; its schema list (`0002:215-222`) and normative type-model clause (`0002:340-352`) each enumerate the same four fields. Route: `/rdr-refine 0002` — add the marker to the `[tags.<tag>]` enumeration and the type-model field list, travelling with A14; then re-enter `/rdr-reconcile 0003`. |
| A10 — RDR 0006 accepts the row-group division of labour | 2 | DOWNGRADED | The cycle is broken *in this document* (the group is defined here). Peer confirmation scheduled: RDR 0006 is `Draft` and its Refinement Context Direction item 2 carries it. |
| A12 — predecessor reachability is decidable | 1, 2 | DOWNGRADED | Narrowed by 3amigo iter-3 to the owned-tag clause alone. RDR 0006's Direction item 3 carries it; the demotion already removed the reachability quantifier from assertion. |
| A14 — normalization preserves each atom's `block` | 2 | DOWNGRADED | Caught by the first normalization test, now named MVV Scenario 7. No peer asserts the opposite. RDR 0002's refine, with A16. |
| A15 — a single published cardinality bound is sufficient | 2 | DOWNGRADED | MVV Scenario 3, which measures whether the proof representation completes per shape and compares that to the bound's verdict — predicted vs. observed provability. |
| **A17** — escape-row overlap is checked in two populations | 1, 2, 4 | DOWNGRADED | Mechanics settled on shipped code (`internal/resolve/resolve.go::Row.rescues`, membership match). Only the cluster's confirmation of the two-population reading is owed. RDR 0006's refine (which owes the §JD-14 invariant 3/4 repair) + a §JD-14 correction. |
| **A18** — a conforming-view check has a producer | 1, 2, 4 | DOWNGRADED | Confirmed against source: `resolve.go::assemble` reads no declaration and `conform` occurs nowhere in the non-test kernel. The premise is stated as an explicit conditional, so nothing here asserts an unciteable enforcement. RDR 0007's next touch and/or RDR 0006's refine — either discharges it. |
| **A19** — a bare escape row's coverage closure is observable | 1, 2 | DOWNGRADED | Observability erosion, not a false green. RDR 0006 owns the verdict shape; its refine, with A17. |
| **A20** — atom-naming gets a carrier on both surfaces | 1, 2 | DOWNGRADED | Confirmed against source: `resolve.go::Refusal` carries one opaque `Guard string` beside `Rows []RowRef`. Lint half is §JD-4's decided duty; runtime half rides RDR 0007's payload reshape. Scenario 8's other assertions stand without it. |
| **A21** — requiring `single_valued` leaves the target flows authorable | 1, 2, 4 | DOWNGRADED | MVV Scenario 4's negative control + Phase 3 declaring the marker on the six named dimensions. Measurable only after A16 lands. |

## Why lock is blocked — the second hard rule

> no spike or assumption deferred past lock if the MVV depends on it.

A16 was booked as a scheduling dependency ("gates MVV Scenario 5's authoring,
not this RDR's contract") **before** the projection clause existed. The
repeatability iteration-3 pass then made `eq`, `in` and the integer comparisons
**single-value operators**: a value atom over a dimension the model does not
declare single-valued "cannot be projected at all" and MUST take the blocking
inability-to-prove outcome. That converts the marker from a product-size
optimization into a precondition for projecting any value atom.

The consequence the earlier scoping missed is that **the MVV's green cases go
with it**. MVV Scenario 2 requires "one complete partition" that *passes* — the
single positive exhaustiveness verdict the MVV exists to produce — and Scenario
8's paired negative control requires a group that "certifies green". Neither is
authorable while the marker cannot be written, so the MVV degrades to refusals
only and can no longer distinguish a tight narrowing from a blanket one, which
is the failure mode it was written to catch. That is an assumption pinning the
property under proof, not a fixture convenience.

The aggregate assessment's own test — "what would change the verdict from
*sequence carefully* to *do not lock* is any of them turning out to require a
*different contract here* rather than a field elsewhere" — is the wrong test,
and A16 is the case that proves it: A16 wants only a field elsewhere and still
empties the MVV. The corrected test, written into the RDR, is the hard rule's:
does the missing input pin the property the MVV proves.

**Cheapness argues for running it, not deferring it.** RDR 0002 is `Draft`, the
request is one additive field, no peer asserts the opposite, and its carriage
clause already covers the field generically.

## Strong consult (rdr-common §strong-consult)

One fresh-context sub-agent, given the clauses in tension and both directions
argued, without prior verdicts. Returned **BLOCK / blocking: yes**, reaching the
same discriminator independently (Scenario 2's passing partition is the MVV's
load-bearing green case and is unauthorable without the marker), and separately
flagged the internal contradiction between A16's own "now blocking" text and the
checklist's stale "Not lock-blocking". No tiebreaker escalated to the user: the
consult and this pass agree, and the route is a peer edit rather than a design
fork.

## Refutation check

**No spike or source search refuted an assumption this RDR relies on.** The
three source searches run this pass (RDR 0002's field census, `resolve.go::assemble`,
`resolve.go::Refusal`) each *confirmed* a gap the draft already books rather
than contradicting a claim it makes. A16 is a missing producer, not a refutation
— which is why the route is a peer refine and not a return to Propose or Refine.

A §punt-ledger row **is** owed and was appended: this gate reopens Stage 5's
obligation to sweep a clause's consequences across the MVV authorability and
gate surfaces, not only onto the assumption record.
Ledger: `docs/rdr/0003-guard-predicate-exhaustiveness-postmortem.md`.

## Completeness check

- No `_Draft placeholder._` in any body section — grep clean.
- No `this is a seed skeleton` header — grep clean.
- No surviving template brackets (`[Required`, `[Conditional`, `[Resource]`,
  `[Capability]`) — grep clean.
- `## References` carries real citations (peer RDRs, JDR 0001 §D1/§D4/§JD-4,
  the output contract, the prior-art corpus), no template placeholders.

## Dispositions landed in the RDR

Every row above is written into the RDR's Critical Assumptions section, not only
here. A16 carries a **Stage 6 disposition — BLOCKER** paragraph; A17, A18, A19,
A20 and A21 each gained a **Stage 6 disposition — DOWNGRADED** paragraph naming
the venue, the survivability argument, and the plan; A10, A12, A14 and A15
carried theirs from prior iterations unchanged.

**Amendment sweep (rdr-common §amendment-sweep)** for the A16 consequence:

- `Not lock-blocking` — 3 sites. A16's checklist entry amended (it asserted the
  pre-projection-clause scoping); A17's and A14's re-read and left standing,
  both still correct.
- `no record now gates lock on a peer decision` — 1 site (Assumption
  Verification, the A8 paragraph), amended to scope its claim to §JD-4.
- **Authorability paragraph** — amended: its "Every other case is authorable
  today" predated the projection clause and is now false for every green-verdict
  case. Now names Scenario 2's partition and Scenario 8's negative control.
- **Aggregate open-assumption assessment** — count corrected (nine of twenty →
  ten of twenty-one), A21 added to the third group, A16 marked as blocking lock
  rather than implementation, and the section's own lock test replaced with the
  hard rule's.
- **Assumption Verification roster** — `Pending` list was missing A21; corrected
  and the census stated explicitly (21 records, 11 Verified, 10 Pending).
- **Prerequisites checklist** — was short by four records; A18, A19, A20 and A21
  entries added with route and done-condition.
- **Metadata `Status`** — records the iteration-5 verdict and route; the
  "Remaining open" list, which named six records, now names all ten.

## Next

```
/rdr-refine 0002        # land the single-valued marker (with A14's block-retention clause)
/rdr-reconcile 0003     # re-entry once RDR 0002's authoring schema carries the field
```

RDR 0006's refine (A10, A12, A17, A19, A20's lint half) and RDR 0007's next
touch (A18, A20's runtime half) are independent of the block and can proceed in
any order; none of them gates this RDR's lock.
