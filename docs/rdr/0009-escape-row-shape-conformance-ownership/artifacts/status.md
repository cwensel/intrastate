RDR: 0009-escape-row-shape-conformance-ownership | phase: 3c — fixup | state: COMPLETE
last: Phase 3c hardened escapeShapeBreaches (ADV-2) and removed the three text-renderer Count tests (ADV-1/ADV-3 adjudicated invalid); suite green, golangci-lint 0 issues
blocker: none
changed: internal/resolve/resolve.go, internal/table/*, internal/cli/flow_resolve.go, internal/cli/flow_input.go, internal/cli/clierr/clierr.go, docs/cli-output-contract.md, 9 test files, artifacts/*.md
validate: go test ./... && golangci-lint run
next: roborev triage, then land (rdr-implement-triage --close-and-flight)
session: 2026-08-26 | rdr-implement-triage rdr-0009
artifacts: req-list.md coverage.md verification.md deviations.md

COMPLETE — 107 executable REQs each carry >=1 green test (REQ-90/REQ-91 DEFERRED by the
record's Testing-Strategy scope note); coverage.md reports 0 orphans in either direction;
REQ-MVV (REQ-70) output recorded in coverage.md. One deviation remains open as an
author-decision item carried per the unattended override, not a blocker: D3 (SPEC-UNDER)
— clierr.Finding gains `Count int json:"count,omitempty"`, a new public surface on a
record RDRs 0005/0006/0008 co-own. D1/D2 dispositioned; D4 (IMPL-DECISION) closed.
