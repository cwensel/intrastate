Model: claude-fable-5
Date: 2026-08-24
Task: B — scoped answer-vs-fences check, sibling 0006 (graph-lint-authority-and-guarantees)
Inputs: home `docs/jdr/0001-resolve-kernel-seam.md` §D5 (222-255), §D7 (300-380), interface record JD-15 (534-540), JD-17 (545-551); RDR 0006 at `526c481` read in full (1740 lines); iter-2 ledger `../iter-2/pairwise-0002-0006.md` G1/G2/G3 + peer-status sweep.

Scope note: 0006's Status qualifier tolerates JD-17 only. JD-15 lists siblings 0002/0004/0009; 0006 is named in §D5's *Lands in* ("**0006** reads a `<clear>` write as removal by citation") and in the JD-15 ledger trail (0002×0006 G3), so it is checked here as well. Read-only: no RDR, JDR, or README edited.

## 1. JD-17 (§D7) — answer vs 0006's fenced text

The answer touches 0006 at: the root/stop-set declarations (A6; invariant 1, 2, 7; the "no initial" fence at 1014-1018; input contract 505-506; disposition row 813; scenario 15/21; Prerequisites 1418-1422); write-replaces (A10; invariant 5 at 584-599; the invariant-class fence at 843-847); the type-model / load-category clauses (821, 1198); accessor references (invariant 1 at 535-536, input contract 507-508).

| # | Clause (0006) | Fenced? | Answer text (§D7) | Verdict |
| --- | --- | --- | --- | --- |
| 17.1 | 1014-1018: "A model that declares no initial owned state MUST be rejected with a blocking finding; lint MUST NOT treat an absent root as an empty reachable set and report a clean model." | fenced | (i) "Root `[initial]` is a table of owned `tag = value` assignments (the root must be a node)" | **CONSISTENT** — the answer supplies the declaration the fence presupposes; meaning unchanged. |
| 17.2 | 843-847: "Graph lint MUST check at least these blocking invariant classes: dangling edge, dead end, determinism/overlap, guard exhaustiveness/gap, single-valued state, owned-set-before-match, and declared terminal/escape handling." | fenced | (i) root/terminal; (iv) write-replaces "a write assigns a tag's whole value and supplants what was held … Closes 0006 A10 and G7" | **CONSISTENT** — invariant 5's per-row reading (592-599: "Dropping the path-accumulation check rests on **write-replaces semantics** … no preceding `clear` required") is exactly the rule (iv) states; the fence's "single-valued state" class is now sound. |
| 17.3 | 1040-1046 (Load-Bearing, unfenced but MUST-bearing): "the root is the declared initial owned state (A6); an edge is a normalized non-escape row … producing the node with that row's writes applied and clears removed" | unfenced | (i) `[initial]` assignments; (iv) write assigns whole value | **CONSISTENT** — "writes applied" under replace semantics is a set-assignment per key, which is what the dataflow already models (per-tag value *set* replaced by the written value on that edge). |
| 17.4 | 584-586 invariant 5: "no row's write block may assign a single-valued tag two values" | unfenced | (iv) "for `set` kind the array literal is the whole new set; no accumulate form … (one value per tag per write)" | **CONSISTENT** — the answer makes the per-row multi-value write a load-level impossibility for scalar kinds; invariant 5 stays as the model-level check over `set`-kind vs single-valued declarations. Home says this closes G7. |
| 17.5 | 542-548 invariant 2: "every reachable owned-state node that satisfies no declared terminal … A terminal is a predicate over owned tags, and a node *satisfies* it when **every** value … meets it" | unfenced | (i) "root `terminal` is a list of **context ids** (the stop set is a predicate, and contexts are already the model's named predicates — no new grammar)" | **CONSISTENT** — a context is a named predicate, so "satisfies a declared terminal" is evaluable unchanged. Note: contexts may carry non-owned atoms, which is why (i) adds the non-owned finding (see §3). |
| 17.6 | 260-268 A6 "Shape requested": "an `initial` table of owned `tag = value` assignments … and a `terminal` list of predicates over owned tags in the same atom shape rules already use. RDR 0002 owns the final spelling; what this RDR requires is that neither is a bare identifier a row references by name" | unfenced (Critical Assumption) | (i) `[initial]` table — identical; `terminal` = context-id list | **CONSISTENT** — `[initial]` matches verbatim. `terminal` differs in *spelling* (context ids, not inline atoms) but 0006 explicitly cedes spelling to 0002 and its underlying requirement (a predicate over tags, not a state name) is met: a context id names a predicate, not a state. Meaning does not change. Flag: the literal phrase "not a bare identifier" now reads awkwardly against a context-id list — prose repair at next touch, not a contradiction. |
| 17.7 | 535-541 invariant 1: "every context reference, tag key, tag value, outcome, and accessor reference named by a row must resolve … the model must declare an initial owned state … Terminal declarations are checked the same way — as tag predicates, not as state names." | unfenced | (i) terminal context ids; (ii) `[tags.<tag>].accessor` removed, `keys` on `[read.<id>]`/`[write.<id>]` is the binding; loader refuses unbound/duplicated keys | **CONSISTENT** — terminal context ids fall under "every context reference … must resolve". The "accessor reference named by a row" check has no site left in the layout (§D7 "Blank": a rule referencing a gate accessor has no site yet); 0006:821 already routes schema non-conformance to 0002 ("refused before normalization; lint never runs") and §D7(iii) says "0006 mints nothing". Net: the accessor-reference arm of invariant 1 becomes vacuous, not wrong. |
| 17.8 | 813 disposition row: "Model with no declared initial owned state … `graph-dangling-edge` naming the missing `[model]` declaration" | unfenced (table) | (i) root `[initial]` is a root-level table, not a `[model]` key | **CONSISTENT** on meaning (the missing declaration is named); the `[model]` locator is stale spelling — citation repair. Scenario 15 (1642-1644) says "naming the absent declaration" and is fine. |
| 17.9 | 246-259 A6 Producer / 253-256: the edit "touches … the closed layout enumeration (root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[dump]`) and '`[model]` MUST contain `id` and `version`'" | unfenced | (ii) `[accessors.<id>]` replaced by `[read.<id>]`/`[write.<id>]`/`[gate.<id>]`; (v) `[model.metadata]` | **CONSISTENT** — 0006 quotes 0002's *pre-answer* enumeration as the thing to be amended; it does not depend on `[accessors.<id>]` surviving. Stale as a description of 0002's current layout — citation repair. |
| 17.10 | 821 / 1198: "Model unreadable / not conforming to RDR 0002's schema \| per RDR 0002 \| refused before normalization; lint never runs \| none of this RDR's codes" | unfenced | (iii) two load categories `malformed tag declaration`, `malformed predicate atom`; "0006 mints nothing for either" | **CONSISTENT** — verbatim agreement in substance. |
| 17.11 | 838-841: "Graph lint MUST consume the normalized candidate-row graph … MUST NOT define a second sparse-source parser" | fenced | whole of §D7 lands in 0002's layout | **CONSISTENT** — 0006 reads every new key through 0002's normalization; silent in a way the answer fills. |
| 17.12 | 862-874: lint "MUST read finite domain, optionality, and single-valuedness from that declaration" | fenced | (iii) keys `kind`/`domain`/`min`/`max`/`elements`/`single_valued`/`required` | **CONSISTENT** — 0006 names the fields abstractly (0003's model, 499-504) and never spells wire keys. |

No fenced clause in 0006 is falsified by §D7. No meaning change anywhere; three unfenced spelling/locator staleness items (17.6, 17.8, 17.9).

## 2. JD-15 (§D5) — answer vs 0006

0006 contains no fenced clause mentioning `clear`/`<clear>`. Touch points are all unfenced.

| # | Clause (0006) | Fenced? | Answer text (§D5) | Verdict |
| --- | --- | --- | --- | --- |
| 15.1 | 494-497 input contract: rows carry "writes, clears" as distinct inputs; 483 "writes and clears" | unfenced | (a) sentinel reserved; 0002 refuses `<clear>` as an authored value at load | **CONSISTENT** — with the sentinel reserved, a `<clear>` entry in 0002's normalized writes re-pairs unambiguously to the row's clear list, which is exactly what iter-2 G3 said was missing. |
| 15.2 | 1044-1046: an edge produces "the node with that row's writes applied and clears removed"; 1114-1115 "A row *preserves* a tag when it neither writes nor clears it" | unfenced (Load-Bearing) | "a `<clear>` write removes the key … clearing an absent key succeeds" | **CONSISTENT** — "clears removed" is the removal reading §D5 names; removing an absent key from a node's per-tag map is a no-op, matching idempotent-clear. This is the "reads by citation" the home describes: 0006 already treats a clear as removal and cites no accessor-level definition. |
| 15.3 | 626-628 invariant 6: "the declared initial owned state counts as a write" | unfenced | (a) unchanged | **CONSISTENT** — silent on clear; the answer merely fills. |
| 15.4 | 592-596 invariant 5: "a write supplants the prior value, no preceding `clear` required" | unfenced | (a) + §D7(iv) | **CONSISTENT**. |

No `<clear>` citation exists in 0006 (it never quotes 0002's `<clear>` clause), so "removal by citation" is currently removal by *description* (15.2) — the home's characterisation holds in substance; nothing to add unless the gate wants an explicit 0004 cite.

## 3. "*Lands in 0006*" items — present / absent / contradicted

| Item (home) | Where it would sit in 0006 | Status | Evidence |
| --- | --- | --- | --- |
| A6 flips by citation (§D7 *Lands in 0006*: "A6/A10 flip by citation") | A6 Status line 232; Prerequisites 1418-1422; Capability row 1191 | **absent** — A6 is still `Status: Pending` (232); 1191 still "Pending (A6)"; 1418 still an unchecked box. 0006's own text pre-authorises the flip: 257-258 "This assumption flips to `Verified` by citation once that clause lands", and 0002's Status line now carries "wire keys for initial/terminal" (0002:9). Not contradicted. |
| A10 flips by citation | A10 Status line 372; Prerequisites 1423-1425; 597 "booked as **A10**" | **absent** — still `Pending`; 379 "Flips when RDR 0002 states write-replaces for single-valued tags" is the flip condition §D7(iv) satisfies. Not contradicted. |
| Terminal-non-owned finding ("A terminal context matching a non-owned tag is a 0006 blocking finding") | invariant 2 (542-548), finding-code table 707-722, disposition table 805-821, advisory-tier fence 1004-1012 | **absent** — no code, disposition row, or scenario names it. Closest text: 544 "A terminal is a predicate over owned tags" (premise, not a check) and invariant 1's context-reference resolution. Not contradicted; the blocking taxonomy fence (843-847) says "at least these", so a new blocking code is admissible without reopening a fence. Needs: a code (or a `graph-dangling-edge` sub-case), a disposition row, one MVV scenario. |
| Root `terminal` as a context-id list | A6 "Shape requested" 263-268; invariant 1 540-541; invariant 2 542-548; scenario 21 (1684-1689) | **absent** as a spelling; semantics present — 0006 says "a `terminal` list of predicates over owned tags in the same atom shape rules already use" and "RDR 0002 owns the final spelling" (265). Not contradicted (a context is a named predicate). Prose "neither is a bare identifier a row references by name" (266-267) wants a repair to "names a context, not a state". |
| `[initial]` assignments (table of owned `tag = value`) | A6 263-264; invariant 5 root paragraph 611-621 ("an `initial` table of owned `tag = value` assignments"); invariant 6 628; fence 1014-1018 | **present** — verbatim shape match (263-264, 618-619). |
| Write-replaces (a write assigns the tag's whole value; `set` array is the whole new set; equality read-back) | invariant 5 592-599; A10 370-399 | **present as the consumed premise** ("write-replaces semantics … a write supplants the prior value"), **absent as a citation** — 594-596 still says "No peer states that rule today", 379-381/390-392 "Flips when RDR 0002 states write-replaces". Not contradicted; becomes a citation to 0002 once repaired. |

## 4. Stale "RDR 0002 is `Draft` / scheduled edit" claims — still present?

All seven iter-2 sites are unchanged at `526c481` (only the Status qualifier moved):

| Line | Text (quoted) | Present? |
| --- | --- | --- |
| 13-15 (Status) | "A6 and A10 are scheduled edits on RDR 0002 (`Draft`), A7 is a route-back on RDR 0003 A18, and A9 flips when RDR 0005 lands transition-model config discovery." | **yes** |
| 256-257 | "RDR 0002 is `Draft` with every assumption Verified, so this reopens no locked document." | **yes** (still false twice — 0002 is Final and its A9–A14 are Pending, 0002:1304) |
| 280-281 | "RDR 0002's authoring schema, which is `Draft` — a scheduled edit on an open peer, not a route-back." | **yes** |
| 380-381 | "RDR 0002 is `Draft`, so this is a scheduled edit on an open peer, not a route-back." (A10) | **yes** |
| 391-392 | "RDR 0002 is `Draft`, so this is a scheduled edit on an open peer, now booked in **Prerequisites**." (A10 Stage 6) | **yes** |
| 1421-1422 | "RDR 0002 is `Draft`, so this schedules an edit on an open peer rather than reopening a locked one." (Prereq A6) | **yes** |
| 1425 | "RDR 0002 is `Draft`, so this too is a scheduled peer edit." (Prereq A10) | **yes** |

Additional stale peer-status text not in the iter-2 seven: 141-142 "a Draft cannot flip a locked peer's record" (0006 calling itself a Draft — now Final; harmless self-reference); 813 "`[model]` declaration" (see 17.8); 253-256 pre-answer layout enumeration (see 17.9). All are citation repairs; none is fenced.

## 5. Status line qualifier (0006:9)

"Final [joint decision → JDR 0001 §JD-17: initial/terminal and write-replaces keys in 0002's layout] [locked 2026-08-23 — Gate PASS. …]"

Names JD-17 only. JD-15 is not tolerated here, consistent with the home's sibling list for JD-15 (0002, 0004, 0009). The trailing bracket still carries the "A6 and A10 are scheduled edits on RDR 0002 (`Draft`)" sentence (13-14), so the qualifier and its own lock note disagree about 0002's status on one line. README:17 row matches the qualifier.

## 6. Summary for the gate

- JD-17 vs fenced 0006 text: 0 CONTRADICTS / 12 CONSISTENT (3 unfenced spelling/locator staleness notes: 17.6, 17.8, 17.9).
- JD-15 vs 0006: 0 CONTRADICTS / 4 CONSISTENT; no fenced touch point exists.
- Lands-in-0006: `[initial]` present; write-replaces present-as-premise/absent-as-citation; A6 flip absent; A10 flip absent; terminal-non-owned finding absent; terminal-as-context-ids absent-as-spelling.
- Seven stale `Draft`/scheduled-edit sites unchanged; plus 141-142, 813, 253-256.
- No single-RDR blocks-impl or risks-impl finding arises from the answers; everything outstanding is a citation repair or an additive landing (one new blocking finding) at 0006's next touch.
