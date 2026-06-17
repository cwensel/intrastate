# intrastate — agent guide

intrastate is a Go CLI. The binary in `cmd/intrastate` is a thin shim;
all wiring lives in `internal/cli`.

## Layout

- `cmd/intrastate/` — entry point (shim over `internal/cli`).
- `internal/cli/` — Cobra command tree and verbs.
  - `respond/` — output gateway (text/json); every verb's I/O goes here.
  - `clierr/` — structured `CLIError` + exit-code mapping.
  - `config/` — `intrastate.toml` discovery + loader.
- `internal/version/` — build metadata set via `-ldflags`.

## Conventions (follow these; don't re-decide)

- **Never print to stdout/stderr directly.** Route success through
  `respond.OK`, failure through `respond.Fail(cmd, &clierr.CLIError{…})`,
  advisories through `respond.Note` / `respond.Warn`.
- Start every `RunE` with `respond.ValidateMode(cmd)`.
- Set `SilenceErrors` + `SilenceUsage` on every command.
- Every user-facing failure is a `*clierr.CLIError` with a stable
  `Code`, a `Message`, and a `Group` that maps to an exit code. Add new
  codes as needed; keep envelope fields `omitempty`.
- `internal/cli/version.go` + `version_test.go` are the copy-me example
  for a new verb and its test.

## Commands

```sh
make build    # ./bin/intrastate
make test     # race + coverage
make check    # fmt-check + vet + lint + test (mirror of CI)
```

After changing Go code, run `make check` before declaring done.

## Output contract

See [docs/cli-output-contract.md](docs/cli-output-contract.md).
