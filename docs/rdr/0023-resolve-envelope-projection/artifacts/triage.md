# Triage — RDR 0023 resolve-envelope projection

Single-pass roborev triage after the Stage-8 implementation launch.

- **Window:** `BASE=d098dbd` … `HEAD=2e8227c` (frozen at Phase 0; fix
  commits `4cd4c53`, `520af02` land past it and were swept once).
- **Batch label:** `batch:rdr-0023`
- **Launch status:** `status.md` COMPLETE, worktree clean, 0 open
  needs-author-decision deviations.

## Findings × verdict

| # | Job | Anchor | Sev | Verdict | odc-type | odc-trigger | Outcome |
|---|-----|--------|-----|---------|----------|-------------|---------|
| 1 | 6462 | `flow_fixtures_0023_test.go:165` `payloadKeyOrder` reports values as keys | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | fixed pre-HEAD by `1b64a9a` (deviation D-2) |
| 2 | 6462 | `flow_mvv_0023_test.go:51` absolute model path in golden | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | D-1; golden relative + root folded |
| 3 | 6462 | `flow_mvv_0023_test.go:344` release-grammar unseeded, missing tags | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | D-5 |
| 4 | 6462 | `flow_partition_0023_test.go:668` `read-state` second role unbound | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | D-3 |
| 5 | 6462 | `flow_partition_0023_test.go:595` reader test records no executions | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | D-4; hook installed at `flow_adversarial_0023_test.go:269` |
| 6 | 6462 | `flow_projection_0023_test.go:710` HTML values only in dropped `labels` | Med | DROP `superseded-at-HEAD` | n/a-not-a-defect | — | covered by Phase 3a byte-identity probes |
| 7 | 6463 | `flow_exec.go:277` `plan-only` literal fails structural oracle | Med | DROP `superseded-at-HEAD` | environment | design-conformance | literal absent at HEAD |
| 8 | 6463 | `flow_exec.go:329` `readerExecutionHook` never installed | Low | DROP `superseded-at-HEAD` | test-oracle | test-coverage | installed + captured by Phase 3b |
| 9 | 6465 | `flow_mvv_0023_test.go:397` step-6 measures absolute path | Med | **FIX-NOW** | test-oracle | boundary | `4cd4c53` |
| 10 | 6466 | `flow_adversarial_0023_test.go:414` exit-3 case is not a reader refusal | Med | KATA-BUG | test-oracle | test-coverage | kata `0gm0` |
| 11 | 6466 | `flow_adversarial_0023_test.go:555` REQ-49 assertion is vacuous | Med | KATA-BUG | test-oracle | test-coverage | kata `0gm0` (collapsed with #10 + ADV-1) |
| 12 | 6468 | range-net: raw-only prefix fold breaks on JSON-escaped separators | Med | **FIX-NOW** | environment | backward-compat | `520af02` |
| 13 | 6469 | fix-commit sweep: `json.Marshal` HTML-escapes vs wire `SetEscapeHTML(false)` | Low | KATA-BUG | environment | boundary | appended to kata `0gm0` |

Job `6467` (`4cd4c53`) returned **No issues found**.

## Notes

**Six of thirteen findings are `n/a-not-a-defect`.** Job 6462 reviews the
Phase 1 tests-first commit, which is committed *deliberately red* by the
red-before-green gate. Its findings are the fixture defects Phase 2 then
classified as deviations D-1…D-5 and fixed before HEAD. This is the flow's
own red-green rhythm, not escaped defects, and they are excluded from the
distribution below.

**#9 was material, not cosmetic.** The step-6 pass bar measured the raw
emitted line including the absolute `--model` path. `model` is an echo
field the projection drops, so the prefix inflated the *default* side only
and the reported saving grew with the checkout path length: release-grammar
read ~50.6% on this checkout against a true 42.2%, an 8.4-point boost
against a bar the true value clears by 2.2 points. A real regression below
40% could have passed. After the fix all three shapes reproduce
`evidence/spikes/a1-byte-width.md` byte-for-byte.

**#12 is the honest cost of #9's remedy** — the fold inherited D-1's
raw-prefix approach, which is a silent no-op on any platform whose
separator is JSON-escaped. Windows is a shipped release target
(`.goreleaser.yaml`) though CI is Linux-only, so it was latent, not live.
Fixed because the remedy is small and this run introduced the second site.

**#13 was filed rather than fixed by contract.** A finding on a
fix-commit is filed, never re-fixed — the bounded sweep is what keeps
triage from becoming the refine loop it replaced.

## Defect distribution (excludes `n/a-not-a-defect`)

- `odc-type`: test-oracle=4, environment=3
- `odc-trigger`: test-coverage=3, boundary=2, design-conformance=1,
  backward-compat=1

The concentration in `test-oracle` / `test-coverage` is the signal worth
carrying forward: every genuine finding this run was about an *oracle that
did not discriminate*, not about the shipped projection, which Phase 3a
probed independently and found correct on every C1/C2 obligation.
