model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability reconstruction — RDR 0024, run 1

Read via projector only. Selected `C1`–`C4`, `MVV`, `D-identity` /
`D-naming` / `D-selection-predicate`, `S1`–`S9`, plus the illustrative
code block.

**Spans widened past** (contract text left a signature, type, or step
order under-determined):

- `§approach` — C2 fixes the check position by step name but not the
  enforcement-tier rationale or what "load pipeline, not graphlint"
  means for the call site.
- `§phase-1-declaration-grammar-and-load-proof` — C1/C2 name no
  function for line attribution; Phase 1 supplies `emitHeaderLine`,
  the `atLine` stamping rule, and the `conformKind`/`conformDomain`
  reuse path. Without this the error-mode reconstruction would have
  invented locator plumbing.
- `§phase-2-normalized-carry` — C3 names the carrier fields but not
  where the sort is applied; Phase 2 fixes it as load-time union
  construction.
- `§phase-3-envelope-surfacing` — C4 fixes the field's position
  relative to `emit` but not the concrete struct field count or the
  join's source; Phase 3 plus S5 supply `NumField() == 15` and the
  fourteen-key wire order.
- `§illustrative-code` — the TOML/JSON shapes, which no contract
  element carries.

Not widened: `§critical-assumptions` beyond the A-refs the contracts
cite inline, `§trade-offs`, `§decision-rationale`,
`§alternatives-considered`.

---

## 1. Public API

Package `internal/table`.

```go
// Categories

const (
    CatMalformedEmitDeclaration = "malformed_emit_declaration"
    CatUnknownEmitKey           = "unknown_emit_key"
    CatEmitValueOutOfDomain     = "emit_value_out_of_domain"
)

// Categories() gains all three entries (D-naming, S6).
func Categories() []string
```

```go
// Carried declaration (C3).

// EmitDecl is the normalized form of one [emit.<key>] declaration.
// It is deliberately NOT a TagDecl: an emit key is not a tag key.
type EmitDecl struct {
    Kind         string            // "enum" | "bool" | "int" | "scalar"
    Domain       []string          // enum only; bytewise-sorted union; nil otherwise
    Dispositions map[string]string // member -> its disposition token; nil/empty when none
}

// Model gains, beside Model.Tags:
type Model struct {
    // ...
    Tags      map[string]TagDecl
    EmitDecls map[string]EmitDecl // keyed by emit key
    // ...
}
```

GUESS: the exact field ordering inside `Model` and whether `EmitDecls`
is non-nil-but-empty vs nil for a zero-declaration model. The RDR
fixes the payload's `{}`-not-`null` for `dispositions` (C4) but says
nothing about the carrier's nil-ness. I reconstruct it as **always
non-nil, possibly empty**, mirroring `Model.Tags`.

GUESS: `EmitDecl` is exported and its fields are exported (C3 names
`Kind`, `Domain`, `Dispositions` in exported spelling, so this is
near-certain, but the RDR never writes the `type` line).

Source schema (`internal/table/source.go`):

```go
type sourceModel struct {
    // ...
    Emit map[string]sourceEmitDecl `toml:"emit"`
    // ...
}

type sourceEmitDecl struct {
    Kind   string `toml:"kind"`
    Domain any    `toml:"domain"` // two-shaped: []string or map[string][]string
}
```

`Domain any` is forced by C1: it is the one field in the source schema
that cannot be concretely typed, so `DisallowUnknownFields` stops
descending at it — the same carve-out `sourceModel.Metadata`
documents.

GUESS: the struct name `sourceEmitDecl` and the `any` spelling (vs
`interface{}` or `toml.Primitive`). C1 fixes the *behavior* — the
decoder stops descending, both nil-domain authorings decode to
`Domain == nil` — which only `any` (or an equivalent) delivers.

**Error modes.** No new exported error type. All three categories are
raised through the existing load-refusal path and stamped with a
source line via `atLine`:

- `malformed_emit_declaration` — stamped on the `[emit.<key>]` header
  line (C1's arms).
- `unknown_emit_key` — stamped on the offending rule's `[rule.emit]`
  block.
- `emit_value_out_of_domain` — likewise, the offending rule's
  `[rule.emit]` block.

Load is **fail-fast**: exactly one categorized error, never a list,
never `Unwrap() []error` (C2, pinned by
`internal/table/reserved_key_0008_test.go`). Order among independent
defects is deliberately unspecified and is not asserted.

Envelope mapping (not this RDR's to change): all three map to
`flow-model-invalid` under the `flow` verbs (`codeModelInvalid` in
`internal/cli/flow_input.go`) and to `model-invalid` under
`intrastate lint` (`internal/cli/lint.go`'s load-refusal arm). The
category slug travels in the inner finding's `code` on both surfaces.

Payload surface (`internal/cli/flow_resolve.go`):

```go
type resolvePayload struct {
    // ... model, revision, observed, owned, readers, outcome,
    //     rule, gates, emit,
    Dispositions map[string]string `json:"dispositions"` // no omitempty
    // ... next, writes, clear, escaped, escape_class (omitempty)
}
```

15 struct fields; 14 wire keys on a non-escaped plan, in order:
`model, revision, observed, owned, readers, outcome, rule, gates,
emit, dispositions, next, writes, clear, escaped` (S5).

`flow next` carries no `dispositions` (C4).

---

## 2. Three most important internal helpers

**`loadEmitDecls(l *loader) error`** — C1's grammar. Reads
`l.src`'s `[emit]` table, validates each `[emit.<key>]` declaration,
and builds the `map[string]EmitDecl` carrier. Refuses
`malformed_emit_declaration` on: unknown `kind` token; an `enum` with
no usable domain (one arm covering all three indistinguishable
authorings — absent `domain`, empty `[emit.<key>.domain]` sub-table,
`domain = []`); empty-string member; duplicate member, including
duplicates across disposition lists; a `domain` on a non-enum kind;
empty-string disposition token; and the four arms strict decoding
cannot reach because `Domain` is `any`-typed — a `domain` that is
neither array nor table, a non-array disposition value, nesting below
the disposition level, and a non-string member. Builds `Domain` as
the **bytewise-sorted union** and `Dispositions` as member→token.

**`checkRuleEmit(l *loader) error`** — C2's two cross-checks over
every `sourceRule.Emit` (source rules, not normalized rows), on
ordinary and escape rules alike. Refuses `unknown_emit_key` for any
authored key absent from the carrier, and `emit_value_out_of_domain`
for any authored value failing its key's kind/domain. Value
conformance delegates to the existing `ConformValue` with a throwaway
adapter: `ConformValue(TagDecl{Kind: d.Kind, Domain: d.Domain}, value)`
— `Min`/`Max`/`Elements` left nil so the `int` bounds arm and the
`set` arm never fire. This reuse is **safe-by-omission and depends on
`loadEmitDecls` having already refused** unknown kinds and empty enum
domains; if those arms are relaxed, `emit_value_out_of_domain`
silently stops firing.

**`emitHeaderLine(src, key string) int`** — sibling of
`tagHeaderLine`, locating the `[emit.<key>]` declaration header (and
a rule's `[rule.emit]` block) in the raw source the loader holds, so
all three categories can be stamped through `atLine`. Without it every
emit refusal carries no source line, and a fail-fast pipeline that
reports one defect per run without a line number makes the author
search the model by hand each round.

GUESS: the exact signatures. The RDR names all three functions
(`loadEmitDecls`, `checkRuleEmit`, `emitHeaderLine`) but writes no
parameter list or return type for any of them. I reconstruct
`loadEmitDecls`/`checkRuleEmit` as loader-step-shaped (matching
`loadTags`, whose step-slice membership C2 fixes) and `emitHeaderLine`
as `tagHeaderLine`-shaped (`(l.src, key)` per Phase 1's description of
the HEAD call at `load.go:214`).

GUESS: the payload-join helper is a fourth, unnamed function. C4 fixes
the join's semantics precisely (row-keyed, off `Plan.RuleID`, one
entry iff the selected row authors `k` AND `k`'s declaration lists
that value under a disposition) but names no function. I reconstruct
it as inline in `resolvePayload` or a small unexported
`joinDispositions(decls map[string]EmitDecl, emit map[string]string)
map[string]string`.

---

## 3. Data model across the boundary

**Authored (TOML source):**

```toml
[emit.next]
kind = "enum"

[emit.next.domain]          # partitioned form
route = ["/rdr-propose", "/rdr-refine"]
stop  = ["stopped:joint-decision", "stopped:propose-order"]

[emit.dpa]
kind = "enum"
domain = ["required", "waived"]   # flat form: no dispositions
```

`kind` is one of `enum | bool | int | scalar` — RDR 0003's spellings
verbatim, minus `set` (an emit value is one authored string). `bool`
fixes the implicit domain `true | false`; `int` constrains the value
to a base-10 integer literal (`ConformValue` reuse makes `03`, `+5`,
`-0` admissible, inherited from `strconv.Atoi`); `scalar` is
declared-but-unvalidated. None of the three takes a `domain`;
dispositions attach only to declared enum members. Disposition tokens
are **model-authored** — intrastate fixes no vocabulary.

Authored emit values are always TOML **strings** (`sourceRule.Emit` is
`map[string]string`), so under `kind = "int"` an author writes
`count = "42"`; a bare `count = 42` refuses as `malformed_toml` from
the decoder, not as `emit_value_out_of_domain`.

A bare `[emit]` table with zero sub-tables is a **zero-declaration
model**, observationally identical to omitting the table. The opt-in
trigger is the *count of declared keys*, not the table's presence.

**Carried (normalized model):** `Model.EmitDecls map[string]EmitDecl`
as above. Value-preserving, **not** order-preserving: `Domain` is the
bytewise-sorted union, and the partition grouping is not carried — it
is fully recoverable from `Dispositions`. The authored order of a flat
`domain` array is also not preserved. The sort is load-bearing:
without it the union built from map iteration over the domain
sub-table is non-deterministic across loads, and nothing existing
catches that (`TestReq146` is scoped to `table.Row`).

**Wire (`flow resolve` success payload):**

```json
{"rule":"joint-fire","gates":[],
 "emit":{"next":"stopped:joint-decision"},
 "dispositions":{"next":"stop"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

`dispositions` is a JSON object, emit key → disposition token, keys in
byte order, present as `{}` — never `null`, never omitted. No
`omitempty`: the field distinguishes "this model declared nothing"
(`{}`) from "this row's answers carry dispositions" (populated), a
distinction `omitempty` would erase. It is **not** a binary-version
signal.

The map is keyed off the **selected row's** authored emit, never off
the declaration set: a declared key the selected row does not emit
contributes no entry, and the map is never padded with nulls, empty
strings, or absent-markers. A plan rescued by an escape row joins from
that escape row's own authored values.

Text mode renders through the generic payload renderer as
`dispositions.<key>: <token>` / `dispositions: (none)`, with no
per-verb special case.

**Not carried anywhere:** the kernel (`internal/resolve`) sees no
declarations and no emit. `internal/graphlint` reads none of this.
`dump` renders a declared emit value byte-identically to its
undeclared rendering — `renderEmit` keeps bypassing `renderValue`.

---

## 4. Top-level pseudo-code

```
// internal/table/load.go::run — step slice, unchanged except for the
// two inserted steps immediately after loadTags and ahead of
// normalizeRules.
//
//   ... loadTags,
//       loadEmitDecls,        // NEW (C1)
//       checkRuleEmit,        // NEW (C2) — must follow loadEmitDecls
//       loadAccessors, loadDump, loadContexts, loadInitial,
//       loadTerminal,
//       normalizeRules, ...

func loadEmitDecls(l) error:
    decls = {}
    for key, sd in l.src.Emit:                 // decoder already
                                               // refused dup keys
        if sd.Kind not in {enum, bool, int, scalar}:
            return refuse(CatMalformedEmitDeclaration,
                          atLine(emitHeaderLine(l.src, key)))
        if sd.Kind != "enum" and sd.Domain != nil:
            return refuse(CatMalformedEmitDeclaration, ...)

        domain, disp = nil, {}
        if sd.Kind == "enum":
            switch shape of sd.Domain:
            case []any:                        // flat form
                members = each element, must be non-empty string
                                               // else refuse
            case map[string]any:               // partitioned form
                for token, list in it:         // token must be
                                               // non-empty string
                    list must be []any of non-empty strings
                    for m in list:
                        if m already seen: refuse   // one member,
                                                    // one disposition
                        disp[m] = token
                    members += list
            default:                           // nil, or neither
                                               // array nor table
                refuse(CatMalformedEmitDeclaration, ...)
            if len(members) == 0: refuse       // "no usable domain",
                                               // ONE arm
            domain = sortBytewise(members)     // load-bearing (C3)

        decls[key] = EmitDecl{sd.Kind, domain, disp}

    l.model.EmitDecls = decls
    return nil

func checkRuleEmit(l) error:
    if len(l.model.EmitDecls) == 0:            // opt-in leg (C2):
        return nil                             // today's pipeline,
                                                // byte-for-byte
    for r in l.src.Rules (ordinary AND escape):
        for key, value in r.Emit:              // map[string]string
            d, ok = l.model.EmitDecls[key]
            if !ok:
                return refuse(CatUnknownEmitKey,
                              atLine(emitBlockLine(l.src, r)))
            if d.Kind == "scalar": continue     // never refused
            if ConformValue(TagDecl{Kind: d.Kind,
                                    Domain: d.Domain}, value) != nil:
                return refuse(CatEmitValueOutOfDomain,
                              atLine(emitBlockLine(l.src, r)))
    return nil

// internal/cli/flow_resolve.go::resolvePayload — after emit is set:
disp = {}                                       // never nil (C4)
for key, value in selectedRow.Emit:             // row-keyed, not
                                                // declaration-keyed
    if d, ok := model.EmitDecls[key]; ok:
        if token, ok := d.Dispositions[value]; ok:
            disp[key] = token
payload.Dispositions = disp                     // marshals byte-ordered
```

GUESS: `emitBlockLine` as a distinct locator for a rule's
`[rule.emit]` block. Phase 1 says `emitHeaderLine` is stamped "keyed
on the `[emit.<key>]` header for a declaration defect and on the
offending rule's `[rule.emit]` block for `unknown_emit_key` /
`emit_value_out_of_domain`" — one named function serving two keying
strategies. I reconstruct it as either one function taking a
discriminating argument or a second unnamed helper; the RDR does not
say which.

GUESS: `loadEmitDecls` returning early vs collecting. Fail-fast (C2)
makes the first refusal terminal, so I reconstruct `return` on the
first defect, and correspondingly the per-key iteration order is
unspecified — which is exactly why S1 requires one defect per fixture
and S2 forbids asserting *which* defect a two-defect model reports.

GUESS: whether `checkRuleEmit`'s zero-declaration early return is an
explicit guard or falls out of the empty-map loop. C2 fixes the
observable ("byte-for-byte today's; no new refusal is reachable"),
which both spellings satisfy. I show the explicit guard because it
makes the opt-in predicate legible at the call site.
