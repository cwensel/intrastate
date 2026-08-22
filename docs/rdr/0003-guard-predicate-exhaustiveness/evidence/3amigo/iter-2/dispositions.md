Model: claude-opus-5[1m]

# 3amigo Dispositions — iteration 2 (re-entry)

Origin ledger = this iteration's `consolidation.md` (10 entries over 22 persona
findings). Every finding traces to a ledger entry; none was net-new scope.

## Grounding gate

All 22 findings passed. Verified verbatim before any edit:

- `internal/resolve/resolve.go` — still ships `Row.Guard string` (:185) and
  `GuardEvaluator.Evaluate(guard string, view TagSet)` (:91). **P2-7 grounds.**
- `docs/rdr/0007-guard-predicate-totality.md:1249-1266` — the SEAM clause fixes
  four atom fields (`Key`, `Operator`, `Literal`, `Block`). **No position
  field: P2-1 grounds.** `:1568` confirms A5's conjoined value row (two atoms,
  one key), and row 17 (`:2245`) sorts payloads on the TOTAL tuple
  `(RuleID, SourceLocator, key, block, operator, literal)` — never atom order.
- `docs/rdr/0006-graph-lint-authority-and-guarantees.md` — zero occurrences of
  `guard_unevaluable` or "narrow" (**P2-3, P3-5 ground**); `:344-350` scopes its
  exhaustiveness clause to the non-finite case only; `:302` defines
  `graph-unprovable-coverage` as the blocking inability-to-prove code; `:696`
  gates on a contract "that claims closed coverage" without defining the claim
  (**P2-6 grounds**).
- `docs/rdr/0002-*.md:214-222` — tag declarations carry name, provenance, value
  kind, optional accessor reference; **no optionality field and no element
  universe** (**P1-4, P2-2, P2-8, P3-3 ground**). Status is `Draft`.
- `docs/jdr/0001-resolve-kernel-seam.md:238-241` — §JD-4 open as to recording
  document and "whether lint gains a warning category" (**P2-4, P3-1 ground**).

No finding was dismissed for failing the gate. No finding re-litigated an
already-decided option.

## Dispositions

| Finding(s) | Ledger entry | Disposition | Section touched |
| --- | --- | --- | --- |
| P1-1, P1-2, P3-2, P3-6 | 1 | **fixed** | Normative Contracts (2 new clauses: observable withheld claim; "can refuse" decision procedure); `disposition` table rows 4 & 11; MVV; Testing Strategy Scenario 6 |
| P2-4, P3-1 | 2 | **fixed** | Normative Contracts (new clause: refuse/downgrade are one blocking outcome; this RDR MUST NOT mint a warning category §JD-4 leaves open); Technical Design (reading pinned at first use) |
| P2-3, P2-5, P3-5 | 3 | **fixed** | Normative Contracts (peer-binding MUST removed from the narrowing clause); A8 Plan (closure is a §JD-4 assignment of ONE recording document, not matching prose); `authority` + `trace` step 10; Assumption Verification |
| P1-3, P2-2, P3-4 | 4 | **fixed** (scope bound stated; A7 stays `Pending`) | Trade-offs / Consequences (new `exists` carve-out naming the fixture group and RDR 0007's two-row pattern); MVV operator list gains `exists` |
| P1-4, P3-3 | 5 | **fixed** | MVV (positive assertion + negative control); Prerequisites (A7 given an explicit route, done-condition, and verifier) |
| P2-8 | 5 | **fixed** — new assumption **A9** | Critical Assumptions (A9, `Pending`, Method: Peer RDR); Capability Dependencies (own row); Assumption Verification |
| P2-1, P3-7 | 6 | **fixed** | Load-Bearing Decisions › Identity — position-based handle replaced with the total tuple; cites why (no position field; reorder invariance; conjoined row) |
| P2-6 | 7 | **fixed** | Technical Design — states default-on-for-finitely-declared-groups reading, with the opt-in alternative routed to cluster reconcile |
| P2-7, P3-8 | 8 | **fixed** | Prerequisites — split A7/A8 by closure route with done-conditions; RDR 0007 sequencing item un-checked and scoped to implementation, not lock |
| P1-5 | 9 | **fixed** | Decision Rationale — promise qualified by the two classes A7/A8 currently exclude; names RDR 0006/0005 as inheritors |
| P1-6 | 10 | **fixed** | A5 "If wrong" — trigger restated as RDR 0007's unshipped reshape dropping a named field, not a normalization step that no longer exists |

All 22 findings **fixed**. None dismissed, none charted to a successor.

## Tiebreaker reduction

Three apparent forks collapsed on evidence rather than escalating:

1. **What lint outcome does a withheld claim take?** Apparent fork: mint a new
   code vs. stay silent. Collapsed — §JD-4 leaves "whether lint gains a warning
   category" open, so this RDR may *not* mint one; RDR 0006 already defines
   exactly one blocking inability-to-prove code (`:302`). Reusing it is the only
   arm that neither invents a category nor leaves the claim unobservable. The
   open question (reuse vs. own code) is handed to the A8 route-back, not decided
   here.
2. **Which document's wording wins for the narrowing?** Apparent fork: 0003
   restates vs. 0006 restates. Collapsed — both arms in A8 reduce to "exactly one
   document records it, the other cites it", because a shared restatement is the
   duplication A8's own "If wrong" forbids. So the route-back asks §JD-4 to
   *assign*, and the clause here stops issuing a MUST at a `Final` sibling.
3. **How is per-atom identity handled without a position field?** Apparent fork:
   slice index vs. request a fifth field from `Final` RDR 0007. Collapsed —
   RDR 0007's own determinism row already sorts on a total tuple that is
   position-free, so this RDR adopts that tuple. No change request to a `Final`
   peer, and reorder invariance (Scenario 5) is preserved.

## Needs (re)verification (carried to Stage 6)

- **A9** (new, `Pending`, Method: Peer RDR) — a declared set-element universe for
  `contains`. Travels with A7 as one RDR 0002 tag-declaration producer request.
- **A7** (carried, `Pending`) — now gates **two** things, not one: the presence
  dimension *and* the new "can refuse" decision procedure, which reads presence
  from declarations. Recorded in A7's Plan.
- **A8** (carried, `Pending`) — closure narrowed to a §JD-4 assignment of one
  recording document; the route-back must also settle whether the withheld-claim
  case reuses `graph-unprovable-coverage` or earns its own code.
- **Row-group claim semantics** (P2-6) — this RDR now states default-on; if RDR
  0006 intends an explicit annotation, that is a divergence for cluster reconcile.
  Not registered as its own assumption: it is a reading of a peer clause, folded
  into the A8 venue.
- No previously `Verified` assumption was invalidated. A5 stays `Verified` — the
  fixes to its "If wrong" and to Identity sharpen the record without disturbing
  its evidence.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `refuse or downgrade` returned 6 sites; all reviewed
   against the new one-outcome clause, all remain literally accurate, and the
   reading is pinned at first use (Technical Design) so none reads as a soft arm.
   `MUST carry this same`/`state it differently` → 0 sites after the edit
   (peer-binding MUST fully removed).
2. **Subject-token grep** — `withhold|withheld` (9 sites) re-read: the two
   `disposition` rows that named only RDR 0007's *runtime* payload for a *lint*
   outcome were the actual defect and are corrected; `trace` step 8 and the MVV
   now agree.
3. **Added clause amends its contract block** — the four new/edited normative
   clauses re-read against their 10 siblings. The finite-domain clause, the
   too-large clause, and the owned-tag clause remain compatible; the new "can
   refuse" clause reuses the *existing* `reachable predecessor` notion (line 512,
   RDR 0006 via A6) rather than minting a second reachability authority.
4. **New order/reachability/emission claim** — the "can refuse" clause makes a
   declaration-decidable presence claim. Grounded at its authority (the owned-tag
   reachability clause) **and** registered as a dependency on A7's blocked
   producer request rather than asserted as already true.

## Mini-checks

Cue read performed at iteration 1 / the grounding pass; `authority`,
`disposition`, and `trace` tables persist in the draft and were all three
updated by this pass. No fix added a new cue.

## Charted to successor

None.

## Tiebreakers escalated

None — all three collapsed on evidence (above).
