Model: claude-opus-5[1m]

# Grounding Findings — iteration 3 (re-entry)

Scope: delta against the iteration-2 grounding commit (`624d211`), i.e. every
codebase/peer claim added or edited by the declaration-model rehoming
(`17fc26d`, `158164b`), the cluster gate (`21322fb`), and the Stage-4 resolve
re-entry (`894d183`, `d0f99b1`). There is no `Ground-sweep:` verdict line in
Decision Rationale, but iteration 1 and 2 already swept the pre-rehoming text;
their CONFIRMED claims are not re-swept. Iteration 2's sole finding (G2-1) is
closed — see CONFIRMED. Three findings this pass, all stale-citation defects
introduced by amendment sweeps that did not reach every site.

## REFUTED

### G3-1 — Six sites still describe A8 as open, after the resolve re-entry flipped it to `Verified`

The re-entry commit (`894d183`) closed **A8** on JDR 0001 §JD-4's 2026-08-22
disposition and updated the Metadata Status line (`:14`), the ownership census
row for the withheld-claim artifact (`:932`), the Prerequisites checklist
(`:1381-1390`), and the phase-gating table (`:1461`, `:1463`). The amendment
sweep did **not** reach six further sites, each of which still asserts A8 is
Pending, contested, or awaiting a §JD-4 disposition. Verified verbatim against
the closed entry at `docs/jdr/0001-resolve-kernel-seam.md:238-264`.

**What §JD-4 now says** (`:240-242`, verbatim):

> **CLOSED 2026-08-22** by the `0003-0006-0007` cluster gate, which is the venue
> 0003 A8, 0006 Direction 1, and this entry's own recommended arm all routed to.
> **RDR 0003 is the recording document; RDR 0006 cites it and mints no code.**

And RDR 0006's Status line (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:10-13`)
records the same assignment from the other side, including "reusing
`graph-unprovable-coverage` with no new code and no non-blocking tier".

The six stale sites:

| # | Line | Stale text | Why refuted |
| --- | --- | --- | --- |
| 1 | `:627` | "that is a divergence to settle at cluster reconcile alongside A8" | The cluster reconcile it defers to has already run (2026-08-22) and closed A8. The *default-on vs. per-group-annotation* divergence it names is real and still open, but it can no longer travel with A8. |
| 2 | `:666` | "RDR 0006 is now `Draft` and carries the §JD-4 recording question in its Refinement Context. A8 tracks the agreement the two documents owe each other." | §JD-4's recording question is closed, not carried. What RDR 0006 now owes at refine is a *citation* plus the atom-level field — a consequent duty, not an open agreement. |
| 3 | `:818-826` | "**Peer obligation, not yet agreed (A8).** … a payload extension this RDR requests of a `Final` document … The route-back in A8 must carry it explicitly" | Two errors. (a) The obligation *is* agreed — §JD-4:260-263 records it as a decided consequent duty on RDR 0006, discharged at its refine. (b) RDR 0006 is no longer `Final`; it was demoted to `Draft` on 2026-08-21 (`0006:9`), so this is a live request on an open peer, not a route-back against a locked document. A8's own record at `:293-300` already states both facts correctly — this block contradicts it. |
| 4 | `:926` | Ownership census, exhaustiveness-verdict row: "**Contested — A8.** … A8 closes on §JD-4 assigning ONE recording document, not on both restating it." | §JD-4 made exactly that assignment. The row two lines below it (`:932`, withheld-claim artifact) was updated to cite the closed §JD-4; this one was not. The two rows now disagree with each other about the same decision. |
| 5 | `:972-978` | Desk-trace step 10: "**GAP — booked as A8.** … a divergence with an open peer … that needs a §JD-4 assignment, not duplicated prose", and the summary "One GAP row remains — A8, the cross-document assignment — carrying a named plan and not asserted as already true." | The §JD-4 assignment it says is needed has been made. Step 10's own predicate ("exactly one document records the narrowing; the other cites it") is now **satisfied**, so this is a passing row misreported as a gap — and the summary sentence asserts a GAP count of one that no longer holds. |
| 6 | `:1071-1080` | Decision Rationale: "§JD-4 decides the substance but leaves the recording document open, and names RDR 0006's proof in its head clause; A8 carries the resulting sibling obligation." … "One class is not yet provable: a group whose narrowing verdict depends on the cross-document agreement A8 tracks." | Both sentences are refuted. §JD-4 no longer leaves the recording document open, and the "not yet provable" class was created solely by that openness — with the assignment made, the class closes. This is a live understatement of the RDR's own promise. |
| 7 | `:1630-1636` | Contradiction Check: "The strength of the exhaustiveness claim is decided in substance but not yet agreed across documents. … §JD-4 remains open as to which document records it and names RDR 0006's proof … A8 tracks it, and it must close before lock." | Directly contradicts `:1381-1390`, which records A8 as closed on the same evidence. One of the two is wrong; the JDR says this one is. |
| 8 | `:1645` | Assumption Verification: "Verified: A1, A2, A3, A4, A5, A6 (labels), A7, A9, A11, and A13. **Pending: A8**, A10, A12, A14, A15" | A8's own record (`:271-302`) reads `Status: Verified`. This roster contradicts it, and it is the roster the Finalization Gate reads. |
| 9 | `:1701-1713` | "**A8 is Pending and is the one record that still gates lock.** … Closure needs a §JD-4 disposition, not an edit here. **Recommended arm:** record the narrowing here …" | The recommended arm was adopted verbatim by the cluster gate. This paragraph proposes as future work a decision already made, and names a lock gate that no longer exists. |

**Severity.** This is not a design defect — the substance the RDR states is
correct and matches §JD-4's disposition on every point. It is a stale-citation
tail: the draft asserts *about its own state* something its own assumption
record and its own checklist contradict. Left standing, the Finalization Gate
reads `:1645` and blocks on a Pending A8 that is Verified, and a reader of
`:1701` concludes the RDR still gates on a peer decision that is closed.

### G3-2 — §JD-4's line anchor into this RDR points at the wrong clause

§JD-4 cites the narrowing clause as `0003:781-786`
(`docs/jdr/0001-resolve-kernel-seam.md:244-245`: "0003 records the narrowing
normatively (`0003:781-786`, already written)"). Lines 781-786 of this RDR are
inside the **presence-dimension projection** clause (A7's `exists` projection),
not the narrowing. The narrowing clause is at `:806` ("the lint promise narrows;
the runtime veto MUST NOT be weakened") and its withheld-claim form follows at
`:843`.

Likewise §JD-4:259 cites `0003:818-826` for the warning-tier question; that
range is the stale A8 peer-obligation block (finding G3-1 row 3), not a
warning-tier statement.

**Whose defect.** The stale anchors live in the JDR, not in this RDR, and this
skill does not edit peer documents. Recorded here because (a) a reader
following §JD-4's cite lands on the wrong clause, and (b) the fix for G3-1 row 3
deletes the block §JD-4:259 points at, which would leave a dangling anchor.
Both belong in the same route-back to the JDR. Not lock-blocking for this RDR:
the clause §JD-4 assigns here is present and correct; only the pointer to it is
off.

### G3-3 — Seven sites route A10/A12 to a cluster-reconcile venue that has already run and dispositioned them elsewhere

A10 and A12 carry `Plan` lines and Stage-6 dispositions that send them "to
cluster reconcile, with A8" — the same venue A8 named. That gate **has run**
(2026-08-22) and did not close them there. It dispositioned both as **standing
tolerances homed at RDR 0006's refine**
(`docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md:162-175`,
verbatim):

> ### Standing tolerances — F4, F5 (0003 A10, A12)
> Not re-decided here; they already have a home and a scheduled producer. Both
> are recorded obligations on 0006's refine …

with a table naming each open question and `0006's refine` as the Home, and both
answered `No`. The report's routing table records the same
(`:85-86`): "JOINT-DECISION → tolerance (0003 A10)" / "(0003 A12)".

The stale sites: `:332` (A10 Plan), `:336` (A10 Stage-6 disposition header),
`:345` (A10 "Travels with A8 and A12"), `:386` (A12 Plan), `:394` (A12 Stage-6
disposition header), `:396` (A12 "Travels with A8 and A10"), `:1374` and `:1368`
(the two Prerequisites checklist entries, both "Venue: cluster reconcile").

**Live consequence.** A reader following these plans schedules a second cluster
reconcile that would re-decide nothing — the gate already declined to re-decide
them — while the actual producer (RDR 0006's refine, which carries both in its
Direction) goes unnamed in the plan. The gate also recorded a substantive detail
none of these sites carries: for A12, RDR 0006 frames the check as *graph
reachability* while this RDR requires a *syntactic decision over declarations*,
and RDR 0006's refine must state which governs (`reconcile-report.md:174-176`).

## NOT-FOUND

None. Every cited `path::Symbol` and every peer/JDR line anchor cited by *this
RDR* resolves on `main`. (The two wrong anchors in G3-2 are cites *into* this
RDR from the JDR, and they resolve to real lines — they simply name the wrong
clause.)

## New rule with an existing sibling

None. See Inverse search.

## CONFIRMED

### Codebase claims (all resolve on `main`)

- `internal/resolve/resolve.go::Row` carries `RuleID` and `SourceLocator`
  (`resolve.go:172-174`), doc-commented as "the source identity RDR 0002
  requires every normalized row to retain". A5's shipped-identity claim holds.
- `Row.Guard string` (`resolve.go:185`) — still the opaque predicate string.
  Confirms A5's caveat that RDR 0007's atom reshape is specified, not shipped.
- `internal/resolve/resolve.go::GuardEvaluator` declared at `:91` with
  `Evaluate(guard string, view TagSet) GuardResult` (`:93`); called from
  `resolve.go::evaluateGuard` (`:425`) via `resolve.go::gate` (`:376`). The
  Investigation's "no implementation of it" claim still holds — `fixtureGuards`
  in `internal/resolve/fixtures_test.go` remains the only implementation.
- The desk-trace step 9 claim that "`RuleID` + `SourceLocator` ship on
  `internal/resolve/resolve.go::Row`" — confirmed at `:173-174`.

### JDR 0001 claims

- **§D1 fixes the parsed-atom row shape** — confirmed
  (`docs/jdr/0001-resolve-kernel-seam.md:86-90`, option (d) recommended): "Row
  carries a slice of atoms (key, operator token, literal, block) instead of a
  string. No reconstruction step". Matches A5.
- **§JD-4 is CLOSED 2026-08-22, naming this RDR the recording document, with
  RDR 0006 citing and minting no code** — confirmed verbatim (`:240-264`),
  including the "no warning category and no non-blocking tier" half and the
  consequent atom-level-field duty on RDR 0006 (`:260-263`). A8's evidence is
  accurate on every point. **This closes iteration 2's finding G2-1** — the
  citation G2-1 refuted has been repaired at its source, not just in this RDR.
- **§JD-13 homes the single-valued marker here** — confirmed (`:306-320`),
  including the reasoning A8-adjacent text reuses ("it is a property of a tag
  class … every sibling property already lives there") and the discharge
  routing ("Discharged at 0003's refine, with 0006's citation at its own").
  The RDR's §JD-13 account (`:1395-1401`) matches.
- **§JD-14 rules this RDR's escape-row reading governs** — confirmed
  (`:321-333`), including that RDR 0006 is split against itself (its
  Load-Bearing Decision at `0006:412-414` sides with this RDR against its own
  invariant 3). The Status line's "no edit owed" for §JD-14 is correct: the
  clauses §JD-14 quotes (`0003:771-779`) are the ones it upholds.

### Peer-RDR claims

- **RDR 0007's atom shape, normative** — confirmed (`0007:1249-1266`): "key,
  operator token, literal, block ∈ {all, unless}", with `Block` an exported
  named string type carrying `BlockAll`/`BlockUnless`. A5 and A14's dependency
  on per-atom `block` are both accurately grounded.
- **RDR 0007 requires this RDR to make the set-literal spelling canonical** —
  confirmed (`0007:1585-1595`): "0003's declaration MUST make that spelling
  canonical". A13's claim to own the canonical byte spelling is correctly
  assigned.
- **RDR 0002:317 requires source rule id and locator retention** — confirmed
  verbatim. A5 and A14's framing of what RDR 0002 does and does not guarantee
  (container-level combination, not per-atom `block`) is accurate.
- **RDR 0006 reuses `graph-unprovable-coverage`, already in its blocking table
  and MVV matrix** — confirmed at `0006:317` (blocking table) and `0006:643`
  (MVV matrix). A8's "reuse costs a widened trigger and no new test scaffolding"
  holds.
- **RDR 0006's finding contract carries no atom-level field** — confirmed
  (`0006:363-367`): stable code, model identity, severity, message, and "the
  source rule/context id or source span". A8's consequent-duty clause is
  accurate. (The *stale* framing of this same fact at `:818-826` is G3-1 row 3
  — the fact is right, the disposition around it is stale.)
- **RDR 0006 is `Draft`, demoted 2026-08-21, with §JD-4/§JD-13/§JD-14 and this
  RDR's A10/A12 in its refine Direction** — confirmed (`0006:8-16`). A10's and
  A12's Stage-6 dispositions ("RDR 0006 is now `Draft` and its Direction carries
  this item") are accurate.
- **The cluster gate report exists at the cited path** — confirmed:
  `docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md` (26 Aug 22
  12:30), alongside `critique-set.md` and three pairwise files. A8's "Gate
  evidence" cite resolves.

### Evidence-artifact claims

- `evidence/research/iter-2-projection-derivation.md` — exists; A7's cite
  resolves.
- `evidence/spikes/iter-2/a1-eval-harness/` — exists with `main.go` and
  `output.txt`; A1's and A14's corroborating cites resolve.

## Inverse search

The delta adds one genuinely new declaration field (the §JD-13 single-valued
marker) and widens one existing clause's trigger (the narrowing). Sibling-path
check for each:

```sh
rg -n 'single-valued|single_valued|singleValued' docs/rdr/000[2367]-*.md internal/
rg -n 'graph-unprovable-coverage|narrow' docs/rdr/0006-graph-lint-authority-and-guarantees.md
```

- **Single-valued marker** — searched RDR 0002 (authoring/normalization owner),
  RDR 0006 (the consumer), RDR 0007 (totality owner), and `internal/`.
  **Searched, none exists as a declaration.** RDR 0006 *consumes* it in four
  places (`:267` lint input contract, `:296` invariant, `:318` blocking code,
  `:643` MVV) but declares nothing; RDR 0002 has zero occurrences; RDR 0007's
  single hit (`:1554`) is the unrelated `single-valued Refusal.Guard text`;
  `internal/` has none. This is the unhomed-producer defect §JD-13 names, and
  homing it here creates no duplicate.
- **The narrowing's widened trigger** — searched RDR 0006 for an existing
  `guard_unevaluable` narrowing. Still none: `0006:363-366` narrows only for the
  non-finite-dimension case. The two triggers differ exactly as §JD-4 records,
  so this RDR's clause duplicates nothing.
