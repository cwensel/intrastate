# Deviations — RDR 0003 Guard predicate exhaustiveness

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Build order: step 4** — consumes 0007's atom shape and 0002's normalized
model; 0006's lint input contract reads the tag declaration model this RDR
owns. See [`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

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

## D3 — §JD-22 answered: cite §D13 for the set-value encoding

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0003's next touch)
- **Source**: JDR 0001 §D13 (`d937eec`, settled `5c2b96b`, both 2026-08-24) —
  answered *after* the `0002-0009` iteration-4 gate recorded §JD-22 as
  unanswered. Status qualifier and README row corrected 2026-08-24 (§JD-22
  dropped; §JD-18 retained).
- **The answer**: `Tag.Value` stays `string`; a set crosses the kernel seam as
  its canonical JSON array — members sorted, duplicate-free, compact encoding —
  and **RDR 0002 declares it**. Read-back equality is byte equality; no third
  encoding exists.
- **Scoped answer-vs-fences check: CONSISTENT.** §D13 governs *carriage* of a
  set tag value across `Tag.Value`; this RDR owns the set **literal** spelling
  (A13, `0003:1241-1248`) and the **element universe** (A9) — distinct things,
  and §D13 says so ("0003 … declares spelling and universe only"; "Encoding is
  carriage, not declaration semantics"). No fenced clause here contradicts it:
  A13's unordered/duplicate-free literal and §D13's sorted/duplicate-free
  carriage share one normal form, and `grep -n 'Tag\.Value'` over this RDR
  returns no fenced tag-value byte claim. The silence §D13 fills is real
  silence, so §D13 *narrows* what 0007 had assigned here rather than
  contradicting a clause.
- **Effect on peers**: 0007's `contains` contract-test leg (`0007:2159-2170`,
  "Blocked on one RDR 0003 declaration") and Testing Strategy scenario 8's
  present-key half are unblocked — the encoding is 0002's, and it is stated.
- **Check** (Stage 8): A13's Evidence and the References list cite JDR 0001
  §D13, noting the tag-value encoding is RDR 0002's while the literal spelling
  stays A13's. `grep -n '§D13' docs/rdr/0003-*.md` → ≥1. If the citation cannot
  be added without a fenced change, escalate.
