model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability reconstruction — run-1

Reconstructed from RDR 0025 only. Elements read: C1–C6, MVV, S2–S5, D-identity,
D-wire-byte-format, D-naming, D-selection-predicate, A2, A7, A9, A11–A14.

**Spans widened past the contracts** (recorded per instruction):
- `§load-bearing-decisions` (D-*) — C1–C6 never state the TOML→Go decode shape
  of `command`/`env`/`env_pass` beyond the normative block, nor the wire-format
  rationale; D-wire-byte-format fixes "flat map of strings, presence is the
  answer" and D-selection-predicate fixes the constructor selection order.
- `§illustrative-code` — the only place the three entry kinds appear together
  as authored TOML; needed to fix which fields co-occur per capability.
- `§approach` — C2/C3 name `{artifact}` and the stdin envelope but not that the
  command bindings implement the *existing* `ReadBinding`/`GateBinding`/
  `WriteBinding` seam; the Approach states it.
- `§minimum-viable-validation` (MVV) — the end-to-end step order (lint → gate →
  invoke → read-back) is nowhere in C1–C6 as a sequence.
- A2, A7, A9, A11–A14 — C2's argv0 resolution, C4's env allowlist and `Detail`
  threading, C6's gate site, and C1's carrier-less residue are all stated in the
  contracts as *conclusions*; the assumptions carry the signatures
  (`func Registry(m *table.Model) accessor.Registry`), the call-site counts, and
  the "If wrong" fallbacks. Without them the reconstruction would have invented
  the registry signature.
- S2–S5 — C3/C4 state the ordering rule; the scenarios supply the exhaustive
  case list (nine ordering arms, three deadline arms) that fixes the classifier's
  branch structure.

Not widened past: `§failure-modes`, `§trade-offs`, `§research-findings`,
`§alternatives-considered`, `§finalization-gate` — read only as outline titles.

---

## 1. Public API of the module

New package (Phase 2): `internal/cli/cmdbind` — a **sibling package to
`flowbind`** (the RDR says "a sibling package to `flowbind`" and names no
package path). **GUESS**: the package name `cmdbind`; the RDR fixes only that it
is a sibling of `internal/cli/flowbind`.

```go
package cmdbind

// Constructor. Returns bindings for one accessor entry, discriminated by the
// caller (the registry) having already selected `command` over `path`.
//
// baseDir is filepath.Dir of the ABSOLUTIZED model path, supplied by the CLI
// caller that opened the model file (C2/A12). May be "" for an in-memory model
// with no source file: a separator-bearing argv0 then refuses at construction.
//
// allowCommands is the C6 gate. When false the constructor returns a REFUSING
// binding — it does not return an error; the refusal is deferred to invocation
// so that lint and load stay ungated.
func NewReader(acc table.Accessor, baseDir string, allowCommands bool) accessor.ReadBinding
func NewGate(acc table.Accessor, baseDir string, allowCommands bool) accessor.GateBinding
func NewWriter(acc table.Accessor, baseDir string, allowCommands bool) accessor.WriteBinding
```

**GUESS**: the three-constructor split and their exact names/parameter order.
The RDR fixes that the constructor takes the base dir "as a parameter" and that
the flag "reaches `flowbind.Registry` as a new parameter", and that the
constructor (not the executor) is the gate site — it does not name the
functions. **GUESS**: constructors return the interface, not an error; the RDR's
"builds refusing command bindings" and "an entry with neither carrier builds a
refusing binding" both describe construction that cannot fail.

The bindings satisfy the **existing, unchanged** seam in
`internal/accessor/binding.go` (Approach; A11 confirms the bare `error` returns):

```go
type ReadBinding  interface { Read(ctx context.Context, artifact string) (map[string]string, error) }
type GateBinding  interface { Gate(ctx context.Context, artifact string) (accessor.Verdict, string, error) }
type WriteBinding interface { Apply(ctx context.Context, artifact string, tags map[string]string) error }
```

**GUESS** — the method *signatures*. The RDR fixes only the method **names**
(`Read`, `Gate`, `Apply`), that they return a bare `error`, and that the
artifact path is caller-bound per role. Parameter lists, the read return type,
and the gate's two-value verdict/reason return are reconstructed from C3's
envelope shapes (a gate answers `verdict` + `reason`; a read answers a flat
map of strings; a write receives the planned tags).

### Error type (C4, normative — the one API addition this RDR fixes exactly)

```go
package accessor

// ExecError is the typed carrier for a command binding's diagnostic tail.
// Returned through the existing bare `error` return; no interface changes.
type ExecError struct {
    Detail string // last 4 KiB of the child's stderr
}
func (e *ExecError) Error() string
```

C4 fixes: one exported `Detail string` field, `*accessor.ExecError`, returned
through the existing `error` return. **GUESS**: the `Error() string` method body
and whether `ExecError` wraps an inner error (C4 says `errors.As`, which needs
no wrapping).

### Modified existing API

```go
// internal/accessor/model.go — ADDITIVE field, normative (C4)
type Refusal struct {
    Class, Accessor, Capability, Role string
    Timeout time.Duration
    Keys []string
    Expected, Observed, Reason string
    Detail string // NEW: bounded stderr tail; EMPTY on every refusal carrying no error
}

// internal/accessor/executor.go — BOTH constructors gain one error parameter (C4, A11)
func refusalOf(..., err error) Refusal        // errors.As-es *ExecError into Detail
func refusalWithKeys(..., err error) Refusal  // wraps refusalOf, forwards err

// internal/cli/flowbind/registry.go — signature change on a FREE function (C6, A13)
func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry
```

**GUESS**: the `Refusal` field types beyond `Detail string` (C4 lists the field
*names* on `main`, not their types) and the exact parameter lists of `refusalOf`
/ `refusalWithKeys` apart from the appended `err error`. **GUESS**: `Registry`'s
new parameter order; A13 fixes the current signature
`func Registry(m *table.Model) accessor.Registry` and that both the base dir and
the flag arrive as parameters.

### Carrier fields (C1, normative — exactly six, with Go types named)

```go
// internal/table/source.go
type sourceAcc struct {
    // ...existing: role, keys, timeout, read_back, path...
    Command      []string          `toml:"command"`
    Output       *string           `toml:"output"`        // read entries only; pointer so omitted ≠ explicit "json"
    ExitAbsent   []int             `toml:"exit_absent"`   // read entries only
    ExitVerdicts map[string]string `toml:"exit_verdicts"` // gate entries only
    Env          map[string]string `toml:"env"`
    EnvPass      []string          `toml:"env_pass"`
}
```

C1 fixes all six names, their Go types, and that `output` is a pointer.
**GUESS**: the Go field names and struct tags; C1 gives the TOML spellings only.
Decoding is strict (`source.go::decodeStrict`, `DisallowUnknownFields`), so an
undeclared field is a document-level `CatUnknownSchemaField`, not a C5 category.

### Error modes

**Load-time** — six new `internal/table/category.go::Category` constants,
appended in this order to `table.Categories()` after the existing members (C5,
normative):

| category | condition |
|---|---|
| `command_and_path_conflict` | both or neither of `path`/`command` |
| `command_empty` | empty vector or empty argv element |
| `command_unknown_placeholder` | unknown or non-whole-element `{…}` token |
| `command_shell_interpreter` | argv0 + inline-code flag (`sh -c`, `bash -c`, `python -c`, `env` chains) — OPEN deny-list |
| `command_output_shape` | `output="raw"` with `len(keys) != 1`; `exit_absent` on a non-read; `exit_verdicts` on a non-gate or naming a non-verdict |
| `command_env_conflict` | an `env` key or `env_pass` name matching the reserved `INTRASTATE_` prefix |

Precedence: **within one entry**, fail-fast in the order above (one categorized
error per document, `0002:C24`). **Across entries in one table**: unspecified
(Go map range) — no test may assert it. **Across tables**: read → write → gate,
first error returns (inherited from 0002).

**Runtime** — RDR 0004's refusal *classes* unchanged:
- `ClassTimeout` — classified from `ctx.Err()`, never from the wait error.
- `ClassExecutionFailure` — spawn failure (no exit code at all), non-zero exit
  not listed in the entry's exit map, malformed stdout envelope, stdout over
  1 MiB, relative or `-`-prefixed substituted artifact path (before spawn),
  separator-bearing argv0 with no base dir (before spawn), non-Unix `GOOS`
  (before spawn), C6 gate unset (before spawn), carrier-less entry.
- `read_back_mismatch` / `read_back_incomplete` — from the role's reader (0004).
- Gate `deny` / `indeterminate` are verdicts, **not** errors (C3/C4 ordering).

---

## 2. The three most important internal helpers

### `buildArgv(acc, artifact, baseDir) ([]string, error)` — C2's authority bound

Whole-element placeholder substitution over the declared vector: an element
**exactly equal** to `{artifact}` is replaced by the caller-bound artifact path
for the entry's declared role; every other element is passed byte-for-byte.
Refuses (`*ExecError`, before spawn) when the substituted path is relative or
begins with `-` — it **refuses rather than absolutizes**, because
`flow_input.go::parseArtifacts` stores `--artifact` verbatim and intrastate's
cwd is not the tool's frame of reference. Resolves argv0: a separator-free name
is left for `LookPath` against the parent's `PATH` at spawn; a name containing a
path separator is joined against `baseDir`; **never** falls back to process cwd
(the `exec.ErrDot` reversal). Empty `baseDir` + separator-bearing argv0 refuses
at construction. **GUESS**: the function name and that substitution and argv0
resolution live in one helper rather than two.

### `childEnv(acc, role, capability, name) []string` — C4's env allowlist

Composes the child environment in exactly four layers, **later wins**:
1. parent allowlist: `PATH`, `HOME`, `TMPDIR`, `LANG`, and any `LC_*`
   (literal prefix match on `LC_` — the *one* prefix rule);
2. `env_pass` named whole variables lifted from the parent (no globbing);
3. the entry's literal `env` table;
4. the overlay `INTRASTATE_ROLE`, `INTRASTATE_CAPABILITY`,
   `INTRASTATE_ACCESSOR`, `INTRASTATE_PROTOCOL=1`.

Nothing else is inherited. Overlay shadowing is impossible: an `env` key or
`env_pass` name with the `INTRASTATE_` prefix is a load-time
`command_env_conflict`. **GUESS**: the function name and signature.

### `run(ctx, argv, env, stdin) (stdout []byte, stderrTail string, exitCode int, err error)` — C4's deadline triple

The A1 spike's necessary-and-sufficient triple, with no injection seam (S4 is
explicit that adding one would put a way to weaken the deadline into the
shipping binding): `SysProcAttr.Setpgid = true`; `Cmd.Cancel` sends `SIGKILL` to
`-pgid` at the ctx deadline; `Cmd.WaitDelay = 500ms` bounds both the stdin write
and the post-kill pipe drain. Caps stdout at 1 MiB (overflow ⇒
`execution_failure` from the bound, and the truncated stdout is *non-empty*, so
it is never exit-mapped). Retains the last 4 KiB of stderr. Reports the deadline
from `ctx.Err()`, never from the `*exec.ExitError` the kill produces. Refuses
before spawn when the injectable package-level `goos` var (default
`runtime.GOOS`) is not Unix — a runtime check, not a build constraint, so the
refusal is observable and lint stays platform-neutral in one binary.
**GUESS**: the function name and return tuple.

Honourable mention (not in the top three but normative): `classify` — the
envelope-before-exit ordering, described in §4 below.

---

## 3. Data model across the boundary

### Declared (persisted in the TOML model, read by humans and lint)

```toml
[read.<id> | gate.<id> | write.<id>]
role          = "<role>"
keys          = ["<key>", ...]
timeout       = "<duration>"
read_back     = true                       # mandatory on every write (0004)
command       = ["<argv0>", "<arg>", ...]  # []string; EXACTLY ONE of path|command
output        = "json" | "raw"             # *string, read only; default "json"
exit_absent   = [<code>, ...]              # []int, read only
exit_verdicts = { "<code>" = "<verdict>" } # map[string]string, gate only
env           = { "<KEY>" = "<value>" }    # map[string]string
env_pass      = ["<VAR>", ...]             # []string
```

Identity is **unchanged**: `(flow, name, capability)` (0004). Carrier kind does
not enter identity.

### Crossing the process boundary (per invocation, ephemeral — nothing persisted)

**stdin — always sent, in every capability:**
- write: `{"<key>": "<value>", ...}` — the planned tags. A planned `<clear>`
  crosses as the **literal reserved value**; the tool or wrapper performs the
  removal and read-back verifies absence (`0004:C11`).
- read/gate: `{}` — an empty object. Nothing per-invocation beyond `{artifact}`.
- The requested key set is deliberately **not** on stdin (D-wire-byte-format):
  a list is not a string value, and the keys are already in the model.

**stdout:**
- read, `output = "json"` (default): flat JSON object of strings. A declared key
  the object **omits** is `UNREADABLE`, never established-absent — the two are
  different answers in `executor.go::classify` and only absence is unsafe to
  guess.
- read, `output = "raw"`: the single declared key's value is stdout minus **one**
  trailing `"\n"`. Valid only when `keys` has exactly one entry.
- gate: `{"verdict": "allow"|"deny"|"indeterminate", "reason": "<text>"}` —
  the verdict strings are `internal/accessor/model.go::Verdict`'s (`0004:C9`).
- write: stdout is not a result channel; success is **never** taken from exit
  status. Verification is read-back through the role's declared reader.

**Exit code**, consulted **only** when stdout is empty and only for a process
that ran and exited:
- read: a code in `exit_absent` establishes **every** declared key absent.
- gate: a code in `exit_verdicts` yields that verdict.
- anything else — unlisted code, spawn failure (no exit code exists, so no map
  entry can match it however written), malformed stdout — is
  `execution_failure`. Execution failure is never laundered into a verdict.

**Environment** (out-of-band, one direction): the four-layer allowlist of §2.
The protocol version rides here as `INTRASTATE_PROTOCOL=1`, deliberately keeping
the stdout map flat.

**Back to the caller:** `*accessor.ExecError{Detail}` (≤4 KiB stderr tail) →
`accessor.Refusal.Detail` → `internal/cli/clierr.CLIError.Detail`, which has
exactly **one** `Detail string` slot rendering one `detail:` line. Where both
senses are present the **applied-but-unverified** text leads and the stderr tail
is appended after it — losing "the write may have applied" would drop the more
consequential fact.

**Nothing new is persisted.** The model file is the only durable artifact this
RDR touches, and only by adding six declared fields.

---

## 4. Pseudo-code of the main operation

```
# One command-backed invocation, from the registry to the classified outcome.
# C1 selection, C6 gate, C2 argv, C4 spawn/deadline/env, C3+C4 classification.

Registry(model, baseDir, allowCommands):                     # C6/A13: ONE production site
  for each accessor entry acc:
    if acc.command non-empty:  binding = cmdbind.New*(acc, baseDir, allowCommands)
    elif acc.path non-empty:   binding = flowbind.New*(acc)          # unchanged
    else:                      binding = refusingBinding(execution_failure,
                                           Detail="malformed entry: no carrier")
                                                             # C1/A14: never Path:""

invoke(ctx, binding, artifact, plannedTags):
  if not allowCommands:  return ExecError("execution refused: allow_commands")   # before spawn
  if goos not in UNIX:   return ExecError("command entries are Unix-only")       # before spawn

  argv = declared vector, each element == "{artifact}" replaced whole by artifact
  if argv0 contains a path separator:
      if baseDir == "": return ExecError("no model dir for separator-bearing argv0")
      argv[0] = join(baseDir, argv[0])                       # never cwd
  if artifact was substituted and (not absolute or starts with "-"):
      return ExecError("artifact path must be absolute for a command entry")

  stdin = plannedTags if WRITE else {}                       # C3: always sent
  env   = allowlist(PATH,HOME,TMPDIR,LANG,LC_*) + env_pass + acc.env + INTRASTATE_*

  ctx = WithTimeout(ctx, acc.timeout)                        # 0004's declared bound
  cmd = exec.CommandContext(ctx, argv...); cmd.Env = env
  cmd.SysProcAttr.Setpgid = true
  cmd.Cancel    = func { kill(-pgid, SIGKILL) }              # the TREE, not the child
  cmd.WaitDelay = 500ms                                      # bounds stdin write + drain
  stdout, stderrTail, exitCode, waitErr = run(cmd, cap=1MiB, tail=4KiB)

  if ctx.Err() == DeadlineExceeded: return Refusal{Timeout}  # from ctx, NOT waitErr
  if spawn failed (no exit code):   return ExecError(stderrTail)   # never exit-mapped
  if stdout over 1 MiB:             return ExecError("stdout exceeds 1 MiB")

  # ---- C4 ordering is NORMATIVE: parse the envelope, THEN the exit code ----
  if stdout non-empty:
      if READ and output=="raw":  return {theOneKey: trimOneTrailingNewline(stdout)}
      parsed = parseFlatJSONObjectOfStrings(stdout)
      if parse fails:             return ExecError(stderrTail)     # whatever the exit
      if GATE:  return parsed.verdict, parsed.reason              # deny beats non-zero exit
      if READ:  return parsed        # a declared key OMITTED is UNREADABLE, not absent
  else:                                                            # empty stdout only
      if READ and exitCode in acc.exit_absent:      return allDeclaredKeysAbsent
      if GATE and exitCode in acc.exit_verdicts:    return that verdict
      return ExecError(stderrTail)                                 # incl. exit 0 w/ no envelope

# Write completion (0004:C12, unchanged): exit status decides NOTHING.
apply(...):
  invoke(WRITE, plannedTags)                                  # <clear> crosses literally
  reader = readerFor(role)                                    # own bounded timeout, 0004:C15
  observed = reader.Read(...)                                 # clearIsUnreadable = FALSE here
  if observed != plannedTags:  return read_back_mismatch
  if read-back could not complete: return applied-but-unverified (0004:C14)
```

**GUESS**: that reads under `output = "raw"` short-circuit before JSON parsing
(C3 implies it — raw mode has no envelope — but no clause orders the two).
**GUESS**: that an empty stdout with exit 0 and no exit-map hit is
`execution_failure` for a read; C3 states it explicitly only for `output="raw"`
("an empty stdout with an unlisted exit is `execution_failure`, not an empty
string"), and by C4's "any unlisted exit... is execution_failure" it generalizes.
**GUESS**: `WaitDelay = 500ms` applies to the stdin write and the drain as one
budget rather than two (C4 says "bounds the stdin write and the pipe drain").
**GUESS**: the refusal for a carrier-less entry is deferred to invocation rather
than raised at construction; C1 says the constructor "builds a refusing binding",
which implies deferral.
