Model: claude-opus-5[1m]

# Repeatability-Lite DIFF — RDR 0003 (Guard Predicate Exhaustiveness)

One alternate-model reconstruction (`run-1.md`, Claude Sonnet 5) diffed against
the RDR. This is a lite pass: no same-model consistency sample, no three-run
disagreement count. Findings below are contract silences only — places where a
specific RDR passage (or a specific RDR silence) left an algorithmic contract
under-determined and the run had to invent a resolution.

## 1. Findings

| ID | RDR anchor | The silence | How run-1 resolved it | Why it is a determinacy gap |
| --- | --- | --- | --- | --- |
| D-1 | Technical Design → `Normative Contracts`, clause "Coverage and overlap checks MUST be scoped to a normalized row group … MUST evaluate the participating guard dimensions as one product"; A2 Evidence | **"Participating dimension" is never defined.** A2's derivation says finite domains are needed for "the guard dimensions **that vary** inside that group"; the normative clause and the default-on clause both say "participating"/"all participating dimensions" with no definition. The two readings differ: *varying* excludes a key every row constrains identically; *referenced* includes it. | `participatingDimensions(rows)` — "every guard dimension (tag key) any row in the group references, across all/unless" (step 1). The *referenced* reading. | Changes the scoped product, therefore changes the coverage verdict and the set of dimensions that can trigger the blocking inability-to-prove outcome. A non-varying dimension with no declared domain blocks the whole group under the run's reading and is irrelevant under A2's. Not style — two conforming implementations return different verdicts on the same model. |
| D-2 | `Normative Contracts`, "If a guard dimension lacks a finite declared domain, lint MUST refuse or downgrade an exhaustiveness claim **for that dimension**"; `disposition` row "Guard dimension lacks a finite declared domain" | Scope of the withholding: **per-dimension or whole-group?** The clause says "for that dimension"; the `disposition` table's Lint outcome column says "Refuse or downgrade the claim" unqualified; the `authority` census records one "Exhaustiveness verdict **for a row group**". Nothing says whether a group with one unprovable dimension still yields coverage/overlap findings over the provable remainder. | Whole-group, short-circuit: on the first non-finite dimension the run appends one finding and `return v` immediately — no product, no overlap check, no coverage check. Comment: "whole-group claim withheld; do not proceed to a false green". | Load-bearing on emitted diagnostics. Under the run's reading an author with one undeclared dimension sees exactly one finding and never learns about a real overlap in the same group; under a per-dimension reading they see both. The RDR fixes the *outcome class* (blocking) but not its *scope* or whether it terminates the rest of the proof. |
| D-3 | `Normative Contracts`, "`"Too large to prove"` MUST be a declared, model-independent bound … The refusal diagnostic MUST report the product size it computed and the bound it exceeded" | The RDR **requires** a published bound and a diagnostic carrying `(computed size, bound)`, but names **no bound and no unit** — cardinality of the product? bitset width? assignment count? — and gives no place for the value to be declared. | `buildScopedProduct` returns a bare `error`; the run mints `Finding{Reason: "product too large to prove deterministically"}` carrying neither the computed size nor the bound, and never names a bound value or unit anywhere. | The RDR's own clause says "the same model MUST receive the same verdict on every conforming implementation" — that is unachievable while the bound is undeclared. This is the RDR asserting cross-implementation determinism over a quantity it does not fix. The run silently dropped the two required payload fields because there was nothing to fill them with. |
| D-4 | `Normative Contracts`, "Overlap and coverage diagnostics MUST name the source rule id or context id that contributed each predicate involved in the finding" | The clause is total over *contributing predicates*, but a **coverage gap has no contributing predicate** — it is an absence. The RDR never says what a `graph-coverage-gap` finding names (the group? the uncovered assignments? every row in the group?). The `disposition` row for "Zero rows qualify" names the code and nothing else; the overlap row explicitly says "naming both source rule ids", the coverage row says nothing. | `v.Gaps = append(v.Gaps, Finding{Reason: "coverage gap"})` — `RuleIDs` left empty, no uncovered assignment reported. The run produced a finding that its own `Finding` contract says must carry rule ids. | The asymmetry is in the RDR: overlap's naming duty is spelled out, coverage's is not. An implementer either emits an unattributable gap finding (the run's choice) or invents an attribution rule. Both are conforming; they are not equally useful, and MVV Scenario 2 asserts gaps "fail with source rule/context ids" — which the run's resolution cannot satisfy. |
| D-5 | A13 (`Pending`), Load-Bearing Decisions → Identity (`literal` as an identity-tuple component), operator matrix rows for `in` / `contains` | **The canonical byte spelling of a set literal is unfixed** — the RDR books this as A13 and states the duty is owned here, but the clause does not exist. Set-literal ordering, duplicate handling, and typed-scalar encoding are all open. | `Literal{Scalar any; Set []any}` — an ordered Go slice of `any`, with no canonicalization, no de-duplication, and no ordering rule. `AtomIdentity.Literal` embeds this value directly, so identity-tuple equality inherits slice order. | Named by the RDR itself as a gap that affects `in`, an operator inside the subset this RDR treats as proven. The run's resolution makes `["mid","large"]` and `["large","mid"]` distinct atoms under the identity tuple, which collides with the RDR's source-order-independence requirement (MVV Scenario 5 reorders atoms and requires unchanged findings). Contract-level, not representation taste. |
| D-6 | `disposition` table, rows for "Value atom over an absent key" and "Row group is domain-exhaustive but a participating row can refuse"; `trace` step 8 | **Ordering/compositionality of the narrowing check** relative to coverage and overlap is not fixed. `trace` step 8 places the narrowing after coverage (step 6) and overlap (step 7), but the RDR nowhere says whether a withheld claim suppresses, coexists with, or precedes gap/overlap findings for the same group. | Narrowing runs last (step 7 of the run) but `return v` immediately on the first refusing row — so overlap and gap findings already collected are returned, while any *later* refusing row is never reported. Comment: "green result must not be emitted alongside a withheld one". | The run invented both a position and a short-circuit. The consequence is observable: a group with two refusing rows emits one finding, and which one depends on row iteration order — which the RDR forbids as a behavior determinant ("Successful row matching MUST NOT depend on source order"; MVV Scenario 5). The RDR pins *that a claim is withheld*, never *how many findings that produces*. |
| D-7 | `Normative Contracts`, "A withheld exhaustiveness claim … MUST name the participating row and the atom that can refuse"; the "Peer obligation, not yet agreed (A8)" blockquote | The RDR states the atom-naming half is a **payload extension with no producer** (RDR 0006's finding contract has no atom-level field) and says "the atom-naming half … is a stated requirement with no producer". It does not say what an implementation does *in the meantime*. | `Finding{ RuleIDs []string; Atom *Atom; Reason string }` — the run minted a package-local atom-bearing finding type and noted "RDR 0006 owns the actual finding code contract; this is this package's internal representation before RDR 0006 maps it." | Admissible but lower-weight: the RDR acknowledges the gap and A8 books it. Recorded because the run had to choose an interim representation and chose to carry the field locally rather than drop it — a defensible reading, but the RDR does not state which. |

## 2. Detail — load-bearing and cross-RDR

**D-1 (product membership) is the highest-consequence finding and is
cross-RDR.** The scoped product is the object A2's coverage identity
(`union(row_i accepted) == scoped product`) is stated over, and RDR 0006
consumes the resulting verdict for its blocking findings. If "participating"
means *referenced*, then a group where every row carries `status.eq = "Draft"`
puts `status` in the product as a one-element dimension — harmless for coverage,
but fatal if `status` has no declared domain, since D-2's clause then blocks the
whole group. If "participating" means *varies*, that same group proves fine.
A11's producer request is sized differently under each reading: the *referenced*
reading needs a declared domain for every key any guard touches; the *varies*
reading needs one only for keys that discriminate. The RDR uses both phrasings
in adjacent passages (A2 Evidence line 163 vs. the normative clause at line 599
and the default-on clause at line 507) without reconciling them.

**D-2 and D-6 together determine what a lint run actually prints**, which is the
RDR's own stated observability requirement ("A withheld exhaustiveness claim MUST
be observable, not silent … Emitting nothing MUST NOT satisfy this clause"). The
RDR is careful that *something* is emitted and careful that refuse and downgrade
are one outcome, but it never fixes the cardinality or completeness of the
finding set for a group with multiple defects. The run's two short-circuit
returns are a reasonable engineering choice and an arbitrary one — nothing in the
RDR selects them over accumulate-and-continue. This bears directly on MVV
Scenario 2, which requires "one intentional gap **and** one intentional overlap"
to be detected; under a short-circuiting implementation those must be placed in
separate groups or the test is order-dependent.

**D-3 is the sharpest self-contradiction the run exposed.** The clause demands
model-independent verdicts and a diagnostic reporting `(product size, bound)`,
yet the RDR declares no bound, no unit for "size", and no location for the
declaration (implementation-published? a model field? a constant in this RDR?).
The run complied with the shape and produced an empty payload. Note the
over-spec boundary: the *representation* of the proof (symbolic vs. bitset) is
explicitly left open by the RDR and the run correctly treated `Subset` as opaque
— that is legitimate freedom. The **bound** is not, because the RDR itself
promises cross-implementation determinism over it.

**D-5 is cross-RDR by assignment.** RDR 0007 assigned the canonical set-literal
spelling to this RDR (A13 cites
`docs/rdr/0007-guard-predicate-totality.md:1585-1595`), and A13 is the one
`Pending` record the RDR says is "closable on this RDR's own initiative". The run
independently landed on an ordered, non-canonicalized slice — exactly the
divergence A13's "If wrong" predicts ("two implementations spell the same set
literal differently"). The lite run confirms the gap is reachable by an ordinary
reader, not hypothetical.

## 3. GUESS clusters

The run marks GUESS aggressively; most land on implementation detail the RDR
deliberately leaves open. Grouped by the contract they touch.

### Admissible — GUESSes that land on an algorithmic/normative contract

- **Set-literal encoding** (`Literal.Set []any`, "GUESS concrete shape; RDR only
  fixes 'typed scalar' / 'typed scalar set' / 'typed element set'") → **D-5**.
  The RDR fixes the literal *shape class* but not its canonical spelling, and the
  identity tuple depends on the spelling.
- **Domain shape** (`Domain{EnumValues, ElementSet, IntMin/IntMax}`, "GUESS shape
  — A11/A7/A9 are Pending producer requests") → touches **D-1**. The run had to
  invent the discriminator between "has a finite domain" and "does not"
  (`decl.HasFiniteDomain()`, itself a GUESS), which is the exact predicate D-2's
  clause branches on. The RDR names the four finite kinds but never states the
  predicate over a declaration.
- **Product-too-large error carrying no size/bound** (implicit GUESS — the run's
  `buildScopedProduct` returns bare `error`) → **D-3**.
- **`Verdict` shape** ("GUESS name/shape — RDR fixes only that the outcome is one
  of…") → touches **D-2/D-6**. The run's decision that `Exhaustive` is a single
  boolean over the whole group, rather than per-dimension, is the whole-group
  reading of D-2 encoded in a type.
- **`Finding.Atom *Atom`** ("present when the finding must name a specific
  refusing atom") → **D-7**, interim representation for a payload field with no
  producer.

### Inadmissible — GUESSes on impl detail the RDR deliberately leaves open

Reported here only to show they were considered and dismissed; none is a finding.

- Package path `internal/guard`, and the split of lint proof into a sibling
  package. The RDR states it "owns the evaluator semantics, not command I/O"
  (Existing Infrastructure Audit) and delegates lint authority to RDR 0006;
  package layout is free.
- All Go identifier choices: `Operator`/`OpEq`/`ValueKind`/`KindEnum`/
  `TypedValue`/`Subset`/`Block`/`BlockAll`. Naming is out of scope by rule.
- `EvaluateAtom(atom Atom, value TypedValue) (bool, error)` exact signature. The
  RDR fixes the *contract* ("value semantics over a present value only … MUST NOT
  read the tag view") and the run reproduced it faithfully; the Go shape is free.
- `ParseAtom` existence and signature. The RDR names the well-formedness rules
  and the error kinds, not an entry point. Helper decomposition is free.
- The three named internal helpers (`buildScopedProduct`,
  `rowAcceptedAssignments`, `canRefuse`) as *functions*. Decomposition taste —
  though the semantics they implement are pinned by the RDR and were rendered
  correctly (see §4).
- `SourceLocator{File, Line, Col}` shape — explicitly RDR 0002's, cited not owned.
- `Subset` as an opaque type. The RDR states the representation is an
  implementation choice constrained only to "deterministic symbolic or
  bitset-equivalent"; the run said so and stopped. Correct restraint.
- The run's note that `projectExistsAtom` "is explicitly not yet specifiable"
  because A7 is `Pending`. That is the run reading the RDR correctly, not a gap
  the RDR left — A7 books it with a plan.

## 4. Agreement — contracts the run rendered as the RDR states them

The RDR is determinate on all of the following; the run reconstructed each
without inventing:

- **Seam scope.** The evaluator owns value semantics over a present value only,
  never reads the tag view, and never sees `exists` — reproduced verbatim,
  including the consequence that `guard_unevaluable` is not this package's to
  mint and that per-atom combination is the kernel's (strong Kleene).
- **Atom shape and source identity.** Four fields (`Key`, `Operator`, `Literal`,
  `Block`), with source identity explicitly *not* an atom field and living on the
  enclosing row (`RuleID`, `SourceLocator`). The run flagged this as stated, not
  guessed.
- **Identity vs. semantic equality split.** The run reproduced the total identity
  tuple `(RuleID, SourceLocator, key, block, operator, literal)`, its use for
  diagnostics/dedup, semantic equality `(tag, operator, literal)` for domain
  computation only, and the "never an index within all/unless" prohibition.
- **`unless` semantics.** Conjunctive exclusion block, subtracted as one
  assignment set — not per-atom negation, not source-order priority. Rendered
  correctly in `rowAcceptedAssignments`.
- **Escape rows.** No guard atoms ⇒ denotes the whole product; participates in
  the union like any other row; never excluded from overlap. The run stated all
  three, including "never as a side fact".
- **Product-not-per-dimension.** The run's step 3 comment ("across ALL
  participating dimensions at once — never dimension-by-dimension") is the
  normative clause, and it identified this as the reason
  `buildScopedProduct` is load-bearing.
- **Row-group definition owned here.** Same source state + same recognized
  outcome, the set RDR 0001 resolves exact-one over — taken from this RDR, not
  read back from RDR 0006. The cycle-break in the Technical Design worked.
- **Refuse and downgrade are one outcome.** The run refused to split them into a
  hard/soft tier in both the type (`InabilityToProve []Finding` with "never split
  into hard/soft") and the prose, and correctly enumerated all three input
  classes that route there.
- **The narrowing is real and its test is syntactic.** `canRefuse` decided over
  declarations, "not from a runtime trace or witness input", total, and blocked
  today on A12 + A7 — the run tracked the blocking assumptions accurately.
- **Emitting nothing is not the artifact.** Reproduced as a closing constraint.
- **No wire format owned here.** RDR 0002 owns the TOML container; no
  encode/decode pair; the run cited the illustrative fixture as illustrative.
- **Assumption ledger.** The run correctly identified A7/A9/A11/A12 as Pending
  producer gaps and did not fabricate the missing RDR 0002 fields as existing —
  it modelled `TagDeclaration.Domain` as `*Domain` nil-today with the A11
  citation. It also correctly recorded that the kernel reshape is unshipped
  (`Row.Guard string`).

That list covers the RDR's central contracts. The determinacy of the seam split,
the identity tuple, `unless` semantics, escape-row participation, and the
one-product rule is strong evidence those passages are doing their job.

## 5. Escalation verdict

**Escalate.** D-1 is a load-bearing, cross-RDR silence — "participating
dimension" is used in the normative coverage clause but defined only as "varies
inside that group" in A2's derivation, and the two readings change both the
scoped product and the size of the A11 producer request that RDR 0002 owes this
RDR. D-3 compounds it: the RDR promises cross-implementation determinism over a
bound it never declares. One reconstruction is enough to show these are reachable
by an ordinary reader, but not enough to show which reading a second reader would
pick — the full x3 lens is warranted on the coverage/withholding clauses
specifically.
