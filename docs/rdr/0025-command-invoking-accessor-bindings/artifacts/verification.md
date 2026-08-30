
## Phase 3b — Adversarial

Reviewed from the RDR's Failure Modes section and the implementation
source only; the Phase 1 tests and Phase 3a's findings were not read. Three
failure modes, five tests, all five currently FAIL against the
implementation.

### ADV-1 — the two exit maps claim a KILLED child

**Mode.** `spawn` sets `inv.exited = true` whenever `cmd.Wait` returns an
`*exec.ExitError`, and a SIGKILL at the deadline is reported as one, with
`ExitCode()` == -1. `Reader.Read` and `Gate.Gate` then consult
`exit_absent` / `exit_verdicts` for that -1. So an entry declaring
`exit_absent = [-1]` or `exit_verdicts = { "-1" = "allow" }` — authorable
and lint-green, since no `command_output_shape` arm rejects a negative
code — converts a child that was KILLED WITHOUT ANSWERING into a success:
every declared key established ABSENT, or the verdict `allow`.

**Anchor.** `0025:F1` ("runtime refusals carry the 0004 diagnosis tuple …
with class `timeout` or `execution_failure`") and `0025:F2`/`0025:F3`
(execution failure is never laundered into a verdict the model did not
decide). `0025:C3` scopes the maps to a process that ran and exited — "a
spawn failure, which has no exit code at all, can never be claimed by them
however broadly they are written"; a killed process has no exit code
either.

**Why it is not already closed.** `executor.go::invokeRead` and
`Executor.Gate` reclassify when `bounded.Err()` is
`context.DeadlineExceeded`, which masks the DEADLINE arm one layer up.
That is exactly the shape `0025:C6` forbids for a safety property: "a
coincidence between two independent decisions". The mask does not cover
`context.Canceled` — an interrupted CLI — where the binding's answer
stands unreclassified. `0004:C12`'s read-back verifies a `<clear>` on
ABSENCE (`executor.go::verifyReadBack`), so a killed read-back command
confirms a removal that never happened: the `0025:F5` "confirm an
unapplied write" shape, reached through a door C1's runtime residue arm
does not cover.

**Tests.** `internal/cli/cmdbind/adversarial_0025_test.go`:
- `TestAdvReadExitAbsentMustNotClaimAKilledChild` — FAILS
  (`values=[{state.phase "" Absent:true}] err=nil`).
- `TestAdvGateExitVerdictsMustNotClaimAKilledChild` — FAILS
  (verdict `allow`, err nil).
- `TestAdvKilledChildOnParentCancelIsNotAnAnswer` — FAILS; the
  unmasked path, `context.Canceled`.

### ADV-2 — the stderr pipe is drained only AFTER stdout completes

**Mode.** `spawn` calls `readBounded(outPipe, …)` to completion and only
then `readBounded(errPipe, …)`. A child that writes more than one pipe
buffer (64 KiB) to stderr before closing stdout blocks in `write(2)` on
stderr while the parent blocks in `read(2)` on stdout. Neither moves until
the deadline. A correct, merely VERBOSE tool — one that emits diagnostics
and then a well-formed envelope — burns its whole timeout and refuses,
with the diagnosis "produced no stdout" naming a tool that produced a
correct envelope.

**Anchor.** `0025:F4` bounds the accepted residual risk to "a child that
IGNORES termination at deadline". This child ignores nothing. `0025:C4`
calls the deadline triple necessary AND sufficient; WaitDelay bounds the
stdin write and the post-kill drain, not the pre-kill one, so nothing in
the triple addresses this.

**Test.** `TestAdvVerboseStderrMustNotDeadlockTheDrain` — FAILS: 128 KiB
of stderr then `{"state.phase":"review"}` refuses after the full 3s
deadline instead of returning `review` in milliseconds.

### ADV-3 — the process group is signalled ONLY on the failure path

**Mode.** The group SIGKILL lives in `cmd.Cancel`, which the runtime
invokes only when the context ends. On the SUCCESS path the group is never
signalled, so a grandchild the tool backgrounded (`( … ) &`, the shape
every wrapper script has) survives the CLI holding the stdout/stderr pipes
it inherited. `readBounded`'s `io.Copy(io.Discard, r)` drain does not see
EOF until the last holder exits, so one backgrounded helper converts a
read whose tool exited 0 immediately into a full-deadline refusal — and
leaves the grandchild running past the CLI either way.

**Anchor.** `0025:C4` states the group signal's purpose as precisely this:
"dropping the group signal ORPHANS A GRANDCHILD holding the pipe".
`0025:F4` admits process leakage only for a child that outlives its
`timeout` refusal; a grandchild outliving a SUCCESSFUL invocation is a
mode the RDR admits nowhere.

**Test.** `TestAdvBackgroundGrandchildMustNotStallASuccessfulRead` —
FAILS: the tool prints a valid envelope and exits 0, and the read still
refuses `context deadline exceeded` after the full 3s.

### Modes probed and found ALREADY HANDLED (not counted)

- Env allowlist bypass — `childEnv` applies the `INTRASTATE_*` overlay
  LAST, so an in-memory entry `env` naming a reserved key is overwritten,
  not honoured. The C5 load arm and the runtime layering agree.
- `errors.Is(err, exec.ErrNotFound)` survives to the refusal site through
  `ExecError.Unwrap` (`0025:C4` detail).
- Stdout overflow is detected at cap+1 and refuses without stalling
  (~400ms on a 3 MiB writer).
- Empty / relative / leading-`-` artifact paths refuse in `substitute`
  before spawn (`0025:C2`).
- A deadline kill DOES reap the whole process group, grandchildren
  included (`cmd.Cancel` → `Kill(-pgid)`), so `0025:F4`'s stated residue
  is genuinely bounded on the failure path.

### Recorded finding, no edit made

`internal/probe0025/` is present and untracked in this worktree. It is
Phase 3a's concurrent artifact, outside RDR 0025's failure-mode surface
for this pass; left in place and not committed here.
