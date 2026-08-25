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

## D11 — The successor join groups on the presence footprint, superseding D9

- **Type**: TEST-FIXTURE
- **Status**: mechanical translation
- **REQ**: REQ-101, REQ-106, REQ-108, REQ-110, REQ-112; supersedes D9.
- **Finding**: D9 chose a FUNCTIONAL successor relation — every edge
  leaving one node folds into ONE successor — reading "never
  path-sensitive" as making a path-sensitive enumeration unrepresentable.
  Phase 3a's FAIL-2 and Phase 3b's ADV-1 independently showed that reading
  unsound. `joinNodes` drops a key absent on either side ("absence
  dominates"), so two rows out of one node writing DIFFERENT owned keys
  annihilate each other: the relation loses both keys and stands for a
  concrete view no path produces. That UNDER-approximates presence, against
  REQ-106's over-approximation contract and REQ-110's premise that a merged
  node admits a SUPERSET of concrete views. The user-visible consequence is
  a false green — live downstream rows are reported `graph-unreachable-rule`
  (advisory), so `checkGroups` skips `checkCoverage` for their groups and a
  whole arm of the model goes unproven at exit 0.
- **Evidence**: REQ-108 licenses the join only between "two edges reaching
  **the same successor**"; a fold across edges that reach DIFFERENT
  successors is not the widening the clause describes. REQ-101 defines a
  node as an abstract owned-state "per owned tag: **absent**, or held with
  its set of possible declared values", which makes the presence footprint
  constitutive of a node's identity rather than incidental to it.
- **Chosen**: group a source node's edges by their successor's PRESENCE
  FOOTPRINT before joining. Edges producing the same held-key set are one
  successor and their per-tag value sets union (REQ-108's merged fixpoint,
  and the shape `TestReq108_...` asserts); edges establishing different
  keys are different successors and stay distinct. Within a group every key
  is held on both sides, so absence-dominates never fires and no reachable
  key is erased. This is still a fixpoint over merged nodes and never
  path-sensitive: a node is keyed by its owned-state, not by the path that
  reached it, so every path arriving at one owned-state arrives at one
  node, and the lattice stays finite so the widening still terminates.
  Successor multiplicity is a property of the transition relation, not of
  path sensitivity — a DFA state has many successors and enumerates no
  paths.
- **Fixture corrected**:
  `TestReq110And112_OwnedSetBeforeMatchReadsMergedFixpointNodes` asserted
  `mids == 1`, counting reachable nodes carrying `status = mid`. Its two
  edges write `{status=mid, opt=p}` and `{status=mid}`, two distinct
  footprints, so the assertion encoded D9's functional reading as its
  oracle. It now counts nodes per OWNED-STATE IDENTITY — the actual
  merged-versus-path-sensitive distinction — and additionally asserts the
  `opt`-absent node exists, since that is the node invariant 6 reads. The
  REQ-112 obligation the test exists for is untouched and still passes: the
  `graph-owned-before-write` finding naming `reads-opt` fires from the
  `{status=mid}` node the `skips-opt` edge genuinely produces. The record's
  own scenario 10 requires exactly that finding and states no node count,
  calling the result "an accepted false positive whose cure is an explicit
  write, clear, or terminal declaration on the model".

## D12 — Invariant 2's outgoing-row test stays bounded by REQ-37

- **Type**: TEST-FIXTURE
- **Status**: needs author decision (implemented on the most defensible
  reading; the suite is green and the REQ-122 census is zero)
- **REQ**: REQ-34, REQ-36, REQ-37, REQ-111, REQ-117, REQ-122.
- **Finding**: Phase 3a's FAIL-4 and Phase 3b's ADV-3 report that
  `checkDeadEnd` splits on terminal-participating keys and then tests
  outgoing rows EXISTENTIALLY on the merged node, so a node multi-valued on
  a key participating in NO terminal is rescued whole by an exit serving
  only one of its values. The miss is real: in the ADV-3 fixture
  `{status=b, phase=q}` satisfies no terminal, has no exit, and lint says
  nothing.
- **Evidence — why the reporters' fix is not available**: REQ-37 states the
  bound as a MUST — "The split is bounded by the declared domains of
  terminal-participating keys only, **never the whole lattice**." The wider
  quantifier was implemented and MEASURED rather than argued away. Ranging
  the outgoing-row test over match-participating keys manufactures concrete
  views the merged node never correlated: in `models/rdr.toml` every
  `stage = "dropped"` write also writes `status = "abandoned"` (seven
  rows, verified), but the merge decorrelates the two keys, so the cross
  product invents `{stage=dropped, status=draft}` and mints NINE false
  `graph-dead-end` findings on the model REQ-122 requires lint clean. The
  record anticipates exactly this: "A merged node holds every value some
  path brings, so a check against it can accuse a path the runtime never
  walks. **A path-sensitive reading would be exponential and is rejected**"
  (LBD, *Join rule and termination*). The ADV-3 node and the conforming
  model's node are structurally identical after the terminal split — one
  singleton terminal key beside one multi-valued non-terminal key — so no
  local rule separates the true positive from the false positives; the
  separation needs the correlation tracking the record rejects.
- **The genuine tension**: REQ-111 ("universal checks must not [read merged
  nodes]") and REQ-117 ("the cure is a clearer model, never a weaker lint")
  favour closing the miss; REQ-37's explicit bound and REQ-122's
  zero-blocking census on the conforming model forbid the only available
  way of closing it. Both readings have record support and the record does
  not rank them.
- **Chosen**: honour REQ-37's explicit bound. It is a stated MUST, whereas
  the wider quantifier is an inference from REQ-111's general principle,
  and adopting it demonstrably breaks REQ-122 on the checked-in model.
  `hasOutgoingOrdinaryRow` keeps its existing scope, with the quantifier
  reasoning and the accepted miss documented at the call site.
- **Fixture corrected**: `TestAdvDeadEndExistentialOnMergedNode` asserted
  the finding must fire. It now PINS the accepted miss — so any later
  widening of the split fails here and forces the REQ-122 census to be
  re-run against `models/rdr.toml` before the widening is accepted — and
  adds a positive arm asserting the bounded split still catches a dead half
  living on a terminal-participating key, so the fixture continues to
  exercise invariant 2 rather than merely recording an absence.
- **For the author**: if invariant 2 should be exact over non-terminal
  keys, the record needs a successor clause that either relaxes REQ-37's
  bound and accepts the false positives on `models/rdr.toml` (re-running
  SC-23's census), or adopts the correlation tracking it currently rejects
  as exponential. This entry is the trigger.

## D13 — Overlap is decided on the projectable dimensions

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-25, REQ-27, REQ-83, REQ-84, REQ-85, REQ-110, REQ-117.
- **Finding**: `emitOverlaps` skipped any row whose accepted-assignment set
  is not projectable. A row carrying one value atom over an OPTIONAL key can
  refuse `guard_unevaluable`, so RDR 0003 gives it no accepted set at all —
  and adding one such atom to each of two genuinely-overlapping rows made
  the `graph-overlap` finding vanish, while the same two rows overlap
  loudly without it (Phase 3a FAIL-3, Phase 3b ADV-2).
- **Evidence**: REQ-84 is explicit — "Withholding a group's exhaustiveness
  claim MUST NOT suppress overlap, coverage, or further withholding
  findings for that group" — and REQ-85 makes them independent findings:
  "each unprovable dimension, each refusing row, each overlapping pair …
  is its own finding". RDR 0003 already makes the same move one level down:
  `acceptedIn` computes over the group's DECIDABLE sub-product because "an
  atom lint cannot project belongs to the row that carries it, not to every
  row sharing its group".
- **Chosen**: `decidableAccepted` re-expresses a non-projectable row by
  intersecting its projectable atom denotations (`guard.Denotation`) into
  the group's scoped product. Where a row projects whole this is
  `guard.AcceptedAssignments` unchanged, so the provable path is untouched.
  The relaxation is scoped to OVERLAP and sound only because overlap is
  EXISTENTIAL (REQ-110): dropping an undecidable conjunct WIDENS the row,
  so a witness found here may be separated by the dropped dimension at
  runtime — a false positive, never a missed defect, which is the direction
  REQ-117 accepts. It is deliberately NOT reused for coverage, which is
  universal and where a widened union is the false-green direction;
  coverage still compares the full product against `guard.CoverageUnion`.
  A row carrying an `unless` block is left undecided rather than decided
  wrongly, since the block matches as ONE conjunction and cannot be
  partially subtracted. `internal/guard` is RDR 0003's surface and was not
  modified.

## D14 — Invariant 7 is decided per node

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-43, REQ-36, REQ-83; REQ-MVV.
- **Finding**: `checkTerminalEscape` opened with
  `if len(a.model.Terminal) > 0 { return }`, so `graph-terminal-escape`
  could fire only on a model declaring no terminal WHATSOEVER. Every model
  declaring even one terminal was exempt, however many of its reachable
  nodes relied on an inferred one (Phase 3a FAIL-1, which had no test).
- **Evidence**: REQ-43 states the condition per NODE — "a reachable
  non-terminal node with no outgoing non-escape row and **no terminal
  declaration covering it** (an implied terminal)" — a per-node coverage
  test, not a per-model emptiness test. The record also prescribes the
  fixture shape: "a fixture is authored by omitting the declaration the
  model depends on", i.e. omitting ONE terminal. Under the global gate that
  fixture could not mint the code, and REQ-MVV's illegal matrix entry was
  satisfiable only by the degenerate no-terminal model.
- **Chosen**: walk per-node terminal coverage, testing the SPLIT node for
  the same reason invariant 2 splits (REQ-36) — terminal satisfaction is
  universal over a node's per-tag value sets, so a merged node is
  anti-monotone in it and would accuse converging flows.
  `satisfiesSomeTerminal` returns false for an empty terminal list, so the
  no-terminal-at-all case is unchanged and the pre-existing REQ-43 test
  still passes. The reliance stays ONE finding against the missing
  declaration rather than one per node, as before. `graph-dead-end` and
  `graph-terminal-escape` are separate obligations over the same node
  condition and both are emitted, per REQ-83's complete-emission clause.

## D15 — A tag with no finite domain is carried as the opaque value

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-101, REQ-106, REQ-109; D11 precedent.
- **Finding**: `OpaqueValue` was declared in `reach.go` and READ by
  `ownedAtomSatisfiable`, but never assigned anywhere — root construction
  and `successor` both called `canonicalValues`, so a tag with no finite
  declared domain (a `scalar`, or an `int` lacking `min`/`max`) was carried
  as the concrete string the author happened to write. A value atom over
  such a tag was then decided against that one string and PRUNED whenever
  it named anything else, so the edge behind it was never taken. The
  user-visible consequence is the same false green D11 recorded on the
  presence axis: an owned `scalar` initialized to `"start"` gating an edge
  on `free eq "other"` yields three advisory `graph-unreachable-rule`
  findings at exit 0, while the duplicate rows behind the gate — a
  BLOCKING `graph-overlap` — go unchecked. Changing the gate literal to
  the initial value surfaces `graph-overlap` and `graph-coverage-gap` at
  exit 2, so the defect was real and only the abstraction hid it.
- **Evidence**: REQ-101 and REQ-109 both state the abstraction verbatim —
  a tag with no finite domain abstracts to held/absent — which makes its
  held half a single representative rather than an enumeration over
  authored strings. REQ-106 fixes the direction: the traversal
  over-approximates the runtime, and pruning an edge lint cannot decide is
  the false-green direction. Nothing adjudicates this away; the string
  "opaque" appears in no RDR 0006 artifact, so the constant was written to
  the record's intent and simply never wired up. `ownedAtomSatisfiable`'s
  `slices.Contains(held, OpaqueValue)` arm already implements the reading
  and needed no change.
- **Chosen**: `heldValues` normalizes a value a node comes to HOLD,
  returning `[]string{OpaqueValue}` when `guard.AssignmentCount` reports no
  finite domain and `canonicalValues` otherwise. It is called at the two
  sites that INTRODUCE a held value — the root's initial assignments and a
  row's writes in `successor`. `guard.AssignmentCount` is the one authority
  on finiteness, so the abstraction is gated on the DECLARATION rather than
  on the declared kind: an `int` carrying a bound stays finite and keeps
  its concrete value, exactly as an enum does. `joinNodes` is deliberately
  untouched — it unions sets both of whose sides already came through
  `heldValues`, and the opaque value is idempotent under that union, so a
  finite tag never mixes with it. `internal/guard` is RDR 0003's surface
  and was not modified.
- **Premortem**: the wider relation could swallow a universal finding
  (`graph-unreachable-rule`, `graph-dead-end`) on a genuinely dead
  scalar-gated arm. That is the direction REQ-106 mandates and the one
  `atomAdmitsValue`'s UNEVALUABLE arm already takes; coverage stays honest
  because `coverage.go`'s `ReasonDimensionNotFinite` arm gates on the very
  same `guard.AssignmentCount` predicate, so a dimension abstracted here is
  still reported unprovable there. The REQ-122 census over `models/rdr.toml` remains zero
  findings after the change.

## D16 — Composite keys escape their fields before joining

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-87, REQ-106, REQ-108, REQ-114; RDR 0002 REQ-126, REQ-127.
- **Finding**: three composite keys joined authored strings on unescaped
  structural delimiters, and each collision is false-green. (1)
  `presenceKey` appended `k + ";"`, so the footprints `{"a;b"}` and
  `{"a", "b"}` both render `a;b;`. That key is the successor-join group id
  in `successorsOf`, so two DIFFERENT owned-states land in one group and
  `joinNodes`'s absence-dominates arm erases every key held on only one
  side: a fixture declaring owned tags `a;b`, `a`, `b` and two rows writing
  `{status=b, "a;b"=l}` and `{status=b, a=l, b=l}` yields the reachable set
  `{status:[a]}, {status:[b]}` — all three owned keys gone, so invariant 6
  will certify a read of a tag the runtime does hold. (2) `Node.key`
  appended `k + "=" + join(vals, ",") + ";"`, so `{a: "b;c=d"}` and
  `{a: "b", c: "d"}` both render `a=b;c=d;`. That key is the `seen`
  de-duplication identity BOTH terminal walks in `analysis.go` close on and
  the fixpoint equality in `reach`, so the collision emits ONE
  `graph-dead-end` where two are owed and can terminate the widening early.
  (3) `Fingerprint` and `compareAtoms` joined atom fields and set members on
  `,`/`;`/`|`/`#`, so a guard `in` over `["a,b", "c"]` and one over
  `["a", "b,c"]` both render `mark|all|in|a,b,c;` — two rows, one finding
  identity.
- **Evidence**: reachability is not hypothetical. `internal/table` imposes
  NO charset on a tag name or a set member — `tagDecl` in
  `internal/table/load.go` validates provenance, kind, per-kind fields, and
  bounds only — and a fixture declaring a tag literally named `a;b`
  (TOML-quoted key `["tags"."a;b"]`) loads clean. RDR 0002 REQ-126/127 name
  case (3)'s pair verbatim as a MUST: `["a,b", "c"]` and `["a", "b,c"]`
  "MUST yield two atoms in the normalized set", and the loader's own
  fixtures already ship them at `internal/table/testdata/delim/`. REQ-114
  closes the finding-identity tuple on the fingerprint, so (3)'s collision
  costs determinism as well as separation. REQ-106 fixes the direction for
  (1) and (2): the traversal over-approximates the runtime, and a rendering
  that merges two owned-states under-approximates it.
- **Chosen**: one unexported `escapeField` beside `canonicalValues` in
  `reach.go`, applied to every field of all three composites — `presenceKey`
  escapes the tag name, `Node.key` escapes the key and each value,
  `Fingerprint` escapes key, block, operator, each literal member, and each
  next-state key and value, and `compareAtoms` escapes through the same
  `literalKey` the fingerprint writes so the sort key and the rendering can
  never diverge. The encoding is backslash escaping: `\` → `\\` first, then
  each of `;` `=` `,` `|` `#`. Canonicalization runs BEFORE escaping, so
  `canonicalValues` still sorts and compacts the AUTHORED values and the
  canonical form stays a property of what the author wrote.
  `internal/table` was not touched: a charset restriction on tag names is a
  loader decision RDR 0002 owns, and lint must be sound over every model the
  loader accepts today.
- **Premortem**: REQ-87 requires the fingerprint be "readable, never a hash,
  sortable", which is what ruled out the two alternatives. A length prefix
  is injective but destroys lexicographic sortability outright — `10:` sorts
  before `2:` — so a sort over fingerprints would stop being meaningful.
  `encoding/json` is injective and reversible but reorders nothing legibly
  and needs its own sortability argument. Backslash escaping keeps all
  three: it is injective (prefix-free and reversible), readable (`a\;b;` is
  eyeball-distinguishable from `a;b;`), and sort-preserving, because `\`
  (U+005C) sorts above every delimiter (the highest is `=` at U+003D), so
  escaping only ever moves a delimiter-bearing member LATER and never
  transposes two members already correctly ordered. The risk that remains is
  a fourth renderer added later that joins without escaping; `escapeField`
  and `escapeJoin` are the one escape in the package precisely so a sort key
  and the composite it orders cannot diverge. The REQ-122 census over
  `models/rdr.toml` remains zero findings after the change, and REQ-85's
  order-independence assertion is unaffected — no value in that fixture
  carries a delimiter, so every escaped rendering is byte-identical to the
  old one.

## D17 — The escape carries injectivity; `compareAtoms` carries the order

- **Type**: IMPL-DECISION
- **Status**: mechanical translation
- **REQ**: REQ-87, REQ-106, REQ-108, REQ-114; RDR 0002 REQ-126, REQ-127,
  `0002:C19`.
- **Finding**: two defects in D16's encoding, both found by review of that
  record's own premortem.

  (1) D16's premortem claims backslash escaping is sort-preserving because
  `\` (U+005C) sorts above every structural delimiter, so escaping "only
  ever moves a delimiter-bearing member LATER". That reasoning compares the
  escape against the delimiter it replaces, but a sort compares against
  whatever the NEIGHBOUR carries at that position. `["a,b"]` precedes
  `["a0"]` as authored — `,` (U+002C) < `0` (U+0030) — and follows it once
  escaped, because the inserted `\` (U+005C) now occupies that position and
  U+005C > U+0030. The two transpose. `compareAtoms` had been changed to
  sort on `literalKey`, the escaped rendering, so a row's atoms were ordered
  by an artifact of the encoding rather than by RDR 0002's canonical order —
  which is precisely the order REQ-87 names as the one the fingerprint must
  be in.

  (2) `escapeJoin` joined members on the separator, which is not injective
  over member-sequence CARDINALITY: `[]` and `[""]` both render the empty
  string, so `k=;` stood for two different states. Both are authorable — a
  `set` declaring no `elements` constrains no member (`conformDomain` guards
  on `len(d.Elements) > 0`), and `conform` loops zero times over an empty
  member sequence, so `write = { k = [] }` and `write = { k = [""] }` both
  load. They are not the same state: `[]` satisfies no value atom while
  `[""]` satisfies `eq ""`.

- **Decision**: separate the two concerns D16 had conflated onto the escape.
  `escapeField`'s contract is INJECTIVITY only, and its doc comment now says
  so rather than claiming an order guarantee it does not have.
  `compareAtoms` compares the literal element by element over the canonical
  AUTHORED members — `slices.Compare(canonicalValues(a.Literal), ...)` —
  which is what `internal/table.compareAtoms` does and what `0002:C19`
  defines as normative. The atoms are therefore ordered BEFORE any escaping,
  and the escape is an encoding applied after the order is fixed.
  `escapeJoin` TERMINATES each member with the separator instead of joining
  on it, so `[]` renders `` and `[""]` renders `,`.
- **Premortem**: the terminator keeps everything REQ-87 asks. It stays
  readable — `a,b,` is as legible as `a,b`, and the written value is still
  eyeball-visible in the next-state half — and it is never a hash. It is
  injective over cardinality as well as content, because every member is
  followed by its terminator, so no member sequence's rendering is a prefix
  of a different one's. Sortability is no longer asked of the encoding at
  all, which is the point: the sort is over authored members and cannot be
  perturbed by a later change to the escape.

  The reachability of the two arms differs, and the fix is scoped to the one
  that is live. `Reach` abstracts an unconstrained `set` to `OpaqueValue`
  under D15 — a `set` with no `elements` carries no finite domain — so the
  `[]`/`[""]` pair never reaches `Node.key` or `presenceKey`; and a `set`
  WITH `elements` is finite but cannot admit `""`, which `conformDomain`
  refuses. `Fingerprint` has no such abstraction between the authored value
  and the rendering: it closes over the row's raw `NextTags`, so that is
  where the collision was reachable and where the regression test pins it.

  Both regression tests were checked against the pre-fix code and both fail
  there, so neither is vacuous. The REQ-122 census over `models/rdr.toml`
  remains zero findings. One existing assertion in
  `TestReq87_FingerprintIsACanonicalSortableSerializationNeverAHash` moved
  from `status=c;` to `status=c,;`: it asserts the written value is READABLE
  and not a hash, which the terminator preserves — the expected substring
  had incidentally encoded the old infix-join form.

  **Scope of the order guarantee.** REQ-87's "canonical sortable
  serialization" is about the atoms WITHIN one fingerprint, and that is
  what `compareAtoms` now delivers. It is NOT a claim that two fingerprints
  sort against each other in authored-literal collation: `identityKey`
  sorts findings on the fingerprint STRING, so two identities differing by
  the transposing pair (`a,b` vs `a0`) still compare in escaped order. No
  requirement asks otherwise. REQ-86 orders findings by the REQ-114 tuple —
  model id, invariant code, rule/context or element id, THEN fingerprint —
  where the fingerprint is the final tie-break and what is required of it
  is injectivity and determinism, both of which hold. REQ-89 asks for
  determinism asserted as full-list golden equality, not a collation
  property, and REQ-88's cross-finding clause is about identity NAMESPACES
  rather than literals. Making the emitted order track authored collation
  would mean either the encodings D16 already rejected on REQ-87, or
  sorting on a structured key that is not the readable fingerprint — which
  would weaken "never a hash" with no clause demanding it.
