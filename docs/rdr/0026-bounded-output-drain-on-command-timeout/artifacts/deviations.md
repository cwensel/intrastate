# Deviations — RDR 0026 bounded-output-drain-on-command-timeout

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

## D1 — Negative witness for the `Applied()`-keyed rendering

**Type**: TEST-FIXTURE. **Status**: OPEN (pre-seeded, 7.1 pairwise 0026×0028 PW3; critique C-9) → CLOSED at Stage 8: the named check ran and passed, contradicting no contract.

**Disposition (Stage 8).** The check is
`internal/cli/held_pipe_0026_test.go::TestReq138_AnUnappliedWritePhaseExecutionFailureRendersNoAppliedProse`
(write phase, `ClassExecutionFailure`, `Applied()` false → no
`detailMayHaveApplied`, and `flow-accessor-failed` kept) and
`::TestReq34_AReadOrGatePhaseHeldPipeRefusalNeverRendersTheAppliedProse`
(read/gate refusal wrapping the typed held-pipe error → `Applied()` false),
with the executor-side halves at
`internal/accessor/held_pipe_0026_test.go::TestReq116` and `::TestReq34`.
All pass. The `errors.As` for the held-pipe error was landed in
`executor.go::Write`'s `err != nil` arm and NOT in the shared `refusalOf`,
so the read and gate paths never reach it; the CLI arm keys on `Applied()`
and never on the phase, so 0028's pre-mutation `edit` refusals keep
`flow-accessor-failed` with no applied prose.

**Situation.** `0026:S1` asserts only the positive: `Applied()` true on both
write legs renders `detailMayHaveApplied`. No scenario asserts the negative —
a WRITE-phase `execution_failure` with `Applied()` false (every refusal 0028's
`edit` binding mints, and today's non-zero-exit command writer) must render
WITHOUT `detailMayHaveApplied`; and a read or gate refusal wrapping the typed
held-pipe error stays `Applied()` false (`0026:C1` `refusal:` "as today"). The
neighbouring `ClassTimeout` arm in `flow_exec.go::accessorFailureOf` keys on
phase; copying that shape would label 0028's pre-mutation refusals applied.

**Check.** A unit test on `accessorFailureOf`: (write phase,
`ClassExecutionFailure`, `Applied()` false) → no `detailMayHaveApplied`;
(read phase, refusal wrapping the typed held-pipe error) → `Applied()` false.
The `errors.As` for the held-pipe error lives in `Executor.Write`'s
`err != nil` arm, never in the shared `refusalOf` — `0026:A8`'s "either
placement compiles" is the trap.

## D2 — The `Applied()` key on the direct write arm is a distinct CLI code

**Type**: IMPL-DECISION. **Status**: decided at JDR 0003 §D4 (b), 2026-09-03 — recorded so Stage 8's additive-surface rule reads it as an author decision, not a SPEC-UNDER.

**Situation.** `0026:A8`/`0026:S1` name the gain as "the CLI's 'may have been
applied' rendering keys on `Applied()`" and assert the `detailMayHaveApplied`
prose. §D4 fixes the form: a distinct exit-3 CLI code for the applied write
refusal (spelling non-normative, e.g. `flow-write-failed-applied`), keyed on
`Applied()` in `flow_exec.go::accessorFailureOf`'s write arm; `detail` text
unchanged; not-applied refusals keep `flow-accessor-failed`. JDR 0001 §D10's
table gains the code by citation repair.

**Check.** S1's rendered half asserts the new code AND the prose on the
direct write leg; D1's negative witness asserts `flow-accessor-failed` with
no prose on the `Applied()`-false write refusal.

## D3 — `Detail` field ordering: held-pipe reason leads, stderr tail follows

**Type**: TEST-FIXTURE. **Status**: mechanical translation (resolved at Phase 0
from the record's own text; no author decision required).

**Situation.** Phase 0 surfaced a contradiction inside the record. `0026:C1`
`refusal:` orders the held-pipe reason "composed AHEAD of the stderr tail
collected up to the bound", and the unfenced Technical Design agrees (`Detail`
"leading with the held pipe(s) and the child's exit status ahead of the stderr
tail"). Scenario S10 states the inverse — the tail "survives into `Detail`
ahead of the held-pipe reason". The two readings give materially different
assertion text on S10, the only scenario pinning the single-pipe `Detail` form.
Neither predecessor resolves it: 0025:C4 orders `Detail` generally but not this
pair, and 0004 does not reach it.

**Resolution.** The contract wins: the held-pipe reason leads, the stderr tail
follows. Grounds: (1) `C1 refusal:` is the normative fence and a scenario is
evidence about a contract, not a peer of it — the precedence idiom this record
itself applies against 0025:C4; (2) the Technical Design corroborates the fence
independently, so S10 stands alone against two agreeing statements; (3) S10's
load-bearing claim is that the tail SURVIVES into `Detail` (its own stated
purpose: "without this row the grace's load-bearing justification has no test"),
which holds under either order — the ordering words are incidental to the
assertion S10 was written to pin.

**Binding on.** REQ-26, REQ-61, REQ-83 keep the contract ordering. S10's test
asserts the tail's SURVIVAL and the bare single-pipe name, not the tail's
position ahead of the reason.

---

## D4 — The per-read grace alone cannot end a drain that keeps making progress, so S9's companion and C1 `precedence:` cannot both hold as written

**Type**: SPEC-DEFECT.

Status: needs author decision → RESOLVED (accepted 2026-09-03: superseded by D8+D9+D10 — the fixture is repaired, the sibling halt is the named terminator, and the residual both-pipes shape is an admitted residue; no ceiling, no second bound)

**Situation.** Two normative statements in this record select opposite
mechanisms on the same input.

C1 `precedence:` fixes the mechanism: "The deadline is re-armed at `now +
DrainGrace` BEFORE EACH read of the final drain, so the grace bounds the
IDLE GAP between reads and never the size or total duration of the tail",
and "a drain making progress keeps its deadline ahead of it and runs to EOF,
and only a gap longer than the grace ends it as held." MVV row 6(b) is the
normative witness for exactly that: 1 MiB paced across 32 chunks with idle
gaps of 0.4·DrainGrace, a final drain measured HERE at ~605 ms — about
12·DrainGrace — which must be recovered whole.

S9's companion (TS scenario 9, REQ-131, and MC `oracle`'s S9 row) fixes the
opposite outcome for a writer whose gaps never exceed the grace: "an escaped
writer emitting a byte every 10 ms forever is REFUSED (never blocked, never
EOF)." A 10 ms gap is a FIFTH of DrainGrace, so under C1's own predicate that
drain re-arms forever, never reports held, and — because C1 `precedence:` (b)
also forbids the parent proceeding past a live drain — `spawn` never returns.
That is the hang this record exists to end, reintroduced by the record's own
mechanism.

C1's prose names the discriminating case as "a writer trickling bytes SLOWER
than the grace", which is consistent with the mechanism; the 10 ms fixture
S9 specifies trickles five times FASTER. The two cannot both be satisfied by
any predicate keyed on the idle gap.

**Grounding.** `evidence/spikes/a9-a10-regrace/` is the mechanism's own
evidence and it measures only the terminating shapes (burst and paced); no
spike exercises a writer that never reaches EOF, so the evidence base does
not decide this. Re-measured here: the paced final drain runs 605-608 ms
against a 500 ms `WaitDelay`, so no cap of one `WaitDelay` can hold row 6(b).
Prior art does not resolve it either — `os/exec/exec.go::awaitGoroutines`
hard-closes at `WaitDelay` and accepts the loss, which is precisely the
disposition `D-the-drain-grace-is-a-mechanism-constant` records this record
as diverging from, and the divergence is what creates the gap: Go's mechanism
is unconditionally live, this one is not.

**What was implemented (Phase 2, SUPERSEDED at Phase 3c).** The per-read
re-arm exactly as C1 states it, plus ONE liveness ceiling on the final drain
of `2·WaitDelay` measured from the drain's own entry into the grace. Phase 3a
and Phase 3b, working independently, both reproduced the same two defects in
it and the record refutes it on both counts:

- It bounds the tail's TOTAL DURATION, which C1 `precedence:` rejects by name
  ("the grace bounds the IDLE GAP between reads and never the size or total
  duration of the tail"; "A single absolute deadline would instead bound the
  whole remaining tail against the clock, which turns a slow-arriving tail
  into a false held-pipe report on a pipe with no writer at all"). The
  ceiling IS that single absolute deadline with a larger constant. Witness:
  a paced tail, every idle gap 40 ms against the 50 ms grace, reaching a real
  EOF with NO holder, truncated at 131072 of 163840 bytes and reported held
  (ADV-2 / FAIL-1). Phase 2's claim that the ceiling "never binds on a drain
  that reaches EOF" was measured against one shape only (A9's 1 MiB / 32
  chunks, 0.605 s); nothing in the contract bounds a tail's total duration.
- It STACKS on the parent's join timer rather than being carved out of it —
  `WaitDelay + DrainGrace + 2·WaitDelay` — so a trickling escapee returned
  at 1.825 s (unit) / 3.562 s (with a declared timeout) against C1 `bound:`'s
  committed `timeout + 2·WaitDelay` + 100 ms (ADV-3 / FAIL-2).

MVV row 6(b) as literally written (~640 ms) passes only because it sits under
the 1000 ms ceiling, so the normative witness never exercised the defect.

Decisive: the record admits no such constant. C1 `bound:` says `WaitDelay` is
"the ONE bound" and `DrainGrace` "is NOT a second bound". The ceiling was an
unnamed THIRD bound the contract forbids.

**What is implemented now (Phase 3c).** The ceiling is REMOVED. `readBounded`
reads no clock beyond the grace itself. In its place the two drains share one
`atomic.Bool` HALT, raised by a drain that has ITSELF REPORTED a refusing
condition (held, or overflow) and never on elapsed time. Once any drain has
reported, the invocation's outcome is already a refusal, so no byte a sibling
could still deliver would be parsed; the sibling therefore ends and reports
the condition it has (held — it did not reach EOF), and the join the parent
must complete returns.

Why this satisfies the constraints the ceiling broke:

1. **No total-duration bound.** Nothing keyed on a clock. On a drain with no
   holder no sibling ever reports held, so the halt is never raised and a
   paced tail runs to EOF however long it takes — MVV rows 6(a) and 6(b), S9
   and ADV-2 all recover the whole payload.
2. **Inside the committed bound.** The halt fires one `DrainGrace` after the
   join timer, so the drain leg costs `WaitDelay + DrainGrace` rather than
   `3·WaitDelay`. Measured: the trickler returns in 1.07 s against 1.1 s
   (ADV-3), and FAIL-2's shape — trickler plus a child sleeping past a 2 s
   declared timeout — in 2.61 s against 3.1 s (was 3.56 s).
3. **`WaitDelay` stays the ONE bound.** No constant is added.
4. The halt is checked BEFORE the grace re-arm and independently of it: a
   drain still making progress never sees a deadline error, so gating it on
   the grace would leave the trickler running (this was verified — gating it
   reproduced the hang).

**How S9's companion terminates.** By the second, independent reason Phase 2
already noted: the fixture's escapee inherits BOTH pipes and drips on stdout
only, so the STDERR drain is empty with a live writer and reports held under
the grace alone. Its halt then ends the stdout drain, which the idle-gap
predicate alone could not (a 10 ms gap is a fifth of `DrainGrace`).

**Residual gap — the part of S9 that is NOT closed.** A writer trickling
under the grace on BOTH pipes leaves no drain to report held first, so no
halt is raised and `spawn` does not return. The spec defect at the head of
this entry is therefore narrowed, not eliminated: C1's idle-gap predicate
still cannot, by itself, terminate a drain that makes progress forever, and
the halt only propagates a refusal that some OTHER drain already established.
No shipped fixture or scenario exercises the both-pipes trickler; S9's
companion is the record's only trickle witness and it is single-pipe.

This is the reading Phase 3c took because the constraints are jointly
unsatisfiable without amending the record: (1)'s witness is a SCENARIO, while
(2), (3) and (4) are the contract's published fences — C1 `precedence:` and
C1 `bound:` — and this record's own precedence idiom is that the fence
outranks a scenario (see D3, decided that way at Phase 0). The alternative,
keeping a ceiling, fails two normative clauses to satisfy one scenario's
fixture shape.

**Alternative not taken.** Bound the final drain at one `WaitDelay`, keeping
the record's single constant. Rejected because it truncates MVV row 6(b) at
~500/605 of its payload, failing a Normative row to satisfy a scenario — and
it is the same total-duration mechanism C1 `precedence:` rejects, merely at a
smaller constant.

**Recommendation (unchanged in substance, narrowed in scope).** Either:

- (a) Correct S9's companion fixture to a gap LONGER than the grace, which is
  what C1's own prose describes ("a writer trickling bytes SLOWER than the
  grace") and what the idle-gap predicate discriminates. Under (a) the
  residual gap closes with no mechanism change at all, and the halt remains
  only as the thing that keeps a sibling from outliving an established
  refusal. This is the recommended option: it repairs a fixture that
  contradicts the contract it was written to witness.
- (b) Admit a terminating mechanism for a drain that makes progress forever
  and NAME it in C1 `bound:` beside `DrainGrace`, since every shipped
  constant of this family is named in its contract. Note this cannot be a
  duration without reintroducing the total-duration bound C1 `precedence:`
  rejects; a byte-progress or read-count predicate would be needed instead.

Until the author decides, the shipped behaviour is the one above: every MVV
row, every scenario and every adversarial witness is green, and the residue
is the both-pipes trickler, which no fixture exercises.


### D4 addendum — grounding sweep and the standing recommendation (Stage 8, post-3c)

A grounding pass over the record, JDR 0001-0003, the codebase and Go prior art
was run after Phase 3c. It does not change D4's status; it settles WHICH of the
two options the author should take and adds one measured fact neither option
accounted for.

**Measured: the residual gap is real, and S9's companion passes incidentally.**
A `setsid` escapee trickling BOTH pipes every 10 ms (each gap 0.2x DrainGrace)
was driven against the shipped drain discipline. Neither drain ever reports a
refusing condition, so neither raises the halt, and the join blocks
indefinitely — the hang is unchanged in that shape. The single-pipe control
(S9's actual `drip` fixture, `fixtures_0026_test.go:213`, which writes to
STDOUT ONLY) shows why the shipped suite is green: the idle stderr drain
reports held at ~602 ms under the grace alone and its halt is what ends the
stdout drain. The stdout drain never reports on its own. So S9's companion is
passing for a reason its own text does not name, and no shipped fixture
exercises the both-pipes shape.

**The `DrainBound` criterion, applied to the sibling halt.** LBD
`D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` states the
test a new mechanism must pass: `DrainBound` was refused for "answering the
SAME question with a second value"; `DrainGrace` was admitted because it
"answers a different one ... and invents no second signal (the Sibling-path
check's test)". The halt invents no constant and reads no clock — it carries
the refusal decision the drain ALREADY made to the drain that cannot reach one.
It answers "this invocation's outcome is already a refusal", not "how long to
wait". It therefore passes the record's own criterion, and it is the only
candidate examined that does: every duration-shaped alternative reintroduces
the total-duration bound `precedence:` rejects by name.

**Codebase conceptual integrity.** The halt is novel here: `atomic.Bool` occurs
nowhere else in `internal/`, there is no errgroup, and `cmdbind.go` is the only
non-test file spawning goroutines. There is no house pattern to match or
diverge from. It is, however, the canonical Go broadcast-cancellation shape
(TGPL 8.9: one flag/closed channel observed by every worker), and it is what
`os/exec` itself does in spirit — `awaitGoroutines` reacts to a decision
already taken rather than each goroutine re-deciding. The comment at
`cmdbind.go:434` ("no shared flag is needed") predates the halt and is now
inaccurate; it describes the JOIN's mechanism, not the drain's, and should be
corrected when the record is amended.

**Prior art, restated precisely.** Go bounds this case with one absolute
`WaitDelay` and then HARD-CLOSES the read ends (`awaitGoroutines` ->
`closeDescriptors`), accepting truncation to guarantee liveness. This record
diverges deliberately (LBD, recorded) to preserve the stderr tail that names
WHICH helper held the pipe. That divergence is about preserving buffered bytes
on the final read; it is NOT a claim about how a drain terminates when the
writer never stops. Nothing in the divergence forbids a non-durational
terminator, which is why the halt is compatible with it and a ceiling was not.

**Recommendation (unchanged in direction, now grounded): take option (a), and
add the both-pipes case to it.**

1. Correct S9's companion fixture to a gap LONGER than `DrainGrace` — C1's own
   prose already names the discriminating case as "a writer trickling bytes
   SLOWER than the grace", so the fixture as written contradicts the contract
   it was authored to witness. This is the same fence-over-scenario idiom the
   record applies to 0025:C4 and that D3 applied at Phase 0.
2. Name the sibling halt in C1 `precedence:` as the terminator, stated as a
   state predicate and explicitly NOT a second bound — one sentence, e.g. "a
   drain that has itself reported a refusing condition ends its sibling, so the
   join completes without any duration governing the tail." It is a consequence
   of `precedence:` (b)'s "MUST NOT read until every drain goroutine has
   returned", not new policy, and naming it keeps `bound:`'s "WaitDelay is the
   ONE bound" literally true.
3. Record the both-pipes-trickle shape as a named residue in C1 `residue:`,
   in F5's shape — "recorded at the check rather than discovered at the join".
   F5 already accepts exactly this: on a non-pollable host "the join stays
   unbounded — the old hang, confined ... and recorded". The both-pipes trickle
   is the second member of that class and deserves the same honest treatment.
   Closing it instead would require either the hard-close Go takes (surrendering
   the tail the whole divergence exists to preserve) or a duration bound
   (`precedence:` forbids it) — so naming it is the only option consistent with
   the record's own commitments.

Option (b) — naming a terminating mechanism in `C1 bound:` — remains
not-recommended for the reason already given: it cannot be a duration without
reintroducing the rejected mechanism, and as a state predicate it belongs in
`precedence:` (where step 2 puts it), not in `bound:`.

## D5 — The first deadline error is the mark, not the verdict

**Type**: IMPL-DECISION. **Status**: mechanical translation (derived from
C1 `precedence:` and the A2 spike; no author decision required).

**Situation.** C1 `precedence:` (a) fixes the parent's `SetReadDeadline` on
the read ends as THE cross-goroutine mark and forbids a separate `graceOn`
flag, while the Illustrative Code puts the re-arm inside the drain's loop.
The drain must therefore learn from the read itself that the grace is on: a
loop that re-armed from its first iteration would deadline the PRE-timer
path, which C1 forbids in the same clause.

**Resolution.** The first `os.ErrDeadlineExceeded` a drain sees is the
parent's mark. The drain records it, re-arms `now + DrainGrace`, and takes
its FINAL read under its own grace; that read's result is the report. This
is C1's "refuse on no EOF after a FINAL BOUNDED READ, not on elapsed time"
(P-2) read literally, and it is what makes the pre-timer path deadline-free
without any flag. Grounded in `evidence/spikes/a2-unblock-read/` finding 4:
a FUTURE deadline delivers buffered bytes on the first read and reports the
timeout only on a subsequent read of an empty pipe, so one extra read is
exactly what the mechanism needs. Cost: at most one extra `DrainGrace`
(50 ms) per drain on the refusing path, well inside the committed total
(MVV rows 1-4 measured at 2.884 s against 3.1 s).

## D6 — The mandated package-level seams are unsafe under `t.Parallel()`

**Type**: TEST-FIXTURE. **Status**: mechanical translation.

**Situation.** REQ-128 (S7) mandates the stall seam as "one unexported
package-level HOOK set by the test and nil in production", and REQ-126 (S5)
simulates F5 the same way. Six Phase-1 tests both mutate one of those
globals and declare `t.Parallel()`. A concurrent invocation then observes a
seam it did not set: with `nonPollableForTest` set by a parallel peer, every
other invocation in flight takes F5's documented UNBOUNDED join and hangs
for the fixture's 60 s, and `TestReq128`'s "nil in production" read observes
a peer's stall value. The suite failed non-deterministically for that reason
alone; every one of those tests passed in isolation against the same build.

**Resolution.** `t.Parallel()` removed from the six tests that read or write
a mandated global seam: `TestReqMVV6a`, `TestReqMVV6b`, `TestReq21`,
`TestReq51`, `TestReq128` and `TestReq126`. Go pauses parallel tests until
the sequential ones in the package have finished, so serialising them is
sufficient and no assertion, expectation, fixture or bound was touched. The
alternative — making the seam per-invocation — would contradict REQ-128's
"one unexported package-level hook" by name.

**Cost.** The package's wall clock is unchanged in practice (15.4 s with the
change; the parallel run failed rather than finishing faster).

## D7 — F5's `Detail` line rides every refusal the invocation produces

**Type**: IMPL-DECISION. **Status**: mechanical translation.

**Situation.** REQ-107 / REQ-126 fix F5's observable as "a `Detail` line on
ANY refusal that invocation produces", and S5's fixture refuses at the
BINDING's parse layer (`Reader.Read`'s "produced no stdout and exited 0"
arm), not inside `spawn`. Composing the line only at `spawn`'s own refusal
arms would leave S5's row — and every other post-`spawn` refusal — without
it.

**Resolution.** One composition point, `invocation.detail`, which orders the
ONE `Detail` slot as C1 `refusal:` and deviations D3 fix it: the held-pipe
reason LEADS, F5's host-condition line follows, and the stderr tail
collected up to the bound comes last. Every `wrap` site that carries this
invocation's tail now calls it, so the F5 line is a property of the
invocation rather than of the refusal arm that happened to fire.

## D8 — S9's companion fixture contradicts the contract it was written to witness

**Type**: TEST-FIXTURE.

Status: needs author decision → RESOLVED (accepted 2026-09-03: correct the companion's cadence to a gap LONGER than DrainGrace; fixture repaired in the suite, contract text unchanged — code is the source of truth and the deviation is the record)

**Situation.** TS scenario 9's companion reads "an escaped writer emitting a
byte every 10 ms forever is REFUSED (never blocked, never EOF)." `DrainGrace`
is 50 ms, so a 10 ms cadence is 0.2x the grace: under C1 `precedence:`'s own
predicate that writer is MAKING PROGRESS, keeps its deadline ahead of it, and
is never reported held. C1's prose names the discriminating case in the
opposite direction — "a writer trickling bytes SLOWER than the grace is
exactly the case that must refuse". The fixture and the fence select opposite
outcomes on the same input, and the fence is what the fixture was authored to
witness.

**Measured (Stage 8, post-3c).** The shipped fixture is `drip`
(`internal/cli/cmdbind/fixtures_0026_test.go:213`) and it writes to STDOUT
ONLY. Its inherited stderr is never written, so the stderr drain sits empty
with a live writer, reports held under the grace alone at ~602 ms, and raises
the sibling halt that ends the stdout drain. The stdout drain never reports a
refusing condition on its own. S9's companion therefore passes for a reason
its own text does not name; a two-pipe trickler at the same cadence hangs
(see D10).

**Recommended amendment.** Change the companion's cadence to a gap LONGER
than `DrainGrace` — the shape C1 `precedence:` already describes — and keep
the assertion (refused, named pipe, inside the committed bound) unchanged.
The 10 ms figure is the defect; the scenario's purpose is sound.

**Precedent.** Fence-over-scenario, the idiom this record applies to 0025:C4
and that D3 applied at Phase 0: a scenario is evidence about a contract, not
a peer of it.

**Code impact.** None. `TestReq105`/`TestReq131` assert the refusal, which the
shipped implementation delivers; only the fixture's cadence and the scenario
prose change. Re-pointing the fixture at a slower cadence keeps the row green
for the reason the row states.

## D9 — The terminator that ends a drain making progress is unnamed in C1

**Type**: SPEC-UNDER.

Status: needs author decision → RESOLVED (accepted 2026-09-03: the sibling halt is the named terminator, a state predicate and NOT a second bound; captured at the mechanism in cmdbind.go rather than by amending the locked record)

**Situation.** C1 `precedence:` (b) requires that the parent "MUST NOT read
`inv.stdout`/`inv.stderr` until every drain goroutine has returned", and
forbids a design "in which the parent proceeds past a live drain". C1 `bound:`
states `WaitDelay` is "the ONE bound" and `DrainGrace` "is NOT a second bound".
Between them, nothing in the contract says what ENDS a drain that is making
progress into output no caller will ever parse. The grace cannot: it bounds
the idle gap, and a progressing drain never trips it. So the contract as
written admits a shape with no terminator, which is the D4 contradiction's
mechanical root.

**What ships.** A sibling halt: the two drains share one `atomic.Bool`; a
drain raises it when it has ITSELF reported a refusing condition (held or
overflow), and a drain observing it set returns held (or overflow, preserving
the rank). Nothing reads a clock. From the moment either drain reports a
refusing condition the invocation's outcome is already a refusal, so no byte a
sibling could still deliver would be parsed.

**Why this shape and not a bound.** LBD
`D-the-drain-grace-is-a-mechanism-constant-not-the-rejected-knob` fixes the
test a new mechanism must pass: `DrainBound` was refused for "answering the
SAME question with a second value"; `DrainGrace` was admitted because it
"answers a different one ... and invents no second signal (the Sibling-path
check's test)". The halt introduces no constant and answers a different
question again — "is this invocation's outcome already a refusal", not "how
long to wait" — so it passes that criterion. Every duration-shaped alternative
fails it by reintroducing the total-duration bound `precedence:` rejects by
name; that is exactly how the Phase 2 ceiling failed (ADV-2/ADV-3, FAIL-1/2).

**Prior art.** The canonical Go broadcast-cancellation shape (TGPL 8.9: one
flag observed by every worker), and the same posture `os/exec` takes in
`awaitGoroutines` — react to a decision already made rather than have each
goroutine re-derive it. Go's own resolution of this case is a hard close at
`WaitDelay` (`closeDescriptors`), which this record deliberately diverges from
to preserve the stderr tail; that divergence concerns preserving buffered
bytes on the final read and says nothing about how a drain terminates, so it
does not forbid a non-durational terminator.

**Recommended amendment.** One sentence in C1 `precedence:`, stated as a state
predicate and explicitly not a second bound: a drain that has itself reported
a refusing condition ends its sibling, so the join completes without any
duration governing the tail. It is a consequence of `precedence:` (b), not new
policy, and naming it keeps `bound:`'s "WaitDelay is the ONE bound" literally
true.

**Code impact.** None — the implementation already has this shape. One stale
comment to correct when the record is amended: `cmdbind.go` at the join site
reads "no shared flag is needed", which predates the halt and describes the
JOIN's mechanism, not the drain's.

## D10 — A writer trickling under the grace on BOTH pipes has no terminator

**Type**: SPEC-UNDER.

Status: needs author decision → RESOLVED (accepted 2026-09-03: the both-pipes trickle is an ADMITTED, named liveness residue in F5's shape — recorded at the check, not discovered at the join; captured at the mechanism in cmdbind.go)

**Situation.** The sibling halt (D9) terminates a drain only once some drain
has reported a refusing condition on its own. When an escapee trickles BOTH
pipes with every gap under `DrainGrace`, neither drain ever reaches a refusing
condition, so the halt is never raised and the join blocks indefinitely — the
unbounded hang this record exists to end, in a shape the record does not name.

**Measured (Stage 8, post-3c).** Driven against the shipped drain discipline
with a real `setsid` escapee writing one byte to each of stdout and stderr
every 10 ms (0.2x `DrainGrace`), no drain reported and the join did not
complete. The single-pipe control (S9's `drip`) terminated at ~602 ms via the
idle stderr drain's held report. No shipped fixture exercises the two-pipe
shape; S9's companion, the record's only trickle witness, is single-pipe by
construction.

**Why it cannot simply be closed.** The two available mechanisms are both
already ruled out by this record. A duration bound on the drain is the
total-duration bound C1 `precedence:` rejects by name (and is what ADV-2/ADV-3
caught). A hard close of the read ends is Go's answer
(`awaitGoroutines` -> `closeDescriptors`) and surrenders the stderr tail
naming WHICH helper held the pipe — the sole justification for this record's
deliberate divergence from that prior art. Closing D10 therefore means
reopening a decision the record already took, which is an author call and not
an implementation one.

**Recommended amendment.** Name it as a residue in C1 `residue:`, in F5's
shape — "recorded at the check rather than discovered at the join". F5 already
accepts precisely this trade for a non-pollable read end: "on that host the
join stays unbounded — the old hang, confined to a non-pollable host and
recorded at the check". The two-pipe trickle is the second member of that
class and warrants the same honest treatment: an admitted, named liveness
residue rather than a silent one. `residue:` today admits leaking a process
and its output while insisting "a withheld refusal is not" admitted — a
withheld refusal is exactly what this shape produces, so the clause needs the
carve-out stated, not implied.

**Alternative, if the residue is judged unacceptable.** Re-open the
hard-close decision for the both-pipes case only: once BOTH drains are known
to be progressing past the bound there is no stderr tail worth preserving
(the tail is being actively overwritten by the trickle), so Go's close could
be taken there without losing the diagnostic the divergence protects. That is
a larger amendment touching `precedence:`, `bound:` and the LBD, and is not
recommended without a fresh spike.

**Code impact.** None for the residue option. The alternative would require a
new terminator path and a fixture for the two-pipe shape.
