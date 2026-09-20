model: claude-sonnet-5
variant: full (profile: foundational)

# Reconstruction run 3 — RDR 0030 (computed write-value grammar / step write)

Source: contracts C1–C3, MVV, S1–S11 (`0030:§normative-contracts`,
`0030:§mini-check-tables`, `0030:§load-bearing-decisions`,
`0030:§minimum-viable-validation`, `0030:§testing-strategy`), plus the
illustrative TOML and metadata (Overrides on `0002:C4`, `0002:C13`).
Widened past the contract spans into `§mini-check-tables` (authority /
oracle / fidelity / disposition / trace tables) and
`§load-bearing-decisions` because the contract text leaves several
signatures underdetermined (exact Go types/return shapes for the moved
`intWidth`, the new `renderWrites`/`expand` parameter widening, and the
exact identity-ordering function) — those tables pin the decisions C1's
prose only implies. Also widened into `§illustrative-code` for concrete
TOML shape and into `§metadata` for the two RDR-0002 overrides that bound
this record's blast radius.

---

## 1. Public API — signatures, types, error modes

This RDR does not add a CLI verb; it widens the **load-time grammar** for
a rule's `write` block and widens two **internal** function signatures
that the loader/normalizer pipeline already has. No package outside
`internal/table` (plus the `internal/resolve` shim it creates) is a
public API surface here — "public" in this record's terms means
"in-repo cross-package contract."

### 1.1 TOML surface (the actual "API" an author sees)

```
[rule.write]
<tag> = { step = <n> }      # n: non-zero TOML integer (not float)
```

- Admitted only when `<tag>`'s declared kind is `int` with both `min` and
  `max` set, or `enum` with a non-empty `domain`.
- Refused (load-time) on `bool`, `set`, `scalar`, an `int` missing a
  bound, a zero step, a float step (`1.0`), any other table key/shape,
  or a bare table with `add`/`next` keys — all under
  `malformed_tag_declaration`.
- `[initial]` block: refused under `malformed_initial_declaration`.
- Predicate literal position (`match`/`guard.all` atom, or `in` member):
  refused under `malformed_predicate_atom`.

### 1.2 Widened Go signatures (internal/table/normalize.go)

GUESS on exact parameter names/order (RDR states the widenings' *shape*
and *reason*, not the literal diff):

```go
// today (inferred baseline, not directly quoted by the RDR):
// func renderWrites(rule *Rule) ([]TagValue, []string, error)
// func expand(writes []TagValue, predicates []Atom, ...) ...

// after this record (two widenings, both local to normalize.go):
func renderWrites(
    rule *Rule,
    predicates []Atom, // NEW: merged predicate set, needed by the admits filter
) ([]TagValue, []string, error) // GUESS: []TagValue widened to carry a
                                  // non-literal step spec alongside literals;
                                  // RDR says "widening that pair" but not the
                                  // concrete carrier type/field name.

func expand(
    writes []TagValue, // now may contain step specs, not just literals
    predicates []Atom,
    // ...unchanged remainder, GUESS on exact remaining params
) (rows []Row, err error)
```

C1 pins: (a) `renderWrites` returns `([]TagValue, []string, error)` today
and cannot carry a non-literal without widening; (b) `expand` takes
`writes []TagValue`; (c) the merged predicate set must reach
`renderWrites`, which it does not today; (d) both widenings are local to
`normalize.go` and add no new call site. The RDR does not name the new
field/type used to carry a step spec inside `TagValue` — GUESS: an
optional `*StepSpec` field alongside the existing literal, since C1 says
"the spec itself rides alongside as one more per-key record."

### 1.3 Relocated / new internal/resolve shim (Phase 2, per the
authority mini-check-table)

```go
// internal/resolve (NEW package or extended existing one — GUESS on
// whether resolve already exists pre-record; C1 says "the relocated
// Evaluator" implying resolve.Evaluator predates this record)

// MOVED as-is, type-free core (was internal/guard/declaration.go):
func intWidth(minV, maxV int) (int, bool)

// STAYS in internal/guard, thin wrapper over the moved core:
func IntDomain(d table.TagDecl) []int
```

`guard::IntDomain` cannot move because `resolve` must not import `table`
(cyclic: `table` imports `resolve`). `intWidth` moves because it is
already type-free (`(minV, maxV int) (int, bool)` — no `table` type in
its signature).

### 1.4 Error modes (exhaustive, per the `disposition` mini-check-table)

| Condition | Category | Loud/quiet |
|---|---|---|
| step on int+bounds / enum+domain | admitted, N rows | loud (visible in dump) |
| step on bool/set/scalar/unbounded int | `malformed_tag_declaration` | loud |
| int width not representable (`max-min` negative or `== math.MaxInt`) | `malformed_tag_declaration` | loud |
| stepped value overflows int range (checked add) | `malformed_tag_declaration`, names cell+bound | loud |
| zero step / float step / other keys | `malformed_tag_declaration` | loud |
| key both stepped and in `clear` | `malformed_tag_declaration` | loud |
| `#` in domain member / duplicate domain member (step-scoped) | `malformed_tag_declaration` | loud |
| zero admitted cells | `malformed_tag_declaration` (NEW check) | loud |
| stepped value leaves domain | `malformed_tag_declaration`, names first failing cell in domain order | loud |
| `{step=n}` under `[initial]` | `malformed_initial_declaration` | loud |
| `{step=n}` as predicate literal | `malformed_predicate_atom` | loud |
| `unless` over-admits interior cell | admitted; dead row; advisory `graph-redundant-row` at lint | quiet at load, loud-advisory at lint |
| excluded cell claimed by nobody | admitted; `graph-coverage-gap` at lint | loud at lint |

The refusal envelope is the existing `internal/table/category.go::Failure`
struct: `Category, Detail, Offending, Remedy, Rule, Line`. This record
populates none of the structured fields beyond `Category`+`Detail` text;
`Offending`/`Remedy` stay reserved for `0008:C3`'s declaration-name pair.

---

## 2. Three most important internal helper functions

1. **The admits filter (inside `renderWrites`, C1's core mechanism)**
   Responsibility: for a stepped tag, compute the set of "admitted
   cells" — domain members (`{min..max}` for int, authored `domain`
   order for enum) satisfying the CONJUNCTION of every positive atom the
   rule authors on that tag (`guard.all eq/in/lt/lte/gt/gte`, `match`
   `eq`/`in`), evaluated per member as the runtime evaluator would,
   against the authored literal (not against a sibling choice point's
   already-chosen member). `unless` atoms are never consulted. This is
   the single most load-bearing function in the record — everything else
   (row count, refusal timing, subsumption) derives from its output.
   GUESS at exact name — RDR describes it as "C1's admits filter" but
   does not give it a Go identifier; likely a new unexported function in
   `normalize.go` alongside `renderWrites`.

2. **`intWidth(minV, maxV int) (int, bool)`** (moving from
   `internal/guard/declaration.go` to the `internal/resolve` shim)
   Responsibility: compute whether an int declaration's width
   (`max-min+1`) is representable without wrapping, returning the width
   and a representability bool. Single source of truth shared by the
   loader (C1's kind arm, deciding whether a stepped int declaration can
   even be enumerated) and lint's `Cardinality` computation — must stay
   the SAME rule so the two never drift (a declaration lint saturates to
   a ceiling must be one the loader also refuses to enumerate).

3. **The suffix/identity emission logic in `normalize.go::expand`**
   (existing function, extended)
   Responsibility: decide, per choice point, whether to append a suffix
   element to a row's identity tuple. For the pre-existing `in`
   expansion, `suffixed := expanding && len(members) > 1` (single-member
   `in` mints no suffix). This record adds a **second, independent**
   local for the step point that is unconditionally true — a stepped row
   is never the authored row, so it always carries its cell, even at one
   admitted cell. Also owns the ordering rule for multiple suffix
   elements: step points (and other choice points) order by KEY alone
   (not `(Key,Block,Operator,Literal)` — a step point isn't an `Atom` and
   its `guard.all eq` literal is per-cell/row-dependent, so that 4-tuple
   doesn't apply).

---

## 3. Data model — persisted / cross-boundary

Nothing new is **persisted** (no new file format, no schema migration).
The record widens two **in-memory, cross-package carriers**:

### 3.1 `TagValue` (internal/table) — widened

Today: literal-only. After this record it must be able to carry either
a literal write value OR a step spec (magnitude `n`, sign, target kind).
GUESS at fields (RDR does not give the struct literally):

```go
type TagValue struct {
    Key     string
    // exactly one of:
    Literal any        // existing: int/string/bool literal
    Step    *StepSpec  // GUESS: NEW, non-literal write descriptor
}

type StepSpec struct { // GUESS — name/shape not given by RDR
    N int // non-zero
}
```

### 3.2 `Row` (internal/table/model.go) — carriers, unchanged shape,
new population rule

- `Row.Writes` and `Row.NextTags` — BOTH populated per expanded row from
  the SAME stepped value, never by aliasing one to the other (`0002:C15`
  invariant, reaffirmed not changed). Consumers reading `NextTags`
  without `Writes` (`graph` per-row `next`, `flow next`, `guard`'s
  written-key lint, `dump`'s next column) depend on this.
- `Row.Atoms` — retained authored atoms, MINUS a subsumed `match in`
  (rewritten to `match eq = <cell>` per row), PLUS the expansion's own
  `guard.all eq = <cell>` atom. Exactly two atoms are contributed by
  this mechanism per row on the stepped tag (the rewritten `eq` where an
  `in` was subsumed, and the generated `eq`); every other authored atom
  on the tag is retained unchanged.
- Row identity: existing tuple `(model id, rule id, suffix…)`, cell
  appended as one more suffix element. Suffix elements for multiple
  choice points on one rule (e.g. a step on `attempt` + unsubsumed `in`
  on `mode`) order by key name: `retry#<attempt>#<mode>`.

### 3.3 `Failure` (internal/table/category.go) — unchanged shape

`Category, Detail, Offending, Remedy, Rule, Line` — no new field. All
four "parts" of a bound-refusal detail (rule, tag, cell, result) ride
the `Detail` string, not structured fields.

### 3.4 Wire/graph document — unchanged

`intrastate.graph/1` gains no member. Per-row findings/graph rows publish
suffixed identity `rule#cell` in the `rule`/`Element`/`identity` slots
(rendering-only change); the in-package join key for
`graphlint/coverage.go::groupHasOverlap` stays the bare `Span` (or
recovered authored id) — never the suffixed `Rule` field, to avoid
silently vacuous-passing the `ambiguous_match` coverage arm.

---

## 4. Top-level pseudo-code of the main operation (step-write expansion)

```
// Entry: internal/table/normalize.go, extending normalizeRule/expand
// for a rule whose write block contains a `{ step = n }` table.

func normalizeRule(rule, declaredTags) (rows []Row, err error):
    predicates := assemblePredicateSet(rule)          // existing
    writes, clears, err := renderWrites(rule, predicates)  // WIDENED: now takes predicates
    if err: return err                                  // malformed_tag_declaration etc.

    // --- inside renderWrites, per write entry ---
    for each writeEntry in rule.write:
        if writeEntry is `{ step = n }`:
            tag := declaredTags[writeEntry.Key]
            if not (tag.Kind == int && tag.Min/Max set) &&
               not (tag.Kind == enum && len(tag.Domain) > 0):
                refuse malformed_tag_declaration  // wrong kind/unbounded
            if n == 0 || not integer:
                refuse malformed_tag_declaration
            if tag.Kind == int:
                width, ok := intWidth(tag.Min, tag.Max)   // moved shim call
                if !ok: refuse malformed_tag_declaration   // unrepresentable width
            if tag.Kind == enum && (hasDuplicateMember(tag.Domain) || hasHashMember(tag.Domain)):
                refuse malformed_tag_declaration  // step-scoped guard only
            if writeEntry.Key in rule.clear:
                refuse malformed_tag_declaration  // stepped + cleared

            stepSpec := StepSpec{N: n}
            track stepSpec alongside literalWrites   // placeholder in assignments map
        else:
            // existing literal handling, unchanged

    if predicates.hasOwned() [existing 0002:C14 gate]:
        RequiresOwned := derive from maps.Keys(assignments)  // stepped key present via placeholder

    return writes, clears, nil

// --- expand: mint literal rows from any step spec ---
func expand(writes []TagValue, predicates []Atom) (rows []Row, err error):
    for each stepped TagValue in writes:
        admitted := admittedCells(tag, predicates)   // C1's admits filter:
                                                        // conjunction of AUTHORED
                                                        // positive atoms per member,
                                                        // `unless` never consulted
        if len(admitted) == 0:
            refuse malformed_tag_declaration  // zero admitted cells

        emit ONE choice point for this stepped tag (never two,
             even if a match `in` on same tag also present):
        if a match `in` on this tag is present (authored or inherited):
            subsume it: drop as separate choice point,
            intersect its members into admitted (already done via
                conjunction), rewrite surviving occurrences to `eq` per row

        for each cell in admitted (domain order):
            value := checkedAdd(cell, n)             // int arm
                     or domainMember(cell, n)         // enum arm (position offset)
            if overflow or value not representable:
                refuse malformed_tag_declaration (names cell + bound)  // BEFORE render
            if value outside tag domain (past max/min, past/before enum ends):
                refuse malformed_tag_declaration (first failing cell in domain order)

            row := cloneAuthoredRule()
            row.Atoms := retainedAuthoredAtoms(tag)   // minus subsumed `in`, plus eq
            row.Atoms += {Block: guard.all, Op: eq, Key: tag, Literal: cell}
            row.Writes[tag] = value                    // BOTH carriers
            row.NextTags[tag] = value                   // BOTH carriers, same value
            row.Identity.Suffix += cell                 // ALWAYS appended,
                                                          // even for single-cell case
            rows = append(rows, row)

    // order multiple suffix elements (if >1 stepped/branching choice
    // point on this rule) by KEY name, not by block/operator/literal
    sort rows' suffix elements by key

    return rows, nil
```

Line count target 20–40 achieved (~40 lines above, condensed).

---

## Notes on silences / GUESSes made explicit

- GUESS: exact Go field/type names for the widened `TagValue` write
  carrier and the step-spec struct — RDR describes the *need* and the
  *constraint* (must ride alongside per-key, must not break
  `RequiresOwned`'s `maps.Keys(assignments)` derivation) but never
  spells the struct.
- GUESS: the admits-filter function's identifier — RDR references it
  only descriptively ("C1's admits filter").
  RDR names concrete comparison/render functions it reuses
  (`compareAtoms`, `compareRows`, `renderSetLiteral`, `seamValue`) but
  not a name for the new filter itself.
- GUESS: whether `internal/resolve` is a brand-new package or already
  exists holding `Evaluator` pre-record — the mini-check-table implies
  `Evaluator` is being relocated INTO resolve as part of this same
  record ("Phase 2, which is also where Evaluator moves"), so resolve
  may be newly created by this work rather than pre-existing. Treated as
  new-or-newly-populated either way; does not change the reconstruction.
- Not a silence: refusal categories, admitted-cell semantics, both-carrier
  population, suffix/identity rules, and the lint attribution
  (suffixed `rule#cell` vs bare join key) are all fully pinned by C1–C3
  and the mini-check tables and reproduced above without guessing.
