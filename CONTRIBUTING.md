# Contributing to intrastate

We welcome bug reports, feature requests, and documentation feedback through
[GitHub Issues](https://github.com/cwensel/intrastate/issues).
**We do not accept pull requests.**

Search existing issues before opening a new one.

## Bug reports

Include:

- A minimal, reproducible test case: the model, required input artifacts,
  and exact commands needed to reproduce the problem.
- Actual and expected results, including **why you believe the current
  behavior is wrong**. Reference documentation or explain the expectation.
- Your `intrastate version` output and operating system.

## Feature requests

Describe the capability's overall value: what it would enable, who
benefits, and why it matters.

Include a short, **non-normative example** to illustrate the capability.
The example explains the intended use; it does not prescribe an interface
or implementation.

## Documentation feedback

Link the relevant section and explain what is incorrect, unclear, or missing.

## Maintainer and local development reference

The following notes cover project conventions and local development.
Deeper topics live under [`docs/`](docs/).

### Project layout

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

### The contract every verb follows

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

#### Exit codes

`clierr.ExitCodeFor` maps an error to a process exit code a script can
branch on:

| exit | meaning                                    |
| ---- | ------------------------------------------ |
| 0    | success / warning                          |
| 1    | unexpected / unclassified error            |
| 2    | user or internal error (bad input, IO)     |
| 3    | environment unavailable                    |
| 130  | interrupted                                |

### Build & test

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

CI runs six jobs (see `.github/workflows/ci.yml`): `test` (`make
test-ci`), `lint` (`make fmt-check` + golangci-lint), `vuln`
(govulncheck), `graph-lint`, which builds the binary and lints
`models/rdr.toml` — the RDR 0006 acceptance gate — `docs` (`make
docs-check`), and `snapshot`, described below.

#### Distribution

Binaries reach users by two paths, both driving the same
`.goreleaser.yaml`. intrastate is pure Go, so one Ubuntu runner
cross-compiles every target; there is no per-OS matrix.

**Tagged releases** are the installable artifact. Pushing a `v*` tag runs
`.github/workflows/release.yml`, which re-runs `make check` against the
tagged commit and then publishes a GitHub Release with a tarball per
platform plus `checksums.txt`:

```sh
git tag -a v0.1.0 -m 'v0.1.0'
git push origin v0.1.0
```

A tag is permanent once published, so the release workflow verifies the
commit itself rather than assuming it passed CI on the way through `main`.

**Snapshots** are per-commit builds for consumers pinning an exact
revision. The `snapshot` job runs on green `main` only and uploads
`intrastate-snapshot-<sha>` as an Actions artifact (7-day retention):

```sh
gh run download <run-id> --repo cwensel/intrastate \
  --name intrastate-snapshot-<sha>
```

Preview the full artifact set locally without publishing:

```sh
make snapshot       # builds dist/ exactly as CI does
make release-check  # asserts the built binary reports its build identity
```

`release-check` guards a silent failure: the `-X` ldflags paths in
`.goreleaser.yaml` name the module path, and if they stop matching
`go.mod` nothing fails to build — the binary just falls back to VCS
stamps and can no longer name its release. Both CI paths run it.

#### Reference docs are generated

`docs/cli-reference.md` and `llms.txt` are generated from the command
tree by `make docs`; `make check` fails when they are stale. Never edit
them by hand — change the command's `Long` or its `withExtendedHelp`
body and regenerate. The binary is the source of truth for the flag
grammar, the finding taxonomy, and the refusal codes.
