Model: claude-opus-5[1m]
Lens: 3amigo — resolution (iteration 1)

Dispositions against the `T-N` origin ledger in `consolidation.md`. Every finding
passed the grounding gate (code on `main` @ `d82f1b5`, `{RDR_RESOURCES}`, and the
RDR's own decided text) before editing.

| # | Disposition | Section touched |
| --- | --- | --- |
| T-1 count carrier | **fixed** — count pinned structurally as `Count int` on each per-identity error, per-identity, counting pre-collapse rows; A8's "If wrong" reworded so a count field is no longer framed as the alternative to collapsing | Normative Contracts (multi-breach); A8 |
| T-2 unnamed identifiers | **fixed** — `ErrEscapeShapeBreach`, `*EscapeShapeBreachError{Ref, Count}`, `Table.CheckValid() error` pinned; `CheckValid` chosen over `Validate` on the pprof precedent; stale "sharpened at Pre-Lock" deferral removed | Normative Contracts (typed error); Load-Bearing Decisions → Naming; Phase 1; MVV preamble |
| T-3 CLI scope silence | **fixed** — scenario 11 marked DEFERRED with `Code`/`Hint` literals and a committed wire carrier; Testing Strategy now states which scenarios are this RDR's Done criteria | Testing Strategy sc.11 + preamble; Normative Contracts (CLI clause) |
| T-4 traversal contract | **fixed** — aggregate pinned to exactly one level deep, elements `*EscapeShapeBreachError`, `errors.Is` holding per element, and `Resolve` returning `CheckValid`'s error verbatim (no `fmt.Errorf` wrap) | Normative Contracts (typed error); sc.9 |
| T-5 ownership overclaim | **fixed** — "no write-bearing escape plan by any producer path" qualified as staged: Phases 1–2 close the reachable population, the authored-path guarantee arrives with RDR 0002 | Consequences |
| T-6 vacuous mutation check | **fixed** — scenario 6 replaced with two named mutants and an explicit exclusion of this RDR's own tests from the oracle; records honestly that the frozen suite alone cannot detect the check's removal | Testing Strategy sc.6 |
| T-7 no oracle for "discriminates" | **fixed** — sc.5 baseline re-pointed at `a3-conformed.out` restricted to then-existing tests; the "discriminates" claim moved to sc.6 where a mutant supplies the oracle | Testing Strategy sc.5 |
| T-8 placement + frozen doc list | **fixed** — placement recorded as a cheapness preference, not an observable contract (`assemble` is pure); the REQ-1 doc list amendment explicitly authorized and added to Phase 1 | Load-Bearing Decisions → Placement / Doc-contract amendment; Phase 1 |
| T-9 Done in codebase terms only | **fixed** — "Done means" restated in user terms first, with the codebase properties as its checkable form | Testing Strategy preamble |
| T-10 author friction unweighed | **fixed** — friction weighed explicitly; notes A5's protection expires exactly when the load-time diagnostic arrives, so strictness never lands without its remedy | Consequences |
| T-11 empty-vs-nil non-discriminating | **fixed** — sc.3's discriminating assertion moved onto the input row; records that no A3 run covered the non-nil-empty variant | Testing Strategy sc.3 |
| T-12 precedence 1-of-5 | **fixed** — sc.4 made exhaustive over all five refusal kinds, matching the MVV's house style | Testing Strategy sc.4 |
| T-13 `Refused()` trap unscenario'd | **fixed** — new scenario 7b asserts the zero-`Result` shape and `Refused() == false` | Testing Strategy sc.7b |
| T-14 sc.7 comparison relation | **fixed** — "identical" defined as equal extracted `[]RowRef` + equal `Count`, explicitly not `reflect.DeepEqual` over two `errors.Join` values | Testing Strategy sc.7 |

## Net-new assumption (flag-as-you-go)

**A9 (Pending, Source Search)** — the T-3 decision to carry row identities in a
NEW `omitempty` field on `clierr.CLIError` rather than in `Detail` runs against a
shipped doc comment at `internal/cli/clierr/clierr.go::CLIError.Cause`, which
states "the wire-visible cause surface is Detail." The RDR now names the new
field as the carrier, so this is a load-bearing claim touching an **Implemented**
peer (RDR 0005) and must be verified before lock rather than absorbed. Both the
normative CLI clause and scenario 11 are marked *pending A9*. Stage 6 closes it.

No previously-`Verified` assumption was invalidated: A1–A8 are untouched by these
edits. A8's wording was clarified (its "If wrong" branch), not its verdict.

## Dispatcher-side grounding that changed a disposition

The frozen package-boundary guards (`resolve_test.go` REQ-25/26/28/36/37,
`mvv_test.go` leg 3) were read before pinning T-2's names, because Phase 1 adds
exported symbols and two imports to a package whose surface is test-frozen. All
banned lists are exact-match and name none of `CheckValid`, `ErrEscapeShapeBreach`,
or `EscapeShapeBreachError`; `errors`/`fmt` are forbidden by no import guard. Had
any collided, the naming decision would have gone the other way — so this check is
load-bearing, not ceremonial. Its corollary also fixed T-3's carrier question:
REQ-37 and REQ-28 forbid the kernel importing `encoding/json` or any CLI package,
so the wire carrier is necessarily CLI-side and the kernel cannot serialize its
own breach.

## Convergence

No open ledger entries. Fourteen fixed, none dismissed, none charted to a
successor — no finding was net-new scope; every one landed inside the contract
this RDR already owns. The edits sharpen forms the RDR had left unnamed and
correct two overclaims; none reopens a decided option, so no re-run of the lens
is owed. Iteration 1 converged.
