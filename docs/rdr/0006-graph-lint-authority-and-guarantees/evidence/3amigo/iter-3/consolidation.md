Model: claude-opus-5[1m]

# 3amigo Consolidation — iteration 3

Three isolated persona passes (PM 8 findings, Implementer 15, QA 14). Isolation
verified: no file references another persona's output. Hotspots below are
computed mechanically — a passage two or more personas named independently. Overlap
marks a hotspot passage, not a validated finding; a single-persona finding is not
thereby weaker.

Origin ledger for this iteration = H1–H8 plus the single-persona entries S1–S14.

## Hotspots (multi-persona passages)

**H1 — Advisory tier: two of three members have no code and no definition.**
`Technical Design` "the advisory tier is scoped to redundant rows, unreachable
rules, and the bare-escape coverage closure" + finding-codes table (only
`graph-coverage-closed-by-escape` is `info`) + disposition row "Redundant row /
unreachable rule".
PM 2, Implementer 13, QA 1 — **all three personas**. Highest-consensus passage in
the pass. PM adds the peer-delegation ground: `0002::Normative Contracts` routes
the unreachable-rule finding *here* ("a rule whose predicate set requires
`recognized` to be absent is dead, which is RDR 0006's unreachable-rule finding,
not a load failure here"). QA adds that MVV scenario 1 asserts `findings` content
on the legal fixture, so a test must name codes that do not exist. Neither
"redundant row" nor "unreachable rule" is defined.
Blocks: the informational half of the taxonomy; MVV scenario 1's assertion set.

**H2 — Invariant 5's second conjunct is undecidable under this RDR's own
reachability abstraction.**
`Technical Design` invariant 5, "no reachable path leaves a second value held
beside one no write or clear on that path removes", against `Load-Bearing
Decisions` → reachability relation ("held with its set of possible declared
values").
Implementer 3, QA 2. The node abstraction merges paths by construction, so it
cannot answer a per-path question; a node abstracting to `{a,b}` is the normal
representation of a legal two-path merge, not a violation. Implementer adds the
unstated write-replaces-vs-clear-first question, under which the conjunct may be
vacuous. Both note the first conjunct is a trivial per-row syntactic check that
satisfies scenario 2 and the MVV verbatim while omitting half the invariant.
Blocks: whether Phase 2 implements a second traversal; the fixture for the
path-accumulation half.

**H3 — A7 discharges RDR 0003 A18 with no soundness argument, and only half of
A18's conformance definition.**
`Critical Assumptions` A7; `Technical Design` invariant 5's "discharges that
premise on its own".
Implementer 14, QA 3; PM 8 hits the same assumption from the flip-condition side.
`0003::Normative Contracts` defines a conforming view as one where **every
always-present key is present** *and* **every single-valued tag holds at most one
declared value**. Invariant 5 addresses only the second conjunct. QA adds that
sufficiency is a soundness claim (model-level clean ⇒ no non-conforming view
assemblable) that is asserted and never argued, with no scenario in the direction
that matters. PM adds that A7's flip condition requires an edit to RDR 0003,
which is `Final [locked 2026-08-22]`.
Blocks: invariant 5's scope; whether A7 can flip honestly.

**H4 — A6's declarations are Pending and nothing specifies lint's behavior in the
window before they land.**
`Critical Assumptions` A6; `Implementation Plan > Prerequisites`; `Capability
Dependencies` row "Initial owned state and terminal declarations".
PM 3, Implementer 5, QA 5 — **all three personas**. PM frames the outcome cost:
invariants 2, 6, 7 inert while `respond.OK` with empty findings is defined as
"the proof's receipt" — a receipt for a proof never attempted. Implementer needs
the *shape* (tag=value set in `[model]`? `[initial]` table? rule id? named states
vs. a predicate) and notes invariant 1's "terminal state named by a row must
resolve" implies named states while invariant 2's "satisfies no declared
terminal" implies a predicate — contradictory. QA needs the negative test for a
rootless model, for which the document supports three different outcomes.
Blocks: Phase 2 entirely; Phase 4 fixtures; the root-absent disposition.

**H5 — A5 / scenario 7: the CI gate has a disjunctive target and no mechanical
oracle.**
`Critical Assumptions` A5; Validation scenario 7; `Implementation Plan` Phase 3.
PM 4, Implementer 7, QA 4 — **all three personas**. PM: "checked-in transition
model **or** fixture corpus" lets the gate be satisfied without ever linting the
production model — the user outcome A5 exists to guarantee. Implementer: A5 says
assert on `Code`, scenario 2 asserts code *and* exit 2, and `Group` is
`json:"-"`, so a gate reading `$?` cannot distinguish lint failure from an
internal error; four candidate gate surfaces, none designated. QA: A5 is the
RDR's only `MVV Test` assumption and flips to Verified on a scenario whose pass
condition is a human reading two files.
Blocks: Phase 3 CI wiring; A5's verification.

**H6 — Overlap/emission cardinality is undefined, and RDR 0003's
report-every-defect clause is not adopted.**
`Technical Design` "Each finding carries a stable code…"; disposition table (one
"Finding minted" cell per class); `Load-Bearing Decisions` → Identity; desk-trace
steps 4 and 7.
PM 1, QA 9; Implementer 11 overlaps on the identity tuple's tiebreak.
`0003::Normative Contracts` carries a clause this RDR is the consumer of and
never cites: "Lint MUST report every defect it can decide in one pass over a row
group, not the first one it encounters… each unprovable dimension, each refusing
row, each overlapping pair, and any coverage gap over a provable product is its
own finding." That clause also settles QA's cardinality question. As written the
disposition table reads as one verdict per class and a first-failure engine is
conformant.
Blocks: collect-all vs. first-failure engine; `assert len(findings) == N`; whether
the MVV needs a multi-defect fixture (it currently has one defect class per
fixture).

**H7 — The success/failure `findings` wire shape is specified two ways.**
`Technical Design` "extend `internal/cli/clierr.CLIError` with an optional typed
`Findings` field… On success the `respond.OK` payload carries the same typed
`findings` list"; disposition row "`findings` present and empty"; `Validation >
Testing Strategy` naming `respond.go::Success`'s `Data any` as the host.
Implementer 6, QA 11. Two different wire shapes are live
(`{"type":"ok","data":{"findings":[…]}}` vs `{"type":"ok","findings":[…]}`), and
"present and empty" as the proof's receipt collides with `omitempty` on the
failure envelope and with an untyped `Data any` on the success one. Implementer
adds type ownership: `clierr` is documented as the leaf package other internal
packages import without pulling `internal/cli`, so hosting a graph-lint `Finding`
type there inverts that, and defining it in the lint package makes `clierr`
import it — a cycle.
Blocks: the Phase 3 envelope change; every MVV JSON assertion; golden-JSON tests.

**H8 — Command input contract is illustrative while the command name is
normative.**
`Illustrative Code` (`--flow rdr --model ./path`); `Proposed Solution` command
authority; Validation scenarios 6 and 7.
PM 6, Implementer 10; QA 14 overlaps on scenario 6's untestable alias-equivalence.
Since `lint` is deliberately outside RDR 0005's `flow` group, the flag vocabulary
has no owner: `--flow` (RDR 0005's config discovery) vs `--model` (explicit path),
both-given behavior, multi-model/corpus support (scenario 7 says "corpus", and the
ordering key leads with model id, implying many), and whether the finding's
`model` field is `[model].id` or the path. QA adds that scenario 6 names no symbol
a test could assert against and is vacuously satisfiable at lock.
Blocks: the Phase 3 command signature; every MVV invocation string.

## Single-persona findings

Not weaker for being single-persona — these were simply reached from one angle.

- **S1 (Impl 1) — "source state" is undefined, yet it keys every row group.**
  `Technical Design` "Row groups are RDR 0003's" adopts "the same source state and
  the same recognized outcome". RDR 0002's authored shape has no `state` field, so
  source state must be derived and no clause says from what; three incompatible
  readings are live (match-pattern over owned tags / the abstract reachability node
  / a distinguished status-like tag). `0003::A10` (`Status: Pending`) explicitly
  books RDR 0006's confirmation of this division of labour as still open.
  Blocks: the group-keying function — the load-bearing data structure of Phase 2.
- **S2 (Impl 2) — RDR 0003's class-scoping half of the escape-row clause is not
  adopted; `owned_state_unavailable` appears nowhere in this RDR.**
  `0003` states an escape row "cannot rescue either class, however bare its guard"
  (`guard_unevaluable`, `owned_state_unavailable`). Invariant 4 adopts the
  participation half; whether a `no_match`-only escape row closes coverage for the
  whole group or only that arm is unstated.
  Blocks: whether the union is per-group or per (group × rescuable class).
- **S3 (Impl 4) — the reachability fixpoint states no join rule or node identity
  over cyclic graphs.** Union-on-merge (widening) vs. distinct nodes
  (path-sensitive, exponential on the self-loops RDR 0002 models) give different
  answers for invariant 6, and the "never a false green" guarantee holds only under
  one. Blocks: the Phase 2 worklist algorithm.
- **S4 (Impl 8) — "reachable predecessors" (per-path, over rows) and invariant 6
  ("every reachable owned-state") are different questions with different complexity
  classes.** Also unstated: what "preserves" means operationally.
- **S5 (Impl 9) — `0003::A21` makes the single-valued marker a "precondition for
  provability" and `0003` requires a "declared, model-independent bound" the
  implementation MUST publish for too-large products.** Neither is carried here; no
  code is minted for the bound case.
- **S6 (Impl 12) — invariant 1's "transition target" has no counterpart in RDR
  0002's normalized row** (next-state tags and writes, not a named target) — either
  redundant with the "tag" term or dependent on A6's unwritten schema.
- **S7 (Impl 15) — scenario 11's "declared always-present" is a third marker state**
  against `0003`'s assignment-count table, where an *unmarked* key takes the ×2
  presence row (unmarked defaults to optional).
- **S8 (Impl 11) — the deterministic-order key bottoms out in an unspecified
  "normalized predicate/write fingerprint"** (hash vs. canonical string; no tiebreak
  across the rule-id/graph-element-id disjunction scenario 5 produces).
- **S9 (QA 6) — invariant 2 defines no satisfaction relation between an abstract
  node and a declared terminal** — exactly the case the over-approximation produces.
- **S10 (QA 7) — invariant 7 is a prohibition on lint's own reasoning, not an
  input-observable defect**, yet the MVV mandates a `graph-terminal-escape` illegal
  fixture that therefore cannot be authored.
- **S11 (QA 8) — desk-trace step 7's precedence rule has no scenario**; 8 and 9c each
  test one half, and nothing asserts the negative (info code absent from the failure
  payload).
- **S12 (QA 10) — scenario 4's text-mode "equivalent semantics" is undefined** while
  Technical Design makes text findings permissive (`may`).
- **S13 (QA 12) — scenario 3's non-finite-dimension finding requires no fields**;
  atom attribution is left to the implementation.
- **S14 (QA 13) — scenario 10's second expectation is a statement about fixture
  comments**, and the dead-end half of the over-approximation contract has no
  scenario.
- **S15 (PM 5) — "rejected before it can be used by the resolver" has no enforcer**:
  nothing couples `flow resolve` to a lint verdict, so "before use" is really
  "before merge".
- **S16 (PM 7) — over-approximation false positives are accepted with no
  suppression, waiver, or sizing**, so the lint can dictate model structure to work
  around its own imprecision.

## Review gate

**Healthy.** Every finding is anchored to a named passage and names the decision
or test it blocks. No generic advice, no unanchored findings, no isolation leak.
The three-persona convergence on H1, H4, and H5 is genuine agreement between
contexts that never saw each other.
