model: claude-opus-5[1m]
variant: full (profile: foundational)

Model: claude-opus-5[1m]

# Repeatability run-1 — RDR 0007 guard predicate totality

Reconstruction from `docs/rdr/0007-guard-predicate-totality.md` alone. Every
point the RDR leaves undetermined is marked **GUESS**.

## 1. Public API — function signatures, types, error modes

### 1.1 The constrained seam (pre-existing, RDR 0001)

The RDR states it constrains exactly one seam and changes no kernel data flow.
It quotes the seam as `Evaluate(guard string, view TagSet)`.

```go
package resolve // internal/resolve

// GuardEvaluator is the seam this RDR binds. The kernel calls it; it does not
// implement operator semantics.
type GuardEvaluator interface {
    Evaluate(guard string, view TagSet) GuardResult   // GUESS: return type name is
                                                      // stated ("GuardResult" appears in
                                                      // Infrastructure Audit) but the
                                                      // RDR never shows the full
                                                      // signature — whether it returns
                                                      // (GuardResult, error) or a bare
                                                      // GuardResult is not stated.
}

// GuardResult is the three-valued verdict. The RDR names exactly three
// constants and states the type "carries no missing-tag payload", so it is a
// bare enum.
type GuardResult int

const (
    GuardTrue GuardResult = iota  // GUESS: ordering/zero value never stated.
    GuardFalse
    GuardUnevaluable
)
```

The RDR pins the three verdict names (`GuardTrue`, `GuardFalse`,
`GuardUnevaluable`) and the refusal kind `guard_unevaluable`, and explicitly
rejects a sixth refusal kind `guard_input_missing`. It states `GuardResult`
carries no missing-tag payload — so no `MissingTags` field.

`TagSet`: the RDR names one exported accessor, provenance-blind:

```go
type TagSet interface {   // GUESS: interface vs struct is not stated; the RDR
                          // calls TagSet.Lookup "the only exported accessor",
                          // and elsewhere writes view.tags / view.has(key,
                          // ProvenanceOwned) as unexported internals, which
                          // reads more like a struct with methods.
    Lookup(key string) (value TagValue, ok bool)   // provenance-blind ok
}
```

The RDR's normative "PRESENCE IS PROVENANCE-BLIND" clause makes `Lookup`'s `ok`
the presence predicate for guards, and contrasts it with the kernel-internal
`missingOwned`, which tests `view.has(key, ProvenanceOwned)`.

**GUESS (significant):** `TagValue`'s shape is never given. The domain rule
distinguishes *scalar* tags from *set-valued* tags (`contains` operates on "the
set-valued domain"; an absent set-valued tag "is unevaluable, not the empty
set"), and bounded integer comparison implies an integer-typed value. So
`TagValue` must carry at least a type discriminant + scalar/set payload, but the
RDR states nothing about it:

```go
type TagValue struct {   // GUESS: entire shape
    // some type tag + scalar (string/int) or set-valued payload
}
```

Provenance constants are named but not enumerated exhaustively; the RDR names
three provenances in prose ("owned, observed, or recognized") and one constant
literally:

```go
const ProvenanceOwned = ...       // named literally in the RDR
const ProvenanceObserved = ...    // named literally ("ProvenanceObserved is
                                  // documented as 'a non-owned context tag
                                  // supplied by the caller'")
const ProvenanceRecognized = ...  // GUESS: prose only ("recognized"), constant
                                  // name never given
```

### 1.2 The refusal (pre-existing, RDR 0001)

```go
type Refusal struct {
    Kind         RefusalKind   // one of a closed five-kind taxonomy
    Rows         []Row         // GUESS: element type. The RDR says
                               // "Refusal.Rows carries every undecidable row"
                               // — could be []Row, []RuleID, or []RowRef.
    Guard        string        // single-valued; the lowest row by
                               // (RuleID, SourceLocator)
    MissingOwned []string      // GUESS: named as the sibling field that "names
                               // its keys"; type never stated.
}
```

Refusal kinds: the RDR states the taxonomy is **closed at five kinds** and names
four of them across the text — `guard_unevaluable`, `owned_state_unavailable`,
`no_match`, `ambiguous_match`. **GUESS:** the fifth kind is never named in this
RDR.

`Refusal.Guard` is normatively "a diagnostic convenience, NOT the
discriminator": any contract or test distinguishing *which* row went unevaluable
MUST assert on `Rows`.

### 1.3 The one genuinely new public surface — Phase 2's conformance harness

This is the only new exported API the RDR requires, and it is deliberately
under-specified: the RDR gives the shape and defers the details.

```go
package resolveconform   // GUESS: package name — the RDR writes
                         // "resolveconform.RunGuardVectors(t, eval)" but says
                         // "the exact package name and vector encoding are
                         // Phase 2's to fix".

// RunGuardVectors runs the normative golden vector suite against any
// GuardEvaluator implementation.
func RunGuardVectors(t TestingT, eval resolve.GuardEvaluator)  // GUESS: exact
                                                               // signature.
```

Normative constraints the RDR *does* fix on this surface:
- MUST live in a **non-`_test` package** so an external caller (RDR 0003's
  implement stage) can import it.
- MUST be an **exported function taking the seam** — "caller supplies the
  evaluator", not "suite constructs one".
- Takes "a `*testing.T`-like reporter" — so a `TestingT` interface, not a
  concrete `*testing.T`. **GUESS:** the reporter interface's method set.
- Placement caveat: if RDR 0003's evaluator lands outside module
  `github.com/cwensel/intrastate`, the `internal/` placement must be revisited.

**GUESS (significant):** the *vector encoding* — file format, whether vectors
are Go literals or `testdata/` golden files, and the suite's version stamp
("named, versioned golden vector suite"). The RDR explicitly defers all of it.

### 1.4 Error modes

The RDR frames error modes as **refusals**, not Go `error`s:

| Condition | Result |
| --- | --- |
| Guard atom over an absent tag key, value-comparing operator | atom unevaluable → `GuardUnevaluable` unless another present-tag atom decides FALSE |
| Absent tag key, `exists` operator | decided TRUE/FALSE — total, never unevaluable |
| Absent set-valued tag under `contains` | unevaluable — **NOT** the empty set |
| Any surviving candidate row unevaluable | resolution refuses `guard_unevaluable` |
| Survivor set with both absent owned state and undecidable guard | `owned_state_unavailable` (pinned precedence) |
| `guard_unevaluable` | **not** a modelable escape class (RDR 0002) |

**GUESS:** whether `Evaluate` may return a Go `error` at all (e.g. for a guard
string the evaluator cannot parse). The RDR's atom vocabulary is closed and
unknown operators are rejected at parse/lint before resolution, which suggests
no runtime parse error — but the RDR never says so, and A10 (Pending) is exactly
the question of how structure reaches `Evaluate`.

## 2. The three most important internal helper functions

### 2.1 `atomVerdict(op Operator, key string, literal Literal, view TagSet) GuardResult`

**GUESS: name and signature entirely.** The RDR never names this function; it is
forced by the normative clauses. Responsibility: the domain rule at atom
granularity — look the referenced tag key up in `view` provenance-blind; if the
operator is `exists`, decide TRUE/FALSE from `ok` alone; for the four
value-comparing operators, return `GuardUnevaluable` when `!ok`, else compare.

### 2.2 `combineKleene(all []GuardResult, unless []GuardResult) GuardResult`

**GUESS: name and signature entirely.** Responsibility: strong-Kleene
three-valued combination across the `all` and `unless` blocks, implementing
`all ∧ ¬(unless_conj)`. Load-bearing behaviors the RDR fixes:
- `F ∧ U = F` — a present-tag atom decided FALSE yields `GuardFalse` for the
  whole `all` block even beside an unevaluable atom. This is what makes D8 sound.
- A verdict is TRUE/FALSE only when present tags **alone** decide it; any
  *unresolved* dependence on an unevaluable atom → `GuardUnevaluable`.
- Empty `all` block = TRUE (empty conjunction).
- Empty `unless` block = **absent**, NOT a vacuously-true conjunction.
- No guard at all (empty `Row.Guard`) → `GuardTrue` before the seam is reached,
  by `resolve.go::evaluateGuard`. Outside the domain rule's scope entirely.

**GUESS:** the exact shape of the `unless` negation under three-valued logic —
the RDR states `¬` is applied to the `unless` conjunction and that strong Kleene
governs, so `¬U = U`, but never writes the negation table.

### 2.3 `gate(rows []Row, view TagSet) (survivors []Row, refusal *Refusal)`

Named literally in the RDR (`resolve.go::gate`), along with its two collaborators
`missingOwned` and `assemble`, and its normative ordering:

1. Prune rows whose guard is decided `GuardFalse` (D8 — "guard-FALSE prunes
   first; a pruned row contributes neither candidacy nor an owned-state
   obligation").
2. Among survivors, run `missingOwned` → refuse `owned_state_unavailable`.
3. Then the undecidable loop → refuse `guard_unevaluable`.

The ordering is normatively pinned **to shipped behavior**, not to its original
rationale, which this RDR's `RequiresOwned` narrowing invalidates. Phase 2 must
assert the ordering as behavior.

**GUESS:** `gate`'s actual signature — the RDR names the function and its
ordering but never its parameters or returns.

Honorable mention (named, not reconstructed): `assemble` (builds the view;
writes `Input.Observed` into `view.tags` unconditionally), `escapeOrRefuse`
(applies the same gate to escape rows as a **separate row set**, per D5),
`evaluateGuard` (special-cases `""` → `GuardTrue`), `missingOwned` (tests
`view.has(key, ProvenanceOwned)` — owned-provenance only).

## 3. Data model crossing the boundary

### 3.1 `Row` (RDR 0001's shape — this RDR narrows one field's *doc contract* only)

```go
type Row struct {
    RuleID        ...   // GUESS: type. Used with SourceLocator as the
                        // tie-break key for Refusal.Guard selection.
    SourceLocator ...   // GUESS: type.
    Guard         string  // OPAQUE at the kernel seam. Stated normatively.
    RequiresOwned []string // GUESS: element type; the RDR calls them "owned tag
                           // keys", so []string is the natural read.
    Writes        ...   // GUESS: type. Includes an authored clear — RDR 0002
                        // normalization "renders a <clear> write", so Row
                        // carries NO separate clear list.
}
```

**The normative narrowing:** `Row.RequiresOwned` names the owned tag keys the
row's **post-guard transition** depends on — the keys its `Writes` require. It
is explicitly **not** about guard decidability. Listing a guard-read owned key
there remains legal and yields the more precise `owned_state_unavailable`
diagnosis among survivors.

### 3.2 The guard's atom structure (RDR 0003's model — **the largest silence**)

The RDR quantifies its entire normative block over guard **atoms** — operator,
referenced tag key, literal, `all`/`unless` placement — but states that this
structure is RDR 0003's to define, and that the mapping from authored guard to
`Row.Guard string` "MUST be stated before that evaluator is built" (A10,
Pending). All the RDR fixes here is that the mapping be **total and lossless**
for those four components.

```go
// GUESS: this entire model. The RDR requires it to exist and be recoverable,
// and explicitly declines to define it.
type Guard struct {
    All    []Atom
    Unless []Atom
}

type Atom struct {
    Op      Operator
    TagKey  string
    Literal Literal
}
```

Operator vocabulary — **closed and typed**, five operators, no runtime
extension point (RDR 0003):

| Operator | Kind | Totality under this RDR |
| --- | --- | --- |
| equality (`eq`) | scalar vs typed literal | partial |
| membership (`in`) | scalar vs typed literal set | partial |
| bounded integer comparison (`lt`/`lte`/`gt`/`gte`) | bounded integer domain | partial |
| set containment (`contains`) | set-valued domain | partial |
| existence (`exists`) | optional scalar or optional set-valued tag | **total** |

**GUESS:** the concrete Go spelling of every operator constant; the RDR quotes
authored TOML forms (`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`) but
gives no in-memory representation.

### 3.3 `Input` / the assembled view

```go
type Input struct {
    Observed ...   // caller-supplied; assemble writes it into view.tags
                   // unconditionally. GUESS: type.
    // GUESS: the owned-snapshot and recognized fields — the RDR names the
    // provenances but not the Input fields carrying them.
}
```

A13 (Pending) turns on exactly this: `Input.Observed` is caller-supplied and
presence is provenance-blind, so a caller can convert a `guard_unevaluable` into
a decided verdict.

## 4. Top-level pseudo-code of the main operation

```
resolve(input) -> Plan | Refusal:

  view := assemble(input)              # owned snapshot + observed + recognized
                                       # observed written in unconditionally

  candidates, refusal := gate(view, table.candidateRows)
  if refusal != nil:
      # candidate set could not be gated: owned_state_unavailable or
      # guard_unevaluable. D5: escape rows are a SEPARATE set and must not
      # mask an unevaluable candidate.
      return refusal

  if len(candidates) == 1:  return plan(candidates[0])
  if len(candidates) == 0:
      escapes, escRefusal := gate(view, table.escapeRows)   # separate row set
      if escRefusal is guard_unevaluable:
          return escRefusal        # NORMATIVE: must not claim the escape
                                   # failed to rescue when it could not be
                                   # decided at all
      if len(escapes) == 1: return plan(escapes[0])
      return refusal{no_match}     # GUESS: exact escape-exhaustion mapping
  return refusal{ambiguous_match}  # GUESS: exact-one selection is stated;
                                   # the >1 mapping to ambiguous_match is
                                   # inferred from the closed escape classes.


gate(view, rows) -> (survivors, refusal):

  survivors    := []
  undecidable  := []

  for row in rows:
      if row.Guard == "":                 # no atoms: outside the domain rule
          verdict := GuardTrue            # evaluateGuard short-circuits
      else:
          verdict := evaluator.Evaluate(row.Guard, view)

      switch verdict:
      case GuardFalse:        continue                    # D8: prune first;
                                                          # contributes neither
                                                          # candidacy nor an
                                                          # owned obligation
      case GuardTrue:         survivors.append(row)
      case GuardUnevaluable:  survivors.append(row)       # GUESS: whether an
                              undecidable.append(row)     # unevaluable row also
                                                          # joins `survivors` for
                                                          # the missingOwned scan,
                                                          # or only `undecidable`.
                                                          # The RDR says the
                                                          # owned check runs
                                                          # "among surviving rows"
                                                          # and that a survivor set
                                                          # carrying both must
                                                          # refuse owned_state_
                                                          # unavailable — which
                                                          # only bites if it does.

  # ORDER IS NORMATIVE — pinned to shipped behavior.
  missing := missingOwned(survivors, view)          # ProvenanceOwned only
  if missing not empty:
      return nil, Refusal{owned_state_unavailable, MissingOwned: missing}

  # Resolution-level aggregation veto:
  if undecidable not empty:
      return nil, Refusal{
          guard_unevaluable,
          Rows:  undecidable,                       # THE discriminator
          Guard: lowest(undecidable, by (RuleID, SourceLocator)).Guard,
      }                                             # single-valued; diagnostic
                                                    # convenience only

  return survivors, nil


Evaluate(guard, view) -> GuardResult:            # RDR 0003 builds this

  g := recoverStructure(guard)   # A10, PENDING — the RDR does not say how.
                                 # Must be total and lossless for operator,
                                 # tag key, literal, block placement.

  allV := TRUE                                   # empty conjunction
  for atom in g.All:      allV = kleeneAnd(allV, atomVerdict(atom, view))

  if g.Unless is absent:  unlessV := FALSE       # empty unless is ABSENT,
                                                 # not vacuously true
  else:
      unlessV := TRUE
      for atom in g.Unless: unlessV = kleeneAnd(unlessV, atomVerdict(atom, view))

  return kleeneAnd(allV, kleeneNot(unlessV))     # F ∧ U = F preserved


atomVerdict(atom, view) -> GuardResult:
  value, present := view.Lookup(atom.TagKey)     # PROVENANCE-BLIND
  if atom.Op == exists:
      return boolVerdict(present)                # sole TOTAL operator
  if !present:
      return GuardUnevaluable                    # partial: never false, never
                                                 # true. Absent set-valued tag
                                                 # under `contains` is
                                                 # unevaluable, NOT empty set.
  return boolVerdict(compare(atom.Op, value, atom.Literal))
```

## Guess inventory (the RDR's silences, ranked)

1. **Guard atom structure at the seam (A10, Pending)** — the entire normative
   block quantifies over a structure the RDR explicitly does not define and
   routes to RDR 0003. Everything in §2.1, §2.2, and §3.2 is reconstruction.
2. **`TagValue` / tag typing** — scalar vs set-valued vs bounded-integer is
   load-bearing for three of the five operators (`contains` absence semantics is
   *normative*), and the type is never given.
3. **Whether an unevaluable row joins `survivors`** for the `missingOwned` scan.
   The pinned owned-before-unevaluable precedence is only observable if it does,
   yet the RDR describes prune/survive in terms of `GuardFalse` only.
4. **`Evaluate`'s full signature** — bare `GuardResult` vs `(GuardResult, error)`.
5. **The fifth refusal kind** — taxonomy stated closed at five; four named.
6. **Conformance harness package name, reporter interface, vector encoding** —
   deliberately deferred to Phase 2, but three of the RDR's own constraints
   (importability, exported-function shape, caller-supplies-evaluator) are
   normative while the encoding is not.
7. **The `¬` table for `unless` under strong Kleene** — stated by reference,
   never written out.
8. **Ambiguity/escape-exhaustion mapping** to `no_match` / `ambiguous_match` —
   inferred from the closed escape-class list, not stated here.
9. **`Refusal.Rows` element type** and `Row`'s field types generally.
10. **`ProvenanceRecognized`** constant name — prose only.
