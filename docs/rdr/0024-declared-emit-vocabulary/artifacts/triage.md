# Triage — RDR 0024 (roborev, unattended)

`docs/rdr/0024-declared-emit-vocabulary.md` @ `bb7c8fc` (base `1cd38d3`)
batch label: `batch:rdr-0024`

16 jobs in window (15 per-commit + 1 cross-commit range net), 26 findings.
Every finding is in exactly one terminal state and its job is closed.

## Verdicts

| # | Job | Finding | Verdict | ODC type / trigger | Outcome |
|---|-----|---------|---------|--------------------|---------|
| 1 | 6475 | `set` fixture carries two defects | DROP superseded-at-HEAD | n/a-not-a-defect | D1 removed the `elements` line |
| 2 | 6475 | narrowing leg never narrows | KATA-BUG | test-oracle / test-coverage | `4v92` |
| 3 | 6476 | disposition-table arm only under `bool` | KATA-BUG | test-oracle / test-coverage | `hpnv` |
| 4 | 6476 | non-`verdict` arms not single-defect | KATA-BUG | test-oracle / test-coverage | `hpnv` |
| 5 | 6477 | locator should return `[rule.emit]` header | DROP rdr-adjudicated | n/a-not-a-defect | REQ-32 inverted by reviewer |
| 6 | 6477 | ordering not proven by `emit_proof:55` | DROP superseded-at-HEAD | n/a-not-a-defect | witnessed at `emit_steps:66`; residual → `5hy1` B(ii) |
| 7 | 6478 | domain-member byte-exactness untested | KATA-BUG | test-oracle / boundary | `jrp2` |
| 8 | 6478 | flat-enum `nil Dispositions` unpinned | KATA-BUG | test-oracle / boundary | `jrp2` |
| 9 | 6478 | kernel-isolation field-name scan | DROP over-engineering | test-oracle / design-conformance | `TestReq56`'s `Dump` oracle is stronger |
| 10 | 6479 | witnesses outside the category floor | KATA-BUG | test-oracle / design-conformance | `5hy1` leg C |
| 11 | 6480 | `dtDecl0024` omits `code` | DROP superseded-at-HEAD | n/a-not-a-defect | D9 |
| 12 | 6481 | `declaredDT0024` clean models refuse | DROP superseded-at-HEAD | n/a-not-a-defect | D9 (same root as #11) |
| 13 | 6481 | `JDR 0002` section under-specified | FIX-NOW | test-oracle / design-conformance | `bb7c8fc` |
| 14 | 6482 | pricing must join non-empty dispositions | FIX-NOW | test-oracle / design-conformance | `bb7c8fc` (model was the root cause) |
| 15 | 6482 | member checks substring-match TOML | KATA-BUG | test-oracle / test-coverage | `5hy1` leg D |
| 16 | 6482 | routing detection substring-matches | KATA-BUG | test-oracle / test-coverage | `5hy1` leg D |
| 17 | 6483 | MVV locator accepts any non-`:1` line | KATA-BUG | test-oracle / test-coverage | `5hy1` leg A |
| 18 | 6484 | step binding ×3 (receiver, order, `atLine(` count) | KATA-BUG | test-oracle / test-coverage | `5hy1` leg B |
| 19 | 6485 | `emitRuleLine` returns the id line | DROP rdr-adjudicated | n/a-not-a-defect | REQ-32 mandates it (same inversion as #5) |
| 20 | 6486 | dispositions fold removes populated objects | KATA-BUG | test-oracle / backward-compat | `av9m` |
| 21 | 6487 | doc promises `[rule.emit]` block line | DROP rdr-adjudicated | n/a-not-a-defect | same inversion as #5/#19 |
| 22 | 6487 | pricing partitions `dpa` against `PH4` | FIX-NOW | documentation / design-conformance | `bb7c8fc` |
| 23 | 6488 | ADV-3 outside the recorded contract | DROP rdr-adjudicated | n/a-not-a-defect | D10 (Phase 3c reached this independently) |
| 24 | 6488 | hash-id case not a legal single-defect model | DROP unreachable | n/a-not-a-defect | premise inverted — D11 step order |
| 25 | 6488 | ADV-1 asserts only a positive line | KATA-BUG | test-oracle / test-coverage | `wkqs` |
| 26 | 6489 | rule discovery requires literal `[[rule]]` | FIX-NOW | algorithm / boundary | `bb7c8fc` |
| 27 | 6489 | `scalarAssignment` drops TOML escapes | KATA-BUG | algorithm / rare-situation | `wkqs` |
| 28 | 6491 | quoted rule headers still collapse to `:1` | KATA-BUG | algorithm / boundary | `m6fz` (bounded sweep — filed, not re-fixed) |

Job 6490 (range net) raised four findings, all deduped against the per-commit
spine (#26, #3, #22, #20) and charged there.

## Counts

- findings: 28 (spine 24, range-net 4 deduped, +1 fix-commit sweep)
- dropped: 9 — superseded-at-HEAD 4, rdr-adjudicated 4, unreachable 1,
  over-engineering 1 (#9 counted once)
- fixed-now: 4 findings in one commit `bb7c8fc`
- kata-bug: 13 findings → 6 kata (`5hy1` `jrp2` `hpnv` `av9m` `wkqs` `4v92`)
  + 1 from the bounded fix-commit sweep (`m6fz`)
- rdr-seed: 0

ODC (defects only, excluding `n/a-not-a-defect`): test-oracle 11,
algorithm 3, documentation 1.

## The one finding class that mattered most

Three separate findings (#5, #19, #21) asserted the rule-side locator should
return the `[rule.emit]` header. **REQ-32 says the opposite** — it keys "on the
offending RULE ID, **not on its `[rule.emit]` header**". Acting on any of them
would have reversed an RDR-mandated decision and broken the shipped contract.
This is what grounding findings against the record exists to catch.

Conversely #26 was real and severe in effect: a legal `[[ rule ]]` spelling
dropped the rule from the census and collapsed a later rule's refusal locator to
`:1`, which `0024:MVV` step 2 explicitly forbids. Verified by A/B probe before
and after the fix.

## Open questions carried into kata

- `5hy1` — legs B(i)/B(iii) are low-value and could be scoped out; B(ii) must land.
- `av9m` — D8 adjudicates the fold's existence, not its precision.
- `wkqs` / `m6fz` — the hand-written-TOML-scanning shape is now evidenced twice.
  A third divergence escalates to an RDR-SEED on decoder-provided positions.

## Next

`/kata-flight --label batch:rdr-0024 --drain` ships the `type:bug` children.
