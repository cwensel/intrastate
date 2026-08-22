Model: claude-opus-5[1m]

# Existence-atom projection — derivation against the 0003 spike fixture

Date: 2026-08-21
Question routed here by RDR 0007 A12 (DOWNGRADED) / JDR 0001 §JD-4:
how does an existence atom project onto the declared-domain product that
coverage and overlap are proved over?

## Why the row-group branch of A12 cannot be the answer

A12 offers a disjunction: "either row-group scoping keeps the rows out of one
compared product, or absence has standing in the declared-domain product."

RDR 0003 scopes a row group by SELECTION CONTEXT — "one
source-state/recognized-outcome selection group" (Technical Design). RDR 0007's
two-row absence pattern (test row 9) is
`X exists = false` row  ‖  `X exists = true` + `X eq v` row.
Both rows answer the SAME source state and outcome; that is what makes them a
disjunction rather than two unrelated edges. So they necessarily land in ONE
row group. The first branch is unavailable, and the projection branch is forced.

## Why the presence axis is derived, not chosen

RDR 0007 fixes `exists` as the sole TOTAL operator: "it MUST decide true or
false from presence or absence of its referenced key alone. Its verdict is
`presence == literal`." A total predicate partitions its domain — it never
abstains.

Lint's proof requires every atom to denote a SUBSET of the scoped product
(RDR 0003 A2: "A row with `all` atoms denotes the intersection of each atom's
allowed subset of the scoped product"). An `exists` atom denotes
`{present}` or `{absent}`. Those are subsets of a presence axis. If the product
carries no presence axis, an `exists` atom denotes nothing in it, and A2's
coverage identity `union(row_i accepted) == scoped product` is evaluated with
one row's constraint silently dropped — the union is computed over an
under-constrained space.

So the presence axis is not a design preference; it is what makes an `exists`
atom denotable at all. Absence has standing in the product because a total
operator ranges over it.

## The fixture makes the gap concrete

`evidence/spikes/guard-fixture.toml` declares:

```toml
[tags.cluster_eligible]
provenance = "observed"
kind = "bool"
domain = [true, false]
```

and rule `foundational-to-cove` guards it:

```toml
[rule.guard.all.cluster_eligible]
exists = true
```

Under RDR 0003 as written today, `cluster_eligible`'s declared domain is
`{true, false}` — COMPLETE. There is no element for `exists = true` to select
and none for `exists = false` to select. The atom is unprojectable: lint either
drops it (unsound — it would certify a group exhaustive that refuses at
runtime, breaching §JD-4/P5) or refuses to prove anything about any group
containing an `exists` atom (sound but useless — it disqualifies the RDR's own
representative fixture).

Note the fixture ALSO shows why the axis must be per-key and declaration-driven:
`status`, `profile`, `stage`, `lens`, `prelock_iterations`, `rewind_target` are
all matched or written on every path and are never absent. Adjoining a presence
axis to all seven keys multiplies the product by 2^7 = 128 for no proof value.

## Two candidate semantics, tested against the fixture

**(P1) Presence axis adjoined per possibly-absent key.**
A key declared as possibly-absent contributes `{present, absent}` to the scoped
product. `exists = true` selects `{present}`; `exists = false` selects
`{absent}`; a value atom over that key selects a subset of `present × D`, and
is unevaluable on `absent` — which under §JD-4 withholds the exhaustiveness
claim exactly when RDR 0007's veto would refuse. The two-row pattern's rows are
then provably disjoint (`{absent}` vs `{present}`) and provably cover the axis
— which is the disjointness RDR 0007 A5 already asserts on the runtime side
("no assignment has X both absent and present").

**(P2) Rows carrying existence atoms are excluded from the product proof.**
Lint refuses to certify any row group containing an `exists` atom. Sound, and
consistent with the existing refuse-or-downgrade valve. But it disqualifies
`foundational-to-cove` — a row from this RDR's own representative fixture — and
it makes RDR 0007's sanctioned two-row absence pattern permanently unprovable,
pushing authors back to sentinel-stamping, which RDR 0007 names as THE
anti-pattern this whole seam exists to prevent.

P1 is the only option that keeps the sanctioned authoring pattern provable.
P1 also reproduces P2's safety as a special case: where a key's optionality is
NOT declared, there is no presence axis, the value atom is unevaluable on
absence, and §JD-4 withholds the claim.

## Blocked on one RDR 0002 declaration

P1 needs a per-tag optionality declaration. RDR 0002's tag declaration carries
"tag name, provenance (`owned`, `observed`, `recognized`), value kind, and
optional accessor reference" — RDR 0007 A12's evidence says it plainly: "RDR
0002's tag-declaration schema has no optionality field."

This is the same shape as the request RDR 0007 Phase 4 already hands to RDR
0003 for set-valued element encoding: a declaration RDR 0003's semantics need
and RDR 0002 must carry. It is a producer request, not a contradiction.

## Verdict

The SUBSTANCE is derivable and effectively settled by RDR 0007's totality
clause plus RDR 0003's own product definition; what is genuinely open is
(a) recording the rule in RDR 0003 and (b) the RDR 0002 optionality field it
depends on. That is a Propose-shaped edit (a new normative clause plus a
producer request), not a Resolve verification.

## Prior-art pass — clean negative

Bounded corpus pass (4 queries, budget exhausted) over `PapersFast`,
`StateMachineLit`, `StateMachineRes`:

```
decision table completeness and consistency verification missing attribute values
guard exhaustiveness proof optional undefined attribute three-valued
rule base verification completeness redundancy conflict detection condition attribute domain partition
exhaustive mutually exclusive guards absent optional attribute null coverage proof
```

negative: no corpus evidence in PapersFast, StateMachineLit, StateMachineRes

The negative is clean rather than weak-signal: PapersFast hits clustered at
0.73–0.81 on bibliography/index fragments (the signature of no on-topic chunk);
StateMachineLit's best was 0.46 on undefined-behaviour-as-termination in a
refinement framework; StateMachineRes returned only repo-local markdown. Two
VLDB property-graph normalization hits ("missing properties should not affect
the validity of a business rule") were REJECTED as citations — that is a
functional-dependency semantics choice, not a coverage/overlap proof treatment.

These corpora do not carry decision-table verification, DMN, or rule-base
completeness literature. No Prior Art citation is defensible for this
assumption, so the projection rule must rest on the internal derivation above
(RDR 0007's totality clause + RDR 0003's own product definition) and be
recorded as a Design Decision, never as Prior Art.
