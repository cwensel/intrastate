Model: Fable 5.1 (claude-fable-5-1)
scope: C1 and C3 only (C2 and §mini-check-tables read for cross-checks); iteration-1 anchors covered: the six disagreements and four GUESS clusters that the rewrite pinned as P1–P7 (site ownership, `expand` totality, crossing-record contents, collision-vs-bound order, `IntWidth` export, shim signature, lint join key).

Method: read C1, C3, C2 and the mini-check tables once through `recs inspect --select`; located the rewrite's touched ranges with `git diff -U0` on the record (hunks at C1 627/650/655/667/747–757/814–821/839–857, C3 1005–1018, authority-table rows 1043–1048; the trace/oracle/disposition tables were NOT touched); grounded P2, P6, P7 against `main` with scoped ripgrep.

## P1 — site ownership (renderWrites owns declaration-dependent work; expand mints rows)

VERDICT: PINNED

The paragraph headed "The two widened sites own different halves" enumerates the owned set exhaustively ("the admitted-cell evaluation, the `int` width check, the `#`-member and duplicate-member guards, the step/clear collision, the per-cell stepped literal and its `conform` bound check (C2)") and grounds it in a structural fact rather than preference ("`::expand` is a package-level function with no `*loader` receiver and so no `TagDecl`"). No sibling paragraph gives `expand` any of these: the earlier "placeholder the expansion replaces per cell" is consistent (expand substitutes a literal it is handed, it computes nothing), and the authority table's "Writer: `expand`'s per-cell loop" for the two carriers describes who SETS `Row.Writes`/`Row.NextTags`, which is what `expand` does today. One pre-existing, non-blocking wrinkle: the untouched `trace` table lists "emit one choice point → 5 rows" as step 3 and "C2 bound" under step 5, an assertion walk whose numbering reads as rows-before-bound; C1's normative sentence "every refusal this record mints has already fired before it is called" overrides it, and the diff shows the trace predates the rewrite, so this is not a new disagreement.

## P2 — `expand` stays total

VERDICT: PINNED

"It therefore stays TOTAL — it keeps its `[]Row` return and gains no error" plus the explicit refusal of the alternative ("a third signature change this clause does not admit") leaves no second reading. Grounded: `internal/table/normalize.go::expand` returns `[]Row`, and its only write-side work is `rows[i].NextTags = cloneTagValues(writes)` / `rows[i].Writes = cloneTagValues(writes)` per row — substituting a per-cell literal from a widened per-key record at that point needs no error path, and nothing else in the function refuses. Nothing in `normalize.go` makes P2 impossible.

## P3 — the crossing record carries resolved cells and per-cell literals

VERDICT: PINNED

"`::renderWrites` resolves the step spec into the per-key list of admitted cells and the literal each cell writes, and that — not the step magnitude — is what the widened per-key record carries, so `expand` needs neither the `TagDecl` nor `n`." The earlier sentence "The spec itself rides alongside as one more per-key record" still uses the word "spec", but the widened-sites paragraph defines the record's content explicitly and excludes the magnitude; that is a wording residue, not an under-determination.

## P4 — collision refusal precedes the per-cell walk

VERDICT: PINNED

"the collision refusal PRECEDES the per-cell walk: a key both stepped and cleared is refused before any cell is evaluated, so it never reaches C2's bound. Where a rule would trip both, the collision is what the author is told" — the order and the observable consequence are both stated. C2 does not claim precedence for its own refusal, and the disposition table lists the collision as a distinct `malformed_tag_declaration` row, so no sibling contradicts it. (The order of the collision relative to the grammar/kind arm for a rule that is wrong in both ways is not stated; that is outside this pin, pre-existing, and same-category.)

## P5 — `resolve.IntWidth` exported

VERDICT: PINNED

C1 states the name and signature ("EXPORTED at its new home as `resolve.IntWidth(minV, maxV int) (int, bool)`") and the reason ("left lowercase the moved core is unreachable from both callers"). The authority table's width row was updated in the same rewrite to "the EXPORTED `resolve.IntWidth`" and the sibling `IntDomain` row names `resolve.IntWidth` for the loader; the `guard` wrapper/`domainSize`/`valueAssignments` repoint sentence is consistent. No paragraph still says `intWidth` at the new home.

## P6 — admits shim signature

VERDICT: PINNED

"the shim takes a `resolve.GuardAtom` and the held value as a STRING — the cell, never the declared kind — and returns `resolve.GuardResult`, whose members are `GuardTrue`, `GuardFalse` and `GuardUnevaluable`." Grounded: `internal/guard/grammar.go::Evaluator.Evaluate` is `func (Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult` and the three members are already in `internal/resolve` (`guardcontract.go`), so "the relocation carries no retyping" is true of `main`. The authority table's UNDECIDED row (loader EXCLUDES) agrees with C1's three-valued fence.

## P7 — lint read-back joins on the authored id recovered from the suffixed form

VERDICT: PINNED

C3 fixes the key ("the AUTHORED ID RECOVERED from the suffixed form — the published identity truncated at the first `#`"), rejects the alternative by name ("The finding's own `Span` is NOT the join key"), and gives the structural reason (`graph-overlap` carries no right-hand span). The authority table's lint-naming row was rewritten to the same shape. Grounded: `internal/graphlint/coverage.go::groupHasOverlap` joins `f.Rule` and `f.Element` against `ruleIDsOf(g)` (bare `row.RuleID`); `internal/graphlint/groups.go` emits `graph-overlap` with `Element: right.RuleID` and only `Span: rows[i].SourceLocator`, and `graph-redundant-row` likewise — so truncating both slots at the first `#` before the `slices.Contains` leaves `ruleIDsOf` untouched, as C3 claims. The truncation is exact because `internal/table/normalize.go` refuses a rule id containing `suffixSep` (`"#"`). `graphlint` has no non-test `Row.Identity()` caller, matching C3. `engine.go::identityKey` reads `f.Rule` only as a sort tie-break, not a join.

## Grounding check (source repo, `main`)

- P7 implementable without changing `ruleIDsOf`: CONFIRMED — `internal/graphlint/coverage.go::groupHasOverlap`, `::ruleIDsOf`, `::firstRuleID`; emission sites `internal/graphlint/groups.go` (overlap, redundant-row), `analysis.go`, `coverage.go` (`Span: row.SourceLocator` per row).
- P2 possible: CONFIRMED — `internal/table/normalize.go::expand` (`[]Row`, no error), `::renderWrites` (`([]TagValue, []string, error)`), call site in `::normalizeRule` assembling `predicates` immediately before `renderWrites` and handing it only to `expand`, exactly as C1 describes.
- P6 surface on `main`: CONFIRMED — `internal/guard/grammar.go::Evaluator.Evaluate`.

## New silences opened by the rewrite

none
