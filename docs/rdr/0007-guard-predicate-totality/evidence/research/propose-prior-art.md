Model: claude-fable-5

# RDR 0007 — Stage 2 prior-art research cache

Problem class: what verdict a state-machine guard/condition yields when it
references data absent from the runtime view (undecidable-guard semantics),
and whether that verdict may be silently folded into "false."

## Accepted citations

### C1 — SCXML: cond-evaluation error is 'false' PLUS a mandatory observable error (never silent)

Source: W3C SCXML Recommendation, §5.9.1 Conditional Expressions
(https://www.w3.org/TR/scxml/), fetched 2026-08-11 via WebFetch.

> "If a conditional expression cannot be evaluated as a boolean value ('true'
> or 'false') or if its evaluation causes an error, the SCXML Processor MUST
> treat the expression as if it evaluated to 'false' and MUST place the error
> 'error.execution' in the internal event queue."

Related (§5.3.3 Late Data Binding): accessing data before it is loaded "MUST
yield the same execution-time behavior as accessing non-existent data
substructure" — i.e. absent data routes through the same error machinery.

Reading: the only major standard in this problem class that maps
unevaluable→false does so ONLY in tandem with a mandated, observable error
signal on a second channel (the internal event queue). Silent absence-as-false
has no prior-art support. Intrastate's kernel has no second channel — the
verdict itself is the only observable — so an honest transplant of the SCXML
rule surfaces the error in the verdict (the unevaluable verdict), not in the
boolean.

### C2 — In-repo sibling signal: the "cannot decide ⇒ distinct verdict, non-escapable refusal" decision already exists

`internal/resolve/resolve.go::GuardUnevaluable` (three-valued `GuardResult`),
`internal/resolve/resolve.go::KindGuardUnevaluable` (closed refusal kind),
`internal/resolve/resolve.go::evaluateGuard` — a guarded row with a nil seam
returns `GuardUnevaluable`, "never evaluated by the kernel itself." The
nil-seam case is the same class of undecidability (the evaluator cannot
answer) and it is NOT folded into false.

## Queries run (budget ledger)

Corpus queries (arc search semantic, --limit 5):
1. StateMachineRes: "SCXML cond attribute evaluation error treated as false
   error.execution" — hits were state-machine-cat/scxmlcc import-export docs;
   no spec semantics. Rejected.
2. StateMachineRes: "guard condition references undefined variable missing
   data evaluation semantics" — workflow-tool docs (inngest, Archon, awf-cli);
   off-class. Rejected.
3. StateMachineRes: "what happens when a transition guard cannot be evaluated
   error false statechart" — stateless CHANGELOG entries about guard
   precedence bugs; off-class. Rejected.
4. StateMachineLit: "three-valued logic unknown missing value predicate
   evaluation" — LLM-agent papers, off-class. Rejected.

⚠ no prior-art coverage in the arc corpora for the undecidable-guard /
three-valued-verdict problem class; corpus results above are all rejected.
External coverage came from one bounded web read of the W3C SCXML spec (C1).

## Demoted to Resolve assumptions (not load-bearing here)

- SQL/Kleene K3 three-valued logic as precedent for "definite results from
  present data dominate; unknown propagates otherwise" — known from the model
  prior; not opened/quoted this pass. Carried as a Prior Art assumption for
  Stage 4 (any standard SQL reference or Kleene logic text).
