Model: claude-opus-5[1m]

# 3amigo — Persona 1 (Product Manager), RDR 0010

Read: `0010:§problem-statement`, `0010:§approach`, `0010:§decision-rationale`,
`0010:MVV`. Widened where noted.

### P1-1 — The outcome the Problem Statement names as the reason to build this is guaranteed by no contract  [High]
**Anchor**: `0010:§problem-statement`, `0010:C5`, `0010:MVV`
**Finding**: The Problem Statement makes coverage the *justifying* outcome —
"Exhaustiveness/coverage lint over observed dimensions must keep working
unchanged — the PoC surfaced a 208/512 coverage hole ... which is the property
that makes stateless tables worth supporting at all." But the delivered
guarantee is conditional on an authoring convention that this RDR states only
in prose. `0010:C5` requires coverage to "run unchanged over the authored
observed dimensions"; over a decision table whose dimensions are authored as
`[rule.match.<key>]` atoms, `guard.Dimensions` returns none and
`internal/graphlint/coverage.go:35` takes the `len(dims) == 0` branch, closing
the claim by row membership alone — verified on main. Such a table lints exit 0
with an empty findings list while incomplete. So the user who writes the model
"as what it is" and discriminates with the more obvious `match` shape gets a
green lint over an incomplete table, which is the exact failure the Problem
Statement says the feature exists to prevent, delivered *by* the feature.
The RDR knows this: the risk is written out at `0010:§risks-and-mitigations`,
`0010:§key-discoveries`, and `0010:§illustrative-code`. But every mitigation
named there is documentation, an illustration, or a test of the RDR's own
fixture (`0010:MVV` step 1, `0010:S4`) — none of them is a fence an author's
model has to pass. `0010:BR6` separately rejects a class-announcing finding as
noise; the rejection is sound for *announcing the class*, and does not reach
the different question of announcing a decision-table group with zero
participating dimensions.
**Blocks**: whether this RDR owes a sixth contract (or an amendment to `0010:C5`)
making a `"decision-table"` group with zero participating guard dimensions a
finding rather than a silent pass — and, if not, whether the RDR may claim the
coverage outcome at all or must state the outcome as conditional on the guard
authoring shape. Locking `0010:C5` as written locks the conditional guarantee.

### P1-2 — Phase 4 does not commit to the one document the outcome depends on  [Medium]
**Anchor**: `0010:§phase-4-model-and-docs` (widened), `0010:§risks-and-mitigations`
**Finding**: If P1-1 stands as documentation-only, the docs are load-bearing for
the user outcome, and Phase 4's intent names only "`docs/cli-output-contract.md`
and the model authoring docs updated." That is a generic docs line; it does not
say the authoring doc must carry the guard-vs-match warning, and it is the sole
delivery vehicle the mitigation names ("the Illustrative Code authors guard
atoms and says why" is inside this RDR, which the author never reads). `0010:MVV`
and `0010:§testing-strategy` verify the *implementation*; nothing verifies the
author-facing artifact that the guarantee rests on.
**Blocks**: whether Phase 4 is done — i.e. what the implementer must write for the
phase to be complete, and whether a docs assertion belongs in `0010:MVV`.

### P1-3 — The Problem Statement's user is a specific consumer; the RDR never states what that consumer can do afterward that it cannot do now  [Low]
**Anchor**: `0010:§problem-statement`, `0010:§approach`, `0010:MVV`
**Finding**: The motivating user is rdr#tmxk's navigator model (`0010:§background`).
The Problem Statement describes the workaround's cost (dummy tag, `[initial]`,
accessors, scratch artifact) and the design question, but never states the
consumer-side acceptance: what rdr#tmxk's model looks like after this lands and
what call it makes. `0010:MVV`'s end-state is written entirely in intrastate's
own terms ("a stateless model lints for coverage and resolves to `rule` + `emit`"),
over a synthetic fixture, with `models/rdr.toml` explicitly untouched
(`0010:§phase-4-model-and-docs`). That is correct for the generic-engine
constraint, but it means nothing in the RDR closes the loop back to the user
whose refusal opened it: no step confirms the motivating model becomes
authorable, and `0010:A6` ("a flat, string-valued `[rule.emit]` table is
sufficient for the motivating consumer's answer") is the only place the
consumer's need is characterized at all.
**Blocks**: what counts as "shipped" for intrastate#zdat — whether the tracker
closes on the generic fixture passing, or on rdr#tmxk's model being rewritten
without the scratch artifact. Also bears on `0010:A6`'s sufficiency claim, which
has no consumer-side check behind it.

## Widening

Sent out of the owned set twice. P1-1 could not be settled inside
`§problem-statement` — the outcome is asserted there but the guarantee lives in
`0010:C5` and in `internal/graphlint/coverage.go`, so I read the contract and the
source to see whether the prose claim holds unconditionally (it does not).
P1-2 followed the mitigation named in `§risks-and-mitigations` to its only
delivery vehicle, `§phase-4-model-and-docs`, since a silence about docs has no
line range in the owned set.
