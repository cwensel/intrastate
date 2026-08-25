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

---

# Stage 8 Phase 2 (implementation)

Recorded by the Phase 2 implementer. Every entry below was resolved
against the record's own evidence base and the landed peers; none blocked
progress and none weakened a Phase 1 assertion.

## D5 — `respond.FindingCarrier` is a new public surface the record does not name

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice recorded)
- **REQ**: REQ-97 ("Text mode MUST enumerate every finding's code and
  message, which requires extending … the `respond.OK` text branch
  (success, which renders only `Notes`/`Warnings` today and drops
  `Data`)").
- **Finding**: the record books the edit but names no mechanism. The text
  branch must read the verb's `Data` payload, and `respond` must not
  import `internal/graphlint` — the same leaf-package discipline
  `0006:C14` fixes for `clierr`. Three shapes were available: a
  `Findings` field on `respond.Success` (a verb-specific field the record
  explicitly declines — "rather than growing `respond.Success` a
  verb-specific field"), a type switch on the concrete payload (which
  imports the producing package), or an interface the payload satisfies.
- **Evidence**: `0006:C14` fixes both constraints — success findings
  travel "under the existing `respond.Success.Data` payload" and the
  gateway gains no dependency on the producing package. An interface is
  the only one of the three that satisfies both.
- **Chosen**: `respond.FindingCarrier` with one method
  `LintFindings() []clierr.Finding`; `internal/cli.lintPayload`
  implements it. `clierr.EmitFindingsText` is the shared renderer both
  branches call, so the failure and success surfaces cannot drift.
- **Ships untested by a REQ-N of its own**: REQ-97/REQ-98 exercise it end
  to end through the CLI in both modes, so the behaviour is covered even
  though the surface is not named.

## D6 — The `ambiguous_match` arm's coverage union excludes ordinary rows

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice recorded)
- **REQ**: REQ-62, REQ-63; SC-19a/SC-19b.
- **Finding**: `0006:C10` requires the union be computed per (group ×
  declared rescuable class) and that the `ambiguous_match` arm be checked
  for a group carrying a `graph-overlap` finding. It does not say WHICH
  rows close that arm. Reading "every row participates" uniformly makes
  the arm vacuously closed for every group whose ordinary rows cover the
  product — which is every overlapping group, since overlapping rows
  cover at least as much as one of them alone. SC-19a's fixture would
  then emit no `graph-coverage-gap`, contradicting REQ-62's own scenario.
- **Evidence**: `0006:C10` states the mechanism — "`escapeOrRefuse`
  selects escapes by `row.rescues(r.Kind)` for the kind that actually
  occurred". The kernel reaches `escapeOrRefuse` on `ambiguous_match`
  PRECISELY because none of the ordinary rows was the exact-one match, so
  an ordinary row cannot rescue an ambiguity it caused. `no_match` is the
  opposite condition — it arises exactly where no row accepts — so an
  ordinary row accepting the assignment is what keeps the refusal from
  arising, and the ordinary population does close that arm.
- **Chosen**: `no_match`'s union is the ordinary rows plus the escape rows
  declaring the class (RDR 0003's `CoverageUnionFor`); `ambiguous_match`'s
  union is the escape rows declaring the class alone.
- **Note**: `internal/guard`'s own `coverageFindings` carries a NOTE
  declining to implement REQ-63's vacuous-closure rule, citing a conflict
  between the two locked records. This RDR owns the graph-lint verdict, so
  the arms are computed here rather than delegated, and `guard.Lint` is
  not called. No RDR 0003 test changes.

## D7 — `graph-overlap` excludes subsumption pairs

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-27, REQ-77.
- **Finding**: two rows whose accepted assignments intersect can be either
  an overlap (blocking) or a redundant row (advisory). The record draws
  the line in the advisory tier's own definition: overlap "is a partial
  intersection between two rows neither of which subsumes the other", and
  a redundant row is one "whose accepted assignments are a *proper
  subset* of a sibling's". The two are therefore disjoint by construction,
  and a pair where one side properly subsumes the other takes the
  advisory code alone.
- **Consequence**: `advisoryBody` / `mvvLegalAdvisory` (a bare row and a
  guarded row in one group) lints CLEAN with `graph-redundant-row`, which
  is what REQ-75 and the MVV's legal advisory matrix require. Two rows
  with EQUAL accepted sets are not a proper subset either way and stay
  blocking overlap, which is what `mvvOrdinaryOverlap` requires.

## D8 — The bare-escape closure is reported whatever else closes the arm

- **Type**: SPEC-UNDER
- **Status**: mechanical translation (derived choice recorded)
- **REQ**: REQ-65, REQ-66; SC-19b.
- **Finding**: `0006:C9` says a group "whose coverage is closed by a bare
  escape row MUST emit `graph-coverage-closed-by-escape`". RDR 0003's
  landed implementation additionally suppresses the report when the
  group's ordinary rows already close the product alone. SC-19b's fixture
  is exactly that shape — `advance-on`/`advance-off` partition `flag` AND
  a bare escape row declares `no_match` — and REQ-63 requires the closure
  code be emitted there.
- **Evidence**: the clause's purpose is stated in the same fence: "a bare
  green MUST NOT satisfy this clause", i.e. the reader must be able to
  tell a group carrying a catch-all from one that does not WITHOUT
  inspecting the model. Suppressing the report exactly where the ordinary
  rows close makes a bare green satisfy it for that group.
- **Chosen**: the finding fires whenever the closing population contains a
  bare escape row. It is `info`, so it never changes the disposition.

## D9 — The reachability successor relation is functional, not per-write

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-101, REQ-108, REQ-110, REQ-112.
- **Finding**: "two edges reaching the same successor produce one node" does
  not fix WHEN two edges reach the same successor. Keying the successor on
  its own value assignment never merges anything (two writes of different
  values are two assignments), which REQ-108 and REQ-110/112 both refute.
- **Evidence**: `0006:LBD` states the traversal is "a fixpoint over merged
  nodes, NEVER path-sensitive", and rejects the path-sensitive reading as
  exponential. The strongest reading of "never path-sensitive" is the one
  under which a path-sensitive enumeration is unrepresentable: every edge
  leaving one node reaches ONE successor node, whose per-tag value sets are
  the join of what those edges produce.
- **Chosen**: a functional successor relation, with the join taking the
  union of per-tag value sets and ABSENCE DOMINATING — a key one path never
  established is absent in the join, since the node stands for every
  concrete view some path reaches it with. Recording it as held would let
  invariant 6 certify a read the runtime finds unavailable, which is the
  false-green direction the record's soundness clause forbids. The fixpoint
  folds a successor into an existing node whenever that node already stands
  for every view it stands for, which is what terminates on cycles.

## D10 — `graph-coverage-gap` attribution splits by arm

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-88, REQ-111, REQ-127; ASSUMPTION-8.
- **Finding**: REQ-111 requires the gap finding NAME the row (`f.Rule ==
  "only-on"`); REQ-88 requires a run mixing rule-id and graph-element-id
  namespaces, and the multi-defect fixture's only non-rule-scoped finding
  is a coverage gap.
- **Chosen**: a gap the group's OWN ROWS leave uncovered names the row; a
  gap that is the ABSENT rescue arm — the group's ordinary rows close the
  product but no escape row declares the class — names the selection
  context as a graph element id and carries no rule, because there is no
  authored row to name for a row that was never written. REQ-127's
  fallback is what keeps the second arm actionable.

## D3 disposition (Phase 2)

The single-valued check IS implemented as a total function over its input
(`internal/graphlint/analysis.go::checkSingleValuedState`): it scans every
row's write block for a single-valued tag assigned more than one value.
The loader still refuses every authorable spelling upstream, so no fixture
reaches it — the entry stays OPEN on RDR 0002's terms, not on this one's.

## D4 disposition (Phase 2)

Discharged. `models/rdr.toml` is authored and homed, `.github/workflows/ci.yml`
carries the `graph-lint` job, and the `Makefile` carries the `graph-lint`
target plus the `build` edge on `check`.
