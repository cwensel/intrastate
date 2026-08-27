Model: claude-opus-5

# RDR 0009 — Stage 4 (Resolve) prior-art pass for assumption A6

Assumption under test (A6, explicitly NOT load-bearing; citation-only):

> "Load-time document/schema rejection is the standard enforcement point
> for table-shape conformance in peer state-machine systems (SCXML
> conformance model)."

Stage 2 (`propose-prior-art.md`) ran 3 corpus queries, all rejected, and
recorded the SCXML load-time-rejection claim as unverified model-prior.
This pass extends that ledger with new queries and does not repeat them.

## Query ledger (budget 4; 2 corpus queries run, 2 unspent)

1. `arc search semantic --corpus StateMachineRes --limit 8 --json "SCXML
   conformance: processor must reject a document that is not well-formed
   before execution begins"` — 8 hits, scores 0.80–0.86. Hit bodies were
   section headers only (`state-machine-cat/docs/SCXML.md`,
   `scxmlcc/doc/user-manual.md`, `uscxml/docs/DEVELOPERS.md`), but the hits
   correctly identified the source-level locations worth opening.
   **Partially accepted** — used as navigation, not as citation.
2. `arc search semantic --corpus PapersFast --limit 8 --json "static
   well-formedness checking of statechart models at load time rejecting
   malformed transition declarations"` — 8 hits, top score 0.757, all
   off-topic (LLM test generation, runtime enforcement of reactive systems,
   autoregressive reasoning). No coverage of the problem class.
   **Rejected.**

Corpora not queried this pass: `StateMachineLit` (covered by Stage 2 query
3, rejected), `DevRef` (budget stopped early — a decisive hit was found).

Source-level follow-up (sibling repo `../state-machines`,
literal sweep, not a corpus query): scanned peer-engine markdown for
load-time rejection language across `repos/xstate`, `repos/statewright`,
`repos/StateSmith`, `repos/scxmlcc`, `study/uscxml`.

## Hits opened (5 of 6 budget)

- `../state-machines/repos/scxmlcc/doc/user-manual.md`
  — 11 `#### Valid Children` blocks. Documents a per-element child grammar
  (a load-time document shape), but states no rejection obligation.
  Rejected as a citation for A6.
- `../state-machines/study/uscxml/test/w3c/TESTS.md`
  — the W3C SCXML IRP conformance test matrix, whose `alt` attributes quote
  the normative spec assertion each test covers. **Accepted.** Decisive.
- `../state-machines/repos/state-machine-cat/docs/SCXML.md`,
  `repos/xstate/packages/core/CHANGELOG.md`,
  `repos/StateSmith/docs/plantuml-input.md` — surveyed; no normative
  statement about the enforcement point for structural conformance.
  Rejected.

## Accepted citation — and it FALSIFIES A6 as written

`study/uscxml/test/w3c/TESTS.md:1609`, the conformance assertion attached
to W3C SCXML IRP **test 313**:

> "The SCXML processor MAY reject documents containing syntactically
> ill-formed expressions at document load time, or it MAY wait and place
> error.execution in the internal event queue at runtime when the
> expressions are evaluated."

This is the SCXML conformance model speaking directly to the exact
enforcement-point question A6 raises, and it declines to standardize an
answer. Load-time rejection is **permitted (MAY), not required**, and
deferring to a runtime error channel is explicitly and equally conformant.
So SCXML is not evidence that load-time rejection is *the standard*
enforcement point for structural conformance; it is evidence that the
enforcement point is deliberately left to the implementation.

Contrast — the one place SCXML does mandate rejection is narrow, and it is
about an unresolvable external resource rather than document shape
(`TESTS.md:1470`, test 301):

> "If the script specified by the `src` attribute of a script element
> cannot be downloaded within a platform-specific timeout interval, the
> document is considered non-conformant, and the platform MUST reject it."

Adjacent load-time MUSTs found in the same matrix are bindings, not
shape validation (`TESTS.md:906` `_sessionid`, `:936` `_name`,
`:1485` evaluate `script` children at document load time).

negative: A6 — no corpus evidence in StateMachineRes, StateMachineLit,
PapersFast that load-time rejection is the *conventional* enforcement point
for structural conformance of a transition/state definition. The one real
citable source located (W3C SCXML IRP test 313, via uSCXML's conformance
matrix) states the opposite of A6's "standard": the choice between
load-time rejection and runtime error is left open as MAY.

## Consequence for the RDR

A6 is citation-only and explicitly NOT load-bearing, so this does not
disturb the design choice. The correct action is to restate A6 rather than
cite it: SCXML does not establish load-time rejection as conventional — it
declines to fix the enforcement point at all. Any RDR prose asserting a
peer-standard convention for load-time shape rejection should be dropped or
rewritten to say that peer systems leave the enforcement point to the
implementation, which if anything *strengthens* RDR 0009's freedom to
choose its own enforcement point on in-repo grounds (RDR 0001's reserved Go
error path, RDR 0002's normative escape-rule contract) rather than by
appeal to external convention.
