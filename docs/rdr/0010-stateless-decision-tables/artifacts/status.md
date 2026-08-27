RDR: 0010-stateless-decision-tables | phase: 3c — fixup | state: IN-PROGRESS
last: FAIL-1 and ADV-1..ADV-3 all resolved; suite and make check green
blocker: none
changed: internal/table/load.go, internal/table/class_0010_test.go, internal/graphlint/analysis.go, internal/cli/flow_next_0011_test.go, docs/cli-output-contract.md, docs/model-authoring.md, docs/README.md, docs/rdr/0010-stateless-decision-tables/artifacts/{verification,coverage,deviations,status}.md
validate: go test ./... && make check
next: completion gate
session: phase-3c
artifacts: req-list.md coverage.md verification.md deviations.md
