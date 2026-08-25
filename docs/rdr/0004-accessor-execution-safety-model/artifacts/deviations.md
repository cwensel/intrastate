# Deviations — RDR 0004 Accessor execution safety model

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: 5 of 8** — after 0007, 0002 and 0006. Read-back equality is byte
equality over 0002's canonical JSON array; §D9 settles the `gate denied`
refusal/non-refusal split (D3 below). See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

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

## D3 — JD-19/20/22 answered: cite §D9, §D11, §D13

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0004's next touch)
- **Source**: JDR 0001 §D9/§D11/§D13 (`d937eec`, settled `5c2b96b`, both
  2026-08-24) — answered *after* the `0002-0009` iteration-4 gate. Status
  qualifier and README row corrected 2026-08-24. No joint decision is now open
  against this RDR.
- **The answers**: §D9 (§JD-19) — gates run after exact-one selection and
  before the plan, only the selected row's; `deny` is a **refusal**, not an
  escape class; deny-overrides with every gate reported. This settles the
  split this RDR currently shows, which lists `gate denied` once as a refusal
  (`0004:804`) and once as a typed non-refusal result (`0004:502`); §D9 picks
  the refusal reading. §D11 (§JD-20) — `--clear <key>`, JSON-array set writes,
  unbound `--write` keys refused at the CLI **before any accessor runs**, which
  is the boundary this RDR left unstated. §D13 (§JD-22) — read-back equality is
  byte equality over RDR 0002's canonical JSON array, which is the "unspecified
  form" this RDR's read-back assertion rested on.
- **Check** (Stage 8): `grep -n '§D9\|§D11\|§D13' docs/rdr/0004-*.md` → ≥1
  each; the `0004:502` non-refusal reading of `gate denied` agrees with §D9 or
  is repaired. If the `502`/`804` split cannot be reconciled by citation,
  escalate — §D9 decides it, so a surviving contradiction is a fenced conflict.

---

## D9 — REQ-126 is a status disclosure, not a testable obligation

- **Type**: TEST-FIXTURE
- **Status**: RESOLVED at Phase 1 (no author decision needed)
- **Source**: Stage 8 Phase 1 coverage audit; `artifacts/coverage.md`.
- REQ-126 reads: "Prerequisites carried unchecked at Gate PASS: RDR
  0001's stateless resolver returning planned owned-tag writes, RDR
  0002's accessor references / provenance / artifact roles through
  normalization, RDR 0003's consumption of accessor-produced values
  without executing accessors." — (IP Prerequisites; deviations D2 Note)
- **Disposition**: it names the *state of the record's checkboxes at
  lock*, not a behaviour the accessor boundary can exhibit or violate.
  There is no input for which an implementation could take a wrong
  branch, so no test can fail if a future change "broke" it. Phase 1
  therefore records it as covered-by-construction rather than authoring
  a test that asserts a documentation fact.
- **What IS tested**: the three peer capabilities REQ-126 names are
  exercised where this RDR actually consumes them —
  `TestReq118_…` and `TestReq34_…` drive RDR 0001's shipped kernel
  (`resolve.Resolve`, `Refusal.MissingOwned`, `TagSet.has`);
  `TestReq8_…` and `TestReq15_…` consume RDR 0002's shipped
  `table.Accessor` carrier and its `Keys`/`ReadBack` fields;
  `TestReq44_…` consumes RDR 0002's `table.ClearSentinel`. RDR 0003's
  clause (guards do not execute accessors) is a property of the guard
  package, outside this RDR's surface — no accessor test can witness it.
- **Basis**: §scope-discipline — a clause with no reachable input at
  this boundary is not this RDR's to assert.

## D10 — the eighth validation arm's code spelling

- **Type**: NAMING
- **Status**: RESOLVED at Phase 1 under the req-list's standing
  ASSUMPTION (no author decision needed)
- **Source**: REQ-17 / REQ-84; req-list ASSUMPTION "the eighth
  validation arm's code (REQ-17/REQ-84) follows the seven witnessed
  names' spelling convention (snake_case, defect-named); the record
  names the arm but not its code."
- **Disposition**: Phase 1 fixes the spelling as
  `missing_requested_key_set` — snake_case, defect-named, in the family
  of the seven witnessed spellings (`missing_accessor`,
  `missing_write_read_back`). ORA 1 requires each arm to assert "its
  *own named* code" and fixes only seven spellings, so the eighth is
  chosen at implementation exactly as the ASSUMPTION anticipated.
- The tests bind the spelling through the exported constant
  `accessor.CodeMissingRequestedKeySet` rather than a literal, so a
  later rename is a one-line change that does not silently pass.

## D11 — `gate denied` reconciled on the layered reading (Q1 / D3 check)

- **Type**: CONTRACT-READING
- **Status**: RESOLVED at Phase 1 by citation — D3's Stage 8 check ran
  and did NOT escalate
- **Source**: req-list Q1 and REQ-37/REQ-38; deviations D3's check ("the
  `0004:502` non-refusal reading of `gate denied` agrees with §D9 or is
  repaired … a surviving contradiction is a fenced conflict").
- **Disposition**: the split is reconcilable by citation and no fenced
  conflict survives. JDR 0001 §D9 states both halves explicitly — "Deny
  is a refusal at the CLI and a typed result at the accessor" — so
  `0004:502` (typed gate result) and `0004:804` (listed among the
  caller-visible typed refusals) are each true AT THEIR OWN LAYER.
  Phase 1 implements the accessor half only: `TestReq37_…` asserts a
  deny is a typed `Verdict` carrying a reason and that no `gate_denied`
  member exists in the accessor refusal-class set. The CLI half is RDR
  0005's and is not asserted here.
- **Escalation condition, not met**: §D9's own note says escalation is
  required "only if an implementation attempt shows the layered reading
  unbuildable". The reading is buildable — the two halves live on
  different return types (`GateResult.Verdict` vs the CLI's refusal
  mapping) and never collide in one value.

## D12 — the gate SITE is not implemented here (Q2)

- **Type**: SCOPE
- **Status**: RESOLVED at Phase 1 under the req-list's reading (b)
- **Source**: req-list Q2 and REQ-38; JDR 0001 §D9 ("0004 never states
  the site"), which assigns the site to RDR 0005.
- **Disposition**: Phase 1 ships gate INVOCATION semantics only —
  allow / deny / indeterminate, bounded timeout, accessor-refusal
  classification of a gate's timeout or execution failure, and
  "reported, never applied". It ships NO gate scheduling: no test drives
  when gates run, which gates run for a selected row, deny-overrides
  aggregation, or `set-state`'s exemption. `gate_0004_test.go`'s header
  records the boundary.
- **Basis**: §D9 lands the site in 0005; this RDR's own Approach says
  "The resolver stays stateless" and CA A5's If-wrong is "The accessor
  layer becomes an orchestrator" — both forbid reading (a).
- **Not validatable yet**: RDR 0005 is unimplemented, so no shipped
  caller exists to check the invocation surface against. Recorded, not
  blocking.

## D13 — the pre-seeded D1/D2/D3 checks, run at Phase 1

- **Type**: TEST-FIXTURE
- **Status**: RESOLVED for the test half; the record half is NOT Stage
  8's to touch
- Ran each pre-seeded entry's own named check. Results:

| Entry | Check | Result | Disposition |
| --- | --- | --- | --- |
| D1 | author the escaped-plan / empty-write-set test in 0004's suite | **DONE** — `TestReq42_NoWriteRunsWithoutASuccessfulPlansWrites/escaped_plan_carries_an_empty_write_set_and_writes_nothing`, with a control asserting a NON-escaped plan does write, so an executor that never writes cannot pass vacuously | the obligation's testable half is discharged (REQ-43). The write accessor honours it without a fenced change, so D1's escalation condition is **not met** |
| D1 | `grep -c 'NextTags' docs/rdr/0004-*.md` → ≥1 once cited | **0** | record-side citation, see below |
| D2 | `grep -n 'escapable .no_match\|A6b (Pending'` → 0 | **1** | record-side citation repair, see below |
| D3 | `grep -n '§D9\|§D11\|§D13'` → ≥1 each | **1 / 1 / 1 — PASS** | closed; the `gate denied` split reconciles by citation (D11 above), so D3 does not escalate |

- **On the two record-side residues (D1's grep, D2)**: both ask for
  edits to `docs/rdr/0004-accessor-execution-safety-model.md` itself.
  **RDRs are never amended** — they are prompts and a design-decision
  history, and code is the source of truth. Stage 8 therefore records
  the disposition here rather than editing the record, which is exactly
  what D2's own framing anticipates ("Status: OPEN (citation repair
  pending 0004's next touch)").
- **What Phase 1 implemented against**, per the req-list's standing
  ASSUMPTION on D2: RDR 0007's fenced `guard_unevaluable` and RDR 0002's
  shipped provenance scoping — never the stale prose. `TestReq34_…`
  is worded to that effect: it asserts the match-only residual lands on
  `no_match` as the SHIPPED kernel decides it
  (probed against `resolve.Resolve` directly), and pins that the
  accessor layer must not "fix" it by synthesizing a value.
