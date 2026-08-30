model: claude-sonnet-4-5-20250929 (Claude Sonnet 5)
variant: full (profile: foundational)

Widened past: `§approach` (368-408) for architecture framing (existing
`flowbind`/`ReadBinding`/`GateBinding`/`WriteBinding` seam, RDR 0004 executor
reuse); `§problem-statement` (28-54) for the read/write asymmetry this closes
and the call-site file (`internal/cli/flow_state.go`); `§illustrative-code`
(860-887) for concrete TOML shape since C1-C6 give grammar/rules but only
the illustrative block shows a filled-in `[write.stage]`/`[read.*]`/`[gate.*]`
entry with wrapper-separator semantics; `§capability-dependencies` (889-898)
for what's available vs. deferred (RDR 0016 reader-uniqueness is NOT yet
landed — `readerFor` is first-match, an explicit open hazard, not a GUESS);
`§prerequisites` (1268-1277) confirming A11-A14 are still open at time of
this record. All S-scenarios (S1-S7, including 5b/6b) read in full since the
C-clauses alone underdetermine ordering/oracle details the scenarios pin
down explicitly (e.g., the nine-way exit-code table, the two-reader arity
split). Package/file names for pieces C1-C6 don't state explicitly (e.g.
exact struct field names beyond what's quoted) are GUESSed and marked below.

# Reconstruction: command-invoking accessor bindings (RDR 0025)

## 1. Public API

Package: `internal/accessor` (existing seam from RDR 0004), extended;
new command-carrier support lives beside the existing path-carrier
(`internal/cli/flowbind`).

```go
// internal/accessor/binding.go (existing interfaces, unchanged signatures —
// command bindings implement them, they do not add new ones)

type ReadBinding interface {
    Read(ctx context.Context, artifact string) (map[string]string, error)
}

type GateBinding interface {
    Gate(ctx context.Context, artifact string) (Verdict, error)
}

type WriteBinding interface {
    Write(ctx context.Context, artifact string, tags map[string]string) error
}
```

```go
// internal/accessor/refusal.go — GUESS at exact file, RDR states the field
// addition only: "a NEW `Detail string` field on `accessor.Refusal`"

type Refusal struct {
    Category string // one of the closed ValidationCodes / execution refusal classes
    // ... existing fields (0004) ...
    Detail string // NEW (C4): last 4 KiB of stderr tail, or empty when no
                  // underlying error exists (timeout, read-back mismatch, gate-off)
}
```

```go
// internal/accessor/executor.go (existing file, per C4/S3 narrative)

// refusalOf is the base constructor — GUESS at exact signature, RDR states
// it "gains an error parameter" and does an errors.As on it.
func refusalOf(category string, err error) Refusal

// refusalWithKeys wraps refusalOf; Executor.Read routes every read refusal
// through it, so the error parameter threads through both — RDR is explicit
// that the read path (capability binding directly) was the one that used to
// silently drop the tail.
func refusalWithKeys(category string, keys []string, err error) Refusal

// readOutcome.classify — cited directly in S5 for the reserved `<clear>`
// literal's UNREADABLE classification on the read path.
type readOutcome struct { /* GUESS: fields for classification state */ }
func (o readOutcome) classify() /* GUESS: return type, likely a verdict/class enum */
```

```go
// internal/cli/flowbind/registry.go (existing file — S5b names it directly)

// Registry is a FREE FUNCTION (C6: "Registry is a free function, not a
// method"), the single production construction site, called from
// flow_exec.go:830. Signature changes to accept the --allow-commands gate.
func Registry(
    // ...existing params (GUESS at exact list — RDR does not enumerate)...
    allowCommands bool, // NEW (C6): when false, every constructed command
                        // binding refuses execution_failure before spawn,
                        // Detail naming the gate "allow_commands"
) *FlowbindRegistry // GUESS at return type name
```

Errors / refusal categories (closed set per C5, load-time — these are the
public contract a lint consumer sees):

```
command_and_path_conflict     // both or neither of path/command declared
command_empty                 // empty vector or empty argv element
command_unknown_placeholder   // unknown or non-whole-element {...} token
command_shell_interpreter     // argv0 + inline-code flag (sh -c, bash -c, python -c, env chains)
command_output_shape          // output="raw" with keys != 1; exit_absent on non-read; exit_verdicts on non-gate or naming a non-verdict
command_env_conflict          // env key or env_pass name matching reserved INTRASTATE_ prefix
```

All six are appended to `table.Categories()` at the tail, in the order
declared above, after existing (0004-era) members — `accessor.ValidationCodes`
count stays at 8 total categories (existing 2 + these... GUESS on the prior
count's composition, RDR only pins the post-change total via S1's REQ-84
reference).

Execution-time refusal classes (from RDR 0004, unchanged, reused): at minimum
`execution_failure` and `timeout`, plus `read_back_mismatch` and
`read_back_incomplete` (S5) for the write path's post-write verification.
GUESS: the exact enum/const names for these classes beyond the string labels
the RDR quotes; RDR never spells the Go identifier, only the observable
string.

## 2. TOML entry grammar (public surface authors write against)

```toml
[read.<id> | gate.<id> | write.<id>]
command       = ["<argv0>", "<arg>", ...]   # []string
output        = "json" | "raw"              # *string, read entries only; default "json"
exit_absent   = [<code>, ...]               # []int, read entries only
exit_verdicts = { "<code>" = "<verdict>" }  # map[string]string, gate entries only
env           = { "<KEY>" = "<value>" }     # map[string]string
env_pass      = ["<VAR>", ...]              # []string
```

Exactly one of `path` (existing) / `command` (new) per entry — C1's mutual
exclusion; an entry declaring neither is a REFUSING binding at runtime, never
a `Path: ""` file binding that silently reads an empty artifact (S7).

Placeholder vocabulary (C2): `{artifact}` only in v1 — substituted
whole-element (never mid-string; mid-string is a C5 load defect) with the
caller-bound artifact path for the entry's role.

## 3. Three most important internal helpers

**a. `load.go::loadAccessors` (existing file, extended)** — walks the model's
`read`/`gate`/`write` tables in that fixed order and returns on the first
error (C5's cross-table ordering is inherited from 0002, not established
here). Within one entry it runs C5's six new validators plus the existing
carrier/placeholder/output-shape checks fail-fast, reporting only the first
defect hit in the declared clause order. Responsibility: turn a raw TOML
entry into either a validated binding descriptor or one categorized
`ConfigError` — the single point where "does this model even lint" is
decided, and per S1 the six new categories must show up in
`table.Categories()` at the tail for a per-mutant test to mean anything
(a bare "validation returned non-empty" assertion would pass without the
registration).

**b. command-binding constructor (new, GUESSED file:
`internal/accessor/command_binding.go` or similar — RDR names the behaviour,
not the file)** — given a validated entry, builds the concrete argv (after
`{artifact}` substitution), the layered env (overlay > entry `env` >
`env_pass` > parent allowlist {PATH, HOME, TMPDIR, LANG, LC_*} per C4), and
wraps `exec.CommandContext` with `Setpgid` for a process group, `Cancel =
SIGKILL to -pgid` at context deadline, and `WaitDelay = 500ms` bounding the
stdin write and pipe drain. Responsibility: the single place that turns a
declared entry into an actual bounded child process — the load-bearing
deadline triple (S4) lives here with no injection seam (the ablations in S4
are cited as spike evidence, not shipped test knobs). Also does the
`runtime.GOOS` platform check (via injectable package-level `goos` var,
defaulting to `runtime.GOOS`) and refuses `execution_failure` before spawn on
non-Unix (C4 `platform`).

**c. `flow_exec.go::buildRequest` (existing file, four call sites: two verbs
plus `flow_state.go` twice)** — the sole place `--allow-commands` is read and
passed to `Registry` (C6's ONE gate site, `flow_exec.go:830`). Responsibility:
translate CLI invocation flags into the gated `flowbind.Registry` that every
executor shares — an unregistered verb's flag lookup returns `false`, which
REFUSES rather than silently bypassing (fail-closed by construction, per
C6's "drift costs execution, never a bypass").

## 4. Data model (persisted / crossing the process boundary)

Persisted (TOML, in the model file, author-facing):

```
[read.<id>]   command, output?, exit_absent?, env?, env_pass?, keys, timeout, role
[gate.<id>]   command, exit_verdicts?, env?, env_pass?, timeout, role
[write.<id>]  command, env?, env_pass?, keys, timeout, role, read_back?
```

Crossing the child-process boundary (wire format, C3):

```
stdin (write):      {"<key>": "<value>", ...}   // planned tags; <clear> literal for a planned clear
stdin (read/gate):  {}                            // always sent, even though empty of per-invocation data
                                                    // beyond {artifact} substitution already in argv

stdout (read, output="json", default):
    flat JSON object of strings; a declared key the object omits is
    UNREADABLE (never established-absent)

stdout (read, output="raw"):
    the single declared key's value = stdout minus one trailing "\n"
    (valid only when keys has exactly one entry)

stdout (gate):
    {"verdict": "allow" | "deny" | "indeterminate", "reason": "<text>"}

exit_absent   = [<code>, ...]                         // read: listed exit + empty stdout = every declared key absent
exit_verdicts = {"<code>": "allow"|"deny"|"indeterminate"}  // gate: listed exit + empty stdout = that verdict
```

Child environment (constructed per invocation, layered, later wins on
collision): `overlay {INTRASTATE_ROLE, INTRASTATE_CAPABILITY,
INTRASTATE_ACCESSOR, INTRASTATE_PROTOCOL=1}` > entry `env` > `env_pass`
named vars > parent allowlist `{PATH, HOME, TMPDIR, LANG, LC_*}` (LC_ is a
literal prefix match; `env_pass` is whole-variable-name only, no patterns).
Nothing else is inherited — an `env`/`env_pass` entry colliding with the
`INTRASTATE_` prefix is a load-time C5 defect, not a silent shadow.

Same wire shape (flat map-of-strings, presence-is-the-answer) as
`flowbind.go::store` already uses for intrastate's own artifact — one shape,
two transports (file vs. child stdout), decoded by separate decoders because
they decode different things.

## 5. Top-level pseudo-code (main operation: executing a command binding)

```
func (b *commandBinding) invoke(ctx context.Context, artifact string, stdinPayload map[string]string) (Result, Refusal) {
    // --- gate check (C6) ---
    if !b.allowCommands {
        return nil, refusal("execution_failure", detail="allow_commands")
    }

    // --- platform check (C4 platform) ---
    if goos != "linux" && goos != "darwin" /* non-Unix */ {
        return nil, refusal("execution_failure", detail="unsupported platform")
    }

    // --- build argv (C1/C2) ---
    argv := substituteWholeElement(b.command, "{artifact}", artifact)
    // load-time lint already proved: exactly one carrier, no unknown/mid-string
    // placeholders, no interpreter form, output shape valid, no env conflict

    // --- build layered env (C4) ---
    env := allowlistedParentEnv()         // PATH, HOME, TMPDIR, LANG, LC_*
    env = merge(env, b.envPass)           // named pass-through vars
    env = merge(env, b.env)               // entry's literal env map
    env = merge(env, overlayVars(b))      // INTRASTATE_* — always wins

    // --- spawn under bounded context (C4 deadline) ---
    cctx, cancel := context.WithTimeout(ctx, b.timeout)
    defer cancel()
    cmd := exec.CommandContext(cctx, argv[0], argv[1:]...)
    cmd.Env = env
    cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    cmd.Cancel = func() error { return sendSIGKILL(-processGroupOf(cmd)) }
    cmd.WaitDelay = 500 * time.Millisecond

    stdin, _ := json.Marshal(stdinPayload)
    cmd.Stdin = bytes.NewReader(stdin)
    var stdout, stderr boundedBuffer  // stdout capped at 1 MiB
    cmd.Stdout = &stdout
    cmd.Stderr = tailBuffer(4 * 1024) // last 4 KiB retained

    err := cmd.Run()

    // --- classify: envelope BEFORE exit code (C4 order) ---
    if stdout.overflowed() {
        return nil, refusalWithDetail("execution_failure", stderr.tail())
    }
    if env, ok := parseEnvelope(stdout.Bytes(), b.kind); ok {
        // well-formed envelope wins regardless of exit code
        return resultFrom(env), nil
    }
    if stdout.Len() > 0 {
        // non-empty but malformed: always execution_failure, whatever exit
        return nil, refusalWithDetail("execution_failure", stderr.tail())
    }

    // stdout is empty: consult exit code maps
    if cctx.Err() != nil {
        return nil, refusal("timeout")
    }
    switch b.kind {
    case Read:
        if contains(b.exitAbsent, cmd.ExitCode()) {
            return allKeysAbsent(), nil
        }
    case Gate:
        if verdict, ok := b.exitVerdicts[cmd.ExitCode()]; ok {
            return Result{Verdict: verdict}, nil
        }
    }
    return nil, refusalWithDetail("execution_failure", stderr.tail())
}

// Write path additionally, after invoke() succeeds:
func (b *commandWriteBinding) Write(ctx context.Context, artifact string, tags map[string]string) error {
    _, refusal := b.invoke(ctx, artifact, tags)
    if refusal != nil {
        return refusal
    }
    if b.readBack {
        got, rerr := b.reader.Read(ctx, artifact) // role's reader, first-match
                                                     // pending RDR 0016 (A4)
        if rerr != nil {
            return readBackIncomplete(rerr)
        }
        if !matches(got, tags) {
            return readBackMismatch(got, tags)
        }
    }
    return nil // write's own exit status decides NONE of this
}
```

## Notes on RDR silences / underdetermined spans

- Exact Go file names for the new command-binding constructor, the
  `Refusal.Detail` addition's file, and `readOutcome`/`classify`'s full
  signature are GUESSed — the RDR states behaviour and cites some file:line
  anchors (`executor.go::refusalOf`, `executor.go::refusalWithKeys`,
  `executor.go::readOutcome.classify`, `flow_exec.go:830`) but not every
  struct/field name.
- The exact `Registry(...)` parameter list beyond the new `allowCommands`
  bool is a GUESS; the RDR states only that it "gains" the gate as a
  signature change.
- Cross-entry-defect ordering within one table is explicitly UNSPECIFIED by
  the RDR itself (Go map ranging) — not a reconstruction gap, a documented
  non-guarantee (C5 precedence clause).
- Reader selection under multiple readers on one role is explicitly
  first-match/name-sorted pending RDR 0016, not yet fail-closed — carried
  through in the pseudo-code comment rather than guessed as fixed.
