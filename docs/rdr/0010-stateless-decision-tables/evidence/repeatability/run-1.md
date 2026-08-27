model: claude-sonnet-5
variant: full (profile: foundational)

# Reconstruction of RDR 0010: Stateless Decision Tables

Widened past the contract spans into: `§technical-design` lead-in (the
three-seam summary and data-flow line), `§illustrative-code` (the worked
TOML/mermaid example), `§load-bearing-decisions`, and `§authority` (the
class-decision ownership table). These are non-normative but were needed to
fix concrete Go symbol names, field types, and call order that the contract
prose refers to but does not fully spell out in one place (e.g. the exact
`sourceModel`/`Row`/`table.Model` field shapes, and which four
`len(Initial)==0` sites exist and which two are augmented).

---

## 1. Public API

### 1.1 Model declaration (TOML surface)

```toml
[model]
id = "navigator"
version = 1
class = "decision-table"        # NEW optional key; absent => "state-machine" (C1)

[tags.status]
provenance = "observed"         # decision-table models declare ZERO "owned" tags (C1)
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[[rule]]
id = "draft-propose"
[rule.match.recognized]         # match atoms scope the group only (0006:C7)
eq = "locate"
[rule.guard.all.status]         # guard atoms are the product dimensions (C5)
eq = "Draft"
[rule.emit]                     # NEW: flat string->string answer block (C3)
next = "propose"
```

Grammar-level contract (C1, C2, C3):
- `[model].class : string`, admitted values `{"state-machine", "decision-table"}` (GUESS on exact TOML type beyond "string" — RDR only says "a string"); absent reads as `"state-machine"`. Any other value: load failure, category `malformed model declaration`.
- A `"decision-table"` model MUST declare zero `provenance = "owned"` tags. Non-zero owned count under that class: load failure `malformed model declaration`, detail containing literal token `owned=<n>`.
- A `"decision-table"` model MUST NOT declare `[initial]`, `terminal`, or an accessor whose `keys` name an owned tag; its ordinary rules MUST NOT carry a write block or clear list. These surface as pre-existing refusal categories (`unknown tag`, `malformed accessor binding`, `malformed rule shape`) rather than new ones — condition is only on class.
- Any rule (ordinary or escape), either class, MAY carry `[rule.emit]`: flat TOML table, values MUST be strings (nested table or non-string value refuses as `malformed TOML`, decoder type-error category). Duplicate keys are unauthorable (TOML decoder refuses first).

### 1.2 Go types (internal/table)

```go
// internal/table/source.go (sourceModel, sourceRule — authored/decoded shape)
type sourceModel struct {
    // ... existing fields ...
    Class string `toml:"class"` // GUESS: exact field name; RDR fixes only that
                                  // sourceModel "gains class" and default is
                                  // "state-machine" when absent
}

type sourceRule struct {
    // ... existing fields ...
    Emit map[string]string `toml:"emit"` // C3: "sourceRule.Emit is typed map[string]string"
}
```

```go
// internal/table/model.go (normalized row + model, post-load)
type EmitValue struct { // C3: "a NEW row-level type {Key, Value string} — not TagValue"
    Key   string
    Value string
}

type Row struct {
    // ... existing fields: RequiresOwned, Writes, NextTags, ... ...
    Emit []EmitValue // C3: sorted by key; empty (not nil-vs-absent-distinguished) when unauthored
}

// table.Model gains a class accessor (Technical Design §1):
func (m *Model) Class() string // zero value / default "state-machine" (GUESS on exact
                                 // method name and receiver; RDR specifies only "an
                                 // accessor whose zero value is state-machine" on
                                 // table.Model — the identifier Class() is inferred)
```

### 1.3 CLI surface

- `intrastate flow resolve --model <path> --outcome <o> --tag k=v... [--artifact ...] [--as=json|text]`
  - **Unchanged signature.** `--outcome` still required in both classes (C4). `--artifact` for a role no invoked reader needs is ignored (this already held pre-RDR per 0005:C1; restated, not new).
  - Success JSON payload (`resolvePayload`, `internal/cli/flow_resolve.go`) gains one field, `emit`, positioned **immediately after `rule`** in the (now) fourteen-field declaration order (C4 says the existing payload was thirteen fields; `emit` makes fourteen — GUESS on the exact prior count's other field names, RDR only asserts position relative to `rule`):
    ```json
    {
      "type": "ok",
      "data": {
        "rule": "draft-propose",
        "emit": {"next": "propose"},
        "next": {}, "writes": {}, "clear": [], "owned": {}, "readers": [],
        "...": "... remaining 0005:C1 fields, empty for decision-table ..."
      }
    }
    ```
  - `emit` is **always present**, `{}` (never `null`, never omitted) when the selected row authored none.
  - Text mode: no per-verb special case. Renders via 0005's generic renderer as `emit.<key>: <value>` per pair, and `emit: (none)` for an empty block.
- `intrastate flow next --model <path> ...` — **unchanged** (this RDR does not touch it; A11). Over a decision-table model every candidate has empty `required`/`next`/`writes`. `flow next` carries no `emit` field (rejected as BR4).
- `intrastate flow read-state`, `intrastate flow set-state` — unchanged.
- `intrastate lint --model <path> --as=json` — unchanged signature; over a decision-table model, reachability is rooted at the empty owned-state node (C5) instead of `[initial]`.
- `intrastate dump --model <path> [--dump col,col,...]` — gains one column, `emit`, appended **last** (after `escape`) in the default column set for both classes; rendered as `key=value` pairs in key order, bracketed like `writes` (e.g. `[next=propose]`), empty as `[]`. A `[dump]` list that omits `emit` is unaffected unless the model/fixture explicitly enumerates columns (S3: omitting it from an explicit `[dump]` refuses `malformed dump declaration` — GUESS-adjacent: this refusal is stated as an S3 test expectation, treated here as normative behavior since it is asserted, not merely illustrative).

### 1.4 Error modes (new or newly-reachable)

| Condition | Category | Notes |
|---|---|---|
| `class` outside `{state-machine, decision-table}` | `malformed model declaration` | new arm |
| `decision-table` class with non-zero owned tag count | `malformed model declaration` | detail: `owned=<n>` literal token |
| `decision-table` declares `[initial]` / `terminal` / accessor on owned tag not possible (owned set is empty) — falls through to pre-existing | `unknown tag` or `malformed accessor binding` | no new category |
| `decision-table` ordinary rule with write/clear block | `malformed rule shape` (state-machine only) — decision-table loads with empty `Writes`/`NextTags`/`RequiresOwned` instead of refusing | class-conditioned arm of existing category |
| `[rule.emit]` value not a string, or nested table | `malformed TOML` (decoder type-error category) | asserted on category only, not message text |
| `[dump]` explicitly listed and omits `emit` | `malformed dump declaration` | GUESS: inferred from S3's expected-refusal language |
| Graph lint: decision-table group with zero participating guard dimensions | finding `graph-unprovable-coverage`, `reason = "no-participating-dimension"` (new 4th reason value, append to 0006's closed set — contingent on A14) | not a load error; a lint finding |

---

## 2. Three most important internal helper functions

1. **`internal/table/load.go::run` (loader orchestrator / fail-fast pipeline)**
   Responsibility: fixes the ordered sequence `loadModelHeader, loadOutcomes, loadTags, ...` and runs each step fail-fast. It is the reason the class/owned-set agreement check cannot live in the `[model]` header step: `loadModelHeader` executes before `l.model.Tags` is populated, so the check must be a step at or after `loadTags`. This ordering is itself normative (C1) — an undeclared-tag refusal must precede a class-disagreement refusal.

2. **`internal/table/normalize.go::normalizeRule` (per-rule normalization / shape validation)**
   Responsibility: validates rule shape and conditions the no-write-block arm of `malformed rule shape` on the model's class — refusing for `state-machine`, admitting (with empty `Writes`/`NextTags`/`RequiresOwned`) for `decision-table`. It is also where `Row.Kind()` (ordinary vs. escape) is presence-inferred per the Authority table, contrasted explicitly with the class decision, which is never re-derived downstream. `expand` (same file) is the sibling responsibility that mints one row per match-block member from an eight-field seed literal (not by copying `base`), which is why `Emit` must be threaded explicitly through `expand` — an omission would silently drop the answer for every non-first expanded row of a multi-member `in` match atom (S3's mandatory expanding-rule case).

3. **`internal/graphlint/coverage.go::checkCoverage` (coverage/overlap analysis, decision-table zero-dimension arm)**
   Responsibility: computes the guard-dimension scoped product and emits coverage/overlap findings. Gains a class-keyed arm inside its existing `len(dims) == 0` branch: for a `decision-table` group with zero participating guard dimensions, it must emit `graph-unprovable-coverage` (reason `no-participating-dimension`) and return **before** falling into `emitCoverageArms` (whose `bareEscapeFor` path would otherwise close the group as covered via an escape row). This ordering is load-bearing because the empty scoped product carries a single vacuous assignment that any one ordinary row's coverage trivially satisfies — the defect this RDR exists to prevent.

   (Runner-up, named because it is the other augmented site: `internal/graphlint/reach.go::reach`, whose root-seeding guard becomes `class == "decision-table" || len(Initial) > 0` — augmenting, not replacing, the existing test, so a rootless state-machine still traverses nothing and still takes `checkDanglingEdge`'s missing-root finding.)

---

## 3. Data model (persisted / boundary-crossing)

### 3.1 Persisted (authored TOML → loaded model)

```
sourceModel
 ├─ ...existing fields...
 └─ Class string                 // "state-machine" | "decision-table" | absent(->state-machine)

sourceRule
 ├─ ...existing fields...
 └─ Emit map[string]string       // decoded; refuses on non-string / nested value

table.Model (post-load, in-memory)
 ├─ ...existing fields...
 ├─ Tags map[string]Tag          // provenance incl. "owned" | "observed" | "recognized"
 └─ Class() string                // accessor; zero value "state-machine" — SINGLE
                                   //   source of truth for "what class is this";
                                   //   never re-derived from len(owned)==0 elsewhere

Row (normalized, per expanded rule)
 ├─ RequiresOwned []...           // empty for decision-table rows
 ├─ Writes []...                  // empty for decision-table rows
 ├─ NextTags []...                // empty for decision-table rows
 ├─ Emit []EmitValue              // NEW; key-sorted; empty sequence if absent
 └─ Suffix ...                    // distinguishes expanded rows sharing one rule id

EmitValue { Key string; Value string }   // NEW type; NOT TagValue (no member-sequence,
                                          // no clone/render/clear-sentinel machinery)
```

Byte/wire shapes (Load-Bearing Decisions):
- Dump column `emit`: `key=value` pairs, key-sorted, bracketed like `writes` (e.g. `[next=propose]`), empty as `[]`; appended last after `escape` in `dumpColumns`.
- JSON wire (`resolvePayload`): `emit` is a JSON object, string→string, keys in byte order, HTML-escaping disabled (consistent with other 0005:C1 payload maps); present as `{}` never `null`/omitted.
- `emit` is explicitly **not** carried by `internal/graphlint/engine.go::Fingerprint` (A7) — it does not affect finding identity, and is **not** added to the kernel row (`internal/resolve/resolve.go::Row` stays untouched); the CLI joins `Emit` back onto the kernel's selected `Plan.RuleID` after selection, safe because emit is authored per-rule and shared by all of a rule's expanded rows.

### 3.2 Not persisted / not crossing the boundary

- `emit` keys are **not tag keys**: undeclared, uninterpreted, compared by exact byte equality, never matched/guarded/written/read by any accessor.
- The model's `class` is never inferred from the owned-tag set at any read site (four call sites: `normalizeRule`, `reach`, `checkDanglingEdge`'s root arm, `checkCoverage`'s zero-dim arm) — always read from the one accessor.

---

## 4. Top-level pseudo-code: `flow resolve` over a decision-table model

```
function FlowResolve(modelPath, outcome, tags map[string]string, artifact *string) Payload:
    model := table.Load(modelPath)
    //   load.go::run, fixed order: loadModelHeader, loadOutcomes, loadTags, ...
    //   loadModelHeader reads sourceModel.Class (raw string, may be absent)
    //   at/after loadTags: agreement check —
    //     if model.Class() == "decision-table" and len(ownedTags(model)) != 0:
    //         fail("malformed model declaration", detail="owned=" + count)
    //     if model.Class() not in {"state-machine","decision-table"}:
    //         fail("malformed model declaration")
    //   normalizeRule per rule:
    //     if model.Class() == "state-machine" and rule has no write block:
    //         fail("malformed rule shape")
    //     else if decision-table:
    //         row.RequiresOwned, row.Writes, row.NextTags = [], [], []  // empty, not refused
    //     row.Emit = sortByKey(decodeEmitBlock(rule))  // [] if absent or present-but-empty
    //   expand() mints one Row per match-block member, carrying Emit on every
    //     expanded row (not just the first)

    demand := computeReaderDemand(model)          // empty set: decision-table rules
                                                    // declare no RequiresOwned
    if artifact != nil and role(artifact) not in demand:
        // ignored per 0005:C1 — no invocation, no refusal

    view := assembleView(tags, invokeReaders(demand))  // no readers invoked when
                                                         // demand is empty
    plan := kernel.Resolve(view, outcome)
        // exact-one hit policy over kernel rows (owned view is empty for
        // decision-table); ambiguous match -> flow-ambiguous-match in
        // either class; no match -> escape "otherwise" row rescue, or
        // no_match

    selectedRow := rowByID(model, plan.RuleID)     // first-match join; safe because
                                                     // Emit is identical across all
                                                     // expanded rows of one rule
    payload := resolvePayload{
        Rule:    plan.RuleID,
        Emit:    toJSONObject(selectedRow.Emit),    // {} if empty, never null
        Next:    {}, Writes: {}, Clear: [], Owned: {}, Readers: [],
        // ...remaining 0005:C1 fields, empty for decision-table...
    }
    return ok(payload)
```

```
function LintDecisionTableModel(model) []Finding:
    // reach.go::reach — root-seeding guard AUGMENTS len(Initial)==0, doesn't replace it:
    root := (model.Class() == "decision-table") ? EMPTY_OWNED_STATE_NODE
                                                  : (len(model.Initial) > 0 ? model.Initial : NONE)
    if root == NONE:
        emit checkDanglingEdge missing-root finding (element=model)  // 0006:C18, unchanged path
        return  // rootless model traverses nothing

    reachable := traverseFrom(root)   // for decision-table: exactly {root}, every
                                        // selection context reachable

    for group in guardScopedGroups(model):     // dims from guard.all/guard.unless ONLY;
                                                 // match atoms never contribute (0006:C7)
        dims := guard.Dimensions(group)
        if model.Class() == "decision-table" and len(dims) == 0:
            emit graph-unprovable-coverage(group, reason="no-participating-dimension")
            continue   // MUST precede emitCoverageArms; a bareEscapeFor row must
                        // not close this group as covered
        runOverlapAndCoverage(group, dims)      // unchanged invariant 3/4 machinery
        // graph-dead-end, always-present-owned, owned-before-write,
        // single-valued-state: vacuous by construction for decision-table, MUST NOT emit
        // dangling-edge terminal-predicate arm: stays live (enforces C2's ban on
        // `terminal` over an observed tag)
```

---

GUESS markers used: 6 (sourceModel.Class field name/tag; table.Model.Class() method identifier; the prior 13-field payload's other field names/order beyond position relative to `rule`; the `malformed dump declaration` refusal on an explicit `[dump]` list omitting `emit` — inferred from S3's test-expectation language rather than a direct contract clause; exact TOML type of `class` beyond "string"; and the `resolvePayload` field count arithmetic (13 -> 14) is an inference, not a directly stated total).
