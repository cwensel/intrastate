# Resolve-stage research — the owned/readers partition question (ECHO vs PLAN)

Date: 2026-08-29. Question raised by the A7 source read: `owned` (artifact-state
snapshot) and `readers` (invoked reader identities) are classed ECHO
(projectable) by C2; strictly caller-reconstructible only via the caller's own
prior calls — should JDR 0002 §D1's unclear-joins-PLAN tiebreak move them?

## Corpus findings

**H1 — kernel results carry outputs, not echoed inputs: MIXED (supporting, with
a named gap).**

- DevRef, *Designing Data-Intensive Applications* (Kleppmann), p481 (Ch 11):
  command processing validates against current state inside one transaction and
  its result is the emitted event(s) — prior state is not echoed back.
- DevRef, *DDIA*, p480: current-state snapshots are "derived from the log of
  events… only a performance optimization" — the state side is reconstructible
  by the consumer.
- Gap (three passes now, propose + this stage): no corpus coverage for the
  DMN-engine result-echo question or the statechart
  `transition(state, event) → next state + actions` result-shape idiom. The
  decision does not rest on it; the local kernel signature answers the same
  question first-hand (below).

**H2 — caller-controlled result projection, input side droppable: SUPPORTED.**

- DevRef, *REST API Design Rulebook* (Masse), p92: the client trims response
  data with a `fields` query parameter, excluding "fields whose values are
  known to be sizable and unused."
- DevRef, *Microservices From Day One*, p220: JSON API sparse fieldsets —
  client-requested per-type field subsets for large-by-default responses.

**H3 — provenance rides as bounded identity, never re-echoed content:
SUPPORTED.**

- DevRef, *REST API Design Rulebook*, p53: ETag is "an opaque string that
  identifies a specific version" of the representation — the identity token
  stands in for content; freshness is checked by token comparison, not
  re-transfer.
- DevRef, *RESTful Web APIs* (Richardson/Amundsen/Ruby), p358: same doctrine,
  importance "very high"; revalidation via If-None-Match (*Microservices From
  Day One*, p105).

**H4 — chained-agent token economy of tool results: MIXED (premise grounded).**

- PapersFast, *ReAct* (Yao et al. 2022), p2: action selection "requires complex
  reasoning over the trajectory context (Question, Act 1-3, Obs 1-3)" — tool
  observations accumulate and are re-paid as context across turns.
- No passage on deliberate tool-response minimization; PapersFast's fast index
  returned citation-page fragments, so this absence is weak evidence.

## Queries (rejected branches recorded per the stage contract)

- H1 StateMachineRes: "transition function returns next state and actions
  without echoing prior state" — miss; "machine.transition is a pure function
  that returns the next state given current state and event" — miss
  (transition-config docs, not result shape). DevRef: "DMN decision table
  evaluation result output entries business rules engine" — miss; "command
  handler validates against current state and emits new events event
  sourcing" — hit.
- H2 DevRef: "partial response client selects which fields are returned" — hit.
- H3 DevRef: "version number ETag token identifies state for concurrency
  control instead of resending content" — hit.
- H4 PapersFast: "LLM agent tool output length context window cost
  truncation" — miss; "summarizing or compressing intermediate observations to
  reduce tokens in multi-step agent trajectories" — miss; "context engineering
  for LLM agents managing tool results in the context window" — miss; "ReAct
  agent interleaves reasoning with actions and receives observations from
  tools" — hit.

## Record-set and code precedent (audited this stage)

- Kernel purity: `func Resolve(in Input) (Result, error)`
  (`internal/resolve/resolve.go::Resolve`; `Input{Flow, Table, Owned,
  Observed, Recognized, Guards}`) — `owned` is a kernel INPUT; the plan is the
  output. JDR 0001 §D8 rejected caller-carried owned state (write-skew) and
  ordered the payload to gain "the read accessor identities and the assembled
  owned snapshot" when resolve began reading for itself — that is how
  owned/readers entered the envelope (0005 pins them only as minimum-shape
  members; no transparency/audit rationale is recorded).
- `readers` never reaches the kernel: `internal/cli/flow_exec.go::invokedReaders`
  is a pure function of the normalized model + the requested outcome — derivable
  from the caller's own request alone.
- Roll-forward: the normalizer renders write block + clear list as one
  assignment set cloned into both `NextTags` and `Writes`
  (`internal/table/normalize.go`, clear = `<clear>` sentinel write); the
  payload's `next` covers exactly writes ∪ clear, so prior-owned overlaid with
  `next` (sentinel = removal) is the successor owned view; untouched keys never
  change ("absence… never implies deletion"). `flow set-state`'s read-back
  refuses when any unplanned tag changed (`internal/cli/flow_exec.go`) — foreign
  drift is caught on the write path. An escaped plan mutates nothing.
- User-facing partition already shipped: `flowResolveExtendedDesc` "Reading a
  successful plan" lists exactly the plan group; `docs/cli-output-contract.md`
  never names owned/readers in prose.
- Peer stances: 0021:C5 — no caller-controlled projection on the export verb,
  and any later one MUST conform to JDR 0002 §D1. 0010:C4 — `emit` present as
  `{}` (empty-container distinction is load-bearing). 0011:D-identity — a mode
  bit joins request identity; reporting width may vary, decision may not.

## Disposition input

The corpus and precedent agree with keeping `owned`/`readers` in ECHO: ECHO =
the kernel's inputs (supplied or caused by the caller), PLAN = the kernel's
outputs plus bounded identity attestations (`revision`, the ETag form). A
future drift-detection ask is served by adding a bounded state-identity field
to PLAN under §D1's field-addition rule, not by re-echoing state content in
every projected payload. Final call recorded in the RDR (author's round).
