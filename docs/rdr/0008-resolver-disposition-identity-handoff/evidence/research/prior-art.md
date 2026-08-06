# RDR 0008 prior-art selection record

## Accepted anchors

- RDR 0002 `Normative Contracts`: "Each candidate row MUST retain its source
  rule id and source locator." `Load-Bearing Decisions / Identity` adds model id
  and deterministic expansion suffix.
- RDR 0005 `Technical Design`: the `flow resolve` payload includes "matched
  rule identity."
- Stateless `src/Stateless/Transition.cs::Transition`: one successful
  transition value exposes get-only `Source`, `Destination`, and `Trigger`.
- uscxml `src/uscxml/interpreter/LargeMicroStep.h::Transition` and
  `LargeMicroStep.cpp` take-transitions block: the selected transition retains
  its source element/context and the same element reaches before/after monitor
  callbacks. `WrappedInterpreterMonitor.cpp::afterTakingTransition` derives an
  XPath, source, and targets from that selected element.

## Queries

1. Semble, `../state-machines`: `selected transition result carries rule
   identity source location diagnostics audit resolution disposition transition
   plan`.
2. Semble, Stateless clone: `Transition object source destination trigger
   passed to state change notification execution context`.
3. Semble, uscxml clone: `selectTransitions returns enabled transition elements
   passed to microstep exitStates execute transition content enterStates`.
4. Arc, attempted `StateMachineOS`: `selected transition object identity
   metadata source state event execution result transition context`; rejected
   because that corpus is no longer configured.
5. Arc, `StateMachineLit`: `selected transition identity source metadata carried
   from transition selection into execution result diagnostics audit trace`.
6. Arc, `StateMachineLit`: `SCXML selectTransitions enabled transition set
   transition element source state execution microstep`.

## Rejected branches

- The two `StateMachineLit` queries found general next-configuration,
  observability, and event-history material but no load-bearing carrier design;
  none is cited as support for the choice.
- Returning the entire selected transition/Edge was not adopted from uscxml:
  intrastate callers do not need guard or normalized-table internals.
- Reconstructing rule identity in the CLI was rejected because it introduces a
  second lookup that can disagree with the edge actually selected.

## Stage 4 cold-path re-resolution record

### Accepted anchors

- `internal/resolver/resolver.go::Edge`, `::TransitionPlan`, and `::plan`: the
  current resolver has no selected-rule carrier, and `plan` already owns the
  defensive action copy.
- `internal/resolver/resolver.go::Resolve`: ordinary and modeled-escape
  selections use the same `plan` boundary.
- RDR 0002 `Normative Contracts` and `Load-Bearing Decisions / Identity`: every
  normalized candidate row retains source identity and locator together;
  modeled escapes remain normalized candidate rows.
- RDR 0005 `Technical Design` and `Normative Contracts`: `flow resolve` is a
  translating consumer of matched rule identity, not a second resolver.
- `evidence/spikes/command.txt` → `evidence/spikes/output.txt`: the live Go
  spike reran green and preserved surrogate string identity and cloned action
  fields after input mutation and reuse.

### Queries

1. Semble, intrastate: `resolver plan copies selected rule identity source
   locator action TransitionPlan from matched normalized row without mutable
   aliases`.
2. Semble, intrastate: `validated normalized table constructor binds match
   predicates selected rule identity source locator action and production
   Resolve accepts opaque table`.
3. Semble, intrastate: `modeled escape disposition semantics selected escape
   action next tags writes success refusal resolver`.
4. Semble, RDR peers: `Each candidate row MUST retain source rule id source
   locator model id deterministic expansion suffix escape`.
5. Semble, RDR peers: `flow resolve minimum payload matched rule identity`.
6. Semble, RDR peers: `RDR 0001 resolution kernel normative contracts modeled
   escape successful plan downstream values`.

### Negative, incomplete, and rejected branches

- Reuse audit negative: no production selected-rule identity carrier or
  discriminator exists in `internal/resolver`; extending `plan` does not
  duplicate an adjacent capability.
- A2 remains incomplete: the live spike uses a local string-locator surrogate;
  production has no selected-rule or locator carrier and accepts raw `[]Edge`.
- A3 remains incomplete: Final RDR 0005 describes matched identity and action
  in `Technical Design`, but its `Normative Contracts` require only next tags
  or refusal.
- A4 remains incomplete: RDR 0002 retains source identity but names no typed
  locator constructor or mismatch-validation boundary.
- A5 remains incomplete: no production normalizer, opaque validated-table
  value, or construction boundary exists; Draft RDR 0007 records the same gap.
- A6 is contradicted, not merely uncited: RDR 0001 and current `Resolve` make an
  exactly-one escape successful, while RDR 0002 forbids escape action fields.
  Resolve must return to Stage 2/3 to choose one cross-RDR meaning.
