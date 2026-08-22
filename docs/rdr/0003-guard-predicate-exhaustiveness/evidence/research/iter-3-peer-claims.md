Model: claude-opus-5[1m]

# Stage 4 Resolve (scoped re-entry, iter-3) — accepted citations

Date: 2026-08-22
Scope: A5 (the ID the re-entry qualifier names) + the anchors the
`0003-0006-0007` cluster gate's edits touched — A8 (§JD-4 closure), the
declaration model (§JD-13), and the escape-row clauses (§JD-14).

No corpus search was run, and none was owed: every in-scope claim is a
peer-document, JDR-home, or own-source claim. `StateMachineRes` /
`StateMachineLit` were not queried — the cluster gate changed no
external-behavior claim. `iter-2-peer-claims.md` remains valid for A4 and the
evaluator-scoping clause; this file records only what iter-3 re-verified or
newly cited.

negative: none — every in-scope claim landed on a local authority.

## Reuse audit (RDR_ENV paths)

Ran against `internal/cli/`, `internal/resolve/`, `cmd/`. No existing capability
duplicates what this RDR introduces:

- No tag declaration model, value-kind type, or finite-domain representation
  exists in `internal/` (grep: `TagDecl`, `Declaration`, `finite domain`,
  `declaredDomain` — zero hits).
- No `lint` verb exists in `internal/cli/` or `cmd/` (grep `"lint"` — zero
  hits), so RDR 0006's surface is unbuilt and nothing is being rebuilt here.
- `internal/resolve/resolve.go:89-93` carries `GuardEvaluator.Evaluate(guard
  string, view TagSet)` — the pre-reshape string seam. Confirms the A5 finding
  below rather than contradicting it.

Finding: nothing to fold in. This RDR builds on empty ground.

## A5 — re-verified (the named re-entry ID)

All four anchors resolve and none moved under the 2026-08-21 declaration-model
rehoming, which touched declaration *semantics* while A5 is about row
*carriage*:

- JDR 0001 §D1 — `docs/jdr/0001-resolve-kernel-seam.md:64` ("How does the guard
  seam carry a predicate?"), resolved to parsed atoms on the row.
- RDR 0007 `Normative Contracts` —
  `docs/rdr/0007-guard-predicate-totality.md:1252`: a row's guard is a slice of
  parsed atoms, "key, operator token, literal, block ∈ {all, unless}". No
  reconstruction step exists, so no mapping step can lose identity.
- RDR 0002 `Normative Contracts` —
  `docs/rdr/0002-transition-table-as-reviewable-data.md:317`: each normalized
  candidate row MUST "retain its source rule id and source locator".
- Shipped code — `internal/resolve/resolve.go:173-174` carries `RuleID` and
  `SourceLocator`, commented as the source identity RDR 0002 requires. The
  subject of A5 is therefore already true in built code; the unshipped half is
  RDR 0007's atom reshape (`resolve.go:185` still `Guard string`), which is a
  sequencing risk recorded in Prerequisites, not an A5 defect.

Verdict: **Verified**, unchanged. Evidence line pinned to the resolved
locations so a later sweep does not have to re-find them.

## A8 — CLOSED by the cluster gate (was the last lock blocker)

`docs/jdr/0001-resolve-kernel-seam.md` §JD-4 reads **CLOSED 2026-08-22**, with
both halves decided:

- **Recording assignment**: "RDR 0003 is the recording document; RDR 0006 cites
  it and mints no code." Ground given at the home: the clause must name the
  refusing atom, and `atom` occurs once in RDR 0006, only to delegate atoms to
  RDR 0003.
- **Warning category — no**: RDR 0006 reuses `graph-unprovable-coverage`, which
  already sits in its mandatory blocking table
  (`docs/rdr/0006-graph-lint-authority-and-guarantees.md:317`) and its MVV
  matrix (`:642`). Its advisory tier stays scoped to "redundant rows or
  unreachable rules" and MUST NOT absorb the withheld claim.
- **Corroboration from the peer itself**: RDR 0006's Status line (`:12`) records
  the same assignment — reuse `graph-unprovable-coverage` "with no new code and
  no non-blocking tier".

A8's plan said "close on the assignment, not on matching prose"; the assignment
exists at the home, so the record closes. The consequent duty (RDR 0006 extends
its finding contract with an atom-level field) is recorded at the home as RDR
0006's refine item and does not gate this RDR.

Gate evidence: `docs/rdr/cluster-reconcile/0003-0006-0007/reconcile-report.md`.

## §JD-13 — the single-valued marker (edit owed HERE, applied this pass)

`docs/jdr/0001-resolve-kernel-seam.md` §JD-13: RDR 0006 attributes
"single-valued grouping" to this RDR's declaration model and gates a mandatory
blocking code (`graph-single-valued-state`) plus an MVV assertion on it, while
the token occurred zero times here and zero times in RDR 0002. Verified this
pass: `single-valued` occurs 11× in RDR 0006, 0× in RDR 0002. Decision: this
RDR's declaration model gains the field.

Applied: the declaration-model clause now lists a single-valued marker; a new
normative clause states what it licenses (the domain is a partition, not
independent dimensions) and forbids inferring it from a name, a value spelling,
or a fixture; the domain/kind agreement clause rejects it on a `set` kind or a
kind with no finite domain. Swept sites: Approach alphabet, ownership census
(new row), MVV authorability list, Testing Strategy scenario 4.

## §JD-14 — escape rows (no edit owed here; verified)

§JD-14 decides "0003's reading governs". Checked against the clause in place:
this RDR already states that an escape row "participates in the coverage
identity like any other row", that lint MUST NOT treat its existence as a
coverage-satisfying fact outside the union, and MUST NOT exclude escape rows
from overlap checks. The clause matches the decision verbatim in substance; the
repair is RDR 0006's (its invariants 3 and 4). No edit applied, correctly.

## Rejected / corrected branches

- **"A8 still Pending because RDR 0006 has not yet written its citation" —
  REJECTED.** A8's own plan forbids this reading: closure is on the assignment,
  not on matching prose. Requiring RDR 0006 to carry the sentence would recreate
  the duplication §JD-4 exists to prevent. The consequent duty is a *payload
  extension* on RDR 0006, tracked at the home, not an open half of A8.
- **"§JD-13 could be discharged by citing RDR 0006's use of the field" —
  REJECTED.** That is the unhomed-producer defect itself: the consumer's
  attribution is what created the gap. Only a declaration clause here closes it.

## Still open after this pass (unchanged by it)

- **A10, A12** — tolerances on RDR 0006's refine; a peer cannot confirm a
  division of labour from inside this document. Both named at the home and on
  RDR 0006's Direction list.
- **A14** — per-atom `block` retention, a normalization-carriage request on RDR
  0002 (`Draft`). Deliberately not rehomed: carriage is RDR 0002's charter.
- **A15** — discharged by MVV Scenario 3 at implementation; not runnable at
  draft time because it needs finite products built and measured.
