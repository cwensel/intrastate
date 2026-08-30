# A1 spike — deadline bounds the invocation (Go 1.26.6, darwin/arm64)

Run: `./run.sh` → `output.txt`. Source: `main.go`.

## Quoted from `$(go env GOROOT)/src/os/exec/exec.go` (go1.26.6)

`Cancel` (lines 270-276):

```
270:	// If Cancel is set to nil, nothing will happen immediately when the command's
271:	// Context is done, but a nonzero WaitDelay will still take effect. That may
272:	// be useful, for example, to work around deadlocks in commands that do not
273:	// support shutdown signals but are expected to always finish quickly.
...
276:	Cancel func() error
```

`WaitDelay` (lines 278-303):

```
278:	// If WaitDelay is non-zero, it bounds the time spent waiting on two sources
279:	// of unexpected delay in Wait: a child process that fails to exit after the
280:	// associated Context is canceled, and a child process that exits but leaves
281:	// its I/O pipes unclosed.
...
296:	// If pipes are closed due to WaitDelay, no Cancel call has occurred,
297:	// and the command has otherwise exited with a successful status, Wait and
298:	// similar methods will return ErrWaitDelay instead of nil.
299:	//
300:	// If WaitDelay is zero (the default), I/O pipes will be read until EOF,
301:	// which might not occur until orphaned subprocesses of the command have
302:	// also closed their descriptors for the pipes.
303:	WaitDelay time.Duration
```

Stdout/Stderr and Wait — descendants holding the pipe (lines 216-220, 906-908, 917-918):

```
216:	// Otherwise, during the execution of the command a separate goroutine
217:	// reads from the process over a pipe and delivers that data to the
218:	// corresponding Writer. In this case, Wait does not complete until the
219:	// goroutine reaches EOF or encounters an error or a nonzero WaitDelay
220:	// expires.
906:			// Close the child process's I/O pipes, in case it abandoned some
907:			// subprocess that inherited them and is still holding them open
908:			// (see https://go.dev/issue/23019).
917:	// If any of c.Stdin, c.Stdout or c.Stderr are not an [*os.File], Wait also waits
918:	// for the respective I/O loop copying to or from the process to complete.
```

## Method notes

- macOS `/bin/sh -c 'sleep 30'` execs `sleep` in place (no grandchild), so the
  shell string is `sleep 30; : A1SPIKE_<S>_<pid>` to force a real fork; the
  marker sits in sh's argv for `pkill -f`.
- Survivors are found by snapshotting the child's ppid-walked tree at 200ms and
  re-checking those pids (same command line) after Wait. This catches orphans
  re-parented to launchd (ppid=1), which a pgid scan misses in the non-Setpgid
  scenarios.
- `waitErr_is_DeadlineExceeded` is false in every killed scenario: per the
  `Cancel` doc, a child that exits with a non-success status (SIGKILL) makes
  Wait return the usual `*ExitError` ("signal: killed"), not ctx.Err(). The
  caller must consult `ctx.Err()` (printed as `ctxErr=`) to classify the failure
  as a timeout. Same reason `waitDelayHit` is false in S3: ErrWaitDelay only
  replaces a nil error.
- S6 elapsed (254ms) includes the 200ms snapshot sleep; the echo itself
  completes in single-digit ms.

## Results (see output.txt)

| S | fixings | elapsed | Wait returned | tree survivors |
|---|---|---|---|---|
| S1 | none (CommandContext only) | 5000ms (watchdog) | NO — blocked on stdout pipe held by orphaned `sleep 30` | 1 (`sleep 30`, ppid=1) |
| S2 | Setpgid + group-kill Cancel + WaitDelay | 1004ms | yes | 0 |
| S3 | WaitDelay only | 1501ms (= timeout + WaitDelay) | yes, via pipe close | 1 (`sleep 30`, ppid=1) |
| S4 | S2 + 1MiB stdin to non-reading child | 1002ms | yes, no stdin-write hang | 0 |
| S5 | S2, direct child | 1002ms | yes | 0 |
| S6 | S2, fast child | ms | yes, exit 0, output captured | 0 |

## Conclusion

1. A1 holds: `context.WithTimeout` + `Setpgid` + `Cancel` = `kill(-pgid, SIGKILL)` + `WaitDelay` bounds the invocation at timeout (+WaitDelay worst case) with the whole tree dead and no stdin/stdout pipe hang (S2, S4, S5), and does not disturb a fast child (S6).
2. Omit everything (S1): Wait never returns within bound — the grandchild inherits the stdout pipe, Go reads until EOF, the CLI hangs for the grandchild's lifetime; grandchild orphaned.
3. Omit process group (S3): WaitDelay closes the pipes so Wait returns at timeout+WaitDelay, but only the direct child is killed; `sleep 30` survives as an orphan — the CLI is unblocked, the tree is not stopped.
4. Omit WaitDelay but keep group kill (S2 minus WaitDelay, not run): the group kill closes all pipe holders in the group, so it would likely still return; WaitDelay remains the necessary backstop for holders outside the group (setsid'd or double-forked descendants) and for the stdin copier.
5. Timeout classification must use `ctx.Err()`, not `errors.Is(waitErr, context.DeadlineExceeded)`; the kill surfaces as `*ExitError` "signal: killed".
