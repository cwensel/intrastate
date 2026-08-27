# cove spike — stray [initial] on a pure decision-table shape

Question: C2 claims `[initial]` is "already unauthorable once the owned set
is empty". Executed rather than argued.

## Arm A — baseline: pure decision-table shape, no [initial]
```
{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","findings":[{"code":"graph-dangling-edge","message":"the model declares no initial owned state; add an `[initial]` table assigning every always-present owned tag","model":"navigator","severity":"blocking","element":"model"}]}
exit status 2
```

## Arm B — same + stray [initial] over the DECLARED OBSERVED tag `status`
```
{"code":"model-invalid","message":"model does not conform to the transition-model schema","param":"model","detail":"malformed_accessor_binding: written tag status is served by 0 writers; want exactly one","findings":[{"code":"malformed_accessor_binding","message":"malformed_accessor_binding: written tag status is served by 0 writers; want exactly one","locator":"docs/rdr/0010-stateless-decision-tables/evidence/spikes/cove-dt-initial/dt-pure-initial.toml:1"}]}
exit status 2
```

Verdict: CLAIM HOLDS. The refusal arrives as `malformed_accessor_binding`
(`[initial]` keys are treated as written tags demanding exactly one writer;
a decision table declares none), not via `loadInitial`, which refuses only an
UNDECLARED key (`internal/table/load.go:546`) and defers the declared-non-owned
case to the writer binding (comment at `load.go:574-577`). C2 does not name this
category, but the prohibition is enforced. Candidate finding dismissed.
