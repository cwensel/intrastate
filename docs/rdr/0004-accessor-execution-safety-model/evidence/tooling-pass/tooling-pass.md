Model: claude-opus-5[1m]

# Tooling Pass — RDR 0004 (mechanical adherence sweep)

RDR: `0004-accessor-execution-safety-model`
Status line at entry: `Draft [revised from Final 2026-08-12; re-verified A8 — …]`
Run: iteration 1 for this re-entry (no prior report under `evidence/tooling-pass/`,
so the LOOP-BREAKER has no prior finding to compare against — no routing loop possible).

This is a **re-lock**: the RDR locked once at `c149031`, was demoted at `76a9e67`
when JDR 0001's decisions were applied across the 0002–0009 cluster, and has since
run refine → resolve → grounding → 3amigo → critique → reconcile.

## Verdict

**PASS** — after one MECHANICAL fix applied in-pass and the sweep re-run.

## CHECK 1 — Template section coverage

Every **Required** (spine) section is Present-substantive. No verbatim template
bracket (`[Required`, `[Conditional`, `[Resource]`, `[Capability]`), no
`_Draft placeholder._`, no `TBD`/"see above", no seed-skeleton header.
`Phase 2: Operational Activation` is a **Conditional** section cleanly omitted —
PASS per the false-positive guard, not Missing.

Stage-5 mini-check tables (Disposition Table, Oracle Discriminability, Fidelity
Table, Desk Trace) are present and substantive — legitimate instance additions.

### C1 finding — FIXED IN-PASS (mechanical)

- **C1 — "Method vocabulary" block, former lines 205–247 — surviving template
  instructional text. DELETED this pass; sweep re-run clean.**
  The 43-line block reproduced the eight-Method list from
  `$RDR_HOME/README.md:194–212` near-verbatim, plus three author-directed gate
  rules (Source-Search self-reference, symbol-resolves-on-`main`, exactness
  claims). It is guidance, not instance content: nothing in it is specific to
  0004. Three independent proofs it does not belong:
  1. The **current** `TEMPLATE.md` ships no such block (`grep 'Method vocabulary'`
     → no hit); it survived from an older template version — exactly the class the
     Stage 7 prompt names as MECHANICAL ("including blocks an older TEMPLATE.md
     shipped and the current one doesn't").
  2. `$RDR_HOME/README.md:196–198` states it directly: *"This section is the
     authoritative Method vocabulary — the TEMPLATE's Evidence Record points here;
     guidance never ships inside the template body, because the template is copied
     verbatim to make each instance."*
  3. The repo's own convention has already moved: RDRs 0007/0008/0009 (the three
     most recently authored) carry **no** such block; only 0001–0006 do.

  Fixed in-pass, not routed to Refine — conformance is outside Refine's contract,
  so the finding would return unchanged. No content was lost: the vocabulary is
  authoritative at `$RDR_HOME/README.md`, and every one of this RDR's ten Method
  labels remains in the sanctioned set (CHECK 2).

### C1 note — NOT a finding, recorded deliberately

- **`### Critical Assumptions` nested under `## Research Findings` (line 92),
  where TEMPLATE.md places `## Critical Assumptions` at top level between Problem
  Statement and Proposed Solution.**

  Considered and **not** raised as a BLOCK. This is a **repo-wide authoring
  convention**, not a regression this RDR's rounds introduced: *every* RDR in
  `docs/rdr/` carries it nested — 0001 (Implemented), 0002, 0003, 0004, 0005
  (Final), 0006 (Final), 0007, 0008 (Final), 0009 (Final). The tooling pass is a
  *post-mutation regression sweep*; reshaping the section spine on the one RDR
  that happens to be at lock would make 0004 structurally inconsistent with its
  own cluster while fixing nothing the rounds broke, and the content is fully
  substantive either way (A1–A10, each with Status/Method/Evidence/If wrong).

  The placement rationale is real (`$RDR_HOME/README.md:112–117` — contract-bearing
  sections lead because mid-document retrieval is weakest), so this is worth a
  cluster-wide normalization pass, but that is a separate decision across ten
  documents and is **not** this lock's business.

## CHECK 2 — Method-label vocabulary

**PASS.** All ten Critical Assumption records carry a Method from the sanctioned
set, exactly one each:

| Method | Count | Records |
| --- | --- | --- |
| Spike | 5 | A1, A2, A6, A7, A8 |
| Source Search | 1 | A3 |
| MVV Test | 3 | A4, A9, A10 |
| Design Decision | 1 | A5 |

No missing, paraphrased, or off-vocabulary label. No record gained or lost a label
during the rounds.

## CHECK 3 — Source Search self-reference

**PASS.** One Source Search record (A3). Its Evidence resolves to
`internal/cli/clierr::CLIError`, `internal/cli/clierr::ExitCodeFor`, and
`internal/cli/respond::Fail` — production source, not this RDR and not any path
under its artifact directory. The record also cites the spike transcript, but as
corroboration alongside the source anchors, not as the anchor itself. No offender.

## CHECK 4 — Docs Only on load-bearing claims

**PASS.** Zero records carry `Method: Docs Only`. (The only `Docs Only` strings in
the file were inside the deleted vocabulary block and the Finalization Gate
preamble — neither is an Evidence Record.)

## CHECK 5 — Symbol resolution of anchors

**PASS.** All 29 in-repo `path::Symbol` anchors resolve **in their cited file** —
no phantoms, no moves, no renames, no bare `file:line` anchors standing in for a
symbol. Verified by delegated read; full resolution table in that pass. Highlights:

- `internal/resolve/resolve.go::Tag` (`:20`) — exactly two fields `Key`/`Value`,
  independently corroborating A8's load-bearing claim that the absent/unread
  distinction cannot survive the seam type.
- `internal/resolve/resolve.go::Input.Owned` (`:224`/`:227`) — `[]Tag`, no error
  channel, as A8 and JDR §D3 both state.
- `internal/resolve/resolve.go::Refusal.MissingOwned` (`:269`) — doc comment reads
  "names owned tag keys absent from the snapshot", which the RDR quotes verbatim
  and correctly.
- `internal/resolve/adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
  (`:114`) — exact name match; the Disposition Table's conditional-loudness note
  cites real asserted kernel behavior, not a supposition.
- 14 spike anchors all resolve in `evidence/spikes/main.go`, including two
  non-function symbols (`refusal` type `:22`, `absentValue` const `:36`) — both
  legitimate `path::Symbol` anchors.
- External prior-art symbols (`StateRepresentation::ExecuteEntryActions`,
  `TaintCommand::Run`, `Client::transitionWorkItem`) do not resolve in this repo —
  **expected and out of scope**: they are cited as `Prior Art` in Investigation and
  References, never as `Source Search` evidence for an assumption.

Referenced documents exist and say what the RDR claims: `docs/cli-output-contract.md`;
`docs/jdr/0001-resolve-kernel-seam.md` §D3 (line 141, resolved **(b)** — complete
tag set or refuse, with option **(c)** explicitly kept available, which is exactly
what A9's "If wrong" cites as the fallback) and §JD-3 (line 188, `RequiresOwned`
has no obliged producer — listed as still open, which the RDR's Capability
Dependencies row characterizes correctly without overstating).

**Witness-citation accuracy.** All 23 cited `output.txt` lines match the actual
transcript content verbatim, including the two quotes the Desk Trace reproduces
inline. Two witness gaps are **disclosed rather than hidden**, and both are
accurate as written:
- The Disposition Table's "Definition invalid (8 shapes)" row cites `output.txt:17-23`
  (seven named codes) and annotates the eighth arm as "not witnessed — see A9".
- The "Gate returns deny" row cites `output.txt:2` with the parenthetical
  "(allow arm)"; deny is not separately witnessed. Self-labeled, not overclaimed.

## CHECK 6 — Status consistency

**PASS.** No settled-fact prose leans on an unsettled assumption. Every property
belonging to the two Pending assumptions is marked at its site: the Disposition
Table rows carry "not witnessed — see A9"/"see A10", the Oracle Discriminability
rows are tagged `(A9)`/`(A10)`, and Testing Strategy scenarios 6/7/8 name their
assumption. Scenario 2's prose explicitly separates what the spike witnesses from
what is "new and unwitnessed — they are A9's".

The A9/A10 clauses in Normative Contracts state **requirements this RDR
legislates**, not verified empirical facts — which is the correct register for an
unverified normative clause and does not constitute a settled-fact contradiction.

Checklist-vs-gate agreement: the Prerequisites checkbox ("A1-A8 Verified; **A9 and
A10 Pending, DOWNGRADED at Stage 6**") and the Assumption Verification section
state the same disposition. No disagreement.

## CHECK 9 — Evidence-field budget (advisory, never blocks)

**0 fields over budget.** Longest Evidence fields: A9 at 13 lines, A10 at 7 lines;
A1–A8 are 1 line each. Budget is 30 lines. Nothing to report, no truncation
proposed.

## Determinacy / repeatability trigger

**Discharged.** `Profile: large` and the Normative Contracts do touch
ambiguous-field-ownership-adjacent territory (absent vs unreadable key encoding),
so the trigger was correctly evaluated rather than skipped.
`evidence/grounding/iter-2/dispositions.md:40-45` carries the written disposition:

> `Determinacy trigger: **n/a**` — the contract legislates execution safety and
> refusal classing at an I/O boundary. It carries no parse/deparse, import/export,
> compose/decompose, hashing, identity, or migration behavior, no data-model field
> whose ownership could be read more than one way (tag ownership is fixed by RDR
> 0001's owned/observed/recognized split), and the MVV rests on per-call branch
> dispositions rather than multi-step transformation fidelity.

A written `determinacy: n/a — <reason>` disposition is exactly what the Stage 7
prompt accepts in place of `repeatability/run-1.md` + `diff.md`. No repeatability
row is owed, so no lens evidence is missing from the `large` row
(grounding → 3amigo → critique, all complete for iter-2).

## Cluster re-entry note

**Absent — PASS.** `grep 'Refinement Context'` → no hit. The 07.1 gate's
re-entry block (if one was ever present) is gone, and the cross-RDR defect that
demoted this RDR is dispositioned: the `Row.RequiresOwned` producer question is
recorded as an **Open** row in Capability Dependencies and routed to JDR 0001
§JD-3 / RDR 0002, not silently dropped.

## Findings summary

| Check | Verdict | Findings |
| --- | --- | --- |
| C1 template coverage | PASS (after fix) | 1 mechanical, fixed in-pass; 1 recorded note (repo-wide convention, not a regression) |
| C2 Method vocabulary | PASS | none |
| C3 Source Search self-reference | PASS | none |
| C4 Docs Only load-bearing | PASS | none |
| C5 symbol resolution | PASS | none (29/29 resolve; 23/23 witnesses accurate) |
| C6 status consistency | PASS | none |
| C9 evidence budget (advisory) | PASS | 0 over budget |

**PASS — no findings outstanding; proceed to the Gate's written responses.**
