# Run scorecard

session: `9e981bb5-bac6-4eac-b871-8cb0101862fc.jsonl`  started: 2026-09-14 00:24 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 89 | 63 | 47 | 2,518 | 7,559,380 | 403,440 | 118,491 | 135,250 |
| Tail triage leaf for RDR 0021 | 183 | 128 | 29 | 5,016 | 21,463,786 | 1,297,676 | 72,001 | 192,467 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,537 | 4 | 12,540 |
| ↳ Ground empty-set null export finding | 59 | 44 | 4 | 1,798 | 2,750,958 | 292,265 | 12,878 | 81,337 |
| ↳ Ground initial-node marker finding | 67 | 51 | 5 | 2,054 | 3,547,618 | 333,141 | 21,390 | 85,819 |
| ↳ Ground export.go edge-join finding | 50 | 35 | 6 | 1,510 | 2,847,789 | 347,855 | 24,161 | 95,521 |
| ↳ Ground round-trip oracle finding | 44 | 33 | 4 | 1,318 | 1,625,670 | 258,246 | 12,966 | 62,658 |
| ↳ Phase 3c fixup | 136 | 97 | 11 | 4,262 | 11,614,325 | 553,643 | 40,025 | 140,862 |
| ↳ Phase 3c respawn: finish REQ-18 demo | 33 | 21 | 2 | 966 | 989,315 | 112,994 | 7,534 | 41,658 |
| Tail land leaf for RDR 0021 | 43 | 24 | 7 | 926 | 1,749,788 | 361,203 | 16,784 | 66,910 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,508 | 4 | 12,511 |
| ↳ Ship rdr-0021 branch to main | 25 | 15 | 3 | 740 | 727,320 | 77,675 | 6,937 | 39,395 |
| ↳ Flip RDR 0021 docs to Implemented | 17 | 10 | 1 | 454 | 660,651 | 112,772 | 3,212 | 53,373 |
| **total** | **748** | **521** | | | | | **336,387** | |

## Totals

- context tokens processed (in + cache read + cache write): **59,734,123**
- output tokens: 336,387
- orchestrator share of context: **13%** at 135,250 peak
- agents: 13 (12 sub-agent legs)
- wall: **47 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **3 min** over 6 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **37 min** over 2 Agent waits
- guard refusals: 1 (harness blocks 1 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate: no gate ask seen (looked for `rdr-gate shard`, then the first Phase 2 leg spawn); at the final turn: 135,250
- parent cd-prefixed calls: 6 (into a worktree: 2)
- parent full-suite runs: 1 beyond the baseline (baseline: no)

## Budget asks per Phase 2 leg

| leg | budget asks | commits |
|---|--:|--:|
| Tail triage leaf for RDR 0021 | 2 | 0 |
| Phase 3c fixup | 3 | 2 |
| Phase 3c respawn: finish REQ-18 demo | 1 | 1 |

## Leg discipline (helper calls against the raw commands they replace)

| leg | leg-read | raw cat | leg-test | raw go test | leg-commit | raw git commit |
|---|--:|--:|--:|--:|--:|--:|
| orchestrator | 0 | 1 | 0 | 1 | 0 | 2 |
| Tail triage leaf for RDR 0021 | 1 | 0 | 1 | 0 | 1 | 0 |
| Ground empty-set null export finding | 0 | 0 | 0 | 3 | 0 | 0 |
| Ground initial-node marker finding | 0 | 3 | 0 | 0 | 0 | 0 |
| Ground export.go edge-join finding | 0 | 3 | 0 | 7 | 0 | 0 |
| Ground round-trip oracle finding | 0 | 1 | 0 | 8 | 0 | 0 |
| Phase 3c fixup | 31 | 0 | 6 | 2 | 3 | 0 |
| Phase 3c respawn: finish REQ-18 demo | 3 | 0 | 0 | 0 | 1 | 0 |
| Ship rdr-0021 branch to main | 0 | 0 | 0 | 3 | 0 | 0 |

## Human prompts after launch

None -- the run was unattended end to end.

## Digest

```
scorecard: 59,734,123 context tokens over 13 agents
  orchestrator: 13% of context, 135,250 peak
  turns/tools:  748 / 521   output: 336,387
  wall:         47 min active, 0 min idle
  test waits:   3 min over 6 runs   waiting on legs: 37 min over 2 Agent waits
  budget asks:  6 asks / 3 commits
  human prompts after launch: 0
  orchestrator ctx: gate - → final 135,250   legs: helpers 47 / raw 34
  guard refusals: 1 (harness 1 · leg-guard 0)   parent: cd-prefixed 6 (worktree 2) · full-suite runs 1 (+0 baseline)
```
