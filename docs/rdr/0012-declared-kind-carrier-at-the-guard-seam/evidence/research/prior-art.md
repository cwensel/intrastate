# RDR 0012 — Stage 2 prior-art search record

Model: claude-fable-5
Date: 2026-08-28

Corpus budget per prompt: ≤3 queries per claim, ≤5 opened hits per claim.

## Claim 1 (instance): what peer state-machine systems do with a guard
value the evaluator cannot type

Queries (arc, corpus `StateMachineRes`):
1. `SCXML conditional expression evaluation error treated as false error.execution` — hits: state-machine-cat SCXML docs, scxmlcc manual, uscxml DEVELOPERS (thin; led to direct open of the uscxml W3C manifest below).
2. `guard expression type checking declared types evaluation environment` — hit accepted: BehaviorTree.CPP `docs/PORT_CONNECTION_RULES.md` §"Type Conversion via convertFromString".

Accepted citations (quoted from source):

- **W3C SCXML conformance tests 309/344** (via
  `state-machines/study/uscxml/test/w3c/TESTS.md`):
  test 309 — "If a conditional expression cannot be evaluated as a boolean
  value (`true` or `false`) or if its evaluation causes an error, the SCXML
  processor MUST treat the expression as if it evaluated to `false`."
  test 344 — same premise, "MUST place the error `error.execution` in the
  internal event queue."
  ⇒ SCXML folds an unevaluable guard into FALSE plus a loud side-channel
  event. Intrastate's kernel (0007) deliberately rejects the fold-to-false
  half; the instance read confirms the disposition is a real design fork in
  peer systems, and BOTH dispositions require the evaluator to *detect* the
  type failure — which requires a kind carrier. SCXML's carrier is the
  datamodel: values are typed by the ECMAScript environment the expression
  runs in, not by a per-call parameter.

- **BehaviorTree.CPP port typing**
  (`state-machines/repos/BehaviorTree.CPP/docs/PORT_CONNECTION_RULES.md`):
  a string blackboard entry "can be connected to ports of **any type** that
  has a `convertFromString<T>()` specialization … Note that this may cause a
  run-time error if the string is not convertible."
  ⇒ declared type lives on the PORT (the declaration), conversion of the
  string-carried runtime value happens at read against that declaration, and
  an inconvertible value is a loud run-time error — never a silent false.

## Claim 2 (class): where type information travels for typed evaluation

Query (arc, corpus `StateMachineRes`):
3. `CEL expression evaluation typed environment declaration runtime type
   mismatch error` — no CEL coverage in corpus (hits were the same
   BehaviorTree/awf-cli/fizz docs). ⚠ no prior-art coverage for
   compile-checked expression environments (CEL-class) in the available
   corpora; the class claim rests on the two anchored instances above plus
   the in-repo prior art (`0007:A17` Evidence: "The declared kind needed to
   *parse* the value is known to the evaluator from the table it was built
   for, not from the kernel"), which is quotable and decisive on its own.

Rejected branches:
- SQL/K3 three-valued treatment — already load-bearing in closed `0007:A4`;
  not re-searched (budget).
- scxmlcc / state-machine-cat docs — opened via hit list only; no statement
  about cond-evaluation typing; rejected as non-load-bearing. (2 opened hits
  against Claim 1's budget of 5.)
