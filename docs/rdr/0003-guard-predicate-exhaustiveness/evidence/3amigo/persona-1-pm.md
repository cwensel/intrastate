# Persona 1 - Product Manager

1. **Set-containment scope is unclear.**
   - Passage: `Problem Statement`, `Key Discoveries`, `Normative Contracts`,
     `Minimum Viable Validation`.
   - Concern: the stated user outcome is conditional edges for cap/profile/lens
     routing and rewinds. The Resolve spike proves `eq`, `in`, `lt`, `gte`, and
     `exists`, while the proposed contract also includes set containment. The
     RDR should say whether set containment is part of the initial user outcome
     or an implementation-time MVV requirement not yet covered by the Resolve
     fixture.

2. **The user-facing pass/fail outcome for finite-domain lint is too abstract.**
   - Passage: `Approach`, `Technical Design`, `Validation`.
   - Concern: the RDR says lint proves coverage and overlap, but a reader cannot
     tell which finite-domain cases must be accepted, rejected, or downgraded to
     know the user outcome is delivered.

3. **The exact shape authors write is only illustrative.**
   - Passage: `Technical Design`, `Illustrative Code`.
   - Concern: RDR 0002 owns the TOML container, but this RDR owns the embedded
     atom grammar. The draft names conceptual fields, then shows one TOML shape;
     it should make the atom spelling/operator-kind contract explicit enough
     that a flow author knows what is legal.
