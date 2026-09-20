model: claude-sonnet-5
scope: delta re-run (iter-2) — 0012:C1,C2,C3,C4,C5,S4

## 1. Go declarations

```go
package guard // internal/guard

// NewEvaluator is the sole non-test constructor. kinds maps tag key ->
// declared kind token (enum|bool|int|set|scalar). A nil-mapping
// Evaluator{} (zero-value literal, var, new, or embedded field) answers
// GuardUnevaluable for every eq/in atom.
func NewEvaluator(kinds map[string]string) Evaluator

type Evaluator struct {
        kinds map[string]string // unexported; holds nothing else
}

func (e Evaluator) Evaluate(atom GuardAtom, value string) GuardResult // 0007:C1, unchanged

// DeclaredKinds is the single exported producer of the kind mapping.
// Total over m.Tags; a Kind == "" entry is OMITTED, not mapped to "".
// Nil model -> empty non-nil map, never panics.
func DeclaredKinds(m *table.Model) map[string]string

package resolve // internal/resolve

type GuardResult int

const (
        GuardFalse GuardResult = iota
        GuardTrue
        GuardUnevaluable
)

type GuardAtom struct {
        Key      string
        Operator string
        Literal  string
        Block    resolve.Block // BlockAll | BlockUnless | BlockMatch
}

// ContractKinds publishes the kernel-owned declaration fixture as a
// fresh map each call (not a package-level var).
func ContractKinds() map[string]string

// TestGuardEvaluatorContract widens to take a seam constructor, not a
// constructed seam, so a caller cannot compile against the wrong kinds.
func TestGuardEvaluatorContract(t *testing.T,
        newSeam func(kinds map[string]string) GuardEvaluator)

// conformingContractSeam (guard_mvv_test.go) gains a kinds field and a
// matching constructor to satisfy newSeam.
type conformingContractSeam struct {
        kinds map[string]string
        // ... existing fields
}

package graphlint // internal/graphlint

// atomAdmitsValue is renamed to reflect its widened three-valued return;
// callers (reach.go, analysis.go::nodeMeetsAll) must handle
// GuardUnevaluable as its own case, not collapse via != GuardFalse /
// == GuardTrue.
func atomVerdictForValue(/* unchanged params */) resolve.GuardResult
```

## 2. Step order of the eq/in comparison at the seam

1. Look up `atom.Key` in the constructed `kinds` map (nil map -> miss).
2. Key absent from the mapping (including a declared key whose `Kind`
   was the empty string, since C1 omits those) -> GuardUnevaluable.
   Dispatch over the found kind token is exhaustive across the five-kind
   vocabulary; any unrecognized token also falls through `default` ->
   GuardUnevaluable. The kind lookup/dispatch runs BEFORE any literal or
   held-value parsing — an unknown/absent kind never reaches step 3.
3. Per found kind, dispatch:
   - `int`: parse held value and literal (each `in` member) via
     `strconv.Atoi`. Parse-ALL-then-compare (not per-member
     parse-and-compare). Any unparseable held value, or any single
     unparseable `in` member, -> GuardUnevaluable — even if the held
     value equals a member that did parse. Compare over parsed values.
   - `bool`: held value and literal must each be exactly `"true"` or
     `"false"`; any other held-value spelling -> GuardUnevaluable;
     otherwise token equality.
   - `enum` | `scalar`: exact string equality (share one dispatch arm).
   - `set`: eq/in not admitted by the operator/kind matrix ->
     GuardUnevaluable (defensive; matrix rejection at load is primary
     guard).
4. Unparseable literal is a per-kind obligation folded inside step 3's
   `int`/`bool` arms (not a pre-dispatch gate) -> GuardUnevaluable.
5. Return the verdict (GuardFalse | GuardTrue | GuardUnevaluable); no
   aggregation across atoms happens at the seam (that's
   `evaluateAtoms`/`kleeneAnd`, out of scope for C1-C3).

## 3. Determinacy per anchor

0012:C1 — DETERMINATE. `NewEvaluator(kinds map[string]string) Evaluator`,
`DeclaredKinds(m *table.Model) map[string]string`, both homed in
internal/guard, nil-mapping/nil-model dispositions, TOTAL-and-OMIT
semantics for empty Kind — all explicit with signatures and package
homes.

0012:C2 — DETERMINATE. Full dispatch table, parse-all-then-compare
rule, GuardResult iota order (fixed, not re-derivable), step order
(kind lookup before literal parse) are all explicit. The `int` width
(`strconv.Atoi`, platform-native, unpinned) is explicitly flagged as
open but the clause tells the implementer exactly what to do about it
(nothing — out of scope, suite pins in-range only), so it does not
block implementation.

0012:C3 — DETERMINATE. `ContractKinds() map[string]string` and the
widened `TestGuardEvaluatorContract(t, newSeam func(map[string]string)
GuardEvaluator)` signatures are exact; the three migrating call sites
are named with file:line; re-keying assignments (enum/scalar keeps
"subject"-style keys, gte->int, contains->set) are explicit.

0012:C4 — DETERMINATE. Three construction sites named with exact
file/function; CLI threading (`guardSeam`/`probeRow`) and graphlint
threading (`atomAdmitsValue`/`matchSatisfiable`) both specified as
"construct once, thread down." The rename `atomAdmitsValue` ->
`atomVerdictForValue` and its new `resolve.GuardResult` return type are
explicit. GUESS: the clause does not spell the renamed function's
parameter list verbatim (only says the bool return widens to
GuardResult) — I reconstructed it as "unchanged params" by inference
from "widens... rather than to a lint-local enum," which is a safe but
unverified guess since the surrounding signature isn't quoted.

0012:C5 — DETERMINATE. Canonicalization rule (`strconv.Itoa(n) ==
authored string`), exact ingress site (`(*loader).atom`, calling
`conform(decl, operator, members)` at normalize.go:172), scope
(predicate ingress only, CLI explicitly out), error code reuse
(`malformed_predicate_atom`), and message shape
(`<site>: "00" is not the canonical spelling of int tag n; write 0`)
are all explicit enough to implement without invention.

0012:S4 — DETERMINATE. The scenario specifies: assertion form (matches
composite-literal, var, new, and embedded-field zero values of
guard.Evaluator outside NewEvaluator's body), implementation technique
(go/ast or golang.org/x/tools/go/packages pass over non-test files, not
regexp), location (a Go test under internal/guard), and wiring (runs
under `make check` via the existing `test` leg, no separate make
target). No invention required to write this check.
