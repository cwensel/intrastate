# Deviations — RDR 0002 Transition table as reviewable data

Pre-seeded by the `0002-0009` cluster gate, iteration 4 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-4/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

---

## D1 — Canonical fixtures vs RDR 0003's declaration vocabulary

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0002-0003.md` F1; `iter-4/critique-set.md` Q-1.
- **Conditions carried**: (a) the offending text is fixture/evidence/
  illustrative — `evidence/spikes/iter-2/rdr-fixture.toml:65` (`kind =
  "string"`), Illustrative Code `0002:1698`, A13 evidence `0002:340`,
  Round-Trip prose `0002:1551`, and the fixture's `single_valued` /
  assignment shape on `iter` / `cluster_ready`; 0002's own fence defers
  the kind vocabulary to RDR 0003 (`0003:791-793`: exactly five tokens,
  `scalar` not `string`) and the assignment table to `0003:948-952`,
  `1273-1290`; (b) checks below; (c) no clause's meaning changes — the
  fences already agree, the artifacts lag them.
- **Checks** (Stage 8):
  1. `grep -rn 'kind = "string"' docs/rdr/0002-*/evidence/spikes/iter-2/`
     → 0 after the fixtures are rewritten to 0003's five tokens; the loader
     refuses an unknown kind token as `malformed tag declaration`.
  2. Run RDR 0006's lint (or, before it exists, 0003's assignment table by
     hand) over the promoted fixtures: `iter` / `cluster_ready` declared
     `single_valued` must satisfy `0003:948-952`. If a fixture cannot be
     made to conform without a fenced change in 0002 or 0003, escalate as
     SPEC-DEFECT.

## D2 — Empty write block on an escape row has no fixture

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0009-0002.md` F1 (JD-11 row 3 residual).
- **Conditions carried**: (a) both fences agree (`0002:833-836`,
  `0009:843-846`: `writes = []` is a presence-keyed rejection); the gap
  is in evidence — the promoted spike checks `len(rule.Write) > 0`
  (`iter-2/main.go:546`) and `gen-cases.py:60-62` populates
  `neg-escape-with-write`; (b) check below; (c) no meaning change.
- **Check** (Stage 8): mint `neg-escape-with-empty-write` (`write = []` on
  an escape row) and assert `malformed escape declaration`; the loader must
  key on key presence, not length. RDR 0009 Scenario 8 is the second oracle.

## D3 — Write-block value kind/domain conformance names no load category

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: `iter-4/pairwise-0002-0006.md` F1.
- **Conditions carried**: (a) unfenced on both sides — 0002's arity rule
  claims the check "at minimum" but lists no category; 0006 checks it in
  invariants 1 and 5; (b) check below; (c) additive — no clause changes.
- **Check** (Stage 8): one Scenario-3 fixture writing a literal outside the
  declared domain (and one of the wrong kind) → a named 0002 load category
  (`malformed tag declaration` / literal-outside-domain per §D7(iii)).
  If the category set must widen in fenced text, escalate.

## D4 — Stale peer-status sentences inside fences (citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0002's next touch)
- Sites: `0002:1053-1058` ("RDR 0004 does not yet carry that clause … not
  end-to-end until RDR 0004 lands §D5") — 0004 re-locked 2026-08-24 with the
  clause fenced at `0004:357-364`; `0002:1229-1233` ("RDR 0003 is `Draft
  [revised from Final 2026-08-24]`") — 0003 re-locked the same day. Both are
  self-scoped status claims; no obligation changes. Also `0002:2073-2077`
  ("the reshape has no owner") vs 0007 Phase 1 (`0007:2117-2126`), and
  `0002:2297-2298` Perf Expectations "suffix … an alphabet member" vs the
  fenced sequence-suffix rule. Artifact of record: `docs/rdr/README.md`.
- **Check**: `grep -n 'does not yet carry\|has no owner' docs/rdr/0002-*.md`
  → 0, and the `0003 is Draft` sentence at `0002:1229` gone.
