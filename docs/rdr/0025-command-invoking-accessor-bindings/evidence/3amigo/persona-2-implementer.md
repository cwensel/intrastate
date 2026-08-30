Model: claude-opus-5[1m]

# Persona 2 — Implementer

Question answered: *if I started coding this Monday, what would I ask in the first hour?*

Owned starting set: `C1`–`C6`, `D-identity` / `D-wire-byte-format` / `D-naming` /
`D-selection-predicate`, and the `source-anchor` edges. **Widened** to
`§critical-assumptions` (A4, A9, A10), `§implementation-plan`, `§testing-strategy`,
`§capability-dependencies`, and `§failure-modes` — because five of the eight findings
below are *silences* in C1/C4/C6 (an unstated TOML surface, an unstated carrier for a
new field, an unstated gate injection point, an unstated absolutizer) and a silence has
no line range. What sent me there is named per finding.

Code grounded in `/Users/cwensel/sandbox/newcoinc/intrastate` at HEAD `b614406`.

---

## High

### H1 — `0025:C3` / `0025:C1`: five new TOML keys are normative but never declared on the carrier, and the loader decodes strictly

`C3` and `C4` make `output`, `exit_absent`, `exit_verdicts`, `env`, and `env_pass`
normative per-entry keys, and `C5`'s `command_output_shape` defect is *defined in terms
of* `output`/`exit_absent`/`exit_verdicts`. But `C1`'s normative block declares exactly
one new carrier field:

> ```normative
> [read.<id> | gate.<id> | write.<id>]
> command = ["<argv0>", "<arg>", ...]   # TOML array of strings
> ```

and Phase 1 (`0025:§implementation-plan`) says only "Extend `internal/table`
(`sourceAcc`, `Accessor`, `accessorTable`) with `command`".

The loader decodes strictly. `internal/table/source.go::decodeStrict` calls
`dec.DisallowUnknownFields()` and maps a strict miss to `CatUnknownSchemaField`;
`internal/table/source.go::sourceAcc` today carries only
`Role/Path/Keys/Timeout/ReadBack`. So on Monday, an entry carrying `output = "raw"`
fails the *whole document* with `unknown_schema_field` before `accessorTable` ever
runs — which is not the C5 category the Testing Strategy S1 asserts. The Go types are
also unstated: is `exit_verdicts` `map[string]string` keyed by a decimal string (as the
Illustrative Code's `{ "0" = "allow" }` implies) or `map[int]Verdict`? Is `output`'s
default the literal string `"json"` or the empty string treated as json? Is
`output = "json"` explicitly writable or only implicit?

**Decision blocked:** the `sourceAcc` / `Accessor` field list and their TOML types —
i.e. all of Phase 1, and with it the C5 category-vs-`unknown_schema_field` split that
S1 asserts on.

### H2 — `0025:C4`: `Detail` is specified on `accessor.Refusal`, but no plumbing exists to carry a stderr tail from a binding to a refusal, and the CLI drops it

C4 states:

> `bounds:   stdout capped at 1 MiB (overflow = execution_failure); a NEW `Detail string` field on `accessor.Refusal` carries the last 4 KiB of stderr`

The field addition is indeed additive (`internal/accessor/model.go::Refusal` carries
`Class, Accessor, Capability, Role, Timeout, Keys, Expected, Observed, Reason` plus
unexported `applied`). But there are two unstated links in the chain:

1. **The binding cannot deliver it.** `internal/accessor/binding.go` declares
   `Read(...) (values []KeyValue, unreadable []string, err error)`,
   `Gate(...) (Verdict, string, error)`, `Apply(...) error`. The executor *discards*
   the error: `internal/accessor/executor.go::invokeRead` does
   `if err != nil { return readOutcome{class: ClassExecutionFailure} }`, and
   `Executor.Gate` / `Executor.Write` likewise return
   `refusalOf(def, timeout, ClassExecutionFailure)` with no reference to `err`.
   Nothing in `refusalOf` takes a detail argument. So carrying the tail requires either
   a typed error the executor unwraps (`errors.As` on a new exported error type — which
   RDR 0004's package would then own), or a widened binding signature — a real
   0004-surface decision the RDR does not make.
2. **The CLI drops it.** `internal/cli/flow_exec.go::accessorFailure` renders
   `ClassExecutionFailure` as `envErr(codeAccessorFailed, id, "the accessor `"+id+"`
   could not be executed")` and never reads `refusal.Detail`. Worse, `clierr.CLIError`
   already has its own `Detail` field which that same function assigns
   `detailMayHaveApplied` on the timeout and read-back-incomplete arms — so "put the
   stderr tail in `CLIError.Detail`" collides with an existing, contract-bearing use.

**Decision blocked:** the mechanism by which a command binding's stderr reaches a
refusal (typed error vs. interface change), and whether/where the CLI renders it —
which S3 asserts on ("each refusal carrying … the 4 KiB stderr tail").

### H3 — `0025:C6` / `0025:A10`: the execution gate has no named injection point, and three executor construction sites exist

C6 requires that with the gate off, "a command invocation refuses before spawn
(`execution_failure`, `Detail` naming `allow_commands`)", and S7 asserts this **by
absence of a spawn**. But the RDR never says *where* the gate value enters the object
graph. The candidates are materially different builds:

- `internal/cli/flowbind/registry.go::Registry(m *table.Model)` takes only the model —
  gating here means changing the sole registry constructor's signature (it is called
  from three sites) or constructing a refusing stub binding.
- `internal/accessor/NewExecutor(reg Registry, arts Artifacts)` — gating here puts a
  CLI-policy boolean inside the RDR-0004 accessor package, which `0004:C3`'s
  no-ambient-authority posture arguably resists.
- The command binding itself, holding the flag as a constructed field.

There are three live `NewExecutor` call sites (`internal/cli/flow_exec.go:256`,
`:644`, `internal/cli/flow_state.go:306`), so whichever seam is chosen must be threaded
to all three or the gate leaks on one path. Compounding this, **A10 is the only
assumption still `Status: Pending`**, and it explicitly leaves open "whether v1 ships
the flag alone or also the file". Widened here from C6 to `§critical-assumptions` and
`§capability-dependencies` (which marks the config capability `Build`, source
"**none — no config subsystem exists**").

**Decision blocked:** where the gate is enforced and how it is threaded — i.e. whether
Phase 3 is a one-line discriminator or a signature change rippling through three call
sites; and whether Phase 4 builds a config reader at all.

---

## Medium

### M1 — `0025:C2`: argv0 must resolve against "the model file's directory", but `table.Model` does not carry one

C2 fixes argv0 resolution:

> a name containing a path separator resolves against the **model file's directory**
> (the git-hook / pre-commit convention)

The Illustrative Code relies on this (`command = ["tools/flowstate-write",
"{artifact}"]   # separator ⇒ resolved against the model file's dir`). But
`internal/table/model.go::Model` (lines 440–479) carries
`ID, Version, Description, Class, Metadata, Outcomes, Initial, Terminal, Tags,
Readers, Writers, Gates, EmitDecls, DumpOrder, Rows` — **no source path**, and
`internal/table/model.go::Accessor` carries `Role, Path, Keys, Timeout, ReadBack` —
also none. `internal/cli/flowbind/registry.go::Registry` receives only `*table.Model`
and so has no way to compute the base directory. The CLI does hold it
(`--model` flag, `internal/cli/lint.go:130`), but nothing carries it into the model or
the registry.

Two further under-specifications inside the same clause: is the model file's directory
resolved *at load* (baked into the stored `Accessor.Command[0]`, per "argv0 resolution
is fixed at load") or *at spawn*? And a multi-model document (`DumpAll`/multi-model
loads exist) has one file but many models — same answer either way, but the storage
site differs.

**Decision blocked:** whether `Model` (or `Accessor`) gains a source-path field, or
whether the registry constructor gains a base-dir parameter — a Phase 1-vs-Phase 3
scoping question that changes which struct the round-trip and dump tests must tolerate.

### M2 — `0025:C2`: the absolute-path precondition has no named absolutizer, and today's CLI passes `--artifact` paths through verbatim

C2's substitution contract:

> the substituted value must be an absolute path — the executor refuses the invocation
> (`execution_failure`) when the caller-bound artifact path is relative or begins with
> `-`

Grounded: `internal/cli/flow_input.go::parseArtifacts` splits `role=path` and stores
`out[role] = path` with **no** absolutization or validation beyond non-empty;
`internal/cli/flow_exec.go::artifactMap` copies it straight into
`accessor.Artifact{Role: role, Path: path}`. So `--artifact repo=./state.ini` — a form
that works today for path-backed entries — refuses for every command-backed entry. The
RDR does not say whether that is intended user-facing friction, or whether the CLI
should `filepath.Abs` the binding before it reaches the executor (which would make the
refusal unreachable except for `-`-prefixed values), or whether absolutization belongs
in the binding. The clause also says "the **executor** refuses", but the executor
(`internal/accessor/executor.go`) is binding-agnostic and would have to learn a
command-specific precondition to do so — placing the check there contradicts C4's claim
that the executor is inherited "with its refusal *classes* unchanged and one additive
type change".

**Decision blocked:** who absolutizes/validates the artifact path (CLI, executor, or
binding) — and therefore whether existing relative-path invocations break.

### M3 — `0025:C5`: the clause says "all five" register but names five codes whose count is inconsistent with the prose, and the registration target is a hand-maintained list with a closure test

C5's registration line reads:

> `registration: all five are appended to `table.Categories()`, whose hand-maintained
> list IS the closed set`

but the same section's prose says "Decision (C5): the **four** defects are …" in A6's
Evidence, and the block lists five constants (`command_and_path_conflict`,
`command_empty`, `command_unknown_placeholder`, `command_shell_interpreter`,
`command_output_shape`). Five is almost certainly right — the MVV and S1 both enumerate
five mutants — but the mismatch with `0025:A6` is a real coin-flip at the point of
writing the constant block. Grounded: `internal/table/category.go` does have exactly the
two-place shape C5 describes (a `const` block plus a `Categories()` returning a literal
slice), and the RDR 0024 comment there already states the "wire slug and the registered
category are ONE decision" rule, so the mechanism is real.

Unstated: whether the constant names follow the existing `Cat…` prefix convention
(`CatCommandAndPathConflict`?) and where in the declaration order they append —
`Categories()` returns "in declaration order" and there are order-sensitive golden
tests (`internal/table/dump_test.go:993` maps each category to a slug string).

**Decision blocked:** the constant block's contents and ordering — trivially reversible
but asserted on by S1 ("all five categories are members of `table.Categories()`") and
by existing golden maps.

### M4 — `0025:C1` + `0025:C5`: `command_and_path_conflict` overlaps an existing refusal that fires first, and `accessorTable`'s validation is a fail-fast `switch`

C1 relaxes the path-required rule. Grounded at
`internal/table/load.go::accessorTable` (lines 942–995): validation is a single
`switch` returning on the *first* failing arm, in order role → path → keys → timeout,
with `case a.Path == nil || *a.Path == "": return nil, bad("path is absent or empty")`
carrying `CatMalformedAccessorDeclaration`. Two implementer questions the record does
not answer:

1. The "neither declared" half of `command_and_path_conflict` is exactly the case that
   today mints `CatMalformedAccessorDeclaration` — and `internal/table/accessors_test.go`
   has existing arms (lines 43–90) asserting that category for accessor-shape defects.
   C5 says the neither-case is now `command_and_path_conflict`. That is a **behaviour
   change to an existing asserted category** for at least one existing negative fixture
   (`neg/neg-accessor-*`, referenced at `dump_test.go:837`), and the record does not
   name it as such or say which existing tests move.
2. Load is fail-fast and yields exactly one failure (`Failure`'s doc: "a document
   yields exactly one of these and never a list"), yet the five C5 defects can co-occur
   on one entry. The precedence among them is unstated, and S1 asserts each mutant
   "is rejected with its own named C5 `table.Category`" — which only holds if each
   mutant triggers exactly one, or if precedence is fixed.

**Decision blocked:** whether existing `CatMalformedAccessorDeclaration` fixtures
re-categorize, and the C5 defect precedence order inside `accessorTable`.

---

## Low

### L1 — `0025:C4`: the env allowlist's `LC_*` is the only glob in a policy that explicitly bans globs, and `HOME`'s presence is unreconciled with the A7 evidence

C4's env line: `child env = allowlisted parent vars {PATH, HOME, TMPDIR, LANG, LC_*}`.
`env_pass` is described as "the named, **no-glob** escape hatch", so `LC_*` is the one
prefix-match in the design — implementable, but the record does not say whether it is
prefix or shell-glob, nor whether `LC_ALL` alone or every `LC_` var passes.

More substantively: A7's spike (`0025:A7`) is cited as showing "the same binary and argv
answering from `~/.gitconfig`" — i.e. `HOME` is precisely the vector — yet `HOME` is on
the allowlist. C4's own reasoning ("the A7 spike showed the same binary and argv
answering from `~/.gitconfig`, `GIT_CONFIG_GLOBAL` or `GIT_CONFIG_COUNT` injection when
those were inherited") names `~/.gitconfig` as a defect but keeps its enabling variable.
The Illustrative Code compensates per-entry (`env = { GIT_CONFIG_NOSYSTEM = "1" }`),
which suggests `HOME` is deliberate — but an implementer cannot tell whether keeping
`HOME` is the decision or an oversight, and dropping it would break `git`, `ssh`, and
most wrappers.

**Decision blocked:** the literal allowlist contents and the `LC_*` match semantics —
low risk, but the value is a normative constant that S4b asserts on.

### L2 — `0025:C4` / `0025:§failure-modes`: the Windows refusal has no home and no defect class

Failure Modes states: "a command entry invoked on Windows refuses `execution_failure`
naming the unsupported platform, **before spawn**", while "lint stays platform-neutral".
Nothing in C1–C6 carries that rule — it exists only in prose in `§failure-modes` — and
no `table.Category` or refusal-detail spelling is fixed for it. Since it is a runtime
(not load-time) refusal it cannot live in C5, and since C6's gate refusal is the only
other before-spawn refusal specified, the implementer must invent the check site and
its message. `internal/accessor/model.go`'s `RefusalClass` set is explicitly closed
(`refusalClasses` / `RefusalClasses()`), so `execution_failure` is right — but the
build-tag strategy (`//go:build unix` on the binding? a runtime `runtime.GOOS` check?)
is unstated, and the two differ in whether the package compiles on Windows at all.

**Decision blocked:** build-constraint vs. runtime check for platform support, which
determines whether `Setpgid`/`Cancel` code can be written unguarded.

---

## Non-findings (checked, held up)

- `0025:A1`'s executor-owns-the-deadline claim: confirmed —
  `internal/accessor/executor.go::invokeRead` wraps in `context.WithTimeout(ctx, timeout)`
  and classifies from `bounded.Err()`, exactly as C4 requires. A binding carrying no
  timer of its own is implementable as written.
- `0025:A4`'s `readerFor` claim: confirmed — `internal/accessor/model.go:227–231`,
  `func (reg Registry) readerFor(role string) (Definition, bool)`, selects by
  `Capability == CapRead && Accessor.Role == role`, binding-agnostic. A command-backed
  reader is admissible with no change. S5b's "pin the inherited first-match" instruction
  is concrete enough to build.
- `0025:A9`'s locator sweep: spot-confirmed — `internal/table/dump.go` is row-level and
  emits no accessor `Path`; `Registry` is the only non-test constructor consuming
  `acc.Path`; the file binding's own `Path` reads
  (`flowbind.go::verdictFor`, `::unreachable`) are on the path-constructed branch only.
- `0025:D-selection-predicate`'s "no precedence order exists to get wrong": correct given
  C1's exactly-one rule — the discriminator at `registry.go::Registry` is a genuine
  two-way switch with no third state.
