# Post-Mortem: RDR 0021 Lint's Normalized-Graph Export

Started at a Stage 6 route-back (rdr-common §punt-ledger), before
implementation. Only the Escaped-Defect Ledger is populated; the remaining
post-mortem sections are authored after implementation.

## Escaped-Defect Ledger

One row per finding in the implementation's `<art>/triage.md`, enriched with
the two fields triage cannot assign. Route-back rows (rdr-common
§punt-ledger) append here at the moment a stage reopens a completed stage —
odc columns n/a until implementation triage.

| Finding | odc-type | odc-trigger | Expected-catching stage | Precursor lens | Escape distance |
| --- | --- | --- | --- | --- | --- |
| Stage 6 found A9's unforgeability limb refuted against the loader. C2 asserts the `<opaque>` abstraction sentinel "is NOT an authored value, and a consumer distinguishes it by that exact spelling", but only `<clear>` is reserved at the three authored-value ingress sites; nothing refuses an authored `<opaque>`, so it reaches `Node.Values` byte-identical to the synthesized sentinel. Because C2 also fixes consumer terminal-marking to the opaque-admits reading, a forged value marks a node terminal that lint reports as a dead end. The critique lens raised A9 and the repeatability lens widened it — both treated the collision question as owed verification and left it Pending, but neither ran the loader search that would have refuted it, and the draft continued to build C2's marking rule on top of the unverified limb. | n/a | n/a | `4-resolve` (A9's own Verification step named the declaration/atom validation search; resolve owed that search before the clause it licenses was written into C2) | `critique` (raised A9), `repeatability` (widened it into C2's terminal predicate) | 2 stages (caught at 6-reconcile; owed at 4-resolve) |
| Stage 6 found `0021:C2` restating a merged-node terminal-satisfaction quantifier that JDR 0001 §JD-23 homes, and picking the opposite reading. §JD-23 states the doctrine binds both liveness-shaped invariants, that "neither record widens or narrows the split unilaterally", and that both records cite the entry rather than restate the other's contract; its siblings are 0015 and 0022. C2 instead writes its own existential/opaque-admits rule in 0021's own contract block, justified on 0021's over-approximation soundness rule. Whether the published merged relation is genuinely a different object from the split-node check is exactly the question a joint-decision home exists to settle, so the clause was authored unilaterally where the registry required a hoist-and-cite. The joint-decision check recorded `joint_checks=1` with `joint_check_home=clear`, which did not reach §JD-23 because 0021 is not in JDR 0001's cluster frontmatter. | n/a | n/a | `2-propose` (the joint-decision check owed §JD-23 the moment C2 adopted a merged-node terminal quantifier; the entry predates this RDR's draft) | `3amigo` / `critique` (both reviewed C2's marking clause without testing it against §JD-23) | 4 stages (caught at 6-reconcile; owed at 2-propose) |
