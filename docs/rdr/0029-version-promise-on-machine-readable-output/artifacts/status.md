RDR: 0029-version-promise-on-machine-readable-output | phase: 3c — fixup | state: IN-PROGRESS
last: Phase 3 done — 3a clean (58 REQ + MVV probed against a built binary), 3b BLOCK with ADV-1/2/3; verification.md committed
blocker: none — zero open author decisions
changed: internal/** (phase 2 + 3b tests), artifacts/{coverage,deviations,verification,status}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: Phase 3c clears six reds — ADV-1/2/3 plus the three code-change resolutions (D2 runner re-scope, D3 census sweep, D4 0023 floor); then the completion gate with suite_green, then triage
reads: artifacts/verification.md, artifacts/deviations.md, artifacts/req-list.md, artifacts/coverage.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
