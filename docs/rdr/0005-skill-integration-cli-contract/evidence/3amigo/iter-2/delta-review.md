# 3amigo Delta Review — Iteration 2

Scope: original ledger entries O1-O4 only.

## O1 next payload contract

Closed. Technical Design now names the minimum `next` payload fields
(`outcomes[]`, `candidates[]`, selected flow/model identity, supplied tags) and
Normative Contracts define candidate summaries with source rule identity,
required supplied facts, unresolved guard/accessor facts, and preview next tags
or write targets when available without evaluating missing facts.

## O2 request input grammar

Closed. Technical Design now pins MVP `--flow <id>`, repeated
`--tag name=value`, duplicate tag refusal, `--artifact role=path`, and deferred
structured input.

## O3 stable refusal code list

Closed for the draft. Failure Modes now lists minimum `flow-*` code strings and
the exit-code group for each.

## O4 set-state semantics

Closed. Technical Design, Normative Contracts, Illustrative Code, Risks, and
Validation now distinguish context `--tag` values from planned owned-tag
`--write name=value` mutations and read-back-confirmed values.

## Remaining Items

No open 3amigo findings. A6 remains Pending for reconcile/MVV verification.
