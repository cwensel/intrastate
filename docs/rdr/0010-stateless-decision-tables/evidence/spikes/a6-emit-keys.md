Model: claude-opus-5[1m]

# A6 — the motivating consumer's actual emit keys (Stage 6 re-verification)

A6 was flipped Verified → Pending by the critique lens (D-7): its *supporting
reading* was wrong ("a flat string carries `Next: /rdr-prelock 0046 critique`
whole" — it cannot; `0046` is caller data). The narrowed claim is the checkable
one, and this is its re-verification against **rdr#tmxk's actual emit keys**,
not against memory.

## What the consumer is

`kata show tmxk` (repo `~/sandbox/newcoinc/rdr`):

> `rdr/models/rdr-status.toml`: the routing table as an intrastate model, linted
> for coverage. Port rdr-status's "How it decides next" + rdr-common §lens-row to
> rules over N2's facts; **the selected row (or an `emit`) is the answer.**

So the answer this table must produce is one **routing decision**: which stage
command to run next, and (for Stage 5) which lens.

## The answer's actual shape, read from the skill that renders it

`skills/rdr-status/SKILL.md` §Output, item 3 (line 168):

> 3. **Next** — the exact command to run, e.g. `Next: /rdr-prelock 0046 critique`.

Decomposed, that rendered line has exactly three parts:

| Part | Example | Known when the table is authored? | Supplied by |
| --- | --- | --- | --- |
| stage verb | `/rdr-prelock` | **yes** — closed set of stage skills | the table |
| lens | `critique` | **yes** — closed set (`grounding`, `cove`, `3amigo`, `critique`, `repeatability`) | the table |
| record number | `0046` | **no** — the record under inspection | the **caller** |

The caller already holds the record number: it is the argument the user typed
(`/rdr-status 0046`). It is never a table value, and the table never needs to
interpolate it. This is exactly the split A6's narrowed claim asserts.

## Every input the table matches on is a closed enum of fixed strings

`models/rdr-facts.toml` — the fact vocabulary the routing table matches against
is the only contract between the two files ("the fact names are the only contract
between rdr-facts.toml and this model"). Every declared fact carries a closed
`domain` of authoring-time string literals:

```
79:  domain = ["Draft", "Final", "Implemented", "Deferred", "Demoted", "Abandoned", "Superseded", "Reverted", "Rejected"]
99:  domain = ["none", "demoted-target", "revised-from", "joint-decision", "revisit-when", "bracketed", "dash", "parenthetical"]
144: domain = ["small", "mid", "large", "foundational"]
257: domain = ["all-pending", "mixed", "all-terminal", "unknown-plan"]
552: domain = ["COMPLETE", "INCOMPLETE", "IN-PROGRESS"]
```

Facts are `probe` (path exists: bool) or `field` (a projection path into
`rdr inspect --json`). Neither kind produces a value the *table* must render
into an emit string — they are match inputs, not emit outputs.

## Verdict

**VERIFIED.** Every value the motivating consumer needs *from the table* is a
fixed string known at authoring time:

- the stage verb — a closed set of skill names;
- the lens — a closed set of lens names;
- (and for judgment cells, per the tracker, a stop-packet row — again a fixed
  string, the `stopped:<code>` shape).

Per-invocation data — the RDR number, the slug, a path — is held by the caller
and substituted by it, never interpolated by the table. A flat, string-valued
`[rule.emit]` table is therefore sufficient for this consumer.

The widening A6 defers (a template language *inside* an emit value) is not
needed here and is not provided: C3 makes emit values uninterpreted and compared
by exact byte equality, so a `{id}` in a value would be literal text.

## Evidence pointers

- `kata show tmxk` → the consumer's own statement of the answer shape.
- `rdr/skills/rdr-status/SKILL.md:168` → the rendered `Next:` line, the answer.
- `rdr/skills/rdr-status/SKILL.md:131-150` → the routing inputs (`Profile` latch,
  §lens-row, per-lens Stage-5 fork) — all closed vocabularies.
- `rdr/models/rdr-facts.toml:79,99,144,257,552` → the closed fact domains.
