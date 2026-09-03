Model: claude-opus-5

# Spike: A3 (FIFO / no-loss) and A6 (the total bound), over FX-deadline-escape

Two Critical Assumptions of RDR 0026, run live against one shared fixture.

## What was run

```
# darwin (host)
cd <scratch>
go build -o fixture fixture.go
go build -o spike   spike.go
./spike ./fixture

# linux (container)
docker run --rm -v <this-dir>:/w:ro -w /w golang:1.26 sh /w/run-linux.sh
```

`run-linux.sh` copies the two sources into the container, builds them with the
container's toolchain, and runs the same driver. Verbatim output of both runs
is in `output.txt`.

## Platforms

| leg    | platform    | toolchain |
| ------ | ----------- | --------- |
| darwin | darwin/arm64 | go1.26.6 |
| linux  | linux/arm64  | go1.26.8 (`golang:1.26`) |

Both legs ran for real. Neither is UNRUN.

## WaitDelay

`WaitDelay = 500` milliseconds, read from `internal/cli/cmdbind/cmdbind.go`
(the `WaitDelay` constant, declared as `WaitDelay = 500 // milliseconds` and
applied as `cmd.WaitDelay = WaitDelay * time.Millisecond`). The spike mirrors
it as `WaitDelayMS = 500`.

The driver also mirrors the rest of the real `spawn` order: `exec.CommandContext`
under a context deadline, `Setpgid`, a `Cancel` that sends `kill(-pgid, SIGKILL)`,
manual `os.Pipe` ends on `cmd.Stdout`/`cmd.Stderr`, two concurrent bounded drain
goroutines, `cmd.Wait()`, `reapGroup`, then the drain join. The one deliberate
deviation is that the join here is BOUNDED by one more `WaitDelay` timer with
`SetReadDeadline(now)` on the read ends when it fires — which is exactly what
RDR 0026 proposes and what A6 measures. The shipped `drains.Wait()` is unbounded.

## The fixture

`fixture.go` is a Go helper because macOS ships no `setsid(1)`, so a shell
fixture cannot leave the process group. It re-execs itself with
`SysProcAttr{Setsid: true}`, handing the grandchild BOTH inherited output pipes,
and does not wait on it.

- **shape-A (deadline)**: spawn the escapee, write nothing, sleep past the deadline.
- **shape-B (success)**: spawn the escapee, write the whole envelope to stdout,
  exit 0 — leaving the escapee as the only remaining writer.

The envelope is a deterministic byte pattern framed `ENVELOPE-BEGIN:` ...
`:ENVELOPE-END`, reconstructed independently by the driver and compared with
`bytes.Equal`, so completeness is asserted byte-for-byte rather than by length.

A side channel on fd 3 carries the pgid evidence and a `child-write-returned`
marker, which distinguishes a child that finished its write from one still
blocked in `write(2)`.

## Escape-is-real confirmation

Both platforms, from `output.txt`:

```
darwin: child pid=37707 pgid=37707 | grandchild pid=37709 pgid=37709
linux:  child pid=828   pgid=828   | grandchild pid=834   pgid=834
grandchild alive AFTER kill(-pgid, SIGKILL)? true      (both)
```

The grandchild's pgid differs from the child's on both platforms, and it is
still alive after the group SIGKILL that `reapGroup` sends. The escape is real,
not simulated: the drain sees no EOF, and the bounded join fires on its timer.

## A3 — the FIFO / no-loss claim

Result: **holds within its scope, and the spike fixes what that scope is.**

Identical on darwin and linux:

| payload | prompt drain | adversarial delayed drain | child exit |
| ------- | ------------ | ------------------------- | ---------- |
| 1 KiB   | complete     | **complete**              | exit 0 |
| 63 KiB  | complete     | **complete**              | exit 0 |
| 128 KiB | complete     | 65536 of 131072 (prefix)  | **SIGKILL, blocked in write** |
| 1 MiB   | complete     | 65536 of 1048576 (prefix) | **SIGKILL, blocked in write** |

The adversarial ordering is the load-bearing case: the stdout drain does not
read a single byte until after `cmd.Wait()` has returned AND the group SIGKILL
has been sent, so every byte it recovers was read strictly after the writing
child was gone. For payloads that fit the pipe buffer it still recovers the
complete envelope byte-for-byte. That is the FIFO property A3 rests on,
observed directly rather than inferred.

**The >pipe-buffer case, explicitly.** At 128 KiB and 1 MiB with the drain
delayed, the child **cannot exit**: it blocks in `write(2)` once the ~64 KiB
pipe buffer fills, is still blocked when the deadline arrives, and is killed by
the group signal (`sig:killed`, and its `child-write-returned` marker never
appears). Exactly 65536 bytes — one pipe buffer — survive, as a correct prefix.

This is a real limit on A3's scope, and it is a limit on the ANTECEDENT, not a
counterexample to the claim. A3 speaks of "every byte the EXITED direct child
wrote". A child blocked in `write(2)` never exited and never finished writing,
so there is no complete answer for FIFO to have delivered. Nothing was lost that
the child had successfully written; the 65536 bytes recovered are precisely what
it had managed to put into the pipe. What the case does establish is that A3's
guarantee is conditional on the drain keeping up: if the drain is stalled past
the point where the pipe buffer fills, a child with a large answer is converted
from a success into a deadline kill. Note the prompt-drain column, where the
same 1 MiB payload completes cleanly — the shipped concurrent-drain design is
what keeps the antecedent satisfiable.

So A3's reasoning is sound where it applies: after `Wait` + `reapGroup`, a drain
blocked on an EMPTY pipe past the bound has already consumed everything the
exited child wrote, and the only remaining writer is outside the group. Refusing
at that point is a policy choice, not loss of the direct child's answer.

## A6 — the total bound

Result: **the total bound (c) HOLDS with wide margin on both platforms; the
per-leg claim about (b) is TIGHT and marginally exceeded.**

`timeout` = 1.5 s, `WaitDelay` = 500 ms, budget = 2.5 s, n = 10 per shape.

### darwin/arm64

| measure | shape | min | median | max |
| --- | --- | --- | --- | --- |
| (a) deadline -> `Wait()` returned | A | 1.349 ms | 2.472 ms | 12.231 ms |
| (b) `Wait()` -> bounded join done | A | 500.075 ms | 500.431 ms | 502.955 ms |
| (c) total elapsed                 | A | 2.0019 s | 2.0026 s | 2.0129 s |
| (a) deadline -> `Wait()` returned | B | -1.4949 s | -1.4939 s | -1.4906 s |
| (b) `Wait()` -> bounded join done | B | 500.309 ms | 501.086 ms | 501.177 ms |
| (c) total elapsed                 | B | 506.148 ms | 507.177 ms | 509.762 ms |

### linux/arm64

| measure | shape | min | median | max |
| --- | --- | --- | --- | --- |
| (a) deadline -> `Wait()` returned | A | 1.218 ms | 6.996 ms | 78.834 ms |
| (b) `Wait()` -> bounded join done | A | 500.617 ms | 501.439 ms | 503.962 ms |
| (c) total elapsed                 | A | 2.0027 s | 2.0095 s | 2.0817 s |
| (a) deadline -> `Wait()` returned | B | -1.4888 s | -1.4692 s | -1.4296 s |
| (b) `Wait()` -> bounded join done | B | 500.277 ms | 502.525 ms | 508.213 ms |
| (c) total elapsed                 | B | 512.824 ms | 531.098 ms | 575.689 ms |

### Reading the numbers

**(a) is comfortable.** In shape-A the child ignores nothing but is silent, and
`Wait` returns 1-79 ms after the context deadline — two orders of magnitude
inside one `WaitDelay`. `Wait` does not pay a `WaitDelay` here because the
direct child dies promptly to the `Cancel` group kill; the delay only bounds the
case where it does not. In shape-B, (a) is NEGATIVE (about -1.49 s): the child
writes and exits ~10 ms in, so `Wait` returns long before the deadline exists to
be measured from. "Within one WaitDelay of the deadline" is vacuously true
there, and the honest statement is that `Wait` returned early.

**(b) is tight, and is exceeded by a small margin in every single sample.**
The bounded join takes 500.1-508.2 ms against a 500 ms `WaitDelay`. It is never
under. This is structural, not noise: the join waits the full timer (the escaped
grandchild holds the pipes, so the drain never reaches EOF on its own), and then
pays timer-fire latency plus the goroutine scheduling to observe
`SetReadDeadline` and unwind. A specification that says the join is bounded by
"one more `WaitDelay`" is therefore false as a strict inequality — the true
bound is one `WaitDelay` PLUS a small scheduling epsilon (worst observed 8.2 ms
on linux, 3.0 ms on darwin).

**(c) holds with wide margin.** Worst case 2.013 s (darwin) and 2.082 s (linux)
against the 2.5 s budget — margins of 487 ms and 418 ms. The margin exists
because (a) consumes almost none of its allowance, which absorbs (b)'s small
overshoot many times over.

### Verdict on A6

The claim as a TOTAL bound is sound and was never close to being violated. The
claim as stated per-leg is imprecise in one place: the drain join is not
strictly "within one `WaitDelay`", it is one `WaitDelay` plus scheduling
latency. If A6 is to be relied on as a hard deadline rather than a budget, the
prose should say `timeout + 2*WaitDelay + epsilon`, or the bound should be
described as approximate. The generous slack in (a) means the arithmetic
conclusion survives either way.

Caveat on scope: these numbers are for a child that dies promptly to the group
kill. A child that BLOCKS termination (uninterruptible, or ignoring the signal
in a way SIGKILL cannot) would push (a) toward its full `WaitDelay`, and only
then does the 2.5 s budget become tight rather than roomy. That case was not
exercised here.

## Files

- `fixture.go` — FX-deadline-escape, the child that spawns a setsid escapee.
- `spike.go` — the driver, mirroring `cmdbind.spawn` with a bounded join.
- `run-linux.sh` — container entry point for the linux leg.
- `output.txt` — verbatim output of both runs, including the timing tables.

Both sources carry `//go:build ignore` so they are inert to the project build,
and are compiled directly by file name.

## Cleanup

Every grandchild is SIGKILLed by pid after each case, and the fixture's escapee
also self-exits on its own short sleep as a backstop. No stray sleepers remain.
