RDR: 0008-recognized-tag-key-ownership | phase: 3c — fixup | state: COMPLETE
last: Phase 3c fixed ADV-3 (sorted-key scan), adjudicated ADV-1/ADV-2 to RDR 0005/0006 as D9, closed D3 citation half; suite green, vet + golangci-lint clean
blocker: none
changed: internal/resolve/precondition.go, internal/resolve/resolve.go, internal/table/load.go, internal/table/advisory.go, internal/cli/*_0008_test.go, internal/{resolve,table}/*_0008_test.go, docs/rdr/0008-recognized-tag-key-ownership.md, <art>/{req-list,coverage,verification,deviations}.md
validate: go test ./... && go vet ./... && golangci-lint run
next: Close — write post-mortem (docs/rdr/0008-recognized-tag-key-ownership-postmortem.md)
session: 2026-08-26 | rdr-implement-triage rdr-0008
artifacts: req-list.md coverage.md verification.md deviations.md
COMPLETE — every REQ-N has >=1 green test, coverage.md has no open orphans (D3's record-citation orphan closed this run), REQ-MVV output recorded, deviations.md has no open needs-author-decision entries.
