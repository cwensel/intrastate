RDR: 0029-version-promise-on-machine-readable-output | phase: 3c — fixup | state: IN-PROGRESS
last: 3c fixed 4 of 6 (ADV-3, D2, D3, D4); ADV-1/ADV-2 held:contract, recorded as D5 — a C4-vs-0006 namespace fork
blocker: none — D5 in Phase 3d grounding
changed: internal/table/advisory.go, internal/graphlint/**, internal/cli/**, artifacts/{coverage,deviations,verification,status}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: ground D5; if it survives to ask, the launch ends stopped:open-decisions with the branch unmerged and the landing flags held
reads: artifacts/deviations.md, artifacts/verification.md, artifacts/req-list.md, artifacts/coverage.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
