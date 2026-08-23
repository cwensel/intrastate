Model: claude-opus-5[1m]

# Reconcile Report — RDR 0002 (Stage 6)

Supersedes the earlier report in this file, which was written before the cove,
critique iter-2, and repeatability rounds and still claimed "A1-A7 are terminal
… no `Status: Pending` remains." Those rounds mutated the draft and booked
fifteen open claims; this pass forces each to a terminal disposition.

## Stage 5 preflight

`Profile: foundational` → lens row `cove → 3amigo → critique → repeatability`.
All four have completed evidence:

| Lens | Evidence | Verdict |
| --- | --- | --- |
| cove | `evidence/cove/{findings,dispositions}.md` | complete (subsumes grounding Step 0) |
| 3amigo | `evidence/3amigo/` + `iter-2/` (3 personas each) | complete |
| critique | `evidence/critique/` + `iter-2/` (`critique.md`, `critique-modelB.md`, `diff.md`) | complete, dual-model |
| repeatability | `run-1` `claude-opus-5[1m]`, `run-2` `Claude Sonnet 5`, `run-3` `claude-fable-5`, all `variant: full (profile: foundational)`, + `diff.md` | complete, three distinct stamps |

No variant mismatch: `run-1.md` records `variant: full`, which is what
`foundational` owes. Determinacy is not a `mid`/`large` trigger case here.
Stage 5 is complete — no return.

## Absorption audit (delegated)

A sub-agent audited all eight lens rounds for findings claimed fixed/pinned but
absent from the current RDR. Result: **no residue**. Every fixed/pinned row —
grounding G-001, cove CV-001..019, 3amigo HOT-001..004 / HOT2-001, critique
C-001..003 and M-1..M-21, repeatability R-1..R-9 and C-2/C-3/C-6 — resolves to
present RDR text in its named section.

## Dispositions

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| Pre-lock lists: grounding (i1, i2), 3amigo (i1, i2), critique i1 | 1 | VERIFIED | Each `dispositions.md` reports `Needs verification: None`; the absorption audit confirms their fixes are present in the RDR. |
| A1 TOML represents sparse rules | 2 | VERIFIED (unchanged) | Spike rerun today; see spike row below. |
| A2 Row order not part of edge selection | 2 | VERIFIED (unchanged) | `internal/resolve/resolve.go::Resolve` / `::gate` / `::missingOwned`; `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`. |
| A3–A8 | 2 | VERIFIED (unchanged) | Spike / source search / peer RDR as recorded; no round disturbed them after their last verification. |
| A9 unified block-retaining atom set conforms at handoff | 1 (repeatability) | DOWNGRADED | `Pending` / MVV Test. Structurally unverifiable before RDR 0007's kernel reshape (Prerequisites) — no `Atom` type, no block field, `Row.Guard` is a `string` on `main`. Plan: Testing Strategy sc.2 normalization assertions + sc.4 `Resolve` run against one normalized value. "If wrong" stated and survivable. |
| A10 `duplicate model id` decidable only across documents | 1 (repeatability) | VERIFIED (in part) + DOWNGRADED | Decidability half **verified at source**: `evidence/spikes/main.go::Model` declares `Model` as a singular struct field, not a slice, mirroring `[model]` as a singular TOML table (`rdr-fixture.toml` `id = "rdr"`), so one document structurally cannot carry two model ids. Behavioral half owed: sc.3 paired-document fixture. RDR record updated from "Owed". |
| A11 escape rule expands under `in`, one rescuing row per member | 4 (critique i2 item 7) | DOWNGRADED | **New assumption written into the RDR.** Kernel half verified (`::escapeOrRefuse` filters on `row.Outcome`); expansion unwitnessed — `rdr-fixture.toml` `draft-no-match-escape` binds `eq = "round-clean"`. Plan: sc.2 assertion now names an escape rule binding `in` over two members → two rows, distinct suffixes, empty write sets. |
| A12 handoff routes every atom to exactly one of `Match`/`Guard` | 4 (critique i2 item 6) | DOWNGRADED | **New assumption.** A9 asserted conformance, not totality/disjointness. Plan: sc.2 assertion that `Match` ∪ guard atoms = normalized set and intersection empty, with the kata fixture's `status.eq=closed@unless` as the discriminating case. Blocked on RDR 0007 like A9. |
| A13 merge idempotent on identical atoms | 4 (critique i2 item 1, mirror half) | DOWNGRADED | **New assumption.** sc.2 asserted only the *differing*-literal half, so a normalizer emitting duplicates passed every control. Plan: sc.2 assertion that a byte-identical `(block, key, operator, literal)` atom from rule + context collapses to count one. |
| A14 version gate precedes strict decode | 4 (critique i2 item 5) | DOWNGRADED | **New assumption.** sc.3's existing `unsupported version` fixture is v1-shaped with a bad value and trips the category under either ordering — it witnesses nothing about precedence. Plan: a **v2-shaped** fixture (`version = 2` + a v2-only key) asserted to refuse `unsupported version`, not `unknown schema field`. |
| Critique i2 items 2, 3, 4, 8 (set literals as member sequences; outcome binding reads match blocks only; `#` reserved; extended fixture sibling rows) | 1 | VERIFIED absorbed | Already carried as explicit Testing Strategy obligations: sc.2 (`["needs work"]` vs `["needs","work"]` distinct), sc.3 (`recognized` under `guard.all`/`guard.unless`; rule id and alphabet member containing `#`), and the MVV's extend-then-promote clause (two sibling candidate rows binding the same outcome). No new record owed. |
| Named spike: parse-normalize-dump (Resolve spike) | 3 | VERIFIED | Rerun today: `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes && GOCACHE=… GOFLAGS=-mod=mod GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml`. Output SHA-256 `c4be7447a241fb724632c53a9e5f39c7a9b5c7cc78e0a1e74ae4132271432a1b` — byte-identical to the recorded digest and to `evidence/spikes/output.txt`. Discharges "Exact field layout needs spikes before lock" (Alternative 1 Cons) and the Decision Rationale premortem. |
| `duplicate model id` framed as a single-rule check | 4 | VERIFIED + FIXED | Real inconsistency the repeatability amendment sweep left behind: the category-list framing said "These are the single-rule checks the arity split assigns to this RDR" while the normative clause says the category is decidable only across documents and a single-document loader MUST NOT report it. Both framing sites corrected (category block; `disposition` mini-check). |
| CV-020 — RDR 0008 asserts "nothing in this RDR or RDR 0002 forbids that alphabet entry" | 1 (cove) | ACCEPTED (routed) | Confirmed stale at `0008-recognized-tag-key-ownership.md:1234`. RDR 0002 now forbids the empty-string alphabet member and A8 leans on that clause as its sole barrier. Decisions are compatible (0008 scoped itself out); this is stale peer *fact*, and we never amend RDRs. Recorded as a Risk in this RDR and routed to `/rdr-cluster-reconcile`. |
| CV-021 — 0002 `Draft` while its consumers are `Final` | 1 (cove) | ACCEPTED (routed) | Confirmed: 0001 `Implemented`; 0003–0009 `Final`; RDR 0009 twice cites "RDR 0002's **Final** escape-rule prohibition". Producer locks after consumers. Recorded as a Risk with `/rdr-cluster-reconcile` named as the mechanism, so the obligation survives lock. |
| Exactness-word delta sweep (delegated) | 4 | VERIFIED | Sub-agent swept every exactness word against the A-records and MVV, restricted to round-touched claims. It returned exactly four uncovered round-introduced claims — now A11–A14. Its separate pre-existing list is Stage 7 mechanical-sweep scope (CHECK 4), not Stage 6's. |
| Completeness grep | — | VERIFIED | No `_Draft placeholder._`, no seed-skeleton header, no surviving template bracket. `## References` is populated with real citations. Finalization Gate sub-sections are template *spec* bodies replaced at lock by the gate.md pointer, not hollows. |

## Hard rules

**Refutation.** No spike or source search refuted an assumption the RDR relies
on. The one claim that could have — A12's routing against RDR 0007's `Block`
type carrying "exactly two constants" (`BlockAll`, `BlockUnless`) while this RDR
retains three authored blocks (`match`/`all`/`unless`) — resolves cleanly:
0007's `Block` is scoped to the **guard** atom, and this RDR's split routes
match-block atoms to `Row.Match`, so a `match` block value never reaches 0007's
two-constant type. No route-back; no §punt-ledger row owed.

**MVV deferral floor.** No downgraded item pins what the MVV proves *at lock*.
A9 and A12 are blocked on RDR 0007's unimplemented reshape, which Stage 6 cannot
unblock and which the Prerequisites already gate implementation sequencing on;
A10's load-bearing half is source-verified; A11, A13, A14 are fixture extensions
executed during the MVV itself, each with a named scenario and a failing control.
The MVV's one genuinely deferred assertion — ambiguous overlap — was already
dispositioned in-draft as RDR 0006's by the arity split and explicitly not
counted as satisfied here.

## Verdict

**RECONCILED.** All twenty items reach a terminal disposition, no BLOCKER, and
every disposition is written into the RDR itself (A10 evidence upgraded; A11–A14
added; Testing Strategy sc.2/sc.3 carry the four new oracles; Prerequisites names
why each Pending is a deliberate downgrade; Risks carries the cross-RDR routing).
Ready for Finalize.
