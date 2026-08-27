# Triage — RDR 0011 `flow next` match-conditioned candidates

Single unattended pass. Window frozen once: `BASE 58f8e41` → `HEAD 6638741`
(24 commits). Batch label `batch:rdr-0011`.

Spine: 12 in-window per-commit auto-reviews (jobs 6295–6308), 20 findings.
Cross-commit net: job 6309 over the same range — 3 findings, all duplicates
of spine findings, charged to the per-commit job. Fix-commit sweep: job 6310
over `1973323` — clean, no second round.

Every finding reached exactly one terminal state and its job is closed.

## Findings × verdict

| Job | Location | Sev | Verdict | odc-type / trigger | Outcome |
|---|---|---|---|---|---|
| 6295 | `flow_fixtures_0011_test.go:884` S8 arms | Med | DROP rdr-adjudicated | test-oracle / rare-situation | DEV-4 created `flowMatchOnlyOwnedSoloModel` for this exact collision |
| 6296 | `flow_next_0011_test.go:156` gate exclusion | Med | KATA-BUG | test-oracle / test-coverage | `g1pa` |
| 6296 | `flow_next_0011_test.go:513` S5 key set | Med | DROP rdr-adjudicated | test-oracle / design-conformance | DEV-3: `assemble`/`TagSet` key set unexported, S5 forbids exporting |
| 6296 | `flow_fixtures_0011_test.go` kernel diff | Med | KATA-BUG | environment / rare-situation | `8ka9` |
| 6296 | `flow_next_0011_test.go:464` empty string | Low | KATA-BUG | test-oracle / boundary | `yn1v` |
| 6297 | `flow_all_0011_test.go:253` gated + `--all` | Low | DROP over-engineering | test-oracle / test-coverage | C2 scopes `--all` to the match predicate; cross-product unrequired |
| 6297 | `flow_all_0011_test.go:526` duplicate flag map | Low | DROP rdr-adjudicated | test-oracle / doc-code-drift | REQ-65: "These are TWO oracles, not one"; C3 forbids editing 0005's map |
| 6298 | `flow_unknown_0011_test.go:520` element shape | Med | DROP rdr-adjudicated | test-oracle / boundary | DEV-7: `stringsAt` returns `([],true)` on empty in every build |
| 6298 | `flow_unknown_0011_test.go:319` pair dedup/sort | Med | KATA-BUG | test-oracle / test-coverage | `yn1v` |
| 6298 | `flow_unknown_0011_test.go:464` reason leaf | Med | FIX-NOW | test-oracle / test-coverage | `1973323` |
| 6299 | `flow_demand_0011_test.go:174` reader set | Med | DROP rdr-adjudicated | test-oracle / rare-situation | DEV-6 — found and fixed during the launch; written against superseded state |
| 6299 | `flow_demand_0011_test.go:127` cross-outcome | Med | KATA-BUG | test-oracle / test-coverage | `gq75` |
| 6299 | `flow_demand_0011_test.go:414` shipped fixtures | Med | KATA-BUG | test-oracle / test-coverage | `gq75` |
| 6300 | `flow_mvv_0011_test.go:426` kernel untouched | Med | KATA-BUG | environment / rare-situation | `8ka9` |
| 6300 | `flow_rehome_0011_test.go:72` census count | Med | DROP rdr-adjudicated | test-oracle / design-conformance | C3's `file:line` are lock-time provenance; the COUNT is the contract |
| 6300 | `flow_rehome_0011_test.go:395` CLI vs kernel | Med | DROP rdr-adjudicated | test-oracle / design-conformance | Same DEV-3 limit as 6296; reported twice |
| 6300 | `flow_rehome_0011_test.go:260` help contract | Med | KATA-BUG | test-oracle / design-conformance | `m4f2` |
| 6301 | `flow_fixtures_0011_test.go:1496` prod scope | Med | KATA-BUG | environment / rare-situation | `8ka9` |
| 6301 | `flow_next_0011_test.go:652` probe shape | Med | KATA-BUG | test-oracle / logic-flow | `g1pa` |
| 6302 | `flow_next.go:93` unused constant | Med | DROP superseded-at-HEAD | build-package-merge / logic-flow | referenced at HEAD; `make check` green |
| 6304 | `cli-output-contract.md:182` example order | Low | FIX-NOW | documentation / doc-code-drift | `1973323` |
| 6306 | ADV tests fail CI | High | DROP superseded-at-HEAD | n/a-not-a-defect | deliberately-red Phase 3b probes; adjudicated at `fd84ab7`/`f0233e3` |
| 6307 | ADV-2 contradictory | High | DROP superseded-at-HEAD | n/a-not-a-defect | resolved at `f0233e3` exactly as recommended |
| 6308 | ADV-2 misleading name | Low | DROP superseded-at-HEAD | documentation / doc-code-drift | renamed in the same commit family |

## Terminal-state counts

- **DROP 11** — rdr-adjudicated 7, superseded-at-HEAD 3 (+1 low), over-engineering 1
- **FIX-NOW 2** — both in commit `1973323`
- **KATA-BUG 5 katas** from 9 findings (consolidated by root cause)
- **RDR-SEED 0** — no finding was contract-level

## Fix-now detail — `1973323`

1. `docs/cli-output-contract.md` — the `unknown` example placed
   `gate_passed` before `approval`, contradicting the `(key, reason)` sort
   the contract states two lines above and `flow_next.go::dedupeUnknown`
   implements. Example reordered; the implementation was already correct.
2. `internal/cli/flow_unknown_0011_test.go` — `TestReq72And73` scanned for
   the bare substrings `wanted`/`absent`, but the fixture's rule ids
   (`absent-row`, `all-absent`) carry "absent" on their own, so a build
   emitting no `reason` leaf would pass. Replaced with path-qualified
   `key:`/`reason:` leaf assertions and mutation-checked (the assertion
   fails when the expected leaf is wrong). Strengthened, never relaxed.

`make check` green after both. Fix-commit auto-review (job 6310): no issues.

## Katas filed — all `batch:rdr-0011`, `src:roborev`, `lifecycle:queued`

| id | title | sev | pri |
|---|---|---|---|
| `8ka9` | Fail closed on the kernel-untouched and production-diff-scope git guards | Med | 1 |
| `yn1v` | Harden the unknown-payload oracles (pair dedup, reason sort, empty-string) | Med | 2 |
| `gq75` | Pin the demand-set boundary conditions (cross-outcome, shipped fixtures) | Med | 2 |
| `g1pa` | Make the probe-shape and gate-exclusion oracles observe what they claim | Med | 3 |
| `m4f2` | Assert `flow next` help states REQ-51's semantic clauses, not substrings | Low | 3 |

No kata carries an open question — every finding was settled from evidence.

## Notes

- No finding indicated broken shipped behaviour. All 20 are oracle-strength
  or documentation claims; the suite is green at HEAD and 130/130 REQs hold.
- 7 of 20 DROP as `rdr-adjudicated`. The launch's `deviations.md` (DEV-3,
  DEV-4, DEV-6, DEV-7) and the REQ-54/60 non-deviation note each pre-answer
  a finding; DEV-6 is a finding the launch itself found and fixed in Phase 2.
- Two findings (6296:513 and 6300:395) are the same oracle reported twice.
- `8ka9` is the highest-value item: it is the only finding where the RDR's
  own named verification mechanism (`0011:S5`, MVV step 9) is inert in
  every CI mode.
- Line drift: the kernel-diff finding cites `flow_fixtures_0011_test.go:1298`,
  which is `unknownKeys`; the real sites are `gitDiffStat` (:1462) and
  `changedFilesSince` (:1497). Re-anchored by symbol per source-location
  discipline.
