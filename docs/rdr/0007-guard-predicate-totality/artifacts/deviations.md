# Deviations — RDR 0007 Guard predicate totality over an incomplete evaluation view

Opened 2026-08-24 to record the joint decisions JDR 0001 answered after the
`0002-0009` cluster gate, iteration 4
(`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

---

## D1 — JD-8/JD-21/JD-22 answered: two landings and one citation

- **Type**: TEST-FIXTURE
- **Status**: OPEN (pending 0007's next touch)
- **Source**: JDR 0001 §D10/§D12/§D13 (`d937eec`, settled `5c2b96b`, both
  2026-08-24) — answered *after* the iteration-4 gate recorded them
  unanswered. Status qualifier and README row corrected 2026-08-24.
  **§JD-18 remains open** (conforming-view enforcer, sibling 0003) and stays
  in the qualifier.
- **§D12 (§JD-21) — a landing in this RDR.** One Go type: `Block` gains
  `BlockMatch`, **added in this RDR's Phase 1**. The fence at `0007:1264-1266`
  ("exactly two constants") is the clause §D12 overrides; Phase 1's reshape is
  the named site. This is a normative change to a Final RDR's fenced
  cardinality — it cannot be done as a citation. Escalate at Stage 8 and take
  the scoped re-entry rather than widening the type silently.
- **§D13 (§JD-22) — unblocks this RDR.** The set-value encoding is RDR 0002's
  canonical sorted, duplicate-free JSON array in `Tag.Value`. This RDR's
  `contains` contract-test leg (`0007:2159-2170`, "Blocked on one RDR 0003
  declaration") and Testing Strategy scenario 8's present-key half were
  blocked on exactly this and are now writable. Note §D13 assigns the
  declaration to **0002**, not to 0003 as `0007:2182` states — that sentence
  is a citation repair.
- **§D10 (§JD-8) — citation.** This RDR's A19 is **accepted**: the envelope
  carries one `omitempty` `findings` field, and §D4's routing is reversed.
- **Check** (Stage 8): `grep -n '§D1[023]' docs/rdr/0007-*.md` → ≥1 each;
  `BlockMatch` exists in the Phase 1 type; the `contains` contract-test leg is
  written against 0002's JSON-array form; `0007:2182` names 0002 as the
  declarer. Escalate the §D12 fence change as above.
