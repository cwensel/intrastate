# RDR 0008 — Stage 4 scoped re-entry: A9 / A4 re-verification

Model: claude-fable-5
Date: 2026-08-23
Scope: `Draft [revised from Final 2026-08-23; re-verify A9, A4 — …]` — the
re-verify set is {A9, A4} plus the anchors the demotion's finding set touched
(block 1's stale peer fact, block 2's predicate-position and lower-bound
clauses, scenarios 3 and 5, Phase 2, the Risks/Consequences/Failure Modes
sites that restated A4/A9, and the "sweep on the way" items the demotion
named). Every other assumption's `Verified` stamp carries forward and was NOT
re-derived. `prior-art.md` (Stage 2 cache) and `reverify-a6-a11.md` (first
re-entry) are reused, not re-run.

## The defect, confirmed

A9 and A4 were verified on 2026-08-11 against RDR 0002's normalizer spike as
it then stood (`a9-normalization.md`). 0002's Stage 4 rebuilt the spike to
its dump contract with outcome lifting (`8cf1dbc`, 2026-08-22), renamed its
fixtures per JDR 0001 §JD-10, and re-locked (`5bed32a`, 2026-08-23) with
fenced text that (a) makes the `recognized` match atom the mandatory outcome
binding lifted into the row's outcome field, (b) refuses a guard-position
`recognized` atom at load, and (c) makes load fail-fast. 0008's A9 asserted
the negation of (a), its block 2 described a shape (b) rejects, and its
scenario 5 counted two failures against (c).

Same escape class as the first re-entry's A6/A11: a verdict on a moving
`Final` peer's artifact stamped with no re-verification trigger. Ledger row
appended to `0008-…-postmortem.md`.

## Verdicts (Method + evidence pointer)

| ID | Verdict | Method | Evidence |
| --- | --- | --- | --- |
| A9 | **Verified, restated** (previous reading withdrawn) | Peer RDR + Spike | RDR 0002 → Technical Design → Normative Contracts, outcome-binding block; Load-Bearing Decisions / Outcome binding; JDR 0001 §JD-10 (ANSWERED 2026-08-23). Spike: `spikes/a9-lifted-outcome.md` §1 — re-run reproduces `0002…/evidence/spikes/output.txt` byte-for-byte; `main.go::liftOutcome`, `::Row.Outcome`. Kernel: `internal/resolve/resolve.go::Row.Outcome`, gated in `::Resolve` / `::escapeOrRefuse`; `::assemble` injects the view tag. |
| A4 | **Verified, restated** (census corrected) | Source Search + Spike | On disk: `0002…/spikes/rdr-fixture.toml` `[tags.recognized]` + four `[rule.match.recognized]`; `kata-fixture.toml` `[tags.recognized]` + two. Remaining: `0003…/spikes/guard-fixture.toml:36-37,72` — refused at load (`spikes/a9-lifted-outcome.md` §5). Second clause unchanged: `fixtures_test.go::recognizedTagSensitiveTable` match pattern only. |

Recommendation pressure: **none against the approach — strengthened.** Final
0002 adopted this RDR's naming rule and category into its own fenced text and
cites 0008 as their owner; the seam this RDR was seeded on is closed on the
declaration side. What stays 0008's: the kernel-side ratification (block 1),
the naming rule and category as their single home (blocks 2–3), the
failure-payload contract (block 3 — not absorbed by 0002; rides JDR 0001
§JD-8), and the `Input` / `RequiresOwned` predicate (blocks 4–5).

## Observation, not a finding

**0002's spike lifts from any block.** `main.go::liftOutcome` lifts a
`recognized` atom from *any* block; 0002's fenced text requires refusal for
`guard.all` / `guard.unless`. Mutants in `spikes/a9-lifted-outcome.md` §2
normalize successfully. Checked against 0002's own Testing Strategy before
routing: scenario 3 already names both guard-position mutants as owed and
records the spike as witnessing nine categories with the remainder owed at
implementation. So this is a booked witness gap on 0002's side, not an
unrecorded defect — nothing to route, no edit to 0002 (Final; RDRs are not
amended from a sibling's stage). 0008 restates the fenced rule.

## Sweep on the way (demotion's "Direction")

| Site | Was | Now |
| --- | --- | --- |
| Overrides | "invalidates its three canonical fixtures" | two fixtures, renamed at 0002's re-lock |
| Problem Statement | "all three committed … chose an author name" | historical (before 0002's re-lock); two since renamed |
| Key Discoveries | 0002 "never binds the name" | superseded on 0002's side; new bullet for outcome lifting / guard refusal / fail-fast |
| A1 (A1d + handoff sentence) | both fixtures `[tags.outcome]`; pointer missing in 0002 | `[tags.recognized]`; pointer now in 0002's fenced text |
| A7 | quoted 0007 A13 sentence; JD-4 "open" | 0007's A13 no longer carries it (re-anchored to 0007 Joint-check / A13 → §JD-9); JD-4 answered by §D4 |
| A12 | handoff form "undetermined, open at JD-10" | cross-RDR-edit form has happened; payload/advisory residual rides §JD-8 |
| Approach step 3 | JD-8 "open there" | widened at iteration 2; also carries payload + advisory channel |
| Block 1 (fenced) | "nothing in this RDR or RDR 0002 forbids that alphabet entry" | 0002 forbids it at load; residual is RDR 0001's for other producers |
| Block 2 (fenced) | predicate positions "resolve to" the declaration; "no lower bound … lints clean and refuses `no_match`" | 0002's outcome binding decides predicate positions; lower bound is 0002's; residual scoped to non-0002 producers |
| Block 2 fall-through | conforming model's row "never fires" | narrowed: the rule must also bind `recognized` or fail load |
| JD-10 boundary | "(open)" | answered 2026-08-23; blocks 1–2 do not restate |
| Enforcement locus | "0009 … still `Draft`" | `Final`; symbol names 0009's (A11) |
| Decision Rationale (row + joint-check) | rename of three fixtures; stale 0007 quote | rename landed at 0002; 0007 re-anchored; 0002 inbound adoption noted |
| Consequences / Risks | rename inventory; "mechanical, no gating semantics" | discharged on 0002's side; rename *is* the gating semantics |
| Risks (JD-4) | "open there" | answered by §D4 |
| Failure Modes | "Silent, residual — not a narrow path" | narrowed by 0002's outcome binding; stray-predicate case and non-0002 producers remain |
| Prerequisites | handoff gate unchecked, form owned by JD-10 | checked (cross-RDR edit landed); new residual gate on §JD-8 carrier |
| Phase 2 | rename three declarations + six sites | no rename remains; 0003's fixture is 0003's, refused at load |
| Testing Strategy header | iteration-1 cluster state only | iteration-2 closures named |
| Scenario 3 | — | note: 0002 mints no guard atom on `recognized`; kernel-plumbing test |
| Scenario 5 | "two failures, not one" | exactly one refusal, by category; normative fixture named |

Not swept (out of scope, unchanged): A2, A3, A5, A6, A8, A10, A11, A13; blocks
3–6; scenarios 1, 2, 4, 6–9; Alternatives; the `## Refinement Context`
section (Stage 8 deletes it on re-lock).

## Author's round (Stage 4's single interaction)

Fixture rendered: **one** — F-1, scenario 5's single-refusal witness
(`spikes/a9-lifted-outcome.md` §3: one line, `reserved_tag_key`, exit 1).
Read from the spike, not invented; recorded as scenario 5's normative fixture
pending approval.

Questions put, each with grounding:

1. **Phase 2 scope** — drop the fixture rename entirely (0002's done; 0003's
   `guard-fixture.toml` is not a 0002 model and is refused at load) rather
   than routing a rename to RDR 0003. Grounding: §5 of the spike file; 0002
   Testing Strategy scenarios 1–2; JDR 0001 §JD-16 lists 0002×0003 as the
   sibling pair for layout reconciliation.
2. **0002 spike lift-from-any-block** — record as an observation (chosen)
   vs. file a kata against 0002. Grounding: 0002's Testing Strategy scenario 3
   already owes both guard-position mutants at implementation and records the
   spike as a nine-category witness; a kata would duplicate 0002's own
   obligation.

Dispositions are recorded in the close packet once answered.

## Profile

Unchanged: **`foundational`**, one contract (who owns the *name* of the
recognized-provenance tag key in the assembled view). Recount from the
Normative Contracts section: six ```normative``` blocks, all clauses of the
single reserved-key identity rule and its enforcement; no second independent
contract (no distinct hash, wire format, taxonomy, or destructive-op policy),
so the ≥2 split signal does not fire. The cross-RDR-producer trigger is
stronger than at either lock: 0002 has now adopted the rule into its own
fenced text, and 0008 co-resides with 0009 at `Resolve` entry (§JD-5).
Accretion floor: not applicable — `Seam Lineage` records no prior accretion.
