# cove spike — stray `terminal` on a pure decision-table shape

Question: does `checkDanglingEdge`'s TERMINAL arm fire independently of its
MISSING-ROOT arm? C2 depends on it firing; C5 requires the invariant silent.

## Model: zero owned tags, escape row, no write/clear, `terminal` over an observed tag
```
{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","findings":[{"code":"graph-dangling-edge","message":"the model declares no initial owned state; add an `[initial]` table assigning every always-present owned tag","model":"navigator","severity":"blocking","element":"model"},{"code":"graph-dangling-edge","message":"terminal predicate reads the observed tag \"status\"; a terminal is a predicate over owned tags, so a non-owned key resolves to no owned-state element","model":"navigator","severity":"blocking","element":"terminal[0]","key":"status","operator":"eq","literal":"Final","block":"match"}]}
exit status 2
```

Verdict: BOTH arms fire, distinguished by `element`:
  - `element: "model"`       — the missing-root arm (C5 overrides this one)
  - `element: "terminal[0]"` — the terminal-predicate arm (C2 depends on this one)

One code, two arms. C5's "declared-terminal handling ... MUST NOT emit" and
A10's wholesale `graph-dangling-edge` silence would suppress the arm C2 relies on.
See findings.md finding 1.
