# Recommendation 0023: Resolve-envelope projection opt-out

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
- **Profile**: foundational — one independent contract, the
  caller-controlled projection axis on the resolve envelope (C1
  surface/wire + C2 partition instance), a cross-RDR producer: it
  narrows `0005:A6`, instances JDR 0002 §D1, and composes with
  `0024:C4`.
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
- **Related Issues**: kata `intrastate#srz2` (1604); kata `rg0e`
  (sibling resolve-envelope decision, seeded in the same batch —
  independent asks, one natural instrument)
- **Predecessors**: 0005-skill-integration-cli-contract,
  0011-flow-next-match-conditioned-candidates
- **Overrides**: `0005:A6`'s pinned minimum success payload fields for
  `resolve` (the 0005 Technical Design "Envelope" list: `model`,
  `revision`, `observed`, `owned`, `readers[]`, `outcome`, `rule`,
  `gates[]`, `next`, `writes`, `clear[]`, `escaped`, `escape_class`) —
  narrowed to the DEFAULT mode: under the opt-in `--plan-only`
  projection (C1 below) the echo-group fields `model`, `observed`,
  `owned`, `readers`, `outcome` are absent from the success payload
  (`revision` stays — plan-side identity slot, C2). Default-mode output is byte-identical to today's; every
  other 0005 obligation on `flow resolve` — one terminal envelope,
  refusal codes, exit groups, gate ordering, both-modes agreement — is
  unchanged. No clause of 0011 is overridden: `0011:C2`'s fence is
  `--all`'s non-registration on the other three verbs, which this RDR
  leaves intact, and its `--select` rejected-spelling pin is on `flow
  next` (`internal/cli/flow_all_0011_test.go::TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand`);
  what is deliberately reopened is only 0011's rejection GROUND
  ("the defect is the default") as applied to a projection flag —
  inapplicable here because the echo default is agreed correct and
  stays.
- **Seam Lineage**: no prior accretion

## Problem Statement

An agent-driven caller of `flow resolve` pays for every output byte as
transcript tokens re-paid on each subsequent turn — the dominant cost
of the call (wall time is 5ms). The seed measured a 48-fact
decision-table call at 1529 bytes of full JSON against an
`emit` + `rule` answer of ~430 — about 72% of that resolve's output
was the caller's own input echoed back under `observed.*`. The echo is
the **right default** and is settled: a live consumer deliberately
dropped its own status call because the resolve already returns every
fact it renders. The ask is only an opt-out for chained or secondary
calls — the second `--outcome` call in a chain receives the same 48
facts again as pure duplication — and every additional model a
consumer migrates multiplies the per-run echo.

Two measurement notes, so the seed's headline is not read as the
shipped claim. (1) BYTES ARE THE MEASURAND throughout this RDR,
including in A1's verification and the Performance Expectations. The
cost is incurred in tokens, but no clause, oracle, or assumption here
states a token figure: bytes are used as the proxy, the tokenizer is
the consumer's and is not ours to pin, and the proxy is monotone —
deleting whole keys and their values never increases a token count.
An oracle asserting a token reduction is deliberately NOT specified;
A1's bar is a byte bar and is meant to be discharged by byte evidence.
(2) The 1529 B above is the SEED's own motivating call and is retired
as a headline by A1, which re-measured the class on a synthetic
48-fact model (`evidence/spikes/a1-byte-width.md`, S1b) at 708 B
default → 146 B projected. The two numbers are different calls, not a
discrepancy: A1's is the shipped-partition measurement on a model
authored for the spike, and it — not 1529 — is the figure every later
section quotes.

The fork is the mechanism and its blast radius against two locked
records. (a) A field-projection flag (`--select emit[,rule]`): but
`0011:D-naming` rejected the `--select` spelling by name, and
`0011:A14` makes `--all`'s non-registration on
`resolve`/`read-state`/`set-state` a load-bearing structural negative
enforced by absence oracles in `internal/cli/flow_all_0011_test.go` —
fences that shape, without by themselves blocking, a new flag on
`resolve`; 0011's rejection ground ("the defect is the default") does
not bind a surface whose default is agreed correct (Background carries
the source read). (b) A verbosity tier (`--as json`
omits `observed` unless `--explain`): collides with `0005:C1`, which
requires text output derive from the same verb-specific result — the
`respond` gateway renders marshalled JSON through one `flatten`, so a
json-only omission is structurally impossible, and a mode-dependent
one mints the per-verb text template 0011 already declined as
"editing a predecessor's surface for cosmetics." (c) A quiet/suppress
flag dropping only the echoed input. The routable question: should the
resolve payload have a caller-controlled projection axis at all — and
if so, does it live on the flag surface `0011:C2` fenced, in the
`respond` text/JSON gateway `0005:C1` requires to agree, or nowhere?
Each answer also sets the terms on which a future verb could inherit
the projection — and what this RDR settles about inheritance is
bounded, deliberately: it decides the MECHANISM a later verb would
reuse (a fixed RDR-owned partition projected before `respond.OK`,
never a caller-supplied field list, never a gateway tier) and it
registers the flag on `flow resolve` ALONE, by a structural negative
that closes over verbs added later (C1). It does NOT oblige, permit,
or forbid any particular later verb to acquire a projection: that
question belongs to JDR 0002 §D1's later-verb clause, which obliges
the unproposed 0021 to nothing. So a future verb author reads this
record for the shape a projection must take if they add one, and
reads §D1 for whether they may — the two are separate answers and
this RDR gives only the first.

## Critical Assumptions

- **A1 The saving survives the SHIPPED partition: measured with the
  full plan group (not the seed's `emit`+`rule` subset), across at
  least two fixture shapes with different fact/gate ratios, the
  projected width is materially smaller on the chained-call class
  the RDR targets — the seed's 1529→~430 figure is the echo SHARE of
  the motivating call, not the shipped reduction.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: live-run byte table (author-approved 2026-08-29;
    `evidence/spikes/a1-byte-width.md`): pricing 2×2 290→152 B
    (47.6% saved); synthetic 48-fact motivating class 708→146 B
    (79.4%); release ship-clean (writes+clear) 392→191 B (51.3%);
    release begin (gate+writes) 412→238 B (42.2%) — materially
    smaller on every shape with the full shipped PLAN group
    retained. The seed's 72% was the echo share against
    `emit`+`rule` alone; 79.4% is the honest shipped reduction on
    that class. MVV step 6 re-records this table at implementation.
  - **If wrong**: the flag saves too little to justify a new surface
    on a locked verb, and the RDR routes back to "nowhere" (no
    projection axis).
- **A2 Go's encoder supports the wire shape C1 fixes, escaping the
  premortem's trilemma (P-1): the echo fields can be made ABSENT
  under the flag — candidate mechanism: pointer-valued echo fields
  with `omitempty`, always populated in default mode, nilled at
  projection — while default-mode bytes stay byte-identical
  (including `observed`/`owned` present as `{}` when empty, which
  bare-map `omitempty` would drop), surviving keys keep declaration
  order, and text mode stays deterministic (the flattener already
  sorts keys from the decoded map, `respond/text.go::flatten`, so no
  map-iteration-order path exists).**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: encoder spike
    (`evidence/spikes/a2-encoder-mechanism.md`): pointer-valued echo
    fields with `omitempty`, populated by decoding a live CLI
    reference line, marshal byte-identical to shipped default output
    — including `"owned":{}` (a non-nil pointer to an empty map
    survives `omitempty`) — and, nilled, drop all five echo keys:
    absent, declaration order kept, surviving bytes equal to jq
    key-deletion of the reference. The bare non-pointer `omitempty`
    hazard reproduced (drops empty `{}`/`[]`). Text mode: 10 runs,
    one hash (`respond/text.go::flatten` sorts keys). Normative
    fixture (author-approved 2026-08-29): the projected wire record
    of the pricing 2×2 call, recorded in the spike artifact and
    cited by Testing Strategy S1/S2 Expected.
  - **If wrong**: C1's absence-not-null and byte-identity clauses are
    unimplementable as specced and the mechanism (or the clause) must
    change before lock.
- **A3 No shipped oracle or fixture asserts an echo field's presence
  in a way an OPT-IN flag moves: every existing resolve-payload
  assertion runs default-mode, so the 0005/0010/0011 suites do not
  move and "default unchanged" holds without editing a predecessor's
  tests.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: swept all 22 test files invoking `flow resolve`
    (count re-executed at the critique lens, 2026-08-29:
    `rg -l 'flow.*resolve|"resolve"' internal/cli --glob '*_test.go'`
    → 22; the earlier "16" and "21" were both miscounts, the
    substantive claim unchanged across all three): none passes a
    projection flag (`rg plan-only internal/` → no matches); every success-payload
    assertion runs default mode — key-set
    (`internal/cli/flow_resolve_0005_test.go`), byte-identity and
    wire key-order (`internal/cli/decision_table_0010_test.go`).
    Caveat carried to implementation: `decision_table_0010_test.go`
    pins `resolvePayload` at exactly 14 struct fields — the A2
    mechanism keeps the count (pointer conversion adds no field).
  - **If wrong**: the change stops being additive — moved predecessor
    oracles are the "editing a predecessor's surface" cost this RDR
    claims to avoid.
- **A4 0024's `dispositions` composes with this projection because
  its derivation reads PLAN-GROUP inputs only: `0024:C4` joins each
  token from the `[emit]` declaration to "the selected row's authored
  value" by `Plan.RuleID` — never from `observed`/`owned` — so
  keeping `dispositions` while projecting the echo group strips no
  derivation input, and whichever RDR lands second changes nothing in
  the other's contract (the premortem's seed-3 check: this cites the
  quoted join rule, not a summary of the sibling).**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: `0024:C4` re-read whole at Resolve (2026-08-29):
    `dispositions` joins from the selected row's OWN authored values
    by "the same single `Plan.RuleID` join path as `emit`" —
    plan-group inputs only, `observed`/`owned` never read — and
    JDR 0002 §D1's Resolved clause records the same assignment
    ("`dispositions` is a PLAN-group field… consistent by
    citation"). Neither contract moves in either landing order.

    Scope of this stamp, narrowed at the critique lens (2026-08-29):
    it covers the CONTRACTS only. It does NOT claim 0024 is
    oracle-free against this RDR's baseline — `0024:C4` appends
    `dispositions` LAST in `internal/cli/flow_resolve.go::resolvePayload`,
    and two shipped 0010 oracles are extensional over that struct:
    `decision_table_0010_test.go:435` asserts `NumField() != 14` and
    `:427` a `slices.Equal` exact wire-key-set comparison. Both go red
    when 0024 lands, in either order. That cost is 0024's own additive
    append (the same one `0010:C4` took) and is not created by this
    projection — A2's conversion changes field TYPES, not the count,
    so this RDR alone moves neither assertion, which is what keeps A3
    true as written. Recorded here because A3's "no predecessor
    oracle moves" holds for THIS RDR and would be misread as covering
    the pair.
  - **If wrong**: the joint decision homed at `cli/0023:C2` reopens
    and one of the two envelopes must move.
- **A5 The whole-tree registration oracle C1 mandates is
  implementable in the house idiom without colliding with the `--all`
  oracles: a walk of the command tree from the root can enumerate
  every command registering `plan-only` (verb-local registration is
  visible to a per-command `Lookup`), the four-verb vacuity guard
  transfers, and nothing in
  `internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup`
  or `registerSelectionFlags` registers `plan-only` under any other
  name. Scoped at the cove sweep to the flag name this RDR adds: the
  general claim that no name is registered twice across the tree is
  FALSE (`help --all`), which is why C1 mandates a total walker rather
  than reuse of an existing partial one.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: the collision half is VERIFIED — the `--all` probes
    are exact-name `Lookup("all")`
    (`internal/cli/flow_all_0011_test.go`, three probe sites), so a
    resolve-local `plan-only` is invisible to them; the root builds
    via plain `AddCommand` (`internal/cli/root.go`); the four-verb
    vacuity guard transfers as a pattern; `registerSelectionFlags`
    registers only `model`/`flow`/`artifact`. The IMPLEMENTABILITY
    half was re-opened by the cove sweep (2026-08-29) and its original
    evidence retracted: no shipped test walks `Commands()` recursively
    — every sweep is a fixed two-level loop over the `flow` group
    (`flow_all_0011_test.go`, `flow_surface_0005_test.go`,
    `flow_input_0005_test.go`) — and the repo's only recursive walker,
    `internal/cli/help_all.go::walkCommandTree`, SKIPS children named
    `help` or `completion` before recursing, with its `docs.go` callers
    further gating on `c.Hidden`. So `docs` enumerates because that
    caller admits hidden commands, not because the walker is total:
    reusing it would assert a strictly weaker negative than C1's
    "whole command tree". The sweep also found a live counterexample
    to the generalization — `internal/cli/help_all.go::wireHelpSubcommandAll`
    registers a second `--all` on the auto-generated `help` command,
    which today's 0011 oracle misses precisely because it descends the
    `flow` group only. `plan-only` is unaffected (different name), but
    the whole-tree idiom C1 mandates has no shipped exemplar.
  - **Verification plan**: write the total walker (no name-skip, no
    `Hidden` gate) against the real root and confirm it enumerates
    `help` and `completion`; assert the `plan-only` registrant set is
    exactly `{flow resolve}` under it. Method: Spike, at implementation
    Phase 2 — C1/S4 now fix the scope it must cover. Scope correction
    (3amigo, 2026-08-29): a bare `NewRootCmd()` tree does NOT contain
    `completion` — probed, its children are `docs`, `flow`, `help`,
    `lint`, `version`, and `help` is present only because
    `help_all.go` force-inits it; cobra creates `completion` in
    `InitDefaultCompletionCmd` from `ExecuteC`. So the spike must
    materialize both auto-generated commands via cobra's own
    initializers before walking and assert they are IN the walked
    set, else it verifies a walker over a tree lacking the command
    class the clause exists to cover.
  - **If wrong**: C1's whole-tree structural negative cannot be pinned
    in the house idiom and needs its own form.
- **A6 A refusal path carries nothing to project: `respond.Fail`
  renders `clierr.CLIError` (+ `findings`) with no request-echo
  member, so `--plan-only` cannot change a single refusal byte and
  "report-only" holds on failures by construction.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: the refusal envelope's fields are `code`,
    `message`, `param`, `detail`, `hint`, `findings`
    (`internal/cli/clierr/clierr.go`) — no echo-group member;
    `respond.Fail` marshals the `*CLIError` directly and never the
    verb payload struct. (`Finding`'s `model,omitempty` is RDR 0006
    lint-finding metadata, not a resolve request echo.)
  - **If wrong**: the flag changes refusal output, breaking C1's
    report-only clause and the caller's error handling.
- **A7 The echo group really is an echo: `observed` renders the
  caller's `--tag` set verbatim (`0005:C1` refuses rather than
  coerces — a malformed or nonconforming tag literal is
  `flow-tag-invalid`, an owned/reserved key is refused before any
  accessor), and `model`
  and `outcome` are the caller's own flag values copied through
  (`runFlowResolve`: `payload.Outcome` is `--outcome`'s string) — no
  CLI-side defaulting, coercion, or expansion produces an observed
  entry the caller did not supply. One value-preserving exception is
  in scope, not a refutation: a set-kind literal is re-encoded to
  the canonical bytes of the same members (the form a conforming
  caller already speaks).**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/flow_input.go::parseTags` refuses
    owned/reserved/duplicate/malformed/nonconforming tags; the sole
    rewrite is `canonicalSet`'s sort/dedup/compact re-encode of a
    VALID set literal — exactly the named value-preserving exception
    (the claim's wording was narrowed to match at Resolve);
    `runFlowResolve` copies `--model`/`--outcome` through verbatim.
    Round disposition (author-approved 2026-08-29): `owned` and
    `readers` stay ECHO — `readers` is a pure function of
    model + outcome (`internal/cli/flow_exec.go::invokedReaders`,
    derivable from the caller's own request), and `owned` rolls
    forward as prior-owned overlaid with `next` (a clear is a
    `<clear>` sentinel write; untouched keys never change), with
    `flow set-state`'s read-back refusing unplanned foreign changes;
    corpus + precedent grounding in
    `evidence/research/resolve-partition-prior-art.md`.
  - **If wrong**: `observed` is derived output, its projection hides
    what the kernel actually matched against, and the field moves to
    the plan group before lock.
- **A8 The strict-width clause is satisfiable in TEXT mode on the
  same unit as JSON: the projected `--as=text` rendering is strictly
  fewer total bytes than the default rendering of the same request,
  not merely a line subset.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: not yet measured. A2 established text-mode
    DETERMINISM only (10 runs, one sha256,
    `evidence/spikes/a2-encoder-mechanism.md` §Text-mode determinism);
    the desk trace's "default 15 lines → projected 9" is a LINE count,
    and no artifact records text byte widths. The claim is very likely
    true — the six dropped lines (`model`, `observed.region`,
    `observed.tier`, `outcome`, `owned`, `readers`) carry non-empty
    content — but C1 now asserts it as a MUST in both modes, so it is
    booked rather than assumed.
  - **Verification plan**: measure total rendered bytes of
    `--as=text` ± the flag on the pricing 2×2 call and on one
    gate/write-heavy shape; record both in the A1 table's unit.
    Method: Spike, at implementation Phase 2 alongside the S5 oracle.
  - **If wrong**: C1's both-modes width clause is unsatisfiable as
    written and the clause narrows to JSON, leaving S5's subset
    assertion as text mode's only width guard.
- **A9 Default-mode byte identity under A2's pointer conversion rests
  on an invariant no contract states: every container writer feeding
  the echo group returns a non-nil value, so a pointer to it renders
  `{}`/`[]` and never `null`. Under the conversion this incidental
  property becomes load-bearing — a nil container reached by any
  writer renders `"readers":null` (or `"owned":null`) under
  `omitempty`, silently changing DEFAULT-mode bytes for callers who
  never opted in.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: the four container writers all allocate
    unconditionally — `tagMap` and `observedTagMap`
    (`internal/cli/flow_exec.go`), `emitMap`
    (`internal/cli/flow_resolve.go`) and `readerIDs`
    (`internal/cli/flow_next.go`) each open with `make(...)` — so no
    shipped path produces a nil container today and the claim's
    premise holds on `main`. What is NOT established is that the
    property is contracted or guarded: only `emitMap`'s is deliberate
    (`0010:C4`'s never-`null` clause), the rest are incidental, and
    the critique lens verified by execution that a pointer to a nil
    slice does render `null` under `omitempty`. The gap is that A2's
    spike proved byte identity on POPULATED reference values, which
    cannot exercise the nil arm.
  - **Verification plan**: enumerate every assignment reaching the five
    echo fields at the `resolvePayload{…}` site and confirm each source
    is unconditionally allocated; then assert it — a default-mode
    oracle over a request whose containers are empty (no `--tag`, no
    owned keys, no readers), asserting the rendered bytes carry
    `{}`/`[]` and no `null`. Method: Source Search + oracle, at
    implementation Phase 1, before the conversion lands.
  - **If wrong**: A2's mechanism changes default-mode output for an
    untested request shape, breaking C1's byte-identity clause and
    §Consequences' "no consumer changes" promise — and S1 cannot see
    it, since it compares projected against default within one build.

## Proposed Solution

### Approach

Give `flow resolve` — and only `flow resolve` — a caller-controlled
projection axis, as one boolean flag: `--plan-only`, default false.
Default-mode output is byte-identical to today's; the echo default is
settled and stays. Under the flag the success payload carries exactly
the **plan group** the verb's own extended help already names under
"Reading a successful plan" (`internal/cli/flow_resolve.go::flowResolveExtendedDesc`:
`rule`, `next{}`, `writes{}`, `clear[]`, `emit{}`, `escaped`,
`escape_class`, `gates[]`), plus `revision` as the plan-side identity
slot (loader-produced, not request echo, and presently constant-empty
— C2 carries the rationale),
and omits the **echo group** (`model`,
`observed`, `owned`, `readers`, `outcome`) as absent keys.
The projection is applied to the verb-specific result *before*
`respond.OK`, so both output modes render the same projected content
through the shipped generic renderer
(`internal/cli/respond/text.go::flatten`) — `0005:C1`'s two-modes
agreement holds by construction, not by a second template. ⇒ the axis
is a fixed, RDR-owned partition (the `gh pr diff --name-only` form),
not a caller-supplied field list (the `gh --json <fields>` form):
the partition is what makes the projection safe (a caller can never
drop `escaped`) and cheap (no field universe to police, one payload
shape per mode-width).

The scored matrix (Decision Rationale) ran over four answers to the
routable question — fixed projection flag, named-field selection,
gateway verbosity tier, no CLI axis — and the fixed flag wins on
safety, fence compatibility, and parse stability.

### Technical Design

`runFlowResolve` builds `resolvePayload` exactly as today; when
`--plan-only` is set it hands `respond.OK` the projected form of the
same value — same field values, echo fields absent. Nothing upstream
of payload assembly reads the flag: readers, kernel call, gate run,
and every refusal path are flag-blind. The respond gateway is
untouched; the projection is verb-owned, which is what lets a later
verb adopt (or decline) the same axis without a shared mechanism
(C2's envelope-general clause).

#### Normative Contracts

**C1**

```normative
flow resolve MUST accept a boolean flag --plan-only, default false,
long form only (no shorthand). Default-mode output is byte-identical
to today's: the request echo is the correct default and this RDR does
not change it.

Under --plan-only a SUCCESS payload MUST carry exactly the plan
group — rule, gates, emit, next, writes, clear, escaped,
escape_class, revision, plus every field a successor assigns to the
plan group under C2 — and MUST NOT carry the echo group: model,
observed, owned, readers, outcome. An omitted key is ABSENT, never
null and never an empty placeholder. Every carried field's rendering
MUST be byte-identical to its default-mode rendering (same encoder,
HTML escaping disabled), and the carried keys keep their
default-mode relative order. The projected TOP-LEVEL KEY SET is
itself normative: an oracle MUST assert it as an explicit literal
list, so a later payload field reaches the projected width only by a
conscious C2 assignment, never by falling through an omit-list.

--plan-only is REPORT-ONLY, and report-only is ORACLE-ENFORCED, not
aspirational: it MUST NOT change what is decided or how the run
fails, and the build MUST carry a differential oracle asserting,
over the same request with and without the flag: identical exit
codes; byte-identical refusal envelopes (CLIError, findings — never
projected, and carrying no echo group to project); the projected
success payload a strict key-subset of the default payload with
byte-identical values on every carried key; an identical
invoked-reader set; and the projected encoding STRICTLY SHORTER than
the default. The width assertion is not implied by the subset one — a
later plan-group field can grow the projected payload past today's
full width with the key-set and partition oracles still green — so
the flag's whole reason for existing is itself oracle-enforced. A later change that skips work whose only
consumer is a projected-away field breaches this clause.

Two of those assertions need their measurand fixed, because the
clause is otherwise satisfiable by an oracle that measures the wrong
thing.

STRICTLY SHORTER is measured on the FULL EMITTED LINE — the complete
{"type":"ok","data":{…}} NDJSON record as written, excluding the
trailing newline — not on the .data payload alone, and it is asserted
in BOTH output modes. The TEXT-mode half of that "both" rests on A8,
presently Pending and unmeasured: JSON-mode width is measured (A1),
text-mode width is not. The clause binds both modes as written; if
A8's Phase 2 measurement refutes it, the remedy is A8's own — the
clause narrows to JSON before lock, not after, since this RDR does
not amend a locked contract. The envelope-inclusive unit is the one the
caller actually pays for, and it is the stricter reading (a constant
wrapper makes any payload-level saving a smaller proportion of the
line, never a larger one). In text mode the measurand is the total
rendered byte count of the emitted lines; a text projection that
drops no bytes fails this clause even though the S5 subset assertion
would still pass, a subset being satisfied by the equal set.

The IDENTICAL INVOKED-READER SET cannot be read off the projected
run's payload, because readers is an ECHO field this flag projects
away, and the shipped helper reading it fails hard on absence
(internal/cli/flow_fixtures_0011_test.go::readersOf).

The measurand is therefore OBSERVED READER EXECUTION, not a recomputed
derivation. Comparing invokedReaders(model, outcome) across the two
runs is NOT an admissible form of this assertion: it is a pure function
of two arguments the flag does not touch
(internal/cli/flow_exec.go::invokedReaders takes *table.Model and
outcome), so it returns equal sets on every implementation — including
one that never runs a reader at all. Such an oracle is green by
construction and cannot fail, which makes it worse than absent: it
reads as coverage of exactly the breach this clause names.

The oracle MUST instead observe which readers ACTUALLY RAN in each
run — recording execution at the reader invocation site (a counting or
recording seam around the reader pass, the run's own side effects,
never a re-call of the planning function) — and assert the two
observed sets are equal. The DEFAULT-mode run of the same request
supplies the expected set. The discriminating property is that a
projected run which skips the reader pass MUST turn this assertion
red; an assertion that cannot distinguish that case does not satisfy
this clause. Asserting it from the projected payload is not a weaker
form of this check — it is unwritable; asserting it by recomputation
is not weaker either — it is vacuous.
Presence-rule fields keep their own contracts inside the plan group,
unchanged and not restated here: emit stays present as {} (0010:C4),
and escape_class keeps whatever presence rule its producer already has
(0005:A-3 — omitted when unescaped, and also when an escaped row's
class is unprobeable, `flow_resolve.go::escapeClassOf`). This
projection neither widens nor narrows those rules: a field the default
mode omits is omitted under the flag for the same reason, so an oracle
MUST assert presence-rule fields by comparison against the SAME run's
default output, never against an unconditional literal.

Text mode renders the projected result through the same generic
payload renderer as every other success; the projection is applied
to the verb-specific result BEFORE respond.OK, never in the respond
gateway and never per output mode, so the two modes cannot disagree
(0005:C1 held by construction). The projected text lines MUST be a
subset of the default-mode text lines, byte-identical per line.

--plan-only MUST NOT be accepted by any other command — satisfied by
NON-REGISTRATION (registered on resolve only, never the shared
registrars, the flow group's own or persistent set, or the root's
persistent set), pinned by a STRUCTURAL oracle that walks the WHOLE
command tree from the root and asserts the set of commands
registering plan-only is exactly {flow resolve} — closed over verbs
added later, deliberately stronger than the enumerated-sibling sweep
the --all oracle shipped
(internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup)
— vacuity-guarded by requiring the four flow verbs to exist.

WHOLE means total: the walk descends every child of the root with no
name-based skip and no Hidden gate, INCLUDING the auto-generated help
and completion commands. It therefore MUST NOT reuse
internal/cli/help_all.go::walkCommandTree, which skips children named
help or completion before recursing and whose docs.go callers further
gate on Hidden — that walker asserts a strictly weaker negative than
this clause requires. The scope is load-bearing, not pedantic: a
registration reachable only on an auto-generated command is exactly
what an enumerated sweep misses, and the repo already contains one
(help --all, internal/cli/help_all.go::wireHelpSubcommandAll) that
the shipped 0011 oracle does not see because it descends the flow
group only. The behavioural command-error run is corroboration, never
the assertion. This contract mints no new refusal code and no new
exit group.

This clause's IMPLEMENTABILITY is A5, presently Pending: the scope
above is authored as a MUST, but no shipped test walks the tree
totally and the walker is written for the first time at Phase 2. The
clause is stated at full strength deliberately — narrowing it to what
an existing partial walker can assert would discard the cove sweep's
finding — and A5's verification plan is the spike that writes it. If
that spike cannot reach this scope, C1's structural negative needs
its own form and this paragraph moves; every downstream statement of
it (§disposition, §oracle-discriminability, S4, MVV step 5, the desk
trace) inherits that condition rather than restating it.

Totality is a property of the TREE UNDER TEST, not only of the
walker, and the two auto-generated commands are not symmetric: cobra
materializes them lazily, so a walker that never had them to skip
asserts the same weaker negative as one that skips them. A bare
NewRootCmd() tree carries help — force-inited at
internal/cli/help_all.go::InitDefaultHelpCmd — but carries NO
completion command, which cobra creates in InitDefaultCompletionCmd
from ExecuteC. The oracle MUST therefore materialize both
auto-generated commands on the root before walking (calling cobra's
own initializers, never hand-constructing a stand-in), and MUST
assert their presence in the walked set as a precondition. The
control asserting the walk reaches help is NOT sufficient on its own:
help is present by default, so it cannot distinguish a total walker
from an untested tree — completion is the discriminating case and
the control MUST name it.

That precondition is a DEPENDENCY ON COBRA'S OWN INITIALIZERS
(InitDefaultCompletionCmd, reachable from ExecuteC), and the
assertion is what makes the dependency safe rather than silent: if a
cobra upgrade renames, relocates or stops materializing the command,
the precondition fails LOUDLY and the oracle goes red, instead of the
walk quietly covering a smaller tree and passing vacuously. This is
why presence is asserted as a precondition rather than assumed — the
one failure mode a registration-absence oracle cannot survive is
passing over a tree that never contained the command class the clause
exists to reach.
```

**C2**

```normative
The projection axis is the ECHO/PLAN partition doctrine of
JDR 0002 §D1, applied to flow resolve — a fixed partition, not a
caller-supplied field list. The doctrine (one-of-two assignment at
field addition, reflective enforcement, the always-keep rule,
unclear-joins-PLAN, the later-verb conformance rules) is normative
THERE and is cited, not restated; this contract owns the verb's
instance.

The flow resolve assignments: the ECHO group is the model
reference, the observed tags (--tag, echoed unchanged), the
assembled owned view and the invoked reader identities, and the
requested outcome (--outcome, echoed unchanged). The PLAN group is
rule identity, gate results, authored answers and their
interpretations, planned next/writes/clear, the escape disposition,
and revision. revision rides the PLAN side deliberately: it is
produced by the loader, never restated from the request, so it is
structurally a plan-side attestation rather than an echo, and it is
the slot a model-identity signal occupies when one exists.

revision is CONTRACTED BUT PRESENTLY VACANT, and this RDR records
that rather than assuming a live value: 0002's [model] block admits
no revision key, so no model can declare one and the accessor returns
"" on every payload this CLI can emit
(internal/cli/flow_exec.go::revision, REQ-25's consequent, DEV-7).
The assignment is therefore made on the field's DEFINITION — a
loader-produced identity token — not on any provenance it delivers
today; a chained caller cannot currently detect a model change from
it. Projecting it costs 14 bytes and keeps the wire shape stable for
the day 0002 admits the key. If a successor instead retires the
field, that RDR removes it from this group and from the always-keep
core together.
--plan-only reports the PLAN group and omits the ECHO group (C1).

The assignment ground (recorded at Resolve, 2026-08-29): ECHO is
the kernel's input side — supplied by the caller or assembled on
its behalf (internal/resolve/resolve.go::Resolve takes Owned and
Observed as Input; internal/cli/flow_exec.go::invokedReaders
derives readers from model + outcome alone) — and PLAN is the
kernel's output side plus bounded identity attestations, of which
revision is the declared instance: provenance rides as an identity
token, never re-echoed content — the slot is contracted and stable
even while its value is empty. A future drift-detection ask joins
PLAN as a bounded state-identity field under the §D1 field-addition
rule, not by moving raw owned content into the projected width.

Per JDR 0002 §D1's enforcement rule, this verb's reflective oracle
MUST assert every field of the resolve success payload is assigned
to exactly one group, so an unassigned new field is a test failure,
not a silent default — the projection MUST NOT be implemented as a
bare omit-list whose complement is "whatever else exists".

The oracle reflects the payload struct against an assignment source
that is INDEPENDENT of the projection code, and the independence is
the whole content of the check: an oracle that derives the ECHO set
by observing what the projection drops is tautological — it restates
the implementation and cannot fail. So the assignment is DECLARED (a
per-field marker on `resolvePayload` or a table keyed by field name,
one entry per field, carrying `echo` or `plan`), the projection is
implemented FROM that declaration, and the oracle asserts three
things against it: every struct field has exactly one entry; the
entry set and the field set are equal (neither a field without an
entry nor an entry without a field); and the keys a projected run
actually emits equal the declaration's `plan` side. A new field with
no entry then fails at the first assertion rather than defaulting
into either width.

The carrier is CONSTRAINED, not free: it MUST be a standalone table
keyed by field name, in its own declaration site, NOT a per-field
marker on `resolvePayload`. The two carriers are not equivalent under
A2's mechanism — A2 rewrites the echo fields' types and tags on that
struct, so a marker carried there is edited in the very commit that
changes the projection, and an implementer who nils a field and
retypes its tag in one edit has moved the assertion and its subject
together. That is the tautological oracle this clause exists to
forbid, reached by the carrier the clause would otherwise permit.
Independence has to be structural — a separate site a projection edit
does not open — or S3 asserts only that the implementation agrees
with itself.

This verb's ALWAYS-KEEP core (JDR 0002 §D1) is rule, escaped,
escape_class, revision: if the boolean axis ever generalizes to an
enum or a field list, no mode may PROJECT them away, so a projected
payload can never launder a rescued plan into an ordinary one or
detach a plan from the model revision that produced it. Always-keep
is projection-invariance, not unconditional presence: a core field
whose producer already has a presence rule (escape_class) appears
under the flag exactly when it appears by default — the rule the
laundering guard actually rests on is escaped, which is
unconditionally present.

The core's two halves carry different weight today, and the clause
says so rather than implying both are live. rule, escaped and
escape_class are load-bearing now — they are what a caller reads to
tell a rescued plan from an ordinary one. revision's clause is
FORWARD-BINDING: the field is presently constant-empty (above), so
what always-keep buys is that a future projection mode cannot drop
the model-identity slot at the moment 0002 gives it a value. Binding
it now is the cheap half of the trade — 14 bytes against a silent
provenance loss in a mode nobody has designed yet.

How always-keep is enforced today, stated so it is not mistaken for
an unasserted claim: with ONE projection mode, always-keep and the
S2 key-set literal have the same extension — S2 pins all four core
fields as carried, so any change dropping one turns S2 red. That is
enforcement by coincidence of scope, not by an oracle that knows
about the core, and the distinction becomes real the moment a second
mode exists. This RDR therefore does NOT mint a separate always-keep
oracle (it would today assert exactly what S2 asserts); the
obligation it creates is on the successor: an RDR adding a second
projection mode owes an always-keep oracle quantified over MODES —
for every mode, the four core fields are carried — because at that
point S2's single literal no longer covers the claim. Recorded here
rather than in that RDR because this is where the core is declared.

Under JDR 0002 §D1, 0024's dispositions, if that RDR lands, is a
PLAN-group field: each token is assigned by the [emit] declaration
to the selected row's authored emit value (0024:C4's join) —
derived from plan-group inputs only, never from observed or owned —
and its never-omitted clause is untouched by this projection.
```

#### Load-Bearing Decisions

- **Identity** — `--plan-only` is part of *request* identity (the
  same request in the same mode reports the same bytes) but NOT of
  *decision* identity: the same invocation ± the flag selects the same
  rule, runs the same gates, and refuses identically — only report
  width differs. This is the deliberate contrast with `0011:D-identity`,
  where `--all` changes the reported set; here the flag may never
  change what is decided, and an oracle asserts the ±-flag plan-group
  equality.
- **Wire / byte format** — projection is key deletion, never
  re-encoding: the projected payload is the default payload minus the
  echo keys, byte-for-byte on every surviving field, keys in surviving
  declaration order, same non-HTML-escaping encoder. Absent means
  absent — never `null`, `{}`, or `""` stand-ins for a projected key.
  The in-code mechanism (pointer-`omitempty` nilling vs a projected
  struct) is deferred to Resolve (A2); the wire result is fixed here.
- **Naming** — the flag is `--plan-only`, boolean, long form only.
  Precedent for the `<result>-only` fixed projection:
  `gh pr diff --name-only` ("Display only names of changed files",
  `langref/gh-cli/pkg/cmd/pr/diff/diff.go`), mirroring
  `git diff --name-only`; and the verb's own help already calls the
  kept group "a successful plan" (`flowResolveExtendedDesc`), so the
  name states the result in house vocabulary. Rejected: `--select`
  (a `0011:D-naming` named-rejected spelling; even though its pin
  test is scoped to `next`, reusing the spelling for a different axis
  would make the two rejections ambiguous); `--json <fields>` /
  `--fields` (imports the field-universe policing of gh's general
  form — Alternative 1); `--quiet` (peer meaning is "suppress output"
  or "ids only" — `grep -q` emits nothing — not "answer without
  echo"); `--no-echo` (names the mechanism, not the result — the same
  test `0011:D-naming` applied to `--ignore-match`); `--short` /
  `--brief` (no fixed partition meaning in any surveyed peer); an
  inverted default (`--explain` to get the echo) — the defect-is-the-
  default inversion: it would put a flag on every one-call consumer
  to preserve behaviour that is already correct.
- **Selection / predicate** — which side a field joins is decided by
  the partition doctrine at JDR 0002 §D1 (cited, not restated:
  echo/plan assignment at field addition, unclear → plan); C2
  carries this verb's assignments. The reflective
  partition-completeness oracle (C2) makes skipping that assignment a
  test failure rather than a silent omit-list default.

#### Source-authority census

Every success-payload field, its writer, and its C2 side. The set is
FOURTEEN fields — `reflect.TypeOf(resolvePayload{}).NumField()` is 14,
the count `internal/cli/decision_table_0010_test.go` already pins — laid
out in thirteen rows below because `writes` and `clear` share one
`plan.Writes` writer. All fourteen are written at one assembly site
(`internal/cli/flow_resolve.go` `resolvePayload{…}` literal) — there is
no fallback arm and no second writer, which is what makes the partition
a single-site edit. The three numbers must agree: this census is the
input to C2's declaration and to S3's bidirectional equality assertion,
so a miscount here becomes a red S3 on the implementer's first run.

| Field | Writer | Derived from | C2 side |
| --- | --- | --- | --- |
| `model` | `req.modelRef` | `--model`, verbatim | ECHO |
| `observed` | `observedTagMap(req.observed)` | `--tag`, refused-not-coerced (`parseTags`) | ECHO |
| `owned` | `tagMap(owned)` | reader pass over the caller's own request | ECHO |
| `readers` | `readerIDs(readers)` | pure function of model + outcome (`flow_exec.go::invokedReaders`) | ECHO |
| `outcome` | `outcome` | `--outcome`, verbatim | ECHO |
| `revision` | `req.revision()` — constant `""` today (`flow_exec.go`, no `[model]` revision key exists) | loader, not the request | PLAN (identity slot) |
| `rule` | `plan.RuleID` | kernel | PLAN |
| `gates` | `req.runGates` | gate run on the selected row | PLAN |
| `emit` | `emitMap(row.Emit)` | selected row's authored block (`0010:C4`) | PLAN |
| `next` | `tagMap(plan.NextTags)` | kernel | PLAN |
| `writes`, `clear` | `plan.Writes` split on `ClearSentinel` | kernel | PLAN |
| `escaped` | `plan.Escaped` | kernel | PLAN |
| `escape_class` | `escapeClassOf(...)` | re-probe of the ordinary rows | PLAN |

Sibling arms: none. The reflective oracle (C2) is what keeps this table
from going stale — a field added without a side is a test failure.

#### Disposition

| Input class | Exit | Envelope | Artifact | Loud / silent |
| --- | --- | --- | --- | --- |
| `--plan-only` on `flow resolve`, success | 0 | projected success payload | — | loud (narrower output is the request) |
| `--plan-only` on `flow resolve`, refusal | unchanged | refusal, byte-identical (A6) | — | loud, flag-blind |
| `--plan-only` on `next`/`read-state`/`set-state` | 2 | `command-error` shared bucket | — | loud (parse failure, `0011:A14` limit) |
| `--plan-only` on a binary predating this RDR | 2 | `command-error` | — | loud — an absent `observed` always means projection, never an old binary |
| echo group under the flag | — | keys ABSENT | — | silent by design; the caller supplied every one |

The flag mints no refusal code and no exit group (C1).

#### Oracle discriminability

Each Testing Strategy scenario, what makes it fail, and its negative
control — the guard against an oracle that passes by absence-of-error.

| Oracle | Fails if X is wrong, because Y | Negative control |
| --- | --- | --- |
| S1 ± flag differential | a projected value drifts from its default rendering — the comparison is against the same run's default bytes, not a literal; the width clause additionally fails if the projection stops saving bytes on the full emitted line; the reader-set check observes which readers actually RAN, so it fails on a real execution change rather than on the field's absence | mutate one carried value under the flag; separately, add a large plan-group field and confirm the width clause goes red where the key-set oracle stays green; separately — the control that discriminates this clause — SKIP the reader pass under the flag and confirm the reader assertion goes red. That control is the reason the measurand is observed execution: a recomputed `invokedReaders(model, outcome)` comparison stays green through it (neither argument is flag-dependent), and reading the set from the projected payload aborts instead |
| S2 explicit key set + absent-not-null | the key list is an authored literal, so a field falling through into the projected width fails it; and a projected key rendered as `null`/`{}`/`""` fails the absence assertion rather than passing as "not carried" | add a field to the payload without a C2 side; S2 and S3 must both go red. Absence control: render one echo key as `null` (the bare non-pointer `omitempty` hazard A2 reproduced) — S2 must go red where a key-presence-only check would stay green |
| S3 partition completeness | reflective over the struct against an INDEPENDENT declared assignment source, so an unassigned field is a failure rather than a default, and the check cannot restate the projection code | add a field to the payload without a C2 side; S2 and S3 must both go red. Discriminating control (not shared with S2): add a field, declare it `plan`, but omit it from the projection — S3 must go red on the emitted-keys-equal-`plan` assertion while S2 stays green |
| S4 registration walk | vacuity-guarded on the four flow verbs AND on the walked set containing `help` and `completion`, so neither an empty walk nor an unmaterialized tree can pass | register `plan-only` on a second command; S4 must go red. Second control: register it on the auto-generated `completion` command — S4 must go red there too, which is the control a bare `NewRootCmd()` tree cannot run |
| S5 text subset | asserts per-line byte identity against the default run's lines, not merely "fewer lines", plus a strict text-width reduction so the equal set cannot pass | drop a plan-group line under the flag; S5 must go red. Second control: make the text projection a no-op (project nothing) — the subset assertion stays green and only the width assertion goes red |

S4 is the one structural-absence oracle here, and the vacuity guard is
what stops it passing on an empty command tree.

#### Illustrative Code

Illustrative only — tests must not assert these bytes; normative
fixtures land at Resolve.

```sh
# Chained second call: the caller already holds the 48 facts it passed.
intrastate flow resolve --model pricing.toml \
  --outcome decide --tag tier=free --tag region=eu \
  --plan-only --as=json
```

```json
{"revision":"","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

The same run without `--plan-only` additionally carries `model`,
`observed`, `owned`, `readers`, and `outcome`, exactly as
`docs/cli-output-contract.md`'s worked payload shows today.

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Success payload assembly | `internal/cli/flow_resolve.go::resolvePayload` | Fields are concrete (non-pointer), so conditional absence needs a mechanism (A2) | Extend | Echo fields become projectable; plan fields untouched |
| Both-modes rendering | `internal/cli/respond` (`OK`, `text.go::flatten`) | None — renders whatever the verb hands it | Reuse unchanged | Projection lands before `respond.OK`, so no gateway change |
| Flag registration | `newFlowResolveCmd` (verb-local), `registerSelectionFlags` (shared) | Shared registrar reaches all verbs | Extend verb-local only | The flag never enters the shared registrar (C1) |
| Absence-oracle idiom | `internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` | Enumerated two-level sweep over the `flow` group, one flag name per sweep; no shipped test recurses | Reuse as pattern, strengthened | New whole-tree oracle: commands registering `plan-only` == exactly `{flow resolve}` (C1, A5) |
| Recursive tree walk | `internal/cli/help_all.go::walkCommandTree` | Skips `help`/`completion` by name; callers gate on `Hidden` — not total | Do NOT reuse; write a total walker | C1 fixes the scope: no name-skip, no `Hidden` gate (A5 spike) |
| Verb help | `flow_resolve.go::flowResolveExtendedDesc` ("Reading a successful plan") | None | Extend | The help's plan list is the partition's user-facing statement; gains the flag line |

### Decision Rationale

The routable question — should the resolve payload have a
caller-controlled projection axis, and where does it live — was scored
as a matrix (profile: foundational), approaches × deciding criteria:

| Criterion | O1 fixed plan-projection flag | O2 named-field selection | O3 gateway verbosity tier | O4 no CLI axis (caller-side jq) |
| --- | --- | --- | --- | --- |
| Correctness fit (drop the chained-call echo; refusals untouched) | full — the echo group is exactly the measured duplication | full, plus generality nothing asked for | partial — omission differs by mode | partial — needs a per-skill recipe carried in every prompt |
| Prior-art alignment | `gh pr diff --name-only` (fixed projection) | `gh --json <fields>` (field selection + universe policing) | none surveyed — gh renders text from the same selected fields, not a tier | jq is gh's *escape hatch* (`--jq`), offered beside, not instead of, projection |
| Safety of the projected payload | `escaped`/`escape_class` never PROJECTABLE, by construction | caller can drop `escaped` unless an always-keep set is added — which IS this partition | n/a | a recipe can drop anything, including refusal fields on the shared stdout |
| Blast radius vs locked fences | one verb, one flag; `0005:C1` held by construction; no shipped oracle moves | same, plus a declared field universe and adjacency to the rejected `--select` spelling | breaks `0005:C1`'s one-`flatten` gateway or mints per-verb text templates | zero repo change; the burden exports to every consumer |
| Consumer parse stability | one struct parses both widths (absent echo keys decode to zero values) | N shapes, one per field list | two shapes per mode | consumer-defined, drifts as the payload grows (0024 adds a field) |
| Extensibility (a later verb, e.g. an export surface) | one rule: project the verb result pre-`respond.OK`, echo group only | generalizes, dragging the field universe along | — | — |
| Cost / reversibility | one flag + oracles; removable by deregistration | exporter machinery + per-verb field lists | per-verb templates, the cost 0011 already declined | zero now; N-skill drift later |

O1 wins on the safety, fence, and stability rows; the correctness row
is a tie between O1 and O2 and the extensibility row does not need
O2's generality (no caller asked for sub-plan selection — the ask is
"answer without my own input back"). O2 is the strongest rejected
alternative and is written up in full; O3 falls to the shipped
gateway's one-marshal structure, O4 to exported drift.

Premortem consequence (hardened critic, `claude-opus-5[1m]`, one
query; ledger P-1…P-18 in
`docs/rdr/0023-resolve-envelope-projection/evidence/propose-premortem/critic.md`):
the recommendation survived; no finding forced a switch. Its
mitigations are live in the sections that own them — C1's
reflective partition-completeness, explicit projected-key-set,
whole-tree registration, and differential report-only oracles
(P-1/P-5/P-9/P-12/P-13); C2's always-keep core and plan-side
`revision` (P-3/P-17); A1's shipped-partition scope across two
fixture shapes (P-7/P-8); A4's quoted join inputs (P-15); the
inverted-default rejection on its own ground (P-16); the
refusal-width scope and version-skew loudness under
Trade-offs/Failure Modes (P-6/P-11); P-14's generic-bucket message
limit accepted as `0011:A14`'s own. One finding was refuted by
source: P-3's claim that `outcome` is produced rather than echoed —
`runFlowResolve` copies the `--outcome` flag verbatim, so `outcome`
stays echo; A7 pins that read.

Joint-decision check — open peers at depth 1 with Status Draft/Final:
0012–0022 and 0024 (13, all Draft; 0001–0011 are Implemented). FIRE
on 0024: shared modify-anchor
`internal/cli/flow_resolve.go::resolvePayload` (0024:C4 appends
`dispositions`; this RDR projects the echo group off the same
payload). Homed, not paused: `JDR 0002 §D1` (the chartered envelope
registry) is the normative home; C2 carries the verb's instance —
under the doctrine `dispositions` is a
never-projected plan-group field, and
0024:C4's quoted `Plan.RuleID` join reads plan-group inputs only, so
the two compose in either landing order with neither contract
moving; the same fire and home are written symmetrically into
0024's Joint-check line. Context beside the fire: 0012 touches the
same file at a disjoint symbol (`flow_resolve.go::guardSeam`);
0014/0017/0019/0020 cite `docs/cli-output-contract.md`, which this
RDR's Phase 3 also edits, in disjoint sections (0017 writes its
findings/code-registry sections; this RDR adds the projected worked
payload beside the resolve section) — co-citation context, no shared
decision. Roster peer 0021 (unproposed, later in order) would be a
consumer of JDR 0002 §D1's later-verb clause if its export surface
ever wants a projection; that clause obliges it to nothing. Absence arm: n/a —
this proposal converts no refusal into an acceptance (`--plan-only`
on the other verbs stays refused; success reporting only narrows).
Bridge sub-check: n/a — no Cluster membership; no plan here
introduces a surface a sibling schedules for deletion (no
bridge/retire/replace cues).

Premortem: hardened (hardened)

Ground-sweep: clean (24 anchors)

Joint-check: fired → 0024 (home: JDR 0002 §D1)

## Alternatives Considered

### Alternative 1: Named-field selection (`--json <fields>` form)

**Description**: The caller names the fields to keep —
`--fields emit,rule` — against a per-command declared field universe,
the general form `gh` ships through one funnel
(`langref/gh-cli/pkg/cmdutil/json_flags.go::addJsonFlag`:
`f.StringSlice("json", nil, "Output JSON with the specified fields")`,
wired per command with a command-owned field list).

**Pros**:

- Strictly more expressive; any future narrowing ask is already
  answered.
- The strongest peer precedent by adoption (`gh`'s entire list/view
  surface).

**Cons**:

- Unsafe without an always-keep set: a caller who omits `escaped`
  reads a rescued plan as an ordinary one — and adding the always-keep
  set reduces the design to C2's partition with extra machinery on
  top.
- Imports the field-universe obligation: a declared pickable-field
  list per verb, validation, and its own refusal for an unknown field
  — surface the ask ("answer without my echo") never needed.
- N payload shapes instead of two; every consumer struct is
  per-field-list.
- Adjacent to `--select`, a `0011:D-naming` named-rejected spelling;
  in `gh` the same flag also switches output mode, a duty `--as`
  already owns here.

**Reason for rejection**: all of its extra generality is cost —
the one measured need is the echo group, and the safe version of this
alternative contains the chosen design as its mandatory core.

### Alternative 2: Verbosity tier in the output gateway

**Description**: `--as json` omits `observed` (or the echo group)
unless `--explain`/`--verbose` restores it — a mode/tier decided in
`internal/cli/respond` rather than in the verb.

**Pros**:

- No new per-verb flag; one policy for every verb at once.

**Cons**:

- A json-only omission is impossible in the shipped gateway without
  restructuring it: `respond/text.go::writeTextPayload` marshals the
  same verb result once for both modes through one `flatten`, and
  `0005:C1` requires text "derived from the same verb-specific
  result" — so the tier needs either a restructured gateway or the
  per-verb text template 0011 already declined as "editing a
  predecessor's surface for cosmetics".
- The gateway is verb-agnostic; a projection policy there must know
  each verb's partition, importing per-verb knowledge into the one
  layer that has none today.
- A gateway-level policy is exactly the silent surface-widening the
  shipped absence oracles exist to catch (a persistent flag reaching
  all verbs at once).

**Reason for rejection**: everything a gateway tier could share is
already shared by C2's two rules (project before `respond.OK`, echo
group only) without restructuring the gateway or teaching it
per-verb partitions; the residual difference is only where the
partition knowledge lives, and it belongs to the verb that owns the
payload.

### Alternative 3: No CLI axis — caller-side projection (jq)

**Description**: Ship nothing; each consumer pipes
`... --as=json | jq '{rule, emit}'` when it wants the narrow form.

**Pros**:

- Zero repo change; zero new contract on a locked verb.

**Cons**:

- The recipe is a second description of the payload carried in every
  consumer skill prompt, drifting as the payload grows (0024 adds
  `dispositions`); the tokens the recipe and its maintenance cost are
  the currency this RDR is trying to save.
- Refusals share stdout as the same single JSON line, so a naive
  projection corrupts the CLIError envelope into `{"rule":null,...}`;
  a safe recipe must branch on envelope type — per skill, forever.
- Exports the safety problem: nothing stops a recipe from dropping
  `escaped`.

**Reason for rejection**: relocates the cost from one audited
contract to N unaudited copies of it, and pays tokens to save tokens.

### Briefly Rejected

- **Inverted default** (echo only under `--explain`): rejected on
  its own ground, not on the settled-consumer citation alone (the
  premortem's seed-3(b) point — that consumer's decision speaks to a
  rendering consumer, not a chained one): one-call consumers are the
  common case and should carry no flag, and an opt-out is reversible
  where a default flip is a breaking change to every existing
  caller.
- **Config/env-controlled projection** (`INTRASTATE_PLAN_ONLY=1`):
  report width must be request identity, visible in the invocation —
  ambient config makes two identical commands report differently.
- **Projecting refusal envelopes too**: refusals are already minimal
  and every field is diagnosis; nothing measured motivates it.
- **A separate verb** (`flow answer`): duplicates `resolve`'s whole
  decision surface to change report width — the `flow candidates`
  rejection (`0011:BR2`) transposed.
- **Registering the flag but having other verbs refuse it with a
  typed code**: C1 follows `0011:C2`'s settled disposition —
  non-registration with the shared `command-error` bucket — and mints
  no new refusal code.

## Context

### Background

Filed while shrinking a consumer's navigator skill; tracked as kata
`intrastate#srz2`, raised to P1 in the 2026-08-27 garden pass because
the transcript token, not the millisecond, is the cost every skill run
pays. RDR 0011's fences are narrower than a first read suggests:
it rejected `--select` as
a named spelling with a pin test **on `flow next`**
(`internal/cli/flow_all_0011_test.go::TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand`),
and its A14/C2 structural absence oracle fences `--all` specifically —
not any flag — on the other three verbs. So no shipped oracle blocks a
new flag on `resolve`; what this RDR reopens deliberately is 0011's
rejection *ground* ("the defect is the default"), which does not bind
a surface whose default is agreed correct. No projection precedent
exists to reuse:
`--all` runs the opposite direction (widens `next`'s row set), and the
related `excluded_by`-on-`--why` idea was Briefly Rejected in 0011
with the note that reopening must be "deliberate rather than
rediscovered" — the same disposition this RDR discharges. Sibling
kata `rg0e` (seeded in this batch) adds *interpretation* to the same
envelope where this seed *removes* fields — independent asks touching
different locked clauses, deliberately not merged. Sequencing: land
before or with the consumer migrations that add more models per skill
run.

### Technical Environment

Go module `github.com/cwensel/intrastate`. Surfaces:
`internal/cli/respond/text.go` (the `flatten` gateway — "the two
modes cannot disagree about what the run reported"),
`internal/cli/flow_all_0011_test.go` (structural absence oracles for
the flag fence). Governing records: RDR 0005 (C1, minimum data
shape), RDR 0011 (C2, A14, D-naming).

## Research Findings

### Investigation

Prior art was read before enumerating (bounded pass recorded, with
queries and rejected branches, in
`docs/rdr/0023-resolve-envelope-projection/evidence/research/propose-prior-art.md`):
peer-CLI output reduction over the `langref` checkout set the
resources index names, and the governing fences via the projector
(`0005:C1`, `0005:§technical-design` "Envelope", `0011:C2`,
`0011:D-naming`, `0024:C4`). The shipped surfaces were read directly:
`internal/cli/flow_resolve.go` (payload assembly, extended help),
`internal/cli/respond/text.go` (the one-`flatten` gateway),
`internal/cli/flow_all_0011_test.go` (the absence-oracle idiom and
the scope of the `--select` pin). ⚠ no prior-art coverage was found
for the decision-engine instance question ("does a DMN-class engine
echo inputs in its evaluation result") — two `StateMachineRes`
queries missed; the choice does not rest on it.

### Key Discoveries

- **Documented** — `gh` offers caller-controlled output reduction in
  two forms through one funnel: named-field selection
  (`pkg/cmdutil/json_flags.go::addJsonFlag`, "Output JSON with the
  specified `fields`") and jq escape hatch; and, separately, the
  boolean fixed projection `gh pr diff --name-only` ("Display only
  names of changed files", `pkg/cmd/pr/diff/diff.go`). ⇒ both
  candidate shapes have live peer precedent; the boolean form carries
  no field-universe obligation.
- **Documented** — `0011:C2`'s fence and its shipped oracles pin
  `--all`'s absence by NAME on the other three verbs
  (`flow_all_0011_test.go::TestReq46And47And65And98_...`), and the
  rejected-spelling pin (`--select` among them) sweeps `flow next`'s
  flag set only (`TestReq36And50_...`). ⇒ no shipped oracle blocks a
  new flag on `resolve`; the reopening is of a rejection *ground*,
  not of a fence.
- **Documented** — `respond/text.go::writeTextPayload` marshals the
  verb result once and `flatten`s it for text, the same encoder as
  the wire. ⇒ any projection must land before `respond.OK`; placed
  there, `0005:C1`'s both-modes agreement is free.
- **Documented** — the verb's own help
  (`flow_resolve.go::flowResolveExtendedDesc`, "Reading a successful
  plan") already partitions the payload: it lists `rule`, `next{}`,
  `writes{}`, `clear[]`, `emit{}`, `escaped`, `escape_class`,
  `gates[]` and none of the echo fields. ⇒ C2's partition is the
  house vocabulary made normative, not a new taxonomy — plus
  `revision` on the plan side as the identity slot (C2 — contracted,
  presently empty). The
  writer of that state is `runFlowResolve` itself: every plan-group
  field is computed by the selection/gate path, every echo-group
  field is copied from the request or the reader pass.
- **Documented** — the only conditional success-payload fields
  shipped are producer presence rules (`resolvePayload.EscapeClass`
  `omitempty`; `flow next` candidate `gates,omitempty` under
  `--evaluate-gates`, whose writer is the `--evaluate-gates` gate
  run). ⇒ no caller-controlled projection exists to reuse (sibling-
  path check: searched, none exists).
- **Assumed** — the seed's measurement (1529 → ~430 bytes on the
  48-fact call, ~72% echo share) transfers to the shipped flag (A1),
  and the encoder mechanics of conditional absence (A2) hold; both
  are Resolve's to verify.

## Trade-offs

### Consequences

- Positive: a chained or secondary `resolve` call stops re-paying the
  echo as transcript tokens. The seed's ~72% is RETIRED as a headline
  (Problem Statement): it was the echo *share* of the seed's own
  motivating call measured against `emit`+`rule` alone, and that call
  was never re-run under the shipped partition. The shipped figures
  are A1's: 42.2–51.3% on the three checked-in fixtures and 79.4% on
  a synthetic 48-fact shape authored for the spike and not checked in.
  Quote the range with its shape, never a single number — the saving
  scales with fact count, the floor being the retained PLAN group
  (`gates`, `next`, `writes`, `clear`, `escaped`, `revision`). The
  saving multiplies with every model a consumer migrates.
- Scope of the guarantee: the token reduction covers successes only.
  Refusal envelopes are untouched — they carry diagnosis
  (`findings`), not echo (A6) — so a probing caller that mostly
  refuses saves nothing, by design. Unsized here: what share of the
  motivating consumer's traffic refuses. The expected value scales
  with that mix and this RDR does not measure it; the flag is
  strictly non-negative on every call either way, so the mix changes
  how much is saved, never whether the change is worth landing.
- Scope of the deliverable: this RDR ships the CAPABILITY, not its
  adoption. The flag is opt-in and no consumer passes it on landing,
  so a green build delivers zero measured saving on day one — every
  oracle can pass while the motivating cost is still being paid. That
  is the correct boundary (a projection flag and a consumer's
  call-site migration are separate changes with separate blast
  radii), but it means "done" here is deliberately weaker than
  "benefit realized": Phase 2 green plus the MVV step-6 table is the
  completion bar, and the saving is realized only when a caller adds
  the flag. Adoption is not scheduled by this RDR and is not one of
  its phases; the seam-mate ask (kata `rg0e`) and the consumer that
  dropped its status call are the natural first adopters, and
  whoever migrates them owns measuring the realized reduction against
  A1's table.
  Consequence for the motivating tracker, stated so it is not
  discovered at close: `intrastate#srz2` (P1) is the ADOPTION ask, so
  this RDR landing green does not close it — it unblocks it. Closing
  it needs a follow-on migration change that adds `--plan-only` at a
  real call site and records the realized reduction. That successor
  is not seeded here (it is a consumer-side change, not this RDR's
  surface), and the RDR's own completion bar stays Phase 2 green plus
  the MVV step-6 table; what changes is only that nobody should read
  a green build as having delivered the P1 saving.
- Positive: the default is untouched — no consumer changes, no 0005/
  0010/0011 test moves (A3), no respond-gateway change.
- Negative: `flow resolve` now has two success-payload widths; docs,
  help, and any payload-shape oracle must say which width they speak
  of.
- Negative: the partition (JDR 0002 §D1; this verb's instance in C2)
  is a standing obligation — every future payload field must be
  assigned a side by the RDR adding it, and a mis-assignment is a
  spec defect, not a style choice.
- Accepted asymmetry: `--plan-only` on the other three verbs is the
  shared `command-error` usage bucket, not a typed refusal — the
  same accepted limit `0011:A14` recorded for `--all`.

### Risks and Mitigations

- **Risk**: partition drift — a successor appends a payload field
  without assigning it a side, and an omit-list implementation
  defaults it into the projected width silently (the premortem's
  P-12: a large derived field could make `--plan-only` output
  *bigger* than today's full payload with every test green).
  **Mitigation**: C2 forbids the bare omit-list and mandates the
  reflective partition-completeness oracle plus the explicit
  projected-key-set literal — an unassigned field is a test failure,
  and a field reaching the projected width is a conscious edit to a
  normative list. Those two catch an unassigned field but not a
  legitimately-assigned large one, so C1's differential also asserts
  the projected encoding is strictly shorter than the default — the
  clause that actually fails on P-12's "bigger output, every test
  green" scenario.
- **Risk**: the boolean later generalizes (a second projection
  profile is asked for) and drifts toward the rejected field-list
  design without its safety.
  **Mitigation**: C2's always-keep core (`rule`, `escaped`,
  `escape_class`, `revision`) binds every FUTURE projection mode
  now, so the safety invariant survives any widening of the axis.
- **Risk**: mechanism drift — a projected-struct implementation (A2)
  duplicates `resolvePayload` and the two fall out of sync when a
  field lands (e.g. 0024's `dispositions`).
  **Mitigation**: the ±-flag oracle asserts plan-group byte equality
  between widths on every run, so a field added to one struct only is
  caught by the suite, not by a consumer.
- **Risk**: a consumer treats plan-only output as the full record and
  loses the model reference — `model` is the projected-away field that
  actually names which model produced the plan, since `revision` is
  constant-empty today (C2).
  **Mitigation**: opt-in with an unchanged default; the help line
  states the omitted group by name; the caller that opts in is the
  caller that already holds the request, `model` included — which is
  precisely why `model` is ECHO and not PLAN.
  **Residual, accepted and named rather than mitigated away**: "holds
  the request" is an assumption about the caller, and the motivating
  caller is the case that strains it — an agent projects precisely
  because it is dropping transcript, so a chain that keeps only
  projected bodies has no `model` and an empty `revision`, and cannot
  answer "which model decided this" from the record alone. The
  partition is not reopened for it: `model` is a verbatim `--model`
  copy (A7), so promoting it to PLAN would put echoed input in the
  decision group and contradict JDR 0002 §D1's ground, and `revision`
  is named in that JDR's always-keep core at decision level, so this
  RDR cannot drop it. What the residual actually costs is bounded by
  the always-keep core: a projected plan can never be MISREAD (`rule`,
  `escaped`, `escape_class` are all carried) — it can only be
  unattributed. A consumer needing attribution in the projected width
  either records the model on its own side at call time, or waits for
  `revision` to carry a value (C2's forward-binding clause), and an
  RDR that gives `revision` a live value closes this without moving a
  field.

### Failure Modes

- Visible: `--plan-only` on `next`/`read-state`/`set-state` fails the
  parse — `command-error`, exit 2, flag name only in pflag's free
  text (the `0011:A14` shared-bucket limit, accepted again here).
- Visible: a projected payload missing a plan-group field (mechanism
  bug) fails the ±-flag equality oracle and any consumer struct
  relying on `rule`/`emit`.
- Silent-by-design, guarded: a projection accidentally applied to a
  refusal path would change refusal bytes — A6 plus a
  refusal-unchanged oracle make it a test failure, not a field
  report.
- Silent-by-consumer-choice, accepted: a consumer struct or jq
  recipe written against the FULL width and fed projected output
  reads zero values for the echo fields. Accepted because the flag
  is opt-in on the caller's own invocation — the caller that sets it
  is the caller that owns the parse — and the always-keep core means
  no decision field is ever among the zeros.
- Version skew is loud, not silent: a binary predating this RDR
  refuses `--plan-only` at the parse (`command-error`, exit 2), so a
  caller can never hold a full-width payload believing it was
  projected, and an absent `observed` always means projection, not
  an old binary.
- Diagnosis: `--as=json` ± `--plan-only` over the same request
  diffs to exactly the echo keys; any other diff is a defect in the
  projection.
- Recovery: drop the flag — the default width is the full record and
  is contractually unchanged.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified, with ONE carried exception:
      A5's implementability half is discharged by writing the total
      walker itself, which is Phase 2 work. It is a Pending
      assumption whose verification IS a build step, not a blocker on
      starting — Phase 1 does not depend on it. What it does gate is
      lock: if the walker cannot be written to C1's scope (whole
      tree, both auto-generated commands materialized and present in
      the walked set), C1's structural negative needs its own form
      and the contract moves. Every other assumption is Verified
      before Phase 1.
- [ ] Ordering tolerance with cli/0024 confirmed (A4): either RDR may
      land first; the second lands with `dispositions` already/newly
      in the plan group and no contract in either moves.

### Minimum Viable Validation

Over the checked-in decision-table fixture
(`models/examples/pricing-decision-table.toml` or the 0005 fixture
family):

1. Run `flow resolve` with a recognized outcome and discriminating
   `--tag`s, `--as=json`, without the flag → today's full payload,
   byte-identical to the pre-change fixture. The fixture is a GOLDEN
   CAPTURED FROM THE PRE-CHANGE BINARY and checked in BEFORE Phase 1
   edits the struct (Phase 1's first step, `git stash`-clean tree) —
   this is the only comparison in the whole battery that has a
   pre-change side. S1–S5 all compare the new build against itself,
   so none of them can see a default-mode regression that moves both
   sides together; without a captured golden, C1's headline
   "byte-identical to today's" is asserted nowhere and every existing
   consumer rides on a manual step.
2. Re-run the same invocation with `--plan-only` → exactly the
   normative projected key set (`revision`, `rule`, `gates`, `emit`,
   `next`, `writes`, `clear`, `escaped` — plus `escape_class` iff
   step 1 carried it); every carried field byte-identical to step 1's;
   `model`/`observed`/`owned`/`readers`/`outcome` absent (not null,
   not empty).
3. Re-run step 2 with `--as=text` → the projected lines are a
   byte-identical subset of step 1's text lines, stable across
   repeated runs.
4. Re-run with an unrecognized outcome, ± the flag → byte-identical
   refusal envelopes and exit codes.
5. Run `flow next --plan-only` → `command-error`, exit 2; and the
   whole-tree structural oracle (the set of commands registering the
   flag is exactly `{flow resolve}`, over a root with the
   auto-generated commands materialized and present in the walked
   set) passes. One sibling verb is run, not all three, deliberately:
   C1 makes the behavioural run CORROBORATION and the structural walk
   the assertion — the walk already covers `read-state`, `set-state`
   and every command added later, so running the other two would add
   no coverage the walk does not have.
6. Record default vs projected byte counts on the motivating-model
   shape and one gate/write-heavy fixture, measured on the full
   emitted line (C1's unit), and compare against A1's baseline in
   `evidence/spikes/a1-byte-width.md` — 79.4% saved on the 48-fact
   class, 42–51% on the state-machine shapes.

   PASS BAR, stated on shapes the implementer can actually run: every
   shape measured saves bytes (the S1 width oracle already enforces
   this per-run), and the CHECKED-IN fixtures
   (`models/examples/pricing-decision-table.toml`,
   `release-grammar.toml`, `review-state-machine.toml`) each save at
   least 40% — A1 measured 42.2–51.3% on them, so a landing below 40%
   there means the plan group is carrying materially more than A1
   measured, which is A1's "saves too little to justify a new
   surface" condition and routes back rather than recording a number.

   The 79.4% figure is NOT a pass bar, because its shape is not
   checked in: `pricing-48.toml` was authored inside the A1 spike as
   scratch (`a1-byte-width.md` §S1b) and no fixture of that width
   exists in `models/examples/`. It is recorded as the high-fact-count
   DEMONSTRATION — the saving scales with fact count by construction,
   the floor being the retained PLAN group — and any headline quoting
   79.4% must say which shape produced it. To gate on it, the step
   re-materializes the 48-fact model from the spike's recorded
   generator first; otherwise the checked-in bar above is the gate.

End-state: one flag, two widths, one decision; refusals and default
mode untouched.

### Phase 1: Projection on the verb

Register `--plan-only` on `resolve` only and hand `respond.OK` the
projected verb result when set — the mechanism A2 verified, applied
after payload assembly, before the gateway.

The projection site sits on the SUCCESS path only, after the last
`respond.Fail` return: refusal flag-blindness (A6, C1's report-only
clause) is then structural — a refusal returns before the projection
is reachable — and NOT a defensive branch. Do not add a "if refusing,
skip projection" guard; a guard would mean the projection site is
wrongly placed, and it would make the flag readable on a path C1
requires it cannot influence. The flag is read once, at the
projection site, and nowhere in gate evaluation, rule selection, or
refusal construction — which is what makes S1's identical-decision
assertions hold by construction rather than by test.

The `resolvePayload` struct keeps exactly its current field count:
`decision_table_0010_test.go` pins `NumField() == 14` (:435) and the
13-key wire list (:422-426), so A2's pointer conversion changes field
TYPES only. Adding a field here would falsify A3's "no predecessor
oracle moves" claim, so a field addition is out of scope for this
RDR by construction, not by preference.

### Phase 2: Oracles

Five oracles, one per Testing Strategy scenario S1–S5, and the
enumeration here is that list — not a sixth: the whole-tree
registration oracle (S4: commands registering the flag ==
exactly `{flow resolve}`, vacuity-guarded, over a tree with the
auto-generated commands materialized); the differential
report-only oracle (S1: ± flag — exit codes, refusal bytes, strict
key-subset with byte-identical carried values, invoked-reader set
via the projection-independent channel, full-line width); the
reflective partition-completeness oracle (S3: every payload field in
exactly one C2 group, against the declared assignment source); the
explicit projected-key-set oracle (S2), which CARRIES the
absent-not-null assertion — a projected key must be absent, never
`null` and never an empty placeholder, the live hazard A2's spike
reproduced with bare non-pointer `omitempty`; and the text-subset
oracle (S5, including the text width assertion).

### Phase 3: Docs and help

`docs/cli-output-contract.md` gains the projected worked payload
beside the full one; `flowResolveExtendedDesc` gains the flag under
"Reading a successful plan"; the partition statement (C2) lands where
the payload fields are documented.

The help text is not only prose: `flowResolveExtendedDesc` and the
flag's own usage string are inputs to the GENERATED artifacts
`docs/cli-reference.md` and `llms.txt`, which `make docs-check`
verifies as up to date and which `make check` runs (`Makefile`,
`DOCS_FILES` / `docs-check` / `check`). So this phase regenerates
those files and commits them in the same change. The flag's usage
string is authored here, once, since it ships into a committed
artifact rather than staying an implementation detail.

Phase ordering, stated as a constraint rather than a hazard: because
the generated artifacts derive from the flag's registration, the
commit that registers the flag MUST also carry the regenerated
`docs/cli-reference.md` and `llms.txt`, or `make check` fails on
`docs-check` at that commit. "Phase 3" names the authoring work, not
a later commit boundary — regeneration rides with Phase 1's
registration commit, and the prose/worked-payload edits may follow
separately. A phase split that leaves the flag registered and the
artifacts stale is a red build by construction, so the plan does not
schedule one.

## Validation

### Testing Strategy

The Phase 2 oracle battery is the strategy; done = every oracle
green and the MVV run recorded. Every C1 MUST maps to at least one
oracle below or an MVV step; A1's byte table is recorded evidence,
not a test assertion.

1. **Scenario**: same request ± `--plan-only`, `--as=json`
   (differential report-only oracle).
   **Expected**: identical exit codes; identical invoked-reader set,
   measured as OBSERVED READER EXECUTION per C1 (the default run
   supplies the expected set; the projected run's set is recorded at
   the reader invocation site, never from its own payload —
   `readers` is projected away and
   `flow_fixtures_0011_test.go::readersOf` fails hard on absence —
   and never by re-calling `flow_exec.go::invokedReaders`, whose
   arguments the flag does not touch, which would make the assertion
   unfailable);
   refusal envelopes byte-identical; the projected success payload a
   strict key-subset of the default with byte-identical values on
   every carried key; and the projected encoding strictly shorter
   than the default, measured on the FULL EMITTED LINE including the
   `{"type":"ok","data":{…}}` envelope (C1's width clause and its
   unit — the guard P-12 otherwise escapes; unconditional because the
   five echo keys always render, `model`/`outcome` being
   non-`omitempty` and the three containers rendering `{}`/`[]`).
   Normative fixture: the projected wire record of the pricing 2×2
   call, whose exact bytes are the "Projected reference" line in
   `evidence/spikes/a2-encoder-mechanism.md` §Reference output; the
   corresponding widths (290 B full line → 152 B projected, envelope
   included — the same unit this scenario asserts) are row S1 of
   `evidence/spikes/a1-byte-width.md`. The record and its measurement
   live in the two spikes respectively; neither artifact carries both.
2. **Scenario**: projected top-level key set (explicit-literal
   oracle).
   **Expected**: exactly `revision`, `rule`, `gates`, `emit`,
   `next`, `writes`, `clear`, `escaped`, plus `escape_class` exactly
   when the same request's DEFAULT output carries it (the producer's
   own presence rule, C1 — not an unconditional "when escaped");
   omitted keys absent — never null or empty placeholders (the
   absent-not-null assertion rides here, Phase 2). The same normative
   fixture pins the exact bytes — the "Projected reference" line in
   `evidence/spikes/a2-encoder-mechanism.md` §Reference output.
3. **Scenario**: reflective partition-completeness over the resolve
   success payload, against the declared assignment source (C2).
   **Expected**: every field assigned to exactly one C2 group; the
   declared entry set and the struct field set are equal in both
   directions; the keys a projected run emits equal the declaration's
   `plan` side; an unassigned new field is a test failure. The
   assignment source is independent of the projection code — an
   oracle deriving the ECHO set from what the projection drops is
   tautological and does not satisfy this scenario.
4. **Scenario**: whole-tree flag-registration walk from the root, over
   a root with both auto-generated commands materialized.
   **Expected**: the set of commands registering `plan-only` is
   exactly `{flow resolve}`; vacuity-guarded on the four flow verbs.
   The walk is TOTAL — no name-skip, no `Hidden` gate, `help` and
   `completion` included (C1). Two preconditions, both asserted: the
   oracle materializes the auto-generated commands via cobra's own
   initializers before walking, and the walked set CONTAINS
   `completion` (and `help`). `completion` is the discriminating
   member — `help` is force-inited already
   (`internal/cli/help_all.go`, `InitDefaultHelpCmd`) so its presence
   proves nothing, whereas a bare `NewRootCmd()` tree has no
   `completion` at all (cobra creates it in `InitDefaultCompletionCmd`
   from `ExecuteC`), which would let the oracle pass vacuously on the
   exact command class C1 widened the walk to reach.
   `help_all.go::walkCommandTree` is not reusable here (it skips both).
   This scenario also carries the POSITIVE registration-shape
   assertion on the single registrant — `plan-only` is boolean,
   defaults false, no shorthand (C1, `D-naming`) — on the `--all`
   precedent
   (`flow_all_0011_test.go::TestReq36And50_NextRegistersABooleanAllFlagDefaultFalseWithNoShorthand`);
   the shape is otherwise stated twice in this RDR and asserted
   nowhere.
5. **Scenario**: `--as=text` ± the flag, repeated runs.
   **Expected**: projected text lines a byte-identical, stable
   subset of the default-mode lines — set membership, not a
   subsequence: the wire keeps struct declaration order while text is
   sorted by the shipped `flatten` (`respond/text.go`), and no clause
   requires the two orders to agree. Additionally the projected text
   is STRICTLY SHORTER in total rendered bytes (C1's width clause
   binds both modes); a subset assertion alone is satisfied by the
   equal set, so a zero-saving text projection would otherwise pass.
   The two measurands are JOINTLY SUFFICIENT only because the subset
   is asserted per-line and byte-identically: every projected line
   must appear in the default set unchanged, so the width reduction
   can only come from DROPPED lines, never from a shortened carried
   value. Stated because the pair is otherwise separable — a
   projection that kept every default line and shortened one carried
   value would satisfy a set-subset read loosely and still shrink
   the total, and it is the per-line byte identity, not the width
   clause, that refuses it.
   Witness: 15 default lines → 9 projected on the normative fixture.

### Desk trace

The MVV walked stepwise with every assertion in force at that step, on
the normative fixture (the pricing 2×2 call,
`evidence/spikes/a2-encoder-mechanism.md`). Witnesses are that call's
real bytes, not restatements.

| MVV step | Assertions in force | Witness | Verdict |
| --- | --- | --- | --- |
| 1 — no flag, json | C1 default byte-identity (against the pre-change golden, the battery's only pre-change side); `0005:A6` envelope list; `0010:C4` `emit` present; A9 nil-vs-empty | 14 keys, `"owned":{}` and `"readers":[]` present — non-nil containers, the A9 invariant — `"emit":{"dpa":"required","plan":"pro"}` | OK |
| 2 — `--plan-only`, json | C1 key set + absent-not-null + declaration order + carried-value byte identity; C2 always-keep core; S2 literal | `{"revision":"","rule":"paid-eu","gates":[],"emit":{…},"next":{},"writes":{},"clear":[],"escaped":false}` — 8 keys, equal to S2's list | OK |
| 3 — `--plan-only`, text | C1 text-subset; `0005:C1` both-modes agreement | default 15 lines → projected 9; the 6 dropped are exactly `model`, `observed.region`, `observed.tier`, `outcome`, `owned`, `readers` | OK |
| 4 — unrecognized outcome ± flag | C1 report-only; A6 | refusal fields `code`/`message`/`param`/`detail`/`hint`/`findings` — no echo member to project | OK |
| 5 — `flow next --plan-only`; tree walk | C1 non-registration + vacuity guard + walked-set precondition; A5 | `--all` probes are exact-name `Lookup("all")`; a resolve-local `plan-only` is invisible to them. Walked set on a bare root is `docs, flow, help, lint, version` — `completion` absent until `InitDefaultCompletionCmd`, so the oracle materializes it and asserts its presence first | OK (with the materialization precondition; without it the walk is vacuous on `completion`) |
| 6 — byte counts | A1; MVV step-6 pass bar | checked-in fixtures 290→152 (47.6%), 392→191 (51.3%), 412→238 (42.2%) — all clear the 40% bar; 708→146 B (79.4%) is the not-checked-in 48-fact demonstration, not a gate | OK |

Two orderings coexist and are not the same ordering: the **wire** keeps
struct declaration order (C1), while **text** is alphabetically sorted
by the shipped `flatten` (`respond/text.go`, `sort.Strings`). Step 3's
subset property holds under either, because projection only deletes
keys; no clause requires the two orders to match, and none should.

### Performance Expectations

- Wall time unchanged: projection is key deletion on the assembled
  payload before `respond.OK` — no extra I/O, no second marshal
  path; the 5ms call the Problem Statement measures is unaffected.
  This is a structural claim, not a measured budget: no oracle or
  scenario asserts a wall-time threshold, and none is specified,
  because the work removed (encoding five fewer keys) cannot make the
  call slower and a millisecond bar on a 5ms call would measure the
  harness rather than the change. If a future projection mode does
  work per field rather than deleting keys, that RDR owes the budget
  this one declines.
- Output width (A1's live table,
  `evidence/spikes/a1-byte-width.md`): 79.4% saved on the 48-fact
  motivating class (708→146 B), 47.6% on the 2×2 table, 42–51% on
  the gate/write state-machine shapes. The floor is the retained
  PLAN group; the saving grows with fact count.
- Determinism checklist (C1 claims byte-stability; A2 spike
  results): no hashing; one wire encoder
  (`internal/cli/clierr/clierr.go::WriteJSONLine` — `json.Encoder`,
  `SetEscapeHTML(false)`, compact one-line NDJSON); wire key order
  is struct declaration order, text mode sorts decoded keys
  (`respond/text.go::flatten`; 10 repeated runs, one hash); empty
  vs absent distinguishable (non-nil pointer to empty map renders
  `{}`, nil renders absent — never `null`); the version marker is
  `revision` (plan-side provenance, C2).

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

- JDR 0002 §D1 (`docs/jdr/0002-success-envelope-projection.md`) —
  the partition doctrine C2 instances; JDR 0001 §D8
  (`docs/jdr/0001-resolve-kernel-seam.md`) — why `owned`/`readers`
  entered the envelope.
- `docs/cli-output-contract.md` — the worked payload default mode
  keeps byte-identical.
- Peer-CLI prior art (langref checkouts):
  `gh-cli/pkg/cmd/pr/diff/diff.go` (`--name-only` fixed
  projection); `gh-cli/pkg/cmdutil/json_flags.go` (`--json
  <fields>` general form — Alternative 1).
- Literature (Resolve research,
  `evidence/research/resolve-partition-prior-art.md`): Kleppmann,
  *Designing Data-Intensive Applications*, Ch 11 pp480–481 (command
  results are emitted events; state is derived); Masse, *REST API
  Design Rulebook*, p92 (caller-side field trimming), p53 (ETag as
  opaque version identity); Richardson/Amundsen/Ruby, *RESTful Web
  APIs*, p358; Yao et al., *ReAct*, p2 (observations re-paid as
  trajectory context). Propose research:
  `evidence/research/propose-prior-art.md`.
- Spike evidence: `evidence/spikes/a1-byte-width.md`,
  `evidence/spikes/a2-encoder-mechanism.md`. Related: kata
  `intrastate#srz2`.
