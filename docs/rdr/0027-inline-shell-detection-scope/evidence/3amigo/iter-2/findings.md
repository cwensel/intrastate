Model: claude-sonnet-5

# 3amigo delta-scoped review — iter 2 — RDR 0027

## Scope

Checked exactly the seven edited/added passages against the working-tree diff
(`git diff 9c66bf0 -- docs/rdr/0027-inline-shell-detection-scope.md`):
A6 (new), C1's `interpreter set:` line, D-selection-predicate's new worked
example, §phase-3-promise-wording, §testing-strategy (Scenario 3/4 detail
additions, new Scenario 7, framing-line update), MVV (step 3 addition, new
step 5, End-state line), and the `oracle` mini-check's new S7 row.

## Verification performed

- Read `internal/table/load.go::shellInterpreters` (lines 1035-1047) and
  confirmed `sh`/`python` → `{-c}`, `ruby` → `{-e}`, matching C1's new
  membership pointer exactly.
- Confirmed via metadata's Overrides line that C1 widens the predicate from
  "argv0-after-an-env-walk" (today's shipped `interpreterForm`) to "any argv
  position under any prefix" — so the new D-selection-predicate worked example
  (`["python","sh","-c","echo"]` → `python -c`) is evaluated against the
  *proposed* predicate, not today's code. Hand-traced the proposed rule
  (lowest i, then lowest j) against both candidate pairs (i=0 python/-c@2,
  i=1 sh/-c@2): lowest i wins, giving `python -c`. Matches the stated
  outcome and does not contradict the "lowest i, then lowest j" predicate
  language, which was already present (unedited) in C1 before this pass.
- Cross-checked the four sites that gate on A6 (Phase 3 wording, Testing
  Strategy intro, MVV framing, A6's own If-wrong) — all explicitly hedge on
  A6 being unresolved rather than asserting settled fact, satisfying the
  Gate's own "Pending must not have settled-fact prose depending on it" bar.
- Compared Scenario 7 / MVV step 5 / oracle S7 row / A6's If-wrong for a
  common pass bar: description exists, is reachable without a refusal, and
  names both out-of-scope forms non-trivially (not a non-empty check). All
  four state this identically; no drift.
- Verified the "Scenarios 1-3 / 4-5 and 7" framing line's omission of
  Scenario 6 predates this delta (present in the 9c66bf0 baseline before
  the "and 7" append) — not an edit-introduced staleness.
- Verified S4 / MVV step 3 / oracle S4 row's new "argv unchanged" control
  reads identically in all three edited sites.

## Findings

None. The delta is clean.

Considered and ruled out as non-findings:
- Apparent code/RDR mismatch on `interpreterForm` scanning only one `i`
  candidate: resolved — the shipped code implements the pre-C1 narrower
  predicate; C1 is the not-yet-implemented widened predicate the new example
  correctly illustrates.
- Scenario 6 absent from the "Scenarios 1-3 / 4-5 and 7" framing sentence:
  pre-existing gap, not touched by this delta.
- Node (`["node","-"]`) absent from the oracle's "discriminating members"
  list: pre-existing, unedited by this pass.

---
verdict: PASS
blocking: no
evidence_paths: [/Users/cwensel/sandbox/newcoinc/intrastate/docs/rdr/0027-inline-shell-detection-scope.md, /Users/cwensel/sandbox/newcoinc/intrastate/internal/table/load.go]
changed_paths: []
next_action: none
summary_50w: Zero findings. All seven delta sites (A6, C1's shellInterpreters pointer, the two-interpreter tie-break example, Phase 3 create-wording, Scenario 7, MVV step 5, oracle S7 row) are mutually consistent and match source. The worked tie-break example correctly targets C1's proposed (not-yet-shipped) widened predicate.
