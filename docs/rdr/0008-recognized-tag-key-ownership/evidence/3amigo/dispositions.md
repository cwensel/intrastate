Model: claude-opus-5[1m]

# 3amigo lens — resolve pass (iteration 1)

Origin ledger = `consolidation.md` (PM-1…5, IMP-1…5, QA-1…5; 15 entries).
Every entry exits exactly once. Grounding gate applied per finding against
(1) code on `main`, (2) `{RDR_RESOURCES}` design docs, (3) the RDR's own
decided text — including cove's F-1/F-5 dismissals, so re-raises are cited
rather than re-fixed.

| # | Hotspot | Disposition | Section touched |
| --- | --- | --- | --- |
| IMP-5(a) / QA-3 | H-1 | **fixed** | Normative block 2; scenario 5 |
| IMP-1 | H-2 | **fixed** | Normative block 5; A7; Approach 5–6; Phase 2; scenario 9 |
| IMP-2 / QA-1 | H-3 | **fixed** (fork collapsed) | Load-Bearing Decisions / Enforcement locus; block 4; scenario 6; A6 Residual; new A8 |
| IMP-3 / QA-4 | H-4 | **fixed** | Normative block 3; scenario 7 |
| PM-1 | — | **fixed** | Problem Statement; Alternative 2 Cons |
| PM-2 | — | **fixed** (in-scope half) + **charted** (C-2) | Approach 3; `charted.md` |
| PM-3 | — | **fixed** (scoping honesty) + **charted** (C-1) | Failure Modes; `charted.md` |
| PM-4 | — | **fixed** | Consequences (naming-freedom bullet) |
| PM-5 | — | **fixed** | Decision Rationale deciding rows |
| QA-2 | — | **fixed** | Scenario 3 |
| QA-5 | — | **fixed** | Scenario 8 |
| IMP-4 | — | **fixed** | Testing Strategy "Done"; MVV |
| IMP-5(b) | — | **fixed** | Normative block 2 (checked positions) |

No finding was dismissed. No finding traced outside the ledger.

## The two findings that changed the design

**IMP-1 — block 5's lint locus was wrong.** Grounded by reading peer
contracts rather than the RDR: RDR 0002's *normative* field layout enumerates
every authored construct (root `outcomes`, `[model]`, `[tags.<tag>]`,
`[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, `[dump]`) and carries no
`requires_owned` key; the token has zero occurrences repo-wide. `RequiresOwned`
is a *derived* normalizer output, so a source-lint rule over authored TOML has
nothing to inspect. Worse, RDR 0007 narrows the field to post-guard
write-dependency keys (`0007…md:2075`), i.e. derived from `[rule.write]` —
where a `recognized` entry is already rejected by RDR 0002's pre-existing
`write to non-owned tag` category (`0002…md:345`). Block 5 as written was dead
code in the wrong category.

This is the clause **cove added last pass** (F-8), which is the useful part:
the kernel-side residual F-8 identified is real and stays (hand-constructed
rows still yield `owned_state_unavailable` naming the reserved key), but its
enforcement locus was assumed, not verified. Block 5 now binds the normalizer
as a *producer* and states that the current derivation path discharges it by
construction. Scenario 9 loses its unwritable lint half.

**IMP-2 / QA-1 — the Enforcement locus, collapsed rather than escalated.**
The RDR itself said "Pick the final form at Pre-Lock," and Pre-Lock is this
pass. The tiebreaker-reduction gate applies: the evidence dissolved the fork.
Candidate (a) was refuted by in-repo evidence the RDR had already surfaced but
not applied — `internal/cli/respond/respond.go:22-25` reserves the terminal
`"ok"`/`"failed"` names with zero enforcement, a live instance of how (a) ages
— combined with A6's Residual conceding that under (a) the shadowing is
"detected by nothing." Between (b) and (c): (c) alone cannot see a producer
that never calls it; (b) alone duplicates the predicate at the boundary. RDR
0009's "one predicate, two call sites" pairing is the established peer shape
and is what stops the enforcement points drifting. Settled: **(b)+(c) paired.**

Surface conceded explicitly rather than hidden — Consequences now carries a
negative bullet saying the kernel gains one exported predicate plus its entry
call, and that "no kernel change" is scoped to disposition for conforming
input.

## Net-new scope — charted, not absorbed

Two items in `charted.md`: **C-1** the advisory intent-mismatch heuristic
(PM-3) — a second contract with an unsized false-positive profile, suggested
successor RDR after 0002's normalizer exists; **C-2** a user-facing surface for
data-level validation failures (PM-2) — belongs to RDR 0005's extensible
code-string table, carried as a cluster-reconcile watch item. Neither was
folded into this RDR.

## Needs (re)verification (Stage 6 closes these)

- **A8 (new, `Verified` this pass)** — the settled locus introduced a
  load-bearing claim no assumption covered: that an exported `Input` predicate
  plus a `Resolve`-entry call disturbs neither RDR 0001's locked contracts nor
  the existing suite. Verified immediately because both halves were cheap and
  would have changed the design if false: REQ-6 quoted verbatim
  (`0001-…/artifacts/req-list.md:34`) *affirmatively reserves* the Go error path
  for programmer mistakes, and the only `Key: "recognized"` in the test tree is
  a `Row.Match` pattern (`fixtures_test.go:237`), never an `Input` field — so
  no existing fixture trips the precondition.
- **A7 (`Pending`, unchanged status; basis materially corrected)** — the
  ownership half still routes to Stage 7.1 cluster-reconcile. Two corrections
  landed: the enforcement locus (above), and a citation precision fix — 0007's
  "single normative home" claim scopes to the *guard-decidability domain rule*
  (`0007…md:1098-1104`), with `RequiresOwned` narrowed as a consequence; the
  draft had paraphrased it as ownership of "the field's meaning."
- **No previously-`Verified` assumption was invalidated.** A1–A6 stand; A6's
  Residual is closed by the locus decision rather than falsified.
- **Carried watch item (unchanged from cove)**: RDR 0007 A13's
  negative-existential remains falsifiable by a future peer — cluster-reconcile.

## Tiebreakers escalated to the driver

None. The one genuine fork (Enforcement locus) was collapsed with evidence per
the tiebreaker-reduction gate, as the RDR's own text instructed for Pre-Lock.
