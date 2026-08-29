Model: claude-opus-5[1m]

# 3amigo resolve — cli/0023 (iteration 1)

Origin ledger = `consolidation.md` (H1–H7 hotspots + the single-persona
findings). Every entry below traces to one. No net-new scope was folded in.

## Dispositions

| # | Origin | Disposition | Section touched |
|---|---|---|---|
| H1 | IMP F1 + QA S4 (`C1`/`S4`/`A5`) | **fixed** | C1 WHOLE clause; S4; S4 discriminability row; desk trace step 5; A5 verification plan |
| H2 | IMP F3 + QA MED (`C2`/`S3`) | **fixed** | C2 enforcement clause; S3; S3 discriminability row |
| H3 | IMP F6 + QA MED (`C1`/`S1`/`S5`) | **fixed** | C1 width clause (unit + both modes); S1; S5; S1/S5 discriminability rows |
| H4 | PM MED + IMP F4 + QA MED (Phase 2/3) | **fixed** | Phase 2 (six→five, absent-not-null homed on S2); Phase 3 (generated docs); S2 + its discriminability row |
| H5 | PM HIGH + IMP F6 + QA LOW (denomination) | **fixed** | Problem Statement (two measurement notes); Performance Expectations (wall-time non-budget) |
| H6 | IMP F2 (`§prerequisites` vs `A5`) | **fixed** | Prerequisites |
| H7 | PM LOW + QA LOW (`MVV` 5/6) | **fixed** | MVV steps 5 and 6 (pass bar; corroboration rationale) |
| PM-a | `§problem-statement` vs `A1` (1529 vs 708) | **fixed** | Problem Statement note (2) — different calls, seed's headline retired by A1 |
| PM-b | routable-question inheritance promise | **fixed** | Problem Statement — bounded to mechanism + registration; §D1 owns permission |
| PM-c | no adoption step / named adopter | **fixed** (boundary stated, scope NOT absorbed) | Consequences — "scope of the deliverable" |
| PM-d | refusal carve-out unsized | **fixed** | Consequences — unsized, and why it does not change the landing call |
| QA-a | `S1` reader set unwritable | **fixed** | C1 (projection-independent channel); S1; S1 discriminability row |
| QA-b | "normative fixture" citation | **fixed** | S1/S2 — record lives in A2 spike §Reference output, widths in A1 row S1 |
| IMP F5 | order semantics wire vs text | **fixed** | S5 — set membership, not subsequence |
| IMP F7 | always-keep unenforced | **fixed** | C2 — enforcement-by-coincidence-of-scope stated; successor owes a mode-quantified oracle |
| IMP F8 | refusal flag-blindness structural | **fixed** | Phase 1 — projection site after last `respond.Fail`; no defensive guard |
| IMP F9 | 14-field pin | **fixed** | Phase 1 — field count fixed; a field addition would falsify A3 |
| IMP F10 | flag shape unoracled | **fixed** | S4 — positive registration-shape assertion on the `--all` precedent |

No finding was dismissed and none was charted — all seventeen grounded.

## Grounding performed (executed, not argued)

- **H1** reproduced independently by the dispatcher: a probe test on the real
  `NewRootCmd()` printed children `docs, flow, help, lint, version`; after an
  explicit `InitDefaultCompletionCmd()` the set gains `completion`. `help` is
  present only because `help_all.go:135` force-inits it. The probe file was
  removed after the run.
- **QA-a** confirmed at `internal/cli/flow_fixtures_0011_test.go::readersOf`
  (:1769, `t.Fatalf` on absence) with 9 call sites in
  `flow_demand_0011_test.go`; the derivation is
  `internal/cli/flow_exec.go::invokedReaders` (:86).
- **IMP F4** confirmed in `Makefile`: `DOCS_FILES = docs/cli-reference.md
  llms.txt` (:103), `docs-check` (:105), and `check: … docs-check test` (:34).
- **IMP F9** confirmed at `decision_table_0010_test.go:435` (`NumField() != 14`)
  and :422-426 (13 wire keys).
- **QA-b** confirmed: `a2-encoder-mechanism.md` contains the projected record
  (§Reference output) but neither "290"/"152" nor the word "normative"; the
  widths are row S1 of `a1-byte-width.md`. Citation split, not absent.
- **H3 unit check, computed**: piping both reference lines through `wc -c`
  returns exactly 290 and 152 — A1's figures already measure the FULL emitted
  line including the `{"type":"ok","data":{…}}` envelope, so pinning that unit
  in C1 matches the shipped evidence rather than re-defining it.
- **A2 re-read** for the text-mode width claim: the spike measured text
  DETERMINISM only (10 runs, one sha256). It records no text byte widths —
  which is why A8 was booked rather than closed.

## Needs (re)verification — carried to Stage 6

- **A8 (NEW, Pending)** — the strict-width clause is satisfiable in TEXT mode
  on the same unit as JSON. C1 now asserts strict width in BOTH modes; A2
  measured text determinism only, and the desk trace's "15 lines → 9" is a
  line count, not bytes. Method: Spike at Phase 2, alongside the S5 oracle.
  Booked because the clause is a MUST and no artifact measures it.
- **A5 (already Pending, scope corrected)** — its verification plan now
  requires the spike to materialize both auto-generated commands and assert
  their presence in the walked set. The pre-existing plan would have verified
  a walker against a tree lacking `completion`, i.e. the H1 defect.

No previously-Verified assumption was invalidated. A1, A2, A3, A4, A6, A7 are
untouched by these edits: no fix changed a measured figure, the encoder
mechanism, the predecessor-oracle claim, the 0024 join, the refusal-envelope
read, or the echo-group characterization.

## Tiebreakers

None escalated. Two forks were collapsed with evidence rather than handed up:

- **Always-keep oracle (IMP F7)** — the fork was "mint a sixth oracle" vs
  "leave it". Collapsed by observing the two options have the same extension
  today: with one projection mode, S2's literal already pins all four core
  fields, so a new oracle would assert exactly what S2 asserts. Recorded the
  real obligation (a mode-quantified oracle) on the successor that creates the
  second mode, where it first has content.
- **Adoption (PM-c)** — the fork was "add an adoption phase" vs "leave the gap
  unstated". Both rejected: an adoption phase is net-new scope (a consumer
  call-site migration has its own blast radius), and silence was the actual
  finding. Stated the boundary explicitly instead, naming the completion bar
  and the natural first adopters without scheduling them.

## Convergence

Iteration 1 fixed all seventeen ledger entries; none remains open. The edits
are additive precision on existing clauses — no clause was rewritten in a way
that opens a new gap — so no delta-scoped re-run is owed. `rdr lint` returns
PASS with no `resolution` findings; the surviving `placeholder:survived` hits
are all in the Finalization Gate block (Stage 7's section, not yet authored),
and the `prose:exactness` hits are `byte-identical` claims carried by the
normative fixture and the MVV.
