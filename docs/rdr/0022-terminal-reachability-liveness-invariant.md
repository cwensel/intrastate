# Recommendation 0022: Terminal-reachability liveness invariant

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
  <!--
  - `Deferred` is the parked-with-a-revisit-trigger status for a
    Draft that cannot proceed because **no acceptable mechanism
    exists yet** — every in-our-control path is ruled out and the
    one that would work is outside our control. It is a *pause in
    the lifecycle*, not an exit from it: the RDR stays intact and
    re-enters at the stage it stopped when the trigger fires.
    Carry the condition on the live value:
    `Deferred [revisit when <condition>]`, and say in the same
    field what was ruled out and why (Alternatives Considered
    carries the long form). Distinct from `Abandoned`, which is
    terminal — an Abandoned RDR is closed, owes a post-mortem, and
    never re-enters. A Deferred RDR owes **no** post-mortem
    (nothing was implemented), and its `Priority` records what the
    fix is *worth*, not what is scheduled. Do not defer merely to
    park work that is possible but unfunded — that is a Priority,
    not a Status.
  - `Demoted` is the terminal status for an RDR judged
    *not RDR-shaped* — the decision was never a real
    design fork, so it leaves the RDR lifecycle and is
    refiled as a plain issue. Carry the destination on the
    live value: `Demoted [→ <issue link>]`, and record the
    same link under **Related Issues**. A `Demoted` RDR runs
    no further stages. (Distinct from the 07.1 *demotion*
    below, which is a `Final → Draft` flip that keeps the
    RDR in the lifecycle — that flip never writes
    `Status: Demoted`; see the disambiguation note there.)
  - A Draft demoted from Final by the 07.1 cluster gate
    carries a qualifier on the live value:
    `Draft [revised from Final YYYY-MM-DD; re-verify A2,A4
    — <one-line reason>]`. It is still a `Draft` for every
    binary Draft/Final gate; only Stage 4 (scoped
    re-verify) and Stage 7 (re-lock) parse the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Final tolerated at the 07.1 gate under a JOINT-DECISION
    carries `Final [joint decision → <home §-anchor>: <the
    open question>]`. It is still a `Final` for every binary
    gate. The qualifier is an **open obligation, not a
    coherence claim**: it says the named question is
    unanswered here, not that this RDR agrees with the answer.
    So it does not self-clear. When the home answers, this RDR
    owes a scoped check of that answer against its own
    normative fences before it re-locks or implements —
    consistent → drop the qualifier and record the clearing;
    contradicts fenced text → a 07.1 SPEC-DEFECT. A re-lock
    that comes first carries the qualifier forward unchanged;
    it is never silently dropped.
  -->
- **Type**: Feature
- **Profile**: large — provisional: one contract (the
  `graph-terminal-unreachable` finding class) extending locked
  RDR 0006's finding taxonomy — an enum surface.
  <!-- Do not paste the matrix below into the field; it is the
  Stage 5 routing latch, provisional on `Draft`, made
  authoritative by Resolve.
  Sized by BLAST RADIUS — the MAX of two axes, not
  contract count or word count.
  (1) contract axis: small = one contract, no user-facing
  surface (skips Stage 5); mid = one contract + user-facing
  surface OR locks a contract; large = locks an enum/hash/
  format/grammar/destructive-op; foundational = cross-RDR
  producer / spans modules.
  (2) accretion axis (HARD floor): if `Seam Lineage` below
  carries ≥2 closed prior point-fixes at this locus, Profile
  is floored at FOUNDATIONAL regardless of the contract axis
  — a seam with prior point-fixes is never small/mid (it
  spans the prior RDRs/patches = the matrix's cross-RDR
  trigger). The only escape is a written accretion disposition
  in the Seam Lineage field. This floor is what stops a
  "one contract → mid" sizing from under-gating an accreting
  seam.
  Matrix: rdr/stages/README.md. Seed estimates from the design
  shape; Resolve overwrites from the verified count; Stage 8
  Gate locks it at Draft → Final. Never skip lenses off a
  Draft Profile until Resolve has run. -->
- **Priority**: Low
- **Related Issues**: kata `intrastate#4hps` (1603); kata `jjkh`
  (graph export — this invariant needs only the in-memory
  relation, no dependency on the export landing)
- **Predecessors**: 0006-graph-lint-authority-and-guarantees
- **Seam Lineage**: no prior accretion

## Problem Statement

A model author whose graph legally contains cycles (revise self-loops,
rewind backedges — both present in `models/rdr.toml`) needs lint to
catch livelock: a reachable state, or cycle of states, from which no
declared terminal is reachable. `graphlint` proves safety invariants
over the declared finite product — determinism, totality, dead-end,
unreachable-rule, terminal-escape (`internal/graphlint/taxonomy.go`) —
but not the liveness property "from every reachable state some
terminal is reachable" (CTL `AG EF terminal`). `graph-dead-end` cannot
detect it: every member of a livelocked cycle has an outgoing row.
This is the one property class an external model checker (TLC/NuSMV)
would add; review recommended adding it as one more in-house
reachability pass — reverse-reachability from the terminal set over
the graph `reach.go` already builds — rather than adopting a second,
drifting formal artifact.

The fork is real on four axes: (1) **quantifier/abstraction** — `EF`
over the merged fixpoint graph is an existential path claim on a
widened relation, the false-green direction for this check (a merged
node may show a terminal path no concrete run has); choose between
splitting on terminal-participating keys (invariant 2's remedy),
accepting a pinned miss (the D12 pattern, cf. RDR 0015), or rejecting
path-sensitivity as exponential; (2) **granularity** — per reachable
state vs SCC-collapsed, so one finding names the cycle rather than
every member; (3) **authority** — blocking vs advisory, given cycles
are legal by design and the advisory tier is closed at four members;
(4) **refusal posture** — when a participating dimension is unprovable
or the product too large, downgrade to the existing
`graph-unprovable-coverage`/`graph-product-too-large` path rather than
guess. RDR 0006 is Implemented and never amended, so this lands as a
successor RDR extending the finding taxonomy. In scope is existential
terminal reachability only; fairness/eventuality under specific
outcome sequences (`AF`, true LTL liveness) is explicitly out, since
outcome choice is caller-driven.

## Critical Assumptions

- **A1 A split node's row-successor re-anchors into the merged fixpoint:
  for every reachable split node `s` and satisfiable ordinary row `r`,
  some merged fixpoint node with `successor(s, r)`'s presence footprint
  subsumes it, so the split closure stays inside the split forms of the
  already-computed node set and is finite.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: [needed: a spike driving
    `internal/graphlint/reach.go::successor` from split nodes of a cyclic
    fixture and asserting `subsumes` finds a host node for every produced
    successor]
  - **If wrong**: the closure escapes the computed node set — the pass
    either fails to terminate on the split space or accuses states the
    fixpoint never admitted; surfaces as a hang or a phantom finding on
    the legal-cycle fixture.
- **A2 `models/rdr.toml` passes the new check unchanged: every reachable
  split node of its revise self-loops and rewind backedges reaches a
  node satisfying `landed` or `abandoned` under the split relation.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: [needed: `intrastate lint` over `models/rdr.toml` with
    the pass enabled, findings unchanged from today's run]
  - **If wrong**: the blocking gate turns red on the checked-in model —
    either a genuine model defect (fix the model) or a false accusation
    (reopen the quantifier choice); surfaces as CI failure on an
    untouched model.
- **A3 Reusing invariant 2's terminal-satisfaction semantics verbatim
  over a terminal-participating key with no finite declared domain (the
  opaque carrier; writer: `internal/graphlint/reach.go::heldValues`,
  gated on `guard.AssignmentCount`) does not make this check reject a
  model class invariant 2 passes today — the evaluator's verdict for
  `OpaqueValue` under a value atom decides whether such terminals are
  satisfiable.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: [needed: a fixture whose terminal carries a value atom
    over an unbounded scalar key, linted with the pass enabled; record
    whether the model is accused, and the
    `internal/guard::Evaluator.Evaluate` verdict for the opaque value]
  - **If wrong**: value-atom terminals over unbounded keys empty the
    terminal set and the whole cyclic graph is accused; Resolve routes
    that case to the axis-4 refusal (an unprovable-dimension downgrade)
    instead of the verbatim reuse.
- **A4 The accused population and `graph-dead-end`'s population are
  disjoint and jointly cover the bad stuck states: the closure's node
  universe is a subset of the `splitNode` forms `checkDeadEnd` already
  iterates, so no split node exists that only this pass sees (no orphan
  claimed by neither check); and a trapping sink component's every
  ordinary edge stays inside the component (an edge to a good node
  would contradict badness), so requiring at least one ordinary-row
  edge excludes exactly the no-outgoing-row nodes invariant 2 owns.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: [needed: an adversarial fixture holding both a dead-end
    node and a livelocked cycle, asserting each defect takes exactly one
    code]
  - **If wrong**: one defect takes two codes; surfaces as a double
    finding in the fixture corpus and a burying diagnostic for authors.
- **A5 No shipped consumer closes the blocking-code enum: the CLI help
  and envelopes enumerate via `graphlint.BlockingCodes()`
  (`internal/cli/lint.go`, `internal/cli/help_all_test.go`), and the
  manual doc surfaces (`docs/cli-reference.md`'s code table,
  `docs/model-authoring.md`'s inline code mentions) are updated in the
  same change.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: [needed: a sweep for consumers switching or validating
    over the blocking-code set beyond the dynamic enumerators above]
  - **If wrong**: an enum-closing consumer breaks on the append;
    surfaces as a failing test or an envelope rejecting the new code.
- **A6 The split closure for realistic models stays under
  `NodeCeiling()`: `models/rdr.toml`'s terminal-participating keys are
  finite and small, so the split product does not trip C3's size
  refusal.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: [needed: split-node count for `models/rdr.toml`
    measured by the spike, against the 4096 ceiling]
  - **If wrong**: the check refuses (`graph-product-too-large`) on the
    checked-in model; surfaces as a blocking refusal where a verdict was
    expected.
- **A7 The property class is the literature's: "from every reachable
  state some terminal is reachable" is CTL `AG EF terminal`, decidable
  by backward reachability from the terminal set, and distinct from
  fairness-dependent `AF` liveness (Clarke/Emerson/Sistla 1986;
  Newcombe 2015).**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: [needed: citation with section — the local corpora
    yielded no quotable source at Propose
    (`docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/research/prior-art.md`);
    demoted here rather than leaned on]
  - **If wrong**: the scope framing misnames the property class; the
    mechanism and fixtures are unaffected, but the out-of-scope line
    (`AF`) must be restated.

## Proposed Solution

### Approach

Add one new blocking finding class, `graph-terminal-unreachable`
(invariant 8: terminal reachability), computed in-house inside
`internal/graphlint` over the traversal `analysis.go::newAnalysis`
already runs — no second formal artifact, no new CLI surface.

The check answers the four axes as follows. **Quantifier**: run over
split nodes — every reachable merged node is split on the
terminal-participating keys exactly as invariant 2 splits
(`analysis.go::splitNode` over `analysis.go::terminalKeys` ⇒ the
universal half of the claim never reads a merged node, per RDR 0006's
soundness-direction rule); terminal satisfaction is universal per split
node (`analysis.go::satisfiesSomeTerminal`, reused verbatim); path
existence is existential over ordinary-row edges, the direction merged
reads are licensed for. A backward closure from the terminal-satisfying
split nodes marks good nodes; the residual miss on keys outside the
terminal-participating set is the same pinned false-negative invariant
2 already books (`analysis.go::hasOutgoingOrdinaryRow` doc, D12
pattern). **Granularity**: bad nodes are grouped into strongly
connected components; one finding accuses each *trapping* component — a
sink of the bad-region condensation containing at least one
ordinary-row edge — so the author sees the cycle, not a member census.
**Authority**: blocking — the advisory tier is closed at four
(`0006:C17`) and the blocking tier is a floor ("at least these blocking
invariant classes", `0006:C3` ⇒ a successor may append), and the
accusation is sound over the declared relation. **Refusal posture**:
size-scoped — skip when the forward traversal was already truncated
(the ceiling finding blocks in its place), refuse with
`graph-product-too-large` when the split closure itself would exceed
the published ceiling; never a partial verdict.

### Technical Design

The pass lives in `internal/graphlint` beside the other invariants and
consumes only the derived view `analysis` already carries: the merged
reachable node set (`reach.go::Reach` fixpoint), the split machinery,
and RDR 0003's evaluator via `atomAdmitsValue`. Data flow: (1) split
every reachable merged node on `terminalKeys()`; (2) mark the split
nodes satisfying some declared terminal (`satisfiesSomeTerminal`); (3)
build the split-node edge relation — for each split node and each
ordinary row satisfiable there (`matchSatisfiable`), take
`successor()`, re-anchor it into the merged fixpoint by presence
footprint and subsumption (`subsumes`), and re-split it on the
terminal-participating values the produced node holds; (4) backward
closure from the marked nodes gives the good set; (5) SCC-condense the
bad remainder (Tarjan, introduced by this RDR — no SCC utility exists
in `internal/`; searched, none exists) and accuse each trapping sink
component per C1. Escape rows are self-loops carrying no write and no
clear (`0006:D-reachability-relation`; writer: RDR 0002's normalizer) ⇒
they add no edge and no progress here, exactly as invariant 2 excludes
them.

#### Normative Contracts

**C1**

```normative
`graph-terminal-unreachable` is a new BLOCKING finding code (invariant
8, terminal reachability), appended to the blocking tier under
`0006:C3`'s "at least" floor; the advisory tier stays closed at four
(`0006:C17`). It reports: a reachable region of the declared owned-state
graph from which no declared terminal is reachable. The accused unit is
the TRAPPING COMPONENT: a strongly connected component of bad split
nodes (no path to any terminal-satisfying split node) that has no edge
to another bad component and contains at least one ordinary-row edge (a
self-loop counts). Every trapping component in a run is reported
(`0006:C16`); bad nodes upstream of a trap are not separately accused.
A split node with no outgoing ordinary row is `graph-dead-end`'s (or
invariant 7's) population and MUST NOT take this code.
```

**C2**

```normative
Quantifier and abstraction (shared doctrine homed at JDR 0001
§JD-23): the check runs over split nodes — each
reachable merged node split on the terminal-participating keys exactly
as invariant 2 (`splitNode` over `terminalKeys`). Terminal satisfaction
is UNIVERSAL per split node and reuses `satisfiesSomeTerminal`
verbatim; no second satisfaction semantics is introduced. Path
existence is EXISTENTIAL over ordinary-row edges: an edge exists from
split node `s` for each non-escape row whose match is satisfiable at
`s` (`matchSatisfiable`), and its target set is the split forms — over
the terminal-participating values actually held — of the merged
fixpoint node subsuming `successor(s, row)`. The imprecision of that
edge set is ONE-DIRECTIONAL: guard atoms are never consulted and owned
match atoms are decided existentially per held value
(`reach.go::matchSatisfiable` doc, `ownedAtomSatisfiable`), so a split
node's edge set is a superset of the edges any concrete view it stands
for takes — splitting removes no concrete edge, and extra edges only
widen `EF` (the miss direction, never a false accusation). The anchor
step is total and deterministic: when several fixpoint nodes subsume
the produced node, the first in the fixpoint's construction order is
taken (`reach.go::indexOf` precedent); when none does, the produced
node's own split forms are the targets. The closure's node universe is
the precomputed finite set of split forms of fixpoint nodes plus any
fallback forms, each drawn from declared domains, so the closure
terminates without a widening argument. Escape rows contribute no
edge. The guarantee is therefore: a model whose declared relation
reaches a terminal from every reachable split node is never accused.
The misses — livelock separable only on merged
non-terminal-participating keys, or masked by an unpruned
guard/observed edge — are ACCEPTED BY THIS RDR IN ITS OWN RIGHT as a
path-scoped residual, strictly larger than invariant 2's local
one-step residual: one over-approximated edge anywhere on a path can
clear every node upstream of it. `AF`/fairness properties are out of
scope.
```

**C3**

```normative
Refusal posture (size-scoped): when the forward traversal is incomplete
(`reach.go::reach` returned `complete == false`; `checkNodeCeiling`
already emits blocking `graph-product-too-large`, element `traversal`),
this check MUST NOT run and MUST NOT emit — the ceiling finding blocks
in its place, so the silence is never a green. When the traversal is
complete but this check's split closure would exceed `NodeCeiling()`
split nodes, the check MUST emit one blocking `graph-product-too-large`
finding whose element is `traversal` and whose message names the
terminal-reachability split and the authored remedy (narrow a declared
domain of a terminal-participating key), and no
`graph-terminal-unreachable` finding — a refusal, never a partial
verdict. Decision tables are
silent by class (`0010:C5`'s machine-only rule, the same
`table.IsDecisionTable` key `checkDeadEnd` uses).
```

**C4**

```normative
Finding shape and surfaces: a `graph-terminal-unreachable` finding
carries code, model identity, severity `blocking`, element, and message
per `0006:C13`; it names no rule. Its element is the `nodeElement`
rendering of the trapping component's lexicographically least member
split node; the message states that the claim is over the DECLARED
graph (the relation over-approximates the runtime), names the
component's member count, an ordered, bounded enumeration of member
states sorted by node key, and — when the trap is entered from outside
— the rule id of one entering row, so the author has an authored site
to edit. Ordering and identity follow `0006:C15` unchanged
(element-namespace id): findings are sorted by finding identity
(`engine.go::sortFindings`) and every enumerated member list is sorted
by canonical node key, so the SCC algorithm's internal visit order is
unobservable. The code is
appended to `taxonomy.go::blockingCodes` after `CodeTerminalEscape` in
taxonomy order, which publishes it through the existing help output and
text/JSON envelopes with no new CLI surface; the aggregate failure code
stays `graph-lint-failed` (`0006:C14`). Per `0017:C1` the code is the
finding's own disposition drawn from graphlint's producer-declared
vocabulary; `docs/cli-reference.md`'s finding-code table gains its row,
and `docs/model-authoring.md`'s inline finding-code prose is extended
where it enumerates the blocking classes, in the same change.
```

#### Load-Bearing Decisions

- **Identity** — two findings are the same iff they agree on
  `(model id, graph-terminal-unreachable, element)`, where the element
  is the lexicographically least member split node's `nodeElement`
  rendering — the component representative, stable across runs because
  node keys are canonical and injective (`reach.go::key`).
- **Naming** — `graph-terminal-unreachable`: names the violated
  guarantee (a terminal is unreachable), parallel to
  `graph-unreachable-rule`'s direction. Rejected: `graph-livelock`
  (names the symptom, not the violated declaration) and
  `graph-no-terminal-path` (reads as a per-path claim; the check is
  per-region).
- **Selection / predicate** — when several trapping components exist,
  all are reported in one pass (`0006:C16` ⇒ complete emission, never
  first-failure); bad nodes outside any trapping sink component are
  never accused — the cure is the trap, not a census (the same
  one-defect economy `checkAlwaysPresentOwned` and invariant 7 apply).

#### Illustrative Code

Illustrative — shape only, not load-bearing:

```
good := splitNodes satisfying some terminal        // universal, reused
work := good
for s := pop(work):
    for p in splitPredecessors(s):                 // ordinary rows only
        if p not in good: good += p; work += p
bad := allSplitNodes - good
for scc in tarjan(bad, ordinaryEdges within bad):
    if scc is sink && scc has an ordinary edge:
        emit(graph-terminal-unreachable, element=leastMember(scc))
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Declared root and terminal predicates | Predecessor (0006:A6 → RDR 0002 schema) | Available | Root and stop set exist; missing root is already blocking (`0006:C18`) |
| Merged reachability fixpoint | Existing (`reach.go::reach`) | Available | Reused as-is; ceiling semantics unchanged |
| Split machinery and terminal satisfaction | Existing (`analysis.go::splitNode`, `terminalKeys`, `satisfiesSomeTerminal`) | Available | Reused verbatim — no second semantics |
| SCC condensation over split nodes | This RDR | Introduced | Internal helper in `graphlint`; not exported |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Edge relation | `reach.go::matchSatisfiable` + `successor` | Over-approximates (guards, observed atoms unpruned) | Reuse | Miss direction pinned in C2 |
| Terminal test | `analysis.go::satisfiesSomeTerminal` | Opaque carrier on unbounded keys (A3) | Reuse | Verbatim reuse; A3 spikes the opaque arm |
| One-step liveness | `analysis.go::checkDeadEnd` / `hasOutgoingOrdinaryRow` | One-step only | Extend (new sibling pass) | Population disjointness per C1 |
| Finding taxonomy | `taxonomy.go::blockingCodes` | Blocking tier is a floor (`0006:C3`) | Extend (append) | New code row in help/docs; enum append precedent `0010`'s reason append |
| Finding order/identity | `engine.go::sortFindings` (`0006:C15`) | Element-namespace ids only | Reuse | Element = component representative |

### Decision Rationale

Scored matrix (approaches × deciding criteria; one clause per cell):

| Criterion | A: merged-EF, per node | B: split-EF + SCC (chosen) | C: refined forward fixpoint | D: external checker |
| --- | --- | --- | --- | --- |
| Correctness fit | false-green on the named axis-1 direction | universal half exact on terminal keys; misses pinned like invariant 2 | precision ≈ B | strongest, plus `AF` |
| Prior-art alignment | contradicts `0006:D-soundness` (universal over merged) | exactly invariant 2's shipped remedy | second relation beside `reach` | rejected by `0006:ALT4` |
| Reversibility | easy | easy — additive pass | hard — re-sizes ceiling semantics | hard — new toolchain |
| Blast radius | one pass | one pass + one enum append | traversal core touched | build/CI + model translation |
| Cost | lowest | low — reuses split machinery | medium | highest, drifting duplicate |

The deciding rows are correctness fit and prior-art alignment: the
problem statement names merged-EF the false-green direction, and RDR
0006's soundness rule already prescribes the split remedy for the
universal half of exactly this shape — invariant 2 is the shipped
precedent, so B introduces no new abstraction doctrine. C buys no
precision over B but duplicates the traversal (the in-house edition of
the drifting-artifact failure the review rejected); B is distinct from
C on exactly that axis — its closure is a VIEW over the one fixpoint's
node universe and successor semantics, minting no second join rule and
no second ceiling population, where C runs an independent traversal
with its own. D re-litigates `0006:ALT4` with the same cons.
Granularity: SCC-collapse because a trapping component's members share
one cure (Selection decision above). Authority: blocking, on the
finding's nature — the accusation direction is sound over the declared
relation, the same posture `graph-dead-end` already takes on the same
relation — with the closed advisory tier (`0006:C17`) as the secondary
constraint.

Premortem: hardened (hardened) — critic PASS with mitigations, none
forcing a switch
(`docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/propose-premortem/critic.md`):
P-2/P-3 (edge-direction under the split) answered from shipped
existential match semantics and folded into C2; P-5/P-6 (anchor
tie-break, closure termination) folded into C2 and A1; P-8/P-21
(entry-row diagnostic, SCC-order determinism) folded into C4; P-12
(orphan split-node population) folded into A4; P-13 (skip window)
folded into Failure Modes; P-14 cleared at source (`0006:C3`/`C14`
carry no uncited stability half; findings are append-only); P-17/P-18
folded as the Alt-1 descope guard and the view-not-second-relation
clause; P-4/P-9 folded as C2's own written residual acceptance;
P-3/P-16's narrative case added to the fixture plan.

Ground-sweep: clean (34 anchors)
(`docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/propose-premortem/ground-sweep.md`;
33 CONFIRMED, one cosmetic doc-surface miss — `docs/model-authoring.md`
carries inline finding-code prose, not a table — fixed inline in C4 and
A5).

Joint-check: fired → 0015 (home: JDR 0001 §JD-23). Context: RDR 0015 (Draft,
unproposed) owns the ranking of D12's two record-supported readings of
invariant 2's merged-node quantifier and anticipates "a successor
clause naming the quantifier explicitly"; this RDR's C1 carves its
population against `graph-dead-end`'s and its C2 codifies the
terminal-participating-key split bound with the D12-class residual
accepted for invariant 8 — the same doctrine 0015 may relax
(REQ-37→REQ-111), replace (correlation tracking), or confirm. The
joint question: one quantifier doctrine must govern both liveness-
shaped invariants, and its normative home is the umbrella clause both
records cite: JDR 0001 §JD-23. Non-fires adjudicated as context: 0021 (reads the same
relation for export; both records declare mutual independence),
0012/0019 (read `atomAdmitsValue` / own `checkDanglingEdge` root
semantics — adjacent surfaces, no shared decision), 0014 (cites
`graph-lint-failed`/`AggregateCode` as an existing oracle literal),
0024 (declares no graphlint change). 11 other open peers grep-clean on
every distinctive anchor and on `graph-terminal-unreachable`. Absence
arm: no `Final` peers exist (0001–0011 are Implemented), and this RDR
adds a refusal rather than converting one to an acceptance.

Joint-check: fired → 0029 (home: `cli/0029 §Normative Contracts` C3).
Recorded 2026-09-11 from 0029's propose-stage check; symmetric write.
Context: RDR 0029 (Draft, `mid`) settles what a released version promises
an agent about machine-readable output, and its C3 rules that a new lint
finding code enters at severity `info` and is promoted to `blocking` only
in a later, disclosed release. This RDR's C1 introduces
`graph-terminal-unreachable` directly as a BLOCKING code and restates
"the advisory tier stays closed at four (`0006:C17`)" — the exact rule
0029:C3 changes. Disposition: 0029:C3 is the single normative home; this
record CITES it rather than restating the C17 closure, and at
implementation either enters the code at `info` or claims C3's
re-attribution clause (a code may enter at `blocking` when it
re-attributes a condition some other code already refused). Which of the
two applies is 0022's call, made against 0029:C3 — not re-decided here.
Note both records are pre-1.0.0, where `0.x` carries no compatibility
guarantee, so the rule binds as declared intent until 1.0.0.

Joint-check: fired → 0030 (home: `cli/0030 §Normative Contracts` C3).
Context: A1 rests on `internal/graphlint/reach.go::successor` applying
literal writes; 0030 (Draft) admits a computed write form and its C3 keeps
every row's successor a literal by load-time expansion — cite it there
rather than assume it. Recorded symmetrically by 0030's propose,
2026-09-19.

## Alternatives Considered

### Alternative 1: Merged-node backward reachability (pinned miss)

**Description**: The review's literal suggestion — compute the terminal
set existentially over merged nodes and run one backward closure over
the merged relation; accept every merging-induced miss as pinned (the
D12 pattern).

**Pros**:

- Cheapest: no split space, no SCC pass needed for correctness.
- Uses only relations already materialized.

**Cons**:

- The universal half (`AG` — *every* reachable state) is evaluated
  against widened nodes, which `0006:D-soundness-direction` forbids for
  universal checks: a merged node whose value set mixes a terminating
  and a livelocked value shows a terminal path only some of its
  concrete views have.
- Incoherent with invariant 2: the same node could pass this check yet
  fail the one-step dead-end test's split form, or vice versa.

**Reason for rejection**: it is precisely the false-green direction the
problem statement names for this check, on the property the RDR exists
to catch; the pinned-miss budget belongs on the keys the split cannot
reach, not on the terminal keys the shipped remedy already covers.
Guard: if implementation descopes the split, the shipped behavior
becomes this alternative — in that event the severity question reopens
(a merged-`EF` pass misses more and may not carry the blocking
rationale unchanged); descoping is a route back to this fork, not a
silent simplification.

### Alternative 2: Refined forward fixpoint (never merge terminal keys)

**Description**: Run a second `reach`-style traversal whose join never
widens terminal-participating keys, then a plain backward closure over
its exact-on-terminal-keys node set.

**Pros**:

- Same precision class as the chosen approach, with a simpler closure.

**Cons**:

- A second reachability relation beside `reach.go::reach` — the
  in-house edition of the "second, drifting formal artifact" the review
  rejected; two fixpoints to keep in agreement.
- Re-sizes the node-ceiling semantics: the primary published ceiling
  would gate a different node population depending on the pass.

**Reason for rejection**: buys no precision over splitting the existing
fixpoint's nodes, at the cost of a duplicate traversal that can drift.

### Alternative 3: Per-node findings, no SCC collapse

**Description**: Emit one `graph-terminal-unreachable` finding per bad
split node, like invariant 2 emits per node.

**Pros**:

- No SCC machinery.

**Cons**:

- A livelocked cycle of *n* states with *m* trapped predecessors emits
  `n + m` findings for one authored defect, burying the cycle.

**Reason for rejection**: the actionable unit is the trap — one finding
naming the cycle is the granularity the problem statement's axis 2
asks for, and the one-defect economy invariant 7 and
`checkAlwaysPresentOwned` already practice.

### Briefly Rejected

- **External model checker (TLC/NuSMV)**: `0006:ALT4`'s rejection
  stands — a second graph language whose duplicated semantics drifts;
  native lint over the model intrastate actually consumes.
- **Advisory severity**: on the finding's nature first — a livelock
  breaks the same reach-a-terminal guarantee the blocking
  `graph-dead-end` enforces one step at a time, with the same declared-
  relation soundness, so the confidence that justifies blocking there
  justifies it here; the closed advisory tier (`0006:C17`) is the
  secondary constraint, not the argument.
- **Path-sensitive enumeration**: unrepresentable in this engine by
  design — the fixpoint is keyed on owned-state, not path
  (`reach.go::reach` doc), and the exponential reading is already
  rejected by `0006:D-join-rule-and-termination`.
- **`AF`/fairness liveness**: out of scope by the problem statement —
  outcome choice is caller-driven, so eventuality under outcome
  sequences is not a model property here.

## Context

### Background

Tracked as kata `intrastate#4hps` (`release:post-1.0`), from the
state-machines research review (the property class per
Clarke/Emerson/Sistla 1986 and Newcombe 2015; verification belongs in
the validator, not a runtime). Not a facet of kata `jjkh` (RDR 0021,
the graph export): this invariant runs over the in-memory relation
and neither work item blocks the other. Fixture obligations carried
from the seed: a model with a legal cycle that *does* reach a terminal
must pass, a model with a cycle that cannot must fail, and
`models/rdr.toml` must keep passing unchanged.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/graphlint/reach.go` (reachability graph),
`internal/graphlint/taxonomy.go` (finding codes),
`internal/graphlint/analysis.go`. Governing record: RDR 0006 (A6 —
declared initial/terminal give reachability a root and stop set;
finding taxonomy and soundness-direction rule).

## Research Findings

### Investigation

Read at source: `internal/graphlint/reach.go` (the merged fixpoint, its
join and subsumption, the opaque carrier `heldValues`),
`analysis.go` (invariant 2's split remedy, `hasOutgoingOrdinaryRow`'s
pinned-residual doc, invariant 7's one-defect emission,
`checkNodeCeiling`'s refusal), `taxonomy.go` and `engine.go` (tier
vocabulary, `severityFor`, `sortFindings` identity), `models/rdr.toml`
(cycles present; terminals `landed`/`abandoned`); RDR 0006 via the
projector (`0006:C3` floor, `0006:C12`–`C17`, `0006:ALT4`,
`D-reachability-relation`, `D-soundness-direction`), RDR 0010's
append precedent (`ReasonNoParticipatingDimension`), and PROPOSED
peers `0017:C1` (per-finding code identity) and `0024:§approach`
(emit checks land in the load pipeline; emit values stay uninterpreted
⇒ no interaction with the owned-state relation this pass traverses).
External pass: ⚠ no prior-art coverage for this problem class in the
indexed corpora (StateMachineRes/StateMachineLit/PapersFast — queries
and rejected branches in
`docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/research/prior-art.md`);
the literature framing is demoted to A7 rather than leaned on.

### Key Discoveries

- **Documented** — `0006:C3` makes the blocking tier a floor ("at
  least these blocking invariant classes") ⇒ a successor RDR may append
  a code without amending 0006; RDR 0010 already appended to a closed
  append-only vocabulary (`taxonomy.go::ReasonNoParticipatingDimension`).
- **Documented** — `0006:D-soundness-direction` classifies checks by
  quantifier: universal claims must not read merged nodes; invariant
  2's shipped remedy is the terminal-key split
  (`analysis.go::splitNode`) ⇒ the `AG` half of this check must split,
  and the machinery to do it already exists.
- **Documented** — the residual false negative on
  non-terminal-participating keys is deliberately pinned
  (`analysis.go::hasOutgoingOrdinaryRow` doc, "See D12") ⇒ this check
  inherits, not widens, that budget.
- **Documented** — escape rows normalize to write-free, clear-free
  self-loops (`0006:D-reachability-relation`; producer: RDR 0002's
  normalizer) ⇒ they are never progress and add no edge here.
- **Documented** — help output and envelope enumerate codes dynamically
  (`internal/cli/lint.go` over `graphlint.BlockingCodes()`) ⇒ the
  append publishes itself; only the two doc tables are manual.
- **Assumed** — re-anchoring split successors by subsumption is total
  (A1); `models/rdr.toml` stays green (A2); opaque terminal keys do not
  mass-accuse (A3); the closure stays under the ceiling (A6).

## Trade-offs

### Consequences

- Positive: livelock — the one property class only an external checker
  would otherwise add — becomes a native blocking invariant; legal
  cycles that can terminate stay legal.
- Positive: no new abstraction doctrine — quantifier, split, escape
  handling, refusal, ordering, and emission economy are all reused from
  the shipped invariants.
- Negative: a new blocking code can newly reject models that lint green
  today (that is its purpose); the fixture obligations bound the risk
  on the checked-in model (A2).
- Negative: misses are structural, not incidental — livelock reachable
  only through correlations on non-terminal-participating keys, or
  masked by unpruned guard/observed edges, is pinned out of reach (C2).
- Negative: the split closure adds work proportional to reachable
  nodes × terminal-key split product; bounded by the existing ceiling
  (C3), never unbounded.

### Risks and Mitigations

- **Risk**: the re-anchor step (A1) hides a footprint the fixpoint
  never minted (e.g. a clear on a split node), leaving a successor with
  no host.
  **Mitigation**: A1's spike drives `successor` from split nodes of
  cyclic fixtures before the pass is built; if unhosted successors
  exist, the target set falls back to the produced node's own split
  forms — still finite, since footprints and per-key values come from
  declared domains.
- **Risk**: element identity churns across runs if the least-member
  representative changes when the merged fixpoint widens a value set
  (a node's key changes as its value sets grow).
  **Mitigation**: the fixpoint is run to completion before this pass
  reads it, and node keys are canonical and injective
  (`reach.go::key`); the determinism scenario in the fixture corpus
  asserts identical findings across repeated runs.
- **Risk**: opaque terminal keys (A3) empty the terminal set and accuse
  every cyclic model that terminates on an unbounded key.
  **Mitigation**: A3's spike decides the evaluator's opaque verdict at
  Resolve; the axis-4 downgrade (refusal instead of accusation) is the
  named fallback.

### Failure Modes

- **Breaks visibly**: a livelocked model takes one blocking
  `graph-terminal-unreachable` per trap inside the aggregate
  `graph-lint-failed` envelope; the element names the trap's
  representative state and the message lists members.
- **Fails silently**: pinned misses (C2) — merged non-terminal keys,
  guard/observed edges taken as traversable. Diagnosis path: the
  model's own runtime refusals; the pinned budget is documented in the
  check's doc comment like invariant 2's.
- **Skip window**: when the traversal ceiling trips, this check is
  skipped behind a blocking `graph-product-too-large` — the model
  cannot pass lint, so the skip is never a green; a separate
  machine-readable "liveness not evaluated" marker is declined (the
  envelope stays `0006:C14`'s shape), and a team that ships past a red
  gate has left the tool's guarantee surface — the ceiling finding's
  remedy text is the recovery path.
- **False accusation relative to runtime** (not to the declared
  relation): a trap reachable only via an infeasible abstract path —
  the same accepted posture `graph-dead-end` already carries, with the
  same remedy (narrow the declared domains or the match patterns).
- **Recovery**: the finding is advisory-free and blocking-only; a false
  accusation is fixed in the model or falsifies A2/A3 and reopens the
  quantifier choice — never a suppression flag.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1/A3/A6 spikes land at
  Resolve; A2/A4 are MVV scenarios)

### Minimum Viable Validation

1. Author a livelock fixture: two states cycling on ordinary rows with
   no row leaving the cycle toward the declared terminal — run
   `intrastate lint` — expect failure with exactly one
   `graph-terminal-unreachable` finding whose element names the cycle's
   representative state and whose message lists both members.
2. Author a legal-cycle fixture: the same cycle plus one exit row
   reaching a terminal-satisfying state — run lint — expect success
   with no `graph-terminal-unreachable` finding.
3. Run lint over `models/rdr.toml` — expect the findings unchanged from
   today's run (green, revise/rewind cycles pass).
4. Author a mixed fixture holding one dead-end node and one livelocked
   cycle — expect exactly one `graph-dead-end` and one
   `graph-terminal-unreachable`, no double accusation (A4).
5. Run lint twice on the livelock fixture (step 1) — expect
   byte-identical findings output across runs: this is the proof
   obligation behind C2's deterministic anchor selection and C4's
   canonical-node-key member ordering (no SCC visit order or map
   iteration order may reach the output).

### Phase 1: Code Implementation

#### Step 1: Taxonomy and documentation append

Append `CodeTerminalUnreachable` to `taxonomy.go` and
`blockingCodes` (after `CodeTerminalEscape`); add the code's row to
`docs/cli-reference.md` and `docs/model-authoring.md`.

#### Step 2: The liveness pass

Implement the split-closure pass in `analysis.go` (C1–C4): split,
mark, edge relation with re-anchor, backward closure, SCC condensation,
trapping-component accusation, decision-table silence, and both refusal
arms.

#### Step 3: Fixture corpus and adversarial tests

The MVV fixtures plus: a merged node mixing terminating and livelocked
values on a terminal key (the axis-1 false-green witness — must be
accused); a node rescued only by an escape self-loop (must not count as
progress); a cycle whose only exit row is guarded on a
non-terminal-participating key (the guard is unpruned, so the edge is
traversable — must NOT be accused; the pinned-miss direction made
concrete); a cyclic model declaring terminals no state satisfies (every
trap reported); a decision table (silent); repeated-run determinism.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Conditional — omit (don't N/A-bullet) this section unless
comparing alternatives on empirical performance grounds.
Do not include effort estimates or speculative
throughput targets. Rough performance metrics are
appropriate only when comparing alternatives — note
empirical data or obvious gains that support the
chosen approach over a rejected one.]

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0006 (`0006:C3`, `0006:C12`–`C18`, `0006:ALT4`,
  `0006:D-reachability-relation`,
  `0006:D-soundness-direction-is-per-invariant-not-global`) — the
  governing taxonomy, relation, and quantifier doctrine.
- RDR 0010 (`0010:C5`) — decision-table silence and the append-only
  vocabulary precedent; RDR 0017 (`0017:C1`) — per-finding code
  identity; RDR 0024 (`0024:§approach`) — emit checks stay in the load
  pipeline.
- Source reviewed: `internal/graphlint/{reach,analysis,taxonomy,engine}.go`,
  `internal/cli/lint.go`, `models/rdr.toml`.
- Kata `intrastate#4hps`; state-machines research review
  (Clarke/Emerson/Sistla 1986, "Automatic Verification of Finite-State
  Concurrent Systems"; Newcombe et al. 2015, "How Amazon Web Services
  Uses Formal Methods") — pending A7.
- Stage-2 search record:
  `docs/rdr/0022-terminal-reachability-liveness-invariant/evidence/research/prior-art.md`.
