model: claude-opus-5[1m]

# Repeatability DIFF — RDR 0026 (post-barrier, three runs)

Runs compared:
- run-1 — `claude-opus-5[1m]`, variant: full (profile: foundational)
- run-2 — `claude-sonnet-5-20250929`, variant: full (profile: foundational)
- run-3 — `claude-haiku-4-5-20251001`, variant: full (profile: foundational)

run-1's `variant:` header reads `full`, so three runs is the intended coverage, not a
`lite`-variant mismatch. The model boundary is run-3 (haiku) and secondarily run-2
(sonnet) against run-1 (opus); a split that follows that boundary is called out below.

---

## 1. Disagreements

### D-1 — `graceOn` / where the deadline is set. **3-way, and run-3 renders a mechanism C1 forbids by name.**
**Element: `0026:C1` (`precedence:` clause) — the strongest finding in this pass.**

- run-1: no `graceOn` flag exists; the parent's `SetReadDeadline` on the read ends IS the
  mark (edge (a)), and the re-arm lives inside the read loop. run-1 then still prints
  `func graceOn(outR, errR *os.File)` as a GUESS signature for a thing it just said is not
  introduced — internally inconsistent, but it gets the mechanism right.
- run-2: `graceOn(outR, errR)` is a real helper called from the parent after
  `waitBounded` fails; it "sets each read end's deadline to now + DrainGrace and lets the
  still-running drain goroutines discover that deadline on their NEXT read." run-2 then
  states, correctly, that the one-shot parent form is the rejected mechanism — so run-2's
  prose and its own pseudo-code disagree about whether `graceOn` is one-shot.
- run-3: no `graceOn` at all. Its Phase-4 pseudo-code spawns a watcher goroutine that on
  timer-fire calls `outR.SetReadDeadline(now + DrainGrace)` and `errR.SetReadDeadline(...)`
  **once, from the parent**, with no per-read re-arm anywhere in the loop it draws.

The RDR passage: C1 `precedence:` is explicit — "The deadline is re-armed at
`now + DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the IDLE GAP
between reads and never the size of the tail... A single absolute deadline would instead
bound the whole remaining tail against the clock, which turns a large buffered payload into
a false held-pipe report on a pipe with no writer at all — the failure MVV row 6 exists to
catch." The Illustrative Code repeats it ("the re-arm lives INSIDE that read loop, not in
the parent"). **Run-3 nonetheless reconstructed exactly the forbidden shape.**

The silence that let it: the Illustrative Code block *shows* `graceOn(outR, errR)` as a
parent-level call taking the two read ends, and C1 says "no separate `graceOn` flag is
introduced." Those two readings are reconcilable only by reading the sentence three
paragraphs later. A reader who takes the Illustrative Code at face value builds the
one-shot. The contract narrates the correct mechanism but never draws it.

This is a model-split finding (the weakest model built the rejected mechanism) AND the
one place where the RDR's own illustrative shape actively points at the failure its MVV
row 6 exists to catch.

### D-2 — `readBounded`'s new signature and how a drain reports out. **3-way, all GUESS.**
**Element: `0026:C1` (`precedence:`), supported by §illustrative-code.**

- run-1: `readBounded(r *os.File, limit int, grace *graceFlag) ([]byte, drainEnd)` — a new
  three-valued `drainEnd` enum (EOF / held / overflow) plus a grace flag parameter.
- run-2: `readBounded(r *os.File, limit int) (data []byte, held bool, overflow bool)` —
  two booleans, no enum, no grace parameter.
- run-3: `readBounded(r io.Reader, limit int64) ([]byte, error)` — a plain error return;
  also the only run to type the reader as `io.Reader` rather than `*os.File`, which is
  wrong on its face since the function must call `SetReadDeadline`.

RDR passage: C1 says the rewrite is "a rewrite of the function, not the signature change
alone," and the Illustrative Code says "each drain must report its terminal condition out
for `heldPipes()` to have anything to read." Three terminal conditions are Normative
(held / EOF / overflow). The **carrier** is nowhere. run-1 names this silence explicitly;
run-2 names it; run-3 does not name it as a gap at all (see D-8).

Consequence of the divergence: run-2's `(held, overflow bool)` cannot express the
precedence rule "overflow outranks held on the same drain" without the caller re-deriving
it; run-3's bare `error` cannot distinguish overflow from held at all, since overflow is
not an error condition inside the loop. Two of three reconstructions produce a type that
**cannot carry the contract's own precedence**.

### D-3 — `heldPipes()`'s input and return. **3-way.**
**Element: `0026:C1` (`precedence:`), §illustrative-code, `0026:S10`.**

- run-1: `heldPipes(outEnd, errEnd drainEnd) []string` — takes the two reported conditions.
- run-2: `heldPipes() []string` — no parameters, exactly as the Illustrative Code prints
  it; run-2 flags the return type itself as a GUESS.
- run-3: `heldPipes(drainStdoutErr, drainStderrErr)` used in the pseudo-code as a
  **boolean** (`if heldPipes(...) { ... }`), not as a slice — so run-3 cannot produce
  S10's single-pipe `Detail` form (`stderr` alone) from its own reconstruction.

RDR passage: C1 pins only the consulted shape `len(held) != 0` and the ordering guarantee
(`stdout, stderr`, fixed order; bare name in the single-pipe form per S10 + MVV row 5).
Whether the fold takes the drain conditions as parameters or reads closed-over state is
silent — and the Illustrative Code's zero-argument print is what run-2 copied.

### D-4 — `spawn`'s signature and return kind. **3-way.**
**Element: §technical-design / §approach (no contract id owns `spawn`'s signature).**

- run-1: `spawn(ctx, cmd) -> (inv, error)`; declines to print a Go signature; marks the
  parameter list a GUESS.
- run-2: `func spawn(ctx context.Context, /* GUESS: cmd *exec.Cmd or equivalent */) (invocation, error)` — value return.
- run-3: `func spawn(ctx context.Context, stdin []byte, cmd *exec.Cmd) (*invocation, error)` —
  **pointer** return, and a third `stdin []byte` parameter run-3 inferred from A7's
  `cmd.Stdin = bytes.NewReader(stdin)`.

RDR passage: A1 cites `spawn` at `cmdbind.go:205`/`:282`/`:302` and the callers' form
`inv, err := spawn(...)`, which fixes neither arity nor whether `inv` is a pointer. Nothing
in the record prints the signature. Value-vs-pointer matters here because C1's caller
obligation ("the field is readable on every error path") reads differently against a nil
pointer than against a zero value.

### D-5 — the `invocation` field set. **2-way, run-3 alone.**
**Element: `0026:C1` (`refusal:`).**

- run-1 and run-2: `stdout []byte`, `stderr []byte`, `exitCode int` — exactly the three
  spellings C1 names (`inv.stdout`, `inv.stderr`, `inv.exitCode`), both flagging the struct
  name as a GUESS.
- run-3: adds a fourth field, **`killed bool`** ("whether child was signaled"), present in
  both its type sketch and its return literal (`killed: extractSignal(werr)`).

RDR passage: C1 `refusal:` names three fields and no more. But C1 also requires the
`Detail` to carry "the direct child's exit status (exited N / **killed by signal**)", and
the record never says where that distinction is carried on `inv`. run-3 invented a field
to hold it; runs 1 and 2 folded it into a `Status string` on the error type instead
(run-1: `Status string`; run-2: `ExitStatus string`). So the silence is real: **the record
requires "exited N / killed by signal N" to reach `Detail` and never says which value
carries it across the join.**

### D-6 — `Applied()` on the write's two legs. **2-way, and it is a correctness split.**
**Element: `0026:C1` (`refusal:`), with `0026:S1` and `0026:A8` as the collision.**

- run-1: `Applied()` true on the **direct** write; the read-back leg "already today" true;
  false on read/gate; explicitly notes "one site changes, not two."
- run-2: "`Applied: bool // true on **both write legs** once a held pipe is hit post-run`"
  — run-2 flattens the two legs into one rule, which is precisely what C1 says the clause
  "does not flatten."
- run-3: "`Applied()` true on write path after child was reaped" — silent on the leg split;
  its GUESS 3 then correctly locates the change at `executor.go:367-371`.

RDR passage: C1 `refusal:` — "The write's TWO legs already differ in class and this clause
does not flatten them... So one site changes, not two." But **S1** states the expected
outcome as "on BOTH write legs `Applied()` is true AND the CLI renders the 'may have been
applied' text." Both statements are true (the read-back leg is already true today, the
direct leg gains it), yet S1's phrasing read alone yields run-2's flattening. The
disagreement is manufactured by S1's summary sentence sitting apart from C1's "one site,
not two." **The gain-vs-regression distinction lives only in S1's later sentences and in
A8; the headline of S1 contradicts the headline of C1 for a fast reader.**

### D-7 — the write-path method name and the journey cost. **2-way / 2-way.**
**Element: `0026:A1` (naming), `0026:C1` (`bound:`) (journey cost).**

- Method name: run-1 guessed `Writer.Write`; run-3 gives `(*Writer).Apply` at `:722`
  (which is what A1 actually states); run-2 does not name the caller methods at all. run-1
  marked its guess as a GUESS, so this is a widening difference, not an RDR silence — A1
  determines it and run-1 did not read A1 closely enough.
- Journey cost: run-1 reports the write journey costs up to `3·(timeout + 2·WaitDelay)`
  and that the baseline leg's refusal is DISCARDED. **Neither run-2 nor run-3 mentions the
  three-invocation journey at all** — both publish only the per-invocation
  `timeout + 2·WaitDelay`. C1 `bound:` states the 3× figure and warns "a caller sizing a
  watchdog from a single invocation's bound is killed mid-write." Two of three
  reconstructions dropped the exact fact that clause exists to prevent being dropped —
  it sits at the end of a very long `bound:` clause, behind the per-invocation figure.

### D-8 — `closeReads` as obligation vs shape. **2-way.**
**Element: `0026:C1` (`precedence:`, close-ownership sentence).**

- run-1: no explicit close is added; the existing defers satisfy it; an explicit close
  "would double-close." Correct.
- run-2: prints `closeReads(outR, errR)` in its pseudo-code with the comment "belt-and-
  suspenders over the existing defers" — i.e. exactly the double-close C1 rules out.
- run-3: omits `closeReads` and keeps only `defer: close(outR), close(errR)`.

RDR passage: C1 says it plainly, but only in the final sentence of the very long
`precedence:` clause, while the Illustrative Code shows `closeReads(outR, errR)` as a live
call on the path. Same failure shape as D-1: the drawn shape contradicts the prose that
disowns it.

### D-9 — arm ordering inside `spawn` after the join. **3-way.**
**Element: `0026:C1` (`stdin:` + `precedence:`), `0026:D-selection-predicate`.**

- run-1: werr switch (ExitError / non-ExitError) → overflow → held → nil.
- run-2: overflow → held → **werr non-ExitError** → nil. This puts the `exec.ErrWaitDelay`
  arm after the held check.
- run-3: held → success. No overflow arm and no non-ExitError arm in the pseudo-code at
  all, despite naming `exec.ErrWaitDelay` in its error-mode list.

RDR passage: C1 `stdin:` fixes that the non-ExitError arm refuses "AFTER reapGroup and the
join" and that "no arm of `spawn` returns between Wait and the join";
D-selection-predicate fixes overflow > held. Neither fixes the werr arm's position
relative to overflow and held. run-1 flagged this as a GUESS explicitly; run-2 and run-3
did not notice they had ordered it differently.

### D-10 — `DrainGrace` read as a polling interval. **run-3 alone; model-split.**
**Element: `0026:C1` (`bound:`) / `0026:D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob`.**

run-3's pseudo-code annotates both drain goroutines as "polls every 50 ms grace" and its
constants section says the grace is "re-armed before each read of final drain" — the two
statements are not the same mechanism and run-3 carries both. Runs 1 and 2 both get it
right (an offset on the read deadline, not a poll period). The RDR is not silent here —
D-the-drain-grace states it is "an offset on the single existing timer's expiry" — so this
is a comprehension failure at the weakest model rather than an RDR gap. Reported because
it is the second-order consequence of D-1: once the re-arm is hoisted out of the loop into
a parent watcher, "50 ms" has nowhere to live except as a poll interval.

### D-11 — the non-pollable (F5) fallback in the main flow. **3-way in prominence.**
**Element: `0026:F5`, `0026:S5`.**

- run-1: `pollable := probeDeadline(outR) == nil && probeDeadline(errR) == nil` as an
  explicit line-3 step, with an `if !pollable { drains.Wait() }` unbounded fallback branch
  drawn in the main pseudo-code.
- run-2: mentions the probe only in a `// ... existing pipe setup ...` comment; **its
  pseudo-code has no fallback branch**, so run-2's `spawn` would deadline a non-pollable
  pipe it cannot deadline.
- run-3: treats placement as an open GUESS ("in `readBounded` or at pipe creation?",
  "cached or panic-on-error"), and floats **panic** as a disposition — which the
  `disposition` mini-check contradicts (the row is "unbounded join — the old hang...
  recorded on the invocation at the check").

RDR passage: F5 states the check runs "in `spawn` at pipe creation" and that "on that host
the join stays unbounded." The silence run-1 named and the others fell into: **where the
`pollable` result is stored on `inv` so it can surface as the F5 `Detail` line.** Nothing
in the record says. Since S5's whole assertion is that `Detail` line, the carrier for it is
undetermined by the same silence as D-5.

### D-12 — whether the record has any silences at all. **run-3 alone; model-split.**
**Element: §technical-design (whole).**

run-3's §5 is titled "Widened Spans (if contract silences)" and answers **"None required.
The contracts are explicit on the mechanism and the precedence."** It then carries three
GUESS notes of its own and a fabricated `killed` field. run-1 widened into four spans and
listed seven unresolved silences; run-2 widened into six spans and attributed its GUESS
density to the record's deliberate "names for the shape, not the contract" stance. The
disagreement is about the record's determinacy itself, and it splits on the model
boundary: the weakest model read C1's length and confidence as completeness.

---

## 2. GUESS clusters (two or more runs marked GUESS — the candidate rewrites)

### G-1 — the drain's terminal-condition carrier. **3/3 GUESS.** → `0026:C1` `precedence:`
All three runs guessed a different type (enum / two bools / error) and all three marked it.
The record states the report must exist and what it must distinguish, and stops. This is
the single highest-value rewrite: naming the carrier resolves D-2 and D-3 together, and
without it the "overflow outranks held" precedence is unimplementable in two of three
reconstructions.

### G-2 — the held-pipe error type's identifier and field set. **3/3 GUESS.** → `0026:C1` `refusal:`
run-1 `HeldPipeError{Pipes []string; Status string; Bound time.Duration}`; run-2
`HeldPipeError{Pipes []string; ExitStatus string}` and separately a lowercase
`heldPipeError`; run-3 describes it prose-only, no struct. All three correctly place it in
`internal/accessor` and all three correctly cite the import-cycle argument — the *package*
is determinate, the *type* is not. Note the divergence on exported-vs-unexported (run-2
prints both), which matters because `executor.go` must `errors.As` it from another package.

### G-3 — `spawn`'s signature. **3/3 GUESS.** → §technical-design (no owning contract id)
See D-4. Worth hoisting into C1 or the Technical Design because the caller obligation in
C1 `refusal:` is stated in terms of a value (`inv`) whose kind is undetermined.

### G-4 — the public binding-method signatures (`Reader.Read`, `Gate.Gate`, `Writer.Apply`). **2/3 GUESS (run-1, run-3 partial); run-2 omits.** → §technical-design / `0026:A1`
run-1 names this as the reason "reconstructing the module's public API from this RDR alone
is only partially possible." A1 gives the receiver/method names and line numbers but no
parameter or return lists. Low implementation risk (the record does not change them) but
it is the reason no run could answer the prompt's question 1 completely.

### G-5 — the `Detail` string's exact composition. **2/3 GUESS (run-1 explicit, run-2 by omission of separators); run-3 asserts an order.** → `0026:C1` `refusal:`
Fixed by the record: the three parts' order (applied-sense, held-pipe reason, stderr tail)
and the pipe-name ordering. Not fixed: separators, the `exited N` / `killed by signal N`
spelling, the remediation's exact phrasing. run-3 additionally inserts "ordered per
0025:C4" as if 0025:C4 supplied the ordering, which C1 does not say.

### G-6 — the test-only stall seam's shape. **1/3 named (run-1 only).** → `0026:S7`
Below the two-run threshold, but recorded because S7 states the seam is shared by S7, S9
and MVV row 6 and is "injected through the package's existing test surface (not a build
tag, not a sleep in production code)" without saying what that surface is. run-2 and run-3
did not reach it; run-1 reached it and could not resolve it. It is a real silence with only
one witness.

### G-7 — where `pollable` is recorded on the invocation. **2/3 (run-1 explicit silence; run-3 GUESS 2).** → `0026:F5` / `0026:S5`
See D-11.

---

## 3. Agreement (confirms the RDR is determinate here — not findings)

All three runs rendered these identically: the fixed order in `spawn`
(`Wait` → `reapGroup` → bounded join); `WaitDelay = 500 ms` as the ONE bound, shared with
Cmd's stdin/post-Cancel timer, not model-declarable; `DrainGrace = 50 ms` as a mechanism
constant and **not** a second bound; the held-pipe error declared in `internal/accessor`
and never in `cmdbind`, with the import-cycle reason; `errors.As` at the executor as the
match idiom; the executor — never the binding — assigns the class, deadline-first, so
`timeout` outranks `execution_failure`; `timeout` carries no `Err`/`Detail` (0025:C4);
overflow outranks held on the same drain; the per-invocation bound `timeout + 2·WaitDelay`;
stdout on a held drain is never parsed on any path; the caller obligation to check `err`
before touching `inv` is documented, not structural, and deferred to 0028; the two drains
run concurrently; `cmdbind` has no logger, so F5 surfaces as a `Detail` line; `StdoutCap`
1 MiB with cap-exactly a value and cap+1 an overflow; `StderrTailCap` last 4 KiB;
`exec.ErrWaitDelay` on the stdin leg refuses `execution_failure` after `reapGroup` and the
join; nothing is persisted to disk.

That list is long, and it is the honest headline: C1 is a highly determining contract for
policy, ordering, precedence, placement, and class. Every finding above is about **type
shape** or about **a drawn shape contradicting the prose that disowns it** (D-1, D-8).
