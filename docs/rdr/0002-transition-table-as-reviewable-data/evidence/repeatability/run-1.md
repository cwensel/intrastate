model: claude-opus-5[1m]
variant: full (profile: foundational)

# RDR 0002 Reconstruction — Transition Table As Reviewable Data

Reconstructed from the RDR text alone. Every point the RDR does not fix is
marked **GUESS**.

---

## 1. Public API

**GUESS** — the RDR names no package path, no exported identifiers, and no Go
signatures. It fixes only responsibilities ("the loader", "the normalizer", "the
dump"), the row target (`internal/resolve::Row`), and the category taxonomy. The
Existing Infrastructure Audit says a "New internal package can own sparse source
structs and normalization" without naming it. Package name below is a GUESS:
`internal/table`.

### Entry points

```go
package table // GUESS: package name and path (internal/table)

// Load reads a sparse transition-model TOML document and returns the
// normalized model. Two passes: permissive [model] read for `version`, refuse
// on any value but 1, then strict decode rejecting unmapped keys.
// (Two-pass ordering and strict decoding are normative in the RDR.)
func Load(path string) (*Model, error)              // GUESS: signature/name

// LoadBytes is the same over an in-memory document; `name` supplies the
// locator's file part.
func LoadBytes(name string, data []byte) (*Model, error) // GUESS

// Rows returns the normalized candidate-row set as kernel rows, already in the
// dump's total order. Order is presentational only: the RDR's normative clause
// forbids any positional field a consumer could tiebreak on.
func (m *Model) Rows() []resolve.Row                // GUESS: name/receiver

// Table returns the kernel-facing table: the rows plus the declared
// recognized-outcome alphabet the kernel's `models` check reads.
func (m *Model) Table() resolve.Table               // GUESS: name; resolve.Table
                                                    // is named in the RDR (A8)

// Dump writes the expanded transition table. Every field of the normalized row
// value must be present; [dump] settings may reorder columns, never omit a
// field.
func (m *Model) Dump(w io.Writer) error             // GUESS: signature

// Rendered locator for diagnostics; identifies at least model id + rule id.
func (r RowInfo) Locator() Locator                  // GUESS
```

### Types

```go
type Model struct {          // normalized model (GUESS: name)
    ID        string         // [model].id
    Version   int            // always 1
    Desc      string         // human description
    Outcomes  []string       // root `outcomes` alphabet, non-empty, dup-free,
                             // no empty string, no `#`
    Tags      map[string]TagDecl // keyed by exact byte spelling
    Accessors map[string]AccessorRef
    Rows      []Row          // normalized candidate rows, pre-sorted
    Dump      DumpSettings   // presentation only; never reaches the value
}
```

Every field above is named in the RDR's source-schema list; the Go spelling and
the container choices are **GUESS**.

### Error modes

The RDR fixes the *categories*, and fixes that they are "stable data-level
categories **before** CLI mapping". It explicitly says the CLI code names,
`Group` assignments, and exit mapping are **new design work at implementation**
(A5) — so those are not reconstructible.

Load-time categories (normative minimum, verbatim from the RDR):

| Category | Trigger |
| --- | --- |
| `malformed_toml` | TOML syntax failure |
| `unknown_schema_field` | decoder hit an unmapped key (strict decode) |
| `unsupported_version` | `[model].version != 1`; refused in pass 1 |
| `missing_recognized_outcome_alphabet` | root `outcomes` absent |
| `malformed_recognized_outcome_alphabet` | empty alphabet, duplicate member, empty-string member, or member containing `#` — **each distinguishable** |
| `unknown_tag` | reference (match/write/clear) resolves to no `[tags.<tag>]` |
| `unknown_context` | `use` names a missing context |
| `cyclic_context_inheritance` | `inherits` cycle |
| `write_to_non_owned_tag` | write/clear target is not provenance `owned` |
| `unknown_accessor` | tag declares an accessor id with no `[accessors.<id>]` |
| `malformed_predicate_atom` | unknown operator, or a literal ill-formed for its operator — including any existence literal other than the kernel's two boolean forms |
| `malformed_escape_declaration` | escape rule carrying a write block or a clear list (even empty) |
| `malformed_outcome_binding` | zero or >1 `recognized` atom in the match extent; literal outside the alphabet; a `recognized` atom authored under `guard.all`/`guard.unless` |
| `malformed_rule_id` | rule id containing `#` |
| `duplicate_rule_id` | byte-equal rule id twice in one model |
| `duplicate_model_id` | two models sharing an id |
| `reserved_tag_key` | owned/observed declaration named `recognized`, or a `recognized`-provenance declaration named otherwise (RDR 0008) |

The snake_case spellings are **GUESS** — the RDR only quotes `reserved_tag_key`
and `unknown tag` / `unknown schema field` as prose. **GUESS**: the error is a
single `*LoadError` carrying `Category`, `Locator`, and detail, and the
implementation exposes `errors.Is`-able sentinels per category, because the
Testing Strategy demands assertion "on the category, not the message text".

**Not this API's errors** (explicitly routed elsewhere):

- overlap, gap, dead row, read-before-write, ambiguous expansion → RDR 0006 lint
  (cross-row, by the arity split).
- `unmodeled_outcome`, `owned_state_unavailable`, `guard_unevaluable`, zero/
  multiple survivors → kernel refusals, RDR 0001.
- Expansion counts per rule → a lint *diagnostic*, non-blocking, no threshold.

---

## 2. Three most important internal helpers

### `flattenContexts(src, rule) (atomSet, error)`
Resolves the rule's `use` list transitively through `inherits`, detecting cycles,
and merges every inherited `[context.<id>.match.<tag>]` block with the rule's
local `[rule.match]` into one explicit predicate set. The merge key is the **full
atom identity tuple** `(key, block, operator token, literal)` — never a proper
prefix. Two atoms agreeing on `(key, block, operator)` but differing in literal
are distinct and both survive; identical atoms collapse (idempotent). Inheritance
accumulates, never overrides. The RDR calls out keying on `(block, key, operator)`
as the specific defect this helper must not have, because the loss happens before
the dump's sort can certify anything.

### `bindOutcome(matchAtoms) (outcome string, suffix string, expanded []row, error)`
Scans **only the match extent** — the rule's local `match` plus the inherited
contexts' `match` blocks, *not* `guard.all`/`guard.unless` — for exactly one atom
on key `recognized`, operator `eq` or `in`, whose literal(s) are alphabet members.
Lifts that atom out of the predicate set into the row's outcome field. An `in`
atom expands to one row per member, each suffixed with the outcome literal; a
single-member `in` and an `eq` both produce **one unsuffixed row**, so two
spellings of one edge cannot mint two identities. A `recognized` atom found in a
guard block is a load failure, never lifted — lifting from `unless` would invert
the author's intent silently. Zero, two, or off-alphabet → load failure.

### `emitRow(rule, atoms, outcome, suffix) resolve.Row`
Builds the kernel row and performs the two derivations and the one split the RDR
makes this RDR's:

- **`RequiresOwned` derivation** — the sorted, duplicate-free union of the keys
  named by the write block and the clear list. No authored form exists. Empty on
  escape rows. Stated as a MUST on the normalizer because it is the *only*
  enforcement point (the kernel gates escape rows identically to ordinary ones).
- **Two write-side fields populated independently** — next-state tags and the
  accessor-facing writes, both carrying rendered `<clear>` entries, both
  populated explicitly from the rule and **never by aliasing one to the other**,
  so a future divergence is a compile-time change rather than a silent one.
- **Predicate split by operator** — equality atoms on a declared tag go to
  `Match []Tag`; every other atom (set membership, comparison, existence, and
  *every* `unless` atom regardless of operator) goes to the guard predicate. No
  atom appears in both. Existence atoms are emitted with the kernel constants
  `OpExists` / `LiteralTrue` / `LiteralFalse` verbatim.

**GUESS**: helper names and signatures. **GUESS**: a fourth helper `sortRows` /
`canonicalize` exists to place the pre-sorted sequence before any encoder, since
the RDR forbids emitting from a map or delegating key order to an encoder.

---

## 3. Data model across the boundary

### On-disk (the wire format — normative, fully fixed by the RDR)

Field layout is the Resolve spike layout, locked normatively:

```toml
outcomes = ["round-clean", "verdict-flapping", "reconcile-block", "finalized"]

[model]
id = "rdr-flow"
version = 1                    # only accepted version
# description                  # GUESS: key name for the human description

[tags.status]
provenance = "owned"           # owned | observed | recognized
# + RDR 0003's type model: value kind, optional finite domain, optionality,
#   set-element universe, single-valued marker (cited, not restated here)
accessor = "..."               # optional, observed or owned read-back

[tags.recognized]
provenance = "recognized"      # MUST be named exactly `recognized`

[accessors.<id>]               # accessor reference; execution is RDR 0004's

[context.<id>]
inherits = "<other-context>"
[context.<id>.match.<tag>]     # e.g. status.eq = "Draft"

[[rule]]
id = "prelock-flapping-cap"    # unique per model, byte-equal compare, no `#`
use = ["prelock"]
clear = ["..."]                # rule-level explicit clear list
escape = ["no_match"]          # escape rules only; classes limited to
                               # `no_match` and `ambiguous_match`
[rule.match.<tag>]
[rule.guard.all.<tag>]
[rule.guard.unless.<tag>]
[rule.write]                   # one or more tag assignments

[dump]                         # presentation only; MUST NOT reach the value
```

Key identity: **exact byte equality at every stage** — no case folding, no
trimming, no namespace rewriting. Canonical spelling is the `[tags.<tag>]`
declaration key. Literals carry the same byte-exact identity; a set literal
normalizes to an **ordered member sequence** (members sorted
byte-lexicographically), never a delimiter-joined string.

### Normalized candidate row (the semantic object crossing to the kernel)

| Field | Content |
| --- | --- |
| row identity | `(model id, rule id, expansion suffix)`; suffix present only where the rule expanded, rendered as `rule#suffix` |
| source locator | model id + rule id at minimum; line/column optional diagnostic detail |
| outcome | the single lifted `recognized` literal; kernel filters on it before matching |
| predicate atoms | `(key, operator token, literal, block)` — block ∈ `match` / `all` / `unless`, retained, never folded |
| next-state tags | write block + clear list, including `<clear>` entries |
| writes | the owned-tag writes the accessor layer applies; same rendered set, populated independently |
| required-owned | sorted, dup-free keys from writes ∪ clears; empty on escape rows |
| escape classes | modeled failure-class list; **the sole discriminator** of escape identity |

Row kind (`transition` / `escape`) is a **derived view property, never a row
field** — computed from a non-empty escape class list. Internal indexing (trie,
decision tree, decision DAG) is an implementation detail; the normative object is
the candidate-row set plus locators.

### Expanded table dump (the review surface)

Every field of the row value plus the derived kind column. `[dump]` may reorder
columns, never omit a field. Row order is the identity tuple compared field by
field, byte-lexicographically, absent suffix sorting before any present one;
within a row, atoms sort by `(key, block, operator, literal)`, and next tags,
writes, required-owned keys, and escape classes sort by key. The locator is
**excluded** from ordering. Emitted from a pre-sorted sequence — never a map
iteration, never an encoder's key order.

Round-trip: `parse ∘ normalize` = value identity over the full field list, with
only the locator's line/column detail excluded. **No dump grammar is defined**, so
no textual inverse is claimed; three named lossy sites (`<clear>` unreserved in
the value space, unescaped separators, unreserved suffix separator) would each
have to be closed by a later re-readable-dump RDR.

**Nothing is persisted at runtime.** The RDR states "No runtime persistent
resource is introduced." The on-disk model files are source-controlled; dumps are
regenerable.

---

## 4. Top-level pseudo-code of `Load`

```
function Load(path):
  bytes := readFile(path)                       # GUESS: file read is the caller's
                                                # in the CLI; RDR 0005 owns binding

  # --- pass 1: version gate BEFORE strict field validation (normative) ---
  head := decodePermissive(bytes, onlyTable="model")
  if head.model.version is absent or != 1:
      fail(unsupported_version, locator{path, model?})

  # --- pass 2: strict decode; unmapped keys are a stable refusal ---
  src, err := decodeStrict(bytes)               # decoder MUST reject unmapped keys
  if err: fail(malformed_toml | unknown_schema_field, err.locator)

  # --- document-level validation (single-rule arity: this RDR's) ---
  checkAlphabet(src.outcomes)                   # non-empty, dup-free, no "",
                                                # no member containing "#"
                                                # each violation its own category
  for each tagName, decl in src.tags:
      checkProvenance(decl)                     # owned|observed|recognized
      checkReservedKey(tagName, decl)           # RDR 0008: recognized <-> `recognized`
      carryTypeModel(decl)                      # RDR 0003's fields, forwarded losslessly
      if decl.accessor and not in src.accessors: fail(unknown_accessor)
  detectContextCycles(src.contexts)             # cyclic_context_inheritance
  requireUniqueModelID(src.model.id)            # duplicate_model_id

  rows := []
  seenRuleIDs := {}
  for each rule in src.rules:                   # source order is irrelevant
      if rule.id contains "#":       fail(malformed_rule_id)
      if rule.id in seenRuleIDs:     fail(duplicate_rule_id)
      seenRuleIDs.add(rule.id)

      matchAtoms := flattenContexts(src, rule)  # unknown_context on a bad `use`;
                                                # merge on FULL atom identity
      guardAtoms := atoms(rule.guard.all, block="all")
                  + atoms(rule.guard.unless, block="unless")

      for each atom in matchAtoms + guardAtoms:
          if atom.key not in src.tags: fail(unknown_tag)
          checkOperatorAndLiteral(atom)         # malformed_predicate_atom;
                                                # exists literal must be one of
                                                # the kernel's two boolean forms
          if atom.key == "recognized" and atom.block != "match":
              fail(malformed_outcome_binding)   # never lift from a guard block

      outcome, members := bindOutcome(matchAtoms)   # exactly one recognized atom,
                                                    # eq or in, alphabet members
      predicate := (matchAtoms - recognizedAtom) + guardAtoms

      if rule.escape is non-empty:
          if rule.write or rule.clear present:  fail(malformed_escape_declaration)
          checkEscapeClasses(rule.escape)       # only no_match | ambiguous_match
          writes, next, requiresOwned := {}, {}, {}
      else:
          if rule.write absent:                 fail(malformed_escape_declaration)  # GUESS: category
          for each target in rule.write.keys + rule.clear:
              if target not declared:           fail(unknown_tag)
              if provenance(target) != "owned": fail(write_to_non_owned_tag)
          next  := render(rule.write, rule.clear)      # `<clear>` entries included
          writes := render(rule.write, rule.clear)     # populated independently,
                                                       # NOT aliased to `next`
          requiresOwned := sortedUnique(keys(rule.write) + rule.clear)

      for each m in members:                    # 1 member -> unsuffixed row
          suffix := (len(members) > 1) ? m : ""
          match, guard := splitByOperator(predicate)   # eq-on-declared-tag -> Match;
                                                       # everything else + all
                                                       # `unless` -> guard
          rows.append(Row{
              identity: (src.model.id, rule.id, suffix),
              locator:  locatorOf(rule),
              outcome:  m,
              match:    match, guard: guard,
              next: next, writes: writes,
              requiresOwned: requiresOwned,
              escape: rule.escape,
          })

  sortRows(rows)                                # by identity tuple, byte-lex;
                                                # locator excluded; atoms and
                                                # every set field sorted by key
  return Model{..., Rows: rows}                 # cross-row checks (overlap, gap,
                                                # dead row, read-before-write)
                                                # are RDR 0006's, not run here
```

---

## GUESS index

| # | Guess |
| --- | --- |
| G1 | Package path/name `internal/table`; all Go identifiers, signatures, receivers |
| G2 | `Load`/`LoadBytes`/`Rows`/`Table`/`Dump` as the public surface — the RDR names capabilities, not functions |
| G3 | snake_case spellings of the load-time category constants |
| G4 | Single `*LoadError` type with `Category` + `Locator`, plus per-category sentinels |
| G5 | Struct shapes of `Model`, `TagDecl`, `AccessorRef`, `DumpSettings`, and container choices |
| G6 | Helper names `flattenContexts` / `bindOutcome` / `emitRow` / `sortRows` |
| G7 | TOML key name for the human description in `[model]` |
| G8 | Category assigned to an ordinary rule missing a write block (RDR says a write block is required but does not name the category) |
| G9 | Whether `Table()` composes rows + alphabet, and whether `resolve.Table` is a struct or interface |
| G10 | That file reading lives outside the loader (CLI/RDR 0005 owns path binding) |
| G11 | Dump output syntax entirely — the RDR fixes the field list and row order but explicitly defines no grammar |
| G12 | That the alphabet, tag declarations, and accessor map survive into the returned model as addressable state for RDR 0006 lint |
