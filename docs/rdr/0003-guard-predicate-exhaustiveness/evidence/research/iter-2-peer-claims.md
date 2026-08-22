Model: claude-opus-5[1m]

# Stage 4 Resolve (scoped re-entry, iter-2) — accepted citations

Date: 2026-08-21
Scope: A5 (named) + anchors the demotion edit touched (A4; the two new
normative clauses citing JDR 0001 §D4 / §JD-4).

No corpus search was needed: every in-scope claim is a peer-document or
own-source claim. Corpora (`StateMachineRes`, `StateMachineLit`) were not
re-queried — the demotion changed no external-behavior claim.

negative: none — every in-scope claim landed on a local authority.

## Accepted citations

### A5 — parsed-atom shape and source identity

- JDR 0001 §D1, resolved option (d): "Row carries a slice of atoms (key,
  operator token, literal, block) instead of a string. No reconstruction step,
  so no mapping failure."
- RDR 0007 Normative Contracts (SEAM): "A candidate row carries its guard as a
  slice of parsed atoms — key, operator token, literal, block ∈ {all, unless} —
  the shape JDR 0001 §D1 fixes."
- RDR 0002 Normative Contracts: "Each candidate row MUST retain its source rule
  id and source locator."
- RDR 0002 Validation/Testing Strategy scenario 2: candidate rows retain source
  rule ids/locators through `all`/`unless` expansion.
- RDR 0006 Technical Design: findings carry "rule/context id when available,
  source span when available".
- `internal/resolve/resolve.go::Row` — `RuleID` / `SourceLocator` ship TODAY,
  commented as "the source identity RDR 0002 requires every normalized row to
  retain". A5's actual subject is therefore already true in built code.

### A4 — error ownership split

- JDR 0001 §D4: "Lands in 0007, the single normative home of the seam, the
  domain rule, and the payload."
- RDR 0007 Normative Contracts: `Refusal` carries `Undecided []UndecidedRow`;
  the refusal names the missing thing per row and per atom.
- RDR 0006: `graph-overlap`, `graph-coverage-gap` finding codes.
- RDR 0005 Normative Contracts: failures MUST use the existing `clierr`
  envelope.

### Evaluator-scoping clause (kernel owns presence/existence/combination)

RDR 0007 Normative Contracts, one block, all five sub-claims:
- kernel decides presence, provenance-blind, via the `TagSet.Lookup` `ok` test;
- "an existence atom is decided by the KERNEL from presence alone, and the
  evaluator MUST NOT be consulted for it";
- "a value-comparing atom whose key is ABSENT is marked unevaluable by the
  KERNEL, and the evaluator MUST NOT be consulted for it";
- "Atom verdicts combine in the KERNEL under strong-Kleene three-valued logic …
  The tables are normative; the evaluator holds no part of them.";
- "`Evaluate(atom, value) GuardResult` … It never sees the view."

### Lint-promise-narrowing clause

- RDR 0007 SURVIVOR MEMBERSHIP: "if any surviving candidate row's guard is
  GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`" and "Where
  RDR 0006's exhaustiveness proof and this veto disagree, lint's promise narrows
  (JDR 0001 §JD-4)."
- JDR 0001 §JD-4 + P5: "A green lint means resolution succeeds."
- §JD-4 is open only as to WHICH document records the narrowing. RDR 0006 line
  65 assigns finite-domain exhaustiveness to RDR 0003, so 0003 recording it for
  its own semantics is a legitimate discharge, not a land-grab.

## Rejected / corrected branches

- **"landed in RDR 0007" read as shipped code — REJECTED.** `docs/rdr/README.md`
  lists 0007 as `Final`, 0001 as `Implemented`. JDR 0001 §D1 uses "Lands in
  0007" to mean the document that absorbs the change; 0007 itself re-enters at
  propose. RDR 0002 words the same citation correctly: "0007 is the landing
  document". See `../spikes/iter-2/atom-shape.md`.
- **Positional predicate identity vs 0007's sort tuple — NOT a conflict.**
  RDR 0003's Load-Bearing "Identity" entry is authoring identity; RDR 0007's
  `(RuleID, SourceLocator, key, block, operator, literal)` is payload ordering,
  deliberately non-positional. Different purposes; 0007 does not claim 0003's.

## Open obligation routed HERE by RDR 0007

RDR 0007 A12 is DOWNGRADED and names RDR 0003 as its closing site: RDR 0003
"is silent on the projection" of an existence atom onto the declared-domain
product, and the rule "lands when 0003 states its existence-atom projection,
and rides that document's lock rather than this one."
