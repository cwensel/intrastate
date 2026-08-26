Model: claude-fable-5

# cli/0011 — Stage 2 prior-art ledger

Stage: 2-propose. Date: 2026-08-26. Budget: ≤3 corpus queries per claim, ≤5
opened hits per claim; instance reads via the `../state-machines` checkouts.

## Claim 1 — peers condition a "what next" query on the current state

Corpus queries (arc, `--corpus StateMachineRes --limit 5 --json`):

1. `query the transitions available from the current state: permitted triggers, next events, get_triggers`
   → hit: `repos/qmuntal-stateless/README.md` §Introspection (accepted, opened).
2. `list of events or transitions enabled from the current state with guard conditions evaluated versus declared regardless of guards (nextEvents, can, may_)`
   → no new accepted hit (scxmlcc `event`, awf-cli events, StateSmith changelog — rejected: event vocabulary, not an enumerate-candidates operator).
3. `SCXML selectTransitions enabled transition: event matches and cond evaluates true; decision table applicable rules hit policy`
   → hit: `repos/scxmlcc/doc/user-manual.md` `cond` (accepted); state-machine-cat SCXML docs (rejected: serialization, not selection).

Accepted citations (opened and quoted):

- qmuntal-stateless `README.md` §Introspection: "The state machine can provide a
  list of the triggers that can be successfully fired within the current state
  via the `StateMachine.PermittedTriggers` property."
  Source: `statemachine.go::StateMachine.PermittedTriggersCtx` →
  `states.go::stateRepresentation.PermittedTriggers`, which appends a trigger
  only when `len(tb.UnmetGuardConditions(ctx, unmet[:0], args...)) == 0`.
  ⇒ conditioned on current state AND guards; no unconditioned form offered.
- pytransitions `study/transitions/transitions/core.py::Machine.get_triggers`:
  "Collects all triggers FROM certain states" — `[t for (t, ev) in
  self.events.items() if any(name in ev.transitions for name in names)]`, no
  condition evaluation; `Machine._can_trigger` (the `may_<trigger>` family)
  iterates `get_triggers(state)` and evaluates conditions.
  ⇒ declared-shape enumeration exists but as a distinct operator from the
  evaluated one.
- xstate `repos/xstate/packages/core/CHANGELOG.md`: "d3d6149c7: Removed
  `MachineSnapshot['nextEvents']`." (v5); `packages/core/src/State.ts::machineSnapshotCan`
  computes `this.machine.getTransitionData(this, event)` and requires a
  non-forbidden transition.
  ⇒ the unconditioned enumeration was retired; the conditioned query survived.
- scxmlcc `repos/scxmlcc/doc/user-manual.md` `cond`: "The transition is only
  executed if the condition evaluates to true."
  ⇒ enablement is event- and condition-conditioned.

Verdict: class frame = "next" is state-conditioned; instance disposition =
conditioned is the default operator, enumeration is secondary (a separate
operator in pytransitions, removed in xstate, absent in stateless).

## Claim 2 — a peer distinguishes an ABSENT state key from a mismatch at the match seam

No query budget spent beyond the three above; none of the opened hits treats
an absent current state (every peer machine always has one).
⚠ no prior-art coverage for absent-vs-mismatch at the match seam; C1's
three-valued rule rests on this codebase's posture
(`internal/cli/flow_next.go::excluded` doc; `0007:C8`) and is verified by
A2/A3/A5.

## Claim 3 — flag naming for the enumeration surface

Not searched (not load-bearing; naming is a Load-Bearing Decision with the
rejected spellings recorded in the RDR). Model prior.

## Rejected branches

- `--all` as a filter on `outcomes` (alphabet) — 0005 DEV-4 already refused
  filtering the alphabet; not re-opened.
- javascript-state-machine `transitions()` / `can()` — checkout has no
  source-level `transitions()` accessor under `src/StateMachine.js` at the
  grep'd names; not opened further (budget).
- xstate `nextEvents` semantics in v4 — the removal note suffices for the
  disposition; v4 source not opened.
