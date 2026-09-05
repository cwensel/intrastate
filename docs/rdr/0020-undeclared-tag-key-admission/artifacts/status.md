RDR: 0020-undeclared-tag-key-admission | phase: 2 — implementation | state: IN-PROGRESS
last: Phase 1 PASS — 28 tests, red gate confirmed (15 new REQs red), no orphans @886e600
blocker: none
changed: internal/cli/flow_carrier*0020_test.go, artifacts/coverage.md
validate: go test ./...
baseline: green @d97f2c5
next: Phase 2 leg 1 — Phase 1 tests, then TestAdv0008, TestMVV0008, then full suite
session: 2026-09-04
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
