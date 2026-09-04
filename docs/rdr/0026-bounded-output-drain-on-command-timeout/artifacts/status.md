RDR: 0026-bounded-output-drain-on-command-timeout | phase: COMPLETION GATE | state: COMPLETE
last: D8/D9/D10 landed — S9 fixture repaired to a 100 ms gap (2x grace), sibling halt documented as a state predicate, both-pipes residue admitted and pinned by a bounded test
blocker: none
changed: internal/cli/cmdbind/cmdbind.go, internal/accessor/{model,executor}.go, internal/cli/{flow_exec,flow_input}.go, docs/cli-output-contract.md, 11 _test.go files, artifacts/*
validate: go test -race ./... (green); gofmt/go vet/golangci-lint clean
next: roborev triage, then land (ship + RDR docs flip + close tracker 936e + flight drain)
session: 2026-09-03
artifacts: req-list.md coverage.md verification.md deviations.md
