
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

## Phase 3a — CoVe

Probes derived from the RDR and the implementation source only; the Phase 1
tests and `coverage.md` were not read. 28 spec-derived probes were run
against the real implementation — throwaway probe packages against the
`cmdbind` / `flowbind` / `accessor` seams, plus a real end-to-end MVV model
driven through the built `intrastate` binary. Probe packages were removed
after the pass; nothing permanent was added.

Three actual violations reproduced. All three were reached independently
from the spec, and each was then confirmed a second time with a fresh
minimal probe before being recorded here.

### FAIL-1 — the two exit maps claim a child that was KILLED (`0025:C3`, REQ-8/F1)

**Violates.** `0025:C3` — "the maps are consulted only for a process that
ran and exited"; a killed child never exited. `0025:F1` promises `timeout`
or `execution_failure` here, never established absence or a verdict.

**Exact failing input.**

```go
// read: an entry declaring exit_absent = [-1] (lint-green: no
// command_output_shape arm rejects a negative code)
r := cmdbind.Reader{Accessor: table.Accessor{
    Role: "repo", Command: []string{"<script that runs `sleep 30`>"},
    Keys: []string{"k"}, ExitAbsent: []int{-1}}, Name: "p",
    Config: cmdbind.Config{AllowCommands: true}}
ctx, _ := context.WithTimeout(context.Background(), 400*time.Millisecond)
values, unreadable, err := r.Read(ctx, art, []string{"k"})
```

**Observed.** `values=[{Key:k Value: Absent:true}] unreadable=[] err=<nil>`
— the killed child established every declared key ABSENT.

The gate arm is the same input with
`ExitVerdicts: map[string]string{"-1": "allow"}`: observed
`verdict="allow" reason="" err=<nil>` — execution failure laundered into a
verdict the model never decided.

**Required.** `execution_failure` (or, one layer up, `timeout`); the exit
maps must not be consulted for a process the parent killed.

**Reachability.** `spawn` sets `inv.exited = true` for any `*exec.ExitError`,
and `Wait` reports a SIGKILL as one with `ExitCode() == -1`. At the
DEADLINE the executor masks this by checking `bounded.Err()` first; it does
not mask `context.Canceled`. Probed at the executor seam with a cancelled
context and a command-backed reader whose accessor declares
`exit_absent = [-1]`:

**Observed.** `refusal=<nil> values=[{Key:k Value: Absent:true}]` — an
INTERRUPTED read returned a believed absence with no refusal at all.
`0004:C12` verifies a `<clear>` on absence, so a killed read-back command
confirms a removal that never happened. Not reachable from today's CLI
(cobra `Execute()` supplies `context.Background()`), but reachable at the
seam `C3`/`C7` govern and through any caller that cancels.

### FAIL-2 — stderr is drained only AFTER stdout completes (`0025:C4`, `0025:F4`)

**Violates.** `0025:C4` calls the deadline triple necessary AND sufficient;
`WaitDelay` bounds the stdin write and the POST-KILL drain, not the
pre-kill one. `0025:F4` bounds the accepted residue to "a child that
IGNORES termination at deadline" — this child ignores nothing and is
correct.

**Exact failing input.** A script that writes diagnostics to stderr, then a
well-formed envelope to stdout, then exits 0:

```sh
#!/bin/sh
i=0; while [ $i -lt 5000 ]; do printf 'noise-noise-noise-noise\n'; i=$((i+1)); done 1>&2
printf '{"k":"v"}'
exit 0
```

read under an 8s deadline with no `exit_absent` declared.

**Observed / threshold** (same script, stderr line count varied):

| stderr lines (~24 B each) | elapsed | result |
| --- | --- | --- |
| 1 000 (~24 KiB) | 229ms | `values=[{k v}] err=<nil>` — correct |
| 5 000 (~120 KiB) | 8.002s | `err=the read command produced no stdout and exited -1` |
| 20 000 (~480 KiB) | 8.005s | same |

**Required.** The value `v`, in milliseconds. **Observed.** The full
deadline burned, then a refusal whose text names "no stdout" for a tool
that produced a correct envelope.

**Mechanism.** `spawn` runs `readBounded(outPipe, StdoutCap+1)` to
completion and only then `readBounded(errPipe, StdoutCap)`. The child
blocks in `write(2)` on a full 64 KiB stderr pipe before it ever reaches
its stdout write; the parent blocks in `read(2)` on stdout. It deadlocks
far below the 1 MiB cap, so the cap does not bound it.

### FAIL-3 — the process group is signalled only on the FAILURE path (`0025:C4`, `0025:F4`)

**Violates.** `0025:C4` states the group signal's purpose as exactly this —
"dropping the group signal orphans the grandchild [holding the pipe]".
`0025:F4` admits process leakage only for a child outliving its `timeout`
refusal; a grandchild outliving a SUCCESSFUL invocation is admitted
nowhere.

**Exact failing input.** A wrapper of the shape every script has — one
backgrounded helper — that answers correctly and exits 0:

```sh
#!/bin/sh
( sleep 30 ) &
printf '{"k":"v"}'
exit 0
```

read under a 3s deadline.

**Observed.** `elapsed=3.001s values=[] err=context deadline exceeded` —
the tool exited 0 with a valid envelope in milliseconds and the read still
refused after the full deadline.

**Required.** `values=[{k v}]`, promptly.

**Isolated.** The same script with the grandchild's fds closed
(`( exec 1>/dev/null 2>/dev/null; sleep 30 ) &`) returns
`elapsed=244ms values=[{k v}] err=<nil>`. The stall is therefore the
inherited stdout pipe: `cmd.Cancel` fires only when the context ends, so on
the success path the group is never signalled and `readBounded`'s
`io.Copy(io.Discard, r)` drain never sees EOF while the grandchild holds
the write end.

### Clauses probed and found CONFORMING

Recorded so the pass's negative results are as legible as its findings.

- **C1** — carrier selection is total: an in-memory entry with NEITHER
  carrier builds `cmdbind.Reader`/`*cmdbind.Writer`/`cmdbind.Gate` (never a
  `Path: ""` file binding) and refuses `execution_failure` with a Detail
  naming the malformed entry, on read and gate alike.
- **C2** — declared argv → executed argv is byte-for-byte except the
  whole-element `{artifact}` substitution; a repeated placeholder
  substitutes every occurrence; a relative path, a leading-`-` path, and a
  separator-bearing argv0 with no BaseDir each refuse BEFORE spawn; a bare
  argv0 resolves through the parent `PATH` and `errors.Is(err,
  exec.ErrNotFound)` survives to the refusal.
- **C3** — verified across 20+ inputs: raw mode strips exactly one trailing
  `\n` (`draft\n\n` → `"draft\n"`); a JSON-mode omitted key is
  `unreadable`, never absent; `exit_absent` establishes
  `KeyValue{Absent:true}` IN `values`; empty stdout with exit 0 and no
  `exit_absent` refuses; a non-empty stdout always beats the exit map,
  including a `{"k":"real"}` with exit 1 listed in `exit_absent`; a gate
  deny envelope with exit 1 is a deny; a malformed gate envelope, a
  non-verdict string, and an unlisted exit each refuse rather than becoming
  a verdict; stdin is the flat object with HTML escaping off
  (`{"b":"x&y<z>","flow.stage":"<clear>"}`).
- **C4** — stdout cap is exact (1 048 576 B is a value, +1 refuses); a
  20 MB writer refuses in ~1s without deadlock; stderr tail keeps the LAST
  4 KiB (len 4096, ends with the terminal marker); the env is composed, not
  inherited (`SECRET_LEAK` absent unless named in `env_pass`; layering
  overlay > entry `env` > `env_pass` > allowlist confirmed by overriding
  `HOME`; an `env_pass` name unset in the parent is not synthesized);
  deadline elapses at 1.011s against a 1s timeout with zero survivors
  including a double-backgrounded grandchild; a non-reading child does not
  hang a 400 KB stdin write; `Detail` is empty on refusals carrying no
  error.
- **C5** — all six categories fire with their contract spellings via
  `intrastate lint` on real mutants; `Categories()` carries them at the
  TAIL in clause order; intra-entry precedence holds (conflict before
  empty, empty before placeholder); an unknown field is a document-level
  `unknown_schema_field`, not a per-entry category; `output` on a write,
  `exit_absent` on a write, `exit_verdicts` on a read, a non-verdict map
  value, and a bad `output` literal each land on `command_output_shape`;
  the `env`-chain interpreter form (`env FOO=1 /bin/bash -c`) is caught.
- **C6** — `--allow-commands` is registered exactly once, on the `flow`
  group's `PersistentFlags()`; `lint` does not carry it and validates
  ungated; all four verbs (`next`, `resolve`, `read-state`, `set-state`)
  refuse `execution_failure` naming the gate with the flag absent, and no
  child spawns (verified with a filesystem marker).
- **C7** — `Invocations()` counts Apply ENTRIES: 1 after a successful
  apply, 1 after a failed spawn, 1 after a gate-off refusal; no internal
  respawn.
- **REQ-MVV** — the full round trip runs green: a command-backed model
  lints clean; six mutants each reject with their own C5 defect; gated off
  the same invocation refuses before spawn; gated on, `flow set-state`
  drives `git config --file` through a declared wrapper and the artifact
  really changes (`[flow] stage = shipped`); read-back verifies through the
  declared `git config --get` raw-mode reader; a `<clear>` round-trips as
  ABSENCE (artifact emptied, verified via `exit_absent = [1]`); a sleeper
  refuses `flow-accessor-timeout` at 1.011s with no orphan; a failing write
  surfaces its stderr tail in `detail:`; and when the write applied but
  read-back could not complete, `flow-write-readback-incomplete` leads with
  the applied sense and appends the tail after it.

### Observation — not a FAIL

`command_shell_interpreter` keys the deny-list on the exact argv0 basename,
so `python3 -c` and `nodejs -e` are ADMITTED while `python -c` and
`node -e` are denied (`/usr/bin/python3 -c 'print(1)'` lints green). This is
NOT recorded as a violation: `0025:C5` fixes the interpreter set as OPEN —
"deny-listed, not closed … an unlisted interpreter is admitted, so the list
grows by amendment" — and an unlisted spelling is exactly that case. Noted
because the spelling gap is one an amendment would likely want to close.
