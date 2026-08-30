# Recommendation 0025: Command-Invoking Accessor Bindings

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Implemented
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
the table was built to remove. This RDR fixes the *declaration* half: the edit
becomes a linted, declared command. The caller still carries the decision from
resolve to apply by hand — today's `--write`-flag set-state path
(`internal/cli/flow_state.go`) is the only write entry point, and automatic
resolve→apply wiring is a successor's scope, not this record's. What changes
is that the edit it drives is declared and bounded rather than ad hoc. The read side has the same hole: state only a
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
  envelope. For the tool shapes the spike covers — key-value query, file
  probe, status verb, key-value write — that wrapper is small and mechanical
  (1–8 lines). The claim is scoped to those shapes and is NOT a claim about
  every integration: a tool with transactional, locking or partial-failure
  semantics is outside what the spike tested and outside what a flat
  string-map envelope can express]**
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
    Sample bound: four tools of one family (`git config` ×3 modes, `test -f`,
    `git diff --quiet`) — favorable shapes, chosen for being established and
    locally consistent. The spike does not evidence tool classes with
    transactional or partial-failure semantics (a DB CLI, `kubectl patch`, a
    state-locking tool); for those the wrapper is unmeasured, and the flat
    string-map envelope carries no way to express partial application.
  - **If wrong**: every integration needs a bespoke wrapper script — the
    delegate-to-established-tools thesis fails and the choice tips toward the
    adapter-registry alternative. The bounded version of being wrong, and the
    likelier one: the wrapper stays mechanical for key-value-shaped tools and
    becomes bespoke past them, which costs the *breadth* of the thesis, not
    its core — the QOC "integration cost" row is scored on the shapes the
    spike measured and would drop for a corpus weighted toward transactional
    tools.
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
- **A11 [A NET-NEW typed `*accessor.ExecError` (C4 specifies it; no such type
  exists on `main`) returned through the existing bare `error` returns carries
  the stderr tail to a NET-NEW `Refusal.Detail` field without changing the
  binding interfaces or any path-backed binding]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `binding.go`'s `Read`/`Gate`/`Apply` return a
    bare `error`, which the executor discards, keeping only
    `ClassExecutionFailure` (confirmed on `main`). Note the executor is
    **asymmetric**: only `invokeRead` exists as a helper (also reused for
    read-back); `Executor.Gate` and `Executor.Write` handle their binding
    errors inline. The critique lens refuted the single-hook form of this
    claim: `executor.go::refusalOf` is the base constructor but **not** the
    sole post-selection one — `executor.go::refusalWithKeys` wraps it and is
    what `Executor.Read` uses for *every* read refusal, so an error parameter
    on `refusalOf` alone reaches gate and write but silently drops the tail on
    reads (`executor.go::selects` builds its own pre-selection refusal
    directly, which can carry no invocation error). The claim is now that
    threading the parameter through **both** constructors reaches every
    invocation refusal and leaves every existing binding and all production
    `NewExecutor` callers behaving identically. **Verified**: `refusalOf` has
    16 call sites, all in `executor.go`; `refusalWithKeys` wraps it and has
    exactly two, both the refusal returns in `Executor.Read`. Every site can
    supply the offending error or `nil` — the pre-selection sites in
    `::Executor.selects` pass `nil`, and the unbound/mismatch refusal there is
    a bare `&Refusal{...}` composite literal bypassing both constructors, which
    can carry no invocation error (structural: no binding is selected yet).
    No other consumer branches on the discarded error — it never leaves
    `executor.go`: `internal/accessor/disposition.go`'s three dispositions read
    only `Refusal.Class`, and `internal/cli/flow_exec.go::accessorFailure`
    switches on `Class` and reads `Keys`/`Expected`/`Observed`. All three
    production `accessor.NewExecutor` call sites
    (`internal/cli/flow_exec.go` ×2, `internal/cli/flow_state.go`) are
    unaffected.
    **Three net-new facts the build must carry.** (1) `accessor.Refusal`
    (`internal/accessor/model.go`) has fields `Class, Accessor, Capability,
    Role, Timeout, Keys, Expected, Observed, Reason` and unexported `applied` —
    there is **no `Detail` field**; it is additive (Capability Dependencies,
    C4). (2) `accessor.ExecError` **does not exist anywhere in the repo**, and
    no production code imports `os/exec` — C4 specifies the type this RDR
    mints, it is not an existing type being threaded. (3) `readOutcome` drops
    the error today (`return readOutcome{class: ClassExecutionFailure}`), so it
    gains an `err` field for the read path to reach `refusalWithKeys` — the one
    structural change inside the executor beyond the two signatures.
  - **If wrong**: the stderr tail needs an interface change on all three
    binding interfaces, which contradicts C4's "0004 executor unchanged"
    framing and widens the blast radius to every path-backed binding.

- **A12 [The CLI callers that open the model file can pass its directory to
  the command-binding constructor for C2's argv0 resolution, leaving `table`
  and in-memory model construction untouched]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the original form of this assumption — "the loader records
    the directory on `table.Model`" — was **refuted** by the critique lens:
    `internal/table/load.go::Load` takes `(src []byte, sourceID string)`,
    performs no file I/O and no path resolution by charter (`0002:EIA`), and
    `sourceID` is a display string for line locators, never stored as a path.
    The loader never sees a path, so it cannot record one. Re-scoped to the
    layer that does hold it: `internal/cli/flow_input.go::selectModel` and
    `internal/cli/lint.go::runLint` each `os.ReadFile` a `--model` path
    (confirmed on `main`), and `flow_exec.go::buildRequest` is where the
    registry is built from the model. The claim is that the base dir threads
    from that caller to the command-binding constructor as a parameter, so no
    `table.Model` field is added and no in-memory construction or fixture
    changes.
    **Verified, and the thread is shorter than claimed.** The caller and the
    constructor are the SAME function body: `buildRequest` opens with
    `model, ref, ce := selectModel(cmd)` — `ref` IS the `--model` path, from
    `selectModel`'s third return — and constructs the registry
    (`registry: flowbind.Registry(model)`) three statements later. The path is
    already persisted as `flowRequest.modelRef` and read by four downstream
    verbs, so it is a pre-existing carried value, not one this RDR introduces;
    the four `buildRequest` callers need no change, since the threading runs
    INWARD from `buildRequest` to `Registry`. Confirmed alongside: `Load` is
    `func Load(src []byte, sourceID string) (*Model, error)` with no file I/O
    (no `os.ReadFile`/`Open`/`Stat` anywhere in `internal/table`), `sourceID`
    lives only on the unexported `loader`, and `table.Model` carries no path
    field — the only `Path` in the package is `Accessor.Path`, the per-entry
    declared locator. No layering obstacle: `flowbind` already imports
    `internal/accessor`, `internal/table`, `os` and `path/filepath`, and
    neither imported package imports `internal/cli` (pinned by a test in
    `internal/table`), so `filepath.Abs` then `filepath.Dir` are available
    with no new dependency and no cycle.
  - **If wrong**: separator-bearing argv0 is dropped from v1 in favour of
    PATH-resolved bare names and absolute argv0 only — C2's resolution rule
    loses one arm, and no other clause changes.
- **A13 [The `--allow-commands` gate reaches every execution path through one
  construction site, and the verbs that must resolve the flag are exactly
  the callers of `flow_exec.go::buildRequest` — covered by ONE registration
  on the `flow` group, whose subtree is exactly that set]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: opened by the critique lens, which found the
    RDR, and both critique passes, undercounting this surface. Established on
    `main`: `flowbind.Registry` is a free function
    (`func Registry(m *table.Model) accessor.Registry`) with exactly ONE
    production call site, `flow_exec.go:830` inside `buildRequest`, and it
    constructs `Reader{Path: acc.Path}` / `&Writer{...}` / `Gate{...}` for
    every entry with no carrier discriminator; `buildRequest` has FOUR
    production callers (`flow_next.go:207`, `flow_resolve.go:227`,
    `flow_state.go:135`, `flow_state.go:274`) across three verb files. The
    claim is that gating at that one construction site covers every executor
    built from the shared `flowRequest.registry`, so flag registration is the
    only drift surface.
    **Verified, and one arm corrected.** No other production path constructs
    an `accessor.Registry` or reaches a binding: the only `accessor.Registry{}`
    literals are `registry.go` and test fixtures; the only
    `Reader{`/`&Writer{`/`Gate{` constructions are `registry.go` ×3 plus
    `flowbind/persistence_0005_test.go`; all four `flowRequest` literals live
    inside `buildRequest`, three of them the zero-value error return, so
    `flowRequest.registry` can only ever hold a `flowbind.Registry(model)`
    value, and all three production `NewExecutor` sites read that shared field.
    The corrected arm is the ORIGINAL claim's last clause — "the flag lookup on
    a verb that never registered it returns false rather than erroring". That
    is imprecise: pflag's `GetBool` on an unregistered flag returns
    `(false, *NotExistError)` — an error, never a panic — whose text is "flag
    accessed but not defined", i.e. the library calls it a PROGRAMMER ERROR,
    not an answer about user intent. The repo's idiom happens to discard it
    (`planOnly, _ := cmd.Flags().GetBool(...)`, `internal/cli/flow_projection.go`),
    so a miss would read `false` in practice — but that makes fail-closed a
    coincidence between two independent decisions, not a property. C6 therefore
    no longer relies on it: ONE registration on
    `internal/cli/flow.go::newFlowCmd`'s `PersistentFlags()` — whose subtree is
    exactly the four `buildRequest` call sites and nothing else, with `lint` at
    root outside it — and the lookup error CHECKED at `buildRequest` as a
    wiring bug. Precedent is in-repo (`internal/cli/root.go` registers the
    shared output-mode flag persistently and reads it via `Lookup` + nil check)
    and external (`evidence/reconcile/prior-art-absence-and-gating.md`: eleven
    mature Go CLIs, none gating on a lookup miss; etcd registers persistently
    AND checks the error).
  - **If wrong**: the gate needs a second site per bypassing path, and C6's
    "no execution path can bypass the gate" claim narrows to the paths
    actually covered — the refusal stays correct, its completeness does not.
- **A14 [The command-binding constructor can refuse a carrier-less entry
  instead of constructing a `Path: ""` file binding, and nothing today
  depends on that fail-open behaviour]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the hazard is established on `main`:
    `flowbind.go::load` maps a non-existent file to an EMPTY store by design
    (so a flow's first write can establish state) and `os.ReadFile("")`
    returns `ErrNotExist`, so an entry that reached the file binding with an
    empty path would read every declared key as absent and confirm an
    unapplied write — silently. C1's load-time exactly-one rule makes that
    residue unreachable through the loader, but a binding built from an
    in-memory model never passed the loader. The claim is that the
    constructor can refuse it without disturbing the deliberate
    empty-artifact semantics for genuinely path-backed entries.
    **Verified.** No production path constructs an accessor entry with an empty
    `Path`: `flowbind/registry.go::Registry` sets `Path: acc.Path` from
    `table.Accessor.Path`, and `internal/table/load.go::accessorTable` refuses
    `a.Path == nil || *a.Path == ""` as `malformed accessor declaration` before
    a model exists. A repo-wide sweep for `Path: ""` / `Path = ""` literals in
    `internal/` returns ZERO hits in production and in tests alike — the only
    literal `Path:` values in tests are named strings — so no test depends on
    empty-`Path` construction. The two cases are cleanly separable in source: a
    first-write flow carries a NON-EMPTY declared `Path` (the model's locator,
    which `unreachable()` inspects and which is never the filesystem) while the
    caller's `accessor.Artifact.Path`, set from `--artifact role=path`, names a
    file that may not yet exist. `load` operates on the CALLER's `art.Path`, so
    refusing on the DECLARED `Path` in the constructor cannot disturb the
    deliberate empty-artifact semantics. Note `internal/accessor/validate.go::Validate`
    has no empty-path arm among its eight codes, so `accessorTable` is the sole
    existing refusal — exactly the loader-only coverage C1 assumes, which is
    what makes the constructor-side refusal load-bearing rather than redundant.
  - **If wrong**: the residue arm is dropped and C1 relies on the loader
    alone, leaving in-memory models able to build a fail-open binding — which
    must then be named as a Failure Mode rather than refused.

- **A10 [No user-scope configuration surface exists in this repo, so
  `--allow-commands` alone is C6's complete v1 gate and the config file is a
  successor's scope]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the absence is established — no `internal/cli/config/`
    directory; zero `intrastate.toml` or `allow_commands` hits in non-test
    `*.go`; no user-scope TOML loader of any kind. A flag is therefore the
    only opt-in v1 can carry without building a new discovery/precedence/
    malformed-file surface, which is net-new scope (charted: `evidence/3amigo/
    Charted:` — user-scope config for `allow_commands`). Sufficiency: the gate
    is a boolean and the flag sets it, so absence refuses and presence
    executes — the property C6 needs. What v1 gives up is the *one-time*
    quality of the opt-in, recorded as a Consequence, not the gate.
  - **If wrong**: nothing in C6 changes; a later config file composes
    disjunctively with the flag (either grants, neither revokes), which is why
    the successor inherits an already-gated seam.
- **A15 [C3's read envelope maps onto `ReadBinding.Read`'s existing split
  return without an interface change: a parsed key becomes a `KeyValue` in
  `values`, an omitted declared key a name in `unreadable`, and `exit_absent`
  absence a `KeyValue{Absent: true}` in `values` — whole-invocation, every
  declared key at once, the invocation being the unit; per-key absence needs
  the JSON envelope — and `WriteBinding.Invocations()` has a defensible
  answer for a command binding]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: to verify — opened by the repeatability lens, where all
    three runs invented a `Read` signature (two rendered it
    `(map[string]string, error)`) because the RDR named the seam without
    quoting it. Established on `main`:
    `internal/accessor/binding.go::ReadBinding` is
    `Read(ctx context.Context, art Artifact, requested []string) (values []KeyValue, unreadable []string, err error)`,
    `::GateBinding` is `Gate(ctx, art) (Verdict, string, error)`, and
    `::WriteBinding` is `Apply(ctx, art, planned []resolve.Tag) error` plus
    `Invocations() int`. C7 now quotes them. UNREADABLE rides the `unreadable`
    slice (not an error, not a sentinel value) — confirmed.
    **The omission arm was REFUTED and is repaired here.** The claim as first
    written — established-absent expressible as omission from BOTH slices —
    is false against source: `internal/accessor/executor.go::readOutcome.classify`
    iterates the DEFINITION's requested set and sends a key found in neither
    slice to `unread` ("A key the binding classified as neither read nor
    unreadable defaults to UNREADABLE: guessing absence is the failure this
    contract exists to prevent", `0004:C6`), whereupon
    `::Executor.Read` refuses `ClassIncompleteRead`. Pinned by
    `internal/accessor/read_completeness_0004_test.go::TestReq29_UnclassifiedKeyDefaultsToUnreadableNotAbsent`,
    which drives a binding returning a key in neither list and asserts the
    refusal. Shipped as first worded, every `exit_absent` read would have
    refused `incomplete_read` — the opposite of the intent.
    The repair, consult-PASSed: absence crosses as `KeyValue{Key: k, Absent:
    true}` in `values`. `internal/accessor/binding.go::ReadBinding.Read`'s own
    doc says it returns "one KeyValue per key it could read (present or
    **established-absent**)", and `::KeyValue.Absent` is documented as "a VALUE
    (the binding read the artifact successfully and the key was not there)".
    An Absent record survives `classify` intact: the `clearIsUnreadable` arm is
    guarded by `!v.Absent`, and the tail propagates `Absent: v.Absent` verbatim.
    This is the channel the shipped path-backed binding already uses —
    `internal/cli/flowbind/flowbind.go::Reader.Read` emits
    `KeyValue{Key: key, Value: v, Absent: !held}` for every requested key — so
    the command binding matches it rather than inventing a second convention.
    **Two boundaries, both intact**: at the BINDING return absence is the flag;
    at the RESOLVER seam it becomes omission, converted by
    `internal/accessor/model.go::ReadResult.OwnedSnapshot` ("Omission, never a
    placeholder: a sentinel would make `TagSet.has` true and silently retire
    the refusal") and again at `internal/cli/flow_exec.go` for reported tags.
    `0004:C8`'s "what crosses the seam omits the key entirely" governs that
    OUTER boundary only — `KeyValue.Absent`'s doc states the split in terms
    ("the accessor layer's INTERNAL representation; what crosses the seam omits
    the key entirely (`0004:C8`)"). Adjudication:
    `evidence/reconcile/two-boundary-absence.md`.
    Granularity: `exit_absent` is a WHOLE-INVOCATION signal (C3 — "a listed
    exit with empty stdout establishes every declared key absent"), sound only
    because C3 admits the exit maps solely on EMPTY stdout, the one case the
    per-key envelope structurally cannot express; a per-key absence requires
    the envelope. That restriction is load-bearing, not incidental.
    Prior art (`evidence/reconcile/prior-art-absence-and-gating.md`): the
    explicit-flag form is the established pattern (`sql.Null[T].Valid`; cty's
    `IsNull()`/`IsKnown()` split; Consul KV's two-channel 404; gorm's
    `Nullable() (nullable, ok bool)`), and omission-as-absence is the
    documented failure mode proto3 reversed with `optional` in v3.12/v3.15.
    The fail-closed default matches OpenTofu's `ErrIsNotExist`, which admits
    absence only on "an affirmative response".
    `Invocations()` has **zero production consumers** — declared at
    `internal/accessor/binding.go` and implemented once at
    `internal/cli/flowbind/flowbind.go::Writer.Invocations`; all other uses are
    test assertions of `== 0` / `== 1` proving the accessor layer performed no
    retry, undo, or re-derivation (`0004:C14`). Counting `Apply` ENTRIES —
    one per call, with no internal respawn — satisfies every existing assertion
    shape (C7 states the constraint normatively).
  - **If wrong**: if the two slices cannot express the three read outcomes, C3
    needs a third outcome channel and C7's "no interface change" claim fails —
    widening the blast radius to the path-backed bindings C4 promises to leave
    untouched. If `Invocations()` has no sound answer, the write binding needs
    a documented constant and a Failure Mode.

## Proposed Solution

### Approach

The undecided contract at this seam: **what carries a declared accessor's
real external behaviour**. Today `flowbind` *simulates* command behaviour
through magic `path` suffixes (`internal/cli/flowbind/flowbind.go::verdictFor`,
`::unreachable`) because no declared-command carrier or authority bound
exists — the gap both accreted point-fixes (Background) patched around.

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
command       = ["<argv0>", "<arg>", ...]   # []string
output        = "json" | "raw"              # *string, read entries only; default "json" when omitted
exit_absent   = [<code>, ...]               # []int,            read entries only
exit_verdicts = { "<code>" = "<verdict>" }  # map[string]string, gate entries only
env           = { "<KEY>" = "<value>" }     # map[string]string
env_pass      = ["<VAR>", ...]              # []string

exactly one of `path` / `command` per entry. Load-time this is C5 `command_and_path_conflict`; at RUNTIME the constructor refuses a carrier-less entry with a refusing binding (`execution_failure`, Detail naming the malformed entry) and MUST NOT fall through to a `Path: ""` file binding, which would read every declared key as absent and confirm an unapplied write. In-memory models never pass the loader, so the runtime arm is the one that closes this — not a duplicate of C5
```

The TOML spellings above are the contract; the Go field names and struct tags
on `internal/table/source.go::sourceAcc` are unconstrained implementation
choice, as is the decomposition of the binding into helpers and the package or
file that holds it. This RDR fixes the wire, the carrier grammar, the load-time
vocabulary and the seam it implements — not the module layout.

An accessor entry carries **exactly one** of `path` or `command`; both or
neither is a load-time defect. `command` must be non-empty and contain no
empty element. `role`, `keys`, `timeout`, and `read_back` rules are unchanged
from RDR 0002/0004.

The selection in `flowbind.Registry` must be **exhaustive, and fail closed on
the residue**. Today it constructs `Reader{Path: acc.Path}` /
`&Writer{Path: acc.Path}` / `Gate{Path: acc.Path}` for every entry with no
carrier discriminator, and the file binding **fails open** on an empty path:
`flowbind.go::load` maps a non-existent file to an EMPTY store by design (so
a flow's first write can establish state), and `os.ReadFile("")` returns
`ErrNotExist`. A command entry that reached the file binding would therefore
read every declared key as *absent* and confirm an unapplied write as
verified — silent state loss, never a crash. So the constructor selects
`command` first, `path` second, and **panics-free refuses** the third case:
an entry with neither carrier builds a refusing binding
(`execution_failure`, `Detail` naming the malformed entry), never a
`Path: ""` file binding. C1's load-time exactly-one rule makes that residue
unreachable through the loader; the constructor refuses it anyway, because a
binding built from an in-memory model never passed the loader. S1 asserts the
load-time rule; the constructor's refusal is asserted by S7's sibling arm.

The six fields above are the complete set this RDR adds to
`internal/table/source.go::sourceAcc`, with the Go types named — the loader
decodes strictly (`source.go::decodeStrict` sets `DisallowUnknownFields`), so
a field absent from the struct is a **document-level**
`CatUnknownSchemaField`, not one of C5's per-entry categories. Declaring them
on the carrier is therefore a precondition of C5 reporting anything: the
mutants S1 asserts (`output = "raw"` with two keys, `exit_absent` on a
non-read) only reach `accessorTable` once the fields decode. `output` is a
pointer so that omitted and explicit-`"json"` are distinguishable at load,
where the C5 `command_output_shape` arms are judged.

The C5 conflict/empty arms **re-categorize an existing refusal**: today
`load.go::accessorTable` emits `CatMalformedAccessorDeclaration` for "path is
absent or empty" (A9). A command entry is now legal there, so the
path-absent-and-command-absent case moves to `command_and_path_conflict`.
Existing fixtures asserting the old category on a path-less entry change
category, not verdict — S1 asserts the new spelling.

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
absolute path — the **command binding** refuses the invocation
(`execution_failure`, before spawn) when the caller-bound artifact path is
relative or begins with `-`, so a path can never be parsed as a flag by the
invoked tool. The check sites in the binding, not in `executor.go`, which
keeps C4's "0004 executor unchanged" claim true; and it **refuses rather than
absolutizes**, because `internal/cli/flow_input.go::parseArtifacts` stores
`--artifact` paths verbatim and intrastate's cwd is not the tool's frame of
reference — silently resolving against the wrong base is the failure the
refusal exists to prevent. The cost is stated: a relative `--artifact` that
works for a path-backed entry today refuses for a command entry, and the
refusal `Detail` says so ("artifact path must be absolute for a command
entry").

argv0 resolution is fixed at load, never cwd-relative: a bare name resolves
through the parent's `PATH` at spawn (A2); a name containing a path
separator resolves against the **model file's directory** (the git-hook /
pre-commit convention); an empty argv0 is a C5 defect, and the binding never
restores implicit current-directory lookup — Go itself shipped cwd-relative
`PATH` resolution for a decade and reversed it (`exec.ErrDot`, Go 1.19;
reversal ledger, `evidence/research/prior-art-resolve-devref-go.md`).

The model file's directory is not reachable from the loader, and cannot be:
`internal/table/load.go::Load` takes `(src []byte, sourceID string)` and
performs no file I/O and no path resolution (`0002:EIA`); `sourceID` is a
display string for line locators, never a resolved path. The **CLI** owns the
read — `internal/cli/flow_input.go::selectModel` and
`internal/cli/lint.go::runLint` each `os.ReadFile` a `--model` path they hold
— so the base dir is carried at that layer, not on `table.Model`: the
command-binding constructor receives `filepath.Dir` of the **absolutized**
model path from the caller that opened the file, alongside the model. This
keeps `table` free of path resolution and leaves in-memory model construction
untouched — no `table.Model` field, so no constructor or fixture changes.
`--model` is stored verbatim (it is not absolutized today), so the caller
absolutizes before taking `Dir`; a relative `--model` must not make argv0
resolution cwd-dependent. **This threading is A12, still `Pending`**: the
clause states the shape the fix takes, not a verified reachability — if the
model path does not reach the registry construction site without a new field,
A12's "If wrong" applies and separator-bearing argv0 leaves v1, costing this
paragraph and no other clause.

A model with no source file (built in memory, or an embedder's) has no base
dir. A separator-bearing argv0 in one is refused `execution_failure` at
**binding construction**, before any spawn — not a load-time defect, because
C5's registered set has no category for it and the base dir's absence is not
a property of the declaration the loader validates. What it must never do is
fall back to the process cwd, which is the reversal this clause exists to
honour.

**C3**

```normative
stdin (write): {"<key>": "<value>", ...}   # the planned tags; a planned `<clear>` crosses as the literal reserved value
stdin (read/gate): {}                       # nothing per-invocation beyond {artifact}; the object is always sent
stdout (read, output = "json", default): flat JSON object of strings; a declared key the object omits is UNREADABLE, never established-absent
stdout (read, output = "raw"):            the single declared key's value = stdout minus one trailing "\n"; valid only when `keys` has exactly one entry
stdout (gate): {"verdict": "allow" | "deny" | "indeterminate", "reason": "<text>"}
exit_absent   = [<code>, ...]                # read entries: a listed exit with empty stdout establishes every declared key absent
empty stdout, no exit-map match: execution_failure, on BOTH read modes and on gate — never UNREADABLE and never established-absent. A silent tool is a broken tool, not an answer: `exit_absent`/`exit_verdicts` are the ONLY way an empty stdout carries meaning, so a read whose exit is unlisted (exit 0 included) refuses, exactly as `output = "raw"` already does
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
laundered into a verdict. A spawn failure has no exit code at all
(`exec.ErrNotFound`, fixture E1h), so no map entry can match it however the
map is written — the maps are consulted only for a process that ran and
exited. A non-empty stdout is always parsed first (C4
ordering); the exit maps apply only to an empty stdout — including a stdout
truncated at the 1 MiB cap, which is *non-empty* and therefore an
`execution_failure` from the bound (C4), never an exit-map consultation.

Omission is unreadable, not absent, because the two are different answers in
`internal/accessor/executor.go::classify` and only one of them is safe to
guess: a tool that fails to report a key has not established that the key has
no value. Establishing absence requires a positive declaration — `exit_absent`
(this clause) — which is why that map exists for exit-speaking tools at all.
The declaration is positive on the wire AND in the return: absence reaches the
executor as `KeyValue{Absent: true}` in `values`, never as omission from both
slices, which `classify` reads as unreadable (C7).
Under `output = "raw"` the single declared key is established from stdout,
and an empty stdout with an unlisted exit is `execution_failure`, not an empty
string. The inherited reserved-literal rule composes unchanged **on the read path
only**: a value that reads back as the literal `<clear>` is UNREADABLE
(`internal/accessor/executor.go::readOutcome.classify`, the
`clearIsUnreadable` arm) — newly reachable now that raw-mode reads carry real
tool output, and asserted by S5. The rule is parameterized, not universal:
`Executor.Read` passes `clearIsUnreadable = true`, while `Executor.Write`'s
read-back passes `false`, where the literal's presence is the defect
(`0004:C11`). A command-backed reader inherits whichever sense its caller
already used, so this RDR adds no arm and changes no site. The protocol's
version rides out-of-band as `INTRASTATE_PROTOCOL` in the child env (C4
overlay), keeping the stdout map flat — the envelope-versioning lesson of
kubebuilder's external-plugin `apiVersion` and client-go's exec-credential
`v1alpha1`→`v1` migration (reversal ledger).

**C4**

```normative
order:    parse the stdout envelope, THEN classify the exit code
deadline: Setpgid; at ctx deadline Cancel = SIGKILL to -pgid; WaitDelay = 500ms bounds the stdin write and the pipe drain; timeout is classified from ctx.Err()
bounds:   stdout capped at 1 MiB (overflow = execution_failure); a NEW `Detail string` field on `accessor.Refusal` carries the last 4 KiB of stderr
env:      child env = allowlisted parent vars {PATH, HOME, TMPDIR, LANG, LC_*} + named `env_pass = ["VAR", ...]` vars + the entry's literal `env = { KEY = "value" }` + the overlay {INTRASTATE_ROLE, INTRASTATE_CAPABILITY, INTRASTATE_ACCESSOR, INTRASTATE_PROTOCOL=1}; nothing else is inherited. On a key collision the LATER layer wins, in exactly that order: overlay > entry `env` > `env_pass` > parent allowlist. `LC_*` is a literal prefix match on `LC_` (the one prefix rule; `env_pass` names whole variables and admits no pattern). An `env` key or `env_pass` name matching the `INTRASTATE_` prefix is a C5 `command_env_conflict` defect at load, so the overlay is never shadowed silently
detail:   the stderr tail reaches the refusal as a typed `*accessor.ExecError` (Detail string) returned by the command binding through the existing `error` return; `executor.go::refusalOf` gains an error parameter and `errors.As`-es it. It is the base constructor, NOT the only post-selection one: `executor.go::refusalWithKeys` wraps it and is what `Executor.Read` routes every read refusal through, so the error parameter threads through BOTH or the read path — the capability that binds directly — silently drops its tail. `Detail` is empty on every refusal carrying no error (timeout, read-back, gate-off) and set only where a binding returned one. `ExecError` is `struct { Detail string; Err error }` with `Error() string` and `Unwrap() error`: it WRAPS the offending `os/exec` error rather than flattening it, so `errors.Is(err, exec.ErrNotFound)` survives the trip to the refusal site and the argv0-not-found case stays distinguishable from a non-zero exit. `Detail` is the 4 KiB stderr tail, not `Err.Error()`
write:    read_back required; verified only through the role's reader, never by exit status
platform: a runtime `runtime.GOOS` check in the command binding refuses `execution_failure` before spawn on non-Unix — the predicate is `goos == "windows" || goos == "js" || goos == "plan9"` (refuse-listed, not allow-listed), so every Unix that supports C4's process-group mechanism, including the BSDs, is admitted without enumerating it; a build constraint would make the refusal unbuildable-on-Windows rather than observable, and lint must stay platform-neutral in the same binary. The check reads an injectable package-level `goos` var (defaulting to `runtime.GOOS`) so the refusal is testable on Unix CI — unlike C4's deadline triple, overriding this one weakens nothing at runtime
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
gate deny.

The field alone is inert: the binding **interfaces** carry their failure in a
bare `error` slot (`internal/accessor/binding.go` — `Read`, `Gate`, `Apply`;
C7 quotes the full returns), and the executor discards that error, keeping
only `ClassExecutionFailure`. So the carrier is a **typed error**, not an
interface change: the command binding returns `*accessor.ExecError` in that
existing slot (C4 fixes its shape).

The hook is the **refusal constructors**, not the invocation helpers. There is
no symmetric `invokeRead`/`invokeGate`/`invokeWrite` triple to edit — only
`invokeRead` exists (also reused for read-back at two further sites);
`Executor.Gate` and `Executor.Write` handle their binding errors inline. What
all of them share is that every refusal raised *after* accessor selection is
constructed by `refusalOf` — **or by `refusalWithKeys`, which wraps it** and
which `Executor.Read` uses for every read refusal, including the
`ClassExecutionFailure` a command read mints. So the parameter threads
through both: `refusalOf` takes the error and `errors.As`-es it,
`refusalWithKeys` forwards it. Passing it to `refusalOf` alone would leave
reads — the one capability established tools bind *directly*, and so the
likeliest source of a stderr tail — silently tail-less. With both, every
present and future invocation site is covered, and refusals raised without an
error — timeout from `ctx.Err()`, read-back mismatch, the C6 gate — get an
empty `Detail` by construction rather than by each call site remembering to
leave it blank. The one refusal built outside it
(`executor.go::selects`, for an unknown accessor or capability mismatch) is
pre-selection: no binding has been chosen, so no invocation error can exist to
carry, and its `Detail` is empty for the same structural reason.
Path-backed bindings return untyped errors and are unaffected, so every
production `NewExecutor` caller and every existing binding behave unchanged —
the "0004 executor unchanged" claim survives as *classes* unchanged, with
**two** constructors gaining one parameter (`refusalOf` and the
`refusalWithKeys` that wraps it). **This is A11, still `Pending`**: the two
constructors are named from `main`, but that every call site can supply the
offending error or `nil` is the claim A11 verifies. If it fails, A11's "If
wrong" applies — the tail needs an interface change on the three binding
interfaces, and the "0004 executor unchanged" framing goes with it.

Rendering is a second hop and a **name collision to avoid**:
`internal/cli/flow_exec.go::accessorFailure` builds the CLI error and already
sets `CLIError.Detail` from `detailMayHaveApplied` (the 0004 applied-but-
unverified sense). There is exactly **one** slot to land in —
`internal/cli/clierr/clierr.go::CLIError` has a single `Detail string` and
`EmitText` renders a single `detail:` line — so "renders beneath" is not
available and no second carrier is added here. The field is documented
multi-line and `lint.go::runLint` already passes multi-line text through it,
so where both senses are present they occupy that one slot in order: the
applied-sense text **first** (losing "the write may have applied" to a
diagnostic tail would drop the more consequential fact), the stderr tail
appended after it. S3 asserts the tail is present on an `execution_failure`
that carries no applied sense, and that the applied sense leads when both are. With that field, the declared timeout
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
command_env_conflict          # an `env` key or `env_pass` name matching the reserved `INTRASTATE_` prefix (C4 overlay)

registration: all six are appended to `table.Categories()`, whose hand-maintained list IS the closed set; appended in the order declared above, after the existing members, so a consumer enumerating the list sees additions only at the tail. Each gets a typed `Category` constant of the existing `Cat…` form beside the others; the wire STRINGS above are the contract and the constant identifiers are not. The list's total size is not a contract at any point — it is append-only, and no clause or test may assert a count
precedence: WITHIN one entry — load is fail-fast (`0002:C24` — one categorized error for the whole document, never a list), so an entry carrying several of these defects reports the first in the order declared above. ACROSS entries in the SAME table there is no order: `accessorTable` ranges a Go map, so which of two defective entries is reported is unspecified, and no test may assert it. ACROSS tables the order IS fixed — `load.go::loadAccessors` runs read, then write, then gate, returning on the first error — so a defective read entry always masks a defective gate entry; that ordering is `0002`'s, inherited not established here, and a test may rely on it only as 0002's contract
interpreter set: OPEN (deny-listed, not closed) — an unlisted interpreter is admitted, so the list grows by amendment
```

Static validation. **Six** new load-time defect classes, surfaced by the same
validation that rejects malformed accessor entries today and reported by
`intrastate lint`. (A6's "four defects" counts the classes fixed at Resolve;
`command_output_shape` came with C3's raw mode and exit maps, and
`command_env_conflict` with C4's env composition. Six is the count this
contract registers, and S1 asserts six.) Home fixed at Resolve (A6): these
are
`internal/table/category.go::Category` constants with the spellings above —
the family `load.go::accessorTable` already uses for accessor-entry defects —
not `accessor.ValidationCode`, whose eight-member closure is RDR 0004's own
test contract. A `Category` constant is not in the closed set until it is
appended to `table.Categories()`, whose returned literal list *is* the set
(RDR 0024's "ONE decision") — an unregistered constant compiles and refuses
correctly at the call site while remaining invisible to every consumer that
enumerates categories, so S1 must assert membership, not merely the refusal.
Ordering is normative because `Categories()` returns a literal *slice* and
any golden output keyed on its order would otherwise move under an insertion:
the six append at the tail, in clause order.
Fail-fast precedence matters for the same reason S1 asserts per-mutant
categories — `accessorTable` reports one defect per entry, so a mutant
carrying two defects reports the earlier clause, and each S1 mutant carries
exactly one. That guarantee is **intra-entry only**, and deliberately so:
`load.go::accessorTable` ranges `map[string]sourceAcc`, whose iteration order
Go randomizes, and it runs once per capability table. Making cross-entry
reporting deterministic would mean sorting that map — a change to every
existing category's behaviour, not this RDR's scope. So a model with two
defective entries reports one of them, unspecified which, and S1's mutants
carry one defect in one entry each.
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
--allow-commands        # v1's ONLY opt-in: a per-invocation flag; never the model file; `lint` does not carry it
registration: ONE registration, on the `flow` GROUP — `internal/cli/flow.go::newFlowCmd`'s `PersistentFlags()`, whose subtree is exactly the four `buildRequest` call sites across three verbs (`flow_next.go`, `flow_resolve.go`, `flow_state.go` ×2) and nothing else. `lint` sits at ROOT, outside the group, so it does not carry the flag — which is this clause's first line, now structural rather than remembered. The repo already sets this precedent: `internal/cli/root.go` registers the shared output-mode flag on `PersistentFlags()` and reads it via `Lookup` + nil check, never a discarded `GetBool`
lookup: `buildRequest` reads the flag with the error CHECKED, not discarded. Because persistent registration guarantees presence, a lookup miss can only be a wiring bug and MUST surface as one — pflag's own text is "flag accessed but not defined", a programmer error, never a statement about user intent. A default-deny that rests on an unregistered lookup returning the zero value is NOT the gate: it would make the safety property a coincidence between two independent decisions (no verb registers it; this site discards the error), and the dangerous edit — a parent registering it, or a verb registering a different default — silently converts the miss into a hit and fails OPEN. Prior art (`evidence/reconcile/prior-art-absence-and-gating.md`): eleven mature Go CLIs were swept and NONE gates a dangerous operation on a lookup miss; etcd is the direct precedent (security booleans on root `PersistentFlags`, read by shared helpers that check the error and exit), and hugo's exec allowlist makes its zero value deny by explicit construction rather than by absence
absent ⇒ every command invocation refuses execution_failure before spawn, Detail naming the gate; lint (C1/C5) validates regardless
gate site: ONE — `buildRequest` reads the flag and passes it to `flowbind.Registry`, the single production construction site (`flow_exec.go:830`), which builds refusing command bindings when it is unset. `Registry` is a free function, not a method, so this is a signature change on it plus the `flowRequest` field; every executor built from that shared registry inherits the gate, so no execution path can bypass it by reaching the executor another way
signature: `Registry` today is `func Registry(m *table.Model) accessor.Registry`. It gains BOTH new inputs and no others, in this order: `func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry`. `baseDir` is A12's model-file directory (C2's argv0 resolution root), `allowCommands` this gate; the return type is unchanged and `Registry` does not gain an error return — a carrier-less or gate-refused entry yields a refusing binding (C1), not a construction failure
```

Execution gate. A model-declared command is code that runs when the model is
used, and the model travels with the repo — so execution requires an opt-in
**outside the model**. **v1 ships the flag alone** (A10): the repo has no
`internal/cli/config/`, no `intrastate.toml` reader, and no user-scope
configuration of any kind, so a config file means building a whole new
discovery/precedence/malformed-file surface — net-new scope this RDR charts to
a successor rather than absorbs. The flag is a complete gate: with no config
surface, `--allow-commands` is the only opt-in and absence refuses. The cost
is named honestly — per-invocation friction, not the one-time opt-in a config
file would give (Consequences); the successor that adds the file inherits an
already-gated seam and owes only precedence (a file that grants, a flag that
also grants; neither can revoke, so the two compose disjunctively) and a
fail-closed malformed-file rule.

The gate sites **in the command binding's constructor**, not in the executor:
`NewExecutor` has several production call sites and gating there would need
the same check written at each, every one able to drift. The registry builds
bindings once from the model — `flowbind.Registry` is a free function with
exactly ONE production call site (`flow_exec.go:830`, inside `buildRequest`),
whose result is the shared `flowRequest.registry` every executor is built
from — so a constructor that returns a refusing binding when the flag is
unset gates every present and future caller by construction. The flag
therefore reaches `flowbind.Registry` as a new parameter, which is a
signature change on a free function, not a method (A13).

Its registration follows the `flow` GROUP, whose subtree is exactly the four
`buildRequest` callers (`flow_next.go`, `flow_resolve.go`, `flow_state.go` ×2)
and nothing else, so ONE `PersistentFlags()` registration covers the caller set
without a hand-kept verb list. `intrastate lint` sits at root, outside the
group, and does **not** carry the flag — because validation is ungated and a
flag that changes nothing is a false affordance; that exclusion is now
structural rather than remembered. Registration was considered as the drift
surface and deliberately removed as one: an earlier form of this clause rested
on an unregistered verb's lookup returning false, which is imprecise (pflag
returns `(false, NotExistError)` — "flag accessed but not defined", a
programmer error) and would make fail-closed a coincidence between
non-registration and a discarded error rather than a property. The dangerous
edit under that shape fails OPEN: a parent registering the flag, or a verb
registering a different default, silently converts the miss into a hit. With
one registration on the group and the lookup error CHECKED at `buildRequest`,
a verb added outside the group loses command execution loudly and a wiring bug
surfaces immediately — which is why the coupling is stated to the group that
contains `buildRequest`'s callers rather than to a per-verb list.
With the gate off, a command invocation refuses
before spawn (`execution_failure`, `Detail` naming
`allow_commands`), while `intrastate lint` still validates the entries — the
model is reviewable before it is trusted. This is the shape every shipped
peer converged on after gating too late: Consul script checks flipped
off-by-default and re-gated twice more, Hugo's deny-by-default
`security.exec.allow`, Go's `GOVCS` allowlist, beads refusing repo-persisted
commands pending a trust gate (reversal ledger).

**C7** — the implemented seam, quoted as it exists.

```normative
the three command bindings implement `internal/accessor/binding.go`'s interfaces UNCHANGED — this RDR adds no method, changes no signature, and every clause above is expressed within them:

  Read(ctx context.Context, art Artifact, requested []string) (values []KeyValue, unreadable []string, err error)
  Gate(ctx context.Context, art Artifact) (Verdict, string, error)
  Apply(ctx context.Context, art Artifact, planned []resolve.Tag) error
  Invocations() int          # on WriteBinding, beside Apply

the command write binding counts `Apply` ENTRIES — one per call — and performs NO internal respawn: one `Apply`, one spawn, one increment. `Invocations()` has zero production consumers (declared in `binding.go`, implemented once at `internal/cli/flowbind/flowbind.go::Writer.Invocations`); every use is a `== 0` / `== 1` test assertion that the accessor layer performed no retry, undo, or re-derivation (`0004:C14`). Those assertions are about the LAYER's re-entry, not process count, so an internally-retrying binding would both break them and misreport against C14

the write method is `Apply`, NOT `Write`. `Gate` returns verdict + reason + error, so C3's `reason` field crosses back through the Go return and is not dropped. C3's read envelope maps onto `Read`'s SPLIT return: a parsed key becomes a `KeyValue` in `values`, an omitted declared key a name in `unreadable` — UNREADABLE is that second slice, never an error and never an absent value. `exit_absent` absence crosses as `KeyValue{Key: k, Absent: true}` IN `values` — a present record carrying the flag, the same channel `internal/cli/flowbind/flowbind.go::Reader.Read` already uses for the path-backed binding. Returning the key in NEITHER slice does NOT establish absence: `internal/accessor/executor.go::readOutcome.classify` defaults a requested key found in neither to UNREADABLE ("guessing absence is the failure this contract exists to prevent", `0004:C6`), so an omitted key refuses the whole read `incomplete_read`. The two boundaries are distinct and both hold: at the BINDING return absence is the flag; at the RESOLVER seam it becomes omission, converted by `internal/accessor/model.go::ReadResult.OwnedSnapshot`, which is what `0004:C8` governs. A command binding that signalled absence by omission would be refused, not believed
```

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
| User-scope configuration discovery (`allow_commands`) | **none — no config subsystem exists** (no `internal/cli/config/`, no `intrastate.toml` reader) | Not built | C6 — v1 gates on `--allow-commands` alone (A10); the config surface is charted to a successor |
| Free-text diagnosis on a refusal | `accessor.Refusal` — no `Detail` field today | Extend | C4 — additive `Detail string`; class set unchanged |
| Typed carrier for the stderr tail | **none — `accessor.ExecError` does not exist**; no production code imports `os/exec` | Not built | C4 specifies it: `struct { Detail string; Err error }` with `Error()`/`Unwrap()`; A11 |
| Invocation error reaching the read refusal | `accessor.readOutcome` — drops the binding error today | Extend | C4/A11 — additive `err` field; `refusalOf`/`refusalWithKeys` gain the parameter |

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
| Binds established tools without wrappers | 2 — artifact path cannot reach argv; `0004:C3` bars ambient discovery as the workaround | 4 — `{artifact}` covers the positional case: reads and gates bind directly, writes still need a wrapper (A3) | 5 — anything spliceable | 2 — only shipped adapters |
| Prior-art alignment | 4 — Terraform external (PA-5) | 5 — Docker exec form (PA-3) + git difftool closed vocabulary (PA-4) + PA-5 envelope | 2 — awf-cli needed an injection linter after the fact (PA-1) | 3 — plugin-registry precedent, none in the peer family |
| Injection / blast radius | 5 — nothing varies | 4 — one path-valued hole, whole-element | 1 — caller data spliced into text | 5 — no author-declared argv |
| Integration cost per new tool | 2 — a protocol wrapper per tool | 4 — direct for path-positional tools, thin wrapper otherwise | 5 — none | 1 — an intrastate release per tool |
| Reversibility | 4 — can add placeholders later | 4 — vocabulary can grow or shrink by validation | 2 — authors embed arbitrary templates immediately | 2 — adapter API becomes public surface |
| **Total** | **22** | **26** | **16** | **18** |

O2 wins on the deciding rows: it is the only option scoring ≥4 on both static
reviewability and binding established tools — the two clauses the Problem
Statement conjoins. The second row is scored on the *read and gate* half,
where O2 binds directly; established-tool **writes** need a wrapper under both
O1 and O2, so that half separates them not at all (A3, Approach) and O2's
margin comes from the positional read/gate case plus integration cost. O1 concedes the product thesis (every tool needs a
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
- Negative: command entries require an **absolute** `--artifact` path (C2),
  while path-backed entries accept the relative paths every fixture and doc
  example in this repo uses today (`flow_input.go::parseArtifacts` stores the
  value verbatim). So converting an existing entry from `path` to `command`
  can break callers that lint green, and the two carriers differ in what they
  accept from the same flag. Accepted rather than absolutized, because
  intrastate's cwd is not the tool's frame of reference (C2) — but it is a
  migration cost, not merely a per-entry one.
- Negative: nothing executes without `--allow-commands` on the invocation
  (C6) — and in v1 that is **per invocation**, not the one-time opt-in a
  user-scope config would give, because no config surface exists to hold it
  (A10). Accepted as the price every gated-too-late peer ended up paying
  anyway, with the standing cost named: a scripted caller repeats the flag,
  and the successor that adds the config file makes a one-time opt-in
  available again (it composes disjunctively — it does not retire the flag).

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
- **Risk**: the per-invocation flag is the one gate shape none of the peers in
  this RDR's own reversal ledger *kept* — each ended at a persistent surface
  (config key, env var, allowlist file). A caller that invokes intrastate from
  CI or a Makefile repeats the flag at every site, and the cheap workaround is
  an alias or wrapper that always appends it, which defeats the gate silently
  rather than loudly.
  **Mitigation**: none available in v1 — A10 establishes there is no config
  surface to hold a persistent opt-in, and building one is the successor's
  scope (charted). What v1 buys is that the seam ships *gated*: the successor
  adds a second grant path to an already-refusing default, which is the cheap
  direction. Expect C6 to be amended, not rewritten, when that lands; the
  failure this ordering prevents is shipping ungated and re-gating later,
  which is what the ledger's peers each paid for.
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
- Silent risk: two readers on one role with **disjoint key sets** load clean
  (`load.go::checkAccessorBindings` counts readers per tag key, not per role)
  and `model.go::readerFor` selects by role, first match — so a migration
  shape that keeps a path-backed reader beside a new command-backed one
  verifies the write through whichever sorts first, with no ambiguity error.
  Diagnosis: the refusal names no reader, so the tuple looks correct; compare
  the role's reader declarations. This is inherited first-match behaviour,
  not something this RDR introduces (A4) — 0016's fail-closed `readerFor`
  closes it, and S5b arm (b) pins today's selection so that landing is
  visible rather than silent.
- Visible: a relative `--artifact` path that works for a path-backed entry
  refuses `execution_failure` for a command entry (C2's absolute-path
  guard) — `flow_input.go::parseArtifacts` stores the value verbatim and
  every fixture and doc example in this repo is relative today, so
  **converting an existing entry to `command` breaks invocations that lint
  green**. Diagnosis: the refusal `Detail` says the artifact path must be
  absolute for a command entry. Recovery is the caller's, not the model's:
  pass an absolute `--artifact`.
- Silent risk — **data loss, and the sharpest edge in this design**: a write
  tool that sinks its stdin into the artifact ingests the C3 envelope as
  content. The A3 spike's `tee {artifact}` overwrote the artifact with
  `{"state.phase":"review"}` and exit 0; it surfaces only at read-back
  (`git config` then fails `bad config line 1` →
  `read_back_incomplete`), never at the write, and the blast radius is
  whatever else reads that artifact. Recovery is the user's backup — this
  RDR's bindings hold no undo (`0004:C14`). **It is not statically
  detectable**: whether a tool reads stdin is not visible in its argv, so no
  C5 arm can catch it and none is claimed; the deny-listed interpreter check
  (C5) does not reach it either. What bounds it in v1 is the wrapper contract
  (envelope in, tool-native invocation out) — a convention, stated here as a
  convention and not as a guarantee — plus read-back, which turns silent
  corruption into a loud `read_back_incomplete` on the *next* invocation
  rather than at the write. Charted, not built: a declared `stdin =
  "none" | "envelope"` field would let the binding withhold the envelope from
  a tool that never asked for it, making the footgun unauthorable rather than
  discouraged; it is a carrier-shape change and so a successor's
  (premortem P-1/P-2, `evidence/critique/Charted:`).
- Visible: a command entry invoked on Windows refuses `execution_failure`
  naming the unsupported platform, before spawn (C4 `platform`). v1 supports
  the Unix platforms the process-group mechanism requires (A2, A1/C4), and
  Windows re-quotes the vector into one command line — which would void C2's
  authority bound silently. Refusing is what keeps that bound honest rather
  than platform-dependent; lint stays platform-neutral, so a model
  authored on Windows still validates.
- Recovery: fix the model entry, re-run `intrastate lint`, re-invoke;
  bindings hold no state between invocations (`0004:C14` — no retry, no
  undo).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A10 verified; A11/A12 opened by
      the 3amigo pass and re-scoped by the critique pass, A13/A14 opened by
      the critique pass — Stage 6 closes all four, and all four are
      source-search reads of `main`, not spikes)
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
2. `intrastate lint` accepts it; six mutated copies (path+command conflict,
   empty element, unknown placeholder, `["sh", "-c", …]` interpreter form,
   `output = "raw"` with two keys, an `env` key shadowing `INTRASTATE_ROLE`)
   are each rejected with their C5 defect.
3. With `--allow-commands` on the invocation (C6; without it, the same
   invocation refuses `execution_failure` naming the gate, before spawn),
   drive a state change end to end through the existing write path
   (`internal/cli/flow_state.go`'s set-state verb, whose planned tags come
   from its `--write` flags): the declared write **command** — not intrastate
   — applies the edit to the artifact. Automatic resolve→apply wiring is not
   this RDR's scope and no Phase builds it (see Problem Statement); what this
   validates is that a declared command can *carry* the apply, which is the
   half that is missing today.
4. Read-back runs through the declared command reader and verifies the
   written value; the result reports the written tags.
5. A command that sleeps past its declared timeout refuses `timeout` and
   leaves no orphan process.

End state: a linted model applies and verifies a real state change through a
declared command, with the executed argv readable in the model. "Through a
declared command" is exact for the read (`git config --get` binds directly)
and means *the declared wrapper* for the write — the wrapper's argv is in the
model and linted, its body is script the model does not carry (Approach,
Consequences).

### Phase 1: Carrier

Extend `internal/table` (`sourceAcc`, `Accessor`, `accessorTable`) with
`command` and the C1/C5 load rules.

### Phase 2: Binding family

Implement command-backed `ReadBinding`/`GateBinding`/`WriteBinding` over
`exec.CommandContext` honoring C2–C4 (a sibling package to `flowbind`).

### Phase 3: Selection

Discriminate the constructor by carrier field at
`internal/cli/flowbind/registry.go::Registry` (a signature change on a free
function with one production call site), refusing the carrier-less residue
rather than binding `Path: ""` (C1). Register `--allow-commands` on the `flow`
group's `PersistentFlags()` and thread it (C6) from `buildRequest` — reading it
with the lookup error checked — in the SAME phase: selection is the step that
first makes a command reachable, so the gate must not lag it by even one
commit.

Phases 1–2 are independently revertible — the carrier is inert without this
phase, and the binding package is unreferenced until it. From Phase 3 the
seam is live, which is why the gate lands here.

### Phase 4: Surface and proof

Wire lint reporting for the C5 defects and land the MVV scenario plus the
timeout and read-back failure scenarios as tests.

## Validation

### Testing Strategy

Coverage goal: every C1–C6 clause has a test that fails when its rule is
dropped; the MVV scenario is the integration proof.

1. **Scenario**: load-time validation (C1/C5) — a valid command entry plus
   the six MVV mutants (path+command conflict, empty element, unknown or
   partial `{…}` token, `["sh", "-c", …]`, raw output with two keys, an `env`
   key named `INTRASTATE_ROLE`).
   **Expected**: the valid entry loads; each mutant is rejected with its own
   named C5 `table.Category` — asserted by category, never by "validation
   returned non-empty"; all six categories are members of
   `table.Categories()`, at the tail and in clause order (the registration is
   what puts them in the closed set — a per-mutant refusal assertion passes
   without it); `accessor.ValidationCodes` stays at eight (REQ-84). Plus the
   **open-deny-list negative**: an entry whose argv0 is an unlisted
   interpreter (`perl -e`) **loads**, asserting that C5's interpreter set is
   deny-listed and not closed — without it the suite reads as "inline shell is
   impossible", which is the guarantee C5 explicitly declines to make.
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
   and one exiting 128; a 1 MiB + 1 byte stdout; and an argv0 that does not
   resolve, which fails with no exit code at all.
   **Expected**, in the order listed: deny; `execution_failure`;
   `execution_failure`; deny; `execution_failure`; established-absent
   (asserted as `KeyValue{Absent: true}` in `values` for every declared key,
   NOT as omission — omission would refuse `incomplete_read`, C7);
   `execution_failure`; `execution_failure`; `execution_failure`. Two
   different oracles, and the split is normative because the fixture covers
   only one of them.
   The **empty-stdout** arms (cases 3–7 — unmapped non-zero exit, the
   `exit_verdicts` hit and miss, `exit_absent = [1]`, exit 128) follow
   normative fixture **FX-exit-codes**
   (`evidence/spikes/a3-a7-a8-real-tools/output.txt` E1a–E1h, R6), every case
   of which is empty-stdout: E1b/E1e supply the mapped-verdict codes, E1c/E1f
   the error codes, R6/E1g the established-absent read, and **E1h the ninth
   case** — `exec.ErrNotFound`, no exit code at all, so no exit map can claim
   it however it is written (both C3 maps are keyed on a listed exit and there
   is none). That arm is what proves spawn failure cannot be laundered into a
   verdict by an over-broad map.
   The **non-empty-stdout** arms (cases 1, 2, 8 — deny envelope with a
   non-zero exit, malformed envelope with exit zero, 1 MiB + 1 byte stdout)
   have **no fixture oracle and need none**: their expected values are fixed
   by C3/C4's ordering rule directly (parse first, so a well-formed envelope
   wins over any exit; a malformed one is `execution_failure` whatever the
   exit; an over-cap stdout is `execution_failure` from the bound and is
   never exit-mapped, being non-empty). These are contract-derived assertions,
   not fixture-derived — asserting them against FX-exit-codes would be reading
   an oracle that says nothing about them.
   Each refusal carries the 0004 diagnosis tuple and, per C4, the 4 KiB stderr
   tail in `Refusal.Detail` — asserted non-empty on case 3, where no
   applied-sense text competes for `CLIError.Detail`.
4. **Scenario**: deadline (C4, A1) — a child sleeping past `timeout`, a
   wrapper whose grandchild holds the stdout pipe, and a child that never
   reads stdin (1 MiB stdin); plus the ablations: no process group, no
   `WaitDelay`.
   **Expected**: `timeout` refusal within `timeout + WaitDelay` in all three,
   no hang, and no process from the child's group surviving the refusal —
   normative fixture **FX-deadline** (`evidence/spikes/a1-deadline/output.txt`
   S1–S5: S2/S4/S5 ≈1.0s with 0 survivors). The two ablations are **not test
   arms**: C4 states the `Setpgid`/`Cancel`/`WaitDelay` triple as fixed
   behaviour with no injection seam, and adding one purely to disable it would
   put a way to weaken the deadline into the shipping binding. They are cited
   from the spike (S3's orphan, S1's blocked `Wait`) as the evidence that each
   leg is load-bearing; the suite asserts the assembled triple only.
4b. **Scenario**: raw read and env policy (C3/C4, A3/A7) — `git config --file
   {artifact} --get <key>` in raw mode; the same entry with an inherited
   `GIT_CONFIG_COUNT` injection in the parent env and no `env` declaration.
   **Expected**: the value is stdout minus one trailing newline — normative
   fixture **FX-raw-read** (`output.txt` R1: stdout `"draft\n"` → `draft`);
   the injection never reaches the child (its env holds only the allowlist),
   so the read is unchanged. Plus the composition arms C4 now orders: a
   variable named in `env_pass` reaches the child; an unnamed one does not; a
   key set in both the parent allowlist and the entry's `env` arrives with the
   entry's value; an entry `env` key shadowing `INTRASTATE_ROLE` is a C5
   defect at load; `LC_ALL` passes and `LCFOO` does not (the `LC_` prefix is
   literal). Asserted on the child's observed environment, not on the read
   result, since an env defect can leave the value correct.
5. **Scenario**: write read-back (C3/C4) — a command write whose wrapper
   applies the planned tags; one whose tool exits zero without applying them
   (`git config` with the value on stdin, spike R3b); one whose argv0 sinks
   stdin into the artifact (`tee {artifact}`, spike R4b); one planning
   `<clear>`.
   **Expected**: success reporting the written tags; `read_back_mismatch`;
   `read_back_incomplete` (the corrupted artifact no longer parses for the
   reader); the key reads back absent — the write's exit status decides none
   of them. Plus the inherited reserved-literal arm (C3): a read whose tool
   emits the literal `<clear>` is UNREADABLE
   (`internal/accessor/executor.go::readOutcome.classify`), not a cleared key
   — newly reachable through raw-mode reads of real tools, and asserted by
   class, since a cleared key and an unreadable one are indistinguishable by
   absence alone. Asserted on the **read** path, where `clearIsUnreadable` is
   true; the write's read-back passes false, where the same literal is a
   mismatch instead — the two senses are asserted separately and neither test
   stands in for the other.
5b. **Scenario**: read-back reader selection under two readers on one role
   (A4, Prerequisites). The arity rules are keyed differently and that
   difference is the scenario: `load.go::checkAccessorBindings` counts
   readers **per tag key** (an owned tag needs exactly one), while
   `model.go::readerFor` selects **per role**, first match. So on a command
   write's own owned key the reader is loader-forced to be unique and
   first-match is unobservable — the fixture must NOT be built there, or it
   asserts nothing. The authorable case is two readers on one role with
   **disjoint key sets** (the live migration shape: a path-backed reader kept
   for rollback beside a new command-backed one), which loads clean and where
   `readerFor` silently takes the first.
   **Expected**: two arms. (a) The owned-key arm: a second reader serving a
   command write's owned key is REFUSED at load
   (`malformed_accessor_binding`) — so read-back's reader is unique by
   construction, not by selection order. (b) The disjoint-keys arm: pending
   0016's fail-closed uniqueness, the selection is `registry.go`'s
   name-sorted first match — asserted **by the selected reader's identity**,
   never by "read-back succeeded", since either reader can verify a correct
   write and the test would pass under either. The property is emergent
   across two packages — `internal/cli/flowbind/registry.go` sorts the names,
   `internal/accessor/model.go::readerFor` takes the first match by slice
   order — so arm (b) asserts **both halves**: that the registry emits
   readers name-sorted, and that `readerFor` returns the first. A test on the
   composed outcome alone passes while either half moves. Arm (b) is what
   makes 0016's landing a visible change rather than a silent one; it is
   also, until 0016 lands, a real hazard the model author can author — named
   in Failure Modes.
6. **Scenario**: the MVV end to end against an established tool present in
   CI.
   **Expected**: the declared write command, not intrastate, edits the
   artifact, and the declared command reader verifies it.
6b. **Scenario**: platform refusal (C4 `platform`, F7) — a valid command
   entry invoked with the binding's `goos` set to `windows`.
   **Expected**: `execution_failure` naming the unsupported platform, with no
   child spawned (asserted by absence of a spawn, as in scenario 7); the same
   model passes `intrastate lint` under that setting, since lint is
   platform-neutral. CI runs `ubuntu-latest` on all six jobs, so the injected
   `goos` is the only way this branch is ever exercised — without it the
   clause ships untested.

7. **Scenario**: execution gate (C6) — a valid command model invoked without
   `--allow-commands`; with it; `intrastate lint` under both. The executor
   callers are the wrong axis: `flowbind.Registry` has ONE production call
   site (`flow_exec.go:830`), so every `NewExecutor` shares one gated
   registry and testing them separately tests one path repeatedly. With C6's
   ONE registration on the `flow` group, the drift surface is **containment**,
   not per-verb registration: the scenario derives the set of verbs reaching
   `buildRequest` and asserts each RESOLVES the flag through the group's
   persistent set — a structural `Flags().Lookup("allow-commands") != nil`
   walk, never a text match on pflag's error. A fifth verb added INSIDE the
   group inherits the flag and must pass; one added OUTSIDE it must fail this
   test rather than silently inherit a lookup miss. The complementary arm
   asserts `lint`, at root, does NOT resolve the flag — C6's first line, now
   testable. The oracle derives both sets rather than restating a literal
   four.
   Plus the C1 residue arm: an entry with neither carrier builds a
   **refusing** binding, never a `Path: ""` file binding.
   **Expected**: refusal `execution_failure` naming `allow_commands` with no
   child process spawned — asserted by **absence of a spawn** (a sentinel
   argv0 that would leave an observable trace if executed), not merely a
   non-zero exit; normal execution under the flag; lint passes in both
   (validation is ungated) and does not resolve the flag; every in-group verb
   resolves it; a verb outside the group fails the containment assertion, and
   the carrier-less entry refuses rather than reading an empty artifact. v1
   has no config file to malform (A10).

### Pre-Lock Mini-Checks

Cue-fired at Stage 5. Four of five fired; **test-discriminability did not**
(every Expected line above names a discriminating value, and S1 already bans
the absence-of-error oracle).

**`authority`** — source-of-truth census (cue: ≥2 candidates; a resolver; sibling arms).

| Input / decision | Canonical writer | Readers | Sibling arm | Which is canonical |
| --- | --- | --- | --- | --- |
| Load-time defect codes | `table.Category` + `Categories()` | `intrastate lint`, loader | `accessor.ValidationCode` (8, closed) | `table.Category` — A6; the 0004 set is 0004's own test contract and does not grow |
| Child environment | C4 allowlist | executor at spawn | `env_pass`, entry `env`, `INTRASTATE_*` overlay | C4's four-part composition is total and ordered (later layer wins: overlay > entry `env` > `env_pass` > allowlist); nothing else is inherited, and shadowing the overlay is a C5 defect |
| Execution permission | `--allow-commands` flag (v1's only opt-in) | command-binding constructor, pre-spawn | user-scope config (charted, not built) | The flag — A10; absence refuses. A successor's config composes disjunctively |
| Read result | child stdout envelope (C3) | command binding | `flowbind.go::store` (artifact file) | Disjoint transports, one shared map-of-strings rule (Wire LBD). The carrier selects the transport and the selection must be TOTAL: `flowbind.go::load` treats a missing file as an empty store, so a command entry reaching the file binding would read absent instead of refusing — C1 refuses the residue at the constructor (A14) |
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
| empty stdout, exit in `exit_absent` | every declared key established absent — `KeyValue{Absent: true}` in `values` (C7) | loud |
| empty stdout, unlisted exit (zero or not) | `execution_failure` — read and gate alike (C3) | loud |
| stdout > 1 MiB | `execution_failure` | loud |
| deadline exceeded | `timeout` (from `ctx.Err()`), group killed | loud |
| spawn failure / gate off / Windows | `execution_failure`, `Detail` names the cause | loud |
| write applied, read-back unavailable | `read_back_incomplete`, `Applied()` true | loud — never silent success |

**`trace`** — desk trace of the MVV end state (cue: C2/C3/C4 all bear on the invocation result).

| MVV step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1 author model | C1 exactly-one; C5 six defects registered | `read.branch_state` with `command`, no `path` | consistent |
| 2 lint accepts, 6 mutants rejected | C5 + `Categories()` membership (S1) | each mutant → its own category | consistent |
| 3 gate off → refuse | C6; A10 (flag alone is v1's gate) | `execution_failure`, no spawn | consistent |
| 3 gate on → invoke write | C2 absolute-path substitution; C4 env allowlist | argv = declared vector + absolute `{artifact}` | consistent |
| 4 read-back through declared reader | C3 raw mode; C4 write rule; A4 `readerFor` | `git config --get` → `draft\n` → `draft` (FX-raw-read R1) | consistent — **single reader, loader-forced**: `checkAccessorBindings` counts readers per tag key, so the write's owned key admits exactly one and selection order cannot apply here. The two-reader case needs disjoint keys (S5b arm b), unpinned until 0016 |
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

Responses: 0025-command-invoking-accessor-bindings/artifacts/gate.md (Gate PASS 2026-08-30)

### Cross-Cutting Concerns

[Gate key: cross-cutting]

**Secret / credential lifecycle.** Declared `command` argv is reviewable data
by design, so anything placed in it is visible in the model file, in `git`
history, and in host process listings. This RDR's policy: **credentials never
appear in `command`, in `env` literals, or in stdin values** — a bound tool
reads its own credential from its own configuration, and where a variable must
reach the child it is named in `env_pass` (C4), which forwards a value without
recording it. The C4 allowlist is what makes this enforceable rather than
advisory: the child environment is composed, not inherited (A7 establishes that
inherited env silently steers the same binary and argv), so a credential cannot
arrive by accident. A lint advisory on credential-shaped argv elements is
**not** in v1 and is charted to the successor that adds the config surface
(A10), alongside the persistent `allow_commands` opt-in.

**Deployment model.** Command-backed accessors require a POSIX process model:
`Setpgid` plus group `SIGKILL` plus `WaitDelay` is the necessary-and-sufficient
termination triple (A1), and it has no Windows equivalent in this design. A
command entry invoked on Windows refuses `execution_failure` with `Detail`
naming the cause (F9) — a loud, declared refusal, never a silent unbounded
child. Path-backed accessors are unaffected, so Windows keeps the whole
path-backed surface. Widening this is a successor's scope.

**Incremental adoption.** The carrier is additive and the selection is total:
`path` → today's file binding, `command` → the new binding, exactly one
declared (C1). Every existing model keeps working untouched, and no path-backed
entry changes meaning. The one migration cost is real and named rather than
hidden: converting an entry from `path` to `command` requires an **absolute**
`--artifact` path (C2), while path-backed entries accept the relative paths this
repo's fixtures and docs use today, so a converted entry can break a caller that
lints green (Consequences, F7). Adoption is also gated at the invocation
(`--allow-commands`, C6), so the feature is opt-in twice over — per model entry
and per invocation.

**Concurrency model.** This RDR adds a child-process boundary, not concurrency:
each invocation spawns one child in its own process group and waits under the
declared timeout. The concurrency-adjacent obligation it does own is **draining
without deadlock** — stdout is bounded at 1 MiB and the drain is bounded by
`WaitDelay`, so a child that holds a pipe open cannot hang the parent past the
deadline (C4, A1). Ordering across accessors is unchanged and stays RDR 0004's.

**Character encoding.** Two encodings meet at the child boundary and the RDR
fixes both. The `json` mode is a flat JSON object of strings (UTF-8 by the JSON
spec), where an omitted key means not-carried and never empty-string. The `raw`
mode is defined bytewise: the value is the child's stdout minus **exactly one**
trailing `\n` — no trimming, no case folding, no whitespace normalization, so a
value with meaningful leading or internal whitespace survives. intrastate never
re-emits the child's bytes (the artifact's byte form belongs to the bound tool,
not to intrastate — Performance Expectations), which is why no canonical-form
obligation attaches.

**Canonical-form / determinism.** This RDR claims **no** byte-identical output,
content-addressed identity, or replay-stable hash, so the hash/pre-image
checklist does not apply. What it does claim is a fidelity invariant of a
different kind, and that one is pinned: declared argv → executed argv is
byte-for-byte equality except the single whole-element `{artifact}`
substitution, and planned tags → stdin → tool → read-back is **typed tag-value**
equality, explicitly not byte equality, with `<clear>` round-tripping as absence
(Pre-Lock Mini-Checks, `fidelity`). Read-back comparison-set semantics are not
this RDR's to fix — **RDR 0016 owns them** (`cli/0016:C4`, reader-per-role
resolution and comparison sets; kata intrastate#p63c), and A4 declares this RDR
a consumer of that policy rather than a second author of it.

**Versioning.** The placeholder vocabulary is versioned as a closed set: v1 is
`{artifact}` and is complete (C2). Because it is closed and validated at load,
an unknown placeholder is a C5 defect rather than a silently-passed literal —
which is what lets the vocabulary grow later without ambiguity about what an
older model meant. C6's per-invocation gate is expected to be **amended, not
rewritten**, when a config surface lands: the successor adds a second grant path
disjunctively to an already-refusing default (Risks).

Not applicable, and omitted rather than N/A-bulleted: build tool compatibility,
licensing, IDE compatibility, memory management.

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
