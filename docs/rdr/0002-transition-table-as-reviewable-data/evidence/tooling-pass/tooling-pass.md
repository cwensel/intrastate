Model: claude-opus-5[1m]

# Tooling Pass — RDR 0002 (Stage 7 mechanical pre-sweep, run 1)

Post-mutation regression sweep run after the cove / 3amigo / critique (dual-model)
/ repeatability (×3) rounds and the Stage 6 reconcile. Findings only; the
mechanical share is fixed in-pass and the sweep re-run.

## CHECK 1 — Template section coverage

Every **Required** spine section of `TEMPLATE.md` is Present-substantive:
Metadata, Problem Statement, Critical Assumptions, Proposed Solution (Approach,
Technical Design, Normative Contracts, Load-Bearing Decisions, Round-Trip /
Inverse Invariants, Illustrative Code, Capability Dependencies, Existing
Infrastructure Audit, Decision Rationale), Alternatives Considered (5 named +
Briefly Rejected), Context, Research Findings, Trade-offs (Consequences, Risks
and Mitigations, Failure Modes), Implementation Plan (Prerequisites, Minimum
Viable Validation, Phases 1–5, Day 2 Operations, New Dependencies), Validation
(Testing Strategy, Performance Expectations), Finalization Gate, References.

`#### Conditional Mini-Checks` is an extra substantive section (fidelity /
disposition / oracle / trace tables), not a template section — no finding.

The Finalization Gate section still carries the template's instructional text
and its five bracketed sub-sections. This is the **pre-lock** state the template
specifies ("At lock, replace this section's body with the one-line pointer to
gate.md"), not a hollow section — no finding.

- **C1 — `## Critical Assumptions` (lines 311–353): surviving template guidance
  text.** The instance carries a verbatim copy of the eight-Method vocabulary
  list, the Source-Search self-reference rule, and the exactness-claim rule.
  `README.md` (engine root) states this explicitly: *"This section is the
  authoritative Method vocabulary — the TEMPLATE's Evidence Record points here;
  guidance never ships inside the template body, because the template is copied
  verbatim to make each instance."* The current `TEMPLATE.md` does not ship the
  block; an older template version did, and this instance carries it forward.
  Peers 0004, 0007, 0008, and 0009 have it removed. **MECHANICAL** — delete in
  this pass, re-run the sweep.

## CHECK 2 — Method-label vocabulary

14 Evidence Records (A1–A14). Every Method is exactly one of the eight
sanctioned labels:

| Records | Method |
| --- | --- |
| A1, A3, A6, A7 | `Spike` |
| A2, A5, A8 | `Source Search` |
| A4 | `Peer RDR` |
| A9, A10, A11, A12, A13, A14 | `MVV Test` |

No missing, paraphrased, or off-vocabulary label. **PASS.**

## CHECK 3 — Source Search self-reference

Three `Source Search` records. A2 cites
`internal/resolve/resolve.go::Resolve`, `::gate`, `::escapeOrRefuse`,
`::missingOwned`, and
`adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`.
A5 cites `internal/cli/clierr/clierr.go::CLIError`,
`internal/cli/respond/respond.go::Fail`, and
`internal/cli/config/config.go::Load`. A8 cites
`internal/resolve/resolve.go::Resolve`, `::escapeOrRefuse`, `::assemble`,
`::Table.models`, and `const recognizedTagKey`. None resolves to
`$RDR_PATH` or to any path under this RDR's artifact directory. **PASS.**

(A1/A3/A6/A7 cite `evidence/spikes/` paths, and A10 cites
`evidence/spikes/main.go::Model` — all under `Spike` / `MVV Test`, which C3 does
not govern.)

## CHECK 4 — Docs Only on load-bearing claims

No record uses `Docs Only`. **PASS.**

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Every cited symbol resolves on the working tree:

| Anchor | Resolves |
| --- | --- |
| `internal/resolve/resolve.go::Row` | `resolve.go:169` `type Row struct` |
| `::Resolve` | `resolve.go:318` |
| `::assemble` | `resolve.go:148` |
| `::gate` | `resolve.go:376` |
| `::escapeOrRefuse` | `resolve.go:473` |
| `::missingOwned` | `resolve.go:444` |
| `::Table.models` | `resolve.go:216` |
| `const recognizedTagKey` | `resolve.go:110` |
| `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder` | `adversarial_test.go:266` |
| `internal/cli/clierr/clierr.go::CLIError` | `clierr.go:48` |
| `internal/cli/respond/respond.go::Fail` | `respond.go:133` |
| `internal/cli/config/config.go::Load` | `config.go:74` |
| `config-not-found` / `config-read-error` / `config-invalid` | `config.go:61,72,85` — matches A5's claim that `config-invalid` exists only as a doc comment |

No bare `file:line` anchor in any Evidence field. No phantom symbol. **PASS.**

The Prerequisites' negative claims also hold: no `Atom` type, no `OpExists` /
`LiteralTrue` / `LiteralFalse` in the repo, `Row.Guard` is a `string` and
`Row.Match` a flat `[]Tag` — the reshape RDR 0007 owes is genuinely
unimplemented, exactly as the RDR states.

## CHECK 6 — Status consistency

Six records are `Pending` (A9–A14), every one a deliberate Stage 6 DOWNGRADE
with a named MVV assertion, enumerated per-assumption in Prerequisites and
carried into Testing Strategy scenarios 2 and 3 with the owed fixture named.
No settled-fact prose leans on an unsettled assumption: A11's escape-expansion
clause, A12's routing totality, A13's merge idempotence, and A14's version-gate
precedence appear in Normative Contracts as **normative obligations** (MUST/MAY
clauses this RDR is *making*), not as verified statements about existing code,
and each is paired with its owed oracle in Testing Strategy. A9/A12 are
additionally flagged unsatisfiable until RDR 0007's reshape lands.

The RDR's `Status: Draft` line carries no 07.1 qualifier while the README index
row still reads `Draft [revised from Final 2026-08-12; re-verify A2, A7]`. Per
TEMPLATE.md the qualifier self-clears at the Stage 7 flip to `Final`, which
overwrites the whole value and updates the row — not a contradiction, and the
re-verification it named is discharged in Prerequisites ("A2, A7 re-verified;
A8 resolved"). **PASS.**

## CHECK 9 — Evidence-field budget (ADVISORY, never blocks)

| Assumption | Evidence lines |
| --- | --- |
| A8 | 32 |
| A7 | 22 |
| A2 | 15 |

**1 field over budget, 32 lines.** A8's mass is the residual-scope analysis —
the alphabet ban is a load-time check in the normalizer, the kernel imports no
normalizer, so the barrier covers exactly one producer. That is load-bearing
reasoning the grounding sweep reads, and the anchors
(`::Table.models`, `::escapeOrRefuse`, `::assemble`, `recognizedTagKey`) stay
findable within it. No truncation proposed; no relocation warranted. Advisory
report only.

## Cluster re-entry note

No `## Refinement Context (cluster re-entry)` block survives in the RDR.
**PASS.**

## Repeatability requirement (Profile `foundational`)

The Normative Contracts name parse/normalize/dump, identity, total ordering,
and multi-step MVV fidelity, so repeatability is owed and present:
`evidence/repeatability/run-1.md` (`claude-opus-5[1m]`), `run-2.md`
(`Claude Sonnet 5`), `run-3.md` (`claude-fable-5`) — all
`variant: full (profile: foundational)` — plus `diff.md` (health verdict
**Healthy**) and `dispositions.md` (R-1..R-9, C-1..C-7 all terminal).
**Satisfied.**

---

## Verdict

**BLOCK — 1 finding (C1), MECHANICAL.**

Surviving template guidance text in `## Critical Assumptions`. Fixed in this
pass per the Stage 7 contract (conformance is outside Refine's contract); sweep
re-run below.

C9's over-budget report does not contribute to the verdict.

---

# Re-run (after mechanical fix)

The 43-line guidance block (`**Method vocabulary**` header through the
exactness-claim paragraph) was deleted from `## Critical Assumptions`. The
section now ends with A14's `If wrong` line and flows directly into
`## Proposed Solution`. No Evidence Record, Method label, Status, or anchor was
touched, so C2–C6 and C9 are unchanged by construction.

- CHECK 1 — no surviving template instructional text outside the pre-lock
  Finalization Gate body. **PASS.**
- CHECK 2 — **PASS** (unchanged).
- CHECK 3 — **PASS** (unchanged).
- CHECK 4 — **PASS** (unchanged).
- CHECK 5 — **PASS** (unchanged).
- CHECK 6 — **PASS** (unchanged).
- CHECK 9 — 1 field over budget, 32 lines. Advisory.

**PASS — no blocking findings; proceed to the Gate's written responses.**
