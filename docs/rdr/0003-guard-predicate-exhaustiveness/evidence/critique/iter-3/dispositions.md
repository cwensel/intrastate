Model: claude-opus-5[1m]

# Critique Dispositions — iteration 3 (re-entry, dual-model)

Origin ledger = `diff.md`'s reconciled ledger (R-01..R-20), built by a barrier
sub-agent that authored neither pass. `C-N` IDs do not correspond across
`critique.md` (Model: claude-opus-5, 17 rows) and `critique-modelB.md`
(Model: claude-sonnet-5, 10 rows); reconciliation was by passage anchor.

## Dual-model draw

Two distinct stamps, fanned out in parallel per rdr-common §auto-fanout, diffed
behind the barrier by a third context that authored neither. `large` profile:
dual-model is strongly recommended and was achieved; no single-model fallback
recorded.

**The disagreement was the signal, and it was asymmetric.** Pass A found five
BLOCKING computational defects that pass B did not raise at all — the uncomputed
coverage verdict, the `set` cardinality exponent, the Kleene `unless` mismatch,
and the `gate`/`escapeOrRefuse` control flow. Pass B's reading was
governance-weighted and its two strongest rows (R-05, R-06) were also found by A
from a different surface. Where they contradicted, the evidence favoured A both
times (see CT-1, CT-2 in `diff.md`). A single-model pass would have shipped four
false-greens.

## Grounding gate

Every actioned finding was verified against source before any edit. The four
blocking arithmetic/control-flow claims were verified **by execution or by
reading the control flow**, not by argument (COMPUTE-DON'T-ARGUE):

- **R-01 — executed.** The Draft group's product is
  `profile`(4) × `prelock_iterations`(4) × `cluster_eligible`(2) × presence(2)
  = **64**. `profile-to-grounding` accepts 24, `foundational-to-cove` accepts 8,
  union **32**, **uncovered 32**. `profile = small` is accepted by no row.
  Overlap is empty, so trace step 7's existing claim stands. **R-01 grounds —
  the trace stamped OK on a group that is half uncovered.**
- **R-02 — computed against the RDR's own clauses.** The bound clause said "the
  product of every participating dimension's *declared domain size*"; the kind
  table gives `set`'s domain as its element universe (size N). But the RDR's own
  single-valued clause already writes `2^|domain|` for co-occurring values, and
  the kind table says `set` values "co-occur by construction" — so `set` must be
  `2^N`. The document contradicted itself by an exponent. **R-02 grounds.**
- **R-03 — control flow read verbatim.** `resolve.go::Resolve`:
  `selected, blocked := gate(...)` then `if blocked != nil { return refuse(in,
  *blocked), nil }`. `escapeOrRefuse` is reachable only from the `case 0` and
  `default` arms **after** that early return. `gate` returns `blocked` for both
  `KindGuardUnevaluable` and `KindOwnedStateUnavailable`. So an escape row cannot
  rescue either class. **R-03 grounds.**
- **R-04 — peer clause read verbatim.** RDR 0007 `:1412-1414` fixes `¬U = U`;
  `:1435-1436` fixes the verdict as `all_result ∧ ¬(unless_conj)`. RDR 0003
  modelled `unless` as a decided set subtraction and its disposition table had
  **no row** for an unevaluable `unless` atom. **R-04 grounds.**
- **R-06 — shipped struct read.** `resolve.go::Refusal` carries `Guard string`
  (one opaque predicate) and `Rows []RowRef` (all implicated rows). RDR 0006's
  finding contract has no atom field (`0006:363-367`). Atom-naming has no carrier
  on **either** surface. **R-06 grounds, and B's runtime half is net-new.**
- **R-08 — fixture read.** `status` is `[rule.match.status]` in
  `profile-to-grounding` and `[rule.guard.unless.status]` in
  `reconcile-rewind-legality`. **R-08 grounds** (severity PARTIAL per the diff:
  the two rules are in different groups, so no group currently collides).
- **R-17 — absence verified.** `conform` occurs nowhere in the kernel;
  `resolve.go::assemble` merges three tag sources into a `TagSet` and reads no
  declaration. **R-17 grounds.**
- **R-10 — peer text read.** RDR 0006's default-on rider is real
  (`0006:913-915`). **R-10 grounds as a real consequence**, though it is RDR
  0006's rider and not this RDR's to reverse.

Nothing was dismissed for failing the code gate. Three rows were dismissed as
already-adjudicated (below).

## Dispositions

| Ledger | Disposition | Sections touched |
| --- | --- | --- |
| R-01 (trace step 6 uncomputed) | **fixed** — step 6 rewritten with the executed arithmetic and a GAP verdict; step 3 now counts the presence dimension; closing note corrected from "no GAP row"; Phase 3 prerequisite recorded so Scenario 2 does not reuse the fixture unchanged | `trace` steps 3/6 + closing note; Testing Strategy; phase-gating table |
| R-02 (`set` cardinality off by an exponent) | **fixed** — bound clause now carries a per-kind **assignment count** table (`enum`/`int` `\|domain\|`; unmarked `2^\|domain\|`; `set` `2^\|universe\|`; optional ×2) and the worked 8,000-vs-`2^60` case | Normative Contracts (bound clause) |
| R-03 (escape row cannot rescue `guard_unevaluable`) | **fixed** — coverage clause narrowed to the classes an escape row can actually rescue, citing the gate's early return | Normative Contracts (escape clause); Contradiction Check; `authority` (new row) |
| R-04 (`unless` two-valued vs. Kleene) | **fixed** — new normative clause; "can refuse" now reads **both blocks**; disposition table gains the unevaluable-`unless` row | Normative Contracts (2 clauses); `disposition`; `authority` (new row); Contradiction Check |
| R-05 (composed worst-case defaults) | **fixed** — the same per-kind table makes both multipliers explicit and computable in one place | Normative Contracts (bound clause) |
| R-06 (atom-naming has no carrier) | **fixed as a booked record** — consequent-duty note extended with the runtime half; new **A20** | Normative Contracts (note); Critical Assumptions |
| R-07 (bare escape row = silent opt-out) | **fixed** — bare-escape coverage closure must be observable in the verdict; new **A19** books RDR 0006's producer | Normative Contracts (escape clause); Critical Assumptions |
| R-08 (match/guard split per key vs. per atom) | **fixed** — participation clause states the split is per atom **per group**, not a per-key property of the model | Normative Contracts (participation clause) |
| R-09 (non-discriminating key demands a domain) | **dismissed-with-cite** — this is the participation clause working as decided: a key that bounds the product must be finite or the group is unprovable. Re-litigating it is the settled "rows-must-differ" reading the clause explicitly rejects, and `trace` step 3 records why | — |
| R-10 (docs-only PR breaks CI) | **charted-to-successor** — real, but it is RDR 0006's default-on rider (`0006:913-915`), not this RDR's clause; reversing it here would re-decide a peer's rider unilaterally | `Charted.md` |
| R-11 (§JD-14 reversal) | **dismissed-with-cite** — the Contradiction Check already adjudicates this verbatim: "it does not re-decide §JD-14 unilaterally", the coverage half stands, the overlap half is booked as A17 with the cross-document repair named. B's *process* charge is answered by the draft; A's substantive attack is R-03, which is fixed | — |
| R-12 (A14 contradiction vs. request) | **fixed** — A14 already carries the producer request and MVV Scenario 7; no clause change owed, but Scenario 8 now also exercises the `unless` block so a block-collapse defect surfaces on two scenarios rather than one | Testing Strategy |
| R-13 (A15's test cannot falsify) | **fixed** — Scenario 3 now measures whether the proof representation actually completes per shape and compares that to the bound's verdict; A15's Evidence records that the consistency comparison alone cannot close it | Testing Strategy; A15 |
| R-14 (A10/A12 deferred to a Draft peer) | **fixed as aggregate disclosure** — no new clause (the records already carry plans); folded into the new aggregate assessment, which names Phase 2 group construction as genuinely gated | Assumption Verification |
| R-15 (aggregate never assessed) | **fixed** — new aggregate open-assumption table grouping all nine Pending records by who can close them, with an explicit honest verdict | Assumption Verification |
| R-16 (illustrative code disowned) | **dismissed-with-cite** — the block already closes with "the field's authoring location is RDR 0002's to fix (**A16**), which is why this block shows a shape rather than asserting one." The disclosure the finding asks for is the text that is already there | — |
| R-17 (conformance has no enforcer) | **fixed** — conformance clause now states the gap explicitly; new **A18** books the producer | Normative Contracts (conformance clause); Critical Assumptions; `authority` (new row) |
| R-18 (type system absorbed without alternatives) | **charted-to-successor** — a real proportionality gap, but authoring an alternatives pass for the declaration model is net-new scope for a document already at lock; the aggregate assessment names it as the thing to watch | `Charted.md` |
| R-19 (A1 verified on a pre-model harness) | **fixed** — A1's Scope-of-record now states the harness predates the declaration model and does **not** prove `contains` over a *declared* element universe | A1 |
| R-20 (governance crowds out engineering) | **dismissed-with-cite** — ranked NOISE by the diff; the Metadata qualifier is the re-entry contract Stage 5 requires, and this pass moved arithmetic *into* the body rather than the ledger | — |

**14 fixed, 4 dismissed-with-cite, 2 charted.**

## Compute, don't argue

- **Coverage** — executed, not asserted: 64 product, 32 union, 32 uncovered,
  overlap 0. The witness assignment written into step 6
  (`profile=small, prelock_iterations=0, cluster_eligible=true, present`) was
  confirmed present in the uncovered set.
- **Cardinality** — the per-kind table was derived from the RDR's own
  single-valued clause rather than invented: `2^|domain|` for co-occurring
  values was already written there, so `set` (values co-occur by construction)
  must be `2^|universe|`. The correction reconciles two clauses that disagreed.
- **Assumption census** — recomputed by parsing `- **A<N>` blocks for their
  `Status` line: **20 records, 11 Verified, 9 Pending**. The roster and the
  aggregate table read that computed set, not a carried-forward count.
- **Cross-lens check** — the `authority`, `disposition`, and `trace` tables were
  re-read against this pass's clause edits, since the row edits a shared draft.
  All three needed changes and got them: `trace` steps 3/6 and its closing note,
  `disposition`'s new unevaluable-`unless` row, and four new `authority` rows.
  The 3amigo iteration-3 pass had corrected `trace` step 3's *match-key* reading
  and this pass corrected the same step's *dimension count* — two lenses, one
  table, no conflict between the edits.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `declared domain size` now occurs once, inside the
   clause that corrects it. `closes coverage` (5 sites) re-read: all consistent
   with the narrowed escape rule. `can refuse` (16 sites) re-read: the
   declaration-only test still holds; the two clauses that define it now both say
   "either block".
2. **Subject-token grep** — `unless` re-read against the Kleene correction at
   the derivation (`A2`), the disposition table, the "can refuse" clause, and
   MVV Scenarios 1/7/8. The A2 derivation's "subtracted" language is correct as
   the *lint projection over decided atoms* and is now explicitly scoped that way
   by the new clause, so it was left standing rather than reworded.
3. **Added clause amends its contract block** — the four new/amended normative
   clauses were re-read against their siblings: the Kleene clause against the
   narrowing it strengthens, the escape-classes clause against the two-population
   overlap clause (compatible — one is about coverage, the other about overlap),
   the assignment-count table against the single-valued clause it derives from,
   and the observability clause against the anti-opt-out paragraphs it enforces.
4. **New order/reachability/emission claim** — three new claims about *who
   produces what* (conformance enforcement, bare-escape observability,
   atom-level carriage) were each registered as assumptions **A18/A19/A20**
   rather than asserted as agreed.

## Mini-checks

The cue read was performed at iteration 1 / the grounding pass; the `authority`,
`disposition`, and `trace` tables persist in the draft and all three were
updated by this pass. No fix added a **new** cue — the per-kind assignment-count
table is a normative clause, not a mini-check table — so no new table is owed.

## Needs (re)verification (carried to Stage 6)

New records, all `Pending`:

- **A18** — a conforming-view check has a producer. Unhomed today: the clause
  names view assembly, which performs no check. Request at RDR 0007's next touch
  and/or RDR 0006's refine.
- **A19** — a bare escape row's coverage closure is observable in lint output.
  RDR 0006 refine, travelling with A17's population question.
- **A20** — atom-naming gets a carrier on both surfaces. RDR 0006 refine (lint,
  already §JD-4's duty) + RDR 0007 (runtime payload).

No previously-`Verified` record was demoted by this pass. A1's stamp stands with
its scope tightened (R-19): it proves evaluation, not exhaustiveness over a
declared universe.

## Charted to successor

Two — see `Charted.md` (R-10, R-18).

## Tiebreakers escalated

None. Two forks collapsed on evidence:

1. **R-11 — is the §JD-14 reopening a defect or a disclosure?** Collapsed on the
   draft's own text: the Contradiction Check states the reopening is deliberate,
   keeps the coverage half, refutes only the overlap half's premise against
   shipped code, and books A17. B's process charge fails; the residual
   cross-document staleness is A17's named repair, already scheduled.
2. **R-03 vs. R-11 — does A17 go far enough?** It did not, and this was the
   pass's most valuable finding: A17 reads `Resolve` step 2 correctly and stops
   one function short of `gate`, so the RDR cited the kernel as authority for an
   escape-row behaviour the kernel does not have. Fixed at the clause rather than
   re-litigated at the assumption.
