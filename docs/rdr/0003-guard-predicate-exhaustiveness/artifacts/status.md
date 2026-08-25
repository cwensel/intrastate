RDR: 0003-guard-predicate-exhaustiveness | phase: 3c — fixup | state: COMPLETE
last: D15 resolved (stale ADV-1 assertion retired per 0007:C1); 150/150 guard tests green, race-clean, lint 0 issues
blocker: none
changed: internal/guard/*.go, internal/table/{model,normalize}.go, internal/guard/guard_*_0003_test.go, internal/table/fixup_0002_test.go
validate: go build ./... && go test ./... && go test -race ./internal/guard/ && golangci-lint run
next: roborev triage (batch:rdr-0003), then land per BUILD-ORDER (0003 is run 3 of 8)
session: 2026-08-25 | rdr-implement-triage
artifacts: req-list.md coverage.md verification.md deviations.md
