# Recommendation 0030: Counting and stepping a tag without writing out every value by hand

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-09-19
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
    @refine — <one-line reason>]` — the one place this
    grammar is spelled: the date, `;`, `re-verify <IDs>` or
    `re-verify none`, then `@<stage>` naming the target
    re-entry stage (`@propose` = Stage 2, `@refine` = 3,
    `@resolve` = 4, `@finalize` = a RE-LOCK-ONLY re-entry
    whose note-listed wording fixes land in the lock pass),
    then `— <reason>`. Every clause after
    the date is optional; `/rdr-status` routes on `@<stage>`
    and falls back to `/rdr-resolve` when it is absent. It
    is still a `Draft` for every binary Draft/Final gate;
    only Stage 4 (scoped re-verify) and Stage 7 (re-lock)
    parse the rest of the qualifier.
    The Stage 7 flip to `Final` overwrites the whole value,
    so the qualifier self-clears at re-lock — no separate
    cleanup. This 07.1 "demotion" is a *verb* describing the
    Final→Draft flip; it is **not** the `Demoted` status
    above (which exits the lifecycle to an issue) — do not
    conflate the two. (`Reverted` above is the unrelated
    terminal "implementation rolled back" status — also do
    not conflate.)
  - A Draft sent BACKWARD by a mid-flow stage carries the
    sibling qualifier:
    `Draft [routed back from resolve YYYY-MM-DD; re-verify
    A5,A6 @propose — <one-line reason>]` — the origin stage
    that sent it back, the date, `;`, `re-verify <IDs>` or
    `re-verify none`, then the SAME `@<stage>` slot spelled
    above (here also `@prelock` = Stage 5 and `@reconcile`
    = 6, which a demotion can never name because a Final has
    closed them), then `— <reason>`. Every clause after the
    date is optional, and `/rdr-status` routes it through the
    same rules. The **origin** is carried because the rework
    owed differs by where it came from; without the qualifier
    the record's evidence reads FORWARD and the receiving stage
    refuses it as already-passed. Unlike the demotion form it
    does not wait for a re-lock — **the stage named by
    `@<stage>` clears it as its first act**. A second
    route-back overwrites the value rather than stacking.
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
- **Profile**: foundational — C1, the step write and its expansion into literal rows; user-facing yes; locks cross-rdr
- **Priority**: Medium
- **Related Issues**: intrastate#jwgr
- **Predecessors**: 0002-transition-table-as-reviewable-data, 0012-declared-kind-carrier-at-the-guard-seam
- **Overrides**: 0002-transition-table-as-reviewable-data:C4 — extends the
  write-value grammar admitted at load; the replace-not-accumulate
  application semantics C4 states are left untouched.
- **Seam Lineage**: no prior accretion.

## Problem Statement

A model author writing a transition table for a retry, escalation, or
attempt-counting flow wants to say "one more attempt" or "step to the next
tier" once, and have the table mean it for every value the tag can hold.
Today they cannot: a write value must be a literal, so the author writes out
one row per (current value, outcome) cell by hand. They discover this the
moment the flow has a counter in it — a ladder that is one sentence of intent
("attempt+1 while attempt is under the cap; on a deterministic failure step
the tier") becomes a few dozen rows that differ only in the literal being
written, and every change to the cap or to the tier domain means re-typing the
unrolled block. The author is also the one who must review the table, so the
cost lands twice: the intent they wanted reviewable as data is spread across
rows that no longer read as a single decision.

The system-internal requirement is that a write value is a literal **by
construction**. `internal/table/load.go::valueMembers` admits
string/bool/int64/float64/[]any and refuses anything else as
`malformed_tag_declaration`, and `internal/table/normalize.go::renderWrites`
conforms a write as an `eq` literal with no other arm. So the question this
record must answer is what a write value may *be*: does it stay a literal, or
may it be an authorable computed form over the tag's own declared kind?

Three things ride on that answer and are not statable before it. What a
computed write does at the tag's declared bound — an `int` that could exceed
its `max`, and an ordered `enum` already at its last domain member. And what
the normalized graph enumerates for a successor that is no longer a constant:
`internal/graphlint/reach.go` builds each node's successor by cloning the
row's writes, `engine.go` canonicalizes `NextTags` the same way, and
`groups.go::checkIdempotentWrites` compares write values for a
set-back-to-same test that has no meaning for a value computed at resolve
time. RDR 0021 publishes those concrete values as the graph export, and RDR
0022 and RDR 0015 are both `Draft` against a successor-is-a-constant
assumption.

This record does not reopen what RDR 0002:C4 decided. C4 states that a write
**replaces** — it supplants the whole held value and does not pile a second
value onto a single-valued tag, which is what 0004's equality read-back and
0006's single-valued-state invariant rest on. That is the semantics of
*application*. The grammar of the write *value* is not addressed by any
record; the refusal in `valueMembers` is the absence of a decision, not a
decision. The **Overrides** field above keeps the two apart.

## Critical Assumptions

Pending at Propose; Stage 4 verifies. Each names the artifact that decides it.

- **A1 The loader's existing per-member row expansion can host a guard-block
  expansion keyed on the stepped tag without disturbing 0002:C13's
  match-only rule for `in` atoms.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/normalize.go::expand` — its choice-point
    loop over sorted candidates appends the member to `Row.Suffix`, and its
    `case expanding:` arm MINTS a fresh `eq` atom rather than rewriting the
    authored one, so the comment's reason for refusing a guard `in` (a
    conjunction would become a disjunction, and `0002:C7` forbids folding an
    `unless` into a per-member `eq`) does not reach a step point. Two
    structural facts bound the work: the choice-point discriminant is
    `c.atom.Operator == "in" && c.atom.Block == BlockMatch`, so a step point
    is a third kind whose suffix element must slot into
    `internal/table/normalize.go::compareAtoms`' sort; and `renderWrites`
    runs BEFORE `expand`, whose tail clones one write set onto every row, so
    the per-cell write must move inside the loop.
  - **If wrong**: the step expansion is a second mechanism beside `expand`
    with its own suffix and ordering rules, and Phase 1 doubles.
- **A2 An enum's authored `domain` order survives load, normalization, the
  dump and the graph export unsorted, so it can carry step order.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/load.go::tagDecl` carries `Domain:
    src.Domain` unsorted; `0021:C2` "`domain` carries the declaration's
    AUTHORED members"; a sweep of every `sort.`/`slices.Sort` in
    `internal/table`, `internal/cli/graph_document.go` and
    `internal/graphlint` finds none over `TagDecl.Domain` — the export
    clones it, `conformDomain` only reads membership. The one domain sort
    is over `EmitDecl.Domain` (RDR 0024's emit alphabet), a different type
    whose carry is documented value-preserving, not order-preserving.
  - **If wrong**: `step` on an enum has no stable order and needs an
    explicit ordering facet, widening the declaration model RDR 0003 owns.
- **A3 An enum domain member is NOT among the strings today's `#` ban
  covers, so a member carrying `#` can reach a row's expansion suffix
  unless this record refuses it.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the load-time `#` ban has exactly three sites, and the one
    that bans a member is `internal/table/normalize.go::normalizeAtom`'s
    `case "in":` arm, gated `if b == BlockMatch` with the comment "The `#`
    reservation is match-only: only match blocks expand"; the other two are
    `::normalizeRules` (rule ids) and `internal/table/load.go::loadOutcomes`
    (the recognized alphabet). `tagDecl` performs NO `#` check on
    `src.Domain`. So a domain member carrying `#` is authorable today and
    reaches `internal/table/model.go::Identity`'s suffix, and C1's refusal
    is a new guard, not a restatement. `0002:C17` scopes the ban to
    match-block `in` members, which a domain member is not.
  - **If wrong** (the ban already covers domain members): C1's refusal of
    a `step` over a domain with a `#` member is redundant with an existing
    load check and drops to a restatement, not a new guard.
- **A4 The per-cell admits filter is decidable at load with the same
  comparison the runtime guard evaluator uses, for `eq`/`in`/`lt`/`lte`/
  `gt`/`gte` on `int` and `eq`/`in` on `enum`.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/guard/grammar.go::Evaluator.Evaluate` is the sole
    implementation of `internal/resolve/resolve.go::GuardEvaluator` (0012's
    typed evaluator, `0012:C2`); `internal/resolve/guardcontract.go` is its
    cross-RDR conformance harness, not the evaluator. `lt`/`lte`/`gt`/`gte`
    parse both sides with `strconv.Atoi` and delegate to
    `internal/guard/grammar.go::compare`; `eq` is string equality and `in` is
    `slices.Contains` over a parsed set. The type is `struct{}` — documented
    as carrying no state so that "never reads the tag view" is a property of
    the type — so each comparison is a pure function of (operator, atom
    literal, candidate value) and replays at load over a declared finite
    domain. `Evaluate` is THREE-valued, and C1 now disposes of the third
    arm explicitly (exclude, never admit) while recording that it is
    unreachable for the kinds C1 admits: `internal/table/model.go`'s
    operator/kind matrix confines `lt`/`lte`/`gt`/`gte` to `int` and
    `internal/table/load.go::conform` kind-checks their bound, so both
    `Atoi` calls succeed over an `IntDomain()` member; `in`'s literal is
    always a rendered set and always parses; `eq` is total;
    `internal/table/normalize.go` refuses an unknown operator at load
    against the frozen vocabulary; `contains` is `set`-only and `exists`
    is kernel-decided (`0007:C1`). `0007:C2`'s partiality is scoped to a
    tag key absent from a runtime VIEW, which a load-time domain
    enumeration has no analogue for.
  - **If wrong**: load and runtime disagree on which cells a stepping row
    admits; the bound refusal fires on a cell the runtime never reaches or
    misses one it does.
- **A5 Expanded rows introduce no overlap the hand-unrolled table did not
  have: two rows in one outcome group that step the same tag partition its
  cells by their own atoms exactly as literal rows do.**
  - **Status**: Pending — resolves at the MVV; the form is unimplemented, so
    the pair cannot be run before Phase 2.
  - **Method**: MVV Test
  - **Evidence**: the MVV's fixture pair — `intrastate lint --model
    <fixture> --as json` on each, finding sets compared as sets over
    `findings[]` (normalizing `model`), including zero `graph-overlap`. All
    four cited finding codes exist today:
    `internal/graphlint/taxonomy.go::CodeOverlap`, `::CodeCoverageGap`,
    `::CodeOwnedBeforeWrite` (blocking) and `::CodeIdempotentWrite`
    (advisory tier — so "identical finding sets" spans both tiers).
  - **If wrong**: a computed ladder lints `graph-overlap` where the
    literal one did not, and the "same table, fewer rows" claim fails.
- **A6 For every (state, outcome) cell the `step` ladder resolves to the
  same plan as the unrolled ladder.**
  - **Status**: Pending — resolves at the MVV, with A5.
  - **Method**: MVV Test
  - **Evidence**: `intrastate flow resolve --model <fixture> --outcome <o>
    --artifact <role>=<statefile> --as json --plan-only` swept over every
    (state, outcome) cell of the pair; `internal/cli/flow_resolve.go`'s
    payload carries `next`, `writes` and `clear` as named fields, and
    `--plan-only` drops the request echo so the two models' plans compare
    directly. The state cell is driven by rewriting the bound artifact file,
    the pattern `internal/cli/flow_exit_0005_test.go` already uses.
  - **If wrong**: load-time expansion is a semantics change, not a grammar
    change, and Alternative 1 reopens.
- **A7 The step expansion's row growth is bounded by the declaration
  alone — at most the admitted-cell count, itself at most the domain
  width — and needs no load-time ceiling at the ladder sizes reported;
  the lint enumeration limit is RDR 0013's caller-supplied analysis scope
  (`0013:C1`, `0013:C3`) and gates nothing at load.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a7-row-growth.md` — a table of N literal
    rows (what C1 says the expansion mints) at the A7 declaration `min = 0,
    max = 100000` loads and normalizes in **0.96 s / 584 MB at 100,001
    rows**, exit 0, and stays linear (~10 us and ~5.8 KB per row) to 500,000
    rows / 5.08 s. No row-count gate exists in `internal/table`, and
    `internal/guard/declaration.go::AssignmentCount` is not merely unreached
    but UNREACHABLE from load: `go list -deps ./internal/table` names
    neither `internal/guard` nor `internal/graphlint`, so the package graph
    enforces it. Both published bounds are lint-side
    (`internal/guard/product.go` product bound, `internal/graphlint/taxonomy.go`
    node ceiling), and exceeding one degrades the verdict —
    `reach.go` sets `complete=false`, reported as a finding — rather than
    refusing the model. The joint check below records the 0013 coupling.
  - **If wrong**: C1 gains a load-time ceiling of its own — a 0030
    decision, never a reuse of 0013's lint limit, which the analysis scope
    it parametrises cannot reach.
- **A8 An expanded row's suffixed identity is what `flow next`, `flow
  resolve` and the dump already show for `in`-expanded rows, so no new
  rendering surface is owed.**
  - **Status**: Refuted — the dump half holds, the flow half does not.
  - **Method**: Source Search
  - **Evidence**: the suffix lives in a separate `Row.Suffix` field and is
    joined to the rule id ONLY by `internal/table/model.go::Identity`;
    `expand` leaves `RuleID` the bare authored id. `internal/table/dump.go`
    does render it (`case "identity": return r.Identity()`, beside
    `SourceLocator`), so the dump claim stands. But
    `internal/table/model.go::Row.KernelRow` builds the kernel row as
    `resolve.Row{RuleID: r.RuleID, …}`, dropping `Suffix` entirely, and
    `internal/cli/flow_next.go::summarize` and `internal/cli/flow_resolve.go`
    then publish that bare id: an `in`-expanded row surfaces TODAY as
    `retry`, never `retry#3`. A related join is ambiguous for the same
    reason — `internal/cli/flow_resolve.go::rowByID` matches on
    `row.RuleID == ruleID`, so N expanded rows share one key and the first
    wins (sound by `0010:C3`; see the audit row).
    The refutation narrows the claim rather than opening a gap: the split is
    per SURFACE, not per verb — every surface that names a ROW publishes
    `identity` (`internal/table/dump.go`'s column;
    `internal/cli/graph_document.go`'s `graphRowDoc.Identity`), every surface
    that names a RULE publishes `rule` (`flow_next.go`, `flow_resolve.go`,
    and the graph document's edge). The published vocabulary agrees —
    `0005:C1` says "matched rule identity", `docs/cli-output-contract.md`
    says "a candidate is a rule". So no rendering surface is owed; C3 states
    the split.
  - **If wrong** (as verified, and dispositioned): the narrowed claim is what
    C3 now carries — the dump/graph precedent only. Were the flow verbs
    instead owed the suffix, it would widen the `resolve.Plan` seam RDR
    0002/0007 owns and reach the `intrastate.graph/1` edge vocabulary C3
    holds fixed — a successor record, not a clause here.
- **A9 No open peer relies on the inline-table write refusal
  (`malformed_tag_declaration: … is not a tag value`) as a contract.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: joint check arm 3 read the `§normative-contracts` fence of
    every open (Draft) peer — 0012, 0013, 0014, 0015, 0016, 0017, 0018,
    0022 — for `is not a tag value` and `malformed_tag_declaration`: no
    citing element id, clear. 0019 is Implemented, not an open peer, and
    names `valueMembers` only for canonical-form rendering. 0015 and 0022
    consume `reach.go::successor` over an ALREADY-NORMALIZED row, so
    expanding to literal rows at load preserves the shape their fences rest
    on — more rows, same literal `Writes`.
  - **If wrong**: a peer's fenced refusal silently becomes an acceptance
    — a joint decision, not an edit here.

## Proposed Solution

### Approach

Admit one computed write form, `<tag> = { step = <n> }`, and have the
**loader expand it into literal rows** — one row per current value the rule
admits, each carrying a `guard.all` `eq` atom on the tag, a literal write of
the stepped value, and the current value as its row-identity suffix. This
is the mechanism the loader already uses to expand a match-block `in` atom
into one row per member (`internal/table/normalize.go::expand`,
`0002:C13`), pointed at a guard dimension instead of a match scope. The
form is admitted on an `int` declared with both bounds and on an `enum`
with a domain, whose authored member order is the step order; on every
other kind and on an unbounded `int` it is refused at load. A stepped value
that leaves the domain at any admitted cell is a **load refusal** naming
the rule, the tag and the cell — no saturation, no wrap. Everything after
normalization (kernel resolve, write application and read-back, lint
coverage and reachability, the graph export, and the two Draft consumers
that assume a constant successor) sees only literal writes and is
unchanged. `0002:C4`'s "a write replaces" holds per expanded row, which is
why this record extends the value grammar C4 admits without reopening the
application semantics C4 states.

### Technical Design

The write value enters at `renderWrites`, which today hands every value to
`internal/table/load.go::valueMembers`; a TOML inline table has no arm
there and is refused as `malformed_tag_declaration`. The design intercepts
the table shape on the **write-block path only** — `[initial]` values and
predicate literals keep the literal-only grammar; `valueMembers` itself is
unchanged, since the interception precedes the call. A rule with a
`step` write is recorded as a **step point** on the rule, and `expand`
treats it as one more choice point in its Cartesian product: the members
are the tag's admitted cells, the emitted atom is `guard.all` `eq =
<cell>`, the row's write for that key is the stepped literal, and the
suffix element is the cell. Two step writes in one rule, or a step write
beside an `in` match atom, compose as the product already does.

**Admitted cells.** The rule's own positive atoms on the stepped tag decide
which domain members expand: `eq`/`in`/`lt`/`lte`/`gt`/`gte` under
`guard.all` and `eq`/`in` under `match`, each evaluated per member with the
same comparison the runtime evaluator performs (`0012:C2`) — which is
three-valued, so C1 fixes the undecided arm to EXCLUDE the member. That
arm is unreachable for the kinds C1 admits (A4); fixing it is a soundness
fence, and the direction is chosen because a loader that over-admits turns
an undecidable predicate into C2's refusal, while one that under-admits
leaves an ordinary `graph-coverage-gap` the author can see and answer.
`unless` atoms
are not consulted — an `unless` block is negated as a whole and may mention
other keys, so per-cell evaluation is not decidable from this tag alone;
ignoring it can only over-admit, and an over-admitted cell surfaces as the
bound refusal with the `guard.all` complement as the named remedy, which is
the authoring guide's existing advice for a needed conjunction
(`docs/model-authoring.md` §Guard and match atoms). A rule with zero
admitted cells is refused: it is the unsatisfiable-atom case `0003:C6`
already refuses at load rather than yielding a never-matching row.

**The bound.** For each admitted cell the stepped literal passes the same
`conform` check a literal write passes today (`internal/table/load.go::
conform` → `conformDomain`: `int` against `min`/`max`, `enum` against
`domain`). A failure is the bound refusal. ⇒ The bound is therefore not a
new analysis: "guard it below the cap" makes the rule admit only cells
whose successor exists, and the cell at the cap is then an ordinary
uncovered cell in the outcome group's product, which the existing coverage
proof reports as `graph-coverage-gap` unless another row (the exhausted
row) claims it. The bound is not a guard-conditioned interval analysis
because the interval is enumerated, not analysed.

**Why guard atoms, not match atoms.** "A **match** atom *scopes* the row
group … and contributes **no dimension** to the product lint proves
exhaustiveness over. Only `guard.all` and `guard.unless` atoms collect as
dimensions" (`docs/model-authoring.md` §Discriminate with guard atoms). ⇒ An
expansion into match atoms would leave the cap cell unprovable — a silent
green over the exact hole the author is most likely to leave. The emitted
atom goes to `guard.all`. The state read this expansion depends on — the
tag's current value — is written by `[initial]` and by prior rows' writes,
which is why an unstepped-before-initialised key surfaces through the
existing `graph-owned-before-write` finding, not a new one.

#### Normative Contracts

**C1**

```normative
STEP WRITE. A rule's write block MAY assign an owned tag the inline table
`{ step = <n> }`, `n` a non-zero TOML integer, and no other table shape:
any other key set, a zero step, or a non-integer step (a TOML float `1.0`
included) is refused at load under `malformed_tag_declaration` with a
detail naming the one admitted form. The form is admitted only on a tag whose declared kind is `int` with
both `min` and `max`, or `enum` with a non-empty `domain`; on `bool`,
`set`, `scalar`, or an `int` missing a bound it is refused under the same
category. An `int` whose declared `min..max` spans the full integer width
is refused too: its domain is unenumerable, so the expansion has no cell
set to compute over. `[initial]` values and predicate literals keep the literal-only
grammar; the table shape is admitted on the write-block path alone.

The loader expands a rule carrying a step write into literal rows, one per
ADMITTED CELL of the stepped tag, composed into the rule's existing
expansion product (`0002:C13`) as one more choice point. An admitted cell
is a domain member — `{min..max}` for `int`, the authored `domain` order
for `enum` — that satisfies the CONJUNCTION of every positive atom the rule
authors on that tag, however many per block (`guard.all`
`eq`/`in`/`lt`/`lte`/`gt`/`gte`; `match` `eq`/`in`), evaluated per member
as the runtime evaluator would and against the atoms as already chosen in
the expanded row — a match `in` on the stepped tag has expanded to its
`eq` member first, so the two compose to exactly its members; `unless`
atoms are not consulted. The runtime evaluator is three-valued; an atom
answering neither true nor false at a member does NOT admit that member.
For the kinds this clause admits that arm is UNREACHABLE — the operator/
kind matrix confines the ordered operators to `int`, their bounds are
kind-checked at load, an unknown operator is already refused against the
frozen operator vocabulary, `exists` is decided in the kernel and never
reaches this seam, `contains` is `set`-only and excluded here, and `eq`
and `in` are total over a declared domain member — so the rule is a
soundness fence, not a live branch. It is written to exclude rather than
admit because the two directions fail differently at LOAD: excluding
leaves the cell to the outcome group's ordinary `graph-coverage-gap`,
which is visible, non-blocking and actionable, whereas admitting mints a
row the author did not ask for and lets C2's bound report "I cannot tell"
as "your model is broken". The authored atoms are retained; the expansion
adds its own atom beside them. Each expanded row carries a `guard.all`
atom `eq = <cell>` on the tag, the literal write `cell + n` (`int`) or the
domain member `n` positions from the cell (`enum`), and the cell appended
to its expansion suffix. The stepped value is computed without overflow
at every admitted cell, whatever the step's magnitude; whether it lands
inside the domain is C2's bound, not an admission question. A rule
admitting zero cells is refused at load. A step write over a domain
containing a member with the suffix separator `#` is refused at load.

After normalization no surface distinguishes an expanded row from an
authored literal row: `Row.Writes` holds literals only, and `0002:C4`'s
write-replaces clause applies per row unchanged.
```

**C2**

```normative
Surface — of C1; the refusal the expansion's conformance check implies.
BOUND. At every admitted cell the stepped literal MUST conform to the tag's
declaration exactly as an authored literal write does (`int` within
`min..max`; `enum` within `domain`). A cell whose stepped value leaves the
domain in either direction — past `max` or below `min`, past the last
member or before the first — is a load refusal under `malformed_tag_declaration` — the category a
non-conforming literal write already takes — whose detail names the rule,
the tag, the cell, and the stepped value. There is no saturating and no
wrapping form. The remedy is the author's: a positive atom excluding the
cell, or a literal row for it; the cell it leaves unclaimed is then the
outcome group's ordinary coverage obligation (`graph-coverage-gap`).
```

**C3**

```normative
Surface — of C1; the invariance C1's last paragraph implies on shipped
surfaces. INVARIANCE. `flow next`, `flow resolve`, `flow set-state`,
`lint`, `dump`, and `graph` read expanded rows through `Row.Writes` and
`Row.Atoms` alone; none learns a computed value shape, and the graph
document's `intrastate.graph/1` vocabulary gains no member. The
authoring guide (`docs/model-authoring.md`) documents the form, the
admitted kinds, the admitted-cell rule, and the bound refusal beside the
existing write-block section; `intrastate --help-all` is regenerated. No
emitted vocabulary — `internal/table::Categories()`, the CLIError codes,
the finding codes, the graph document's members — gains a member, so no
`0029:C4` stability tier is owed by this record. Every surface that names a
ROW publishes the suffixed identity `rule#cell` — the shape `in` expansion
already emits there (`dump`'s `identity` column, the graph document's
`identity` row field) — and every surface that names a RULE publishes the
authored rule id (`flow next`'s and `flow resolve`'s `rule`, the graph
document's edge `rule`). That split is the existing design, not a gap: a
flow payload names the rule it matched, so this record mints no obligation
on it.
```

#### Load-Bearing Decisions

- **Identity** — an expanded row's identity is the existing tuple
  `(model id, rule id, suffix…)` with the cell appended as one more suffix
  element in the choice-point sort order; no new identity, and a
  single-cell expansion still appends its element (a stepped row is never
  the authored row, so it always carries the cell — the opposite of the
  single-member `in` rule, whose one row IS the authored row). Two step
  points in one rule order their suffix elements by the choice-point
  sort, as `in` atoms do; the element is the bare cell, not
  tag-qualified, on the `in` precedent — and `#` is banned in rule ids,
  so `retry#3` can never be an authored id.
- **Wire / byte format** — the TOML inline table `{ step = <n> }`, exactly
  one key. A string sentinel (`"+1"`) is rejected: `strconv.Atoi("+1")` is
  `1`, so it already parses as a literal on an `int` tag. A bare table with
  the reporter's `add`/`next` keys is rejected as a second spelling of one
  operation; the detail names `step`.
- **Naming** — `step`, one operator over both ordered kinds (a negative
  step is the predecessor; `next = true` carried an operation in a
  boolean and no distance; `add` reads as arithmetic and has no enum
  meaning). The refusal category is reused, not minted: a write value's
  kind/domain conformance already files under `malformed_tag_declaration`
  (`renderWrites`, deviations D3 precedent), and the bound is that check
  on the expanded literal.
- **Selection / predicate** — which cells expand is decided by the rule's
  own positive atoms on the tag, never by `unless`, never by other rows;
  the cells it leaves are the group's to cover.

#### Illustrative Code

Illustrative — the ladder from the Background, both forms; tests assert
the fixture pair the MVV names, not this excerpt.

```toml
# today: one row per attempt value, differing only in the literal
[[rule]]
id = "retry-2"
[rule.match.recognized]
eq = "transient"
[rule.guard.all.attempt]
eq = 2
[rule.write]
attempt = 3
# … retry-0, retry-1, retry-3, retry-4 …

# proposed: one row; the loader mints retry#0 … retry#4
[[rule]]
id = "retry"
[rule.match.recognized]
eq = "transient"
[rule.guard.all.attempt]
lt = 5
[rule.write]
attempt = { step = 1 }

[[rule]]
id = "escalate"
[rule.match.recognized]
eq = "deterministic"
[rule.guard.all.tier]
in = ["small", "medium"]
[rule.write]
tier = { step = 1 }
attempt = 0
```

The `retry` rule over `attempt` declared `min = 0, max = 9` admits cells
0–4 (the `lt 5` atom), writes 1–5, and leaves cells 5–9 to the exhausted
row; without the `lt 5` atom it admits 0–9 and cell 9 refuses at load. The
`escalate` rule admits `small` and `medium` and leaves `large` to the
exhausted row; `tier = { step = 1 }` with no atom would admit `large` and
refuse.

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Mint one row per member with a suffixed identity | `internal/table/normalize.go::expand` | match-block `in` only, by design (`0002:C13`) | Extend | one more choice-point kind; guard atom emitted instead of match |
| Refuse a write value outside its declaration | `internal/table/load.go::conform` / `conformDomain` | literal only | Reuse | the bound refusal is this check on the expanded literal (C2) |
| Prove the cap cell is claimed | `internal/guard` product coverage, `graph-coverage-gap` | guard atoms only | Reuse | why the expansion emits `guard.all` (C1) |
| Admit a write value | `internal/table/load.go::valueMembers` | no inline-table arm | Reuse, unchanged | intercepted before it on the write path only (C1) |
| Successor per node, export, fingerprint | `reach.go::successor`, `graph_document.go::graphRows`, `engine.go::Fingerprint` | read `Row.Writes` literals | Reuse, unchanged | C3 |
| Flag a write that moves nothing | `groups.go::checkIdempotentWrites` | `eq`-match only | Reuse, unchanged | a non-zero step never lands on its own cell |
| Decide whether one atom admits one candidate value | `internal/guard/product.go::valueSatisfies` (render-then-`Evaluate` shim, ~13 lines) | lives in `internal/guard`, which imports `internal/table`, so `table` cannot import it (the cycle `0012:A5` identifies) | Reuse (extract) | extract the shim to `internal/resolve` — the leaf both reach, which already owns `GuardAtom`, `GuardResult` and `GuardEvaluator` and imports nothing intra-repo — and repoint `guard` and `graphlint`. It is already duplicated three ways (`guard/product.go::valueSatisfies`, `graphlint/reach.go::atomAdmitsValue`, a test-local copy in `internal/resolve`), plus the set renderer twice (`guard/assignment.go::renderSet`, `graphlint/reach.go::renderSetLiteral`); this record would mint a FOURTH, and the first in the loader, where a divergence is a load refusal rather than a lint finding. C1's "as the runtime evaluator would" then holds structurally |
| Enumerate an int tag's declared domain | `internal/guard/declaration.go::IntDomain` | same package placement | Cited, not reused | the loader enumerates `min..max` directly; no cycle incurred |
| Intersect an atom's denotation across a rule's guard block | `internal/guard/product.go::Denotation` / `::acceptedIn` | carries `AssignmentSet`, presence dimensions, group products and `unless` subtraction; yields the unprojectable empty set on an undecidable dimension — lint's false-green guard, the wrong direction for a loader (C1's third arm) | Cited, not reused | the loader's conjunction is a `&&` over the shim above, per cell |
| Join a rule's `emit` block back after selection | `internal/cli/flow_resolve.go::rowByID` (first-match on bare `RuleID`) | one key per rule, not per expanded row — contracted, not a defect (`0010:C3`) | Reuse, unchanged | none: 0030 INHERITS `0010:C3`'s soundness condition and preserves it — the step expansion rides the same `expand` loop and copies `Emit`, `Gate`, `RequiresOwned` and `Escape` identically to every row it mints, so every row a stepped rule expands to carries the same block |
| Cap row growth | `internal/guard/declaration.go::AssignmentCount` ceiling | a lint analysis scope (`0013:C1`); unreachable from load — `go list -deps ./internal/table` names neither `guard` nor `graphlint` | Not reusable (A7) | none: the A7 spike measured load linear and cheap, so no load ceiling is owed |

### Decision Rationale

**Q: what may a write value be?** Scored matrix — approaches × criteria;
`+` favourable, `0` neutral, `−` unfavourable, one clause per cell.

| Criterion | 1 Load-time expansion (chosen) | 2 Symbolic write to the kernel | 3 External unroll generator | 4 Expression language |
| --- | --- | --- | --- | --- |
| Correctness fit: bound decided declaratively | `+` bound = existing `conform` + coverage cell | `0` new guard-conditioned interval analysis, or a runtime refusal | `+` generator emits literals | `−` undecidable in general |
| Prior-art alignment | `+` the Ragel review's option (a), made cheap | `0` option (b) without the host language | `+` option (a) by tooling | `+` every peer, and the drift they document |
| Reversibility | `+` remove the form; expanded rows were never persisted | `−` a second value shape across ten consumers and an export member | `+` delete the tool | `−` models come to depend on it |
| Blast radius (seams touched) | `internal/table` + docs, plus a behaviour-preserving move of one shim into `internal/resolve` | table, resolve, accessor, cli ×4, graphlint ×3, export, 0022/0015 | none in-repo; a new artifact | table, resolve, graphlint, export, 0003's operator model |
| Cost | one choice-point kind in `expand`; one shim extracted and two callers repointed; fixtures | value shape, evaluator, lint analysis, export `/1` addition, preview semantics | a generator + a drift check | a parser and evaluator |
| Review legibility (the user's problem) | `+` source is one row; dump shows the cells | `+` source and dump both one row | `−` the reviewed file is the unrolled one | `+` one row, opaque semantics |
| Open-peer safety (0022, 0015 assume constant successor) | `+` unchanged | `−` both re-authored | `+` unchanged | `−` both re-authored |

Row 1 and row 7 decide it. Approach 1 turns the bound question the seed
could not state into two things that already exist — a literal conformance
check and an uncovered coverage cell — and leaves every consumer of a
row's writes, including the two Draft records, reading literals. Approach
2 is the only other one that solves the user's problem without a second
artifact, and it loses on rows 3–4: ten shipped read sites of `Row.Writes`
(`internal/guard/lint.go::writtenKeys`, `internal/accessor/executor.go`,
`flow_resolve.go`, `flow_next.go`, `flow_state.go`, `graph_document.go`,
`reach.go::successor`, `engine.go::Fingerprint`,
`groups.go::checkIdempotentWrites`, `dump.go::renderValue`) would each
learn a computed shape, and the bound would need a guard-conditioned
interval analysis. Approach 3 fails row 6, which is the problem statement.
Approach 4 fails rows 1 and 3 and is what RDR 0002 chartered against. The
cost Approach 1 pays is row 6's second half: the normalized dump and the
graph document's rows show `retry#0 … retry#4`, not `retry` — the same view
an `in` expansion already gives *there*, accepted on that precedent
(`internal/table/dump.go`'s `identity` column;
`internal/cli/graph_document.go`'s `graphRowDoc.Identity`). The flow verbs
are a different surface and keep naming the authored rule
(`internal/cli/flow_next.go`, `internal/cli/flow_resolve.go` carry `rule` ←
`plan.RuleID`, which `0023:C2` binds into every projection mode), so a
consumer wanting per-row identity from a flow payload reads `graph` or
`dump`. Publishing the suffix on `rule` was weighed and rejected: the kernel
plan carries only `RuleID` and one `SourceLocator` shared by every expanded
row (`internal/table/normalize.go::expand`), so the suffix is genuinely
absent at a seam RDR 0002/0007 owns; `rule` is also the graph document's
EDGE field, which C3 promises not to touch; and
`TestReq92_ANonFirstExpandedRowResolvesToTheAuthoredEmitBlock`
(`internal/cli/decision_table_0010_test.go`) asserts the authored id on a
non-first expanded row today. Per-row identity IN a flow payload is a
possible successor record — disclosed here as a non-obligation, not an owed
fix.

**Q: where does the admits comparison live?** Extracted, not duplicated.
The comparison C1 specifies already ships as a ~13-line render-then-
`Evaluate` shim, and it already exists three times — in lint's projection,
in the reachability traversal, and test-locally — with its set renderer
twice more. A fourth copy would be the first one in the LOADER, where a
divergence from the runtime is a load refusal rather than a lint finding,
and the precedent for catching that early is in the repo: `flow_next.go`
carries a comment about a hand-rolled filter that "answered a question it
does not own, and its answer disagreed with the kernel silently", and
`internal/table/model.go`'s `operatorKinds` claims to mirror
`guard.Accepts` with no cross-package test holding the two together — so
the duplicate-plus-contract-test option is already shown not to hold here.
House doctrine points the same way (`0029`: the shared constant lives in
the leaf both sides reach). `internal/resolve` is that leaf: it imports
nothing intra-repo and already owns `GuardAtom`, `GuardResult` and the
`GuardEvaluator` interface, and the shim needs no `table` type, so the
cycle `0012:A5` identifies is not incurred. `Denotation`/`acceptedIn` are
CITED and not reused — they carry lint's `AssignmentSet`, presence
dimensions and `unless` subtraction, and return the unprojectable empty
set on an undecidable dimension, which is the false-green direction a
loader must not take (C1's third arm, from the other side).

Premortem: hardened — the draft-free critic returned PASS with
fourteen rules the brief omitted (conjunction of atoms per key, `in`+step
composition, lower-bound and oversize steps, float/zero steps, the
write-block-only interception on every path, the untiered-surface class);
each is folded into C1–C3, D-identity, A7, A8, the MVV and Failure Modes
(`evidence/propose-premortem/critic.md`).

Ground-sweep: clean (33 anchors) — 32 confirmed at source or by the
projector, 1 cosmetic (the sweep brief, not the record, said `Identity`
joins the model id with `#`; it is `.`, and only the suffix elements take
`#`); no load-bearing miss (`evidence/research/ground-sweep.md`).

Joint-check: fired → 0013 (home: `cli/0013 §Normative Contracts` C1).
Arm 1: `internal/guard/declaration.go::AssignmentCount` is shared with
0013, uncited before this line. Neither record modifies it; 0013 makes the
enumeration limit a caller-supplied lint scope (`0013:C1`, `0013:C3`), so
A7 cites 0013 as the owner and any load-time ceiling is 0030's own.

Joint-check: fired → 0022 (home: `cli/0030 §Normative Contracts` C3).
Arm 1: `internal/graphlint/reach.go::successor` is shared with 0022,
cited. `0022:A1` rests on a constant literal-applied successor; C3 is
where this record preserves it, for 0022 to cite rather than assume. Arm
2 also intersects 0022 and 0015 on `lint` and the authoring-guide path —
invariance claims on both sides, no shared decision.

Joint-check: fired → 0012 (home: `cli/0003 §Normative Contracts` C6).
Arm 2: the kind tokens `bool`/`enum`/`int`/`scalar`/`set` are shared with
0012 inside both records' fences — RDR 0003's closed vocabulary, which
neither extends.

Arm 3 (absence, manual): clear — no peer fence relies on the inline-table
refusal this record converts; 0019 names `valueMembers` for float
rendering only. All three arms ran. Ground-before-ask resolved `apply`
(searched=cluster, found), so the homes are recommended from the peers'
own clauses, not asked.

## Alternatives Considered

### Alternative 1: Symbolic step carried into the kernel

The reporter's framing: `Row.Writes` gains a computed value shape (`step
n`), `internal/resolve` computes the concrete plan value from the tag view
at `planOf`, graphlint maps the step over each node's held values in
`successor`, and lint gains a guard-conditioned interval check for the
bound.

- **Pros**: the normalized table stays one row per authored rule; the dump
  and `flow next` name the authored rule id; no row-count growth.
- **Cons**: a second value shape at every one of the ten `Row.Writes` read
  sites listed in the rationale; `Fingerprint`/`literalKey` and the
  `intrastate.graph/1` export need a spelling for it (an additive member,
  but a consumer-visible one that RDR 0029's tier vocabulary would oblige);
  `checkIdempotentWrites` and the 0004 read-back stay meaningful only
  because the plan is concrete — `flow next`'s preview, which has no view,
  cannot show a value; the bound is either the new analysis or a runtime
  `write_out_of_domain` refusal for what is statically decidable; RDR 0022
  and RDR 0015 are re-authored against a non-constant successor.
- **Prior art**: the Ragel review's option (b) — "count in an action
  variable" — minus the host language; every peer does this, and the
  review records the cost as "the .rl is no longer the single source of
  truth".
- **Rejected because** the chosen approach yields the same runtime plans
  and the same lint verdicts while leaving the analyzability seam
  untouched; its one cost (the suffixed row view) already exists for `in`.

### Alternative 2: Keep writes literal; ship an unroll generator

A tool (or `intrastate expand`) reads a compact ladder in a side format and
writes the unrolled literal TOML.

- **Pros**: zero change to the loader, the seam, or any record.
- **Cons**: the unrolled file is what lint loads and what a reviewer diffs,
  so the review cost the problem statement names is untouched; two files
  for one intent drift; a side format is a second grammar.
- **Rejected because** it solves the typing cost and not the review cost,
  which the problem states lands twice.

### Alternative 3: Expression language in write values

`attempt = { expr = "attempt + 1" }`, the SCXML/sismic class answer.

- **Pros**: unbounded expressiveness; familiar to authors of every peer
  tool.
- **Cons**: coverage and reachability proofs need the value set of every
  expression — undecidable in general and out of RDR 0003's closed
  operator model; the table stops being reviewable data (RDR 0002's
  charter); the Ragel review names the drift: "Everything interesting about
  this model lives in the actions, not in the regular expression".
- **Rejected because** it is the failure mode the data-only table exists
  to avoid, and one operator over two ordered kinds covers the reported
  need.

### Briefly Rejected

- **Saturate or wrap at the bound**: a write that moves nothing is what
  `graph-idempotent-write` already flags as a smell; silent saturation
  hides an exhausted counter, against `0002:C3`'s "never a silent no-op".
- **String sentinel spelling (`"+1"`, `"next"`)**: `+1` already parses as
  the literal `1` on an `int`; `"next"` could be a legal enum member.
- **An explicit `order = [...]` facet for enums**: a second copy of the
  domain that can disagree with the first; authored `domain` order is
  carried today (A2).
- **Computed values in the next-state block (`advance`)**: out of scope;
  the write block is the value grammar this record decides.
- **`bool` toggle**: two literals; write the literal.

## Context

### Background

Raised as a model-authoring complaint from a real ladder: a retry/escalation
table whose whole intent is "increment the attempt while it is under the cap;
on a deterministic failure step the tier small → medium → large" had to be
unrolled into rows differing only in the literal next value — transient-retry
rows, one deterministic row per (tier, attempt) cell, and exhausted rows. That
ladder is roughly 24 rows unrolled, against 8 with both step forms or 10 with
the integer form alone.

In this repo today, `attempt = 1` is accepted; `attempt = { add = 1 }` and
`tier = { next = true }` are both refused `malformed_tag_declaration: … is
not a tag value`. The refusal comes from `internal/table/load.go`'s value
admission, which has no arm for a TOML inline table, and
`internal/table/normalize.go` then renders a write as an `eq` literal. That
refusal is conformant with RDR 0002:C4 as written, so nothing here is a bug
report — the proposal is to add a normative choice the records have not made.

The proof is grounded in-repo: `internal/table/testdata/rdr-fixture.toml`
already carries a ladder-shaped `int` tag with `lt` guards and enum
domains, which is the fixture base the MVV builds on.

### Technical Environment

Go CLI. The affected seams are `internal/table` (write-value parsing,
admission, and normalization), `internal/resolve` (write application at
resolve time), and `internal/graphlint` (`reach.go` successor construction,
`engine.go` `NextTags` canonicalization, `groups.go` idempotent-write check).
The authoring surface documented in `docs/model-authoring.md` is what a model
author reads. Guards already admit `lt`/`lte`/`gt`/`gte`/`eq` on `int`, so the
guard grammar is not in scope.

Neighbouring records: RDR 0002 (transition table as reviewable data, C4's
write-replaces clause), RDR 0004 (accessor execution safety, the equality
read-back), RDR 0006 (graph lint authority, A10 `write-replaces`), RDR 0012
(declared-kind carrier at the guard seam — the same carrier-vs-shape question
one seam over), RDR 0021 (the normalized graph export), and the two Draft
consumers RDR 0022 and RDR 0015.

## Research Findings

### Investigation

Prior art was read before the approaches were named; the record is at
`0030-computed-write-value-grammar/evidence/research/prior-art.md`. The
class read (Ragel POC review, sismic `code.rst`, SCXML `<assign>` via
uscxml/scxmlcc) shows every peer with extended state answering "count in
an action" with a host expression language, and the Ragel review naming
the alternative intrastate forces today — "unroll the loop into three
explicit states … the cap becomes hard-wired topology". The instance read
found no peer with a data-only computed write and none deciding the bound
declaratively, so the bound disposition rests on in-repo principle
(`0002:C3` via the `valueMembers` comment) rather than external prior art.
⚠ no literature coverage for analysis-time successor enumeration
(`StateMachineLit`, one query); the choice does not rest on it. In-repo,
the choice rests on four shipped facts: `expand` already mints suffixed
rows per member; coverage is proved over guard atoms only; a literal write
outside its declaration already refuses at load under
`malformed_tag_declaration`; and `TagDecl.Domain` is carried in authored
order.

### Key Discoveries

- **Verified** — `attempt = { add = 1 }` and `tier = { next = true }`
  refuse today as `malformed_tag_declaration: … is not a tag value`
  (Background probe); `valueMembers` has no inline-table arm.
- **Documented** — coverage dimensions come from `guard.all`/`guard.unless`
  atoms only; a match atom scopes and proves nothing
  (`docs/model-authoring.md` §Discriminate with guard atoms).
- **Documented** — `expand` is a match-block mechanism by `0002:C13`, and
  its comment records why a guard `in` is not expanded; a step expansion
  emits a fresh `eq` atom per cell and never rewrites an authored guard
  atom, so that reason does not apply to it.
- **Documented** — `0021:C2` exports `domain` as the AUTHORED members and
  its vocabulary is append-only; the chosen approach adds no member.
- **Documented** — `0003:C6` refuses an unsatisfiable atom at load "rather
  than a silently-never-matching row"; a zero-cell step rule takes the same
  disposition.
- **Assumed** — the reach of today's `#` ban over domain members, the
  admits filter, and the row growth ceiling (A3, A4, A7).

## Trade-offs

### Consequences

- A ladder is one row per intent in source; the normalized view shows the
  cells as `rule#cell` rows, as `in` expansion already does.
- The bound is authored, not inferred: a step rule must exclude the cap
  cell with its own positive atom (or a literal row must claim it), and an
  unexcluded cap refuses at load with the cell named.
- A stepped tag must be initialised: the expanded `guard.all` atom reads
  the key, so an uninitialised counter surfaces as the existing
  `graph-owned-before-write` finding.
- Authored enum `domain` order acquires meaning (step order) where today
  it has none; reordering a domain that a rule steps changes the model.
- A guard bound and the declared `max` can drift apart: tightening `max`
  below an admitted cell's successor refuses at load (C2); widening it
  leaves new cells to the group's coverage proof (`graph-coverage-gap`).
  Both are existing refusals; no advisory is minted.
- A `step` ladder and a hand-unrolled ladder are two authorings of one
  behaviour but not of one dump: their row identities differ (suffix vs
  authored id). `0002:§round-trip-inverse-invariants` is not extended to
  the pair; the MVV compares plans and findings, not dumps.

### Risks and Mitigations

- **`unless` on the stepped tag over-admits cells** → the bound refusal
  fires and its detail names the `guard.all` complement; the authoring
  guide already prescribes that rewrite for a needed conjunction.
- **Wide `int` bounds multiply rows at load** (A7, spike-settled) → load is
  linear and cheap (100,001 rows in 0.96 s), so C1 takes no load-time
  ceiling. The cost lands on `lint`, whose enumeration is roughly quartic
  and whose bounds are 0013's caller-supplied analysis scope, unreachable
  from load by the package graph: a wide step domain yields a model that
  loads instantly and that lint cannot finish analysing, degrading the
  verdict (`complete=false`, reported as a finding) rather than refusing
  the model.
- **A rule with two atoms on the stepped tag in one block** (the 0011
  escape: the atom builder nests operators inside keys) → the admits filter
  is defined over EVERY positive atom on the tag, not "the" atom, so two
  atoms conjoin; the MVV fixture carries such a rule.
- **Suffix collision via a `#` domain member** (A3) → C1 refuses the step
  over such a domain at load.

### Failure Modes

- A step write on a `scalar`, `set`, `bool`, or unbounded `int` — load
  refusal, `malformed_tag_declaration`, detail names the admitted kinds.
  An `int` bounded across the full integer width refuses the same way: the
  domain is unenumerable (`internal/guard/declaration.go::intWidth` already
  reports that span as unusable), so this is an explicit refusal rather
  than an accidental fall-through to the zero-cell rule.
- A stepped value past `max` or past the last member at an admitted cell —
  load refusal (C2), detail names rule, tag, cell, value.
- A step rule admitting no cell — load refusal, the `0003:C6`
  disposition. A zero or float step — load refusal, C1's grammar arm. A
  step whose magnitude exceeds the domain width refuses at every admitted
  cell under C2, naming the first.
- A cap cell excluded by the step rule and claimed by no other row —
  `graph-coverage-gap` at lint, the existing finding.
- A stepped tag never initialised — `graph-owned-before-write` at lint.
- A consumer needing per-row identity from a flow verb — reads `graph` or
  `dump`, whose row surfaces publish `identity`. `flow next` and `flow
  resolve` name the authored rule on `rule` by design (C3).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] RDR 0002 is `Implemented` on `main` (it is); no open peer edit is
  scheduled — the joint check below records the peers consulted.

### Minimum Viable Validation

1. Author two in-repo fixtures of the Background's ladder over
   `internal/table/testdata/rdr-fixture.toml`'s `iter`-shaped `int` tag and
   an ordered `enum` tier: `ladder-literal.toml` (unrolled, as today) and
   `ladder-step.toml` (`attempt = { step = 1 }`, `tier = { step = 1 }`,
   each step rule carrying at least one row with two positive atoms on the
   stepped tag).
2. `intrastate lint` on both: identical finding sets, zero
   `graph-overlap`, zero `graph-coverage-gap` (A5).
3. `flow resolve` over every (state, outcome) cell of both: identical
   `writes`/`next`/`clear` (A6).
4. `ladder-step-unguarded.toml` (the `lt 5` atom removed): load refuses
   `malformed_tag_declaration` naming `retry`, `attempt`, cell `9`, value
   `10`; the enum sibling names `tier`, `large`. Two more refusal
   fixtures: `[initial] attempt = { step = 1 }` and a predicate literal
   `{ step = 1 }` both refuse — the table shape is write-block only.
5. `graph` export of `ladder-step.toml` decodes under the shipped
   `intrastate.graph/1` document type with no new member (C3).

End-state: the step ladder is the literal ladder to every consumer, and
the unguarded cap is a named load refusal.

### Phase 1: Admit the form

Intercept the inline-table shape on the write-block path, validate the
one admitted key and the admitted kinds, and record the step point on the
rule (C1's grammar half).

### Phase 2: Expand into cells

Extract `internal/guard/product.go::valueSatisfies` to `internal/resolve`
(the leaf `table`, `guard` and `graphlint` all reach — it needs only
`resolve.GuardAtom` and a `[]string` literal, so no `table` type is named
and the cycle `0012:A5` identifies does not arise), and repoint
`guard::valueSatisfies` and `graphlint::atomAdmitsValue` at it. The loader's
admits filter then calls the same function the runtime does, rather than a
fourth copy of it. Accepted cost: this widens the change beyond
`internal/table`, taken deliberately.

Add the step choice point to `expand`: admitted-cell evaluation over the
rule's positive atoms, the `guard.all` `eq` atom, the stepped literal, the
suffix element, the zero-cell and `#`-member refusals (C1's expansion
half), and the bound refusal through `conform` (C2).

### Phase 3: Authoring surface and proof

Document the form beside the write-block section of the authoring guide,
regenerate the help, land the MVV fixture pair and the refusal fixtures,
and assert C3's invariance on the export and the CLI previews.

## Validation

### Testing Strategy

Fixtures live in `internal/table/testdata/` beside `rdr-fixture.toml`
(whose `[tags.iter]` `int` with `min = 0, max = 9` and enum `domain`
declarations are the shapes the ladder pair is built on); the assertions
are CLI-level, on the JSON envelopes, in `internal/cli/lint_mvv_0030_test.go`
and `internal/cli/flow_mvv_0030_test.go` per the `<topic>_mvv_<rdr>_test.go`
precedent. Done = every scenario below green with no new finding code and no
new envelope member.

1. **Scenario**: `lint --as json` over `ladder-literal.toml` and
   `ladder-step.toml`, the same ladder authored unrolled and stepped, each
   step rule carrying a rule with two positive atoms on the stepped tag.
   **Expected**: identical finding sets as sets over `findings[]`
   (`model` normalized), spanning both the blocking codes and the advisory
   `graph-idempotent-write`; zero `graph-overlap` and zero
   `graph-coverage-gap` on both. Backs A5.
2. **Scenario**: `flow resolve --as json --plan-only` swept over every
   (state, outcome) cell of the pair, the state driven by rewriting the
   bound artifact file.
   **Expected**: `writes`, `next` and `clear` equal cell by cell. Backs A6.
3. **Scenario**: `ladder-step-unguarded.toml` — the `lt 5` atom removed —
   loaded.
   **Expected**: refusal under `malformed_tag_declaration` naming rule
   `retry`, tag `attempt`, cell `9`, value `10`; the enum sibling names
   `tier` and `large`. Backs C2, whose check is
   `internal/table/load.go::conform` → `::conformDomain` (the existing
   message names value and bound, and `renderWrites` already prefixes rule
   and key, so only the cell is new text).
4. **Scenario**: `[initial] attempt = { step = 1 }` and a predicate literal
   `{ step = 1 }`, each loaded.
   **Expected**: both refuse — the table shape is write-block only. Backs
   C1's interception arm.
5. **Scenario**: zero-cell, zero-step, float-step, and `#`-in-domain step
   models loaded; a step write on `bool`, `set`, `scalar`, an unbounded
   `int`, and an `int` bounded across the full integer width.
   **Expected**: each a load refusal under `malformed_tag_declaration`
   whose detail names the admitted form or kinds. Backs C1's grammar arm
   and, for the `#` case, the new guard A3 established is not redundant.
6. **Scenario**: `graph` export of `ladder-step.toml` decoded against the
   shipped `intrastate.graph/1` document type.
   **Expected**: decodes with no new member; `internal/table::Categories()`,
   the CLIError codes and the finding codes each gain none. Backs C3.
7. **Scenario**: `flow resolve` selecting a NON-FIRST expanded row of a
   stepped rule that authors an `emit` block —
   `TestReq34_TheRuleIDJoinIsSoundOverAnExpandingRule`'s shape
   (`internal/cli/decision_table_0010_test.go`) extended from an `in`
   expansion to a `step` one.
   **Expected**: the payload's `emit` and `dispositions` are the authored
   block, and `rule` is the authored id — the first-match `rowByID` join
   stays sound because a stepped rule's rows share one `Emit`. Pins the
   `0010:C3` inheritance the audit row records.

### Performance Expectations

Measured, not estimated — `evidence/spikes/a7-row-growth.md`, the A7 spike,
over tables of N literal rows (what the expansion mints).

Load and normalize are **linear** in the expanded row count: 0.96 s and
584 MB resident at 100,001 rows, 5.08 s at 500,000, ~10 us and ~5.8 KB per
row throughout. So the declaration bounds the cost and no load-time ceiling
is owed at the ladder sizes this record serves (tens of rows) or four orders
of magnitude past them.

Lint is the opposite curve and is **0013's analysis scope, not this
record's**: it is roughly quartic — 7.5 s at 100 rows, 116 s at 200, beyond
600 s at 1,000 — so a step write over a wide domain yields a model that
loads instantly and that `lint` cannot finish analysing. Exceeding a lint
bound degrades the verdict (`reach.go` sets `complete=false`, surfaced as a
finding) rather than refusing the model, which is why the bound cannot be
answered at load: a load-time ceiling would refuse models the loader handles
in under a second for a cost the loader never pays. Were a ceiling ever
wanted it would be memory-motivated and first bite near 1,000,000 rows
(~6 GB extrapolated), an order of magnitude past the A7 declaration.

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

Re-validate the **Profile** Metadata field: re-run
rdr-write's `--outcome profile` with the clause's own
dispositions and confirm the value Resolve wrote is
what it emits (a stop is not a match).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. (The row already
excludes `Transient`-marked contracts.) Also confirm form:
value + one clause naming the contract(s) and its two
dispositions; strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- RDR 0002 `0002:C4` (write replaces), `0002:C13` (match-block expansion),
  `0002:C17` (`#` ban scope), `0002:C22`/`0002:C24` (declaration keys and
  load categories); RDR 0003 `0003:C3`/`0003:C6` (finite domains, load-time
  refusal of an unsatisfiable atom); RDR 0012 `0012:C2` (typed comparison);
  RDR 0021 `0021:C2` (export vocabulary, authored `domain`).
- `internal/table/load.go::valueMembers`, `::conform`, `::conformDomain`,
  `::tagDecl`; `internal/table/normalize.go::renderWrites`, `::expand`;
  `internal/table/model.go::Identity`; `internal/graphlint/reach.go::
  successor`, `engine.go::Fingerprint`, `groups.go::checkIdempotentWrites`;
  `internal/cli/graph_document.go::graphRows`;
  `internal/guard/declaration.go::AssignmentCount`.
- `docs/model-authoring.md` §Discriminate with guard atoms, §Guard and
  match atoms, §Rules advance the state.
- Prior art: `../state-machines/contrast/poc-rdr-ragel/REVIEW.md`
  §Faithfulness; `../state-machines/study/sismic/docs/code.rst`;
  `../state-machines/MODEL-transition.md` §3a; scxmlcc user manual §Assign.
  Record: `0030-computed-write-value-grammar/evidence/research/prior-art.md`.
- Related issue: intrastate#jwgr.
