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
- **Profile**: foundational — C1, the declared-kind carrier at the value seam, which RDR 0030 and every `GuardEvaluator` implementer consume; user-facing yes; locks cross-rdr
- **Priority**: Medium
- **Related Issues**: kata `intrastate#cq5p` (1528)
- **Predecessors**: 0007-guard-predicate-totality,
  0003-guard-predicate-exhaustiveness
- **Seam Lineage**: no prior accretion

## Problem Statement

A reader handing back an owned tag value onto an `int` or `bool`
dimension expects a value the declared kind cannot parse to refuse
loudly — `flow-guard-unevaluable` — the same three-valued honesty the
seam guarantees for absent facts. Today `eq` and `in` compare the held
value as a raw string (`internal/guard/grammar.go::Evaluator`), so a
present-but-malformed value decides `GuardFalse`: the runtime silently
prunes a candidate row or selects a sibling/`unless` row it never
actually compared. For a consumer like the rdr navigator that is a
silent misroute — worse than the loud refusal the seam was designed to
prefer.

The defect's home is the OWNED ingress, and only it. Of the three ways
a value reaches a comparison, two are already conformed to their
declaration before resolution — model literals at load
(`internal/table/load.go::conformKind`) and CLI `--tag`/`--write`
values (`internal/cli/flow_input.go::canonicalValue`, which is why
`--tag iter=many` refuses `flow-tag-invalid` at input and never reaches
the seam, fixture S8). Reader-supplied owned values are the sole
unconformed door: `internal/accessor/model.go::ReadResult.OwnedSnapshot`
builds `resolve.Tag` straight from the reader's bytes and
`internal/resolve/resolve.go::assemble` merges them
`ProvenanceOwned` with no kind check anywhere on the path. That door is
not a narrow one — every persisted owned value and every hand edit to
the artifact enters through it (`JDR 0004 §JD-2`, whose ingress census
reaches the same table).

This is the record's own O1-vs-O4 argument in miniature: boundary
validation is discipline at one ingress, and the owned ingress proves
the discipline is not total. A seam honest on its own is the only
structural guarantee.

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
  - **Status**: Verified
  - **Method**: Peer RDR + Source Search
  - **Evidence**: `0007:A7` (Verified) plus
    `internal/table/normalize.go::(*loader).atomsFromBlock`, whose contract is
    block-agnostic — "the tag-key declaration rule are enforced
    identically in all three blocks" (match, `guard.all`,
    `guard.unless`) — and refuses an undeclared key with
    `CatUnknownTag` ("references the undeclared tag"). The refusal
    therefore fires for guard atoms, not only match keys.
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
  - **Status**: Verified
  - **Method**: Peer RDR + Source Search
  - **Evidence**: `0003:C9`/REQ-34 fence reads "MUST decide value
    semantics over a present value only … the evaluator MUST NOT read
    the tag view" — neither "stateless" nor "zero fields" appears in
    0003's normative or req-list text. The field check exists only in
    the test body
    (`internal/guard/guard_evaluator_0003_test.go::TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView`,
    a `reflect` `NumField() != 0` assertion), confirming it is a
    structural proxy rather than a fenced obligation.
  - **If wrong**: statelessness is itself fenced and the choice must
    move to the atom-stamp carrier (Alternative 2), reopening JDR 0001
    §D1's atom shape instead.
- **A3 Every non-test construction site of the evaluator has the loaded
  model (or its declarations) in scope:
  `internal/cli/flow_resolve.go::guardSeam`,
  `internal/guard/product.go::valueSatisfies`,
  `internal/graphlint/reach.go::atomAdmitsValue`.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `rg` over non-test sources returns exactly the three
    named `Evaluator{}` sites (`internal/cli/flow_resolve.go::guardSeam`,
    `internal/guard/product.go::valueSatisfies`,
    `internal/graphlint/reach.go::atomAdmitsValue`) — no fourth site on
    `main`. In all three the `*table.Model` is in scope at an ancestor
    frame, never at the immediate site: `valueSatisfies` is called from
    `Denotation(m *table.Model, …)` one frame up; `atomAdmitsValue` sits
    two frames below `matchSatisfiable(m *table.Model, …)`; `guardSeam()`
    is zero-arg, and of its three call sites two hold `req.model`
    directly (`internal/cli/flow_resolve.go` at the resolve path and at
    `escapeClassOf`, which resolves a probe derived from the same
    model's `KernelTable()`), while the third —
    `internal/cli/flow_next.go::probeRow` — takes only a `table.Row`,
    with the model one frame up at its caller. Of the three measured
    sites C4 migrates TWO: the guard package needs one added local
    parameter, and the CLI additionally widens `probeRow`'s signature
    or passes the constructed evaluator down. Both are intra-package —
    no cross-package signature surgery. The graphlint site is measured
    here but deliberately NOT migrated (C4, match path): both its
    callers are reached only through `BlockMatch` atoms, where the
    kernel byte-compares, so typing it would create the divergence this
    record closes. This assumption measures the site population; which
    of them C4 obliges is C4's call.
  - **If wrong**: a consumer without a model needs a declaration-free
    degraded evaluator, re-introducing the raw-string arms the RDR
    exists to remove — the constructor design must change (e.g. thread
    the model through the caller).
- **A4 With `int` spellings canonicalized at the PREDICATE ingress (C5)
  and parsed comparison confined to the GUARD path (C2, C4),
  lint and the runtime agree on every
  `eq`/`in` GUARD atom over an `int` dimension: lint's product renders
  held values
  from the declared domain with `strconv.Itoa`, the authored literal is
  admitted only in that same canonical spelling, and a model spelling an
  int non-canonically refuses at load naming the canonical form rather
  than silently changing a verdict. MATCH atoms are outside this claim:
  both deciders byte-compare them, so they agree without typing.**
  - **Status**: Verified — under the 2026-09-20 author ruling (option
    (ii)) the claim is scoped to the guard path, which is the only
    path on which lint and the kernel share a decision. The held-side
    leak this assumption was re-opened for lives entirely on the MATCH
    path (`reach.go::heldValues` → `canonicalValues`, no integer
    round-trip), and that path is no longer typed by this record — C4
    excludes `reach.go::atomAdmitsValue`, which keeps byte-comparing
    exactly as the kernel's `TagSet.matches` does. A non-canonical
    `[initial]`/`[rule.write]` int cell therefore reaches lint and the
    kernel in the SAME spelling and is read the same way by both; it
    is not a divergence this record introduces or leaves. It remains a
    canonicalization gap at that ingress, recorded as the Residual in
    §Load-Bearing Decisions and carried by kata `intrastate#ch99`
    (the wider
    rule contradicts RDR 0024's locked REQ-7/REQ-9, so it is not
    resolvable here).
  - **Method**: Source Search + Spike
  - **Evidence**: HELD side —
    `internal/guard/assignment.go::valueAssignments` (called from
    `internal/guard/product.go` with `m.Tags[key]`) renders the declared
    domain with `strconv.Itoa` in its `int` arm, which emits no leading
    zero, `+`, or whitespace, so only canonical spellings reach the lint
    side by THAT path (spike
    `evidence/spikes/a4-render-path.md`). That spike measured the
    product path, which under the narrowed scope is the whole held
    side of this claim: `product.go::valueSatisfies` is reached from
    `Denotation`, and `product.go::selectionOf` routes `BlockMatch`
    atoms away from it. The other held-side path —
    `internal/graphlint/reach.go::heldValues` → `canonicalValues` (no
    integer round-trip) — feeds `reach.go::atomAdmitsValue` ONLY,
    which C4 excludes from typing, so lint and the kernel read those
    cells identically and this claim does not quantify over them.
    Reconcile verified that exclusion is structural, not incidental
    (both callers are `BlockMatch`-only; derivation in
    `evidence/reconcile/report.md` §A4). LITERAL side — today the
    literal is joined from authored bytes (`valueSatisfies`) and the
    loader's kind check is `strconv.Atoi`
    (`internal/table/load.go::conformKind`), which ACCEPTS `"00"`,
    `"01"`, `"+1"`; measured, an authored `n eq "00"` flips lint from
    blocking `graph-coverage-gap` (exit 2) to clean (exit 0) — the
    blocking→clean direction (fixture S7; spike
    `evidence/spikes/a4-lint-diff.md`). C5 closes that leg by admitting
    only `strconv.Itoa(n) == authored`, so the divergence the spike
    measured cannot be authored at all. Bare floats are covered: a TOML
    `eq = -0.0` is rendered to the literal `"-0"` by
    `internal/table/load.go::valueMembers`
    (`strconv.FormatFloat(t, 'g', -1, 64)`) BEFORE the kind check, and
    `Atoi("-0")` succeeds while `Itoa(0)` is `"0"` — C5 refuses it,
    `conformKind` alone does not. The 123-model corpus shows no verdict
    flip, but it contains ZERO `eq`/`in` atoms over an `int` tag, so
    that result measures coverage, not safety (S6).
  - **If wrong**: lint and runtime disagree on a canonical-form edge the
    load refusal does not cover, and the fix narrows to
    validate-then-byte-compare (parse to prove comparability, compare
    raw bytes) — the fallback recorded in §Load-Bearing Decisions.
    Surfaces as a model whose `lint` verdict and `flow resolve` verdict
    disagree over an `int` guard, caught by S6's before/after corpus
    diff and S7's regression case.
- **A5 The conformance suite can publish the declaration fixture its
  kind-aware cases assume as a kernel-owned `map[string]string`
  (key → declared kind) without an import cycle: `internal/table`
  imports `internal/resolve` (`table/model.go::KernelTable` returns
  `resolve.Table`), so `internal/resolve` cannot import `table` and the
  fixture must not mention `table.TagDecl`.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/model.go` imports `internal/resolve`
    and `model.go::KernelTable` returns `resolve.Table`; a grep for
    `internal/table` across `internal/resolve`'s non-test sources
    returns zero hits, so the direction is one-way as stated and a
    kernel-owned `map[string]string` fixture mentions no `table` type.
  - **If wrong**: the suite cannot state its own fixture and the
    contract cases move to `internal/guard` tests — losing the
    cross-implementer conformance property `0007:C1` placed in the
    kernel package.
- **A6 The installed base of persisted owned values contains no value
  whose guard verdict this RDR changes SILENTLY — that is, no held
  value that is non-canonical but parseable (`"07"` against
  `n eq 7`), which today prunes as a raw-string mismatch and after C2
  matches as a parsed equality.**
  - **Status**: Verified — census run at reconcile, over the population
    the 2026-09-20 ruling narrows this to: GUARD atoms only, since the
    match path is no longer typed (C4) and so changes no verdict at
    all.
  - **Method**: Spike (corpus scan)
    S6 measured the AUTHORED side (123 models,
    zero `eq`/`in` atoms over an `int` tag). The held side is a
    different population, and the flip there runs BOTH ways —
    malformed → `GuardUnevaluable` is the loud fix this RDR wants, but
    non-canonical-yet-parseable → `GuardFalse` becoming `GuardTrue` is
    a silent verdict change, precisely the harm the RDR exists to
    remove. C5 does not reach reader values (`trace` step 3), so
    nothing upstream prevents it.
  - **Evidence**: the census is bounded because an owned value can only
    reach a typed comparison through a tag that is BOTH
    `provenance = "owned"` AND declared `int` or `bool`. Across the
    committed models exactly one such dimension exists:
    `models/rdr.toml::[tags.gate_passed]` (`owned`, `bool`). The other
    two `int`/`bool` declarations in the corpus —
    `models/examples/release-grammar.toml::[tags.risk]` (`int`) and
    `::[tags.urgent]` (`bool`) — are `provenance = "observed"`, which
    carries no persisted owned value. Every authored `gate_passed`
    value is the canonical token: `[initial] gate_passed = "false"`,
    the two partitioning guards `[rule.guard.all.gate_passed] eq =
    "true"` and `eq = "false"`, and one `[rule.write] gate_passed =
    "false"`. Canonical `bool` tokens compare identically under
    raw-string and token equality, so no committed model's guard
    verdict changes — silently or loudly. The population of
    non-canonical-yet-parseable held values on a typed dimension is
    EMPTY on `main`, which is the claim.
  - **If wrong**: the change needs a release note and possibly a scan
    or `--fix` affordance before rollout, rather than shipping as a
    behavior fix — and the silent-match cases must be enumerated, since
    an RDR premised on removing silent verdict changes cannot introduce
    an unbounded set of them. Residual risk is a deployment whose
    persisted artifact holds an out-of-corpus owned value on a typed
    dimension; S9's owned-reader fixture is the standing regression
    case for that direction.
- **A7 RETIRED — its subject no longer exists in this design.** The
  assumption quantified over a widening of graphlint's atom check from
  `bool` to `resolve.GuardResult`. Under the 2026-09-20 author ruling
  (option (ii)) C4 excludes `reach.go::atomAdmitsValue` from the typed
  construction sites, so that widening is not proposed, and neither is
  the `atomVerdictForValue` rename it was attached to.
  - **Status**: Verified — as a retirement: the claim's subject is
    absent from the design, which is a checkable fact about this
    record rather than a deferral. The rejected alternative is the
    widening itself (option (i)'s three-valued `atomAdmitsValue`),
    declined in C4 with its reason.
  - **Method**: Design Decision
  - **Evidence**: retirement is verified, not asserted. The claim's
    subject was a change to `atomAdmitsValue`'s return type; the
    narrowed C4 leaves that function at its `main` shape (`bool`,
    `guard.Evaluator{}`, byte comparison), and no other clause in this
    record proposes touching it — `atomVerdictForValue` appears
    nowhere in the record after the narrowing. The assumption is
    therefore RETIRED rather than re-scoped: there is no residual
    claim to re-verify at a narrower home, because the population it
    quantified over (lint verdicts changed by the widening) is empty
    by construction. Reconcile also confirmed WHY the exclusion is
    safe, which is the reason the retirement is not a loss of
    coverage: both callers of `atomAdmitsValue`
    (`reach.go::ownedAtomSatisfiable` via `matchSatisfiable`, and
    `analysis.go::nodeMeetsAll` via `satisfiesSomeTerminal` over
    `model.Terminal`) are reached only through `BlockMatch` atoms, so
    lint's match admission stays byte-equal to the kernel's
    `resolve.go::TagSet.matches` with no widening at all.
  - **If wrong**: if a later change does widen that check, the
    opposite-polarity hazard this assumption named is real and must be
    re-opened — the two callers read the result with OPPOSITE polarity
    (existential at `reach.go`, exclusionary at `analysis.go`), so a
    sign error SUPPRESSES a finding rather than producing a wrong one,
    which no corpus diff over firing findings detects. C4 records that
    hazard as the stated reason the widening is declined, so the
    knowledge survives the retirement.

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
intended delivery path. Note the scope this narrows: A17's headline
reads "needs neither the view, nor provenance, nor sibling atoms, nor
tag declarations **at evaluation time**", and C2's arm does consult a
declared kind at evaluation time — keyed by the atom, off state fixed at
construction. Read strictly, the headline's last clause is narrowed to
"needs no declaration it was not built over"; A17's own Evidence and its
If-wrong ("the seam takes whatever RDR 0003 shows the evaluator needs")
are what license the narrowing. Recorded here as a deviation against a
Verified peer assumption's wording so a cluster pass reads it as
adjudicated, not as a cross-record contradiction. The only other
casualty is 0003's zero-field
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
`enum | bool | int | set | scalar`, RDR 0003's spelling). Its producer
is EXPORTED and singular:

func DeclaredKinds(m *table.Model) map[string]string

exported because the two C4 construction sites span two packages
(`internal/cli`, `internal/guard`), so an
unexported helper cannot serve them and two hand-rolled loops would
be two chances to disagree. It is TOTAL over `m.Tags`: every declared
key gets an entry, and a key whose `Kind` is the empty string is
OMITTED rather than mapped to `""` — an empty-string entry would be
indistinguishable at C2's arms from an absent key, collapsing two
diagnostics into one. It takes the model rather than a
`map[string]TagDecl` (the shape `product.go::CanRefuse` uses) because
both construction sites hold a `*table.Model` at an ancestor
frame (A3) and none holds a bare decls map; the decls-map form is not
the pattern here. The evaluator
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
§D1: Key, Operator, Literal, Block) are UNCHANGED, and the kernel's
RESOLUTION PATH — `Resolve` and its callees — continues to read no
declarations. The rule governs that path, not `internal/resolve` as a
namespace, which may host the seam implementation with its declaration
state (`JDR 0004 §JD-1`).

The atom's fields keep their `main` types, which this clause reads and
does not restate as new: `Key`, `Operator` and `Literal` are bare
`string`, `Block` is the named string type `resolve.Block`
(`BlockAll | BlockUnless | BlockMatch`) that `internal/table` aliases.
`DeclaredKinds` is homed in `internal/guard` — the package that owns
the evaluator it feeds, and the one both C4 sites already import;
`internal/table` is import-legal but would put a guard-shaped producer
in the loader's package. Handed a nil model it returns an empty
non-nil map rather than panicking: that is the same LOUD-not-crash
disposition as the nil-mapping evaluator above, and it makes the two
failure paths converge on one diagnosis — every `eq`/`in` answers
GuardUnevaluable — instead of one refusing and one crashing.
```

**C2**

```normative
Surface — of C1; states how the carried kind is applied at the comparison.
TYPED COMPARISON APPLIES ON THE GUARD PATH. `eq` and `in` GUARD atoms
resolve the atom's key against the
constructed kind mapping and compare under the declared kind. MATCH
atoms do NOT: the kernel byte-compares them at
`internal/resolve/resolve.go::TagSet.matches`, and lint's match
admission (`internal/graphlint/reach.go::atomAdmitsValue`) keeps doing
the same, so the two agree without typing (C4 excludes that site).

This split is deliberate and is named here so it is discoverable
rather than surprising: ONE declared kind, TWO comparison rules,
divided by whether lint and the kernel share a decision path. They
share it for guards and not for matches, which is a structural fact
about the code rather than a convention. The reasoning, the rejected
alternative, and the residual it leaves are in §Load-Bearing
Decisions ("The guard/match split"). Everything below governs the
guard path:

- kind `int`: held value AND literal (each member, for `in`) must parse
  as integers; comparison is over the PARSED values. A present held
  value that does not parse is GuardUnevaluable — never GuardFalse.
  For `in`, ONE unparseable member poisons the whole list:
  GuardUnevaluable even when the held value equals a member that did
  parse. Parse-all-then-compare, not parse-and-compare-per-member —
  the permissive reading would let lint and the runtime disagree over
  a mixed list, which is the divergence `D-canonical-spelling-at-load`
  exists to prevent. The parse is `strconv.Atoi`, whose width is Go's
  platform-native `int` — 64-bit on the supported build targets, and
  nothing on `main` pins it (no width contract exists in the guard
  package or the conformance suite). The conformance suite therefore
  states no expected verdict for a value that overflows on one target
  and not another: an out-of-range spelling is GuardUnevaluable on the
  target under test, and the suite pins the in-range cases only. Fixing
  the width is out of scope here — it is `conformKind`'s reading as
  much as the seam's, and changing one without the other re-opens the
  divergence this RDR closes.
- kind `bool`: held value and literal must be the boolean tokens
  (`true` | `false`, the `0007:A27` spellings); comparison is token
  equality; any other held value is GuardUnevaluable. The rejected set
  is every other spelling — `"1"`, `"0"`, `"True"`, `"yes"` — and it is
  reachable only from the owned door: `internal/cli/flow_input.go::
  canonicalValue` already refuses non-`true`/`false` bool values before
  resolution, so the CLI cannot produce one. A6's census covers this
  population as well as the `int` one, and it is a LOUD change
  (previously-comparing values start refusing), not a silent one.
- kind `enum` | `scalar`: every string parses; comparison is exact
  string equality (unchanged behavior). These two SHARE one arm — no
  input distinguishes them, so the conformance suite's
  one-case-per-kind-token requirement discriminates each against
  `int`/`bool`, never against the other.
- kind `set`: the operator/kind matrix does not admit `eq`/`in` over
  `set`; handed one anyway, the seam answers GuardUnevaluable
  (defensive — the matrix rejection at load stays the primary guard).
- key absent from the mapping: GuardUnevaluable — the seam cannot type
  the comparison. This arm is DEFENSIVE for the undeclared key (A1 /
  `0007:A7` make that unreachable when C4 holds) but LIVE for a second
  population A1 does not cover: a key that IS declared with `Kind` the
  empty string. C1 omits those from the mapping deliberately, and on
  `main` nothing rejects them — `internal/table/load.go::conformKind`
  switches on `decl.Kind` with cases for `int` and `bool` and NO
  `default`, so an empty-Kind tag's members pass with no validation at
  all (`conformDomain` likewise). For such a key this is the ordinary
  arm, not a defensive one: `eq`/`in` over it answers GuardUnevaluable
  where today it compares raw strings. That is the intended loud
  disposition — an undeclared kind cannot type a comparison — but it is
  a behavior change on a reachable input, so the conformance suite
  pins it as a named case rather than filing it under the unreachable
  defensive arms. Undeclared-key ADMISSION policy — whether such a key
  may enter the view at all — is upstream and deliberately not decided
  here (see Joint-check).
- unparseable LITERAL: GuardUnevaluable (defense in depth; load-time
  rejection per `0003:C8` — "A predicate whose literal cannot be parsed
  as the declared tag kind MUST be rejected before resolution" — stays
  the primary guard, whose writer is the model loader).

Dispatch over the kind token is EXHAUSTIVE across the five-kind
vocabulary, and the fallthrough `default` means GuardUnevaluable — the
key-absent and `set` arms above both want that verdict. A `default`
meaning string equality would satisfy the `enum`/`scalar` cases while
silently disarming both defensive arms, so the suite asserts an
unknown kind token answers GuardUnevaluable rather than comparing.

`lt/lte/gt/gte` keep their operator-inferred integer parse (the matrix
admits only `int`); `contains` and the §D13 set arms are unchanged;
`exists` never reaches the seam (`0007:C1`). Their reading is fixed by
the OPERATOR, not the kind: when C3's re-keying moves the five `gte`
cases onto the `int` key and the four `contains` cases onto the `set`
key, their verdicts are unchanged — the kind lookup is not consulted
for those operators at all, so `contains` over a `set` key does NOT
become GuardUnevaluable under the `set` arm above (which governs
`eq`/`in` only). Operator-inferred and declared-kind readings never
compete: the declared kind changes the reading of `eq` and `in`, and
nothing else.

SHAPES THIS CLAUSE READS AND DOES NOT CHANGE. The arms above quantify
over per-member parsing, a three-valued verdict and an operator
dispatch, each of which already has a fixed shape on `main`. Naming
them here is what makes the clause reconstructible; none is amended.

- `in` MEMBERS. The atom's `Literal` is a single `string`
  (`internal/resolve/guard.go::GuardAtom`), and an `in` list travels
  through it as a JSON array, in whatever form the two producers on
  `main` already render — `internal/guard/assignment.go::renderSet`
  and `internal/graphlint/reach.go::renderSetLiteral`, both sorting,
  deduping and compact-encoding the members per JDR 0001 §D13. That
  encoding is those functions' to define and this clause makes no
  exactness claim of its own about it; the conformance suite's
  existing `in/member` case (`internal/resolve/guardcontract.go`)
  is the fixture that pins it.
  "Each member, for `in`" above means the members decoded from that
  array, and parse-all-then-compare quantifies over them. No new
  members field and no delimiter convention is introduced.
- `GuardResult` ENCODING. `resolve.GuardResult` is an `int` whose
  constants are declared `GuardFalse = iota, GuardTrue,
  GuardUnevaluable` (`internal/resolve/resolve.go`), so the ZERO value
  is `GuardFalse`. That order is deliberately NOT the strong-Kleene
  truth order (`internal/resolve/guard.go` documents the difference);
  `kleeneAnd`/`kleeneNot` implement the truth table directly and never
  compare raw constant values. This RDR fixes no ordering and adds no
  verdict — a reconstruction that re-derives the iota order, or
  reorders it to match the truth order, is wrong.
- OPERATOR DISPATCH. `Operator` is a bare `string`, not a named type,
  and only `exists` has a declared constant (`resolve.OpExists`). The
  kind lookup is reached from the `eq`/`in` arms; an operator the seam
  does not recognize — `exists` included, which `0007:C1` keeps off
  this path — answers GuardUnevaluable like any other atom it cannot
  type. The seam NEVER panics: C1 fixes the nil-mapping disposition as
  LOUD-not-crash, and that posture is the seam's generally, so an
  unreachable-by-contract operator is a verdict, not an
  `unreachable` panic.
- STEP ORDER. The unparseable-LITERAL arm is a per-kind obligation
  inside the `int` and `bool` arms, not a pre-dispatch gate above the
  kind switch. The kind lookup runs FIRST: an unknown or absent kind
  answers GuardUnevaluable through the `default` without the literal
  being parsed at all. Hoisting the literal parse above the switch
  would make the literal's shape decide a verdict the declared kind
  owns, and reverses this clause's order of authority.
```

**C3**

```normative
Surface — of C1; the cross-implementer gate proving the carrier types comparisons.
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
crosses — A5). It is published as a FUNCTION returning a fresh map:

func ContractKinds() map[string]string

not a package-level `var`. Nothing like it exists on `main`, so the
form is this RDR's to fix, and an exported map var is mutable shared
state every importing test package could write — one package's mutation
would silently re-type another's comparisons, which is precisely the
vacuous-pass mode F5 names. A constructor hands each caller its own
copy. Delivery is STRUCTURAL, not a documented obligation —
the suite takes a constructor and builds the seam itself:

func TestGuardEvaluatorContract(t *testing.T,
        newSeam func(kinds map[string]string) GuardEvaluator)

so a caller cannot construct over the wrong map and still compile.
A prose "construct over the published fixture first" obligation beside
the existing `(t, seam)` signature would be unenforceable, and F5's
vacuous-pass mode is exactly a caller supplying its own kinds. This
widening breaks both in-tree call sites
(`internal/resolve/guard_mvv_test.go` at the `conformingContractSeam`
and `recordingSeam` drivers) and they migrate in the same phase.

Because that fixture types comparisons BY KEY, one key carries exactly
one kind, so each case names its own key drawn from the fixture, chosen
so the case's operator is one the matrix admits for that kind. Today
every case is built with `Key: "subject"`; re-keying routes the
existing `eq`/`in` cases onto an `enum`/`scalar` key, the five `gte`
cases onto the `int` key, and the four `contains` cases onto the `set`
key, with the new legs on the `int`, `bool` and `set` keys. Without
per-case keys the one-case-per-kind-token requirement above is
unsatisfiable; keying an operator onto a kind the matrix rejects would
make the fixture itself unauthorable.

Both in-tree implementers are bound, not only `guard.Evaluator`: the
kernel package's own reference seam
`internal/resolve/guard_mvv_test.go::conformingContractSeam` is driven
by this suite directly and again wrapped in `recordingSeam` for
`0007`'s REQ-71 tests. It is a raw-string `struct{}` today, so it is
made kind-aware over the same published fixture in the phase the new
cases land, or those tests fail. Concretely it gains a
`kinds map[string]string` field and a matching constructor to satisfy
the `newSeam` parameter above; `recordingSeam` keeps delegating to its
`inner` and needs no kinds of its own.

`guard.Evaluator`'s existing suite call is RE-POINTED, not duplicated.
On `main` the contract test has THREE call sites: two in
`internal/resolve/guard_mvv_test.go` (`recordingSeam` at :257, bare
`conformingContractSeam` at :462) and a third in another package —
`internal/guard/guard_evaluator_0003_test.go:67`, which drives
`var ev guard.Evaluator`, the ZERO VALUE. So the real implementer is
already gated, but by a nil-mapping seam: under C1's LOUD disposition
that site's `eq`/`in` cases invert from GuardTrue to GuardUnevaluable,
which is exactly what scenario 5 demands. It therefore migrates to a
call over the published fixture in the SAME phase as the signature
widening — it is the third member of that migration list, and the only
one crossing a package boundary.

The constructor is passed WRAPPED, not bare. Go function types are
invariant in the return position, so `guard.NewEvaluator` — typed
`func(map[string]string) guard.Evaluator` by C1 — does not assign to
the `newSeam` parameter above, which returns the INTERFACE
`GuardEvaluator`; measured, `go build` refuses it (spike
`evidence/spikes/c3-newseam-invariance.md`). Each migrating site
therefore supplies a one-line adapter:

func(k map[string]string) resolve.GuardEvaluator {
        return guard.NewEvaluator(k)
}

`NewEvaluator`'s return type is NOT widened to the interface to avoid
the wrapper: C1 fixes the concrete value shape deliberately, and
widening it would hand every non-test construction site an interface
where it holds a struct — trading a three-line closure in three test
call sites for an indirection on the production path.

The widening also breaks a normative pin this RDR must amend
explicitly: `internal/resolve/guard_mvv_test.go:222` holds
`var fn func(*testing.T, resolve.GuardEvaluator) = resolve.TestGuardEvaluatorContract`
inside `TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface`,
whose doc marks the name AND signature as normative BOUNDARY surface
for `0007:REQ-71`. Taking a constructor changes that signature, so the
assignment stops compiling. `0007:REQ-71` is amended here: its
importable-cross-RDR-surface guarantee is preserved — same name, still
declared in the kernel's non-test sources (the `kernelDeclaresFunc`
half at :227 is untouched) — and only the parameter shape changes,
from a constructed seam to a seam constructor. This is a recorded
deviation against an implemented test, the same treatment A2 gives
`0003:C9`/REQ-34, and the `var fn` line is re-typed to the new
signature rather than deleted, so the pin keeps pinning.
```

**C4**

```normative
Surface — of C1; the construction-site obligation keeping the carrier's mapping correct.
PRODUCER OBLIGATION. EVERY non-test evaluator construction site —
present or later added — constructs over the declarations of the SAME
loaded model whose rows it evaluates: one model, one evaluator, no
mixed-model evaluation. The obligation is scoped to the GUARD path —
the decision path lint and the kernel genuinely share. On `main` that
set is exactly two:
`internal/cli/flow_resolve.go::guardSeam`,
`internal/guard/product.go::valueSatisfies` (A3). At the CLI the
constructed evaluator is THREADED DOWN, not re-derived: `guardSeam`
takes the model and is called once per request, and
`internal/cli/flow_next.go::probeRow` receives the constructed
`resolve.GuardEvaluator` rather than widening to take a `*table.Model`.
`probeRow` runs in a loop, so widening it would rebuild the kinds map
per row and give one request several evaluators — the opposite of
this clause's "one model, one evaluator". (A3 left this as an
either/or; the universal above decides it.) The enumeration is
documentation of today's set, not the guarantee — the guarantee is the
universal quantifier, with Validation scenario 4 as its mechanical
backstop ("no zero-value construction of `guard.Evaluator` outside
`NewEvaluator` survives in non-test code"), whose scope narrows with
this clause's: the graphlint match site below is the one admitted
`Evaluator{}` construction, and scenario 4 names it as such rather
than failing on it.

THE MATCH PATH IS EXCLUDED, DELIBERATELY.
`internal/graphlint/reach.go::atomAdmitsValue` is NOT a construction
site under this clause and keeps its `main` shape: it constructs
`guard.Evaluator{}` with no mapping and compares raw bytes. That is
required, not tolerated. Its two callers are both reached only through
MATCH atoms — `reach.go::matchSatisfiable` filters
`a.Block != table.BlockMatch`, and `analysis.go::nodeMeetsAll` reads
`model.Terminal`, whose atoms are built by
`internal/table/load.go::loadContexts` with `BlockMatch` — so typing
them would make lint's match admission diverge from the kernel's, which
byte-compares at `internal/resolve/resolve.go::TagSet.matches`. C1
requires the kernel to keep byte-comparing there; a typed match arm in
lint would be the very lint/runtime divergence this record exists to
close, inverted. Typed comparison therefore applies exactly where lint
and the kernel share a decision path — the guard path, via
`product.go::valueSatisfies` under `Denotation` — which is a
structural fact about the code, stated here and pinned by C3's
conformance suite rather than left to a reader's inference.

That split is the record's one asymmetry and it is named on purpose
(C2 states it on the comparison side, and the user-facing docs carry
it): one declared kind, two comparison rules, divided by whether the
two deciders meet. Alternative (i) — canonicalizing the held ingress
so one rule could serve both — is refused because the wider rule
contradicts RDR 0024's locked REQ-7/REQ-9, and this project does not
re-open locked records; see the Residual note in §Load-Bearing
Decisions for the kata that carries it.

CONSTRUCTION IS NOT ENOUGH — the verdict must also be CONSUMED
three-valued, on the guard path. Constructing over the right model
makes a consumer's verdict correct; it does not make that consumer's
READING of the verdict correct, and the two are separately falsifiable.
`internal/guard/product.go::valueSatisfies` already returns the
three-valued `resolve.GuardResult` on `main`, and its caller switches
on it — so the guard path needs no widening and this clause requires
none. Every consumer of a
seam verdict handles GuardUnevaluable as its OWN case, distinct from
both GuardTrue and GuardFalse; a `!= GuardFalse` / `== GuardTrue`
two-valued collapse at any consumer is a defect this clause forbids.
Lint's stated property — "the same decision the runtime evaluator
makes" (`valueSatisfies` doc) — is a claim about the verdict AND its
reading; on the guard path both already hold on `main`, and this
clause keeps them holding as the seam becomes typed.

Lint's MATCH admission is outside that property by construction, and
`atomAdmitsValue` keeps its `bool` return and its name. Widening it to
`resolve.GuardResult` was considered and is REJECTED: it would type a
comparison the kernel byte-compares, and its failure mode is the worst
shape available — the two callers read the current `bool` with
OPPOSITE polarity (`reach.go::ownedAtomSatisfiable` existentially,
`analysis.go::nodeMeetsAll` exclusionarily), so a sign error at the
negated caller SUPPRESSES a finding rather than producing a wrong one,
which no corpus diff over firing findings detects. Declining the
widening removes that failure mode rather than testing for it.

This clause reaches the CONSTRUCTION and the CONSUMPTION of a verdict,
never its AGGREGATION: how a row's several atoms combine is already
fixed in the kernel by `internal/resolve/guard.go::evaluateAtoms`,
which evaluates every atom with no short-circuit and combines under
strong-Kleene AND, where GuardFalse dominates GuardUnevaluable
(`kleeneAnd`). `probeRow` does not aggregate at all — it builds a
single-row table and calls `resolve.Resolve` once. So MVV step 2's
singular "the `iter` atom" names the one unevaluable atom in that
fixture, not a first-unevaluable-wins rule, and nothing here changes
which atom a multi-atom row reports.

Dependency: `JDR 0004 §JD-3` obliges the loader-side admitted-cell shim
RDR 0030's Phase 2 extracts to construct through `NewEvaluator` with the
loader's own `l.model.Tags`. That shim does NOT exist on `main` — no
`internal/table/normalize.go::renderWrites` is present (verified by
source search at pre-lock critique) — so it is not in today's two-site
set and this clause names no present site for it. It becomes a
construction site when 0030's Phase 2 lands, at which point the
universal above already governs it, whatever the extracted function is
named.
```

No new refusal kind, error code, or envelope field for the SEAM: a typed
`GuardUnevaluable` reaches the user through the existing
`guard_unevaluable` refusal and its payload (`0007:C8`), reason
`uncomparable` ("the key was present and its value was not compared to
a verdict") — the CLI already maps it to `flow-guard-unevaluable`. This
sentence is about the seam ONLY and does not generalize to the record:
C5 adds a user-visible load-time refusal, which is a new message under
an existing code. No envelope FIELD is added anywhere, but a consumer
matching on message text at `malformed_predicate_atom` sees a shape it
did not see before — the record claims code and envelope stability, not
message stability.

**C5**

```normative
Surface — of C1; the ingress canonicalization that keeps the literal
side of a typed comparison from drifting.
CANONICAL INT SPELLING. A value authored against an `int`-declared tag
is admitted only in its canonical decimal spelling: after
`strconv.Atoi` succeeds, `strconv.Itoa(n)` MUST equal the authored
string. Non-canonical spellings (`"00"`, `"01"`, `"+1"`, `"-0"`) are
REFUSED at the ingress that admits them. The rule's scope is the
PREDICATE ingress and only it: guard and match atom literals, the
values a typed comparison reads. It is NOT homed in
`internal/table/load.go::conformKind`'s shared int arm, because that
arm also serves `[emit]`, `[initial]` and `[rule.write]` values, which
`0024:REQ-7` fixes as LEXICAL ("no value is parsed into a typed
representation, canonicalized, or converted") and `0024:REQ-9` settles
in the permissive direction. Two BOUNDARY-marked tests pin that
disposition on `main` and pass today —
`internal/table/emit_grammar_0024_test.go::TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted`
(admits `03`, `+5`, `-0`) and
`::TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive`
(authored `007` rides onto the row verbatim) — so tightening the shared
arm would silently retire a neighbouring RDR's locked call. The
canonicality check is therefore applied at the atom-building site
(`internal/table/normalize.go::(*loader).atom`, which already calls
`conform(decl, operator, members)` at `normalize.go:172`), leaving
`conformKind` itself unchanged. It also covers the bare-float path
there, where `internal/table/load.go::valueMembers` renders a TOML
`eq = -0.0` to the literal `"-0"` via `strconv.FormatFloat` before any
kind check.

The CLI is OUT of scope: `--tag n=07` / `--write n=+1` stay accepted.
They reach `conformKind` through `table.ConformValue`, the same shared
arm, and no measurement covers them (S5 counted guard atoms). Whether
the CLI should also demand canonical spellings is a separate question
this RDR does not decide — see the Charted note in
`evidence/3amigo/`.

The refusal is user-visible and DIAGNOSTIC: it names the offending site
and the canonical rewrite, e.g.
`<site>: "00" is not the canonical spelling of int tag n; write 0`.
It travels as whatever error value `(*loader).atom`'s existing
`conform` call already returns to its caller — the check is folded in
beside that call and returns the same shape, adding no new error type
and no new sentinel. Whether the comparison is factored into a private
helper is implementation latitude, but the refusal must arrive on the
path `conform`'s already does, or it reaches the envelope differently
from the code it reuses.
It reuses the one refusal code its ingress already carries —
`malformed_predicate_atom` (guard and match atoms, via
`internal/table/normalize.go::(*loader).atom`) — and adds no envelope
field and no sixth refusal kind (REQ-7's taxonomy is untouched, in
this RDR's scope and in `0024`'s). The code is broader than its name
reads: a non-canonical spelling is well-formed to `conformKind`, so
"malformed" describes the atom's rejection at this ingress, not the
value's syntax. Reuse is deliberate — a new code for one spelling rule
would fragment the ingress's taxonomy for a distinction the MESSAGE
already draws, and the message is where a user reads the fix ("is not
the canonical spelling … write 0"). The code routes; the diagnostic
explains.

"Canonical" here is exactly the `strconv.Itoa` round-trip stated above
and nothing wider — no case folding, whitespace or unicode
normalization is claimed or performed. Evidence Record: normative
fixture S7 (the `cover07` model and its refusal), plus MVV step 5;
the held side's canonical rendering is A4's verified
`assignment.go` int arm (`evidence/spikes/a4-render-path.md`).

Rationale for admitting no violator: the corpus contains zero
`eq`/`in` atoms over an `int` tag (S6), so the invariant costs nothing
to establish AT THE PREDICATE INGRESS and is what makes C2's parsed
comparison agree with lint's canonical product (A4). S6's zero counts
guard atoms only — it measures nothing about `[emit]`, `[initial]`,
`[rule.write]` or CLI values, which is precisely why the rule is
scoped to the ingress the measurement covers rather than to the shared
`conformKind` arm those populations also reach. A `--fix` affordance
may be layered on later; it does not substitute for the refusal.
```

#### Pre-lock mini-checks

Five structural cues fired at pre-lock; each table is the decision, not
a note about it.

**`authority` — where a comparison's declared kind comes from.** The cue
is three arms deciding value/kind questions over the same models.

| Input / decision | Writer (canonical) | Readers | Call sites | Sibling arms | Canonical? |
| --- | --- | --- | --- | --- | --- |
| Declared kind of a key | `internal/table/load.go` tag declarations (`m.Tags`) | `declaration.go::DeclarationOf`, `NewEvaluator`'s mapping | 2 construction sites (C4) | `Conforms` (presence only, no kind, uncalled) | YES — the loaded model is the sole kind authority |
| Whether an authored value is well-spelled (KIND) | `load.go::conformKind` | every authoring ingress | `normalize.go::(*loader).atom`, `[initial]`, `[rule.write]`, `[emit]`, `flow_input.go::canonicalValue` | none | YES — the only kind authority, shared by all five ingresses, and UNCHANGED by this RDR. It validates KIND only, and only for `int`/`bool`: its switch has no `default`, so an empty-`Kind` tag's members pass unvalidated |
| Whether a PREDICATE literal is canonically spelled | `normalize.go::(*loader).atom` (C5's round-trip) | the guard/match atom builder | `normalize.go::(*loader).atom` only | `conformKind`'s permissive int arm still serves `[emit]`/`[initial]`/`[rule.write]`/CLI, per `0024:REQ-7` | YES for predicate atoms — deliberately NOT total over the other ingresses (S5 measured only this one) |
| Whether an owned READER value is well-formed | nobody today; C1/C2's seam after this change | the seam | `accessor::OwnedSnapshot` → `resolve::assemble` | none | The seam — the gap this RDR closes |
| Guard-atom verdict | `guard.Evaluator` (typed after C2) | `flow resolve`, `flow next`, lint product | `product.go::valueSatisfies`, `flow_resolve.go::guardSeam` | — | YES |
| Match-atom verdict | `resolve.go::TagSet.matches` (byte compare) | kernel resolution | `resolve.go` | `graphlint reach.go::atomAdmitsValue` byte-compares too (C4 excludes it from typing) | YES — both sides byte-compare, so they agree without typing. The held-cell canonicalization gap (`heldValues`/`canonicalValues`, no round-trip) is read the SAME way by both and is therefore not a divergence; it is the Residual carried by a kata (A4) |

**`disposition` — input class → outcome at the seam.** The cue is C2's
per-class verdicts and the load refusals.

| Input class | Verdict / exit | Refusal code | Loud or silent |
| --- | --- | --- | --- |
| Held value parses as declared kind, differs from literal | `GuardFalse` (row prunes) | — | silent by design |
| Held value present, does not parse as declared kind | `GuardUnevaluable` → `flow resolve` exit non-zero | `flow-guard-unevaluable`, reason `uncomparable` | LOUD (the fix) |
| Same, via `flow next` | pair reported `unknown`/`uncomparable` | `0011`'s closed reason vocabulary | LOUD |
| Key absent from the constructed mapping | `GuardUnevaluable` (defensive; unreachable under A1) | `flow-guard-unevaluable` | LOUD |
| `eq`/`in` over a `set` key | `GuardUnevaluable` (defensive) | `flow-guard-unevaluable` | LOUD; load matrix is the primary guard |
| Nil-mapping evaluator (missed C4 site) | `GuardUnevaluable` for every `eq`/`in` | `flow-guard-unevaluable` | LOUD (the zero-value-survives mode) — never a raw-string revert |
| Authored non-canonical int value, PREDICATE atom | refused at the atom builder | `malformed_predicate_atom` | LOUD, with the canonical rewrite named |
| Authored non-canonical int value, `[emit]`/`[initial]`/`[rule.write]` | ADMITTED unchanged | — | silent by design — `0024:REQ-7` fixes these checks as lexical; out of this RDR's scope |
| Same, CLI (`--tag n=07`, `--write n=+1`) | ADMITTED unchanged | — | unchanged — C5 no longer reaches the CLI, so this RDR adds no second visible break |
| Absent key | unchanged three-valued honesty | — | unchanged |

**`fidelity` — C5's round-trip.** The cue is an inverse/normalization
claim.

| Operation | Invariant | Lossy exemptions |
| --- | --- | --- |
| Authored PREDICATE int literal → parsed → re-rendered | `strconv.Itoa(strconv.Atoi(s)) == s`, else refuse at the atom builder | `[emit]`/`[initial]`/`[rule.write]`/CLI int values — exempt by `0024:REQ-7`, which fixes those checks as lexical |
| Declared int domain → held value rendering | `assignment.go::valueAssignments` emits `strconv.Itoa` only (A4) | none |
| TOML bare float → literal | `valueMembers` renders via `FormatFloat` BEFORE the kind check, so `-0.0` becomes `"-0"` and C5 refuses it | none — this is the path `conformKind` alone misses |
| "Canonical" scope | decimal `Itoa` round-trip ONLY | no case folding, whitespace or unicode normalization claimed |

**`oracle` — what each MVV step fails on.** The cue is steps whose
surface reading is "refuses" or "is green".

| MVV step | Fails if X is wrong, because Y | Negative / failing control |
| --- | --- | --- |
| 2 (`many` on `int`) | fails if the seam still compares raw strings — the row would prune silently and the command would exit 0 with a plan | today's behavior IS the control: same invocation plans the fallback |
| 3 (`07` vs `eq 7`) | fails if comparison is not parsed — `"07" != "7"` prunes the row | the `4` leg must still prune (GuardFalse), proving step 3 is not "always true" |
| 4 (suite green) | fails if the fixture is scalar-typed or the seam ignores kinds — the int/bool want-Unevaluable legs cannot pass vacuously | the suite/fixture drift mode; the per-kind-token meta-check |
| 5 (`n eq "00"` refused) | fails if the atom builder keeps the bare `Atoi` reading — the model loads and lint reports clean | measured before-state (S7): exit 2 → exit 0 flip |

**`trace` — the MVV walked stepwise against every clause in force.**

| Step | Clauses in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1 author model + owned reader | C5 (load admits `7` canonical), A1 (`iter` declared) | `iter` int, single-valued; rows `iter eq 7` + fallback — the one literal every step below reads (S8/S9's `mvv-int`) | consistent |
| 2 reader returns `many` | C1 (constructed seam), C2 int arm, C4 (CLI site holds model) | `many` fails `Atoi` → `GuardUnevaluable` | consistent — and `flow next` reports `uncomparable` on the same probe (Failure Modes) |
| 3 reader returns `07` vs `eq 7` | C2 int arm (parsed both sides), C5 (literal `7` is canonical; `07` is unauthorable AS A LITERAL but arrives as a READER value, which C5 does not reach) | parsed `7 == 7` → match | consistent — this is the one leg load cannot reach, which is why it is the MVV's parsed-comparison witness |
| 3b reader returns `4` | C2 int arm | parsed `4 != 7` → `GuardFalse` | consistent — prunes, proving no blanket-true |
| 4 suite green | C3 (per-case keys, both implementers), A5 (no `internal/table` type crosses) | kernel-owned `map[string]string` fixture | consistent |
| 5 `n eq "00"` at load | C5 round-trip, C2 (never reached — load refuses first) | `Itoa(0) == "0" != "00"` → refuse | consistent |
| end-state | C1 LOUD disposition, the visible-break and zero-value-survives modes | seed defect refuses instead of misrouting | no CONTRADICTION row |

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
- **Canonical spelling at load** — A4's literal leg is closed by
  admitting one spelling (C5), not by weakening the comparison.
  Rejected: (a) accepting the lint flip as documented behavior — it
  ratifies lint going silent on a genuine dead end (measured: typed
  lint exit 0 while the runtime is `flow-no-match` at the root); (b)
  validate-then-byte-compare — it reintroduces the silent prune this
  record exists to remove, and splits `int` semantics against
  `lt/gte`, which already parse both sides. A third rejection —
  excluding `atomAdmitsValue` from typed comparison, on the ground
  that it institutionalizes two comparison semantics for one declared
  kind (Meyer, OOSC §11.7 Non-Redundancy pp.354-355; Raymond, AoUP
  SPOT p.124) — was REVERSED at pre-lock; see the next decision, which
  owns that call and its reasoning. Establishing the invariant now is cheapest with zero
  violators (SWE-at-Google pp.161-162, gofmt vs buildifier;
  Refactoring Databases p.186) — but "zero violators" holds only at
  the PREDICATE ingress, which is what scopes C5 there. The shared
  `conformKind` arm has violators and they are sanctioned ones: two
  BOUNDARY-marked `0024` tests admit `03`/`+5`/`-0`/`007` for
  `[emit]`, so the cheap-invariant argument does not reach that arm
  and C5 does not either.

- **The guard/match split, and why typing stops at the guard path** —
  the Non-Redundancy objection above is real and is answered, not
  dismissed. C1 leaves the kernel's MATCH atoms byte-compared
  (`internal/resolve/resolve.go::TagSet.matches`) while C2 makes GUARD
  atoms parsed. That IS two comparison rules for one declared kind, and
  the original defence — that C5 makes the difference unobservable —
  holds for LITERALS only: graphlint's held values arrive via
  `internal/graphlint/reach.go::heldValues` from `[initial]`/
  `[rule.write]` cells that `conformKind` admits non-canonically and
  `canonicalValues` does not round-trip.

  Two repairs were available. (i) Extend canonicalization to the held
  ingress, giving one rule everywhere. (ii) Make lint's match arm
  byte-compare, so the kernel and lint share one reading on the match
  path and typing applies only where they share a decision.
  **(ii) is chosen.** (i) is not available to this record: the wider
  rule is refuted by two BOUNDARY-marked tests enforcing RDR 0024's
  locked REQ-7/REQ-9, and this project does not re-open locked
  records — so (i) would make this record depend on a successor
  landing first while leaving the divergence live in the interim.

  (ii) is also the cheaper repair in both directions. For the USER it
  costs nothing: match atoms compare exactly as they do on `main`, so
  there is no new refusal, no model edit, no migration — and guards
  get the fix this record exists to deliver. For the implementation it
  is a narrowing rather than new surface, and it DELETES work the lens
  row found defective: the three-valued widening of graphlint's atom
  check (former A7), whose failure mode was a SUPPRESSED finding
  invisible to a diff over firing findings, along with C4's
  consumption obligation on that path and the `atomVerdictForValue`
  signature.

  The boundary (ii) draws is principled rather than arbitrary: typing
  applies exactly where lint and the kernel share a decision path.
  That is a structural fact about the code, not a convention — both
  callers of `reach.go::atomAdmitsValue` are reached only through
  `BlockMatch` atoms (`matchSatisfiable` filters on it;
  `analysis.go::nodeMeetsAll` walks `model.Terminal`, built with
  `BlockMatch` by `load.go::loadContexts`), while the guard path runs
  through `product.go::valueSatisfies` under `Denotation`. C2 and C4
  state the split, and the user-facing docs carry it, so it is
  discoverable rather than surprising.

  **Residual, stated not hidden.** A non-canonical int in an
  `[initial]`/`[rule.write]` cell still reaches `heldValues` →
  `canonicalValues` with no integer round-trip. Under (ii) this is no
  longer a lint/runtime DIVERGENCE — both sides byte-compare those
  cells and read them identically — but it remains a
  canonicalization gap at that ingress. It is out of scope here
  because closing it is exactly repair (i), which RDR 0024's locked
  REQ-7/REQ-9 forbid this record from making. Kata
  `intrastate#ch99` carries it with
  the pre-lock critique evidence and that conflict recorded as the
  reason.
- **C5 is not Alternative 3** — ALT3 was rejected as a *runtime*
  producer obligation over caller-supplied values, which leaves every
  non-CLI `resolve.Resolve` caller unprotected and pre-empts JDR 0001
  §JD-18. C5 constrains *authored predicate literals* at load, a venue
  that already conforms kinds (the atom builder's existing `conform`
  call) and whose population is
  disjoint from the runtime's (`JDR 0004 §JD-2`). It tightens an
  existing load-time check rather than adding a runtime validation
  layer, so the seam still carries its own honesty. The distinction
  survives the narrowing: ALT3's defect was the VENUE (runtime,
  caller-supplied), not the reach, so scoping C5 to one authoring
  ingress leaves it on the correct side of that line.

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
| Declared-kind lookup | `internal/guard/declaration.go::DeclarationOf` | Takes `*table.Model`; seam constructor needs a plain map (A5) | Reuse, behind the exported `guard.DeclaredKinds(m)` (C1) — exported because the two C4 sites span two packages | No parallel kind model invented |
| Value seam | `internal/guard/grammar.go::Evaluator` | Zero-value `struct{}`; raw-string `eq`/`in` arms | Extend (constructor + typed arms, C1/C2) | The RDR's core change |
| View-freedom proof | `internal/guard/guard_evaluator_0003_test.go::TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView` | Zero-field check is stricter than REQ-34's text (A2) | Amend (field check narrows to "no view-typed/runtime state") | Recorded deviation against 0003's implemented proxy |
| Cross-implementer contract | `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` | No kind carrier, so no kind-aware cases exist | Extend (C3) | Suite grows fixture + cases |
| Malformed-set handling | `internal/guard/grammar.go::parseHeldSet` | Already answers unevaluable without declarations (the `072c7a0` inline fix) | Reuse unchanged | Proves the disposition; doesn't reach `eq`/`in` |
| Load-time kind conformance | `internal/table/load.go::conformKind` | Int arm is `strconv.Atoi` only, so it admits `"00"`/`"+1"`/`"-0"`; SHARED by every authoring site (guard and match atoms via block-agnostic `normalize.go::(*loader).atom`, `[initial]`, `[rule.write]`, `[emit]`) and by the CLI via `ConformValue` | Reuse UNCHANGED — the shared arm's permissiveness is `0024:REQ-7`/`REQ-9`'s settled call, pinned by two BOUNDARY tests | C5 lands beside `(*loader).atom` instead, so no non-predicate ingress changes |
| Predicate-atom canonicality | `internal/table/normalize.go::(*loader).atom` | Already calls `conform(decl, operator, members)` (`normalize.go:172`); does not check canonical spelling | Extend (add C5's `Itoa`-round-trip beside the existing conform) | Keeps the refusal under the existing `malformed_predicate_atom` category — no new code |
| CLI value conformance | `internal/cli/flow_input.go::canonicalValue` | Already conforms every declared `--tag`/`--write` before resolution (`flow-tag-invalid` / `flow-write-invalid`) | Reuse unchanged | Confirms the CLI is not the defect's ingress (S8) |

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
extends it. Recorded symmetrically by 0030's propose, 2026-09-19. Its
re-fire opened `JDR 0004` (settled 2026-09-20, binds 0012 and 0030),
whose three entries land here as: C1's closing sentence scoped to the
resolution path (§JD-1), the ingress disjointness cited in §Problem
Statement and A4 (§JD-2), and C4's forward dependency on the
admitted-cell shim 0030 Phase 2 extracts (§JD-3). The
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
kind-lookup path exists. The nearest non-duplicate is
`internal/guard/declaration.go::Conforms(m *table.Model, v View)`, the
view-conformance premise 0003 states its lint claims over: it holds
both the model and the view and returns an error, but checks only
required-key presence and single-valuedness — never kind — and has no
non-test caller on `main`. So it does not decide what C1/C2 decide and
is not a parallel implementation; it is the natural venue for
`JDR 0001 §JD-18`'s open view-conformance question, named here so a
later answer lands on it rather than minting a third path.

### Key Discoveries

- **Documented** — `0007:A17`'s Evidence already names construction as
  the delivery path (quoted in §Approach): A17 froze the signature,
  not the constructor, so the chosen carrier is the one the verifying
  record itself anticipated. Its headline's "nor tag declarations at
  evaluation time" is narrowed to "no declaration it was not built
  over" — the deviation §Approach records.
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
  `Evaluator{}` construction sites (A3's measurement); C4 migrates the
  two on the guard path and leaves the graphlint match site
  byte-comparing. No other consumer exists to inherit the constructor
  change.
- **Verified** — typed numeric equality agrees with lint's rendered
  product ONLY once the PREDICATE literal side is canonicalized at
  load: the
  held side is already canonical (`strconv.Itoa`), but `conformKind`
  admits `"00"`/`"+1"`/`"-0"`, and the measured effect is a lint flip
  from blocking to clean (A4, S6). C5 closes it at the atom builder,
  leaving `conformKind` shared and permissive for the non-predicate
  ingresses (`0024:REQ-7`); the fallback
  (validate-then-byte-compare) stays recorded in Load-Bearing
  Decisions.
- **Verified** — the seam decides MATCH atoms too, on lint's side only.
  The kernel byte-compares MATCH atoms
  (`internal/resolve/resolve.go::TagSet.matches`, plain `!=`) and C1
  leaves that unchanged, but graphlint's reachability runs each owned
  MATCH atom through the same seam
  (`internal/graphlint/reach.go::matchSatisfiable` keeps
  `a.Block == table.BlockMatch` and consults no guard atom;
  `atomAdmitsValue` calls `guard.Evaluator{}`). So at that third
  construction site C4's "the same decision the runtime evaluator
  makes" is measured against `TagSet.matches`, not the guard seam, and
  C2's parsed comparison makes the two differ in principle. C5 closes
  the LITERAL leg of that difference — match literals are canonicalized
  at load, `internal/table/normalize.go::(*loader).atom` being
  block-agnostic — but NOT the held leg, and the difference is therefore
  observable. `internal/graphlint/reach.go::heldValues` sources
  `[rule.write]`/`[initial]` cells and passes them through
  `reach.go::canonicalValues`, which sorts and dedupes and applies no
  integer round-trip; those cells are conformed only by
  `internal/table/load.go::conformKind`, whose int arm is a bare
  `strconv.Atoi` admitting `"007"`. So a model CAN author a held value
  the two read differently — lint parses it at the guard seam, the
  kernel byte-compares it at `TagSet.matches` — which is the `matchlit`
  shape the Q1 ruling rejected a design for reopening. Canonicalization
  as narrowed does not close the divergence; A4 carries the open
  question of whether the held leg must be closed too (a second C5
  ingress) or the seam's match arm must byte-compare to match the
  kernel.
- **Verified** — the owned ingress is the sole door that conforms
  neither KIND nor SPELLING
  (`internal/accessor/model.go::ReadResult.OwnedSnapshot` →
  `internal/resolve/resolve.go::assemble`, no conform call on the
  path), which is why the MVV lives there and not on `--tag`. It is not
  the only door admitting non-canonical SPELLING: the model-literal and
  CLI doors both run `conformKind`, whose permissive int arm admits
  `"07"`, so they conform kind while leaving spelling open. C5 closes
  spelling at the predicate-literal door only; `--tag iter=07` stays
  admitted (§Failure Modes).

## Trade-offs

### Consequences

- Positive: the fixed semantics become computable and suite-enforced —
  a malformed caller value on a typed dimension refuses
  `flow-guard-unevaluable` instead of silently pruning or misrouting;
  every `GuardEvaluator` consumer inherits the fix, not just the CLI.
- Positive: kernel, atom shape, refusal taxonomy, and envelope are
  byte-untouched. The PRODUCTION surface is small and counted: three
  construction sites across two importing packages (`internal/cli`,
  `internal/graphlint`) plus `internal/guard` itself. The TEST surface
  is not small — eighteen zero-value construction sites in
  `internal/guard`'s own suite migrate with the constructor (Phase 1's
  count: eleven `Evaluator{}` literals plus seven `var ev` forms), plus
  `internal/cli`'s four `guardSeam()` fixtures, `internal/graphlint`'s,
  C3's conformance cases and C5's loader fixtures. "One package plus
  three call sites" describes the production diff only; the landed
  change is several times that, and Phase 1 carries most of it.
- Negative: zero-value `Evaluator{}` construction is retired with no
  compatibility shim (project convention: no back-compat); all sites
  and fixtures migrate in one change.
- Negative: 0003's statelessness proxy (the zero-field reflection
  check) narrows — a recorded deviation against a closed RDR's
  implemented test, justified by A2.
- Negative: verdicts flip for previously "deciding" comparisons, and
  the two legs are bounded differently. Malformed values
  (False → Unevaluable) are the fix, and loud. Non-canonical numerals
  under C2's parsed comparison (`"07" eq 7`: False → True) are a
  SILENT flip, and are bounded only on the authored side: A4 covers
  authored predicate literals (and is itself Pending after C5's
  narrowing). The held side — persisted owned values, the population
  this RDR's own ingress analysis calls the wide door — is UNMEASURED
  and tracked by A6. So the blast radius of the silent leg is not yet
  known, which is a rollout input, not a design defect.

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

- **Visible break**: a reader returning an unparseable owned value on an
  `int` dimension now exits `flow-guard-unevaluable`, payload naming the
  row and atom with reason `uncomparable` — previously a silent prune or
  a wrong sibling selection. Diagnosis reads the refusal payload;
  recovery is fixing the reader's value or the persisted artifact. A
  second, load-time break: a model authoring a non-canonical `int`
  spelling (`n eq "00"`) now refuses at load under C5 with the canonical
  rewrite named. Zero committed models are affected (S6), so the break
  is prospective. There is NO third break: C5 is scoped to the
  predicate ingress, so the CLI and the `[emit]`/`[initial]`/
  `[rule.write]` sites keep accepting what they accept today.

  The load-time break has a SECOND-ORDER form worth naming, because
  following the diagnostic is what triggers it. A model authored
  consistently non-canonical — `[match] n = "007"` with callers passing
  `--tag n=007` — works today: the kernel byte-compares match atoms
  (`internal/resolve/resolve.go::TagSet.matches`), so both sides read
  `"007"`. C5 refuses that model at load and names the rewrite; once the
  author rewrites the literal to `7`, the SAME callers' `--tag n=007`
  stops matching, because the CLI ingress still admits `007` and the
  kernel still byte-compares. The refusal is loud, but the breakage it
  invites is not — it lands at the caller, after the fix. This is the
  held-leg divergence (Key Discoveries) seen from the author's side, and
  it is the strongest argument that A4's open question must be settled
  before implementation rather than after.

  That scoping carries a known, accepted residual, and it is the same
  one that argued for the wider rule: a CLI `--tag iter=07` against
  `iter eq 7` is admitted at input and then, under C2's parsed
  comparison, matches where today's raw-string seam prunes — a verdict
  moving GuardFalse → GuardTrue silently (measured,
  `evidence/spikes/a4-typed-compare.md`). The wider C5 would have
  refused it loudly. It is accepted here because closing it at
  `conformKind` would refuse literals two BOUNDARY-marked `0024` tests
  admit, retiring a neighbouring RDR's locked REQ-7/REQ-9 disposition
  from inside this one — a cross-RDR contract change this RDR has no
  mandate to make and no measurement to justify (S5 counted guard
  atoms only). The honest shape is: this RDR fixes the owned ingress
  it set out to fix, and leaves the CLI canonicality question open for
  a successor that can weigh it against `0024` properly. Charted in
  `evidence/3amigo/`. The break reaches a second verb:
  `flow next` derives
  its `unknown` pairs from the same probe seam
  (`internal/cli/flow_next.go::probeRow`), where the `guard_unevaluable`
  payload is the only source of `uncomparable`, so a malformed owned int
  is now reported `uncomparable` there too, and a non-canonical held
  value (`"07"` against `n eq 7`) moves a row between excluded and
  candidate. Same honesty, second surface — and one governed by `0011`'s
  closed reason vocabulary, which the new refusals reuse rather than
  extend.
- **Undeclared-key guard flip (A1's If-wrong)**: A1 holds that every
  guard-referenced key is declared, so the constructed mapping always
  answers. Were a guard to reference an undeclared key, its kind lookup
  misses and a comparison that decided before this change refuses
  `flow-guard-unevaluable` instead — loud, but a behavior change A1 did
  not predict. Detection is the refusal payload naming a key absent from
  the model's declarations; A1 is Verified against
  `internal/table/normalize.go::(*loader).atomsFromBlock`, which refuses an
  undeclared key at load under `CatUnknownTag` before any guard runs, so
  the scenario needs that load check to regress first.
- **Silent failure guarded against**: the misroute itself — a
  present-but-malformed value deciding `GuardFalse` and pruning a row
  the caller believes was compared.
- **Residual silent mode**: a construction site violating C4 (wrong
  model's kinds) can still type a comparison wrongly without refusing —
  guarded by the site audit (A3) and the conformance suite's
  fixture-bound cases, not by the kernel.
- **Suite/fixture drift**: an implementer constructs over its own
  kinds instead of the published fixture and the kind-aware cases pass
  vacuously (everything scalar → string compare). The suite's
  discriminating cases (int/bool want-Unevaluable legs) fail against a
  scalar-typed fixture, which is the detection.
- **Zero-value evaluator survives migration**: a missed construction
  site keeps `Evaluator{}` (nil mapping) — every `eq`/`in` atom it sees
  answers GuardUnevaluable, so the miss surfaces as loud
  `guard_unevaluable` refusals on previously-deciding guards (C1's
  fixed disposition), never as a silent raw-string revert. Diagnosis:
  the refusal payload's `uncomparable` reasons on well-formed values;
  recovery is constructing per C4. The loudness binds `eq`/`in` ONLY:
  `lt/lte/gt/gte` and `contains` are operator-inferred and read no
  mapping, so a missed site whose guards use only those operators
  migrates with no signal at all. Detection there falls to S4's typed
  construction check and the A3 site audit, not to a refusal — narrower
  than "the miss surfaces loudly" reads.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A6). A6 is IN the gate, not
      an optional tail: it is the held-side census, and the population it
      measures is the one silent verdict change this RDR would otherwise
      introduce (non-canonical-but-parseable `"07"` against `n eq 7`,
      today GuardFalse, after C2 GuardTrue). An RDR whose premise is that
      silent verdict changes are the harm cannot ship with that set
      unmeasured — so A6's outcome gates rollout shape (release note,
      scan, or `--fix` affordance), not merely the record's tidiness.

### Minimum Viable Validation

1. Author a model declaring `iter` as `int`, single-valued, with a row
   guarded `iter eq 7` and an unguarded fallback row on the same
   outcome, plus an accessor reading `iter` as an OWNED value — the
   same `mvv-int` model and literal fixture S8 and S9 use, so every
   step below reads one guard. The
   owned ingress is the venue because it is the only one that reaches
   the seam unconformed (§Problem Statement); the CLI `--tag` path
   cannot host this step — `--tag iter=many` refuses `flow-tag-invalid`
   upstream at `canonicalValue`, measured both before and after the
   change (fixture S8). The reader is driven by the persisted artifact,
   not by a test double: steps 2 and 3 set `iter`'s owned value in the
   artifact the accessor reads (`many`, then `07`, then `4`) and re-run
   the command, so the MVV exercises the real
   `OwnedSnapshot` → `resolve::assemble` path end to end. A seam-level
   unit test over `NewEvaluator` would satisfy the assertions while
   testing none of the ingress, which is the whole claim.
2. With the reader returning `many` for `iter`: `flow resolve
   --outcome <o>` → refuses `flow-guard-unevaluable`; the payload names
   the guarded row and the `iter` atom with reason `uncomparable`.
   Today the same invocation silently prunes the guarded row and plans
   the fallback. The envelope carries the existing
   `code`/`message`/`schema_version` fields and NO new top-level key —
   the row, the atom and the reason are named in `message`, the same
   place S8's control envelope names its offending value. Naming them
   as new sibling keys would add envelope surface this RDR elsewhere
   declares it adds none of (C5), and the reason vocabulary is
   `0011`'s existing closed set, reused rather than minted.
3. With the reader returning `07` against the guard `iter eq 7`: the
   guarded row MATCHES under parsed comparison where it was pruned
   under raw-string comparison — the parsed-comparison leg, on the one
   path load cannot reach. With the reader returning `4`: the guarded
   row prunes (GuardFalse) and the fallback plans, which is the control
   proving step 3 is not blanket-true.

   This leg sits on the GUARD path and is therefore unaffected by C4's
   exclusion of the match site: `iter eq 7` is a row GUARD, so the
   owned value reaches `evaluateAtom` and the typed seam, never
   `reach.go::atomAdmitsValue`. The narrowing removes no MVV step —
   every step above is a guard-path or load-path assertion.
4. `TestGuardEvaluatorContract`, extended per C3 and driven against
   `guard.NewEvaluator` constructed over the published fixture kinds,
   is green — including the int/bool want-Unevaluable legs that fail
   against today's raw-string arms.
5. `lint --model` on a model authoring `n eq "00"` over an `int` tag
   refuses with C5's canonical-spelling diagnostic naming the
   site and the rewrite — where today it loads and reports clean
   (fixture S7). The refusal fires at the atom builder, so it reaches
   predicate literals only; an `[emit]` or CLI value spelled `"00"` is
   still admitted (C5).

End-state: step 2's invocation is the seed defect
(kata `intrastate#cq5p`) refusing loudly instead of misrouting.

Acceptance for the reporting consumer: any caller reaching the seam
through an OWNED value — a `resolve.Resolve` library caller, or the CLI
where the value comes from the artifact rather than `--tag` — sees a
`flow-guard-unevaluable` refusal where a plan previously appeared. That
is the intended outcome: a misroute becomes a refusal it can report.
The acceptance does NOT rest on the rdr navigator specifically
exercising it: a consumer that reaches the CLI only through `--tag`,
with no `int`-declared dimension, is conformed upstream at
`canonicalValue` and cannot reach this defect or this fix at all. The
MVV's venue is the owned door for exactly that reason, and the
beneficiary of the change is the class of owned-value callers, not one
named program. It is a NEW terminal state for every such caller.
Callers that treated "no plan" as impossible now have a case to
handle; the refusal reuses `0011`'s existing closed reason vocabulary
(`uncomparable`) rather than minting a new one, so the handling is
additive, not a re-shape.

The delta is not only misroute → refusal. A held value that is
non-canonical but parseable (`"07"` against `iter eq 7`) moves the
other way — it prunes today and MATCHES after C2, silently (scenario
9 specifies this as intended). On an ingress this record exists to
de-silence, that leg deserves naming: it is correct under the declared
kind, it is the point of parsed comparison, and its population is
unmeasured. A6 carries the scan, and whether that scan is a pre-lock
or pre-rollout obligation is Stage 6's call.

### Phase 1: Seam constructor and typed arms

Intent: `guard.NewEvaluator` over a key→kind mapping (with its
`DeclaredKinds` producer, C1); C2's typed
`eq`/`in` arms; zero-value construction retired; REQ-34 test's field
check narrowed per A2 (the recorded deviation).

Phase 1 also carries `internal/guard`'s own test migration, and the
set is BOTH zero-value spellings, not just the composite literal:
eleven `Evaluator{}` literals (`guard_evaluator_0003_test.go` 7,
`guard_narrowing_0003_test.go` 3, `guard_fixtures_0003_test.go` 1)
AND seven `var ev guard.Evaluator` declarations in four files —
`guard_evaluator_0003_test.go:26`, `guard_declaration_0003_test.go:267`,
`guard_grammar_0003_test.go` (:31, :177, :426) and
`guard_scope_0003_test.go` (:297, :353). Eighteen sites; the two
`var`-form files are ones the composite-literal grep does not reach,
and `guard_grammar_0003_test.go` is the densest `eq`/`in` user of all.
It has to be this phase: once the typed arms land, every one of those
evaluates `eq` against a nil mapping and each test asserting
GuardTrue/GuardFalse on an `eq` goes red. Deferring them to Phase 3
would leave `main` red between the two, which the project's
run-`make check`-before-done convention does not allow. Phase 3's
"test fixtures migrate" is therefore the OTHER packages' fixtures
(`internal/cli`'s four `guardSeam()` call sites in
`escape_shape_0009_test.go` :83/:624 and
`escape_shape_adv_0009_test.go` :40/:159, plus `internal/graphlint`),
not these.

### Phase 2: Conformance suite extension

Intent: publish the fixture kinds map in `internal/resolve`, add the
kind-discriminating cases (C3), keeping the suite free of
`internal/table` types (A5).

### Phase 3: Producer alignment

Intent: the two C4 sites construct over their own loaded model's
declarations via `DeclaredKinds`; at the CLI the evaluator is threaded
into `probeRow` rather than widening it to take a model (C4);
`internal/cli` test fixtures migrate to the constructor.
`internal/graphlint::atomAdmitsValue` is deliberately NOT migrated — it
is the match path, which keeps byte-comparing (C4).

### Phase 4: Canonical int spelling at the predicate ingress

Intent: apply C5's `strconv.Itoa(n) == authored` rule at
`internal/table/normalize.go::(*loader).atom`, beside its existing
`conform(decl, operator, members)` call, with the diagnostic naming
the site and the canonical rewrite under the ingress's existing
`malformed_predicate_atom` category. `load.go::conformKind` is left
UNCHANGED: it is shared with `[emit]`/`[initial]`/`[rule.write]`,
whose permissive spelling `0024:REQ-7`/`REQ-9` fix and two
BOUNDARY-marked tests pin (see C5). A green `make check` after this
phase is the expected outcome — no existing corpus or testdata model
should start failing; one that does means the check landed on the
shared arm rather than the predicate site. Lands with S7's regression
and after Phase 1, so the literal side is closed before typed
comparison ships — the ordering A4 depends on.

## Validation

### Testing Strategy

The matrix the verified assumptions imply. Each row names the code path
the verification reviewed, so a regression lands against the same
evidence that admitted the design.

**Known coverage gaps** (accepted, not oversights). Two obligations
carry no automated oracle and are guarded by review instead.
First, C4's "SAME loaded model" half: scenario 4 asserts only that no
zero-value construction survives, so a site constructing over the
WRONG model's kinds still compiles and types comparisons silently
(F4). There is no observable — no panic, no refusal code — because
both models are well-formed; it is guarded by the A3 site audit and by
there being exactly two migrated sites, each with one model in scope.
Second, the `enum`-vs-`scalar` distinction: they share one arm by
construction (C2), so no input discriminates them, and the suite's
per-kind-token cases separate each from `int`/`bool` only. What the
suite DOES pin is that an unknown kind token answers
GuardUnevaluable, which is the regression that would matter.

1. **Scenario**: Shared conformance suite, kind-discriminating cases
   (C3), driven against `guard.NewEvaluator` constructed over the
   published `map[string]string` fixture.
   **Expected**: green, including the `int`/`bool` want-Unevaluable legs
   that fail against today's raw-string arms. Backed by A5's verified
   import direction (`internal/table` → `internal/resolve`, one-way), so
   the kernel-owned fixture crosses no `table` type.
2. **Scenario**: Per-kind dispatch at the seam — `int`, `bool`, `enum`,
   `scalar`, `set`, and key-absent — over `eq` and `in`.
   **Expected**: the C2 matrix verbatim. The `set` and key-absent arms
   are defensive; A1's verified load-time refusal
   (`internal/table/normalize.go::(*loader).atomsFromBlock`, block-agnostic) is
   what makes the key-absent arm unreachable in production.
3. **Scenario**: The zero-field reflection assertion in
   `guard_evaluator_0003_test.go` narrowed to "no view-typed or
   runtime-valued state".
   **Expected**: passes with the declaration mapping present, fails if a
   `resolve.TagSet` or runtime value is added. A2 verified this is a
   test-body proxy, not a 0003 fence, so the narrowing is a recorded
   deviation against an implemented test rather than a contract change.
4. **Scenario**: Construction-site migration (C4) — the two guard-path
   sites A3 verified on `main`, each constructing over the same model
   whose rows it evaluates.
   **Expected**: no zero-value construction of `guard.Evaluator` outside
   `NewEvaluator` survives in non-test code, with ONE allow-listed
   exception: `internal/graphlint/reach.go::atomAdmitsValue`, which C4
   deliberately leaves on the match path byte-comparing. The assertion
   matches the
   composite-literal form (`Evaluator{}`) *and* the `var` / `new` /
   embedded-field zero values Go equally admits, since each yields the
   same nil-mapping evaluator. It is C4's mechanical backstop, so a site
   added later (the shim `JDR 0004 §JD-3` obliges, once 0030 Phase 2
   lands) is caught without re-enumerating. The allow-list is ONE named
   function, asserted by name rather than by pattern, so a second
   untyped site cannot slip in under it and moving or renaming that
   function fails the check — which is the point: the exception is as
   auditable as the rule. Implement it as a typed
   check, not a text grep: a `go/ast` or `golang.org/x/tools/go/packages`
   pass over non-test files in the module, failing on any composite
   literal, `var` declaration, `new`, or struct field of type
   `guard.Evaluator` outside `NewEvaluator`'s own body. A regexp over
   source cannot separate those forms from a constructor call and would
   pass vacuously; the scenario is only a backstop if it runs, so it is
   wired into `make check` rather than left as prose. It lives as a Go
   test under `internal/guard` — the package whose invariant it
   guards, so it moves with the type it pins and needs no separate
   `make` target beyond the existing `test` leg. A `cmd/`-hosted tool
   would be a second build artifact and a second thing to remember to
   wire.
   **Why the backstop matters more than C1's LOUD disposition reads.**
   C1 fixes a nil-mapping evaluator's disposition as LOUD, but loudness
   is a property of the CONSUMER, not the verdict: at
   `internal/graphlint/reach.go::atomAdmitsValue` a
   GuardUnevaluable is folded into "admissible" by
   `verdict != resolve.GuardFalse`. On the match path that folding is
   harmless — the site is allow-listed and never types a comparison —
   but it is exactly why a GUARD-path site accidentally left
   zero-valued would not surface: the LOUD fallback would be swallowed
   somewhere that looks like lint staying green. This scenario is the
   mechanical backstop for that, which is why it must run rather than
   sit as prose.
5. **Scenario**: Nil-mapping disposition (C1) — the guarantee S4's
   typed construction check explicitly defers to, asserted directly.
   Construct a bare
   `guard.Evaluator{}` (and the `var` / `new` zero values) and evaluate
   an `eq` atom whose held value equals its literal, plus an `in` atom
   whose held value is a member.
   **Expected**: `GuardUnevaluable` for BOTH — never `GuardTrue`. The
   held-equals-literal input is the discriminating one: a raw-string
   fallback answers `GuardTrue` and only this assertion catches it.
   This is the negative control for the zero-value-survives failure
   mode, and after the migration it is the one case the conformance
   suite cannot supply, since the suite will drive only
   `NewEvaluator`-constructed seams. Today the opposite holds and that
   is the conflict this scenario resolves:
   `internal/guard/guard_evaluator_0003_test.go:67` drives the suite
   over `var ev guard.Evaluator`, so the suite's `eq`/`in` cases
   currently assert GuardTrue on the exact zero value this scenario
   requires to answer GuardUnevaluable. That site is RE-POINTED to
   `guard.NewEvaluator` over the published fixture (C3), which is what
   makes the two consistent; a run that leaves it pointed at the zero
   value fails one or the other by construction.
   **Scope note**: the LOUD disposition binds `eq`/`in` only.
   `lt/lte/gt/gte` and `contains` are operator-inferred and consult no
   mapping, so a nil-mapping evaluator still answers them normally —
   a missed C4 site whose guards use only those operators migrates
   silently. That residual is F4's, not F6's; F6 is narrowed to match.
6. **Scenario**: Lint/runtime agreement over the committed corpus —
   before/after verdict diff.
   **Expected**: no committed model's verdict flips. Measured: 123
   models, no diff (`evidence/spikes/a4-lint-diff.md`). This corpus
   contains zero `eq`/`in` atoms over an `int` tag, so the result is a
   coverage statement; the authored-literal case below is what the suite
   must actually pin.
7. **Scenario**: Authored non-canonical `int` literal — model `cover07`
   declares `[tags.n]` int min 0 max 1, row `r-zero` guarded
   `n eq "00"` emitting `zero`, row `r-one` guarded `n eq 1` emitting
   `one`; invoked `intrastate lint --model cover07.toml --as json`.
   **Expected** (normative fixture, C5): exit 2 with the
   canonical-spelling load refusal naming the site and the rewrite —
   `"00" is not the canonical spelling of int tag n; write 0` — under
   the ingress's existing `malformed_predicate_atom` code. This is the
   regression guarding the blocking→clean flip measured on `main`
   (`evidence/spikes/a4-lint-diff.md`): today the model LOADS and lint
   reports blocking `graph-coverage-gap` (exit 2), and under typed
   comparison alone it would have gone clean (exit 0). C5 makes the
   model unauthorable instead, so neither verdict is reachable. Also
   covers the bare-float spelling `eq = -0.0`, which
   `valueMembers` renders to the literal `"-0"` (A4).
8. **Scenario**: CLI `--tag` conformance fires upstream of the seam —
   same `mvv-int` model (`iter` int min 0 max 9, row `guarded` guarded
   `iter eq 7`, unguarded row `fallback`); invoked `intrastate flow
   resolve --model mvv-int.toml --outcome step --tag iter=many
   --plan-only --as json`.
   **Expected** (normative fixture), unchanged before and after: exit 2,
   `{"code":"flow-tag-invalid","message":"the value for `iter` does not
   conform to its declaration: \"many\" is not an int",
   "schema_version":"0.1","param":"iter"}`
   (`evidence/spikes/a4-typed-compare.md`). This pins where `--tag`
   conformance fires and is the standing reason the CLI tag path cannot
   host the MVV's unevaluable step (§MVV step 1).
9. **Scenario**: Owned-reader parsed comparison — a reader returning
   `07` for `iter` against the guard `iter eq 7`, on the `mvv-int`
   model.
   **Expected**: the guarded row MATCHES (parsed 7 == 7) where today's
   raw-string seam prunes it. This is the parsed-comparison leg C2
   asserts, carried here and by C3's conformance-suite cases because the
   CLI route to it closes under C5 — an owned value is the one path load
   does not conform (§Problem Statement).

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

Peer records and joint decisions
- RDR 0003 — guard predicate exhaustiveness: `C8` (a literal
  unparseable for its declared kind is rejected before resolution),
  `C9`/REQ-34 (the seam never reads the tag view), the operator/kind
  matrix, and the frozen eight-operator set.
- RDR 0007 — guard predicate totality: `C1`/REQ-10 (the seam
  signature), `A7` (every guard-referenced key is declared), `A17`
  (the evaluator reads no declaration at evaluation time), REQ-71.
- RDR 0011 — the refusal reason vocabulary reused by the
  `flow-guard-unevaluable` payload (`uncomparable`).
- RDR 0024 — `REQ-7`/`REQ-9`, the locked permissive int spelling at
  the `[emit]`/`[initial]`/`[rule.write]` ingress; the reason the
  held-ingress canonicalization residual is out of scope here.
- RDR 0030 — Phase 2 extracts the loader-side admitted-cell shim that
  becomes a construction site under C4.
- JDR 0001 — `§D1` (atom shape), `§D12` (the match block), `§D13`
  (set member encoding), `§JD-18`.
- JDR 0004 — `§JD-1` (the resolution path reads no declarations),
  `§JD-2` (the widening's reach), `§JD-3` (the shim's construction
  obligation).

Source paths reviewed (this repo, at `main`)
- `internal/resolve/resolve.go` — `TagSet.matches` (match-atom byte
  comparison), `GuardResult` constants.
- `internal/resolve/guard.go` — `GuardAtom`, `Block` vocabulary,
  `evaluateAtoms`, `kleeneAnd`/`kleeneNot`.
- `internal/resolve/guardcontract.go` — the shared conformance suite.
- `internal/guard/product.go` — `valueSatisfies`, `Denotation`,
  `selectionOf`, `CanRefuse`.
- `internal/guard/assignment.go` — `valueAssignments`, `renderSet`.
- `internal/guard/declaration.go` — `DeclarationOf`, `Conforms`.
- `internal/graphlint/reach.go` — `matchSatisfiable`,
  `ownedAtomSatisfiable`, `atomAdmitsValue`, `heldValues`,
  `canonicalValues`, `renderSetLiteral`.
- `internal/graphlint/analysis.go` — `nodeMeetsAll`,
  `satisfiesSomeTerminal`.
- `internal/table/load.go` — `conformKind`, `valueMembers`,
  `loadContexts`, `loadTerminal`, `contextAtoms`, `loadInitial`.
- `internal/table/normalize.go` — `(*loader).atom`, `renderWrites`.
- `internal/cli/flow_resolve.go` — `guardSeam`, `escapeClassOf`.
- `internal/cli/flow_next.go` — `probeRow`, the reason vocabulary.
- `internal/cli/flow_input.go` — `canonicalValue`, `codeInvalidFor`.
- `internal/accessor/model.go` — `OwnedSnapshot`.
- `models/rdr.toml`, `models/examples/release-grammar.toml` — the
  owned/observed `int`/`bool` declarations A6's census scanned.

Spike artifacts
- `evidence/spikes/a4-render-path.md` — the lint product's `int` render
  path.
- `evidence/spikes/a4-lint-diff.md` — the authored `n eq "00"` lint
  flip and the 123-model corpus diff.
- `evidence/spikes/a4-typed-compare.md` — `--tag iter=07` under typed
  comparison.
- `evidence/spikes/c3-newseam-invariance.md` — Go function-type
  invariance on the `newSeam` parameter.

Prior art
- Meyer, *Object-Oriented Software Construction*, §11.7
  Non-Redundancy, pp. 354-355.
- Raymond, *The Art of Unix Programming*, SPOT, p. 124.
- Winters/Manshreck/Wright, *Software Engineering at Google*,
  pp. 161-162 (gofmt vs buildifier).
- Ambler/Sadalage, *Refactoring Databases*, p. 186.
- `evidence/research/prior-art.md` — the literature pass behind the
  canonical-spelling decision.

Related issues
- kata `intrastate#cq5p` (1528) — the seed defect.
- kata `intrastate#ch99` — the held-ingress canonicalization
  residual (`[initial]`/`[rule.write]` → `heldValues` →
  `canonicalValues`, no int round-trip), blocked on RDR 0024
  `REQ-7`/`REQ-9`; filed at reconcile.
