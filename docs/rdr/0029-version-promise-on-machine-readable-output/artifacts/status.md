RDR: 0029-version-promise-on-machine-readable-output | phase: 2 — implementer | state: IN-PROGRESS
last: Phase 2 leg returned NEEDS_DECISION; suite green except 3 tests, all recorded as open deviations
blocker: none — D2/D3/D4 in Phase 3d grounding
changed: internal/cli/**, internal/resolve/**, internal/guard/**, internal/table/**, internal/accessor/**, internal/graphlint/**, docs/cli-output-contract.md, artifacts/{coverage,deviations}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: Phase 3d on D2/D3/D4, then Phase 3a+3b, then 3c if either finds anything
reads: artifacts/deviations.md, artifacts/coverage.md, artifacts/req-list.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
