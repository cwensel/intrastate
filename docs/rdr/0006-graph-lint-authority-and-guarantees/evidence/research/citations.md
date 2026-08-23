# RDR 0006 — Stage 4 Resolve research/verification record

Model: claude-opus-5

Date: 2026-08-23
Scope: SCOPED RE-ENTRY per `Status: Draft [demoted from Final 2026-08-21 —
re-verify A2, A5]`, plus the anchors the demotion's refine (`f69802e`) touched:
A3 (evidence rewritten to quantify over the new reachability relation) and the
new A6 the refine introduced.

Verification was delegated per stage doctrine; verdicts held here, edits
authored in the main context.

## A2 — predicate lint can decide overlap and coverage (Peer RDR)

**Verdict: VERIFIED.** All eight sub-claims of A2's evidence line check out
verbatim against RDR 0003.

| Sub-claim | Verdict | Anchor |
| --- | --- | --- |
| coverage `union(row_i accepted assignments) == scoped product`; overlap = non-empty intersection over enum/boolean, declared set-universe, bounded-int | TRUE | `0003::Critical Assumptions` A2 |
| finite domains required; blocking inability-to-prove outcome | TRUE | `0003::Normative Contracts` "cannot carry an exhaustiveness claim, and lint MUST take the blocking inability-to-prove outcome" |
| tag declaration model = value kind, finite domain, optionality, single-valued marker, element universe | TRUE | `0003::Normative Contracts` "MUST carry a value kind, and … MAY carry a finite domain, an optionality marker, a single-valued marker" |
| scoped row group = rows sharing source state + recognized outcome | TRUE | `0003::Technical Design` "**The row group is defined here**, not deferred" |
| closes 0003 A10 and A12 | TRUE | `0006::Normative Contracts` row-group clause; `0006::Load-Bearing Decisions` reachability bullet |
| escape-row participation clause quoted verbatim | TRUE | `0003::Normative Contracts` "A declared escape row participates in the coverage identity" |
| narrowing clause quoted verbatim; "can refuse" is syntactic over optionality | TRUE | `0003::Normative Contracts` "An exhaustiveness claim MUST NOT be stronger than the runtime it describes" |
| coverage claim default-on, never opt-in | TRUE | `0003::Technical Design` "default-on for every scoped row group whose participating dimensions are all finitely declared" |

§JD-13 landed: `single-valued` now occurs 71x in RDR 0003 with its own normative
clause; RDR 0002 carries it as an authoring-location mention citing 0003 as the
normative home. The field 0006's A2 depends on has a producer.

A12's graph-vs-syntactic distinction is explicitly drawn by the consumer:
`0006::Load-Bearing Decisions` — "It is a different predicate from RDR 0003's
'can refuse' test, which is decided over the optionality field alone and never
consults this relation."

## A5 — CI can run `intrastate lint` as the blocking authority (MVV Test)

**Verdict: gate surfaces VERIFIED; two evidence-line imprecisions found.**

| Claim | Verdict | Anchor |
| --- | --- | --- |
| `Makefile::check` is the local aggregate gate | TRUE, with caveat | `check: fmt-check vet lint test` — does **not** depend on `build` |
| `Makefile::build` builds `./bin/intrastate` from `cmd/intrastate` | TRUE | `BIN := $(BIN_DIR)/intrastate`, `PKG := ./cmd/intrastate` |
| `.github/workflows/ci.yml::jobs` runs gates on push and PR | TRUE, imprecise | jobs `test`/`lint`/`vuln` run discrete targets; CI never invokes `make check`; `push` restricted to `main` |
| `internal/cli/root.go::NewRootCmd` registers through `ExecuteAndEmit` | TRUE | `cmd.AddCommand(newVersionCmd())`; `ExecuteAndEmit(NewRootCmd(), os.Args[1:])` |

Corroborating facts for the implementation:

- No `intrastate lint` verb exists; `newVersionCmd` is the only registered verb.
- No `Findings` field exists in `clierr`. `CLIError` carries
  `Code/Message/Param/Detail/Hint/Group/Cause`; all optional fields use
  `omitempty`, so the append-only typed field 0006 requires is wire-compatible.
  `respond.Success` already carries `Data any` — a host on the success side.
- `clierr::GroupUserEnv` exists; `ExitCodeFor` maps `GroupUserEnv, GroupInternal
  → 2`, matching Validation scenario 2's asserted exit code. Note exit 2 is
  shared with `GroupInternal`, so a CI gate should assert on `Code`, not the
  exit integer alone.

A5 remains `Pending` by `MVV Test` as designed — the gate wiring is proven
*possible*, not *done*.

## A6 — declared initial owned state and terminal states (new at refine)

**Verdict: NOT VERIFIABLE — genuine producer gap; author decision required.**

- `initial` occurs zero times as a state concept in RDR 0002, 0001, and 0003.
- `terminal` occurs exactly twice in RDR 0002, both naming the fixture rule id
  `terminal-archive` in Validation scenario 2 — confirming A6's own claim.
- Live fixtures corroborate: both `[model]` blocks in
  `0002-…/evidence/spikes/rdr-fixture.toml` and `kata-fixture.toml` carry only
  `id`, `version`, `description`. The single `terminal` token is a rule id
  (`id = "terminal-archive"`) and a source label (`source = "rdr:terminal"`).
- No checked-in production transition model exists anywhere in the repo.

RDR 0002's schema is closed by two normative clauses that a route-(a) fix must
amend together: the layout enumeration (`root outcomes, [model], [tags.<tag>],
[accessors.<id>], [context.<id>], [[rule]], [dump]`) and "`[model]` MUST contain
`id` and `version`".

Structural note: RDR 0001's kernel is stateless and single-step — it never
traverses a graph. Multi-step reachability exists only in this RDR's lint, which
is why no peer declares a root today.

## Reuse audit (`$RDR_ENV` reuse-audit paths)

No existing capability duplicates what this RDR introduces. `internal/` +
`cmd/` hold 8 non-test Go files; the only registered verb is `version`. Searches
for a graph, lint verb, invariant engine, or reachability dataflow return
nothing — the two `graph|invariant|reachab` hits are incidental comments
(`root.go:116` "never-silent invariant", `resolve.go:404` "Unreachable: pruned
above"). Nothing to fold in; the introduction stands.

## Code-vs-spec state (context for lock decisions, not an A-record)

`internal/resolve/resolve.go::Row` still declares `Guard string`; `Refusal`
still declares `Guard string`; no atom type (`Key`/`Operator`/`Literal`/`Block`)
exists in the repo. The last commit touching `internal/` or `cmd/` is `1c9c0ca`;
162 commits have landed since, all specification. This RDR's finding contract
cites the atom shape as RDR 0007's — normative in the spec, absent from code.

## Cluster survey (0001-0009 + JDR 0001) — cohesion check

Run because the author asked for a cluster-wide review before choosing an A6
arm. Verdicts below are the load-bearing ones; they change the recommendation.

### Cluster status

| RDR | Status | Assumptions (V / P / other) |
| --- | --- | --- |
| 0001 | Implemented | 4 V, 1 Accepted |
| 0002 | **Draft** | 8 V, 0 P |
| 0003 | Final (locked 2026-08-22) | 13 V, 8 P |
| 0004 | Final | 8 V, 2 P (DOWNGRADED) |
| 0005 | Final | 5 V, 1 Accepted |
| **0006** | **Draft** | 4 V, 2 P (A5, A6) |
| 0007 | Final | 21 V, 4 DOWNGRADED |
| 0008 | Final | 12 V, 1 Refuted-as-stated |
| 0009 | Final | 10 V |

### The atom-shape citation is sound — earlier concern retired

RDR 0007 is `Final` and **fully absorbed JDR 0001 §D1 and §D4**, re-entering at
propose (its evidence tree carries `iter-2` under `propose-premortem/`, `cove/`,
`3amigo/`). It normatively fixes the four field names this RDR cites:

> The atom's four fields are named `Key`, `Operator`, `Literal`, `Block`
> (`0007::Normative Contracts`, SEAM clause)

So 0006's `Method: Peer RDR` for the guard atom shape is well-founded, and the
§JD-4 consequent duty (carry the refusing atom) is discharged in 0006's
normative finding contract. 0002 absorbed §D2 (gate-then-count, with the
Scenario 4 unevaluable-sibling qualifier) and §D4's normalizer duties; 0004
absorbed §D3 (read completeness clause + MVV partial-read refusal).

### Every inbound obligation on 0006 is discharged

§JD-4 (citation + atom field + no new tier), §JD-13 (single-valued citation),
§JD-14 (invariants 3/4 repair), and RDR 0003's A10, A12, A17, A19, A20 all
appear discharged in 0006's current body. RDR 0003 A18 is discharged at the
**model level** (invariant 5, single-valued half) with the view-level half
explicitly routed to RDR 0007 — which is one of the two venues 0003 A18 itself
names ("either discharges it, both is better").

### The one asymmetry

0006's only unmet **outward** request is A6 → RDR 0002. The two `Draft` nodes
in the cluster (0002 and 0006) are mutually gating on the reachability root:
0006 re-locking is what retires five of 0003's eight open records, but 0006
cannot close A6 unless 0002 declares the root and stop set.

### Open JDR entries vs 0006

- **§JD-5** (precondition precedence) — does not touch 0006; no precondition
  ordering invariant exists here.
- **§JD-3** (`RequiresOwned` producer) — `RequiresOwned` occurs zero times in
  0006; invariant 6 reasons over declared writes and this RDR's own
  reachability relation, not that field.
- **§JD-10** (recognized-tag totality) — the matching-time half is already
  answered in `0002::Normative Contracts` ("The `recognized` tag is total over
  matching"). 0006 is silent, not contradicted; a citation duty at most.

### Structural fact that decides the A6 arm

RDR 0001's kernel is **stateless and single-step** — it resolves one transition
and never traverses a graph. Multi-step reachability exists nowhere in the
cluster except this RDR's lint. That is why no peer declares a root today: no
peer needs one. The declarations are lint-only inputs.

RDR 0002's schema is closed by two normative clauses a route-(a) fix must amend
together: the layout enumeration (root `outcomes`, `[model]`, `[tags.<tag>]`,
`[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[dump]`) and "`[model]` MUST
contain `id` and `version`". 0002 is Draft with all 8 assumptions Verified and
one open Prerequisite that gates *implementation sequencing*, not lock.

### Minor cohesion gaps found (not A-records; for the pre-lock lenses)

1. **Presence dimension** — `0003::Normative Contracts` states "Lint MUST NOT
   drop `exists` atoms from the product: a group carrying one stays provable,
   and certifying it exhaustive while ignoring the presence dimension is the
   false-green this RDR's narrowing forbids." The words "presence" and "exists
   atom" occur nowhere in 0006's body; its invariant 4 speaks only of
   "participating guard dimensions".
2. **Vacuous `exists`** — `0003::Technical Design` operator table says an
   `exists` atom over an always-present key is "well-formed but vacuous — lint
   reports it as such rather than rejecting it". 0006 has no advisory finding
   for it.
3. **Always-present conformance half** — `0003::Normative Contracts` defines a
   conforming view as (a) every always-present key present AND (b) every
   single-valued tag holding at most one value. 0006's invariant 5 covers (b);
   "always-present" occurs zero times in 0006. Note 0003 defaults a declaration
   with no optionality marker to *optional*, so always-present is opt-in and
   rare, and 0003 routes this half to RDR 0007's view assembly as well.
