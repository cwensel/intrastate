Model: claude-opus-5

# Persona 1 — Product Manager: does 0029 deliver the user outcome?

Owned set read in full: `0029:§problem-statement`, `0029:§approach`,
`0029:§decision-rationale`, `0029:MVV`.

**Widened, and what sent me.** The problem statement's user outcome is a
*discovery* outcome — the agent "needs to know, **before upgrading**, what
that version number said was allowed to change." Nothing in the owned set
says how the agent learns it; the only mechanism is one line of
`0029:§phase-2-operational-activation` (Activation Step 1). A silence has no
line range, so I widened to the Implementation Plan, `0029:C1`–`0029:C4`,
`0029:§consequences`, `0029:§failure-modes`, `0029:§briefly-rejected`,
`0029:§existing-infrastructure-audit`, and out of the record to the two
artifacts the RDR names as the delivery vehicle (`llms.txt`,
`docs/cli-output-contract.md`) plus `internal/cli/version.go`, because the
claim "the two documents already in an agent's read path" is a product claim
about a real file and is checkable.

Verdict is that the RDR solves the *governance* problem (what the project is
allowed to change) very well, and under-delivers the *consumer-side* problem
it opened with (how an agent finds out, and what it does with the answer).
Findings are severity-ranked.

---

## S1 — The stated user outcome is "know before upgrading," but every
## mechanism delivered is observable only *after* upgrading

Anchors: `0029:§problem-statement`, `0029:§approach`,
`0029:§phase-2-operational-activation` (Activation Step 1), `0029:MVV`.

`0029:§problem-statement` names the moment of need precisely: the agent
"needs to know, **before upgrading**, what that version number said was
allowed to change," and it "discovers the gap the hard way." That is a
pre-upgrade, decision-support outcome.

Every mechanism `0029:§approach` lists is post-upgrade:

- `schema_version` (`0029:C1`) is a field on an envelope the agent can only
  read by *running the new binary* — i.e. after the upgrade it was trying to
  decide about.
- the tier declarations (`0029:C2`, `0029:C4`) live in a markdown doc, not on
  any wire.
- the promotion-disclosure rule (`0029:C3`) binds release notes, which are
  the one artifact the RDR never says anyone writes, checks, or generates.

`0029:MVV` confirms the gap rather than closing it: its scenario is "an agent
pins version N, parses the envelope, and **survives** N+1" — survival after
the fact. Step 1 and step 3 both *run the binary*. No MVV step models an
agent asking "may I upgrade?" and getting an answer.

**Decision blocked.** Whether this RDR owes a pre-upgrade affordance at all.
Two live options it never puts on the table: (a) accept that the promise is a
document read by a human or a skill author at pin time, and say so — which
demotes the "agent hard-fails" framing in `0029:§problem-statement` to
motivation rather than requirement; or (b) owe a machine-readable surface
that answers the question before the upgrade (see S2). Until this is settled
an implementer cannot tell whether Activation Step 1 is the whole delivery or
a placeholder for it, and a reviewer cannot tell whether shipping C1–C4 with
no pre-upgrade path is the RDR working or the RDR missing its own target.

---

## S2 — The promise is published only as prose, in a document whose own
## index tells agents to prefer the binary over it

Anchors: `0029:§phase-2-operational-activation` (Activation Step 1),
`0029:BR4`, `0029:C2`, `0029:C4`, `0029:§existing-infrastructure-audit`.

Activation Step 1 is the entire user-facing delivery: "Write the tier table
and the C1/C3 rules into `docs/cli-output-contract.md`, and point `llms.txt`
at it — the two documents already in an agent's read path."

I checked both. They exist, and the read path is real —
`llms.txt` already links `docs/cli-output-contract.md`. But two product facts
cut against the claim:

1. `llms.txt` opens with a standing instruction *away* from the docs:
   "The binary is authoritative and self-describing. Prefer asking it over
   reading any file here — a file is a snapshot, the binary is the build you
   actually have," followed by `--help-all` and `--as=json`. The RDR is
   placing its single user-facing artifact in the one part of the read path
   the project tells agents to deprioritize, and `0029:§approach` asserts the
   opposite ("stated somewhere an agent can find it") without engaging this.
2. `llms.txt`'s one-line description of that file is "worked JSON payloads
   and the rationale behind the envelope shape." An agent scanning for
   compatibility, versioning, upgrade, or breaking-change guidance has no
   term to match on. The RDR requires the doc's *body* to change and never
   requires the *pointer* to change.

`0029:§existing-infrastructure-audit` states the doc "Contains no
release/upgrade/breaking language today" and decides "Extend" — which is
correct about the body and silent about discoverability.

`0029:BR4` rejects `docs/compatibility.md` on the grounds that the contract
doc "is already named authoritative for the envelope and already in the
agent's read path." That rejection is sound on the "third competing document"
argument and does not address the ordering instruction above it.

**Decision blocked.** What Activation Step 1's acceptance criterion is. As
written it is unfalsifiable — "write the rules into the doc" passes if one
paragraph lands anywhere in a 541-line file. It blocks a reviewer from saying
the promise shipped, and it blocks the choice between: amending the
`llms.txt` entry so the pointer names the promise; hoisting a short
compatibility section to a named anchor; or exposing the tier table from the
binary (`intrastate version --as=json`, or a `--help-all` section), which
would satisfy `llms.txt`'s own precedence rule and is the option `0029:BR4`
never considered because it only weighed document-vs-document.

---

## S3 — The RDR's own success test is deferred past the release that would
## prove it, leaving no in-scope evidence the outcome was delivered

Anchors: `0029:§phase-2-operational-activation` (Activation Step 3),
`0029:MVV`, `0029:§scope-verification`.

Activation Step 3 says: at 1.0.0 the schema moves to `"1.0"`, the tiers stop
being intent and become guarantees, and "No code change is expected at that
point — **that is the test of whether this RDR worked**."

The RDR therefore names its own success test and places it at an event
(1.0.0) that `0029:A2` establishes has not happened — no tag has ever been
cut, so this is an indefinite deferral past the entire `0.x` series. Nothing
in `0029:MVV` substitutes: the MVV proves the *mechanism* (a new `info` code
does not flip a verdict), which is a governance property, not the stated user
outcome that an agent "gets a stated, per-surface promise... in the document
they are already directed to read" (`0029:§consequences`, first bullet).

`0029:§scope-verification` is still the unfilled template block, so the gate
that exists to say "the MVV is in scope and executed, not deferred" has not
been answered for a record whose own worked-test sentence points at a
deferred event.

**Decision blocked.** Whether anything in this RDR's implementation is
verifiable as delivering the user outcome, versus only the plumbing. It
blocks writing the Scope Verification gate response honestly, and it blocks
an implementer from knowing whether "done" means C1–C4 in code plus a doc
edit, or something a consumer can be shown.

---

## S4 — The `0.x` reframing hollows out the promise for the whole period in
## which the RDR actually ships, and the problem statement is not reconciled
## to it

Anchors: `0029:§approach`, `0029:C2` (final paragraph), `0029:§problem-statement`.

`0029:§approach` is unambiguous: "this RDR is **not** writing a promise that
binds today," the series is `0.x`, "anything may move in any release," and
tiers are "declared and honored as intent, not guaranteed." `0029:C2`'s last
paragraph makes it normative: a consumer "MUST NOT" read a tier as a
guarantee already in force.

Set that against the opening user in `0029:§problem-statement`: an agent
author who "pins a released `intrastate` version," upgrades, and gets a
silent mis-parse or a red CI. Through `0.x` that user gets, from this RDR:
a version string they cannot rely on, tiers that are explicitly not promises,
and one clause that does bind — C3's disclosure obligation. The RDR knows
this and handles it honestly in `0029:C3` ("The disclosure obligation binds
during `0.x` as well... a consumer's pipeline going red is equally disruptive
at 0.4.0"). But `0029:§problem-statement` was never rewritten to match: it
still poses the fork as a live consumer-breakage problem, and
`0029:§approach` answers a different, later one — what the promise *will* be.
That is the "solves a different problem" shape, made visible by the two
sections disagreeing about when the user is helped.

The gap has a product consequence the record does not draw: during `0.x` the
only binding thing this RDR delivers to the agent author is a release note.
Release notes appear in `0029:C3`, `0029:§consequences`,
`0029:§capability-dependencies` ("Release-notes discipline naming promotions
| This RDR | Introduced"), and `0029:F2` — and in none of them is there a
file, a template, a generator, or a check. `0029:F2` books the consequence
("a consumer's CI goes red on a model nobody edited") as a residual risk;
`0029:§risks-and-mitigations` accepts it as "a known limit."

**Decision blocked.** Whether the sole user-visible deliverable of the `0.x`
period — the release note naming a promotion — needs any implementation at
all. Today it is an unowned process obligation introduced by a record whose
own Implementation Plan has no step that creates it. Settling this also
settles whether `0029:§problem-statement` should be re-scoped to the
governance problem the RDR actually solves, so that the record's opening user
and its approach are pointed at the same moment.

---

## S5 — `0029:C4`'s tier table is the consumer-facing product of this RDR,
## and is specified as a snapshot with no stated home or upkeep owner

Anchors: `0029:C4`, `0029:C2`, `0029:§failure-modes` (F3),
`0029:§phase-2-operational-activation`.

`0029:C4` opens "The tier assignments **at this RDR's implementation**" — a
point-in-time list. `0029:C2` says the tier "is recorded in
`docs/cli-output-contract.md` beside that vocabulary," while C4's own list is
a flat inventory in a decision record. Those are two different artifacts with
two different shapes (per-vocabulary annotation vs. a table), and the RDR
does not say which one the consumer reads, nor that they must agree.

`0029:C4` closes with the right rule — "A surface added later takes a tier
assignment in the same document as part of the change that adds it; an
unassigned machine-readable surface is a defect" — but `0029:F3` concedes the
detection story is manual: a missing tier is "Diagnosed by auditing C4's list
against the emitted vocabularies." `0029:A4`'s own history is the evidence
that this fails: the census found sixteen sets where C4 had named seven, and
the grounding lens caught two more inside `findings[]` after the first pass.
The gap recurred *twice inside the fix* by the record's own account.

**Decision blocked.** Which artifact is normative for a consumer — C4 in the
RDR, or the annotations in `docs/cli-output-contract.md` — and whether the
sixteen-surface census becomes an enforced artifact (an enumeration test in
the spirit of the `slices.Equal` assertions C4 cites for `Operators` and
`AdvisoryCodes`) or stays a prose list. It blocks Activation Step 1, which
cannot be executed without knowing whether it copies C4 or annotates
per-vocabulary, and it blocks the consumer promise, because a tier table that
silently drifts is worth less than no table.

---

## S6 — `0029:MVV` step 5 discloses a consumer-visible shape break that the
## RDR's promise never mentions to the consumer

Anchor: `0029:MVV` (step 5 and the `disposition` table row "a model firing a
`blocking` code").

MVV step 5 records that on promotion, "the `ok` envelope is replaced by the
bare `CLIError` refusal, which carries `code` and no `type` key," and calls
this "a more disruptive change for a consumer than a `type` value flipping to
`"failed"`." The `disposition` table repeats it: `bare CLIError, no type key`.

This is the single most consumer-relevant fact in the record — the two
terminal records are structurally asymmetric, so an agent branching on `type`
gets nothing to branch on when a run refuses. `0029:C1` states the asymmetry
and cites `0005:C1` for it, and instructs the consumer to reject an
unsupported major and ignore unknown keys. Neither C1 nor `0029:C2` tells the
consumer the thing it most needs: *discriminate on the presence of `code`,
not on `type`*. That rule appears once, in the MVV's `trace` table
("`code`'s presence discriminates"), which is validation scaffolding rather
than a published contract clause, and is not among the things Activation
Step 1 publishes (it publishes "the tier table and the C1/C3 rules").

**Decision blocked.** Whether the published promise includes a consumer
parsing rule for telling the two terminal records apart. It blocks the
content of Activation Step 1 and, downstream, blocks the S1 question — if the
promise is prose for a consumer, this is the sentence that prose most needs,
and it is currently only in the MVV.

---

## S7 — `0029:§scope-verification` and `0029:§proportionality` are unfilled
## template blocks on a `Profile: large` record

Anchors: `0029:§scope-verification`, `0029:§proportionality`,
`0029:§cross-cutting-concerns`, `0029:§metadata` (Profile).

Three of the five gate elements are still the template's bracketed guidance
rather than authored responses (consistent with the lint run in this evidence
directory, which reports five `placeholder:survived` findings). Two of them
bear directly on my question:

- `0029:§scope-verification` is where the record must say the MVV executes
  during implementation rather than being deferred — which is exactly the
  claim S3 puts in doubt.
- `0029:§proportionality` is where the "sole author of at most one
  independent load-bearing contract" split test is answered. From a product
  standpoint this record plausibly carries two separable deliverables: a wire
  field with an independent version line (`0029:C1`), and a stability-tier
  vocabulary plus its promotion policy (`0029:C2`/`C3`/`C4`). `0029:§metadata`
  already records `Profile: large` and enumerates three things in one breath.
  The split question is unanswered on a record that flags itself as large.

**Decision blocked.** Lock. A reviewer cannot accept the scope claim S3
questions, nor the one-seam claim, because neither has been written. Noting
this is a gate-mechanics finding rather than a prose-outcome one — I raise it
only where it gates my own question.

---

## Not findings (checked, and they hold)

- `0029:§decision-rationale`'s reframing of the seed's fork ("not *does the
  code exist* but *does it fire by default*") is the strongest product move in
  the record: it converts a would-be every-release-is-a-major policy into a
  cheap one, and it is grounded in four verbatim-fetched peers rather than
  preference. `0029:ALT1`'s rejection against "four shipped growth events" is
  a real fact on the ground, not a taste argument.
- The `schema_version` naming argument in `0029:D-naming` — rejecting
  `version` because "the `version` verb's payload already means build
  identity, and a key meaning two things on one wire is the defect this RDR
  exists to prevent" — is a correct user-facing call. I verified
  `internal/cli/version.go` emits build identity under the same `ok` envelope,
  so the collision would have been real.
- `0029:C1`'s tolerant-reader/strict-input asymmetry, and its explicit
  "A surface MUST NOT be given the input rule and the output rule at once,"
  is the right consumer-facing rule and is argued from audience rather than
  convenience.
