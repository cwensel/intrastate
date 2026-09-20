model: claude-opus-5[1m]
variant: full (profile: foundational)

## Spans widened past

Contracts alone (C1, C2, C3) left the following under-determined, so I read
wider before guessing. Each is named here so the diff does not score a read
as a GUESS:

- **Function signatures.** C1 names `renderWrites` returning
  `([]TagValue, []string, error)` and `expand` taking `writes []TagValue`,
  and says BOTH widen — but never gives the widened shapes. Widened to
  `§phase-2-expand-into-cells`, which supplies `expand`'s full current
  signature and states the split of responsibility between `renderWrites`
  (declaration-dependent half) and `expand` (row minting). The widened
  record type's name and fields are still unnamed anywhere → GUESS.
- **The extracted shim's signature.** C1 and the `authority` mini-check
  name the shim but not its spelling. Widened to `§phase-2` +
  `§prerequisites`, which fix it as taking a `resolve.GuardAtom` and the
  held value (a string, per 0012 prerequisite), returning the three-valued
  verdict. Name of the function → GUESS.
- **Step-spec type.** C1 says the stepped key is placed in `assignments`
  carrying "a placeholder the expansion replaces per cell" and the spec
  "rides alongside as one more per-key record". Neither the placeholder
  sentinel nor the record type is named. Widened to §illustrative-code and
  the `trace` mini-check; both are behavioural, not structural → GUESS.
- **Step order of the load pipeline.** C1 fixes exactly two orderings
  (checked add BEFORE render; step spec derived BEFORE the clear-list pass)
  and C2 fixes the domain-order-first reporting rule. The rest of the
  ordering came from the `trace` mini-check table (8 steps), which is the
  record's own walk and is treated here as fixed, not guessed.
- **Enum bound arithmetic.** C2 says an `enum` has "no stepped value" out
  of range and the detail names a signed offset. Widened to the `trace`
  and `disposition` tables to confirm the `int` and `enum` arms share one
  refusal category and differ only in the detail's result slot.
- **Lint suffix/join split.** C3 states the emission/read-back split but
  not the recovery function. Widened to `§mini-check-tables` `authority`
  row, which names `Span` as the canonical join key with "the authored id
  recovered from the suffixed form" as the alternative. Which of the two is
  chosen → GUESS (I take `Span`, the row the table calls canonical).

Not widened: I did not read §consequences, §risks, §testing-strategy S1–S11
bodies, or §decision-rationale. The four items below are reconstructable
without them.

---

## 1. The module's public API

This record adds no CLI verb and no exported package API to `internal/table`.
Its public API is (a) the authoring grammar on the wire, (b) the refusal
vocabulary, and (c) one newly public function in `internal/resolve`.

### 1a. Authoring surface (the TOML grammar — the real public API)

```toml
[rule.write]
<tag> = { step = <n> }      # n: non-zero TOML integer, exactly one key
```

Admitted only where the tag's declared kind is:
- `int` with BOTH `min` and `max` declared, and a representable width
  (`max - min` non-negative and not equal to `math.MaxInt`);
- `enum` with a non-empty `domain`.

Refused on `bool`, `set`, `scalar`, and `int` missing a bound. Admitted on
the write-block path only — never under `[initial]`, never as a predicate
literal.

### 1b. Newly public in `internal/resolve` (the relocation)

```go
// package internal/resolve — MOVED from internal/guard/grammar.go, unchanged.
type Evaluator struct{}                       // stateless; field-count assertion follows it
func (Evaluator) Evaluate(...) GuardVerdict    // signature unchanged from guard's

// MOVED from internal/guard/declaration.go, type-free, relocates as-is.
func intWidth(minV, maxV int) (int, bool)

// NEW, extracted from internal/guard/product.go::valueSatisfies.
// Takes a resolve.GuardAtom and the held value as a STRING (not a
// table.Atom, not a typed value) because resolve must not import table.
func AtomAdmits(atom GuardAtom, value string) GuardVerdict   // GUESS: name
```

`intWidth` is exported or package-visible as the phase requires; the record
fixes only that it MOVES and that loader and lint share the one rule.
**GUESS:** it is exported as `resolve.IntWidth` since `internal/guard` and
`internal/table` both call it across a package boundary.

`guard::IntDomain(d table.TagDecl) []int` does NOT move. It stays in
`guard` as a thin wrapper unpacking the decl and calling the moved core.

### 1c. Error modes (all load-time refusals; no runtime error mode added)

Every refusal this record mints on the write-block path files under
`malformed_tag_declaration`:

| Condition | Detail names |
| --- | --- |
| table shape other than exactly `{ step = <n> }` | the one admitted form |
| `step = 0`, or a non-integer step (TOML float `1.0` included) | the one admitted form |
| kind is `bool`/`set`/`scalar`/unbounded `int` | the kind and the admitted kinds |
| `int` width not representable (`max-min` negative or `== math.MaxInt`) | the declaration |
| stepped value leaves the domain at an admitted cell | rule, tag, cell, result |
| stepped value computed past `int` range (checked add, BEFORE render) | rule, tag, cell, the bound passed |
| zero admitted cells | rule, tag (new check, `0003:C6` disposition) |
| `#` in an `enum` domain member (step-scoped) | the member |
| duplicate `enum` domain member (step-scoped) | the member |
| key BOTH stepped and named in the rule's `clear` list | the key |

Two refusals keep their OWN pre-existing categories and are deliberately
NOT merged:

| Condition | Category |
| --- | --- |
| `{ step = n }` under `[initial]` | `malformed_initial_declaration` |
| `{ step = n }` as a predicate literal | `malformed_predicate_atom` |

The refusal carrier is `internal/table/category.go::Failure` with fields
`Category`, `Detail`, `Offending`, `Remedy`, `Rule`, `Line`. **This record
adds no field.** All four parts of a bound failure ride the `Detail`
STRING, with the cell named FIRST so assertions anchor at a stable
position. `Failure.Rule` is NOT populated with a rule id — it is a
direction identifier (`RuleKernelOwned` / `RuleAuthorMustRename`) that
`internal/cli/flow_input.go` gates a rename hint on.

No new CLIError code, no new graph-lint finding code, no new
`intrastate.graph/1` member, no new `internal/table::Categories()` member.
Therefore no `0029:C4` stability tier is owed.

---

## 2. The three most important internal helper functions

### (i) `(*loader).renderWrites` — widened; the declaration-dependent half

Current signature `renderWrites(...) ([]TagValue, []string, error)`. Two
widenings, both local to `normalize.go`, neither adding a call site:

1. It takes the rule's MERGED predicate set as a new parameter. Today
   `::normalizeRule` assembles `predicates` one line above the call and
   hands it only to `::expand`; the admits filter needs it.
2. Its return widens so the write list can carry a non-literal — a per-key
   step record rather than only a `TagValue`.

**GUESS (shape):**
```go
func (l *loader) renderWrites(
    src ruleSource, predicates []Atom,
) (writes []writeSpec, clears []string, err error)

type writeSpec struct {
    Key     string
    Literal string   // set when the write is an authored literal
    Step    *stepSpec
}
type stepSpec struct {
    Cells    []string // admitted cells, in the tag's own domain order
    Literals []string // parallel: the rendered stepped value per cell
}
```

Responsibilities it owns (all needing the `TagDecl` only it holds):
grammar validation of the inline table; the kind gate; the `int` width
check via the moved `intWidth`; admitted-cell evaluation (the conjunction
over authored positive atoms); the zero-cell refusal; the `#`-member and
duplicate-member guards; the checked `int` add; the per-cell `conform`
bound check (C2). Critically, it derives the step spec BEFORE the
clear-list pass so the step/clear collision is detectable, and it places
the stepped key in `assignments` with a placeholder — because `0002:C14`
derives `RequiresOwned` from `maps.Keys(assignments)`, and a spec tracked
outside the map would silently drop the key from the kernel's owned-state
gate.

### (ii) `expand` — the row minter, one more choice point

Current: `expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row`.
Only `writes` widens to the record above. `expand` does for a step point
exactly what it does for an `in` today: take a per-key member list, mint
one row per combination in the existing expansion product (`0002:C13`),
carry the supplied literal, and append the suffix element.

Step-point specifics it owns:
- **Subsumption.** One choice point per stepped tag, never two. A match
  `in` on the same tag — authored locally OR inherited from a context —
  is subsumed, and its members join the admitting conjunction instead of
  expanding separately. A subsumed `in` is REWRITTEN per row to
  `match eq = <cell>`, exactly as the existing `case expanding:` arm does;
  never retained in `in` form (a retained multi-member `in` reaches the
  kernel as `members[0]` via `::seamValue` and every later-cell row would
  never match).
- **Per-row atoms.** The stepped tag carries the retained authored set
  PLUS exactly two: the rewritten `match eq = <cell>` (only where an `in`
  was subsumed) and the expansion's own `guard.all eq = <cell>`. Every
  other authored atom on the tag, bounds included, is retained unchanged.
- **Suffix.** Appended at EVERY admitted cell, including when exactly one
  is admitted. `normalize.go`'s `suffixed := expanding && len(members) > 1`
  is a per-choice-point local, so the step point computes its own as
  constant-true; the `in` arm is unchanged.
- **Both carriers.** `Row.Writes` AND `Row.NextTags` are each populated
  per row from the same stepped value, never by aliasing one to the other
  (`0002:C15`).
- **Ordering.** Step points order by KEY alone (a step point is not an
  `Atom`, so `compareAtoms`' `(Key, Block, Operator, Literal)` tuple does
  not apply, and its literal is row-dependent anyway). Key alone is total
  because subsumption leaves at most one suffix-emitting step point per
  key. The same key comparison places step points among the other
  candidates — unsubsumed `in` points on other keys, and the outcome
  point.

### (iii) The extracted admits shim in `internal/resolve`

The two-valued core of `internal/guard/product.go::valueSatisfies`:
render the atom's literal, call the moved `Evaluator`, return the
three-valued verdict. It carries the canonicalizing set renderer
(`reach.go::renderSetLiteral`'s sorted/compacted form, per the `0003`
set-literal clause), not `assignment.go::renderSet`'s as-authored form.

**The shim is the comparison, never the policy.** Each of its three
callers keeps its own undecided-arm disposition at its own call site:
`guard::valueSatisfies` passes the verdict through; `graphlint::atomAdmitsValue`
keeps its `!= GuardFalse` collapse (false-green); the loader EXCLUDES the
member (C1). The loader's exclusion is a soundness fence, not a live
branch — for the kinds C1 admits the undecided arm is unreachable.

Three non-test call sites repoint: `guard/product.go::valueSatisfies`,
`graphlint/reach.go::atomAdmitsValue`, and `cli/flow_resolve.go::guardSeam`
(which constructs the evaluator directly). Each renders its own
`table.Atom` into a `resolve.GuardAtom` at the call site.

---

## 3. The data model across the boundary

### Persisted / on the wire — UNCHANGED

Nothing new is persisted. After normalization, no surface distinguishes an
expanded row from an authored literal row: both carriers hold literals
only. The `intrastate.graph/1` document type gains no member and decodes
unchanged.

### `Row` (the normalized model, in-memory across the loader→kernel boundary)

| Field | Value on an expanded row |
| --- | --- |
| `Atoms` | authored atoms retained, plus `guard.all eq = <cell>` on the stepped tag, plus a rewritten `match eq = <cell>` where an `in` was subsumed |
| `Writes` | the per-cell stepped literal |
| `NextTags` | the SAME per-cell stepped literal, independently populated |
| identity suffix | the cell appended as one more element |
| `RuleID` | the AUTHORED rule id — unchanged by expansion |
| `Span` | `model:rule` — unchanged by expansion; the canonical lint join key |
| `Fingerprint` | hashes `Row.Atoms` before `NextTags`, so each expanded row is distinct via its `guard.all eq` |

### Row identity

The existing tuple `(model id, rule id, suffix…)`, with the cell appended
as one more suffix element. No new identity kind. The element is the BARE
cell, not tag-qualified (the `in` precedent). `#` is banned in rule ids, so
`retry#3` can never be an authored id. Two step points on distinct keys
yield `retry#<attempt>#<tier>`, ordered by key name; a step on `attempt`
with an `in` on `mode` yields `retry#<attempt>#<mode>`.

This NARROWS `0002:C13`'s "suffix non-empty exactly when the rule produced
more than one row" to the `in` expansion it was written about.

### Published identity, per surface

| Surface | Publishes |
| --- | --- |
| `dump` `identity` column | suffixed `rule#cell` |
| graph document `identity` row field | suffixed `rule#cell` |
| graph-lint per-row finding `Rule` | suffixed `rule#cell` |
| graph-lint per-row finding `Element` (`graph-overlap`, `graph-redundant-row`) | suffixed `rule#cell` — BOTH slots or neither |
| graph-lint group-level findings (`coverage.go::firstRuleID`) | BARE authored id — unchanged |
| `flow next` `rule`, `flow resolve` `rule` | BARE authored id |
| graph document edge `rule` | BARE authored id |

Rule: a surface that names a ROW publishes the suffixed identity; a surface
that names a RULE publishes the authored id.

### The lint suffix / join-key split

The suffixed form is a RENDERING applied at emission. The in-package
read-back — `coverage.go::groupHasOverlap`, which matches an already-emitted
`graph-overlap` finding's `Rule` and `Element` against `ruleIDsOf(g)` to
decide whether the group's `ambiguous_match` arm is reachable — joins on a
key that does NOT carry the suffix: the row's `Span` (**GUESS**: `Span` is
taken over the alternative of recovering the authored id from the suffixed
form; the record's `authority` table names `Span` first and calls it the
stable key). No in-package consumer may resolve a per-row finding's `Rule`
against authored `[[rules]]` ids by equality. Out-of-repo consumers joining
on `rule` are affected identically; `span` is the stable key.

### Data model NOT extended

`Failure` gains no field. `0002:§round-trip-inverse-invariants` is NOT
extended to the literal↔step fixture pair: the invariant is plan-equality
per (state, outcome) cell and finding-set equality under the
`(code, key, dimension, class, reason)` projection — never byte-equality.
Row identity (`rule`, `span`, `element`, `fingerprint`, dump `identity`)
differs by construction and is the lossy-exemption set.

---

## 4. Top-level pseudo-code of the main operation

The load-time expansion of one rule carrying a step write.

```
normalizeRule(src, decls):
  predicates := mergePredicates(src, inheritedContexts)      # existing
  writes, clears, err := renderWrites(src, predicates)       # widened: takes predicates
  if err != nil: return err
  outcome := outcomeAtom(src)
  rows := expand(baseRow(src), predicates, outcome, writes)
  return rows

renderWrites(src, predicates):                               # *loader — holds TagDecl
  specs := []
  for key, raw in src.Write:                                 # derived BEFORE the clear pass
    if raw is not an inline table:
      specs += literalSpec(key, conform(decls[key], raw))     # existing path, unchanged
      continue
    if raw.keys != {"step"} or raw.step is not an integer or raw.step == 0:
      refuse malformed_tag_declaration: "the one admitted form is { step = <n> }"
    d := decls[key]
    switch d.Kind:
      case int:  if d.Min or d.Max unset: refuse malformed_tag_declaration
                 if _, ok := resolve.IntWidth(d.Min, d.Max); !ok: refuse   # shared rule
                 cells := d.Min .. d.Max                                   # ascending
      case enum: if len(d.Domain) == 0: refuse
                 if d.Domain has a duplicate or a "#"-bearing member: refuse  # step-scoped
                 cells := d.Domain                                         # AUTHORED order
      default:   refuse malformed_tag_declaration                          # bool/set/scalar
    admitted := []
    for c in cells:                                          # domain order, for C2 reporting
      ok := true
      for a in predicates where a.Key == key and a.Block != "unless" and a.Positive:
        if resolve.AtomAdmits(a.AsGuardAtom(), c) != GuardTrue: ok = false  # EXCLUDE undecided
      if ok: admitted += c
    if len(admitted) == 0:
      refuse malformed_tag_declaration                       # zero-cell, 0003:C6 disposition
    literals := []
    for c in admitted:                                       # first failure wins, domain order
      v, ok := checkedAdd(index(c), raw.step)                # BEFORE any render (A13)
      if !ok: refuse malformed_tag_declaration: cell c, the bound passed
      if v outside d's domain: refuse: rule, tag, cell c, result(v or signed offset)
                                       + any unless atom on key, "did not exclude the cell"
      literals += render(v)                                  # byte-identical to valueAssignments
    if key in src.Clear: refuse malformed_tag_declaration     # step/clear collision
    specs += stepSpec(key, admitted, literals)
    assignments[key] = placeholder                           # so 0002:C14 RequiresOwned sees it
  return specs, clearList(src, assignments), nil

expand(base, predicates, outcome, writes):
  points := []
  for spec in writes where spec.Step != nil:
    predicates = rewriteSubsumedIn(predicates, spec.Key)     # in -> per-row match eq
    points += stepPoint(spec)                                # ONE per key; always suffixes
  points += inPoints(predicates, excluding stepped keys) + outcomePoint(outcome)
  sortByKey(points)                                          # step points: KEY alone
  for combo in cartesianProduct(points):                     # 0002:C13, unchanged
    r := base
    r.Atoms  = retainedAtoms + combo.guardAllEq + combo.rewrittenMatchEq
    r.Writes, r.NextTags = combo.literal, combo.literal      # BOTH, never aliased (0002:C15)
    r.Suffix = base.Suffix + combo.element                   # appended at every cell
    rows += r
  return rows
```

**GUESSES in the pseudo-code**, beyond those already flagged: the helper
names `checkedAdd`, `rewriteSubsumedIn`, `literalSpec`, `stepPoint`,
`AsGuardAtom`, and `placeholder`'s spelling are all invented — the record
fixes the behaviours and the two ordering constraints (checked add before
render; step spec before the clear pass) but names none of these.
**GUESS:** the step/clear collision check is placed after the cell walk
here; the record requires only that the spec be DERIVED before the clear
pass, so an earlier placement is equally conformant.
