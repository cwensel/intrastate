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

---

# Stage 4 re-entry (scoped) — 2026-08-21

Scope: `Draft [revised from Final 2026-08-12 — re-verify A3, A16, A17,
A18, A19, A21, A22]`. A1/A2/A4/A5/A6a/A6b/A7/A8/A9/A11/A13/A14/A20 carry
forward as already Verified; A12 stays Deferred (§JD-4, out of scope).

The re-proposed approach moved the enforcement site to the kernel
(JDR 0001 §D4/§JD-12), so the in-scope assumptions are peer-RDR and
own-source claims, not corpus claims.

`negative: A16/A17/A18/A19/A21/A22 — no corpus evidence sought in
DevRef/PapersFast/StateMachineLit/StateMachineRes.` These are peer-RDR
(0002/0003/0005/0009), JDR, and own-kernel source-search claims; the
owning documents are the domain that owns the behavior, so a corpus pass
would be the wrong source. C3 (PostgreSQL 16 §9.1) continues to carry A4
unchanged and is not re-searched.

## Accepted anchors this pass

- **A16** — RDR 0003 operator/kind matrix: the `exists` row is the sole
  absence-inspecting operator ("Tests presence or absence, not value
  equality"), literal shape `boolean`; there is no `absent` counterpart
  token. JDR 0001 §D4(b) states the kernel cost as "one grammar fact (the
  existence operator token and its two boolean literal forms), exported as
  constants 0002's normalizer MUST emit." RDR 0003 and RDR 0002 each carry
  the identical joint-check line naming the kernel-exported token.
- **A17** — RDR 0003 Load-Bearing Decisions → Operator semantics: every
  operator reads one value against one literal. Decisive passage for the
  "bounded" worry, RDR 0003 Technical Design: "A guard may still compare an
  unbounded integer at runtime, but lint must report that it cannot prove
  exhaustive coverage for that dimension" ⇒ declared bounds are a LINT
  input, never an evaluation input. `contains`'s declared element universe
  is likewise a lint-proof input. Identity tuple is `(tag, operator,
  literal)` with no kind field.
- **A18** — `internal/resolve/resolve.go::evaluateGuard` nil-seam branch
  returns `GuardUnevaluable` ("a guarded row with no seam is undecidable,
  never evaluated by the kernel itself"); `KindGuardUnevaluable` is
  documented as "the guard seam reported a predicate it could not decide"
  — the shipped contract names the seam as the reporter. RDR 0003's only
  parse-rejection clause is about the authored LITERAL ("A predicate whose
  literal cannot be parsed as the declared tag kind MUST be rejected before
  resolution"), not a runtime value, so the runtime case is genuinely open
  in 0003 and A18's non-presumptuous stance is correct.
- **A21** — RDR 0009 (Final) Normative Contracts: "a Row with a non-empty
  Escape list MUST have an empty Writes slice." `RequiresOwned` appears
  ZERO times in RDR 0009. `resolve.go::missingOwned` over an empty slice
  never enters the loop, so an escape row raises no
  `owned_state_unavailable`. Composition valid.

## Rejected / downgraded this pass

- **A19 (second half)** — REFUTED as written.
  `internal/cli/clierr/clierr.go::CLIError` is flat strings only (`Code`,
  `Message`, `Param`, `Detail`, `Hint`); `Detail` is `string`, documented
  "May be multi-line" — multi-line TEXT, not a nested object. RDR 0007's
  payload contract requires a two-level array-of-objects (per row
  `(RuleID, SourceLocator)`; per atom key, block, reason). RDR 0005's
  normative split puts structure on the success side only ("Failures MUST
  use the existing CLIError JSON/text envelope"). So the payload cannot
  ride JSON structurally without a new envelope field, which §JD-8
  pre-authorization explicitly does not cover.
  Secondary: `guard_unevaluable` HAS a code row (RDR 0005 Failure Modes:
  `flow-guard-unevaluable` / `GroupUserEnv`) — but scoped "from supplied
  facts", which the cluster gate already flagged against the
  provenance-blind view. `owned_state_unavailable` has NO code row; that
  gap is pre-authorized by §JD-8 ("need `Code` values").
- **A22 (second half)** — DOWNGRADED to Pending. Kernel half verified:
  `resolve.go::assemble` keys `view.tags` by `Tag.Key` verbatim and
  `TagSet.Lookup` is a bare map index — no canonicalization anywhere. But
  RDR 0002 has NO canonicalization or tag-name grammar clause; its
  declaration-completeness rule ("The model MUST declare every tag it
  matches or writes, including each tag's provenance") is a different rule
  and fixes spelling only under the unstated premise that declaration
  lookup is itself exact. The supporting sentence lives in JDR 0001 §D4's
  landing note ("the normalizer ... canonicalizes key spellings before a
  row exists") — an obligation FILED on 0002, and 0002's Refinement
  Context Direction list does not carry it.

## Cross-cutting observation

A19 and A22 fail the same way: each cites a JDR §D4/§JD-8 sentence as
though it were a peer-RDR contract, where the receiving RDR has not
absorbed the obligation. JDR 0001's own Problem Statement names this
failure mode — "An obligation filed on a document that never received it
binds nobody." Before lock, every §D4 landing-note duty should appear in
the receiving RDR's Refinement Context Direction list; 0002's
canonicalization duty currently does not.
