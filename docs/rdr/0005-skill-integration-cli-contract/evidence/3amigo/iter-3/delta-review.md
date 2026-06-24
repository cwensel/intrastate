# 3amigo Delta Review - Iteration 3

Scope: cluster re-entry defect for RDR 0004/RDR 0005 gate accessor handoff.

## Persona 1 - Product Manager

The cluster fix correctly prevents `flow read-state` from returning gate results
as tag values, but the user outcome remains incomplete unless skill callers know
which command maps gate allow/deny/indeterminate. A skill author still needs one
deterministic CLI contract for gate-conditioned transitions.

## Persona 2 - Implementer

The current draft says `read-state` invokes only read accessors and says gate
results surface as candidate facts or typed refusals, but it does not explicitly
say whether `next`, `resolve`, or some later command invokes declared gate
accessors, nor how `--artifact role=path` reaches that invocation path.

## Persona 3 - QA / Tester

The validation matrix names gate indeterminate during accessor fact evaluation,
but the pass/fail test cannot identify the command under test. If `read-state`
must not invoke gates, then a gate-indeterminate CLI refusal needs to be tested
through the verb that evaluates gate facts.

## Consolidation

O5 gate accessor invocation surface - PM, implementer, and QA all flagged that
the post-cluster draft distinguishes read accessors from gate accessors but does
not pin which CLI verb invokes gates and how artifact bindings are supplied.

Highest-priority rewrites: A3/A4, Technical Design, Normative Contracts,
Implementation Plan, and Validation.
