Model: claude-opus-5[1m]

# Grounding Dispositions — iteration 3 (re-entry)

Origin ledger = this iteration's `findings.md`. Three entries.

| Finding | Disposition | Origin | Sections touched |
| --- | --- | --- | --- |
| G3-1 — nine sites still describe A8 as open/contested/Pending after the resolve re-entry flipped it to `Verified` | **fixed** | G3-1 (origin ledger entry 1) | Technical Design (default-on paragraph, provenance paragraph); Normative Contracts (peer-obligation block); `authority` census; `trace` desk trace step 10 + summary; Decision Rationale; Contradiction Check; Assumption Verification roster; Scope Verification |
| G3-2 — §JD-4's line anchors into this RDR (`0003:781-786`, `0003:818-826`) name the wrong clauses | **charted-to-successor** | G3-2 (origin ledger entry 2) | none — the defect is in `docs/jdr/0001-resolve-kernel-seam.md`, a peer document this skill does not edit |
| G3-3 — seven sites route A10/A12 to a cluster-reconcile venue that already ran and homed them at RDR 0006's refine | **fixed** | G3-3 (origin ledger entry 3) | Critical Assumptions A10 (Plan, Stage-6 disposition, travels-with); A12 (same three); Prerequisites checklist (A10 and A12 entries) |

## Grounding gate

**G3-1.**
1. **Code on `main`** — n/a (document-state finding). Every cited passage read
   verbatim: `docs/jdr/0001-resolve-kernel-seam.md:238-264` (§JD-4 CLOSED, the
   recording assignment, the no-warning-tier half, the consequent duty),
   `docs/rdr/0006-graph-lint-authority-and-guarantees.md:8-16` (Status: `Draft`,
   demoted 2026-08-21, carrying the §JD-4 citation duty), `:317` and `:643`
   (`graph-unprovable-coverage` in the blocking table and MVV matrix), `:363-367`
   (finding contract, no atom-level field).
2. **{RDR_RESOURCES}** — no principle bears on it. The fixes change no design
   claim, no error shape, and nothing `docs/cli-output-contract.md` governs;
   they align the draft's account of its own assumption state with the record.
3. **RDR's own decided text** — not a re-raise. The finding contradicts the
   draft's *claims about its own state*, which its own A8 record (`Verified`)
   and its own Prerequisites checklist (`[x] A8 closed`) already refute. Two
   halves of one document disagreeing is a defect, not a settled call.

**G3-2.** Same passages. Charted rather than fixed because the wrong anchors
live in the JDR; this skill authors RDR documents and does not edit peers. It is
not lock-blocking here — the clause §JD-4 assigns to this RDR is present and
correct, only the JDR's pointer to it is off.

**G3-3.**
1. **Code on `main`** — n/a. Read
   `docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md:85-86`
   (routing table: "JOINT-DECISION → tolerance") and `:162-176` (Standing
   tolerances section, the Home column reading `0006's refine` for both, and the
   A12 decision-procedure note).
2. **{RDR_RESOURCES}** — no principle bears on it.
3. **RDR's own decided text** — not a re-raise. A10 and A12 keep their substance
   and their `Pending` status untouched; only the venue pointer, which the gate
   superseded, is corrected.

## Compute, don't argue

The A8 roster fix is an executed subtraction, not an argument. Before the fix
the Assumption Verification roster read `Pending: A8, A10, A12, A14, A15` while
the assumption records themselves read:

```sh
grep -n -A1 '^- \*\*A[0-9]' docs/rdr/0003-guard-predicate-exhaustiveness.md \
  | grep -E 'Status'
```

returning `Verified` for A1–A9, A11, A13 and `Pending` for A10, A12, A14, A15 —
eleven Verified, four Pending. The roster claimed ten Verified and five Pending.
The difference is exactly A8. The roster now reads the computed set.

Likewise the desk-trace step 10 verdict is the executed comparison, not a claim
about it: its predicate is "exactly one document records the narrowing; the
other cites it", and §JD-4:242 assigns the recording to this RDR while
`0006:10-13` records the citation from the other side — both halves present, so
the row returns OK.

Cross-lens check: the `authority` census row and the `trace` step-10 row were
both written by earlier lenses in this row and both pinned the same fact
(A8/§JD-4) as contested. They are now pinned one way, consistent with the A8
record and with the `:932` census row the resolve re-entry had already updated.

## Tiebreaker reduction

No fork survived. The apparent one — "is A8 really closed, or did the resolve
re-entry close it prematurely?" — collapsed on the evidence: §JD-4 states
`CLOSED 2026-08-22` with both halves decided in text, names the venue, and RDR
0006 records the same assignment independently on its own Status line. Two
documents agreeing, plus the gate report that produced them, is not an
indeterminate call.

The second apparent fork — "should the peer-obligation block be deleted or
rewritten?" — also collapsed. The block states a true fact (RDR 0006's finding
contract has no atom-level field, confirmed at `0006:363-367`) under a false
disposition (not-yet-agreed, against a `Final` document). Deleting it would lose
a fact an implementer needs; rewriting it under the decided disposition keeps
the fact and fixes the frame.

## Edits applied

- **Technical Design** — the default-on paragraph's divergence now routes to RDR
  0006's refine rather than a spent cluster-reconcile venue; the provenance
  paragraph states §JD-4's assignment instead of "A8 tracks the agreement".
- **Normative Contracts** — the peer-obligation block reframed from "not yet
  agreed" to §JD-4's decided consequent duty on RDR 0006, with the correct
  `Draft` status and a line stating it does not gate this RDR's lock.
- **`authority` census** — the exhaustiveness-verdict row's Canonical cell
  changed from "Contested — A8" to this RDR, citing closed §JD-4 and naming why
  both triggers survive. Now agrees with the withheld-claim-artifact row below it.
- **`trace` desk trace** — step 10 flipped GAP → OK with the executed
  comparison as its witness; the summary sentence updated from "one GAP row
  remains" to no GAP row.
- **Decision Rationale** — §JD-4 described as closed and assigning; the
  "one class is not yet provable" sentence replaced by both classes closing.
- **Contradiction Check** — the exhaustiveness-strength paragraph rewritten to
  decided-and-agreed, matching the Prerequisites checklist it contradicted.
- **Assumption Verification** — A8 moved from the Pending roster to the Verified
  roster.
- **Scope Verification** — the "A8 is Pending and is the one record that still
  gates lock" paragraph rewritten to A8 `Verified` with no peer-decision lock
  gate; the A10/A12 paragraph no longer routes them through A8's venue. The
  stale "grounding iteration 2 is complete" line updated to iteration 3.
- **Critical Assumptions A10 / A12** — Plan lines now name RDR 0006's refine as
  the home the cluster gate assigned, A12's carrying the gate's decision-procedure
  note (graph reachability vs. syntactic decision over declarations);
  Stage-6 disposition headers and travels-with lines updated. Status and
  substance unchanged — both stay `Pending`.
- **Prerequisites** — the A10 and A12 checklist entries' Venue lines updated to
  RDR 0006's refine.

## Amendment sweep (rdr-common §amendment-sweep)

1. **Pre-edit token grep** — `rg 'A8'` returned 20 sites before the edit. All 20
   reviewed: 9 stale (fixed), 11 already correct (Status line, the A8 record
   itself, the `:932` census row, the checklist entry, both phase-table rows, and
   the cross-references inside A10/A12/A14). `rg 'cluster reconcile'` returned 8
   further sites; 7 stale (fixed), 1 correct (the A8 record's own account of the
   venue that closed it).
2. **Subject-token grep** — `rg 'guard_unevaluable'` and `rg 'narrow'` across the
   draft; producers and consumers re-read. No semantic disagreement: the
   narrowing normative clause (`:809-813`), the withheld-claim clause
   (`:816-822`), the `disposition` table's absent-key row (`:954`), the Failure
   Modes entry (`:1339`), and MVV scenario 6 all remain consistent, and none
   asserted an A8 state. `rg '`Final`'` confirmed no surviving claim that RDR
   0006 is Final — the two remaining hits are RDR 0007, which is Final.
3. **Added clause amends its contract block** — no normative clause was added or
   reworded. The peer-obligation block is a non-normative note beside the
   withheld-claim clause; its siblings were re-read and neither the narrowing
   clause nor the report-every-defect clause changes meaning under the corrected
   note.
4. **New order/reachability/emission claim** — none added. G3-3 corrected the
   venue of an existing reachability assumption (A12) without changing what it
   claims or its `Pending` status.

## Mini-checks

Cue read performed against the current draft. The three tables the draft already
carries — `authority`, `disposition`, `trace` — stay fired and were updated in
place by this pass. No new cue fired: the fixes add no fallback path, derived
output, oracle, or encode/decode pair. `round-trip / fidelity` remains unfired
by design (`Round-Trip / Inverse Invariants`: "This RDR introduces no
encode/decode pair"), and `test-discriminability` stays discharged by the MVV's
existing named negative controls rather than an absence-of-error oracle.

**Mini-checks fired: `authority`, `disposition`, `trace` (all pre-existing;
updated, not newly fired).**

## Needs (re)verification (carried to Stage 6)

None new. No fix touched or added a load-bearing claim: every edit corrected the
draft's account of an assumption state that was already decided elsewhere, and
no normative signature, format, exactness word, or external-behavior claim
changed. No previously `Verified` assumption was invalidated — A8 was
*restored* to the Verified state its own record already carried.

Carried unchanged from prior passes: A10, A12 (tolerances, homed at RDR 0006's
refine), A14 (carriage request on RDR 0002), A15 (discharged by MVV Scenario 3).

## Charted to successor

- **G3-2** — `Charted:` §JD-4's two line anchors into RDR 0003 name the wrong
  clauses: `0003:781-786` points at the presence-dimension projection (A7), not
  the narrowing (now at `:809-813`), and `0003:818-826` points at the
  peer-obligation block this pass rewrote, not a warning-tier statement.
  Out of scope because the defect is in `docs/jdr/0001-resolve-kernel-seam.md`,
  a peer document. Suggested successor: fold into the next JDR 0001 touch or the
  RDR 0006 refine that discharges §JD-4's citation duty — same document set,
  same edit session. Not lock-blocking for this RDR.

## Tiebreakers escalated

None.
