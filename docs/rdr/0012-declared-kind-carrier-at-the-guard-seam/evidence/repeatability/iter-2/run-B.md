model: claude-fable-5-1
scope: delta re-run (iter-2) — 0012:C1,C2,C3,C4,C5,S4

---

## 1. Go declarations the six clauses fix

```go
// ---- internal/resolve (read from main; C2 names the shapes, amends none) ----

type GuardResult int

const (
	GuardFalse GuardResult = iota // zero value
	GuardTrue
	GuardUnevaluable
)

type Block string // BlockAll | BlockUnless | BlockMatch; aliased by internal/table

type GuardAtom struct {
	Key      string
	Operator string // bare string; only OpExists has a constant
	Literal  string // `in` list travels as a compact sorted/deduped JSON array
	Block    Block
}

const OpExists = "exists"

type GuardEvaluator interface {
	Evaluate(atom GuardAtom, value string) GuardResult
}

// ---- internal/resolve, non-test source (C3; new) ----

// Fresh copy per call; key -> kind token. Values are exactly
// {enum,bool,int,set,scalar} (meta-check).
func ContractKinds() map[string]string

// Widened from (t, seam GuardEvaluator). Suite builds the seam over
// ContractKinds() itself.
func TestGuardEvaluatorContract(t *testing.T,
	newSeam func(kinds map[string]string) GuardEvaluator)

// ---- internal/resolve/guard_mvv_test.go (C3; amended) ----

type conformingContractSeam struct {
	kinds map[string]string // was struct{}
}

func newConformingContractSeam(kinds map[string]string) resolve.GuardEvaluator

// recordingSeam unchanged: keeps delegating to `inner`, no kinds field.

// Re-typed pin inside TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface:
var fn func(*testing.T, func(map[string]string) resolve.GuardEvaluator) = resolve.TestGuardEvaluatorContract

// ---- internal/guard (C1) ----

type Evaluator struct {
	kinds map[string]string // unexported; the ONLY state
}

func NewEvaluator(kinds map[string]string) Evaluator

// Total over m.Tags; keys whose decl.Kind == "" are OMITTED.
// nil model -> empty non-nil map (never panics).
func DeclaredKinds(m *table.Model) map[string]string

// Unchanged signature (0007:C1). Nil/empty kinds -> every eq/in answers
// GuardUnevaluable. Never panics.
func (e Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult

// ---- internal/guard/guard_evaluator_0003_test.go (C3; migrated) ----
// resolve.TestGuardEvaluatorContract(t, <adapter over guard.NewEvaluator>)
// -- see C3 GUESS in §3: guard.NewEvaluator's return type is Evaluator,
//    which does not assign to func(map[string]string) resolve.GuardEvaluator.

// ---- internal/guard (S4; new test) ----
// A go/ast or golang.org/x/tools/go/packages pass over non-test files in
// the module failing on any composite literal, var decl, new(), or struct
// field of type guard.Evaluator outside NewEvaluator's body.
func TestNoZeroValueEvaluatorConstruction(t *testing.T) // name is latitude

// ---- internal/cli (C4) ----
// guardSeam takes the model, called once per request, returns the
// constructed evaluator; probeRow receives it rather than a *table.Model.
func guardSeam(m *table.Model) resolve.GuardEvaluator   // return type: Evaluator or the interface — latitude
func probeRow(ev resolve.GuardEvaluator, /* existing params */) // widened by the evaluator, NOT by *table.Model

// ---- internal/guard/product.go (C4) ----
// valueSatisfies constructs via NewEvaluator(DeclaredKinds(m)) over the
// model whose rows it evaluates.

// ---- internal/graphlint/reach.go (C4) ----
// Renamed from atomAdmitsValue; bool return widens to the kernel verdict.
// Evaluator constructed once where *table.Model is in scope
// (matchSatisfiable) and threaded two frames down.
func atomVerdictForValue(ev resolve.GuardEvaluator, /* existing params */) resolve.GuardResult
// Both callers (reach.go site and analysis.go::nodeMeetsAll) switch on
// all three verdicts; `!= GuardFalse` / `== GuardTrue` collapses forbidden.

// ---- internal/table/normalize.go::(*loader).atom (C5) ----
// Beside the existing conform(decl, operator, members) call at :172, for
// decl.Kind == "int": for each member m,
//   n, err := strconv.Atoi(m); if err == nil && strconv.Itoa(n) != m -> refuse
// Refusal: same error shape conform already returns; code
// malformed_predicate_atom; message
//   <site>: "00" is not the canonical spelling of int tag n; write 0
// conformKind unchanged. CLI (--tag/--write via table.ConformValue) untouched.
```

---

## 2. Step order of the eq/in comparison at the seam

1. Operator dispatch on `atom.Operator` (bare string). `lt|lte|gt|gte`
   take the existing operator-inferred integer parse; `contains` and the
   §D13 set arms are unchanged; the kind mapping is NOT consulted for
   any of these. `exists` and any unrecognized operator answer
   `GuardUnevaluable` — no panic. Only `eq` and `in` continue below.
2. Kind lookup FIRST: `kind, ok := e.kinds[atom.Key]`. A nil map, an
   absent key, or a key declared with empty `Kind` (omitted by
   `DeclaredKinds`) all fall to the switch `default` →
   `GuardUnevaluable`. The literal is NOT parsed at this point.
3. Exhaustive switch over the five tokens; `default` →
   `GuardUnevaluable`:
   - `int`: (a) `strconv.Atoi(value)`; failure (including overflow on
     the target) → `GuardUnevaluable`. (b) Decode the literal — for
     `eq` the single string, for `in` every member of the JSON array —
     with `strconv.Atoi`; ANY failure → `GuardUnevaluable`
     (parse-all-then-compare; one bad member poisons the list, even if
     the held value equals a member that parsed). (c) Compare parsed
     ints: `eq` → equal ? `GuardTrue` : `GuardFalse`; `in` → any member
     equal ? `GuardTrue` : `GuardFalse`. Sub-order of (a) vs (b) does
     not change the verdict (both failures yield the same result).
   - `bool`: held value must be exactly `true` or `false`, else
     `GuardUnevaluable` (`"1"`, `"0"`, `"True"`, `"yes"` rejected);
     literal / each member likewise, else `GuardUnevaluable`; then token
     equality (`eq`) or membership (`in`).
   - `enum`, `scalar` (one shared arm): exact string equality /
     membership; no parse can fail.
   - `set`: `GuardUnevaluable` (defensive; matrix rejects at load).
4. Return the verdict. Aggregation across a row's atoms is untouched
   (`evaluateAtoms`, no short-circuit, strong-Kleene AND via
   `kleeneAnd`; `GuardFalse` dominates `GuardUnevaluable`).

---

## 3. Per-anchor determinacy

**0012:C1 — DETERMINATE.** `NewEvaluator(map[string]string) Evaluator`,
exported `DeclaredKinds(*table.Model) map[string]string` homed in
`internal/guard`, total over `m.Tags` with empty-Kind keys omitted,
nil model → empty non-nil map, unexported single field, nil-mapping →
`GuardUnevaluable` on every `eq`/`in`, `Evaluate` signature and atom
field types unchanged. Nothing to invent. (Receiver kind for
`Evaluate` is read from `main`, not restated — acceptable.)

**0012:C2 — DETERMINATE.** Five arms plus `default`, verdict per arm,
parse-all-then-compare for `in`, `strconv.Atoi` width disposition,
bool token set, enum/scalar shared arm, set/absent/unknown → default,
operator dispatch before kind lookup, kind lookup before literal parse,
no-panic posture, `in` member encoding pinned by the existing fixture,
`GuardResult` iota order preserved. §2 above was written without a
guess.

**0012:C3 — GUESS.** Everything about the suite is fixed (constructor
form of `ContractKinds`, `newSeam` parameter, per-case keys and their
routing, the meta-check, the three migrating call sites, the re-typed
`var fn` pin, `conformingContractSeam` gaining `kinds` + a constructor,
`recordingSeam` unchanged). The one silence: the clause instructs the
guard-package site to migrate to
`TestGuardEvaluatorContract(t, guard.NewEvaluator)`, but Go function
types are invariant in their return type, so
`func(map[string]string) guard.Evaluator` does NOT assign to
`func(map[string]string) resolve.GuardEvaluator` and that line will not
compile with C1's `NewEvaluator` return type. The implementer must
invent either a closure adapter at the call site
(`func(k map[string]string) resolve.GuardEvaluator { return guard.NewEvaluator(k) }`)
or change `NewEvaluator`'s return type to the interface (which C1 fixes
as `Evaluator`). I would choose the closure adapter, as it leaves C1
intact; the record does not say which.

**0012:C4 — DETERMINATE** for what it claims to fix: the universal
"one model, one evaluator", the three present sites, thread-down (not
rebuild) at `probeRow` and `atomAdmitsValue`, three-valued consumption
at both graphlint callers, return type widened to `resolve.GuardResult`
(not a lint enum or a pair), the rename to `atomVerdictForValue`, and
the explicit non-reach into aggregation and into the not-yet-present
0030 shim. One declared deferral remains that an implementer must
resolve to ship — recorded as an extra GUESS row below, not as a
silence in the clause.

**0012:C5 — DETERMINATE.** Rule (`Itoa(Atoi(s)) == s`), scope (predicate
ingress only: guard + match atom literals over `int`-declared tags, all
operators, including the bare-float `-0` path), home (`(*loader).atom`
beside the existing `conform` call, `conformKind` untouched), CLI
explicitly out of scope, refusal code `malformed_predicate_atom`,
same error shape/path as `conform`, message shape given with a worked
example, no wider normalization claimed. Nothing to invent.

**0012:S4 — DETERMINATE.** Typed `go/ast`/`go/packages` pass (not a
regexp) over non-test files in the module, failing on composite
literal, `var`, `new`, or struct field of type `guard.Evaluator`
outside `NewEvaluator`'s body; lives as a Go test under
`internal/guard`; runs via the existing `make test` leg of `make
check`. Minor unlisted forms (a named result parameter of type
`Evaluator`, a zero-appended slice element) are edge cases the same
typed pass can cover by "any zero-valued expression of that type"; I
do not count them as a silence.

Extra GUESS rows (declared deferrals an implementer must still resolve):

- 0012:C4-lint-policy — what graphlint DOES with a `GuardUnevaluable`
  atom (report a finding vs over-approximate as admissible and stay
  silent) is stated as not decided; the implementer must pick one and
  apply it identically at `reach.go` and `nodeMeetsAll`. I would
  over-approximate (admit) and emit no finding, preserving today's
  green-lint behaviour while keeping the reading explicit.
