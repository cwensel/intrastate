# Model authoring

A model is one TOML document. This document covers the choices an author
makes above the grammar; the loader and `intrastate lint` are the
authoritative checks, and code is the source of truth.

## Model class

`[model]` may carry `class`, whose only admitted values are
`"state-machine"` and `"decision-table"`. An absent `class` is read as
`"state-machine"`, so every existing model keeps its meaning.

- **`state-machine`** — the model carries owned tags, declares an
  `[initial]` owned state, and its rules advance that state. Every ordinary
  rule contains a write block.
- **`decision-table`** — the model carries **zero** tags of provenance
  `owned`. It holds no state and advances nothing: it maps a supplied
  situation to an answer. An ordinary rule needs no write block, and the
  model declares no `[initial]`, no `terminal`, and no owned accessors.

The class is **declared, not inferred** from the owned set. A
`decision-table` model that declares an owned tag is refused at load as a
malformed model declaration, whose detail names the class and the count as
the token `owned=<n>`. The check is one-directional: a `state-machine`
declaring zero owned tags is not refused — it is rootless, which lint
reports as `graph-dangling-edge`.

## Authoring a decision table

A decision table's rows are cells. Each row binds the outcome with a
`[rule.match.recognized]` atom and discriminates the cell with **guard**
atoms:

```toml
outcomes = ["decide"]

[model]
id = "pricing"
version = 1
class = "decision-table"

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[tags.tier]
provenance = "observed"
kind = "enum"
domain = ["free", "paid"]
single_valued = true
required = true

[tags.region]
provenance = "observed"
kind = "enum"
domain = ["eu", "us"]
single_valued = true
required = true

[[rule]]
id = "free-eu"
[rule.match.recognized]
eq = "decide"
[rule.guard.all.tier]
eq = "free"
[rule.guard.all.region]
eq = "eu"
[rule.emit]
plan = "basic"
dpa = "required"
```

This snippet shows one row of the table's shape, not a complete model:
`tier × region` has four cells and only `free-eu` is claimed, so linting
it as written yields a blocking `graph-coverage-gap`. A real table authors
all four rows, or closes the rest with the escape row below.

`[rule.emit]` is the row's answer: a flat block of string values, returned
by `flow resolve` under the `emit` key. Its keys are not tags — nothing
declares them, nothing writes them, and they take no part in selection.

### Discriminate with guard atoms, not `[rule.match.<key>]`

**A decision table's discriminating dimensions must be authored as
guard atoms — `[rule.guard.all.<key>]` or `[rule.guard.unless.<key>]` — not
as `[rule.match.<key>]`.** This is the one authoring choice that
lint cannot forgive quietly, and it is easy to get wrong because both block
kinds accept the same operators over the same tags.

The reason is what each block does to the coverage claim. A **match** atom
*scopes* the row group — it decides which rows are compared against one
another — and contributes **no dimension** to the product lint proves
exhaustiveness over. Only `guard.all` and `guard.unless` atoms collect as
dimensions. A table discriminated by match atoms therefore ranges over an
empty product, which any single row's coverage equals, so the group closes
clean **whether or not the table is complete** — a silent green over a
table with holes in it.

Lint refuses that silence for this class. A `decision-table` group whose
scoped product has zero participating dimensions takes
`graph-unprovable-coverage` carrying `reason = no-participating-dimension`,
naming the group and pointing at the guard-atom remedy. Over a
`state-machine` the same shape is legitimate and stays silent — a machine's
group may genuinely range over no guard dimension.

With the dimensions authored as guards, lint proves the table for real:
`graph-coverage-gap` reports, per unclosed failure-class arm, how many of
the product's cells are uncovered and over which dimensions; `graph-overlap`
names two rows that both claim one, and the table is only green when every
cell is claimed exactly once.

### The escape row: the "otherwise" idiom

An ordinary catch-all row overlaps every other row and yields
`flow-ambiguous-match`. The way to author a default is an **escape row**
rescuing `no_match`:

```toml
[[rule]]
id = "otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "decide"
[rule.emit]
plan = "unclassified"
```

An escape row carries no write block, no `clear`, and no gate — but it may
carry its own `[rule.emit]`, and a plan it rescues returns that row's
answer with `escaped: true` and the `escape_class` that was rescued.

Two things to know about it:

- **Rescue is scoped per outcome.** A table over an N-outcome alphabet
  needs N escape rows — or one rule whose match block uses an `in` atom
  over the outcomes, which expands to N rows. An outcome with no rescue row
  is not silently defaulted: its group takes its ordinary
  `graph-coverage-gap`.
- **Closing by escape is reported.** A group whose coverage is closed by a
  bare escape row rather than proved over its declared domains takes the
  `graph-coverage-closed-by-escape` advisory. It is informational and does
  not fail the lint, but a bare green is deliberately not available: the
  advisory is how an author sees that the table leans on its default
  instead of enumerating.

## Invoking a decision table

A decision table declares no owned tag, so no reader is invoked and
`--artifact` is not required. The discriminating dimensions arrive as
`--tag`, and `--outcome` is required exactly as it is for a state machine:

```sh
intrastate flow resolve --model pricing.toml \
  --outcome decide --tag tier=free --tag region=eu --as=json
```

The answer is `rule` plus `emit`; `next`, `writes`, `clear`, `owned`, and
`readers` are all empty over this class. See
[cli-output-contract.md](cli-output-contract.md) for the payload.
