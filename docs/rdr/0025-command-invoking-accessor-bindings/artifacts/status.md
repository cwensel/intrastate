RDR: 0025-command-invoking-accessor-bindings | phase: COMPLETION GATE | state: COMPLETE
last: Phase 3c closed all three converged Phase-3 defects; suite, -race, and make check green
blocker: none
changed: internal/table, internal/accessor, internal/cli/{cmdbind,flowbind}, internal/cli/{flow,flow_exec}.go, docs/cli-reference.md, llms.txt, 9 *_0025_test.go, artifacts/*.md
validate: go test ./... && make check
next: Close — post-mortem (docs/rdr/0025-command-invoking-accessor-bindings-postmortem.md)
session: 2026-08-30T23:52Z
artifacts: req-list.md coverage.md verification.md deviations.md
