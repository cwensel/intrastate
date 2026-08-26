# Recommendation 0007: Guard predicate totality over an incomplete evaluation view

> Revise during planning; lock at implementation.
> If wrong, abandon code and iterate RDR.

## Metadata

- **Date**: 2026-08-09
- **Status**: Implemented
- **Type**: Architecture
- **Profile**: foundational — cross-RDR producer at the
  0001↔0003 seam, consumed by RDR 0003's value evaluator and
  RDR 0002's normalizer. Three coupled contracts, all
  carried by the single `Row` reopening JDR 0001 §D1 grants:
  the guard-evaluation domain rule and its kernel
  enforcement site; the per-row/per-atom `guard_unevaluable`
  payload replacing `Refusal.Guard`; and the narrowed
  meaning of `Row.RequiresOwned` (kata `xg7p`). Not split:
  none can lock without the same type change, so splitting
  yields RDRs that cannot lock independently.
- **Priority**: Medium
- **Related Issues**: kata `xg7p` — "RDR 0001:
  RequiresOwned conflates guard-input state with post-guard
  transition state" (`src:roborev`, `severity:medium`,
  `area:internal-resolve`).
- **Predecessors**: 0001-resolution-kernel,
  0003-guard-predicate-exhaustiveness; JDR 0001
  (`docs/jdr/0001-resolve-kernel-seam.md`) §D1, §D2, §D3.
- **Overrides**: None superseded. This RDR carries the
  `resolve.Row` change JDR 0001 §D1 assigns to it (the guard
  seam carries parsed atoms, not a string), *fixes* the
  meaning of `Row.RequiresOwned` (post-guard write-dependency
  keys), and is the single normative home of the
  guard-evaluation domain rule and its enforcement site.
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
a tag absent from the assembled evaluation view, nothing
stops the answering component from returning a two-valued
verdict — treating absence as non-satisfaction — rather than
refusing to answer.
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

Constraints this rule must hold. RDR 0001's
deviation **D8** (guard-FALSE prunes first; a pruned row
contributes neither candidacy nor an owned-state
obligation) is pinned by frozen tests ADV-1 and ADV-1b.
Deviation **D5** scopes `RequiresOwned` and guard evaluation
to candidate rows only. RDR 0002 closes escape lists to
`no_match` and `ambiguous_match`, which is the asymmetry that
makes verdict ordering load-bearing.

**JDR 0001 decisions this RDR inherits.** Three cluster
decisions fix the ground this RDR builds on:

- **§D1** — the row carries **parsed atoms** (key, operator
  token, literal, block); no reconstruction step, so no
  mapping failure exists and no panic or error channel is
  needed for one. This RDR takes the `Row` change.
- **§D2** — **gate, then count**: guards evaluate over
  surviving candidates, an unevaluable candidate vetoes,
  and only then exact-one matching and escape reachability
  apply. Stated here as the kernel's ordering; RDR 0002
  restates its resolver flow to match.
- **§D3** — a read accessor MUST return the complete tag
  set or refuse. RDR 0004 carries the clause as its own;
  this RDR cites it (A6b).

RDR 0003 and RDR 0002 are both **Final** and cite the §D1
atom shape rather than restating it. 0003 owns operator
vocabulary and typed semantics; its evaluator is a *value*
evaluator under this RDR.

### Technical Environment

Go; the resolution kernel `internal/resolve` — one non-test
file, `resolve.go::Resolve` as the single pure entry point,
no production importer (only in-package tests), which is why
JDR 0001 P7 prefers the clean shape over a compatible one.
The contract sits at the boundary between:

- **RDR 0001 — Resolution kernel** (`Implemented`). Owns
  the closed five-kind refusal taxonomy — `no_match`,
  `ambiguous_match`, `owned_state_unavailable`,
  `guard_unevaluable`, `unmodeled_outcome`
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
  refusals; **RDR 0006** (`Final`, locked 2026-08-23)
  narrows lint's promise where this RDR's veto and its proof
  disagree — the narrowing is recorded by RDR 0003
  (JDR 0001 §JD-4, closed 2026-08-22) and cited by RDR 0006.

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

## Critical Assumptions

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
  - **Status**: Verified — re-spiked against the SPECIFIED
    shape (`Refusal.Guard` deleted, view-free per-atom seam,
    explicit existence literals). The bound holds:
    `escapeOrRefuse` has ZERO diff, `gate`'s prune →
    owned-state → undecidable order is byte-identical,
    `boundary_test.go` and `mvv_test.go` are untouched, and
    the frozen suite is 154 PASS / 0 FAIL.
    **`Fixup1d` is a RE-DECIDE, as this assumption predicted**,
    and for a deeper reason than the deleted field: its
    fixture guard `iterations >= 3` names a key the view does
    not carry, so under the reshape the KERNEL decides that
    atom unevaluable on ABSENCE (captured verbatim:
    `Reason:absent`) and the nil seam is never consulted — the
    property the test's NAME asserts is gone, and it would go
    green while testing nothing about the nil seam. The SEAM
    clause's prescribed fix (move the atom to the PRESENT key
    `reviews`, expect `uncomparable`) restores it; Phase 1
    carries it. Two other tests needed compile-forced
    re-encodes only (`TestReq33` non-emptiness check;
    `TestAdv3`'s `t.Errorf` format string).
    Coverage gap CONFIRMED and sharpened: every frozen guard
    is a single atom, so the frozen suite cannot specify the
    combinator. A FALSE-vs-UNEVALUABLE dominance mutant
    survives all 154 and is killed only by the multi-atom
    probe — Phase 1 must add combination tests.
    Verified independently of the shape: the touch-list is
    additive beyond the five names above — the reshape also
    introduces `GuardAtom`, `evaluateAtom` (where the
    presence/existence rule lands), and
    `truthRank`/`kleeneAnd`/`kleeneNot`/`sortUndecided`/
    `compareUndecidedAtoms` to keep the slice-valued payload
    tuple-deterministic under REQ-1 (proposed names; none
    exists under `internal/` today). Note the "replace" shape
    needs NO copy helper: the payload is built fresh from the
    verdict pass rather than echoing the row's guard slice, so
    it cannot alias the caller's table.
    The "zero production importers" bound was a point-in-time
    fact, so it carried a re-verification trigger; the sweep
    was re-run this stage (cove iteration 2) and is still
    CLEAN — 154 PASS / 0 FAIL, no production importer — so
    the `Row` change is a free change, not a coordinated
    migration. The trigger stands for Phase 1 start.
  - **Method**: Spike
  - **Evidence**: Re-spike record (AUTHORITATIVE — the
    specified shape)
    `docs/rdr/0007-guard-predicate-totality/evidence/spikes/a3-reshape-respike.md`,
    with the proved kernel and probes captured under
    `evidence/spikes/a3-reshape/respike_*.txt`. It supersedes
    the first spike record
    `evidence/spikes/a3-reshape-probe.md`, which proved a
    SUPERSET shape (it kept `Refusal.Guard`, passed the view
    to the seam, and used a literal-less existence atom) and
    whose existence case is confirmed WRONG against this
    draft. Both are kept: the first is the record of what the
    re-spike had to correct.
    (Reshape applied to a scratchpad copy; `internal/`
    untouched, `go vet ./...` exit 0.) Baseline 154 PASS / 0
    FAIL; post-reshape 154 PASS / 0 FAIL with identical
    dispositions (180 PASS with the probes added).
    `escapeOrRefuse` needed ZERO changes and `gate`'s
    prune → owned-state → undecidable partition is
    byte-identical — its only hunk builds the refusal
    payload. Importer sweep
    `rg -n -t go --glob '!internal/resolve/**' 'newcoinc/intrastate/internal/resolve' .`
    returns no production importer; `internal/cli/root.go::NewRootCmd`
    registers only `newVersionCmd()`. Inspection: `resolve.go::gate`
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
    guard's representation. All six pass re-encoded, and NO
    frozen contract had to be reopened: the only two sites
    reading `Refusal.Guard` assert no guard-text identity —
    `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
    compares against the value the test itself set, and
    `resolve_test.go::TestReq33_UnavailableOwnedStateAndUnevaluableGuardAreValueRefusals`
    is a non-emptiness check.
  - **Coverage gap (mutation-verified, load-bearing for
    Phase 1)**: the frozen suite proves "dispositions
    unchanged" but does NOT cover the new combinator. A
    mutation swapping strong-Kleene FALSE-dominance for
    UNEVALUABLE-dominance in `evaluateGuard` SURVIVES all
    154 frozen tests, because every frozen guard is a single
    atom and a one-atom conjunction cannot observe
    combination order. Phase 1 MUST add multi-atom
    combination tests as genuinely new surface; the 6-case
    probe at `evidence/spikes/a3-reshape/spike_kleene_test.go.txt`
    kills the mutant and is the shape to start from **for the
    conjunction cases only**. Its existence case is WRONG
    against this draft and MUST NOT be copied: the probe
    asserts "an exists atom over an absent key is DECIDED
    false, so the row is pruned" from an atom
    `{Key: "guard.missing", Op: resolve.OpExists}` carrying no
    `Literal`. Under the existence clause the verdict is
    `presence == literal`, so an absent key decides TRUE for
    `exists = false` and FALSE for `exists = true`; and a
    missing literal is `uncomparable`, never presence-decided
    (A24). Phase 1 rewrites that case to carry an explicit
    `LiteralTrue`/`LiteralFalse` and asserts both polarities
    (Testing Strategy row 7).
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
  - **Method**: Peer RDR
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
  - **Status**: Verified (kernel half); the agreement with RDR
    0003's subtractive lint algebra is DOWNGRADED — 0003
    states no empty-block identity, so the concord is an
    inference (see *Downgraded* below).
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
    subtracted from the row's accepted assignments") is
    COMPATIBLE with the term-dropping reading but does not
    state it: RDR 0003 mentions the empty or omitted `unless`
    case zero times, so the agreement is this RDR's inference,
    not 0003's clause. Read literally over an empty block,
    0003's subtraction takes "the single conjunctive
    assignment set matched by the row's full `unless` block" —
    and an empty conjunction matches the whole scoped product,
    which would subtract everything and disable the row. That
    is the same wrong identity this clause rejects, arrived at
    from 0003's side. The two algebras therefore agree only
    once 0003 states the empty-block identity. Both canonical
    fixtures omit `unless` on most rules; the 0002 spike
    normalizer contributes zero predicates for an absent
    block. Now kernel-implemented, so Scenario 9(b) becomes a
    kernel test.
  - **Downgraded**: the kernel half (this RDR's own verdict
    rule) is Verified; the cross-document agreement with RDR
    0003's lint algebra is NOT, and rides to 0003's re-lock as
    a Prerequisite. If 0003 reads its subtraction literally
    over an empty block, lint and the kernel select different
    row sets — lint proves a group exhaustive and the kernel
    returns `no_match`, which is JDR 0001 P5's promise
    ("a green lint means resolution succeeds") failing in the
    direction that makes lint worthless.
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
  - **Status**: DOWNGRADED — RDR 0003 is silent on the
    projection, so the question is open, not refuted; an
    inherited obligation with home JDR 0001 §JD-4 (lint's
    promise narrows), re-confirmed open this stage (§JD-4 is
    not among the Closed entries; its SUBSTANCE is already
    decided — "the *promise* narrows — P5 decides" — leaving
    open only which document records the narrowing and whether
    lint gains a warning category).
    Survivable, and not this RDR's to close: the projection
    rule belongs to RDR 0003/0006, both Draft. Named plan: it
    lands when 0003 states its existence-atom projection, and
    rides that document's lock rather than this one. Nothing
    in the kernel's domain rule depends on the answer — as
    "If wrong" records, the rule survives and only the
    AUTHORING story for the disjunction case degrades.
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
  - **Status**: Verified
  - **Method**: Source Search (RDR 0003 vocabulary and RDR
    0002 normalization) + Design Decision
  - **Evidence**: RDR 0003's operator/kind matrix spells
    presence as exactly ONE token — the `exists` row ("Tests
    presence or absence, not value equality", literal shape
    `boolean`); there is no `absent` counterpart token,
    absence being the `false` literal. The other four rows
    (`eq`, `in`, `lt`/`lte`/`gt`/`gte`, `contains`) are
    stated purely as domain-narrowing value comparisons, and
    RDR 0003's vocabulary clause closes the set at five
    operator CLASSES (eight tokens — bounded comparison
    spells `lt`, `lte`, `gt`, `gte`).
    JDR 0001 §D1 fixes the atom as "(key, operator token,
    literal, block)" and says the kernel "carries
    *cardinality*, not *meaning*"; §D4(b) prices the
    exception exactly — the kernel learns "one grammar fact
    (the existence operator token and its two boolean literal
    forms), exported as constants 0002's normalizer MUST
    emit." RDR 0003 and RDR 0002 each carry the identical
    joint-check line naming that kernel-exported token.
    The KERNEL half — one token, two literal forms, exported
    as constants — is fully verified and is what this RDR
    owns. The PRODUCER half is not yet a duty on 0002: its
    Refinement Context Direction list names the
    `RequiresOwned` producer and the fixture rename, and
    carries neither the existence-token emission duty nor
    A22's canonicalization duty. The joint-check line is an
    acknowledgement, not a clause; both duties are added to
    0002's direction list by the same Prerequisite.
  - **If wrong**: the kernel cannot separate `exists` from
    value operators and the enforcement site falls back to
    the evaluator (Alternative 1) for that distinction.
- **A17 A value operator over a present key is decidable
  from the atom and the present value alone: the seam
  `Evaluate(atom, value)` needs neither the view, nor
  provenance, nor sibling atoms, nor tag declarations at
  evaluation time.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: RDR 0003's operator semantics — "equality
    compares a tag value to one typed literal; membership
    checks a scalar tag against a typed literal set; bounded
    integer comparison uses `lt`, `lte`, `gt`, and `gte`;
    existence checks presence of an optional tag value";
    set containment "Narrows
    the set-valued domain to assignments containing every
    listed element" — each reads one value against one
    literal. The declared kind needed to *parse* the value
    is known to the evaluator from the table it was built
    for, not from the kernel. The "bounded" worry is closed
    by RDR 0003's Technical Design: "A guard may still
    compare an unbounded integer at runtime, but lint must
    report that it cannot prove exhaustive coverage for that
    dimension" ⇒ declared bounds are a LINT input, never an
    evaluation input. `contains`'s declared element universe
    is likewise a lint-proof input, and RDR 0003's semantic
    identity tuple is `(tag, operator, literal)` with no kind
    field. The seam signature holds as stated: the atom need
    not carry declared kind or bounds.
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
  - **Status**: Verified
  - **Method**: Design Decision + Source Search (RDR 0003's
    parse/lint rules for runtime values)
  - **Evidence**: An unparseable present value is the same
    class as the nil-seam case `resolve.go::evaluateGuard`
    already answers `GuardUnevaluable`: the evaluator cannot
    decide, and collapsing that into `false` would be the
    masking path one layer down. Whether RDR 0003 treats it
    as a lint-time impossibility (declared kinds validated
    at load) or a runtime verdict is 0003's; the kernel
    accepts the verdict either way. That deferral is
    confirmed correct, not merely polite: RDR 0003's only
    parse-rejection clause is about the authored LITERAL ("A
    predicate whose literal cannot be parsed as the declared
    tag kind MUST be rejected before resolution") and its
    Testing Strategy lists "literal parse mismatch" among
    PARSE-time scenarios only — so 0003 does not make an
    unparseable RUNTIME value a load-time impossibility, and
    the gap this assumption anticipates is real. The shipped
    kernel already names the seam a legitimate reporter of
    undecidability: `resolve.go::KindGuardUnevaluable` is
    documented "the guard seam reported a predicate it could
    not decide".
  - **If wrong**: the seam must be two-valued and a
    mis-typed runtime value becomes a programmer-error
    return from `Resolve` — a kernel contract change.
- **A19 The `guard_unevaluable` refusal can name the absent
  tag keys directly (the kernel now knows them). Carrying
  that payload to the CLI as STRUCTURE requires ONE new
  optional envelope field on RDR 0005's failure envelope —
  the original "without a new envelope field" reading is
  refuted.**
  - **Status**: DOWNGRADED — blocked on a peer-document
    change, not on evidence; nothing further is verifiable
    here before lock. The kernel half is Verified.
    Survivable, and the fallback is named and costed below:
    if the reopening is declined, the payload flattens into
    `Detail` exactly as §D4 originally ruled. The direction is
    reversible EITHER way — an `omitempty` field is
    append-only and breaks no existing consumer — so lock does
    not foreclose the choice. Named plan: the reopening runs
    as a Prerequisite against JDR 0001, and blocks RDR 0005's
    JSON rendering only, not this RDR's implementation.
    **This is a REVERSAL of a closed joint decision, not an
    open question.** JDR 0001 §D4 already routed this payload:
    it "reaches the CLI through §JD-8's `Detail`" — flattened
    text — and §JD-12 records §D4 as Closed. The author's
    2026-08-21 round chose the option §D4 rejected. This RDR
    must therefore reopen §D4, not merely amend §JD-8.
    Why the reversal is nonetheless right, on engineering
    grounds independent of precedence: `clierr.CLIError`'s own
    type doc invites exactly this extension — "Extend with new
    optional fields as needed — keep them `omitempty` so the
    envelope stays append-only and stable for tools" — and an
    `omitempty` field is invisible to every existing consumer.
    `Detail` is documented human prose ("the underlying syscall
    reason, a parser diagnostic"); serializing a two-level
    sorted array into it puts a second, undocumented encoding
    inside a string, which defeats the total-order
    determinism whose only purpose is machine consumption.
    Further, §JD-8's "not new envelope fields" is a
    SUFFICIENCY claim about the `Code` values RDR 0005
    enumerated — inherited verbatim from 0005's A5 evidence,
    written before this payload existed — not a prohibition.
    Named plan: reopen §D4 at JDR 0001 with the §JD-8
    amendment granting ONE `omitempty` structured field
    (scoped to the machine-readable refusal payload; no new
    exit group), then this assumption cites that decision and
    flips to Verified. RDR 0005 is `Final`, so 0007 MUST NOT
    unilaterally widen its envelope.
    **Fallback if the reopening is declined**: flatten the
    payload into `Detail` as §D4 originally ruled, and record
    here that the JSON envelope carries the payload as text
    only — machine consumers re-parse it. The choice is
    reversible in both directions (an `omitempty` field is
    append-only and breaks nothing later), which is why this
    does NOT need to block lock; see Prerequisites.
  - **Method**: Source Search (RDR 0005 envelope; JDR 0001
    §JD-8; `internal/cli/clierr` source)
  - **Evidence**: The refutation is in shipped source:
    `internal/cli/clierr/clierr.go::CLIError` carries only
    flat strings — `Code`, `Message`, `Param`, `Detail`,
    `Hint` — where `Detail` is a `string` documented "May be
    multi-line", i.e. multi-line TEXT, not a nested object.
    RDR 0007's payload contract needs a two-level
    array-of-objects (per undecidable row its
    `(RuleID, SourceLocator)`; per unevaluable atom its key,
    block, and reason), sorted on the total tuple. RDR 0005's
    normative split puts structure on the SUCCESS side only
    ("Failures MUST use the existing CLIError JSON/text
    envelope defined by `docs/cli-output-contract.md` and
    `internal/cli/clierr`"), so the structure this RDR needs
    is on the wrong side of that split. JDR 0001 §JD-8's
    pre-authorization ("only require new `Code` constants or
    literals, **not new envelope fields or exit groups**")
    therefore does NOT cover this payload.
    What IS pre-authorized and needs nothing further: the
    `Code` values. `guard_unevaluable` already has a row in
    RDR 0005's stable-code table (`flow-guard-unevaluable`,
    `GroupUserEnv`) — though scoped there to "supplied
    facts", which the cluster gate already flagged against
    this RDR's provenance-blind view; `owned_state_unavailable`
    has NO row, a gap §JD-8 names explicitly ("need `Code`
    values"). The kernel-side field replaces
    `Refusal.Guard string` (one row's guard text —
    meaningless once guards are atoms); `Refusal.Rows` keeps
    row identity.
    The chosen shape (author's round) is one `omitempty`
    structured field; `CLIError`'s own doc comment invites
    exactly this — "Extend with new optional fields as
    needed — keep them `omitempty` so the envelope stays
    append-only and stable for tools" — so the cost is in
    JDR consistency, not in code.
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
  - **Status**: Verified (composition); the PRODUCER question
    stays open at §JD-3 and is carried as a Capability
    Dependency, not as a claim of this RDR.
  - **Method**: Source Search (RDR 0009 Final; RDR 0002's
    §JD-3 producer clause)
  - **Evidence**: RDR 0009 Normative Contracts, quoted with
    its governing scope: escape-row shape conformance "is a
    PRODUCER obligation on every constructor of `resolve.Row`
    values: a Row with a non-empty Escape list MUST have an
    empty Writes slice." That clause binds producers — but it
    is NOT 0009's only enforcement: a separate 0009 normative
    block puts the same obligation on the KERNEL as an entry
    precondition ("The kernel enforces the same obligation as
    an entry precondition of Resolve … A table containing a
    row that breaches the conformance predicate above MUST
    cause Resolve to return a non-nil Go error identifying the
    offending row by RuleID and SourceLocator, with no Result
    disposition"), evaluated over the WHOLE table before any
    evaluation step and travelling the Go error path, never a
    refusal kind. So the composition holds for conforming rows
    AND a non-conforming escape row is rejected at `Resolve`
    entry rather than reaching the gate.
    JDR 0001 §JD-3 names this the open composition (0007
    derives from `Writes`, 0009 empties `Writes`), notes the
    field "appears **zero** times in 0002, which owns the
    normalized row", and calls the D1 `Row` reopening "the
    natural occasion to settle it"; RDR 0002's re-entry
    direction takes the producer question. Under this RDR's
    definition the composition yields the empty set, so an
    escape row raises no `owned_state_unavailable` of its
    own — confirmed at the kernel: `resolve.go::missingOwned`
    ranges over `row.RequiresOwned`, so an empty slice never
    enters the loop and `missing` stays nil. RDR 0009's
    clause is verbatim as quoted and is `Final`; note it
    mentions `RequiresOwned` ZERO times, so it supplies only
    the `Writes`-empty premise and neither supports nor
    contradicts the derivation. §JD-3 remains NOT closed
    (unlike §JD-1/§JD-2/§JD-6/§JD-7/§JD-12), and RDR 0002
    carries the producer only as a re-entry Direction ("Name
    the `RequiresOwned` producer"), not yet a clause. The
    composition holds regardless of who populates the field;
    the deferral must survive into 0002's re-lock or the
    derive-from-`Writes` definition has no implementer.
  - **Fixture collision with 0009's precondition (A25).** The
    shipped `fixtures_test.go::escapeRow` sets
    `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}`
    on a row with a non-empty `Escape` list, so it breaches
    0009's kernel precondition on its face. `escapeRow` is the
    constructor for every escape fixture in the frozen suite,
    so once 0009 lands, ADV-1b, ADV-2, Fixup-1d and
    Fixup-1e resolve through a non-nil error rather than the
    dispositions they assert. Phase 1 re-encodes those
    fixtures for atoms and MUST repair `escapeRow` in the same
    pass — otherwise an implementer cannot tell whether 0007's
    reshape or 0009's precondition broke them. Which RDR
    sequences first is a cluster question this RDR does not
    settle; it records the collision (A25).
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
  - **Status**: DOWNGRADED — kernel half Verified; the
    producer half is an UNSTATED INFERENCE against RDR 0002's
    current text, re-confirmed this stage (the token
    `canonicaliz` appears ZERO times in RDR 0002, and its
    Refinement Context Direction list names only the
    `RequiresOwned` producer (JD-3) and the fixture rename
    (JD-10)). Survivable: the failure direction is CLOSED —
    a spelling the normalizer lets through inconsistently
    yields an unevaluable refusal, never a plan — and the
    fallback below (narrow A22 to a kernel-side input
    precondition, which is fully verified today) requires no
    code change, only a change in where the obligation is
    filed. Named plan: bind it at RDR 0002's re-lock
    as one normative clause (tag keys are canonicalized —
    naming case, namespace, whitespace — during
    normalization, so a `Row` carries only canonical keys),
    and add that duty to 0002's Refinement Context Direction
    list, which does not currently carry it. Fallback if 0002
    declines: narrow A22 to kernel-side only (the kernel
    performs no canonicalization and requires canonical keys
    as an INPUT PRECONDITION), which is fully verified today,
    and re-file canonicalization as an obligation on 0002
    rather than an assumption of 0007.
  - **Method**: Source Search (RDR 0002 tag declaration and
    normalization contracts; `resolve.go::assemble`)
  - **Evidence**: `resolve.go::assemble` keys `view.tags` by
    `Tag.Key` verbatim and `TagSet.Lookup` indexes the map
    by the string it is given; the kernel has never
    canonicalized — verified. The producer half does NOT
    hold against 0002's text: RDR 0002 has no
    canonicalization clause and no tag-name grammar or
    charset rule anywhere; its declaration-completeness rule
    ("The model MUST declare every tag it matches or writes,
    including each tag's provenance") is a DIFFERENT rule,
    and fixes spelling only under the unstated premise that
    declaration lookup is itself exact — if 0002's normalizer
    ever matched declarations case-insensitively, two
    authored spellings would resolve to one declaration yet
    reach the kernel as distinct `Tag.Key` strings, which is
    precisely this assumption's "If wrong" branch. The only
    supporting sentence is JDR 0001 §D4's landing note ("the
    normalizer ... canonicalizes key spellings before a row
    exists") — an obligation FILED on 0002, which 0002 has
    not received as a clause. Raised by premortem P-17.
  - **If wrong**: a spelling the normalizer lets through
    differently in a guard and in a write is "absent" to the
    kernel — an unevaluable refusal, never a plan (fail
    closed), but one that misdiagnoses a producer defect as
    missing state.
- **A23 The named payload surface — `Refusal.Undecided
  []UndecidedRow`, its `UndecidedAtom` element, and the closed
  `Reason` constant set — is additive on `Refusal` and reopens
  no frozen kernel contract: no shipped test asserts `Refusal`'s
  field set exhaustively, and the closed-set pattern it copies
  (`RefusalKind`/`RefusalKinds()`) is the kernel's own.**
  - **Status**: Verified — closed by A3's re-spike, which
    landed the surface with `Refusal.Guard` actually deleted.
  - **Method**: Spike (folded into A3's re-spike — the same
    reshape with `Refusal.Guard` actually removed is where
    these names land)
  - **Evidence**: `evidence/spikes/a3-reshape-respike.md`
    §Q-A23. `Refusal.Undecided` lands as a bare field and
    `UndecidedRow` / `UndecidedAtom` / `Reason` / `Reasons()`
    as new top-level surface, with NO frozen contract
    reopened. Three grounds: (1) no test constructs a
    `Refusal` composite literal at all (`grep 'Refusal{'` over
    `internal/resolve/*_test.go` → zero hits), so no test can
    break through that route; (2) all four
    `reflect.DeepEqual`-over-`Refusal` sites
    (`adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`,
    `::TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`;
    `fixup_test.go::TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder`,
    `::TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder`)
    are SELF-comparisons of forward vs. reversed row order —
    structurally agnostic to the field set, and they needed
    zero edits. They are also what would catch a
    non-deterministic payload, so their passing is independent
    evidence that `sortUndecided` makes `Undecided` a function
    of the input tuple; (3)
    `resolve_test.go::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`
    PASSES UNEDITED — it pins the KIND set, and the new
    `Reason` closed set is a separate type that never enters
    `RefusalKinds()`. The payload is a field, not a sixth
    kind, exactly as the Load-Bearing Decision states.
    Only two tests needed edits and both were compile-forced
    by the deletion, neither by `Refusal`'s shape being
    frozen.
  - **If wrong**: a frozen test does constrain `Refusal`'s
    shape, and the payload lands as a nested type behind an
    accessor rather than a bare field.
- **A24 The kernel can compare an existence atom's literal
  against two exported constants without parsing: the literal
  reaches the kernel as an already-typed boolean-shaped token
  the normalizer emitted, so verbatim comparison is total over
  well-formed input and every other value is `uncomparable`.**
  - **Status**: Verified (kernel half) — closed by A3's
    re-spike. The PRODUCER half rides on the same unlanded
    0002 duty as A16 and A22 (Prerequisites), and fails
    CLOSED if it never lands, which is why it does not block.
  - **Method**: Spike (folded into A3's re-spike) + Design
    Decision
  - **Evidence**: `evidence/spikes/a3-reshape-respike.md`
    §Q-A24. Verbatim comparison against the two exported
    constants is TOTAL over well-formed input, and the kernel
    needs no literal-parse step: implementation is
    `switch atom.Literal` over `LiteralTrue`/`LiteralFalse`
    with `default:` → `(GuardUnevaluable, ReasonUncomparable)`
    — no parsing, case-folding, or coercion.
    Both polarities asserted with an EXPLICIT literal, all
    four PASS: `exists=true` over an absent key decides FALSE
    (row pruned); `exists=false` over an absent key decides
    TRUE (row selected); both present-key mirrors hold. Each
    case also asserts the seam was NEVER consulted, per the
    clause's "the evaluator MUST NOT be consulted for it".
    The fail-closed leg PASSES over six cases — three literal
    shapes (`""` zero-value, `"True"` case-drift, `"yes"`
    foreign) × present/absent key — each yielding
    `guard_unevaluable` with reason `uncomparable`, asserted
    by `reflect.DeepEqual` over the full payload. Critically,
    the absent-key/missing-literal case does NOT decide TRUE,
    which is what treating a missing literal as `LiteralFalse`
    would have produced. The first spike's literal-less
    existence atom is confirmed WRONG against this draft and
    is superseded.
  - **If wrong**: the kernel needs a literal-parse step, which
    reopens "cardinality, not meaning" wider than the one
    grammar fact §D4(b) prices.
- **A25 RDR 0009's kernel entry precondition and this RDR's
  fixture migration can be sequenced without either reopening
  the other's frozen contracts.**
  - **Status**: DOWNGRADED — the collision is Verified on
    source (re-confirmed this stage: `fixtures_test.go::escapeRow`
    sets `Writes: []resolve.Tag{{Key:"status", Value:"Blocked"}}`
    AND `RequiresOwned: []string{"status"}` beside a non-empty
    `Escape` list, breaching 0009's predicate on both fields).
    What stays open is only the SEQUENCING, which is a cluster
    decision this RDR cannot settle alone — it is 0009's
    re-lock that must order itself against this RDR's Phase 1,
    and the Prerequisites carry it as an explicit 0009 item.
    Named plan: settle the order at cluster-reconcile
    (`/rdr-cluster-reconcile`) before either implementation
    starts; whichever lands second re-runs a suite the first
    already re-encoded, so the cost of getting it wrong is
    rework, not a wrong contract.
    Survivable because both migrations rewrite the SAME four
    frozen dispositions (ADV-1b, ADV-2, Fixup-1d, Fixup-1e)
    toward the SAME end state — an escape row with empty
    `Writes` and empty `RequiresOwned`. Neither changes what
    the fixtures must become; only who edits them first.
    A3's re-spike further narrows the exposure: it migrated
    `escapeRow` and ran all four to green without touching
    `boundary_test.go` or `mvv_test.go`, so the migration is
    known to be mechanical.
  - **Method**: Source Search (RDR 0009 kernel-precondition
    block; `internal/resolve/fixtures_test.go`)
  - **Evidence**: 0009 (`Final`) states "The kernel enforces
    the same obligation as an entry precondition of Resolve …
    MUST cause Resolve to return a non-nil Go error", over the
    whole table, before any evaluation. `fixtures_test.go::escapeRow`
    builds every frozen escape fixture with
    `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}`
    beside a non-empty `Escape` list, so it breaches that
    predicate. ADV-1b, ADV-2, Fixup-1d and Fixup-1e all route
    through `escapeRow`. Both changes therefore rewrite the
    same fixtures, and whichever lands second re-runs a suite
    the first already re-encoded.
  - **If wrong**: the two migrations collide mid-phase and an
    implementer re-decides a frozen disposition under schedule
    pressure — the failure A3's "re-encode, don't re-decide"
    discipline exists to prevent.
- **A26 The new exported kernel surface clears RDR 0001's
  frozen boundary-symbol tests: none of `GuardAtom`,
  `UndecidedRow`, `UndecidedAtom`, `Undecided`, `OpExists`,
  `LiteralTrue`, `LiteralFalse`, `ReasonAbsent`,
  `ReasonUncomparable`, `Reasons`, `Block`, `BlockAll`,
  `BlockUnless`, or `TestGuardEvaluatorContract` collides with
  a banned name.**
  - **Status**: Verified — closed by A3's re-spike. All three
    frozen boundary tests PASS against the compiled surface,
    and no proposed name collides under EITHER matching rule
    (see the correction below).
  - **Method**: Spike (the re-spike compiles the surface, then
    runs REQ-25/REQ-36/REQ-37)
  - **Evidence**: `internal/resolve/boundary_test.go::exportedKernelSymbols`
    walks every exported top-level name; three frozen tests
    match against banned substring lists —
    `resolve_test.go::TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence`
    (`Run`, `Execute`, `Apply`, `Persist`, `Save`, `Commit`,
    `Start`, `Orchestrate`),
    `::TestReq36_KernelExposesNoEncodeDecodeOrInverseOperation`,
    and `::TestReq37_KernelIntroducesNoHashOrCanonicalSerialization`
    (`Hash`, `Checksum`, `Canonicalize`, `Fingerprint`,
    `Digest`). SWEEP RUN (`evidence/spikes/a3-reshape-respike.md`
    §Q-A26): all three PASS against the compiled surface, and
    a standalone sweep clears every one of the fourteen names
    against the union of all three banned lists.
    **Correction — these are EXACT-match lists, not substring
    lists.** All three tests compare `if n == b` over
    `exportedKernelSymbols(t)`; nothing tests containment. The
    re-spike swept under both readings and found no collision
    either way, so the verdict is unchanged — but the
    consequence for A22 is weaker than previously stated: only
    the bare name `Canonicalize` is banned, so a kernel-side
    `CanonicalizeKey` would PASS REQ-37. That constraint on
    A22's fallback is therefore a house-style preference, not
    a frozen-test prohibition.
    **Scope note.** `Undecided` is a FIELD of `Refusal`, not a
    top-level declaration, and `exportedKernelSymbols` walks
    top-level decls only — it never descends into struct
    fields. The three frozen tests are structurally incapable
    of seeing it, so its clearance rests on the standalone
    sweep, not on REQ-25/36/37. The same holds for every field
    name on the new payload types (`Key`, `Block`, `Operator`,
    `Literal`, `Reason`, `Atoms`, `RuleID`, `SourceLocator`).
    `TestGuardEvaluatorContract` is Phase 3 surface, swept as
    a name only; it also lives in a `_test.go` file, which
    `parseKernelPackage` filters out.
  - **If wrong**: a normative name this RDR fixes by spelling
    collides with a frozen boundary test, and Phase 1 must
    either rename the normative surface or reopen an RDR 0001
    contract.
- **A27 The literal byte values this RDR now spells —
  `OpExists = "exists"`, `LiteralTrue = "true"`,
  `LiteralFalse = "false"`, `BlockAll = "all"`,
  `BlockUnless = "unless"` — are consistent with the token
  spellings RDR 0002's normalizer and RDR 0003's grammar
  already use, so fixing them here pins a shared spelling
  rather than forking one.**
  - **Status**: Verified (kernel half) — the five spellings
    are confirmed consistent with every authored surface that
    exists today. The PRODUCER half rides on the same unlanded
    RDR 0002 duty as A16, A22 and A24, and fails CLOSED on
    mismatch, which is why it does not block. Named by the
    repeatability iteration-2 pass (G-2/G-3); the values were
    previously unstated, which is why three independent
    reconstructions each guessed them.
  - **Method**: Source Search (RDR 0003 operator vocabulary;
    RDR 0002 normalization; any authored fixture table)
  - **Evidence**: All five values check out against the only
    surfaces that spell them. Block names: RDR 0003 uses
    `all` / `unless` throughout its normative text — "Positive
    `all` predicates are conjunctive requirements. Negative
    `unless` predicates are also conjunctive" — and RDR 0002's
    `all`/`unless` split is the authored shape both cite.
    Operator token: RDR 0003's operator/kind matrix names
    `exists` ("Tests presence or absence, not value
    equality"). Boolean literals: the only authored guard
    table in the repo,
    `docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`,
    spells `[rule.guard.all.cluster_eligible]` / `exists = true`
    and a declared `domain = [true, false]` — lower-case,
    matching `LiteralTrue`/`LiteralFalse` exactly, and using
    the `all` block name in the same path. RDR 0002 contains
    no competing spelling (it has no tag-name grammar at all —
    see A22). So this RDR pins a shared spelling rather than
    forking one. The contract is byte equality
    performed by the kernel with no case-folding, so a
    spelling mismatch against what 0002 emits fails closed
    (a foreign token becomes a value atom; a foreign literal
    becomes `uncomparable`) rather than producing a plan —
    the direction of failure is safe either way, which is why
    the producer half is carried as a downgrade rather than a
    blocker.
  - **If wrong**: 0002/0003 already spell one of these
    differently and the kernel constants adopt that spelling
    instead; no rule changes, only the five byte values.

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
existence, and combination; `Refusal.Guard` is replaced by
the per-row, per-atom payload. `assemble`, `gate`,
`missingOwned`, `escapeOrRefuse`, and `Resolve` are
untouched in control flow (the payload is a data change
`gate` populates, not a reordering). Test fixtures
migrate from text-keyed verdicts to atoms.

#### Normative Contracts

**C1**
```normative
SEAM. A candidate row carries its guard as a slice of parsed
atoms — key, operator token, literal, block ∈ {all, unless}
— the shape JDR 0001 §D1 fixes; this RDR cites it and does
not restate the grammar. An empty slice is an unguarded row.
There is no opaque guard string and no reconstruction step,
so no mapping failure exists and no panic or error channel is
needed for one.

Two spellings of that shape are fixed HERE because this RDR's
payload depends on them, not because it is re-deciding §D1's
grammar. (1) The atom's four fields are named `Key`,
`Operator`, `Literal`, `Block` — the payload's `UndecidedAtom`
mirrors them by name, so divergent spellings would make the
mirror a mapping. (2) `Block` is an exported named STRING type
carrying exactly two constants, `BlockAll = "all"` and
`BlockUnless = "unless"`. The payload clause below requires the
payload's `Block` to be "the same exported constant type the
atom carries"; that requirement is unstatable while the type
and its members are unnamed. Nothing else about the atom —
operator vocabulary beyond `OpExists`, literal typing, the
parse — is this RDR's, and §D1 continues to own it.

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
  `Evaluate(atom, value) GuardResult` — a single-method
  INTERFACE named `GuardEvaluator`, not a func type, so the
  nil-seam rule above is a nil interface value and RDR 0003's
  evaluator satisfies it by declaring the method. The seam decides
  true or false under RDR 0003's typed operator semantics,
  and MAY answer unevaluable for a present value it cannot
  compare (A18). It never sees the view.

The kernel therefore enforces the domain rule by structure:
an evaluator cannot fold absence into false, because it is
never asked about an absent key.

A NIL seam does not make a row unevaluable by itself. The
kernel decides existence atoms and absent-key atoms without
consulting the seam, so a row whose atoms are all
kernel-decidable resolves under a nil seam exactly as it
would under a present one. A nil seam yields unevaluable
only for an atom that WOULD have been handed to it — a
value-comparing atom over a present key. Such an atom is
reported in the payload like any other unevaluable atom,
with reason `uncomparable` (its key is present, so `absent`
would be a lie); the payload never omits it and never names
the nil seam as a reason. This narrows the
shipped whole-guard behavior (`evaluateGuard` returns
GuardUnevaluable for any non-empty guard when `seam == nil`)
to a per-atom rule.

Frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
keeps its verdict but changes its REASON, and Phase 1 must
re-read it rather than re-encode it: its guard
`iterations >= 3` references a key the fixture view does not
carry (`legalInput` supplies `status`, `reviews`,
`recognized`), so after the reshape the KERNEL decides that
atom unevaluable on absence and the nil seam is never
consulted. The refusal kind the test asserts still holds; the
proposition its name states — that an absent seam is what
refuses — no longer does on this fixture. To keep testing the
nil-seam rule the fixture needs a value atom over a PRESENT
key (e.g. `reviews >= 3`).
```

**C2**
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

**C3**
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
the normalizer MUST emit for existence atoms (A16). The
constants are `OpExists` (the operator token) and
`LiteralTrue` / `LiteralFalse` (the two boolean literal
forms); the kernel compares an atom's operator and literal
against these values verbatim, performing no parsing,
case-folding, or coercion of its own.

Because that comparison IS byte equality, the bytes are
normative and stated here: `OpExists = "exists"`,
`LiteralTrue = "true"`, `LiteralFalse = "false"`, all
lower-case. A16's cross-component MUST ("the normalizer MUST
emit exactly these values") has no content without them — a
producer cannot emit what the contract does not spell, and
`"exists"` versus `"Exists"` is the whole difference between a
decided atom and a fail-closed value atom.

Drift at this boundary fails CLOSED: an existence atom carrying
a foreign TOKEN is a value atom to the kernel (unevaluable on
absence, handed to the seam on presence). A foreign LITERAL on
an `OpExists` atom — any value that is neither `LiteralTrue` nor
`LiteralFalse`, the empty literal included — is UNEVALUABLE at
the kernel, reason `uncomparable`; the kernel MUST NOT decide
such an atom from presence, and MUST NOT treat a missing literal
as `LiteralFalse`. The normalizer additionally MUST reject a
foreign literal at load (RDR 0002's typed validation), so the
kernel rule is the fail-closed backstop for a producer that did
not, not a duplicate of it: load rejection is upstream and
`internal/resolve` has no load path of its own. Neither route
can yield a plan from absence.
```

**C4**
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

**C5**
```normative
A row with no atoms is decided TRUE without consulting the
view (shipped: `evaluateGuard`'s empty-guard branch). Empty
atom blocks take the conventional identities: an empty `all`
block is TRUE (empty conjunction). An empty or omitted
`unless` block is ABSENT: the `¬(unless_conj)` TERM DROPS OUT
of the row verdict, which reduces to `all_result`. It is NOT
a vacuously-true conjunction — `unless_conj = T` would give
`¬T = F` and disable every row carrying an empty block — and
it is NOT `unless_conj = F` either, which would be a claim
about atoms that do not exist. The block contributes no
operand, not a neutral one; the reduction is stated on the
verdict, not on `unless_conj`, because a conjunction over
zero atoms has no value to state here.
```

**C6**
```normative
Atom verdicts combine in the KERNEL under strong-Kleene
three-valued logic. Conjunction is `min` under the TRUTH
order `F < U < T`; negation is `¬T = F`, `¬F = T`, `¬U = U`:

| ∧ | T | F | U |
|---|---|---|---|
| **T** | T | F | U |
| **F** | F | F | F |
| **U** | U | F | U |

The TABLES are normative; `min` names the lattice, NOT an
implementation. Shipped `resolve.go::GuardResult` declares
`GuardFalse, GuardTrue, GuardUnevaluable` (`iota` 0, 1, 2) —
the order `F < T < U`, which is NOT the truth order. Integer
`min` over those constants reproduces 7 of the table's 9
cells and fails exactly the commutative pair `T ∧ U`, which it
computes as `T` where the table requires `U`. The kernel MUST implement
combination against the tables (or against an explicit
truth-order rank), and MUST NOT `min` the raw constant values;
the constant order is a shipped fact this RDR does not
renumber (renumbering would be a silent behavior change for
any existing comparison).

The row verdict is `all_result ∧ ¬(unless_conj)`, where
`unless_conj` is the conjunction of the `unless` block's atoms
(`unless` is block-level negation, NOT per-atom negation).
When the `unless` block is empty or omitted the `¬(unless_conj)`
term drops out and the verdict is `all_result` alone (empty-block
clause above). A guard is GuardTrue or GuardFalse only when
decided atoms alone determine it; any unresolved dependence on an
unevaluable atom yields GuardUnevaluable. A present-tag atom
decided FALSE still yields GuardFalse beside an unevaluable
atom (`F ∧ U = F`): the falsity is witnessed by present
state and holds under every resolution of the unevaluable
atom — this is what makes D8 sound. The tables are
normative; the evaluator holds no part of them.
```

**C7**
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
(D5) — "identically" by DELEGATION, not by a parallel
implementation: shipped `escapeOrRefuse` filters the escape
rows and then calls the same `gate` function
(`gate(escapes, in.Guards, view)`), which is why it needs no
change here and inherits the payload for free. An unevaluable
CANDIDATE is never masked by a decidable
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
Combined reporting — one refusal carrying BOTH payloads, for
which this RDR's per-row/per-atom surface would be the
vehicle — is REJECTED here: `Refusal.Kind` is RDR 0001's
single stable discriminator that RDR 0005 maps to one CLI
code, so a both-ways refusal would need either a new kind or
a kind whose meaning depends on which payload fields are
populated. Both reopen a frozen taxonomy (REQ-7 pins the kind
set at exactly five) for a diagnosis the second round already
delivers. The two-round loop is the accepted cost of keeping
one refusal, one kind.
Authoring docs MUST NOT repeat the superseded rationale.
```

**C8**
```normative
The kernel maps GuardUnevaluable to `guard_unevaluable`,
which RDR 0002 excludes from escape lists. Missing artifact
state MUST NOT be maskable behind an escapable refusal class.

The refusal MUST name what was missing, per row and per
atom: for every undecidable row, its `(RuleID,
SourceLocator)` and, for each of its unevaluable atoms, the
referenced key, the block, and a reason drawn from a closed
set — `absent` (the key was not in the view) or
`uncomparable` (the key was present and its value was not
compared to a verdict).

`uncomparable` is defined by the OUTCOME, not by which
component produced it, and covers all three ways a present
key fails to decide: the seam answered unevaluable (A18); the
atom carries a foreign or missing literal on `OpExists`; or
the seam is NIL, so no comparison could be attempted. The
closed set stays at two members and every unevaluable atom of
an undecidable row therefore carries a reason — the payload's
per-atom completeness obligation above admits no third state
and no omitted entry. A nil seam is not itself a reason (it
is a wiring fact, not a property of the input tuple), which is
why the present-key/nil-seam atom reports `uncomparable`
rather than a reason of its own.

The payload is named surface, not shape-by-description; every
assertion below and RDR 0005's renderer type against it.
`Refusal` carries the field `Undecided []UndecidedRow`, where
an `UndecidedRow` names its `RuleID`, `SourceLocator`, and
`Atoms []UndecidedAtom`, and an `UndecidedAtom` names its
`Key`, `Block`, `Operator`, `Literal`, and `Reason`. `Reason`
is a kernel-owned closed constant set — `ReasonAbsent`,
`ReasonUncomparable` — spelled and enumerated the way
`resolve.go::RefusalKind` / `RefusalKinds()` already spell the
refusal taxonomy: a named STRING type with exported constants
plus an exported enumerator `Reasons() []Reason` returning the
set in declaration order, mirroring `RefusalKinds()`. The
enumerator is required surface, not an implementation choice —
"a test can pin the set exactly" is unsatisfiable without it.
RDR 0005 maps it the same way. `Block` is the same exported
constant type the atom carries, not a separately-spelled
payload value.

The reason values are PRODUCED BY THE VERDICT PASS, not
re-derived: the per-atom evaluation that decides an atom
unevaluable emits that atom's payload entry at the same step,
and the kernel carries the entries forward with the row
verdict. A second walk over the atoms after the row is known
undecidable is FORBIDDEN — it would consult the seam twice for
every present-key atom, and the seam is not required to be
pure or cheap. This fixes the kernel's internal data flow only;
the exported surface is the payload above.
This is the kernel's first two-level payload; it is additive on
`Refusal` and reopens no closed taxonomy (Load-Bearing
Decisions: a field, not a sixth kind). The spike's
`UndecidedAtoms` is NOT this shape — it sat beside a retained
`Refusal.Guard` (A3).

This payload replaces the
single-valued `Refusal.Guard` text — and with it the shipped
selection of ONE representative row
(`gate`'s `slices.MinFunc` over the undecidable set, which
picked the lowest `(RuleID, SourceLocator)` to fill that
field). Every undecidable row appears in the payload, so
there is no representative to choose; the `MinFunc` call is
retired, not re-typed. `Refusal.Guard` has no referent
once guards are atoms; `Refusal.Rows` continues to carry
every undecidable row. The kernel MUST evaluate every atom
of every survivor — no short-circuit on a decided block —
and MUST sort the payload so it is a function of the input
tuple (RDR 0001 REQ-1), never of atom or row order.

The sort key MUST be TOTAL over payload entries. Row identity
then key is NOT total: this RDR's own sanctioned idioms put two
atoms on one key in one row — A5's conjoined value row
(`X exists = true` beside `X eq v`; Testing Strategy row 9) and
the `unless` idiom in Failure Modes (`legal_hold exists = true`
beside `legal_hold eq true`) — and one key may go unevaluable in
both blocks of the same row. Two entries would then tie on
`(row, key)` and differ only in block, operator, or reason,
leaving the order to an unstable tie-break — the atom-order
dependence this clause exists to forbid, one level below the
row order ADV-3 already freezes. The ordering is therefore the
tuple `(RuleID, SourceLocator, key, block, operator token,
literal)`, compared field by field in that order. Those are the
row's identity plus the atom's four §D1 fields, so the tuple is
total by construction: two entries equal on all six name the
same atom, which the kernel MUST NOT report twice. Any contract
or test distinguishing WHICH row went unevaluable MUST assert on
the row entries.

Two conditions the totality claim rests on, both stated rather
than assumed. (1) The literal must have a byte-comparable
spelling. For scalar literals it does; for the SET-valued
literals RDR 0003 gives `in` ("non-empty typed scalar set") and
`contains` ("non-empty typed element set") it does not yet —
`resolve.Tag.Value` is a bare `string` and no element encoding
is declared (Phase 3). Until 0003 declares it, two `in` atoms
on one key differing only by literal-set ORDER are not
distinguishable by this tuple, so the kernel MUST compare the
literal as the exact bytes the normalizer emitted and 0003's
declaration MUST make that spelling canonical. This is the same
0002/0003 canonicalization duty A22 and A24 name, applied to
literals rather than keys. (2) "MUST NOT report twice" is a
kernel obligation, not a property of the sort: the kernel emits
one entry per unevaluable atom, and an atom is identified by
the six-tuple, so duplicate suppression is by construction of
the emit loop — the sort orders entries, it does not dedupe
them.

What the payload does NOT promise: a row pruned as
GuardFalse under `F ∧ U = F` is not a survivor, so its
absent keys are not reported. The guarantee is "never a plan
from absence," not "every absence reported"; an absence is
named exactly when it blocked a verdict.
```

**C9**
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

**C10**
```normative
SURVIVOR MEMBERSHIP. A row whose guard is GuardTrue or
GuardUnevaluable is a SURVIVOR; only GuardFalse rows are
pruned. The owned-state scan runs over survivors, so an
unevaluable row's RequiresOwned keys DO raise
owned_state_unavailable, and the undecidable check runs over
the same set (shipped partition in `gate`).

The owned scan runs ONCE over the whole survivor set, and
`MissingOwned` is the deduplicated, sorted union of the
missing keys across survivors — shipped `missingOwned` already
takes `[]Row`, dedupes through a `seen` map, and
`slices.Sort`s before returning. Stated because REQ-1 purity is
claimed for the whole result and this RDR pins determinism
only for the `Undecided` payload: a per-row scan whose results
were concatenated would be row-order-dependent and
duplicate-bearing. This is shipped behavior restated, not a
change.

Aggregation is resolution-level: if any surviving candidate
row's guard is GuardUnevaluable, the resolution MUST refuse
`guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be
selected while an unevaluable candidate exists. This veto is
the cost the RDR accepts: one unreadable row refuses a table
whose other rows decide cleanly. Narrowing it would decide
that an undecided edge is a non-edge — absence-as-false at
resolution scope. Where RDR 0006's exhaustiveness proof and
this veto disagree, lint's promise narrows (JDR 0001 §JD-4).

The veto carries NO prior-art support, and the RDR does not
claim any: SQL:2003's strict routines and PostgreSQL's K3
tables (A4, research C3) govern the PER-ATOM rule only. Read
at resolution scope, SQL is in fact the counter-example — a
`WHERE` clause evaluating to UNKNOWN does not qualify the
row and the statement proceeds, which is the
absence-as-non-selection this clause rejects. The veto rests
on this RDR's own argument (an undecided edge is not a
non-edge) and on the closed-escape asymmetry RDR 0002 fixes,
not on precedent. Stated because the two citations sit one
clause away and could otherwise be read as covering both
levels.
```

**C11**
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
  harness is unenforced; held by the kernel it cannot drift,
  and the refusal can name the absent keys.
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
| Existence-operator token | RDR 0003 (Final) | Introduced | Kernel exports the constant; normalizer emits it (A16) |
| Value evaluator over present values | RDR 0003 implementation (future) | Deferred | Narrow seam; value-contract tests (Phase 3) |
| `RequiresOwned` producer | RDR 0002 (Final, §JD-3 closed) | Deferred | Escape rows carry none (A21) |
| Read completeness | RDR 0004 (§D3) | Decided | Cited; no obligation here |
| Escape-class closure | RDR 0002 | Available | `guard_unevaluable` non-escapable by construction |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Presence predicate | `internal/resolve/resolve.go::TagSet.Lookup` | None for this use (provenance-blind `ok`) | Reuse | The kernel's atom presence test |
| Undecidability verdict + refusal mapping | `resolve.go::GuardUnevaluable`, `resolve.go::gate` | `Refusal` is flat scalars plus `[]RowRef`/`[]string`; it carries guard text, not keys, and has never held a nested value | Reuse the verdict + kind; EXTEND `Refusal` with the kernel's first two-level payload | Per-row/per-atom array-of-objects — the reason the A19 envelope question exists |
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
| Drift enforcement (premortem P-10) | ✗ — a conformance harness nothing compels 0003 to run leaves the residual unmitigated | ✓ — no evaluator code path exists for absence; only value semantics can drift, and those are 0003's to test | ✓ |
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
granularity. Its cost over A is bounded kernel logic that A would
otherwise need anyway as a "vector-local, non-normative"
view-reading evaluator — B makes that code the normative one
instead of a stopgap. The
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
Ledger: `docs/rdr/0007-guard-predicate-totality/evidence/propose-premortem/iter-2/critic.md`.

Ground-sweep: clean (24 anchors) — 15 code anchors
(`resolve.go`, test files, `root.go`), 9 peer-document
passages.

Joint-check: fired → 0003, 0002, 0009 (home: JDR 0001 §D4 /
§JD-12), 0008, 0004 (home: JDR 0001 §D1). §D1 settles the
*transport* (atoms on `Row`, the seam change 0008/0004 cite
as shipped code); §D4 settles the *enforcement site*: the
value-only per-atom seam, the kernel-exported existence
constants RDR 0002's normalizer must emit (A16), key identity
(A22), and the per-atom payload replacing `Refusal.Guard`
(which Final RDR 0009 cites; stale, rides to its re-lock per
§JD-12). This RDR is §D4's landing document; 0003 and 0002
cite. Absence arm: n/a (no refusal is converted to an
acceptance). Bridge sub-check: n/a — no sibling plan retires
a surface this plan introduces.

## Alternatives Considered

### Alternative 1: Evaluator-enforced domain rule over parsed atoms

**Description**: Take §D1's atom slice but keep a row-level
seam — `Evaluate(atoms, view) GuardResult` — so RDR 0003's
evaluator applies the domain rule and the strong-Kleene
tables itself; the kernel stays fully grammar-agnostic and
only maps the verdict.

**Pros**:

- Kernel learns no operator token; the cleanest reading of
  "cardinality, not meaning."
- Smallest kernel diff: `Row` and the seam signature only.

**Cons**:

- The domain rule is held by discipline: a conformance
  harness that nothing compels RDR 0003's build to run, so
  the residual is unmitigated.
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
- **The domain rule is scoped to guards, so the match
  pattern still folds absence into non-match.**
  `resolve.go::TagSet.matches` returns false on `!ok`, so a
  row whose MATCH pattern references an absent key drops out
  of candidacy and yields the escapable `no_match` — the
  Background probe's shape, reachable with no guard involved.
  This is deliberate (the match pattern is closed-world by
  design, RDR 0001) and NOT closed by this RDR: predicate
  PLACEMENT decides whether an absent key refuses
  non-escapably or escapes. The only control is authoring
  guidance — a predicate whose absence must refuse belongs in
  a guard atom, never in the match pattern — which Phase 4
  hands to RDR 0003. Named here so the "masking path closed"
  claim is read at its true scope: closed for guard atoms,
  unchanged for match tags.
- **A pruned row reports nothing.** Under `F ∧ U = F` a row
  with a decided-FALSE atom beside an unevaluable one is
  pruned, so its absent key appears in no payload (the
  payload clause's stated non-promise). The author asking
  "why didn't my row fire?" gets no surface answering it.
  Accepted: reporting absences on pruned rows would report
  every absence in the table on every resolution. The
  mitigation is authoring-side (Phase 4) — the risk is that
  it pushes authors toward sentinel-stamping to make the
  kernel talk, which is the named anti-pattern.
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
- [ ] **JDR 0001 §D4 is reopened and §JD-8 amended to grant
      ONE `omitempty` structured field on `CLIError` (A19).**
      §D4 (Closed) routed this payload through §JD-8's
      `Detail` as flattened text; the author's round
      2026-08-21 chose a structured field instead, so this is
      a REVERSAL of that ruling and must be taken back to the
      JDR, not folded in here. This RDR MUST NOT state the
      field itself. **Does not block lock**: the direction is
      reversible either way (an `omitempty` field is
      append-only), and A19 carries a `Detail`-flattening
      fallback if the reopening is declined. Blocks RDR
      0005's implementation of the JSON rendering only.
- [ ] Two duties are added to RDR 0002's Refinement Context
      Direction list, which currently carries neither —
      it names only the `RequiresOwned` producer (JD-3) and
      the fixture rename (JD-10):
      (a) the normalizer emits the kernel's exported
      existence token and its two boolean literal forms
      (A16); (b) canonicalization of authored tag-key
      spellings is bound as a normative clause (A22; JDR 0001
      §D4 landing note). Until (a) lands, A16's producer half
      is acknowledged only by a joint-check line; until (b)
      lands, A22 stays an inference and the §D4 obligation
      dangles. A22's fallback is to narrow it to kernel-side
      only and re-file canonicalization as an obligation on
      0002.
- [ ] **RDR 0009 re-locks against the loss of
      `Refusal.Guard`.** 0009 is `Final` and names
      `Refusal.Kind` / `Refusal.Guard` as the asserted
      properties carrying its non-vacuity evidence; this RDR
      deletes that field, so 0009's evidence goes stale on
      implementation. Re-run its conformance evidence against
      the atom-shaped refusal. **Also sequence 0009's kernel
      entry precondition against Phase 1's fixture migration
      (A25):** both rewrite `fixtures_test.go::escapeRow` and
      the four frozen escape tests routed through it, and
      whichever lands second re-runs a suite the first already
      re-encoded.
- [ ] **RDR 0003 states the empty/omitted `unless` identity
      (A9).** 0003 mentions the case zero times; read
      literally, its subtractive lint algebra over an empty
      block subtracts the whole scoped product and disables
      the row, which contradicts this RDR's term-dropping
      verdict rule. Until 0003 states it, lint and the kernel
      may select different row sets — JDR 0001 P5's promise
      failing in the direction that makes lint worthless.
- [x] A3, A16–A19, A21–A27 dispositioned (Stage 6 reconcile,
      2026-08-21). A16, A17, A18, A21 verified 2026-08-21.
      **A3, A23, A24, A26 VERIFIED** by one re-spike — the
      reshape with `Refusal.Guard` actually removed
      (`evidence/spikes/a3-reshape-respike.md`; 154/154 frozen,
      180/180 with probes, `go vet` clean). It supersedes the
      first spike, which proved a superset shape. **A27
      VERIFIED** (kernel half) against RDR 0003's vocabulary
      and the authored `guard-fixture.toml`. **A19, A22, A25
      and A12 DOWNGRADED** — each blocked on a peer document
      or a cluster decision, not on evidence this RDR can
      produce; each carries a named plan and a survivable,
      fail-closed fallback recorded at the assumption. No
      assumption remains Pending and none is MVV-critical.
- [ ] This RDR Final **before** RDR 0003's implementation
      begins: 0003's evaluator is built against this seam.

### Minimum Viable Validation

Against the real kernel with a real atom — no stub — the
masking probe inverts: one candidate row whose guard is a
value atom over an absent key, whose `RequiresOwned` is
SATISFIED, and a modeled `no_match` escape row yield
`Refusal.Kind == guard_unevaluable`, a nil `Plan`, the absent
key in the payload, and the row in `Refusal.Rows`; the same
table with the key present and the value decided FALSE
prunes and escapes (D8 preserved).

The satisfied `RequiresOwned` is load-bearing, not
incidental: an unevaluable row is a SURVIVOR, so the owned
sweep runs over it and reports first (GATE, THEN COUNT;
Testing Strategy row 13). A row carrying BOTH an absent
owned key and an unevaluable atom yields
`owned_state_unavailable`, not `guard_unevaluable` — that
pairing is scenario 13's subject, and putting it in the
masking probe would test the owned sweep while claiming to
test the domain rule. Third scenario:
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
guard in the reference fixtures as safe-or-migration. Retire
`gate`'s `slices.MinFunc` representative-row selection with
the field it fed. Rewrite `gate`'s doc comment, which still
states D8's superseded rationale verbatim ("absent owned
state … is frequently the reason the seam could not decide
the predicate") — the clause above forbids repeating it, and
the kernel comment is the authoring doc closest to the code.

### Phase 3: Value-seam contract tests

Export a small contract-test function for the value seam
(present value × literal × operator; unparseable value →
unevaluable, never false) that RDR 0003's implement stage
instantiates against its evaluator. It is
`resolve.TestGuardEvaluatorContract(t *testing.T, seam
GuardEvaluator)` — named here because it is a CROSS-RDR
surface 0003 must call by name, and a caller cannot
instantiate a function whose name and signature the
contract leaves open.

**Blocked on one RDR 0003 declaration.** `resolve.Tag.Value`
is a bare `string`, so a set-valued tag reaches the seam as
one opaque string with no stated element encoding. The
`contains` leg of this contract test therefore cannot be
written from any current document, and neither can Testing
Strategy scenario 8's present-key half. This RDR
deliberately does not invent the encoding — element
representation is part of RDR 0003's typed operator
semantics. What this RDR fixes is unchanged and needs no
encoding: an ABSENT set-valued tag is unevaluable, decided
by the kernel on presence alone. Phase 4 carries the
declaration request.

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
sentinel-stamping as the anti-pattern; and the REQUEST that
0003 declare the element encoding of a set-valued tag value,
without which Phase 3's `contains` contract test cannot be
written.

## Validation

### Testing Strategy

Every scenario is a kernel test in `internal/resolve`
(package `resolve_test`), driven THROUGH atoms and a value
stub — never by injecting row verdicts, which would leave
the kernel combinator off the tested path (premortem P-20).
The A3 RE-SPIKE (`evidence/spikes/a3-reshape-respike.md`)
established the baseline this matrix extends, running the
reshape as specified — `Refusal.Guard` DELETED, the seam
view-free, existence literals explicit. The frozen 154-test
suite passes under it (154/154; 180/180 with the probes), so
every row below is ADDED surface and the frozen suite is
expected to pass re-encoded EXCEPT Fixup-1d.
`fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
is RE-DECIDED, not re-encoded, and the re-spike showed the
reason runs deeper than the deleted field: that fixture's
guard names an ABSENT key, so the kernel decides it on
absence (`Reason:absent`, captured verbatim) and the nil seam
is never consulted — the test would go green while testing
nothing about the nil seam. Phase 1 moves the atom to the
PRESENT key `reviews` (expect `uncomparable`), per the SEAM
clause, and records the disposition. Phase 1's exit condition
is that count with the field actually removed — not the 154
the superseded superset spike
(`evidence/spikes/a3-reshape-probe.md`) returned.
One coverage gap the re-spike PROVED rather than inferred:
every frozen guard is a single atom, so the 154 cannot
specify the combinator — a FALSE-vs-UNEVALUABLE dominance
mutant survives all of them and is killed only by a
multi-atom conjunction (row 21).

**The coverage floor this matrix exists to hold.** The spike
mutation-tested the combinator: swapping strong-Kleene
FALSE-dominance for UNEVALUABLE-dominance in `evaluateGuard`
**survives all 154 frozen tests**, because every frozen
guard is a single atom and a one-atom conjunction cannot
observe combination order. Scenarios 3–5 below are what kill
that mutant. Any future edit to the combination tables MUST
be re-mutation-tested against them.

| # | Scenario | Backing assumption | Code path / evidence |
| --- | --- | --- | --- |
| 1 | Value atom over an absent key + modeled `no_match` escape ⇒ `guard_unevaluable`, nil `Plan`, key in payload, row in `Refusal.Rows` — the masking probe inverted | A2, A6a | `evaluateAtom` presence step; MVV scenario 1 |
| 2 | Same table, key PRESENT and value decided FALSE ⇒ row prunes and the escape plans (D8 preserved) | A2, A20 | `gate` prune branch, unchanged partition (spike: byte-identical) |
| 3 | `F ∧ U = F` — a present-tag atom decided FALSE beside an unevaluable atom still prunes | A2, A4 | strong-Kleene table; **mutation-killing** (spike probe case 1) |
| 4 | `T ∧ U = U` — an unevaluable atom beside all-TRUE siblings ⇒ `guard_unevaluable` naming exactly that atom | A2, A4 | **mutation-killing** (spike probe case 2) |
| 5 | Full strong-Kleene matrix across `all` and `unless`, incl. `¬U = U` on a block-level negation | A2, A4 | normative tables; `unless` is block-level, not per-atom |
| 6 | unevaluable-blocks-true-sibling: row A unevaluable beside row B decided TRUE ⇒ refusal naming A, no plan (the aggregation veto) | A2 | `gate` undecidable branch; MVV scenario 3 |
| 7 | `exists` total across BOTH polarity channels: `exists = true` over a present key and `exists = false` over an absent key each decide, never refuse | A14, A16 | kernel-decided, seam never consulted; MVV scenario 4 |
| 8 | `contains` over an ABSENT set-valued tag ⇒ unevaluable, NOT the empty set | A1, A17 | the clause A1 pins against RDR 0003's silence |
| 9 | Two-row absence pattern (`X exists = false` row ‖ `X exists = true` + `X eq v` row), plus a bare-value-row negative control that must refuse | A5, A14 | disjointness; authoring grammar handed to 0003 (Phase 4) |
| 10 | Unguarded row (empty atom slice) ⇒ TRUE without consulting the view; empty `all` ⇒ TRUE; omitted `unless` ⇒ absent, contributes nothing | A9 | shipped empty-guard branch in `evaluateGuard`; 9(b) is now a kernel test |
| 11 | Foreign existence token fails CLOSED: an atom carrying a non-exported token is a VALUE atom (unevaluable on absence, seam on presence). Foreign/missing LITERAL on an `OpExists` atom ⇒ kernel-unevaluable, reason `uncomparable` — asserted in `internal/resolve`; the load-rejection leg is RDR 0002's and is NOT a kernel test (no load path exists under `internal/resolve`) | A16 | the drift boundary §D4 prices |
| 12 | Seam answers unevaluable on a PRESENT value it cannot compare ⇒ payload reason `uncomparable`, not `absent` | A18 | the reason field is what keeps skew from reading as absence |
| 13 | Owned-before-unevaluable precedence within ONE survivor set (a row failing both ways surfaces the write-dependency first) | A20, A21 | `gate` order pinned to shipped behavior, not to D8's superseded rationale |
| 14 | Escape-set scoping, both legs: an unevaluable CANDIDATE is never masked by a decidable escape row; an unevaluable ESCAPE row yields `guard_unevaluable` | A20 | `escapeOrRefuse` (spike: ZERO diff); frozen ADV-2 + Fixup-1d |
| 15 | Escape row raises no `owned_state_unavailable` of its own (empty `RequiresOwned` by composition). The test MUST build the escape row explicitly with empty `Writes` and empty `RequiresOwned` and assert the kernel consequence; the shipped `fixtures_test.go::escapeRow` is non-conforming on BOTH fields and must not be reused as-is. Note the `Writes` half is not merely unconforming to A21's composition — it BREACHES RDR 0009's kernel entry precondition, so once 0009 lands the fixture yields a non-nil `Resolve` error (A25) | A21, A25 | `missingOwned` over an empty slice never enters the loop; who POPULATES `RequiresOwned` is RDR 0002's producer duty (§JD-3), untested here, but the `Writes`-empty half is kernel-enforced by 0009 |
| 16 | Guard-path observed-tag substitution: a caller-supplied observed tag turns `guard_unevaluable` into a SPECIFIC verdict — assert the exact disposition (`--tag X=v` against `X eq v` ⇒ row selected and a plan; against `X eq w` ⇒ row pruned, no plan), never merely "not `guard_unevaluable`", which the masking path would also satisfy | A11, A13 | no existing test contends the GUARD path (ADV-4/ADV-5 defend only the owned path) |
| 17 | Payload determinism: sorted on the TOTAL tuple `(RuleID, SourceLocator, key, block, operator, literal)`, so it is a function of the input tuple and never of atom or row order (RDR 0001 REQ-1). Includes the tie case the total order exists for: one row carrying two atoms over ONE key (A5's conjoined row; the `unless` idiom), permuted in the slice — both payloads MUST compare equal | A3 | `undecidedAtoms`/`compareAtoms` (both PROPOSED — modeled on shipped `resolve.go::compareRefs`/`copyTags`); frozen ADV-3 is the order-independence precedent |
| 18 | Provenance-blindness: an atom decides identically whether its key arrived owned, observed, or recognized | A11 | `TagSet.Lookup` `ok`, not `TagSet.has` |
| 19 | Nil seam is per-atom, not whole-guard. The DISCRIMINATING leg is the `exists` atom over a PRESENT key: shipped `evaluateGuard` returns GuardUnevaluable for any non-empty guard when `seam == nil`, whereas this rule decides it TRUE/FALSE — a nil-seam row that now yields a PLAN. This leg is REQUIRED, not optional: it is a widening of shipped behavior (previously any non-empty guard refused without a seam; now a pure-`exists` row plans and its `Writes` land), and no frozen test asserts it, so without this row a wiring bug that leaves `Input.Guards` unset ships as silent plan-production. Assert both that the plan is produced AND that a value atom over a present key still refuses under the same nil seam. (The absent-key leg is non-discriminating: unevaluable either way, though the REASON becomes `absent` rather than want-of-a-seam, and a nil seam is never itself a payload reason.) | A3, A17 | narrows shipped `evaluateGuard`'s `seam == nil` branch; Fixup-1d must be re-read (its `iterations` key is absent from the fixture view, so the kernel now decides it) |
| 20 | Every undecidable row appears in the payload — no representative row is chosen (shipped `gate` picked one via `slices.MinFunc` to fill `Refusal.Guard`; that selection retires with the field) | A3 | `gate`'s `MinFunc` call; ADV-3 order-independence still applies to the sorted payload |
| 21 | Combination is table-driven, NOT integer `min` over the shipped constants: assert `T ∧ U = U` AND, on a MULTI-ATOM conjunction, that `F` dominates `U` (`F ∧ U = F` ⇒ the row is PRUNED, not undecidable). Shipped `GuardResult` declares `F, T, U` as `iota` 0,1,2. The re-spike mutation-tested both legs: a raw `min` is caught by the frozen suite already (`T <= U` folds even a single-atom `U` to `T`), but swapping FALSE-dominance for UNEVALUABLE-dominance **survives all 154** and is observable only across two atoms — that leg is this row's real work | A4 | the K3 clause's explicit non-`min` statement; **mutation-killing (verified: `evidence/spikes/a3-reshape-respike.md` §Q-A3)** |
| 22 | A value atom over a PRESENT key under a NIL seam appears in the payload with reason `uncomparable` — not omitted, and not `absent`. Pairs with row 19's nil-seam legs: assert the payload entry, not just the refusal kind | A18 | the widened `uncomparable` definition (outcome, not producer); the per-atom completeness obligation |
| 23 | `MissingOwned` is the deduplicated, sorted union across survivors: two survivor rows requiring the same absent owned key yield ONE entry, and permuting the rows yields an equal slice | A20 | shipped `missingOwned` `seen`-map + `slices.Sort`; REQ-1 purity over the whole result |

**Not covered here, by ownership.** Value-operator semantics
over a present value are RDR 0003's; Phase 3 exports a
contract-test function (present value × literal × operator;
unparseable value ⇒ unevaluable, never false) that 0003's
implement stage instantiates against its evaluator. The CLI
rendering of the payload is RDR 0005's and is blocked on the
A19 envelope-field grant — or, if that reopening is
declined, on A19's `Detail`-flattening fallback. The kernel
side is unblocked either way: the payload exists on
`Refusal` regardless of how the CLI renders it.

#### Mini-checks

Four structural cues fire on this draft (source-authority,
test-discriminability, disposition, desk trace); round-trip
/ fidelity does not — the draft's "migrate" is fixture
migration, and no inverse invariant is claimed.

**`authority`** — who decides what, and which decision is canonical.

| Input / decision | Writer | Readers | Sibling arms | Canonical |
| --- | --- | --- | --- | --- |
| Atom presence in the view | kernel (`TagSet.Lookup` `ok`) | `evaluateAtom` | none — the seam never sees the view | **kernel**, provenance-blind |
| Existence-atom verdict | kernel, from presence alone | combination step | seam is NOT consulted | **kernel** |
| Value-atom verdict over a PRESENT key | evaluator seam `Evaluate(atom, value)` | combination step | kernel decides nothing here | **seam** (RDR 0003 semantics) |
| Value-atom verdict over an ABSENT key | kernel (unevaluable) | combination step | seam is NOT consulted | **kernel** — the domain rule |
| Atom combination (K3) | kernel | `gate` | evaluator holds no part of the tables | **kernel** |
| Tag-key canonicalization | RDR 0002 normalizer, pre-`Row` | kernel compares exact strings | kernel performs none (A22) | **normalizer** (duty not yet a 0002 clause) |
| Existence token + literal forms | kernel exports constants | 0002 normalizer emits them | foreign token ⇒ value atom (fails closed) | **kernel constants** (A16) |
| Guard-input coverage | the domain rule | — | `RequiresOwned` is post-guard writes only | **domain rule**, NOT `RequiresOwned` |
| Refusal payload CLI transport | JDR 0001 §D4 (`Detail`) vs A19 (structured field) | RDR 0005 renderer | **unresolved — see A19** | §D4 stands until reopened |

**`oracle`** — each MVV row's discriminating power and its failing control.

| MVV row | Fails if X is wrong because Y | Negative / failing control |
| --- | --- | --- |
| 1 — value atom over absent key ⇒ `guard_unevaluable`, key in payload, row in `Rows` | asserts on refusal KIND + payload CONTENT, so absence-folded-to-false surfaces as a plan instead of a refusal | same table, key present + value TRUE ⇒ plan; the pre-RDR kernel returns a plan here |
| 2 — key present, value FALSE ⇒ prunes and escape plans (D8) | asserts a PLAN is produced, so a regression that refuses on decided falsity fails loudly | row 1 is the control: same shape, absent key, must NOT plan |
| 3 — unevaluable A beside decided-TRUE B ⇒ refusal naming A | asserts B is NOT selected; the aggregation veto is the only thing stopping it | drop the veto ⇒ B is selected and the test fails |
| 4 — `exists = false` over absent key ⇒ row SELECTED | asserts the one sanctioned absence⇒verdict route still works; guards against over-refusing | row 1 is the control: a VALUE atom over the same absent key must refuse |

Not absence-of-error oracles: every row asserts a specific
kind or a specific selection, and each has a sibling that
must come out the other way. The spike's mutation result is
the standing evidence that verdict-only assertions are
insufficient — scenarios 3–5 exist to kill that mutant.

**`disposition`** — input class ⇒ outcome.

| Input class | Refusal kind / outcome | Exit group | Payload minted | Silent or loud |
| --- | --- | --- | --- | --- |
| All atoms decided TRUE | row selected → plan | 0 | none | loud (plan) |
| Any atom decided FALSE from present state | row PRUNED (D8) | n/a | none — pruned rows report no absences | silent by design |
| Survivor missing an owned `RequiresOwned` key | `owned_state_unavailable` | **no code row yet** (JDR §JD-8) | `MissingOwned` + `Rows` | loud |
| Survivor with an unevaluable atom (absent key) | `guard_unevaluable` | `GroupUserEnv` | per-row/per-atom, reason `absent` | loud |
| Survivor with an unevaluable atom (present, uncomparable) | `guard_unevaluable` | `GroupUserEnv` | per-row/per-atom, reason `uncomparable` | loud |
| Row failing BOTH ways in one survivor set | `owned_state_unavailable` first | `GroupUserEnv` | owned payload only | loud (the accepted cost) |
| Unevaluable ESCAPE row | `guard_unevaluable` replaces the candidate refusal | `GroupUserEnv` | escape set's payload | loud (ADV-2 / Fixup-1d) |
| Foreign existence token | treated as VALUE atom ⇒ unevaluable on absence | `GroupUserEnv` | reason `absent` | loud (fails closed) |
| Foreign boolean literal on an `OpExists` atom | kernel: unevaluable, reason `uncomparable`; normalizer: rejected at load | `GroupUserEnv` (kernel leg) / producer defect (load leg) | per-row/per-atom | loud, both legs |

Two exit-group notes, both upstream of this RDR and neither
fixable here. (1) `owned_state_unavailable` has NO row in RDR
0005's stable-code table — the gap JDR 0001 §JD-8 names
explicitly ("`owned_state_unavailable` … need `Code` values").
The cell above says so rather than asserting a mapping that
does not exist. (2) `flow-guard-unevaluable` DOES have a row,
but scoped there to "guard cannot be evaluated from supplied
facts" — while this RDR's flagship case is missing *artifact*
state, which `GroupUserEnv` (exit 2, "bad input") mis-signals
as a caller-input problem. Both are RDR 0005's to fix at its
re-lock; recorded here so the payload's remedy story is not
read as settled.

**`trace`** — the MVV walked stepwise against every assertion in force.

| Step | Assertions in force | Witness |
| --- | --- | --- |
| 1. Assemble view | provenance precedence owned > observed > recognized; presence is provenance-blind | `legalInput`: `status`(owned), `reviews`(observed), `recognized` |
| 2. Evaluate each atom | domain rule; existence totality; seam only on present keys; nil seam is per-atom (row 19) | absent-key value atom ⇒ U, seam never called |
| 3. Combine per row | K3 tables (`min` names the truth order, not an implementation over the shipped `F,T,U` constants); `all ∧ ¬(unless_conj)`; empty `unless` contributes nothing | one U atom, no F ⇒ row U |
| 4. Prune | only GuardFalse prunes; F∧U=F is witnessed falsity | MVV row 2 prunes; MVV row 1 does not |
| 5. Owned sweep over SURVIVORS | unevaluable rows ARE survivors, so their `RequiresOwned` counts; owned-before-unevaluable | MVV row 1's `RequiresOwned` is satisfied ⇒ passes through to step 6. That witness is non-discriminating by construction (it passes under either ordering); the ordering's witness is Testing Strategy row 13, where one row fails BOTH ways and `owned_state_unavailable` MUST win |
| 6. Undecidable veto | any surviving U ⇒ refuse; no decided-TRUE sibling escapes | MVV row 3: B not selected |
| 7. Mint payload | every undecidable row, sorted on the total tuple `(RuleID, SourceLocator, key, block, operator, literal)`; no representative row (row 20) | `MinFunc` retired |
| 8. Escape gating | escape set gated identically; unevaluable escape replaces the candidate refusal | Fixup-1d |

**One CONTRADICTION found and fixed in this pass.** Step 5
caught it: the MVV previously gave scenario 1 an absent
`RequiresOwned` key AND an unevaluable atom while asserting
`guard_unevaluable`. Under this RDR's own SURVIVOR
MEMBERSHIP and GATE-THEN-COUNT clauses — and under shipped
`resolve.go::gate`, which returns `KindOwnedStateUnavailable`
before reaching the undecidable loop — that input yields
`owned_state_unavailable`. The primary validation scenario
contradicted three clauses of its own draft and would have
failed on first implementation. The MVV now requires the
row's `RequiresOwned` to be satisfied; the both-ways pairing
stays where it belongs, in Testing Strategy row 13.

### Performance Expectations

No performance dimension. Evaluation stays linear in the
atom count per candidate row — the kernel replaces one seam
call per row with one presence lookup per atom plus a seam
call per present-key value atom — and the payload sort is
over undecidable rows only, on the refusal path. `Resolve`
is pure and in-process; no new I/O.

## Finalization Gate

Responses: 0007-guard-predicate-totality/artifacts/gate.md (Gate PASS 2026-08-21)

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
