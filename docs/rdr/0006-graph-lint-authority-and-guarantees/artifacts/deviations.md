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

---

# Stage 8 Phase 1 (test authoring)

Recorded by the Phase 1 test author. Neither entry changes a REQ; both
record a gap between a clause and the input surface the landed peers
actually expose, so Phase 2 does not read the test shape as a shortcut.

## D3 — `graph-single-valued-state` has no authorable input surface

- **Type**: SPEC-vs-LANDED-PEER
- **Status**: OPEN (Phase 2 to confirm the check is implemented as a total
  function over its input, even though no fixture can trigger it)
- **REQ**: REQ-40 (invariant 5, "no row's write block may assign a
  single-valued tag two values"); REQ-31 lists the class as mandatory.
- **Finding**: RDR 0002's loader refuses every TOML spelling of the defect
  *before* normalization, so lint never receives it. Probed against the
  landed `internal/table` loader:
  - `enum` / `int` / `bool` / `scalar` declared `single_valued = true` and
    written a member sequence → `malformed_tag_declaration: rule r write
    <k>: kind <kind> holds one value, not a member sequence`.
  - `set` declared `single_valued = true` → `malformed_tag_declaration:
    tag <k>: kind set admits no single_valued marker`, so the one kind
    that admits a member sequence cannot carry the marker.
  - `[initial]` assigning two members to a single-valued key →
    `malformed_initial_declaration`.
- **Disposition**: `TestReq40_SingleValuedStateIsDecidedPerRowSyntactically`
  asserts (a) the code is a declared, blocking member of the taxonomy, and
  (b) the per-row reading does not fire on the legal multi-row shape,
  rather than driving the defect through a fixture that cannot exist.
  REQ-41's negative half (a merged node holding two values is NOT a
  violation) is fully testable and is tested.
- **Not escalated**: the clause is not wrong — it is the model-level half
  of RDR 0003's single-valued conformance conjunct, and stating it keeps
  the invariant set total. It is simply discharged upstream today. If a
  later RDR 0002 revision admits a `set` with a single-valued marker, the
  fixture becomes authorable and this entry closes.

## D4 — A8's checked-in transition model does not exist yet

- **Type**: PREREQUISITE
- **Status**: OPEN (Phase 2 authors the model, the CI job, and the
  Makefile edge)
- **REQ**: REQ-119, REQ-120, REQ-121, REQ-122.
- **Finding**: no transition model is checked into this repo; the only
  `.toml` files are `internal/table/testdata/` fixtures, which SC-7
  explicitly excludes ("not a hook wrapper, unit-test-only engine path, or
  fixture-only corpus"). `.github/workflows/ci.yml` carries no
  `graph-lint` job, and the `Makefile`'s `check` target lacks the `build`
  edge IP Phase 3 requires.
- **Disposition**: per ASSUMPTION-9 the model is authored during this
  implementation. `internal/cli/lint_gate_0006_test.go` pins its home at
  `models/rdr.toml` and is red until it exists, so the prerequisite is a
  failing gate rather than a silently dropped REQ. REQ-MVV proper is
  fixture-backed and does not consume the CI gate, so the MVV is
  independently satisfiable.
