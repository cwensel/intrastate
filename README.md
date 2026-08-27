# intrastate

intrastate is a Go CLI for making workflow state transitions explicit,
reviewable, and deterministic.

The project is built around a simple idea: a flow should be navigated from
declared state and recognized outcomes, not reimplemented ad hoc in every skill,
script, or agent. intrastate provides a small resolver kernel, a reviewable
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

## Two Model Classes

A model declares its class in `[model]`, and the class says what kind of
question the model answers.

**State machine** (`class = "state-machine"`, the default) — the model owns
state and advances it. It declares owned tags, a root `[initial]` state,
accessors that read and write that state, and terminal predicates that say
where it stops. Its rules carry writes, and `flow resolve` answers with the
plan for the next state. This is the classical transition system: a
tag-valued state, a recognized-outcome alphabet, and a transition relation
authored as rows.

**Decision table** (`class = "decision-table"`) — the model owns no state
and advances nothing. It maps a supplied situation to an answer. It
declares zero owned tags, no root, no terminals, and no accessors; its
rules carry `[rule.emit]` instead of writes, and `flow resolve` answers
with that row's emit block. This is the classical decision table: condition
dimensions spanning a product, one row per cell, one answer per cell.

Both ship, both are authored in one grammar, and both are checked by the
same analysis. What they share is the part worth stating plainly:

- **Coverage is proved, per outcome, for both.** Rules are grouped by the
  outcome they match; the group's guard atoms span a product of their
  declared domains; every assignment must be claimed. An unclaimed
  assignment is `graph-coverage-gap` and the model does not lint. This is
  the decision table's completeness check, and a state machine gets it too.
- **Exactly one row, or a refusal.** Two rows enabled by the same
  assignment is `graph-overlap` at lint and `flow-ambiguous-match` at
  runtime. Nothing resolves the tie by row order, priority, or authoring
  position — picking one would be a decision the model never made.
- **The escape row is the only default.** An ordinary catch-all row
  overlaps everything it means to catch. An escape row is the construct
  that closes the remainder, and a group closed that way carries an
  advisory saying so.

Worked, CI-linted examples of both live in
[`models/examples/`](models/examples/), and
[docs/model-authoring.md](docs/model-authoring.md) walks through authoring
each.

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
`intrastate --help-all` prints the whole surface at once.

```sh
intrastate --help-all                    # every command's full reference
intrastate flow resolve --help-all       # one verb's refusal vocabulary
intrastate lint --help-all               # the live finding taxonomy
```

Those code lists are spelled through the same constants the wire is
emitted from, so they cannot drift from what a caller receives.
[docs/cli-reference.md](docs/cli-reference.md) is the same content as a
file, generated by `make docs` and gated by `make check`.

## Docs

- [docs/cli-reference.md](docs/cli-reference.md) — every command and
  flag (generated; do not hand-edit).
- [docs/model-authoring.md](docs/model-authoring.md) — authoring a
  transition model in TOML.
- [docs/cli-output-contract.md](docs/cli-output-contract.md) — worked
  JSON payloads and envelope rationale.
- [llms.txt](llms.txt) — agent-oriented index.

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
