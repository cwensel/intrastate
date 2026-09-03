Model: claude-opus-5[1m]

# Persona 3 — QA / Tester, RDR 0026

Read: `0026:C1`, `0026:MVV`, `0026:S1`–`0026:S9` (owned). WIDENED to
`0026:§illustrative-code`, `0026:§conditional-mini-checks`, `0026:§existing-infrastructure-audit`,
`0026:§technical-design`, `0026:§failure-modes` and `0026:§consequences`. What sent me:
three of my scenarios assert against a mechanism the scenario text never states (the drain
loop's behaviour after the deadline is armed), and one asserts a CLI-rendered string that
exists only in the executor layer — neither pass criterion is recoverable inside the
S/C/MVV set.

Overall: this record is unusually testable. Most rows carry a named failing control
(`oracle` mini-check), a measured basis, and a discriminating negative. The findings below
are the residue — rows where I can build the fixture but cannot decide, from the record,
whether the run passed.

---

## HIGH

### H1 — `0026:S9` (companion row) — "held" is undefined for a pipe that is non-empty *and* deadlines

`C1 precedence:` defines a held pipe twice, and the two definitions disagree on exactly the
input S9's companion feeds:

- as a *condition*: "a deadline error is a HELD pipe — **empty**, with a writer the group
  signal could not reach";
- as a *mechanism* (`0026:§illustrative-code`, `heldPipes()`): held is whatever a drain
  **reports** as its terminal condition, and the terminal condition is
  `os.ErrDeadlineExceeded`.

S9's companion is "an escaped writer emitting a byte every 10 ms forever", expected REFUSED.
Under a `DrainGrace` of 50 ms that pipe is **not empty** during the grace window — the drain
loop takes roughly five successful reads and only then hits the absolute deadline. So the
drain terminates on `ErrDeadlineExceeded` with bytes in hand. Under the mechanism reading
that is held (refuse — matches S9's expectation); under the condition reading it is not
("empty" is false), and a strict implementation of the prose could accept.

Missing pass/fail criterion: whether `heldPipes()` keys on the terminal error alone, or on
terminal error AND zero bytes read after the deadline was armed.

Test it prevents: the S9 companion assertion itself. I cannot write the trickle-writer arm,
because a refusal and an acceptance are each defensible against a different sentence of the
same clause — and this arm is the one S9 says "pins refuse-on-no-EOF, not on elapsed time".
The prompt-drain arm of S9 is unaffected and writable.

**Fix shape**: one clause in `C1 precedence:` — held is the drain's reported terminal
condition (`os.ErrDeadlineExceeded`), regardless of bytes delivered during the grace; "empty"
is a description of the common case, not the predicate.

### H2 — `0026:C1` `precedence:` / `0026:MVV` row 6 — the deadline is armed once, but the drain loop reads many times

`C1` and MVV row 6 both speak of "each drain's **FINAL** read" (singular). `readBounded` is a
loop (`0026:§existing-infrastructure-audit` cites `cmdbind.go:352-356`). With one absolute
deadline at `now + DrainGrace`, a large buffered payload takes many reads: MVV row 6 is 1 MiB
against a ~64 KiB pipe buffer, so at minimum 16 pipe-buffer-fills must transit, each gated by
the *writer* refilling — and the writer here is a child that already exited, so the bytes are
whatever the pipe holds plus whatever the OS pipe can stage. Whether 1 MiB drains inside
50 ms is a *measurement*, and the record has none for this path: A2 measured 24 bytes, A3
measured a delayed drain (unbounded grace) at 1 KiB / 63 KiB / 1 MiB, never a 1 MiB drain
under a 50 ms absolute deadline.

Missing pass/fail criterion: (a) whether the deadline is re-armed per read or set once for the
whole tail, and (b) if once, what payload size the grace is committed to draining. If neither
is stated, MVV row 6 has no deterministic pass — it may pass on a fast machine and refuse on a
loaded one, and it would refuse by reporting a *held pipe* on a pipe with no writer at all,
which is a false refusal, not a slow one.

Test it prevents: MVV row 6 (late-but-unheld, 1 MiB) — the row the record calls "the MVV's
guard on C1 `precedence:`, not a formality". Also S9's prompt arm at 1 MiB.

Note this is the same clock-sensitivity `F2` accepts for `WaitDelay`, but `F2` was closed by
A5's 160-trial margin measurement; `DrainGrace` has no equivalent, and 50 ms was sized against
*scheduling latency* (~28 ms, A5), not against *drain throughput*.

**Fix shape**: state the re-arm rule in `C1 precedence:` (per-read `now + DrainGrace` bounds
the *idle* gap, which is what the "empty pipe with a live writer" test actually needs), or
add a measured basis for a whole-tail 50 ms drain and shrink MVV row 6's payload to it.

---

## MEDIUM

### M1 — `0026:S8` vs `0026:C1` `bound:` — two different totals, no stated precedence

`C1 bound:` commits "the COMMITTED return bound is `timeout + 2·WaitDelay`" with no tolerance.
`S8` expects "within `timeout + 2·WaitDelay + 100 ms` (the premortem's P-8 tolerance)" for a
child that resists termination. `C1 bound:` names S8 as the case where "the total tightens" —
it does not say the total may be *exceeded* by 100 ms.

Missing pass/fail criterion: is `timeout + 2·WaitDelay` a hard contract that S8 must meet
(in which case S8's assertion is wrong and too loose), or does the contract carry a
scheduling tolerance (in which case `C1 bound:` should carry the +100 ms and the MVV's
"elapsed ≤ 3 s, asserted with margin" should say what the margin is)?

Test it prevents: S8's timing assertion, and any CI-hardening of the MVV row 3 bound. A run at
`timeout + 2·WaitDelay + 60 ms` passes S8 and violates C1 — I cannot tell whether to file that
as a bug.

**Fix shape**: put the tolerance in `C1 bound:` once, and have both S8 and MVV row 3 cite it.

### M2 — `0026:S1` — "the CLI renders the 'may have been applied' text" has no oracle

S1's expected end-state, direct-write leg, is the *proof the re-key landed*. It asserts a
rendered string. The record names the producing symbol (`flow_exec.go::accessorFailureOf`'s
`ClassExecutionFailure` arm) and the helper (`detailMayHaveApplied`) but never states what is
observable at the CLI boundary: no exit code, no output stream, no stable substring, no
golden-file convention.

Missing pass/fail criterion: the assertion surface. `Applied()` is a Go-level assertion I can
write today (`Refusal.applied`, per the `authority` mini-check's canonical key); the *rendered
text* is a second, unspecified one, and S1 explicitly requires BOTH ("`Applied()` is true AND
the CLI renders...").

Test it prevents: the direct-write-leg end-to-end assertion in S1 — the only row that proves
the A8 gain reached a user. I can write the unit half and not the CLI half.

**Fix shape**: either name the observable (helper + stream, or "asserted at `Applied()`, with
rendering covered by `flow_exec`'s existing golden"), or drop the rendering conjunct and let
S1 assert the field.

### M3 — `0026:C1` `refusal:` / `0026:S1` / `0026:MVV` row 5 — only the both-pipes-held case has an oracle

`C1 refusal:` says `Detail` names the pipe(s) held — `"stdout"`, `"stderr"`, or both. Every
scenario that asserts on it uses the both case: MVV row 5 asserts `Err` naming `"stdout,
stderr"`; S1 says "no envelope parsed" without a per-pipe form. The single-pipe forms are
contract surface with no scenario and no stated rendering (is it `"stdout"`, `"stdout only"`,
what separator, what order?).

More consequential for testing: **stderr-held-alone**. `C1 refusal:` composes the held-pipe
reason "AHEAD of the stderr tail collected up to the bound" — but if it is *stderr* that is
held, the tail is exactly the truncated thing, and `D-the-drain-grace` justifies `DrainGrace`
on precisely this ground ("what the grace preserves is the stderr tail... the text that tells
a user WHICH helper held the pipe"). No scenario exercises stdout-clean/stderr-held, so the
grace's stated purpose has no test.

Missing pass/fail criterion: the `Detail` shape for a single held pipe, and the expected tail
content when stderr is the held one.

Test it prevents: a stderr-held-alone case, and any assertion on the single-pipe `Detail`
string. Also leaves `D-the-drain-grace`'s load-bearing rationale unverified by any row.

**Fix shape**: add a scenario (or a row to S1) for stdout-EOF/stderr-held, asserting the tail
survives the grace; state the single-pipe `Detail` form in `C1 refusal:`.

---

## LOW

### L1 — `0026:S5` — "logged at the check" names no observable

S5's expected end-state is "the condition is logged at the check, not discovered at the join;
the fallback is the documented unbounded path, and the log names the host condition." No
level, no sink, no message key, no assertion surface. `F5` repeats the same prose. The
`disposition` mini-check calls this row "**loud at the check, by design**" — loudness is the
whole point of the row, and it is the one property I cannot assert.

Missing pass/fail criterion: what a test observes to conclude the log happened.

Test it prevents: S5's positive assertion. I can assert the fallback (join stays unbounded on
a simulated non-pollable end) but not the "loud" half, which is the design intent.

**Fix shape**: name the sink and a stable key in `C1` or S5 — this is a one-clause fix and the
row becomes fully writable.

### L2 — `0026:S3` — "byte-for-byte unchanged from today" has no named baseline

S3 and the `fidelity` mini-check both commit byte-equality against *today's* behaviour. The
record names no baseline artifact: no golden file, no pre-change capture, no existing test
whose output is the reference. The `oracle` mini-check flags S3 as an absence-of-change
oracle and discriminates it with the cap boundaries — which is a real control, but it tests
the cap, not the byte-equality claim.

Missing pass/fail criterion: what "today" is, mechanically, at assertion time.

Test it prevents: the byte-equality half of S3. The cap-boundary half is writable as
specified (exactly-cap is a value; cap+1 is a detected overflow).

**Fix shape**: cite the existing test(s) whose expectations are the baseline, or state that
S3 is satisfied by the unchanged existing suite passing.

### L3 — `0026:MVV` row 6 / `0026:S7` / `0026:S9` — "artificially delayed / stalled drain" is unspecified

Three rows depend on a *stalled drain* as the input, and the mechanism is never named: a test
hook, a build tag, an injected sleep, a seam in `readBounded`? S7 is explicitly a
regression trap for a future refactor ("fails here rather than silently voiding A3"), so the
seam is long-lived production-adjacent surface, not scaffolding.

Missing pass/fail criterion: none, strictly — the *expectations* are precise (65536 bytes as a
correct prefix; whole 1 MiB when prompt). What is missing is the *input* construction, and the
three rows likely need the same seam.

Test it prevents: nothing outright — it costs an implementer a design decision that three
scenarios silently share, with a risk that MVV row 6 and S7 build two different hooks.

**Fix shape**: one line naming the seam, cited by all three rows.

---

## Rows I checked and found fully testable (no finding)

- `0026:MVV` rows 1–5: fixture, timings, class, `Detail` emptiness, envelope-not-parsed, and
  the residue assertion (grandchild alive, killed by the test) are all stated with controls.
- `0026:S2`: the absence-of-change risk is explicitly acknowledged and paired with S7.
- `0026:S4`: `exec.ErrWaitDelay` → non-`ExitError` arm → `execution_failure`, ordering stated
  in `C1 stdin:`; A7 verified.
- `0026:S6`: `timeout` + empty `Detail`, and `F6` records the quiet-reason consequence as
  deliberate.
- `0026:S7` prompt arm and the 65536-byte prefix assertion.
- Read-end close: `§existing-infrastructure-audit` settles that the defers already satisfy
  `C1 precedence:` and that the shape's `closeReads` would cost a swallowed `ErrClosed`, not a
  panic — no double-close test is owed.
