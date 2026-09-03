model: claude-sonnet-5-20250929
variant: full (profile: foundational)

# Repeatability Reconstruction — RDR 0026 (Bounded Output Drain on Command Timeout)

Source: RDR 0026 only (contracts C1, MVV, S1–S10; widened into §Technical Design
[normative-contracts, conditional-mini-checks, load-bearing-decisions, illustrative-code],
§Existing Infrastructure Audit, §Problem Statement, §Failure Modes — the contract text
itself repeatedly cross-references shape, boundary, and rationale that live in these
surrounding spans, and reconstructing signatures/types from C1 prose alone would guess
at names the Illustrative Code and Infra Audit already give explicitly). No code outside
the RDR was consulted.

---

## 1. Public API — function signatures, types, error modes

Package: `internal/cli/cmdbind` (binding layer), with one type declared in
`internal/accessor` (import-direction constraint stated explicitly in C1 `refusal:` —
`cmdbind` imports `accessor`, never the reverse).

```go
// internal/cli/cmdbind/cmdbind.go

// spawn is the single spawn site for every command-binding invocation: read, gate,
// write, and a write's read-back all funnel through it (A1). Its signature is not
// given verbatim in the RDR; reconstructed from call-site evidence (builds `inv`
// before the werr switch, returns (invocation, error) on every arm) — GUESS on the
// exact parameter list.
func spawn(ctx context.Context, /* GUESS: cmd *exec.Cmd or equivalent build args */) (invocation, error)

// readBounded is rewritten by C1 precedence: from an io.ReadAll/io.Copy pair into
// an explicit per-read deadline loop. Today's signature returns []byte and NO error
// (Infra Audit row 4; Illustrative Code final paragraph) — this is the one signature
// C1 is explicit is changing in *behavior*, and the record says a rewrite of the body,
// "not the signature change alone," implying the signature DOES change to expose the
// terminal condition. The RDR never states the new signature literally.
// GUESS at the concrete shape, grounded in what heldPipes() must be able to read:
func readBounded(r *os.File, limit int) (data []byte, held bool, overflow bool)
// Rationale for the GUESS: C1 precedence says "each drain must report its terminal
// condition out for heldPipes() to have anything to read" (Illustrative Code note) —
// EOF-whole, deadline-held, or overflow are the three reported outcomes.

// heldPipes reports, after the join, which of stdout/stderr ended on the deadline.
// Name and call shape are Illustrative Code only ("names for the shape, not the
// contract"); RDR states its return is consulted as `len(held) != 0`.
func heldPipes() []string // GUESS at concrete return type; RDR only pins usage shape

// reapGroup — existing, reused unchanged; one kill(2) to -pgid, non-blocking.
func reapGroup(p *os.Process) // signature is a GUESS; RDR names call shape only

// wrap — existing, composes Detail (held-pipe reason + stderr tail) and Err.
func wrap(tail []byte, err error) error // GUESS at param order/types; RDR shows
                                          // `wrap(tail(stderr), errHeldPipe(held, exitStatus))`
                                          // as illustrative call shape only

// errHeldPipe — constructs the typed held-pipe error.
func errHeldPipe(held []string, exitStatus /* GUESS type */) error
```

```go
// internal/accessor/model.go (or sibling file) — held-pipe error type
// C1 refusal: "declared in internal/accessor, NOT in cmdbind... follows the carrier
// already in use: accessor.ExecError (model.go:423)... the held-pipe error is its
// sibling on the same side of the edge."
// GUESS at exact field names/type name; RDR states required CONTENT, not the Go type:
type HeldPipeError struct {
    Pipes      []string // "stdout"/"stderr"/both, fixed order stdout, stderr
    ExitStatus string   // "exited N" or "killed by signal"
}
func (e *HeldPipeError) Error() string // GUESS — errors.As target, so must implement error
```

Constants (both explicitly named/reused per C1 `bound:` and Infra Audit):

```go
const WaitDelay  = 500 * time.Millisecond // existing, reused as the ONE bound —
                                            // stdin write, post-Cancel wait, AND drain join
const DrainGrace = 50 * time.Millisecond   // NEW — offset inside the one bound, not a
                                            // second knob; re-armed before EACH read
```

Error/refusal modes (from C1 `refusal:`, `class:`, disposition table):

- `timeout` — ctx `DeadlineExceeded` wins first (deadline-first rule); carries no
  `Err`/`Detail` (0025:C4's rule, restated not amended here).
- `execution_failure` — the binding's `*accessor.ExecError` / held-pipe error, matched
  via `errors.As`; `Detail` = held pipe name(s) + bound + child exit status + remediation
  text, composed AHEAD of the stderr tail; `Err` wraps the typed held-pipe error.
- Overflow (existing, unchanged) — `len(inv.stdout) > StdoutCap` outranks a held
  classification on the same drain (same class, `execution_failure`, but different
  `Detail` text — the overflow message, not the held-pipe message).
- `exec.ErrWaitDelay` (stdin leg, unchanged) — non-reading child on stdin write path,
  refused as `execution_failure` after `reapGroup` and the join (S4).
- Non-pollable pipe (`os.ErrNoDeadline`, F5) — recorded at pipe-creation-time probe,
  not at the join; surfaces as a `Detail` line on whatever refusal that invocation
  later produces; on that host the join falls back to the OLD unbounded wait.

No documented panic paths. `cmdbind` has no logger (stated explicitly under F5) — all
diagnosis is via the returned error/Detail, never a log line.

---

## 2. Three most important internal helper functions

1. **The rewritten `readBounded` per-drain loop** (unnamed in the RDR beyond its
   illustrative shape). Responsibility: read one pipe (stdout or stderr) up to its cap,
   re-arming `SetReadDeadline(now + DrainGrace)` before EACH read once grace mode is on,
   and reporting its OWN terminal condition (EOF = whole, `os.ErrDeadlineExceeded` = held,
   cap exceeded = overflow) rather than swallowing the distinction as today's code does.
   This is the one function C1 says gets rewritten, not just re-signed — replaces
   `io.ReadAll(io.LimitReader(r, limit))` + `io.Copy(io.Discard, r)` with an explicit
   `for { SetReadDeadline(...); Read(...) }` loop. Two invariants it must preserve:
   the cap boundary (`StdoutCap` is a value, `StdoutCap+1` is overflow) and a
   deadline-free read path before the timer fires.

2. **`heldPipes()`** — consulted once, immediately after the join returns. Responsibility:
   collect which of stdout/stderr (zero, one, or both) reported a deadline-held terminal
   condition, in fixed order stdout-then-stderr, so `wrap`/`errHeldPipe` can compose the
   `Detail` text and the write-path `Applied()` logic can key off it. Its exact
   implementation is illustrative-only; the RDR pins only its consulted shape
   (`len(held) != 0`) and the ordering guarantee on its output.

3. **The grace-arming step inside the join (`graceOn` in Illustrative Code, name is
   shape not contract)**. Responsibility: when the `WaitDelay` timer fires and the join
   has not yet returned, set each read end's deadline to `now + DrainGrace` and let the
   still-running drain goroutines discover that deadline on their NEXT read — this is
   what turns an in-flight blocking `Read` into one that can report "buffered bytes
   delivered, then either EOF (whole) or deadline (held)" instead of just hanging
   forever. Critically this lives INSIDE each drain's read loop (each read re-arms its
   own deadline), not as a single one-shot deadline set from the parent before
   `drains.Wait()` — the RDR is explicit that the one-shot form is a DIFFERENT, rejected
   mechanism (it would bound the whole remaining tail against the clock and falsely
   report "held" on a large buffered payload with no live writer — MVV row 6's control).

---

## 3. Data model — persisted / passed-across-boundary state

Nothing here is persisted to disk; everything is in-process state on one `spawn` call
and the `invocation` value it returns across the binding→executor boundary.

```go
// invocation — GUESS at exact field set/name (RDR calls it "inv" and cites specific
// fields by name: inv.stdout, inv.stderr, inv.exitCode); built BEFORE the werr switch
// so it is populated on every return path, error or not.
type invocation struct {
    stdout   []byte // set at cmdbind.go:275-279 today (drain goroutines' writes);
                     // MUST NOT be read by the parent until every drain has returned —
                     // this is the happens-before edge C1 preserves via edge (b)
    stderr   []byte
    exitCode int     // or exit status representation; GUESS at type
    // GUESS: no explicit "held" or "err" field on invocation itself — the contract
    // states stdout/exitCode are readable "on every error path" and the caller
    // obligation is to check the separately-returned `error` FIRST. This is stated
    // as a documented CALLER OBLIGATION, not a structural/type-level guarantee —
    // the RDR explicitly declines to make it structural (no error-only return, no
    // zeroing of inv on refusal) and defers that hardening to RDR 0028.
}
```

Cross-boundary refusal payload (accessor-side, crosses into the executor / CLI
rendering layer):

```go
// *accessor.ExecError (existing) — the held-pipe error is described as its "sibling
// on the same side of the edge," i.e. a distinct type in the same package, not a new
// field bolted onto ExecError itself. GUESS at exact relationship (sibling type vs.
// embedded) since RDR names only the constraint, not the struct diff.
type ExecError struct { /* existing, unmodified per this RDR */ }

// Held-pipe error — content is Normative (C1 refusal:), Go shape is a GUESS:
type heldPipeError struct {
    Pipes      []string // "stdout" | "stderr" | both, fixed order
    ExitStatus string   // "exited N" / "killed by signal N"
}
```

Rendered/observable surface (crosses to the CLI/agent caller, per JDR 0003 §D1 and
docs/cli-output-contract.md, both external to this RDR but cited by it):

```
Detail: "<held pipe(s)>, <bound>, <exit status>, <remediation>" + stderr tail (≤4 KiB)
Err:    wraps the typed held-pipe error (errors.As-matchable)
Class:  "timeout" | "execution_failure" | "read_back_incomplete"
Applied: bool  // true on both write legs once a held pipe is hit post-run;
                // false on read/gate invocations, always
detail field: "detailMayHaveApplied" constant (flow_exec.go:423, existing) rendered
              on both write legs when Applied() is true
```

Byte-level invariants pinned as data-model contracts (fidelity table):
- `StdoutCap` = 1 MiB — exact cap is a valid value; cap+1 is overflow.
- `StderrTailCap` = 4 KiB — LAST 4 KiB kept, earlier bytes dropped, unchanged from today.
- A drain past the ~64 KiB pipe buffer, if concurrent and prompt: whole 1 MiB recoverable.
- A drain stalled past the pipe buffer: exactly 65536 bytes recovered as a correct
  PREFIX (not the whole), remainder admitted lost.

---

## 4. Top-level pseudo-code of the main operation (~30 lines)

Reconstructed directly from the RDR's own Illustrative Code block plus the
Technical Design's four-step prose (this is the closest the RDR comes to giving
literal pseudo-code, and I have only lightly expanded it to match the surrounding
prose steps):

```go
func spawn(ctx context.Context, cmdSpec ...) (invocation, error) {
    // ... existing pipe setup: os.Pipe() for stdout/stderr, pollability probe (F5),
    //     two drain goroutines started, each doing readBounded on its pipe ...
    // defer outR.Close(); defer errR.Close()   // existing, fires at function return

    werr := cmd.Wait()          // direct child reaped (bounded by ctx deadline + WaitDelay)
    reapGroup(cmd.Process)      // release the group FIRST: one non-blocking kill(2) to -pgid;
                                 // closes the common backgrounded-grandchild case immediately

    if !waitBounded(&drains, WaitDelay) {
        // WaitDelay elapsed and the join has NOT returned: a pipe is still held
        graceOn(outR, errR)     // each drain's NEXT read re-arms deadline = now+DrainGrace
                                 // BEFORE that read — re-armed again before every
                                 // subsequent read, so grace bounds the IDLE GAP
                                 // between reads, never the size of the remaining tail
        drains.Wait()           // now returns for real: each drain has exited its loop
                                 // reporting EOF (whole), os.ErrDeadlineExceeded (held),
                                 // or cap-exceeded (overflow)
    }

    closeReads(outR, errR)      // release read-end fds on every path (belt-and-suspenders
                                 // over the existing defers)

    inv := invocation{stdout: ..., stderr: ..., exitCode: ...}  // built here, before
                                 // any error return, so it is populated on every arm

    if overflow { // checked ahead of held — overflow outranks held on the same drain
        return inv, existingOverflowError(...)   // unchanged today's overflow error
    }
    if held := heldPipes(); len(held) != 0 {
        return inv, wrap(tail(stderr), errHeldPipe(held, exitStatus))
    }
    if werr != nil && !isExitError(werr) {
        // e.g. exec.ErrWaitDelay from a non-reading child's stdin leg
        return inv, execution_failure(werr)
    }
    return inv, nil             // EOF on both drains: whole output, unchanged
}

// Downstream, at the executor (unchanged by this RDR, cited for completeness):
//   ctx.Err() == DeadlineExceeded  => class = "timeout"      (Err/Detail suppressed)
//   else if err != nil             => class = "execution_failure" (Err/Detail carried)
//   else                            => parse envelope / derive gate verdict / confirm write
```

---

## Notes on widening and GUESS density

Widened past the four selected C-kind/MVV/S elements into: §Technical Design in full
(normative-contracts prose, conditional-mini-checks tables, load-bearing-decisions,
illustrative-code), §Existing Infrastructure Audit, §Problem Statement, §Failure Modes.
Reason: C1's prose repeatedly cites concrete file:line locations, existing function
names, and an Illustrative Code block that only appear in the surrounding sections —
reconstructing "the module's public API" from C1 alone would either omit `readBounded`,
`heldPipes`, `wrap`, `reapGroup`, and the accessor-side error type entirely, or invent
them from nothing. Widening let GUESS marks attach to real named seams (documented
call shape, unstated exact Go signature) rather than to invented names.

Everything marked GUESS above is GUESS because the RDR is deliberately an
Illustrative/Normative-prose record, not a diff: it states "waitBounded/deadlineIn/
closeReads/heldPipes/errHeldPipe are names for the shape, not the contract" — so any
literal Go signature, struct field name, or exact type for these is inherently
unspecified by design, not an oversight this reconstruction can resolve further by
reading more of the record.
