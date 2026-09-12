Model: claude-opus-5

# Grounding — iteration 2 (delta-scoped)

Scope: the two tier assignments iteration 1 ADDED to C4. The loop row returned
`rerun` (fix=substantial), and the delta is the claim iteration 1 asserted by
reasoning rather than read from source — that `findings[].operator` and
`findings[].block` are both `frozen`.

## Findings

**GF5 — REFUTED. `findings[].block` is not `frozen`; C4 contradicted a Final
joint decision.**

Iteration 1 assigned `block` → `frozen`, citing `0007:C1` as fixing its members
and `0007:C6` as naming `all`/`unless`. Neither says that:

- `0007:C1` fixes the TYPE SHAPE (an exported named string type), and explicitly
  cedes vocabulary ownership: "Nothing else about the atom … is this RDR's".
- `0007:C6` fixes the row-verdict FORMULA's operands, not the set's cardinality.

`JDR 0001 §D12` ("`Block` cardinality") already adjudicated exactly this and went
the other way, verbatim: "**Resolved: one exported type, three constants.**
`resolve.Block` gains `BlockMatch` in 0007 Phase 1; 0007's 'exactly two' is the
fence that gives." The set has GROWN under record, from two members to three.
`internal/resolve/guard.go::isGuardBlock` is written for exactly that, treating
"BlockMatch …, the zero value, and any future token" as non-operands.

A `frozen` assignment here would have put this RDR in contradiction with a Final
JDR at Stage 7.1. Corrected to `append-only` with §D12 cited.

**GF6 — CONFIRMED-with-corrected-reason. `findings[].operator` is `frozen`, but
iteration 1's justification was factually wrong.**

The verdict survives; the reasoning does not. Iteration 1 wrote that a new
operator "is admitted by the operator/kind matrix (`::Accepts`) rather than
appended to a list". Backwards: `operators` IS a list and IS what gets appended
to; `matrix` is a SECOND edit site and `Evaluate`'s switch a third.

The real authority is a contract iteration 1 never cited: `0003:C7` — "The
initial operator vocabulary MUST be closed and typed: equality, membership,
bounded integer comparison, existence, and set containment. Unknown operators
MUST be rejected during parse or lint before resolution." Enforced exhaustively
by `internal/guard/guard_grammar_0003_test.go::TestReq2_OperatorVocabularyIsClosedAndTyped`
(`slices.Equal` over the 8-token literal, plus negative controls rejecting
`neq`/`matches`/`regex`), with a second test requiring every operator be
classified in `SingleValueOperator`'s partition.

Rewritten to cite `0003:C7` and the test.

## Why this iteration was owed

Both findings are defects the first pass INTRODUCED while fixing GF1 — a fix that
assigned tiers from local code reading without checking whether a settled peer
record already governed the vocabulary. The loop row's `rerun` on
`fix=substantial` is what surfaced them; one of the two contradicted a Final JDR.

## Still open

None. Both delta claims now rest on cited contract text (`0003:C7`,
`JDR 0001 §D12`) read first-hand in the parent context, not on inference.
