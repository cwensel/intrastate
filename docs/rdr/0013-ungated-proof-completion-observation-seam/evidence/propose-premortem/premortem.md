# RDR 0013 — Stage 2 premortem (variant: paragraph)

Date: 2026-08-28. Profile: mid.

Assume the seam shipped and failed in production. Eighteen months on, the
proof representation was swapped for a bitset engine whose cost tracks
dimension count, and the C2 differential kept passing — because
"completes" still meant "the retired cross-product enumerator exhausts",
so the published bound's verdict silently stopped describing the
machinery that actually proves, and nobody noticed until two conforming
implementations disagreed on a model. Separately, a `domainSize`
off-by-one moved the cap and the enumerator's ranges together, the
differential agreed with itself, and the defect surfaced only when a
user's model was refused at a cardinality the docs called admitted.
And the definitional close of the budget question (C3) was read a year
later as license to lower the bound to 64 without anyone treating it as a
product regression, since no test fires.

Does the recommendation survive? Yes, with one hardening. The
representation-swap failure is real but is a spec event by 0003's own
terms — a new representation publishes its own bound — so the answer is
to make the representation-dependence explicit rather than to switch
approaches: one line folds into Failure Modes naming a representation
swap as an event that reopens C2 (the observation is a claim about the
enumerating representation, not about proving in general). The
common-mode `domainSize` failure was already named in Risks and
Mitigations and is bounded by the distinct code paths the two halves run
(multiplication vs actual loop iteration with dedup); the alternative
that fully removes it — a dual independent proof representation — stays
rejected as over-scoped for a Low-priority spec-integrity risk. The
bound-lowering read of C3 is the recorded intent, not a failure: REQ-86
makes the value an implementation's published scope, and treating a
smaller scope as a product regression is a product-priority question no
test oracle admissible under the REQ-86/ADV-2 exclusions could decide
anyway. Verdict: hardened (one Failure Modes line folded).
