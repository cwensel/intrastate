# RDR 0021 Stage-2 prior-art record

Scope: Propose-stage reads only (selection, not verification). Budget:
≤3 corpus queries per claim, ≤5 opened hits per claim.

## Claim 1 — the seed's assigned role and rejected alternative are real

Freshness re-validation of the Problem Statement's citations, read
directly in the `../state-machines` sibling repo (Domain priors row):

- `../state-machines/AUDIT-rdr-flow.md`, CW8 row: "Spec the legal graph
  once and *verify no illegal edge*. This is the sole place an FSM tool
  ... fits — as a **spec/validator, not a runtime**". CONFIRMED — the
  spec/validator role assignment is verbatim.
- `../state-machines/research/state-machines-research.md` §7d ("Runtime
  enforcement, verification & production state-machine patterns for
  agents") exists as cited. CONFIRMED.

⇒ Both seed citations resolve; Problem Statement stands.

## Claim 2 — peer state-machine CLIs export the machine as a
## format-selected document on stdout (class + instance)

Corpus: `StateMachineRes` (preferred prior-art corpus).

Queries run (semantic, limit 5, --json):

1. "export state machine definition as JSON serialization for diagrams
   and external tools" — hits: jfsm README (references only),
   StateSmith diagram-features, state-machine-cat SCXML doc,
   javascript-state-machine intro, statewright spec. Opened:
   state-machine-cat docs (accepted, see below).
2. "generate graphviz dot output visualization from state machine CLI"
   — hits: statewright pipeline, ms-conductor web-ui brainstorm, Archon
   CLI docs, state-machine-cat hacking doc. No stronger instance than
   state-machine-cat; stopped.
3. "machine definition serialized to JSON as the canonical interchange
   for visualization and analysis tools" — hits thin (scxmlcc manual
   stub, Archon prompts, ragel-colm CLAUDE.md, fizz). Rejected branch:
   none of these carries a quotable export contract; did not widen.

Accepted instance citation (opened and quoted):

- `../state-machines/repos/state-machine-cat/README.md` (usage block):
  `-T --output-type <type>  svg|eps|ps|ps2|dot|smcat|json|ast|scxml|…`
  and the piped invocation `smcat -T dot docs/sample.smcat -o - | dot
  -T svg …`. ⇒ The peer pattern is ONE exporter surface with a
  format-selector flag, emitting the chosen document as a raw stream to
  stdout/file; JSON and DOT are sibling output types of one surface,
  not two surfaces.

Rejected branches:

- StateSmith / scxmlcc: code generators (input → generated code), not
  a validator exporting its own normalized view; wrong direction for
  this problem.
- ms-conductor `--web` brainstorm: a server/UI surface, out of scope
  for a CI-diffable artifact.

⚠ No corpus coverage found for "graph document embedded in a JSON
response envelope" — no opened peer wraps DOT inside a JSON envelope;
the observed pattern is always the bare document on stdout. Treated as
a negative result, not widened.

## Claim 3 — in-repo priors (Domain priors, read directly)

- `docs/cli-output-contract.md`: `--as=json` → exactly one terminal
  envelope on stdout; `--as=text` → verb-defined human output.
- `0005:C1`: "Other command groups (lint, dump, parse) are outside this
  contract and are owned by the RDR that names them" — a new root
  export verb is anticipated, unowned, and free to claim.
- `0002:§round-trip-inverse-invariants`: the dump "does not define a
  dump grammar", its set-valued rendering is lossy by decision, and "a
  re-readable dump grammar is a follow-up RDR, seeded from this pass" —
  the JSON export is that follow-up's machine-readable half.
- `internal/cli/respond/respond.go::TextLiner`: the existing gateway
  seam for a payload whose whole content is one canonical string the
  text branch prints verbatim ("an identity string a caller pipes or
  greps") — the seam a bare document rides in text mode.
