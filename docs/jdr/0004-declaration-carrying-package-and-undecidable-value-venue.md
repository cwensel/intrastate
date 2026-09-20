---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: open
cluster: 0012, 0030
labels: intrastate, guard-seam, declaration-carrier, undecidable-value, rdr-cluster
---

# JDR 0004 Which Package May Carry Tag Declarations, and Where Does a Value Admitted at Load but Undecidable at Runtime Refuse?

> One question over one seam. Entries are append-only and never merge, split
> or delete — a cited anchor never moves. Convention lives in
> `$RDR_HOME/jdr/README.md` and is never repeated here; process history lives
> in git and is never narrated here.

*Go identifiers and package names below are **non-normative**; the RDRs own
exact contracts. Source: the 0012×0030 joint-check re-fire against 0030 as a
Final peer, 2026-09-20.*

## Problem statement

Two records restructure one type — `internal/guard/grammar.go::Evaluator`, the
sole implementation of the kernel's `resolve.GuardEvaluator` seam — in ways
that cannot both stand, and each fences its own answer.

`cli/0030` (Final, unimplemented) MOVES the evaluator into `internal/resolve`
so an extracted admitted-cell shim can reach it from the loader without the
`guard → table` import cycle. `cli/0012` (Draft) gives the evaluator STATE — a
key → declared-kind mapping delivered by `NewEvaluator` — retires the
zero-value construction form, and closes its own contract with the sentence
"the kernel (`internal/resolve`) continues to carry no declarations."

After both land, `internal/resolve` holds a type whose field is declaration
state, which `0012:C1` fences against. Neither record can yield unilaterally:
0030's relocation is what makes its loader-side filter reachable at all, and
0012's carrier is what makes a present-but-unparseable value refuse honestly
rather than compare as a raw string.

The second fork rides the same seam. 0012 WIDENS what answers
`GuardUnevaluable` at runtime — an `int`-declared tag holding `many` becomes
undecidable where it was decidably false. 0030 fixed the undecided arm of its
load-time admitted-cell filter to EXCLUDE, and argued (`cli/0030:A4`) that the
arm is unreachable for the kinds it admits. That argument is sound against
today's evaluator and silent about 0012's. A value can now be admitted at load
and undecidable at runtime, and no record says which venue owns it: 0030 routes
load-time failures to `malformed_tag_declaration` / `malformed_predicate_atom`
/ `malformed_initial_declaration`, 0012 routes runtime ones to
`guard_unevaluable`, and the new overlap belongs to neither by any written
rule.

Answered per-RDR, the drift is concrete: the loader gets a declaration-aware
comparison its own record never specified, or 0030's zero-cell refusal fires on
correct models because a nil-mapping evaluator excludes every cell.

## Principles

1. **The kernel decides no value semantics.** `0007:C1` places presence and
   atom combination in the kernel and delegates value comparison to the seam:
   the evaluator "decides value semantics over a PRESENT value and never reads
   the tag view" (`0003:C9`). Any resolution that moves comparison policy into
   kernel control is out.
2. **A present value an operator cannot parse is unevaluable, never false.**
   `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` enforces it
   as a cross-RDR obligation; `0007:C2` fixes the direction "never to false and
   never to true." No resolution may collapse an undecided arm to a decided
   verdict on either side of the seam.
3. **One comparison, three third-arm policies.** `cli/0030:C1` establishes that
   the extraction shares the two-valued core only — `guard` passes the verdict
   through, `graphlint` admits on undecided, the loader excludes — because "what
   'as the runtime evaluator would' buys is a shared comparison, never a shared
   undecided-arm policy." A resolution that unifies the third arm breaks a
   locked clause.
4. **A locked record's clause is not amended by a peer's proposal.** 0030 is
   Final. Its contracts change only by a Final → Draft flip, never by 0012
   asserting a different rule (`$RDR_HOME/stages/07.1-cluster-reconcile.md`
   HARD RULE).

---

## D1 — May `internal/resolve` carry tag declarations?

`cli/0012:C1` fences that it may not ("the kernel (`internal/resolve`)
continues to carry no declarations") while giving `Evaluator` a declaration
field. `cli/0030:C1` relocates that same `Evaluator` into `internal/resolve`.
The two clauses are jointly unsatisfiable, and each is load-bearing for its own
record's central mechanism.

- **(a) 0012 keeps the no-declarations fence; 0030's relocation is reverted.**
  Restores 0012's clause verbatim, but the loader can no longer reach the
  shared comparison, so `cli/0030:C1`'s single-source extraction fails and the
  record mints the fourth copy of the shim it exists to prevent — in the
  loader, where a divergence is a load refusal rather than a lint finding.
  Requires a Final → Draft flip on 0030. Ruled out by nothing, but it pays
  0030's whole reuse argument to preserve a sentence about package hygiene.

- **(b) Narrow the fence to its intent: the kernel's RESOLUTION PATH carries no
  declarations; a seam implementation co-located in the package may — recommended.**
  0012's sentence is a statement about the kernel as a decider, not about the
  package as a namespace: its own C1 preserves "never reads the tag view" as a
  property of the type by admitting only load-time model data, never a
  `resolve.TagSet` or runtime value. Relocating a seam implementation that
  already satisfies that restriction does not make the kernel declaration-aware
  — `resolve.Resolve` and its callees still read none. Satisfies Principle 1
  because comparison policy stays in the seam type, not in kernel control flow.
  Costs 0012 a narrowed clause at propose, before it locks.

- **(c) Relocate to a fourth package below both, leaving `resolve` untouched.**
  Preserves both clauses literally. But `internal/resolve` imports nothing
  intra-repo and already owns `GuardAtom`, `GuardResult`, `GuardEvaluator` and
  the conformance harness, so a new package would hold one type and split the
  seam's definition from its implementation across two homes for a naming
  reason. Ruled out by no principle; rejected on conceptual integrity — it
  makes the seam harder to find in order to keep a sentence true.

**Resolved: (open).**

*Lands in **0012**: C1's no-declarations sentence is narrowed to the resolution
path. **0030** cites the decision rather than restating it; if the resolution
is (a) or (c), 0030 flips Final → Draft and re-reconciles C1.*

---

## D2 — Where does a value admitted at load but undecidable at runtime refuse?

`cli/0012:C2` widens the runtime undecidable set (an `int`-declared tag holding
an unparseable value becomes `GuardUnevaluable` rather than decidably false).
`cli/0030:A4` argues the undecided arm is unreachable for the kinds its C1
admits, and `cli/0030:C1` fixes that arm to EXCLUDE as a soundness fence. Under
0012 the arm becomes reachable, and the two records route the same condition to
different venues — 0030 to its three load-refusal categories, 0012 to
`guard_unevaluable` — with no rule saying which owns the overlap.

- **(a) Load refuses it: the loader rejects a model whose declared domain
  contains a member the typed comparison cannot decide.** Catches the defect
  at authoring time, but converts an undecidable predicate into "your model is
  broken" — the direction `cli/0030:C1` explicitly rejected for this arm, and
  it would make 0012's runtime widening unreachable by construction.

- **(b) Runtime refuses it; load excludes the cell and says so — recommended.**
  Keeps each venue deciding what it can see: the loader excludes an undecidable
  cell (0030's settled direction, leaving the cell to the outcome group's
  ordinary `graph-coverage-gap`), and a runtime value the typed comparison
  cannot decide surfaces as `guard_unevaluable` (0012's). Satisfies Principle 2
  on both sides — neither venue collapses undecided to decided. Costs 0030's A4
  its "unreachable" premise, which becomes "reachable under 0012, and the fence
  is live" — a strengthening of the record, not a refutation, since the arm's
  disposition was already written.

- **(c) Leave it to each record.** The status quo, and the reason this entry
  exists. Ruled out by Principle 4's corollary: the venue is a shared decision,
  and a shared decision declared in two places drifts.

**Resolved: (open).**

*Lands in **0030**: A4's reachability premise is restated as conditional on the
evaluator's typing, with C1's exclusion unchanged. **0012** cites the venue
split rather than restating 0030's arm.*

---

## D3 — Does the zero-value `Evaluator{}` construction form survive?

`cli/0012:C1` retires it: "every non-test construction goes through
`NewEvaluator`", and a bare literal answers `GuardUnevaluable` for every
`eq`/`in` atom. `cli/0030`'s extracted shim constructs `Evaluator{}` at its
call site. Under 0012, 0030's admitted-cell filter would receive
`GuardUnevaluable` for every candidate and — per its own EXCLUDE arm — admit
zero cells, firing `cli/0030:C1`'s zero-cell load refusal on correct models.

- **(a) The shim constructs through `NewEvaluator` with the loader's own
  declaration map — recommended.** The loader holds the declarations at the
  point it runs the filter, so the mapping is available without a new carrier.
  Makes the load-time comparison declaration-aware, which is what 0012 argues it
  should have been all along. Requires 0030's Phase 2 step to name the
  construction, which is implementation detail under a clause that already
  fences the shim's signature.
- **(b) The shim keeps a nil-mapping evaluator and the loader's filter stays
  raw-string.** Preserves 0030's step literally but reintroduces exactly the
  untyped comparison 0012 exists to remove, in the one venue where a wrong
  answer is a load refusal. Ruled out by Principle 2.

**Resolved: (open).**

*Lands in **0030**: the Phase 2 extraction step names the construction. **0012**
gains the loader as a declared construction site.*

## Interface record

- **JD-1** `open` — Whether `internal/resolve` may host a seam implementation
  holding declaration state, and how `0012:C1`'s no-declarations sentence is
  scoped. A reader of either record must be able to tell whether the kernel
  package is declaration-free as a namespace or as a decision path.
  *(Provenance: D1, opened by the 0012×0030 joint-check re-fire, 2026-09-20.)*
  **Binds:** 0012, 0030.
- **JD-2** `open` — The refusal venue for a value admitted at load and
  undecidable at runtime. A model author needs one answer to "why did this
  refuse, and where do I fix it."
  *(Provenance: D2, same re-fire.)* **Binds:** 0012, 0030.
- **JD-3** `open` — The construction form the loader-side admitted-cell shim
  uses, and whether the load-time comparison is declaration-aware.
  *(Provenance: D3, same re-fire.)* **Binds:** 0012, 0030.

## What this does not decide

- **The admitted-cell rule itself** — which domain members a step write expands
  to, and the EXCLUDE direction of its undecided arm. `cli/0030:C1` owns it;
  this registry decides only when that arm is reachable (D2).
- **The typed-comparison table** — which kinds parse how, and what each answers
  on an unparseable held value. `cli/0012:C2` owns it.
- **The import-cycle fact** — that `guard` imports `table` and `table` cannot
  import `guard`. Settled in source and cited by both records
  (`cli/0012:A5`, `cli/0030:A12`); not a fork.
- **The shared two-valued core and its three call-site policies.**
  `cli/0030:C1` owns it; Principle 3 protects it rather than reopening it.
- **The five kind tokens** (`bool`/`enum`/`int`/`scalar`/`set`). Homed at
  `cli/0003:C6` by the existing 0012×0030 joint-check; neither record extends
  the vocabulary.
