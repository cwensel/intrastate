# Recommendation 0016: Reader cardinality of the artifact-role read binding

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
- **Type**: Bug Fix
- **Profile**: foundational — provisional: one contract (the
  role read-binding cardinality plus its dependent key-coverage
  obligation) spanning the RDR 0002/RDR 0004 seam.
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
- **Related Issues**: kata `intrastate#p63c` (1544); roborev
  jobs 6215 (ref 7932928), 6220
- **Predecessors**: 0002-transition-table-as-reviewable-data,
  0004-accessor-execution-safety-model
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author whose model passes validation expects a correct write to
read back verified — never a spurious `read_back_incomplete` — and
expects reader resolution to be deterministic rather than
registry-order-dependent. Today RDR 0002 counts reader coverage per
*key* (each owned tag served by exactly one reader, each observed tag
by at most one — `internal/table/load.go`), while RDR 0004's read-back
resolves a reader per *artifact role* and assumes that mapping is a
function (`0004:C12`/`C13` speak of "the role's read definition",
singular). Role is carried but never counted
(`internal/accessor/validate.go` has no role-uniqueness arm), so a
valid model can reach `readerFor` (`internal/accessor/model.go`) with
two same-role candidates and get first-match-over-registry-order
selection. The owned half is the real hazard: nothing binds an owned
key's reader to the *writer's* role, so the post-write compare can
request keys its resolved reader cannot serve, yielding
`read_back_incomplete` after a correct write.

The decision, made once: does **role become an arity-bearing
coordinate of the read binding** — at most one reader per role, a new
RDR 0002 binding validation making `readerFor` total by construction —
or does it stay non-counted, in which case `0004:C12`/`C13` must be
restated over a reader *set*, with read-back selecting and aggregating
by declared keys? The writer-role-must-cover-owned-keys question is
the totality half of the same contract, not an independent one: under
either cardinality arm it is the condition that makes read-back
resolvable for every planned owned key, and its statement form (a
single reader's keys vs a union over the reader set) is determined by
the cardinality answer — one RDR, one fork with a dependent clause.

## Critical Assumptions

- **A1 [No shipped model or test fixture binds two read accessors to
  one artifact role, or a writer role its reader's keys do not cover]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: sweep every `[read.*]`/`[write.*]` table in
    `models/**/*.toml` and the in-repo test fixtures for a repeated
    `role` value among readers and for a writer key set not contained
    in its role's reader key set; the Propose spot-check of
    `models/rdr.toml` and `models/examples/*.toml` found one reader
    per role with reader keys equal to writer keys.
  - **If wrong**: C1/C2's load refusal invalidates a previously-valid
    in-repo model or fixture; the implementing change must migrate it
    (merge same-role readers, or widen the role reader's keys) before
    tests pass.
- **A2 [A write's planned keys are always bounded by the writer
  binding's declared key set at execution]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/accessor/executor.go::nonOwnedPlanKeys`
    refuses any planned key outside the writer definition's own
    declared keys before the write command runs, so C2's coverage
    over declared keys bounds every key the post-write compare can
    request on the planned side.
  - **If wrong**: a plan can carry an owned key the writer never
    declared, C2 under-constrains, and the spurious
    `read_back_incomplete` survives validation exactly as today.
- **A3 [The post-write compared key set is exactly the planned keys
  plus the resolved reader's declared protected keys]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/accessor/executor.go` builds the compared
    set from `plannedKeys` united with `protectedKeys(reader, …)`
    (the reader's own declared keys minus planned), fixed before the
    write command runs; no third source — in particular not the keys
    the external write command actually emits — feeds the re-read's
    requested keys.
  - **If wrong**: keys outside the single reader's declared obligation
    enter the compare, and single-reader coverage (C2) is not the
    complete totality condition — a reader-set statement would be
    needed after all.
- **A4 [C1/C2 leave RDR 0002's per-key authorability mitigation
  intact]**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: `0002:§risks-and-mitigations`, final risk, quoted:
    "A `keys` binding that must cover every declared key would make
    `recognized` and caller-supplied `--tag` keys unauthorable. /
    Mitigation: The binding validation is provenance-scoped — exactly
    one reader per owned key, at most one per observed key, none for
    `recognized`" — the recorded scope forbids *whole-key coverage
    obligations* and says nothing about role cardinality; C1 counts
    read bindings per role and C2 constrains only the keys a write
    binding declares (owned by construction via
    `write_non_owned_tag`), so neither imposes reader coverage on
    recognized or observed keys. Resolve confirms the quote still
    anchors and that no other 0002 element widens the mitigation.
  - **If wrong**: this RDR silently re-fights a decided mitigation and
    models carrying reader-less recognized/observed keys stop
    loading.
- **A6 [Every registry the executor acts on has passed
  `accessor.Validate`]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: sweep every construction/entry path that hands a
    `Registry` to the executor (flow execution, resume/replay if any,
    lint/export) and show each routes through
    `internal/accessor/validate.go::Validate` before a write can run;
    also confirm `Validate` has no early return that can starve a
    later arm (at Propose it appends findings and never returns
    early).
  - **If wrong**: a bypass path reaches `readerFor` with an
    unvalidated registry — C4's fail-closed resolver is the designed
    backstop, refusing rather than reverting to first-match.
- **A5 [`readerFor` is the sole role→reader resolution point]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: sweep call sites of `Registry.readerFor` and any
    other selection keyed on `Accessor.Role` over read definitions;
    at Propose the only role-keyed selection found is
    `internal/accessor/model.go::readerFor`, called from the
    executor's post-write read-back.
  - **If wrong**: another site keeps first-match-over-registry-order
    selection and C1's determinism guarantee does not reach it.

## Proposed Solution

### Approach

Artifact role becomes an **arity-bearing coordinate of the read
binding**, decided at the RDR 0002 load seam, with the totality half
stated over the single reader.

Cardinality: at most one read accessor binding per artifact role,
enforced by a new walk in
`internal/table/load.go::checkAccessorBindings` under the existing
`malformed_accessor_binding` category ⇒ the phrase "the role's read
definition" in `0004:C12`/`0004:C13` denotes a unique binding, and
`internal/accessor/model.go::readerFor` — today first-match over
registry order — becomes deterministic because at most one candidate
can exist.

Totality: every write accessor binding's role must be served by a
read binding whose declared keys include every key the write binding
declares. The executor already refuses any planned key outside the
writer's own declared set
(`internal/accessor/executor.go::nonOwnedPlanKeys`) ⇒ under this
clause every key the post-write compare can request is one the
resolved reader declares, so the spurious `read_back_incomplete`
after a correct write is closed at validation, before any write runs.

Mirror: `internal/accessor/validate.go::Validate` gains two arms
deciding the same two defects over the registry the executor is
actually handed, per the established load/execution split (its own
doc comment: the loader refuses these shapes at load; execution
re-decides them at the boundary).

RDR 0004 is not overridden: `0004:C12`/`0004:C13` stand as written —
this RDR makes their singular reading true by construction rather
than restating them over a reader set.

### Technical Design

Two walks at the load seam (model admission), two mirror arms at the
execution boundary (registry admission), a fail-closed resolution
site, no change to executor control flow or refusal classes. Data
flow: `[read.*]`/`[write.*]` tables → `loadAccessors` fills
`Model.Readers`/`Model.Writers` (writer of the state the walks read)
→ `checkAccessorBindings` walks them → at execution, the registry
carries the same bindings and `Validate` re-decides.

**Consistency corollary (premortem P-4).** Role-designation cannot
diverge from RDR 0002's per-key designation; agreement is a theorem
of the walks, not a fourth arm. For any owned key K a write binding
W (role R) declares: C2 forces R's reader to declare K; 0002's
existing walk forces *exactly one* reader to serve K
(`internal/table/load.go::checkAccessorBindings`, whose writer is
`loadAccessors` filling `Model.Readers` from the `[read.*]` tables);
therefore R's reader **is** K's unique server. For an observed key in
R's reader's declared set, the at-most-one walk makes it the only
server likewise. A model in which the role-designated reader and the
per-key server differ cannot pass both validations.

**Compare-set membership (premortem P-2/P-5).** The post-write
compared set is fixed *before* the write command runs — planned keys
plus the resolved reader's declared protected keys — so keys an
external write command *emits* beyond its declaration never enter the
compare (A3 verifies the construction at Resolve). C2's coverage over
the writer's declared keys is therefore the complete totality
condition; recognized and observed keys stay reader-optional exactly
as 0002 scoped them (A4).

#### Normative Contracts

**C1**

```normative
An artifact role MUST be named by at most one read accessor binding
in a model. A model in which two or more read bindings carry the same
role MUST fail load under category `malformed_accessor_binding`
(`internal/table/load.go::checkAccessorBindings`), before any
execution.
```

**C2**

```normative
Every write accessor binding's role MUST be served by a read accessor
binding, and that reader's declared key set MUST include every key
the write binding declares. A model violating either clause MUST fail
load under category `malformed_accessor_binding`. Coverage is stated
over the single role reader (per C1), never over a union of readers.
```

**C3**

```normative
`internal/accessor/validate.go::Validate` MUST decide the same two
defects over the registry handed to the executor: two or more read
definitions sharing a role → finding code `multiply_bound_role`; a
write definition whose role no read definition serves with full
declared-key coverage → finding code `unserved_writer_role`. One
defect, one code; a registry with either finding never reaches a
write.
```

**C4**

```normative
`0004:C12` and `0004:C13` stand as written; this RDR overrides
nothing. Under C1, `internal/accessor/model.go::readerFor` resolves
"the role's read definition" uniquely — and the resolution site MUST
fail closed: presented with two or more candidate read definitions
for one role, it MUST refuse (no reader found for read-back purposes,
surfacing as the existing `read_back_incomplete` fail-safe), never
select by registry order. Validation (C1/C3) is the primary
mechanism; the fail-closed resolver is defense-in-depth so any future
registry path that bypasses validation fails loudly instead of
silently reverting to first-match.
```

Failure output rides the existing surfaces unchanged: C1/C2 report
through RDR 0002's load-failure envelope (category + message naming
the offending role, binding ids, and — for C2 — the uncovered keys);
C3 reports through RDR 0004's validation-finding envelope. Exact
message wording and finding-field payloads are Resolve/Pre-Lock
detail.

#### Load-Bearing Decisions

- **Identity** — an accessor's *registration* identity stays the
  `(flow, name, capability)` triple (`0004`'s multiply-bound arm,
  `internal/accessor/validate.go`); this RDR adds the read-back
  *resolution* identity: the artifact role. Two read bindings sharing
  a role are one collision defect, never two candidates.
- **Naming** — execution-boundary finding codes
  `multiply_bound_role` and `unserved_writer_role` (rejected:
  overloading `multiply_bound_accessor`, which names the
  identity-triple defect); load-time failures reuse
  `malformed_accessor_binding` (rejected: a new category — per-role
  arity is the same defect class as the existing per-key arity
  walks).
- **Selection / predicate** — when N read bindings share a role,
  *none* is chosen: the model is refused at load and the registry at
  the execution boundary. When N=1, the sole binding is the role's
  read definition. Registry order never participates in selection.

#### Illustrative Code

Illustrative only — intent, not fixture text:

```toml
# Refused by C1: two read bindings, one role.
[read.status-a]
role = "rdr"
keys = ["stage"]

[read.status-b]
role = "rdr"
keys = ["status"]
```

```toml
# Refused by C2: the writer's role is served, but coverage is short —
# "status" is plannable yet unservable by the role's reader, which is
# exactly the shape that today yields read_back_incomplete after a
# correct write.
[read.status]
role = "rdr"
keys = ["stage"]

[write.status]
role = "rdr"
keys = ["stage", "status"]
read_back = true
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Load-time binding arity walk | `internal/table/load.go::checkAccessorBindings` | walks per-key counts, not roles | Extend | C1 role walk + C2 coverage walk, same function and category |
| Execution-boundary validation | `internal/accessor/validate.go::Validate` | eight arms, none role-keyed | Extend | C3's two arms and finding codes |
| Role→reader resolution | `internal/accessor/model.go::readerFor` | first-match over registry order | Reuse | mechanics unchanged; C1 bounds candidates at ≤1 |
| Post-write compare | `internal/accessor/executor.go` | may request keys the resolved reader never declared | Reuse | unchanged; C2 makes every requestable key servable |

### Decision Rationale

Scored QOC matrix (0 = fails the criterion, 1 = partial, 2 = full;
approaches from the prior-art read in Research Findings):

| Criterion | A: role arity + writer coverage (chosen) | B: reader-set semantics | C: deterministic tiebreak | D: execution-time refusal |
| --- | --- | --- | --- | --- |
| Correctness fit (closes nondeterminism AND spurious `read_back_incomplete`) | 2 — both, statically, pre-execution | 2 — both, but decided inside the executor | 0 — determinism only; an ordered pick can still not serve the compared keys | 1 — refusal is deterministic but the coverage hazard needs C2 anyway |
| Prior-art alignment | 2 — registration-refusal (ServeMux) plus three in-repo arity walks | 0 — no read peer aggregates same-key handlers across a set | 1 — SCXML-style document order exists, but its conflicting candidates are interchangeable; same-role readers are not | 1 — right refusal posture, wrong site |
| Reversibility | 2 — a validation can be relaxed to B later; every A-valid model is B-valid | 0 — aggregation and refusal-attribution semantics become locked contract | 2 — an ordering rule is droppable | 1 — a runtime refusal class, once shipped, is contract |
| Blast radius | 1 — refuses model shapes that load today (none shipped; A1) | 0 — executor rewrite plus an override of `0004:C12`/`C13` | 2 — admits everything | 1 — new runtime refusal class on the write path |
| Cost | 2 — two load walks, two validate arms | 0 — multi-reader read aggregation, per-reader refusal attribution, contract restatement | 2 — one sort | 1 — resolution-site rework plus C2 anyway |
| **Total** | **9** | **2** | **7** | **5** |

The deciding rows are correctness fit and reversibility: C's total is
close but its correctness score is 0 — determinism without coverage
leaves the exact user-visible defect (spurious `read_back_incomplete`)
that motivates the RDR, so its cheapness buys nothing; B is the only
other approach that closes both hazards, and it loses on every other
row while requiring an override of two implemented contracts. A is
also the codebase's own answer everywhere the same question has come
up (the three arity walks in the sibling-path exhibit) — choosing B
or C would introduce a second candidate-resolution philosophy beside
them.

Sibling-path exhibit (step-5 check — does an adjacent path already
make this decision?): yes, three do, and A reuses their signal shape
rather than inventing a parallel one:
`internal/accessor/validate.go` (`multiply_bound_accessor` — same
identity triple bound twice is refused, not tiebroken),
`internal/table/load.go::checkAccessorBindings` ("owned tag … served
by %d readers; want exactly one"), and
`internal/cli/flow_state.go::writerFor` (exactly one `[write.<id>]`
per key, kata 8dg3). No adjacent path resolves N candidates by
ordering.

The hardened premortem
(`evidence/propose-premortem/critic.md`, findings P-1…P-11) did not
overturn the choice but forced four hardenings, folded above: C4 now
requires a fail-closed resolution site (P-7/P-10 — validation-bypass
paths must refuse, never first-match; the rejected execution-time
refusal returns as defense-in-depth, not primary mechanism); the
consistency corollary and compare-set membership paragraphs answer
P-2/P-4/P-5/P-11 from the shipped walk structure; A6 was added for
registry-path coverage; and P-3/P-6 (load-refusal of never-writes
models carrying duplicate-role readers) is answered by in-repo
precedent rather than scoping C1 to write-referenced roles —
`internal/table/load.go` already holds arity defects malformed
regardless of use ("a two-writer key outside `written` is malformed
even though nothing in the model writes it"), and A1's corpus sweep
is the evidence obligation that the refused set is empty, with
fixture migration in-plan if it is not.

Premortem: hardened (hardened)
Ground-sweep: clean (16 anchors)
Joint-check: clear (12 peers)

Joint-check context (not a fire): of the fourteen anchor/literal
tokens swept, the sole peer hit is 0024's Technical Environment
naming bare `internal/table/load.go` for a different locus ("the tag
declaration grammar to mirror" — `loadTags`, not
`checkAccessorBindings`), with no shared contract literal.
Absence arm: not applicable — this proposal converts no refusal into
an acceptance; it adds refusals.

## Alternatives Considered

### Alternative 1: Reader-set semantics (B)

**Description**: Keep role non-counted. Override `0004:C12`/`C13`
with a restatement over the *set* of read definitions sharing the
write binding's role: read-back selects per compared key the reader
whose declared keys serve it, aggregates the reads, and attributes a
per-reader refusal on incompleteness. The totality half becomes a
union obligation: the same-role readers' declared keys must jointly
cover the writer's declared keys.

**Pros**:

- Admits every model shape the current loader admits — zero admission
  blast radius.
- Lets one role split cheap and expensive reads across two commands.

**Cons**:

- Overrides two contracts of an Implemented RDR (`0004:C12`,
  `0004:C13`) that are otherwise correct as written.
- Adds multi-reader aggregation, partial-read ordering, and
  per-reader refusal-attribution semantics to the executor — new
  locked contract surface for a shape no in-repo model uses (A1
  sweeps model *shapes*, not file presence, precisely to keep this
  rejection evidence-based rather than asserted).
- Diverges from every adjacent candidate-arity decision in the
  codebase (sibling-path exhibit in Decision Rationale), introducing
  a second resolution philosophy.
- Near-irreversible once shipped: aggregation behavior becomes
  observable contract; the cardinality arm stays freely relaxable in
  the other direction.

**Reason for rejection**: loses every QOC row except correctness;
the only shape it preserves is one nothing uses, at the price of
overriding implemented contracts.

### Alternative 2: Deterministic tiebreak (C)

**Description**: Keep multiple same-role readers legal; define a
total selection order (e.g. lexicographic binding id) at `readerFor`
so selection is deterministic and registry order never shows through.

**Pros**:

- Cheapest change; admits everything; droppable later.

**Cons**:

- Determinism alone leaves the RDR's motivating defect: the
  deterministically chosen reader can still lack the compared keys,
  so the spurious `read_back_incomplete` survives — and the writer
  key-coverage obligation would still have to be stated against
  *some* reader, which is the cardinality question again.
- The statecharts-style document-order precedent does not transfer:
  its conflicting candidates are interchangeable resolutions of one
  event, while same-role readers with different key sets are not
  interchangeable for a compare.

**Reason for rejection**: scores 0 on correctness fit — it answers
the nondeterminism symptom while preserving the hazard the RDR
exists to close.

### Briefly Rejected

- **Execution-time-only ambiguity refusal (D)**: defers to the write
  moment a defect fully decidable from declarations before any
  command runs (both C1 and C2 are static); it survives *as
  defense-in-depth* in C4's fail-closed resolver, but as the primary
  mechanism it reports the model's defect only when a write finally
  exercises the role.
- **Codify registry order as normative**: locks an accident of load
  order into contract and still leaves coverage undecided.
- **Documentation-only ("authors should avoid duplicate roles")**:
  the loader is the project's established enforcement point for
  binding arity; advice is not a validation.

## Context

### Background

Raised by roborev jobs 6215 and 6220; tracked as kata
`intrastate#p63c`. The current code is a correct implementation of the
contract as written — the contract is what is underspecified: RDR 0004
explicitly delegates reader-cardinality to RDR 0002, whose binding
validations are deliberately provenance-scoped per key (with a
recorded mitigation rationale: a whole-key-coverage binding would make
`recognized` and caller-supplied `--tag` keys unauthorable). Deviation
D14, read fairly, is narrow — key *width* under an assumed single
reader — and never counts readers per role, so reader
selection/cardinality is genuinely unadjudicated. The
"cheap validation arm" resolution is not cheap: role-uniqueness in
`internal/accessor/validate.go` narrows which models load, against
0002's deliberate per-key scope, and the spurious
`read_back_incomplete` survives it unless the writer-owned-keys
obligation is also decided. Fail-safe direction holds (`0004:C14`): a
correct write is rejected, never a wrong write accepted — severity
medium.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/table/load.go` (per-key binding validations),
`internal/accessor/validate.go` (identity `(flow, name, capability)`;
no role arm), `internal/accessor/model.go::readerFor`,
`internal/accessor/executor.go` (post-write compare). Governing
records: RDR 0002 (binding validations), RDR 0004 (C5, C12, C13, C14,
REQ-50, REQ-102, deviation D14).

## Research Findings

### Investigation

Prior art was read before enumeration (queries and rejected branches
in `evidence/research/propose-prior-art.md`). ⚠ no prior-art coverage
in the indexed corpora (`StateMachineRes`, `DevRef`,
`StateMachineLit`) for handler-cardinality-at-registration — the
corpus hits were heading-only or off-domain, so the class read rests
on an openable local source and in-repo precedent instead. Go's
stdlib mux refuses a conflicting registration at bind time —
`net/http/server.go::ServeMux.registerErr` (go1.26.6): "pattern %q
(registered at %s) conflicts with pattern %q (registered at %s)" —
ambiguity no specificity rule resolves is refused at registration,
never resolved by registration order at serve time ⇒ supports
refusing same-role duplicates at model load over a runtime tiebreak.
In-repo, the same question has three prior answers, all refusal:
`internal/accessor/validate.go` (`multiply_bound_accessor`),
`internal/table/load.go::checkAccessorBindings` per-key arity walks
(state written by `loadAccessors` from the `[read.*]`/`[write.*]`
tables), and `internal/cli/flow_state.go::writerFor` (kata 8dg3) ⇒
role arity as a fourth refusal walk is signal reuse, not invention.
Code paths shaping the totality half:
`internal/accessor/executor.go::nonOwnedPlanKeys` (plan keys bounded
by the writer's declared keys) and the pre-invocation compared-set
construction ⇒ single-reader declared-key coverage is the complete
condition.

### Key Discoveries

- **Documented** — `internal/accessor/model.go::readerFor` selects
  first-match over `Definitions` order on
  `Capability==CapRead && Accessor.Role==role`; nothing counts
  same-role readers anywhere.
- **Documented** — `internal/accessor/validate.go::Validate` runs
  eight arms, none role-keyed, and never early-returns (findings
  accumulate), so mirror arms cannot be starved by ordering.
- **Documented** — `internal/accessor/executor.go::nonOwnedPlanKeys`
  refuses planned keys outside the writer's declared set before the
  command runs; the compared set is fixed pre-invocation from
  planned ∪ reader-declared keys.
- **Documented** — every shipped model (`models/rdr.toml`,
  `models/examples/*.toml`) binds one reader per role with reader
  keys equal to writer keys — C1/C2 refuse none of them.
- **Documented** — `internal/table/load.go` enforces binding arity
  regardless of use ("a two-writer key outside `written` is
  malformed even though nothing in the model writes it") ⇒ global C1
  (not scoped to write-referenced roles) is the established
  admission philosophy.
- **Assumed** — no *test fixture* relies on same-role duplicate
  readers or uncovered writer roles (A1); and no registry path
  reaches the executor unvalidated (A6).

## Trade-offs

### Consequences

- Positive: the spurious `read_back_incomplete` after a correct
  write becomes a load-time message naming the binding and the
  uncovered keys — caught at authoring, not at the first write.
- Positive: `readerFor` is deterministic (≤1 candidate) and, for
  every write binding's role, resolvable with full coverage — the
  singular reading of `0004:C12`/`C13` holds by construction, with
  no executor or contract change.
- Negative: two model shapes that load today stop loading — same-role
  duplicate readers, and writer roles whose reader under-covers.
  Both are breaking admission changes (project convention: no
  back-compat); the second shape was already broken at run time.
- Negative: a role's single reader must declare the union of keys
  its writers need — one larger read command instead of several
  narrow ones.

### Risks and Mitigations

- **Risk**: an in-repo fixture depends on a shape C1/C2 refuse.
  **Mitigation**: A1's corpus sweep at Resolve; Step 4 migrates any
  hit in the implementing change.
- **Risk**: the coverage clause creeps onto recognized/observed keys,
  re-fighting 0002's recorded scoping mitigation.
  **Mitigation**: C2 is stated over the keys the write binding
  declares, owned by construction (`write_non_owned_tag`); A4 quotes
  the 0002 clause and Resolve re-anchors it.
- **Risk**: a registry path bypasses validation and reaches the
  resolver with duplicates.
  **Mitigation**: A6 sweeps construction paths; C4's fail-closed
  resolver refuses rather than first-matching.

### Failure Modes

- Visible at load: `malformed_accessor_binding` naming the role and
  the duplicate read bindings (C1), or the writer id, role, and
  uncovered keys (C2). Diagnosis is the message; recovery is a model
  edit (merge same-role readers; widen the role reader's keys). No
  state migration — validation only.
- Visible at the execution boundary: `multiply_bound_role` /
  `unserved_writer_role` findings when a registry that skipped load
  validation is handed to the executor (C3), and a fail-closed
  refusal at the resolution site under C4 if both nets are somehow
  bypassed.
- Residual runtime refusal (unchanged, by design): a reader that
  *declares* a key its external command cannot actually produce
  still yields `read_back_incomplete` at run time — static
  validation cannot prove external-tool behavior; `0004:C13` remains
  the fail-safe backstop for declared-but-unreadable keys.
- Partial-implementation hazard: C1 without C2 removes the
  nondeterminism but leaves the spurious `read_back_incomplete`; the
  MVV's coverage leg exists to make that half-ship fail its own
  validation.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A6)

### Minimum Viable Validation

1. Author a fixture model: reader role `r` with keys `["a"]`; writer
   role `r` with keys `["a", "b"]`, `read_back = true`. Load MUST
   fail with the C2 coverage message naming the writer and key `b`.
2. Add `b` to the reader's keys; the model loads. Drive a flow write
   of the owned keys through stub accessor commands; the write
   completes with a green read-back — no `read_back_incomplete`.
3. Duplicate the reader under a second binding id with the same
   role; load MUST fail with the C1 role-arity message.

End state: both refusals observed at load with messages naming the
offending binding, and the corrected model's write verifies green
end-to-end — the pre-fix behavior (step 1's shape loading fine and
refusing at run time) is no longer reachable.

### Phase 1: Code Implementation

#### Step 1: Role-arity walk at the load seam

Extend `internal/table/load.go::checkAccessorBindings` with the
≤1-reader-per-role walk (C1).

#### Step 2: Writer-role coverage walk

Same function: resolve each write binding's role to its (now unique)
reader and require declared-key coverage of the writer's declared
keys (C2).

#### Step 3: Execution-boundary mirror arms and fail-closed resolver

Add the two arms and finding codes to
`internal/accessor/validate.go::Validate` (C3); make the resolution
site refuse on plural candidates (C4).

#### Step 4: Fixture migration

Repair any in-repo model or test fixture A1's sweep surfaces (merge
same-role readers; widen role-reader key sets).

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

- RDR 0002 (`0002:§risks-and-mitigations` final risk — the per-key
  scoping mitigation A4 quotes); RDR 0004 elements `0004:C5`,
  `0004:C12`, `0004:C13`, `0004:C14`.
- Source reviewed: `internal/table/load.go::checkAccessorBindings`,
  `internal/accessor/validate.go::Validate`,
  `internal/accessor/model.go::readerFor`,
  `internal/accessor/executor.go` (`nonOwnedPlanKeys`,
  compared-set construction), `internal/cli/flow_state.go::writerFor`,
  `models/rdr.toml`, `models/examples/*.toml`.
- Prior art: Go stdlib `net/http/server.go::ServeMux.registerErr`
  (go1.26.6) — registration-time conflict refusal. Search record:
  `evidence/research/propose-prior-art.md`.
- Premortem: `evidence/propose-premortem/critic.md` (P-1…P-11).
- Related: kata `intrastate#p63c` (1544); roborev jobs 6215
  (ref 7932928), 6220.
