RDR: 0029-version-promise-on-machine-readable-output | phase: 3c — fixup (leg 2) | state: IN-PROGRESS
last: D5 resolved by strong-consult — emit site predates 0029 (RDR 0008 on main), so reading 1 governs and the ADV predicates are over-narrow; zero open author decisions
blocker: none
changed: internal/** (phases 2, 3b, 3c), artifacts/{coverage,deviations,verification,status}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: 3c leg 2 broadens the two ADV predicates to every declared seam; then the completion gate with suite_green=true, then triage
reads: artifacts/deviations.md, artifacts/verification.md, artifacts/req-list.md, artifacts/coverage.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
