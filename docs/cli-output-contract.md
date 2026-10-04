# CLI output contract

Worked payloads and the reasoning behind the envelope shape.

**The binary states the contract itself.** `intrastate --help-all` lists
the output modes, the refusal codes each verb can return, and the
exit-code contract, spelled through the same constants the wire is
emitted from. This document does not repeat that reference — it carries
what help output cannot: real payloads to parse against, and the *why*
behind the shapes.

For the mode/stream summary and the exit table, run:

```sh
intrastate --help-all      # vocabulary, output modes, exit codes
intrastate lint --help-all # the live finding taxonomy
```

## Vocabulary stability tiers

Every machine-readable vocabulary this CLI emits carries exactly one
declared stability tier, recorded here beside that vocabulary.

- **`frozen`** — no member added or removed within a major. A consumer MAY
  treat an unrecognized member as a defect.
- **`append-only`** — members MAY be added in a minor; none removed or
  renamed within a major. A consumer MUST tolerate an unrecognized member,
  and MUST NOT assert on the set's cardinality, a member's ordinal
  position, or a tail position.
- **`growing`** — as `append-only`, and additionally a new member MAY fire
  on input that previously produced no such finding. A consumer that counts
  findings, gates on an empty findings list, or diffs output across releases
  MUST expect movement on identical input. The verdict still cannot change;
  the payload can.

**While the binary's version is `0.x`, a tier is a DECLARATION OF INTENT.**
It states what the surface is expected to promise at 1.0.0 and MUST NOT be
read as a guarantee already in force. That qualification is carried by the
schema major, not by this prose: while `schema_version` reads `0.x` the
tiers are intent, and from `"1.0"` they are guarantees. The tier names and
assignments are authored and maintained from the first release, so reaching
1.0.0 requires no reclassification.

| vocabulary | tier | seam |
| --- | --- | --- |
| envelope `type` discriminator | frozen | `respond::Types()` |
| `schema_version` | frozen (the field; its VALUE moves per the rules below) | `clierr.SchemaVersion` |
| severity (`blocking`, `info`) | frozen | `graphlint::Severities()` |
| exit-code classes | frozen | `clierr::ExitCodes()` |
| stderr advisory `level` (`note`, `warning`) | frozen | `respond::Levels()` |
| gate `verdict` (`allow`, `deny`, `indeterminate`) | frozen | `accessor::Verdicts()` |
| `data.escape_class` | frozen | `resolve::RefusalKinds()` |
| `findings[].operator` | frozen | `guard::Operators()`, `table::Operators()` |
| the `graph-lint-failed` aggregate code | frozen | single `const` |
| the `version` payload field names (`version`, `commit`, `date`) | frozen | `version::Info` |
| model load categories | append-only | `table::Categories()` |
| the CLIError `code` vocabulary | append-only | none (prose-only) |
| the `graph-unprovable-coverage` `reason` set | append-only | `graphlint::Reasons()` |
| the `flow next` unknown-`reason` set | append-only | `cli::UnknownReasons()` |
| graph-lint BLOCKING finding codes | append-only | `graphlint::BlockingCodes()` |
| `findings[].block` | append-only | `resolve::Blocks()` |
| graph-lint ADVISORY finding codes | growing | `graphlint::AdvisoryCodes()` |

Two emitted namespaces take NO tier, deliberately: `findings[].class` and
`data.dispositions`. Both carry model-authored tokens to which this CLI
ascribes no meaning — author surface, not CLI vocabulary.

### `schema_version` and how the wire moves

Both terminal records — the `ok` envelope and the bare `CLIError` refusal —
carry a top-level `schema_version` string of the form `MAJOR.MINOR`. It is
never `omitempty`, is versioned independently of the binary's release
version, and is never projected into `data`.

It begins at `"0.1"` and tracks the wire, not the binary. While its major is
`0` the schema is explicitly unstable and MAY change incompatibly in any
release; **the minor component still increments on every change**, so a
consumer can detect movement even while it cannot rely on compatibility.

From `"1.0"` onward the minor increments for backward-compatible additions —
a new optional field, a new member of an append-only or growing vocabulary —
and the major increments for changes that are not backward-compatible: a
removed or renamed field, a removed or renamed member of any vocabulary, or
a change to a frozen surface.

A consumer MUST ignore object properties with unrecognized names, and MUST
reject an envelope reporting an unsupported major.

### Introducing and promoting finding codes

A lint finding code is **introduced at severity `info`**. Introduction is a
minor-release change and MUST NOT alter the success disposition of any
input.

**Promotion** of a finding code from `info` to `blocking` is a distinct
release event. It MUST be disclosed in the release notes for the release
that carries it, naming the code, and MUST NOT occur in a patch release.
The disclosure obligation binds during `0.x` as well as after.

A code MAY be introduced directly at `blocking` only when it reports a
condition that was **already refused** by some other code — that is, when it
re-attributes an existing refusal rather than creating a new one. The
exception is discharged by evidence, not by judgement: the change must name
the pre-existing code it re-attributes from and show that the input refused
by the new code was already refused — same input, same verdict, different
code.

## Structured findings

`findings` is a top-level array on the `CLIError` envelope — a **sibling**
of `code`, not nested under an `error` wrapper. It is `omitempty`, so a
refusal that carries none omits the key entirely.

The carrier is fixed by the failure FAMILY, not by how many subjects a
given run happens to produce:

- **Scalar families** — tag, write, and artifact validation, and model
  selection — name their one offending subject in `param`.
- **Aggregate families** — gate results, model-load categories,
  read-back mismatches, and declared-line-edit refusals — report in
  `findings`, one entry per subject, **regardless of runtime
  cardinality**. A single denied gate and a model-load error naming one
  category each emit a one-element `findings` array and no top-level
  `param`.

A declared-line-edit refusal is an aggregate family for a reason of
CARDINALITY: one write entry can carry several line rules, so one refusal
can have several subjects, and none of them is a flag or an argument.
`findings` is therefore the carrier even when a given run produces one.
Each finding's `message` is the refusal's `Detail` verbatim. These refusals
are decided BEFORE any byte is written and take the **exit-2** group: they
are about the request, not the environment, so re-running the same request
unchanged cannot help. An unreadable target is the contrasting case and
keeps exit 3.

**Do not assume the subject is a rule.** It often is, but three shapes
ship, and a consumer splitting a `message` on the rule-scoped form
mis-slices two of them:

- **rule-scoped** — `<token>: <accessor id>.edit.<key>: <detail>`. A
  cardinality or stability defect belongs to one rule and names it, which
  is what lets a caller trace the refusal to one declaration among an
  entry's siblings.
- **entry-scoped** — `<token>: <accessor id>: <detail>`. The multiline-value
  scan spans the entry's whole plan and the cross-rule collision sweep is a
  property of a PAIR, so neither has a single rule to name.
- **pair-scoped** — an entry-scoped subject whose `<detail>` names the two
  colliding rules, `edit_anchor_collision` being the one case with two
  rules to report and no way to attribute the defect to either alone.

A fourth case carries **no** `edit_*` reason token at all. The
**entry-level preconditions** name the gate or the placeholder instead,
having no rule to name. An unbound `{tag.<key>}` in a rule's own *anchor*
is still a declared line edit, so it takes this code and this `findings`
array with the placeholder — not a token, not a rule — as its subject. The
other entry-level preconditions are refusals of the *request* rather than
of a line edit: the gate-off read-back check and the argv preconditions of
any command-backed accessor take the generic request-refused code, carry no
`findings` array at all, and name the accessor in `param`. The refusal's
ORIGIN selects the envelope; the phase it was raised in does not.

The declared-line-edit carrier names two DISJOINT sets of six. Confusing
them is the mistake this paragraph exists to prevent: only the first set
registers in the load-category list, and a consumer matching an apply-time
token against that list will never find it.

Load-time categories, reported by `lint` when a model is read:
`edit_carrier_conflict`, `edit_key_mismatch`, `edit_anchor_invalid`,
`edit_template_invalid`, `edit_clear_invalid`, `edit_tag_argv0`. The last
fires on a `command` entry of any kind, including one carrying no `edit`
table, because it belongs to the argv placeholder family rather than to
the carrier.

Apply-time reason tokens, carried on a refusal's `Detail` beside its
subject — rule-scoped or entry-scoped, per the shapes above — and
registering in no category list: `edit_anchor_unmatched`,
`edit_anchor_ambiguous`, `edit_anchor_collision`, `edit_anchor_unstable`,
`edit_value_multiline`, `edit_clear_undeclared`. Lint does not decide any
of them — a stale anchor is a property of the artifact at apply time, not
of the model.

A consumer therefore selects the carrier from the refusal `code`, never
from the number of subjects it observes.

A scalar family names the offending flag, key, role, outcome, or accessor
id in `param`. An aggregate family uses `findings[]` because reporting only
the first subject would hide the rest — and for an ambiguous match it would
amount to a tie-break the model did not author.

Each finding is one flat record. Optional fields are `omitempty`, so a
field that is present is one its producer populated:

| field | meaning |
| --- | --- |
| `code` | the finding's own discriminator — for a model-load failure, the load category slug |
| `message` | self-sufficient prose; it renders the failure readably with no other field consulted |
| `param` | the offending flag or argument |
| `locator` | source position, as `file:line` |
| `hint` | a one-line remedy |
| `severity`, `model`, `rule`, `span`, `element`, `reason`, `dimension` | graph-lint attribution |
| `fingerprint` | the canonical sortable predicate/write serialization used in finding identity |
| `key`, `operator`, `literal`, `block` | a guard atom's fields, flat — never nested under an `atom` object |
| `class` | a declared failure class |
| `count` | how many underlying subjects collapsed into this one finding, when the producer collapses equal identities — the record's one number; absent when it is zero |

No producer nests its own fields in a sub-object, and each populates only
the fields it owns.

On a per-row graph-lint finding, `rule` (and `element`, where the finding
names a second row) carries the **row's identity**: the authored rule id,
joined with the row's expansion suffix by `#` when the rule expanded — so a
row a step or an `in` minted reads `retry#3`, not `retry`. A finding that
names a group rather than a row keeps the bare authored id. A consumer
joining findings to source should use `span`, the stable key, which always
names the authored rule.

In `--as=text` every finding renders on its own line under the message,
with its identity fields appended. The set of finding codes in text output
equals the set in JSON output: a renderer that summarised would drop a
defect the caller needs to see.

## Set values on the wire

A set-valued tag crosses the CLI as its canonical JSON array — members
sorted, duplicate-free, compact — rendered with HTML escaping **disabled**,
so `<`, `>`, and `&` serialize as themselves and never as `\u003c`,
`\u003e`, or `\u0026`.

Every site that emits or compares one builds it through the same
canonicalization — sorting, deduplicating, and encoding with HTML escaping
disabled — which is what makes plan-to-request copy-through and read-back
equality byte equality. `clierr.WriteJSONLine` is the non-HTML-escaping
JSON line writer the envelope itself is emitted through; it does not sort
or deduplicate, so it is not by itself the canonical-set encoder.

## Why exit 2 and exit 3 are separate

`intrastate --help-all` carries the exit table; the mapping lives in
`clierr.ExitCodeFor`. What matters here is the reasoning, because it
constrains every future refusal:

**Exit 3 means the environment could not be consulted, and the same
request may be re-run unchanged.** An accessor timed out, failed to
execute, returned an incomplete key set, or a post-mutation read-back did
not complete. The remedy is to repair the environment and re-issue the
identical request.

Every other failure is exit 2 — the request or the model is wrong, or the
model said no. The split is what makes a caller's retry loop safe: a
refusal about the request must never exit 3, or the caller spins on an
input it must instead fix. A command-backed accessor invoked without the
`--allow-commands` opt-in is one: it refuses before spawn as
`flow-accessor-request-refused`, exit 2.

A read-back that did not complete carries a `detail` saying the write may
have been applied and was not verified. It is never reported as a write
that did not occur.

A write whose command already ran and then failed carries that same
`detail` under its own code, `flow-write-failed-applied`. It is exit 3 like
its neighbours, and it is deliberately distinct from `flow-accessor-failed`:
a write that could not start may be retried unchanged, and one that may
already have landed must be inspected first.

A `steps` writer runs one command per step, and the split falls per step.
A first step that fails — a non-zero exit, a signal, or a failure to start
— applied nothing, even though its command ran, so it is
`flow-accessor-failed`: a compare-and-set placed first refuses cleanly when
it loses. A step that fails after an earlier step exited 0 is
`flow-write-failed-applied`, exit 3: the earlier steps' mutations stand,
and nothing is retried or undone. A held pipe keeps the applied sense on
any step, the first included. A removal the model declares no `clear` arm
for refuses before any step runs, as `flow-accessor-request-refused` at
exit 2, since re-running it unchanged cannot help.

## Held output pipes on a command accessor

A command accessor's stdout and stderr are drained under a bound. If a
process still holds one of those pipes open when the bound expires — the
usual cause being a helper that backgrounds a daemon inheriting its stdio —
the invocation **refuses** rather than waiting. The refusal names the pipe
or pipes that were held, the child's own exit status, and the remedy:
*close or redirect the helper's inherited stdio*. Redirecting the
background process's stdout and stderr (to a file, or to `/dev/null`) is
the fix; there is no flag that accepts the partial output instead.

Output read from a held pipe is never parsed. No read envelope, no gate
verdict and no write read-back is derived from a stream whose writer the
CLI could not reach, so a held pipe is always a refusal and never a
half-trusted answer.

Two consequences are admitted rather than hidden. A process outside the
child's process group may **leak** — it outlives the refusal, and whatever
it writes after that is never read; it stays visible to `ps` until its next
write fails. And when the accessor's declared `timeout` had already elapsed,
the refusal is classified `timeout`, which by contract carries no `detail`
— so the held-pipe reason is reported on `flow-accessor-failed` and not on
the timeout. The bounded refusal itself is never withheld either way.

## Illustrative invocations

The `flow` command group is the skill-integration surface. These shapes
exercise the flag grammar; the output shown by `--as=json` is one terminal
envelope per invocation.

```sh
# Report the candidate rules the supplied state can take, and the legal
# outcome alphabet behind them.
intrastate flow next --model flow.toml \
  --artifact state=./state.json --tag profile=mid --as=json

# Report every row the guards do not exclude, regardless of match.
intrastate flow next --model flow.toml \
  --artifact state=./state.json --all --as=json

# Map one recognized outcome to exactly one plan, or refuse.
intrastate flow resolve --model flow.toml \
  --artifact state=./state.json --outcome advance --as=json

# The same verb over a decision-table model. The table declares no owned
# tag, so no reader is invoked and --artifact is not required; the
# discriminating dimensions arrive as --tag. The answer is `rule` plus
# `emit`.
intrastate flow resolve --model pricing.toml \
  --outcome decide --tag tier=free --tag region=eu --as=json

# Report what the declared read accessors see. Runs EVERY declared reader,
# so every declared role must be bound.
intrastate flow read-state --model flow.toml \
  --artifact state=./state.json --artifact notes=./notes.json --as=json

# Apply planned owned-tag writes, verified by read-back. A set value is a
# JSON array literal; a removal uses --clear, never --write k=<clear>.
intrastate flow set-state --model flow.toml \
  --artifact state=./state.json \
  --write status=final --write 'labels=["cli","final"]' \
  --clear stale --as=json

# Seed owned state from the model's [initial] root, once. It carries no
# write grammar: the model's own root is the plan.
intrastate flow init-state --model flow.toml \
  --artifact state=./state.json --as=json
```

## `flow resolve`, the `emit` answer, and its `dispositions`

The `flow resolve` success payload carries `emit`, a JSON object of string
values with keys in byte order. It is the row's authored answer, joined to
the selected rule after selection. It is present as `{}` — never `null`,
never omitted — when the selected row authored none, so a consumer parses
one shape either way. `emit` sits immediately after `gates`: the row
selected and what it says are one answer.

```json
{"model":"pricing.toml","revision":"","observed":{"region":"eu","tier":"free"},
 "owned":{},"readers":[],"outcome":"decide","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},"dispositions":{"dpa":"gate"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

Gates are unchanged and run before the answer is built, so a denied gate is
`flow-gate-denied` — a refusal, never a payload — and `emit` is never
computed on a denied selection.

A plan rescued by an escape row carries **that escape row's own** `emit`,
alongside `escaped: true` and its `escape_class`. An escape row authoring no
answer renders `{}` like any other row.

In `--as=text` the field renders through the same generic payload renderer as
every other: one path-qualified leaf per line, `emit.<key>: <value>`, and an
unauthored block as `emit: (none)`.

`flow next` carries no `emit`. Its candidate preview reports what a row would
require and write without evaluating anything; the answer is what `flow
resolve` selects.

### Declaring the emit vocabulary

A model MAY declare its emit vocabulary in a top-level `[emit]` table, one
sub-table per emit key. Declaring is opt-in and the trigger is the COUNT of
declared keys, never the presence of the table: a model with zero
declarations behaves exactly as it does today, and a bare `[emit]` table
with no sub-tables is a zero-declaration model.

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
not admitted, because an emit value is one authored string. Only an `enum`
takes a `domain`; `bool` fixes the implicit domain `true | false`, `int`
constrains the value to a base-10 integer literal, and `scalar` is the
declared-but-unvalidated escape hatch — the key is admitted, the value
unconstrained.

The `domain` key has two spellings, and they are two spellings of ONE key,
not two fields. `domain = [...]` is a flat member array with no
dispositions. `[emit.<key>.domain]` is a sub-table whose keys are
MODEL-AUTHORED disposition tokens and whose values are member arrays; the
key's domain is the union, and each member carries the one disposition it
is listed under.

Every check is LEXICAL, on the authored string. No value is parsed into a
typed representation, canonicalized, or converted anywhere: the payload and
the dump carry the authored bytes. The authored value is always a TOML
string, so under `kind = "int"` an author writes `count = "42"` and a bare
`count = 42` is a decoder refusal, `malformed_toml`, exactly as it is today.

Once one key is declared, the whole model is strict, and three load
categories become reachable — refusals, never advisories, so `intrastate
lint` exits nonzero and the `flow` verbs answer `flow-model-invalid`:

| category | when |
| --- | --- |
| `malformed_emit_declaration` | the declaration itself is ill-formed — an unknown `kind`, an `enum` with no usable domain, an empty or duplicated member, a `domain` on a non-enum kind, an empty disposition token |
| `unknown_emit_key` | a rule, ordinary or escape, emits a key `[emit]` does not declare |
| `emit_value_out_of_domain` | an authored value is outside its key's declared domain, is not a `bool` token, or is not an `int` literal |

Each carries the offending block's source line in its `locator`. Load is
fail-fast, so exactly one refusal comes back per run.

A declaration is author-owned and unversioned: widen or narrow a domain by
editing the model. A declared key no rule emits, and a declared member no
rule authors, are findings of no tier — that is authoring headroom.

### `dispositions`

The success payload carries `dispositions` immediately after `emit`: a JSON
object mapping emit key → the disposition token the declaration assigns the
SELECTED ROW's authored value.

An entry exists for key `k` exactly when the selected row authors `k` AND
`k`'s declaration lists that authored value under a disposition. A declared
key the row does not emit contributes no entry, and the object is never
padded to the declared key set. It is present as `{}` — never `null`, never
omitted — when no selected value carries one: an undeclared model, a
non-enum kind, a flat-array domain and an empty emit block all render `{}`.
Keys arrive in byte order.

`models/examples/routing-decision-table.toml` is the worked example. It
declares one key, `next`, partitioned into `route` and `stop`:

```toml
[emit.next]
kind = "enum"
[emit.next.domain]
route = ["escalate", "notify"]
stop = ["park", "close"]
```

Its `low-assigned` row answers `next = "park"`, which the declaration lists
under `stop`:

```json
{"rule":"low-assigned","gates":[],
 "emit":{"next":"park"},"dispositions":{"next":"stop"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

The token is surfaced verbatim. intrastate fixes no disposition vocabulary
and never interprets a token — `route` and `stop` above are the model
author's words, and the tool only guarantees they arrive unchanged.

A plan rescued by an escape row joins `dispositions` from that escape row's
own authored values, the same single join path `emit` takes. In `--as=text`
the field renders through the generic payload renderer as
`dispositions.<key>: <token>`, and as `dispositions: (none)` when empty.
`flow next` carries no `dispositions`, for the same reason it carries no
`emit`.

## The `graph` document and what it does not claim

`intrastate graph` exports a model's normalized graph as a document. It is
a DOCUMENTATION artifact, never a verdict: the command runs load,
normalization, grouping, and the reachability traversal only. It never runs
the graph invariants, emits no findings, and exports any model that loads —
including one `intrastate lint` refuses, which is what keeps a failing
model inspectable as a graph. `lint` remains the acceptance surface.

The JSON document is versioned by its required LEADING `schema` field,
initial value `intrastate.graph/1`. That marker versions the DOCUMENT.
The envelope's own `schema_version` versions the ENVELOPE carrying it and
is never projected into `data`, where it would collide with this field:
two markers at different levels, neither substituting for the other.

`--emit` selects the DOCUMENT and `--as` selects the ENVELOPE, and all
four combinations are defined:

```sh
# The bare JSON document on stdout, plus one trailing newline.
intrastate graph --model flow.toml

# The same relation as a DOT digraph, ready to pipe.
intrastate graph --model flow.toml --emit dot | dot -Tsvg > graph.svg

# One terminal ok envelope whose `data` embeds the same document.
intrastate graph --model flow.toml --as=json | jq .data

# Under --emit dot the `data` object carries the DOT text as its single
# required string member, so the unwrap is exactly:
intrastate graph --model flow.toml --emit dot --as=json | jq -r .data.dot
```

Under `--as=text` stdout carries the selected document verbatim and
nothing else, so the DOT stream pipes to `dot` and the JSON document diffs
raw in CI. A consumer that parses stdout as an envelope MUST use
`--as=json`: the text-mode stream is the bare document by contract.

Evolution within `/1` is ADDITIVE. A member MAY be added in a minor; none
is removed or renamed within a major. A consumer MUST therefore tolerate
an unrecognized field and MUST NOT assert on the field set's cardinality,
on a member's ordinal position, or on a tail position.

For one model and one build, emission is byte-for-byte identical across
invocations in every `--emit` and `--as` combination. The promise is scoped
to (model, build) and NOT across builds: a build bump may change the JSON
only by the additive rule above, and may re-baseline DOT diffs freely,
since DOT styling is explicitly non-normative. Identity is
invocation-independent — `model` carries the authored `[model] id`, never
the `--model <path>` argument or any path-derived string — so the same
model exported from two checkouts is byte-identical.

### The `abstraction` marker and the soundness rule

The `reach` block carries a REQUIRED `abstraction` marker with the token
value `declared-over-approximation`. It states what the relation IS: the
DECLARED over-approximation of the runtime, with nodes merged and
guard/observed atoms left unpruned.

**The soundness rule: universal claims ("no path does X") proved over this
relation hold at runtime; existence claims ("some path reaches X") may be
spurious.** A reader who takes the merged relation for the runtime one
over-counts edges and will "verify" a property the runtime lacks. The
marker is required precisely so that no consumer reads the diagram as a
certificate.

The document carries NO verdict field and NO finding field — an export is
never a lint pass — and no per-node terminal marking: which merged nodes
satisfy a `terminal` predicate set is a quantifier this document does not
decide. The declared `initial` and `terminal` sets travel in the document
for a consumer to evaluate for itself.

This record claims no content-addressed identity and no replay-stable
hash. The promise is byte identity of the emitted stream, not a digest, so
the document carries no hash, digest, or checksum member.

When the traversal does not complete under the published node ceiling the
command refuses with `graph-export-too-large` (exit 2), naming the ceiling
and the narrow-a-domain remedy. It never emits a partial document: a
partial graph diffs as a graph change.

## `flow resolve --plan-only` and the two halves of the payload

Every field of the `flow resolve` success payload belongs to exactly one of
two groups, and the assignment is fixed:

| group | fields | what it is |
| --- | --- | --- |
| **echo** | `model`, `observed`, `owned`, `readers`, `outcome` | the request, read back — the caller already holds all of it |
| **plan** | `revision`, `rule`, `gates`, `emit`, `dispositions`, `next`, `writes`, `clear`, `escaped`, `escape_class` | what the call DECIDED |

`--plan-only` omits the echo group. It is report-only: the same rule is
selected, the same gates run, the same readers are invoked, and refusals are
byte-identical with and without it. Every field it does carry is carried
byte-for-byte, in the same order, through the same encoder — the projection
deletes whole keys and never rewrites a value.

An omitted key is **absent**. It is never `null`, never `{}`, and never
`""`: those would be stand-ins, and a consumer could not tell a projected
run from one whose reader returned nothing.

The full payload:

```json
{"model":"pricing.toml","revision":"","observed":{"region":"eu","tier":"free"},
 "owned":{},"readers":[],"outcome":"decide","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},"dispositions":{"dpa":"gate"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

The same call with `--plan-only`:

```json
{"revision":"","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},"dispositions":{"dpa":"gate"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

Fields with a presence rule keep it, unchanged, inside the plan group:
`emit` and `dispositions` are still present as `{}` when nothing joins, and
`escape_class` still appears exactly when it would appear by default —
omitted on an unescaped plan in both widths. The flag neither widens nor
narrows a producer's own rule.

`--plan-only` rides `flow resolve` alone. On `next`, `read-state` or
`set-state` it fails the parse as `command-error`, exit 2 — which is also
what a binary predating the flag does, so a caller can never hold a
full-width payload believing it was projected.

The default width is unchanged and remains the full record; dropping the
flag is the whole recovery procedure.

## `flow set-state`, and how a no-op is told from an applied write

The `flow set-state` success payload carries `writers`, the accessors that
actually applied keys, alongside the `writes` and `clear` it was asked for.
A plan that plans nothing — empty `writes` and an absent or empty `clear`,
which is what a `stopped:*` or `none` row resolves to — is a **no-op success
at exit 0**, not a refusal. `writers` and `writes` emptiness is what tells
the two apart, and it is the sanctioned check:

```json
{"model":"flow.toml","revision":"","artifacts":{"state":"./state.json"},
 "writers":[],"writes":{},"clear":[],"owned":{}}
```

`writers` is `[]` exactly when nothing was applied and non-empty on any real
write, because it is appended to only as each writer's read-back confirms.
It cannot be `[]` on a write that was meant to land: a key no declared write
accessor serves refuses with `write-unbound` rather than applying nothing,
so an empty `writers` never stands in for a silently dropped key.

`dispositions` is deliberately **absent** from this payload. A disposition
is joined from the SELECTED ROW's authored `emit`, and
`set-state` selects no row — echoing a caller-supplied plan's blob back
under that name would give one wire key two provenances. `flow next` omits
it for the same reason. A caller that wants the disposition reads it off the
`flow resolve` envelope it piped in, where it is joined to a real row; a
caller that wants to know whether the write landed reads `writers`.

## `flow next` candidates and the `unknown` list

A candidate is a rule whose match and guard both HOLD or are UNDECIDED
over the supplied state. A match or guard key the state does not carry
does not exclude the row — it leaves the row a candidate with that key
named under `unknown`. A candidate is what the state does not exclude, not
what `flow resolve` will select: a candidate carrying an `unknown` entry
may still be refused by `flow resolve` over the same state.

`--all` reports every row the guards do not exclude, regardless of match.
In that mode a match atom takes no part in the verdict and contributes no
`unknown` entry.

Each candidate carries `rule`, `outcome`, `required`, `unknown`, `next`,
`writes`, `clear`, and — only under `--evaluate-gates` — `gates`. The
`unknown` list is present in both modes, as `[]` rather than omitted, so
one consumer struct parses either.

Each `unknown` entry is a `{key, reason}` pair, deduplicated on the pair
and sorted by `(key, reason)`. The reason names the remedy and comes from
a closed set:

| reason | source | remedy |
| --- | --- | --- |
| `absent` | a match or guard atom over a key the state does not carry, or an owned key no invoked reader established | bind the reader, or supply the tag |
| `uncomparable` | a guard atom over a key present at a value its operator cannot compare | fix the value, or the atom's literal |
| `not-evaluated` | one of the row's own gate ids, un-run because `--evaluate-gates` was not passed | pass `--evaluate-gates` |

```json
{"rule":"prelock","outcome":"advance","required":["stage"],
 "unknown":[{"key":"approval","reason":"not-evaluated"},
            {"key":"gate_passed","reason":"absent"}],
 "next":{"stage":"prelocked"},"writes":{"stage":"prelocked"},"clear":[]}
```

## `flow init-state`, and how a seed is told from a no-op

`flow init-state` writes the model's `[initial]` assignments through the
declared write accessors, verified by read-back — the same route `set-state`
takes. It has no write grammar of its own, so nothing a caller types can
become a seeded value.

Seeding is ALL-OR-NOTHING over an EMPTY store. It seeds if and only if every
bound artifact carries no key, and then it seeds every `[initial]` key. The
quantifier is over the whole bound set: one surviving key ANYWHERE blocks the
whole seed, so a cleared key is never re-established while the store retains
another. A store that carries any key — torn, partially seeded, post-clear, or
fully seeded alike — is a no-op SUCCESS at exit 0 with zero writes.

The success payload carries `seeded` and `absent`, and which arm ran is read
off them:

| field | meaning |
| --- | --- |
| `seeded` | the `[initial]` keys this run wrote. Non-empty on exactly the seeding arm. Key NAMES only — `flow read-state` is the surface that reports values. |
| `absent` | the `[initial]` keys the store does not carry. It belongs to the no-op arm; on the seeded arm every key was just written, so it is necessarily empty. |

Both fields are present as `[]` rather than omitted, so one consumer struct
parses either arm.

```json
{"model":"flow.toml","revision":"","artifacts":{"state":"./state.json"},
 "seeded":["note","stage"],"absent":[]}
```

The scope is the FILE-BACKED carrier. An edit-carried or command-backed
accessor has no JSON store, so "the store carries no key" is undefined for
it, and a model whose needed write OR read binding is not file-backed refuses
at exit 2 naming the capability and the accessor. That refusal is decided at
registry construction, before any accessor is invoked, so it preempts the
`--allow-commands` refusal: a command-backed accessor yields the carrier code
whether or not the opt-in was passed.

Two refusal classes are this verb's own, both exit 2 and both distinct from
each other and from the shared classes:

| condition | disposition |
| --- | --- |
| the model's class declares no `[initial]` root to materialize (a decision table) | exit 2, before any accessor runs and with no artifact touched |
| a needed role's write or read binding is not the file-backed carrier | exit 2, naming which capability and which accessor failed, with zero writes |

Everything else reuses the shared `flow-*` codes unchanged, at the groups
those codes already carry: an unbound needed role is `flow-artifact-missing`,
a model routing one `[initial]` key to two writers refuses at LOAD as
`flow-model-invalid`, and the read-back splits the way it does for
`set-state` — completed-and-disagreed is exit 2, could-not-complete is exit 3.
Both leave a NON-EMPTY store, so a re-run is a no-op that does NOT repair
either; recovery is an explicit `set-state`, or discarding the artifact and
re-running init.

One boundary is disclosed rather than designed around: clearing the LAST
remaining key empties the store, and an emptied store is indistinguishable at
the content level from a never-written one — a clear is a key REMOVAL, not a
tombstone. A subsequent `init-state` therefore RESEEDS such a store.
