# intrastate

intrastate is a Go CLI for making workflow state transitions explicit,
reviewable, and deterministic.

The project is built around a simple idea: a flow should be navigated from
declared state and recognized outcomes, not reimplemented ad hoc in every skill,
script, or agent. intrastate will provide a small resolver kernel, a reviewable
transition-table format, static graph lint, and a CLI surface that lets callers
ask what can happen next, resolve a recognized outcome, and safely read or
persist owned state.

## What It Does

intrastate models workflows as tag-based transition graphs. Flow authors define
legal outcomes, guards, state writes, and accessor bindings in data. The tool
then normalizes and lints that model before runtime, so incomplete, ambiguous,
or illegal graphs are caught during review and CI.

At runtime, intrastate refuses unsafe guesses. Given the same transition table,
state snapshot, observed tags, and recognized outcome, the resolver returns the
same result: one legal transition plan or one typed refusal.

## Install

```sh
make install        # builds ./bin/intrastate and installs to ~/.local/bin
```

Or build locally:

```sh
make build          # ./bin/intrastate
```

## Usage

```sh
intrastate version                       # build version, commit, date
intrastate version --as=json             # same, as a JSON envelope

intrastate lint --model flow.toml        # check a model's graph invariants

intrastate flow next       --model flow.toml --artifact state=state.json
intrastate flow resolve    --model flow.toml --artifact state=state.json \
    --outcome approved
intrastate flow read-state --model flow.toml --artifact state=state.json
intrastate flow set-state  --model flow.toml --artifact state=state.json \
    --write status=approved
```

Every model-taking command requires exactly one of `--model <path>` or
`--flow <id>`. `--flow` is reserved for config discovery; this build
registers no ids and refuses with `flow-model-not-found`, so pass
`--model`.

Every command accepts the global `--as text|json` flag. Under `--as=json`
stdout carries a single terminal envelope discriminated by a `type`
field (`ok` | `failed`); under `--as=text` it is human-readable.

The CLI is self-describing. `--help` on any command is the terse
orientation; `--help-all` adds its extended reference — the vocabulary,
the refusal codes it can return, and the exit contract. A single
`intrastate --help-all` (or `intrastate help --all`) prints the whole
surface at once.

```sh
intrastate --help-all                    # every command's full reference
intrastate flow resolve --help-all       # one verb's refusal vocabulary
intrastate lint --help-all               # the live finding taxonomy
```

The finding and refusal code lists in that output are generated from the
same constants the wire is emitted from, so they cannot drift from what a
caller actually receives.

## Design Goals

- Keep transition logic in reviewable data, not scattered code.
- Make resolver behavior deterministic and replay-safe.
- Use symbolic guards so lint can prove coverage and overlap where domains are
  finite.
- Treat graph lint as the blocking design-time authority.
- Keep artifact access explicit through declared read, gate, and write
  accessors.
- Preserve a thin CLI over the kernel rather than building a workflow
  orchestrator.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for layout, build/test commands,
and the output/error contract every verb follows. The locked architecture
decisions live under [docs/rdr](docs/rdr).

```sh
make check          # fmt-check + vet + lint + test (local mirror of CI)
make test           # race + coverage
```

## License

[MIT](LICENSE)
