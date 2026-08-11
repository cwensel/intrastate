Model: claude-opus-5[1m]

# RDR 0007 — Stage 4 citation cache

Extends `propose-prior-art.md` (Stage 2). That pass found ⚠ no corpus
coverage for the undecidable-guard problem class and demoted the
SQL/Kleene K3 claim to a Stage 4 assumption (A4). This pass resolves it.

## Accepted citation

### C3 — PostgreSQL 16 §9.1: SQL three-valued logic, full strong-Kleene truth tables (A4)

Source: PostgreSQL 16 Documentation, Chapter 9 "Functions and Operators",
§9.1 "Logical Operators", p. 231 (PDF page 272). Corpus `DevRef`;
local path `~/Documents/References/PostgresOfficial/postgresql-16-US.pdf`.

> SQL uses a three-valued logic system with true, false, and null, which
> represents "unknown". Observe the following truth tables:

| a | b | a AND b | a OR b |
|---|---|---|---|
| TRUE | TRUE | TRUE | TRUE |
| TRUE | FALSE | FALSE | TRUE |
| TRUE | NULL | NULL | TRUE |
| FALSE | FALSE | FALSE | FALSE |
| FALSE | NULL | FALSE | NULL |
| NULL | NULL | NULL | NULL |

| a | NOT a |
|---|---|
| TRUE | FALSE |
| FALSE | TRUE |
| NULL | NULL |

> The operators AND and OR are commutative, that is, you can switch the left
> and right operands without affecting the result.

The two rows load-bearing for A4 are explicit: `FALSE AND NULL = FALSE` and
`TRUE OR NULL = TRUE` — definite results from present data dominate; unknown
propagates otherwise. The commutativity note rules out a left-to-right
short-circuit reading: the dominance is semantic, not evaluation-order
dependent, which matters for an evaluator free to visit atoms in any order.

This is strong Kleene (K3), where conjunction is `min` under `F < U < T`, and
is the discriminating case against *weak* Kleene (Bochvar), in which `U` is
infectious and `F ∧ U = U`.

## Query ledger (budget: 2 of 4 queries, 1 of 6 hits opened)

1. `DevRef` — "SQL NULL three-valued logic unknown truth value comparison" —
   **ACCEPT**. Rank-5 hit (0.730) is PostgreSQL 16 §9.1; snippet truncated
   mid-table, so the hit was opened once to confirm the full tables.
2. `PapersFast` — "null values relational databases three-valued logic
   semantics" — **REJECT**. All five hits were bibliography/reference-list
   pages (Libkin 2020, Toussaint VLDB 2022, Gyssens JACM 1989) plus one page
   arguing *against* 3VL; titles corroborate 3VL as the standard SQL
   treatment but no hit carried truth-table semantics in body text.
   Superseded by C3.

Stopped on verification per the sufficiency bar. `StateMachineLit` and
`StateMachineRes` were not re-queried — the Stage 2 pass already searched
both for this problem class and found nothing.

## Not cited from corpora (recorded as negative)

`negative: A1/A5/A7 — no corpus evidence sought in DevRef/PapersFast/
StateMachineLit/StateMachineRes.` These are peer-RDR source-search claims
verified against `docs/rdr/0003-guard-predicate-exhaustiveness.md` and
`docs/rdr/0002-transition-table-as-reviewable-data.md` directly, which is the
domain that owns the behavior; a corpus pass would be the wrong source.

`negative: A6 (second half) — no corpus or peer-RDR evidence found for the
"read failed vs genuinely absent" distinction.` Searched
`docs/rdr/0004-accessor-execution-safety-model.md` (the owning contract) and
swept `docs/rdr/` for read-failure-vs-absence language: the only hits are
inside A6's own text in RDR 0007. RDR 0004 is silent.

## Scope note on C3

A4 is verified for conjunction and disjunction — the connectives RDR 0007
uses. K3 negation (`¬U = U`) is also in the cited table. K3 has no designated
-value tautologies (`p ∨ ¬p` is `U` when `p` is `U`); RDR 0007 does not rely
on excluded middle over guard atoms, so this does not bite, but a future
exhaustiveness proof over guard atoms would need that stated separately.
