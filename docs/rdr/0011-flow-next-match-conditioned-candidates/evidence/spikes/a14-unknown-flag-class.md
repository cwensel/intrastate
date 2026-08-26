Model: claude-opus-5

# A14 — the unknown-flag class and exit code for `--all` on the flow verbs

A14 claims: "Non-registration of `--all` on `flow resolve` / `flow read-state` /
`flow set-state` yields a deterministic error class and exit code through
cobra's unknown-flag path as `internal/cli/root.go` maps it, so C2's negative is
assertable on that class rather than on a `flow-*` code."

Verified live against `2d3fb77` (worktree clean), `cobra v1.10.2` /
`pflag v1.0.10`.

## The four source conditions

Each was confirmed by reading source, not by inference from the run.

1. **No persistent `--all`/`-a` on `flow` or on root.** The only persistent flag
   in the tree is `--as`: `internal/cli/root.go:51`
   (`cmd.PersistentFlags().String(respond.FlagName, …)`), and it is the only
   `PersistentFlags()` registration outside `respond`'s lookup helper
   (`internal/cli/respond/respond.go:203`) and two test assertions. The `flow`
   group registers **no** flags of its own — `newFlowCmd`
   (`internal/cli/flow.go:36-88`) carries none, and the shared registrars
   `registerSelectionFlags` (`internal/cli/flow.go:92-97`: `--model`, `--flow`,
   `--artifact`) and `registerTagFlag` (`internal/cli/flow.go:102-105`: `--tag`)
   both write to `cmd.Flags()`, not `PersistentFlags()`. The complete per-verb
   flag surface is `--evaluate-gates` (`internal/cli/flow_next.go:79`),
   `--outcome` (`internal/cli/flow_resolve.go:72`), and `--write` / `--clear`
   (`internal/cli/flow_state.go:137,139`). No `all` name and no `a` shorthand
   anywhere: a repo-wide `grep -rn '"all"' --include='*.go' internal cmd` hits
   only `internal/resolve/guard.go:16` (`BlockAll Block = "all"`),
   `internal/table/source.go:88` (a TOML tag), and test fixtures — nothing in
   `internal/cli`.

2. **No `FParseErrWhitelist.UnknownFlags` anywhere.**
   `grep -rn 'FParseErrWhitelist\|UnknownFlags' --include='*.go' .` → no
   matches. Unknown flags therefore stay hard parse errors.

3. **No `DisableFlagParsing`.**
   `grep -rn 'DisableFlagParsing' --include='*.go' .` → no matches. Every
   command parses its flags.

4. **`TraverseChildren` unset.**
   `grep -rn 'TraverseChildren' --include='*.go' .` → no matches; it is
   false everywhere by zero value.

## The mapping path

`ExecuteAndEmit` (`internal/cli/root.go:82-95`) receives the pflag error,
fails `errors.As(err, &ce)` because a pflag error is not a `*clierr.CLIError`,
and hands it to `cobraErrorToCLIError` (`internal/cli/root.go:123-130`), which
stamps a fixed envelope:

```go
	return &clierr.CLIError{
		Code:    "command-error",
		Message: err.Error(),
		Group:   clierr.GroupUserEnv,
		Hint:    "run with `--help` to list supported commands and flags",
	}
```

That is emitted through **`respond.Fail`** (`internal/cli/root.go:94`), i.e.
the structured gateway — not cobra's own printer. `SilenceErrors: true` on
root (`internal/cli/root.go:46`) and on every verb
(`internal/cli/flow.go:52`, `flow_next.go:74`, `flow_resolve.go:65`,
`flow_state.go:65,130`) is what keeps cobra's `Error: …` text off the wire.
`respond.Fail` (`internal/cli/respond/respond.go:157-168`) branches on mode:
`clierr.EmitJSON` to **stdout** under `--as=json`, `clierr.EmitText` to
**stderr** in text mode. `Group: GroupUserEnv` maps to **exit 2** in
`clierr.ExitCodeFor` (`internal/cli/clierr/clierr.go:131-149`), which
`Execute` passes to `os.Exit` (`internal/cli/root.go:69-71`). This matches
`docs/cli-output-contract.md:110` ("Every other failure is exit 2") and its
`--as=json` clause at lines 21-28 (the `CLIError` envelope goes to stdout with
no `type` discriminator).

`primeAsFlag` (`internal/cli/root.go:100-119`) is what makes `--as=json` still
honoured on this path: it commits `--as` onto the persistent flag by scanning
argv *before* `cmd.Execute()`, so the JSON envelope is produced even though
pflag never completed a parse.

## Commands run

```sh
$ go build -o "$SCRATCH/intrastate" ./cmd/intrastate
$ printf '{"stage":"resolved","status":"draft","gate_passed":"false"}' \
    > "$SCRATCH/rdr-resolved.json"
```

The bare `--all` form is sufficient and is what is recorded below: pflag
rejects the unknown flag during `cmd.Execute()`'s parse, strictly **before**
`RunE` runs, so no verb ever reaches its own argument requirements. This is
demonstrated rather than assumed — see the control block, where the same argv
minus `--all` reaches the verb and fails with `flow-model-not-found` instead.

### Unknown flag, text mode

```sh
$ intrastate flow resolve --all      ; echo $?
$ intrastate flow read-state --all   ; echo $?
$ intrastate flow set-state --all    ; echo $?
$ intrastate flow next --all         ; echo $?
```

All four: **exit 2**, stdout **empty**, stderr identical:

```
error: command-error: unknown flag: --all
  hint: run with `--help` to list supported commands and flags
```

### Unknown flag, `--as=json`

```sh
$ intrastate flow resolve --all --as=json     ; echo $?
$ intrastate flow read-state --all --as=json  ; echo $?
$ intrastate flow set-state --all --as=json   ; echo $?
$ intrastate flow next --all --as=json        ; echo $?
```

All four: **exit 2**, stderr **empty**, stdout exactly one NDJSON line,
byte-identical across the four verbs:

```json
{"code":"command-error","message":"unknown flag: --all","hint":"run with `--help` to list supported commands and flags"}
```

The envelope **is** produced on this path. Note the shape: no `type` key (the
contract's documented asymmetry, `docs/cli-output-contract.md:21-28`), no
`param`, no `findings`, and the offending flag name is recoverable only from
the free-text `message`.

### Control — same argv WITHOUT `--all`

```sh
$ intrastate flow resolve      ; echo $?
$ intrastate flow read-state   ; echo $?
$ intrastate flow set-state    ; echo $?
$ intrastate flow next         ; echo $?
```

All four: **exit 2**, stderr:

```
error: flow-model-not-found: model selection requires exactly one of --model <path> or --flow <id>
```

So the *code* difference is attributable to `--all`; the exit code is 2 either
way, because both are `GroupUserEnv`.

### Control — fully-valid argv, with and without `--all`

```sh
$ intrastate flow next       --model models/rdr.toml --artifact rdr=$A --as=json ; echo $?
$ intrastate flow resolve    --model models/rdr.toml --artifact rdr=$A --outcome advance --as=json ; echo $?
$ intrastate flow read-state --model models/rdr.toml --artifact rdr=$A --as=json ; echo $?
$ intrastate flow set-state  --model models/rdr.toml --artifact rdr=$A --write status=draft --as=json ; echo $?
```

All four **exit 0** with a `{"type":"ok",…}` envelope (e.g. `flow resolve`
selects `"rule":"prelock"`, `"outcome":"advance"`). Re-running each of those
four with `--all` inserted: all four **exit 2** with the same
`{"code":"command-error","message":"unknown flag: --all",…}` line and no `ok`
envelope. The flag is rejected even when everything else about the request is
valid.

### Variants

```sh
$ intrastate flow resolve --all=true ; echo $?   # exit 2, "unknown flag: --all"
$ intrastate flow resolve -a         ; echo $?   # exit 2
$ intrastate flow next    -a         ; echo $?   # exit 2
```

The shorthand carries a **different message text**:

```
error: command-error: unknown shorthand flag: 'a' in -a
  hint: run with `--help` to list supported commands and flags
```

Same code (`command-error`), same exit (2), different `message`. An oracle that
matches on message substring must therefore pick the form it asserts, or match
the code alone.

## Limits, recorded honestly

- **`command-error` is not specific to `--all`.** The identical code, group,
  exit, and hint are produced for any unknown flag, any unknown subcommand,
  and any `Args` violation reaching `ExecuteAndEmit`
  (`internal/cli/root.go:123-130` is unconditional), and the bare-`flow`
  RunE (`internal/cli/flow.go:71-77`) hand-authors the same `command-error` +
  `GroupUserEnv` pair. A negative asserted on the code alone does not by
  itself prove the subject was `--all`; it must be paired with the `message`
  to name the flag.
- **The flag name is only in free text.** `param` is empty on this path, so an
  oracle wanting to name `--all` must substring-match
  `"unknown flag: --all"` — text owned by pflag, not by this repo, and hence
  liable to change on a pflag upgrade. This is the weakest link in the
  assertion.
- **Exit 2 is not discriminating on its own.** Every control above also exits
  2. Only the code+message pair separates the unknown-flag disposition from a
  verb-level refusal.
- **This is a snapshot, not an invariant.** Nothing in the build enforces the
  four conditions; a future persistent `--all` on `flow` or root would silently
  convert every one of these runs from exit 2 to a verb-level outcome. A14's
  own note that "one oracle per verb … guards a future persistent flag" is the
  right remedy, and no such oracle exists today — no test in
  `internal/cli/*_test.go` asserts the `--all` case.
- Runs were single-shot per invocation, not repeated; output was
  byte-identical across the four verbs within each mode, which is the
  determinism claim actually checked.

## Verdict

A14 holds as stated. Non-registration of `--all` on `flow resolve`,
`flow read-state`, and `flow set-state` — and on `flow next`, which today
behaves identically and is the pre-change baseline — yields one deterministic
disposition: code `command-error`, group `GroupUserEnv`, **exit 2**, message
`unknown flag: --all`, emitted through the structured `respond.Fail` gateway
(not cobra's printer), with the full JSON envelope produced on stdout under
`--as=json` and the `error:` line on stderr in text mode. The four conditions
the claim rests on are all confirmed in source. C2's negative is therefore
assertable on that class, and specifically **not** on a `flow-*` code — no
`flow-*` code is reachable on this path, because pflag fails the parse before
any `RunE` executes. The caveat the RDR should carry forward is that
`command-error` + exit 2 is a shared usage bucket: the oracle must pin the
`message` substring `unknown flag: --all` to attribute the refusal to the flag,
and that substring is pflag's text rather than this repo's.
