RDR: 0011-flow-next-match-conditioned-candidates | phase: 1 — tests first | state: IN-PROGRESS
last: Phase 1 wrote 62 tests in 7 *_0011_test.go files citing 130/130 REQ; red gate PASSED (53 change-behaviour red, 9 preservation green)
blocker: none
changed: internal/cli/flow_{fixtures,next,all,unknown,demand,mvv,rehome}_0011_test.go, artifacts/{req-list,coverage,deviations}.md
validate: make check
next: Phase 2 — implementer (minimum code to green; start at flow_exec.go::invokedReaders)
session: 2026-08-26 | rdr-implement-triage rdr-0011
artifacts: req-list.md coverage.md verification.md deviations.md
