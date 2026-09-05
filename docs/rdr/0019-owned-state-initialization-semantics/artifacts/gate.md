# Finalization Gate — cli/0019 owned-state-initialization-semantics

Date: 2026-09-05 · Verdict: READY (Gate PASS) · Tooling pass: PASS
(`evidence/tooling-pass/tooling-pass.md`) · Joint-decision fence: clear
after re-run (see the `Joint-check:` line in Decision Rationale).
Cross-Cutting Concerns stays in the record as `cli/0019:G-cross-cutting`.

## Contradiction Check

No contradictions found between research findings, design principles,
and the proposed solution.

- Research → solution: the prior-art class read (initialization is an
  explicit session-start act, whole-snapshot, never a read-time
  fallback) is exactly D's shape; Key Discoveries' load-time writer
  coverage (`checkAccessorBindings`, A3) is what C1 relies on instead
  of a lint arm, and the Risks section records the earlier lint-arm
  reading as discharged, not retained; the docs already promising
  start-state semantics are delivered, not retracted.
- Principles: 0004:C3 read purity and REQ-107 cleared ≠ unseeded both
  hold by construction (no read path synthesizes; the predicate is
  store emptiness, never per-key absence). The one tension is disclosed,
  not contradicted: emptying a store by clears returns it to the
  initializable class (C1's residual), pinned by MVV step 9 and named in
  Risks and Consequences.
- Cross-record: the addition to 0005:C1's closed verb enumeration is a
  recorded additive override (Overrides, A4) on the 0010 → 0002:C19
  precedent; cli/0028's fire on this record is homed at cli/0019:C1 and
  now carried symmetrically here.

## Assumption Verification

All eight Critical Assumption records are internally consistent:
Status `Verified` on each, Method in the sanctioned vocabulary (Source
Search: A1, A3, A7, A8; Spike: A2, A5; Source Search + Spike: A6; Peer
RDR: A4), Evidence agreeing with the Method (spike commands and
read-back output where Spike; `path::Symbol` anchors where Source
Search; `0005:C1` / `0005:D-naming` elements where Peer RDR), and a
non-empty "If wrong" on every record. No `Docs Only` record. No
`Pending` or `Unverified` record, so status consistency is vacuous. No
self-referential `Verified` stamp: every Source Search anchor resolves
into `internal/…` on `main` (lint `--locking`: 0 unresolved, 0
unlooked), and A4's evidence lands on peer elements, not a record.
Five Evidence fields exceed the soft budget (A2, A4, A6, A7, A8): each
opens on its load-bearing anchor and the balance is verification
content — kept as is.

## Scope Verification

The Minimum Viable Validation is in scope and executes in Phase 1, not
deferred. The proof is the ten-step MVV over a fixture state-machine
model (one always-present and one plain owned key) plus step 10 over
the repo's own `models/rdr.toml`: `flow init-state` seeds an empty
store and `flow read-state` returns the `[initial]` values in canonical
form (RT1); the re-run commits zero writes (RT2); a cleared key stays
absent through a later `init-state` (RT4 / A5); the decision-table,
unbound-role and non-file-backed-carrier invocations refuse with zero
writes (S7, S6, S10); and step 9 asserts the emptied-store boundary in
the failing direction. Testing Strategy scenarios S1–S13 are the
executable form, including S4's per-kind byte-identity table against
the manual `set-state` path (RT3).

## Proportionality

Right-sized on the split test: `contracts=1`, one durable contract
(C1), this record its sole author, one seam (`[initial]` runtime
semantics and its carrier). Profile re-validated against that count:
`foundational — … (C1); user-facing yes; locks cross-rdr` is value +
one clause naming the contract, and the two judgements hold (a new
user-facing verb; an override of a Final peer's closed enumeration).
The lens battery that Profile requires ran to completion — cove,
3amigo with consolidation, critique on two models with diff,
repeatability variant full (three runs + diff) — reconcile closed
every spike, and the accretion floor does not apply (Seam Lineage: no
prior accretion). Nothing to trim before lock: the length is in C1's
ordering, predicate and carrier clauses, each of which a lens forced
by naming the defect its absence admits, and in the spike-bearing
Evidence fields judged above. The four judged gate items leave the
record for this file at lock.
