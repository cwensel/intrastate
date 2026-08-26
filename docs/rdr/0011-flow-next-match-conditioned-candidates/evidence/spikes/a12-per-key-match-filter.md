Model: claude-opus-5

# A12 — the match filter is per-KEY, and the probe's guard/owned verdict is `flow resolve`'s

A12 claims: "Omitting the absent-key match tags from a one-row probe changes
only whether the row reaches `gate`, and does so per KEY: the kernel filters
on `matches` before `gate`, `gate` evaluates `row.Guard` and `RequiresOwned`
independently of `row.Match`, and because every tag on one key shares that
key's presence the filter hands a key's tags over together or omits them
together — so a probe's guard verdict and owned-state disposition equal the
full-table `flow resolve`'s for that row, and a key carrying several match
tags (the `eq`+`in` pairing `atomsFromBlock` admits; `0002:C13`'s dead row
when the literals differ) is decided by the kernel's conjunction exactly as
`resolve` decides it, never split by the CLI."

Verdict up front: **(i) PASS, (ii) PASS, (iii) PASS, (iv) PASS**, with one
correction to the RDR's arithmetic recorded under (i) and one scope limit
recorded at the end.

As in A6, the second and third parts are verified against a **prototype** of
`0011:C1`, not against `flow next`: C1 is not implemented on this build —
`internal/cli/flow_next.go::excluded` still sets `probe.Match = nil`, so a
live `flow next` cannot show a match-driven exclusion at all. The prototype is
C1's probe builder verbatim (row's own match tags filtered to keys the
assembled view HAS, escape stripped, `Recognized: row.Outcome`), run through
the same `internal/resolve.Resolve` and the same `guardSeam()` the shipped
`excluded` uses.

Binary: `./bin/intrastate` at commit `1a7bddf`. No tracked file changed.

## The fixtures

`pair.toml` — one match key carrying BOTH `eq` and `in`:

```toml
outcomes = ["go"]

[model]
id = "pair"
version = 1
description = "One match key carrying both eq and in."

[tags.k]
provenance = "observed"
kind = "enum"
domain = ["x", "y", "z"]
single_valued = true

[tags.answer]
provenance = "owned"
kind = "enum"
domain = ["a1"]
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
id = "pairing"
[rule.match.k]
eq = "x"
in = ["x", "y"]
[rule.match.recognized]
eq = "go"
[rule.write]
answer = "a1"
```

`pairg.toml` is `pair.toml` plus an observed tag `g` (`domain = ["on","off"]`)
and a guard on the same rule, so the guard verdict can be varied
independently of the match:

```toml
[rule.guard.all.g]
eq = "on"
```

`nav.json` is `{}`.

## (i) — the pairing loads; the expansion is 2 rows, not 3

```
$ ./bin/intrastate flow next --model $S/pair.toml --artifact nav=$S/nav.json --as=json
{"type":"ok","data":{"model":".../pair.toml","revision":"","observed":{},"owned":{},"readers":["nav"],"outcomes":["go"],"candidates":[{"rule":"pairing","outcome":"go","required":["answer"],"unresolved":["k","answer"],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[]},{"rule":"pairing","outcome":"go","required":["answer"],"unresolved":["k","answer"],"next":{"answer":"a1"},"writes":{"answer":"a1"},"clear":[]}]}}
EXIT=0
```

**No loader tier refuses the pairing.** `atomsFromBlock` iterates the
operators under one key (`for _, op := range slices.Sorted(maps.Keys(predicates))`)
and emits an atom per operator; `eq` and `in` are both admitted under
`BlockMatch`, and nothing downstream rejects a key carrying both. Exit 0, two
candidate entries, both `"rule":"pairing"`.

`lint` reports one finding, and it is not about the pairing:

```
$ ./bin/intrastate lint --model $S/pair.toml --as=json
{"code":"graph-lint-failed","message":"the model carries blocking graph-lint findings","findings":[{"code":"graph-dangling-edge","message":"the model declares no initial owned state; add an `[initial]` table assigning every always-present owned tag","model":"pair","severity":"blocking","element":"model"}]}
EXIT=2
```

That is the same `graph-dangling-edge` the A5 spike recorded for a fixture
with no `[initial]` table, and `lint` does not gate `flow next`.

The normalized rows, from the prototype (`A12_MODEL=$S/pair.toml`):

```
== (i) expansion of eq="x" + in=["x","y"] ==
rows: 2
row[0] rule=pairing outcome=go suffix=[x] matchTags=1 [{k x}] guardAtoms=0 requiresOwned=[answer]
row[1] rule=pairing outcome=go suffix=[y] matchTags=2 [{k x} {k y}] guardAtoms=0 requiresOwned=[answer]
```

and with the guard (`A12_MODEL=$S/pairg.toml`):

```
rows: 2
row[0] rule=pairing outcome=go suffix=[x] matchTags=1 [{k x}] guardAtoms=1 requiresOwned=[answer]
row[1] rule=pairing outcome=go suffix=[y] matchTags=2 [{k x} {k y}] guardAtoms=1 requiresOwned=[answer]
```

**A correction to A12's parenthetical.** A12 speaks of "a key carrying several
match tags" as though the pairing produced a two-tag row in every case. It
does not, uniformly:

- **row `[x]` carries ONE tag**, not two. `expand` turns the `in` member `"x"`
  into an `eq` atom on key `k` with literal `"x"` — which is byte-identical to
  the authored `eq = "x"` atom, so the two collapse to one entry in
  `Row.Atoms`. That is `0002:C13`'s own rule ("`eq = "x"` and `in = ["x"]` are
  one spelling of one edge and cannot mint two identities") applying to the
  member that coincides with the `eq` literal.
- **row `[y]` carries TWO tags on key `k`** — `{k x}` from the authored `eq`
  and `{k y}` from the `in` member — and is the dead row: no single value of
  `k` satisfies both.

So the pairing yields **exactly 2 expanded rows** (one per `in` member), of
which **1 is live** (`[x]`) and **1 is dead** (`[y]`), and only the dead one
exercises the several-tags-on-one-key case. A12's substantive claim is
unaffected — the dead row is the one it is about — but the "`eq`+`in` pairing"
does not double every row's tag count, and a spike or test asserting two tags
on the live row would be asserting something false.

## (ii) — the dead row, live over the prototyped C1 probe

From `A12_MODEL=$S/pair.toml` (no guard, so nothing but the match moves):

```
== view: k ABSENT, g ABSENT  (view keys = [answer]) ==
row[0] suffix=[x] kept=[] absentMatchKeys=[k] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[] absentMatchKeys=[k k] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE

== view: k=x, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=plan:pairing suffix=[x]        guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=plan:pairing suffix=[x]        guard/owned AGREE

== view: k=y, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE

== view: k=z, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
```

Read off the dead row (`row[1]`, suffix `[y]`):

- **`k` ABSENT** — `kept=[]`, `absentMatchKeys=[k k]` (BOTH of the key's tags
  reported absent, which is the `{k, absent}` C1 would report once, deduped);
  `C1probe=plan:pairing`. **The dead row is a candidate while `k` is absent.**
- **`k` present at ANY value** — `k=x`, `k=y`, `k=z` all give
  `C1probe=refused:no_match` on the dead row. At `k=x` the tag `{k y}` fails;
  at `k=y` the tag `{k x}` fails; at `k=z` both fail. The exclusion is the
  kernel's `TagSet.matches` conjunction over the row's tags — **no CLI literal
  comparison exists**: the prototype's only decision is presence
  (`if _, present := view[t.Key]; present`), never value.

`z` is worth keeping: it is a value the author never wrote on either side of
the pairing, and the dead row is excluded there too, which is what "excluded
once `k` is present at any value" means.

## (iii) — the probe's guard verdict and owned disposition equal `flow resolve`'s

The guarded fixture varies the guard (`g` absent / `off` / `on`) and the owned
snapshot (`answer` supplied / withheld) independently of the match key, over
both expanded rows. 28 row×view observations. `pureGuardOwned` is the SAME
probe with `Match` set to `nil` — i.e. today's shipped `excluded()`, and the
row's pure guard/owned verdict with match removed entirely. The check applied
mechanically: *whenever the C1 probe is not refused BY the match filter, its
disposition must equal `pureGuardOwned`.*

```
$ A12_MODEL=$S/pairg.toml go test ./internal/cli -run TestScratchA12 -v
== view: k ABSENT, g ABSENT  (view keys = [answer]) ==
row[0] suffix=[x] kept=[] absentMatchKeys=[k] | C1probe=refused:guard_unevaluable    pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[] absentMatchKeys=[k k] | C1probe=refused:guard_unevaluable    pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE

== view: k=x, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:guard_unevaluable    pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:guard_unevaluable      guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:guard_unevaluable      guard/owned AGREE

== view: k=y, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE

== view: k=z, g ABSENT  (view keys = [answer k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:guard_unevaluable    fullResolve=refused:no_match               guard/owned AGREE

== view: k ABSENT, g=off (guard FALSE)  (view keys = [answer g]) ==
row[0] suffix=[x] kept=[] absentMatchKeys=[k] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[] absentMatchKeys=[k k] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE

== view: k=x, g=off (guard FALSE)  (view keys = [answer g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE

== view: k=y, g=off (guard FALSE)  (view keys = [answer g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:no_match             fullResolve=refused:no_match               guard/owned AGREE

== view: k ABSENT, g=on  (view keys = [answer g]) ==
row[0] suffix=[x] kept=[] absentMatchKeys=[k] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[] absentMatchKeys=[k k] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE

== view: k=x, g=on  (view keys = [answer g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=plan:pairing                 pureGuardOwned=plan:pairing                 fullResolve=plan:pairing suffix=[x]        guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=plan:pairing suffix=[x]        guard/owned AGREE

== view: k=y, g=on  (view keys = [answer g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE

== view: k=z, g=on  (view keys = [answer g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=plan:pairing                 fullResolve=refused:no_match               guard/owned AGREE

== view: k ABSENT, g=on, OWNED MISSING  (view keys = [g]) ==
row[0] suffix=[x] kept=[] absentMatchKeys=[k] | C1probe=refused:owned_state_unavailable pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[] absentMatchKeys=[k k] | C1probe=refused:owned_state_unavailable pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:no_match               guard/owned AGREE

== view: k=x, g=on, OWNED MISSING  (view keys = [g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:owned_state_unavailable pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:owned_state_unavailable guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:owned_state_unavailable guard/owned AGREE

== view: k=y, g=on, OWNED MISSING  (view keys = [g k]) ==
row[0] suffix=[x] kept=[{k x}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:no_match               guard/owned AGREE
row[1] suffix=[y] kept=[{k x} {k y}] absentMatchKeys=[] | C1probe=refused:no_match             pureGuardOwned=refused:owned_state_unavailable fullResolve=refused:no_match               guard/owned AGREE
--- PASS: TestScratchA12PerKey (0.00s)
```

```
$ grep -c 'DISAGREE\|SPLIT' guarded.out
0
```

**Zero disagreements across all 28 observations.** Every difference between
`C1probe` and `pureGuardOwned` is the match filter turning the disposition
into `no_match`, and nothing else: the guard never flips between
`guard_unevaluable`, `no_match` (guard FALSE prunes to zero survivors), and
`plan`; the owned refusal never appears or disappears because of the match.

Against the full table, per row, the three dispositions the RDR names line up:

- **guard undecidable** (`g` absent, `k=x`): probe `guard_unevaluable`,
  `fullResolve` `guard_unevaluable`. Same.
- **guard FALSE** (`g=off`): probe `no_match` (the row is pruned at `gate`,
  leaving zero survivors), `fullResolve` `no_match`. Same — and this is the
  disposition today's `excluded()` acts on.
- **owned missing** (`k=x, g=on`, no `answer`): probe
  `owned_state_unavailable`, `fullResolve` `owned_state_unavailable`. Same.
- **live** (`k=x, g=on`, `answer` supplied): probe `plan:pairing`,
  `fullResolve` `plan:pairing suffix=[x]` — and `fullResolve` selects the LIVE
  row `[x]`, not the dead `[y]`, exactly as the probe decides row-by-row.

Note `fullResolve` is a whole-table verdict for the recognized outcome, so it
is the same string on both rows of a view; the per-ROW comparison that A12
actually asserts is `C1probe` vs `pureGuardOwned`, which is the column pair
checked mechanically above.

## (iv) — the filter is per KEY, and a fully-filtered row reaches `gate`

Two mechanical checks, both clean:

**No split.** The prototype asserts, for every row and view, that a key's tags
are all kept or all dropped — `if kn := kept[k]; kn != 0 && kn != n` flags a
partial keep. `grep -c 'SPLIT' = 0` over both fixtures. The dead row's two
tags on `k` are visible moving together in the output: `kept=[]` with
`absentMatchKeys=[k k]` when `k` is absent, `kept=[{k x} {k y}]` when `k` is
present. **Never one of two.** This is structural, not incidental: the filter
predicate reads only `t.Key`, and both tags carry the same key, so the
membership test returns the same answer for both by construction.

**A zero-tag match matches unconditionally and reaches `gate`.** Directly, over
a hand-built row with `Match: nil`:

```
zero-tag match row, view={answer:a0}: err=<nil> refused=false plan=&{RuleID:zero SourceLocator: NextTags:[] Writes:[{Key:answer Value:a1}] Revision: Escaped:false}
```

`TagSet.matches` is a `for` over `want` returning `true` on an empty slice, so
the row survives the candidacy filter and is handed to `gate`. The guarded
fixture confirms the consequence: at `k ABSENT, g ABSENT` both rows have
`kept=[]` and BOTH answer `guard_unevaluable` — a refusal only `gate` can
produce, so both rows demonstrably reached it.

## Prototype probe code

The scratch test lived at `internal/cli/zzscratch_a12_perkey_test.go` and was
deleted after the run (`git status` clean). Its C1 probe builder, verbatim:

```go
// c1Probe prototypes 0011:C1's probe builder: the row's own match tags,
// filtered to keys the assembled view HAS; escape stripped; Recognized
// bound to the row's outcome (bound at the Resolve call below).
func c1Probe(row table.Row, view map[string]string) resolve.Row {
	probe := row.KernelRow()
	kept := make([]resolve.Tag, 0, len(probe.Match))
	for _, t := range probe.Match {
		if _, present := view[t.Key]; present {
			kept = append(kept, t)
		}
	}
	probe.Match = kept
	probe.Escape = nil
	return probe
}

func a12Run(row table.Row, probeRow resolve.Row, owned []resolve.Tag, observed []resolveTag) string {
	res, err := resolve.Resolve(resolve.Input{
		Table: resolve.Table{
			Outcomes: []string{row.Outcome},
			Rows:     []resolve.Row{probeRow},
		},
		Owned:      owned,
		Observed:   kernelTags(observed),
		Recognized: row.Outcome,
		Guards:     guardSeam(),
	})
	if err != nil {
		return "ERR " + err.Error()
	}
	if res.Refused() {
		return "refused:" + string(res.Refusal.Kind)
	}
	return "plan:" + res.Plan.RuleID
}
```

The per-key and agreement assertions:

```go
// (iii)/(iv): whenever the C1 probe is not refused BY the match filter,
// its disposition must equal the pure guard/owned verdict.
agree := "guard/owned AGREE"
if probeVerdict != baseline && probeVerdict != "refused:no_match" {
	agree = "*** DISAGREE ***"
}
// (iv): per-key atomicity — a key's tags are kept together or dropped
// together, never split.
have, kept := map[string]int{}, map[string]int{}
for _, tg := range kr.Match {
	have[tg.Key]++
}
for _, tg := range pr.Match {
	kept[tg.Key]++
}
split := ""
for k, n := range have {
	if kn := kept[k]; kn != 0 && kn != n {
		split = fmt.Sprintf(" *** SPLIT key=%s kept %d of %d ***", k, kn, n)
	}
}
```

The source half of the claim, for the record:

- `internal/resolve/resolve.go::Resolve` filters candidacy on
  `if !view.matches(row.Match) { continue }` and only then calls
  `gate(candidates, in.Guards, view)` — match strictly before gate.
- `internal/resolve/resolve.go::gate` reads `row.Guard` (via `evaluateAtoms`)
  and `row.RequiresOwned`. It never reads `row.Match`.
- `internal/resolve/resolve.go::TagSet.matches` is a conjunction over `want`;
  an absent or conflicted key, or an unequal value, is `false`. Over zero tags
  it returns `true`.
- `internal/table/normalize.go::atomsFromBlock` emits one atom per operator
  under a key, so `eq` + `in` on one key is two atoms; `expand` turns each
  `in` member into an `eq` atom carrying the same key.

## Limits

1. **C1 is unimplemented, so (ii) and (iii) are not attestable through
   `flow next`.** Everything above the `flow next`/`lint` runs in (i) comes
   from a prototype of C1's probe builder, not from shipped CLI behaviour.
   Today `excluded()` sets `probe.Match = nil`; the `pureGuardOwned` column is
   literally today's build, and it over-reports the dead row as a candidate at
   `k=y` and `k=z` where C1 would drop it. A12 describes the POST-C1 world.

2. **The `{k, absent}` report itself is not exercised.** The prototype computes
   the absent-key set (`absentMatchKeys`) but nothing in this spike renders it
   into a `unresolved`/`absent` payload — that surface does not exist yet. The
   spike establishes that both of a key's tags land in that set together
   (`[k k]`, deduping to one `k`), not what the eventual payload looks like.

3. **Single-key, single-valued, observed.** The fixture pairs `eq`+`in` on one
   observed enum tag. A `set_valued` tag, an OWNED match key, or two keys each
   carrying multiple tags were not exercised. The per-key argument is
   structural (the filter reads only `t.Key`), so the risk is low, but this
   spike does not measure those cases. In particular a CONFLICTED key —
   present but carrying no single value, which `TagSet.matches` treats as
   non-matching — is present-for-filter-purposes and therefore hands its tags
   to the kernel; that path is untested here.

4. **A concurrent sibling spike was editing this working tree.** During the
   run `internal/cli/flow_exec.go` briefly carried an uncommitted `A15 spike`
   edit (`matchOwnedKeys`) from another agent, and a scratch test file was
   swept mid-run. The edit did not touch `resolve`, `normalize`, or
   `excluded`, and the guarded sweep above was re-run to completion after
   recreating the scratch file. `git status` is clean at the end of this
   spike apart from this evidence file.
