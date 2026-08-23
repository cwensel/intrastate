model: claude-fable-5
variant: full (profile: foundational)

# Reconstruction of RDR 0002 — Transition Table As Reviewable Data

Everything below is derived only from the RDR text. Items the RDR does not fix
are marked GUESS.

## 1. Public API

The RDR mandates a new internal package owning "sparse source structs and
normalization" (Existing Infrastructure Audit), targeting the shipped kernel row
`internal/resolve::Row`. Package name is not fixed.

GUESS: package `internal/transition` (the RDR names the artifact "transition
model"; alternatives like `internal/table` or `internal/model` fit equally).

```go
package transition

// Load reads a transition model from sparse TOML bytes.
// Two-pass, per the version-gate contract: pass 1 reads [model] permissively
// enough to obtain `version` and refuses any value but 1 with
// CategoryUnsupportedVersion; only then does pass 2 decode strictly,
// rejecting unmapped keys (CategoryUnknownSchemaField).
// All single-rule validation (the load-time category set) runs here;
// cross-row checks (overlap, gap, dead row, read-before-write) do NOT —
// those are RDR 0006 lint.
func Load(src []byte) (*Model, error)
    // GUESS: []byte rather than a path — config discovery / path binding is
    // explicitly deferred to CLI integration (RDR 0005), so the format layer
    // should not touch the filesystem. A LoadFile convenience wrapper is
    // plausible but not required by the RDR.

// Normalize expands the sparse model into the deterministic candidate-row
// set: contexts flattened, all/unless merged with per-atom block retained,
// the single `recognized` match atom lifted into the outcome field, `in`
// atoms expanded one row per member (suffix only when >1 row results),
// RequiresOwned derived from writes ∪ clears, next-state tags and writes
// populated independently (never aliased), existence atoms carrying the
// kernel's OpExists / LiteralTrue / LiteralFalse verbatim.
// [dump] settings are ignored entirely — they must not reach the value.
func Normalize(m *Model) (*Expanded, error)
    // GUESS: Normalize can still return an error (some checks, e.g. outcome
    // binding, are naturally discovered during expansion). The RDR only
    // requires that every listed category is refused "at load"; whether
    // Load internally calls Normalize to surface them, or the two stages
    // each own some categories behind one facade, is unspecified.
    // GUESS: alternatively Load returns (*Model, *Expanded, error) in one
    // call; the RDR treats load+normalize as one pipeline.

// Rows converts the normalized candidate rows into kernel rows for
// internal/resolve.Resolve. No new kernel surface: the target shape is
// resolve.Row {RuleID, SourceLocator, Outcome, Match, Guard(→atom-shaped
// after the RDR 0007 reshape), RequiresOwned, Writes, Escape}.
// Equality atoms on declared tags route to Match; every other atom —
// set membership, comparison, existence, and every unless atom — routes
// to the guard field; no atom appears in both.
func (e *Expanded) Rows() []resolve.Row
    // GUESS: method name and exact signature; the RDR fixes only that the
    // normalizer "produces kernel rows" and the routing rule.

// Dump renders the expanded table view: every field of the row value
// (identity, locator, kind column derived from the escape-class list,
// outcome, atoms with block, next tags, writes incl. <clear>,
// requires_owned, escape classes), rows in total identity-tuple order,
// emitted from a pre-sorted slice, never from map iteration.
// [dump] settings may reorder columns only — never omit a field.
func Dump(e *Expanded, w io.Writer) error
    // GUESS: io.Writer and the concrete text grammar — the RDR explicitly
    // defines NO dump grammar (no delimiter/escaping/record separator) and
    // claims no textual round-trip, so the rendering format is free.
```

### Error mode

```go
// LoadError is the stable data-level failure, categorized before any CLI
// mapping. The category set is normative (see §3); CLI code names, Group
// assignment, and exit mapping are explicitly NEW design work at
// implementation (A5: these would be the first parse-validation codes in
// the clierr envelope) and are not part of this package's contract.
type LoadError struct {
    Category Category  // stable, asserted by tests by category not message
    Model    string    // model id when known
    Rule     string    // offending rule id when applicable
    Detail   string    // human diagnostic
    Locator  Locator   // routes back to the authored rule
}
func (e *LoadError) Error() string
// GUESS: struct field set beyond Category; the RDR fixes only that
// categories are stable and distinguishable and that diagnostics carry
// enough source identity to reach the authored sparse rule.
```

Failure taxonomy at the boundary (fixed by the disposition table):
- load failure (this package): the category list in §3 — loud, stable category.
- lint finding (RDR 0006, not this package): overlap, gap, dead row,
  read-before-write, ambiguous expansion; expansion counts reported
  non-blocking.
- runtime refusal (kernel, RDR 0001): unmodeled_outcome,
  owned_state_unavailable, guard_unevaluable, zero/multiple survivors —
  never this package's to emit.

## 2. Three most important internal helpers

1. `flattenContexts(rule, contexts) (atomSet, error)` — resolves the rule's
   `use` references transitively, refusing unknown-context and
   cyclic-context-inheritance; merges inherited match atoms with the rule's
   local match, guard.all, and guard.unless atoms into one predicate set keyed
   on the FULL atom identity `(key, block, operator token, literal)` — never a
   prefix of it, so two atoms differing only in literal both survive; merge is
   idempotent on identical atoms (accumulate-only, no override); contradictory
   atoms are kept (dead rule is RDR 0006's finding, not pruned here).

2. `liftOutcome(rule, matchAtoms, alphabet) (outcomes []string, rest atomSet, error)` —
   outcome binding over the NARROW extent (local match + inherited context
   match blocks only, never guard blocks): finds exactly one `recognized` atom
   using `eq` or `in` with literal(s) inside the alphabet, removes it from the
   predicate set, and returns the outcome list that drives expansion. Refuses
   malformed-outcome-binding on zero atoms, more than one atom, an
   out-of-alphabet literal, or a `recognized` atom authored under guard.all or
   guard.unless (never lifted from a guard — lifting from unless would invert
   the author's intent). Expansion suffix attaches only when the rule yields
   more than one row (single-member `in` == `eq`, no suffix).

3. `sortExpanded(rows)` (with its comparator) — establishes the format-contract
   total order: rows by the identity tuple `(model id, rule id, expansion
   suffix)` compared field-by-field, byte-lexicographically, absent suffix
   before any present one, source locator excluded; within a row, atoms by
   `(key, block, operator token, literal)`; next tags, writes, requires_owned,
   escape classes by key; set-literal members byte-sorted. No positional
   tiebreak exists to fall back on.

(Runner-up, worth naming: `deriveRequiresOwned(rule)` — sorted duplicate-free
union of write-block and clear-list keys, empty on escape rows; the normalizer
is the ONLY enforcement point for the write-free escape row.)

GUESS: helper names and exact signatures; the responsibilities themselves are
normative.

## 3. Data model

### On-disk source (sparse TOML, the persisted artifact) — normative layout

- Root `outcomes = [...]` — closed recognized-outcome alphabet: non-empty,
  duplicate-free, no empty string, no member containing `#`.
- `[model]` — required `id` and `version`; only `version = 1` accepted.
- `[tags.<tag>]` — declaration per tag: `provenance` ∈ {owned, observed,
  recognized}; RDR 0003's type model fields (value kind, optional finite
  domain, optionality, set-element universe, single-valued marker — carried
  through without loss, semantics not owned here); optional accessor
  reference. The recognized-provenance declaration must be named `recognized`
  and nothing else may take that name (reserved_tag_key, RDR 0008).
- `[accessors.<id>]` — named accessor references (bindings only; execution is
  RDR 0004's). GUESS: contents of an accessor entry (the RDR never shows one).
- `[context.<id>]` — shared match contexts; optional `inherits`; predicates
  under `[context.<id>.match.<tag>]`.
- `[[rule]]` — `id` (unique per model, byte-compared, no `#`), optional
  `use = [context ids]`, `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`,
  `[rule.guard.unless.<tag>]`, and either `[rule.write]` (one or more tag
  assignments) plus optional rule-level `clear` list, or a rule-level
  `escape` list (∈ {no_match, ambiguous_match}) with NO write block or clear
  list, not even empty.
- `[dump]` — presentation only: column reorder for the rendered view; never
  reaches the normalized value.
- Predicate atom spelling `tag.<op> = literal`; operator set is RDR 0003's
  (eq, in, lt, exists appear in evidence). GUESS: the full operator token
  list — the RDR defers the grammar to RDR 0003.

### In-memory typed source (post-parse)

GUESS at struct shapes; contents are fixed by the schema above:

```go
type Model struct {
    ID, Version   // version already gated to 1
    Description   // flow metadata
    Outcomes []string
    Tags     map[string]TagDecl   // provenance, type-model fields, accessor ref
    Accessors map[string]AccessorRef
    Contexts map[string]Context   // inherits, match atoms
    Rules    []Rule               // id, use, match, all, unless, write, clear, escape
    Dump     DumpSettings         // column order only
}
```

### Normalized candidate row (the boundary value peers consume)

Field list is normative (dump contract + round-trip invariant):

```go
type Atom struct {
    Key      string  // declared spelling, byte-exact
    Block    string  // "match" | "all" | "unless" — authored block, retained
    Op       string  // operator token; OpExists verbatim for existence
    Literal          // byte-exact; set literal = member-sorted SEQUENCE,
                     // never a joined string; LiteralTrue/LiteralFalse for exists
}

type CandidateRow struct {
    ModelID, RuleID, Suffix string // identity tuple = the dump sort key;
                                   // Suffix "" unless the rule expanded >1 row
    Locator  Locator               // model id + rule id at minimum; optional
                                   // line/col is diagnostic only — excluded
                                   // from ordering and from round-trip compare
    Outcome  string                // exactly one, lifted from match blocks
    Atoms    []Atom                // sorted (key, block, op, literal)
    Next     map[string]string     // next-state tags incl. "<clear>" entries
    Writes   map[string]string     // accessor-facing writes incl. "<clear>";
                                   // equal to Next today but populated
                                   // independently, never aliased (RDR 0009 A4)
    RequiresOwned []string         // sorted dedup keys of write ∪ clear;
                                   // empty on escape rows
    Escape   []string              // modeled failure classes; non-empty list is
                                   // the SOLE escape discriminator — row kind
                                   // is a derived view column, never a field
}
```

Mapping to the kernel row: equality atoms → `resolve.Row.Match`; all other
atoms (in, comparisons, exists, every unless atom) → the guard field;
Outcome → `Row.Outcome`; RequiresOwned, Writes, Escape as named. GUESS: the
concrete guard encoding handed to the kernel — `Row.Guard` is a `string` on
main and RDR 0007's unimplemented reshape replaces it with atom-shaped data;
this RDR fixes only which atoms are the guard's.

### Load-time category enum (persisted contract, pre-CLI-mapping)

malformed TOML; unknown schema field; missing recognized outcome alphabet;
malformed recognized outcome alphabet (empty / duplicate member / empty-string
member / member containing `#` — each distinguishable); unknown tag; unknown
context; cyclic context inheritance; write to non-owned tag; unknown accessor;
unsupported version; malformed predicate atom (unknown operator, ill-formed
literal, non-boolean existence literal); malformed escape declaration
(including a write block or clear list on an escape rule); malformed outcome
binding (zero/multiple recognized atoms, out-of-alphabet literal, recognized
atom in a guard block); malformed rule id (contains `#`); duplicate rule id;
duplicate model id; reserved_tag_key (RDR 0008's category, in this set).

No runtime persistent resource: model files are source-controlled, dumps are
regenerate-and-delete.

## 4. Top-level pseudo-code: LoadAndNormalize

```
LoadAndNormalize(src bytes) -> Expanded | LoadError:
  # pass 1: version gate BEFORE strict decoding
  head   := permissive-decode only [model] from src        # tolerate unknown keys
  if head missing id or version        -> LoadError(malformed TOML)   # GUESS: category for missing [model] fields
  if head.version != 1                 -> LoadError(unsupported version)

  # pass 2: strict decode, unmapped keys rejected
  m := strict-decode src into Model    # decoder configured to refuse unknown keys
  on parse failure                     -> LoadError(malformed TOML)
  on unmapped key                      -> LoadError(unknown schema field)

  # single-rule / declaration validation (arity: one rule + declarations)
  validate m.Outcomes: non-empty, dedup, no "", no "#"      -> alphabet categories
  validate tags: recognized provenance named "recognized",
                 no other decl named "recognized"           -> reserved_tag_key
  validate accessor refs resolve                            -> unknown accessor
  check duplicate model id (multi-model dump input)         -> duplicate model id  # GUESS: where multi-model arrives
  seen := {}
  rows := []
  for rule in m.Rules:
      if rule.id contains "#"                -> LoadError(malformed rule id)
      if rule.id in seen (byte-equal)        -> LoadError(duplicate rule id)
      if rule has escape: must lack write block AND clear list; classes ⊆
         {no_match, ambiguous_match}         -> malformed escape declaration
      else: rule must have write block; writes/clears target declared owned
         tags only                           -> unknown tag / write to non-owned tag
      atoms := flattenContexts(rule, m.Contexts)   # cycle/unknown ctx refused;
                                                   # merge on full atom identity
      for a in atoms: a.key declared? operator known? literal well-formed for
         operator (existence literal boolean)?     -> unknown tag / malformed predicate atom
      (outcomes, atoms) := liftOutcome(rule, matchExtent(atoms), m.Outcomes)
                                                   # guard-block recognized atom refused
      for oc in outcomes:
          r := CandidateRow{
              ModelID: m.ID, RuleID: rule.id,
              Suffix:  oc if len(outcomes) > 1 else "",
              Locator: locator(m.ID, rule.id [, line/col]),
              Outcome: oc,
              Atoms:   sortAtoms(atoms),           # exists → OpExists/LiteralTrue/False verbatim
              Next:    render(rule.write) + {k: "<clear>" for k in rule.clear},
              Writes:  same rendering, populated independently (no alias),
              RequiresOwned: sortedDedup(keys(rule.write) ∪ rule.clear),  # empty if escape
              Escape:  rule.escape or [] }
          rows.append(r)
  sortExpanded(rows)                    # identity tuple, byte-lexicographic,
                                        # locator excluded; never map iteration
  return Expanded{rows}                 # [dump] settings never consulted here
```

Dump(e) then renders every field from this pre-sorted slice with the derived
kind column; Rows(e) routes atoms into resolve.Row Match/guard and hands the
set — unordered as far as selection is concerned — to the kernel.
