model: claude-fable-5-1
variant: full (profile: foundational)

# Repeatability reconstruction — run 3 — RDR 0019 (`flow init-state`)

Read set: `inspect --json --filter elements,outline`, then `C1`, `MVV`,
`S1`–`S12` by id. Widened past the contract spans into: the Technical
Design preamble (before Normative Contracts), `D-identity`,
`D-wire-byte-format`, `D-naming`, `D-selection-predicate`, `RT1`–`RT4`,
Illustrative Code, Implementation Plan Steps 1–3, Mini-checks (all five
tables), and assumptions `A6`, `A7`. Widening reasons, per item:

- C1 fixes the payload's *content* (seeded-all vs no-op, absent
  `[initial]` keys) but names no field; `D-wire-byte-format` confirms
  field names are deferred ("owned here", not yet fixed). Payload field
  names below are therefore GUESSES, not RDR silences I manufactured.
- C1 says the emptiness read has an undecided carrier; `A6` enumerates
  the three candidates. The helper signature below picks one and is a
  GUESS.
- C1 names the Go types for the carrier gate; `A7` confirms the
  pointer/value shape and that `Definition.Binding` is the field to
  switch on.
- Step 1–3 give the function name `newFlowInitStateCmd` and the
  `runFlowSetState` execution shape; the internal helper names below
  are GUESSES modeled on those.

## 1. Public API

The module is a new verb in the existing `flow` command group, package
`internal/cli`. Its "public API" is the CLI surface plus one exported
constructor; everything else is package-private.

### CLI surface (normative per C1, Step 1, Illustrative Code)

```
intrastate flow init-state (--flow <path> | --model <path>)
    --artifact <role>=<path> [--artifact <role>=<path> ...]
    [--as json|text]            # shared output mode (0005:C1 envelope)
    [--allow-commands]          # accepted for uniformity; preempted by
                                # the carrier refusal (C1 ORDERING)
```

- No write grammar: no `--write`, `--clear`, `--plan`. The plan source
  is `Model.Initial` (C1 CARRIER).
- Selection flags are the shared `--flow`|`--model` pair (C1).
- `--artifact role=path` is explicit and per-invocation; an unbound
  needed role is a refusal, never an artifact treated as empty (C1).

### Go constructor (Step 1)

```go
// newFlowInitStateCmd registers the fifth flow verb beside next,
// resolve, read-state, set-state.
func newFlowInitStateCmd() *cobra.Command   // GUESS: cobra, matching the
                                            // four shipped verbs' shape
```

### Exit groups and refusal codes (C1, disposition table)

| Outcome | Exit | Writes | Code |
| --- | --- | --- | --- |
| Every bound artifact empty → seeded all `[initial]` keys | 0 | all | (success envelope) |
| Any bound artifact non-empty (incl. sealed, torn, post-clear, fully seeded) | 0 | zero | (success envelope, no-op arm) |
| `decision-table` model (`table.IsDecisionTable`) | 2 | zero, no accessor run | `flow-init-class-unsupported` (spelling NON-normative; group + distinctness normative) |
| Any needed role's write binding not `*flowbind.Writer`, or read binding not `flowbind.Reader` | 2 | zero, no accessor run | `flow-init-carrier-unsupported` (spelling NON-normative), names capability + accessor |
| Needed artifact role unbound | 2 | zero | existing artifact-binding `flow-*` code (unchanged) |
| `[initial]` key with ≠1 writer | 2 | zero | `flow-model-invalid` at LOAD — never reaches the verb |
| Model selection failure | 2 (GUESS: same group as set-state) | zero | existing model-selection `flow-*` code |
| Read-back mismatch after a write | non-zero (GUESS: exit 3, the read-back family's group) | partial, per-key | existing read-back code; payload names key PRESENT-AND-UNVERIFIED |

Ordering fixed by C1: model load → class refusal → registry
construction → carrier refusal (before any accessor is invoked, and
preempting `Registry.AllowCommands`) → role-binding validation → emptiness
read → plan validation → per-writer apply with read-back.

### Success payload (content normative per C1; FIELD NAMES ARE A GUESS —
`D-wire-byte-format` defers them)

```jsonc
// rides the 0005:C1 envelope, e.g. {"ok": true, "verb": "init-state", ...}
{
  "arm": "seeded" | "noop",            // GUESS name; C1 requires the two
                                       // arms be distinguishable
  "seeded": ["stage", "status", ...],  // keys written this run; empty on noop
  "absent": ["gate_passed", ...],      // [initial] keys the store does NOT
                                       // carry (informational; on the
                                       // seeded arm this is [])
  "artifacts": {"state": "state.json"} // GUESS: bound roles echoed back
}
```

Scope is exactly the `[initial]` key set; the payload claims nothing
about owned keys `[initial]` does not assign (C1).

## 2. Three most important internal helpers

Names are GUESSES; responsibilities are drawn from C1, D-selection-
predicate, Steps 2–3, and the Mini-checks `authority` table.

### 2a. `initCarrierCheck(reg *accessor.Registry, needed []string) error`

Responsibility: the carrier gate over BOTH capabilities. For every role a
needed writer names, fetch the constructed write binding and the
role-resolved read binding (`Registry.readerFor` over `CapRead`) and
type-switch on `accessor.Definition.Binding`:

- write side admits ONLY `*flowbind.Writer` (pointer); refuses
  `*flowbind.EditWriter` and `*cmdbind.Writer`;
- read side admits ONLY `flowbind.Reader` (value); refuses
  `cmdbind.Reader`. A case spelled `*flowbind.Reader` matches nothing and
  refuses every model — S1 is the detector.

Runs AFTER registry construction and BEFORE any accessor invocation,
preempting the `--allow-commands` refusal. Returns a refusal naming the
capability (read/write) and the accessor id, exit group 2. MUST NOT
re-derive the carrier from `table.Accessor.Edit`/`.Command` and MUST
NOT call the unexported `flowbind.commandBacked`.

### 2b. `storeIsEmpty(ctx, reader flowbind.Reader, art Artifact) (bool, error)`

Responsibility: the second new data flow — answer "does this bound
artifact carry ANY key" for the file-backed carrier only, counting STORE
keys (so `flowbind.sealedKey` counts) not owned keys, and never
substituting "every `[initial]` key reads absent". The carrier is
UNDECIDED in the RDR (A6, Pending). GUESS: implement as an exported
`flowbind` cardinality probe, e.g. `func (r Reader) Empty(art) (bool,
error)` loading the store and returning `len(store) == 0` (an absent
file loads as a zero-key store, so absent and emptied-by-clears are
indistinguishable — C1's disclosed residual). The verb calls it for
EVERY bound needed role and ANDs the answers (ALL quantifier, S11).

### 2c. `planInitial(m *table.Model, req flowRequest) (map[writerID][]resolve.Tag, error)`

Responsibility: build and validate the WHOLE plan before any write.
For each `Model.Initial` assignment (`table.TagValue{Key, Value
[]string}`): resolve the writer via `writerFor` (exactly-one is already
enforced at load by `checkAccessorBindings`, so the verb owes no arity
check, only the MUST-NOT-route-to-two statement); confirm the writer's
role is bound in `req.artifacts[role]` (S6 arm, refusal in the
artifact-binding family); encode the value kind-dispatched — `decl.Kind
== "set"` → `canonicalSet(members)` (sort, compact, JSON array), every
scalar kind → `members[0]` verbatim — into `resolve.Tag{Key, Value
string}`. NOT `canonicalValue` (re-runs argv admission, refusing
loader-admitted spellings) and NOT a re-conform pass (cannot fail, buys
nothing). Group by writer (`groupByWriter`). Any error → zero writes
committed.

## 3. Data model at the boundary

Persisted: nothing new. Artifacts keep the `internal/cli/flowbind` JSON
shape (`{"key":"value",...}`, values stored verbatim, key PRESENCE
distinguishes empty-set from cleared; a `sealedKey` NUL-prefixed marker
may be present). An empty store is `{}`+newline, 3 bytes (trace table).

Passed across the boundary:

```go
// Source of the plan (shipped, loader-owned):
table.Model.Initial  []table.TagValue     // TagValue{Key string; Value []string}
table.Model.Writers  []accessor.Definition // GUESS on element type; per-writer Role
table.IsDecisionTable(m) bool             // class discriminator; sole key for the class arm

// Write seam (shipped):
resolve.Tag{Key string; Value string}     // one string per key, canonical form

// Carrier discrimination (shipped, exported):
accessor.Definition.Binding accessor.Binding // dynamic type: *flowbind.Writer |
                                             // *flowbind.EditWriter | *cmdbind.Writer
                                             // (write); flowbind.Reader | cmdbind.Reader (read)

// New (GUESS — A6 undecided): emptiness answer for the file-backed reader
type emptiness struct{ Role string; Empty bool }

// Payload (GUESS on field names; content fixed by C1):
type initStatePayload struct {
    Arm       string   `json:"arm"`       // "seeded" | "noop"
    Seeded    []string `json:"seeded"`
    Absent    []string `json:"absent"`    // [initial] keys the store lacks
    Artifacts map[string]string `json:"artifacts"`
}
```

Value encoding is JDR 0001 §D13 canonical form dispatched on declared
kind. Numeric spellings are the loader's (`1` not `1.0`, `5` not `+5`),
which RT3 scopes its byte-identity claim to. An empty scalar `""` seeds
and reads back PRESENT (F5).

## 4. Top-level pseudo-code of `runFlowInitState`

```
func runFlowInitState(cmd, req):
  m, err := loadAndSelectModel(req.flow, req.model)        // shared flow path
  if err: return respond(refusal(existing model-selection code, exit 2))
  // [initial] keys with ≠1 writer already refused here as flow-model-invalid (S5)

  if table.IsDecisionTable(m):                              // class arm, keyed on the
    return respond(refusal(codeInitClassUnsupported, 2))    // discriminator ALONE; no
                                                            // accessor constructed or run
  reg, err := flowbind.Registry(m, req.allowCommands)       // construction always succeeds
  needed := rolesNamedBy(writersServing(m.Initial))
  if err := initCarrierCheck(reg, needed); err != nil:      // BOTH capabilities; before any
    return respond(refusal(codeInitCarrierUnsupported, 2,   // accessor invocation; preempts
                           capability, accessorID))        // Registry.AllowCommands refusal
  for role in needed:
    if req.artifacts[role] == "":                           // S6: reached at the verb
      return respond(refusal(existing artifact-binding code, 2))   // zero writes

  allEmpty := true
  for role in needed:                                       // ALL quantifier (S11)
    empty, err := storeIsEmpty(ctx, reg.readerFor(role), req.artifacts[role])  // A6 carrier
    if err: return respond(refusal(existing read family))
    allEmpty = allEmpty && empty                            // STORE keys, sealedKey counts (S9)

  if !allEmpty:                                             // torn / post-clear / seeded / sealed
    absent := initialKeysNotPresent(ctx, reg, m.Initial, req)   // informational only
    return respond(ok(payload{Arm: "noop", Seeded: nil, Absent: absent}))  // ZERO writes;
                                                            // never repair a torn seed (S12)

  plan, err := planInitial(m, req)                          // encode kind-dispatched, group
  if err: return respond(refusal(..., 2))                   // by writer; ZERO writes on error

  seeded := []
  for writer, tags in plan:                                 // runFlowSetState execution shape
    result, err := executor.Write(ctx, writer, req.artifacts[writer.Role], tags)
    if err or result.readBackMismatch:                      // read-back is the only commit
      return respond(refusal(existing read-back code,       // check; no cross-writer atomicity
                             keys PRESENT-AND-UNVERIFIED))  // re-run is a no-op, NOT recovery
    seeded += keys(tags)
  return respond(ok(payload{Arm: "seeded", Seeded: seeded, Absent: []}))
```

## GUESS ledger (for the diff)

1. Payload field names (`arm`, `seeded`, `absent`, `artifacts`) — RDR
   defers them (D-wire-byte-format).
2. Emptiness carrier: exported `flowbind` cardinality probe — RDR leaves
   it undecided among three (A6).
3. Exit group of read-back mismatch (non-zero, guessed 3) — C1 says
   "distinct terminal refusal", fixes no number.
4. Exit group of model-selection failure (guessed 2, "shared classes
   reuse existing codes").
5. Internal helper names and the `Model.Writers` element type.
6. cobra as the command framework (inferred from "beside the four
   shipped verbs").
7. `--allow-commands` accepted on this verb's flag set (C1 discusses its
   preemption, implying the flag exists here; not stated outright).
