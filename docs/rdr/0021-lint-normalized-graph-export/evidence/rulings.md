# Author rulings — 0021-lint-normalized-graph-export

## 2026-09-12 — rdr-draft-to-lock (orchestrator)

- **fork: refine advance on an inherited Ledger PASS** — RULED: Accept the prior run's committed Ledger PASS for Stage 3 refine (run-plan.md, commit b9d8c3b; refine commit aed9efe) as satisfying the draft-to-lock advance carve-out, and enter this run at Stage 4. The carve-out is scoped to "this run's Ledger"; the author authorized treating the prior run's committed PASS as equivalent. The router's repeated `/rdr-refine` (rule `locate-draft-refine`, guard `ca=all-pending`) is the known no-done-signal property of refine, not an outstanding refine debt. — absorbed @resolve 2026-09-12
