# Recommendation 0024: Declared emit vocabulary

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
- **Profile**: foundational — one contract (the emit-vocabulary
  declaration and what it proves, stated as C1–C4 by surface:
  grammar, load proof, normalized carry, resolve envelope) that
  locks a grammar, spans `internal/table` and `internal/cli`, and
  produces a declaration carry peer RDRs 0021/0023 consume.
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
- **Priority**: High
- **Related Issues**: kata `intrastate#vt9n` (1638); kata
  `intrastate#rg0e` (1605, folded into this record by
  seed-triage — the stop/route disposition facet); kata `srz2`
  (RDR 0023 — a projection would naturally select the declared
  emit keys)
- **Predecessors**: 0002-transition-table-as-reviewable-data,
  0005-skill-integration-cli-contract,
  0006-graph-lint-authority-and-guarantees,
  0010-stateless-decision-tables
- **Seam Lineage**: no prior accretion

## Problem Statement

A model author's bargain with lint is that the table is the single
routing authority *because lint proves it* — every routing cell
claimed exactly once, a guarantee prose cannot give — and the
motivating consumer's skills therefore execute `emit.next` verbatim
instead of re-deriving routing. The guard half keeps that bargain; the
answer half does not: emit keys and values are arbitrary strings
nobody validates. Proven at HEAD (`bc9f2a0`): a two-rule table whose
first rule emits a misspelled command value and whose second rule
typos the key itself (`nxet`) lints with exit 0 and zero findings. The
table's routing is proven; its answers are unreviewed — a model can
route perfectly to a command that does not exist, and a verbatim
executor executes it. A consumer-side golden table covers only the
rows in-flight work exercises, so a rare row's typo ships and fires at
a close-out.

The folded facet (kata `rg0e`): some rows deliberately answer with
`stopped:` tokens — a judgment a fact cannot make — and the
stop-vs-route distinction is load-bearing: a stop printed as a command
to run is a skipped check reading as a passed one. Today nothing in
the envelope marks a stop (resolve exits 0 for both; `0010:C3` fixes
emit values as uninterpreted, compared by byte equality), so the
distinction lives in a string-prefix naming discipline every consumer
re-implements and no lint checks.

The decision, made once: does a model declare its emit vocabulary, and
what does a declaration prove? The candidate weighed: an optional
per-key `[emit.<key>]` declaration table mirroring the tag grammar's
declare-then-prove move, with domain members carrying a disposition
that surfaces on the resolve payload (the Proposed Solution is its
adopted form; C1–C4 are authoritative). Predicate: with any
declaration present, an undeclared emit key and an out-of-domain
value each refuse at load and surface as blocking `intrastate lint`
findings; a model with zero declarations lints exactly as today
(opt-in, no corpus breakage). The RDR must weigh the folded facet's
cheaper competing shape — a documented reserved prefix enforced by
lint — and decide the contract once across grammar
(`internal/table`), proof (the load pipeline), and surfacing
(`internal/cli`), keeping intrastate generic: domains and dispositions
are model-authored.

## Critical Assumptions

- **A1 The top-level `[emit]` TOML key is free: the source schema
  declares no `emit` field, so today's `decodeStrict` refuses it as an
  unknown schema field and adding it collides with nothing.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/source.go::sourceDoc` declares
    exactly `Outcomes, Terminal, Model, Initial, Tags, Read, Write,
    Gate, Context, Rule, Dump` (toml tags `outcomes, terminal, model,
    initial, tags, read, write, gate, context, rule, dump`) — **no
    `emit` field**, so the top-level key is free. The refusal leg is
    real and carries the claimed slug:
    `internal/table/source.go::decodeStrict` sets
    `dec.DisallowUnknownFields()` and maps a `*toml.StrictMissingError`
    to `fail(CatUnknownSchemaField, …)`, where
    `internal/table/category.go` spells `CatUnknownSchemaField` as
    `"unknown_schema_field"`. It is called unconditionally from
    `internal/table/load.go`'s pass-2 strict decode, so an older binary
    handed a declared model refuses loudly rather than ignoring the
    declaration.
  - **If wrong**: the declaration table needs another authoring locus
    and C1's grammar moves; surfaces as a decode error on the first
    declared model.
- **A2 Appending `dispositions` after `emit` in
  `internal/cli/flow_resolve.go::resolvePayload` breaks no payload
  CONSUMER (every reader parses by key, none positionally), and its
  cost is a bounded, enumerable test-corpus diff — 31 inline
  whole-payload assertion sites, 3 `resolveGolden` byte-identity
  goldens, the `shippedFixtureGolden` sweep, and one peer REQ test
  that pins the struct's field count and wire key order.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: **Consumer leg — clean.** Every Go-side reader
    parses by key (`data["emit"]`, `map[string]any`); `wireKeyOrder`
    and `keysOf` in `internal/cli/decision_table_0010_test.go` are
    test helpers asserting order, not production parsers, and the
    consumer skills under the RDR engine's `skills/` do not parse the
    payload positionally. `encoding/json` emits struct fields in
    declaration order — the property `resolvePayload`'s own comment
    already relies on ("the declaration order IS the JSON key order
    Go's encoder emits") — so inserting after `Emit` does place
    `dispositions` there.
    **Cost leg — enumerated, and larger than a fixture sweep.** The
    binding item is
    `internal/cli/decision_table_0010_test.go::TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire`,
    which asserts BOTH the exact 13-key wire order (`slices.Equal`)
    and `reflect.TypeOf(resolvePayload{}).NumField() != 14` — it
    encodes `0010:C4`'s normative position clause as a REQ test, so
    editing it amends a peer RDR's spec verification, not a fixture.
    C4 takes the struct from fourteen fields to fifteen and the wire
    from thirteen keys to fourteen. Beside it sit the inline
    whole-payload assertions and byte-identity goldens in
    `internal/cli/flow_demand_0011_test.go` (RDR 0011's guard against
    exactly this class of change). **The full enumerated edit set —
    counts, files, and what each asserts — is the licensed diff in
    Implementation Plan Phase 3**, named there so it is authorized
    rather than discovered mid-change.
    **Citation correction.** `0005:C1` does not state a general
    append-only discipline for the SUCCESS payload — it extends the
    CLIError *failure* envelope by "exactly one omitempty structured
    field, `findings`". `0010:C4` is the operative payload precedent
    (it appended `emit` after `Gates` for this same reason: "the repo
    asserts payload JSON inline in Go tests, so an unfixed position is
    an unlicensed diff"). `0010:A4` is a real precedent for the
    *shape* of this obligation — its census "counted TOML only, and
    that was the gap", missing a Go slice literal — which is why this
    Evidence enumerates the Go-side assertions rather than only the
    inline JSON.
  - **If wrong**: the `dispositions` field breaks a consumer parsing
    positionally or an inline test corpus larger than one change can
    carry; surfaces as red payload-equality tests.
- **A3 No checked-in model, fixture, or test source anywhere in the
  repo authors an `[emit.<key>]` declaration, so every existing emit
  author takes C2's zero-declaration leg and the whole corpus lints
  exactly as today.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-corpus-sweep.md`.
    `rg --glob '*.toml' '^\s*\[emit'` over the repo returns **zero
    hits** — no model, no fixture, no evidence artifact declares an
    emit vocabulary, so the opt-out leg covers the entire corpus.
    Baseline confirmed green at HEAD: `go build ./...` clean,
    `go test ./...` all packages `ok`, and all four `models/*.toml`
    lint exit 0 (two carry pre-existing non-blocking
    `graph-coverage-closed-by-escape` advisories, unrelated to emit).
    **Enumeration corrected**: `[rule.emit]` is authored in more
    places than `models/` — `models/examples/pricing-decision-table.toml`
    is the only *model* (4 blocks), but ten `_test.go` files under
    `internal/{table,cli,graphlint}` author `[rule.emit]` as inline
    TOML string literals, and one 0010 evidence artifact authors one
    block. The Go test sources are load-bearing for this assumption
    because they feed the loader directly; all of them likewise
    declare zero `[emit.<key>]`, so the opt-out predicate holds over
    every one. The post-change leg (full suite green WITH the loader
    change and no model edited) is unrunnable at Draft — that code
    does not exist at HEAD — and is carried to implementation as
    Testing Strategy scenario 3, not claimed here.
  - **If wrong**: strictness fires on an existing model and the
    opt-in predicate is broken; surfaces as a load refusal in an
    untouched fixture.
- **A4 Kind conformance for an emit value is a value-level lexical
  check (enum membership, bool/int literal form) implementable in the
  load pipeline without importing `internal/guard`'s predicate-atom
  machinery.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/guard` exports **no value-level
    conformance checker**. Its closest export,
    `guard.Conforms(m *table.Model, v View) error`, checks only
    required-key presence and single-valued cardinality, skips
    `ProvenanceRecognized`, and performs no kind or domain-membership
    check; every other export is keyed to predicate atoms, assignment
    sets, or product cardinality. So the assumption's operative claim
    — implementable **without importing `internal/guard`** — holds.
    The reusable checker is in this package, not guard:
    `internal/table/load.go::conform` / `::conformKind` /
    `::conformDomain` are unexported, receiver-free functions taking
    `(TagDecl, string)`, and `conformDomain`'s enum arm is exactly the
    membership test C2 needs (`slices.Contains(decl.Domain, member)`).
    They are keyed to `TagDecl`, so an emit declaration reuses them
    through a small adapter or a narrower helper extracted from the
    two; the `set` arm simply never fires for emit (`0010:C3`: an emit
    value is one authored string).
    **Evidence-pointer correction**: the write-block conformance
    CHECKER is at `internal/table/load.go::conform`, not in
    `normalize.go` — `internal/table/normalize.go` holds only the call
    site (`conform(decl, "eq", members)`).
  - **If wrong**: the check either drags guard's atom semantics into
    emit (the wall `0010:C3` builds) or duplicates kind rules that can
    drift from RDR 0003's; surfaces as divergent refusals for the same
    literal.
- **A5 The motivating consumer's emit answers are enumerable at
  authoring time — routing commands are fixed strings and judgment
  cells answer with fixed `stopped:<code>` tokens — so an enum domain
  with a disposition partition can actually be authored for each of
  its routing keys.**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: `0010:A6` (Status `Verified`) states it verbatim:
    the stage verb and the lens "are fixed strings known when the
    table is authored", the record number "is the argument the caller
    already holds", and "Judgment cells route to a stop-packet row,
    again a fixed `stopped:<code>` string" — per-invocation data is
    supplied by the caller, never the table.
    **Empirical confirmation over the live consumer models.**
    `rdr-status.toml` carries 54 `[rule.emit]` blocks whose `next` key
    holds **22 distinct values, every one a fixed literal**, with no
    interpolation and no per-invocation data: 16 route-like
    (`none`, `/rdr-refine`, `/rdr-resolve`, `/rdr-prelock <lens>`, …)
    and 6 `stopped:`-like (`stopped:no-profile`,
    `stopped:determinacy-trigger-unjudged`, …). The partition is
    clean, so a disposition-partitioned enum domain is authorable as
    the model stands today.
    **Generalized from the Draft's "its `next` key"**: the consumer
    has TWO routing keys under two names — `rdr-write.toml` (30
    `[rule.emit]` blocks) routes on **`op`**, not `next`, with 10
    distinct values partitioning the same way (`claim`, `lock`,
    `demote`, `readme-add`, … / `stopped:not-lockable`,
    `stopped:no-gate-written`, …). C1's per-key `[emit.<key>]` grammar
    accommodates this by construction; a key-level or single-key
    design would not have. `rdr-write.toml`'s `edit` key holds a
    multi-line shell script — an unambiguously open-valued key that
    takes `kind = "scalar"`, which is the per-key escape C1 provides
    and a live instance of why it is needed.
  - **If wrong**: the consumer cannot close its domain and declares
    `scalar`, which admits the misspelled value again; surfaces as the
    seed defect persisting under a declaration.
- **A6 DMN decision tables carry an allowed-values list on output
  clauses that conformant tooling checks output entries against — the
  external alignment claim for declare-then-prove over answers.**
  - **Status**: Pending — **not load-bearing**; carried unverified by
    decision, not by omission (see plan below).
  - **Method**: Prior Art
  - **Evidence**: **Searched and not found**
    (`evidence/research/resolve-prior-art.md`). No DMN or BPM-suite
    material exists in any corpus available to this project — DevRef,
    StateMachineLit, StateMachineRes, and PapersFast were each
    queried, and the sibling `../state-machines` repo returns zero
    word-boundary hits for `DMN` / `Decision Model and Notation` /
    `outputValues`. The two legs separate: leg 1 (the standard
    declares an allowed output-values list on the output clause) has
    no section-anchored source; leg 2 — that conformant tooling
    *checks* authored output entries against it, which is the half
    that would actually corroborate declare-then-prove over answers —
    has **zero** support. Camunda's output-clause docs describe
    `typeRef`, a type constraint on evaluated output, not an
    allowed-values enumeration; rejected as non-analogous.
    **Verification plan if ever needed**: read OMG DMN 1.3 §8 directly
    from the spec PDF (publisher sites block automated fetch, so this
    is a manual acquisition into a corpus first). Not scheduled — the
    claim is carried as an unverified aside, and this Draft does not
    rest on it: the alignment argument runs on the in-repo mirror
    (`0002:C22`) and the two opened peer citations (ms-conductor
    `type: terminate`, scxmlcc `<final>`), which are Verified prior
    art for the structural-declaration move.
  - **If wrong**: an alignment citation drops; the approach stands on
    the tag-grammar mirror and the opened peer citations.

## Proposed Solution

### Approach

An **optional per-key `[emit.<key>]` declaration table** that mirrors
the tag grammar's declare-then-prove move, whose checks land in the
**load pipeline** (not graphlint), and whose domain members may carry
a **model-authored disposition** that `flow resolve` surfaces as one
appended payload field.

A declaration answers the contract question the Problem Statement
poses — *what does a declaration prove?* — as: **a declaration
constrains authoring, never evaluation.** With any declaration present, the
loader proves every authored `[rule.emit]` key is declared and every
authored value conforms to its key's declared kind/domain, before
candidate rows are yielded. Evaluation is untouched: emit values stay
uninterpreted and byte-compared (`0010:C3`), the kernel never sees
emit, and intrastate never interprets a disposition token — domains
and dispositions are model-authored vocabulary, carried through and
surfaced verbatim. A model with zero declarations loads and lints
exactly as today; its one observable delta is the resolve payload's
appended `dispositions: {}` (C4) — the same class of additive,
fixture-updated change `0010:C4` shipped for `emit` itself, never a
behavioral one.

The enforcement point is the load pipeline, not lint: both new checks
are single-rule — decidable from one rule plus the model's
declarations — and `0002:C24`'s load/lint arity split assigns
single-rule checks to load ("Cross-row findings … carry RDR 0006's
lint categories, not these"); a parallel lint-side single-rule tier
would contradict the split the repo is built on. The Problem
Statement's predicate holds verbatim: `intrastate lint` loads first
and maps a load refusal to blocking findings (`internal/cli/lint.go`
builds them via its `loadFindings` arm), and `flow resolve` maps
every load category to `flow-model-invalid` (`0005:C1`). No graphlint
code changes; the advisory tier stays closed (`0006:C17`) and the
blocking tier's "at least" floor (`0006:C3`) is not touched.

### Technical Design

Data flow: `[emit.<key>]` declarations decode with the source
document, load beside the tag declarations, are proven against every
rule's authored emit block during the source-to-candidate-rows
pipeline, and are carried losslessly onto the normalized model —
mirroring `0002:C22`'s carry-through clause — so downstream readers
(the resolve payload join today; later lint finding classes or
exports tomorrow) read the same declaration the author wrote. This
RDR adds exactly one downstream reader: the `dispositions` join in
`internal/cli/flow_resolve.go::resolvePayload`.

#### Normative Contracts

These four blocks are facets of the one contract this RDR owns — the
emit-vocabulary declaration and what it proves — stated separately by
surface (grammar, proof, carry, envelope), not four independent seams.

**C1**

```normative
GRAMMAR. A model MAY declare its emit vocabulary in a top-level
`[emit]` TOML table, one sub-table per emit key: `[emit.<key>]`.
Each declaration carries:

- `kind` (required): one of `enum | bool | int | scalar` — RDR 0003's
  token spellings reused verbatim (`internal/table/model.go::
  declaredKinds` is the adjacent vocabulary; `set` is excluded because
  an emit value is one authored string, `0010:C3`).
- a domain, `enum` only, in exactly one of two forms:
  - `domain = [ ... ]` — a flat member array, no dispositions; or
  - `[emit.<key>.domain]` — a sub-table whose keys are MODEL-AUTHORED
    disposition tokens (e.g. `route`, `stop`, `terminal` — intrastate
    fixes no vocabulary) and whose values are member arrays. The
    key's domain is the union; each member carries the one disposition
    it is listed under.
- `bool` fixes the implicit domain `true | false`; `int` constrains
  the value to a base-10 integer literal; `scalar` is the
  declared-but-unvalidated escape hatch — the key is admitted, the
  value unconstrained. None of the three takes a `domain`, and
  dispositions attach only to declared enum members.

Every kind check is a LEXICAL check on the authored string: no value
is parsed into a typed representation, canonicalized, or converted
anywhere downstream — evaluation and the payload carry the authored
bytes (`0010:C3`). Whether a non-canonical `int` literal (`03`) is
admitted is a Resolve-level detail of the lexical rule, not a
matching question — emit values are answers and are never compared
against the domain at evaluation.

Refused as `malformed_emit_declaration`: an unknown `kind` token; an
`enum` with no domain, an empty domain, an empty-string member, or a
duplicate member (duplicates across disposition lists included — one
member, one disposition); a `domain` on a non-enum kind; a
disposition token that is the empty string. Emit declarations are NOT
tag declarations: no provenance, no accessor reference, no
`min/max/elements/single_valued/required`, and an emit key remains
barred from match, guard, write, and accessor use (`0010:C3`).
```

**C2**

```normative
PROOF (opt-in, whole-model). With ZERO `[emit.*]` declarations the
load pipeline is byte-for-byte today's: no new refusal is reachable.
With ONE OR MORE declarations present, the source-to-candidate-rows
pipeline (`0002:C24`'s "load") MUST refuse, before yielding rows:

- `unknown_emit_key` — any `[rule.emit]` key of any rule, ordinary or
  escape, not declared under `[emit]`. Strictness is whole-model, not
  per-key, because a typo'd KEY is indistinguishable from an
  intentionally undeclared one — per-key checking cannot catch the
  seed's `nxet` defect.
- `emit_value_out_of_domain` — an authored value that is not a member
  of its key's declared enum domain, not a `bool` token, or not an
  `int` literal, per C1's kinds (`scalar` values are never refused).

These two categories and C1's `malformed_emit_declaration` join
`0002:C24`'s data-level set (that list is "at minimum", and
`reserved_tag_key` is the append precedent) and therefore map to
`flow-model-invalid` under `0005:C1`, and to a blocking finding under
`intrastate lint` via its load-refusal arm
(`internal/cli/lint.go`). All three categories are
refusals — nonzero exit, never advisory, under every surface that
loads the model.

**Reporting is FAIL-FAST, one refusal per run — inherited from the
load tier, not decided here.** `0002:C24` routes these checks to load
by arity, and the load tier's cardinality is already contracted:
`0008:C4` fixes it — "the predicate reports **the first breach it
finds and returns a single non-nil error**; it does not aggregate …
a producer holding a programmer mistake is not owed an exhaustive
list" — and `internal/table/load.go`'s package contract states the
same for the pipeline as a whole ("the refusal is singular: load is
fail-fast and returns one categorized error, never a list"). So a
model carrying an emit defect AND a coexisting structural one
surfaces whichever the pipeline reaches first; the author fixes and
re-runs, exactly as for every other load category today.

**Why fail-fast is the correct tier behavior, not a limitation.** The
project accumulates findings precisely where a TOTAL ORDER over them
exists, and only there — the term is used in its exact sense, and each
instance names the comparator that provides it: graph lint reports
every defect in one pass (`0003:C19`, restated `0006:C16`) because
`internal/graphlint/engine.go::sortFindings` orders findings by their
identity tuple, making the emitted set independent of iteration order;
`0009`'s escape-shape report aggregates because
`internal/resolve/resolve.go::compareRefs` orders `RowRef`
lexicographically on `(RuleID, SourceLocator)`. Load
has no such comparator — `internal/table/load.go` states that "the order in
which independent defects are checked is **deliberately
unspecified**" — so accumulating there would publish a
non-deterministic set, which is the property that unspecified order
exists to keep unobservable. Where two refusals coexist the house
answer is PRECEDENCE, not aggregation (`0018:C1` makes the entry-check
order an observable contract; `0018:ALT1` rejected `errors.Join` on
the ground that "aggregation serves user-facing declarative
validation; this is the programmer-mistake path"). This RDR fixes no
precedence among the three emit categories and the existing ones, and
takes none: order stays unspecified.

`0005:C1`'s "one findings[] entry per category hit" is satisfied
vacuously, because no load path yields more than one hit to map — this
RDR introduces no accumulation and takes no `Overrides` entry against
it. **What this RDR does NOT claim** is the lint tier's
report-everything guarantee: an emit refusal MAY mask a coexisting
structural one, and vice versa. Multi-defect load reporting would
require inventing the total order `load.go` declines to fix; it is a
pipeline-wide change out of scope here, left to a later RDR that owns
it.
Proving the AUTHORED form is proving the
executed form: `0010:C3` fixes that normalization carries the emit
block through `expand` unmutated ("Emit is never mutated after
normalization") — no emit value exists post-normalization that load
did not check. Nothing is checked at resolve time, and an emit
value's evaluation semantics — uninterpreted, byte-compared — are
unchanged (`0010:C3`). A declared key no rule emits, and a declared
member no rule authors, are NOT findings of any tier (authoring
headroom; the advisory tier is closed, `0006:C17`) — a
never-emitted stop token is a liveness question over the graph, left
to a later consumer of C3's carry, not a load check.
```

**C3**

```normative
CARRY. Every declared field — key, kind, domain members, and each
member's disposition — is carried losslessly through normalization
onto the normalized model, the same clause `0002:C22` states for tag
declarations ("lint (RDR 0006) and the guard proof (RDR 0003) read
the same declaration the author wrote"). The carried form is a new
model-level type, deliberately NOT `TagDecl` (the mirror of
`EmitValue` not being `TagValue`, `0010:C3`). The kernel
(`internal/resolve`) continues to carry no declarations and no emit.
This RDR adds exactly one reader of the carried declarations — C4's
payload join — and `internal/graphlint` reads none of it here; the
carry exists so a later RDR can add finding classes or exports over
the declared vocabulary without reopening the grammar.
```

**C4**

```normative
ENVELOPE. The `flow resolve` success payload gains `dispositions`: a
JSON object mapping emit key → the disposition token the declaration
assigns the selected row's authored value, keys in byte order,
present as `{}` — never `null`, never omitted — when no selected
value carries one (undeclared model, non-enum kind, flat-array
domain, or empty emit block alike). It is inserted at one fixed
position in `internal/cli/flow_resolve.go::resolvePayload`,
immediately after `emit` — the same additive insertion `0010:C4`
used (`emit` after `Gates`), for the same reason: struct declaration
order is the emitted field order and the repo asserts payload JSON
inline. The token is surfaced verbatim;
intrastate never interprets it. Text mode renders through the generic
payload renderer as `dispositions.<key>: <token>` / `dispositions:
(none)` with no per-verb special case (the `0010:C4` clause). A plan
rescued by an escape row joins `dispositions` from that escape row's
OWN authored values — the same single `Plan.RuleID` join path as
`emit` (`0010:C4`), so no defaulted or merged value can reach the
payload unjoined: there is no default-row or merge path for emit in
this model family. `flow next` carries no `dispositions`, for
`0010:C4`'s reason: the answer is what `resolve` selects.
```

#### Load-Bearing Decisions

- **Identity** — an emit declaration is identified by its emit key,
  byte-exact (the same identity emit keys already have, `0010:C3`);
  one declaration per key, duplicates refused by the TOML decoder as
  the tag table's are. A domain member's identity is its byte-exact
  string; one disposition per member.
- **Naming** — the table is `[emit]`, mirroring `[tags]` /
  `[rule.emit]` (rejected: `[emits]`, `[vocabulary]`,
  `[declare.emit]` — the grammar sits beside the block it constrains
  and the tag grammar it mirrors). The new load categories are
  `malformed_emit_declaration`, `unknown_emit_key`, and
  `emit_value_out_of_domain`, following `0002:C24`'s
  `malformed_*`/`unknown_*` house scheme. The payload field is
  `dispositions` (rejected: `emit_dispositions` — it sits adjacent to
  `emit`; `kinds` — the field carries dispositions, not kinds).
- **Selection / predicate** — dispositions are per **domain member**,
  not per key: the folded facet's whole point is that one key's
  values split between routes and stops, so a key-level disposition
  cannot express the model that motivates the field. When the
  selected row authors a value for a declared enum key, the
  disposition surfaced is the one the declaration lists that member
  under — exactly one exists by C1's duplicate refusal.

#### Illustrative Code

Illustrative only — shapes, not fixtures.

```toml
[emit.next]
kind = "enum"

[emit.next.domain]
route = ["/rdr-propose", "/rdr-refine"]
stop  = ["stopped:joint-decision", "stopped:propose-order"]

[emit.dpa]
kind = "enum"
domain = ["required", "waived"]   # flat form: no dispositions
```

```json
{"rule":"joint-fire","gates":[],
 "emit":{"next":"stopped:joint-decision"},
 "dispositions":{"next":"stop"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Declaration authoring locus (`[emit]` top-level table) | This RDR | Introduced | C1; A1 verifies the key is free |
| Kind vocabulary and value-conformance rules | Predecessor (RDR 0003 via `0002:C22`) | Available | C1 reuses the token spellings; A4 decides reuse vs local check |
| Load-category refusal + findings mapping | Predecessor (`0002:C24`, `0005:C1`) | Available | C2 appends three categories to an "at minimum" set |
| Lint surfacing of load refusals as blocking findings | Predecessor (RDR 0006 CLI arm) | Available | C2; no graphlint change |
| Append-only resolve payload | Predecessor (`0005:C1`, precedent `0010:C4`) | Available | C4 inserts `dispositions` after `emit` |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Declaration parse/load | `internal/table/load.go::loadTags` + `sourceTagDecl` | Tag-specific: provenance, accessor ties, reserved-key rules | Extend the pattern, not the type — a parallel `loadEmit` with its own decl type | C1's "NOT tag declarations" clause |
| Value conformance | `internal/table/load.go::conform`/`conformKind`/`conformDomain` (`normalize.go` is the call site) | Keyed to `TagDecl` and tag kinds incl. `set` | REUSE — verified separable (A4): unexported, receiver-free, `(TagDecl, string)`; adapt via a narrow helper, no `internal/guard` import | C2 |
| Emit carry-through | `internal/table/model.go::EmitValue`, `normalize.go::emitSequence` | Carries authored pairs only, no declarations | Extend model with a declaration carrier | C3 |
| Payload join | `internal/cli/flow_resolve.go::resolvePayload` + `emitMap` | Fixed field order asserted inline in tests | Extend by append (A2) | C4 |

### Decision Rationale

Scored QOC matrix (5 best; criteria weighted equally, the choice
falls out of the correctness and blast-radius rows):

| Criterion | A: per-key declaration (chosen) | B: reserved stop-prefix lint | C: declare via tag grammar | D: validate at resolve time |
| --- | --- | --- | --- | --- |
| Correctness fit (catches misspelled value AND typo'd key; marks stops) | 5 — whole-model strictness catches both seed defects; member dispositions mark stops | 2 — marks stops only; both seed defects still lint clean | 4 — same checks, but drags tag semantics onto emit | 2 — un-exercised rows stay unproven; the rare-row typo still fires at a close-out |
| Prior-art alignment | 5 — mirrors `0002:C22`; peers declare stops structurally (SCXML `<final>`, ms-conductor `type: terminate`) | 2 — the prefix convention is the discipline the facet complains about | 3 — reuses the grammar but against its own wall | 2 — no peer validates answers only at evaluation |
| Reversibility / opt-in | 5 — zero declarations = today, byte-for-byte | 4 — opt-in flag or reserved token | 2 — reopens a locked contract to later back out | 3 — runtime behavior change to back out |
| Blast radius on locked contracts | 4 — additive: new grammar, three appended categories, one appended field; `0010:C3`/`C4` untouched | 5 — smallest | 1 — contradicts `0010:C3` ("emit keys are NOT tag keys"), an amendment path that does not exist | 2 — strains `0010:C3`'s uninterpreted-evaluation clause |
| Cost | 3 — grammar + checks + carry + field | 5 — trivial | 3 — comparable to A | 3 — comparable, plus evaluation-path risk |
| **Total** | **22** | **18** | **13** | **12** |

The deciding rows are correctness and blast radius: B is the only
cheaper option and cannot catch either proven defect — it ships the
convention the declaration exists to replace; C buys nothing over A
while amending a Final contract; D re-creates the exact failure the
seed documents (a defect that fires when the rare row is finally
exercised). A is the only approach that closes both defects, keeps
every locked contract intact, and costs nothing to models that do not
opt in. The load-not-lint enforcement point is the Approach's
`0002:C24` arity argument; the acceptance predicate holds either way.

**House precedent (verified at Resolve,
`evidence/research/resolve-precedent.md`).** The chosen shape is the
established pattern, named as such by a peer: `0020` decided the
sibling question for `--tag` keys and its rationale states "the house
already adjudicated this shape once … **carrier-by-default,
declare-to-tighten is the established pattern**", with `0020:C1`
closing "declaring a key later is the **opt-in tightening**" and
`0020:BR1` naming "**the 0024 shape**" as compatible layered work.
The two records govern opposite sides of one seam and do not collide:
0020 admits an undeclared CALLER key at runtime *because* the
model-side closed world stays shut ("a rule atom, accessor key, write
target, clear target, or `[initial]` assignment naming an undeclared
tag still refuses at load"), and emit values are model-authored, so
they sit on the closed side. The loader census bears this out — tag
keys, tag kinds, operators, the outcome alphabet, dump columns,
reserved keys, and schema fields all refuse unknown members at load;
`[rule.emit]` is the **sole** admitted-unconditionally vocabulary.
Closing it restores the rule the author-facing doc already states:
"a misspelled facet is a stable refusal rather than a silent no-op"
(`docs/model-authoring.md`), whose sibling clause — "declaring the
class rather than deriving it is what makes the refusal possible" —
is this RDR's argument in the project's own words. `ALT3` is
independently foreclosed by `0006:C1`: coupling the resolver to a
design-time proof "would put a design-time proof on the runtime path,
which this RDR's authority split rejects".

The hardened premortem (`evidence/propose-premortem/critic.md`,
17-finding ledger) returned PASS with mitigations; the recommendation
survives hardened. Its accepted findings are live in C1–C4, the
Risks and Mitigations, and the Failure Modes.

Premortem: hardened (hardened)

Ground-sweep: clean (23 anchors)

Joint-check: fired → 0023 (home: JDR 0002 §D1; the verb's instance
assignments ride `cli/0023:C2`). The fire: 0023's proposal shares the modify-anchor
`internal/cli/flow_resolve.go::resolvePayload` — this RDR appends
`dispositions` (C4), 0023 projects the echo group off the same
payload under an opt-in `--plan-only`. Homed and compositional, not a
contradiction: under `JDR 0002 §D1`'s partition doctrine
`dispositions` is a never-projected PLAN-group field (the instance
assignment rides `cli/0023:C2`) because C4's join reads plan-group
inputs only (the `[emit]` declaration and the selected row's authored
value, by `Plan.RuleID`), so C4's append-last and never-omitted
clauses hold under the projection in either landing order and neither
contract moves. Peer sweep: no other peer
shares a modify-anchor at symbol level — 0012 touches
`internal/cli/flow_resolve.go::guardSeam` and `internal/guard/lint.go`
where this RDR touches `resolvePayload` and `internal/cli/lint.go`
(same or like-named files, disjoint symbols); 0016 touches
`internal/table/load.go::checkAccessorBindings` where this RDR adds a
sibling emit-declaration load step. No peer names this RDR's authored
literals (`malformed_emit_declaration`, `unknown_emit_key`,
`emit_value_out_of_domain`, `[emit.<key>]`, payload `dispositions`);
the word "dispositions" in 0017/0018 is prose about finding/escape
dispositions, a different referent. Roster peers 0021/0023 name
`emit` (0010's payload surface) as what their export/projection would
carry — the consumer side of this seam; this RDR is the owner side
and C3's lossless carry is the interface they conform to. Absence
arm: n/a — this proposal converts acceptances into refusals, never a
refusal into an acceptance; the predecessor reliance on emit being
"undeclared" (`0010:C3`) is the widening `0010:A6` explicitly
deferred to a new RDR — lineage, not a fire. Bridge sub-check: n/a —
no surface here is scheduled for deletion by a sibling plan and this
plan retires none.

## Alternatives Considered

### Alternative 1: Reserved stop-prefix enforced by lint

**Description**: The folded facet's cheaper shape — document one
reserved value prefix (e.g. `stopped:`) or let the model declare a
prefix string, and have lint (or load) mark any emit value carrying
it as a stop; everything else routes. No key admission, no domains.

**Pros**:

- Smallest possible change; no new grammar, one check.
- Directly encodes the stop-vs-route distinction the facet needs.

**Cons**:

- Neither proven defect is caught: the misspelled command value and
  the typo'd `nxet` key both still lint clean — the main kata's
  bargain stays broken.
- It standardizes exactly the string-prefix naming discipline the
  Problem Statement calls out ("every consumer re-implements and no
  lint checks") instead of replacing it.
- A binary stop/route split is baked in; the peer instance read shows
  peers carrying richer declared end-state semantics (ms-conductor's
  `status: success|failed` on `type: terminate`).

**Reason for rejection**: solves only the folded facet, at the cost
of shipping the convention as contract; the QOC correctness row is a
2.

### Alternative 2: Declare emit keys through the tag grammar

**Description**: Add a provenance (or sibling arm) to `[tags.<key>]`
so emit keys reuse `TagDecl`, the existing kind/domain load checks,
and the declaration carry that already exists.

**Pros**:

- Maximum reuse: parse, conformance, and carry-through all exist.
- One declaration grammar in the language instead of two.

**Cons**:

- `0010:C3` builds the wall this tears down: "Emit keys are NOT tag
  keys … MUST NOT be matched, guarded, written, or read by any
  accessor", and `EmitValue` was deliberately not `TagValue` to stay
  clear of set semantics, `renderValue`, and clear-sentinel paths.
- A widening on a locked Final contract is a new-RDR override of
  0010's core type decision, with every tag feature (accessors,
  matching, min/max) now needing an "except emit" carve-out.

**Reason for rejection**: contradicts the locked contract it would
build on; the mirror (same *move*, separate grammar) keeps the reuse
where it is safe and the wall where it is load-bearing.

### Alternative 3: Validate at resolve time

**Description**: Keep authoring free; at `flow resolve`, check the
selected row's emit values against declarations (or against the
recognized command surface) and refuse or annotate on mismatch.

**Pros**:

- No load/lint changes; the check runs where the answer is produced.

**Cons**:

- Proves only exercised rows — the rare row's typo still ships and
  fires at a close-out, which is the seed's documented failure, now
  with a runtime refusal instead of a wrong command.
- Strains `0010:C3`'s evaluation clause (values uninterpreted at
  resolve) and puts a proof obligation on the hot path lint exists to
  keep offline.

**Reason for rejection**: converts an authoring-time proof into a
runtime surprise; the QOC correctness row is a 2.

### Briefly Rejected

- **Consumer-side golden table only**: the status quo the Problem
  Statement documents failing — covers only the rows in-flight work
  exercises.
- **One global closed emit-value list (not per-key)**: collapses
  every key into one namespace, so a valid value on the wrong key
  passes; loses kinds and dispositions.
- **Making declarations mandatory**: breaks every existing
  emit-carrying model and the opt-in acceptance criterion carried
  from the seed.

## Context

### Background

Tracked as kata `intrastate#vt9n`; kata `intrastate#rg0e` is folded
in: its declared-kind shape is one field of this declaration grammar,
and its lint-enforced reserved-prefix shape is the degenerate form of
the same move — deciding them separately would either reopen this
contract or ship the convention the declaration exists to replace.
No STOP/ROUTE discriminator exists anywhere today — `EmitValue` is
`{Key, Value string}` ("undeclared, uninterpreted, compared by exact
byte equality"), the resolve payload copies it verbatim, and
`internal/graphlint` never reads `Row.Emit`. RDR 0010 adjudicated the
principle while naming this exact widening as deferred: `0010:A6`
verified string-valued emit "as sufficient for the motivating
consumer" and its Consequences state "structured answers wait for a
widening"; a widening on a locked contract is a new RDR, never an
amendment. Consumer stakes: every prose decision table migrated to a
model multiplies the unlinted emit surface, and a cross-file seam
test wants the leg "every declared emit domain member is a real
command or stop token." The acceptance predicate lives in the Problem
Statement.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/table/load.go` (the tag declaration grammar to mirror),
`internal/table/normalize.go::emitSequence` (pass-through, no
validation), `internal/table/model.go::EmitValue`,
`internal/graphlint/` (unchanged — the emit checks land in the load
pipeline), `internal/cli/flow_resolve.go` (resolve payload `Emit`
map).
Governing records: RDR 0002 (tag grammar), RDR 0005 (C1 append-only
envelope), RDR 0006 (lint authority), RDR 0010 (C3, C4, A6 — emit
block ownership and the deferred widening).

## Research Findings

### Investigation

Prior art was read before enumeration (queries and rejected branches:
`evidence/research/propose-prior-art.md`). The in-repo mirror is the
primary prior: `0002:C22`'s tag type-model declaration (`kind`,
`domain`, authored beside the block it constrains, carried losslessly)
and `0002:C24`'s load/lint arity split ⇒ single-rule emit checks
belong to load. The instance question — how peers mark a stop —
opened two citations: ms-conductor `examples/README.md` §Explicit
Termination (declared `type: terminate` steps with a `status` field)
and scxmlcc `doc/user-manual.md` §Final State (`<final>` as a
declared element) ⇒ peers declare the stop structurally, which
disqualifies the prefix convention as the contract. ⚠ the DMN
output-values class claim was SEARCHED AT RESOLVE AND NOT FOUND —
four corpora plus the sibling analysis repo return nothing
(`evidence/research/resolve-prior-art.md`), and the half that would
corroborate declare-then-prove over answers (that conformant tooling
CHECKS authored output entries) has zero support. It stays A6,
unverified and not leaned on; the alignment argument rests on
`0002:C22` and the two peer citations above. Code paths read:
`internal/table/load.go::loadTags` (the declare-then-prove shape to
mirror), `internal/table/model.go::EmitValue` and `declaredKinds`,
`internal/table/normalize.go::emitSequence`,
`internal/cli/flow_resolve.go::resolvePayload`/`emitMap`, and
`internal/cli/lint.go`'s load-refusal arm.

### Key Discoveries

- **Documented** — `0006:C3` phrases the blocking invariant classes
  as "at least", and `0002:C24`'s category list is "at minimum" with
  `reserved_tag_key` (RDR 0008) as the append precedent ⇒ appending
  load categories is a licensed widening; appending an advisory tier
  member is not (`0006:C17` closes it at four) — so no check in this
  RDR may be advisory-shaped.
- **Documented** — `internal/cli/lint.go` loads via
  `table.LoadWithAdvisories` and maps a load refusal to blocking
  findings before any graph proof runs ⇒ the seed's acceptance
  ("each a blocking finding" under lint) is satisfied by load-time
  enforcement; graphlint needs no change.
- **Verified** (repo sweep at Propose; A3 re-verifies as a spike) —
  no checked-in model authors `[emit.` and only
  `models/examples/pricing-decision-table.toml` authors
  `[rule.emit]` ⇒ the opt-in predicate holds over the corpus.
- **Documented** — `internal/table/source.go::sourceDoc` declares no
  `emit` field ⇒ the authoring locus is free (A1), and an older
  binary handed a declared model refuses loudly as
  `unknown_schema_field` rather than silently ignoring the
  declarations (its writer: `decodeStrict`'s unknown-key arm).
- **Documented** — `0010:A6`: the motivating consumer's answers
  decompose into closed authoring-time vocabularies, and judgment
  cells answer with fixed `stopped:<code>` strings ⇒ an enum domain
  with a disposition partition is authorable (A5).

## Trade-offs

### Consequences

- Positive: a declared model's emit surface is proven at authoring
  time — both seed defects become load refusals, and every prose
  table migrated to a model can bring its answer vocabulary with it.
- Positive: the stop-vs-route distinction becomes declared data on
  the envelope; consumers stop re-implementing a string-prefix
  discipline, and intrastate stays generic (tokens are model-authored
  and uninterpreted).
- Positive: the lossless declaration carry (C3) leaves a seam later
  RDRs can consume — finding classes over the declared vocabulary,
  normalized-graph exports — without touching this grammar.
- Negative: opt-in means an undeclared model keeps today's silent
  gap; the guarantee exists only where an author spends the
  declaration effort.
- Negative: whole-model strictness makes adopting the first
  declaration a step: every emit key in the model must be declared at
  once (`scalar` is the deliberate low-cost escape per key).
- Negative: a declared-key-never-emitted or member-never-authored
  situation is deliberately unreported (closed advisory tier); stale
  vocabulary can accrete in a declaration.

### Risks and Mitigations

- **Risk**: domain drift — the declared domain and the consumer's
  real command surface diverge, so the proof holds against a stale
  vocabulary.
  **Mitigation**: out of intrastate's scope by design (generic core);
  the consumer's cross-file seam test ("every declared emit domain
  member is a real command or stop token") closes the loop, and the
  declaration gives that test one authoritative list to read.
- **Risk**: the disposition partition grammar (`[emit.<key>.domain]`
  sub-table) proves awkward to author or decode against real TOML
  tooling.
  **Mitigation**: MVV authors the motivating consumer's `next` key
  first; the flat-array form is the fallback shape and dispositions
  are the only casualty of a grammar retreat, not the vocabulary
  proof.
- **Risk**: inline payload-JSON assertions across the flow test
  corpus make the appended `dispositions` field a wide mechanical
  diff.
  **Mitigation**: A2 sizes it at Resolve; `0010:C4` walked the same
  path with `emit`, and the fixed insertion point (immediately after
  `emit`) confines each inline edit to one field at a known position.
- **Risk**: scalar hollowing — whole-model strictness pushes a bulk
  adopter to declare everything `scalar`, and "declared" reads as
  "checked" (premortem P-2/P-11).
  **Mitigation**: `kind = "scalar"` is visible in the one reviewable
  declaration table (the `0002` bargain: reviewers read the
  declaration, not 500 rows); an advisory finding for scalar-declared
  keys is deliberately NOT taken — the advisory tier is closed
  (`0006:C17`) — so the check is review- and consumer-seam-side.
- **Risk**: disposition-token drift — a model renames `stop` to
  `halt` and every consumer's `dispositions[k] == "stop"` test goes
  quietly false (premortem P-12).
  **Mitigation**: same seam as domain drift: intrastate is generic by
  design, and the consumer pins its own token vocabulary with a
  contract test against its model — one authoritative list to pin is
  what the declaration provides.

### Failure Modes

- **Breaks visibly**: a declared model with an undeclared key,
  out-of-domain value, or malformed declaration refuses at load —
  `intrastate lint` exits nonzero with one blocking finding — the
  first refusal the fail-fast pipeline reaches (category slug in the
  finding's `code`, offending file in `locator`) — and
  `flow resolve`/`flow next` refuse `flow-model-invalid` with the same
  finding. A model carrying several defects is fixed and re-run one
  refusal at a time, as it is today for every other load category. Diagnose from the
  finding's category + locator; recover by fixing the declaration or
  the rule. Blast radius is authoring-time, not run-time: every check
  reads only the model file against itself (never the consumer's
  command surface), so a refusal can only be introduced by an edit to
  that file — and fires at the lint gate on that edit, not later at a
  close-out (the premortem's whole-table-offline scenario requires
  shipping a model whose own lint was red).
- **Fails silently (accepted, bounded)**: a `scalar`-declared key
  admits any value — declared-but-unvalidated is the documented
  escape hatch, visible in the model source. A model with zero
  declarations is today's world: nothing new fires. A stale-but-valid
  domain member routes a wrong-but-declared answer, and a declaration
  that faithfully copies a typo proves consistency with the mistake
  (premortem P-1); both classes are external-truth questions the
  consumer seam test ("every declared member is a real command or
  stop token") exists to catch — the declaration's contribution is
  giving that test one authoritative list.
- **Old binary, new model**: a model carrying `[emit]` refuses on an
  older binary as `unknown_schema_field` (strict decode) — loud,
  never a silent ignore of the declaration.

## Implementation Plan

### Prerequisites

- [x] Critical Assumptions verified: A1–A5 Verified at Stage 4. A6
      (DMN prior art) is carried Pending and NOT load-bearing — no
      corpus reaches DMN material; the alignment argument runs on
      `0002:C22` and the two opened peer citations instead.

### Minimum Viable Validation

1. Re-author the seed's two-rule adversarial table — first rule's
   emit value misspells a command, second rule types the key `nxet` —
   and confirm it still lints exit 0 with zero findings (the defect,
   reproduced).
2. Add an `[emit.next]` enum declaration with a
   route/stop-partitioned domain; run `intrastate lint`.
   **Expected**: exit nonzero with ONE blocking finding — the load
   pipeline is fail-fast (C2), so the run reports whichever defect it
   reaches first, with its category slug in `code` and a locator.
   Fix that rule and re-run: the second defect now surfaces, likewise
   as one finding. Across the two runs both categories are observed —
   `emit_value_out_of_domain` for the misspelled value and
   `unknown_emit_key` for the typo'd `nxet` key. Do NOT expect both in
   one run; one run, one finding is the locked cardinality
   (`0008`, asserted).
3. Fix both rules; lint exits 0. Run `flow resolve` selecting a rule
   whose value is listed under `stop`. **Expected**: payload carries
   the authored `emit` unchanged and `dispositions` mapping the key
   to `stop`, positioned immediately after `emit`.
4. Delete the `[emit]` table; re-run lint and resolve over the
   original (defective) table. **Expected**: exit 0, no findings,
   payload carries `dispositions: {}` — byte-identical behavior to
   today apart from the appended empty field.

### Phase 1: Declaration grammar and load proof

Extend the source schema with the `[emit]` table, load declarations
beside tags, and refuse `malformed_emit_declaration` /
`unknown_emit_key` / `emit_value_out_of_domain` in the load pipeline
(C1, C2) — the whole opt-in predicate lives here. Reuse
`internal/table/load.go::conformKind`/`conformDomain` through a narrow
helper (A4); no `internal/guard` import.

**Amend three stale comments in the same change.** Each states emit's
openness as a positive commitment and reads false once a declaration
can exist — leaving them is the real inconsistency risk, not the
grammar:
`internal/table/model.go::EmitValue` ("undeclared, uninterpreted, and
compared by exact byte equality"),
`internal/table/normalize.go::emitSequence` ("they are undeclared and
uninterpreted"), and `internal/table/dump.go::renderEmit` ("an emit
key has neither: it is undeclared, so there is no kind to consult").
The amendment is narrow and must preserve what `0010:C3` actually
fixed: "undeclared" there means **not a tag** — no provenance, no
accessor, no match/guard/write participation — every clause of which
C1 keeps. `[emit]` is a SEPARATE NAMESPACE from `[tags]`, so an emit
key still never enters `Model.Tags` and
`internal/table/emit_0010_test.go::TestReq27_AnEmitKeyIsNotATagKey`
keeps passing unchanged; only the "no kind to consult" and
"uninterpreted" wordings narrow to "no TAG kind" and "never parsed,
canonicalized, or converted" (C1's lexical-check clause).

### Phase 2: Normalized carry

Carry declarations losslessly onto the normalized model as a new
non-TagDecl carrier type (C3); dump and kernel untouched.

### Phase 3: Envelope surfacing

Insert `dispositions` into `resolvePayload` after `emit`, joined from
the selected row's authored values against the carried declarations
(C4). **The licensed diff, enumerated by A2** — all of it lands in
this one change, and nothing outside it:
(i) the 28 inline whole-payload assertions in
`internal/cli/flow_demand_0011_test.go`, 7 of which pin the
`"emit":{…},"next"` adjacency;
(ii) the 3 `resolveGolden` byte-identity goldens and the
`shippedResolveGoldens` sweep in the same file — RDR 0011's guard
against exactly this class of change, so their update is authorized
HERE by name;
(iii) `internal/cli/decision_table_0010_test.go::TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire`,
whose 13-key `slices.Equal` list becomes 14 keys and whose
`NumField() != 14` becomes 15 — it encodes `0010:C4`'s position clause,
so this is an amendment to a peer RDR's spec test and is called out
rather than absorbed as a fixture edit.
`TestReq43` (escaped between clear and escape_class) asserts relative
indices and needs no edit. Any assertion outside this list going red
is a signal the append did more than C4 licenses — investigate, do not
update it.

**The §D1 registration obligation, whichever record lands second.**
`JDR 0002 §D1` requires every success-payload field of a projecting
verb to be assigned to ECHO or PLAN **when it is added**, enforced by
a reflective oracle for which "an unassigned new field is a test
failure, never a silent default". `dispositions` is assigned PLAN
there by name. That oracle does not exist at HEAD — 0023 is Final and
unimplemented — so: if 0023 lands first, this change registers
`dispositions` with the oracle it shipped; if 0024 lands first, 0023's
implementation registers it when the oracle arrives. Either order
works and neither contract moves (`0023:A4`), but the registration is
not optional and is named here so it is not discovered as a red
reflective test.

### Phase 4: Example and docs

Declare the pricing example's emit keys
(`models/examples/pricing-decision-table.toml`), and extend
`docs/cli-output-contract.md`'s resolve section with `dispositions`
and the declaration grammar — the doc is the contract's worked-payload
home.

## Validation

### Testing Strategy

The MVV is the acceptance spine — the seed's adversarial table driven
from reproduced defect to green declaration to opt-out identity. Done
= every scenario below is green and the MVV run recorded; every C1
refusal arm and every C2 category maps to at least one scenario.

**Code paths under test, and the verification each rests on.** The
load-side scenarios (1–3) exercise `internal/table/load.go`'s refusal
arms beside the existing `loadTags` step, reusing
`load.go::conform`/`conformKind`/`conformDomain` for value conformance
(A4 — verified reusable in-package, no `internal/guard` import). The
grammar sits on the free `[emit]` key in
`internal/table/source.go::sourceDoc`, refused today by
`decodeStrict`'s `CatUnknownSchemaField` arm (A1). Scenario 3's
opt-out baseline is the corpus sweep and green suite in
`evidence/spikes/a3-corpus-sweep.md` (A3). Scenario 5 exercises
`internal/cli/flow_resolve.go::resolvePayload`'s join, whose licensed
test-corpus diff is enumerated in A2 and Phase 3. The single-finding
cardinality every load scenario asserts is the fail-fast behavior
captured in `evidence/spikes/c2-finding-multiplicity.md`.

1. **Scenario**: table-driven load of malformed declarations — an
   unknown `kind` token; an `enum` with no domain, an empty domain,
   an empty-string member, or a duplicate member (across disposition
   lists included); a `domain` on a non-enum kind; an empty-string
   disposition token (C1).
   **Expected**: each case, loaded on its own, refuses
   `malformed_emit_declaration` at load with exactly one blocking
   finding, before any row is yielded. One defect per fixture — the
   pipeline is fail-fast, so a fixture carrying two defects proves
   nothing about the second.
2. **Scenario**: a declared model with an undeclared `[rule.emit]`
   key on an ordinary rule and on an escape rule, and with
   out-of-domain values per kind — a non-member enum value, a
   non-`true|false` bool, a non-integer-literal int, an arbitrary
   `scalar` value (C2).
   **Expected**: `unknown_emit_key` on both rule classes;
   `emit_value_out_of_domain` for the enum/bool/int cases; `scalar`
   never refused — each from its own single-defect fixture, one
   blocking finding per run. A model carrying an emit defect AND a
   coexisting structural one asserts the fail-fast contract instead:
   exactly one finding, and the test does NOT assert WHICH — the
   pipeline order is unspecified and pinning it would fabricate a
   guarantee (the same "deliberately unasserted" move
   `internal/table/reserved_key_0008_test.go` makes).
3. **Scenario**: zero-declaration opt-out (C2's opt-in leg; A3's
   corpus sweep as the fixture check).
   **Expected**: full suite green with the loader change and no model
   changed; no new refusal reachable; the only observable delta is
   C4's appended `dispositions: {}`. This is A3's post-change leg,
   which could not run at Draft — the A3 spike establishes the
   pre-change baseline (build clean, suite `ok`, all four
   `models/*.toml` lint exit 0) that this scenario re-runs after the
   change to prove byte-for-byte equivalence.
4. **Scenario**: declaration carry through normalization (C3).
   **Expected**: key, kind, domain members, and member dispositions
   read identically off the normalized model; kernel and dump
   surfaces carry none of it.
5. **Scenario**: `flow resolve` payloads across the join matrix — an
   enum member carrying a disposition, a flat-array domain, non-enum
   kinds, an empty emit block, an undeclared model, and an escape-row
   rescue (C4).
   **Expected**: `dispositions` maps only member-disposition hits and
   is `{}` in every other case — never `null`, never omitted — sits
   immediately after `emit`, joins from the escape row's own authored
   values by `Plan.RuleID`, and renders in text mode through the
   generic payload renderer; `flow next` payloads carry no
   `dispositions`. The post-change wire key order is exactly these
   fourteen keys, in this order (the pre-change thirteen from
   `decision_table_0010_test.go::TestReq39`, with `dispositions`
   inserted at index 9):

   ```
   model, revision, observed, owned, readers, outcome, rule, gates,
   emit, dispositions, next, writes, clear, escaped
   ```

   and `reflect.TypeOf(resolvePayload{}).NumField() == 15`
   (`escape_class` is `omitempty` and absent from the wire on a
   non-escaped plan). These two values are the normative fixture this
   scenario asserts; both are read from the current test's asserted
   values, not invented.

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

- Peer RDR elements: `0002:C22`, `0002:C24` (declaration grammar and
  the load/lint arity split), `0005:C1` (envelope + findings
  mapping), `0006:C3`, `0006:C13`, `0006:C17` (finding tiers),
  `0010:C3`, `0010:C4`, `0010:A6` (emit ownership and the deferred
  widening this RDR takes up).
- Source reviewed: `internal/table/load.go`,
  `internal/table/model.go`, `internal/table/normalize.go`,
  `internal/table/source.go`, `internal/table/category.go`,
  `internal/cli/flow_resolve.go`, `internal/cli/lint.go`,
  `internal/graphlint/taxonomy.go`.
- Prior art opened: ms-conductor `examples/README.md` §Explicit
  Termination; scxmlcc `doc/user-manual.md` §Final State (`<final>`)
  — search record at `evidence/research/propose-prior-art.md`.
- Stage 4 evidence: `evidence/research/resolve-prior-art.md` (A6
  negative — DMN searched across four corpora, not found);
  `evidence/spikes/a3-corpus-sweep.md` (A3 — zero `[emit.` repo-wide,
  green baseline); `evidence/spikes/c2-finding-multiplicity.md` (C2 —
  the load pipeline is fail-fast, one refusal per run).
- Related issues: kata `intrastate#vt9n`, kata `intrastate#rg0e`
  (folded), kata `srz2` (RDR 0023 projection consumer).
