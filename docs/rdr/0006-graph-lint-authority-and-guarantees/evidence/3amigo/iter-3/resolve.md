Model: claude-opus-5[1m]

# 3amigo Resolve — iteration 3

Origin ledger = `consolidation.md` H1–H8 (multi-persona hotspots) and S1–S16
(single-persona). Every entry exits exactly one way below.

## Hotspots

- **fixed** — origin H1 (advisory tier undefined); section touched:
  `Technical Design` finding-codes table + advisory paragraph, disposition
  table, `Normative Contracts` (new advisory clause), MVV, scenario 14.
  Grounding: `0002::Normative Contracts` routes the unreachable-rule finding
  here by name ("which is RDR 0006's unreachable-rule finding, not a load
  failure here"), so the class was owed, not optional. Minted
  `graph-redundant-row` (accepted assignments a *proper subset* of a sibling's —
  distinguished from overlap, which is a partial intersection and blocking) and
  `graph-unreachable-rule` (no reachable owned-state node satisfies the row,
  including 0002's dead-rule case). Tier closed at three members.

- **fixed** — origin H2 (invariant 5's second conjunct undecidable); section
  touched: `Technical Design` invariant 5. Grounding: the RDR's own reachability
  relation abstracts a node as a per-tag *set* of possible values and merges on
  join, so it cannot answer a per-path question — the personas were right that
  the conjunct is undecidable as written, and checking it at node level would
  flag every legal two-path merge. Resolved by deciding the semantics rather
  than the algorithm: a write to a single-valued tag **replaces**, so no path can
  accumulate a second value and the conjunct is vacuous. Invariant 5 is now a
  per-row syntactic check, stated as such.

- **fixed** — origin H3 (A7 discharges only half of A18); section touched:
  invariant 5, A7 (restated), finding-codes table, `Normative Contracts`,
  scenario 17, Finalization Gate. Grounding: `0003::Normative Contracts` defines
  a conforming view with **two** conjuncts — every always-present key present
  *and* every single-valued tag holding at most one declared value — and
  invariant 5 covered only the second. Added `graph-always-present-owned` for
  the owned half; the observed/recognized half is not model-decidable
  (`resolve.go::assemble` merges them reading no declaration), so A7 now claims
  a **partial** discharge. Also recorded that A7's flip edits `0003`, which is
  `Final [locked 2026-08-22]` — a route-back, not a Draft amendment.

- **fixed** — origin H4 (A6 Pending, behavior in the window unspecified);
  section touched: A6 (`Shape requested`, `Behavior before it lands`),
  invariant 1, disposition table, `Normative Contracts`, scenario 15. Grounding:
  the implementer's contradiction is real — invariant 1's "terminal state named
  by a row" implied named states while invariant 2's "satisfies" implied a
  predicate. RDR 0001 models state as a tag-set and RDR 0002 emits *next-state
  tags*, never a named target, so the predicate reading is the cluster's; the
  named-state phrasing was the outlier and is gone. A rootless model is now
  rejected with a blocking finding rather than reading as clean.

- **fixed** — origin H5 (A5/scenario 7 disjunctive, no oracle); section touched:
  A5 (`Gate target`), scenario 7, Finalization Gate. Grounding: A5's own
  evidence already recorded that `Makefile::check` lacks a `build` edge and CI
  runs discrete jobs never calling `make check` — the facts to decide with were
  present, only the decision was missing. Pinned: a new `lint` job in
  `.github/workflows/ci.yml` running the built binary over the **checked-in
  model** (not the fixture corpus), asserting the JSON `code`, since
  `GroupUserEnv` and `GroupInternal` share exit 2.

- **fixed** — origin H6 (no report-everything duty; overlap cardinality);
  section touched: `Technical Design` emission paragraph, `Normative Contracts`,
  desk-trace step 8, MVV, scenario 13. Grounding: `0003::Normative Contracts`
  carries a clause this RDR consumes and never cited — "Lint MUST report every
  defect it can decide in one pass over a row group, not the first one it
  encounters… each overlapping pair… is its own finding." Adopted by citation;
  it also settles the cardinality question (one finding per pair naming both
  rows, one per shared class for escape pairs).

- **fixed** — origin H7 (two wire shapes; `Finding` type cycle); section
  touched: `Technical Design` envelope paragraph, `Normative Contracts`,
  `Load-Bearing Decisions` wire format, Testing Strategy, disposition table,
  scenario 16. Grounding on `main`: `CLIError`'s optional fields are all
  `omitempty` and `Group` is `json:"-"`; `respond.Success` carries `Data any`.
  The two depths are kept and stated (`error.findings` vs `data.findings`)
  rather than forced into one shape. `Findings` is the one field that must
  **not** be `omitempty`, since the empty list is the receipt. `Finding` is
  defined in `clierr` as a subsystem-agnostic record, so the leaf package gains
  no dependency on the lint package.

- **fixed** — origin H8 (input contract illustrative while name is normative);
  section touched: `Illustrative Code` (now a normative input contract),
  deterministic-order clause, scenario 6. Grounding: RDR 0005 owns `--flow`
  config discovery; nothing owned the pairing. Pinned `--flow` and `--model` as
  mutually exclusive, exactly one model per invocation, `model` field = the
  `[model].id`. Scenario 6 rewritten to assert something that exists at lock —
  one exported engine entry point and one request builder — instead of an alias
  that does not.

## Single-persona

- **fixed** — S1 ("source state" undefined): pinned to the abstract owned-state
  node; group identity is `(node, recognized outcome)`. This is the highest-value
  fix in the pass — it also grounds H4, S3, S4, S6, S9. Grounding: "source state"
  occurs 3× in this RDR, 0× in RDR 0002; `0003::A10` (Pending) explicitly books
  this RDR's confirmation of the division of labour, which the edit now gives.
- **fixed** — S2 (escape class-scoping half not adopted): `0003` states an escape
  row "cannot rescue either class, however bare its guard"; union now computed
  per (group × declared rescuable class). `owned_state_unavailable` recorded as a
  runtime class no lint code mints, explicitly not conflated with invariant 6.
  Scenario 19 added.
- **fixed** — S3 (no join rule / termination): merged-node fixpoint with lattice
  widening; finite lattice, terminates on cycles. Recorded that merging is *what
  makes* the no-false-green guarantee true.
- **fixed** — S4 (predecessors vs. nodes): node form declared normative,
  "reachable predecessors" demoted to descriptive prose; "preserves" defined.
- **fixed** — S5 (`0003::A21` + product bound not carried): both adopted;
  `graph-product-too-large` minted, bound must be published in help output.
  Scenario 18.
- **fixed** — S6 (invariant 1's "transition target"): term removed; RDR 0002
  emits next-state tags, so the destination is checked as tag keys/values.
- **fixed** — S7 (scenario 11's control): control now requires the *explicit*
  always-present marker, since `0003`'s unmarked default is optional.
- **fixed** — S8 (fingerprint unspecified): canonical sortable serialization
  (never a hash), reusing `0002`'s atom sort and `0003`'s set canonicalization;
  namespace tiebreak stated.
- **fixed** — S9 (invariant 2 satisfaction relation): a node satisfies a terminal
  only when **every** value in each per-tag set meets it; partial satisfaction
  rejected, with the reason (the one direction over-approximation must not fail).
- **fixed** — S10 (invariant 7 unobservable): named the input-observable defect
  (a model *relying* on inference — implied terminal, undeclared escape closing
  coverage), so the MVV's mandated fixture can be authored.
- **fixed** — S11 (step 7 precedence untested): scenario 12, asserting the
  negative half (info code absent).
- **fixed** — S12 (text-mode oracle): text mode must enumerate every finding;
  oracle is code-set equality between text and JSON.
- **fixed** — S13 (scenario 3 field presence): dimension arm must name the
  dimension; no arm may carry the code alone.
- **fixed** — S14 (scenario 10 fixture-comment expectation): replaced with the
  paired dead-end case, covering the other half of the over-approximation
  contract.
- **fixed** — S15 ("before use" unenforced): stated that enforcement is the merge
  boundary, and that the kernel deliberately does not consult a lint verdict.
- **fixed** — S16 (no waiver mechanism): decided explicitly — no suppression, no
  allowlist, with the cost named and the escape route (a successor RDR on
  guard-aware pruning) rather than a flag.

Charted to successor: none. Every entry traced to a ledger concern; no net-new
scope absorbed. The two candidates for charting — guard-aware pruning (S16) and
the observed/recognized conformance residue (H3) — were recorded as named future
work inside existing clauses rather than expanding this RDR's contract.

## Needs verification

- **A7** (`Pending`, Peer RDR) — **re-scoped by this pass**, not merely carried:
  now claims a *partial* discharge over owned keys, and its flip requires RDR
  0003 A18 to record two codes plus the open residue. RDR 0003 is `Final`, so
  this is a route-back on a locked peer for Stage 6 to route.
- **A5** (`Pending`, MVV Test) — unchanged in status; its gate target is now
  determinate, so scenario 7 has a mechanical oracle.
- **A6** (`Pending`, Peer RDR) — unchanged in status; the requested *shape* (tag
  predicates, not names) is now stated as a consumer requirement on RDR 0002.
- **New load-bearing claims added by this pass, all carried by existing
  assumptions rather than new ones**: the merged-node fixpoint's no-false-green
  property (a derivation, stated inline); write-replaces semantics for
  single-valued tags (a Design Decision stated in invariant 5, and a consumer
  requirement on RDR 0002's write/clear semantics — Stage 6 should confirm RDR
  0002 does not specify clear-before-write); the published product bound (an
  implementation constant, asserted by scenario 18).

## Tiebreakers

- None. The two forks that looked genuine both collapsed on cluster evidence:
  invariant 5's path-sensitivity question collapsed once write-replaces made the
  conjunct vacuous, and the "source state" three-way ambiguity collapsed once
  RDR 0001's tag-set model and RDR 0002's next-state-tags output were on the
  table — there was never a named state to key on.

## Review gate

Edits stayed brief and carry no change-history narration. Every fix is grounded
in peer-RDR text or `main` source, quoted at the point of use. The
needs-verification list is honest: three Pending assumptions carried, one of
them re-scoped by this pass and now requiring a route-back on a locked peer. No
net-new scope folded in.
