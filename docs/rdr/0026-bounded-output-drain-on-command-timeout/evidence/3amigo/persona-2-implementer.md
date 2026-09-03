Model: claude-opus-5[1m]

# Persona 2 — Implementer, RDR 0026

Question asked: if I started coding this Monday, what would I ask in the first hour?

Owned starting set: `0026:C1`, the four `D-*` decisions, and the `source-anchor` edges.
**Widened** to `0026:§implementation-plan` (its four phases), `0026:A8`, `0026:F5`/`0026:S5`,
and `0026:§existing-infrastructure-audit` — because three of the five findings below are
about what the contract does NOT say, and a silence has no line range. What sent me:
C1 `refusal:` names a type (`cmdbind.HeldPipeError`) and an `errors.As` site in a package
that cannot see it; that forced me to check import direction and then the phase list that
should have scheduled the work.

Verified against source at `internal/cli/cmdbind/cmdbind.go`,
`internal/accessor/executor.go`, `internal/cli/flow_exec.go`, and `go list -deps`.

---

## S1 — HIGH. C1 `refusal:` places an `errors.As` match on a type the matching package cannot import

**Anchor**: `0026:C1` (`refusal:` clause), corroborated by `0026:A8`.

**Passage that triggered it**: C1 `refusal:` — "`Err` wraps the typed held-pipe error
`cmdbind.HeldPipeError` (pipes held, exit status), which the executor matches with
`errors.As` at the refusal site." A8 restates it: "`internal/accessor/executor.go::Write`
has, or can gain, a refusal site that matches `cmdbind.HeldPipeError` with `errors.As`".

**The problem**: the dependency runs the other way. `go list -deps ./internal/cli/cmdbind`
includes `internal/accessor`; `go list -deps ./internal/accessor` does not include
`cmdbind` (its import block is `context`, `errors`, `slices`, `time`, `internal/resolve`).
`cmdbind` already imports `accessor` to return `*accessor.ExecError`. An `errors.As` against
a concrete `cmdbind.HeldPipeError` from inside `executor.go` is an import cycle and will not
compile.

**Clarification request**: which of these does C1 intend, and where does the type live?
(a) a sentinel/interface declared in `internal/accessor` (e.g. `accessor.ErrHeldPipe` or an
`interface{ HeldPipes() []string }`) that `cmdbind` satisfies — keeps the direction legal;
(b) a structural match on `*accessor.ExecError` plus a new typed field, no new package;
(c) a third package both can import. C1's `authority` mini-check row ("Applied sense on a
write … `Refusal.applied`") does not settle it because it names the writer, not the type's
home.

**Decision it blocks**: where the held-pipe type is declared, and therefore the first file
created on Monday. Every other change in the record is downstream of this — the executor
cannot set `applied` for a held pipe until it can recognise one.

---

## S2 — HIGH. The Implementation Plan schedules no phase for the executor and CLI edits C1 and A8 require

**Anchor**: `0026:§implementation-plan` (Phases 1–4), against `0026:C1` (`refusal:`) and
`0026:A8`.

**Passage that triggered it**: Phase 1 — "Bound the drain join in `spawn` … **No other line
of `spawn` moves**." Phase 2 is comment re-citation in `cmdbind.go`. Phase 3 is fixtures.
Phase 4 is authoring-surface docs. Neither `internal/accessor/executor.go` nor
`internal/cli/flow_exec.go` appears anywhere in the plan (grep of the projected section
returns one incidental prose hit, no phase).

**The problem**: C1 `refusal:` and A8 commit two edits outside `cmdbind`:
1. `executor.go::Write` — the `if err != nil` arm at `:368` must gain the held-pipe match
   and set `r.applied = true` (today that arm returns
   `refusalOf(def, timeout, ClassExecutionFailure, err)` with no applied sense; confirmed at
   source).
2. `flow_exec.go::accessorFailureOf` — the `ClassExecutionFailure` arm at `:410-412` must
   gain the `Applied()` key. Confirmed: it carries no phase check and sets no `Detail`;
   `Applied()` has zero non-test callers repo-wide.

`0026:S1` (Testing Strategy) asserts the CLI renders "may have been applied" on the direct
write leg — an assertion with no phase that lands the code it tests.

**Clarification request**: are the executor and `flow_exec` changes a fifth phase, or are
they meant to be read as inside Phase 1? If Phase 1, its "no other line of `spawn` moves"
sentence should not be the only scoping statement, since the change is in a different
package. Also: is the `Applied()` re-key in scope for THIS record, or does JDR 0003 §D1
rule 3 land it in 0028? A8 says "JDR 0003 §D1 rule 3, which assigns that gain to this
record", and `0026:JC1` says the joint-check "fired → 0028 (home: JDR 0003 §D1)" — I read
those as consistent (constraint homed in the JDR, gain implemented here) but a Monday
implementer should not have to infer it.

**Decision it blocks**: the branch's file scope and the commit boundary — whether this is a
one-package change with a test, or a three-package change touching the CLI's error
rendering.

---

## S3 — MEDIUM. F5/S5 require a "creation-time pollability check" that does not exist and is not scheduled

**Anchor**: `0026:F5` and `0026:S5`, against `0026:§existing-infrastructure-audit`.

**Passage that triggered it**: F5 — "a pipe whose poller registration failed
(`os.ErrNoDeadline` **at the pollability check**) cannot be deadlined". S5 — "the
creation-time pollability check returns `os.ErrNoDeadline` (F5), simulated at the check …
the condition is **logged at the check**, not discovered at the join".

**The problem**: both clauses use the definite article for a surface that does not exist.
`grep -rn 'NoDeadline\|SetReadDeadline' internal/` returns nothing outside the evidence
spikes. The Existing Infrastructure Audit has a "Read-deadline grace" row marked **New**
(`DrainGrace`), but no row for a pollability check — so the audit does not schedule it
either, and neither does any phase.

**Clarification request**: three sub-questions, all first-hour:
- Where does the check run? C1 `precedence:` does not mention it; F5 says "creation-time",
  which would be right after `os.Pipe()` in `spawn` — a probe `SetReadDeadline` that is then
  cleared? Or is it lazy, at the timer-fire, checking the error `SetReadDeadline` returns?
- What does "logged" mean here? `spawn` has no logger today (its only outputs are the
  `invocation` and a typed error). Is this a stderr write, a `Detail` line on a later
  refusal, or a new logging dependency in `cmdbind`?
- Is the fallback silently unbounded, or does the check refuse? The `disposition` row in
  `0026:§conditional-mini-checks` says "unbounded join — the old hang … **loud at the
  check, by design**" — loud how, when nothing consumes it?

**Decision it blocks**: whether `spawn` gains a logging dependency, and whether S5 is a unit
test against a seam that must be extracted (a `deadliner` interface) or an assertion on a
log line. S5 says "simulated at the check", which implies an injectable seam that C1 does
not name.

---

## S4 — MEDIUM. `readBounded` must report its terminal condition, but no clause fixes WHAT it reports on the overflow path

**Anchor**: `0026:§illustrative-code` (the closing note) and the `readBounded` row of
`0026:§existing-infrastructure-audit`, against `0026:C1` (`precedence:`).

**Passage that triggered it**: the Illustrative Code note — "`readBounded` … returns its
bytes on ANY error, so today a deadline error is indistinguishable from a clean EOF at the
join … each drain must report its terminal condition out (EOF or `os.ErrDeadlineExceeded`)
for `heldPipes()` to have anything to read. That plumbing is part of this change."

**The problem**: `readBounded` has three exits, not two, and the record names two.
Current source:

```go
b, err := io.ReadAll(io.LimitReader(r, int64(limit)))
if err != nil {
    return b
}
_, _ = io.Copy(io.Discard, r)   // overflow path — this Copy can ALSO hit the deadline
return b
```

On the stdout overflow path the function has already read cap+1 bytes and then blocks in
`io.Copy(io.Discard, r)` draining the remainder. If the escapee holds the pipe, that Copy is
where the goroutine is parked when the timer fires — and it returns a deadline error that
the current code discards with `_, _ =`. So a drain can be simultaneously "overflowed" (C1
`whole:` and 0025:C4 say overflow is an `execution_failure`) and "held" (C1 `refusal:` says
held refuses with `HeldPipeError`). C1's `Selection / predicate` decision (`0026:D-selection-predicate`)
resolves EOF-vs-bound and `timeout`-vs-`execution_failure`, but not
**overflow-vs-held** — both are `execution_failure`, so the class does not disambiguate and
the two carry different `Err`/`Detail`.

**Clarification request**: when a drain exceeds `StdoutCap` AND is still held at the bound,
which refusal is returned — the existing "stdout exceeded the 1 MiB bound" (`cmdbind.go:326-329`,
which runs *after* the join and would win by position) or `HeldPipeError`? And does a
deadline hit inside the discard-Copy count as "held", given that the pipe is demonstrably
not empty and C1 defines held as "empty, with a writer the group signal could not reach"?

**Decision it blocks**: the shape of the value each drain goroutine reports out (a bare
`error`, or a small struct carrying bytes + terminal condition + overflow flag), and the
ordering of the two post-join refusal arms in `spawn`. This is the first data structure I'd
write and I cannot write it from C1.

---

## S5 — LOW. C1 `precedence:` says "after the join every read end is closed"; the audit says the existing defers already do it, and the Illustrative Code shows an explicit `closeReads`

**Anchor**: `0026:C1` (`precedence:` final sentence), `0026:§illustrative-code`,
and the close-ownership row of `0026:§existing-infrastructure-audit`.

**Passage that triggered it**: C1 — "After the join every read end is closed (deadline to
unblock, close to release the fd)". Illustrative Code — `closeReads(outR, errR)` as a step.
Audit — "`C1 precedence:` … is ALREADY satisfied by the defers; the Illustrative Code's
`closeReads` is shape, not a new obligation. The defers stay — an explicit close before them
would double-close".

**The problem**: this is resolved, but only in the audit table, and it resolves *against*
the two places an implementer reads first (the contract and the code shape). The audit's
resolution is correct — `spawn` defers both closes at pipe creation (`cmdbind.go` `defer func() { _ = outR.Close() }()`
and the `errR` pair), and they fire at function return, after the join on every path.

**Clarification request**: none blocking — but C1's "after the join every read end is
closed" is a statement about *fd release timing*, and the defers close at *function return*,
which is strictly later than "after the join" (the `werr` switch and the overflow check run
in between). If any future arm returns bytes to a caller that assumes the fd is already
released, the contract over-promises. Confirm C1 means "no read end outlives `spawn`" rather
than "closed at the join point".

**Decision it blocks**: whether to write `closeReads` at all. Low severity because the audit
answers it; flagged because an implementer following the Illustrative Code literally adds a
double-close that the audit then excuses as "a swallowed error, not a crash" — dead code
worth not writing.

---

## Checked and found sound (no finding)

- **`0026:A1`** — verified `spawn` is the only `exec.CommandContext`/`cmd.Wait()`/`drains.Wait()`
  site in the package and that all three callers check `err` before touching `inv`. One join
  does cover every verb.
- **`0026:A4`** — verified deadline-first at `invokeRead:199`, `Gate:231`, `Write:356-359`
  ahead of `:368`, and the read-back's transitive route via `raw.class == ClassTimeout`.
  C1 `class:` is safe.
- **`0026:A7`** — verified `cmd.Stdin` is a `bytes.Reader` and `cmd.WaitDelay` is set, so the
  third pipe stays Cmd-owned and `ErrWaitDelay` reaches the `default:` arm after the join.
- **`0026:D-naming`** / **`0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`**
  — the WaitDelay-reuse vs DrainGrace-offset distinction is stated with a usable criterion
  (same question vs different question). I could implement both constants from this text.
- **`0026:D-the-drains-stay-concurrent`** — pinned by S7 with a discriminating oracle
  (65536 bytes vs 1 MiB). Nothing to ask.
- **`0026:A6`** — the per-leg claim is honestly narrowed and the contract states the total,
  not the leg. S8 covers the unexercised resisting-child case.
