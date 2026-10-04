# Run scorecard

session: `1ddab04d-5c4c-4e0a-b1cb-ed035eccd301.jsonl`  started: 2026-09-13 16:48 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 104 | 275 | 107 | 2,728 | 23,737,060 | 368,411 | 178,336 | 385,102 |
| Phase 0 spec audit for RDR 0029 | 11 | 31 | 4 | 322 | 502,258 | 87,307 | 17,330 | 88,006 |
| Phase 1 test author for RDR 0029 | 45 | 139 | 19 | 1,410 | 7,416,780 | 258,046 | 92,272 | 268,514 |
| Phase 0 demote re-brief for RDR 0029 | 10 | 26 | 2 | 290 | 395,610 | 50,633 | 9,235 | 56,163 |
| Phase 2 implementer leg 1 for RDR 00 | 49 | 201 | 17 | 1,508 | 8,441,623 | 262,253 | 74,773 | 271,954 |
| Ground D2 open decision | 18 | 53 | 4 | 546 | 1,086,307 | 88,546 | 17,026 | 94,192 |
| Ground D3 open decision | 10 | 28 | 2 | 290 | 293,879 | 37,764 | 7,589 | 43,666 |
| Ground D4 open decision | 10 | 27 | 2 | 290 | 338,767 | 50,492 | 9,332 | 56,022 |
| Strong consult on D2 survivor | 7 | 16 | 1 | 194 | 171,521 | 33,572 | 6,025 | 39,102 |
| Phase 3a CoVe verifier for RDR 0029 | 17 | 59 | 6 | 514 | 1,139,936 | 115,175 | 25,163 | 122,810 |
| Phase 3b adversarial reviewer for RD | 18 | 58 | 5 | 546 | 1,174,261 | 101,686 | 24,198 | 108,560 |
| Phase 3c fixup for RDR 0029 | 30 | 117 | 10 | 930 | 2,685,682 | 135,654 | 46,949 | 142,759 |
| Ground D5 open decision | 12 | 48 | 4 | 354 | 565,223 | 69,784 | 13,904 | 75,314 |
| Strong consult on D5 survivor | 8 | 29 | 2 | 226 | 240,149 | 45,053 | 8,100 | 51,448 |
| Phase 3c leg 2 — ADV predicate corre | 17 | 31 | 5 | 514 | 782,921 | 56,572 | 16,867 | 63,126 |
| Tail triage leaf for rdr-0029 | 24 | 42 | 7 | 528 | 1,974,430 | 228,105 | 30,689 | 134,690 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,512 | 1 | 12,515 |
| ↳ Ground finding F3 seam literals | 10 | 9 | 1 | 20 | 379,452 | 56,465 | 4,031 | 56,467 |
| ↳ Ground finding F4 snapshot reader | 11 | 11 | 1 | 22 | 350,133 | 28,364 | 4,472 | 39,362 |
| ↳ Ground finding F6 ADV-3 header check | 9 | 8 | 1 | 18 | 273,295 | 27,671 | 3,754 | 38,669 |
| ↳ Ground finding F7 UnknownReasons sna | 7 | 7 | 1 | 14 | 230,400 | 33,766 | 4,606 | 44,764 |
| Tail land leaf for rdr-0029 | 10 | 18 | 4 | 230 | 375,070 | 134,315 | 12,154 | 74,638 |
| ↳ Nested spawn probe | 1 | 0 | 0 | 3 | 0 | 12,512 | 174 | 12,515 |
| ↳ Ship rdr-0029 branch | 7 | 10 | 2 | 194 | 173,484 | 27,367 | 5,560 | 33,215 |
| Clear pre-existing modernize lint de | 9 | 20 | 2 | 258 | 234,141 | 29,776 | 5,817 | 35,375 |
| Tail land leaf for rdr-0029 (retry p | 6 | 10 | 1 | 132 | 179,565 | 43,817 | 1,935 | 49,347 |
| ↳ Nested spawn probe | 1 | 0 | 0 | 3 | 0 | 12,505 | 74 | 12,508 |
| **total** | **462** | **1273** | | | | | **620,366** | |

## Totals

- context tokens processed (in + cache read + cache write): **55,562,157**
- output tokens: 620,366
- orchestrator share of context: **43%** at 385,102 peak
- agents: 27 (26 sub-agent legs)
- wall: **107 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **10 min** over 22 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **97 min** over 18 Agent waits
- guard refusals: 2 (harness blocks 2 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate (rdr-gate shard, 16:53): 101,345; at the final turn: 385,102
- parent cd-prefixed calls: 56 (into a worktree: 0)
- parent full-suite runs: 7 beyond the baseline (baseline: yes)

## Budget asks per Phase 2 leg

| leg | budget asks | commits |
|---|--:|--:|
| orchestrator | 2 | 0 |
| Phase 1 test author for RDR 0029 | 2 | 1 |
| Phase 2 implementer leg 1 for RDR 00 | 3 | 2 |
| Phase 3c fixup for RDR 0029 | 1 | 1 |
| Phase 3c leg 2 — ADV predicate corre | 2 | 2 |

## Leg discipline (helper calls against the raw commands they replace)

| leg | leg-read | raw cat | leg-test | raw go test | leg-commit | raw git commit |
|---|--:|--:|--:|--:|--:|--:|
| orchestrator | 1 | 11 | 1 | 12 | 1 | 18 |
| Phase 0 spec audit for RDR 0029 | 0 | 2 | 0 | 0 | 0 | 0 |
| Phase 1 test author for RDR 0029 | 1 | 9 | 2 | 12 | 2 | 0 |
| Phase 2 implementer leg 1 for RDR 00 | 35 | 0 | 11 | 0 | 3 | 0 |
| Ground D2 open decision | 0 | 3 | 0 | 1 | 0 | 0 |
| Ground D4 open decision | 0 | 0 | 0 | 2 | 0 | 0 |
| Strong consult on D2 survivor | 0 | 0 | 0 | 1 | 0 | 0 |
| Phase 3a CoVe verifier for RDR 0029 | 0 | 8 | 0 | 2 | 0 | 0 |
| Phase 3b adversarial reviewer for RD | 21 | 0 | 0 | 4 | 0 | 1 |
| Phase 3c fixup for RDR 0029 | 20 | 2 | 13 | 0 | 1 | 0 |
| Ground D5 open decision | 0 | 2 | 0 | 2 | 0 | 0 |
| Strong consult on D5 survivor | 0 | 1 | 0 | 3 | 0 | 0 |
| Phase 3c leg 2 — ADV predicate corre | 4 | 0 | 2 | 3 | 2 | 0 |
| Tail triage leaf for rdr-0029 | 0 | 0 | 0 | 0 | 0 | 1 |
| Ground finding F3 seam literals | 0 | 0 | 0 | 1 | 0 | 0 |
| Ground finding F4 snapshot reader | 0 | 0 | 0 | 3 | 0 | 0 |
| Ground finding F6 ADV-3 header check | 0 | 0 | 0 | 2 | 0 | 0 |
| Ship rdr-0029 branch | 0 | 0 | 0 | 2 | 0 | 0 |
| Clear pre-existing modernize lint de | 0 | 0 | 0 | 1 | 0 | 1 |

## Human prompts after launch

- 09-13 18:29 clear it in a subagent so we can land the implementation

## Digest

```
scorecard: 55,562,157 context tokens over 27 agents
  orchestrator: 43% of context, 385,102 peak
  turns/tools:  462 / 1273   output: 620,366
  wall:         107 min active, 0 min idle
  test waits:   10 min over 22 runs   waiting on legs: 97 min over 18 Agent waits
  budget asks:  10 asks / 6 commits
  human prompts after launch: 1   agent hand-backs: 0
  orchestrator ctx: gate 101,345 → final 385,102   legs: helpers 120 / raw 110
  guard refusals: 2 (harness 2 · leg-guard 0)   parent: cd-prefixed 56 (worktree 0) · full-suite runs 7 (+1 baseline)
```
