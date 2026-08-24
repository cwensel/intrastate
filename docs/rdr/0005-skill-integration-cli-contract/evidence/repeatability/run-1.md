model: claude-sonnet-5
variant: lite (profile: mid)

# RDR 0005 — Skill Integration CLI Contract: API Reconstruction (run 1)

Reconstructed solely from
`docs/rdr/0005-skill-integration-cli-contract.md`. No project source, peer
RDRs, or JDR 0001 text were read; every place the RDR cites but does not
restate a peer/JDR clause is filled with a best guess, marked **GUESS**.

## 1. Public API

The RDR specifies a Cobra command-group surface (`flow`), not a Go package
API. The "public API" is the CLI verb contract: four subcommands, their
flags, and their JSON `data` response shapes. Go-level function signatures
for the underlying handlers are not given by the RDR (implementation detail
left to Phase 1–3), so those are marked GUESS where included.

### Command group: `flow`

Global mode flag (from the existing, cited-not-restated output contract):
`--as text|json` (persistent flag; default presumed `text` — **GUESS**,
RDR only says the flag "selects the wire format").

#### `flow next`

```
intrastate flow next
  (--flow <id> | --model <path>)
  [--tag name=value ...]
  [--artifact role=path ...]
  [--evaluate-gates]
  [--as text|json]
```

- Purpose: "what outcomes are legal from this state?"
- Runs declared read accessors over `--artifact` bindings to assemble
  `owned` state; `--tag` supplies `observed` context only.
- Gates run only under `--evaluate-gates`; otherwise gate ids are listed
  as unresolved facts.
- Never calls a language model. Effect-free by default (reads only,
  opt-in gates).
- Success `data`: `model`, `revision`, `observed`, `owned`, `readers[]`,
  `outcomes[]` (recognized outcome tags only), `candidates[]` (each with
  source rule identity, outcome, required facts, unresolved guard/gate
  facts, evaluated gate results when `--evaluate-gates`, preview next
  tags/write targets/clear keys when derivable without evaluating missing
  facts).
- Error modes: `flow-tag-*` family, `flow-artifact-*` family,
  `flow-model-not-found`, `flow-model-invalid`, accessor exit-3 codes if
  a read accessor is invoked and fails.

#### `flow resolve`

```
intrastate flow resolve
  (--flow <id> | --model <path>)
  [--tag name=value ...]
  [--artifact role=path ...]
  --outcome <tag>
  [--as text|json]
```

- Purpose: "given this recognized outcome and this state, what exact
  plan or refusal results?"
- Assembles `owned` from read accessors (same as `next`), delegates
  selection to the kernel, then runs the selected row's gates
  (all gates on the row; deny overrides; indeterminate does not override
  deny) after exact-one selection and before emitting a plan.
- Returns exactly one plan or exactly one `CLIError` refusal.
- Success `data`: `model`, `revision`, `observed`, `owned`, `readers[]`,
  `outcome`, `rule`, `gates[]` (`id`, `result`), `next` (next tag-set),
  `writes` (planned owned-tag writes, set values as canonical JSON
  arrays), `clear[]`, `escaped` (bool), `escape_class` (present when
  `escaped`).
- Error modes: kernel refusals map one-to-one to `flow-unmodeled-outcome`,
  `flow-no-match`, `flow-ambiguous-match`, `flow-owned-state-unavailable`,
  `flow-guard-unevaluable`; gate refusals `flow-gate-denied` /
  `flow-gate-indeterminate`; plus the input/model/accessor families
  shared with `next`.

#### `flow read-state`

```
intrastate flow read-state
  (--flow <id> | --model <path>)
  [--artifact role=path ...]
  [--as text|json]
```

- Purpose: diagnostic read without resolution. Invokes only declared
  read accessors; never gate accessors.
- Success `data`: `model`, `revision`, artifact role bindings, `readers[]`
  each with its declared `keys` and the tags it returned (so "absent
  from artifact" is distinguishable from "not requested").
- Error modes: `flow-artifact-*`, `flow-model-*`, accessor exit-3 codes
  (`flow-accessor-timeout`, `flow-accessor-failed`,
  `flow-read-incomplete`).

#### `flow set-state`

```
intrastate flow set-state
  (--flow <id> | --model <path>)
  [--artifact role=path ...]
  [--write name=value ...]
  [--clear <key> ...]
  [--as text|json]
```

- Purpose: persist a decided owned-tag mutation via declared write
  accessors, with read-back verification as the only commit-time check.
- Pre-accessor validation (before any accessor runs): each `--write` /
  `--clear` key must be a declared owned tag served by exactly one
  `[write.<id>]`; each value must be well-formed for its declared kind;
  `--write k=<clear>` refused; a key given twice (across either flag)
  refused.
- `--tag` accepted only as non-written context; never persisted.
- Success only after read-back confirms planned owned-tag values and
  that non-owned tags are unchanged.
- Success `data`: `model`, `revision`, artifact role bindings,
  `writers[]`, requested `writes` and `clear[]`, read-back-confirmed
  owned-tag values.
- Error modes: `flow-write-invalid`, `flow-write-duplicate`,
  `flow-write-unbound`, `flow-clear-unbound`,
  `flow-write-readback-mismatch` (exit 2), `flow-write-readback-incomplete`
  / `flow-write-readback-timeout` (exit 3), plus shared input/model
  families.

### Shared failure envelope

`CLIError` (existing type, extended by exactly one new field):

```
CLIError {
  Code     string            // stable, e.g. "flow-no-match"
  Param    string            // omitempty
  Detail   string            // omitempty
  Hint     string            // omitempty
  Findings []Finding         // omitempty — new in this RDR
}
```

`Finding` (one flat record, `omitempty` optional fields, shared with
RDR 0006 — this RDR populates only the first five):

```
Finding {
  Code     string
  Message  string
  Param    string    `omitempty`
  Locator  string    `omitempty`
  Hint     string    `omitempty`
  // populated by other producers (e.g. RDR 0006), not by this RDR:
  Severity string    `omitempty`
  Model    string    `omitempty`
  Rule     string    `omitempty`
  Key      string    `omitempty`
  Operator string    `omitempty`
  Literal  string    `omitempty`
  Block    string    `omitempty`
  Class    string    `omitempty`
}
```

Exit codes: **3** = environment could not be consulted, safe to re-run
unchanged (`GroupEnvUnavailable`); **2** = request/model wrong or model
said no (`GroupUserEnv` or `GroupInternal`, same exit); **130** = signal
cancel (existing, cited not redefined); **1** = non-`CLIError` failure
(existing). No new exit group.

### GUESS: Go-level handler signatures

The RDR never states Go function signatures for verb handlers — only the
Cobra/`respond`/`clierr` gateway contract and that each verb's `RunE`
calls `respond.ValidateMode(cmd)` then `respond.OK` or
`respond.Fail(cmd, *clierr.CLIError)`. A plausible internal shape
(**GUESS**, not sourced from the RDR):

```go
// GUESS — not specified by the RDR
func runFlowNext(cmd *cobra.Command, req NextRequest) (NextResult, *clierr.CLIError)
func runFlowResolve(cmd *cobra.Command, req ResolveRequest) (ResolveResult, *clierr.CLIError)
func runFlowReadState(cmd *cobra.Command, req ReadStateRequest) (ReadStateResult, *clierr.CLIError)
func runFlowSetState(cmd *cobra.Command, req SetStateRequest) (SetStateResult, *clierr.CLIError)
```

## 2. Three most important internal helper functions

The RDR describes responsibilities, not internal names, for everything
below the verb boundary (it explicitly forbids the CLI from owning table
parsing, guard semantics, or accessor execution). The three most
load-bearing internal *responsibilities* implied by the text, named here
as helpers the implementation must have (functions are **GUESS**-named;
responsibilities are RDR-derived):

1. **`assembleOwnedState` (GUESS name)** — runs each declared read
   accessor over the caller's explicit `--artifact role=path` bindings
   and produces `Input.Owned` for the kernel call. Responsibility:
   enforce that owned state is *never* caller-supplied via `--tag`
   (refuse `flow-tag-owned` before any accessor runs), and that nothing
   is discovered ambiently — only explicit bindings. Used by both
   `flow next` and `flow resolve`.

2. **`runSelectedRowGates` (GUESS name)** — after the kernel selects
   exactly one row (or one escape row), runs every gate declared on that
   row, applies deny-overrides-indeterminate, and reports every gate
   result. Responsibility: keep gate execution strictly post-selection
   (never pre-selection, never used to prune candidates into `no_match`),
   and translate a deny/indeterminate outcome into
   `flow-gate-denied`/`flow-gate-indeterminate` rather than a plan. Used
   by `flow resolve` unconditionally and by `flow next` only under
   `--evaluate-gates`.

3. **`mapRefusal` / error-code translator (GUESS name)** — maps kernel
   `RefusalKinds` (`no_match`, `ambiguous_match`, `owned_state_unavailable`,
   `guard_unevaluable`, `unmodeled_outcome`) one-to-one onto
   `flow-<kind>` `CLIError.Code`s, and maps model-load categories
   (RDR 0002 load categories, RDR 0003 predicate categories, RDR 0008
   `reserved_tag_key`) onto the single code `flow-model-invalid` with one
   `findings[]` entry per category hit. Responsibility: keep exit-code
   assignment stable and centralized (never let a verb decide ad hoc
   whether something is exit 2 or 3) and never launder a gate deny into
   an escapable class.

A fourth strong candidate, not selected only because the task caps the
list at three: a canonical-set-literal encoder (`SetEscapeHTML(false)`,
shared by `clierr.EmitJSON` and `respond.writeJSONLine`) that guarantees
byte-identical set rendering across every emit site — this is explicitly
called out as required (A6, §D13) but is infrastructure the two named
emit functions already own, not a new verb-level helper.

## 3. Data model (persisted / crossing the boundary)

Nothing in this RDR introduces a persistent store of its own — persistence
is delegated to RDR 0004's write accessors over caller-owned artifacts.
The data crossing the CLI boundary is:

### Request-side (flags → typed request)

- Model selection: `--flow <id>` (config-resolved, **GUESS**: resolved
  via `internal/cli/config`, itself called only a "parser placeholder")
  XOR `--model <path>` (explicit file path). Exactly one required.
- `--tag name=value` (repeatable): one *observed* tag; scalar or, for a
  set-kind key, a JSON array literal (e.g. `labels=["a","b"]`), sorted /
  duplicate-free / compact / HTML-escaping-disabled.
- `--artifact role=path` (repeatable): binds a model-declared artifact
  role to a caller-owned file path. No ambient discovery.
- `--outcome <tag>` (resolve only): one recognized-outcome tag.
- `--evaluate-gates` (next only): boolean opt-in.
- `--write name=value` (repeatable, set-state only): scalar or JSON
  array literal for a set-kind key.
- `--clear <key>` (repeatable, set-state only): key-only flag; the
  sentinel value `<clear>` is unauthorable as a `--write` value.

### Identity (Load-Bearing Decisions section)

A CLI request is identified by: verb, model selection (`--flow` id or
`--model` path) + model revision, observed tags, artifact role bindings,
recognized outcome (when applicable), planned writes/clears (when
applicable). Same request over same model revision and same artifact
contents ⇒ same success/refusal, excluding exit-3 environment failures.

### Response-side (verb-specific `data`, inside envelope `{"type":"ok","data":...}`)

Already enumerated per-verb in section 1. Common building blocks:

- `model` (id or path), `revision`
- `observed` (tag map), `owned` (tag map)
- `readers[]` / `writers[]` (accessor identities invoked)
- Set-valued tags always rendered as the canonical JSON array form:
  members sorted, duplicate-free, compact, HTML-escaping disabled (so
  `<`, `>`, `&` serialize literally, never `<`/`>`/`&`).

### Failure-side

`CLIError` + `Finding` as specified in section 1.

### Round-trip / persisted-equivalence invariants (not literally persisted by this RDR, but contractually binding across the boundary)

- `set-state` → `read-state`: value-for-value equality over the owned
  tags the write planned to mutate; byte equality for set values in
  canonical form; a cleared key reads back absent; an empty set reads
  back `[]`.
- `resolve` → (skill) → `set-state`: copy-through byte-identical, because
  one shared encoder renders every canonical set literal.
- `read-state` → `--tag`: only observed tags may be re-paired this way;
  an owned tag read by `read-state` cannot be handed back through
  `--tag` (refused `flow-tag-owned`).

### GUESS: on-disk transition model format

The RDR references a `--model <path>` TOML-like file (e.g.
`./flows/rdr.toml` in the illustrative invocations) loaded by "RDR 0002's
loader" which "takes bytes plus a source id." The RDR explicitly does not
restate RDR 0002's row/table schema, so any field-level shape of the
model file itself is **GUESS** and out of scope for this reconstruction
beyond: it has load categories, a recognized-outcome alphabet, normalized
rows, and `[read.<id>]` / `[write.<id>]` / `[gate.<id>]` capability
blocks with a `keys` list.

## 4. Top-level pseudo-code of the main operation

The RDR gives four verbs; `flow resolve` is the single richest operation
(it exercises read-assembly, kernel selection, gating, and refusal
mapping in one call) and is used here as "the main operation."

```
func FlowResolve(cmd, flags) Result:
    // 1. Output-mode gate (RDR: every verb MUST start here)
    respond.ValidateMode(cmd)                        // text|json

    // 2. Parse + validate request shape, before any I/O or accessor call
    req := parseResolveFlags(flags)
    if req.tag names an owned key or "recognized":
        return respond.Fail(flow-tag-owned | flow-tag-reserved)
    if duplicate --tag name:
        return respond.Fail(flow-tag-duplicate)
    if --outcome missing or empty:
        return respond.Fail(flow-tag-invalid)
    if any --tag set-kind value malformed (not canonical JSON array)
       or scalar/array kind mismatch:
        return respond.Fail(flow-tag-invalid)
    if any --artifact binding malformed:
        return respond.Fail(flow-artifact-invalid)
    if neither/both of --flow/--model given:
        return respond.Fail(flow-model-not-found)

    // 3. Load model (CLI performs file I/O; loader does not)
    bytes := readFile(modelPath)
    model, loadFindings := RDR0002.Load(bytes, sourceID)
    if loadFindings not empty:
        return respond.Fail(flow-model-invalid, findings=loadFindings)

    // 4. Assemble owned state (load→decide fence lives HERE, not in kernel)
    for each declared read accessor the model needs:
        if required --artifact role missing:
            return respond.Fail(flow-artifact-missing, param=role)
    owned, readErr := runDeclaredReadAccessors(model, artifactBindings)
    if readErr is timeout:
        return respond.Fail(flow-accessor-timeout)      // exit 3
    if readErr is execution failure:
        return respond.Fail(flow-accessor-failed)        // exit 3
    if readErr is incomplete key set:
        return respond.Fail(flow-read-incomplete)         // exit 3

    // 5. Pure kernel call — selection only, no accessors, no I/O
    input := Input{Owned: owned, Observed: req.tags, Recognized: req.outcome}
    plan, refusal := internal_resolve.Resolve(input, model.Guards)
    if refusal != nil:
        return respond.Fail(mapKernelRefusalToCode(refusal))
        // flow-unmodeled-outcome | flow-no-match | flow-ambiguous-match |
        // flow-owned-state-unavailable | flow-guard-unevaluable

    // 6. Gates run only after exact-one selection, before the plan is emitted
    gateResults, gateErr := runGatesOnSelectedRow(plan.selectedRow)
    if gateErr is timeout/execution-failure:
        return respond.Fail(flow-accessor-timeout|failed)  // exit 3, not a gate result
    if any gateResults == deny:
        return respond.Fail(flow-gate-denied, findings=onePerGate)
    else if any gateResults == indeterminate:
        return respond.Fail(flow-gate-indeterminate, findings=onePerGate)

    // 7. Success — plan or escaped plan, never both a plan and a refusal
    result := ResolveResult{
        model, revision, observed, owned, readers,
        outcome: req.outcome, rule: plan.ruleID,
        gates: gateResults, next: plan.nextTags,
        writes: plan.writes, clear: plan.clearKeys,
        escaped: plan.isEscape, escape_class: plan.escapeClass,
    }
    return respond.OK(cmd, result)   // renders JSON or text from same struct
```

Notes carried directly from the RDR text (not inferred): step 6 must run
*after* exact-one selection and *before* the plan is emitted (JDR 0001
§D9, cited not restated — the ordering itself is explicit in this RDR).
Deny is never treated as an escape class; it is always a refusal.
`flow resolve` must never print directly, must never initiate skill work,
and must never choose among multiple matching rows (that ambiguity is
`flow-ambiguous-match`, a refusal).
