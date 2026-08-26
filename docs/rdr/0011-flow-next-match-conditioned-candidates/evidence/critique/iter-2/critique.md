Model: claude-opus-5

# Critique — RDR 0011 (iter-2, pass A)

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0011:C1` ("a single probe row holds at most one match tag per key, and the restriction is per-key with no cross-atom interaction (A12)") | FALSE on the shipped normalizer. `atomsFromBlock` loops over OPERATORS within one key, so `[rule.match.k] eq="a"` + `in=["a","b"]` yields a row carrying TWO match tags on key `k` (`{k:a}` and `{k:b}`) — verified live. C1's per-key filter drops both when `k` is absent and keeps both when present, so a present `k` makes the row permanently `no_match` and an absent `k` makes it a candidate. The "no cross-atom interaction" license that makes the filter sound is void. | A rule silently vanishes from `flow next` the moment its match key becomes readable, and appears when the key is unbound — the exact inversion of "supplying state narrows". `--all` shows it; nothing says why. | §1, §4, AT-1 |
| C-2 | `0011:§approach` ("The change lands in the **CLI**… without disturbing `flow resolve`"), `0011:A15` (Status **Pending**), `0011:JC1` | `internal/cli/flow_exec.go::invokedReaders` is called by BOTH `flow next:97` and `flow resolve:98`. Widening its demand set with match-block owned keys widens `flow resolve`'s invoked reader set too. `runReaders` pre-checks every invoked reader's artifact ROLE and returns `flow-artifact-missing` (exit 2), and aborts on `result.Refused()` (exit 3). So a `flow resolve` that succeeds today refuses tomorrow. Consequences' "`flow resolve`'s dispositions are unchanged in every case" (`0011:C1` final para, `0011:§consequences`) is false. | A previously-working `flow resolve --outcome advance` starts failing `flow-artifact-missing: the artifact role \`orphan\`…` or exit 3 on a reader that has nothing to do with the requested outcome. Nothing in the RDR's help, Consequences, or migration note warns them. | §1, §2, §4, AT-2 |
| C-3 | `0011:A15` Status `Pending`; `0011:§prerequisites` ("[ ] A15 verified") | The RDR is being routed to lock with its ONLY behaviour-changing edit outside `flow_next.go` unverified, while `0011:C1`'s normative body, `0011:JC1`, the `authority` mini-check, S8, MVV 8, and Phase 1 all state the extension as settled fact. The Finalization Gate's own "Status consistency" rule (`0011:§assumption-verification`) forbids exactly this. | Whatever A15 turns up lands as a post-lock surprise; if it refutes, C1's demand-set clause is unimplementable and the default degrades to `--all` on any match-only owned key — the failure C1 itself names. | §1, §2 |
| C-4 | `0011:A4` ("BREAKS **0**", "PASSES-UNCHANGED 28") | The classification is scoped to C1's predicate only and explicitly disclaims the payload rename — but it is then cited in `0011:§overrides` and `0011:§risks-and-mitigations` as if it cleared the suite. Independently, `TestReq44_NextInventsNoGuardFactAbsentFromEveryChannel` (`flow_next_0005_test.go:383-410`) iterates EVERY candidate and requires each to list `flag`; it survives only because both rows of `flowGatedNextModel` happen to guard on `flag`. Any C3 fixture row added to that model without a `flag` guard breaks it. C3 mandates adding a discriminating row to `flowMVVModel`, not this one — but the coupling is undocumented and the "BREAKS 0" number invites nobody to check. | A test failure during implementation that the RDR predicted could not happen, spent on re-deriving which oracle A4's census actually covered. | §1, §5, AT-3 |
| C-5 | `0011:C3` ("FIVE TEST reads on this build") | The census counts *reads* and excludes `flow_next_0005_test.go:103`, which it lists as one of the three in that file — but `:103` is `if _, present := c["unresolved"]; !present`, a PRESENCE probe on the map key, not a `stringsAt` read. It re-homes for a different reason (the key is renamed) and needs no ok-bool. Meanwhile `flow_mvv_0005_test.go:46,73` and `flow_next_0005_test.go:23,66,135,137,141,169-171,404-408` carry `unresolved` in COMMENT and MESSAGE text that will become factually wrong. The "count is of test reads only" caveat covers the production sites but not the prose. | Post-implementation the suite is green while a dozen comments and four error messages describe a field that no longer exists. The next reader trusts them. | §1, §5, AT-4 |
| C-6 | `0011:§problem-statement` / `0011:A6` (the `21→3` witness) | The 3 survivors at `stage=resolved` are `prelock`(advance), `resolve-route-back`(revise), `resolve-abandon`(abandon) — verified against `models/rdr.toml:205-236`. That is EXACTLY ONE ROW PER DECLARED OUTCOME, and `outcomes` is `["advance","revise","abandon"]`. The headline result reproduces the full alphabet with a rule id attached. The RDR's own success criterion ("no row is listed that the supplied state excludes") is met while the caller's actual question ("which one") is no better answered than by reading `outcomes[]`. | The skill author who filed kata `1mv1` runs the fixed build, sees three candidates, and is exactly as unable to choose as before. The RDR's motivating defect is not closed by its motivating witness. | §3, §4, AT-5 |
| C-7 | `0011:D-undecided-reporting-shape` (text mode via `respond/text.go::flatten`) | Confirmed by reading `flatten`: a `{key, reason}` object emits TWO path-qualified lines. Today one absent fact is one line (`candidates[0].unresolved[0]: stage`); after C1 it is `candidates[0].unknown[0].key: stage` + `candidates[0].unknown[0].reason: absent`. On 0010's no-tag decision table — which A5 says reports EVERY row with EVERY observed key absent — text output roughly doubles. The RDR notices the flattener and declines to act ("no scenario asserts a text-mode string"). | `flow next --as=text` on a partially-supplied model becomes an unreadable wall. The RDR's Risks section names "looks like the old wall of candidates" and mitigates with the very field that makes the wall taller. | §1, §2, §4, AT-6 |
| C-8 | `0011:C2` (the `--all` `BlockMatch` filter), `0011:F1` | Under `--all`, C2 mandates stripping match facts from `unknown`. `0011:F1` then names `--all` as THE diagnostic for a missing candidate — but a row absent by default and present under `--all` could be excluded by (a) present-and-unequal, or (b) C-1's two-tags-one-key conjunction, and `--all` distinguishes neither. F1 admits the remedy is eyeballing the model file. `intrastate` ships no model-inspection verb (asserted at `flow_surface_0005_test.go:105`). | "Why is `prelock` missing?" has no answer inside the tool. The user greps TOML by hand — the workflow the CLI exists to replace. | §1, §2, §4, AT-7 |
| C-9 | `0011:A3` / `0011:C1` ("the presence test is exact, not approximate") | A3's proof rests on `checkAccessorBindings` refusing `!= 1` readers per owned tag. Verified live: it also fires at **0** readers (`owned tag o is served by 0 readers; want exactly one`). A3 quotes the `%d readers` message but never states the zero arm — which is the arm that makes C1's own "match-only owned key" scenario (S8/A15) load at all. The invariant A3 leans on is stated one case short of what C1's new fixture needs. | S8's fixture is written, fails to load, and the implementer discovers the constraint by trial. | §1, AT-8 |
| C-10 | `0011:§consequences` ("callers scripted against the enumeration must add `--all`"), `0011:C2` | The migration is described as one flag, then contradicted two sentences later ("`--all` restores the candidate SET, not the payload SHAPE"). But it is worse than two migrations: C2 makes `--all` output STRICTLY LESS information than 0005 (match atoms filtered from `unknown`), which `0011:D-naming` concedes is an "asymmetry against the cited precedents". So `--all` is not a compatibility flag at all — no invocation reproduces 0005's payload. | A consumer told "add `--all`" adds it, parses the renamed field, and still loses match-key facts they were reading. Third migration, undocumented as such. | §2, §4, AT-9 |
| C-11 | `0011:§metadata` Profile `large`; `0011:§proportionality` | The record is 807 lines with C1 alone running ~100 normative lines covering: the probe predicate, the demand-set extension, the payload rename, the reason vocabulary, dedup, and sort order. `0011:JC1` concedes the demand-set edit is a shared modify-anchor with a Draft peer. That is a second seam (reader invocation) riding a verb-predicate contract, and it is the one seam whose assumption is `Pending`. | The gate's own split test ("sole author of at most one independent load-bearing contract") should route this back; instead the widest-blast-radius clause is the least-verified one. | §2 |
| C-12 | `0011:A12` ("omitting an absent-key atom relaxes the row rather than inverting a negation") | True per-atom, false per-row once C-1 is admitted: with two tags on one key the row's match pattern is a CONTRADICTION, and omitting both does not relax it — it converts an always-false row into an always-true one. A12's proof ("every match atom is an equality… no negated or absence-testing match atom") examines operators and never asks whether two atoms on one key can coexist. | Same as C-1; recorded separately because A12 is the assumption that licensed C1's filter and is independently wrong. | §1, §4, AT-1 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The probe filter is unsound because a row can carry two match tags on one key (C-1, C-12)

**Root cause in the RDR.** `0011:C1` licenses its filter with an explicit soundness argument:

> "Filtering after conversion is sound because `resolve.Tag` carries `Key`, and the `in`-expansion is row-level (`0002:C13` emits per-member `eq` rows), so **a single probe row holds at most one match tag per key** and the restriction is per-key with no cross-atom interaction (A12)."

`0011:A12` backs it from the operator side:

> "every match atom is an equality — `internal/table/normalize.go::atomsFromBlock` refuses any other operator outright… So omitting an absent-key atom relaxes the row rather than inverting a negation; there is no negated or absence-testing match atom whose omission would change meaning."

**The passage is false, and I verified it live.** `internal/table/normalize.go::atomsFromBlock` iterates keys, and then **iterates the operators declared under that key**:

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

A match block admits `eq` and `in`. Nothing anywhere refuses BOTH on one key. `mergeAtoms` dedupes on `atomIdentity` = `key\0block\0operator\0literal`, so `{k, match, eq, [a]}` and `{k, match, eq, [b]}` are distinct identities and both survive. I loaded this model against the shipped `internal/table`:

```toml
[rule.match.k]
eq = "a"
in = ["a", "b"]
```

and got, on the shipped build:

```
row r1 suffix=[a]
  atom key=k block=match op=eq lit=[a]
  kernel match=[{Key:k Value:a}]
row r1 suffix=[b]
  atom key=k block=match op=eq lit=[a]
  atom key=k block=match op=eq lit=[b]
  kernel match=[{Key:k Value:a} {Key:k Value:b}]
```

The second row carries **two match tags on key `k`**. `TagSet.matches` is an all-must-hold loop, so that row can never match: `k` cannot be both `a` and `b`. It is a dead row — today harmless, because `flow next` strips `Match` entirely and `flow resolve` correctly folds it to `no_match`.

Under C1 it stops being harmless, and it fails **in the wrong direction**:

- `k` PRESENT (any value): both tags are retained, the conjunction is unsatisfiable, probe refuses `no_match`, row **excluded**.
- `k` ABSENT: both tags are dropped, the probe's match pattern is empty, "which matches unconditionally" (C1's own words), row is a **candidate** carrying `{k, absent}`.

So supplying more state makes a row disappear, and supplying less makes it appear. That is precisely the inversion C1 was written to prevent, and it is produced *by* C1's mechanism. C1's remedy for the empty-pattern case ("A row ALL of whose match atoms are omitted is probed with an empty match pattern, which matches unconditionally; it is a candidate") is the clause that converts a contradiction into an acceptance.

The soundness argument's specific error is that it reasons about `in`-expansion (row-level, correct) and about operator kinds (equality-only, correct) and never asks the actual question: *can one row hold two match atoms naming the same key?* It can. `0002:C13`'s per-member `eq` expansion is cited as if it were an exclusivity guarantee; it is a statement about ONE atom's expansion, not about the number of atoms per key.

Yes, such a model is arguably an authoring error — but the loader ACCEPTS it, `flow resolve` handles it correctly today, and C1 is the change that makes the CLI's behaviour on it non-monotonic in supplied state. A contract whose stated soundness condition does not hold on the shipped normalizer is not ready to lock.

**Symptom.** A rule is reported when its key is unbound and silently vanishes once a reader is bound for it. The user binds a reader to *get better answers* and gets fewer, with no `unknown` entry to explain the loss (an excluded row is not reported at all). `--all` shows the row but C2 strips the match facts, so it names nothing.

### 1.2 The demand-set extension changes `flow resolve`, which the RDR says it does not touch (C-2, C-3)

**Root cause in the RDR.** `0011:§approach` states the placement decision in a sentence that is the whole reason this is a re-proposal:

> "The change lands in the **CLI**, and `internal/resolve` is not touched."

and then, in the same paragraph, enumerates the diff as "a probe-shape change in `internal/cli/flow_next.go::excluded`, **one added term in `internal/cli/flow_exec.go::invokedReaders`' demand set**…". `0011:C1`'s normative body makes the extension MUST-level. `0011:§consequences` promises "the kernel's match seam stays two-valued and **`flow resolve`'s dispositions are unchanged in every case**" (`0011:C1` closing paragraph).

**"internal/resolve is not touched" and "flow resolve is unchanged" are two different claims, and the RDR proves the first and asserts the second.** `invokedReaders` has exactly two callers:

```
internal/cli/flow_next.go:97:    readers, owned, ce := req.runReaders(cmd.Context(), invokedReaders(req.model, ""))
internal/cli/flow_resolve.go:98:        invokedReaders(req.model, outcome))
```

`flow resolve` is a caller. Widening the demand set widens which readers `flow resolve` invokes. And `runReaders` is not a passive widening — it has two refusal arms that fire on the invoked set:

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

So a `flow resolve --outcome X` that succeeds today refuses tomorrow with `flow-artifact-missing` (exit 2) if the newly-demanded reader's role is unbound, or exit 3 if that reader refuses. `flow resolve`'s payload also reports `readers` (`flow_resolve.go:145`), so the widening is directly observable in its output.

`0011:A15` — the assumption that is supposed to cover this — actually concedes it, in item (ii):

> "a model whose match-only owned key is served by a refusing reader now fails where it previously succeeded, which is the CORRECT behaviour… but is a behaviour CHANGE"

and item (iii) asserts the widening "is uniform rather than `next`-special-cased", which is a statement that `flow resolve` changes too. The Approach, Consequences, and `0011:JC1` do not carry that forward. `0011:JC1`'s joint-check reasons entirely about 0010's zero-owned-tag class and concludes "Not a collision" — true for that class, and silent on every `state-machine`-class model where `flow resolve`'s reader set moves.

And A15 is **`Status: Pending`**, `Evidence: "To verify."` (C-3). `0011:§prerequisites` marks it unchecked. Yet C1's normative body states the extension as MUST, the `authority` mini-check tabulates it as settled ("three terms after C1"), `0011:JC1` reasons from it, S8 and MVV 8 test it, and Phase 1 implements it. The Finalization Gate's own instruction (`0011:§assumption-verification`) says: "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it." Six places depend on it.

**Symptom.** A user's `flow resolve` pipeline — untouched, same model, same argv — starts refusing after the upgrade, naming an artifact role or reader that has nothing to do with the outcome they asked for. The release note says "flow next now selects by match".

### 1.3 The `--all` escape hatch does not restore anything, and the diagnosis path is a dead end (C-8, C-10, C-7)

**Root cause in the RDR.** Three passages, each individually defensible, compose into a trap.

`0011:C2` requires that under `--all` the match facts be filtered out of `unknown`:

> "Under --all that entry MUST disappear — reporting a fact the mode ignores would be reporting on a predicate it does not apply — so the --all branch MUST filter BlockMatch atoms out of `unknown`, and an oracle MUST assert their absence."

`0011:F1` then makes `--all` the diagnostic:

> "**Visible**: a candidate the caller expected is absent. Diagnose with `--all`: if it appears there, the row was match-excluded… `--all` localizes the cause to match, **and no further**… the row's authored pattern is read from the model file, compared by eye against `owned`/`observed`."

And `0011:§consequences` sells the migration as one flag before withdrawing it:

> "callers scripted against the enumeration must add `--all`. `--all` restores the candidate SET, not the payload SHAPE"

`0011:D-naming` concedes the deeper problem outright: "here the widened rows carry strictly LESS" information than the default's — an asymmetry against every CLI precedent it cites (`docker ps -a`, `git branch -a`).

**Composed, these mean: there is no invocation of the new build that reproduces 0005's output.** Not `--all` (payload renamed, match facts stripped), not the default (predicate changed). A consumer is told "add `--all`" and discovers they owe (1) the flag, (2) the `unresolved`→`unknown` field rename, (3) the element type change `string`→`{key,reason}`, and (4) permanent loss of match-key facts they may have been reading. Consequences names (1) and (2)-(3); (4) is buried in a Load-Bearing Decision.

Meanwhile the diagnosis loop is closed. A row is missing. `--all` tells you "match excluded it" and, per C2, deliberately refuses to say which key or why. `intrastate` ships no model-inspection verb — `flow_surface_0005_test.go:105` asserts `dump` is FOREIGN to the flow group, and `0011:BR3` rejects `excluded_by`. So F1's honest answer is: open the TOML and compare by eye. That is the workflow the tool exists to replace, and under C-1 the eye-comparison will not even find the cause, because the two-tags-on-one-key contradiction is invisible in the authored source (it reads as two ordinary predicates).

Layer C-7 on top. `0011:D-undecided-reporting-shape` correctly identifies that text mode goes through `respond/text.go::writeTextPayload` → `flatten`, which emits one path-qualified line per LEAF. I read `flatten`: a `{key, reason}` object is two leaves, so every undecided fact becomes two lines instead of one. On 0010's no-tag decision table — where `0011:A5` establishes that EVERY row is reported with EVERY observed match key `absent` — text output roughly doubles. The RDR sees this and declines: "no scenario asserts a text-mode string." `0011:§risks-and-mitigations` opens with "**Risk**: an absent match key makes the list look like the old wall of candidates" and mitigates with the field that makes the wall twice as tall in the default output mode.

**Symptom.** The user upgrades to fix "too many candidates", gets the same candidate count on their decision table with twice the text output, and when a candidate they wanted goes missing, the tool tells them to read the model file by hand.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0011:C1`'s demand-set paragraph** — lines beginning "The assembled view MUST actually carry the keys the predicate reads, so `internal/cli/flow_exec.go::invokedReaders` MUST add each row's MATCH-block owned keys to its demand set…" (C-2, C-3, C-11).

It is the only MUST in the record whose supporting assumption is `Pending` (`0011:A15`), and it is the only clause that reaches a function `flow resolve` shares. Its blast radius is not `flow next`'s payload — it is which subprocesses run, which artifact roles must be bound, and therefore which invocations of a **second, unrelated verb** refuse.

The rewrite arrives by one of three routes, all short:

1. **A15 verifies but the fixture bites.** `0011:C3` mandates a match-key fixture that "MUST match on a key it does NOT write or clear, writing some OTHER key instead". Combined with C1's demand-set clause, every such key now demands its reader. The MVV corpus is built around `[read.orphan]` (`flow_fixtures_0005_test.go:122-126`) — a reader deliberately left with an UNBOUND role, whose non-invocation is the ONLY oracle separating the narrowed reader set from "every declared reader" (the fixture header says so verbatim: "the only assertion that separates the narrowed invoked read-accessor set (REQ-35)"). The first fixture that matches on a key that reader serves turns the corpus's central discriminator into a `flow-artifact-missing`. Someone adds a carve-out.

2. **A `flow resolve` regression lands in the field.** The demand-set term applies to both verbs; the RDR's Consequences promises `flow resolve` is unchanged; a user proves otherwise. The clause is narrowed to `next` — at which point `0011:A15`'s item (iii) ("the widening is uniform rather than `next`-special-cased") is retracted, and the RDR's own uniformity argument no longer holds.

3. **A15 refutes.** `0011:A15`'s "If wrong" already writes the replacement: "C1 owes a stated limit instead of the extension — the match-only owned key stays absent, the row reports `{key, absent}` with no caller remedy… and the help C3 mandates must say so." That is C1's demand-set paragraph deleted and C3's help text rewritten.

Adjacent and near-certain: **`0011:C3`'s five-read census** (C-5). It is precise about the wrong denominator. `:103` is counted as a read but is `if _, present := c["unresolved"]; !present` — a map-key presence probe, which re-homes trivially and needs no ok-bool, unlike the three the contract flags. And the census explicitly covers "test reads only", leaving `unresolved` in comments and `t.Errorf` message text at `flow_next_0005_test.go:23,66,104,135,137,141,169-171,404-408` and `flow_mvv_0005_test.go:46,73` — prose that becomes false the day the field is renamed. A census that names line numbers gets trusted as exhaustive.

---

## 3. The one assumption that will not survive first contact with a real user

**That narrowing the candidate set to "what the state can take" answers the caller's question.**

`0011:§problem-statement` states it directly:

> "A skill author (or the human driving one) runs `flow next` to ask 'what is legal from here?' so the next action can be chosen without re-deriving the transition table by hand."

and defines success as narrowing:

> "Success is therefore 'no row is listed that the supplied state excludes, and every row listed that the supplied state could not decide names the undecided key', not a target cardinality."

The headline witness is `0011:A6`: 21 → 3 on `models/rdr.toml` at `stage=resolved`.

**I checked the three.** From `models/rdr.toml:205-236`:

| rule | outcome |
|---|---|
| `prelock` | `advance` |
| `resolve-route-back` | `revise` |
| `resolve-abandon` | `abandon` |

The model's alphabet (`models/rdr.toml:12`) is `outcomes = ["advance", "revise", "abandon"]`.

**The three surviving candidates are exactly one per declared outcome — the full alphabet, with a rule id attached.** The payload already carried `outcomes[]` unconditionally, in both modes, in 0005, and `0011:C1` requires it keep doing so ("`outcomes` MUST remain the model's full declared alphabet in every mode").

So on the very model that motivated the RDR, at the state that motivated it, the new candidate list conveys no information the old payload's `outcomes[]` field did not already carry. The user's question — filed as kata `1mv1`, "cannot constrain a skill's next step" — is answered with the same three-way choice they had before.

This is not a fluke of `stage=resolved`. The model is a lifecycle: at every stage, each declared outcome has exactly one responding rule. Scan the `match.stage` census — `seeded` has 3 rules (`propose`/`seed-abandon`/`seed-revise`), `proposed` has 3, `refined` has 3, `prelocked` has 3, `final` has 3. `reconciled` has 4 (`finalize-pass`, `finalize-blocked`, `reconcile-route-back`, `reconcile-abandon`) — two of which bind `advance` and are separated by a GUARD (`gate_passed`), which the RDR leaves untouched. **The 21→3 result generalizes to "n rules → one per outcome" across this entire model, which is the alphabet.**

`0011:§problem-statement` half-anticipates this — "The outcome is a candidate list narrowed to what the supplied state can take — NOT a single next action" — and pre-defends it: "a caller facing several candidates is being told the state genuinely admits several." That defence is correct in principle. It is also the answer the user already had, and it means the RDR's headline metric measures the removal of rows the caller was never going to pick (`stage=seeded` rules while at `stage=resolved`) rather than any gain in decidability.

The user's real need — visible in the problem statement's own framing, "so the next action can be chosen" — is discrimination *among the outcomes*, which is exactly what `flow next` is contractually forbidden from providing (`0011:D-selection-predicate`, unchanged from `0005`). The first real user runs the fixed build, sees three candidates where they wanted one, and files the same kata again.

---

## 4. The premortem

*Written from twelve weeks after ship.*

We shipped 0011 in the first week. `make check` was green, the S1 oracle asserted `candidates[]` was exactly `prelock`, `resolve-route-back`, `resolve-abandon` on `models/rdr.toml` at `stage=resolved`, and the `21→3` line went in the release note. Nobody noticed that those three rules are one per declared outcome and that `outcomes[]` in the same payload had been saying `["advance","revise","abandon"]` since 0005. We had built a filter whose best-case output on our own flagship model is the field that was already there.

**Week 2 — the skill author comes back.** The same person who filed kata `1mv1` runs `flow next --model models/rdr.toml --artifact rdr=./0011.md`. Three candidates. They still cannot pick. They re-read the help text that `0011:C3` had us write — "reports the candidates the supplied state can take" — and it is accurate, which is what makes it useless. They file the follow-up: "next still doesn't tell me what to do." We point at `0011:§problem-statement`'s "NOT a single next action" and close it as working-as-designed. It reopens in month two under a different title.

**Week 3 — `flow resolve` starts refusing.** A consumer model in the `state-machine` class matched on an owned key that no rule wrote, cleared, or guarded — the exact case `0011:C1`'s demand-set paragraph was added to fix, and the case `0011:A15` was still `Pending` about when we locked. `invokedReaders` now demanded that key, `flow_exec.go::runReaders` pre-checked the serving reader's artifact role, found it unbound, and returned `flow-artifact-missing` at exit 2. Not from `flow next` — from `flow resolve --outcome advance`, in a CI pipeline that had not changed. We had written in `0011:§approach` that "the change lands in the CLI, and `internal/resolve` is not touched", and in `0011:C1` that "`flow resolve`'s dispositions are unchanged in every case". Both were true about `internal/resolve` and false about `flow resolve`, and nobody had run `grep -n 'invokedReaders(' internal/` — which returns `flow_next.go:97` and `flow_resolve.go:98`. `0011:A15` item (ii) had said in so many words that a newly-invoked refusing reader "now fails where it previously succeeded"; it sat in a `Pending` assumption while `0011:JC1` reasoned about 0010's zero-owned-tag class and concluded "Not a collision."

The fix took three days because the demand set is mode-independent by contract (`0011:C1`: "The demand set is a property of the MODEL, not of the mode") and verb-independent by A15 item (iii) ("the widening is uniform rather than `next`-special-cased"). We special-cased `next`, contradicting both, and the RDR's uniformity argument stopped applying to the code it described.

**Week 5 — the rule that vanishes when you bind a reader.** A consumer authored, on one key, both an `eq` and an `in` under `[rule.match]`. Two operators, one key. `internal/table/normalize.go::atomsFromBlock` loops over operators within each key, `mergeAtoms` dedupes on `key\0block\0operator\0literal` so both survive, and `expand` emits a row carrying `{k:a}` AND `{k:b}` in `KernelRow().Match`. Unsatisfiable. Harmless for two years, because `flow next` stripped `Match` and `flow resolve` folded it to `no_match` correctly.

Then C1's filter arrived. With `k` absent, `excluded` dropped both tags and the probe matched unconditionally — C1 says so explicitly: "A row ALL of whose match atoms are omitted is probed with an empty match pattern, which matches unconditionally; it is a candidate." With `k` present, both tags rode the probe, the conjunction failed, `no_match`, row excluded and unreported.

So binding a reader made a rule disappear. The user's bug report was "flow next shows fewer rules the more state I give it", which is the inverse of the sentence at the top of `0011:§approach`. `--all` showed the row, and `0011:C2`'s mandatory `BlockMatch` filter meant it named nothing about why. `0011:F1` sent them to the model file to compare by eye, where the two atoms read as two ordinary predicates and the contradiction is invisible. We had asserted the impossibility twice: `0011:C1` ("a single probe row holds at most one match tag per key… the restriction is per-key with no cross-atom interaction (A12)") and `0011:A12` ("every match atom is an equality… no negated or absence-testing match atom whose omission would change meaning"). Both reasoned from `0002:C13`'s per-member `in`-expansion, which is a claim about ONE atom, and neither asked whether one key can carry two.

**Week 6 — text mode.** 0010 shipped its decision-table class. `0011:A5` had already established the shape: no `--tag`, every ordinary rule reported, every observed match key `absent`. Under 0005 that was one text line per fact. Under `0011:D-undecided-reporting-shape`, `unknown` is a list of `{key, reason}` objects, and `respond/text.go::flatten` emits one path-qualified line per leaf — so `candidates[0].unknown[0].key: status` and `candidates[0].unknown[0].reason: absent`. Two lines where there was one, across every row of the table. The RDR had read `flatten` and decided not to act: "no scenario asserts a text-mode string." `0011:§risks-and-mitigations` had opened with "**Risk**: an absent match key makes the list look like the old wall of candidates" — and the mitigation doubled the wall in the default output mode.

**Week 9 — the migration nobody could complete.** An out-of-repo consumer read our note: "callers scripted against the enumeration must add `--all`." They added it. Their parser broke on `unresolved` → `unknown`. They fixed that. It broke on `string` → `{key, reason}`. They fixed that. Then they found the match-key facts they had been reading were gone permanently, because `0011:C2` requires `--all` to filter `BlockMatch` atoms out of `unknown` and `0011:D-naming` had already conceded that "here the widened rows carry strictly LESS" than the default's. There is no invocation of the shipped binary that reproduces 0005's payload. `0011:§consequences` had named two of the four migrations and called `--all` a restoration.

**What the postmortem said.** The record was 807 lines, ground-swept clean at 58 anchors, with a hardened premortem at 15 findings and PASS. Every load-bearing failure came from a passage that *cited real symbols and drew the wrong conclusion from them*: `atomsFromBlock` was read for its operator whitelist and not for its operator loop; `invokedReaders` was read for its demand terms and not for its second caller; `flatten` was read and then discounted; `A6`'s three surviving rules were counted and never compared against the three-member alphabet printed beside them. Anchor density was never the problem. Nobody wrote down what each anchor would look like if the claim were false.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

### AT-1 — one key, two match atoms (C-1, C-12)

```gherkin
Given a model whose rule authors, under one [rule.match.<k>] block,
      both `eq = "a"` and `in = ["a", "b"]`
  And the model LOADS (no refusal from internal/table)
 When I enumerate the normalized rows' KernelRow().Match
 Then some row carries TWO resolve.Tag entries whose Key is "<k>"
```
*Discriminates:* passes today on the shipped normalizer — I ran it. C1's "at most one match tag per key" and A12's "no cross-atom interaction" are refuted before any code is written.

```gherkin
Given that model, and the C1 probe builder
 When key "<k>" is ABSENT from the assembled view
 Then the row is reported as a candidate with {<k>, absent}
 When a reader is bound establishing "<k>" to ANY value
 Then the row is NOT reported, and no `unknown` entry explains it
```
*Discriminates:* the monotonicity violation. Supplying state must never remove a row that absence kept. Negative control: a row with ONE match atom on `<k>` behaves monotonically.

### AT-2 — the demand-set extension is not `flow next`-local (C-2)

```
Step 1. Run: grep -n 'invokedReaders(' internal/
Step 2. Assert the result set is exactly one call site.
```
*Discriminates:* fails today — two sites, `flow_next.go:97` and `flow_resolve.go:98`. One command, run at review, refutes "the change lands in the CLI [and disturbs no other verb]".

```gherkin
Given a model with an owned key some row MATCHES but no row writes, clears, or guards
  And that key's serving reader binds an artifact role left UNBOUND
 When I run `flow resolve --outcome <o>` with the pre-0011 demand set
 Then it exits 0 with a plan
 When I run the same argv with the C1 demand set
 Then it refuses `flow-artifact-missing` at exit 2
```
*Discriminates:* directly contradicts `0011:C1`'s "`flow resolve`'s dispositions are unchanged in every case". Negative control: 0010's zero-owned-tag `decision-table` class — identical disposition both ways, which is the only case `0011:JC1` actually checked.

```gherkin
Given the flowMVVModel corpus, where [read.orphan] serves `note`
      and its role is deliberately left unbound
 When any rule is amended to MATCH on `note`
 Then `flow next` refuses `flow-artifact-missing`
  And the REQ-35 reader-narrowing oracle no longer discriminates
```
*Discriminates:* the C3 fixture mandate colliding with the MVV corpus's central invariant.

### AT-3 — A4's census is predicate-scoped, not suite-scoped (C-4)

```gherkin
Given TestReq44_NextInventsNoGuardFactAbsentFromEveryChannel
 Then it iterates EVERY candidate and requires each to list "flag"
  And it passes only because both flowGatedNextModel rows guard on `flag`
 When any row is added to that model without a `flag` guard
 Then it fails
```
*Discriminates:* "BREAKS **0**" is conditional on a fixture property A4 never states. Would have forced A4 to publish its per-oracle classification, not just the totals.

### AT-4 — the `unresolved` census counts reads, and the record is not only reads (C-5)

```
Step 1. grep -rn 'unresolved' internal/cli/ | wc -l    → 30
Step 2. Of those, classify: production field/tag/comment;
        test READ via stringsAt; test PRESENCE probe; comment; message text.
Step 3. Assert every occurrence has a stated disposition in 0011:C3.
```
*Discriminates:* `flow_next_0005_test.go:103` is a PRESENCE probe (`if _, present := c["unresolved"]`), not a `stringsAt` read, yet C3 counts it among the three needing ok-bool treatment. And `:23,66,104,135,137,141,169-171,404-408` plus `flow_mvv_0005_test.go:46,73` carry the name in prose that goes stale silently.

### AT-5 — the narrowing witness reproduces the alphabet (C-6)

```gherkin
Given models/rdr.toml seeded at stage=resolved
 When `flow next --as=json` runs under the C1 predicate
 Then candidates[] is prelock, resolve-route-back, resolve-abandon
  And the set of their `outcome` fields equals `outcomes[]` exactly
 Then assert: |distinct candidate outcomes| < |outcomes[]|
```
*Discriminates:* the final assertion FAILS on the RDR's headline witness — 3 distinct outcomes, 3 declared. This is the test that turns "21→3" from a success metric into the question "and how is that better than `outcomes[]`?" Run it at review and MVV 3 stops being evidence for the problem statement.

Generalization check, also at review:
```
For each declared value v of `stage`:
  seed stage=v; assert |distinct candidate outcomes| < |outcomes[]|
```
Fails at `seeded`, `proposed`, `refined`, `resolved`, `prelocked`, `final`.

### AT-6 — text mode doubles (C-7)

```gherkin
Given a 0010-class decision table run with no --tag
 When `flow next --as=text` runs under 0005
 Then each undecided fact is ONE line
 When it runs under 0011's {key, reason} shape via respond/text.go::flatten
 Then each undecided fact is TWO lines (…unknown[i].key, …unknown[i].reason)
  And total output line count is ≥ 1.8× the 0005 baseline
```
*Discriminates:* makes the flattener consequence a measured number instead of a Load-Bearing Decision's parenthetical. `0011:§risks-and-mitigations`' first risk is "looks like the old wall of candidates"; this test shows the mitigation makes it worse in the default mode.

### AT-7 — `--all` cannot diagnose what it is prescribed for (C-8)

```gherkin
Given a row excluded by default
 When I run with --all and the row appears
 Then the payload names WHICH match key excluded it
```
*Discriminates:* fails by construction — `0011:C2` MANDATES stripping those facts. Stating it as a failing acceptance test forces `0011:F1`'s "read the model file by eye" to be evaluated as the shipped remedy it is, against `0011:BR3`'s rejection of `excluded_by`.

### AT-8 — the presence-test invariant is stated one arm short (C-9)

```gherkin
Given a model declaring an owned tag served by ZERO readers
 When it loads
 Then it is refused: "owned tag <k> is served by 0 readers; want exactly one"
```
*Discriminates:* passes (I ran it), and `0011:A3` quotes only the `%d readers` format string without stating the zero arm. That arm is what governs whether S8/A15's match-only-owned-key fixture can be authored at all.

### AT-9 — no invocation reproduces 0005 (C-10)

```gherkin
Given a caller scripted against 0005's flow next payload
 When they run the 0011 build with --all
 Then their parser must handle: the field rename `unresolved`→`unknown`,
      the element type change string→{key,reason},
      and the permanent loss of BlockMatch facts (0011:C2)
 Then assert: some invocation of the 0011 binary is byte-compatible with 0005
```
*Discriminates:* the final assertion has no witness. Forces `0011:§consequences` to state four migrations rather than describing `--all` as a restoration, and forces `0011:D-naming`'s conceded asymmetry ("the widened rows carry strictly LESS") into the migration note where a consumer will read it.
