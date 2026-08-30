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

Worked, CI-linted examples live under
[`models/examples/`](../models/examples/): a minimal
[decision table](../models/examples/pricing-decision-table.toml), a minimal
[state machine](../models/examples/review-state-machine.toml), and a
[grammar-surface model](../models/examples/release-grammar.toml) carrying
every construct [The grammar](#the-grammar) *documents*. All three are
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

## The grammar

Everything below is common to both classes: the same tag declarations, the
same operators, the same rule shape. What differs is which parts a class
uses, which the sections after this one cover.

A model's root admits a closed set of keys — `outcomes`, `terminal`,
`[model]`, `[initial]`, `[tags.<key>]`, `[read.<id>]`, `[write.<id>]`,
`[gate.<id>]`, `[context.<id>]`, `[[rule]]`, and `[dump]`. Decoding is
strict: an unrecognised key is refused as `unknown_schema_field`, never
ignored. A misspelled facet is a stable refusal rather than a silent no-op.

The worked model for this section is
[`release-grammar.toml`](../models/examples/release-grammar.toml), linted
by `make graph-lint`. The snippets in this section are taken from it, so
none of them describes a construct the tool rejects. It carries every
construct this section documents, with two deliberate exceptions: `[dump]`
and the rule-level `source` key are each named here — `[dump]` in the
root-key set above, `source` in the placement list below — but neither is
otherwise covered, and the example authors neither.

### Declaring a tag

A tag declaration carries a **provenance**, a **kind**, and the facets its
kind admits.

Provenance is one of three, and it decides who supplies the value:

| provenance | supplied by | notes |
| --- | --- | --- |
| `owned` | the model's `[read.*]` accessors | state the model owns; a `decision-table` declares none |
| `observed` | the caller, as `--tag key=value` | may also be served by at most one reader |
| `recognized` | the kernel, from `--outcome` | exactly the reserved `recognized` key; no accessor may name it |

Kind is one of exactly **five** tokens, and each admits a different facet.
Authoring a *populated* facet on the wrong kind is refused at load as
`malformed_tag_declaration`:

| kind | facet it carries | `single_valued` | finite dimension? |
| --- | --- | --- | --- |
| `enum` | `domain = [...]` | admitted | yes, given a `domain` |
| `bool` | none — its two literals are the domain | admitted | yes, always |
| `int` | `min` / `max` | admitted | only with **both** bounds |
| `set` | `elements = [...]` | **refused** | yes, given `elements` |
| `scalar` | none | **refused** | **never** |

`single_valued = true` says the tag holds **at most one** of its domain's
values — one when the key is present, and it constrains nothing when the
key is absent. Pairing it with `required = true` is what yields exactly
one. Without it the values may co-occur, and the dimension becomes one
independent boolean per value — `2^|domain|` rather than `|domain|`. It is
refused outright on `set` and `scalar`: a set holds any *subset* of its
element universe, so its space is `2^|elements|` by construction and there
is no single-valued reading to select; a scalar has no declared domain to
partition at all.

Two of those refusals key on the facet being **non-empty**, not on it
being written: `domain` and `elements` are checked by length, so an empty
`domain = []` or `elements = []` on a kind that admits neither loads
without complaint. It declares nothing and changes nothing — the kind's
own dimension rules still decide — but do not read a clean load of an
empty array as the kind having accepted the facet. `single_valued` is
checked by presence instead, so `single_valued = false` on a `set` or
`scalar` *is* refused.

`required = true` says the key is always present. An optional key doubles
its dimension with a `{present, absent}` factor, and — more consequentially
— any *value* atom over it can refuse `flow-guard-unevaluable` at runtime
when the key is missing, which lint reports as `graph-unprovable-coverage`
with `reason = row-can-refuse`. **A dimension a decision table
discriminates on therefore wants `required = true`**; a sentinel member in
the domain, always emitted by the producer, is the idiom for a value that
is legitimately sometimes unknown.

`single_valued = true` is the companion recommendation, but only for the
**single-value** operators — `eq`, `in`, and the four integer
comparisons — on an `enum`, `bool`, or `int` dimension. It keeps the
dimension `|domain|` rather than `2^|domain|`, and a guard row over a
multi-valued tag is what lint reports as unprovable. The other two
operators are deliberate exceptions:

- **`contains`** discriminates on a `set`, which *refuses*
  `single_valued` by construction. A `contains` dimension wants
  `required = true` alone; its `2^|elements|` space is the point, not a
  defect.
- **`exists`** discriminates on *presence*, so the key must be
  **optional**. `required = true` collapses the `{present, absent}`
  factor to one, leaving the atom constant either way — `exists = true`
  denotes the whole dimension, `exists = false` denotes none of it — and
  lint reports both as `graph-vacuous-atom`.

```toml
# an enum: the only kind that carries a domain
[tags.phase]
provenance = "owned"
kind = "enum"
domain = ["idle", "building", "shipped", "held"]
single_valued = true
required = true

# an int: both bounds, or the dimension is not finite
[tags.risk]
provenance = "observed"
kind = "int"
min = 0
max = 2
single_valued = true
required = true

# a set: `elements`, and no single-valued marker
[tags.checks]
provenance = "observed"
kind = "set"
elements = ["tests", "signoff"]
required = true

# a scalar: opaque, no facets, no exhaustiveness claim
[tags.build-id]
provenance = "owned"
kind = "scalar"
```

A dimension lint cannot enumerate — a `scalar`, an `int` missing a bound, an
`enum` with no `domain` — takes `graph-unprovable-coverage` with
`reason = dimension-not-finite`. A `scalar` is therefore fine to *carry* and
never usable to *discriminate on*.

### Guard and match atoms

An atom is one operator applied to one tag, authored as a sub-table named
for the tag:

```toml
[rule.guard.all.risk]
gte = 1
```

The operator vocabulary is **closed at eight**, and an unadmitted token is
refused as `malformed_predicate_atom` in every block. Which *kinds* each
operator accepts is decided by the matrix below, and an unadmitted
operator/kind pair is refused as `malformed_predicate_atom` too:

| operator | accepts kinds | literal shape |
| --- | --- | --- |
| `eq` | `enum` `bool` `int` `scalar` | one typed scalar |
| `in` | `enum` `bool` `int` `scalar` | non-empty typed scalar set |
| `lt` `lte` `gt` `gte` | `int` | one typed integer |
| `contains` | `set` | non-empty typed element set |
| `exists` | all five | `true` or `false` |

**The matrix binds guard atoms only** — `[rule.guard.all.<key>]` and
`[rule.guard.unless.<key>]`. Guard atoms are what become product
dimensions, so the matrix is what keeps a dimension's operator meaningful
over its kind. A match block is scoped separately: it admits only `eq` and
`in`, and over **any** declared kind, `set` included. A
`[rule.match.checks]` block carrying `eq = "tests"` is therefore
well-formed where `[rule.guard.all.checks]` carrying the same atom is
not.

What still binds every block is the **literal**. A member is checked
against the tag's kind — an `int` tag takes only parseable integers, a
`bool` tag only `true`/`false` — and against its declared domain: an
`enum`'s `domain`, a `set`'s `elements`, an `int`'s `min`/`max`. An
out-of-domain member is `malformed_predicate_atom` wherever it is
authored. A comparison bound is kind-checked but not domain-checked, so
`gte = 9` over `min = 0, max = 2` loads and simply denotes nothing.

`contains` asks whether the held set is a **superset** of the literal, and
it is total: a held `[]` answers false rather than refusing. `exists` reads
presence alone and never the value, so it is the one operator that decides
something over an optional key. Over an always-present key it is constant,
and which constant depends on the literal: `exists = true` holds
everywhere and denotes the whole dimension, `exists = false` holds nowhere
and denotes none of it. Neither discriminates, and lint reports both as
`graph-vacuous-atom` — the advisory keys on the operator and the key's
`required` marker, not on the literal. A set literal is order-insensitive
and a repeated member is refused, not collapsed.

Three blocks take atoms, and the block decides what the atom does:

- **`[rule.guard.all.<key>]`** — a conjunct. Every atom must hold.
- **`[rule.guard.unless.<key>]`** — the negated block. Negation is
  **block-level, not per atom**: every `unless` atom is conjoined and the
  conjunction is negated *once*.
- **`[rule.match.<key>]`** — scoping, not discrimination. A match block
  admits **only `eq` and `in`**; any other operator is refused. An `in`
  atom here *expands* into one row per member.

The row's guard verdict is therefore

```
all₁ ∧ all₂ ∧ … ∧ ¬(unless₁ ∧ unless₂ ∧ …)
```

with the `unless` term contributing nothing when the block is absent.
**With two or more `unless` atoms this is weaker than negating each one.**
The row is enabled where *at least one* `unless` atom fails — not where
they all fail.

If what you want is `¬u₁ ∧ ¬u₂`, express each negation in `guard.all`
with the operator's complement — `gte` against `lt`, `eq` against an `in`
over the rest of the domain — since `guard.all` conjoins. Splitting the
atoms across two rules does **not** work: that yields two rows, one
enabled by `¬u₁` and one by `¬u₂`, which is the disjunction, and where
both hold the rows overlap into a `graph-overlap` finding or an
`ambiguous_match` refusal. Where no complement is expressible the
conjunction needs the model restructured, not a second `unless` block.

Only `guard.all` and `guard.unless` atoms become dimensions in the product
lint proves coverage over. That distinction is the single most consequential
authoring choice for a decision table, and
[its own section](#discriminate-with-guard-atoms-not-rulematchkey) covers
why.

### Shared contexts

A `[context.<id>]` block is a named match block. A rule pulls one in with a
rule-level `use` list, and a context chains onto another with `inherits`:

```toml
[context.shipping]
[context.shipping.match.recognized]
eq = "ship"

[context.expedited]
inherits = "shipping"

[[rule]]
id = "ship-clean"
use = ["expedited"]
[rule.match.recognized]
eq = "ship"
```

Inheritance **accumulates** — it never overrides. Merging is by full atom
identity, so an atom a rule and a context both contribute collapses to one,
and two atoms differing only in literal both survive. An `inherits` cycle is
refused as `cyclic_context_inheritance`; naming a context that does not
exist is `unknown_context`.

One ordering rule is deliberate and catches authors out: **`use` cannot
supply the rule's match block, only atoms within it.** The "every rule
carries a match block" check runs *before* inheritance and keys on the
block's *presence* — a rule with no `[rule.match.*]` sub-table at all is
refused as `malformed_rule_shape` however much its contexts would have
contributed.

What that check does **not** demand is a local `recognized` atom. Outcome
binding runs *after* the merge and reads the merged match atoms, local and
inherited alike, so a rule carrying some other local match atom may take
its `recognized` binding entirely from a context. A `recognized` atom
authored under `guard.all` or `guard.unless` is a different matter: it is
refused as `malformed_outcome_binding` rather than lifted, because only
match blocks bind outcomes.

The example rule above repeats its `recognized` atom regardless, for a
reason unrelated to shape: a match atom other than `recognized` *scopes*
the coverage group, and scoping that rule separately would break the
partition it forms with its siblings.

Contexts serve a second job unrelated to rules: `terminal` names them, which
[the terminals section](#terminals-are-predicates-not-state-names) covers.

### Rule-level keys and where they go

`id`, `use`, `gate`, `clear`, `escape`, and `source` belong to the
`[[rule]]` table itself, so they **must precede the first `[rule.*]`
sub-table**. TOML reads
a key written after one as belonging to *that* table — `clear` placed after
`[rule.guard.all.build-id]` parses as a guard atom and refuses with
`unknown operator`. It is a plain TOML rule with a confusing error, and it
is worth knowing before it happens:

```toml
[[rule]]
id = "ship-clean"
use = ["expedited"]
clear = ["build-id"]
[rule.match.recognized]
eq = "ship"
[rule.guard.all.phase]
eq = "building"
[rule.write]
phase = "shipped"
```

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

## Where the two classes come from

Neither class is invented here. Both are long-standing specification
techniques, and intrastate's invariants are the classical checks on each,
enforced before runtime rather than discovered in production.

### The state machine

A model of this class is a finite transition system: a finite owned-state
space, a recognized-outcome alphabet, and a transition relation authored
as rows.

Where the output is attached distinguishes the two classical variants.
**Mealy** machines associate output with the *transition* — the (state,
input) pair — while **Moore** machines make it a function of the state
alone [1][2]. intrastate's rules are Mealy-shaped: `[rule.write]` and
`[rule.emit]` hang off the row, which is a (state, outcome) pair, not off
the state.

Determinism is the property that a machine has **exactly one** transition
per (state, input) pair; a nondeterministic one admits a *set* of next
states [3][4]. intrastate refuses nondeterminism statically rather than
resolving it: two rows enabled by the same assignment is `graph-overlap`
at lint, and `flow-ambiguous-match` if it is ever reached at runtime.

The guard blocks are a guarded-command construct in Dijkstra's sense,
where nondeterminacy arises precisely when several guards are
simultaneously true and "the order in which the guarded commands of a set
appear in our text is semantically irrelevant" [5]. intrastate takes the
same position on order — authoring position never breaks a tie — but
diverges on the consequence: Dijkstra's alternative construct selects an
arbitrary true-guarded list, where intrastate refuses. An arbitrary choice
is not reviewable, and a workflow that silently picked one of two rules
would be unauditable.

Two classical reachability notions are separate findings here, and the
distinction is the standard one. An **inaccessible** state has no path
from the start state; a **dead** (or trap) state is reachable but cannot
reach an accepting state [4]. `graph-unreachable-rule` reports the first —
"no reachable owned-state satisfies the selection context of row X" — and
`graph-dead-end` the second: a *reachable* state that satisfies no
declared terminal and is the source of no non-escape row.

Statecharts extend flat machines with hierarchy, concurrency, and
communication, motivated by the difficulty of describing reactive
behaviour in ways that are "clear and realistic, and at the same time
formal and rigorous" [6]. intrastate deliberately does not take that
extension: a model is one flat rule set over a tag-valued state.
Hierarchy is the natural next step if flat models stop scaling, not
something the current design provides.

### The decision table

Decision tables are a specification technique from commercial data
processing, first reported in 1957 [7]. General Electric's TABSOL made the
form executable in 1960 [8], and IBM's 1962 manual gave the technique its
canonical name and treatment [9]. The technique was standardized late and
repeatedly — a CODASYL Decision Table Task Group ran from 1973 and
reported in 1982 [10].

The classical structure is a condition part that is "a truth table (from
propositional logic) that has been rotated 90°," which "guarantees that we
consider every possible combination of condition values" [11]. That is
what intrastate's coverage proof checks. The classical vocabulary also
distinguishes **limited-entry** tables, where every condition is binary,
from **extended-entry** tables, where a condition may take several values
[11][12]. intrastate's tables are extended-entry: a dimension is a tag's
declared domain, of any finite size, not a yes/no.

The scaling consequence is the classical one: a limited-entry table over
n conditions has 2^n rules [11]. intrastate's product bound (published in
`intrastate lint --help-all`) is the point past which the analysis
declines to prove coverage rather than enumerate indefinitely.

Two defects have been the subject of formal analysis since the 1960s:
**ambiguity**, where more than one rule covers a case [13], and
**redundancy**. Jorgensen puts the consequence directly: when two rules
apply to the same transaction, "1. Rules 4 and 9 are inconsistent. 2. The
decision table is nondeterministic" [11]. intrastate reports the first as
`graph-overlap` and the second as `graph-redundant-row`, and treats
ambiguity as blocking.

**The modern standard is DMN** (Decision Model and Notation), whose
decision tables carry a *hit policy* saying what to do when several rules
match: single-hit **Unique, Any, Priority, First** and multiple-hit
**Output order, Rule order, Collect**, defaulting to Unique [14].

This is where intrastate deliberately parts company with common practice.
A `First` or `Priority` policy resolves ambiguity by authoring order or by
a declared ranking; intrastate provides no such policy, because a
tie-break the model did not author is a decision no reviewer approved.
Overlap is refused at lint instead. The nearest DMN analogue to
intrastate's design is the default policy, `Unique` — the one policy under
which overlap cannot arise.

DMN defines completeness the same way — "a decision table will be
considered complete if its rules cover all combinations of expected input
values for all input expressions" — and provides a `defaultOutputEntry`
for the unmatched case, with a rule that is exactly intrastate's advisory:
"A complete decision table SHALL NOT specify a default output value" [14].
An escape row and a proved product are alternatives, not companions, which
is why closing coverage by escape is reported rather than silent.

The **ELSE** mechanism has a long history and a long-standing critique.
Vanthienen and Dries argue against the ELSE-column specifically because it
becomes "a waste basket in which all kinds of hidden combinations of
conditions are put," and prefer an explicit `OTHER` state per condition,
which "eliminates 'rule ambiguity' ... and simplifies testing for
completeness" [7]. intrastate's escape row is closer to their `OTHER` than
to the ELSE-column: it is scoped per outcome, it is declared with the
failure class it rescues, and a plan that came through it is marked
`escaped: true` rather than being indistinguishable from an ordinary hit.

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

Every accessor entry needs all four of `role`, `path`, `keys`, and a
positive Go-duration `timeout` — any one absent or empty is
`malformed_accessor_declaration`. `read_back` is a **write-entry key**:
a write entry must carry `read_back = true`, and a read or gate entry
carrying it at all is refused. No accessor of any capability may name the
reserved `recognized` key.

#### Gates

`[gate.<id>]` is the third accessor capability. It answers **allow, deny,
or indeterminate** — it reads nothing into the tag view and writes nothing.
A rule opts into one with a rule-level `gate` list:

```toml
[gate.change-window]
role = "release"
path = "release.window"
keys = ["phase"]
timeout = "5s"

[[rule]]
id = "begin"
gate = ["change-window"]
[rule.match.recognized]
in = ["build"]
[rule.guard.all.phase]
eq = "idle"
[rule.write]
phase = "building"
build-id = "pending"
```

Gates run **after** exact-one selection and before the plan is emitted.
That placement is the whole point: a gate never prunes a candidate the
model would otherwise have chosen, it denies a plan already chosen. Under
`flow resolve` a deny is the refusal `flow-gate-denied` — never a plan,
never a silent fall-through to a second row. `flow next` invokes no gate at
all unless `--evaluate-gates` is passed, and reports each id as
not-evaluated otherwise.

Naming a gate no `[gate.<id>]` block declares is refused as
`unknown_accessor`. An escape row may carry no gate list at all, even an
empty one.

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

A write **replaces**; it never merges. For a `set` key the array literal is
the whole new set, and on every other kind a multi-member value is refused
rather than truncated. Every key a `[rule.write]` block names must be an
owned tag — writing an observed or recognized one is `write_to_non_owned_tag`
— and a value outside its declared domain is refused at load, so the
declaration is never advisory on the path that persists.

#### Removing an owned tag: `clear`

Absence from a write block never implies deletion. Removing an owned key is
a rule-level `clear` list, and normalization renders each named key as a
`<clear>` write:

```toml
[[rule]]
id = "ship-clean"
clear = ["build-id"]
[rule.match.recognized]
eq = "ship"
[rule.guard.all.phase]
eq = "building"
[rule.write]
phase = "shipped"
```

`clear` carries the same obligations a write does: an undeclared key is
`unknown_tag`, a non-owned one is `write_to_non_owned_tag`. A rule may
carry both blocks, as above, and the cleared keys join the written ones in
the plan's owned-key set.

`<clear>` is a **reserved TAG value**, refused in the three places a tag
value is authored: a `[rule.write]` block, `[initial]`, and any atom
literal. The `clear` list is the only way to produce it as a tag value,
and it surfaces the same way at the CLI: `flow set-state --clear <key>`,
never `--write key=<clear>`.

The reservation does not reach `[rule.emit]`. Emit keys are not TAG keys:
`[tags.<key>]` never declares one, nothing writes one, and their values
are never parsed or canonicalized. So `plan = "<clear>"` in an emit block
carries no clearing meaning there — it is the literal text, and it is
answered as that text.

Whether it LOADS is a separate question, and the model's own to answer. An
emit key MAY be declared in the top-level `[emit]` table
([below](#declaring-the-emit-vocabulary)); if `plan` is declared `enum`
over a domain that does not list `<clear>`, that row refuses at load as
`emit_value_out_of_domain`. Declare nothing and it loads, as it always
has.

Note the placement — `clear` is a rule-level key, so it precedes the first
`[rule.*]` sub-table. See
[rule-level keys](#rule-level-keys-and-where-they-go).

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
by `flow resolve` under the `emit` key. Its keys are not TAG keys —
`[tags.<key>]` never declares one, nothing writes them, and they take no
part in selection. They MAY carry their own declaration, in the top-level
`[emit]` table.

#### Declaring the emit vocabulary

Left undeclared, an emit block is free text: a misspelled key or a typo'd
value loads clean and is answered verbatim, because nothing knows what the
key was supposed to be. Declaring the vocabulary makes it a checked fact.

```toml
[emit.plan]
kind = "enum"
domain = ["basic", "pro"]

[emit.dpa]
kind = "enum"
[emit.dpa.domain]
gate = ["required"]
clear = ["none"]
```

`kind` is required and is one of `enum`, `bool`, `int`, `scalar`. `set` is
not admitted, because an emit value is one authored string rather than a
member sequence. Only an `enum` takes a `domain`; `scalar` admits any
value and is the documented escape hatch for a key you want ADMITTED but
not checked.

The `domain` key has two spellings. `domain = [...]` is a flat member
array. `[emit.<key>.domain]` is a sub-table whose keys are your own
disposition tokens and whose values are member arrays — the domain is the
union, each member carries the one disposition it is listed under, and
`flow resolve` surfaces that token in its `dispositions` object. intrastate
fixes no disposition vocabulary and never interprets a token.

Declaring is **opt-in and whole-model**. The trigger is the COUNT of
declared keys, so a model with no `[emit]` table — or a bare `[emit]` with
no sub-tables — behaves exactly as it does today and no new refusal is
reachable. Declare ONE key and the whole model is strict: every rule's
emit keys must be declared, on ordinary and escape rows alike.

Three refusals become reachable, all at load, so `intrastate lint` exits
nonzero and `flow resolve` answers `flow-model-invalid`:

- `malformed_emit_declaration` — the declaration itself is ill-formed: an
  unknown `kind`, an `enum` with no usable domain, an empty or duplicated
  member, a `domain` on a non-enum kind, an empty disposition token.
- `unknown_emit_key` — a rule emits a key `[emit]` does not declare.
- `emit_value_out_of_domain` — an authored value is outside its key's
  declared domain, is not a `bool` token, or is not an `int` literal.

An emit declaration is not a tag declaration: it takes no `provenance`, no
`min`/`max`/`elements`/`single_valued`/`required`, and a declared emit key
is still barred from match, guard, write, and accessor use. It is also
author-owned and unversioned — widen or narrow a domain by editing the
model, and nothing checks the edit against a history. A declared key no
rule emits, and a declared member no rule authors, are findings of no
tier.

`models/examples/pricing-decision-table.toml` and
`models/examples/routing-decision-table.toml` both ship declared; the
second is the partitioned-domain worked example.

### Discriminate with guard atoms, not `[rule.match.<key>]`

**A decision table's discriminating dimensions must be authored as
guard atoms — `[rule.guard.all.<key>]` or `[rule.guard.unless.<key>]` — not
as `[rule.match.<key>]`.** This is the one authoring choice that
lint cannot forgive quietly, and it is easy to get wrong because the two
blocks share a syntax and overlap on the operators most tables reach for:
`eq` and `in` are admitted in either. They are not interchangeable beyond
that overlap — a guard atom must satisfy the
[operator/kind matrix](#guard-and-match-atoms) while a match atom takes
only `eq` and `in`, over any kind — so a block swapped for the other can
turn from silently-accepted into a loader refusal depending on the
operator. The failure this section is about is the silent one.

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

### Passing a set-valued tag

A `set` tag's value crosses as a JSON array literal, and the model's
declaration is what tells the parser to expect one:

```sh
intrastate flow resolve --model triage.toml \
  --outcome triage --tag 'labels=["security"]' --tag severity=3 --as=json
```

**An undeclared tag whose value is an array is refused.** A `--tag` naming
a key the model does not declare passes through harmlessly *as a scalar* —
but it is validated against the zero declaration, whose kind is not `set`,
so an array literal there returns `flow-tag-invalid` ("the tag `x` is not
set-valued") and the call exits 2.

This bites a caller that pipes a producer's whole tag output through one
invocation: **every set-valued key the caller passes must be declared, even
one no row guards on.** Declaring it costs the coverage product nothing —
only guard atoms contribute dimensions, so a declared-but-unguarded tag adds
no cell to prove:

```toml
# received and ignored; declared only so the array literal parses
[tags.labels]
provenance = "observed"
kind = "set"
elements = ["security", "docs"]
required = true
```

## References

1. Mealy, G. H. (1955). "A method for synthesizing sequential circuits."
   *The Bell System Technical Journal*, 34(5), 1045–1079.
   doi:[10.1002/j.1538-7305.1955.tb03788.x](https://doi.org/10.1002/j.1538-7305.1955.tb03788.x)
2. Moore, E. F. (1956). "Gedanken-experiments on sequential machines." In
   C. E. Shannon & J. McCarthy (Eds.), *Automata Studies*, Annals of
   Mathematics Studies 34 (pp. 129–153). Princeton University Press.
   Reprinted De Gruyter, doi:10.1515/9781400882618-006.
3. Sipser, M. (2013). *Introduction to the Theory of Computation* (3rd
   ed.). Cengage Learning. Definition 1.5, p. 35; p. 36; Definition 1.37,
   p. 53.
4. Hopcroft, J. E., Motwani, R., & Ullman, J. D. (2007). *Introduction to
   Automata Theory, Languages, and Computation* (3rd ed.).
   Pearson/Addison-Wesley. §2.2, p. 45 (determinism); §2.2.3, p. 48
   (transition tables); p. 44 (inaccessible states); p. 67 (dead/trap
   states).
5. Dijkstra, E. W. (1975). "Guarded commands, nondeterminacy and formal
   derivation of programs." *Communications of the ACM*, 18(8), 453–457.
   doi:[10.1145/360933.360975](https://doi.org/10.1145/360933.360975)
6. Harel, D. (1987). "Statecharts: a visual formalism for complex
   systems." *Science of Computer Programming*, 8(3), 231–274.
   doi:[10.1016/0167-6423(87)90035-9](https://doi.org/10.1016/0167-6423(87)90035-9)
7. Vanthienen, J., & Dries, E. (1992). *Developments in Decision Tables:
   Evolution, Applications and a Proposed Standard.* Onderzoeksrapport
   9227, Katholieke Universiteit Leuven.
8. Kavanagh, T. F. (1960). "TABSOL: a fundamental concept for
   systems-oriented languages." In *Proceedings of the Eastern Joint
   IRE-AIEE-ACM Computer Conference*, 117–136.
   doi:[10.1145/1460512.1460522](https://doi.org/10.1145/1460512.1460522)
9. IBM (1962). *Decision Tables: A Systems Analysis and Documentation
   Technique.* IBM General Information Manual, form F20-8102.
10. CODASYL (1982). *A Modern Appraisal of Decision Tables: A CODASYL
    Report.* Report of the Decision Table Task Group. ACM. (The task group
    was initiated in 1973 and reported in 1982; the report is often
    misdated to the 1960s.)
11. Jorgensen, P. C. (2014). *Software Testing: A Craftsman's Approach*
    (4th ed.), Ch. 7 "Decision Table-Based Testing", pp. 117–131. CRC
    Press.
12. Reinwald, L. T., & Soland, R. M. (1966). "Conversion of limited-entry
    decision tables to optimal computer programs I: minimum average
    processing time." *Journal of the ACM*, 13(3), 339–358. (Part II:
    *JACM* 14(4), 1967, 742–755.)
13. King, P. J. H. (1968). "Ambiguity in limited entry decision tables."
    *Communications of the ACM*, 11(10), 680–684.
    doi:[10.1145/364096.364113](https://doi.org/10.1145/364096.364113)
14. Object Management Group (2024). *Decision Model and Notation (DMN),
    Version 1.5.* OMG document formal/24-01-01. §8.1, §8.2.4, §8.2.9,
    §8.2.11. <https://www.omg.org/spec/DMN/1.5/>

### Prior work in this project

The design research behind intrastate lives in the sibling
`state-machines` repository, whose `research/state-machines-research.md`
carries a reference list for the transition-system side — Harel, the
actor and CSP lineage, state-machine replication, and the
agent-as-state-machine literature.

That list has **no decision-table entries**. The design research was
state-machine-first, and the decision-table class arrived later without a
literature pass. The decision-table references above are therefore new
work rather than a restatement of that survey, which is also why the
FSM-versus-decision-table contrast below could not be sourced from it.

### A note on what is *not* cited

No source here contrasts decision tables with finite state machines
directly. The literatures developed separately — decision tables in
commercial data processing, state machines in automata theory and circuit
design — and the closest authoritative statement is DMN's, that decision
modeling "complements process modeling," which is a contrast with BPMN
process models rather than with FSMs. The framing in this document that
the two classes answer different questions is therefore intrastate's own,
not a claim borrowed from the literature.
