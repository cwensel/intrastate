# Recommendation 0005: Skill Integration CLI Contract

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft
  <!--
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 08.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 08.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 8 (re-lock) parse the qualifier.
    The Stage 8 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 08.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  -->
- **Type**: Feature
- **Profile**: mid — exposes the resolver contract through a user-facing CLI surface.
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: small = one contract, no user-facing
  surface (skips Stage 5); mid = one contract + user-facing
  surface OR locks a contract; large = locks an enum/hash/
  format/grammar/destructive-op; foundational = cross-RDR
  producer / spans modules.
  (2) accretion axis (HARD floor): if `Seam Lineage` below
  carries ≥2 closed prior point-fixes at this locus, Profile
  is floored at FOUNDATIONAL regardless of the contract axis
  — a seam with prior point-fixes is never small/mid (it
  spans the prior RDRs/patches = the matrix's cross-RDR
  trigger). The only escape is a written accretion disposition
  in the Seam Lineage field. This floor is what stops a
  "one contract → mid" sizing from under-gating an accreting
  seam.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
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
`respond.OK` convention, so Resolve must decide whether to update the example
or allow `respond.OK` to carry text payloads for new verbs.

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
- **Documented** - RDR 0004 owns read/write/gate accessor execution and
  read-back verification; this RDR should expose those outcomes, not redefine
  accessor safety.
- **Documented** - sibling prior art names the thin CLI as
  `resolve`/`next`/`read-state`/`set-state` and explicitly rejects a runtime,
  driver, or framework.
- **Assumed** - a single Cobra command group can expose all four verbs without
  making the CLI the owner of table parsing, guard semantics, or accessor
  execution.
- **Assumed** - `next` can return a compact legal-outcome alphabet and
  condition summary that is useful to skills in text mode and complete enough
  for tools in JSON mode.
- **Assumed** - resolver/accessor refusal classes can be mapped to stable
  `CLIError.Code` values without requiring new exit-code groups beyond the
  current `clierr.ExitCodeFor` taxonomy.

### Critical Assumptions

- **A1 The existing output gateway can carry all resolver CLI success and
  failure dispositions without direct stdout/stderr writes.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Pending: Resolve must confirm
    `internal/cli/respond::OK`, `internal/cli/respond::Fail`,
    `internal/cli/respond::ValidateMode`, and
    `internal/cli/clierr::CLIError` cover text/json success, advisories, and
    structured failures for the four verbs.
  - **If wrong**: The CLI surface would need a separate output contract or would
    violate the project's no-direct-print convention.
- **A2 The four-verb surface (`next`, `resolve`, `read-state`, `set-state`) is
  the minimal complete skill integration contract.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: Pending: Resolve must name an MVV that uses one representative
    RDR or kata fixture to ask legal outcomes, resolve one outcome, read state,
    and persist the decided next tags.
  - **If wrong**: Skill authors would still re-implement part of transition or
    state-binding logic outside intrastate.
- **A3 `next` can expose the legal-outcome alphabet without evaluating
  conditionals owned by guards or accessors.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Pending: RDR 0002 must provide recognized-outcome alphabets
    and normalized conditional rows, while RDR 0003 owns guard semantics.
  - **If wrong**: The CLI would either under-inform constrained decoding or
    incorrectly become the owner of conditional evaluation.
- **A4 `resolve` can remain pure over supplied tags and model data, with
  artifact binding handled by `read-state` / `set-state`.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: Pending: RDR 0001 must preserve location-free resolution
    inputs and RDR 0004 must provide accessor-read values and persistence
    disposition.
  - **If wrong**: The CLI would have to mix artifact discovery, accessor
    execution, and kernel selection in one command, weakening replay safety.
- **A5 Stable CLI error codes can distinguish bad input, unknown outcome,
  zero-match, multi-match, unavailable accessor, gate-indeterminate, and
  write-readback-mismatch failures.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Pending: Resolve must confirm `internal/cli/clierr::CLIError`
    and `internal/cli/clierr::ExitCodeFor` can carry these codes with the
    current exit groups, or name the minimal extension.
  - **If wrong**: Scripted skill calls could not branch deterministically on
    resolver failure classes.

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
caller-supplied artifact roles. `read-state` returns the current tag-set read
from the artifact boundary. `set-state` persists a decided owned-tag mutation
through RDR 0004's accessor executor and returns success only after read-back
verification succeeds.

### Technical Design

The CLI layer is a translator, not a decision engine. Each verb starts with
`respond.ValidateMode(cmd)`, validates arguments into typed request values,
loads the configured transition model, calls the owning internal package, and
renders one terminal disposition. Success uses `respond.OK`; user-visible
failure uses `respond.Fail(cmd, &clierr.CLIError{...})`. The text mode contract
is human-scannable but still derived from the same request/response structs as
JSON mode.

State binding is explicit. Commands accept a model or flow identifier plus
state tags supplied as flags or a structured input file, and accept artifact
role bindings only for accessor verbs. The chosen approach deliberately avoids
ambient artifact discovery in `resolve`; location belongs to accessor I/O, not
to the pure resolver.

The four verbs share a small error taxonomy. Input and config errors remain in
`GroupUserEnv`. External/accessor unavailability uses `GroupEnvUnavailable` when
the environment, not the model, is unavailable. Internal parse or invariant
failures use the existing internal/user-env mapping unless Resolve finds a need
for a narrower group. Error envelopes stay append-only through `CLIError`
fields.

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
state tag-set, plus conditional next summaries when more than one next tag-set
depends on unresolved guard/accessor facts. It MUST NOT evaluate guard facts
that were not supplied.

flow resolve MUST return exactly one resolved next tag-set or exactly one
CLIError refusal mapped from the resolver kernel. It MUST NOT discover
artifacts, run accessors, print directly, initiate skill work, or choose among
multiple matching rows.

flow read-state MUST invoke only declared read or gate accessors over
caller-supplied artifact role bindings and return the read tag-set or a stable
CLIError failure.

flow set-state MUST invoke only declared write accessors over caller-supplied
artifact role bindings and report success only after the accessor layer's
read-back verification confirms the expected owned-tag values.
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
intrastate flow read-state --flow rdr --artifact state:RDR_FILE
intrastate flow set-state --flow rdr --artifact state:RDR_FILE --tag status=Final
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Text/json output gateway | Existing `respond` / `clierr` | Available, pending Resolve source check | Reuse; no new terminal envelope. |
| Stateless resolver kernel | RDR 0001 | Pending peer RDR | `flow resolve` delegates selection and refusal classes. |
| Transition model and outcome alphabet | RDR 0002 | Pending peer RDR | `flow next` depends on model-provided outcome and row data. |
| Guard predicate semantics | RDR 0003 | Pending peer RDR | CLI reports conditions but does not own guard truth. |
| Accessor execution and read-back | RDR 0004 | Pending peer RDR | `read-state` / `set-state` expose accessor outcomes. |
| Graph lint authority | RDR 0006 | Pending peer RDR | CLI may later expose lint, but this RDR does not own lint guarantees. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Output mode validation | `internal/cli/respond::ValidateMode` | Current verbs must call it manually. | Reuse | Every new verb starts with it. |
| Success rendering | `internal/cli/respond::OK` | Text mode currently emits only advisories, while `version` uses `cmd.Println`. | Extend or adjust usage | Resolve must settle text success rendering for verb payloads. |
| Failure rendering | `internal/cli/respond::Fail` and `internal/cli/clierr::CLIError` | No resolver-specific codes yet. | Extend | Add stable refusal codes as needed. |
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
refusal codes, and read-back-gated success normative; Resolve must verify those
points before lock rather than leaving them as implementation taste.

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
  **Mitigation**: It only persists supplied next tags through declared write
  accessors and depends on RDR 0004 read-back verification.

### Failure Modes

Visible failures include unknown command arguments, model/config not found,
unknown outcome, zero matching row, multiple matching rows, missing supplied
guard facts, accessor unavailable, gate indeterminate, and write read-back
mismatch. Each should surface as one `CLIError.Code` with a clear message,
detail, and hint where useful.

The main silent-failure risk is treating a failed accessor or ambiguous row as a
successful transition. The recovery path is refusal-first: the command exits
non-zero, emits the structured error envelope, and leaves skill judgment or
model repair to the caller. Diagnosis starts from the error code, model/rule id
when available from RDR 0002, and accessor identity/artifact role when the
failure comes from RDR 0004.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0001, 0002, 0003, and 0004 expose enough internal contracts for the
      CLI to call without owning their semantics.
- [ ] Resolve settles whether text success payloads should be added to
      `respond.OK` or emitted through a small response helper behind `respond`.

### Minimum Viable Validation

Implement one fixture-backed flow and prove all four verbs through the
production Cobra path: `flow next` returns the legal outcome alphabet,
`flow resolve` maps one outcome to one next tag-set, `flow read-state` reads the
fixture artifact tags, and `flow set-state` persists a planned owned-tag write
then read-back-verifies it. Run the same happy path and at least one typed
refusal in `--as=text` and `--as=json`.

### Phase 1: Command Contract Skeleton

Add the `flow` command group and the four verbs with Cobra argument validation,
`respond.ValidateMode`, `respond.OK`, `respond.Fail`, and focused command tests.

### Phase 2: Kernel And Model Binding

Wire `next` and `resolve` to the RDR 0001/0002/0003 interfaces using a fixture
model before broad config discovery.

### Phase 3: Accessor Binding

Wire `read-state` and `set-state` to RDR 0004's accessor executor and preserve
read-back failure as a typed CLI refusal.

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

[Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Do not include effort estimates or speculative
throughput targets. Rough performance metrics are
appropriate only when comparing alternatives — note
empirical data or obvious gains that support the
chosen approach over a rejected one.]

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

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

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

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

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
past the lens battery undetected. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

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
