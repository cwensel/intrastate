# Run scorecard

session: `09736056-427e-41a5-a0e2-95e278560827.jsonl`  started: 2026-09-23 19:26 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 49 | 42 | 92 | 98 | 4,624,985 | 107,235 | 21,242 | 122,782 |
| Phase 0 spec audit 0012 | 43 | 49 | 7 | 86 | 4,522,612 | 184,363 | 1,809 | 190,137 |
| Phase 1 tests-first 0012 | 80 | 81 | 23 | 160 | 11,020,376 | 254,630 | 2,171 | 260,431 |
| Phase 2 leg 1 implement 0012 | 87 | 88 | 13 | 174 | 15,340,965 | 257,667 | 2,722 | 263,468 |
| Phase 3a CoVe verify 0012 | 60 | 59 | 10 | 120 | 5,506,069 | 140,612 | 1,678 | 146,914 |
| Phase 3b adversarial review 0012 | 51 | 50 | 9 | 102 | 4,517,919 | 137,900 | 1,506 | 143,687 |
| Phase 3c fixup 0012 | 47 | 46 | 10 | 94 | 3,607,236 | 104,395 | 1,388 | 110,169 |
| Tail triage leaf rdr-0012 | 34 | 31 | 10 | 76 | 1,631,461 | 115,700 | 1,552 | 67,120 |
| ↳ Nested spawn probe | 3 | 1 | 0 | 9 | 30,910 | 13,983 | 369 | 17,184 |
| ↳ Ground roborev job 8229 | 17 | 16 | 2 | 34 | 603,758 | 40,035 | 723 | 46,152 |
| ↳ Ground roborev job 8247 | 12 | 11 | 2 | 24 | 450,814 | 44,845 | 932 | 50,637 |
| ↳ Phase 3c fixup rdr-0012 | 24 | 24 | 4 | 48 | 1,138,275 | 56,746 | 1,119 | 62,520 |
| Tail land leaf rdr-0012 | 21 | 22 | 5 | 42 | 1,186,545 | 69,183 | 1,089 | 75,616 |
| ↳ Nested spawn probe | 3 | 1 | 0 | 9 | 30,628 | 13,879 | 512 | 17,102 |
| ↳ Ship rdr-0012 branch to main | 20 | 20 | 2 | 40 | 1,540,703 | 80,114 | 754 | 91,589 |
| ↳ Flip RDR 0012 to Implemented | 12 | 13 | 1 | 24 | 672,433 | 60,360 | 407 | 71,691 |
| **total** | **563** | **554** | | | | | **39,973** | |

## Totals

- context tokens processed (in + cache read + cache write): **58,108,476**
- output tokens: 39,973
- orchestrator share of context: **8%** at 122,782 peak
- agents: 16 (15 sub-agent legs)
- wall: **92 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **11 min** over 16 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **88 min** over 8 Agent waits
- guard refusals: 1 (harness blocks 1 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate (rdr-gate shard, 19:59): 94,926; at the final turn: 122,782
- parent cd-prefixed calls: 1 (into a worktree: 0)
- parent full-suite runs: 0 beyond the baseline (baseline: yes)

## Budget asks per Phase 2 leg

| leg | budget asks | commits |
|---|--:|--:|
| Phase 1 tests-first 0012 | 2 | 1 |
| Phase 2 leg 1 implement 0012 | 3 | 2 |
| Phase 3b adversarial review 0012 | 1 | 1 |
| Phase 3c fixup 0012 | 4 | 1 |
| Phase 3c fixup rdr-0012 | 5 | 2 |

## Leg discipline (helper calls against the raw commands they replace)

| leg | leg-read | raw cat | leg-test | raw go test | leg-commit | raw git commit |
|---|--:|--:|--:|--:|--:|--:|
| orchestrator | 0 | 1 | 0 | 1 | 0 | 0 |
| Phase 0 spec audit 0012 | 6 | 1 | 0 | 0 | 0 | 0 |
| Phase 1 tests-first 0012 | 20 | 2 | 10 | 0 | 2 | 0 |
| Phase 2 leg 1 implement 0012 | 29 | 1 | 10 | 0 | 3 | 0 |
| Phase 3a CoVe verify 0012 | 10 | 1 | 21 | 2 | 0 | 0 |
| Phase 3b adversarial review 0012 | 16 | 1 | 9 | 0 | 1 | 0 |
| Phase 3c fixup 0012 | 6 | 1 | 5 | 0 | 4 | 0 |
| Tail triage leaf rdr-0012 | 0 | 2 | 0 | 0 | 0 | 0 |
| Ground roborev job 8229 | 0 | 0 | 0 | 4 | 0 | 0 |
| Ground roborev job 8247 | 0 | 0 | 0 | 5 | 0 | 0 |
| Phase 3c fixup rdr-0012 | 5 | 3 | 3 | 0 | 5 | 0 |
| Ship rdr-0012 branch to main | 0 | 0 | 0 | 1 | 0 | 0 |

## Human prompts after launch

None -- the run was unattended end to end.

11 agent hand-backs re-entered the orchestrator; those are not prompts.

## Digest

```
scorecard: 58,108,476 context tokens over 16 agents
  orchestrator: 8% of context, 122,782 peak
  turns/tools:  563 / 554   output: 39,973
  wall:         92 min active, 0 min idle
  test waits:   11 min over 16 runs   waiting on legs: 88 min over 8 Agent waits
  budget asks:  15 asks / 7 commits
  human prompts after launch: 0   agent hand-backs: 11
  orchestrator ctx: gate 94,926 → final 122,782   legs: helpers 165 / raw 26
  guard refusals: 1 (harness 1 · leg-guard 0)   parent: cd-prefixed 1 (worktree 0) · full-suite runs 0 (+1 baseline)
```
