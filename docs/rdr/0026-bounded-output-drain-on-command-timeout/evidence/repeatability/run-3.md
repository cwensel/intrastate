model: claude-haiku-4-5-20251001
variant: full (profile: foundational)

## Reconstruction Run 3: Bounded Output Drain on Command Timeout

### Summary

RDR 0026 fixes the unbounded hang in `internal/cli/cmdbind/cmdbind.go::spawn` when a detached grandchild holds inherited stdout/stderr write ends. It amends 0025:C4 (deadline triple, "necessary and sufficient") to insert precedence: a bounded join outranks a whole read. The solution bounds the drain join by `WaitDelay` (500 ms), terminates held drains by setting a read deadline of `now + DrainGrace` (50 ms) on the parent's read ends, and refuses through `*accessor.ExecError` with the held pipes and tail in `Detail`.

---

## 1. Public API: Function Signatures, Types, Error Modes

### `spawn` function (package-internal, not public)
**Location**: `internal/cli/cmdbind/cmdbind.go::spawn`

```go
func spawn(
    ctx context.Context,
    stdin []byte,
    cmd *exec.Cmd,
) (*invocation, error)
```

**Return type**: `*invocation` struct
- `stdout []byte` — up to 1 MiB (StdoutCap)
- `stderr []byte` — last 4 KiB (StderrTailCap)
- `exitCode int` — direct child's exit status
- `killed bool` — whether child was signaled

**Error modes** (C1 refusal):
1. **Held pipe**: `*accessor.ExecError` with:
   - `Err`: typed held-pipe error carrying pipes held ("stdout"/"stderr"/"both") and exit status
   - `Detail`: held-pipe reason + pipes held + bound + exit status + remediation + stderr tail (ordered per 0025:C4)
   - `Applied()` true on write path after child was reaped

2. **Timeout (executor-classified)**: No `Err`/`Detail` returned from `spawn`; executor reads context's `DeadlineExceeded` first (A4, deadline-first rule)

3. **Stdin copy timeout**: `exec.ErrWaitDelay` from `cmd.Wait()` on non-reading child (A7, stdin bounded by Cmd.WaitDelay)

### Related binding types
- **Reader.Read**: calls `spawn`, checks `err` before `inv.stdout`
- **Gate.Gate**: calls `spawn`, checks `err` before `inv.exitCode`
- **Writer.Apply**: calls `spawn`, checks `err` before `inv.stdout`/read-back

### Callers (A1)
- `internal/cli/cmdbind/cmdbind.go::Reader.Read:566`
- `internal/cli/cmdbind/cmdbind.go::Gate.Gate:657`
- `internal/cli/cmdbind/cmdbind.go::(*Writer).Apply:722`

All three follow: `inv, err := spawn(...); if err != nil { return ... }` — no fast path reads held streams.

### Held-pipe error type
Declared in `internal/accessor` (not `cmdbind`, avoiding import cycle A4/C1 `refusal:`):
- Sibling to existing `*accessor.ExecError`
- Carries: pipes held (stdout/stderr), child exit status (exited N / killed by signal)
- Matched by `executor.go::Write` via `errors.As` after direct child reaped

---

## 2. Three Most Important Internal Helper Functions

### A. `readBounded(r io.Reader, limit int64) ([]byte, error)`
**Location**: `internal/cli/cmdbind/cmdbind.go:352` (rewrite required)

**Responsibility**: Drain a single pipe end (stdout or stderr) with bounded grace and overflow detection.

**Current implementation**: `io.ReadAll(io.LimitReader(r, limit)) + io.Copy(io.Discard, r)`

**Required change** (C1 `precedence:`): Rewrite to explicit loop with per-read deadline re-arm:
```
for {
    SetReadDeadline(now + DrainGrace)  // 50 ms, re-armed each iteration
    n, err := Read(buf)
    append to []byte
    if err == os.ErrDeadlineExceeded {
        return bytes, deadline error
    }
    if err == io.EOF {
        return bytes, nil
    }
}
```

**Invariants** (C1 `whole:`, `precedence:`):
- Read stays bounded at `StdoutCap+1` (1 MiB cap, cap+1 overflow failure)
- Discard-remainder loop must itself carry the deadline (grace bounds idle gap, not tail)
- Pre-timer path must stay deadline-free (no deadline set by caller)
- Overflow outranks held (drain past `StdoutCap` refuses overflow error regardless of terminal condition)
- Grace re-arm bounds the IDLE GAP between reads, never the tail size (A9)

### B. `reapGroup(proc *os.Process)`
**Location**: `internal/cli/cmdbind/cmdbind.go:297`

**Responsibility**: Release every writer the group signal can reach (non-blocking).

**Implementation** (unchanged from 0025): One `kill(2)` to `-pgid` with `SIGKILL`, build-tagged for Unix.

**Contract** (C1 `residue:`, A1):
- Kills the direct child's process group
- Non-blocking call; does not wait for processes to exit
- Escaping processes (via `setsid(2)`) survive; this is intentional residue
- Called BEFORE drain join in `spawn` (order fixed in C1)

### C. `heldPipes()` — the join's terminal condition reader (NEW site)
**Location**: `internal/cli/cmdbind/cmdbind.go` (must be created)

**Responsibility**: Determine if a drain ended because the pipe was held (deadline) or whole (EOF).

**Input**: The error returned by `readBounded`

**Predicate** (C1 `precedence:`, A2):
- If `err == os.ErrDeadlineExceeded` → held pipe
- If `err == io.EOF` → whole output (may have been partial, but pipe reached EOF)
- Any other error → propagate as execution failure

**Critical behavior** (A2 spike evidence):
- `SetReadDeadline(now + DrainGrace)` returns nil and unblocks pending `Read` in **33–101 µs**
- A read loop accumulating bytes before inspecting `err` keeps buffered data (24 bytes test: kept 24/24)
- Deadline of `now` short-circuits BEFORE the syscall → loses buffered bytes (falsified as mechanism)

---

## 3. Data Model: Persisted/Boundary Data

### Invocation struct (`*invocation`)
**Definition**: `internal/cli/cmdbind/cmdbind.go` (exact location TBD by grep)

Fields:
- `stdout []byte` — caller obligation: DO NOT READ if `err != nil` (C1 `refusal:`)
- `stderr []byte` — last 4 KiB preserved (0025:C4)
- `exitCode int` — child's wait status (0–255 for clean exit)
- `killed bool` — child was signaled (not exited cleanly)

**Handoff guarantees** (C1 `whole:`, A11 — PENDING):
- Parent MUST NOT read `inv.stdout`/`inv.stderr` until AFTER `drains.Done()` (all goroutines finished)
- Synchronization: `drains.Wait()` is the current ONLY happens-before edge (C1)
- Bounded join replaces unconditional wait → edge must be re-established (A11)
- Two edges required:
  1. Parent's `SetReadDeadline` on read ends (mark; internally synchronized via `os.File` poller) — A2
  2. Parent MUST still NOT proceed past a live drain (join remains a real join, bounds only delay before timeout, not whether to wait)
- Race-checked by S3, S7 under `-race`

### Hold-pipe error carrier
**Type**: `accessor.ExecError` (home in 0026)

**Structure**:
- `Err`: typed held-pipe error (new, declared in `internal/accessor`, C1 `refusal:`)
  - Carries: pipes held string, child exit status (exited N or killed by signal)
  - Matched by `executor.go::Write` via `errors.As`
- `Detail`: human-readable held-pipe reason
  - Ordered: applied sense + pipes held + bound (500 ms) + exit status + remediation ("close or redirect helper stdio") + stderr tail (C1 `refusal:`, ordered per 0025:C4)
  - Composed AHEAD of the tail (unlike overflow, which uses the tail as its diagnostic)

**Applied sense** (A8 — PENDING):
- `Applied() bool` on write path after direct child reaped = true
- Derived by `executor.go::Write` from held-pipe error via `errors.As` at `if err != nil` arm (C1 `refusal:`)
- Already set at six other refusal sites in `executor.go` (`:363`, `:375`, `:400`, `:409`, `:423`, `:472`)
- Flow: `cmdbind` populates `*accessor.ExecError`; `executor.go` matches and sets `applied`; `flow_exec.go` renders via `detailMayHaveApplied`

### Constants
- **WaitDelay**: 500 ms (C1 `bound:`, A5, A6) — same for stdin copy, post-Cancel wait, drain join
- **DrainGrace**: 50 ms (C1 `precedence:`, D-naming) — mechanism constant, NOT a second bound
  - Re-armed before each read of final drain (C1)
  - Bounds idle gap between reads, not tail size (C1, A9)
  - NOT model-declarable; cannot move refuse/accept boundary (D-naming)
  - Drives the final read to happen (A2 spike: `now` deadline short-circuits; `now+grace` deadline returns buffered bytes first)

---

## 4. Top-Level Pseudo-Code: Main Operation (20–40 lines)

**Function**: `spawn(ctx, stdin, cmd)`

```
--- Phase 0: Command setup (unchanged)
set cmd.Stdin to bytes.Reader(stdin)
set cmd.Stdout/Stderr to os.Pipe() read/write ends
set cmd.WaitDelay to WaitDelay (500 ms)
defer: close(outR), close(errR)

--- Phase 1: Start two concurrent drain goroutines
go drain_stdout = readBounded(outR, StdoutCap+1) — polls every 50 ms grace (PENDING A9)
go drain_stderr = readBounded(errR, StderrTailCap) — polls every 50 ms grace (PENDING A9)
var drains sync.WaitGroup
drains.Add(2)
for each goroutine: defer drains.Done()

--- Phase 2: Invoke and wait for direct child
werr := cmd.Start()
if werr != nil: return nil, werr

deadline, ok := ctx.Deadline()
werr := cmd.Wait()  --- waits for direct child and its stdin copy goroutine

--- Phase 3: Release reachable writers (direct child gone, but group may have escapees)
reapGroup(cmd.Process)  --- kill(-pgid, SIGKILL) — non-blocking
                        --- escapees survive; A1 proves no private wait elsewhere

--- Phase 4: Join drains with ONE timer of WaitDelay covering BOTH (C1, A6)
timer := time.NewTimer(WaitDelay)
go func() {
    select {
    case <-timer.C:
        --- Timer fired: set read deadline of now + DrainGrace on both read ends
        --- (A2 spike: deadline of now would lose buffered bytes; grace lets final read happen)
        outR.SetReadDeadline(now + DrainGrace)  --- measured 33–101 µs to unblock on linux/darwin (A2)
        errR.SetReadDeadline(now + DrainGrace)
        
        --- Each drain's final Read then reports pipe state, not clock:
        --- - bytes buffered: delivered with err == nil (A2)
        --- - EOF: whole output (A3)
        --- - deadline error: held pipe (C1 precedence)
        
    case <-ctx.Done():
        --- Context deadline elapsed; deadline-first rule (A4, C1 class:)
        --- Set same deadline on read ends to unblock drains
        outR.SetReadDeadline(now + DrainGrace)
        errR.SetReadDeadline(now + DrainGrace)
    }
}()

drains.Wait()  --- wait for both goroutines to return — real join, not abandoned (C1, A11 PENDING)
               --- timeout on context deadline triggers the deadline-set-on-read above

--- Phase 5: Classify and return
if heldPipes(drainStdoutErr, drainStderrErr) {
    --- One or both drains reported deadline error
    detail := format held-pipe reason (pipes, bound, exit status, tail)
    return inv, wrap(detail, errHeldPipe)  --- *accessor.ExecError per C1 refusal:
    --- Executor reads deadline first: ctx.Err() == context.DeadlineExceeded -> timeout
    --- Otherwise: errors.As(err, &HeldPipeError) -> execution_failure (A4, C1 class:)
}

--- Success: output whole, child closed pipes or drain reached EOF before bound
return &invocation{
    stdout: drainStdout,      --- caller obligation: check err first (A1)
    stderr: drainStderr,      --- last 4 KiB or complete (0025:C4)
    exitCode: extractCode(werr),
    killed: extractSignal(werr),
}, nil

--- Total bound: timeout + 2·WaitDelay (C1 bound:)
--- - cmd.Wait() leg: <1 ms to 79 ms (well inside one WaitDelay) (A6)
--- - drain join leg: 500.1–508.2 ms in 40/40 trials (A6) ← one full WaitDelay + scheduling
--- Worst measured: 2.082 s on linux (A6) against 2.5 s budget (timeout=1.5s, reserve=100ms, 2·500ms)
```

---

## 5. Widened Spans (if contract silences)

**None required.** The contracts are explicit on the mechanism and the precedence:

- **C1 `precedence:`** fully specifies the order (Wait → reapGroup → bounded join)
- **C1 `refusal:`** fully specifies error carriers and Detail structure
- **A2 spike** (verified) confirms `SetReadDeadline` unblocks and preserves buffered bytes
- **A6 spike** (verified) measures the total bound within tolerance
- **A4 source search** (verified) confirms deadline-first classification in executor

**Pending investigations** do NOT widen the contracts; they measure unexercised cases:
- **A9** (pending spike): 1 MiB tail under re-armed 50 ms deadline → whole output, no false held (MVV row 6)
- **A10** (pending spike): probe-and-clear on first pipe is a no-op for ordinary path (S3 preservation)
- **A11** (pending race test): bounded join preserves happens-before edge; S3, S7 under `-race`
- **S8** (pending test): resisting child that ignores SIGTERM → tests leg (a) under full budget

These are UNRUN tests, not gaps in API or data flow.

---

## 6. Outstanding Gaps & Decisions

**PENDING items** (status from contracts):
- A9: Re-armed deadline loop in `readBounded` — adds explicit `for { SetReadDeadline; Read }` (C1 precedence)
- A10: Pollability probe at pipe creation — `SetReadDeadline(far_future)` immediately cleared, catches `os.ErrNoDeadline` (F5)
- A11: Happens-before edge preservation — no extra flag, `drains.Done()` still runs (A11), S3/S7 under `-race`
- A8: Applied-sense derivation in `executor.go::Write` — match held-pipe error via `errors.As`, set `r.applied = true` in `if err != nil` arm

**No structural gaps.** All constants, callers, and error carriers are named or existing.

---

## Evidence Inventory

- **A2** (verified): Spike a2-unblock-read confirms deadline unblock (33–101 µs) preserves buffered bytes (darwin/linux)
- **A3** (verified): Spike a3-a6-drain-bound confirms direct child's whole output survives with only escapee left (1 KiB, 63 KiB tested)
- **A5** (verified): MVV oracles + spike a5-ingroup-margin measure group teardown margin (p99 15–19 ms, worst 27.7 ms = 5.5% of 500 ms)
- **A6** (verified): Spike a3-a6-drain-bound measures escaped-grandchild bound (total 2.013–2.082 s, worst headroom 418 ms against 2.5 s budget)
- **A1, A4, A7** (verified): Source search confirms one spawn site, deadline-first executor, stdin bounded by Cmd.WaitDelay
- **MVV** (verified): Oracles pass 75/75 under `-race -count=25`

---

## GUESS Notes

**GUESS 1**: `heldPipes()` function location. The RDR does not name where to create this decision site. Reading C1 `precedence:` and tracing the join, it must exist at the point where `spawn` reads the drain errors after `drains.Wait()`. This is likely a helper that examines the returned error values from both `readBounded` calls. **Placement guess**: immediately after `drains.Wait()` in `spawn`, before the `werr` switch at `:304`.

**GUESS 2**: Pollability probe in `readBounded` or at pipe creation? C1 names it as a check at pipe creation (F5, "checked once at creation"). The current `readBounded` is at `:352`. The probe likely sits in pipe setup (`os.Pipe` creation), not in the drain loop itself. **Placement guess**: immediately after pipe creation, before goroutines start, with the result cached or panic-on-error (S5 tests non-pollable fallback).

**GUESS 3**: Applied-sense derivation sits in `executor.go::Write` as per A8, but the exact location within the error-handling arm is not line-cited. **Placement guess**: `executor.go::Write` at `:367–371` (`if err != nil` arm), after the deadline check at `:359` but before any return, so the field assignment happens for the held-pipe error case only.

---

## Cross-Check vs RDR Structure

- **Problem Statement** (§23): unbounded hang on escapee grandchild ✓
- **Approach** (§118): decide precedence (bound wins); reuse WaitDelay ✓
- **Technical Design** (§152): spawn sequence, join with timer, read deadline, classification ✓
- **C1 normative contract** (§188): precedence, refusal, class, stdin, bound, whole, residue ✓
- **Assumptions A1–A11**: all cited ✓
- **Testing strategy S1–S10**: all cited, two (S8, A9) unrun ✓
- **Failure modes F1–F6** (§651): all cited ✓
- **MVV** (§722): 6 rows, 2 oracles unverified (A9, A10), 1 test unrun (S8) ✓

