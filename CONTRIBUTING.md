# Contributing to intrastate

Short, dense, and meant to be re-read on every change. Deeper topics
live under [`docs/`](docs/).

## Project layout

- **`cmd/intrastate/`** — main entry point; a thin shim over
  `internal/cli`.
- **`internal/cli/`** — the Cobra command tree and every verb.
  - **`internal/cli/respond/`** — output gateway. Every verb's
    stdout/stderr routes through this package; honors `--as text|json`.
  - **`internal/cli/clierr/`** — structured `CLIError` type and
    exit-code mapping.
  - **`internal/cli/config/`** — `intrastate.toml` discovery and loader.
- **`internal/version/`** — build-identity metadata set via `-ldflags`.

New domain packages live under `internal/` (or `pkg/` if they become a
public API). Keep `cmd/` a shim.

## The contract every verb follows

1. Register the verb with `cmd.AddCommand(newXxxCmd())` in
   `NewRootCmd` (`internal/cli/root.go`).
2. In `RunE`, call `respond.ValidateMode(cmd)` first.
3. Route success through `respond.OK` and failure through
   `respond.Fail(cmd, &clierr.CLIError{…})`. **Never** print to
   stdout/stderr directly — the gateway owns both streams.
4. Set `SilenceErrors` + `SilenceUsage` on every command so cobra's
   plain-text errors don't stack above the structured envelope.

`internal/cli/version.go` is the worked example; copy it. Its test in
`version_test.go` is the harness pattern for new verb tests.

### Exit codes

`clierr.ExitCodeFor` maps an error to a process exit code a script can
branch on:

| exit | meaning                                    |
| ---- | ------------------------------------------ |
| 0    | success / warning                          |
| 1    | unexpected / unclassified error            |
| 2    | user or internal error (bad input, IO)     |
| 3    | environment unavailable                    |
| 130  | interrupted                                |

## Build & test

```sh
make build     # ./bin/intrastate
make install   # ~/.local/bin/intrastate
make test      # race + coverage
make check     # fmt-check + vet + lint + test — local mirror of CI
make vuln      # govulncheck (its own CI job; not in `make check`)
make hooks     # install .githooks/pre-commit (gofmt + go vet)
```

CI runs `make test-ci`, `make fmt-check` + golangci-lint, and
govulncheck (see `.github/workflows/ci.yml`).
