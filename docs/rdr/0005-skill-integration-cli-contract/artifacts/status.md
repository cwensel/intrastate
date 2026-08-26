RDR: 0005-skill-integration-cli-contract | phase: 3c — fixup | state: COMPLETE
last: Phase 3c fixed all 7 defects (FAIL-1, ADV-1/2/3, DEV-2/3/5) as separate green increments; full suite green
blocker: none
changed: internal/cli/{flow,flow_input,flow_exec,flow_next,flow_resolve,flow_state}.go, internal/cli/flowbind/, internal/cli/respond/, internal/cli/clierr/, internal/table/model.go, docs/cli-output-contract.md, artifacts/*.md
validate: go test ./... && golangci-lint run
next: roborev triage, then land (rdr-implement-triage --close-and-flight)
session: 2026-08-26 | rdr-implement-triage rdr-0005
artifacts: req-list.md coverage.md verification.md deviations.md
