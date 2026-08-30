model: claude-fable-5
variant: full (profile: foundational)

# Repeatability reconstruction — RDR 0025 (command-invoking accessor bindings)

Read path: started from elements C1–C6, MVV, S1–S7; widened past the bare C
spans to the whole `§normative-contracts` section (C spans are ~10 lines each
inside a 389-line section whose surrounding prose carries the signature,
threading, and error-mode detail), plus `§approach`, `§load-bearing-decisions`,
`§illustrative-code`, and `§problem-statement`. Widened spans are noted
inline where they supplied a fact the contract block alone left open.

## 1. Public API

### Declaration surface (TOML, on the existing accessor entry — C1)

```toml
[read.<id> | gate.<id> | write.<id>]
command       = ["<argv0>", "<arg>", ...]   # []string; exactly one of path/command per entry
output        = "json" | "raw"              # *string; read entries only; default "json" when omitted
exit_absent   = [<code>, ...]               # []int; read entries only
exit_verdicts = { "<code>" = "<verdict>" }  # map[string]string; gate entries only
env           = { "<KEY>" = "<value>" }     # map[string]string
env_pass      = ["<VAR>", ...]              # []string
# role, keys, timeout, read_back unchanged from RDR 0002/0004
```

Placeholder vocabulary (C2): closed, tool-defined, v1 complete at exactly
`{artifact}`, substituted whole-element only with the caller-bound artifact
path for the entry's declared role. A `{...}` token that is not exactly a
known whole-element placeholder is a load-time defect, never literal text.

### Go surface

- `flowbind.Registry` — a free function (not a method), ONE production call
  site (`flow_exec.go:830` inside `buildRequest`). Gains two inputs: the
  `--allow-commands` gate value (A13) and the model base directory —
  `filepath.Dir` of the absolutized `--model` path, threaded from the CLI
  caller that read the file, never stored on `table.Model` (A12, Pending).
  GUESS at the exact signature (the RDR fixes the inputs, not the order or
  spelling): `func Registry(m *table.Model, modelDir string, allowCommands bool) (…bindings…, error)`.
- Command bindings implement the existing three interfaces in
  `internal/accessor/binding.go` — `ReadBinding`, `GateBinding`,
  `WriteBinding` (`Read`, `Gate`, `Apply`, each returning a bare `error`);
  the interfaces do NOT change.
- `accessor.Refusal` gains one additive field: `Detail string`, carrying the
  last 4 KiB of child stderr (empty on every refusal that carries no error —
  timeout, read-back, gate-off). Existing fields unchanged: Class, Accessor,
  Capability, Role, Timeout, Keys, Expected, Observed, Reason.
- `accessor.ExecError` — NEW exported typed error with one exported field
  `Detail string`; returned by the command binding through the existing
  `error` return; the executor `errors.As`-es it. GUESS: it also wraps or
  stores an underlying cause and implements `Error() string` conventionally.
- `internal/table/category.go` — six NEW `Category` constants, appended to
  the literal slice `table.Categories()` at the tail, in this order:
  `command_and_path_conflict`, `command_empty`,
  `command_unknown_placeholder`, `command_shell_interpreter`,
  `command_output_shape`, `command_env_conflict`.
  (GUESS at Go identifier spellings, e.g. `CatCommandAndPathConflict` —
  the RDR fixes the category strings and the `Cat…` family precedent, not
  the constant names.)
- `internal/table/source.go::sourceAcc` gains exactly the six fields above
  with the Go types shown; `decodeStrict` (DisallowUnknownFields) means an
  undeclared field is a document-level `CatUnknownSchemaField`.
- CLI: `--allow-commands`, a per-invocation flag, v1's ONLY opt-in — never
  the model file; registered on every verb whose command path reaches
  `flow_exec.go::buildRequest` (today four call sites across three verbs:
  `flow_next.go`, `flow_resolve.go`, `flow_state.go` ×2). `intrastate lint`
  does not carry it; lint validates regardless.
- `executor.go::refusalOf` gains an error parameter (`errors.As` →
  `*accessor.ExecError` → `Refusal.Detail`); `refusalWithKeys` wraps it and
  forwards the error — BOTH thread the parameter (A11, Pending), or the read
  path silently drops its stderr tail.
- Non-Unix: a runtime `runtime.GOOS` check in the command binding (an
  injectable package-level `goos` var defaulting to `runtime.GOOS`) refuses
  `execution_failure` before spawn; no build constraint; lint stays
  platform-neutral in the same binary.

### Error modes (refusal classes unchanged from RDR 0004)

- `execution_failure`: spawn failure (no exit code — `exec.ErrNotFound`),
  malformed or oversized (>1 MiB) stdout, unlisted non-zero exit with empty
  stdout, relative or `-`-prefixed artifact path (refused before spawn),
  carrier-less entry reaching the constructor (refusing binding, never a
  `Path: ""` file binding), gate off (`Detail` naming `allow_commands`,
  before spawn), non-Unix platform, missing base dir with a
  separator-bearing argv0 (refused at binding construction).
- `timeout`: classified from `ctx.Err()`, never from the wait error.
- Read: a declared key the JSON object omits is UNREADABLE, never
  established-absent; absence is established only by `exit_absent` (listed
  exit + empty stdout) or a key genuinely reading back absent; a value that
  reads back as the literal `<clear>` is UNREADABLE on the read path
  (`clearIsUnreadable = true`) and a mismatch on write read-back (false).
- Write: success never taken from exit status; read-back through the role's
  reader is mandatory (`read_back = true`); mismatch → `read_back_mismatch`,
  unparsable artifact → `read_back_incomplete`; applied-but-unverified sense
  (`0004:C14`) preserved when read-back cannot complete.
- Gate: verdict strings are `internal/accessor/model.go::Verdict`'s
  (`allow` | `deny` | `indeterminate`); execution failure is never laundered
  into a verdict.

## 2. Three most important internal helpers

1. **Load-time command-entry validator** (arms added to
   `internal/table/load.go::accessorTable`): enforces exactly-one of
   `path`/`command`; non-empty vector with no empty element; whole-element
   known placeholders; the open (deny-listed) interpreter check (`sh -c`,
   `bash -c`, `python -c`, env chains — unlisted interpreters admitted);
   shape rules (`raw` ⇒ exactly one key; `exit_absent` read-only;
   `exit_verdicts` gate-only, verdict values only); `INTRASTATE_` prefix ban
   on `env` keys and `env_pass` names. Reports ONE category per entry,
   fail-fast in clause order (intra-entry only; cross-entry order is
   unspecified map iteration; cross-table order is 0002's read→write→gate).
2. **Spawn/invocation helper in the command binding** (GUESS at name, e.g.
   `invoke`): performs whole-element `{artifact}` substitution with the
   absolute-path / no-leading-`-` refusal before spawn; resolves argv0
   (bare name → parent `PATH`; path separator → model file's dir; never
   cwd); composes the child env in layer order
   overlay {INTRASTATE_ROLE, INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR,
   INTRASTATE_PROTOCOL=1} > entry `env` > `env_pass` > parent allowlist
   {PATH, HOME, TMPDIR, LANG, LC_* (literal `LC_` prefix)} — later layer
   wins, nothing else inherited; writes the stdin JSON; runs under the ctx
   deadline with the A1 triple `Setpgid` + `Cancel` = SIGKILL to `-pgid` +
   `WaitDelay = 500ms`; caps stdout at 1 MiB; keeps the last 4 KiB of
   stderr for `ExecError.Detail`.
3. **Result classifier** (GUESS at name, e.g. `classifyResult`): normative
   ordering — parse the stdout envelope FIRST, then classify the exit code;
   exit maps (`exit_absent` / `exit_verdicts`) consulted only for a process
   that ran and exited with EMPTY stdout (a truncated-at-cap stdout is
   non-empty ⇒ `execution_failure`); timeout from `ctx.Err()`; raw mode
   strips exactly one trailing `"\n"`; everything unclassifiable becomes
   `*accessor.ExecError` carrying the stderr tail.

(The registry's carrier selector — `command` first, `path` second, refusing
binding on the residue, gate applied at construction — is the fourth
load-bearing helper; listed here because it is where fail-closed lives.)

## 3. Data model (persisted / boundary-crossing)

- **Model file (persisted TOML)**: accessor entries as in §1; identity stays
  RDR 0004's `(flow, name, capability)` triple — the carrier kind does not
  enter identity. The file binding's magic-suffix vocabulary (`verdictFor`,
  `unreachable`) stays unchanged for path-backed entries.
- **Child stdin (write)**: `{"<key>": "<value>", ...}` — the planned tags as
  a flat JSON object of strings; a planned clear crosses as the literal
  reserved value `<clear>`; the tool/wrapper performs the removal.
- **Child stdin (read/gate)**: `{}` — always sent; nothing per-invocation
  beyond `{artifact}` (the declared keys are already in the model).
- **Child stdout (read, json)**: flat JSON object of strings; omitted
  declared key = UNREADABLE. **(read, raw)**: the single declared key's
  value = stdout minus one trailing `"\n"`. **(gate)**:
  `{"verdict": "allow"|"deny"|"indeterminate", "reason": "<text>"}`.
- **Exit maps**: `exit_absent` (read) — listed exit + empty stdout
  establishes every declared key absent; `exit_verdicts` (gate) — listed
  exit + empty stdout is that verdict; verdict codes only.
- **Child env overlay**: `INTRASTATE_ROLE`, `INTRASTATE_CAPABILITY`,
  `INTRASTATE_ACCESSOR`, `INTRASTATE_PROTOCOL=1` (protocol version rides
  out-of-band, keeping the stdout map flat).
- **In-process**: `accessor.Refusal` + new `Detail string`;
  `accessor.ExecError{Detail string}`; `table.Category` six new members in
  `Categories()`; NO new field on `table.Model` (base dir threaded by the
  CLI caller). Rendering: one `CLIError.Detail` slot — applied-sense text
  first, stderr tail appended after.

## 4. Top-level pseudo-code — one command-backed invocation

```
invoke(entry, capability, boundArtifactPath, ctx):
    # constructed by flowbind.Registry: command→command binding, path→file
    # binding, neither→refusing binding (never Path:"")
    if not allowCommands:                        # gate sits in the constructor
        refuse execution_failure "allow_commands" (before spawn)
    if goos != unix-like:
        refuse execution_failure "unsupported platform" (before spawn)

    argv = copy(entry.command)
    for i, elem in argv:
        if elem == "{artifact}":
            if !isAbs(boundArtifactPath) or startsWith(boundArtifactPath, "-"):
                refuse execution_failure "artifact path must be absolute for a command entry"
            argv[i] = boundArtifactPath          # whole-element only; C5 caught partial tokens at load
    argv0 = resolve(argv[0]): bare → parent PATH; has separator → modelDir
            (no modelDir → refused at binding construction); never cwd

    env = allowlist{PATH,HOME,TMPDIR,LANG,LC_*} ⊕ env_pass ⊕ entry.env
          ⊕ {INTRASTATE_ROLE, INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR,
             INTRASTATE_PROTOCOL=1}              # later layer wins
    stdin = (write ? json(plannedTags incl "<clear>" literal) : "{}")

    cmd = exec.CommandContext(ctx, argv0, argv[1:]...)
    cmd.SysProcAttr.Setpgid = true
    cmd.Cancel = kill(-pgid, SIGKILL); cmd.WaitDelay = 500ms
    out, errTail, exitInfo = run(cmd, stdin, stdoutCap=1MiB, stderrTail=4KiB)

    if ctx.Err() deadline:  refuse timeout       # never from the wait error
    if spawnFailed:         refuse execution_failure (ExecError{Detail: errTail})
    if len(out) > 1MiB:     refuse execution_failure  # non-empty: never exit-mapped

    if out != "":                                # parse FIRST, exit second
        result = parseEnvelope(out)              # raw: strip one "\n"; json: flat map; gate: verdict envelope
        if malformed: refuse execution_failure (ExecError{Detail: errTail})
        return result                            # well-formed deny + non-zero exit = deny
    else:
        if read  and exit in exit_absent:    return allDeclaredKeysAbsent
        if gate  and exit in exit_verdicts:  return mappedVerdict
        if exit == 0 and read(json):         return UNREADABLE per omitted key
        refuse execution_failure (ExecError{Detail: errTail})

    # write path only, in Executor (not the binding):
    #   re-read via the role's declared reader under its own timeout (0004:C15);
    #   mismatch → read_back_mismatch; unparsable → read_back_incomplete;
    #   never success from exit status; applied-but-unverified preserved
```

GUESS ledger (where the RDR is silent and I chose): the exact Go signature
and parameter order of `flowbind.Registry`; the Go identifier spellings of
the six `Category` constants; `ExecError`'s methods beyond the exported
`Detail` field; the internal helper names (`invoke`, `classifyResult`); the
exact zero-exit empty-stdout json-read behaviour shown above (the RDR fixes
omission-is-UNREADABLE for a parsed object; an entirely empty stdout with
exit 0 and no map is, by the "exit maps apply only to empty stdout" and
"unlisted exit" clauses, most consistently an `execution_failure` for
gates — for reads I marked the UNREADABLE reading as a GUESS).
