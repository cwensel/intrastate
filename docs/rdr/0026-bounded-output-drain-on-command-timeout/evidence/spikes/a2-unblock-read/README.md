Model: claude-opus-5

# Spike: A2 — ending a pending blocking `Read` on an `os.Pipe` read end

Verifies Critical Assumption A2 of RDR 0026: a joining goroutine can end a
pending blocking `Read` on an `os.Pipe` read end from another goroutine, and the
drain goroutine returns promptly with the bytes it had read.

This is an order/reachability claim, so everything below is a measured run, not a
symbol read.

## Platforms

- darwin/arm64, go1.26.6 (Darwin 25.6.0, arm64) — run natively
- linux/arm64, go1.26.8 — run under `docker run --rm golang:1.26`

## Commands

```sh
# darwin (native)
cd <scratch>/a2spike  && go build -o a2spike . && ./a2spike
cd <scratch>/a2spike2 && go build -o a2b     . && ./a2b

# linux (docker; image golang:1.26 -> go1.26.8 linux/arm64)
docker run --rm -v <scratch>/a2spike:/w  -w /w golang:1.26 sh -c 'go build -o /tmp/a2  . && /tmp/a2'
docker run --rm -v <scratch>/a2spike2:/w -w /w golang:1.26 sh -c 'go build -o /tmp/a2b . && /tmp/a2b'
```

Sources: `a2_unblock_read.go` (S1–S6), `a2_buffered_bytes.go` (T1–T3).
Verbatim captured stdout/stderr of all four runs: `output.txt`.

## Fixture

Every scenario keeps the pipe's WRITE end held open by a `setsid(2)` grandchild
(`SysProcAttr{Setsid: true}`, write end inherited as fd 3 via `ExtraFiles`), so
the read end never sees EOF: a blocked `Read` can end only because of the
deadline or the `Close`, never because the writer went away. macOS ships no
`setsid` binary, so the holder is the spike binary re-exec'd, not a shell.

The holder signals readiness through a temp file before the measurement starts,
and refuses a non-positive lifetime. This guard exists because the first draft
passed a bare `4000` to a `time.Duration` parameter — 4000 **nanoseconds**, not
4s — so every holder exited instantly and three scenarios reported a spurious
`io.EOF`. The guard makes "the writer already died" impossible to mistake for a
result. All grandchildren are killed by process group at scenario end.

## Findings, keyed to points 1–6

1. **A blocked `Read` does return when another goroutine sets a read deadline —
   CONFIRMED, both platforms.** In S1 the reader is verified parked (still
   blocked after 300ms) with the grandchild holding the write end; then
   `SetReadDeadline(time.Now())` releases it.

2. **Latency: microseconds.** 101.5µs (darwin), 33.4µs (linux) from
   `SetReadDeadline` returning to `Read` returning. `SetReadDeadline` itself
   returned `nil`.

3. **Error: `os.ErrDeadlineExceeded`, confirmed via `errors.Is`.** Concretely
   `&fs.PathError{Op:"read", Err: os.ErrDeadlineExceeded}`, printing as
   `read |0: i/o timeout`. `errors.Is(err, os.ErrDeadlineExceeded) == true`.

4. **Buffered bytes — the load-bearing result. A deadline already in the past
   does NOT deliver buffered bytes; it returns the timeout with `n=0`.**

   - **(a) writer wrote N bytes then went quiet, write end still open, deadline
     set, read after:** `n=0`, `err=i/o timeout` — the 27 buffered bytes were
     **skipped**, on both darwin and linux (S2). The deadline is checked before
     the read is attempted, so a past deadline short-circuits.
     **This contradicts a literal reading of C1's "bytes still buffered are
     delivered" if the implementation sets the deadline to `time.Now()`.**
   - Those bytes are **DEFERRED, not LOST** (T1): after
     `SetReadDeadline(time.Time{})` clears the deadline, the very next `Read`
     returns all 27 bytes. They stay in the pipe buffer.
   - **The C1-compatible discipline is a FUTURE deadline** (T2): with
     `SetReadDeadline(time.Now().Add(750ms))` and bytes already buffered, `Read`
     #1 delivers all 27 bytes with `err=nil`, and only the following `Read` (on
     a now-empty pipe, writer alive) returns the timeout, after ~751ms.
     A short grace deadline therefore both drains what is buffered and still
     bounds the wait.
   - A drain loop that accumulates into a buffer before inspecting the error
     keeps everything it had already read (S3): 24/24 bytes preserved, loop
     ended 48.8µs (darwin) / 34.4µs (linux) after the deadline was set.
   - **(b) pipe genuinely empty with a live writer:** `n=0`,
     `err=os.ErrDeadlineExceeded`, returned in 21µs (darwin) / 2µs (linux) —
     as expected (T3, and S1).

5. **F5 / `os.ErrNoDeadline` — trigger NOT exercised on a pipe.** An `os.Pipe`
   read end always accepted the deadline (`nil`) in every run, so a natural
   non-pollable pipe was not constructed. Recorded instead:
   - a regular file returns `os.ErrNoDeadline` ("file type does not support
     deadline"), `errors.Is` true — this is the observed error *value*;
   - `os.NewFile(dup-of-pipe-read-end)` also returned `os.ErrNoDeadline`, since
     the dup'd fd is adopted without the poller registration the original had.
     This is a real, reachable way to hold a pipe fd that rejects deadlines, but
     it is not how the RDR's drain obtains its read end.

   The F5 trigger is therefore **unexercised for the actual code path**; the
   error value is confirmed, the pipe-specific registration failure is not.

6. **`Close` on the read end also unblocks, and is strictly worse for C1.**
   `Close()` returned `nil` in 20µs (darwin) / 59µs (linux); the blocked `Read`
   returned 232µs (darwin) / 104µs (linux) later with
   `&fs.PathError{Op:"read", Err: os.ErrClosed}` — `read |0: file already
   closed`, `errors.Is(err, os.ErrClosed) == true` (S4). But with 27 bytes
   sitting unread, `Close` makes them **unreachable** — no clearing step can
   recover them (S5), unlike the deadline case. **Deadline dominates `Close`**
   for a drain that must preserve partial output.

## Coverage

darwin/arm64: RUN. linux/arm64: RUN (docker, go1.26.8). Results identical on
both platforms for every point, including the point-4a byte-skipping behavior,
so that is portable Go semantics rather than a darwin artifact.

Not covered: linux/amd64 (only arm64 was available); a genuinely non-pollable
`os.Pipe` (point 5).
