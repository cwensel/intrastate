# RDR 0012 — Stage 2 grounding micro-sweep (step 7.5)

Model: claude-fable-5 (parent record of a factored fresh-context
sub-agent's returned packet; the checker saw anchors only, never the
argument)
Date: 2026-08-28

Verdict: 24 anchors — 23 CONFIRMED, 1 REFUTED (cosmetic, fixed in
place before commit).

Confirmed (anchor → claim): `internal/guard/grammar.go::Evaluator`
(struct{}; raw-string eq/in arms); `::parseHeldSet` (held null →
unevaluable); `internal/resolve/guardcontract.go::TestGuardEvaluatorContract`
(unparseable present value unevaluable-never-false; `(t, seam)`
signature); `0007:C1` (seam fence, MAY-answer-unevaluable, four atom
field names mirrored by UndecidedAtom); `0007:A17` (Evidence contains
"The declared kind needed to *parse* the value is known to the
evaluator from the table it was built for, not from the kernel");
`0007:A7`; `0007:A27`; `0007:C2` ("never to false and never to
true"); `0007:C8` (reason `uncomparable`); `0003:C8` (literal
parse-rejection before resolution); `0003:C9`;
`guard_evaluator_0003_test.go::TestReq34…` (arity + no-TagSet +
zero-field checks); `table/model.go::KernelTable` (no declarations
cross); `flow_resolve.go::guardSeam`, `product.go::valueSatisfies`,
`reach.go::atomAdmitsValue` (the three `Evaluator{}` sites);
`declaration.go::DeclarationOf`; JDR 0001 §D1 (atom shape) and JD-18
(still Open, no DECIDED line); `resolve.go` REQ-7 five-kind refusal
pin (comment at `TagSet.matches`); W3C SCXML tests 309/344 quotes
(uscxml `test/w3c/TESTS.md`); BehaviorTree.CPP
`PORT_CONNECTION_RULES.md` convertFromString quote;
`resolve.go::GuardEvaluator` interface.

Refuted (fixed): the sweep brief attributed "the seam takes whatever
RDR 0003 shows the evaluator needs (e.g. the declared kind on the
atom) — but never the view" to `0007:A18`'s If-wrong; it is
`0007:A17`'s If-wrong. `0007:A18`'s actual If-wrong is the two-valued
seam / programmer-error return. Non-load-bearing (the carrier choice
rests on A17's Evidence, independently confirmed); the Key
Discoveries bullet was re-attributed to A17 in the same pass.
