# Critique dual-model diff — RDR cli/0011 (iter-2)

Pass A: `critique.md` (Model: claude-opus-5, 12 rows)
Pass B: `critique-modelB.md` (Model: claude-sonnet-5, 9 rows)

Reconciled by PASSAGE ANCHOR — `C-N` ids do not correspond across files.
Barrier: both passes landed before this diff was written; this diffing
context authored neither (rdr-common §auto-fanout).

---

## Agreement (both models, same passage)

| Passage anchor | Pass A | Pass B | What both say | Depth |
|---|---|---|---|---|
| `0011:A15` Status `Pending` + `0011:C1`'s demand-set MUST | C-2, C-3 | C-1 | The demand-set extension is the RDR's riskiest mechanism and its backing assumption is unverified at lock time, while normative prose depends on it as settled fact. | **A deeper.** B stops at "unverified, and a match-only owned key silently degrades to `--all`" (the *If wrong* arm the RDR already writes). A goes past it to the mechanism B missed: `invokedReaders` has a SECOND caller in `flow_resolve.go`, so the widening changes a different verb. |
| `0011:C2` `--all` strips `BlockMatch` + `0011:F1` "diagnose with `--all`" + `0011:D-naming`'s conceded asymmetry | C-8, C-10 | C-4, and named as B's §2 "rewritten within 6 weeks" | `--all` is prescribed as THE diagnostic for a missing candidate but is contractually forbidden from naming why; the remedy degrades to eyeballing the model TOML, which an automated skill caller cannot do. | **Even.** A adds the migration-count argument (no invocation reproduces 0005's payload — four migrations, not one). B adds the automatability argument (a skill caller cannot open a TOML). |
| `0011:§consequences` / `0011:C3` — the `unresolved`→`unknown` rename bundled with the predicate change | C-10 | C-5, premortem Ticket 4 | The wire-format rename is an independent axis bundled into the same release with no decoupling path. | **B deeper on the consumer axis** (out-of-repo parsers, `intrastate` is by design a generic library, so out-of-repo consumers structurally exist and the in-repo census cannot see them). A deeper on the *count* (the consumer owes four changes, and Consequences names two). |
| `0011:§metadata` Profile `large` / `0011:§proportionality` — record size vs. change size | C-11 | C-2, and the closing paragraph of B's premortem | ~800 lines, three contracts, 15 assumptions for a change whose code delta is one function's control flow + one demand-set term + a field rename. | **Even, and both are soft.** A frames it as the gate's own split test (two seams, one of them `Pending`); B frames it as implementer working-memory overflow. Neither is independently load-bearing. |

**Hotspot note.** The `0011:A15` / demand-set passage is the one place BOTH isolated passes independently landed, and the one place my verification found a claim the RDR's own prose contradicts. That is agreement functioning as signal, not as validation — the two passes agreed on the *location* and disagreed on the *defect*, and only one of them found the live one.

---

## A-only findings

| ID | Passage | Claim | Verdict (see §Independent verification) |
|---|---|---|---|
| A C-1 / C-12 | `0011:C1` "a single probe row holds at most one match tag per key"; `0011:A12` "no cross-atom interaction" | `atomsFromBlock` loops over OPERATORS within a key, so one `[rule.match.k]` block with both `eq` and `in` yields TWO match tags on key `k`; C1's per-key filter then makes the predicate non-monotonic in supplied state. | **CONFIRMED (reproduced independently)** |
| A C-2 | `0011:§approach`, `0011:C1` closing, `0011:JC1` | `invokedReaders` has a second caller at `flow_resolve.go:98`; widening the demand set widens `flow resolve`'s invoked reader set, and `runReaders` has two refusal arms on that set. | **CONFIRMED** |
| A C-6 | `0011:§problem-statement` / `0011:A6` (the 21→3 witness) | The 3 survivors are exactly one per declared outcome — the alphabet `outcomes[]` already carried, unconditionally, since 0005. Generalizes across `stage=`. | **CONFIRMED, and generalizes further than A claimed** |
| A C-9 | `0011:A3` | `checkAccessorBindings` fires at 0 readers too, not only `!= 1` as A3's quoted message implies. | **PARTIAL — mechanically true, but A3 is not actually wrong** |
| A C-5 | `0011:C3` "FIVE TEST reads" | `:103` is a presence probe not a `stringsAt` read; and the census leaves ~15 comment/message occurrences of `unresolved` uninventoried. | **PARTIAL — first half REFUTED as stated, second half CONFIRMED** |
| A C-4 | `0011:A4` "BREAKS 0" | `TestReq44` iterates every candidate requiring `flag`; survives only because both `flowGatedNextModel` rows guard on `flag`. Any added row without a `flag` guard breaks it. | **CONFIRMED (the coupling), REFUTED (the "cited beyond its scope" half)** |
| A C-7 | `0011:D-undecided-reporting-shape` | `respond/text.go::flatten` emits one line per leaf, so `{key, reason}` doubles text output on 0010's no-tag decision table. | **NOT INDEPENDENTLY RE-VERIFIED** (outside my priority list; A reports reading `flatten`, the mechanism is consistent with the RDR's own Load-Bearing Decision, which concedes the flattener behaviour and declines to act) |

## B-only findings

| ID | Passage | Claim | Assessment |
|---|---|---|---|
| B C-3 | `0011:A11` | A11 is true today but is a coincidence of the current call graph (one `owned` slice handed to both `assembledView` and `excluded`); nothing structurally enforces it, and no oracle in S1–S8 fails specifically because A11 broke. | **CONFIRMED as a durability gap.** Verified: `runFlowNext` builds `owned` once (`flow_next.go:97`) and passes the same slice to `assembledView` (`:104`) and `excluded` (`:123`). S5 asserts the three *load/parse* pins A3 rests on and `git diff --stat internal/resolve` empty — it does NOT assert the assembledView-vs-`assemble` key-set agreement. B's proposed AT (a direct key-set property test) is a real, cheap addition. |
| B C-6 | `0011:C1` "never by rebuilding match tags from `row.Atoms`" | The filter-after-`KernelRow()` constraint exists only as RDR prose; no guard comment and no test distinguishes a compliant build from one that filters `row.Atoms` before conversion. Set-kinded match keys would silently mis-compare. | **CONFIRMED as a durability gap.** `Row.KernelRow` (`internal/table/model.go:232-264`) routes `BlockMatch` atoms through `seamValue(a.Literal, r.isSet(a.Key))`; the set-canonicalization is real and re-deriving it CLI-side would fork. No S1–S8 scenario uses a set-kinded match key, so the two implementations are observationally identical across the mandated fixtures. |
| B C-8 | `0011:§minimum-viable-validation` step 6 / S3 | No fixture exercises MULTIPLE match keys on the SAME row in mixed states (one present-equal, one present-unequal, one absent). A12's per-atom independence is asserted, not tested against a multi-key row. | **CONFIRMED.** S3 is specified as "the three reachable match cases, one row each"; S4 covers all-absent; nothing covers a mixed multi-key row. This is the exact fixture that would have caught A C-1 — a multi-key-on-ONE-key fixture is one step away from a two-atoms-on-one-key fixture. B found the fixture-shaped hole without finding the defect underneath it. |
| B C-7 | `0011:§technical-design` "the one signature change in the file" | The `bool` → `Result` change on `excluded` is a bigger refactor than framed: three refusal kinds plus a success path where today there is one boolean, with A13's ordering constraint. | **PARTIAL.** The mechanism is real — `excluded` today (`flow_next.go:276-295`) returns `result.Refusal.Kind == resolve.KindNoMatch` and discards the `Result`; C1 requires `summarize` to additionally read `Refusal.Undecided`. But this is implementation risk, not a defect in the record: S7 exists precisely to pin the `owned_state_unavailable`-without-payload vs `guard_unevaluable`-with-payload precedence, and the RDR states the limit rather than hiding it. B's own AT for this is a reasonable S7 sharpening, not a route-back. |
| B C-9 | `0011:§decision-rationale`, "O3 wins on the two deciding rows" | The permanent `flow next` / `flow resolve` disagreement on an absent-match-key row is a structural split the "report-vs-select" analogy understates. | **RE-RAISE — see Unhealthy-pass check.** |

---

## Independent verification

### 1. Pass A C-1 / C-12 — two match tags on one key — **CONFIRMED**

I authored my own probe model and harness against the shipped `internal/table`,
without reading Pass A's harness. Model: one owned enum tag `k` with domain
`["a","b","c"]`, one reader, one writer, one rule authoring BOTH operators under
one match block:

```toml
[rule.match.k]
eq = "a"
in = ["a", "b"]
```

The mechanism is at `internal/table/normalize.go::atomsFromBlock` (`:23-40`) —
an outer loop over keys and an **inner loop over the operators declared under
that key**:

```go
for _, key := range slices.Sorted(maps.Keys(block)) {
    ...
    predicates := block[key]
    for _, op := range slices.Sorted(maps.Keys(predicates)) {
        atom, err := l.atom(decl, key, op, predicates[op], b, owner)
        ...
        out = append(out, atom)
    }
}
```

`internal/table/normalize.go::atom` (`:77`) refuses any operator other than
`eq`/`in` under a match block. Nothing refuses BOTH on one key.

Observed (`table.Load` on the shipped build, then `Row.KernelRow()`):

```
LOAD OK, rows = 2
row r1 outcome=advance
  atom key=k block=match op=eq lit=[a]
  kernel Match=[{k a}]
row r1 outcome=advance
  atom key=k block=match op=eq lit=[a]
  atom key=k block=match op=eq lit=[b]
  kernel Match=[{k a} {k b}]
  >>> KEY "k" CARRIES 2 MATCH TAGS ON ONE ROW
```

**C1's stated soundness condition — "a single probe row holds at most one match
tag per key" — is FALSE on the shipped normalizer.** So is `0011:A12`'s
"no cross-atom interaction". A12's Evidence reasons correctly about the operator
whitelist and correctly about `0002:C13`'s row-level `in`-expansion, and never
asks whether one key can carry two atoms.

**What C1's filter then does.** I built the probe exactly as C1 mandates
(`probe := row.KernelRow()`, drop every `resolve.Tag` whose `Key` the view
lacks, `probe.Escape = nil`, `Recognized: row.Outcome`) and ran the shipped
`resolve.Resolve`:

```
k ABSENT         probe.Match=[]              =>  owned_state_unavailable -> CANDIDATE
k PRESENT=a      probe.Match=[{k a} {k b}]   =>  no_match -> EXCLUDED (unreported)
k PRESENT=b      probe.Match=[{k a} {k b}]   =>  no_match -> EXCLUDED (unreported)
k PRESENT=c      probe.Match=[{k a} {k b}]   =>  no_match -> EXCLUDED (unreported)
```

The inversion is exactly as Pass A describes, and it is produced BY C1's own
mechanism: C1's clause "A row ALL of whose match atoms are omitted is probed
with an empty match pattern, which matches unconditionally; it is a candidate"
is the clause that converts an unsatisfiable row into an accepted one.
Supplying state removes a row that absence kept — the inverse of
`0011:§approach`'s framing. `no_match` is unreported by construction, so no
`unknown` entry explains the loss, and `0011:C2` strips the match facts under
`--all`, so `--all` names nothing.

**Is such a row AUTHORABLE per the loader's validation?** Yes, with one
qualification Pass A did not check and which I add here because verification
required it:

- `table.Load` **accepts** the model outright — no refusal (output above).
- `graphlint.Run` emits exactly one finding: `graph-unreachable-rule`
  ("no reachable owned-state satisfies the selection context of row `r1`, so
  the row can never be a candidate") — and `graphlint.IsBlocking` returns
  **false** for it (`internal/graphlint/taxonomy.go:112`). `internal/cli/lint.go:116`
  gates on `report.Blocking()`, so the model **passes the blocking lint gate**.

So the row is authorable, loadable, and ships past `intrastate lint`'s blocking
tier with an advisory the author is free to ignore — which is precisely the
population in which C1's non-monotonicity would be discovered in the field
rather than at authoring time. Today the row is harmless: `excluded` sets
`probe.Match = nil` (`flow_next.go:278`), and `flow resolve` folds it to
`no_match` correctly.

Two honest limits on this finding, stated so the fix half does not over-correct:
(i) the row IS an authoring error and lint already flags it, advisory-tier;
(ii) the same two-tags-on-one-key shape is reachable from two `eq`-only blocks
only if TOML permits a duplicate key, which it does not — the reachable
authoring is `eq` + `in` under one block, as tested. Neither limit rescues C1's
stated soundness condition, which is a claim about what the normalizer can
produce, not about what a careful author would write.

### 2. Pass A C-2 — `invokedReaders`' second caller — **CONFIRMED**

```
$ grep -rn 'invokedReaders(' internal/
internal/cli/flow_next.go:97:	readers, owned, ce := req.runReaders(cmd.Context(), invokedReaders(req.model, ""))
internal/cli/flow_resolve.go:98:		invokedReaders(req.model, outcome))
internal/cli/flow_exec.go:86:func invokedReaders(m *table.Model, outcome string) []string {
```

Two callers, exactly as Pass A states. `flow next` passes `outcome == ""`
(every row contributes); `flow resolve` passes a concrete outcome, so
`invokedReaders` skips rows whose `Outcome != outcome` (`flow_exec.go:89-105`)
— but the surviving rows' match-block owned keys would join the demand set
under C1's MUST, and the demand set is contractually mode- and verb-independent
("The demand set is a property of the MODEL, not of the mode", C1;
A15 item (iii): "the widening is uniform rather than `next`-special-cased").

`internal/cli/flow_exec.go::runReaders` (`:175-205`) has both refusal arms
Pass A names, and the role pre-check runs over the whole invoked set BEFORE any
accessor executes:

```go
for _, name := range names {
    def, ok := r.registry.Lookup(name, accessor.CapRead)
    if !ok { continue }
    if _, bound := r.artifacts[def.Accessor.Role]; !bound {
        return nil, nil, userErr(codeArtifactMissing, def.Accessor.Role, ...)
    }
}
...
result := exec.Read(ctx, name)
if result.Refused() {
    return nil, nil, accessorFailure(*result.Refusal, phaseRead)
}
```

So a `flow resolve --outcome X` that succeeds today can refuse
`flow-artifact-missing` (exit 2) or on reader refusal (exit 3) after the
extension, over a `state-machine`-class model with a match-only owned key on a
row bound to X.

**On the RDR's wording.** `0011:C1`'s closing paragraph reads:
"This contract changes nothing in `internal/resolve`: the kernel's match seam
stays two-valued and `flow resolve`'s dispositions are unchanged in every case."
The first clause is TRUE and verified. The second clause is a claim about the
VERB, and the verb's reader-invocation set moves. `0011:§approach`'s "The change
lands in the **CLI**, and `internal/resolve` is not touched" is likewise true
about the package and silent about the second caller. `0011:A15` item (ii)
concedes the refusal-change in its own text ("now fails where it previously
succeeded … but is a behaviour CHANGE"), and `0011:§consequences` does not carry
it forward — its Negative bullets name the `--all` migration, the payload
rename, the next/resolve absent-key disagreement, and the A3 dependency, and
name no `flow resolve` reader-set change.

`0011:JC1` reasons entirely about 0010's zero-owned-tag `decision-table` class
("the added term is empty over that class") and concludes "Not a collision" —
correct for that class, and silent on every `state-machine`-class model.

### 3. Pass A C-6 — the 21→3 witness reproduces the alphabet — **CONFIRMED, and generalizes further**

I ran C1's probe shape over `models/rdr.toml` for EVERY declared value of
`stage`, with `status=draft`, `gate_passed=false`:

```
outcomes = [advance revise abandon]

stage=seeded       cands=3 [propose seed-abandon seed-revise]                                  distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=proposed     cands=3 [propose-abandon propose-again refine]                              distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=refined      cands=3 [refine-abandon refine-again resolve-assumptions]                   distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=resolved     cands=3 [prelock resolve-abandon resolve-route-back]                        distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=prelocked    cands=3 [prelock-abandon prelock-route-back reconcile]                      distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=reconciled   cands=4 [finalize-blocked finalize-pass reconcile-abandon reconcile-route-back]  distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=final        cands=3 [final-abandon final-route-back implement]                          distinct-outcomes=[abandon advance revise]  ==alphabet? true
stage=implemented  cands=0 []                                                                  distinct-outcomes=[]                        ==alphabet? false
stage=dropped      cands=0 []                                                                  distinct-outcomes=[]                        ==alphabet? false
```

`stage=resolved` matches `0011:A6` exactly (3 candidates: `prelock`,
`resolve-route-back`, `resolve-abandon`) and matches the A6 spike at
`evidence/spikes/a6-rdr-toml-narrowing.md`. **At every non-terminal stage the
distinct candidate outcomes equal the full declared alphabet.** The two
terminals yield zero candidates (genuinely informative — but they are terminals,
where the caller already knows there is nothing to do).

The A6 spike's own captured baseline payload prints
`"outcomes": ["advance", "revise", "abandon"]` in the SAME object as the 21
candidates — so `outcomes[]` was already carrying the answer the narrowed
candidate list reproduces, and `0011:C1` requires it keep doing so
("`outcomes` MUST remain the model's full declared alphabet in every mode").

**Scope of this finding, stated precisely.** It does NOT refute A6, whose
narrowing claim is arithmetically correct. It does not refute C1 either. It is a
finding against the PROBLEM STATEMENT's use of A6 as the motivating witness:
the RDR's stated success criterion ("no row is listed that the supplied state
excludes") is met, while the caller's stated need ("so the next action can be
chosen") is not advanced on the flagship model — the rule ids are new
information, the outcome set is not. The RDR half-anticipates this
(`0011:§problem-statement`: "NOT a single next action"; "a caller facing several
candidates is being told the state genuinely admits several") and that defence
is sound in principle. What is missing is the record acknowledging that its
headline metric measures removal of rows for OTHER stages, which is real but is
not what `21→3` is being read to mean.

### 4. Pass A C-9 — `checkAccessorBindings` at 0 readers — **PARTIAL**

Mechanically CONFIRMED. `internal/table/load.go::checkAccessorBindings`
(`:360-383`):

```go
case ProvenanceOwned:
    if readerCount[key] != 1 {
        return fail(CatMalformedAccessorBinding,
            fmt.Sprintf("owned tag %s is served by %d readers; want exactly one", key, readerCount[key]))
    }
```

`readerCount` is a map with zero default, so the arm fires at 0. Verified live
by removing the `[read.k-reader]` block from my probe model:

```
zero-reader owned tag load err = malformed_accessor_binding: owned tag k is served by 0 readers; want exactly one
```

**But A3 is not wrong.** A3's Evidence says `checkAccessorBindings` "refuses at
MODEL LOAD any model where an owned tag is served by `!= 1` readers" — `!= 1`
INCLUDES 0, and A3 quotes the format string with `%d`, not a hardcoded `2`.
Pass A's ledger row ("A3's proof rests on … refusing `!= 1` readers … but never
states the zero arm") reads a gap that the quoted predicate already covers.
`0011:C1` also states the zero-arm consequence explicitly in its own text:
"`checkAccessorBindings` does not prevent this — it requires every DECLARED
owned tag to have exactly one reader, so the reader exists and is simply never
invoked."

Residual value: the implementer authoring S8's fixture needs to know the key
must declare a reader (and, per `checkAccessorBindings`' writer walk, that an
owned key nothing writes owes at most one writer). That is a one-clause
clarification to S8, not a defect in A3. **Severity: minor / editorial.**

### 5. Pass A C-5 — the five-read census — **PARTIAL (first half REFUTED, second half CONFIRMED)**

`internal/cli/flow_next_0005_test.go:103` is indeed a presence probe:

```go
if _, present := c["unresolved"]; !present {
    t.Errorf("candidate[%d] carries no unresolved guard/gate facts", i)
}
```

**But Pass A's ledger row claims C3 "counts it among the three needing ok-bool
treatment". That is REFUTED.** `0011:C3` names the three ok-bool sites
explicitly and `:103` is not among them: "THREE of the five discard the comma-ok
of a helper (`flow_harness_0005_test.go::stringsAt` …) — `flow_next_0005_test.go:404`,
`flow_mvv_0005_test.go:66`, and `flow_adversarial_0005_test.go:446`". C3
separately and correctly labels `:103` "the presence check `:103`". The census
is accurate as written; the only quibble left is whether a presence probe should
be called a "read", which is a word choice, not a defect.

**The second half CONFIRMED.** `grep -rn 'unresolved' internal/cli/` returns 33
occurrences (30 in tests, 7 in `flow_next.go` — overlapping counts by line).
Beyond the five reads C3 inventories, `unresolved` appears in COMMENT text at
`flow_next_0005_test.go:23,66,135,137,141,383` and `flow_mvv_0005_test.go:46`,
and in `t.Errorf`/`t.Fatalf` MESSAGE text at `flow_next_0005_test.go:104,165,169,171,407,408`,
`flow_mvv_0005_test.go:73`, and `flow_adversarial_0005_test.go:448,450`. C3's
caveat "The count is of test reads only" is honest about the production sites
(which it then enumerates in detail) but does not inventory the comment and
message prose, which goes factually false the day the field is renamed and which
a green suite will never flag. **Severity: minor, resolvable in-draft** — one
clause in C3 requiring the same mechanical sweep over comment/message text.

### 6. Pass A C-4 — A4's "BREAKS 0" scope — **CONFIRMED (coupling), REFUTED (mis-citation)**

The coupling is real. `TestReq44_NextInventsNoGuardFactAbsentFromEveryChannel`
(`internal/cli/flow_next_0005_test.go:384-411`) iterates **every** candidate and
requires each to list `flag`:

```go
for _, c := range candidates {
    unresolved, _ := stringsAt(c, "unresolved")
    if !slices.Contains(unresolved, "flag") {
        t.Errorf("candidate %v is guarded on `flag` …", c["rule"])
    }
}
```

It survives only because both rows of `flowGatedNextModel`
(`flow_fixtures_0005_test.go:304-326`) carry `[rule.guard.all.flag]`. Any row
added to that model without a `flag` guard breaks it. This coupling is
undocumented, and `0011:S2` names `flowGatedNextModel` as a target fixture
("`flowGatedNextModel` and `flowMVVModel` … plus the discriminating row C3
requires") without saying which model the discriminating row lands in.

**But Pass A's claim that A4's "BREAKS 0" is "cited in `0011:§overrides` and
`0011:§risks-and-mitigations` as if it cleared the suite" is REFUTED.** A4's own
Evidence carries the scoping caveat verbatim: "This classification is
predicate-scoped and does NOT clear the suite for the build as a whole: the
`unresolved` → `unknown` payload change is an independent edit that mechanically
breaks 5 assertions in the same suite (C3 names them, with line numbers)." `0011:C3`
repeats the fence ("MUST NOT be counted against or excused by A4's
predicate-scoped BREAKS-0 result") and `0011:S6` repeats it a third time
("A green 0005 suite is NOT evidence C1 shipped — S2/S3 are"). The RDR fences
this number three separate times. **Severity: minor — the residual defect is the
undocumented `flowGatedNextModel` / `flag` coupling, which is one clause in S2,
not a scope-creep problem with A4.**

### 7. Pass B C-3 / C-6 / C-8 — the durability findings — **CONFIRMED (all three)**

- **B C-3 (A11 held by prose, not oracle).** Verified: `runFlowNext`
  (`flow_next.go:97`) obtains `owned` once and hands the same slice to
  `assembledView(owned, req.observed)` (`:104`) and to
  `excluded(row, owned, req.observed)` (`:123`), which passes it as
  `resolve.Input.Owned`. A11 is true, and true *by that coincidence*.
  `0011:S5` asserts the three load/parse pins A3 rests on and the empty
  `internal/resolve` diff — it does NOT assert the key-set agreement itself. No
  S1–S8 scenario would fail specifically because A11 broke. B's proposed AT
  (compare `assembledView`'s key set against the kernel's, modulo `recognized`)
  is a cheap, direct pin.
- **B C-6 (filter-after-`KernelRow()` invisible to a refactor).** Verified:
  `internal/table/model.go::Row.KernelRow` (`:232-264`) renders each `BlockMatch`
  atom as `resolve.Tag{Key, Value: seamValue(a.Literal, r.isSet(a.Key))}` — the
  set-canonicalization C1 is protecting. No S1–S8 scenario uses a set-kinded
  match key, so a build that filtered `row.Atoms` before conversion would be
  observationally identical across every mandated fixture. The constraint lives
  only in RDR prose.
- **B C-8 (no multi-key mixed-state fixture).** Verified against the Testing
  Strategy: S3 is "the three reachable match cases, **one row each**"; S4 covers
  the all-absent row and the `recognized`-only row. Nothing exercises one row
  carrying two match keys in different presence/equality states. This is the
  fixture-shaped hole immediately adjacent to Pass A's C-1 — a two-keys-on-one-row
  fixture is one authoring step from a two-atoms-on-one-key fixture, and B found
  the hole without finding what was in it.

---

## Divergence assessment

**Where they agree, the agreement is about LOCATION, not about the defect.**
Both passes converged on `0011:A15` / the demand-set clause, and on `--all`'s
diagnostic dead end. On the first, they found different things: B found the
`Pending`-at-lock process defect (which is real, and which the record's own
Finalization Gate rule would catch); A found the substantive one — the second
caller. On the second they found the same thing from two angles (A: no
invocation reproduces 0005's payload; B: an automated caller cannot open a
TOML). Treating the A15 convergence as validation would have surfaced only the
process complaint and missed the verb-crossing behaviour change entirely.

**Where A goes deeper: the source.** A's three load-bearing findings (C-1/C-12,
C-2, C-6) each rest on a live observation against the shipped tree, and all
three reproduce. A's characteristic move — read the cited symbol for the thing
the RDR did NOT quote (the operator loop under `atomsFromBlock`'s key loop; the
second caller of `invokedReaders`; the alphabet printed beside A6's three) — is
the move that found every defect that matters here. A's own premortem names this
pattern explicitly and it is accurate as self-description.

**Where B goes deeper: durability and the consumer boundary.** B's C-3, C-6, and
C-8 are all "true today, unenforced tomorrow" findings, and all three are
CONFIRMED and are things A did not raise at all. B's C-5 also carries a point A
misses: `intrastate` is by design a generic library with out-of-repo consumers
(the RDR's own Context says so), so the in-repo census is structurally blind to
the population that bears the rename's cost. These are genuine additions.

**Where they contradict.** One real contradiction, on `0011:A4`. Pass A C-4
asserts A4's "BREAKS 0" is being cited as if it cleared the suite; Pass B's
premortem closing paragraph cites the OPPOSITE reading approvingly, quoting S6's
"a green 0005 suite is NOT evidence C1 shipped" as evidence the RDR is honest
about its safety net. **B is right and A is wrong** on this narrow point — the
RDR fences the number three times (A4's own Evidence, C3, S6). A's underlying
observation about the `flowGatedNextModel` / `flag` coupling stands
independently.

A second, softer divergence: A's C-9 reads a gap in A3 that A3's own quoted
predicate (`!= 1`) already covers, while B independently reports having
confirmed A3 as sound. B is closer to right here too.

**Net.** A carries the route-back-forcing findings; B carries the best of the
in-draft additions and is the more accurate of the two on the two points where
they conflict. Neither pass alone is sufficient.

---

## Unhealthy-pass check

**Pass A: healthy.** Every row anchors to a `path::Symbol`, a line number, or a
projector-addressable element, and the load-bearing rows survive independent
re-verification. Two rows are over-stated relative to what the record says
(C-9's A3 gap, which A3's `!= 1` already covers; C-4's "cited beyond its scope",
which the record fences three times) — both are over-reach on a real underlying
observation, not fabrication, and both are cheap to correct. No row re-litigates
Alternatives or Decision Rationale.

**Pass B: mostly healthy, with one clear re-raise and a tendency to generic
framing.**

- **B C-9 is a re-raise, not a defect.** It reopens the O3-vs-O4 fork. The
  permanent `flow next` / `flow resolve` disagreement on an absent-match-key row
  is (i) decided in `0011:§decision-rationale`, on the two rows marked deciding,
  with `0007:C10`'s resolution-level veto and `Resolve`'s control flow as the
  reason O4 loses; (ii) recorded as a cost in `0011:§consequences` in its own
  Negative bullet, in the same terms B uses ("the same report-vs-select
  relationship the verbs already have for `guard_unevaluable`, in the escapable
  direction"); and (iii) explicitly left open as `0007:REQ-78`. B's row adds the
  observation that the analogy is "inexact", which the record itself qualifies
  ("in the escapable direction"). The fix half MUST NOT edit the draft to satisfy
  this row.
- **B C-2 and the premortem's closing paragraph are generic.** "Volume this
  large is exactly the condition under which an implementer's code diverges" is
  unanchored — it names no specific clause an implementer would get wrong that
  is not already covered by B's own C-7. It overlaps A C-11 without adding
  anchoring.
- **B C-5's premortem Ticket 4 partly speculates.** Its concrete mechanism
  ("an off-by-one in `guardOwnedKeys`'s sibling function, or a missed edge case
  where the key also appears in an escape row") invents an implementation bug in
  code that does not exist yet. The underlying finding (out-of-repo consumers are
  structurally invisible to the census) stands on its own and does not need it.

Neither pass is generic overall. Neither pass fabricated a symbol. Two Pass A
rows and one Pass B row should be dropped or downgraded rather than acted on.

---

## Verdict

**Route back to propose/resolve. Two findings force it.**

**Forcing finding 1 — Pass A C-1 / C-12 (verified independently).** `0011:C1`'s
filter is licensed by a stated soundness condition — "a single probe row holds
at most one match tag per key … no cross-atom interaction (A12)" — that is FALSE
on the shipped normalizer, and I reproduced both the two-tag row and the
resulting monotonicity inversion (`k` absent ⇒ candidate; `k` present at any
value ⇒ excluded and unreported) with my own model and harness. This is not an
editing defect. It touches:

- `0011:C1`'s soundness paragraph, whose stated licence must be replaced (either
  by a per-key conjunction-aware filter rule — drop a key's match tags only when
  the key is absent AND the key carries exactly one tag, or keep the row's whole
  match pattern when a key carries more than one tag — or by a stated limit
  naming the shape);
- `0011:A12`, a `Verified` assumption whose Evidence must be re-derived (its
  operator-whitelist and `in`-expansion reasoning are correct; its conclusion is
  not);
- `0011:A3`, whose "the presence test is exact, not approximate" claim is now
  one case short for the same reason;
- `0011:S3` / `0011:S4`, which need the discriminating fixture (and this is
  exactly the fixture Pass B C-8 independently asked for);
- `0011:F1`, since the shape is invisible in the authored TOML, so
  "compare by eye against the model file" cannot diagnose it.

That is a contract, two Verified assumptions, two scenarios, and a failure mode
— a re-proposal of C1's central mechanism, not a draft edit.

**Forcing finding 2 — Pass A C-2 (verified independently).** `invokedReaders`
has two callers (`flow_next.go:97`, `flow_resolve.go:98`), and `runReaders`
refuses on an unbound artifact role (exit 2) or a reader refusal (exit 3) over
the whole invoked set. C1's demand-set MUST is contractually verb-independent
(A15 item (iii)), so it changes which readers `flow resolve` invokes and can
turn a passing `flow resolve` into a refusal. `0011:C1`'s "`flow resolve`'s
dispositions are unchanged in every case" and `0011:JC1`'s "Not a collision" are
both true only for the class each actually examined. The backing assumption
`0011:A15` is `Status: Pending` with `Evidence: "To verify."` while six places
in the record state the extension as settled fact — which the Finalization
Gate's own Status-consistency rule (`0011:§assumption-verification`) forbids
independently. Resolving this needs A15 actually verified and its answer folded
into the Approach, Consequences, and JC1 — and if the answer is "narrow the
extension to `next`", A15 item (iii)'s uniformity argument retracts and C1's
demand-set paragraph is rewritten.

**What does NOT force a route-back, and should be absorbed in-draft on the way
back through:** A C-6 (the 21→3 witness reproduces the alphabet — a
problem-statement honesty edit, not a design change; the narrowing is real, the
metric is over-read); A C-5's prose-staleness half and A C-9 (both one-clause
additions to C3 and S8); A C-4's `flowGatedNextModel` / `flag` coupling (one
clause in S2); B C-3, C-6, and C-8 (three named oracles to add — the A11 key-set
pin, a set-kinded match-key fixture that discriminates filter-before from
filter-after `KernelRow()`, and the multi-key mixed-state row, which doubles as
the C-1 regression fixture); B C-5's out-of-repo-consumer point (one Consequences
bullet). B C-9 must NOT be acted on — it re-raises the decided O3-vs-O4 fork.
A C-7 (text-mode doubling) I did not independently re-verify; it is consistent
with the record's own Load-Bearing Decision and is in-draft either way.
