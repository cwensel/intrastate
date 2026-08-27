RDR: 0011-flow-next-match-conditioned-candidates | phase: 3c — fixup | state: COMPLETE
last: Phase 3c resolved ADV-1/ADV-3 in production (demand term scoped to consultable rows) and re-anchored ADV-2 as TEST-FIXTURE; make check green
blocker: none
changed: internal/cli/{flow_next,flow_exec}.go, 8 *_0011*_test.go, 5 re-homed 0005 assertions, docs/cli-output-contract.md, README.md, artifacts/*.md
validate: make check
next: roborev triage, then land (rdr-implement-triage --close-and-flight)
session: 2026-08-26 | rdr-implement-triage rdr-0011
artifacts: req-list.md coverage.md verification.md deviations.md
