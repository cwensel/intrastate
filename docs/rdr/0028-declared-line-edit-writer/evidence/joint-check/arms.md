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
