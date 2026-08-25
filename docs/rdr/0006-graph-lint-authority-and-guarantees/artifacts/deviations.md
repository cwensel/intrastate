# Deviations — RDR 0006 Graph lint authority and guarantees

Pre-seeded by the `0002-0009` cluster gate, iteration 3 (2026-08-24 —
`docs/rdr/cluster-reconcile/0002-0009/iter-3/reconcile-report.md`). Stage 8
opens this file; running each named check IS the entry's disposition. An entry
escalates only if its check contradicts a contract.

**Run order: 4 of 8** — after 0003, whose tag declaration model IS this RDR's
lint input contract; ahead of 0004/0005. D2's stale citations include sites
that 0002's and 0003's re-locks moved. See
[`../../BUILD-ORDER.md`](../../BUILD-ORDER.md).

---

## D1 — §D7 landings absent from 0006's text (A6/A10 flip; terminal-non-owned finding; terminal spelling)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (deferred at cluster gate; DEFER-TO-IMPLEMENTATION)
- **Source**: JDR 0001 §D7 (i)/(iv), §JD-17 — answered 2026-08-24;
  `iter-3/answer-check-0006.md`.
- **Conditions carried**: (a) all sites unfenced — A6 (`0006:232`, `263-268`),
  A10 (`0006:372`), invariant 2 prose (`542-548`), the "no peer states that
  rule" sentence (`594-596`); (b) mechanical checks below; (c) additive —
  no clause's meaning changes (fenced `1014-1018`, `843-847`, `862-874` are
  CONSISTENT with §D7).
- **Checks** (Stage 8):
  1. `grep -n 'Pending' docs/rdr/0006-*.md` at A6 and A10 → both read
     `Verified — by citation to RDR 0002 §<layout clause>` after 0002 re-locks
     under §JD-17 (`[initial]` root assignments; `terminal` context-id list;
     write-replaces clause). If 0002's re-locked layout lacks any of the three,
     escalate as SPEC-DEFECT against 0002.
  2. Lint test: a root `terminal` entry naming a context whose predicate
     reads a **non-owned** tag → one blocking finding (§D7(i)). Add the code
     row and a scenario; if 0006's finding taxonomy cannot admit the code
     without a fenced change, escalate.
  3. `0006:263-268` "not a bare identifier" wording vs §D7's context-id list:
     repair the prose to cite §D7(i); no test.

## D2 — Stale peer-status claims (cosmetic, citation repair)

- **Type**: TEST-FIXTURE
- **Status**: OPEN (citation repair pending 0006's next touch)
- Sites: `0006:14, 141-142, 253-256, 256, 280, 380, 391, 813, 1421, 1425` —
  "RDR 0002 is `Draft` … scheduled edit on an open peer" and the pre-§D7
  layout description. Artifact of record: `docs/rdr/README.md` Index.
- **Check**: `grep -n 'Draft' docs/rdr/0006-*.md` returns no peer-status
  claim that disagrees with the README Index.
