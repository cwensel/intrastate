Model: claude-opus-5

# Spike: A9 (the per-read deadline re-arm) and A10 (the pollability probe)

Two Critical Assumptions of RDR 0026, run live on both target platforms.

## What was run

```
# darwin (host)
cd <scratch>
go build -o fixture fixture.go
go build -o a9      a9_regrace.go
go build -o a10     a10_probe.go
./a9 ./fixture
./a10

# linux (container)
docker run --rm -v <this-dir>:/w:ro -w /w golang:1.26 sh /w/run-linux.sh
```

`run-linux.sh` copies the three sources into the container, builds them with
the container's toolchain, and runs the same two drivers. Verbatim output of
every run is in `output.txt`.

## Platforms

| leg    | platform     | toolchain                |
| ------ | ------------ | ------------------------ |
| darwin | darwin/arm64 | go1.26.6                 |
| linux  | linux/arm64  | go1.26.8 (`golang:1.26`) |

Both legs ran for real. Neither is UNRUN. Measured pipe buffer: 65536 bytes
on both.

## DrainGrace

`DrainGrace = 50ms`, the RDR's proposed constant: the short FUTURE deadline
re-armed before each read of the final drain. The spike mirrors it as
`DrainGraceMS = 50`.

---

## A9 — the re-arm bounds the idle gap, not the tail

### The shape that matters

The first version of this spike used only a "burst" child — write 1 MiB in
one `Write`, exit. **That produced a null result**, and it is recorded in
`output.txt` rather than discarded, because the null is itself instructive:

A 1 MiB tail from a burst writer drains in **1.6 ms on darwin, 9.5 ms on
linux** — 32 reads, none of them waiting. Every byte is already queued (or in
a parked `write(2)` that completes the instant the drain makes room). Under
that shape the re-armed arm and the single-absolute-deadline control are
**indistinguishable: both recover all 1 MiB and both end on EOF.**

So the danger to a single absolute deadline is **not tail SIZE**. It is
tail DURATION. A big tail that is already sitting in the pipe costs
milliseconds. The assumption "1 MiB is a lot of bytes, therefore it will
stress a 50 ms deadline" is false, and the burst rows disprove it.

The shape with teeth is a child that **emits output as it works** — the
realistic case. The `paced` fixture mode writes the payload in 32 chunks with
a 20 ms idle gap between them:

- each individual gap is 20 ms = **0.4x DrainGrace** (a re-armed deadline
  clears every one of them);
- the gaps SUM to ~620 ms = **12x DrainGrace** (a single absolute deadline
  set once at the start cannot survive them).

That is precisely the gap-vs-total distinction the claim is about.

### Result — CONFIRMED, both OSes

n=12 trials at 1 MiB, drain started 200 ms late (past the point a real
timeout timer would have fired):

| OS     | arm          | complete | EOF   | deadline | fewest bytes    | reads | max per-read GAP | total elapsed (max) |
| ------ | ------------ | -------- | ----- | -------- | --------------- | ----- | ---------------- | ------------------- |
| darwin | **re-armed** | 12/12    | 12/12 | 0/12     | 1048576/1048576 | 32    | 26.806ms         | 611.708ms = 12.2x   |
| darwin | absolute     | 0/12     | 0/12  | 12/12    | 163840/1048576  | 5     | 22.021ms         | 51.659ms            |
| linux  | **re-armed** | 12/12    | 12/12 | 0/12     | 1048576/1048576 | 32    | 31.569ms         | 662.905ms = 13.3x   |
| linux  | absolute     | 0/12     | 0/12  | 12/12    | 163840/1048576  | 5     | 24.044ms         | 51.908ms            |

The decisive pair of numbers, on both OSes:

- **max per-read gap ~27–32 ms < DrainGrace 50 ms** — the grace bounds THIS.
- **total elapsed ~612–663 ms = 12–13x DrainGrace** — the grace does NOT
  bound this, and the drain still reached EOF with all 1048576 bytes.

Completeness is `bytes.Equal` against an independently reconstructed envelope,
not a length check. The envelope is framed `ENVELOPE-BEGIN:` / `:ENVELOPE-END`
so a truncation anywhere, not just at the tail, would show.

### The control's failure point

The single absolute deadline fails at **exactly the same place every time on
both OSes**: 163840 bytes of 1048576 (15.6%), after 5 reads, at ~51 ms. The
terminal error is `os.ErrDeadlineExceeded` — a **false held-pipe report** on a
pipe whose writer is alive, well, and about to deliver the other 84%.

Per payload size, the control truncates at (darwin and linux identical):

| payload  | recovered by control | terminal |
| -------- | -------------------- | -------- |
| 1 KiB    | 384                  | deadline |
| 64 KiB   | 24576                | deadline |
| 256 KiB  | 90112                | deadline |
| 1 MiB    | 163840               | deadline |

Note the 1 KiB row: the control loses 62% of a **one-kilobyte** payload. Size
was never the variable.

### One flake, and what it was

RUN 3 (linux) showed `burst`/`absolute` failing **1 of 12** — 393216 bytes,
59.681 ms. That is the control's 50 ms window being blown by container
scheduling jitter on a tail that normally drains in 9 ms. RUN 4 repeated the
1 MiB block three more times: burst/absolute came back 12/12 in all three, so
the failure rate is 1 in 48 and the cause is jitter, not mechanism.

It is worth keeping, because it shows the single absolute deadline is fragile
even on the fast path — it can false-report on a tail it should trivially
handle, whenever the host hiccups.

**The re-armed arm did not flake in any configuration, on either OS, across
48+ trials per OS.**

---

## A10 — the pollability probe is safe and detecting

The probe, run once at pipe creation:

```go
armErr := rd.SetReadDeadline(time.Now().Add(probeWindow))
_ = rd.SetReadDeadline(time.Time{})   // clear immediately
```

### (a) DETECTING — CONFIRMED, both OSes

Both non-pollable shapes are detected, with an identical error value on
darwin and linux:

| shape                                                | arm error                             | `errors.Is(err, os.ErrNoDeadline)` |
| ---------------------------------------------------- | ------------------------------------- | ---------------------------------- |
| regular file (`os.Open` of a temp file)              | `"file type does not support deadline"` | **true**                           |
| `os.NewFile` over `syscall.Dup` of a pipe read end   | `"file type does not support deadline"` | **true**                           |
| ordinary `os.Pipe` read end (contrast)               | `<nil>`                               | false                              |

Error type is `*errors.errorString` in every failing case. `os.ErrNoDeadline`
declares itself as `"file type does not support deadline"`.

The dup'd-pipe shape is the load-bearing one: the fd is a **genuine pipe**,
but `syscall.Dup` + `os.NewFile` bypasses the runtime poller registration that
`os.Pipe` performs. So the probe is detecting *poller registration*, not
*file kind* — which is the property the design actually needs.

Both non-pollable files still **read correctly** after the failed probe
(21 bytes from the regular file, 15 through the dup): a failed probe costs
nothing but the information it returns.

Note: the **clear** call also returns `ErrNoDeadline` on a non-pollable file.
An implementation may check either call's error; checking the arm is the
clearer expression of intent.

### (b) SAFE / NO-OP — CONFIRMED, both OSes

The weak version of this test — round-trip a payload and see it arrive —
would pass even if the clear had silently left a deadline armed, because the
bytes were already there. So the strong test is run instead.

After probe+clear, a blocking `Read` on an **empty** pipe with a **live**
writer must still park indefinitely. A leaked 10 ms probe deadline would fire
**30x sooner** than the 300 ms observation window.

n=12 trials each:

| OS     | arm                  | still parked at 300ms | woke on write | correct bytes | errors | worst wake latency |
| ------ | -------------------- | --------------------- | ------------- | ------------- | ------ | ------------------ |
| darwin | WITH probe+clear     | 12/12                 | 12/12         | 12/12         | 0/12   | 56µs               |
| darwin | WITHOUT (baseline)   | 12/12                 | 12/12         | 12/12         | 0/12   | 45µs               |
| linux  | WITH probe+clear     | 12/12                 | 12/12         | 12/12         | 0/12   | 799µs              |
| linux  | WITHOUT (baseline)   | 12/12                 | 12/12         | 12/12         | 0/12   | 332µs              |

The probed and unprobed arms are behaviourally identical: the read parks, and
wakes only when bytes actually arrive.

Full-payload round-trip, `bytes.Equal`, n=12 per cell — **12/12 identical in
every cell, on both OSes**, with and without the probe, at 1 KiB / 64 KiB /
1 MiB.

---

## Verdict

Both claims **CONFIRMED on both OSes**. The two platforms agree on every
point; there is no disagreement to report.

The one caveat worth carrying into the RDR is a scoping correction, not a
refutation: **the re-arm earns its keep against slow-arriving tails, not
large ones.** A large tail already resident in the pipe drains in
milliseconds and would survive a single absolute deadline. What breaks the
single deadline is elapsed delivery time — a child that emits output as it
works, which is the ordinary case.

## Files

- `fixture.go` — the child: `tail` (burst) and `paced` (chunked with idle
  gaps) modes, writing a framed, independently reconstructible envelope.
- `a9_regrace.go` — the A9 driver: re-armed vs single-absolute drain loops,
  both child shapes, 4 payload sizes, n=12 trials at 1 MiB.
- `a10_probe.go` — the A10 driver: detection on 2 non-pollable shapes, and
  the parked-read + round-trip no-op safety tests.
- `run-linux.sh` — container build-and-run for the linux leg.
- `output.txt` — verbatim output of all four runs.
