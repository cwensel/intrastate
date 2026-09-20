---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: settled
seam:
  - internal/guard/grammar.go::Evaluator
  - internal/guard/product.go::valueSatisfies
  - internal/graphlint/reach.go::atomAdmitsValue
  - internal/cli/flow_resolve.go::guardSeam
  - internal/resolve/guardcontract.go::TestGuardEvaluatorContract
cluster: 0012, 0030 — members as of iteration 1; a VIEW, not the identity
labels: intrastate, guard-seam, declaration-carrier, evaluator-placement, rdr-cluster
---

# JDR 0004 Where Does the Value-Comparison Evaluator Live, and What May It Carry?

> One question over one seam. Entries are append-only and never merge, split
> or delete — a cited anchor never moves. Convention lives in
> `$RDR_HOME/jdr/README.md` and is never repeated here; process history lives
> in git and is never narrated here.

*Go identifiers and package names below are **non-normative**; the RDRs own
exact contracts. Source: the 0012×0030 joint-check re-fire against 0030 as a
Final peer, 2026-09-20.*

## Problem statement

Two records restructure one type — `internal/guard/grammar.go::Evaluator`, the
sole implementation of the kernel's `resolve.GuardEvaluator` seam — and each
fences its own answer without seeing the other's.

`cli/0030` (Final, unimplemented) MOVES the evaluator into `internal/resolve`
so an extracted admitted-cell shim can reach it from the loader without the
`guard → table` import cycle. `cli/0012` (Draft) gives the evaluator STATE — a
key → declared-kind mapping delivered by `NewEvaluator` — retires the
zero-value construction form, and closes its contract with the sentence "the
kernel (`internal/resolve`) continues to carry no declarations."

Read as package purity, those two cannot both stand: after both land,
`internal/resolve` holds a type whose field is declaration state. Neither
record can yield unilaterally — 0030's relocation is what makes its
loader-side filter reachable at all, and 0012's carrier is what makes a
present-but-unparseable value refuse honestly rather than compare as a raw
string.

The stake is concrete and asymmetric. Under 0012, a bare `Evaluator{}` answers
`GuardUnevaluable` for every `eq`/`in`; under 0030's EXCLUDE arm that admits
zero cells and fires `cli/0030:C1`'s zero-cell load refusal on CORRECT models,
with no diagnostic pointing at the cause. Answered per-RDR, whichever
implements second inherits that trap silently.

## Principles

1. **The kernel decides no value semantics.** `0007:C1` places presence and
   atom combination in the kernel and delegates value comparison to the seam;
   `0003:C9` fences the evaluator as deciding "value semantics over a PRESENT
   value" that "never reads the tag view." Any resolution moving comparison
   policy into kernel control flow is out.
2. **A present value an operator cannot parse is unevaluable, never false.**
   `0007:C2` fixes the direction — "never to false and never to true" — and
   `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` enforces it
   across implementers. No resolution may collapse an undecided arm to a
   decided verdict on either side of the seam.
3. **One comparison, three third-arm policies.** `cli/0030:C1` shares the
   two-valued core only — `guard` passes the verdict through, `graphlint`
   admits on undecided, the loader excludes — because "what 'as the runtime
   evaluator would' buys is a shared comparison, never a shared undecided-arm
   policy." A resolution unifying the third arm breaks a locked clause.
4. **A locked record's clause is not amended by a peer's proposal.** 0030 is
   Final. Its contracts change only by a Final → Draft flip, never by 0012
   asserting a different rule (`$RDR_HOME/stages/07.1-cluster-reconcile.md`,
   HARD RULE).

---

## D1 — May `internal/resolve` host a seam implementation that carries declarations?

`cli/0012:C1` fences that it may not ("the kernel (`internal/resolve`)
continues to carry no declarations") while giving `Evaluator` a declaration
field. `cli/0030:C1` relocates that same `Evaluator` into `internal/resolve`.
The clauses are jointly unsatisfiable under a package-purity reading, and each
is load-bearing for its own record's central mechanism.

- **(a) 0012 keeps the fence as package purity; 0030's relocation is reverted.**
  Preserves the sentence literally, but the loader can no longer reach the
  shared comparison, so `cli/0030:C1`'s single-source extraction fails and the
  record mints the fourth copy of the shim it exists to prevent — in the
  loader, where a divergence is a load refusal rather than a lint finding.
  Requires a Final → Draft flip on 0030.

- **(b) The fence scopes to the kernel's RESOLUTION PATH; a co-located seam
  implementation may carry declarations — recommended.** `JDR 0001 §D1` already
  fixed what the fence means: "The kernel carries *cardinality*, not *meaning*
  — it stays as grammar-agnostic as it already is for `Match []Tag` and
  `Writes []Tag`." That is a claim about what the kernel DECIDES, not what the
  package CONTAINS. `JDR 0001 §D4` resolved the matching half — "the evaluator
  decides value comparisons over present keys and never sees the view" — making
  the evaluator's independence a property of the TYPE, which relocation does not
  touch. Satisfies Principle 1: comparison policy stays in the seam type, and
  `resolve.Resolve` and its callees read no declaration.

- **(c) Relocate to a fourth package below both, leaving `resolve` untouched.**
  Preserves both clauses literally, but `internal/resolve` imports nothing
  intra-repo and already owns `GuardAtom`, `GuardResult`, `GuardEvaluator` and
  the conformance harness. A new package would hold one type and split the
  seam's definition from its implementation across two homes for a naming
  reason.

**Resolved: (b).** `internal/resolve` may host the value-comparison seam's
implementation, declaration state included. The no-declarations rule governs the
kernel's resolution path — `Resolve` and its callees — not the package as a
namespace. The package already hosts non-kernel code for precisely this reason:
`guardcontract.go` is a `testing`-importing contract harness placed there so two
records could reach it across this same cycle, so hosting-by-reachability is the
established pattern rather than an exception minted here.

*Lands in **0012**: C1's closing sentence is scoped to the resolution path.
**0030** cites this entry; its C1 relocation stands unchanged and needs no
re-lock.*

---

## D2 — Where does a value admitted at load but undecidable at runtime refuse?

`cli/0012:C2` widens the runtime undecidable set: an `int`-declared tag holding
an unparseable value becomes `GuardUnevaluable` where it was decidably false.
`cli/0030:A4` argues the undecided arm is unreachable for the kinds its C1
admits, and `cli/0030:C1` fixes that arm to EXCLUDE as a soundness fence. The
fork as filed asked which venue owns a value that is admitted at load and
undecidable at runtime.

**Resolved: no conflict — the two venues quantify over disjoint populations.**
Tracing every ingress by which a value reaches a comparison:

| Ingress | Held to the declaration? | Where |
| --- | --- | --- |
| Model literals (`[[rule]]` atoms, `[initial]`) | yes — kind and domain | `internal/table/load.go::conform`, `::conformKind`, `::conformDomain` |
| CLI `--write` / `--tag` | yes — kind and domain | `internal/table/load.go::ConformValue`, via `internal/cli/flow_input.go::canonicalValue` |
| Reader-supplied owned values | **no** | nothing in `internal/accessor` conforms |

0030's admitted-cell filter enumerates the DECLARED DOMAIN and never sees a
held value; 0012's widening bites only on the third row, a reader handing back
a value never held to its declaration. No predicate is therefore both
load-admitted and runtime-undecidable: the loader decides over declared
members, the runtime over supplied values. Each venue refuses what it can see —
the loader excludes an undecidable cell, leaving it to the outcome group's
ordinary `graph-coverage-gap`; the runtime answers `guard_unevaluable` —
and Principle 2 holds on both sides because neither collapses undecided to
decided.

`cli/0030:A4` is therefore correct AS WRITTEN and is not narrowed by 0012: its
scope is declared domain members, which are conformed at load. The entry is
recorded rather than withdrawn so a later reader does not re-derive the
question, and so the disjointness is citable when a fourth ingress is proposed.

*Lands in **neither record**. Both stand unchanged; 0012 may cite this entry
where it states the widening's reach.*

---

## D3 — Does the zero-value `Evaluator{}` construction form survive?

`cli/0012:C1` retires it: "every non-test construction goes through
`NewEvaluator`", and a bare literal answers `GuardUnevaluable` for every
`eq`/`in` atom. `cli/0030`'s extracted shim constructs `Evaluator{}` at its
call site. Under 0012, 0030's admitted-cell filter would receive
`GuardUnevaluable` for every candidate and — per its own EXCLUDE arm — admit
zero cells, firing `cli/0030:C1`'s zero-cell load refusal on correct models.

- **(a) The shim constructs through `NewEvaluator` with the loader's own
  declaration map — recommended.** `internal/table/normalize.go::renderWrites`
  is a `*loader` method and holds `l.model.Tags` at the call site — the same
  access it already uses to read a write key's declaration — so the mapping is
  available with no new carrier and no plumbing. Makes the load-time comparison
  declaration-aware, which is what 0012 argues it should be.
- **(b) The shim keeps a nil-mapping evaluator and the loader's filter stays
  raw-string.** Preserves 0030's step literally but reintroduces the untyped
  comparison 0012 exists to remove, in the one venue where a wrong answer is a
  load refusal rather than a lint finding. Ruled out by Principle 2.

**Resolved: (a).** The loader-side admitted-cell shim constructs through
`NewEvaluator`, supplying the declarations the loader already holds. The
load-time comparison is declaration-aware, and no construction site anywhere
relies on the retired zero-value form.

This is the entry with a correctness trap rather than a definitional question:
resolved the other way, the failure is a load refusal on a correct model with
nothing naming the cause. This constrains 0030's implementation whichever way D1 had
landed.

*Lands in **0030**: the Phase 2 extraction step names the construction.
**0012** gains the loader as a declared construction site.*

## Interface record

- **JD-1** `decided` — `internal/resolve` may host the value-comparison seam's
  implementation, declaration state included; the no-declarations rule governs
  the kernel's RESOLUTION PATH, not the package as a namespace. A reader of
  either record can tell that the kernel is declaration-free as a decision
  path while the package hosts the seam that types comparisons.
  *(Provenance: §D1, opened by the 0012×0030 joint-check re-fire, 2026-09-20;
  scope fixed by `JDR 0001 §D1` and `§D4`.)*
  **Binds:** 0012, 0030.
- **JD-2** `decided` — Load-time admission and runtime comparison quantify over
  disjoint populations: the loader over declared domain members (conformed at
  load), the runtime over supplied held values (reader values unconformed).
  Each venue refuses what it can see; neither inherits the other's undecided
  arm. A model author gets one answer to "why did this refuse, and where do I
  fix it."
  *(Provenance: §D2, same re-fire; ingress census against `main`.)*
  **Binds:** 0012, 0030.
- **JD-3** `decided` — The loader-side admitted-cell shim constructs through
  `NewEvaluator` with the loader's own declarations; the retired zero-value
  form has no remaining construction site. Prevents a zero-cell load refusal
  firing on a correct model.
  *(Provenance: §D3, same re-fire.)*
  **Binds:** 0012, 0030.

## What this does not decide

- **The admitted-cell rule itself** — which domain members a step write expands
  to, and the EXCLUDE direction of its undecided arm. `cli/0030:C1` owns it;
  this registry decides only where that comparison is constructed (§D3) and
  that the arm's population is disjoint from the runtime's (§D2).
- **The typed-comparison table** — which kinds parse how, and what each answers
  on an unparseable held value. `cli/0012:C2` owns it.
- **Whether reader-supplied values should be conformed at all.** §D2 records
  that they are not; whether that is a defect is nobody's decision here.
  `cli/0012` reaches it at runtime by construction, and a record proposing
  load-time or accessor-time conformance would own the change.
- **The import-cycle fact** — that `guard` imports `table` and `table` cannot
  import `guard`. Settled in source and cited by both records
  (`cli/0012:A5`, `cli/0030:A12`); not a fork.
- **The shared two-valued core and its three call-site policies.**
  `cli/0030:C1` owns it; Principle 3 protects it rather than reopening it.
- **The five kind tokens** (`bool`/`enum`/`int`/`scalar`/`set`), homed at
  `cli/0003:C6` by the existing 0012×0030 joint-check; neither record extends
  the vocabulary.
