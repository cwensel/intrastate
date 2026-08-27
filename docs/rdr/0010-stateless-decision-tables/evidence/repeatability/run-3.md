model: claude-haiku-4-5-20251001
variant: full (profile: foundational)

## 1. Public API (Function Signatures, Types, Error Modes)

### Modified Types
- **`Model.Class`**: New optional string field on `table.Model`. Accessor with zero value `"state-machine"` (RDR silent: assumed via `model.go::Model` accessor pattern, not explicitly in contracts, but C1 specifies the accessor exists). Admitted values: `"state-machine"`, `"decision-table"`. Any other value: `malformed model declaration` error.

- **`Row.Emit`**: New field `[]EmitValue` on `table.Row`, where `EmitValue` is a **GUESS** new struct `{Key string, Value string}` (structured differently from `TagValue` to avoid drag-in of set-related operations per C3). Carries the authored `[rule.emit]` block as a key-sorted sequence.

- **`resolvePayload.Emit`**: New JSON field `emit` on CLI payload, ordered immediately after `rule`. Type: `map[string]string`, keys in byte order, empty `{}` when no block authored (never `null`). Error modes: none new at resolve; the block is optional per C3.

- **`sourceModel.Class`**: New optional TOML field `class` on source model (GUESS: assumed from C1's requirement that the class is read in `loadModelHeader`). Type: string.

- **`sourceRule.Emit`**: New optional TOML field `emit` on source rule (GUESS: `map[string]string`, per C3's decoder constraint). Refuses `malformed TOML` on non-string values (nested table included).

### New Error Modes
- **`malformed model declaration`** (existing category, new arm): `decision-table` declaring non-zero owned tags. Detail: token `owned=<n>`.
- **`malformed TOML`** (existing category, applied to emit): non-string value in `[rule.emit]` block (decoder type arm).
- Existing errors apply unchanged: `unknown tag` (emit key in match/guard), `write to non-owned tag` (attempt to write undeclared tag), `graph-dangling-edge` (terminal predicate over non-owned tag in decision table).

### Modified Error Modes (Conditional)
- **`malformed rule shape`**: No-write-block refusal is now scoped to `"state-machine"` class only (C2). Decision-table rules load with empty `RequiresOwned`, `Writes`, `NextTags`.

### Lint Finding Code Changes
- **`graph-unprovable-coverage`**: New discriminator: `reason = "no-participating-dimension"` for decision-table groups with zero guard dimensions (C5). Existing reasons: `dimension-not-finite`, `tag-not-single-valued`, `row-can-refuse` (all preserved).
- **`graph-dangling-edge`**: Terminal-predicate arm stays live for decision tables to enforce C2's `terminal` prohibition. Missing-root arm augmented: `(class == "decision-table" || len(Initial) > 0)` — the silenced root arm for decision tables only (C5).

### Functions (GUESS: Inferred from Contracts)
- **`internal/table/load.go`**: 
  - Class agreement check runs in a step at or after `loadTags` (after `loadModelHeader`, `loadOutcomes`, `loadTags` per C1). 
  - `loadModelHeader`: Reads `class` field.
  - `normalizeRule`: Conditions no-write-block refusal on class.
- **`internal/graphlint/reach.go::reach`**: Seeds traversal at empty owned-state node when `class == "decision-table"` (C5).
- **`internal/graphlint/analysis.go::checkDanglingEdge`**: Missing-root arm checks `(class == "decision-table" || len(Initial) > 0)`.
- **`internal/graphlint/coverage.go::checkCoverage`**: Inside `len(dims) == 0` branch, emits `graph-unprovable-coverage` with `reason = "no-participating-dimension"` ahead of `emitCoverageArms` (C5).
- **`internal/cli/flow_resolve.go::rowByID`**: Joins emit back by rule id (first-match unchanged); emit is per-rule, shared by all expanded rows from that rule.
- **`internal/table/dump.go::dumpColumns`**: Appends `emit` column last (C3).
- **`internal/table/normalize.go::expand`**: Carries `Emit` through for every expanded row from multi-member match atoms (C3).

---

## 2. Three Most Important Internal Helper Functions

**GUESS: Inferred from RDR narrative**

### 1. **`internal/table/load.go::classOwnershipAgreement` (Conjectured)**
**Responsibility**: Validates C1's one-directional agreement after tag table is loaded. Reads `model.Class` (set in `loadModelHeader`) against `len(model.Tags[owned])`. Refuses `malformed model declaration` with detail `owned=<n>` if class is `"decision-table"` and owned count > 0. No refusal if class is `"state-machine"` or absent (defaults to state-machine).

**Load Position**: Called at or after `loadTags`, before `checkAccessorBindings` (run's last step), so undeclared-tag refusals precede class disagreement.

### 2. **`internal/graphlint/coverage.go::emitUnprovableDimensionOrNone` (Conjectured)**
**Responsibility**: Emitted inside the `len(dims) == 0` branch of `checkCoverage`. For decision-table models, emits one `graph-unprovable-coverage` finding with `reason = "no-participating-dimension"`. For state machines, does nothing (zero-dimension groups are legitimate and stay silent per A13). Precedes `emitCoverageArms` so the unprovable group is reported instead of being closed by a bare escape row.

**Keying**: Checks `model.Class == "decision-table"` and `len(dims) == 0`.

### 3. **`internal/cli/flow_resolve.go::selectRowWithEmit` (Conjectured)**
**Responsibility**: After kernel selection, joins the selected row's `emit` block back from the normalized table by rule id using `rowByID` (first-match). Populates `resolvePayload.Emit` with the key-sorted `map[string]string`, or empty `{}` if no block authored. Emit is never passed to the kernel; it is table-scoped.

**Position**: Runs after kernel returns the selected row; no artifact invocation needed for emit to populate.

---

## 3. Data Model for Persisted and Boundary-Crossed Data

### Table-Layer Persisted Model
**TOML (wire) → Normalized (internal) → Kernel (resolve boundary)**

#### Source (`sourceModel`, `sourceRule`)
```
[model]
id = "navigator"
version = 1
class = "decision-table"  # NEW: optional, admitted {"state-machine", "decision-table"}

[[rule]]
id = "draft-propose"
[rule.emit]               # NEW: flat TOML table, string values only
next = "propose"
key2 = "value2"
```

#### Normalized (`table.Model`, `table.Row`)
```go
type Model struct {
  ID    string
  Class string  // NEW: "state-machine" or "decision-table"; zero value "state-machine"
  Rows  []Row
}

type Row struct {
  ModelID       string
  RuleID        string
  Suffix        []string
  Outcome       string
  Atoms         []Atom
  Gate          []string
  NextTags      []TagValue
  Writes        []TagValue
  RequiresOwned []string
  Escape        []string
  Emit          []EmitValue  // NEW: key-sorted, each {Key string, Value string}
}

type EmitValue struct {  // NEW type (GUESS)
  Key   string
  Value string
}
```

#### Kernel Boundary (`resolve.Row`, `resolve.Table`)
- **Unchanged**: Kernel row has no emit field. Emit never crosses to kernel.
- Normalized row's `Emit` is joined back by rule id (`RuleID`) at resolve time in CLI layer.

#### CLI Payload (`resolvePayload`)
```json
{
  "model": "navigator",
  "revision": "1",
  "observed": {"status": "Draft"},
  "owned": {},
  "readers": [],
  "outcome": "locate",
  "rule": "draft-propose",
  "emit": {"next": "propose", "key2": "value2"},  // NEW: map[string]string, keys in byte order
  "next": {},
  "writes": {},
  "clear": [],
  "escaped": false
}
```

### Lint Model (Reachability Graph)
**Decision-table root: Augmented per C5**
- Seeding predicate for `reach()`: `class == "decision-table" || len(Initial) > 0`
- Reachable set for decision-table: exactly the empty owned-state node; all selection contexts reachable from there.
- Product dimensions: only guard atoms (`guard.all`/`guard.unless`) participate; match atoms scope the group without contributing to product.
- Zero-dimension group handling: NEW finding `graph-unprovable-coverage` with `reason = "no-participating-dimension"` for decision tables; state machines silently accept zero-dimension groups (legitimate).

### Dump Column Order (Presentation)
`internal/table/dump.go::dumpColumns` appended order:
- Existing columns: `[Model, Rule, Outcome, Gate, Initial, Terminal, Match, Guard, Writes, Clear, Escape]`
- **NEW: `Emit` appended last** (C3 requirement: append last, after escape, so existing indices stable).

#### Emit Column Rendering
- Row without `[rule.emit]` block: empty bracket `[]`
- Row with block: `key1=value1; key2=value2; …` in key order, bracketed like `writes`.
- Example: `[next=propose; key2=value2]`

---

## 4. Top-Level Pseudo-Code of Main Operation (20–40 lines)

### Main Operation: `intrastate flow resolve --model <file> --outcome <o> --tag k=v… [--artifact role]`

```
// Load and Lint (shared across classes)
1.  TOML → decodeStrict → sourceModel, sourceRules
2.  IF model.version ≠ 1: REFUSE version mismatch
3.  loadModelHeader() → read class (default "state-machine"), id, version, metadata
4.  loadOutcomes() → parse outcomes
5.  loadTags() → parse [tags.*] declarations
6.  classOwnershipAgreement() → IF class="decision-table" AND owned.count > 0: REFUSE malformed model declaration, detail "owned=<n>"
7.  loadAccessors() → read [read], [write], [gate] bindings (no writer for decision-table by C2)
8.  normalizeRules() → FOR each rule:
        IF rule has no write block:
          IF class="state-machine": REFUSE malformed rule shape
          IF class="decision-table": CONTINUE (empty NextTags, Writes, RequiresOwned)
        normalize atoms, emit, escape → Row with Emit []EmitValue (sorted by key)
        FOR each multi-member match atom: expand(rule) → multiple rows, each carrying same Emit
9.  rows → sorted by (ModelID, RuleID, Suffix)
10. lint() → lint the model (returns early if missing-root and not decision-table):
      reach() uses seed predicate: class="decision-table" OR len(Initial) > 0
      checkCoverage(): IF len(guard.Dimensions()) == 0:
                        IF class="decision-table": emit graph-unprovable-coverage(reason="no-participating-dimension")
                        ELSE: CONTINUE (zero-dimension group legitimate for state-machine)
      checkDanglingEdge(): missing-root arm uses class="decision-table" OR len(Initial) > 0

// Resolve (decision-table specific behavior)
11. assembled_view = {observed tags, recognized outcome, no owned, no readers invoked}
12. kernel.Exact(table, view) → SELECT row matching all guards and match atoms (or no_match)
13. IF matched row: selected_rule_id = row.RuleID
    ELSE IF escape row rescues no_match: selected_rule_id = escape.RuleID, escaped=true
    ELSE: RETURN flow-no-match refusal
14. runGates(selected_rule, gates) → IF gate denied: RETURN flow-gate-denied
15. rowByID(selected_rule_id) → SCAN normalized rows, FIRST MATCH (all expanded rows from same rule share Emit)
16. resolvePayload = {
      model: model.ID,
      outcome: recognized_outcome,
      rule: selected_rule_id,
      emit: selected_row.Emit as map[string]string (empty {} if absent),  // NEW
      observed: view.observed,
      owned: {} (decision-table always empty),
      next: {} (decision-table always empty),
      writes: {} (decision-table always empty),
      clear: [] (decision-table always empty),
      readers: [],
      gates: gate results
    }
17. text_render(payload) → generic payload renderer from 0005:
      FOR each path-qualified leaf: emit one line
      emit.<key>: <value> for each key; emit: (none) if empty (C4, 0005 behavior unchanged)
18. json_render(payload) → Marshal to JSON (emit field position: immediately after rule)
19. RETURN exit 0, payload
```

### Key Decisions (Load-Bearing)
- **Line 6** (classOwnershipAgreement): One-directional, fires only on decision-table + owned conflict; state-machine + zero owned loads cleanly for lint to report as rootless.
- **Line 8** (no-write-block refusal): Conditioned on class; decision-table rows bypass refusal and normalize with empty writes/next/owned.
- **Line 10** (seed predicate): Augments existing `len(Initial) > 0` test; does not replace it.
- **Line 12** (kernel exact-one): Unchanged; decision-table has no owned demand, so every row is a candidate (per C2, C3 compatibility).
- **Line 15** (rowByID first-match): Joined from normalized row set, not kernel; emit is rule-scoped, shared across all expanded rows.
- **Line 16–18** (emit population and rendering): Emit joins payload as a JSON object beside rule; text renderer treats it as a generic leaf path (no per-verb special case per C4).

### Widening Notes
- **Widened past C1**: Read full load ordering in `internal/table/load.go::run` to confirm class agreement step position.
- **Widened past C3**: Read `model.go` to confirm Row structure and understand that emit is not passed to kernel; read `normalize.go` to confirm expand() carries emit through all expanded rows.
- **Widened past C4**: Read `flow_resolve.go` to confirm resolvePayload field order and that text rendering uses generic 0005 renderer (no special case added).
- **Widened past C5**: Read `graphlint` sections to confirm reach(), checkCoverage(), and checkDanglingEdge() integration points.

---

### Footnote: Guesses Marked
**5 GUESS markers** identified in the run:
1. EmitValue struct definition and format (inferred from C3's key-sorted requirement and contrast with TagValue).
2. Assumed Model.Class accessor exists and has zero-value semantics (not explicitly stated in contracts, but required by C1's narrative and standard Go patterns).
3. classOwnershipAgreement function name and position (RDR specifies behavior, not implementation).
4. selectRowWithEmit or equivalent join function (narrative describes joining emit back by rule id; implementation name conjectured).
5. emitUnprovableDimensionOrNone or equivalent zero-dimension handler (C5 specifies emission site and condition; function name conjectured).

These represent silent portions of the RDR where implementation details are not prescribed, only the observable behavior and load order are normative.
