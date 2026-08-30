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

## Structured findings

`findings` is a top-level array on the `CLIError` envelope — a **sibling**
of `code`, not nested under an `error` wrapper. It is `omitempty`, so a
refusal that carries none omits the key entirely.

The carrier is fixed by the failure FAMILY, not by how many subjects a
given run happens to produce:

- **Scalar families** — tag, write, and artifact validation, and model
  selection — name their one offending subject in `param`.
- **Aggregate families** — gate results, model-load categories, and
  read-back mismatches — report in `findings`, one entry per subject,
  **regardless of runtime cardinality**. A single denied gate and a
  model-load error naming one category each emit a one-element `findings`
  array and no top-level `param`.

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
input it must instead fix.

A read-back that did not complete carries a `detail` saying the write may
have been applied and was not verified. It is never reported as a write
that did not occur.

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
```

## `flow resolve` and the `emit` answer

The `flow resolve` success payload carries `emit`, a JSON object of string
values with keys in byte order. It is the row's authored answer, joined to
the selected rule after selection. It is present as `{}` — never `null`,
never omitted — when the selected row authored none, so a consumer parses
one shape either way. `emit` sits immediately after `gates`: the row
selected and what it says are one answer.

```json
{"model":"pricing.toml","revision":"","observed":{"region":"eu","tier":"free"},
 "owned":{},"readers":[],"outcome":"decide","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},
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

## `flow resolve --plan-only` and the two halves of the payload

Every field of the `flow resolve` success payload belongs to exactly one of
two groups, and the assignment is fixed:

| group | fields | what it is |
| --- | --- | --- |
| **echo** | `model`, `observed`, `owned`, `readers`, `outcome` | the request, read back — the caller already holds all of it |
| **plan** | `revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`, `escape_class` | what the call DECIDED |

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
 "emit":{"dpa":"required","plan":"basic"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

The same call with `--plan-only`:

```json
{"revision":"","rule":"free-eu","gates":[],
 "emit":{"dpa":"required","plan":"basic"},
 "next":{},"writes":{},"clear":[],"escaped":false}
```

Fields with a presence rule keep it, unchanged, inside the plan group:
`emit` is still present as `{}` when the row authored none, and
`escape_class` still appears exactly when it would appear by default —
omitted on an unescaped plan in both widths. The flag neither widens nor
narrows a producer's own rule.

`--plan-only` rides `flow resolve` alone. On `next`, `read-state` or
`set-state` it fails the parse as `command-error`, exit 2 — which is also
what a binary predating the flag does, so a caller can never hold a
full-width payload believing it was projected.

The default width is unchanged and remains the full record; dropping the
flag is the whole recovery procedure.

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
