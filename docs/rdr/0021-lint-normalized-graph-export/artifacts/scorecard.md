# Run scorecard

session: `910616ea-7d94-41da-bed7-b6b16b0b820d.jsonl`  started: 2026-09-13 19:09 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 108 | 71 | 78 | 2,796 | 9,846,170 | 431,043 | 143,942 | 149,189 |
| Phase 0 spec auditor | 53 | 42 | 5 | 1,606 | 3,430,970 | 725,389 | 28,414 | 123,778 |
| Phase 1 test author | 169 | 126 | 27 | 5,258 | 36,485,120 | 1,600,498 | 136,821 | 370,425 |
| Phase 2 implementer leg 1 | 226 | 179 | 18 | 7,142 | 37,308,002 | 2,173,397 | 69,110 | 289,488 |
| Phase 3a CoVe verifier | 95 | 75 | 9 | 2,950 | 6,140,229 | 563,334 | 44,593 | 116,530 |
| Phase 3b adversarial reviewer | 106 | 80 | 10 | 3,302 | 10,759,983 | 679,415 | 41,494 | 160,759 |
| Phase 3c fixup | 138 | 107 | 9 | 4,326 | 12,229,320 | 657,178 | 37,102 | 144,665 |
| **total** | **895** | **680** | | | | | **501,476** | |

## Totals

- context tokens processed (in + cache read + cache write): **123,057,428**
- output tokens: 501,476
- orchestrator share of context: **8%** at 149,189 peak
- agents: 7 (6 sub-agent legs)
- wall: **78 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **4 min** over 7 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **78 min** over 6 Agent waits
- guard refusals: 0 (harness blocks 0 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate (rdr-gate shard, 19:46): 107,175; at the final turn: 149,189
- parent cd-prefixed calls: 14 (into a worktree: 0)
- parent full-suite runs: 0 beyond the baseline (baseline: yes)

## Budget asks per Phase 2 leg

| leg | budget asks | commits |
|---|--:|--:|
| Phase 1 test author | 1 | 1 |
| Phase 2 implementer leg 1 | 2 | 1 |
| Phase 3c fixup | 2 | 2 |

## Leg discipline (helper calls against the raw commands they replace)

| leg | leg-read | raw cat | leg-test | raw go test | leg-commit | raw git commit |
|---|--:|--:|--:|--:|--:|--:|
| orchestrator | 0 | 2 | 0 | 1 | 0 | 0 |
| Phase 0 spec auditor | 0 | 6 | 0 | 0 | 0 | 0 |
| Phase 1 test author | 6 | 7 | 2 | 11 | 1 | 0 |
| Phase 2 implementer leg 1 | 85 | 0 | 4 | 12 | 2 | 0 |
| Phase 3a CoVe verifier | 3 | 1 | 0 | 2 | 0 | 0 |
| Phase 3b adversarial reviewer | 17 | 2 | 0 | 17 | 0 | 1 |
| Phase 3c fixup | 24 | 4 | 4 | 0 | 2 | 0 |

## Human prompts after launch

None -- the run was unattended end to end.

## Digest

```
scorecard: 123,057,428 context tokens over 7 agents
  orchestrator: 8% of context, 149,189 peak
  turns/tools:  895 / 680   output: 501,476
  wall:         78 min active, 0 min idle
  test waits:   4 min over 7 runs   waiting on legs: 78 min over 6 Agent waits
  budget asks:  5 asks / 4 commits
  human prompts after launch: 0
  orchestrator ctx: gate 107,175 → final 149,189   legs: helpers 150 / raw 66
  guard refusals: 0 (harness 0 · leg-guard 0)   parent: cd-prefixed 14 (worktree 0) · full-suite runs 0 (+1 baseline)
```
