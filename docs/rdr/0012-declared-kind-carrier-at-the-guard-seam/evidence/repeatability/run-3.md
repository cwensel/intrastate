model: claude-fable-5-1
variant: full (profile: foundational)

# Repeatability reconstruction — run 3

Source elements read: C1, C2, C3, C4, C5, MVV, S1–S9 (via the
projector, each once). Widened past the contract spans into §Approach,
§Load-Bearing Decisions, §Illustrative Code, Phases 1–4, and assumptions
A3 and A4 — for the construction-site call frames (A3), the held-side
render path (A4), the `Evaluator` field/helper naming (Illustrative
Code, marked "shape only"), and the phase ordering (Phase 4 after
Phase 1). §Illustrative Code is explicitly non-load-bearing, so names
lifted from it are marked GUESS below even where the record sketches
them.

---

## 1. Public API

### `internal/guard`

```go
// Evaluator is the runtime guard value seam. Its only state is the
// key → declared-kind mapping fixed at construction; it holds no
// resolve.TagSet, no runtime tag value, no view-typed state (C1).
type Evaluator struct {
    kinds map[string]string // unexported; GUESS: field name "kinds" (Illustrative Code sketch)
}

// NewEvaluator is the sole non-test construction form (C1).
// kinds maps tag key → kind token in {enum, bool, int, set, scalar}.
func NewEvaluator(kinds map[string]string) Evaluator

// DeclaredKinds is the exported, singular producer of the mapping (C1).
// TOTAL over m.Tags; a declared key whose Kind == "" is OMITTED.
// GUESS: a nil *table.Model returns an empty non-nil map (record silent).
func DeclaredKinds(m *table.Model) map[string]string

// Evaluate is UNCHANGED in signature (0007:C1, REQ-10); value-receiver
// satisfies resolve.GuardEvaluator.
func (ev Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult
```

Error modes of `Evaluate` (no Go `error` return — the seam is
single-valued; BR2's two-valued seam is rejected):

- Nil mapping (bare `Evaluator{}`, `var`, `new`, embedded zero): every
  `eq`/`in` atom → `GuardUnevaluable` (LOUD; C1). Operator-inferred
  arms (`lt/lte/gt/gte`, `contains`, set arms) are unaffected.
- `int` kind: held or literal (any member, for `in`) fails
  `strconv.Atoi` → `GuardUnevaluable`; parse-all-then-compare, one bad
  member poisons the list. Platform `int` width; overflow →
  `GuardUnevaluable` on that target (C2).
- `bool` kind: held or literal not in {`true`,`false`} →
  `GuardUnevaluable`; otherwise token equality.
- `enum` | `scalar`: exact string equality, one shared arm.
- `set` kind under `eq`/`in`: `GuardUnevaluable` (defensive).
- Key absent from mapping (undeclared, or declared with empty Kind):
  `GuardUnevaluable`.
- `default` over an unknown kind token: `GuardUnevaluable` (never
  string equality).

### `internal/resolve` (kernel; unchanged shapes, one widened test surface)

```go
type GuardAtom struct {   // JDR 0001 §D1 — UNCHANGED, no Kind member
    Key      string
    Operator string
    Literal  string       // GUESS: string; for `in`, members are held
                          // in some list form the record does not fix —
                          // GUESS: []string members alongside Literal,
                          // or Literal split by a fixed delimiter.
    Block    string       // GUESS: type of Block (record names it only)
}

type GuardResult int      // GUESS: int enum
const (
    GuardFalse GuardResult = iota   // GUESS: ordering
    GuardTrue
    GuardUnevaluable
)

type GuardEvaluator interface {
    Evaluate(atom GuardAtom, value string) GuardResult
}

// Conformance suite — signature WIDENED (C3, amends 0007:REQ-71).
// The suite constructs the seam over its own published fixture, so a
// caller cannot construct over the wrong map and compile.
func TestGuardEvaluatorContract(t *testing.T,
        newSeam func(kinds map[string]string) GuardEvaluator)

// Published declaration fixture, kernel-owned, no internal/table type.
// GUESS: identifier name and whether it is a var or a func returning a
// fresh copy per call.
var ContractKinds = map[string]string{ /* one key per kind token, five entries minimum */ }
```

### `internal/table` — Phase 4 ingress refusal (C5)

No new exported function. `(*loader).atom` gains a canonicality check
after `conform(decl, operator, members)` succeeds for an `int`-declared
key. Refusal surfaces under the EXISTING `malformed_predicate_atom`
code with message
`<site>: "00" is not the canonical spelling of int tag n; write 0`.
GUESS: the Go error value is whatever type `(*loader).atom` already
returns for `malformed_predicate_atom` (record names the code, not the
type); no new envelope field, no new refusal kind.

### `internal/cli`, `internal/graphlint` — construction sites (C4)

```go
// internal/cli/flow_resolve.go
func guardSeam(m *table.Model) resolve.GuardEvaluator   // was zero-arg
// internal/cli/flow_next.go — receives the constructed evaluator,
// does NOT widen to *table.Model.
func probeRow(/* existing params */, ev resolve.GuardEvaluator) /* existing returns */
//   GUESS: parameter position and whether other existing params are reordered.
```

Graphlint: evaluator constructed at `matchSatisfiable(m *table.Model, …)`
and threaded down two frames to `atomAdmitsValue`.

---

## 2. Three most important internal helpers

1. **`guard.DeclaredKinds(m *table.Model) map[string]string`** — the
   single mapping producer. Iterates `m.Tags`, emits `key → decl.Kind`
   for every declaration with a non-empty `Kind`, omits empty-Kind
   declarations (so the absent-key arm, not a `""` arm, handles them).
   Exported because three packages construct; one loop, not three.
   GUESS: iteration is over the map without ordering concerns since the
   result is itself a map.

2. **The `eq`/`in` kind-dispatch arm inside `Evaluator.Evaluate`**
   (GUESS: a private method such as `(ev Evaluator) compareTyped(kind,
   op string, held string, members []string) resolve.GuardResult`).
   Responsibilities: look up `ev.kinds[atom.Key]`; exhaustive `switch`
   over the five tokens with `default → GuardUnevaluable`; `int` arm
   parses ALL sides first (`strconv.Atoi`) then compares; `bool` arm
   checks token membership on both sides; `enum`/`scalar` share the
   string-equality arm; `set` → `GuardUnevaluable`. Must never be
   consulted by `lt/lte/gt/gte`/`contains`/set-operator arms, whose
   reading stays operator-inferred.

3. **`(*loader).atom` canonical-int check in
   `internal/table/normalize.go`** (GUESS: factored as a private
   `canonicalIntSpelling(decl, members, site) error`). After the
   existing `conform(decl, operator, members)` passes and `decl.Kind ==
   "int"`, for each member `n, _ := strconv.Atoi(s)`; if
   `strconv.Itoa(n) != s` return the `malformed_predicate_atom`
   refusal naming the site and `write <Itoa(n)>`. Covers the bare-float
   path (`eq = -0.0` → `"-0"` from `valueMembers`) because it runs on the
   rendered member string. Leaves `conformKind` untouched (0024 REQ-7/9
   BOUNDARY tests pin the permissive shared arm).

Honourable mention (test-side but structurally load-bearing): the
`go/ast` / `go/packages` zero-value check wired into `make check`
(S4) — fails on any composite literal, `var`, `new`, or struct field
of type `guard.Evaluator` in non-test files outside `NewEvaluator`'s
body. GUESS: shipped as a small `cmd`/`tools` program or a Go test
under `internal/guard` with a build tag; the record fixes only that it
is typed and runs in `make check`.

---

## 3. Data model across the boundary

Nothing new is persisted. Cross-boundary shapes:

| Shape | Owner | Direction | Notes |
|---|---|---|---|
| `map[string]string` (tag key → kind token) | `internal/guard` produces; `internal/cli`, `internal/guard`, `internal/graphlint` consume | model → evaluator at construction | Kind tokens: `enum`, `bool`, `int`, `set`, `scalar` (0003 spelling). Empty-Kind keys omitted. Held privately in `Evaluator`. |
| `map[string]string` published fixture | `internal/resolve` (kernel-owned) | suite → `newSeam` constructor | Crosses no `internal/table` type (A5). Meta-check asserts its token set is exactly the five-kind vocabulary. Keys chosen so each case's operator is matrix-admitted for that key's kind. |
| `resolve.GuardAtom{Key, Operator, Literal, Block}` | JDR 0001 §D1 | caller → seam | UNCHANGED; no kind stamped (D-identity). Atom identity tuple `(tag, operator, literal)`. |
| `resolve.GuardResult` ∈ {`GuardTrue`, `GuardFalse`, `GuardUnevaluable`} | kernel | seam → consumer | Every consumer handles `GuardUnevaluable` as its own case; `!= GuardFalse` collapses are forbidden (C4). |
| `guard_unevaluable` refusal payload, reason `uncomparable` | 0007:C8 / 0011 vocabulary | runtime → user | Unchanged shape; CLI maps to `flow-guard-unevaluable`. Names the guarded row and the atom. |
| `malformed_predicate_atom` refusal | `internal/table` load | lint/load → user | Existing code; NEW message text for non-canonical int spelling. No envelope field added. |
| `flow-tag-invalid` (`--tag iter=many`) | `internal/cli` `canonicalValue` | CLI → user | Unchanged before/after (S8); fires upstream of the seam. |

Persisted artifact (`OwnedSnapshot`) is read, not changed: the owned
value reaches the seam unconformed as a raw string (MVV step 1), which
is exactly why the seam must type it.

---

## 4. Top-level pseudo-code — `flow resolve` request through the seam

```
resolveRequest(req):                                   # internal/cli/flow_resolve.go
  m      := req.model                                  # *table.Model, loaded; C5 already refused
                                                       #   any non-canonical int predicate literal at
                                                       #   (*loader).atom during load
  ev     := guardSeam(m)                               # ONE evaluator per request (C4)
           = guard.NewEvaluator(guard.DeclaredKinds(m))
  view   := assemble(OwnedSnapshot, --tag values)      # --tag already conformed by canonicalValue
  for row in candidateRows(m, req.outcome):
      verdict := probeRow(row, view, ev)               # ev threaded down, not rebuilt per row
      switch verdict:
        GuardTrue:        keep row
        GuardFalse:       prune row
        GuardUnevaluable: return Fail(guard_unevaluable,
                             reason=uncomparable, row, atom)   # → flow-guard-unevaluable, exit 2
  return plan(kept rows)                               # or flow-no-match if none kept

probeRow(row, view, ev):
  for atom in row.guardAtoms:                          # exists never reaches the seam (0007:C1)
      held := view.value(atom.Key)                     # absent key → handled upstream by 0007's
      v := ev.Evaluate(atom, held)                     #   key-absent verdict; GUESS: probeRow passes ""
      if v != GuardTrue: return v                      # GUESS: first non-True verdict short-circuits,
  return GuardTrue                                     #   Unevaluable preferred over False? — record silent

Evaluator.Evaluate(atom, value):                        # internal/guard
  switch atom.Operator:
    "lt","lte","gt","gte":  return intCompare(value, atom.Literal)      # operator-inferred, unchanged
    "contains", set arms:   return existing arms                         # unchanged, no kind lookup
    "eq","in":
      kind, ok := ev.kinds[atom.Key]                   # nil map → ok=false → Unevaluable (LOUD)
      if !ok: return GuardUnevaluable
      members := literalMembers(atom)                  # GUESS: single element for eq
      switch kind:
        "int":
          held, err := strconv.Atoi(value); if err: return GuardUnevaluable
          parsed := []; for m in members: n, err := Atoi(m); if err: return GuardUnevaluable; parsed += n
          return boolResult(held ∈ parsed)             # parse-all-then-compare
        "bool":
          if value ∉ {true,false}: return GuardUnevaluable
          for m in members: if m ∉ {true,false}: return GuardUnevaluable
          return boolResult(value ∈ members)
        "enum","scalar": return boolResult(value ∈ members)   # raw-string equality
        "set":           return GuardUnevaluable       # matrix rejects eq/in over set; defensive
        default:         return GuardUnevaluable       # unknown token / never string-equality
```

Lint path (graphlint) mirrors the construction step once at
`matchSatisfiable(m)` and threads `ev` to `atomAdmitsValue`, whose
return becomes three-valued; both callers (`reach.go` positive,
`analysis.go::nodeMeetsAll` negated) branch on `GuardUnevaluable`
explicitly. GUESS: lint over-approximates an unevaluable atom as
admissible AND emits a finding — the record states this policy is
deliberately undecided.

---

## Silences the record leaves (all marked GUESS above)

- `GuardAtom.Literal` representation for `in` (list vs delimited string) and `Block`'s type.
- `GuardResult`'s underlying type / constant ordering.
- The published fixture's Go identifier and var-vs-func form.
- `Evaluator`'s field name and the name/shape of the private dispatch helper (Illustrative Code is non-load-bearing).
- `probeRow`'s exact new parameter position; graphlint's intermediate frame between `matchSatisfiable` and `atomAdmitsValue`.
- `DeclaredKinds` on a nil model.
- Per-row verdict aggregation order in `probeRow` (Unevaluable vs False precedence when a row has several atoms).
- The C5 refusal's Go error type (code and message are fixed; type is not).
- The S4 typed zero-value check's packaging (tool vs test) — only "typed, in `make check`" is fixed.
- Lint's disposition of an unevaluable atom (explicitly deferred by C4).
- The name of `conformingContractSeam`'s new constructor (C3 requires one; name unstated).
