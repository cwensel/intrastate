model: claude-sonnet-5
variant: full (profile: foundational)

# Repeatability Reconstruction — RDR 0019 (Owned-state initialization semantics)

Source: RDR 0019 read via the projector only (`inspect --json --filter elements,outline`,
then `--select` on `0019:C1`, `0019:MVV`, `0019:D-identity`, `0019:D-wire-byte-format`,
`0019:D-naming`, `0019:D-selection-predicate`, `0019:RT1-4`, `0019:JC1`, `0019:S1-S12`,
`0019:F1-F5`, `0019:A1-A8`, `0019:BR1-5`, plus the `§illustrative-code`,
`§load-bearing-decisions`, and `§approach` outline sections). The record's only `C`
element is `C1`; `MVV` and `S1`-`S12` are one element apiece. I widened past the
contract span into `§approach` (for the `groupByWriter → executor → read-back` pipeline
statement, which C1 references but does not restate) and into the full `A1`-`A8` set
(for substrate facts C1's prose leans on but does not re-derive: no runtime reader of
`Model.Initial` exists yet, the writer-arity guarantee lives at load not lint, the
carrier-discrimination types are exported field/exact types, and A6 — the emptiness-read
carrier — is **Pending**, i.e. genuinely undecided by this record). I did not widen into
Background/Research Findings/Trade-offs prose beyond what A1-A8, BR1-5 already restate;
those sections looked like narrative context, not additional constraints.

---

## 1. Public API

### CLI surface

A new verb, `init-state`, added to the existing `flow` group (extends `0005:C1`'s
verb enumeration and `0005:D-naming`'s verb list by exactly this spelling — the group
goes from four verbs to five). Signature (illustrative form the record gives):

```
intrastate flow init-state --model <path> \
    --artifact <role>=<path> [--artifact <role>=<path> ...] \
    [--as json]
```

- **Selection flags**: the shared `--flow`/`--model` flags every `flow` verb takes.
- **Binding flags**: explicit `--artifact role=path` bindings, one per role the model's
  declared write/read accessors need.
- **No write grammar**: unlike `set-state`, `init-state` takes no `--write`/`--clear`/`--plan`
  — its planned writes are exactly the model's `[initial]` assignments, sourced from the
  loader-normalized `Model.Initial`, not from argv.
- Inherits `0005:C1`'s per-verb MUSTs: `respond.ValidateMode` first, `respond.OK`,
  `respond.Fail`, `SilenceErrors`/`SilenceUsage` — same envelope every `flow` verb uses.

### Return / exit contract

| Condition | Exit | Payload |
|---|---|---|
| Empty store (every bound artifact carries no key) | 0 | seeded-all: every `[initial]` key written, canonical form |
| Non-empty store (any key present, including a read-back seal) | 0 | no-op: zero writes; lists the `[initial]` keys the store does not carry |
| `[initial]` key with 0 or >1 declared writer | 2 (`flow-model-invalid`, at LOAD, before the verb runs) | writer-arity refusal — never reaches `init-state`'s own code paths |
| Required artifact role unbound | 2 (existing artifact-binding family) | zero writes committed |
| Model is `decision-table` (`IsDecisionTable` true) | 2, new code `flow-init-class-unsupported` | class refusal, before any accessor is invoked, no artifact touched |
| Any bound write or read accessor for a needed role is not file-backed (`*flowbind.Writer` write / `flowbind.Reader` read) | 2, new code `flow-init-carrier-unsupported` | carrier refusal, names the accessor and which capability failed, zero writes, preempts `--allow-commands`' own refusal |
| Read-back after a seed write disagrees with what was written | 2 (existing read-back-mismatch family) | terminal refusal naming the key present-and-unverified — distinct from an incomplete-read-back exit 3 |
| Read-back reports a requested key UNREADABLE (e.g. sealed artifact) | inherited `set-state` semantics (exit 3 class) — but note: a sealed, all-cleared store is itself a **no-op success (exit 0)** at the predicate level, since the seal makes the store non-empty *before* any write is attempted | — |

Both new exit-2 codes are in the `flow-<subject>-<condition>` family; per the record,
their literal spellings are explicitly **non-normative** (carve-out precedent
`0028:C1.3`) — only the exit group (2) and their mutual distinctness are fixed. No test
may pin the string.

### Types (Go-level, as the record cites them — GUESS on anything not literally spelled)

- `internal/table/model.go::Model.Initial` — `[]TagValue` (existing type, read not written
  by this RDR); `TagValue.Value []string` (member sequence).
- `internal/table/model.go::IsDecisionTable(*Model) bool` — existing discriminator, reused
  verbatim as the class gate; **MUST NOT** be re-derived from `len(owned)`.
- `internal/accessor/model.go::Definition.Binding` — exported field, exported `Binding`
  interface; the verb type-switches on its dynamic type. Write side admits only
  `*flowbind.Writer` (pointer); read side admits only `flowbind.Reader` (value). GUESS:
  the type-switch itself (`switch b := def.Binding.(type) { case *flowbind.Writer: ... }`)
  is not given verbatim in the record — I infer the shape from the pointer/value
  distinction the record states as load-bearing.
- `internal/cli/flowbind/flowbind.go::sealedKey` — the NUL-prefixed
  `flow.readback-unreachable` marker (existing, from `0004:C13`/`C14`).

### Error modes summary

Five refusal classes: (1) load-time writer-arity (existing, reused), (2) artifact-binding
unbound-role (existing, reused), (3) new class refusal for decision-table models, (4) new
carrier refusal for non-file-backed write or read accessors, (5) existing read-back-mismatch
family. All five commit zero writes. No repair/recovery path exists inside the verb itself
for a torn multi-writer seed or a present-and-wrong read-back key — recovery is explicit
`set-state`, or discarding and re-running `init-state` (only safe if the artifact holds
nothing but `[initial]` keys).

---

## 2. Three most important internal helpers

1. **The empty-store predicate (unnamed in the record — GUESS at a name:
   `isStoreEmpty` / `allBoundArtifactsEmpty`)**. Responsibility: decide, per invocation,
   whether the WHOLE bound set of artifacts is eligible to be seeded. Semantics fixed by
   C1/`§load-bearing-decisions`: ALL bound artifacts must carry **zero store keys** (not
   zero *owned* keys — a read-back-sealed artifact with all owned keys cleared still
   counts as one key, via `sealedKey`), never ANY/per-artifact. This is explicitly the
   safety property the whole design rests on (A5) — a sibling artifact being empty must
   never license seeding over one that still holds a key (the A5/resurrect hazard). The
   record discloses this predicate is **not computable on the shipped seam as-is** (A6,
   status Pending): `ReadBinding.Read` is key-scoped with no cardinality answer, and
   `flowbind::store` is package-private. GUESS: since A6 is unresolved, the actual
   carrier this helper reads through (a new `ReadBinding` capability vs. an exported
   `flowbind` cardinality probe vs. a `read-state`-family surface) is an RDR silence —
   I mark it UNDECIDED rather than guessing a specific signature.

2. **The carrier gate (unnamed — GUESS: `checkAccessorCarrier` or similar)**.
   Responsibility: for every needed role, type-switch the constructed `Definition.Binding`
   on BOTH the write binding and the read binding, admitting only `*flowbind.Writer`
   (write) and `flowbind.Reader` (value, read); refuse with `flow-init-carrier-unsupported`
   naming the accessor and the failing capability otherwise. Runs AFTER registry
   construction but BEFORE any accessor is invoked, and PREEMPTS the `--allow-commands`
   refusal (so a command-backed accessor always yields the carrier code, regardless of
   whether `--allow-commands` was passed). This ordering is normative and tested by S10.

3. **The seed encoder (unnamed — GUESS: reuses `set-state`'s SET arm as
   `internal/cli/flow_input.go::canonicalSet` for `decl.Kind == "set"`, and
   `members[0]` verbatim for every scalar kind)**. Responsibility: render each
   `Model.Initial` assignment (a `[]string` member sequence) into the same canonical wire
   form `set-state --write` would produce, WITHOUT re-running argv admission checks
   (`::canonicalValue` is explicitly the wrong call — it would re-refuse values the
   loader's `[initial]` admission set allows but argv's narrower admission set does not:
   bare scalar for a set tag, array for a scalar tag, empty scalar). This is the function
   that makes read-back equality byte equality; the record explicitly says a re-conform
   pass buys nothing because loader conformance is already held (`::loadInitial`'s
   `conform`).

---

## 3. Data model (persisted / boundary-crossing)

**Persisted (artifact/store, file-backed JSON carrier only — the only carrier this RDR
admits)**:

- The artifact is a flat key→value JSON object (existing `flowbind` shape, unchanged by
  this RDR): `internal/cli/flowbind/flowbind.go::store` — package-private, GUESS shape
  `map[string]string` given the record's description of key PRESENCE distinguishing
  empty-set from cleared, and values stored VERBATIM.
- Presence/absence is the only signal: a cleared key is a map deletion (`delete(s, t.Key)`),
  never a tombstone — clearing the last key yields the same zero-key store as a
  never-written artifact (byte-identical per A5's spike: `{}` + newline, 3 bytes).
- `sealedKey` — a NUL-prefixed `flow.readback-unreachable` marker string, persisted as a
  regular key in the same store when the writer's declared locator is unreachable
  (`::unreachable(w.Path)`); dropped by the next write whose locator IS reachable.
- The seeded-value form for each key is the canonical wire form: JSON array (sorted,
  compact) for `kind="set"`, the single member string verbatim for every scalar kind
  (`enum`, `bool`, `int`, `scalar`). No new wire format is introduced (`D-wire-byte-format`:
  "none new").

**Passed across the CLI boundary (payload, not persisted)**:

- Success envelope rides the existing `0005:C1` envelope (`respond.OK`/`respond.Fail`).
- Seeded-all payload: which `[initial]` keys were written and their values (canonical
  form) — GUESS at exact field names, since `D-wire-byte-format` explicitly defers field
  naming to "Resolve/Pre-Lock, owned here" (i.e. undecided in this record on purpose).
- No-op payload: MUST distinguish itself from the seeded-all case, and must list exactly
  the `[initial]` keys the store does not carry (informational — makes a torn or
  post-clear state visible without repairing it). Scope is exactly the `[initial]` key
  set; the verb claims nothing about other owned keys.
- Carrier/class refusal payloads: name the failing accessor and (for carrier refusals)
  which capability — read or write — failed.

**In-memory model data consumed (not persisted by this verb, read-only input)**:

- `Model.Initial []TagValue` — the loader-normalized seed source (loader already ran
  `conform(decl, "eq", members)` over every assignment).
- `Model.Writers` — existing per-role writer declarations, iterated the same way
  `runFlowSetState` already does.

---

## 4. Top-level pseudo-code (20-40 lines)

GUESS throughout for exact function/variable names not given verbatim by the record;
the STEP ORDER and each gate's placement are load-bearing and taken directly from C1,
`§approach`, and the MVV/S-series scenarios.

```
func runFlowInitState(req) (Result, error):
    # 1. Standard flow-verb entry (0005:C1 inherited)
    respond.ValidateMode(req)
    model := loadModel(req.model)          # existing loader; refuses at LOAD if any
                                            # [initial] key has writerCount != 1
                                            # (CatMalformedAccessorBinding) or if a
                                            # decision-table model somehow carries
                                            # [initial] (unconstructible - refused earlier)

    # 2. Class gate - BEFORE any accessor construction/invocation
    if IsDecisionTable(model):
        return Fail(flow-init-class-unsupported)   # exit 2, zero accessors touched

    # 3. Resolve bindings for every role a needed writer/reader names
    registry := buildRegistry(model, req.artifacts)
    for role in neededRoles(model):
        if role not in req.artifacts:
            return Fail(existing artifact-binding-unbound family)  # exit 2, zero writes

    # 4. Carrier gate - AFTER registry construction, BEFORE any accessor invoked,
    #    PREEMPTS --allow-commands refusal
    for role, binding in registry.bindings():
        if writeBinding(role) is not *flowbind.Writer:
            return Fail(flow-init-carrier-unsupported, role, "write")
        if readBinding(role) is not flowbind.Reader:
            return Fail(flow-init-carrier-unsupported, role, "read")

    # 5. Emptiness predicate - ALL bound artifacts, STORE keys not owned keys
    #    (carrier for this read is UNDECIDED per A6 - GUESS: routed through
    #    whichever ReadBinding/flowbind surface A6 eventually lands)
    if not allBoundArtifactsEmpty(registry):     # counts sealedKey too
        missing := initialKeysNotInAnyStore(model, registry)
        return OK(no-op: zero writes, missing)   # exit 0, NO write attempted

    # 6. Build the seed plan: every [initial] key -> its one declared writer
    plan := groupByWriter(model.Initial, model.Writers)  # existing set-state pipeline

    # 7. Encode each seed value canonically (no argv re-validation)
    for assignment in plan:
        assignment.wire = encodeSeed(assignment)  # canonicalSet for kind=set,
                                                    # members[0] otherwise

    # 8. Execute exactly on set-state's writer-routing / read-back path
    result := executor.Apply(plan)   # per-writer commit + read-back, no
                                      # cross-writer atomicity (torn seed possible)
    for key, outcome in result:
        if outcome.readBackMismatch:
            return Fail(existing read-back-mismatch family, key)  # exit 2,
                                                                    # present-and-unverified

    return OK(seeded-all: plan.keys, plan.values)  # exit 0
```

Notes load-bearing to this shape: step 5 (emptiness) is a plan-level, all-or-nothing gate
over the WHOLE bound set — there is no per-key merge anywhere in the pipeline. Step 8
reuses `set-state`'s existing writer-routing/read-back machinery verbatim (no new commit
mechanism); a failure partway through step 8 can leave a torn (partially seeded) artifact,
which is an accepted, disclosed residual (F2) — a subsequent re-run is a no-op reporting
the missing keys, never a repair.

---

## Explicit GUESSes (consolidated)

- Exact function/type names not given verbatim (`isStoreEmpty`, `checkAccessorCarrier`,
  `encodeSeed`, the type-switch code shape) — the record names some real symbols
  (`groupByWriter`, `writerFor`, `canonicalSet`, `sealedKey`, `IsDecisionTable`) but not
  a full call graph for `init-state` itself, since it does not exist yet (A5's note:
  "`init-state` does not exist yet").
- Payload field names — explicitly deferred by the record itself to Resolve/Pre-Lock,
  not merely an omission on my part.
- The carrier for the emptiness read (A6, status **Pending**) — genuinely undecided in
  the record; I did not invent a specific API for it.
- Exact literal spellings of the two new exit-2 codes — explicitly non-normative per the
  record's own carve-out.
- The in-memory `store` type as `map[string]string` — inferred from described behavior
  (presence/absence, verbatim values), not a cited type declaration.
