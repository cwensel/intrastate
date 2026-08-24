# Recommendation 0005: Skill Integration CLI Contract

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft [revised from Final 2026-06-24; re-verify A3, A4, A5, A6; STAGE-SCOPED — JDR 0001 §D8–§D11 answered JD-8/9/19/20 and land here; §D13 by citation]
- **Type**: Feature
- **Profile**: mid — one user-facing CLI integration contract over resolver, accessor, and output seams.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel, 0002-transition-table-as-reviewable-data, 0003-guard-predicate-exhaustiveness, 0004-accessor-execution-safety-model, JDR 0001 (resolve-kernel seam)
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill author needs to ask "where next", "what outcomes are legal", and "read or persist state" in one deterministic CLI contract instead of re-implementing transition logic. The system-internal requirement is to define the CLI surface over the resolver kernel, including output shape, state binding, and ownership of conditional evaluation.

## Context

### Background

The CLI is the surface over the kernel. `next` must expose the legal-outcome alphabet for constrained decoding, while keeping the resolver semantics in the kernel rather than in skill code.

The real design fork is the integration contract: text versus JSON versus exit-code output, workspace marker versus config versus per-call args for state binding, and whether callers or the CLI evaluate conditional next edges.

### Technical Environment

intrastate is a Cobra-based Go CLI in `internal/cli`. Every verb must start with `respond.ValidateMode(cmd)`, route success through `respond.OK`, route failure through `respond.Fail`, and keep text/json behavior inside `internal/cli/respond`. The resolver kernel is implemented in `internal/resolve` (RDR 0001); no `flow` verbs exist yet.

## Research Findings

### Investigation

No resolver CLI, table loader, or accessor executor exists under `internal/cli`;
the implemented surface is the root command, `version`, `respond`, `clierr`, and
config loading. The kernel does exist: `internal/resolve::Resolve` takes an
`Input` with `Owned`, `Observed`, `Recognized`, and a `Guards` seam and returns
one plan or one refusal from the closed `internal/resolve::RefusalKinds` set
(`no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`,
`unmodeled_outcome`). The proposal is constrained by the project CLI contract
and by that kernel seam, not by a competing implemented resolver surface.

Prior art was read before naming approaches. `docs/cli-output-contract.md`
states that the persistent `--as text|json` flag selects the wire format, that
JSON stdout carries exactly one terminal record, and that exit-code mapping
lives in `clierr.ExitCodeFor`. `internal/cli/respond::OK`,
`internal/cli/respond::Fail`, `internal/cli/respond::ValidateMode`,
`internal/cli/clierr::CLIError`, and `internal/cli::ExecuteAndEmit` already
provide the failure and output gateway this RDR reuses rather than bypasses.
The current `internal/cli::newVersionCmd` text path still uses `cmd.Println`;
that is a source-level drift from the `respond.OK` convention, so new resolver
verbs must add text payload rendering inside `respond` instead of copying the
version shortcut.

The peer RDR split is load-bearing. RDR 0001 owns a stateless resolver kernel
that returns one transition plan or one typed refusal and explicitly delegates
the CLI mapping to this RDR. RDR 0002 owns the closed TOML transition-model
layout, the recognized-outcome alphabet, normalized rows, rule identity, and the
load-category taxonomy. RDR 0003 owns symbolic guard predicates and
exhaustiveness. RDR 0004 owns declared read/gate/write accessor capabilities and
read-back verification. RDR 0006 owns static graph lint authority. This RDR
exposes those capabilities through a deterministic CLI contract without
absorbing their data formats, predicate grammar, accessor safety model, or lint
invariant set.

Four decisions no single peer owns were hoisted to JDR 0001 and decided there;
this RDR cites them and never restates them: §D8 (which verb assembles owned
state and calls `Resolve`), §D9 (where a gate runs and what `deny` means),
§D10 (what the error envelope carries and which exit code means what), and
§D11 (how the CLI carries a clear, a set value, and an unbound write). §D13
fixes the byte form of a set-valued tag at the kernel seam, which this RDR
reuses as the CLI's set literal.

The strongest external prior is the sibling resolver design. It frames the CLI
as four verbs: `next`, `resolve`, `read-state`, and `set-state`; states that the
flow is "navigated, not orchestrated"; and says `next` emits the
legal-outcome alphabet while an optional constrained decoder remains
downstream. `../state-machines/MODEL-transition.md` sharpens the same split:
`next(state-tags)` returns legal outcome tags, `resolve(state-tags)` is pure and
location-free, and `set-state` persists only the already decided next tags with
read-back verification. Peer CLIs decide the shape of the rest: a command
handler owns load→decide with only *decide* pure (§D8); a second,
side-effecting check runs after selection on its own channel (§D9); exit codes
distinguish only repair-the-environment classes and error envelopes stay flat
plus one structured diagnostics list (§D10); "clear" is a key-only flag and a
set value is an explicit container literal (§D11).

### Key Discoveries

- **Documented** - `docs/cli-output-contract.md` makes `--as text|json`,
  JSON terminal records, advisory streams, and `clierr.ExitCodeFor` the existing
  user-facing output contract.
- **Documented** - `internal/cli/respond::OK`,
  `internal/cli/respond::Fail`, and `internal/cli/respond::ValidateMode` are
  the existing output gateway for new verbs.
- **Verified** - `internal/resolve::Resolve` is pure over an `Input` the caller
  assembles; the closed `RefusalKinds` set is the kernel half of this RDR's
  error taxonomy.
- **Documented** - RDR 0002 provides the legal recognized-outcome alphabet,
  normalized rows, `[read.<id>]` / `[write.<id>]` / `[gate.<id>]` capability
  declarations with `keys`, and the load-category taxonomy; its loader takes
  bytes plus a source id, so file I/O sits at this CLI seam.
- **Documented** - RDR 0004 owns read/write/gate accessor execution, gate
  result shape, and read-back verification; this RDR exposes those outcomes
  through verb-appropriate result shapes, not a redefined accessor safety
  model.
- **Documented** - JDR 0001 §D8–§D11 decide the verb/accessor binding, the gate
  site, the envelope carrier and exit mapping, and the write grammar; this RDR
  is the landing document for all four.
- **Documented** - sibling prior art names the thin CLI as
  `resolve`/`next`/`read-state`/`set-state` and explicitly rejects a runtime,
  driver, or framework.
- **Verified** - a single Cobra command group can expose all four verbs without
  making the CLI the owner of table parsing, guard semantics, or accessor
  execution.
- **Verified** - `next` can return a compact legal-outcome alphabet and
  candidate summary that is useful to skills in text mode and complete enough
  for tools in JSON mode.

### Critical Assumptions

- **A1 The existing output gateway can carry all resolver CLI success and
  failure dispositions without direct stdout/stderr writes.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/respond::Success`,
    `internal/cli/respond::OK`, `internal/cli/respond::Fail`,
    `internal/cli/respond::ValidateMode`, `internal/cli/clierr::CLIError`,
    and `internal/cli::ExecuteAndEmit` provide the shared success,
    advisory, validation, structured failure, and Cobra-error gateway. Source
    search also found `internal/cli::newVersionCmd` still writes text
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
    recognized outcome into a plan, `flow read-state` reads fixture artifact
    tags, and `flow set-state` persists planned owned-tag writes with
    read-back. The sibling prior `../state-machines/MODEL-transition.md` uses
    the same four operations and assigns no fifth required operation to the
    skill integration loop. `read-state` is diagnostic once `resolve` reads for
    itself (§D8) but remains the only verb that exposes a read without a
    resolution.
  - **If wrong**: Skill authors would still re-implement part of transition or
    state-binding logic outside intrastate.
- **A3 `next` can expose the legal-outcome alphabet without owning guard
  semantics or hiding unresolved facts.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Technical Design` defines the recognized outcome
    alphabet ("the closed set of outcome tags that the recognizer may emit for
    the flow") and normalized candidate rows carrying "the source identity
    (`(model id, rule id)` plus any expansion suffix, and the source locator)".
    The selection procedure — gate-then-count, exact-one survivor, which
    refusals are escapable — is JDR 0001 §D2's ("**Resolved: (b).** Gate, then
    count.") and the kernel's (`internal/resolve::RefusalKinds`);
    RDR 0002 `Normative Contracts` states the disclaimer itself — the procedure
    "is **JDR 0001 §D2's and the kernel's, and MUST NOT be restated here**" —
    so 0002 relies on it rather than defining it. RDR 0003 `Approach` makes
    every guard "a symbolic tag-set predicate that lint can reason about", owned
    outside this CLI; presence and unevaluability are decided in the kernel
    (JDR 0001 §D4, "**Resolved: (b).** The kernel enforces the guard domain").
    Therefore `next` can enumerate the alphabet, surface supplied and unresolved
    facts, and map kernel refusals without becoming the owner of guard truth.
    Re-verified at the re-entry Stage 4: the citation now anchors on §D2 / the
    kernel, and RDR 0002 `Normative Contracts` carries no competing selection
    procedure.
  - **If wrong**: The CLI would either under-inform constrained decoding or
    incorrectly become the owner of conditional evaluation.
- **A4 The resolver kernel stays pure while each verb owns load→decide: reads
  and gates bind in the verb over explicit artifact bindings; selection stays
  in `internal/resolve`.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: JDR 0001 §D8 resolves that "`flow resolve` and `flow next` MAY
    run declared read accessors over explicit `--artifact` bindings and MUST
    assemble `Input.Owned` from them", and that "the purity fence moves to where
    it belongs: `internal/resolve`, not the verb". §D9 resolves "**(c)**" — "Only
    the selected row's `gate` list runs", after exact-one selection and before
    the plan, with "**All gates on the selected row run; deny-overrides; every
    result reported.**" and "**Deny is not an escape class.**" RDR 0004
    `Normative Contracts` keeps the capabilities distinct over caller-supplied
    artifact roles ("Every accessor definition MUST declare exactly one
    capability: read, gate, or write"; "Accessors MUST operate on caller-supplied
    artifact roles"): reads return typed tags, gates return
    allow/deny/indeterminate, writes apply planned owned-tag mutations with
    same-role read-back that also verifies "observed and recognized tag values
    present before the write are unchanged". The kernel half is confirmed in
    source: `internal/resolve::Resolve` is pure (the package imports only
    `slices` and `strings`, enforced by
    `internal/resolve/resolve_test.go::TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility`)
    and never executes accessors to fill missing owned state
    (`::TestReq13_KernelDoesNotExecuteAccessorsToFillMissingOwnedState`).
    Therefore the CLI binds reads and gates in `next` / `resolve`, reads only in
    `read-state`, and writes only in `set-state`, without making the kernel
    stateful or coercing gate results into tag values. Re-verified at the
    re-entry Stage 4: the June fencing (`resolve` barred from read accessors,
    gates before the kernel call) is superseded by §D8/§D9. Implementation note:
    the kernel's unexported pre-selection guard filter also spells itself
    `gate`; it is a different mechanism from this RDR's post-selection accessor
    gates and must not be wired to them.
  - **If wrong**: Either every write-bearing row refuses `owned_state_unavailable`
    (caller-supplied context can never satisfy an owned dependency), or a
    stale caller-carried snapshot mints a plan `set-state` applies blindly.
- **A5 Stable CLI error codes plus one structured `findings` field can carry
  every input, model-load, kernel, gate, accessor, and write refusal class with
  no new exit group.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: The exit-mapping half is verified in source.
    `internal/cli/clierr::CLIError` carries a stable string `Code` plus optional
    `omitempty` `Param`, `Detail`, and `Hint`;
    `internal/cli/clierr::ErrorCode` exposes the code for branching;
    `internal/cli/clierr::ExitCodeFor` maps `GroupUserEnv` and `GroupInternal`
    to exit 2 (one case arm — the distinction is carried by `Code`, not by
    group), `GroupEnvUnavailable` to exit 3, `GroupSignalCancel` to 130, and a
    non-`CLIError` to 1. JDR 0001 §D10 assigns "**Exit 3 = the environment could
    not be consulted; repair it and re-run the same request unchanged.**" and
    "**Exit 2 = the request or the model is wrong, or the model said no.**",
    with "No new exit group (0005's A-block)", and adds exactly one `omitempty`
    structured field — `Findings []clierr.Finding` — as the carrier for
    per-row, per-key, per-atom, per-gate, and per-load-category discriminators.
    §D10 is the *sole* authority for the exit-3 rule: RDR 0004 mints the
    accessor refusal classes but its `Disposition Table` states "`exit code` is
    out of scope here — RDR 0005 owns the CLI mapping", so no second leg
    supports it. Re-verified: `Findings []Finding` does not exist yet in
    `clierr`, as this RDR states.
    **The `Finding` field set — settled at the re-entry Stage 4 author's
    round.** This RDR's five (`Code`, `Message`, `Param`, `Locator`, `Hint`)
    are the subset it populates, not the whole type: `clierr.Finding` is shared
    by five producers, and RDR 0006 — Final — normatively requires more of that
    same type. `0006::Normative Contracts`: "Every blocking finding MUST carry a
    stable code, model identity, severity, human-readable message, and the
    source rule/context id or source span … A finding attributed to one guard
    atom MUST carry that atom's `Key`, `Operator`, `Literal`, and `Block`; a
    finding scoped to an escape population MUST carry the failure class."
    **Resolved: one flat record, widened with `omitempty` fields; each producer
    populates the subset it owns.** The atom's four fields sit flat on the
    record, not in a nested `atom` object. Grounding: the shipped type's own doc
    comment already fixes the project convention — "Extend with new optional
    fields as needed — keep them `omitempty` so the envelope stays append-only
    and stable for tools" (`internal/cli/clierr::CLIError`) — and `0009::A9`
    closed the same additive question by source search, finding that no consumer
    or test reads an exact field set. Prior art agrees: `golangci-lint`'s
    `Issue` carries ~100 producers on one flat record and met this exact case
    with nolintlint by **promoting** its producer-specific fields flat rather
    than nesting them (`golangci-lint/pkg/result/issue.go:16`); `beads`'
    `Issue` does the same across ~30 banner-grouped fields
    (`beads/internal/types/types.go:19`); and no system examined nests a typed
    per-producer sub-object on a finding — Kubernetes nests typed `Details` on
    the *envelope*, one per response keyed by reason, and still flattens
    structured detail into the message for generic consumers
    (`k8s.io/apimachinery/.../api/errors/errors.go::NewInvalid`). Flattening the
    other way — forcing 0006's data into the five — is refused: `0006::Normative
    Contracts` sorts findings "by model id, invariant code, source rule/context
    id … then normalized predicate/write fingerprint", and under five fields
    `model` and `severity` have no field while `rule` fuses into `locator`, so
    the sort key must be re-parsed out of a string and collides whenever a rule
    id contains the delimiter (`evidence/spikes/finding-shape-options.out`).
    Carried obligation, from every prior-art system that kept structured detail:
    `Message` MUST stand alone without the structured fields
    (`evidence/spikes/finding-shape-tiebreak.out`). Note `0006::Technical
    Design`'s "`Key`, `Operator`, `Literal`, `Block`, and `Class` are `string`"
    agrees with this shape; 0006's illustrative JSON showing a nested
    `"atom": {…}` does not, and is the side that needs a citation repair — it is
    illustrative, not normative. The `omitempty` half was never open: §D10 item
    3 resolves it — the `CLIError` field is `omitempty` and never empty on a
    failure, while 0006's "empty list is the receipt" binds the success
    payload's own non-omitempty `data.findings`.
  - **If wrong**: Scripted skill calls could not branch deterministically on
    resolver failure classes, or model repair could not locate the defect.
- **A6 The pinned request grammar, minimum success payload fields, and
  stable `flow-*` code spellings are sufficient for the first fixture-backed
  CLI implementation.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: The grammar (`--flow <id>` | `--model <path>`, repeated
    `--tag name=value`, `--artifact role=path`, `--outcome <tag>`,
    `--evaluate-gates`, `--write name=value`, `--clear <key>`), the minimum
    JSON `data` fields in Technical Design, and the code table in Failure Modes
    are this RDR's implementation contract. The rejected alternative is leaving
    grammar, payload fields, or error-code spelling to implementation-time
    invention. Re-accepted at the re-entry Stage 4 against the grammar as
    pinned upstream: JDR 0001 §D11 fixes "**`--write k=v` ↔ the write block; a
    repeatable `--clear <key>` ↔ the rule-level `clear` list.**", the refusal of
    `--write k=<clear>` with the hint "use `--clear`", and "**A `set`-kind key's
    `--write` value is a JSON array literal**"; §D10 item 6 fixes model entry as
    "an explicit file path … beside the config-resolved `--flow <id>`; mutually
    exclusive; `flow-model-not-found` either way". Three spellings here are
    this RDR's own, beyond the JDR's non-normative table, and are recorded as
    RDR-owned rather than cited: `flow-write-duplicate`; the set-kind `--tag`
    array literal with `flow-tag-invalid` on kind mismatch; and
    `flow-model-not-found` covering the selection *arity* error as well as
    non-resolution.
    **Set-member string escaping — settled at the re-entry Stage 4 author's
    round.** JDR 0001 §D13 fixes "members sorted, duplicate-free, compact
    encoding" but is silent on JSON string escaping inside a member, and no RDR
    or JDR text supplied it. Because this RDR claims "copy-through and read-back
    equality are byte equality", the encoder is pinned here: the canonical set
    literal is rendered with HTML escaping **disabled**, so `<`, `>`, and `&`
    serialize as themselves. Grounding: RFC 8785 (JSON Canonicalization Scheme)
    §3.2.2.2 requires a Unicode value outside the ASCII control range to be
    serialized "as is" unless it is a backslash or a quote, which makes a
    `\uXXXX` spelling of `<`/`>`/`&` non-conformant rather than merely verbose;
    peer Go CLIs disable it on exactly this kind of path
    (`gh-cli/pkg/cmdutil/json_flags.go:228`, the single funnel for every
    `gh --json` output; `beads/cmd/bd/protocol/corpus.go:343::marshalCanonical`,
    "without HTML-escaping, so the bytes are stable and minimal"); and Go's own
    `SetEscapeHTML` doc names non-HTML readability as the reason to disable it,
    which is this case — no set literal is ever embedded in a `<script>` tag.
    Normative fixtures approved at the round, from spike
    `evidence/spikes/escaping-surfaces.out`: `["cli","final"]` (the common case,
    identical under either encoder) and `["a<b","p>q","x&y"]` (the divergent
    case). Spike `evidence/spikes/d13-canonical-set.out` carries the sorting,
    duplicate-free, compact, empty-set, and delimiter-collision legs. The
    divergence surface is exactly three characters: U+2028, U+2029, quote,
    backslash, tab, newline, NUL, non-ASCII, and emoji are byte-identical under
    both encoders, so disabling HTML escaping is minimal over ASCII but not
    strictly JCS-conformant — a recorded deviation, since JCS conformance is not
    claimed. Cost accepted: both shipped emit sites
    (`internal/cli/clierr::EmitJSON`, `internal/cli/respond::writeJSONLine`)
    move off bare `json.Marshal` onto one shared helper.
  - **If wrong**: Implementers would still need to invent request grammar,
    payload fields, or error-code spellings during code work.

### Reconciliation Report

A3–A6 re-verified at the re-entry Stage 4; the two items A5 and A6 opened on
were settled at that stage's author's round. This table is rewritten at the
re-entry Stage 6.

| Item | Source | Disposition | Evidence pointer or plan |
| --- | --- | --- | --- |
| A1 output gateway | 1 | VERIFIED | Source Search recorded in A1; anchors re-checked at the re-entry Stage 4 and all still resolve. |
| A2 four-verb minimality | 1 | VERIFIED | MVV Test recorded in A2. |
| A3 `next` alphabet without guard ownership | 4 | VERIFIED | Peer RDR recorded in A3 (JDR 0001 §D2/§D4, `0002::Normative Contracts` disclaimer, `0003::Approach`). |
| A4 verb owns load→decide, kernel stays pure | 4 | VERIFIED | Peer RDR recorded in A4 (JDR 0001 §D8/§D9, `0004::Normative Contracts`), kernel purity confirmed in source. |
| A5 error codes + one `findings` field | 4 | VERIFIED | Exit mapping verified in `internal/cli/clierr::ExitCodeFor`; the `clierr.Finding` field set settled at the re-entry Stage 4 round as one flat `omitempty` record covering all five producers. |
| A6 pinned grammar + code spellings | 4 | VERIFIED | Grammar re-accepted against JDR 0001 §D10 item 6 / §D11; set-member escaping pinned to HTML-escaping-disabled at the round, with normative fixtures from `evidence/spikes/escaping-surfaces.out`. |
| Reuse audit | 4 | NO FINDING | `internal/cli` registers only `newVersionCmd`; no TOML loader, accessor executor, `Finding` type, or `respond.OK` text-payload path exists. Greenfield; nothing to fold in. |

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
maps user flags and a model selection into typed calls, then renders results
through `respond`. It does not own transition-model representation, guard
semantics, accessor safety, graph lint, skill execution, or constrained
decoding. It is the deterministic integration contract a skill can call when it
needs the legal outcome alphabet, one resolved next transition, or a state I/O
operation.

Each verb owns load→decide; only *decide* is pure. `flow next` and
`flow resolve` run the model's declared read accessors over the caller's
explicit `--artifact` bindings, assemble the owned snapshot, and call the kernel
in the same invocation (JDR 0001 §D8). Caller-supplied `--tag` values are
observed context and never satisfy an owned dependency. Nothing discovers
artifacts ambiently; location comes only from explicit bindings.

`flow next` answers "what outcomes are legal from this state?" It returns the
recognized-outcome alphabet plus candidate summaries — source rule identity,
required and unresolved facts, preview next tags and write targets — read from
normalized model data without evaluating missing facts. It never calls a
language model. Gates run only on opt-in (`--evaluate-gates`); by default gate
ids are listed as unresolved facts so the verb stays cheap and effect-free
(§D9).

`flow resolve` answers "given this recognized outcome and this state, what
exact plan or refusal results?" It delegates selection to the kernel, runs the
selected row's gates after exact-one selection and before the plan (§D9), and
maps kernel and gate refusals to stable `CLIError.Code` values. A denied gate is
a refusal, not a plan and not an escape class. An escaped plan is a success
marked `escaped` on the payload (§D10).

`flow read-state` invokes declared read accessors only and returns the read
tag-set. It is the diagnostic verb — a read without a resolution. It does not
invoke gate accessors, because gates return allow, deny, or indeterminate
rather than tag values.

`flow set-state` persists a decided owned-tag mutation — `--write` values and
`--clear` keys mirroring the plan's write block and clear list (§D11) —
through RDR 0004's write accessors and returns success only after read-back
verification succeeds. It never runs gates; the read-back is the commit-time
check, and the window between `resolve` and `set-state` is a stated
non-guarantee of the two-verb design (§D9).

### Technical Design

The CLI layer is a translator, not a decision engine. Each verb starts with
`respond.ValidateMode(cmd)`, validates arguments into typed request values,
loads the selected transition model, calls the owning internal package, and
renders one terminal disposition. Success uses `respond.OK`; user-visible
failure uses `respond.Fail(cmd, &clierr.CLIError{...})`. The text mode contract
is human-scannable but still derived from the same request/response structs as
JSON mode.

**Model selection.** `--flow <id>` selects a config-resolved model;
`--model <path>` names a model file explicitly. They are mutually exclusive;
either failing to resolve is `flow-model-not-found`. The CLI reads the file and
hands RDR 0002's loader bytes plus the path as source id — the loader performs
no file I/O. Every load category RDR 0002 defines, RDR 0003's two predicate
categories, and RDR 0008's `reserved_tag_key` map to one code,
`flow-model-invalid`, with one `findings[]` entry per category hit
(`code` = category slug, `locator` = file:line). The remedy for all of them is
"fix the model at this locator", one branch.

**State input.** `--tag name=value` supplies one observed tag. A `--tag` naming
an owned key is refused `flow-tag-owned`; one naming `recognized` is refused
`flow-tag-reserved`; both before any accessor runs (§D8 — refuse rather than
silently shadow under owned-over-observed precedence). A duplicate tag name is
`flow-tag-duplicate`. A set-kind tag's value is a JSON array literal
(`'labels=["a","b"]'`), the same canonical form the kernel seam carries (JDR
0001 §D13) and `--write` accepts — sorted, duplicate-free, compact, and
rendered with HTML escaping disabled; a bare scalar for a set key, or an array for
a scalar key, is `flow-tag-invalid`. Owned state is never caller-supplied: it
is assembled from declared read accessors over `--artifact role=path`
bindings, where `role` is the model-declared artifact role and `path` the
caller-owned artifact. A role an invoked accessor needs that no binding
supplies is `flow-artifact-missing`.

**Outcome.** `flow resolve` requires `--outcome <tag>`; an absent or empty
value is `flow-tag-invalid` at the CLI, and a value outside the model's
alphabet is the kernel's `flow-unmodeled-outcome`.

**Gates.** `flow resolve` runs every gate on the selected row, after
exact-one selection. Deny-overrides; indeterminate does not override deny; every
gate's result is reported, as `gates[]` on a success or as one `findings[]`
entry per gate on `flow-gate-denied` / `flow-gate-indeterminate`. `flow next`
runs gates only under `--evaluate-gates`. A gate's timeout or execution failure
is an accessor refusal (exit 3), not a gate result.

**Writes.** `flow set-state` takes `--write name=value` (scalar, or a JSON
array literal for a set-kind key) and a repeatable `--clear <key>`. The
sentinel `<clear>` is unauthorable: `--write k=<clear>` is refused
`flow-write-invalid` with the hint "use `--clear`". The same key under both
flags, or twice under either, is `flow-write-duplicate`. Before any accessor
runs the CLI checks that each `--write` / `--clear` key is a declared owned tag
served by exactly one `[write.<id>]` whose `keys` names it, and that each value
is well-formed for its declared kind — otherwise `flow-write-unbound`,
`flow-clear-unbound`, or `flow-write-invalid`. `[]` (empty set) and `--clear`
(absent) stay distinct through read-back. Each writer applies its own keys and
reads back; no cross-writer atomicity is promised. `--tag` on `set-state` is
context only and is never written.

`set-state` trusts the caller to transcribe the plan; nothing links a
`set-state` request to a prior `resolve`. Because `resolve` renders a set write
as the same JSON array `set-state` accepts, copy-through from plan to request is
byte-identical. A carried plan artifact (`--plan <file|->`) is deferred as a
seed, not part of this contract.

**Envelope.** Successful JSON payloads use the existing `{"type":"ok","data":...}`
envelope. The minimum `data` shape is:

- `next`: `model` (id or path) and `revision`, `observed` tags, `owned` tags
  assembled from the readers, `readers[]` (accessor identities invoked),
  `outcomes[]` (recognized outcome tags — nothing else; RDR 0002's layout has
  no per-outcome recognizer text and `[model.metadata]` is not its carrier),
  and `candidates[]`. Each candidate carries source rule identity, the outcome
  it belongs to, required facts, unresolved guard/gate facts, evaluated gate
  results when `--evaluate-gates` was given, and preview next tags, write
  targets, and clear keys when the normalized model exposes them without
  evaluating missing facts.
- `resolve`: `model`, `revision`, `observed`, `owned`, `readers[]`,
  `outcome`, `rule` (matched rule identity), `gates[]` (`id`, `result`),
  `next` (the next tag-set), `writes` (planned owned-tag writes, set values as
  JSON arrays), `clear[]` (planned clears), `escaped` (bool) and
  `escape_class` when the plan came from an escape row.
- `read-state`: `model`, `revision`, artifact role bindings, `readers[]` each
  with its declared `keys` (the requested key set) and the tags it returned —
  so "absent from the artifact" and "not requested" are distinguishable from
  the payload alone.
- `set-state`: `model`, `revision`, artifact role bindings, `writers[]`,
  requested `writes` and `clear[]`, and the read-back-confirmed owned-tag
  values.

Failure payloads use the existing `CLIError` envelope plus exactly one new
`omitempty` structured field, `findings` — a list of `Finding`,
subsystem-agnostic and shared with RDR 0006's lint findings (§D10). `Finding`
is one **flat** record whose optional fields are `omitempty`; each producer
populates the subset it owns. This RDR populates `code`, `message`, `param`,
`locator`, and `hint`; RDR 0006 additionally populates `severity`, `model`,
`rule`, the guard atom's `key` / `operator` / `literal` / `block`, and `class`.
The atom's four fields sit flat on the record, not in a nested `atom` object.
`message` must stand alone: it renders the failure readably without any
structured field being consulted. Exit codes: **3** means the environment
could not be consulted (accessor timeout, execution failure, incomplete read,
read-back incomplete, post-mutation timeout) — repair it and re-run the same
request unchanged; **2** means the request or the model is wrong, or the model
said no. No new exit group.

Text mode renders the same result content in a human-scannable order; it does
not invent fields absent from the JSON payload, and it renders `findings`
one-per-line under the message.

#### Normative Contracts

```normative
The CLI MUST expose one command group for skill integration with these verbs:
next, resolve, read-state, and set-state. Other command groups (lint, dump,
parse) are outside this contract and are owned by the RDR that names them.

All four verbs MUST start RunE by calling respond.ValidateMode(cmd), MUST route
success through respond.OK, MUST route user-facing failure through
respond.Fail(cmd, *clierr.CLIError), and MUST set SilenceErrors and
SilenceUsage.

Under --as=json, each successful invocation MUST emit exactly one stdout JSON
terminal envelope with type "ok" and verb-specific data. Under --as=text, each
successful invocation MUST emit human output derived from the same verb-specific
result. Failures MUST use the existing CLIError JSON/text envelope defined by
docs/cli-output-contract.md and internal/cli/clierr, extended by exactly one
omitempty structured field, findings, per JDR 0001 §D10. Finding MUST be one
flat record with omitempty optional fields, carrying at least code, message,
param, locator, hint, severity, model, rule, key, operator, literal, block, and
class; each producer populates only the fields it owns, and no producer nests
its own fields in a sub-object. Finding.message MUST be self-sufficient — it
MUST render the failure readably with no structured field consulted.

Exit 3 MUST mean the environment could not be consulted and the same request
may be re-run unchanged; every other failure MUST exit 2. No new exit group.

Model selection MUST accept exactly one of --flow <id> or --model <path>. The
CLI MUST perform the model file I/O and hand RDR 0002's loader bytes plus a
source id. Every load-time category MUST map to flow-model-invalid with one
findings[] entry per category hit carrying its locator.

flow next and flow resolve MAY run declared read accessors over explicit
--artifact role=path bindings and MUST assemble owned state only from them
(JDR 0001 §D8). They MUST NOT discover artifacts and MUST NOT run write
accessors. --tag values MUST enter as observed context. A --tag naming an owned
key or the reserved recognized key MUST be refused at the CLI before any
accessor runs.

flow next MUST return the legal recognized-outcome alphabet for the supplied
state, plus candidate summaries containing source rule identity, required
facts, unresolved guard/gate facts, and preview next tags, write targets, and
clear keys when those can be read from normalized model data without
evaluating missing facts. It MUST run gate accessors only when --evaluate-gates
is given; otherwise it MUST list gate ids as unresolved facts. It MUST NOT
invent guard facts that were neither supplied, read, nor produced by a declared
gate accessor.

flow resolve MUST return exactly one plan or exactly one CLIError refusal.
Kernel refusals MUST map one-to-one onto flow-unmodeled-outcome,
flow-no-match, flow-ambiguous-match, flow-owned-state-unavailable, and
flow-guard-unevaluable. It MUST run the selected row's gates only after
exact-one selection and before emitting the plan (JDR 0001 §D9); all gates on
the row MUST run, deny MUST override allow and indeterminate, and every gate
result MUST be reported. A denied or indeterminate gate MUST surface as
flow-gate-denied or flow-gate-indeterminate with one findings[] entry per gate;
it is never a plan and never an escape class. An escaped plan MUST be a success
carrying escaped=true and escape_class. flow resolve MUST NOT print directly,
initiate skill work, or choose among multiple matching rows.

flow read-state MUST invoke only declared read accessors over caller-supplied
role=path artifact bindings and return, per reader, its declared keys and the
tag-set read, or a stable CLIError failure. It MUST NOT invoke gate accessors
or coerce gate allow, deny, or indeterminate results into tag values.

flow set-state MUST invoke only declared write accessors over caller-supplied
role=path artifact bindings and the planned owned-tag mutations given as
--write name=value and --clear <key> (JDR 0001 §D11). Before any accessor runs
it MUST refuse a --write or --clear key that is not a declared owned tag served
by exactly one [write.<id>] whose keys list names it, a --write value malformed
for its declared kind, the literal value <clear>, and a key given twice. It
MUST never run gate accessors. It MUST report success only after the accessor
layer's read-back verification confirms the planned owned-tag values and that
non-owned tags are unchanged (RDR 0004). It MUST NOT treat --tag values as
writes.

A set-valued tag crossing the CLI in --tag, --write, or any payload MUST use
the canonical JSON array form JDR 0001 §D13 fixes — members sorted,
duplicate-free, compact — rendered with HTML escaping DISABLED, so <, >, and &
serialize as themselves and never as \u003c, \u003e, or \u0026. Every site
that emits or compares a canonical set literal MUST use the same encoder, so
plan-to-request copy-through and read-back equality are byte equality.
```

#### Load-Bearing Decisions

- **Identity** - a CLI request is identified by verb, model selection
  (`--flow` id or `--model` path) and model revision, observed tags, artifact
  role bindings, recognized outcome when applicable, and planned writes/clears
  when applicable. The same request over the same model revision and the same
  artifact contents must produce the same success or refusal, excluding
  exit-3 environment failures.
- **Wire / byte format** - JSON and text rendering conform to
  `docs/cli-output-contract.md`; this RDR introduces verb-specific JSON `data`
  payloads and exactly one new `omitempty` `CLIError` field (`findings`), not a
  new terminal envelope or exit group. Set values are JDR 0001 §D13's
  canonical JSON array everywhere they cross the CLI, rendered with HTML
  escaping disabled so `<`, `>`, and `&` serialize as themselves.
- **Naming** - the user-facing group is `flow`, with verbs `next`, `resolve`,
  `read-state`, and `set-state`; flags `--flow`, `--model`, `--tag`,
  `--artifact`, `--outcome`, `--evaluate-gates`, `--write`, `--clear`. Rejected
  group names: `state` because only two verbs perform state I/O, and `run`
  because the CLI does not orchestrate skill work. Kernel-mapped codes mirror
  the kernel refusal kinds one-to-one (`flow-<kind>`, underscores to hyphens).
- **Selection / predicate** - `flow resolve` succeeds only when the kernel
  reports exactly one matching row (or exactly one escape row for an escapable
  class) and every gate on that row allows. Zero matches, multiple matches,
  unmodeled outcomes, unavailable owned state, unevaluable guards, and gate
  deny/indeterminate become typed refusals rather than tie-breaks.

#### Round-Trip / Inverse Invariants

`set-state` followed by `read-state` must return the planned owned-tag values
for the artifact role that was written. The equality is value-for-value over
the owned tags the write planned to mutate — byte equality for set values in
canonical form — not byte-identical artifact content; a cleared key reads back
absent, an empty set reads back `[]`. RDR 0004 owns the read-back semantics,
including that non-owned tags present before the write are unchanged; this RDR
owns surfacing the outcome through the CLI contract.

`resolve` → skill → `set-state` is copy-through: every `writes` value and
`clear[]` key in a plan is accepted verbatim by `set-state`'s `--write` /
`--clear` grammar, so a set write survives the round trip byte-identical. This
holds because one encoder renders every canonical set literal; two emit sites
disagreeing on HTML escaping would report `flow-write-readback-mismatch` on a
write that actually succeeded, and only for members containing `<`, `>`, or `&`
— a data-dependent failure ASCII-clean fixtures never catch
(`evidence/spikes/escaping-surfaces.out`).

`read-state` → `--tag` re-pairs only for observed tags; an owned tag read by
`read-state` cannot be handed back through `--tag` (refused `flow-tag-owned`)
because `resolve` reads it itself.

#### Illustrative Code

Illustrative invocation shapes:

```sh
intrastate flow next --model ./flows/rdr.toml --artifact rdr=docs/rdr/0005-skill-integration-cli-contract.md --tag profile=mid --as=json
intrastate flow resolve --model ./flows/rdr.toml --artifact rdr=docs/rdr/0005-skill-integration-cli-contract.md --tag profile=mid --outcome round-clean
intrastate flow read-state --flow rdr --artifact rdr=docs/rdr/0005-skill-integration-cli-contract.md
intrastate flow set-state --flow rdr --artifact rdr=docs/rdr/0005-skill-integration-cli-contract.md --write status=Final --write 'labels=["cli","final"]' --clear prelock_lens
```

Non-normative payload and refusal shapes are the examples under JDR 0001 §D8,
§D9, and §D11.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Text/json output gateway | Existing `respond` / `clierr` | Verified | Reuse; one new `omitempty` `findings` field, no new terminal envelope or exit group. Text success payloads route through `respond.OK` or a respond-owned helper. |
| Stateless resolver kernel | RDR 0001 (`internal/resolve`) | Implemented | `flow resolve` delegates selection; refusal kinds map one-to-one to codes. |
| Verb/accessor binding, gate site, envelope carrier, write grammar | JDR 0001 §D8–§D11, §D13 | Decided | This RDR is the landing document; cited, not restated. |
| Transition model, outcome alphabet, capability declarations, load categories | RDR 0002 | Final peer RDR | `next` / `resolve` depend on model rows; `set-state` reads `[write.<id>].keys`; load categories ride `flow-model-invalid` findings. |
| Guard predicate semantics | RDR 0003 | Final peer RDR | CLI reports conditions but does not own guard truth; set-literal spelling reused. |
| Accessor execution and read-back | RDR 0004 | Final peer RDR | Reads assemble owned state; gates surface per §D9; writes and read-back through `set-state`. |
| Graph lint authority | RDR 0006 | Out of scope for this contract | Shares the `Finding` record; lint verbs are 0006's. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Output mode validation | `internal/cli/respond::ValidateMode` | Current verbs must call it manually. | Reuse | Every new verb starts with it. |
| Success rendering | `internal/cli/respond::OK` | Text mode currently emits only advisories, while `version` uses `cmd.Println`. | Extend `respond.OK` or a respond-owned helper | New verbs must not use direct Cobra printing for text payloads. |
| Failure rendering | `internal/cli/respond::Fail` and `internal/cli/clierr::CLIError` | No resolver-specific codes; no structured diagnostics field. | Extend codes; add `Findings []Finding` `omitempty` (one flat record, shared with RDR 0006) | Append-only; no new exit groups. |
| Canonical JSON rendering | `internal/cli/clierr::EmitJSON` and `internal/cli/respond::writeJSONLine` | Both call bare `json.Marshal`, which HTML-escapes `<`, `>`, `&` | **Extend** — route both through one helper with `SetEscapeHTML(false)` | Set-literal bytes become stable across emit sites; touches RDR 0006's envelope path. |
| Root wiring | `internal/cli::NewRootCmd` / `ExecuteAndEmit` | Only `version` is registered today. | Extend | Register `flow` group and keep Cobra errors structured. |
| Kernel call | `internal/resolve::Resolve` / `Input` | Takes assembled `Owned`/`Observed`; performs no I/O. | Reuse | Verbs assemble `Input` from readers and `--tag`. |
| Config discovery | `internal/cli/config` | Parser placeholder; no transition-model config yet. | Extend later | `--model <path>` ships first; `--flow <id>` config lookup may follow. |

### Decision Rationale

Choose the four-verb thin CLI because it matches the user outcome without
moving ownership out of the peer seams. Skills get one deterministic contract
for legal outcomes, resolution, and state I/O. The kernel remains pure,
accessors remain the only state-binding authority, and all user-visible output
still conforms to the existing CLI envelope. The sibling prior art also points
to this exact split: a table plus four verbs gives the wall-clock win of
removing repeated position re-derivation without turning intrastate into a
workflow runner.

`resolve` reads for itself because the alternative — caller-carried owned
state — is sound only with a compare-and-set precondition on `set-state`, and
without one a stale snapshot mints a plan that is applied blindly (write skew);
the handler owning load→decide with only decide pure is the standard shape
(JDR 0001 §D8). Gates run after selection because a guard is pure and belongs
to selection while a gate is an effect with a timeout; running effects before
knowing the survivor wastes calls, and pruning a deny into `no_match` launders
it into an escapable class (§D9). One code per caller-branchable remedy plus a
structured `findings` list, and exit 3 reserved for repair-the-environment,
follows every peer CLI examined (§D10). `--clear` as a key-only flag and a JSON
array for set values make the `<clear>` sentinel unauthorable and give `[]` and
absent distinct spellings (§D11).

Rejected alternatives fail one of those boundaries. A JSON-only API is simpler
for scripts but violates the CLI's established text/json mode contract. A single
`resolve` mega-command hides state I/O and makes conditional alphabet queries
harder for constrained recognition. A higher-level `run` or `advance` command
would initiate workflow action and collapse the project distinction between
recognizing an outcome, resolving it, and performing the next skill.

Premortem: this ships and fails because the four verbs are present but
semantically thin: `next` omits enough condition detail to constrain a skill,
`resolve` maps all refusals to generic errors, and `set-state` looks successful
when read-back did not prove the owned tags changed. The contract survives that
failure mode by making candidate summaries, one-to-one refusal codes with
findings, and read-back-gated success normative rather than implementation
taste.

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
- Collapses the plan/apply split that lets skill judgment sit between
  resolution and persistence.

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
- **Caller-carried owned state (`--owned k=v` or piped `read-state`)**: needs
  an `If-Match`-style precondition on `set-state` to be safe; a larger design
  than the one avoided (JDR 0001 §D8).
- **A `flow step` that resolves and writes in one call**: collapses the
  plan/apply split; addable later as sugar over `resolve` + `set-state` (§D8).
- **One `flow-*` code per model-load category or per gate**: the remedy is the
  same branch; discriminators ride `findings` instead (§D10).
- **`--write key-` / `key=null` for clear**: one keystroke from a wrong write,
  or a value that collides with the sentinel (§D11).

## Trade-offs

### Consequences

- Skills get one deterministic surface for the closed-world transition question
  instead of copying transition logic into prompts.
- The CLI remains mode-aware and scriptable through the existing text/json
  contract; `findings` is the one envelope addition, shared with lint.
- `flow resolve` performs reads, so it is no longer effect-free; `flow next`
  stays effect-free by default and gates only on opt-in.
- `set-state` trusts the caller's transcription of the plan and re-checks
  nothing from it; the read-back is the only commit-time check. This is a
  stated non-guarantee, and a carried-plan flag is a later seed.
- Text-mode summaries for conditional `next` output and per-gate findings need
  careful design so they stay readable without becoming a second grammar.

### Risks and Mitigations

- **Risk**: `flow next` becomes a second guard evaluator by trying to resolve
  conditionals without facts.
  **Mitigation**: The contract requires candidate summaries read from
  normalized data and forbids evaluating missing facts; gates run only on
  opt-in.
- **Risk**: Resolver refusals collapse into generic `command-error`, or a gate
  deny is laundered into `no_match` and escaped.
  **Mitigation**: Kernel codes mirror the kernel kinds one-to-one; gates run
  after selection and deny is its own refusal; both covered in the MVV.
- **Risk**: A stale owned snapshot mints a plan that `set-state` applies.
  **Mitigation**: `resolve` reads owned state itself in the same invocation;
  the residual `resolve`→`set-state` window is accepted and named.
- **Risk**: Text and JSON outputs diverge semantically.
  **Mitigation**: Derive both renderings from one typed result per verb and
  test both modes.
- **Risk**: `set-state` is mistaken for an orchestration command.
  **Mitigation**: It only persists planned `--write` / `--clear` mutations
  through declared write accessors, never gates, and depends on RDR 0004
  read-back verification.
- **Risk**: A caller writes an owned key no writer serves, or a set value in a
  private spelling.
  **Mitigation**: Unbound keys and malformed values are refused at the CLI
  before any accessor runs; set values have one canonical form.

### Failure Modes

Visible failures include bad arguments, model not found or invalid, the kernel's
five refusal kinds, gate denied or indeterminate, accessor environment
failures, and write read-back failures. Each surfaces as one `CLIError.Code`
with a clear message and hint, `param` for single-subject failures, and
`findings[]` where several rows, keys, atoms, gates, or load categories must be
named.

Stable code strings (the spellings are normative; the group fixes the exit):

| Family | Failure | Code | Group / exit | Carrier |
| --- | --- | --- | --- | --- |
| input | malformed, empty, or wrong-kind `--tag` / `--outcome` value | `flow-tag-invalid` | `GroupUserEnv` / 2 | `param` |
| input | duplicate `--tag` name | `flow-tag-duplicate` | `GroupUserEnv` / 2 | `param` |
| input | `--tag` names the reserved `recognized` key | `flow-tag-reserved` | `GroupUserEnv` / 2 | `param` |
| input | `--tag` names an owned key | `flow-tag-owned` | `GroupUserEnv` / 2 | `param` |
| input | malformed or wrong-kind `--write` value, or the literal `<clear>` | `flow-write-invalid` | `GroupUserEnv` / 2 | `param` |
| input | `--write` / `--clear` key given twice across either flag | `flow-write-duplicate` | `GroupUserEnv` / 2 | `param` |
| input | `--write` key served by no single `[write.<id>].keys` | `flow-write-unbound` | `GroupUserEnv` / 2 | `param` |
| input | `--clear` key served by no single `[write.<id>].keys` | `flow-clear-unbound` | `GroupUserEnv` / 2 | `param` |
| input | malformed `--artifact` binding | `flow-artifact-invalid` | `GroupUserEnv` / 2 | `param` |
| input | a role an invoked accessor needs has no binding | `flow-artifact-missing` | `GroupUserEnv` / 2 | `param` = role |
| model | neither / both of `--flow`, `--model`, or the selection does not resolve | `flow-model-not-found` | `GroupUserEnv` / 2 | `param` |
| model | any RDR 0002 load category, RDR 0003 predicate category, or RDR 0008 `reserved_tag_key` | `flow-model-invalid` | `GroupUserEnv` / 2 | `findings[]` (category, locator) |
| kernel | `unmodeled_outcome` | `flow-unmodeled-outcome` | `GroupUserEnv` / 2 | `param` |
| kernel | `no_match` | `flow-no-match` | `GroupUserEnv` / 2 | `findings[]` (rows considered) |
| kernel | `ambiguous_match` | `flow-ambiguous-match` | `GroupUserEnv` / 2 | `findings[]` (matching rows) |
| kernel | `owned_state_unavailable` | `flow-owned-state-unavailable` | `GroupUserEnv` / 2 | `findings[]` (keys) |
| kernel | `guard_unevaluable` | `flow-guard-unevaluable` | `GroupUserEnv` / 2 | `findings[]` (rows; atoms when the kernel names them) |
| gate | a selected-row gate denied | `flow-gate-denied` | `GroupUserEnv` / 2 | `findings[]` one per gate |
| gate | no gate denied and at least one indeterminate | `flow-gate-indeterminate` | `GroupUserEnv` / 2 | `findings[]` one per gate |
| accessor | accessor timed out | `flow-accessor-timeout` | `GroupEnvUnavailable` / 3 | `param` = accessor id |
| accessor | accessor execution failed | `flow-accessor-failed` | `GroupEnvUnavailable` / 3 | `param` = accessor id |
| accessor | read returned an incomplete key set | `flow-read-incomplete` | `GroupEnvUnavailable` / 3 | `param` = accessor id |
| accessor | unknown accessor, capability mismatch, write to a non-owned tag (unreachable past load and CLI checks) | `flow-accessor-unknown`, `flow-accessor-capability-mismatch`, `flow-write-non-owned` | `GroupInternal` / 2 | — |
| write | read-back disagrees with the plan, or a non-owned tag changed | `flow-write-readback-mismatch` | `GroupUserEnv` / 2 | `findings[]` per key |
| write | read-back returned an incomplete key set | `flow-write-readback-incomplete` | `GroupEnvUnavailable` / 3 | `detail`: may have been applied |
| write | post-mutation read-back timed out | `flow-write-readback-timeout` | `GroupEnvUnavailable` / 3 | `detail`: may have been applied |

An escaped plan is not a failure: `flow resolve` exits 0 with `escaped: true`
and `escape_class` on the payload.

The main silent-failure risk is treating a failed accessor, a denied gate, or
an ambiguous row as a successful transition. The recovery path is
refusal-first: the command exits non-zero, emits the structured error envelope,
and leaves skill judgment or model repair to the caller. Exit 3 tells a script
to repair the environment and re-run the same request unchanged; exit 2 tells
it someone must look. Diagnosis starts from the code, then `findings[]`:
rule ids and locators from RDR 0002, gate ids, and accessor identity/artifact
role from RDR 0004.

## Implementation Plan

### Prerequisites

- [x] RDR 0001 is implemented (`internal/resolve`); RDR 0002, 0003, and 0004
      are Final and expose enough contract for the CLI to call without owning
      their semantics.
- [x] A3 and A4 re-verified at the re-entry Stage 4 against JDR 0001 §D2/§D4/
      §D8/§D9 and the code as it stands.
- [x] A5's `clierr.Finding` field set and A6's set-member escaping encoder
      settled at the re-entry Stage 4 author's round: one flat `Finding` record
      with `omitempty` fields, and HTML escaping disabled on the canonical set
      literal.
- [ ] Add text success payload rendering through `respond.OK` or a
      respond-owned helper used by `respond.OK`; do not print directly from
      resolver verbs.
- [ ] Add `Findings []Finding` (`json:"findings,omitempty"`) and the flat
      `Finding` record to `internal/cli/clierr`, shared with RDR 0006 — fields
      `Code`, `Message`, plus `omitempty` `Param`, `Locator`, `Hint`,
      `Severity`, `Model`, `Rule`, `Key`, `Operator`, `Literal`, `Block`,
      `Class`.
- [ ] Route `internal/cli/clierr::EmitJSON` and
      `internal/cli/respond::writeJSONLine` through one shared helper using
      `SetEscapeHTML(false)`, so canonical set literals render `<`, `>`, `&` as
      themselves and every emit site agrees byte-for-byte. Bare `json.Marshal`
      must not remain on a path a set value crosses.

### Minimum Viable Validation

Implement one fixture-backed flow and prove all four verbs through the
production Cobra path: `flow next` returns the legal outcome alphabet with
gate ids as unresolved facts; `flow resolve` reads owned state from a fixture
artifact, maps one outcome to one plan, and runs one gate on the selected row;
`flow read-state` reads the fixture artifact tags per reader; and
`flow set-state` persists one scalar write, one set write as a JSON array, and
one `--clear`, then read-back-verifies them. The set write's members MUST
include one containing `<` and one containing `&`, asserted byte-identical
through plan → request → read-back, so the encoder rule is covered by the
validation rather than only by review. Run the happy path, one escaped
plan, one gate deny, one kernel refusal, and one exit-3 accessor failure in
`--as=text` and `--as=json`, asserting `findings[]` where the table names it.

### Phase 1: Command Contract Skeleton

Add the `flow` command group and the four verbs with Cobra argument validation,
`respond.ValidateMode`, `respond.OK`, `respond.Fail`, the `findings` field, and
focused command tests for the input-family refusals.

### Phase 2: Kernel And Model Binding

Wire `next` and `resolve` to `internal/resolve` and the RDR 0002 loader over
`--model <path>` using a fixture model; map the five kernel kinds and the
load categories.

### Phase 3: Accessor Binding

Wire the verbs to RDR 0004's accessor executor only for the capability each
may invoke: reads in `next` / `resolve` / `read-state`, gates after selection
in `resolve` (and `next` under `--evaluate-gates`), writes and read-back in
`set-state`. Preserve gate deny/indeterminate, exit-3 accessor failures, and
read-back failures as typed refusals from the verb path that invoked them.

### Phase 4: Error Taxonomy And Docs

Land the full code table, update `docs/cli-output-contract.md` for the
`findings` field and the exit-3 rule, and document illustrative invocations.

### Day 2 Operations

This RDR does not create a persistent store. It exposes existing transition
model files and caller-owned artifacts through CLI operations.

| Resource | List | Info | Delete | Verify | Backup |
| --- | --- | --- | --- | --- | --- |
| Transition model files | Covered by repository tools | Covered by repository tools | Covered by repository tools | In scope through RDR 0006 lint and CLI MVV | Covered by version control |
| Caller artifacts mutated by `set-state` | Caller-owned | Caller-owned | Caller-owned | In scope through read-back verification | Caller-owned |

### New Dependencies

No new third-party dependency. Cobra, `respond`, and `clierr` already exist;
JSON array literals parse with the standard library. Any TOML dependency is RDR
0002's.

## Validation

### Testing Strategy

Command tests exercise the `flow` command group through the production Cobra
path rather than calling renderers or kernel functions directly. Coverage must
include argument validation, output mode validation, success rendering, and
typed refusal mapping for each verb.

1. **Scenario**: `flow next` over the fixture model in `--as=json` and
   `--as=text`, with and without `--evaluate-gates`.
   **Expected**: both modes report the same recognized-outcome alphabet and
   candidate summaries through the standard success path; JSON includes
   `outcomes[]` and `candidates[]`; without the flag no gate accessor runs and
   gate ids appear as unresolved facts.
2. **Scenario**: `flow resolve` over the fixture model with one recognized
   outcome that matches exactly one row, whose owned dependency is served by a
   fixture reader and which carries one gate.
   **Expected**: the reader runs, `owned` is assembled from it, the plan carries
   `rule`, `next`, `writes`, `clear[]`, and `gates[]` with `allow`; no write
   accessor runs; a `--tag` on the owned key is refused `flow-tag-owned`
   before the reader runs.
3. **Scenario**: `flow resolve` with an unmodeled outcome, an empty `--outcome`,
   a zero-match row, a multi-match row, and an escape row for `no_match`.
   **Expected**: each refusal maps to its one-to-one code with non-zero exit
   under both modes; the escape case exits 0 with `escaped: true` and
   `escape_class`.
4. **Scenario**: `flow resolve` where the selected row's gates return deny +
   allow, and separately indeterminate + allow.
   **Expected**: `flow-gate-denied` then `flow-gate-indeterminate`, each with
   one `findings[]` entry per gate; deny overrides indeterminate when both
   occur; no plan is emitted.
5. **Scenario**: `flow read-state` and `flow set-state` over a fixture artifact
   and declared accessor roles, with a scalar write, a set write as a JSON
   array, and a `--clear`.
   **Expected**: `--artifact` bindings are validated; `read-state` invokes only
   declared readers and reports each reader's `keys` beside its tags;
   `set-state` reports success only after read-back proves the scalar, the
   set (byte-equal canonical array), and the cleared key absent; `--write
   k=<clear>`, an unbound key, and a wrong-kind value are refused before any
   accessor runs. One set member carries `<` and one carries `&`: both render
   as themselves at every emit site — `["a<b","x&y"]`, never
   `["a\u003cb","x\u0026y"]` — and plan, request, and read-back agree
   byte-for-byte (normative fixture, `evidence/spikes/escaping-surfaces.out`).
6. **Scenario**: read accessor timeout during `resolve`, gate execution
   failure, write read-back incomplete, and write read-back mismatch.
   **Expected**: the first three exit 3 with `GroupEnvUnavailable` and a
   "may have been applied" detail on the write case; the mismatch exits 2 with
   `findings[]` per key; none is a successful transition or a coerced
   `read-state` payload.
7. **Scenario**: a refusal carrying `findings[]` is rendered in both modes, and
   the same `Finding` record is populated by a kernel refusal, a gate result,
   and a model-load category.
   **Expected**: each producer populates only its own fields and unset optional
   fields are absent from the JSON (`omitempty`); no producer nests its fields
   in a sub-object; and each finding's `message` alone renders the failure
   readably with no structured field consulted.

### Performance Expectations

No throughput target is load-bearing for this RDR. The command path should stay
single-invocation deterministic: parse inputs, load the selected model, run the
declared readers, call one kernel operation, run the selected row's gates, and
render one terminal result. Any accessor latency belongs to RDR 0004's
execution contract; this RDR only requires that an unavailable accessor
surfaces as an exit-3 refusal and an indeterminate gate as an exit-2 refusal.

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

Rewritten at the re-entry Stage 7 re-lock.

### Assumption Verification

Rewritten at the re-entry Stage 7 re-lock, after A3–A6 re-verify at Stage 4.

### Scope Verification

The Minimum Viable Validation is in scope for implementation: one fixture-backed
flow must prove `flow next`, `flow resolve`, `flow read-state`, and
`flow set-state` through the production Cobra path in both output modes,
including a read-assembled plan, a gated row, an escaped plan, a set write
with a clear, and at least one exit-2 and one exit-3 refusal.

### Cross-Cutting Concerns

- **Versioning**: verb-specific JSON `data` payloads and the `findings` field
  are append-only under the existing CLI output envelope.
- **Incremental adoption**: `--model <path>` and fixture-backed model loading
  ship before config-resolved `--flow <id>`.
- **Secret/credential lifecycle**: this RDR does not introduce credentials;
  accessor execution and external availability belong to RDR 0004.
- **Canonical-form / determinism**: the deterministic claim is request-level
  semantic determinism over the same model revision and artifact contents, not
  byte-identical output; the one byte-level claim is the canonical set literal
  (JDR 0001 §D13) surviving plan→request→read-back, which holds only because
  one encoder — HTML escaping disabled — renders it at every emit site.

### Proportionality

This RDR is right-sized for one load-bearing contract: the user-facing CLI
integration surface over the resolver, accessor, and output seams. It does not
own kernel selection semantics, transition-model format, guard predicate
meaning, accessor safety, or graph lint invariants; the four cross-seam
decisions it depends on live at JDR 0001 and are cited. The `mid` profile
remains appropriate because the contract is user-facing but not foundational
and carries no prior accretion in Seam Lineage.

## References

- `docs/cli-output-contract.md`
- `docs/jdr/0001-resolve-kernel-seam.md` §D2, §D4, §D8, §D9, §D10, §D11, §D13
- `docs/rdr/0005-skill-integration-cli-contract/evidence/spikes/d13-canonical-set.out`
  (§D13 canonical set-value byte form; the escaping edge is open)
- `docs/rdr/0005-skill-integration-cli-contract/evidence/research/citations.md`
- `internal/cli/respond::OK`
- `internal/cli/respond::Fail`
- `internal/cli/respond::ValidateMode`
- `internal/cli/clierr::CLIError`
- `internal/cli/clierr::ExitCodeFor`
- `internal/cli::ExecuteAndEmit`
- `internal/cli::newVersionCmd`
- `internal/resolve::Resolve`
- `internal/resolve::RefusalKinds`
- RDR 0001, 0002, 0003, 0004, and 0006
- `../state-machines/attic/RESOLVER-CLI.md`
- `../state-machines/attic/RESOLVER-DESIGN.md`
- `../state-machines/MODEL-transition.md`
