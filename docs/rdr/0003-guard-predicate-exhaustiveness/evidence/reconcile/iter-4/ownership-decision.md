Model: claude-opus-5[1m]

# Ownership decision — the tag declaration model (RDR 0003, iteration 4)

## Decision

The **tag declaration model** — value kind, finite domain, per-tag optionality,
set-element universe — is homed in **RDR 0003**, not RDR 0002.

RDR 0002 keeps tag *provenance* (its consumers RDR 0004 and RDR 0006 read it),
the authoring location under `[tags.<tag>]`, and normalization carriage. It
cites RDR 0003 for what the type fields mean.

**RDR 0006 is demoted `Final` → `Draft`** to discharge the three peer agreements
RDR 0003 cannot close against a locked document (A8, A10, A12).

## Why the question was reopened

Stage 6 (iteration 3) blocked lock on A11/A7/A9 and routed the fix as a producer
request to RDR 0002. That routing inherited the cluster's habit — every peer
cites "RDR 0002 declares tag name, provenance, value kind, and an optional
accessor reference" — without testing whether RDR 0002 actually owns it.

It does not. Four findings, each independently sufficient:

1. **RDR 0002's only normative tag clause requires provenance alone.** Of its 13
   `normative` blocks, exactly one concerns tags: "The model MUST declare every
   tag it matches or writes, including each tag's provenance: owned, observed,
   or recognized." Nothing about kind, domain, or optionality.
2. **`value kind` is normative nowhere in the cluster.** It appears in RDR
   0002's *prose* schema list (part 2) and in a spike fixture. Every peer citing
   the four-field declaration is quoting that prose as if it were a contract —
   RDR 0003 (5 sites), RDR 0006, RDR 0008. The tag type model was never homed.
3. **RDR 0002's normative validation vocabulary has no type axis.** Its
   categories are malformed TOML, unknown schema field, missing outcome
   alphabet, unknown tag, unknown context, write to non-owned tag, unknown
   accessor, unsupported version, malformed escape, ambiguous overlap — all
   structural. There is no "value outside declared domain" or "domain/kind
   mismatch" category, and adding one imports a type system into a document
   whose charter is "the on-disk sparse representation and the normalized
   expanded-table view".
4. **RDR 0007 (`Final`) already records the opposite assignment.** In its
   rejected-alternative analysis: "Typed literals and value kinds (bounded
   integers, set universes) are **RDR 0003's declarations**; the kernel has none
   and would have to grow a declaration model."

Supporting signal: the token `finite` appears **zero** times in RDR 0002 and 58
times in RDR 0003.

## The argument that was rejected

A fresh-context adjudication argued for RDR 0002 on the grounds that "kind
already lives in 0002 and is read by 0003's operator/kind matrix with no
ownership dispute, so domain belongs beside it — kind is `integer`, domain is
`{0..3}`; they are the same class of fact at different resolutions."

The premise is false, and its falsity reverses the conclusion. Kind lives in RDR
0002's *prose*, is normative nowhere, and a `Final` peer assigns it to RDR 0003.
So the precedent does not show a settled home for kind that domain should join;
it shows the tag type model was never homed at all, and the domain field is the
same gap one layer deeper. The correct repair is to home the whole model once,
with the algebra it feeds.

The adjudication's cohesion point survives and is honored: RDR 0003 must not
specify TOML key placement, dump ordering, or normalization mechanics. It does
not — those stay RDR 0002's, and the new RDR 0003 clause says so explicitly.

## Why not a new RDR

It would split one tag's declaration across two documents (provenance in 0002,
type model in the new one) and add an edge from each of 0002/0003/0006/0007.
JDR 0001 §D1 faced the same choice — absorb into an existing owner, or mint a
new RDR superseding 0001's row shape — and chose absorption. Same reasoning
here.

## Precedent applied

JDR 0001 P6 ("one home per contract; cite, never restate") and the §D1/§D4
landing pattern: the home is the document that owns the *surrounding contract*,
not the document that feels the pain. RDR 0007 owns the atom shape while RDR
0002's normalizer produces it, and nobody reads that as incoherent — ownership
of a contract and ownership of the parser are already separate throughout this
cluster. The declaration model is the guard algebra's input alphabet, so it
belongs with the algebra.

P7 ("pre-release, prefer the clean shape") also applies: `internal/resolve` has
**zero** non-test consumers, so there is no migration cost to putting the
contract where it belongs.

## What did NOT move

**A14 — per-atom `block` retention** stays a request on RDR 0002. Retention
through normalization is *carriage*, which is RDR 0002's charter proper, not
declaration *semantics*. Keeping this distinction is what makes the rehoming a
boundary correction rather than a land grab.

**Provenance** stays in RDR 0002. It is the one tag field that document
normatively owns, and its consumers are RDR 0004 and RDR 0006, not the guard
algebra.

## Consequences

| Record | Before | After |
| --- | --- | --- |
| A11 finite-domain field | Pending — BLOCKER, no producer | **Verified** — Design Decision, stated here |
| A7 existence-atom projection | Pending — BLOCKER (A11 subfield) | **Verified** — normative clause on the recorded derivation |
| A9 set-element universe | Pending — BLOCKER (A11 subfield) | **Verified** — Design Decision |
| A2 finite-domain derivation | Pending — BLOCKER behind A11 | **Verified** — inputs now owned |
| A14 atom `block` retention | Pending (rode the four-field request) | Pending — the one genuine RDR 0002 request |
| A8/A10/A12 peer agreements | Blocked against `Final` RDR 0006 | Pending — scheduled edits on `Draft` RDR 0006 |
| MVV authorability gate | Cases unwritable | **Closed** — every case authorable |

Lock now turns on **A8 alone** (the §JD-4 recording assignment), down from five
blockers.

## Dependency direction

The rehoming **reverses** the 0003→0002 edge. RDR 0002 now cites RDR 0003 for
the declaration model, which matches the direction already stated in RDR 0002's
own Approach: "RDR 0003 owns the fixed predicate operators."

## §JD-4 recommended arm (recorded in the JDR)

RDR 0003 records the narrowing; RDR 0006 cites it and reuses
`graph-unprovable-coverage`.

- RDR 0006's clause fires on a *non-finite dimension*; RDR 0003's case is a
  *fully-finite product whose participating row can still refuse* — a different
  trigger RDR 0006's text does not reach.
- The clause must name the **refusing atom**. The token `atom` occurs once in
  RDR 0006, only to delegate atoms to RDR 0003.
- `graph-unprovable-coverage` already means "Required finite-domain proof
  unavailable", which covers a withheld claim; minting a second code would imply
  a non-blocking tier RDR 0003's own contracts forbid.
