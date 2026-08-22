model: claude-sonnet-5
variant: full (profile: foundational)

Model: claude-sonnet-5

# Reconstruction from RDR 0007 alone

Source: `docs/rdr/0007-guard-predicate-totality.md`, read in full. No other
file was read. Everything below is reconstructed strictly from that
document's prose; every place the RDR is silent or only implies a shape is
marked **GUESS**.

## 1. Public API

Package: `internal/resolve` (Go). The RDR names these exported symbols
directly; I've assembled them into signatures. Items marked GUESS fill in
field types/order the RDR describes narratively but doesn't spell as a
literal struct.

```go
// Row — a candidate transition-table row. Guard is now a slice of parsed
// atoms (JDR 0001 §D1), not a string. RequiresOwned is narrowed to
// post-guard write-dependency keys only (this RDR's core fix).
type Row struct {
    RuleID         string   // GUESS: type; identity field, named throughout
    SourceLocator  string   // GUESS: type; identity field, named throughout
    Guard          []GuardAtom
    RequiresOwned  []string // owned tag keys row.Writes depends on; empty on escape rows (A21)
    Writes         []Tag    // GUESS: shape; RDR 0009's territory, cited not defined here
    Escape         bool     // GUESS: representation; RDR 0009 speaks of "non-empty Escape list"
    // match-pattern fields exist (TagSet.matches) but are out of this RDR's scope
}

// GuardAtom — key, operator token, literal, block ∈ {all, unless}. Shape
// fixed by JDR 0001 §D1; this RDR cites it, does not restate the grammar.
type GuardAtom struct {
    Key      string
    Operator OperatorToken // GUESS: named type; RDR calls it "operator token"
    Literal  string        // GUESS: type — RDR says literal is byte-comparable scalar; for OpExists it's LiteralTrue/LiteralFalse
    Block    BlockKind     // "all" | "unless" — exported constant type (RDR is explicit this is exported)
}

// GuardEvaluator — narrowed to a per-atom value seam. Only called for a
// value-comparing atom whose key is PRESENT. Never sees the view.
type GuardEvaluator interface {
    Evaluate(atom GuardAtom, value string) GuardResult // GUESS: value's Go type (RDR says "the present value")
}

// GuardResult — three-valued verdict. Named directly in the RDR
// ("resolve.go::GuardResult is three-valued").
type GuardResult int // GUESS: underlying type
const (
    GuardFalse GuardResult = iota
    GuardTrue
    GuardUnevaluable
)

// Exported kernel constants (explicitly named as exported in the RDR):
const (
    OpExists     OperatorToken = "..." // GUESS: literal spelling; RDR only says "one operator token"
    LiteralTrue  Literal       = "..." // GUESS: literal spelling
    LiteralFalse Literal       = "..."
)

const (
    ReasonAbsent      Reason = "absent"       // GUESS: exact string; RDR names the constant and its meaning, not its literal spelling
    ReasonUncomparable Reason = "uncomparable"
)

// Refusal — carries the new payload. Refusal.Guard (a string) is REMOVED
// (not retained beside the new field — that was the spike's wrong shape).
type Refusal struct {
    Kind         RefusalKind   // pre-existing, closed 5-kind taxonomy, unchanged
    Rows         []RowRef      // pre-existing; continues to carry every undecidable row
    MissingOwned []string      // pre-existing (owned_state_unavailable payload)
    Undecided    []UndecidedRow // NEW — this RDR's payload
}

type UndecidedRow struct {
    RuleID        string
    SourceLocator string
    Atoms         []UndecidedAtom
}

type UndecidedAtom struct {
    Key      string
    Block    BlockKind
    Operator OperatorToken
    Literal  Literal
    Reason   Reason // ReasonAbsent | ReasonUncomparable
}

// Resolve — the single pure entry point (named directly: "resolve.go::Resolve
// as the single pure entry point"). Signature is GUESS beyond what's implied.
func Resolve(in Input) (Result, error) // GUESS: exact signature; RDR only
    // shows call-site fragments: `gate(candidates, in.Guards, view)`,
    // `switch len(selected)`, and says Resolve returns "a non-nil Go error"
    // only for RDR 0009's malformed-table precondition (out of this RDR's
    // scope) — otherwise it returns a Result carrying either a Plan or a
    // Refusal.
```

**Error modes** (from the RDR's own vocabulary, not GUESSed):

- No panics, no Go `error` return for guard-domain problems — "no
  reconstruction step, so no mapping failure exists and no panic or error
  channel is needed for one." The *only* Go `error` return named anywhere in
  this RDR belongs to RDR 0009's malformed-table entry precondition
  (out of scope here, just cited).
- All guard-domain failure is communicated through `Result`/`Refusal` with
  `Refusal.Kind`. Five closed refusal kinds exist pre-RDR (`RefusalKinds()`);
  this RDR reuses `guard_unevaluable` and `owned_state_unavailable`, adds no
  sixth kind — it adds a *field* (`Undecided`) to the existing kind, not a
  new kind.
- `guard_unevaluable` is explicitly non-escapable (RDR 0002's closure
  excludes it from escape lists).

## 2. Three most important internal helper functions

The RDR names these as existing or to-be-extended kernel internals. I pick
the three most load-bearing for the totality rule specifically:

1. **`evaluateGuard`** (extended, not new) — the per-row verdict hook. Its
   responsibility: for each atom in a row's `Guard` slice, decide presence
   of the referenced key (`TagSet.Lookup`'s `ok`), route existence atoms to
   a kernel-only presence check, route absent-key value atoms to
   `unevaluable` without calling the seam, route present-key value atoms to
   the seam, then combine all atom verdicts under strong-Kleene logic
   (`all_result ∧ ¬(unless_conj)`) into one `GuardResult` for the row. This
   is where "presence, existence, and combination" all live per the RDR's
   own Load-Bearing Decisions section.

2. **`gate`** (existing, structurally unchanged, payload-populating hunk
   added) — implements "gate, then count" (JDR 0001 §D2): partitions
   candidates into pruned (GuardFalse) / survivors; among survivors runs the
   owned-state scan (`missingOwned`) before the undecidable scan; returns
   `blocked` before `Resolve` ever counts `selected`. The RDR is explicit
   that `gate`'s control flow ("prune → owned-state → undecidable order")
   is byte-identical pre/post-RDR — only the refusal-payload-building hunk
   inside it changes (replacing the retired `slices.MinFunc`
   representative-row pick with an exhaustive, sorted `Undecided` list).

3. **`TagSet.Lookup`** (existing, reused as-is) — the sole presence
   predicate: "returns `ok` on any key present regardless of provenance."
   This is explicitly *not* `TagSet.has` (which is owned-provenance-only and
   feeds only `missingOwned`); the RDR is emphatic the two tests must never
   be conflated, since `Lookup` is what makes presence provenance-blind
   (owned/observed/recognized all count).

**GUESS** (helper names the RDR proposes as new but doesn't fully specify):
`evaluateAtom` (per-atom decision — presence, existence-by-presence,
seam-dispatch), `undecidedAtoms`/`compareAtoms`/`copyAtoms` (build/sort/copy
the new payload, modeled on shipped `compareRefs`/`copyTags`). The RDR names
these as "proposed names; none exists under `internal/` today," so their
exact signatures are not stated — I'd guess something like:

```go
func evaluateAtom(atom GuardAtom, view *TagSet, seam GuardEvaluator) (GuardResult, *UndecidedAtom) // GUESS
func compareAtoms(a, b UndecidedAtom) int // GUESS, ordered on (RuleID, SourceLocator, key, block, operator, literal)
```

## 3. Data model (persisted / passed across the boundary)

Nothing here is disk-persisted per the RDR (`Resolve` is described as "pure
and in-process; no new I/O"). The "boundary" is the in-process contract
between the kernel (`internal/resolve`) and two neighbors: the RDR 0003
value evaluator (per-atom seam) and the RDR 0005 CLI renderer (refusal
payload leaving the kernel toward the CLI envelope).

**Across the kernel ↔ evaluator-seam boundary**, per atom:
- Input: `(atom GuardAtom, value string)` — "the atom and the present
  value," explicitly *never* the view, provenance, sibling atoms, or
  declared tag kind (A17).
- Output: `GuardResult` (three-valued; MAY also answer `GuardUnevaluable`
  for a present-but-unparseable value, A18).

**Across the kernel ↔ CLI boundary** (the refusal payload, `Refusal.Undecided`):
- One entry per row: `RuleID`, `SourceLocator`, `Atoms []UndecidedAtom`.
- One entry per unevaluable atom: `Key`, `Block`, `Operator`, `Literal`,
  `Reason` (closed 2-value set: `absent` | `uncomparable`).
- Total, deterministic sort key: `(RuleID, SourceLocator, key, block,
  operator token, literal)` — required so the payload is a pure function of
  input, never of atom/row iteration order (REQ-1).
- **Not yet wired**: how this payload crosses into the CLI's JSON/text
  envelope is explicitly *unresolved in this RDR* — it's blocked on
  reopening JDR 0001 §D4/§JD-8 to grant one new `omitempty` structured field
  on `clierr.CLIError`, with a fallback of flattening into the existing
  `Detail string` field if that reopening is declined (A19). So the
  cross-process/cross-boundary wire shape is explicitly **not fixed** by
  this document — I mark the wire encoding itself (JSON field names, for
  instance) as **GUESS**, since the RDR defers it entirely to RDR 0005.

**Producer-side data this RDR depends on but does not itself define**
(GUESS-flagged because the RDR explicitly says these are owned elsewhere):
- The normalizer (RDR 0002) emits `Row` values from authored TOML — the
  authored-table-to-`Row` mapping is out of scope here.
- `resolve.Tag.Value` is "a bare string" today; for set-valued literals
  (`in`, `contains`) no element encoding is declared yet — the RDR flags
  this as a real gap (Phase 3 blocked on RDR 0003 declaring it), not
  something I should guess a shape for.

## 4. Top-level pseudo-code of the main operation

This reconstructs the guard-evaluation portion of `Resolve`/`gate` as the
RDR narratively specifies it (SEAM clause, existence-operator clause,
combination clause, GATE-THEN-COUNT clause, payload clause). Line-level
function names beyond what the RDR states are marked GUESS inline.

```
function Resolve(input) -> Result:
    view := assemble(input)                        # existing; provenance-blind merge
    candidates := input.TableRows                   # GUESS: exact source of candidate set
    selected, blocked := gate(candidates, input.Guards, view)
    if blocked != nil:
        return refuse(input, *blocked)
    switch len(selected):
        case 1: return plan(selected[0])
        case 0: return escapeOrRefuse(...)           # existing no_match path
        default: return refuse(ambiguous_match)

function gate(candidates, guards, view) -> (selected, blocked):
    survivors := []
    for row in candidates:
        verdict, undecidedAtoms := evaluateGuard(row.Guard, view)  # per-row
        if verdict == GuardFalse:
            continue                                  # D8: prune, no obligation raised
        survivors.append((row, verdict, undecidedAtoms))

    # owned-state scan runs over ALL survivors before the undecidable scan
    missing := []
    for (row, _, _) in survivors:
        missing += missingOwned(row, view)            # unchanged; owned-provenance only
    if missing != nil:
        return nil, &Refusal{Kind: owned_state_unavailable, MissingOwned: missing, Rows: ...}

    # undecidable scan, over the same survivor set
    undecided := []
    for (row, verdict, atoms) in survivors:
        if verdict == GuardUnevaluable:
            undecided.append(UndecidedRow{row.RuleID, row.SourceLocator, atoms})
    if undecided != nil:
        sort(undecided, by=totalTuple)                 # (RuleID, SourceLocator, key, block, op, literal)
        return nil, &Refusal{Kind: guard_unevaluable, Undecided: undecided, Rows: ...}

    selected := [row for (row, verdict, _) in survivors if verdict == GuardTrue]
    return selected, nil

function evaluateGuard(atoms []GuardAtom, view) -> (GuardResult, []UndecidedAtom):
    allResult := T                                    # empty `all` = TRUE
    unlessConj := UNSET                                # empty/omitted `unless` drops out
    undecided := []
    for atom in atoms:
        present := view.Lookup(atom.Key)               # provenance-blind
        if atom.Operator == OpExists:
            verdict := (present == boolOf(atom.Literal))  # presence == literal; total
        elif !present:
            verdict := Unevaluable                     # kernel marks it; seam NEVER called
            undecided.append(UndecidedAtom{atom.Key, atom.Block, atom.Operator, atom.Literal, ReasonAbsent})
        else:
            verdict := seam.Evaluate(atom, view.ValueOf(atom.Key))  # MAY answer Unevaluable (A18: uncomparable)
            if verdict == Unevaluable:
                undecided.append(UndecidedAtom{..., ReasonUncomparable})

        if atom.Block == all:
            allResult := kleeneAnd(allResult, verdict)
        else: # unless
            unlessConj := kleeneAnd(unlessConj, verdict)

    rowVerdict := allResult
    if unlessConj != UNSET:
        rowVerdict := kleeneAnd(allResult, kleeneNot(unlessConj))
    return rowVerdict, undecided

function kleeneAnd(a, b) -> min(a, b) under F < U < T
function kleeneNot(v) -> {T:F, F:T, U:U}[v]
```

Notes on fidelity to the RDR text:
- The `assemble`/candidate-set extraction lines are GUESS shells around
  named-but-undescribed existing functions (`assemble` is named; the exact
  candidate-row source is not).
- Everything inside `evaluateGuard` and the ordering inside `gate`
  (prune → owned-state → undecidable → count) is stated directly and
  repeatedly in the RDR's normative contracts, not guessed.
- The escape-set path (`escapeOrRefuse`) is named as gated "identically, as
  its own row set" but its internals are not re-derived here since the RDR
  says it needs "ZERO changes."
