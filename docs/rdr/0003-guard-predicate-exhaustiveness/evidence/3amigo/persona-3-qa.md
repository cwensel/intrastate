# Persona 3 - QA / Tester

1. **Cannot write a complete operator matrix test.**
   - Passage: `Normative Contracts`, `Validation`.
   - Blocker: the RDR says every operator has accepted tag kinds, but the matrix
     is not enumerated. QA cannot tell which combinations should pass and which
     should fail.

2. **Cannot decide whether missing `contains` coverage is a failure.**
   - Passage: `A1`, `Minimum Viable Validation`, `Testing Strategy`,
     `Performance Expectations`.
   - Blocker: the checked Resolve fixture does not use `contains`, while the
     validation scenario says representative rows use set containment. QA needs
     the RDR to say that `contains` is required in the implementation MVV even
     though it was not part of the Resolve subset.

3. **Cannot test `unless` without an explicit positive and negative case.**
   - Passage: `Testing Strategy`.
   - Blocker: the validation names mixed `all`/`unless` guards but not the
     expected disabled-row case. QA needs at least one assertion where all
     positive predicates match and the row is disabled solely because the full
     `unless` set matches.

4. **Cannot verify source-order independence fully.**
   - Passage: `Validation`.
   - Blocker: the current scenario says reorder rows and guard atoms. QA also
     needs a negative case proving two matching rows become an ambiguity refusal
     rather than first-match success.
