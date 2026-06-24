# Persona 1 — Product Manager

1. The Problem Statement says a skill author needs to ask "where next", but
   `flow next` is specified mostly as a legal-outcome alphabet. The user outcome
   is unclear when a skill needs the actual next-work instruction or recognizer
   prompt shape, because the RDR leaves the `next` result payload at "plus
   conditional summaries".

2. The Technical Design says commands accept a "model or flow identifier plus
   state tags supplied as flags or a structured input file". The user outcome is
   unclear for first use because the command contract does not say which input
   form is normative for the MVV or how a skill author discovers the selected
   model revision.

3. The Normative Contracts require stable CLI refusals, and Failure Modes names
   the failure classes, but the RDR does not give the user a stable code list.
   The outcome for scripted callers is therefore under-specified: callers know
   codes exist, not what to branch on.

4. `flow read-state` / `flow set-state` promise state I/O through caller-supplied
   artifact role bindings, but the RDR does not define the minimum binding
   grammar beyond examples like `state:RDR_FILE`. A skill author cannot tell
   whether artifact binding is a CLI flag convention, a model identifier, or an
   accessor reference.
