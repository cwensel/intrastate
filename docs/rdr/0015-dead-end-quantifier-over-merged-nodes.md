# Recommendation 0015: Dead-end quantifier over merged nodes

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
- **Type**: Technical Debt
- **Profile**: mid — provisional: one contract (invariant 2's
  quantifier over merged nodes), a successor clause to locked
  RDR 0006's guarantee text.
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
- **Priority**: Medium
- **Related Issues**: kata `intrastate#pz9z` (1538); roborev job
  6155; RDR 0006 deviation D12
- **Predecessors**: 0006-graph-lint-authority-and-guarantees
- **Seam Lineage**: no prior accretion

## Problem Statement

A model author running graph lint needs to know what invariant 2's
dead-end guarantee actually promises: exactness over the whole
owned-state view, or exactness only over the terminal-participating
projection with a declared, bounded miss outside it. Today a reachable
node that is dead only on a key participating in no terminal is
silently unaccused — the outgoing-row test runs existentially over the
merged node, so an exit serving one of its values rescues the whole
node — and the record supports both readings: REQ-37's explicit MUST
bounds the *split* to terminal-participating keys ("never the whole
lattice"), while REQ-111 states invariant 2's exactness in terms a
wider quantifier can be inferred from. Deviation D12 says it directly:
both readings have record support and the record does not rank them.

The decision, made once, with measured costs weighed rather than
inherited: (1) **relax REQ-37** toward REQ-111 exactness — widening
the split decorrelates keys the model actually correlates, minting 9
false `graph-dead-end` findings on `models/rdr.toml` and breaking
REQ-122's zero census unless the model is reworked and SC-23's census
re-run (defensible under REQ-117: "the cure is a clearer model, never
a weaker lint"); (2) **correlation tracking** — rejected by the record
as exponential for the general traversal, but not costed at the
narrower scope of a per-node correlation record restricted to
terminal-adjacent keys; or (3) **confirm the bound** and codify the
accepted false negative as a normative limit of the guarantee.
Whichever arm wins, RDR 0006 needs a successor clause naming the
quantifier explicitly — a soundness-versus-precision,
author-values decision the evidence narrows but does not close.

## Critical Assumptions

- **A1 Widening the outgoing-row quantifier beyond
  terminal-participating keys still mints at least one false
  `graph-dead-end` finding on `models/rdr.toml` at HEAD** — D12's
  nine-finding Phase 3c measurement holds in direction (the exact
  count is not load-bearing; the sign is).
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: to produce at Resolve — re-run the D12 widening
    experiment (range `hasOutgoingOrdinaryRow` over all multi-valued
    keys in a throwaway branch) and record the finding count against
    `models/rdr.toml`.
  - **If wrong**: Alternative 1's measured cost evaporates and the
    ranking must re-run — the relax arm re-enters at Resolve before
    any carrier surface is edited.
- **A2 The REQ-122 census is zero at HEAD**: `lint --model
  models/rdr.toml --as=json` returns zero findings, so the residual
  class has no member on the conforming model.
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: to produce at Resolve — the lint invocation and its
    empty findings array.
  - **If wrong**: the residual class (or another defect) is already
    populated on a maintained model — C1's re-rank trigger fires
    before the clause lands and the choice reopens.
- **A3 Invariant 8 adopts the same split bound and accepts the
  path-scoped D12-class residual, homed at JDR 0001 §JD-23**, so
  confirming here binds both liveness invariants with no 0022 rework.
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: cli/0022:C2 (quantifier and abstraction clause,
    "shared doctrine homed at JDR 0001 §JD-23").
  - **If wrong**: the two records diverge on the one doctrine and the
    joint decision reopens at §JD-23 — a Stage 7.1 contradiction, not
    a silent drift.
- **A4 REQ-111 is the only record support for the wider quantifier**
  — no other RDR 0006 clause states invariant 2's exactness in
  concrete-lattice terms, so C1's subordination of REQ-111 leaves no
  live contradiction elsewhere in the governing record.
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to produce at Resolve — sweep
    `docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md`
    and the 0006 record for exactness claims naming invariant 2.
  - **If wrong**: the successor clause contradicts an unsubordinated
    clause — a 7.1 SPEC-DEFECT against this RDR's fence.
- **A5 The pin test's positive arm still exercises invariant 2** —
  `TestAdvDeadEndExistentialOnMergedNode` asserts the bounded split
  catches a dead half on a terminal-participating key, so
  re-documenting it as the declared-limit witness removes no coverage.
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**:
    `internal/graphlint/adversarial_0006_test.go::TestAdvDeadEndExistentialOnMergedNode`
    (D12 "Fixture corrected": "adds a positive arm asserting the
    bounded split still catches a dead half living on a
    terminal-participating key").
  - **If wrong**: the tripwire pins only an absence and an invariant 2
    regression could ship green — the positive arm must be added, not
    merely re-documented.

## Proposed Solution

### Approach

Confirm the bound — arm (3) — and rank the residual, as the answer to
JDR 0001 §JD-23's charter. Invariant 2's quantifier stays exactly what
ships: the dead-end check runs over split nodes bounded by the
terminal-participating keys (REQ-37's MUST), universal terminal
satisfaction per split node, existential outgoing-row test over the
keys that stay merged. No lint behavior changes. What changes is the
record and its carrier surfaces: the guarantee gains an explicit
successor clause naming the quantifier and its declared residual (C1);
the residual class gets a canonical name, classification, and a
normative re-rank trigger; the pin test is re-documented from open
question to declared-limit witness, D12 closes, and the authoring docs
state the limit with the REQ-117 cure at the point the guarantee is
described (C2).

The ranking — the charter's deliverable — is confirm > scoped
correlation tracking > relax, on measured evidence. Relaxing REQ-37
toward REQ-111 exactness converts a declared one-directional miss into
measured false accusations: nine false `graph-dead-end` findings on the
conforming `models/rdr.toml`, breaking REQ-122's zero census — the
false-accuse direction the record's soundness doctrine forbids while it
tolerates misses. Scoped correlation tracking closes the miss soundly
but is priced (Alternative 2) against a residual population with zero
observed members outside the adversarial pin fixture. Confirming costs
nothing, changes no behavior, and binds both liveness-shaped invariants
through the §JD-23 anchor with cli/0022:C2 standing unchanged.

### Technical Design

This is record-and-carrier surgery, not a code change. The shipped
quantifier is already the confirmed reading:
`internal/graphlint/analysis.go::checkDeadEnd` splits each reachable
merged node via `analysis.go::splitNode` over `analysis.go::terminalKeys`
⇒ the universal terminal test (`analysis.go::satisfiesSomeTerminal`)
never reads a merged value on a deciding key; the outgoing-row test
(`analysis.go::hasOutgoingOrdinaryRow`) is existential via
`matchSatisfiable` over the keys left merged inside the split node ⇒
an exit serving one held value of a non-terminal key rescues the split
node — the residual this RDR declares rather than closes. The merged
node set the check reads is written by the reachability fixpoint
(`reach.go::reach`), whose join decorrelates keys by construction —
that writer is why the residual exists at all.

The residual class is classified once, for both liveness invariants:

- **Population**: reachable split nodes whose deadness is separable
  only on the values of a merged key participating in no terminal —
  the population cli/0022:C1 carves against `graph-dead-end`'s.
- **Direction**: one-directional — a miss (false negative), never a
  false accusation; the same soundness direction 0006's doctrine
  prescribes per invariant.
- **Extent**: local for invariant 2 (one node's one-step outgoing-row
  test) versus path-scoped for invariant 8 (cli/0022:C2's strictly
  larger closure residual) — the rank is per-invariant extent under
  one shared bound, cited from §JD-23, not restated.
- **Cure** (REQ-117 direction): model clarity — declare the correlated
  key in a terminal predicate, or write/clear it explicitly so the
  correlation is authored rather than inferred.

Three shipped surfaces carry the clause (C2): the call-site doc
comments in `analysis.go`, the pin test
`TestAdvDeadEndExistentialOnMergedNode` (the widening tripwire — any
widening of the split fails it and forces SC-23's census re-run before
acceptance), and `docs/model-authoring.md`'s dead-end passage. D12
(`docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/deviations.md`,
`Status: needs author decision`) closes by an appended resolution line;
its entry — including the trigger language — is not rewritten.

#### Normative Contracts

One independent load-bearing contract: the successor clause (C1). C2
carries C1's declaration onto shipped surfaces and is not separately
holdable.

**C1**

```normative
Invariant 2's quantifier is CONFIRMED as REQ-37's bound, adopted as
this record's answer to the shared merged-node liveness doctrine homed
at JDR 0001 §JD-23: the dead-end check runs over split nodes — each
reachable merged node split on the terminal-participating keys only,
never the whole lattice; terminal satisfaction is UNIVERSAL per split
node; the outgoing-row test is EXISTENTIAL over the keys that stay
merged inside the split node. The guarantee this clause succeeds
therefore reads: `graph-dead-end` is exact over the
terminal-participating projection; outside it the guarantee carries a
DECLARED RESIDUAL — the D12 residual: a reachable split node whose
deadness is separable only on the values of a merged
non-terminal-participating key is not accused. The residual is
normative, one-directional (a miss, never a false accusation), and
local to the one-step outgoing-row test. REQ-111's exactness statement
is SUBORDINATED to this clause: "recovering exactness" reads as
exactness of the universal terminal half over the split view, not
concrete-lattice exactness of the whole invariant. Re-rank trigger:
the first member of the D12 residual on a maintained (non-fixture)
model — a shipped dead state separable only on non-terminal keys, or a
nonzero SC-23-class census attributable to the residual — REOPENS this
ranking through §JD-23, with scoped correlation tracking (Alternative
2) as the named successor arm; until then neither liveness record
widens or narrows the split.
```

**C2**

```normative
The declared residual is carried on three shipped surfaces, with no
lint behavior change and no new I/O surface — `lint`'s findings on a
conforming model are unchanged by this RDR (MVV step 2 asserts the
zero census): (a)
`analysis.go::hasOutgoingOrdinaryRow`'s and `::checkDeadEnd`'s doc
comments cite cli/0015:C1 and JDR 0001 §JD-23 as the ranking's settled
home, replacing the open-question "See D12" framing; (b)
`TestAdvDeadEndExistentialOnMergedNode` remains the executable witness
of the residual and the widening tripwire — its doc comment cites
cli/0015:C1, and any widening of the split MUST fail this pin and
re-run SC-23's census against `models/rdr.toml` before acceptance; (c)
`docs/model-authoring.md`'s dead-end passage states the limit and its
cure where the guarantee is described: a state kept live only by
values of a key no terminal reads is outside the guarantee, and the
cure is declaring that key in a terminal predicate or writing/clearing
it explicitly (REQ-117 direction). D12 closes as resolved by cli/0015
via an appended resolution line; the entry body is not rewritten.
```

#### Load-Bearing Decisions

- **Naming** — the residual class's canonical name is the **D12
  residual** (the population cli/0022:C1 and :C2 call "D12-class"):
  reachable split nodes whose deadness is separable only on merged
  non-terminal-participating keys. Rejected: "existential-rescue
  defect" (implies a bug where the clause declares a limit) and
  "false-negative class" (names the direction without the
  population).
- **Selection / predicate** — of the two record-supported readings,
  REQ-37's explicit MUST wins over the quantifier inferred from
  REQ-111's principle. The deciding rule: an explicit MUST outranks an
  inference, and adopting the inference measurably breaks a second
  MUST (REQ-122's zero census — nine false findings, D12's Phase 3c
  measurement). REQ-111 is not discarded; its exactness claim is
  scoped to the universal half the split already makes exact (C1's
  subordination).

#### Illustrative Code

Illustrative — shape only, not load-bearing. The D12 residual in the
smallest form (the ADV-3 fixture's shape):

```
terminal:  status = "done"            # status is the only terminal key
node:      {status: [open], phase: [p, q]}   # reachable, merged on phase
rows:      advance: match phase = "p" → write status = "done"
split on terminalKeys() = [status]:   # phase stays merged
  {status=open, phase=[p,q]}          # one split node
outgoing?  advance satisfiable (SOME held value: phase=p) → rescued
concrete:  {status=open, phase=q} has no exit → the declared miss
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Dead-end quantifier | `analysis.go::checkDeadEnd` + `hasOutgoingOrdinaryRow` | D12 residual (the declared miss) | Reuse — confirm; doc-cite update only | C1, C2(a) |
| Widening tripwire | `adversarial_0006_test.go::TestAdvDeadEndExistentialOnMergedNode` | pins the miss; positive arm on a terminal key (A5) | Reuse — re-document as declared-limit witness | C2(b) |
| Guarantee prose | `docs/model-authoring.md` dead-end passage | states the guarantee without its limit | Extend — one limit-and-cure passage | C2(c) |
| Deviation ledger | 0006 `artifacts/deviations.md` D12 | `Status: needs author decision` | Extend — append resolution line | C2 |

### Decision Rationale

Four factors decide, and the measured ones dominate. (1) **Soundness
direction**: 0006's doctrine already ranks the two error directions —
a universal check's miss is tolerable, a false accusation is not — and
relaxing REQ-37 lands nine false accusations on the conforming model
(D12's measured experiment), while confirming books only the declared
miss. (2) **Explicit MUST over inference**: REQ-37 states the bound as
a MUST with "never the whole lattice"; the wider quantifier is an
inference from REQ-111's principle, and honouring the inference breaks
REQ-122's MUST — the Selection decision above. (3) **Measured cost
against measured benefit**: Alternative 2 closes the miss soundly but
its priced costs (a second, tuple-aware semantics beside the split
view; a fixpoint join and ceiling re-sized by pack products) buy
precision on a population with zero observed members outside the
adversarial fixture (A2). (4) **Doctrine unity**: §JD-23 requires one
quantifier doctrine over both liveness invariants; confirming binds
both through the anchor with cli/0022:C2 unchanged, while either other
arm forces 0022's rework and SC-23's census re-run. The sibling-path
check resolves to reuse: invariant 8 already makes this exact decision
through §JD-23 (cli/0022:C2, `splitNode` over `terminalKeys` verbatim)
— the discriminator exists, and this RDR confirms it rather than
minting a parallel one. Rejections in one line each: relax converts a
declared miss into measured false accusations; scoped correlation
tracking is measurable cost against an empty measured population,
retained as the named successor arm; external checking re-litigates
0006:ALT4; advisory demotion is foreclosed by 0006:C17's closed tier.

Premortem: hardened (paragraph). Shipped and failed: a production
workflow model livelocks on a key no terminal reads — a retry-loop tag
correlated with nothing terminal — lint said clean, the trap surfaced
in production, and the incident review read "dead-end guarantee" as
concrete exactness; the declared limit existed only inside the record,
invisible at authoring time. The choice survives because the failure
indicts the clause's *carriage*, not the ranking — closing the miss
(arm 1) would have traded this incident for nine false accusations on
day one, and arm 2's machinery was priced against a then-empty
population. Two hardenings folded in response: C2(c) moves the limit
and its cure into `docs/model-authoring.md` where the guarantee is
read at authoring time, not only into the record; and C1's re-rank
trigger makes exactly this incident the normative reopening condition
— the first real-model member of the D12 residual reopens the ranking
through §JD-23 with Alternative 2 pre-costed as the named successor,
so the failure escalates by rule instead of by re-litigated debate.

Ground-sweep: clean (16 anchors)
(`docs/rdr/0015-dead-end-quantifier-over-merged-nodes/evidence/propose-premortem/ground-sweep.md`;
all CONFIRMED by a fresh-context checker from an anchors-only brief,
including the seven correlated `stage = "dropped"` rows each writing
`status = "abandoned"`, and cli/0022:C1/C2 read via the projector).

Joint-check: fired → 0022 (home: JDR 0001 §JD-23). Context (recorded at 0022's
Stage-2 joint-decision check; symmetric fire): RDR 0022 codifies the
terminal-participating-key split bound with a D12-class accepted
residual for its new `graph-terminal-unreachable` invariant, and
carves that finding's population against `graph-dead-end`'s — the
same quantifier doctrine this RDR's arms may relax, replace, or
confirm. One doctrine must govern both liveness-shaped invariants;
its normative home is the umbrella clause both records cite: JDR 0001
§JD-23. This peer may not later record `clear` against 0022.

Re-run at this RDR's own Stage-2 close (12 open peers, all Draft, no
Final): the 0022 fire above stands ALREADY DISPOSED at §JD-23 — this
proposal is the chartered answer, written through that home. Every
distinctive modify-anchor and contract literal (`checkDeadEnd`,
`hasOutgoingOrdinaryRow`, `splitNode`, `terminalKeys`,
`satisfiesSomeTerminal`, `matchSatisfiable`,
`TestAdvDeadEndExistentialOnMergedNode`, `graph-dead-end`, `D12`,
`REQ-37`, `REQ-111`, `REQ-117`, `REQ-122`, `SC-23`, `JD-23`) hits only
0022; the other 11 peers are grep-clean on all of them. Non-fires
adjudicated as context: 0014 copies `models/rdr.toml` as its gate
oracle's lint-subject fixture (a read, and the chosen arm re-runs no
census, so the Background's arm-(1) coupling never activates);
0016 lists the model among shipped artifacts read-only; 0019 cites it
as a conforming model and edits `docs/model-authoring.md`'s
start-state prose — a different passage of the shared doc than
C2(c)'s dead-end passage, a merge concern, not a shared decision.
Absence arm: no `Final` peers exist, and this proposal converts no
refusal into an acceptance — it confirms an already-shipped
acceptance.

## Alternatives Considered

### Alternative 1: Relax REQ-37 toward REQ-111 exactness

**Description**: Widen the quantifier — range the outgoing-row test
over all multi-valued keys (or all match-participating keys), so a
split node dead on any key is accused. This is the reading REQ-111's
principle supports when read alone, and it closes the D12 residual
entirely. It was implemented and measured during the 0006 launch
rather than argued away (D12's Evidence).

**Pros**:

- Closes the miss: no reachable dead state escapes on any key.
- Aligns with REQ-111's "universal checks must not [read merged
  nodes]" read as a global principle, and with REQ-117's
  clearer-model bias if the resulting findings were genuine.

**Cons**:

- Measured: nine false `graph-dead-end` findings on `models/rdr.toml`
  — the merge decorrelates `stage`/`status` (seven rows write
  `stage = "dropped"` together with `status = "abandoned"`), so the
  cross product manufactures views no path produces and accuses them.
- Inverts the soundness direction: the doctrine tolerates misses and
  forbids false accusations; this arm trades a declared miss for
  measured accusations on a conforming model.
- Breaks REQ-122's zero census; SC-23 must re-run and the checked-in
  model be reworked to merge genuinely distinct keys — REQ-117's
  "clearer model" cure fails here because the accused states are
  abstraction artifacts, not model defects.
- Unilaterally widens the §JD-23 bound both records share, forcing
  cli/0022:C2's rework in the same stroke.

**Reason for rejection**: converts a declared, one-directional miss
into measured false accusations on the conforming model, breaking one
explicit MUST (REQ-122) to honour an inference from a principle
(REQ-111).

### Alternative 2: Scoped correlation tracking

**Description**: Keep REQ-37's split bound, but close the miss by
tracking correlation at a bounded scope — per merged fixpoint node, a
record of the value tuples that actually co-occur over a *pack* of
keys (the terminal-participating keys plus keys written in the same
rows, the "terminal-adjacent" scope D12 leaves uncosted). The
outgoing-row test then ranges over recorded tuples only, so widening
never manufactures decorrelated views: `{stage=dropped, status=draft}`
is absent from every tuple set and cannot be accused. This is the
relational-domain "packing" shape from the abstract-interpretation
tradition (uncorroborated in the local corpora — Research Findings —
so that framing is context, not load-bearing support).

**Pros**:

- Closes the miss soundly: exactness extends to the pack with no
  false accusations, satisfying REQ-111 and REQ-122 simultaneously.
- Bounded by construction: tuple sets are capped by distinct write
  footprints per node, worst case the product of the pack's declared
  domains — not the whole lattice.

**Cons** (the costing D12 left undone):

- **State cost**: the fixpoint join must union tuple sets, so node
  identity and the published node ceiling re-size by pack products —
  a second ceiling population, the blast radius that re-opens
  `reach.go::reach`'s core rather than adding a pass beside it.
- **Semantics cost**: a second, tuple-aware satisfaction/edge
  semantics beside the split view; every consumer of the fixpoint
  (invariants 2, 6, 8, and the 0021 export) must then decide
  tuple-aware or not — the dual-semantics drift the record already
  rejects for path sensitivity ("a path-sensitive reading would be
  exponential and is rejected", D12 quoting 0006's join rule).
- **Pack definition is itself a new contract**: model-dependent,
  unbounded in general, and silently wrong when a correlation spans
  rows outside the pack — the miss returns one scope further out.
- **Benefit is measured at zero**: the residual population has no
  member on the conforming model (REQ-122 census zero, A2) — all cost
  is spent on the adversarial fixture's class.

**Reason for rejection**: measurable cost against an empty measured
population — retained as the named successor arm C1's re-rank trigger
points at, so the pricing above is the head start if the class ever
populates.

### Briefly Rejected

- **External model checker over the concrete space**: exact but
  re-litigates 0006:ALT4 — the authoritative gate stays native lint
  over the model intrastate actually consumes, not a drifting
  duplicate toolchain.
- **Demote the miss to an advisory finding**: the advisory tier is
  closed at four (0006:C17), and an advisory over states the
  abstraction cannot separate would accuse falsely at lower severity
  — the same soundness inversion at a discount.
- **Delete the pin and stay silent**: hides a shipped limit and
  removes the widening tripwire REQ-122's census discipline depends
  on.

## Context

### Background

Deviation D12 from the RDR 0006 launch (`Status: needs author
decision`) — the one open author decision from that implementation —
routed out as a seed rather than left as an unclosed review finding;
tracked as kata `intrastate#pz9z` (roborev job 6155; odc-type:
test-oracle). Phase 3a (FAIL-4) and Phase 3b (ADV-3) independently
reported the miss; Phase 3c measured it (widening the split mints 9
false findings; `lint --model models/rdr.toml --as=json` returns zero
findings at HEAD). The implemented reading follows REQ-37's MUST, and
the accepted miss is pinned by `TestAdvDeadEndExistentialOnMergedNode`
— the pin stays and is the trigger until the successor clause lands.
SC-23 already names a guard-aware-pruning successor RDR as the
expected route. Sequencing: this RDR decides the *predicate*; RDR 0014
(kata `9en6`, Proposed 2026-08-28) decides the gate's *testable
surface* — deciding either
leaves the other open, though arm (1) would re-run the SC-23 census
RDR 0014 also touches.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces: the graph-lint
dead-end check (invariant 2) in `internal/graphlint`
(`analysis.go::checkDeadEnd`), the checked-in
`models/rdr.toml` (seven rows correlate `stage = "dropped"` with
`status = "abandoned"`), pin test
`TestAdvDeadEndExistentialOnMergedNode`. Governing record: RDR 0006
(REQ-37, REQ-111, REQ-117, REQ-122, SC-23, deviation D12).

## Research Findings

### Investigation

Prior art was read before enumeration
(`docs/rdr/0015-dead-end-quantifier-over-merged-nodes/evidence/research/prior-art.md`).
⚠ no prior-art coverage in the local corpora for either external
claim: the relational-packing framing behind Alternative 2 (two
queries, citation-fragment noise) and the peer-tool instance question
(one query here plus sibling 0022's four recorded negatives for the
shared problem class) — both stay Resolve-side context, and the choice
rests entirely on in-repo quotable evidence: D12's measured widening
experiment and REQ texts
(`docs/rdr/0006-graph-lint-authority-and-guarantees/artifacts/req-list.md`,
`artifacts/deviations.md`), the shipped quantifier
(`analysis.go::checkDeadEnd`, `::hasOutgoingOrdinaryRow`), the pin
test, and the §JD-23 doctrine with cli/0022:C2. Constraints that
shaped the choice: the no-amend rule on locked 0006 (hence a successor
clause here, not an edit there), §JD-23's one-doctrine requirement
across invariants 2 and 8, and the advisory tier closed at four
(0006:C17).

### Key Discoveries

- **Verified** (source read at HEAD): `checkDeadEnd` splits every
  reachable node via `splitNode` over `terminalKeys()` and tests
  `hasOutgoingOrdinaryRow` existentially per split node; the call-site
  doc already books the residual as "the accepted false-NEGATIVE the
  record books against invariant 2 … See D12".
- **Verified** (rg count at HEAD): `models/rdr.toml` carries seven
  `stage = "dropped"` rows, matching D12's correlation claim.
- **Documented** (D12, artifacts/deviations.md): the widening was
  implemented and measured — nine false `graph-dead-end` findings on
  the conforming model; "Both readings have record support and the
  record does not rank them"; "This entry is the trigger."
- **Documented** (req-list.md): REQ-37's MUST ("never the whole
  lattice"), REQ-111's split remedy, REQ-117's cure direction,
  REQ-122's zero-census pass condition.
- **Assumed** (A1/A2 spikes at Resolve): the nine-finding direction
  and the zero census still hold at HEAD — the two measurements the
  ranking's arithmetic rests on.

## Trade-offs

### Consequences

- Positive: the guarantee becomes honest — the quantifier is named,
  the residual has a canonical class and a published cure, and the
  D12 open author decision closes with zero behavior change and the
  REQ-122 census untouched.
- Positive: both liveness invariants are bound through one anchor
  (§JD-23); cli/0022 proceeds with its C2 unchanged.
- Negative: the miss stays — a real model can livelock on a
  non-terminal key with lint silent. Accepted deliberately, with C1's
  re-rank trigger as the escalation path.
- Negative: REQ-111's broad reading is narrowed by subordination — a
  recorded reduction in what "exactness" promises.

### Risks and Mitigations

- **Risk**: the D12 residual populates on a real model after lock.
  **Mitigation**: C1's normative re-rank trigger reopens the ranking
  through §JD-23 with Alternative 2 pre-costed as the named successor.
- **Risk**: the clause drifts across its three carrier surfaces.
  **Mitigation**: all three carry the same `cli/0015:C1` citation
  (greppable), and the pin test fails on any widening regardless of
  what the prose says.

### Failure Modes

Visible: any widening of the split fails
`TestAdvDeadEndExistentialOnMergedNode` — the tripwire — and C2(b)
forces SC-23's census re-run before acceptance. Silent: the declared
miss itself — a dead state separable only on a non-terminal key ships
unaccused; diagnosis path is `docs/model-authoring.md`'s limit-and-cure
passage (C2(c)), which names the class and the authored cure; recovery
is C1's re-rank trigger. A developer who suspects the miss reproduces
it by checking whether the stuck state's distinguishing key appears in
any terminal predicate (`analysis.go::terminalKeys` is the deciding
set).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified

### Minimum Viable Validation

1. `go test ./internal/graphlint -run
   TestAdvDeadEndExistentialOnMergedNode` — green, with the pin's doc
   comment now citing `cli/0015:C1` as the settled home.
2. `intrastate lint --model models/rdr.toml --as=json` — zero
   findings: the REQ-122 census is byte-unchanged by this RDR.
3. `rg "cli/0015:C1"` — hits on all three carrier surfaces:
   `internal/graphlint/analysis.go`,
   `internal/graphlint/adversarial_0006_test.go`, and
   `docs/model-authoring.md`.
4. The D12 entry in 0006's `artifacts/deviations.md` carries the
   appended resolution line naming cli/0015.

End state: the record ranks the two readings, the ranking is readable
at every surface an author or implementer meets it, and lint behavior
is unchanged.

### Phase 1: Record and Carrier Surgery

#### Step 1: Cite the clause at the call site

Update `analysis.go::checkDeadEnd`'s and `::hasOutgoingOrdinaryRow`'s
doc comments to cite `cli/0015:C1` and JDR 0001 §JD-23 as the
ranking's settled home, replacing the open-question "See D12" framing
(C2(a)).

#### Step 2: Re-document the pin

Rewrite `TestAdvDeadEndExistentialOnMergedNode`'s doc comment from
needs-author-decision framing to declared-limit witness plus widening
tripwire, citing `cli/0015:C1`; confirm the positive arm's coverage
(A5) rather than assuming it (C2(b)).

#### Step 3: State the limit where the guarantee is read

Extend `docs/model-authoring.md`'s dead-end passage with the
limit-and-cure text of C2(c).

#### Step 4: Close D12

Append the resolution line to the D12 entry in 0006's
`artifacts/deviations.md` — resolved by cli/0015, entry body
untouched.

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

- RDR 0006 artifacts: `artifacts/req-list.md` (REQ-37, REQ-111,
  REQ-117, REQ-122), `artifacts/deviations.md` (D12),
  `artifacts/verification.md` (SC-23).
- JDR 0001 `docs/jdr/0001-resolve-kernel-seam.md` §JD-23 — the
  merged-node liveness-quantifier doctrine's home.
- cli/0022:C1, cli/0022:C2 — invariant 8's population carve and shared
  split bound (via projector).
- Source reviewed: `internal/graphlint/analysis.go`,
  `internal/graphlint/adversarial_0006_test.go`, `models/rdr.toml`,
  `docs/model-authoring.md`.
- Prior-art search record:
  `docs/rdr/0015-dead-end-quantifier-over-merged-nodes/evidence/research/prior-art.md`.
- Related: kata `intrastate#pz9z` (1538), roborev job 6155.
