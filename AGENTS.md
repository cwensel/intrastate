# intrastate — agent guide

intrastate is a Go CLI. The binary in `cmd/intrastate` is a thin shim;
all wiring lives in `internal/cli`.

## Layout

- `cmd/intrastate/` — entry point (shim over `internal/cli`).
- `internal/cli/` — Cobra command tree and verbs.
  - `respond/` — output gateway (text/json); every verb's I/O goes here.
  - `clierr/` — structured `CLIError` + exit-code mapping.
- `internal/version/` — build metadata: `-ldflags`, falling back to Go's
  embedded VCS stamps so an unstamped build still names its revision.

## Conventions (follow these; don't re-decide)

- When seeding RDRs, do not include source/provenance file paths or paths
  outside this repo. Keep seed context portable; current-repo paths are fine
  when they are implementation-relevant.
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
make check    # fmt-check + vet + lint + build + graph-lint
              #   + docs-check + test (mirror of CI)
make docs     # regenerate docs/cli-reference.md + llms.txt
```

After changing Go code, run `make check` before declaring done.

## Docs

`docs/cli-reference.md` and `llms.txt` are GENERATED from the command
tree (`make docs`, gated by `make check`). Never hand-edit them: change
the command's `Long` or its `withExtendedHelp` body instead. The help
bodies spell finding and refusal codes through the same constants the
wire uses, so the published vocabulary cannot drift from the emitted one.

## Output contract

See [docs/cli-output-contract.md](docs/cli-output-contract.md) for worked
payloads and rationale; run `intrastate --help-all` for the live
reference.
