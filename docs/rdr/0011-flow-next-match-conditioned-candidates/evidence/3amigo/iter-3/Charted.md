# Charted to successor — 3amigo iter-3, RDR cli/0011

Findings that were real but net-new scope. Each is dismissed from this loop and
recorded here rather than absorbed into the RDR (rdr-common: never edit the
current RDR to absorb net-new scope).

## 1. No diagnostic surface for a match-excluded row

**Origin**: L4 (PM-H1) / L5 (PM-H2). **Why out of scope**: the RDR's contract is
the candidate PREDICATE. A caller who wants to know *why* a row was excluded
needs either an `excluded_by` field on a reported non-candidate or a
model-inspection verb, and both are payload/surface additions this record
deliberately does not take — C1's shape is "an excluded row is one the state
cannot take", and reporting non-candidates inverts it. What this loop owed and
paid: correcting F1's remedy (it named a `dump` verb that does not exist and a
`--all` mode that C2 strips of match facts) and correcting BR3's rejection
ground (which claimed `--all` yields information C2 forbids it to yield). The
cost is now stated honestly in both places instead of papered over.

**Suggested successor**: a `--why`-style opt-in carrying `excluded_by` (the
failing atom) on match-excluded rows, or a model-inspection verb. BR3 now names
this as the deliberate reopening path. Note the shipped `flow_surface_0005_test.go:105`
asserts `dump` is FOREIGN to the flow group, so a model-inspection verb is a
surface decision against 0005's four-verb closure, not a small addition.

## 2. `absent` does not name which of its two remedies applies

**Origin**: L19 (PM-L1). **Why out of scope**: distinguishing "bind a reader"
from "supply a `--tag`" requires the entry to carry the key's `provenance`,
which is a payload widening on every `unknown` entry. The RDR now states the
limit in Consequences rather than implying the reason vocabulary carries the
whole remedy.

**Suggested successor**: a `provenance` member on the `unknown` entry, or a
documented convention that the caller reads `[tags.<key>]` from the model. Worth
folding into whatever successor takes item 1, since both are `unknown`-payload
questions.

## 3. Text-mode rendering of `unknown` is the shared gateway's, not this verb's

**Origin**: L7 (Impl-S3 / QA-3). **Why out of scope**: `D-undecided-reporting-shape`
asked for `key (reason)` in text mode. `flow next` declares no `FindingCarrier`,
so `respond.OK` renders through the generic `writeTextPayload`/`flatten`, whose
own doc pins the two-modes-cannot-disagree invariant to `0005`'s
REQ-7/REQ-9/REQ-11/REQ-120 — clauses this RDR does not override. Minting a
per-verb text template for cosmetics would edit a predecessor's surface. The
decision now says the pairs are carried in text mode (both members, so the
diagnosis survives) without claiming a bespoke layout.

**Suggested successor**: if path-qualified leaf lines prove unreadable for list-
of-struct payloads generally, that is a `respond` gateway RDR covering every
verb — not a `flow next` clause.
