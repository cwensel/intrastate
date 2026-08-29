Model: claude-opus-5[1m]

# A8 spike: strict text-mode byte width

Date: 2026-08-29

A8 claims the strict-width clause holds in TEXT mode on the same unit as
JSON: the projected `--as=text` rendering is strictly fewer TOTAL BYTES
than the default rendering of the same request, not merely a line subset.

`--plan-only` is not implemented; the projected rendering is derived from
the shipped default rendering by the filter below.

## Filter rule

`internal/cli/respond/text.go::flatten` decodes the success payload's
canonical JSON and emits one path-qualified leaf line per value, with map
keys sorted (`sort.Strings`) at every level. Path segments join with `.`
for map keys and `[i]` for array elements; `label` appends `": "` after
the path. A leaf at top level therefore renders as `key: value`, a nested
map leaf as `key.sub: value`, an array element as `key[0]: value`, and an
empty container or null as `key: (none)`.

Dropping the five ECHO keys (`model`, `observed`, `owned`, `readers`,
`outcome`) from the payload map therefore deletes exactly the lines whose
path root is one of those keys — i.e. lines matching:

```
^(model|outcome)(: )|^(observed|owned|readers)([.[]|: )
```

The alternation on `.`, `[`, and `: ` is what makes the match root-exact:
it matches `observed.tier: paid`, `readers[0]: release-state`, and the
whole-container form `owned: (none)`, while never matching a PLAN key that
merely shares a prefix. No retained line is rewritten, because sorted key
order among the retained keys is unchanged by deleting others and each
line's path is independent of its siblings.

Byte counts are `wc -c` over the full rendering INCLUDING the trailing
newline of every line (the renderer uses `fmt.Fprintln` per line, so the
final line is newline-terminated). Both default and projected counts use
the same convention.

## Byte table

| call | shape | default | projected | saved | saved % | default lines | projected lines |
|------|-------|--------:|----------:|------:|--------:|--------------:|----------------:|
| S1  | decision table, 2 facts (pricing 2x2)           | 267 | 130 | 137 | 51.3% | 15 | 9  |
| S2  | state machine, writes+clear (release ship-clean) | 338 | 154 | 184 | 54.4% | 16 | 9  |
| S2b | state machine, gate+writes (release begin)       | 368 | 207 | 161 | 43.7% | 18 | 11 |

Strictly fewer total bytes on every shape measured.

## S1 — pricing decision table (models/examples/pricing-decision-table.toml)

```
intrastate flow resolve --model models/examples/pricing-decision-table.toml \
    --outcome decide --tag tier=paid --tag region=eu --as text
```

Default (267 bytes, 15 lines):

```
clear: (none)
emit.dpa: required
emit.plan: pro
escaped: false
gates: (none)
model: models/examples/pricing-decision-table.toml
next: (none)
observed.region: eu
observed.tier: paid
outcome: decide
owned: (none)
readers: (none)
revision: 
rule: paid-eu
writes: (none)
```

Projected (130 bytes, 9 lines):

```
clear: (none)
emit.dpa: required
emit.plan: pro
escaped: false
gates: (none)
next: (none)
revision: 
rule: paid-eu
writes: (none)
```

## S2 — release grammar, ship-clean (models/examples/release-grammar.toml)

Owned state seeded to `phase=building`, `build-id=pending` via:

```
intrastate flow set-state --model models/examples/release-grammar.toml \
    --artifact release=release-state.artifact \
    --write phase=building --write build-id=pending --as json
```

```
intrastate flow resolve --model models/examples/release-grammar.toml \
    --artifact release=release-state.artifact \
    --outcome ship --tag 'checks=["tests","signoff"]' --tag risk=0 --as text
```

Default (338 bytes, 16 lines):

```
clear[0]: build-id
emit: (none)
escaped: false
gates: (none)
model: models/examples/release-grammar.toml
next.build-id: <clear>
next.phase: shipped
observed.checks: ["signoff","tests"]
observed.risk: 0
outcome: ship
owned.build-id: pending
owned.phase: building
readers[0]: release-state
revision: 
rule: ship-clean
writes.phase: shipped
```

Projected (154 bytes, 9 lines):

```
clear[0]: build-id
emit: (none)
escaped: false
gates: (none)
next.build-id: <clear>
next.phase: shipped
revision: 
rule: ship-clean
writes.phase: shipped
```

## S2b — release grammar, begin (gate on the selected row)

Owned state seeded to `phase=idle`, `build-id=none` via `flow set-state`
(same form as S2).

```
intrastate flow resolve --model models/examples/release-grammar.toml \
    --artifact release=release-idle.artifact \
    --outcome build --tag risk=0 --tag 'checks=[]' --as text
```

Default (368 bytes, 18 lines):

```
clear: (none)
emit: (none)
escaped: false
gates[0].id: change-window
gates[0].result: allow
model: models/examples/release-grammar.toml
next.build-id: pending
next.phase: building
observed.checks: []
observed.risk: 0
outcome: build
owned.build-id: none
owned.phase: idle
readers[0]: release-state
revision: 
rule: begin
writes.build-id: pending
writes.phase: building
```

Projected (207 bytes, 11 lines):

```
clear: (none)
emit: (none)
escaped: false
gates[0].id: change-window
gates[0].result: allow
next.build-id: pending
next.phase: building
revision: 
rule: begin
writes.build-id: pending
writes.phase: building
```

## Subset corollary

Every projected line is byte-identical to a line in the default rendering,
and the projected rendering is an order-preserving subsequence of it: on
all three shapes a line-diff of default against projected reports only
deletions and zero additions, and zero projected lines fail an exact
whole-line match against the default. The retained PLAN lines carry
through unchanged in content and in relative order.

Lines deleted per shape:

- S1: `model`, `observed.region`, `observed.tier`, `outcome`, `owned`,
  `readers` (6 lines).
- S2: `model`, `observed.checks`, `observed.risk`, `outcome`,
  `owned.build-id`, `owned.phase`, `readers[0]` (7 lines).
- S2b: `model`, `observed.checks`, `observed.risk`, `outcome`,
  `owned.build-id`, `owned.phase`, `readers[0]` (7 lines).

## Reading

A8 is confirmed. Text mode reduces by 43.7-54.4% across the three shapes,
against 42.2-51.3% for the same shapes in JSON (A1's S1, S2, S2b), so the
strict-width clause is not merely satisfiable in text — the text reduction
is slightly larger on every shape. The reason is that flattening prices
each echo leaf at its own full path plus a newline, whereas JSON amortizes
the container over its members; the echo group is leaf-dense (`observed`,
`owned`, `readers` are all containers of scalars) while the retained PLAN
group carries the `(none)` empty-container lines that set the floor.

The reduction is a strict byte reduction and not a reformat: the projected
rendering is a byte-identical subsequence of the default, so the same
clause can be stated once over both modes rather than per-mode.
