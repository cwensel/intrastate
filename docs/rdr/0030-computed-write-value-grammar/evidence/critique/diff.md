Model: claude-fable-5

# Critique barrier diff — cli/0030

Inputs: `critique.md` (pass A, claude-opus-5, C-1..C-20) and
`critique-modelB.md` (pass B, claude-sonnet-5, C-1..C-10). Ledger ids are
per file; rows are matched by RDR passage and claim, never by id. Every
source-grounded claim below was checked against the tree at the time of
the diff; where a check failed the row is kept with the correction in the
`Grounded` column.

## 1. Convergent findings

| Passage | A row | B row | The shared claim (one line) | Both-grounded? |
|---|---|---|---|---|
| `0030:A7`, `0030:§performance-expectations` | C-14, C-15 | C-1 | Load is linear and takes no ceiling; lint is quartic (116 s at 200 rows, >600 s at 1,000) and the record assigns the cliff to 0013 — an ordinary `int` counter (`max` in the hundreds) yields a model that loads in ms and lint cannot finish. | Yes — both read the record's own numbers; A adds the package-graph point (A7: `go list -deps ./internal/table` reaches neither `guard` nor `graphlint`, confirmed) that nothing between load and lint can stop it. Strongest convergence in the set. |
| `0030:§consequences` (enum residual) | C-8 | C-3 | The `enum` arm says the value once but not the cap: excluding the terminal member needs `in = [<every member but the last>]`, re-typed on every domain change — the Problem Statement's own grievance. | Yes on the defect (record text). **Both symptoms are wrong** — see §3 contradiction 2. |
| `0030:A2`, `0030:§consequences`, `0030:C3` | C-20 | C-6 | Authored `enum` `domain` order becomes semantic the moment a rule steps it; reordering silently changes the flow; the record dispositions this to `docs/model-authoring.md` with no lint. | Yes — C3 text says "with no lint to catch it" verbatim; A2 verifies order survives load. Both models independently name it the assumption most likely to cause an incident. |
| `0030:A10`–`0030:A14` (+ A5, A6 per A) | C-11 | C-4 | Five (B) or six (A) Critical Assumptions are `Pending`, four marked "new at this pre-lock pass", and C1 states subsumption, both-carriers, bound ordering and step+clear as settled normative prose that depends on them; A10's own "If wrong" reopens C1's expansion half. | Yes — statuses confirmed in the record. A's count of six is correct (A5, A6 are also `Pending — resolves at the MVV`). |
| `0030:A12`, `0030:§phase-2-expand-into-cells` | C-10 | C-5 | Phase 2 bundles a cross-package relocation of `guard/grammar.go::Evaluator` + `declaration.go::intWidth` into `internal/resolve`, repointing `guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`, `cli/flow_resolve.go::guardSeam` and ~25 test refs across seven files, under an Approach that says "one choice-point kind in `expand`"; A12 is `Pending`, verified by "it compiles". A adds: RDR 0012 is `Draft` and lands on the same file, unsequenced. | Yes — all three call sites confirmed (`product.go:541`, `reach.go:502`, `flow_resolve.go:561`); seven `internal/guard/*_test.go` files carry 25 `Evaluator` refs; 0012 `Status: Draft` confirmed. |
| `0030:A5`, `0030:C1` (`unless` not consulted), `0030:S5b` | C-6 | C-7 | `unless tier = "large"` — the intuitive exclusion — is inert at admission; the dead row surfaces only as advisory `graph-redundant-row`, never a refusal or blocking finding. | Partially. B's framing (advisory-only) holds. A's grounding ("an `unless`-carrying row is the canonical unprojectable row … the advisory does not fire") is **wrong against source** — see §3 contradiction 1. |
| `0030:A8`, `0030:C3`, `0030:§decision-rationale`, `0030:F6` | C-19 | C-2 | `flow next`/`flow resolve` publish the bare authored `rule`; per-row identity is on `dump`/`graph` only and disclaimed as a "possible successor record"; an on-call reader cannot tell which cell fired. | Yes — A8 is `Refuted` and narrowed in the record itself; `rowByID` first-match join confirmed as disclosed. Disclosed trade-off, but both passes call it the wrong surface to give up. |
| `0030:§proportionality` | C-13 | C-10 | The gate's "sole author of at most one independent load-bearing contract" is failed: a grammar, an expansion mechanism, a load-time conformance check, an identity-ordering rule, and a shared-evaluator relocation ride one contract id. | Yes — record structure. Gate text is unfilled template (expected pre-lock), so this is a prediction of what Finalize will hit. |

## 2. Divergent findings

| Passage | Which pass | The claim | Why the other pass plausibly missed it (scope, depth, or it is wrong) | Keep / drop |
|---|---|---|---|---|
| `0030:C3` | A (C-4) | C3 lists `lint` among the readers, then draws a ROW/RULE split covering only `dump`, `graph`, `flow next`, `flow resolve`. Every `graphlint` finding populates `Rule` from `row.RuleID`; none calls `Row.Identity()`. A stepped ladder with one bad cell emits findings all reading `rule: "retry"`. | B stayed at record level and did not open `graphlint`. Confirmed at source: `groups.go:80,207,318,363,399`, `analysis.go:213,567`, `coverage.go:330` set `Rule: row.RuleID`; `coverage.go:91,140` use `firstRuleID(g)`; zero `Identity()` calls in the package. `Finding.Fingerprint` IS per-row (hashes `Row.Atoms`) but is opaque; `Span` is the shared `SourceLocator`. This is the sharpest unique finding in either pass. | **Keep** |
| `0030:S1` | A (C-5) | The MVV oracle projects to `(code, key, dimension, class, reason)` and drops `rule`, `span`, `element`, `fingerprint` — exactly the fields C-4 breaks — so the MVV cannot observe the lint-attribution defect. | Follows from C-4; B did not have C-4. Record text confirms the projection and its rationale. | **Keep** |
| `0030:C1`, `0030:§phase-2-expand-into-cells` | A (C-2, plus residual of C-1) | `expand` is a package-level function `expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row` with no `*loader` receiver and no `TagDecl`; Phase 2 assigns admitted-cell evaluation, `intWidth`, per-cell `conform` and the `#`/duplicate guards to it. "Widening that pair is the one signature change" understates the plumbing. | B read the Approach, not `normalize.go`. Signature confirmed at `normalize.go:730`; `renderWrites` (`:580`) is the `*loader` method that has the decl. Correction to A: threading the model through `expand` is one option, not forced — `renderWrites` can enumerate the cells and hand them through the widened per-key record. Either way the record does not say which, and Phase 2's text implies `expand` does work it cannot do unaided. | **Keep** (with correction) |
| `0030:C2` | A (C-9) | The bound refusal fires on a declaration-level interaction (tightening `max`) naming a rule the author did not touch; the prescribed remedy (positive-atom exclusion) then hands them a `graph-coverage-gap` — two errors from one edit. | B folded this into its enum row. Record-grounded (C2 + §consequences "tightening `max` … refuses at load"). Adjudication sharpens it: `CodeCoverageGap` is in `blockingCodes` (`taxonomy.go:85-93`), yet C1 calls it "visible, non-blocking and actionable" — the record's own remedy chain ends in a blocking finding it describes as advisory. | **Keep** (sharpened) |
| `0030:§consequences` (round-trip) | A (C-16) | The motivating journey is migrating an unrolled ladder to step form; the record declines to extend `0002:§round-trip-inverse-invariants` and ships no equivalence aid — the MVV has the oracle, the user does not. | B's premortem is post-migration; it never asks how the author gets there. Record-grounded. | **Keep** |
| `0030:C2` (detail string) | A (C-18) | The four-part refusal (rule, tag, cell, result) plus the "unconsulted `unless`" clause ride `Failure.Detail` as free text; tests assert by substring position; no structured fields because C3 forbids envelope members. | B did not read C2's mechanics. Source: `category.go:201` `Failure` carries `Category`, `Detail`, `Offending`, `Remedy`, `Rule`, `Line` — the record's "Category, Detail, Offending and Rule" omits `Remedy` and `Line`, a minor inaccuracy the fix should correct. The disposition is disclosed and reasoned in C2; keep as a flagged trade-off. | **Keep** |
| `0030:S5b` | A (C-7) | S5b asserts "EXACTLY ONE advisory" and then asserts its own projectability precondition — proving the outcome on an engineered fixture, not on authored ladders. | B accepted S5b. Record-grounded, but A's premise that authored ladders commonly fail projectability is wrong (see §3.1): non-projectable means an optional-key atom or a refused product, not an `unless`. Downgraded. | **Keep** (downgraded) |
| `0030:§problem-statement`, `0030:§decision-rationale` row 6 | B (C-8) | `dump` and `graph` still show `retry#0 … retry#4`; the "review legibility" win is source-only, halving rather than removing the cost the Problem Statement says lands twice. | A treated the `dump`/`graph` behaviour as the stable good surface. Record discloses this ("the cost Approach 1 pays is row 6's second half"). Low severity, accepted precedent from `in` expansion. | **Keep** |
| `0030:A4`, `0030:C1` (third arm) | B (C-9) | The admits filter's "UNDECIDED arm is unreachable" rests on the frozen operator/kind matrix; the record itself names widening that matrix as the natural successor, after which the soundness argument silently lapses. | A read A4 as settled. Record-grounded, but C1 chose the safe collapse (exclude → coverage gap), so a future reachable arm degrades to a visible gap rather than a wrong admit. Forward-fragility, low. | **Keep** (low) |
| `0030:C1` (placeholder / `conform`) | A (C-1) | `conform(decl, "eq", members)` runs inside `renderWrites` on the authored value; a placeholder either conforms (indistinguishable from a literal) or forces a branch that skips conformance so an out-of-domain step loads clean. | **Wrong against the record.** The audit table says the step form is "intercepted before `valueMembers`" on the write path, and C2 places the conformance check on the EXPANDED per-cell literal, not on the placeholder — so no cell escapes `conform`. `conform` call confirmed at `normalize.go:624` (A cites `:623`). The one residual — the placeholder literal is unnamed — is folded into D-9. | **Drop** |
| `0030:D-identity` | A (C-3) | Step points "order by KEY alone" while other candidates use `compareAtoms`' 4-tuple; a mixed comparator over one `slices.SortFunc` is not a total order, so `Identity()` varies by build. | **Wrong against the record and source.** D-identity states the comparator: key-first for every candidate ("the same key comparison places a step point among the OTHER candidates"), and `compareAtoms`' first component IS key (`normalize.go:263-266`), so key-only comparison agrees with the tuple wherever keys differ; subsumption removes the same-key match `in`, `guard.all` atoms never expand (`0002:C13`, match-block only), and `RecognizedTagKey` is never owned — so no candidate shares a step point's key and no tie exists for instability to act on. Residual worth the fix half's note: no scenario in S1–S10 pins the mixed `step`+`in` suffix order D-identity's `retry#<attempt>#<mode>` example promises. | **Drop** |
| `0030:§finalization-gate` | A (C-12) | The gate is unfilled template text, so the Assumption Verification rule that would catch C-11 has not run. | Not a defect at this stage: the gate header says it is completed "before marking this RDR as Final" and the record is `Draft`; the substantive half (Pending assumptions under settled prose) is D-5. | **Drop** (folded into D-5) |
| `0030:C1` (narrows `0002:C13`) | A (C-17) | A Draft narrows a locked record's clause inside its own prose; nothing in 0002 points at 0030. | **Answered by the record.** `0030:§metadata` Overrides names `0002-…:C13 — narrows the suffix-iff … to the `in` expansion it was written about`. Records are never amended in this workflow, so the Overrides field is the pointer mechanism; 0002 cannot be edited to point forward. | **Drop** |

## 3. Contradictions

**3.1 — `0030:A5` / `0030:S5b`: does the dead row's advisory fire?**
A C-6: "an `unless`-carrying row is the canonical unprojectable row —
`decidableAccepted`'s own comment says so … on the common shape it does not
fire at all." B C-7: it fires, but only as an advisory the author can
ignore.

Winner: **B**. `checkRedundantRows` (`internal/graphlint/groups.go:299`)
calls `guard.AcceptedAssignments` directly, not `decidableAccepted`.
`AcceptedAssignments` → `acceptedIn` (`internal/guard/product.go:568-640`)
handles `unless` by collecting the block's denotations and subtracting
their intersection (`:609-635`); it returns the unprojectable set only for
a can-refuse row, an undecidable product, an atom over an undecidable key,
or an unprojectable denotation — never for the presence of `unless`.
`decidableAccepted`'s `unless` → undecided branch (`groups.go:278-283`) is
on its FALLBACK path, reached only when `AcceptedAssignments` was already
unprojectable, and it serves the overlap check, not the redundant-row one.
So on a fully-declared ladder the dead row's accepted set is `{cell} minus
{cell}` = empty, a proper subset of the literal sibling's, and the advisory
fires. A's residual that survives: the gate is still same-group, same-`Kind()`,
and requires a projectable sibling that claims the cell; where no row claims
it the signal is instead a **blocking** `graph-coverage-gap`.

**3.2 — `0030:§consequences`: what happens when a fourth tier is added?**
A C-8: the model "refuses at load with a bound error pointing at a rule
they did not touch." B C-3 / Incident 1: the step "silently began admitting
a cell into `critical`" and `graph-coverage-gap` "did not fire because the
cell was, technically, claimed."

Winner: **neither**. Per C1 an admitted cell must satisfy the conjunction
of the rule's positive atoms; with `in = ["small","medium"]` unchanged after
`critical` is appended, `large` stays excluded, no cell steps into
`critical`, and no stepped value leaves the domain — so no C2 refusal (A is
wrong) and no widening (B is wrong). `critical` is admitted by no step cell
and claimed by no row, which §consequences already names: "widening it
leaves new cells to the group's coverage proof (`graph-coverage-gap`)". The
actual symptom is a lint `graph-coverage-gap` on `critical` — blocking per
`taxonomy.go:89` — and the ladder never escalating past `large` until the
`in` list is re-typed. The shared defect (cap must be re-listed) stands;
the fix half must not carry either pass's symptom into the record.

**3.3 — `0030:C1` / `0030:§illustrative-code`: C1 calls `graph-coverage-gap`
"non-blocking".** Not a pass-vs-pass contradiction but surfaced by 3.1 and
3.2: C1's third-arm paragraph says excluding "leaves the cell to the outcome
group's ordinary `graph-coverage-gap`, which is visible, non-blocking and
actionable", while A5 lists `CodeCoverageGap` as blocking and
`taxonomy.go:85-93` puts it in `blockingCodes`. The record contradicts
itself and source; folded into D-11.

No other pair of rows makes incompatible claims at one passage. A C-11's
"six Pending" vs B C-4's "five" is a counting difference (A includes A5,
A6), both correct.

## 4. Merged ledger

| ID | RDR passage | Failure mode | Symptom user sees | Source (A/B) | Grounded |
|---|---|---|---|---|---|
| D-1 | `0030:C3` | C3 names `lint` as a reader and then omits it from the ROW/RULE split; every `graphlint` finding sets `Rule` from `row.RuleID`, none from `Identity()`, and the record is silent on it. The MVV cannot see it (D-10). | N findings all reading `rule: "retry"` on a stepped ladder with one bad cell; the unrolled table named `retry-3`. The record's headline benefit inverts on the review tool. | A | source:`internal/graphlint/groups.go:80,207,318,363,399`; `analysis.go:213,567`; `coverage.go:91,140,330` — all `RuleID`/`firstRuleID`, zero `Identity()` callers in package |
| D-2 | `0030:A7`, `0030:§performance-expectations` | No load-time ceiling; lint is quartic and 0013's scope; package graph guarantees nothing between load and lint can stop a wide step. | `attempt = { step = 1 }` over `max = 200`: loads in ms, `lint` never returns or reports `complete=false`; two stepped tags on one rule multiply cells. | both | record:`0030:A7`, `0030:§performance-expectations`; source: `go list -deps ./internal/table` names only `internal/resolve` (confirmed) |
| D-3 | `0030:§consequences` (enum residual), `0030:C3` | `enum` cap exclusion is `in = [<all but last>]`, re-typed on every domain change; the Problem Statement's grievance is reduced, not removed, on the tier example the record leads with. | Adding a tier: the ladder silently keeps stopping at the old terminal and lint reports a blocking `graph-coverage-gap` on the new member (NOT a load refusal, NOT a silent widening — see §3.2); author must find and edit the `in` list in a rule they did not open. | both | record:`0030:§consequences`, `0030:C1` (admitted-cell conjunction); source:`internal/graphlint/taxonomy.go:89` (`CodeCoverageGap` blocking). Both passes' stated symptoms corrected. |
| D-4 | `0030:A2`, `0030:§consequences`, `0030:C3` | Authored `domain` order becomes step order with no load-time or lint-time signal; dispositioned to the authoring guide. | Alphabetizing a `domain` array silently changes successor states; no error, no finding, no diff signal. | both | record:`0030:C3` ("with no lint to catch it"), `0030:A2` |
| D-5 | `0030:A10`, `A11`, `A13`, `A14` (+ `A5`, `A6`) | Six Critical Assumptions `Pending`; C1 states subsumption, both-carriers, bound ordering and step+clear as settled; A10's "If wrong" reopens C1's expansion half. The gate's status-consistency rule will refuse this at Finalize. | Contract locks on unverified mechanism; C1 rewritten post-lock when the MVV disagrees. | both | record:`0030:A5`,`A6`,`A10`–`A14` statuses; `0030:G-assumptions` rule text |
| D-6 | `0030:A12`, `0030:§phase-2-expand-into-cells`, `0030:§approach` | Cross-package relocation of `Evaluator` + `intWidth` with three call-site repoints and ~25 test refs, framed as "one choice-point kind in `expand`"; A12 `Pending`, verified by compile only; RDR 0012 is `Draft` on the same file and unsequenced. | Phase 2 slips; a guard-evaluation regression unrelated to `step`; merge conflict with 0012 whichever lands second. | both | source:`internal/guard/product.go:541`, `internal/graphlint/reach.go:502`, `internal/cli/flow_resolve.go:561`, `internal/guard/declaration.go:190`; 7 `internal/guard/*_test.go` files, 25 `Evaluator` refs; record:`0012:§metadata` Status Draft |
| D-7 | `0030:A5`, `0030:C1` (`unless`), `0030:C2`, `0030:D-selection-predicate` | `unless` on a stepped tag is inert at admission; the intuitive `unless tier = "large"` yields a dead row whose only signal is advisory `graph-redundant-row` (and only when a same-group, same-`Kind()`, projectable sibling claims the cell); C2 patches it in the refusal text. | Model loads clean; a suppressible advisory; the row the author thinks guards the boundary does nothing. | both | source:`internal/graphlint/groups.go:299-313` (gate confirmed), `internal/guard/product.go:568-640` (advisory DOES fire on projectable rows — A's non-firing claim corrected); record:`0030:C2` last paragraph |
| D-8 | `0030:A8`, `0030:C3`, `0030:§decision-rationale`, `0030:F6` | Flow payloads name the bare rule; per-row identity on a flow payload is disclaimed as a possible successor. With D-1, the two surfaces a debugger uses are the two that cannot name the cell. | `flow resolve --as json` says `rule: "retry"`; engineer must cross-reference `dump`/`graph` by hand to learn which of five cells fired. | both | record:`0030:A8` (Refuted, narrowed), `0030:§decision-rationale`; source:`internal/cli/flow_resolve.go::rowByID` first-match (as disclosed) |
| D-9 | `0030:C1` (last paragraphs), `0030:§phase-2-expand-into-cells` | `expand` is declaration-free (`expand(base Row, predicates []Atom, outcome Atom, writes []TagValue) []Row`, no receiver); Phase 2 assigns admitted-cell evaluation, `intWidth`, per-cell `conform` and the `#`/duplicate guards to it. The record neither threads the model through `expand` nor relocates that work to `renderWrites` (the `*loader` method that has the decl), and never names the placeholder literal it puts in `assignments`. | Phase 2 over-runs; the placeholder shape and the home of the admitted-cell computation are decided ad hoc. | A | source:`internal/table/normalize.go:730` (signature), `:580` (`renderWrites` receiver), `:624` (`conform` in write loop). A's "must thread the model" corrected to "must say which". |
| D-10 | `0030:S1` | The MVV finding oracle projects to `(code, key, dimension, class, reason)`, dropping `rule`, `span`, `element`, `fingerprint` — the only fields that could carry cell attribution — so an implementation with D-1 passes green. | MVV passes on an implementation whose lint findings cannot be attributed. | A | record:`0030:S1`, `0030:§consequences` last bullet |
| D-11 | `0030:C2`, `0030:C1` (third-arm paragraph) | Tightening `max` refuses at load naming a rule the author did not touch; the prescribed remedy (exclude with a positive atom) then produces a `graph-coverage-gap`, which C1 calls "visible, non-blocking" but which is in `blockingCodes` — two errors from one edit, the second mislabelled by the record. | Edit `max`; load refuses on `retry`; add the exclusion; lint now BLOCKS on the uncovered cell until a literal row is written. | A | source:`internal/graphlint/taxonomy.go:85-93` (`CodeCoverageGap` blocking); record:`0030:C1` "non-blocking" vs `0030:A5` "blocking" — record contradicts itself |
| D-12 | `0030:§proportionality` | One contract id carries a grammar, an expansion mechanism, a conformance check, an identity-ordering rule and a shared-evaluator relocation; the gate's split test will flag it. | Scope creep chartered; next RDR touching identity, the evaluator, or the admits filter re-litigates all of 0030. | both | record:`0030:G-proportionality` (template text), `0030:§normative-contracts` structure |
| D-13 | `0030:§consequences` (round-trip bullet) | No migration/equivalence aid for converting an unrolled ladder; `0002:§round-trip-inverse-invariants` not extended; the MVV oracle is not a product surface. | Author converts 24 rows to 8, sees a different-looking finding set (D-1), has no supported way to confirm the conversion was faithful. | A | record:`0030:§consequences`, `0030:S1`/`S2` (test-only oracles) |
| D-14 | `0030:C2` (detail string) | Four-part refusal plus the unconsulted-`unless` clause ride `Failure.Detail` as free text asserted by substring position; no structured fields because C3 forbids envelope members. Record's field census omits `Remedy` and `Line`. | Message frozen by tests or tests broken by rewording; machine consumers get no structure on the record's most important new refusal. | A | source:`internal/table/category.go:201-235` (`Failure`: `Category`, `Detail`, `Offending`, `Remedy`, `Rule`, `Line`); record:`0030:C2` |
| D-15 | `0030:S5b` | S5b asserts its own projectability precondition, so it proves the advisory only on a fixture that projects; the precondition is real but the scenario does not say what a NON-projecting authored ladder gets instead. | Fixture passes; a ladder with an optional-key atom on the stepped rule gets no advisory and no explanation. | A | record:`0030:S5b`; source:`internal/guard/product.go:591-614` (what actually defeats projectability — not `unless`). Downgraded from A's framing. |
| D-16 | `0030:§problem-statement`, `0030:§decision-rationale` row 6 | `dump`/`graph` still show the N expanded rows; the review win is source-only. | Reviewer using `dump` sees the unrolled explosion the record promised to remove. | B | record:`0030:§decision-rationale` ("the cost Approach 1 pays is row 6's second half") — disclosed |
| D-17 | `0030:A4`, `0030:C1` (third arm) | "Unreachable" third arm rests on the frozen operator/kind matrix the record itself names as the natural successor to widen; the fence is safe (exclude → coverage gap) but the argument lapses silently. | Future grammar extension: cells excluded unexpectedly, surfacing as coverage gaps nobody traces to this clause. | B | record:`0030:A4`, `0030:C1`, `0030:§consequences` — claim-only on the future; the exclude direction is source-safe |

Dropped (see §2 for the showing): A C-1 (placeholder/`conform`), A C-3
(mixed comparator), A C-12 (unfilled gate, folded into D-5), A C-17
(`0002:C13` narrowing, answered by `0030:§metadata` Overrides).

## 5. Prose divergence

Pass A's premortem is a story about the record being wrong about the code:
the loader cannot host the expansion as written, `Identity()` is unstable,
lint cannot name a cell. Two of those three mechanism claims (C-1, C-3) do
not survive the record and `normalize.go`; the third (C-4, lint attribution)
is fully source-grounded and is the single most important finding in either
pass, unique to A. Pass B's premortem is a story about the record being
right about the code and shipping its disclosed gaps anyway — enum cap,
domain order, lint cliff, bare `rule` on flow payloads — with every incident
traced to a §consequences bullet. B's closing note ("noticed all nine and
chose to ship") is the more accurate shape of the risk; A's row on lint
attribution is the one defect B's shape cannot see, because the record does
not disclose it. Both passes mis-predicted the add-a-tier symptom in
opposite directions (§3.2), and A's `unless` grounding inverts what
`AcceptedAssignments` actually does (§3.1) — the fix half should treat
source-cited rows from A with the corrections above, not at face value.
