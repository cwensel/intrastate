Model: claude-opus-5

# A6 — `models/rdr.toml` narrowing at `stage=resolved`

A6 claims: "`models/rdr.toml` at `stage=resolved` narrows from 21 reported
rules to exactly the rows whose `match.stage` holds over the artifact's owned
`stage`, each ordinary row's `recognized` match holding trivially because the
probe binds `recognized` to the row's own outcome."

Verified in two parts, because the match-conditioned predicate (C1) is not yet
implemented — today `internal/cli/flow_next.go::excluded` sets `probe.Match =
nil`.

## Part 1 — the baseline, live against the current build

Seeded artifact (the `flowbind` on-disk format is a flat JSON object of tag key
to tag value; the model's `[read.rdr-status]` declares `role = "rdr"` and
`keys = ["stage", "status", "gate_passed"]`):

```
$ go build -o bin/intrastate ./cmd/intrastate

$ cat $SCRATCH/rdr-resolved.json
{"stage":"resolved","status":"draft","gate_passed":"false"}

$ ./bin/intrastate flow next \
    --model models/rdr.toml \
    --artifact rdr=$SCRATCH/rdr-resolved.json \
    --as=json
```

Exit 0. Payload head:

```json
{
  "type": "ok",
  "data": {
    "model": "models/rdr.toml",
    "revision": "",
    "observed": {},
    "owned": {
      "gate_passed": "false",
      "stage": "resolved",
      "status": "draft"
    },
    "readers": ["rdr-status"],
    "outcomes": ["advance", "revise", "abandon"],
    "candidates": [ ... ]
  }
}
```

Candidate roll (`rule -> outcome`), 21 entries:

```
final-abandon        -> abandon
final-route-back     -> revise
finalize-blocked     -> advance
implement            -> advance
prelock              -> advance
prelock-abandon      -> abandon
prelock-route-back   -> revise
propose              -> advance
propose-abandon      -> abandon
propose-again        -> revise
reconcile            -> advance
reconcile-abandon    -> abandon
reconcile-route-back -> revise
refine               -> advance
refine-abandon       -> abandon
refine-again         -> revise
resolve-abandon      -> abandon
resolve-assumptions  -> advance
resolve-route-back   -> revise
seed-abandon         -> abandon
seed-revise          -> revise
```

**Counts.** `grep -c '^\[\[rule\]\]' models/rdr.toml` = **22**. Normalized
non-escape rows = **22** (no rule declares `escape`, and no `match ... in`
expansion multiplies a row). Reported candidates = **21**.

The single excluded row is `finalize-pass`: its `[rule.guard.all.gate_passed]
eq = "true"` is decided FALSE over the seeded `gate_passed = "false"`, so the
kernel probe answers `no_match`. The RDR's "21 of 22" is **confirmed** at this
artifact — but note the exclusion is guard-driven, so 21 is a property of
`gate_passed = "false"`, not of `stage = "resolved"`. A seed carrying
`gate_passed = "true"` would report 21 with `finalize-blocked` excluded
instead; the count is stable at 21 either way because the two finalize rows
partition the `gate_passed` dimension.

## Part 2 — the predicted C1 set

Prototyped in a throwaway `internal/cli` test (deleted after the run; `git
status` clean) applying C1's predicate: probe carries the row's match atoms
over keys PRESENT in the assembled view, drops match atoms over absent keys,
drops the escape list, and binds `Recognized: row.Outcome`.

```
total rows (normalized): 22; non-escape rows: 22
match.recognized atoms across non-escape rows: 0
BASELINE (match stripped) count=21
PREDICTED (C1, match over present keys) count=3
   prelock            -> advance
   resolve-abandon    -> abandon
   resolve-route-back -> revise
```

Cross-checked by hand: `grep -n 'eq = "resolved"' models/rdr.toml` yields
exactly three hits (lines 212, 222, 232) — the `[rule.match.stage]` blocks of
`prelock`, `resolve-route-back`, and `resolve-abandon`. Every one of the 22
rules carries a `[rule.match.stage]` block, so there is no row with "no
`match.stage` atom at all" and no row left a candidate by C1's
absent-key-is-undecided provision. `stage` is `required = true` and served by
the invoked `rdr-status` reader, so it is always present in the view for this
model.

**Predicted narrowing: 21 -> 3.**

## The `recognized` sub-claim

`grep -n "RecognizedTagKey" internal/table/normalize.go` and
`internal/table/normalize.go:428-436`:

```go
	// Normalization lifts that atom OUT of the predicate set into the row's
	// outcome field.
	predicates := make([]Atom, 0, len(matchAtoms))
	for _, a := range matchAtoms {
		if a.Key != RecognizedTagKey {
			predicates = append(predicates, a)
		}
	}
```

Every rule in `models/rdr.toml` authors a `[rule.match.recognized]` block, and
`outcomeBinding` requires exactly one per rule (`0002:C13`). But RDR 0002
**consumes** it: the atom becomes `Row.Outcome` and never reaches `Row.Atoms`,
so `Row.KernelRow().Match` carries **zero** `recognized` tags. Measured:
`match.recognized atoms across non-escape rows: 0`.

The probe binding is real —
`internal/cli/flow_next.go::excluded` passes `Recognized: row.Outcome`, and
`internal/resolve/resolve.go::assemble` (line 222) writes that into the view
under `recognizedTagKey = "recognized"`. It is load-bearing for the kernel's
`row.Outcome != in.Recognized` candidacy filter and for the one-row probe
table's `Outcomes: []string{row.Outcome}`, which is what keeps the probe from
refusing `unmodeled_outcome`. It is **not**, however, discharging any
`match.recognized` atom, because no such atom survives normalization.

The sub-claim is therefore **vacuous, not false**: there is nothing for it to
hold over.
