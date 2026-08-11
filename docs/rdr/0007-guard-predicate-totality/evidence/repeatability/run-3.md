model: glm-5.2:cloud
variant: full (profile: foundational)

Model: glm-5.2:cloud

# Repeatability run-3 — RDR 0007 guard predicate totality

Reconstruction from `docs/rdr/0007-guard-predicate-totality.md` alone. Every
point the RDR leaves undetermined is marked **GUESS**. Where the RDR is silent
the silence itself is the finding — the domain rule's load-bearing decisions
are recoverable, but the evaluator's shape, the `Row.Guard` carrier, and the
atom/value representations are not fixed by this RDR.

## 1. Public API — function signatures, types, error modes

### 1.1 The constrained seam (pre-existing, RDR 0001)

The RDR constrains exactly one seam and changes no kernel data flow. It names
the seam `GuardEvaluator.Evaluate` and quotes the signature fragment
`Evaluate(guard string, view TagSet)` (via `fixtureGuards`'
`Evaluate(guard string, _ resolve.TagSet)`).

```go
package resolve // internal/resolve

// GuardEvaluator is the seam this RDR constrains. The kernel calls it; it
// does not implement operator semantics. RDR 0003's evaluator implements it.
type GuardEvaluator interface {
    // Evaluate recovers the atom structure of guard and applies the domain
    // rule. Returns a three-valued verdict; GuardUnevaluable is the honest
    // third value, never folded to false or true.
    Evaluate(guard string, view TagSet) GuardResult // GUESS: return type
}
```

**GUESS — return type.** The RDR quotes only the parameter list
`Evaluate(guard string, view TagSet)` and never the return. The seam is
verdict-only — "the verdict is the only observable" and "`GuardResult`
carries no missing-tag payload" — so a bare `GuardResult` is the most
defensible reading. Whether a Go `error` accompanies it for a malformed guard
is unspecified; the RDR pushes all parse/lint failure to RDR 0003's load-time
(`unknown tag` rejected before resolution), so at Evaluate time the guard is
well-formed and an `error` return has no defined trigger. I reconstruct a bare
`GuardResult`; a `(GuardResult, error)` shape is equally consistent with the
text.

### 1.2 The three-valued verdict (pre-existing, frozen)

```go
type GuardResult int // three-valued, already shipped: resolve.go::GuardResult
const (
    GuardFalse       GuardResult = iota // prunes the row (D8), present-tag falsity only
    GuardTrue                           // row survives; may still be refused by aggregation
    GuardUnevaluable                    // honest third value; maps to guard_unevaluable
)
```

The RDR names exactly three values and pins them to shipped names
(`GuardUnevaluable` / `guard_unevaluable` are kept deliberately; a new
`guard_input_missing` kind is rejected). `GuardUnevaluable` carries no
missing-tag payload — a named diagnostic gap handed to RDR 0003 (Phase 3).

### 1.3 The refusal (kernel consumer side, pre-existing)

```go
type Refusal struct {
    Kind  RefusalKind          // one of the closed five; guard_unevaluable here
    Guard string               // single-valued diagnostic convenience, NOT the discriminator
    Rows  []RowID             // every undecidable row — the discriminator
    MissingOwned []TagKey     // owned_state_unavailable only; absent on guard_unevaluable
}
```

**GUESS — `RowID` / `Rows` element type.** The RDR says `Refusal.Rows` "carries
every undecidable row" and `Refusal.Guard` is "the lowest row by
`(RuleID, SourceLocator)`". So a row identifier exists with at least `RuleID`
and `SourceLocator` fields and a defined ordering, but its Go type, package,
and whether it is a struct or a string handle are not stated. I reconstruct a
`RowID` struct holding `RuleID` + `SourceLocator`; the real type could be a
`Row` pointer or a stringified locator.

**GUESS — `RefusalKind` enum shape.** The RDR names `guard_unevaluable`,
`owned_state_unavailable`, `no_match`, `ambiguous_match` (four of the closed
five) and says the set is "pinned closed" by
`TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`. The fifth kind is named
only by test, not by this RDR; I do not reconstruct it.

### 1.4 `Row` surface (doc-contract narrowing only)

```go
type Row struct {
    Guard         string   // opaque guard text; kernel attaches no meaning ("" → GuardTrue pre-seam)
    RequiresOwned []TagKey // NARROWED: post-guard write-dependency keys (the keys Writes require)
    Writes        []Write  // GUESS: includes an authored clear per RDR 0002 normalization
    // ... selection fields (RuleID, SourceLocator, match pattern) — GUESS, not fixed here
}
```

`RequiresOwned` is narrowed in **doc contract only** (Phase 1 = doc comments):
it names the owned tag keys the row's *transition writes* depend on once the
guard holds — *not* guard inputs. Guard decidability is carried by the guard's
own referenced-tag set under the domain rule. Listing a guard-read owned key
in `RequiresOwned` stays legal and yields the more precise
`owned_state_unavailable` diagnosis among survivors. No `Row.Reads` field is
added (Alternative 2 is rejected).

**GUESS — `Write` / `Writes` shape.** The RDR says RDR 0002's normalization
"renders a `<clear>` write" so the kernel `Row` "carries no separate clear
list". The `Write` type and whether it carries an op/kind are RDR 0002's, not
fixed here.

### 1.5 Operator vocabulary (RDR 0003's, closed)

The evaluator must implement five operators; this RDR fixes each one's
absent-operand behavior:

```go
type Operator string
const (
    OpEq       Operator = "eq"        // value: tag == literal          — PARTIAL (absent ⇒ unevaluable)
    OpIn       Operator = "in"        // value: scalar in literal set   — PARTIAL
    OpCmp      Operator = "lt|lte|gt|gte" // value: bounded int compare — PARTIAL
    OpContains Operator = "contains"  // value: set-valued ⊇ literals  — PARTIAL (absent set ⇒ unevaluable, NOT empty set)
    OpExists   Operator = "exists"    // presence/absence              — TOTAL (sole total operator)
)
```

**GUESS — operator spellings / arity.** The RDR names operators in prose
("equality, membership, bounded integer comparison, existence, and set
containment") and in code as `eq`/`in`/`lt,lte,gt,gte`/`contains`/`exists`.
The exact token form (`eq` vs `=`), and whether `exists` takes a literal
(`exists <key>` vs `exists <key> = true`) are RDR 0003's grammar, not fixed
here. I GUESS `exists` is unary over a key with no literal, since its "whole
job is deciding presence".

### 1.6 Error modes (the verdict is the only observable)

There is no Go `error` channel for guard evaluation. The "error modes" *are*
the three verdict values, mapped by the kernel:

| Guard verdict | Kernel refusal kind | Escapable? |
| --- | --- | --- |
| `GuardTrue` | (row survives; selection/aggregation decides) | — |
| `GuardFalse` | pruned (D8); may reach `no_match` | yes (`no_match`) |
| `GuardUnevaluable` | `guard_unevaluable` | **no** (RDR 0002 closes escape lists to `no_match`, `ambiguous_match`) |

The non-escapability of `guard_unevaluable` is the user outcome the RDR exists
to secure: missing artifact state can never be masked behind an escapable
refusal class.

## 2. The three most important internal helper functions

These are the evaluator's internals (RDR 0003's future build); the RDR names
none of them, so all three are **GUESS** as to shape — only their
*responsibilities* are fixed by the Normative Contracts.

### 2.1 `recoverAtoms(guard string) → atoms` — GUESS (hinge on A10)

Recover the atom structure — operator, referenced tag key, literal, and
`all`/`unless` block placement — from the opaque `Row.Guard string`. The RDR
requires this mapping to be "total and lossless" but does **not** state it:
**A10 is Pending** and is the hinge finding. The RDR explicitly says "no Go
representation of any of that exists: `Row.Guard` is an opaque `string`,
there is no atom type, and no `testdata/` exists repo-wide". Whether `guard`
carries a parsed representation, a key into an evaluator-held structure, or
raw grammar text is unspecified and is the blocking prerequisite for every
atom-level clause. I GUESS a parse into an `atoms{all, unless}` struct; the
real carrier is whatever RDR 0003's implement stage states.

### 2.2 `evalAtom(atom, view) → verdict` — GUESS (shape); fixed (semantics)

Evaluate one atom against the view under the domain rule:
- `exists` → total: `view.Lookup(key)` present (provenance-blind) ⇒ `TRUE`,
  absent ⇒ `FALSE`.
- value-comparing (`eq`/`in`/`cmp`/`contains`) → partial: absent key ⇒
  `UNEVALUABLE` (never false, never true); an absent set-valued tag under
  `contains` ⇒ `UNEVALUABLE`, *not* the empty set.

**GUESS — `view.Lookup` return shape.** The RDR says `TagSet.Lookup` is
provenance-blind (`ok` on any key present) and is the only exported accessor.
I GUESS `Lookup(key) (TagValue, bool)`; the `TagValue` type is not fixed here.
**GUESS — set-valued representation** for `contains`; absent-as-unevaluable is
fixed, but the set model is RDR 0003's.

### 2.3 `combine(allVerdict, unlessVerdict) → GuardResult` — strong-Kleene

Combine atom verdicts under strong-Kleene (K3) three-valued logic:
conjunction = `min` under `F < U < T`; negation `¬T=F, ¬F=T, ¬U=U`. The guard
shape is `all_result ∧ ¬(unless_conj)`.

- `all` block: conjunction of its atoms; empty `all` = `TRUE` (empty conjunction).
- `unless` block: conjunction of its atoms, then negated; **empty/omitted
  `unless` = absent** (NOT a vacuously-true conjunction), contributing `TRUE`.
- `F ∧ U = F` (witnessed falsity from present tags prunes; this is what makes
  D8 sound); `T ∧ U = U`; `U ∧ U = U`.

**GUESS — A9 (empty-`unless` identity) is Pending**: the identity itself is
fixed here ("absent, not vacuous-true") but its *ownership* (this RDR vs RDR
0002's normalization) is open; if RDR 0002 already fixes it, the clause is
struck here. The semantics I implement are as written.

## 3. Data model — persisted or passed across the boundary

### 3.1 Across the kernel↔evaluator boundary

```go
// Guard input (opaque carrier) — GUESS on internal form (A10)
type Guard string   // opaque at the kernel seam; evaluator must recover atoms from it

// Evaluation view — pre-existing, assembled by resolve.go::assemble
type TagSet struct { tags map[TagKey]taggedValue }
type taggedValue struct { value TagValue; provenance Provenance } // owned|observed|recognized
// Readers: Lookup(key) (TagValue, bool) [provenance-blind, exported]
//          has(key, provenance) bool   [unexported; only missingOwned uses it]
//          matches(...); Len() int
```

Presence is monotone in the input (A6a, verified): `assemble` never drops a
key present in the input tuple; owned-over-observed shadowing overwrites a
*value*, never a key's presence. `Input.Owned []Tag` carries no error channel
(A6b, open): a truncated accessor snapshot and genuine absence are
byte-identical at the boundary.

**GUESS — `TagValue`.** The RDR references scalar values, bounded integers,
and set-valued tags but fixes no Go type. I reconstruct a sum type; the real
type is RDR 0003's.

### 3.2 Across the authoring↔`Row.Guard` boundary — GUESS (A10, the hinge)

Authored structurally as `[rule.guard.all.<tag>]` /
`[rule.guard.unless.<tag>]` (RDR 0003's spike fixture), normalized by RDR
0002 into one candidate-row predicate set, then handed to the kernel as
`Row.Guard string`. **No RDR states what that combination yields at the
`Row.Guard` boundary** (A10). The only in-tree consumer (`fixtureGuards`)
treats it as an opaque *name* keyed into `decided map[string]bool`. So the
boundary data model is: authored structured guard → (unknown normalization)
→ opaque string → (unknown recovery) → atoms. Both unknowns are GUESS.

### 3.3 The refusal payload out

```go
type Result struct {
    Plan   *Plan      // nil on refusal
    Refusal Refusal    // populated on refusal
}
func (r Result) Refused() bool

type Plan struct {
    Writes  []Write   // GUESS: includes rendered <clear>
    Escaped bool      // unreadable on a refusal
}
```

`guard_unevaluable` populates `Refusal.Kind`, `Refusal.Guard` (lowest row by
`(RuleID, SourceLocator)` — single-valued diagnostic), and `Refusal.Rows`
(every undecidable row — the discriminator). `MissingOwned` is empty on
`guard_unevaluable` (it belongs to `owned_state_unavailable`).

## 4. Top-level pseudo-code — `Evaluate` (20–40 lines)

The main operation this RDR adds is the evaluator's domain-rule `Evaluate`.
The kernel side (`gate`/`missingOwned`/`escapeOrRefuse`) is pre-existing and
unchanged; its flow is summarized at the end.

```text
Evaluate(guard, view):
    atoms := recoverAtoms(guard)              # GUESS: A10 — opaque string → atom structure
    if atoms is empty:                        # no guard / kernel already returned GuardTrue
        return GuardTrue                      #   pre-seam for guard == ""
    # all block: empty conjunction = TRUE
    allV := TRUE
    for atom in atoms.all:
        allV := k3And(allV, evalAtom(atom, view))
    # unless block: omitted/empty = ABSENT (A9) → contributes TRUE, NOT ¬(TRUE)=FALSE
    unlessV := TRUE
    if atoms.unless is not absent:
        conj := TRUE
        for atom in atoms.unless:
            conj := k3And(conj, evalAtom(atom, view))
        unlessV := k3Not(conj)                # strong-Kleene negation
    return k3And(allV, unlessV)               # GuardTrue | GuardFalse | GuardUnevaluable

evalAtom(atom, view):                         # the domain rule
    if atom.op == exists:                      # sole TOTAL operator
        _, ok := view.Lookup(atom.key)        # provenance-blind (any provenance)
        return ok ? TRUE : FALSE               # absence ⇒ decided, not unevaluable
    _, ok := view.Lookup(atom.key)            # value-comparing: PARTIAL
    if !ok:
        return UNEVALUABLE                     # never false, never true
    return compare(atom.op, val, atom.literal) ? TRUE : FALSE
    # contains over absent set handled by the !ok branch → UNEVALUABLE, not empty-set-false

# Kernel side (pre-existing, unchanged) — for context, not added by this RDR:
#   assemble(input) → view          # monotone; never drops a present key (A6a)
#   gate(rows, view):
#     survivors := [r for r in rows if Evaluate(r.Guard, view) != GuardFalse]   # D8 prune-first
#     if missingOwned(survivors, view): return owned_state_unavailable          # BEFORE undecidable
#     if any survivor GuardUnevaluable: return guard_unevaluable                # aggregation veto
#     return selected
#   escapeOrRefuse applies the same gate to the escape rows as a separate set (D5)
```

### Where this reconstruction is silent / guessed (the run's findings)

- **A10 (hinge, Pending)** — the `Row.Guard` ↔ atom-structure mapping is
  unspecified anywhere; every atom-level clause and Phase 2 vectors 4–8 hang on
  it. The RDR hands it to RDR 0003's implement stage but does not name a
  concrete destination form.
- **`Evaluate` return type** — quoted params only; bare `GuardResult`
  reconstructed, `(GuardResult, error)` equally consistent.
- **Atom / `TagValue` / `Write` Go types** — none exist; RDR 0003's.
- **`Refusal.Rows` element / `RowID`** — `RuleID`+`SourceLocator` ordering is
  fixed, the type is not.
- **A9 (Pending)** — empty-`unless` identity is fixed in *semantics* here but
  its *ownership* (this RDR vs RDR 0002) is open.
- **A6b (open, carried)** — read-failed vs genuinely-absent is
  indistinguishable at the seam; `Input.Owned` has no error channel.
- **`exists` arity / operator spellings** — RDR 0003's grammar.