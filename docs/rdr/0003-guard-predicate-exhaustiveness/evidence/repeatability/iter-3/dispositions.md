Model: claude-opus-5[1m]

# Repeatability iteration 3 — dispositions

Lens: `repeatability` (lite, profile `large`). Origin ledger = `diff.md` D-1..D-5.
Every entry exits exactly one way. Disposition vocabulary is the resolve prompt's
REPEATABILITY DIFF clause: pin / cut / single-source / leave-non-normative /
tiebreaker.

## Per-finding

| ID | Disposition | Origin | Section touched |
| --- | --- | --- | --- |
| D-1 | **fixed** (pin + compute) | diff.md D-1 | `Normative Contracts` — assignment-count table; `trace` steps 3/4/6/7 + summary; `Testing Strategy` (2 sites) |
| D-2 | **fixed** (pin) | diff.md D-2 | `Normative Contracts` — `unless`-Kleene clause |
| D-3 | **fixed** (pin) | diff.md D-3 | `Normative Contracts` — bound clause |
| D-4 | **fixed** (pin) | diff.md D-4 | `Normative Contracts` — escape-population clause; A17 |
| D-5 | **fixed** (pin) | diff.md D-5 | `Normative Contracts` — narrowing clause; escape-rescue passage |

## Detail

**D-1 — pin the table's inputs, then compute the trace.** The table and the two
default clauses were each individually correct; nothing joined them, and the
document's only worked arithmetic silently used a third reading (single-valued,
always-present) the fixture never declared. Pinned: the table reads the
declaration defaults, and reading an unmarked declaration as single-valued or
always-present is the same exponential understatement the `set` row forbids.
Then computed rather than argued — the fixture declares no marker of either
kind, so `profile` = `2^4·2` = 32, `prelock_iterations` = `2^4·2` = 32,
`cluster_eligible` = `2^2·2` = 8, product **8,192**, not the 64 the trace
asserted.

**D-1 follow-on — a second silence the pin exposed.** With the defaults pinned,
a guard dimension is a power set, and the RDR turned out never to state how a
value atom projects onto one. Three readings (exact / superset / subset) differ
by up to 6.5x in union cardinality and disagree on whether the fixture's two
rows overlap at all — the cross-implementation divergence the published-bound
clause exists to prevent, one level down. Collapsed on the RDR's own evidence
rather than escalated: the operator matrix already restricts `eq`/`in` to the
single-valued-shaped kinds and already has a co-occurring-value operator with
fixed semantics (`contains`, "assignments containing every listed element").
So `eq`/`in`/comparisons are stated as **single-value operators**: over a tag
not declared single-valued they do not project at all and take the blocking
inability-to-prove outcome. Lint MUST NOT pick a reading. Consequence for the
fixture: its Draft group is **unprovable at step 4**, not gapped at step 6 —
trace steps 4/6/7 and the summary rewritten accordingly, and both
`Testing Strategy` sites corrected.

**D-2 — pin the contribution as nothing, and suppress the derived gap.** A row
that can refuse has no decidable accepted-assignment set: crediting its
`all`-intersection unsubtracted is the false-green the narrowing forbids, and
contributing nothing while still emitting a coverage gap manufactures a witness
assignment a row would in fact accept. Resolved by reading report-every-defect's
own qualifier literally — the separate gap finding is scoped to a gap **over a
provable product**, which a withheld group does not have. Overlap findings among
the group's decidable rows are explicitly unaffected.

**D-3 — pin cardinality as undefined over an unprovable dimension.** The table
is total over the finite kinds and defines no value for `scalar` or for a finite
kind declaring no domain, so the too-large comparison is defined only over a
fully-provable product. Stated so it does not read as short-circuiting between
the two input classes: each unprovable dimension still draws its own blocking
finding; the size-bearing refusal is simply unavailable when its quantity does
not exist.

**D-4 — pin the class read and the partition, grounded on shipped code.** The
class is read from the row's declared rescue list, which is plural in RDR 0002
(`0002:328-334`) and matched membership-wise by the kernel
(`internal/resolve/resolve.go::Row.rescues`, `slices.Contains(r.Escape, kind)`).
So: one population per declared class, a multi-class row joining each, the
pairwise check run within each population independently, and one finding per
class for a pair colliding in several. A flat pairing over-reports (the kernel
filters by `rescues` before the exact-one count); a per-row partition
under-reports. A17 amended to record that the mechanics are settled and only the
cluster's confirmation of the two-population reading is still owed.

**D-5 — pin the narrowing's population as every participating row.** Collapsed
on shipped code: `escapeOrRefuse` runs escape candidates through the same
viability gate as ordinary ones and returns `refuse(in, *blocked)` on an
undecidable guard, so a guarded escape row can refuse on its own path. The
narrowing's own justification — a claim MUST NOT be stronger than the runtime —
therefore reaches escape rows. The `:1045` "ordinary row" passage was not stale
in substance (it is about a bare escape row failing to *discharge* another row's
withholding) but read as an exhaustive quantifier; tightened to say so and to
point at the narrowing clause. The two-population split is now explicitly scoped
to overlap only.

## Amendment sweep

Per rdr-common §amendment-sweep, for each amended clause: grep the pre-edit
token, update or list every surviving site, re-read producers/consumers.

- `ordinary row can refuse` — 1 site (`:1046`), amended.
- `participating row` — 8 sites re-read; all consistent with the widened
  population (they name the narrowing generically, not the ordinary population).
- `same failure class` — 6 sites re-read; A17, the disposition row, and the
  authority census agree with the per-class partition. Trace step 7 rewritten
  for the projection change, not for this.
- `32 of 64` / `= **64**` — 4 sites, all corrected (trace step 6, trace summary,
  `Testing Strategy` fixture row, `Testing Strategy` partition paragraph).
- `denotes a subset` / `the value dimension` — 2 sites; the `exists`-clause
  sentence at `:1045` presupposed a singular value dimension and was amended to
  defer to the operator/kind agreement clause, noting `exists` is unaffected.
- `single-valued` / `single_valued` — 25 sites re-read. The declaration clause's
  own "operational content" paragraph already stated the `|domain|` vs
  `2^|domain|` split and needed no change; A16 amended (its request is now
  blocking, not completing).

## Needs (re)verification — carried to Stage 6

- **A21 (new, Pending)** — requiring `single_valued` for a provable `eq`/`in`
  atom leaves the target flows authorable. Method: MVV Test. The projection
  clause is a new load-bearing claim that makes the marker a precondition rather
  than an optimization; no fixture declares it today, so it is unverified. Plan
  and the inversion fork it opens are recorded on the assumption.
- **A16 (already Pending, consequence raised)** — its authoring request is now
  blocking: an unwritable marker means no `eq`/`in` guard can carry an
  exhaustiveness claim at all.
- **A17 (already Pending, unchanged in status)** — mechanics settled by D-4; the
  cluster's confirmation of the two-population reading is still owed.
- No previously `Verified` assumption was invalidated. A2's derivation consumes
  finite declared domains and is unaffected by the projection clause; A1's
  harness evaluates over tag *values* and never projects onto a product.

## Tiebreakers escalated

None. Both genuine forks were collapsed with evidence per the
tiebreaker-reduction gate: D-1's follow-on projection fork on the RDR's own
operator matrix (`contains` is the co-occurring-value operator), and D-5's
population question on `resolve.go::escapeOrRefuse`.

## Convergence

Ledger entries D-1..D-5: all **fixed**, none open, none charted. The D-1
follow-on is in-scope — it is the same contract D-1 names (the scoped product's
arithmetic), reached by computing D-1's own fix, not net-new scope. Converged at
iteration 3 of the lens's 3-iteration cap; the fixes are pins to existing
clauses rather than a rewrite, so no re-run of the lens is owed.
