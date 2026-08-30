Model: claude-opus-5[1m]

# 3amigo iter-2 consolidation — 0025 (delta)

Delta-scoped re-run over the clauses iteration 1's resolve rewrote. Each
persona reviewed only its share of the delta, so the id-intersection hotspot
computation of iteration 1 does not apply — the scopes were disjoint by
construction. Findings are listed per persona.

## PM (delta: problem statement, MVV, decision rationale, consequences, C6/A10)

All six iteration-1 PM findings verified CLOSED against the rewrite and
grounded true against source. One new Low:

- Low | 0025:§consequences — the execution-gate bullet claimed a config-file
  successor "converts [the flag] back to one-time", but A10's composition rule
  is disjunctive (either grants, neither revokes), so the file makes a
  one-time opt-in *available*; it does not retire the flag. **fixed.**

## Implementer (delta: C1–C6, A10/A11/A12)

C1, C3, A10, A12 verified closed. Five findings, all from the iteration-1
rewrite itself:

- High | 0025:C4 + 0025:A11 — named `invokeRead`/`invokeGate`/`invokeWrite` as
  three existing sites; only `invokeRead` exists
  (`internal/accessor/executor.go:156`, also reused for read-back).
  `Executor.Gate` and `Executor.Write` handle binding errors inline, and
  several execution-failure arms have no error to `errors.As`. **fixed** — the
  hook moved to `executor.go::refusalOf`, the single refusal constructor every
  site already passes through, which reaches them all with one edit and leaves
  `Detail` empty by construction where no error exists.
- Medium | 0025:C4 — routed an `INTRASTATE_*` env collision to
  `command_output_shape`, whose C5 definition has only three output/exit-map
  arms, while C5 pinned the count at five. **fixed** — new sixth category
  `command_env_conflict`; count and mutant lists swept to six.
- Medium | 0025:C2 — called a base-dir-less separator-bearing argv0 a
  "`command_empty`-family" defect; `command_empty` means empty vector/element
  and `table.Category` has no family mechanism. **fixed** — it is an
  `execution_failure` at binding construction, not a load-time category.
- Medium | 0025:C5 — the `precedence:` line claimed first-match-in-clause-order
  without scoping it; `load.go::accessorTable` ranges a Go map (random order)
  and runs once per capability table. **fixed** — guarantee is intra-entry
  only, cross-entry explicitly unspecified and untestable by rule.
- Low | 0025:C6 — "every verb that can execute an accessor" enumerated no
  verbs. **fixed** — registration follows executor construction (set-state and
  the two flow-exec paths); `lint` explicitly excluded.

## QA (delta: scenarios 1/3/4/4b/5/5b/7, C3–C6, F7)

Every scenario-3 fixture citation verified line-by-line against FX-exit-codes;
the empty-stdout split is exact. Scenarios 1, 4, 4b, 5, 5b, 7 carry real
assertions. Two residual gaps:

- Medium | 0025:F7 + 0025:C4 `platform` — the Windows refusal had no covering
  scenario and no injection seam, so the branch is unreachable on Unix CI,
  defeating the stated reason for choosing a runtime check over a build
  constraint. **fixed** — injectable package-level `goos` var plus new
  scenario 6b. Distinct from scenario 4's refused seam: overriding `goos`
  weakens nothing at runtime, whereas a seam on the deadline triple would ship
  a way to disable the deadline.
- Low | 0025:S3 — fixture arm E1h (spawn error, no exit code) was cited but
  consumed by no case, and C3's exit maps were silent on the no-exit case.
  **fixed** — ninth arm added to S3; C3 now states a spawn failure has no exit
  code so no map can match it.
