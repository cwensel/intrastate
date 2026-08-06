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

## Stage 4 verification record

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
  spike passed and preserved selected identity and action after input mutation
  and reuse.

### Negative and rejected branches

- Reuse audit negative: no production selected-rule identity carrier or
  discriminator exists in `internal/resolver`; extending `plan` does not
  duplicate an adjacent capability.
- No normalizer or edge-construction implementation exists yet. That absence
  is an implementation prerequisite, not contrary evidence: Final RDR 0002
  owns the complete normalized-row identity contract consumed here.
