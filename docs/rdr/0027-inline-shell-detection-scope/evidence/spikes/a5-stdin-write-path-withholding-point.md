Model: claude-sonnet-5

# Spike: A5 — stdin write-path withholding point

## Question

RDR 0027's A5 claims the stdin-fed-interpreter hazard (`sh -s`, bare `sh`,
`python -`, `node -`) is owned by the charted `stdin = "none" | "envelope"`
successor, on the theory that withholding the envelope from a command that
never declared it leaves a stdin-reading interpreter nothing to execute. The
record's Status already marks this Pending with the withholding premise
FALSIFIED at `internal/cli/cmdbind/cmdbind.go:207`. This spike's job (the
record's own Method line): establish what a declared-stdin command actually
receives on the write path, whether any withholding point exists before
`cmd.Stdin` is set, whether adding one is local or structural, and whether
C1 actually relies on the successor closing the axis or merely names it as
out of scope.

## Method

Files/commands read:

- `internal/cli/cmdbind/cmdbind.go` (full file, 741 lines) — the `spawn`
  function (lines 158-331), the three capability bindings (`Reader.Read`
  548-596, `Gate.Gate` 638-699, `Writer.Apply` 701-736), and `stdinObject`
  (517-544).
- `internal/table/model.go:145-174` — the `Accessor` struct definition.
- `docs/rdr/0025-command-invoking-accessor-bindings/evidence/critique/Charted.md`
  (via `find ... -iname Charted.md`) — the successor's charter text.
- `rdr inspect --records .../intrastate/docs/rdr --select 0027:A5
  --select 0027:C1 --select 0027:§approach 0027` — the record's own A5,
  C1, and Approach text (the first `rdr inspect 0027` without `--records`
  hit a different repo's default records dir and returned an unrelated
  "engine_version" RDR; re-run with the explicit `--records` flag against
  this repo's `docs/rdr` returned the correct inline-shell-detection-scope
  record).
- `grep -rn "Stdin\b" internal/table/*.go` — confirms no `Stdin` field
  exists on `table.Accessor` today.
- `kata list --status open | grep -i stdin` — confirms no open kata
  tracks the successor (exit 1, no match).

## Findings

**1. `cmd.Stdin` assignment is unconditional, one line, no existing gate.**

`internal/cli/cmdbind/cmdbind.go:207`:

```go
cmd := exec.CommandContext(ctx, argv0, argv[1:]...)
cmd.Env = childEnv(acc, name, capability)
cmd.Stdin = bytes.NewReader(stdin)
```

No branch, no config lookup, no field read precedes it. The `stdin []byte`
value is a plain function parameter to `spawn` (signature at line 158-162);
nothing between `spawn`'s entry and line 207 inspects or gates it. The
pre-spawn refusal ladder (lines 163-202: `AllowCommands`, `Unsupported()`,
carrier-less-entry, `substitute`, `resolveArgv0`) never touches stdin at
all — none of those checks is a withholding point for this axis.

**2. Every command binding — read, gate, write — sends the same JSON
envelope, unconditionally.**

All three capability bindings call `spawn` with `stdinObject(...)`:

- `Reader.Read` (line 567): `stdinObject(nil)` — a read sends the EMPTY
  object `{}`.
- `Gate.Gate` (line 658): `stdinObject(nil)` — same, empty object.
- `Writer.Apply` (line 723): `stdinObject(planned)` — the planned tags as
  a flat `map[string]string`.

`stdinObject` (lines 519-544) always JSON-encodes a `map[string]string`
(possibly empty) and always returns non-nil bytes (`{}` at minimum, never
literally empty/absent). So every declared command — regardless of
capability, regardless of whether it ever named a stdin appetite — gets a
JSON object on stdin. There is no per-entry declaration today that could
already gate this: `grep -rn "Stdin\b" internal/table/*.go` returns no
hits, and the `Accessor` struct (`internal/table/model.go:145-174`) has no
`Stdin`-shaped field among `Role, Path, Keys, Timeout, ReadBack, Command,
Output, ExitAbsent, ExitVerdicts, Env, EnvPass`.

**3. Adding a withholding point is a LOCAL change.**

`stdin []byte` reaches line 207 as an opaque parameter with exactly one
consumer in the whole function. A future `stdin = "none" | "envelope"`
field on `Accessor` would gate at that exact site:

```go
if acc.Stdin != "none" {
    cmd.Stdin = bytes.NewReader(stdin)
}
```

Nothing else in `spawn` depends on whether `cmd.Stdin` is set — the
pipe-owning stdout/stderr plumbing (`os.Pipe`, the drain goroutines, `Wait`,
`reapGroup`) is independent machinery keyed on `cmd.Stdout`/`cmd.Stderr`,
untouched by this axis. This is a one-branch, one-field addition, not a
restructuring. A5's If-wrong ("the successor never ships or withholds
nothing, and the stdin forms stay unowned admitted interpreters") is
therefore survivable as a scoping/durability downgrade, not a structural
refutation — when the successor does ship, the fix it needs is exactly
this local branch, not a rewrite of `spawn`.

**4. The charted successor's remit, per `Charted.md`, is `tee {artifact}`
— a write-corruption footgun — not a shell hazard.**

`docs/rdr/0025-command-invoking-accessor-bindings/evidence/critique/Charted.md`:

> Declared stdin appetite (`stdin = "none" | "envelope"`) on a command
> entry (diff D-13; pass A C-13, pass B C-6). The spike-proven `tee
> {artifact}` hazard destroys the user's artifact and is not statically
> detectable from argv, so no C5 arm can catch it. Withholding the envelope
> from a write command that never declared it would make the footgun
> unauthorable. Out of scope here: it changes the carrier shape C1 fixes,
> and C1's field set is this record's locked surface. Suggested successor:
> a carrier-shape RDR that adds the field and the C5 arm together.

The successor's own charter never mentions shell interpreters or stdin
scripting; its motivating case is a WRITE command silently consuming and
destroying the artifact via `tee`. RDR 0027's A5/C1 text extends this to
cover the non-shell stdin-fed interpreter forms by arguing the mechanism
(withhold the envelope) is the same regardless of motivating case — a
reasonable extrapolation, but not something the successor's charter itself
commits to. `kata list --status open | grep -i stdin` returns no matches:
no kata tracks this successor at all today.

**5. C1's scoping question: does it RELY on closure, or does it NAME the
forms as out of scope?**

C1's `out of scope, BY NAME` clause:

> an interpreter that reads its script from STDIN (`sh -s`, bare `sh`,
> `sh -es`, `python -`, `node -`) — owned by the charted `stdin = "none" |
> "envelope"` successor, whose remit is what a command entry CONSUMES and
> so covers the non-shell spellings too. The channel is the scope: any
> listed interpreter taking its code on stdin rather than as a later argv
> word is admitted, however spelled.

And the Approach section, directly:

> That routing is an OWNERSHIP claim, not a closure: the successor would
> have to introduce a withholding point, and none exists today
> (`cmdbind.go:207` sets `cmd.Stdin` unconditionally — A5).

This is explicit and load-bearing for the answer: C1 NAMES the stdin-fed
forms as out of scope and ADMITTED (the lint predicate — argv-word
matching only — passes them by construction; C1 says so directly: "any
listed interpreter taking its code on stdin ... is admitted, however
spelled"). C1 does not claim the successor closes the gap; it explicitly
calls the routing "an OWNERSHIP claim, not a closure." The predicate
itself (`command_shell_interpreter`) never reads stdin and is defined
without reference to whether a future field withholds it — C1's normative
text is self-contained and correct today regardless of whether the
successor ever ships. The record's own A5 entry reaches the same
conclusion in its Evidence field: "C1's channel-scoped wording describes
that remit rather than extending it."

## Verdict

A5 is DOWNGRADABLE, not a blocker. The withholding-mechanism premise is
falsified as a present-tense fact (`cmdbind.go:207` is unconditional, and
no `Accessor.Stdin` field exists), but C1 does not rely on that mechanism
existing — C1 names the stdin-fed forms as out-of-scope/admitted by the
predicate's own definition (argv-word matching, no stdin read) and
attributes ownership of closing the gap to a successor without asserting
that ownership is already discharged. The successor's own charter
(`Charted.md`) supports routing this concern there (same carrier-shape
change, same mechanism — a withholding point on a declared stdin appetite
field), even though its motivating case was `tee {artifact}` rather than a
shell interpreter, and even though no kata currently tracks it. Adding the
withholding point when the successor lands is a one-branch, one-field
local change at the single existing call site (line 207), not a
restructuring — so A5's If-wrong is survivable. What A5 needs is a
downgrade from "Verified" (self-referential, already rejected) to a
Pending/Accepted-with-caveat that states plainly: the successor has not
shipped, no kata tracks it, and until it does, C1 is honest that these
forms are admitted and the reviewer reads argv — which is exactly what
C1's own promise clause already says ("what a reviewer may rely on is the
predicate line and nothing more").
