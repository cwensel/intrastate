# RDR 0019 — Stage 2 prior-art research

Scope: the class question (how do state-machine engines treat the declared
initial state at runtime) and the instance question (what does each peer do
when state is externalized / restored). Budget: 3 semantic queries against
`StateMachineRes`, opened hits read in the `state-machines` sibling checkout.

## Queries run (arc, corpus `StateMachineRes`)

1. "initial state at runtime: does the interpreter enter the declared initial
   state on start, versus restoring persisted state" — surfaced
   `repos/qmuntal-stateless/README.md` (Activation/Deactivation around state
   storage) and `study/transitions/README.md`.
2. "XState initial state createActor start persisted snapshot restore" —
   surfaced the xstate checkout; the load-bearing read was taken from source
   (`repos/xstate/packages/core/src/createActor.ts`), not the changelog hits.
3. "SCXML initial attribute configuration document start" — surfaced
   `repos/scxmlcc/doc/user-manual.md` (`<scxml initial="hello" ...>`).

## Accepted citations (quoted from source)

- **qmuntal/stateless** — `repos/qmuntal-stateless/statemachine.go`:
  `func NewStateMachine(initialState State) *StateMachine` and
  `func NewStateMachineWithExternalStorage(stateAccessor func(context.Context)
  (State, error), stateMutator func(context.Context, State) error, ...)`.
  The initial state is supplied by an explicit act at machine construction;
  with externalized state the *accessor* supplies current state and the
  library never synthesizes a fallback on read.
- **XState v5** — `repos/xstate/packages/core/src/createActor.ts`:
  `this._initState(options?.snapshot ?? options?.state)` — the
  restored-vs-initial fork is decided once, explicitly, at actor creation;
  a supplied persisted snapshot restores, otherwise the declared initial is
  materialized at start. No per-read fallback exists; "cleared" and
  "unseeded" are distinguished by which creation path the caller invoked.
- **SCXML** — `repos/scxmlcc/doc/user-manual.md`:
  `<scxml initial="hello" version="1.0" ...>` — the initial attribute names
  the configuration the interpreter enters when the session starts; entering
  it is a session-start event, not a read-time default.
- **In-repo prior** — `docs/model-authoring.md` ("`[initial]` declares the
  owned state a model starts from") already frames `[initial]` as a start
  state to authors, while no `internal/cli` code reads `Model.Initial`
  (grep over `internal/cli` hits tests only).

## Class conclusion

Prior art is uniform: the declared initial state is materialized by an
explicit lifecycle step at session start (constructor / `createActor` /
document start), never by an implicit fallback applied when a read finds
state missing. Engines with externalized state (stateless's external
storage; XState persisted snapshots) push current-state supply to an
explicit accessor or an explicit restore argument — the shape intrastate's
0004 accessor model already has.

## Rejected branches

- Changelog/migration-guide hits from query 2 (release notes, not
  semantics) — discarded without opening beyond snippets.
- `study/transitions/README.md` hit from query 1 (pytransitions state
  *object* attribute persistence, not initial-state entry semantics) —
  read, judged off-topic.
- No further sweep for workflow-engine start events (Inngest hit from
  query 1 was off-topic); the three engine reads above already converge,
  and the budget rule stops at confirmation.
