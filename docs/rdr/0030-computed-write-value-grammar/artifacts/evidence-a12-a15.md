# cli/0030 — relocated Evidence detail for A12 and A15

Written at the Stage 7 lock (2026-09-20). The two Evidence fields below ran
32 and 48 lines against TEMPLATE.md's one-sentence-plus-stable-anchor form.
For a `foundational` record `rdr lint --locking` marks that BLOCKING, and the
tool's own fix is to relocate the balance here while the field keeps the
load-bearing anchors and a pointer. **Nothing is truncated** — every word of
both verification narratives is reproduced verbatim below; the record's
fields now carry the claim, the anchors, and this pointer.

Primary verification artifacts these narratives were derived from:
`evidence/reconcile/source-verify.md` §A12 and
`evidence/reconcile/source-verify-a15.md`.

## A12 — `internal/guard/grammar.go`'s `Evaluator` can move to `internal/resolve` unchanged

**Status**: Verified · **Method**: Source Search

`internal/guard/grammar.go` imports only `encoding/json`, `slices`, `strconv`
and `internal/resolve`, and names no `table` type, so it carries no dependency
`resolve` lacks. `internal/guard/declaration.go`'s
`::intWidth(minV, maxV int) (int, bool)` moves in the same relocation (C1) —
it takes no `table` type either — so the width rule the loader needs is
single-sourced rather than copied into `table`. It is unexported today and
both post-move callers are outside `resolve`, so it lands as the exported
`resolve.IntWidth`; verify the two callers compile against that spelling.
`::IntDomain` does NOT move: its parameter is a `table.TagDecl` and `resolve`
cannot import `table`, so it stays in `guard` as the decl-unpacking wrapper it
already is. Verify the move by compiling, and verify the wrapper by checking
`::domainSize` and `assignment.go::valueAssignments` both reach the moved core
rather than keeping the in-package re-inline of the counted loop that exists
today. `resolve` declares `::GuardEvaluator` (the interface `Evaluator`
implements, `0012:C2`) and hosts `guardcontract.go`'s conformance harness, so
the move co-locates the interface, its sole implementation (A4) and its
contract test.

**Verified at Stage 6 against `main`** (`evidence/reconcile/source-verify.md`
§A12): `grammar.go` imports exactly `encoding/json`, `slices`, `strconv` and
`internal/resolve` and names no `table` identifier; `::intWidth`'s signature
is `int`-only and its three callers are `::domainSize`, `::IntDomain` and
`assignment.go::valueAssignments`, all in `guard` and so all outside
`resolve`, which is why the move exports it; `::IntDomain(d table.TagDecl)`
and `::domainSize(d table.TagDecl)` carry the `table` type and stay;
`internal/resolve/resolve.go::GuardEvaluator` and
`guardcontract.go::TestGuardEvaluatorContract` are both already there.
`go list -deps ./internal/table` names `internal/resolve` and neither
`internal/guard` nor `internal/graphlint`, so the cycle direction is
`table → resolve` as the clause requires. The remaining obligation is the
compile itself, which Phase 2 discharges.

## A15 — a per-row graph-lint finding can carry the expanded row's suffixed identity

**Status**: Verified · **Method**: Source Search

Every per-row site sets `Rule: row.RuleID` (`internal/graphlint/groups.go`,
`analysis.go`, `coverage.go`); the package calls `Row.Identity()` nowhere
outside tests. Verify that those per-row sites all hold the `Row` (not just
its id) at the point of construction, so the suffixed form is reachable
without threading new state, and that the group-level sites using
`firstRuleID` are disjoint from them. **And enumerate every READ-BACK of
those slots**, not only the writes: `coverage.go::groupHasOverlap` joins an
emitted finding's `Rule` and `Element` against `ruleIDsOf(g)`'s bare ids to
gate the `ambiguous_match` coverage arm, so it must move to a key the suffix
does not carry — C3 fixes that key as the authored id RECOVERED from the
suffixed form (truncate at the first `#`), because the join is two-slot and
`Element`'s row carries no span to join on. Verify both slots recover, and
that `ruleIDsOf` needs no change. The search is for consumers of the field,
not just producers — the producer census is the half that misses this.

**Verified at Stage 6 against `main`**
(`evidence/reconcile/source-verify-a15.md`). EIGHT per-row producer sites,
not the three files' worth this text first implied — five in `groups.go`
(`::checkOwnedBeforeMatch`, `::emitOverlaps`, `::checkRedundantRows`,
`::checkIdempotentWrites`, `::checkVacuousAtoms`), two in `analysis.go`
(`::checkSingleValuedState`, `::checkUnreachableRules`), one in `coverage.go`
(`::emitWithholdings`) — and every one holds the `Row` as a loop variable, so
the suffixed form is reachable without new state. The four group-level sites
read `::firstRuleID` or `::bareEscapeFor` over a `guard.Group` with no `Row`
in scope, and are disjoint from the eight. `::ruleIDsOf` reads `row.RuleID`
directly and needs no change; no production caller of `Row.Identity()` exists
in the package. Recovery is exact: `normalize.go::normalizeRules` refuses a
`#` in a rule id and `load.go` refuses one in a candidate member, so
truncating at the FIRST `#` recovers the authored id for any number of suffix
elements.

The read-back census ran to TWO, not one — the second is the reason this
assumption insists on consumers over producers, since its own first pass
named only `::groupHasOverlap`. `internal/graphlint/engine.go::identityKey`,
the sort key `::sortFindings` orders on, also reads both slots. It is
verified SAFE and needs no change: it branches on `Rule` being NON-EMPTY to
pick a namespace and then uses the raw string only as a within-namespace
tie-break, so a suffixed value stays in the same bucket and the
rule-before-element ordering holds; no test pins the literal string, only
that the key forms a stable total order. It never resolves the field against
authored `[[rules]]` ids by equality, which is what C3's fence forbids — so
the fence is total over both readers and only `::groupHasOverlap` moves.
`export.go::compareEdges` reads an `Edge.Rule`, a different struct in the
graph-edge vocabulary C3 holds fixed, and is not a finding read-back at all.
