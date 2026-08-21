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

---

# iter-2 (2026-08-21) — re-propose after JDR 0001 §D1

Model: claude-fable-5

New question opened by §D1 (Row carries parsed atoms): **where is a
partial-domain rule enforced — by the host kernel, or by the pluggable
operator evaluator?** The domain rule itself (absence ⇒ unevaluable) is
unchanged and keeps C1/C2 above plus A4's K3 citation.

## Accepted citations

### C3 — SQL:2003: the HOST short-circuits null arguments; the routine is never invoked

Source: ISO/IEC 9075 SQL:2003 working draft 5WD-13-JRT-2003-09, p.169
(corpus `DevRef`, `sql_2003_standard/5WD-13-JRT-2003-09.pdf:page169`).

> "When an SQL function is called whose CREATE FUNCTION statement specifies
> RETURNS NULL ON NULL INPUT, then if the runtime value of any argument is
> null, the result of the function call is set to null, and the function
> itself is not invoked."

Corroborated by PostgreSQL 13 CREATE FUNCTION, p.1708 (`DevRef`,
`postgresql-13-US.pdf:page1708`): `CALLED ON NULL INPUT` (default) vs
`RETURNS NULL ON NULL INPUT` / `STRICT`.

Reading: the standard treatment of a partial function over possibly-missing
input is that the **engine** enforces the domain (null in ⇒ unknown out,
routine not called); the routine author is not trusted to re-implement
null propagation. Maps onto the instance question directly: the kernel
decides presence and never calls the evaluator for an absent key.
⇒ supports kernel-side enforcement (chosen approach B) over
evaluator-side discipline (alternative A).

## Queries run (budget ledger, iter-2)

1. DevRef: "function declared STRICT is not called when any argument is
   null; null result assumed automatically" — hits: PostgreSQL tutorial
   STRICT example pages (pg15/16/17 ~p.1269–1343), pg17/18 STRICT+VARIADIC
   note. Confirms the mechanism; not the definition. Partial.
2. DevRef: "null handling done by the engine instead of the user-defined
   function; three-valued logic unknown propagation enforced by host" —
   hits: Date's anti-null argument, SQL programmer's handbook ("NULLs
   propagate in calculations"), Datalog 3-valued instances. Off-instance.
   Rejected.
3. DevRef: "RETURNS NULL ON NULL INPUT or STRICT ... the function is not
   executed when there are null arguments" — hit SQL:2003 p.169 (C3) and
   PostgreSQL 13 p.1708. Accepted. Budget stopped here.
