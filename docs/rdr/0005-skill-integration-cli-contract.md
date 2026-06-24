# Recommendation 0005: Skill Integration CLI Contract

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Final
- **Type**: Feature
- **Profile**: mid — one user-facing CLI integration contract over resolver, accessor, and output seams.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel, 0002-transition-table-as-reviewable-data, 0003-guard-predicate-exhaustiveness
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill author needs to ask "where next", "what outcomes are legal", and "read or persist state" in one deterministic CLI contract instead of re-implementing transition logic. The system-internal requirement is to define the CLI surface over the resolver kernel, including output shape, state binding, and ownership of conditional evaluation.

## Context

### Background

The CLI is the surface over the kernel. `next` must expose the legal-outcome alphabet for constrained decoding, while keeping the resolver semantics in the kernel rather than in skill code.

The real design fork is the integration contract: text versus JSON versus exit-code output, workspace marker versus config versus per-call args for state binding, and whether callers or the CLI evaluate conditional next edges.

### Technical Environment

intrastate is a Cobra-based Go CLI in `internal/cli`. Every verb must start with `respond.ValidateMode(cmd)`, route success through `respond.OK`, route failure through `respond.Fail`, and keep text/json behavior inside `internal/cli/respond`.

## Research Findings

### Investigation

The seed is still current: no resolver CLI, table loader, accessor executor, or
graph lint command exists under `internal/`; the implemented CLI surface is the
root command, `version`, `respond`, `clierr`, and config loading. The proposal
is therefore constrained by the project CLI contract, not by a competing
implemented resolver surface.

Prior art was read before naming approaches. `docs/cli-output-contract.md`
states that the persistent `--as text|json` flag selects the wire format, that
JSON stdout carries exactly one terminal record, and that exit-code mapping
lives in `clierr.ExitCodeFor`. `internal/cli/respond::OK`,
`internal/cli/respond::Fail`, `internal/cli/respond::ValidateMode`,
`internal/cli/clierr::CLIError`, and `internal/cli/root::ExecuteAndEmit`
already provide the failure and output gateway this RDR should reuse rather
than bypass. The current `internal/cli/version::newVersionCmd` text path still
uses `cmd.Println`; that is a source-level drift from the AGENTS guide and the
`respond.OK` convention, so new resolver verbs must add text payload rendering
inside `respond` instead of copying the version shortcut.

The peer RDR split is load-bearing. RDR 0001 owns a stateless resolver kernel
that returns one transition plan or one typed refusal and explicitly delegates
the CLI mapping to this RDR. RDR 0002 owns sparse TOML transition-model data,
legal recognized-outcome alphabets, normalized candidate rows, and rule
identity. RDR 0003 owns symbolic guard predicates and exhaustiveness. RDR 0004
owns declared read/gate/write accessor capabilities and read-back verification.
RDR 0006 owns static graph lint authority. This RDR should expose those
capabilities through a deterministic CLI contract without absorbing their data
formats, predicate grammar, accessor safety model, or lint invariant set.

The strongest external prior is the sibling resolver design. It frames the CLI
as four verbs: `next`, `resolve`, `read-state`, and `set-state`; states that the
flow is "navigated, not orchestrated"; and says `next` emits the
legal-outcome alphabet while an optional constrained decoder remains
downstream. `../state-machines/MODEL-transition.md` sharpens the same split:
`next(state-tags)` returns legal outcome tags and a recognizer prompt,
`resolve(state-tags)` is pure and location-free, and `set-state` persists only
the already decided next tags with read-back verification.

Sibling-path check:

```sh
rg -n "next|legal|outcome|json|exit|code|stdin|stdout|state|artifact|condition|guard|resolver|resolve" \
  ../state-machines/attic/RESOLVER-CLI.md ../state-machines/attic/RESOLVER-DESIGN.md \
  ../state-machines/MODEL-transition.md docs/cli-output-contract.md internal/cli internal/version docs/rdr/000*.md
```

The search found prior-art CLI vocabulary in `../state-machines`, existing
output/failure symbols in `internal/cli`, and peer-RDR contracts, but no
implemented resolver CLI surface or sibling discriminator that already decides
the verb set.

### Key Discoveries

- **Documented** - `docs/cli-output-contract.md` makes `--as text|json`,
  JSON terminal records, advisory streams, and `clierr.ExitCodeFor` the existing
  user-facing output contract.
- **Documented** - `internal/cli/respond::OK`,
  `internal/cli/respond::Fail`, and `internal/cli/respond::ValidateMode` are
  the existing output gateway for new verbs.
- **Documented** - RDR 0001 keeps resolution stateless and delegates CLI
  success/refusal mapping to this RDR.
- **Documented** - RDR 0002 provides the legal recognized-outcome alphabet and
  normalized model rows consumed by `next` and `resolve`.
- **Documented** - RDR 0004 owns read/write/gate accessor execution, gate
  result shape, and read-back verification; this RDR should expose those
  outcomes through verb-appropriate result shapes, not redefine accessor
  safety.
- **Documented** - sibling prior art names the thin CLI as
  `resolve`/`next`/`read-state`/`set-state` and explicitly rejects a runtime,
  driver, or framework.
- **Verified** - a single Cobra command group can expose all four verbs without
  making the CLI the owner of table parsing, guard semantics, or accessor
  execution.
- **Verified** - `next` can return a compact legal-outcome alphabet and
  condition summary that is useful to skills in text mode and complete enough
  for tools in JSON mode.
- **Verified** - resolver/accessor refusal classes can be mapped to stable
  `CLIError.Code` values without requiring new exit-code groups beyond the
  current `clierr.ExitCodeFor` taxonomy.

### Critical Assumptions

- **A1 The existing output gateway can carry all resolver CLI success and
  failure dispositions without direct stdout/stderr writes.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/respond::Success`,
    `internal/cli/respond::OK`, `internal/cli/respond::Fail`,
    `internal/cli/respond::ValidateMode`, `internal/cli/clierr::CLIError`,
    and `internal/cli/root::ExecuteAndEmit` provide the shared success,
    advisory, validation, structured failure, and Cobra-error gateway. Source
    search also found `internal/cli/version::newVersionCmd` still writes text
    success with `cmd.Println`, so implementation must extend `respond.OK` (or a
    respond-owned helper called by `OK`) to render text payloads for new verbs;
    no separate output path is required.
  - **If wrong**: The CLI surface would need a separate output contract or would
    violate the project's no-direct-print convention.
- **A2 The four-verb surface (`next`, `resolve`, `read-state`, `set-state`) is
  the minimal complete skill integration contract.**
  - **Status**: Verified
  - **Method**: MVV Test
  - **Evidence**: Minimum Viable Validation scenario "fixture-backed flow
    proves all four verbs through the production Cobra path" covers the complete
    loop: `flow next` asks legal outcomes, `flow resolve` resolves one
    recognized outcome, `flow read-state` reads fixture artifact tags, and
    `flow set-state` persists planned owned-tag writes with read-back. The
    sibling prior `../state-machines/MODEL-transition.md` uses the same four
    operations and assigns no fifth required operation to the skill integration
    loop.
  - **If wrong**: Skill authors would still re-implement part of transition or
    state-binding logic outside intrastate.
- **A3 `next` can expose the legal-outcome alphabet without owning guard
  semantics or hiding unresolved facts.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Technical Design` defines a recognized outcome
    alphabet and normalized candidate rows with rule ids/source locators; RDR
    0002 `Normative Contracts` requires exact-one matching and refusals for zero
    or multiple matches. RDR 0003 `Approach` and `Technical Design` make guard
    facts symbolic tag-set predicates owned outside this CLI. Therefore `next`
    can enumerate the alphabet, surface supplied or unresolved facts, and map
    accessor refusals without becoming the owner of guard truth.
  - **If wrong**: The CLI would either under-inform constrained decoding or
    incorrectly become the owner of conditional evaluation.
- **A4 The resolver kernel can remain pure while the CLI binds accessors only at
  verb boundaries.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0001 `Technical Design` and `Normative Contracts` keep
    artifact discovery and accessor execution outside the resolver kernel; the
    kernel accepts accessor-produced facts and returns structured
    success/refusal values. RDR 0004 `Technical Design` and `Normative
    Contracts` define separate read, gate, and write accessor capabilities over
    caller-supplied artifact roles: reads return typed tags, gates return
    allow/deny/indeterminate, and writes apply planned owned-tag mutations with
    same-role read-back. Therefore the CLI can bind read accessors in
    `flow read-state`, declared gate accessors in `flow next` / `flow resolve`,
    and write accessors in `flow set-state` without making the kernel stateful
    or coercing gate results into tag values.
  - **If wrong**: The CLI would have to mix artifact discovery, accessor
    execution, and kernel selection in one command, weakening replay safety.
- **A5 Stable CLI error codes can distinguish bad input, unknown outcome,
  zero-match, multi-match, unavailable accessor, gate-indeterminate, and
  write-readback-mismatch failures.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr::CLIError` carries a stable string
    `Code` plus optional `Param`, `Detail`, and `Hint`; `ErrorCode` exposes the
    code for branching; `ExitCodeFor` maps existing groups to exit 2 for
    user/model/input refusals, exit 3 for environment unavailability, and 130
    for interruption. Resolver/accessor-specific values only require new
    `Code` constants or literals, not new envelope fields or exit groups.
  - **If wrong**: Scripted skill calls could not branch deterministically on
    resolver failure classes.
- **A6 The pinned MVP request grammar, minimum success payload fields, and
  stable flow error-code spellings are sufficient for the first fixture-backed
  CLI implementation.**
  - **Status**: Accepted
  - **Method**: Design Decision
  - **Evidence**: Stage 6 Reconcile classifies the pinned `--flow`, repeated
    `--tag name=value`, `--artifact role=path`, and `--write name=value` grammar,
    the minimum JSON `data` fields in Technical Design, and the minimum
    `flow-*` code strings in Failure Modes as this RDR's implementation contract.
    The rejected alternative is leaving grammar, payload fields, or error-code
    spelling to implementation-time invention. The Implementation Plan and
    Validation sections keep the fixture-backed production-Cobra MVV as the test
    that proves an implementation conforms to this chosen contract.
  - **If wrong**: Implementers would still need to invent request grammar,
    payload fields, or error-code spellings during code work.

### Reconciliation Report

| Item | Source | Disposition | Evidence pointer or plan |
| --- | --- | --- | --- |
| A4 verb-boundary accessor binding after the gate-accessor re-entry fix | 1, 2, 4 | VERIFIED | Peer RDR evidence recorded in A4: RDR 0001 keeps the kernel pure and RDR 0004 distinguishes read, gate, and write capability result shapes over caller-supplied artifacts. |
| A6 pinned MVP grammar, payload fields, and stable `flow-*` code spellings | 1, 2, 4 | ACCEPTED | Design Decision recorded in A6. The contract pins the MVP grammar, success payload minima, and refusal-code spellings; the rejected alternative is implementation-time invention. |

**Method vocabulary** (pick exactly one per assumption):

- **Source Search** — verified against dependency
  source code. Evidence: a greppable `path::Symbol`
  (function/type/const name), **not a bare `file:line`**;
  a commit-SHA permalink only for audit/traceability.
  Standard for libraries. (Why symbol not line: flow
  README *Doctrine*.)
- **Spike** — verified by running code against a live
  service or fixture. Evidence: command run + path to
  captured output.
- **Prior Art** — same property holds in ≥1 named
  external system. Evidence: system + section/page.
- **Derivation** — pure math or proof. Evidence: the
  derivation, shown inline.
- **Design Decision** — a scoping choice this RDR is
  *making* (not *verifying*). Evidence: the decision
  and the alternative explicitly rejected.
- **Peer RDR** — relies on a property defined in
  another RDR. Evidence: RDR ID + section.
- **MVV Test** — the property is testable via the
  Minimum Viable Validation, and the test
  is named in this RDR's Validation section (pending
  implementation at lock time). Evidence: test name.
- **Docs Only** — documentation reading alone.
  **Insufficient** for load-bearing assumptions; allowed
  only when paired with a Spike or Source Search plan
  in the Evidence line.

A `Method: Source Search` whose Evidence cites this
same RDR file — or any path under the RDR's artifact
directory — is self-reference and not Verified. The
cited proof must also support **the specific claim**,
not an adjacent one: confirming a neighboring fact and
stamping the assumption `Verified` is not verification.
The cited symbol must resolve on `main` (a renamed,
deleted, or never-built symbol fails the check).

Any exactness claim such as all/every, first/nearest,
byte-identical, lossless, canonical, deterministic, or
stable order must be covered by a Critical Assumption
Evidence Record or by the Minimum Viable Validation.

## Proposed Solution

### Approach

Expose a thin `flow` command group with four verbs over the resolver seams:
`flow next`, `flow resolve`, `flow read-state`, and `flow set-state`. The CLI
maps user flags and config-discovered model names into typed calls, then renders
results through `respond`. It does not own transition-model representation,
guard semantics, accessor safety, graph lint, skill execution, or constrained
decoding. It is the deterministic integration contract a skill can call when it
needs the legal outcome alphabet, one resolved next transition, or a state I/O
operation.

`flow next` answers "what outcomes are legal from this supplied state tag-set?"
It returns recognized outcome tags plus any conditional summaries needed to
understand why multiple outcomes or next tags are possible. It may include the
recognizer prompt supplied by the transition model when RDR 0002 provides one,
but it never calls a model.

`flow resolve` answers "given this recognized outcome and supplied state tags,
what exact next tag-set or refusal results?" It delegates edge selection to the
RDR 0001 kernel and maps kernel refusals to stable `CLIError.Code` values.

`flow read-state` and `flow set-state` bind the model's accessor names to
caller-supplied artifact roles. `read-state` invokes declared read accessors
and returns the current tag-set read from the artifact boundary. It does not
invoke gate accessors, because gate accessors return allow, deny, or
indeterminate rather than tag values. `set-state` persists a decided owned-tag
mutation through RDR 0004's accessor executor and returns success only after
read-back verification succeeds.

### Technical Design

The CLI layer is a translator, not a decision engine. Each verb starts with
`respond.ValidateMode(cmd)`, validates arguments into typed request values,
loads the configured transition model, calls the owning internal package, and
renders one terminal disposition. Success uses `respond.OK`; user-visible
failure uses `respond.Fail(cmd, &clierr.CLIError{...})`. The text mode contract
is human-scannable but still derived from the same request/response structs as
JSON mode.

State binding is explicit. Commands accept a model or flow identifier plus
state tags supplied as repeated `--tag name=value` flags in the MVP. Commands
that need accessor execution accept caller-supplied artifact role bindings as
`--artifact role=path`: `read-state` for read accessors, `next` / `resolve` for
declared gate accessors needed to evaluate candidate facts, and `set-state` for
write accessors plus read-back. A later structured input file may add bulk tag
input, but it is not required for the first implementation slice. The chosen
approach deliberately avoids ambient artifact discovery in `resolve`; location
comes only from explicit artifact bindings, and the resolver kernel remains
pure over model data, supplied facts, and accessor results.

The four verbs share a small error taxonomy. Input and config errors remain in
`GroupUserEnv`. External/accessor unavailability uses `GroupEnvUnavailable` when
the environment, not the model, is unavailable. Internal parse or invariant
failures use the existing internal/user-env mapping unless Resolve finds a need
for a narrower group. Error envelopes stay append-only through `CLIError`
fields.

The MVP request grammar is intentionally narrow. `--flow <id>` selects the model
or flow definition. `--tag name=value` supplies one scalar tag fact; duplicate
tag names in the same request are refused until the transition-model layer
exposes a first-class structured literal for set-valued tags. Verbs that invoke
accessors accept `--artifact role=path`, where `role` is the model-declared
artifact role and `path` is the caller-owned artifact path supplied to the
accessor executor. `flow set-state` accepts planned owned-tag mutations as
`--write name=value`; those writes are distinct from `--tag`, which remains
context already known to the caller.

Successful JSON payloads use the existing `{"type":"ok","data":...}` envelope.
The minimum `data` shape is:

- `next`: selected flow/model identity, supplied tags, `outcomes[]`, and
  `candidates[]`. Each outcome has the recognized outcome tag and optional
  recognizer text supplied by the model. Each candidate summary has source rule
  identity, the recognized outcome it belongs to, required supplied facts,
  unresolved guard/accessor facts, evaluated gate facts when artifact bindings
  were supplied, and preview next tags or write targets when the normalized model
  exposes them without evaluating missing facts.
- `resolve`: selected flow/model identity, supplied tags, recognized outcome,
  artifact role bindings used for gate facts, matched rule identity, next tags,
  and planned owned-tag writes.
- `read-state`: selected flow/model identity, artifact role bindings, read
  accessor identities, and read tags.
- `set-state`: selected flow/model identity, artifact role bindings, requested
  owned-tag writes, and read-back-confirmed owned-tag values.

Text mode renders the same result content in a human-scannable order; it does
not invent fields absent from the JSON payload.

Sibling-path check for a new discriminator or identity rule found no existing
resolver CLI signal in `internal/`. The adjacent decision signal is the
documented `--as` output gateway and peer-RDR split, so this proposal reuses
those instead of introducing a parallel output or state-binding mechanism.

#### Normative Contracts

```normative
The CLI MUST expose one command group for skill integration with these verbs:
next, resolve, read-state, and set-state.

All four verbs MUST start RunE by calling respond.ValidateMode(cmd), MUST route
success through respond.OK, MUST route user-facing failure through
respond.Fail(cmd, *clierr.CLIError), and MUST set SilenceErrors and
SilenceUsage.

Under --as=json, each successful invocation MUST emit exactly one stdout JSON
terminal envelope with type "ok" and verb-specific data. Under --as=text, each
successful invocation MUST emit human output derived from the same verb-specific
result. Failures MUST use the existing CLIError JSON/text envelope defined by
docs/cli-output-contract.md and internal/cli/clierr.

flow next MUST return the legal recognized-outcome alphabet for the supplied
state tag-set, plus candidate summaries containing source rule identity,
required supplied facts, unresolved guard/accessor facts, and preview next tags
or write targets when those can be read from normalized model data without
evaluating missing facts. When the selected model references declared gate
accessors and the caller supplies matching `role=path` artifact bindings, `flow
next` MAY invoke those gate accessors through the RDR 0004 executor and surface
allow/deny/indeterminate as candidate facts or stable CLIError refusals. It MUST
NOT invent guard facts that were neither supplied nor produced by a declared
gate accessor.

flow resolve MUST return exactly one resolved next tag-set or exactly one
CLIError refusal mapped from the resolver kernel. It MUST NOT discover
artifacts, run read/write accessors, print directly, initiate skill work, or
choose among multiple matching rows. It MAY invoke declared gate accessors from
caller-supplied `role=path` artifact bindings before the pure kernel call when
the matched candidate requires a gate fact; gate denied or indeterminate results
MUST surface as stable CLIError refusals.

flow read-state MUST invoke only declared read accessors over caller-supplied
`role=path` artifact bindings and return the read tag-set or a stable CLIError
failure. It MUST NOT invoke gate accessors or coerce gate allow, deny, or
indeterminate results into tag values.

flow set-state MUST invoke only declared write accessors over caller-supplied
`role=path` artifact bindings and planned owned-tag `--write name=value`
mutations. It MUST report success only after the accessor layer's read-back
verification confirms the expected owned-tag values. It MUST NOT treat
context-only `--tag` values as writes.
```

#### Load-Bearing Decisions

- **Identity** - a CLI request is identified by verb, flow/model identifier,
  supplied state tags, supplied recognized outcome when applicable, and supplied
  artifact role bindings when applicable. The same request over the same model
  revision must produce the same success or refusal, excluding external accessor
  unavailability.
- **Wire / byte format** - JSON and text rendering conform to
  `docs/cli-output-contract.md`; this RDR introduces verb-specific JSON `data`
  payloads but not a new terminal envelope.
- **Naming** - the user-facing group is `flow`, with verbs `next`, `resolve`,
  `read-state`, and `set-state`. Rejected group names: `state` because only two
  verbs perform state I/O, and `run` because the CLI does not orchestrate skill
  work.
- **Selection / predicate** - `flow resolve` succeeds only when the kernel
  reports exactly one matching row. Zero matches, multiple matches, unknown
  outcomes, missing facts, or unevaluable guards become typed refusals rather
  than tie-breaks.

#### Round-Trip / Inverse Invariants

`set-state` followed by `read-state` must return the expected owned-tag values
for the artifact role that was written. The equality is value-for-value over the
owned tags the write planned to mutate, not byte-identical artifact content.
RDR 0004 owns the accessor read-back semantics; this RDR owns surfacing the
success/failure through the CLI contract.

#### Illustrative Code

Illustrative invocation shapes:

```sh
intrastate flow next --flow rdr --tag stage=resolve --tag profile=small --as=json
intrastate flow resolve --flow rdr --tag stage=prelock --tag iter=2 --outcome verdict-flapping
intrastate flow read-state --flow rdr --artifact state=./docs/rdr/0005-skill-integration-cli-contract.md
intrastate flow set-state --flow rdr --artifact state=./docs/rdr/0005-skill-integration-cli-contract.md --write status=Final
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Text/json output gateway | Existing `respond` / `clierr` | Verified | Reuse; no new terminal envelope. Implementation must route text success payloads through `respond.OK` or a respond-owned helper. |
| Stateless resolver kernel | RDR 0001 | Verified peer RDR | `flow resolve` delegates selection and refusal classes. |
| Transition model and outcome alphabet | RDR 0002 | Verified peer RDR | `flow next` depends on model-provided outcome and row data. |
| Guard predicate semantics | RDR 0003 | Verified peer RDR | CLI reports conditions but does not own guard truth. |
| Accessor execution and read-back | RDR 0004 | Verified peer RDR | `read-state` exposes read-accessor tag values; `set-state` exposes write/read-back outcomes; gate results surface through candidate summaries or typed refusals. |
| Graph lint authority | RDR 0006 | Out of scope for this contract | CLI may later expose lint, but this RDR does not own lint guarantees. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Output mode validation | `internal/cli/respond::ValidateMode` | Current verbs must call it manually. | Reuse | Every new verb starts with it. |
| Success rendering | `internal/cli/respond::OK` | Text mode currently emits only advisories, while `version` uses `cmd.Println`. | Extend `respond.OK` or a respond-owned helper | New verbs must not use direct Cobra printing for text payloads. |
| Failure rendering | `internal/cli/respond::Fail` and `internal/cli/clierr::CLIError` | No resolver-specific codes yet. | Extend codes only | Add stable refusal codes as needed; no new envelope fields or exit groups required. |
| Root wiring | `internal/cli/root::NewRootCmd` / `ExecuteAndEmit` | Only `version` is registered today. | Extend | Register `flow` group and keep Cobra errors structured. |
| Config discovery | `internal/cli/config` | Parser placeholder; no transition-model config yet. | Extend later | CLI can start with explicit flags; config-backed model lookup may follow. |

### Decision Rationale

Choose the four-verb thin CLI because it matches the user outcome without
moving ownership out of the peer seams. Skills get one deterministic contract
for legal outcomes, resolution, and state I/O. The resolver remains stateless,
accessors remain the only state-binding authority, and all user-visible output
still conforms to the existing CLI envelope. The sibling prior art also points
to this exact split: a table plus four verbs gives the wall-clock win of
removing repeated position re-derivation without turning intrastate into a
workflow runner.

Rejected alternatives fail one of those boundaries. A JSON-only API is simpler
for scripts but violates the CLI's established text/json mode contract. A single
`resolve` mega-command hides state I/O and makes conditional alphabet queries
harder for constrained recognition. A higher-level `run` or `advance` command
would initiate workflow action and collapse the project distinction between
recognizing an outcome, resolving it, and performing the next skill.

Premortem: this ships and fails because the four verbs are present but
semantically thin: `next` omits enough condition detail to constrain a skill,
`resolve` maps all refusals to generic errors, and `set-state` looks successful
when read-back did not prove the owned tags changed. The recommendation survives
that failure mode because the contract makes condition summaries, stable
refusal codes, and read-back-gated success normative; Resolve verified those
points so they are not left as implementation taste.

## Alternatives Considered

### Alternative 1: JSON-only machine API

**Description**: Add resolver commands that always emit JSON records and ignore
the global `--as` mode for this command family.

**Pros**:

- Simplest contract for skill callers.
- Avoids text formatting decisions for conditional outcome summaries.

**Cons**:

- Contradicts `docs/cli-output-contract.md` and the root `--as text|json`
  convention.
- Makes this command family an exception every future CLI verb has to remember.
- Still needs structured CLI errors and exit codes, so it does not avoid the
  existing gateway.

**Reason for rejection**: The project already has a mode-aware output contract;
resolver verbs should reuse it.

### Alternative 2: Single `resolve` Command With Modes

**Description**: Put legal-outcome enumeration, resolution, state reads, and
persistence behind one command with action flags such as `--next`,
`--read-state`, and `--set-state`.

**Pros**:

- One command name for skill authors to learn.
- Shared setup flags could be centralized.

**Cons**:

- Blurs separate questions: "what can happen", "what did happen", "what is the
  stored state", and "persist this decision".
- Makes argument validation branch-heavy and easier to misuse.
- Encourages `resolve` to know about artifact binding, which RDR 0001 keeps out
  of the pure kernel.

**Reason for rejection**: Separate verbs preserve the conceptual split and make
invalid combinations easier to refuse.

### Alternative 3: `run` / `advance` Workflow Driver

**Description**: Add a command that reads state, asks or accepts an outcome,
resolves the next state, persists it, and possibly invokes the next skill or
shell action.

**Pros**:

- Attractive single-call experience for happy-path automation.
- Could hide state I/O details from skills entirely.

**Cons**:

- Crosses into orchestration, which RDR 0001 and the sibling resolver prior art
  explicitly reject.
- Would make intrastate responsible for judgment gates and next-work execution
  rather than deterministic resolution.
- Raises blast radius by combining reads, decisions, writes, and side effects.

**Reason for rejection**: The user outcome is deterministic transition support
for skills, not an orchestrator.

### Briefly Rejected

- **Exit-code-only output**: cannot carry legal outcome alphabets, conditional
  summaries, or typed tag sets.
- **Library-only API with no CLI**: misses the skill integration surface and
  would force each harness to rebuild command-line wiring.
- **Let skills call the kernel package directly**: ties prompt authors to Go
  package structure and bypasses the stable output/error contract.

## Trade-offs

### Consequences

- Skills get one deterministic surface for the closed-world transition question
  instead of copying transition logic into prompts.
- The CLI remains mode-aware and scriptable through the existing text/json
  contract.
- The verb set introduces user-facing surface area before the peer packages are
  implemented, so Resolve must keep the normative contract high-level enough not
  to over-specify internal types.
- Text-mode summaries for conditional `next` output need careful design so they
  stay readable without becoming a second grammar.

### Risks and Mitigations

- **Risk**: `flow next` becomes a second guard evaluator by trying to resolve
  conditionals without supplied facts.
  **Mitigation**: The contract requires conditional summaries and forbids
  evaluating missing guard/accessor facts.
- **Risk**: Resolver refusals collapse into generic `command-error`.
  **Mitigation**: Add stable `CLIError.Code` values for each refusal class and
  cover them in the MVV.
- **Risk**: Text and JSON outputs diverge semantically.
  **Mitigation**: Derive both renderings from one typed result per verb and
  test both modes.
- **Risk**: `set-state` is mistaken for an orchestration command.
  **Mitigation**: It only persists planned owned-tag `--write` mutations through
  declared write accessors and depends on RDR 0004 read-back verification.

### Failure Modes

Visible failures include unknown command arguments, model/config not found,
unknown outcome, zero matching row, multiple matching rows, missing supplied
guard facts, accessor unavailable, gate indeterminate, and write read-back
mismatch. Each should surface as one `CLIError.Code` with a clear message,
detail, and hint where useful.

Minimum stable code strings:

| Failure | Code | Group |
| --- | --- | --- |
| malformed `--tag` / `--write` | `flow-tag-invalid` | `GroupUserEnv` |
| duplicate supplied tag name | `flow-tag-duplicate` | `GroupUserEnv` |
| malformed `--artifact` binding | `flow-artifact-invalid` | `GroupUserEnv` |
| selected flow/model not found | `flow-model-not-found` | `GroupUserEnv` |
| recognized outcome absent from model alphabet | `flow-outcome-unknown` | `GroupUserEnv` |
| resolver found no matching row | `flow-zero-match` | `GroupUserEnv` |
| resolver found multiple matching rows | `flow-multi-match` | `GroupUserEnv` |
| required supplied fact is missing | `flow-fact-missing` | `GroupUserEnv` |
| guard cannot be evaluated from supplied facts | `flow-guard-unevaluable` | `GroupUserEnv` |
| accessor cannot run because its environment is unavailable | `flow-accessor-unavailable` | `GroupEnvUnavailable` |
| accessor gate result is indeterminate | `flow-gate-indeterminate` | `GroupUserEnv` |
| write read-back does not confirm expected owned tags | `flow-write-readback-mismatch` | `GroupUserEnv` |

The main silent-failure risk is treating a failed accessor or ambiguous row as a
successful transition. The recovery path is refusal-first: the command exits
non-zero, emits the structured error envelope, and leaves skill judgment or
model repair to the caller. Diagnosis starts from the error code, model/rule id
when available from RDR 0002, and accessor identity/artifact role when the
failure comes from RDR 0004.

## Implementation Plan

### Prerequisites

- [x] All Critical Assumptions have terminal dispositions; A1-A5 are verified,
      and A6 is accepted as the pinned implementation contract.
- [x] RDR 0001, 0002, 0003, and 0004 expose enough internal contracts for the
      CLI to call without owning their semantics.
- [ ] Add text success payload rendering through `respond.OK` or a
      respond-owned helper used by `respond.OK`; do not print directly from
      resolver verbs.

### Minimum Viable Validation

Implement one fixture-backed flow and prove all four verbs through the
production Cobra path: `flow next` returns the legal outcome alphabet,
`flow resolve` maps one outcome to one next tag-set, `flow read-state` reads the
fixture artifact tags, and `flow set-state` persists a planned owned-tag write
then read-back-verifies it. Run the same happy path and at least one typed
refusal in `--as=text` and `--as=json`, including the MVP `--tag`,
`--artifact`, and `--write` grammar.

### Phase 1: Command Contract Skeleton

Add the `flow` command group and the four verbs with Cobra argument validation,
`respond.ValidateMode`, `respond.OK`, `respond.Fail`, and focused command tests.

### Phase 2: Kernel And Model Binding

Wire `next` and `resolve` to the RDR 0001/0002/0003 interfaces using a fixture
model before broad config discovery.

### Phase 3: Accessor Binding

Wire `next`, `resolve`, `read-state`, and `set-state` to RDR 0004's accessor
executor only for the capability each verb may invoke. Restrict `read-state` to
declared read accessors, allow `next` / `resolve` to request declared gate
accessors from explicit artifact bindings, and preserve gate-indeterminate and
read-back failures as typed CLI refusals from the verb path that invoked them.

### Phase 4: Error Taxonomy And Docs

Add resolver/accessor CLI error codes, update the CLI output contract only if
the existing envelope needs an append-only field, and document illustrative
invocations.

### Day 2 Operations

This RDR does not create a persistent store. It exposes existing transition
model files and caller-owned artifacts through CLI operations.

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Transition model files | Covered by repository tools | Covered by repository tools | Covered by repository tools | In scope through RDR 0006 lint and CLI MVV | Covered by version control |
| Caller artifacts mutated by `set-state` | Caller-owned | Caller-owned | Caller-owned | In scope through read-back verification | Caller-owned |

### New Dependencies

No new third-party dependency is selected at Propose. Cobra, `respond`, and
`clierr` already exist. Resolve may identify a parsing dependency through RDR
0002, but this RDR should not introduce one for the CLI surface alone.

## Validation

### Testing Strategy

Command tests should exercise the `flow` command group through the production
Cobra path rather than calling renderers or kernel functions directly. Coverage
must include argument validation, output mode validation, success rendering, and
typed refusal mapping for each verb.

1. **Scenario**: `flow next` over the fixture model in `--as=json` and
   `--as=text`.
   **Expected**: both modes report the same recognized-outcome alphabet and
   candidate summaries through the standard success envelope/rendering path;
   JSON includes `outcomes[]` and `candidates[]`.
2. **Scenario**: `flow resolve` over the fixture model with one recognized
   outcome that matches exactly one row.
   **Expected**: the command returns the expected next tag-set and no read/write
   accessor is attempted; when the matched row needs a gate fact, only a
   declared gate accessor over a supplied artifact binding may run.
3. **Scenario**: `flow resolve` with an unknown outcome, zero-match row, and
   multi-match row.
   **Expected**: each refusal maps to a stable `CLIError.Code` and non-zero
   exit behavior under both output modes.
4. **Scenario**: `flow read-state` and `flow set-state` over a fixture artifact
   and declared accessor roles.
   **Expected**: `--artifact role=path` bindings are validated, `read-state`
   invokes only declared read accessors and returns the artifact tag-set, and
   `--write name=value` mutations report success only after read-back
   verification proves the expected owned-tag values.
5. **Scenario**: read accessor unavailable, gate indeterminate during `next` /
   `resolve` accessor fact evaluation, and write read-back mismatch.
   **Expected**: each failure remains a typed CLI refusal, not a successful
   transition or a coerced `read-state` tag payload.

### Performance Expectations

No throughput target is load-bearing for this RDR. The command path should stay
single-invocation deterministic: parse inputs, load the selected model, call one
owning package operation, and render one terminal result. Any accessor latency
belongs to RDR 0004's execution contract; this RDR only requires that an
unavailable or indeterminate accessor surfaces as a typed refusal.

## Finalization Gate

> Complete each item with a written response before
> marking this RDR as **Final**. Written responses
> prevent rubber-stamping and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses below.

### Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution. The research points to a thin four-verb CLI, reuse of the
existing output gateway, pure resolver delegation, and accessor-owned state
binding. The scoped re-entry clarified that `read-state` binds only read
accessors, while gate accessors keep their RDR 0004 allow/deny/indeterminate
shape and surface through candidate summaries or typed refusals.

### Assumption Verification

A1-A5 are Verified and each has a non-empty "If wrong" branch. A4 was
rechecked after the 3amigo iter-3 gate-accessor fix against RDR 0001 and RDR
0004: the kernel remains pure, read/gate/write accessor result shapes remain
distinct, and gate results surface only through `flow next` / `flow resolve`
candidate facts or typed refusals. A6 is Accepted as a design decision because
the 3amigo pass pinned exact request grammar, success payload, and error-code
spellings that the implementation must satisfy. No assumption uses `Docs Only`;
A1 and A5 cite source search against the CLI gateway and error taxonomy, A2
cites the named MVV, A3-A4 cite peer RDR contracts, and A6 names the rejected
alternative. No evidence cites this RDR or its artifact directory as proof.
Resolve found one implementation requirement rather than a refutation: text
success payloads must be rendered through `respond.OK` or a respond-owned helper
instead of direct Cobra printing.

### Scope Verification

The Minimum Viable Validation is in scope for implementation: one fixture-backed
flow must prove `flow next`, `flow resolve`, `flow read-state`, and
`flow set-state` through the production Cobra path in both output modes,
including at least one typed refusal.

### Cross-Cutting Concerns

- **Versioning**: verb-specific JSON `data` payloads are append-only under the
  existing CLI output envelope.
- **Incremental adoption**: explicit flags and fixture-backed model loading can
  ship before broader config discovery.
- **Secret/credential lifecycle**: this RDR does not introduce credentials;
  accessor execution and external availability belong to RDR 0004.
- **Canonical-form / determinism**: the deterministic claim is request-level
  semantic determinism over the same model revision, not byte-identical output
  or content-addressed identity.

### Proportionality

This RDR is right-sized for one load-bearing contract: the user-facing CLI
integration surface over the resolver, accessor, and output seams. It does not
own kernel selection semantics, transition-model format, guard predicate
meaning, accessor safety, or graph lint invariants. The `mid` profile remains
appropriate because the contract is user-facing but not foundational and carries
no prior accretion in Seam Lineage.

## References

- `docs/cli-output-contract.md`
- `internal/cli/respond::OK`
- `internal/cli/respond::Fail`
- `internal/cli/respond::ValidateMode`
- `internal/cli/clierr::CLIError`
- `internal/cli/clierr::ExitCodeFor`
- `internal/cli/root::ExecuteAndEmit`
- `internal/cli/version::newVersionCmd`
- RDR 0001, 0002, 0003, 0004, and 0006
- `../state-machines/attic/RESOLVER-CLI.md`
- `../state-machines/attic/RESOLVER-DESIGN.md`
- `../state-machines/MODEL-transition.md`
