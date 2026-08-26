Model: claude-opus-5

# 0011:A5 spike — `flow next` over a decision-table-shaped model, no `--tag`

Purpose: establish what the CURRENT build reports for `unresolved` /
`required` over 0010's decision-table class, and isolate which half of A5
depends on the unimplemented `0011:C1` probe.

Binary: `./bin/intrastate` (working tree clean; no tracked file changed).

## Constraint discovered first: a true 0010 decision table is UNLOADABLE today

`0010` is unimplemented. `class = "decision-table"` is not a known
`[model]` key, and — more decisively — an ordinary rule with no write block
is refused at load, so a ZERO-OWNED model cannot be built on this binary.
Both attempts, verbatim:

```
$ ./bin/intrastate flow next --model dt.toml --as=json
{"code":"flow-model-invalid","message":"the selected model could not be loaded","findings":[{"code":"malformed_rule_shape","message":"malformed_rule_shape: ordinary rule draft-small carries no write block",...}]}
EXIT=2

$ ./bin/intrastate flow next --model dt-obswrite.toml --as=json     # write to an OBSERVED tag instead
{"code":"flow-model-invalid","message":"the selected model could not be loaded","findings":[{"code":"write_to_non_owned_tag","message":"write_to_non_owned_tag: rule draft-small writes the observed tag status",...}]}
EXIT=2
```

That is `0010`'s own Alternative 1 ("the dummy-tag convention") observed as
a hard load refusal, and it is why the spike fixture below carries one
dummy owned tag `answer` plus the reader/writer pair it drags in. The
consequence for A5: on this build `required` is `["answer"]`, NEVER empty.
The "empty `required`" half of A5 / `0010:A11` / `0010:MVV` step 5 is NOT
attestable on the current binary — it is downstream of 0010 shipping.

Lint on both fixtures returns `graph-dangling-edge` ("the model declares no
initial owned state"), which is exactly `0010:MVV` step 6's stated negative
control for a `class`-omitted model. `lint` does not gate `flow next`, so
the `flow next` observations below stand.

## Fixture A — four-row decision table (`dt-write.toml`)

```toml
outcomes = ["locate"]

[model]
id = "navigator"
version = 1
description = "Two-dimension stateless decision table."

[tags.status]
provenance = "observed"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true

[tags.size]
provenance = "observed"
kind = "enum"
domain = ["small", "large"]
single_valued = true

# --- dummy owned tag: present ONLY because an ordinary rule must write ---
[tags.answer]
provenance = "owned"
kind = "enum"
domain = ["a1", "a2", "a3", "a4"]
single_valued = true

[read.nav]
role = "nav"
path = "nav.answer"
keys = ["answer"]
timeout = "2s"

[write.nav]
role = "nav"
path = "nav.answer"
keys = ["answer"]
timeout = "2s"
read_back = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[[rule]]
id = "draft-small"
[rule.write]
answer = "a1"
[rule.match.status]
eq = "Draft"
[rule.match.size]
eq = "small"
[rule.match.recognized]
eq = "locate"

[[rule]]
id = "draft-large"
[rule.write]
answer = "a2"
[rule.match.status]
eq = "Draft"
[rule.match.size]
eq = "large"
[rule.match.recognized]
eq = "locate"

[[rule]]
id = "final-small"
[rule.write]
answer = "a3"
[rule.match.status]
eq = "Final"
[rule.match.size]
eq = "small"
[rule.match.recognized]
eq = "locate"

[[rule]]
id = "final-large"
[rule.write]
answer = "a4"
[rule.match.status]
eq = "Final"
[rule.match.size]
eq = "large"
[rule.match.recognized]
eq = "locate"
```

Artifact `nav.json` (empty, so the dummy owned key stays unresolved):

```json
{}
```

### A.1 — no `--tag`

```
$ ./bin/intrastate flow next --model dt-write.toml --artifact nav=nav.json --as=json
```
```json
{
    "type": "ok",
    "data": {
        "model": "/private/tmp/claude-501/-Users-cwensel-sandbox-newcoinc-intrastate/1e287595-8a01-4cf0-a9e8-809d96ea1be1/scratchpad/dt-write.toml",
        "revision": "",
        "observed": {},
        "owned": {},
        "readers": [
            "nav"
        ],
        "outcomes": [
            "locate"
        ],
        "candidates": [
            {
                "rule": "draft-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "status",
                    "answer"
                ],
                "next": {
                    "answer": "a2"
                },
                "writes": {
                    "answer": "a2"
                },
                "clear": []
            },
            {
                "rule": "draft-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "status",
                    "answer"
                ],
                "next": {
                    "answer": "a1"
                },
                "writes": {
                    "answer": "a1"
                },
                "clear": []
            },
            {
                "rule": "final-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "status",
                    "answer"
                ],
                "next": {
                    "answer": "a4"
                },
                "writes": {
                    "answer": "a4"
                },
                "clear": []
            },
            {
                "rule": "final-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "status",
                    "answer"
                ],
                "next": {
                    "answer": "a3"
                },
                "writes": {
                    "answer": "a3"
                },
                "clear": []
            }
        ]
    }
}
```
EXIT=0

Attested by this build:
- exit 0 — no refusal over the decision-table shape;
- every ordinary row reported (4 of 4);
- the row's OBSERVED match keys `status` and `size` DO appear in
  `unresolved`;
- `required` is `["answer"]` — NOT empty (the dummy-tag tax, above).

Not attested (depends on the unimplemented `0011:C1` probe): that the
candidate SET is match-conditioned. Today `excluded()` in
`internal/cli/flow_next.go` sets `probe.Match = nil`, so match plays no
part in exclusion and the candidate set is the plain enumeration.

### A.2 — partial `--tag status=Draft`

```
$ ./bin/intrastate flow next --model dt-write.toml --artifact nav=nav.json --tag status=Draft --as=json
```
```json
{
    "type": "ok",
    "data": {
        "model": "/private/tmp/claude-501/-Users-cwensel-sandbox-newcoinc-intrastate/1e287595-8a01-4cf0-a9e8-809d96ea1be1/scratchpad/dt-write.toml",
        "revision": "",
        "observed": {
            "status": "Draft"
        },
        "owned": {},
        "readers": [
            "nav"
        ],
        "outcomes": [
            "locate"
        ],
        "candidates": [
            {
                "rule": "draft-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "answer"
                ],
                "next": {
                    "answer": "a2"
                },
                "writes": {
                    "answer": "a2"
                },
                "clear": []
            },
            {
                "rule": "draft-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "answer"
                ],
                "next": {
                    "answer": "a1"
                },
                "writes": {
                    "answer": "a1"
                },
                "clear": []
            },
            {
                "rule": "final-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "answer"
                ],
                "next": {
                    "answer": "a4"
                },
                "writes": {
                    "answer": "a4"
                },
                "clear": []
            },
            {
                "rule": "final-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "size",
                    "answer"
                ],
                "next": {
                    "answer": "a3"
                },
                "writes": {
                    "answer": "a3"
                },
                "clear": []
            }
        ]
    }
}
```
EXIT=0

TODAY: `status` leaves every row's `unresolved` (it is now in the view),
but NO ROW DROPS — `final-small` and `final-large` are still reported even
though they match `status = "Final"` and `Draft` was supplied. That is the
stripped-match behaviour, and it is precisely the defect `0011:C1` exists
to fix.

UNDER `0011:C1`: `status` is PRESENT in the assembled view, so each row's
`status` atom is carried into the one-row probe and must hold under the
kernel's match comparison. `draft-small` and `draft-large` hold and remain
candidates with `unresolved = ["size", "answer"]`; `final-small` and
`final-large` are supplied-but-MISMATCHED, the kernel refuses the probe
`no_match`, and they DROP — they must not be reported, and under
`--evaluate-gates` their gates must not run. So yes: a row whose match key
is supplied-but-mismatched drops out.

### A.3 — full `--tag status=Draft --tag size=small`

```
$ ./bin/intrastate flow next --model dt-write.toml --artifact nav=nav.json --tag status=Draft --tag size=small --as=json
```
```json
{
    "type": "ok",
    "data": {
        "model": "/private/tmp/claude-501/-Users-cwensel-sandbox-newcoinc-intrastate/1e287595-8a01-4cf0-a9e8-809d96ea1be1/scratchpad/dt-write.toml",
        "revision": "",
        "observed": {
            "size": "small",
            "status": "Draft"
        },
        "owned": {},
        "readers": [
            "nav"
        ],
        "outcomes": [
            "locate"
        ],
        "candidates": [
            {
                "rule": "draft-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "answer"
                ],
                "next": {
                    "answer": "a2"
                },
                "writes": {
                    "answer": "a2"
                },
                "clear": []
            },
            {
                "rule": "draft-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "answer"
                ],
                "next": {
                    "answer": "a1"
                },
                "writes": {
                    "answer": "a1"
                },
                "clear": []
            },
            {
                "rule": "final-large",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "answer"
                ],
                "next": {
                    "answer": "a4"
                },
                "writes": {
                    "answer": "a4"
                },
                "clear": []
            },
            {
                "rule": "final-small",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "answer"
                ],
                "next": {
                    "answer": "a3"
                },
                "writes": {
                    "answer": "a3"
                },
                "clear": []
            }
        ]
    }
}
```
EXIT=0

TODAY: all four rows still reported; `unresolved` has narrowed to
`["answer"]` on every row. Under `0011:C1` only `draft-small` survives.

## Fixture B — the CRUX: a row whose only match atom is `recognized`

A5 claims every ordinary rule is reported "with its observed match keys in
`unresolved`". A row that matches on NO observed dimension — only on
`recognized` — is the falsifier. `recognized` is lifted out of `Row.Atoms`
at normalization (`internal/table/normalize.go`, the outcome binding is
moved into `Row.Outcome`), and `0011:C1` binds `recognized` to the row's
own outcome in the assembled view. So a `recognized` atom is NEVER an
unresolved fact, by construction.

`dt-crux.toml`:

```toml
outcomes = ["locate", "fallback"]

[model]
id = "crux"
version = 1

[tags.status]
provenance = "observed"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true

[tags.answer]
provenance = "owned"
kind = "enum"
domain = ["a1", "a2"]
single_valued = true

[read.nav]
role = "nav"
path = "nav.answer"
keys = ["answer"]
timeout = "2s"

[write.nav]
role = "nav"
path = "nav.answer"
keys = ["answer"]
timeout = "2s"
read_back = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[[rule]]
id = "row-status"
[rule.match.status]
eq = "Draft"
[rule.match.recognized]
eq = "locate"
[rule.write]
answer = "a1"

# Its ONLY match atom is `recognized` — no observed dimension at all.
[[rule]]
id = "row-recognized-only"
[rule.match.recognized]
eq = "fallback"
[rule.write]
answer = "a2"
```

### B.1 — no `--tag`, empty artifact (dummy owned key unresolved)

```
$ ./bin/intrastate flow next --model dt-crux.toml --artifact nav=nav.json --as=json
```
```json
{
    "type": "ok",
    "data": {
        "model": "/private/tmp/claude-501/-Users-cwensel-sandbox-newcoinc-intrastate/1e287595-8a01-4cf0-a9e8-809d96ea1be1/scratchpad/dt-crux.toml",
        "revision": "",
        "observed": {},
        "owned": {},
        "readers": [
            "nav"
        ],
        "outcomes": [
            "locate",
            "fallback"
        ],
        "candidates": [
            {
                "rule": "row-recognized-only",
                "outcome": "fallback",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "answer"
                ],
                "next": {
                    "answer": "a2"
                },
                "writes": {
                    "answer": "a2"
                },
                "clear": []
            },
            {
                "rule": "row-status",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "status",
                    "answer"
                ],
                "next": {
                    "answer": "a1"
                },
                "writes": {
                    "answer": "a1"
                },
                "clear": []
            }
        ]
    }
}
```
EXIT=0

`row-recognized-only` has `unresolved: ["answer"]` — non-empty ONLY because
the dummy owned tag is unread. Strip the dummy tag (which is exactly what
`0010` does) and nothing is left.

### B.2 — the same model with the owned key RESOLVED, simulating 0010's empty `required`

`nav-full.json`:

```json
{"answer":"a1"}
```

```
$ ./bin/intrastate flow next --model dt-crux.toml --artifact nav=nav-full.json --as=json
```
```json
{
    "type": "ok",
    "data": {
        "model": "/private/tmp/claude-501/-Users-cwensel-sandbox-newcoinc-intrastate/1e287595-8a01-4cf0-a9e8-809d96ea1be1/scratchpad/dt-crux.toml",
        "revision": "",
        "observed": {},
        "owned": {
            "answer": "a1"
        },
        "readers": [
            "nav"
        ],
        "outcomes": [
            "locate",
            "fallback"
        ],
        "candidates": [
            {
                "rule": "row-recognized-only",
                "outcome": "fallback",
                "required": [
                    "answer"
                ],
                "unresolved": [],
                "next": {
                    "answer": "a2"
                },
                "writes": {
                    "answer": "a2"
                },
                "clear": []
            },
            {
                "rule": "row-status",
                "outcome": "locate",
                "required": [
                    "answer"
                ],
                "unresolved": [
                    "status"
                ],
                "next": {
                    "answer": "a1"
                },
                "writes": {
                    "answer": "a1"
                },
                "clear": []
            }
        ]
    }
}
```
EXIT=0

**`"rule":"row-recognized-only" … "unresolved":[]`** — an ordinary,
reported candidate with an EMPTY `unresolved` list, on the current build,
with no `--tag` supplied. Under `0011:C1` over a genuine `0010` decision
table (zero owned tags, `required = []`), this row's `unresolved` is empty
too: its only match atom is `recognized`, which C1 binds to the row's own
outcome and which is therefore PRESENT in the assembled view, hence not an
unresolved fact.

This does not refute `0010:A11` or `0010:MVV` step 5 — both speak only of
exit 0 and empty `required`, and say nothing about `unresolved` being
non-empty. It refutes the UNIVERSAL quantifier in `0011:A5`'s own wording.
