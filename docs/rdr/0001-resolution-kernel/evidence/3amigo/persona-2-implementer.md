# Persona 2 - Implementer

Question: if I started coding this Monday, what would I ask in the first hour?

Findings:

1. In `Context / Technical Environment`, `Approach`, `Technical Design`, `Normative Contracts`, and `Illustrative Code`, should the resolver package execute injected accessors, or should an accessor executor produce the owned snapshot before calling the resolver? The RDR currently says both "rely on injected accessors" and "owned tag snapshot read from those artifacts", while RDR 0004 says the resolver receives accessor-read values.

2. In `Validation / Testing Strategy`, what concrete refusal taxonomy must the kernel expose for `no match`, `ambiguous match`, `missing owned state`, `unevaluable guard`, and `unmodeled recognized outcome`? The RDR says typed refusal, but the implementer needs the unmodeled-outcome case named in the validation list to avoid collapsing it into no-match.
