model: claude-opus-5[1m]
variant: full (profile: foundational)

# Repeatability run-1 — RDR 0009 reconstruction

Reconstructed from `docs/rdr/0009-escape-row-shape-conformance-ownership.md`
alone. Every point the RDR leaves open is marked **GUESS**.

## 1. Public API

Package `internal/resolve` (the RDR 0001 kernel; this RDR adds to its existing
surface).

```go
// Table.CheckValid returns nil if the table's rows are valid, or else an
// error describing a problem.
//
// It checks exactly one property: escape-row shape conformance — a row with a
// non-empty Escape list must have an empty Writes slice. Nothing else about
// the table is checked; in particular the Escape-class restriction to
// no_match/ambiguous_match that Row's doc records for RDR 0002 is not checked,
// so a nil return must not be read as general table validity.
func (t Table) CheckValid() error

// ErrEscapeShapeBreach is the category sentinel for an escape row that carries
// owned-tag writes. Classify with errors.Is.
var ErrEscapeShapeBreach = errors.New("escape-row shape breach")   // GUESS: message text; the RDR pins the identifier and its exportedness, not the string

// EscapeShapeBreachError reports one breaching row identity.
type EscapeShapeBreachError struct {
    Ref   RowRef // the offending row's (RuleID, SourceLocator)
    Count int    // pre-collapse breaching rows sharing Ref; >= 1
}

func (e *EscapeShapeBreachError) Error() string  // GUESS: format string unpinned
func (e *EscapeShapeBreachError) Unwrap() error  // returns ErrEscapeShapeBreach
```

Unchanged, but now with a documented error case:

```go
// Resolve evaluates the input tuple against the table.
//
// Step 0 (new): the whole table is checked for escape-row shape conformance
// before any evaluation step. On breach, Resolve returns the zero Result and a
// non-nil error.
func Resolve(t Table, in Input) (Result, error)
```

`Row`, `Table`, `Plan`, `Result`, `RowRef`, `Refusal`, `RefusalKind`, `TagSet`,
`Tag` keep their existing shapes. **No** new type split, **no** row-kind field,
**no** sixth refusal kind.

### Error modes

| Condition | Return |
| --- | --- |
| ≥1 row with `len(Escape) != 0 && len(Writes) != 0` | `Result{}` (Plan nil, Refusal nil) + non-nil error, always an `errors.Join` aggregate |
| conforming table | exactly as today: one of five modeled refusals or a plan, **nil** error |

Breach precedence: the error precedes **every** modeled disposition, including
`unmodeled_outcome`, and fires for dormant rows no resolution path reaches.

The returned error's shape is fixed:

- always the `errors.Join` aggregate, even for a single breach (uniform caller
  shape — never the bare per-row error in the one-breach case);
- `Unwrap() []error` is **exactly one level deep**; every element is a
  `*EscapeShapeBreachError`; never a nested join;
- `Resolve` returns `CheckValid`'s error **verbatim** — no `fmt.Errorf` wrap,
  which would expose `Unwrap() error` at the outermost level and break the flat
  traversal;
- callers classify/extract with `errors.Is` / `errors.As` / `errors.AsType`,
  never a direct type assertion or `==` against the sentinel.

### Caller hazard (contract, not just a note)

On breach `Result{}` reports `Refused() == false` with a nil `Plan`. Callers
MUST check the error before reading the disposition.

### CLI side (RDR 0005 surface, deferred to the future `flow` verb)

The verb wraps a breach into `*clierr.CLIError`:

- `Code: "escape-row-shape-breach"` (pinned by the RDR)
- `Group: clierr.GroupInternal` → exit 2
- `Cause`: the kernel error (not wire-visible; `json:"-"`)
- `Hint`: "fix the table producer: an escape row must carry no writes"
- a **new** `omitempty` field carrying the offending row identities in a
  clierr-local representation (`clierr` must not import `internal/resolve`).

**GUESS**: the new field's name and type — the RDR explicitly defers name, type,
and whether `Count` serializes to assumption A9 (still Pending). I would write
`Rows []ErrorRow \`json:"rows,omitempty"\`` with
`type ErrorRow struct { RuleID, SourceLocator string; Count int \`json:",omitempty"\` }`.

## 2. Three most important internal helpers

1. **The per-row breach predicate.** `len(row.Escape) != 0 && len(row.Writes) != 0`
   — length-based, never nil-sensitive; `Writes`-only, does not extend to
   `NextTags` (A4 settled closed).
   **GUESS**: whether this is a named unexported helper (e.g. `breachesShape(row Row) bool`)
   or inlined in `CheckValid`'s loop. The RDR pins only the exported
   `Table.CheckValid`.

2. **Collapse-and-order over the breaching rows.** Gathers each breaching row's
   `RowRef`, collapses equal identities into a single entry while counting
   pre-collapse rows, sorts by the existing `compareRefs` (RuleID then
   SourceLocator), and never uses `Table.Rows` position as a tiebreak.
   **GUESS**: its name and signature. I would write
   `func collapseBreaches(rows []Row) []*EscapeShapeBreachError`. The RDR fixes
   the *behavior* (collapse + per-identity `Count` + compareRefs order)
   completely, but names no function.
   **GUESS**: the mechanism — map keyed on `RowRef` then sort, vs. sort-then-scan.
   Both satisfy the contract; the RDR pins the output, not the algorithm.
   Sorting must use a comparison consistent with `compareRefs`; because
   `compareRefs` returns 0 for equal identities and collapsing removes duplicates
   before sorting, stability is not required post-collapse.

3. **Reuse of the existing `rowRefs` / `compareRefs` pair.** The kernel already
   owns row-identity extraction and ordering for REQ-1/REQ-10 refusal payloads;
   the breach report reuses the same ordering rather than minting one.
   **GUESS**: whether `rowRefs` is reused directly or only `compareRefs` is. The
   RDR names both as reused surfaces but doesn't say which call sites change.
   `rowRefs` sorts with the non-stable `slices.SortFunc`, which is safe here
   only after collapsing.

Unchanged existing helpers: `assemble`, `gate`, `escapeOrRefuse`, `planOf`,
`rescues`, `copyTags`, `refuse`. `planOf` in particular stays unconditional — no
stripping logic, since the precondition makes a write-bearing escape plan
unconstructable.

## 3. Data model across the boundary

Nothing is newly persisted. The values that cross boundaries:

```go
type Row struct {
    RuleID        string   // source identity RDR 0002 requires every normalized row to retain
    SourceLocator string
    Escape        []...    // non-empty ⇒ this is an escape row (the sole discriminator)
    Writes        []Tag    // MUST be empty when Escape is non-empty  ← the invariant
    NextTags      ...      // NOT constrained by this RDR
    // predicates, guard, etc. — unchanged
}

type RowRef struct {
    RuleID        string
    SourceLocator string
}
```

**GUESS**: `Escape`'s element type (the RDR says "modeled failure class list"
and `[no_match]`, so plausibly `[]RefusalKind` or a dedicated failure-class
type). Only its emptiness matters here.

Breach report as a value (what tests traverse):

```
errors.Join(
   &EscapeShapeBreachError{Ref: RowRef{"a.rule","f.toml:1"}, Count: 1},
   &EscapeShapeBreachError{Ref: RowRef{"",""},              Count: 3},
)   // one level, ordered by compareRefs, equal identities collapsed
```

Rows built without source identity all collapse to a single `RowRef{"",""}`
entry with `Count` = the number of such rows — the degenerate-producer case the
RDR requires to stay diagnostic via the count.

**GUESS**: where `RowRef{"",""}` sorts relative to populated identities.
`compareRefs` compares strings, so the empty pair sorts first — but the RDR
never states this, and no scenario asserts it.

Wire model (CLI, deferred): the `CLIError` envelope plus the new `omitempty`
identities field (A9-pending, above).

## 4. Top-level pseudo-code

```
func (t Table) CheckValid() error:
    breaches := empty map[RowRef]int          # identity -> pre-collapse count
    for row in t.Rows:
        if len(row.Escape) != 0 and len(row.Writes) != 0:
            ref := RowRef{row.RuleID, row.SourceLocator}
            breaches[ref] += 1
    if len(breaches) == 0:
        return nil

    entries := []*EscapeShapeBreachError{}
    for ref, count in breaches:
        entries = append(entries, &EscapeShapeBreachError{Ref: ref, Count: count})
    sort entries by compareRefs(a.Ref, b.Ref)   # never by Table.Rows position

    errs := []error{}
    for e in entries: errs = append(errs, e)
    return errors.Join(errs...)                 # aggregate even when len == 1

func Resolve(t Table, in Input) (Result, error):
    # step 0 — whole-table precondition, first statement, before assemble.
    # Cheapness preference only: assemble is pure, so placing it after would
    # yield identical dispositions. Nothing may come to depend on the order.
    if err := t.CheckValid(); err != nil:
        return Result{}, err                    # VERBATIM — no fmt.Errorf wrap

    # everything below is unchanged from RDR 0001
    view := assemble(in)
    if unmodeled outcome: return refuse(KindUnmodeledOutcome, ...), nil
    candidates, escapes := partition(t.Rows) by len(row.Escape) != 0
    survivors := gate(candidates, view)
    switch len(survivors):
        case 1: return Result{Plan: planOf(survivors[0], view)}, nil
        default: return escapeOrRefuse(escapes, survivors, view), nil
```

**GUESS**: the exact body of `Resolve` below step 0 — the RDR names `assemble`,
the `len(row.Escape) != 0` partition, `gate`, `escapeOrRefuse`, `planOf`, and
`rescues`, and states the five refusal kinds and that `unmodeled_outcome` is
`Resolve`'s first branch, but does not reproduce the control flow. The shape
above is inferred from those names and is not pinned by the RDR.

## Notes on RDR determinacy (observations while reconstructing)

Determinate — reconstructed with no choice left to me:

- ownership split (RDR 0002 load-time / kernel entry backstop), the channel (Go
  error, never a refusal kind), and that no type split occurs;
- the four exported spellings and their fields;
- aggregate-always, one-level `Unwrap() []error`, verbatim return, per-element
  `errors.Is`;
- collapse + per-identity `Count` counting pre-collapse rows, `compareRefs`
  order, no positional tiebreak;
- whole-table scope, precedence over all five kinds, length-based not
  nil-sensitive, `Writes`-only (not `NextTags`);
- zero-`Result` on breach and the `Refused() == false` caller trap;
- the CLI `Code`, `Group`, `Hint` text, and that identities need a serialized
  carrier that is not `Detail`.

Silences I had to guess (candidates for the diff pass):

- **G1** — the `errors.Is`-visible sentinel's message text and
  `EscapeShapeBreachError.Error()` format. Low stakes (the RDR forbids
  string-matching assertions everywhere), but it is the one thing an operator
  reads.
- **G2** — the serialized `CLIError` carrier field's name and type, and whether
  `Count` rides the wire. Openly deferred to A9, which is **Pending** — this is
  a known-open contract, not an unnoticed silence.
- **G3** — internal helper names/decomposition for the collapse-and-order step,
  and whether `rowRefs` is reused or only `compareRefs`. Taste, not contract.
- **G4** — sort position of the zero-value `RowRef{"",""}` relative to populated
  identities in a mixed multi-breach report. Falls out of `compareRefs`'s string
  comparison, but no clause or scenario states it; scenario 10b uses an
  all-degenerate table and scenario 10 an all-distinct one, so no test pins the
  mixed case.
- **G5** — `Escape`'s element type.
- **G6** — the map-vs-sort-then-scan collapse mechanism (output is pinned, so
  this is genuinely free).
