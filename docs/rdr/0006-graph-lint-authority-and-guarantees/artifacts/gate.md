Model: claude-opus-5[1m]

# Finalization Gate — RDR 0006 Graph Lint Authority And Guarantees

- **RDR**: `cli/0006` — `0006-graph-lint-authority-and-guarantees`
- **Date**: 2026-08-23
- **Profile**: large
- **Mechanical sweep**: PASS (`evidence/tooling-pass/tooling-pass.md`)
- **Verdict**: **READY — Gate PASS**

## 1. Contradiction Check

No contradiction survives, and the two that were found were resolved against
the authority rather than papered over.

**Escape rows (JDR 0001 §JD-14).** The draft's invariants 3 and 4 once granted
escape rows an overlap exemption and let them satisfy coverage outside the
union — the opposite of locked RDR 0003, and split against this RDR's own
Load-Bearing Decision that escape rows are modeled graph edges, not
tie-breakers. Repaired at refine in favour of RDR 0003's two-population reading.
Verified in the current text: escape rows "MUST participate in the coverage
union and MUST be overlap-checked in one population per declared failure class"
(`Normative Contracts`), and the disposition table records the ordinary-row
pairing as silent by design because it is never a runtime ambiguity — which
agrees with `internal/resolve/resolve.go::Resolve`, where escape rows are
consulted only to rescue `no_match` / `ambiguous_match`.

**Group identity (cluster gate `0003-0006-0007` F4).** `Technical Design` had
asserted group identity as `(owned-state node, recognized outcome)` while this
RDR's own contracts forbid defining a second grouping predicate and RDR 0003
(`Final [locked 2026-08-22]`) fixes that a row's match pattern selects its
group. Resolved in favour of RDR 0003: membership is the authored match
pattern; reachability decides only which groups are proven. That is exactly the
one-line citation the cluster gate scoped `0003::A10`'s discharge to, "not a
redefinition".

That resolution exposed and closed a soundness defect in the same
neighbourhood: the over-approximation guarantee had been stated globally ("only
a false positive, never a false green") when it holds only for *existential*
checks — coverage and dead-end are universal, and a ∀-claim over a widened
domain gets easier to satisfy. Now scoped per invariant in `Load-Bearing
Decisions`: existential checks read merged nodes, coverage never reads a node,
dead-end splits nodes on terminal-participating keys before testing.

Research Findings and Proposed Solution agree: the Key Discoveries name the
blocking root command, the peer-supplied graph structure, and the
deterministic-check expressibility that the Approach and Normative Contracts
then specify. No planned feature contradicts a stated principle — notably, the
"no suppression or waiver mechanism" consequence is consistent with the
blocking-authority premise rather than in tension with it.

**JDR duties discharged.** §JD-4 (lint's promise narrows; 0003 records, 0006
cites and mints no code) — verified: the narrowing clause cites
`0003::Normative Contracts` verbatim and reuses `graph-unprovable-coverage`
with a widened trigger. Its consequent duty on 0006 — carry the refusing atom
on the finding record — is discharged at the finding contract, which carries the
atom (`Key`, `Operator`, `Literal`, `Block`) when attributed to one guard atom.
§JD-13 (single-valued marker rehomed to 0003) — verified on both sides: 0003
now carries it and 0006 cites it as a field of the tag declaration model.
§JD-14 — as above.

## 2. Assumption Verification

Ten records, A1-A10, sequential. Every one is internally consistent:
Status / Method / Evidence agree and no `If wrong` is empty. Method vocabulary
conforms (7 Peer RDR, 2 Source Search, 1 MVV Test). **No `Docs Only` record
exists**, so the load-bearing-Docs-Only prohibition cannot fire. **No
`Source Search` self-reference**: A4 cites the shipped CLI failure gateway by
symbol, A8 cites the repository `.toml` census as the absence proof it asserts.

**A1-A4 Verified.** A1-A3 rest on peer-RDR contracts (RDR 0002's normalized
rows with retained source identity; RDR 0003's finite-domain coverage/overlap
derivation, tag declaration model, and scoped row group; RDR 0004's accessor
boundary). A2 was re-verified clause by clause at Stage 4 against RDR 0003's
locked text and again at pre-lock, when its Evidence was corrected: its claim to
close `0003::A10` now rests on the citation paragraph, and its claim that four
further RDR 0003 records "close by citation" was withdrawn — RDR 0003 is locked,
so this draft supplies their *producer* and each closes on the peer as a tracked
route-back. A4 is source-verified: `clierr.go::CLIError`, `::ExitCodeFor`,
`respond.go::Fail`, and `root.go::ExecuteAndEmit` all resolve, and the RDR
correctly records that lint needs new stable codes, not a new envelope.

**A5-A10 Pending, all six DOWNGRADED at Stage 6**, each with a named plan and a
written survivability rationale. Stage 6 re-verified each against the working
tree and refuted none:

| ID | Why it cannot flip pre-lock | Named plan | Survivable because |
| --- | --- | --- | --- |
| A5 | The gate proves a command that does not exist, over a model that does not exist (A8) | Validation scenario 7 — the `graph-lint` job asserting the JSON `code` | The MVV is fixture-backed and does not consume this gate; what defers is *enforcement*, not proof |
| A6 | No peer declares an initial owned state or terminals; confirmed absent at source | Additive clause on RDR 0002's authoring schema (`Draft` — scheduled peer edit) | The window is loud: a rootless model takes a blocking `graph-dangling-edge`, so an unlanded schema cannot read as clean |
| A7 | The discharge lands on `0003::A18`, and RDR 0003 is locked | Route-back at RDR 0003's next touch | `0003::A18` is itself DOWNGRADED as not lock-blocking, with this RDR named as a sufficient producer; and this RDR states the discharge as *partial*, never claiming the view-level enforcement |
| A8 | No transition model is checked into this repo — `.toml` census re-verified | Authored and homed at implementation, booked in `Prerequisites` | Gates only scenarios 7 and 23; the MVV proper needs no shipped model |
| A9 | RDR 0005 defers transition-model config discovery (its own audit row) | Flips when that discovery lands | `--model <path>` is sufficient for every scenario — ergonomics, not authority |
| A10 | No peer states write-replaces for single-valued tags; confirmed absent | RDR 0002 states the rule (`Draft` — scheduled peer edit) | Bounded and named in invariant 5's own text: a re-widening of one invariant, not a redesign. The `graph-single-valued-state` fixture proves the per-row check regardless |

**None is MVV-critical**, so the reconcile's HARD RULE 2 does not fire. The
Minimum Viable Validation is fixture-backed: every blocking invariant class is
proven against authored fixtures through the production command path, without a
shipped model (A8), the CI wiring (A5), `--flow` discovery (A9), or either peer
clause (A6, A10).

Two of the six route to peers that are `Draft` (A6, A10 → RDR 0002) — scheduled
edits on an open document, not route-backs. One routes to a locked peer (A7 →
RDR 0003 A18), tracked alongside A10/A12/A17/A19/A20 for that document's next
touch. **CHECK 6 confirms no settled-fact prose depends on any Pending
record**: the passages that read present-tense are normative contracts
specifying a requirement on the pending artifact, each carrying its `(ANN)`
cross-reference — which is what an RDR is for.

## 3. Scope Verification

**In scope, not deferred.** The Minimum Viable Validation is a fixture-backed
`intrastate lint` invocation through the production command path — the same
command shape intended for CI — passing one legal model and failing one illegal
model per blocking invariant class.

The proof is named and mechanical, not gestural. The illegal matrix must assert
`graph-dangling-edge`, `graph-dead-end`, `graph-overlap` (an ordinary-row pair
*and* an escape-row pair sharing a failure class), `graph-coverage-gap`,
`graph-unprovable-coverage` (a non-finite dimension *and* a withheld claim
naming row and atom), `graph-single-valued-state`,
`graph-always-present-owned`, `graph-owned-before-write`,
`graph-terminal-escape`, and `graph-product-too-large`. It includes a
**multi-defect group** asserting every expected code from one run, so a
first-failure engine cannot pass — that is the complete-emission clause's
oracle. The legal matrix asserts the three advisory codes on success. Every
blocking run must return aggregate `CLIError.Code = graph-lint-failed` with
`GroupUserEnv` exit behavior; every run must emit the `findings` key even when
empty.

**Determinacy disposition.** `Profile: large` and the Normative Contracts name
deterministic finding order and a canonical (never hashed) sortable fingerprint,
so the Stage 7 determinacy trigger fires. No `repeatability` run is owed:
determinacy here is a property of the *specified output*, not of the
RDR-authoring process that lens re-runs, and it already has a mechanical oracle
— Validation scenario 4 asserts JSON finding order "as full-list equality
against a golden file, not merely as a sorted property", against the identity
tuple `(model id, invariant code, source rule/context id or graph element id,
normalized predicate/write fingerprint)`. `Cross-Cutting Concerns` scopes the
claim to semantic finding identity, explicitly not byte-identical output or
content-addressed hashes. **determinacy: satisfied by MVV scenario 4** — see
`evidence/tooling-pass/tooling-pass.md`.

What scenario 7 (the CI gate) and scenario 23 (the false-positive census) prove
is *enforcement and residual rate*, and both wait on A8's model. Neither is the
MVV, and the MVV does not consume them.

## 4. Cross-Cutting Concerns

- **Versioning** — finding codes and JSON payload fields must be stable and
  append-only under the existing CLI output envelope. The advisory tier is
  *closed* at four named members, so growth is a deliberate contract change.
- **Build-tool compatibility** — the authoritative check is the same root
  `intrastate lint` command CI runs after `make build`. Two costs are named
  rather than assumed: `Makefile::check` lacks the `build` edge a
  binary-invoking gate needs, and `.github/workflows/ci.yml`'s existing `lint`
  job key belongs to golangci-lint and must survive untouched — hence the new
  `graph-lint` job.
- **Incremental adoption** — hooks and resolver-local flags may call the engine,
  but only the command/CI gate defines acceptance, and every such path must
  share the one request builder and engine. Scenario 6 makes that assertable
  now: exactly one exported engine entry point and one request builder, so a
  later alias has no second constructor to diverge through.
- **Determinism / canonical form** — deterministic claims mean semantic finding
  identity and stable invariant results for the same normalized model, never
  byte-identical output or content-addressed hashes. Covered above (§3).
- **Bounded work** — the node ceiling and the product bound are both published
  implementation constants in the command's help output, asserted by the MVV,
  with `graph-product-too-large` naming the traversal rather than running
  unboundedly.
- **Peer-owned, not addressed here** — source table format and normalization
  (RDR 0002), predicate grammar and row-group semantics (RDR 0003), accessor
  execution (RDR 0004), CLI verb surface (RDR 0005), runtime refusal (RDR
  0001), guard atom shape (RDR 0007).

## 5. Proportionality

Right-sized; nothing to trim.

This RDR owns exactly one load-bearing contract: **blocking static graph-lint
authority and the mandatory invariant set over the normalized model**,
including the reachability relation the invariants quantify over — the one
graph property no peer can state for itself, since RDR 0001's kernel is
normatively single-step. It does not own the source format, predicate grammar,
row-group definition, runtime resolver, accessor execution, or output envelope,
and its contracts cite rather than restate RDR 0003's clauses throughout.

The `large` profile re-validates at lock: the document locks graph-acceptance
invariants and CI authority with `Seam Lineage: no prior accretion`, and it ran
the profile's full lens row (grounding ×2, 3amigo ×3, critique ×2 dual-model
with strong consult and a delta pass). The 1858-line length is carried by the
invariant taxonomy, the disposition table, and the 23-scenario validation
matrix — all load-bearing acceptance criteria the implementation reads directly.
CHECK 9 finds zero Evidence fields over budget, so the mass is in the contract,
not in the records.

---

**Verdict: READY.** No open blocker. The six DOWNGRADED assumptions are
terminal dispositions with named plans, none MVV-critical; the mechanical sweep
is PASS; the MVV is in scope and fixture-backed; all three JDR 0001 duties
(§JD-4, §JD-13, §JD-14) are discharged. Locking.
