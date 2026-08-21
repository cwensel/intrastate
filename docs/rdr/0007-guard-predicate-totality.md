# Recommendation 0007: Guard predicate totality over an incomplete evaluation view

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-08-09
- **Status**: Draft [revised from Final 2026-08-12; re-proposed
  2026-08-21 under JDR 0001 §D1/§D2/§D3 — re-verify A3, A16, A17,
  A18, A19, A21, A22]
- **Type**: Architecture
- **Profile**: foundational — one contract, the domain of
  guard-predicate evaluation over an incomplete view, now
  enforced by the kernel; cross-RDR producer consumed by RDR
  0003's value evaluator and RDR 0002's normalizer.
- **Priority**: Medium
- **Related Issues**: kata `xg7p` — "RDR 0001:
  RequiresOwned conflates guard-input state with post-guard
  transition state" (`src:roborev`, `severity:medium`,
  `area:internal-resolve`).
- **Predecessors**: 0001-resolution-kernel,
  0003-guard-predicate-exhaustiveness; JDR 0001
  (`docs/jdr/0001-resolve-kernel-seam.md`) §D1, §D2, §D3.
- **Overrides**: None superseded. This RDR *absorbs* the
  `resolve.Row` change JDR 0001 §D1 assigns to it (the guard
  seam carries parsed atoms, not a string), *fixes* the
  meaning of `Row.RequiresOwned` (post-guard write-dependency
  keys), and is the single normative home of the
  guard-evaluation domain rule. RDR 0001 (`Implemented`) has
  no backward edge; its `Row` shape changes here by JDR
  decision, not by amendment.
- **Seam Lineage**: `area:internal-resolve` (guard-verdict
  semantics at the RDR 0001 ↔ 0003 boundary) — no prior
  accretion.

## Problem Statement

Someone authoring a transition table — and the guard
predicates its rows depend on — needs to know what happens
when a guard asks about state the runtime view does not
carry. They want a refusal they can act on: told plainly
that the artifact state needed to decide was missing, not
handed a plan that quietly routed around it. They discover
they need this the first time a table with a modeled escape
produces a plan in a situation where the honest answer was
"I could not tell." At that point the table looks like it
worked, and nothing in the output distinguishes "this row
does not apply" from "I could not read what I needed to
find out."

Internally, the resolution kernel must fix the **domain of
guard-predicate evaluation**. When a guard atom references
a tag absent from the assembled evaluation view, the seam
today may answer with a two-valued verdict — treating
absence as non-satisfaction — or refuse to answer at all.
The two answers are not interchangeable: a `GuardFalse`
prunes the row before it can raise an owned-state
obligation, and the resulting `no_match` is escapable per
RDR 0002, while `guard_unevaluable` is not. So the choice of
domain decides whether missing artifact state can be masked
behind an escapable refusal class.

Nothing normative states which answer is required. The only
statement of `Row.RequiresOwned`'s meaning anywhere is a Go
doc comment — `internal/resolve/resolve.go`'s "RequiresOwned
names owned tag keys the row's evaluation needs" — which
reads as covering both the state a guard reads to be
decidable and the state a transition's writes depend on once
the guard holds. This RDR must state, normatively and in one
place, whether the guard predicate is a total function over
the view or a partial function whose domain is the set of
tags it references — and, since JDR 0001 §D1, **which
component enforces that rule**: the kernel that now sees
every atom, or the evaluator that decides operator meaning.

## Context

### Background

Raised by a roborev finding against the RDR 0001
implementation (kata `xg7p`), then grounded by live probes
against `internal/resolve/resolve.go` during seed triage.
Observed current behavior, not a design position:

- An overlap key present as *observed* with a matching value
  yields `owned_state_unavailable` with `MissingOwned:[gate]`.
- An overlap key absent entirely yields `no_match` — the row
  genuinely does not match (the match pattern is closed-world
  by design).
- A FALSE guard plus an absent `RequiresOwned` key plus a
  modeled `no_match` escape yields a plan via the escape,
  `Escaped:true`. This is the path where missing artifact
  state is masked.

Constraints inherited from prior rounds. RDR 0001's
deviation **D8** (guard-FALSE prunes first; a pruned row
contributes neither candidacy nor an owned-state
obligation) is pinned by frozen tests ADV-1 and ADV-1b.
Deviation **D5** scopes `RequiresOwned` and guard evaluation
to candidate rows only. RDR 0002 closes escape lists to
`no_match` and `ambiguous_match`, which is the asymmetry that
makes verdict ordering load-bearing.

**JDR 0001 re-entry (2026-08-12).** The cluster gate found
the prior Final of this RDR built on the opaque
`Row.Guard string` seam, with a mandated panic as the only
surface for an evaluator that could not parse its guard.
JDR 0001 resolved three decisions this RDR now absorbs:

- **§D1** — the row carries **parsed atoms** (key, operator
  token, literal, block); no reconstruction step, so no
  mapping failure and no panic clause. Lands in this RDR,
  which takes the `Row` change. The prior A3 ("no kernel
  change") is superseded; A10 and A15 (string-encoding
  obligations routed to RDR 0003) are dropped.
- **§D2** — **gate, then count**: guards evaluate over
  surviving candidates, an unevaluable candidate vetoes,
  and only then exact-one matching and escape reachability
  apply. Stated here as the kernel's ordering; RDR 0002
  restates its resolver flow to match.
- **§D3** — a read accessor MUST return the complete tag
  set or refuse. Closes the prior A6b; RDR 0004 carries the
  clause as its own.

RDR 0003 and RDR 0002 are both **Draft** again (re-entry at
refine) and cite the §D1 atom shape rather than restating
it. 0003 owns operator vocabulary and typed semantics; its
evaluator becomes a *value* evaluator under this RDR.

### Technical Environment

Go; the resolution kernel `internal/resolve` — one non-test
file, `resolve.go::Resolve` as the single pure entry point,
no production importer (only in-package tests), which is why
JDR 0001 P7 prefers the clean shape over a compatible one.
The contract sits at the boundary between:

- **RDR 0001 — Resolution kernel** (`Implemented`). Owns
  the closed five-kind refusal taxonomy
  (`resolve.go::RefusalKinds`, pinned by
  `resolve_test.go::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`)
  and the selection rule ("exactly one matching edge after
  guard evaluation"). Its `resolve.go::gate` is already the
  gate-then-count shape §D2 names: prune GuardFalse, report
  missing owned state among survivors, report undecidable
  guards, and only then `Resolve` counts `selected`.
- **RDR 0003 — Guard predicate exhaustiveness** (`Draft`,
  re-entry at refine). Owns the operator vocabulary
  ("equality, membership, bounded integer comparison,
  existence, and set containment"), the tag-kind/operator
  matrix, and the value evaluator this RDR's seam calls.
- **RDR 0002 — Transition table as reviewable data**
  (`Draft`, re-entry at refine). Owns normalization and the
  escape-class closure; names the `RequiresOwned` producer
  (JDR 0001 §JD-3).
- **RDR 0004** (`Final`) carries read completeness (§D3);
  **RDR 0005** (`Final`, tolerance §JD-8/§JD-9) renders
  refusals; **RDR 0006** (`Final`, tolerance §JD-4) narrows
  lint's promise where this RDR's veto and its proof
  disagree.

## Research Findings

### Investigation

Prior art was read before enumerating approaches; the query
ledger and accepted/rejected citations are cached at
`docs/rdr/0007-guard-predicate-totality/evidence/research/propose-prior-art.md`
(iteration 1 for the domain question, iteration 2 for the
enforcement-site question §D1 opened).

**Domain question (unchanged).** ⚠ no prior-art coverage in
the arc corpora (`StateMachineRes` ×3, `StateMachineLit` ×1)
for the undecidable-guard class; the decisive external
anchor is W3C SCXML §5.9.1: "If a conditional expression
cannot be evaluated as a boolean value ('true' or 'false')
or if its evaluation causes an error, the SCXML Processor
MUST treat the expression as if it evaluated to 'false' and
MUST place the error 'error.execution' in the internal event
queue." ⇒ the only standard that folds unevaluable into
`false` does so only with a mandatory observable error on a
second channel; the kernel has one channel (the verdict), so
the honest transplant surfaces the error in the verdict.
SCXML continues after the error where this RDR halts; the
citation supports "never silently false," nothing more. The
in-repo house rule is the same one seam over: RDR 0004's
"Indeterminate MUST be a refusal-class result, not a false
allow and not a false deny" (A8).

**Enforcement-site question (new).** Once the row carries
parsed atoms, the kernel can see each atom's referenced key
and decide presence itself. The instance precedent is
SQL:2003's strict-routine rule (research C3): "When an SQL
function is called whose CREATE FUNCTION statement specifies
RETURNS NULL ON NULL INPUT, then if the runtime value of any
argument is null, the result of the function call is set to
null, and the function itself is not invoked." ⇒ the *host*
enforces the partial domain and never trusts the routine to
re-implement null propagation — which maps directly onto
"the kernel marks an absent-key atom unevaluable and never
calls the evaluator for it." PostgreSQL `STRICT` is the same
rule (`CALLED ON NULL INPUT` vs `RETURNS NULL ON NULL
INPUT`/`STRICT`, pg13 CREATE FUNCTION). Corpus hits for
"engine vs UDF null handling" beyond this were off-instance
(Date's anti-null argument; Datalog 3-valued instances) and
were rejected.

**Sibling-path check.** The discriminator exists:
`resolve.go::GuardResult` is three-valued, `resolve.go::KindGuardUnevaluable`
is a closed refusal kind, and `resolve.go::evaluateGuard`
already answers the nil-seam case `GuardUnevaluable`, never
`GuardFalse` ⇒ the kernel already refuses on one class of
undecidability; this RDR extends that to absent keys. The
presence test exists: `resolve.go::TagSet.Lookup` returns
`ok` provenance-blind and is the only exported accessor ⇒
it is the kernel's presence predicate for atoms. The
owned-only test `resolve.go::TagSet.has` (unexported, single
caller `resolve.go::missingOwned`, written by `resolve.go::assemble`
with `ProvenanceOwned`) answers a different question and is
deliberately not reused. No sibling path combines
three-valued verdicts across atoms — searched
`internal/resolve`; none exists — so the strong-Kleene
combination is new kernel logic, not a parallel of an
existing one.

**Dropped at re-entry.** A10 (string-encoding mapping) and
A15 (conjoined atoms through one string) verified premises
§D1 removed; both are struck, not re-verified. Their
authoring-grammar residue (may one tag carry two operator
keys in TOML; may a tag sit in both blocks) is RDR 0003's,
now Draft, and is handed over in Phase 4 rather than routed
as an obligation onto a Final peer.

### Key Discoveries

- **Documented** — SCXML §5.9.1: unevaluable → `false` is
  legitimate *only* with a mandated observable error; silent
  masking has no prior-art support.
- **Documented** — SQL:2003 strict routines: the host, not
  the routine, enforces the null domain; the routine "is not
  invoked" on a null argument (research C3).
- **Documented** — The kernel is already gate-then-count:
  `resolve.go::Resolve` calls `gate(candidates, in.Guards, view)`,
  returns on `blocked != nil`, and only then switches on
  `len(selected)`; `resolve.go::escapeOrRefuse` gates the
  escape set the same way. §D2 names shipped behavior.
- **Documented** — Zero production importers of
  `internal/resolve` (only its own `_test.go` files import
  it), so reopening `Row` migrates nothing.
- **Documented** — RDR 0002's escape closure ("An `escape`
  list MUST contain only resolver failure classes that RDR
  0001 allows the table to model: `no_match` and
  `ambiguous_match`") makes `guard_unevaluable` structurally
  non-escapable.
- **Assumed** — the evaluator seam needs nothing but the
  atom and the present value to decide a value operator
  (A17); and the kernel can recognize the existence operator
  by one token the normalizer emits (A16).

### Critical Assumptions

- **A1 RDR 0003's closed operator vocabulary needs no new
  operator for this rule: `exists` is the only operator
  whose declared semantics name absence; the remaining four
  are value-comparing and RDR 0003 is silent on their
  absent-operand behavior, which this RDR fixes.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0003 Normative Contracts — "The
    initial operator vocabulary MUST be closed and typed:
    equality, membership, bounded integer comparison,
    existence, and set containment." Its operator/kind
    matrix marks `exists` alone as absence-inspecting
    ("Tests presence or absence, not value equality") while
    `eq`/`in`/`lt,lte,gt,gte`/`contains` are stated in value
    terms. A sweep of RDR 0003 for `absent|absence|optional`
    returns only the `exists` clauses. This RDR's contracts
    supply the absent-operand rule and pin `contains`: an
    absent set-valued tag is *unevaluable*, not the empty
    set.
  - **If wrong**: an operator that is neither
    presence-deciding nor value-comparing has no domain
    rule; the kernel's `exists`-or-value dichotomy (A16)
    misclassifies it.
- **A2 Under strong-Kleene combination, no combination of
  present-tag verdicts and absent-tag unevaluability yields
  `GuardFalse` where absence contributed the falsity.**
  - **Status**: Verified
  - **Method**: Derivation
  - **Evidence**: Conjunction is `min` under `F < U < T`;
    negation `¬T=F, ¬F=T, ¬U=U`. The row-verdict shape
    `all_result ∧ ¬(unless_conj)` is RDR 0003's ("Positive
    `all` predicates are conjunctive requirements. Negative
    `unless` predicates are also conjunctive within the
    excluded predicate set: if all `unless` predicates
    hold, the candidate row is disabled"; "`unless` is not
    per-atom negation"), stated there over complete
    assignments; lifting it to three values is this RDR's
    extension. `GuardFalse` arises only via (i) some `all`
    atom decided `F` — a present-tag comparison or an
    existence atom — since `U` never yields `F` in a
    conjunction without an `F` operand; or (ii)
    `unless_conj = T`, which needs every `unless` atom `T`,
    which a `U` atom blocks rather than drives. The
    `F ∧ U = F` cell is witnessed falsity: it holds under
    every substitution for the unevaluable atom, so pruning
    discards no row that could have survived — this is what
    keeps D8 sound. Cross-checked against PostgreSQL 16 §9.1
    (`FALSE AND NULL = FALSE`, the *strong*-Kleene
    discriminator; A4).
  - **If wrong**: the masking path reopens through
    combination even though single atoms are honest.
- **A3 The kernel change is bounded: replacing
  `Row.Guard string` with a parsed atom slice, narrowing
  `GuardEvaluator` to a per-atom value seam, and adding
  presence + strong-Kleene combination touch `Row`,
  `GuardEvaluator`, `evaluateGuard`, `Refusal`, and every
  test fixture — but leave `gate`'s prune → owned-state →
  undecidable order, `escapeOrRefuse`'s uniform gating, and
  the frozen ADV/Fixup dispositions unchanged.**
  - **Status**: Pending — supersedes the prior A3 ("no
    kernel behavior change"), which JDR 0001 §D1 voided.
    The "zero production importers" bound is a
    point-in-time fact about callers sibling RDRs (0003's
    evaluator, 0005's CLI) have not written yet — the
    escaped-defect class in this project's ledger — so it
    carries a re-verification trigger: re-run the importer
    sweep at Phase 1 start, and if an importer exists, the
    `Row` change is a coordinated migration, not a free one.
  - **Method**: Spike
  - **Evidence**: Inspection so far: `resolve.go::gate`
    consumes only the `GuardResult` from
    `resolve.go::evaluateGuard` and never reads `row.Guard`
    except to copy it into `Refusal.Guard`; the per-row
    verdict is therefore replaceable behind `evaluateGuard`
    without touching the partition. The frozen tests —
    `adversarial_test.go::TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch`,
    `::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`,
    `::TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`,
    `::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`,
    `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`,
    `fixup_test.go::TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates` —
    decide verdicts through `fixtures_test.go::fixtureGuards`
    keyed on guard *text*, so each must be re-expressed with
    atoms; their asserted dispositions do not depend on the
    guard's representation. The spike re-encodes the
    fixtures and confirms every disposition holds.
  - **If wrong**: a frozen disposition depends on guard-text
    identity (e.g. `Refusal.Guard` equality), and Phase 1
    must reopen that test's contract rather than re-encode it.
- **A4 Prior art supports "definite results from present
  data dominate; unknown propagates otherwise" as the
  standard three-valued treatment of missing data
  (SQL/Kleene K3).**
  - **Status**: Verified
  - **Method**: Prior Art
  - **Evidence**: PostgreSQL 16 Documentation §9.1 "Logical
    Operators" — "SQL uses a three-valued logic system with
    true, false, and null, which represents 'unknown'", with
    the AND/OR/NOT tables; `FALSE AND NULL = FALSE` and
    `TRUE OR NULL = TRUE` identify strong Kleene, and the
    commutativity note rules out a short-circuit reading.
    Quote cached in `evidence/research/resolve-citations.md`
    (C3).
  - **If wrong**: the combination rule stands on A2 alone.
- **A5 Authors who need "row applies when tag X is absent"
  express it with an existence atom; the disjunction "absent
  OR equals v" needs two rows (RDR 0003 rejects intra-guard
  disjunction), and the value row MUST itself carry an
  `X exists = true` atom beside `X eq v`, or it survives
  unevaluable when X is absent and the aggregation veto
  refuses the absent leg.**
  - **Status**: Verified (runtime); the load-time half rides
    on A12.
  - **Method**: Design Decision + Source Search
  - **Evidence**: RDR 0003's placement rule is "Positive
    guard atoms MUST live in `all`; negative guard atoms
    MUST live in `unless`" — placement only; its A3 rejects
    "Inline `not` and nested boolean expressions", so the
    two-row pattern is the only disjunction route; the rows
    are disjoint (no assignment has X both absent and
    present). The conjoined value row — two atoms over one
    key — is admitted by the kernel natively once `Row`
    carries an atom slice (cardinality, not grammar).
    Whether RDR 0003's TOML grammar lets an author *write*
    it is 0003's (Draft) and is handed over in Phase 4.
  - **If wrong**: authors cannot express a legitimate
    absence-conditional row and fall back to
    sentinel-stamping upstream — the anti-pattern whose only
    sanctioned alternative *is* this pattern.
- **A6a Kernel-side assembly is faithful: no tag key present
  in the input tuple is ever dropped from the assembled view
  by merge precedence; shadowing changes a value, never a
  key's presence.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `resolve.go::assemble` — every write is an
    unconditional map assignment keyed by tag key; no
    `delete`, filter, or conditional skip; the one
    conditional gates an *insert*. The four readers
    (`TagSet.Lookup`, `TagSet.has`, `TagSet.matches`,
    `TagSet.Len`) remove nothing. Owned-over-observed
    precedence overwrites a value in place.
  - **If wrong**: spurious non-escapable refusal storms, or
    silent plans through a negative existence atom.
- **A6b Read completeness at the accessor→kernel boundary is
  stated by RDR 0004, not here.**
  - **Status**: Verified — closed by JDR 0001 §D3 / §JD-7.
  - **Method**: Prior Decision
  - **Evidence**: JDR 0001 §D3: "A read accessor MUST return
    the complete tag set for the keys it was asked for, or
    take the refusal branch; a partial read is a refusal,
    not a value." Lands in RDR 0004 as one normative clause
    plus an MVV scenario. This RDR cites it and carries no
    obligation of its own.
  - **If wrong**: a truncated snapshot and a genuine absence
    are byte-identical at `Input.Owned`, and a negative
    existence atom decides TRUE against state that exists;
    operators must then treat an unevaluable refusal as
    "inspect the accessor," never "retry."
- **A7 Every guard-referenced key is declared, so narrowing
  `RequiresOwned` leaves no guard input undeclared: a typo'd
  key fails table load.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0003: "A guard predicate MUST be a
    symbolic atom over a declared tag"; "The tag name must
    resolve to a declared tag." The hard MUST-reject is RDR
    0002's validation categories ("The model MUST declare
    every tag it matches or writes, including each tag's
    provenance"; `unknown tag` MUST be retained) — keyed on
    matching, independent of ownership.
  - **If wrong**: a misspelled key is forever absent —
    permanently unevaluable, or shadowed under
    FALSE-domination — and invisible at review time.
- **A8 The no-collapse principle is a house rule settled at
  the adjacent accessor seam.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0004 Normative Contracts — "A gate
    accessor MUST return allow, deny, or indeterminate.
    Indeterminate MUST be a refusal-class result, not a
    false allow and not a false deny." Carried by RDR 0005
    and shipped as `flow-gate-indeterminate`. 0004 reports
    indeterminate as an execution outcome; this RDR derives
    it from an absent key — same principle, distinct seams.
  - **If wrong**: the rule stands on SCXML/K3 alone; the
    choice survives.
- **A9 The empty-atom-block identities do not contradict RDR
  0002's normalized row shape: an omitted `unless` block is
  absent, not a vacuously-true conjunction, and normalization
  never synthesizes one.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0002's only `unless` normative fixes
    representation and combination timing ("Guard predicates
    MUST be represented as positive `all` predicates and
    negative `unless` predicates. Normalization MUST combine
    both into one candidate-row predicate set before
    ambiguity checks"), never a zero-atom truth value; it
    fixes the analogous absence identity for writes ("Absence
    from both the write block and the clear list MUST NOT
    imply deletion"), showing it states such identities when
    it means to. RDR 0003's subtractive algebra ("its
    `unless` block denotes an excluded intersection that is
    subtracted from the row's accepted assignments") requires
    the same reading. Both canonical fixtures omit `unless`
    on most rules; the 0002 spike normalizer contributes zero
    predicates for an absent block. Now kernel-implemented,
    so Scenario 9(b) becomes a kernel test.
  - **If wrong**: a table with an omitted `unless` disables
    every row carrying one.
- **A11 Deciding guard presence provenance-blind conflicts
  with no shipped kernel behavior.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `resolve.go::TagSet.Lookup` returns `ok`
    on any key present regardless of provenance and is the
    only exported accessor; `TagSet.has` takes a provenance
    and serves only `missingOwned`; `TagSet.matches` tests
    key and value with no provenance test. RDR 0003's fixture
    guards `cluster_eligible` (declared `provenance =
    "observed"`) with `exists = true`.
  - **If wrong**: `exists` over an observed tag decides
    differently than over an owned one and RDR 0003's own
    fixture breaks.
- **A12 The two-row absence pattern survives RDR 0003's
  load-time overlap check: either row-group scoping keeps
  the rows out of one compared product, or absence has
  standing in the declared-domain product.**
  - **Status**: Verified — indeterminate by RDR 0003's
    silence; an inherited obligation with home JDR 0001
    §JD-4 (lint's promise narrows), not a refutation.
  - **Method**: Source Search
  - **Evidence**: RDR 0003 defines overlap over a product of
    declared finite domains and is silent on how an
    existence atom projects onto it; RDR 0002's
    tag-declaration schema has no optionality field. RDR
    0002 lists "ambiguous overlap" among categories
    validation MUST retain, with no downgrade valve. RDR 0003
    and RDR 0006 are Draft/tolerance-qualified under §JD-4
    and own the projection and row-group membership rules.
  - **If wrong**: the domain rule survives; its authoring
    story for the disjunction case does not.
- **A13 Provenance-blind presence is an accepted exposure:
  a caller-supplied observed tag can turn
  `guard_unevaluable` into a decided verdict.**
  - **Status**: Verified — accepted exposure; home JDR 0001
    §JD-9 (`--tag` enters as `ProvenanceObserved`, never
    satisfies an owned dependency).
  - **Method**: Source Search
  - **Evidence**: `resolve.go::assemble` writes
    `Input.Observed` unconditionally; RDR 0005 treats
    `--tag` as "context already known to the caller"; RDR
    0001 names "caller-supplied observed tags" in its kernel
    contract. Unreachable today — no production importer,
    `internal/cli/root.go` registers only `version` — live
    at RDR 0005's CLI handoff. The owned path is already
    defended by
    `adversarial_test.go::TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`
    and `::TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement`;
    no test contends the guard path — Phase 2 pins it.
  - **If wrong**: an operator unsticks a refusal from the
    command line and the plan's `Writes` derive from caller
    input.
- **A14 An `exists` atom carries a boolean literal whose
  value is the expected presence (`presence == literal`), so
  `all … exists = false` is a legal single-atom absence
  test.**
  - **Status**: Verified — with a recorded coupling: the
    cluster gate (pairwise 0007×0003 F2) read RDR 0003's
    placement MUST as contradicted by a `false` literal in
    `all`; JDR 0001 routes it to §JD-1, and RDR 0003 is Draft
    and can state which reading it means. The kernel rule
    here is literal-polarity regardless.
  - **Method**: Source Search
  - **Evidence**: RDR 0003 never defines "positive"/
    "negative" atom; its Identity decision keys an atom by
    "source rule/context id plus its position within `all`
    or `unless`" with no polarity field; "`unless` is not
    per-atom negation". The `exists` row's proof role "Tests
    presence **or absence**" has content only if both
    literals are authorable; `false` parses as boolean.
  - **If wrong**: absence is expressible only via
    `unless … exists = true`; the two channels still reach
    the same verdict under this RDR's rule.
- **A16 The kernel can recognize the existence operator by
  one operator token the normalizer emits, and this is the
  only grammar token the kernel needs: every other operator
  is a value comparison the seam decides.**
  - **Status**: Pending
  - **Method**: Source Search (RDR 0003 vocabulary and RDR
    0002 normalization) + Design Decision
  - **Evidence**: JDR 0001 §D1 fixes the atom as "(key,
    operator token, literal, block)" and says the kernel
    "carries *cardinality*, not *meaning*". Recognizing one
    token as presence-deciding is the minimum meaning the
    domain rule needs; the alternative — an atom field
    marking presence operators — moves the same fact one
    level up without removing it. The token's spelling is
    RDR 0003's; the kernel pins it as one exported constant
    the normalizer MUST use.
  - **If wrong**: the kernel cannot separate `exists` from
    value operators and the enforcement site falls back to
    the evaluator (Alternative 1) for that distinction.
- **A17 A value operator over a present key is decidable
  from the atom and the present value alone: the seam
  `Evaluate(atom, value)` needs neither the view, nor
  provenance, nor sibling atoms, nor tag declarations at
  evaluation time.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: RDR 0003's operator semantics — "equality
    compares a tag value to one typed literal; membership
    checks a scalar tag against a typed literal set; bounded
    integer comparison uses `lt`, `lte`, `gt`, and `gte`;
    existence checks presence"; set containment "Narrows
    the set-valued domain to assignments containing every
    listed element" — each reads one value against one
    literal. The declared kind needed to *parse* the value
    is known to the evaluator from the table it was built
    for, not from the kernel.
  - **If wrong**: the seam takes whatever RDR 0003 shows the
    evaluator needs (e.g. the declared kind on the atom) —
    but never the view; the domain rule stays kernel-side.
    The seam's ceiling is stated as a constraint on RDR
    0003: one tag, one literal, one operator. A relational
    (tag-to-tag), temporal, or aggregate operator needs a
    seam extension and is out of this RDR's scope; RDR
    0003's vocabulary is closed to the five named operators
    (A1), so none exists today.
- **A18 The seam MAY answer `GuardUnevaluable` for a present
  value it cannot compare (a value that does not parse as
  the declared kind), and this does not reopen masking:
  the verdict combines under strong Kleene exactly as an
  absent-key atom does, and the refusal still names the row.**
  - **Status**: Pending
  - **Method**: Design Decision + Source Search (RDR 0003's
    parse/lint rules for runtime values)
  - **Evidence**: An unparseable present value is the same
    class as the nil-seam case `resolve.go::evaluateGuard`
    already answers `GuardUnevaluable`: the evaluator cannot
    decide, and collapsing that into `false` would be the
    masking path one layer down. Whether RDR 0003 treats it
    as a lint-time impossibility (declared kinds validated
    at load) or a runtime verdict is 0003's; the kernel
    accepts the verdict either way.
  - **If wrong**: the seam must be two-valued and a
    mis-typed runtime value becomes a programmer-error
    return from `Resolve` — a kernel contract change.
- **A19 The `guard_unevaluable` refusal can name the absent
  tag keys directly (the kernel now knows them), and RDR
  0005's envelope carries that payload without a new
  envelope field.**
  - **Status**: Pending
  - **Method**: Source Search (RDR 0005 envelope; JDR 0001
    §JD-8)
  - **Evidence**: JDR 0001 §JD-8: resolver-specific values
    "only require new `Code` constants or literals, **not
    new envelope fields or exit groups**"; "`Detail` and
    `Hint` ship and render in text and JSON." The kernel
    field replaces `Refusal.Guard string` (one row's guard
    text — meaningless once guards are atoms) with the
    sorted, deduplicated absent keys across every
    undecidable row; `Refusal.Rows` keeps row identity.
  - **If wrong**: the refusal names the rows but not the
    keys, and the Problem Statement's "told plainly … was
    missing" is met only at row granularity.
- **A20 Gate-then-count (JDR 0001 §D2) is already the
  kernel's shape; stating it here changes no code.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `resolve.go::Resolve` — `selected, blocked
    := gate(candidates, in.Guards, view); if blocked != nil
    { return refuse(in, *blocked), nil }` precedes `switch
    len(selected)`; `resolve.go::escapeOrRefuse` gates the
    escape set before counting `viable`. RDR 0002's re-entry
    restates its resolver flow to this shape.
  - **If wrong**: n/a — the code is the evidence.
- **A21 `Row.RequiresOwned` on an escape row is empty by
  composition: RDR 0009 requires an escape row's `Writes`
  to be empty, and this RDR defines `RequiresOwned` as the
  keys `Writes` depend on.**
  - **Status**: Pending
  - **Method**: Source Search (RDR 0009 Final; RDR 0002's
    §JD-3 producer clause)
  - **Evidence**: RDR 0009 Normative Contracts: "a Row with a
    non-empty Escape list MUST have an empty Writes slice."
    JDR 0001 §JD-3 names this the open composition (0007
    derives from `Writes`, 0009 empties `Writes`), notes the
    field "appears **zero** times in 0002, which owns the
    normalized row", and calls the D1 `Row` reopening "the
    natural occasion to settle it"; RDR 0002's re-entry
    direction takes the producer question. Under this RDR's
    definition the composition yields the empty set, so an
    escape row raises no `owned_state_unavailable` of its
    own — consistent with `resolve.go::missingOwned` over
    an empty list.
  - **If wrong**: RDR 0002's producer populates
    `RequiresOwned` on escape rows from some other source,
    and the ordering clause must say what an escape row's
    owned obligation is.
- **A22 Tag-key identity at the kernel is exact string
  equality, and every canonicalization of authored key
  spellings (case, namespace, whitespace) happens in RDR
  0002's normalizer before a `Row` exists — so an atom's key
  and the assembled view's key agree byte-for-byte whenever
  the author meant the same tag.**
  - **Status**: Pending
  - **Method**: Source Search (RDR 0002 tag declaration and
    normalization contracts; `resolve.go::assemble`)
  - **Evidence**: `resolve.go::assemble` keys `view.tags` by
    `Tag.Key` verbatim and `TagSet.Lookup` indexes the map
    by the string it is given; the kernel has never
    canonicalized. RDR 0002 requires every matched or
    written tag to be declared, which is where spelling is
    fixed. Raised by premortem P-17.
  - **If wrong**: a spelling the normalizer lets through
    differently in a guard and in a write is "absent" to the
    kernel — an unevaluable refusal, never a plan (fail
    closed), but one that misdiagnoses a producer defect as
    missing state.

## Proposed Solution

### Approach

**The guard predicate is a partial function whose domain is
the set of tags it references; outside that domain the
verdict is "unevaluable," never "false" — and the KERNEL
enforces that domain.** The row carries parsed guard atoms
(JDR 0001 §D1). For each atom the kernel itself decides
presence of the referenced key in the assembled view,
provenance-blind. An existence atom is decided by the kernel
from presence alone (`presence == literal`). A
value-comparing atom over an absent key is marked
unevaluable by the kernel and the evaluator is never called
for it — the strict-routine rule (research C3). Only a
value-comparing atom over a *present* key reaches the seam,
which compares the present value against the literal under
RDR 0003's typed operator semantics. The kernel then combines
atom verdicts under strong-Kleene three-valued logic across
`all` and `unless` into the row verdict, and the shipped
prune → owned-state → undecidable → count pipeline runs on
it unchanged. Because the kernel knows which keys were
absent, the `guard_unevaluable` refusal names them.

Ownership: **this RDR is the single normative home of the
domain rule and its enforcement site.** RDR 0003's evaluator
implements *operator* semantics over present values; RDR
0001's kernel implements presence, existence, combination,
and the refusal. `Row.RequiresOwned` is narrowed to its
post-guard role: the owned keys the row's `Writes` depend on
once the guard holds (kata `xg7p`). D8 (guard-FALSE prunes
first) is ratified as safe *because of* this rule: falsity
can only come from decided atoms.

### Technical Design

The change is confined to `internal/resolve/resolve.go`:
`Row.Guard` becomes an atom slice; `GuardEvaluator` narrows
to a per-atom value seam; `evaluateGuard` grows presence,
existence, and combination; `Refusal` names absent keys.
`assemble`, `gate`, `missingOwned`, `escapeOrRefuse`, and
`Resolve` are untouched in control flow. Test fixtures
migrate from text-keyed verdicts to atoms.

#### Normative Contracts

```normative
SEAM. A candidate row carries its guard as a slice of parsed
atoms — key, operator token, literal, block ∈ {all, unless}
— the shape JDR 0001 §D1 fixes; this RDR cites it and does
not restate the grammar. An empty slice is an unguarded row.
There is no opaque guard string and no reconstruction step,
so no mapping failure exists and no panic or error channel is
needed for one.

The kernel evaluates a guard ATOM BY ATOM. For each atom it
MUST first decide presence of the referenced key in the
assembled view, provenance-blind (the `TagSet.Lookup` `ok`
test). Then:
- an existence atom is decided by the KERNEL from presence
  alone, and the evaluator MUST NOT be consulted for it;
- a value-comparing atom whose key is ABSENT is marked
  unevaluable by the KERNEL, and the evaluator MUST NOT be
  consulted for it;
- a value-comparing atom whose key is PRESENT is handed to
  the evaluator seam with the atom and the present value:
  `Evaluate(atom, value) GuardResult`. The seam decides
  true or false under RDR 0003's typed operator semantics,
  and MAY answer unevaluable for a present value it cannot
  compare (A18). It never sees the view.

The kernel therefore enforces the domain rule by structure:
an evaluator cannot fold absence into false, because it is
never asked about an absent key.
```

```normative
Value-comparing guard operators (equality, membership,
bounded integer comparison, set containment) are PARTIAL
over the assembled evaluation view: an atom whose referenced
tag key is absent from the view MUST evaluate to unevaluable
— never to false and never to true. This rule is fixed HERE;
RDR 0003 names absence only for `exists`. An absent
set-valued tag MUST be treated as unevaluable under set
containment, NOT as the empty set.
```

```normative
The existence operator is the sole TOTAL operator: it MUST
decide true or false from presence or absence of its
referenced key alone. Its verdict is `presence == literal`:
`exists = true` decides TRUE when the key is present,
`exists = false` decides TRUE when it is absent. Polarity
lives in the literal; `all`/`unless` placement composes with
it. Absence tests MUST be expressed as existence atoms; no
value-comparing operator may act as an implicit existence
test. The kernel recognizes the existence operator by one
operator token, and its literal by exactly two boolean
literal forms, all three exported as kernel constants that
the normalizer MUST emit for existence atoms (A16). Drift at
this boundary fails CLOSED: an existence atom carrying a
foreign token is a value atom to the kernel (unevaluable on
absence, handed to the seam on presence), and a foreign
literal is a producer defect the normalizer MUST reject at
load — neither can yield a plan from absence.
```

```normative
PRESENCE IS PROVENANCE-BLIND. "Present in the assembled
view" means the key is in the view under ANY provenance —
owned, observed, or recognized. An atom MUST NOT be decided
differently according to how a tag reached the view. Key
identity is exact string equality on the key as assembled;
canonicalization of authored key spellings is the
normalizer's (RDR 0002), upstream of the kernel, and the
kernel performs none (A22). This is deliberately NOT the
kernel's owned-state predicate:
`missingOwned` tests owned provenance because
`owned_state_unavailable` is about the owned snapshot; guard
decidability takes the broader test. The two MUST NOT be
conflated.
```

```normative
A row with no atoms is decided TRUE without consulting the
view (shipped: `evaluateGuard`'s empty-guard branch). Empty
atom blocks take the conventional identities: an empty `all`
block is TRUE (empty conjunction); an empty or omitted
`unless` block is ABSENT and contributes nothing — NOT a
vacuously-true conjunction, which under `all ∧ ¬(unless_conj)`
would disable every row carrying one.
```

```normative
Atom verdicts combine in the KERNEL under strong-Kleene
three-valued logic. Conjunction is `min` under `F < U < T`;
negation is `¬T = F`, `¬F = T`, `¬U = U`:

| ∧ | T | F | U |
|---|---|---|---|
| **T** | T | F | U |
| **F** | F | F | F |
| **U** | U | F | U |

The row verdict is `all_result ∧ ¬(unless_conj)`, where
`unless_conj` is the conjunction of the `unless` block's atoms
(`unless` is block-level negation, NOT per-atom negation). A
guard is GuardTrue or GuardFalse only when decided atoms
alone determine it; any unresolved dependence on an
unevaluable atom yields GuardUnevaluable. A present-tag atom
decided FALSE still yields GuardFalse beside an unevaluable
atom (`F ∧ U = F`): the falsity is witnessed by present
state and holds under every resolution of the unevaluable
atom — this is what makes D8 sound. The tables are
normative; the evaluator holds no part of them.
```

```normative
GATE, THEN COUNT (JDR 0001 §D2, stated as the kernel's
ordering). Over the candidate set: guard-FALSE rows prune
first (D8); absent owned state among survivors is reported
next as `owned_state_unavailable`; then any surviving row
whose guard is unevaluable vetoes the resolution as
`guard_unevaluable`; only then does RDR 0001's exact-one
count run over the rows the gate returned, and only then is
escape reachability consulted. `no_match` and
`ambiguous_match` sit downstream of all three gate facts.
The escape set is gated identically, as its own row set
(D5): an unevaluable CANDIDATE is never masked by a decidable
escape row (the gate returns before the escape path is
reached), and an unevaluable ESCAPE row yields
`guard_unevaluable` in place of the candidate-set refusal
(shipped: `escapeOrRefuse` returns the escape set's blocking
refusal; frozen by ADV-2 "guard UNEVALUABLE must not rescue"
and Fixup-1d).

The owned-before-unevaluable precedence is pinned to shipped
behavior, not to D8's original rationale ("absent owned state
is frequently the reason the seam could not decide"), which
this RDR's narrowing invalidates: once `RequiresOwned` names
post-guard write dependencies, the two refusals diagnose
independent problems. The honest cost is that a row failing
both ways surfaces the write-dependency problem first.
Authoring docs MUST NOT repeat the superseded rationale.
```

```normative
The kernel maps GuardUnevaluable to `guard_unevaluable`,
which RDR 0002 excludes from escape lists. Missing artifact
state MUST NOT be maskable behind an escapable refusal class.

The refusal MUST name what was missing, per row and per
atom: for every undecidable row, its `(RuleID,
SourceLocator)` and, for each of its unevaluable atoms, the
referenced key, the block, and a reason drawn from a closed
set — `absent` (the key was not in the view) or
`uncomparable` (the key was present and the seam could not
compare the value, A18). This payload replaces the
single-valued `Refusal.Guard` text, which has no referent
once guards are atoms; `Refusal.Rows` continues to carry
every undecidable row. The kernel MUST evaluate every atom
of every survivor — no short-circuit on a decided block —
and MUST sort the payload by row identity then key, so the
payload is a function of the input tuple (RDR 0001 REQ-1),
never of atom or row order. Any contract or test
distinguishing WHICH row went unevaluable MUST assert on the
row entries.

What the payload does NOT promise: a row pruned as
GuardFalse under `F ∧ U = F` is not a survivor, so its
absent keys are not reported. The guarantee is "never a plan
from absence," not "every absence reported"; an absence is
named exactly when it blocked a verdict.
```

```normative
Row.RequiresOwned names the owned tag keys the row's
post-guard transition depends on — the keys its `Writes`
require, which per RDR 0002 includes an authored clear
(normalization renders it as a `<clear>` write). Guard
decidability is not RequiresOwned's job: guard-input coverage
is enforced by the domain rule above. Listing a guard-read
owned key in RequiresOwned remains legal and yields the more
precise owned_state_unavailable diagnosis among survivors —
and is the ONLY way a guard over an owned tag is protected
from a caller-supplied observed tag satisfying presence
(A13): the guard's key set is NOT required to be a subset of
RequiresOwned, because guards legitimately read observed and
recognized tags. On an escape row, whose `Writes` MUST be
empty (RDR 0009), RequiresOwned is therefore empty (A21). Who
populates the field is RDR 0002's (JDR 0001 §JD-3).
```

```normative
SURVIVOR MEMBERSHIP. A row whose guard is GuardTrue or
GuardUnevaluable is a SURVIVOR; only GuardFalse rows are
pruned. The owned-state scan runs over survivors, so an
unevaluable row's RequiresOwned keys DO raise
owned_state_unavailable, and the undecidable check runs over
the same set (shipped partition in `gate`).

Aggregation is resolution-level: if any surviving candidate
row's guard is GuardUnevaluable, the resolution MUST refuse
`guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be
selected while an unevaluable candidate exists. This veto is
the cost the RDR accepts: one unreadable row refuses a table
whose other rows decide cleanly. Narrowing it would decide
that an undecided edge is a non-edge — absence-as-false at
resolution scope. Where RDR 0006's exhaustiveness proof and
this veto disagree, lint's promise narrows (JDR 0001 §JD-4).
```

```normative
Deviation D8 (guard-FALSE prunes first; a pruned row
contributes neither candidacy nor an owned-state obligation)
is ratified as normative, conditional on the domain rule:
pruning is safe exactly because GuardFalse can only arise
from decided atoms — value comparisons over present keys, or
existence atoms — never from absence folding into a value
comparison.
```

#### Load-Bearing Decisions

- **Enforcement site** — the kernel enforces presence,
  existence, and combination; the evaluator decides only
  value comparisons over present keys. Why: a domain rule
  held by evaluator discipline plus an optional conformance
  harness is unenforced (the prior draft's own "UNMITIGATED"
  residual); held by the kernel it cannot drift, and the
  refusal can name the absent keys.
- **Selection / predicate** — strong-Kleene combination in
  the kernel: a `GuardFalse` derivable from decided atoms
  prunes (D8); otherwise any unevaluable atom forces
  `GuardUnevaluable`.
- **Seam shape** — per-atom `Evaluate(atom, value)
  GuardResult`; rejected: passing the view (lets the
  evaluator re-decide presence), a row-level seam (moves the
  combination tables out of their normative home), a
  two-valued seam (A18).
- **Naming** — verdicts keep `GuardUnevaluable` /
  `guard_unevaluable`; rejected: a sixth refusal kind
  (reopens RDR 0001's closed taxonomy). The absent-key
  payload is a field on `Refusal`, not a new kind.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Parsed guard atoms on the row | JDR 0001 §D1 | Decided; lands here | `Row.Guard` becomes an atom slice; fixtures migrate |
| Three-valued `GuardResult` + `guard_unevaluable` kind | RDR 0001 (implemented) | Available | Reused unchanged |
| Existence-operator token | RDR 0003 (Draft) | Introduced | Kernel exports the constant; normalizer emits it (A16) |
| Value evaluator over present values | RDR 0003 implementation (future) | Deferred | Narrow seam; value-contract tests (Phase 3) |
| `RequiresOwned` producer | RDR 0002 (Draft, §JD-3) | Deferred | Escape rows carry none (A21) |
| Read completeness | RDR 0004 (§D3) | Decided | Cited; no obligation here |
| Escape-class closure | RDR 0002 | Available | `guard_unevaluable` non-escapable by construction |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Presence predicate | `internal/resolve/resolve.go::TagSet.Lookup` | None for this use (provenance-blind `ok`) | Reuse | The kernel's atom presence test |
| Undecidability verdict + refusal mapping | `resolve.go::GuardUnevaluable`, `resolve.go::gate` | Refusal carries guard text, not keys | Reuse + extend | Payload becomes absent keys |
| Per-row verdict hook | `resolve.go::evaluateGuard` | Delegates whole guard to seam | Extend | Presence, existence, K3 live here |
| Owned-state precheck | `resolve.go::missingOwned` | Owned-provenance only | Reuse | Unchanged; escape rows yield empty |
| Verdict stub for tests | `internal/resolve/fixtures_test.go::fixtureGuards` | Keyed on guard text; discards the view | Replace | Atom-level fixture evaluator |
| Golden vector harness | None under `internal/resolve/` | — | Build (smaller) | Domain-rule vectors are kernel tests; only value-operator contract tests are exported |

### Decision Rationale

The undecided contract: **"when a guard atom references a
tag absent from the assembled view, what verdict results —
and which component enforces that rule now that the kernel
sees every atom?"** Approaches were enumerated as answers to
that question, with the domain rule (absence ⇒ unevaluable)
held fixed per JDR 0001's direction.

Scored matrix (criteria weighted by the user outcome — an
actionable refusal, never a masked one):

| Criterion | A: Evaluator-enforced over parsed atoms | B: Kernel-enforced; per-atom value seam (chosen) | C: Kernel evaluates everything, no seam |
| --- | --- | --- | --- |
| Correctness fit (masking closed?) | ✓ by discipline — a conforming evaluator closes it; a drifting one reopens it silently | ✓ by structure — the evaluator is never asked about an absent key | ✓ by structure |
| Prior-art alignment | ~ — SCXML/K3 satisfied; no host-side precedent claimed | ✓ — SQL:2003 strict routines: host short-circuits, routine "is not invoked" (C3); §D1 "cardinality, not meaning" holds except one token (A16) | ✗ — the kernel would own typed semantics RDR 0003 declares; no declarations in the kernel |
| Drift enforcement (premortem P-10) | ✗ — conformance harness nothing compels 0003 to run; prior draft's UNMITIGATED residual | ✓ — no evaluator code path exists for absence; only value semantics can drift, and those are 0003's to test | ✓ |
| Refusal names the missing key (JDR P4) | ~ — possible via a richer seam return, but the key set is then the evaluator's word | ✓ — kernel decided presence, so the payload is kernel truth | ✓ |
| Blast radius | `Row` + seam signature; `evaluateGuard` unchanged | `Row` + seam + `evaluateGuard` + `Refusal` payload; `gate` order unchanged; ~all test fixtures re-encoded | `Row` + kernel grows operator/kind semantics; 0003's evaluator ownership dissolves |
| Reversibility | ✓ | ✓ — kernel-side enforcement can be relaxed to A by widening the seam; the reverse migration (A → B) costs the same | ✗ — re-separating semantics later is a rewrite |
| Cost | Low | Medium — presence, existence, two conjunctions, one payload; the domain-rule vectors become ordinary kernel tests (no stopgap view-reading evaluator) | High |

Deciding rows: **correctness-by-structure, drift
enforcement, and the refusal payload.** B is the only
approach that closes the masking path without trusting a
component this RDR does not build, and it is the only one
whose refusal meets the Problem Statement's "told plainly
that the artifact state needed to decide was missing" at key
granularity. Its cost over A is bounded kernel logic that the
prior draft would otherwise have built anyway as a
"vector-local, non-normative" view-reading evaluator — B
makes that code the normative one instead of a stopgap. The
one grammar token the kernel learns (A16) is the honest
price of "never call the evaluator on absence"; the
alternative is the evaluator re-deciding presence, which is
A. C is rejected because typed operator semantics and
declared kinds are RDR 0003's and the kernel has no tag
declarations to check literals against.

B does not remove grammar drift; it moves it from the
seam to the kernel/normalizer boundary (token, literal
form, key identity — A16, A22), and that move is accepted
because the direction of failure flips: drift at the seam
under A can yield a plan from absence, while drift at the
kernel boundary under B yields an unevaluable refusal or a
load rejection — fail-closed, never masked. That asymmetry,
not the absence of drift, is the argument.

Premortem: hardened — critic verdict PASS, no switch
forced; "kernel-owned absence is the right trust boundary."
Folded: P-16 (Alternative 1's rejection reason corrected
from capability to trust); P-6/P-9/P-12 (refusal payload
made per-row, per-atom, with reason `absent|uncomparable`,
order-independent); P-3/P-4/P-15/P-17 (kernel-owned grammar
facts — token, literal form, key identity — pinned as
exported constants with the fail-closed direction stated;
A22 added); P-5/P-13 (`unless`-over-absent veto and
pruned-row absences not reported — stated honestly under
Failure Modes / payload clause); P-7/P-8 (seam ceiling
stated as a constraint on RDR 0003 in A17); P-9/P-10
(malformed caller value denial; evaluator skew — Failure
Modes); P-11, P-18 (policy surface; load-time parse —
accepted, one line each); P-14 (the "zero importers" seed
class — re-verify trigger on A3); P-20 (suite must drive
atoms, Phase 1); P-1/P-2/P-19 (observed tags satisfying
presence — already A13; RequiresOwned clause now states the
non-subset rule and the owned-protection route).
Ledger: `docs/rdr/0007-guard-predicate-totality/evidence/propose-premortem/iter-2/critic.md`
(iteration 1 at `…/propose-premortem/critic.md`).

Ground-sweep: clean (24 anchors) — 15 code anchors
(`resolve.go`, test files, `root.go`), 9 peer-document
passages; two cosmetic wording notes fixed inline (JD-3
phrasing; C3 blockquote wrap).

Joint-check: fired → 0003, 0002 (home: OPEN), 0009 (home:
OPEN), 0008, 0004 (home: JDR 0001 §D1). Context beside the
fire: JDR 0001 §D1 settles the *transport* (atoms on `Row`,
the seam change 0008/0004 cite as shipped code) but not the
*enforcement site* — the value-only per-atom seam, the kernel
holding the existence token and literal constants that RDR
0002's normalizer must emit (A16), exact key identity (A22),
and the per-atom refusal payload that removes `Refusal.Guard`
(which Final RDR 0009 cites as a refusal property). Those
need one home: hoist into JDR 0001 as a §D4 (cite-don't-
restate), or declare the seam interface in both 0007 and
0003. Final peers are not edited; their coupling rides to
7.1. Absence arm: n/a (no refusal is converted to an
acceptance). Bridge sub-check: n/a — no sibling plan retires
a surface this plan introduces.

## Alternatives Considered

### Alternative 1: Evaluator-enforced domain rule over parsed atoms

**Description**: Take §D1's atom slice but keep a row-level
seam — `Evaluate(atoms, view) GuardResult` — so RDR 0003's
evaluator applies the domain rule and the strong-Kleene
tables itself; the kernel stays fully grammar-agnostic and
only maps the verdict. This is the prior Final's design
transplanted onto atoms.

**Pros**:

- Kernel learns no operator token; the cleanest reading of
  "cardinality, not meaning."
- Smallest kernel diff: `Row` and the seam signature only.

**Cons**:

- The domain rule is held by discipline: a conformance
  harness that nothing compels RDR 0003's build to run —
  the prior draft recorded this residual as UNMITIGATED.
- The refusal payload could be carried (a seam returning
  `(GuardResult, unevaluable atoms)`), so capability is not
  the objection — **trust** is: whoever decides presence
  decides whether absence can become false, and under A
  that is the component this RDR does not build.
- The normative truth tables live in this RDR but execute
  in another RDR's code.

**Reason for rejection**: fails the drift-enforcement row;
closes masking only conditionally, by discipline.

### Alternative 2: Kernel evaluates every operator; no seam

**Description**: With atoms in hand, implement equality,
membership, comparison, containment, and existence in the
kernel and delete `GuardEvaluator`.

**Pros**:

- Everything this RDR states is one package; no cross-RDR
  conformance question at all.

**Cons**:

- Typed literals and value kinds (bounded integers, set
  universes) are RDR 0003's declarations; the kernel has
  none and would have to grow a declaration model.
- Dissolves RDR 0003's evaluator ownership — a cluster-level
  reassignment this RDR has no mandate for.

**Reason for rejection**: highest blast radius and an
ownership change JDR 0001 did not make.

### Briefly Rejected

- **Opaque guard string + error channel or mandated panic**:
  closed by JDR 0001 §D1 — no reconstruction step, so no
  mapping failure to surface.
- **Total closed-world guard (absence ⇒ false)**: *is* the
  masking path; SCXML permits the fold only with a mandatory
  error event the kernel cannot emit.
- **Declared guard-reads precheck (`Row.Reads`)**: duplicates
  the key set the atoms now carry natively; drift is silent
  unsoundness; reopens ADV-1's poisoning.
- **Narrowed veto (unevaluable row degrades when a decided
  sibling exists)**: absence-as-false at resolution scope.
- **New refusal kind `guard_input_missing`**: reopens the
  closed five-kind taxonomy for a distinction the payload
  carries.
- **Seam receives the view**: invites the evaluator to
  re-decide presence; the point of B is that it cannot.
- **Lint-only totality**: observed-tag presence is a runtime
  fact; RDR 0003's lint already rejects the one static case.

## Trade-offs

### Consequences

- Positive: missing artifact state becomes unmaskable behind
  an escape, and the guarantee is a property of kernel code,
  not of a future evaluator's conformance. The probe under
  *Background* inverts on the shipped kernel once Phases 1–2
  land — kata `xg7p` closes against behavior, not against a
  doc comment.
- Positive: the refusal names the absent keys and the rows.
- Positive: the domain-rule test matrix (operators ×
  presence × blocks × combination) is ordinary kernel
  testing; no stopgap view-reading evaluator and no exported
  domain-rule harness.
- Positive: absence remains expressible via existence atoms.
- Negative: `Row`, the seam, and every test fixture change.
  Bounded by zero production importers (JDR 0001 P7); the
  ADV/Fixup dispositions must be re-encoded, not re-decided
  (A3).
- Negative: tables that assumed closed-world reads refuse
  where they previously pruned or escaped. Bounded — no
  evaluator exists yet — but the reference fixtures
  (`0003-…/evidence/spikes/guard-fixture.toml`,
  `0002-…/evidence/spikes/rdr-fixture.toml`) are authored
  intent and the first migration case; Phase 2 classifies
  each authored guard as safe-or-migration.
- Negative: `guard_unevaluable` is non-escapable by design;
  recovery is fixing the artifact state or accessor, never
  the table.
- Negative: the kernel knows one grammar token (A16).

### Risks and Mitigations

- **Risk**: RDR 0003's value evaluator drifts — e.g. treats
  an unparseable present value as `false`.
  **Mitigation**: partial. The domain rule cannot drift
  (kernel-side). Value-operator semantics can; Phase 3 ships
  exported contract tests for the value seam that RDR 0003
  instantiates. Residual: value-level drift is 0003's to
  test, and A18 bounds its worst case to an unevaluable
  refusal, never a masked plan.
- **Risk**: strong-Kleene combination has a subtle leak
  (`unless` negation, empty blocks).
  **Mitigation**: A2's derivation plus Phase 2's exhaustive
  kernel matrix over {T, F, U}³ across both blocks — now
  runnable, since the combination is kernel code.
- **Risk**: view assembly or a truncated accessor read drops
  state the artifact has.
  **Mitigation**: A6a verifies assembly; read completeness
  is RDR 0004's clause (JDR §D3). Until 0004 implements it,
  operators treat an unevaluable refusal as "inspect the
  accessor," never "retry."
- **Risk**: authors defeat partiality by stamping sentinels
  upstream.
  **Mitigation**: partial. Single-atom absence tests are
  legal (A14); the disjunction pattern rests on A12's
  load-time projection, now homed at §JD-4 with 0003/0006
  able to state it. Authoring guidance names
  sentinel-stamping as the anti-pattern.
- **Risk**: the existence token the kernel pins and the
  token the normalizer emits diverge.
  **Mitigation**: the constant is exported from the kernel
  and the normalizer imports it; a Phase 2 vector asserts an
  existence atom with a foreign token is NOT decided by
  presence (it reaches the seam as a value atom and is
  unevaluable on absence — the safe side).

### Failure Modes

- **Visible break**: a previously escaping table refuses
  `guard_unevaluable`, naming the absent keys and rows.
  Recovery: make the state readable, or rewrite the guard
  with an existence atom if absence was intended.
- **Silent failure guarded against**: absence folded into
  `GuardFalse`. No evaluator code path exists for an absent
  key; the kernel matrix is the tripwire.
- **Review-time blind spot guarded against**: a typo'd key
  is permanently absent; RDR 0002/0003's declared-tag
  rejection catches it at load (A7).
- **Refusal defeated by caller-supplied state** (A13,
  accepted; home §JD-9): an operator can supply the missing
  key as an observed tag. Unreachable until RDR 0005 wires
  `--tag`; Phase 2 pins the guard-path behavior.
- **Seam unevaluable on a present value** (A18): the refusal
  names the row, the key, and reason `uncomparable` — never
  `absent`. Because caller-supplied observed tags are
  unconstrained (A13), one malformed observed value whose
  key any guard mentions can DENY every resolution touching
  it, non-escapably; the payload makes the cause visible
  (`uncomparable`, the key), and the fix is the caller's
  input. Whether lint makes the artifact-side case
  unreachable is RDR 0003's.
- **Evaluator/grammar version skew** (operator token the
  seam does not recognize): the seam answers unevaluable
  and the payload says `uncomparable` on a present key —
  the author is not told a tag is missing, but is not told
  "skew" either. Acceptable under the closed five-kind
  taxonomy; the reason field is what keeps it from reading
  as absence.
- **`unless` over a rarely-set tag** (premortem P-5): an
  `unless` block whose only atom is a value atom over an
  absent key is `¬U = U` and vetoes — correct by the rule
  and a footgun in practice ("skip when `legal_hold` is
  true" refuses on every artifact that never had
  `legal_hold`). The sanctioned idiom conjoins an existence
  atom (`legal_hold exists = true` beside `legal_hold eq
  true` in `unless`: `F ∧ U = F`, so the block is decided).
  Phase 4 hands this idiom to RDR 0003's authoring guidance.
- **Policy surface visible in payloads**: the refusal names
  guard keys to whoever can trigger a resolution. Accepted:
  the table is reviewable data (RDR 0002), not a secret.
- **Load-time parse failure has no refusal kind**: a guard
  the normalizer cannot parse never becomes a `Row`; it is
  RDR 0002's typed validation failure and reaches the user
  through RDR 0005's envelope as a load error, not as a
  resolution refusal.
- **Conflated recovery signal**: until RDR 0004 implements
  §D3, transient read failure and genuine absence are
  indistinguishable at the kernel boundary.

## Implementation Plan

### Prerequisites

- [ ] JDR 0001 §JD-3 landed in RDR 0002 (who populates
      `RequiresOwned`; escape rows carry none — A21).
- [ ] RDR 0003's re-entry cites the §D1 atom shape and the
      kernel's existence token (A16), and states its reading
      of the placement MUST against `exists = false` in
      `all` (A14; pairwise 0007×0003 F2 → §JD-1).
- [ ] A3, A16–A19, A21, A22 verified (Resolve).
- [ ] This RDR Final **before** RDR 0003's implementation
      begins: 0003's evaluator is built against this seam.

### Minimum Viable Validation

Against the real kernel with a real atom — no stub — the
masking probe inverts: one candidate row whose guard is a
value atom over an absent key, an absent `RequiresOwned`
key, and a modeled `no_match` escape row yield
`Refusal.Kind == guard_unevaluable`, a nil `Plan`, the absent
key in the payload, and the row in `Refusal.Rows`; the same
table with the key present and the value decided FALSE
prunes and escapes (D8 preserved). Third scenario:
*unevaluable-blocks-true-sibling* — row A unevaluable beside
row B decided TRUE → refusal naming A, no plan. Fourth: an
`exists = false` atom over the absent key decides TRUE and
the row is selected — the one sanctioned route from absence
to a verdict.

### Phase 1: Row and seam reshape

Replace `Row.Guard string` with the §D1 atom slice; narrow
`GuardEvaluator` to the per-atom value seam; export the
existence token and literal constants; narrow
`Row.RequiresOwned`'s doc contract; migrate every fixture so
the frozen suite drives verdicts THROUGH atoms and a value
stub — never by injecting row verdicts, which would leave the
kernel combinator off the tested path (premortem P-20) — and
re-run it (A3 spike).

### Phase 2: Kernel domain rule and payload

Implement presence, existence, and strong-Kleene combination
in `evaluateGuard`; replace `Refusal.Guard` with the
absent-key payload; pin owned-before-unevaluable within one
survivor set; encode the domain-rule matrix as kernel tests
(operators × present/absent × `all`/`unless` × {T,F,U}
combinations, empty blocks, two-row pattern, escape-set
scoping, foreign existence token); classify each authored
guard in the reference fixtures as safe-or-migration.

### Phase 3: Value-seam contract tests

Export a small contract-test function for the value seam
(present value × literal × operator; unparseable value →
unevaluable, never false) that RDR 0003's implement stage
instantiates against its evaluator.

### Phase 4: Authoring and diagnostics handoff to RDR 0003

Hand over, with 0003 in Draft: the existence token and
`presence == literal` rule; the placement-MUST reading for
`exists = false`; the two-row absence pattern and its
conjoined value row (authoring grammar: may one tag-table
carry two operator keys); the projection of existence atoms
onto the declared-domain product (A12, §JD-4); the
predicate-PLACEMENT guidance (the match pattern is
closed-world by design, so a predicate whose absence must
refuse belongs in guard atoms, never in the match pattern);
sentinel-stamping as the anti-pattern.

## Validation

### Testing Strategy

_Draft placeholder — Resolve re-derives the matrix against
the atom-level seam._ Carried forward as candidate scenarios
from the prior Final: (1) unevaluable-not-no_match,
(2) D8 preserved with the key present, (3)
unevaluable-blocks-true-sibling, (4) strong-Kleene matrix
across `all`/`unless`, (5) escape-set scoping both legs,
(6) `exists` total across both polarity channels, (7)
`contains` over an absent set, (8) two-row absence pattern
plus bare-value-row negative control, (9) unguarded and
empty-block rows — all now kernel-level; plus (10) foreign
existence token, (11) seam-unevaluable on a present value,
(12) owned-before-unevaluable within one survivor set, (13)
guard-path observed-tag substitution (A13).

### Performance Expectations

_Draft placeholder._

## Finalization Gate

Reopened 2026-08-21 — the prior Gate PASS (2026-08-11,
`0007-guard-predicate-totality/artifacts/gate.md`) is voided
by the JDR 0001 re-entry; re-run at Stage 7.

## References

- JDR 0001 — `docs/jdr/0001-resolve-kernel-seam.md` §D1
  (guard seam carries parsed atoms), §D2 (gate, then count),
  §D3 (read completeness), §JD-1/§JD-3/§JD-4/§JD-8/§JD-9
- ISO/IEC 9075 SQL:2003 (5WD-13-JRT-2003-09) p.169 —
  `RETURNS NULL ON NULL INPUT`: "the function itself is not
  invoked" (research C3); PostgreSQL 13 CREATE FUNCTION
  `STRICT`
- W3C SCXML Recommendation §5.9.1 Conditional Expressions
  (research C1)
- PostgreSQL 16 Documentation §9.1 "Logical Operators" —
  strong-Kleene truth tables (A4;
  `evidence/research/resolve-citations.md`)
- RDR 0001 — Resolution kernel (refusal taxonomy; selection
  rule; deviations D5/D8 in
  `docs/rdr/0001-resolution-kernel/artifacts/deviations.md`)
- RDR 0002 — Transition table as reviewable data (escape
  closure; normalization; §JD-3 producer)
- RDR 0003 — Guard predicate exhaustiveness (operator
  vocabulary and semantics; placement rule)
- RDR 0004 — Accessor execution safety model (gate
  indeterminate house rule; §D3 read completeness)
- RDR 0009 — Escape-row shape (empty `Writes` on escape rows)
- `internal/resolve/resolve.go` (`Row`, `GuardEvaluator`,
  `GuardResult`, `TagSet.Lookup`, `assemble`, `gate`,
  `evaluateGuard`, `missingOwned`, `escapeOrRefuse`,
  `Resolve`, `Refusal`)
- `internal/resolve/adversarial_test.go`,
  `internal/resolve/fixup_test.go`,
  `internal/resolve/fixtures_test.go` (frozen dispositions
  and the text-keyed stub Phase 1 replaces)
- Prior-art cache:
  `docs/rdr/0007-guard-predicate-totality/evidence/research/propose-prior-art.md`
- Stage 4 (iteration 1) aggregation spike:
  `docs/rdr/0007-guard-predicate-totality/evidence/spikes/aggregation-probe.md`
- kata `xg7p` (originating finding)
