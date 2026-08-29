Model: claude-opus-5[1m] (dispatcher; persona stamps in their own files)

# 3amigo consolidation — cli/0023

Three isolated persona passes (no cross-persona visibility). Consolidation is
mechanical: hotspots are the SET INTERSECTION on anchored element ids, not a
re-judgement by a reader who has seen all three. Overlap marks a hotspot
passage, not a validated finding; a single-persona finding is not thereby
weaker.

Sources: `persona-1-pm.md` (PM), `persona-2-implementer.md` (IMP),
`persona-3-qa.md` (QA).

## Anchor index (id → personas)

| Element id | PM | IMP | QA | Personas |
|---|---|---|---|---|
| `0023:C1` | — | F1, F5, F6, F8, F9 | S4-half, S1-half | **2** |
| `0023:S4` | — | F1 (via C1) | HIGH | **2** |
| `0023:C2` | — | F3, F7 | MEDIUM (S3) | **2** |
| `0023:S1` | — | F6 (width unit) | HIGH, MEDIUM | **2** |
| `0023:A5` | — | F1, F2 | widened-cite | **2** |
| `0023:MVV` | LOW (step 6) | — | LOW (step 5) | **2** |
| `0023:§performance-expectations` | HIGH (token/byte) | F6 (width unit) | LOW (wall time) | **3** |
| `0023:§problem-statement` | HIGH ×2, MED | — | — | 1 |
| `0023:§approach` | MED | — | — | 1 |
| `0023:§implementation-plan` Phase 2/3 | MED | F4 | MEDIUM (6-vs-5) | **3** |
| `0023:§prerequisites` | — | F2 | — | 1 |
| `0023:§consequences` | LOW | — | — | 1 |
| `0023:S3` | — | F3 | MEDIUM | **2** |
| `0023:S5` | — | F5 | MEDIUM | **2** |
| `0023:S2` | — | — | LOW | 1 |
| `0023:D-naming` | — | F10 | — | 1 |
| `0023:D-wire-byte-format` | — | F7 | — | 1 |

## Hotspots (≥2 personas, independently)

**H1 — `0023:C1` / `0023:S4` / `0023:A5`: the WHOLE-tree walk names a command
the tree does not contain.** IMP F1 and QA S4 independently probed the real
root and found `NewRootCmd().Commands()` = `docs, flow, help, lint, version`
— no `completion`. C1 mandates the walk descend "INCLUDING the auto-generated
help and completion commands"; cobra creates `completion` lazily
(`InitDefaultCompletionCmd`, called from `ExecuteC`), so a walk built against
`NewRootCmd()` asserts a negative over a tree missing the very command class
the clause was widened to cover. QA adds that S4's own stated negative control
("a control asserting the walk actually reaches `help`") is blind to this,
because `help` IS present by default — the control cannot distinguish a total
walker from one that merely never had `completion` to skip. Blocks Phase 2's
one piece of new machinery with no shipped exemplar; A5's Pending
implementability half is measured against this same clause.

**H2 — `0023:C2` / `0023:S3`: reflective partition-completeness names no
artifact to reflect against.** IMP F3 (no in-code home for the assignment —
struct tag, package map, or test map) and QA (the reflective oracle is
tautological if ECHO is derived from the projection code itself, and S3's
negative control is shared with S2 so it does not discriminate). C2 states the
assignments in prose only. Blocks Phase 1, and is entangled with A2's
deliberately deferred mechanism choice.

**H3 — `0023:S1` / `0023:C1`: the width clause's unit and mode scope.** IMP F6
(the "STRICTLY SHORTER" unit is ambiguous — S1's fixture cites the 290→152 B
NDJSON *line*, `§performance-expectations` cites 708→146 B *payload*; mode
scope unstated) and QA MEDIUM (the width clause is enforced in JSON only; S5
asserts a text *subset*, and a subset is satisfied by the equal set, so a
zero-saving text projection passes). C1 itself calls this clause "the flag's
whole reason for existing" and it is the sole guard against premortem P-12.

**H4 — `0023:§implementation-plan` Phase 2 oracle count and reach.** All three:
PM (no adoption step / no named adopter — "done" is defined purely as oracles
green), IMP F4 (Phase 3 edits `flowResolveExtendedDesc`, which feeds the
GENERATED `docs/cli-reference.md` gated by `make docs-check`, Makefile:34,105 —
Phase 3 names only the hand-written `cli-output-contract.md`, so the Phase 3
commit fails CI), QA MEDIUM (Phase 2's done-criterion counts SIX oracles,
Testing Strategy enumerates FIVE — absent-not-null has no S element, no
discriminability row, and no negative control, and it is the one property A2
found a live hazard for).

**H5 — `0023:§performance-expectations`: measurement denomination.** PM HIGH
(the outcome is denominated in transcript TOKENS; every measurement in the
record is in BYTES, with no conversion and no stated proxy justification —
A1's "If wrong" therefore sets a token bar that byte evidence cannot fail),
IMP F6 (same section, unit ambiguity), QA LOW ("wall time unchanged" carries
no threshold, oracle, or scenario).

**H6 — `0023:A5` / `0023:§prerequisites`: a contradiction on the starting
gun.** IMP F2: Prerequisites says "All Critical Assumptions verified"; the
projector reports A5 **Pending**, its implementability half deferred to a
Phase-2 spike whose `If wrong` reopens C1. Read with H1 this is not a
bookkeeping nit — the deferred spike is the one that H1 says is mis-specified.

**H7 — `0023:MVV`: oracle thresholds.** PM LOW (step 6 says "record" byte
counts with no threshold; A1's failure condition is qualitative — blocks the
`G-scope` gate response) and QA LOW (step 5 runs one sibling verb where C1's
non-registration spans three; consistent with C1's "corroboration, never the
assertion", recorded only).

## Single-persona findings (carried, not weaker)

- **PM HIGH `0023:§problem-statement` vs `0023:A1`** — headline motivating
  payload is 1529 B; A1's verified 48-fact measurement is 708 B. Only the
  *percentage* discrepancy is reconciled downstream (§consequences retires the
  72%); the absolute width is never retired. Dispatcher note: A1's 708 B is a
  **synthetic** model authored for the spike
  (`evidence/spikes/a1-byte-width.md`, S1b "pricing-48.toml, scratch"), while
  1529 B is the seed's real motivating call — different calls, so not a
  contradiction, but the record never says so.
- **PM MED `0023:§problem-statement` (routable question) vs
  `§decision-rationale`** — the Problem Statement promises this RDR decides
  "whether future verbs inherit the projection"; the body defers it to JDR
  0002 §D1's later-verb clause, which "obliges [0021] to nothing."
- **PM MED `0023:§approach` / Phases 1–3** — no adoption step, no named
  adopter; an opt-in flag landing with zero adopters yields zero of the
  motivating benefit with every oracle green.
- **PM LOW `0023:§consequences`** — the success-only carve-out is disclosed
  but never sized against the motivating consumer's traffic mix.
- **QA HIGH `0023:S1`** — "identical invoked-reader set" is asserted over a run
  where `readers` is projected away (it is an ECHO field). The only shipped
  channel is the payload field, and the precedent helper
  (`internal/cli/flow_fixtures_0011_test.go::readersOf`, :1769) `Fatal`s on
  absence; no substitute channel is named. This is the assertion that would
  catch C1's own "skips work whose only consumer is a projected-away field"
  breach.
- **QA LOW `0023:S1`/`S2`** — the cited "normative fixture" file contains
  neither 290/152 nor the word "normative"; the figures live only in A1's
  table.
- **IMP S3 `0023:C1` F5** — "carried keys keep relative order" is wire-only
  (text `flatten` sorts alphabetically). `§desk-trace` disambiguates; C1 does
  not.
- **IMP S3 `0023:C2` F7** — the ALWAYS-KEEP core is declared and leaned on by
  Risks, but no S1–S5 oracle enforces it; S2's literal key set is a
  today-assertion only.
- **IMP S3 `0023:C1` F8** — confirm refusal-path flag-blindness is STRUCTURAL
  (projection site after the last `respond.Fail`) rather than behavioural, so
  no defensive branch is added.
- **IMP LOW `0023:C1` F9** — `decision_table_0010_test.go:435` pins
  `NumField() == 14` and :422-426 pins 13 wire keys; A3's caveat is correct
  for both A2 arms — recorded as a hard Phase 1 constraint.
- **IMP LOW `0023:D-naming` F10** — "long form only, boolean, default false"
  is stated twice with no oracle; the `--all` precedent has exactly that test.

## Isolation check

No persona file references another persona's output. The H1 convergence (IMP
F1 / QA S4) is two independent probes of the real command tree reaching the
same result — agreement, not conformity. Dispatcher re-probed independently
and reproduced it (`NewRootCmd().Commands()` → docs, flow, help, lint,
version; `completion` appears only after an explicit
`InitDefaultCompletionCmd()`).
