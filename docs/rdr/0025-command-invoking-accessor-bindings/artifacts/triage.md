# Triage — RDR 0025 command-invoking-accessor-bindings

Window: `2f49f0c`..`d845b59` (21 commits) · batch `batch:rdr-0025`
Spine: 15 per-commit auto-reviews (14 with open findings; #6529 clean).
47 findings → 8 root-cause clusters, each grounded in an isolated sub-agent.

## Verdicts

| # | Job | Cluster | Verdict | ODC type / trigger | Outcome |
|---|-----|---------|---------|--------------------|---------|
| 34, 38 | 6525, 6526 | A portability | KATA-BUG (high) | build-package-merge / design-conformance | windows build break |
| 35,36,39,40,43,44 | 6525-6527 | B drain/tail | DROP superseded-at-HEAD | timing-serialization / concurrency | fixed by 6464b51 + 76c121b |
| 47 | 6530 | C drain wait | RDR-SEED (high) | timing-serialization / concurrency | unbounded hang, live |
| 31 | 6523 | D placeholder | KATA-BUG (high) | checking / boundary | whitespace escape |
| 32 | 6523 | D interpreter | RDR-SEED | algorithm / design-conformance | wrapper-aware detection |
| 33, 41 | 6524, 6526 | E snapshot | DROP rdr-adjudicated | interface / design-conformance | 0004:C13/C14 by design |
| 37, 42 | 6525, 6526 | F json null | DROP rdr-adjudicated | checking / boundary | Phase-0 Q2 reading |
| 1–30 | 6516-6522 | G red-rhythm | DROP superseded-at-HEAD | n/a-not-a-defect | deliberately-red Phase-1 commits |
| 45 | 6528 | H cancel oracle | KATA-BUG (med/prio-high) | test-oracle / test-coverage | weak guard |
| 46 | 6528 | H leak oracle | KATA-BUG (med/prio-high) | test-oracle / timing-serialization | weak guard |

## Notes on the load-bearing calls

**C (RDR-SEED, the most serious).** `reapGroup` SIGKILLs only `-pid`, so a
grandchild that leaves the group via `setsid(2)` survives and holds the
manually-owned pipes; those pipes are outside `Cmd.WaitDelay`'s reach, so
`drains.Wait()` (cmdbind.go:302) is unbounded. **Reproduced at HEAD through
the real `Reader.Read`**: blocked >15s with the declared 2s timeout and
500ms WaitDelay long expired, no refusal delivered. C4's "If wrong" names
this exact condition ("the CLI hangs"); F4 admits only a leaked *process*,
never a lost *refusal*. Seed and not bug because the fix must reconcile two
C4 guarantees now in conflict — "drains reach a true EOF, nothing truncated"
vs "no unbounded wait". Introduced by the pipe-ownership change in 76c121b
that fixed cluster B.

Caveat for the reproduction: macOS ships no `setsid` binary, so a shell
fixture using it does NOT escape the group and the hang does not appear. The
confirming probe used a Go helper with `SysProcAttr{Setsid: true}`.

**D-31 (KATA-BUG).** C5 clause 3 skips any brace-bearing element containing
whitespace, handing it to clause 4 — but clause 4 fires only for a listed
interpreter argv0, so with any other argv0 the element is checked by nothing.
Verified by lint: `["reader","--file={artifact}"]` and `["sh","-c","cat
{artifact}"]` refuse correctly, while `["reader","prefix {artifact}"]` and
`["echo","hello {role}"]` are admitted. Contradicts C2 verbatim ("never
silently-literal text"). Distinct from D6, which decided which *category* a
form reports, not that a class escapes all six.

**G (DROP).** All 30 sit on deliberately-RED Phase-1 commits where the
implementation was declaration-only by design — the flow's own red-green
rhythm, not escapes. The 8 TEST-FIXTURE deviations were checked individually
(#12→D7, #21→D11, #23→D13, #24→D14): each strengthened its oracle rather
than weakening it for green.

**H (KATA-BUG ×2).** Both confirmed by constructing the mutant implementation
and observing the suite stay green: #45's cancellation arm passes against an
empty silent success; #46 passes with `reapGroup` entirely deleted (20.4s,
blocking until the grandchild exits naturally). Implementations are correct at
HEAD — these are guard-strength defects, so the risk is silent regression.
