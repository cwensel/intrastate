RDR: 0023-resolve-envelope-projection | phase: 3 — self-verification | state: COMPLETE
last: Phase 3a CoVe (no FAIL-N) + Phase 3b adversarial (4 tests added, all pass) committed
blocker: none
changed: internal/cli/{flow_resolve,flow_exec}.go, internal/cli/{flow_projection,flow_partition}.go, internal/cli/flow_{fixtures,projection,partition,mvv,adversarial}_0023_test.go, docs/{cli-output-contract,cli-reference}.md, artifacts/{req-list,coverage,verification,deviations,status}.md, artifacts/mvv-step1-default-golden.json
validate: go test ./... (exit 0, 10 pkgs) ; make check
next: COMPLETE — Stage 8 done; proceed to roborev triage, then Close/post-mortem
session: 2026-08-29T16:00Z | rdr-implement-triage
artifacts: req-list.md coverage.md verification.md deviations.md

COMPLETE — every REQ-N (137/137) has >=1 green test, coverage.md reports no
orphans in either direction, REQ-MVV output is recorded (all six steps pass;
pricing 290->152 B / 47.6%, release-grammar 412->238 B / 42.2%, both matching
the A1 spike and clearing the 40% bar), and deviations.md has no open
needs-author-decision entries (D-1..D-7 are resolved, mechanical translation,
or recorded).
