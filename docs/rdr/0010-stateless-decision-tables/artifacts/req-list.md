# REQ List — RDR 0010 Stateless Decision Tables

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0010-stateless-decision-tables.md`. Quotes are verbatim — copied from
the projector (`rdr inspect --select <id>`) for fenced elements and read from
the record for testable prose outside the fences — never transcribed by hand.

Element ids (`0010:C4`, `0010:MVV`, `0010:S3`) are carried wherever a REQ
derives from a labelled element, so a later stage can trace the REQ back to its
contract.

**Shape.** `counts.elements.C = 5`, and the five fences (C1–C5) carry the bulk
of the normative surface — but each fence carries many independent obligations,
and a material remainder is unfenced PROSE: `Approach`, `Technical Design`
(the three-seam list, the four-`len(Initial)`-site partition, the data flow),
`Authority`, `Load-Bearing Decisions`, `Existing Infrastructure Audit`,
`Trade-offs / Failure Modes`, `Cross-Cutting Concerns`, `Implementation Plan`
phases, and the six numbered `Testing Strategy` scenarios (`0010:S1`–`0010:S6`)
plus its `Done`/licensed-diff rule. The `Critical Assumptions` (A1–A15) are all
`Verified`; they are read as *evidence*, not as REQs, except where an assumption
fixes an implementation obligation the contracts delegate to it (A4, A10, A13,
A14, A15) — those are recorded and cited.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `AUTH` = Technical Design / Authority
- `LBD` = Technical Design / Load-Bearing Decisions
- `IC` = Technical Design / Illustrative Code (**non-normative**, per its own
  "tests must not assert these literally")
- `CAP` = Proposed Solution / Capability Dependencies
- `AUDIT` = Proposed Solution / Existing Infrastructure Audit
- `CONS` = Trade-offs / Consequences
- `RISK` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes
- `IP` = Implementation Plan (prerequisites, phases)
- `MVV` = Implementation Plan / Minimum Viable Validation
- `TS` = Validation / Testing Strategy (Done clause, Oracle, Trace)
- `SC-n` = Validation / Testing Strategy, numbered scenario *n* (`0010:Sn`)
- `XC` = Finalization Gate / Cross-Cutting Concerns (`0010:G-cross-cutting`)

---

## A. Model class declaration (C1)

- [REQ-1] "`[model]` MAY carry `class`, a string whose only admitted values are
  `\"state-machine\"` and `\"decision-table\"`" — (NC, `0010:C1`)

- [REQ-2] "an absent `class` MUST be read as `\"state-machine\"`" — (NC, `0010:C1`)

- [REQ-3] "A `class` outside that set is `malformed model declaration`." —
  (NC, `0010:C1`)

- [REQ-4] "The agreement check is **one-directional**: a `\"decision-table\"`
  model MUST declare zero tags of provenance `owned`, and a non-zero owned count
  under that class is a `malformed model declaration` load failure whose detail
  names the class and the owned count, rendered as the token `owned=<n>`." —
  (NC, `0010:C1`) — the detail assertion is on the literal token, not an
  invented substring (SC-1).

- [REQ-5] "A `\"state-machine\"` model MUST NOT be refused for declaring zero
  owned tags — that model is rootless, and `0006:C18` reports it as
  `graph-dangling-edge` at lint, which is the pre-existing behaviour this RDR
  leaves untouched" — (NC, `0010:C1`) — negative REQ: no second direction of the
  agreement check may be implemented.

- [REQ-6] "The class is model data and is carried on the loaded model; nothing
  downstream infers it from the owned set." — (NC, `0010:C1`) — negative REQ.

- [REQ-7] "The class is carried as an **exported `Class string` field on
  `internal/table/model.go::Model`**, whose zero value — the empty string —
  MUST be read everywhere as `\"state-machine\"`, so a hand-constructed
  `table.Model` is a state machine without naming a class." — (NC, `0010:C1`;
  grounded A15) — an exported field, not an accessor.

- [REQ-8] "A reader MAY spell the zero value through a small helper rather than
  repeating the empty-string comparison, but the field is the storage and no
  reader re-derives the class from the owned set." — (NC, `0010:C1`)

- [REQ-9] "The check MUST run after the tag table is loaded, not in the
  `[model]` header step" ... "The class is read in `loadModelHeader` (it is
  `[model]` data) and the agreement is checked in a step at or after `loadTags`"
  — (NC, `0010:C1`)

- [REQ-10] "The position is fixed at **both** ends: the check MUST run at or
  after `loadTags` and **before `checkAccessorBindings`**, `run`'s last step."
  — (NC, `0010:C1`)

- [REQ-11] "a decision-table model that declares both an owned tag and
  `[initial]` refuses as `malformed model declaration` naming the class and
  `owned=<n>`, never with the accessor-binding writer-arity diagnostic" —
  (NC, `0010:C1`) — the doubly-malformed determinacy case.

- [REQ-12] "Any step in that window satisfies this clause; the RDR fixes the
  window, not the step's name." — (NC, `0010:C1`) — licence: the implementer
  chooses the step name/position within the six available insertion points (A15).

- [REQ-13] "The class is **declared, not inferred** from the empty owned set."
  — (AP) — negative REQ, restating REQ-6 at the design level.

## B. Class-conditioned rule shape and prohibitions (C2)

- [REQ-14] "A decision-table model MUST NOT declare `[initial]`, `terminal`, or
  any accessor whose `keys` name an owned tag, and its ordinary rules MUST NOT
  carry a write block or a clear list." — (NC, `0010:C2`)

- [REQ-15] "None of these is a new refusal: each is already unauthorable once
  the owned set is empty" — (NC, `0010:C2`) — negative REQ: no new refusal arm
  may be added for `[initial]`/`terminal`/write-block on a decision table; the
  existing `unknown tag` / `malformed accessor binding` / `write to non-owned
  tag` arms carry it.

- [REQ-16] "`0002:C4`'s \"an ordinary transition rule MUST contain a write
  block\" is conditioned on class: it binds the `\"state-machine\"` class only,
  and the `malformed rule shape` arm that enforces it MUST NOT fire for a
  `\"decision-table\"` model." — (NC, `0010:C2`) — the one edit C2 requires.

- [REQ-17] "Rule ids, the match-block obligation, shared contexts, guards,
  gates, escape rules, and outcome binding are unchanged in both classes." —
  (NC, `0010:C2`) — negative REQ.

- [REQ-18] "a terminal predicate over a non-owned tag is 0006's
  `graph-dangling-edge` terminal arm, which C5 keeps live for exactly this
  reason" — (NC, `0010:C2`; see REQ-53) — the `terminal` prohibition is enforced
  at lint by the live terminal arm, not by a load refusal.

## C. `[rule.emit]` grammar, normalization, and dump (C3)

- [REQ-19] "Any rule, ordinary or escape, in either class MAY carry
  `[rule.emit]`: a flat TOML table whose values MUST be strings." — (NC,
  `0010:C3`) — both classes, both rule kinds.

- [REQ-20] "`sourceRule.Emit` is typed `map[string]string`, so a non-string
  value — a nested table (`[rule.emit.sub]`) included — is refused by the
  decoder as a type error, which `internal/table/source.go::decodeStrict` maps
  to `malformed TOML`" — (NC, `0010:C3`)

- [REQ-21] "The assertion is on the category, never on upstream message text
  (`0002:C24`), and no hand-written type check is added." — (NC, `0010:C3`) —
  negative REQ: no hand-rolled type validation.

- [REQ-22] "The normalized row carries `Emit []EmitValue`, a **new** row-level
  type `{Key, Value string}` — *not* `TagValue`" — (NC, `0010:C3`)

- [REQ-23] "Duplicate keys are unreachable — a TOML table refuses them in the
  decoder before this code runs — so normalization sorts by key and asserts no
  dedup pass" — (NC, `0010:C3`) — negative REQ: no dedup arm, and SC-3 records
  no duplicate-key case.

- [REQ-24] "`Row.Emit` MUST be carried through
  `internal/table/normalize.go::expand` for every expanded row, which builds
  each row from the seed literal rather than copying `base`." — (NC, `0010:C3`)
  — the Trace's step 4′, flagged "**unwitnessed — new in this pass**"; an
  implementation obligation, not a verified property.

- [REQ-25] "The carried sequence MAY be **shared** across the rows one rule
  expands to — no clone is required." ... "the clause is stated so an
  implementer does not add a defensive copy and read it as contract." — (NC,
  `0010:C3`) — negative REQ: do not add a clone.

- [REQ-26] "`Emit` is never mutated after normalization: it is read by the dump
  column and by the `resolvePayload` join and by nothing else, and both are
  readers." — (NC, `0010:C3`)

- [REQ-27] "Emit keys are NOT tag keys: they are undeclared, uninterpreted,
  compared by exact byte equality, and MUST NOT be matched, guarded, written, or
  read by any accessor" — (NC, `0010:C3`) — a `[rule.match.<emit-key>]` refuses
  `unknown tag` (SC-3, RISK).

- [REQ-28] "Normalization MUST carry the block on the row as a key-sorted
  sequence; an absent block normalizes to an empty sequence, and a **present but
  empty** `[rule.emit]` normalizes to the same empty sequence rather than
  refusing — `emit` keys on length, not on key presence" — (NC, `0010:C3`) — the
  deliberate divergence from the write/clear/gate blocks' "EVEN AN EMPTY ONE".

- [REQ-29] "`emit` MUST join the closed dump column vocabulary (`0002:C19`)
  **appended last**, after `escape`" — (NC, `0010:C3`) — the only insertion that
  leaves every existing column at its existing index.

- [REQ-30] "Its cell is rendered as `key=value` pairs in key order, bracketed
  like `writes` (`[a=x; b=y]`), with the value emitted as the raw authored
  string — **not** through `renderValue`" — (NC, `0010:C3`)

- [REQ-31] "An empty block renders as the empty bracket the `writes` column
  already uses for an empty sequence." — (NC, `0010:C3`)

- [REQ-32] "A `[dump]` column list MUST name it, and the default column set used
  when `[dump]` is absent MUST include it for both classes." — (NC, `0010:C3`) —
  a `[dump]` list omitting `emit` refuses `malformed dump declaration` (FM, SC-3).

- [REQ-33] "`emit` MUST NOT be added to the kernel row; it is table data joined
  back by rule id after selection" — (NC, `0010:C3`) — negative REQ; the kernel
  row (`internal/resolve/resolve.go::Row`) is untouched (TD item 1).

- [REQ-34] "the existing first-match `internal/cli/flow_resolve.go::rowByID` is
  correct as it stands" — (NC, `0010:C3`) — negative REQ: no per-row uniqueness
  is claimed or needed; `rowByID` is not changed.

- [REQ-35] "`emit` on the wire is a JSON object of string→string, keys in byte
  order, HTML escaping disabled like every other payload map (`0005:C1`); in the
  dump it is one column, `key=value` pairs in key order, bracketed like
  `writes`." — (LBD, `0010:D-wire-byte-format`)

- [REQ-36] "`Row.Emit` is sorted by key and the dump renders `key=value` pairs
  in key order (C3), and the payload's `emit` object is serialized with the same
  byte order and HTML-escaping policy every other payload map already uses" —
  (XC, `0010:G-cross-cutting`, Canonical-form / determinism)

- [REQ-37] "`emit` does not join `internal/graphlint/engine.go::Fingerprint`, so
  editing an emit block never changes a finding's identity" — (XC; LBD
  `0010:D-identity`; A7) — negative REQ, asserted by a fingerprint test (SC-4).

## D. Resolve payload (C4)

- [REQ-38] "The `flow resolve` success payload MUST carry `emit`: a JSON object
  of string values, keys in byte order, present as `{}` — never `null`, never
  omitted — when the selected row authored none." — (NC, `0010:C4`)

- [REQ-39] "Its **position is fixed**: `emit` is declared immediately after
  `Gates` in `internal/cli/flow_resolve.go::resolvePayload`, taking that struct
  from thirteen fields to fourteen; the declaration order is the JSON key order
  Go's encoder emits." — (NC, `0010:C4`)

- [REQ-40] "`Gates`, not `Rule`, is `emit`'s predecessor" — (NC, `0010:C4`) —
  the pre-edit order `Model, Revision, Observed, Owned, Readers, Outcome, Rule,
  Gates, Next, Writes, Clear, Escaped, EscapeClass` is not otherwise disturbed.

- [REQ-41] "`next`, `writes`, `clear`, `owned`, and `readers` keep their
  `0005:C1` shapes and are empty over a decision-table model." — (NC, `0010:C4`)

- [REQ-42] "Gate evaluation is **unchanged and prior to the emit join**" ... "a
  deny is `flow-gate-denied` — a refusal, never a payload, so `emit` is never
  computed on a denied selection." — (NC, `0010:C4`) — negative REQ: no gate
  behaviour is added.

- [REQ-43] "`escaped` is an existing payload field this RDR does not move: it
  stays at its declared position between `clear` and `escape_class`." — (NC,
  `0010:C4`)

- [REQ-44] "A plan rescued by an escape row carries that escape row's **own**
  `emit`, joined on the same `Plan.RuleID` path as any other selection" ... "so
  `rowByID` returns the escape row and C3's per-rule join needs no
  escape-specific arm. An escape row authoring no `[rule.emit]` renders `{}`
  like any other." — (NC, `0010:C4`) — negative REQ: no escape-specific join arm.

- [REQ-45] "Text mode is **not** specified by this RDR: `emit` is a payload
  field like any other and renders through the generic payload renderer 0005
  owns (`internal/cli/respond/text.go::writeTextPayload`), which emits one
  path-qualified leaf per line — `emit.<key>: <value>` — and renders an empty
  block as `emit: (none)` via `flatten`'s empty-container arm." — (NC, `0010:C4`)

- [REQ-46] "this RDR MUST NOT introduce a per-verb text special case for `emit`;
  no `Overrides` entry against 0005's text rendering is taken, and Phase 3 does
  not touch that file." — (NC, `0010:C4`) — negative REQ, naming a file that
  must not be edited.

- [REQ-47] "`--artifact` is not required when the invoked reader set is empty —
  a consequence of `0005:C1`'s reader scoping, restated here, not a new rule —
  and an `--artifact` binding for a role no invoked reader needs MUST be
  ignored." — (NC, `0010:C4`; A2) — ignored, not refused.

- [REQ-48] "`--outcome` remains required in both classes." — (NC, `0010:C4`)

- [REQ-49] "`flow next`, `flow read-state`, and `flow set-state` are unchanged
  by this RDR." — (NC, `0010:C4`) — negative REQ.

- [REQ-50] "`flow next` carrying no `emit` is **deliberate, not an
  oversight**" ... "Adding `emit` to the candidate preview is `0010:BR4`,
  rejected." — (NC, `0010:C4`; `0010:BR4`) — negative REQ.

- [REQ-51] "`invokedReaders` and `runReaders` are unchanged — the empty demand
  set already makes `--artifact` unrequired" and "`flow next` is not touched by
  this RDR (its candidate predicate is cli/0011:C1; A11)" — (TD item 3) —
  negative REQ naming the sites that must not change.

## E. ∅-rooted graph lint (C5)

- [REQ-52] "Graph lint over a `\"decision-table\"` model MUST root the
  reachability relation at the empty owned-state node; the reachable set is
  exactly that node, and every selection context is reachable." — (NC, `0010:C5`)

- [REQ-53] "The root is keyed on the declared class read from the loaded model,
  **augmenting** the existing `len(Initial)` test rather than replacing it: a
  state-machine model with an empty `[initial]` MUST still traverse nothing and
  take `0006:C18`'s missing-root finding, so the seeding predicate is
  \"decision-table class OR a non-empty declared root\", never \"class alone\"."
  — (NC, `0010:C5`)

- [REQ-54] "Overlap (invariant 3), coverage and withholding (invariant 4), the
  redundant-row, vacuous-atom, and unreachable-rule advisories, the node
  ceiling, and the product bound MUST run unchanged over the authored **guard**
  dimensions — the ones `internal/guard/product.go::Dimensions` collects from
  `guard.all`/`guard.unless` alone (`0006:C7`, `0003:C13`), not the tags'
  `observed` provenance" — (NC, `0010:C5`)

- [REQ-55] "A `\"decision-table\"` group whose scoped product has **zero
  participating dimensions** MUST take `graph-unprovable-coverage` naming the
  group, and MUST NOT be reported as covered." — (NC, `0010:C5`)

- [REQ-56] "This RDR therefore requires a **fourth reason value**,
  `no-participating-dimension`, appended to that set." — (NC, `0010:C5`; A14) —
  appended to `internal/graphlint/taxonomy.go`'s closed, append-only `reason`
  set; `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` must stay
  exhaustive and green.

- [REQ-57] "`reason` is normative for this fence: an emission carrying `\"\"`
  fails REQ-80 outright, and one borrowing a dimension-scoped member states a
  remedy that does not apply." — (NC, `0010:C5`)

- [REQ-58] "the arm MUST be emitted inside the `len(dims) == 0` branch and MUST
  precede `emitCoverageArms`" — (NC, `0010:C5`;
  `internal/graphlint/coverage.go::checkCoverage`)

- [REQ-59] "Where both would apply, this finding is reported and
  `graph-coverage-closed-by-escape` MUST NOT be" — (NC, `0010:C5`)

- [REQ-60] "Precedence is delivered by **returning** from the class-keyed arm,
  not by ordering alone: emitting and then falling through to
  `emitCoverageArms` yields *both* findings over a zero-dimension group carrying
  a bare escape row, which is the double-report this clause forbids" — (NC,
  `0010:C5`; A13) — the early `return` is normative, not stylistic.

- [REQ-61] "The arm MUST supply its **own message** rather than routing through
  `internal/graphlint/coverage.go::unprovableMessage`, whose `default` limb
  would phrase the remedy as \"declare the domain\"" — (NC, `0010:C5`)

- [REQ-62] "which is why the clause binds the decision-table class only" — (NC,
  `0010:C5`) — negative REQ: an unconditional arm would redden 151
  zero-dimension groups across 37 of 37 loadable fixtures (RISK).

- [REQ-63] "The missing-root arm of invariant 1 (`0006:C18`), dead end (2),
  always-present-owned (5), owned-set-before-match (6), and single-valued state
  are vacuous by construction over this class and MUST NOT emit." — (NC,
  `0010:C5`; A10)

- [REQ-64] "Lint MUST NOT report the class itself as a finding: no finding's
  `Code`, `Element`, `Message`, or `Detail` may name the class or the string
  `decision-table`" — (NC, `0010:C5`; `0010:BR6` rejected) — negative REQ, in
  assertable form.

- [REQ-65] "Invariant 1's **terminal-predicate arm** MUST stay live" ... "so
  silencing the code wholesale would remove C2's enforcement of the `terminal`
  prohibition." — (NC, `0010:C5`) — the two arms of `graph-dangling-edge` are
  distinguished by `element` (`model` vs `terminal[i]`), not by code (A12, SC-4).

- [REQ-66] "Declared-terminal handling (7) proper, and the escape arm of it,
  stay silent by construction: they key on a declared terminal, which C2
  forbids." — (NC, `0010:C5`)

- [REQ-67] "Over a `\"state-machine\"` model every 0006 invariant, `0006:C18`
  included, is unchanged." — (NC, `0010:C5`) — negative REQ, the regression
  guarantee SC-6 sweeps.

- [REQ-68] "Of the four shipped `len(Initial) == 0` sites, exactly two become
  class-aware: `internal/graphlint/reach.go::reach` and `checkDanglingEdge`'s
  missing-root arm, both *augmenting* that test per C5. The other two,
  `checkAlwaysPresentOwned` and `checkUnreachableRules`
  (`internal/graphlint/analysis.go`), keep the bare `len(Initial)` test" —
  (TD item 1; A10; AUTH) — negative REQ: converting either of the latter two to
  the class test would flip it live over every decision table.

- [REQ-69] "Any site asking \"is this a decision table\" reads the `Class` field
  (empty = `state-machine`, C1) and never re-derives it from `len(owned) == 0`."
  — (TD item 1) — the four class readers are `normalizeRule`, `reach`,
  `checkDanglingEdge` root arm, `checkCoverage` zero-dim arm (AUTH).

## F. Cross-cutting / compatibility

- [REQ-70] "`class` is optional and its absence reads as `\"state-machine\"`
  (C1), so every checked-in model loads, lints, and resolves unchanged; the
  agreement check is one-directional, so no existing model can be newly
  refused." — (XC, `0010:G-cross-cutting`)

- [REQ-71] "The one non-silent widening is the dump vocabulary: an explicit
  `[dump]` list must gain `emit` or fail at load (A4, 103 fixtures)" — (XC; A4)
  — every checked-in model or fixture carrying an explicit `[dump]` column list
  is updated in this change.

- [REQ-72] "no concurrency surface is added. The class is read from the loaded
  model, immutable after load; lint and resolve are single-pass over that
  value." — (XC) — negative REQ.

- [REQ-73] "`[model] version` stays the model's own version, per `0002:C3`" and
  "The model grammar carries no version negotiation and this RDR adds none" —
  (XC, Versioning) — negative REQ.

## G. Failure modes and diagnosis (assertable forms)

- [REQ-74] "class/owned-set disagreement → `flow-model-invalid` with a
  `malformed model declaration` finding naming the class and count; a write
  block on a decision-table row → `write to non-owned tag` / `unknown tag`; a
  `[dump]` list omitting `emit` → `malformed dump declaration`." — (FM, Visible)

- [REQ-75] "`intrastate lint --as=json` over a deliberately partial table must
  report a coverage finding, and over a match-discriminated table must report
  `graph-unprovable-coverage` (C5) — neither may exit 0 with `[]`." — (FM,
  Silent (guarded))

- [REQ-76] "drop `class` (the model still loads and lint reports it as a
  rootless machine — `graph-dangling-edge`, exactly as today)" and "nothing
  persistent is written by this RDR." — (FM, Recovery)

## H. Minimum Viable Validation (`0010:MVV`)

- [REQ-77] "Author a decision-table model: `class = \"decision-table\"`, two
  observed enum dimensions of two values each — each declared `single_valued`
  and `required`, and **discriminated by `[rule.guard.all.<key>]` atoms**, with
  `[rule.match.recognized]` binding the outcome — one recognized outcome, three
  ordinary rules each carrying `[rule.emit]`, no `[initial]`, no accessors." —
  (MVV step 1)

- [REQ-78] "`intrastate lint --model m.toml --as=json` → exit 2 with one
  coverage finding naming the uncovered cell (a *positive* finding, proving
  coverage ran with no root declared)." — (MVV step 2)

- [REQ-79] "Add the fourth rule; lint → exit 0, findings exactly `[]` — no code
  from the 0006 taxonomy present (A10)." — (MVV step 3)

- [REQ-80] "Variant: replace the fourth rule with an escape row rescuing
  `no_match` carrying `[rule.emit]` → exit 0 with only the
  `graph-coverage-closed-by-escape` advisory (A9)." — (MVV step 3)

- [REQ-81] "`intrastate flow resolve --model m.toml --outcome <o> --tag a=x
  --tag b=y --as=json` with no `--artifact` → exit 0; `data.rule` is the fourth
  rule's id, `data.emit` equals its authored block, `data.next`, `data.writes`,
  `data.owned` are empty, `data.readers` is `[]`." — (MVV step 4)

- [REQ-82] "`intrastate dump --model m.toml` renders the `emit` column for every
  row with no `[dump]` declared; `flow next --model m.toml --as=json` with no
  `--artifact` → exit 0, every candidate with empty `required` (A11)." — (MVV
  step 5)

- [REQ-83] "Negative controls: the same model with `class` omitted → lint
  `graph-dangling-edge` naming the absent `[initial]` (0006:C18 unchanged); with
  one `provenance = \"owned\"` tag added → load refusal `malformed model
  declaration`; `models/rdr.toml` lints and resolves as before." — (MVV step 6)

- [REQ-84] "End-state: a stateless model lints for coverage and resolves to
  `rule` + `emit` with no owned state, no accessor, and no artifact, while every
  state-machine model behaves as it did." — (MVV)

- [REQ-85] "That rewrite happens in the consumer repo, not here — intrastate
  stays generic and `models/rdr.toml` is untouched" — (MVV, consumer-side
  acceptance) — negative REQ: `models/rdr.toml` MUST NOT be edited, and the
  generic fixture, not the navigator model, is what closes intrastate#zdat.

## I. Testing Strategy — Done clause and licensed diffs

- [REQ-86] "Done = the MVV passes end to end under `make check`, every scenario
  below has a green test, and every diff to a checked-in expectation is
  **licensed** by the mechanical rule." — (TS)

- [REQ-87] "A changed line is licensed when it differs only by one of exactly
  four shapes, and by nothing else: (i) an added trailing ` emit=[…]` dump cell;
  (ii) an added `\"emit\":` payload member; (iii) an added `\"emit\"` member
  inside a `[dump]` `order = [ … ]` list — the 103 fixtures under
  `internal/table/testdata/` (A4); or (iv) an added `\"emit\"` in a
  dump-vocabulary expectation **in Go**, namely
  `internal/table/dump_test.go::TestReq95_DumpColumnVocabularyIsClosedAndVerbatim`"
  — (TS)

- [REQ-88] "Any other changed line is a regression, not a licensed diff." — (TS)
  — negative REQ; this is the scope gate for the whole implementation.

## J. Testing Strategy — scenarios (`0010:S1`–`0010:S6`)

- [REQ-89] "loader table tests over `class` — absent, each admitted value, an
  unknown value, the one disagreement direction (`decision-table` declaring an
  owned tag), and its counterpart control — a `state-machine` model (or one with
  `class` omitted) declaring zero owned tags, which MUST load." — (SC-1,
  `0010:S1`) — Expected: "absent reads as `state-machine`; the unknown value and
  the decision-table-with-owned-tag disagreement refuse `malformed model
  declaration`, the disagreement's detail containing the literal token
  `owned=<n>` ...; the zero-owned state-machine **loads** and is left to lint
  (C1)."

- [REQ-90] "rule-shape normalization over an ordinary rule with no write block,
  in each class." Expected: "`malformed rule shape` for `state-machine`; loads
  with empty `Writes`/`NextTags`/`RequiresOwned` for `decision-table` (C2)." —
  (SC-2, `0010:S2`)

- [REQ-91] "`[rule.emit]` normalization — unordered keys, a non-string value, a
  nested table (`[rule.emit.sub]`), an absent block, a **present but empty**
  block; dump with and without an explicit `[dump]` naming `emit`;
  `[rule.match.<emit-key>]`" — (SC-3, `0010:S3`)

- [REQ-92] "an emit block on an **expanding** rule — a multi-member `in` match
  atom, so `normalize.go::expand` mints several rows from one rule. That case is
  mandatory, not optional" — (SC-3, `0010:S3`) — Expected: "**every** row
  expanded from the `in`-atom rule carries the authored block byte-for-byte,
  asserted on the expanded rows and again through `flow resolve` on a selection
  that lands on a non-first expanded row (C3, A4)."

- [REQ-93] "No duplicate-key case: a TOML table refuses duplicates in the
  decoder, so the case is unauthorable and no dedup arm exists to test." —
  (SC-3, `0010:S3`) — negative REQ.

- [REQ-94] "graph-lint fixture pair (partial → one coverage finding naming the
  cell; complete, closed by a fourth *ordinary* rule → `[]`), the
  escape-\"otherwise\" variant, the class-omitted control, a fingerprint test
  editing only an emit block, a **match-only negative control** ... a
  **two-outcome escape control** ... and a **stray-`terminal` control** — a
  decision table declaring `terminal` over an observed tag." — (SC-4, `0010:S4`)

- [REQ-95] "the match-only control reports exactly `graph-unprovable-coverage`
  naming the group at exit 2 **while incomplete** — a *positive* assertion" —
  (SC-4, `0010:S4`)

- [REQ-96] "the two-outcome escape control reports the unrescued outcome's
  coverage gap, pinning that \"the otherwise row\" is per-outcome (A9)" —
  (SC-4, `0010:S4`; `0002:C5` per-outcome rescue scoping)

- [REQ-97] "the stray-`terminal` control reports `graph-dangling-edge` at
  `element = terminal[0]` — the arm C5 keeps live, distinguishing it from the
  silenced root arm by `element` rather than by code (C2, C5)." — (SC-4,
  `0010:S4`)

- [REQ-98] "`flow resolve` and `flow next` over the fixture with no
  `--artifact`, and with an `--artifact` bound to an uninvoked role, in JSON and
  text modes. The fixture MUST include one selectable row authoring **no**
  `[rule.emit]` block, and the `{}`-never-`null` assertion MUST be taken on that
  row's resolve" — (SC-5, `0010:S5`)

- [REQ-99] "Expected: exit 0; `emit` present (`{}` when unauthored, never
  `null`); `next`/`writes`/`clear`/`owned` empty, `readers` `[]`; the stray
  binding ignored; text renders one path-qualified line per pair
  (`emit.<key>: <value>`) and the line `emit: (none)` for an unauthored block —
  the generic renderer's output, asserted as built, with no per-verb special
  case (C4, A2, A11)." — (SC-5, `0010:S5`)

- [REQ-100] "regression sweep — every checked-in model and fixture
  (`models/rdr.toml` included) under `make check`." Expected: "lint, resolve, and
  dump output changed only by licensed diffs under the Done clause's rule, and
  every other line byte-identical." — (SC-6, `0010:S6`)

- [REQ-101] "The old-binary arm of A8 is **not** in this scenario" ... "A8's
  forward-compatibility claim is discharged by its one-time spike ... and is not
  promoted to a standing CI obligation" — (SC-6, `0010:S6`) — negative REQ: do
  not write a two-toolchain test.

## K. Implementation plan and docs

- [REQ-102] Phase 1 — "make the class a declared, checked fact and let a row
  carry an answer — `class` in `sourceModel`, `emit` in `sourceRule`, the
  agreement check, the class-conditioned rule-shape arm, `Row.Emit`, and the
  `emit` dump column, with every `[dump]`-carrying fixture updated." — (IP)

- [REQ-103] Phase 2 — "`reach` seeds by class, `checkDanglingEdge` keys its root
  arm on class, and the graph-lint fixtures gain a decision-table pair (partial
  → coverage finding; complete → clean) plus the class-omitted control." — (IP)

- [REQ-104] Phase 3 — "`emit` on `resolvePayload`, joined by rule id after
  selection, rendered in text mode; the output contract doc gains the field and
  a decision-table invocation." — (IP)

- [REQ-105] Phase 4 — "a generic decision-table fixture under
  `internal/*/testdata`, `docs/cli-output-contract.md` gaining the `emit` field
  and a decision-table invocation, and the model authoring docs gaining the
  class." — (IP)

- [REQ-106] Phase 4 doc obligation — "it states that a decision table's
  discriminating dimensions are authored as `[rule.guard.all.<key>]` atoms and
  says why (match atoms scope the group and contribute no dimension,
  `0006:C7`), and it documents the escape-row \"otherwise\" idiom (A9)." — (IP)

## L. Naming and rejected shapes (negative REQs)

- [REQ-107] "**Naming** — `class` with values `state-machine` /
  `decision-table`. Rejected: `kind` ..., `type` ..., `stateless = true` ....
  `emit` for the answer block. Rejected: `output` ..., `result`, `answer`." —
  (LBD, `0010:D-naming`)

- [REQ-108] "**Selection / predicate** — unchanged: the kernel's exact-one match
  is the decision table's hit policy (A5); two matching rules is
  `flow-ambiguous-match` in both classes, and an escape row rescues it the same
  way." — (LBD, `0010:D-selection-predicate`; A5) — negative REQ: no first-hit,
  priority, or collect policy.

- [REQ-109] Briefly rejected, all negative REQs: "Answer = row identity only"
  (`0010:BR1`); "Answer = pseudo-owned `next` tag (the PoC's shape)"
  (`0010:BR2`); "`any`-typed emit values" (`0010:BR3`); "`emit` in the `flow
  next` candidate preview" (`0010:BR4`); "`--outcome` optional when the model
  declares one outcome" (`0010:BR5`); "A new lint finding announcing the class"
  (`0010:BR6`). — (BR)

- [REQ-110] "Illustrative — a two-dimension decision table and the resolve
  envelope it yields; tests must not assert these literally." — (IC) — negative
  REQ: the Illustrative Code and the mermaid figure are **not** normative here
  (contrast 0006, whose IC input contract *is* normative).

---

## ASSUMPTIONS

- **ASSUMPTION-1 (REQ-9/REQ-10/REQ-12)** — the agreement check lands as a new
  named step in `internal/table/load.go::run`'s slice, inserted somewhere in the
  window at/after `loadTags` and before `checkAccessorBindings`. C1 explicitly
  leaves the step's name and exact position free ("Any step in that window
  satisfies this clause"), and A15 confirms six insertion points all satisfy it.
  The audit does not fix a name; a later phase may choose one.

- **ASSUMPTION-2 (REQ-7)** — `Model.Class` is a plain `string` field, not a
  named string type. C1 says "an exported `Class string` field", which reads as
  the builtin type; a defined type would still satisfy "exported field" but
  would change the zero-value comparison's spelling at four reader sites for no
  stated benefit.

- **ASSUMPTION-3 (REQ-22)** — `EmitValue` is declared in `internal/table` beside
  `TagValue` (`internal/table/model.go`), since C3 calls it "a **new**
  row-level type" on `Row` and every peer row type lives there. The record does
  not name the file.

- **ASSUMPTION-4 (REQ-30/REQ-31)** — the dump cell's pair separator is `; `,
  matching the `writes` column C3 says it is "bracketed like" and the example
  `[a=x; b=y]` it gives. The exact separator is read from the existing `writes`
  renderer, not invented.

- **ASSUMPTION-5 (REQ-56)** — the `no-participating-dimension` append is made
  directly in `internal/graphlint/taxonomy.go` in this RDR's implementation
  (rather than being blocked on a 0006 re-lock). IP prerequisites record 0006's
  assent as "structural, not solicited" and check it off; A14 demonstrated the
  append green across the suite. Grounded, not a free reading.

- **ASSUMPTION-6 (REQ-55/REQ-61)** — the group-level `graph-unprovable-coverage`
  message names the group and states the remedy "author the discriminators as
  guard atoms", the remedy C5's own prose gives. C5 fixes only that the arm
  "MUST supply its **own** message"; the exact wording is the implementer's,
  asserted on `Code`/`Element`/`reason` rather than on message text.

- **ASSUMPTION-7 (REQ-63/REQ-64)** — "MUST NOT emit" is asserted as an exact
  finding-list equality over the decision-table fixtures (SC-4's "exact finding
  lists against the full taxonomy"), not as a per-code negative unit test at
  each silent site. This matches 0006's own assertion style.

- **ASSUMPTION-8 (REQ-71/REQ-87)** — the 103 `[dump]`-carrying fixtures under
  `internal/table/testdata/` are updated mechanically (append `"emit"` to each
  `order` list, append the ` emit=[…]` cell to each expected dump line) and the
  diff is verified against the four licensed shapes rather than re-derived. A4
  verified the census; the count is evidence, not a target to hit exactly.

- **ASSUMPTION-9 (REQ-77/REQ-105)** — the MVV fixture and the Phase 4 generic
  fixture are the same artifact where practical; the record asks for "a generic
  decision-table fixture under `internal/*/testdata`" and the MVV table is
  generic by construction (no consumer knowledge). SC-5's extra requirement (a
  selectable row authoring **no** emit block) means that fixture is a superset
  of the MVV's three-rule table, not identical to it.

- **ASSUMPTION-10 (REQ-45)** — text-mode assertions are taken on the renderer's
  actual output for the new field, "asserted as built" per SC-5, rather than on
  a hand-written expected string derived from reading `flatten`. The REQ is that
  no special case is added; the exact leaf spelling comes from the generic
  renderer.

---

## QUESTIONS

None blocking.

Four clauses were candidates and each is answered by the record, its evidence,
or predecessor precedent — so none is escalated:

1. **Does the `no-participating-dimension` append require 0006 to be reopened
   before Phase 2 can land?** Two readings (block on a 0006 re-lock vs. append
   under 0006's own append-only licence) would produce materially different
   sequencing. Answered by IP Prerequisites, which checks the item off and calls
   the assent "structural, not solicited" (0006 declares the set "closed,
   **append-only**", and `0006:C17` reserves "closed **at** …" for sets that may
   not grow), plus A14's green-suite demonstration and peer precedent in 0008 and
   0011. Recorded as REQ-56 + ASSUMPTION-5.

2. **Is `emit`'s payload position after `Rule` or after `Gates`?** C4's first
   sentence pairs `emit` with the selection block ("`emit` sits next to the
   selection block"), which reads toward `Rule`; its explicit clause fixes
   `Gates`. Answered inside C4 itself — "`Gates`, not `Rule`, is `emit`'s
   predecessor" — with the reason (displacing `Gates` would rewrite a key order
   this RDR does not own). Recorded as REQ-39/REQ-40; no ambiguity survives.

3. **Which of the four `len(Initial) == 0` sites become class-aware?** A
   plausible reading class-keys all four for consistency. Answered by TD item 1
   and A10 in the same terms: exactly two (`reach`, `checkDanglingEdge`'s root
   arm); `checkAlwaysPresentOwned` and `checkUnreachableRules` keep the bare
   test, and converting either "would flip it live over every decision table".
   Recorded as REQ-68.

4. **Does C5's zero-dimension arm bind both classes?** An unconditional arm is
   the simpler implementation. Answered by C5 ("the clause binds the
   decision-table class only") and quantified by RISK: 151 zero-dimension groups
   across 37 of 37 loadable fixtures, all state-machine class, would redden.
   Recorded as REQ-62.
