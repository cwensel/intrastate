RDR: 0010-stateless-decision-tables | phase: 1 — tests first | state: IN-PROGRESS
last: 81 tests across table/graphlint/cli; 105/110 REQs covered; RED confirmed (3 pkgs do not compile against main); 10 characterization tests verified green on main in isolation
blocker: none
changed: internal/table/{fixtures,class,emit,emit_shape}_0010_test.go, internal/graphlint/class_0010_test.go, internal/cli/{decision_table,mvv,failure_modes}_0010_test.go, docs/rdr/0010-stateless-decision-tables/artifacts/{coverage,deviations,status}.md
validate: go test ./...
next: Phase 2 — implementation
session: phase-1
artifacts: req-list.md coverage.md verification.md deviations.md
