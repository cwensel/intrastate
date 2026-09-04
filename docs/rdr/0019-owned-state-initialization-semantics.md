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
- **Profile**: foundational — one contract (owned-state first-run
  initialization, C1/C2) whose clauses condition or extend surfaces
  across RDRs 0002/0004/0005/0006/0010.
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
- **Priority**: Low
- **Related Issues**: kata `intrastate#7nqv` (1574)
- **Predecessors**: 0002-transition-table-as-reviewable-data,
  0004-accessor-execution-safety-model,
  0005-skill-integration-cli-contract,
  0006-graph-lint-authority-and-guarantees,
  0010-stateless-decision-tables
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
`set-state` for values the model already declares.

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
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: grep over `internal/cli` and `internal/resolve` for
    `Initial` hits tests only (`lint_gate_0006_test.go`,
    `mvv_0010_test.go`); re-verify at Resolve.
  - **If wrong**: a hidden reader already assigns `[initial]` different
    runtime semantics and the clause here contradicts shipped behavior.
- **A2 The `set-state` write path can carry model-sourced values: a
  `table.TagValue` from `Model.Initial` renders to the same canonical wire
  form `parseWrites` produces for a request value, so read-back equality
  (REQ-107's value-for-value form) holds for seeded keys.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: needed — table-driven over EVERY value kind `[initial]`
    admits: seed each through
    `internal/cli/flow_state.go::groupByWriter` and the executor into a
    fresh artifact, `set-state --write` the equivalent argv value into a
    second fresh artifact, assert the two artifacts byte-identical.
    `--write` values enter as argv strings through `parseWrites`
    coercion while `Model.Initial` values arrive loader-typed, so the
    premise that both reach one canonical wire form is a
    normalization premise — verified, not assumed.
  - **If wrong**: init's read-back mismatches on lint-clean models — or,
    worse, passes while persisting a form a manual `set-state` would not
    — and the verb needs an explicit route through the same coercion
    stage `--write` uses, changing the design's reuse claim.
- **A3 Every `[initial]` key of a lint-certified state-machine model is
  served by EXACTLY ONE declared `[write.<id>].keys` writer — existence
  AND uniqueness — so init can route every seed through
  `internal/cli/flow_state.go::writerFor` (which refuses both the
  zero-writer and the multi-writer arm) without a bypass.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: needed — lint certifies REACHABILITY of `[initial]`
    keys, not WRITABILITY, and no cited clause ties `[initial]` (or
    `graphlint` `ownedRequiredKeys`) to writer coverage: the veto
    clause this assumption needs may simply not exist. Find and quote
    the governing lint/loader clause, or establish that none exists.
  - **If wrong**: a lint-certified model refuses at seed time
    (`flow-write-unbound`/multi-writer arm) — caught with zero writes
    committed by C2's plan-level validation, but still a wall; the
    clause then needs a companion lint arm making writer
    existence-and-uniqueness for `[initial]` keys a certification
    requirement (the preferred cure — fail at lint time), or a decided
    direct-write bypass.
- **A4 Extending 0005:C1's closed verb enumeration by one verb is a
  recordable override of a Final peer's clause.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: needed — an override's EXISTENCE (`0010:C5`'s
    override of `0006:C18`, per `0010:A10`) does not show that an
    override may extend a closed ENUMERATION; citing it for this case
    borrows authority it may not carry. Resolve must quote the override
    provision's actual text against the enumeration case, and enumerate
    the enumeration's consumers this override must update — `flow`
    command registration, `--help-all`, docs generation, and the
    `flow-*` code taxonomy for the verb's new failure classes.
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
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: needed — MVV steps 6–8 (clear one key of a two-key
    seeded artifact, read-back absent, re-init is a no-op success whose
    payload reports the key absent-from-`[initial]` and whose artifact
    bytes are unchanged), AND step 9, which clears the remaining key and
    asserts the reseed actually happens — the boundary is verified in
    the direction that fails, not only the direction that holds.
  - **If wrong**: either some path seeds into a store that still carries
    a key — the resurrect hazard the per-key variant was rejected for —
    or the emptied-store boundary sits somewhere other than "zero keys",
    in which case the predicate cannot be stated in terms of store
    emptiness at all. The predicate, not the disclosure, is the safety
    property, so either way the fix is in the predicate.

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
init seeds if and only if the bound artifact is EMPTY, and then seeds
every `[initial]` key; a non-empty artifact is a no-op success whose
payload reports the `[initial]` keys the store lacks — so no
re-invocation resurrects a cleared key while the store still carries
one, a torn seed is visible without being silently half-repaired, and
the key-added-after-seeding case is answered by explicit `set-state`
guided by that report. The predicate's boundary is disclosed rather
than hidden: emptying the store by clearing its last key returns it to
the initializable class, because an emptied store and a never-written
one are the same store (C2). Class-keyed like the lint:
a `decision-table` model (`internal/table/model.go::IsDecisionTable`
⇒ the existing class discriminator is reused, no parallel predicate)
has no `[initial]` (`0010:C2` forbids it) and the verb refuses.

### Technical Design

The verb sits entirely on shipped surfaces: model selection and request
assembly via the shared `flow` path (`internal/cli/flow_exec.go`),
mutation routing via `writerFor`/`groupByWriter`, application and
commit-time verification via the accessor executor's write + read-back
(RDR 0004's model — read-back is the only commit check, no cross-writer
atomicity, and this verb inherits both statements verbatim). The only new
data flow is the *source* of the planned writes: `Model.Initial`
(`internal/table/model.go::Model.Initial`, today written by the loader
and read by lint only) instead of `--write`/`--clear` request flags.
That difference in plan source is why this is a new verb rather than a
`set-state` flag — see Load-Bearing Decisions (Naming).

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
that class rather than succeed vacuously.
```

**C2**

```normative
CARRIER. The `flow` group gains one verb, `init-state` — an override
extending 0005:C1's verb enumeration and 0005:D-naming's verb list by
exactly this spelling. It takes the shared selection flags
(`--flow`|`--model`), explicit `--artifact role=path` bindings, and no
write grammar: its planned writes are the model's `[initial]`
assignments. It MUST route every seed through the declared write
accessors with commit-time read-back, on the same
writer-routing/no-cross-writer-atomicity terms as `set-state`
(0005:C1); it MUST NOT write an artifact directly.

Seeding is ALL-OR-NOTHING over an EMPTY store, never a per-key merge:
init-state seeds if and only if the bound artifact carries NO key, and
then it seeds every `[initial]` key. A non-empty artifact — torn,
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
The payload MUST distinguish the seeded-all case from the no-op case,
and its scope is exactly the `[initial]` key set — the verb claims
nothing about owned keys `[initial]` does not assign.

The ENTIRE plan is validated before any write: every `[initial]` key
must route to exactly one declared writer and every needed artifact
role must be bound, or the verb refuses with ZERO writes committed —
per-key commit has no atomicity across writers, so plan-level
validation is where the all-or-nothing property lives. A read-back
mismatch is a distinct terminal refusal naming the key as
PRESENT-AND-UNVERIFIED; the store is then non-empty, so a re-run is a
no-op that does NOT repair it and MUST NOT be documented as its
recovery — recovery is an explicit `set-state` (or discarding the
artifact and re-running init). Against a `decision-table` model the
verb refuses before any accessor runs — a refusal (exit 2), not a
no-op success — with a dedicated code in the `flow-*` family (spelling
sharpened pre-lock); shared refusal classes (model selection, artifact
binding, writer routing, read-back) reuse the existing `flow-*` codes
unchanged, and the classes new to this verb get codes recorded in the
0005:C1 taxonomy extension this override carries.
```

Exact payload field names, the refusal-code spelling, and help text are
deliberately deferred to Resolve/Pre-Lock; C1/C2 fix the semantics and
the carrier, which is the fork this RDR exists to close.

#### Load-Bearing Decisions

- **Identity** — a seeded value's read-back equality is REQ-107's
  form: value-for-value over the keys init planned to seed, canonical
  wire form for set values — the identical equality `set-state`
  already verifies, not a new one (A2 verifies the rendering path).
- **Wire / byte format** — none new; artifacts keep
  `internal/cli/flowbind` shape, and the success payload rides the
  0005:C1 envelope. Field names deferred to Resolve/Pre-Lock, owned
  here.
- **Naming** — `init-state`, completing the `read-state`/`set-state`
  family. Rejected: `init` (reads as scaffolding a model file, not
  seeding state), `seed` (outside the `*-state` family the group
  established), and a `--from-initial` flag on `set-state` (a
  `set-state` payload echoes the request's planned writes; init's plan
  comes from the model — one verb with two write-plan sources muddies
  both grammars and REQ-3's flag list anyway).
- **Selection / predicate** — per store, not per key: seed (all
  `[initial]` keys) iff the bound artifact carries no key (the store's
  own presence answer, `internal/cli/flowbind/flowbind.go::store` —
  whose writer is any prior committed `set-state`/init write; an
  empty store means never-written or emptied by clears, which are
  indistinguishable at the content level, and both read as
  "initializable" — the cleared-vs-unseeded ambiguity is not removed by
  this predicate, it is confined to the single state where the store
  carries nothing at all, since one surviving key blocks seeding;
  C2 records that residual); class: refuse iff
  `internal/table/model.go::IsDecisionTable` — the shipped
  discriminator (its writer is the loader materializing `[model]
  class`, `0010:C1`), never a re-derivation from `len(owned)`.

#### Round-Trip / Inverse Invariants

- `read-state ∘ init-state = [initial] on the empty-store class`:
  after a seeding init over role R, `flow read-state` over R reports
  every `[initial]` key with its declared value — value-for-value in
  canonical form, the REQ-107 equality, not merely exit 0.
- `init-state ∘ init-state = init-state` (idempotence on any store):
  the second run writes nothing, succeeds, and the artifact bytes are
  unchanged.
- Equivalence with the manual path: `init-state` into a fresh
  artifact and the `set-state --write` transcription of the same
  `[initial]` assignments into a second fresh artifact produce
  byte-identical artifacts, for every value kind `[initial]` admits
  (the A2 spike's table; catches both the hard read-back divergence
  and the silent canonical-form divergence).
- Cleared-key preservation, on a store that retains a key:
  `init-state ∘ (set-state --clear k)` on a seeded artifact carrying at
  least one key besides k writes nothing, and `read-state` still
  reports k absent — REQ-107's reader guarantee stays untouched, and
  the writer-side resurrect hazard is closed for this class by the
  empty-store predicate; re-establishment is explicit `set-state`. The
  inverse holds at the boundary and is asserted as such: clearing the
  last key empties the store, so the next `init-state` reseeds every
  `[initial]` key (C2's disclosed residual).

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
| Class discrimination | `internal/table/model.go::IsDecisionTable` | none | Reuse | C2's class refusal keys on it |
| Writer routing | `internal/cli/flow_state.go::writerFor`, `::groupByWriter` | requires exactly one writer per key (A3) | Reuse | init refuses writerless `[initial]` keys unless A3 forces a lint arm |
| Write + read-back | accessor executor (`accessor.NewExecutor`, `Write`) | no cross-writer atomicity (stated, inherited) | Reuse | C2 inherits `set-state`'s terms verbatim |
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
predicate, and only within the boundary C2 discloses:

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
residual C2 discloses.

Premortem: switched (hardened) — the init-verb approach survived with
its seeding rule switched to the empty-store predicate; the remaining
findings are folded into C2, A2–A5, and Failure Modes.
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
  residual (C2). Rejected because it changes the artifact wire format
  every reader and RDR 0004's read-back comparison touch — a
  cross-verb blast radius out of proportion to a first-run verb. The
  residual is accepted and disclosed instead; revisit this if the wire
  format opens for another reason.
- **`--from-initial` flag on `set-state`**: one verb with two
  write-plan sources (request vs model) muddies the payload's
  writes-echo-the-request grammar; see Load-Bearing Decisions
  (Naming).
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
Low on the blocking axis: the shipped `models/rdr.toml` is a
decision table with zero owned tags, so no current first run is
blocked.

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
- **Assumed** — canonical-form fidelity of model-sourced values through
  the set-state write path (A2) and writer coverage of `[initial]` keys
  (A3) need Resolve verification before the reuse claim is load-bearing.

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
  initializable class and a later init reseeds it (C2). Operators who
  clear keys individually to reset must discard the artifact or expect
  the reseed; the boundary is contract text and a test, not a hidden
  edge.

### Risks and Mitigations

- **Risk**: a lint-certified model whose `[initial]` names a writerless
  (or multi-writer) key makes init refuse, re-creating the wall it
  exists to remove (A3).
  **Mitigation**: C2's plan-level validation guarantees the refusal
  lands with zero writes committed; Resolve verifies writer coverage,
  and if uncovered the clause gains a companion lint arm (writer
  existence-and-uniqueness for `[initial]` keys joins certification)
  rather than a direct-write bypass.
- **Risk**: operators read init as "reset to initial".
  **Mitigation**: the empty-store predicate makes init incapable of
  touching a store that carries any key; reset stays `set-state`'s job.
- **Risk**: an operator clears owned keys one at a time to "start
  over", empties the store, and a scheduled `init-state` reseeds —
  reading as the resurrect hazard the design rejects elsewhere.
  **Mitigation**: this is C2's disclosed residual, not a defect; the
  no-op payload and the docs must state the boundary in store terms
  ("while any key remains"), and MVV step 9 pins it so it cannot move
  silently. Closing it entirely would require artifact tombstones,
  rejected on blast radius (Briefly Rejected).

### Failure Modes

- **Visible**: init on a decision-table model refuses with the C2 class
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
- **Scope gap (disclosed)**: init's claim covers exactly the
  `[initial]` key set (C2); an owned key `[initial]` does not assign
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
8. `flow init-state` against a decision-table model refuses (exit 2)
   with the C2 class code; against the fixture with an unbound artifact
   role it refuses in the existing artifact-binding family with zero
   writes committed.
9. The boundary, asserted in the failing direction: clear the
   REMAINING key so the store carries none, then `flow init-state` —
   it seeds every `[initial]` key again. This is C2's disclosed
   residual and the test exists so a future change cannot silently
   move the boundary.

### Phase 1: Code Implementation

#### Step 1: Verb skeleton

`newFlowInitStateCmd` beside the four shipped verbs in `internal/cli`:
shared selection/tag/artifact registration, `ValidateMode`/respond
gateway, class refusal via `table.IsDecisionTable` before any accessor.

#### Step 2: Predicate and plan

Read the bound store; non-empty → the no-op arm (report absent
`[initial]` keys, write nothing). Empty → plan every `[initial]`
assignment and validate the WHOLE plan (`writerFor` per key, artifact
roles bound) before any accessor runs.

#### Step 3: Apply and report

Execute per writer with read-back (the `runFlowSetState` execution
shape); assemble the payload distinguishing the seeded arm from the
no-op arm.

### Phase 2: Operational Activation

#### Activation Step 1: Contract and reference docs

`docs/cli-output-contract.md` (verb I/O, refusal codes),
`docs/cli-reference.md`, `docs/model-authoring.md` (the start-state
sentence gains its runtime carrier), `--help-all` text.

## Validation

### Testing Strategy

The MVV sequence above is the acceptance spine and runs as an
integration test over a fixture model; "done" is every one of its
steps green plus the unit coverage below. Coverage goals are stated as
arms, not percentages — each normative arm of C2 owes at least one
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
   transcription, table-driven over every value kind `[initial]`
   admits.
   **Expected**: the two artifacts are byte-identical (RT3; this is
   A2's spike promoted to a standing test).
5. **Scenario**: `[initial]` names a key with zero declared writers,
   and separately one with more than one.
   **Expected**: refusal from the existing writer-routing family with
   ZERO writes committed — asserted on artifact bytes, not just exit
   code.
6. **Scenario**: A required artifact role is unbound.
   **Expected**: refusal in the existing artifact-binding family,
   zero writes committed.
7. **Scenario**: `decision-table` model.
   **Expected**: exit 2 with the C2 class code, raised before any
   accessor runs (asserted by a writer that fails if invoked).
8. **Scenario**: A writer whose read-back disagrees with the value
   written.
   **Expected**: terminal refusal naming the key present-and-unverified;
   a subsequent re-run is a no-op that does NOT repair it.

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
