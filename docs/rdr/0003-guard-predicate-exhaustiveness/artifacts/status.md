RDR: 0003-guard-predicate-exhaustiveness | phase: 3c — fixup | state: INCOMPLETE
last: Phase 3c closed FAIL-1 and ADV-1/2/3 plus two latent defects; 149/150 guard tests green, race-clean, lint 0 issues
blocker: D15 open (needs author decision) — ADV-1's second assertion is stale against its own premise, leaving 1 red test; COMPLETE requires no open author-decision entries and every REQ green
changed: internal/guard/*.go, internal/table/{model,normalize}.go, internal/guard/guard_*_0003_test.go, internal/table/fixup_0002_test.go
validate: go build ./... && go test ./... && go test -race ./internal/guard/ && golangci-lint run
next: author decides D15 (retire the stale assertion, or reject D14's REQ-57 reading and direct a different ADV-1 resolution), then re-run /rdr-implement 3
session: 2026-08-25 | rdr-implement-triage
artifacts: req-list.md coverage.md verification.md deviations.md
