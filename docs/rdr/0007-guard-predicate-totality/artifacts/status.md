RDR: 0007-guard-predicate-totality | phase: 3c — fixup | state: COMPLETE
last: Phase 3c closed FAIL-1/2 and ADV-1/2/3; suite green, race-clean, lint 0 issues
blocker: none
changed: internal/resolve/guard.go, internal/resolve/guardcontract.go, internal/resolve/resolve.go, internal/resolve/guard_*_test.go
validate: go build ./... && go test ./... && go test -race ./internal/resolve/ && golangci-lint run
next: roborev triage (batch:rdr-0007), then land per BUILD-ORDER (0002 is run 2 of 8)
session: 2026-08-25 | rdr-implement-triage
artifacts: req-list.md coverage.md verification.md deviations.md
