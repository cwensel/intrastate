Model: claude-opus-5

# 0011:A15 spike — extending `invokedReaders`' demand set with match-block owned keys

A15 claims the extension invokes strictly MORE readers and never fewer, for
BOTH callers of `invokedReaders`; changes no reader's contract; is a no-op
over a zero-owned-tag model; and changes `flow resolve` in exactly one
class — a model with a match-only owned key — in two directions.

Verified LIVE, by building both arms: the extension was applied as a
temporary uncommitted edit, both binaries kept, every arm run against each,
and the edit then reverted. `git status` at the end shows no modified
tracked file.

Baseline: branch `main`, commit `1a7bddf`.

## The local edit, verbatim

`git diff -- internal/cli/flow_exec.go` while the spike ran:

```diff
diff --git a/internal/cli/flow_exec.go b/internal/cli/flow_exec.go
index 94f1654..6cfceaf 100644
--- a/internal/cli/flow_exec.go
+++ b/internal/cli/flow_exec.go
@@ -107,6 +107,9 @@ func invokedReaders(m *table.Model, outcome string) []string {
 		for _, key := range guardOwnedKeys(m, row) {
 			demanded[key] = true
 		}
+		for _, key := range matchOwnedKeys(m, row) {
+			demanded[key] = true
+		}
 	}
 
 	var out []string
@@ -144,6 +147,24 @@ func guardOwnedKeys(m *table.Model, row table.Row) []string {
 	return out
 }
 
+// matchOwnedKeys returns the OWNED tag keys one row's MATCH atoms
+// reference (A15 spike, uncommitted).
+func matchOwnedKeys(m *table.Model, row table.Row) []string {
+	var out []string
+	for _, atom := range row.Atoms {
+		if atom.Block != table.BlockMatch {
+			continue
+		}
+		if m.Tags[atom.Key].Provenance != table.ProvenanceOwned {
+			continue
+		}
+		if !slices.Contains(out, atom.Key) {
+			out = append(out, atom.Key)
+		}
+	}
+	return out
+}
+
 // declaredReaders returns every declared reader, sorted. `flow read-state`
 // is the single exception to narrowing: a diagnostic read has no candidate
 // set to narrow by (REQ-37).
```

It is a THIRD union term alongside `Row.RequiresOwned` and
`guardOwnedKeys`. Adding to a set can only grow it, so "strictly more,
never fewer" is structural; the arms below measure where the growth is
observable.

Both binaries:

```
$ go build -o $SCRATCH/pre  ./cmd/intrastate     # at 1a7bddf, unmodified
$ go build -o $SCRATCH/post ./cmd/intrastate     # with the diff above applied
BUILD_PRE_OK
BUILD_POST_OK
```

## Fixture 1 — the match-only owned key (`mo.toml`)

`mode` is owned, MATCHED on by the one rule, and never written, cleared, or
guarded. `answer` is the written owned key that keeps the ordinary rule
legally shaped (an ordinary rule with no write block is refused
`malformed_rule_shape`, and a write to an observed tag `write_to_non_owned_tag`
— see the a5 spike). `checkAccessorBindings` requires exactly one reader per
owned tag, so `mode` gets `modereader` and `answer` gets `ansreader`; the two
sit on DIFFERENT artifact roles, which is what makes the unbound arm
reachable.

```toml
outcomes = ["go"]

[model]
id = "matchonly"
version = 1

# The MATCH-ONLY owned key: some row MATCHES on it; no row writes, clears,
# or guards it.
[tags.mode]
provenance = "owned"
kind = "enum"
domain = ["fast", "slow"]
single_valued = true

# The written owned key, so the ordinary rules have a legal write block.
[tags.answer]
provenance = "owned"
kind = "enum"
domain = ["a1", "a2"]
single_valued = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

# Exactly one reader for `mode` — the reader whose invocation is at stake.
[read.modereader]
role = "mode"
path = "mode.state"
keys = ["mode"]
timeout = "2s"

[read.ansreader]
role = "ans"
path = "ans.state"
keys = ["answer"]
timeout = "2s"

[write.answriter]
role = "ans"
path = "ans.state"
keys = ["answer"]
timeout = "2s"
read_back = true

[[rule]]
id = "match-on-mode"
[rule.match.mode]
eq = "fast"
[rule.match.recognized]
eq = "go"
[rule.write]
answer = "a1"
```

Artifacts:

```
$ cat mode.json      -> {"mode":"fast"}
$ cat ans.json       -> {}
$ cat ans-full.json  -> {"answer":"a2"}
```

`mo-refuse.toml` is `mo.toml` with one line changed, so `modereader`
REFUSES rather than answering — `flowbind.Reader.Read` returns an error on
an `unreachable` declared locator:

```
$ sed 's|^path = "mode.state"|path = "mode.state-unreachable"|' mo.toml > mo-refuse.toml
$ grep -n 'unreachable' mo-refuse.toml
31:path = "mode.state-unreachable"
```

## Fixture 2 — the BREAKING arm (`mo-breaking.toml`)

`mo.toml` plus a SECOND ordinary row on the same outcome `go` that does NOT
match on `mode`. Everything else is identical. Only the rules differ:

```toml
# Row 1: matches on the match-only owned key `mode`.
[[rule]]
id = "match-on-mode"
[rule.match.mode]
eq = "fast"
[rule.match.recognized]
eq = "go"
[rule.write]
answer = "a1"

# Row 2: SAME outcome, does NOT match on `mode`. This is the row whose plan
# the pre-extension build yields when `mode`'s reader role is unbound.
[[rule]]
id = "plain-row"
[rule.match.recognized]
eq = "go"
[rule.write]
answer = "a2"
```

## Fixture 3 — the zero-owned-tag approximation (`dt-write.toml`)

The a5 spike's Fixture A, reused verbatim: 0010's four-row decision table
with the dummy owned `answer` tag. Its MATCH atoms are all over OBSERVED
tags (`status`, `size`) plus `recognized`, so the added term must find
nothing.

**Limit, carried over from a5 and unchanged here:** a GENUINE zero-owned
model is UNLOADABLE on this build — an ordinary rule with no write block is
refused `malformed_rule_shape`, and writing an observed tag is refused
`write_to_non_owned_tag`. `dt-write.toml` is therefore the nearest
reachable approximation: zero owned MATCH atoms, one dummy owned WRITTEN
tag. What (i) attests is that the added term contributes nothing when no
row matches on an owned key — which is the property A15's "no-op over
0010's `decision-table` class" actually turns on, since a genuine 0010
table has no owned tag to match on at all.

---

## (i) — the added term is EMPTY over the zero-owned-tag model

```
$ ./pre  flow next --model dt-write.toml --artifact nav=nav.json --as=json > zero.pre.json   # EXIT=0
$ ./post flow next --model dt-write.toml --artifact nav=nav.json --as=json > zero.post.json  # EXIT=0
$ diff zero.pre.json zero.post.json && echo "IDENTICAL"
IDENTICAL
```

The payload both builds emit (`readers` bolded by attention, not markup):

```json
{"type":"ok","data":{"model":"dt-write.toml","revision":"","observed":{},"owned":{},"readers":["nav"],"outcomes":["locate"],"candidates":[{"rule":"draft-large","outcome":"locate","required":["answer"],"unresolved":["size","status","answer"],"next":{"answer":"a2"},"writes":{"answer":"a2"},"clear":[]},{"rule":"draft-small","outcome":"locate","required":["answer"],"unresolved":["size","status","answer"],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[]},{"rule":"final-large","outcome":"locate","required":["answer"],"unresolved":["size","status","answer"],"next":{"answer":"a4"},"writes":{"answer":"a4"},"clear":[]},{"rule":"final-small","outcome":"locate","required":["answer"],"unresolved":["size","status","answer"],"next":{"answer":"a3"},"writes":{"answer":"a3"},"clear":[]}]}}
```

Confirmed directly on the reader set too (probe below): `dt-write.toml`
yields `readers=[nav]` at BOTH outcomes on BOTH arms.

**(i) HOLDS** (against the nearest-reachable model; see the limit above).

## (ii) — the match-only-owned-key fixture

### `flow next`

PRE — `modereader` is NOT invoked, and `mode` is absent from the view:

```
$ ./pre flow next --model mo.toml --artifact mode=mode.json --artifact ans=ans.json --as=json
{"type":"ok","data":{"model":"mo.toml","revision":"","observed":{},"owned":{},"readers":["ansreader"],"outcomes":["go"],"candidates":[{"rule":"match-on-mode","outcome":"go","required":["answer"],"unresolved":["mode","answer"],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[]}]}}
EXIT=0
```

`"readers":["ansreader"]`, `"owned":{}`, and `mode` sits in `unresolved`.

POST — `modereader` IS invoked and `mode` is in the view:

```
$ ./post flow next --model mo.toml --artifact mode=mode.json --artifact ans=ans.json --as=json
{"type":"ok","data":{"model":"mo.toml","revision":"","observed":{},"owned":{"mode":"fast"},"readers":["ansreader","modereader"],"outcomes":["go"],"candidates":[{"rule":"match-on-mode","outcome":"go","required":["answer"],"unresolved":["answer"],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[]}]}}
EXIT=0
```

`"readers":["ansreader","modereader"]`, `"owned":{"mode":"fast"}`, and
`mode` has LEFT `unresolved`. That is S8: the reader runs and the key
becomes a fact the row's match can be decided against.

### `flow resolve --outcome go`, three arms

PRE — all three arms return the SAME thing, because the reader is never
invoked in any of them and `mode` is never a fact:

```
$ ./pre flow resolve --model mo.toml --outcome go --artifact mode=mode.json --artifact ans=ans-full.json --as=json
{"code":"flow-no-match","message":"no rule matches the recognized outcome `go` over the assembled state","findings":[{"code":"flow-no-match","message":"no rule in the model responds to the recognized outcome `go`"}]}
EXIT=2

$ ./pre flow resolve --model mo.toml --outcome go --artifact ans=ans-full.json --as=json           # mode role UNBOUND
{"code":"flow-no-match","message":"no rule matches the recognized outcome `go` over the assembled state","findings":[{"code":"flow-no-match","message":"no rule in the model responds to the recognized outcome `go`"}]}
EXIT=2

$ ./pre flow resolve --model mo-refuse.toml --outcome go --artifact mode=mode.json --artifact ans=ans-full.json --as=json   # reader REFUSES
{"code":"flow-no-match","message":"no rule matches the recognized outcome `go` over the assembled state","findings":[{"code":"flow-no-match","message":"no rule in the model responds to the recognized outcome `go`"}]}
EXIT=2
```

`flow-no-match` in ALL THREE — including the arm where the reader is bound,
present, and would have answered `mode=fast`. The refusal is for want of a
reader the build never invoked. This is the "refused for want of a reader it
never invoked" direction A15 names.

POST — the three arms SEPARATE, exactly as claimed:

```
$ ./post flow resolve --model mo.toml --outcome go --artifact mode=mode.json --artifact ans=ans-full.json --as=json
{"type":"ok","data":{"model":"mo.toml","revision":"","observed":{},"owned":{"answer":"a2","mode":"fast"},"readers":["ansreader","modereader"],"outcome":"go","rule":"match-on-mode","gates":[],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[],"escaped":false}}
EXIT=0

$ ./post flow resolve --model mo.toml --outcome go --artifact ans=ans-full.json --as=json          # mode role UNBOUND
{"code":"flow-artifact-missing","message":"the artifact role `mode` that the read accessor `modereader` needs has no --artifact binding","param":"mode"}
EXIT=2

$ ./post flow resolve --model mo-refuse.toml --outcome go --artifact mode=mode.json --artifact ans=ans-full.json --as=json   # reader REFUSES
{"code":"flow-accessor-failed","message":"the accessor `modereader` could not be executed","param":"modereader"}
EXIT=3
```

- reader bound and answering -> **a plan** (`rule: match-on-mode`), exit 0;
- role unbound -> **`flow-artifact-missing`**, exit 2, naming the role and
  the accessor;
- reader refusing -> **the reader's OWN refusal** (`flow-accessor-failed`),
  exit 3 — an environment class, not a `flow-no-match` about the request.

And `flow-no-match` for `mode` appears in **NONE** of the three post
arms, while it IS what all three pre arms returned.

**(ii) HOLDS** in every particular.

## (ii-b) — the BREAKING arm

Pinned, not discovered: `mo-breaking.toml`, the `mode` role UNBOUND,
`answer` established so `plain-row` can plan.

PRE — the second row's plan, exit 0:

```
$ ./pre flow resolve --model mo-breaking.toml --outcome go --artifact ans=ans-full.json --as=json
{"type":"ok","data":{"model":"mo-breaking.toml","revision":"","observed":{},"owned":{"answer":"a2"},"readers":["ansreader"],"outcome":"go","rule":"plain-row","gates":[],"next":{"answer":"a2"},"writes":{"answer":"a2"},"clear":[],"escaped":false}}
EXIT=0
```

POST — `flow-artifact-missing`, exit 2:

```
$ ./post flow resolve --model mo-breaking.toml --outcome go --artifact ans=ans-full.json --as=json
{"code":"flow-artifact-missing","message":"the artifact role `mode` that the read accessor `modereader` needs has no --artifact binding","param":"mode"}
EXIT=2
```

A working plan becomes a refusal. The selected row (`plain-row`) does not
need `mode` at all — but the demand set is a UNION over the outcome's rows,
so `match-on-mode`'s demand pulls `modereader` in, and `runReaders` raises
`flow-artifact-missing` for the unbound role BEFORE any accessor runs. This
is the DEV-8 class, reproduced exactly as A15 predicts.

**(ii-b) HOLDS. This is a real BREAKING change**, observed live, not
inferred.

## (iii) — every shipped fixture and `models/rdr.toml` byte-identical

Seed (the a6 spike's, reused):

```
$ cat rdr-resolved.json
{"stage":"resolved","status":"draft","gate_passed":"false"}
```

`flow resolve` over `models/rdr.toml` at every declared outcome plus one
outside the alphabet, and `flow next`, compared byte for byte with `cmp`:

```
$ for o in advance revise abandon bogus; do
    ./pre  flow resolve --model models/rdr.toml --outcome $o --artifact rdr=rdr-resolved.json --as=json > rdr.$o.pre.json
    ./post flow resolve --model models/rdr.toml --outcome $o --artifact rdr=rdr-resolved.json --as=json > rdr.$o.post.json
    cmp -s rdr.$o.pre.json rdr.$o.post.json && echo "$o BYTE-IDENTICAL" || diff rdr.$o.pre.json rdr.$o.post.json
  done
advance  EXIT=0  BYTE-IDENTICAL
revise   EXIT=0  BYTE-IDENTICAL
abandon  EXIT=0  BYTE-IDENTICAL
bogus    EXIT=2  BYTE-IDENTICAL
next: BYTE-IDENTICAL
```

The payloads compared (pre == post):

```json
advance: {"type":"ok","data":{...,"readers":["rdr-status"],"outcome":"advance","rule":"prelock","gates":[],"next":{"stage":"prelocked"},...}}
revise:  {"type":"ok","data":{...,"readers":["rdr-status"],"outcome":"revise","rule":"resolve-route-back","gates":[],"next":{"stage":"refined"},...}}
abandon: {"type":"ok","data":{...,"readers":["rdr-status"],"outcome":"abandon","rule":"resolve-abandon","gates":[],"next":{"stage":"dropped","status":"abandoned"},...}}
bogus:   {"code":"flow-unmodeled-outcome","message":"the recognized outcome `bogus` is outside the model's declared alphabet","param":"bogus"}
```

The RDR's stated reason holds mechanically: `models/rdr.toml` declares ONE
reader, `rdr-status`, over `["stage","status","gate_passed"]`, and every
owned key matched there is also written, so the added term names only keys
`RequiresOwned` already demanded. The reader set is `[rdr-status]` before
and after, at every outcome.

The shipped test fixtures, via the suite. First, a leftover UNTRACKED
scratch file from an earlier spike (`internal/cli/zz_scratch_a12_spike_test.go`,
pre-existing, not this spike's) fails on both arms for its own reason
(`open : no such file or directory`) and was moved aside for the run, then
restored:

```
$ go test ./...                       # POST-extension, only the diff above applied
?   	github.com/newcoinc/intrastate/cmd/intrastate	[no test files]
ok  	github.com/newcoinc/intrastate/internal/accessor	0.935s
ok  	github.com/newcoinc/intrastate/internal/cli	0.469s
ok  	github.com/newcoinc/intrastate/internal/cli/clierr	0.446s
ok  	github.com/newcoinc/intrastate/internal/cli/flowbind	0.581s
ok  	github.com/newcoinc/intrastate/internal/graphlint	4.265s
ok  	github.com/newcoinc/intrastate/internal/guard	0.886s
ok  	github.com/newcoinc/intrastate/internal/resolve	1.059s
ok  	github.com/newcoinc/intrastate/internal/table	0.975s

$ go vet ./...
VET_OK
```

**Fully green post-extension**, `internal/cli` included. No shipped fixture
changes behaviour.

**(iii) HOLDS.** What was compared: `models/rdr.toml` under `flow resolve`
at all three declared outcomes plus an unmodeled one, and under `flow next`,
byte-for-byte across the two binaries; plus the whole `go test ./...` suite
on the post-extension build. **Not reached:** no shipped fixture was diffed
payload-by-payload across the two binaries — the suite's green is the
evidence for those, so a fixture asserting only a subset of a payload would
not have caught a change in an unasserted field. `models/rdr.toml` is the
one model compared at full byte level.

## Both callers

`invokedReaders` is ONE function with exactly two call sites, confirmed:

```
internal/cli/flow_next.go:97     invokedReaders(req.model, "")
internal/cli/flow_resolve.go:98  invokedReaders(req.model, outcome)
```

`flow_next.go` passes `""` (union over ALL rows); `flow_resolve.go` passes
the `--outcome`. The added term sits inside the shared per-row loop, above
both, so it binds both by construction.

Measured directly, with a throwaway `internal/cli` probe test calling
`invokedReaders` at `outcome=""` (the `next` call) and at each declared
outcome (the `resolve` call), run against each arm and DELETED after:

PRE-extension:

```
mo.toml           outcome=""        readers=[ansreader]
mo.toml           outcome="go"      readers=[ansreader]
mo-breaking.toml  outcome=""        readers=[ansreader]
mo-breaking.toml  outcome="go"      readers=[ansreader]
dt-write.toml     outcome=""        readers=[nav]
dt-write.toml     outcome="locate"  readers=[nav]
../../models/rdr.toml outcome=""        readers=[rdr-status]
../../models/rdr.toml outcome="advance" readers=[rdr-status]
../../models/rdr.toml outcome="revise"  readers=[rdr-status]
../../models/rdr.toml outcome="abandon" readers=[rdr-status]
```

POST-extension:

```
mo.toml           outcome=""        readers=[ansreader modereader]
mo.toml           outcome="go"      readers=[ansreader modereader]
mo-breaking.toml  outcome=""        readers=[ansreader modereader]
mo-breaking.toml  outcome="go"      readers=[ansreader modereader]
dt-write.toml     outcome=""        readers=[nav]
dt-write.toml     outcome="locate"  readers=[nav]
../../models/rdr.toml outcome=""        readers=[rdr-status]
../../models/rdr.toml outcome="advance" readers=[rdr-status]
../../models/rdr.toml outcome="revise"  readers=[rdr-status]
../../models/rdr.toml outcome="abandon" readers=[rdr-status]
```

Every POST row is a SUPERSET of its PRE row; no row shrinks. Over
`mo.toml` at outcome `go` the two callers agree exactly — `[ansreader
modereader]` for both — confirmed independently from the CLI payloads:

```
$ ./post flow next    --model mo.toml             --artifact mode=mode.json --artifact ans=ans-full.json --as=json \
    | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["readers"])'
['ansreader', 'modereader']
$ ./post flow resolve --model mo.toml --outcome go --artifact mode=mode.json --artifact ans=ans-full.json --as=json \
    | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["readers"])'
['ansreader', 'modereader']
```

**Both-callers check HOLDS.**

## "changes no reader's contract"

Nothing in the edit touches `readerOutput`, `runReaders`, `RequestedKeys`,
or any accessor seam — it changes only WHICH readers the demand set names.
`go vet ./...` clean and the full suite green corroborate. A reader that
newly runs reports through the same `readers[]`/`owned{}` shape it always
did, visible in the POST payloads above.

## Verdicts

| Clause | Verdict |
| --- | --- |
| (i) added term empty over zero-owned-tag model | **HOLDS** (nearest-reachable model; see limit) |
| (ii) match-only-owned-key fixture, `next` + three `resolve` arms | **HOLDS** |
| (ii-b) BREAKING arm — plan becomes `flow-artifact-missing` exit 2 | **HOLDS** (breaking change confirmed live) |
| (iii) `models/rdr.toml` + shipped fixtures unchanged | **HOLDS** |
| both callers bound; strictly more, never fewer | **HOLDS** |

A15 is **not refuted in any clause**.

## Limits recorded

1. **A genuine zero-owned-tag model is unloadable on this build.** Carried
   over from the a5 spike and re-confirmed as the reason (i) is attested
   against `dt-write.toml`, which carries one dummy owned WRITTEN tag and
   zero owned MATCH atoms, rather than against a true 0010 decision table.
   The property (i) needs — the added term finds nothing when no row
   matches on an owned key — is fully exercised; the literal "zero owned
   tags" model is not constructible until 0010 ships.
2. **Shipped fixtures were verified by suite-green, not by payload diff.**
   Only `models/rdr.toml` was compared byte-for-byte across the two
   binaries. A test asserting a subset of a payload would not catch a
   change in an unasserted field.
3. **The refusing-reader arm's code is `flow-accessor-failed` (exit 3).**
   A15 says "the reader's own refusal (exit 3)" without naming the code;
   `flowbind`'s unreachable locator produces an EXECUTION FAILURE, so that
   is the exit-3 class observed. Another refusal mechanism (timeout,
   incomplete read) would take a different exit-3 code. The exit code and
   the "not a `flow-no-match`" property are what was verified.
4. **A pre-existing untracked scratch test**
   (`internal/cli/zz_scratch_a12_spike_test.go`, from an earlier spike)
   fails on BOTH arms for an unrelated reason and was moved aside during
   the suite runs, then restored. It is untracked and unchanged by this
   spike.

## Tree state

The edit was reverted and the pre-existing untracked file restored:

```
$ git checkout -- internal/cli/flow_exec.go
$ git status --porcelain
?? internal/cli/zz_scratch_a12_spike_test.go
```

No tracked file modified. The only new file is this evidence record.
