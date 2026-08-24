# Cluster Reconcile Report — 0002-0009 — Iteration 4

Model: claude-fable-5
Date: 2026-08-24
Iteration: **N = 4** (iter-1 `../report.md` 2026-08-11; iter-2
`../iter-2/reconcile-report.md` 2026-08-23; iter-3
`../iter-3/reconcile-report.md` 2026-08-24). **Past the cap: this iteration
may not demote.** It ran because iter-3's three SPEC-DEFECTs were the only
open ledger entries and the three re-locks landed today; the question it
answers is whether they discharged, and whether anything the re-locks
introduced is a *ledgered* open finding (→ stop) or net-new (→ record, home,
defer — no demotion).

## Membership

Cluster = the Final-and-unimplemented RDRs under `docs/rdr` (0001 Implemented,
excluded). All eight Final at start and at end.

| RDR | Status at start | Revision examined | Moved since iter-3 (`rev` recorded there)? |
| --- | --- | --- | --- |
| 0002 | Final | `365bd66` | **yes** — Stage 3→4→6→7 re-entry, STAGE-SCOPED (iter-3 `526c481`) |
| 0003 | Final [§JD-18] | `31bb9e4` | **yes** — Stage 3→7 re-entry, RE-LOCK-ONLY (`526c481`) |
| 0004 | Final | `ffaf897` | **yes** — Stage 3→7 re-entry, RE-LOCK-ONLY (`526c481`) |
| 0005 | Final [§JD-8/9] | `76a9e67` | no (byte-identical since 2026-08-12) |
| 0006 | Final | `567600f` | iter-3 status edit only |
| 0007 | Final [§JD-8/18] | `526c481` | no |
| 0008 | Final [§JD-5/8/9] | `567600f` | iter-3 status edit only |
| 0009 | Final [§JD-5/8] | `567600f` | iter-3 status edit only |
| **Home — JDR 0001** | open | `567600f` | iter-3's own edits only (sharpenings, sibling-check notes, D7(ii) open note). **No answer** — JD-5, JD-8, JD-9, JD-18 unanswered, so no answer-vs-fences check is owed this pass. |

Delta scope per the iteration contract: re-scan every pair with a moved member
(13), re-run the whole-set critique, run a discharge check per re-locked member
against the home's §D5/§D6/§D7 (the answered-tolerance path those demotions
rode). Eight pairs CARRY.

## Inputs run

| Check | Status | Output |
| --- | --- | --- |
| Whole-set critique | RE-SCANNED | `critique-set.md` (16 of 45 iter-3 rows closed; 9 new) |
| 0002 discharge (STAGE-SCOPED; re-verify A1, A9, A12; D7(ii) note) | DISCHARGE CHECK | `discharge-check-0002.md` |
| 0003 discharge (RE-LOCK-ONLY; `0003:1065-1067` route) | DISCHARGE CHECK | `discharge-check-0003.md` |
| 0004 discharge (RE-LOCK-ONLY; read-back prose, `<clear>`, `0004:371`) | DISCHARGE CHECK | `discharge-check-0004.md` |
| 0002-0003 | RE-SCANNED | `pairwise-0002-0003.md` |
| 0002-0004 | RE-SCANNED | `pairwise-0002-0004.md` |
| 0002-0005 | RE-SCANNED | `pairwise-0002-0005.md` |
| 0002-0006 | RE-SCANNED | `pairwise-0002-0006.md` |
| 0007-0002 | RE-SCANNED | `pairwise-0007-0002.md` |
| 0008-0002 | RE-SCANNED | `pairwise-0008-0002.md` |
| 0009-0002 | RE-SCANNED | `pairwise-0009-0002.md` |
| 0003-0005 | RE-SCANNED | `pairwise-0003-0005.md` |
| 0003-0006 | RE-SCANNED | `pairwise-0003-0006.md` |
| 0007-0003 | RE-SCANNED | `pairwise-0007-0003.md` |
| 0004-0005 | RE-SCANNED | `pairwise-0004-0005.md` |
| 0007-0004 | RE-SCANNED | `pairwise-0007-0004.md` |
| 0009-0004 | RE-SCANNED | `pairwise-0009-0004.md` |
| 0007-0008, 0008-0002 (iter-3 verdicts), 0008-0005-0006, 0008-0009 | CARRIED from iter-3 (no member moved; home unanswered) | `../iter-3/pairwise-*.md` |
| 0005-0006, 0007-0005, 0007-0006, 0007-0009, 0009-0005 | CARRIED from iter-2 via iter-3 | `../iter-2/pairwise-*.md` |

## The set-level finding that governs this report

**The three re-locks discharged every iter-3 SPEC-DEFECT in fenced text**
(`discharge-check-000{2,3,4}.md`: every named clause LANDED or LANDED
DIFFERENTLY without contradicting the home; Refinement Context absent in
all three; Status and README agree). No ledger entry stands open, so the cap
does not force a stop.

What the re-locks exposed is the critique's headline: **the members were
reconciled to the home, not to each other.** Landing §D5/§D6/§D7 in 0002 and
0004 created four fenced member-vs-member seams the home had never
adjudicated — gate placement, owned-state assembly for `flow resolve`,
`Block` cardinality, and the set-value carrier — every one landing on 0005
(never re-entered since June) or on a type 0007 fixed before 0002 widened
it. All four are shared decisions neither RDR solely owns; per the
disposition rules they are hoisted, not demoted. The unmoved member most
likely rewritten is **0005**.

## Origin ledger trace (iter-3 table → this pass)

| Iter-3 row | Now |
| --- | --- |
| 0002 × home SPEC-DEFECT (JD-16 routing; JD-17 layout; JD-15 residue) | **RESOLVED** — routing block-keyed `0002:1168-1192`; three-valued block `923-938`; `in`-expansion + suffix `1103-1114`; closed layout `606-621`; capability tables `661-676`; `[initial]`/`terminal` `716-748`; type-model keys `1384-1387`; `[model.metadata]`; write-replaces `838-841`; two categories `1392-1396`; `<clear>` reserved `1049-1061`; A1 Verified (evidence `929e931`). A9/A12 **Pending, gate-accepted at Stage 6** (blocked on 0007's reshape) — see Observations. |
| 0003 × home SPEC-DEFECT (`1065-1067` route onto 0006 findings) | **RESOLVED** — `0003:915-925` routes both rejections to 0002's load categories, "RDR 0006 mints nothing"; §JD-16 cleared; §JD-18 carried. One unfenced MVV residue (`0003:2094`) → 0003 `deviations.md` D2. |
| 0004 × home SPEC-DEFECT (read-back presence; zero `<clear>`; `0004:371`) | **RESOLVED** — absence arm `0004:235-238`, `571`; fenced `<clear>` clause `357-364`, A11, Scenario 9; "cannot be rebound" retired `408-411`. Prerequisites still unchecked (Observations). |
| Home JD-17 "Open at (ii)" (critique N-4; owner 0002) | **RESOLVED at 0002** — provenance-scoped binding validation fenced at `0002:675-707`; home §D7(ii) and JD-17 note now record it (this pass). 0004's unscoped restatement → 0004 D2 citation repair. |
| 0006 DEFERRED D1 (§D7 landings) | CARRIED — check 1's precondition now met on 0002's side (`pairwise-0002-0006.md`); check 2 stays 0006's. Not re-deferred. |
| 0006 DEFERRED D2 (stale "0002 Draft" ×7) | CARRIED — still present; Status line too (`0006:9` "RDR 0002 (`Draft`)"). |
| 0008 DEFERRED D1 (scenario-3 view seam) | CARRIED (0007-0008 pair unmoved). |
| 0008 DEFERRED D2 (A10 on retired layout) | CARRIED — **its grep was run this pass and PASSES** (`0002:678-680` fenced; `negative-cases.txt:48` refuses `recognized` in `keys`); the text residue closes at 0008's next touch. |
| 0008-0009 → JD-5; 0005/0007/0008/0009 → JD-8; 0005/0008 → JD-9; 0003/0007 → JD-18 | STILL-OPEN tolerances (home unanswered); JD-8 and JD-9 **sharpened** below. |
| 0008-0006 dead-rule fixture (routed into 0002's re-entry) | RESOLVED on 0002's side (`0002:1076-1078` names the authorable form); 0006's fixture/quote at `0006:670-682` stale → 0006 D2 scope (cosmetic). |
| 0007 / 0009 sentence-entangled peer-status claims | STILL-OPEN, cosmetic; 0007's now stale in the *opposite* direction (all four duties it routes to "0002's Refinement Context" landed fenced in 0002: `950-957`, `959-973`, `1125-1142`, fixture rename). |
| 0002 restates §D2; 0003/0006 `graph-product-too-large` | CARRIED, cosmetic. |

No prior `DEFERRED` row was re-found by any scan, so no double-deferral is
possible; iter-3's three deferrals carry unexamined per the contract except
0008 D2, whose check ran incidentally and passed.

## Findings table

Ledger tag `NET-NEW` = traces to no prior report row; per the contract such a
finding demotes nobody and blocks no implement on its own.

| Pair | Finding | TYPE | SEV | Own | Ledger | Disposition |
| --- | --- | --- | --- | --- | --- | --- |
| 0002-0005 | `flow resolve` fenced with no read accessors (`0005:380-383`) + JD-9 "wire `--tag` to Observed" + 0002's every ordinary row write-bearing (`0002:830-833`, `1125-1130`) ⇒ every ordinary row refuses `owned_state_unavailable` | contradiction | blocks | joint | §JD-9 (sharpens iter-2 0007-0005 row 4); critique Q-3 | **JOINT-DECISION → §JD-9 sharpened** (which verb assembles owned state and calls `Resolve`); 0002 gains the §JD-9 qualifier |
| 0002-0004 / 0004-0005 | gate `deny` meaning and timing delegated by 0002 (`709-714`) to 0004, which is silent and lists `gate denied` as both refusal (`804`) and typed result (`502`); 0005 `resolve` has no gate carrier | gap | risks | joint | NET-NEW (critique Q-2) | **JOINT-DECISION → new §JD-19**; siblings 0002, 0004, 0005 |
| 0004-0005 / 0002-0005 | `--write name=value` cannot carry `<clear>` or a set-kind sequence; no boundary refuses a `--write` key outside a writer's `keys`; §D5 landing list omits 0005 | gap / round-trip | risks | joint | NET-NEW (0004-0005 F7/F8, 0002-0005 F7) | **JOINT-DECISION → new §JD-20**; siblings 0002, 0004, 0005 |
| 0002 × 0007 | `Block` "exactly two constants" fenced in 0007 (`1264-1266`) vs third `match` member fenced in 0002 (`923-938`) | contradiction | risks | joint | NET-NEW (critique Q-4; 0007-0002 scan missed it, verified by the parent) | **JOINT-DECISION → new §JD-21**; siblings 0002, 0003, 0007 |
| 0007-0003 / 0002-0004 | set-valued tag *value* encoding across `Tag.Value string` undeclared; 0007's `contains` test leg blocked on 0003 (`0007:2157-2168`) | gap | risks | joint | NET-NEW (0007-0003 F5, 0002-0004 F-A, critique Q-5/C-10) | **JOINT-DECISION → new §JD-22**; siblings 0002, 0003, 0004, 0007 |
| 0004-0005 | `read_back_incomplete` collapses into `flow-write-readback-mismatch` (0004 fenced MUST NOT `376-390`) + 7 refusal classes unmapped + exit-3 arm | contradiction / gap | risks | joint | iter-2 F1-F3 → §JD-8 | STILL-OPEN tolerance; **§JD-8 note gains the ordering root** (0005's re-entry is the next move) |
| 0002-0005, 0003-0005 | 25 load categories / predicate kinds with no envelope code; escape plan rendered as success; no path form; `flow-guard-unevaluable` "supplied facts" | gap | risks | joint | iter-2 → §JD-8 / §JD-9 | STILL-OPEN tolerances (0005 byte-identical) |
| 0002-0003 | 0002's normative fixture `kind = "string"` (`rdr-fixture.toml:65`, `0002:1698`, `340`, `1551`) vs 0003's five fenced kind tokens (`0003:791-793`); critique Q-1 `single_valued` on `iter`/`cluster_ready` vs `0003:948-952` | contradiction | risks | single (0002 artifacts) | NET-NEW | **DEFERRED** → 0002 `deviations.md` D1 — (a) fixture/evidence/illustrative, 0002's fence defers the vocabulary to 0003; (b) grep + lint; (c) fences agree |
| 0009-0002 | `writes = []` presence-keyed in both fences; spike checks length (`main.go:546`); no empty-block fixture | gap | risks | single (0002 evidence) | JD-11 row 3 residual | **DEFERRED** → 0002 D2 — (a) evidence only; (b) mint `neg-escape-with-empty-write`; (c) no meaning change |
| 0002-0006 | write-value kind/domain conformance claimed "at minimum", no load category; 0006 checks it in invariants 1/5 | gap | risks | single (0002) | NET-NEW | **DEFERRED** → 0002 D3 — (a) unfenced both sides; (b) one Scenario-3 fixture; (c) additive |
| 0009-0004 | 0009's write-accessor obligation + escaped-plan/`NextTags` test (`0009:1515-1527`) bind nobody in 0004; `0004:353-354` forbids the regression | gap | risks | joint (JD-7 closed, residue) | iter-2 F2 / iter-1 F1 | **DEFERRED** → 0004 D1 — (a) unfenced; (b) named test; (c) additive. Home JD-7 notes the residue |
| 0007-0003 | empty/omitted `unless` identity fenced in 0007 (`1395-1409`); 0003's fence silent (`1143-1170` scoped to decided atoms) | gap | risks (critique C-16) / cosmetic (scan) | single (0003) | iter-2 F3 / C-16 | **DEFERRED** → 0003 D1 — (a) 0003 silent, prose unfenced; (b) 0003 Scenario 2 on an `unless`-less group; (c) fills silence, changes no clause |
| 0002 | fenced peer-status sentences: `0002:1053-1058` "0004 does not yet carry [§D5]"; `0002:1229-1233` "0003 is `Draft`"; `2073-2077` "reshape has no owner"; `2297-2298` suffix wording | contradiction | cosmetic | single (0002) | NET-NEW (critique Q-7) | **CITATION REPAIR** (artifact of record: README) — written as 0002 D4 for its next touch; no design-body edit of a Final RDR |
| 0003 | `0003:2094` "0006 can map to a lint finding"; A20 "0007 states no atom carrier"; `1530` "0002 Pending"; 0005 as lint enveloper | contradiction | cosmetic | single (0003) | iter-3 answer-check residue / iter-2 F2/F5 | **CITATION REPAIR** → 0003 D2 |
| 0004 | guard-absent → `no_match` prose (`495`, `523-527`, `828-831`; regressed 1→3 sites); A6b "Pending" cite (`1015`); §D7(ii) unscoped restatement (`254-257`); vacuous 0008 cite | contradiction / duplication | cosmetic | single (0004) | iter-2 F1/F2, iter-3 N-4, critique Q-8/Q-9/C-21 | **CITATION REPAIR** → 0004 D2 |
| 0007 | `0007:687, 939, 951, 2040` cite 0002's cut Refinement Context (duties landed fenced); `148/463/600/613/645/2172` call 0002/0003/0006 Draft; `904-907` "§JD-3 NOT closed"; A12 claims false vs `0002:1384-1387`, `545-547`; A9 quote drifted; pre-§D7 fixture path | contradiction | cosmetic | single (0007) | iter-2/iter-3 residue + NET-NEW | CITATION REPAIR — pend 0007's next touch (sentence-entangled; unchanged policy) |
| 0006 | `0006:9` Status, `141-148`, `307-311`, `592-597`, `670-682`, `946-947` stale after 0002's and 0003's re-locks (A7 flip condition met, A21 citation) | contradiction | cosmetic | single (0006) | iter-3 D2 scope | CARRIED under 0006 D2 (citation repair) |
| 0008 | A4/A9 evidence cites superseded spike; `keys` position omitted; block-2 parenthetical; `2359-2360` qualifier list | contradiction | cosmetic | single (0008) | NET-NEW / iter-3 N1-N5 | recorded; 0008's next touch |
| 0009 | `Refusal.Guard` `421`; "implemented" `1184`; retired "row kind `escape`" ×3 | contradiction | cosmetic | single (0009) | iter-2/3 | CARRIED; 0009's next touch |
| 0003-0006 | `MUST NOT suppress` vs `MUST NOT additionally emit graph-coverage-gap` — qualifier "over a provable product" reconciles; shared Scenario 8 | duplication | cosmetic | joint | NET-NEW | NO CONFLICT |
| 0009-0002, 0002-0006, 0007-0004, 0009-0004 | round-trips (JD-11 nine fields; initial/terminal/writes; accessor→kernel presence facts; `<clear>` write→absent) | round-trip | — | — | iter-1/2/3 rows | RESOLVED / NO CONFLICT — asserted value-for-value in each scan |
| README | all eight rows agree with Status lines at pass start; five rows gain qualifiers this pass | — | — | — | — | CITATION REPAIR — done |

## Dispositions

### No SPEC-DEFECT

No finding is a ledgered open entry, and none is a single-RDR contradiction
in fenced text or in meaning. The cap is therefore not tested: nothing
qualified for demotion, and net-new findings could not have driven one.

### JOINT-DECISIONS — four new homes (JD-19…JD-22), two sharpened (JD-8, JD-9)

Each is genuinely joint: the seam is between a re-locked member's landing and
a peer that is *silent or fixed earlier*, and neither RDR can decide it alone
(gate site spans 0002's layout, 0004's accessor model, 0005's verb; CLI
carriage spans §D5/§D7 and 0005's grammar; `Block` arity is a type two fences
name; set-value encoding is a kernel-seam byte form four RDRs consume). Each
home entry is one paragraph — decision, stake, siblings, evidence pointer —
with derivation in this tree. Qualifiers written (Status + README):

| RDR | Qualifier after this pass |
| --- | --- |
| 0002 | `Final [joint decision → JDR 0001 §JD-9, §JD-19, §JD-20, §JD-21, §JD-22]` |
| 0003 | `Final [joint decision → JDR 0001 §JD-18, §JD-22]` |
| 0004 | `Final [joint decision → JDR 0001 §JD-19, §JD-20, §JD-22]` |
| 0005 | `Final [joint decision → JDR 0001 §JD-8, §JD-9, §JD-19, §JD-20]` |
| 0007 | `Final [joint decision → JDR 0001 §JD-8, §JD-18, §JD-21, §JD-22]` |
| 0006, 0008, 0009 | unchanged |

### DEFER-TO-IMPLEMENTATION — six entries, three RDRs — written

| RDR | Entry | Mechanical check | Conditions |
| --- | --- | --- | --- |
| 0002 | D1 fixtures vs 0003's kind tokens / assignment table | grep `kind = "string"` → 0; 0006 lint (or 0003 table) over promoted fixtures | (a) artifacts, (b) grep+lint, (c) fences agree |
| 0002 | D2 empty write block on escape row | mint `neg-escape-with-empty-write` → `malformed escape declaration` | (a) evidence, (b) fixture, (c) additive |
| 0002 | D3 write-value kind/domain category | one Scenario-3 out-of-domain fixture → named category | (a) unfenced, (b) fixture, (c) additive |
| 0003 | D1 empty `unless` identity | Scenario 2 on an `unless`-less group = full product | (a) silence/unfenced, (b) scenario, (c) no clause changes |
| 0004 | D1 0009's write-accessor obligation/test | author escaped-plan/`NextTags` test; `grep NextTags` ≥1 | (a) unfenced, (b) test, (c) additive |
| 0002 D4, 0003 D2, 0004 D2 | citation repairs (not deferrals) | named greps | facts of record |

Nothing fenced and nothing meaning-changing was deferred; no prior `DEFERRED`
row was deferred a second time (0006 D1/D2, 0008 D1/D2 carry unexamined,
except 0008 D2's check which ran and passed).

### HOME edits (JDR 0001) — done

- §D7(ii): provenance-scoped binding validation recorded (settled by 0002).
- JD-7: residue note → 0004 D1. JD-8: ordering root (0005's re-entry is the
  next move). JD-9: sharpened — the open question is owned-state assembly
  for `flow resolve`; 0002 added as sibling. JD-17: "Open at (ii)" closed;
  iteration-4 sibling re-locks recorded.
- JD-19, JD-20, JD-21, JD-22 appended (one paragraph each).

### CITATION REPAIRS — README rows updated (done); in-text repairs recorded

0002 D4, 0003 D2, 0004 D2, 0006 D2 (carried); 0007/0008/0009 residues pend
their next touch. No fenced text of a Final RDR was edited.

## Verdict

**RECONCILED WITH TOLERANCES.**

- No SPEC-DEFECT stands open; iter-3's three are discharged.
- Standing tolerances:

| Tolerance | Home | Open question | Answered? |
| --- | --- | --- | --- |
| 0008, 0009 | §JD-5 | precedence of the two `Resolve`-entry preconditions; error wrapping | No |
| 0005, 0007, 0008, 0009 | §JD-8 | full code table / exit-3 / envelope carriers; **ordering root — 0005 re-entry** | No |
| 0002, 0005, 0008 | §JD-9 | which verb assembles owned state and calls `Resolve`; `--tag recognized=` classification | No (sharpened) |
| 0003, 0007 | §JD-18 | conforming-view enforcer | No |
| 0002, 0004, 0005 | §JD-19 | gate evaluation site and `deny` semantics | No (new) |
| 0002, 0004, 0005 | §JD-20 | CLI carriage of `<clear>` / set writes; unbound `--write` key | No (new) |
| 0002, 0003, 0007 | §JD-21 | `Block` cardinality | No (new) |
| 0002, 0003, 0004, 0007 | §JD-22 | set-value encoding at `Tag.Value` | No (new) |

- Deferred: 0002 D1–D3, 0003 D1, 0004 D1 (this pass); 0006 D1, 0008 D1–D2
  (carried). Citation repairs: 0002 D4, 0003 D2, 0004 D2, 0006 D2.
- The cluster may implement. **Practical ordering**: the eight tolerances
  concentrate on 0005 (four of eight) and on 0002's seams with 0007
  (JD-21/22). A sibling with an ANSWERED tolerance runs its scoped
  answer-vs-fences check before it implements. JD-8/JD-9/JD-19/JD-20 all
  re-walk 0005 — answering them at the home and then re-entering 0005 once
  is cheaper than four scoped checks; JD-21/JD-22 are type-level and are
  cheapest to answer with the first Go type (0007 Phase 1 / 0002's
  normalizer), i.e. at Stage 8.

## Observations outside this gate's authority

- **0002 A9 and A12 are `Pending`, not `Verified`**, despite iter-3's
  `re-verify A1, A9, A12` — 0002's Stage 6 downgraded them (blocked on
  0007's atom reshape, recorded unowned at `0002:2073-2081`) and the
  finalize gate accepted that. Recorded as a sequencing fact: 0007 Phase 1
  is the reshape owner (`0007:2117-2126`); 0002's D4 repairs the "no owner"
  sentence.
- **0004 is Final with all four Prerequisites unchecked** for the third
  lock (`0004:841-849`); its Status discloses A9/A10/A11 Pending. A per-RDR
  finalize question, noted in 0004 `deviations.md` for Stage 8's opener.
- **0005 has not been touched since 2026-08-12** and now holds four
  tolerances; the critique names it the most likely rewrite.
- **Code has not moved** (critique C-31): ~75.7k markdown lines vs 3,971 Go;
  no `internal/` commit since 2026-08-09.

## Review Gate

- **Cluster correct?** Yes — same eight members as iter-1..3; 0001 excluded;
  every member cites or is cited by another; the home lists all eight.
- **N>1 anchored to history?** Every check is a RE-SCANNED, DISCHARGE CHECK,
  or CARRIED row naming its source iteration; re-scans are exactly the 13
  pairs whose member moved plus the three discharge checks; no demotion
  occurred, so none is unledgered.
- **Cap held?** N=4 demoted nobody; no ledger entry stood open, so the
  `stopped:cluster-flapping` packet was not owed. NET-NEW findings were
  recorded, homed, or deferred — none blocked implement on its own.
- **No design-body edit of a Final RDR?** Correct — edits: five Status
  lines, five README rows, the home, three new `deviations.md` files.
- **Each JOINT-DECISION genuinely joint?** Yes — each names why neither RDR
  owns it (above); each home entry is one paragraph; each sibling's
  qualifier names home and question. 0007-0002's scan missed JD-21; the
  parent verified both fences before homing it.
- **Facts of record?** README Index is the peer-status artifact; every
  disagreeing peer figure is typed CITATION REPAIR (0002 D4, 0003 D2,
  0004 D2, 0006 D2), not deferral or demotion.
- **Every home answer diffed, every sibling checked?** The home's only
  movement since `4f9639c` was iter-3's own edits — diffed (`git diff
  4f9639c 567600f`), no answer found; the three siblings that re-entered on
  iter-3's answers each have a discharge check here.
- **Deferrals clear all three conditions and are written?** Yes — six
  entries in three files, each with its check and its carried conditions;
  rows read `DEFERRED`; RDRs stay Final; no prior DEFERRED row re-deferred.
- **Both prompts run?** Yes — `critique-set.md` and thirteen `pairwise-*.md`
  under this directory; carried pairs cite `../iter-3/` or `../iter-2/`.

## Next

`Next: /rdr-implement NNNN` for each member — subject to the ordering note
in the Verdict. Cheapest first move for the set: answer JD-8/9/19/20 at the
home, then `/rdr-refine 0005` → `/rdr-finalize 0005` (one re-walk clears
four tolerances); JD-21/22 answer with the first Go types.
