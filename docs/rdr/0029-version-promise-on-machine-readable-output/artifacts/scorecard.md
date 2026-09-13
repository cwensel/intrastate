# Run scorecard

session: `1ddab04d-5c4c-4e0a-b1cb-ed035eccd301.jsonl`  started: 2026-09-13 16:48 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 379 | 275 | 107 | 10,388 | 78,797,767 | 1,411,951 | 704,139 | 385,102 |
| Phase 0 spec audit for RDR 0029 | 42 | 31 | 4 | 1,284 | 1,759,197 | 340,539 | 17,395 | 88,006 |
| Phase 1 test author for RDR 0029 | 184 | 139 | 19 | 5,828 | 27,770,439 | 1,284,380 | 92,754 | 268,514 |
| Phase 0 demote re-brief for RDR 0029 | 36 | 26 | 2 | 1,062 | 1,391,603 | 185,262 | 9,277 | 56,163 |
| Phase 2 implementer leg 1 for RDR 00 | 250 | 201 | 17 | 7,820 | 39,935,670 | 1,514,674 | 75,280 | 271,954 |
| Ground D2 open decision | 71 | 53 | 4 | 2,212 | 4,047,866 | 363,802 | 17,136 | 94,192 |
| Ground D3 open decision | 38 | 28 | 2 | 1,126 | 1,043,814 | 146,593 | 7,645 | 43,666 |
| Ground D4 open decision | 37 | 27 | 2 | 1,094 | 1,207,872 | 194,445 | 9,368 | 56,022 |
| Strong consult on D2 survivor | 23 | 16 | 1 | 646 | 556,857 | 106,808 | 6,067 | 39,102 |
| Phase 3a CoVe verifier for RDR 0029 | 76 | 59 | 6 | 2,342 | 4,976,119 | 532,654 | 25,335 | 122,810 |
| Phase 3b adversarial reviewer for RD | 76 | 58 | 5 | 2,342 | 4,808,437 | 437,781 | 24,315 | 108,560 |
| Phase 3c fixup for RDR 0029 | 147 | 117 | 10 | 4,614 | 12,525,438 | 673,432 | 47,378 | 142,759 |
| Ground D5 open decision | 60 | 48 | 4 | 1,830 | 2,615,652 | 361,039 | 14,031 | 75,314 |
| Strong consult on D5 survivor | 37 | 29 | 2 | 1,094 | 1,070,542 | 196,496 | 8,156 | 51,448 |
| Phase 3c leg 2 — ADV predicate corre | 48 | 31 | 5 | 1,446 | 2,011,530 | 199,168 | 16,949 | 63,126 |
| Tail triage leaf for rdr-0029 | 66 | 42 | 7 | 1,482 | 5,112,748 | 698,006 | 30,754 | 134,690 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,512 | 1 | 12,515 |
| ↳ Ground finding F3 seam literals | 17 | 9 | 1 | 34 | 673,493 | 86,270 | 4,049 | 56,467 |
| ↳ Ground finding F4 snapshot reader | 19 | 11 | 1 | 38 | 569,405 | 67,207 | 4,491 | 39,362 |
| ↳ Ground finding F6 ADV-3 header check | 14 | 8 | 1 | 28 | 434,194 | 38,511 | 3,766 | 38,669 |
| ↳ Ground finding F7 UnknownReasons sna | 14 | 7 | 1 | 28 | 428,985 | 80,047 | 4,623 | 44,764 |
| Tail land leaf for rdr-0029 | 28 | 18 | 4 | 656 | 893,267 | 387,597 | 12,190 | 74,638 |
| ↳ Nested spawn probe | 1 | 0 | 0 | 3 | 0 | 12,512 | 174 | 12,515 |
| ↳ Ship rdr-0029 branch | 17 | 10 | 2 | 454 | 407,235 | 74,044 | 5,584 | 33,215 |
| Clear pre-existing modernize lint de | 29 | 20 | 2 | 838 | 739,257 | 93,689 | 5,852 | 35,375 |
| Tail land leaf for rdr-0029 (retry p | 16 | 10 | 1 | 362 | 460,152 | 111,917 | 1,948 | 49,347 |
| ↳ Nested spawn probe | 1 | 0 | 0 | 3 | 0 | 12,505 | 74 | 12,508 |
| **total** | **1727** | **1273** | | | | | **1,148,731** | |

## Totals

- context tokens processed (in + cache read + cache write): **203,910,437**
- output tokens: 1,148,731
- orchestrator share of context: **39%** at 385,102 peak
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
scorecard: 203,910,437 context tokens over 27 agents
  orchestrator: 39% of context, 385,102 peak
  turns/tools:  1727 / 1273   output: 1,148,731
  wall:         107 min active, 0 min idle
  test waits:   10 min over 22 runs   waiting on legs: 97 min over 18 Agent waits
  budget asks:  10 asks / 6 commits
  human prompts after launch: 1
  orchestrator ctx: gate 101,345 → final 385,102   legs: helpers 120 / raw 110
  guard refusals: 2 (harness 2 · leg-guard 0)   parent: cd-prefixed 56 (worktree 0) · full-suite runs 7 (+1 baseline)
```
