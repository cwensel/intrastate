# RDR 0008 — Stage 2 prior-art research cache

Model: claude-fable-5
Date: 2026-08-11

Problem class: who owns the *name* of an engine-injected datum inside a
namespace the author otherwise controls (declared tag vocabulary vs
kernel-injected recognized-outcome tag).

## Queries run (budget: ≤3 corpus queries/claim, ≤5 opened hits/claim)

Claim A — "established statechart/workflow engines fix the name of the
engine-provided event/outcome datum; author vocabulary cannot rebind it."

1. `arc search semantic --corpus StateMachineRes --limit 5 --json "reserved
   system variable event name SCXML underscore prefix author cannot rebind"`
   → top hits state-machine-cat/docs/SCXML.md, scxmlcc/doc/user-manual.md —
   no quotable system-variable passage. Rejected.
2. `arc search semantic --corpus StateMachineRes --limit 5 --json "_event
   _sessionid system variables protected read-only data model"`
   → hit: `repos/inngest/pkg/api/v2/README.md` — "Read-Only Fields: Some
   fields are computed or system-managed and cannot be modified through
   write operations." ACCEPTED (system-managed names reserved from user
   writes; workflow-engine peer).
3. `arc search semantic --corpus StateMachineLit --limit 5 --json "SCXML
   system variables _event reserved names data model"` → no SCXML spec text
   in corpus. Rejected; SCXML claim demoted to a Resolve assumption.

Follow-up (semble tier-2, then scoped rg literal sweep on the sibling
checkout named by resources.md as the interim source-level fallback):

4. `semble search "_event system variable reserved read-only scxml"
   ../state-machines` → located `repos/README.md` engine roster; no uscxml
   checkout present. XState checkout present.
5. `rg -n "context, event" repos/xstate/packages/core/src/guards.ts` →
   ACCEPTED anchor: `packages/core/src/guards.ts::GuardArgs` — guards are
   invoked with `{ context, event }: GuardArgs<any, any>`; the injected
   event occupies the framework-fixed field name `event`, and
   author-controlled data lives in `context`, a separate namespace. The
   author names event *types*, never the carrier slot.

## Accepted citations

- **XState** — `../state-machines/repos/xstate/packages/core/src/guards.ts::GuardArgs`
  (`{ context, event }` destructure at lines 132/194/267): the engine owns
  the name of the injected-event slot; author vocabulary is a disjoint
  namespace (`context`).
- **Inngest** — `../state-machines/repos/inngest/pkg/api/v2/README.md`
  "Read-Only Fields": system-managed fields are name-reserved against user
  writes.
- **In-repo (RDR 0002)** — the `<clear>` sentinel: "normalization renders as
  a `<clear>` write" (Normative Contracts) — an already-locked precedent for
  a carrier-owned reserved token adjacent to author-named vocabulary,
  syntactically outside declarable tag-name space.

## Demoted (not openable locally — Resolve assumption, not load-bearing here)

- W3C SCXML §5.10 system variables: the Processor MUST bind the processed
  event to `_event`; names beginning `_` are reserved; authors must not
  rebind system variables. Matches the accepted XState posture but the spec
  text is not in any local corpus/checkout — verify at Resolve (Method:
  Prior Art, fetch spec section).

## Negative results

- No local prior art found for the *model-carried key declaration* branch
  (an engine letting authors rename the injected event/outcome slot) — no
  surveyed engine does this. Recorded as a negative result, not widened.
