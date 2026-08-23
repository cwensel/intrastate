model: Claude Sonnet 5
variant: full (profile: foundational)

# Reconstruction of the RDR 0002 module from RDR text alone

Scope per the RDR: a Go internal package (unnamed in the text; GUESS name
`internal/transitiontable`, mirroring the sibling `internal/resolve` and
`internal/cli` naming) that parses sparse TOML transition models, normalizes
them into `internal/resolve.Row` values, and renders an expanded-table dump.
The RDR explicitly disclaims package/API naming ("New internal package can own
sparse source structs and normalization" — Existing Infrastructure Audit), so
every symbol name below is a GUESS unless flagged otherwise.

## 1. Public API

The RDR states three operations by role (parse, normalize/load, dump) but
never names exported functions, types, or an error type. All signatures below
are GUESSes shaped strictly by the normative contracts (two-pass version gate,
strict decoding, load-time category list, deterministic dump).

```go
package transitiontable // GUESS package name

// LoadModel is the two-pass entry point the version-gate contract requires:
// pass 1 reads [model] permissively enough to get Version; pass 2, only if
// Version == 1, strict-decodes the full document (unmapped keys rejected).
// GUESS signature; RDR fixes only the two-pass ordering and strict decoding.
func LoadModel(path string) (*Model, error)

// LoadModelBytes is the byte-slice variant, since the RDR discusses parsing
// as a pure function of TOML content, not file I/O specifically. GUESS.
func LoadModelBytes(data []byte, sourceID string) (*Model, error)

// Normalize turns a parsed Model into the deterministic candidate-row set.
// The RDR calls this step "normalize" throughout and fixes its output shape
// (Normative Contracts: "normalized candidate row is the kernel row") but
// never names the function. GUESS signature.
func Normalize(m *Model) ([]resolve.Row, error)

// Dump renders the expanded-table view from the normalized row set. Column
// order is controlled by [dump] settings in the source; the RDR is explicit
// that [dump] is presentation-only and MUST NOT affect Normalize's output.
// GUESS signature and return shape (some serialized/tabular type).
func Dump(rows []resolve.Row, settings DumpSettings) (ExpandedTable, error)

// Model is the typed, decoded (but not yet normalized) source document.
// Field presence driven by the six conceptual schema parts (Technical
// Design �1-6) and the root-key layout normative contract. GUESS field
// names; the TOML key names (outcomes, [model], [tags.<tag>], etc.) are
// RDR-fixed, the Go struct field names are not.
type Model struct {
    Outcomes []string                    // root `outcomes`: recognized alphabet
    Meta     ModelMeta                   // [model]: id, version
    Tags     map[string]TagDecl          // [tags.<tag>]
    Accessors map[string]AccessorRef     // [accessors.<id>]
    Contexts map[string]Context          // [context.<id>]
    Rules    []Rule                      // [[rule]]
    Dump     DumpSettings                // [dump]
}

type ModelMeta struct {
    ID      string
    Version int
}

// TagDecl carries provenance + RDR 0003's type model, cited not owned here.
type TagDecl struct {
    Provenance    Provenance // owned | observed | recognized
    Kind          string     // RDR 0003's value-kind vocabulary; opaque here. GUESS
    Domain        []string   // optional finite domain. GUESS
    Optional      bool       // GUESS
    SetUniverse   []string   // optional set-element universe. GUESS
    SingleValued  bool       // GUESS
    Accessor      string     // optional accessor id for read-back
}

type Provenance int // GUESS enum
const (
    ProvOwned Provenance = iota
    ProvObserved
    ProvRecognized
)

type AccessorRef struct {
    ID string // GUESS shape; RDR only says "accessor references," execution is RDR 0004's
}

type Context struct {
    ID       string
    Inherits string            // context id this inherits from, or ""
    Match    map[string][]Atom // [context.<id>.match.<tag>]
}

// Rule is the sparse authored rule. Exactly one of (Write != nil, Escape
// != nil) is populated per the normative "ordinary XOR escape" contract.
type Rule struct {
    ID      string
    Use     []string            // shared-context references
    Match   map[string][]Atom   // [rule.match.<tag>]
    GuardAll    map[string][]Atom // [rule.guard.all.<tag>]
    GuardUnless map[string][]Atom // [rule.guard.unless.<tag>]
    Write   map[string]string  // [rule.write]; nil for escape rules
    Clear   []string           // rule-level `clear`; nil/empty for escape rules
    Escape  []string           // rule-level `escape` failure classes; nil for ordinary rules
}

type Atom struct {
    Key      string
    Block    Block  // match | all | unless
    Operator string // eq, in, lt, exists, ... (RDR 0003's grammar; opaque here)
    Literal  any    // scalar or ordered sequence for set-valued literals
}

type Block int // GUESS enum
const (
    BlockMatch Block = iota
    BlockAll
    BlockUnless
)

type DumpSettings struct {
    ColumnOrder []string // MAY reorder columns; MUST NOT omit a field
}

type ExpandedTable struct {
    Rows []ExpandedRow // pre-sorted by (model id, rule id, expansion suffix)
}

type ExpandedRow struct {
    Identity        RowIdentity
    SourceLocator   SourceLocator
    Outcome         string
    Atoms           []Atom            // sorted (key, block, operator, literal)
    NextStateTags   map[string]string // includes <clear> entries
    Writes          map[string]string // includes <clear> entries
    RequiresOwned   []string          // sorted, dedup
    EscapeClasses   []string          // sorted; empty for ordinary rows
    Kind            RowKind           // derived: transition | escape
}

type RowIdentity struct {
    ModelID         string
    RuleID          string
    ExpansionSuffix string // "" if the rule did not expand
}

// SourceLocator must identify at least model id + rule id; line/column
// are optional diagnostic detail and excluded from ordering/round-trip
// comparison. GUESS field names.
type SourceLocator struct {
    ModelID string
    RuleID  string
    Line    int // optional, 0 if unknown
    Column  int // optional, 0 if unknown
}

type RowKind int // derived only, never stored as an authored field
const (
    KindTransition RowKind = iota
    KindEscape
)
```

### Error modes

The RDR fixes a closed set of **load-time data-level categories** (Normative
Contracts, "Load-time validation failures MUST retain stable data-level
categories"), which is the load-bearing part of the public error contract.
Category **names** below are quoted/paraphrased from the RDR where given;
where the RDR names a concept but not an exact category string, it is marked
GUESS spelling.

- `malformed_toml`
- `unknown_schema_field`
- `missing_recognized_outcome_alphabet`
- `malformed_recognized_outcome_alphabet` (empty alphabet, duplicate member,
  empty-string member, member containing `#`)
- `unknown_tag`
- `unknown_context`
- `cyclic_context_inheritance`
- `write_to_non_owned_tag`
- `unknown_accessor`
- `unsupported_version`
- `malformed_predicate_atom` (unknown operator; ill-formed literal for its
  operator, including non-canonical existence literals)
- `malformed_escape_declaration`
- `malformed_outcome_binding` (zero or >1 `recognized` atom in match blocks;
  a `recognized` atom authored under a guard block)
- `malformed_rule_id` (contains `#`)
- `duplicate_rule_id`
- `duplicate_model_id`
- `reserved_tag_key` (RDR 0008's category, participates in this RDR's set)

GUESS: exact string spellings (e.g. whether it's `unknown_schema_field` or
`unknown-schema-field`) — the RDR names categories in prose only. GUESS: these
surface as a Go `error` implementing `internal/cli/clierr.CLIError` (the RDR
says load/lint failures "can flow through the existing structured error
gateway," A5), most plausibly via a `LoadError` type carrying `Category`,
`Detail`, and a `SourceLocator`.

Two categories are explicitly **not** this module's: cross-row lint findings
(overlap, gap, dead row, read-before-write, ambiguous expansion) are RDR
0006's, and resolver-time refusals (`unmodeled_outcome`,
`owned_state_unavailable`, `guard_unevaluable`, zero/multiple survivors) are
RDR 0001's kernel, not this module's error surface.

## 2. Three most important internal helper functions

1. **Version-gate pre-decode** (`readVersionOnly` — GUESS name). Responsibility:
   implement the mandatory two-pass load ("read `[model]` permissively enough
   to obtain `version`, refuse on any value but `1`, and only then decode the
   document strictly"). This exists specifically so a future v2 file fails with
   an `unsupported_version` diagnostic instead of an arbitrary
   `unknown_schema_field` on whatever v2-only key the strict decoder hits
   first. Must run before the strict/unmapped-key-rejecting decode.

2. **Outcome-binding lifter** (`liftOutcome` — GUESS name). Responsibility:
   scan only a rule's match blocks (local `match` plus inherited contexts'
   `match`, explicitly *not* `guard.all`/`guard.unless`) for exactly one atom
   on `recognized` using `eq` or `in`; reject rules with zero, two, or an
   out-of-alphabet `recognized` atom; reject a `recognized` atom found under a
   guard block as `malformed_outcome_binding` rather than lifting it. On an
   `in` atom, this function is also what performs 1→N row expansion, minting
   the expansion suffix per outcome member and leaving single-member
   bindings (whether authored `eq` or single-member `in`) unsuffixed.

3. **Context/atom merger** (`mergeContexts` — GUESS name). Responsibility:
   flatten inherited context chains and rule-local blocks into one predicate
   set keyed on the **full atom identity** `(key, block, operator, literal)` —
   never a proper prefix of it — so two atoms sharing key/block/operator but
   differing in literal both survive, merging is idempotent on truly identical
   atoms, and an atom authored in both `all` and `unless` survives as two
   distinct atoms (a dead-rule condition for RDR 0006, not a load failure
   here). This is the function the RDR calls out by name as the one place a
   naive `(block, key, operator)` merge key would silently and undetectably
   drop a constraint.

(A close fourth, not counted in the top three: the **deterministic sorter**
that produces the pre-sorted row/atom/tag/write/escape-class sequences the
dump requires — RDR is explicit the dump "MUST be emitted from a pre-sorted
sequence... never by iterating a map... never by delegating key order to an
encoder.")

## 3. Data model for anything persisted or passed across the boundary

**Persisted (on disk, source of truth, hand-authored):** the sparse TOML
transition model file — root `outcomes` array; `[model]` (`id`, `version`);
`[tags.<tag>]` (provenance, RDR-0003 type fields, optional accessor);
`[accessors.<id>]`; `[context.<id>]` (optional `inherits`, `match.<tag>`
blocks); `[[rule]]` array (`id`, `use` list, `match`/`guard.all`/
`guard.unless` blocks, either `write`+optional `clear` or `escape`);
`[dump]` (column-order/presentation only, no semantic effect).

**Passed across the boundary to the kernel (`internal/resolve`):** the
normalized candidate row, which *is* `resolve.Row` per the RDR ("The
normalized candidate row is the kernel row"). Per-row fields, all fixed by
the Normative Contracts:

| Field | Content |
|---|---|
| Identity | `(model id, rule id, expansion suffix)` — total order key, byte-lexicographic |
| Source locator | model id + rule id required; line/column optional/diagnostic, excluded from ordering and round-trip comparison |
| Outcome | single lifted `recognized` value |
| Match / Guard | unified predicate set split by operator: equality-on-declared-tag atoms → `Match []Tag`; everything else (set membership, comparison, existence, all `unless` atoms regardless of operator) → `Guard` |
| Predicate atoms | `(key, operator token, literal, block)`; existence atoms use kernel's exact `OpExists`/`LiteralTrue`/`LiteralFalse` constants verbatim |
| Next-state tags | rule's write block + clear list, rendered with `<clear>` entries |
| Writes | same rendered set as next-state tags (equal under this RDR's authoring surface, but kept as two distinct fields per RDR 0009 A4) |
| RequiresOwned | derived: sorted, dedup tag keys from write block ∪ clear list; empty for escape rows |
| Escape | failure-class list (`no_match`, `ambiguous_match`); empty for ordinary rows; row kind (`transition`/`escape`) is a *derived view property*, never a stored field |

**Passed to review/tooling (expanded-table dump):** every field above, plus
the derived `Kind` column, rendered as a total-ordered sequence — no dump
grammar/serialization format is specified (explicitly out of scope; "no
delimiter, escaping, quoting, or record separator" — the RDR calls the
rendered text lossy at three named sites: `<clear>` sentinel not reserved in
value space, set/multi-entry fields rendered with unescaped separators,
expansion suffix has no reserved separator). GUESS: concrete output format
(text table vs JSON vs CSV) is left to implementation; RDR fixes only field
list, presence, and row/atom/sub-list ordering, not serialization.

**Round-trip invariant** (value-level, not textual): `parse ∘ normalize`
must be identity over the full field list above (locator's positional
line/column excluded, locator's rule-identifying presence included) for
differently-key-ordered, differently-rule-ordered, or `eq`-vs-single-member-
`in` reauthored variants of one semantic model.

## 4. Top-level pseudo-code of the main operation (load → normalize → dump)

```text
function LoadNormalizeDump(path) -> (ExpandedTable, error):
    raw := readFile(path)

    // Pass 1: version gate, permissive read
    meta := decodePermissive(raw, onlyField="[model]")
    if meta.version != 1:
        return error(unsupported_version)

    // Pass 2: strict decode, now that version is known-good
    model := decodeStrict(raw)          // rejects unmapped keys -> unknown_schema_field
    if err: return err                  // malformed_toml / unknown_schema_field

    // Load-time validation (arity: single-rule, decidable from one rule
    // plus declarations only)
    validateTagDeclarations(model.Tags)          // provenance, recognized-name reserved
    validateAlphabet(model.Outcomes)              // non-empty, dedup, no "", no "#"
    validateAccessorRefs(model.Accessors)
    resolveContextInheritance(model.Contexts)     // unresolvable/cyclic -> error
    for rule in model.Rules:
        validateRuleID(rule.ID)                   // unique, no "#"
        validateWriteTargets(rule, model.Tags)     // must be owned tags
        validateEscapeShape(rule)                  // escape XOR write/clear, never both
        validateOutcomeBinding(rule, model.Contexts) // exactly one recognized atom in match blocks only
        validatePredicateAtoms(rule, model.Tags)   // known ops, canonical existence literals
    if any validation error: return firstOrAggregated(errors)  // stable category per finding

    // Normalize: sparse rules -> explicit candidate rows
    rows := []
    for rule in model.Rules:
        ctxAtoms := mergeContexts(rule.Use, model.Contexts)   // full-identity merge, idempotent
        allAtoms := union(ctxAtoms, rule.Match, rule.GuardAll, rule.GuardUnless)
        outcomeMembers, predicateAtoms := liftOutcome(allAtoms)  // match-blocks-only scan

        for i, outcome in enumerate(outcomeMembers):            // 1 row unless `in` expands
            suffix := "" if len(outcomeMembers) == 1 else outcome
            row := Row{
                Identity: (model.Meta.ID, rule.ID, suffix),
                Locator:  rule.sourceLocator(),
                Outcome:  outcome,
                Match, Guard: splitByOperator(predicateAtoms),  // eq-on-tag -> Match; else -> Guard
            }
            if rule.isEscape():
                row.Escape = rule.EscapeClasses
                row.NextStateTags, row.Writes, row.RequiresOwned = {}, {}, {}
            else:
                writes := renderWrites(rule.Write, rule.Clear)   // includes <clear>
                row.NextStateTags = writes
                row.Writes = writes
                row.RequiresOwned = sortedDedupKeys(rule.Write, rule.Clear)
            rows.append(row)

    // Cross-row structural checks this RDR still owns (duplicate identity)
    checkNoDuplicateRowIdentities(rows)   // duplicate rule id / model id -> error
    // NOTE: overlap/gap/dead-row/read-before-write are RDR 0006 lint, NOT here

    sortedRows := sortByIdentityTuple(rows)   // (model id, rule id, suffix), byte-lexicographic
    for row in sortedRows:
        sortAtomsAndSubLists(row)             // (key, block, op, literal); tags/writes/etc by key

    table := renderDump(sortedRows, model.Dump.ColumnOrder)  // presentation only; ignores nothing semantic
    return table, nil
```

## Where the RDR was silent (summary of GUESS load)

The RDR is a data/normalization *contract* spec, not an API spec: it never
names a package, an exported function, an error type, or a concrete dump
serialization. Every Go identifier, struct field name, and the dump's
wire/text format above is a GUESS. What is *not* guessed — because the RDR
states it normatively — is: the row/field shape (matches `resolve.Row`
1:1), the full load-time error category list, the two-pass version-gate
sequencing, the outcome-lifting rule (match-blocks-only, `#`-free, suffix-
only-on-expansion), the atom-merge-by-full-identity rule, and the total
deterministic ordering key for the dump.
