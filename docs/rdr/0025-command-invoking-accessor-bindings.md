# Recommendation 0025: Command-Invoking Accessor Bindings

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
- **Type**: Feature
- **Profile**: foundational — Seam Lineage ≥2 floors it; Resolve recount:
  one seam, the declared-command carrier + authority bound (C1/C2/C5/C6)
  with its invocation envelope (C3/C4), inseparable — no split.
- **Priority**: High
- **Related Issues**: intrastate#v0hb (tracker);
  intrastate#p63c (reader cardinality per role, owned by
  RDR 0016)
- **Seam Lineage**: `internal/accessor/binding.go::WriteBinding`
  — 3rd point-fix; trail: dkcf (read-back seal) + 742a
  (0004 suite oracles). Count ≥2 → Profile floored at
  `foundational`; no accretion disposition.

## Problem Statement

A model author declares accessors so that `intrastate` can determine external
state by delegating a read to an established tool and apply a state change by
delegating a write to one — the product's own thesis. Today the second half is
missing: a model can DECIDE a state change (from a linted, exhaustiveness-proven
table) and can persist the decision into intrastate's own flat-JSON artifact,
but no binding invokes a command, so the change cannot be APPLIED to the
artifact the state actually lives in. The caller discovers this when it takes
the resolved decision and must execute the edit itself — an edit no model
declared and no lint covered, relocating exactly the prose-and-drift problem
the table was built to remove. The read side has the same hole: state only a
command can report (a VCS query, a file probe, a tool's own status verb) cannot
be a declared read.

The system-internal requirement: a binding family for
`internal/accessor/binding.go`'s `ReadBinding`/`GateBinding`/`WriteBinding`
that invokes a DECLARED command under the capability, timeout, read-back and
disposition rules RDR 0004 fixed, with the executed authority bounded and
reviewable statically from the model alone (`intrastate lint` validates the
declaration) — not a shell runner, and not a re-litigation of RDR 0004's
rejected raw shell-out (Alt 2) or global executable allowlist (Alt 3).

## Critical Assumptions

- **A1 [The declared timeout can genuinely bound a command invocation —
  including the child's process *tree*, a non-reading child's stdin pipe, and
  post-kill pipe drain — building on the executor's existing context deadline
  plus process-group termination and a drain bound this RDR adds in the
  binding]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `internal/accessor/executor.go::invokeRead` wraps every
    invocation in `context.WithTimeout` over the declared bound (`0004:C16`).
    Spike `evidence/spikes/a1-deadline/` (`run.sh` → `output.txt`, Go 1.26,
    darwin/arm64, timeout 1s, `WaitDelay` 500ms): with `Setpgid` +
    `Cancel = kill(-pgid, SIGKILL)` + `WaitDelay`, the grandchild-holds-pipe
    (S2, 1004ms), non-reading-stdin with a 1 MiB stdin (S4, 1002ms) and
    direct-sleeper (S5, 1002ms) cases all return within the bound with
    **zero** surviving processes in the group; the fast child (S6) is
    unaffected. Controls: naive `CommandContext` (S1) blocks `Wait` past the
    5s watchdog because the orphan holds stdout; `WaitDelay` alone (S3)
    returns at 1501ms but leaves `sleep 30` orphaned. Normative fixture
    **FX-deadline** (author-approved): S2/S3/S1 lines of `output.txt`.
    Mechanism fixed in C4. Two mechanism facts for the build: `Wait` returns
    `*ExitError` ("signal: killed"), never `context.DeadlineExceeded` — the
    timeout is classified from `ctx.Err()`; `$GOROOT/src/os/exec/exec.go`
    `WaitDelay` doc ("bounds the time spent waiting on … a child process that
    exits but leaves its I/O pipes unclosed") and `Cancel` doc anchor the
    semantics.
  - **If wrong**: the CLI hangs or a runaway child outlives the `timeout`
    refusal while still mutating the artifact.
- **A2 [`os/exec` passes the argv vector verbatim to the child: no shell, no
  word-splitting, no glob expansion, on the supported platforms]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `$GOROOT/src/os/exec/exec.go` package doc (Go 1.26.6):
    "the os/exec package intentionally does not invoke the system shell and
    does not expand any glob patterns or handle other expansions, pipelines,
    or redirections"; `Cmd.Args` doc: "Args holds command line arguments,
    including the command as Args[0]"; `exec.Command` resolves a
    separator-free name through `LookPath` in the **parent's** `PATH` before
    spawn. Caveat recorded: on Windows the vector is re-quoted into a single
    command line (`SysProcAttr.CmdLine` overrides) — v1's supported platforms
    are the Unix ones the process-group mechanism (A1, C4) needs.
    (`evidence/research/source-search.md` §A2.)
  - **If wrong**: the authority bound of C2 is fiction — a declared argument
    containing shell metacharacters executes something the model never
    declared.
- **A3 [Established tools bind directly for *reads and gates* —
  path-positional argv plus the raw single-value read mode or A8's declared
  exit mapping; established-tool *writes*, and any read whose output is
  multi-valued, are served by a thin declared wrapper conforming to the C3
  envelope, and that wrapper class is small and mechanical, not bespoke per
  integration]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: narrowed from "a useful class binds directly" by spike
    `evidence/spikes/a3-a7-a8-real-tools/` (`run.sh` → `output.txt`,
    `notes.md`, git 2.55): R1 `git config --file {artifact} --get state.phase`
    emits `draft\n` — not a JSON object — so it binds only via the raw mode C3
    now fixes (normative fixture **FX-raw-read**, author-approved); R2 the
    multi-valued `--list` needs a 5-line `k=v → JSON` wrapper
    (`wrapper-list.sh`); R3 **every** established-tool write needs a wrapper:
    the fixed argv cannot carry the planned value and `git config` ignores
    stdin (R3b leaves the artifact unchanged), so `wrapper-write.sh` (stdin
    JSON → `git config` per key, 8 lines) is the write form, and read-back
    through R1 then observes `review`. Hazard P-1 confirmed: a stdin-sinking
    argv0 (`tee {artifact}`, R4b) overwrites the artifact with the envelope and
    the next read-back fails `exit=128 bad config line 1` — named in Failure
    Modes. P-2: a 2 MiB envelope captures in 88ms (R5) — C4 caps stdout at
    1 MiB. Per tool class (notes.md table): VCS query → raw; file probe /
    status verb → exit-mapped (A8); VCS write → wrapped.
  - **If wrong**: every integration needs a bespoke wrapper script — the
    delegate-to-established-tools thesis fails and the choice tips toward the
    adapter-registry alternative.
- **A4 [Read-back for a command-backed write completes through the role's
  declared reader, which may itself be command-backed, preserving the
  applied-but-unverified sense when it cannot]**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: `internal/accessor/model.go::readerFor` — `func (reg
    Registry) readerFor(role string) (Definition, bool)`, predicate
    `Capability == CapRead && Accessor.Role == role` — selects by role only and
    returns a `Definition` whose binding is reached through the `ReadBinding`
    interface, so a command-backed reader is admissible without change;
    called from `internal/accessor/executor.go::Executor.Write` for read-back,
    whose failure is independent of the write (`0004:C13`, `0004:C14`).
    `0016:C4` pinned: "`internal/accessor/model.go::readerFor` resolves 'the
    role's read definition' uniquely — and the resolution site MUST fail
    closed". Landing order: `rdr status --tags 0016` reports `status=Draft`
    (`gate_stale=false`, `gate_written=false`) — today's `readerFor` is still
    first-match; 0025 lands after 0016 or inherits first-match until it does
    (Prerequisites).
  - **If wrong**: command writes are unverifiable by construction and every
    such write refuses `read_back_incomplete`.
- **A5 [Kubernetes `ExecAction` documents its probe command as an argv array
  executed without a shell — the demoted prior-art instance of the class
  claim]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `k8s.io/api@v0.36.2/core/v1/types.go` `type ExecAction
    struct`, `Command []string` doc comment: "The command is simply exec'd, it
    is not run inside a shell, so traditional shell instructions ('|', etc)
    won't work." (`evidence/research/source-search.md` §A5; promoted from the
    prior-art record's demoted list.)
  - **If wrong**: nothing structural — the class claim already rests on the
    quoted Docker exec-form and Terraform external-program citations.
- **A6 [RDR 0004's closed validation-code set closes 0004's own defect
  vocabulary, not the seam: a successor RDR may add codes for surfaces 0004
  did not define]**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: no 0004 contract clause closes the set; the only closure
    statement is `0004:§testing-strategy` S1 ("eight validation arms, each
    asserting its own named code"), realized in own code as the
    `internal/accessor/model.go::ValidationCode` comment ("closed eight-member
    validation-code set") and
    `internal/accessor/validation_0004_test.go::TestReq84_TheValidationCodeSetIsClosedAtEight`
    (exact-list assertion). So the closure is 0004's own test contract, not
    the seam's — and extending `ValidationCode` would break REQ-84. Decision
    (C5): the four defects are **`internal/table/category.go::Category`**
    constants, the family that already owns accessor-entry load defects
    (`CatMalformedAccessorDeclaration` — `load.go::accessorTable` emits it for
    "path is absent or empty" today).
  - **If wrong**: C5's defects must ship as a separate lint family instead of
    extending `accessor.ValidationCode` — carrier shape unchanged, surfacing
    relocated.
- **A7 [A child process's inherited environment does not silently widen the
  authority bound: an explicit env policy (inherit-as-is vs scrubbed) can be
  fixed at Resolve without reshaping the carrier]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: spike `evidence/spikes/a3-a7-a8-real-tools/` V1–V4: `git
    config --file` answers correctly with an **empty** child env (V1 — argv0
    resolves through the parent's `PATH` before spawn, A2) and with `PATH`
    only (V2); `--file` isolates from `~/.gitconfig`, `GIT_CONFIG_GLOBAL`,
    `GIT_DIR`, `GIT_CONFIG_SYSTEM`, `GIT_CONFIG_PARAMETERS` and
    `GIT_CONFIG_COUNT` injection (V3a–V4e all read `draft`), but the controls
    without `--file` (V3c, V4a', V4c'') return `FROM_HOME_GITCONFIG` /
    `FROM_GIT_CONFIG_GLOBAL` / `FROM_ENV_INJECTION` — the same binary and argv
    answer differently under inherited env, which is exactly the ambient
    widening the claim excludes. Policy fixed in C4: allowlist inheritance
    (`PATH`, `HOME`, `TMPDIR`, `LANG`, `LC_*`) plus named `env_pass` vars, a
    literal per-entry `env` table, and the `INTRASTATE_*` overlay; no ambient
    pass-through (the sudo `env_reset` / Bazel strict-action-env / Hugo
    `security.exec.osEnv` end-state — reversal ledger).
  - **If wrong**: ambient process state reaches the command invisibly — the
    same ambient-authority class `0004:C3` exists to exclude, one seam over.
- **A8 [Established gate/read tools that answer in exit codes (`test -f`,
  `git diff --quiet` conventions) can be admitted via a *declared* per-entry
  exit-code mapping — deny / established-absent as model-visible
  declarations — without laundering execution failure into a verdict]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: spike `evidence/spikes/a3-a7-a8-real-tools/` E1a–E1h and
    R6, normative fixture **FX-exit-codes** (author-approved): verdict codes
    and error codes are disjoint for the bound tools — `test -f` 0/1 (verdict)
    vs 2 (bad operator); `git diff --quiet` 0/1 vs 128 (bad path); `git config
    --get` missing key → 1 with empty stdout (established absent) vs a missing
    executable → `exec.ErrNotFound`, no exit code at all. A declared map that
    lists only verdict codes and sends every unlisted code, every spawn error
    and every malformed envelope to `execution_failure` therefore admits these
    tools without laundering failure into a verdict (C3 `exit_verdicts` /
    `exit_absent`).
  - **If wrong**: exit-speaking gates always need wrappers — narrows A3's
    direct-binding class but changes no contract.
- **A9 [No consumer besides `internal/table/load.go::accessorTable`, the
  registry's constructors, and the file binding's own runtime use of its
  copied locator requires `path` presence on an accessor entry — a
  locator-less command entry breaks no other partition]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: sweep of every non-test reader of the accessor `Path` field
    (`evidence/research/source-search.md` §A9): declaration
    `internal/table/source.go::sourceAcc`; validation + copy in
    `internal/table/load.go::accessorTable` (the "path is absent or empty"
    arm C1 relaxes); constructors `internal/cli/flowbind/registry.go::Registry`;
    and — the one site the claim now names — the file binding's runtime reads
    of its own copied `Path` (`internal/cli/flowbind/flowbind.go::Reader.Read`,
    `::Writer.Apply`, `::Gate.Gate` via `unreachable`/`verdictFor`), which
    only a path-constructed binding ever executes. Every other consumer of
    `Readers`/`Writers`/`Gates` (`load.go`, `normalize.go`, `flow_exec.go`,
    `flow_state.go`) reads `Role`/`Keys` only, and `dump.go` does not emit
    `Path`. Premortem P-7/P-11's pre-dispatch site does not exist.
  - **If wrong**: each such site gains the same carrier discrimination the
    registry gets — mechanical, but it must be enumerated before build.
- **A10 [The C6 execution gate needs a user-scope configuration surface this
  repo does not have, so its discovery and precedence rules are this RDR's to
  fix — and `--allow-commands` alone is a complete gate without them]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to verify — the absence is established (no
    `internal/cli/config/`; zero `intrastate.toml` or `allow_commands` hits in
    non-test `*.go`; no user-scope TOML loader). What remains is the positive
    call: whether v1 ships the flag alone or also the file, and if the file,
    its search order, its precedence against the flag, and its
    absent/malformed behaviour — a malformed config must not fail *open*.
  - **If wrong**: C6 is unimplementable as written and S7's config arms have
    nothing to exercise; the gate degrades to the flag, which changes the
    "one-time opt-in" property the reversal ledger argues for.

## Proposed Solution

### Approach

The undecided contract at this seam: **what carries a declared accessor's
real external behaviour**. Today `flowbind` *simulates* command behaviour
through magic `path` suffixes (`internal/cli/flowbind/flowbind.go::verdictFor`,
`::unreachable`) because no declared-command carrier or authority bound
exists — the gap both accreted point-fixes (Seam Lineage) patched around.

Answer: an accessor entry may declare `command`, a **fixed argv vector**
carried on the same TOML entry as its role, keys, and timeout. The vector is
validated at load, executed with **no shell**, and admits exactly one form of
variation: a **closed, tool-defined placeholder vocabulary** (v1: `{artifact}`)
substituted **whole-element only** with the caller-bound artifact path for the
entry's declared role. All other per-invocation data crosses on **stdin as a
JSON object of strings**; read results return on **stdout as a flat JSON
object of strings** — the same shape intrastate's own artifact and
Terraform's external-program protocol already use — and gate results as a
verdict envelope (C3). Command-backed bindings
implement the existing `ReadBinding`/`GateBinding`/`WriteBinding` seam and run
under RDR 0004's executor unchanged: declared timeout, refusal classes,
post-write read-back through the role's reader.

The executed authority is therefore reviewable statically: the reviewer reads
the argv in the model; lint proves the vector is well-formed, the placeholders
are known and whole-element, the entry carries exactly one of
`path`/`command`, and `argv0` is not an inline-code interpreter form (C5).
Nothing a caller supplies at invocation time can change *which words* the
command runs — only which artifact path fills the one declared hole and what
typed data arrives on stdin.

Integration honesty (premortem P-1/P-8/P-9): a tool binds *directly* when it
is path-positional and its semantics fit the envelope, the raw single-value
read mode, or the declared exit maps (A3, A8) — in practice reads and gates;
established-tool **writes** always take the other form, because a fixed argv
cannot carry the planned value and the tools do not read stdin (A3 spike).
That form is a **thin declared wrapper that conforms to the same C3
envelope** — its invocation argv is still in the model and linted, and the
envelope bounds what it must do, but its body is script the model does not
carry (accepted, recorded as a Consequence).

### Technical Design

Components:

- **Carrier** — `internal/table`: `sourceAcc`/`Accessor` gain `command`
  (`[]string`); `internal/table/load.go::accessorTable` today rejects a
  missing `path` ⇒ its rule relaxes to *exactly one of* `path`/`command`
  (C1), with the C5 defects joining the existing load-time rejections.
- **Binding family** — a sibling package to `flowbind` implementing the three
  binding interfaces over `os/exec`. The executor already owns timeout
  enforcement: `internal/accessor/executor.go::invokeRead` wraps invocation in
  `context.WithTimeout` (writer: the executor, from the declared `timeout`
  metadata) ⇒ the binding builds on `exec.CommandContext` and carries no
  timer of its own (A1).
- **Selection** — `internal/cli/flowbind/registry.go::Registry` constructs
  `Reader`/`Writer`/`Gate` from `acc.Path` per capability table ⇒ the carrier
  discriminator slots exactly there: a `command`-carrying entry constructs the
  command binding, a `path`-carrying entry constructs today's file binding.
  Sibling-path check: searched for an existing path-vs-command discriminator —
  none exists (`os/exec` appears nowhere outside tests; `flowbind` is the sole
  non-test binding implementation).
- **Data flow** — per capability, under C3/C4: a read command receives an
  empty object on stdin (its key set is the entry's declared `keys`, already
  in the model) and answers with the flat string map (an omitted key is
  not-carried, mirroring `0004:C8`) or, in raw mode, the one value; a gate
  command answers a verdict envelope or a declared exit code (three-valued,
  refusal never laundered into a verdict — the rule
  `internal/cli/flowbind/flowbind.go::Gate` already keeps); a write command
  receives the planned tags on stdin, `<clear>` literal, and is verified only
  by read-back. Field grammars are C3's, fixed by the A3/A8 spikes.

#### Normative Contracts

**C1** — carrier and mutual exclusion.

```normative
[read.<id> | gate.<id> | write.<id>]
command = ["<argv0>", "<arg>", ...]   # TOML array of strings
```

An accessor entry carries **exactly one** of `path` or `command`; both or
neither is a load-time defect. `command` must be non-empty and contain no
empty element. `role`, `keys`, `timeout`, and `read_back` rules are unchanged
from RDR 0002/0004.

**C2**

```normative
{artifact}   # the closed placeholder vocabulary, v1 complete: replaced whole-element by the caller-bound artifact path for the entry's declared role
```

Authority bound. The executed argv is exactly the declared vector after
whole-element placeholder substitution. The placeholder vocabulary is closed
and defined by intrastate, not the model author; v1 is exactly the one token
above. A placeholder is recognized only as a whole argv element. An element that
contains a `{...}` token without being exactly a known placeholder is a
load-time defect (C5) — never silently-literal text. Execution invokes no
shell and performs no other rewriting of the vector.

Substitution contract (premortem P-3): the substituted value must be an
absolute path — the executor refuses the invocation (`execution_failure`)
when the caller-bound artifact path is relative or begins with `-`, so a
path can never be parsed as a flag by the invoked tool.

argv0 resolution is fixed at load, never cwd-relative: a bare name resolves
through the parent's `PATH` at spawn (A2); a name containing a path
separator resolves against the **model file's directory** (the git-hook /
pre-commit convention); an empty argv0 is a C5 defect, and the binding never
restores implicit current-directory lookup — Go itself shipped cwd-relative
`PATH` resolution for a decade and reversed it (`exec.ErrDot`, Go 1.19;
reversal ledger, `evidence/research/prior-art-resolve-devref-go.md`).

**C3**

```normative
stdin (write): {"<key>": "<value>", ...}   # the planned tags; a planned `<clear>` crosses as the literal reserved value
stdin (read/gate): {}                       # nothing per-invocation beyond {artifact}; the object is always sent
stdout (read, output = "json", default): flat JSON object of strings; an omitted key is not-carried
stdout (read, output = "raw"):            the single declared key's value = stdout minus one trailing "\n"; valid only when `keys` has exactly one entry
stdout (gate): {"verdict": "allow" | "deny" | "indeterminate", "reason": "<text>"}
exit_absent   = [<code>, ...]                # read entries: a listed exit with empty stdout establishes every declared key absent
exit_verdicts = { "<code>" = "allow" | "deny" | "indeterminate", ... }   # gate entries: a listed exit with empty stdout is that verdict
```

Invocation envelope. Per-invocation data beyond the artifact path crosses
only on stdin; read and gate results return only on stdout, in the shapes
above. A write command's success is never taken from its exit status alone —
verification is read-back (`0004:C12`). Grammar fixed at Resolve: the gate
`verdict` strings are `internal/accessor/model.go::Verdict`'s
(`0004:C9`); `<clear>` is carried unchanged as the planned value — the tool
or wrapper performs the removal, and read-back verifies absence
(`0004:C11`). `output = "raw"` and the two exit maps exist because established
tools speak single values and exit codes, not envelopes (A3, A8 spikes); the
maps list **verdict** codes only — any unlisted exit, spawn failure, or
malformed stdout is `execution_failure` (C4), so execution failure is never
laundered into a verdict. A non-empty stdout is always parsed first (C4
ordering); the exit maps apply only to an empty stdout. The protocol's
version rides out-of-band as `INTRASTATE_PROTOCOL` in the child env (C4
overlay), keeping the stdout map flat — the envelope-versioning lesson of
kubebuilder's external-plugin `apiVersion` and client-go's exec-credential
`v1alpha1`→`v1` migration (reversal ledger).

**C4**

```normative
order:    parse the stdout envelope, THEN classify the exit code
deadline: Setpgid; at ctx deadline Cancel = SIGKILL to -pgid; WaitDelay = 500ms bounds the stdin write and the pipe drain; timeout is classified from ctx.Err()
bounds:   stdout capped at 1 MiB (overflow = execution_failure); a NEW `Detail string` field on `accessor.Refusal` carries the last 4 KiB of stderr
env:      child env = allowlisted parent vars {PATH, HOME, TMPDIR, LANG, LC_*} + named `env_pass = ["VAR", ...]` vars + the entry's literal `env = { KEY = "value" }` + the overlay {INTRASTATE_ROLE, INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR, INTRASTATE_PROTOCOL=1}; nothing else is inherited
write:    read_back required; verified only through the role's reader, never by exit status
```

Execution-safety inheritance. Command bindings run under RDR 0004's executor
with its refusal *classes* unchanged and one additive type change — a `Detail
string` field on `accessor.Refusal`, which today carries no free-text
diagnosis (`internal/accessor/model.go` Refusal: Class, Accessor, Capability,
Role, Timeout, Keys, Expected, Observed, Reason). The stderr tail has no
existing carrier: the `Detail` fields that do exist belong to
`accessor.Finding` (validation) and `table.Failure` (load), neither of which
is on the refusal path. The addition is additive-only — 0004 pins the refusal
*class* set, not the struct's field set (contrast `ValidationCodes`' eight-member
closure, REQ-84) — and `Reason` is not reused because `0004:C7` reserves it for a
gate deny. With that field, the declared timeout
bounds the invocation, and at the deadline the binding terminates the child's
**process group**, not only the direct child, with a bounded stdin write and a
bounded post-kill pipe-drain wait so a non-reading or slow-draining child can
never hang the CLI (premortem P-4). The mechanism is the A1 spike's
necessary-and-sufficient triple: `Setpgid`, a `Cancel` that signals the
group, and `WaitDelay` — the spike shows dropping the group signal orphans
the grandchild and dropping `WaitDelay` hangs `Wait`; and because `Wait`
reports the kill as an `ExitError`, the deadline is classified from
`ctx.Err()`, never from the wait error. The deadline mints
`ClassTimeout`; spawn failure, a malformed or oversized stdout envelope, and —
by default — a non-zero exit are `ClassExecutionFailure`, whose `Detail`
carries the bounded stderr tail. Ordering is normative (premortem P-12):
the stdout envelope is parsed **before** exit-code classification, so a
well-formed deny envelope with a non-zero exit is a deny, not a failure. A
non-zero exit is a gate verdict or an established absence **only** when the
entry's C3 exit map lists that code and stdout is empty. The env policy is an
allowlist because the A7 spike showed the same binary and argv answering from
`~/.gitconfig`, `GIT_CONFIG_GLOBAL` or `GIT_CONFIG_COUNT` injection when those
were inherited; `argv0` still resolves through the parent's `PATH` before
spawn (A2), which the allowlist passes on unchanged. `env_pass` is the named,
no-glob escape hatch (an `SSH_AUTH_SOCK`-class need must be declared, per
entry, in the model where review sees it); the `INTRASTATE_*` overlay gives
wrappers their context without new placeholders — the `CONSUL_INDEX` /
`ETCD_WATCH_*` / `HELM_PLUGIN_*` convention. A write entry carries
`read_back = true` (mandatory for every write since RDR 0004), its success is
never taken from exit status, and it is
verified through the role's declared reader — the re-read runs under its own
bounded timeout (`0004:C15`), never the write's residue — with the
applied-but-unverified sense (`0004:C14`) preserved when read-back cannot
complete.

**C5**

```normative
command_and_path_conflict     # both or neither of path/command declared
command_empty                 # empty vector or empty argv element
command_unknown_placeholder   # unknown or non-whole-element {…} token
command_shell_interpreter     # argv0 + inline-code flag (sh -c, bash -c, python -c, env chains); no opt-in in v1
command_output_shape          # output = "raw" with keys ≠ 1; exit_absent on a non-read; exit_verdicts on a non-gate or naming a non-verdict

registration: all five are appended to `table.Categories()`, whose hand-maintained list IS the closed set
interpreter set: OPEN (deny-listed, not closed) — an unlisted interpreter is admitted, so the list grows by amendment
```

Static validation. New load-time defect classes, surfaced by the same
validation that rejects malformed accessor entries today and reported by
`intrastate lint`. Home fixed at Resolve (A6): these are
`internal/table/category.go::Category` constants with the spellings above —
the family `load.go::accessorTable` already uses for accessor-entry defects —
not `accessor.ValidationCode`, whose eight-member closure is RDR 0004's own
test contract. A `Category` constant is not in the closed set until it is
appended to `table.Categories()`, whose returned literal list *is* the set
(RDR 0024's "ONE decision") — an unregistered constant compiles and refuses
correctly at the call site while remaining invisible to every consumer that
enumerates categories, so S1 must assert membership, not merely the refusal.
`command_shell_interpreter` is what keeps "not a shell
runner" enforced rather than aspirational (premortem P-3/P-10): a fixed
`["bash", "-c", …]` vector is still a shell runner, so the known
interpreter-with-inline-code forms are load-time defects. That list is a
**deny-list, explicitly not closed** — unlike C3's exit maps or 0004's
eight-member code closure, and the distinction is normative because the two
fail in opposite directions: an unlisted interpreter (`perl -e`, a
`busybox sh` alias, a novel runtime) is *admitted*, so the check raises the
cost of an inline-shell carrier without claiming to make one impossible. The
guarantee C2/C5 actually enforce is that the argv is fixed and reviewable;
the interpreter deny-list is defense in depth over that, not the barrier
itself. Adding a form is an amendment to this clause, not a lint-rule tweak.
There is no opt-in
in v1: inline code belongs in a wrapper *file* named as `argv0`, which is the
form A3's spike wrappers take — and the defect's report says so ("inline
shell is not a declared command; put it in a script and declare the script
as argv0"), the remediation-in-the-refusal pattern ralph-tui's metacharacter
ban uses. Every peer that shipped a shell-string carrier later deprecated it
(Helm subprocess v1, Consul `script`/`handler`→`args`) or bolted heuristic
lint over it (awf-cli's security-validator: "an opaque shell string cannot
be validated structurally"); no argv carrier reversed (reversal ledger,
`evidence/research/prior-art-resolve-{cli-a,cli-b,sm}.md`).

**C6**

```normative
allow_commands = true   # user-scope config (NEW surface, A10) or --allow-commands; never the model file; absent or false ⇒ every command invocation refuses execution_failure naming the gate; lint (C1/C5) validates regardless
```

Execution gate. A model-declared command is code that runs when the model is
used, and the model travels with the repo — so execution requires a one-time
opt-in **outside the model**, in the invoking user's configuration; a
`--allow-commands` invocation flag forces on for one run. **The user-scope
config surface does not exist yet and this RDR builds it** (A10): the repo
has no `internal/cli/config/`, no `intrastate.toml` reader, and no user-scope
configuration of any kind today — so the discovery rules (search order,
precedence against the flag, behaviour when the file is absent or malformed)
are this RDR's to fix, not an existing surface to cite. The flag alone is a
complete gate for a v1 that ships before the file: with no config surface,
`--allow-commands` is the only opt-in and absence refuses. With the gate off,
a command invocation refuses
before spawn (`execution_failure`, `Detail` naming
`allow_commands`), while `intrastate lint` still validates the entries — the
model is reviewable before it is trusted. This is the shape every shipped
peer converged on after gating too late: Consul script checks flipped
off-by-default and re-gated twice more, Hugo's deny-by-default
`security.exec.allow`, Go's `GOVCS` allowlist, beads refusing repo-persisted
commands pending a trust gate (reversal ledger).

#### Load-Bearing Decisions

- **Identity** — unchanged: the accessor identity stays RDR 0004's
  `(flow, name, capability)` triple; whether an entry is path- or
  command-backed does not enter identity.
- **Wire / byte format** — stdin/stdout JSON-object-of-strings envelopes
  plus the raw single-value read mode and the two exit maps (C3), field
  grammar fixed by the A3/A8 spikes against `git config`, `git diff --quiet`
  and `test -f`. Rejected: carrying the requested key set on stdin (a list is
  not a string value, and the keys are already declared in the model).
  The read envelope is deliberately the same flat map-of-strings, with the
  same presence-is-the-answer rule, that `flowbind.go::store` already uses for
  intrastate's own artifact — one wire shape, two transports (a file, a
  child's stdout), not a second format. They stay separate decoders because
  they decode different things; the shared rule is what must not drift.
- **Naming** — the field is `command`; the family is "command-backed
  accessors". Rejected: `exec` (implies a shell surface), `run` (verb
  collision with CLI verbs), `argv` (implementation jargon in an
  author-facing carrier).
- **Selection / predicate** — the registry selects the binding constructor by
  carrier field: `path` → file binding (today's `flowbind`), `command` →
  command binding. C1's exactly-one rule makes the selection total; no
  precedence order exists to get wrong. The file binding keeps its magic-suffix
  vocabulary (`verdictFor`, `unreachable`) **unchanged** — so this RDR bounds
  the accretion rather than ending it: it removes the *pressure* to grow that
  vocabulary by giving real external behaviour a declared carrier, but the
  simulated suffixes remain for path-backed entries and their retirement is a
  successor's decision, not this one's.

#### Illustrative Code

Illustrative — intent only; tests must not assert it literally.

```toml
[read.branch_state]
role        = "repo"
command     = ["git", "config", "--file", "{artifact}", "--get", "flow.stage"]
output      = "raw"          # stdout is the one value, not an envelope
exit_absent = [1]            # `--get` of a missing key: exit 1, empty stdout
keys        = ["flow.stage"]
timeout     = "5s"

[gate.clean_tree]
role          = "repo"
command       = ["git", "-C", "{artifact}", "diff", "--quiet"]
exit_verdicts = { "0" = "allow", "1" = "deny" }   # 128 (bad path) stays execution_failure
timeout       = "5s"

[write.stage]
role      = "repo"
command   = ["tools/flowstate-write", "{artifact}"]   # separator ⇒ resolved against the model file's dir; wrapper: stdin JSON → git config per key
env       = { GIT_CONFIG_NOSYSTEM = "1" }
env_pass  = ["SSH_AUTH_SOCK"]                          # named pass-through; nothing else inherited beyond the C4 allowlist
keys      = ["flow.stage"]
timeout   = "10s"
read_back = true
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Bounded child execution under a context deadline | Go stdlib `os/exec` (`CommandContext`) | Available | A1/A2 verify the bound is real |
| Timeout enforcement and refusal classes at the seam | RDR 0004 executor (`internal/accessor`) | Available | C4 inherits, adds no class |
| Reader-per-role read-back resolution | RDR 0016 (`0016:C4`; Draft at Resolve, intrastate#p63c) | Deferred | A4 — `readerFor` is binding-agnostic today (first-match); fail-closed uniqueness arrives with 0016 |
| Process-group termination + drain bound | Go stdlib `os/exec` (`Setpgid`, `Cancel`, `WaitDelay`) | Available | C4 mechanism, A1 spike |
| User-scope configuration discovery (`allow_commands`) | **none — no config subsystem exists** (no `internal/cli/config/`, no `intrastate.toml` reader) | Build | C6 — this RDR fixes discovery + precedence (A10); `--allow-commands` alone is a complete v1 gate |
| Free-text diagnosis on a refusal | `accessor.Refusal` — no `Detail` field today | Extend | C4 — additive `Detail string`; class set unchanged |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Declared-accessor carrier | `internal/table` (`sourceAcc`, `Accessor`, `load.go::accessorTable`) | `path` is required today | Extend | C1 exactly-one rule |
| Execution safety + validation | `internal/accessor` (executor, `Validate`, `ValidationCode`) | validation-code set closed at eight by 0004's REQ-84 test | Reuse unchanged | C5 homes in `table.Category` (A6) |
| Binding construction | `internal/cli/flowbind/registry.go::Registry` | constructs path-backed bindings only | Extend | Selection LBD |
| Simulated command behaviour | `flowbind::verdictFor` / `::unreachable` path suffixes | magic-string simulation — the accreted seam | Reuse unchanged | command family removes the pressure to grow it |

### Decision Rationale

Scored QOC matrix — options answer the contract question "what carries a
declared command and what exactly may vary at invocation time":
O1 fixed argv + stdin only; **O2 closed-placeholder argv (chosen)**; O3
templated/shell strings; O4 built-in adapter registry. 5 = best.

| Criterion | O1 | O2 | O3 | O4 |
| --- | --- | --- | --- | --- |
| Static reviewability (lint from model alone) | 5 — argv literal | 5 — argv + closed whole-element holes | 1 — executed text is runtime-computed (PA-1/PA-2) | 5 — argv owned by intrastate |
| Binds established tools without wrappers | 2 — artifact path cannot reach argv; `0004:C3` bars ambient discovery as the workaround | 4 — `{artifact}` covers the positional case | 5 — anything spliceable | 2 — only shipped adapters |
| Prior-art alignment | 4 — Terraform external (PA-5) | 5 — Docker exec form (PA-3) + git difftool closed vocabulary (PA-4) + PA-5 envelope | 2 — awf-cli needed an injection linter after the fact (PA-1) | 3 — plugin-registry precedent, none in the peer family |
| Injection / blast radius | 5 — nothing varies | 4 — one path-valued hole, whole-element | 1 — caller data spliced into text | 5 — no author-declared argv |
| Integration cost per new tool | 2 — a protocol wrapper per tool | 4 — direct for path-positional tools, thin wrapper otherwise | 5 — none | 1 — an intrastate release per tool |
| Reversibility | 4 — can add placeholders later | 4 — vocabulary can grow or shrink by validation | 2 — authors embed arbitrary templates immediately | 2 — adapter API becomes public surface |
| **Total** | **22** | **26** | **16** | **18** |

O2 wins on the deciding rows: it is the only option scoring ≥4 on both static
reviewability and binding established tools — the two clauses the Problem
Statement conjoins. O1 concedes the product thesis (every tool needs a
wrapper) to gain no review advantage over O2; O3 concedes review to gain
expressiveness O2 mostly retains; O4 relocates authorship, not authority, and
prices every integration at a release.

Premortem: hardened — the critic's 14-row P-ledger
(evidence/propose-premortem/critic.md) forced into the draft: the
interpreter-argv0 lint defect and dash-guarded absolute-path substitution
(P-3/P-10 → C2/C5), process-group + stdin/drain bounding (P-4 → C4/A1),
envelope-before-exit ordering and the declared exit-mapping question
(P-5/P-12 → C4/A8), the wrapper-honesty restatement and inverted-rejection
repairs (P-1/P-8/P-9 → Approach, Alt 1), locator-presence sweep (P-7/P-11 →
A9), and the read-back staleness/reader-less rows (P-6 → Risks, A4). The
recommendation survives with those folded.
Ground-sweep: clean (18 anchors) — 11 code anchors and 7 peer elements
(0004:C3/C8/C12–C16) all CONFIRMED by a factored checker; one cosmetic note
(role-routing of read-back reads from C12/code, independent-failure from
C13), no load-bearing miss.
Joint-check: fired → 0016 (home: cli/0016 §Normative Contracts C4);
fired → 0020 (home: cli/0025 §Normative Contracts C1 / cli/0020 §Normative
Contracts C1 — mutual tolerance: disjoint rule families at `accessorTable`,
entry shape here, tag admission there; composes) — all three arms ran on the
written proposal. Arm 1 (modify-anchors, `--repo`-resolved): pair with 0020
on `internal/table/load.go::accessorTable` (uncited — 0025's C1 relaxes the
path-required rule there; 0020's tag-admission touch is a disjoint rule
family in the same loader and composes, each record homing its own C1) and
pair with 0016 on `internal/accessor/model.go::readerFor`
(cited — A4 declares the coupling: 0016:C4 is the normative home for
reader selection; 0025's command-write read-back rides it as consumer,
pinning the element id via A4). Arm 2 (contract
literals): no pair involving 0025. Arm 3 (absence, manual): C1 converts the
load-time path-required refusal into acceptance for command entries; grepped
the locked peers (0001–0011) for reliance on that refusal — the line-level
hits are unrelated senses of "path" (resolver disposition, SCXML masking,
output path), and 0004's Alt-2/Alt-3 rejections bound shell and allowlist
shapes, both of which C2/C5 preserve — no locked peer relies on the refusal.

## Alternatives Considered

### Alternative 1: Fixed argv + stdin-only protocol (no placeholders)

**Description**: The Terraform-pure form (PA-5): the declared vector is fully
literal; every per-invocation datum, including the artifact path, crosses on
stdin as JSON.

**Pros**:

- Strictest possible authority bound — nothing varies, ever.
- Exact match to a quoted prior-art protocol.

**Cons**:

- The caller-bound artifact path cannot reach an established tool's argv, so
  no established tool is directly bindable — each needs a protocol-aware
  wrapper, and `0004:C3` (no ambient discovery) rightly bars "the tool
  defaults to cwd" as the workaround.
- Pushes authors toward wrapper scripts, which are exactly the undeclared
  prose surface the product removes.

**Reason for rejection**: it forces the wrapper for *every* tool, including
the path-positional class O2 binds directly (and the exit-speaking class A8
may admit); some tools need wrappers under O2 too (P-8), but O1 pays that cost
universally while O2's one
path-valued, whole-element, dash-guarded hole keeps the review property O1
was buying (P-3 mitigations in C2/C5).

### Alternative 2: Built-in adapter registry

**Description**: The model names a shipped adapter (`adapter = "git-config"`)
plus typed params; intrastate owns every argv shape.

**Pros**:

- Tightest authority: no author-declared argv at all.
- Adapters can be individually hardened and tested.

**Cons**:

- Every new tool integration is an intrastate release — the author cannot
  declare delegation to a tool intrastate has not met.
- Adapter params become a second, ad-hoc command language with its own
  validation burden.

**Reason for rejection**: relocates authorship, not authority; contradicts the
product thesis that the *model author* declares the delegation.

### Briefly Rejected

- **Shell-string / open-templated commands**: re-litigates RDR 0004 Alt 2;
  the peer evidence is decisive — awf-cli carries `command:` shell strings
  with `{{.inputs...}}` interpolation and needed a security-validator plugin
  to warn about "unquoted variable references in commands" after the fact
  (PA-1), and ms-conductor's full Jinja2 templating makes the executed text a
  runtime function (PA-2).
- **Global executable allowlist**: RDR 0004 Alt 3, rejected there —
  "constrains mechanism but not semantic power"; nothing here changes that.
- **Host-language callbacks**: RDR 0004 Alt 5, rejected there — authority
  hidden in code, invisible to lint.

## Context

### Background

Filed as kata intrastate#v0hb after RDR 0004 (Accessor Execution Safety Model)
shipped the capability-bounded execution seam — read, gate, write, each
declared, artifact-role-scoped, timeout-bounded, and (for write) verified by
read-back. No RDR 0001–0024 adjudicates a command-invoking binding (RDR 0014
uses `exec.Command` only in illustrative snippets); RDR 0004's deviations
D1–D18 fence nothing out here, and D17's "write command" language anticipates
process semantics.

The design fork: whether a declared external command is carried as a fixed
argv vector on the accessor's capability-table entry — no shell
interpretation, no interpolation of caller/artifact data — versus a
parameterized/templated command shape that admits bounded substitution. The
lint-validatable-from-the-model-alone authority bound is the load-bearing
fork; carrier placement (generalizing the per-accessor `path` locator vs a
sibling field vs a new accessor shape) is its dependent clause, decided by the
same answer.

### Technical Environment

Go CLI (`intrastate`). The seam: `internal/accessor/binding.go`
(`ReadBinding`/`GateBinding`/`WriteBinding`, `CapRead`/`CapWrite`/`CapGate`);
sole non-test implementation `internal/cli/flowbind` (flat-JSON artifact,
self-described as "deliberately the thinnest thing that satisfies the seam").
Governing peers: RDR 0002 owns the TOML model carrier and the per-accessor
`path` locator; RDR 0004 owns the capability/timeout/read-back/disposition
classes (0004:C12–C14; `ClassTimeout`, `ClassExecutionFailure`) and the
applied-but-unverified sense; RDR 0016 owns read-back comparison-set semantics
and reader-per-role resolution (kata intrastate#p63c). Deviations ledger:
`docs/rdr/0004-accessor-execution-safety-model/artifacts/deviations.md`.

## Research Findings

### Investigation

Prior art was read before enumeration (full record with quotes, queries, and
rejected branches: `evidence/research/prior-art.md`). The peer
state-machine/workflow family splits into shell-string templating (awf-cli,
ms-conductor — PA-1/PA-2) and host-language callbacks (xstate, temporal,
restate — 0004 Alt 5's shape); none carries a statically lint-validatable
command bound, so the bounded forms come from outside the family: Docker's
exec form ("doesn't automatically invoke a command shell … variable
substitution doesn't happen" — PA-3) ⇒ argv-array carriers get no-shell by
construction; git difftool's `$LOCAL`/`$REMOTE` ("evaluated in shell with the
following variables available…" — PA-4) ⇒ when per-invocation data must reach
a declared command, established practice is a closed tool-defined vocabulary
carrying paths, never content interpolation; Terraform's external program
("a list of strings…", "does not execute the program through a shell", JSON
object of strings on stdin/stdout — PA-5) ⇒ the envelope shape for everything
that is not path-positional. Code paths grounding the seam:
`internal/accessor/executor.go::invokeRead` (executor-owned timeout),
`internal/cli/flowbind/registry.go::Registry` (the constructor-selection
slot), `internal/table/load.go::accessorTable` (the `path`-required rule C1
relaxes).

### Key Discoveries

- **Documented** — awf-cli's model carries shell strings with Go-template
  interpolation and ships a lint plugin warning about injection and missing
  timeouts after the fact (PA-1) — the templated-carrier end state.
- **Documented** — Docker exec form and Terraform external-program protocol
  both fix no-shell argv execution; Terraform additionally fixes the
  JSON-object-of-strings stdin/stdout envelope (PA-3, PA-5).
- **Documented** — git difftool demonstrates the closed placeholder
  vocabulary carrying file paths (PA-4).
- **Verified** — the executor already owns timeout enforcement via
  `context.WithTimeout` (`internal/accessor/executor.go::invokeRead`), and no
  path-vs-command discriminator exists anywhere in the tree (`os/exec` absent
  outside tests; reuse audit re-run at Resolve).
- **Verified (spike)** — a context deadline alone does not bound a command:
  `Setpgid` + group `SIGKILL` + `WaitDelay` is the necessary-and-sufficient
  triple (A1; `evidence/spikes/a1-deadline/`).
- **Verified (spike)** — established tools speak single values and exit
  codes, never the envelope: reads and gates bind directly through the raw
  mode and exit maps, every established-tool write needs a wrapper because a
  fixed argv cannot carry the value and the tool ignores stdin (A3, A8;
  `evidence/spikes/a3-a7-a8-real-tools/`).
- **Verified (spike)** — inherited env silently steers the same binary and
  argv (`~/.gitconfig`, `GIT_CONFIG_GLOBAL`, `GIT_CONFIG_COUNT`), so the child
  env is an allowlist, not inheritance (A7).
- **Documented** — Kubernetes `ExecAction.Command` is "simply exec'd, it is
  not run inside a shell" (A5, promoted from the prior-art demoted list).
- **Documented** — the reversal ledger
  (`evidence/research/prior-art-resolve-{cli-a,cli-b,sm,devref-go}.md`):
  Helm removed shell-string plugin commands for argv (subprocess v1); Consul
  deprecated `script`/`handler` for `args` and re-gated script checks
  off-by-default; Hugo added deny-by-default exec + env allowlists post hoc;
  Go reversed cwd-relative `PATH` lookup (`exec.ErrDot`, 1.19) and narrowed
  VCS binaries (`GOVCS`); of 32 peer state-machine tools only four let a
  model name a command, and both shell-string carriers accrued corrective
  ADRs/lints while neither argv carrier did. String→argv, inherit→allowlist,
  open→gated — no project moved the other way; C1/C2, C4's allowlist and
  C6's gate each sit on the settled side.

## Trade-offs

### Consequences

- Positive: a decided transition can finally be APPLIED to the artifact the
  state lives in, by a declared, linted, timeout-bounded command — closing
  the product-thesis gap with authority visible in the model.
- Positive: `flowbind`'s magic-path simulation stops being the only way to
  express command-like behaviour, removing the pressure that accreted at this
  seam. The suffixes themselves stay (Selection LBD) — bounded, not retired.
- Negative: a new process boundary enters the trust surface — child env, PATH
  resolution of `argv0`, and process-listing visibility of argv become
  reviewable concerns (A7, Failure Modes).
- Negative: every established-tool *write*, and any multi-valued read, needs
  a thin declared wrapper — the A3 spike's wrappers are 5–8 lines of shell,
  but their bodies are script the model does not carry.
- Negative: the child env is an allowlist (C4), so a tool that needs another
  variable declares it literally in `env` or names it in `env_pass` — an
  explicit, reviewable cost.
- Negative: nothing executes until the invoking user sets
  `allow_commands = true` once in their own config (C6) — one boolean of
  adoption friction, accepted as the price every gated-too-late peer ended
  up paying anyway.

### Risks and Mitigations

- **Risk**: command stdout nondeterminism (ordering, trailing whitespace)
  makes read-back equality flaky.
  **Mitigation**: read-back compares typed tag values from the reader's
  envelope, not raw bytes; comparison-set semantics ride RDR 0016 (A4).
- **Risk**: eventually-consistent external state (a remote-backed query)
  makes a synchronous read-back observe staleness, training operators to
  ignore applied-but-unverified refusals (premortem P-6).
  **Mitigation**: the MVV and A3's spike bind *locally consistent* tools;
  remote-backed accessors are a declared-tool choice the model author makes,
  and the refusal's diagnosis tuple names which reader disagreed — revisit at
  the cross-cutting gate if a retry/settle policy is ever wanted.
- **Risk**: secrets placed in declared argv are visible in process listings
  and in the model.
  **Mitigation**: the carrier is reviewable data by design — document that
  credentials belong in the tool's own config, never in `command`; revisit at
  the cross-cutting gate (secret/credential lifecycle).
- **Risk**: `argv0` resolves via `PATH`, so the same model executes different
  binaries on different hosts.
  **Mitigation**: accepted as-declared semantics — `PATH` is in C4's
  allowlist and `LookPath` runs in the parent before spawn (A2); an author who
  wants host-independence writes an absolute `argv0`. A lint advisory on bare
  `argv0` is deferred to the cross-cutting gate.

### Failure Modes

- Visible: load-time C5 defects name the entry and the offending element;
  runtime refusals carry the 0004 diagnosis tuple (accessor, capability,
  role, timeout) with class `timeout` or `execution_failure`.
- Visible: with the C6 gate off, every command invocation refuses
  `execution_failure` naming `allow_commands` before any child spawns; lint
  still validates, so the model is reviewable before it is trusted.
- Visible: a gate command with a malformed envelope refuses
  `execution_failure` — never a deny the model didn't decide.
- Silent risk: a child that ignores termination at deadline can outlive its
  `timeout` refusal (A1's spike bounds this; diagnosis: refusal class plus a
  lingering process on the artifact).
- Silent risk: a command write whose role has no command-capable reader
  refuses `read_back_incomplete` on every invocation — correct but
  surprising; diagnosis starts at the role's reader declaration. Load-time
  detection of a reader-less write role is reader-cardinality territory —
  RDR 0016's contract, not this one (A4).
- Silent risk: a write tool that sinks its stdin into the artifact ingests
  the C3 envelope as content — the A3 spike's `tee {artifact}` overwrote the
  artifact with `{"state.phase":"review"}` and exit 0; it surfaces only at
  read-back (`git config` then fails `bad config line 1` → 
  `read_back_incomplete`), never at the write. The wrapper contract (envelope
  in, tool-native invocation out) is the answer for every stdin-reading tool
  (premortem P-1/P-2).
- Visible: a command entry invoked on Windows refuses `execution_failure`
  naming the unsupported platform, before spawn. v1 supports the Unix
  platforms the process-group mechanism requires (A2, A1/C4), and Windows
  re-quotes the vector into one command line — which would void C2's
  authority bound silently. Refusing is what keeps that bound honest rather
  than platform-dependent; lint stays platform-neutral, so a model
  authored on Windows still validates.
- Recovery: fix the model entry, re-run `intrastate lint`, re-invoke;
  bindings hold no state between invocations (`0004:C14` — no retry, no
  undo).

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions verified (A1–A9)
- [x] A3 spike verdict on the envelope vs raw single-value read mode — raw
      mode adopted (C3 `output = "raw"`)
- [ ] RDR 0016 lands first (fail-closed `readerFor`), or 0025 records that
      command-write read-back inherits first-match until it does (A4)

### Minimum Viable Validation

1. Author a model declaring one command-backed read and one command-backed
   write (`read_back = true`) over a single artifact role, delegating to an
   established tool present in CI: the read binds `git config --file
   {artifact} --get <key>` directly in raw mode; the write is a declared
   wrapper over `git config` (the A3 spike's `wrapper-write.sh` shape).
2. `intrastate lint` accepts it; five mutated copies (path+command conflict,
   empty element, unknown placeholder, `["sh", "-c", …]` interpreter form,
   `output = "raw"` with two keys) are each rejected with their C5 defect.
3. With `allow_commands = true` in the invoking configuration (C6; without
   it, the same invocation refuses `execution_failure` naming the gate),
   drive a state change end to end: the resolved decision invokes the
   declared write command; the tool — not intrastate — applies the edit to
   the artifact.
4. Read-back runs through the declared command reader and verifies the
   written value; the result reports the written tags.
5. A command that sleeps past its declared timeout refuses `timeout` and
   leaves no orphan process.

End state: a linted model applies and verifies a real state change through
declared commands only, with the executed argv readable in the model.

### Phase 1: Carrier

Extend `internal/table` (`sourceAcc`, `Accessor`, `accessorTable`) with
`command` and the C1/C5 load rules.

### Phase 2: Binding family

Implement command-backed `ReadBinding`/`GateBinding`/`WriteBinding` over
`exec.CommandContext` honoring C2–C4 (a sibling package to `flowbind`).

### Phase 3: Selection

Discriminate the constructor by carrier field at
`internal/cli/flowbind/registry.go::Registry`.

### Phase 4: Surface and proof

Wire lint reporting for the C5 defects and land the MVV scenario plus the
timeout and read-back failure scenarios as tests.

## Validation

### Testing Strategy

Coverage goal: every C1–C6 clause has a test that fails when its rule is
dropped; the MVV scenario is the integration proof.

1. **Scenario**: load-time validation (C1/C5) — a valid command entry plus
   the five MVV mutants (path+command conflict, empty element, unknown or
   partial `{…}` token, `["sh", "-c", …]`, raw output with two keys).
   **Expected**: the valid entry loads; each mutant is rejected with its own
   named C5 `table.Category` — asserted by category, never by "validation
   returned non-empty"; all five categories are members of
   `table.Categories()` (the registration is what puts them in the closed
   set — a per-mutant refusal assertion passes without it);
   `accessor.ValidationCodes` stays at eight (REQ-84).
2. **Scenario**: substitution guard (C2) — `{artifact}` bound to an absolute
   path, a relative path, and a `-`-prefixed path; an element embedding the
   token mid-string.
   **Expected**: the absolute path is substituted whole-element and the argv
   the child observes equals the declared vector otherwise byte-for-byte; the
   relative and `-` cases refuse `execution_failure` before spawn; the
   mid-string case is a C5 defect at load.
3. **Scenario**: envelope-before-exit ordering (C3/C4) — a gate command that
   writes a well-formed deny envelope and exits non-zero; one that writes a
   malformed envelope and exits zero; one that exits non-zero with no
   envelope and no exit map; the same with `exit_verdicts = { "1" = "deny" }`;
   a gate exiting 2 under that map; a read exiting 1 with `exit_absent = [1]`
   and one exiting 128; a 1 MiB + 1 byte stdout.
   **Expected**: deny, `execution_failure`, `execution_failure`, deny,
   `execution_failure`, established-absent, `execution_failure`,
   `execution_failure` respectively — the verdict/error code split follows
   normative fixture **FX-exit-codes** (`evidence/spikes/a3-a7-a8-real-tools/
   output.txt` E1a–E1h, R6) — each refusal carrying the 0004 diagnosis tuple
   and the 4 KiB stderr tail.
4. **Scenario**: deadline (C4, A1) — a child sleeping past `timeout`, a
   wrapper whose grandchild holds the stdout pipe, and a child that never
   reads stdin (1 MiB stdin); plus the ablations: no process group, no
   `WaitDelay`.
   **Expected**: `timeout` refusal within `timeout + WaitDelay` in all three,
   no hang, and no process from the child's group surviving the refusal —
   normative fixture **FX-deadline** (`evidence/spikes/a1-deadline/output.txt`
   S1–S5: S2/S4/S5 ≈1.0s with 0 survivors; the ablations reproduce S3's
   orphan and S1's blocked `Wait`, which the test asserts as failures of the
   ablated build, not of the binding).
4b. **Scenario**: raw read and env policy (C3/C4, A3/A7) — `git config --file
   {artifact} --get <key>` in raw mode; the same entry with an inherited
   `GIT_CONFIG_COUNT` injection in the parent env and no `env` declaration.
   **Expected**: the value is stdout minus one trailing newline — normative
   fixture **FX-raw-read** (`output.txt` R1: stdout `"draft\n"` → `draft`);
   the injection never reaches the child (its env holds only the allowlist),
   so the read is unchanged.
5. **Scenario**: write read-back (C3/C4) — a command write whose wrapper
   applies the planned tags; one whose tool exits zero without applying them
   (`git config` with the value on stdin, spike R3b); one whose argv0 sinks
   stdin into the artifact (`tee {artifact}`, spike R4b); one planning
   `<clear>`.
   **Expected**: success reporting the written tags; `read_back_mismatch`;
   `read_back_incomplete` (the corrupted artifact no longer parses for the
   reader); the key reads back absent — the write's exit status decides none
   of them.
5b. **Scenario**: read-back reader selection under two readers on one role
   (A4, Prerequisites) — a role with a path-backed and a command-backed
   reader, both admissible to today's first-match `readerFor`.
   **Expected**: pending 0016's fail-closed uniqueness, the selection is
   `registry.go`'s name-sorted first match — asserted **by the selected
   reader's identity**, never by "read-back succeeded", since either reader
   can verify a correct write and the test would pass under either. This
   pins the inherited behaviour so 0016's landing is a visible change, not a
   silent one.
6. **Scenario**: the MVV end to end against an established tool present in
   CI.
   **Expected**: the declared write command, not intrastate, edits the
   artifact, and the declared command reader verifies it.
7. **Scenario**: execution gate (C6) — a valid command model invoked with the
   gate unset; with the `--allow-commands` flag alone; `intrastate lint`
   under both. Once A10 fixes the config surface, the same scenario adds
   `allow_commands = true` in config, and the flag-vs-config precedence case
   A10 decides.
   **Expected**: refusal `execution_failure` naming `allow_commands` with no
   child process spawned — asserted by **absence of a spawn**, not merely a
   non-zero exit; normal execution under the flag; lint passes in both
   (validation is ungated). A malformed config must refuse, never fail open.

### Pre-Lock Mini-Checks

Cue-fired at Stage 5. Four of five fired; **test-discriminability did not**
(every Expected line above names a discriminating value, and S1 already bans
the absence-of-error oracle).

**`authority`** — source-of-truth census (cue: ≥2 candidates; a resolver; sibling arms).

| Input / decision | Canonical writer | Readers | Sibling arm | Which is canonical |
| --- | --- | --- | --- | --- |
| Load-time defect codes | `table.Category` + `Categories()` | `intrastate lint`, loader | `accessor.ValidationCode` (8, closed) | `table.Category` — A6; the 0004 set is 0004's own test contract and does not grow |
| Child environment | C4 allowlist | executor at spawn | `env_pass`, entry `env`, `INTRASTATE_*` overlay | C4's four-part composition is total and ordered; nothing else is inherited |
| Execution permission | user-scope config (A10, new) | binding, pre-spawn | `--allow-commands` flag | Either grants; absence refuses. Precedence is A10's to fix |
| Read result | child stdout envelope (C3) | command binding | `flowbind.go::store` (artifact file) | Disjoint transports, one shared map-of-strings rule (Wire LBD) |
| Gate verdict | child stdout envelope + C3 exit map | command binding | `flowbind.go::verdictFor` (path suffix) | Disjoint by carrier: suffixes are path-backed only (Selection LBD) |
| Write success | read-back through the role's reader | executor | child exit status | Read-back only — exit status decides no write (`0004:C12`) |

**`fidelity`** — round-trip / inverse invariants (cue: parse/deparse, write→read-back).

| Operation | Invariant | Exemption |
| --- | --- | --- |
| declared argv → executed argv | byte-for-byte equality except the one whole-element `{artifact}` substitution (C2) | none — this is the authority bound |
| planned tags → stdin JSON → tool → read-back | typed **tag-value** equality, not byte equality; `<clear>` round-trips as absence | the artifact's byte form is the tool's, not intrastate's (Perf Expectations) |
| child stdout → tag map | `raw`: value = stdout minus exactly one trailing `\n`; `json`: omitted key ⇒ not-carried (never empty-string) | non-inverse by design: intrastate never re-emits the child's bytes |

**`disposition`** — input class → outcome (cue: exit codes, refusal classes).

| Input class | Refusal / verdict | Loud? |
| --- | --- | --- |
| non-empty stdout, well-formed | parsed result — wins over exit code (C4 ordering) | loud |
| non-empty stdout, malformed | `execution_failure` | loud |
| empty stdout, exit in `exit_verdicts` | that verdict | loud |
| empty stdout, exit in `exit_absent` | keys established absent | loud |
| empty stdout, unlisted non-zero exit | `execution_failure` | loud |
| stdout > 1 MiB | `execution_failure` | loud |
| deadline exceeded | `timeout` (from `ctx.Err()`), group killed | loud |
| spawn failure / gate off / Windows | `execution_failure`, `Detail` names the cause | loud |
| write applied, read-back unavailable | `read_back_incomplete`, `Applied()` true | loud — never silent success |

**`trace`** — desk trace of the MVV end state (cue: C2/C3/C4 all bear on the invocation result).

| MVV step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1 author model | C1 exactly-one; C5 five defects registered | `read.branch_state` with `command`, no `path` | consistent |
| 2 lint accepts, 5 mutants rejected | C5 + `Categories()` membership (S1) | each mutant → its own category | consistent |
| 3 gate off → refuse | C6; A10 (flag alone suffices) | `execution_failure`, no spawn | consistent |
| 3 gate on → invoke write | C2 absolute-path substitution; C4 env allowlist | argv = declared vector + absolute `{artifact}` | consistent |
| 4 read-back through declared reader | C3 raw mode; C4 write rule; A4 `readerFor` | `git config --get` → `draft\n` → `draft` (FX-raw-read R1) | consistent — **single reader**; the two-reader case is S5b, unpinned until 0016 |
| 5 sleeper past timeout | C4 deadline triple; A1 | ≈1.0s, 0 survivors (FX-deadline S2/S4/S5) | consistent |
| end state | all of the above jointly | linted model applied + verified a real change | **no CONTRADICTION** |

### Performance Expectations

Spike-measured, darwin/arm64, Go 1.26.6, git 2.55 (`evidence/spikes/`):
a `git config --file` read or write invocation costs 7–15ms wall including
spawn; a 2 MiB stdout envelope captures in 88ms (C4 caps at 1 MiB); the
deadline overshoot is ≤4ms past `timeout` when the group dies on signal and
at most `timeout + WaitDelay` (500ms) when a descendant holds a pipe. No
byte-stable output is claimed: read-back compares typed tag values, not
bytes (Risks), so the determinism checklist does not apply.

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- Prior-art record with quotes and search budget:
  `docs/rdr/0025-command-invoking-accessor-bindings/evidence/research/prior-art.md`;
  Resolve source-search record: `…/evidence/research/source-search.md`;
  Resolve prior-art reversal ledger:
  `…/evidence/research/prior-art-resolve-{cli-a,cli-b,sm,devref-go}.md`
  (helm, gh-cli, goreleaser, kubebuilder, hugo, golangci-lint, consul,
  opentofu, etcd, roborev, beads, semgrep rules; 32 peer state-machine
  tools; Go toolchain `generate`/`ErrDot`/`GOVCS`/cgo-flag allowlist;
  corpus-cited Unix/security references).
- Spikes: `…/evidence/spikes/a1-deadline/` (A1; `run.sh`, `main.go`,
  `output.txt`, `notes.md`) and `…/evidence/spikes/a3-a7-a8-real-tools/`
  (A3/A7/A8; `run.sh`, `main.go`, `wrapper-*.sh`, `output.txt`, `notes.md`).
- Go 1.26.6 `os/exec` package doc, `Cmd.Args`, `Cmd.Cancel`, `Cmd.WaitDelay`
  (`$GOROOT/src/os/exec/exec.go`); `k8s.io/api@v0.36.2/core/v1/types.go`
  `ExecAction.Command`.
- Docker Dockerfile reference (exec form vs shell form); git-difftool
  documentation (`difftool.<tool>.cmd`); Terraform `external` data source
  documentation (hashicorp/terraform-provider-external).
- Peer checkouts read: awf-cli (`workflow-syntax.md`, security-validator
  plugin README), ms-conductor (`yaml-schema.md` §Template Syntax) — via the
  StateMachineRes corpus / `state-machines` sibling checkout set.
- Source reviewed: `internal/accessor/binding.go`, `internal/accessor/model.go`,
  `internal/accessor/executor.go`, `internal/cli/flowbind/flowbind.go`,
  `internal/cli/flowbind/registry.go`, `internal/table/source.go`,
  `internal/table/load.go`.
- Related: RDR 0002 (TOML carrier), RDR 0004 (execution safety), RDR 0016
  (reader cardinality, intrastate#p63c); kata intrastate#v0hb.
