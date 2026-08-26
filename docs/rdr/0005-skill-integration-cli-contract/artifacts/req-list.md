# REQ List — RDR 0005 Skill Integration CLI Contract

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0005-skill-integration-cli-contract.md`. Quotes are verbatim — copied
from the projector (`rdr inspect --select <id>`) for fenced elements and read
from the record for testable prose outside the fences — never transcribed by
hand.

Element ids (`0005:C1`, `0005:MVV`, `0005:S3`, `0005:D-naming`) are carried
wherever a REQ derives from a labelled element, so a later stage can trace the
REQ back to its contract.

**Why the REQ count far exceeds `counts.elements.C = 1`.** This record has
exactly one `normative` fence (`0005:C1`), and that single fence carries eleven
independent paragraphs of MUST-clauses. Everything else testable is PROSE:
`Technical Design` (model selection, state input, outcome, gates, writes,
envelope, text mode), `Load-Bearing Decisions`, `Round-Trip / Inverse
Invariants`, `Trade-offs / Failure Modes` (the 25-row stable-code table plus the
gate-disposition mini-check), `Implementation Plan / Prerequisites`, and
`Validation / Testing Strategy`. Each of those is mined below.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced, `0005:C1`)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `LBD` = Technical Design / Load-Bearing Decisions
- `RT` = Technical Design / Round-Trip / Inverse Invariants
- `IC` = Technical Design / Illustrative Code
- `CAP` = Proposed Solution / Capability Dependencies
- `EIA` = Proposed Solution / Existing Infrastructure Audit
- `CONS` = Trade-offs / Consequences
- `RM` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes (code table + gate-disposition mini-check)
- `PRE` = Implementation Plan / Prerequisites
- `MVV` = Implementation Plan / Minimum Viable Validation
- `PH` = Implementation Plan / Phases
- `TS` = Validation / Testing Strategy
- `PE` = Validation / Performance Expectations
- `CA` = Research Findings / Critical Assumptions

**Standing note on what has already shipped.** RDR 0001 (`internal/resolve`),
0002 (`internal/table`), 0003/0007 (`internal/resolve` guard atoms), 0004
(`internal/accessor`), and 0006 (`internal/graphlint`, `internal/cli/lint.go`)
are all implemented in this worktree. In particular `internal/cli/clierr`
**already carries** `Findings []Finding` and the flat `Finding` record, landed
by RDR 0006. A REQ below that a shipped peer already discharges is satisfied by
citing that surface plus a test at *this* RDR's boundary; it is not a licence to
re-implement or to restate the peer's semantics. REQs in that class are marked
**[peer-delivered]**. Two of them (REQ-14, REQ-16) require an *amendment* to the
shipped type and are marked **[peer-delivered, amend]**.

---

## Command group and verb surface

- [REQ-1] "The CLI MUST expose one command group for skill integration with these verbs: next, resolve, read-state, and set-state." — (NC, `0005:C1`)

- [REQ-2] "Other command groups (lint, dump, parse) are outside this contract and are owned by the RDR that names them." — (NC, `0005:C1`) — the `flow` group MUST NOT absorb the shipped `lint` verb (`internal/cli/lint.go`, RDR 0006) or `dump` (RDR 0002).

- [REQ-3] "the user-facing group is `flow`, with verbs `next`, `resolve`, `read-state`, and `set-state`; flags `--flow`, `--model`, `--tag`, `--artifact`, `--outcome`, `--evaluate-gates`, `--write`, `--clear`." — (LBD, `0005:D-naming`) — the exact spellings are normative; `state` and `run` are named rejected group names in the same decision.

- [REQ-4] "Kernel-mapped codes mirror the kernel refusal kinds one-to-one (`flow-<kind>`, underscores to hyphens)." — (LBD, `0005:D-naming`)

## Output gateway and envelope

- [REQ-5] "All four verbs MUST start RunE by calling respond.ValidateMode(cmd), MUST route success through respond.OK, MUST route user-facing failure through respond.Fail(cmd, *clierr.CLIError), and MUST set SilenceErrors and SilenceUsage." — (NC, `0005:C1`)

- [REQ-6] "Under --as=json, each successful invocation MUST emit exactly one stdout JSON terminal envelope with type \"ok\" and verb-specific data." — (NC, `0005:C1`)

- [REQ-7] "Under --as=text, each successful invocation MUST emit human output derived from the same verb-specific result." — (NC, `0005:C1`)

- [REQ-8] "Failures MUST use the existing CLIError JSON/text envelope defined by docs/cli-output-contract.md and internal/cli/clierr, extended by exactly one omitempty structured field, findings, per JDR 0001 §D10." — (NC, `0005:C1`) — **[peer-delivered, amend]** the field exists as `Findings []Finding \`json:"findings"\`` (`internal/cli/clierr/clierr.go:69`); see REQ-16 for the `omitempty` amendment and ASSUMPTION A-6.

- [REQ-9] "Add text success payload rendering through `respond.OK` or a respond-owned helper used by `respond.OK`; do not print directly from resolver verbs." — (PRE) — the shipped `respond.OK` text branch renders only Notes, Warnings, and a `FindingCarrier` payload's findings (`internal/cli/respond/respond.go:125`); it renders no verb payload today.

- [REQ-10] "New verbs must not use direct Cobra printing for text payloads." — (EIA, *Success rendering* row) — `internal/cli::newVersionCmd`'s `cmd.Println` is the named drift and MUST NOT be copied.

- [REQ-11] "Text mode renders the same result content in a human-scannable order; it does not invent fields absent from the JSON payload, and it renders `findings` one-per-line under the message." — (TD, *Envelope*)

- [REQ-12] "The `--as` flag itself — its persistence at the root, its `text|json` domain, and its default — stays the root contract's (`docs/cli-output-contract.md`, `internal/cli::NewRootCmd`); this RDR inherits it unchanged and does not redefine it for the `flow` group." — (TD, *Envelope*)

- [REQ-13] "this RDR introduces verb-specific JSON `data` payloads and exactly one new `omitempty` `CLIError` field (`findings`), not a new terminal envelope or exit group." — (LBD, `0005:D-wire-byte-format`)

## The `Finding` record

- [REQ-14] "Finding MUST be one flat record with omitempty optional fields, carrying at least code, message, param, locator, hint, severity, model, rule, key, operator, literal, block, and class" — (NC, `0005:C1`) — **[peer-delivered, amend]** the shipped `clierr.Finding` (`internal/cli/clierr/clierr.go:257`) carries `code`, `model`, `severity`, `message`, `rule`, `span`, `element`, `reason`, `dimension`, `key`, `operator`, `literal`, `block`, `class`, `fingerprint`; it is **missing `param` and `locator`**, which this RDR MUST add as `omitempty` fields.

- [REQ-15] "each producer populates only the fields it owns, and no producer nests its own fields in a sub-object." — (NC, `0005:C1`) — restated in TD: "The atom's four fields sit flat on the record, not in a nested `atom` object."

- [REQ-16] "Add `Findings []Finding` (`json:\"findings,omitempty\"`) and the flat `Finding` record to `internal/cli/clierr`, shared with RDR 0006 — fields `Code`, `Message`, plus `omitempty` `Param`, `Locator`, `Hint`, `Severity`, `Model`, `Rule`, `Key`, `Operator`, `Literal`, `Block`, `Class`." — (PRE) — **[peer-delivered, amend]** see ASSUMPTION A-6 for the reconciliation with `0006:REQ-94`.

- [REQ-17] "Finding.message MUST be self-sufficient — it MUST render the failure readably with no structured field consulted." — (NC, `0005:C1`)

- [REQ-18] "This RDR populates `code`, `message`, `param`, `locator`, and `hint`; RDR 0006 additionally populates `severity`, `model`, `rule`, the guard atom's `key` / `operator` / `literal` / `block`, and `class`." — (TD, *Envelope*)

## Exit codes

- [REQ-19] "Exit 3 MUST mean the environment could not be consulted and the same request may be re-run unchanged; every other failure MUST exit 2. No new exit group." — (NC, `0005:C1`)

- [REQ-20] "Exit codes: **3** means the environment could not be consulted (accessor timeout, execution failure, incomplete read, read-back incomplete, post-mutation timeout) — repair it and re-run the same request unchanged; **2** means the request or the model is wrong, or the model said no." — (TD, *Envelope*) — the exit-3 population is exactly those five classes.

- [REQ-21] "an unavailable accessor surfaces as an exit-3 refusal and an indeterminate gate as an exit-2 refusal." — (PE)

## Model selection and loading

- [REQ-22] "Model selection MUST accept exactly one of --flow <id> or --model <path>." — (NC, `0005:C1`) — "neither / both of `--flow`, `--model`, or the selection does not resolve" is `flow-model-not-found`, `GroupUserEnv` / 2, carrier `param` (FM).

- [REQ-23] "The CLI MUST perform the model file I/O and hand RDR 0002's loader bytes plus a source id." — (NC, `0005:C1`) — TD adds "the loader performs no file I/O"; `internal/table::Load(src []byte, sourceID string)` is the shipped seam.

- [REQ-24] "Every load-time category MUST map to flow-model-invalid with one findings[] entry per category hit carrying its locator." — (NC, `0005:C1`) — TD scopes the population: "Every load category RDR 0002 defines, RDR 0003's two predicate categories, and RDR 0008's `reserved_tag_key` map to one code, `flow-model-invalid`, with one `findings[]` entry per category hit (`code` = category slug, `locator` = file:line)."

- [REQ-25] "`revision` on every payload is the loaded model's own revision identity, carried through verbatim … The CLI never derives, hashes, or synthesizes it, so a model that declares none renders `revision` empty rather than a CLI-invented value." — (TD, *Model selection*)

## State input: `--tag` and `--artifact`

- [REQ-26] "--tag values MUST enter as observed context." — (NC, `0005:C1`)

- [REQ-27] "A --tag naming an owned key or the reserved recognized key MUST be refused at the CLI before any accessor runs." — (NC, `0005:C1`) — the codes are `flow-tag-owned` and `flow-tag-reserved` respectively (TD, FM), both `GroupUserEnv` / 2, carrier `param`.

- [REQ-28] "A duplicate tag name is `flow-tag-duplicate`." — (TD, *State input*; FM)

- [REQ-29] "A set-kind tag's value is a JSON array literal (`'labels=[\"a\",\"b\"]'`), the same canonical form the kernel seam carries (JDR 0001 §D13) and `--write` accepts — sorted, duplicate-free, compact, and rendered with HTML escaping disabled; a bare scalar for a set key, or an array for a scalar key, is `flow-tag-invalid`." — (TD, *State input*)

- [REQ-30] "Owned state is never caller-supplied: it is assembled from declared read accessors over `--artifact role=path` bindings, where `role` is the model-declared artifact role and `path` the caller-owned artifact." — (TD, *State input*)

- [REQ-31] "Nothing discovers artifacts ambiently; location comes only from explicit bindings." — (AP) — restated fenced: "They MUST NOT discover artifacts" (NC, `0005:C1`).

- [REQ-32] "malformed `--artifact` binding" is `flow-artifact-invalid`, `GroupUserEnv` / 2, carrier `param`. — (FM)

- [REQ-33] "A role an invoked accessor needs that no binding supplies is `flow-artifact-missing`." — (TD, *State input*; FM, carrier `param` = role)

- [REQ-34] "`--tag` on `set-state` is context only and is never written." — (TD, *Writes*) — restated fenced: "It MUST NOT treat --tag values as writes." (NC, `0005:C1`)

## Narrowed read-accessor set (A7)

- [REQ-35] "The invoked read-accessor set MUST be exactly those readers serving an owned key some candidate row of the requested model requires — not every declared reader." — (NC, `0005:C1`)

- [REQ-36] "A reader no candidate row needs MUST NOT run, and its unbound artifact role MUST NOT raise flow-artifact-missing; flow-artifact-missing is scoped to the roles the invoked set needs." — (NC, `0005:C1`)

- [REQ-37] "flow read-state is the exception and runs every declared reader, because a diagnostic read has no candidate set to narrow by." — (NC, `0005:C1`) — restated in TD: "`read-state` runs them all."

- [REQ-38] "flow next and flow resolve MAY run declared read accessors over explicit --artifact role=path bindings and MUST assemble owned state only from them (JDR 0001 §D8)." — (NC, `0005:C1`)

- [REQ-39] "They … MUST NOT run write accessors." — (NC, `0005:C1`)

## `flow next`

- [REQ-40] "flow next MUST return the legal recognized-outcome alphabet for the supplied state, plus candidate summaries containing source rule identity, required facts, unresolved guard/gate facts, and preview next tags, write targets, and clear keys when those can be read from normalized model data without evaluating missing facts." — (NC, `0005:C1`)

- [REQ-41] "It MUST run gate accessors only when --evaluate-gates is given; otherwise it MUST list gate ids as unresolved facts." — (NC, `0005:C1`)

- [REQ-42] "Under --evaluate-gates it MUST run the gates of every candidate row it reports and no others — a row whose guards already exclude it is not a candidate, so its gates MUST NOT run — and it MUST report each gate result on the candidate that carries it." — (NC, `0005:C1`)

- [REQ-43] "A deny on one candidate constrains only that candidate: flow next reports it and still exits 0, because next enumerates rather than selects, and only flow resolve turns a deny into flow-gate-denied." — (NC, `0005:C1`)

- [REQ-44] "It MUST NOT invent guard facts that were neither supplied, read, nor produced by a declared gate accessor." — (NC, `0005:C1`)

- [REQ-45] "It never calls a language model." — (AP)

- [REQ-46] The `next` gate-disposition row set (FM mini-check): "gate allows on a reported/selected row" → "exit 0; result on the candidate"; "gate denies" → "exit 0; result on the candidate, no refusal"; "gate indeterminate, none deny" → "exit 0; result on the candidate"; "gate times out / fails to execute" → "exit 3 `flow-accessor-timeout` / `flow-accessor-failed`"; "gate on a guard-excluded row" → "not run, not reported"; "gates not requested (`next` default)" → "not run; gate ids listed as unresolved facts". — (FM)

- [REQ-47] "no gate result is ever dropped silently, and a gate that could not be consulted is exit 3 rather than a deny." — (FM) — binding in both verbs.

## `flow resolve`

- [REQ-48] "flow resolve MUST return exactly one plan or exactly one CLIError refusal." — (NC, `0005:C1`)

- [REQ-49] "Kernel refusals MUST map one-to-one onto flow-unmodeled-outcome, flow-no-match, flow-ambiguous-match, flow-owned-state-unavailable, and flow-guard-unevaluable." — (NC, `0005:C1`) — the shipped closed set is `internal/resolve::RefusalKinds()`.

- [REQ-50] "It MUST run the selected row's gates only after exact-one selection and before emitting the plan (JDR 0001 §D9); all gates on the row MUST run, deny MUST override allow and indeterminate, and every gate result MUST be reported." — (NC, `0005:C1`)

- [REQ-51] "A denied or indeterminate gate MUST surface as flow-gate-denied or flow-gate-indeterminate with one findings[] entry per gate; it is never a plan and never an escape class." — (NC, `0005:C1`) — FM fixes the precedence: `flow-gate-indeterminate` fires only when "no gate denied and at least one indeterminate".

- [REQ-52] "An escaped plan MUST be a success carrying escaped=true and escape_class." — (NC, `0005:C1`) — restated in FM: "An escaped plan is not a failure: `flow resolve` exits 0 with `escaped: true` and `escape_class` on the payload."

- [REQ-53] "flow resolve MUST NOT print directly, initiate skill work, or choose among multiple matching rows." — (NC, `0005:C1`)

- [REQ-54] "`flow resolve` requires `--outcome <tag>`; an absent or empty value is `flow-tag-invalid` at the CLI, and a value outside the model's alphabet is the kernel's `flow-unmodeled-outcome`." — (TD, *Outcome*)

- [REQ-55] "A gate's timeout or execution failure is an accessor refusal (exit 3), not a gate result." — (TD, *Gates*)

- [REQ-56] "`flow resolve` succeeds only when the kernel reports exactly one matching row (or exactly one escape row for an escapable class) and every gate on that row allows. Zero matches, multiple matches, unmodeled outcomes, unavailable owned state, unevaluable guards, and gate deny/indeterminate become typed refusals rather than tie-breaks." — (LBD, `0005:D-selection-predicate`)

- [REQ-57] "`flow resolve` performs reads, so it is no longer effect-free; `flow next` stays effect-free by default and gates only on opt-in." — (CONS)

## `flow read-state`

- [REQ-58] "flow read-state MUST invoke only declared read accessors over caller-supplied role=path artifact bindings and return, per reader, its declared keys and the tag-set read, or a stable CLIError failure." — (NC, `0005:C1`)

- [REQ-59] "It MUST NOT invoke gate accessors or coerce gate allow, deny, or indeterminate results into tag values." — (NC, `0005:C1`)

## `flow set-state`

- [REQ-60] "flow set-state MUST invoke only declared write accessors over caller-supplied role=path artifact bindings and the planned owned-tag mutations given as --write name=value and --clear <key> (JDR 0001 §D11)." — (NC, `0005:C1`)

- [REQ-61] "Before any accessor runs it MUST refuse a --write or --clear key that is not a declared owned tag served by exactly one [write.<id>] whose keys list names it, a --write value malformed for its declared kind, the literal value <clear>, and a key given twice." — (NC, `0005:C1`) — the four codes are `flow-write-unbound` / `flow-clear-unbound`, `flow-write-invalid`, `flow-write-invalid` (with hint "use `--clear`"), and `flow-write-duplicate` (TD, *Writes*; FM).

- [REQ-62] "The sentinel `<clear>` is unauthorable: `--write k=<clear>` is refused `flow-write-invalid` with the hint \"use `--clear`\"." — (TD, *Writes*)

- [REQ-63] "The same key under both flags, or twice under either, is `flow-write-duplicate`." — (TD, *Writes*)

- [REQ-64] "It MUST never run gate accessors." — (NC, `0005:C1`) — restated in AP: "It never runs gates".

- [REQ-65] "It MUST report success only after the accessor layer's read-back verification confirms the planned owned-tag values and that non-owned tags are unchanged (RDR 0004)." — (NC, `0005:C1`)

- [REQ-66] "`[]` (empty set) and `--clear` (absent) stay distinct through read-back." — (TD, *Writes*)

- [REQ-67] "Each writer applies its own keys and reads back; no cross-writer atomicity is promised." — (TD, *Writes*)

- [REQ-68] "`set-state` trusts the caller to transcribe the plan; nothing links a `set-state` request to a prior `resolve`." — (TD, *Writes*) — a stated non-guarantee, not a defect (CONS).

- [REQ-69] "A carried plan artifact (`--plan <file|->`) is deferred as a seed, not part of this contract." — (TD, *Writes*) — no `--plan` flag ships.

## Canonical set literal and the shared encoder

- [REQ-70] "A set-valued tag crossing the CLI in --tag, --write, or any payload MUST use the canonical JSON array form JDR 0001 §D13 fixes — members sorted, duplicate-free, compact — rendered with HTML escaping DISABLED, so <, >, and & serialize as themselves and never as \\u003c, \\u003e, or \\u0026." — (NC, `0005:C1`)

- [REQ-71] "Every site that emits or compares a canonical set literal MUST use the same encoder, so plan-to-request copy-through and read-back equality are byte equality." — (NC, `0005:C1`)

- [REQ-72] "Route `internal/cli/clierr::EmitJSON` and `internal/cli/respond::writeJSONLine` through one shared helper using `SetEscapeHTML(false)` … Bare `json.Marshal` must not remain on a path a set value crosses." — (PRE) — both call sites are shipped and both use bare `json.Marshal` today (`clierr.go:146`, `respond.go:203`).

- [REQ-73] "Set values are JDR 0001 §D13's canonical JSON array everywhere they cross the CLI, rendered with HTML escaping disabled so `<`, `>`, and `&` serialize as themselves." — (LBD, `0005:D-wire-byte-format`)

## Success payload shapes

- [REQ-74] `next` data minimum: "`model` (id or path) and `revision`, `observed` tags, `owned` tags assembled from the readers, `readers[]` (accessor identities invoked), `outcomes[]` (recognized outcome tags — nothing else; RDR 0002's layout has no per-outcome recognizer text and `[model.metadata]` is not its carrier), and `candidates[]`." — (TD, *Envelope*)

- [REQ-75] `next` candidate shape: "Each candidate carries source rule identity, the outcome it belongs to, required facts, unresolved guard/gate facts, evaluated gate results when `--evaluate-gates` was given, and preview next tags, write targets, and clear keys when the normalized model exposes them without evaluating missing facts." — (TD, *Envelope*)

- [REQ-76] `resolve` data minimum: "`model`, `revision`, `observed`, `owned`, `readers[]`, `outcome`, `rule` (matched rule identity), `gates[]` (`id`, `result`), `next` (the next tag-set), `writes` (planned owned-tag writes, set values as JSON arrays), `clear[]` (planned clears), `escaped` (bool) and `escape_class` when the plan came from an escape row." — (TD, *Envelope*)

- [REQ-77] `read-state` data minimum: "`model`, `revision`, artifact role bindings, `readers[]` each with its declared `keys` (the requested key set) and the tags it returned — so \"absent from the artifact\" and \"not requested\" are distinguishable from the payload alone." — (TD, *Envelope*)

- [REQ-78] `set-state` data minimum: "`model`, `revision`, artifact role bindings, `writers[]`, requested `writes` and `clear[]`, and the read-back-confirmed owned-tag values." — (TD, *Envelope*)

- [REQ-79] "`param` for single-subject failures, and `findings[]` where several rows, keys, atoms, gates, or load categories must be named." — (FM)

## Stable code table (FM — "the spellings are normative; the group fixes the exit")

- [REQ-80] `flow-tag-invalid` — "malformed, empty, or wrong-kind `--tag` / `--outcome` value" — `GroupUserEnv` / 2 — carrier `param`. — (FM)
- [REQ-81] `flow-tag-duplicate` — "duplicate `--tag` name" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-82] `flow-tag-reserved` — "`--tag` names the reserved `recognized` key" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-83] `flow-tag-owned` — "`--tag` names an owned key" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-84] `flow-write-invalid` — "malformed or wrong-kind `--write` value, or the literal `<clear>`" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-85] `flow-write-duplicate` — "`--write` / `--clear` key given twice across either flag" — `GroupUserEnv` / 2 — `param`. — (FM) — RDR-owned spelling (CA A6).
- [REQ-86] `flow-write-unbound` — "`--write` key served by no single `[write.<id>].keys`" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-87] `flow-clear-unbound` — "`--clear` key served by no single `[write.<id>].keys`" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-88] `flow-artifact-invalid` — "malformed `--artifact` binding" — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-89] `flow-artifact-missing` — "a role an invoked accessor needs has no binding" — `GroupUserEnv` / 2 — "`param` = role". — (FM)
- [REQ-90] `flow-model-not-found` — "neither / both of `--flow`, `--model`, or the selection does not resolve" — `GroupUserEnv` / 2 — `param`. — (FM) — CA A6 records that this code "covering the selection *arity* error as well as non-resolution" is RDR-owned.
- [REQ-91] `flow-model-invalid` — "any RDR 0002 load category, RDR 0003 predicate category, or RDR 0008 `reserved_tag_key`" — `GroupUserEnv` / 2 — "`findings[]` (category, locator)". — (FM)
- [REQ-92] `flow-unmodeled-outcome` — kernel `unmodeled_outcome` — `GroupUserEnv` / 2 — `param`. — (FM)
- [REQ-93] `flow-no-match` — kernel `no_match` — `GroupUserEnv` / 2 — "`findings[]` (rows considered)". — (FM)
- [REQ-94] `flow-ambiguous-match` — kernel `ambiguous_match` — `GroupUserEnv` / 2 — "`findings[]` (matching rows)". — (FM)
- [REQ-95] `flow-owned-state-unavailable` — kernel `owned_state_unavailable` — `GroupUserEnv` / 2 — "`findings[]` (keys)". — (FM)
- [REQ-96] `flow-guard-unevaluable` — kernel `guard_unevaluable` — `GroupUserEnv` / 2 — "`findings[]` (rows; atoms when the kernel names them)". — (FM)
- [REQ-97] `flow-gate-denied` — "a selected-row gate denied" — `GroupUserEnv` / 2 — "`findings[]` one per gate". — (FM)
- [REQ-98] `flow-gate-indeterminate` — "no gate denied and at least one indeterminate" — `GroupUserEnv` / 2 — "`findings[]` one per gate". — (FM)
- [REQ-99] `flow-accessor-timeout` — "accessor timed out" — `GroupEnvUnavailable` / 3 — "`param` = accessor id". — (FM)
- [REQ-100] `flow-accessor-failed` — "accessor execution failed" — `GroupEnvUnavailable` / 3 — "`param` = accessor id". — (FM)
- [REQ-101] `flow-read-incomplete` — "read returned an incomplete key set" — `GroupEnvUnavailable` / 3 — "`param` = accessor id". — (FM)
- [REQ-102] `flow-accessor-unknown`, `flow-accessor-capability-mismatch`, `flow-write-non-owned` — "unknown accessor, capability mismatch, write to a non-owned tag (unreachable past load and CLI checks)" — `GroupInternal` / 2 — no carrier. — (FM)
- [REQ-103] `flow-write-readback-mismatch` — "read-back disagrees with the plan, or a non-owned tag changed" — `GroupUserEnv` / 2 — "`findings[]` per key". — (FM)
- [REQ-104] `flow-write-readback-incomplete` — "read-back returned an incomplete key set" — `GroupEnvUnavailable` / 3 — "`detail`: may have been applied". — (FM)
- [REQ-105] `flow-write-readback-timeout` — "post-mutation read-back timed out" — `GroupEnvUnavailable` / 3 — "`detail`: may have been applied". — (FM)

- [REQ-106] "The main silent-failure risk is treating a failed accessor, a denied gate, or an ambiguous row as a successful transition. The recovery path is refusal-first: the command exits non-zero, emits the structured error envelope, and leaves skill judgment or model repair to the caller." — (FM)

## Round-trip / inverse invariants

- [REQ-107] "`set-state` followed by `read-state` must return the planned owned-tag values for the artifact role that was written. The equality is value-for-value over the owned tags the write planned to mutate — byte equality for set values in canonical form — not byte-identical artifact content; a cleared key reads back absent, an empty set reads back `[]`." — (RT)

- [REQ-108] "`resolve` → skill → `set-state` is copy-through: every `writes` value and `clear[]` key in a plan is accepted verbatim by `set-state`'s `--write` / `--clear` grammar, so a set write survives the round trip byte-identical." — (RT)

- [REQ-109] "`read-state` → `--tag` re-pairs only for observed tags; an owned tag read by `read-state` cannot be handed back through `--tag` (refused `flow-tag-owned`) because `resolve` reads it itself." — (RT)

## Request identity and determinism

- [REQ-110] "a CLI request is identified by verb, model selection (`--flow` id or `--model` path) and model revision, observed tags, artifact role bindings, recognized outcome when applicable, and planned writes/clears when applicable. The same request over the same model revision and the same artifact contents must produce the same success or refusal, excluding exit-3 environment failures." — (LBD, `0005:D-identity`)

- [REQ-111] "The command path should stay single-invocation deterministic: parse inputs, load the selected model, run the declared readers, call one kernel operation, run the selected row's gates, and render one terminal result." — (PE)

- [REQ-112] "No new third-party dependency. Cobra, `respond`, and `clierr` already exist; JSON array literals parse with the standard library. Any TOML dependency is RDR 0002's." — (Implementation Plan / New Dependencies)

## Ownership fences (negative REQs)

- [REQ-113] "It does not own transition-model representation, guard semantics, accessor safety, graph lint, skill execution, or constrained decoding." — (AP)

- [REQ-114] "Caller-supplied `--tag` values are observed context and never satisfy an owned dependency." — (AP)

- [REQ-115] "`flow read-state` … does not invoke gate accessors, because gates return allow, deny, or indeterminate rather than tag values." — (AP)

- [REQ-116] "A denied gate is a refusal, not a plan and not an escape class." — (AP)

- [REQ-117] "the read-back is the commit-time check, and the window between `resolve` and `set-state` is a stated non-guarantee of the two-verb design (§D9)." — (AP)

- [REQ-118] "The contract requires candidate summaries read from normalized data and forbids evaluating missing facts; gates run only on opt-in." — (RM) — the mitigation for "`flow next` becomes a second guard evaluator".

- [REQ-119] "Kernel codes mirror the kernel kinds one-to-one; gates run after selection and deny is its own refusal" — (RM) — the mitigation for "a gate deny is laundered into `no_match` and escaped".

- [REQ-120] "Derive both renderings from one typed result per verb and test both modes." — (RM) — the mitigation for "Text and JSON outputs diverge semantically".

- [REQ-121] "the kernel's unexported pre-selection guard filter also spells itself `gate`; it is a different mechanism from this RDR's post-selection accessor gates and must not be wired to them." — (CA A4, implementation note)

## Minimum Viable Validation

- [REQ-MVV] `0005:MVV` — "Implement one fixture-backed flow and prove all four verbs through the production Cobra path: `flow next` returns the legal outcome alphabet with gate ids as unresolved facts; `flow resolve` reads owned state from a fixture artifact, maps one outcome to one plan, and runs one gate on the selected row; `flow read-state` reads the fixture artifact tags per reader; and `flow set-state` persists one scalar write, one set write as a JSON array, and one `--clear`, then read-back-verifies them. The set write's members MUST include one containing `<` and one containing `&`, asserted byte-identical through plan → request → read-back, so the encoder rule is covered by the validation rather than only by review. Run the happy path, one escaped plan, one gate deny, one kernel refusal, and one exit-3 accessor failure in `--as=text` and `--as=json`, asserting `findings[]` where the table names it.

  The fixture model MUST also declare one reader no candidate row needs, with its artifact role left unbound: `flow next` and `flow resolve` MUST succeed without invoking it or raising `flow-artifact-missing`, while `flow read-state` on the same model MUST invoke it and refuse the missing binding. That pair is the only assertion that distinguishes the narrowed invoked set from \"every declared reader\". `flow next --evaluate-gates` over a model with one gated candidate and one guard-excluded row MUST report the gate result on the reported candidate, run no gate for the excluded row, and exit 0 even when that gate denies." — (MVV, `0005:MVV`)

## Testing Strategy scenarios

- [REQ-122] "Command tests exercise the `flow` command group through the production Cobra path rather than calling renderers or kernel functions directly. Coverage must include argument validation, output mode validation, success rendering, and typed refusal mapping for each verb." — (TS)

- [REQ-123] `0005:S1` — "`flow next` over the fixture model in `--as=json` and `--as=text`, with and without `--evaluate-gates`." Expected: "both modes report the same recognized-outcome alphabet and candidate summaries through the standard success path; JSON includes `outcomes[]` and `candidates[]`; without the flag no gate accessor runs and gate ids appear as unresolved facts." — (TS, `0005:S1`)

- [REQ-124] `0005:S2` — "`flow resolve` over the fixture model with one recognized outcome that matches exactly one row, whose owned dependency is served by a fixture reader and which carries one gate." Expected: "the reader runs, `owned` is assembled from it, the plan carries `rule`, `next`, `writes`, `clear[]`, and `gates[]` with `allow`; no write accessor runs; a `--tag` on the owned key is refused `flow-tag-owned` before the reader runs." — (TS, `0005:S2`)

- [REQ-125] `0005:S3` — "`flow resolve` with an unmodeled outcome, an empty `--outcome`, a zero-match row, a multi-match row, and an escape row for `no_match`." Expected: "each refusal maps to its one-to-one code with non-zero exit under both modes; the escape case exits 0 with `escaped: true` and `escape_class`." — (TS, `0005:S3`)

- [REQ-126] `0005:S4` — "`flow resolve` where the selected row's gates return deny + allow, and separately indeterminate + allow." Expected: "`flow-gate-denied` then `flow-gate-indeterminate`, each with one `findings[]` entry per gate; deny overrides indeterminate when both occur; no plan is emitted." — (TS, `0005:S4`)

- [REQ-127] `0005:S5` — "`flow read-state` and `flow set-state` over a fixture artifact and declared accessor roles, with a scalar write, a set write as a JSON array, and a `--clear`." Expected: "`--artifact` bindings are validated; `read-state` invokes only declared readers and reports each reader's `keys` beside its tags; `set-state` reports success only after read-back proves the scalar, the set (byte-equal canonical array), and the cleared key absent; `--write k=<clear>`, an unbound key, and a wrong-kind value are refused before any accessor runs. One set member carries `<` and one carries `&`: both render as themselves at every emit site — `[\"a<b\",\"x&y\"]`, never `[\"a\\u003cb\",\"x\\u0026y\"]` — and plan, request, and read-back agree byte-for-byte (normative fixture, `evidence/spikes/escaping-surfaces.out`)." — (TS, `0005:S5`)

- [REQ-128] `0005:S6` — "read accessor timeout during `resolve`, gate execution failure, write read-back incomplete, and write read-back mismatch." Expected: "the first three exit 3 with `GroupEnvUnavailable` and a \"may have been applied\" detail on the write case; the mismatch exits 2 with `findings[]` per key; none is a successful transition or a coerced `read-state` payload." — (TS, `0005:S6`)

- [REQ-129] `0005:S7` — "a refusal carrying `findings[]` is rendered in both modes, and the same `Finding` record is populated by a kernel refusal, a gate result, and a model-load category." Expected: "each producer populates only its own fields and unset optional fields are absent from the JSON (`omitempty`); no producer nests its fields in a sub-object; and each finding's `message` alone renders the failure readably with no structured field consulted." — (TS, `0005:S7`)

- [REQ-130] `0005:S8` — "a fixture model declaring one reader no candidate row needs (its role unbound) and, separately, `flow next --evaluate-gates` over one gated candidate plus one guard-excluded gated row." Expected: "`next` and `resolve` neither invoke the unneeded reader nor raise `flow-artifact-missing`, while `read-state` on the same model invokes it and refuses the missing binding; under `--evaluate-gates` only the reported candidate's gate runs, its result rides that candidate, and a deny there still exits 0 (`next` enumerates, it does not select)." — (TS, `0005:S8`)

## Documentation

- [REQ-131] "Land the full code table, update `docs/cli-output-contract.md` for the `findings` field and the exit-3 rule, and document illustrative invocations." — (PH, Phase 4)

- [REQ-132] The four illustrative invocation shapes in `Illustrative Code` must parse under the shipped flag grammar (they are non-normative as *output*, but they exercise `--model`, `--flow`, `--artifact`, `--tag`, `--outcome`, `--write` scalar, `--write` set literal, `--clear`, and `--as=json`). — (IC)

---

## ASSUMPTIONS

- **A-1 — `--flow <id>` config lookup may be stubbed in this RDR's shipped
  surface.** `EIA`'s *Config discovery* row says "`--model <path>` ships first;
  `--flow <id>` config lookup may follow", while `NC` requires "exactly one of
  --flow <id> or --model <path>" and `FM` gives `flow-model-not-found` for a
  selection that "does not resolve". **Reading taken:** both flags are
  *registered* and the mutual-exclusion / neither-given arms are fully
  enforced (REQ-22), but a `--flow <id>` that config discovery cannot resolve
  refuses `flow-model-not-found` rather than being a missing feature. This
  mirrors the shipped `lint` verb, which registers `--flow` alongside `--model`
  with the same exclusion check (`internal/cli/lint.go:56-57`).

- **A-2 — the "reserved `recognized` key" is a literal key name.** `TD` says
  "one naming `recognized` is refused `flow-tag-reserved`". Read as the exact
  tag key `recognized`, consistent with `internal/resolve::Input.Recognized`
  being the kernel's separate channel and with RDR 0008's `reserved_tag_key`
  load category. `flow-tag-reserved` is CLI-side and precedes any accessor run.

- **A-3 — `escape_class` is derived at the CLI, not read off `Plan`.** The
  shipped `internal/resolve::Plan` carries `Escaped bool` (`resolve.go:332`) but
  no escape-class field. **Reading taken:** the CLI recovers the class by
  looking the plan's `RuleID` up in the loaded table and reporting the escape
  class the request's primary refusal kind matched (`Row.Escape`, `resolve.go:277`).
  This adds no kernel change and keeps REQ-52 satisfiable. If a row models
  several classes, the reported class is the one that was rescued — the
  request's own refusal kind — not the whole `Escape` list.

- **A-4 — "candidate row" for the narrowing rule (REQ-35) means a row surviving
  guard evaluation over what is known *before* reads run, not a row the kernel
  selected.** Reads must run before the kernel call (REQ-38), so the narrowing
  set is computed from normalized model data: rows whose `Outcome` matches the
  request (for `resolve`) or any modeled outcome (for `next`), taking each such
  row's `RequiresOwned` (`internal/table`, `0002::Row.RequiresOwned`) as its
  owned demand and invoking exactly the readers whose declared `keys` serve
  those demands. CA A7 grounds both legs as "dumpable columns of the normalized
  value" available "at load and before any evaluation".

  **SUPERSEDED IN PART by DEV-8** (`artifacts/deviations.md`). The
  candidate-row leg stands; the OWNED-DEMAND leg does not. Reading the demand
  set as `RequiresOwned` ALONE under-read REQ-35's "an owned key some candidate
  row **requires**": RDR 0002 derives `RequiresOwned` from the write block and
  clear list only, so a row that GUARDS on an owned key it never writes demands
  that key without naming it there. Phase 3b's ADV-2 showed the consequence — the
  serving reader is never invoked and `flow resolve` refuses
  `flow-guard-unevaluable` over an artifact that HOLDS the fact. The demand set
  is `RequiresOwned` unioned with the row's guard-block (`all` / `unless`) atom
  keys the model declares owned. Both legs remain dumpable normalized model data
  per CA A7.

- **A-5 — "preview next tags, write targets, and clear keys" on a `next`
  candidate are the normalized row's `NextTags`, `Writes`, and clear list read
  verbatim.** REQ-40/REQ-75 qualify them "when those can be read from
  normalized model data without evaluating missing facts". Since RDR 0002
  normalizes an authored clear into a `<clear>` write
  (`internal/resolve::Row.RequiresOwned` doc comment), the clear keys are the
  write entries carrying the `<clear>` sentinel, split out for the payload.

- **A-6 — `CLIError.Findings` moves to `json:"findings,omitempty"`, and RDR
  0006's shipped tests stay green.** REQ-8/REQ-16 require `omitempty`; RDR
  0006's shipped code deliberately declares it non-`omitempty`
  (`internal/cli/clierr/clierr.go:62-69`) and `0006:REQ-94` says "the `findings`
  field MUST NOT be `omitempty`". **Reading taken — 0005's own A5 settles
  this** and I verified it against the shipped tests: §D10 item 3 makes the
  *`CLIError`* field `omitempty` while 0006's "empty list is the receipt" binds
  "the success payload's own non-omitempty `data.findings`". The two shipped
  failure-envelope assertions (`internal/cli/lint_0006_test.go:482`,
  `:593`, `internal/cli/lint_mvv_0006_test.go:466`) all run against a lint
  failure whose findings list is non-empty, so `omitempty` never elides the key
  on those paths and they remain green. The success-path assertions
  (`lint_0006_test.go:573`, `lint_mvv_0006_test.go:81,447`) read
  `data.findings` on `respond.Success.Data`, which this change does not touch.
  Stage 8 must run 0006's suite to confirm.

- **A-7 — the `Finding` fields this RDR adds are `Param` and `Locator` only.**
  REQ-14's list is "at least" thirteen fields; twelve are shipped. `param` and
  `locator` are absent. `locator` is *not* the shipped `span` or `element`:
  REQ-24 fixes its content as "file:line", and REQ-18 assigns `locator` to this
  RDR and `rule`/`span`/`element` to 0006. Both are added as `omitempty`
  strings; no shipped field is renamed or removed (append-only, per the type's
  own doc comment).

- **A-8 — the shared canonical encoder is a new `clierr`-level (or lower)
  helper both emit sites call.** REQ-72 names `clierr::EmitJSON` and
  `respond::writeJSONLine`. `respond` already imports `clierr`, so a helper
  exported from `clierr` (or a small shared package) satisfies "one shared
  helper" without inverting the dependency. `SetEscapeHTML(false)` on
  `json.Encoder` also appends a newline, so the helper must account for the
  existing `Fprintln` / newline behaviour at both sites rather than
  double-writing it.

- **A-9 — accessor refusal classes map to `flow-*` codes by class, and
  read-back classes are disambiguated by the invoking verb.**
  `internal/accessor::RefusalClass` ships eight values. The mapping this RDR's
  FM table implies: `timeout` → `flow-accessor-timeout` (or
  `flow-write-readback-timeout` when raised by `set-state`'s post-mutation
  read-back), `execution_failure` → `flow-accessor-failed`, `incomplete_read` →
  `flow-read-incomplete` (or `flow-write-readback-incomplete` from `set-state`),
  `read_back_mismatch` → `flow-write-readback-mismatch`, `gate_indeterminate` →
  `flow-gate-indeterminate`, `unknown_accessor` → `flow-accessor-unknown`,
  `capability_mismatch` → `flow-accessor-capability-mismatch`. The verb-phase
  split is forced: FM lists both `flow-read-incomplete` and
  `flow-write-readback-incomplete` at exit 3 with distinct descriptions, and
  only the invoking phase distinguishes them.

- **A-10 — `flow-write-non-owned` has no accessor `RefusalClass` and is raised
  from the CLI's own pre-accessor check or an internal invariant.** FM groups it
  with the two "unreachable past load and CLI checks" internal codes
  (`GroupInternal` / 2). Read as a defensive arm, not a reachable user path;
  REQ-61's `flow-write-unbound` is what a caller actually hits.

- **A-11 — text mode for the four success payloads is human-scannable but
  field-faithful.** REQ-7/REQ-11 forbid inventing fields and require `findings`
  one-per-line, but fix no layout. Read as: layout free, content complete —
  every field named in REQ-74..REQ-78 is representable in text, and REQ-120's
  "one typed result per verb" makes the two renderings share a source struct.
  `clierr::EmitFindingsText` is the shipped precedent for the one-per-line rule.

- **A-12 — `flow next` reports `outcomes[]` as the model's recognized-outcome
  alphabet filtered to outcomes with at least one candidate row surviving the
  supplied state.** REQ-40 says "the legal recognized-outcome alphabet **for the
  supplied state**", which is narrower than the model's whole declared alphabet.
  REQ-74's "nothing else" governs the *shape* of each entry (a bare tag, no
  recognizer text), not the population. AP agrees: "It returns the
  recognized-outcome alphabet plus candidate summaries". If no fact is supplied,
  no row is excluded and the two coincide.

---

## QUESTIONS

One clause admits two readings that produce materially different behaviour and
no predecessor artifact settles it. This run is unattended and proceeds under
the best-supported reading; the alternative is recorded so the orchestrator can
overrule.

- **Q1 — Does `flow next --evaluate-gates` narrow the invoked *read*-accessor
  set the same way `flow resolve` does, given that `next` may report candidates
  for every modeled outcome?** REQ-35 fixes the invoked reader set as "exactly
  those readers serving an owned key some candidate row of the requested model
  requires", and REQ-37 makes `read-state` the single named exception. But
  `next` takes no `--outcome`, so *every* row of the model is potentially a
  candidate, and the union of every row's `RequiresOwned` may equal "every
  declared reader" — collapsing the distinction the MVV calls "the only
  assertion that distinguishes the narrowed invoked set from 'every declared
  reader'". The two readings are (a) `next` narrows by the union over all rows,
  so a reader is skipped only when it serves no owned key any row requires, and
  (b) `next` narrows by the union over rows still viable under the supplied
  `--tag` facts, so supplying facts shrinks the reader set. **Proceeding on (a)**
  — it is the literal reading of REQ-35 ("some candidate row of the requested
  model requires", with no state qualifier), it makes the MVV's unbound-reader
  fixture discriminating (that reader serves an owned key *no* row requires, so
  it is skipped under either reading), and it keeps `next` deterministic in the
  sense REQ-110 requires without making the reader set depend on how much
  observed context a caller happened to pass. Reading (b) would additionally
  make `flow-artifact-missing` fire or not fire based on `--tag` content, which
  REQ-110's identity clause does not list as an input to reader selection.
  Escalation is required only if (a) proves to run readers the MVV's fixture
  asserts must not run.
