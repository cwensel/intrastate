Model: claude-opus-5[1m]

# Charted to successor — 3amigo lens, RDR 0008 (iteration 1)

Net-new scope surfaced by the lens, recorded here rather than absorbed into
this RDR (the scope-expansion wormhole guard). Each line: finding + why
out-of-scope + suggested successor.

## C-1 — Advisory lint heuristic for *intent* mismatch (origin: PM-3)

**Finding.** The reserved-key rule closes the naming channel but not the
intent channel. An author who declares `[tags.outcome]` with
`provenance = "observed"` while intending recognizer semantics passes every
check this RDR adds, and their row still silently never fires. PM-3 argues
this is plausibly the *majority* of the Problem Statement's population, not a
narrow residual — supported by this RDR's own A1 evidence, where both
canonical RDR 0002 fixtures name the recognized-provenance tag `outcome`.

**Why out-of-scope here.** This RDR's single contract is the *identity of the
recognized-provenance tag key* — who owns the name. An intent heuristic is a
different rule over a different signal: declaration *shape* (a model whose
rows gate on outcomes but which declares no recognized-provenance tag at
all). It carries a false-positive profile this RDR has no evidence to size,
and would be advisory rather than a validation failure. Folding it in would
add a second independent contract and trip the Normative Contracts split
signal.

**Suggested successor.** A new RDR at the same `area:internal-resolve` seam,
scoped to advisory authoring diagnostics over the normalized model — probably
seeded after RDR 0002's normalizer exists, since the heuristic needs the
declaration corpus to tune against. Failure Modes now states the boundary
explicitly so a reader does not mistake silence for coverage.

## C-2 — A user-facing surface for data-level validation failures (origin: PM-2)

**Finding.** The RDR claims the fix "converts the seed's silent no-match
discovery into an authoring-time diagnosis," but names no command the author
runs that surfaces a `reserved_tag_key` failure. RDR 0005's stable code-string
table (`0005…md:599-616`) contains no entry for any RDR 0002 data-level
validation category — its rows are CLI-argument and kernel-refusal failures.
RDR 0006 owns a lint command surface but is explicitly held at arm's length in
Phase 2.

**Why out-of-scope here.** The gap is real but belongs to the CLI layer, not
to this RDR's contract. This RDR's obligation is that the failure *data*
carries the guidance (block 3); which verb renders it, and under which exit
code, is RDR 0005's mapping decision — and 0005's table is an explicitly
extensible "Minimum" floor, so the addition is additive there exactly as the
category is additive in 0002. Deciding it here would author a contract this
RDR does not own.

**Suggested successor.** A cluster-reconcile watch item against RDR 0005: when
0002's validation categories become mappable, confirm `reserved_tag_key` gets
a stable code string and reaches a verb the author actually runs. Recorded as
a needs-verification item for Stage 6/7.1 rather than a new RDR, since it is
plausibly a one-row addition to an existing extensible table.
