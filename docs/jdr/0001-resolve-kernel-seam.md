---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: open
cluster: 0002, 0003, 0004, 0005, 0006, 0007, 0008, 0009
labels: intrastate, internal-resolve, refusals, guard-evaluation, rdr-cluster
---

# JDR 0001 What the Resolver Can Trust, and What the User Sees When It Can't

*All Go identifiers, code spellings, and CLI transcripts below are
**non-normative** illustrations. The decisions are normative once resolved; the
RDRs own exact contracts. Extracted from the `0002-0009` cluster gate
(`docs/rdr/cluster-reconcile/0002-0009/report.md`, 13 pairwise scans plus a
whole-set critique, 2026-08-11).*

## Problem Statement

Eight RDRs jointly decide one thing: **when the resolver cannot decide, what
happens.** One RDR owns the table format (0002), one the guard vocabulary
(0003), one accessor execution (0004), one the CLI surface (0005), one graph
lint (0006). Three later RDRs — 0007 (guard evaluation domain), 0008 (recognized
tag key), 0009 (escape-row shape) — each landed on the `area:internal-resolve`
seam between them, two months after the first five locked.

The gate found the structural cause of every finding, and it is mechanically
checkable:

```
grep -cE 'RDR 0007|RDR 0008|RDR 0009' 0002…md 0003…md 0004…md 0005…md 0006…md
→ 0, 0, 0, 0, 0
```

**No older member references any newer one.** Citation traffic is entirely
one-directional. Meanwhile 0007 names `Obligation destination (NAMED): RDR
0003's implement stage` three times, routes A6b to 0004's implement stage, and
0008/0009 route further obligations to 0002/0004/0005. RDR 0003's Prerequisites
still read `- [x] All Critical Assumptions verified`, checked before those
obligations existed.

An obligation filed on a document that never received it binds nobody. That is
what this registry fixes: the shared halves live here, the RDRs cite them, and
the decisions that need a human get made once, below.

## Scope

Members are the eight Final-and-unimplemented RDRs above. RDR 0001 is
`Implemented` and out of scope — code is the source of truth there.

One JDR, not several, because these ten entries span a single connected set of
RDRs: nine of ten touch 0007, and the tenth (JD-10) reaches the rest through
0008. There is no clean partition, and per `README.md` §Identity a JDR is keyed
to the sharing set rather than to the theme or the gate run.

## How to use this document

Cite `JDR 0001 §JD-n`. **Never restate the mechanism** in an RDR — the cluster
gate treats a restatement as the defect. When an entry moves to `decided`, the
siblings it spans may carry `Final [joint decision → JDR 0001 §JD-n]`. While an
entry is `open`, no tolerance qualifier is warranted: the sibling is not free to
implement across it.

## Guiding principles

Standing commitments every resolution below must satisfy. Each is **derived from
normative text the members already lock** — none is invented here. They exist so
ten entries resolve into one coherent behavior instead of ten locally-reasonable
answers that do not compose. Where a principle rules an option out, the entry
says so.

**P1 — Missing artifact state is never masked.** RDR 0007's whole design turns
on this: the choice of guard domain "decides whether missing artifact state can
be masked behind an escapable refusal class," and it closes that path
deliberately. Any resolution that lets absent owned state reach the user as an
escapable `no_match`, or as a plan computed around it, violates the reason 0007
exists. *Binds:* D3, JD-7, JD-9, JD-2.

**P2 — The kernel refuses rather than guesses.** The typed-refusal taxonomy is
the kernel's answer to every undecidable case; RDR 0004 states the same shape
for accessors — "Indeterminate MUST be a refusal-class result, not a false allow
and not a false deny." A resolution that produces a confident answer from
incomplete input fails this even when no contract forbids it in terms. *Binds:*
JD-1, JD-4, JD-7, JD-10.

**P3 — Every user-visible failure is structured.** RDR 0005 routes every verb's
failure through `respond.Fail(cmd, *clierr.CLIError)` and sets `SilenceErrors`,
because the consumer is a skill parsing JSON. A failure that reaches the user as
anything else is a breach of the contract 0005 exists to provide, regardless of
which layer produced it. *Binds:* D1, JD-6, JD-8.

**P4 — A refusal names what was missing and who can fix it.** The "what" is
0007's stated user outcome — told plainly that the state needed to decide was
missing, not handed a plan that routed around it. The "who" is 0005's exit-code
grouping, which already encodes remedy: `GroupUserEnv` (exit 2) means *you* can
fix it, `GroupEnvUnavailable` (exit 3) means the environment could not answer.
A refusal whose code implies the wrong remedy is wrong even if it is
well-formed. *Binds:* D2, JD-8.

**P5 — A green lint means resolution succeeds.** RDR 0006 locks blocking
authority; the value of that authority is the promise it makes to an author
before runtime. If lint can prove a partition exhaustive and the runtime can
still refuse an assignment inside that proof, the promise is void and the
authority is theater. *Binds:* JD-4, JD-1.

**P6 — One home per contract; cite, never restate.** The engine doctrine —
"Two copies of one contract drift into self-contradiction… the cure is deletion,
not a lint to keep them in sync." This is why this registry exists, and it
constrains resolutions too: a resolution that asks two RDRs to each state the
rule has failed, however carefully worded. *Binds:* every entry.

**P7 — Pre-release, prefer the clean shape over the compatible one.** A
resolution should not carry a bridge, a shim, or a deprecation path to protect a
caller that does not exist. Verified: `internal/resolve` has one non-test file
and zero production importers — nothing under `cmd/` or `internal/cli/` reaches
it. *Binds:* D1, D2, JD-3, JD-10. **Provenance differs from P1–P6:** those quote
locked RDR text; this one rests on a standing project convention (unreleased, no
backward compatibility) that is *not* written down in the repo. Confirm it holds
before leaning on it, or promote it to repo text.

---

## Decisions requiring a call

Three entries are genuine forks that cannot be settled by re-reading the RDRs.
The rest of the registry records decisions already implied by the evidence, or
defers them to the implementation that will settle them honestly.

### D1 — How does a producer defect surface: panic, or the error envelope?

RDR 0007 makes this normative for the guard seam:

> A crash is the honest surface precisely because the case is unreachable for
> lint-passed guards: it is a programmer defect, and the kernel MUST NOT recover
> the panic into a verdict or refusal.

RDR 0005 makes the opposite normative for everything the user touches: every
verb routes failure through `respond.Fail(cmd, *clierr.CLIError)` under a
never-silent contract, and `internal/cli/root.go` sets `SilenceErrors`. Neither
RDR mentions the other; 0005 contains the string `panic` zero times.

The user-visible difference (non-normative):

```
$ intrastate flow next --model m.toml --tag stage=review
panic: guard mapping failed for row 4
goroutine 1 [running]: ...stack trace...
exit status 2
```

versus the structured envelope every other failure produces.

- **(a) Panic stands** — 0007's reasoning is that `Evaluate` returns a bare
  `GuardResult` with no error channel, so a mapping failure has no in-band
  surface. Cost: the one failure mode that indicates a *programmer* defect is
  the one that breaks the skill-integration contract, because a panic is
  unparseable by the JSON consumer 0005 exists to serve.
- **(b) Widen the seam** — give the evaluator an out-of-band error channel so
  the kernel can map it to a refusal. Cost: 0007 rejected this deliberately —
  "Adding a second channel would give the evaluator a way to report
  undecidability that the kernel's refusal taxonomy does not model." Reopening
  it reopens 0007's design.
- **(c) Panic at the seam, recover at the boundary — recommended** — 0007's
  clause binds the *kernel*, and the CLI is not the kernel. The `flow` verb
  recovers at its own boundary and renders a `GroupInternal` envelope with a
  stable code. 0007's "MUST NOT recover into a verdict or refusal" is honored
  literally: the recovery produces neither — it produces a crash report in
  structured form. Cost: one recover site, which must be narrowly scoped so it
  never swallows genuine panics into silence.

*Principles:* **P3** rules out (a) — a stack trace is not a structured failure,
and the skill consuming JSON cannot parse it. **P7** weakens the case for (b):
widening the seam is the compatible shape, but there is no caller to protect, so
the cost is design churn rather than migration. (c) satisfies P3 without
reopening 0007.

**Resolved:** _pending — needs a call._

### D2 — Is `flow-guard-unevaluable` about supplied facts, or about the view?

RDR 0005's code table declares:

> | guard cannot be evaluated from supplied facts | `flow-guard-unevaluable` | `GroupUserEnv` |

`GroupUserEnv` maps to exit 2 — the bad-input bucket, whose remedy is "fix your
input." RDR 0007 defines the identical condition over the *assembled view*,
which includes accessor-produced owned state the caller cannot supply, and its
whole premise is that the user is told the artifact state needed to decide was
missing. Under 0007, "supply the missing fact" is frequently not an available
remedy.

The sibling kind is worse off: `owned_state_unavailable` — which 0007 pins as
the *preferred* refusal whenever a survivor set carries both — has **no row at
all** in 0005's twelve-row table. A sweep for `owned.state|state-unavailable`
returns one unrelated prose hit.

- **(a) Keep `GroupUserEnv` for both** — cheapest, and wrong for the
  accessor-produced half: it tells the user to fix input they cannot reach.
- **(b) Split by cause — recommended** — `flow-guard-unevaluable` keeps
  `GroupUserEnv` when the missing key was caller-suppliable; a new
  `flow-owned-state-unavailable` row takes `GroupEnvUnavailable` (exit 3, the
  bucket that already means "the environment could not answer"). This uses
  0005's existing taxonomy rather than widening it, and it makes the exit code
  carry the remedy.
- **(c) One new group** — add `GroupInternal` and route both there. Simplest to
  state; loses the distinction between "you can fix this" and "you cannot," and
  RDR 0009 independently needs `GroupInternal` for its breach (JD-8), so this
  overloads one group with two unrelated meanings.

*Principles:* **P4** rules out (a) and (c) — both hand the user an exit code
whose implied remedy is wrong for at least one cause. (b) is the only option
where the code carries the remedy, and **P7** favors it further: adding a row to
an unreleased taxonomy costs nothing.

**Resolved:** _pending — needs a call._

### D3 — May `--tag` satisfy owned state?

RDR 0005 specifies `--tag name=value` as

> context already known to the caller.

with no allowlist, validation, or accessor-origin requirement. RDR 0007's
presence rule is provenance-blind by design: a key present under *any*
provenance decides a guard. Composed, an operator who hits a
`guard_unevaluable` refusal can re-run with `--tag` and get a plan — writes
computed from typed-in state, with nothing in the output marking it.

RDR 0007 records this as assumption A13, an "accepted exposure," on the ground
that the exposure is currently unreachable. Implementing 0005 is the event that
makes it reachable, so the bound expires exactly when the code lands.

- **(a) Accept** — `--tag` is a power-user affordance; document it. Cost: the
  cluster's headline guarantee (missing artifact state can never be masked) is
  bypassable from the command line, and the kernel already refuses the
  equivalent on the owned-state path.
- **(b) Provenance-tag the CLI channel — recommended** — tags arriving via
  `--tag` enter as `ProvenanceObserved` and may never satisfy an owned-state
  dependency; a guard needing owned state still refuses. Preserves 0007's
  guarantee at the only producer that can violate it, and matches the kernel's
  existing owned-state stance. Cost: `read-state` → `--tag` round-tripping stops
  working as a state-restoration trick.
- **(c) Allowlist** — enumerate which keys `--tag` may set. Cost: a list that
  must track every model's tag vocabulary; brittle, and it fails open for
  anything new.

*Principles:* **P1** rules out (a) outright — accepting the bypass is precisely
"missing artifact state masked," here by the user rather than by the resolver,
and the kernel already refuses the equivalent on the owned-state path. It also
rules out (c), which fails open for any key not yet listed. (b) is the only
option that holds P1 at the producer that can violate it.

**Resolved:** _pending — needs a call._

---

## Interface record

Each line is the joint decision plus the user-visible stake, its status, and its
provenance. Edit in place; normative once `decided`.

- **JD-1 Guard atom transport.** `Row.Guard` is one `string` per row, but RDR
  0003's Identity key is per-*atom* ("its position within `all` or `unless`"), so
  an ordinary multi-atom row needs N keys through a one-slot channel. RDR 0003
  disclaims the bridge: "This RDR introduces no encode/decode pair." RDR 0007's
  A10 verified this mapping as sufficient on a key whose cardinality it did not
  check. Stake: every guard with two conditions — the common case, and the case
  both RDRs' own example TOML shows. *Status:* `deferred → 0003 Phase 1
  (Predicate Model)` — the encoding is a property of the evaluator, unresolvable
  on paper. *(0007×0003 F1, blocks-impl)*

- **JD-2 Resolver flow: gate-then-count, or match-then-count?** RDR 0002 states
  the escape-reachability condition as a total function of the match count over
  non-escape rows; RDR 0007's aggregation veto refuses *before* any counting,
  and can kill a table where exactly one row matches and is decided true, purely
  because an unevaluable sibling also matched. An implementer reading 0002
  builds two stages; reading 0007, three. Stake: whether a table with one good
  match resolves. *Status:* `open` — 0002 and 0007 specify different control
  flow for the same call. *(0007×0002 F1, blocks-impl)*

- **JD-3 `RequiresOwned` producer and preservation.** RDR 0007 redefines
  `Row.RequiresOwned` as post-guard write-dependency keys and pins an ordering
  rule over its contents. The string appears **zero** times in RDR 0002, which
  owns the normalized row and enumerates what the dump must preserve — so the
  producer is unnamed and the field may be dropped across the round trip. On
  escape rows the two new RDRs disagree by construction: 0007 derives the field
  from `Writes`, 0009 empties `Writes`. Stake: whether 0007's owned-before-guard
  ordering quantifies over anything at all. *Status:* `open`. *(0007×0002 F3,
  0007×0009 F2, blocks-impl)*

- **JD-4 Static/runtime correspondence.** RDR 0006's lint proves a partition
  exhaustive; RDR 0007's veto can refuse an assignment *inside* that proven set.
  The forward proof holds and the inverse fails, so "lint is green" stops
  meaning "resolution succeeds." No invariant category in 0006 detects
  unevaluability risk, and 0003's coverage algebra has no image for the absent
  case under `contains`, `eq`, `in`, or comparison. Stake: what a green
  `intrastate lint` actually promises. *Status:* `open` — the promise must be
  restated or the lint must gain a category; both are cross-RDR. *(0007×0006 F2,
  0007×0003 F3, 0008×0006 F3, blocks-impl)*

- **JD-5 Precedence of the two whole-table entry preconditions.** RDR 0009 adds
  a breach check at `Resolve` entry and orders it "before every modeled
  disposition"; RDR 0008 adds a reserved-key input precondition at the same
  entry. Neither orders itself against the other — 0009's precedence clause
  ranks itself against *dispositions*, not against a peer precondition, and is
  silent on one. Stake: which error a table with both defects reports, and
  therefore which one the author fixes first. *Status:* `open`. *(0007×0009 F1,
  0008×0009 F2, blocks-impl)*

- **JD-6 Producer-defect surface.** Governed by §D1. Until D1 resolves, the
  kernel seam and the CLI boundary specify contradictory behavior for the same
  event. Stake: whether a programmer defect reaches the user as a stack trace or
  an envelope. *Status:* `open → §D1`. *(0007×0009 F3, whole-set critique C-8/C-9)*

- **JD-7 Read completeness at the accessor→kernel boundary.** RDR 0004 says only
  "A read accessor MUST return typed tag values or a typed refusal" — a
  disjunction with no completeness requirement, so a partially-read artifact may
  conformantly return what it got. `Input.Owned` has no error channel, so a
  truncated snapshot and a genuine absence are the same input. RDR 0007 names
  this as A6b, marks it `Pending`, and routes it to RDR 0004's implement stage;
  RDR 0004 mentions 0007 zero times. Stake: this is the floor under every
  presence/absence contract in the cluster — the critique's "each RDR verified
  its own edge; nobody verified the floor." A truncated read makes `exists =
  false` decide true, and a plan is produced against state that exists.
  *Status:* `open` — the highest-consequence entry here. *(0007×0004 F1/F3,
  0009×0004 F1, blocks-impl)*

- **JD-8 CLI surface for the new refusals.** Three separate gaps at one seam:
  `owned_state_unavailable` has no code row (JD-2's preferred refusal, absent
  from 0005's twelve-row table); `reserved_tag_key` (0008) has no code or exit
  mapping; and 0009 requires `GroupInternal`, which 0005's declared mapping does
  not contain — 0005 enumerates exactly three buckets and assigns all twelve
  codes to two of them. Compounding it, 0005's audit row says "no new envelope
  fields," but 0007 needs `Refusal.Rows`/`MissingOwned` and 0009 needs row
  identity to reach the user. Stake: whether the operator learns *which* rows
  and *what* state, or just that something failed. Scoping half governed by §D2.
  *Status:* `open → §D2` for the scoping; the missing rows are additive and
  settle with the `flow` verb. *(0007×0005 F1/F2/F3, 0009×0005 F1/F2,
  0008×0005 F1, blocks-impl)*

- **JD-9 `--tag` provenance.** Governed by §D3. Also reaches 0008: a caller
  passing `--tag recognized=…` collides with the reserved key, which 0008 classes
  a programmer mistake and 0005 would class user input. Stake: whether a
  guarantee the kernel enforces can be bypassed from the command line.
  *Status:* `open → §D3`. *(0007×0005 F4, 0008×0005 F2, blocks-impl)*

- **JD-10 Recognized-tag totality and fixture ownership.** RDR 0008's
  declaration channel is total (a declared `recognized` tag always exists in the
  model) while the emit channel is conditional (`assemble` injects only for a
  non-empty `Input.Recognized`), so on the empty-outcome input class a row
  matching the reserved tag sees no key rather than a declared-empty one. And
  0008's name constraint invalidates all three of RDR 0002's canonical fixtures
  — which 0002 declares normative: "The RDR and kata spike fixtures are the
  canonical examples implementation tests must promote." 0008's `Overrides`
  claims it "narrows nothing in either peer"; that claim is false for the
  name-constraint half. Stake: whether existing authored tables still load, and
  who owns renaming the fixtures. *Status:* `open`. *(0008×0002 F1/F2/F3,
  0008×0009 F3, blocks-impl)*

### Withdrawn

- **JD-11 Escape-row identity across the dump** — *withdrawn, re-triaged as a
  single-RDR defect.* RDR 0002's round-trip invariant requires the dump to
  preserve "row kind … and escape failure classes"; its dump-derivation contract
  sixty lines away enumerates "model id, row identity, source locator,
  predicates, and writes." Both are `normative` blocks **inside RDR 0002**, so
  the contradiction is visible reading 0002 alone — RDR 0009 only escalates it
  from a fidelity issue to breach-laundering (under 0009's `len(Escape) != 0`
  sole discriminator, a dumped-and-reloaded escape row becomes a conforming
  ordinary row). Not joint. Routed as a SPEC-DEFECT against 0002 at
  **RE-LOCK-ONLY** scope (`re-verify none` — reconcile two field lists). Entry
  retained so the anchor never dangles. *(0009×0002 F3)*

## What this does not decide

Local items stay with their RDRs:

- **RDR 0008's re-entry** — A6/A11 are `Verified` on a false negative-existential
  about Final 0009 ("the string `Input` appears **zero** times"; it appears five
  times, including inside 0009's own Normative Contracts). Already demoted to
  `Draft [revised from Final 2026-08-11; re-verify A6, A11]` with a re-entry note;
  re-enters at Stage 4, STAGE-SCOPED. That is 0008's defect, not a joint one.
- **RDR 0002's dump/round-trip contradiction** — see JD-11 above; a 0002 re-lock.
- **RDR 0004's Prerequisites** — `Final` with every box unchecked, including
  `- [ ] All Critical Assumptions verified`. A finalize-gate question for 0004.
- **RDR 0009's frozen API surface** — `ErrEscapeShapeBreach`, `Unwrap() []error`
  "EXACTLY ONE LEVEL", "return CheckValid's error VERBATIM". Internal to 0009,
  and flagged by the critique as the most likely near-term rewrite.
- **Exact code spellings, exit integers, and envelope field names** — the RDRs
  own them once the decisions above fix the shape.

## Provenance

Cluster gate `0002-0009`, 2026-08-11: whole-set critique (34 rows) plus 13
pairwise scans, all under
`docs/rdr/cluster-reconcile/0002-0009/`. 31 findings, 24 `blocks-impl`. Verdict
NOT RECONCILED — one SPEC-DEFECT (0008, demoted) and these entries, which had no
home until this registry.

**A standing note on why these are hard to resolve on paper.**
`internal/resolve` contains exactly one non-test file and **zero production
consumers** — nothing in `cmd/` or `internal/cli/` imports it. Every entry above
sits at a seam no shipping code has ever exercised. Entries marked `deferred`
are marked so deliberately: manufacturing a paper answer for a question the
first vertical slice would settle is what produced this registry's contents in
the first place.
