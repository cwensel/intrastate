# Persona 3 - QA / Tester

Question: how do I test this? What are the pass/fail criteria?

Findings:

1. I cannot write a boundary test for stateless replay if the RDR leaves accessor execution inside the resolver. A test with the same logical tuple but a different artifact read result would pass or fail depending on whether the resolver is allowed to perform hidden reads. The RDR needs one boundary: resolver consumes an explicit owned snapshot and returns planned writes, while accessor execution is tested by RDR 0004.

2. I cannot write the full refusal matrix from the current validation section because unmodeled recognized outcome is normative but not listed as its own scenario. Pass/fail should say an outcome tag absent from the model returns the typed unmodeled-outcome refusal and performs no persistence side effect.
