Model: claude-fable-5

# 0028 propose — joint-decision check, three arms (run on the written proposal)

## Arm 1 — `rdr index --anchor-intersect --json --record 0028-declared-line-edit-writer` (repo-resolved)

```json
{
  "open_only": true,
  "overlaps": [
    {
      "records": [
        "0027",
        "0028"
      ],
      "anchors": [
        "internal/table/load.go::carrierDefect"
      ],
      "cited": false
    }
  ],
  "record": "0028",
  "schema": "2"
}
```

## Arm 2 — `rdr index --literal-intersect --json --record 0028-declared-line-edit-writer`

```json
{
  "open_only": true,
  "overlaps": [
    {
      "records": [
        "0014",
        "0028"
      ],
      "anchors": [
        "intrastate lint"
      ],
      "cited": false
    },
    {
      "records": [
        "0016",
        "0028"
      ],
      "anchors": [
        "read_back_incomplete"
      ],
      "cited": false
    },
    {
      "records": [
        "0021",
        "0028"
      ],
      "anchors": [
        "intrastate lint"
      ],
      "cited": false
    },
    {
      "records": [
        "0026",
        "0028"
      ],
      "anchors": [
        "execution_failure"
      ],
      "cited": false
    },
    {
      "records": [
        "0027",
        "0028"
      ],
      "anchors": [
        "table.Categories()"
      ],
      "cited": false
    }
  ],
  "record": "0028",
  "schema": "2"
}
```

## Arm 3 — absence (manual)

Refusals this proposal converts to acceptances: 0025:C5 `command_and_path_conflict` "neither" arm (an entry carrying only `edit`); `command_unknown_placeholder` for a declared `{tag.<key>}` argv element (C6).

Grep over depth-1 `*.md` whose first `- **Status**:` starts with Draft or Final, excluding 0028, for `command_and_path_conflict|command_unknown_placeholder|v1 complete|exactly one of \`path\` / \`command\`|{artifact}`:

```
0027-inline-shell-detection-scope.md [Draft] hits=8
```

Only 0027 hits; its reliance is 0027:C1 "clause 3 coupling … never command_unknown_placeholder", preserved by C6. Terminal peers carrying the tokens: 0025 (Implemented; the declared Overrides target — rides to 7.1).

## Verdict

fired → 0027, 0026, 0016 (home: OPEN). See the record's Decision Rationale `Joint-check:` line for the per-pair reading and candidate dispositions.

## Dispositions (2026-08-31)

Decided by the human; applied to this record only.

- **0027 ↔ 0028** (`internal/table/load.go::carrierDefect`, `table.Categories()`) — cite-don't-restate, home `cli/0025:C5`. Both records fence only their own deltas against 0025:C5 (Implemented owner of the clause map and within-entry precedence); 0027's line is already symmetric. Changed here: `Joint-check:` names the home as `cli/0025:C5 for cli/0027:C1`; no fence edit.
- **0026 ↔ 0028** (`execution_failure`) — hoisted to `JDR 0003 §D1` (graded constraint). Changed here: `0028:C3` `order:` drops `applied: false` from `Detail` (rule id `<id>.edit.<key>` and reason token remain; applied sense is `Applied()`, executor-set) and cites the entry; `0028:A8` claim aligned (Detail carries rule id + reason token, applied sense rides `Applied()`), Evidence set to the constraint entry, If-wrong "owned here" voided (a typed carrier field is the registry's decision); References gains the JDR line. Status stays Pending for Stage 4.
- **0016 ↔ 0028** (`read_back_incomplete`) — cite-don't-restate, home `cli/0016:C4` (the pre-mutation gate check relies on 0016:C4's fail-closed `readerFor`). 0016 (Final) untouched. Changed here: `Joint-check:` names the home as `cli/0016:C4`.

`Joint-check:` head now reads `(home: JDR 0003 §D1 for cli/0026:C1; cli/0025:C5 for cli/0027:C1; cli/0016:C4)`; the OPEN candidate-dispositions sentence is replaced by `Dispositions above.` No other `applied: false` occurrences remained in the body.

## Extension — 2026-09-03 (finalize, Stage 7)

The Stage-7 joint-decision fence (`--outcome fence`) returned
`stopped:overlap-uncited`. `rdr index --anchor-intersect --record 0028` named two
in-flight peers sharing a source anchor with 0028 and citing neither direction:

| pair | shared anchor | disposition |
|---|---|---|
| 0019 ↔ 0028 | `internal/cli/flowbind/flowbind.go::load` | cite-don't-restate |
| 0020 ↔ 0028 | `internal/cli/flow_input.go::parseTags`   | cite-don't-restate |

Arm 1 (modify-anchors) re-run on both, against the written proposals:

**0019 — `flowbind.go::load`.** Neither record changes the function. Both read its
absent-file-is-empty rule identically, as existing so a first write is possible
(0019 §problem-statement, §approach, §investigation; 0028 S27b). 0019 RELIES on it —
it is what lets `init-state` seed an empty store — while 0019:C1 forbids any READ
path synthesizing `[initial]`. 0028 fences it OUT for the `edit` carrier: a bound
target that cannot be read refuses with `execution_failure` before mutation rather
than treating absent as empty, because `edit`'s `input:`/`select:` fences creation
out. Complementary positions on one unchanged rule; no shared answer to hoist.

**0020 — `flow_input.go::parseTags`.** 0020:C1 answers the UNDECLARED key at `--tag`
admission (pure carrier, byte-verbatim, no new refusal code minted) and explicitly
preserves the `flow-tag-owned` (REQ-27) arm "unchanged, and still preceding any
accessor". 0028:A7 relies on exactly that preserved arm: `parseTags` refusing a
`--tag` naming an owned key is what enforces the context-only semantic, which is
what makes `{tag.nnnn}` bindable on the invocation with no reader run and no
caller-typed planned value. 0020 owns the undeclared arm; 0028 consumes the
declared-and-owned refusal 0020 leaves untouched. Disjoint arms of one function.

No hoist warranted for either: a JDR exists to home an answer two records would
otherwise both give, and in both pairs each record answers a different question.
The JC element in the record carries these dispositions.

Also re-confirmed from the same intersect run: `0014/0021 ↔ 0028` on the literal
`intrastate lint` remains a tool name, not a decision (disposed 2026-08-31, unchanged).
