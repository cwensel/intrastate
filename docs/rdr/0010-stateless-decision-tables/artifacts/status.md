RDR: 0010-stateless-decision-tables | phase: completion gate | state: COMPLETE
last: Phase 3c fixed FAIL-1 (docs) + ADV-1/ADV-2 (machine-only invariants) + ADV-3 (class-check order); suite + make check green
blocker: none
changed: internal/table/{model,source,load,normalize,dump}.go, internal/graphlint/{reach,analysis,coverage,taxonomy}.go, internal/cli/flow_resolve.go, docs/{cli-output-contract,model-authoring,README}.md, 100 internal/table/testdata fixtures, 0010 test files
validate: go test ./... && make check
next: roborev triage, then land
session: completion-gate
artifacts: req-list.md coverage.md verification.md deviations.md
