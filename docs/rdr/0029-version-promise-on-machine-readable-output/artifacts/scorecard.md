# Run scorecard

session: `3882fd5d-8482-4c1f-a549-d563a4354f57.jsonl`  started: 2026-09-13 16:46 UTC

| agent | turns | tools | wall min | in | cache r | cache w | out | peak ctx |
|---|--:|--:|--:|--:|--:|--:|--:|--:|
| orchestrator | 19 | 14 | 1 | 518 | 862,431 | 142,903 | 19,438 | 58,449 |
| **total** | **19** | **14** | | | | | **19,438** | |

## Totals

- context tokens processed (in + cache read + cache write): **1,005,852**
- output tokens: 19,438
- orchestrator share of context: **100%** at 58,449 peak
- agents: 1 (0 sub-agent legs)
- wall: **1 min active**, 0 min idle
- test waits (a Bash call that runs a test -- rdr-leg-test or a raw test/lint command; never an Agent/TaskOutput wait): **0 min** over 0 blocking runs (>20s)
- waiting on legs (the orchestrator's own wall time between spawning a sub-agent and that leg going quiet; never summed with test waits): **0 min** over 0 Agent waits
- guard refusals: 0 (harness blocks 0 · leg-guard denies 0)
- orchestrator context at the Phase 1 gate: no gate ask seen (looked for `rdr-gate shard`, then the first Phase 2 leg spawn); at the final turn: 58,449
- parent cd-prefixed calls: 0 (into a worktree: 0)
- parent full-suite runs: 0 beyond the baseline (baseline: no)

## Budget asks per Phase 2 leg

No `rdr-leg-commit` budget asks recorded. A Phase 2 leg that never
asks is the 965K-context failure mode -- check whether one ran.

## Leg discipline (helper calls against the raw commands they replace)

No helper calls or raw reads/tests/commits recorded.

## Human prompts after launch

None -- the run was unattended end to end.

## Digest

```
scorecard: 1,005,852 context tokens over 1 agents
  orchestrator: 100% of context, 58,449 peak
  turns/tools:  19 / 14   output: 19,438
  wall:         1 min active, 0 min idle
  test waits:   0 min over 0 runs   waiting on legs: 0 min over 0 Agent waits
  budget asks:  0 asks / 0 commits
  human prompts after launch: 0
  orchestrator ctx: gate - → final 58,449   legs: helpers 0 / raw 0
  guard refusals: 0 (harness 0 · leg-guard 0)   parent: cd-prefixed 0 (worktree 0) · full-suite runs 0 (+0 baseline)
```
