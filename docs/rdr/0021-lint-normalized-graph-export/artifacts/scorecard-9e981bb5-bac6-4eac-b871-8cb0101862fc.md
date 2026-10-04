# Run scorecard

session: `9e981bb5-bac6-4eac-b871-8cb0101862fc.jsonl`  started: 2026-09-14 00:24 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 26 | 63 | 47 | 742 | 2,371,038 | 119,772 | 35,305 | 135,250 |
| Tail triage leaf for RDR 0021 | 55 | 128 | 29 | 1,430 | 6,537,328 | 521,076 | 71,654 | 192,467 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,537 | 4 | 12,540 |
| ↳ Ground empty-set null export finding | 15 | 44 | 4 | 450 | 748,714 | 75,037 | 12,778 | 81,337 |
| ↳ Ground initial-node marker finding | 16 | 51 | 5 | 482 | 879,293 | 79,184 | 21,279 | 85,819 |
| ↳ Ground export.go edge-join finding | 15 | 35 | 6 | 450 | 935,407 | 88,689 | 22,823 | 95,521 |
| ↳ Ground round-trip oracle finding | 11 | 33 | 4 | 322 | 452,289 | 57,128 | 12,888 | 62,658 |
| ↳ Phase 3c fixup | 39 | 97 | 11 | 1,218 | 3,624,976 | 132,121 | 39,623 | 140,862 |
| ↳ Phase 3c respawn: finish REQ-18 demo | 12 | 21 | 2 | 354 | 383,228 | 35,773 | 7,478 | 41,658 |
| Tail land leaf for RDR 0021 | 19 | 24 | 7 | 428 | 787,471 | 164,159 | 16,708 | 66,910 |
| ↳ Spawn probe | 1 | 0 | 0 | 3 | 0 | 12,508 | 4 | 12,511 |
| ↳ Ship rdr-0021 branch to main | 10 | 15 | 3 | 290 | 291,310 | 33,161 | 6,912 | 39,395 |
| ↳ Flip RDR 0021 docs to Implemented | 7 | 10 | 1 | 194 | 278,748 | 47,843 | 3,192 | 53,373 |
| **total** | **227** | **521** | | | | | **250,648** | |

## Totals

- context tokens processed (in + cache read + cache write): **18,675,156**
- output tokens: 250,648
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
scorecard: 18,675,156 context tokens over 13 agents
  orchestrator: 13% of context, 135,250 peak
  turns/tools:  227 / 521   output: 250,648
  wall:         47 min active, 0 min idle
  test waits:   3 min over 6 runs   waiting on legs: 37 min over 2 Agent waits
  budget asks:  6 asks / 3 commits
  human prompts after launch: 0   agent hand-backs: 0
  orchestrator ctx: gate - → final 135,250   legs: helpers 47 / raw 34
  guard refusals: 1 (harness 1 · leg-guard 0)   parent: cd-prefixed 6 (worktree 2) · full-suite runs 1 (+0 baseline)
```
