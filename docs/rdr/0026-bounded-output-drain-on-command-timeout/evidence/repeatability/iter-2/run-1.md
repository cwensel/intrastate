model: claude-haiku-4-5-20251001
variant: full (profile: foundational) — iter-2 delta re-run

---

## Reconstruction Results

### 1. Read Deadline Set and Loop Discipline

**Status**: DETERMINATE

The read deadline is re-armed PER READ inside the drain loop, never once in the parent. C1 `precedence:` explicitly states: "The deadline is re-armed at `now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP between reads and never the size of the tail."

A2 confirmed via spike that `SetReadDeadline` at `now+grace` allows buffered bytes to be delivered (first read returns whole payload with `err==nil`), while `now` or a past deadline short-circuits with `n=0` and discards buffered data.

**Pseudo-code of drain's read loop:**
```
for {
    r.SetReadDeadline(time.Now().Add(DrainGrace))
    n, err := r.Read(buf)
    if err == os.ErrDeadlineExceeded {
        // Drain ended on deadline; report held condition
        return heldCondition, n
    }
    if err == io.EOF {
        // Drain ended cleanly; report whole condition
        return wholeCondition, n
    }
    if err != nil {
        return err, n
    }
    // Accumulate and continue
    out = append(out, buf[:n]...)
    if len(out) > limit {
        return overflowCondition, len(out)
    }
}
```

---

### 2. Drain Report Type and Shape

**Status**: DETERMINATE

Each drain returns a THREE-VALUED condition report alongside its bytes. C1 `precedence:` is explicit: "The terminal condition is carried as ONE three-valued report per drain, not as a pair of booleans and not as a bare `error`: the three states are mutually exclusive and totally ordered here (overflow > held > whole)."

**The three terminal conditions in order of precedence:**
1. **overflow** — bytes read exceeded `StdoutCap` (stdout) or `StdoutCap` (stderr during read). The drain goroutine finished reading past the cap.
2. **held** — final read ended on `os.ErrDeadlineExceeded`. Pipe has a live writer (unreachable group member).
3. **whole** — final read ended on `io.EOF`. All writers exited cleanly.

**Type shape required:**
The shape must express total ordering such that overflow outranks held, which outranks whole. A two-boolean pair or a bare error cannot express this ordering; the proposed implementation carries this as a single three-valued enum or discriminated-union type. C1 `precedence:` confirms that reshaping it to booleans or an `error` "forces the caller to re-derive the precedence this clause fixes."

**heldPipes() fold input and return:**
- **Input**: the two drains' REPORTED terminal conditions (not byte counts or closed-over state).
- **Return**: held pipe NAMES in fixed order [stdout, stderr], e.g., `[]string{"stdout"}`, `[]string{"stderr"}`, `[]string{"stdout", "stderr"}`, or empty slice if no drain held. Single-pipe form is the bare name; both-pipes form is comma-separated in the fixed order (per S10).

---

### 3. Held-Pipe Fold Signature

**Status**: DETERMINATE

```go
func heldPipes(outEnd, errEnd <ThreeValuedCondition>) []string {
    // Checks ONLY the terminal conditions from each drain
    // Input: neither byte count nor any closed-over state
    // Returns: held pipe names in order [stdout, stderr] if held,
    //          empty if not held or if whole/overflow
    //
    // Overflow outranks held: a drain reporting overflow cannot
    // report held on the same drain.
    // Only a drain reporting held contributes to the result.
}
```

The fold derives nothing from byte counts or from closed-over state—it reads ONLY the reported terminal condition of each drain.

---

### 4. Explicit Close of Read Ends

**Status**: DETERMINATE

There is NO explicit close of the read ends on this path. The existing defers at pipe creation handle all closes. C1 `precedence:` states: "No read end outlives `spawn` (deadline to unblock, close to release the fd). The existing defers at pipe creation already satisfy this — they fire at function return, which is after the join but also after the `werr` switch and the overflow check, and that is the guarantee: no fd leaks, not that the fd is released at the join point. An explicit close before them would double-close."

The illustrative code explicitly draws NO `closeReads` call and explains: "the existing defers already close both read ends, so an explicit close would double-close."

---

### 5. Write Path Applied-Sense Gain

**Status**: DETERMINATE

Only the DIRECT write invocation leg GAINS the applied-sense change. C1 `refusal:` is explicit: "The DIRECT write invocation refuses `execution_failure`, and it is the leg that needs the gain: `executor.go::Write`'s `err != nil` arm sets no applied sense today, and `flow_exec.go::accessorFailureOf`'s `ClassExecutionFailure` arm carries neither a phase check nor a `Detail`, so a held-pipe write slips past the 'may have been applied' rendering. The READ-BACK leg after the write ran is already whole and gains nothing: a held pipe there returns the write's own refusal re-classed `read_back_incomplete` with `applied` already set."

A8 confirms: the read-back arm already sets `detailMayHaveApplied` and the executor already sets `applied`, so the re-key (adding the applied-sense check and `Detail` population) is one site (the direct `err != nil` write arm), not two.

---

### 6. werr Non-ExitError Arm Relative to Overflow and Held

**Status**: DETERMINATE

The `werr` non-ExitError arm (when `err == exec.ErrWaitDelay` from Cmd's stdin pipe binding a non-reading child) sits AFTER the overflow and held checks. C1 `precedence:` and D-selection-predicate agree: "When the `werr` non-ExitError arm (`exec.ErrWaitDelay`) and a drain condition both qualify, the drain conditions are selected FIRST — overflow, then held — and the non-ExitError arm refuses only if neither fired."

**Position in spawn's post-Wait path:**
```
1. werr := cmd.Wait()
2. reapGroup(cmd.Process)
3. drains.Wait() [bounded by per-read deadlines]
4. if held := heldPipes(outEnd, errEnd); len(held) != 0 { return held-pipe refusal }
5. if len(inv.stdout) > StdoutCap { return overflow refusal }
6. if werr != nil && is non-ExitError { return ErrWaitDelay refusal }
7. return inv, nil (success)
```

---

### 7. Held-Pipe Error Type Export

**Status**: DETERMINATE

The held-pipe error type is EXPORTED. C1 `refusal:` explains: "Being matched from `executor.go` and populated from `cmdbind`, both outside its own package, the type is EXPORTED and its held-pipes and exit-status fields with it — an unexported type would be unnameable at the `errors.As` site, which is the same constraint that puts the declaration in `accessor` rather than `cmdbind`. It follows the carrier already in use: `accessor.ExecError`, which `executor.go:82` already matches with `errors.As`; the held-pipe error is its sibling on the same side of the edge."

The type is declared in `internal/accessor`, not `cmdbind`, because `cmdbind` imports `accessor` and `accessor` imports no `cli` package, so a `cmdbind`-defined type named at an `errors.As` site in `executor.go` would create an import cycle.

---

### 8. Drain-Start Stall Seam Existence

**Status**: GUESS

The drain-start stall seam that S7, S9, and MVV row 6 need does NOT exist in the package today. A12 confirms: "The drain-start stall seam S7/S9/MVV row 6 share is a NEW production-code addition... S7 previously asserted the seam was 'injected through the package's existing test surface'; a repo-wide search found no such surface — `cmdbind`'s only pre-existing test-settable var is `goos` (platform gating)."

**This record MUST add it.** A12 specifies the shape: "one unexported package-level hook set by the test and nil in production (not a build tag, not a sleep in production code)."

---

## Post-Wait Pseudo-Code

```go
// After cmd.Start() completed and stdout/stderr goroutines launched:

werr := cmd.Wait()                        // Direct child reaped (≤ WaitDelay past ctx deadline)
reapGroup(cmd.Process)                    // Release group FIRST; kills escapees

// drains.Wait() is now BOUNDED: each drain's read loop re-arms
// the deadline before each read, reporting one of three conditions.
drains.Wait()

// At this point:
// - Both drain goroutines have finished writing stdout/stderr slices
// - Each drain has reported its terminal condition (overflow, held, or whole)
// - The happens-before edge is preserved: SetReadDeadline from parent
//   and drain Done() ensure the parent's slice reads do not race

inv := invocation{stdout: stdout, stderr: tail(stderr)}

switch {
case werr == nil:
    inv.exited = true
case errors.As(werr, &ee):
    if status, held := ee.Sys().(syscall.WaitStatus); held && status.Signaled() {
        inv.signaled = true
        inv.exitCode = ee.ExitCode()
        break
    }
    inv.exited = true
    inv.exitCode = ee.ExitCode()
default:
    // werr is non-ExitError (e.g., exec.ErrWaitDelay from stdin binding)
    // But check drain conditions FIRST before returning this error
    if held := heldPipes(outEnd, errEnd); len(held) != 0 {
        return inv, wrap(inv.stderr, errHeldPipe(held, exitStatus))
    }
    if len(inv.stdout) > StdoutCap {
        return inv, wrap(inv.stderr, errors.New("stdout exceeded "+strconv.Itoa(StdoutCap)+" bytes"))
    }
    // Now the non-ExitError arm, only if neither overflow nor held fired
    return inv, wrap(inv.stderr, werr)
}

// Check for held pipes (even on apparent success, a drain may have ended
// on the deadline if a group member still holds the write end)
if held := heldPipes(outEnd, errEnd); len(held) != 0 {
    return inv, wrap(tail(stderr), errHeldPipe(held, exitStatus))
}

// Overflow already detected during the bounded read in the drain goroutine,
// or checked at the join if the discard arm ended on deadline
if len(inv.stdout) > StdoutCap {
    return inv, wrap(inv.stderr, errors.New("stdout exceeded "+strconv.Itoa(StdoutCap)+" bytes"))
}

// All write ends closed, drains finished, no held pipes, no overflow
return inv, nil
```

---

## Summary

The rewritten record closes the gap from the first reconstruction. The read deadline re-arm is explicitly inside the drain loop (C1 `precedence:` line); the drain report is now a three-valued terminal condition carrying the discipline needed to avoid false held reports (C1 `precedence:`, D-selection-predicate); the held-pipe fold takes ONLY the reported conditions and returns pipe names; the write arm gains applied sense on the direct leg only; the non-ExitError arm sits after overflow/held checks; the held-pipe type is exported and in `accessor` to avoid import cycles; and the drain-start stall seam is a new hook this record must add.

---

## Appendix: Point-by-Point Reconstruction

1 | DETERMINATE | Read deadline re-armed per read inside drain loop (now + DrainGrace before each); bounded idle gap, not tail size.
2 | DETERMINATE | Three-valued terminal condition per drain (overflow > held > whole); shape must express total order.
3 | DETERMINATE | Held-pipe fold takes ONLY reported conditions, returns pipe names in order [stdout, stderr] or empty.
4 | DETERMINATE | Existing defers close read ends after join; no explicit close; double-close would error, not crash.
5 | DETERMINATE | Only DIRECT write invocation leg gains applied-sense; read-back leg already has it (C1 refusal:).
6 | DETERMINATE | Non-ExitError arm (ErrWaitDelay) sits AFTER overflow and held checks; drain conditions selected first.
7 | DETERMINATE | Held-pipe error type is EXPORTED in accessor pkg; unexported type unnameable at errors.As site.
8 | GUESS | Drain-start stall seam does NOT exist today; this record must add unexported nil-in-production hook.

