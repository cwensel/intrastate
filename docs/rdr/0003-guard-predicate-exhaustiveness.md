# Recommendation 0003: Guard Predicate Exhaustiveness

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-06-19
- **Status**: Draft [revised from Final 2026-08-12; re-verify A5 —
  JDR 0001 §D1 fixes the normalized atom shape this RDR must cite rather than
  restate; JD-4 narrows the lint promise where it and runtime disagree, though
  it names RDR 0006 and leaves the recording document open — see A8]
- **Type**: Architecture
- **Profile**: large — locks one guard-predicate contract: symbolic atom grammar plus finite-domain exhaustiveness semantics.
- **Priority**: High
- **Related Issues**: None
- **Predecessors**: 0001-resolution-kernel
- **Overrides**: None
- **Seam Lineage**: no prior accretion

## Problem Statement

A flow author needs conditional edges such as cap-3 handling, profile-to-lens routing, and rewind-target legality to be expressed in a way the lint can prove exhaustive and mutually exclusive. The system-internal requirement is to define how tag predicates are represented without giving up static verification.

## Context

### Background

The model treats guards as predicates on tags, not side inputs. The two target flows need equality, integer comparison, and set-membership; the guard shape has to remain strong enough for the lint to prove determinism.

The real design fork is expressiveness versus verifiability: a fixed enum of operators, a tiny expression grammar, or embedded host predicates.

### Technical Environment

intrastate is a Go CLI wired through `internal/cli`. Guard predicates are consumed by the transition table, the resolver kernel, and the lint that proves the graph is safe before runtime.

## Research Findings

### Investigation

RDR 0001's kernel now ships as `internal/resolve`; there is still no transition
table or guard-predicate package under `internal/`, and no guard evaluator. The
kernel calls the guard seam this RDR owns (`resolve.go::GuardEvaluator`) but
supplies no implementation of it. Remaining consumers are the peer RDRs and the
future lint seam. RDR 0001 requires guard evaluation to feed
exact-one edge selection, RDR 0002 names positive `all` and negative `unless`
guard lists but delegates the operator grammar here, and RDR 0006 will depend on
this RDR for static determinism and exhaustiveness checks.

The Domain priors do not include an in-repo competitor document for guard
predicate exhaustiveness, so the bounded prior-art pass used the sibling
`../state-machines` corpus. The strongest local prior is
`MODEL-transition.md`: state is a tag-set, a transition is predicate matching
over that tag-set, and "guard" is not a side input but a predicate over machine
data. That model also names tag provenance (`owned`, `observed`,
`recognized`) and design-time lint failures for overlapping rows,
non-exhaustive predicates, and owned-tag read-before-write.

The tool corpus sharpens the options. `transitions` has the closest authoring
shape to RDR 0002: positive `conditions` and negative `unless` lists. enetx/fsm,
stateless, and qmuntal-stateless prove the runtime pattern, but their guards are
host callbacks, so they are useful implementation priors and poor static-lint
contracts. qmuntal-stateless explicitly requires same-state guards to be
mutually exclusive, which matches this RDR's target property but leaves the proof
to humans or runtime behavior. SCXML/scxmlcc, Sismic, and StateSmith demonstrate
that extended state, guards, bounded counters, rewinds, and conditional skips are
statechart-shaped and can be specified for validation; the `../state-machines`
contrast docs are clear that intrastate should use that class as
specification/verification prior art, not as a runtime orchestrator. Statewright
is useful contrast for a small external guard DSL (`field`, `op`, `value`) and
per-state `max_iterations`, but its MCP harness model is explicitly outside this
project's "no orchestrator" lens.

Arc corpus searches over `StateMachineOS`, `StateMachineLit`, `DevRefOS`, and
`SpecDrivenDev` did not change the option set, but they strengthened the
boundary. `StateMachineOS` surfaced stateless guard-discrimination tests,
`transitions`' `Condition` helper, Sismic guard evaluation, and SCXML `cond`
fixtures. `StateMachineLit` surfaced Symbolic Guardrails and AgentSpec as the
research form of the same split: symbolic predicates are valuable where the
policy is concrete, and an explicit escape remains necessary where judgment is
not symbolically enforceable. StateFlow confirms the finite-state framing and
finite alphabets, but it is orchestration-adjacent rather than a local predicate
grammar for intrastate.

Sibling-path check for an existing guard identity or predicate signal:

```sh
rg -n "resolve|resolver|transition|state|guard|tag|predicate|recognized|outcome|next legal|illegal|refus" internal cmd docs
```

The search found no implemented guard predicate evaluator under `internal/` —
`internal/resolve` declares the `GuardEvaluator` seam and calls it, but no type
satisfies it. The adjacent design signal is RDR 0002's `all`/`unless` split and
exact-one row selection, so this RDR extends that signal instead of inventing a
parallel guard model.

### Key Discoveries

- **Documented** — RDR 0001 requires guard predicates to expose enough
  structure for deterministic single-edge selection.
- **Documented** — RDR 0002 owns sparse transition-table authoring and
  normalization, and delegates the fixed predicate operator set to this RDR.
- **Documented** — prior-art FSM libraries commonly support guards, but
  callback-style guards are only runtime predicates; they do not give lint a
  symbolic domain to prove coverage or overlap.
- **Documented** — `MODEL-transition.md` already frames guards as tag-set
  predicates and requires design-time rejection for overlapping predicates,
  non-exhaustive predicates, and owned-tag read-before-write.
- **Documented** — statechart/SCXML/Sismic prior art covers the hard workflow
  cases in grammar: extended state, bounded counters, rewinds, and conditional
  skips; the project lens keeps these as validators/specs, not as orchestrators.
- **Documented** — arc `StateMachineLit` results support the symbolic/escape
  boundary: concrete policies can be enforced by symbolic predicates, while
  ambiguous requirements still need a model or human escape path.
- **Verified** — the target RDR and kata flows only need equality, bounded
  integer comparison, enum/set membership, existence, and declared negation via
  `unless`; the Resolve spike encoded representative guard rows with that
  vocabulary.
- **Verified** — requiring finite domains for exhaustive checks is acceptable:
  bounded integer tags such as cap counters can be proven, while unbounded
  values can be evaluated but cannot carry an exhaustiveness guarantee.

### Critical Assumptions

- **A1 The target RDR and kata flows fit a closed typed predicate vocabulary.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `cd docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes && sh check.sh guard-fixture.toml` validated four representative guard rows covering status/profile routing, cap-3 handling, prelock lens sets, cluster eligibility, and rewind legality with only `eq`, `in`, `lt`, `gte`, and `exists`; transcript captured in `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`.
  - **If wrong**: The fixed operator set is too small, and authors will need an
    expression grammar or host predicates that weaken static lint.
- **A2 Every exhaustiveness claim can be reduced to scoped finite declared
  domains.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: For every exhaustiveness-eligible row group, lint receives the candidate rows that share a selection context from the normalized model plus finite domains for the guard dimensions that vary inside that group: enum/boolean values as declared sets, set-valued tags as a symbolic set over the declared element universe, and bounded integers as `{min..max}`. A row with `all` atoms denotes the intersection of each atom's allowed subset of the scoped product; its `unless` block denotes an excluded intersection that is subtracted from the row's accepted assignments. Coverage is `union(row_i accepted assignments) == scoped product` for that row group, and overlap is any non-empty `row_i accepted assignments intersect row_j accepted assignments`. If any participating dimension lacks a finite domain, or the finite product cannot be represented by the implementation's symbolic/bitset-equivalent proof, lint must refuse or downgrade the exhaustiveness claim.
  - **If wrong**: Lint may falsely claim guard coverage or miss legal gaps in
    cap/profile/lens routing.
- **A3 `all` plus `unless` is enough polarity; inline `not` operators are not
  required.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: This RDR chooses separate positive and negative guard lists: `all` is the required conjunctive predicate set, `unless` is the conjunctive exclusion set, and RDR 0002's normalized candidate-row contract combines both before ambiguity checks. Inline `not` and nested boolean expressions are rejected to keep each diagnostic tied to an authored atom and to preserve finite-domain coverage/overlap derivation.
  - **If wrong**: Authors will duplicate rows or encode confusing inverse
    predicates that make overlap diagnostics harder to understand.
- **A4 Guard predicate errors can use the existing structured CLI failure
  gateway.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/clierr/clierr.go::CLIError` carries stable `Code`, human `Message`, optional `Param`, `Detail`, `Hint`, and exit-code `Group`; `internal/cli/respond/respond.go::Fail` emits the envelope in text/json modes; `internal/cli/config/config.go::Load` already demonstrates stable parse/read error codes. This RDR owns predicate semantic kinds such as unknown operator, type mismatch, literal parse failure, and unsupported operator/tag-kind pairing; the per-row/per-atom `guard_unevaluable` payload is RDR 0007's under JDR 0001 §D4; RDR 0006 owns graph-lint finding codes for non-exhaustive and overlapping row groups; RDR 0005 owns the CLI command/envelope mapping onto the existing gateway.
  - **If wrong**: This RDR or RDR 0005 must add a separate user-facing error
    contract before implementation.
- **A5 The normalized predicate representation can retain source identity for
  actionable diagnostics.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: The normalized atom shape is fixed at the kernel seam by JDR 0001 §D1; RDR 0007 is the landing document and states it normatively — a row carries parsed atoms (key, operator token, literal, block), not an opaque predicate string, so there is no reconstruction step that could lose identity. RDR 0007 is `Final`, not yet implemented: the shipped kernel still carries `Row.Guard string`, and 0007's A3 scopes the reshape. The source identity this assumption is about already ships — `internal/resolve/resolve.go::Row` carries `RuleID` and `SourceLocator`. RDR 0002 `Normative Contracts` require each normalized candidate row to retain source rule id and source locator, and its `Validation / Testing Strategy` carries those through `all`/`unless` expansion; RDR 0006 consumes the same source rule ids/spans for graph lint findings.
  - **If wrong**: Lint may detect an error but fail to point reviewers at the
    guard to fix.
- **A6 Tag provenance is available to predicate lint.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: RDR 0002 `Normative Contracts` require the model to declare every matched or written tag, including `owned`, `observed`, or `recognized` provenance. RDR 0006 `Technical Design` consumes tag provenance and owned-tag write effects from normalized rows for owned-set-before-match and coverage checks.
  - **If wrong**: Predicate lint can still evaluate runtime truth, but it cannot
    prove that owned tags are set before they are matched.
- **A7 An existence atom projects onto the declared-domain product as a
  per-key presence dimension, so a row group carrying `exists` atoms stays
  provable.**
  - **Status**: Pending
  - **Method**: Design Decision
  - **Evidence**: Derivation recorded in
    `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/research/iter-2-projection-derivation.md`.
    RDR 0007 A12 routes the question here (`DOWNGRADED`; "it lands when 0003
    states its existence-atom projection"), and its Phase 4 hands it over with
    six sibling items. The row-group branch of A12's disjunction is closed:
    RDR 0007's two-row absence pattern answers one source-state/outcome pair,
    so both rows share one selection group and scoping cannot separate them.
    RDR 0007 fixes `exists` as the sole TOTAL operator (`presence == literal`),
    and A2 above requires every atom to denote a subset of the scoped product;
    an `exists` atom denotes `{present}` or `{absent}`, which are subsets of a
    presence dimension. Without that dimension the atom denotes nothing and
    A2's coverage identity drops a live constraint. The bounded prior-art pass
    is a recorded negative — no decision-table/DMN/rule-base verification
    material in `PapersFast`, `StateMachineLit`, or `StateMachineRes` — so this
    resolves as a Design Decision, never as Prior Art.
  - **Plan**: state the projection rule as a normative clause before lock, and
    carry the RDR 0002 producer request it depends on (below). Blocked on that
    request: RDR 0002's tag declaration carries name, provenance, value kind,
    and optional accessor reference, with no optionality field, so no producer
    can today declare which keys may be absent. RDR 0002 is itself `Draft`.
  - **If wrong**: lint either drops `exists` atoms from the product — certifying
    a row group exhaustive that refuses `guard_unevaluable` at runtime, which
    breaches the §JD-4 narrowing this RDR just adopted — or refuses every group
    containing an `exists` atom, which disqualifies this RDR's own representative
    fixture row `foundational-to-cove` and leaves RDR 0007's sanctioned two-row
    absence pattern permanently unprovable, pushing authors back to the
    sentinel-stamping anti-pattern.
- **A8 RDR 0006 can carry the same `guard_unevaluable` narrowing this RDR
  states, so the two documents agree on one lint promise.**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence needed**: JDR 0001 §JD-4 decides the substance ("the *promise*
    narrows — P5 decides") but is an open entry as to which document records
    it, and its head clause names RDR 0006's proof, not this RDR's. RDR 0007
    reads it the same way: its Predecessors block names "**RDR 0006** (`Final`,
    tolerance §JD-4) narrows lint's promise", and RDR 0006 carries
    `Final [joint decision → JDR 0001 §JD-4]` on its Status line. But RDR
    0006's exhaustiveness contract narrows only for a non-finite dimension
    ("If a required dimension is not finite, lint MUST emit a blocking
    inability-to-prove finding"); it says nothing about a fully-finite product
    whose participating row can still refuse `guard_unevaluable`. This RDR now
    states that case. Verification is a route-back to RDR 0006 confirming it
    adopts the same wording, or a §JD-4 disposition assigning the recording to
    one document.
  - **Plan**: raise the divergence at cluster reconcile before either document
    locks; RDR 0006 is `Final`, so closing this needs a route-back, not a
    silent edit here.
  - **If wrong**: the two documents state one lint promise two ways — the
    single-source failure §JD-4 exists to prevent — and an implementer reading
    only RDR 0006 certifies a row group exhaustive that this RDR forbids.

**Method vocabulary** (pick exactly one per assumption):

- **Source Search** — verified against dependency
  source code. Evidence: a greppable `path::Symbol`
  (function/type/const name), **not a bare `file:line`**;
  a commit-SHA permalink only for audit/traceability.
  Standard for libraries. (Why symbol not line: flow
  README *Doctrine*.)
- **Spike** — verified by running code against a live
  service or fixture. Evidence: command run + path to
  captured output.
- **Prior Art** — same property holds in ≥1 named
  external system. Evidence: system + section/page.
- **Derivation** — pure math or proof. Evidence: the
  derivation, shown inline.
- **Design Decision** — a scoping choice this RDR is
  *making* (not *verifying*). Evidence: the decision
  and the alternative explicitly rejected.
- **Peer RDR** — relies on a property defined in
  another RDR. Evidence: RDR ID + section.
- **MVV Test** — the property is testable via the
  Minimum Viable Validation, and the test
  is named in this RDR's Validation section (pending
  implementation at lock time). Evidence: test name.
- **Docs Only** — documentation reading alone.
  **Insufficient** for load-bearing assumptions; allowed
  only when paired with a Spike or Source Search plan
  in the Evidence line.

A `Method: Source Search` whose Evidence cites this
same RDR file — or any path under the RDR's artifact
directory — is self-reference and not Verified. The
cited proof must also support **the specific claim**,
not an adjacent one: confirming a neighboring fact and
stamping the assumption `Verified` is not verification.
The cited symbol must resolve on `main` (a renamed,
deleted, or never-built symbol fails the check).

Any exactness claim such as all/every, first/nearest,
byte-identical, lossless, canonical, deterministic, or
stable order must be covered by a Critical Assumption
Evidence Record or by the Minimum Viable Validation.

## Proposed Solution

### Approach

Use a closed, symbolic predicate-atom model over declared tags. This is the
middle path between table-as-opaque-callback and full statechart adoption: keep
RDR 0002's sparse TOML table as the source, but make every guard a symbolic
tag-set predicate that lint can reason about. Authors keep the RDR 0002 shape of
positive `all` guards and negative `unless` guards; each guard entry is a typed
atom: tag reference, fixed operator, and typed literal or literal set. The
operator set is deliberately small: equality, membership, bounded integer
comparison, existence, and set containment. There is no embedded host predicate
and no free-form expression grammar.

Enforcement of the guard domain is split, per JDR 0001 §D4 (normative in RDR
0007, its landing document): the kernel decides key presence, decides existence atoms from presence
alone, marks a value atom over an absent key unevaluable without consulting
this RDR's evaluator, and combines per-atom verdicts. This RDR's evaluator owns
**value semantics over a present value** — nothing more; it never sees the tag
view.

Each tag declaration supplies the value kind and, when lint must prove
exhaustiveness, the finite domain: enum values, boolean values, set element
universe, or bounded integer range. Runtime guard evaluation is simple predicate
evaluation over the assembled tag-set. Static lint reasons over the same atoms
inside row groups supplied by the normalized model: within one selection
context, it checks the scoped product of participating guard domains, proves
whether candidate rows are mutually exclusive, and proves whether they cover the
declared product.

The initial operator/kind matrix is:

| Operator | Accepted tag value kinds | Literal shape | Lint proof role |
| --- | --- | --- | --- |
| `eq` | enum, boolean, integer, string-like scalar | one typed scalar | Narrows the tag domain to one value. |
| `in` | enum, boolean, integer, string-like scalar | non-empty typed scalar set | Narrows the tag domain to the listed values. |
| `lt`, `lte`, `gt`, `gte` | integer | one typed integer | Narrows a bounded integer domain by comparison; remains runtime-only for an unbounded integer. |
| `exists` | optional scalar or optional set-valued tag | boolean | Tests presence or absence, not value equality. Decided by the kernel from presence alone (JDR 0001 §D4); it never reaches this RDR's evaluator. |
| `contains` | set-valued tag with a declared element universe | non-empty typed element set | Narrows the set-valued domain to assignments containing every listed element. |

This intentionally aligns with prior art that separates positive and negative
guard lists, and deliberately diverges from callback-driven FSM libraries.
Callbacks are ergonomic for application code, but intrastate needs the guard
surface to be reviewable, serializable, and mechanically analyzable.

### Technical Design

Guard predicates are part of the normalized transition model, not a separate
runtime plug-in system. RDR 0002 owns the sparse table source, tag declarations,
provenance, and normalization; this RDR owns the guard predicate grammar and
semantics that normalized rows carry. RDR 0001 consumes evaluated predicates for
exact-one edge selection, and RDR 0006 consumes the symbolic predicate
constraints for graph lint.

A predicate atom has four conceptual fields: tag name, operator, expected value,
and source identity. The tag name must resolve to a declared tag. The operator
must be allowed by the operator/kind matrix above. The expected value must parse
to the operator's literal shape. Source identity is inherited from the RDR 0002
row/context so diagnostics can point back to the authored guard. The atom shape
itself is fixed at the kernel seam by JDR 0001 §D1 and stated normatively in RDR
0007; this RDR cites it rather than restating it.

Because the kernel decides presence, a value atom over an absent key is
unevaluable rather than false, and an existence atom over an absent key is
decided — not unevaluable. This RDR's requirement that every guard resolve to a
declared tag is therefore read against `exists = false`: an author expressing
"this row applies when tag X is absent" writes an existence atom, not a value
atom that happens to miss.

Positive `all` predicates are conjunctive requirements. Negative `unless`
predicates are also conjunctive within the excluded predicate set: if all
`unless` predicates hold, the candidate row is disabled. Normalization may
combine them into one internal constraint object, but authoring keeps the
separation because it reads better and mirrors prior art.
For lint, RDR 0002 supplies the normalized candidate rows and RDR 0006 supplies
the graph-lint grouping context, such as one source-state/recognized-outcome
selection group. Within that group, a row's accepted assignments are the
intersection of all positive `all` atom domains minus the single conjunctive
assignment set matched by the row's full `unless` block. `unless` is not
per-atom negation, and it does not create source-order priority.

Exhaustiveness is only asserted over finite declared domains. Enum, boolean,
and set-element universes are finite by declaration; set-valued tags must be
represented symbolically or as an implementation-equivalent bitset rather than
requiring literal powerset materialization. Integer tags are exhaustive only
when they declare a bounded range; this covers cap counters and iteration limits
without pretending arbitrary integers can be fully partitioned. A guard may
still compare an unbounded integer at runtime, but lint must report that it
cannot prove exhaustive coverage for that dimension. If a finite scoped product
is too large for the implementation to prove deterministically, lint must also
refuse or downgrade the exhaustiveness claim instead of silently capping the
proof.

Where this exhaustiveness proof and RDR 0007's aggregation veto disagree, the
**promise narrows**: a green exhaustiveness result must mean resolution
succeeds. Lint therefore cannot certify a row group exhaustive when any
participating row could refuse `guard_unevaluable` at runtime; the runtime veto
is not weakened to make lint's claim true.

JDR 0001 §JD-4 decides that substance ("the *promise* narrows — P5 decides")
but leaves open **which document records it**, and its head clause names RDR
0006's proof rather than this RDR's. This RDR records the narrowing for the
proof it defines, because the constraint binds what its own finite-domain
product may claim. RDR 0006 is `Final` under the §JD-4 tolerance and its
exhaustiveness clause covers only the non-finite-dimension refusal, so it does
not yet carry this case; A8 tracks the agreement the two documents owe each
other. Provenance affects lint: recognized tags are fresh event inputs, observed
tags are re-read before matching, and owned tags must have a reachable
predecessor write before a row may match them.

#### Normative Contracts

```normative
A guard predicate MUST be a symbolic atom over a declared tag, not a host
language callback and not a free-form expression string.
```

```normative
The initial operator vocabulary MUST be closed and typed: equality,
membership, bounded integer comparison, existence, and set containment. Unknown
operators MUST be rejected during parse or lint before resolution.
```

```normative
Each operator MUST declare which tag value kinds it accepts. A predicate whose
literal cannot be parsed as the declared tag kind MUST be rejected before
resolution.
```

```normative
This RDR's guard evaluator MUST decide value semantics over a present value
only. Key presence, existence atoms, absent-key unevaluability, and the
combination of per-atom verdicts belong to the kernel (JDR 0001 §D4, normative
in RDR 0007); the evaluator MUST NOT read the tag view.
```

```normative
Positive guard atoms MUST live in `all`; negative guard atoms MUST live in
`unless`. Successful row matching MUST NOT depend on source order or
first-match priority.
```

```normative
Lint MAY claim guard exhaustiveness only for finite declared domains: enum
values, booleans, declared set element universes, or bounded integer ranges.
```

```normative
If a guard dimension lacks a finite declared domain, lint MUST refuse or
downgrade an exhaustiveness claim for that dimension rather than treating the
covered examples as complete.
```

```normative
Coverage and overlap checks MUST be scoped to a normalized row group supplied by
the transition/lint model, and MUST evaluate the participating guard dimensions
as one product rather than as independent one-dimensional checks.
```

```normative
An exhaustiveness claim MUST NOT be stronger than the runtime it describes: lint
MUST NOT certify a row group exhaustive when a participating row can refuse
`guard_unevaluable` under RDR 0007's aggregation veto. Where the two disagree
the lint promise narrows; the runtime veto MUST NOT be weakened. RDR 0006, which
holds the JDR 0001 §JD-4 tolerance, MUST carry this same narrowing for the
findings it emits; the two documents MUST NOT state it differently.
```

```normative
Set-valued guard domains MUST be proved with a deterministic symbolic or
bitset-equivalent representation. If the finite product is too large for that
proof, lint MUST refuse or downgrade the exhaustiveness claim rather than
silently capping enumeration.
```

```normative
Overlap and coverage diagnostics MUST name the source rule id or context id
that contributed each predicate involved in the finding.
```

```normative
Predicate lint MUST distinguish owned, observed, and recognized tags. A row
that matches an owned tag MUST be rejected unless every reachable predecessor
sets or preserves that tag before the match.
```

#### `authority` — source-authority census

Cue: the exhaustiveness verdict has more than one candidate source of truth
(this RDR and RDR 0006), and guard-domain enforcement is split across the kernel
and this RDR's evaluator.

| Input / decision | Writer (canonical) | Readers | Sibling arms | Which is canonical |
| --- | --- | --- | --- | --- |
| Key presence | Kernel (RDR 0007, JDR 0001 §D4) | This RDR's lint; RDR 0001 selection | This RDR's evaluator — **explicitly not an arm**; it never reads the tag view | Kernel |
| Existence-atom verdict | Kernel, from presence alone (`presence == literal`) | Lint's presence dimension (A7) | Value atoms over the same key (unevaluable on absence) | Kernel |
| Value-atom verdict over a present value | This RDR's evaluator | Kernel combinator | None | This RDR |
| Per-atom verdict combination | Kernel (strong Kleene) | RDR 0001 exact-one selection | None | Kernel |
| Exhaustiveness verdict for a row group | **This RDR** (finite-domain product proof) | RDR 0006 lint findings; RDR 0005 envelope | RDR 0006's own exhaustiveness clause — non-finite-dimension case only | **Contested — A8.** This RDR states the `guard_unevaluable` narrowing; RDR 0006 holds the §JD-4 tolerance but does not carry it. Must converge before lock. |
| Guard atom shape | Kernel seam (JDR 0001 §D1, normative in RDR 0007) | This RDR's grammar; RDR 0002's normalizer | Shipped `Row.Guard string` — the pre-reshape form | Kernel seam (specified); the shipped string form is superseded, not an arm |
| Source identity for diagnostics | RDR 0002 normalization (`RuleID`, `SourceLocator`) | This RDR's diagnostics; RDR 0006 findings | None | RDR 0002 |
| Per-tag optionality (which keys may be absent) | RDR 0002 — **not yet declared** | A7's presence dimension | None | RDR 0002 (requested; A7 blocked on it) |

#### `disposition` — input class × outcome

Cue: the draft sets outcomes over input classes (refuse, downgrade, prune,
unevaluable).

| Input class | Runtime outcome | Lint outcome | Diagnostic minted | Silent or loud |
| --- | --- | --- | --- | --- |
| Every `all` atom decides true; `unless` block not fully true | Row qualifies | Row contributes its accepted assignments | none | — |
| `all` atom decides false | Row pruned | Row contributes nothing | none | Silent by design — a decided false is not a defect |
| Full `unless` block decides true | Row disabled | Excluded intersection subtracted | none | Silent by design |
| Value atom over an absent key | `guard_unevaluable` refusal (RDR 0007 veto) | Exhaustiveness claim **withheld** for that group | Per-row/per-atom payload — RDR 0007's (key, block, reason `absent`) | Loud |
| Existence atom over an absent key | Decided (`presence == literal`) — never unevaluable | Selects `{absent}` on the presence dimension (A7) | none | — |
| Zero rows qualify | RDR 0001 refuses | Coverage gap if the product is provable | `graph-coverage-gap` (RDR 0006) | Loud |
| Two or more rows qualify | RDR 0001 refuses — never first-match | Overlap finding | Overlap code (RDR 0006), naming both source rule ids | Loud |
| Guard dimension lacks a finite declared domain | Evaluable at runtime | **Refuse or downgrade** the claim — never treat examples as complete | Inability-to-prove finding | Loud |
| Finite product too large to prove deterministically | Evaluable at runtime | **Refuse or downgrade** — never silently cap enumeration | Inability-to-prove finding | Loud |
| Unknown operator / unknown tag / operator–kind mismatch / literal parse failure | Rejected before resolution | Rejected at load | Predicate semantic kind (this RDR) → RDR 0006 finding → RDR 0005 envelope | Loud |
| Row group is domain-exhaustive but a participating row can refuse | Refusal stands | Claim withheld — the narrowing | RDR 0007 payload; no green certification | Loud |

#### `trace` — desk trace over the MVV

Cue: ten normative clauses bear on one output surface — the exhaustiveness
verdict for a row group. Walked stepwise against the MVV, with witnesses from
`evidence/spikes/guard-fixture.toml`.

| Step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1. Load the RDR/kata slice as normalized candidate rows | atoms-not-callbacks; closed typed operator vocabulary; operator declares accepted kinds | `profile-to-grounding` parses to `profile in [mid,large]` + `unless prelock_iterations gte 3`; operators all in the matrix | OK |
| 2. Group rows by selection context | coverage/overlap scoped to a normalized row group, evaluated as one product | `profile-to-grounding` and `foundational-to-cove` share `match status eq Draft` → one group | OK |
| 3. Build the scoped product from declared domains | exhaustiveness only over finite declared domains | `profile` = 4 enum values; `prelock_iterations` = `{0..3}`; both finite | OK |
| 4. Project each atom onto the product | every atom denotes a subset of the scoped product (A2) | `profile eq "foundational"` → `{foundational}`; `profile in [mid,large]` → `{mid,large}` | OK |
| 5. Project the `exists` atom | A7 presence dimension — **Pending** | `foundational-to-cove` carries `cluster_eligible exists = true`; `cluster_eligible`'s declared domain `{true,false}` is complete, so no element selects presence | **GAP — booked as A7, not a contradiction.** Without the presence dimension this step has no defined result. A7 is `Pending` with a plan; the draft does not claim the step succeeds. |
| 6. Compute coverage | `union(row_i accepted) == scoped product` | Step 5 unresolved for this group; groups with no `exists` atom compute normally (`continue-prelock-lenses`, `reconcile-rewind-legality`) | OK for `exists`-free groups; blocked on A7 otherwise |
| 7. Compute overlap | any non-empty pairwise intersection; no source-order priority | `profile in [mid,large]` ∩ `profile eq foundational` = ∅ → rows disjoint on that dimension | OK |
| 8. Apply the runtime-veto narrowing | claim MUST NOT be stronger than the runtime; withhold if a participating row can refuse | MVV Scenario 6's row group: domain-exhaustive, one value atom over a possibly-absent key → claim withheld | OK — and the reason A7's presence dimension must not silently drop `exists` atoms |
| 9. Emit the verdict | diagnostics name the contributing source rule/context id | `RuleID` + `SourceLocator` ship on `internal/resolve/resolve.go::Row` | OK |
| 10. Cross-document agreement | the two documents MUST NOT state the narrowing differently | RDR 0006 carries no `guard_unevaluable` narrowing | **GAP — booked as A8.** Not a contradiction inside this draft; a divergence with a `Final` sibling that needs a route-back. |

No CONTRADICTION row. The two GAP rows are the draft's two `Pending`
assumptions, each carrying a named plan; neither is asserted as already true.

#### Load-Bearing Decisions

- **Identity** — a guard predicate is identified by its source rule/context id
  plus its position within `all` or `unless`; semantic equality is the normalized
  tuple `(tag, operator, literal)`.
- **Wire / byte format** — RDR 0002 owns the TOML container; this RDR owns the
  guard atom grammar embedded in that container.
- **Naming** — the canonical name is "guard predicate"; rejected alternatives
  are "condition callback" and "guard expression" because both invite opaque
  host logic.
- **Operator semantics** — equality compares a tag value to one typed literal;
  membership checks a scalar tag against a typed literal set; bounded integer
  comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks presence of an
  optional tag value; set containment checks declared set-valued tags against a
  typed element set.
- **Finite-domain proof** — exhaustive coverage is a lint claim over the
  declared tag-domain product for a normalized row group, not over examples
  observed in fixtures. Unbounded dimensions and finite products that cannot be
  represented deterministically remain runtime-evaluable but cannot satisfy an
  exhaustiveness proof.
- **Error ownership** — this RDR names predicate semantic kinds over present
  values; RDR 0007 owns the per-row/per-atom `guard_unevaluable` payload (JDR
  0001 §D4), RDR 0006 maps graph-level predicate failures to lint finding codes,
  and RDR 0005 maps those findings/refusals to the CLI envelope.
- **Selection / predicate** — a row qualifies only when every `all` atom is true
  and the `unless` predicate set is not fully true; if multiple rows qualify,
  RDR 0001's exact-one resolver refuses instead of choosing by priority.

#### Round-Trip / Inverse Invariants

This RDR introduces no encode/decode pair. Parse/render fidelity for the sparse
TOML source belongs to RDR 0002.

#### Illustrative Code

Illustrative predicate shape only:

```toml
[rule.guard.all]
status.eq = "Draft"
profile.in = ["mid", "large", "foundational"]

[rule.guard.unless]
prelock_iterations.gte = 3
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Sparse transition-table container | RDR 0002 | Pending | This RDR assumes guards live inside RDR 0002's `all`/`unless` blocks. |
| Deterministic exact-one resolver | RDR 0001 | Pending | Runtime selection refuses zero or multiple matching rows. |
| Guard predicate grammar and finite-domain semantics | This RDR | Introduced | Lint and runtime share one symbolic predicate model. |
| Accessor read/write safety | RDR 0004 | Pending | Guard evaluation consumes tag values after accessor binding; it does not execute accessors. |
| Graph lint authority | RDR 0006 | Pending | Exhaustiveness and overlap findings become blocking lint there. |
| Kernel guard-domain enforcement and `guard_unevaluable` payload | RDR 0007 | Pending | The kernel decides presence and existence atoms; this RDR's evaluator narrows to value semantics over a present value. |
| Per-tag optionality declaration (which keys may be absent) | RDR 0002 | Requested | A7's presence dimension needs it; RDR 0002's tag declaration carries name, provenance, value kind, and accessor reference only. Requested alongside the set-valued element encoding RDR 0007 Phase 4 routes here. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| User-facing failures | `internal/cli/clierr`, `internal/cli/respond` | Must verify error-code coverage | Reuse | Predicate parse/lint failures should use existing CLI envelopes. |
| Guard evaluator | None found under `internal/` | New semantic surface | Introduce | This RDR owns the evaluator semantics, not command I/O. |
| Transition model source | Pending RDR 0002 | Not implemented | Extend peer | Guard atoms are embedded in table rows/contexts. |

### Decision Rationale

Joint-check: fired → RDR 0007 is the normative home of the guard seam (JDR 0001
§D1, §D4). It fixes the parsed-atom row shape, the value-only per-atom
`GuardEvaluator` seam, the kernel-exported existence operator token and boolean
literal forms RDR 0002's normalizer emits, exact key identity at the kernel, and
the per-atom `guard_unevaluable` payload. This RDR cites those rather than
restating them. It also records the §JD-4 narrowing for its own proof — §JD-4
decides the substance but leaves the recording document open, and names RDR
0006's proof in its head clause; A8 carries the resulting sibling obligation.

The fixed symbolic atom model best matches the user's outcome: flow authors can
write conditional edges, and lint can still prove whether those edges are
complete and mutually exclusive. It preserves RDR 0002's readable `all`/`unless`
authoring form while giving RDR 0006 finite-domain constraints to analyze. It
also follows the `MODEL-transition.md` thesis directly: a guard is a predicate
over tags, not a second decision channel.

This choice aligns with `transitions` on positive/negative guard factoring,
with statechart/SCXML/Sismic on treating counters and rewinds as extended-state
guarded transitions, and with Statewright-style guard objects on keeping
predicate data small. It deliberately diverges from those systems where they
would undermine the local contract: callback guards are not statically
exhaustive, a full statechart runtime would replace RDR 0002's sparse table
source, and an external guardrail harness would violate the no-orchestrator
project lens.

The chosen approach survived the premortem. If it shipped and failed, the most
likely failure would be an operator set too small for a real RDR/kata guard,
causing authors to request escape hatches. The mitigation is to verify the
target-flow fixtures in Resolve before lock and add only concrete operators
proven necessary there. The recommendation still holds because the alternatives
that solve expressiveness up front lose the static proof this RDR exists to
protect.

## Alternatives Considered

### Alternative 1: Fixed Typed Predicate Atoms Inside RDR 0002

**Description**: Represent every guard as a symbolic atom over declared tags,
with a closed operator set, tag provenance, and finite-domain declarations for
exhaustive lint. Author predicates in RDR 0002's positive `all` and negative
`unless` blocks.

**Prior-art alignment**: `MODEL-transition.md` defines transition resolution as
tag-set predicate matching. `transitions` supplies the `conditions`/`unless`
factoring. Statewright's guard objects show the useful small shape (`field`,
`op`, `value`) while the project rejects its external harness model.

**Pros**:

- Keeps guards reviewable and serializable as data.
- Lets lint reason about overlap, coverage, and owned-tag read-before-write
  without executing user code.
- Covers the seed's equality, integer comparison, and membership needs.
- Aligns with RDR 0002's sparse table shape and RDR 0001's exact-one resolver.
- Leaves SCXML/statechart tools available as external validators rather than
  replacing the local source format.

**Cons**:

- Requires finite-domain declarations for strong exhaustiveness claims.
- May need later operator additions if Resolve finds a target-flow gap.
- Requires RDR 0002 to carry tag provenance and source identity through
  normalization.

**Reason for selection**: This is the smallest model that can express the known
guards while preserving static verification and the no-orchestrator project
lens.

### Alternative 2: Table-Level Tag Predicates Without Separate Guard Syntax

**Description**: Remove the conceptual distinction between "match" and "guard":
every transition row is just a set of tag predicates on the left-hand side and
tag writes on the right-hand side. `all`/`unless` becomes documentation or
syntactic sugar over row predicates.

**Prior-art alignment**: This is the purest reading of `MODEL-transition.md`:
there are no separate guards, only predicates on tags.

**Pros**:

- Conceptually clean: one predicate language for all row matching.
- Makes lint rules uniform for outcome selection, profile routing, cap
  counters, and rewind legality.
- Avoids a second nested guard namespace in the table format.

**Cons**:

- Conflicts with RDR 0002's existing `all`/`unless` authoring split.
- Loses the readability of positive requirements versus negative exclusions.
- Makes it harder to align with `transitions` prior art and with the way users
  usually discuss guards.

**Reason for rejection**: This is semantically attractive but too disruptive to
the neighboring RDR. The chosen approach keeps `all`/`unless` as the authoring
surface while normalizing to tag predicates internally.

### Alternative 3: Tiny Boolean Expression Grammar

**Description**: Put a small expression language in guard fields, such as
`status == "Draft" && profile in ["mid", "large"] && iterations < 3`.

**Prior-art alignment**: SCXML `cond` expressions, XState parameterized guards,
and Statewright guard objects all demonstrate guard expressions or structured
guard references, but they do not all provide static exhaustiveness over the
local tag domains.

**Pros**:

- Compact for authors who already know expression syntax.
- Can express nested boolean logic without row factoring.
- Familiar from SCXML/statechart tools.

**Cons**:

- Requires parser, precedence, type-checking, formatting, and diagnostics that
  are larger than the known problem.
- Makes proof/debug output harder to map back to individual authored atoms.
- Encourages clever predicates instead of simple reviewable guard rows.
- Reintroduces much of SCXML's expression machinery without adopting SCXML's
  full validation model.

**Reason for rejection**: The extra expressiveness is not needed for the target
flows and increases the surface that lint must prove.

### Alternative 4: Embedded Host Predicates / FSM Callback Guards

**Description**: Let table rows reference Go functions or callbacks that return
true or false at runtime.

**Prior-art alignment**: enetx/fsm, stateless/qmuntal-stateless, looplab-fsm,
and XState all prove this as a runtime pattern. qmuntal-stateless even states
same-state guards must be mutually exclusive.

**Pros**:

- Maximum expressiveness.
- Matches many FSM libraries' guard pattern.
- Easy to add without designing a predicate grammar.
- Can reuse Go types and helper functions directly.

**Cons**:

- Static lint cannot prove exhaustiveness or mutual exclusion without executing
  arbitrary code.
- Reviewers must inspect scattered Go functions instead of one transition
  artifact.
- Replay determinism depends on callback purity and hidden dependencies.
- The prior art's mutual-exclusion requirement becomes a runtime convention, not
  a design-time proof.

**Reason for rejection**: It gives up the static verification requirement in
the problem statement.

### Alternative 5: Full Statechart/SCXML/Sismic Guard Semantics

**Description**: Adopt a statechart-style guard language or runtime/spec with
hierarchy, datamodel expressions, Design-by-Contract invariants, and conditional
transitions.

**Prior-art alignment**: The `../state-machines` compiler eval found scxmlcc
and Sismic are strong at expressing the RDR hard parts in grammar: cap counters,
rewinds, coupled status, and conditional skips.

**Pros**:

- Strong prior art for guarded transitions and graph tooling.
- Can model hierarchy and extended state directly.
- Could validate the legal graph outside intrastate's own implementation.
- Captures bounded counters and rewind edges more naturally than a flat FSM.

**Cons**:

- Imports a larger conceptual model than intrastate needs.
- Risks making intrastate a workflow engine instead of a thin deterministic
  resolver/lint CLI.
- Still needs local restrictions to make guard predicates statically provable.
- Adds another source-of-truth format beside RDR 0002's sparse TOML model.

**Reason for rejection**: RDR 0002 already chooses a sparse local table format;
this RDR should define the predicate seam inside that format, not replace the
format with a statechart runtime. SCXML/Sismic remain good validator or
cross-check prior art for RDR 0006.

### Alternative 6: External Guardrail DSL / Harness Model

**Description**: Move guard decisions into a workflow harness that constrains
tools, commands, state transitions, and approvals around the agent.

**Prior-art alignment**: Statewright exposes conditional transitions as
structured guards over context data and supports per-state max iterations.

**Pros**:

- Strong for external enforcement of phase behavior and tool access.
- Makes guard conditions explicit and inspectable.
- Already includes workflow-level controls such as approval gates and iteration
  caps.

**Cons**:

- It is an orchestrator/harness, which the `../state-machines` lens explicitly
  rejects for intrastate.
- Moves the source of truth outside the RDR transition model.
- Solves agent/tool enforcement, not the internal guard predicate grammar needed
  by RDR 0001/0002/0006.

**Reason for rejection**: Useful contrast only. Intrastate needs a local
predicate contract inside its resolver/lint data, not a harness wrapped around
skills.

### Briefly Rejected

- **First-match priority guards**: Rejected because source order would become
  behavior, conflicting with exact-one selection and review-safe reordering.
- **Regular-language FSM compilers such as Ragel/re2c**: Rejected by the sibling
  POC as the wrong category; they recognize byte/token streams and push guards,
  counters, and rewinds into host code.

## Trade-offs

### Consequences

- Guard authoring stays data-first and reviewable.
- Static lint can produce gap and overlap diagnostics from the same predicate
  model runtime resolution uses.
- Unbounded integer or free-form string dimensions cannot receive silent
  exhaustiveness claims; authors must declare finite domains or accept a lint
  limitation.

### Risks and Mitigations

- **Risk**: The closed operator set misses a real target-flow guard.
  **Mitigation**: Resolve must encode representative RDR and kata guards before
  lock; add only operators proven necessary by that fixture.
- **Risk**: Finite-domain declarations feel like boilerplate.
  **Mitigation**: Keep declarations close to tag definitions in RDR 0002's
  model and make diagnostics explain when a missing domain blocks proof.
- **Risk**: Multi-dimensional coverage is accidentally checked one dimension at
  a time.
  **Mitigation**: Require MVV cases with a gap and overlap that only appear in a
  scoped row-group product.
- **Risk**: Set-valued domains are implemented by naive powerset enumeration.
  **Mitigation**: Require symbolic or bitset-equivalent proof and explicit
  refusal/downgrade for finite products that are too large to prove.
- **Risk**: `unless` semantics are misunderstood as per-atom negation.
  **Mitigation**: Normatively define it as an excluded predicate set and require
  fixtures that demonstrate mixed `all`/`unless` behavior.

### Failure Modes

Visible failures should be typed parse or lint failures: unknown operator,
operator/tag-kind mismatch, literal parse failure, unknown tag, non-exhaustive
finite domain, overlapping candidate rows, guard dimension not provable because
it lacks a finite domain, or finite scoped product too large to prove. A guard
that cannot be decided at runtime surfaces as RDR 0007's `guard_unevaluable`
refusal, not as a predicate parse kind owned here. Silent
failure would be a false exhaustiveness claim; the recovery path is to keep
every exactness claim tied to A2 and the MVV fixture before Final.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified — A1-A6 hold (A5 re-verified against the
  kernel-fixed atom shape, JDR 0001 §D1); **A7 is Pending**, blocked on RDR
  0002's optionality declaration; **A8 is Pending**, blocked on a route-back to
  `Final` RDR 0006 to agree one lint promise
- [x] RDR 0002's sparse table/container contract is stable enough to host guard
  atoms.
- [x] RDR 0007 is the normative home of the guard seam, the domain rule, and the
  `guard_unevaluable` payload this RDR cites. RDR 0007 is `Final`; its kernel
  reshape is specified, not yet implemented, so this RDR's implementation
  sequences after it.

### Minimum Viable Validation

Encode one RDR flow slice and one kata flow slice as normalized candidate rows,
including equality, enum membership, set containment, bounded integer
comparison, and mixed `all`/`unless` guards. Lint must prove one exhaustive and
mutually exclusive scoped row group, then detect one intentional gap and one
intentional overlap that only appear in the multi-dimensional product, with
source rule/context ids in the diagnostic. The fixture must also include one row
group that is domain-exhaustive yet contains a possibly-absent guard key, and
assert that lint withholds the exhaustiveness claim there rather than
contradicting the runtime refusal.

### Phase 1: Predicate Model

Define the tag-kind/operator compatibility matrix and normalized predicate atom
shape used by resolver and lint.

### Phase 2: Finite-Domain Lint Semantics

Define how enum, boolean, set-universe, and bounded-int domains are converted
into scoped row-group coverage and overlap checks, including
refusal/downgrade behavior for unbounded dimensions and finite products too
large to prove deterministically.

### Phase 3: Target-Flow Fixture

Build the MVV fixture against representative RDR and kata guards and use the
result to confirm or adjust the initial operator set. The Resolve spike already
covered the target-flow subset `eq`, `in`, `lt`, `gte`, and `exists`; the
implementation MVV must add at least one `contains` predicate over a declared
set-valued tag before the full operator vocabulary is accepted.

### Phase 4: Integration With Peer RDRs

Connect predicate diagnostics to RDR 0002 source identities, RDR 0001 exact-one
selection, RDR 0005 CLI output, and RDR 0006 lint authority.

### Day 2 Operations

This RDR creates no persistent resource. Day-2 management belongs to the
transition model artifact owned by RDR 0002.

### New Dependencies

No new third-party dependency is proposed. The predicate grammar and finite
domain checks should be implemented with local Go code unless Resolve proves a
small parsing or set library is necessary.

## Validation

### Testing Strategy

The MVV should become production tests that exercise both runtime predicate
evaluation and lint-time finite-domain reasoning. Done means the same normalized
predicate atoms drive exact-one row selection, overlap detection, and
exhaustiveness proof/refusal without host callbacks or source-order priority.
Resolve evidence for the representative authoring shape lives in
`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/`: the fixture
covers profile routing, cap-3 handling, prelock lens sets, cluster eligibility,
and rewind legality with the target-flow subset `eq`, `in`, `lt`, `gte`, and
`exists`. The implementation MVV must extend that coverage with `contains`
before the full closed operator vocabulary is accepted.

1. **Scenario**: Evaluate representative RDR and kata rows that use equality,
   membership, set containment, bounded integer comparison, existence, and mixed
   `all`/`unless` guards.
   **Expected**: Exactly one qualifying row resolves for the legal input; a row
   whose `all` predicates match is disabled when its full conjunctive `unless`
   block also matches; zero or multiple qualifying rows become typed refusals.
2. **Scenario**: Lint finite enum, boolean, set-universe, and bounded-int tag
   domains inside one normalized row group with one complete partition, one
   intentional gap visible only in the multi-dimensional product, and one
   intentional overlap visible only in the multi-dimensional product.
   **Expected**: Complete partitions pass; product-level gaps and overlaps fail
   with source rule/context ids.
3. **Scenario**: Lint an otherwise valid guard over an unbounded integer or
   undeclared finite domain, and lint a declared finite product too large for
   the deterministic proof representation.
   **Expected**: Runtime evaluation remains available, but lint refuses or
   downgrades the exhaustiveness claim for that dimension/product.
4. **Scenario**: Parse malformed guard atoms: unknown tag, unknown operator,
   unsupported operator/tag-kind pair, and literal parse mismatch.
   **Expected**: Each failure is rejected before resolution with a predicate
   semantic kind that RDR 0006 can map to a lint finding and RDR 0005 can map to
   the structured CLI gateway.
5. **Scenario**: Reorder authored rows and guard atoms without changing their
   semantics.
   **Expected**: Successful matching and lint findings are unchanged because
   source order is not a selection mechanism; an ambiguous pair remains a
   multiple-match refusal instead of becoming a first-match success.
6. **Scenario**: Evaluate a row group whose guards are exhaustive over their
   declared domains but where one participating row carries a value atom over a
   key that can be absent.
   **Expected**: The value atom is unevaluable rather than false, resolution
   refuses `guard_unevaluable` under RDR 0007's veto, and lint does not certify
   that row group exhaustive — the green-lint-implies-resolution-succeeds
   promise holds. An existence atom over the same absent key decides instead of
   refusing, and this RDR's evaluator is never consulted for either.

### Performance Expectations

Resolve measured representative fixture size rather than setting a throughput
target: the spike uses four rows and the target-flow operator subset. The
intended implementation uses local typed comparisons and deterministic symbolic
or bitset-equivalent finite-domain proof over declared domains; no callback
invocation, expression parser, or external engine is part of the hot path. If
implementation later indexes predicates for speed, the optimization must
preserve the normalized atom semantics, scoped product proof, and exact-one
refusal behavior.

## Finalization Gate

> Complete each item with a written response before
> marking this RDR as **Final**. Written responses
> prevent rubber-stamping and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses below.

### Contradiction Check

The research supports a symbolic guard-atom model over declared tag domains, and
the proposed solution keeps host callbacks and free-form expression strings out
of the contract so lint can prove finite-domain coverage and overlap.

Guard-domain enforcement is resolved rather than open: the kernel decides
presence and existence atoms, so this RDR's evaluator is scoped to value
semantics over a present value (JDR 0001 §D4, a Closed entry).

The strength of the exhaustiveness claim is decided in substance but not yet
agreed across documents. §JD-4 settles that the lint promise narrows and the
runtime veto stands, and this RDR now states that for its own proof. §JD-4
remains open as to which document records it and names RDR 0006's proof; RDR
0006 is `Final` under that tolerance and carries no `guard_unevaluable`
narrowing. A8 tracks the resulting divergence, which must close before lock.

### Assumption Verification

All Critical Assumption records are internally consistent: A1 through A6 are
`Verified`, each uses an allowed Method label, each has concrete Evidence, and
each has a non-empty "If wrong" consequence. No record uses `Docs Only`. Two
records remain `Pending` — A7 and A8, each with a named plan below; neither is
`Unverified`.

A1 is backed by the Resolve spike transcript in
`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/output.txt`. A2
is a derivation of the finite scoped-product proof. A3 is an explicit design
decision that rejects inline negation and nested boolean expressions. A4 cites
source-search anchors that resolve now:
`internal/cli/clierr/clierr.go::CLIError`,
`internal/cli/respond/respond.go::Fail`, and
`internal/cli/config/config.go::Load`. A5 rests on the kernel-fixed parsed-atom shape (JDR 0001 §D1, normative in RDR
0007) plus RDR 0002's source-identity contract, and on
`internal/resolve/resolve.go::Row`, which already carries `RuleID` and
`SourceLocator`; A6 relies on peer RDR contracts in RDR 0002 and RDR 0006.
Neither is self-reference. No `Source Search` Evidence cites this RDR or its
artifact directory.

**A7 is Pending and this RDR is NOT lockable until it resolves.** It carries the
existence-atom projection RDR 0007 A12 routed here, and it is blocked on a
producer request to RDR 0002 (a per-tag optionality declaration) that RDR 0002,
itself `Draft`, does not yet carry. The derivation is recorded in
`evidence/research/iter-2-projection-derivation.md`; the bounded prior-art pass
returned a clean negative, so it resolves as a Design Decision. Until it is
stated as a normative clause, this RDR's exhaustiveness proof is silent on the
one operator RDR 0007 makes total.

**A8 is Pending and also blocks lock.** This RDR now records the §JD-4
narrowing for its own proof, but §JD-4 is an open ledger entry that names RDR
0006's proof and leaves the recording document unassigned, and `Final` RDR 0006
carries no `guard_unevaluable` narrowing. Closing it needs a route-back to RDR
0006 or a §JD-4 disposition, not an edit here.

### Scope Verification

The Minimum Viable Validation is in scope for implementation, not deferred. The
specific proof is a production test fixture that encodes one RDR flow slice and
one kata flow slice as normalized candidate rows, including equality, enum
membership, set containment, bounded integer comparison, existence, mixed
`all`/`unless`, one complete partition, one intentional product-level gap, one
intentional product-level overlap, and refusal/downgrade cases for unbounded or
too-large finite domains.

### Cross-Cutting Concerns

- **Incremental adoption**: guards over unbounded dimensions remain runtime
  evaluable, but lint must refuse or downgrade exhaustiveness for those
  dimensions rather than blocking all use of predicates.
- **Memory management**: set-valued finite domains must use deterministic
  symbolic or bitset-equivalent proof; the RDR explicitly rejects naive powerset
  materialization when the finite product is too large to prove.
- **Canonical-form / determinism**: successful matching and lint findings cannot
  depend on source order. Exact-one selection belongs to RDR 0001, source-row
  identity comes from RDR 0002, CLI envelopes belong to RDR 0005, and blocking
  graph-lint findings belong to RDR 0006.

### Proportionality

The document is right-sized for lock. It owns one independent load-bearing
contract: the guard-predicate grammar and finite-domain exhaustiveness
semantics. Neighboring surfaces are explicitly delegated to peer RDRs rather
than locked here: sparse table shape to RDR 0002, exact-one resolution to RDR
0001, CLI envelope mapping to RDR 0005, and graph-lint authority to RDR 0006.

The `large` Profile still matches because this RDR locks a grammar and
finite-domain proof semantics. The required large-profile lenses ran
(`grounding`, `3amigo`, and `critique`); every finding they raised is either
resolved in the live text above or carried as a named obligation in the Minimum
Viable Validation and Testing Strategy. The re-entry re-runs that row against
the revised text: `grounding` iteration 2 is complete, and its `authority`,
`disposition`, and `trace` mini-check tables are carried above.

## References

- RDR 0001, Resolution Kernel Contract.
- RDR 0002, Transition Table As Reviewable Data.
- RDR 0006, Graph Lint Authority and Guarantees.
- RDR 0007, Guard Predicate Totality — normative home of the guard seam, the
  domain rule, and the `guard_unevaluable` payload.
- JDR 0001, Resolve Kernel Seam — §D1 (parsed-atom row shape), §D4 (kernel
  enforces the guard domain), §JD-4 (the lint promise narrows — substance
  decided, recording document still open, head clause names RDR 0006).
- `docs/cli-output-contract.md`.
- Resource index: `.rdr/resources.md`.
- Seed prior: `../state-machines/BUILD-SEEDS.md`, especially the guard
  expression and transition-table seeds.
- Transition model prior: "The transition model — inputs, outputs, error
  conditions."
- Prior-art corpus: `../state-machines` audits, evals, contrasts, and checked
  repositories for `transitions`, stateless/qmuntal-stateless, SCXML/scxmlcc,
  Sismic, StateSmith, and Statewright.
