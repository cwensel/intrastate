Model: claude-opus-5[1m]

# Repeatability resolve — dispositions (RDR 0007)

Origin ledger = `diff.md`'s D-1..D-10 and G-1..G-8. Each entry exits exactly
once. Per the REPEATABILITY DIFF clause, a divergence is dispositioned as
**pin** / **cut** / **single-source** / **leave non-normative** / **tiebreaker**
before it becomes an edit.

## Disagreements

| ID | Disposition | Action | Section touched |
| --- | --- | --- | --- |
| D-1 | **pin** → fixed | Stated survivor membership constructively: GuardTrue or GuardUnevaluable ⇒ survivor; only GuardFalse prunes. Names the consequence runs 2/3 derived silently — an unevaluable row's `RequiresOwned` keys DO raise `owned_state_unavailable`. Grounded on `resolve.go::gate` (appends unless GuardFalse). | Normative Contracts (aggregation block) |
| D-2 | **pin** → fixed | `Refusal.Rows` element type pinned to `[]resolve.RowRef` (shipped `(RuleID, SourceLocator)` pair). run-2 was right; runs 1/3 invented `[]Row` / `[]RowID`. Required because the adjacent clause makes `Rows` a normative assertion target. | Normative Contracts (refusal-payload block) |
| D-3 | **pin** → fixed | Mapping failure is NOT GuardUnevaluable. `fable-5`'s invented fallback would make an evaluator defect indistinguishable from missing artifact state — the RDR's core conflation, one layer up. Pinned as unreachable post-lint (A1/A7) and a Go-error-path programmer error. | Normative Contracts (atom-vocabulary scope) |
| D-4 | **pin** → fixed | `Evaluate` returns a bare `GuardResult`; no second channel. All three runs agreed on the value and all three marked it GUESS — unanimous-but-unstated is exactly what the lens is for. Grounded on `resolve.go::GuardEvaluator`. | Normative Contracts (atom-vocabulary scope) |
| D-5 | **pin** → fixed | Seam is the `GuardEvaluator` **interface**, not a func type. Low severity on its own; pinned because the Phase 2 harness signature — which the RDR *does* fix normatively — takes it as a parameter. | Normative Contracts (atom-vocabulary scope) |
| D-6 | **pin** → fixed | `TagSet` is a concrete struct with an unexported map and no exported constructor, so **no out-of-package caller can build a view**. Verified: no `NewTagSet`, `assemble` unexported, every shipped test drives views through `resolve.Input`. Phase 2 must choose Input-driven vectors or add a constructor. | Phase 2 (new *View construction* bullet) |
| D-7 | **single-source** → fixed | Zero/multiple counting is RDR 0001's and happens in `Resolve` after `gate` returns (`resolve.go` L347–356). Cited rather than restated, and the RDR now says explicitly it fixes three ordering facts and no others. | Normative Contracts (aggregation block) |
| D-8 | **pin** → fixed | Named all five refusal kinds where the taxonomy is first called closed: `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, `unmodeled_outcome`. `fable-5`'s `invalid_input` confirmed a confabulation against `resolve.go::RefusalKinds`. | Technical Environment (RDR 0001 summary) |
| D-9 | **pin** → fixed | Phase 2's reference evaluator is unexported, and its guard encoding is vector-local and non-normative — it must not pre-empt A10's mapping. Only `fable-5` surfaced this bootstrapping tension. | Phase 2 (item 1) |
| D-10 | **pin** → fixed (+ new A14) | `exists` carries a boolean literal; verdict is `presence == literal`. Grounded on RDR 0003's operator matrix (literal shape "boolean") and its fixture (`exists = true`); `glm-5.2`'s unary reading refuted. Consequence: a single-row absence test exists, so A5's two-row pattern is needed only for the *disjunction* — A5 rewritten accordingly. | Normative Contracts (existence clause); A5; Scenario 6 |

## GUESS clusters

| ID | Disposition | Reason |
| --- | --- | --- |
| G-1 | **dismissed-with-cite** (already tracked) | The `Row.Guard` ↔ atom mapping is A10, already Pending and already a lock-blocking Prerequisite. The lens adds force, not a new defect: all 3/3 runs ranked it first, confirming every atom-level clause is unrenderable until it closes. No edit — acting again would duplicate an open assumption. |
| G-2 | **leave non-normative** → fixed (scoping note) | Tag value typing is RDR 0003's declared-kind system, not this RDR's to model — specifying it here would be the over-spec trap and take a second contract into a single-contract RDR. Added an explicit statement that the clauses bind regardless of representation. |
| G-3 | folded into D-2 | Same contract; resolved there. The remaining `Row` field-type guesses are RDR 0001's shipped struct, not this RDR's silence. |
| G-4 | folded into D-4 | Same contract; resolved there. |
| G-5 | folded into D-8 | Same contract; resolved there. |
| G-6 | **leave non-normative** | The harness package name and vector encoding are *deliberately* deferred ("Phase 2's to fix") and both runs read the deferral correctly. Legitimate under-specification by design; the three adjacent constraints (importability, exported shape, caller-supplies-evaluator) stay normative. No edit. |
| G-7 | **pin** → fixed | Named `ProvenanceRecognized` alongside the two constants already quoted. One-word gap in a normative clause quantifying over all three. |
| G-8 | **pin** → fixed | Lifted the K3 conjunction and negation tables from A2's Evidence into the normative combination block. The real finding is placement: the tables are load-bearing but lived outside the ```normative``` block, so 2/3 runs had to guess them. |

## Needs (re)verification — carried to Stage 6

- **A14 (new, Pending)** — `exists` carries a boolean literal and
  `all … exists = false` is authorable. Method: Source Search against RDR 0003's
  grammar, specifically whether its placement rule ("positive atoms in `all`,
  negative in `unless`") forbids a negative-polarity literal in `all`. The
  fixture only exercises `exists = true`. Does not block the domain rule; it
  fixes whether a single-atom absence test exists, which sets how hard A12 binds.
- **A5 (Pending, unchanged status; wording sharpened)** — the claim is now
  narrower and more accurate: the two-row pattern is required for the
  *disjunction* only. Its load-time half still rides on A12. If A14 fails, A5
  reverts to its former stronger dependence on A12.
- **A9, A10, A12, A13 (pre-existing, Pending)** — unchanged by this lens. G-1
  reinforces A10's blocking status without altering it.
- **A6b (pre-existing, open by decision)** — unchanged.

No previously-Verified assumption was invalidated. A3's "no kernel change" claim
is untouched: every D-* fix states shipped behavior or binds the future
evaluator; none asks the kernel to change.

## Tiebreakers

None escalated. Two forks looked load-bearing and collapsed on evidence:

- **D-10 (`exists` arity)** — resolved against RDR 0003's own operator/kind
  matrix (boolean literal shape) plus its fixture, not by preference. The
  residual uncertainty is narrower than the fork and is recorded as A14 rather
  than escalated.
- **D-3 (mapping-failure verdict)** — the apparent either/or (route it to
  GuardUnevaluable vs. leave it undefined) dissolved once the totality
  requirement was read alongside the load-time rejection rules: a well-loaded
  guard that cannot be parsed is unreachable, and routing it to the verdict
  would re-create the RDR's own target conflation.

## Lens health (review gate)

**Healthy.** Findings localize to specific interfaces with named passages; GUESS
markers cluster on the same contracts across models (G-1 3/3, G-2 3/3, G-3 3/3);
the Agreement paragraph is substantial — all three models reproduced the entire
domain rule, strong-Kleene shape, provenance-blind presence, `gate` ordering, the
aggregation veto, and the `RequiresOwned` narrowing identically. That is the
determinacy signal: the RDR's *content* is reproducible, and the divergences are
its unstated Go surface. Not the unhealthy "identical-but-confidently-wrong"
shape — every run carried explicit GUESS markers, and the two confabulations
(`invalid_input`, the GuardUnevaluable fallback) came from one model and were
refuted against `main`.
