# Run scorecard

session: `09736056-427e-41a5-a0e2-95e278560827.jsonl`  started: 2026-09-23 19:26 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 79 | 42 | 92 | 158 | 7,246,512 | 198,089 | 36,762 | 122,782 |
| Phase 0 spec audit 0012 | 79 | 49 | 7 | 158 | 8,378,323 | 356,358 | 2,254 | 190,137 |
| Phase 1 tests-first 0012 | 141 | 81 | 23 | 282 | 19,723,509 | 491,133 | 2,641 | 260,431 |
| Phase 2 leg 1 implement 0012 | 148 | 88 | 13 | 296 | 26,730,208 | 472,884 | 3,172 | 263,468 |
| Phase 3a CoVe verify 0012 | 104 | 59 | 10 | 208 | 9,766,029 | 257,033 | 1,984 | 146,914 |
| Phase 3b adversarial review 0012 | 88 | 50 | 9 | 176 | 7,918,498 | 246,449 | 1,820 | 143,687 |
| Phase 3c fixup 0012 | 82 | 46 | 10 | 164 | 6,350,621 | 174,723 | 1,628 | 110,169 |
| Tail triage leaf rdr-0012 | 55 | 31 | 10 | 120 | 2,589,651 | 169,887 | 2,028 | 67,120 |
| ↳ Nested spawn probe | 3 | 1 | 0 | 9 | 30,910 | 13,983 | 369 | 17,184 |
| ↳ Ground roborev job 8229 | 29 | 16 | 2 | 58 | 1,018,441 | 76,035 | 844 | 46,152 |
| ↳ Ground roborev job 8247 | 19 | 11 | 2 | 38 | 687,801 | 77,354 | 996 | 50,637 |
| ↳ Phase 3c fixup rdr-0012 | 42 | 24 | 4 | 84 | 1,965,407 | 105,344 | 1,271 | 62,520 |
| Tail land leaf rdr-0012 | 43 | 22 | 5 | 86 | 2,406,111 | 139,891 | 1,330 | 75,616 |
| ↳ Nested spawn probe | 3 | 1 | 0 | 9 | 30,628 | 13,879 | 512 | 17,102 |
| ↳ Ship rdr-0012 branch to main | 37 | 20 | 2 | 74 | 2,745,518 | 174,794 | 848 | 91,589 |
| ↳ Flip RDR 0012 to Implemented | 25 | 13 | 1 | 50 | 1,354,714 | 146,524 | 491 | 71,691 |
| **total** | **977** | **554** | | | | | **58,950** | |

## Totals

- context tokens processed (in + cache read + cache write): **102,059,211**
- output tokens: 58,950
- orchestrator share of context: **7%** at 122,782 peak
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
scorecard: 102,059,211 context tokens over 16 agents
  orchestrator: 7% of context, 122,782 peak
  turns/tools:  977 / 554   output: 58,950
  wall:         92 min active, 0 min idle
  test waits:   11 min over 16 runs   waiting on legs: 88 min over 8 Agent waits
  budget asks:  15 asks / 7 commits
  human prompts after launch: 0   agent hand-backs: 11
  orchestrator ctx: gate 94,926 → final 122,782   legs: helpers 165 / raw 26
  guard refusals: 1 (harness 1 · leg-guard 0)   parent: cd-prefixed 1 (worktree 0) · full-suite runs 0 (+1 baseline)
```
