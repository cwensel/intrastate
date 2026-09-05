# REQ List — RDR 0019 Owned-state initialization semantics

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0019-owned-state-initialization-semantics.md`. Quotes are verbatim —
copied from the projector (`rdr inspect --select <id>`) for fenced elements and
read from the record for testable prose outside the fences — never transcribed
by hand.

Element ids (`0019:C1`, `0019:MVV`, `0019:S10`) are carried wherever a REQ
derives from a labelled element, so a later stage can trace the REQ back to its
contract.

**Shape.** `counts.elements.C = 1`. The single fence `0019:C1` (25.4K, lines
607–1022) carries almost the whole normative surface, and it is unusually dense:
one fence, many dozens of independent obligations, organised internally by
run-in headings (SEMANTICS, CARRIER, CARRIER SCOPE, ORDERING, EMPTINESS). The
REQ set below decomposes that fence rather than treating it as one clause. A
material remainder is unfenced or structurally-labelled: the four
`Load-Bearing Decisions` (`0019:D-identity`, `D-wire-byte-format`, `D-naming`,
`D-selection-predicate`), the four `Round-Trip / Inverse Invariants`
(`0019:RT1`–`RT4`), the `Minimum Viable Validation` (`0019:MVV`, 10 steps), the
thirteen `Testing Strategy` scenarios (`0019:S1`–`S13`), the five `Failure
Modes` (`0019:F1`–`F5`), `Cross-Cutting Concerns` (`0019:G-cross-cutting`), the
`Technical Design` prose head, the `Implementation Plan` phases, and the
`Overrides` metadata field.

The eight `Critical Assumptions` (A1–A8) are all `Verified` and are read as
*evidence*, not as REQs — except where an assumption fixes an implementation
obligation the contract delegates to it (A6's surviving carrier, A7/A8's
type-switch spellings). Those are recorded and cited.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced,
  `0019:C1`)
- `TD` = Proposed Solution / Technical Design (unfenced prose head)
- `LBD` = Technical Design / Load-Bearing Decisions (`0019:D-*`)
- `RT` = Technical Design / Round-Trip / Inverse Invariants (`0019:RT1`–`RT4`)
- `AUDIT` = Proposed Solution / Existing Infrastructure Audit
- `META` = Metadata (Overrides field)
- `CONS` = Trade-offs / Consequences
- `FM` = Trade-offs / Failure Modes (`0019:F1`–`F5`)
- `IP` = Implementation Plan (prerequisites, Phase 1 steps, Phase 2)
- `MVV` = Implementation Plan / Minimum Viable Validation (`0019:MVV`)
- `TS` = Validation / Testing Strategy (header "done" clause)
- `SC-n` = Validation / Testing Strategy, numbered scenario *n* (`0019:Sn`)
- `XC` = Finalization Gate / Cross-Cutting Concerns (`0019:G-cross-cutting`)

---

## A. `[initial]` semantics and the no-synthesis prohibition (C1)

- [REQ-1] "`[initial]` is BOTH the lint-time reachability root (0006:C18
  unchanged — a state-machine model declaring none stays a blocking finding) AND
  the runtime bootstrap source for owned state." — (NC, `0019:C1`)

- [REQ-2] "Materializing it is an EXPLICIT, PERSISTING act: no read verb, no
  artifact load, and no accessor read path may synthesize, default, or fall back
  to `[initial]` values when a key is absent" — (NC, `0019:C1`) — the
  prohibition S13 asserts as an invariant, not a snapshot.

- [REQ-3] "a cleared key reads back absent for every reader (REQ-107
  unchanged), and owned state is assembled only from caller-bound artifacts
  (0004:C3 unchanged)." — (NC, `0019:C1`) — conformance, not restatement (XC
  "Read purity / never-fill").

## B. Class refusal (C1 SEMANTICS)

- [REQ-4] "A `decision-table` model has no `[initial]` (0010:C2) and therefore
  no bootstrap to materialize; initialization MUST refuse for that class rather
  than succeed vacuously." — (NC, `0019:C1`)

- [REQ-5] "The refusal keys on `internal/table/model.go::IsDecisionTable`
  ALONE, never on \"declares `[initial]` but is decision-table\": that state is
  UNCONSTRUCTIBLE" — (NC, `0019:C1`); LBD restates the ban on the alternative:
  "never a re-derivation from `len(owned)`" — (LBD, `0019:D-selection-predicate`).

- [REQ-6] "the verb's refusal is a CLASS refusal, not an emptiness one — an
  ordinary state-machine model whose `[initial]` is empty is a different case,
  already blocked at lint by 0006:C18." — (NC, `0019:C1`)

- [REQ-7] "the verb refuses before any accessor runs — a refusal (exit 2), not
  a no-op success." — (NC, `0019:C1`) — asserted by the absence of any accessor
  invocation at all, since a decision-table fixture can declare no owned-tag
  writer (SC-7, `0019:S7`).

## C. Verb carrier, grammar, and the 0005 override (C1 CARRIER; META)

- [REQ-8] "The `flow` group gains one verb, `init-state` — an override
  extending 0005:C1's verb enumeration and 0005:D-naming's verb list by exactly
  this spelling." — (NC, `0019:C1`)

- [REQ-9] "**appended, not conditioned**, and additive on its owner's grammar:
  no existing verb changes and nothing is withdrawn. The new verb inherits C1's
  per-verb MUSTs (`respond.ValidateMode` first, `respond.OK`/`Fail`,
  `SilenceErrors`/`SilenceUsage`) unchanged" — (META, Overrides)

- [REQ-10] "It takes the shared selection flags (`--flow`|`--model`), explicit
  `--artifact role=path` bindings, and no write grammar: its planned writes are
  the model's `[initial]` assignments." — (NC, `0019:C1`)

- [REQ-11] "It also carries `--allow-commands`, not by declaring it but by
  INHERITING it: the flag is registered once on the `flow` group's persistent
  set (`internal/cli/flow.go::newFlowCmd`), so every child verb takes it
  structurally and this one cannot decline it." — (NC, `0019:C1`) — and "The
  enumeration above is the verb's OWN grammar, not the closed set of flags it
  accepts".

- [REQ-12] "It MUST route every seed through the declared write accessors with
  commit-time read-back, on the same writer-routing/no-cross-writer-atomicity
  terms as `set-state` (0005:C1); it MUST NOT write an artifact directly." —
  (NC, `0019:C1`)

- [REQ-13] "`init-state`, completing the `read-state`/`set-state` family." —
  (LBD, `0019:D-naming`) — the verb spelling is normative; `init`, `seed`, and a
  `--from-initial` flag on `set-state` are rejected.

- [REQ-14] "`internal/cli/flow_surface_0005_test.go` fixes `flowVerbs = {next,
  resolve, read-state, set-state}` and asserts SET EQUALITY against the
  registered group, so it must gain `init-state` here" — (IP, Phase 1 Step 1) —
  registration lands A4's override in the SAME step.

- [REQ-15] "Six code/doc sites hard-code the cardinal and change with it:
  `internal/cli/flow.go`'s package comment, `::newFlowCmd`'s `Long` body,
  `::flowExtendedDesc`, `internal/cli/flow_exec.go`'s header comment,
  `docs/cli-reference.md` (two strings), and `docs/model-authoring.md`." —
  (IP, Phase 1 Step 1)

- [REQ-16] "Regenerate and commit `llms.txt` (`internal/cli/docs.go`); it
  enumerates the verbs as bullets and hard-codes no cardinal, as does
  `docs/cli-output-contract.md`, so neither needs a cardinal edit" — (IP,
  Phase 1 Step 1)

## D. Seed encoding (C1)

- [REQ-17] "Seed values are taken from the LOADER-NORMALIZED `Model.Initial`
  assignments, NOT by transcribing them to their `--write` argv spelling" —
  (NC, `0019:C1`)

- [REQ-18] "each value is rendered in JDR 0001 §D13's canonical form,
  DISPATCHED ON THE DECLARED KIND, which is what makes read-back equality byte
  equality and leaves no third encoding." — (NC, `0019:C1`)

- [REQ-19] "the seed encoder is exactly the SET arm of `set-state`'s: for
  `decl.Kind == \"set\"` it is `internal/cli/flow_input.go::canonicalSet` (sort,
  compact, JSON array); for every scalar kind it is the single member verbatim,
  `members[0]`." — (NC, `0019:C1`)

- [REQ-20] "`::canonicalSet` alone is NOT the whole encoder — it takes
  `[]string` and renders a JSON array, so applying it to a scalar seed would
  persist `[\"draft\"]` where `set-state --write status=draft` persists `draft`"
  — (NC, `0019:C1`)

- [REQ-21] "The argv encoder `::canonicalValue` is the correct *reference* for
  both arms' output but is NOT the call: it re-runs the argv admission checks
  this clause rules out below." — (NC, `0019:C1`)

- [REQ-22] "Conformance to the declaration is ALREADY HELD by the loader and is
  not re-established here … A re-conform pass on a loader-normalized value
  cannot fail and buys nothing; the byte-equality property comes from the
  encoder alone." — (NC, `0019:C1`) — no re-conform pass at seed time.

- [REQ-23] "the loader's `[initial]` admission set is a proper SUPERSET of the
  argv route's: a bare scalar for a set-valued tag and a SINGLE-MEMBER array
  literal for a scalar tag both load and normalize … and are refused only on the
  argv surface … Transcribing through argv would therefore refuse at seed time
  models that load clean, re-creating the first-run wall this verb exists to
  remove (A2)." — (NC, `0019:C1`)

## E. The empty-store predicate (C1; LBD selection/predicate)

- [REQ-24] "Seeding is ALL-OR-NOTHING over an EMPTY store, never a per-key
  merge: init-state seeds if and only if EVERY bound artifact carries NO key,
  and then it seeds every `[initial]` key." — (NC, `0019:C1`)

- [REQ-25] "The quantifier is ALL, not ANY and not per-artifact, and it is
  load-bearing rather than stylistic … a per-artifact or ANY reading would seed
  into a store that still carries a key whenever a SIBLING artifact happened to
  be empty" — (NC, `0019:C1`) — SC-11 (`0019:S11`) is its only test.

- [REQ-26] "Under ALL, one surviving key anywhere in the bound set blocks the
  whole seed, so the cleared-key guarantee below holds per model rather than
  merely per artifact." — (NC, `0019:C1`)

- [REQ-27] "The count is over STORE keys, not owned keys, and the difference is
  reachable: a read-back-SEALED artifact carries
  `internal/cli/flowbind/flowbind.go::sealedKey` … so an artifact whose owned
  keys have all been cleared but whose last write declared an unreachable
  read-back locator is a ONE-key, NON-EMPTY store that init-state declines to
  seed." — (NC, `0019:C1`)

- [REQ-28] "The sealed store is nonetheless a NO-OP SUCCESS at exit 0 here, not
  an exit-3 refusal … What it needs is a key COUNT, which the predicate obtains
  without reading key values" — (NC, `0019:C1`) — SC-9 (`0019:S9`) asserts it on
  artifact bytes.

- [REQ-29] "whichever carrier lands MUST answer cardinality on a sealed store
  rather than inheriting `::Reader.Read`'s unreadable short-circuit, or this arm
  degrades from no-op success to an exit-3 refusal and S9 fails." — (NC,
  `0019:C1`; delegated to A6, whose VERIFIED disposition names candidate (b), an
  exported `flowbind` cardinality probe, as "the surviving carrier", reading
  `len(load(path))` so the seal counts as 1 without entering the seal branch).

- [REQ-30] "no implementation may substitute \"every `[initial]` key reads
  absent\" for it, since that is the per-key variant this contract rejects (it
  cannot see a non-`[initial]` key, including `::sealedKey`, and so would seed
  into a non-empty store)." — (NC, `0019:C1`)

- [REQ-31] "\"Bound\" is not a filter that can shrink the quantifier's domain …
  An unbound needed role is therefore a REFUSAL, never an artifact treated as
  empty" — (NC, `0019:C1`)

## F. Ordering of the gates (C1)

- [REQ-32] "the ROLE-BINDING half of plan validation runs BEFORE the emptiness
  read, so an invocation that leaves a needed role unbound refuses at exit 2
  REGARDLESS of what the bound stores contain — it does not reach the predicate
  and cannot take the no-op arm." — (NC, `0019:C1`)

- [REQ-33] "It runs AFTER the carrier gate, which is likewise observable: an
  invocation against a non-file-backed model that ALSO omits a needed artifact
  flag refuses under the carrier code, not the artifact-binding family." — (NC,
  `0019:C1`)

- [REQ-34] "The WRITER-ARITY half is not ordered here at all, having already
  been settled at load … Only encoding and the write follow the predicate." —
  (NC, `0019:C1`)

- [REQ-35] "ORDERING: the CLASS refusal precedes the CARRIER refusal." — (NC,
  `0019:C1`) — "the class arm keys on `::IsDecisionTable`, which is answerable
  from the loaded model alone, while the carrier arm needs the constructed
  registry."

- [REQ-36] "The carrier refusal is decided AFTER registry construction (it
  reads the constructed binding) but BEFORE any accessor is invoked, and it
  PREEMPTS the `--allow-commands` refusal
  (`internal/accessor/model.go::Registry.AllowCommands`, also exit 2). A
  command-backed accessor therefore yields the carrier code whether or not the
  opt-in was passed" — (NC, `0019:C1`) — SC-10 runs the command-backed rows with
  `--allow-commands` UNSET and asserts the carrier code.

- [REQ-37] "Gate the carrier on both capabilities (Step 1's class refusal has
  already fired); non-file-backed → refuse (exit 2). Then validate that every
  role a needed writer names is bound; unbound → refuse (exit 2) without reading
  any store. Then read the bound stores; non-empty → the no-op arm (report
  absent `[initial]` keys, write nothing). Empty → plan and encode every
  `[initial]` assignment (`writerFor` per key) before any accessor runs." —
  (IP, Phase 1 Step 2)

## G. Carrier scope and detection (C1 CARRIER SCOPE)

- [REQ-38] "Against a model any of whose bound write accessors is
  non-file-backed, init-state MUST refuse — a distinct terminal refusal in the
  `flow-*` family naming the accessor and its carrier — rather than seed on an
  emptiness answer it cannot compute." — (NC, `0019:C1`)

- [REQ-39] "The carrier gate is over BOTH capabilities, not the write side
  alone … The refusal therefore fires when EITHER the bound write binding or the
  bound read binding for a needed role is not the file-backed type, and it names
  which capability and which accessor failed." — (NC, `0019:C1`)

- [REQ-40] "The verb detects the carrier by the CONSTRUCTED BINDING'S TYPE, not
  by re-reading the accessor declaration … The verb type-switches on
  `internal/accessor/model.go::Definition.Binding` … on the write and read
  bindings alike: it admits ONLY `*flowbind.Writer` on the write side, refusing
  `*flowbind.EditWriter` and `*cmdbind.Writer`, and ONLY `flowbind.Reader` on
  the read side, refusing `cmdbind.Reader`." — (NC, `0019:C1`; A7, A8)

- [REQ-41] "The POINTER/VALUE spelling differs between the two capabilities and
  is normative here, because a type-switch case on the wrong one matches
  nothing … A read-side case written `*flowbind.Reader` would therefore refuse
  EVERY model, file-backed ones included" — (NC, `0019:C1`) — S1 is the second
  control (SC-10).

- [REQ-42] "Re-deriving the carrier in `internal/cli` from `table.Accessor`'s
  `Edit`/`Command` fields is BANNED … Exporting `commandBacked` is not required
  and is not authorized here." — (NC, `0019:C1`)

- [REQ-43] "The verb selects over the exported `::Registry.Definitions` slice,
  admitting the FIRST definition whose `Identity.Capability` is
  `accessor.CapRead` and whose `Accessor.Role` is the needed role. FIRST MATCH
  IN REGISTRY ORDER is normative, not incidental" — (NC, `0019:C1`) — a gate
  that refused on ambiguity would be WIDER than `::readerFor`; selecting the
  last would disagree.

- [REQ-44] "A no-match is the flat `false` `::readerFor` returns, which reaches
  the unbound-needed-role arm above; it MUST NOT be split into a separate
  missing-reader code." — (NC, `0019:C1`)

- [REQ-45] "The gate reads `Identity`, `Accessor`, and `Binding` as FIELDS and
  invokes no `Binding` method … which is what makes S10's no-accessor-invoked
  assertion true by construction rather than by timing." — (NC, `0019:C1`)

- [REQ-46] "the gate must not sort, filter, or re-order `Definitions` before
  scanning, since the equivalence rides on sharing the registry's own slice
  order." — (NC, `0019:C1`) — "Exporting `readerFor` would serve equally but is
  an `internal/accessor` API change this RDR does not authorize".

## H. The no-op arm and the cleared-key guarantee (C1; RT4; F2–F4)

- [REQ-47] "A non-empty artifact — torn, partially seeded, post-clear, or fully
  seeded alike — is a NO-OP SUCCESS: zero writes, and the payload reports which
  `[initial]` keys the store does not carry (informational, so a torn or
  post-clear state is visible without being repaired, resurrected, or failed
  on)." — (NC, `0019:C1`)

- [REQ-48] "while the store retains AT LEAST ONE key, a cleared key is never
  re-established by init-state under any invocation pattern, including
  unconditional automated re-runs (the artifact is non-empty, so nothing
  writes)." — (NC, `0019:C1`)

- [REQ-49] "clearing the LAST remaining key empties the store, and an emptied
  store is indistinguishable at the content level from a never-written one …
  A subsequent init-state therefore RESEEDS such a store. This is the accepted
  residual of rejecting tombstones … no payload or doc may describe init-state
  as unable to re-establish a cleared key without this qualification." — (NC,
  `0019:C1`) — MVV step 9 asserts the boundary in the failing direction.

- [REQ-50] "A key added to `[initial]` after seeding is re-established by
  explicit `set-state`, not by init — the payload's absent-key report names it."
  — (NC, `0019:C1`; FM, `0019:F4`)

- [REQ-51] "the no-op re-run reports the absent `[initial]` keys and never
  repairs them (F2) — visibility over silent merge is the posture, and the
  repair route is explicit `set-state` of the listed keys." — (XC, "Concurrency
  model"; FM, `0019:F2`) — SC-12 (`0019:S12`) is the only test of F2.

- [REQ-52] "init's claim covers exactly the `[initial]` key set (C1); an owned
  key `[initial]` does not assign can still surface `unknown[].reason: absent`
  in `flow next` after a successful init." — (FM, `0019:F3`)

## I. The empty scalar (C1; F5)

- [REQ-53] "An `[initial]` value that is an EMPTY SCALAR (`note = \"\"`) seeds
  like any other: the scalar arm renders `members[0]`, the artifact carries
  `{\"note\":\"\"}`, and the key reads back PRESENT — distinct from a cleared
  key, which leaves the object entirely." — (NC, `0019:C1`)

- [REQ-54] "It creates no third encoding and is not refused at seed time,
  because refusing it would require re-conforming a loader-normalized value,
  which this clause rules out above." — (NC, `0019:C1`)

- [REQ-55] "The ARGV route's contrary refusal (`::canonicalValue` rejects
  `--write note=`) is a loader/write-surface asymmetry predating this RDR and is
  NOT resolved here in either direction (F5, routed to RDR 0002)" — (NC,
  `0019:C1`; FM, `0019:F5`) — out of scope; no change to `::canonicalValue`.

## J. Payload content (C1; LBD wire/byte format)

- [REQ-56] "The payload MUST distinguish the seeded-all case from the no-op
  case, and its scope is exactly the `[initial]` key set — the verb claims
  nothing about owned keys `[initial]` does not assign." — (NC, `0019:C1`)

- [REQ-57] "the payload reports seeded keys by NAME and does not echo their
  values — `read-state` is the surface that reports values (RT1), and echoing
  them here would create a second place a seeded value can be read from" — (NC,
  `0019:C1`)

- [REQ-58] "The absent-key report belongs to the NO-OP arm, which is the only
  arm where the set can be non-empty; on the seeded arm every `[initial]` key
  was just written, so an absent list there is necessarily empty and MUST NOT be
  reported as if it carried information." — (NC, `0019:C1`)

- [REQ-59] "none new; artifacts keep `internal/cli/flowbind` shape, and the
  success payload rides the 0005:C1 envelope. Field names deferred to
  implementation, owned here; the payload's CONTENT is not deferred" — (LBD,
  `0019:D-wire-byte-format`)

## K. Plan validation and read-back refusals (C1; F1)

- [REQ-60] "The ENTIRE plan is validated before any write: every `[initial]`
  key must route to exactly one declared writer and every needed artifact role
  must be bound, or the verb refuses with ZERO writes committed — per-key commit
  has no atomicity across writers, so plan-level validation is where the
  all-or-nothing property lives." — (NC, `0019:C1`)

- [REQ-61] "the writer-arity half is already enforced at LOAD …  Such a model
  refuses as `flow-model-invalid` at load, never in the writer-routing family,
  and S5 asserts that ordering. The verb states the requirement (it MUST NOT
  route a key to two writers) but owes no runtime check of its own for it." —
  (NC, `0019:C1`)

- [REQ-62] "The role-binding half IS reachable, since roles are bound per
  invocation (`::runFlowSetState` resolves `req.artifacts[def.Accessor.Role]`),
  and S6 is its test." — (NC, `0019:C1`)

- [REQ-63] "A read-back mismatch is a distinct terminal refusal naming the key
  as PRESENT-AND-UNVERIFIED; the store is then non-empty, so a re-run is a no-op
  that does NOT repair it and MUST NOT be documented as its recovery — recovery
  is an explicit `set-state` (or discarding the artifact and re-running init)."
  — (NC, `0019:C1`)

## L. Refusal codes and exit groups (C1)

- [REQ-64] "It carries a dedicated code in the `flow-*` family, sharpened here
  as `flow-init-class-unsupported`; the carrier refusal above is
  `flow-init-carrier-unsupported`. Both follow the shipped
  `flow-<subject>-<condition>` convention" — (NC, `0019:C1`)

- [REQ-65] "What this clause fixes NORMATIVELY is the exit GROUP (2 on both)
  and that the two are DISTINCT codes and distinct from each other and from the
  shared classes — the literal spellings are non-normative on the precedent of
  `0028:C1.3`'s `codeWriteEditRefused`, so no test may pin the string." — (NC,
  `0019:C1`) — a deliberate CARVE-OUT from the code table's general rule, which
  says the opposite.

- [REQ-66] "Shared refusal classes reuse the existing `flow-*` codes unchanged,
  at the groups those codes already carry — this verb introduces no new group
  for them and may not re-map one. Model selection
  (`::selectModel`/`::selectModelPath`), artifact binding and writer routing
  (`internal/cli/flow_state.go::writerFor`) are all `clierr.GroupUserEnv`,
  exit 2." — (NC, `0019:C1`)

- [REQ-67] "a read-back that COMPLETED and disagreed is exit 2
  (`GroupUserEnv`), while a read-back that could not complete — the unreachable-
  locator and timeout arms — is exit 3 (`GroupEnvUnavailable`) … So the
  PRESENT-AND-UNVERIFIED refusal above is exit 2, and a seed whose read-back is
  incomplete exits 3 with the same partial-write disposition; both leave a
  non-empty store, so both make a re-run a no-op that does not repair." — (NC,
  `0019:C1`; FM, `0019:F1`)

- [REQ-68] "The classes new to this verb get codes recorded in the 0005:C1
  taxonomy extension this override carries." — (NC, `0019:C1`)

## M. Round-trip / inverse invariants (RT1–RT4)

- [REQ-69] "after a seeding init over role R, `flow read-state` over R reports
  every `[initial]` key with its declared value — value-for-value in canonical
  form, the REQ-107 equality, not merely exit 0." — (RT, `0019:RT1`)

- [REQ-70] "the second run writes nothing, succeeds, and the bytes of EVERY
  bound artifact are unchanged (the byte assertion is scoped to the file-backed
  carrier, the only one this RDR admits — C1)." — (RT, `0019:RT2`)

- [REQ-71] "`init-state` into a fresh artifact and the `set-state --write`
  transcription of the same `[initial]` assignments into a second fresh artifact
  produce byte-identical artifacts, for every value kind BOTH routes admit AND
  every VALUE whose two routes carry the same spelling" — (RT, `0019:RT3`)

- [REQ-72] "the three kinds only the loader admits (bare scalar for a set tag,
  array for a scalar tag, empty scalar) have no `--write` transcription to
  compare against and are outside this invariant" — (RT, `0019:RT3`)

- [REQ-73] "on the numeric kinds the two routes agree on the value and may
  DISAGREE on its spelling, so the invariant is scoped to the spelling the
  loader normalizes to … `[initial] threshold = 1.0` seeds `\"1\"` where
  `--write threshold=1.0` writes `\"1.0\"`, and `--write n=+5` writes `\"+5\"`
  where the loader seeds `\"5\"`." — (RT, `0019:RT3`)

- [REQ-74] "`init-state ∘ (set-state --clear k)` on a seeded artifact carrying
  at least one key besides k writes nothing, and `read-state` still reports k
  absent — REQ-107's reader guarantee stays untouched" — (RT, `0019:RT4`)

- [REQ-75] "a seeded value's read-back equality is REQ-107's form:
  value-for-value over the keys init planned to seed, canonical wire form for
  set values — the identical equality `set-state` already verifies, not a new
  one" — (LBD, `0019:D-identity`)

## N. Minimum Viable Validation (MVV)

- [REQ-MVV] "Over a fixture state-machine model declaring `[initial]` with one
  always-present owned key and one plain owned key, both writer-served:" —
  (MVV, `0019:MVV`) — the acceptance spine, ten steps, run as an integration
  test; "\"done\" is every one of its steps green plus the unit coverage below"
  (TS).

  Steps, each individually gating:

  - [REQ-MVV.1] "`flow lint` certifies the model (0006 arms all green)."
  - [REQ-MVV.2] "`flow init-state --artifact state=<fresh path>` exits 0;
    payload reports the empty-store seeding arm with both keys seeded."
  - [REQ-MVV.3] "`flow read-state` reports both keys at their `[initial]`
    values (canonical form) — round-trip invariant 1."
  - [REQ-MVV.4] "`flow next` over the same artifact reports candidates with no
    `unknown[].reason: absent` entry for either key — the first-run wall is
    gone."
  - [REQ-MVV.5] "Re-run `flow init-state`: exit 0, zero writes, artifact bytes
    unchanged — invariant 2."
  - [REQ-MVV.6] "`flow set-state --clear <plain key>`; `flow read-state`
    reports it absent — REQ-107 held."
  - [REQ-MVV.7] "`flow init-state` again (the unconditional-automation
    pattern): exit 0, ZERO writes, payload reports the cleared key as an absent
    `[initial]` key, `flow read-state` still reports it absent — invariant 4 /
    A5."
  - [REQ-MVV.8] "`flow init-state` against a decision-table model that LOADS
    (no `[initial]`, no owned tags — see S7) refuses (exit 2) with the C1 class
    code, with no accessor run and no artifact touched; against the fixture with
    an unbound artifact role it refuses in the existing artifact-binding family
    with zero writes committed; against a model whose write accessor is
    edit-carried or command-backed it refuses with the carrier code and zero
    writes (S10)."
  - [REQ-MVV.9] "The boundary, asserted in the failing direction: clear the
    REMAINING key so the store carries none, then `flow init-state` — it seeds
    every `[initial]` key again."
  - [REQ-MVV.10] "`flow init-state` against `models/rdr.toml` … bound to a
    fresh artifact seeds exactly `stage`, `status`, `gate_passed` at their
    declared values, and `flow next` then reports no `unknown[].reason: absent`
    for them." — the user outcome on a SHIPPED model, not a fixture. Its
    `[tags.gate_passed]` bool-vs-string spelling "must not be read as" covering
    the encoder arms — "S4's kind table is what covers the arms".

## O. Testing Strategy scenarios (S1–S13)

- [REQ-76] "each normative arm of C1 owes at least one test that fails if the
  arm is removed." — (TS) — coverage goals are stated as arms, not percentages.

- [REQ-77] SC-1 (`0019:S1`): "Empty store, every `[initial]` key
  writer-served. **Expected**: exit 0; every key seeded; `read-state` returns
  each at its declared value in canonical form (RT1)."

- [REQ-78] SC-2 (`0019:S2`): "Re-run against the store just seeded.
  **Expected**: exit 0; zero writes; artifact bytes unchanged (RT2); payload
  reports the no-op arm."

- [REQ-79] SC-3 (`0019:S3`): "Seeded store, one plain key cleared, then
  re-run. **Expected**: exit 0; zero writes; the cleared key still reads absent
  and is listed as an absent `[initial]` key in the payload (RT4 / A5)."

- [REQ-80] SC-4 (`0019:S4`): "table-driven over every value kind BOTH routes
  admit (9 kinds: enum, scalar string, bool, int, float, set array, set with
  duplicates and HTML characters, empty set, single-member set) … **Expected**:
  the two artifacts are byte-identical (RT3; this is A2's spike promoted to a
  standing test)."

- [REQ-81] SC-4 argv spelling: "Each row's `--write` argv MUST use the spelling
  the loader normalizes to (`1`, not `1.0`; `5`, not `+5`) — the numeric rows
  otherwise diff on a transcription difference RT3 does not claim, and a test
  written the other way fails on day one for the wrong reason (A2)."

- [REQ-82] SC-4 boundary row: "One ADDITIONAL row asserts the boundary rather
  than the invariant: `[initial] threshold = 1.0` seeded, against `--write
  threshold=1.0`, produces artifacts that DIFFER (`\"1\"` vs `\"1.0\"`) and both
  read back their own value — the divergence is pinned as decided behavior".

- [REQ-83] SC-4 loader-only kinds: "the 2 kinds only the loader admits (bare
  scalar for a set tag, array for a scalar tag) seed successfully via the
  normalized path and read back value-for-value … The empty scalar is the
  loader's THIRD such divergence and IS a row here on the same terms: `[initial]
  note = \"\"` seeds, the artifact carries `{\"note\":\"\"}`, and the key reads
  back PRESENT with an empty value — asserted distinct from a cleared key".

- [REQ-84] SC-5 (`0019:S5`): "**Expected**: refusal at LOAD
  (`flow-model-invalid`), BEFORE the verb's own plan validation runs — NOT the
  writer-routing family … and asserts ZERO writes committed: artifact
  bytes/mtime unchanged, or the artifact path still absent."

- [REQ-85] SC-6 (`0019:S6`): "A required artifact role is unbound.
  **Expected**: refusal in the existing artifact-binding family … zero writes
  committed, asserted on artifact bytes/mtime unchanged, or the artifact path
  still absent."

- [REQ-86] SC-7 (`0019:S7`): "**Expected**: exit 2 with the C1 class code,
  raised before any accessor runs … Assert instead that no accessor process ran
  and no artifact was created or modified (bytes/mtime unchanged, or the
  artifact path still absent), which is the discriminating control".

- [REQ-87] SC-8 (`0019:S8`): "the fixture binds the READ accessor to a
  different artifact path than its writer, pre-seeded with a conflicting value
  for the same key … **Expected**: terminal refusal naming the key
  present-and-unverified, distinct from the read-back-INCOMPLETE class an
  unreachable locator raises; a subsequent re-run is a no-op that does NOT
  repair it — asserted, since C1 forbids documenting the re-run as recovery."

- [REQ-88] SC-9 (`0019:S9`): "**Expected**: NO-OP SUCCESS with zero writes —
  the store carries the seal and is therefore non-empty, even though it holds no
  owned key. Asserted on artifact bytes, not just exit code. This pins the
  predicate's count to STORE keys".

- [REQ-89] SC-9 setup obligation: "The test MUST fail loudly rather than skip
  if its setup does not produce exactly `{sealedKey}` — assert the post-setup
  bytes before the scenario runs, so a representation change surfaces as a
  failure naming this coupling. And when A6's carrier lands, if it can construct
  the sealed state directly, this fixture MUST be rewritten onto it".

- [REQ-90] SC-10 (`0019:S10`): "run three times over an otherwise S1-shaped
  fixture. Twice on the WRITE side: once with an edit-carried accessor … and
  once with a command-backed one … Once on the READ side: a file-backed
  `[write.x]` paired with a command-backed `[read.x]` on the same role, which
  the write-side type-switch alone would ADMIT".

- [REQ-91] SC-10 expectation: "**Expected**: exit 2 with the carrier refusal
  code, naming the accessor and its carrier, with ZERO writes and no accessor
  invocation … the discriminating control is that removing the carrier check
  makes the edit case attempt a write. For the edit carrier the refusal is
  over-determined … and the test must still see the carrier code, not a
  missing-reader error. The command-backed rows are run with `--allow-commands`
  UNSET and assert the carrier code rather than the allow-commands refusal … The
  read-side row asserts the refusal names the READ capability and its accessor".

- [REQ-92] SC-11 (`0019:S11`): "Role A's artifact carries a key; role B's is
  empty. Run `init-state`. **Expected**: NO-OP SUCCESS with ZERO writes to BOTH
  artifacts — asserted on both artifacts' bytes, not exit code alone. This is
  the only test of C1's ALL quantifier".

- [REQ-93] SC-12 (`0019:S12`): "The torn state is built by writing role A's
  artifact directly during SETUP, never through the verb … **Expected**: NO-OP
  SUCCESS whose payload names exactly B's un-seeded `[initial]` keys, with A's
  committed values NOT rewritten (asserted on A's artifact bytes/mtime). This is
  the only test of F2".

- [REQ-94] SC-13 (`0019:S13`): "bind an artifact in which that key is ABSENT —
  never seeded, as distinct from cleared — and read it back on every path that
  loads state: `read-state`, and `next`'s candidate computation. **Expected**:
  the key reports ABSENT on every path, and `next` reports it under
  `unknown[].reason: absent`. No read path returns the `[initial]` value." —
  the only scenario pinning REQ-2. "The discriminating control is that the same
  fixture with the key PRESENT reads back its stored value".

## P. Phase 2 — operational activation (IP)

- [REQ-95] "`docs/cli-output-contract.md` (verb I/O, refusal codes),
  `docs/cli-reference.md`, `docs/model-authoring.md` (the start-state sentence
  gains its runtime carrier), `--help-all` text." — (IP, Phase 2 Activation
  Step 1)

- [REQ-96] "it states the FILE-BACKED scope (C1 carrier scope), not an
  unqualified \"every state-machine flow\"; and it names the verb from the
  surface where an operator meets the wall — `flow next`'s
  `unknown[].reason: absent` report" — (IP, Phase 2)

- [REQ-97] "Pointing an operator from `flow next` to `init-state` in-band would
  mean adding a field or a reason-string to another verb's payload, which is an
  override of `0005:C1`'s I/O for `next` that this record does not carry — its
  override is the verb enumeration, nothing more." — (IP, Phase 2) — no
  `next`-side payload change.

## Q. Prerequisites (IP)

- [REQ-98] "All Critical Assumptions verified (A2's canonical-form spike and
  A3's writer-coverage audit gate the design's reuse claim)" — (IP,
  Prerequisites) — satisfied: A1–A8 all carry `Status: Verified`.

## R. Non-goals / explicit non-obligations

Recorded so a later phase does not manufacture work the record forbids.

- [REQ-99] "no existing verb, flag, artifact shape or model schema changes." —
  (XC, "Incremental adoption")

- [REQ-100] "no lint arm owed; a model that loads already satisfies it" —
  (AUDIT, writer-routing row) — the verb owes no writer-arity lint arm.

- [REQ-101] "Extending the predicate to those carriers is deliberately out of
  scope here (see A6, and the successor noted in Consequences)" — (NC,
  `0019:C1`) — edit-carried (0028) and command-backed (0025) write accessors,
  and the command-backed reader, stay refused.

- [REQ-102] "A repair path that preserved the cleared-key guarantee (a seed
  scoped to the failed writer's role, gated on that role's artifact being empty)
  is a successor's, not this record's." — (FM, `0019:F2`)

- [REQ-103] "per-finding code identity against it is cli/0017's question, not
  this record's." — (XC, "Versioning")

---

## ASSUMPTIONS

Implicit choices made where the record's wording was imprecise but a single
reading is defensible. Each is recorded so a later phase can challenge it
against the record rather than against this audit.

- **ASSUMPTION-1 (A6's carrier is candidate (b)).** C1 leaves the
  emptiness-read carrier open ("Which carrier serves it … is UNDECIDED here and
  is booked as A6"), but A6's `Status: Verified` evidence closes it: "VERIFIED
  at reconcile: candidate (b), an exported `flowbind` cardinality probe, is the
  surviving carrier". Phase 1 therefore implements an exported cardinality probe
  on `internal/cli/flowbind`, reading `len(load(path))` so a read-back-sealed
  store counts 1 (REQ-29). C1's "UNDECIDED" is read as the state at drafting,
  superseded by A6's verification — not as a licence to pick freely.

- **ASSUMPTION-2 (payload field names).** `0019:D-wire-byte-format` defers field
  NAMES to implementation while fixing CONTENT. The implementation chooses names
  consistent with the shipped `flow` payloads; no test may pin a name that C1
  does not fix, and REQ-56/57/58 are asserted on content (seeded key names on
  the seed arm; the absent list on the no-op arm only).

- **ASSUMPTION-3 (refusal-code spellings).** REQ-65 makes the two new code
  spellings non-normative. The implementation uses C1's own sharpened spellings
  `flow-init-class-unsupported` and `flow-init-carrier-unsupported` as the
  concrete constants, since they are what the record names, while tests assert
  only the exit GROUP and code DISTINCTNESS — never the literal string.

- **ASSUMPTION-4 ("needed" role).** C1 speaks of "every role a needed writer
  names" and "the bound read binding for a needed role" without a standalone
  definition. Read as: a role named by a write accessor that serves at least one
  `[initial]` key (the keys the verb plans to seed), plus that role's read
  binding for the emptiness read and the commit-time read-back. A write accessor
  serving no `[initial]` key is not "needed" and neither its binding nor its
  role participates in the carrier gate, the role-binding check, or the ALL
  quantifier.

- **ASSUMPTION-5 (`--allow-commands` remains registered).** REQ-11 fixes that
  the verb inherits the flag structurally. Since the carrier gate refuses every
  command-backed binding before invocation (REQ-36), the flag is inert on this
  verb; the implementation neither removes it from the group's persistent set
  nor special-cases it on this child.

- **ASSUMPTION-6 (verb registration site).** IP Phase 1 Step 1 names
  `newFlowInitStateCmd` "beside the four shipped verbs in `internal/cli`". The
  constructor name is taken as the record's own spelling and is used verbatim;
  it is naming guidance, not a contract element, so a later rename is not a
  contract breach.

- **ASSUMPTION-7 (MVV step 10 fixture is the repo's live model).** Step 10 names
  `models/rdr.toml` and describes it as declaring three owned tags with a
  file-backed `[write.rdr-status]`. The test binds that shipped file rather than
  a copy, so the step keeps proving the user outcome if the model changes; the
  step's claim is the user outcome, not encoder coverage (REQ-MVV.10).

---

## QUESTIONS

None. Two clauses initially read as ambiguous; both resolve inside the record,
so neither is escalated:

1. **Is the emptiness-read carrier still open at implementation time?** C1 says
   "UNDECIDED here and is booked as A6", which reads as a genuine fork the
   implementer must pick. Answered by A6 itself, whose `Status` is `Verified`
   and whose evidence names candidate (b) as "the surviving carrier", with the
   sealed-store sub-claim resolved on where the short-circuit lives
   (`::Reader.Read`'s body, not `::load`). Recorded as REQ-29 + ASSUMPTION-1.

2. **Do the two new refusal codes' spellings bind tests?** C1 both sharpens the
   spellings and declares them non-normative, and the code table's own header
   comment says the opposite ("The spellings are normative and the group fixes
   the exit"). Answered inside C1, which names the conflict explicitly and calls
   itself "a deliberate CARVE-OUT from the code table's general rule, not an
   application of it", on the `0028:C1.3` `codeWriteEditRefused` precedent.
   Recorded as REQ-64/REQ-65 + ASSUMPTION-3.
