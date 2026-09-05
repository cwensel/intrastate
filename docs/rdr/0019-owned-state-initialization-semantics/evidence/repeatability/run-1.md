model: claude-opus-5[1m]
variant: full (profile: foundational)

## Read path (what was widened, and why)

Contracts read by id: `0019:C1` (300 lines, the whole normative block),
`0019:MVV`, `0019:S1`–`0019:S12`, plus `0019:D-identity`,
`0019:D-wire-byte-format`, `0019:D-naming`, `0019:D-selection-predicate`
and `0019:RT1`–`0019:RT4`.

WIDENED past the contract spans, deliberately, at four points where C1
left a signature/type/step-order under-determined:

1. `§approach` + `§technical-design` prose — C1 names `groupByWriter`,
   `writerFor` and "the `runFlowSetState` execution shape" without
   giving the verb's own entry point; the Technical Design prose is
   where the two new data flows (plan source, emptiness read) are
   separated.
2. `§mini-checks` (`authority`, `disposition`, `trace` tables) — C1
   fixes payload SEMANTICS ("distinguish the seeded-all case from the
   no-op case") but no field names; the `disposition` table is the only
   place every input class is mapped to an exit/writes/payload triple,
   and `trace` is the only place the step ORDER (bind → class → carrier
   → emptiness → plan → encode/write → read-back) is written down.
   Without it the ordering of the class check versus the carrier check
   would have been a GUESS; it is not.
3. `§implementation-plan` Phase 1 Steps 1–3 — the only place the
   constructor name `newFlowInitStateCmd` appears, and the only
   statement of which step registers the verb.
4. `§failure-modes` (F2, F5) and `§consequences` — F2 is the source for
   the torn-seed payload obligation (list the un-seeded keys, do not
   rewrite committed ones) and F5 for the empty-scalar seed behavior.

`0019:D-wire-byte-format` explicitly says "Field names deferred to
Resolve/Pre-Lock, owned here" — so every payload field NAME below is a
GUESS by the record's own admission, while the payload's required
DISTINCTIONS are contract.

---

## 1. The module's public API

Package: `internal/cli` (verb), with a new emptiness capability whose
home is **UNDECIDED in the record** (booked as assumption A6).

### 1.1 Command surface (contract)

```
intrastate flow init-state
    (--flow <id> | --model <path>)      # shared selection flags, C1 CARRIER
    --artifact <role>=<path> ...        # explicit bindings, repeatable
    [--as json|text]                    # shared output flag
    [--validate]                        # shared ValidateMode gateway
```

No write grammar: `--write`, `--clear`, `--plan` are NOT accepted — the
planned writes ARE the model's `[initial]` assignments (C1 CARRIER).
GUESS: passing `--write`/`--clear` is an unknown-flag refusal from the
flag parser rather than a `flow-*` code — the record does not say.

### 1.2 Go entry point

```go
// internal/cli/flow_init.go   -- file name is a GUESS
func newFlowInitStateCmd() *cobra.Command   // name from Impl Plan Step 1
```

GUESS (shape only, not the name): the record names `newFlowInitStateCmd`
"beside the four shipped verbs", so the signature is whatever the four
shipped constructors use. `*cobra.Command` is a GUESS — the record never
names the CLI framework.

```go
// GUESS: entire signature. C1 names only "the runFlowSetState execution shape".
func runFlowInitState(ctx context.Context, req flowRequest, out responder) error
```

### 1.3 Seed encoder (contract on behavior, GUESS on signature)

C1 fixes the encoder exactly: kind-dispatched, set arm is
`internal/cli/flow_input.go::canonicalSet`, scalar arm is `members[0]`.
It explicitly BANS calling `::canonicalValue` (re-runs argv admission)
and BANS a re-conform pass (the loader already ran `conform`).

```go
// GUESS: name and signature. Behavior is C1-fixed.
func canonicalSeed(decl table.TagDecl, members []string) string {
    // decl.Kind == "set" -> canonicalSet(members)
    // otherwise          -> members[0]
}
```

### 1.4 Emptiness probe — CARRIER UNDECIDED (A6, contract-acknowledged)

C1: "Which carrier serves it — a new `ReadBinding` capability, an
exported `flowbind` cardinality probe, or a `read-state`-family surface
— is UNDECIDED here and is booked as A6". The SEMANTICS are fixed; the
signature is not. This is a genuine RDR silence, not a truncation.

```go
// GUESS: one of C1's three named options, picking the flowbind probe.
// internal/cli/flowbind
func (r Reader) KeyCount() (int, error)   // count of STORE keys, incl. sealedKey
```

C1 forbids one implementation explicitly: substituting "every
`[initial]` key reads absent" for the probe is BANNED (it cannot see a
non-`[initial]` key, including `::sealedKey`).

### 1.5 Types crossed (contract)

- `internal/table/model.go::Model.Initial` — the plan source.
- `internal/table/model.go::TagValue{ Key string; Value []string }` —
  member sequence.
- `internal/resolve/resolve.go::Tag{ Value string }` — the write seam,
  single string. The `[]string` → `string` narrowing IS the encoder.
- `internal/accessor/model.go::Definition.Binding` — the exported field
  the carrier type-switch reads.
- `internal/table/model.go::IsDecisionTable` — the class discriminator.
  C1 BANS re-deriving it from `len(owned)`.

### 1.6 Error modes

Contract-fixed (exit GROUP and distinctness are normative; the literal
code SPELLINGS are explicitly NON-normative — C1 carves these out from
`internal/cli/flow_input.go`'s "spellings are normative" header, on the
`0028:C1.3::codeWriteEditRefused` precedent. **No test may pin the
string.**):

| Condition | Exit | Writes | Code family |
| --- | --- | --- | --- |
| `decision-table` model (`IsDecisionTable`) | 2 | zero, **no accessor runs, no artifact touched** | class refusal, sharpened as `flow-init-class-unsupported` (spelling non-normative) |
| Bound WRITE binding not `*flowbind.Writer` | 2 | zero, no accessor invocation | carrier refusal, `flow-init-carrier-unsupported`, naming the accessor AND the WRITE capability |
| Bound READ binding not `flowbind.Reader` (value) | 2 | zero | same carrier code, naming the READ capability |
| Needed artifact role unbound | 2 | zero (plan-level) | existing artifact-binding `flow-*` family, reused unchanged |
| `[initial]` key with ≠1 declared writer | 2 | zero | `flow-model-invalid` **at LOAD**, before the verb — NOT the writer-routing family (S5 pins this ordering) |
| Read-back disagreement | non-zero (2) | partial, per-key | terminal refusal naming the key PRESENT-AND-UNVERIFIED; re-run does NOT repair and MUST NOT be documented as recovery |
| Read-back incomplete (unreachable locator) | 3 | write possibly applied | inherited `set-state` semantics, distinct from the mismatch class |

Ordering (contract): the carrier refusal is decided AFTER registry
construction (it reads the constructed binding) but BEFORE any accessor
is invoked, and it PREEMPTS the `--allow-commands` refusal
(`Registry.AllowCommands`, also exit 2). A command-backed accessor
yields the carrier code whether or not `--allow-commands` was passed.

Pointer/value spelling is NORMATIVE, because a wrong case matches
nothing: writers are constructed as POINTERS (`&Writer{}`,
`&cmdbind.Writer{}`), readers as VALUES (`Reader{}`, `cmdbind.Reader{}`).
A read case spelled `*flowbind.Reader` refuses EVERY model — S1 goes red.

GUESS: the relative order of the CLASS check and the CARRIER check.
`trace` lists class before carrier; C1 fixes only that class refuses
"before any accessor runs" and carrier "after registry construction".
Both are pre-invocation, so a decision-table model with a command-backed
writer could take either code — the record does not resolve which.

---

## 2. The three most important internal helpers

### 2.1 The store-emptiness predicate (`allBoundStoresEmpty`) — GUESS: name

Responsibility (contract): answer "does EVERY bound artifact carry NO
key". Load-bearing details all fixed by C1:
- Quantifier is **ALL**, never ANY and never per-artifact. A model may
  declare N writers over N roles; an ANY reading would seed into a store
  that still carries a key whenever a sibling artifact was empty — the
  A5 resurrect hazard. S11 is the only test of this.
- Count is over **STORE** keys, not owned keys. A read-back-SEALED
  artifact carries `flowbind::sealedKey` (NUL-prefixed
  `flow.readback-unreachable`), so an artifact whose owned keys are all
  cleared but which is sealed is a ONE-key, NON-EMPTY store that
  init-state DECLINES to seed. S9 is the test.
- The domain is exactly the roles plan-validation already established;
  an unbound needed role is a REFUSAL, never "an empty artifact".
- Scoped to the file-backed carrier; undefined elsewhere.

GUESS: signature `func allBoundStoresEmpty(bindings map[string]accessor.Definition) (bool, error)`.

### 2.2 The carrier gate (`checkFileBackedCarrier`) — GUESS: name

Responsibility (contract): type-switch on
`accessor.Definition.Binding` for BOTH capabilities of every needed
role. Admit ONLY `*flowbind.Writer` on write and ONLY the VALUE
`flowbind.Reader` on read; refuse `*flowbind.EditWriter`,
`*cmdbind.Writer`, `cmdbind.Reader`. Runs after registry construction,
before any invocation, preempting `AllowCommands`. Names which
capability and which accessor failed.

C1 BANS the two tempting alternatives: calling
`flowbind::commandBacked` (it is UNEXPORTED and `internal/cli` cannot
reach it — the citation identifies the registry's branch, not the
verb's call), and re-deriving the carrier from `table.Accessor`'s
`Edit`/`Command` fields (duplicates a discriminator whose writer is the
registry). Exporting `commandBacked` is "not required and is not
authorized here."

### 2.3 The seed encoder (`canonicalSeed`) — GUESS: name

Responsibility (contract): turn each `Model.Initial` `TagValue`
(`[]string`) into the write seam's single `string`, dispatched on the
DECLARED kind. Set → `canonicalSet` (sort, compact, JSON array); every
scalar kind → `members[0]` verbatim. This is the byte-equality source:
applying `canonicalSet` to a scalar would persist `["draft"]` where
`set-state --write status=draft` persists `draft` (six of S4's nine
kinds are scalars). NO re-conform (the loader already ran it and a
re-conform on a normalized value cannot fail), and NOT via
`::canonicalValue` (it re-runs argv admission, which would refuse at
seed time models that load clean — re-creating the very wall the verb
removes).

Honourable mentions the record also fixes but which are reuse, not new:
plan validation (`writerFor` per key + role-bound check, whole plan
before any write) and the per-writer execute+read-back loop (the
`runFlowSetState` shape).

---

## 3. Data model across the boundary

### 3.1 Persisted: the artifact (contract — NO new format)

`D-wire-byte-format`: "none new; artifacts keep
`internal/cli/flowbind` shape". A flat JSON object, key → canonical
string value, values stored VERBATIM. Two contract properties (the rest
of the format is an `flowbind` IMPL-DECISION, explicitly not contract):
key PRESENCE distinguishes an empty set from a cleared key; values are
stored verbatim.

- Empty store: `{}` plus newline, 3 bytes (A5 spike, cited in `trace`).
- Sealed store: exactly `{sealedKey}` — `len == 1`, NON-empty.
- Empty scalar `[initial] note = ""` → `{"note":""}`, reads back
  PRESENT (`Absent:false`), distinct from cleared (F5). This is DECIDED
  behavior for the seed route; the argv asymmetry (`--write note=` is
  refused) stays undecided and is routed to RDR 0002.
- A clear is a key REMOVAL, not a tombstone — so an emptied store and a
  never-written one are the same bytes, and a later init RESEEDS.
  Disclosed residual, not a defect.

### 3.2 In-process: plan source

`Model.Initial []TagValue{ Key string; Value []string }` → per key, the
declared writer (`writerFor`) → `resolve.Tag{ Key, Value string }` with
`Value = canonicalSeed(decl, members)`.

### 3.3 Payload (envelope contract; FIELD NAMES ARE A GUESS)

Contract: rides the 0005:C1 envelope; MUST distinguish the seeded-all
arm from the no-op arm; its scope is EXACTLY the `[initial]` key set
(nothing is claimed about owned keys `[initial]` does not assign); in
the no-op arm it MUST report which `[initial]` keys the store does not
carry. `D-wire-byte-format` says field names are owned here but
deferred — so:

```json
// GUESS: every field name and the nesting. Only the distinctions are contract.
{
  "outcome": "seeded" | "noop",
  "seeded":  [{"role":"state","key":"stage","value":"draft"}],
  "absent":  ["status","gate_passed"]
}
```

GUESS: whether `absent` is emitted (empty) in the `seeded` arm.
GUESS: whether seeded values are echoed at all, or only key names.

---

## 4. Top-level pseudo-code of `runFlowInitState`

Step order below is the `trace` table's, read verbatim; the class/carrier
relative order is the GUESS noted in 1.6.

```
runFlowInitState(req):
  model = selectModel(req.flow | req.model)          # shared flow path; load
                                                     # refuses ≠1-writer [initial]
                                                     # keys as flow-model-invalid,
                                                     # before we run at all (S5)

  if table.IsDecisionTable(model):                   # CLASS arm, C1
      refuse exit 2, flow-init-class-unsupported     # keyed on IsDecisionTable
      # no accessor constructed, no artifact touched (S7)
      # never keyed on len(owned) or on "declares [initial]"

  registry = buildRegistry(model, req.artifacts)     # construction always succeeds

  for role in rolesNeededByWriters(model):           # CARRIER gate, both caps
      if registry.write(role).Binding is not *flowbind.Writer:
          refuse exit 2, carrier code, naming role + WRITE capability
      if registry.read(role).Binding is not flowbind.Reader:   # VALUE, not pointer
          refuse exit 2, carrier code, naming role + READ capability
  # preempts --allow-commands (also exit 2); nothing invoked yet

  plan = []                                          # PLAN, validated whole
  for tv in model.Initial:
      w = writerFor(tv.Key)                          # MUST be exactly one
      if req.artifacts[w.Role] unbound:
          refuse exit 2, existing artifact-binding family, ZERO writes (S6)
      plan.append(Tag{Key: tv.Key,
                      Value: canonicalSeed(decl(tv.Key), tv.Value)})
                                                     # set -> canonicalSet
                                                     # scalar -> members[0]
                                                     # no re-conform, not canonicalValue

  empty = true                                       # PREDICATE: ALL, store keys
  for role in boundRoles(plan):                      # domain = plan-validated roles
      if storeKeyCount(role) != 0:                   # counts sealedKey too (S9)
          empty = false                              # one surviving key anywhere
                                                     # blocks the whole seed (S11)

  if not empty:                                      # NO-OP SUCCESS arm
      absent = [k for k in initialKeys if not present(k)]
      respond exit 0, outcome=noop, absent          # zero writes; torn seeds are
      return                                        # reported, NOT repaired (F2/S12)

  for (writer, tags) in groupByWriter(plan):         # APPLY, per writer
      apply(writer, tags)                            # no cross-writer atomicity
      readback(writer, tags)                         # the only commit check
      if mismatch: refuse, key PRESENT-AND-UNVERIFIED  # re-run will NOT repair
      if incomplete: exit 3, write possibly applied

  respond exit 0, outcome=seeded, plan               # every [initial] key seeded
```

---

## RDR silences reached (each guessed above, listed once)

1. The emptiness probe's carrier and signature — record says UNDECIDED
   (A6). Contract-acknowledged, not an omission.
2. Every payload field name — record says "deferred to
   Resolve/Pre-Lock, owned here".
3. The Go signatures of the verb's run function and helpers, and every
   helper name except `newFlowInitStateCmd`.
4. The file the verb lives in.
5. The relative order of the class arm and the carrier arm when a model
   would trip both.
6. What happens when a write-grammar flag (`--write`/`--clear`/`--plan`)
   is passed to `init-state`.
7. Whether the seeded arm echoes values or only key names, and whether
   `absent` appears in the seeded arm.
