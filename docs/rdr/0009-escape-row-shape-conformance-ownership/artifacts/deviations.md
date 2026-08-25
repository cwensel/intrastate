# Deviations — RDR 0009 Ownership of escape-row shape conformance

Opened 2026-08-24 to record the joint decision JDR 0001 answered after the
`0002-0009` cluster gate, iteration 4
(`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: 8 of 8, last** — §JD-5's ordering resolves here, once both this
RDR's and 0008's `Resolve`-entry checks exist. See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

---

## D1 — JD-8 answered: cite §D10

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0009's next touch)
- **Source**: JDR 0001 §D10 (`d937eec`, settled `5c2b96b`, both 2026-08-24) —
  answered *after* the iteration-4 gate. Status qualifier and README row
  corrected 2026-08-24. **§JD-5 remains open** (precondition precedence vs
  0008) and stays in the qualifier.
- **The answer**: one `Code` per caller-branchable failure with kernel-aligned
  names; one `omitempty` `findings` list on the envelope, which is the carrier
  for this RDR's row identities; exit 3 = the environment could not be
  consulted; an escaped plan is a **success**, not a failure. `GroupInternal`,
  which this RDR needs, already ships in `clierr.go`.
- **Check** (Stage 8): `grep -n '§D10' docs/rdr/0009-*.md` → ≥1; row
  identities are carried in `findings`, and the escaped-plan path is asserted
  as a success. Escalate only if the `findings` shape cannot carry a row
  identity without an envelope field §D10 does not grant.

## D2 — Carried citation residues (from iteration 2/3)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (carried; citation repair pending 0009's next touch)
- **Source**: `iter-4/reconcile-report.md` — typed CITATION REPAIR, carried.
- Sites: `Refusal.Guard` at `0009:421`; "implemented" at `0009:1184`; the
  retired row kind `escape` (×3).
- **Check** (Stage 8): each site agrees with the shipped kernel and the
  current 0002/0007 fences, or is repaired.
