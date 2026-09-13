RDR: 0029-version-promise-on-machine-readable-output | phase: 1 — tests first | state: IN-PROGRESS
last: Phase 1 red suite committed; five unassertable REQ demoted to EXCLUDED in both ledgers
blocker: none
changed: internal/cli/schema_*_0029_test.go, internal/graphlint/tier_assertions_0029_test.go, artifacts/{req-list,impact,coverage}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: Phase 2 — implementer legs over impact.md families in file order, then the one full suite
reads: docs/rdr/0029-version-promise-on-machine-readable-output/artifacts/req-list.md, .../coverage.md, .../impact.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
