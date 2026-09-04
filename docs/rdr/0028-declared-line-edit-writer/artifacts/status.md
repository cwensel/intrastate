RDR: 0028-declared-line-edit-writer | phase: 2 — Implementation | state: GREEN
last: Phase 2 implementer landed C1.1-C1.6 over four commits; suite and lint green; MVV recorded (8/8 at the accessor seam)
blocker: none
open: D9 (DEPENDENCY-LIMIT) — RDR 0002's one-reader/one-writer-per-owned-tag arity fences out the record's Illustrative Code shape for a shared artifact; Phase 2 took the two-key resolution and all eight MVV steps pass under it
changed: internal/table/{model,source,load,category,dump,edit}.go internal/table/testdata/neg/neg-edit-*.toml internal/accessor/{model,executor}.go internal/cli/flowbind/{edit,registry}.go internal/cli/cmdbind/cmdbind.go internal/cli/{flow_exec,flow_input}.go internal/cli/edit_envelope_0028_test.go docs/cli-output-contract.md
validate: go test ./... && golangci-lint run
next: Phase 3 — roborev triage over the four Phase 2 commits
session: 2026-09-04 | rdr-implement-triage unattended
artifacts: req-list.md coverage.md deviations.md status.md
