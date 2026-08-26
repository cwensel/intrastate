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

---

## D3 — `clierr.Finding` gains a `Count int` field (new public surface on a shared record)

- **Type**: SPEC-UNDER
- **Status**: needs author decision (proceeding under the unattended
  override with the reading below)
- **REQs**: REQ-63 (`0009:C7` "the per-identity Count MUST serialize
  alongside the identities"), REQ-44, REQ-60/61/62. Raised as `req-list.md`
  Q1 and `coverage.md` Q-B, both of which pre-select reading (a).
- **The gap**: `0009:C7` requires the per-identity `Count` to reach the
  wire, but names no carrier field for it. `clierr.Finding` — co-owned by
  RDRs 0005/0006/0008 — carries `Code, Message, Model, Severity, Param,
  Locator, Hint, Rule, Span, Element, Reason, Dimension, Key, Operator,
  Literal, Block, Class, Fingerprint` and **no `Count`**. The RDR's
  Normative Contracts name no such field, so the addition is new public
  surface on a shared record with no REQ-N of its own beyond REQ-63's
  "must serialize".
- **Resolution taken**: add `Count int \`json:"count,omitempty"\`` to
  `clierr.Finding`. The two alternatives the req-list enumerates are both
  excluded by shipped REQs: repeating the `Finding` per pre-collapse row
  contradicts REQ-42's collapse, and rendering the count into
  `Message`/`Hint` prose contradicts REQ-44 ("never only in formatted
  prose") and REQ-27/REQ-86 (no structure recovered from message text).
  `TestReq63_ThePerIdentityCountSerializesAlongsideTheIdentities` asserts a
  `count` key on the SERIALIZED finding object, so no equally-conformant
  approach avoiding the widening exists.
- **Evidence base** (`.rdr/resources.md` → Design docs →
  `docs/cli-output-contract.md`): the contract describes the finding record
  as "one flat record. Optional fields are `omitempty`, so a field that is
  present is one its producer populated", and its field table is
  illustrative of what producers populate rather than a closed set. The
  envelope's own doc comment (`clierr.go::CLIError`) grants "Extend with
  new optional fields as needed — keep them `omitempty` so the envelope
  stays append-only and stable for tools", which is the allowance
  `Findings` itself landed under (RDR 0008 per JDR 0001 §D10, recorded here
  as D1). `Count == 1` on every non-degenerate breach means the key elides
  on every path a shipped RDR asserts.
- **Blast radius checked**: no exact-field-set assertion on `Finding`
  exists in the tree. `clierr/finding_0005_test.go` asserts named-field
  PRESENCE, wire spellings, and `omitempty` over an explicit `want` list;
  its only whole-record loop rejects nested `Struct`/`Map` fields, which an
  `int` is not.
  `reserved_key_0008_test.go::TestReq102_TheFindingCarrierHasTheFiveFieldsIncludingHint`
  likewise asserts presence, not exhaustiveness. No shipped 0005/0006/0008
  test was edited or weakened.
- **Recommendation to the author**: accept the field as permanent shared
  surface and give it a REQ in whichever RDR next touches `Finding`, or
  reject it and re-open `0009:C7` to name a breach-local carrier. Until
  then the field is documented in `docs/cli-output-contract.md`'s
  finding-field table alongside the other optional keys.

---

## D1 / D2 — Stage 8 dispositions (Phase 2)

Both entries carry a **Check** whose running IS the disposition. Both were
run at the Phase 2 implementation commit.

**D1 (JD-8, cite §D10) — DISCHARGED on its substantive half.**

- `grep -c '§D10' docs/rdr/0009-*.md` → **1** (≥1 satisfied).
- Row identities ARE carried in `findings`: `kernelResolveFailure` converts
  each `*EscapeShapeBreachError`'s `RowRef` into a `clierr.Finding`
  (`Rule`/`Locator`), and the serialized envelope recorded in `coverage.md`
  shows them on `findings[]` with no identity in `detail` and no `cause`
  key on the wire.
- The escaped-plan path is asserted as a SUCCESS, not a failure:
  `TestMVV0009_…/2_conformed_table_resolves_value_for_value` and
  `TestReq74_EmptyNotNilWritesConforms` both pin `Result.Plan` non-nil with
  `Escaped: true` and a nil error.
- **No escalation.** The `findings` shape carried the row identity without
  any envelope field §D10 does not grant — `Findings` itself is §D10's
  grant. The one field that WAS added is `Finding.Count`, which is a field
  on the finding RECORD rather than on the envelope, and it is recorded
  separately as D3.

**D2 (carried citation residues) — NOT ACTIONED; recorded as out of scope.**

- **Type**: TEST-FIXTURE (unchanged)
- **Status**: OPEN — carried forward, deliberately not repaired here.
- The three named sites (`Refusal.Guard` at `0009:421`; "implemented" at
  `0009:1184`; the retired row kind `escape`, now at `0009:249` and
  `0009:497`) are all **inside the locked RDR record itself**, not in code
  or in an artifact. Repairing them means editing a Final RDR, which the
  project's standing convention forbids: RDRs are never amended, code is
  the source of truth, and the record is the history. A Stage 8
  implementer therefore has no authority to close D2 by editing it.
- The substantive half — that the shipped kernel agrees with what the
  sites describe — is satisfied independently: the kernel carries no row
  kind field at all (`TestReq6_RowKeepsItsSingleShapeWithNoKindField`
  passes, and `resolve.Row` discriminates escape identity solely by a
  non-empty `Escape` list), so the retired `escape` row kind is absent from
  the code the record governs.
- **Recommendation**: close D2 at the record's next authored revision, or
  retire it as a permanent known-residue of the locked text. It blocks
  nothing in the implementation.
