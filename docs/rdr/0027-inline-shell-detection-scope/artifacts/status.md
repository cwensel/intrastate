RDR: 0027-inline-shell-detection-scope | phase: 3 — self-verification | state: COMPLETE
last: Phase 3c closed 3 coverage orphans with mutation-verified guards; REQ-32 withdrawn (D8)
blocker: none
changed: internal/table/load.go, internal/table/category.go, internal/cli/lint.go, docs/cli-reference.md, 4 new *_test.go, artifacts/*
validate: go test ./... && golangci-lint run
next: roborev triage, then land (rebase, ff-merge, Status→Implemented)
session: 2026-09-04T03:55:28Z
artifacts: req-list.md coverage.md verification.md deviations.md
