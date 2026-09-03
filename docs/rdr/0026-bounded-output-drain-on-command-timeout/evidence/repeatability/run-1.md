model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability reconstruction — RDR 0026, run 1

## Read log (what I selected, what I widened past)

Selected, in order:
- `inspect --json --filter elements,outline 0026` (id list + ranges)
- `0026:C1` (alone; ~13 KB) — the sole normative contract, six clauses:
  `precedence:`, `refusal:`, `class:`, `stdin:`, `bound:`, `whole:`, `residue:`
- `0026:MVV` and `0026:S1`–`0026:S10` (batched)

WIDENED deliberately past the contract spans (each cost a read, each was necessary):
1. **`§technical-design` (152–367)** — C1 names `readBounded`, `heldPipes()`, `graceOn`,
   `wrap`, `errHeldPipe` but gives no signature for any of them. The section's
   Illustrative Code block and the `authority` mini-check table are the only places the
   drain's *reported terminal condition* is described as a value rather than a behavior
   ("`readBounded` returns `[]byte` and NO error at all"). Without this widen I would have
   GUESSED the drain-report shape and then reported that guess as an RDR silence — it is
   not one; the record states it twice.
2. **`§approach` (118–150)** — C1 never names the module. The approach names it:
   `internal/cli/cmdbind/cmdbind.go::spawn`, the single funnel.
3. **`F1`–`F6`** — C1's `class:` clause defers the `timeout`/`execution_failure` split's
   consequences to F6, and the pollability fallback (`os.ErrNoDeadline`) appears in C1
   only as a parenthetical in the technical design; F5 is where its placement (at pipe
   creation, surfaced as a `Detail` line, no logger) is determined.
4. **`§implementation-plan` (714–829)** — Phases 1 and 4 fix *where* the held-pipe error
   type is declared (`internal/accessor`), the exact executor arm that changes
   (`executor.go:368-371`), and the negative-control requirement. C1 `refusal:` asserts
   the package but not the phase boundary.

Not read (judged non-determining for items 1–4): `§context`, `§research-findings`,
`§trade-offs/§consequences`, `§alternatives-considered`, `§finalization-gate`,
`§references`. Alternatives ALT1–ALT3 and BR1–BR4 are rejected paths and by construction
do not constrain the reconstruction.

Overall: the contract set is unusually determining. C1 alone fixes ordering, the
predicate, the bound arithmetic, the refusal carrier and its package, the applied-sense
split across the write's two legs, and the concurrency antecedent. What C1 does NOT fix is
almost entirely *type shape* — signatures, field names, the drain-report value — which is
where nearly every GUESS below sits.

---

## 1. The module's public API

Module: `internal/cli/cmdbind` (package `cmdbind`), with one type declared across the
boundary in `internal/accessor`.

### 1a. The command binding surface (unchanged by this record)

The record states `spawn` is the one funnel every binding invocation passes through
(A1), and that it has three callers today: the reader, the gate, and the writer. It does
not print their signatures.

```go
// cmdbind — unchanged shapes; this record changes only what spawn does after Wait.
func (r *Reader) Read(ctx context.Context, ...) (..., error)   // GUESS: exact params/returns
func (g *Gate) Gate(ctx context.Context, ...) (..., error)     // GUESS
func (w *Writer) Write(ctx context.Context, ...) (..., error)  // GUESS
```

GUESS: the receiver names `Reader`/`Gate`/`Writer` and the method names `Read`/`Gate`/
`Write` come from the record's prose references (`Reader.parse`, `Gate.Gate`, "the
writer's read-back"); their parameter and return lists are nowhere in the RDR.

### 1b. Constants (Normative — C1 `bound:`, D-naming, D-the-drain-grace)

```go
// WaitDelay is the ONE bound. Not model-declarable. Shared, by design, with Cmd's own
// stdin-write / post-Cancel timer: two timers, one value, one question.
const WaitDelay = 500 * time.Millisecond

// DrainGrace is a mechanism constant, NOT a second bound. It is how far ahead the read
// deadline must sit for the kernel to deliver already-buffered bytes. Never
// model-declarable; moves no refuse/accept boundary.
const DrainGrace = 50 * time.Millisecond

// Unchanged from 0025:C4.
const StdoutCap    = 1 << 20 // 1 MiB
const StderrTailCap = 4 << 10 // 4 KiB
```
GUESS: `1 << 20` / `4 << 10` spellings and the exact identifier `StderrTailCap` (the
record names "StdoutCap and StderrTailCap" in C1 `bound:`, so those two identifiers are
Normative; the literals' spelling is mine).

### 1c. The held-pipe error type — declared in `internal/accessor`, NOT in `cmdbind`

Normative placement (C1 `refusal:`): `cmdbind` imports `accessor`; `accessor` imports no
`cli` package (`go list -deps ./internal/accessor` ⇒ `internal/resolve`, `internal/table`
only). A concrete `cmdbind` type named at an `errors.As` site inside `executor.go` is an
import cycle that does not compile. It is the sibling of `accessor.ExecError`
(`model.go:423`), which `executor.go:82` already matches with `errors.As`.

Normative content: it carries **the pipes held** and **the direct child's exit status**.

```go
package accessor

// HeldPipeError — GUESS on the name; the RDR calls it only "a typed held-pipe error"
// and "the held-pipe sibling", and the Illustrative Code's constructor is errHeldPipe.
type HeldPipeError struct {
    Pipes  []string // GUESS: shape. Normative: the pipe names, fixed order "stdout", "stderr".
    Status string   // GUESS: shape. Normative: "exited N" or "killed by signal".
    Bound  time.Duration // GUESS: field; Normative only that the bound is in Detail.
}

func (e *HeldPipeError) Error() string { ... } // GUESS: message text
```

GUESS: every identifier above. The RDR fixes the *package*, the *carried data*, and the
*match idiom* (`errors.As` at the refusal site in `executor.go`), and nothing else.

The `Pipes` rendering IS Normative (S10 + MVV row 5): single-pipe form is the bare name
(`stderr`); both-pipes form is the two names comma-separated in the fixed order
`stdout, stderr`.

### 1d. Error modes of a `spawn`-backed invocation

Normative, exhaustive as far as the record determines:

| Condition | Returned error | Executor class | `Applied()` |
| --- | --- | --- | --- |
| Child closes pipes, exits 0 | nil | success | n/a |
| Held pipe, ctx deadline NOT elapsed | `*accessor.ExecError` wrapping the held-pipe error | `execution_failure` | true on a **direct write**; false on read/gate |
| Held pipe, ctx deadline ALSO elapsed | same error, but discarded downstream | `timeout`, empty `Err`/`Detail` (0025:C4; F6, S6) | as above |
| Held pipe on the write's **read-back** | write's own refusal, re-classed | `read_back_incomplete`, `applied` already set | true (already today) |
| stdout past `StdoutCap` (overflow), whatever the terminal condition | existing overflow error (`cmdbind.go:326-329`, unchanged) | `execution_failure` | — |
| Non-reading child holding stdin's read end | `exec.ErrWaitDelay` from `Wait`, via the non-`ExitError` arm | `execution_failure` | — |
| Non-pollable read end (`os.ErrNoDeadline` at the creation probe) | no new error — join stays UNBOUNDED (the old hang) | whatever the old path produced | — |

**Precedence among these is Normative and total:**
`timeout` > `execution_failure` (the executor's existing deadline-first rule, A4) ·
overflow > held **on the same drain** (D-selection-predicate: the bytes were read, so the
pipe's terminal state is not what the invocation turns on) · EOF > bound if EOF arrives
first.

The binding **never classifies `timeout` itself** (C1 `class:`). `spawn` returns errors;
the executor assigns classes.

**Caller obligation (Normative, and explicitly NOT structural):** the returned
invocation's `stdout`/`exitCode` fields are populated on every error path — `spawn`
builds `inv` before the `werr` switch (`cmdbind.go:304`) and both error arms return that
populated value (`:322`, `:326`). Every caller MUST check `err` before touching `inv`.
The record ships this as a *comment* at the `inv` construction site naming the two error
returns, and explicitly declines to make it structural (an error-only return, or zeroing
`inv`); it routes that to 0028's command-reader read-back leg.

---

## 2. The three most important internal helper functions

The Illustrative Code names `waitBounded`, `graceOn`, `closeReads`, `heldPipes`,
`errHeldPipe` and states in terms: *"names for the shape, not the contract."* So the
names below are GUESSES; the **responsibilities** are Normative.

### H1. `readBounded` — REWRITTEN, and the only rewrite in the record

```go
// Today (cmdbind.go:352-364): io.ReadAll(io.LimitReader(r, limit)) then
// io.Copy(io.Discard, r); returns []byte and NO error.
// After: an explicit read loop that re-arms per read and reports its terminal condition.
func readBounded(r *os.File, limit int, grace *graceFlag) ([]byte, drainEnd) // GUESS: signature
```

Responsibility (Normative): own the read loop so there is a per-read point to re-arm the
deadline, and **report its terminal condition out** — held / EOF / overflow. Today it
reports nothing, so held is indistinguishable from clean EOF at the join; the record calls
this out twice as the fact the current code does not expose.

Three preserved behaviors, pinned by S3 and called a contract violation rather than a
refactor if broken:
- the read stays bounded at `StdoutCap+1`, so exactly-at-cap is a value and cap+1 is a
  failure (the check at `cmdbind.go:326` reads `len(inv.stdout) > StdoutCap`);
- the **discard-remainder exit must itself carry the deadline** (the third exit,
  `io.Copy(io.Discard, r)`, whose error is discarded today);
- the **pre-timer path stays deadline-free**.

The re-arm is `SetReadDeadline(now + DrainGrace)` **before each read**, so the grace
bounds the IDLE GAP between reads and never the size of the tail. A single absolute
deadline is forbidden by name: it would bound the whole remaining tail against the clock
and turn a large buffered payload into a false held report on a writerless pipe — the
failure MVV row 6 exists to catch.

### H2. `graceOn` / the deadline-set — the synchronisation mark

```go
func graceOn(outR, errR *os.File) // GUESS: name and signature
```

Responsibility (Normative): after the bounded join expires, mark the drains so each
subsequent read arms its own `now + DrainGrace` deadline. The parent's `SetReadDeadline`
on the read ends **is** happens-before edge (a); `os.File`'s poller is internally
synchronised, so `SetReadDeadline` from the joining goroutine against a `Read` in flight
is the sanctioned cross-goroutine call, and **no separate `graceOn` flag is introduced**.

Explicitly forbidden: a one-shot `SetReadDeadline` from the parent before `drains.Wait()`.
The re-arm lives INSIDE the read loop, not in the parent.

The grace is what MAKES the final read happen: Go checks the deadline before attempting
the syscall (`internal/poll/fd_unix.go::(*FD).Read` calls `prepareRead` ahead of
`syscall.Read`), so a deadline of `now` short-circuits with n=0 and reports held without
ever looking at the pipe.

### H3. `heldPipes` — the refuse/accept predicate

```go
func heldPipes(outEnd, errEnd drainEnd) []string // GUESS: signature
```

Responsibility (Normative): fold the two drains' reported terminal conditions into the
refusal decision. **The predicate is the drain's REPORTED terminal condition, not the byte
count and not the clock.** A drain ending on `os.ErrDeadlineExceeded` reports held however
many bytes it delivered during the grace; one ending on EOF reports whole. "Empty, with a
writer the group signal could not reach" describes the common case but is expressly *not*
the test — a writer trickling bytes slower than the grace is never empty and must refuse
(S9's companion). Overflow outranks held on the same drain.

Honorable mentions the record fixes but which are not among the top three:
`waitBounded(&drains, WaitDelay)` (one timer covering BOTH drains — one bound, not two);
`wrap`/`withStderrTail` (composes `Detail`: held-pipe reason AHEAD of the stderr tail, in
one slot, applied-sense first); the **pollability probe** at pipe creation (probe
`SetReadDeadline` on each read end, immediately cleared, error is the signal — NEW, and
nothing in `internal/` calls `SetReadDeadline` or names `os.ErrNoDeadline` today).

`closeReads` is explicitly **shape, not an obligation**: the existing defers at pipe
creation already satisfy "no read end outlives `spawn`", and an explicit close before them
would double-close, so none is added.

---

## 3. Data model across the boundary

### 3a. The invocation value (`inv`) — internal to `cmdbind`, returned to the executor

```go
type invocation struct { // GUESS: name and every field name
    stdout   []byte // Normative: populated on every error path; NEVER parsed on a held drain
    stderr   []byte // Normative: LAST 4 KiB preserved
    exitCode int    // Normative: readable on error paths (hence the caller obligation)
}
```
Normative: the fields `inv.stdout`, `inv.stderr`, `inv.exitCode` are named in C1 by those
spellings; the struct name and its declaration site are GUESSES.

### 3b. The drain's reported terminal condition — NEW, and the record's own admission of absence

```go
type drainEnd int // GUESS: entire type
const (
    drainEOF drainEnd = iota // whole output
    drainHeld                // ended on os.ErrDeadlineExceeded
    drainOverflow            // past StdoutCap; outranks held
)
```
Normative: three and only three terminal conditions exist (held / EOF / overflow), they
are reported *out* of the drain goroutine, and `heldPipes()` reads them. Everything about
the encoding is a GUESS — the record says explicitly the distinction "is not a fact the
current code exposes."

### 3c. The refusal across the accessor→CLI boundary

```go
// accessor side
type ExecError struct { ... }            // existing, model.go:423
type Refusal struct {                     // GUESS: name
    applied bool                          // Normative field-sense: model.go:324's Applied()
    Detail  string
    Err     error
    class   string
}
func (r *Refusal) Applied() bool // Normative: model.go:324; ZERO non-test callers today
```

**`Detail` composition (Normative, ordered):**
1. applied-sense first,
2. the held-pipe reason — pipe(s) held (`stdout` / `stderr` / `stdout, stderr`), the bound,
   the direct child's exit status (`exited N` / `killed by signal N`), and the remediation
   text *"close or redirect the helper's inherited stdio"*,
3. then the stderr tail collected up to the bound.
Plus, on any invocation whose pollability probe failed, an F5 `Detail` line naming the
host condition. `cmdbind` has no logger and gains none, so F5 is a `Detail` line, never a
log line.

**Class vocabulary crossing to the CLI envelope:** `timeout`, `execution_failure`,
`read_back_incomplete` (and the envelope's `detail` field carrying
`flow_exec.go:423`'s existing `detailMayHaveApplied` constant, per
`docs/cli-output-contract.md`).

**Persisted:** nothing. Nothing in this record is written to disk. The only durable
artifact is the doc text added to `docs/cli-output-contract.md` in Phase 5 (the residue
and remediation, beside the command-entry `timeout` text).

### 3d. Bound arithmetic published across the boundary

- Per invocation, the COMMITTED return bound is `timeout + 2·WaitDelay`, asserted with a
  scheduling tolerance of **+100 ms**. Inside the tolerance is a pass; past it is a
  contract violation.
- The two legs are NOT separately bounded. The drain-join leg runs one full `WaitDelay`
  plus timer-fire and goroutine-scheduling latency (measured 500.1–508.2 ms, 40/40); the
  deadline leg's unused allowance absorbs it. That absorption is *why* the total holds and
  is stated rather than left as a coincidence.
- A command **write** runs up to THREE invocations through the same reader: the pre-write
  baseline read when `protectedKeys` returns non-empty (`executor.go:320`), the write
  (`:355`), and the read-back (`:397`) — so a write journey costs up to
  `3·(timeout + 2·WaitDelay)`. The baseline leg's refusal is DISCARDED (`:326` sets
  `baselineUnread` and continues); this record does not change that, but names the cost so
  the published total is not read as the journey's.
- Moving `WaitDelay` is a **contract change requiring an amendment**, not a one-line
  follow-up, because C1 `bound:` publishes it. S4 and S8 are the pinning pair (S4 fails if
  it grows past what the stdin path tolerates; S8 fails if it shrinks past what a
  resisting child needs).

---

## 4. Top-level pseudo-code of the main operation

`internal/cli/cmdbind/cmdbind.go::spawn` — the one funnel. Ordering is Normative and
fixed; no arm of `spawn` returns between `Wait` and the join.

```
spawn(ctx, cmd) -> (inv, error):
 1  outR, outW := os.Pipe(); errR, errW := os.Pipe()
 2  defer outR.Close(); defer errR.Close()          # existing defers; NO explicit close added
 3  pollable := probeDeadline(outR) == nil && probeDeadline(errR) == nil   # NEW, F5
 4      # probe SetReadDeadline, immediately cleared; os.ErrNoDeadline => non-pollable.
 5      # Recorded on inv HERE, surfaced as a Detail line on any refusal this invocation makes.
 6  cmd.Stdout, cmd.Stderr = outW, errW              # os.File ends, so the group is
 7  cmd.Stdin = bytes.Reader(payload)                # releasable before the drains are joined
 8  cmd.WaitDelay = WaitDelay                        # Cmd's own stdin / post-Cancel bound
 9  cmd.Cancel = one signal, returns immediately     # so Cmd's timer starts at the deadline
10  start cmd; start TWO CONCURRENT drain goroutines:  # concurrency is load-bearing (A3):
11      stdout: b, end := readBounded(outR, StdoutCap+1); record; drains.Done()
12      stderr: b, end := readBounded(errR, StderrTailCap); record; drains.Done()
13  werr := cmd.Wait()                               # direct child reaped, <= WaitDelay past deadline
14  reapGroup(cmd.Process)                           # ONE kill(2) to -pgid, non-blocking; RELEASE FIRST
15  if !pollable:                                    # F5 fallback: the old unbounded hang,
16      drains.Wait()                                # recorded at the check, not discovered here
17  else if !waitBounded(&drains, WaitDelay):        # ONE timer, covering BOTH drains
18      graceOn(outR, errR)                          # edge (a): the deadline-set IS the mark.
19          # Each drain now arms SetReadDeadline(now+DrainGrace) BEFORE EACH read, so the
20          # grace bounds the IDLE GAP, never the tail. A drain making progress runs to EOF
21          # however many reads it takes; only a gap > grace ends it as held.
22      drains.Wait()                                # edge (b): a REAL join. A drain ending on
23          # the deadline still runs Done(). Proceeding past a live drain is FORBIDDEN by name.
24  # -- only past here may the parent read inv.stdout / inv.stderr --
25  inv := invocation{stdout, stderr, exitCode}      # built here; readable on EVERY error path,
26      # so the caller obligation "check err before touching inv" is carried as a COMMENT here.
27  switch:
28    werr is *ExitError:                 exitCode = its code; fall through
29    werr non-nil, not ExitError:        return inv, wrap(tail(stderr), werr)   # e.g. exec.ErrWaitDelay
30  if len(inv.stdout) > StdoutCap:       return inv, overflowErr   # OVERFLOW OUTRANKS HELD
31  if held := heldPipes(); len(held) != 0:
32      return inv, wrap(detail(held, bound, exitStatus, remediation), tail(stderr),
33                       &accessor.HeldPipeError{Pipes: held, Status: exitStatus})
34  return inv, nil                                  # EOF before OR after the timer: whole output
35
36 # Downstream, unchanged in location: the EXECUTOR classifies.
37 #   ctx.Err() == DeadlineExceeded            -> timeout (empty Err/Detail per 0025:C4) [F6]
38 #   else errors.As(err, &HeldPipeError{})    -> execution_failure; on the DIRECT write arm
39 #     (executor.go:368-371) ALSO set the applied sense -- ONLY for the held-pipe sibling,
40 #     never for *accessor.ExecError generally (the required negative control).
```

GUESS: lines 1–12 (setup order, `os.Pipe` call sites, `cmd.WaitDelay` assignment, the
`bytes.Reader` construction) reconstruct today's code, which the RDR describes but does
not print. Line 30's placement of the overflow check before the held check follows
D-selection-predicate's "overflow wins"; the RDR cites the existing check at
`cmdbind.go:326-329` as standing unchanged *after the join*, so its position relative to
line 31 is a GUESS consistent with the stated precedence.

GUESS: the `switch` shape at 27–29 and the `wrap(...)` argument order.

---

## Explicit RDR silences found (things I had to guess even after widening)

1. **Every identifier** of the new type and helpers: `HeldPipeError`, `drainEnd`,
   `readBounded`'s new signature, `graceOn`, `heldPipes`, `waitBounded`. The record says
   in terms these are "shape, not the contract." This is a *deliberate* silence, not a
   defect — but it means two implementers will not produce the same signatures.
2. **How a drain reports out.** Return value vs shared struct vs channel: undetermined.
   The record establishes only that the report must exist and what it must distinguish.
3. **The `Detail` string's exact wording/formatting.** Fixed: the pipe-name ordering
   (`stdout, stderr`) and the ordering of the three parts. Not fixed: separators, the
   remediation's exact phrasing beyond "close or redirect the helper's inherited stdio",
   the `exited N` / `killed by signal` spelling.
4. **The public binding method signatures** (`Reader.Read`, `Gate.Gate`, `Writer.Write`).
   Nowhere in the record. Reconstructing the *module's public API* from this RDR alone is
   therefore only partially possible — which is itself worth reporting, since C1 is
   otherwise the most determining contract in the set.
5. **The test-only stall seam's shape.** S7 fixes that there is ONE seam shared by S7, S9
   and MVV row 6, injected "through the package's existing test surface (not a build tag,
   not a sleep in production code)" — but not what that surface is.
6. **The pollability probe's "far-future deadline" value.** `§technical-design` says a
   far-future deadline, F5 says a probe immediately cleared. Undetermined.
7. **Where `pollable` is stored on `inv`** so it can surface as an F5 `Detail` line.

## Silences I did NOT report because widening resolved them
- Whether `heldPipes` turns on byte count or the clock (C1 `precedence:` — neither; the
  reported condition).
- Where the held-pipe type lives (C1 `refusal:` — `internal/accessor`, with the
  import-cycle argument given).
- Whether one or two executor sites change (C1 `refusal:` — one; the read-back leg is
  already correct and is named only to stop a reader re-keying a correct arm).
- Whether `DrainGrace` is a second bound (C1 `bound:` — no; a mechanism constant).
- Whether the parent may proceed past a live drain (C1 `precedence:` — forbidden by name).
