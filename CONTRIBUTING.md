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
- **`internal/version/`** — build-identity metadata: `-ldflags` where
  present, Go's embedded VCS stamps as the fallback, so a plain
  `go build` / `go install` binary still reports a source revision.
  Same scheme as retrofit, roborev, and kata.

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
`version_test.go` is the harness pattern for new verb tests. Note it
routes BOTH modes through `respond.OK` — a payload whose whole content is
one line implements `respond.TextLiner` rather than printing directly.

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
make check     # fmt-check + vet + lint + build + graph-lint + docs-check
               #   + test — local mirror of CI
make docs      # regenerate docs/cli-reference.md + llms.txt from the CLI
make vuln      # govulncheck (its own CI job; not in `make check`)
make hooks     # install .githooks (pre-commit: gofmt + go vet;
               #   commit-msg: Conventional Commits; and post-* hooks)
```

CI runs four jobs (see `.github/workflows/ci.yml`): `test` (`make
test-ci`), `lint` (`make fmt-check` + golangci-lint), `vuln`
(govulncheck), and `graph-lint`, which builds the binary and lints
`models/rdr.toml` — the RDR 0006 acceptance gate.

### Reference docs are generated

`docs/cli-reference.md` and `llms.txt` are generated from the command
tree by `make docs`; `make check` fails when they are stale. Never edit
them by hand — change the command's `Long` or its `withExtendedHelp`
body and regenerate. The binary is the source of truth for the flag
grammar, the finding taxonomy, and the refusal codes.
