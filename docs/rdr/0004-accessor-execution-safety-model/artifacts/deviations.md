# Deviations — RDR 0004 Accessor execution safety model

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

---

## D1 — RDR 0009's write-accessor obligation and test bind nobody here

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0009-0004.md` F1 (iter-2 F2 / iter-1 F1;
  JDR 0001 §JD-7 residue).
- **Conditions carried**: (a) unfenced — 0009 files the obligation at
  `0009:1515-1527`; 0004 mentions 0009 zero times; 0004's fence
  `0004:353-354` already forbids the regression; (b) check below; (c) no
  clause's meaning changes — the test is additive.
- **Check** (Stage 8): author the escaped-plan / `NextTags` test 0009 names
  (an escaped plan reaches the write accessor with an empty write set and
  performs no write) in 0004's suite; `grep -c 'NextTags' docs/rdr/0004-*.md`
  → ≥1 once cited. If the write accessor cannot honour it without a fenced
  change, escalate.

## D2 — Stale peer claims and restatements (citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0004's next touch)
- Sites: `0004:495, 523-527, 828-831` — a guard-consumed absent key "falls to
  escapable `no_match`" vs 0007's fenced `guard_unevaluable`
  (`0007:1329-1331`; JDR 0001 §D4) — pre-§D1 "guard string" vocabulary;
  `0004:1015` cites 0007 A6b as "Pending, downgraded at Stage 6" vs
  `0007:485` Verified/closed by §D3; `0004:254-257` restates §D7(ii)'s
  zero-readers rule in its pre-provenance-scope form (settled by 0002 at
  `0002:675-707`, now carried at §D7(ii)); the `0004:371`-successor line
  cites RDR 0008 vacuously (`grep -c 'RDR 0008'` → 0).
- **Check**: `grep -n 'escapable .no_match\|A6b (Pending' docs/rdr/0004-*.md`
  → 0.
- Note (not a deviation): Prerequisites `0004:841-849` are all unchecked at
  Gate PASS; the Status line discloses A9/A10/A11 Pending. A per-RDR finalize
  question, recorded here so Stage 8's opener sees it.
