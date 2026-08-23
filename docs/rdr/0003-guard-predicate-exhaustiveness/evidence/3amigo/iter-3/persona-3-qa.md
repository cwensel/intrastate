Model: claude-opus-5[1m]

# Persona 3 — QA / Tester — RDR 0003, re-entry iteration 3

Delta-scoped to the passages the re-entry disturbed: the rehomed tag declaration
model (value kind, finite domain, optionality, single-valuedness, element
universe), the RDR 0002 split, A8/§JD-4, §JD-13's single-valued marker, §JD-14's
escape-row confirmation, A5's re-verification, and the still-open A10/A12/A14/A15.

Severity-ranked. Each finding names the test it prevents and the missing
pass/fail criterion.

---

## P3-1 — HIGH — `Normative Contracts`, optionality clause (line 704-712) and single-valued clause (line 714-725): "conforming evaluation view" has no pass/fail definition and no violation outcome

Both clauses the re-entry added or touched are stated as constraints on a
**"conforming evaluation view"**:

- optionality: "A key declared always-present MUST NOT be absent from a
  conforming evaluation view"
- single-valued: "at most one of its declared domain values holds in any
  conforming evaluation view"

`conforming evaluation view` occurs exactly twice in the whole cluster — these
two lines — and nowhere else in this RDR, RDR 0002, RDR 0006, RDR 0007, or JDR
0001 (grep confirms zero other occurrences). Nothing states **who checks
conformance, when, or what happens when a view violates it.** The `disposition`
table (line 944-961) enumerates input classes × outcomes and has no row for
"view violates a declaration"; the trace table (line 963-985) never walks one.

**Test I cannot write:** a runtime/kernel test that supplies an evaluation view
in which a tag declared always-present is absent (or in which a single-valued
tag holds two of its declared domain values), and asserts the system's response.

**Missing pass/fail criterion:** whether such a view is (a) rejected with a
named diagnostic before evaluation, (b) evaluated anyway with declaration-driven
lint results now unsound, or (c) undefined behaviour outside the contract. The
clause says MUST NOT without saying MUST-NOT-**or-else**, so there is no verdict
to assert. Note this is not academic: the whole point of the optionality clause
is that lint decides `exists` projection and the "can refuse" narrowing *from
the declaration and never from a runtime trace* — which makes lint's soundness
strictly conditional on view conformance being enforced somewhere. Nothing
enforces it.

---

## P3-2 — HIGH — `Testing Strategy` Scenario 4 (line 1563-1575) vs. §JD-13 duty: the single-valued marker has a rejection test but no acceptance test, so `graph-single-valued-state`'s producer is untested

§JD-13 (JDR 0001, "Decided 2026-08-22") homed the single-valued field here
*because* RDR 0006 gates a **mandatory blocking** code on it
(`graph-single-valued-state`, `docs/rdr/0006-...:318`) and asserts it in RDR
0006's own MVV illegal-fixture matrix (`:643`). The re-entry discharged this by
adding the declaration clause plus one **negative** test: Scenario 4 asserts the
malformed cases (marker on a `set` kind, marker on a kind with no finite
domain).

There is no positive case anywhere. Grepping `single-valued` across this RDR
(lines 15, 17, 543, 687, 715, 729, 937, 940, 1405-1411, 1455-1456, 1568, 1573)
yields: metadata, the model enumeration, the clause, the domain/kind rejection,
two authority rows, the prerequisite checkbox, one Authorability sentence
("a single-valued enum tag declares the marker"), and Scenario 4's rejection
cases. The Authorability paragraph asserts the field is *declarable*; no
Scenario asserts anything is *decided differently* because it was declared.

**Test I cannot write:** the acceptance-side test for the field §JD-13 homed
here — take one row group over a finite-domain enum tag, lint it once with the
tag declared single-valued and once without, and assert the two runs differ in
the way the clause promises.

**Missing pass/fail criterion:** the clause says the marker "licenses a consumer
to treat the tag's domain as a partition — mutually exclusive values whose
coverage the scoped product can be grouped over — rather than as independent
dimensions." That is the only statement of effect, and it is not operational:
it does not say what the scoped product *is* under each reading, so I cannot
compute the two expected verdicts and assert they differ. Concretely — for an
`enum` tag with domain `{a,b,c}`, is the non-single-valued product 3
assignments, or 2^3 = 8 (each value independently held)? The coverage identity
`union(row_i accepted) == scoped product` (A2, line 185) gives a different
verdict under each, and nothing in this RDR picks one. Without that, the
producer field ships with only a rejection test behind it, and RDR 0006's
mandatory blocking code has an untested input.

---

## P3-3 — HIGH — `Normative Contracts` domain/kind agreement clause (line 727-735), second sentence: the literal-outside-declared-domain rejection has no test in any Testing Strategy scenario

The clause the re-entry added carries two independent rules. The first
(domain/kind disagreement) is tested — Scenario 4 (line 1563-1575) enumerates
the three malformed-declaration cases explicitly. The second is not:

> "A value appearing in a guard literal that lies outside its tag's declared
> domain MUST likewise be rejected before resolution, so an unsatisfiable atom
> is a load-time error rather than a silently-never-matching row."

Scenario 4's stated malformed-guard-atom set is "unknown tag, unknown operator,
unsupported operator/tag-kind pair, and literal parse mismatch" — an
out-of-domain literal is none of these. `profile eq "gigantic"` where `profile`
declares `{foundational,mid,large,huge}` parses cleanly, matches its kind, and
uses a known operator; it fails only the domain check. No scenario covers it.

**Test I cannot write:** the load-time rejection test for a well-typed,
well-parsed guard literal outside its tag's declared domain — the exact case
the clause exists to catch.

**Missing pass/fail criterion:** which predicate semantic kind this rejection
carries. A4 (line 195-201) enumerates the kinds this RDR owns as "unknown
operator, type mismatch, literal parse failure, and unsupported operator/tag-kind
pairing" — the out-of-domain literal is not in that list either, so there is no
named code to assert on. The rule is normative, load-bearing (an unsatisfiable
atom silently narrows a row's accepted assignments to ∅ and can make a group
look non-exhaustive for the wrong reason), and untestable as written.

---

## P3-4 — MEDIUM — `Normative Contracts` tag declaration model clause (line 684-692) and the operator/kind matrix (line 571-577): value-kind vocabularies disagree, so operator/kind rejection tests have two contradictory oracles

The re-entry's declaration clause fixes the kinds as: "`enum`, `bool`, `int`,
`set`, and opaque scalar" (line 689). The operator/kind matrix rows for `eq` and
`in` accept "enum, boolean, integer, string-like scalar" (lines 573-574). Same
five concepts, two spellings, and `opaque scalar` vs `string-like scalar` is not
obviously the same thing — one names representation, the other names opacity.
The `exists` row (line 576) uses a third vocabulary: "optional scalar or
optional set-valued tag", which is a kind crossed with the optionality field
rather than a kind.

**Test I cannot write:** the table-driven operator×kind acceptance/rejection
matrix — the natural shape for Scenario 4's "unsupported operator/tag-kind pair"
case. I cannot enumerate the rows, because I cannot tell whether the kind axis
has four members or five, nor whether `opaque scalar` and `string-like scalar`
are one cell or two.

**Missing pass/fail criterion:** one canonical kind enumeration that both the
declaration clause and the matrix use, and — for `exists` specifically — whether
its acceptance condition is a *kind* constraint or an *optionality* constraint.
The latter matters for a real assert: is `exists` over a key declared
always-present a rejection (the matrix says "optional scalar"), or is it legal
and merely vacuous (the presence-dimension clause at line 785-797 says an
always-present key "contributes no presence dimension", which reads as legal)?
Those are opposite verdicts for the same fixture row, and the MVV's negative
control (line 1446-1449) depends on always-present keys behaving predictably.

---

## P3-5 — MEDIUM — `Minimum Viable Validation` Authorability paragraph (line 1451-1462) and Phase-gating table row 3 (line 1471): the MVV is declared authorable while A14 is still Pending, and no test states what happens if A14 lands negative

The re-entry's central claim is that the rehoming "closes the gate that
previously bound the MVV's schedule to a peer's producer request", leaving "only
sequencing, not authorability". But A14 (line 426-454) is still `Pending`:
per-atom `block` retention through normalization is an open request on RDR 0002,
and the Phase-gating table lists Phase 3 (Target-Flow Fixture) as "Blocked on:
A14" while simultaneously marking it "Startable today — Yes".

A14's "If wrong" (line 452-454) is severe: "`unless` collapses into `all` after
normalization, the identity tuple loses a component that makes it total, and a
row's excluded intersection can no longer be subtracted from its accepted
assignments." Every MVV scenario that uses `unless` — Scenario 1 (mixed
`all`/`unless`), Scenario 2 (product-level gap and overlap), the trace table's
step 1 witness `profile-to-grounding` with `unless prelock_iterations gte 3` —
depends on it.

**Test I cannot write:** the normalization-carriage test A14's disposition names
as its own detector ("the failure is caught by the first normalization test").
No scenario in `Testing Strategy` is that test; Scenario 4 covers parse-level
malformation, not post-normalization field retention.

**Missing pass/fail criterion:** an assertion that a normalized atom authored in
`unless` still reports `Block == unless` after normalization. The document
asserts the failure is "detectable at the first normalization test" without any
scenario making it detectable, so the safety argument for shipping A14 Pending
rests on a test that does not exist in this RDR's Validation section.

---

## P3-6 — MEDIUM — `Testing Strategy` Scenario 3 (line 1550-1562), the A15 discharge: "two equal-cardinality products of differing shape" has no constructibility criterion, so the test cannot be written deterministically

A15 (line 455-489) is discharged by this scenario, and the re-entry restates
that in Metadata (line 21-22) and the prerequisite list (line 1379-1382). The
scenario says: "two equal-cardinality products of differing shape (few wide
dimensions vs. many narrow ones)" receive "the same verdict — a divergence
refutes A15."

Two problems for a test author:

1. **No cardinality is given relative to the bound.** The verdict being asserted
   is presumably "both refuse" (both over the bound) or "both prove" (both
   under). The scenario does not say which side of the published bound these
   products sit on, and the bound is "a single integer constant published by the
   implementation" (line 897-899) — i.e. unknown at fixture-authoring time. A
   fixture with a hardcoded cardinality either passes vacuously (both under the
   bound, both prove, telling us nothing about the shape sensitivity A15 is
   about) or is unbuildable.

2. **"Differing shape" is not quantified.** How few is "few wide", how many is
   "many narrow"? Any two distinct factorizations of the same integer satisfy
   the prose. If the divergence is shape-sensitive at some threshold, an
   arbitrary pair may straddle it or may not.

**Test I cannot write:** the A15-discharging comparison itself — the one test
this RDR nominates to close its last self-owned open assumption.

**Missing pass/fail criterion:** the target cardinality expressed relative to
the published bound (e.g. "both products at bound−1" and "both at bound+1", run
as two cases), and the two factorizations pinned concretely (e.g. `2×2×…` ×n
vs. `2^n × 1`). Without these the scenario is a description of an experiment,
not a test with a verdict.

---

## P3-7 — LOW — `authority` census row "Tag declaration model" (line 937) vs. RDR 0002's authoring clause: the single-valued field is absent from the peer's carriage contract, so no end-to-end declaration test can be written for it

The authority row and the `Capability Dependencies` row "Authoring location +
normalization carriage for tag declarations" (line 1060) both place authoring
and carriage with RDR 0002. RDR 0002's normative clause
(`docs/rdr/0002-transition-table-as-reviewable-data.md:341-350`) enumerates the
carried type model as "value kind, and optionally a finite domain, an
optionality marker, and a set-element universe" — **four fields, no
single-valued marker**. Its schema list (`:218-219`) repeats the same four, and
its ownership-correction note (`:909-910`) repeats them again. RDR 0006's lint
input contract (`docs/rdr/0006-...:264-267`) does name "single-valued grouping
when applicable" separately from the four.

The re-entry added the fifth field here on 2026-08-22; RDR 0002's clause was
last written 2026-08-21 and lists four. RDR 0002's clause does say "carry every
declared field through to the normalized model without loss", which arguably
covers it generically — but its own enumeration contradicts that reading, and
this RDR's own §JD-13 checkbox (line 1403-1414) marks the item **Done** without
noting the peer enumeration still reads four.

**Test I cannot write:** the end-to-end declaration test — author
`[tags.status] single_valued = true` in a TOML fixture, normalize, and assert
RDR 0006's lint input carries the marker. There is no authoring location for the
field to be written at, because RDR 0002 does not enumerate it under
`[tags.<tag>]`.

**Missing pass/fail criterion:** the authored spelling and placement of the
single-valued marker. This is the same unhomed-producer shape §JD-13 was opened
to close, displaced one seam outward: the field is now *declared* here but not
*authorable* there. Lower severity than P3-2 because a single-field addition at
RDR 0002's refine closes it and A14 already establishes that venue is open — but
it means the §JD-13 checkbox is marked Done while the field remains unwritable
in a fixture.

---

## Not findings

- The escape-row clause (line 799-807) confirmed by §JD-14 is testable as
  written: both halves (participation in the union, non-exclusion from overlap)
  have concrete verdicts, and Scenario 2's overlap assertion reaches the second
  half. The repair owed is RDR 0006's, correctly.
- A8's narrowing clause (line 809-814) and the withheld-claim observability
  clause (line 816-823) are testable: MVV Scenario 6 asserts positively, names
  the artifact, and explicitly fails a run that emits nothing. The atom-naming
  half's missing producer is correctly booked as RDR 0006's duty, and the RDR
  says so.
- A5's re-verified anchors all resolve — `internal/resolve/resolve.go:173-174`
  carries `RuleID`/`SourceLocator` with the cited comment, and `Row.Guard string`
  at `:185` is the pre-reshape form the RDR says it is.
