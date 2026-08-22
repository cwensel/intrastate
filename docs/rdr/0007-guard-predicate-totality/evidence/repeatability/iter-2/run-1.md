model: claude-opus-5
variant: full (profile: foundational)

Model: claude-opus-5

# Reconstruction of RDR 0007 — Guard predicate totality (document-only)

Scope: `internal/resolve` (package `resolve`), Go. Reconstructed from
`docs/rdr/0007-guard-predicate-totality.md` alone. Every element not fixed by
the RDR text is marked **GUESS**.

---

## 1. Public API

### 1.1 Entry point (unchanged from RDR 0001, cited by 0007)

```go
// Resolve is the kernel's single pure entry point. No I/O; the result is a
// function of the input tuple (RDR 0001 REQ-1).
func Resolve(in Input) (Result, error)
```

- `error` (non-nil Go error) is reserved for **programmer/table errors**, not
  refusals. The RDR names exactly one such case, inherited from RDR 0009: a row
  with a non-empty `Escape` list and a non-empty `Writes` slice MUST cause
  `Resolve` to return a non-nil error identifying the offending row by `RuleID`
  and `SourceLocator`, with no `Result` disposition, evaluated over the whole
  table before any evaluation step.
- This RDR adds **no new error channel**. JDR 0001 §D1 removes the guard
  reconstruction step, so "no mapping failure exists and no panic or error
  channel is needed for one."
- **GUESS**: exact spelling `Resolve(in Input) (Result, error)` and the field
  `in.Guards` / `in.Owned` / `in.Observed` naming is taken from the quoted
  fragment `gate(candidates, in.Guards, view)` and `Input.Owned` /
  `Input.Observed`; the full `Input`/`Result` struct shape is not stated.

### 1.2 Guard atom (new; JDR 0001 §D1 shape, landed here)

```go
// GuardAtom is one parsed guard predicate: key, operator token, literal, block.
type GuardAtom struct {
    Key      string   // canonical tag key; exact string equality at the kernel
    Op       Operator // operator token (opaque to the kernel except OpExists)
    Literal  string   // verbatim bytes the normalizer emitted
    Block    Block    // BlockAll | BlockUnless
}
```

- The four fields (key, operator token, literal, block ∈ {all, unless}) are
  normative. **GUESS**: the Go field names `Key`/`Op`/`Literal`/`Block` — `Key`
  and `Op` appear in the RDR's quoted spike atom
  `{Key: "guard.missing", Op: resolve.OpExists}`, and `Literal` appears in the
  payload clause; `Block` appears as an "exported constant type the atom
  carries."
- **GUESS**: `Operator` and `Block` as distinct named types (the RDR says
  `Block` is "the same exported constant type the atom carries," which implies a
  named type; it says the same closed-set idiom as `RefusalKind`).

```go
type Operator string          // GUESS: underlying type
type Block string             // GUESS: underlying type

// The one grammar fact the kernel knows (A16, JDR 0001 §D4(b)).
// The normalizer MUST emit these for existence atoms.
const OpExists Operator = "exists"       // GUESS: literal value
const LiteralTrue  = "true"              // GUESS: value + type
const LiteralFalse = "false"             // GUESS: value + type

const (
    BlockAll    Block = "all"     // GUESS: names and values
    BlockUnless Block = "unless"
)
```

### 1.3 Evaluator seam (narrowed)

```go
// GuardEvaluator decides ONE value-comparing atom against ONE present value,
// under RDR 0003's typed operator semantics. It never sees the view, never
// sees provenance, never sees sibling atoms, and never sees tag declarations
// (A17). It is never consulted for an existence atom, nor for an atom whose
// key is absent.
type GuardEvaluator interface {
    Evaluate(atom GuardAtom, value string) GuardResult
}
```

- Signature `Evaluate(atom, value) GuardResult` is normative.
- The seam MAY answer `GuardUnevaluable` for a present value it cannot compare
  (A18). It is three-valued, not boolean.
- A **nil** seam is legal: it yields unevaluable only for an atom that would
  have been handed to it (a value-comparing atom over a present key).
- **GUESS**: `GuardEvaluator` is an interface rather than a func type, and the
  value parameter is `string` (`resolve.Tag.Value` is stated to be a bare
  `string`).

### 1.4 Verdicts and refusal kinds (reused unchanged from RDR 0001)

```go
type GuardResult int   // GUESS: underlying type
const (
    GuardTrue GuardResult = iota
    GuardFalse
    GuardUnevaluable
)

type RefusalKind string  // GUESS: underlying type
func RefusalKinds() []RefusalKind   // closed five-kind set; frozen by REQ-7
```

Named kinds appearing in the RDR: `KindGuardUnevaluable`
(`guard_unevaluable`), `KindOwnedStateUnavailable`
(`owned_state_unavailable`), plus `no_match` and `ambiguous_match`. The fifth is
not named in this document. **GUESS**: the fifth is an internal/table-error
kind. This RDR adds **no sixth kind** — that alternative is explicitly rejected.

### 1.5 Row (reopened)

```go
type Row struct {
    RuleID        string      // GUESS: type
    SourceLocator string      // GUESS: type
    Match         []Tag       // GUESS: shape; closed-world match pattern
    Guard         []GuardAtom // WAS: string. Empty slice = unguarded row.
    RequiresOwned []string    // owned tag keys the row's Writes depend on
    Writes        []Tag
    Escape        []string    // GUESS: refusal-class names; closed to
                              // no_match and ambiguous_match by RDR 0002
}
```

- `Row.Guard` becoming an atom slice is normative. An empty slice is an
  unguarded row, decided TRUE without consulting the view.
- `Row.RequiresOwned` is **narrowed normatively**: the owned tag keys the row's
  **post-guard transition** depends on — the keys its `Writes` require, an
  authored clear included (normalized as a `<clear>` write). It is *not* a
  guard-input declaration. The guard's key set is NOT required to be a subset of
  it. On an escape row it is empty by composition (RDR 0009 empties `Writes`).
- **GUESS**: field names other than `Guard`, `RequiresOwned`, `Writes`,
  `Escape`, `RuleID`, `SourceLocator` (all of which the RDR names).

### 1.6 Refusal payload (new exported surface, A23)

```go
type Refusal struct {
    Kind         RefusalKind
    Rows         []RowRef        // every undecidable row (shipped)
    MissingOwned []string        // shipped
    Undecided    []UndecidedRow  // NEW — replaces Refusal.Guard string
    // Refusal.Guard string is REMOVED.
}

type UndecidedRow struct {
    RuleID        string
    SourceLocator string
    Atoms         []UndecidedAtom
}

type UndecidedAtom struct {
    Key      string
    Block    Block       // the same exported constant type the atom carries
    Operator Operator
    Literal  string
    Reason   Reason
}

type Reason string       // GUESS: underlying type
const (
    ReasonAbsent       Reason = "absent"        // key not in the view
    ReasonUncomparable Reason = "uncomparable"  // present, seam could not compare
)
func Reasons() []Reason  // GUESS: the enumerator, modeled on RefusalKinds()
```

- All the names above (`Undecided`, `UndecidedRow`, `UndecidedAtom`, its five
  fields, `ReasonAbsent`, `ReasonUncomparable`) are normative in the RDR. The
  `Reason` set is kernel-owned and closed, "spelled and enumerated the way
  `RefusalKind`/`RefusalKinds()` already spell the refusal taxonomy."
- **GUESS**: the type name `Reason` and the existence of a `Reasons()`
  enumerator — the RDR says "a test can pin the set exactly," which implies one.
- `RowRef` is shipped (`compareRefs` is cited). **GUESS**: it carries
  `(RuleID, SourceLocator)`.

### 1.7 Error modes (complete list from the document)

| Mode | Channel | Notes |
| --- | --- | --- |
| `guard_unevaluable` | `Refusal.Kind` | Non-escapable by construction (RDR 0002 closes escape lists to `no_match`/`ambiguous_match`). Payload names row + atoms. |
| `owned_state_unavailable` | `Refusal.Kind` | Survivor missing an owned `RequiresOwned` key. Wins over `guard_unevaluable` within one survivor set. |
| `no_match` | `Refusal.Kind` | Escapable. Downstream of all three gate facts. |
| `ambiguous_match` | `Refusal.Kind` | Escapable. |
| Non-conforming escape row | non-nil Go `error` | RDR 0009 entry precondition; whole-table, pre-evaluation; identifies the row by `RuleID` + `SourceLocator`; no `Result` disposition. |
| Load/parse failure | not a kernel mode | RDR 0002's typed validation; reaches the user as an RDR 0005 load error, never a resolution refusal. `internal/resolve` has no load path. |
| Panic | **none** | Explicitly foreclosed by §D1. |

**GUESS**: the fifth refusal kind's name (not stated anywhere in this document).

### 1.8 Exported test-contract surface (Phase 3)

```go
// GUESS: name and signature entirely. The RDR says only "a small
// contract-test function for the value seam" that RDR 0003 instantiates.
func TestGuardEvaluatorContract(t *testing.T, mk func() GuardEvaluator)
```
The `contains` leg is **blocked** until RDR 0003 declares the element encoding
of a set-valued tag value.

---

## 2. The three most important internal helpers

### 2.1 `evaluateAtom(atom GuardAtom, view TagSet, seam GuardEvaluator) GuardResult`

The domain rule's enforcement site — the single function that makes absence
structurally unmaskable. Responsibilities, in order:

1. Decide presence of `atom.Key` in the assembled view via `TagSet.Lookup`'s
   `ok`, **provenance-blind**. Never `TagSet.has` (owned-only, serves
   `missingOwned`, deliberately not reused).
2. If `atom.Op == OpExists`: decide from presence alone, `presence == literal`.
   `exists = true` → TRUE when present; `exists = false` → TRUE when absent. The
   seam MUST NOT be consulted. If the literal is neither `LiteralTrue` nor
   `LiteralFalse` (the empty literal included), return `GuardUnevaluable` with
   reason `uncomparable` — never decide from presence, never treat a missing
   literal as `LiteralFalse`. Comparison is verbatim: no parsing, case-folding,
   or coercion.
3. Else (value-comparing operator) with key **absent**: return
   `GuardUnevaluable`, reason `absent`. The seam MUST NOT be consulted — the
   SQL:2003 strict-routine rule: the host never invokes the routine.
4. Else (value-comparing, key **present**): return `seam.Evaluate(atom, value)`.
   A nil seam here yields `GuardUnevaluable`. A foreign operator token falls
   into this branch, which is why token drift fails closed.

**GUESS**: the name `evaluateAtom` is stated in the RDR as a proposed name; its
exact signature (whether it returns a reason alongside the verdict, or the
reason is recomputed at payload-mint time) is not. I model it as returning
`(GuardResult, Reason)` in the pseudo-code below.

### 2.2 `evaluateGuard(row Row, view TagSet, seam GuardEvaluator) GuardResult`

Shipped name, extended. Owns **strong-Kleene combination** — the tables are
normative and "the evaluator holds no part of them."

- Empty atom slice → `GuardTrue` without consulting the view (shipped branch).
- Partition atoms by `Block`. `all_result` = conjunction (`min` under
  `F < U < T`) over the `all` atoms; empty `all` → `T`.
- `unless_conj` = conjunction over the `unless` atoms. `unless` is **block-level
  negation, not per-atom negation**.
- Row verdict = `all_result ∧ ¬(unless_conj)`, with `¬T=F, ¬F=T, ¬U=U`. When the
  `unless` block is empty or omitted, the `¬(unless_conj)` **term drops out**
  and the verdict is `all_result` alone — it is neither a vacuously-true
  conjunction (which would disable every such row) nor `F`.
- MUST evaluate **every** atom of every survivor — **no short-circuit** on a
  decided block — because the payload must name every unevaluable atom.
- `F ∧ U = F` is witnessed falsity and prunes: this is what keeps D8 sound.

### 2.3 `undecidedAtoms` / `compareAtoms` (payload mint + total order)

Proposed names in the RDR, modeled on shipped `compareRefs`/`copyTags`.

- `undecidedAtoms` collects, per surviving row whose verdict is
  `GuardUnevaluable`, one `UndecidedAtom` **per unevaluable atom** — never
  twice for the same atom (duplicate suppression is by construction of the emit
  loop, not a property of the sort).
- `compareAtoms` orders payload entries on the **total** six-tuple
  `(RuleID, SourceLocator, key, block, operator token, literal)`, compared field
  by field in that order. Row-identity-then-key is *not* total: A5's conjoined
  value row (`X exists = true` beside `X eq v`) and the `unless` idiom
  (`legal_hold exists = true` beside `legal_hold eq true`) put two atoms on one
  key in one row. The literal is compared as the **exact bytes the normalizer
  emitted** (set-valued literal spelling is an open RDR 0003 declaration).
- The result must be a function of the input tuple (REQ-1), never of atom or row
  order. `gate`'s `slices.MinFunc` representative-row selection **retires** with
  `Refusal.Guard`; every undecidable row appears.
- A `copyAtoms` helper keeps the slice-valued payload tuple-deterministic.

Honourable mentions the RDR names but that are unchanged: `TagSet.Lookup` (the
presence predicate), `assemble` (faithful; never drops a key), `missingOwned`
(owned-provenance only; unchanged), `escapeOrRefuse` (ZERO diff in the spike),
`gate` (partition byte-identical; only its payload-building hunk and its
doc comment change).

---

## 3. Data model across the boundary

### 3.1 Nothing is persisted

`Resolve` is pure and in-process; "no new I/O." There is no serialization,
hashing, or canonicalization in the kernel — REQ-37 bans a kernel symbol named
`Canonicalize`, `Hash`, `Checksum`, `Fingerprint`, or `Digest`, and REQ-36 bans
encode/decode/inverse operations.

### 3.2 Normalizer (RDR 0002) → kernel: the `Row`

| Element | Producer | Kernel contract |
| --- | --- | --- |
| Tag keys | 0002 normalizer canonicalizes case, namespace, whitespace | Kernel compares **exact strings**; performs no canonicalization (A22 — duty not yet a 0002 clause) |
| Existence operator token | 0002 normalizer emits the kernel's exported `OpExists` | Kernel compares verbatim; a foreign token becomes a value atom (fails closed) |
| Existence literal | 0002 emits `LiteralTrue`/`LiteralFalse`; rejects a foreign literal at load | Kernel compares verbatim; foreign or missing literal ⇒ unevaluable/`uncomparable` |
| Set-valued literal spelling | RDR 0003 must declare the element encoding | Kernel compares the emitted bytes; `resolve.Tag.Value` is a bare `string` today |
| `RequiresOwned` | RDR 0002 producer (§JD-3, still open) | Kernel only reads it in `missingOwned` |

### 3.3 Kernel → evaluator seam (RDR 0003)

One direction, one shape: `(GuardAtom, value string) → GuardResult`. The seam
receives **no view, no provenance, no sibling atoms, no tag declarations**.
Declared bounds and element universes are LINT inputs, never evaluation inputs.
Ceiling: one tag, one literal, one operator — a relational, temporal, or
aggregate operator would need a seam extension and is out of scope.

### 3.4 Kernel → CLI (RDR 0005): the refusal payload

Logical shape — a two-level sorted array of objects, the kernel's first:

```
Undecided: [
  { RuleID, SourceLocator,
    Atoms: [ { Key, Block, Operator, Literal, Reason } , ... ] },
  ...
]                       # sorted on (RuleID, SourceLocator, Key, Block,
                        # Operator, Literal); no duplicate atoms
```

**Transport is UNRESOLVED in the document** and I do not invent a resolution:

- JDR 0001 §D4 (Closed) routes it through §JD-8's `Detail` as **flattened
  text**.
- A19 (Pending) proposes ONE `omitempty` structured field on
  `internal/cli/clierr/clierr.go::CLIError` (today: flat strings `Code`,
  `Message`, `Param`, `Detail`, `Hint`), which requires reopening §D4. The
  kernel side is unblocked either way — the payload exists on `Refusal`
  regardless.
- `Code` values: `guard_unevaluable` → `flow-guard-unevaluable`,
  `GroupUserEnv`; `owned_state_unavailable` has **no code row yet** (§JD-8 gap).

### 3.5 The assembled view

`TagSet`, built by `assemble` from `Input.Owned` / `Input.Observed` /
recognized tags, with precedence owned > observed > recognized. Precedence
overwrites a **value**, never a key's presence. Readers: `Lookup` (exported,
provenance-blind `ok`), `has` (unexported, owned-only, single caller
`missingOwned`), `matches` (key+value, provenance-free), `Len`.
**GUESS**: `TagSet` is `map[string]Tag`-backed with a `tags` field (the RDR
mentions `view.tags`).

---

## 4. Top-level pseudo-code of `Resolve`

```
func Resolve(in Input) (Result, error):
  # RDR 0009 entry precondition — whole table, before any evaluation
  for row in in.Table.Rows:
      if len(row.Escape) > 0 and len(row.Writes) > 0:
          return zeroResult, error("escape row must have empty Writes: " +
                                   row.RuleID + " @ " + row.SourceLocator)

  view := assemble(in.Owned, in.Observed, in.Recognized)   # faithful; drops no key
  candidates := match(in.Table.Rows, view)   # closed-world: absent key => no match

  selected, blocked := gate(candidates, in.Guards, view)
  if blocked != nil: return refuse(in, *blocked), nil

  switch len(selected):
    case 1: return plan(selected[0], view), nil
    case 0: return escapeOrRefuse(in, view, KindNoMatch), nil
    default: return escapeOrRefuse(in, view, KindAmbiguousMatch), nil

# GATE, THEN COUNT (JDR 0001 §D2)
func gate(rows, seam, view) (survivors []Row, blocked *Refusal):
  verdicts := {}
  for row in rows:
      verdicts[row] = evaluateGuard(row, view, seam)   # every atom; no short-circuit

  # 1. D8 — GuardFalse prunes first; pruned rows report nothing
  survivors := [row for row in rows if verdicts[row] != GuardFalse]

  # 2. owned state among SURVIVORS (unevaluable rows ARE survivors)
  if missing := missingOwned(survivors, view); missing != nil:
      return nil, &Refusal{Kind: KindOwnedStateUnavailable,
                           MissingOwned: missing, Rows: refsOf(...)}

  # 3. undecidable veto — resolution-level, no decided-TRUE sibling escapes
  undecided := [row for row in survivors if verdicts[row] == GuardUnevaluable]
  if len(undecided) > 0:
      payload := undecidedAtoms(undecided, view, seam)   # one entry per atom
      sort(payload, compareAtoms)  # (RuleID, SourceLocator, Key, Block, Op, Literal)
      return nil, &Refusal{Kind: KindGuardUnevaluable,
                           Rows: refsOf(undecided),      # ALL rows, no MinFunc
                           Undecided: payload}

  return [row for row in survivors if verdicts[row] == GuardTrue], nil
```

**GUESS** in the above: that `gate` recomputes or caches per-atom reasons for
`undecidedAtoms` (the RDR does not say whether the verdict pass carries reasons
forward or the payload pass re-walks the atoms); that `match` is a separate step
named that way; that `Result`/`plan`/`refsOf` are spelled so. The **ordering**
(precondition → assemble → match → prune → owned → undecidable → count →
escape), the no-short-circuit rule, the retirement of `MinFunc`, and the
six-field sort key are all normative, not guesses.
