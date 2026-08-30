RDR: 0025-command-invoking-accessor-bindings | phase: 1 — Tests first | state: IN-PROGRESS
last: wrote 79 REQ-cited tests across 7 files; red gate confirmed (79 FAIL, 0 regressions)
blocker: none
changed: internal/table/{model.go,category.go,command_carrier_0025_test.go}, internal/accessor/{model.go,exec_detail_0025_test.go}, internal/cli/cmdbind/*, internal/cli/flowbind/{registry.go,registry_command_0025_test.go}, internal/cli/{flow_exec.go,command_gate_0025_test.go,command_fixtures_0025_test.go,command_mvv_0025_test.go}, artifacts/{coverage.md,deviations.md,status.md}
validate: go test ./...
next: Phase 2 — Implementation
session: 2026-08-30T00:00:00Z
artifacts: req-list.md coverage.md verification.md deviations.md
