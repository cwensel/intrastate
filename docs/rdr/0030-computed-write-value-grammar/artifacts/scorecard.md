# Run scorecard

session: `59719e15-8f8d-4905-b840-e56830958852.jsonl`  started: 2026-09-23 23:03 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 71 | 59 | 124 | 142 | 8,014,345 | 144,706 | 34,282 | 158,596 |
| Phase 0 spec audit rdr-0030 | 45 | 48 | 10 | 90 | 4,705,137 | 200,612 | 2,079 | 206,386 |
| Ground Phase 0 Q1 decision | 17 | 16 | 1 | 34 | 646,888 | 48,277 | 1,018 | 54,051 |
| Phase 1 tests-first rdr-0030 | 107 | 107 | 30 | 214 | 23,840,185 | 414,129 | 9,556 | 419,903 |
| Phase 2 leg 1 rdr-0030 | 83 | 86 | 14 | 166 | 14,304,320 | 248,719 | 2,604 | 254,514 |
| Phase 3a CoVe verifier rdr-0030 | 61 | 61 | 10 | 122 | 6,467,904 | 156,047 | 1,740 | 161,821 |
| Phase 3b adversarial review rdr-0030 | 45 | 48 | 10 | 90 | 4,530,312 | 158,256 | 3,564 | 164,030 |
| Phase 3c fixup rdr-0030 | 69 | 72 | 9 | 138 | 6,566,332 | 132,216 | 4,336 | 137,990 |
| Ground D13 spec-defect decision | 28 | 28 | 2 | 56 | 1,233,869 | 55,128 | 1,150 | 60,902 |
| Strong consult on D13 fork | 47 | 48 | 8 | 94 | 4,035,769 | 134,187 | 2,004 | 139,140 |
| Phase 3c leg 2 apply D13 | 28 | 27 | 4 | 56 | 1,564,513 | 75,907 | 952 | 81,681 |
| Tail triage leaf rdr-0030 | 31 | 29 | 11 | 66 | 1,457,603 | 113,806 | 1,176 | 65,404 |
| ↳ Nested spawn check | 4 | 1 | 0 | 38 | 45,493 | 17,478 | 757 | 17,488 |
| ↳ Ground roborev job 8279 finding | 23 | 23 | 3 | 46 | 1,256,036 | 73,578 | 1,092 | 79,352 |
| ↳ Phase 3c fixup rdr-0030 | 22 | 21 | 3 | 44 | 1,132,656 | 62,984 | 645 | 68,758 |
| Tail land leaf rdr-0030 | 23 | 24 | 5 | 46 | 1,213,835 | 63,005 | 955 | 68,792 |
| ↳ Nested spawn probe | 4 | 1 | 0 | 38 | 44,751 | 17,154 | 399 | 17,164 |
| ↳ Ship rdr-0030 branch | 14 | 15 | 3 | 28 | 622,386 | 50,432 | 298 | 56,226 |
| ↳ Flip RDR 0030 docs | 14 | 14 | 1 | 28 | 796,520 | 71,503 | 455 | 71,505 |
| Tail flight leaf rdr-0030 | 20 | 20 | 4 | 42 | 984,816 | 137,077 | 445 | 77,383 |
| ↳ Spawn probe | 4 | 1 | 0 | 38 | 44,764 | 17,189 | 293 | 17,199 |
| ↳ Scope-review kata c8ge | 16 | 20 | 3 | 32 | 819,375 | 59,427 | 799 | 70,699 |
| **total** | **776** | **769** | | | | | **70,599** | |

## Totals

- context tokens processed (in + cache read + cache write): **86,781,274**
- output tokens: 70,599
- orchestrator share of context: **9%** at 158,596 peak
- agents: 22 (21 sub-agent legs)
- wall: **124 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **11 min** over 17 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **120 min** over 13 Agent waits
- guard refusals: 0 (harness blocks 0 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate (rdr-gate shard, 23:47): 103,135; at the final turn: 158,596
- parent cd-prefixed calls: 4 (into a worktree: 0)
- parent full-suite runs: 0 beyond the baseline (baseline: yes)

## Budget asks per Phase 2 leg

| leg | budget asks | commits |
|---|--:|--:|
| Phase 1 tests-first rdr-0030 | 4 | 1 |
| Phase 2 leg 1 rdr-0030 | 2 | 2 |
| Phase 3b adversarial review rdr-0030 | 2 | 1 |
| Phase 3c fixup rdr-0030 | 3 | 3 |
| Phase 3c leg 2 apply D13 | 1 | 1 |
| Phase 3c fixup rdr-0030 | 1 | 1 |

## Leg discipline (helper calls against the raw commands they replace)

| leg | leg-read | raw cat | leg-test | raw go test | leg-commit | raw git commit |
|---|--:|--:|--:|--:|--:|--:|
| orchestrator | 0 | 0 | 0 | 1 | 0 | 0 |
| Phase 0 spec audit rdr-0030 | 14 | 1 | 0 | 0 | 0 | 0 |
| Ground Phase 0 Q1 decision | 3 | 0 | 0 | 0 | 0 | 0 |
| Phase 1 tests-first rdr-0030 | 36 | 2 | 16 | 0 | 4 | 0 |
| Phase 2 leg 1 rdr-0030 | 26 | 1 | 10 | 0 | 2 | 0 |
| Phase 3a CoVe verifier rdr-0030 | 11 | 1 | 7 | 0 | 0 | 0 |
| Phase 3b adversarial review rdr-0030 | 14 | 1 | 8 | 0 | 2 | 0 |
| Phase 3c fixup rdr-0030 | 11 | 8 | 6 | 0 | 3 | 0 |
| Ground D13 spec-defect decision | 4 | 1 | 0 | 0 | 0 | 0 |
| Strong consult on D13 fork | 0 | 1 | 0 | 9 | 0 | 0 |
| Phase 3c leg 2 apply D13 | 12 | 1 | 5 | 0 | 1 | 0 |
| Tail triage leaf rdr-0030 | 0 | 1 | 0 | 0 | 0 | 0 |
| Ground roborev job 8279 finding | 0 | 1 | 0 | 5 | 0 | 0 |
| Phase 3c fixup rdr-0030 | 6 | 1 | 2 | 0 | 1 | 0 |
| Ship rdr-0030 branch | 0 | 0 | 0 | 2 | 0 | 0 |

## Human prompts after launch

None -- the run was unattended end to end.

16 agent hand-backs re-entered the orchestrator; those are not prompts.

## Digest

```
scorecard: 86,781,274 context tokens over 22 agents
  orchestrator: 9% of context, 158,596 peak
  turns/tools:  776 / 769   output: 70,599
  wall:         124 min active, 0 min idle
  test waits:   11 min over 17 runs   waiting on legs: 120 min over 13 Agent waits
  budget asks:  13 asks / 9 commits
  human prompts after launch: 0   agent hand-backs: 16
  orchestrator ctx: gate 103,135 → final 158,596   legs: helpers 204 / raw 37
  guard refusals: 0 (harness 0 · leg-guard 0)   parent: cd-prefixed 4 (worktree 0) · full-suite runs 0 (+1 baseline)
```
