# Recommendation 0029: What a released version number promises an agent about the machine-readable output

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-09-11
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
- **Type**: Architecture
- **Profile**: large — the version promise over every machine-readable
  vocabulary this CLI emits: the `--as=json` envelope's `schema_version`,
  the three-tier stability vocabulary, and the tier assigned to each of
  the fourteen tiered surfaces; user-facing yes; locks format.
- **Priority**: High
- **Related Issues**: none
- **Seam Lineage**: no prior accretion

## Problem Statement

> Synthesized from a free-form idea rather than a kata — flagged for author
> review.

An agent author pins a released `intrastate` version in a skill, a CI job, or
a harness, and writes code that parses `--as=json`: it branches on the
envelope's `type`, matches a refusal `code`, checks an exit status, and —
if it lints models — expects a model that passes today to pass tomorrow. The
agent then upgrades, because an install step resolved a newer release or a
package manager pulled one. It needs to know, before upgrading, what that
version number said was allowed to change. It discovers the gap the hard way:
its parse silently takes the wrong branch, or its CI goes red on a model
nobody edited.

The audience makes this sharper than it would be for a human-facing CLI. A
human reader degrades gracefully when a field is renamed or a new diagnostic
appears — they read the message and adapt. An agent hard-fails, or worse,
mis-parses without failing. So the promise a version number carries has to be
stated over the *machine-readable* surfaces specifically, and stated somewhere
an agent can find it.

Internally, `intrastate` already has three such surfaces, and they already move
independently in code:

- the `--as=json` envelope (`type`, `data`, `notes`, `warnings`), the refusal
  `code` vocabulary, and the exit-code classes;
- the model TOML's `[model] version` integer, pinned at `1`, with a load
  refusing anything else;
- the lint finding taxonomy, whose per-finding code identity RDR 0017 fixed.

No record says what a release version promises about any of them, nor which
mechanism proves the promise held. The open fork this RDR must settle: a
**new** lint finding code removes nothing and renames nothing, yet a model
that lints clean at one version starts refusing at the next — breaking a
consumer's pipeline. Semantic versioning has no opinion on whether that is a
minor or a major; this project must have one, and every future change to the
envelope or the taxonomy will need to read it.

Out of scope: the release *mechanism* — tag format, goreleaser, and semver
arithmetic. Those are settled by existing tooling and wider convention, and
belong in ordinary documentation rather than a decision record.

In scope, and load-bearing: the fact that the series starts at **0.1.0, not
1.0.0**. This is not release mechanism despite sitting next to it — under
SemVer a `0.x` major carries no compatibility guarantee, so the starting
number decides what the promise can say and when it starts binding. This
RDR therefore settles two things at once: the promise that takes effect at
1.0.0, and what the CLI does during `0.x` so that promise costs nothing to
adopt when it arrives.

## Critical Assumptions

- **A1 The graph-lint advisory tier can accept new members, or `0006:C17`
  can be amended to allow it.**
  - **Status**: Verified — ruled 2026-09-11: `0006:C17` is opened as
    append-only. The code already permits it; the closure was the only
    blocker and it is lifted.
  - **Method**: Peer RDR
  - **Evidence**: `0006:C17` states "The advisory tier is closed at
    `graph-coverage-closed-by-escape`, `graph-redundant-row`,
    `graph-unreachable-rule`, and `graph-vacuous-atom`" — a closure at
    exactly four members; 0006 is Implemented. The IMPLEMENTATION,
    however, imposes no such limit: severity is DERIVED, not carried
    per-code — `internal/graphlint/taxonomy.go::severityFor` returns
    `SeverityInfo` for anything `IsBlocking` rejects, and `IsBlocking` is
    `slices.Contains(blockingCodes, code)`. So a fifth advisory code
    becomes `info` with no severity plumbing, no exhaustiveness switch,
    and no cardinality constant anywhere in `internal/graphlint/` or
    `internal/cli/`; the eleven `AdvisoryCodes()` consumers all iterate.
    Exactly one assertion blocks it:
    `internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
    compares `AdvisoryCodes()` to a four-element literal with
    `slices.Equal` — a by-value set assertion, so the fix is editing the
    literal, not relaxing a number.
  - **Ruling (2026-09-11)**: C17's closure is incidental-as-then-stood, not
    load-bearing, and is opened. This follows `JDR 0001 §D10` rule 3, which
    reopened a closed `CLIError` field on a peer's request and recorded the
    result as "a 0006 citation repair, not a reopening" — the same shape.
    C17's load-bearing guarantee is its last sentence, "Advisory findings
    MUST NOT change the success disposition," which adding a member does
    not touch; C3 restates that guarantee rather than amending it.
  - **Consequent obligations** (both land in the implementation of this
    RDR, not in a separate change): 0006 restates C17 as append-only — the
    `growing` tier C2 defines — and the four-element `want` literal at
    `internal/graphlint/findings_0006_test.go:288-293` is edited in the
    same change that adds any fifth advisory code. Until 0006 is amended,
    this RDR and 0006 disagree on the tier's openness; the amendment is a
    prerequisite of implementation, and 7.1 cluster reconcile over
    {0006, 0029} is where the two records are re-read against each other.
  - **If wrong**: C3 is unimplementable as written for graph-lint codes —
    a new finding could only enter at `blocking`, which is exactly the
    verdict-changing event the policy exists to prevent. Surfaces as a
    contradiction between this RDR and 0006 at Stage 7.1, or as an
    implementer unable to add an `info` code without amending a Final
    peer.
- **A2 `schema_version` ships before 1.0.0, while the `0.x` series still
  permits incompatible change.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: No tag has ever been cut — `git tag --list` is empty and
    `git describe --tags` fails. `internal/version::resolve` defaults to
    `Info{Version: "dev", Commit: "none", Date: "unknown"}`
    (`internal/version/version.go::resolve`), filling commit and date from
    VCS stamps but never moving `version` off `"dev"` without an
    `-ldflags` injection or a module version;
    `internal/cli/version.go`'s verb reports that value and nothing else.
    `.github/workflows/release.yml` releases on a `v*` tag push
    (`on.push.tags: ["v*"]`, alongside a `workflow_dispatch`). No such tag
    exists, so the first release is 0.1.0 and the window is the whole
    `0.x` series rather than one release.
  - **If wrong**: adding `schema_version` after 1.0.0 makes the
    compatibility field itself a breaking change under C1's own rule,
    forcing a 2.0.0 for the field that exists to prevent major bumps.
    Surfaces as a strict-parsing consumer failing on the release that
    announces the compatibility promise. Note this is a *soft* deadline
    for the whole `0.x` series, not a hard one at 0.1.0 — the risk is
    procrastination past 1.0.0, not missing the first tag.
- **A3 No current consumer parses the envelope strictly enough that one
  added field breaks it.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `DisallowUnknownFields` IS used in production, but only
    on the INPUT path — `internal/table/source.go::decodeStrict`, which
    rejects unknown keys in the model TOML under `0002:C3` ("an unknown
    schema field is a stable refusal, never a silent no-op"); its one
    in-repo call site is TOML, at `internal/table/load.go::Load`. There is a real in-repo consumer that
    re-parses the OUTPUT envelope, and it is TOLERANT: `--plan` decodes a
    `flow resolve --as json` envelope through
    `internal/cli/flow_input.go::planEnvelope`, a four-field struct read by
    a plain `json.Unmarshal`, which ignores unrecognized top-level
    properties by Go's default decoder behaviour. Test-site envelope
    decodes are likewise plain structs; no test anywhere sets
    `DisallowUnknownFields`. No `.claude/` skill, hook, Makefile, or CI
    workflow parses the envelope at all, and
    `schema_version`/`format_version` appear nowhere outside this record's
    prose. So the additive field is safe for every DECODER that exists.

    The binding constraint is not a decoder. `internal/cli/flow_mvv_0023_test.go`
    asserts BYTE-IDENTITY over a whole terminal envelope, comparing a
    `flow resolve --as=json` line by string equality against the checked-in
    golden `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`
    (`{"type":"ok","data":{…}}`, no `schema_version`). A non-`omitempty`
    field fails it on the commit that adds it. That is a pinned baseline
    rather than a consumer — 0023 captured it to prove ITS change was
    additive — so the correct response is to re-capture it, not to weaken
    C1. It is an implementation obligation (Phase 1 Step 1), not a
    refutation.
  - **If wrong**: the additive change is breaking for that consumer
    regardless of what C1 declares, and the rollout needs a deprecation
    window rather than a single release. Surfaces as a harness failing to
    parse immediately after upgrade.
- **A4 The three tier names partition every machine-readable surface this
  CLI emits — no surface needs a fourth tier.**
  - **Status**: Verified as to the PARTITION; C4's assignment list is
    incomplete and is corrected below.
  - **Method**: Source Search
  - **Evidence**: An enumeration of every string vocabulary reachable on
    the `--as=json` wire found **sixteen** emitted sets: fourteen that
    take a tier, plus the two model-authored namespaces that deliberately
    take none. C4 originally named seven of the fourteen. Every one of the
    seven it had not named fits an existing tier, so no surface needs
    semantics the three tiers do not provide — the partition holds, which
    is what this assumption claims. Five sit at the envelope's top level
    or under `data`, with the tier each takes (now assigned in C4): the
    stderr advisory
    `level` set `note`/`warning` (`internal/cli/respond/respond.go::Note`,
    `::Warn`) → `frozen`, a `type`-shaped discriminator on the second
    stream; the gate `verdict` set `allow`/`deny`/`indeterminate`
    (`internal/accessor/model.go::Verdicts`, emitted at
    `internal/cli/flow_exec.go::gateResult`) → `frozen`, since it is also
    the accessor-protocol INPUT alphabet and `cmdbind` rejects a stranger,
    so it cannot append without breaking gate authors; `data.escape_class`
    (`internal/resolve/resolve.go::RefusalKinds`, emitted at
    `internal/cli/flow_resolve.go`) → `frozen`, already enforced by
    `internal/resolve/guard_atoms_test.go::TestReq79_NoSixthRefusalKindIsMinted`;
    the `flow next` unknown-`reason` set
    (`internal/cli/flow_next.go::unknownFact.Reason`) → `append-only`,
    a set DISTINCT from the `graph-unprovable-coverage` `reason` set C4
    already names, on a different field with a different producer; and the
    `graph-lint-failed` aggregate code
    (`internal/graphlint/taxonomy.go::AggregateCode`) → `frozen`, the one
    code that envelope ever carries. The remaining two ride the fields of
    a `findings[]` element, which the first pass of this census missed by
    enumerating only the top level and `data`: `findings[].operator`
    (`internal/guard/grammar.go::Operators`) → `frozen` on `0003:C7`'s
    "the initial operator vocabulary MUST be closed and typed", enforced
    by an exhaustive `slices.Equal` test; and `findings[].block`
    (`internal/resolve/guard.go::Block`) → `append-only`, NOT frozen:
    `JDR 0001 §D12` already grew it from two members to three, ruling that
    "0007's 'exactly two' is the fence that gives", and `isGuardBlock`
    tolerates further tokens by construction. The two took different
    tiers for the same reason the tier vocabulary exists — one set is
    closed at its boundary, the other has grown under record. Both are
    emitted at `internal/graphlint/analysis.go`,
    `::coverage.go`, `::groups.go` and `internal/cli/flow_exec.go`. Two
    model-authored open namespaces — `findings[].class` and
    `data.dispositions` — are correctly tier-less: `clierr` ascribes them
    no meaning, so they are author surface, not CLI vocabulary; the same
    holds for `findings[]`'s free-text `key`, `literal` and `dimension`,
    which are drawn from the model under lint and are not vocabularies.
    C4 also said "the graph-lint finding codes" where the blocking and
    advisory tiers take DIFFERENT tiers; it now names each.
  - **If wrong**: C4 is incomplete and an unassigned surface has no stated
    promise — the exact gap this RDR exists to close, reappearing inside
    the fix. Surfaces at implementation as a vocabulary that fits no tier.
    This is what the enumeration above actually found, at seven surfaces —
    five on the first pass and two more the grounding lens caught inside
    `findings[]`; the correction is C4's, not the tier vocabulary's. That
    the gap recurred once inside the fix is why C4 now states where a
    vocabulary must be looked for, not only which ones were found.

## Proposed Solution

### Approach

Draw the compatibility edge at **default-enablement, not at code
existence**, and state the promise **per surface** rather than as one
blanket claim over "the JSON output".

**The series is `0.x`, and that is the governing fact.** The first release
is 0.1.0, not 1.0.0, and under SemVer a `0.x` series carries no
compatibility guarantee: anything may move in any release. So this RDR is
not writing a promise that binds today. It is deciding **what the promise
will be when 1.0.0 arrives, and what the CLI does in the meantime so that
promise is cheap to keep** — the surfaces get their tiers and their
`schema_version` now, while breaking them is still free, precisely so the
1.0.0 commitment is a formality rather than a redesign.

Stated at 1.0.0 and after, a released version promises an agent this: *an
invocation that succeeded at version N, over input you did not change,
still succeeds at N+1 within the same major.* It does **not** promise the
vocabularies are frozen. New refusal codes, new load categories, and new
lint findings may appear in a minor release — because they already have,
four times — but a new finding may not change a verdict that was
previously clean until a release that discloses the promotion.

Through `0.x` the tiers are **declared and honored as intent, not
guaranteed**: they tell a consumer which surfaces are expected to be
stable at 1.0.0, and they tell this project which changes are cheap now
and expensive later. A consumer pinning a `0.x` version is told plainly
that it is pinning, not relying.

Three mechanisms carry that promise, all of which reuse surfaces that ship
today:

1. **A `schema_version` on the `--as=json` envelope**, versioned
   independently of the binary's release version, following the
   Terraform/OpenTofu rule: minor for additive, major for
   non-backward-compatible, consumers instructed to ignore unrecognized
   keys. The binary version answers "what build is this"; the schema
   version answers "can my parser read this", and those two questions
   have never had the same answer.
2. **A declared stability tier per vocabulary**, replacing the four
   divergent code-comment policies with one vocabulary used the same way
   everywhere. The word "closed" is retired from this role, because it is
   already live in two incompatible senses on adjacent surfaces (Key
   Discoveries).
3. **A severity-promotion rule** governing `info` → `blocking`. This is
   the only genuinely new policy. Introduction of a new code at `info` is
   a minor-release event needing no disclosure beyond release notes;
   *promotion* to `blocking` is the disclosed event, because that is the
   one that turns a consumer's passing pipeline red.

The insight the prior art forces: intrastate does not need to choose
between "new codes are breaking" and "new codes are free". Five peer
tools independently split this at the same place — the boundary is
whether the new diagnostic *fires by default*, not whether it exists.

### Technical Design

The three vocabularies get one declared tier each, and the tier is a
property stated in the RDR and asserted by a test, not a comment.

**Stability tiers** (the replacement vocabulary for "closed"). Three
tiers, defined normatively in C2 and assigned to surfaces in C4:
`frozen` (a consumer may exhaustively switch), `append-only` (it must
tolerate an unknown member), and `growing` (a new member may also newly
fire on input that previously passed, subject to the promotion rule
below).

Three is the count because the tiers are ordered by **what a consumer may
assume**, and each step down removes exactly one assumption: `frozen` →
`append-only` drops exhaustiveness, `append-only` → `growing` drops
verdict-stability. A surface needing a fourth tier would have to remove
some other assumption, which is what A4 tests. Each surface is assigned
exactly one tier; that assignment is the contract, and the tier names are
just the shorthand.

**The promotion rule.** A finding code enters at `info`. `0006:C17`
already guarantees `info` cannot change the outcome — "Advisory findings
MUST NOT change the success disposition" — so introduction is verdict-safe
by construction, with no new machinery. Promotion from `info` to
`blocking` is a separate, later, disclosed release event. This is the
Ruff preview→stable ladder built out of parts intrastate already has: the
severity field *is* the preview gate.

**Where the promise lives.** `docs/cli-output-contract.md` — the document
`.rdr/resources.md` already names authoritative for verb I/O, the error
envelope, and exit-code mapping, and which today contains zero occurrences
of "release", "upgrade", or "breaking". The promise goes beside the wire
format it constrains, matching the in-repo precedent at
`.goreleaser.yaml:55`, which names the archive template "a compatibility
surface, not cosmetics" in the file that produces it. `llms.txt` gets a
pointer, since it is what tells agents to bind to the binary in the first
place.

#### Normative Contracts

**C1**

```normative
The `--as=json` terminal envelope carries a `schema_version` string field
of the form `MAJOR.MINOR`, present on both terminal records: the `ok`
envelope (`internal/cli/respond::Success`) and the refusal record
(`internal/cli/clierr::CLIError`). Those two records are structurally
asymmetric and this contract does not change that: `ok` carries a `type`
discriminator, a refusal carries no `type` and no wrapper — it is the bare
`CLIError`, per `0005:C1`. `schema_version` is the one field they share.

It is versioned independently of the binary's release version and MUST NOT
be derived from it.

`schema_version` is NOT `omitempty` on either record, which is what makes
it a version a consumer can rely on rather than one it must handle the
absence of.

This does not breach `0005:C1`'s "exactly one omitempty structured field,
`findings`". That restriction descends from `JDR 0001 §D10` rule 3, whose
subject is a STRUCTURED carrier — a record field on which inner
discriminators ride so that codes do not proliferate (rule 2: "One CLI code
per caller-branchable failure; inner discriminators ride one structured
field"). `schema_version` is a scalar version marker carrying no
discriminator and no payload, so it does not consume that budget. The
refusal envelope still has exactly one structured field, `findings`.

The schema version begins at `"0.1"` and tracks the wire, not the binary.
While its major is `0` the schema is explicitly unstable and MAY change
incompatibly in any release; the minor component still increments on every
change so a consumer can detect movement even while it cannot rely on
compatibility.

From `"1.0"` onward: the minor component increments for backward-compatible
additions — a new optional field, a new member of an append-only or growing
vocabulary. The major component increments for changes that are not
backward-compatible: a removed or renamed field, a removed or renamed member
of any vocabulary, or a change to a Frozen surface.

A consumer MUST ignore object properties with unrecognized names, and MUST
reject an envelope reporting an unsupported major.

This tolerance rule is deliberately the INVERSE of the rule this CLI
applies to its own input: `0002:C3` makes an unknown key in a model TOML a
stable refusal. The asymmetry is intended and follows the audience —
an authored input document is refused so its author learns of the typo,
whereas an emitted envelope is tolerated so a consumer survives the
addition. A surface MUST NOT be given the input rule and the output rule
at once.
```

**C2**

```normative
Every machine-readable vocabulary this CLI emits carries exactly one
declared stability tier, and the tier is recorded in
`docs/cli-output-contract.md` beside that vocabulary:

- `frozen` — no member added or removed within a major. A consumer MAY
  treat an unrecognized member as a defect.
- `append-only` — members MAY be added in a minor; none removed or renamed
  within a major. A consumer MUST tolerate an unrecognized member, and
  MUST NOT assert on the set's cardinality, a member's ordinal position,
  or a tail position.
- `growing` — as `append-only`, and additionally a new member MAY fire on
  input that previously produced no such finding, subject to C3.

The term `closed` MUST NOT be used to describe any of these tiers, in code
comments or documentation, because it is currently live in two
incompatible senses.

While the binary's version is `0.x`, a tier is a DECLARATION OF INTENT: it
states what the surface is expected to promise at 1.0.0, and MUST NOT be
read by a consumer as a guarantee already in force. The tier names and
assignments are nonetheless authored and maintained from the first release,
so that reaching 1.0.0 requires no reclassification.
```

**C3**

```normative
A lint finding code is introduced at severity `info`. Introduction of an
`info` code is a minor-release change and MUST NOT alter the success
disposition of any input (this restates no new rule; it is `0006:C17`'s
existing guarantee, cited not amended).

Promotion of a finding code from `info` to `blocking` is a distinct
release event. It MUST be disclosed in the release notes for the release
that carries it, naming the code. A promotion MUST NOT occur in a patch
release.

The disclosure obligation binds during `0.x` as well as after. It is the
one clause here that does not wait for 1.0.0: a consumer's pipeline going
red is equally disruptive at 0.4.0, and disclosure costs a release note
rather than a design constraint.

A code MAY be introduced directly at `blocking` only when it reports a
condition that was already refused by some other code — that is, when the
promotion re-attributes an existing refusal rather than creating a new one.
```

**C4**

```normative
The tier assignments at this RDR's implementation:

- `frozen`: the envelope `type` discriminator; the severity vocabulary
  (`blocking`, `info`); the exit-code classes emitted by
  `internal/cli/clierr::ExitCodeFor`; the stderr advisory `level` set
  (`note`, `warning`); the gate `verdict` set (`allow`, `deny`,
  `indeterminate`), which is frozen because it is simultaneously the
  accessor-protocol input alphabet; `data.escape_class`
  (`internal/resolve::RefusalKinds`); the `graph-lint-failed` aggregate
  code; `findings[].operator` (`internal/guard::Operators`), frozen on
  `0003:C7`'s authority — "The initial operator vocabulary MUST be closed
  and typed… Unknown operators MUST be rejected" — enforced by the
  exhaustive `slices.Equal` assertion in
  `internal/guard/guard_grammar_0003_test.go::TestReq2_OperatorVocabularyIsClosedAndTyped`.
- `append-only`: `internal/table::Categories()`; the CLIError `code`
  vocabulary; the `graph-unprovable-coverage` `reason` set; the
  `flow next` unknown-`reason` set (a distinct vocabulary on a distinct
  field from the preceding entry); the graph-lint BLOCKING finding codes;
  `findings[].block` (`internal/resolve::Block`). That last is
  append-only rather than frozen on the authority of `JDR 0001 §D12`,
  which already grew the set from two members to three — "`resolve.Block`
  gains `BlockMatch` in 0007 Phase 1; 0007's 'exactly two' is the fence
  that gives." `internal/resolve/guard.go::isGuardBlock` is written to
  tolerate further tokens, treating every block outside `all`/`unless` as
  a non-operand of `0007:C6`'s verdict formula.
- `growing`: the graph-lint ADVISORY finding codes. `0006:C17` is amended
  from "closed at" its four members to append-only as part of this RDR's
  implementation (A1) — the tier assignment here is what that amendment
  records.

Two emitted namespaces take NO tier, deliberately: `findings[].class` and
`data.dispositions` carry model-authored tokens to which `clierr` ascribes
no meaning. They are author surface, not CLI vocabulary, and a consumer
MUST NOT read them as either. The remaining author-supplied `Finding`
fields — `key`, `literal`, `dimension` — are free text drawn from the
model under lint, not vocabularies, and take no tier for that reason.

A vocabulary is tiered wherever it is EMITTED, including on the fields of
a `findings[]` element — not only at the envelope's top level and under
`data`. A surface added later takes a tier assignment in the same document
as part of the change that adds it; an unassigned machine-readable surface
is a defect.
```

#### Load-Bearing Decisions

- **Identity** — two envelopes are "the same shape" iff they share a
  `schema_version` major. Minor differences are additive by C1 and a
  conforming parser cannot distinguish them.
- **Wire / byte format** — `schema_version` is a string
  (`"1.0"`), not a pair of integers and not a number: it follows
  OpenTofu's `format_version` form exactly, and a float would make `1.10`
  sort below `1.9`.
- **Naming** — `schema_version`, rejecting `format_version` (OpenTofu's
  own name) because this envelope carries a schema across every verb
  rather than one command's output format, and rejecting `version`
  outright — the `version` verb's payload already means build identity,
  and a key meaning two things on one wire is the defect this RDR exists
  to prevent.
- **Selection / predicate** — when a new finding could be introduced at
  either severity, `info` is chosen unless C3's re-attribution clause
  applies. The tie goes to the non-breaking tier by default.

**Sibling-path check.** One adjacent version discriminator already exists:
the model TOML's `[model] version`, an integer pinned at `1`, enforced by
`internal/table/load.go` — "`version %d; this RDR accepts version 1 only`",
refusing under `table::CatUnsupportedVersion`. It is deliberately NOT
reused or matched here. It versions the *input* document an author writes,
where an integer suffices because there is one accepted value and no
additive-vs-breaking distinction to express. C1 versions the *output*
schema a consumer parses, where that distinction is the entire point —
hence `MAJOR.MINOR`. The two are different surfaces with different
audiences, so the divergence in form is intended; what they share is the
rule that an unsupported major is rejected outright. Searched for an
existing output-schema version (`grep -rn "format_version\|schema_version"`
over the Go sources): none exists.

#### Illustrative Code

A terminal success envelope as it ships during `0.x` (C1: the schema
version begins at `"0.1"`; it reads `"1.0"` only from the 1.0.0 release
onward):

```json
{"type":"ok","schema_version":"0.1",
 "data":{"findings":[]}}
```

A refusal, showing the schema version on the failure record too. Note the
shape: the failure record is the BARE `CLIError`, with no wrapper and **no
`type` key at all** — not `"type":"failed"`. That is `0005:C1`'s decision
("Failures MUST use the existing CLIError JSON/text envelope"), confirmed
against a live run under {SPIKE_DIR}`envelope-baseline.md`, whose observed
failure key set is `code,detail,findings,message,param`:

```json
{"schema_version":"0.1",
 "code":"graph-lint-failed",
 "message":"the model carries blocking graph-lint findings",
 "findings":[{"code":"graph-overlap","severity":"blocking","rule":"r3"}]}
```

`schema_version` is therefore what a consumer branches on to know which
record shape it holds, since the two records share no discriminator: `ok`
carries `type`, a refusal carries `code`. A consumer distinguishes them by
`code`'s presence, exactly as `internal/cli/flow_input.go::planEnvelope`
already does.

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Severity partition (`blocking`/`info`) gating the exit code | Predecessor (RDR 0006) | Available | C3 rides it; no new machinery. `0006:C17`'s closure at four is opened to append-only as part of this change (A1, ruled). |
| Terminal envelope every verb routes through | Existing (`internal/cli/respond`) | Available | C1 adds one field at the single gateway rather than per verb. |
| Enumerable vocabularies to assign tiers to | Existing (`table.Categories()`, `graphlint.BlockingCodes()`) | Available | C4 assigns tiers to surfaces that already enumerate themselves. |
| Release-notes discipline naming promotions | This RDR | Introduced | C3's disclosure obligation is process, asserted by review not by a test. |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Envelope field carrier | `internal/cli/respond::Success` | `Data`/`Notes`/`Warnings` are `omitempty`; `schema_version` must not be | Extend | One non-omitempty field added at the gateway; C1 fixes its form. |
| Failure envelope carrier | `internal/cli/clierr::CLIError` | Marshals itself; no wrapper to nest under | Extend | Same field, added symmetrically so one consumer struct parses both. |
| Load-category vocabulary | `internal/table::Categories()` | Comment already says "the list's size is not a contract" | Reuse | C4 assigns `append-only`; the existing consumer rule is already correct and gets hoisted, not changed. |
| Lint finding taxonomy | `internal/graphlint/taxonomy.go` | Advisory tier declared CLOSED at four (`0006:C17`); severity is DERIVED (`severityFor`), so the code imposes no limit | Extend | C3 needs `info` open for new codes. Ruled: C17 becomes append-only; the only code-side blocker is one by-value test assertion (A1). |
| Stability prose | `docs/cli-output-contract.md` | Contains no release/upgrade/breaking language today | Extend | Becomes the promise's home. |

### Decision Rationale

The decisive factor is that **the fork as the seed framed it is not the
fork the prior art actually splits on**. The seed asks whether a new lint
finding code is a minor or a major change. Five independent peer tools
answer a different question: not *does the code exist* but *does it fire
by default*. typescript-eslint states the test directly — "A change to the
plugins shall be considered breaking if it will require the user to change
their config" (typescript-eslint.io/users/versioning) ⇒ the breaking edge
is the consumer's opted-in configuration, not the tool's vocabulary. Ruff
makes the same split structural: "New rules should always be added in
preview mode. New rules will remain in preview mode for at least one minor
release before being promoted to stable" (docs.astral.sh/ruff/versioning)
⇒ introduction is verdict-safe by construction, and promotion is the
disclosed event. golangci-lint draws it in the release-tier definition
itself — a minor release "might break your lint build because of newly
found issues" (golangci-lint.run/docs/product/roadmap) ⇒ new-finding
breakage is a named, expected minor-release risk, not a major bump.

That convergence is what makes C3 cheap rather than novel. intrastate
already has the two-valued severity partition, and `0006:C17` already
says "Advisory findings MUST NOT change the success disposition" ⇒ the
`info` tier is already a preview gate; C3 only states the policy that
governs leaving it. The alternative of building a separate `--preview`
flag was rejected precisely because that mechanism is already shipping
under another name.

The second factor is that this project has **already made this decision
four times, inconsistently, in code comments**. `table.Categories()` grew
under RDRs 0002, 0024, 0025 and 0028 ⇒ "a new code appears" is this
project's normal release event, so a policy classifying it as breaking
would make nearly every release a major bump. Meanwhile the word `closed`
is live in two incompatible senses on adjacent surfaces (Key Discoveries)
⇒ a consumer reading either comment and generalizing gets the other
wrong, which is why C2 retires the term rather than defining it.

The third factor is that binary version and wire-schema version answer
different questions and have never had the same answer. OpenTofu, carrying
Terraform's original language, versions the JSON format independently:
"We will increment the minor version, e.g. `1.1`, for backward-compatible
changes or additions… We will increment the major version, e.g. `2.0`, for
changes that are not backward-compatible" (opentofu
`website/docs/internals/json-format.mdx`) ⇒ C1 takes that rule verbatim in
form. Cargo independently corroborates the pattern — "The format is stable
and versioned. When calling `cargo metadata`, you should pass
`--format-version` flag explicitly to avoid forward incompatibility
hazard" (doc.rust-lang.org/cargo/reference/external-tools.html) ⇒ two
unrelated tools reached the same decoupling, which is why C1 is stated as
a schema version rather than a promise about the release number.

**Premortem.** Assume this shipped and failed. The most plausible
narrative: `intrastate` cuts v0.2.0 carrying C1, and the very release that
announces "your parse keeps working" is the one that breaks it — a
consumer validating the envelope against a strict schema (a Go struct
decoded with `DisallowUnknownFields`, or a JSON-schema check in CI) sees
an unexpected `schema_version` key and hard-fails. The promise's first act
is a violation of itself, and the agent audience is exactly the population
that hand-rolls strict parsers. The recommendation survives this, for two
reasons, but not unchanged. First, the exposure is bounded by the release
series itself: the first release is 0.1.0, and a `0.x` series carries no
SemVer compatibility guarantee, so C1's major-bump rule has nothing to
bind until 1.0.0. The whole `0.x` run is the free window, not just the
first tag. Second, C1 already carries the mitigation as a consumer
obligation ("MUST ignore object properties with unrecognized names"),
which is the clause that makes every *later* additive change safe. What
the premortem changes is the sequencing, and that is now a Pending
assumption (A2): `schema_version` must land before 1.0.0, or the project
spends its free window and then owes a major bump for the field that
exists to prevent major bumps. The honest reading is that this failure is
less acute than it first appears — the `0.x` series is a generous runway —
but it is also the kind of deadline a project misses by never treating any
particular release as the one. A second, weaker failure — that C3's disclosure obligation
is process rather than a test, so a promotion ships undisclosed — is
accepted as a known limit rather than designed around; it is recorded as
a failure mode, and the tier assignment in C4 is what a reviewer checks
against.

All three rejected alternatives (see Alternatives Considered) fail against
facts already on the ground rather than against preference: ALT1 against
four shipped growth events, ALT2 against what `llms.txt` already tells
agents, ALT3 against a strictness flag that does not exist yet. That is
why none of them is held open as a live option.

The argument above rests only on typescript-eslint, Ruff, golangci-lint
and OpenTofu, all fetched verbatim. ESLint is corroboration only — its
own new-rule classification could not be opened (Investigation), so
nothing here leans on it.

**Joint-decision check (three arms, all run).** Arm 1 (modify-anchors,
`rdr index --anchor-intersect`): no overlaps. Arm 2 (contract literals,
`rdr index --literal-intersect`): four overlaps against open peers, one of
them a genuine joint decision — RDR 0022 (Draft, `large`) proposes
`graph-terminal-unreachable` as a NEW BLOCKING finding code and states
"the advisory tier stays closed at four (`0006:C17`)", which is precisely
the rule C3 here would change. The other three are incidental vocabulary
sharing on one token each and are not joint decisions: 0021 shares
`--as=json` but adds no envelope field (its export rides `data`), 0014
shares `code` as a CI-oracle diagnosis surface, and 0017 shares `code`
with an existing cross-citation. Arm 3 (absence, manual): C3 converts no
refusal into an acceptance and removes no guard, so no Final peer relies
on a token this proposal stops saying; the C17 closure it does depend on
belongs to 0006, which is Implemented and therefore not edited — that
coupling rides to 7.1.

Premortem: survived (paragraph)
Ground-sweep: reopened → A3's evidence line (`DisallowUnknownFields` is in
production use at `internal/table/source.go::decodeStrict`, on the INPUT
path under `0002:C3`); 15 of 16 anchors CONFIRMED
Joint-check: fired → 0022 (home: `cli/0029 §Normative Contracts` C3)

**Joint decision, settled.** C3 in this record is the single normative home
for the introduce-at-`info` / disclose-on-promotion rule. RDR 0022 drops its
restatement of the C17 closure and CITES `0029:C3` instead, then either
enters `graph-terminal-unreachable` at `info` or claims C3's re-attribution
clause if that code re-attributes an existing refusal. This record is `mid`
against 0022's `large` and its whole subject is the cross-version promise,
so the rule carries less blast radius here; cite-don't-restate also avoids
the copy-drift that forces 7.1 demotions. The `0.x` framing above softens
the immediate stakes — through `0.x` neither record's tier is a guarantee —
but the rule still has to have one home before 1.0.0, and this is it.

## Alternatives Considered

### Alternative 1: Freeze the vocabularies within a major

**Description**: Treat every machine-readable vocabulary as fixed for the
life of a major version. Any new refusal code, load category, or lint
finding waits for the next major bump. The promise to an agent becomes
maximally simple: exhaustively switch on anything, and an unknown member
is always a bug.

**Pros**:

- The strongest possible guarantee, and the easiest one to state in a
  sentence an agent author can act on without reading further.
- A consumer's exhaustive `switch` is correct by construction, with no
  tolerance clause to get wrong.
- No promotion policy needed — C3 disappears entirely.

**Cons**:

- Contradicts four shipped growth events. `table.Categories()` grew under
  RDRs 0002, 0024, 0025 and 0028, all pre-1.0.
- Prices ordinary work as a major bump: under this rule, adding one
  `edit`-carrier category (RDR 0028's six) would have forced v2.0.0.
- Majors would be cut so often that the major number stops carrying the
  signal it exists to carry, and consumers learn to ignore it.

**Reason for rejection**: It is refuted by this repo's own history rather
than by preference. The project's normal release event is exactly the
event this alternative classifies as breaking, and no peer tool surveyed
adopts it — golangci-lint explicitly places new-linter additions in minor
releases and even classifies *linter removal* as non-breaking.

### Alternative 2: Declare all machine-readable output unstable

**Description**: State that `--as=json` carries no cross-version
guarantee at all. Consumers pin an exact version or accept the risk.
This is `gh`'s de facto position — no documented stability guarantee was
found for its `--json` output — and shellcheck's, where the community
answer is "pin a specific version to avoid surprise build breaks".

**Pros**:

- Zero ongoing obligation; no policy to maintain, no promotion discipline,
  no tier assignments to keep current.
- Honest about a pre-1.0 project whose envelope RDRs 0023 and 0028 moved
  recently.
- Matches at least two prominent peers, so it is not an eccentric position.

**Cons**:

- Directly contradicts `llms.txt`, which tells agents "The binary is
  authoritative and self-describing. Prefer asking it over reading any
  file here" — the project has already invited the coupling.
- Pushes the cost onto every consumer, who must each independently
  discover which surfaces move.
- Wastes mechanism already built: the severity partition and the
  enumerable vocabularies exist and already behave well.

**Reason for rejection**: The project cannot simultaneously tell agents to
bind tightly to the binary's output and decline to say what that output
promises. The seed's own framing — that an agent hard-fails or mis-parses
where a human adapts — is the argument against this option.

### Alternative 3: A strictness opt-in flag carrying a stability disclaimer

**Description**: Adopt clippy's shape — add a `--fail-on-advisory` (or
equivalent) flag, and disclaim stability for consumers who opt into it,
exactly as clippy states "we do not guarantee stability under
`#[deny(lintname)]`". Default-path consumers get a strong promise;
strict-mode consumers accept new-finding risk.

**Pros**:

- Clean separation: the consumers who want maximum signal are the ones
  who accept the churn, and they opted in explicitly.
- Precedent in a widely-used tool, with published wording to borrow.
- Would let advisory findings become actionable in CI without forcing
  every consumer to care.

**Cons**:

- The flag does not exist. `internal/cli/lint.go` registers `--model` and
  `--flow` and nothing else, so the disclaimer would attach to a surface
  with no implementation.
- Adds a user-visible flag to settle a documentation question, widening
  the change well past what the problem needs.
- Orthogonal rather than competing: once the flag exists, this disclaimer
  is an addition to C3, not a replacement for it.

**Reason for rejection**: Premature, not wrong. It solves a problem the
CLI does not yet have, and the clause it would add remains available the
day a strictness flag is introduced. Recorded here so that change inherits
the reasoning rather than re-deriving it.

### Briefly Rejected

- **Version the schema per verb rather than per envelope**: each verb's
  payload would carry its own version, but the envelope is the thing
  consumers parse generically, and per-verb versions multiply the state an
  agent must track without answering "can my parser read this".
- **Derive `schema_version` from the binary version**: collapses the two
  questions C1 exists to separate, and would force a schema major on every
  product major even when the wire did not move.
- **Use an integer `schema_version`**: loses the additive/breaking
  distinction that makes the minor component useful; OpenTofu's
  `MAJOR.MINOR` string is the shape with the field experience behind it.
- **Put the promise in a new `docs/compatibility.md`**: a third document
  competing with `docs/cli-output-contract.md` and `llms.txt` for the same
  reader, when the contract doc is already named authoritative for the
  envelope and already in the agent's read path.

## Context

### Background

Raised while judging whether the project is stable enough to cut its first
release. The build and publish machinery is already in place and unexercised:
`.goreleaser.yaml` builds five targets, `.github/workflows/release.yml`
publishes a GitHub Release on a `v*` tag after re-running `make check` against
the tagged commit, `ci.yml` carries a gated snapshot job, and
`internal/version` stamps build identity via ldflags with a VCS fallback that
`make release-check` guards. No tag has ever been cut, so
`intrastate version --as=json` currently reports a pseudo-version.

What is absent is not machinery but a stated promise. The distribution surface
is deliberately consumer-facing already — `.goreleaser.yaml` calls the archive
`name_template` "a compatibility surface, not cosmetics" — and `llms.txt`
tells agents the binary is authoritative and to prefer asking it over reading
any file. Neither says what stays true across versions.

The seed of the fork is in the existing records: RDR 0017 gave each lint
finding its own code identity, which is what makes "a new finding code
appeared" a distinguishable event a consumer can trip over. RDR 0023 and 0028
moved envelope shapes recently enough that the wire contract is visibly still
in motion.

### Technical Environment

Go 1.26.3, `github.com/cwensel/intrastate`, a cobra command tree over a pure-Go
binary (no cgo). The relevant surfaces:

- `internal/cli/respond/` — the output gateway every verb routes through;
  `Success{Type, Notes, Warnings, Data}` is the terminal JSON envelope.
- `internal/cli/clierr/` — structured `CLIError` plus exit-code mapping; the
  refusal-code vocabulary and the `findings` carrier.
- `internal/table/load.go` — refuses a model whose `[model] version` is not
  `1`, reporting the `unsupported_version` category from
  `internal/table/category.go`.
- `internal/version/` — ldflags-stamped build identity, with a
  `debug.ReadBuildInfo` fallback.
- `docs/cli-output-contract.md` — the worked payloads and envelope rationale;
  `docs/cli-reference.md` and `llms.txt` are generated from the command tree
  and gated against staleness by `make docs-check`.

Distribution today is GitHub Releases plus `go install`. Package-manager reach
(a Homebrew tap, an install script, any npm wrapper) is a separate, dependent
question — it widens the audience that would be holding whatever promise this
RDR settles, but it does not decide the promise.

## Research Findings

### Investigation

Two passes, both recorded under
`docs/rdr/0029-version-promise-on-machine-readable-output/evidence/research/`.
The in-repo pass (`in-repo-prior-art.md`) read the four machine-readable
vocabularies at HEAD — `internal/graphlint/taxonomy.go`,
`internal/table/category.go`, `internal/cli/clierr/clierr.go` and
`internal/cli/lint.go` — plus `docs/cli-output-contract.md`, `llms.txt`
and `.goreleaser.yaml`, looking for how growth and closure are declared
today. The external pass (`prior-art.md`) asked what established CLIs
publish about machine-readable output stability, and specifically how each
handles a new lint rule making previously-clean input newly fail: peers
opened were golangci-lint, typescript-eslint, Ruff, clippy (RFC 2476),
OpenTofu's `json-format.mdx`, Cargo's external-tools reference, and
negative results for shellcheck and `gh`. ESLint's own classification page
could not be opened; the local checkout confirms only the operational use
of `semver-minor` at `docs/src/maintain/manage-releases.md:49`.

### Key Discoveries

- **Documented** — The project already decides this per-vocabulary, in
  code comments, and the decisions disagree. The comment heading
  `graphlint::CodeCoverageClosedByEscape`'s const block declares the
  advisory tier "CLOSED at these four"; the doc comment on
  `table::Categories()` calls its set "the closed load-category set"
  while a sibling comment says "The list's size is not a contract". Two
  incompatible senses of `closed` on adjacent machine surfaces.
- **Documented** — A new code appearing is this project's *normal*
  release event, not a hypothetical: `table.Categories()` grew under RDRs
  0002 (25 members), 0024 (+3), 0025 (+6) and 0028 (+6), all before any
  tag was cut.
- **Documented** — The mechanism for verdict-safe introduction already
  ships. `graphlint::Report.Blocking()`/`Advisory()` partition on a
  two-valued severity, and only the blocking half reaches
  `respond.Fail` under `graphlint::AggregateCode`; the advisory half
  rides the SUCCESS payload. `0006:C17` already binds it:
  "Advisory findings MUST NOT change the success disposition."
- **Documented** — Five peer tools split this fork at default-enablement,
  not at code existence. The cleanest statement is typescript-eslint's:
  "A change to the plugins shall be considered breaking if it will
  require the user to change their config." Ruff makes it structural
  ("New rules should always be added in preview mode"), golangci-lint
  puts it in the release-tier definition ("Minor release (might break
  your lint build because of newly found issues)").
- **Documented** — OpenTofu is direct precedent for C1, down to the
  consumer rule. `opentofu/website/docs/internals/json-format.mdx:15-26`
  states both halves verbatim: "**Ignore any object properties with
  unrecognized names** to remain forward-compatible with future minor
  versions" and "**Reject any input which reports an unsupported major
  version**." The field sits on the response envelope beside the binary
  version, not derived from it — `jsonplan/plan.go:31-58` carries
  `format_version` and `terraform_version` as siblings — and is
  **per-surface**: `jsonplan` is at `"1.2"` while `jsonstate` and
  `jsonprovider` are at `"1.0"`, three independently incremented output
  schemas in one binary. Cargo's `cargo metadata --format-version` is the
  second instance.
- **Documented** — The pattern is established but **minority**, which is
  the honest framing of how much precedent C1 carries. Present in tools
  whose JSON is contractually consumed by other programs (OpenTofu;
  kubebuilder's plugin protocol, `pkg/plugin/external/types.go:26`;
  roborev and beads, both integer-valued `schema_version`). Absent from
  gh, goreleaser, golangci-lint, hugo, consul, prometheus, etcd and
  semgrep — tools whose `--json` is treated as a convenience for humans
  and scripts. This CLI's output is agent-consumed, which is what puts it
  in the first group.
- **Documented** — The 0.x instability claim A2 rests on has a normative
  Go-toolchain statement, which is a stronger citation here than SemVer's
  own §4: the Go Modules Reference (`go-modules-reference.pdf`, p.4) —
  "A version is considered **unstable** if its major version is 0… Unstable
  versions are **not subject to compatibility requirements**. For example,
  v0.2.0 may not be compatible with v0.1.0." SemVer 2.0.0's own wording is
  not in the local corpora and is deliberately not quoted from memory.
- **Documented** — C1's tolerant-reader rule is textbook, and its
  asymmetry with `0002:C3`'s strict input is only PARTLY grounded.
  Daigneau, *Service Design Patterns* (p.244) names the Robustness
  Principle / Postel's Law and the Tolerant Reader; DDIA (p.121) gives the
  mechanism ("it can simply ignore that field… This maintains forward
  compatibility"). Both justify tolerant-on-receive / validated-on-send.
  Neither states the specific inversion this RDR asserts — strict
  reject-unknown-keys on *author-written config* so a human learns of the
  typo. That leg is a defensible extension of "senders should use schema
  validation before sending" with the CLI validating on the author's
  behalf, but it is an inference, not a citation, and C1's asymmetry
  paragraph should be read as this project's own decision rather than as
  received practice.
- **Documented** — The authoritative wire-contract document says nothing
  about versions: `docs/cli-output-contract.md` contains zero occurrences
  of "release", "upgrade" or "breaking", while `llms.txt` tells agents
  "The binary is authoritative and self-describing. Prefer asking it over
  reading any file here."
- **Documented** — No strictness opt-in exists to hang a clippy-style
  disclaimer on: `internal/cli/lint.go` registers only `--model` and
  `--flow`.
- **Verified at Stage 4** — All four assumptions now rest on source, a
  ruling, or a live run, not on assumption: `0006:C17`'s closure is opened
  by ruling (A1); no release exists yet (A2, `git tag --list` empty and
  `internal/version::resolve` pinned to `"dev"`); every envelope DECODER
  that exists is tolerant (A3, `flow_input.go::planEnvelope`), the binding
  constraint being `0023`'s byte-identity golden, which this RDR re-captures;
  and the three tiers partition all fourteen tiered vocabularies (A4), seven
  of which C4 had not assigned and now does.

## Trade-offs

### Consequences

- Agents get a stated, per-surface promise where none existed, in the
  document they are already directed to read.
- Four divergent code-comment policies collapse into one vocabulary, and
  the ambiguous word `closed` is retired from this role.
- The project gains an obligation it did not have: every new
  machine-readable surface now owes a tier assignment, and every
  `info` → `blocking` promotion owes a release note.
- `schema_version` is one added field on every envelope — a small,
  permanent output cost on a CLI that RDR 0023 established is measured in
  output bytes.
- Consumers must be told to ignore unknown keys, which is a promise only
  as strong as their compliance; a strict parser is outside this RDR's
  reach.

### Risks and Mitigations

- **Risk**: 0006 is not amended before a fifth advisory code is added, so
  `0006:C17` still reads "closed at" four while this RDR's C4 calls the set
  `growing` — two Final-or-Implemented records disagreeing on one
  vocabulary's openness (A1, ruled but not yet executed).
  **Mitigation**: the amendment and the `want`-literal edit at
  `findings_0006_test.go:288-293` are prerequisites of this RDR's
  implementation, not follow-ups; 7.1 cluster reconcile over {0006, 0029}
  re-reads the two records against each other before either locks further.
  The failure is visible — `TestReq74` goes red on the fifth code — so it
  cannot ship silently.
- **Risk**: `schema_version` is added after 1.0.0, making the
  compatibility field itself a breaking change (A2).
  **Mitigation**: land it during `0.x`, where incompatible change is free;
  the Implementation Plan's Phase 1 puts it in the first tagged release so
  the runway is never the thing being spent.
- **Risk**: The promotion disclosure in C3 is process, not a test, so a
  promotion ships unannounced.
  **Mitigation**: C4's tier assignment is the reviewable artifact — a
  severity change shows up as a diff against a declared tier rather than
  as an invisible constant edit. Accepted as a known limit.
- **Risk**: The tier vocabulary is adopted in docs but not in the code
  comments it replaces, leaving the contradiction live.
  **Mitigation**: Phase 1 includes replacing the `closed` wording at the
  two named sites; the term's absence is checkable by grep.

### Failure Modes

- **Breaks visibly**: a consumer with a strict parser fails on the
  release that adds `schema_version` (A3) — immediate, loud, diagnosable
  from the parse error naming the unexpected key.
- **Fails silently**: an `info` code is promoted to `blocking` without a
  release note. A consumer's CI goes red on a model nobody edited, and the
  only way to diagnose it is to diff the severity constants between two
  builds. This is the residual risk C3's disclosure obligation reduces but
  does not eliminate.
- **Fails silently**: a new machine-readable surface ships without a tier
  assignment, so it has no stated promise and nobody notices until a
  consumer depends on the wrong assumption. Diagnosed by auditing C4's
  list against the emitted vocabularies.
- **Recovery path**: for a wrongly-promoted code, demote it to `info` in a
  patch release — the demotion cannot break anyone, since `0006:C17`
  guarantees `info` does not change dispositions. For a missing tier
  assignment, assign it; the assignment is documentation, not behavior.
- **Diagnosis**: `intrastate lint --as=json` reports each finding's
  `severity` on the wire, so a consumer can compare severities across two
  builds directly rather than inferring from exit codes.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] `0006:C17` amended from "closed at" its four members to append-only,
      with the `want` literal at `findings_0006_test.go:288-293` edited in
      the same change (A1, ruled 2026-09-11). C3's wording depends on it;
      nothing else in this plan does.
- [ ] Still pre-1.0.0 (A2) — the `0.x` runway this plan assumes. The
      first release is 0.1.0; nothing here waits on 1.0.0 except the
      schema's own promotion to `"1.0"`.

### Minimum Viable Validation

An agent pins version N, parses the envelope, and survives N+1 across
exactly the event the seed named — a new lint finding code appearing.

1. Build at the current HEAD with `schema_version` implemented; run
   `intrastate lint --model <clean-model> --as=json` and record the full
   envelope. It reports `"schema_version":"0.1"` and no findings.
2. Add a new graph-lint finding code at severity `info` that fires on the
   model from step 1.
3. Re-run the same command against the same unmodified model.
4. **Expected end state**: the exit code is unchanged (0), the envelope's
   `type` is still `ok`, and the new finding appears in `data.findings`
   carrying `"severity":"info"`. The schema version's MAJOR is unchanged —
   a growing-vocabulary member is an additive change, so C1 moves the
   minor (`"0.1"` → `"0.2"`) and never the major. That is the observable
   the promise rests on: a consumer branching on `type` and the exit code
   observes no difference, a consumer reading findings sees one more
   entry, and a consumer gating on the major keeps parsing.
5. Promote that code to `blocking` and re-run: now the exit code changes
   from 0 to 2 and the record shape changes — the `ok` envelope is replaced
   by the bare `CLIError` refusal, which carries `code` and no `type` key
   (the observed failure key set is `code,detail,findings,message,param`;
   see {SPIKE_DIR}`envelope-baseline.md`). That is a more disruptive
   change for a consumer than a `type` value flipping to `"failed"`, which
   is what makes promotion the verdict-changing event C3 requires be
   disclosed, and introduction not.

This walks the whole promise end to end: C1's field, C2's tolerance rule,
and C3's introduce-at-`info` / disclose-on-promotion split.

**`oracle`** — what each MVV step fails on, and its negative control. Step 4
is the promise's load-bearing observation, so it may not pass by absence of
error:

| MVV step | Fails if X is wrong, because Y | Negative / failing control |
| --- | --- | --- |
| 1 emit baseline | `schema_version` absent or not `"0.1"` — the key set is compared exactly, so a missing or misspelled field is a diff, not a silent pass | pre-change baseline `data,type` has no such key; the run before the change MUST fail this row |
| 4 additive event | exit code ≠ 0, `type` ≠ `ok`, the new finding missing from `data.findings`, or its `severity` ≠ `info` — each asserted by value, not by absence of error | step 5, the same code at `blocking`, MUST flip exit to 2 and replace the envelope; if step 5 does not differ, step 4 proved nothing |
| 4 version movement | major moved, or minor did not — read off the emitted string, not inferred from the change class | a breaking change in the same position MUST move the major; a no-op release MUST move neither |
| 5 promotion | exit stays 0 or the record keeps a `type` key — the shape change IS the disclosed event C3 governs | step 4's `info` run is the control: same input, same code, different severity, different verdict |

**`disposition`** — the input classes this RDR assigns outcomes to. The
promise is exactly that the first row never becomes the second without the
release-note disclosure C3 requires:

| Input class | Exit code | Event / error | Artifact minted | Silent vs loud |
| --- | --- | --- | --- | --- |
| clean model | 0 | none | `ok` envelope, `data.findings` empty | silent |
| model firing an `info` code | 0 (unchanged) | none — advisory rides the success payload | `ok` envelope, one more `findings[]` entry | loud in payload, silent in verdict — C3's guarantee |
| model firing a `blocking` code | 2 | `graph-lint-failed` aggregate | bare `CLIError`, no `type` key | loud |
| an `info` code promoted to `blocking` | 0 → 2 for input that previously passed | `graph-lint-failed` | envelope shape changes | loud, AND MUST be disclosed in release notes naming the code (C3) |
| envelope carrying an unrecognized key | n/a — consumer side | none | — | silent by C1's tolerant-reader rule |
| envelope reporting an unsupported major | n/a — consumer side | consumer MUST reject | — | loud, by C1 |

**`trace`** — the assertions in force at each step of the MVV's end state,
walked for joint satisfiability. Witness values are the emitted envelope's
own bytes:

| Step | Assertions in force | Witness |
| --- | --- | --- |
| emit `ok` at 0.x | C1 (field present, non-`omitempty`, `MAJOR.MINOR`); C1's 0.x clause (major `0` ⇒ unstable) | `{"type":"ok","schema_version":"0.1","data":{…}}` |
| add `info` code | C3 (introduce at `info`); `0006:C17` via C3 (advisory MUST NOT change success disposition); C4 `growing` on advisory codes | exit 0, `type` still `ok`, `"severity":"info"` |
| version moves | C1 (minor increments on additive change); C2 `growing` (new member MAY fire on input that previously produced no finding) | `"0.1"` → `"0.2"`, major unchanged |
| refusal path | `0005:C1` (bare `CLIError`, no wrapper, no `type`); C1 (field on BOTH terminal records) | `{"schema_version":"0.1","code":"graph-lint-failed",…}` |
| consumer reads either record | C1 (tolerate unknown keys; reject unsupported major); the two records share only `schema_version`, so `code`'s presence discriminates | `planEnvelope` branches on `code`, unchanged by the added field |

No CONTRADICTION row: the 0.x clause and the minor-increment rule coexist
because C1 states the minor still moves while the major is `0`, and the
`growing` tier is what licenses the step-4 finding to appear on unchanged
input without that being a breaking change.

### Phase 1: Code Implementation

#### Step 1: Carry `schema_version` on both terminal records

Add the field at the single output gateway — `respond.Success` and
`clierr.CLIError` — so every verb inherits it without per-verb work, and
both terminal records carry it. Neither is `omitempty`. The records stay
structurally asymmetric (a refusal has no `type`); `schema_version` is the
field they share, not a step toward unifying them.

Re-capture `0023`'s envelope golden in this same commit:
`docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`
is asserted by byte-identity at `internal/cli/flow_mvv_0023_test.go`, so a
non-`omitempty` field fails it until the baseline is retaken. The
re-captured golden must differ from its predecessor by exactly the one
`schema_version` key (`0029:S3`).

#### Step 2: Retire the ambiguous `closed` wording

Replace the two conflicting uses with the declared tier names at
`internal/graphlint/taxonomy.go` and `internal/table/category.go`, so the
code comments and the contract doc say the same thing.

#### Step 3: Assert the tier rules the consumer is promised

A test per tier: that append-only vocabularies are never asserted by
cardinality or ordinal, that the frozen sets match their declared members,
and that an `info` finding leaves the success disposition untouched.

### Phase 2: Operational Activation

#### Activation Step 1: Publish the promise where agents already read

Write the tier table and the C1/C3 rules into
`docs/cli-output-contract.md`, and point `llms.txt` at it — the two
documents already in an agent's read path.

#### Activation Step 2: Cut 0.1.0 with the tiers already declared

The field ships in the first `v*` tag (A2). Before tagging, confirm
`make docs-check` passes so the generated reference does not contradict
the new prose, and that the published text says plainly that `0.x` carries
no guarantee — the tiers are the 1.0.0 intent, declared early so reaching
it costs nothing.

#### Activation Step 3: Promote the schema to `"1.0"` at 1.0.0

The one deferred step. When the binary reaches 1.0.0, the schema version
moves to `"1.0"` and the tiers stop being intent and start being
guarantees. No code change is expected at that point — that is the test of
whether this RDR worked.

### New Dependencies

None. Every mechanism this RDR relies on — the severity partition, the
enumerable vocabularies, the output gateway — already ships.

## Validation

### Testing Strategy

The matrix the verified assumptions imply. Each row names the code path
Stage 4 read or the spike output that backs it; the normative fixtures are
the two envelope shapes captured under
{SPIKE_DIR}`envelope-baseline.md`.

1. **Scenario**: The `ok` envelope carries `schema_version`, non-`omitempty`.
   **Expected**: `intrastate lint --model <clean-model> --as=json` emits
   top-level keys `data,schema_version,type` with `"schema_version":"0.1"`.
   The pre-change baseline is `data,type` (captured), so this row is the
   one added key and nothing else. Backed by
   `internal/cli/respond/respond.go::Success`.

2. **Scenario**: The refusal record carries `schema_version` too, and keeps
   its shape.
   **Expected**: a provoked failure emits `code,detail,message,param,
   schema_version` (or with `findings` where the failure aggregates) — the
   captured baseline `code,detail,message,param` plus the one field. No
   `type` key appears; the record stays the bare `CLIError`. Backed by
   `internal/cli/clierr/clierr.go::CLIError` and
   `internal/cli/respond/respond.go::Fail`.

3. **Scenario**: The existing in-repo envelope consumer survives the added
   field.
   **Expected**: `--plan` still decodes a `flow resolve --as json` envelope
   that now carries `schema_version`, unchanged. This is the A3 regression
   guard for the strictest DECODER the repo holds —
   `internal/cli/flow_input.go::planEnvelope`.

   The repo's other envelope coupling is tighter and is not a decoder:
   `internal/cli/flow_mvv_0023_test.go` asserts byte-identity against
   `0023`'s checked-in golden. **Expected**: that golden is re-captured in
   the same commit that adds the field, and the re-captured file differs
   from its predecessor by exactly the one `schema_version` key. A diff
   showing any other key changed means this RDR altered an envelope it
   promised not to touch.

4. **Scenario**: An `info` finding leaves the success disposition untouched
   (C3, riding `0006:C17`'s guarantee).
   **Expected**: exit 0, `type` still `"ok"`, the finding present in
   `data.findings` with `"severity":"info"`. Backed by
   `internal/graphlint/taxonomy.go::severityFor` — severity is derived from
   `!IsBlocking`, so a new advisory code needs no severity plumbing.

5. **Scenario**: Promotion to `blocking` IS the verdict-changing event.
   **Expected**: exit 2 and the record shape changes from the `ok` envelope
   to the bare `CLIError`. This is the observable C3 requires be disclosed.

6. **Scenario**: `frozen` sets match their declared members exactly.
   **Expected**: a by-value assertion per frozen vocabulary in C4 —
   including the five this RDR newly assigned. `data.escape_class` already
   has one
   (`internal/resolve/guard_atoms_test.go::TestReq79_NoSixthRefusalKindIsMinted`);
   that test is the pattern the others follow.

7. **Scenario**: `append-only` vocabularies are never asserted by
   cardinality or ordinal.
   **Expected**: membership-and-uniqueness assertions only. The existing
   pattern to extend is
   `internal/table/command_carrier_0025_test.go::TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`,
   which deliberately avoids count and index goldens — not a golden file.
   No golden/snapshot test of the JSON envelope exists in the repo today,
   and this RDR does not add one: a byte-golden over the envelope would
   itself assert the cardinality C2 forbids consumers from asserting.

8. **Scenario**: `0006:C17`'s amendment lands with the fifth advisory code
   (A1).
   **Expected**: `internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
   is updated in the same change, and 0006's C17 text reads append-only.
   Left unamended, `TestReq74` goes red — which is the intended visible
   failure, not a flake to suppress.

9. **Scenario**: The retired `closed` wording does not survive (C2's last
   clause).
   **Expected**: no occurrence of `closed` describing a tier in
   `internal/graphlint/taxonomy.go` or `internal/table/category.go`. A grep
   assertion is sufficient and is the cheapest guard against the ambiguity
   this RDR exists to remove.

### Performance Expectations

Omitted deliberately: no alternative in this RDR was weighed on empirical
performance grounds. The change adds one short string field to two structs
already marshalled once per invocation, through the single encoder
`internal/cli/clierr::WriteJSONLine`; there is no measurement that would
discriminate between the alternatives considered.

Byte-stability is NOT claimed for the envelope, so the determinism
checklist does not apply to it. That is deliberate and worth stating,
because C2 forbids exactly the assertion a byte-stable envelope would
invite: a consumer MUST NOT assert on a vocabulary's cardinality or a
member's ordinal position. The one place this project does make a
byte-equality claim — set-value read-back, via the non-HTML-escaping
encoder (`JDR 0001 §D13`) — is upstream of this RDR and untouched by it.

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

- [Requirements/standards with section numbers]
- [Dependency docs, source paths reviewed]
- [Dependency repos searched (clone + code search)]
- [Related issues, articles, discussions]
