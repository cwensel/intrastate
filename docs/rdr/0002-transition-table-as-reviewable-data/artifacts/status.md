RDR: 0002-transition-table-as-reviewable-data | phase: 3c — fixup | state: COMPLETE
last: Phase 3c fixed FAIL-1/FAIL-2/ADV-1/ADV-2/ADV-3; full suite green, vet+gofmt clean
blocker: none
changed: internal/table/{category,source,load,normalize,model,dump}.go, internal/table/*_test.go, internal/table/testdata/, go.mod, go.sum, docs/rdr/0002-transition-table-as-reviewable-data/artifacts/*.md
validate: go test ./...
next: Stage 8 done — roborev triage, then land (post-mortem at close)
session: 2026-08-25 | rdr-implement-triage rdr-0002
artifacts: req-list.md coverage.md verification.md deviations.md
