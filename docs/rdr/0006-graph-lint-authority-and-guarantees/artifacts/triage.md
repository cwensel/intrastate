# Triage — RDR 0006 Graph Lint Authority And Guarantees

Unattended single-pass roborev triage after the Stage 8 launch.
Window frozen once: `BASE=3bd7d83` .. `HEAD=c60a8c0` (24 commits).
Batch label: `batch:rdr-0006`.

Spine: 15 in-window per-commit auto-reviews (jobs 6130–6157), 47 findings.
Cross-commit net: job 6160 (`review --since BASE`), 8 findings — **all** duplicates
of spine findings (6145×5, 6146×2, 6147×1), charged to their per-commit jobs.
Net-new from the range review: zero.

17 out-of-window jobs (pre-`BASE` on `main`/other worktrees, plus 4507/6095) were
skipped — they belong to earlier work.

## Routing

| Finding | Verdict | odc-type | odc-trigger | Evidence | Outcome |
|---|---|---|---|---|---|
| 6130 High — stubs make every model report clean | pre-filter | n/a-not-a-defect | — | ref `21dc7e2` IS the intentionally-RED gate commit; `skeleton.go` gone at HEAD, census zero | excluded |
| 6130 Med — `--flow` advertised but unresolved | DROP rdr-adjudicated | interface | design-conformance | RDR 0006 A9 (352-369) + Prerequisites (1429-1432): `--flow` waits on RDR 0005 config discovery | drop |
| 6134 ×4, 6135 ×5, 6136 ×4, 6138 ×5, 6139 ×3, 6140 ×2 — weak oracles | KATA-BUG ×2 | checking | test-coverage | verified at HEAD; split by whether a fixture or an assert is at fault | `py0x`, `aqzt` |
| 6139 #1 — MVV omits `graph-single-valued-state` | DROP rdr-adjudicated | checking | design-conformance | **D3**: RDR 0002's loader refuses every TOML spelling; the fixture is unauthorable | drop |
| 6140 #6 — declaration-model test never invokes graph lint | DROP over-engineering | test-oracle | test-coverage | `soundness_0006_test.go:541-560` is a REQ-45 *agreement* test, bounded by design; behaviour covered by REQ-54/56/80/81 | drop |
| 6134 #2, 6140 #5 — escape self-loop / REQ-112 discriminator | DROP superseded-at-HEAD | test-oracle | test-coverage | Phase 3c's D11 rewrote the test to count nodes per owned-state identity | drop |
| 6135 #3, 6139 #2 — `ProductBound` numeric fragility | DROP over-engineering | test-oracle | test-coverage | `ProductBound() == guard.Bound() == 2048`; 1,000,001 exceeds it ~500×; both tests read the published constant | drop |
| 6140 #1/#2/#3 — repository gate tested by a non-gate path | RDR-SEED | test-oracle | design-conformance | beyond **D4**'s disposition (existence discharged, fidelity not) | `9en6` |
| 6145 #2 — `OpaqueValue` never assigned | KATA-BUG | assignment | rare-situation | false-green **reproduced**: blocking overlap masked at exit 0 | `ffgb` |
| 6145 #1, #3 | DROP superseded-at-HEAD | algorithm | logic-flow | fixed by `6e14a06` (successor join) and `a09174c` (invariant 7 per node) | drop |
| 6145 #4, 6154, 6157 — delimiter collisions | KATA-BUG (collapsed) | algorithm | rare-situation | authorability **verified**: `table.Load` accepts a tag named `a;b`; RDR 0002 REQ-126/127 require `,`-bearing members stay distinct | `04xp` |
| 6145 #5 — `in` missing from `singleValueOperators` | **FIX-NOW** | assignment | logic-flow | `guard.isSingleValueOperator` includes `in`; RDR line 733 names it | `d47d5fb` |
| 6146 #1 — dimensionless groups run coverage arms | DROP unreachable | algorithm | logic-flow | `coverageUnionFor`'s dimensionless branch is an EXISTENTIAL check: any ordinary row returns the full product. The reported blocking `no_match` gap needs a group of only escape rows not declaring `no_match` — which is the correct REQ-62 outcome. Orchestrator verified the control flow directly; the reporting leaf reasoned from code shape without executing | drop |
| 6146 #2 — bare escape reported as closing | DROP rdr-adjudicated | checking | design-conformance | **D8** adjudicates verbatim; `coverage.go:291-296` carries D8's reasoning; REQ-65 "a bare green MUST NOT satisfy this clause" | drop |
| 6147 — identity values unescaped in text mode | KATA-BUG | checking | boundary | live, but RDR 750-752 promises no one-finding-per-line guarantee | `a5an` |
| 6151 #1 — ADV-1 contradicts D9 | DROP superseded-at-HEAD | test-oracle | rare-situation | **D11 explicitly supersedes D9**; ADV-1's oracle is now the correct one | drop |
| 6151 #2 — can-refuse rows and overlap | DROP rdr-adjudicated | test-oracle | rare-situation | **D13**: REQ-84 unqualified, REQ-110 existential; RDR 0003's clause narrows *coverage*, not overlap | drop |
| 6151 #3 — ADV-3 demands non-terminal split | DROP rdr-adjudicated | test-oracle | rare-situation | REQ-37 MUST; the reviewer's prescription is already the test's second arm | drop |
| 6155 — dead-end pin vs FAIL-4 | RDR-SEED | test-oracle | design-conformance | **D12** (`needs author decision`); reviewer's premise stale (D12 landed in `9265396`, after the reviewed commit) | `pz9z` |

## Fix-now applied

`d47d5fb` — `fix(graphlint): classify an in atom as single-value for the stable reason`.
Suite green (`go test ./...`), `golangci-lint run` 0 issues, `models/rdr.toml` census
still zero.

## Notes on reporter accuracy

Two roborev findings and one grounding verdict did not survive checking, and are
recorded as drops with the reason rather than silently discarded:

- 6146 #1 was reported as a live blocking defect and re-reported by the range net.
  Reading the control flow shows the branch is existential; no loader-accepted model
  produces the described symptom.
- 6151 #1 and 6155 both reason from deviation state that a later commit had already
  changed (D9→D11; D12's landing commit).
