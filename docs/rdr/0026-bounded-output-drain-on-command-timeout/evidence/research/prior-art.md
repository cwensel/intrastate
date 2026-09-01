Model: claude-fable-5

# Prior art — bounded output drain on command timeout (Stage 2 read, not a spike)

Problem class: a parent reads a child's stdout/stderr over inherited pipes;
a descendant the parent's termination cannot reach keeps a write end open,
so the read never sees EOF. Question this RDR owns: when "bounded wait" and
"whole read" conflict, which wins, and what is reported.

## Accepted citations (load-bearing; quoted from source)

### Go `os/exec` — the instance the code is built on

`$(go env GOROOT)/src/os/exec/exec.go` (go1.26), `Cmd.WaitDelay` doc:

> If WaitDelay is non-zero, it bounds the time spent waiting on two sources
> of unexpected delay in Wait: a child process that fails to exit after the
> associated Context is canceled, and a child process that exits but leaves
> its I/O pipes unclosed.
> ...
> If pipes are closed due to WaitDelay, no Cancel call has occurred,
> and the command has otherwise exited with a successful status, Wait and
> similar methods will return ErrWaitDelay instead of nil.
>
> If WaitDelay is zero (the default), I/O pipes will be read until EOF,
> which might not occur until orphaned subprocesses of the command have
> also closed their descriptors for the pipes.

`exec.go::awaitGoroutines` doc + body:

> If c.WaitDelay elapses before the goroutines complete, awaitGoroutines
> forcibly closes their pipes and returns ErrWaitDelay.

    case <-timer.C:
        closeDescriptors(c.parentIOPipes)
        ...
        return ErrWaitDelay

`exec.go::ErrWaitDelay`:

> ErrWaitDelay is returned by [Cmd.Wait] if the process exits with a
> successful status code but its output pipes are not closed before the
> command's WaitDelay expires.

⇒ Go's own precedence: the bound outranks the whole read; the PARENT closes
its read ends to end the wait; a success-status child whose pipes stayed
open is reported as an ERROR (`ErrWaitDelay`), never as a silent success.
This only reaches pipes `Cmd` owns; `cmdbind::spawn` hands `*os.File`
ends to `cmd.Stdout`/`cmd.Stderr`, which `Wait` neither copies nor closes
(doc: "If any of c.Stdin, c.Stdout or c.Stderr are not an [*os.File], Wait
also waits for the respective I/O loop"), so the owned drain sits outside
`WaitDelay`'s reach today.

### Go `os` — can a blocked pipe read be ended by the joining goroutine?

`$(go env GOROOT)/src/os/file_unix.go::newFile`:

    pollable := kind == kindOpenFile || kind == kindPipe || kind == kindSock || nonBlocking

`os/pipe2_unix.go::Pipe` / `os/pipe_unix.go::Pipe` construct both ends with
`kindPipe`. `os/file.go`:

> Only pollable files support [File.SetDeadline], [File.SetReadDeadline],
> and [File.SetWriteDeadline].

`internal/poll/fd_poll_runtime.go::(*pollDesc).evict`:

> Evict evicts fd from the pending list, unblocking any I/O running on fd.

called from `internal/poll/fd_unix.go::(*FD).Close`.

⇒ Both mechanisms (`SetReadDeadline` on the read end; `Close` of the read
end from the joiner) are documented to end a pending `Read` on an
`os.Pipe` file on Unix. Which one `spawn` uses is a mechanism choice; that
either ENDS the read in practice (and what error the drain sees) is an
order/reachability claim — demoted to a Resolve spike (A2), not rested on.

### Pipe semantics (class frame) — Kerrisk, *The Linux Programming Interface* §44

DevRef corpus, semantic query "read from a pipe returns end-of-file only
when all processes have closed the write end" (top hit, score 0.739):

> The process reading from the pipe closes its write descriptor for the
> pipe, so that, when the other process completes its output and closes
> its write descriptor, the reader sees end-of-file (once it has read any
> outstanding data in the pipe).

⇒ EOF is a property of the WRITER SET, not of the direct child: any
descendant holding a write end defers EOF. It also gives the FIFO fact
behind the "partial output" question — outstanding data is delivered
before EOF, so a read that is BLOCKED (empty pipe, writer alive) has
already consumed everything the exited direct child wrote. Whether the
implementation observes this ordering is a spike claim (A3), not a quote.

### Peer CLIs — the instance question (what do comparable Go CLIs do?)

Literal sweep `WaitDelay` over the Go/CLI checkout set (gh-cli, helm,
goreleaser, kubebuilder, roborev, beads, ...). Hits outside vendored code:
roborev and beads only; gh-cli / helm / goreleaser / kubebuilder set none
(they read child output with `Cmd`-owned pipes and no drain bound).

- `roborev/internal/agent/process.go::configureSubprocess`:
  `cmd.WaitDelay = subprocessWaitDelay` (5 s); the caller maps
  `errors.Is(runErr, exec.ErrWaitDelay)` to the context error — i.e. a
  bound-expired drain is reported as the timeout it followed.
- `beads/internal/storage/dolt/store.go::cliExecWaitDelay` (10 s):
  > CommandContext kills only the direct dolt child; a grandchild (e.g. a
  > cloud credential helper) that inherited the output pipes would
  > otherwise keep Wait blocked indefinitely after the kill.

⇒ Both peers that met this exact failure took Go's precedence unchanged:
bound wins, parent closes, the caller sees an error. Neither preserves a
"whole read" promise past the bound; neither widens the kill reach.

## Rejected branches

- Session-level / subreaper / cgroup reaping to make the whole read always
  terminate: no portable Unix primitive (`PR_SET_CHILD_SUBREAPER` and
  cgroups are Linux-only; macOS has neither); not searched further — the
  RDR's platform set is "every Unix with process groups" (`0025:C4`
  platform clause), so a Linux-only widening cannot be the contract.
- StateMachineRes / StateMachineLit corpora: not queried — the problem is
  process/pipe lifecycle, not state-machine design; no claim rests on them.

## Queries run (budget: ≤3 per claim, ≤5 opened hits)

1. `arc search semantic --corpus DevRef --limit 5 --json "child process orphan grandchild inherits stdout pipe parent read blocks until all writers close timeout"` — hits: TLPI §44 (pipes), §34 (orphaned process groups); frame confirmed, no quote taken.
2. `arc search semantic --corpus DevRef --limit 3 --json "read from a pipe returns end-of-file only when all processes have closed the write end"` — hit 1 quoted above.
3. `rg -c WaitDelay <langref checkouts> --glob '*.go'` — literal sweep for the peer instance question; roborev + beads read.
4. Go source opened directly: `os/exec/exec.go` (WaitDelay, ErrWaitDelay, awaitGoroutines), `os/file_unix.go::newFile`, `os/file.go` (SetDeadline docs), `internal/poll/fd_unix.go`, `internal/poll/fd_poll_runtime.go`.
