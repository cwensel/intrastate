# Iteration-3 brief — cluster 0002-0009 (Stage 7.1 re-run)

Model: claude-fable-5
Date: 2026-08-24

Read-only against the RDRs. Sub-agents write ONLY their own output file under
this directory. Never edit an RDR, the JDR, or the README.

## Paths (absolute)

- Records dir: `./docs/rdr/`
  - 0002 `0002-transition-table-as-reviewable-data.md`
  - 0003 `0003-guard-predicate-exhaustiveness.md`
  - 0004 `0004-accessor-execution-safety-model.md`
  - 0005 `0005-skill-integration-cli-contract.md`
  - 0006 `0006-graph-lint-authority-and-guarantees.md`
  - 0007 `0007-guard-predicate-totality.md`
  - 0008 `0008-recognized-tag-key-ownership.md`
  - 0009 `0009-escape-row-shape-conformance-ownership.md`
  - (ignore `*-postmortem.md`; 0001 is Implemented — out of cluster)
- **Home**: `./docs/jdr/0001-resolve-kernel-seam.md`
- **Artifact of record for peer status**: `docs/rdr/README.md` Index table.
- **Ledgers**: iteration 1 `../report.md`; iteration 2 `../iter-2/reconcile-report.md`
  (its findings table + tolerance table are what this iteration traces to).
  Prior per-pair scans: `../pairwise-*.md` (iter-1), `../iter-2/pairwise-*.md`.
- Pairwise prompt: `.claude/skills/rdr-cluster-reconcile/pairwise.md`
- Whole-set critique prompt: `.claude/skills/rdr-cluster-reconcile/2-critique.md`

## What changed since iteration 2 (delta scope)

Member revisions recorded at iter-2 vs now:

| RDR | iter-2 rev | now | Moved? |
| --- | --- | --- | --- |
| 0002 | 5bed32a | 526c481 | Status qualifier only (1 line) |
| 0003 | f1e9a57 | 526c481 | Status qualifier only |
| 0004 | 48cb2a4 | 526c481 | Status qualifier only |
| 0005 | 76a9e67 | 76a9e67 | no |
| 0006 | ae2f478 | 526c481 | Status qualifier only |
| 0007 | 21322fb | 526c481 | 4 peer-status citation repairs (0007:124,159,1711,1713) |
| 0008 | ba95813 | 4581223 | **YES — Stage 4 re-entry + re-lock** (418+/301−) |
| 0009 | 76a9e67 | 526c481 | Status qualifier only |
| Home JDR 0001 | 21322fb | 4f9639c | **YES — §D5/§D6/§D7 decided → JD-15, JD-16, JD-17 ANSWERED** |

Therefore in scope:
1. **RE-SCAN** every pair with 0008: 0008-0002, 0008-0009, 0008-0005-0006, 0007-0008.
2. **Answered tolerances** — scoped answer-vs-fences checks for each sibling:
   JD-15 (§D5): 0002, 0004, 0009 (0006 reads by citation — check too).
   JD-16 (§D6): 0002, 0003.
   JD-17 (§D7): 0002, 0003, 0004, 0006.
3. **RE-SCAN** whole-set critique (≥1 member changed).
Everything else CARRIES from iter-2 (0007-000{2,3,4,5,6,9}, 0009-0005, 0004-0005,
0005-0006, 0002-0005, 0003-0005, 0003-0006 — JD-18 unanswered, no member moved).

## Task A — pairwise re-scan (per assigned pair)

1. Run `pairwise.md` verbatim against the two current files. Quote-anchored.
2. Ledger trace: for every iter-2 row for this pair (and iter-1 rows it carries),
   state RESOLVED / ANSWER-CONTRADICTS / STILL-OPEN / SUPERSEDED with quotes.
   Tag each new finding `LEDGER:<row>` or `NET-NEW`.
3. Peer-status sweep + restatement check (homed decisions §D1–§D7, JD-4/13/14).
4. For 0008 specifically: confirm the iter-2 SPEC-DEFECT is discharged — A9/A4
   re-verified against Final 0002 (outcome lifting, fail-fast load, rename), block 2
   and scenario 5 restated; Refinement Context note gone; no `DEFERRED` rows exist.

Write `pairwise-<A>-<B>.md` here (`Model: claude-fable-5` first line).

## Task B — scoped answer-vs-fences check (per sibling)

For the sibling RDR and each answered entry it tolerates: read the home §D5/§D6/§D7
text and the JD entry; read every ```normative fence, Critical Assumption, MVV/
scenario, and Round-Trip/Invariant clause in the sibling that the answer touches.
Report per (sibling, JD): **CONSISTENT** (frozen text agrees or is silent in a way the
answer merely fills — quote the compatible clause) or **CONTRADICTS** (quote the fenced
clause the answer falsifies; say fenced/unfenced; say whether meaning changes). Also
list the "*Lands in <sibling>*" items the home names and whether each already exists
in the sibling's text (present / absent / contradicted). The gate decides
disposition — you report evidence only.

Write `answer-check-<NNNN>.md` here (`Model: claude-fable-5` first line).

## Task C — whole-set critique

Run the whole-set variant in `2-critique.md` over the eight members + the home.
Read `../iter-2/critique-set.md` first; carry its ledger forward (row ID, RDR, status:
closed/open/superseded/new) so the delta is explicit. Write `critique-set.md` here.

## Return packet (all tasks) — rdr-common §return-packet, EXACTLY:

```
verdict: PASS | BLOCK | INCOMPLETE | NEEDS_DECISION
blocking: yes | no
evidence_paths: [...]
changed_paths: [...]
next_action: ...
summary_50w: ...
```
followed by a compact table (one line per finding / per check). PASS = no
blocks-impl/risks-impl single-RDR finding and no CONTRADICTS on fenced text; BLOCK otherwise.
