Model: claude-opus-5[1m]

# 3amigo resolve — disposition ledger (iteration 1)

Origin ledger = the 15 persona findings (PM-1..5, IMP-1..5, QA-1..5). Every
entry exits exactly one way. Grounding gate run against code on `main`
(commit 9ad06a3), `{RDR_RESOURCES}` design docs, and the RDR's own decided text.

| # | Disposition | Origin | Section touched |
|---|---|---|---|
| IMP-1 | **fixed** | H-1 hinge | Normative Contracts — new *Scope of the atom vocabulary* clause; A10 added (Pending) |
| QA-1 | **fixed** | H-1 | Testing Strategy — blocking prerequisite for rows 4–8, pointing at A10 |
| IMP-5 (empty-block landing site) | **fixed** | H-1 | Scenario 9 *Test level*; `RequiresOwned` clause corrected on "clears" |
| QA-2 | **fixed** | H-2 | MVV + Scenario 1 — `Escaped:false` was factually unassertable (Plan field, nil on refusal) |
| QA-4 | **fixed** | H-2 | Scenario 5 — discriminator moved to `Refusal.Rows` |
| IMP-3 | **fixed** | H-2 | Normative Contracts — owned-before-unevaluable pinned; `Refusal.Guard` demoted to diagnostic |
| PM-1 | **fixed** | H-2 | Failure Modes — asymmetry vs. the Problem Statement's promise stated honestly |
| PM-4 | **fixed** | H-3 | Normative Contracts aggregation — veto made explicit; D5 tension answered |
| PM-2 | **fixed** | H-4 | Consequences — exposure bounded (no evaluator exists yet); sequencing named as the mitigation |
| PM-3 | **fixed** | H-4 | Decision Rationale — operator-recovery criterion scored explicitly as B's worst axis, accepted with reason |
| PM-5 | **fixed** | — | Consequences — present-tense masking claim corrected to contract-vs-shipped |
| IMP-2 | **fixed** | — | Normative Contracts — new *Presence is provenance-blind* clause; A11 added (Verified) |
| IMP-4 | **fixed** | — | Phase 2 — importability + signature constraints on the harness |
| QA-3 | **fixed** | — | Scenario 2 — `TestAdv1b` miscitation corrected; A/B control named as new Phase 2 work |
| QA-5 | **fixed** | — | Empty-block normative clause marked conditional on A9 (status consistency) |

No finding was dismissed and none was charted to a successor: all 15 grounded
against `main` and all fell inside this RDR's existing seam (guard-evaluation
domain). The one net-new-scope candidate — the guard-text grammar itself — was
**not** absorbed: it is recorded as A10 and handed to RDR 0003 in Phase 3, which
is the anti-wormhole disposition, not an expansion.

## Grounding gate notes

Verified firsthand against `internal/resolve/resolve.go` on `main`:

- `Escaped` is a `Plan` field (L250, set in `planOf` L509); `Result.Plan` is nil
  on every refusal — QA-2 is a factual defect, not a preference.
- `Refusal.Guard` is a single string set from `slices.MinFunc` over
  `compareRefs` (L409–414); `Refusal.Rows` carries all undecidable rows — IMP-3
  confirmed.
- `TagSet.Lookup` is provenance-blind and the only exported accessor; `has` is
  unexported with `missingOwned` as its sole caller — IMP-2/A11 confirmed.
- `fixtureGuards` is in package `resolve_test` (unimportable); `GuardEvaluator`
  is under `internal/` — IMP-4 confirmed.
- `TestAdv1b` decides `"never"` from guard text with `status:Draft` present and
  never varies presence — QA-3 confirmed.
- `Row` has `NextTags`/`Writes` and no clear field; RDR 0002 renders a clear as
  a write — IMP-5's "clears" objection confirmed and corrected.
- RDR 0003's spike fixture declares `[tags.cluster_eligible] provenance =
  "observed"` and guards it `exists = true` — the witness for IMP-2.

Re-raise check: PM-4's D5 objection was tested against the RDR's decided text
and D5's own wording (`0001-…/artifacts/deviations.md`). D5 scopes evaluation to
*matching* candidates to stop an **unrelated** row poisoning a resolution; an
unevaluable row that matched is not unrelated. Rather than dismiss it as a
re-raise, the implicit reasoning was made explicit in the draft — the corpus's
top flapping cause is leaving such a decision implicit.

## Needs (re)verification — carried to Stage 6

- **A10 (new, Pending)** — guard-structure mapping at the `Row.Guard` boundary.
  Method: Source Search against RDR 0002's normalization contract and RDR 0003's
  grammar. Blocks lock (added to Prerequisites) and blocks encoding Phase 2
  vectors 4/6/7/8. If unspecified anywhere, resolution is to hand it to RDR 0003
  as a named obligation, not to specify a grammar here.
- **A11 (new, Verified)** — provenance-blind presence. Verified inline against
  the four `TagSet` readers; no further work owed.
- **A9 (pre-existing, Pending)** — unchanged by this lens; its normative clause
  is now explicitly conditional on it.
- **A6b (pre-existing, open by decision)** — unchanged; PM-3 stacked it into the
  operator-recovery cost, which is now scored in Decision Rationale.

No previously-Verified assumption was invalidated by these edits.

## Tiebreakers

None escalated. The one load-bearing fork — whether to specify a guard grammar
here or hand it to RDR 0003 — collapsed on evidence: RDR 0003 owns the grammar
by its own Normative Contracts, its fixture already authors guards structurally,
and the kernel seam deliberately attaches no meaning to the text. Specifying a
grammar in 0007 would take a second contract into a single-contract RDR and
reopen a Final peer.
