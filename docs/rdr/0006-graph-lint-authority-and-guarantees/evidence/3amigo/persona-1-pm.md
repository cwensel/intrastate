# Persona 1 - Product Manager

Question: does this RDR actually deliver the user outcome?

Findings:

1. The user outcome is "illegal or incomplete transition graphs caught at design
   time," but the accepted input boundary for `intrastate lint` is still
   illustrative. The RDR names `intrastate lint --flow rdr --model ...` as
   command shape, while RDR 0005 owns flow command placement. A maintainer
   cannot tell which invocation CI must actually run for every transition-model
   change.

2. The RDR says every blocking finding carries a stable code, severity, model
   id, rule/context id when available, source span when available, and a concise
   message, but it does not pin the mandatory code taxonomy. That weakens the
   review outcome because "stable codes" could be satisfied by ad hoc names that
   do not map one-to-one to the invariant classes users need to act on.

3. CI authority is selected, but the enforcement boundary is under-specified.
   The RDR says CI invokes the same command for model changes, yet does not say
   whether a missing model, missing lint invocation, or use of a wrapper with
   different rules is itself a failure of the acceptance contract.

4. The RDR differentiates lint from runtime refusal well, but the success
   surface is vague. A maintainer needs to know whether success promises "no
   blocking findings for the supplied normalized model" only, or also asserts
   warnings/advisories are absent.
