# Cluster Reconcile Report — 0002-0009 — Iteration 3

Model: claude-fable-5
Date: 2026-08-24
Iteration: **N = 3** (iter-1 `../report.md` 2026-08-11 opus; iter-2
`../iter-2/reconcile-report.md` 2026-08-23 fable). **N=3 is the last iteration
that may demote.**

## Membership

Cluster = the Final-and-unimplemented RDRs under `docs/rdr` at pass start (0001
Implemented, excluded). All eight were `Final` at start; this pass leaves
0002, 0003, 0004 `Draft` (below).

| RDR | Status at start | Revision examined | Moved since iter-2 (`rev` recorded there)? |
| --- | --- | --- | --- |
| 0002 | Final [§JD-15/16/17] | `526c481` | qualifier only (iter-2 `5bed32a`) |
| 0003 | Final [§JD-16/18] | `526c481` | qualifier only (`f1e9a57`) |
| 0004 | Final [§JD-15/17] | `526c481` | qualifier only (`48cb2a4`) |
| 0005 | Final [§JD-8/9] | `76a9e67` | no |
| 0006 | Final [§JD-17] | `526c481` | qualifier + lock note only (`ae2f478`) |
| 0007 | Final [§JD-8/18] | `526c481` | 4 pure peer-status repairs (`21322fb`) |
| 0008 | Final [§JD-5/8] | `4581223` | **yes** — Stage 4 re-entry + re-lock (`ba95813`) |
| 0009 | Final [§JD-5/8/15] | `526c481` | qualifier only (`76a9e67`) |
| **Home — JDR 0001** | open | `4f9639c` | **yes** — §D5/§D6/§D7 decided → JD-15, JD-16, JD-17 ANSWERED (`21322fb`) |

Per the iteration contract: re-scan every pair with a moved member (0008's four),
re-run the whole-set critique (≥1 member moved), and — the answered-tolerance
exception — run the scoped answer-vs-fences check on every sibling of JD-15/16/17
even though none of their texts moved. Everything else CARRIES from iter-2.

## Inputs run

| Check | Status | Output |
| --- | --- | --- |
| Whole-set critique | RE-SCANNED | `critique-set.md` (27 iter-2 rows open, 3 closed, 3 superseded, 12 new) |
| 0008-0002 | RE-SCANNED (0008 moved) | `pairwise-0008-0002.md` |
| 0008-0009 | RE-SCANNED (0008 moved) | `pairwise-0008-0009.md` |
| 0008-0005 / 0008-0006 | RE-SCANNED (0008 moved) | `pairwise-0008-0005-0006.md` |
| 0007-0008 | RE-SCANNED (0008 moved) | `pairwise-0007-0008.md` |
| 0002 × home (JD-15/16/17) | ANSWER CHECK | `answer-check-0002.md` |
| 0003 × home (JD-16/17) | ANSWER CHECK | `answer-check-0003.md` |
| 0004 × home (JD-15/17) | ANSWER CHECK | `answer-check-0004.md` |
| 0006 × home (JD-17, JD-15 by citation) | ANSWER CHECK | `answer-check-0006.md` |
| 0009 × home (JD-15) | ANSWER CHECK | `answer-check-0009.md` |
| 0002-0003, 0002-0004, 0002-0006, 0009-0004 | CARRIED from iter-2 (no member text moved); their JD-15/16/17 rows are re-examined via the answer checks above | `../iter-2/pairwise-*.md` |
| 0007-0002, 0007-0003, 0007-0004, 0007-0005, 0007-0006, 0007-0009 | CARRIED from iter-2 | `../iter-2/pairwise-0007-*.md` |
| 0009-0002, 0009-0005, 0004-0005, 0005-0006, 0002-0005, 0003-0005, 0003-0006 | CARRIED from iter-2 (JD-18 unanswered; no member moved) | `../iter-2/pairwise-*.md` |

## The set-level finding that governs this report

Iteration 2 found members locked against answers the home never gave, and
moved the home. This pass finds the converse: **the home answered, and three
members' fenced text predates the answers.** §D6 (block-keyed routing) and §D7
(closed layout) overrule 0002's fenced routing clause and layout enumeration;
§D7(iii) overrules one fenced route in 0003; §D5 (reserved `<clear>`) overrules
0004's unfenced read-back prose in meaning. These are the ordinary-path
SPEC-DEFECTs the JOINT-DECISION disposition promised when it let the siblings
proceed under tolerance: the answer landed, the check ran, three failed. The
other two siblings (0006, 0009) pass and clear their qualifiers. 0008's
iteration-2 demotion is discharged on every pair.

## Origin ledger trace (iter-2 table → this pass)

| Iter-2 row | Now |
| --- | --- |
| 0008-0002 SPEC-DEFECT (A9/A4, block 2, scenario 5, empty alphabet, census) | **RESOLVED** — all seven discharge conditions verified on four scans (`0008:464-510, 802-839, 1355-1364, 2446-2454, 2247-2249`; no Refinement Context; no DEFERRED) |
| 0008-0002 payload + advisory carrier → JD-8 | STILL-OPEN; 0008 now names it (`2191-2196`); **home JD-8 gains the clause** |
| 0008-0002 0002 restates 0008's naming rule | STILL-OPEN, cosmetic; swept by 0002's re-entry |
| 0003-0006 JD-14 / A17 / A18-A20 rows | CARRIED (iter-2 verdicts) |
| 0003-0006 / 0007-0003 → JD-18 | STILL-OPEN; 0003 and 0007 carry the qualifier |
| 0002-0003 → JD-16 | **ANSWERED (§D6)** → 0003 CONSISTENT (qualifier cleared); **0002 CONTRADICTS, fenced** → SPEC-DEFECT |
| 0002-0003 / 0002-0006 / 0002-0004 → JD-17 | **ANSWERED (§D7)** → 0006 CONSISTENT (cleared; landings deferred); 0004 fences CONSISTENT, `0004:371` unfenced contradiction (folded into 0004's re-entry); **0003 one fenced clause CONTRADICTS** → SPEC-DEFECT; **0002 CONTRADICTS, fenced** → SPEC-DEFECT |
| 0002-0006 / 0002-0004 / 0009-0004 → JD-15 | **ANSWERED (§D5)** → 0009 CONSISTENT (cleared); 0002 fences CONSISTENT, unfenced Round-Trip false (folded into 0002's re-entry); **0004 unfenced read-back prose contradicts in meaning** → SPEC-DEFECT |
| 0007-0005, 0004-0005, 0002-0005, 0003-0005, 0005-0006 → JD-8 widened | STILL-OPEN (0005 byte-identical; home unchanged) |
| 0008-0009 → JD-5 | STILL-OPEN, **sharpened at home** (0009's fenced MUSTs admit only 0009-first; wrapping half now written into the entry) |
| 0007-0008 JD-10 | RESOLVED (0008 absorbs the total answer) |
| 0009-0002 / 0009 `Refusal.Guard` / 0009 "implemented" | STILL-OPEN, cosmetic, unfenced — pend 0009's next touch |
| 0007 peer-status ×7 sentence-entangled | STILL-OPEN, cosmetic — pend 0007's next touch |
| 0006 "0002 is Draft" ×7 | STILL-OPEN, cosmetic — written as 0006 `deviations.md` D2 |
| 0008 "0009 still Draft" / stale 0007 quotes | RESOLVED by the re-lock |
| 0002 restates §D2; 0004:468-472; 0003/0006 `graph-product-too-large` | CARRIED, cosmetic |

No prior `DEFERRED` row exists (iter-2 deferred nothing), so no repeat is possible.

## Findings table

| Pair | Finding | TYPE | SEV | Own | Ledger | Disposition |
| --- | --- | --- | --- | --- | --- | --- |
| 0002 × home | fenced routing clause `0002:693-717` operator-keyed; block vocabulary two-valued `539-545`; `in`-expansion / suffix totality `623-631`, `751-755` — vs §D6 block-keyed | contradiction | blocks | single (0002) | JD-16 answered → check failed | **SPEC-DEFECT → 0002 Draft, Stage 3, STAGE-SCOPED (re-verify A1, A9, A12)** |
| 0002 × home | fenced layout `0002:430-441` `[accessors.<id>]` mode+path; tag-side `accessor` `813-814`; all §D7 landings absent; fixture non-conforming — vs §D7 | contradiction | blocks | single (0002) | JD-17 answered → check failed | same SPEC-DEFECT |
| 0002 × home | unfenced Round-Trip "`<clear>` not reserved" `936-938`, `956`, `949-950`; load category absent — vs §D5 | contradiction | cosmetic | single (0002) | JD-15 answered | folded into the same re-entry (not a demotion driver on its own) |
| 0003 × home | fenced `0003:1065-1067` routes declaration / literal-domain rejections "onto RDR 0006 findings"; §D7(iii): 0002 load categories, 0006 mints nothing | contradiction | risks | single (0003) | JD-17 answered → check failed | **SPEC-DEFECT → 0003 Draft, Stage 3, RE-LOCK-ONLY (re-verify none)** |
| 0003 × home | participation / can-refuse clauses already block-keyed; `min`/`max` asserted not cited; illustrative `optional` spelling | — | cosmetic | single (0003) | JD-16/17 | CONSISTENT; sweeps ride the re-entry |
| 0004 × home | unfenced read-back "values are present" `0004:215`, `511` vs §D5 read-back-absent; zero `<clear>` clauses; `0004:371` "cannot be rebound" vs §D7(ii) | contradiction (meaning) | risks | single (0004) | JD-15/17 answered → check failed | **SPEC-DEFECT → 0004 Draft, Stage 3, RE-LOCK-ONLY (re-verify none)** — meaning contradiction, so not deferrable despite being unfenced |
| 0006 × home | every fenced clause CONSISTENT with §D7/§D5; A6/A10 still Pending; terminal-non-owned finding and terminal spelling absent | gap | cosmetic | single (0006) | JD-17 answered → check passed | qualifier CLEARED; **DEFERRED** → `0006/artifacts/deviations.md` D1 (checks: A6/A10 flip by citation; lint test for terminal over non-owned tag; prose cite) |
| 0006 | "RDR 0002 is Draft / scheduled edit" ×7 + 141-142, 253-256, 813 | contradiction | cosmetic | single (0006) | iter-2 | CITATION REPAIR (artifact of record: README) — written as `deviations.md` D2 for the next touch |
| 0009 × home | every fenced clause CONSISTENT with §D5; home's "0009 unchanged" verified | — | — | — | JD-15 answered → check passed | qualifier CLEARED (§JD-5, §JD-8 kept) |
| 0007-0008 | scenario 3 `0008:2390-2413` captures a guard-side **view**; 0007 fenced `1282-1290` "`Evaluate(atom, value)` … never sees the view" | contradiction | risks | single (0008) | NET-NEW (critique C-12; never in a report table) | **DEFERRED** → `0008/artifacts/deviations.md` D1 — (a) MVV text outside every fence, (b) compile/test at Phase 1, (c) contract unchanged |
| 0008-0002 | A10 evidence `0008:882-918` rests on `[accessors.<id>]` `{Mode,Path}` retired by §D7 | contradiction | cosmetic | single (0008) | NET-NEW | **DEFERRED** → `0008/artifacts/deviations.md` D2 (grep after 0002 re-locks) |
| 0008-0002 | 0008 block 2 restates 0002's outcome-binding rule; Phase 3 "(two naming failures)" vs scenario 5 "exactly one" | duplication | cosmetic | single (0008) | NET-NEW | recorded; 0008's next touch |
| 0008 | `0008:1772` "`<clear>` syntactically un-declarable" vs §D5 load category | contradiction | cosmetic | single (0008) | NET-NEW | recorded in deviations D2 note |
| 0008-0009 | JD-5 precedence: 0009 fenced `864-867`, `926-929` admit only 0009-first; 0008 licenses either; wrapping differs | gap | blocks | joint | JD-5 | STILL-OPEN → **home sharpened**; both qualifiers present |
| 0008-0005 | `reserved_tag_key` no code/exit row | gap | blocks | joint | JD-8 | STILL-OPEN; tolerance |
| 0008-0005 | `--tag recognized=` classification (programmer mistake vs `GroupUserEnv`) | contradiction | blocks | joint | JD-9 | STILL-OPEN → **home records the classification arm; 0008 qualifier gains §JD-9** |
| 0008-0002/0006 | payload + near-miss advisory has no carrier | gap | risks | joint | JD-8 | STILL-OPEN → **home JD-8 names 0008's carrier** |
| 0008-0006 | dead-rule / scenario-14 fixture needs a `recognized`-absent atom 0002 refuses at load | gap | risks | joint | NET-NEW | recorded; routed into 0002's re-entry note (state the authorable form) |
| home | §D7(ii) "key served by zero readers refused" would refuse `recognized` / `--tag` keys (§JD-9) | contradiction (home-internal) | risks | home | NET-NEW (critique N-4) | **HOME NOTE** — recorded as open at (ii), owner 0002 at re-entry |
| README | 0008 row bare `Final` vs Status `[§JD-5, §JD-8]` | citation | cosmetic | README | NET-NEW | **CITATION REPAIR — done** (row now `§JD-5, §JD-8, §JD-9`) |
| 0007 / 0009 | sentence-entangled stale peer-status claims | citation | cosmetic | single | iter-2 | pend next touch (unchanged) |

## Dispositions

### SPEC-DEFECT — RDR 0002 → Draft (STAGE-SCOPED)

Ledger-traced through the answered-tolerance path: JD-16 and JD-17 were open
tolerances 0002 carried; the home answered both on 2026-08-24; the scoped check
(`answer-check-0002.md`) fails on fenced text — the routing clause is
operator-keyed where §D6 is block-keyed, and the layout/accessor fences spell
the shape §D7 retires. **Less foundational**: 0002 yields to the home (the
authority the tolerance named). **Re-entry**: Stage 3 (Refine — restate the
fences to the answers), Stage 4 re-verifies **A1, A9, A12**, forward through
6 → 7. **Scope**: STAGE-SCOPED — the approach (sparse TOML + normalizer)
holds; a bounded assumption set moves; FULL-FLOW not earned. The unfenced
§D5 residue and the cosmetic sweeps ride the same re-entry. Flip, README row
and re-entry note **done** (`0002:9`; `## Refinement Context` appended).

### SPEC-DEFECT — RDR 0003 → Draft (RE-LOCK-ONLY)

One fenced clause (`0003:1065-1067`) contradicts §D7(iii) on where the two
rejection rules surface. Rule substance, phases and ownership already agree,
so this is a cross-ref fix: **re-verify none**, Stage 3 reword → Stage 7
re-lock. Fenced text never defers, hence a demotion rather than a deferral.
0003 carries **§JD-18** forward (still open); §JD-16 is cleared (consistent).
Flip, README row and re-entry note **done**.

### SPEC-DEFECT — RDR 0004 → Draft (RE-LOCK-ONLY)

All fences CONSISTENT; the contradiction is in *meaning* (read-back asserts
presence at `0004:215`, `511`; §D5 asserts absence for a `<clear>` write) —
the guardrail says a meaning contradiction never defers, fenced or not. The
§D5 landing clauses and the `0004:371` repair are additive: **re-verify none**,
Stage 3 → Stage 7. The still-unchecked Prerequisites (`0004:771-780`) are
named for the re-lock gate. No qualifier carries forward. Flip, README row
and re-entry note **done**.

### Tolerances CLEARED (answered, consistent)

- **0006 / §JD-17** — 16 clauses checked, 0 contradict; qualifier replaced
  by the answered note; the additive landings are **DEFERRED** (below).
- **0009 / §JD-15** — 9 clauses checked, 0 contradict; the home's "0009
  unchanged" holds; §JD-15 dropped from the qualifier, §JD-5/§JD-8 kept.
- **0003 / §JD-16** — consistent; cleared in the Draft qualifier text.

### DEFER-TO-IMPLEMENTATION (three entries, two RDRs) — written

| RDR | Entry | Mechanical check | Conditions |
| --- | --- | --- | --- |
| 0006 | `artifacts/deviations.md` D1 — §D7 additive landings (A6/A10 flip by citation; terminal-over-non-owned finding; terminal spelling) | grep of A6/A10 stamps after 0002 re-locks; a lint test for a terminal context over a non-owned tag | (a) unfenced/absent, (b) grep + test, (c) additive |
| 0008 | `artifacts/deviations.md` D1 — scenario 3's view-capturing seam vs 0007's fenced `Evaluate(atom, value)` | Phase-1 test authored at 0007's per-atom seam capturing `value`; escalate if same-view cannot be asserted | (a) MVV text outside fences, (b) compile/test, (c) contract unchanged |
| 0008 | `artifacts/deviations.md` D2 — A10 evidence on the retired `[accessors.<id>]` layout | grep of 0002's re-locked layout + fixture for a `keys` entry naming `recognized` | (a) assumption evidence, (b) grep, (c) conclusion survives |

Nothing fenced and nothing meaning-changing was deferred; no prior `DEFERRED`
row exists to repeat. 0006's stale-status sweep is recorded as D2 (citation
repair), not a deferral.

### HOME edits (JDR 0001) — done

- JD-15 / JD-16 / JD-17: sibling check outcomes recorded per entry.
- JD-5: sharpened (0009-first is the only order 0009's fences admit; wrapping
  half written in).
- JD-8: third carrier (0008's payload + advisory) named.
- JD-9: classification arm recorded as still open; siblings 0005, 0008.
- D7(ii): open note — the reader-binding refusal needs a provenance scope
  (`recognized`, `--tag`); owner 0002 at re-entry.

### CITATION REPAIRS — done

README 0008 row → `Final [joint decision → JDR 0001 §JD-5, §JD-8, §JD-9]`;
0008 Status gains §JD-9 (its own text says the repair lands there). README rows
for 0002/0003/0004/0006/0009 updated with the flips and clearings.

## Verdict

**NOT RECONCILED.**

- Three SPEC-DEFECTs, all on the answered-tolerance path:
  - **0002 → Draft**, Stage 3, **STAGE-SCOPED** (re-verify A1, A9, A12).
  - **0003 → Draft**, Stage 3, **RE-LOCK-ONLY** (re-verify none; carries §JD-18).
  - **0004 → Draft**, Stage 3, **RE-LOCK-ONLY** (re-verify none).
- 0008's iteration-2 demotion is **discharged**.
- Standing tolerances after this pass:

| Tolerance | Home | Open question | Answered? |
| --- | --- | --- | --- |
| 0008, 0009 | §JD-5 | precedence of the two `Resolve`-entry preconditions; error wrapping | No (sharpened) |
| 0005, 0007, 0008, 0009 | §JD-8 | full code table / exit-3 / envelope carriers (incl. 0008's) | No |
| 0005, 0008 | §JD-9 | `--tag recognized=` classification | No (provenance half answered) |
| 0003, 0007 | §JD-18 | conforming-view enforcer | No |
| ~~0002, 0004, 0009~~ | §JD-15 | — | **Yes** — 0009 cleared; 0002/0004 re-enter |
| ~~0002, 0003~~ | §JD-16 | — | **Yes** — 0003 cleared; 0002 re-enters |
| ~~0002, 0003, 0004, 0006~~ | §JD-17 | — | **Yes** — 0006 cleared (landings deferred); 0002/0003/0004 re-enter |

- Deferred: 0006 D1; 0008 D1, D2 (checks above).
- Implementation ordering: no member implements over the open SPEC-DEFECTs.
  0005/0006/0007/0008/0009 are typed non-defect and may start Stage 8 once
  0002 re-locks — 0002 is the producer every other member's wire format and
  fixtures read, so its re-lock is the practical gate for the set. 0006's and
  0008's deferred checks run at their Stage 8.

## Cap

N=3 was the last iteration permitted to demote; it did. If the three re-locks
touch a cross-RDR seam (0002's under STAGE-SCOPED likely will — new layout,
new routing), the next gate pass is N=4 and **may not demote**: with findings
still open it stops with `stopped:cluster-flapping` and the decomposition
question. To make N=4 clean: 0002 re-locks first, 0003/0004 re-lock against
the re-locked 0002, and any new joint question is homed at the JDR *before*
the gate re-runs. Owed answer-vs-fences checks survive a stop regardless.

## Required next action

1. **`/rdr-refine 0002`** → `/rdr-resolve 0002` (A1, A9, A12) → `/rdr-reconcile 0002`
   → `/rdr-finalize 0002`. Settle the D7(ii) provenance-scope note on the way.
2. **`/rdr-refine 0003`** → `/rdr-finalize 0003` (carrying §JD-18).
3. **`/rdr-refine 0004`** → `/rdr-finalize 0004` (Prerequisites resolved at the gate).
4. Re-run this gate at N=4 scoped to the three re-locks (no demotion possible).

## Observations outside this gate's authority

- **Code has not moved** (critique C-31): 69,271 markdown lines vs 3,971 Go;
  zero `internal/` commits since 2026-08-09. Every clause across the cluster
  is written against types nobody has written.
- 0005 remains the one member never re-entered since June; JD-8's answer is
  what re-walks it.
- The 0007 and 0009 sentence-entangled peer-status claims pend their next
  touches (unchanged from iter-2).

## Review Gate

- **Cluster correct?** Yes — eight Final-and-unimplemented members at start,
  0001 excluded; the same set as iter-2, no missed peer (every member cites
  or is cited by another; the home lists all eight).
- **N>1 anchored to history?** Every check is a RE-SCANNED, ANSWER CHECK, or
  CARRIED row naming its source; re-scans are exactly the pairs whose member
  moved (0008's four) plus the answered-tolerance checks the contract
  requires; each demotion traces to an open ledger entry (JD-16, JD-17,
  JD-15) via the answered path. NET-NEW findings demoted nobody (C-12, A10,
  N-4, F6 recorded/deferred/homed).
- **Cap held?** N=3 demotes; the report says what N=4 may not do.
- **Right RDR demoted, narrowest scope?** Each yields to the home; 0002
  STAGE-SCOPED with a named `re-verify` set; 0003 and 0004 RE-LOCK-ONLY with
  `re-verify none` — no `none` defect scoped above RE-LOCK-ONLY, no FULL-FLOW.
- **No design-body edit of a Final RDR?** Correct — edits: Status lines,
  README rows, three flips + re-entry notes (now Drafts), two `deviations.md`
  files, the home.
- **Each JOINT-DECISION genuinely joint?** No new joint decisions this pass;
  the three still-open entries (JD-5/8/9) were sharpened, not re-typed, and
  every sibling's qualifier names the home and the question. 0008 gained the
  §JD-9 qualifier its own text asks for.
- **Home form?** One paragraph per decision; this pass appended one-to-three
  sentence sibling-check notes, no derivation (that lives in the
  `answer-check-*.md` files).
- **Facts of record?** README Index is the peer-status artifact; the 0008
  row disagreed and was repaired; no figure was arbitrated.
- **Every home answer diffed, every sibling checked?** Yes — §D5/§D6/§D7
  diffed; five sibling checks written, including 0006 and 0009 whose text
  never moved; each answer's outcome shown here.
- **Deferrals clear all three conditions and are written?** Yes — three
  entries in two `deviations.md` files, each with its check and its
  conditions; rows read `DEFERRED`; RDRs stay Final; nothing fenced or
  meaning-changing deferred; no prior `DEFERRED` row.
- **Both prompts run?** Yes — `critique-set.md` and four `pairwise-*.md`
  under this directory; carried pairs cite `../iter-2/`.
