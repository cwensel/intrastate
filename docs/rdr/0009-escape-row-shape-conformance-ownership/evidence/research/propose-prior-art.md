Model: claude-fable-5

# RDR 0009 — Stage 2 prior-art pass (bounded)

Problem class: ownership of structural row-shape conformance at a
table-producer / table-consumer seam (who guarantees "an escape row bears
no writes and no clears," and in what form).

## Corpus queries (budget: 3 run, 0 accepted)

1. `arc search semantic --corpus StateMachineRes --limit 5 --json "where is
   transition validity enforced: load-time validation of machine definition
   versus runtime checks"` — top hits were section headers only
   (fizz/docs/api.md, ms-conductor plan docs); nothing on the problem
   class. Rejected.
2. `arc search semantic --corpus StateMachineRes --limit 5 --json "make
   illegal states unrepresentable in transition definitions with distinct
   types for different transition kinds"` — hits were statewright/fizz
   design notes unrelated to row-shape conformance. Rejected.
3. `arc search semantic --corpus StateMachineLit --limit 5 --json
   "validating statechart or SCXML document structure before execution:
   schema constraints on transitions"` — scores ≤0.50, runtime-enforcement
   and scenario-based-programming papers; not this problem class. Rejected.

⚠ no prior-art coverage for structural row-shape conformance ownership in
the arc corpora; external prior-art claims (SCXML document-conformance
rejection at load; "parse, don't validate" shape-level doctrine) remain
model-prior and are deliberately NOT load-bearing — demoted to a Resolve
assumption (Method: Prior Art) in the RDR.

## Accepted in-repo citations (load-bearing, quoted in the RDR)

- RDR 0002 Normative Contracts: "An escape rule MUST contain an `escape`
  list and MUST NOT contain a write block or clear list." and
  "Normalization MUST render an escape rule as a candidate row with row
  kind `escape`, its normal predicate set, source rule id, source locator,
  and modeled failure class list." (writes absent from the escape-row
  rendering); validation categories include "malformed escape
  declarations."
- RDR 0001: "Modeled refusal is a value-level resolver disposition, not a
  CLI error and not the Go error path for parser bugs, IO failures, or
  programmer mistakes." — the Go error path is already reserved for
  producer/programmer contract breaches; taxonomy is "exactly" five kinds
  (A5 design decision).
- `internal/resolve/resolve.go::Resolve` doc: "the error return is
  reserved for programmer mistakes, not for modeled refusals."
- `internal/resolve/resolve.go::Row` (`Escape`, `Writes` coexist),
  `resolve.go::planOf` (copies `Writes` unconditionally, `Escaped` flag),
  `resolve.go::rescues` + the candidate loop's `len(row.Escape) != 0`
  partition (escape-ness is already discriminated by the Escape field —
  the existing sibling signal).
- RDR 0001 Existing Infrastructure Audit row: "Transition table | RDR 0002
  | Deferred to peer | Kernel assumes a parsed, reviewable table shape."

## Consequence for enumeration

The seed's branch 3 premise ("runtime predicate … requires a new refusal
kind") is only half true: runtime enforcement via a NEW REFUSAL KIND does
open RDR 0001's closed taxonomy, but runtime enforcement via the RESERVED
GO ERROR PATH does not — RDR 0001 pre-allocated that channel for exactly
this breach class. The candidate set was enumerated with both runtime
sub-branches distinct.
