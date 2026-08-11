model: claude-fable-5
variant: full (profile: foundational)

# Repeatability run-2 — RDR 0009 reconstruction

Reconstructed from
`docs/rdr/0009-escape-row-shape-conformance-ownership.md` alone.
Every point the RDR leaves open is marked **GUESS**.

## 1. Public API

All additions land in `internal/resolve` (the kernel package).
`Resolve`'s signature is unchanged — the RDR reuses the existing
`(Result, error)` return and RDR 0001's reserved
programmer-mistake error slot.

```go
// CheckValid returns nil if t's rows conform to the escape-row
// shape rule, or else an error describing a problem. It checks
// exactly one property — escape-row shape conformance (a row
// with a non-empty Escape list must have an empty Writes slice)
// — and nothing else; in particular it does not check the
// Escape-class restriction Row's doc records for RDR 0002. A
// nil return is not general table validity.
func (t Table) CheckValid() error
```

GUESS: value receiver `(t Table)` rather than pointer — the RDR
pins "a METHOD ON Table taking no arguments" but not the
receiver form; a read-only structural scan suggests a value
receiver, matching a stateless kernel.

```go
// ErrEscapeShapeBreach is the category sentinel every
// per-identity breach error wraps, so errors.Is classifies the
// breach class from outside the package.
var ErrEscapeShapeBreach = errors.New("escape-row shape breach")
```

GUESS: the sentinel's message text. The RDR pins the identifier
(`ErrEscapeShapeBreach`, package-level, exported) and its role,
not its prose.

```go
// EscapeShapeBreachError reports the breaching rows sharing one
// RowRef identity. Count is the number of pre-collapse rows
// carrying that identity (Count == 1 for a single-row breach).
type EscapeShapeBreachError struct {
    Ref   RowRef
    Count int
}

func (e *EscapeShapeBreachError) Error() string
func (e *EscapeShapeBreachError) Unwrap() error // returns ErrEscapeShapeBreach
```

GUESS: `Error()`'s message format. The RDR requires the row
identity and count be recoverable structurally (never only from
prose), so the message is unpinned diagnostic text — something
like `"escape-row shape breach: rule %q at %s carries writes
(%d row(s))"`.

### Error modes

- **Breach** (any row with non-empty `Escape` and non-empty
  `Writes`): `Resolve` returns the zero `Result` (`Plan` nil,
  `Refusal` nil — `Refused() == false`, so callers MUST check
  the error first) and a non-nil error. The error is
  `Table.CheckValid()`'s return **verbatim** — never
  `fmt.Errorf`-wrapped — and is always the `errors.Join`
  aggregate, even for a single breach: `Unwrap() []error` is
  exactly one level deep, every element a
  `*EscapeShapeBreachError`, `errors.Is(elem,
  ErrEscapeShapeBreach)` holds on each element. Elements are
  ordered by `compareRefs` over `Ref` (never table position);
  equal identities collapse to one element whose `Count` is the
  pre-collapse row count. Callers classify/extract with
  `errors.Is` / `errors.As` / `errors.AsType` and traverse the
  aggregate for multi-breach reads (a single `errors.As` finds
  only the first element).
- **Conforming table**: nil error; every disposition identical
  to today — zero behavior change.
- **CLI surfacing** (deferred; binds the future `flow` verb):
  the verb wraps the kernel error into
  `*clierr.CLIError{Code: "escape-row-shape-breach", Group:
  clierr.GroupInternal}` (exit 2) with a remedy `Hint` and the
  offending identities rendered into a NEW `omitempty`
  serialized field (never `Detail` prose; `Cause` is `json:"-"`
  so the Go chain is not wire-visible). GUESS: the new field's
  name and clierr-local type — the RDR explicitly leaves both
  (plus whether `Count` serializes) to A9; `clierr` is a leaf
  and must not import `internal/resolve`. GUESS: the `Hint`
  wording — the Validation section suggests "fix the table
  producer: an escape row must carry no writes".

No new refusal kind, no `Row` type split, no row-kind field, no
change to `Resolve`'s signature, `planOf`, or the escape
discriminator (`len(row.Escape) != 0` / `rescues`).

## 2. Three most important internal helpers

1. **The per-row breach predicate** — `len(row.Escape) != 0 &&
   len(row.Writes) != 0`, length-based on purpose (a non-nil
   empty `Writes` conforms; nil-ness is never tested). GUESS:
   whether this is a named unexported helper (e.g.
   `breachesEscapeShape(row Row) bool`) or inlined in
   `CheckValid`'s loop — the RDR pins the predicate expression
   and its `Writes`-only scope (A4: no `NextTags` widening) but
   not its packaging.
2. **The aggregate builder inside `CheckValid`** — walks the
   whole `Table.Rows` (every row, not only reached rows),
   collects breaching `RowRef{RuleID, SourceLocator}`
   identities, collapses equal identities while counting
   pre-collapse rows, sorts the distinct identities with the
   existing `compareRefs` ordering, and returns
   `errors.Join(per-identity errors...)`. GUESS: whether this
   is a separate unexported function or the body of `CheckValid`
   itself; GUESS: the accumulation structure (a
   `map[RowRef]int` then a sorted key slice is the natural
   shape, since `RowRef` is two strings and comparable).
3. **`compareRefs` (existing, reused)** — the kernel's
   `(RuleID, SourceLocator)` comparator that `rowRefs` already
   uses for refusal payload stability; the breach report adopts
   the same input-tuple determinism, sorted with
   `slices.SortFunc` over distinct identities (collapse makes
   the order total, so non-stable sort is safe).

## 3. Data model

Nothing new is persisted. Across the kernel boundary:

- **`Row`** — unchanged single shape: `Escape`, `Writes`,
  `NextTags`, `RuleID string`, `SourceLocator string` (plus its
  existing predicate/guard fields). Escape identity remains
  discriminated solely by non-empty `Escape`. The producer
  obligation (non-empty `Escape` ⇒ empty `Writes`) is stated on
  the `Row` doc contract, not encoded in the type.
- **`RowRef`** — existing value `{RuleID string, SourceLocator
  string}`; the breach error's identity payload, mirroring the
  kernel's `Refusal.Rows` diagnostics.
- **`EscapeShapeBreachError`** — `{Ref RowRef; Count int}` as
  above; the only new value crossing the boundary, carried on
  the Go error chain inside the `errors.Join` aggregate.
- **`Result`** on breach — the zero struct (`Plan` nil,
  `Refusal` nil); there is no nil `Result` since `Result` is a
  struct.
- **Wire (CLI, deferred)** — `clierr.CLIError` gains a new
  `omitempty` field carrying the offending identities in a
  clierr-local representation. GUESS: shape — e.g. a slice of
  `{RuleID, SourceLocator string}` pairs in a clierr-defined
  struct; pinned only as "not `Detail` prose, not
  `resolve.RowRef`", the rest pending A9.
- Package imports: `internal/resolve` gains `errors` (and `fmt`
  for the message) — both clear the frozen import guards.

## 4. Top-level pseudo-code of the main operation

```
CheckValid(t Table) -> error:
    counts := empty map[RowRef]int                 # GUESS: map accumulation
    for row in t.Rows:                             # whole table, every row
        if len(row.Escape) != 0 and len(row.Writes) != 0:
            counts[RowRef{row.RuleID, row.SourceLocator}] += 1
    if counts is empty:
        return nil
    refs := keys(counts)
    sort refs by compareRefs                       # identity order, never position;
                                                   # collapse already made it total
    errs := []
    for ref in refs:
        errs.append(&EscapeShapeBreachError{Ref: ref, Count: counts[ref]})
    return errors.Join(errs...)                    # aggregate even when len == 1

Resolve(in) -> (Result, error):
    # step 0 — whole-table entry precondition, first statement,
    # before assemble; precedes every modeled disposition
    if err := in.Table.CheckValid(); err != nil:   # GUESS: field spelling in.Table
        return Result{}, err                       # verbatim; zero Result, no wrap
    view := assemble(in)                           # existing pipeline, unchanged:
    partition candidates by len(row.Escape) != 0   #   ordinary vs escape (rescues)
    gate(...)                                      #   guard evaluation
    selection / escapeOrRefuse(...)                #   exact-one plan, refusal, or escape
    return planOf(selected), nil                   #   planOf unconditional, unchanged
```

GUESS: how the input carries the table — the RDR shows
`assemble(in)` and speaks of "the supplied table", so `in.Table`
is inferred, not quoted.

## Notable RDR silences encountered (summary)

- Receiver form of `Table.CheckValid` (value vs pointer).
- Sentinel message text and `EscapeShapeBreachError.Error()`
  format (deliberately unpinned — identity travels structurally).
- Packaging of the per-row predicate and aggregate builder
  (named helpers vs inline) and the accumulation structure.
- The clierr envelope field's name, type, and whether `Count`
  serializes (explicitly deferred to A9).
- Exact `Hint` prose.
- The `Input` field naming (`in.Table`).
