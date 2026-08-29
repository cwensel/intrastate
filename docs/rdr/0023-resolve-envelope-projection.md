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
- **Profile**: foundational — provisional: one contract (a
  caller-controlled projection axis on the resolve envelope)
  spanning fences locked by RDR 0005 and RDR 0011.
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
  (`revision` stays — plan-side provenance, C2). Default-mode output is byte-identical to today's; every
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
of the call (wall time is 5ms). Measured on a 48-fact decision-table
call: full JSON is 1529 bytes, the `emit` + `rule` answer ~430 — about
75% of a resolve's output is the caller's own input echoed back under
`observed.*`. The echo is the **right default** and is settled: a live
consumer deliberately dropped its own status call because the resolve
already returns every fact it renders. The ask is only an opt-out for
chained or secondary calls — the second `--outcome` call in a chain
receives the same 48 facts again as pure duplication — and every
additional model a consumer migrates multiplies the per-run echo.

The fork is the mechanism and its blast radius against two locked
records. (a) A field-projection flag (`--select emit[,rule]`): but
`0011:D-naming` rejected the `--select` spelling by name, and
`0011:A14` makes non-registration of widening flags on
`resolve`/`read-state`/`set-state` a load-bearing structural negative
enforced by absence oracles in `internal/cli/flow_all_0011_test.go` —
though 0011's rejection ground was "the defect is the default," and
here the default is agreed correct, so whether that ground still binds
is 0011's to reopen, not a kata's. (b) A verbosity tier (`--as json`
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
Each answer also decides whether future verbs inherit the projection.

## Critical Assumptions

- **A1 The saving survives the SHIPPED partition: measured with the
  full plan group (not the seed's `emit`+`rule` subset), across at
  least two fixture shapes with different fact/gate ratios, the
  projected width is materially smaller on the chained-call class
  the RDR targets — the seed's 1529→~430 figure is the echo SHARE of
  the motivating call, not the shipped reduction.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: re-run the kata `intrastate#srz2` measurement over
    the shipped `--plan-only` (or a prototype projection) on the
    motivating decision-table shape AND a gate/write-heavy fixture;
    record the default vs projected byte table.
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
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: a minimal encoder spike over
    `internal/cli/flow_resolve.go::resolvePayload` showing default
    bytes unchanged (fixture-diff) and projected bytes = default
    minus echo keys, plus a repeated text-mode run.
  - **If wrong**: C1's absence-not-null and byte-identity clauses are
    unimplementable as specced and the mechanism (or the clause) must
    change before lock.
- **A3 No shipped oracle or fixture asserts an echo field's presence
  in a way an OPT-IN flag moves: every existing resolve-payload
  assertion runs default-mode, so the 0005/0010/0011 suites do not
  move and "default unchanged" holds without editing a predecessor's
  tests.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: sweep `internal/cli/*_test.go` for resolve success
    payload assertions; confirm none invokes the new flag's path and
    none needs edits.
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
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: `0024:C4` (envelope clause: `dispositions` appended
    after `emit`, "never omitted", joined by the "same single
    `Plan.RuleID` join path as `emit`") read against C2's partition
    assignment below; re-check at Resolve if 0024's draft moves.
  - **If wrong**: the joint decision homed at `cli/0023:C2` reopens
    and one of the two envelopes must move.
- **A5 The shipped absence-oracle idiom extends to a second flag
  without colliding with the `--all` oracles: a `Lookup("plan-only")`
  sweep over `next`/`read-state`/`set-state`, the `flow` group, and
  the root's persistent set is independent of
  `internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup`.**
  - **Status**: Pending
  - **Method**: Source Read
  - **Evidence**: the cited test's structure (per-name `Lookup` with
    the four-verb vacuity guard) — confirm nothing in it or in
    `registerSelectionFlags` is name-generic.
  - **If wrong**: C1's structural negative cannot be pinned in the
    house idiom and needs its own form.
- **A6 A refusal path carries nothing to project: `respond.Fail`
  renders `clierr.CLIError` (+ `findings`) with no request-echo
  member, so `--plan-only` cannot change a single refusal byte and
  "report-only" holds on failures by construction.**
  - **Status**: Pending
  - **Method**: Source Read
  - **Evidence**: `internal/cli/respond` failure emission path and
    `internal/cli/clierr` envelope fields.
  - **If wrong**: the flag changes refusal output, breaking C1's
    report-only clause and the caller's error handling.
- **A7 The echo group really is an echo: `observed` renders the
  caller's `--tag` set verbatim (`0005:C1` refuses rather than
  coerces — a non-canonical set literal is `flow-tag-invalid`, an
  owned/reserved key is refused before any accessor), and `model`
  and `outcome` are the caller's own flag values copied through
  (`runFlowResolve`: `payload.Outcome` is `--outcome`'s string) — no
  CLI-side defaulting, coercion, or expansion produces an observed
  entry the caller did not supply. One value-preserving exception is
  in scope, not a refutation: a set-kind literal is re-encoded to
  the canonical bytes of the same members (the form a conforming
  caller already speaks).**
  - **Status**: Pending
  - **Method**: Source Read
  - **Evidence**: `internal/cli/flow_input.go` tag parsing (refusal
    codes, no rewrite path) and
    `internal/cli/flow_resolve.go::runFlowResolve` payload
    population; a fixture asserting `observed` byte-equals the
    supplied tag set.
  - **If wrong**: `observed` is derived output, its projection hides
    what the kernel actually matched against, and the field moves to
    the plan group before lock.

## Proposed Solution

### Approach

Give `flow resolve` — and only `flow resolve` — a caller-controlled
projection axis, as one boolean flag: `--plan-only`, default false.
Default-mode output is byte-identical to today's; the echo default is
settled and stays. Under the flag the success payload carries exactly
the **plan group** the verb's own extended help already names under
"Reading a successful plan" (`internal/cli/flow_resolve.go::flowResolveExtendedDesc`:
`rule`, `next{}`, `writes{}`, `clear[]`, `emit{}`, `escaped`,
`escape_class`, `gates[]`), plus `revision` as the plan's provenance
(the one correction the premortem forced — it is loader-produced,
not request echo), and omits the **echo group** (`model`,
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
byte-identical values on every carried key; and an identical
invoked-reader set. A later change that skips work whose only
consumer is a projected-away field breaches this clause.
Presence-rule fields keep their own contracts inside the plan group:
emit stays present as {} (0010:C4), escape_class stays
omitempty-on-unescaped.

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
— vacuity-guarded by requiring the four flow verbs to exist. The
behavioural command-error run is corroboration, never the assertion.
This contract mints no new refusal code and no new exit group.
```

**C2**

```normative
The projection axis is a PARTITION this RDR owns, not a
caller-supplied field list. Every current and future flow resolve
success-payload field is assigned to exactly one side when it is
added: the ECHO group — what the caller supplied verbatim or can
reconstruct from its own prior calls: the model reference, the
observed tags (--tag, echoed unchanged), the assembled owned view
and the invoked reader identities, the requested outcome (--outcome,
echoed unchanged) — or the PLAN group — what the run decided, its
provenance, and what a caller acts on or must not misread: rule
identity, gate results, authored answers and their interpretations,
planned next/writes/clear, the escape disposition, and revision.
revision rides the PLAN side deliberately: it is produced by the
loader, not restated from the request, and it is the plan's only
provenance — a chained caller cannot otherwise detect that the model
changed under the same path between calls.

The partition is ENFORCED, not prose: a reflective oracle MUST
assert every field of the resolve success payload is assigned to
exactly one group, so an unassigned new field is a test failure, not
a silent default — the projection MUST NOT be implemented as a bare
omit-list whose complement is "whatever else exists".

An ALWAYS-KEEP core — rule, escaped, escape_class, revision — MUST
survive every current and future projection mode of this verb: if
the boolean axis ever generalizes to an enum or a field list, no
mode may omit them, so a projected payload can never launder a
rescued plan into an ordinary one or detach a plan from the model
revision that produced it. When a new field fits neither group's
description cleanly, it joins the PLAN group — over-reporting is
recoverable, silent omission is not. Under this rule 0024's
dispositions, if that RDR lands, is a PLAN-group field: each token
is assigned by the [emit] declaration to the selected row's authored
emit value (0024:C4's join) — derived from plan-group inputs only,
never from observed or owned — and its never-omitted clause is
untouched by this projection.

A LATER verb that adopts a caller-controlled success-payload
projection MUST follow the same two rules: the projection is applied
to the verb-specific result before respond.OK — never in the respond
gateway, never differing by output mode — and the projectable set is
that verb's echo group, never its decision. This clause is the
conformance surface for a future report/export verb; it obliges no
verb to offer a projection, and it reserves nothing about the flag
spelling a later verb picks.
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
- **Selection / predicate** — the partition rule (C2) decides which
  side a field joins: supplied verbatim by the caller or
  reconstructable from its own prior calls → echo; a decision, its
  provenance, an authored answer, or a planned mutation → plan;
  unclear → plan. The rule is applied once, at the moment a field is
  added, by the RDR adding it — and the reflective
  partition-completeness oracle makes skipping that assignment a
  test failure rather than a silent omit-list default.

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
| Absence-oracle idiom | `internal/cli/flow_all_0011_test.go::TestReq46And47And65And98_AllIsAbsentFromTheOtherThreeVerbsAndTheFlowGroup` | Asserts one flag name per sweep | Reuse as pattern | New oracle for `plan-only` over the other three verbs, group, root |
| Verb help | `flow_resolve.go::flowResolveExtendedDesc` ("Reading a successful plan") | None | Extend | The help's plan list is the partition's user-facing statement; gains the flag line |

### Decision Rationale

The routable question — should the resolve payload have a
caller-controlled projection axis, and where does it live — was scored
as a matrix (profile: foundational), approaches × deciding criteria:

| Criterion | O1 fixed plan-projection flag | O2 named-field selection | O3 gateway verbosity tier | O4 no CLI axis (caller-side jq) |
| --- | --- | --- | --- | --- |
| Correctness fit (drop the chained-call echo; refusals untouched) | full — the echo group is exactly the measured duplication | full, plus generality nothing asked for | partial — omission differs by mode | partial — needs a per-skill recipe carried in every prompt |
| Prior-art alignment | `gh pr diff --name-only` (fixed projection) | `gh --json <fields>` (field selection + universe policing) | none surveyed — gh renders text from the same selected fields, not a tier | jq is gh's *escape hatch* (`--jq`), offered beside, not instead of, projection |
| Safety of the projected payload | `escaped`/`escape_class` never omissible, by construction | caller can drop `escaped` unless an always-keep set is added — which IS this partition | n/a | a recipe can drop anything, including refusal fields on the shared stdout |
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
the recommendation survived with mitigations folded, none forcing a
switch. Folded: `revision` moved to the plan group as provenance
(P-3's produced-not-echoed half; its `outcome` half is refuted by
source — `runFlowResolve` copies the `--outcome` flag verbatim, so
`outcome` stays echo, with A7 pinning that read); the reflective
partition-completeness and explicit projected-key-set oracles
(P-9/P-12, the seed-1 class); the whole-tree exactly-one
registration oracle replacing the enumerated-sibling sweep (P-13);
the differential report-only oracle including the invoked-reader set
(P-5); A1 rescoped to the shipped partition across two fixture
shapes (P-7/P-8); A2 rewritten around the encoder trilemma and the
flattener's key-sort (P-1, with the map-order scare answered by
`flatten`'s `sort.Strings`); A4 rewritten to quote 0024:C4's join
inputs instead of summarizing them (P-15, the seed-3 class); the
inverted-default rejection restated on its own ground (P-16); the
always-keep core (P-17); the refusal-width scope and version-skew
loudness recorded in Trade-offs/Failure Modes (P-6, P-11); P-14's
generic-bucket message limit accepted as `0011:A14`'s own.

Joint-decision check — open peers at depth 1 with Status Draft/Final:
0012–0022 and 0024 (13, all Draft; 0001–0011 are Implemented). FIRE
on 0024: shared modify-anchor
`internal/cli/flow_resolve.go::resolvePayload` (0024:C4 appends
`dispositions`; this RDR projects the echo group off the same
payload). Homed, not paused: `cli/0023:C2` is the normative home —
it assigns `dispositions` to the never-projected plan group, and
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
consumer of C2's envelope-general clause if its export surface ever
wants a projection; C2 obliges it to nothing. Absence arm: n/a —
this proposal converts no refusal into an acceptance (`--plan-only`
on the other verbs stays refused; success reporting only narrows).
Bridge sub-check: n/a — no Cluster membership; no plan here
introduces a surface a sibling schedules for deletion (no
bridge/retire/replace cues).

Premortem: hardened (hardened)

Ground-sweep: clean (24 anchors)

Joint-check: fired → 0024 (home: cli/0023:C2)

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
pays. Scope review read the primary shape as already adjudicated;
the Stage-2 source read narrows that: RDR 0011 rejected `--select` as
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
  `revision` on the plan side, the premortem's provenance
  correction. The
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
  echo as transcript tokens. The seed's ~72% is the echo *share* of
  the motivating call measured against `emit`+`rule` alone; the
  shipped reduction is smaller (the plan group also carries `gates`,
  `next`, `writes`, `clear`, `escaped`, `revision`) and A1 measures
  it honestly across fixture shapes before lock. The saving
  multiplies with every model a consumer migrates.
- Scope of the guarantee: the token reduction covers successes only.
  Refusal envelopes are untouched — they carry diagnosis
  (`findings`), not echo (A6) — so a probing caller that mostly
  refuses saves nothing, by design.
- Positive: the default is untouched — no consumer changes, no 0005/
  0010/0011 test moves (A3), no respond-gateway change.
- Negative: `flow resolve` now has two success-payload widths; docs,
  help, and any payload-shape oracle must say which width they speak
  of.
- Negative: the partition (C2) is a standing obligation — every
  future payload field must be assigned a side by the RDR adding it,
  and a mis-assignment is a spec defect, not a style choice.
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
  normative list.
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
  loses provenance (which model/revision produced this plan).
  **Mitigation**: opt-in with an unchanged default; the help line
  states the omitted group by name; the caller that opts in is the
  caller that already holds the request.

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

- [ ] All Critical Assumptions verified
- [ ] Ordering tolerance with cli/0024 confirmed (A4): either RDR may
      land first; the second lands with `dispositions` already/newly
      in the plan group and no contract in either moves.

### Minimum Viable Validation

Over the checked-in decision-table fixture
(`models/examples/pricing-decision-table.toml` or the 0005 fixture
family):

1. Run `flow resolve` with a recognized outcome and discriminating
   `--tag`s, `--as=json`, without the flag → today's full payload,
   byte-identical to the pre-change fixture.
2. Re-run the same invocation with `--plan-only` → exactly the
   normative projected key set (`revision`, `rule`, `gates`, `emit`,
   `next`, `writes`, `clear`, `escaped` — `escape_class` when
   escaped); every carried field byte-identical to step 1's;
   `model`/`observed`/`owned`/`readers`/`outcome` absent (not null,
   not empty).
3. Re-run step 2 with `--as=text` → the projected lines are a
   byte-identical subset of step 1's text lines, stable across
   repeated runs.
4. Re-run with an unrecognized outcome, ± the flag → byte-identical
   refusal envelopes and exit codes.
5. Run `flow next --plan-only` → `command-error`, exit 2; and the
   whole-tree structural oracle (the set of commands registering the
   flag is exactly `{flow resolve}`) passes.
6. Record default vs projected byte counts on the motivating-model
   shape and one gate/write-heavy fixture (A1's table).

End-state: one flag, two widths, one decision; refusals and default
mode untouched.

### Phase 1: Projection on the verb

Register `--plan-only` on `resolve` only and hand `respond.OK` the
projected verb result when set — the mechanism A2 verified, applied
after payload assembly, before the gateway.

### Phase 2: Oracles

The whole-tree registration oracle (commands registering the flag ==
exactly `{flow resolve}`, vacuity-guarded); the differential
report-only oracle (± flag: exit codes, refusal bytes, strict
key-subset with byte-identical carried values, identical
invoked-reader set); the reflective partition-completeness oracle
(every payload field in exactly one C2 group); the explicit
projected-key-set oracle; the text-subset oracle; the
absent-not-null oracle.

### Phase 3: Docs and help

`docs/cli-output-contract.md` gains the projected worked payload
beside the full one; `flowResolveExtendedDesc` gains the flag under
"Reading a successful plan"; the partition statement (C2) lands where
the payload fields are documented.

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
