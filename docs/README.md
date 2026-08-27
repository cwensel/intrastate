# intrastate docs

Reference material for intrastate. Keep documents short and current;
code is the source of truth.

**The binary is self-describing — prefer it over any file here.**
`intrastate --help-all` prints the full reference for every command:
the vocabulary, the flag grammar, the refusal codes a verb can return,
and the exit contract. Those lists are spelled through the same
constants the wire is emitted from, so they cannot drift from what a
caller actually receives.

## Generated

- [cli-reference.md](cli-reference.md) — every command, flag, and
  extended help body. **Generated** from the command tree by
  `make docs`; `make check` fails when it is stale. Do not hand-edit.
  The same content comes from `intrastate --help-all`.

## Hand-written

These carry what the CLI surface cannot state about itself:

- [model-authoring.md](model-authoring.md) — how to author a
  transition model in TOML: the model class
  (`state-machine` / `decision-table`), a decision table's cells, and
  the "otherwise" escape row. This is about the INPUT format, which the
  binary has no help surface for.
- [cli-output-contract.md](cli-output-contract.md) — worked JSON
  payloads and the rationale behind the envelope shape: why an
  aggregate failure reports in `findings[]` while a scalar one names
  `param`, and why closing coverage by escape is only an advisory.

## Decision records

- [rdr/](rdr/) and [jdr/](jdr/) — the locked design decisions behind
  the contracts above. History and rationale, never amended.
