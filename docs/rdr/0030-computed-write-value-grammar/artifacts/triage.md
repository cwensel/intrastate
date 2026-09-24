# Triage — RDR 0030 roborev findings

triage: docs/rdr/0030-computed-write-value-grammar.md @ b9c340f (base 48049f9) batch: batch:rdr-0030

Window: 8 commits `48049f9..b9c340f`. The spine is jobs 8278 and 8279. The range net is job 8294. Grounding is under `/tmp/batch:rdr-0030/`.

| # | Source | Location | Verdict | odc-type | odc-trigger | Outcome |
|---|---|---|---|---|---|---|
| 1 | job 8278 #1 @9fb52c2 | internal/table/step_grammar_0030_test.go:147 | DROP red-phase-commit | n/a-not-a-defect | | drop:red-phase-commit (hunk rewritten at HEAD: the control writes tc.minV) |
| 2 | job 8278 #2 @9fb52c2 | internal/cli/lint_mvv_0030_test.go:903 | DROP red-phase-commit | n/a-not-a-defect | | drop:red-phase-commit (hunk rewritten at HEAD: la#rge appended, large kept) |
| 3 | job 8279 #1 @6ace589 + job 8294 #1 (range, dup) | internal/table/step.go resolveStep | IN-SCOPE (0030:C2; 0002:C11) | checking | design-conformance | fixed:373e1de |

## Row 3 — a step onto `<clear>` bypasses the reserved-value refusal

Reproduces at HEAD b9c340f. Enum domain `["small","<clear>"]` with guard `eq = "small"` and `tier = { step = 1 }` loads. It emits row `m.escalate#small` writing `<clear>`, and resolve and flowbind treat that write as deleting the tag. The literal write `tier = "<clear>"` is refused as `reserved_tag_value`. C2 requires a stepped value to be checked "exactly as an authored literal write" is.

Fix: in `resolveStep`, after `stepValue`, refuse when the value is `ClearSentinel`, using `CatReservedTagValue` (the category a literal `<clear>` write takes). Then add a regression test with its paired negative control. The category choice is noted as the implementer's call in the grounding. A probe in a scratch worktree passed the table, graphlint and cli tests.

## Close-ledger — @373e1de

The fixup commit's review (job 8297 @373e1de) passed with no findings, so round 2 is clean. There were no `held:` rows. Jobs 8279 and 8294 were answered with `fixed:373e1de` and closed. Job 8278 was answered earlier with its drop reasons and closed. The completion gate reads COMPLETE.
