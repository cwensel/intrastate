# Recommendation 0011: flow next selects by match; --all enumerates the alphabet

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-26
- **Status**: Draft
- **Type**: Feature
- **Profile**: mid — one user-facing verb predicate (`flow next` candidate set + `--all`) that overrides a locked 0005 contract clause
- **Priority**: High
- **Related Issues**: kata `1mv1` (defect tracker); umbrella rdr#thsc; consumer model rdr#tmxk
- **Predecessors**: 0005-skill-integration-cli-contract
- **Overrides**: `0005:C1`, the `flow next` candidate clause — "flow next MUST return the legal recognized-outcome alphabet for the supplied state, plus candidate summaries …" read as guard-conditioned, and its stated ground "because next enumerates rather than selects" — narrowed to the match-conditioned predicate of C1 below; 0005's enumeration becomes `--all` (C2). Every other `0005:C1` obligation on `flow next` is unchanged. `internal/cli/flow_next_0005_test.go` assertions that relied on a match-excluded row being reported move under `--all`.
- **Seam Lineage**: no prior accretion

## Problem Statement

A skill author (or the human driving one) runs `flow next --model <model> --artifact <state>` to ask "what is legal from here?" so the next action can be chosen without re-deriving the transition table by hand. Today the answer is wrong for that question: RDR 0005 made `next` enumerate rather than select, stripping each row's `match` (and `escape`) before probing, so every row not pruned by a guard is reported as a candidate — 21 of the 22 rules on `models/rdr.toml` at `stage=resolved`. The caller discovers this when the candidate list is the whole alphabet regardless of the supplied state, and cannot constrain a skill's next step from it (0005's own premortem, "next omits enough condition detail to constrain a skill", materialized).

The system-internal requirement is a single decision on the candidate predicate of `flow next`: whether "the legal recognized-outcome alphabet for the supplied state" means guard-conditioned (0005's predicate: candidate = not guard-excluded, match is the kernel's selection pattern and `next` must not select) or match-conditioned (candidate = match holds AND not guard-excluded, evaluated by handing the kernel the row with its match intact so selection semantics stay in the kernel); which of the two is the default and which rides `--all`; and whether guard verdicts on match-excluded rows are still reported. A sub-question inside the same contract: whether `ambiguous_match` / `guard_unevaluable` on a match-retained probe leaves the row a candidate (today only `no_match` excludes). `--all` is the same contract's second surface, not a second contract. Help text follows the decision.

## Critical Assumptions

- **A1 A one-row probe table handed to the kernel with the row's match atoms retained can only yield a plan, `no_match`, `owned_state_unavailable`, or `guard_unevaluable`; `ambiguous_match` cannot arise, so the sub-question "does `ambiguous_match` leave the row a candidate" is vacuous on a probe.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/resolve/resolve.go::Resolve` counts `selected` over the probe's single candidate — `ambiguous_match` needs more than one — and `internal/resolve/resolve.go::gate` returns `owned_state_unavailable` / `guard_unevaluable` before the count; to be re-read against the exact code path once C1's probe shape (absent-key atoms removed) is built.
  - **If wrong**: a match-retained row could be dropped as "ambiguous" and vanish from the candidate list with no unresolved fact explaining it.
- **A2 The CLI already holds the presence test C1 needs: `internal/cli/flow_next.go::summarize` walks `Row.Atoms` — which carries match-block atoms alongside guard atoms (`internal/table/model.go::Row.KernelRow` routes them by `Block`) — against `internal/cli/flow_next.go::assembledView`, so a match key the view lacks is reported in `unresolved` today, and the same presence lookup decides which match atoms the probe retains.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `summarize`'s loop `for _, atom := range row.Atoms { if _, known := view[atom.Key]; known { continue } … c.Unresolved = append(…) }` has no `Block` filter; confirm with a `next` run over a model whose match key is neither owned nor supplied and assert the key appears in `unresolved`.
  - **If wrong**: an absent match key silently drops out of the report and C1's "MUST appear in that candidate's `unresolved`" needs new code rather than the existing loop.
- **A3 The kernel's match seam folds absence into non-match — `internal/resolve/resolve.go::TagSet.matches` returns false when `tv, ok := s.tags[w.Key]; !ok` — so the three-valued match C1 specifies is obtained only by removing absent-key match atoms from the probe before the kernel sees it; removing an atom the view cannot compare changes no value comparison the kernel would have made.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `TagSet.matches` and `Resolve` step 2 ("candidate rows are the non-escape rows for that outcome whose match pattern holds over the assembled view"); no kernel entry point distinguishes absent from mismatched at the match seam.
  - **If wrong**: the kernel already carries an undecided match verdict and C1 should hand it the full match rather than pre-strip.
- **A4 The 0005 `flow next` tests move rather than rewrite: every oracle that asserts on a named candidate seeds the owned key its rows match on (`status=draft` in `flowMVVModel` / `flowGatedNextModel`), and the alphabet, gate-handling, reader-narrowing, and non-mutation oracles do not depend on a match-excluded row being reported.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: `go test ./internal/cli -run 'Next|Req4[0-7]|Req3[5-9]'` after the default flips; any assertion that fails only because a match-excluded row is absent is re-run under `--all` and moves there (C3).
  - **If wrong**: the 0005 suite splits into default / `--all` variants beyond the moves C3 names, and Testing Strategy must enumerate them.
- **A5 Over 0010's decision-table class, `flow next` with no `--tag` reports every ordinary rule with its observed match keys in `unresolved` and empty `required` — so `0010:A11` and 0010's MVV step 5 hold under C1 — and a partial `--tag` set narrows the list to the rows whose supplied keys hold.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: 0010's MVV step 5 (`flow next --model m.toml --as=json` with no `--artifact` → exit 0, every candidate with empty `required`) run against this RDR's build; `0010:C4` leaves `flow next` to this RDR ("`flow next`, `flow read-state`, and `flow set-state` are unchanged by this RDR").
  - **If wrong**: 0010 owes a class-specific clause on this seam, or Alternative 1's strict predicate is what a decision table needs and the two classes diverge.
- **A6 `models/rdr.toml` at `stage=resolved` narrows from 21 reported rules to exactly the rows whose `match.stage` holds over the artifact's owned `stage`, each ordinary row's `recognized` match holding trivially because the probe binds `recognized` to the row's own outcome.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: the MVV below, step 3; `internal/cli/flow_next.go::excluded` already sets `Recognized: row.Outcome` on the probe.
  - **If wrong**: the observed defect is not closed by a match predicate and the problem statement's diagnosis (match stripped) is incomplete.

## Proposed Solution

### Approach

`flow next` keeps asking the kernel one question per row, but stops taking the row's match pattern out of the question. The probe `internal/cli/flow_next.go::excluded` builds is reshaped: the row's match atoms over keys the assembled view **carries** stay on the probe so the kernel compares them; match atoms over keys the view **lacks** are removed and the key rides the candidate's `unresolved` list exactly as an unknown guard key does today; the escape list stays stripped. The kernel's answer is read as before — only `no_match` excludes; `guard_unevaluable` and `owned_state_unavailable` are undecided, not false, and leave the row a candidate. The verb still never chooses among the rows that survive, never turns a gate deny into a refusal, and still returns the model's full declared `outcomes` alphabet: it narrows the **candidate list** to rows the supplied state can actually take, which is what "for the supplied state" meant to the caller.

`--all` restores 0005's enumeration verbatim — every match atom removed from the probe — for the caller who wants the guard-conditioned alphabet rather than the state-conditioned one. The two modes differ only in what the probe carries; reader narrowing, gate opt-in, the payload shape, and exit semantics are the same in both. ⇒ the whole change is one probe-shape decision inside one function, plus a flag and its help text; no kernel change, no new payload field, no new refusal code.

Sibling-path check (step 5): the presence discriminator this adds — "is this key in the assembled view?" — is the one `internal/cli/flow_next.go::summarize` already applies to every atom of `Row.Atoms` to fill `unresolved`, over the same `internal/cli/flow_next.go::assembledView` map (its writer: `runFlowNext`, from `runReaders`' owned tags plus `--tag` observed tags). C1 reuses that signal; searched `internal/cli` and `internal/resolve` for another "absent key ⇒ undecided" site at the match seam — none exists (the kernel's `TagSet.matches` folds absence into false).

### Technical Design

Data flow, per ordinary row: `table.Row` → `Row.KernelRow()` (match/guard split by block) → probe reshaping in `excluded` (default: drop match atoms whose key is absent from the view; `--all`: drop every match atom; both: drop escape) → `resolve.Resolve` over the one-row table with `Recognized = row.Outcome` → `no_match` ⇒ excluded, anything else ⇒ candidate → `summarize` (unchanged: absent atom keys, missing owned keys, and un-run gate ids become `unresolved`) → gates under `--evaluate-gates` for reported candidates only → payload. The alphabet `outcomes` is copied from the model in both modes (0005's DEV-4 reading stands: the alphabet is what may be requested; the candidates are what varies with state).

Selection semantics stay in the kernel: the CLI decides presence (a lookup it already owns for reporting) and never compares a value; equality, set-literal spelling, the conflicted-key rule (`0007:C8` via `TagSet.matches`), and the guard verdict are all the kernel's over the retained atoms.

#### Normative Contracts

**C1**

```normative
flow next MUST report as a candidate exactly each non-escape row of the
requested model for which (a) every match atom whose key is PRESENT in the
assembled view — owned tags from the invoked readers, observed --tag values,
and `recognized` bound to the row's own outcome — holds under the kernel's
own match comparison, and (b) the kernel does not decide the row's guard
false. A match atom over a key the assembled view does not carry is an
unresolved fact: it MUST NOT exclude the row, and its key MUST appear in
that candidate's `unresolved`. Both verdicts MUST be the kernel's, asked
over a one-row probe table carrying the row's match atoms over present
keys, none of its match atoms over absent keys, and no escape list; only
the kernel refusal `no_match` excludes. `guard_unevaluable` and
`owned_state_unavailable` MUST leave the row a candidate with the
undecided facts reported in `unresolved`. A row excluded under (a) is not a
candidate: it MUST NOT be reported, and under --evaluate-gates its gates
MUST NOT run. `outcomes` MUST remain the model's full declared alphabet in
every mode. flow next MUST NOT choose among the candidates it reports and
MUST NOT turn a gate deny into a refusal.
```

⇒ this is the clause that overrides `0005:C1`'s candidate reading; the exclusion mechanism (kernel probe, `no_match`-only) is 0005's, extended from guard-only to match-and-guard.

**C2**

```normative
flow next MUST accept a boolean flag --all, default false. Under --all the
candidate predicate MUST be exactly the one 0005:C1 specified: every
non-escape row whose guard the kernel does not decide false, with EVERY
match atom removed from the probe; a match atom over a present key is then
neither an exclusion nor an unresolved fact, and a match atom over an
absent key is reported in `unresolved` as today. Everything else 0005:C1
requires of flow next — the data minimum, gate accessors only under
--evaluate-gates and only for reported candidates, a deny reported on its
candidate with exit 0, no invented guard facts, the narrowed invoked reader
set — MUST hold identically in both modes. The success payload shape MUST
be the same in both modes; --all is part of request identity. --all MUST
NOT be accepted by flow resolve, flow read-state, or flow set-state.
```

⇒ 0005's behaviour is preserved one flag away, so a caller that wanted the enumeration loses nothing and the 0005 oracles for it move rather than die.

**C3**

```normative
The command's Short and Long help MUST state that flow next reports the
candidates the supplied state can take (match and guard both undecided or
holding), that an unsupplied match or guard key leaves a row a candidate
with the key listed under unresolved, and that --all reports every row the
guards do not exclude regardless of match. The 0005 flow next tests in
internal/cli/flow_next_0005_test.go MUST keep passing unchanged where they
assert the alphabet, gate handling, reader narrowing, determinism, and
non-mutation; an assertion that depended on a match-excluded row being
reported MUST be re-homed under --all, not deleted, and the file's header
comment MUST name this RDR as the source of the default.
```

⇒ the override is visible where a reader meets it — help text and the test file — not only in this record.

#### Load-Bearing Decisions

- **Identity** — a `flow next` request is identified as `0005:D-identity` says, plus the `--all` bit; the same request in the same mode over the same model revision and artifact contents reports the same candidate set.
- **Naming** — the flag is `--all`. Rejected: `--enumerate` (describes 0005's verb, not the caller's intent), `--ignore-match` / `--no-match` (names the mechanism and collides with the kernel refusal spelling), and flipping the sense (`--select` / `--match` on a 0005 default) because the defect is the default and every skill would carry the flag (Alternative 2).
- **Selection / predicate** — default: candidate ⇔ no match atom over a present key fails ∧ guard not decided false (C1). `--all`: candidate ⇔ guard not decided false (C2). In both modes `next` reports all survivors and chooses none; exactly-one selection and refusal remain `flow resolve`'s (`0005:D-selection-predicate`).

#### Illustrative Code

Illustrative — intent only.

```sh
# Default: rows the supplied state can take. At stage=resolved this lists the
# resolved→* rows, not all 22.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --as=json

# 0005's enumeration: every row the guards do not exclude, match ignored.
intrastate flow next --model models/rdr.toml --artifact rdr=./0011.md --all --as=json

# A decision table with no tags: every rule is a candidate; each observed
# match key is listed under `unresolved`. Supplying --tag a=x narrows it.
intrastate flow next --model table.toml --as=json
```

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Per-row kernel probe | `internal/cli/flow_next.go::excluded` | strips every match atom (`probe.Match = nil`) | Extend | probe keeps present-key match atoms by default; strips all under `--all` |
| Presence test for "absent = unresolved" | `internal/cli/flow_next.go::summarize` over `assembledView` | none — already walks match-block atoms | Reuse | C1's `unresolved` clause is the existing loop |
| Match/guard split on the probe row | `internal/table/model.go::Row.KernelRow` | none | Reuse | no table change |
| Match comparison, guard verdict, refusal kinds | `internal/resolve/resolve.go::Resolve`, `::gate`, `::TagSet.matches` | `matches` folds absence into false | Reuse | no kernel change; the CLI pre-strips absent-key atoms (A3) |
| Flag registration | `internal/cli/flow.go::registerSelectionFlags` + per-verb `cmd.Flags()` | none | Reuse | `--all` registered on `next` only (C2) |
| Reader narrowing | `internal/cli/flow_exec.go::invokedReaders` | union over all rows for `next` | Reuse | unchanged in both modes |

### Decision Rationale

The caller's question is "what can I do from here", and every peer that answers it conditions the answer on the current state: stateless's `PermittedTriggers` returns only triggers whose guards are met in the current state; xstate v5 removed the unconditioned `nextEvents` and kept the conditioned `can()`; pytransitions keeps the declared form (`get_triggers`) but as a separate operator from the evaluated one (`may_*`); SCXML's `cond` gates enablement (Investigation). 0005's predicate was the declared-shape answer; in this table model the analogue of "current state" is the match block, so a state-conditioned `next` is a match-conditioned one. The mechanism is chosen to move the least: the kernel already answers a one-row probe, and the CLI already owns the presence lookup that makes an absent key undecided rather than false — the same posture `excluded`'s own doc states for guards ("treating 'unknown' as 'false' would silently drop rows the caller is entitled to see"). Alternative 1 (strict kernel match) was rejected because it empties the list for any caller that supplies partial state, including 0010's no-tag decision table (A5), and contradicts that posture at the same function. Alternative 2 (keep 0005's default, opt-in flag) was rejected because the defect is the default and the instance read says the conditioned form is the primary one. Alternative 3 (a kernel-side three-valued match entry point) was rejected because it moves a reporting concern into the kernel for no semantic gain — the CLI still compares no value.

Premortem (paragraph): this shipped, and a skill driving a model whose rows match on an owned key its reader does not serve (a typo'd `keys` list) saw every row listed as a candidate with the key under `unresolved` — the same wall-of-candidates symptom as before, now with a hint. A second failure: a caller who scripted against 0005's list found candidates missing and no flag in the error, because there is no error — the list is just shorter. Both are answered by the design rather than forcing a switch: the first is the intended three-valued behaviour and the `unresolved` key is the diagnostic (the alternative, silently dropping the row, is the worse failure); the second is the override's cost, mitigated by C3's help text, the test-file header, and `--all`. Neither shows a case the chosen predicate cannot answer.
Premortem: survived (paragraph)
Ground-sweep: clean (24 anchors) — 23 confirmed, 1 cosmetic correction folded (`models/rdr.toml` has 22 `[[rule]]` tables; the observed 21 reported is the seed's count, now stated as 21 of 22). Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/grounding-sweep/sweep.md`.
Joint-check: fired → 0010 (home: cli/0011:C1; disposition 2026-08-26: cite-don't-restate — C1 is the sole normative home, 0010 cites it and restates nothing). Open peers at depth 1 with Status Draft/Final: cli/0010 only (0001–0009 are `Implemented`). Shared modify-anchor: `internal/cli/flow_exec.go::invokedReaders` (Reuse/unchanged in both RDRs). Shared contract literals: `no_match` (0010: the escape class of its "otherwise" row, `0010:A9`; here: the excluding kernel refusal in C1) and `flow next`. The coupling behind the tokens is real: `0010:A11` and 0010's MVV step 5 assert that `flow next` over a decision table with no `--tag` lists every rule with empty `required`, which holds only under C1's absent-key rule (A5) and fails under Alternative 1; 0010 defers that predicate to this RDR (`0010:C4`), so the decision is homed in C1 by citation. Dispositions for the user: cite-don't-restate (`0010:A11` keeps citing cli/0011:C1; nothing added to 0010) or declare in both (0010 gains a class clause on `flow next` over the decision-table class). Absence arm: this RDR narrows a reported set rather than converting a refusal into an acceptance; no `Final` peer exists, and the `Implemented` predecessor 0005's reliance on the enumeration is named under `Overrides` and rides to 7.1. Bridge sub-check: n/a — neither plan retires a surface the other introduces (0010 stays out of `flow_next.go`; this RDR stays out of `resolvePayload`, `normalizeRule`, `reach`, `checkGroups`).

## Alternatives Considered

### Alternative 1: Strict kernel match (absent key excludes)

**Description**: Hand the kernel the row's full match block and let `TagSet.matches` decide: a match key the view lacks is a non-match, so the row is excluded.

**Pros**:

- Zero CLI logic beyond removing two `nil` assignments; `next`'s candidates equal exactly the rows `resolve` could select.

**Cons**:

- Any partially supplied state — a decision table queried with no `--tag`, a machine whose match key no reader serves — yields an empty list with no unresolved fact to act on.
- Contradicts the "undecided is not false" posture `excluded` already takes for guards and 0007 takes for atoms.

**Reason for rejection**: it breaks 0010's no-tag `flow next` (A5) and turns the verb's most useful degraded answer into silence.

### Alternative 2: Keep 0005's default; add an opt-in `--match` / `--select` flag

**Description**: Leave the guard-conditioned enumeration as the default and let a caller request the match-conditioned list with a flag.

**Pros**:

- No override of `0005:C1`; the 0005 tests do not move.

**Cons**:

- The defect is the default: every skill invocation must carry the flag or get the alphabet.
- The instance read points the other way — xstate retired the unconditioned form, stateless never offered it.

**Reason for rejection**: fixes the symptom for callers who know about the flag and leaves the observed defect in place for everyone else.

### Alternative 3: Kernel-side three-valued match (`resolve.Candidates`)

**Description**: Add a kernel entry point that returns, per row, match holds / fails / undecided, so the CLI never touches the probe shape.

**Pros**:

- One home for match semantics, including absence.

**Cons**:

- A new kernel API and a JDR §D2-adjacent decision for a CLI reporting need; the presence lookup already exists in the CLI and compares nothing.

**Reason for rejection**: larger blast radius (0001/0007 seam) for no change in what the caller sees.

### Briefly Rejected

- **Filter `outcomes` to outcomes with a surviving candidate**: 0005 already refused this reading (DEV-4 in `runFlowNext`); the alphabet is what may be requested.
- **A separate verb (`flow candidates`)**: two verbs for one question; 0005's four-verb closure stands.
- **Report match-excluded rows with an `excluded_by` field instead of dropping them**: grows the payload for the same information `--all` already yields.

## Context

### Background

Observed 2026-08-26 while driving the consumer model (rdr#tmxk): `flow next --model models/rdr.toml --artifact rdr=<state with stage=resolved>` lists all 21 rules. The cause is by design in RDR 0005: `internal/cli/flow_next.go::excluded` nils `probe.Match` / `probe.Escape` before probing, under the comment "dropping a row whose match pattern the supplied facts do not satisfy would be a selection this verb was not asked to make" — the sentence this RDR revisits. 0005 rationale `0005:A3` holds that `next` exposes the alphabet without owning guard evaluation; selection (gate-then-count, exact-one survivor) belongs to the kernel. 0005's Briefly Rejected list also refused exit-code-only output because it "cannot carry legal outcome alphabets, conditional summaries" — the conditional summary is the surface this RDR sharpens. Constraints: RDRs are never amended in content (this is a new RDR that overrides 0005's clause); intrastate stays generic — no consumer (RDR-process) knowledge in code, docs, or fixtures; `make check` must pass. Not a facet of kata `zdat` / cli/0010 (owned-state optionality is a model-class decision; this is a verb-predicate decision) — cross-cite only; 0010 leaves `flow next` to this RDR (`0010:C4`, `0010:A11`).

### Technical Environment

Go CLI (`bin/intrastate`); `internal/cli/flow_next.go` (`excluded`, `summarize`, `assembledView`, candidate/guard reporting), kernel probe API `internal/resolve/resolve.go::Resolve` that today receives rows with match stripped, `internal/cli/flow_next_0005_test.go`. Conventions in `AGENTS.md` (respond gateway, CLIError codes, SilenceUsage). Design history: `docs/rdr/0001–0010`, `docs/jdr/0001`.

## Research Findings

### Investigation

Prior art was read first, class and instance. Class (StateMachineRes corpus, three queries, ledger in `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`): the "what next" operator exists in every peer that has a current state, and is conditioned on it. Instance: qmuntal-stateless's README "Introspection" — "a list of the triggers that can be successfully fired within the current state via the `StateMachine.PermittedTriggers` property" — implemented in `states.go::stateRepresentation.PermittedTriggers` as triggers with `len(tb.UnmetGuardConditions(...)) == 0` ⇒ state- and guard-conditioned, the conditioned form is the only form. pytransitions `transitions/core.py::Machine.get_triggers` returns triggers declared from the given states without evaluating conditions, while `Machine._can_trigger` (the `may_*` family) evaluates them ⇒ the declared-shape enumeration exists but as a distinct operator, which is what `--all` is here. xstate's `packages/core/CHANGELOG.md` records "Removed `MachineSnapshot['nextEvents']`" and `packages/core/src/State.ts::machineSnapshotCan` keeps the guard-evaluated `can(event)` ⇒ when forced to keep one, the peer kept the conditioned query. scxmlcc `doc/user-manual.md` (`cond`): "The transition is only executed if the condition evaluates to true" ⇒ enablement is event- and condition-conditioned. ⚠ no prior-art coverage for a peer that distinguishes an *absent* state key from a mismatching one at the match seam — every peer assumes a current state is always present — so the three-valued absent-key rule in C1 rests on this codebase's own posture (`excluded`, `0007:C8`) and is verified by A2/A3/A5, not by prior art. Code paths: `internal/cli/flow_next.go::excluded` (the probe; `probe.Match = nil`), `::summarize` (presence walk over `Row.Atoms`), `internal/table/model.go::Row.KernelRow` (block routing), `internal/resolve/resolve.go::Resolve` / `::gate` / `::TagSet.matches` (the kernel's two verdicts and its absence rule).

### Key Discoveries

- **Documented** — `internal/cli/flow_next.go::excluded` strips `Match` and `Escape` from the probe and excludes only on `resolve.KindNoMatch`; its doc gives the "unknown is not false" rule this RDR extends to match. ⇒ the change is a probe-shape change in one function.
- **Documented** — `internal/resolve/resolve.go::TagSet.matches` returns false for an absent key (`!ok`) and for a conflicted key. ⇒ absent-key atoms must be removed from the probe to stay undecided (A3); conflicted keys stay the kernel's non-match.
- **Documented** — `internal/cli/flow_next.go::summarize` walks every `Row.Atoms` entry, match-block atoms included, against the view. ⇒ an absent match key is already an `unresolved` entry (A2).
- **Documented** — 0005 Decision Rationale premortem: "`next` omits enough condition detail to constrain a skill". ⇒ 0005 foresaw this failure and answered it with candidate *summaries*; the summaries were not enough because the *set* was unconditioned.
- **Documented** — `0010:C4`: "`flow next`, `flow read-state`, and `flow set-state` are unchanged by this RDR"; `0010:A11` defers the decision-table `next` behaviour to this RDR's predicate. ⇒ the absent-key rule must serve a no-tag decision table (A5).
- **Documented** — prior art (Investigation): stateless conditioned-only; pytransitions two operators; xstate retired the unconditioned one. ⇒ conditioned is the default, enumeration is the secondary surface.
- **Assumed** — every 0005 `next` oracle seeds the owned key its named rows match on (A4).

## Trade-offs

### Consequences

- Positive: a skill can read the next step from the candidate list at `stage=resolved` (A6) instead of the whole table.
- Positive: no kernel, table, payload, or refusal-code change; the diff is one function, one flag, help text, and test moves.
- Negative: an override of a locked 0005 clause — callers scripted against the enumeration must add `--all`.
- Negative: a candidate can still be listed that `resolve` would refuse with `no_match` (absent match key) — the same relationship `next` already has with `guard_unevaluable`.

### Risks and Mitigations

- **Risk**: an absent match key makes the list look like the old wall of candidates.
  **Mitigation**: the key is named under `unresolved` on every such candidate (C1); the fix is on the caller's side (bind the reader / supply the tag) and is visible.
- **Risk**: a 0005 oracle silently passes for the wrong reason after the default flips.
  **Mitigation**: A4's run plus C3's explicit re-homing rule; Resolve's Testing Strategy adds a default-vs-`--all` pair on the MVV fixture.

### Failure Modes

- **Visible**: a candidate the caller expected is absent. Diagnose with `--all`: if it appears there, a match atom over a supplied key failed (compare the row's `match` block in `dump` against `owned`/`observed` in the payload).
- **Visible**: every row is a candidate with the same key under `unresolved` — the match key is neither owned by an invoked reader nor supplied; bind or supply it.
- **Silent**: a match atom over a present key with a conflicted value (`0007:C8`) is a kernel non-match, so the row is excluded rather than reported undecided. This is the kernel's rule and is the same under `resolve`; `--all` shows the row.
- **Recovery**: none needed — `next` is effect-free in both modes; re-run with the corrected inputs.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A2/A3 first — they fix the probe shape; A4 decides the test moves)
- [ ] Joint-check disposition with cli/0010 recorded (Decision Rationale)

### Minimum Viable Validation

1. Build with the C1 default and the `--all` flag.
2. Seed an artifact for `models/rdr.toml` whose owned `stage` is `resolved`; run `flow next --model models/rdr.toml --artifact rdr=<it> --as=json`.
3. Assert: `candidates[]` contains exactly the rules whose `[rule.match.stage]` admits `resolved` (and no rule whose `match.stage` names another value); `outcomes[]` is still the model's full alphabet.
4. Re-run with `--all`; assert `candidates[]` has the set 0005 reported (21 of the 22 rules) and the same `outcomes[]`.
5. Run over the 0005 MVV fixture (`flowMVVModel`, `status=draft`) with and without `--all`; assert the `gated-excluded` row is absent in both (guard-excluded), and that a row whose `match.status` names a value other than `draft` is absent by default and present under `--all`.
6. Run over a model with no reader and no `--tag` whose rows match on an observed key (0010's decision-table shape or a generic fixture); assert every ordinary row is a candidate with that key under `unresolved`; add `--tag key=<value>` and assert the list narrows to the rows whose match holds.
7. `make check` passes; the 0005 `next` suite passes with the moves C3 names and no deletions.

End-state: `flow next` answers "what is legal from here" with the rows the supplied state can take, `--all` yields 0005's enumeration, and every 0005 obligation other than the candidate predicate is unchanged.

### Phase 1: Probe Shape

Intent: make `excluded`'s probe carry the row's match atoms over present keys (absent-key atoms removed, escape removed) so the kernel decides match and guard together; `summarize` is unchanged.

### Phase 2: `--all` and Help

Intent: register `--all` on `next` only, route it to the all-match-stripped probe, and rewrite Short/Long help per C3.

### Phase 3: Tests and Docs

Intent: re-home the 0005 `next` assertions that relied on match-excluded rows under `--all`, add the default/`--all` pair and the absent-key case, update the file header, and add the `--all` invocation to `docs/cli-output-contract.md`.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

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
> At lock, replace this section's body with the
> one-line pointer to gate.md — responses are never
> inlined. The sub-sections below spec gate.md's
> content.

### Contradiction Check

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

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

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

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

- RDR 0005 `0005:C1` (the overridden clause), `0005:A3`, `0005:A7`, `0005:D-selection-predicate`, Decision Rationale premortem; RDR 0010 `0010:C4`, `0010:A11`; RDR 0007 `0007:C8`; JDR 0001 §D2.
- Source reviewed: `internal/cli/flow_next.go`, `internal/cli/flow_exec.go`, `internal/cli/flow.go`, `internal/table/model.go`, `internal/resolve/resolve.go`, `internal/cli/flow_next_0005_test.go`, `internal/cli/flow_fixtures_0005_test.go`, `docs/cli-output-contract.md`.
- Prior art (sibling checkouts under `../state-machines`): `repos/qmuntal-stateless` (`README.md` Introspection; `statemachine.go`, `states.go`), `study/transitions` (`transitions/core.py`), `repos/xstate` (`packages/core/CHANGELOG.md`, `packages/core/src/State.ts`), `repos/scxmlcc` (`doc/user-manual.md`). Ledger: `docs/rdr/0011-flow-next-match-conditioned-candidates/evidence/research/prior-art.md`.
- kata `1mv1`; rdr#thsc; rdr#tmxk.
