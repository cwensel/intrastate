# RDR 0014 — Stage 2 prior-art record

Problem class: verification oracle for repository CI gates — how a test
proves a gate *enforces*, not merely that its text looks right.

## Accepted citations

### In-repo priors (load-bearing)

- `Makefile` header (lines 1–2): "CI invokes the same targets so local
  and remote runs share one source of truth." ⇒ the repo's own stated
  convention is single-executable-carrier; the bespoke inline shell in
  ci.yml's `graph-lint` job violates it.
- `Makefile::release-check` — behavioral gate oracle precedent: runs the
  built binary and asserts the commit-stamp width (its comment records
  that a text-presence check was tried and demonstrated unsound).
- `Makefile::docs-check` — behavioral oracle precedent:
  regenerate-then-diff, never grep.
- `internal/cli/lint_gate_0006_test.go::TestReq120_ADeliberatelyIllegalEditToTheCheckedInModelFailsTheGate`
  — existing in-process behavioral adversarial (invokes `runCmd`, not
  text).
- `internal/cli/clierr/clierr.go::ExitCodeFor` — `GroupUserEnv` and
  `GroupInternal` both map to exit 2 (the RDR 0006:A5 rationale for
  reading the JSON `code`, which any restated oracle must preserve as
  diagnosis).
- `internal/graphlint/taxonomy.go::AggregateCode` = `graph-lint-failed`
  — the refusal code a blocking run returns; the behavioral predicate's
  expected diagnosis.
- RDR `0006:A5` (Gate target decision): "The gate asserts on the JSON
  `code` field, never the exit integer alone, because `ExitCodeFor`
  maps `GroupUserEnv` and `GroupInternal` both to 2."

### Peer-instance reads (langref checkout set; instance question)

- `beads` — `scripts/ci_capability_selector_test.go` (line ~229):
  `exec.Command("bash", filepath.Join("..", ".github", "scripts",
  "ci-capability-selector.sh"))` — a Go test subprocess-invoking the
  exact script CI runs, asserting on its behavior. Direct precedent for
  the single-carrier + behavioral-oracle arm.
- `helm` — `.github/workflows/build-test.yml`: `run: make
  test-source-headers` / `make test-coverage` / `make build` — CI as a
  thin shim over Make targets.
- `roborev` — `.github/workflows/ci.yml`: `run: make markdown-ci` —
  same pattern.
- Negative instance result: a bounded literal sweep of the peer set
  (gh-cli, helm, goreleaser, roborev, beads) found **no** Go test
  asserting on the repo's own CI workflow text; the `.github/workflows`
  hits in gh-cli are its `gh workflow` product feature, not
  self-assertions.

### Literature (class level, supporting — not load-bearing)

- Continuous Delivery (Humble & Farley), DevRef corpus, p.187: the one
  deployment script "used throughout the pipeline … and thus came to
  trust it" — one carrier across environments is what makes a gate
  trustable.
- The DevOps 2.0 Toolkit, DevRef corpus, p.226: moving pipeline logic
  out of CI-server properties into a script in the code repository.

## Queries run (arc, corpus DevRef; 3 total — within budget)

1. "keep CI pipeline configuration thin: run the same scripts locally
   and in continuous integration so the build is testable"
2. "run the same scripts in CI as developers run locally; pipeline
   config should only invoke scripts" → DevOps 2.0 p.226, CD p.187
3. "build and test scripts checked into version control, same script
   run by CI server and by developers on their workstations" → CD
   p.90/101 (context only)

Plus literal (non-corpus) sweeps over the langref checkout set:
workflow files for `make |script/` invocation; `*_test.go` for
`.github/workflows` self-assertions.

## Rejected branches

- Workflow-level e2e (nektos/act, canary branch): no precedent found in
  the peer set; carried into the record as an alternative, rejected on
  cost/flakiness/dependency.
- Continuous Delivery "pretested commit" (p.101): about commit
  workflow, not gate oracles — not used.
- Claims that cannot be quote-confirmed (order/emission): (a) the
  bespoke ci.yml step's crash-path behavior (fail-open), (b) `make
  graph-lint MODEL=x` aborting before the example-model loop — demoted
  to Resolve assumptions (Method: Spike), per stage rule.
