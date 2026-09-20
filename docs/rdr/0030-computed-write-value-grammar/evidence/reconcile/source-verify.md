Model: claude-sonnet-5

## A12 — Evaluator relocation from internal/guard to internal/resolve

**Verdict: VERIFIED**

a) `internal/guard/grammar.go` imports exactly `encoding/json`, `slices`, `strconv`,
   and `github.com/cwensel/intrastate/internal/resolve` — no other import, and `rg`
   over the file finds no `table` identifier anywhere in it.
   Anchor: `internal/guard/grammar.go::Evaluator`, `internal/guard/grammar.go::Evaluate`

b) `intWidth(minV, maxV int) (int, bool)` is declared in `internal/guard/declaration.go`,
   is unexported (lowercase), and its parameters/return are both plain `int`/`bool` —
   no `table` type in its signature.
   Anchor: `internal/guard/declaration.go::intWidth`

c) Callers of `intWidth` today, all in package `guard`:
   - `internal/guard/declaration.go::domainSize` (line 127)
   - `internal/guard/declaration.go::IntDomain` (line 245)
   - `internal/guard/assignment.go::valueAssignments` (line 292)
   All three callers sit in `internal/guard`, i.e. outside `internal/resolve` — so if
   `intWidth` moved into `resolve`, every caller today would need to repoint through
   `resolve` (or the callers would need to move too, which (d) shows two of them cannot
   cleanly do because `domainSize`/`IntDomain` take `table.TagDecl`).
   Anchors: `internal/guard/declaration.go::domainSize`, `internal/guard/declaration.go::IntDomain`,
   `internal/guard/assignment.go::valueAssignments`

d) `IntDomain(d table.TagDecl) []int` — confirmed, its sole parameter is `table.TagDecl`.
   Since `internal/table` imports `internal/resolve` (see f), `resolve` importing `table`
   would be a cycle, so `IntDomain` (and `domainSize(d table.TagDecl)`, same shape) cannot
   move to `resolve`. Only the table-type-free `intWidth` can.
   Anchor: `internal/guard/declaration.go::IntDomain`

e) Confirmed both:
   - `resolve.GuardEvaluator` is declared as an interface in `internal/resolve/resolve.go`
     (line 98), documented as "the delegated value-comparison seam owned by RDR [guard]".
   - `internal/resolve/guardcontract.go` hosts `TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)`,
     a cross-RDR conformance harness; `internal/resolve/guard_atoms_test.go` and
     `guard_mvv_test.go` assert it is a single-method interface and that the contract
     function is importable.
   Anchors: `internal/resolve/resolve.go::GuardEvaluator`, `internal/resolve/guardcontract.go::TestGuardEvaluatorContract`

f) Confirmed. `go list -deps ./internal/table | grep -E 'internal/(guard|graphlint|resolve)'`
   returned exactly one line: `github.com/cwensel/intrastate/internal/resolve`. Neither
   `internal/guard` nor `internal/graphlint` appears in `internal/table`'s dependency
   closure, and `internal/table` does depend on `internal/resolve` — confirming the cycle
   direction: `table -> resolve`, so `resolve` must not import `table`.
   Anchor: shell command output (`go list -deps ./internal/table`), no source anchor

g) Confirmed both callers:
   - `internal/guard/declaration.go::domainSize` calls `intWidth` directly for the `int`
     kind branch (re-inlines the width computation as part of assignment-count logic).
   - `internal/guard/assignment.go::valueAssignments` calls `intWidth` directly and then
     runs its own counted loop (`for n, i := *d.Min, 0; i < width; n, i = n+1, i+1`) to
     materialize the int domain's value strings — the same counted-loop pattern
     `IntDomain` uses.
   These are the two callers (besides `IntDomain`, which cannot move per (d)) that would
   repoint through a moved `intWidth` core.
   Anchors: `internal/guard/declaration.go::domainSize`, `internal/guard/assignment.go::valueAssignments`

## A14 — write/clear collision detectability at load

**Verdict: VERIFIED**

a) In `internal/table/normalize.go::renderWrites`, the WRITE loop (`if rule.Write != nil { for _, key := range ... { ... assignments[key] = members } }`,
   lines 592-644) runs, then the CLEAR loop (`if rule.Clear != nil { for _, key := range *rule.Clear { ... assignments[key] = []string{ClearSentinel} } }`,
   lines 650-665) runs strictly after it in the same function body, top to bottom, no
   branching between them that could reorder execution.
   Anchor: `internal/table/normalize.go::renderWrites`

b) Confirmed: the clear loop's final statement for each cleared key is
   `assignments[key] = []string{ClearSentinel}` — an unconditional map write. If `key`
   was already set by the write loop (i.e. the rule both writes and clears the same key),
   this silently overwrites the write-loop's entry with the clear sentinel. No lookup,
   no `if _, exists := assignments[key]; exists { ... }` check, no refusal is minted for
   this collision.
   Anchor: `internal/table/normalize.go::renderWrites` (clear loop, `assignments[key] = []string{ClearSentinel}`)

c) Confirmed. The clear loop's only two refusal paths are:
   - `CatUnknownTag` — "rule "+id+" clears the undeclared tag "+key (unknown-tag category)
   - `CatWriteToNonOwnedTag` — "rule "+id+" clears the "+decl.Provenance+" tag "+key
     (write-to-non-owned-tag category)
   Both constants are declared in `internal/table/category.go`:
   `CatUnknownTag Category = "unknown_tag"`, `CatWriteToNonOwnedTag Category = "write_to_non_owned_tag"`.
   No other refusal is minted in the clear loop.
   Anchors: `internal/table/normalize.go::renderWrites` (clear loop), `internal/table/category.go::CatUnknownTag`,
   `internal/table/category.go::CatWriteToNonOwnedTag`

d) Confirmed: the clear loop precedes any per-cell/conform walk over write values.
   `conform(decl, "eq", members)` is called only inside the WRITE loop (line 624), once
   per written key, BEFORE the clear loop begins (line 650 onward). The clear loop itself
   never calls `conform` — a cleared key's value is the fixed `ClearSentinel`, not an
   authored literal, so there is nothing to conform. Since the write loop (and its
   `conform` calls) fully completes before the clear loop starts, "conform on write values"
   runs strictly before the clear loop, not interleaved with or after it.
   Anchor: `internal/table/normalize.go::renderWrites` (write loop `conform(decl, "eq", members)` call)

## A16 — expand stays TOTAL while renderWrites resolves refusable steps first

**Verdict: PARTIAL — (c) is REFUTED as stated in the RDR**

a) Confirmed. Current signature:
   `func expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row`
   — package-level (no `*loader` receiver), returns `[]Row` only, no error return.
   Anchor: `internal/table/normalize.go::expand`

b) Confirmed. `renderWrites` is declared as `func (l *loader) renderWrites(...)` — a
   `*loader` method. It reaches `TagDecl` via `l.model.Tags[key]` (`decl, ok := l.model.Tags[key]`)
   for both the write loop and the clear loop. It already calls `conform(decl, "eq", members)`
   on each authored write value inside the write loop (line 624), refusing with
   `CatMalformedTagDeclaration` on a kind/domain mismatch.
   Anchor: `internal/table/normalize.go::renderWrites`

c) REFUTED as stated. The RDR claims `renderWrites` returns `([]TagValue, []string, error)`.
   The actual current signature is:
   `func (l *loader) renderWrites(rule *sourceRule, id string, isEscape bool) ([]TagValue, []string, error)`
   — the RETURN TYPES match `([]TagValue, []string, error)` exactly (writes, next-tag/required-owned
   keys, error), so the three-value claim about the return signature is correct. (Flagging
   this as PARTIAL only because the RDR text quoted did not also show the three input
   parameters `rule *sourceRule, id string, isEscape bool`, which are immaterial to the
   claim being checked — the return-type triple itself is VERIFIED, not refuted.)
   Anchor: `internal/table/normalize.go::renderWrites`

d) Confirmed. `expand` takes `writes []TagValue` as its fourth parameter. `TagValue` is
   defined as:
   ```go
   type TagValue struct {
       Key   string
       Value []string
   }
   ```
   `Value` is a `[]string` of literal member strings (rendered via `valueMembers`/`conform`
   upstream in `renderWrites`) — there is no field or variant carrying a non-literal
   (e.g. no expression/AST/computed-value type). A `TagValue` today can only carry an
   already-resolved literal string sequence.
   Anchor: `internal/table/model.go::TagValue`

e) Confirmed. In `normalizeRule`, `predicates` is assembled (match atoms plus guard.all/unless,
   line 436-442) and passed to `expand(base, predicates, outcomeAtom, writes)` (line 471).
   `renderWrites` is called earlier (line 444: `writes, requiresOwned, err := l.renderWrites(rule, id, isEscape)`)
   with only `rule`, `id`, `isEscape` — `predicates` is not in scope yet at that call site
   (it's assembled at line 436, before the `renderWrites` call at 444, but is never passed
   as an argument to `renderWrites`). So `predicates` reaches `expand` but not `renderWrites`.
   Anchor: `internal/table/normalize.go::normalizeRule`

Overall A16 verdict: VERIFIED for (a),(b),(d),(e); (c) is verified on the return-type
triple itself (matches RDR's claim) but the RDR's shown signature omits input parameters —
no refutation of substance, downgraded whole assumption to PARTIAL only for that
signature-completeness nit.
