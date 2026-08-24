# Deviations — RDR 0003 Guard predicate exhaustiveness

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

---

## D1 — Empty / omitted `unless` identity is stated by 0007, not here

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0007-0003.md` F1; critique C-16 (open since
  iteration 2).
- **Conditions carried**: (a) 0003's fence (`0003:1143-1170`) is scoped to
  decided atoms and is *silent* on a group with no `unless`; the identity
  (`unless = ∅` subtracts nothing) is fenced in 0007 (`0007:1395-1409`),
  and 0003's prose at `741-747` is unfenced; (b) check below; (c) filling
  the silence with 0007's identity changes no 0003 clause.
- **Check** (Stage 8): 0003 MVV Scenario 2 run over a group with no `unless`
  block yields the full product (no subtraction); 0007 row 10 is the kernel-
  side twin. If Scenario 2 cannot pass without a fenced change, escalate.

## D2 — Stale peer claims (citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0003's next touch)
- Sites: `0003:2093-2095` MVV Scenario 4 "a predicate semantic kind that RDR
  0006 can map to a lint finding" vs this RDR's own fence `0003:920-925` and
  JDR 0001 §D7(iii) ("RDR 0006 mints nothing"); `0003:583-586` A20 "RDR 0007
  states no atom carrier" vs `0007:1522-1536` `UndecidedAtom`; `0003:1530`
  "RDR 0002 | Pending" vs `0003:1539`; 0006 as lint-finding enveloper at
  `0003:1481, 1364, 1354, 1567, 1545` vs `0005:447`.
- **Check**: `grep -n '0006 can map' docs/rdr/0003-*.md` → 0.
