# Deviations — RDR 0001 Resolution Kernel Contract

Phase 2 (implementer) artifact. Every entry records a gap between the locked
RDR / frozen Phase 1 tests and what the implementation had to decide, with its
Type, grounding evidence, and status.

No Phase 1 test was weakened, deleted, skipped, or edited. `internal/resolve/resolve.go`
is the only source file changed.

---

## D1 — Escape-edge scope: can an escape rescue `owned_state_unavailable` or `guard_unevaluable`?

- **Type**: SPEC-UNDER
- **Status**: needs author decision (recorded; implementation continued with the
  recommendation below)
- **REQs**: REQ-15, REQ-16

### The gap

RDR 0001 REQ-15 reads: "Zero, multiple, unavailable, or unevaluable candidates
are refusals **unless the table contains a modeled escape edge** that itself
matches exactly once." Read literally, the escape clause governs all **four**
listed conditions, so an escape edge could rescue `owned_state_unavailable` and
`guard_unevaluable` as well as `no_match` and `ambiguous_match`.

Phase 1 deliberately wrote no test either way (`coverage.md`, "Spec tension
noted"): `TestReq15_AllFourListedConditionsRefuseWithoutAnEscapeEdge` asserts
only that all four refuse when **no** escape is modeled — true under both
readings — and the rescue tests exercise only `no_match` / `ambiguous_match`.
So the red-before-green gate is blind to this clause.

### Evidence consulted

- `docs/rdr/0001-resolution-kernel.md` Risks (line ~412): "Make multiple matches
  a refusal; **any priority or escape must be explicit table data owned by RDR
  0002**." RDR 0001 delegates escape-list *content* to RDR 0002.
- `docs/rdr/0002-transition-table-as-reviewable-data.md` normative block (line
  292): "An `escape` list MUST contain only resolver failure classes that RDR
  0001 allows the table to model: `no_match` and `ambiguous_match`."
- `artifacts/req-list.md` ASSUMPTION (escape scope) records the same narrowing.
- `evidence/critique/critique.md` acceptance tests construct escape rescue only
  for the zero/multiple-match conditions.
- `.rdr/resources.md`: no corpus or design doc speaks to escape scope; the
  resolution rests entirely on the two RDR texts above.

### Recommendation (implemented)

Escapes rescue **only `no_match` and `ambiguous_match`**. Because RDR 0001
itself assigns escape-list ownership to RDR 0002, and RDR 0002 closes the list
to those two classes, the narrow reading is the one RDR 0001's own evidence base
supports; the four-item list in REQ-15 is best read as enumerating the refusal
conditions, with the escape clause attaching to the two that RDR 0002 models.

As implemented this needs no special case: `owned_state_unavailable` and
`guard_unevaluable` are decided **before** selection completes, so control never
reaches the rescue phase for them. `escapeOrRefuse` is called only on the
zero-match and multiple-match paths. Should the author choose the broad reading,
the change is to route those two refusals through `escapeOrRefuse` as well — no
restructuring.

**Author decision needed**: confirm the narrow reading, or direct that RDR 0001
REQ-15 be honoured literally over RDR 0002's narrowing.

---

## D2 — Tag key for the freshly recognized outcome in the evaluation view

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-17

REQ-17 requires the freshly recognized tag to be merged into the evaluation view,
but neither the RDR nor `req-list.md` names the key it takes there. The frozen
Phase 1 fixture pins it: `recognizedTagSensitiveTable` (`fixtures_test.go`)
matches on `{Key: "recognized", Value: "successful"}` against an input whose
`Recognized` field is `"successful"`.

Implemented as the unexported constant `recognizedTagKey = "recognized"`. This
adds **no** public surface — the key is observable only through table data the
caller already writes, and the constant is package-private. Recorded because a
future RDR (0002 normalization) must spell the same key for a table row to match
on the recognized outcome as a tag.

---

## D3 — Provenance precedence on tag-key collision

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-17, REQ-3

Three tag sources merge into one keyed view (REQ-17); the RDR does not say which
wins when owned, observed, and recognized supply the same key. REQ-1/REQ-3
require only that the choice be deterministic and independent of input slice
order.

Implemented as **owned > observed > recognized**: the accessor-produced snapshot
is authoritative and can never be shadowed by caller-supplied context, which is
the ordering consistent with RDR 0004 treating owned tags as the artifact's
truth. No frozen test exercises a collision, so this is latitude, recorded
because it is contract-visible to RDR 0002 table authors.

---

## D4 — Refusal precedence order

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-1, REQ-32

`req-list.md` QUESTIONS §2 notes the RDR states no total ordering across refusal
conditions and requires only determinism, with `unmodeled_outcome` before
`no_match` grounded in `evidence/critique/critique.md` acceptance test 3.

Implemented order: `unmodeled_outcome` → `owned_state_unavailable` →
`guard_unevaluable` → exact-one / `no_match` / `ambiguous_match`. Each stage
returns before the next, and rows are visited in table order, so the disposition
is a pure function of the input tuple (REQ-1). Documented on `Resolve`'s doc
comment so future readers do not re-derive it.

---

## D5 — Owned-state and guard checks are scoped to matching candidates

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-5, REQ-33

The RDR says required owned state being unavailable is a refusal but does not say
whether the check spans every table row or only the rows in play. Checking every
row would make an unrelated row's `RequiresOwned` key poison an otherwise legal
resolution.

Implemented as: `RequiresOwned` and guards are evaluated only over **candidate**
rows — non-escape rows whose `Outcome` equals the recognized outcome and whose
`Match` pattern holds over the assembled view. Grounded in
`req-list.md`'s ASSUMPTION that `owned_state_unavailable` arises when "a
candidate row's evaluation requires an owned tag absent from the
accessor-produced owned snapshot". `missingOwnedTable`'s candidate matches the
view and requires the absent `gate` key, so the fixture is satisfied.

---

## D6 — Escape rows are excluded from ordinary candidate selection

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-15, REQ-16

Required by the frozen tests, not chosen freely. In
`TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce/escape modeled for
another class does not rescue`, an escape row whose `Match` pattern **does** hold
over the view is appended to a no-match table, and the expected disposition is
still `no_match`. Were escape rows ordinary candidates, that row would have been
selected as an exact-one match and produced a plan.

Implemented: a row with a non-empty `Escape` list is skipped during candidate
selection and considered only in the rescue phase, for its declared class.

---

## D7 — `slices` standard-library import

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQs**: REQ-26, REQ-37

`golangci-lint` (modernize) flagged two membership loops. Resolved with
`slices.Contains`, which requires importing `slices`. Confirmed benign against
both negative-contract tests: `slices` is standard library, so
`TestReq26_KernelUsesNoThirdPartyDependency` (`isThirdParty` keys on a dot in the
first path element) and `TestReq37_KernelIntroducesNoHashOrCanonicalSerialization`
(explicit hash/encoding import list) both stay green. No third-party dependency
was added; `go.mod` is unchanged.

---

## Summary

| Count | Category |
| --- | --- |
| 6 | mechanical translation (D2–D7) |
| 1 | needs author decision (D1) |
| 0 | halted on |

No SPEC-DEFECT, no DEPENDENCY-LIMIT, and no TEST-FIXTURE defect was found: the
Phase 1 test suite and fixtures were internally consistent and every frozen
assertion was satisfiable without edit.

**New public surface added beyond the RDR's Normative Contracts: none.** The
Phase 1 skeleton already declared every exported symbol (`Tag`, `Provenance`,
`RefusalKind` + 5 constants, `RefusalKinds`, `GuardResult` + 3 constants,
`GuardEvaluator`, `TagSet` + `Lookup`/`Len`, `Row`, `Table`, `Input`, `Plan`,
`Refusal`, `RowRef`, `Result` + `Refused`, `Resolve`); Phase 2 added only
unexported helpers and the unexported `recognizedTagKey` constant.
