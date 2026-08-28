# Recommendation 0017: Per-finding code identity on mixed-disposition rows

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
- **Profile**: foundational — provisional: one contract (what
  `findings[].code` is for a multi-subject producer) binding the
  `Finding` surface co-owned by RDRs 0005/0006/0008/0009.
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
- **Related Issues**: kata `intrastate#sd16` (1549)
- **Predecessors**: 0005-skill-integration-cli-contract,
  0006-graph-lint-authority-and-guarantees
- **Seam Lineage**: no prior accretion

## Problem Statement

A structured consumer — a skill branching on `findings[i].code` —
reading a mixed-disposition refusal sees `flow-gate-denied` on a gate
whose actual verdict was allow, mis-attributes the refusal, and
repairs the wrong accessor. Live at HEAD whenever a selected row
carries both an allowing and a refusing gate:
`internal/cli/flow_exec.go::gateFindingCode`'s default arm re-emits
the envelope's refusal code, falsifying its own doc comment ("Each
finding's own `code` names THAT gate's disposition"). The envelope
code, exit status, and each finding's prose `message` are correct
(REQ-17), so only structured consumers are misled.

The decision, made once: what *is* `findings[i].code` for a producer
that emits one entry per subject where subjects have differing
dispositions — an identity/carrier call, not a bug fix. Two
populations share the key today: the envelope's closed `flow-*`
refusal-code table (exit-group-bearing) and REQ-24's per-entry
category slug (no exit group). The decision fixes (a) whether
per-finding `code` is a same-namespace echo of the envelope code or a
producer-local discriminator from a separate slug vocabulary; (b)
whether naming an allowing gate therefore requires new *table* surface
(subject to the Stage-8 "additive is not exempt" rule) or merely a new
slug; and (c) whether a non-refusing subject inside a refusal envelope
is entitled to a `code` of its own, or is carried by `message` prose
alone under REQ-17. Three live arms: mint `gate_allowed`; carry the
disposition in `class` — refuted on the record (RDR 0005's field
partition assigns `class` to RDR 0006, and its shipped meaning is the
failure class of an *escape-scoped* finding, which REQ-51 says a gate
refusal never is); or leave the behavior and correct the doc comment,
since no REQ requires a per-gate structured code at all. Whichever arm
wins must state the rule generally enough to bind future multi-subject
findings, not just gates.

## Critical Assumptions

- **A1 [No shipped structured consumer depends on the allow-gate echo
  — none branches on, validates `findings[].code` against the closed
  `flow-*` table, or counts findings whose code equals the envelope
  code (premortem P-3/P-9)]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: a consumer census — sweep skills/scripts/tests for
    every reader of `findings[].code` (branch, validate-against-table,
    and count-equal-to-envelope patterns), not only the branch case.
  - **If wrong**: minting `gate_allowed` silently changes that
    consumer's outcome (a table-validator hard-fails on the unknown
    slug; a counter's total shifts); surfaces after upgrade as a skill
    mis-handling or rejecting a mixed-disposition refusal rather than
    as a test failure here.
- **A2 [No locked clause anywhere — REQ text, 0005 contract, schema,
  enum, or lint rule — closes `findings[].code` over the `flow-*`
  table or otherwise pins the allow-gate finding's literal;
  REQ-51/97/98 fix the envelope code and the one-entry-per-gate
  cardinality only (premortem P-1/P-2)]**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: `0005:C1`, and the REQ-51, REQ-97, REQ-98 rows of
    the owning record's req-list artifact, plus a sweep of the same
    record's contracts/verification artifacts and the contract doc for
    any blanket "codes come from the table" clause or machine schema
    constraining finding codes — this is an absence claim, so the
    sweep must be recorded, not asserted.
  - **If wrong**: `gate_allowed` is a deviation against a locked
    contract, not an open cell — the arm needs a 0005 deviation record
    or a different carrier, and the Stage 8 gate fails on the mismatch.
- **A3 [`gate_allowed` collides with no existing per-finding `code`
  population — no load-category slug or other findings producer emits
  that literal]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: grep for the literal across `internal/` and docs;
    2026-08-28 spot-check found its only near-neighbour is the
    `gate_passed` sample tag *key* in `docs/cli-output-contract.md`'s
    `unknown[]` example, a different field in a different surface.
  - **If wrong**: the discriminator stops discriminating — a consumer
    branching on `gate_allowed` cannot tell which producer emitted it;
    pick a differently-spelled slug.
- **A4 [External prior art matches the rule: peer diagnostic/test
  surfaces give a non-failing subject inside a failure report its own
  per-subject disposition marker, separate from the process status
  (TAP-style `ok`/`not ok`)]**
  - **Status**: Pending
  - **Method**: Prior Art
  - **Evidence**: demoted from the Propose external pass — local
    corpora had no quotable instance coverage (see
    `0017-per-finding-code-identity/evidence/research/propose-prior-art.md`);
    the choice rests on the in-repo citations, not on this.
  - **If wrong**: the rule loses external corroboration only; the
    in-repo contract-doc definition and REQ-24 precedent still carry
    the choice.
- **A5 [The identity rule binds later-added finding classes without
  amendment — a future multi-subject producer on this `Finding` record
  can conform by declaring its own per-finding vocabulary]**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: check at Resolve against the then-proposed sibling
    RDRs that add finding classes to this record (terminal-reachability
    and dead-end-quantifier seeds) — their findings must be expressible
    as own-disposition codes under C1.
  - **If wrong**: C1 needs a carve-out per producer, which is the
    per-case ambiguity this RDR exists to close; surfaces as a sibling
    RDR unable to state its finding codes without echoing an envelope.

## Proposed Solution

### Approach

Adopt the identity rule the shipped output contract already states and
the gate producer violates: `findings[i].code` is the finding's **own
discriminator** — it names *that* finding's disposition or category
from a vocabulary its producer declares, and is never inherited from
the envelope or from a sibling subject. `docs/cli-output-contract.md`
§Structured findings already defines the field this way ("`code` | the
finding's own discriminator — for a model-load failure, the load
category slug") ⇒ the allow-gate echo is the defect and the rule is a
fix *toward* the contract, not new doctrine.

Concretely, for the gate-refusal producer: a denying gate's finding
keeps `flow-gate-denied` and an indeterminate one keeps
`flow-gate-indeterminate` — those literals *are* those findings' own
dispositions, so their coincidence with the envelope table is truthful,
not an echo — and an allowing gate's finding takes the new
producer-local slug `gate_allowed`, in REQ-24's category-slug family:
no exit group, no new row in the closed `flow-*` refusal-code table.
The refusal answer to Problem Statement fork (c) is therefore *yes*: a
non-refusing subject inside a refusal envelope is entitled to a `code`
of its own, because the record requires `code` of every finding and a
false code is worse than a prose-only channel.

### Technical Design

One function changes and two surfaces tighten. The producer seam is
`internal/cli/flow_exec.go::gateFindingCode` (sole caller:
`gateRefusal`, which builds one finding per `gateResult`); its writer
is `runGates`, which populates each `gateResult.Result` from the closed
`accessor.Verdict` vocabulary (`allow`/`deny`/`indeterminate`,
`internal/accessor/model.go::VerdictAllow`) and converts every
non-verdict outcome to an accessor refusal (exit 3) before findings
exist ⇒ the disposition mapping below is total; no default-arm echo
remains reachable. The doc table in `docs/cli-output-contract.md`
§Structured findings gains the general rule and the gate vocabulary,
and the per-gate test tightens from `param`-presence to
code-matches-disposition.

#### Normative Contracts

The identity rule — general; binds every `findings[]` producer on the
`clierr.Finding` record:

**C1**

```normative
findings[i].code is the finding's OWN discriminator. It MUST name THAT
finding's disposition or category, drawn from a per-finding vocabulary
its producer declares. It MUST NOT be inherited from the envelope code
or copied from a sibling subject's disposition. An envelope refusal-
table literal (`flow-*`) MAY appear as a finding's code only when the
finding's own disposition IS that refusal. A subject whose disposition
has no refusal code of its own (a non-refusing subject reported inside
a refusal envelope, or a category with no exit group) takes a
producer-local slug in REQ-24's category-slug family: underscore-
spelled, no exit group, no row in the flow-* refusal-code table.
A producer-local slug MUST be declared in docs/cli-output-contract.md
§Structured findings before a producer emits it, and its meaning is
read relative to the envelope code it appears under; the doc section
is the slug registry — an undeclared slug is a contract defect, not
free latitude.
```

A later-added finding class on this record conforms by declaring its
per-finding vocabulary under this rule (A5); C1 does not enumerate
producers. The declaration clause is the governance answer to "no
table row, so no spec review": the slug family is reviewed surface
too — its registry is the doc section rather than the exit-group
table (premortem P-8/P-11).

**C2** — the gate-refusal instance

```normative
Within a gate refusal (envelope `flow-gate-denied` or
`flow-gate-indeterminate`), each findings[] entry's code is a total
function of that gate's own verdict:
  deny          -> "flow-gate-denied"
  indeterminate -> "flow-gate-indeterminate"
  allow         -> "gate_allowed"
`gate_allowed` is a producer-local slug (C1): it carries no exit
group, adds no refusal-code-table row, and never appears as an
envelope code. Envelope code, exit group, findings cardinality
(one entry per gate on the selected row), `param` = gate id, and
`message` prose are unchanged (REQ-50/51/97/98, REQ-17).
```

**C3** — the observable surfaces

```normative
docs/cli-output-contract.md §Structured findings MUST state C1's rule
and the C2 gate vocabulary (the `gate_allowed` literal named), and
MUST state the level semantics: a `flow-*` literal at the envelope is
the row's refusal, the same literal inside findings[] is one gate's
own verdict — a consumer attributes the refusal from the envelope
code, never by scanning findings[] for a flow-* literal.
TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal MUST
assert, on MIXED-disposition fixtures (at least one allowing plus one
refusing gate on the selected row), that every finding's `code` equals
the C2 mapping of that gate's verdict — not only that `param` names
the gate — and that the envelope code and exit group are unchanged by
the per-finding codes.
```

#### Load-Bearing Decisions

- **Identity** — a finding's subject identity stays `param` (the gate
  id; unchanged); `code` is the *disposition/category* discriminator,
  so two findings with equal `code` share a disposition class, never
  necessarily a subject. The envelope code identifies the refusal; a
  finding's code identifies only itself (C1).
- **Naming** — `gate_allowed`, underscore-spelled: underscore marks
  REQ-24's producer-local slug family, dash `flow-*` marks the
  exit-group-bearing envelope table, so the spelling itself carries
  the namespace. Rejected: `flow-gate-allowed` (claims a refusal-table
  row for a non-refusal), `gate-allowed` (dash would visually claim
  the envelope namespace), `allowed` (drops the producer prefix the
  slug family carries), `gate_passed` (already appears in
  `docs/cli-output-contract.md`'s `unknown[]` example as a sample
  authored tag *key*; reusing the spelling would invite conflation).
- **Selection / predicate** — the C2 mapping is total over the closed
  `accessor.Verdict` vocabulary because its writer,
  `internal/cli/flow_exec.go::runGates`, converts every non-verdict
  outcome (timeout, execution failure) to an accessor refusal at exit
  3 before any finding is built ⇒ `gateFindingCode` needs no
  fourth arm and MUST NOT keep an envelope-echoing default.

#### Illustrative Code

Shape of the producer change (Illustrative — the contract is C2's
mapping, not this text):

```go
// gateFindingCode names one gate's own disposition (C2). The envelope
// code is no longer an input: a finding's code is never inherited (C1).
func gateFindingCode(verdict accessor.Verdict) string {
    switch verdict {
    case accessor.VerdictDeny:
        return codeGateDenied
    case accessor.VerdictIndeterminate:
        return codeGateIndeterminate
    default: // VerdictAllow — runGates admits no other verdict here
        return codeGateAllowed // "gate_allowed"
    }
}
```

Illustrative mixed-disposition refusal, envelope `flow-gate-denied`:

```json
{"code":"flow-gate-denied","findings":[
  {"code":"gate_allowed","param":"permits",
   "message":"the gate `permits` answered allow"},
  {"code":"flow-gate-denied","param":"refuses",
   "message":"the gate `refuses` denied the transition: ..."}]}
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Per-gate disposition naming | `internal/cli/flow_exec.go::gateFindingCode` | default arm echoes the envelope code on allow | Extend | drop the `envelope` parameter; mint `codeGateAllowed` (C2) |
| Producer-local per-finding slug family | `internal/cli/flow_input.go::loadFindings` (`Code: string(category)`, REQ-24) | none | Reuse the pattern | `gate_allowed` joins the underscore-slug family; no table row |
| Disposition vocabulary | `internal/accessor/model.go::Verdict` (`allow`/`deny`/`indeterminate`) | verdict strings are not finding codes | Reuse as derivation source | C2 maps verdicts to codes; totality guaranteed by `runGates` (writer) |
| Per-gate structured reporting on the success path | `internal/cli/flow_exec.go::gateResult` (`gates[]` under `flow next --evaluate-gates`) | reports `result`, not `code` | Reuse as agreement oracle | the two verbs keep agreeing about what each gate said (existing test's cross-verb check) |

### Decision Rationale

The deciding rows of the QOC matrix (Alternatives Considered) are
*prior-art alignment* and *contract blast radius*. The shipped contract
doc already defines `code` as "the finding's own discriminator" —
sibling-path check: the discriminator decision *exists*; exhibits are
that field-table row, `internal/cli/flow_input.go::loadFindings`
(populates `Code: string(category)` — the REQ-24 slug precedent), and
`gateResult.Result` (the success path already names an allowing gate's
disposition structurally) ⇒ this RDR reuses an existing identity signal
rather than inventing a parallel one, and the only arm consistent with
all three exhibits is the chosen one. Blast radius: one function arm,
one new slug, a doc-table addition, a test tightening — no `Finding`
field, no refusal-table row, no envelope change. Full-rename (Alt 1)
buys namespace purity C1 does not need at consumer-breaking cost;
`class` (Alt 2) is refuted on the record; status quo (Alt 3) would have
to *amend the contract doc* to bless an echo — a larger contract change
than the slug it avoids. Consumers doing exhaustive matching over
`flow-*` finding codes see a new, unrecognized `gate_allowed`; that is
the intended fail-loud direction (better unrecognized than falsely
attributed), folded into Failure Modes and A1.

Premortem: hardened (hardened) — critic PASS at
`0017-per-finding-code-identity/evidence/propose-premortem/critic.md`;
folded: C1 registry clause (P-8/P-11), C3 level semantics + mixed
fixture (P-6/P-12), A1 census widened (P-3/P-9), A2 made an explicit
absence-sweep (P-1/P-2), Alt-1 rejection de-contradicted (P-5), P-13
competitor answered in Briefly Rejected, P-4 resolved by the
derivation-independence discovery.
Ground-sweep: clean (17 anchors)
Joint-check: clear (12 peers) — context: 0019 and 0020 cite
`docs/cli-output-contract.md` as environment (verb-I/O listing; §Set
values on the wire); neither touches this RDR's modify-anchor
(§Structured findings) or any contract literal; no Final peers, so
the absence arm is vacuous.

## Alternatives Considered

Scored QOC matrix (foundational profile). Scores 0–2 per cell, higher
better; the choice falls out of the totals, and the deciding rows are
prior-art alignment and contract blast radius.

| Criterion | A: own-disposition code, mint `gate_allowed` (chosen) | B: full producer-local rename | C: disposition in `class` | D: status quo + doc fix |
| --- | --- | --- | --- | --- |
| Correctness fit (structured consumer attributes the refusal per gate) | 2 — every finding's code is its own disposition | 2 — same property | 1 — right fact, wrong field: consumers must learn a second channel | 0 — mis-attribution stays live |
| Prior-art alignment | 2 — matches the doc's "finding's own discriminator" + REQ-24 slug family + `gates[]` sibling | 1 — matches the doc definition but renames literals no prior art demands | 0 — refuted: REQ-18 partition assigns `class` to 0006; escape-scope meaning; REQ-51 "never an escape class" | 0 — contradicts the shipped doc definition of `code` |
| Contract blast radius | 2 — one slug, one function arm, doc row, test tightening | 0 — every gate-finding code changes; consumer-visible break; REQ-97/98-adjacent churn | 1 — semantics change on another RDR's field | 2 — doc comment only, but leaves the doc-contract contradiction standing |
| Reversibility | 2 — additive slug; trivially revertible | 1 — rename churns both ways | 0 — cross-RDR field-meaning change is sticky | 2 |
| Cost | 2 — small | 1 — medium | 1 — small code, high contract cost | 2 — trivial |
| **Total** | **10** | **5** | **3** | **6** |

### Alternative 1: Full producer-local vocabulary

**Description**: Per-finding codes never draw from the envelope table:
`gate_denied`, `gate_indeterminate`, `gate_allowed` on findings; the
envelope keeps `flow-gate-denied`/`flow-gate-indeterminate`.

**Pros**:

- Purest namespace separation — a finding code and an envelope code can
  never be confused, even lexically.

**Cons**:

- Renames two shipped literals that are already *truthful* — a denying
  finding's own disposition IS `flow-gate-denied` — so consumers break
  for zero information gain.
- Two spellings for one disposition (deny = `flow-gate-denied` at the
  envelope, `gate_denied` at the finding) invites mapping bugs in every
  consumer.

**Reason for rejection**: C1 does not require lexical disjointness —
only non-inheritance; the mixed list is principled (a table literal
appears iff the finding itself refuses). Consumer exposure is an
*assumption*, not a settled fact (A1): the chosen arm's exposure is
bounded to allow-gate findings in mixed-disposition refusals, while
the rename changes every gate finding in every refusal — it multiplies
the surface A1 must clear while buying purity, not correctness. (This
is deliberately not the contradiction premortem P-5 names: neither arm
claims zero breakage; both are gated on the same census, and the
chosen arm minimizes what the census must clear.)

### Alternative 2: Carry the disposition in `class`

**Description**: Leave `code` echoing the envelope; add
`class: allowed|denied|indeterminate` per finding.

**Pros**:

- No new code vocabulary.

**Cons**:

- Refuted on the record: RDR 0005's field partition (REQ-18) assigns
  `class` to RDR 0006, and its shipped meaning
  (`internal/cli/clierr/clierr.go::Finding.Class`) is the failure class
  of an *escape-scoped* finding — which REQ-51 says a gate refusal
  never is.

**Reason for rejection**: repurposes another RDR's field against its
shipped meaning, and still leaves `code` violating the doc-contract
definition.

### Alternative 3: Status quo plus doc-comment fix

**Description**: Keep the envelope echo; fix `gateRefusal`'s doc
comment; `message` prose (REQ-17) is the supported per-gate channel.
No REQ requires a per-gate structured code.

**Pros**:

- Zero behavior change; zero consumer risk.

**Cons**:

- The structured mis-attribution the Problem Statement opens with stays
  live; parsing verdicts out of `message` prose is the anti-pattern
  structured findings exist to remove.
- `docs/cli-output-contract.md` §Structured findings already defines
  `code` as "the finding's own discriminator" — blessing the echo means
  amending the *contract doc*, a larger contract change than the slug.

**Reason for rejection**: the contract doc already decided the field's
identity; this arm re-decides it in the wrong direction to save one
small function change.

### Briefly Rejected

- **Disposition in `severity`**: same REQ-18 partition objection as
  `class` (`severity` is 0006's), and severity ranks — it does not
  classify.
- **New per-finding `disposition`/`verdict` field reusing the
  `gates[]` vocabulary (premortem P-13's "strongest competitor")**:
  new public surface on a record co-owned by 0005/0006/0008/0009 (the
  Stage-8 additive rule; RDR 0009's D3 shows what one added field
  costs), and it does not fix the defect — `code` would still carry
  the false `flow-gate-denied` echo on an allowing gate, so either the
  echo stays (mis-attribution survives for every consumer reading
  `code`) or `code` changes too (and the field is redundant).
- **New `flow-gate-allowed` refusal-table row**: the `flow-*` table is
  exit-group-bearing refusal surface (REQ-97/98 shape); an allow is not
  a refusal, so the row would mint an exit group for a non-refusal.

## Context

### Background

Observed by roborev and filed per the triage tie-breaker (undecided
drop-vs-file → file); tracked as kata `intrastate#sd16`. Scope review
judged it RDR-shaped: every available fix is a contract change, and
the fork is *which* contract. Triage then refuted the seed's own
recommendation (reusing `class`) against RDR 0005's field-ownership
partition and the `Finding` struct's shipped meaning, and surfaced the
third arm the seed never named — do nothing, and let `message` prose
be the supported per-gate channel under REQ-17. RDR 0005 fixes only
"one `findings[]` entry per gate" (REQ-51/97/98); it never says which
vocabulary `findings[i].code` draws from. If a structured arm wins,
`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal` must
assert each finding's `code` matches that gate's disposition, not only
its `param`. No data-safety impact.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/cli/flow_exec.go::gateFindingCode`, the `Finding` struct in
`internal/cli/clierr.go`, test
`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal`.
Governing records: RDR 0005 (refusal-code table; REQ-17, REQ-18,
REQ-24, REQ-51, REQ-97, REQ-98; field-ownership partition), RDR 0006
(`class` field owner); `Finding` co-owned by RDRs 0005/0006/0008/0009.

## Research Findings

### Investigation

Prior art was read before enumerating (record:
`0017-per-finding-code-identity/evidence/research/propose-prior-art.md`).
In-repo: `docs/cli-output-contract.md` §Structured findings field table
("`code` | the finding's own discriminator — for a model-load failure,
the load category slug") ⇒ the field's identity is already documented
and the gate producer violates it; RDR 0005 req-list rows REQ-17/18/24/
51/97/98 ⇒ the envelope code and cardinality are locked, the per-gate
finding vocabulary is not; source reads of
`internal/cli/flow_exec.go::gateFindingCode`/`gateRefusal`/`runGates`
and `internal/cli/flow_input.go::loadFindings`. External pass (bounded,
foundational): two arc queries (`StateMachineRes`, `DevRef`); best hit
was problem-detail-document framing (RESTful Web APIs, p.328) —
supporting context only. ⚠ no prior-art coverage in local corpora for
the instance question (how peer CLIs code a *non-failing* subject
inside a failure list); the TAP-style analogy is from the model prior
and is demoted to A4 rather than leaned on.

### Key Discoveries

- **Documented** — the contract doc defines per-finding `code` as "the
  finding's own discriminator" (§Structured findings field table) ⇒
  the chosen rule is a fix toward the shipped contract, not new
  doctrine. (Verbatim-quote re-check of this clause and any nearby
  schema/enum is A2's Resolve sweep — premortem P-1.)
- **Documented** — REQ-24's population is the precedent:
  `internal/cli/flow_input.go::loadFindings` sets
  `Code: string(category)` — a producer-local underscore slug with no
  exit group ⇒ the slug family C1 generalizes already ships.
- **Documented** — REQ-51 text fixes the envelope code and one-entry-
  per-gate cardinality; no REQ row read at Propose assigns per-finding
  `code` values for gate findings ⇒ the cell is open (absence claim —
  A2 owns the exhaustive sweep).
- **Documented** — derivation independence (premortem P-4):
  `internal/cli/flow_exec.go::gateVerdictFailure` computes the envelope
  code by counting `gateResult` verdicts *before* `gateRefusal` builds
  findings, and the exit group is set on the `CLIError` (`GroupUserEnv`)
  ⇒ envelope code and exit group are not computed from finding codes;
  changing an allow finding's code cannot move them.
- **Documented** — the success path already names an allowing gate's
  disposition structurally: `runGates` writes
  `gateResult.Result = string(result.Verdict)` from the closed
  `accessor.Verdict` vocabulary and converts non-verdict outcomes to
  exit-3 accessor refusals first ⇒ the C2 mapping is total and the
  sibling `gates[]` surface stays the cross-verb agreement oracle.
- **Assumed** — no consumer depends on the echo (A1); external
  per-subject-disposition prior art (A4).

## Trade-offs

### Consequences

- Positive: a structured consumer attributes the refusal per gate from
  `findings[i].code` alone; the doc-contract contradiction closes; the
  identity rule (C1) binds future multi-subject producers, including
  the finding classes sibling RDRs will add to this record (A5).
- Positive: the per-finding vocabulary gains a registry (the doc
  section), so producer-local slugs are reviewed surface.
- Negative: `gate_allowed` is new consumer-visible vocabulary; a
  consumer with an exhaustive/validating match over finding codes sees
  an unrecognized value (intended fail-loud direction, but a behavior
  change — A1).
- Negative: one findings list mixes dash table literals and an
  underscore slug; without the C3 doc rule the mix reads as accidental.

### Risks and Mitigations

- **Risk**: a consumer validates `findings[].code` against the closed
  `flow-*` table and hard-fails on `gate_allowed` (premortem P-3).
  **Mitigation**: A1's census covers validators and counters, not only
  branchers; the C1 registry clause gives validators a documented
  vocabulary to widen to.
- **Risk**: a locked schema/enum or blanket clause closes finding
  codes over the table and vetoes the slug late (premortem P-2).
  **Mitigation**: A2 is an explicit recorded sweep at Resolve; if it
  fires, the arm needs a 0005 deviation record before Stage 8.
- **Risk**: the two-level `flow-*` homonym — a consumer scanning all
  codes in the payload keys on a finding's `flow-gate-denied` under a
  `flow-gate-indeterminate` envelope and mis-attributes (premortem
  P-6). **Mitigation**: C3's level-semantics doc rule (attribute from
  the envelope, per-gate from findings) and the mixed-fixture test.

### Failure Modes

- Visible: `TestFail1_…GateRefusal` fails when any finding's `code`
  mismatches its gate's verdict under the C2 mapping, or when the
  envelope code/exit group shift — the regression a future echo-style
  change would trip.
- Silent: a downstream consumer treating finding codes as a closed set
  logs/rejects `gate_allowed` without failing this repo's suite; only
  A1's census and release notes cover it. Diagnosis: the finding's
  `message` prose states the verdict (REQ-17) and the envelope code
  still names the refusal, so a developer comparing `message` with
  `code` sees the mismatch immediately.
- Recovery: revert is one function arm (the slug is additive); no
  stored data or wire history is affected.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1's consumer census and
      A2's closed-vocabulary sweep are the gating two)

### Minimum Viable Validation

1. Load a model whose selected row carries one allowing and one
   denying gate; run `flow resolve … --as=json` on state matching that
   row.
2. Expect exit group 2, envelope `code` = `flow-gate-denied`, and
   `findings[]` with exactly two entries: the allowing gate's entry
   carrying `code` = `gate_allowed` and the denying gate's carrying
   `code` = `flow-gate-denied`, each `param` its gate id and each
   `message` prose unchanged.
3. Repeat with an allowing plus an indeterminate gate: envelope
   `flow-gate-indeterminate`, findings `gate_allowed` +
   `flow-gate-indeterminate`.
4. End-state: the two verbs still agree — `flow next --evaluate-gates`
   reports the same gates' verdicts in `gates[]` for the same row.

### Phase 1: Code Implementation

#### Step 1: Mint the slug and totalize the mapping

Rewrite `internal/cli/flow_exec.go::gateFindingCode` to the C2
verdict-total mapping: add `codeGateAllowed = "gate_allowed"`, drop
the `envelope` parameter (no inherited code remains reachable).

#### Step 2: Tighten the per-gate test

Extend `TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal`
per C3: assert each finding's `code` against the C2 mapping on the
existing mixed-disposition fixtures, and assert the envelope code and
exit group are unchanged.

#### Step 3: Declare the vocabulary in the contract doc

Amend `docs/cli-output-contract.md` §Structured findings per C3: the
C1 rule, the gate vocabulary with the `gate_allowed` literal, and the
level-semantics sentence (attribute the refusal from the envelope
code, never by scanning `findings[]` for a `flow-*` literal).

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

- `docs/cli-output-contract.md` §Structured findings (the `code`
  field-table row; the `unknown[]` example carrying the `gate_passed`
  sample tag key)
- RDR 0005 `artifacts/req-list.md` — REQ-17, REQ-18, REQ-24, REQ-47,
  REQ-50, REQ-51, REQ-97, REQ-98
- `internal/cli/flow_exec.go` (`gateFindingCode`, `gateRefusal`,
  `gateVerdictFailure`, `runGates`, `gateResult`),
  `internal/cli/flow_input.go::loadFindings`,
  `internal/cli/clierr/clierr.go` (`Finding`),
  `internal/accessor/model.go` (`Verdict`)
- kata `intrastate#sd16` (1549) — the filed observation
- Stage-2 evidence:
  `0017-per-finding-code-identity/evidence/research/propose-prior-art.md`,
  `0017-per-finding-code-identity/evidence/propose-premortem/critic.md`
