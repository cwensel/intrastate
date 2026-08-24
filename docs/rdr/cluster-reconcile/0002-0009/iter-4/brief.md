# Iteration-4 brief — cluster 0002-0009 (Stage 7.1 re-run, PAST THE DEMOTION CAP)

Model: claude-fable-5
Date: 2026-08-24

Read-only against the RDRs. Sub-agents write ONLY their own output file under
this directory. Never edit an RDR, the JDR, the README, or any deviations.md.

N = 4. The cap (2 re-runs) is spent: this iteration MAY NOT demote. Every
finding is evidence for one of: RECONCILED (WITH TOLERANCES) or
`stopped:cluster-flapping`. Report findings honestly — the gate, not the
sub-agent, decides which.

## Paths (absolute)

- Records dir: `/Users/cwensel/sandbox/newcoinc/intrastate/docs/rdr/`
  - 0002 `0002-transition-table-as-reviewable-data.md` (re-locked 2026-08-24, `365bd66`; STAGE-SCOPED re-entry from iter-3)
  - 0003 `0003-guard-predicate-exhaustiveness.md` (re-locked 2026-08-24, `31bb9e4`; RE-LOCK-ONLY re-entry)
  - 0004 `0004-accessor-execution-safety-model.md` (re-locked 2026-08-24, `ffaf897`; RE-LOCK-ONLY re-entry)
  - 0005 `0005-skill-integration-cli-contract.md` (`76a9e67`, unmoved since 2026-08-12)
  - 0006 `0006-graph-lint-authority-and-guarantees.md` (`567600f`, iter-3 status edit only)
  - 0007 `0007-guard-predicate-totality.md` (`526c481`, unmoved since iter-3)
  - 0008 `0008-recognized-tag-key-ownership.md` (`567600f`, iter-3 status edit only)
  - 0009 `0009-escape-row-shape-conformance-ownership.md` (`567600f`, iter-3 status edit only)
  - (ignore `*-postmortem.md`; 0001 is Implemented — out of cluster)
- **Home**: `/Users/cwensel/sandbox/newcoinc/intrastate/docs/jdr/0001-resolve-kernel-seam.md`
  (`567600f` — moved since iter-3's recorded `4f9639c` ONLY by iter-3's own
  edits: JD-5 sharpened, JD-8 third carrier, JD-9 classification arm still
  open, JD-15/16/17 sibling-check notes, D7(ii) open note. **No new answer.**
  JD-5, JD-8, JD-9, JD-18 remain unanswered.)
- **Artifact of record for peer status**: `docs/rdr/README.md` Index table.
- **Ledgers**: iter-1 `../report.md`; iter-2 `../iter-2/reconcile-report.md`;
  iter-3 `../iter-3/reconcile-report.md` (its findings table, tolerance table,
  DEFERRED rows and the three SPEC-DEFECT dispositions are what this iteration
  traces to). Prior per-pair scans: `../pairwise-*.md`, `../iter-2/pairwise-*.md`,
  `../iter-3/pairwise-*.md`; iter-3 answer checks `../iter-3/answer-check-*.md`.
- Deviations already written (do not re-defer these): `0006-*/artifacts/deviations.md`
  D1, D2; `0008-*/artifacts/deviations.md` D1, D2.
- Pairwise prompt: `/Users/cwensel/sandbox/newcoinc/intrastate/.claude/skills/rdr-cluster-reconcile/pairwise.md`
- Whole-set critique prompt: `/Users/cwensel/sandbox/newcoinc/intrastate/.claude/skills/rdr-cluster-reconcile/2-critique.md`

## Delta scope

Moved since iter-3: **0002, 0003, 0004** (full re-entry + re-lock). Home: no
answer. So:
1. **DISCHARGE CHECK** (Task B) for 0002, 0003, 0004 — did the re-lock land
   what the iter-3 SPEC-DEFECT named, against the home's §D5/§D6/§D7?
2. **RE-SCAN** (Task A) every pair with a moved member (13):
   0002-0003, 0002-0004, 0002-0005, 0002-0006, 0007-0002, 0008-0002, 0009-0002,
   0003-0005, 0003-0006, 0007-0003, 0004-0005, 0007-0004, 0009-0004.
3. **RE-SCAN** whole-set critique (Task C).
CARRIED from iter-3/iter-2 (no member moved, home unanswered): 0005-0006,
0007-0005, 0007-0006, 0007-0008, 0007-0009, 0008-0005-0006, 0008-0009, 0009-0005.

## Task A — pairwise re-scan (per assigned pair)

1. Run `pairwise.md` verbatim against the two CURRENT files. Quote-anchored
   (`NNNN:line` + a grepable snippet).
2. Ledger trace: for every iter-3 row for this pair (and iter-2/iter-1 rows it
   carries or that were CARRIED into iter-3), state RESOLVED / STILL-OPEN /
   SUPERSEDED / REGRESSED with quotes. Tag each new finding `LEDGER:<row>` or
   `NET-NEW`.
3. For each finding say: fenced (inside a ```normative fence) or unfenced;
   does it change a clause's MEANING; is there a mechanical check (grep/compile/
   test the RDRs already specify) that decides it; is it internal to one RDR or
   a decision neither owns (and if joint, which open JD it belongs to, or none).
4. Peer-status sweep + restatement check: stale claims about a peer's status or
   text (e.g. 0007 cites "0002's Refinement Context Direction" at 0007:687,
   939, 951, 2040 — 0002 cut that section at re-entry; report whether the duties
   it names landed in re-locked 0002), and any sibling restating a homed
   decision's (§D1–§D7, JD-*) mechanism prose instead of citing it.
5. Round-Trip / Inverse invariants declared across the seam: assert value-for-
   value, not signature agreement.

Write `pairwise-<A>-<B>.md` here (`Model: claude-fable-5` first line).

## Task B — discharge check (per re-locked member: 0002, 0003, 0004)

Read `../iter-3/reconcile-report.md` §Dispositions for this RDR and
`../iter-3/answer-check-<NNNN>.md`. For each contradicting clause named there,
find its successor in the CURRENT text and report LANDED (quote) / NOT LANDED /
LANDED DIFFERENTLY (quote both; say whether the difference contradicts the home's
§D5/§D6/§D7 text or an open JD). Also: Refinement Context section absent; Status
line and README row agree; for 0002 the re-verify set A1, A9, A12 carries a
fresh `Verified` stamp with evidence dated on/after 2026-08-24; for 0002 the
D7(ii) provenance-scope note (home JD-17 "Open at (ii)", owner 0002) — say
whether 0002 settled it and how (quote), and whether the home text at D7(ii)
now contradicts 0002's settlement. For 0004: Prerequisites checklist state.
Report evidence only; the gate dispositions.

Write `discharge-check-<NNNN>.md` here (`Model: claude-fable-5` first line).

## Task C — whole-set critique

Run the whole-set variant in `2-critique.md` over the eight members + the home.
Read `../iter-3/critique-set.md` first; carry its ledger forward (row ID, RDR,
status: closed/open/superseded/new) so the delta is explicit. Write
`critique-set.md` here.

## Return packet (all tasks) — rdr-common §return-packet, EXACTLY:

```
verdict: PASS | BLOCK | INCOMPLETE | NEEDS_DECISION
blocking: yes | no
evidence_paths: [...]
changed_paths: [...]
next_action: ...
summary_50w: ...
```
followed by a compact table (one line per finding / per check): pair | finding |
fenced? | meaning-change? | mechanical check? | ownership/JD | ledger tag |
severity (blocks-impl / risks-impl / cosmetic). PASS = no blocks-impl or
risks-impl finding on fenced text or on meaning; BLOCK otherwise. Do NOT
propose dispositions.
