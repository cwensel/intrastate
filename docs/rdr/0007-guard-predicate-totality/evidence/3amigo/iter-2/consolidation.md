Model: claude-opus-5[1m]

# 3amigo consolidation — RDR 0007 Guard Predicate Totality (iteration 2)

Delta pass. Re-entry scope from the `Status:` line — re-verify A3, A16, A17,
A18, A19, A21, A22 — plus the passages the JDR 0001 §D1/§D4 reshape touched.
Three personas run independently (PM / Implementer / QA); no persona saw
another's findings, and the isolation check found zero cross-persona references
in any file.

Persona lists: `persona-1-pm.md` (10), `persona-2-implementer.md` (11),
`persona-3-qa.md` (14). 35 findings total.

Overlap is computed from the **anchored passage each finding cites**, not by
re-judging the merged list. Overlap marks a hotspot *passage*; it does not
validate a finding, and a single-persona finding is not thereby weaker.

**Iteration 1 is closed and not carried forward.** Its headline hotspot H-1
(the opaque `Row.Guard string` and the unnamed flattening step) was dissolved
by §D1 putting parsed atoms on the row; H-2's "cannot name the absent tag" is
now the per-row/per-atom payload clause. The iteration-1 ledger is history, not
open scope.

## Hotspot passages (2+ personas)

### H-1 — The refusal payload is asserted on everywhere and typed nowhere (IMP-2 + QA-1 + QA-7 + QA-13)

**Passage**: Normative Contracts, refusal clause — "for every undecidable row,
its `(RuleID, SourceLocator)` and, for each of its unevaluable atoms, the
referenced key, the block, and a reason drawn from a closed set — `absent` …
or `uncomparable`"; Load-Bearing Decisions — "The absent-key payload is a field
on `Refusal`, not a new kind."

The clause fixes the payload's *content* and *sort order* but never names the
field, its element type, or whether `reason` and `block` are closed exported
constant sets (the house pattern being `resolve.go::RefusalKinds`). The only
concrete name in the evidence tree, the spike's `UndecidedAtoms`, is disowned
by A3's own status as the wrong shape — it sat *beside* a retained
`Refusal.Guard`, which the Normative Contracts replace.

Consequences across personas: the implementer cannot write the `Refusal` type
change (IMP-2); QA cannot assert "key in payload" for matrix rows 1, 12, 17,
20, MVV scenario 1, or `oracle` row 1 (QA-1); row 20 offers two candidate
oracles — `Refusal.Rows` versus payload row entries — with no stated invariant
between them (QA-7); and "MUST evaluate every atom of every survivor" has no
observable consequence beyond row 20 (QA-13).

**Priority**: highest. Six assertion sites depend on a name that does not
exist, and RDR 0005's renderer and Phase 3's exported contract test both type
against this field.

### H-2 — "Sorted by row identity then key" is not a total order (IMP-4 + QA-2)

**Passage**: Normative Contracts, refusal clause — "MUST sort the payload by
row identity then key, so the payload is a function of the input tuple (RDR
0001 REQ-1), never of atom or row order."

Both personas derived the same collision independently, from the RDR's own
sanctioned idioms: A5 requires a row carrying **two atoms over one key**
(`X exists = true` beside `X eq v`), Testing Strategy row 9 makes that row a
test subject, and the Failure Modes P-5 fix does the same inside `unless`
(`legal_hold exists = true` beside `legal_hold eq true`). A key can also be
unevaluable in both `all` and `unless` of one row. Two payload entries then
collide on `(row, key)` and differ only in block, operator, or reason.

This re-introduces one level down exactly the order-dependence ADV-3 exists to
forbid, and REQ-1 determinism is the clause's own stated purpose. No in-repo
comparator can be borrowed: `resolve.go::compareRefs` and
`resolve.go::missingOwned` are both total by construction.

**Priority**: highest. It is a self-contradiction against a REQ-1 claim, and it
makes row 17's permutation test unwritable.

### H-3 — A3's spike-shape gap versus the Testing Strategy's baseline claim (PM-7 + IMP-5 + QA-14)

**Passage**: A3 **Status** — "the BOUND is verified; the SHAPE the spike proved
is not the shape this RDR specifies"; against Testing Strategy preamble — "the
frozen 154-test suite passes unchanged under the reshape, so every row below is
ADDED surface, not a re-encoding."

All three personas hit the same contradiction from their own side. The spike
retained `Refusal.Guard` (retyped) and added `UndecidedAtoms` beside it; under
"replace," `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` must be
RE-DECIDED rather than re-encoded — which A3 names as its own "If wrong"
branch. The Testing Strategy nonetheless asserts the 154-PASS baseline as
established for the specified shape.

PM asks whether Phase 1 is a migration or a contract reopening (PM-7); the
implementer asks whether the re-spike gates the work and what replaces the
guard-text assertion at `fixup_test.go:64` (IMP-5); QA notes Phase 1 states no
exit count (QA-14).

**Priority**: high. It decides whether Phase 1 may start, and A3 is inside the
re-verify scope.

### H-4 — `RequiresOwned` narrowing has no enforcer, and its test cannot fail (PM-8 + IMP-10 + QA-3)

**Passage**: Normative Contracts — "On an escape row, whose `Writes` MUST be
empty (RDR 0009), RequiresOwned is therefore empty (A21). Who populates the
field is RDR 0002's (JDR 0001 §JD-3)"; Testing Strategy row 15.

A21 concedes §JD-3 "remains NOT closed" and that 0002 carries the producer
"not yet a clause." So the narrowing that kata `xg7p` exists to fix — and that
Overrides names as *fixed* — has no implementer. The Phase 1 change is
doc-only; row 15 passes identically under either meaning because
`missingOwned` over an empty slice never enters the loop (QA-3 calls it a
fixture tautology and notes `fixtures_test.go::escapeRow` and frozen Fixup-1e
assert the opposite); nothing discriminates the narrowed contract (IMP-10).

**Priority**: high. It decides whether `xg7p` can close on lock.

### H-5 — Phase 3's `contains` leg is blocked, and row 8 passes trivially meanwhile (PM-10 + IMP-9 + QA-9)

**Passage**: Phase 3 — "**Blocked on one RDR 0003 declaration.** `resolve.Tag.Value`
is a bare `string` … The `contains` leg of this contract test therefore cannot
be written from any current document, and neither can Testing Strategy scenario
8's present-key half."

The sole mitigation for the value-evaluator drift risk ships with a
`contains`-shaped hole, unblocked only by a Phase 4 "REQUEST" (PM-10). Row 8 as
written is listed as backed while its present-key control is declared
unwritable, so it passes trivially (QA-9). The implementer notes row 8 has no
present-key half at all, leaving it unclear whether a row is missing or the
cross-reference is stale (IMP-9).

**Priority**: medium-high — partly a stale cross-reference, partly a real
coverage hole.

### H-6 — The existence operator's literal half has no home (IMP-6 + QA-5)

**Passage**: Normative Contracts, existence clause — "The kernel recognizes the
existence operator by one operator token, and its literal by exactly two
boolean literal forms, all three exported as kernel constants that the
normalizer MUST emit"; Testing Strategy row 11 — "a foreign boolean literal is
rejected at load".

Three exported constants are specified with no names, types, or values — the
drift boundary's actual strings are the one thing not written down (IMP-6).
Worse, "rejected at load" names a load path that does not exist in
`internal/resolve`, and the kernel-side rule for a *malformed* existence
literal is absent; QA-5 notes the A3 spike's literal-less `OpExists` atom
decides FALSE, contradicting `presence == literal`.

**Priority**: high. A16 is inside the re-verify scope, and this is the fail-closed
boundary the Decision Rationale rests its drift argument on.

### H-7 — MVV scenario construction against the RDR's own gate ordering (IMP-11 + QA-11)

**Passage**: Minimum Viable Validation — scenario 1's modeled `no_match` escape
row, and scenario 2's "the same table with the key present and the value
decided FALSE prunes and escapes (D8 preserved)".

The escape row in scenario 1 is inert under the RDR's own GATE-THEN-COUNT
ordering, since `gate` returns before `escapeOrRefuse` is reached — unstated
whether it is documentary or load-bearing (IMP-11). Scenario 2's escape leg
cannot be built substantively from an RDR-0009-conforming escape row (empty
`Writes`/`RequiresOwned`), and the shipped fixture is non-conforming on both
(QA-11).

**Priority**: medium. The MVV is the RDR's headline oracle; iteration 1 already
fixed one contradiction here, and the cove pass fixed another.

### H-8 — The refusal's exit group and code coverage (PM-1 + PM-2 + QA-8)

**Passage**: A19 **Evidence** — "`guard_unevaluable` already has a row in RDR
0005's stable-code table (`flow-guard-unevaluable`, `GroupUserEnv`) — though
scoped there to 'supplied facts'"; and the `disposition` mini-check rows
assigning `GroupUserEnv`.

PM reads the group as encoding the wrong remedy: the flagship case is missing
*artifact* state, not supplied facts, and no Prerequisite fixes it (PM-2). QA
notes the `disposition` table assigns `GroupUserEnv` to
`owned_state_unavailable`, which A19 and RDR 0005's shipped code table both say
has **no row** (QA-8). PM-1 additionally challenges "Does not block lock" while
the payload's only route to a human is undecided and A19's own evidence says
`Detail`-flattening "defeats the sorted-by-row-then-key determinism."

**Priority**: high on the factual half (QA-8's table row is checkable);
PM-1's lock question is partly a re-raise of a decision the draft adjudicated
and will be graded as such at resolve.

## Single-persona findings (resolved individually, not hotspots)

- **IMP-1 (HIGH)** — empty/omitted `unless` "contributes nothing" is not a
  truth value, yet `all ∧ ¬(unless_conj)` requires one; three incompatible
  implementations, affecting nearly every fixture row.
- **IMP-3 (HIGH)** — `GuardAtom`'s field types are never fixed; a single scalar
  `literal` cannot represent RDR 0003's set-shaped literals for `in`/`contains`.
- **IMP-7** — A22's normative clause and its fallback produce identical kernel
  code; unclear whether "INPUT PRECONDITION" implies a `Resolve`-entry check.
- **IMP-8** — "no short-circuit" necessarily applies to every candidate
  (survivorship is post-verdict), making seam call-count an unstated observable
  contract on RDR 0003's code.
- **QA-4 (HIGH)** — row 19's nil-seam narrowing is non-discriminating on the
  absent-key leg; "DECIDES" conflates kernel-decidable with T/F, and the
  nil-seam payload reason is neither `absent` nor `uncomparable`.
- **QA-6** — row 5's `{T,F,U}³` matrix names no atom-level fixture-stub
  contract, so U-from-absence and U-from-uncomparable cannot be produced
  systematically; the mutation MUST names no tool, mutant set, or threshold.
- **QA-10** — row 16 asserts only "a decided verdict"; TRUE and FALSE are
  opposite dispositions, and a "not `guard_unevaluable`" assertion would also
  pass under the masking path.
- **QA-12** — `trace` step 5's witness is a *satisfied* `RequiresOwned` row,
  which passes under either ordering — non-discriminating by construction.
- **PM-3** — no scenario validates the Problem Statement's "told plainly"; all
  20 rows and 4 MVV rows assert on Go struct fields.
- **PM-4** — the aggregation veto's cost is unsized, with no dry-run affordance.
- **PM-5** — the `unless` footgun is routed to Draft RDR 0003's authoring prose
  with no gate — the "held by discipline" weakness Alternative 1 was rejected for.
- **PM-6** — a producer defect surfaces as reason `absent`, a misdiagnosis the
  closed reason set cannot express.
- **PM-9** — the both-ways diagnosis costs a two-round loop, accepted on
  incumbency after the RDR invalidates D8's original rationale.
- **IMP-11 / QA-14** — counted in H-7 and H-3 respectively; both marked
  `[out-of-delta]` by their authors.

## Review-gate note

Healthy signal on all three passes. Every finding names a passage and the
decision or test it blocks; the implementer pass explicitly dropped candidates
the shipped kernel already answers (`TagSet.Lookup`'s provenance-blind `ok`,
`missingOwned` over an empty slice, `gate`'s partition order, `Tag.Value`
opacity, the Fixup-1d key-absence premise), and QA grounded assertability
against the shipped types. No persona file references another's output.

The convergence pattern is itself evidence: PM reached H-3 and H-4 through
"can this lock / can `xg7p` close", the implementer through "what do I type on
Monday", QA through "what can I assert" — three routes to the same passages,
with no shared context.
