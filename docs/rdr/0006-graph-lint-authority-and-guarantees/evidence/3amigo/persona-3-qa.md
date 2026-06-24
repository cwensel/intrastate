# Persona 3 - QA / Tester

Question: how do I test this? What are the pass/fail criteria?

Findings:

1. The MVV names illegal fixtures per invariant class but does not enumerate the
   expected stable code for each fixture. QA cannot write pass/fail assertions
   for "stable finding codes" without the exact code list.

2. The MVV says source rule/context identity must appear in JSON mode "when the
   normalized model provides one," but it does not test the fallback path when a
   finding has only a source span or graph element id. That leaves a diagnostic
   guarantee untested.

3. The inability-to-prove scenario is named, but its relationship to guard
   exhaustiveness/gap is not pinned. QA cannot tell whether it is a separate
   mandatory code, a variant of coverage gap, or only a message detail.

4. The text/JSON scenario requires same semantic result and exit behavior, but
   does not say whether finding ordering must be deterministic. The RDR has a
   finding identity decision, but tests need an ordering or order-independent
   assertion rule.

5. CI-shaped invocation is required by MVV, but no pass/fail criterion proves
   CI is invoking the authoritative command rather than a wrapper or hook with
   drifted rules.
