# Run scorecard

session: `59719e15-8f8d-4905-b840-e56830958852.jsonl`  started: 2026-09-23 23:03 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 133 | 59 | 124 | 266 | 14,905,028 | 303,753 | 70,384 | 158,596 |
| Phase 0 spec audit rdr-0030 | 78 | 48 | 10 | 156 | 8,165,699 | 378,091 | 2,304 | 206,386 |
| Ground Phase 0 Q1 decision | 24 | 16 | 1 | 48 | 886,770 | 81,628 | 1,159 | 54,051 |
| Phase 1 tests-first rdr-0030 | 197 | 107 | 30 | 394 | 44,764,717 | 809,531 | 10,170 | 419,903 |
| Phase 2 leg 1 rdr-0030 | 149 | 86 | 14 | 298 | 25,486,809 | 478,766 | 3,214 | 254,514 |
| Phase 3a CoVe verifier rdr-0030 | 113 | 61 | 10 | 226 | 12,029,444 | 295,062 | 2,135 | 161,821 |
| Phase 3b adversarial review rdr-0030 | 84 | 48 | 10 | 168 | 8,639,550 | 313,645 | 3,872 | 164,030 |
| Phase 3c fixup rdr-0030 | 115 | 72 | 9 | 230 | 11,159,663 | 234,980 | 4,768 | 137,990 |
| Ground D13 spec-defect decision | 44 | 28 | 2 | 88 | 1,885,851 | 101,563 | 1,359 | 60,902 |
| Strong consult on D13 fork | 77 | 48 | 8 | 154 | 6,614,106 | 227,241 | 2,207 | 139,140 |
| Phase 3c leg 2 apply D13 | 46 | 27 | 4 | 92 | 2,614,918 | 134,032 | 1,055 | 81,681 |
| Tail triage leaf rdr-0030 | 50 | 29 | 11 | 106 | 2,309,723 | 198,505 | 1,506 | 65,404 |
| ↳ Nested spawn check | 7 | 1 | 0 | 68 | 76,817 | 31,669 | 765 | 17,488 |
| ↳ Ground roborev job 8279 finding | 37 | 23 | 3 | 74 | 1,975,154 | 133,638 | 1,194 | 79,352 |
| ↳ Phase 3c fixup rdr-0030 | 37 | 21 | 3 | 74 | 1,914,135 | 106,801 | 772 | 68,758 |
| Tail land leaf rdr-0030 | 38 | 24 | 5 | 76 | 1,953,884 | 107,483 | 1,171 | 68,792 |
| ↳ Nested spawn probe | 7 | 1 | 0 | 68 | 75,608 | 31,070 | 408 | 17,164 |
| ↳ Ship rdr-0030 branch | 25 | 15 | 3 | 50 | 1,057,896 | 103,866 | 436 | 56,226 |
| ↳ Flip RDR 0030 docs | 28 | 14 | 1 | 56 | 1,548,986 | 155,956 | 494 | 71,505 |
| Tail flight leaf rdr-0030 | 37 | 20 | 4 | 78 | 1,892,309 | 257,067 | 1,056 | 77,724 |
| ↳ Spawn probe | 7 | 1 | 0 | 68 | 75,657 | 31,082 | 301 | 17,199 |
| ↳ Scope-review kata c8ge | 29 | 20 | 3 | 58 | 1,392,733 | 147,065 | 873 | 70,699 |
| **total** | **1362** | **769** | | | | | **111,603** | |

## Totals

- context tokens processed (in + cache read + cache write): **156,090,847**
- output tokens: 111,603
- orchestrator share of context: **10%** at 158,596 peak
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
scorecard: 156,090,847 context tokens over 22 agents
  orchestrator: 10% of context, 158,596 peak
  turns/tools:  1362 / 769   output: 111,603
  wall:         124 min active, 0 min idle
  test waits:   11 min over 17 runs   waiting on legs: 120 min over 13 Agent waits
  budget asks:  13 asks / 9 commits
  human prompts after launch: 0   agent hand-backs: 16
  orchestrator ctx: gate 103,135 → final 158,596   legs: helpers 204 / raw 37
  guard refusals: 0 (harness 0 · leg-guard 0)   parent: cd-prefixed 4 (worktree 0) · full-suite runs 0 (+1 baseline)
```
