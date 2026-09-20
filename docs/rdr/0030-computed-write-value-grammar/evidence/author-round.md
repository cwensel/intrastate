Model: claude-opus-5[1m]

# cli/0030 — author's round

Stage 4 (Resolve) ran delegated, so this round could not be put to the
author in session. Each item carries the grounding it takes to rule on it.
Answers land in `rulings.md` beside this file (rdr-common §run-prompt); the
re-run reads them before the record and finishes the round.

Fixtures: none rendered. No assumption in this record turns on an exact
expected value — the MVV compares two runs to each other, not either to a
literal. The round runs anyway for the four questions below.

## 2026-09-20 — resolve

- **Q1 — A8 is refuted: `flow next` / `flow resolve` publish the AUTHORED
  rule id, never the expansion suffix. What does this record owe?** — The
  record assumed (A8) that an expanded row's suffixed identity is already
  what the flow verbs show, so "no new rendering surface is owed." Source
  says otherwise: the suffix lives in a separate `Row.Suffix` field and is
  joined only by `internal/table/model.go::Identity`;
  `internal/table/model.go::Row.KernelRow` builds `resolve.Row{RuleID:
  r.RuleID, …}` and drops `Suffix`, so `internal/cli/flow_next.go::summarize`
  and `internal/cli/flow_resolve.go` publish the bare id. An `in`-expanded
  row surfaces TODAY as `retry`, not `retry#3`. The dump half of A8 does
  hold (`internal/table/dump.go` renders `Identity()` and `SourceLocator`).
  This bears directly on the Decision Rationale's accepted cost — "the
  normalized dump and `flow next` show `retry#0 … retry#4`, not `retry` —
  the same view an `in` expansion already gives, accepted on that
  precedent" — which is true of the dump and false of `flow next`. Three
  dispositions, and the choice is the author's: (a) C3 gains a clause
  obliging the flow payloads to carry `Identity()`, making this record the
  one that fixes a gap `in` expansion already has; (b) the cost is accepted
  as-is and the Rationale sentence is corrected to claim only the dump,
  leaving the flow verbs naming a rule whose row cannot be identified;
  (c) out of scope — a separate record owns the flow-payload identity.
  — grounding: `internal/table/model.go::Row.KernelRow`,
  `internal/cli/flow_next.go::summarize`, `internal/cli/flow_resolve.go`,
  `internal/table/dump.go::column`; contrast `internal/table/model.go::Identity`.
  Corpus rung ran and found nothing: `searched=corpus; found: none` —
  no prior-art treatment of authored-vs-generated identity in tool output
  across `StateMachineRes`, `DevRef`.

- **Q2 — the same read exposes an ambiguous join:
  `internal/cli/flow_resolve.go::rowByID` matches on bare `RuleID`. In
  scope or out?** — N expanded rows share one authored `RuleID`, so that
  lookup returns the first. This is already true for `in` expansion; this
  record multiplies the rows it applies to. It is a pre-existing defect
  this record's verification surfaced, not one it introduces, and rdr-common
  doctrine is that a record records the past rather than absorbing adjacent
  work. Rule: name it in Failure Modes and file it separately, or fold it
  into C3 alongside Q1(a) if Q1(a) is the answer there.
  — grounding: `internal/cli/flow_resolve.go::rowByID`; the same expansion
  arity that makes `Row.Suffix` load-bearing.

- **Q3 — `Evaluator.Evaluate` is THREE-valued; C1 disposes of two arms.
  What is an admitted cell when the evaluator answers `Unevaluable`?** —
  A4 verified that the load-time admits filter can replay the runtime
  comparison exactly (`internal/guard/grammar.go::Evaluator.Evaluate` is
  stateless, a pure function of operator, atom literal and candidate value).
  But `Evaluate` returns True / False / **Unevaluable** — an unparseable
  literal or value yields the third, not false. C1 says a cell is admitted
  when it "satisfies the CONJUNCTION of every positive atom … evaluated per
  member as the runtime evaluator would," which is silent on the third arm.
  Three readings, each defensible and each a different refusal surface:
  admit the cell (over-admit, and the bound refusal C2 catches the
  consequence), exclude it (under-admit, and the cell becomes an ordinary
  `graph-coverage-gap`), or refuse the model at load (the `0003:C6`
  disposition this record already takes for a zero-cell rule). Note the
  adjacent fact that makes this reachable rather than theoretical:
  `internal/table/load.go::conform` deliberately does NOT domain-check
  `lt`/`lte`/`gt`/`gte` bounds (`0002:C17`), so a guard bound outside the
  declared domain is legal today and the admits filter must not assume
  otherwise.
  — grounding: `internal/guard/grammar.go::Evaluator.Evaluate` (three-valued
  `GuardResult`); `internal/table/load.go::conform` (no domain check on
  ordered operators); `0003:C6` (unsatisfiable atom refused at load).
  Corpus rung ran and found nothing: `searched=corpus; found: none` —
  no treatment of three-valued-evaluator reuse at load, or of the soundness
  direction for an undecided predicate, in `StateMachineRes`, `PapersFast`.

- **Q4 — the admitted-cell evaluator already exists in `internal/guard`.
  Cite it, extract it, or diverge?** — The Technical Design describes
  per-cell evaluation as new work and the Existing Infrastructure Audit has
  no row for it. In fact `internal/guard/product.go::Denotation(m, key,
  atom)` returns exactly "the subset of key's dimensions the atom denotes,"
  enumerating the declared domain and filtering each member through
  `internal/guard/product.go::valueSatisfies`, whose own comment says it "is
  the same decision the runtime evaluator makes, over the same rendered
  value form" — verbatim the semantics C1 specifies. The conjunction step
  exists too: `internal/guard/product.go::acceptedIn` intersects
  `Denotation` across a row's `guard.all` atoms and segregates `unless`
  into a separately-subtracted term, which is precisely C1's "`unless` atoms
  are not consulted." Int domain enumeration is
  `internal/guard/declaration.go::IntDomain`. The obstacle is placement, not
  semantics: these live in `internal/guard`, which consumes a NORMALIZED
  model, and `guard` imports `table`, so `table` cannot import `guard`
  without a cycle (the same package-graph fact that makes A7's
  `AssignmentCount` unreachable from load). `Denotation` also returns lint's
  richer `AssignmentSet` and yields the unprojectable empty set on an
  unevaluable dimension — a lint-specific false-green guard, not a load
  refusal (which is Q3 again, from the other side). Rule: extract the
  shared predicate to a package both can import, duplicate it in `table`
  with a contract test pinning the two together, or diverge deliberately —
  and note that divergence silently breaks C1's "as the runtime evaluator
  would" claim. Whichever, the audit table owes a row naming
  `Denotation` / `valueSatisfies` rather than describing a new evaluator.
  — grounding: `internal/guard/product.go::Denotation`, `::valueSatisfies`,
  `::acceptedIn`; `internal/guard/declaration.go::IntDomain`; the import
  direction `guard` → `table`.
