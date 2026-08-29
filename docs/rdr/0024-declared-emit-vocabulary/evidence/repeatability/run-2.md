model: claude-sonnet-5
variant: full (profile: foundational)

# Reconstruction run 2 — RDR 0024 (Declared emit vocabulary)

Widened past the four normative contracts (C1-C4) into: D-identity /
D-naming / D-selection-predicate (Load-Bearing Decisions), Mini-checks
(disposition/fidelity/trace tables), Illustrative Code, Testing
Strategy S1-S7 (S8/S9 not read in full — inferred from their one-line
titles + C1/C4 by GUESS), MVV, and the Approach section. This was
necessary because C1-C4 leave several signatures (exact function
names beyond `loadEmitDecls`/`checkRuleEmit`, `EmitDecl`'s exact
field tags, `resolvePayload`'s full struct) under-determined — the
Illustrative Code and Testing Strategy sections supply the wire
shapes the contracts only describe in prose.

## 1. Public API

```go
// internal/table package additions

// EmitDecl is the normalized carrier for one declared emit key.
// Not a TagDecl (C3: "deliberately NOT TagDecl").
type EmitDecl struct {
	Kind         string            // "enum" | "bool" | "int" | "scalar"
	Domain       []string          // bytewise-sorted union of all members; nil/empty for non-enum kinds
	Dispositions map[string]string // member -> disposition token; only enum members with a partitioned domain appear
}

// Model gains a new field alongside the existing Tags map (C3).
type Model struct {
	// ... existing fields ...
	EmitDecls map[string]EmitDecl // keyed by emit key; empty (not nil, by GUESS matching Tags convention) map when [emit] is absent or empty
}

// New load-pipeline steps, inserted immediately after loadTags in
// load.go::run's step slice, in this order (C2):
func loadEmitDecls(src *sourceDoc, m *Model) error // C1's grammar checks; builds m.EmitDecls
func checkRuleEmit(src *sourceDoc, m *Model) error // C2's two cross-checks; reads m.EmitDecls + sourceRule.Emit

// New load-refusal categories (internal/table/category.go), house
// scheme malformed_*/unknown_* (D-naming):
const (
	CatMalformedEmitDeclaration = "malformed_emit_declaration"
	CatUnknownEmitKey           = "unknown_emit_key"
	CatEmitValueOutOfDomain     = "emit_value_out_of_domain"
)
// Each registered in table.Categories() (S6) with a testdata/neg/*.toml witness.

// internal/cli/flow_resolve.go: resolvePayload gains a field,
// inserted immediately after Emit (C4). By GUESS (struct tags not
// quoted verbatim in the record, inferred from sibling Emit field
// and 0010:C4 convention):
type resolvePayload struct {
	// ... Model, Revision, Observed, Owned, Readers, Outcome, Rule, Gates ...
	Emit         map[string]string `json:"emit"`
	Dispositions map[string]string `json:"dispositions"` // never omitempty; always {} at minimum (C4)
	// ... Next, Writes, Clear, Escaped ...
	EscapeClass string `json:"escape_class,omitempty"` // pre-existing, unchanged
}
```

Error modes (all load-time refusals, fail-fast: exactly one
categorized error per run, never a list — inherited from
`internal/table/load.go`'s `Load` contract, C2):

- `malformed_emit_declaration` — C1's arms: unknown `kind` token;
  enum with no usable domain (no `domain` key, empty domain
  sub-table, or `domain = []` — one arm, per C1's "not two" ruling);
  empty-string enum member; duplicate member (including across
  disposition lists); `domain` on a non-enum kind; empty-string
  disposition token; `domain` neither flat array nor disposition
  table; a disposition value that isn't an array of strings; nesting
  below the disposition level; a non-string member.
- `unknown_emit_key` — a `[rule.emit]` key (ordinary or escape rule)
  not declared under `[emit]`, when one or more declarations exist.
- `emit_value_out_of_domain` — an authored value not in its key's
  enum domain, not a `bool` token, or not an `int` literal (`scalar`
  never refused).
- Pre-existing, unchanged by this RDR: `malformed_toml` (non-string
  emit value, e.g. `verdict = 42`) and `unknown_schema_field` (a
  declaration-level key typo, e.g. `domaim = [...]`) — both fire from
  the existing decoder, not from the two new steps.

All three new categories map to `flow-model-invalid` (flow verbs,
`0005:C1`) or `model-invalid` (`intrastate lint`'s load-refusal arm) —
the outer envelope code differs by surface; the category slug travels
in the inner finding's `code` on both (C2).

## 2. Three most important internal helpers

1. **`loadEmitDecls(src *sourceDoc, m *Model) error`** — reads the
   top-level `[emit]` table from the source document, validates each
   `[emit.<key>]` sub-table against C1's grammar (kind token, domain
   shape, member/disposition well-formedness), and on success builds
   `m.EmitDecls`: sorts each enum's member union bytewise, and
   populates `Dispositions` from whichever domain form was authored
   (flat array = no dispositions; sub-table = union of its keyed
   arrays). Must run and fully succeed before `checkRuleEmit`, since
   the cross-check reads the carrier this step builds (C2). Refuses
   `malformed_emit_declaration` on any grammar violation, before any
   row is yielded.

2. **`checkRuleEmit(src *sourceDoc, m *Model) error`** — iterates
   every rule's `sourceRule.Emit` (ordinary and escape rules alike)
   against `m.EmitDecls` and performs C2's two cross-checks: (a) every
   authored emit key must be declared (else `unknown_emit_key`); (b)
   every authored value must conform to its key's declared kind/domain
   via a throwaway `TagDecl{Kind: decl.Kind, Domain: decl.Domain}`
   passed to the existing `ConformValue`/`conformKind`/`conformDomain`
   machinery (else `emit_value_out_of_domain`). Runs only when
   `m.EmitDecls` is non-empty (zero declarations = C2's opt-in leg,
   byte-for-byte today's pipeline). Fail-fast: returns on the first
   defect found, in unspecified order relative to other load checks.

3. **`resolvePayload`'s dispositions join** (a join performed inline
   inside `internal/cli/flow_resolve.go::resolvePayload`'s
   construction, not necessarily its own named function — GUESS on
   whether it is factored into a helper; the record only fixes its
   *position and behavior*, not whether it is extracted) — for the
   selected row's (or, on an escape rescue, the escape row's own)
   authored `Emit` map, looks up each key in `m.EmitDecls`, and where
   the declaration lists the authored value under a disposition,
   inserts `key -> disposition` into the `dispositions` output map.
   Keys the row doesn't emit, or whose value carries no disposition
   (non-enum kind, flat-array domain, `scalar`), contribute no entry.
   Result defaults to `{}`, never `nil`/`null`, never omitted.

## 3. Data model (persisted / boundary-crossing)

**Source (authored TOML, on-disk model file):**
```toml
[emit.<key>]
kind = "enum" | "bool" | "int" | "scalar"
# enum only, exactly one of:
domain = ["member1", "member2", ...]              # flat: no dispositions
[emit.<key>.domain]                                # partitioned: keys are author-chosen disposition tokens
route = ["member1", "member2"]
stop  = ["member3"]
```
Decoded via `pelletier/go-toml/v2` with `DisallowUnknownFields`; the
`domain` field is the one place strict decoding cannot descend
further (it is `any`-typed), so C1's hand-written arms cover what the
decoder can't.

**Normalized, in-process (`Model.EmitDecls map[string]EmitDecl`):**
- Key: emit key (byte-exact, one declaration per key).
- `EmitDecl.Kind string` — the authored kind token, verbatim.
- `EmitDecl.Domain []string` — bytewise-sorted union of all declared
  members (order not preserved from authoring; recoverable from
  `Dispositions` if partition grouping is needed).
- `EmitDecl.Dispositions map[string]string` — member -> disposition
  token; empty/nil for non-enum kinds and flat-array-domain enums.
- Carried unmutated through normalization (`0010:C3`'s "Emit is never
  mutated" extended to declarations by C3). Not read by
  `internal/resolve` (kernel) or `internal/graphlint` — GUESS confirmed
  by explicit text: "the kernel continues to carry no declarations."

**Wire (JSON, `flow resolve` success payload):**
```json
{
  "model": "...", "revision": "...", "observed": "...", "owned": "...",
  "readers": [...], "outcome": "...", "rule": "...", "gates": [...],
  "emit": {"next": "stopped:joint-decision"},
  "dispositions": {"next": "stop"},
  "next": {}, "writes": {}, "clear": [], "escaped": false
}
```
14 fields on the wire (up from 13 pre-change), `dispositions`
inserted at index 9 immediately after `emit`. `reflect.TypeOf(
resolvePayload{}).NumField() == 15` (the 15th, `escape_class`, is
`omitempty` and normally absent). `flow next` payloads carry no
`dispositions` field at all (not even empty) — GUESS-by-explicit-text:
the record states this directly rather than leaving it silent.

## 4. Top-level pseudo-code (main operation: model load with emit checks)

```
func Load(sourcePath) (*Model, error):
    src := decodeStrict(sourcePath)          // existing; refuses
                                              // unknown_schema_field,
                                              // malformed_toml, etc.
    m := &Model{}

    loadTags(src, m)                         // existing step
    // --- new steps, inserted here (C2) ---
    err := loadEmitDecls(src, m)             // C1: parse+validate [emit]
    if err != nil:
        return nil, err                      // malformed_emit_declaration
                                              // fail-fast, before rows yielded

    if len(m.EmitDecls) > 0:                 // C2 opt-in gate
        err = checkRuleEmit(src, m)          // C2: cross-check every
        if err != nil:                       // sourceRule.Emit
            return nil, err                  // unknown_emit_key or
                                              // emit_value_out_of_domain
    // zero declarations => byte-for-byte today's pipeline, no-op
    // --- end new steps ---

    loadAccessors(src, m)                    // existing, independent
    loadDump(src, m)
    loadContexts(src, m)
    loadInitial(src, m)
    loadTerminal(src, m)

    normalizeRules(src, m)                   // mints candidate rows;
                                              // emit block carried
                                              // unmutated (0010:C3)
    return m, nil

func resolvePayload(plan *Plan, m *Model) Payload:
    row := plan.SelectedRow                  // or escape row on rescue
    p := Payload{
        // ... Model, Revision, Observed, Owned, Readers, Outcome,
        //     Rule, Gates ...
        Emit: row.Emit,                      // authored bytes, verbatim
    }

    p.Dispositions = map[string]string{}     // never nil (C4)
    for key, value := range row.Emit:
        decl, declared := m.EmitDecls[key]
        if not declared:
            continue                          // no entry; loader already
                                               // proved this can't be
                                               // unknown_emit_key if
                                               // declarations existed
        if disposition, ok := decl.Dispositions[value]; ok:
            p.Dispositions[key] = disposition
        // else: non-enum kind / flat domain / scalar -> no entry

    p.Next = ...; p.Writes = ...; p.Clear = ...; p.Escaped = ...
    return p                                  // dispositions sits
                                               // immediately after emit
                                               // in struct/JSON order
```

## Notes on guesses

- The exact extraction boundary of the "dispositions join" (item 2.3)
  as its own named function vs. inline code is a GUESS — the record
  fixes behavior and position, not internal factoring.
- `resolvePayload`'s exact Go struct tags/field names beyond
  `Emit`/`Dispositions` are GUESSed from convention (existing sibling
  fields named in Testing Strategy scenario 5's 14-key wire list).
- Whether `m.EmitDecls` is nil or an empty (non-nil) map when `[emit]`
  is absent is a GUESS by analogy to `Model.Tags`; the record does not
  state this explicitly.
- S8 and S9 (Phase 4 shipped artifacts; `dump` rendering a declared
  emit value raw) were not read in full text — their titles plus C1/C4
  make their content inferable, but any signature-level detail from
  them is not reflected above and is out of scope for this run's
  guesses.
