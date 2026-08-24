Model: claude-opus-5[1m]

# Tooling Pass — RDR 0002 (Stage 7 mechanical pre-sweep, run 2)

Post-mutation regression sweep for the **re-lock**. Since run 1 (which preceded
the 2026-08-23 lock) `/rdr-cluster-reconcile` iter-3 demoted this RDR to Draft
(SPEC-DEFECT, STAGE-SCOPED, re-verify A1/A9/A12), the draft was re-authored to
the JDR 0001 §D7 closed layout, and Stage 3 refine → Stage 4 resolve → cove
iter-2 → 3amigo iter-3 → critique iter-3 → Stage 6 reconcile iter-2 all ran.
The assumption set grew from 14 to 15 (A15 added at Stage 6) and A11/A13/A14
flipped to `Verified`. Findings only; the mechanical share is fixed in-pass.

## CHECK 1 — Template section coverage  (primary signal)

Every **Required** spine section of `TEMPLATE.md` is Present-substantive:

| Template section | RDR | Verdict |
| --- | --- | --- |
| `## Metadata` | :6 | Present-substantive (Profile carries the value + naming clause) |
| `## Problem Statement` | :25 | Present-substantive |
| `## Critical Assumptions` | :96 | Present-substantive — 15 Evidence Records |
| `## Proposed Solution` / `### Approach` / `### Technical Design` | :466/:468/:495 | Present-substantive |
| `#### Normative Contracts` | :599 | Present-substantive |
| `#### Load-Bearing Decisions` | :1466 | Present-substantive |
| `#### Round-Trip / Inverse Invariants` | :1523 | Present-substantive |
| `#### Illustrative Code` | :1660 | Present-substantive |
| `### Capability Dependencies` | :1759 | Present-substantive |
| `### Existing Infrastructure Audit` | :1774 | Present-substantive |
| `### Decision Rationale` | :1784 | Present-substantive |
| `## Alternatives Considered` | :1833 | Present-substantive — 5 named + Briefly Rejected |
| `## Context` (Background, Technical Environment) | :29 | Present-substantive |
| `## Research Findings` (Investigation, Key Discoveries) | :41 | Present-substantive |
| `## Trade-offs` (Consequences, Risks and Mitigations, Failure Modes) | :1950 | Present-substantive |
| `## Implementation Plan` (Prerequisites, MVV, Phases 1–5, Day 2, New Dependencies) | :2033 | Present-substantive |
| `## Validation` (Testing Strategy, Performance Expectations) | :2242 | Present-substantive |
| `## Finalization Gate` | :2304 | see below |
| `## References` | :2308 | Present-substantive — real citations, no placeholder |

`#### Conditional Mini-Checks` (:1587) is an extra substantive section
(fidelity / disposition / oracle / trace tables), not a template section — no
finding.

The C1 finding run 1 raised — a surviving 43-line template guidance block in
`## Critical Assumptions` — **stays fixed**; the section opens directly at A1
with no vocabulary list, self-reference rule, or exactness-claim paragraph. The
re-authoring did not reintroduce it. Grep for surviving verbatim template
brackets (`[Required`, `[Conditional`, `[Resource]`, `[Capability]`),
`_Draft placeholder._`, and a seed-skeleton header returns **zero hits**.

- **C1 — `## Finalization Gate` carries a STALE lock pointer.** The section
  body reads
  `Responses: 0002-transition-table-as-reviewable-data/artifacts/gate.md (Gate PASS 2026-08-23)`
  — the pointer written at the *previous* lock. That lock was undone by
  cluster-reconcile iter-3, and `artifacts/gate.md` still holds the 2026-08-23
  responses, which describe a 14-assumption pre-§D7 draft that no longer
  exists (no A15; A11/A13/A14 recorded `Pending`). The pointer is not hollow,
  but it names the wrong record for the current draft. **MECHANICAL** — the
  READY action overwrites `gate.md` and rewrites the pointer date in the same
  pass, which is exactly the prescribed lock procedure; no separate fix or
  re-run is owed. Recorded so the staleness is not mistaken for a current
  record.

## CHECK 2 — Method-label vocabulary

15 Evidence Records (A1–A15). Every Method is exactly one of the eight
sanctioned labels:

| Records | Method |
| --- | --- |
| A1, A3, A6, A7, A11, A13, A14, A15 | `Spike` |
| A2, A5, A8 | `Source Search` |
| A4 | `Peer RDR` |
| A9, A10, A12 | `MVV Test` |

No missing, paraphrased, or off-vocabulary label. Note the shift since run 1:
A11, A13, and A14 moved from `MVV Test` to `Spike` when Stage 4 / Stage 6
promoted them on executed evidence, and A15 entered as `Spike`. Both are
legitimate relabels backed by transcripts, not drift. **PASS.**

## CHECK 3 — Source Search self-reference

Three `Source Search` records. A2 cites `internal/resolve/resolve.go::Resolve`,
`::gate`, `::escapeOrRefuse`, `::missingOwned`, and
`adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`.
A5 cites `internal/cli/clierr/clierr.go::CLIError`,
`internal/cli/respond/respond.go::Fail`, `internal/cli/config/config.go::Load`.
A8 cites `internal/resolve/resolve.go::Resolve`, `::escapeOrRefuse`,
`::assemble`, `::Table.models`, `const recognizedTagKey`. None resolves to
`{RDR_PATH}` or to any path under this RDR's artifact directory. **PASS.**

(A1/A3/A6/A7/A11/A13/A14/A15 cite `evidence/spikes/iter-2/` paths and A10 cites
`evidence/spikes/iter-2/main.go::Model` — all under `Spike` / `MVV Test`, which
C3 does not govern.)

## CHECK 4 — Docs Only on load-bearing claims

No record uses `Docs Only`. **PASS.**

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Delegated symbol resolution against the working tree; every cited symbol
resolves in the file it is cited in:

| Anchor | Resolves |
| --- | --- |
| `internal/resolve/resolve.go::Row` | `resolve.go:169` |
| `::Resolve` | `resolve.go:318` |
| `::assemble` | `resolve.go:148` |
| `::gate` | `resolve.go:376` |
| `::escapeOrRefuse` | `resolve.go:473` |
| `::missingOwned` | `resolve.go:444` |
| `::Table.models` | `resolve.go:216` |
| `const recognizedTagKey` | `resolve.go:110` |
| `adversarial_test.go::TestAdv3_Guard…RowOrder` | `adversarial_test.go:266` |
| `internal/cli/clierr/clierr.go::CLIError` | `clierr.go:48` |
| `internal/cli/respond/respond.go::Fail` | `respond.go:133` |
| `internal/cli/config/config.go::Load` | `config.go:74` |

A5's claim about `config-invalid` holds precisely: `Load` emits
`config-not-found` (`config.go:79`) and `config-read-error` (`config.go:85`) as
live `CLIError` codes, while `config-invalid` exists only as a doc-comment
mention (`config.go:72`) and a `TODO` (`config.go:100`) with no construction
site anywhere under `internal/`.

The **negative** claims Prerequisites and A9/A12 rest on also hold on the
working tree: zero hits for an `Atom` type, zero for `OpExists` /
`LiteralTrue` / `LiteralFalse` in non-evidence Go source, `Row.Guard` still
`string` (`resolve.go:186`), `Row.Match` still a flat `[]Tag`
(`resolve.go:180`). RDR 0007's reshape is genuinely unimplemented, exactly as
the RDR states.

No bare `file:line` anchor in any Evidence field. No phantom symbol. **PASS.**

## CHECK 6 — Status consistency

Three records are `Pending` (A9, A10, A12) — down from six at run 1, since
Stage 4 promoted A11/A14 and Stage 6 closed A13 (and closed the newly-added
A15 in the same pass rather than deferring it).

No settled-fact prose leans on an unsettled assumption. The three sites where a
Pending assumption's subject appears outside its own record were read in full:

- **:1170–1185** (`Match`/`Guard` routing, A9/A12) — states the handoff
  obligation in `MUST` form ("the handoff MUST route each atom to exactly one
  of them by its block"). The one factual claim inside it — `Guard string`
  carries no per-atom block — is source-verified above, not assumed.
- **:1330–1342, :1430–1440, :1601** (`duplicate model id`, A10) — all three
  are `MUST` / `MUST NOT` clauses this RDR is *making* about a loader that does
  not exist yet, plus the source-verified singular-`[model]` decidability fact.
  None asserts observed behavior of shipped code.
- **:2259** (A12's oracle) — Testing Strategy, and it states its own blocker
  ("unsatisfiable until RDR 0007's reshape lands").

Prerequisites (`:2037`) names the actual `Pending` set (A9, A10, A12) and its
count, and the `[x]` checkbox is qualified "**except A9, A10, and A12**" —
the run-1-era mismatch (checkbox `[x]` while five were `Pending`) was corrected
at Stage 6 and does not recur. No checklist-vs-gate disagreement: the gate
record is being rewritten in this same pass.

The RDR's `Status: Draft` line carries no 07.1 qualifier while the README index
row reads `Draft [revised from Final 2026-08-24; re-verify A1, A9, A12 —
STAGE-SCOPED, re-enter Stage 3]`. Per TEMPLATE.md that qualifier self-clears at
the Stage 7 flip, which overwrites the whole value. Its named re-verification is
discharged: A1 re-verified at Stage 4 against the §D7 layout, A9 and A12
re-examined at Stage 6 and deliberately DOWNGRADED to `Pending` with named MVV
oracles and stated blockers. **PASS.**

## CHECK 9 — Evidence-field budget  (ADVISORY, never blocks)

| Assumption | Evidence lines |
| --- | --- |
| A13 | 60 |
| A15 | 33 |
| A7 | 31 |

**3 fields over budget, 124 lines.** Up from 1 at run 1 — and the growth is
where Stage 6 did its work. Reviewed per assumption as the check asks:

- **A13 (60)** — the mass is the Stage 6 close: the defect as found
  (`Atom.identity` → `literalString` → `strings.Join`), the fix
  (`memberKey`, `Write.identity`, `renderMembers`), the discriminating
  witnesses, and the five-delimiter parameterization that twice earned its keep.
  Load-bearing anchors (`main.go::Atom.identity`, `::memberKey`,
  `::renderValue`, the `delim/` fixture dir, `negative-cases.txt`) stay findable
  within it. The narrative *is* the verification; a pointer-only field would
  leave "why five delimiters" unrecoverable, and that is the reasoning the
  grounding sweep reads. Keep as written.
- **A15 (33)** — same shape: the guard-block gap as reproduced, the
  `checkMatchBlock` → `checkAtomBlock` generalization, the two rules that stay
  legitimately match-only, the per-operator domain/kind scoping (a real design
  refinement the pass forced), and the six minted controls. Anchors resolve;
  the repro is preserved at `evidence/reconcile/iter-2/a15-guard-block-repro.txt`.
  Keep.
- **A7 (31)** — one line over. Unchanged in character from run 1's report; the
  mass is the determinism checklist and the permutation controls.

Advisory report only; contributes nothing to the verdict. No truncation or
relocation proposed.

## Cluster re-entry note

`grep -n "Refinement Context"` returns **zero hits**. The
`## Refinement Context (cluster re-entry — delete on re-lock)` block that
cluster-reconcile iter-3 left was cut at Stage 3 refine (commit `3c3cb3e`,
"cut Refinement Context"), and its defect was folded into live text: the §D7
layout restatement, block-keyed routing, and the general match-block expansion.
**PASS.**

## Repeatability requirement (Profile `foundational`)

Normative Contracts name parse/normalize/dump, identity, total ordering, and
multi-step MVV fidelity, so repeatability is owed and present:
`evidence/repeatability/run-1.md` (`claude-opus-5[1m]`), `run-2.md`
(`Claude Sonnet 5`), `run-3.md` (`claude-fable-5`) — all
`variant: full (profile: foundational)` — plus `diff.md` (health **Healthy**)
and `dispositions.md` (R-1..R-9, C-1..C-7 terminal). **Satisfied.**

**Caveat carried forward, not a return.** The reconcile report records that
these files predate the §D7 re-authoring while the other three lenses ran after
it. §lens-row's completion test is the run/diff files at the right variant with
distinct stamps, which is met, so the row is satisfied and this is not a
Stage-5 return. Recorded here because a reconstruction over the *pre*-§D7 draft
is weaker evidence about the current draft than folder presence implies; the
reconcile's absorption audit confirmed no clause it touched was re-opened by the
later rounds.

---

## Verdict

**PASS — no blocking findings; proceed to the Gate's written responses.**

C1's stale lock pointer is MECHANICAL and is resolved by the READY action's own
prescribed steps (overwrite `gate.md`, rewrite the pointer line) rather than by
a separate fix-and-re-run. C9's three over-budget fields are advisory and
contribute nothing to the verdict.
