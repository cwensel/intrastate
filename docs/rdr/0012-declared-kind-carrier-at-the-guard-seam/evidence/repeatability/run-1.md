model: claude-opus-5[1m] (Opus 5, 1M context)
variant: full (profile: foundational)

## Read path

Structure once via `inspect --json --filter elements,outline`; then by id:
C1, C2, C3, C4, C5 (the five normative contracts), MVV, S1–S9.

WIDENED past the contract spans, deliberately, at these points (each
because a contract left a signature, type, error mode, or step order
under-determined):

- `§illustrative-code` — C1/C2 name `NewEvaluator`, `DeclaredKinds`, and
  the arms in prose but never show a receiver, a field name, or how the
  arm returns. The sketch gives `ev.kinds[atom.Key]`, `boolResult(...)`,
  and the `guardSeam(m *table.Model) resolve.GuardEvaluator` shape. Marked
  "Illustrative — shape only, not load-bearing", so I treat the names as
  indicative and the shapes as GUESS where they go beyond C1's prose.
- `§approach` — to confirm the seam interface `Evaluate(atom, value)` is
  frozen and that the change is confined to one package.
- `§problem-statement` — for the error-mode vocabulary at the CLI edge
  (`flow-guard-unevaluable`, `flow-tag-invalid`) and for the ingress
  census; C2 names verdicts but not the refusal codes they surface as.
- `§load-bearing-decisions` (D-identity, D-naming, D-selection-predicate,
  D-canonical-spelling-at-load, D-c5-is-not-alternative-3) — for the
  atom identity tuple (no kind member) and the OPEN A4 fork on the held
  ingress, which C2/C5 do not resolve between them.
- `§phase-1` and `§phase-4` — for step order and for the fact that the
  eighteen zero-value test sites migrate in Phase 1, which C3 implies but
  does not sequence.

Not widened past: the `in` arm's poison rule, the five-kind vocabulary,
the `default` disposition, and the `malformed_predicate_atom` reuse are
all stated verbatim in C2/C5 and are reproduced, not guessed.

---

## 1. The module's public API

### Package `internal/guard`

```go
// Evaluator is the guard value seam's runtime implementation.
// It holds a declared-kind mapping and NOTHING else: no resolve.TagSet,
// no runtime tag value, no view-typed state.
type Evaluator struct {
    kinds map[string]string // key -> declared kind token; UNEXPORTED (C1)
                            // field name `kinds` is illustrative-code only  [GUESS on the name]
}

// NewEvaluator is the only non-test construction form. The zero value
// Evaluator{} is retired but still admissible in Go; its disposition is
// fixed LOUD (see Evaluate).
func NewEvaluator(kinds map[string]string) Evaluator

// Evaluate is UNCHANGED from 0007:C1 / REQ-10 — signature frozen.
func (ev Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult
```

Value receiver on `Evaluate` and value (not pointer) return from
`NewEvaluator`: **GUESS**. C1 says only "constructs" and that the zero
value `Evaluator{}` must be a legal-but-loud literal; a `*Evaluator`
return would make `Evaluator{}` a differently-shaped hazard and the
illustrative `guardSeam` returns the call directly into an interface.
A value type also makes the nil-map read (`ev.kinds[k]` on a nil map)
safe, which is what the LOUD-not-panicking disposition requires.

`NewEvaluator` does not copy or validate `kinds`: **GUESS**. C1 states
totality and the empty-Kind omission as obligations of the *producer*,
not of the constructor, and states no error return — so the constructor
is total and infallible.

### Package `internal/guard` — the declaration producer

```go
// DeclaredKinds is EXPORTED (C1) because the three construction sites
// span internal/cli, internal/guard, internal/graphlint.
// TOTAL over m.Tags: every declared key gets an entry, EXCEPT a key whose
// Kind is the empty string, which is OMITTED rather than mapped to "".
func DeclaredKinds(m *table.Model) map[string]string
```

Homed in `internal/guard`: C1 gives the signature without a package
qualifier, but the illustrative site writes `guard.DeclaredKinds(m)`, and
it takes a `*table.Model` while `internal/table` must not import the
kernel. **GUESS** on the home package (the alternative, `internal/table`,
is equally import-legal).

`DeclaredKinds(nil)`: **GUESS** — returns an empty non-nil map rather
than panicking, by the same LOUD-not-crash reasoning C1 applies to the
zero-value evaluator. The RDR is silent.

### Kernel: `internal/resolve` (types crossed at the boundary)

```go
type GuardResult int // 0007; three-valued
const (
    GuardTrue GuardResult = iota
    GuardFalse
    GuardUnevaluable
)
```
Ordering/underlying type of the constants: **GUESS**. The RDR names the
three verdicts and forbids two-valued collapse but fixes no encoding.

```go
// UNCHANGED — JDR 0001 §D1
type GuardAtom struct {
    Key      string
    Operator string
    Literal  string
    Block    string // match | guard.all | guard.unless   [GUESS on spellings]
}

type GuardEvaluator interface {
    Evaluate(atom GuardAtom, value string) GuardResult
}
```
Field types all `string` and `Block`'s domain: **GUESS**. The RDR names
the four members ("Key, Operator, Literal, Block") and names the three
blocks elsewhere, but fixes neither the Go types nor whether Operator and
Block are named string types.

### Conformance suite (exported test helper, `internal/resolve`)

```go
// internal/resolve/guardcontract.go — signature WIDENED by C3.
// Takes a CONSTRUCTOR, not a constructed seam, so a caller cannot build
// over the wrong map and still compile.
func TestGuardEvaluatorContract(t *testing.T,
        newSeam func(kinds map[string]string) GuardEvaluator)

// The kernel-owned published fixture the suite's cases assume.
// map[string]string — key -> kind token; no internal/table type crosses.
func ContractKinds() map[string]string    // name is a GUESS
```
The fixture's *existence*, type, and ownership are normative (C3, A5);
its exported identifier is not named anywhere in the record. **GUESS** on
`ContractKinds`; a package-level `var` is equally likely.

### Error modes

The seam itself has NO error return — the whole design is that a
three-valued verdict replaces one. The error modes are:

| condition | seam verdict | user-visible surface |
|---|---|---|
| `int` key, held value unparseable | `GuardUnevaluable` | `flow-guard-unevaluable`, reason `uncomparable` (0011's existing vocabulary) |
| `int` key, any `in` member unparseable | `GuardUnevaluable` (whole list poisoned) | same |
| `int` key, held value out of platform `int` range | `GuardUnevaluable` on the target under test | same |
| `bool` key, held value not `true`/`false` | `GuardUnevaluable` | same |
| `set` key under `eq`/`in` | `GuardUnevaluable` (defensive) | same |
| key absent from mapping (undeclared, or declared with empty `Kind`) | `GuardUnevaluable` | same |
| unknown kind token (`default` arm) | `GuardUnevaluable` | same |
| literal unparseable for declared kind | `GuardUnevaluable` (defense in depth; `0003:C8` load rejection is primary) | — |
| nil-mapping evaluator, any `eq`/`in` | `GuardUnevaluable` | same |
| authored non-canonical `int` predicate literal | refused at LOAD, never reaches seam | exit 2, `malformed_predicate_atom`, message `<site>: "00" is not the canonical spelling of int tag n; write 0` |
| CLI `--tag iter=many` on an `int` tag | refused at INPUT | exit 2, `flow-tag-invalid` |

Exit code 2 for the lint/resolve refusals is normative (S7, S8). The
exact JSON envelope for `flow-guard-unevaluable` is **GUESS** beyond
"the payload names the guarded row and the `iter` atom with reason
`uncomparable`" (MVV step 2); S8 pins the *tag-invalid* envelope
(`code`, `message`, `schema_version`, `param`) and I assume the
unevaluable one shares `code`/`message`/`schema_version` with row- and
atom-naming fields whose keys the RDR does not name.

---

## 2. The three most important internal helper functions

1. **`func (ev Evaluator) kindOf(key string) (string, bool)`** — the
   single lookup into the constructed mapping, returning the declared
   kind token and whether the key is present. It is the *only* read of
   declaration state in the package, which is what keeps A2's narrowed
   reflection assertion ("no view-typed or runtime-valued state")
   checkable by inspection. It must distinguish absent from empty —
   C1 omits empty-`Kind` keys precisely so these two do not collapse —
   so the comma-ok form, not a bare index. Name and existence as a
   distinct function: **GUESS**; C1 fixes only the behavior.

2. **`func (ev Evaluator) evalEq(atom, value string) resolve.GuardResult`
   / `evalIn(...)`** — the per-kind dispatch. `evalEq` switches
   EXHAUSTIVELY over the five-kind vocabulary with a `default` meaning
   `GuardUnevaluable`; `evalIn` splits the literal on the list separator,
   then **parses every member before comparing any** — parse-all-then-
   compare, so one unparseable member poisons the list even when the held
   value equals a member that parsed. These are two helpers, not one,
   because the poison rule has no analogue in `eq`. Splitting `eq` and
   `in` into separate helpers, and the member separator: **GUESS** (the
   RDR states the semantics, not the decomposition, and never says how an
   `in` literal encodes its members).

3. **`func parseInt(s string) (int, bool)`** — the one place `strconv.Atoi`
   is called for `eq`/`in`, wrapping it so held side and literal side
   cannot drift. C2 pins the parse to `strconv.Atoi`, hence Go's
   platform-native `int` width, and explicitly declines to fix that width
   or to state a suite verdict for a value overflowing on one target and
   not another. Overflow therefore surfaces here as `ok == false` →
   `GuardUnevaluable`. Existence as a named helper: **GUESS**.

Honourable mentions that are normative rather than helper-shaped:
`DeclaredKinds` (public, §1) and C5's canonicality check at
`internal/table/normalize.go::(*loader).atom`, which is a *load-time*
helper in a different package —

```go
// beside the existing conform(decl, operator, members) call at normalize.go:172
func (l *loader) conformCanonicalInt(decl TagDecl, site string, members []string) error
```
signature entirely **GUESS**; C5 and Phase 4 fix the rule
(`strconv.Itoa(n) == authored` after a successful `Atoi`), the site, the
refusal code, and the message, but not the function.

---

## 3. Data model across the boundary

Nothing new is **persisted**. The RDR adds no envelope field, no sixth
refusal kind, no atom field, and no hash input. Three shapes cross a
boundary:

```go
// 1. The carrier. Constructed once per loaded model, passed by
//    construction into the evaluator. Key -> kind token.
//    Vocabulary is EXACTLY five: enum | bool | int | set | scalar
//    (RDR 0003's spelling). Totality: every key in m.Tags, except
//    Kind == "" which is OMITTED. Never nil in non-test code.
map[string]string
```
Kind tokens are plain lowercase `string`, not a named enum type:
**GUESS** — C1 calls them "tokens" and C3's meta-check compares the
fixture's tokens against the vocabulary, which reads as string compare.

```go
// 2. The conformance fixture. Same type, kernel-owned, published by
//    internal/resolve so no internal/table type crosses the import
//    boundary (A5: internal/table -> internal/resolve is one-way).
//    Types comparisons BY KEY, so one key carries exactly one kind and
//    each suite case names its own key drawn from it.
map[string]string
```
Fixture contents are **GUESS** except for what C3 forces: at least one
key per kind token, with the existing `eq`/`in` cases re-keyed onto an
`enum`/`scalar` key, the five `gte` cases onto the `int` key, and the
four `contains` cases onto the `set` key, replacing today's universal
`Key: "subject"`.

```go
// 3. The verdict. UNCHANGED shape; what changes is that
//    GuardUnevaluable now occurs on inputs that previously answered
//    GuardFalse, and C4 forbids any consumer collapsing three into two.
resolve.GuardResult
```

Threading (C4), not a data shape but a lifetime rule: at the CLI,
`guardSeam(m)` builds the evaluator once per request and
`internal/cli/flow_next.go::probeRow` RECEIVES the constructed
`resolve.GuardEvaluator` rather than widening to take a `*table.Model`.
In graphlint, the evaluator is built where `matchSatisfiable(m *table.Model, …)`
holds the model and threaded down two frames to `atomAdmitsValue`.
One model, one evaluator.

The three construction sites on `main`: `internal/cli/flow_resolve.go::guardSeam`,
`internal/guard/product.go::valueSatisfies`,
`internal/graphlint/reach.go::atomAdmitsValue`. A fourth arrives when
RDR 0030 Phase 2 extracts the loader-side admitted-cell shim
(`JDR 0004 §JD-3`), which constructs over `l.model.Tags`.

`atomAdmitsValue`'s post-change return type: **GUESS**. C4 forbids the
present `return verdict != resolve.GuardFalse` (a two-valued collapse
whose two callers read it with opposite polarity) and requires
GuardUnevaluable be its own case at every consumer, but explicitly does
NOT decide what lint then does with it — so whether the function returns
`resolve.GuardResult`, a three-state lint enum, or a `(bool, bool)` is
open. I assume it returns `resolve.GuardResult` and that
`analysis.go::nodeMeetsAll` switches on it. The A4 fork on the HELD
ingress (canonicalize held values, or make the seam's match arm
byte-compare) is likewise recorded as OPEN and must be settled before
implementation — I assume neither, i.e. the design as written.

---

## 4. Top-level pseudo-code of the main operation

The main operation is `Evaluator.Evaluate` — the seam.

```
func (ev Evaluator) Evaluate(atom GuardAtom, value string) GuardResult:
    # exists never reaches the seam (0007:C1); the caller handles it.
    switch atom.Operator:

    case "lt", "lte", "gt", "gte":
        # OPERATOR-inferred, not kind-inferred. The matrix admits only
        # int here, so the mapping is NOT consulted at all — a
        # nil-mapping evaluator still answers these normally (F4).
        return compareOrdered(atom.Operator, value, atom.Literal)

    case "contains", <the §D13 set arms>:
        # Unchanged. Operator-inferred; mapping not consulted, so
        # `contains` over a set key does NOT hit the set arm below.
        return evalSetOperator(atom, value)

    case "eq", "in":
        kind, declared := ev.kindOf(atom.Key)     # nil map -> declared=false
        if not declared:
            # Two populations: the truly undeclared key (defensive —
            # A1/0007:A7 refuse it at load) and the key declared with
            # Kind "" (LIVE, and a named suite case). Also the whole
            # nil-mapping zero-value evaluator. Never fall back to
            # raw-string comparison.
            return GuardUnevaluable

        switch kind:                              # EXHAUSTIVE over five
        case "int":
            held, ok := parseInt(value)           # strconv.Atoi
            if not ok: return GuardUnevaluable    # incl. platform overflow
            members := literalMembers(atom)       # 1 for eq, n for in
            parsed := []
            for m in members:                     # PARSE ALL FIRST:
                n, ok := parseInt(m)              # one bad member poisons
                if not ok: return GuardUnevaluable
                parsed.append(n)
            return boolResult(held in parsed)     # compare PARSED values

        case "bool":
            if value not in {"true","false"}: return GuardUnevaluable
            for m in literalMembers(atom):
                if m not in {"true","false"}: return GuardUnevaluable
            return boolResult(token equality)

        case "enum", "scalar":
            # ONE shared arm — no input distinguishes them.
            return boolResult(value matches any literalMember)  # exact string

        case "set":
            # Matrix does not admit eq/in over set; defensive.
            return GuardUnevaluable

        default:
            # Unknown kind token. MUST mean unevaluable, not string
            # equality — a string-equality default would satisfy
            # enum/scalar while silently disarming both defensive arms.
            return GuardUnevaluable

    default:
        return GuardUnevaluable                   # [GUESS]
```

`literalMembers`, `compareOrdered`, `evalSetOperator`, `boolResult`, and
the final `default` are **GUESS** — the RDR fixes the verdicts for every
branch above but names none of these functions and never states how an
`in` literal carries its members or what an unrecognized operator does.
