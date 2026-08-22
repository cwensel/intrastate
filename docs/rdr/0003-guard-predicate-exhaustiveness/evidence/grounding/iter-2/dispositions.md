Model: claude-opus-5[1m]

# Grounding Dispositions — iteration 2 (re-entry)

Origin ledger = this iteration's `findings.md`. One entry.

| Finding | Disposition | Origin | Sections touched |
| --- | --- | --- | --- |
| G2-1 — §JD-4 names RDR 0006 as the narrowing document, not this RDR; and §JD-4 is an open entry, not a closed assignment | **fixed** | G2-1 (origin ledger entry 1) | Metadata Status; Critical Assumptions (new A8); Technical Design; Normative Contracts; Decision Rationale; Implementation Plan / Prerequisites; Finalization Gate / Contradiction Check; Finalization Gate / Assumption Verification; References |

## Grounding gate — G2-1

1. **Code on `main`** — n/a (document-citation finding). Every cited passage
   read verbatim: `docs/jdr/0001-resolve-kernel-seam.md:238-241` (§JD-4 head
   clause + open-status sentence), `docs/rdr/0007-guard-predicate-totality.md:159-161`
   (Predecessors naming RDR 0006 as the §JD-4 tolerance-holder) and `:593-597`
   (A12 restating §JD-4 as not-Closed),
   `docs/rdr/0006-graph-lint-authority-and-guarantees.md:9` (Status qualifier)
   and `:344-350` (exhaustiveness clause, non-finite case only).
2. **{RDR_RESOURCES}** — no principle forbids the fix. The narrowing itself is
   unchanged in substance and stays consistent with
   `docs/cli-output-contract.md` (the refusal still reaches the CLI as a
   structured envelope; nothing about error shape moved).
3. **RDR's own decided text** — not a re-raise. The RDR's Contradiction Check
   asserted this boundary was "resolved rather than open"; §JD-4 records the
   opposite. The finding contradicts a *claim about the source*, not a design
   option the RDR already weighed and declined.

## Tiebreaker reduction

The apparent fork — "should 0003 or 0006 record the narrowing?" — collapsed on
the evidence rather than escalating:

- §JD-4 decides the **substance** ("the *promise* narrows — P5 decides") and
  leaves open only the recording document. So stating the rule is not
  contested; only its attribution was.
- The narrowing constrains what **this RDR's** finite-domain product may claim,
  and RDR 0006's A2 consumes 0003's coverage/overlap derivation rather than
  defining its own. Recording it here is defensible on the merits.
- RDR 0007 A12 already routes the adjacent projection question here with the
  named plan "it lands when 0003 states its existence-atom projection",
  establishing the same routing direction.

So the fix keeps the clause where it is and repairs the citation, rather than
deleting the clause or escalating the fork. What could **not** be collapsed
unilaterally is the sibling agreement: RDR 0006 is `Final` and carries no
`guard_unevaluable` narrowing, and this skill does not edit peer RDRs. That is
booked as A8 (Pending) with a route-back plan, not resolved here.

## Edits applied

- **Technical Design** — the narrowing paragraph keeps its rule but drops the
  bare `(JDR 0001 §JD-4)` attribution; a following paragraph states what §JD-4
  actually decides, what it leaves open, that its head clause names RDR 0006,
  and why this RDR records it anyway.
- **Normative Contracts** — the exhaustiveness-narrowing clause gains a second
  sentence obliging RDR 0006 to carry the same narrowing and forbidding the two
  documents from stating it differently.
- **Critical Assumptions** — new **A8** (`Pending`, Method: Peer RDR) carrying
  the cross-document agreement, its evidence-needed line, a route-back plan, and
  an "If wrong" naming the single-source failure.
- **Decision Rationale** — "narrows its own lint promise under §JD-4" replaced
  with the accurate account plus an A8 pointer.
- **Contradiction Check** — split: §D4 (guard-domain enforcement) stays
  "resolved rather than open" and is now marked as a Closed entry; the
  exhaustiveness-strength boundary is restated as decided-in-substance but not
  yet agreed across documents, pointing at A8.
- **Assumption Verification** — the stale "no record remains `Pending`" sentence
  corrected to name A7 and A8. (This sentence already contradicted A7 before
  this pass; fixed in the same amendment sweep.)
- **Prerequisites** — the unchecked all-assumptions item now names A8 alongside
  A7.
- **Metadata Status** and **References** — the §JD-4 shorthand qualified so the
  Status line and the reference entry no longer imply a closed assignment.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `rg 'JD-4'` returned 6 sites before the edit
   (Status, A7 "If wrong", Technical Design, Decision Rationale, Contradiction
   Check, References). All 6 reviewed; 5 updated. A7's "the §JD-4 narrowing this
   RDR just adopted" left as-is — still true, since the RDR does adopt the rule;
   only the attribution was wrong, and A7 makes no attribution claim.
2. **Subject-token grep** — `rg 'narrow'` and `rg 'guard_unevaluable'` across
   the draft; producers and consumers re-read. No semantic disagreement found:
   the Failure Modes entry ("A guard that cannot be decided at runtime surfaces
   as RDR 0007's `guard_unevaluable` refusal") and MVV Scenario 6 both remain
   consistent with the narrowed clause.
3. **Added clause amends its contract block** — the exhaustiveness-narrowing
   clause's siblings in `Normative Contracts` re-read for agreement. The
   finite-domain clause and the refuse-or-downgrade clause both remain
   compatible: they gate on *provability*, the new sentence gates on
   *cross-document agreement*, and neither weakens the runtime veto.
4. **New order/reachability/emission claim** — none added. A8 is a document-
   agreement claim, registered as a Pending assumption per this rule rather
   than asserted.

## Needs (re)verification (carried to Stage 6)

- **A8** (new, `Pending`, Method: Peer RDR) — RDR 0006 adopts the same
  `guard_unevaluable` narrowing, or §JD-4 is dispositioned to assign the
  recording to one document. **Requires a route-back:** RDR 0006 is `Final`,
  so this cannot close inside this RDR. Cluster reconcile is the natural venue —
  0003, 0006, and 0007 all touch it.
- **A7** (carried, unchanged by this pass) — still `Pending`, still blocked on
  RDR 0002's per-tag optionality declaration.

No previously `Verified` assumption was invalidated by this pass. A5, the
assumption this re-entry was scoped to re-verify, is confirmed on all its cited
grounds (see `findings.md` CONFIRMED) and stays `Verified`.

## Charted to successor

None. The finding was in scope: it corrects a citation this re-entry introduced.

## Tiebreakers escalated

None.
