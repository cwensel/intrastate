model: glm-5.2:cloud
variant: full (profile: foundational)

# Repeatability run-3 — RDR 0009 reconstruction

Reconstructed from `docs/rdr/0009-escape-row-shape-conformance-ownership.md`
**based only on that RDR**. Every place the RDR is silent is flagged
`GUESS`. The RDR pins identifier spellings (`ErrEscapeShapeBreach`,
`*EscapeShapeBreachError`, `Table.CheckValid`, `Code:
"escape-row-shape-breach"`, `GroupInternal`); everything else about shape,
ordering, message text, field names, and file location is inferred and
marked.

---

## 1. Public API (function signatures, types, error modes)

The change lives in the kernel package `internal/resolve` (the RDR 0001
resolution kernel) and adds one CLI-side wrapping obligation. The RDR
names `internal/resolve` as the owner of `Row`/`Plan`/`Resolve`; I assume
the new exported symbols are declared in that package. **GUESS: exact
file** (`resolve.go` vs a new `breach.go`) — the RDR does not name one.

### Kernel exported surface (what an implementer must reproduce)

```go
package resolve

// Existing, unchanged in shape:
type Tag struct { Key, Value string }   // GUESS: Tag field names not restated in 0009

type RowRef struct { RuleID, SourceLocator string }   // existing; ordered by compareRefs

type Row struct {
    Escape        []EscapeEdge // GUESS: element type name; RDR only says "Escape list"
    Writes        []Tag
    NextTags      []Tag        // GUESS: element type; RDR calls NextTags "the next state"
    RuleID        string
    SourceLocator string
} // single shape — NO ordinary/escape split, NO row-kind field (RDR: "shared resolve.Row type keeps its single shape")

type Table struct { Rows []Row }   // GUESS: whether Table carries more than Rows

type Plan struct {
    Writes        []Tag
    NextTags      []Tag
    Escaped       bool
    RuleID        string
    SourceLocator string
} // GUESS: full Plan field set; RDR names Writes, NextTags, Escaped, RuleID, SourceLocator

type Refusal struct { Kind RefusalKind; Rows []RowRef }
type RefusalKind /* exactly five, closed — unchanged */ // GUESS: the five names are not enumerated in 0009

type Result struct { Plan *Plan; Refusal *Refusal }
func (r Result) Refused() bool   // true iff Refusal != nil

type Input /* GUESS: field set */ // RDR says `assemble(in)`, `in.Table`; I assume Input embeds/holds Table + the input tuple (observed/recognized tags).

func Resolve(in Input) (Result, error)   // existing signature; error return was always-nil at HEAD, now non-nil on breach

// NEW — exported by this RDR (spellings pinned by Normative Contracts):
var ErrEscapeShapeBreach error = errors.New("escape-row shape breach") // GUESS: message text

type EscapeShapeBreachError struct {
    Ref   RowRef
    Count int
}
func (e *EscapeShapeBreachError) Error() string   // GUESS: exact wording (see §2)
func (e *EscapeShapeBreachError) Unwrap() error { return ErrEscapeShapeBreach }

func (t Table) CheckValid() error   // nil if conforming; else the errors.Join aggregate of *EscapeShapeBreachError
```

### Error modes

- **Conforming table** — `Resolve` returns `(Result{…one disposition…}, nil)`.
  `Plan` XOR `Refusal` populated exactly as before; the added scan is
  allocation-free in the conforming case (reads `len(Escape)`/`len(Writes)`
  only).
- **Breaching table** — `Resolve` returns `(Result{}, non-nil error)` where
  `Result{}` is the zero value: `Plan == nil`, `Refusal == nil`, hence
  `Refused() == false`. Callers MUST check the error before branching on
  `Refused()`, or they read a success-shaped value with a nil `Plan`
  (scenario 7b trap). The error is `Table.CheckValid()`'s return **verbatim** —
  never `fmt.Errorf`-wrapped, because an added layer would expose
  `Unwrap() error` at the outermost level and break the flat
  `Unwrap() []error` one-level guarantee.
- **The breach is NOT a modeled refusal** — never a `RefusalKind`, never a
  sixth taxonomy member. It travels the Go error path RDR 0001 reserves for
  programmer mistakes. The five-kind refusal taxonomy stays closed.
- **Multi-breach** — one `Resolve` call reports **every** breaching identity
  in one pass (aggregate via `errors.Join`, not fail-fast). Equal
  `RowRef` identities collapse to one reported `*EscapeShapeBreachError`
  whose `Count` is the pre-collapse row count. `Unwrap() []error` of the
  aggregate is exactly one level deep; every element is a
  `*EscapeShapeBreachError`; `errors.Is(elem, ErrEscapeShapeBreach)` holds
  on each. A single `errors.As`/`AsType` finds only the FIRST breach, so a
  test reading all offending rows must traverse the slice.
- **Authored path** (separate enforcement, owned by RDR 0002's normalizer,
  not by this kernel code): an authored escape rule carrying a write block
  or clear list — including an **empty** block `writes = []` — is rejected
  at table load under RDR 0002's "malformed escape declaration" category,
  naming the source rule. Keys on **presence** of the block, not contents;
  normalized escape rows render write-free. This is stricter than the
  kernel's length-based predicate by design.

### CLI surfacing (a future `flow` verb; none exists at HEAD)

```go
// GUESS: helper name + location. The RDR binds the carrier CLASS (a new
// omitempty field on *clierr.CLIError, NOT Detail, NOT Cause) but leaves
// the field name and clierr-local type to A9 (Pending).
func wrapEscapeBreach(err error) *clierr.CLIError {
    return &clierr.CLIError{
        Code:  "escape-row-shape-breach",
        Group: clierr.GroupInternal,          // exit 2
        Hint:  "fix the table producer: an escape row must carry no writes", // GUESS: exact wording (RDR gives this example)
        Detail: <human sentence>,             // GUESS: prose; identities do NOT live here
        <NewField>: <clierr-local row-identity repr>, // GUESS: name + type — A9; clierr MUST NOT import internal/resolve
    }
}
```
`CLIError.Cause` is `json:"-"`, so the `RowRef` values riding the Go error
chain are not wire-visible; identities MUST reach the envelope through the
new serialized `omitempty` field. Returning the kernel error **unwrapped**
is a defect the exit code does NOT reveal (still exit 2, but misclassified
as `command-error`/`GroupUserEnv` via `cobraErrorToCLIError`), so
conformance is asserted on `Code`, never on exit code alone. This is
**deferred** (scenario 11): no `flow` verb exists at HEAD, and A9 is
Pending.

---

## 2. The three most important internal helper functions

### (a) `Table.CheckValid` — the conformance scan + aggregate build

Responsibility: the single predicate enforced at two call sites —
`Resolve` entry (the backstop) and any producer at construction time — so
the two enforcement points cannot drift. Iterates the **whole** table
(`Table.Rows`), applies the per-row breach predicate
`len(row.Escape) != 0 && len(row.Writes) != 0` (length-based, **not**
nil-ness; does **not** widen to `NextTags`), and builds the aggregate. A
breach surfaces even when no resolution path reaches the row (whole-table
scope, value-level, independent of the input tuple) and precedes every
modeled disposition.

The RDR fixes the *outcome* (every breaching identity reported, ordered by
`compareRefs`, equal identities collapsed, `Count` = pre-collapse rows,
joined with `errors.Join`) but not the mechanism. **GUESS: the collapse is
a map keyed on `RowRef` accumulating counts, with a separate ordered
distinct-identity slice, then `slices.SortFunc(order, compareRefs)`
(non-stable, matching the existing `rowRefs`), then `errors.Join`.** An
equivalent sort-then-adjacent-dedupe is not ruled out.

### (b) `*EscapeShapeBreachError` — per-identity typed carrier

Responsibility: carry exactly one `RowRef` (`Ref`) and the per-identity
pre-collapse `Count`, and wrap the package-level `ErrEscapeShapeBreach`
sentinel so `errors.Is` classifies **every** per-row element (not only the
aggregate). `Unwrap() error` returns `ErrEscapeShapeBreach`. Both the type
and the sentinel are **exported** — intended callers (a future `flow` verb,
non-kernel table producers) are outside the package, so an unexported
sentinel would defeat the `errors.Is` classification this clause exists
to provide. Row identity is recoverable structurally (`errors.As`/`AsType`
on `Ref`), never only from formatted prose. **GUESS: `Error()` message
wording** — something like
`"escape-row shape breach: rule %q at %s carries %d owned-tag write(s) on an escape row"`
— the RDR forbids relying on the text, so the exact string is unspecified.

### (c) CLI breach wrapper (name GUESS: `wrapEscapeBreach`)

Responsibility: translate the kernel's typed error into the wire surface
an operator reads. Produces `*clierr.CLIError{Code: "escape-row-shape-breach",
Group: GroupInternal (exit 2), Hint: <remedy>, <new omitempty field>: <identities>}`,
extending the shipped `config.Load` wrapping pattern (`Code + Detail + Group
+ Cause`) by adding a `Hint` and the identity carrier. Because `Cause` is
`json:"-"`, the wrapper renders the offending `RowRef` identities into the
new serialized field — **not** into `Detail`, which would force consumers
to re-parse prose this RDR forbids re-parsing everywhere else. The carrier
is a clierr-local representation (plain strings/struct), since `clierr` is a
leaf package and MUST NOT import `internal/resolve`. **The field name, its
clierr-local type, and whether `Count` serializes are A9's to settle
(Pending) — GUESS only.** This helper is deferred: no `flow` verb exists
at HEAD to host it.

---

## 3. Data model (persisted / passed across the boundary)

### Boundary-crossing values

- **`Row`** (kernel-internal, but constructed by every table producer —
  the population this RDR targets, per A5 hand-built at HEAD): single shape,
  `Escape []EscapeEdge`, `Writes []Tag`, `NextTags []Tag`, `RuleID string`,
  `SourceLocator string`. The illegal combination this RDR outlaws is
  `len(Escape) != 0 && len(Writes) != 0` in one `Row`. No row-kind field,
  no ordinary/escape type split.
- **`Table`** — `Rows []Row`; the receiver of `CheckValid`, so the
  whole-table scope is carried by the receiver, not restated at each call
  site.
- **`RowRef`** — `{RuleID, SourceLocator string}`; the identity a producer
  fixes by. Ordered by `compareRefs` (RuleID then SourceLocator; returns 0
  when both match — **not** a total order over distinct rows). Rows built
  with no source identity carry the zero value `RowRef{"",""}` and collide.
- **`Plan`** — emitted by `planOf`; `Writes []Tag`, `NextTags []Tag`,
  `Escaped bool`, `RuleID`, `SourceLocator`. `planOf` copies `row.Writes`
  unconditionally — unchanged, because under the precondition an escape row
  no longer carries writes to copy. No stripping logic is added (silent
  divergence was rejected).
- **`EscapeShapeBreachError`** — `{Ref RowRef; Count int}`; the
  per-identity diagnostic. `Count` is per-identity and counts
  **pre-collapse** rows (3 rows sharing `RowRef{"",""}` → one element,
  `Count == 3`; a single-row breach → `Count == 1`, uniform not degenerate).
- **Aggregate** — `errors.Join` of the per-identity errors;
  `Unwrap() []error` exactly one level deep, every element a
  `*EscapeShapeBreachError`, never a nested join.
- **`CLIError`** (clierr side) — existing `{Code, Detail, Group, Cause, Hint}`
  (`Cause` is `json:"-"`) **plus a NEW `omitempty` field** carrying the
  clierr-local row-identity representation. **GUESS: field name + type —
  A9 Pending.** `Detail` carries the human sentence; the new field carries
  the identities; the two must not be conflated.

### What is NOT persisted / NOT crossed

- The kernel is stateless by contract — no cache, no memoization (a cache
  keyed on a caller-supplied value would reintroduce the ambient state RDR
  0001 forbids). `CheckValid` is re-run on every `Resolve` entry.
- The breach error is a programmer-mistake diagnostic, not a modeled
  disposition; it never enters `Refusal`, never reaches RDR 0004's accessor
  layer, never mutates owned tag state.

---

## 4. Top-level pseudo-code of `Resolve` (~35 lines)

```
func Resolve(in Input) (Result, error) {
    // step 0 — whole-table escape-row shape conformance precondition.
    // RDR: prepended to Resolve's numbered evaluation-order doc list as step 0;
    //      first statement, above `view := assemble(in)`. Cheapness preference
    //      (assemble is pure) but the doc list MUST be amended or it ships wrong.
    if err := in.Table.CheckValid(); err != nil {
        return Result{}, err           // zero Result: Plan nil, Refusal nil → Refused()==false
    }                                  // error returned VERBATIM (no fmt.Errorf wrap)

    view := assemble(in)                // candidate partition: len(row.Escape)!=0 → escape candidate

    survivors, refusal := gate(view)   // guard evaluation over candidates
    if refusal != nil {                 // guard_unevaluable / owned_state_unavailable
        return Result{Refusal: refusal}, nil
    }
    if len(survivors) == 0 {            // no_match → escapeOrRefuse
        return escapeOrRefuse(view)     // escape via rescues, or refuse unmodeled_outcome
    }
    if len(survivors) > 1 {             // ambiguous_match
        return Result{Refusal: ambiguousRefusal(survivors)}, nil
    }
    return Result{Plan: planOf(survivors[0])}, nil   // exact-one: copies Writes/NextTags onto Plan
}

// CheckValid — one predicate, two call sites (also called by producers at
// construction time). The RDR fixes the outcome; the collapse mechanism is GUESS.
func (t Table) CheckValid() error {
    counts := map[RowRef]int{}          // GUESS: collapse via map keyed on identity
    var order []RowRef                   // GUESS: distinct-identity insertion order
    for _, row := range t.Rows {          // WHOLE table — not only reached rows
        if len(row.Escape) != 0 && len(row.Writes) != 0 {   // length-based, NOT nil-ness; NOT NextTags
            r := RowRef{row.RuleID, row.SourceLocator}
            if _, seen := counts[r]; !seen { order = append(order, r) }
            counts[r]++                  // pre-collapse row count per identity
        }
    }
    if len(order) == 0 { return nil }     // allocation-free conforming path
    slices.SortFunc(order, compareRefs)  // non-stable; equal identities already collapsed → no positional tiebreak (REQ-2/REQ-10)
    errs := make([]error, 0, len(order))
    for _, r := range order {
        errs = append(errs, &EscapeShapeBreachError{Ref: r, Count: counts[r]})
    }
    return errors.Join(errs...)          // aggregate; Unwrap()[]error one level; wraps even the single-breach case
}

func (e *EscapeShapeBreachError) Unwrap() error { return ErrEscapeShapeBreach }
// GUESS: Error() wording — identity MUST also be recoverable via errors.As on e.Ref, never only from this string:
func (e *EscapeShapeBreachError) Error() string {
    return fmt.Sprintf("escape-row shape breach: rule %q at %s carries %d owned-tag write(s) on an escape row",
        e.Ref.RuleID, e.Ref.SourceLocator, e.Count)
}
```

---

## GUESS ledger (where the RDR is silent — candidate rewrites)

1. **`Input` field set** — RDR references `in.Table` and `assemble(in)` but
   never enumerates `Input`'s fields. GUESS: embeds `Table` plus the input
   tuple (observed/recognized tags).
2. **Element type names for `Escape` and `NextTags`** — RDR names `Tag` for
   `Writes` only. `Escape []EscapeEdge` and `NextTags []Tag` are GUESSes.
3. **`EscapeShapeBreachError.Error()` message string** — RDR forbids
   recovering identity from prose, so the wording is unspecified. GUESS.
4. **`ErrEscapeShapeBreach` sentinel message** — GUESS.
5. **Multi-breach collapse mechanism** — RDR fixes the outcome
   (collapse equal identities, `Count` = pre-collapse rows, `compareRefs`
   order, `errors.Join`) but not the algorithm. GUESS: map-keyed counts +
   ordered distinct slice + non-stable sort. Sort-then-dedupe not ruled out.
6. **New `CLIError` field name + clierr-local type + whether `Count`
   serializes** — explicitly A9-pending ("the clause binds the carrier CLASS
   only"). GUESS only.
7. **`wrapEscapeBreach` helper name + location** — no `flow` verb at HEAD;
   GUESS.
8. **`Hint` exact wording** — RDR gives one example sentence; GUESS it's
   normative.
9. **File location of `ErrEscapeShapeBreach` / `*EscapeShapeBreachError` /
   `CheckValid`** — RDR says "the kernel package"; GUESS `resolve.go` vs a
   new `breach.go`.
10. **`RefusalKind`'s five names** — RDR refers to
    `unmodeled_outcome`, `no_match`, `ambiguous_match`,
    `owned_state_unavailable`, `guard_unevaluable` in scenario 4; GUESS
    that is the closed five (the RDR does not say "these are the five").
11. **Whether `Table` carries fields beyond `Rows`** — GUESS.
12. **`Plan`'s full field set** — RDR names `Writes`, `NextTags`, `Escaped`,
    `RuleID`, `SourceLocator`; GUESS nothing else.

The RDR is otherwise unusually determinate on the load-bearing contracts:
predicate shape, whole-table scope, precedence, the typed-error surface
(sentinel + `*EscapeShapeBreachError{Ref, Count}` + `Unwrap`), the
one-level aggregate, identity collapse + `Count`, the stable CLI `Code`,
and the authored-path presence-keyed rejection are all pinned in
Normative Contracts. The silences above are mostly naming/prose/mechanism,
not contract gaps — except **#6 (A9's envelope field)**, which is a
genuinely open contract the RDR itself marks Pending and owes to Stage 6.