# CLI output contract

Every verb routes stdout and stderr through `internal/cli/respond`. This
document is the authoritative description of what that gateway emits.

## Modes

The persistent root flag `--as text|json` selects the wire format.
Unknown values are refused (`flag-invalid-value`, exit 2).

## `--as=text`

- **stdout** — verb-defined human output.
- **stderr** — advisory and error lines:
  - `note: <message>`
  - `warning: <message>` (with optional `  detail:` / `  hint:` lines)
  - `error: <code>: <message>` (with optional `  detail:` / `  hint:`)

## `--as=json`

- **stdout** — exactly one terminal envelope, discriminated by `type`:
  - `{"type":"ok", ...}` on success
  - `{"type":"failed", ...}` is not emitted by `Fail`; instead the
    structured `CLIError` envelope (`{"code":...,"message":...}`) is
    written to stdout so a single stream carries both dispositions.
    _(Revisit this if/when an intermediate-record `Stream` emitter is
    added — at that point the terminal record should carry an explicit
    `type` discriminator.)_
- **stderr** — advisories only, discriminated by `level`:
  - `{"level":"note","message":...}`
  - `{"level":"warning","message":...,"code":...}`

The terminal record is emitted on every graceful exit. Its absence means
the process was killed.

## Exit codes

See [CONTRIBUTING.md](../CONTRIBUTING.md#exit-codes). The mapping lives
in `clierr.ExitCodeFor`.
