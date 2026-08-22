Model: claude-opus-5

# Hostile Set Critique — RDR 0003 · 0006 · 0007

Stage 7.1 whole-set critique. Target: the three-RDR cluster taken **together**.
Single-model pass (`claude-opus-5`); recorded as a fallback per
`2-critique.md` — no alt model reachable this session. The `large` profile
permits this; it is recorded, not passed off as dual-model.

Assignment is hostile critique. No hedging, no balance.

**Verdict up front.** The three-way division of labour is not implementable as
written, and the reason is structural, not editorial: **the cluster's one open
question (JDR 0001 §JD-4) is the only thing holding the three documents apart,
and every document has routed its closure to a venue rather than answering it.**
The set has a bidirectional Draft-on-Draft dependency (0003 A8/A10/A12 need
0006; 0006's Refinement Context items 1–4 need 0003) whose only stated exit is
"cluster reconcile" — this document. If this pass defers again, the deadlock is
permanent by construction.

Against that: 5,037 lines of RDR body and 19,125 lines of evidence — **24,162
lines** for this three-RDR set — versus **1,236 lines** of non-test Go across
all of `internal/` and `cmd/`, of which the product surface is a `version`
subcommand. Ratio 19.5:1. The last commit touching `internal/` or `cmd/` is
`1c9c0ca` (2026-08-09); **129 commits** have landed since, all specification.

---

## Findings Ledger

| ID | RDR | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-----|-------------|--------------|-------------------|--------|
| S-1 | set | 0003 A8 "Stage 6 disposition — DOWNGRADED to the cluster-reconcile venue"; 0006 Refinement Context Direction 1 | §JD-4's recording assignment is routed to a venue by both parties; neither document may state it unilaterally, so the deadlock survives any number of passes | Implementer finds two Drafts, each citing the other as the blocker; no one can start | §1, §2, premortem, AT-1 |
| S-2 | set | 0003 A10 (`Pending`) + 0006 Refinement Context defect 2; `row group` occurs 32× in 0003, **1×** in 0006 | Circular definition: 0003 defines the row group and books A10 to confirm 0006 accepts; 0006's A2 takes the coverage derivation *from* 0003. Neither is the floor | Lint groups rows one way, the exhaustiveness proof another; green lint over an unproved product | §1, premortem, AT-2 |
| S-3 | set | 0003 A12 (`Pending`); `reachable predecessor` occurs 5× in 0003, 3× in 0006, **0×** in 0007 | Two 0003 clauses quantify over "every reachable predecessor"; the document that owns graph traversal states no reachability contract a peer can cite | `owned-set-before-match` is unimplementable as specified; lint either over- or under-blocks | §1, AT-3 |
| S-4 | 0007 | Predecessors block line 159: "**RDR 0006** (`Final`, tolerance §JD-4) narrows lint's promise" | A `Final` RDR asserts a peer's status that is false since 2026-08-21; 0006 records the staleness but declares it "not this RDR's to fix" | Reader trusts 0007 (Final) and builds against a tolerance no document carries | §1, §2, AT-4 |
| S-5 | set | JDR §JD-4 "Open only as to which document records the narrowing" | The `guard_unevaluable` narrowing is routed through three documents plus a JDR entry to state one sentence. The cut is wrong, not the wording | Four places to read; three can drift; the one-sentence rule is nowhere normative | §1, §3, premortem |
| S-6 | 0006 | `guard_unevaluable`: **0** occurrences in any normative block (only line 10 Status, line 825 Refinement Context); `RDR 0007`: 0 outside the Refinement Context (lines 875/878) | The document that emits the finding has no normative knowledge of the runtime veto that narrows its promise | Lint certifies a group exhaustive; runtime refuses `guard_unevaluable` on that group. The headline guarantee (P5) is false | §1, premortem, AT-5 |
| S-7 | set | JDR §D1 "Resolved: (d) `Row` carries parsed atoms"; `internal/resolve/resolve.go:185` `Guard string` | A joint decision recorded as Resolved has zero code effect after 129 commits; the shipped type still contradicts it | Every atom-based clause in all three RDRs is written against a type that does not exist | §1, §2, premortem, AT-6 |
| S-8 | 0003 | Status line: 15 assumptions, **5 Pending**; 24 normative blocks; 1,755 lines | A `Draft` is the cluster's normative center of gravity: it now owns atom grammar, exhaustiveness, **and** the tag declaration model (rehomed 2026-08-21) | Two `Draft`/`Final` peers cite a Draft's normative clauses as their floor | §1, §2, §3 |
| S-9 | 0007 | Status `Final`; A12 "Status: **DOWNGRADED** — RDR 0003 is silent on the projection" | A Final RDR carries a downgraded assumption whose home is an open JDR entry; `Final` overstates what was verified | Gate PASS recorded on an open question; re-lock is already owed and untracked | §2, AT-7 |
| S-10 | 0006 | A2 "**Re-verify at refine**"; A5 `Pending` (MVV, CI wiring) | The only member with a user-visible surface (`intrastate lint`) has its authority assumption unverified and no `lint` verb in the tree | `intrastate lint` does not exist; grep for `"lint"` in `internal`/`cmd` returns nothing | §2, §3, premortem, AT-8 |
| S-11 | set | 0003 mentions `RDR 0006` 71× and `RDR 0007` 46×; 0006 mentions `RDR 0007` 0× (body) | Citation traffic is one-directional into 0003. The set is a hub-and-spoke around a Draft, not three peers | Changing 0003 invalidates both peers; 0003 cannot lock until both agree | §1, §3 |
| S-12 | set | 24,162 lines RDR+evidence vs 1,236 lines non-test Go for this set | Specification is the deliverable; the cluster optimizes document consistency over a running program | User has a `version` command after 129 commits | §1, §2, premortem |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The deadlock is never broken, because every party routed it to a venue

**Root cause.** §JD-4's substance is decided ("the *promise* narrows — P5
decides"). What is open is only *which document records it*. Both candidate
documents have written down that they cannot decide it alone:

- 0003 A8, `Stage 6 disposition`: "**DOWNGRADED to the cluster-reconcile
  venue.** Not closable inside this RDR alone: JDR 0001 §JD-4 is an open ledger
  entry, so the assignment is a joint decision. … Named plan: the
  cluster-reconcile pass over 0003 · 0006 · 0007 carries §JD-4 and returns the
  recording assignment plus the finding-code question."
- 0006 Refinement Context, Direction 1: "**Record or cite the §JD-4 narrowing,
  and say which.** Recommended: RDR 0003 records it … Whichever arm is taken,
  exactly one document states it."
- JDR §JD-4 itself: "**Recommended arm (2026-08-21, from 0003's Stage 6):** 0003
  records it and 0006 cites it."

Three documents recommend the *same arm* and none of them takes it. A8 is
`Pending` with `Method: Peer RDR`. This is the defining pathology of the set:
**a decision that everyone agrees on, that no one will write down**, because
each document has correctly concluded that unilateral action would duplicate
prose. The routing is individually rational and collectively paralysing.

**Specific passage that enabled it.** 0003 A8's `Plan`: "**Close on the
assignment, not on matching prose.**" That instruction is right, and it means
A8 cannot be closed by *any* edit to 0003 or 0006 — only by an edit to
**JDR 0001 §JD-4**, which neither RDR's flow stage touches.

**Symptom.** An implementer opens the cluster, finds 0003 `Draft` blocked on
A8/A10/A12, opens 0006 to unblock them, finds 0006 `Draft` blocked on 0003's
model, and reports the cluster unimplementable. Which it is.

### 1.2 Row-group membership is defined in a circle, so lint is unimplementable

`row group` occurs **32 times** in 0003, **once** in 0006, **zero times** in
0007. 0003's Technical Design says: "**The row group is defined here**, not
deferred … RDR 0006 supplies the graph traversal that enumerates which selection
contexts are reachable; it does not define the grouping predicate, and this RDR
does not read one back from it. (The term `row group` appears nowhere in RDR
0006, and its A2 delegates the coverage derivation here — so a definition
deferred to RDR 0006 would close a cycle with no floor. A10 books the
confirmation that RDR 0006 accepts this division.)"

0003 correctly diagnoses the cycle and then **books its resolution as a
`Pending` assumption on the peer** (A10). Meanwhile 0006's A2 says: "**Re-verify
at refine**: confirm this RDR's scoped grouping matches RDR 0003's row-group
definition (its A10)". A10 waits on 0006; 0006's A2 waits on A10. This is the
same deadlock as 1.1, on a second axis — and note the prior cluster critique
already logged it as C-20 ("Group membership never defined beyond 'the same
state/outcome'"). It has been open across two reconcile passes.

**Symptom.** Lint's overlap and coverage checks compare different row sets than
the exhaustiveness proof reasons over. A green `intrastate lint` does not imply
resolution succeeds — which is exactly what Principle P5 of the JDR promises it
must imply.

### 1.3 The atom shape the whole set is written against does not exist in code

JDR §D1 records **Resolved: (d)** — `Row` carries parsed atoms, not a string —
and 0007 absorbs the change. All three RDRs are written against that shape:
`atom` occurs 131× in 0003 and 223× in 0007. 0003's Technical Design cites it as
fixed: "The atom shape is fixed at the kernel seam by JDR 0001 §D1 and stated
normatively in RDR 0007 (`docs/rdr/0007-guard-predicate-totality.md:1249-1266`)
— four fields: `Key`, `Operator`, `Literal`, `Block`."

`internal/resolve/resolve.go:185` still reads:

```go
	// Guard is the predicate handed to the guard seam. Empty means no
	// guard.
	Guard string
```

The decision was recorded 2026-08-21; the last commit touching `internal/` is
2026-08-09 (`1c9c0ca`). **129 commits have landed since, none of them code.**
The cluster has been reasoning for two weeks about the semantics of a field
shape that no one has written.

**Symptom.** Phase 1 of either implementation plan begins with a `Row`
reopening that touches the single non-test consumer (`resolve.go`) and all
2,735 lines of test. Every downstream estimate in the set assumes that change
is already made.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**RDR 0003's `#### Normative Contracts` (lines 663–892, 24 normative blocks) —
specifically the tag declaration model rehomed into it on 2026-08-21.**

It will be rewritten because it was moved for a *documentary* reason, not an
implementation one. JDR *Ownership corrections* justifies the rehome as: "The
model was homeless, not homed elsewhere." That is a true finding about the
document graph. But the model's actual consumers are 0002's normalizer (carriage
and authoring location) and 0006's lint input contract (finite-domain metadata)
— neither of which is 0003. 0003 owns it because 0003 is where the *proof* needs
it, which is the weakest of the three claims to ownership.

The mechanical prediction: the first implementer writing the TOML loader needs
value kind and finite domain at parse time, in 0002's normalizer, before a row
exists. They will put the struct there. 0003's normative blocks will then
describe a model that lives in another package, and the rehome will be reversed
or split — for the second time, since it already moved 0002→0003 once.

Corroborating signal: this section is the newest normative surface in the set
and it is already carrying the cluster's largest `Pending` load. 0003 has **5 of
its 15 assumptions Pending**, against 1 of 5 in 0006 and 0 of 26 in 0007.

**What would have to be true for no rewrite:** the tag declaration model would
need a consumer inside 0003's own implementation phase that is not a lint or a
normalizer concern. There isn't one — 0003's evaluator, per §D4, "never sees the
tag view" and only does value comparisons over a present value.

---

## 3. The cross-cutting assumption that will not survive first contact

**"Exactly one document records each contract, and citation is free."** (JDR
Principle P6: "One home per contract; cite, never restate.")

The set treats a citation as costless and a restatement as the only hazard. The
`guard_unevaluable` narrowing shows the assumption failing in both directions at
once. To state one rule — *lint may not certify a row group exhaustive when a
participating row could refuse `guard_unevaluable` at runtime* — the set
currently requires:

1. **JDR §JD-4** to assign the recording document (open);
2. **0003** to state it for the finite-domain case (done, line 652 ff.);
3. **0006** to cite it and reuse `graph-unprovable-coverage` (not done — the
   token `guard_unevaluable` occurs **2×** in 0006, at line 10 (Status line) and
   line 825 (Refinement Context), and **zero times in any normative block**;
   `graph-unprovable-coverage` is scoped to the non-finite case);
4. **0007** to hold the runtime veto that triggers it (done, line 1656).

Four coordination points for one sentence. P6 assumes the cost of a citation is
a pointer; here the cost is a `Pending` assumption, a demotion, and a venue.

**Where it breaks with a real user:** the first person to add a tag with an
unbounded integer domain. They run `intrastate lint`, get
`graph-unprovable-coverage` for the non-finite dimension, bound the integer,
re-run, and get **green**. Then resolution refuses `guard_unevaluable` at
runtime because a participating row has an `exists` atom over an optional key —
0003's case, which 0006's clause "does not reach" (0006's own Refinement Context
defect 1 says exactly this). The user's conclusion: lint lies. That is P5
("A green lint means resolution succeeds") failing on its first non-trivial
input, and it fails because the narrowing lives in a document the lint
implementer had no reason to read.

---

## 4. Premortem — written as if it already happened

*Six months on. The `flow` verb shipped in month four. This is the retrospective.*

We never broke the §JD-4 deadlock. The cluster-reconcile pass in late August
produced a fourth document recommending the same arm the other three already
recommended — 0003 records, 0006 cites — and returned `NEEDS_DECISION` because
no skill stage is authorized to edit `docs/jdr/0001-resolve-kernel-seam.md`. A8,
A10, and A12 stayed `Pending`. 0003 never left `Draft`.

In October, under delivery pressure, an implementer bypassed the cluster
entirely. They read `internal/resolve/resolve.go`, saw `Guard string` (still
`string` — JDR §D1's atom decision had never been written into code), and built
`internal/lint` against the shipped type. They did not read 0003; it was a
`Draft` and the README index said so. They implemented `graph-coverage-gap` and
`graph-unprovable-coverage` from 0006's taxonomy table, which was complete and
unambiguous, and they implemented row grouping the obvious way: rows sharing
`(source state, outcome)` — which happened to match 0003's definition by luck,
not by citation.

What they did not implement was the narrowing. `guard_unevaluable` appears zero
times in 0006's body, so there was nothing to implement. `lintGroup()` proved
the finite product covered and returned no finding.

The failure surfaced in the `foundational-to-cove` fixture — 0003's own
representative row, named in A7's `If wrong` clause. The row uses an `exists`
atom over an optional key. Lint proved the group exhaustive over the declared
enum product, having dropped the `exists` dimension as non-participating.
`Resolve()` then hit `evaluateGuard()`, got `GuardUnevaluable` from the seam for
the absent key, and returned `KindGuardUnevaluable` — `resolve.go:413`.

The user journey: a flow maintainer edits the RDR transition model, runs
`make check`, sees `intrastate lint` pass, commits, and pushes. The next
`intrastate flow next` in CI refuses with `guard_unevaluable` naming a key the
maintainer never mentioned. They re-run lint: still green. They file a bug
against lint. We closed it `working-as-intended` and pointed at RDR 0003 §
Technical Design line 652 — a `Draft` the maintainer had no reason to read, and
which by its own Status line was "the one record still gating lock."

The root cause in one line: **we knew the exact defect, wrote it down in four
places, and never assigned it to one.** 0006's Refinement Context named it as
defect 1 on 2026-08-21 and correctly predicted the trigger mismatch
("non-finite dimension" vs "fully-finite product whose participating row can
still refuse"). We had the analysis six months early. What we did not have was a
document with the authority to close it.

Secondary finding: the rehomed tag declaration model was moved back out of 0003
in November, into `internal/table`, when the normalizer needed value kinds at
parse time. Nobody updated 0003's 24 normative blocks; it is still `Draft`.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — §JD-4 has a recording document (catches S-1, 1.1)**
```gherkin
Given JDR 0001 §JD-4 is an entry in the Interface record
When any RDR in the cluster books an assumption whose closure route is "§JD-4"
Then §JD-4 MUST name exactly one recording document in its Resolved clause
And no cluster RDR may be marked Final or Draft-ready while its closure route
    points at a JDR entry that is still "Open"
```
*Result today: FAILS.* §JD-4 reads "Open only as to which document records the
narrowing". 0003 A8, 0007 A12 both route to it.

**AT-2 — Row-group definition has a floor (catches S-2, 1.2)**
```gherkin
Given RDR 0003 defines "scoped row group"
And RDR 0006 computes overlap and coverage over row groups
When the definition is traced from consumer to producer
Then the trace terminates at exactly one document
And the consuming document contains the defined term at least once normatively
```
*Result today: FAILS.* `row group` occurs 1× in 0006 (inside the Refinement
Context, not a normative clause) and 32× in 0003, whose A10 defers acceptance
back to 0006.

**AT-3 — Every quantifier has a stated contract (catches S-3)**
```gherkin
Given RDR 0003 states clauses quantifying over "every reachable predecessor"
Then the document owning graph traversal (RDR 0006) MUST state a
     predecessor-reachability contract in a normative block
```
*Result today: FAILS.* `reachable predecessor` occurs 3× in 0006, none in a
`normative` block; 0003 A12 is `Pending` on exactly this.

**AT-4 — No RDR asserts a stale peer status (catches S-4)**
```gherkin
Given RDR 0007 is Final and names peers with their statuses
When RDR 0006 is demoted from Final to Draft
Then RDR 0007's Predecessors block MUST be corrected before any implement stage
     reads it
```
*Result today: FAILS.* `0007:159` still says "**RDR 0006** (`Final`, tolerance
§JD-4)"; 0006 logged the staleness and declined to fix it.

**AT-5 — Lint's promise is stated where lint is implemented (catches S-6, §3)**
```gherkin
Given Principle P5 "A green lint means resolution succeeds"
And RDR 0007 defines a runtime veto refusing guard_unevaluable
When an implementer reads only RDR 0006 to build the lint
Then RDR 0006 MUST contain the narrowing or a normative citation to it
```
*Result today: FAILS.* `guard_unevaluable` occurs 2× in 0006 — line 10 (Status
line) and line 825 (Refinement Context, which says "delete on re-lock") — and
**zero times in any normative block**. 0006's own defect 1 states this: "a
Status-line citation standing in for absorbed text".

**AT-6 — A Resolved JDR decision reaches code before dependent RDRs lock
(catches S-7, 1.3)**
```gherkin
Given JDR 0001 §D1 is marked "Resolved: (d)" — Row carries parsed atoms
When 129 commits land after the decision
Then internal/resolve/resolve.go MUST NOT still declare `Guard string`
Or the decision MUST be marked "Resolved — not yet implemented"
```
*Result today: FAILS.* `resolve.go:185` is `Guard string`; last code commit
`1c9c0ca`, 2026-08-09.

**AT-7 — Final means no downgraded assumptions (catches S-9)**
```gherkin
Given RDR 0007 Status is Final
Then no Critical Assumption may carry Status "DOWNGRADED" whose home is an
     open JDR entry
```
*Result today: FAILS.* 0007 A12 is `DOWNGRADED`, home §JD-4, "re-confirmed open
this stage".

**AT-8 — The authority surface exists before the authority is claimed
(catches S-10)**
```gherkin
Given RDR 0006 makes `intrastate lint` the blocking acceptance authority
And A5 is Pending on an MVV that CI invokes it
When the command tree is searched for a "lint" command
Then a registered command MUST exist, or A5 MUST block lock
```
*Result today: FAILS (correctly Pending).* `grep -rn '"lint"' internal cmd`
returns no match; A5 `Pending`.

**AT-9 — Spec-to-code proportionality gate (catches S-12)**
```gherkin
Given a cluster of RDRs targeting one seam
When the cluster's RDR+evidence line count exceeds 10x the non-test code it
     governs
Then no further RDR stage may run until an implement stage lands
```
*Result today: FAILS.* 24,162 lines vs 1,236 — a ratio of 19.5:1.

---

## Appendix — verified measurements

Every number above, with its command.

| Measurement | Value | Command |
|---|---|---|
| 0003 body | 1,755 | `wc -l docs/rdr/0003-guard-predicate-exhaustiveness.md` |
| 0006 body | 880 | `wc -l docs/rdr/0006-graph-lint-authority-and-guarantees.md` |
| 0007 body | 2,402 | `wc -l docs/rdr/0007-guard-predicate-totality.md` |
| Three RDR bodies | 5,037 | `cat 0003 0006 0007 \| wc -l` |
| Evidence, 3 dirs | 19,125 | `find <3 evidence dirs> -type f -exec cat {} + \| wc -l` |
| **Set total** | **24,162** | sum of the two above |
| Evidence files | 36 / 12 / 55 | `find docs/rdr/<slug> -type f \| wc -l` |
| Non-test Go, all | **1,236** | `find internal cmd -name '*.go' -not -name '*_test.go' -exec cat {} + \| wc -l` |
| Non-test Go files | 8 | same, `\| wc -l` on file list |
| Test Go lines | 2,735 | `find internal cmd -name '*_test.go' -exec cat {} + \| wc -l` |
| `internal/resolve/resolve.go` | 576 | `wc -l internal/resolve/resolve.go` |
| Spec:code ratio | 19.5:1 | 24,162 / 1,236 |
| Last code commit | `1c9c0ca` 2026-08-09 | `git log -1 --format='%h %ad' --date=short -- internal cmd` |
| Commits since | 129 | `git log --oneline --since=2026-08-09 \| wc -l` |
| `lint` command in tree | 0 | `grep -rn '"lint"' --include='*.go' internal cmd` |
| `Row.Guard` type | `string` | `sed -n '183,185p' internal/resolve/resolve.go` |

Cross-reference counts (`grep -c` on each RDR body):

| Token | in 0003 | in 0006 | in 0007 |
|---|---|---|---|
| `RDR 0003` | 5 | 32 | 71 |
| `RDR 0006` | 71 | 2 | 3 |
| `RDR 0007` | 46 | **2** (lines 875/878, both in Refinement Context) | 1 |
| `JDR 0001` | 24 | 4 | 41 |
| `JD-4` | 28 | 8 | 8 |
| `guard_unevaluable` | 20 | **2** (line 10 Status, line 825 Refinement Context; **0 normative**) | 29 |
| `graph-unprovable-coverage` | 5 | 4 | 0 |
| `row group` (`-i`) | 32 | 1 | 0 |
| `atom` (`-i`) | 131 | 3 | 223 |
| `reachable predecessor` (`-i`) | 5 | 3 | 0 |

Assumption load (`grep -cE '^- \*\*A[0-9]+'` and `Status` greps):

| RDR | Assumptions | Pending | Verified | `normative` blocks | Status |
|---|---|---|---|---|---|
| 0003 | 15 | 5 | 10 | 24 | Draft |
| 0006 | 5 | 1 | 4 | 10 | Draft (demoted 2026-08-21) |
| 0007 | 26 | 0 | 22 | 11 | Final |

---

## Answers to the five posed questions

1. **Is the division of labour coherent?** No. It is circular on two axes.
   0003 defines the row group but books A10 on 0006's acceptance; 0006's A2
   defers the coverage derivation to 0003 and re-verifies against A10. Same
   shape for §JD-4: 0003 A8 waits on the venue, 0006 Direction 1 waits on the
   assignment. Neither cycle has a floor. (S-1, S-2, S-3)

2. **Is the set over-specified relative to shipped code?** Yes, by 19.5:1 —
   24,162 lines of RDR+evidence against 1,236 lines of non-test Go, none of it
   touched in 129 commits since 2026-08-09. The `intrastate lint` command that
   0006 makes the acceptance authority does not exist. (S-12, S-10)

3. **Has 0003 become the cluster's center of gravity, and is that a problem?**
   Yes and yes. It is not the largest body (1,755 vs 0007's 2,402) but it is the
   normative hub: 24 normative blocks, and inbound citation traffic of 32 (from
   0006) + 71 (from 0007) against its own outbound 5 references to itself. It
   now owns atom grammar, finite-domain exhaustiveness, **and** the tag
   declaration model. The problem is that it is a `Draft` with 5 of 15
   assumptions `Pending`, and a `Final` peer (0007) cites it 71 times. The
   cluster's floor is its least-settled document. (S-8, S-11)

4. **Does the `guard_unevaluable` narrowing need three documents?** No. It is
   one sentence requiring four coordination points (§JD-4 assignment, 0003's
   statement, 0006's citation, 0007's veto). That is the symptom of a bad cut:
   the design-time proof and the runtime veto are **the same predicate evaluated
   at two times**, split across documents by *when* it runs rather than by
   *what* it decides. A cut along "who owns the predicate's semantics" would put
   the narrowing in one place by construction. (S-5, S-6)

5. **Are decisions being routed around rather than made?** Yes — and the set
   documents this itself. Three separate texts recommend the identical arm for
   §JD-4 (JDR "Recommended arm", 0003 A8 `Plan`, 0006 Direction 1) and none
   takes it. 0007 A12 is `DOWNGRADED` rather than answered. 0006 logs 0007's
   stale `Final` claim and declines to fix it. 0003 A8 names *this pass* as its
   closure venue — so a deferral here is the failure, not a step toward it.
   (S-1, S-4, S-9)

---

## The one action that unblocks the set

Edit `docs/jdr/0001-resolve-kernel-seam.md` §JD-4 from **Open** to **Resolved**,
taking the arm all three documents already recommend:

> **Resolved.** RDR 0003 records the `guard_unevaluable` narrowing for the
> finite-domain product it proves; the clause names the refusing atom. RDR 0006
> cites RDR 0003 and reuses `graph-unprovable-coverage` for the withheld claim —
> no second code, no non-blocking tier. *Lands in 0003 (already stated, line
> 652 ff.) and 0006 (one citation in a normative block).*

That single edit closes 0003 A8, unblocks 0006 Direction 1, and lets A10/A12 be
answered by 0006's refine pass rather than by a fourth venue. It is the only
edit in this cluster that no RDR stage is authorized to make, which is precisely
why it has survived four passes.

**It requires a human decision** — §JD-4 is a joint decision entry in a JDR, and
this critique has no authority to close it.
