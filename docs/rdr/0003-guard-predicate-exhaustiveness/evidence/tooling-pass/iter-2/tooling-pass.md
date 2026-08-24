Model: claude-opus-5[1m]

# Tooling Pass — RDR 0003 Guard Predicate Exhaustiveness (iteration 2, re-lock)

Mechanical adherence sweep, run as the Stage 7 pre-step to the Finalization
Gate. This is a **re-lock sweep**: 0003 locked Final 2026-08-22, was demoted to
Draft by cluster-reconcile iter-3 (`567600f`, RE-LOCK-ONLY / `re-verify none`)
over one mis-routed fenced clause, and was refined by `aa39d86`. So this is a
post-mutation regression check over what those two commits changed.

## Loop-breaker

Prior report: `evidence/tooling-pass/tooling-pass.md` (2026-08-22), which closed
**PASS** after fixing two MECHANICAL findings in-pass — a missing greppable
`Premortem:` verdict line, and four bare `file:line` anchors in A5.

**Both prior fixes survive `aa39d86` intact**: `Joint-check: fired →` and
`Premortem: survived.` are both present, and A5's four anchors are still in
stable `path::Symbol` / `NNNN::Section` form. **No finding in this round is a
re-report of a prior finding after its named stage ran — no ROUTING-LOOP.**

Findings 2–4 below are the *same finding class* the prior round raised for A5
(C5 bare-line form), but on **different records the prior round did not name**,
so they are fresh instances rather than a routing loop.

## Round 1 — findings

- **C1 / `## Critical Assumptions` — surviving template instructional text
  (43 lines).** A verbatim block from an **older TEMPLATE.md** sat between the
  A21 record and `## Proposed Solution`: the `**Method vocabulary** (pick
  exactly one per assumption)` legend enumerating all eight Methods, plus the
  `Method: Source Search` self-reference rule and the exactness-claim paragraph.
  Confirmed template-origin, not authored content: the engine repo shipped it in
  `f0fa020` and removed it in `bfb6253` (`rdr/TEMPLATE.md`), and it does not
  occur in the current TEMPLATE.md. Precedent: **RDR 0002's own iter-2 re-lock
  sweep raised and removed the identical block** ("a surviving 43-line template
  guidance block in `## Critical Assumptions`"), and 0002 now opens directly at
  A1. **MECHANICAL** — surviving template text, fixed in-pass per the Stage 7
  split; never routed to Refine, whose contract excludes conformance.
- **C5 / A13 Evidence — bare `file:line` anchor.**
  `docs/rdr/0007-guard-predicate-totality.md:1585-1595` carried no symbol or
  section anchor. The cited content resolves (0007:1583-1597 is the
  canonicalization-duty passage assigning the set-literal spelling to this RDR
  and naming `in`), so this is the *form* finding, not a phantom.
  **MECHANICAL** — C5 anchor rewrite.
- **C5 / Technical Design — bare `file:line` anchor.**
  `docs/rdr/0007-guard-predicate-totality.md:1249-1266` for the SEAM
  parsed-atom clause. Resolves (0007:1248 `#### Normative Contracts`; :1249-1266
  is exactly the `SEAM.` normative block). A5 already cites this same clause in
  the stable form, so the anchor was in the document and the body had not
  adopted it. **MECHANICAL** — C5 anchor rewrite.
- **C5 / `disposition` table — two bare line anchors.** RDR 0007 `` `:1412-1414` ``
  and `` `:1435-1436` ``. Both resolve to `0007::Normative Contracts` — the
  strong-Kleene atom-verdict block and the row-verdict
  `all_result ∧ ¬(unless_conj)` clause. **MECHANICAL** — C5 anchor rewrite.
- **C5 / A14 Evidence — stale verbatim quotations of RDR 0002.** Two of the
  record's three quoted strings were no longer RDR 0002's text. A14 quoted
  "Normalization MUST combine **both** into one candidate-row predicate set …
  the block (`all` or `unless`) it was authored in"; 0002 now reads "combine
  **the match atoms and both guard blocks** … and the block it was authored in
  — **a three-valued domain, `match`, `all`, or `unless`**". The second quote
  ("Block retention is carriage: …") did not exist in 0002 in that wording.
  The **claim A14 makes still holds** (block retained; no folding), so the
  `Verified` stamp was not falsified and this raised no C6 finding — but a Peer
  RDR record's evidence *is* its quotation. The substantive change behind the
  drift is RDR 0002's §D6-authorized widening of the block domain from
  two-valued to three-valued; 0002 explicitly declares the widening as its own
  and instructs implementers reading 0007 alone to "widen to three here", so
  this is a peer-owned widening, not a cross-document contradiction.
  **MECHANICAL** — re-quote from the peer's current text.

**Round 1 verdict: BLOCK — 5 findings.** All five **MECHANICAL**, none
SUBSTANTIVE, so no stage return is owed and no stage prompt owns a fix here.
All were fixed in this pass and the sweep re-run, per the Stage 7 split.

## Fixes applied in-pass

- `Critical Assumptions` — the 43-line older-template block deleted; the section
  now ends at A21's `If wrong` and `## Proposed Solution` follows directly,
  matching RDR 0002's post-fix shape. No authored content removed: every line
  deleted is template instructional prose, and the rules it stated are engine
  doctrine enforced by this sweep, not RDR content.
- `A13` Evidence — `…0007-guard-predicate-totality.md:1585-1595` → ``0007::Normative
  Contracts``, the literal-spelling condition the totality claim rests on.
- `Technical Design` — `…0007-guard-predicate-totality.md:1249-1266` →
  ``0007::Normative Contracts``, the SEAM parsed-atom clause (the stable form A5
  already uses for the same clause).
- `disposition` table — RDR 0007 `` `:1412-1414` ``, `` `:1435-1436` `` →
  ``0007::Normative Contracts``, the strong-Kleene atom-verdict clause and the
  row-verdict clause.
- `A14` Evidence — both stale quotations replaced with RDR 0002's current text,
  verified verbatim-exact by whitespace-normalized string match against
  `docs/rdr/0002-transition-table-as-reviewable-data.md`. One clause added
  recording that the third `match` member is 0002's own §D6 widening and that
  the guard-side carriage this record depends on is unaffected.

## Round 2 — re-run after fixes

- **C1 Template section coverage — PASS.** Every **Required** (spine) section
  Present-substantive: Metadata, Problem Statement, Critical Assumptions,
  Proposed Solution (Approach / Technical Design / Normative Contracts),
  Decision Rationale, Alternatives Considered (six blocks + Briefly Rejected),
  Context, Research Findings, Trade-offs, Implementation Plan (Prerequisites,
  MVV, Phases 1–4), Validation, Finalization Gate, References. Conditional
  sections present are substantive: Load-Bearing Decisions, Round-Trip /
  Inverse Invariants, Illustrative Code, Capability Dependencies, Existing
  Infrastructure Audit, Day 2 Operations, New Dependencies, Performance
  Expectations. Zero `TBD`, zero `_Draft placeholder._`, zero `seed skeleton`
  header, zero surviving verbatim template brackets, and the round-1 legend
  block is gone. Both greppable verdict lines present. The Finalization Gate
  section carries the at-lock one-line pointer to `artifacts/gate.md` — the
  record, not a hollow section. **PASS.**
  - *Structural note, non-blocking:* `Critical Assumptions` is emitted at `###`
    under `## Research Findings`, where TEMPLATE.md places it at `##`. This is
    the cluster-wide house form (RDR 0002, 0006, 0007 all do the same) and the
    section is Present-substantive, so it is not a C1 finding.
- **C2 Method label vocabulary — PASS.** 21 records, every Method exactly one of
  the sanctioned eight. Census: Spike (A1), Derivation (A2), Design Decision
  (A3, A7, A9, A11, A13), Source Search (A4), Peer RDR (A5, A6, A8, A10, A12,
  A14, A16, A17, A18, A19, A20), MVV Test (A15, A21). `aa39d86` flipped four
  Statuses (A10, A12, A17, A19 → `Verified`) and relabeled no Method; no record
  added or removed. **17 Verified / 4 Pending / 0 Unverified.**
- **C3 Source Search self-reference — PASS.** One Source Search record (A4);
  its three anchors are consumer-repo paths, none under this RDR or its artifact
  directory.
- **C4 Docs Only — PASS.** Zero `Docs Only` records.
- **C5 Symbol resolution — PASS.** All four round-1 findings cleared. A4's three
  anchors resolve in the cited files (`clierr.go::CLIError`,
  `respond.go::Fail`, `config.go::Load`). A1's spike output exists
  (`evidence/spikes/iter-2/a1-eval-harness/{main.go,output.txt}`). Body source
  anchors resolve in `internal/resolve/resolve.go`: `Row` (carrying `RuleID`,
  `SourceLocator`), `Row.Guard` (still `string` — the Prerequisites sequencing
  claim is accurate), `Row.rescues`, `Resolve`, `Refusal`, `assemble`,
  `GuardEvaluator`. Negative claims verify (`conform` absent from the non-test
  kernel). **Anchors `aa39d86` introduced all resolve:** JDR 0001 §D6 (:257,
  closes JD-16) and §D7(iii) (:309, type-model wire keys, naming the two RDR
  0002 load categories with "0006 mints nothing for either" — exactly as A4 and
  the Status qualifier represent it). Also resolving: §D1, §D4, §JD-4, §JD-8,
  §JD-13, §JD-14, §JD-16, §JD-18. Peer-RDR section anchors resolve to the right
  section of the right document, including the four newly-closed records'
  clauses in RDR 0006 (row-group division, reachability relation, escape-row
  two-population overlap + `graph-coverage-closed-by-escape`).
- **C6 Status consistency — PASS.** Prerequisites checkboxes and record Statuses
  agree item for item: `[x]` on the 17 `Verified` (including the newly-closed
  A10/A12/A17/A19), `[ ]` on A15, A18, A20, A21 (all `Pending`). No record is
  `Unverified`. Settled-fact sweep over the four `Pending`: A15's bound adequacy
  is stated open inside the clause it gates and routed to MVV Scenario 3; A18's
  conformance premise is an explicit conditional with a "booked gap, not a
  silence" note; A20's MUST carries the runtime-half-open disclosure and the
  `authority` table row; A21 routes to MVV Scenario 4. The four newly-`Verified`
  records each dropped their `Evidence needed` / `Plan` scaffolding cleanly.
  **PASS.**
  - **Reported, not blocking:** `artifacts/gate.md` is the 2026-08-22 record and
    disagrees with the current document (it reads 13 `Verified` / 8 `Pending`;
    the document now reads 17 / 4). That is the expected consequence of the lock
    being demoted and re-refined — gate.md is overwritten by the Gate's
    written-responses step at this re-lock, so it is not authoritative here.
- **C9 Evidence-field budget (ADVISORY) — PASS, report only.** Longest fields:
  A17 (21), A16 (21), A12 (21), A14 (~26 after the re-quote), A8 (19), A1 (19).
  **0 fields over the 30-line budget, 0 lines.** `aa39d86` net-shrank the
  Evidence fields (450 insertions vs 632 deletions).

## Also checked

- **Cluster re-entry note — CLEAN.** No `## Refinement Context (cluster
  re-entry — delete on re-lock)` block survives; the string does not occur in
  the RDR. `567600f`'s re-entry material was folded into the Metadata Status
  qualifier (which self-clears at the Stage 7 flip to `Final`) and `aa39d86`
  removed the rest. Also clean: no `Stage 6 disposition`, `DOWNGRADED`,
  `BLOCKER`, or round-narration strings survive.
- **The demotion defect is discharged, bilaterally.** JDR 0001 §D7(iii) reads
  "Two load categories in 0002, rules supplied by 0003: `malformed tag
  declaration` … and `malformed predicate atom` … 0006 mints nothing for
  either; 0005 maps them under §JD-8." The RDR now routes rejections exactly
  there (A4 Evidence; the `disposition` table; the load-category passage), and
  RDR 0002 (`Final`) reciprocally names both categories as "rejection rules"
  that "RDR 0003's rules reject". §D6 is cited as decided where the
  participation / can-refuse clauses read the block.
- **§JD-16 correctly dropped from the qualifier** (JDR 0001 records "0003
  consistent, qualifier cleared"); **§JD-18 correctly retained** (open; A18,
  A20). §JD-14 states "0003 A17 closes on it", and A17 is now `Verified`.
- **Determinacy requirement — SATISFIED.** Profile is `large` and the trigger
  fires: `Normative Contracts` carry canonical-form and cross-implementation
  determinism claims (canonical unordered duplicate-free set-literal spelling
  entering the identity tuple; the published model-independent product bound;
  source-order independence). Evidence exists —
  `evidence/repeatability/iter-2/` and `iter-3/` each hold `run-1.md` +
  `diff.md` + `dispositions.md`, `variant: lite (profile: large)`, all five
  iteration-3 findings dispositioned `fixed`. No `determinacy: n/a` disposition
  is needed and none is written. The re-entry scope was `re-verify none`, so no
  lens re-run is owed.
- **Every Critical Assumption record internally consistent — PASS.** All 21 have
  Status/Method/Evidence in agreement and a non-empty "If wrong". The four
  `Pending` carry `Evidence needed` + `Plan`; the 17 `Verified` carry
  `Evidence`. No record mixes the two.

## Verdict

**PASS — no findings; proceed to the Gate's written responses.**

All five round-1 findings were MECHANICAL, fixed in-pass, and cleared on re-run.
No SUBSTANTIVE finding was raised, so no stage return is owed.
