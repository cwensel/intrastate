model: claude-fable-5
variant: full (profile: foundational)

Model: claude-fable-5

# Reconstruction from RDR 0007 alone: guard predicate totality

Source: `docs/rdr/0007-guard-predicate-totality.md` only. Everywhere the RDR
is silent, the guess is marked **GUESS**.

## 1. Public API

Package `internal/resolve` (Go), one non-test file `resolve.go`, single pure
entry point. Zero production importers; the package boundary is the module's
public API for this reconstruction.

### Entry point

```go
// Resolve is the single pure entry point. In-process, no I/O, pure over the
// input tuple (RDR 0001 REQ-1: results are a function of the input tuple).
func Resolve(in Input) (Result, error)
```

GUESS: the exact return shape. The RDR shows `return refuse(in, *blocked), nil`
and "a non-nil Go error ... with no Result disposition", so the pair
`(Result, error)` is documented; whether `Result` is a value or `*Result` is a
GUESS (value chosen, matching `refuse(...)` returning something directly).

**Error modes** (three disjoint channels, all documented):

1. **Refusal disposition** — `Result` carries a `Refusal` and a nil `Plan`.
   Closed five-kind taxonomy pinned by
   `TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`. Four kinds are named in
   the RDR: `no_match`, `ambiguous_match`, `owned_state_unavailable`,
   `guard_unevaluable`. GUESS: the fifth kind (never named in this RDR);
   guessed `invalid_input`-style kind. This RDR adds NO sixth kind — the new
   payload is a field on `Refusal`, not a kind.
2. **Go error** — table-shape breach at entry (RDR 0009's kernel
   precondition, cited here): a row with a non-empty `Escape` list and
   non-empty `Writes` causes `Resolve` to return a non-nil error identifying
   the offending row by `RuleID` and `SourceLocator`, before any evaluation,
   with no `Result` disposition. (Owned by 0009; visible at this API.)
3. **No panic channel** — the row carries parsed atoms (JDR 0001 §D1), so no
   guard-reconstruction failure exists and no panic or error channel is
   needed for one.

### Types on the boundary

```go
// Input to Resolve.
// GUESS at exact field names beyond those the RDR names: it documents
// in.Guards (the evaluator seam) and Input.Owned / Input.Observed tag inputs;
// the table/candidate rows and any recognized-tag input are GUESSED fields.
type Input struct {
    Rows     []Row          // GUESS (name): the transition table's rows
    Owned    []Tag          // owned snapshot (documented name Input.Owned)
    Observed []Tag          // caller-supplied context (documented)
    Guards   GuardEvaluator // the value seam; nil is legal (per-atom rule)
    // GUESS: recognized tags enter somewhere; provenance "recognized" exists.
}

type Tag struct {
    Key   string
    Value string // bare string; set-valued encoding is RDR 0003's, undeclared
}

// Row — the §D1 reshape this RDR carries.
type Row struct {
    RuleID        string   // GUESS at type; row identity half 1
    SourceLocator string   // GUESS at type; row identity half 2
    Match         []Tag    // GUESS (name): closed-world match pattern
    Guard         []GuardAtom // WAS string; empty slice = unguarded row
    RequiresOwned []string // NARROWED: owned keys the row's Writes depend on
    Writes        []Tag    // post-guard writes (clear renders as <clear> write)
    Escape        []RefusalKind // GUESS at element type; only no_match /
                                // ambiguous_match are legal (RDR 0002 closure)
}

// GuardAtom — the §D1 shape: (key, operator token, literal, block).
type GuardAtom struct {
    Key     string
    Op      string // operator token; GUESS that it is a string, since the
                   // kernel compares it verbatim against OpExists
    Literal string // GUESS: string, compared verbatim (LiteralTrue/False)
    Block   Block
}

type Block string // GUESS at underlying type; exported constants required

const (
    BlockAll    Block = "all"    // GUESS at spelling
    BlockUnless Block = "unless" // GUESS at spelling
)

// The one grammar fact the kernel owns (§D4(b)): the existence operator
// token and its two boolean literal forms, exported so 0002's normalizer
// MUST emit exactly these. Values GUESSED; names documented.
const (
    OpExists     = "exists" // GUESS at value
    LiteralTrue  = "true"   // GUESS at value
    LiteralFalse = "false"  // GUESS at value
)

// GuardEvaluator — narrowed to a per-atom VALUE seam. It never sees the
// view, provenance, sibling atoms, or tag declarations; only value-comparing
// atoms over PRESENT keys reach it. It may answer GuardUnevaluable for a
// present value it cannot compare (A18).
// GUESS: interface vs func type; func type chosen for a minimal seam.
type GuardEvaluator func(atom GuardAtom, value string) GuardResult

// GuardResult — shipped three-valued verdict, reused unchanged.
type GuardResult int // GUESS at underlying type
const (
    GuardTrue GuardResult = iota // GUESS at ordering/values
    GuardFalse
    GuardUnevaluable
)

// Result / Plan.
// GUESS at Result's exact shape: the RDR names Plan (with Writes derived
// from it), Refusal, and an Escaped:true marker observed in behavior.
type Result struct {
    Plan    *Plan
    Refusal *Refusal
    Escaped bool // GUESS: documented as observed output ("Escaped:true")
}

type Plan struct {
    Writes []Tag // GUESS beyond Writes, which A13 documents ("plan's Writes")
}

// Refusal — Guard string is REMOVED (replaced), Undecided is added.
type Refusal struct {
    Kind         RefusalKind
    Rows         []RowRef  // every implicated row; RowRef GUESSED as
                           // {RuleID, SourceLocator}
    MissingOwned []string  // owned_state_unavailable payload
    Undecided    []UndecidedRow // NEW: the per-row/per-atom payload; the
                                // kernel's first two-level payload
    // GUESS: other flat scalar fields exist ("flat scalars plus ...");
    // unnamed in the RDR.
}

type RowRef struct { // GUESS at exact shape
    RuleID        string
    SourceLocator string
}

// The named payload surface (normative names).
type UndecidedRow struct {
    RuleID        string
    SourceLocator string
    Atoms         []UndecidedAtom
}

type UndecidedAtom struct {
    Key      string
    Block    Block  // the SAME exported constant type the atom carries
    Operator string // GUESS at field type (matches GuardAtom.Op)
    Literal  string // GUESS at field type
    Reason   Reason
}

// Reason — kernel-owned closed constant set, spelled/enumerated the way
// RefusalKind / RefusalKinds() spell the refusal taxonomy.
type Reason string // GUESS at underlying type
const (
    ReasonAbsent       Reason = "absent"       // key not in the view
    ReasonUncomparable Reason = "uncomparable" // present, seam could not compare
)
// GUESS: a Reasons() enumerator mirroring RefusalKinds(), so a test can pin
// the set exactly (the RDR says "a test can pin the set exactly").

// RefusalKind + RefusalKinds() — shipped closed taxonomy, unchanged.
type RefusalKind string // GUESS at underlying type
func RefusalKinds() []RefusalKind
const (
    KindGuardUnevaluable      RefusalKind = "guard_unevaluable"
    KindOwnedStateUnavailable RefusalKind = "owned_state_unavailable"
    KindNoMatch               RefusalKind = "no_match"        // GUESS at name
    KindAmbiguousMatch        RefusalKind = "ambiguous_match" // GUESS at name
    // fifth kind: GUESS, unnamed in this RDR
)

// TagSet — assembled evaluation view. Lookup is the ONLY exported accessor:
// provenance-blind presence + value.
type TagSet struct{ /* unexported map keyed by Tag.Key */ }
func (t TagSet) Lookup(key string) (value string, ok bool) // GUESS at exact
                                                           // return shape
func (t TagSet) Len() int // documented reader; GUESS it is exported

// Phase 3: exported contract-test function for the value seam that RDR
// 0003's implement stage instantiates against its evaluator.
// GUESS at signature:
func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)
```

### API-level guarantees (the contract, condensed)

- Value operators are PARTIAL over the view: absent key ⇒ unevaluable, never
  true/false, and the seam is never consulted for it (kernel-enforced).
- `exists` is the sole TOTAL operator: verdict is `presence == literal`,
  kernel-decided, seam never consulted. Foreign token ⇒ treated as value
  atom; foreign/missing literal on `OpExists` ⇒ unevaluable, `uncomparable`.
- Presence is provenance-blind; key identity is exact string equality
  (canonicalization is the normalizer's, upstream).
- Atom verdicts combine under strong Kleene (`min` over `F < U < T`;
  `¬U = U`); row verdict is `all_result ∧ ¬(unless_conj)`; empty/omitted
  `unless` drops the term; empty atom slice ⇒ TRUE without the view.
- Gate, then count: prune GuardFalse (D8) → owned sweep over survivors →
  undecidable veto → exact-one count → escape reachability. `guard_unevaluable`
  is non-escapable. One unevaluable survivor vetoes the whole resolution.
- The refusal payload is sorted on the total tuple
  `(RuleID, SourceLocator, key, block, operator token, literal)` and is a
  function of the input tuple, never of atom or row order.
- Nil seam is a per-atom fact: kernel-decidable atoms (existence, absent-key)
  resolve identically with or without a seam; only a value atom over a
  present key goes unevaluable for want of a seam (and a nil seam is never
  itself a payload reason).

## 2. Three most important internal helpers

1. **`evaluateAtom(atom GuardAtom, view TagSet, seam GuardEvaluator) (GuardResult, Reason-or-none)`**
   (proposed name, stated in the RDR). The enforcement site of the domain
   rule. Decides presence via `TagSet.Lookup` (provenance-blind); decides
   existence atoms entirely in the kernel (`presence == literal`, verbatim
   comparison against `OpExists`/`LiteralTrue`/`LiteralFalse`, no parsing —
   foreign literal ⇒ unevaluable/`uncomparable`); marks absent-key value
   atoms unevaluable (`absent`) without calling the seam; hands only
   present-key value atoms to `Evaluate(atom, value)`; maps a nil seam to
   unevaluable only for atoms that WOULD have been handed to it. GUESS: the
   exact return shape (result plus a reason for the payload when unevaluable).

2. **`evaluateGuard(row, view, seam) GuardResult`** (shipped name, extended).
   The per-row verdict hook. Empty atom slice ⇒ GuardTrue without consulting
   the view (shipped branch). Otherwise evaluates every atom (no
   short-circuit on a decided block — the payload must name every
   unevaluable atom), combines under the normative strong-Kleene tables:
   conjunction per block, then `all_result ∧ ¬(unless_conj)` with the
   `unless` term dropping out when the block is empty/omitted. GuardFalse
   only when decided atoms alone witness it (`F ∧ U = F`); any unresolved
   dependence on a U atom ⇒ GuardUnevaluable. GUESS: it also returns/collects
   the per-atom unevaluable records for `gate`, since `gate`'s only hunk
   builds the refusal payload.

3. **`gate(candidates, seam, view) (selected []Row, blocked *Refusal)`**
   (shipped name and shape). The pipeline: prune GuardFalse rows first (D8 —
   pruned rows contribute neither candidacy nor owned obligations); run
   `missingOwned` over SURVIVORS (unevaluable rows included) and refuse
   `owned_state_unavailable` first when any survivor's `RequiresOwned` key is
   absent from the owned snapshot; then refuse `guard_unevaluable` if any
   survivor is unevaluable, minting the sorted `Undecided` payload over EVERY
   undecidable row (the `slices.MinFunc` representative-row selection is
   retired with `Refusal.Guard`); otherwise return survivors for the
   exact-one count. Reused identically for the escape row set by
   `escapeOrRefuse` (uniform gating; an unevaluable escape row's refusal
   replaces the candidate-set refusal).

Runners-up, documented but less central: `assemble` (builds the TagSet with
owned > observed > recognized precedence GUESS on full order — owned-over-
observed is documented — never dropping a key), `missingOwned`
(owned-provenance sweep via unexported `TagSet.has`), `escapeOrRefuse`,
and the proposed payload plumbing `undecidedAtoms`/`compareAtoms`/`copyAtoms`
(tuple-deterministic sort/copy, modeled on shipped `compareRefs`/`copyTags`).

## 3. Data model across the boundary

Nothing is persisted — `Resolve` is pure, in-process, no I/O. Three
boundary-crossing shapes:

**a. Row (normalizer → kernel), produced by RDR 0002's normalizer.**
`Row.Guard []GuardAtom` where an atom is the §D1 tuple
`(key, operator token, literal, block ∈ {all, unless})`. Preconditions on the
producer: keys canonicalized upstream (kernel compares bytes); existence
atoms carry exactly `OpExists` + `LiteralTrue`/`LiteralFalse`; every
guard-referenced key is declared (typos fail table load); escape rows have
empty `Writes` (0009) and hence empty `RequiresOwned` (composition).
`RequiresOwned []string` = owned keys the row's `Writes` depend on,
post-guard only. Set-valued literals (`in`, `contains`) reach the kernel as
one opaque string with no declared element encoding (blocked on RDR 0003);
the kernel compares literal bytes verbatim.

**b. Refusal payload (kernel → caller/renderer).**
`Refusal.Undecided []UndecidedRow`; `UndecidedRow = {RuleID, SourceLocator,
Atoms []UndecidedAtom}`; `UndecidedAtom = {Key, Block, Operator, Literal,
Reason ∈ {absent, uncomparable}}`. Additive on `Refusal`; replaces
`Refusal.Guard string`. Sorted on the total six-field tuple; one entry per
unevaluable atom, no duplicates by construction of the emit loop. Explicit
non-promise: pruned rows (`F ∧ U = F`) report no absences — the guarantee is
"never a plan from absence," not "every absence reported."

**c. CLI transport (kernel refusal → RDR 0005's failure envelope).**
Unresolved in this RDR: JDR 0001 §D4 (Closed) routes the payload through
`CLIError.Detail` as flattened text; A19 proposes reopening §D4 to grant ONE
`omitempty` structured field on `CLIError` (two-level sorted array of
objects: per row `(RuleID, SourceLocator)`, per atom key/block/reason).
§D4 stands until reopened; fallback is `Detail` flattening. `Code` values:
`guard_unevaluable` maps to `flow-guard-unevaluable` / `GroupUserEnv` (exit
2); `owned_state_unavailable` has no code row yet. GUESS: the structured
field, if granted, is named something like `Undecided` and mirrors the
kernel payload's field names — the RDR does not name it and forbids this RDR
from stating it.

## 4. Top-level pseudo-code of `Resolve`

```
Resolve(in):
  # 0009 entry precondition, whole table, before any evaluation
  for row in in.Rows:
    if row.Escape non-empty and row.Writes non-empty:
      return no-disposition, error(row.RuleID, row.SourceLocator)

  view = assemble(in)            # owned > observed precedence; keys never dropped
  candidates = rows whose closed-world Match pattern matches view   # absent
                                 # match key => not a candidate (no_match side)

  selected, blocked = gate(candidates, in.Guards, view)
  if blocked != nil: return refuse(in, *blocked), nil

  switch len(selected):
    case 1:  return plan(selected[0]), nil          # exactly-one rule
    case 0:  return escapeOrRefuse(no_match, ...), nil
    default: return escapeOrRefuse(ambiguous_match, ...), nil

gate(rows, seam, view):
  survivors = []; undecided = []
  for row in rows:
    verdict = evaluateGuard(row, view, seam):
      if row.Guard empty: GuardTrue                 # no view consulted
      per atom: presence = view.Lookup(atom.Key)    # provenance-blind
        if atom.Op == OpExists:
          if atom.Literal not in {LiteralTrue, LiteralFalse}: U (uncomparable)
          else: (presence == literal) ? T : F       # kernel-decided, total
        elif not presence: U (absent)               # seam NEVER called
        elif seam == nil:  U                        # per-atom, value atoms only
        else: seam.Evaluate(atom, value)            # T, F, or U (uncomparable)
      combine ALL atoms (no short-circuit), strong Kleene:
        row_verdict = conj(all_atoms) ∧ ¬conj(unless_atoms)   # unless term
                                                    # drops out if block empty
    if verdict == GuardFalse: continue              # D8: prune, no obligations
    survivors += row
    if verdict == GuardUnevaluable: undecided += (row, its U-atoms + reasons)
  missing = missingOwned(survivors, view)           # owned provenance only
  if missing: return nil, Refusal(owned_state_unavailable, rows, missing)
  if undecided: return nil, Refusal(guard_unevaluable,
                    rows = every undecided row,
                    Undecided = sort by (RuleID, SourceLocator,
                                         key, block, operator, literal))
  return survivors, nil

escapeOrRefuse(kind, ...):
  gate the escape row set identically; an unevaluable escape row's
  guard_unevaluable refusal replaces the candidate refusal; a viable
  escape row plans with Escaped = true; else refuse(kind)
```

GUESS: candidate matching happens inside `Resolve` before `gate` (the RDR
says `gate(candidates, ...)` and that the match pattern folds absence into
non-candidacy, but never shows the matching step's location or name).
