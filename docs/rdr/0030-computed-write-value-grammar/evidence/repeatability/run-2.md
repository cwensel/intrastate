model: claude-fable-5-1
variant: full (profile: foundational)

# Repeatability reconstruction — run 2

Source: RDR 0030 (Counting and stepping a tag without writing out every
value by hand). Read set: C1, C2, C3, MVV, S1–S11 (all `C`/`MVV`/`S`
elements). Widened past the contract spans into `§load-bearing-decisions`
(D-identity, D-wire-byte-format, D-naming, D-selection-predicate),
`§illustrative-code`, and `§phase-2-expand-into-cells`, because the
contract text names the two signature widenings and the shim relocation
but leaves the widened record's shape, the shim's signature, and the
step-point / `expand` split under-determined. Anything still open after
that widening is marked GUESS.

## 1. Public API

The record's "module" is the loader's write-block path in
`internal/table` plus one relocated comparison core in `internal/resolve`.
Nothing new is exported to CLI consumers: no new CLI verb, no new
`Categories()` member, no new CLIError code, no new finding code, no new
`intrastate.graph/1` member (C3).

### 1a. Authoring surface (TOML, the wire format)

```toml
[rule.write]
<tag> = { step = <n> }     # n: non-zero TOML integer; exactly one key
```

- Admitted only on a tag whose declared kind is `int` with both `min`
  and `max`, or `enum` with a non-empty `domain`.
- Refused at load, category `malformed_tag_declaration`, when: any other
  key set; `n == 0`; `n` is a non-integer (TOML float `1.0` included);
  tag kind is `bool`/`set`/`scalar`; `int` missing a bound; `int` width
  not representable (`max - min < 0` or `== math.MaxInt`); `enum` domain
  contains a `#` member; `enum` domain repeats a member; the rule admits
  zero cells; the rule both steps the key and names it in `clear`; any
  admitted cell's stepped value leaves the domain (C2); the checked add
  overflows `int` (C2/A13).
- `[initial]` keeps `malformed_initial_declaration`; predicate literals
  keep `malformed_predicate_atom` (C1 interception arm).

### 1b. Go surface in `internal/resolve` (relocated, exported)

```go
// moved unchanged from internal/guard/grammar.go
type Evaluator struct{}
func (Evaluator) Evaluate(atom GuardAtom, held Value) Verdict   // GUESS: param
                                                                  // names; the
                                                                  // three-valued
                                                                  // Verdict type
                                                                  // is the one
                                                                  // guard uses
                                                                  // today

// moved from internal/guard/declaration.go, type-free core, returns the
// inclusive width; ok=false when unrepresentable.
func IntWidth(minV, maxV int) (int, bool)    // GUESS: exported spelling;
                                             // RDR names it ::intWidth and
                                             // says it relocates as-is

// extracted two-valued core of guard/product.go::valueSatisfies:
// render the atom literal, call the evaluator, pass the verdict through.
func AtomSatisfies(atom GuardAtom, held Value) Verdict   // GUESS: name.
                                                         // RDR fixes only
                                                         // that it takes a
                                                         // resolve.GuardAtom
                                                         // + held value, not
                                                         // table.Atom, and
                                                         // carries the
                                                         // canonicalizing
                                                         // set renderer
```

Error modes: the shim returns the three-valued verdict; it never refuses.
Policy (what to do with the undecided arm) stays at each caller.

### 1c. Go surface in `internal/guard` (kept, repointed)

```go
func IntDomain(d table.TagDecl) []int   // thin wrapper over resolve.IntWidth
                                        // (does not move: resolve cannot
                                        // import table)
```
`guard/assignment.go::valueAssignments`, `guard::domainSize`,
`guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`
and `cli/flow_resolve.go::guardSeam` repoint to the `resolve` copies with
unchanged behaviour.

### 1d. Go surface in `internal/table` (unexported, the two widenings)

```go
// before: (rule) ([]TagValue, []string, error)
func (l *loader) renderWrites(rule srcRule, predicates []Atom) ([]WriteSpec, []string, error)
        // GUESS: srcRule name; RDR only fixes "takes the merged predicates as a
        // parameter" and "returns a widened per-key record"

// before: expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row
func expand(base Row, predicates []Atom, outcome Atom, writes []WriteSpec) []Row
```
Both widenings are local to `normalize.go`; neither adds a call site.
The error mode of `renderWrites` is unchanged in kind: a
`*Failure{Category: "malformed_tag_declaration", Detail: …}` returned to
the loader, which surfaces it as a load refusal. No field is added to
`Failure` (C2).

## 2. Three most important internal helpers

1. **`(*loader).renderWrites` — step-spec resolution and the bound.**
   Holds the `TagDecl` and now the merged predicate set. For each written
   key: if the value is the `{ step = n }` table, validate the form and
   the tag kind (grammar arm), compute the admitted-cell list by
   evaluating every positive atom the rule authors on that tag
   (`guard.all` `eq/in/lt/lte/gt/gte`, `match` `eq/in`; never `unless`)
   per domain member through `resolve.AtomSatisfies` with the exclude
   collapse (undecided ⇒ not admitted); refuse on zero cells; for each
   admitted cell compute the stepped literal with a CHECKED add (int) or
   position index (enum), then run `conform` → `conformDomain`, refusing
   on the FIRST failure in domain order with the C2 detail. It derives
   the step spec BEFORE the clear-list pass so step∩clear is detectable.
   It places the stepped key in `assignments` with a placeholder (so
   `0002:C14`'s `RequiresOwned` still sees the key) and returns a
   `WriteSpec` per key.
2. **`expand` — the choice-point product.** Unchanged in role: takes per-key
   member lists and mints one row per combination. Gains a step choice
   point per stepped key that (a) subsumes any `match in` on the same key
   (local or context-inherited) into the admitted-cell conjunction and
   rewrites it per row to `match eq = <cell>`, exactly as the `case
   expanding:` arm does today; (b) emits the row's `guard.all eq = <cell>`
   atom; (c) writes the per-cell literal into BOTH `Row.Writes` and
   `Row.NextTags`; (d) appends the bare cell to the suffix with a
   constant-true `suffixed` local (never shares the `in` arm's
   `len(members) > 1`); (e) orders the step point among other candidates
   by KEY alone (D-identity).
3. **`resolve.AtomSatisfies` (extracted shim) + relocated `Evaluator`.**
   The one comparison the loader, `guard`, and `graphlint` all call so the
   admitted-cell filter agrees with the runtime evaluator by
   construction. Renders the atom literal (with the canonicalizing set
   renderer), evaluates, returns the three-valued verdict; each caller
   applies its own collapse. `resolve.IntWidth` rides alongside as the
   single width rule so the loader can never enumerate a domain lint
   cannot count.

(Honourable mention, because C3 makes it a real change: the graph-lint
per-row emission that renders `Rule`/`Element` as `row.Identity()`
(`rule#cell`) while `coverage.go::groupHasOverlap` joins on `Span` or the
recovered bare id — a helper such as `rowIdentityForFinding(row) string`
plus an `authoredID(identity string) string`. GUESS: names.)

## 3. Data model

### Source-side (TOML → loader)

```go
// GUESS: names. The RDR fixes only: one key, non-zero int, derived before
// the clear pass.
type stepSpec struct {
    Key string
    N   int   // non-zero
}
```

### The widened per-key write record crossing renderWrites → expand

```go
// GUESS: exact shape. The RDR fixes: it replaces []TagValue in both
// signatures, carries a non-literal, and hands expand "a list of members
// per key and the literal each cell writes".
type WriteSpec struct {
    Key     string
    Literal string        // authored literal write; empty when Cells != nil
    Cells   []cellWrite   // nil for a literal write
}
type cellWrite struct {
    Cell    string        // the admitted domain member, rendered (suffix element
                          // and guard.all eq literal; byte-identical to
                          // guard/assignment.go::valueAssignments' rendering)
    Literal string        // the stepped value rendered: cell+n (int) or the
                          // member n positions on (enum)
}
```

### Row (existing, populated per expanded row — no new fields)

```go
type Row struct {
    RuleID   string        // authored id, bare
    Span     string        // model:rule — unchanged by expansion; the stable join key
    Suffix   []string      // + one bare-cell element per stepped key (always, even at one cell)
    Atoms    []Atom        // retained authored atoms + guard.all eq=<cell>
                           // (+ match eq=<cell> where an in was subsumed); never a retained in
    Writes   []TagValue    // per-row stepped literal
    NextTags []TagValue    // per-row stepped literal, populated independently (0002:C15)
    Emit     …             // shared across a rule's rows (S7)
}
// Identity(): rule id + "#" + suffix elements joined; Fingerprint hashes Atoms before NextTags.
```

### Refusal (existing `internal/table/category.go::Failure`, no new fields)

```go
type Failure struct {
    Category  string   // "malformed_tag_declaration" for every write-block refusal here
    Detail    string   // cell FIRST, then rule, tag, result (int literal | enum offset | bound passed)
                       //   plus any unless atom on the stepped tag, noted as not excluding
    Offending string   // untouched (0008:C3 pair)
    Remedy    string   // untouched
    Rule      string   // untouched — direction id, never a rule id
    Line      int
}
```
GUESS: exact detail wording. Reconstructed from S3/S3b/S3c/S5:
`cell 9: rule retry writes attempt 10, is above max 9` /
`cell large: rule escalate steps tier 1 past the last member` /
`cell 0: rule retry steps attempt 1 before the first member` /
`cell <c>: rule r steps attempt past math.MaxInt` (A13 arm) — with an
appended `; unless <atom> did not exclude the cell` where present.

### Graph-lint finding payload (C3)

Per-row finding `rule` and `element` carry `row.Identity()` (`retry#3`);
group-level findings keep the bare authored id; `span` is the stable
key. No new code, no new field.

## 4. Top-level pseudo-code — `normalizeRule` with a step write

```
normalizeRule(src, ctx):
  base       := rowFrom(src)                                    # id, span, emit, …
  predicates := mergePredicates(ctx, src)                       # merged set, as today
  outcome    := outcomeAtom(src)
  writes, clears, err := l.renderWrites(src, predicates)        # widened (C1)
  if err != nil: return refuse(err)
  rows := expand(base, predicates, outcome, writes)             # widened (C1)
  return rows

(*loader).renderWrites(src, predicates):
  specs := []
  for key, val in src.write (authored order):
    decl := l.decls[key]
    if isStepTable(val):                                          # exactly {step=n}
      n := stepOf(val)                                            # refuse: other keys, n==0, float
      if decl.Kind not in {int with min&max, enum with non-empty domain}: refuse(malformed_tag_declaration)
      if key in src.clear: refuse(…)                              # spec derived BEFORE clear pass
      members := domainMembers(decl)                              # int: width,ok := resolve.IntWidth(min,max); !ok → refuse
                                                                  # enum: authored order; refuse on '#' member or duplicate
      positive := atomsOn(predicates, key, positiveOnly=true)     # guard.all eq/in/lt/lte/gt/gte, match eq/in; no unless
      cells := []
      for m in members:                                           # domain order
        admit := true
        for a in positive:
          if resolve.AtomSatisfies(render(a), m) != GuardTrue: admit = false   # exclude collapse
        if admit: cells.append(m)
      if len(cells) == 0: refuse(…)
      for c in cells:                                             # first failure in domain order wins (C2)
        lit, ok := stepValue(decl, c, n)                          # int: checked add; enum: index+n
        if !ok:  refuse(detailBound(c, src.id, key, offsetOrBound, unlessAtomsOn(predicates,key)))
        if f := l.conform(decl, lit); f != nil: refuse(detailBound(c, …, f.result))
        cellWrites.append({Cell: render(c), Literal: render(lit)})
      assignments[key] = placeholder                              # keeps RequiresOwned (0002:C14)
      specs.append(WriteSpec{Key: key, Cells: cellWrites})
    else:
      lit := l.conform(decl, val) …                               # today's literal path, unchanged
      assignments[key] = lit; specs.append(WriteSpec{Key: key, Literal: lit})
  clears := clearList(src)                                        # existing pass, now after step derivation
  return specs, clears, nil

expand(base, predicates, outcome, writes):
  points := []
  for w in writes where w.Cells != nil:                           # one step point per stepped key
    inAtom := matchInOn(predicates, w.Key)                        # local or inherited
    points.append(stepPoint{key: w.Key, cells: w.Cells, subsumedIn: inAtom})
  for a in predicates where a is match/guard.all in and a.Key not stepped:
    points.append(inPoint{…})                                     # existing arm, unchanged
  points.append(outcomePoint(outcome))
  sort points by Key                                              # D-identity: key alone
  for combo in product(points):                                   # 0002:C13
    row := base.clone()
    for p in combo:
      switch p:
        case stepPoint:
          row.Atoms = retained(predicates, p.key) minus p.subsumedIn
          if p.subsumedIn != nil: row.Atoms += matchEq(p.key, cell)   # rewritten, never retained
          row.Atoms += guardAllEq(p.key, cell)
          row.Writes  += TagValue{p.key, cellLiteral}
          row.NextTags += TagValue{p.key, cellLiteral}            # both carriers (0002:C15)
          row.Suffix += cell                                      # suffixed := true, always
        case inPoint:   … today's arm, suffixed := len(members) > 1
        case outcomePoint: …
    rows.append(row)
  sort rows by compareRows; return rows
```

GUESS log (silences the contract text left after widening):
- names of the extracted shim, relocated width function, widened record,
  step-spec struct, and helper for suffixing lint findings;
- exact wording/ordering of the C2 `Detail` string beyond "cell first,
  then rule, tag, result, unless-note";
- whether the `Verdict` type is `GuardTrue/GuardFalse/GuardUndecided` or
  another three-valued spelling;
- whether `renderWrites` inspects `src.clear` directly or receives the
  clear list; the RDR fixes only the order (step spec before clear pass);
- whether `int` cells and `enum` cells share one `cellWrite` string
  rendering or carry typed values — chosen string here because C1 fixes
  byte-identity with `valueAssignments`' rendered cells.
