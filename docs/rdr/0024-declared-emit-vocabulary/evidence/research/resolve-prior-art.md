# RDR 0024 — Stage 4 (Resolve) prior-art search record — A6

Model: claude-opus-5[1m]

Assumption A6: "DMN decision tables carry an allowed-values list on
output clauses that conformant tooling checks output entries against —
the external alignment claim for declare-then-prove over answers."

Method: Prior Art. Required evidence shape: a named external system plus
a section or page reference.

## Result

`negative: A6 — no corpus evidence in DevRef, StateMachineLit, StateMachineRes, PapersFast`

No section-anchored citation for either leg was found within budget. A6
is recorded NEGATIVE (unverified), not refuted: the corpora simply carry
no OMG DMN / BPM-suite material. Absence of coverage is not evidence the
DMN feature does not exist.

## Leg split

- Leg 1 (DECLARATION exists in the standard — an allowed output-values
  list on the decision table's output clause): NOT ESTABLISHED from a
  cited section. Partial, non-qualifying material only — see the
  rejected web branch below, which documents the output clause and its
  `typeRef` type constraint but names no allowed-values list and carries
  no spec section or page.
- Leg 2 (conformant TOOLING CHECKS authored output entries against the
  declared list): NOT ESTABLISHED. No tool documentation describing
  output-entry validation against declared output values was reached.
  This is the load-bearing half for 0024's "declare-then-prove over
  answers" alignment, and it is the leg with zero support.

## Corpora searched

`DevRef` (pdf, 148 software-engineering reference books),
`StateMachineLit` (pdf, state-machine literature), `StateMachineRes`
(markdown, curated state-machine research), `PapersFast` (pdf, ~3623
academic CS papers).

## Queries run

Full-text (`arc search text --corpus <C> --limit 6|8 --json`):

1. `DevRef` — "DMN decision table output values allowed values output
   clause" → 2 hits, both linear-algebra / linear-programming pages from
   *Introduction to Algorithms* 4th ed. Rejected.
2. `PapersFast` — "Decision Model and Notation decision table output
   clause allowed values" → unranked bag-of-words match, top hit an
   Adaptive Resonance Theory survey. Rejected.
3. `"decision table" DMN` — run against all four corpora.
   - `DevRef`: *Object-Oriented Software Engineering* (Jacobson) p351
     table-of-test-cases sense; *developer-testing.pdf* p142–143
     decision tables as a BDD example format; *Code Complete* p472
     table-driven methods. All the generic "decision table" sense, no
     DMN, no output clause. Rejected.
   - `StateMachineLit`: zero hits.
   - `StateMachineRes`: zero hits.
   - `PapersFast`: multi-dimensional decision tables vs decision trees
     (Article file_50), plus incidental "decision. Table II" caption
     splits. Rejected.

Semantic (`arc search semantic --corpus <C> --limit 5`):

4. `DevRef`, `PapersFast` — "business rules decision table output column
   restricted to a declared list of allowed output values validated by
   the rule engine" and the reduced form "decision table output values"
   → best score 55% (*developer-testing.pdf* p142, insurance-premium
   decision-table example), second 49% (PostgreSQL 15 manual p1580,
   `ALTER TABLE ... user_catalog_table`). Both unrelated. Rejected.

No corpus term-frequency signal for `DMN`, `Decision Model and
Notation`, or `outputValues` in any of the four corpora.

## Rejected branches

- Sibling analysis repo `../state-machines` (`AUDIT-*.md`, `EVAL-*.md`,
  `contrast/`, `repos/`): ripgrep for `DMN`, `Decision Model and
  Notation`, `decision table`, `outputValues`, `output clause` returned
  only substring false positives inside unrelated words in
  `ANALYSIS-determinism-topology.md`. A word-boundary re-run for `DMN`,
  `Decision Model and Notation`, and `outputValues` returned zero
  matches. No DMN prior art in the sibling repo.
- Web search (one attempt, per method budget): "DMN decision table
  output clause outputValues allowed output values validation Camunda"
  → Camunda 8 and Camunda 7 DMN decision-table output docs. These
  establish that an output clause carries id, label, name and a
  `typeRef` type, and that an evaluated output entry is checked to
  convert to that declared type — a TYPE constraint, not an
  allowed-VALUES list, and the pages carry no spec section or page
  anchor. Does not satisfy either leg as written. Not accepted as a
  citation; recorded here so the branch is not re-walked. Publisher and
  OMG spec pages block automated fetch and were not attempted.

## Effect on the RDR

A6 was declared non-load-bearing when it was written, and this negative
leaves that unchanged. The choice rests on the in-repo mirror
`0002:C22` (tag type-model declaration) plus the RDR 0003 kind
vocabulary in `internal/table/model.go::declaredKinds`, both openable
and quoted in the RDR body. The DMN sentence in the RDR should be read
as an unverified aside, not as external support; if it is to stay, it
should be marked as such rather than cited.
