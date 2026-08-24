# Cluster Reconcile Report — 0002-0009 — Iteration 2

Model: claude-fable-5
Date: 2026-08-23
Iteration: **N = 2** (iteration 1 = `../report.md`, 2026-08-11, claude-opus-5).

## Membership

Cluster = the Final-and-unimplemented RDRs under `docs/rdr` at pass start. 0001
is `Implemented` (out of scope). All eight members were `Final` when this pass
began; 0008 leaves it `Draft` (below).

| RDR | Status at start | Revision examined | Moved since iter-1? |
| --- | --- | --- | --- |
| 0002 | Final | `5bed32a` (2026-08-23) | yes — refine/resolve/re-lock, folded JDR duties |
| 0003 | Final | `f1e9a57` (2026-08-22) | yes — re-walked, tag model rehomed, re-locked |
| 0004 | Final | `48cb2a4` (2026-08-21) | yes — §D3 clause + re-lock |
| 0005 | Final [§JD-8, §JD-9] | `76a9e67` (2026-08-12) | qualifier only; body June-frozen |
| 0006 | Final | `ae2f478` (2026-08-23) | yes — re-walked, re-locked |
| 0007 | Final | `21322fb` (2026-08-22) | yes — absorbed §D1/§D4, re-locked |
| 0008 | Final [§JD-5] | `ba95813` (2026-08-21) | yes — re-resolved A6/A11, re-locked |
| 0009 | Final [§JD-5] | `76a9e67` (2026-08-12) | qualifier only; body frozen 2026-08-11 |
| **Home — JDR 0001** | open | `21322fb` (2026-08-22) | created after iter-1; D1–D4 resolved; this pass edits it |

Iteration 1 recorded no revisions, so per the iteration contract **every check
re-ran**; nothing is CARRIED. Cluster membership was also widened to the
old-member pairs (0002-0003, 0002-0006, 0003-0006, 0002-0004, 0004-0005,
0005-0006, 0002-0005, 0003-0005) that iteration 1 carried from the `0001-0006`
pass: 0002/0003/0004/0006 have all re-locked since, so those carries were stale.

## Inputs run

| Check | Status | Output |
| --- | --- | --- |
| Whole-set critique | RE-SCANNED (second model of the dual draw: fable vs iter-1 opus) | `critique-set.md` — 33 rows, 13 iter-1 rows closed |
| 0007-0002, 0007-0003 | RE-SCANNED | `pairwise-0007-0002.md`, `pairwise-0007-0003.md` |
| 0007-0004, 0007-0005, 0007-0006 | RE-SCANNED | `pairwise-0007-000{4,5,6}.md` |
| 0007-0008, 0007-0009 | RE-SCANNED | `pairwise-0007-000{8,9}.md` |
| 0008-0002, 0008-0009, 0008-0005-0006 | RE-SCANNED | `pairwise-0008-*.md` |
| 0009-0002, 0009-0004, 0009-0005 | RE-SCANNED | `pairwise-0009-*.md` |
| 0002-0003, 0002-0006 | RE-SCANNED (new pairs) | `pairwise-0002-000{3,6}.md` |
| 0003-0006, 0002-0004 | RE-SCANNED (0003-0006 ledger = `../../0003-0006-0007/`) | `pairwise-0003-0006.md`, `pairwise-0002-0004.md` |
| 0004-0005, 0005-0006, 0002-0005, 0003-0005 | RE-SCANNED (new pairs) | `pairwise-*-0005.md`, `pairwise-0005-0006.md` |

21 pairs + 1 critique. Not run: 0004-0006 (0004 cites 0006 zero times, 0006
cites 0004 five times, none normative), 0003-0004 (9/2 references, all via
0002's carriage — covered by 0002-0003 and 0002-0004).

## The set-level finding that governs this report

Iteration 1 found *silence*. This pass finds the cure working — and its residue.
Every entry the home decided (§D1–§D4, §JD-4, §JD-13, §JD-14-coverage) is now
**RESOLVED in fenced text on both sides of every pair that raised it**: 0007
carries parsed atoms and no panic; 0002 is gate-then-count and emits the kernel
constants; 0004 refuses a partial read; 0003 records the narrowing and 0006
cites it. 18 of iteration 1's 36 rows close cleanly.

The residue has one shape: **members locked against answers the home never
gave.** 0002 decided recognized-totality "in the home's name" (§JD-10 cited,
entry still open); 0006 locked with A6/A10 booked as "scheduled edits on 0002
(`Draft`)" the day 0002 locked without them; 0003 and 0006 both wrote the
two-population overlap rule while the home still said the opposite; iteration
1's `<clear>` tolerance vanished when the JDR renumbered JD-12. Where the home
lagged its members, this pass moves the home. Where a member's fences contradict
a now-ratified answer, the ordinary path applies — once (0008).

## Origin ledger trace (iteration 1 → this pass)

| Iter-1 row | Now |
| --- | --- |
| 0007-0003 JD-1 ×2 | RESOLVED (§D1; 0007:1251-1271 defines atoms, 0003:845-855 cites) |
| 0007-0003 JD-4 | RESOLVED (0003:1162-1169, 0007:1656-1657) |
| 0007-0002 JD-2 ×2 | RESOLVED (0002:791-799 + Scenario 4; 0007:1451-1472) |
| 0007-0002 JD-3 / 0007-0009 JD-3 | **CLOSED at home this pass** — 0002:652-670 is the producer |
| 0007-0008 JD-4 ×2 | RESOLVED by §D4 (mis-routed to JD-4 in iter-1; home note added) |
| 0007-0009 JD-5 | STILL-OPEN; 0007 half textually reconciled; 0008/0009 carry qualifier |
| 0007-0009 JD-6 | RESOLVED (§D1; no panic anywhere in 0007) |
| 0007-0004 JD-7 ×2 / 0009-0004 JD-7 | RESOLVED (§D3; 0004:293-300 fenced + MVV 783-797) |
| 0007-0005 JD-8 ×3, 0009-0005 JD-8 ×2, 0008-0005 JD-8 | STILL-OPEN — **widened at home** (see JD-8) |
| 0007-0005 JD-9, 0008-0005 JD-9 | STILL-OPEN; home text decides provenance only, not the `--tag recognized=` classification |
| 0007-0006 JD-4 ×2 | RESOLVED (0006:891-898 cites 0003; reuses `graph-unprovable-coverage`; atom fields 0006:953-955) |
| 0008-0002 JD-10 ×3, 0008-0009 JD-10 | **ANSWERED** by 0002's re-lock, ratified at home → 0008's check failed (below) |
| 0008-0009 SPEC-DEFECT (A6/A11) | RESOLVED — `Input` count 5 in 0009, negative withdrawn, citations repaired (0008:571-579, 902-907) |
| 0008-0009 JD-5 | STILL-OPEN; base commit for the oracle still unnamed |
| 0008-0005/0006 JD-4 | SUPERSEDED — closed by 0002's load categories (0002:827-860), not by JD-4 |
| 0009-0002 JD-11 ×3 | RESOLVED (withdrawn → 0002 fixed at re-lock: 0002:726-732 dump = 0002:915-922 invariant; `writes = []` refused 0002:481-483) |
| 0009-0004 JD-12 (`<clear>`) | **LOST at the home** — restored as **JD-15**, still open |
| 0009-0004 JD-12 (unsanctioned escape writes) | SUPERSEDED — input class unconstructable under 0009:864-870 |
| 0009-0005 cosmetic (row-identity inverse) | STILL-OPEN, cosmetic, 0009-internal |

## Findings table

Severity/ownership as returned by the scans, verified by the parent where it
drives a disposition. `NET-NEW` = traces to no iteration-1 row (per the
iteration contract it demotes nobody on its own).

| Pair | Finding | TYPE | SEV | Own | Ledger | Disposition |
| --- | --- | --- | --- | --- | --- | --- |
| 0008-0002 | A9 `Verified` on 0002's superseded spike; block 2 says guard-position `recognized` refs "resolve to it", 0002 refuses them at load | contradiction | risks | single (0008) | JD-10 answered → check failed | **SPEC-DEFECT → 0008 Draft, Stage 4, STAGE-SCOPED** |
| 0008-0002 | Scenario 5 expects two failures; 0002 load is fail-fast | contradiction | risks | single (0008) | JD-10 | folded into the same SPEC-DEFECT |
| 0008-0002 | 0008 fenced "nothing in 0002 forbids empty alphabet entry"; 0002 now forbids | contradiction | cosmetic | single (0008) | JD-10 | same re-entry sweeps it |
| 0008-0002 | A4 census / Phase 2 rename scope stale (fixtures already renamed) | contradiction | cosmetic | single (0008) | JD-10 | same re-entry (re-verify A4) |
| 0008-0002 | Loader failure payload (three fields) + near-miss advisory channel: 0002 absorbed category, not carrier | gap | risks | joint | JD-10 residual | JOINT-DECISION → rides **JD-8** (carrier surface) |
| 0008-0002 | 0002 restates 0008's naming rule in fenced text | duplication | cosmetic | single (0002) | NET-NEW | recorded; cite-don't-restate at 0002's next touch |
| 0003-0006 | §JD-14 overlap half: home says escape×ordinary overlap-checked; both members say two populations | ANSWER-CONTRADICTS (home stale) | — | joint | F8 | **HOME CORRECTED** (JD-14 amended) |
| 0003-0006 | 0006 computes coverage per (group × rescuable class); 0003 states one union | gap | risks | joint | F8 rider | JOINT-DECISION → recorded under JD-14 correction; 0003 A17 closes on it |
| 0003-0006 | 0003 cites 0006 as Draft / "no atom-level field" | contradiction | cosmetic | single (0003) | F3 | CITATION REPAIR — pending 0003's next touch (sentence-entangled) |
| 0003-0006 | A10, A12, A19, A20-lint answered by 0006's re-lock | — | — | — | tolerances | RESOLVED; 0003's Pending stamps stale until next touch |
| 0003-0006 | A18 / A20 runtime half → "0007's next touch"; 0007 Final, silent | gap | risks | joint | NET-NEW | JOINT-DECISION → **JD-18** |
| 0007-0003 | conforming-view enforcer (same as above) | gap | risks | joint | NET-NEW | → JD-18 |
| 0002-0003 | Match/Guard routing by operator (0002) vs authored block (0003) | round-trip | blocks | joint | NET-NEW | JOINT-DECISION → **JD-16** |
| 0002-0003 | five type-model fields have no wire key | gap | risks | single (0002) | NET-NEW | → JD-17 (the key set is a shared interface) |
| 0002-0003 | declaration/kind and literal-outside-domain rejections have no pipeline site | gap | risks | joint | NET-NEW | → JD-17 |
| 0002-0006 | 0006 requires initial/terminal declarations; Final 0002 closed layout refuses them | gap | blocks | joint | NET-NEW | JOINT-DECISION → **JD-17** |
| 0002-0006 | invariant 5 rests on write-replaces rule no Final document states | gap | risks | joint | NET-NEW | → JD-17 |
| 0002-0006 / 0002-0004 / 0009-0004 | `<clear>` unreserved in normalized value; 0004 silent; 0006 cannot re-pair | round-trip / gap | blocks | joint | iter-1 JD-12 (lost) | JOINT-DECISION → **JD-15** (restored) |
| 0002-0004 | 0004 needs timeout/role/keys metadata; 0002 accessor entry is mode+path, strict-decoded | contradiction | risks | joint | NET-NEW | → JD-17 |
| 0002-0004 | "typed values" names no kind source | gap | risks | joint | NET-NEW | → JD-17 |
| 0002-0004 | 0004:567 says RequiresOwned producer unassigned | contradiction | cosmetic | single (0004) | JD-3 | CITATION REPAIR — pending 0004's next touch |
| 0007-0005 | 0007 Prerequisite reopens §D4 for a structured envelope field; home silent | gap | risks | joint | JD-8 | recorded in **JD-8 widening** |
| 0007-0005 | `flow-guard-unevaluable` "supplied facts" vs provenance-blind view | contradiction | risks | joint | JD-8 | STILL-OPEN under JD-8/JD-9; 0005 qualifier |
| 0007-0005 | `owned_state_unavailable` no code row | gap | risks | joint | JD-8 | STILL-OPEN; 0005 qualifier |
| 0004-0005 | seven 0004 refusal classes unmapped; `read_back_incomplete` collapses into `flow-write-readback-mismatch` (0004 MUST NOT) | gap / contradiction | risks | single (0005) | NET-NEW | JOINT-DECISION → **JD-8 widened** (see note) |
| 0002-0005 | ~20 load categories delegated to 0005; one code exists; escaped plan rendered as ordinary success; no path form | gap / round-trip | risks | single (0005) | NET-NEW | → JD-8 widened |
| 0003-0005 | predicate semantic kinds + declaration errors routed to 0005; no code | gap | risks | single (0005) | NET-NEW | → JD-8 widened |
| 0005-0006 | `--flow` only, no path form; 0006 A9 keyed to an unscheduled 0005 event; non-omitempty `Findings` on shared `CLIError` | gap / contradiction | risks | single (0005) / joint | NET-NEW / JD-8 | → JD-8 widened |
| 0005-0006 | both book the `respond.OK` text-branch edit | duplication | cosmetic | joint | NET-NEW | recorded |
| 0008-0009 | two error classes on one `Resolve` error; only 0009's `errors.Is`-classifiable | gap | risks | joint | NET-NEW | rides **JD-5** (the precedence answer must name the wrapping) |
| 0008-0009 | precedence between entry preconditions | gap | blocks | joint | JD-5 | STILL-OPEN; both qualifiers present |
| 0007-0009 | fixture sequencing (`escapeRow` Writes + RequiresOwned) | gap | risks | joint | JD-3 | recorded in JD-3 closure as implementation-order residue |
| 0007-0008 | JD-10 "total" flips `exists`/`eq` over `recognized` | gap | risks | joint | JD-10 | ANSWERED — total; 0008 re-resolve absorbs |
| 0009-0002 | 0009 quotes retired 0002 "row kind `escape`" wording ×3; 0009:1184 "implemented" | contradiction | cosmetic | single (0009) | JD-11 | CITATION REPAIR — pending 0009's next touch |
| 0009 | `Refusal.Guard` at 0009:421 | — | cosmetic | single (0009) | JD-12 | unfenced A3 evidence, true against `resolve.go:272` today; rides re-lock |
| 0007 | peer-status claims: 0002/0003/0006 called Draft ×11; "§JD-4 not closed" (594-598); A16/A21/A22/Prereqs say 0002 lacks clauses it has | contradiction | cosmetic | single (0007) | — | **CITATION REPAIR — 4 pure status claims repaired in place** (0007:124, 159, 1711, 1713); the sentence-entangled rest (463, 600, 613, 645, 2172, A12/A16/A21/A22 evidence, Phase 4 handoff "with 0003 in Draft") pend 0007's next touch |
| 0006 | "RDR 0002 is `Draft` … scheduled edit on an open peer" ×7 (14, 256, 280, 380, 391, 1421, 1425) | contradiction | cosmetic | single (0006) | — | CITATION REPAIR — the *fact* is JD-17's substance; the wording pends 0006's next touch |
| 0008 | "0009 still Draft" (1668); stale 0007 quotes | contradiction | cosmetic | single (0008) | — | swept by the 0008 re-entry |
| 0002 | restates §D2 gate-then-count in prose (325-329, 384-393) against its own fenced "MUST NOT be restated here" (791-793) | duplication | cosmetic | single (0002) | JD-2 | recorded; §D2's "lands in 0002 … restates" sanctions it; tidy at next touch |
| 0004 | 0004:468-472 guard-only absent key → escapable `no_match` (pre-§D4 flow) | contradiction | cosmetic | single (0004) | NET-NEW | recorded; unfenced narration; rides next touch |
| 0003/0006 | `graph-product-too-large` second blocking code; 0006 cites 0003::A21 for a fenced rule | duplication / contradiction | cosmetic | — | NET-NEW | recorded |

**DEFER-TO-IMPLEMENTATION: none.** No finding cleared all three conditions
without a cheaper disposition available: the cosmetic items are citation
repairs (facts of record), and every risks/blocks item either changes a clause's
meaning or sits in fenced text. No prior `DEFERRED` row exists to repeat.

## Dispositions

### SPEC-DEFECT — RDR 0008 → Draft (the one demotion)

Ledger-traced: JD-10 was an open tolerance; RDR 0002's re-lock answered it in
fenced text (totality via outcome lifting, guard-position `recognized` refused
at load, fail-fast load). The home entry still read open — a member deciding
in the home's name (critique C-20). This gate **ratifies the answer at the
home** (0002 owns normalization; the answer is consistent with §D4's
kernel-decided presence; 0008 itself defers "not this RDR's to rule on"), then
runs the sibling's scoped answer-vs-fences check. 0008 fails it:

- A9 `Verified` (Method: Spike) rests on 0002's *pre-lift* spike ("no outcome
  field to receive a predicate … none produces `Row.Outcome`"); the current
  spike (`main.go::liftOutcome`, `output.txt`) lifts it.
- Block 2 (fenced, 0008:1283-1290): guard-position reserved-name references
  "need no separate reserved-key check … resolve to it" — 0002 (fenced,
  0002:623-641) refuses them at load. Meaning contradiction, fenced: not
  deferrable.
- Scenario 5 (0008:2340-2352) pins "two failures, not one"; 0002 (fenced,
  0002:455-463) is fail-fast.

**Less foundational**: 0008 — 0002 is the producer of normalization and is
internally consistent. **Re-entry**: Stage 4 (Resolve). **Scope**:
STAGE-SCOPED — bounded `re-verify A9, A4`; block 2 and scenario 5 restate to
the lifted-outcome / fail-fast shape; the approach (reserve the name, enforce
at the kernel input boundary) holds. FULL-FLOW not earned. Flip, README row,
and re-entry note are **done** (0008:14; `## Refinement Context` appended).
0008 re-locks carrying §JD-5 / §JD-8.

Same defect class as iteration 1's A6/A11 (an assumption `Verified` against a
sibling's superseded text). Iteration 1's row is closed; this is a distinct
assumption on a distinct sibling revision — not a re-deferral and not a plank.

### HOME CORRECTIONS (JDR 0001) — done this pass

- **§JD-3 CLOSED** — 0002:652-670 is the producer; sequencing residue noted.
- **§JD-10 ANSWERED** — ratified from 0002; 0008's failed check recorded.
- **§JD-14 corrected** — overlap half amended to the two-population reading
  both members already lock (0003:1206-1217, 0006:901-903). The home lagged
  both siblings; no member moves.
- **§JD-4** — note that the two 0007×0008 rows iteration 1 routed here are
  answered by §D4.
- **§JD-8 widened** — full code table + carrier questions (0007's §D4
  reversal request, 0006's `Findings`).

### JOINT-DECISIONS — four new homes (JD-15…JD-18), tolerances recorded

| ID | Joint decision | Pairs | Siblings qualified |
| --- | --- | --- | --- |
| JD-15 | `<clear>` at the write accessor (restored iter-1 tolerance) | 0009-0004, 0002-0004, 0002-0006 | 0002, 0004, 0009 |
| JD-16 | Match/Guard routing key — operator vs authored block | 0002-0003 | 0002, 0003 |
| JD-17 | What the closed TOML layout must additionally spell (initial/terminal, write-replaces, accessor metadata, type-model keys) | 0002-0006, 0002-0004, 0002-0003 | 0002, 0003, 0004, 0006 |
| JD-18 | Conforming-view enforcer (runtime half of declaration conformance) | 0007-0003, 0003-0006 | 0003, 0007 |

Also qualified on existing entries: 0007 → §JD-8 (its per-atom payload carrier
is the open question); 0009 → §JD-8 (its row-identity field). All qualifiers
name the open question (Status lines + README rows updated — done).

**On 0005.** Four scans return single-RDR risks-impl gaps on 0005 (its
June-frozen code table vs every producer's failure family). All are NET-NEW —
the pairs were never scanned in iteration 1 — so under the iteration contract
they demote nobody on their own. They are also exactly JD-8's question at full
width, so they are homed there (widened) and 0005 keeps its existing §JD-8
tolerance. Stated plainly: **0005 is the member JD-8's answer re-walks**, and
it is the one member never re-entered since June. The whole-set critique
(`critique-set.md` C-6, C-23) says the same.

### CITATION REPAIRS

Applied in place, facts of record only (README index): 0007:124, 159, 1711,
1713. Recorded and pending each RDR's next touch where the stale claim is
sentence-entangled (0007 ×7 more, 0006 ×7, 0003 ×2, 0004 ×1, 0009 ×4, 0002 ×1)
— listed in the findings table and in each `pairwise-*.md` peer-status sweep.
None is fenced; none changes meaning.

## Verdict

**NOT RECONCILED.**

- One SPEC-DEFECT: **RDR 0008 → Draft**, re-entry **Stage 4**, scope
  **STAGE-SCOPED** (re-verify A9, A4; block 2; scenario 5).
- Home corrected on JD-3 (closed), JD-10 (answered), JD-14 (overlap half),
  JD-4 (routing note), JD-8 (widened).
- Four new joint decisions homed (JD-15…JD-18); nine standing tolerances:

| Tolerance | Home | Open question | Answered? |
| --- | --- | --- | --- |
| 0008, 0009 | §JD-5 | precedence of the two `Resolve`-entry preconditions; error wrapping | No |
| 0005, 0007, 0008, 0009 | §JD-8 | full code table / exit-3 / envelope carrier | No (widened) |
| 0005, (0008) | §JD-9 | `--tag recognized=` classification | No (provenance half answered) |
| 0002, 0004, 0009 | §JD-15 | `<clear>` semantics at write + read-back | No |
| 0002, 0003 | §JD-16 | Match/Guard routing key | No |
| 0002, 0003, 0004, 0006 | §JD-17 | wire keys the closed layout must add | No |
| 0003, 0007 | §JD-18 | conforming-view enforcer | No |

No deferrals. The cluster does not implement over the open SPEC-DEFECT; the
seven members left Final are typed non-defect and may proceed once 0008
re-locks — except that JD-16 and JD-17 are `blocks-impl` on 0002's own
normalizer, so 0002 in particular should not start Stage 8 before those two
entries answer.

## Required next action

1. **`/rdr-resolve 0008`** → `/rdr-finalize 0008` (STAGE-SCOPED; carries
   §JD-5/§JD-8 forward).
2. **Answer JD-15, JD-16, JD-17 at the home** (author call; each is a
   paragraph). JD-16 and JD-17 gate 0002's implement; JD-15 gates 0004/0009.
3. **Re-run this gate at N=3** scoped to 0008's re-lock and to the siblings
   holding tolerances the home answers. N=3 is the last iteration that may
   demote.

## Observations outside this gate's authority

- **Code has not moved.** `internal/resolve/resolve.go:185` still declares
  `Guard string`; no commit under `internal/` or `cmd/` since 2026-08-09. Every
  atom-based clause across the cluster is written against a type nobody has
  written. The critique's ratio is now ~16:1 markdown to Go.
- **0004's Prerequisites** remain unchecked under `Final` (iter-1 observation,
  unchanged; critique C-22).
- **0009** has had no content re-lock since 2026-08-11 and predates §D1/§D4;
  its `Refusal.Guard` citation is unfenced and currently true. When 0007
  Phase 2 lands it goes stale — 0009's implement gate should re-read §D4.
- **Whole-set critique delta**: of iteration 1's rows, 13 closed, 7
  superseded, 13 still open — everything that closed was decided by §D1–§D4;
  everything open sits on a blank/open home entry or on 0005/0009.

## Review Gate

- **Cluster correct?** Yes. Eight Final-and-unimplemented members at start;
  0001 excluded. Widened to the old-member pairs iteration 1 carried, since all
  four had re-locked — the carries were stale, and the widening produced JD-16
  and JD-17, both `blocks-impl`.
- **N>1 anchored to history?** Every check is RE-SCANNED (iteration 1 recorded
  no revisions; every member moved). Revisions recorded above for N=3. The
  demotion traces to an open ledger entry (JD-10) through the answered-
  tolerance path; every NET-NEW single-RDR finding is recorded without
  demotion. Cap holds (N=2).
- **SPEC-DEFECT scope narrowest?** Yes — bounded `re-verify A9, A4`, approach
  intact → STAGE-SCOPED, Stage 4.
- **Right RDR demoted?** 0008; 0002 is the producer and consistent.
- **No design-body edit of a Final RDR?** Correct. Edits made: Status lines,
  README rows, 0008's flip + re-entry note, four pure peer-status citation
  repairs in 0007 (facts of record), and the home.
- **Each JOINT-DECISION genuinely joint?** JD-15/16/17/18 each sit between a
  producer and a consumer where neither is solely wrong (silence or a shared
  interface). The 0005 findings are single-RDR in wording but are JD-8's
  question at full width, and are NET-NEW — homed, not demoted, with that
  reasoning stated.
- **Home form?** One paragraph per new entry; derivation stays in the
  `pairwise-*.md` files. Restatements found (0002 restating §D2, 0002
  restating 0008's naming rule) are recorded as cosmetic findings.
- **Facts of record?** Peer status → README index; disagreeing peer figures
  typed CITATION REPAIR, four repaired in place, the rest recorded.
- **Has any home answered since the last report?** Yes — D1–D4, JD-4, JD-13,
  JD-14 (checked on every sibling pair, all RESOLVED except the JD-14 overlap
  half, corrected at the home), and JD-10 (answered by 0002, ratified here;
  0008's check failed → the SPEC-DEFECT).
- **Deferrals?** None; nothing fenced or meaning-changing deferred; no prior
  `DEFERRED` row.
- **Both prompts run?** Yes — `critique-set.md` (second model of the dual
  draw) and 21 `pairwise-*.md` files, all under this directory.
