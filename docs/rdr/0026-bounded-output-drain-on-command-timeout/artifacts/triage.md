# Triage — RDR 0026 bounded-output-drain-on-command-timeout

Single-pass roborev triage after the Stage-8 launch reached COMPLETE.

- Window frozen once: BASE `f85fde7` .. HEAD `2e43b57`.
- Batch label: `batch:rdr-0026`.
- Spine: 6 in-window per-commit auto-reviews (of 60 open jobs repo-wide; the other
  54 are on `main` and unrelated branches, outside the window and not this run's).
- 11 findings across those 6 jobs. Every finding is graded against the state at
  HEAD, not at the commit it was raised on — the launch's red-then-green rhythm
  means an early finding is often already resolved.

## Findings

| # | job | finding | verdict | odc-type | odc-trigger | outcome |
|---|---|---|---|---|---|---|
| 1 | 6717 | sibling halt raised pre-reap stalls a live child's sibling drain | FIX-NOW | timing-serialization | concurrency | fix-now:f9dca2d |
| 2 | 6721 | S9 trickle oracle does not discriminate (fixture inherits idle stderr) | FIX-NOW | test-oracle | test-coverage | fix-now (stage 2) |
| 3 | 6706.1 | trickle cadence 10 ms under a 50 ms grace never reports held | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | drop |
| 4 | 6706.2 / 6710.2 | large payloads mean the late-drain rows never reach the final-read path | DROP superseded-at-HEAD | n/a-not-a-defect | test-coverage | drop |
| 5 | 6706.3 | package-global test seams mutated from parallel tests | DROP rdr-adjudicated (D6) | n/a-not-a-defect | concurrency | drop |
| 6 | 6706.4 / 6710.3 | TestReq79 cannot stage the ErrWaitDelay it claims to outrank | KATA-BUG | test-oracle | logic-flow | kata 048g |
| 7 | 6710.1 | held stderr past 1 MiB laundered into overflow, refusal withheld | DROP superseded-at-HEAD | n/a-not-a-defect | boundary | drop |
| 8 | 6710.4 | CI is Ubuntu-only, so REQ-133's darwin half never runs | KATA-BUG | environment | test-coverage | kata 5qx5 |
| 9 | 6710.5 | ErrWaitDelay child that exited 0 reported as "did not exit" | KATA-BUG | assignment | rare-situation | kata x6dc |
| 10 | 6712.1 | same laundering root cause as 6710.1, at the ceiling | DROP superseded-at-HEAD | algorithm | boundary | drop |
| 11 | 6714.1 | ADV-3 asserts a tighter bound than C1 owes | DROP superseded-at-HEAD | test-oracle | design-conformance | drop |
| 12 | 6714.2 | ADV-1/2 couple to drainReport, a representation C1 leaves free | DROP rdr-adjudicated | interface | test-coverage | drop |
| 13 | 6714.3 | ADV-2's 40 ms gap leaves only 10 ms under the grace (flake risk) | DROP superseded-at-HEAD | timing-serialization | concurrency | drop |

## The two FIX-NOW defects, in full

### 6717 — the sibling halt stalls a live child's drain (fixed: f9dca2d)

`readBounded` raises the shared halt on ANY read error once past the cap — including
EOF — and the drain goroutines start BEFORE `cmd.Wait()`. A child that writes past
`StdoutCap` to stdout and then closes it therefore raises the halt while it is still
live; the sibling stderr drain observes it, returns, and abandons a pipe the live
child is still filling. The child blocks in `write(2)`, `cmd.Wait()` never returns,
and the invocation stalls to the context deadline — misreported as `timeout`, hiding
the stdout-overflow refusal. The parent's read deadline cannot rescue it: it is armed
only after `cmd.Wait()` returns.

Reproduced: 2 MiB stdout, close, 256 KiB stderr — `cmd.Wait()` blocked >5 s, sibling
drain returned held with 32768 of 262144 bytes; the identical control with no halt
completed in 21 ms with every byte. Independently reproduced a second time from the
in-tree regression test (10.0 s against a 1.1 s committed bound).

Violates `C1 whole:` ("the drains run CONCURRENTLY and are never stalled") and REQ-52,
which is an explicit NEGATIVE REQ: the drains may not be serialized. Also REQ-50/51,
REQ-75, REQ-129 (S7).

Not covered by D9 or D10: both frame the halt entirely in the post-Wait regime and
presuppose a reaped child. No fixture reached it either — `emit-n` and
`overflow-and-hold` write stdout only and exit at once, so the child never writes
stderr after overflow.

Fix: gate ONLY the halt's observe site on a new `reaped` flag set immediately after
`cmd.Wait()`; the raise sites are untouched, so a drain still reports its own
condition normally. This cannot reintroduce the trickler hang: in every shape the halt
exists to terminate the holder is a `setsid` escapee and the refusal lands at
~`WaitDelay + DrainGrace`, necessarily after `cmd.Wait()` returned, so `reaped` is
already true. The gate suppresses the halt only in the window where the direct child
is still live — precisely the window `whole:` forbids stalling a drain in.

### 6721 — the S9 trickle oracle does not discriminate

D8 repaired the cadence (10 ms -> 100 ms) so the stdout drain reaches held on its own
idle-gap mechanism. That part is real: an isolated probe with no sibling and no halt
reports held in ~101 ms, and never returns at 10 ms.

But the fixture still inherits stderr and `drip` writes stdout only, so the idle
stderr drain remains an incidental route to green. Proven by mutation: with
`readBounded` altered so a drain that has read bytes never reports held on its own
idle gap — the exact regression S9's companion exists to catch — both trickle rows
still PASSED, and the refusal text was byte-identical mutated vs unmutated. The
oracles assert only `err != nil` plus an elapsed bound, not the named pipe that D8's
own recommended amendment already called for.

This is the non-discriminating-oracle class the Phase-1 red gate exists to prevent,
and it is D8's own stated scope left incomplete rather than a new request.

Fix: add a `trickle-stdout` fixture variant with stderr CLOSED (so "stdout" alone can
only be named by the stdout drain's own report), point the S9 row at it, and add the
named-pipe assertion. The both-pipes `trickle` case stays for D10's residue row and
TestAdv3, which measure the bound rather than the mechanism.

## Notes on the DROPs

Every DROP rests on positive evidence at HEAD, not on the presumption that a
tests-first commit was deliberately red:

- The late-drain rows carry explicit precondition guards (`if elapsed < stall`) that
  FAIL if the drain was not delayed past the join timer, so the tautology the reviewer
  described cannot hide.
- The parallel-seam race is adjudicated by D6; no seam-touching test declares
  `t.Parallel()` at HEAD, and both setters restore via `t.Cleanup`.
- The stderr laundering is fixed at the source: `overflowed()` keys on the refusal
  boundary (`len(kept) > StdoutCap`), and stderr's limit is exactly `StdoutCap`, so it
  can never rank overflow. `TestAdv1` pins it.
- ADV-2's timing margin was measured, not argued: 0/10 failures under `-race`, and
  3/3 under deliberate 8-way CPU saturation. A9's known sensitivity is real in
  principle but did not bind.
- The white-box coupling in ADV-1/ADV-2/D10 is deliberate and recorded; the same
  observable behaviour has end-to-end coverage through `Reader.Read` (TestReq132,
  TestReq26/29/53), so the unit rows are additive precision, not a substitute.

## Filed

| kata | type | severity | title |
|---|---|---|---|
| 048g | bug | medium | TestReq79 cannot stage the ErrWaitDelay it claims to outrank |
| x6dc | bug | low | An ErrWaitDelay child that exited 0 is reported as did-not-exit |
| 5qx5 | bug | medium | CI runs Ubuntu only, so REQ-133's darwin half never executes |

`048g` and `x6dc` are coupled and should land together: `x6dc` is the defect, `048g`
is why nothing caught it, and `x6dc` has no red test available until `048g` builds the
fixture that can stage a genuine `ErrWaitDelay`. `5qx5` carries an open question on CI
policy (every PR vs nightly vs label-gated), since a macOS runner has cost and flake
implications beyond this record's surface.
