RDR: 0006-graph-lint-authority-and-guarantees | phase: 3c — fixup | state: COMPLETE
last: Phase 3c fixed FAIL-1/FAIL-2/FAIL-3 (+ADV-1/ADV-2); ADV-3/FAIL-4 pinned as D12 bounded reading; suite green -race, golangci-lint 0 issues
blocker: none
changed: internal/graphlint/, internal/cli/{lint.go,clierr,respond}, models/rdr.toml, .github/workflows/ci.yml, Makefile, artifacts/*.md
validate: go test -race ./... && golangci-lint run ./...
next: unattended triage (/roborev-triage), then land per /rdr-implement-triage --close-and-flight
session: 2026-08-25 | rdr-implement-triage rdr-0006
artifacts: req-list.md coverage.md verification.md deviations.md

COMPLETE — 130 REQ (REQ-1..129 + REQ-MVV) each have >=1 green test; coverage.md has no orphans either direction; REQ-MVV output recorded; deviations D1..D14 classified, D12 open as `needs author decision` recorded per the unattended-run override (implemented on the most defensible reading, miss pinned by test).
