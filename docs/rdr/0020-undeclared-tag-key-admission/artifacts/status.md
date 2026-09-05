RDR: 0020-undeclared-tag-key-admission | phase: 3 — self-verification | state: INCOMPLETE
last: Phase 3a PASS (no FAIL-N); Phase 3b ADV-3 -> D3; suite green, 0 orphans, REQ-MVV recorded
blocker: stopped:open-decisions — D1 (REQ-40/41 names absent docs/model-schema.md) and D3 (Failure Modes Diagnosis echo absent from refusing exit) are open author decisions
changed: internal/cli/flow_input.go, docs/model-authoring.md, internal/cli/flow_carrier*0020_test.go, artifacts/*
validate: go test ./...
baseline: green @d97f2c5
next: author resolves D1 + D3 in deviations.md (rewrite Status line to RESOLVED), then re-run the completion gate
session: 2026-09-04
artifacts: req-list.md impact.md coverage.md verification.md deviations.md
