Model: claude-opus-5[1m] (dispatcher; run-1 and diff contexts stamped in their own files)

# Repeatability resolve — RDR cli/0011 (lite, iter-1)

Variant: lite (profile `large`, Determinacy add-on). Origin ledger = `diff.md` F1–F6.

| # | Disposition | Ledger entry | Section touched |
| --- | --- | --- | --- |
| F1 | **fixed** (pin) | gate-id source set unowned | `C1` — names `Row.Gate` as the source; excludes model-level and 0005 gate-accessor readings |
| F2 | **fixed** (pin) | `not-evaluated` emission per disposition | `C1` — emitted on every reported disposition incl. `owned_state_unavailable`; fidelity limit scoped to the guard payload |
| F3 | **fixed** (pin) | guard-payload read's disposition guard | `C1` — read MUST be guarded on `guard_unevaluable`, not read-then-filter |
| F4 | **dismissed-with-cite** | `excluded` arity / view ownership | — grounding refutes the premise |
| F5 | **fixed** (pin + single-source) | `--all` filter vs dedup order | `C1`, `C2`, `S2`, `S7`, `authority` table, audit row, Illustrative Code |
| F6 | **dismissed-with-cite** | demand term's outcome scope | — grounding refutes the premise |

## Grounding notes (code on `main`)

- **F4 dismissed.** The diff's premise — and A11's prose — say the shipped shape rebuilds the view per row. It does not. `assembledView` is called exactly ONCE at `flow_next.go:102`, outside the row loop, and `summarize(row, view, gatesRan)` (`flow_next.go:168`) RECEIVES the map. `excluded(row, owned, observed) bool` (`flow_next.go:276`) takes raw slices only because it builds a probe and performs no presence test at all today (`probe.Match = nil`). So C1's "fixed ONCE per invocation, before the row loop" already describes shipped behaviour; there is no fork to pin. Recorded positively as A16 rather than left implicit, since the separability it establishes is what licenses F5's fix.
- **F6 dismissed.** `invokedReaders(m *table.Model, outcome string)` (`flow_exec.go:86`) ALREADY filters by outcome — `if outcome != "" && row.Outcome != outcome { continue }` (`flow_exec.go:88`) — so the new match-owned term inherits that filter by construction; it is not a free choice. Escape rows are likewise already included when their outcome matches (the existing comment at `flow_exec.go:89-101` says a matching-outcome escape row "keeps its demands"). Run-1's unfiltered rendering for both callers is a run error against shipped code, not an RDR silence. C1's "a union over the outcome's rows" is correct as written.
- **F1 confirmed but narrower than framed.** Shipped code reads `row.Gate` (`flow_next.go:203`; field `internal/table/model.go:174`) — row-scoped. Run-1 rendered this correctly; the RDR simply never named the field, leaving the model-level and gate-accessor readings open to an implementer. Pinned.

## Amendment sweep (§amendment-sweep)

Amended tokens: dedup ordering, gate-id source, payload-read guard. Sites swept and updated:
`C1` (dedup-last rule), `C2` (emission-scope rule), `authority` mini-check row, Existing Infrastructure Audit `Candidate reporting` row, Illustrative Code intent line.

**Semantic disagreement caught, not just staleness:** S7 argued dedup-on-pair-vs-key "is NOT discriminated by any reachable input." True as written for the dedup KEY question, but the F5 fix makes the filter/dedup ORDER discriminable — on a matched-AND-written key under `--all`, filter-then-dedup keeps the owned-key fact while dedup-then-filter deletes it. S7's claim corrected to scope it to the dedup-key question, and the new discriminating oracle added to **S2** (where `--all` is tested) rather than S7, with a note on why it and C3's not-written-key oracle are complementary: this one fails on the wrong ORDER, that one fails on a filter that never RUNS.

## Needs (re)verification (Stage 6)

- **A16 — added, Status: Verified** (Source Search, this session). Establishes the three emission sites are separable and `Row.Gate` is the gate-id source. No prior assumption covered it: A2 covers the atom walk only.
- **A11 — no flip, but its Evidence prose is loose.** It says `runFlowNext` "passes that same slice to `assembledView(owned, req.observed)` and to `excluded(row, owned, req.observed)`", which reads as two view builds. Both call sites are real, but `excluded` builds no view — the sentence is about slice provenance, not view construction. A16 now states the shape unambiguously; A11's verdict (key-set agreement) is unaffected and stays Verified.

No assumption flipped to Pending; no spike owed.

## Escalation

**Not escalated.** Criterion (a) does not fire — 7 GUESS markers plus two internal self-contradictions in run-1 show real under-determination, not overfit. Criterion (b): F6 was the only cross-RDR-anchor candidate and it was refuted by grounding, so nothing load-bearing survives on the shared `invokedReaders` anchor; `0010:A2` independently empties the term over 0010's `decision-table` class under every reading. F1–F5 are `flow next` payload-internal.
