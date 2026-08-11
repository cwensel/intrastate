model: claude-fable-5
variant: full (profile: foundational)

Model: claude-fable-5

# Repeatability run-2 — RDR 0007 guard predicate totality

Reconstruction from `docs/rdr/0007-guard-predicate-totality.md` alone. Every
point the RDR leaves undetermined is marked **GUESS**.

## 1. Public API — function signatures, types, error modes

### 1.1 The constrained seam (pre-existing, RDR 0001)

The RDR constrains exactly one seam and changes no kernel data flow. It names
the seam `GuardEvaluator.Evaluate` and quotes the signature fragment
`Evaluate(guard string, view TagSet)` (via `fixtureGuards`'s
`Evaluate(guard string, _ resolve.TagSet)`).

```go
package resolve // internal/resolve

// GuardEvaluator is the seam. "The kernel calls it; it does not
// implement operator semantics."
type GuardEvaluator interface {
    Evaluate(guard string, view TagSet) GuardResult
}
```

- **GUESS**: `Evaluate` returns `GuardResult` alone, with no `error` second
  return. The RDR never shows the return side of the signature. It strongly
  implies the verdict is the whole answer — `GuardResult` is already
  three-valued and `GuardUnevaluable` *is* the failure channel ("the verdict
  is the only observable"), and `evaluateGuard` answers the nil-seam case with
  `GuardUnevaluable`, not an error. But the RDR never rules out
  `(GuardResult, error)`.
- **GUESS**: `GuardEvaluator` is an interface (single-method). The RDR calls
  it a seam the kernel calls and RDR 0003 implements, and shows a test stub
  standing in for it, which fits an interface; a function type
  `type GuardEvaluator func(string, TagSet) GuardResult` would also satisfy
  every quoted sentence.

```go
// Three-valued verdict enum — names are shipped and kept by the
// Load-Bearing Decisions "Naming" entry.
type GuardResult int // GUESS: underlying type int (could be a string enum)

const (
    GuardFalse GuardResult = iota // GUESS: ordering/values of the constants
    GuardTrue
    GuardUnevaluable
)
```

### 1.2 The kernel entry point (pre-existing, RDR 0001)

```go
func Resolve(in Input /* GUESS: , rows []Row or a Table, and the evaluator */) Result
```

- The RDR shows `Resolve` exists (`Resolve` returns on `gate(candidates,…)`'s
  block before `escapeOrRefuse` is reachable) and that `refuse(in, *blocked)`
  is an internal constructor. **GUESS**: the full parameter list — whether the
  table/rows and the `GuardEvaluator` arrive as arguments, as fields of
  `Input`, or on a receiver. The RDR never shows it.
- Error mode: refusal is in-band. `Result.Refused()` is true and
  `Result.Plan` is nil on refusal; the RDR states `Escaped` is a `Plan` field
  "unreadable on a refusal", so the discriminator is `Refusal.Kind` plus a
  nil `Plan`. **GUESS**: `Result` also has no `error` channel — the closed
  refusal taxonomy appears to be the only failure surface at this seam.

### 1.3 Refusal kinds

The five-kind set is closed and pinned by
`TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`. The RDR names four
members and their triggers (RDR 0001's "zero, multiple, unavailable, or
unevaluable"):

| kind | trigger | escapable (RDR 0002) |
| --- | --- | --- |
| `no_match` | zero candidates after guard evaluation | yes |
| `ambiguous_match` | multiple matching candidates | yes |
| `owned_state_unavailable` | a survivor's `RequiresOwned` key absent as owned | no |
| `guard_unevaluable` | a survivor's guard undecidable | no |

- **GUESS**: the fifth kind's name and trigger. The RDR quotes the taxonomy
  as closed at five (`KindGuardUnevaluable` among them, constant-name prefix
  `Kind*`) but only ever discusses these four. Guess: something like an
  input/table-validity refusal (`invalid_input` or similar), never reachable
  from the guard seam.
- **GUESS**: Go constant spellings `KindNoMatch`, `KindAmbiguousMatch`,
  `KindOwnedStateUnavailable` mirroring `KindGuardUnevaluable`; wire/string
  forms are the snake_case names quoted throughout.

### 1.4 The conformance harness (Phase 2 — this RDR's new public surface)

```go
package resolveconform // GUESS: package name — RDR gives
                       // "resolveconform.RunGuardVectors(t, eval)" as "e.g."
                       // and says the exact package name and vector encoding
                       // are Phase 2's to fix

// RunGuardVectors runs the normative golden vector suite against any
// GuardEvaluator. Shape is fixed by the RDR: "caller supplies the
// evaluator", not "suite constructs one"; must live in a non-_test
// package so RDR 0003's build can import it.
func RunGuardVectors(t TestingT, eval resolve.GuardEvaluator)
```

- The RDR fixes: exported, importable (non-`_test` package under the same
  module — `internal/` is importable anywhere inside
  `github.com/newcoinc/intrastate`), takes "the seam and a `*testing.T`-like
  reporter".
- **GUESS**: `TestingT` is a local minimal interface
  (`Errorf/Fatalf/Helper`-shaped) rather than a concrete `*testing.T`, since
  the RDR says "`*testing.T`-like reporter". A concrete `*testing.T` would
  also satisfy the sentence.
- **GUESS**: whether the vectors are compiled-in Go values or golden files
  under a new `testdata/` (the RDR notes "no `testdata/` exists repo-wide"
  and calls the suite "named, versioned golden vectors" — I guess golden
  files with an in-package loader, but structs-in-Go is equally consistent).

### 1.5 The view-reading evaluator (Phase 2 — new, not reuse)

Phase 2 must build "a view-reading evaluator … the domain rule's first
executable expression" because `fixtureGuards` discards its `TagSet`.

- **GUESS**: its exported shape and package. The RDR never says whether it is
  exported at all, or a test-tree-internal reference implementation the
  vectors run against. Guess: unexported reference evaluator beside the
  harness, since the *normative* evaluator is RDR 0003's future build and
  exporting a second one would invite production use of a stopgap.
- **GUESS**: how it recovers atom structure, since A10 (the
  `Row.Guard string` ↔ atom-structure mapping) is open by design. To be
  buildable at all it must sidestep A10 — guess: it defines a private literal
  guard encoding for vector purposes only (or is constructed with an explicit
  atom list per vector, keyed by guard text), without claiming that encoding
  for RDR 0003.

## 2. The three most important internal helpers

Chosen for the module this RDR actually specifies (the evaluator obligation +
harness); the kernel's own helpers (`gate`, `assemble`, `missingOwned`,
`escapeOrRefuse`, `evaluateGuard`) pre-exist, are frozen, and are consumed
unchanged.

1. **`recoverAtoms(guard string) (atoms, ok)`** — GUESS (name and existence
   as a distinct function). Responsibility: produce, for the guard the kernel
   handed over, its atom structure — per atom: operator, referenced tag key,
   typed literal, and `all`/`unless` block placement. The RDR makes this the
   hinge: "an evaluator that cannot recover which tag keys an atom references
   cannot implement the domain rule at all", and requires the mapping be
   total and lossless — but hands the mapping itself to RDR 0003 (A10). So
   the helper's contract is fixed here while its implementation is
   unspecifiable; in Phase 2 it is backed by the vector encoding
   (§1.5 GUESS).

2. **`evalAtom(atom, view TagSet) verdict3`** — GUESS (name). Responsibility:
   the per-atom domain rule. Look the referenced key up
   **provenance-blind** (`TagSet.Lookup`'s `ok` — never `has(key,
   ProvenanceOwned)`). If the operator is `exists`: decide TRUE/FALSE from
   presence alone (the sole total operator). Otherwise (equality, membership,
   bounded integer comparison, set containment): absent key ⇒ `U`
   (unevaluable), never false, never true; present key ⇒ compare the typed
   value against the literal, yielding `T`/`F`. `contains` over an absent
   set-valued tag is `U`, not empty-set-false.

3. **`combineKleene(allVerdicts, unlessVerdicts) verdict3`** — GUESS (name
   and factoring). Responsibility: strong-Kleene aggregation into the row
   verdict, shape `all_result ∧ ¬(unless_conj)` (the shape is RDR 0003's;
   the K3 lifting is this RDR's). Conjunction is `min` under `F < U < T`;
   `¬U = U`. Empty `all` block ⇒ `T` (empty conjunction); an omitted/empty
   `unless` block is treated as **absent** — the `¬(unless_conj)` term is not
   formed — never as a vacuously-true conjunction (which would disable every
   row carrying one). `F ∧ U = F` stands (witnessed falsity; what makes D8
   sound). Map `T/F/U` → `GuardTrue`/`GuardFalse`/`GuardUnevaluable`.

Kernel-side, the load-bearing pre-existing helper is **`gate(rows, view)`**
(GUESS at exact parameters): prune decided-FALSE rows, run
`missingOwned(survivors, view)` (owned-state refusal reported before
undecidable guards), then the undecidable loop — `len(undecidable) > 0`
overwrites the named return `selected` with nil, which is the
resolution-level veto.

## 3. Data model at the boundary

Nothing new is persisted by this RDR except the Phase 2 vectors (encoding =
GUESS, §1.4). Passed across the seam:

```go
// The assembled evaluation view. Built by assemble; read-only thereafter.
type TagSet struct {
    tags map[string]taggedValue // keyed by tag key; GUESS: field names
}

func (s TagSet) Lookup(key string) (Value, bool) // GUESS at exact returns —
                                                 // RDR fixes only: exported,
                                                 // provenance-blind, "ok on
                                                 // any key present"
// unexported: has(key, provenance) — owned-scoped, sole caller missingOwned
// unexported: matches(want …) — candidate selection; Len() int

type taggedValue struct { // GUESS: name and exact fields
    value      Value      // GUESS: representation of typed values
    provenance Provenance
}

type Provenance int // GUESS: underlying type; constants ProvenanceOwned,
                    // ProvenanceObserved, ProvenanceRecognized (first two
                    // quoted; third implied by "owned, observed, or
                    // recognized")
```

```go
type Input struct {
    Owned      []Tag  // no error channel — a truncated snapshot and genuine
                      // absence are byte-identical (A6b, open by design)
    Observed   []Tag  // GUESS: []Tag — caller-supplied context tags
    Recognized string // gates an *insert* in assemble ("if in.Recognized != ''")
    // GUESS: any further fields (table reference? evaluator?) — unstated
}

type Tag struct { Key string; Value Value } // GUESS: whole shape — never shown
```

```go
type Row struct {
    RuleID        string   // GUESS: type — used in the ordered pair
    SourceLocator string   // GUESS: type — (RuleID, SourceLocator) orders rows
    Guard         string   // opaque to the kernel; "" = unguarded ⇒ GuardTrue
                           // before the seam is consulted
    RequiresOwned []string // NARROWED MEANING (this RDR): the owned tag keys
                           // the row's post-guard transition writes depend on
                           // (incl. rendered <clear> writes) — NOT guard inputs
    Writes        []Tag    // GUESS: shape; per RDR 0002 normalization renders
                           // a <clear> write, no separate clear list
    // GUESS: match-pattern fields (outcome/tag pattern) and escape marking —
    // the RDR references "matched the outcome and the tag pattern" and
    // "modeled escape edge" without showing the fields
}
```

```go
type Refusal struct {
    Kind         RefusalKind // GUESS: type name
    Rows         []RowRef    // every undecidable row — THE discriminator;
                             // GUESS: element type (row ids? *Row?)
    Guard        string      // single-valued diagnostic convenience: the
                             // lowest undecidable row by (RuleID,
                             // SourceLocator); MUST NOT be asserted on to
                             // distinguish rows
    MissingOwned []string    // owned_state_unavailable payload, e.g. [gate]
    // NOTE: no missing-tag payload for guard_unevaluable — the absent tag is
    // inferable only from guard text; enrichment is RDR 0003 Phase 3's.
}

type Result struct {
    Plan    *Plan // nil on refusal — half the refusal discriminator
    Refusal *Refusal // GUESS: representation; Refused() bool is quoted
}
func (r Result) Refused() bool

type Plan struct {
    Escaped bool // unreadable on refusal (Plan is nil)
    Writes  []Tag // GUESS: field shape
}
```

**Vector record (Phase 2, persisted)** — GUESS at encoding and fields;
constrained by the RDR to cover: operator × {present, absent} ×
strong-Kleene verdict pairs × `unless` atoms × row aggregation, versioned,
named. Guess at shape: `{name, guard structure (per the §1.5 vector-local
encoding), view contents incl. provenance, expected GuardResult}` plus
aggregation-level vectors `{rows, view, expected Refusal.Kind + Rows}`.

## 4. Top-level pseudo-code — the main operation

Guard evaluation under the domain rule (the seam this RDR binds), in its
kernel context:

```
Evaluate(guard, view):                          # RDR 0003's evaluator MUST behave so
  # kernel already answered guard == "" with GuardTrue before this seam
  atoms ← recoverAtoms(guard)                   # per A10 mapping (open; total+lossless required)
  if atoms unrecoverable:                       # GUESS: this case exists at all —
      return GuardUnevaluable                   #   nothing may silently decide, so U

  allV ← T                                      # empty all-block ⇒ T (empty conjunction)
  for atom in atoms.all:
      allV ← min(allV, evalAtom(atom, view))    # strong-Kleene ∧ over F < U < T

  if atoms has no unless block:                 # omitted/empty unless ⇒ ABSENT,
      rowV ← allV                               #   never a vacuous-true conjunction
  else:
      unlessV ← T
      for atom in atoms.unless:
          unlessV ← min(unlessV, evalAtom(atom, view))
      rowV ← min(allV, ¬unlessV)                # ¬T=F ¬F=T ¬U=U

  return {T: GuardTrue, F: GuardFalse, U: GuardUnevaluable}[rowV]

evalAtom(atom, view):
  (val, present) ← view.Lookup(atom.key)        # provenance-blind, always
  if atom.op == exists:
      return bool3(present == atom.wantPresent) # sole TOTAL operator; absence decides
  if not present:
      return U                                  # partial: never false, never true
                                                #   (contains: absent ≠ empty set)
  return bool3(compare(atom.op, val, atom.literal))   # typed comparison ⇒ T | F

# consumer side, unchanged (frozen; shown for the contract it realizes):
gate(rows, view):
  survivors ← rows minus decided-GuardFalse rows          # D8: prune first
  if missingOwned(survivors, view) ≠ ∅: refuse owned_state_unavailable  # before undecidable
  if any survivor is GuardUnevaluable: refuse guard_unevaluable         # veto: discard
      with Rows = all undecidable survivors               #   accumulated GuardTrue rows
  exactly one GuardTrue survivor → selected; zero → no_match; multiple → ambiguous_match
  # escapes: same gate on the escape rows as their OWN set (D5); unevaluable
  # escape row ⇒ guard_unevaluable REPLACES the candidate-set refusal;
  # only no_match / ambiguous_match are escapable at all
```

- **GUESS** (flagged in the pseudo-code): whether "guard text that cannot be
  mapped to atoms" is a reachable case for a conforming evaluator, and if so
  that it yields `GuardUnevaluable` rather than a panic/error — the RDR is
  silent because A10 defers the mapping wholesale.
- **GUESS**: relative order of `no_match`/`ambiguous_match` determination vs
  the two refusal checks shown — the RDR fixes only: prune first (D8),
  owned-before-undecidable among survivors, and the undecidable veto
  discarding `selected`. Where the zero/multiple counts are computed relative
  to those is unstated (it is unobservable when the refusal checks fire
  first, so I ordered them as `gate`'s quoted control flow suggests).
- **GUESS**: `bool3`/`min` framing — the RDR fixes the truth tables (A2) but
  not that the implementation literalizes them as a min-lattice; any
  encoding with those tables conforms.
