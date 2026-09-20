# Recommendation 0012: Declared-kind carrier at the guard value seam

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
  declared-kind carrier) whose every candidate shape modifies a
  seam pinned by RDR 0007/RDR 0003, and every `GuardEvaluator`
  implementer inherits the answer.
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
- **Related Issues**: kata `intrastate#cq5p` (1528)
- **Predecessors**: 0007-guard-predicate-totality,
  0003-guard-predicate-exhaustiveness
- **Seam Lineage**: no prior accretion

## Problem Statement

A caller supplying runtime tag values (`flow resolve --tag iter=many`)
onto an `int` or `bool` dimension expects a value the declared kind
cannot parse to refuse loudly — `flow-guard-unevaluable` — the same
three-valued honesty the seam guarantees for absent facts. Today `eq`
and `in` compare the held value as a raw string
(`internal/guard/grammar.go::Evaluator`), so a present-but-malformed
value decides `GuardFalse`: the runtime silently prunes a candidate row
or selects a sibling/`unless` row it never actually compared. For a
consumer like the rdr navigator that is a silent misroute — worse than
the loud refusal the seam was designed to prefer.

The semantics are already fixed — `0007:C1` and the shared conformance
suite (`internal/resolve/guardcontract.go`) state that a present value a
typed operator cannot parse is unevaluable, never false — but no carrier
delivers the tag's *declared kind* to the comparison: `guard.Evaluator`
is `struct{}`, deliberately declaration-free per Verified assumption
`0007:A17`, and only the int-arity operators (`lt/lte/gt/gte`) can infer
a kind from the operator alone. The fork `0007:A18` deliberately
deferred must now be decided: which carrier makes the fixed semantics
computable — an evaluator constructed over declarations, a widened seam
signature, the kind stamped on the atom, or an upstream producer
obligation at assemble/merge — given that every candidate collides with
a pinned shape (REQ-10's frozen `(atom, value)` signature, the
stateless-evaluator reflection test, JDR 0001 D1's atom shape, or
Resolve's refusal taxonomy).

## Critical Assumptions

- **A1 Every guard-referenced key is declared in the loaded model, so a
  kind lookup at the seam always answers when the evaluator is
  constructed over the same model the evaluated table came from.**
  - **Status**: Pending
  - **Method**: Peer RDR + Source Search
  - **Evidence needed**: `0007:A7` ("Every guard-referenced key is
    declared … a typo'd key fails table load") plus the load-time check
    that enforces it in `internal/table` — confirm the refusal fires for
    a guard atom naming an undeclared key, not only for match keys.
  - **If wrong**: `eq`/`in` over an undeclared key answers
    `GuardUnevaluable` under the defensive no-declaration arm, and a
    guard that decided before this change now refuses
    `flow-guard-unevaluable` — loud, but a behavior change A1 did not
    predict.
- **A2 The zero-field reflection assertion in
  `internal/guard/guard_evaluator_0003_test.go::TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView`
  is a structural PROXY for REQ-34's normative text ("the evaluator MUST
  NOT read the tag view"), not itself a fenced 0003 contract — so
  narrowing it to "no view-typed or runtime-valued state" preserves
  every normative obligation while admitting declaration state.**
  - **Status**: Pending
  - **Method**: Peer RDR + Source Search
  - **Evidence needed**: `0003:C9`'s fenced text and the REQ-34 quote in
    `0003`'s req-list — confirm neither says "stateless" or "zero
    fields"; the field check appears only in the test body.
  - **If wrong**: statelessness is itself fenced and the choice must
    move to the atom-stamp carrier (Alternative 2), reopening JDR 0001
    §D1's atom shape instead.
- **A3 Every non-test construction site of the evaluator has the loaded
  model (or its declarations) in scope:
  `internal/cli/flow_resolve.go::guardSeam`,
  `internal/guard/product.go::valueSatisfies`,
  `internal/graphlint/reach.go::atomAdmitsValue`.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence needed**: the call chains into all three sites, showing a
    `*table.Model` (or equivalent declaration carrier) reachable without
    signature surgery beyond the site itself.
  - **If wrong**: a consumer without a model needs a declaration-free
    degraded evaluator, re-introducing the raw-string arms the RDR
    exists to remove — the constructor design must change (e.g. thread
    the model through the caller).
- **A4 Typed comparison for `int` (parse both sides, compare parsed
  values) changes no lint verdict: lint's product enumerates held values
  rendered from the declared domain in canonical form, so a
  non-canonical spelling (`"07"`) never appears on the lint side, and
  the only runtime verdicts that flip are malformed-value
  `GuardFalse → GuardUnevaluable` (the fix) and non-canonical-numeral
  `GuardFalse → GuardTrue` — the latter reachable from either side of
  the comparison, since C2 parses held value AND literal (a
  caller-supplied `--tag n=07`, or an authored `n eq 07` in a guard).**
  - **Status**: Pending
  - **Method**: Source Search + MVV Test
  - **Evidence needed**: `internal/guard/lint.go`/`product.go` render
    path showing held values come from declared domains; an MVV case
    exercising `--tag n=07` against `n eq 7`; and a before/after lint
    diff over the checked-in model corpus (`models/`) proving no
    committed model's lint verdict flips — authored-literal flips
    (`n eq 07` in a guard) count, not just held-value flips.
  - **If wrong**: lint and runtime disagree on a canonical-form edge and
    the fix must narrow to validate-then-byte-compare (parse to prove
    comparability, compare raw bytes), which A4's verification decides.
- **A5 The conformance suite can publish the declaration fixture its
  kind-aware cases assume as a kernel-owned `map[string]string`
  (key → declared kind) without an import cycle: `internal/table`
  imports `internal/resolve` (`table/model.go::KernelTable` returns
  `resolve.Table`), so `internal/resolve` cannot import `table` and the
  fixture must not mention `table.TagDecl`.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence needed**: the import direction as stated, and a
    constructor shape for `guard.NewEvaluator` that accepts the plain
    mapping so the suite's fixture can drive it.
  - **If wrong**: the suite cannot state its own fixture and the
    contract cases move to `internal/guard` tests — losing the
    cross-implementer conformance property `0007:C1` placed in the
    kernel package.

## Proposed Solution

### Approach

Carry the declared kind to the value seam **by construction**: the
evaluator is built over the loaded model's tag declarations before any
resolution or lint pass, and looks the kind up by the atom's key at
evaluation time. The seam signature `Evaluate(atom, value)` — frozen by
`0007:C1` and its REQ-10 — is untouched; the kernel stays
declaration-free; JDR 0001 §D1's atom shape is untouched. What changes
is one package: `guard.Evaluator` stops being a zero-value `struct{}`
and gains a constructor taking a key → declared-kind mapping, and the
`eq`/`in` arms consult it — a held value or literal that does not parse
as the declared kind is `GuardUnevaluable`, never `GuardFalse`, which is
exactly the obligation `0007:A18` reserved and the shared conformance
suite states.

This is the carrier `0007:A17`'s own Evidence anticipated: "The declared
kind needed to *parse* the value is known to the evaluator from the
table it was built for, not from the kernel" ⇒ A17 verified the
*signature* needs no kind precisely because construction was always the
intended delivery path. The only casualty is 0003's zero-field
reflection assertion — a structural proxy for "never reads the tag
view", which the design preserves structurally by admitting only
declaration state (load-time model data), never a `resolve.TagSet` or
any runtime value (A2).

### Technical Design

#### Normative Contracts

**C1**

```normative
CARRIER. The declared kind travels to the value seam by CONSTRUCTION,
not by signature, atom field, or upstream refusal.

func NewEvaluator(kinds map[string]string) Evaluator

`kinds` maps tag key → declared kind token (the five-kind vocabulary
`enum | bool | int | set | scalar`, RDR 0003's spelling). The evaluator
holds this mapping and NOTHING else: no `resolve.TagSet`, no runtime
tag value, no view-typed state — "never reads the tag view" (`0003:C9`,
REQ-34) remains a property of the type. The zero-value `Evaluator{}`
construction form is retired; every non-test construction goes through
`NewEvaluator`. The mapping field is unexported, and Go still admits a
bare `Evaluator{}` literal — its disposition is fixed LOUD: a
nil-mapping evaluator answers GuardUnevaluable for every `eq`/`in`
atom (the no-declaration arm), surfacing as `guard_unevaluable`
refusals — never silently reverting to raw-string comparison. The seam signature `Evaluate(atom GuardAtom, value
string) GuardResult` (`0007:C1`, REQ-10) and the atom shape (JDR 0001
§D1: Key, Operator, Literal, Block) are UNCHANGED, and the kernel
(`internal/resolve`) continues to carry no declarations.
```

**C2**

```normative
TYPED COMPARISON. `eq` and `in` resolve the atom's key against the
constructed kind mapping and compare under the declared kind:

- kind `int`: held value AND literal (each member, for `in`) must parse
  as integers; comparison is over the PARSED values. A present held
  value that does not parse is GuardUnevaluable — never GuardFalse.
- kind `bool`: held value and literal must be the boolean tokens
  (`true` | `false`, the `0007:A27` spellings); comparison is token
  equality; any other held value is GuardUnevaluable.
- kind `enum` | `scalar`: every string parses; comparison is exact
  string equality (unchanged behavior).
- kind `set`: the operator/kind matrix does not admit `eq`/`in` over
  `set`; handed one anyway, the seam answers GuardUnevaluable
  (defensive — the matrix rejection at load stays the primary guard).
- key absent from the mapping: GuardUnevaluable — the seam cannot type
  the comparison (defensive arm; unreachable when C4 holds, per A1 /
  `0007:A7`). Undeclared-key ADMISSION policy — whether such a key may
  enter the view at all — is upstream and deliberately not decided
  here (see Joint-check).
- unparseable LITERAL: GuardUnevaluable (defense in depth; load-time
  rejection per `0003:C8` — "A predicate whose literal cannot be parsed
  as the declared tag kind MUST be rejected before resolution" — stays
  the primary guard, whose writer is the model loader).

`lt/lte/gt/gte` keep their operator-inferred integer parse (the matrix
admits only `int`); `contains` and the §D13 set arms are unchanged;
`exists` never reaches the seam (`0007:C1`).
```

**C3**

```normative
CONFORMANCE SUITE. `internal/resolve/guardcontract.go::
TestGuardEvaluatorContract` extends with kind-discriminating cases:
`eq`/`in` over an `int` and a `bool` dimension where the held value
does not parse (want GuardUnevaluable) and where it parses but differs
(want GuardFalse); plus the defensive arms — an `in` list with a
non-integer member on an `int` dimension, an integer-overflow held
value, `eq` over a `set`-kind key, and a key absent from the fixture
(each want GuardUnevaluable). The suite carries at least one
discriminating case per kind token, and a meta-check that the
fixture's kind tokens are exactly the five-kind vocabulary. The suite
publishes the declaration fixture its cases assume as a kernel-owned
`map[string]string` (key → kind token; no `internal/table` type
crosses — A5), and its doc states the caller obligation: construct the
seam under test over that fixture before invoking the suite.
```

**C4**

```normative
PRODUCER OBLIGATION. Every non-test evaluator construction site —
`internal/cli/flow_resolve.go::guardSeam`,
`internal/guard/product.go::valueSatisfies`,
`internal/graphlint/reach.go::atomAdmitsValue` — constructs over the
declarations of the SAME loaded model whose rows it evaluates: one
model, one evaluator, no mixed-model evaluation. This is what keeps
lint's stated property — "the same decision the runtime evaluator
makes" (`valueSatisfies` doc) — true after the seam becomes
declaration-aware.
```

No new refusal kind, error code, or envelope field: a typed
`GuardUnevaluable` reaches the user through the existing
`guard_unevaluable` refusal and its payload (`0007:C8`), reason
`uncomparable` ("the key was present and its value was not compared to
a verdict") — the CLI already maps it to `flow-guard-unevaluable`.

#### Load-Bearing Decisions

- **Identity** — the declared kind is a *function of the key* through
  the model's declaration; the atom's semantic identity tuple stays
  `(tag, operator, literal)` (RDR 0003) with no kind member, and no
  kind is ever stamped into an atom, payload, or hash.
- **Naming** — `NewEvaluator` (the Go-conventional constructor for the
  existing `Evaluator` type); rejected: `NewTypedEvaluator` (there is
  no untyped one left), a `Kind` field on `GuardAtom` (Alternative 2).
- **Selection / predicate** — typed comparison compares PARSED values
  for `int`, boolean TOKENS for `bool`, and raw strings for
  `enum`/`scalar` (C2). The
  narrower "validate-then-byte-compare" variant is the recorded
  fallback if A4's lint-agreement verification fails.

#### Illustrative Code

Illustrative — shape only, not load-bearing.

```go
// construction site (flow_resolve.go)
func guardSeam(m *table.Model) resolve.GuardEvaluator {
    return guard.NewEvaluator(guard.DeclaredKinds(m)) // key → kind
}

// eq arm sketch (grammar.go)
case "eq":
    switch ev.kinds[atom.Key] {
    case "int":
        held, err1 := strconv.Atoi(value)
        lit, err2 := strconv.Atoi(atom.Literal)
        if err1 != nil || err2 != nil { return resolve.GuardUnevaluable }
        return boolResult(held == lit)
    case "bool":  // token comparison over "true"/"false"
    case "enum", "scalar":
        return boolResult(value == atom.Literal)
    default:      // no declaration → cannot type
        return resolve.GuardUnevaluable
    }
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Declared-kind lookup | `internal/guard/declaration.go::DeclarationOf` | Takes `*table.Model`; seam constructor needs a plain map (A5) | Reuse (behind a small kinds-extraction helper) | No parallel kind model invented |
| Value seam | `internal/guard/grammar.go::Evaluator` | Zero-value `struct{}`; raw-string `eq`/`in` arms | Extend (constructor + typed arms, C1/C2) | The RDR's core change |
| View-freedom proof | `internal/guard/guard_evaluator_0003_test.go::TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView` | Zero-field check is stricter than REQ-34's text (A2) | Amend (field check narrows to "no view-typed/runtime state") | Recorded deviation against 0003's implemented proxy |
| Cross-implementer contract | `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` | No kind carrier, so no kind-aware cases exist | Extend (C3) | Suite grows fixture + cases |
| Malformed-set handling | `internal/guard/grammar.go::parseHeldSet` | Already answers unevaluable without declarations (the `072c7a0` inline fix) | Reuse unchanged | Proves the disposition; doesn't reach `eq`/`in` |

### Decision Rationale

Scored matrix (foundational profile). Criteria: **fit** — makes
`0007:C1`/`0007:A18`'s fixed semantics computable at `eq`/`in`;
**pins** — cost against pinned shapes (REQ-10 signature, REQ-34
statelessness proxy, §D1 atom, refusal taxonomy); **prior art** —
alignment with the instance reads and `0007:A17`'s evidence;
**structure** — guarantee holds for every `GuardEvaluator` consumer by
construction, not discipline; **blast** — modules/records touched.

| Option | fit | pins | prior art | structure | blast |
| --- | --- | --- | --- | --- | --- |
| O1 constructed evaluator (chosen) | full | amends one test-body proxy only | A17's evidence names it; BT.CPP binds types at the declared port | seam-level, suite-enforced | one package + 3 call sites |
| O2 widened signature | full | breaks REQ-10 + REQ-34 arity test + `0007:C1` fence | none — no peer passes type per-call | seam-level | kernel must carry declarations it never had |
| O3 kind on atom | full | breaks §D1 shape + `0007:C1` mirror fence | SCXML types values, not predicates | seam-level | every Row producer + payload mirror + fixtures |
| O4 boundary validation | partial — seam stays raw | needs a refusal path REQ-7 pins at five, or CLI-only | BT.CPP converts at READ, not at write | discipline — library callers bypass it | CLI + kernel entry + open JD-18's venue |

The deciding rows are **pins** and **structure**. O1 is the only option
whose collision is with a test-body proxy rather than a fenced contract:
REQ-10's signature, §D1's atom, and REQ-7's five-member refusal
taxonomy all survive verbatim. O4 additionally fails **fit**: with the seam still
raw-string, `TestGuardEvaluatorContract`'s "unevaluable, never false"
stays uncomputable for `eq`/`in`, and any `resolve.Resolve` caller that
is not the CLI (lint, graphlint, library use) keeps the silent
misroute — the guarantee would hold by producer discipline, the exact
class `0007:C1` rejects ("enforces the domain rule by structure").
O4 would also quietly answer JDR 0001 §JD-18 (the open view-conformance
venue, siblings 0003/0007) as a side effect of a bug fix — the wrong
home for that decision. O2 and O3 buy nothing over O1 (same fit,
same structure) at strictly higher pin cost.

Premortem: hardened — critic PASS with five mitigations folded: C2
gained the `set`-kind defensive arm and the undeclared-key-arm scope
note (admission is upstream, not decided here); C3 gained the
overflow/mixed-list/set-kind/absent-key cases plus the per-kind-case
and five-token fixture meta-checks; C1 fixed the zero-value
(nil-mapping) evaluator's disposition LOUD, with F5 recording the
missed-site failure mode; A4 gained the checked-in-model lint-diff
leg including authored-literal flips (see
`evidence/propose-premortem/critic.md`).
Ground-sweep: clean (24 anchors — 23 confirmed, 1 cosmetic
misattribution fixed in place: A17's If-wrong text had been credited
to A18 in Key Discoveries; see
`evidence/propose-premortem/ground-sweep.md`)
Joint-check: fired → 0030 (home: `cli/0003 §Normative Contracts` C6) —
the kind tokens `bool`/`enum`/`int`/`scalar`/`set` appear inside both
records' fences; both cite RDR 0003's closed vocabulary and neither
extends it. Recorded symmetrically by 0030's propose, 2026-09-19. The
other 11 peers are clear; nearest couplings are 0020, which shares the
runtime-admission locus upstream (undeclared *keys*; this RDR decides
declared-key *values* and pins no contract there), and open JDR 0001
§JD-18 (view-conformance venue), deliberately NOT decided here (O4's
rejection preserves it).

## Alternatives Considered

### Alternative 1: Widened seam signature

**Description**: change the seam to
`Evaluate(atom, value, kind string)` (or pass a declaration struct);
the kernel resolves the kind and hands it over per call.

**Pros**:

- Kind is explicit at the call — no hidden constructor state, trivially
  testable per invocation.

**Cons**:

- Breaks REQ-10's frozen `(atom, value)` signature, `0007:C1`'s fenced
  seam text, and REQ-34's arity assertion (`NumIn() != 3`) in one move.
- The kernel has no kind to pass: `resolve.Table` carries no
  declarations (`table/model.go::KernelTable` drops `m.Tags`), so the
  kernel seam itself must widen to carry declarations — the largest
  possible blast radius for the same behavioral result.

**Reason for rejection**: pays two fenced-contract breaks plus a kernel
widening to deliver information the evaluator can already own at
construction (matrix rows **pins**, **blast**).

### Alternative 2: Kind stamped on the atom

**Description**: the normalizer stamps each parsed atom with its key's
declared kind (`GuardAtom.Kind`); the kernel carries it opaquely; the
seam reads it.

**Pros**:

- Seam signature and evaluator statelessness both survive untouched;
  the kernel stays semantics-free (the field is opaque, like
  `Literal`).

**Cons**:

- Reopens JDR 0001 §D1's four-field atom shape and `0007:C1`'s fence
  ("the atom's four fields are named `Key`, `Operator`, `Literal`,
  `Block`" with the payload's `UndecidedAtom` mirroring them by name) —
  the mirror question reopens with it.
- Every `Row` producer must stamp: the normalizer, plus every test
  fixture constructing atoms by hand; an unstamped atom becomes a new
  undefined state the kernel can neither refuse nor default.

**Reason for rejection**: same fit and structure as the chosen carrier
at strictly higher pin cost — it trades one test-body proxy for a
cross-record shape fence plus a producer obligation on every fixture.

### Alternative 3: Upstream producer validation

**Description**: refuse a caller-supplied tag value that does not parse
as its key's declared kind before evaluation — at the CLI
(`--tag`/`--write` parsing) or at kernel entry
(`resolve.go::assemble`/`CheckInput`).

**Pros**:

- Malformed values never reach any comparison; the user error surfaces
  at its source with the offending flag named.

**Cons**:

- The kernel arm needs declarations the kernel does not carry, and a
  refusal path the taxonomy does not have (REQ-7 pins the refusal kinds
  at five; a new CLI error code is an unadjudicated envelope change).
- The CLI-only arm leaves every non-CLI `resolve.Resolve` caller (lint,
  graphlint, library use) unprotected, and it pre-empts JDR 0001
  §JD-18's open view-conformance question — the matrix's **fit** and
  **structure** rows (§Decision Rationale) carry the full argument.

**Reason for rejection**: fails **fit** (the seam contract stays
unenforceable) and **structure**, and pre-empts an open joint decision.
Boundary validation remains available LATER as §JD-18's answer — it
composes with, rather than replaces, a seam that is honest on its own.

### Briefly Rejected

- **Kind inference from literal bytes**: a literal that parses as an
  int does not make the dimension an int (`"7"` on an enum of
  `{"7","8"}`) — inference guesses where a declaration answers.
- **Two-valued seam with an error return**: `0007:A18`'s own If-wrong
  fallback — turns a caller data mistake into a kernel contract change
  and a programmer-error channel `0007:C1` already rejected (JDR 0001
  §D1 dissolved the error-channel question).
- **SCXML's disposition (fold to false + `error.execution` event)**:
  the instance read (W3C SCXML tests 309/344) confirms peers decide
  this fork, but fold-to-false is exactly the masking `0007:C2` forbids
  ("never to false and never to true"), and the kernel has no internal
  event queue to carry the side channel.

## Context

### Background

Raised independently by roborev job 6085 (`3e4e705`) and the
`f4e6b19..fda8e70` range review during `/rdr-implement-triage` for
RDR 0003; tracked as kata `intrastate#cq5p`. Scope review judged it not
shippable as a contained kata: each candidate fix is a cross-RDR
contract change (breaking `0007:A17`'s declaration-free seam and its
reflective boundary test, or adding an unadjudicated refusal path to
Resolve). The same defect class for set containment (held `null`) was
fixed inline during triage (`072c7a0`) because `parseHeldSet` needs no
declarations; the `eq`/`in` arms cannot know the target kind from the
atom alone, so that fix shape does not reach them. Lint never produces
such a value (it renders from the declared domain), so lint and runtime
diverge exactly where the runtime accepts caller-supplied strings with
zero validation (`internal/resolve/resolve.go::assemble`/`merge`).

### Technical Environment

Go module `github.com/cwensel/intrastate`. Seam under decision:
`internal/guard/grammar.go::Evaluator`, implementing the exported
`resolve.GuardEvaluator` interface; shared conformance suite
`internal/resolve/guardcontract.go::TestGuardEvaluatorContract`;
producers `internal/resolve/resolve.go::assemble`/`merge`. Governing
records: RDR 0007 (C1, A17, A18, REQ-10), RDR 0003 (view-conformance
clause), JDR 0001 D1 (atom shape).

## Research Findings

### Investigation

Read before enumerating: the seam and its arms
(`internal/guard/grammar.go::Evaluator` — `eq`/`in` compare raw
strings; `lt..gte` parse; `contains`/`parseHeldSet` already answer
unevaluable), the governing fences (`0007:C1` seam text and REQ-10;
`0007:A17`/`0007:A18` with their If-wrongs; `0003:C8`/`0003:C9`; JDR
0001 §D1), the kernel's declaration-freedom
(`table/model.go::KernelTable` drops `m.Tags`; `resolve.Table` has no
declaration field), the three non-test construction sites (C4), and the
enforcement tests (`guardcontract.go::TestGuardEvaluatorContract`;
`guard_evaluator_0003_test.go::TestReq34…` — the zero-field check).
Prior-art pass (bounded, recorded in `evidence/research/prior-art.md`):
W3C SCXML tests 309/344 (guard evaluation error → treat as `false` +
`error.execution` event — the fold this project's kernel rejects) and
BehaviorTree.CPP port typing (`convertFromString<T>()` converts a
string-carried value against the PORT's declared type at read; an
inconvertible string is a loud run-time error) — both peers bind type
information to a declaration the evaluator was constructed against,
none passes it per call. ⚠ no prior-art coverage in the available
corpora for compile-checked expression environments (CEL-class); the
class claim rests on the two anchored instances plus `0007:A17`'s
quoted evidence.

Sibling-path check (new discriminator: declared-kind lookup at
evaluation): the adjacent lint path already makes this exact decision —
`internal/guard/declaration.go::DeclarationOf` reads a key's declared
kind out of a loaded model, and lint's product renders held values from
the declared domain. The design reuses that signal (audit row 1)
rather than inventing a parallel kind model; searched, no second
kind-lookup path exists.

### Key Discoveries

- **Documented** — `0007:A17`'s Evidence already names construction as
  the delivery path (quoted in §Approach): A17 froze the signature,
  not the constructor, so the chosen carrier is the one the verifying
  record itself anticipated.
- **Documented** — `0007:A18` and `0007:C1` ("MAY answer unevaluable
  for a present value it cannot compare (A18)") fix the semantics;
  `0007:A17`'s If-wrong names the fallback carrier ("e.g. the declared
  kind on the atom … but never the view") — only the carrier was
  deferred.
- **Documented** — `resolve.Table` carries `Revision/Outcomes/Rows`
  only, with no declaration field (the fact Alternative 1's rejection
  turns on).
- **Documented** — the inline `072c7a0` fix (`parseHeldSet` nil check)
  proves the disposition for the declaration-free arms; `eq`/`in`
  cannot reach it without a kind (Background).
- **Verified** — grep over non-test sources: exactly three
  `Evaluator{}` construction sites (C4's list); no other consumer
  exists to inherit the constructor change.
- **Assumed** — typed numeric equality agrees with lint's rendered
  product (A4); the fallback (validate-then-byte-compare) is recorded
  in Load-Bearing Decisions.

## Trade-offs

### Consequences

- Positive: the fixed semantics become computable and suite-enforced —
  a malformed caller value on a typed dimension refuses
  `flow-guard-unevaluable` instead of silently pruning or misrouting;
  every `GuardEvaluator` consumer inherits the fix, not just the CLI.
- Positive: kernel, atom shape, refusal taxonomy, and envelope are
  byte-untouched — the change is one package plus three call sites.
- Negative: zero-value `Evaluator{}` construction is retired with no
  compatibility shim (project convention: no back-compat); all sites
  and fixtures migrate in one change.
- Negative: 0003's statelessness proxy (the zero-field reflection
  check) narrows — a recorded deviation against a closed RDR's
  implemented test, justified by A2.
- Negative: verdicts flip for previously "deciding" comparisons —
  malformed values (False → Unevaluable, the fix) and, under C2's
  parsed comparison, non-canonical numerals (`"07" eq 7`:
  False → True) — A4 bounds the blast.

### Risks and Mitigations

- **Risk**: mixed-model construction — an evaluator built over model
  A's declarations evaluating model B's rows types the wrong kinds.
  **Mitigation**: C4's one-model-one-evaluator obligation at each site;
  the sites all hold exactly one loaded model today (A3).
- **Risk**: typed numeric equality diverges from lint's product on a
  canonical-form edge.
  **Mitigation**: A4's verification before lock; recorded fallback to
  validate-then-byte-compare in Load-Bearing Decisions.
- **Risk**: the suite's fixture kinds drift from the evaluator's
  five-kind vocabulary (a sixth kind added later, suite silent).
  **Mitigation**: the fixture uses RDR 0003's kind tokens verbatim
  (C1); a vocabulary change is a 0003-successor contract change that
  re-enters here via the suite cases.

### Failure Modes

- **Visible break**: `flow resolve --tag iter=many` on an `int`
  dimension now exits with `flow-guard-unevaluable`, payload naming the
  row and atom with reason `uncomparable` — previously a silent prune
  or a wrong sibling selection. Diagnosis reads the refusal payload;
  recovery is fixing the caller's value.
- **Silent failure guarded against**: the misroute itself — a
  present-but-malformed value deciding `GuardFalse` and pruning a row
  the caller believes was compared.
- **Residual silent mode**: a construction site violating C4 (wrong
  model's kinds) can still type a comparison wrongly without refusing —
  guarded by the site audit (A3) and the conformance suite's
  fixture-bound cases, not by the kernel.
- **F4 Suite/fixture drift**: an implementer constructs over its own
  kinds instead of the published fixture and the kind-aware cases pass
  vacuously (everything scalar → string compare). The suite's
  discriminating cases (int/bool want-Unevaluable legs) fail against a
  scalar-typed fixture, which is the detection.
- **F5 Zero-value evaluator survives migration**: a missed construction
  site keeps `Evaluator{}` (nil mapping) — every `eq`/`in` atom it sees
  answers GuardUnevaluable, so the miss surfaces as loud
  `guard_unevaluable` refusals on previously-deciding guards (C1's
  fixed disposition), never as a silent raw-string revert. Diagnosis:
  the refusal payload's `uncomparable` reasons on well-formed values;
  recovery is constructing per C4.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A5)

### Minimum Viable Validation

1. Author a model declaring `iter` as `int`, single-valued, with a row
   guarded `iter eq 3` and an unguarded fallback row on the same
   outcome.
2. `flow resolve --tag iter=many --outcome <o>` → refuses
   `flow-guard-unevaluable`; the payload names the guarded row and the
   `iter` atom with reason `uncomparable`.
3. `flow resolve --tag iter=4 --outcome <o>` → decides: the guarded row
   prunes (GuardFalse) and the fallback plans.
4. `TestGuardEvaluatorContract`, extended per C3 and driven against
   `guard.NewEvaluator` constructed over the published fixture kinds,
   is green — including the int/bool want-Unevaluable legs that fail
   against today's raw-string arms.

End-state: step 2's invocation is the seed defect
(kata `intrastate#cq5p`) refusing loudly instead of misrouting.

### Phase 1: Seam constructor and typed arms

Intent: `guard.NewEvaluator` over a key→kind mapping; C2's typed
`eq`/`in` arms; zero-value construction retired; REQ-34 test's field
check narrowed per A2 (the recorded deviation).

### Phase 2: Conformance suite extension

Intent: publish the fixture kinds map in `internal/resolve`, add the
kind-discriminating cases (C3), keeping the suite free of
`internal/table` types (A5).

### Phase 3: Producer alignment

Intent: the three C4 sites construct over their own loaded model's
declarations (reusing `DeclarationOf`/a kinds-extraction helper);
test fixtures migrate to the constructor.

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

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
