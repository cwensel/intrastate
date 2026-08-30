RDR: 0025-command-invoking-accessor-bindings | phase: 2 — Implementation | state: IN-PROGRESS
last: full suite green tree-wide; golangci-lint clean; REQ-MVV output recorded
blocker: none
changed: internal/table/{source.go,model.go,category.go,load.go,accessors_test.go,dump_test.go,testdata/neg/neg-command-*.toml}, internal/accessor/{executor.go,model.go}, internal/cli/cmdbind/{cmdbind.go,seam_0025_test.go}, internal/cli/flowbind/registry.go, internal/cli/{flow.go,flow_exec.go,command_fixtures_0025_test.go,command_gate_0025_test.go,command_mvv_0025_test.go}, docs/cli-reference.md, llms.txt, artifacts/{coverage.md,deviations.md,status.md}
validate: go test ./...
next: Phase 3 — Self-verification
session: 2026-08-30T23:11:16Z
artifacts: req-list.md coverage.md verification.md deviations.md
