RDR: 0029-version-promise-on-machine-readable-output | phase: 3d — decision grounding | state: IN-PROGRESS
last: D5 grounded through all four rungs to ask; the predicate-broadening shortcut was falsified on provenance (3c68705 authored the range the ADV tests would be widened to match)
blocker: none yet — D5 at one strong-consult; if it survives, the launch ends stopped:open-decisions
changed: internal/** (phases 2, 3b, 3c), artifacts/{coverage,deviations,verification,status}.md
validate: go test -C /Users/cwensel/sandbox/newcoinc/intrastate/.claude/worktrees/rdr-0029 ./...
baseline: green @2f7d22c
next: record the D5 consult outcome; a survivor means stopped:open-decisions, branch UNMERGED, landing flags held
reads: artifacts/deviations.md, artifacts/verification.md, artifacts/req-list.md, artifacts/coverage.md
session: 3882fd5d-8482-4c1f-a549-d563a4354f57
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
