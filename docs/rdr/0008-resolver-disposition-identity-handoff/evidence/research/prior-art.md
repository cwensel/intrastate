# RDR 0008 prior-art selection record

## Accepted anchors

- RDR 0002 `Normative Contracts`: "Each candidate row MUST retain its source
  rule id and source locator." `Load-Bearing Decisions / Identity` adds model id
  and deterministic expansion suffix.
- RDR 0005 `Technical Design`: the `flow resolve` payload includes "matched
  rule identity."
- qmuntal-stateless `statemachine.go::internalFireOne` and
  `::handleTransitioningTrigger`: selection constructs one
  `Transition{Source, Destination, Trigger}` and passes it through execution
  and transition callbacks.
- scxmlcc `doc/user-manual.md::Transition` and
  `::Custom Actions and Conditions`: a transition may omit a target while its
  executable content remains a separately modeled choice.

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
7. Arc, `StateMachineLit`: `SCXML targetless transition no target state
   configuration executable content successful transition semantics`.
8. Web, W3C-only: `targetless transition no target state configuration
   executable content` and `transition without target executable content
   internal transition`.

## Rejected branches

- The two `StateMachineLit` queries found general next-configuration,
  observability, and event-history material but no load-bearing carrier design;
  none is cited as support for the choice.
- Returning the entire selected transition/Edge was not adopted from the prior
  art:
  intrastate callers do not need guard or normalized-table internals.
- Reconstructing rule identity in the CLI was rejected because it introduces a
  second lookup that can disagree with the edge actually selected.
- The additional `StateMachineLit` query returned general statechart semantics
  but no direct carrier contract.
- The formerly accepted Stateless and uscxml clone anchors were removed when a
  cold-path refresh found neither clone in the configured corpora or local
  prior-art checkout. Their exact selected-rule/source-locator claim is not
  carried forward.

## Stage 2 re-entry disposition

The revised candidate set treats selected-rule identity and action as separate
questions. Ordinary selection remains a successful plan with identity and
action. An exactly-one escape preserves its underlying refusal kind and carries
the selected escape row's identity, but no action; unmodeled kernel refusals
carry neither selection nor action. Successful no-change action, a third public
modeled-escape branch, and disposition-level optional identity were rejected by
the scored matrix in the RDR.

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

## Stage 4 external-evidence refresh

### Accepted anchors

- qmuntal-stateless `statemachine.go::internalFireOne` and
  `::handleTransitioningTrigger`: source, destination, and trigger are carried
  from selection into execution and transition callbacks.
- scxmlcc `doc/user-manual.md::Transition` and
  `::Custom Actions and Conditions`: target and executable content are
  independent transition properties; a targetless transition can still have
  an authored action.

### Queries

1. Arc, `StateMachineOS`: `selected transition object retains transition
   identity source context from selection through execution monitoring
   diagnostics`.
2. Arc, `OpenSource`: `state machine selected transition carries source
   destination trigger metadata into transition execution notification`.
3. Arc, `DevRefOS`: `selected transition object retains source destination
   trigger identity context into execution monitor callbacks state machine`.
4. Arc, `StateMachineRes`: `selected transition identity source destination
   trigger context carried through execution diagnostics monitoring`.
5. Arc, `StateMachineLit`: `targetless transition no target state configuration
   executable content semantics actionless refusal`.
6. Arc, `StateMachineRes`: `targetless transition actionless transition refusal
   no state change executable content semantics`.
7. Arc, `DevRef`: `SCXML targetless transition does not change state
   configuration invokes executable content`.
8. Arc, `StateMachineLit`: `targetless internal transition executable content
   action no state change transition semantics`.

### Negative and rejected branches

- Exact selected-rule/source-locator handoff: no surviving corpus evidence.
  `StateMachineOS`, `OpenSource`, and `DevRefOS` were unavailable, and the
  previously recorded Stateless/uscxml clone anchors were absent.
- The surviving prior art supports direct transition-context handoff but does
  not select or refute this RDR's split plan/refusal carrier. That taxonomy
  remains a project design decision grounded in the RDR 0002 escape schema.
