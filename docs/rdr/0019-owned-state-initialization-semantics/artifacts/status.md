RDR: 0019-owned-state-initialization-semantics | phase: 3 — self-verification | state: COMPLETE
last: Phase 3c fixed FAIL-1 (doc cardinal) and FAIL-2/ADV-2 (role-scoped absent report); gate returned COMPLETE
blocker: none
changed: internal/cli/flow_initstate.go, internal/cli/flowbind/flowbind.go, internal/cli/flow.go, internal/cli/flow_exec.go, docs/cli-reference.md, docs/cli-output-contract.md, docs/model-authoring.md, llms.txt, artifacts/*
validate: go test ./... && golangci-lint run ./...
baseline: green @fcb85d6
next: roborev triage, then land (ship + RDR docs flip + tracker close + flight)
session: 2026-09-05 | rdr-implement-triage
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
