# Critique dual-model diff — cli/0021

Pass A `critique.md` — Model: claude-opus-5
Pass B `critique-modelB.md` — Model: claude-sonnet-5

Two distinct stamps: a genuine cross-model draw, not a single-model fallback.
Ledger rows are keyed per file (`C-N`); the join key across files is the
element id (`rdr anchors --record 0021`).

## Anchor sets

| Set | Anchors |
| --- | --- |
| Both (converged) | 0021:A8, 0021:C2, 0021:C4, 0021:C5 |
| Opus only | 0021:A1, 0021:A6, 0021:RT1, 0021:S6, 0021:§metadata, 0021:§pre-lock-mini-checks |
| Sonnet only | 0021:A2, 0021:C1, 0021:§context, 0021:§decision-rationale, 0021:§failure-modes, 0021:§problem-statement, 0021:§references |

Net-new on this pass: none (`comm -13` empty — one pass, so ledger = pass).

## Where the models AGREED

All four converged anchors. Agreement here is signal, not redundancy — two
isolated hostile contexts on different models independently landed on A8, C2,
C4 and C5. Three of the four produced accepted fixes:

- **0021:A8** — both flagged a Pending assumption carrying normative weight.
  Opus framed it as "Pending under normative MUSTs"; Sonnet as "load-bearing for
  C4's MUST-clause, ungated by Prerequisites". Same defect, two framings. FIXED.
- **0021:C5** — both hit the envelope/mode surface. Sonnet's angle was the
  bare-JSON-on-stdout collision (already covered by F2); Opus's was sharper and
  correct: the DOT key is never named. FIXED via Opus's framing.
- **0021:C2** — both raised it; the models disagreed on WHAT. Opus: the
  `<opaque>` sentinel reaches the wire unnamed (founded in code). Sonnet:
  no machine-checkable schema, and reviewer-traceability to authored TOML
  (both weak/out of frame). FIXED on Opus's row only.
- **0021:C4** — both raised it, BOTH WRONG in different directions. Opus
  claimed a code collision with `graph-product-too-large`; Sonnet claimed the
  neutrality oracle is untested against a future observer arm. Convergence on an
  anchor is not convergence on a defect: C4 is the densest clause in the record,
  which draws attention without implying a fault. Both dismissed with cites.

## Where the models DISAGREED — the decorrelated-error signal

The disjoint sets are the payoff of the cross-model draw:

- **Opus went deep on mechanism** — spike/oracle/invariant internals (A6
  spellings, S6 escaping order, RT1's quantifier, A1's TextLiner payload size).
  It found the two real code-grounded defects (`<opaque>`, the DOT key) and
  the one real logic defect (RT1's tautology). It also produced the two most
  confidently wrong rows (C-1/C-9), traceable to a stale doc comment in
  `taxonomy.go:150-152` that misattributes the node ceiling — charted.
- **Sonnet went wide on framing** — problem statement, context, decision
  rationale, failure modes, references. Every row was either already covered by
  the record's own text or a re-raise of a settled ruling. Zero accepted fixes.

Neither model's misses were the other's: Opus never questioned the framing;
Sonnet never opened the code. The draw behaved as intended — differently wrong
rather than uniformly wrong — and the accepted fixes all came from the pass
that read source.

## Outcome

4 fixed (A8, C2, C5, RT1) · 13 dismissed-with-cite · 1 charted-to-successor.
1 new Pending assumption (A9). No tiebreakers. Loop: converged (`loop-small-fix`).
