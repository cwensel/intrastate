# Recommendation 0019: Owned-state initialization semantics

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 7 (re-lock) parse the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Feature
- **Profile**: foundational — owned-state first-run initialization: `[initial]` as both lint root and runtime bootstrap source, materialized by the `flow init-state` verb (C1); user-facing yes; locks cross-rdr
- **Priority**: Medium
- **Related Issues**: kata `intrastate#7nqv` (1574)
- **Predecessors**: 0002-transition-table-as-reviewable-data,
  0004-accessor-execution-safety-model,
  0005-skill-integration-cli-contract,
  0006-graph-lint-authority-and-guarantees,
  0010-stateless-decision-tables
- **Overrides**: 0005:C1 and 0005:D-naming (the closed `flow` verb
  enumeration gains exactly one verb, `init-state`) — **appended, not
  conditioned**, and additive on its owner's grammar: no existing verb
  changes and nothing is withdrawn. The new verb inherits C1's
  per-verb MUSTs (`respond.ValidateMode` first, `respond.OK`/`Fail`,
  `SilenceErrors`/`SilenceUsage`) unchanged, and the refusal classes
  new to it are recorded in C1's `flow-*` taxonomy. Precedent for
  extending a closed enumeration this way: `0010`'s override of
  `0002:C19` (closed dump column vocabulary gains `emit`), recorded as
  an ordinary additive override (A4).
- **Seam Lineage**: no prior accretion

## Problem Statement

An operator standing up a state-machine flow hits a first-run wall:
lint refuses to certify a model without `[initial]`
(`internal/graphlint/analysis.go::checkDanglingEdge`, "declares no
initial owned state"), yet no runtime path reads it — no `internal/cli`
code references `Model.Initial`, an absent artifact loads as an empty
store (`internal/cli/flowbind/flowbind.go::load`), and `flow next`
over the absent artifact reports every candidate
`unknown[].reason: absent` for exactly the keys `[initial]` would have
assigned. The first run of every state-machine flow needs a manual
`set-state` for values the model already declares. The wall this
record removes is scoped to the FILE-BACKED write carrier (C1 carrier
scope): flows whose bound write accessor is edit-carried (0028) or
command-backed (0025) keep it, and the verb refuses against them
rather than seeding on an emptiness answer it cannot compute.

The decision, made once: is `[initial]` a lint-only reachability
declaration, or also the runtime bootstrap source for owned state —
and if the latter, which carrier delivers it? (a) Implicit
read-fallback (an absent artifact reads as `[initial]`) collides with
0005's REQ-107 guarantee that a cleared key reads back absent — a
cleared owned tag would silently resurrect its initial value, making
"cleared" and "unseeded" indistinguishable — and with 0004:C3's ban on
ambient artifact discovery, since a synthesized state is a fact no
accessor established. (b) An explicit opt-in flag on `next`/`resolve`
extends 0005 REQ-3's closed, normative flag list. (c) A persisting
init step — a `set-state`-family verb that writes `[initial]` through
the declared writers — keeps the read path pure and the kernel's
never-fill rule intact. The clause must also decide the per-key case
(an artifact present but missing one `[initial]` key — a partial
write, or a key added after seeding), and must be class-keyed the way
the lint already is: RDR 0010's `decision-table` class has no
`[initial]` to materialize.

## Critical Assumptions

- **A1 No runtime path reads `Model.Initial`, so making the init verb its
  sole runtime consumer conflicts with nothing.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: repo-wide sweep for `.Initial` returns no
    non-test consumer in `internal/cli` or `internal/resolve`. The
    only non-test readers are lint —
    `internal/graphlint/reach.go::reach` (root seeding, and the
    `len(m.Initial) == 0 && !table.IsDecisionTable(m)` arm),
    `internal/graphlint/analysis.go::checkAlwaysPresentOwned`,
    `::checkUnreachableRules` — and the loader itself,
    `internal/table/load.go::(*loader).loadInitial` (the write site)
    and `::checkAccessorBindings` (which reads it to build the
    `written` set, see A3). Every other hit is a test
    (`lint_gate_0006_test.go`, `mvv_0010_test.go`,
    `internal/table/*_test.go`, `internal/graphlint/*_0006_test.go`).
    The runtime-consumer class is empty.
  - **If wrong**: a hidden reader already assigns `[initial]` different
    runtime semantics and the clause here contradicts shipped behavior.
- **A2 The `set-state` write path can carry model-sourced values, over the
  admission set the two routes SHARE: a `table.TagValue` from
  `Model.Initial` renders to the same canonical wire form `parseWrites`
  produces for the equivalent argv value, so read-back equality
  (REQ-107's value-for-value form) holds for seeded keys. The loader's
  `[initial]` admission set is a proper SUPERSET of the argv route's, so
  init-state seeds from the loader-normalized value rather than by
  transcribing to argv (C1).**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `go test ./internal/cli/ -run 'TestSpikeA2' -v`,
    table-driven over every value kind `[initial]` admits, comparing two
    fresh artifacts byte-for-byte — the strong form, not value-level
    (`evidence/spikes/a2-canonical-form.md`, raw output
    `a2-canonical-form-run.txt`, source `a2-spike-source.go.txt`).
    RENDERING agrees byte-identically on all 9 kinds both routes admit:
    enum, scalar string, bool, int, float, set array (sorted), set with
    duplicates and HTML characters (deduped, escaping off), empty set
    `[]`, single-member set. The spike varied the value's KIND, and one
    value per kind; it did not vary the SPELLING of a value within a
    kind, which is where the numeric routes diverge (below).
    ADMISSION diverges on 3, and this is
    decided behavior, not a defect: `internal/cli/flow_input.go::canonicalValue`
    refuses a bare scalar for a `kind="set"` tag and an array literal
    for a scalar tag per **JDR 0001 §D11** ("a bare scalar for a set
    key, or an array for a scalar key, is `flow-write-invalid`"), whose
    stated ground is argv parsing ambiguity — "the model declares the
    kind, so parsing is unambiguous". A model-sourced seed has no argv
    string to disambiguate, so §D11's refusals do not govern it, and
    §D13's canonical form ("the same byte form the CLI accepts in
    `--write` … so 0004's read-back equality is byte equality and no
    third encoding exists") is preserved by rendering
    through the same kind-dispatched encoder `set-state` reaches
    (`::canonicalSet` for a set tag, the member verbatim for a scalar
    — C1) — the encoder alone carries the byte-equality property; conformance to the declaration is already
    held by the loader (`::loadInitial`'s `conform`) and is not
    re-established at seed time (C1). The third divergence, an empty scalar
    (`note = ""`), is unwritable by the ARGV route today —
    loader-admitted, `canonicalValue`-refused unconditionally — while
    the seed route admits and persists it (`{"note":""}`, reading back
    present); the residual argv asymmetry is recorded as out of scope
    here (see Failure Modes).
    SPELLING, as distinct from kind, is the boundary of the claim and
    is verified against source rather than the spike: on the numeric
    kinds the two routes normalize the same VALUE to different BYTES.
    `internal/table/load.go::valueMembers` renders through `strconv`
    (`FormatInt`, `FormatFloat(_, 'g', -1, 64)`) while
    `::canonicalValue` returns the argv string verbatim, and
    `internal/table/load.go::conformKind`'s `int` arm admits a leading
    `+` (`strconv.Atoi`). Probed: `valueMembers(float64(1.0)) == "1"`
    against `--write threshold=1.0` storing `"1.0"`; `+5` against `5`
    on an `int` tag. (`float` is not a declared kind — the closed
    vocabulary in `internal/table/model.go` is
    `{enum, bool, int, set, scalar}` — so a TOML float reaches the
    seed encoder under `scalar`, which `conformKind` passes through
    unchecked.) This does NOT break the claim this assumption makes:
    each route is internally canonical and read-back equality holds
    within it, which is what the seed path needs. It bounds RT3, whose
    equivalence is between the two routes and therefore cannot bind an
    operator's choice of spelling.
  - **If wrong**: init's read-back mismatches on lint-clean models — or,
    worse, passes while persisting a form a manual `set-state` would not
    — and the verb needs an explicit route through the same coercion
    stage `--write` uses, changing the design's reuse claim.
- **A3 Every `[initial]` key of a model that LOADS is served by EXACTLY
  ONE declared `[write.<id>].keys` writer — existence AND uniqueness —
  so init can route every seed through
  `internal/cli/flow_state.go::writerFor` (which refuses both the
  zero-writer and the multi-writer arm) without a bypass.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the guarantee holds, but its site is the LOADER, not
    lint — the assumption originally looked for a lint arm and there is
    none. `internal/table/load.go::(*loader).checkAccessorBindings`
    builds a `written` set from every key a rule writes or clears AND
    every key `[initial]` assigns (`for _, t := range l.model.Initial {
    written[t.Key] = true }`), then refuses any `written` key whose
    writer count is not exactly 1: `CatMalformedAccessorBinding`,
    "written tag %s is served by %d writers; want exactly one". Its
    own comment states the clause: "Every key any rule's write block or
    clear list names, and every key `[initial]` assigns, MUST be served
    by exactly one writer." Confirmed absent from lint:
    `internal/graphlint` has no non-test reference to `m.Writers` at
    all, and `::checkAlwaysPresentOwned` covers a strict SUBSET
    (owned tags with `decl.Required == true`) and checks state
    presence, not writer coverage. Because this runs at load, any model
    the CLI can see already satisfies it by construction — a stronger
    and earlier guarantee than lint certification, so the assumption's
    scope widens from "lint-certified" to "loads at all", and
    `writerFor` needs no bypass. Matching arms confirmed at
    `internal/cli/flow_state.go::writerFor` (`n == 0` and `n > 1` both
    user errors).
  - **If wrong**: a model refuses at seed time
    (`flow-write-unbound`/multi-writer arm) — caught with zero writes
    committed by C1's plan-level validation, but still a wall; the
    clause then needs a companion lint arm making writer
    existence-and-uniqueness for `[initial]` keys a certification
    requirement (the preferred cure — fail at lint time), or a decided
    direct-write bypass.
- **A4 Extending 0005:C1's closed verb enumeration by one verb is a
  recordable override of a Final peer's clause.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: the assumption's own caution was right — `0010:C5`'s
    override of `0006:C18` is a CLASS-CONDITIONING of a substantive
    rule (it narrows scope; nothing is added to a set) and does not
    reach the enumeration case. The governing precedent is a different
    entry on the same `0010` Overrides line: **`0002:C19` (closed dump
    column vocabulary gains `emit`)** — recorded as an ordinary
    additive override of a Final peer, with no assent clause and no
    special provision. `0002:C19` fixes "That closed list is the column
    vocabulary" and refuses an `order` naming an unknown identifier at
    load, so it is a genuinely closed enumeration extended by exactly
    one member. `0010`'s own line marks the contrasting strict case —
    `0006`'s `reason` set, "**appended, not conditioned** … the one
    clause of this RDR that is not additive on its owner's grammar, so
    it needs 0006's assent (A14)" — which does not govern here: adding
    a verb takes nothing from `0005` and changes no existing verb, so
    it is additive on its owner's grammar, sitting with `0002:C19`.
    Shipped code already anticipates the extension:
    `internal/cli/flow.go:87` — "a fifth verb added inside the group
    inherits the gate". Consumers this override must update, verified
    by grep across the tree and split by WHAT changes, since the two
    classes fail differently. Sites hard-coding the CARDINAL, six:
    (1) `internal/cli/flow.go`'s package comment ("One group, four
    verbs"); (2) `::newFlowCmd`'s `Long` body; (3)
    `::flowExtendedDesc`; (4) `internal/cli/flow_exec.go`'s header
    comment; (5) `docs/cli-reference.md`, two strings mirroring
    flow.go; (6) `docs/model-authoring.md` ("The same four verbs.").
    Sites that gain an ENTRY but state no cardinal: the bare-`flow`
    `command-error` message, which enumerates the verbs verbatim;
    `internal/cli/docs.go` generation and the checked-in `llms.txt`
    (bullet-enumerated — regenerate and commit); the `flow-*` taxonomy
    in `docs/cli-output-contract.md` ("each verb", generic);
    `0005:D-naming`'s verb list; and narrative prose in `README.md`
    and `internal/cli/root.go`. (`internal/cli/flow_state.go`'s "The
    four the contract…" is a cardinal about REFUSALS, not verbs, and
    does not change.) The verb also inherits `0005:C1`'s per-verb MUSTs
    (`respond.ValidateMode` first, `respond.OK`, `respond.Fail`,
    `SilenceErrors`/`SilenceUsage`), which 0005's test corpus asserts
    verb-by-verb; the binding constraint is
    `internal/cli/flow_surface_0005_test.go`'s
    `TestReq1And3_FlowGroupExposesExactlyTheFourNormativeVerbs`, whose
    `flowVerbs` slice is asserted SET-EQUAL to the registered group, so
    a fifth verb fails it until that slice is extended — independently
    of any doc string (Phase 1 Step 1).
  - **If wrong**: the verb cannot be added without re-opening RDR 0005,
    and the carrier falls back to the rejected flag-on-`set-state` form.
- **A5 The empty-store predicate preserves cleared keys under
  unconditional automated re-invocation for as long as the store
  retains at least one key: every later `init-state` — including one in
  a deploy script or CI bootstrap that runs it every time — writes
  nothing, and the cleared key stays absent for every reader AND every
  writer path this RDR introduces. The stated exception is exact:
  clearing the LAST key empties the store and the next init-state
  reseeds it.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `init-state` does not exist yet, so the spike verified
    the SUBSTRATE property the predicate rests on, against the shipped
    CLI and at the Go level, and they agree
    (`evidence/spikes/a5-empty-store.md`; steps
    `a5-step-i-absent.txt` … `a5-step-v-diff.txt`,
    `a5-step-go-load-apply.txt`). A clear is key REMOVAL —
    `internal/cli/flowbind/flowbind.go::Writer.Apply` does
    `delete(s, t.Key)`, no tombstone. Clearing the last key leaves
    exactly `{}` + newline (3 bytes; `xxd` captured), with no header,
    version marker, or retained key names, and
    `::load` returns the same zero-key store for it as for an absent
    file (`diff` of the two `read-state` envelopes exits 0 modulo the
    echoed path; `reflect.DeepEqual` true). A subsequent write over the
    emptied store and over a fresh one are byte-identical (`cmp`). A
    store retaining one key reports `len == 1`, so the predicate
    correctly declines. EMPTIED == ABSENT confirmed; the only surviving
    distinction is that the file exists on disk, which no content-level
    predicate can see. The MVV steps 6–9 remain the standing
    regression once the verb ships.
    Boundary found by the spike and now named in C1: a read-back-SEALED
    artifact (`::sealedKey`, the NUL-prefixed
    `flow.readback-unreachable` marker persisted per `0004:C13`/`C14`,
    REQ-104) whose owned keys are all cleared is a ONE-key store — so
    the count that makes the predicate safe is over STORE keys, not
    owned keys. Such an artifact is already exit-3 unreadable to every
    reader (`::Reader` reports every requested key UNREADABLE, not
    absent, per `0004:C7`), and the seal is dropped by the next
    reachable-locator write.
  - **If wrong**: either some path seeds into a store that still carries
    a key — the resurrect hazard the per-key variant was rejected for —
    or the emptied-store boundary sits somewhere other than "zero keys",
    in which case the predicate cannot be stated in terms of store
    emptiness at all. The predicate, not the disclosure, is the safety
    property, so either way the fix is in the predicate.

- **A6 A store-emptiness answer ("does this bound artifact carry ANY
  key") is obtainable for the file-backed carrier without violating
  0004:C3's no-direct-artifact-access rule, by exactly one of: a new
  `accessor.ReadBinding` capability, an exported `flowbind`
  cardinality probe, or a `read-state`-family surface — and the chosen
  carrier does not disturb the shipped `ReadBinding` contract for its
  existing implementations.**
  - **Status**: Pending
  - **Method**: Source Search + Spike
  - **Evidence**: Opened by the cove lens, which established the gap:
    `internal/accessor/binding.go::ReadBinding.Read` is key-scoped
    (`Read(ctx, art, requested)` returns one `KeyValue` per requested
    key plus an `unreadable` list) and reports no cardinality, and
    `internal/cli/flowbind/flowbind.go::store` is package-private with
    no exported surface returning a key count. C1's predicate is
    therefore not computable on the seam as it ships, which the
    Technical Design's first-data-flow framing had to be corrected for.
    To verify: enumerate the three candidate carriers against
    `::ReadBinding`'s implementers, of which there are exactly TWO —
    `internal/cli/flowbind/flowbind.go::Reader.Read` and
    `internal/cli/cmdbind/cmdbind.go::Reader.Read` — and confirm one
    candidate can answer emptiness for the file-backed carrier with no
    change to the other implementer's behavior. The candidate MUST also
    answer cardinality on a READ-BACK-SEALED store — returning 1, not an
    unreadable short-circuit — since C1's sealed arm is a no-op success
    at exit 0 that depends on the count being obtainable there; a
    candidate inheriting `::Reader.Read`'s seal short-circuit degrades
    that arm to an exit-3 refusal and fails S9. Added at the
    repeatability lens, which found the two dispositions stated without
    the constraint that reconciles them. The carrier gate now
    covers the READ capability as well as the write one (C1, A8), so
    the candidate need only answer emptiness for `flowbind.Reader` —
    a command-backed reader is refused before the emptiness question
    is asked, which removes the need for any candidate to define a
    cardinality answer for `cmdbind`. The EDIT carrier is not
    a third implementer: `::EditWriter` declares `CapWrite` and has no
    `Read` method at all, and the registry constructs it only in the
    write loop. That is independent ground for C1's carrier refusal —
    an edit-carried model cannot answer emptiness because it has no
    read binding, not merely because it has no JSON store.
  - **If wrong**: if no carrier can answer emptiness without either
    breaking `ReadBinding`'s existing implementers or reading the
    artifact directly, the empty-store predicate cannot be implemented
    as specified, and the verb must fall back to a different gate
    (per-key absent over `[initial]`, which C1 rejects for missing
    non-`[initial]` keys including `::sealedKey`) or the approach
    reopens. This is the predicate the whole safety argument rests on,
    so a refutation is approach-level, not editorial.
- **A7 The verb can decide a write accessor's carrier from
  `internal/accessor/model.go::Definition.Binding`'s dynamic type
  alone, from package `internal/cli`, without exporting
  `flowbind::commandBacked` and without re-deriving the carrier from
  `table.Accessor` — and the three constructed binding types are
  mutually exclusive and exhaustive over the registry's branches.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Opened by the 3amigo lens (implementer persona),
    which established the gap C1's carrier-scope clause had: the
    discriminator C1 originally named,
    `internal/cli/flowbind/registry.go::commandBacked`, is UNEXPORTED,
    so `internal/cli` cannot call it and the contract named a check the
    verb could not perform. Confirmed:
    `::Definition.Binding` is an exported field of an exported struct
    holding an exported `Binding` interface, and all three constructed
    types are exported (`flowbind::Writer`, `flowbind::EditWriter` via
    `::NewEditWriter`, `cmdbind::Writer`).
    `::Registry`'s write loop yields EXACTLY those three, with no
    fourth arm and no type shared between two branches. The shape is
    an initializer plus a two-arm `switch`, not a switch with a
    `default:` clause — `var binding accessor.Binding =
    &Writer{Path: acc.Path}` precedes `case len(acc.Edit) != 0:`
    (→ `NewEditWriter`) and `case commandBacked(acc):`
    (→ `&cmdbind.Writer{…}`) — so the JSON writer is the residue, not
    a branch. That distinction is what makes the type-switch sound
    rather than accidentally sound, and it turns on a second fact:
    `::commandBacked` is `len(acc.Command) != 0 || acc.Path == ""`, so
    the residue is reachable only when `acc.Path != ""`. A
    `*flowbind.Writer` carrying an empty `Path` — which would `load("")`
    to an empty store and let the verb seed every key into nothing — is
    therefore UNCONSTRUCTIBLE here, and the arm ordering is commented
    in-source to protect exactly that property. Admitting only
    `*flowbind.Writer` is thus exactly the file-backed set, not a proxy
    for it. (The registry as a whole constructs six binding types; the
    other three are `Reader`/`cmdbind.Reader` and `Gate`/`cmdbind.Gate`
    — all VALUE types, where this loop's three are pointers (A8),
    none of them write bindings, so they are outside this type-switch's
    domain.)
  - **If wrong**: if the registry can yield a fourth write binding type
    or two carriers share one type, the type-switch admits or refuses
    the wrong carrier silently — seeding on an emptiness answer C1
    declares UNDEFINED, which is the failure S10 exists to catch. The
    fallback is exporting a carrier predicate from `flowbind` (a
    `flowbind` API change this RDR does not currently authorize), so a
    refutation is a contract change here, not an implementation detail.
- **A8 The same carrier discrimination holds on the READ side: the
  registry's read-construction yields exactly the VALUE types
  `flowbind.Reader` and `cmdbind.Reader` (not pointers, as on the
  write side), mutually exclusive and exhaustive, so admitting
  only `flowbind.Reader` is exactly the file-backed read set — and
  the reader a needed role resolves to
  (`internal/accessor/model.go::Registry.readerFor`) is reachable from
  `internal/cli` at carrier-gate time, before any accessor runs.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Opened by the critique lens, which established that
    C1's carrier gate type-switched on the WRITE binding alone while
    the emptiness read and the commit-time read-back both travel the
    READ binding. Confirmed so far: readers and writers ARE resolved
    independently — `::Registry` loops `m.Readers` and `m.Writers` in
    separate passes, and `::readerFor` selects by ROLE over `CapRead`
    — so a file-backed writer paired with a command-backed reader on
    one role is constructible and the write-side switch would admit
    it. Also confirmed, and the reason this assumption states the
    spelling: the read loop constructs VALUES where the write loop
    constructs POINTERS — `var binding accessor.Binding =
    Reader{Path: acc.Path}`, reassigned to `cmdbind.Reader{…}` under
    `commandBacked`, with `Read` declared on value receivers in both
    packages, against the write loop's `&Writer{…}` and pointer
    receivers on `Apply`. A read-side case spelled `*flowbind.Reader`
    matches neither constructed type and would refuse every model,
    which is why C1 fixes the spelling rather than leaving it to the
    implementer. To verify: read the registry's READ-construction
    branches as A7's evidence read the write loop, confirming the two
    types are the complete set (no third arm, no shared type, and — as
    on the write side — whether the file-backed reader is a residue
    whose reachability depends on a predicate the verb cannot call); and
    confirm `::readerFor`'s result is available to the gate before
    invocation, since a discrimination the verb can only make after
    running an accessor cannot serve a refusal C1 orders before one.
  - **If wrong**: if the read side has a carrier the type-switch does
    not name, or the reader cannot be resolved before invocation, the
    gate cannot be stated over both capabilities as C1 now requires —
    either the read-side refusal moves later (losing the
    no-accessor-invoked property S10 asserts) or the emptiness answer
    is taken from a reader whose carrier was never checked, which is
    the UNDEFINED case C1 exists to refuse. The fallback is scoping
    the verb to models whose read and write accessors share one
    file-backed carrier declaration, a narrower admission set than
    this RDR currently claims.

## Proposed Solution

### Approach

Decide the fork as: **`[initial]` is both the lint-time reachability root
and the runtime bootstrap source for owned state, and the carrier is an
explicit, persisting init verb** — `flow init-state`, a `set-state`-family
verb that writes the model's `[initial]` assignments through the declared
write accessors with read-back, exactly on the existing
`internal/cli/flow_state.go::groupByWriter` → executor → read-back path
⇒ the read path stays pure and the kernel's never-fill rule is untouched:
no read verb, and no artifact load
(`internal/cli/flowbind/flowbind.go::load`, whose absent-file-is-empty
rule is what makes a first write possible at all), ever synthesizes state.

This is the shape every read peer engine uses: initialization is an
explicit lifecycle step at session start, never a read-time fallback,
and the initial-vs-existing fork is decided whole-snapshot, never as a
per-key merge (see Investigation). Per-key semantics follow:
init seeds if and only if EVERY bound artifact is EMPTY — the
quantifier is ALL, over the whole bound set, not per-artifact (C1) —
and then seeds every `[initial]` key; a non-empty artifact is a no-op
success whose
payload reports the `[initial]` keys the store lacks — so no
re-invocation resurrects a cleared key while the store still carries
one, a torn seed is visible without being silently half-repaired, and
the key-added-after-seeding case is answered by explicit `set-state`
guided by that report. The predicate's boundary is disclosed rather
than hidden: emptying the store by clearing its last key returns it to
the initializable class, because an emptied store and a never-written
one are the same store (C1). Class-keyed like the lint:
a `decision-table` model (`internal/table/model.go::IsDecisionTable`
⇒ the existing class discriminator is reused, no parallel predicate)
has no `[initial]` (`0010:C2` forbids it) and the verb refuses.

### Technical Design

The verb sits entirely on shipped surfaces: model selection and request
assembly via the shared `flow` path (`internal/cli/flow_exec.go`),
mutation routing via `writerFor`/`groupByWriter`, application and
commit-time verification via the accessor executor's write + read-back
(RDR 0004's model — read-back is the only commit check, no cross-writer
atomicity, and this verb inherits both statements verbatim). The first new
data flow is the *source* of the planned writes: `Model.Initial`
(`internal/table/model.go::Model.Initial`, today written by the loader
and read by lint only) instead of `--write`/`--clear` request flags.
That difference in plan source is why this is a new verb rather than a
`set-state` flag — see Load-Bearing Decisions (Naming).

There is a SECOND new data flow, and it does not sit on a shipped
surface: the store-emptiness read C1's predicate gates on. The accessor
seam is key-scoped (`internal/accessor/binding.go::ReadBinding.Read`
answers only the keys handed to it and reports no cardinality) and
`internal/cli/flowbind/flowbind.go::store` is package-private, so
nothing shipped can answer "does this artifact carry any key" to a verb
that is also forbidden to read the artifact directly (0004:C3). Its
carrier is undecided and booked as A6 — the one part of this design not
already on a shipped surface, and the reason the verb cannot be
described as pure composition of existing ones. The predicate's
semantics are fixed by C1 independently of which carrier lands.

#### Normative Contracts

**C1**

```normative
SEMANTICS. `[initial]` is BOTH the lint-time reachability root
(0006:C18 unchanged — a state-machine model declaring none stays a
blocking finding) AND the runtime bootstrap source for owned state.
Materializing it is an EXPLICIT, PERSISTING act: no read verb, no
artifact load, and no accessor read path may synthesize, default, or
fall back to `[initial]` values when a key is absent — a cleared key
reads back absent for every reader (REQ-107 unchanged), and owned
state is assembled only from caller-bound artifacts (0004:C3
unchanged). A `decision-table` model has no `[initial]` (0010:C2) and
therefore no bootstrap to materialize; initialization MUST refuse for
that class rather than succeed vacuously. The refusal keys on
`internal/table/model.go::IsDecisionTable` ALONE, never on
"declares `[initial]` but is decision-table": that state is
UNCONSTRUCTIBLE — a decision-table model carrying `[initial]` refuses
at LOAD, so no such model reaches this verb. Every construction tried
refuses on one of three arms: an `[initial]` key with no writer is
`malformed_accessor_binding` (`written tag <k> is served by 0
writers`); adding a writer for it refuses as
`write_to_non_owned_tag` when the tag is OBSERVED
(`internal/table/load.go`, the non-owned writer guard) and as
`malformed_accessor_binding` (`write <w> names the reserved key <k>`)
when it is RECOGNIZED, the reserved-key guard firing first; and
declaring the tag owned is `malformed_model_declaration` (`[model]
class "decision-table" declares owned=<n>`, `::checkClassAgreement`).
The arm that fires is not load-bearing here — the conclusion is, and
it is what the verb keys on: no admitted decision-table model carries
`[initial]`.
Every decision-table model the loader admits therefore has
`len(Model.Initial) == 0`, and the verb's refusal is a CLASS refusal,
not an emptiness one — an ordinary state-machine model whose
`[initial]` is empty is a different case, already blocked at lint by
0006:C18.

CARRIER. The `flow` group gains one verb, `init-state` — an override
extending 0005:C1's verb enumeration and 0005:D-naming's verb list by
exactly this spelling. It takes the shared selection flags
(`--flow`|`--model`), explicit `--artifact role=path` bindings, and no
write grammar: its planned writes are the model's `[initial]`
assignments. It also carries `--allow-commands`, not by declaring it
but by INHERITING it: the flag is registered once on the `flow` group's
persistent set (`internal/cli/flow.go::newFlowCmd`), so every child verb
takes it structurally and this one cannot decline it. The enumeration
above is the verb's OWN grammar, not the closed set of flags it accepts
— a reader who takes it as closed concludes the ORDERING clause below
preempts a flag the verb does not have. It MUST route every seed through the declared write
accessors with commit-time read-back, on the same
writer-routing/no-cross-writer-atomicity terms as `set-state`
(0005:C1); it MUST NOT write an artifact directly.

Seed values are taken from the LOADER-NORMALIZED `Model.Initial`
assignments, NOT by transcribing them to their `--write` argv spelling:
each value is rendered in JDR 0001 §D13's canonical form, DISPATCHED
ON THE DECLARED KIND, which is what makes read-back equality byte
equality and leaves no third encoding. `Model.Initial` holds every
assignment as a member sequence
(`internal/table/model.go::TagValue`, `Value []string`) while the
write seam carries one string (`internal/resolve/resolve.go::Tag`,
`Value string`), so the seed encoder is exactly the SET arm of
`set-state`'s: for `decl.Kind == "set"` it is
`internal/cli/flow_input.go::canonicalSet` (sort, compact, JSON
array); for every scalar kind it is the single member verbatim,
`members[0]`. `::canonicalSet` alone is NOT the whole encoder — it
takes `[]string` and renders a JSON array, so applying it to a scalar
seed would persist `["draft"]` where `set-state --write status=draft`
persists `draft`, breaking the very byte-identity RT3 asserts; six of
S4's nine kinds are scalars. The argv encoder
`::canonicalValue` is the correct *reference* for both arms' output
but is NOT the call: it re-runs the argv admission checks this clause
rules out below. Conformance to the declaration is ALREADY HELD by the
loader and is not re-established here: `::loadInitial` runs
`conform(decl, "eq", members)` over every assignment, which for a
non-operator is exactly `conformKind` + `conformDomain` per member —
the same body `table.ConformValue` calls. A re-conform pass on a
loader-normalized value cannot fail and buys nothing; the byte-equality
property comes from the encoder alone. This is deliberate, because the loader's
`[initial]` admission set is a proper SUPERSET of the argv route's: a
bare scalar for a set-valued tag and a SINGLE-MEMBER array literal for
a scalar tag both load and normalize — `::loadInitial` refuses
`decl.Kind != "set" && len(members) != 1` as
`CatMalformedInitialDeclaration`, so the scalar-from-array admission is
arity-1 only and a two-member row never reaches the verb — and are
refused only on the argv surface
(`internal/cli/flow_input.go::canonicalValue`) under JDR 0001 §D11,
whose stated ground is that the model declares the kind so argv parsing
is unambiguous — a disambiguation duty a model-sourced seed does not
carry. Transcribing through argv would therefore refuse at seed time
models that load clean, re-creating the first-run wall this verb
exists to remove (A2).

Seeding is ALL-OR-NOTHING over an EMPTY store, never a per-key merge:
init-state seeds if and only if EVERY bound artifact carries NO key,
and then it seeds every `[initial]` key. "Bound" is not a filter that
can shrink the quantifier's domain: the plan-validation clause below
requires every role a needed writer names to BE bound, and refuses
otherwise, so the artifacts quantified over here are exactly the ones
that clause already established. An unbound needed role is therefore a
REFUSAL, never an artifact treated as empty — the reading that would
let a half-bound invocation seed a model on the strength of the roles
it happened to name.
That dependence FIXES AN ORDER, and it is stated here because the
outcome differs: the ROLE-BINDING half of plan validation runs BEFORE
the emptiness read, so an invocation that leaves a needed role unbound
refuses at exit 2 REGARDLESS of what the bound stores contain — it does
not reach the predicate and cannot take the no-op arm. It runs AFTER
the carrier gate, which is likewise observable: an invocation against a
non-file-backed model that ALSO omits a needed artifact flag refuses
under the carrier code, not the artifact-binding family. Carrier before
role-binding is the order because the carrier fact is a property of the
model alone while role-binding is a property of the invocation, and the
refusal that names the model's own defect is the more useful one. The WRITER-ARITY
half is not ordered here at all, having already been settled at load
(below). Only encoding and the write follow the predicate. Without this
order the domain is not established when the quantifier runs, and a
half-bound invocation against a non-empty store would exit 0 as a no-op
while an otherwise identical invocation against an empty one refuses —
the same defect reported under two outcomes. The quantifier is ALL, not
ANY and not per-artifact, and it is load-bearing rather than
stylistic: a model may declare N write accessors over N roles
(`internal/table/model.go::Model.Writers`, a `map[string]Accessor`
which `::runFlowSetState` already iterates per writer, checking
`req.artifacts[role]` for each),
so a per-artifact or ANY reading would seed into a store that still
carries a key whenever a SIBLING artifact happened to be empty —
precisely the resurrect hazard A5 exists to exclude and the ground the
per-key variant was rejected on. Under ALL, one surviving key anywhere
in the bound set blocks the whole seed, so the cleared-key guarantee
below holds per model rather than merely per artifact. The count is over STORE keys, not
owned keys, and the difference is reachable: a read-back-SEALED
artifact carries `internal/cli/flowbind/flowbind.go::sealedKey` (the
NUL-prefixed `flow.readback-unreachable` marker persisted in the
artifact per `0004:C13`/`0004:C14`, REQ-104), so an artifact whose
owned keys have all been cleared but whose last write declared an
unreachable read-back locator is a ONE-key, NON-EMPTY store that
init-state declines to seed. That is the correct arm and it is stated
rather than incidental: such an artifact is already exit 3 unreadable
to every reader — the read accessor reports every requested key
UNREADABLE rather than absent (`0004:C7`) — so seeding into it would
write beneath an unverifiable state. The sealed store is nonetheless a
NO-OP SUCCESS at exit 0 here, not an exit-3 refusal, and the two facts
do not conflict: exit 3 is what a READ-BACK on a sealed artifact
returns, and init-state performs no write and therefore no read-back on
this arm. What it needs is a key COUNT, which the predicate obtains
without reading key values — the seal makes the store non-empty at
cardinality 1 and the verb declines before any write is planned. That a
key-scoped reader short-circuiting on the seal cannot itself supply
that count is precisely A6's undecided carrier: whichever carrier lands
MUST answer cardinality on a sealed store rather than inheriting
`::Reader.Read`'s unreadable short-circuit, or this arm degrades from
no-op success to an exit-3 refusal and S9 fails. The seal is dropped by the next
write whose locator IS reachable, which returns the artifact to the
ordinary classes above.

CARRIER SCOPE of the predicate. Everything stated above about STORE
keys, `::sealedKey`, and artifact BYTES is scoped to the FILE-BACKED
write carrier — the JSON `internal/cli/flowbind/flowbind.go::Writer`.
That is not the only carrier: `::Registry` selects among THREE
(`internal/cli/flowbind/registry.go`), branching on
`len(acc.Edit) != 0` to the line-oriented `::NewEditWriter` (0028) and
on `::commandBacked` to `cmdbind.Writer` (0025), falling back to the
JSON writer. An edit-carried or command-backed write accessor has no
JSON store, no `sealedKey`, and no artifact bytes in the sense A5
verified, so "the store carries no key" is UNDEFINED for it. Against a
model any of whose bound write accessors is non-file-backed,
init-state MUST refuse — a distinct terminal refusal in the `flow-*`
family naming the accessor and its carrier — rather than seed on an
emptiness answer it cannot compute.

The carrier gate is over BOTH capabilities, not the write side alone.
Readers and writers are resolved independently — `::Registry` loops
`m.Readers` and `m.Writers` separately, and read-back resolves its
binding by ROLE over `CapRead`
(`internal/accessor/model.go::Registry.readerFor`) — so a model may
declare a file-backed WRITE accessor and a command-backed READ
accessor on the same role, and the write-side type-switch alone would
admit it. It must not: the emptiness answer and the commit-time
read-back both travel the READ binding, so a non-file-backed reader
leaves "the store carries no key" exactly as UNDEFINED as a
non-file-backed writer does, and admitting the model on its write
carrier would refuse late, under an unrelated code, after the gate
already said file-backed. The refusal therefore fires when EITHER the
bound write binding or the bound read binding for a needed role is not
the file-backed type, and it names which capability and which accessor
failed. S10 gains the read-side row.

The verb detects the carrier by the CONSTRUCTED BINDING'S TYPE, not by
re-reading the accessor declaration: `::commandBacked` is UNEXPORTED
(`internal/cli/flowbind/registry.go`), so the verb's package
`internal/cli` cannot call the discriminator named above; that citation
identifies the registry's branch, it does not name the verb's call. The
verb type-switches on `internal/accessor/model.go::Definition.Binding`
(an exported field holding the exported binding type the registry
selected), on the write and read bindings alike: it admits ONLY
`*flowbind.Writer` on the write side, refusing `*flowbind.EditWriter`
and `*cmdbind.Writer`, and ONLY `flowbind.Reader` on the read side,
refusing `cmdbind.Reader`. The POINTER/VALUE spelling differs between
the two capabilities and is normative here, because a type-switch case
on the wrong one matches nothing: the registry constructs writers as
POINTERS (`&Writer{Path: acc.Path}`, `&cmdbind.Writer{…}`, both with
pointer receivers) and readers as VALUES (`Reader{Path: acc.Path}`,
`cmdbind.Reader{…}`, both with value receivers on `Read`). A read-side
case written `*flowbind.Reader` would therefore refuse EVERY model,
file-backed ones included — the failure is total rather than partial,
so S10's read-side row also serves as its detector. Re-deriving the carrier
in `internal/cli` from `table.Accessor`'s `Edit`/`Command` fields is
BANNED on the same ground the class arm bans re-deriving from
`len(owned)` (Selection / predicate): it would duplicate a
discriminator whose writer is the registry. Exporting `commandBacked`
is not required and is not authorized here.

ORDERING: the CLASS refusal precedes the CARRIER refusal. Both are
pre-invocation, so nothing above orders them, and the combination is
constructible against the code even though it is not against the
loader's admitted models: the class arm keys on `::IsDecisionTable`,
which is answerable from the loaded model alone, while the carrier arm
needs the constructed registry. Class-first is therefore the cheaper
order as well as the stated one — a decision-table model refuses
without a registry being built. The rule matters only for a model that
would trip BOTH arms; since no admitted decision-table model carries
`[initial]` or can bind an owned-tag writer, such a model does not
reach the verb today, and the order is fixed here so an implementer
does not have to re-derive that it cannot. The carrier refusal is
decided AFTER registry construction
(it reads the constructed binding) but BEFORE any accessor is invoked,
and it PREEMPTS the `--allow-commands` refusal
(`internal/accessor/model.go::Registry.AllowCommands`, also exit 2). A
command-backed accessor therefore yields the carrier code whether or
not the opt-in was passed — otherwise the same model would refuse under
two different codes depending on a flag irrelevant to the carrier fact,
and S10's command-backed row would pass on the wrong one. The
preemption is a real ordering, not a tautology, and it holds against
BOTH of the gate's refusal sites: `Registry.AllowCommands` is carried
as state through construction, which always succeeds, and the refusal
fires later either at invocation
(`internal/cli/cmdbind/cmdbind.go::spawn`'s first pre-spawn rung) or,
for an `edit` write whose role reader is command-backed, at the
pre-mutation read-back check in `internal/accessor/executor.go`. Both
sites are past construction; the carrier refusal is at it. So S10's
command-backed row discriminates: a carrier check moved after registry
invocation would reach `::spawn` first and fail that row. Extending the predicate to those
carriers is deliberately out of scope here (see A6, and the successor
noted in Consequences); this clause fixes the boundary so the verb
cannot silently do the wrong thing at it.

EMPTINESS IS A NEW DATA FLOW, and it needs a carrier this RDR does not
yet fix. The shipped accessor seam is strictly KEY-SCOPED:
`internal/accessor/binding.go::ReadBinding.Read` answers only the keys
it is handed and reports no store cardinality, and
`internal/cli/flowbind/flowbind.go::store` is package-private with no
exported surface returning a key count. So "the bound artifact carries
no key" cannot be evaluated by any verb sitting on the seam as it
ships, while this contract simultaneously forbids reading or writing
the artifact directly (0004:C3). The emptiness read is therefore a
SECOND new data flow, additional to the plan's source (Technical
Design). Which carrier
serves it — a new `ReadBinding` capability, an exported `flowbind`
cardinality probe, or a `read-state`-family surface — is UNDECIDED
here and is booked as A6; the predicate's SEMANTICS above are fixed
regardless of which lands, and no implementation may substitute
"every `[initial]` key reads absent" for it, since that is the
per-key variant this contract rejects (it cannot see a non-`[initial]`
key, including `::sealedKey`, and so would seed into a non-empty
store).

A non-empty artifact — torn,
partially seeded, post-clear, or fully seeded alike — is a NO-OP
SUCCESS: zero writes, and the payload reports which `[initial]` keys
the store does not carry (informational, so a torn or post-clear state
is visible without being repaired, resurrected, or failed on).
Consequences fixed here: while the store retains AT LEAST ONE key, a
cleared key is never re-established by init-state under any invocation
pattern, including unconditional automated re-runs (the artifact is
non-empty, so nothing writes). The boundary is exact and is a property
of the store, not of intent: clearing the LAST remaining key empties
the store, and an emptied store is indistinguishable at the content
level from a never-written one — `load` yields the same zero-key store
for an absent file and for one emptied by clears, and a clear is a key
REMOVAL, not a tombstone (`internal/cli/flowbind/flowbind.go::load`,
`::store.Apply`, `0004:C11`). A subsequent init-state therefore RESEEDS
such a store. This is the accepted residual of rejecting tombstones
(see Briefly Rejected); it is disclosed, not designed around, and no
payload or doc may describe init-state as unable to re-establish a
cleared key without this qualification. A key added to `[initial]`
after seeding is re-established by explicit `set-state`, not by init —
the payload's absent-key report names it.
An `[initial]` value that is an EMPTY SCALAR (`note = ""`) seeds like
any other: the scalar arm renders `members[0]`, the artifact carries
`{"note":""}`, and the key reads back PRESENT — distinct from a cleared
key, which leaves the object entirely. It creates no third encoding and
is not refused at seed time, because refusing it would require
re-conforming a loader-normalized value, which this clause rules out
above. The ARGV route's contrary refusal (`::canonicalValue` rejects
`--write note=`) is a loader/write-surface asymmetry predating this RDR
and is NOT resolved here in either direction (F5, routed to RDR 0002) —
what is decided here is only what the seed route does, because this
verb is the route that reaches it.
The payload MUST distinguish the seeded-all case from the no-op case,
and its scope is exactly the `[initial]` key set — the verb claims
nothing about owned keys `[initial]` does not assign. Its CONTENT is
fixed here even though its field NAMES are deferred (Wire / byte
format): the payload reports seeded keys by NAME and does not echo
their values — `read-state` is the surface that reports values (RT1),
and echoing them here would create a second place a seeded value can be
read from, with no invariant tying the two. The absent-key report
belongs to the NO-OP arm, which is the only arm where the set can be
non-empty; on the seeded arm every `[initial]` key was just written, so
an absent list there is necessarily empty and MUST NOT be reported as
if it carried information.

The ENTIRE plan is validated before any write: every `[initial]` key
must route to exactly one declared writer and every needed artifact
role must be bound, or the verb refuses with ZERO writes committed —
per-key commit has no atomicity across writers, so plan-level
validation is where the all-or-nothing property lives. The two halves
differ in REACHABILITY, and only the second is testable at this verb:
the writer-arity half is already enforced at LOAD —
`internal/table/load.go::checkAccessorBindings` folds every `[initial]`
key into its `written` set and refuses `writerCount[key] != 1` as
`CatMalformedAccessorBinding` — so no model reaching the verb can
violate it, exactly as no admitted decision-table model carries
`[initial]`. Such a model refuses as `flow-model-invalid` at load,
never in the writer-routing family, and S5 asserts that ordering. The
verb states the requirement (it MUST NOT route a key to two writers)
but owes no runtime check of its own for it. The
role-binding half IS reachable, since roles are bound per invocation
(`::runFlowSetState` resolves `req.artifacts[def.Accessor.Role]`), and
S6 is its test. A read-back
mismatch is a distinct terminal refusal naming the key as
PRESENT-AND-UNVERIFIED; the store is then non-empty, so a re-run is a
no-op that does NOT repair it and MUST NOT be documented as its
recovery — recovery is an explicit `set-state` (or discarding the
artifact and re-running init). Against a `decision-table` model — as
answered by `::IsDecisionTable`, per the class arm above — the
verb refuses before any accessor runs — a refusal (exit 2), not a
no-op success. Because such a model declares zero owned tags and so
can bind no owned-tag writer, this refusal CANNOT be asserted by "a
writer that fails if invoked"; the before-any-accessor ordering is
asserted instead by the absence of any accessor invocation at all
(no artifact is touched and no accessor process runs), which S7 states
in those terms. It carries a dedicated code in the `flow-*` family,
sharpened here as `flow-init-class-unsupported`; the carrier refusal
above is `flow-init-carrier-unsupported`. Both follow the shipped
`flow-<subject>-<condition>` convention
(`internal/cli/flow_input.go`'s code block). What this clause fixes
NORMATIVELY is the exit GROUP (2 on both) and that the two are
DISTINCT codes and distinct from each other and from the shared
classes — the literal spellings are non-normative on the precedent of
`0028:C1.3`'s `codeWriteEditRefused`, so no test may pin the string.
That precedent is a deliberate CARVE-OUT from the code table's general
rule, not an application of it, and the distinction is load-bearing
because the general rule says the opposite: the table's header comment
in `internal/cli/flow_input.go` states "The spellings are normative and
the group fixes the exit", while `::codeWriteEditRefused`'s own doc
comment states its "SPELLING is this stage's to choose and is
explicitly non-normative — no test may pin the string; what the
contract fixes is the exit GROUP and the `findings[]` carriage" (its
sibling `::codeRequestRefused` carries the same carve-out). So the two
codes this verb adds join the carved-out minority, and the reader who
checks the header alone will conclude the opposite of what this clause
fixes;
S7, S10 and MVV steps 7–8 assert the group and the distinctness.
Shared refusal classes reuse the existing `flow-*` codes unchanged, at
the groups those codes already carry — this verb introduces no new
group for them and may not re-map one. Model selection
(`::selectModel`/`::selectModelPath`), artifact binding and writer
routing (`internal/cli/flow_state.go::writerFor`) are all
`clierr.GroupUserEnv`, exit 2. Read-back splits, and the split is
load-bearing rather than incidental: a read-back that COMPLETED and
disagreed is exit 2 (`GroupUserEnv`), while a read-back that could not
complete — the unreachable-locator and timeout arms — is exit 3
(`GroupEnvUnavailable`), per `internal/cli/clierr/clierr.go::ExitCodeFor`
and the raise sites in `internal/cli/flow_exec.go::accessorFailureOf`.
So the PRESENT-AND-UNVERIFIED refusal above is exit 2, and a seed whose
read-back is incomplete exits 3 with the same partial-write disposition;
both leave a non-empty store, so both make a re-run a no-op that does
not repair. The classes new to this verb get codes recorded in the
0005:C1 taxonomy extension this override carries.
```

Exact payload field names, the refusal-code spelling, and help text are
deliberately deferred to Resolve/Pre-Lock; C1 fixes the semantics and
the carrier, which is the fork this RDR exists to close.

#### Load-Bearing Decisions

- **Identity** — a seeded value's read-back equality is REQ-107's
  form: value-for-value over the keys init planned to seed, canonical
  wire form for set values — the identical equality `set-state`
  already verifies, not a new one (A2 verifies the rendering path).
- **Wire / byte format** — none new; artifacts keep
  `internal/cli/flowbind` shape, and the success payload rides the
  0005:C1 envelope. Field names deferred to implementation, owned here;
  the payload's CONTENT is not deferred — C1 fixes it (key names not
  values on the seeded arm, the absent list on the no-op arm only).
- **Naming** — `init-state`, completing the `read-state`/`set-state`
  family. Rejected: `init` (reads as scaffolding a model file, not
  seeding state), `seed` (outside the `*-state` family the group
  established), and a `--from-initial` flag on `set-state`. The
  ground for that last rejection is NOT "two write-plan sources on one
  verb": `set-state` already has two and already merges them cleanly
  through one admit body (`--plan <file|->`,
  `internal/cli/flow_input.go` `planFlagName`;
  `internal/cli/flow_state.go::parseWrites`, whose duplicate refusal
  names "`--plan`, `--write` and `--clear`"), so that argument is
  refuted by shipped code and must not be relied on. The ground that
  survives is the PREDICATE: `--from-initial` would make `set-state`
  conditional on store emptiness — a verb whose whole grammar is
  "commit the writes I named" would silently write nothing on a
  non-empty store, and its no-op success would be indistinguishable
  from a committed write. init's all-or-nothing empty-store gate is a
  different contract from set-state's unconditional commit, and a flag
  cannot carry it without changing what `set-state` means. REQ-3's
  flag list is a secondary cost, not the reason.
- **Selection / predicate** — per store, not per key: seed (all
  `[initial]` keys) iff EVERY bound artifact carries no key — the ALL
  quantifier, not ANY and not per-artifact (C1: a sibling empty
  artifact must not license seeding over one that still holds a key) —
  counting every
  STORE key, not only the owned ones, so a read-back seal counts
  (`::sealedKey`, C1) — (the store's
  own presence answer; its CARRIER is undecided and booked as A6,
  because the shipped seam
  (`internal/accessor/binding.go::ReadBinding.Read`) is key-scoped and
  reports no cardinality, and
  `internal/cli/flowbind/flowbind.go::store` is package-private — the
  emptiness read is a second new data flow, not a free consequence of
  the existing one; scoped to the file-backed carrier on BOTH
  capabilities, with the other write carriers (`::NewEditWriter`,
  `cmdbind.Writer`) and the command-backed reader (`cmdbind.Reader`)
  alike refused, per C1 —
  whose writer is any prior committed `set-state`/init write; an
  empty store means never-written or emptied by clears, which are
  indistinguishable at the content level — verified byte-level at
  Resolve, A5 — and both read as
  "initializable" — the cleared-vs-unseeded ambiguity is not removed by
  this predicate, it is confined to the single state where the store
  carries nothing at all, since one surviving key blocks seeding;
  C1 records that residual); class: refuse iff
  `internal/table/model.go::IsDecisionTable` — the shipped
  discriminator (its writer is the loader materializing `[model]
  class`, `0010:C1`), never a re-derivation from `len(owned)`.

#### Round-Trip / Inverse Invariants

- `read-state ∘ init-state = [initial] on the empty-store class`:
  after a seeding init over role R, `flow read-state` over R reports
  every `[initial]` key with its declared value — value-for-value in
  canonical form, the REQ-107 equality, not merely exit 0.
- `init-state ∘ init-state = init-state` (idempotence on any store):
  the second run writes nothing, succeeds, and the bytes of EVERY
  bound artifact are unchanged (the byte assertion is scoped to the
  file-backed carrier, the only one this RDR admits — C1).
- Equivalence with the manual path: `init-state` into a fresh
  artifact and the `set-state --write` transcription of the same
  `[initial]` assignments into a second fresh artifact produce
  byte-identical artifacts, for every value kind BOTH routes admit
  AND every VALUE whose two routes carry the same spelling (the A2
  spike's table; catches both the hard read-back divergence and the
  silent canonical-form divergence). Two qualifiers, both exact and
  load-bearing. FIRST, the loader's `[initial]` admission set is a
  proper superset of the argv route's, so the three kinds only the
  loader admits (bare scalar for a set tag, array for a scalar tag,
  empty scalar) have no `--write` transcription to compare against and
  are outside this invariant — they are covered instead by the
  conform-and-canonicalize path C1 fixes, and TWO of the three are
  testable there by read-back equality (S4) — including the empty
  scalar, which the seed route persists as `""` and reads back present
  even though argv refuses it (F5).
  SECOND, on the numeric kinds the two routes agree on the value and
  may DISAGREE on its spelling, so the invariant is scoped to the
  spelling the loader normalizes to. `internal/table/load.go::valueMembers`
  renders a TOML numeric through `strconv`
  (`FormatFloat(_, 'g', -1, 64)`, `FormatInt`), while
  `internal/cli/flow_input.go::canonicalValue` returns the argv string
  VERBATIM after conformance — and `::conformKind`'s `int` arm admits
  a leading `+` via `strconv.Atoi`. So `[initial] threshold = 1.0`
  seeds `"1"` where `--write threshold=1.0` writes `"1.0"`, and
  `--write n=+5` writes `"+5"` where the loader seeds `"5"`. This is a
  divergence in the TRANSCRIPTION, not in either route's canonical
  form: each route is internally consistent and read-back equality
  holds within it. The transcription is the operator's, so the
  invariant cannot bind it — S4 compares the two routes on values
  whose spellings coincide, and A2 records the numeric-spelling
  divergence as the boundary of the byte-identity claim rather than a
  defect in either encoder.
- Cleared-key preservation, on a store that retains a key:
  `init-state ∘ (set-state --clear k)` on a seeded artifact carrying at
  least one key besides k writes nothing, and `read-state` still
  reports k absent — REQ-107's reader guarantee stays untouched, and
  the writer-side resurrect hazard is closed for this class by the
  empty-store predicate; re-establishment is explicit `set-state`. The
  inverse holds at the boundary and is asserted as such: clearing the
  last key empties the store, so the next `init-state` reseeds every
  `[initial]` key (C1's disclosed residual).

#### Mini-checks

Five structural cues fired at pre-lock. Each table is the decision, not
a note about it.

**`authority`** — who answers each question the verb branches on.

| Input / decision | Writer (source of truth) | Readers | Sibling arms | Canonical |
| --- | --- | --- | --- | --- |
| Seed values | loader `::loadInitial` → `Model.Initial` | lint (`graphlint`), this verb | argv `--write` transcription | `Model.Initial` — argv route rejected (C1, A2) |
| Model class | loader `[model] class` (0010:C1) | `::IsDecisionTable` callers | `len(owned)==0` re-derivation | `::IsDecisionTable` — never re-derived |
| Store emptiness | the artifact's own key presence | **no shipped reader** — carrier undecided (A6) | per-key absent over `[initial]` | STORE-key count, incl. `::sealedKey` (C1) |
| Write carrier | registry write loop's constructed type (`::Definition.Binding`) | this verb's type-switch | `table.Accessor`'s `Edit`/`Command` fields; `::commandBacked` (unexported) | the constructed type — admits `*flowbind.Writer` only (C1, A7) |
| Read carrier | registry read construction, resolved by role (`::readerFor`) | this verb's type-switch | the WRITE binding, as a proxy for both | the constructed read type — admits the VALUE `flowbind.Reader` only, never the pointer spelling (C1, A8) |
| Value encoding | kind-dispatched: `::canonicalSet` (set) / `members[0]` (scalar), JDR 0001 §D13 | `set-state`, this verb | loader `conform` (kind/domain only); `::canonicalSet` for every kind | the kind-dispatched pair — the byte-equality source (C1) |
| Writer for a key | loader `::checkAccessorBindings` (exactly-one) | `::writerFor` | none | load-time binding (A3) |
| Exit group of a shared refusal | the raise site's `clierr.Group`, not the code table | `::ExitCodeFor` | the code table's header rule ("spellings normative, group fixes the exit") | the raise site — mismatch 2, incomplete/timeout 3 (C1) |

**`oracle`** — each MVV row fails if the right thing is wrong.

| MVV row | Fails if X is wrong because Y | Negative control |
| --- | --- | --- |
| 1–2 seed empty store | a key is missing/mis-encoded → read-back value comparison, not exit 0 | seed with a hand-transcribed argv value → bytes differ |
| 3 idempotent re-run | a second run writes → artifact bytes/mtime compared | remove the predicate → bytes change |
| 4–5 cleared-key preservation | a clear is resurrected → `read-state` reports k present | per-key variant → k reappears |
| 7–8 decision-table refusal | refusal raised after accessor setup → no accessor ran, artifact untouched | move the class check after routing → artifact created |
| 9 sealed artifact declined | seal not counted → store reads empty and seeds | count owned keys only → seeds into sealed store |
| 10 non-file-backed carrier refused | carrier check missing → seeds on an undefined emptiness answer | remove the check → the edit case attempts a write; gate the WRITE side only → the command-backed-reader row is admitted (A8) |
| S11 ALL quantifier (two roles, one empty) | predicate read per-artifact → seeds the empty sibling | switch to ANY/per-artifact → role B gains its key; every single-artifact row still passes |
| S12 torn seed re-run | re-run repairs a partial seed → A's bytes rewritten | complete-on-re-run → B seeded, F2's no-repair property lost |

**`fidelity`** — the byte/value equalities and where they stop.

| Operation | Invariant | Lossy exemption |
| --- | --- | --- |
| `read-state ∘ init-state` | every `[initial]` key reports its declared value, canonical form (RT1) | — |
| `init-state ∘ init-state` | all bound artifact bytes unchanged (RT2) | file-backed carrier only |
| init vs `set-state` transcription | byte-identical artifacts (RT3) | the 3 kinds only the loader admits — no argv spelling exists, so all 3 are covered by seed-side read-back instead (F5); and numeric values whose two routes carry different spellings (`1.0`/`1`, `+5`/`5` — A2) |
| clear then init | k stays absent while any key remains (RT4) | last-key clear → reseed (disclosed residual) |

**`disposition`** — every input class to its outcome.

| Input class | Exit | Writes | Payload |
| --- | --- | --- | --- |
| Every bound artifact empty | 0 | all `[initial]` keys | seeded-all |
| Any bound artifact non-empty | 0 | zero | no-op + which `[initial]` keys are absent |
| Sealed (one-key) store | 0 | zero | no-op (declines; C1) |
| `decision-table` model | 2 | zero, no accessor run | class refusal code |
| Role unbound (reached at the verb; S6) | 2 | zero (plan-level) | existing artifact-binding `flow-*` family |
| `[initial]` key with ≠1 writer | 2 | zero — the model never loads | `flow-model-invalid` at LOAD, before the verb (S5) |
| Torn (partially seeded, multi-writer) | 0 | zero | no-op + the un-seeded `[initial]` keys; committed ones NOT rewritten (F2, S12) |
| Non-file-backed write carrier | 2 | zero | carrier refusal naming the WRITE capability (C1 carrier scope) |
| Non-file-backed read carrier (file-backed writer) | 2 | zero | carrier refusal naming the READ capability (C1 carrier scope, A8) |
| Empty-scalar `[initial]` value | 0 | the key, as `""` | seeded-all; reads back present, distinct from cleared (F5) |
| Read-back mismatch (completed, disagreed) | 2 | partial, per-key | key PRESENT-AND-UNVERIFIED |
| Read-back incomplete / timeout (could not complete) | 3 | partial, per-key | inherited `flow-write-readback-incomplete`/`-timeout` (C1) |

**`trace`** — the MVV walked stepwise against the assertions in force.

| Step | Assertions in force | Witness |
| --- | --- | --- |
| bind model + artifacts | C1 carrier (route via accessors), class arm | fixture: state-machine model, file-backed writer |
| class check | C1 class arm keyed on `::IsDecisionTable` | decision-table fixture → exit 2, no artifact (S7) |
| carrier check | C1 carrier scope, both capabilities | edit/command-backed writer → exit 2; command-backed reader behind a file-backed writer → exit 2 (A8); pointer/value spelling per capability — writers `*Writer`, readers `Reader` (C1) |
| role-binding validation | C1 all-or-nothing, ordered AFTER the carrier gate and BEFORE the predicate | every role a needed writer names is bound; unbound → exit 2, no predicate (S6); non-file-backed + unbound → the carrier code |
| emptiness read | C1 predicate (ALL bound artifacts, STORE keys), **A6 carrier undecided** | empty: `{}`+newline, 3 bytes (A5 spike); sealed: `len==1` |
| plan encode | C1 all-or-nothing; writer arity already settled at LOAD (A3) | every `[initial]` key routes to 1 writer |
| encode + write | C1 canonical form, kind-dispatched (`::canonicalSet` set / `members[0]` scalar) | RT3 byte-identity vs `set-state`, both arms (A2 spike: 9 kinds) |
| read-back | 0004 read-back-is-the-only-commit-check | mismatch → PRESENT-AND-UNVERIFIED |
| second run | RT2 idempotence | bytes unchanged |

No CONTRADICTION row. Two unresolved cells, both under-specified
mechanisms rather than conflicts between assertions: the emptiness
read's carrier (A6, Pending), and the read-side carrier
discrimination the gate now requires (A8, Pending — opened by the
critique lens, which found the gate stated over write bindings alone).

#### Illustrative Code

Illustrative only:

```
intrastate flow init-state --model flow.toml \
    --artifact state=state.json --as json
# empty store  → ok payload: every [initial] key seeded
# non-empty    → ok payload: zero writes; reports [initial] keys
#                the store does not carry
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Initial assignments in the normalized model | `internal/table/model.go::Model.Initial` (loader-validated, `CatMalformedInitialDeclaration`) | today read by lint only (A1) | Reuse | none — no schema change |
| Class discrimination | `internal/table/model.go::IsDecisionTable` | exported, but every current non-test caller is in `internal/table` or `internal/graphlint` — its doc comment names `graphlint` as the reason it is exported; `internal/cli` would become the first CLI-layer class reader | Reuse (new caller layer) | C1's class refusal keys on it |
| Writer routing | `internal/cli/flow_state.go::writerFor`, `::groupByWriter` | requires exactly one writer per key — guaranteed at load by `internal/table/load.go::(*loader).checkAccessorBindings` for every `[initial]` key (A3) | Reuse | none — no lint arm owed; a model that loads already satisfies it |
| Write + read-back | accessor executor (`accessor.NewExecutor`, `Write`) | no cross-writer atomicity (stated, inherited) | Reuse | C1 inherits `set-state`'s terms verbatim |
| Request assembly / selection flags | shared `flow` path (`internal/cli/flow_exec.go::buildRequest`) | none | Reuse | new verb registers like the four shipped verbs |
| Verb surface | `flow` group (0005:C1 closed verb enum) | closed enumeration | Extend | recorded override adding `init-state` (A4) |

### Decision Rationale

The question, scored (approaches: A lint-only status quo, B implicit
read-fallback, C opt-in ephemeral flag on `next`/`resolve`, D explicit
persisting init verb, empty-store predicate). The final criterion row
is load-bearing and applies to every column including the chosen one:
any write predicate keyed on a key's absence resurrects cleared keys
under unconditional automated re-invocation, which is the same axis
that disqualifies B — so D earns its column only under the empty-store
predicate, and only within the boundary C1 discloses:

| Criterion | A lint-only | B read-fallback | C flag on read verbs | D init verb (empty-store) |
| --- | --- | --- | --- | --- |
| Cleared ≠ unseeded (REQ-107) | holds | breaks — cleared keys silently resurrect on read | ambiguous per invocation | holds — no read path synthesizes |
| Read purity (0004:C3, never-fill kernel) | intact | broken — a fact no accessor established | broken inside read verbs | intact — write path only |
| Prior-art alignment | diverges — every peer materializes initial at session start | diverges — no peer read-time fallback | partial — XState's restore fork is at creation, not per read | aligns — explicit whole-snapshot lifecycle step (stateless ctor, `createActor`, SCXML doc start) |
| First-run wall removed | no | yes | partially — state never persists | yes |
| Blast radius on decided contracts | none | 0004:C3 + REQ-107 both violated, silently | REQ-3 flag list + 0005:C1's "read verbs MUST NOT run write accessors" | 0005:C1 verb enum: one recorded addition |
| Key-added-after-seeding case | manual `set-state` | uncontrolled | uncontrolled | decided — explicit `set-state`, guided by init's absent-key report |
| Reversibility | — | hard: implicitly materialized state everywhere | medium | easy: remove the verb; artifacts stay valid |
| Cleared-key survival under unconditional automated re-invocation | holds (nothing writes) | fails — every read resurrects | fails — every flagged call resurrects | holds while the store retains a key — zero writes; the disclosed exception is a store emptied by clears |

The deciding rows are the first two, prior-art alignment, and the
automation axis: B and C lose on contract violations that are the
seed's own collision analysis (REQ-107 resurrect; 0004:C3 synthesis; C
additionally puts write authority or phantom state inside read verbs),
and A answers the fork honestly but leaves the user outcome unmet
while shipping docs (`docs/model-authoring.md`: "`[initial]` declares
the owned state a model starts from") already promise start-state
semantics the runtime does not deliver. D under the empty-store
predicate is the only column that removes the wall, survives the
automation axis within a stated boundary, and leaves every decided
clause intact except one enumerated, recordable-override addition. D's
costs are counted, not waved past: a new payload shape (which is not
why `--from-initial` was rejected — that rejection stands on the
two-plan-sources-in-one-verb grammar muddle), one explicit `set-state`
for a key added to `[initial]` after seeding, and the emptied-store
residual C1 discloses.

Premortem: switched (hardened) — the init-verb approach survived with
its seeding rule switched to the empty-store predicate; the remaining
findings are folded into C1, A2–A5, and Failure Modes.
Ground-sweep: clean (22 anchors)
Joint-check: clear (12 peers) — context beside the verdict: 0021
mentions `[initial]` solely in its lint-root role (the reachability
relation lint already computes), which C1 leaves unchanged; 0016's
read-back and 0023's `set-state` mentions are closed-0004/0005
vocabulary (reader cardinality and the four-verb list), not this RDR's
initialization decision. Absence arm vacuous: no open peer is `Final`,
and every closed record is `Implemented` (outside the peer set); the
refusal this RDR converts (unconsumed `Model.Initial` at runtime) is
relied on by no peer text found.

## Alternatives Considered

### Alternative 1: Lint-only — keep the status quo and document it

**Description**: Answer the fork with "`[initial]` is a reachability
declaration only"; the first run stays a manual `set-state`
transcription of `[initial]`, and docs stop implying start-state
semantics.

**Pros**:

- Zero contract motion: REQ-3, 0005:C1, 0004:C3, REQ-107 all untouched.
- The kernel's never-fill property needs no new argument.

**Cons**:

- The user outcome (remove the first-run wall) is unmet; every
  state-machine flow's first run stays a hand-typed duplication of data
  the model declares, with transcription drift as the failure mode.
- Diverges from all read prior art — every peer engine materializes the
  declared initial state via some explicit runtime step — and from the
  project's own authoring docs.

**Reason for rejection**: it decides the fork by forfeiting the
problem; acceptable only if every materializing carrier were
contract-breaking, and D is not.

### Alternative 2: Implicit read-fallback

**Description**: An absent artifact (or absent key) reads as its
`[initial]` value inside the read path.

**Pros**:

- No new surface at all; first run just works.

**Cons**:

- A cleared owned tag silently resurrects — "cleared" and "unseeded"
  become indistinguishable, inverting REQ-107's read-back guarantee.
- A synthesized value is a fact no accessor established — 0004:C3's
  ambient-discovery ban in spirit and letter.
- No peer engine does read-time fallback; the restored-vs-initial fork
  is uniformly decided once, at session start.

**Reason for rejection**: breaks two decided contracts silently; the
worst reversibility of the four (implicitly materialized state
everywhere, no record of which values were real).

### Alternative 3: Opt-in flag on `next`/`resolve`

**Description**: e.g. `--assume-initial`: the read verbs treat absent
owned keys as holding `[initial]` values for this invocation
(ephemeral), or persist them (write-through).

**Pros**:

- One flag, no verb-surface change; the caller controls when it
  applies.

**Cons**:

- Ephemeral form: the plan is computed over state that does not exist —
  a following `set-state` read-back and every parallel reader see
  different facts, and the cleared-vs-unseeded ambiguity recurs on
  every invocation that passes the flag.
- Persisting form: read verbs gain write authority, directly against
  0005:C1 ("MUST NOT run write accessors" for `next`/`resolve`).
- Either form extends REQ-3's closed flag list on two verbs rather
  than the verb enumeration once.

**Reason for rejection**: both variants trade a one-time explicit act
for a per-invocation semantic mode inside verbs whose purity is a
decided property.

### Briefly Rejected

- **Per-key absent-only seeding**: a write predicate of "key is
  absent" is the very observable a clear produces, so an automated
  re-run resurrects cleared keys — the read-fallback rejection reason
  relocated to the write path. Its re-run recovery story also silently
  skips a present-and-wrong key after a read-back mismatch.
- **Cleared-key tombstones in the artifact**: the only mechanism that
  distinguishes cleared from never-set — it would let per-key seeding
  be safe, and would also close the chosen predicate's emptied-store
  residual (C1). Rejected because it changes the artifact wire format
  every reader and RDR 0004's read-back comparison touch — a
  cross-verb blast radius out of proportion to a first-run verb. The
  residual is accepted and disclosed instead; revisit this if the wire
  format opens for another reason.
- **`--from-initial` flag on `set-state`**: it would make `set-state`
  conditional on store emptiness, so a named write could silently
  commit nothing and read as success. NOT rejected for carrying a
  second write-plan source — `--plan` already is one; see Load-Bearing
  Decisions (Naming).
- **Seed at artifact creation by an external tool/wrapper**: pushes the
  contract outside the CLI where lint and read-back cannot see it —
  ambient by construction.
- **Auto-init on first `set-state`**: makes an unrelated write
  implicitly materialize every `[initial]` key — implicit again, just
  relocated.

## Context

### Background

Tracked as kata `intrastate#7nqv`, scope-reviewed as RDR-shaped: every
fix arm is a contract-level fork, and superseding a locked RDR clause
is an rdr-seed, not a kata fix (batch precedent `d3h2`). The kata's
`type:bug` label was judged a misclassification at triage — the
status-quo refusal contradicts no decided clause: 0002:C2 defines
`[initial]` as a lint-time graph root, not a runtime default, and
0006:A6 scopes initial/terminal as lint-only inputs. RDR 0010 did not
fix it: its Overrides record confirms 0006:C18 (missing root stays a
blocking finding) for the `state-machine` class. Repro caution for
whoever picks this up: `flow resolve --artifact <role>=/nonexistent`
on the checked-in model yields `flow-no-match`, not
`flow-owned-state-unavailable`, because match-pruning fires before the
owned-state check — demonstrate via `flow next` instead. Priority is
Medium on the blocking axis: the shipped `models/rdr.toml` is a
state-machine model (no `class` key ⇒ `ClassStateMachine`) declaring
three `provenance = "owned"` tags and an `[initial]` block seeding all
three, bound to a file-backed `[write.rdr-status]` — so it is in the
population this verb serves, and the repo's own gate model is a live
blocked first run, not a hypothetical one.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/graphlint/analysis.go::checkDanglingEdge`,
`internal/cli/flowbind/flowbind.go::load`,
`internal/cli/flow_state.go`, `internal/cli/flow_exec.go`,
`internal/accessor/binding.go::KeyValue.Absent`. Governing records:
RDR 0002 (C2), RDR 0004 (C3), RDR 0005 (REQ-3 closed flag list,
REQ-107), RDR 0006 (A6, C18), RDR 0010 (model-class split).

## Research Findings

### Investigation

Prior art was read before enumeration (evidence:
`evidence/research/prior-art.md` — queries, accepted citations,
rejected branches). The class read is uniform: a declared initial
state is materialized by an explicit lifecycle step at session start,
never a read-time fallback. Instance reads: qmuntal/stateless
(`repos/qmuntal-stateless/statemachine.go`, `state-machines` checkout)
takes the initial state as an explicit constructor argument —
`func NewStateMachine(initialState State) *StateMachine` — and its
external-storage form hands current-state supply to an explicit
accessor, never a library fallback; XState v5
(`repos/xstate/packages/core/src/createActor.ts`,
`this._initState(options?.snapshot ?? options?.state)`) decides
restored-vs-initial once at actor creation ⇒ "cleared" vs "unseeded"
is resolved by which explicit path the caller invoked, the exact
distinction option B destroys; SCXML (`repos/scxmlcc/doc/user-manual.md`,
`<scxml initial="hello" ...>`) enters the initial configuration at
document start. On the code side, the shaping constraints are the
shipped write path (`internal/cli/flow_state.go::runFlowSetState` —
writer routing, read-back as the only commit check) and
`internal/cli/flowbind/flowbind.go::load`'s absent-file-is-empty rule,
whose stated purpose is making the first write possible ⇒ the
first-run seam was already designed to be crossed by a *write*, and D
is that write with the model as its plan source.

### Key Discoveries

- **Documented** — no `internal/cli` runtime code reads
  `Model.Initial`; the only non-lint consumers are tests
  (`lint_gate_0006_test.go`, `mvv_0010_test.go`). The fork is genuinely
  undecided in code, not decided-by-accident.
- **Documented** — peer engines uniformly make initialization an
  explicit session-start act (citations above); none defaults absent
  persisted state on read.
- **Documented** — `docs/model-authoring.md` already tells authors
  "`[initial]` declares the owned state a model starts from": the
  user-facing framing is bootstrap semantics, so option A would owe a
  doc retraction, not just a no-op.
- **Documented** — lint already class-keys the root arm
  (`internal/graphlint/analysis.go::checkDanglingEdge` augments
  `len(Initial)` with `!table.IsDecisionTable`), and
  `checkAlwaysPresentOwned` enforces that `[initial]` establishes every
  always-present owned key — so for the always-present subset, seeding
  from `[initial]` is guaranteed complete by certification.
- **Documented** — writer coverage of `[initial]` keys is guaranteed at
  LOAD, not at lint: `internal/table/load.go::(*loader).checkAccessorBindings`
  folds every `[initial]` key into its `written` set and refuses any
  whose writer count is not exactly one. `internal/graphlint` has no
  non-test reference to `m.Writers` at all. The guarantee is therefore
  stronger and earlier than the RDR first assumed — it holds for any
  model that loads, certified or not (A3).
- **Documented** — canonical-form fidelity holds for every value kind
  both write routes admit, verified byte-for-byte on the artifact; what
  diverges is ADMISSION, the loader's `[initial]` set being a proper
  superset of the argv route's, and — within the numeric kinds — the
  SPELLING each route normalizes a value to (`1.0` seeds `"1"`; `+5`
  writes `"+5"`), which bounds RT3 without disturbing either route's
  internal canonical form (A2). The two set/scalar-shape divergences are
  JDR 0001 §D11's decided argv-disambiguation refusals, which do not
  govern a model-sourced seed; the empty scalar seeds and reads back
  present on the model route, while the argv route's refusal of it
  stays an undecided asymmetry routed out of this RDR (A2, Failure
  Modes).

## Trade-offs

### Consequences

- Positive: the first run of a state-machine flow becomes one explicit
  command; `[initial]` values are never hand-transcribed, so seed drift
  disappears.
- Positive: every read path keeps its purity argument unchanged —
  reviewers of `next`/`resolve` never need to reason about
  initialization.
- Negative: the `flow` verb surface grows by one, via a recorded
  override of a Final peer's closed enumeration (A4) — precedent-setting
  for how 0005's surface evolves.
- Negative: first-run UX is two commands (init, then work) where option
  B/C would have been zero/one — the explicitness is bought with a step.
- Negative: init helps exactly once per artifact (empty-store
  predicate); a key added to `[initial]` after seeding is a manual
  `set-state` — the disclosed price of cleared-key safety under
  automated re-invocation.
- Negative: that safety is bounded by store emptiness, not by intent —
  clearing the last owned key returns the artifact to the
  initializable class and a later init reseeds it (C1). Operators who
  clear keys individually to reset must discard the artifact or expect
  the reseed; the boundary is contract text and a test, not a hidden
  edge.
- Negative: the verb serves the FILE-BACKED write carrier only. A model
  binding an edit-carried (0028) or command-backed (0025) write
  accessor is refused, because store emptiness is undefined for those
  carriers (C1, carrier scope). Both are shipped and Implemented, so
  this is a live scope gap, not a hypothetical one: those flows keep
  the manual first-run wall this RDR removes elsewhere. Extending the
  predicate to them is a successor, gated on A6's carrier decision.

### Risks and Mitigations

- **Risk**: a model whose `[initial]` names a writerless
  (or multi-writer) key makes init refuse, re-creating the wall it
  exists to remove (A3).
  **Mitigation**: discharged at Resolve — the case cannot arise. The
  loader already refuses such a model outright
  (`internal/table/load.go::(*loader).checkAccessorBindings`,
  `CatMalformedAccessorBinding`), so a model the CLI can load has
  exactly one writer for every `[initial]` key. No companion lint arm
  is owed and no direct-write bypass is needed; C1's plan-level
  validation stays as defense in depth, guaranteeing any residual
  refusal lands with zero writes committed.
- **Risk**: operators read init as "reset to initial".
  **Mitigation**: the empty-store predicate makes init incapable of
  touching a store that carries any key; reset stays `set-state`'s job.
- **Risk**: an operator clears owned keys one at a time to "start
  over", empties the store, and a scheduled `init-state` reseeds —
  reading as the resurrect hazard the design rejects elsewhere.
  **Mitigation**: this is C1's disclosed residual, not a defect; the
  no-op payload and the docs must state the boundary in store terms
  ("while any key remains"), and MVV step 9 pins it so it cannot move
  silently. Closing it entirely would require artifact tombstones,
  rejected on blast radius (Briefly Rejected).

### Failure Modes

- **Visible**: init on a decision-table model refuses with the C1 class
  code before any accessor runs; a writerless or multi-writer
  `[initial]` key refuses at plan validation with zero writes
  committed; read-back disagreement is a terminal refusal naming the
  key present-and-unverified at exit 2; incomplete read-back exits 3
  with the write possibly applied (inherited `set-state` semantics).
- **Torn (multi-writer)**: no cross-writer atomicity — a failed writer
  can leave the artifact seeded for some keys only. The store is then
  non-empty, so a re-run is a NO-OP whose payload lists the missing
  `[initial]` keys — visible rather than silently skipped — but it does
  NOT repair them: recovery is explicit
  `set-state` of the listed keys, or discarding the artifact and
  re-running init. The same holds for a present-and-wrong key after a
  read-back mismatch: no re-run ever overwrites it, and no payload ever
  calls it correct.
  The two recovery routes are not equally available, and the asymmetry
  is the accepted cost of no-repair. Discarding is only safe when the
  artifact holds nothing but `[initial]` keys: an artifact can also
  carry owned keys `[initial]` does not assign, committed by an
  ordinary `set-state` before or after the torn seed, and discarding
  destroys those too. Where they exist, explicit `set-state` of the
  listed keys is the ONLY correct recovery — which on a wide
  `[initial]` block is the hand transcription this verb exists to
  remove, now on the unhappy path. This is accepted rather than solved:
  a repairing re-run would have to write into a non-empty store, which
  is exactly the per-key merge C1 rejects and the resurrect hazard A5
  excludes, so the safety property and the recovery ergonomics trade
  directly against each other. The payload's absent-key list is what
  keeps the manual route mechanical — it names precisely the keys to
  pass — and tearing requires a multi-writer model plus a writer
  failure, so the exposure is narrow. A repair path that preserved the
  cleared-key guarantee (a seed scoped to the failed writer's role,
  gated on that role's artifact being empty) is a successor's, not
  this record's.
- **Scope gap (disclosed)**: init's claim covers exactly the
  `[initial]` key set (C1); an owned key `[initial]` does not assign
  can still surface `unknown[].reason: absent` in `flow next` after a
  successful init. The payload's fixed scope plus
  `flow next`'s own unknown report are the diagnosis surface; the cure
  is authoring the key into `[initial]` (lint's
  `checkAlwaysPresentOwned` already forces exactly that for
  always-present keys).
- **Residual wall (accepted, explicit)**: a key added to `[initial]`
  after seeding is not materialized by any automatic path — one
  explicit `set-state`, guided by init's absent-key report. Accepted
  as the price of cleared-key safety under automation.
- **Empty scalar (disclosed, out of scope)**: an `[initial]` key whose
  value is an empty scalar (`note = ""`) is admitted by the loader —
  no kind or domain rule rejects it and its arity is 1 — yet refused
  unconditionally by `internal/cli/flow_input.go::canonicalValue` for
  non-set tags ("the tag `note` was given an empty value"), so it is
  unwritable by the ARGV route today (A2's spike). No written source
  decides it: JDR 0001 §D11 settled the empty SET against absent
  (`[]` and `--clear` "stay distinct") but never the empty scalar, and
  the corpora are silent. This is a loader/write-surface asymmetry that
  predates this RDR and is NOT decided here — it is routed to RDR 0002
  as a seed.
  What this RDR MUST decide, because it adds the route that reaches it,
  is what init-state does with one. It ARRIVES: `::loadInitial`'s
  guards are undeclared-tag, `valueMembers` error, the `<clear>`
  sentinel, `conform`, and an arity check — an empty string is the
  one-member sequence `[""]`, so arity passes and `conformKind` has no
  `scalar` arm, so `conform` returns nil. Probed: the model loads with
  `Initial = {Key:"note", Value:[""]}`, and C1's SCALAR arm
  (`members[0]`) hands `""` to `flowbind::Writer.Apply`, whose only
  special case is `accessor.IsClear`; the artifact becomes
  `{"note":""}` and reads back `Absent:false` — PRESENT, distinct from
  a cleared key, which leaves the object entirely. So the earlier
  reading that C1's seed path "assumes it cannot arrive" was wrong:
  it can, and it seeds cleanly.
  This RDR takes that as the DECIDED behavior for the seed route
  — the encoder has a defined form for it (`members[0]`), it
  round-trips, and refusing it at seed time would require re-conforming
  a loader-normalized value, which C1 rules out. It creates no third
  encoding: `""` is what the scalar arm renders for a one-member
  sequence, as for any other. What stays undecided in RDR 0002 is the
  ARGV asymmetry — that `--write note=` is refused while the model
  route admits it — which this RDR does not resolve in either
  direction. S4 covers the seeded empty scalar by read-back rather than
  by cross-route diff, since there is no argv counterpart to diff.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A2's canonical-form spike and
      A3's writer-coverage audit gate the design's reuse claim)

### Minimum Viable Validation

Over a fixture state-machine model declaring `[initial]` with one
always-present owned key and one plain owned key, both writer-served:

1. `flow lint` certifies the model (0006 arms all green).
2. `flow init-state --artifact state=<fresh path>` exits 0; payload
   reports the empty-store seeding arm with both keys seeded.
3. `flow read-state` reports both keys at their `[initial]` values
   (canonical form) — round-trip invariant 1.
4. `flow next` over the same artifact reports candidates with no
   `unknown[].reason: absent` entry for either key — the first-run wall
   is gone.
5. Re-run `flow init-state`: exit 0, zero writes, artifact bytes
   unchanged — invariant 2.
6. `flow set-state --clear <plain key>`; `flow read-state` reports it
   absent — REQ-107 held.
7. `flow init-state` again (the unconditional-automation pattern):
   exit 0, ZERO writes, payload reports the cleared key as an absent
   `[initial]` key, `flow read-state` still reports it absent —
   invariant 4 / A5.
8. `flow init-state` against a decision-table model that LOADS (no
   `[initial]`, no owned tags — see S7) refuses (exit 2)
   with the C1 class code, with no accessor run and no artifact
   touched; against the fixture with an unbound artifact
   role it refuses in the existing artifact-binding family with zero
   writes committed; against a model whose write accessor is
   edit-carried or command-backed it refuses with the carrier code and
   zero writes (S10).
9. The boundary, asserted in the failing direction: clear the
   REMAINING key so the store carries none, then `flow init-state` —
   it seeds every `[initial]` key again. This is C1's disclosed
   residual and the test exists so a future change cannot silently
   move the boundary.
10. The outcome on a SHIPPED model, not a fixture: `flow init-state`
   against `models/rdr.toml` — a state-machine model (no `class` key)
   declaring three owned tags, an `[initial]` block seeding all three,
   and a file-backed `[write.rdr-status]` — bound to a fresh artifact
   seeds exactly `stage`, `status`, `gate_passed` at their declared
   values, and `flow next` then reports no `unknown[].reason: absent`
   for them. Steps 1–9 prove the mechanism on a fixture; this step
   proves the user outcome on the repo's own gate model, whose three
   owned keys span two value kinds (enum, bool) the fixture's two keys
   do not.
   One property of that model is an AUTHORING choice this step must not
   silently depend on: `[tags.gate_passed]` declares `kind = "bool"`
   while `[initial]` assigns it the TOML STRING `"false"`, not the TOML
   bool `false`. Both load — `conformKind`'s `bool` arm parses the
   string, and `valueMembers` returns a Go `string` verbatim where it
   would take its `bool` arm for the other spelling — and both
   normalize to `"false"`, so the step passes either way. It is
   recorded because the two spellings reach the seed encoder through
   DIFFERENT arms: if `models/rdr.toml` is later "corrected" to
   `gate_passed = false`, this step still passes but is exercising the
   bool arm rather than the string one. The step's claim is the user
   outcome, not encoder coverage, so the substitution is harmless — but
   S4's kind table is what covers the arms, and this step must not be
   read as doing so.

### Phase 1: Code Implementation

#### Step 1: Verb skeleton

`newFlowInitStateCmd` beside the four shipped verbs in `internal/cli`:
shared selection/tag/artifact registration, `ValidateMode`/respond
gateway, class refusal via `table.IsDecisionTable` before any accessor.

Registration lands A4's override in the SAME step, because the surface
test fails the moment the verb registers and no later step would
explain why. `internal/cli/flow_surface_0005_test.go` fixes
`flowVerbs = {next, resolve, read-state, set-state}` and asserts SET
EQUALITY against the registered group, so it must gain `init-state`
here; it also asserts `0005:C1`'s per-verb MUSTs verb-by-verb, which
the new verb inherits. Six code/doc sites hard-code the cardinal and
change with it: `internal/cli/flow.go`'s package comment,
`::newFlowCmd`'s `Long` body, `::flowExtendedDesc`,
`internal/cli/flow_exec.go`'s header comment,
`docs/cli-reference.md` (two strings), and `docs/model-authoring.md`.
Regenerate and commit `llms.txt` (`internal/cli/docs.go`); it
enumerates the verbs as bullets and hard-codes no cardinal, as does
`docs/cli-output-contract.md`, so neither needs a cardinal edit —
their change is the new verb's own entry (Activation Step 1).

#### Step 2: Predicate and plan

Gate the carrier on both capabilities (Step 1's class refusal has
already fired); non-file-backed → refuse (exit 2). Then validate that
every role a needed writer names is bound; unbound → refuse (exit 2)
without reading any store. Then read the bound stores;
non-empty → the no-op arm (report absent `[initial]` keys, write
nothing). Empty → plan and encode every `[initial]` assignment
(`writerFor` per key) before any accessor runs. The order is C1's, not
this step's, and the role-binding half precedes the predicate there —
see C1 for why the outcome depends on it.

#### Step 3: Apply and report

Execute per writer with read-back (the `runFlowSetState` execution
shape); assemble the payload distinguishing the seeded arm from the
no-op arm.

### Phase 2: Operational Activation

#### Activation Step 1: Contract and reference docs

`docs/cli-output-contract.md` (verb I/O, refusal codes),
`docs/cli-reference.md`, `docs/model-authoring.md` (the start-state
sentence gains its runtime carrier), `--help-all` text. Two constraints
on that text: it states the FILE-BACKED scope (C1 carrier scope), not
an unqualified "every state-machine flow"; and it names the verb from
the surface where an operator meets the wall — `flow next`'s
`unknown[].reason: absent` report — since BR4 and BR5 close both
automatic-invocation routes, so `init-state` is typed by hand and an
undiscoverable manual `init-state` would only move the wall rather than
remove it.

That second constraint is DOCUMENTATION, and this RDR deliberately
leaves it there rather than making it a payload change. Pointing an
operator from `flow next` to `init-state` in-band would mean adding a
field or a reason-string to another verb's payload, which is an
override of `0005:C1`'s I/O for `next` that this record does not carry
— its override is the verb enumeration, nothing more. The disclosed
consequence is that discovery rests on the reference docs and on the
`--help` text for the group, and an operator who reads neither meets
`unknown[].reason: absent` exactly as before. That is a real residual
of the two-command first run already recorded in Consequences, and the
successor that carries a `next`-side pointer is the right home for it:
the pointer is worth having, but not at the price of an unrecorded
second override here.

## Validation

### Testing Strategy

The MVV sequence above is the acceptance spine and runs as an
integration test over a fixture model; "done" is every one of its
steps green plus the unit coverage below. Coverage goals are stated as
arms, not percentages — each normative arm of C1 owes at least one
test that fails if the arm is removed.

1. **Scenario**: Empty store, every `[initial]` key writer-served.
   **Expected**: exit 0; every key seeded; `read-state` returns each at
   its declared value in canonical form (RT1).
2. **Scenario**: Re-run against the store just seeded.
   **Expected**: exit 0; zero writes; artifact bytes unchanged (RT2);
   payload reports the no-op arm.
3. **Scenario**: Seeded store, one plain key cleared, then re-run.
   **Expected**: exit 0; zero writes; the cleared key still reads
   absent and is listed as an absent `[initial]` key in the payload
   (RT4 / A5).
4. **Scenario**: Fresh artifact seeded by `init-state` versus a second
   fresh artifact written by the equivalent `set-state --write`
   transcription, table-driven over every value kind BOTH routes admit
   (9 kinds: enum, scalar string, bool, int, float, set array, set with
   duplicates and HTML characters, empty set, single-member set). Each
   row's `--write` argv MUST use the spelling the loader normalizes to
   (`1`, not `1.0`; `5`, not `+5`) — the numeric rows otherwise diff on
   a transcription difference RT3 does not claim, and a test written
   the other way fails on day one for the wrong reason (A2).
   **Expected**: the two artifacts are byte-identical (RT3; this is
   A2's spike promoted to a standing test). One ADDITIONAL row asserts
   the boundary rather than the invariant: `[initial] threshold = 1.0`
   seeded, against `--write threshold=1.0`, produces artifacts that
   DIFFER (`"1"` vs `"1.0"`) and both read back their own value — the
   divergence is pinned as decided behavior, so a later change that
   silently unified the two encoders would be caught. Separately, the 2 kinds
   only the loader admits (bare scalar for a set tag, array for a
   scalar tag) seed successfully via the normalized path and read back
   value-for-value — they have no argv counterpart to diff against, so
   the assertion is read-back equality, not cross-route byte identity.
   The empty scalar is the loader's THIRD such divergence and IS a row
   here on the same terms: `[initial] note = ""` seeds, the artifact
   carries `{"note":""}`, and the key reads back PRESENT with an empty
   value — asserted distinct from a cleared key, which reads absent
   (F5). Argv refuses `--write note=`, so this row too is read-back
   only.
5. **Scenario**: `[initial]` names a key with zero declared writers,
   and separately one with more than one.
   **Expected**: refusal at LOAD (`flow-model-invalid`), BEFORE the
   verb's own plan validation runs — NOT the writer-routing family,
   which this input never reaches, because
   `internal/table/load.go::checkAccessorBindings` refuses both arities
   as `CatMalformedAccessorBinding` for any `[initial]` key (C1, A3).
   The scenario exists to pin that ordering — a regression moving the
   arity check into the verb would surface a writer-routing code here —
   and asserts ZERO writes committed: artifact bytes/mtime unchanged,
   or the artifact path still absent.
6. **Scenario**: A required artifact role is unbound.
   **Expected**: refusal in the existing artifact-binding family —
   unlike S5 this arm IS reached at the verb, since roles are bound per
   invocation (`::runFlowSetState` resolves
   `req.artifacts[def.Accessor.Role]`) — zero writes committed,
   asserted on artifact bytes/mtime unchanged, or the artifact path
   still absent.
7. **Scenario**: `decision-table` model (one that LOADS — necessarily
   with no `[initial]` and no owned tags; a decision-table model
   carrying `[initial]` is unconstructible, refusing at load on
   `malformed_accessor_binding` / `write_to_non_owned_tag` /
   `malformed_model_declaration`, so it cannot be the fixture here).
   **Expected**: exit 2 with the C1 class code, raised before any
   accessor runs. The ordering CANNOT be asserted by "a writer that
   fails if invoked" — the fixture declares zero owned tags, so no
   owned-tag writer can be declared on it. Assert instead that no
   accessor process ran and no artifact was created or modified
   (bytes/mtime unchanged, or the artifact path still absent), which
   is the discriminating control: a refusal raised AFTER accessor
   setup would fail it.
8. **Scenario**: A writer whose read-back disagrees with the value
   written. On the file-backed carrier writer and reader share one
   store (`flowbind::Writer.Apply` saves what `flowbind::Reader.Read`
   loads), so a disagreement is NOT constructible by writing alone; the
   fixture binds the READ accessor to a different artifact path than
   its writer, pre-seeded with a conflicting value for the same key.
   That shape is legal, which is worth stating because the fixture
   looks like it should refuse at load:
   `internal/table/load.go::checkAccessorBindings` counts READERS and
   WRITERS per KEY (an owned tag owes exactly one of each) and
   constrains no relationship between a reader's `Path` and its
   writer's — nothing requires the two to name one artifact. So one
   reader and one writer over the same key, on differing paths, loads
   clean and is the construction this scenario needs.
   (A command-backed READER is NOT an alternative construction for this
   scenario: C1's carrier gate covers the READ capability too, so such a
   model refuses at exit 2 before any read-back runs and can never
   produce the completed-and-disagreeing read-back this fixture needs.
   That is S10's read-side row, not this one.)
   **Expected**: terminal refusal naming the key present-and-unverified,
   distinct from the read-back-INCOMPLETE class an unreachable locator
   raises; a subsequent re-run is a no-op that does NOT repair it —
   asserted, since C1 forbids documenting the re-run as recovery.
9. **Scenario**: A read-back-sealed artifact whose owned keys have all
   been cleared. The seal is a property of the LAST write's declared
   locator, not of the artifact's history — `flowbind::Writer.Apply`
   sets `::sealedKey` under `if unreachable(w.Path)` and DELETES it in
   the else arm, and `::unreachable` is a pure test on the writer's
   declared `Path` — so a given writer is uniformly sealing or never
   sealing, and clearing through a sealing writer both re-seals and
   exits non-zero (`::Reader.Read` short-circuits on the seal and
   reports every key unreadable). The fixture therefore clears the LAST
   owned key THROUGH the sealing writer itself: `Apply` deletes the key
   and sets the seal in ONE atomic `save`, so the mutation lands even
   though the invocation then exits non-zero (the read-back re-read
   short-circuits on the seal). That non-zero exit is an expected SETUP
   step, not a failure of the scenario — the persisted bytes are what
   the scenario needs, and they are exactly `{sealedKey}`. The end
   state MUST carry NO owned key: a store left holding an owned key
   beside the seal would decline for the ordinary non-empty reason, and
   the owned-key-count control below would not run. (A6's carrier, once
   fixed, may instead permit constructing this state directly; the
   CLI-only route is the one stated here.)
   **Expected**: NO-OP SUCCESS with zero writes — the store carries the
   seal and is therefore non-empty, even though it holds no owned key.
   Asserted on artifact bytes, not just exit code. This pins the
   predicate's count to STORE keys; a future change that counted owned
   keys instead would seed beneath an unverifiable state and fail here
   (A5, C1).
   **Coupling, recorded because it is a liability, not a guarantee**:
   every mechanism this construction rides is PRIVATE to `flowbind` —
   `::unreachable` and its suffix vocabulary, `::sealedKey`'s spelling,
   `::store`/`::load`/`::save`, and `Apply`'s one-save seal-plus-mutate
   ordering are all unexported, and that package's own comment records
   the artifact format as an IMPL-DECISION which "RDR 0005 fixes the
   CLI contract, not a file schema". Only `Reader`/`Writer`/`Gate`,
   their seam methods, and the two format properties (key PRESENCE
   distinguishes empty-set from cleared; values stored VERBATIM) are
   contract. So a `flowbind` change to the seal representation breaks
   this fixture's SETUP, and a fixture that breaks in setup can be
   quietly deleted or skipped — taking with it the only assertion that
   the predicate counts store keys. Two obligations follow. The test
   MUST fail loudly rather than skip if its setup does not produce
   exactly `{sealedKey}` — assert the post-setup bytes before the
   scenario runs, so a representation change surfaces as a failure
   naming this coupling. And when A6's carrier lands, if it can
   construct the sealed state directly, this fixture MUST be rewritten
   onto it: the CLI-only route is stated here because it is the only
   one available today, not because it is the one this scenario wants.
10. **Scenario**: A model whose bound accessor is NOT file-backed on
   one of its two capabilities — run three times over an otherwise
   S1-shaped fixture. Twice on the WRITE side: once with an
   edit-carried accessor (`[write.x]` with an `edit` block, routed to
   `flowbind::NewEditWriter`) and once with a command-backed one
   (`[write.x]` with `command = [...]`, routed to `cmdbind.Writer` via
   `flowbind::commandBacked`). Once on the READ side: a file-backed
   `[write.x]` paired with a command-backed `[read.x]` on the same
   role, which the write-side type-switch alone would ADMIT — the row
   that fails if the carrier gate is written over writers only. S1
   itself is this row's second control: readers are constructed as
   VALUES and writers as POINTERS (C1, A8), so a read-side case
   spelled `*flowbind.Reader` matches nothing and refuses every model
   — S1 goes red, not just this row.
   **Expected**: exit 2 with the carrier refusal code, naming the
   accessor and its carrier, with ZERO writes and no accessor
   invocation. This is the only terminal arm guarding an emptiness
   answer C1 declares UNDEFINED, so its failure mode is seeding on a
   wrong answer rather than a visible error: the discriminating control
   is that removing the carrier check makes the edit case attempt a
   write. For the edit carrier the refusal is over-determined — that
   type implements no `Read` at all (A6) — and the test must still see
   the carrier code, not a missing-reader error. The command-backed
   rows are run with `--allow-commands` UNSET and assert the carrier
   code rather than the allow-commands refusal; no helper binary is
   needed, since no accessor runs. The read-side row asserts the
   refusal names the READ capability and its accessor, distinguishing
   it from the write-side rows' payload — the assertion that would
   fail against a gate reading only `Definition.Binding` on writers.
11. **Scenario**: A model declaring TWO write accessors over two roles
   (`Model.Writers`, which `::runFlowSetState` already iterates per
   writer), each serving its own `[initial]` key. Role A's artifact
   carries a key; role B's is empty. Run `init-state`.
   **Expected**: NO-OP SUCCESS with ZERO writes to BOTH artifacts —
   asserted on both artifacts' bytes, not exit code alone. This is the
   only test of C1's ALL quantifier: under an ANY or per-artifact
   reading the verb seeds role B, resurrecting into a model that still
   holds a key, which is exactly the A5 hazard. The discriminating
   control is that switching the predicate to ANY or per-artifact makes
   role B gain its `[initial]` key here while every single-artifact
   scenario still passes.
12. **Scenario**: A torn seed on the S11 two-writer fixture — role A's
   write commits and role B's fails, leaving the model seeded for A's
   `[initial]` keys only. The torn state is built by writing role A's
   artifact directly during SETUP, never through the verb: no
   declarative lever produces it, since a `-unreachable` writer applies
   its tags before sealing and would leave role B seeded rather than
   un-seeded. Re-run `init-state`.
   **Expected**: NO-OP SUCCESS whose payload names exactly B's
   un-seeded `[initial]` keys, with A's committed values NOT rewritten
   (asserted on A's artifact bytes/mtime). This is the only test of F2:
   an implementation that "helpfully" completes a torn seed on re-run
   passes every other scenario, because on a single-writer fixture torn
   and post-clear are indistinguishable. C1's ban on documenting the
   re-run as recovery is what this asserts.

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

- Peer records: RDR 0002 (C2, `[initial]` as lint-time graph root),
  RDR 0004 (C3, ambient-discovery ban; read-back as the only commit
  check), RDR 0005 (C1 verb enumeration and D-naming; REQ-3 closed
  flag list; REQ-107 read-back), RDR 0006 (A6, C18), RDR 0010 (C1
  model class, C2 no `[initial]` on `decision-table`, C5 override
  precedent).
- Source paths reviewed: `internal/graphlint/analysis.go`
  (`checkDanglingEdge`, `checkAlwaysPresentOwned`),
  `internal/cli/flowbind/flowbind.go` (`load`, `store`),
  `internal/cli/flow_state.go` (`runFlowSetState`, `writerFor`,
  `groupByWriter`), `internal/cli/flow_exec.go` (`buildRequest`),
  `internal/table/model.go` (`Model.Initial`, `IsDecisionTable`),
  `internal/accessor/binding.go` (`KeyValue.Absent`).
- Project docs: `docs/model-authoring.md` (the start-state sentence),
  `docs/cli-output-contract.md`, `docs/cli-reference.md`.
- Prior art searched (full log:
  `evidence/research/prior-art.md`): qmuntal/stateless
  (`statemachine.go::NewStateMachine`), XState v5
  (`packages/core/src/createActor.ts::_initState`), SCXML
  (`scxmlcc/doc/user-manual.md`, `<scxml initial=...>`).
- Related issues: kata `intrastate#7nqv` (1574).
