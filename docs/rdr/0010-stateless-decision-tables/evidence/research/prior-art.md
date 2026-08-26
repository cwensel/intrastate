Model: claude-fable-5

# cli/0010 — Stage 2 prior-art pass (bounded read, not a spike)

Budget: ≤3 corpus queries per claim, ≤5 opened hits per claim. Queries ran
2026-08-26 via `arc search semantic --corpus <C> --limit 5 --json "<q>"`;
source-level reads via grep over the `../state-machines` checkouts the hits
named (semble not needed — the arc hit gave the file).

## Claims the choice rests on

### K1 — class: a decision table is a stateless artifact distinct from a state machine

- Query (DevRef): `decision table hit policy exactly one matching rule output entries`
- Accepted: `developer-testing.pdf` p.142 — "we can use _decision tables_,
  which capture all combinations of variables and possible outcomes." ⇒ the
  decision-table artifact is *defined by* its combination coverage — the
  property the PoC's 208/512 hole exercised — and carries no state variable.
- Rejected: SQL `MATCH` predicate, SARIF policies, Postgres RLS / MERGE
  `WHEN` (score < 0.47; keyword noise on "matching rows").
- Query (StateMachineRes): `decision table versus state machine: a stateless
  lookup from condition values to an action with no state variable` — top hits
  are state-machine intros (javascript-state-machine, statewright, scxmlcc);
  none contrasts against decision tables. ⚠ partial coverage: the corpora
  hold no DMN text, so "unique hit policy = exactly one row" is model prior,
  demoted to a Resolve assumption (A5) rather than leaned on.

### K2 — instance: peer engines admit a transition that changes no state

- Query (StateMachineRes): `transition with actions but no target state,
  internal or targetless transition, output on transition`
- Accepted (1): `study/stateless/README.md` §Internal transitions —
  "Sometimes a trigger needs to be handled, but the state shouldn't change.
  This is an internal transition." ⇒ a row whose effect is *not* a state
  write is an existing construct, not a deformation of the machine model.
- Accepted (2): `repos/scxmlcc/doc/user-manual.md` §transition — "At least
  one of the attributes `event`, `cond` or `target` must be specified" and
  "`D` is omitted if the transition has no target … This is for a transition
  without target". ⇒ SCXML's grammar makes `target` optional per transition;
  the output of a targetless transition is its executable content — the
  analogue of an authored `emit`, not of a pseudo-state write.
- Rejected: StateSmith / stateless CHANGELOG entries (bug-fix notes, no
  contract).

### K3 — instance: peer engines *require* an initial state, and the dummy-state workaround is prior art too

- Query (StateMachineRes): `initial state required for a state machine
  definition; machine without initial state`
- Accepted: `study/transitions/README.md` — "A state machine needs to start
  at some _initial state_." and "If you don't provide an initial state in
  the state machine constructor, `transitions` will create and add a default
  state called `'initial'`. If you do not want a default initial state, you
  can pass `initial=None`." ⇒ (a) a *machine* without a root is not a legal
  machine in the peer either — 0006:C18 is aligned, not idiosyncratic; (b)
  the silently-injected dummy state is exactly Alternative A's
  dummy-owned-tag convention, and the peer needed an explicit opt-out
  (`initial=None`) to escape it — evidence that the class is best *declared*,
  not inferred from what is missing.
- Rejected: jfsm README reference list, statewright specs (no contract text).

### K4 — sibling design doc (Domain prior, `../state-machines`)

- `MODEL-transition.md` §3 Outputs: "same tags → same row → same output
  (pure match)"; §3a: the recognized outcome is "consumed — dropped from the
  output regardless". ⇒ the resolver's answer is already framed as an
  *output* keyed by the selected row, so a row-keyed `emit` block is the
  design doc's own vocabulary, and a decision table is the degenerate
  machine whose output carries no next-state.

## In-repo prior art (source anchors the choice rests on — quoted in the RDR body)

- `internal/graphlint/reach.go::reach` — "No declared root: the traversal
  has nothing to start from."; `internal/graphlint/groups.go::checkGroups`
  — overlap/redundancy/coverage run only for a reachable group. ⇒ without a
  root, coverage over a stateless model is *never computed*; the class needs
  a root (∅) defined, not merely the missing-root finding suppressed.
- `internal/graphlint/analysis.go::nodeSatisfiesMatch` consults only owned
  atoms ⇒ the ∅ root satisfies every owned-atom-free context, so all groups
  of a decision table are reachable and coverage runs unchanged.
- `internal/table/dump.go::dumpColumns` closed vocabulary ⇒ an `emit` row
  field is a dump column (0002:C19 "carries every field").
- `internal/table/model.go::Model.Metadata` "carried through untouched and
  never interpreted (`0002:C2`)" ⇒ the uninterpreted-literal precedent
  `[rule.emit]` reuses.
- `internal/table/normalize.go::CatMalformedRuleShape` arm "ordinary rule …
  carries no write block" (`0002:C4`) ⇒ the clause this RDR conditions.
- `internal/cli/flow_exec.go::invokedReaders` demand = `RequiresOwned` ∪
  guard-owned keys ⇒ empty when no owned tag exists, so `--artifact` is
  already unrequired by the existing scoping rule (0005:C1 "A reader no
  candidate row needs MUST NOT run").
- Sibling discriminator: `Row.Kind` (ordinary/escape) is inferred from the
  *presence* of `escape` (`internal/table/normalize.go`, dump column `kind`);
  no adjacent path infers a class from an *absence*.

## Rejected branches

- pytransitions-style auto-injected default state (Alternative A) — the
  workaround the seed already refused.
- A separate `flow decide` verb — no peer splits "select the row" from
  "apply the transition" at the verb; SCXML keeps targetless transitions in
  the same `<transition>` element.
- Richer (`any`-typed) emit values — `[model.metadata]` precedent exists but
  a dump column needs a deterministic rendering; kept string-valued and
  raised as A6.
