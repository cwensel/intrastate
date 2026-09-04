# Cluster reconcile — 0026-0027-0028 — iteration 1 (2026-09-03)

Model: claude-fable-5-1

## Members and homes, revisions as examined

| Record | Revision | Role |
| --- | --- | --- |
| 0026 | fe79fcc | member (seed) |
| 0027 | 10f19d1 | member (seed; `candidate` from 0026 alone, confirmed by `--cluster-of 26,27`) |
| 0028 | 41612e8 | member (peer-evidence via 0026; cross-cutting via 0028) |
| JDR 0003 | 48f5163 | home — §D1 (0026×0028), §D2–§D4 filed this pass; `revisions.txt` carries the post-filing revision as the next iteration's baseline |
| 0025 | 22805c8 | home — C5 (0027×0028 clause map), C4 (deadline triple); Implemented |

Cycles: none (`index --cycles`). Open joint decisions at entry: none (`index --open-joint`). Cap row: `cap-first` → run. Lint at entry: all three PASS (advisory only).

## Inputs run (scanned 2 of 3 possible pairs)

| Check | Status | Note |
| --- | --- | --- |
| Whole-set critique | RE-SCANNED | `critique-set.md` — 12 ledger rows (4 H, 6 M, 2 L) |
| 0026 × 0028 | RE-SCANNED | `pairwise-0026-0028.md` — 5 findings; round-trip HOLDS; both conform to JDR 0003 §D1 |
| 0027 × 0028 | RE-SCANNED | `pairwise-0027-0028.md` — 4 findings; clause map consistent with 0025:C5; landing order either |
| 0026 × 0027 | TRIMMED | no edge of any kind (`--backlinks`), no shared anchor (`--anchor-intersect`); each JC1 disposed the other disjoint (0026:C1 amends 0025:C4, 0027:C1 amends 0025:C5); the whole-set critique covered the triple |

## Findings and dispositions

| # | Pair | Source | Type | Severity | Disposition | Detail |
| --- | --- | --- | --- | --- | --- | --- |
| R1 | 0026×0028 | PW1, C-6 | gap | cosmetic | CITATION REPAIR | JDR 0003 cited `0028:C3` (4 sites); the passage is `0028:C1.3`. Repaired at the home. `0026:JC1`'s `cli/0028:C3` (Decision Rationale narration, locked record) left as is under the consumer's no-amend rule; meaning intact. |
| R2 | 0026×0028 | PW2 | gap | cosmetic | NO CONFLICT | 0026:C1's "0028's read-back leg should revisit" hand-off has no recipient — 0028 adds no `spawn` caller. Unowned, not contradicted. |
| R3 | 0026×0028 | PW3, C-9 | gap | risks-impl | DEFERRED | 0026 `artifacts/deviations.md` D1 (TEST-FIXTURE): negative witness that the re-keyed `execution_failure` rendering keys on `Applied()`, not phase; `errors.As` placement in `Write`'s err arm. Conditions: unfenced (S1/S27 testing strategy); a unit test decides it; no clause meaning changes. |
| R4 | 0026×0028 | PW4, C-5 | gap | cosmetic | NO CONFLICT + unowned demand | Both records accept that `timeout` swallows the sub-reason (0025:C4 `detail:`); 0026:F6 asks 7.1 to name an owner. No member owns it; the route-back F6 names (Stage 6 on 0025) is unavailable — 0025 is Implemented and code is the source of truth. Vehicle: a new RDR seed against 0025:C4 `detail:` on `timeout`. Author's call; not blocking. |
| R5 | 0026×0028 | PW5 | gap | cosmetic | NO CONFLICT | 0026 MVV row 5 witnesses `Err` fields; S1 carries the `Detail` assertion D1 rule 1 moved the reason to. |
| R6 | 0026×0028 | C-1, C-6 | gap | H (critique) | JOINT-DECISION → JDR 0003 §D4 (open) | Applied sense unreadable at the envelope: no `json:"applied"`, no non-test `Applied()` caller, both producers render `flow-accessor-failed`, discriminated by `detail` prose. Not a member contradiction (both conform to D1); a fork at the seam's caller edge. D1's option-(a) premise describes pre-0026 shipped code and rule 3 amends it — not a self-contradiction. |
| R7 | 0026×0028 | C-7 | gap | M | NO CONFLICT (ordering) | `Executor.Write` touched by both: 0026's `errors.As` in the `err != nil` arm; 0028's gate pre-check before `Apply`. Disjoint arms; R3's check pins the placement. |
| R8 | 0026×0028 | C-10 | gap | M | NO CONFLICT (ordering) | 0028's command read-back relies on 0026's timeout delivery and R3's witness needs 0026's re-key landed: implement 0026 before 0028. |
| R9 | 0027×0028 | PW1, C-11 | gap | risks-impl | DEFERRED | Shipped `TestReq77` asserts 0025:C5's six as `Categories()`'s tail (verified `command_carrier_0025_test.go:853-862`); 0028:C1.4's five appends break it; 0027:S6 must be relative-order. 0028 deviations D1, 0027 deviations D1 (TEST-FIXTURE). Conditions: shipped test + unfenced S6; `go test` decides; no meaning change. |
| R10 | 0027×0028 | PW2 | gap | risks-impl | DEFERRED | 0027 Phase 3's description surface — total over `Categories()` or opt-in — decides whether 0028's five need text. 0027 deviations D2 (IMPL-DECISION), 0028 deviations D2 (SPEC-UNDER). Conditions: unfenced (Phase 3, S7); totality test / docs gate decides; no meaning change. |
| R11 | 0027×0028 | PW3 | duplication | cosmetic | CITATION REPAIR (typed; not applied) | 0028:C1.4's aside "clauses 2–6, which are `command`-only" misdescribes 0025:C5's map — source runs clauses 5–6 for every entry. The normative ordering claim ("after clauses 1–6") is correct and unaffected. Repair would edit fenced text of a locked record; withheld under the consumer's no-amend rule. Authority: 0025:C5 and `load.go::carrierDefect`. |
| R12 | 0027×0028 | PW4, C-8 | gap | cosmetic | NO CONFLICT (ordering) + DEFERRED check | Clause 3 of `carrierDefect` is textually shared (family check + signature change vs exemption comment). Second lander rebases; merged-behaviour fixture pre-seeded as 0027 D3 / 0028 D3. 0027:A1's line anchors go stale on a still-resolving symbol — a non-finding. |
| R13 | 0027×0028 | C-3, C-4, C-12 | contradiction | H (critique) | JOINT-DECISION → JDR 0003 §D2 (open) | 0028:C1.6 admits an invocation-bound argv word at any position with no value rule; 0027:C1's promise names two admitted forms and "argv WORDS only". Verified against `0027:C1`, `0028:C1.6`, `load.go` clause 3, `cmdbind.go::substitute`. Neither record owns the answer (promise wording vs placeholder admission). |
| R14 | 0028 (vs JDR 0001 §D10) | C-2 | gap | H (critique) | JOINT-DECISION → JDR 0003 §D3 (open) | Stale-model refusals ride `execution_failure` → exit 3 ("re-run unchanged"); the CLI contract says a request refusal must never exit 3; §D10 rule 7 is the exit-2 precedent. JDR 0003 §D1 had deferred this "to 0028 by citation"; 0028 carries no §D10 citation. |

Critique rows not listed above: C-5 = R4; C-6 = R1/R6; C-8 = R12; C-9 = R3; C-11 = R9. No finding was a SPEC-DEFECT: every fenced-text tension is a shared decision with a named home, and every single-RDR item is unfenced and mechanically decidable.

## Joint decisions filed (home: JDR 0003; entry state `open`)

Per `docs/jdr/README.md`, an `open` entry blocks the siblings it spans and the
`Final [joint decision → …]` qualifier is stamped when the entry is decided,
not filed. So no Status line was touched this pass, and:

| Entry | Question | Spans | Recommendation recorded |
| --- | --- | --- | --- |
| §D2 | Is an argv word bound at invocation inside the inline-shell promise? | 0027, 0028 | (a) a placeholder rule in 0028: never argv0, `-`-prefixed value refused |
| §D3 | Which exit group do 0028's stale-model refusals take? | 0028, JDR 0001 §D10 | (b) exit 2 via typed `Err` + a distinct CLI code, no new class |
| §D4 | How does the applied sense reach the CLI envelope? | 0026, 0028, JDR 0001 §D10 | (b) a distinct exit-3 code keyed on `Applied()` |

When an entry is decided, each spanned record owes its scoped answer-vs-fences
check before it implements: D2 (a) amends `0028:C1.6` (a contract-scoped
re-lock of 0028); D2 (b) amends `0027:C1`; D3 (b) and D4 (b) touch no fence
(both land as Stage 8 surface plus a §D10 table citation repair, though the
new CLI code is a SPEC-UNDER author decision at Stage 8 by the launch
prompt's additive-surface rule).

## Deferred (ledger closed; ids in `deferred.txt`)

R3 → 0026 deviations D1. R9 → 0027 D1, 0028 D1. R10 → 0027 D2, 0028 D2. R12 → 0027 D3, 0028 D3.

## Verdict

**RECONCILED** — no SPEC-DEFECT stands; no member demoted. Three JOINT-DECISION
entries are filed `open` at JDR 0003 (§D2–§D4) and, under the consumer's JDR
doctrine, gate implementation of the records they span until decided: 0026
(§D4), 0027 (§D2), 0028 (§D2, §D3, §D4).

Implementation-ordering consequences once the entries are decided:
1. 0026 first — 0028's command read-back relies on its timeout delivery (R8) and R3's witness needs its re-key.
2. 0027 and 0028 in either order — second lander rebases `carrierDefect` clause 3 and the doc comment (R12); 0028's `TestReq77` rewrite (R9) lands with whichever appends first.
3. If D2 resolves (a), re-lock 0028 before implementing it.

Open obligation with no cluster owner: 0025:C4 `detail:` on `timeout` (R4) — a seed, at the author's discretion.

## Addendum — homes answered 2026-09-03 (same pass)

JDR 0003 §D2–§D4 decided the same day after a corpus/precedent pass
(`research-d2-d4.md`): §D2 (a), §D3 (b), §D4 (b). Scoped answer-vs-fences
checks, per sibling:

| Sibling | Entry | Check | Outcome |
| --- | --- | --- | --- |
| 0027 | §D2 (a) | `0027:C1` promise and out-of-scope list unchanged by a 0028-side rule | consistent — stays Final, no qualifier |
| 0028 | §D2 (a) | `0028:C1.6` "under 0025:C2's rule unchanged", "no other change to 0025:C1–C1.6" — the rule contradicts fenced text | SPEC-DEFECT → 0028 Draft, re-entry Stage 3, STAGE-SCOPED, re-verify A7, A9; README row flipped; re-entry note appended; punt-ledger row in `0028-…-postmortem.md` |
| 0028 | §D3 (b) | `0028:C1.3` "no new refusal class", class `execution_failure` — (b) adds no class; the CLI code and typed `Err` are outside the fence | consistent; the re-entry note directs a citation of §D3 at C1.3 |
| 0026 | §D4 (b) | `0026:C1` names no CLI code; the gain lives in A8/S1 | consistent — stays Final, no qualifier; `artifacts/deviations.md` D2 records the decided form |
| 0028 | §D4 (b) | 0028's refusals keep `flow-accessor-failed` | consistent |

Verdict after answers: **NOT RECONCILED** for 0028 only — it re-locks at its
named scope and rejoins; 0026 and 0027 may implement. Order: 0026 first (R8),
then 0027; 0028 after re-lock (a re-run of this gate scoped to the 0027×0028
and 0026×0028 pairs is due if the refine touches C1.3/C1.4 beyond the
citation).
