---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: open
cluster: 0002, 0003, 0004, 0005, 0006, 0007, 0008, 0009
labels: intrastate, internal-resolve, refusals, guard-evaluation, rdr-cluster
---

# JDR 0001 What the Resolver Can Trust, and What the User Sees When It Can't

*Go identifiers and CLI transcripts below are **non-normative**; the RDRs own
exact contracts. Source: the `0002-0009` cluster gate —
`docs/rdr/cluster-reconcile/0002-0009/`.*

## Problem statement

Eight RDRs jointly decide one thing: **when the resolver cannot decide, what
happens.** 0002 owns the table format, 0003 the guard vocabulary, 0004 accessor
execution, 0005 the CLI surface, 0006 graph lint. Three later RDRs — 0007
(guard evaluation domain), 0008 (recognized tag key), 0009 (escape-row shape) —
each landed on the `area:internal-resolve` seam between them, two months after
the first five locked.

The structural cause of every finding is one grep:

```
grep -cE 'RDR 0007|RDR 0008|RDR 0009' 0002…md 0003…md 0004…md 0005…md 0006…md
→ 0, 0, 0, 0, 0
```

No older member references any newer one. Meanwhile 0007 names
`Obligation destination (NAMED): RDR 0003's implement stage` three times and
routes A6b to 0004's implement stage. An obligation filed on a document that
never received it binds nobody.

## Next steps

**Three decisions below need a call. Nothing else here blocks.**

1. **Decide D1, D2, D3** — the only genuine forks. Each has options and a
   recommendation.
2. **`/rdr-resolve 0008`** — demoted at the gate for a real defect of its own
   (A6/A11 `Verified` on a false claim about Final 0009). Independent of this
   registry.
3. **Re-lock 0002** — its dump contract and round-trip invariant contradict each
   other (withdrawn JD-11 below). RE-LOCK-ONLY.
4. Then implement. The interface record's remaining entries are answers or
   blanks, not gates.

Nothing in `internal/resolve` has a production consumer today — one non-test
file, no `flow` verb. Entries marked *(blank)* are settled by writing that code,
not by more specification.

## Principles

Steering criteria for every decision here. Each is derived from text the members
already lock.

1. **Missing artifact state is never masked** — 0007's reason for existing: the
   domain choice "decides whether missing artifact state can be masked behind an
   escapable refusal class."
2. **The kernel refuses rather than guesses** — 0004: "Indeterminate MUST be a
   refusal-class result, not a false allow and not a false deny."
3. **Every user-visible failure is structured** — 0005 routes all failure
   through `respond.Fail(cmd, *clierr.CLIError)`; the consumer parses JSON.
4. **A refusal names what was missing and who can fix it** — 0007's user
   outcome, plus 0005's exit-code grouping, which encodes remedy.
5. **A green lint means resolution succeeds** — otherwise 0006's blocking
   authority promises nothing.
6. **One home per contract; cite, never restate** — the engine doctrine, and
   why this registry exists.
7. **Pre-release, prefer the clean shape** — no shims for callers that do not
   exist. *(Rests on a project convention, not repo text — confirm before
   leaning on it.)*

---

## D1 — How does the guard seam carry a predicate?

`resolve.Row.Guard` is a `string` (`resolve.go:185`) — an opaque predicate. The
evaluator must reconstruct structure from it and can fail doing so, and 0007
answers that failure with a mandated panic the kernel "MUST NOT recover." 0005
requires every user-visible failure to be a structured envelope. Neither RDR
mentions the other.

The same opacity causes a second defect: 0003's predicate identity is
*per-atom*, so an ordinary two-condition row needs N keys through a one-slot
channel.

- **(a) Panic stands.** A stack trace reaches a JSON consumer. Violates P3.
- **(b) Error channel on `Evaluate`.** 0007 rejected it: a second channel lets
  the evaluator report undecidability outside the verdict taxonomy. The
  objection conflates a *verdict* (cannot decide — modeled as
  `GuardUnevaluable`) with a *defect* (structure is broken), but under (b) the
  distinction is held by discipline rather than structure.
- **(c) Panic at the seam, recover at the CLI.** `panic`/`recover` must cross
  `internal/resolve` → `internal/cli`; recovering only the mapping panic needs a
  typed value and a re-panic default, and every non-CLI consumer inherits the
  obligation. Was the recommendation while "do not reopen Final RDRs" held.
- **(d) Carry the parsed predicate — recommended.** Row carries a slice of atoms
  (key, operator token, literal, block) instead of a string. No reconstruction
  step, so no mapping failure, so no channel question. The kernel carries
  *cardinality*, not *meaning* — it stays as grammar-agnostic as it already is
  for `Match []Tag` and `Writes []Tag`. **Dissolves JD-1.**

  *Cheap:* the only `GuardEvaluator` implementation in the repo is
  `fixtureGuards` in `fixtures_test.go:20`; `resolve.go` is the sole non-test
  file referencing the interface. *Not free:* reopens 0001's `Row` type, which
  0007 declined to touch (A3), at FULL-FLOW scale for 0007. `Implemented` has no
  backward edge, so this lands either as 0007 absorbing the change (precedent:
  it already redefines `Row.RequiresOwned`) or a new RDR superseding 0001's row
  shape.

*The current fixture is the argument:* given a guard it cannot map it returns
`GuardUnevaluable` — the exact conflation 0007's panic clause exists to forbid.

**Resolved:** _pending._

## D2 — Does the aggregation veto run before or after match counting?

0002 states escape reachability as a function of the match count over non-escape
rows. 0007 refuses *before* counting: "if any surviving candidate row's guard is
GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`." An
implementer reading 0002 builds two stages; reading 0007, three.

- **(a) Count first (0002's shape).** An unevaluable row is pruned, the
  resolution reports `no_match` — which 0002 makes escapable. Violates P1: the
  missing state is masked behind an escapable class. This is precisely the
  outcome 0007 was written to prevent.
- **(b) Gate first (0007's shape) — recommended.** Evaluate guards, veto on
  unevaluable, then count survivors. `guard_unevaluable` is not a modelable
  escape class, so the refusal reaches the user intact. Cost: 0002's Technical
  Design prose and its Scenario 4 expectation need restating at its re-lock —
  0002 states the counting rule in its own voice, so 0007's deferral to 0001
  does not reach it.

**Resolved:** _pending._

## D3 — Must a read accessor return a *complete* tag set?

0004 says only: "A read accessor MUST return typed tag values or a typed
refusal." A disjunction with no completeness requirement — a partially-read
artifact may conformantly return what it got. `Input.Owned` has no error
channel, so a truncated snapshot and a genuine absence are the same input.

0007 names this A6b, marks it `Pending`, and routes it to 0004's implement
stage. 0004 mentions 0007 zero times.

This is the floor under every presence/absence contract in the cluster: a
truncated read makes `exists = false` decide TRUE, and a plan is produced
against state that exists. It is also the easiest entry to wave through, because
0004's clause reads fine until you notice what it omits.

- **(a) Status quo.** Partial reads stay conformant and indistinguishable.
  Violates P2 — the kernel guesses.
- **(b) Require completeness — recommended.** A read accessor MUST return the
  complete tag set for the keys it was asked for, or take the refusal branch.
  States the rule 0007 needs at the layer that owns it. Cost: one normative
  clause in 0004 plus an MVV scenario asserting a partial read refuses.
- **(c) Error channel on `Input.Owned`.** Lets the kernel distinguish
  read-failed from absent. More expressive, more surface, and 0004 still has to
  say when to use it — so (b) is a prerequisite either way.

**Resolved:** _pending._

---

## Interface record

Cite `JDR 0001 §JD-n`; never restate. Entries marked *(blank)* name an owner,
not a negotiation.

- **JD-1 Guard atom transport.** Decided by §D1(d): carrying parsed atoms
  removes the N-atoms-through-one-slot problem. Reverts to an open question at
  0003's Phase 1 if D1 resolves otherwise. *(0007×0003 F1)*
- **JD-2 Resolver control flow.** See §D2. *(0007×0002 F1)*
- **JD-3 `RequiresOwned` producer.** The field appears **zero** times in 0002,
  which owns the normalized row and enumerates what the dump preserves — so no
  layer is obliged to populate it, and 0007's owned-before-guard ordering may
  quantify over an always-empty set. On escape rows 0007 derives it from
  `Writes` while 0009 empties `Writes`. Not answered by D1, but D1 reopens the
  same type and is the natural occasion to settle it. *(0007×0002 F3,
  0007×0009 F2)*
- **JD-4 Lint's promise is what gives.** Where 0006's exhaustiveness proof and
  0007's veto disagree, the *promise* narrows — P5 decides the substance. Open
  only as to which document records the narrowing and whether lint gains a
  warning category. *(0007×0006 F2, 0007×0003 F3)*
- **JD-5 Precondition precedence.** 0009's breach check and 0008's reserved-key
  check both land at `Resolve` entry; neither orders itself against the other.
  Either order is defensible — pick one and pin it with a test on a table that
  breaches both. The silence is the defect, not the choice. *(0007×0009 F1,
  0008×0009 F2)*
- **JD-6 Producer-defect surface.** See §D1. *(0007×0009 F3)*
- **JD-7 Read completeness.** See §D3. *(0007×0004 F1/F3, 0009×0004 F1)*
- **JD-8 Refusal codes.** `owned_state_unavailable` and `reserved_tag_key` need
  `Code` values; 0009 needs `GroupInternal`, which **already ships** in
  `clierr.go`. 0005's own A-block pre-authorizes this: resolver-specific values
  "only require new `Code` constants or literals, **not new envelope fields or
  exit groups**." `Detail` and `Hint` ship and render in text and JSON. One live
  question: `GroupUserEnv` and `GroupInternal` both exit 2, so the only remedy
  distinction an exit code can carry is exit 3 — and an accessor that could not
  read the artifact is what `GroupEnvUnavailable` already means. *(blank, except
  the exit-3 call)* *(0007×0005 F1/F2/F3, 0009×0005 F1/F2, 0008×0005 F1)*
- **JD-9 `--tag` provenance.** Caller-supplied tags enter as
  `ProvenanceObserved` and never satisfy an owned-state dependency. `assemble`
  already pins owned-over-observed precedence so the accessor snapshot is "never
  shadowed by caller-supplied context"; the exposure opens only if the `flow`
  verb wires `--tag` into `Input.Owned`. Wire it to `Observed`.
  *(0007×0005 F4, 0008×0005 F2)*
- **JD-10 Recognized-tag totality.** 0008's name constraint invalidates all
  three of 0002's canonical fixtures, which 0002 declares normative — rename
  them; there are no users to migrate. Still open: whether a declared
  `recognized` tag is total (always present, possibly empty) or partial (absent
  with no outcome in flight), which decides whether an empty-outcome row is
  satisfiable, dead, or a lint error. No normalizer exists yet to observe it.
  *(0008×0002 F1/F2/F3, 0008×0009 F3)*

**Withdrawn — JD-11 Escape-row identity across the dump.** Re-triaged as a
single-RDR defect: 0002's round-trip invariant requires the dump to preserve
"row kind … and escape failure classes" while its dump-derivation contract
sixty lines away omits both. Both are `normative` blocks inside 0002, so the
contradiction is visible reading 0002 alone; 0009 only escalates it to
breach-laundering. Routed as a SPEC-DEFECT against 0002, RE-LOCK-ONLY. Entry
kept so the anchor never dangles. *(0009×0002 F3)*

## What this does not decide

- **RDR 0008's re-entry** — A6/A11 verified on a false claim about Final 0009;
  already demoted, re-enters at Stage 4, STAGE-SCOPED.
- **RDR 0002's dump contradiction** — JD-11 above; a 0002 re-lock.
- **RDR 0004's Prerequisites** — `Final` with every box unchecked, including
  "All Critical Assumptions verified." A finalize-gate question for 0004.
- **RDR 0009's frozen API surface** — `Unwrap() []error` "EXACTLY ONE LEVEL",
  "return CheckValid's error VERBATIM". Internal to 0009.
- **Exact code strings, exit integers, envelope field names** — the RDRs own
  them once the decisions above fix the shape.
