# Model authoring

A model is one TOML document. This document covers the choices an author
makes above the grammar; the loader and `intrastate lint` are the
authoritative checks, and code is the source of truth.

Two model classes ship, and they answer different questions:

- a **state machine** answers *what may happen next, and what does that
  change* — it owns state and advances it;
- a **decision table** answers *what is the answer for this situation* —
  it owns no state and advances nothing.

Both are authored in the same grammar, both are linted by the same
analysis, and both are driven by the same four `flow` verbs. The class is
declared in `[model]`, and it decides which invariants apply.

Worked, CI-linted examples of each live under
[`models/examples/`](../models/examples/): a
[decision table](../models/examples/pricing-decision-table.toml) and a
[state machine](../models/examples/review-state-machine.toml). They are
linted by `make graph-lint`, so an example this document describes cannot
drift from one the tool accepts.

## Model class

`[model]` may carry `class`, whose only admitted values are
`"state-machine"` and `"decision-table"`. An absent `class` reads as
`"state-machine"`.

| | `state-machine` | `decision-table` |
| --- | --- | --- |
| owned tags | one or more | **zero** (declaring one is refused at load) |
| `[initial]` | required — the root owned state | none |
| `terminal` / `[context.*]` | declares where the machine stops | none |
| accessors | `[read.*]` and `[write.*]` per owned tag | none |
| ordinary rule | carries `[rule.write]` | carries `[rule.emit]`, no write |
| `flow resolve` answers | `rule` + `writes` / `next` | `rule` + `emit` |
| `--artifact` at runtime | required — readers must be bound | not required |

The class is **declared, not inferred** from the owned set. A
`decision-table` model that declares an owned tag is refused at load as a
malformed model declaration, whose detail names the class and the count as
the token `owned=<n>`. The check is one-directional: a `state-machine`
declaring zero owned tags is not refused — it is rootless, which lint
reports as `graph-dangling-edge`.

Declaring the class rather than deriving it is what makes the refusal
possible. A tool that inferred "no owned tags, therefore a table" would
silently accept a state machine whose owned declaration an author had
dropped, and answer with an empty plan instead of refusing.

## What both classes share

The class changes which invariants apply. It does **not** buy an exemption
from the analysis, and one shared property surprises authors most:

**Coverage is proved per outcome, for both classes.** Rules are grouped by
the outcome they match, and each group's guard atoms span a product of
their declared domains. Every assignment in that product must be claimed
by some row in the group, or lint reports `graph-coverage-gap` naming how
many assignments are uncovered and over which dimensions.

That is intuitive for a table — it is the table's whole point. It applies
just as strictly to a machine. In the
[state-machine example](../models/examples/review-state-machine.toml),
`submit` is guarded on `status = "draft"` over a four-value domain, so
three of the four assignments are unaccounted for, and the model does not
lint until something covers them. Escape rows are what close it there,
exactly as they close a table's unclaimed cells.

The other shared properties:

- **Exactly one row, or a refusal.** Two rows enabled by the same
  assignment is `graph-overlap` at lint and `flow-ambiguous-match` at
  runtime. Nothing picks between them by order, priority, or authoring
  position.
- **The escape row is the only default.** An ordinary catch-all row
  overlaps everything it is meant to catch; the escape row is the
  construct that does not.
- **The same four verbs.** `flow next`, `flow resolve`, `flow read-state`,
  and `flow set-state` drive both classes — though a decision table has no
  state to read or set, so only `next` and `resolve` are meaningful for it.

## Authoring a state machine

A machine's rules advance owned state. The
[worked example](../models/examples/review-state-machine.toml) is
`draft -> submitted -> approved | rejected`.

### The owned state and its accessors

An owned tag is state the model owns and the CLI persists. Each one needs
exactly one reader and, if any rule writes it, exactly one writer — an
owned tag served by zero readers is refused at load as
`malformed_accessor_binding`:

```toml
[tags.status]
provenance = "owned"
kind = "enum"
domain = ["draft", "submitted", "approved", "rejected"]
single_valued = true
required = true

[read.review-state]
role = "review"
path = "review.status"
keys = ["status"]
timeout = "5s"

[write.review-state]
role = "review"
path = "review.status"
keys = ["status"]
timeout = "5s"
read_back = true
```

`role` is the binding a caller supplies at runtime as
`--artifact review=<path>`; nothing is discovered. `read_back = true` is
what makes `flow set-state` verify a write by reading it back rather than
trusting that it landed.

### The root

`[initial]` declares the owned state a model starts from. Every
always-present owned key must be assigned here; one that is declared
always-present and omitted takes `graph-always-present-owned`.

```toml
[initial]
status = "draft"
```

A machine with no `[initial]` is rootless: lint reports
`graph-dangling-edge` against the model, because a transition graph with
no root has nothing to reach its rules from.

### Rules advance the state

An ordinary rule matches an outcome, guards on the current state, and
writes the next one:

```toml
[[rule]]
id = "submit"
[rule.match.recognized]
eq = "submit"
[rule.guard.all.status]
eq = "draft"
[rule.write]
status = "submitted"
```

`flow resolve --outcome submit` over `status = "draft"` selects this row
and returns the write as a **plan**. It applies nothing: `flow set-state`
is what writes, and nothing links the two calls.

### Terminals are predicates, not state names

`terminal` lists the names of `[context.<name>]` blocks, and each block is
a predicate over owned tags:

```toml
terminal = ["approved", "rejected"]

[context.approved]
[context.approved.match.status]
eq = "approved"

[context.rejected]
[context.rejected.match.status]
eq = "rejected"
```

Naming a value directly — `terminal = ["approved"]` with no matching
`[context.approved]` — is refused at load as `unknown_context`. The
indirection buys expressiveness: a terminal can be a predicate over
several owned tags at once, which a bare state name cannot express.

A reachable non-terminal state with no outgoing rule is a
`graph-dead-end`: the machine can arrive somewhere it can never leave and
was never declared finished.

### Covering the state space

Coverage applies per outcome, as it does for a table. `submit` guarded on
`status = "draft"` accounts for one of `status`'s four values, so lint
refuses until the other three are covered. Two ways to close it:

- **Pin the remaining transitions.** If `submitted` should also accept
  `submit` as a no-op, author that row. Coverage is then *proved*.
- **Add an escape row for the outcome.** The request is legal but the
  state does not admit the transition, so it is rescued rather than
  refused.

```toml
[[rule]]
id = "submit-otherwise"
escape = ["no_match"]
[rule.match.recognized]
eq = "submit"
```

`flow resolve --outcome approve` over `status = "draft"` then selects
`approve-otherwise` and returns `escaped: true` with
`escape_class: "no_match"`, instead of refusing. The group takes the
`graph-coverage-closed-by-escape` advisory — informational, exit 0 — which
is how an author sees the model is leaning on a default rather than
enumerating.

## Authoring a decision table

A decision table's rows claim cells. Each row binds the outcome with a
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
it as written yields a blocking `graph-coverage-gap`. A real table covers
all four cells, whether with one row per cell or with rows whose guard
atoms pin fewer dimensions and so span several, and closes any remainder
with the escape row below.

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
