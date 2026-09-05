RDR: 0020-undeclared-tag-key-admission | phase: 3 — self-verification | state: COMPLETE
last: completion gate COMPLETE — 0 orphans, REQ-MVV recorded, 0 open decisions, suite green
blocker: none
changed: internal/cli/flow_input.go, docs/model-authoring.md, internal/cli/flow_carrier*0020_test.go, artifacts/*
validate: go test ./...
baseline: green @d97f2c5
next: Phase 2 — roborev triage, then land
session: 2026-09-04
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
