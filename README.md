# intrastate

intrastate is a CLI for deterministic decisions and state transitions
inside agent workflows. Prompts supply facts, resolve declared rules,
and act on the result. State stays in existing workflow artifacts.

The project is motivated by
[Determinism Injection](https://chris.wensel.net/post/determinism-injection/):
externalize routine lookups, decisions, and writes so a prompt can focus
on the reasoning its task requires. intrastate provides the decision
models and state access; the prompt remains the orchestrator. Scripts
and other applications can use the same CLI.

## How It Fits into a Prompt

A large prompt may need to choose a review stage, enforce a retry limit,
check completion, or update a status field. Put that policy in a TOML
model, check it with `intrastate lint`, and call it at the relevant point
in the prompt:

1. **Collect facts.** Supply observed values with `--tag` and bind any
   artifacts the model needs to read.
2. **Resolve.** Call `flow resolve` with the model and outcome. It returns
   a selected rule, emitted values, and any planned writes.
3. **Act.** Run the selected task, resolve another decision, or report a
   stop. Model-defined dispositions such as `route`, `chain`, and `stop`
   let callers dispatch on a result's category.
4. **Apply.** When the workflow calls for a state change, pass the plan to
   `flow set-state`.

intrastate has no backing datastore. Declared accessors read and write
state in the artifacts your workflow already uses. You can start with
one decision table and add models as more decisions become explicit.

The [worked examples](docs/examples.md) show this pattern in the RDR
workflow, where models support stage selection, review loops,
implementation checks, delegation, and record updates.

## Capabilities

Both model classes use the same TOML grammar and CLI:

| Model class | Use it to | Result |
| --- | --- | --- |
| Decision table | Map current facts to an answer without owning state | Emitted values, optionally grouped into dispositions |
| State machine | Read owned state and determine a transition or decision | Planned writes and/or emitted values |

- **Check coverage and overlap.** For finite, declared domains, lint
  checks combinations of guarded values per outcome. A rule can cover
  multiple combinations. Gaps and overlaps produce findings; explicit
  escape rows can close gaps with a coverage advisory.
- **Resolve deterministically.** Given the same model and assembled
  inputs, the resolver selects one row or returns a structured refusal.
  Row order and priority do not break ties.
- **Validate outputs.** Declared emit domains constrain the answers a
  model can return, including commands, operations, and stop reasons.
- **Separate planning from writes.** Inspect a resolved plan before
  applying it. Accessors define how state is read and written; gates
  check selected transitions, and read-back can verify applied writes.
- **Handle failures explicitly.** Missing facts, ambiguous matches, and
  denied gates have stable refusal codes. A modeled stop is a successful
  decision with a stop value; a refusal exits nonzero.

See [Model authoring](docs/model-authoring.md) for the grammar, coverage
rules, and accessor options.

## Install

With the Go toolchain:

```sh
go install github.com/cwensel/intrastate/cmd/intrastate@latest
```

Or download a prebuilt binary for your platform from the
[latest release](https://github.com/cwensel/intrastate/releases/latest) —
each release carries a tarball per platform plus `checksums.txt`:

```sh
tar xzf intrastate_<version>_<os>_<arch>.tar.gz
install -m755 intrastate ~/.local/bin/
```

Verify a download before installing it:

```sh
sha256sum -c checksums.txt --ignore-missing
```

From a checkout:

```sh
make install        # builds ./bin/intrastate and installs to ~/.local/bin
make build          # ./bin/intrastate, without installing
```

## Quickstart

The following examples run from the root of a repository checkout and
use its bundled models. Build the local binary first:

```sh
git clone https://github.com/cwensel/intrastate.git
cd intrastate
make build
```

### Resolve a Decision

The pricing model maps `tier` and `region` to a plan and a data-processing
agreement requirement. Lint the model, then supply its inputs:

```sh
./bin/intrastate lint --model models/examples/pricing-decision-table.toml
./bin/intrastate flow resolve \
    --model models/examples/pricing-decision-table.toml \
    --outcome decide --tag tier=paid --tag region=eu
```

The response includes these emitted values:

```text
emit.dpa: required
emit.plan: pro
```

The caller can use `emit.plan` and `emit.dpa` directly. Omitting `region`
returns `flow-guard-unevaluable`, identifying the input needed to decide.

### Plan and Apply a State Change

The review model advances `draft → submitted → approved | rejected`.
Initialize its state in a temporary directory:

```sh
review_dir=$(mktemp -d)
./bin/intrastate flow init-state \
    --model models/examples/review-state-machine.toml \
    --artifact "review=$review_dir/state.json"
```

Resolve `submit` and save the plan. Resolution leaves the state unchanged:

```sh
./bin/intrastate flow resolve \
    --model models/examples/review-state-machine.toml \
    --artifact "review=$review_dir/state.json" \
    --outcome submit --as=json > "$review_dir/plan.json"
```

Apply the plan. This model enables read-back to verify the written value:

```sh
./bin/intrastate flow set-state \
    --model models/examples/review-state-machine.toml \
    --artifact "review=$review_dir/state.json" \
    --plan "$review_dir/plan.json"
```

The response confirms `owned.status: submitted` and
`writes.status: submitted`.

For state stored in an existing document, try the
[Markdown read/write example](docs/examples.md#reading-and-writing-a-markdown-document).
It reads status with `sed`, edits one line, and verifies the result.

### Output and Help

Use `--as=json` for a structured response with a terminal `type` of `ok`
or `failed`. Text output is the default. See the
[output contract](docs/cli-output-contract.md) for payloads and exit codes.

Use `--model <path>` to select a model. `--flow <id>` is reserved for
config discovery and currently has no registered IDs.

```sh
./bin/intrastate --help-all                # full command reference
./bin/intrastate flow resolve --help-all   # resolution options and refusals
./bin/intrastate lint --help-all           # lint options and findings
```

## Documentation

- [Examples](docs/examples.md) — capabilities, diagrams, and prompt
  integration patterns from bundled models and the RDR workflow.
- [Model authoring](docs/model-authoring.md) — define decision tables,
  state machines, guards, outputs, and accessors in TOML.
- [CLI reference](docs/cli-reference.md) — every command and flag,
  generated from the command tree.
- [Output contract](docs/cli-output-contract.md) — JSON payloads,
  refusals, and exit codes for callers.
- [llms.txt](llms.txt) — agent-oriented documentation index.

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for the Issues-only contribution
policy, local development setup, and CLI conventions. The locked
architecture decisions live under [docs/rdr](docs/rdr).

```sh
make check          # full local validation (mirrors CI)
make test           # race + coverage
```

## License

[MIT](LICENSE)
