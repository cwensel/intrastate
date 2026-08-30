Model: claude-opus-5[1m]
# iter-3 delta verification

1. FAIL — `refusalOf` exists (internal/accessor/executor.go:66) and is the constructor for 18 of 19 production refusal sites, but `Executor.selects` builds a `&Refusal{}` composite literal directly at executor.go:42 (the unknown_accessor / capability_mismatch arm). C4's "the single constructor every refusal already passes through" and A11's "the single refusal constructor" are therefore literally false by one site. Impact is presentational, not structural: line 42 mints its refusal before any Definition or Binding is selected, so no `*accessor.ExecError` can exist there and `Detail` is correctly empty — the stderr-tail mechanism still reaches every refusal that could carry one. The wording needs the qualifier ("every refusal raised after selection", or "every refusal built from a Definition"), not a different hook.

2. FAIL — the Pre-Lock Mini-Checks `trace` table, MVV-step row 2, still reads "2 lint accepts, 5 mutants rejected" while its own "Assertions in force" cell in row 1 says "C5 six defects registered". C5, MVV step 2, Testing Strategy S1, and trace row 1 all say six; this one cell says five. (A6's "four defects at Resolve" is historical and explicitly reconciled by C5's parenthetical — not a finding.)

3. PASS — `internal/table/load.go::accessorTable` is `func (l *loader) accessorTable(src map[string]sourceAcc, ...)` and its body opens `for id, a := range src`, i.e. it does range a Go map, and it is fail-fast (every arm `return nil, bad(...)` / `return nil, fail(...)`). C5's split — deterministic within an entry, unspecified across entries — matches the source. No scenario asserts a cross-entry ordering: S1's mutants are "one defect in one entry each", and the only ordering assertion in S1 is over `table.Categories()`' returned slice, which is a literal list, not the map.

4. FAIL — the check's premise does not hold, and the record's actual wording is the contradiction. C2 says a separator-bearing argv0 in a model with no source file "is a load-time `command_empty`-family defect rather than a cwd-relative resolve" — LOAD-time, not binding construction, and not `execution_failure`. That contradicts C5 twice: (a) `command_empty` is defined there as "empty vector or empty argv element", which a separator-bearing argv0 is neither, so "`command_empty`-family" names no registered category; (b) C5's six-member list has no category for it, and C5 states its list IS the closed set appended to `table.Categories()`. A load-time defect with no category cannot be reported. Either C2 must move this refusal to binding construction as `execution_failure` (which also reconciles it with C2's own "argv0 resolution is fixed at load", since the base dir is a load-recorded input the binding consumes), or C5 must register a seventh category.

5. PASS — the three production `NewExecutor` call sites are `internal/cli/flow_exec.go:256`, `internal/cli/flow_exec.go:644`, and `internal/cli/flow_state.go:306`. C6 names them exactly ("`internal/cli/flow_state.go`, `internal/cli/flow_exec.go` x2"), and S7 repeats the same triple. `accessor.NewExecutor` is the sole constructor (executor.go:22); every other occurrence is in `_test.go` files. C6's "the verbs that construct an executor (today: set-state and the two flow-exec paths)" is a correct rendering, and its stated coupling — to executor construction rather than a hand-kept verb list — is consistent with the gate living in the registry constructor.

6. PASS — S6b ("a valid command entry invoked with the binding's `goos` set to `windows`") matches C4's platform line, which specifies a runtime `runtime.GOOS` check in the command binding reading "an injectable package-level `goos` var (defaulting to `runtime.GOOS`)". No contradiction with S4: S4 refuses a deadline seam because "adding one purely to disable it would put a way to weaken the deadline into the shipping binding", and C4's platform line draws exactly that distinction ("unlike C4's deadline triple, overriding this one weakens nothing at runtime") — a `goos` override cannot disarm a safety property, it only selects which platform branch runs. S6b's justification (CI is ubuntu-latest on all six jobs, so the branch is otherwise unreachable) is the same test-reachability argument S4 declines only because the deadline seam carries a runtime cost S6b's does not. Consistent with Failure Modes' Windows bullet as well.

7. FAIL — the case list and the Expected list are both eight, but the scenario asserts nine arms. The dash-list enumerates eight cases (deny envelope + non-zero exit; malformed envelope + exit zero; unmapped non-zero exit; the same under `exit_verdicts = {"1"="deny"}`; gate exiting 2 under that map; read exiting 1 under `exit_absent = [1]`; read exiting 128; 1 MiB + 1 byte stdout) and the "Expected, in the order listed" line gives eight values, which correspond one-to-one. The E1h ninth arm (argv0 that does not resolve, `exec.ErrNotFound`, no exit code) appears only in the following prose paragraph — it is in neither the case enumeration nor the Expected list, so the "in the order listed" contract does not cover it and the arm has no positional expected value. The prose also miscounts its own grouping: it labels the empty-stdout arms "cases 3-7" (five cases) and the non-empty arms "cases 1, 2, 8" (three) = eight, then adds a ninth outside both groups. C3 did gain the required sentence ("A spawn failure has no exit code at all (`exec.ErrNotFound`, fixture E1h), so no map entry can match it however the map is written"), so the clause side is correct; only the scenario's enumeration is short.

---

## Disposition (author, same pass)

All four failed checks fixed and re-verified by grep against the record:

1. **C4 / A11 "every refusal"** — qualified to *post-selection* refusals in
   all three sites (C4 normative line, C4 prose, A11 Evidence), and
   `executor.go::selects`'s pre-selection refusal named explicitly with the
   structural reason its `Detail` is empty (no binding chosen ⇒ no invocation
   error can exist). **fixed.**
2. **Desk-trace row 2 "5 mutants"** — now six, matching C5, MVV step 2, S1 and
   trace row 1. Record-wide sweep for a surviving "five" returns clean.
   **fixed.**
3. **C2 source-less separator argv0** — was a load-time
   "`command_empty`-family" defect naming no registered category. Now an
   `execution_failure` at binding construction, with the reason stated: C5's
   registered set has no category for it, and the base dir's absence is not a
   property of the declaration the loader validates. **fixed.**
4. **S3 ninth arm** — E1h was described in prose only. The spawn-failure case
   is now in the case list and carries a ninth positional Expected value;
   both lists are length nine and correspond. **fixed.**

Note on the loop: iteration 3 surfaced no new *design* defects. All four were
mechanical — three were text replacements from the iteration-2 resolve that
failed to apply (line-wrap mismatches in a batch that aborted partway) and
were not re-verified at the time. This is edit-application error, not
verdict-flapping over a contested call; the cure was applying the intended
edits and verifying each by grep, which is what this pass did.
