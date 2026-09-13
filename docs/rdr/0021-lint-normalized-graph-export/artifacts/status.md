RDR: 0021-lint-normalized-graph-export | phase: 3c — fixup | state: INCOMPLETE
last: ADV-1/ADV-2 fixed (fc1cbbb); FAIL-20 held:contract, recorded as D6
blocker: open author decisions D2, D3, D4 (TEST-FIXTURE), D6 (SPEC-UNDER)
changed: internal/graphlint/export.go, artifacts/verification.md, artifacts/deviations.md
validate: go test ./...
baseline: red @fc1cbbb       # D2/D3/D4 test-fixture reds, unchanged from f61a33f
next: run 3d on D6; the gate then reports stopped:open-decisions
reads: internal/graphlint/export.go:88-114, internal/graphlint/reach.go:200-223, internal/cli/graph_document.go:168-191, internal/guard/declaration.go:80-150, internal/guard/declaration.go:206-231, internal/table/load.go:868-900, internal/cli/testdata/graph_0021_state_machine.json
session: 2026-09-13 | phase-3c-fixup
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
