# Persona 2 - Implementer

Question: if I started coding this Monday, what would I ask in the first hour?

Findings:

1. Which package owns the lint engine boundary? The RDR says "Define the lint
   package boundary over the normalized graph API," but does not name whether
   the command imports a graphlint package, a flow-model package method, or an
   RDR-0005 command helper.

2. What exact input object does lint consume after RDR 0002 normalization? The
   RDR lists source rule ids/spans, tags, outcomes, finite domains, terminals,
   and writes, but does not state the minimum normalized graph interface as a
   contract an implementer can compile against.

3. What are the mandatory `CLIError.Code` values? The RDR requires stable codes
   and says new codes are needed, but implementation cannot write tests or map
   failures until the stable code names are fixed.

4. What is the failure grouping for lint errors? The RDR says codes map to the
   existing CLI failure envelope and exit-code groups, but leaves whether graph
   defects are usage/config/data errors to implementation.

5. How are multiple blocking findings reported through `respond.Fail`, given
   `CLIError` is a single envelope? The RDR requires every blocking finding to
   carry rich identity, but does not say whether the command returns one
   aggregate lint error with details or one error per finding.
