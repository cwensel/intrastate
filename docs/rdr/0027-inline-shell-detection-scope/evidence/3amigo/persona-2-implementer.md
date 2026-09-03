Model: claude-sonnet-5

# Persona 2 (Implementer) — findings for RDR 0027

Scope note: started from owned elements (0027:C1, 0027:D-naming, 0027:D-selection-predicate) plus the source-anchor edges (`internal/table/load.go::interpreterForm`, `::carrierDefect`, `::filepathBase`). Widened to 0027:§approach, §existing-infrastructure-audit, §decision-rationale, §problem-statement, and the S1-S6/A2-A4 evidence because the single sharpest first-hour question — "does the contract match the code I'm about to edit?" — is only answerable by diffing C1 against the live body of `interpreterForm`, which I read directly at `/Users/cwensel/sandbox/newcoinc/intrastate/internal/table/load.go:1162-1184`.

## Findings (severity-ranked)

### 1. (Informational / confirmed-by-design, not a defect — flag for implementer awareness) Shipped `interpreterForm` does not yet implement C1's predicate; the RDR is explicit that this is the intended delta

**Anchor**: 0027:C1, 0027:§existing-infrastructure-audit, 0027:§approach

**What I checked**: I wrote and ran throwaway probe tests against the *current* `interpreterForm` (`internal/table/load.go:1162`) with `["nice","sh","-c","echo hi"]`, `["timeout","5","sh","-c","echo hi"]`, and `["xargs","sh","-c","echo hi"]` — all three inputs that C1's mini-check `disposition` table and the illustrative-code "Sample refusals" list as required refusals. All three currently return `ok=false` (i.e., lint clean) because the shipped body only walks a leading `env`/`NAME=VALUE` chain and then checks exactly `argv[i]` at the position immediately after that chain — it does not scan the rest of argv for a listed interpreter at any position, contrary to C1's stated predicate ("at ANY later argv position... under any prefix... nice, timeout, xargs, doas").

**Why this isn't actually a finding against the RDR**: §Existing-Infrastructure-Audit and §Approach both say this explicitly and correctly — "Replace the body (same name, same signature): position-free scan; the `env` walk is deleted as subsumed," and the RDR's own Metadata Status is `Draft` (not Implemented). So the gap I found empirically is exactly the gap the RDR exists to close, not an inconsistency in the record.

**Decision it blocks (if raised)**: none — this confirms C1's contract is actionable and its "replace the body, same signature" instruction in §existing-infrastructure-audit is sufficient to start Monday morning. I'm recording it only so the next reader doesn't have to re-derive "is this already live" — verified empirically: it is not, by design.

### 2. Minor — message-attribution ambiguity when two distinct listed interpreters both appear in argv and share a flag spelling

**Anchor**: 0027:D-selection-predicate, 0027:C1 (report: clause)

**Clarification request**: D-selection-predicate says "when several pairs qualify, the lowest interpreter index and then the lowest flag index are reported." For `["python","sh","-c","echo"]`: `python` (i=0) is a listed interpreter whose flag set includes `-c`; `sh` (i=1) is also listed with flag `-c`. Under "lowest i, then lowest j," the predicate reports `python -c` (i=0, first `-c` found at j=2) even though a human reading the same argv would very plausibly read `sh -c` as the operative pair, since `-c` sits immediately after `sh`. C1's predicate text ("argv[j] one of **its** listed flags," tying j to i's own flag set) is self-consistent and this is not a refusal-correctness bug — both interpreters trigger a defect regardless of which is named — but the DIAGNOSTIC MESSAGE names a pair a reviewer may find surprising or misleading when two interpreter words are chained. None of S1-S6, the illustrative "Sample refusals," or A2/A3's spikes exercise a two-distinct-interpreter argv, so there's no worked example to check an implementation against.

**Decision it blocks**: none of the load-bearing behavior (refuse/admit), but it is unresolved guidance for writing the exact tie-break logic beyond the single-`i` case the illustrative code shows, and for what the "detail" text should say if this shape ever appears in a bug report. Low severity because the shape (two chained listed-interpreter words) is exotic and any occurrence still refuses correctly.

### 3. Non-blocking observation — Scope Verification gate section is unfilled template text

**Anchor**: 0027:§scope-verification

Not a C/D-owned element and not something that blocks writing code from C1/D-naming/D-selection-predicate, so ranked lowest / informational only. The section still contains the bracketed instruction placeholder ("[Confirm the Minimum Viable Validation is in scope...]") rather than a written response, which is a Finalization Gate completeness matter, not a contract-clarity one.
