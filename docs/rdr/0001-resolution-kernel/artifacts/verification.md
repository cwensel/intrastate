# Verification — RDR 0001 Resolution Kernel

Stage-8 implementation verification evidence.

## Phase 3b — Adversarial Review (independent lens)

Reviewer stance: assume the implementation is wrong and prove it. Scope is the
RDR 0001 **Failure Modes** clause:

> "Visible failures are typed refusal values: `no_match`, `ambiguous_match`,
> `owned_state_unavailable`, `guard_unevaluable`, or `unmodeled_outcome`. Silent
> failure would mean the kernel guessed a transition or executed persistence
> directly; both are prohibited by the normative contract. Diagnosis starts with
> the input tuple, the refusal kind, and the transition table revision used for
> that resolution."

Tests added in `internal/resolve/adversarial_test.go`. No existing test or
source file was modified. The Phase 1 suite (`TestReq*`, `TestMVV`) still passes
in full; every failure below is new signal, not a regression.

**6 of 8 added tests currently FAIL against the implementation — all six are
real defects.** Two pass and are retained as regression guards.

---

### ADV-1 — a guard-FALSE row's `RequiresOwned` poisons an otherwise legal resolution

**Status: FAILS (defect).**
`TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch`

FM clause: "Visible failures are typed refusal values … `owned_state_unavailable`",
read with REQ-15 ("the only successful selection is exactly one matching edge
**after guard evaluation**") and the req-list ASSUMPTION binding
`owned_state_unavailable` to a row "whose *evaluation* requires an owned tag
absent from the accessor-produced owned snapshot".

`Resolve` calls `missingOwned(candidates, view)` at resolve.go:333 — **before**
the guard loop at :341. `candidates` is every row matching on tag pattern alone.
A row the seam decides `GuardFalse` is not a surviving candidate and its
evaluation never needed its owned tags, yet its `RequiresOwned` entry still
converts a clean exact-one match into `owned_state_unavailable`.

Observed: kind `owned_state_unavailable`, `MissingOwned=[never-present]`.
Expected: a plan for `rdr.live`.

This is a **wrong-refusal** defect, not a missing refusal — the kernel refuses a
transition it should have planned, and names a condition that does not hold for
any surviving edge. It is the more damaging direction: the caller sees a
diagnosis pointing at owned state that is irrelevant to the edge they wanted.

**ADV-1b (`TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`) — also FAILS.**
Same root cause, sharper consequence. When *every* ordinary row is guard-FALSE,
the honest condition is zero-match → `no_match`, which RDR 0002 makes an
**escapable** class. Because the premature `missingOwned` check fires first and
returns `owned_state_unavailable` (never escapable), the table author's modeled
`no_match` escape edge is silently unreachable. A correctness bug in the
precedence order removes a modeled table capability.

**Fix direction:** move the `missingOwned` check to run over the rows that
*survive* guard evaluation, not over all pattern-matching candidates. Note the
ordering interaction: a row whose `RequiresOwned` is missing may also be the row
whose guard cannot be evaluated, so the fixup must pick and document a
deterministic precedence between the two.

---

### ADV-2 — an escape edge bypasses the guard seam entirely

**Status: FAILS (defect), both sub-cases.**
`TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`

FM clause: "Silent failure would mean the kernel **guessed a transition**",
read with REQ-15 ("exactly one matching edge **after guard evaluation**") and
REQ-23 ("guard evaluation **delegation**").

`escapeOrRefuse` (resolve.go:404-424) filters escape rows on `rescues(kind)`,
`Outcome`, and `Match` — then emits a `Plan` at :418. It **never calls
`evaluateGuard`**. Consequences:

- *guard FALSE*: the kernel emits a plan for an edge whose own predicate says it
  does not hold. Observed: plan `rdr.escape.guarded`, `escaped=true`. Expected:
  `no_match` refusal. This is literally "guessed a transition".
- *guard UNEVALUABLE*: the kernel emits a plan for an edge it cannot know holds.
  Observed: plan `rdr.escape.guarded`. Expected: `guard_unevaluable`.

The escape path is the one place the kernel is most permissive, and it is
exactly where the guard seam is skipped — a guarded escape edge is strictly
more dangerous than a guarded ordinary edge, because it fires only when the
table has already failed to resolve normally.

**ADV-2b (`TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`) — also FAILS.**
The same path skips `RequiresOwned`. Observed: a plan carrying
`Writes=[{status Escaped}]` derived from owned state the accessor snapshot never
produced. Per FM ("or executed persistence directly") and REQ-19, `Writes` is the
persistence description handed to the RDR 0004 accessor layer — describing a
write off absent owned state is the kernel guessing at persistence.

**Fix direction:** the escape candidate set must pass through the same
`RequiresOwned` and guard gates as ordinary candidates before the exact-one
count is taken. Note this changes the rescue arithmetic: an escape row rejected
by its guard should not count toward the "matches exactly once" test at :416.

---

### ADV-3 — refusal diagnosis payload depends on `Table.Rows` order

**Status: FAILS (defect).**
`TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`

FM clause: "**Diagnosis starts with the input tuple, the refusal kind, and the
transition table revision** used for that resolution", read with REQ-1 (same
input → same disposition) and REQ-3 (value-level replay determinism).

The guard loop returns on the **first** `GuardUnevaluable` row in slice order
(resolve.go:348-354), so `Refusal.Guard` and `Refusal.Rows` name whichever
undecidable row happens to sit earlier in `Table.Rows`:

```
rows [a,b] -> guard="unknown-a" rows=[{rdr.unevaluable.a flows/rdr.toml:10}]
rows [b,a] -> guard="unknown-b" rows=[{rdr.unevaluable.b flows/rdr.toml:20}]
```

Same flow, same revision, same owned snapshot, same observed tags, same
recognized outcome — different reported diagnosis. REQ-2 defines the input tuple
by the transition table **revision**, not by row sequence, so two orderings of
the same row set at the same revision are *the same input* and must replay the
same disposition. Row order is a normalization detail RDR 0002 owns; it must not
reach the reported diagnosis.

**ADV-3b (`TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`) — also FAILS.**
Same class, `owned_state_unavailable` payload:

```
rows [a,b] -> missing=[alpha beta]
rows [b,a] -> missing=[beta alpha]
```

`missingOwned`'s doc comment claims it returns keys "in stable order". It is
stable with respect to *row order*, which is not the same as stable with respect
to the *input tuple* — the property REQ-1 and REQ-10 actually require.

Phase 1's `TestReq3_ValueLevelReplayIsIndependentOfInputSliceOrder` covers
`Owned`/`Observed` slice order only. `Table.Rows` order is the surface that
actually leaks, and it was untested.

**Fix direction:** sort the diagnostic collections (`MissingOwned`, and `Rows`
by `RuleID`/`SourceLocator`) before they enter a `Refusal`; for
`guard_unevaluable`, either report all undecidable rows rather than the first,
or select deterministically by rule identity instead of slice position.

---

### Regression guards (pass today — implementation is correct here)

These two were written expecting a defect and did not find one. Retained so a
later refactor cannot silently regress the provenance boundary.

- **ADV-4** `TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot` — **PASSES.**
  `assemble` (resolve.go:145-161) applies `Observed` before `Owned`, so owned
  wins on key collision. A caller supplying `status:Final` as an observed tag
  cannot steer the kernel off the artifact's real `status:Draft` owned state.
  Correct as written; the ordering in `assemble` is load-bearing and undocumented
  as such at the call site.
- **ADV-5** `TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement` — **PASSES.**
  `missingOwned` checks `view.has(key, ProvenanceOwned)`, so an observed tag
  carrying a required owned key does **not** satisfy the requirement; the kernel
  still refuses `owned_state_unavailable` and names the key. Correct: caller
  context cannot stand in for accessor-produced artifact state.

---

### Note for the fixup phase

ADV-1/1b and ADV-2/2b are the same underlying shape from two directions: the
gating checks (`RequiresOwned`, guard) and the candidate set are applied at
inconsistent points in the pipeline. Ordinary candidates get `RequiresOwned`
checked *too early* (before guards prune them); escape candidates get both
checked *not at all*. A fix that only patches one path will leave the other.
The cleanest repair is a single `viable(row) → (ok, refusal)` gate applied
uniformly to ordinary and escape candidates alike, with the exact-one count
taken over the survivors in both phases.

ADV-3/3b are independent of that and can be fixed separately — they are payload
ordering, not selection semantics.
