Model: claude-fable-5

# Critique — cli/0012 pre-lock (model B)

Verdict: BLOCK. The record ships a lint/runtime divergence it claims to
have closed, institutionalizes two comparison semantics for one declared
kind inside the kernel, and names a beneficiary that cannot reach the
defect. Ledger first; sections 1–5 follow at full length.

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | 0012:§key-discoveries, 0012:C5, 0012:A4 | C5 was narrowed to predicate literals, but graphlint reach feeds `[initial]`/`[rule.write]` int cells — still admitted permissively by `conformKind` (`007`) — into the seam against MATCH atoms. After C2, `atomAdmitsValue` parses (`007 == 7`, reachable); the kernel's `TagSet.matches` byte-compares (`"007" != "7"`, no match). The `matchlit` counterexample the resolve ruling rejected is reopened; Key Discoveries' "no model can author a value the two read differently" is now false; A4's Pending re-verify scope aims at literals, not held cells. | `intrastate lint` exit 0 on a model whose root is a dead end; `flow resolve` at that root refuses `flow-no-match`. | §1, premortem, AT-1 |
| C-2 | 0012:C2, 0012:D-canonical-spelling-at-load | The runtime now compares one declared `int` kind two ways: `[match]` atoms byte-equal in `TagSet.matches`, `[guard]` atoms parse-equal at the seam. The record rejects option (c)/(d) for "institutionalizing two comparison semantics for one declared kind" and then does exactly that across blocks. | Same view `n=07`: a row's `[match] n = "7"` fails while its `[guard.all] n eq 7` passes; `flow next` shows the row excluded-by-match on one line and candidate-by-guard on the next. | §1, §3, premortem, AT-2 |
| C-3 | 0012:F1, 0012:C5 | The record accepts a KNOWN silent GuardFalse→GuardTrue flip (`--tag iter=07` vs `iter eq 7`) as a residual and charts it to a successor, while A6's own If-wrong says an RDR premised on removing silent verdict changes cannot introduce them. The thesis and the accepted residual contradict inside one record. | `flow resolve --tag iter=07` plans the guarded row where the release before planned the fallback; no refusal, no note in output. | §3, premortem, AT-3 |
| C-4 | 0012:A6, 0012:§prerequisites | Prerequisites gate lock on "A1–A5 verified"; A6 (the held-side census, Pending, evidence "none yet") is excluded from the gate by construction, so the unbounded silent leg ships unmeasured. | After upgrade a persisted artifact holding `07` on an int dimension silently selects a different row than before; no release note, no scan. | §2, §3, AT-4 |
| C-5 | 0012:§problem-statement, 0012:MVV | The named beneficiary (rdr navigator) execs the CLI (`rdr-next`, `rdr-gate` pass `--tag k=v`), is not a `resolve.Resolve` library caller, declares zero `int` tags, and its `bool` tags are observed-provenance fed through `canonicalValue`, which refuses non-`true`/`false` upstream. The navigator cannot experience the defect or the fix; the motivating journey is fictional and the "Acceptance for the reporting consumer" paragraph describes a caller that does not exist. | Nothing — which is the defect: the fix's acceptance journey for its stated consumer is unexecutable, and a reviewer cannot tell whether any real consumer reaches the owned door with a typed tag. | §3, premortem, AT-5 |
| C-6 | 0012:C1, 0012:F6, 0012:S5 | "Nil-mapping disposition is LOUD" is false at the graphlint site: `atomAdmitsValue` returns `verdict != GuardFalse`, so GuardUnevaluable is SATISFIABLE. A missed/zero-value construction in reach makes every owned `eq`/`in` match atom reachable — lint goes green on dead ends, silently. | `lint` reports no `graph-*` findings on a model with an unreachable node; runtime dead-ends. | §1, premortem, AT-6 |
| C-7 | 0012:C4, 0012:A3 | C4 forbids per-row rebuilds at the CLI ("opposite of one model, one evaluator") but leaves graphlint's site — `atomAdmitsValue`, two frames below `matchSatisfiable`, called per node × row × held value — with no threading decision. Either `DeclaredKinds(m)` is rebuilt per atom evaluation (allocation storm) or three signatures widen, and the record does not say which. | `lint` measurably slower on large models, or the implementer threads a map through three functions with no contract backing the choice. | §1, AT-7 |
| C-8 | 0012:C5 | C5 also refuses non-canonical MATCH literals, but match atoms stay byte-compared. A model authored consistently non-canonical (`[match] n = "007"`, callers pass `--tag n=007`) works today; after C5 it refuses at load, and once the author rewrites the literal to `7`, the same callers' `--tag n=007` no longer matches. Load-time canonicalization of one side of a byte comparison is a regression pathway, not a safety net. | A model that loaded yesterday refuses `malformed_predicate_atom` today; after the suggested rewrite, `flow resolve` returns `flow-no-match` for the tag spelling every caller already uses. | §1, premortem, AT-8 |
| C-9 | 0012:C2, 0012:C3 | Parsing is `strconv.Atoi` (platform-width int). The C3 "integer-overflow held value" case has a GOARCH-dependent expected verdict; a held `3000000000` is GuardTrue/False on 64-bit and GuardUnevaluable on 32-bit. No width is fixed in the contract. | Conformance suite red on a 32-bit build; verdict differs by build target. | §1, AT-9 |
| C-10 | 0012:A2, 0012:S3 | The narrowed reflection assertion "no view-typed or runtime-valued state" is not expressible by reflection: `map[string]string` of declared kinds and `map[string]string` of held values are the same type. The proxy collapses to "no field of type `resolve.TagSet`", which a cached view trivially evades. S3's "fails if a runtime value is added" cannot be implemented as stated. | The test stays green for an evaluator that caches the view in a string map; REQ-34's structural guarantee is gone with nothing replacing it. | §2, AT-10 |
| C-11 | 0012:§approach | The record narrows a peer's Verified assumption (`0007:A17` "nor tag declarations at evaluation time") in its own prose and calls that an adjudication. Two records now state A17 differently; the joint-check fired only for kind tokens, not for this. | 7.1 cluster gate demotes 0012 or 0007 for declare-in-both; implementation stalls on a re-lock. | §2 |
| C-12 | 0012:C5 | C5's refusal reuses `malformed_predicate_atom` for a value `conformKind` itself calls well-formed. The category lies about the defect class, and 0005's envelope routing puts a canonicality nag in the malformed bucket. | `lint` reports `malformed_predicate_atom` for `n eq "00"`; the user reads "malformed" and checks TOML syntax, not spelling. | §1 |
| C-13 | 0012:S4 | Scenario 4 is a grep ("no zero-value construction survives in non-test code") that must also see `var`, `new`, and embedded-field forms. No executable oracle is named; `make check` has `graph-lint` but nothing for this. | S4 never becomes a test; the "mechanical backstop" is a code-review checklist item. | §5, AT-11 |
| C-14 | 0012:§consequences | "One package plus three call sites" contradicts the plan: 18 test sites in `internal/guard`, 4 `guardSeam()` fixtures in `internal/cli`, graphlint fixtures, 3 contract call sites, the REQ-71 pin, and the REQ-34 test. | Estimate off by an order of magnitude; the "one change" migration lands as a multi-day red-tree branch. | §2 |
| C-15 | 0012:C1 | `DeclaredKinds` OMITS keys whose `Kind` is `""`. Load does not obviously default `Kind`; if an owned tag can load without a kind, every `eq`/`in` over it becomes GuardUnevaluable — a silent behavior change on a previously-deciding guard, routed through the "defensive" key-absent arm A1 calls unreachable. | `flow-guard-unevaluable` on a guard that decided yesterday, with a key the model does declare. | §1 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The narrowed C5 reopens the exact counterexample the design was ruled on (C-1, C-2, C-8)

Root cause. `rulings.md` Q1 ruled option (a) WIDENED: "Reject non-canonical
spellings at every authoring site: guard literals, match literals,
`[initial]` values, and `[rule.write]` values," and Q3 ruled that "lint's
typed admission is safe ONLY because load has made typed-vs-bytes
unobservable." The provenance note is explicit that an adversarial pass
"REFUTED the shared-seam premise for match atoms and produced the
`matchlit` counterexample": the kernel byte-compares match atoms in
`internal/resolve/resolve.go::TagSet.matches` (`tv.value != w.Value`),
graphlint runs the same match atoms through the guard seam
(`internal/graphlint/reach.go::atomAdmitsValue`), and C2 makes those two
readings differ. The whole reason C5 exists is to make that difference
unobservable by canonicalizing EVERY value reach can feed the seam.

Then the 3amigo pass narrowed C5 to `(*loader).atom` — predicate literals
only — to protect `0024:REQ-7`/`REQ-9`'s permissive `[emit]` tests. But
reach's held values do not come from predicate literals. They come from
`m.Initial` (`reach.go:114`, `heldValues(m, tv.Key, tv.Value)`) and from
`writesOf(row)` (`reach.go:326`, `row.Writes` + `row.NextTags`). Both are
conformed by `conform(decl, "eq", members)` → `conformKind` → bare
`strconv.Atoi`, which admits `007`, `+5`, `-0`. `heldValues` sorts and
dedupes (`canonicalValues` is `slices.Compact(slices.Sorted(...))`); it
does not canonicalize numerals.

So after Phase 1 + Phase 4 as written: `[initial] n = "007"` (int, min
0, max 9), row `r` with `[match] n = "7"`. Lint: `atomAdmitsValue` →
`NewEvaluator` int arm → `Atoi("007") == Atoi("7")` → GuardTrue → `r`
reachable → no `graph-*` finding, exit 0. Runtime: `TagSet.matches` →
`"007" != "7"` → no match → `flow-no-match` at the root. That is
`matchlit`, verbatim, the fixture the record's own ruling used to reject
"lint going silent on a genuine dead end." The record still says, in Key
Discoveries, "reach's held values come from load-conformed
`[rule.write]`/`[initial]` cells, so no model can author a value the two
read differently." That sentence was true under the widened C5 and is
false under the narrowed one; nobody revised it when C5 moved.

The passage that enabled it. A4's Pending note says to "re-verify that no
`eq`/`in` int comparison reads a LITERAL that escapes the narrowed check."
Literals were never the leak. Held cells are. The re-verification is
pointed at the wrong population, so it will pass and the leak ships.
D-canonical-spelling-at-load compounds it: it still rejects option (c)
"excluding `atomAdmitsValue` from typed comparison" because that
"institutionalizes two comparison semantics for one declared kind" — and
then C2 institutionalizes exactly that between `[match]` (bytes) and
`[guard]` (parsed) inside the kernel itself. The record argues against a
split it has already built.

Symptom. Two user journeys. First: `intrastate lint --model m.toml`
reports clean on a model whose initial state is a dead end; the first
`flow resolve` refuses `flow-no-match`. Second, the runtime-internal
split: a view with `n=07` (owned reader, or `--tag n=07`, which C5 no
longer reaches) against a row carrying `[match] n = "7"` and a sibling
row carrying `[guard.all] n eq 7` — the match row is excluded, the
guard row is selected, and `flow next --all` reports the same key as
both "excluded" and "candidate" on adjacent lines with no reason token
that explains why one `7` is not the other `7`.

And the match-literal half of C5 makes it worse (C-8). C5 refuses
`[match] n = "007"` at load — but match atoms are byte-compared, so an
author who wrote `007` on both the literal and every caller's `--tag`
has a WORKING model today. After Phase 4 that model refuses
`malformed_predicate_atom`; after the author follows the diagnostic and
writes `7`, the callers' `--tag n=007` no longer matches. Canonicalizing
one side of a byte comparison is not a safety property; it is a
regression with a helpful error message.

### 1.2 "LOUD" is a property of two of the three sites, not the design (C-6, C-7, C-15)

Root cause. C1 fixes the nil-mapping evaluator's disposition as LOUD —
"never silently reverting to raw-string comparison" — and F6, S4 and S5
lean on that as the guarantee that a missed C4 site "surfaces as loud
`guard_unevaluable` refusals." That is true where the seam's verdict
becomes a refusal: `flow resolve` (kernel refuses `guard_unevaluable`)
and lint's product (`product.go::Denotation` maps GuardUnevaluable to
the blocking `AssignmentSet{}`). It is false at the third site.
`reach.go::atomAdmitsValue` returns `verdict != resolve.GuardFalse` —
by its own doc comment, "An UNEVALUABLE verdict is taken as
satisfiable: the relation over-approximates." A zero-value or wrong-model
evaluator in reach therefore makes every owned `eq`/`in` match atom
satisfiable, every node reachable, every coverage-gap invisible. The
"LOUD" mode at that site is a green lint. The record's own F6 scope note
admits a narrower version of this for `lt/gte/contains`; it does not see
that the `eq`/`in` case it calls loud is silent at the site that
matters most for lint's honesty.

The passage that enabled it. C1: "a nil-mapping evaluator answers
GuardUnevaluable for every `eq`/`in` atom … surfacing as
`guard_unevaluable` refusals." "Surfacing as refusals" is a claim about
the kernel path, generalized to all three sites without checking how
each consumes the verdict. S5 tests the seam in isolation and calls it
"the negative control for the zero-value-survives failure mode" — it
controls the evaluator, not the consumer that folds Unevaluable into
true.

Second root, same section. C4 decides threading for the CLI ("`probeRow`
runs in a loop, so widening it would rebuild the kinds map per row") and
is silent about graphlint, where `atomAdmitsValue` is called per node ×
row × held value, two frames below the only function that holds `m`.
The implementer either rebuilds `DeclaredKinds(m)` per atom evaluation —
the allocation the record just forbade at the CLI — or widens
`matchSatisfiable`/`ownedAtomSatisfiable`/`atomAdmitsValue`. A3's
"one added local parameter each" is wrong for graphlint: it is three
signatures or a package-level cache, and a package-level cache is the
mixed-model hazard F4 names.

Third root. C1 says `DeclaredKinds` omits keys whose `Kind` is `""`.
Either `Kind` can never be empty after load — then the clause is dead
prose that invites a reader to believe empty kinds exist — or it can,
and every `eq`/`in` over such a key now routes through the "defensive,
unreachable under A1" key-absent arm. A1 verified that undeclared KEYS
refuse at load; it says nothing about declared keys with no kind.

Symptom. A refactor that adds a fourth reach helper (JD-3's
`renderWrites` shim is the scheduled one) with `guard.Evaluator{}` passes
`make check`; lint stops reporting unreachable nodes; nobody notices
until a runtime `flow-no-match` on a model lint certified.

### 1.3 The fix cannot be demonstrated to its stated beneficiary, and the silent leg ships unmeasured (C-5, C-3, C-4)

Root cause. The Problem Statement motivates the whole record with "for a
consumer like the rdr navigator that is a silent misroute," and the MVV's
acceptance paragraph says "the rdr navigator … is a `resolve.Resolve`
library caller." Neither is true. `rdr/tools/rdr` does not import
intrastate (no `require` in its `go.mod`); `rdr/bin/rdr-next` and
`rdr/bin/rdr-gate` exec `intrastate flow resolve` and pass every fact as
`--tag k=v`. `rdr/models/rdr-status.toml` declares zero `int` tags (its
counts are enums: `["0", "1+"]`), and its `bool` tags
(`clustered`, `cluster_reconciled`, `accretion_disposition`, `lens_*`)
are `provenance = "observed"` — they enter through
`flow_input.go::canonicalValue`, whose `ConformValue` refuses anything
but `true`/`false` with `flow-tag-invalid` before the seam. The
navigator never touches the owned door. It cannot misroute on a
malformed typed value today, and it cannot observe the fix tomorrow.

The passage that enabled it. Q2's ruling correctly moved the MVV off
`--tag` and onto the owned reader, but the Problem Statement's
beneficiary sentence and the MVV's "library caller" paragraph were
carried across unchanged. The record now has a motivating consumer and
an acceptance journey that describe different programs.

The consequence is not cosmetic: with no real consumer identified at the
owned ingress, the only measured population is the 123-model authored
corpus (zero `eq`/`in` int atoms — S6 is "a coverage statement, not a
safety one," the record's words), and the held-side population A6 is
"none yet." Prerequisites say "All Critical Assumptions verified
(A1–A5)." A6 is not in the gate. So the record can lock with the silent
GuardFalse→GuardTrue leg (`"07"` matching `eq 7`) unbounded on the one
ingress it calls "not a narrow one — every persisted owned value and
every hand edit to the artifact enters through it."

And F1 makes the contradiction explicit: the CLI `--tag iter=07` flip is
"accepted here" as a "known, accepted residual," charted to a successor.
A6's If-wrong, three hundred lines earlier, says "an RDR premised on
removing silent verdict changes cannot introduce an unbounded set of
them." The record introduces a known one and defers its measurement.

Symptom. An operator upgrades. A persisted artifact holding `07` (hand
edit, older writer, whatever) on an int dimension now selects the guarded
row it used to skip. No refusal, no note, no scan. The report that
should have caught it — the navigator's — never ran through that path.

---

## 2. The section rewritten within six weeks: C5 (and D-canonical-spelling-at-load with it)

C5 will be rewritten because it is a compromise between two rulings that
cannot both hold, and the first real model with an `[initial]` int cell
proves it. The resolve ruling widened canonicality to every authoring
site precisely because `matchlit` showed lint going green on a dead end;
the 3amigo pass narrowed it to protect 0024's `[emit]` tests, and the
narrowing dragged `[initial]` and `[rule.write]` along even though
neither is what `TestReq7`/`TestReq9` pin — those tests exercise
`[emit.count]` and `[emit.code]` only. C5's "It is NOT homed in
`conformKind`'s shared int arm because that arm also serves `[emit]`,
`[initial]` and `[rule.write]`" treats three ingresses as one because
they share a function, not because 0024 fenced all three. Within six
weeks somebody writes `[initial] n = "007"`, lint says clean, runtime
says `flow-no-match`, and C5 gets a fourth home (`loadInitial`,
`renderWrites`) — or `conformKind` grows an `ingress` parameter, which is
the widened C5 by another name. Either way the "canonical? YES for
predicate atoms — deliberately NOT total" row in the `authority` table
and the "no model can author a value the two read differently" sentence
in Key Discoveries are rewritten, and D-canonical-spelling-at-load's
"zero violators holds only at the PREDICATE ingress" rationale goes with
them.

The second rewrite pressure on C5 is C-8: the match-literal half. Once a
user hits "my model refused at load, I fixed the spelling, now my
callers don't match," C5 either stops covering match atoms (reopening
lint-vs-kernel for match literals) or the kernel starts parsing match
atoms (breaking C1's "the resolution path reads no declarations"). Both
are C5 rewrites; one is a C1 rewrite too.

---

## 3. The assumption that will not survive first contact

A4, as restated: "lint and the runtime agree on every `eq`/`in` over an
`int` dimension." It is Pending, its re-verification scope is "no
`eq`/`in` int comparison reads a literal that escapes the narrowed
check," and it will be marked Verified because no LITERAL escapes. The
first user whose model has an int `[initial]` value spelled with a
leading zero, or a `[rule.write]` of `+1`, falsifies it in one command:
`lint` exit 0, `flow resolve` `flow-no-match`. The RDR already measured
this exact flip (`cover07`, exit 2 → exit 0) and built C5 to make it
"unauthorable" — but only on the literal side. The held side of lint's
own reach traversal is authored too, and it is not covered.

Runner-up: the Problem Statement's "for a consumer like the rdr
navigator." The navigator has no int tag and feeds every bool through
`--tag`. First contact is the implementer trying to reproduce the seed
defect with the beneficiary's model and finding there is no path.

---

## 4. Premortem

Written 2026-11-02, six weeks after lock.

RDR 0012 shipped on schedule. `guard.NewEvaluator` landed with
`DeclaredKinds`, C2's typed `eq`/`in` arms, and C5's canonical-spelling
refusal at `(*loader).atom`. `make check` was green. The seed kata
`intrastate#cq5p` was closed on the MVV: a reader returning `many` for
an owned `iter` int now refuses `flow-guard-unevaluable` with reason
`uncomparable`. That part works and nobody has complained about it —
partly because nobody has hit it: the only consumer anyone could name
at the owned door with a typed tag was the `mvv-int` fixture we wrote
ourselves. The rdr navigator, whose "silent misroute" the Problem
Statement was written around, turned out to pass every fact as `--tag`
through `canonicalValue` and to declare no int tag at all; it neither
had the bug nor saw the fix.

The first real report came from a workflow model, not the navigator.
Its author declared `[tags.retry] kind = "int" min = 0 max = 3`, wrote
`[initial] retry = "00"` because the column lined up, and had a row
`[match] retry = "0"` with `[rule.write] retry = "1"`. Under the old
raw-string reach, `"00" != "0"` and lint reported the row unreachable
from the root; the author had been ignoring that finding for weeks
because "it works when I run it" — the runtime also byte-compared and
also failed, but the author always invoked with `--tag retry=00` from a
script. After 0012, C5 did NOT refuse `[initial] retry = "00"`
(`loadInitial` → `conform` → `conformKind` → `strconv.Atoi` → fine), and
`graphlint/reach.go::atomAdmitsValue` now parsed `00 == 0`: the row
became reachable, the finding disappeared, lint went exit 0. The
runtime's `TagSet.matches` still byte-compared `"00" != "0"`. The author
read "lint clean" as permission, removed the script's spelling
workaround, and shipped a `flow-no-match` at the root to a downstream
consumer. When they filed it, the first triage comment was a link to
`rulings.md` Q1: this is `matchlit`, the fixture the record's own
ruling had used to reject exactly this outcome. The narrowed C5 had
reopened it and the Key Discoveries paragraph asserting "no model can
author a value the two read differently" had never been revisited.

The second report was uglier because it was inside one resolution.
A model with an owned reader supplying `stage=07`; row A carries
`[match] stage = "7"`, row B carries `[guard.all] stage eq 7`. The
kernel byte-compared A's match (`"07" != "7"`: not matched) and
parse-compared B's guard (`7 == 7`: true). `flow next --all` listed
`stage` as the excluding key on A and as a satisfied guard on B, with
0011's closed reason vocabulary offering nothing between `absent`,
`uncomparable` and `not-evaluated` to explain why one `7` was not the
other. The user's bug title was "int equality depends on which TOML
block I put it in." It does. D-canonical-spelling-at-load says the
record rejected option (c) because it "institutionalizes two comparison
semantics for one declared kind." C2 institutionalized them across
blocks and called the kernel's half "unchanged behavior."

The third came from the JD-3 shim. RDR 0030 Phase 2 extracted the
admitted-cell shim out of `renderWrites`, and the extraction — written
against 0030, not 0012 — constructed the evaluator with a bare
`guard.Evaluator{}`. Scenario 4's grep was never turned into a test
(nobody could say how a grep for `var`, `new`, and embedded-field zero
values would run under `make check`), and S5's negative control proved
the seam answers GuardUnevaluable without proving what the consumer does
with it. In the shim the Unevaluable verdict was folded to "admitted"
the same way `atomAdmitsValue` folds it to satisfiable — so a nil
mapping produced no refusal anywhere. It produced a lint that admitted
every write cell. F6 had promised the miss would "surface as loud
`guard_unevaluable` refusals." It surfaced as an unusually green CI.

Two smaller ones filled the tail. A 32-bit CI runner (arm, a contributor's
Pi) went red on C3's "integer-overflow held value" case because `Atoi`'s
width is the platform's and the case's expectation had been written on
amd64. And the REQ-34 reflection test, narrowed per A2 to "no view-typed
or runtime-valued state," was found to assert only that no field has
type `resolve.TagSet`; a reviewer demonstrated a green build with the
evaluator caching the whole view in a `map[string]string` named `kinds`.

The fix branch, six weeks in: `conformKind` grew an ingress parameter so
canonicality reaches `[initial]` and `[rule.write]` (the widened C5 the
resolve ruling had ordered and the 3amigo pass had unwound — 0024's
`[emit]` tests were never the obstacle; they only ever exercised
`[emit.count]`/`[emit.code]`). Match literals were pulled OUT of C5's
scope because canonicalizing one side of a byte comparison had broken
three working models. `atomAdmitsValue` got the evaluator threaded
through `matchSatisfiable`/`ownedAtomSatisfiable` and a comment
explaining why Unevaluable-as-satisfiable means "loud" is not a property
of this site. A6 was finally measured — after rollout, on the artifacts
that had already flipped — and the Problem Statement's navigator
sentence was replaced by the name of a consumer that actually reads an
owned int through an accessor. There was one.

---

## 5. Acceptance tests that would have caught each failure at review time

AT-1 (C-1) — Reach reads held cells C5 does not canonicalize
```
Given a model declaring [tags.n] int min 0 max 9
  And [initial] n = "007"
  And a row r with [match] n = "7"
When intrastate lint --model m.toml --as json
Then exit is 2 with a graph-* finding naming r unreachable
  OR load refuses n = "007" naming the canonical rewrite
And when intrastate flow resolve --model m.toml --outcome <o> is run at the root
Then its verdict agrees with lint's (no-match ⇔ unreachable)
```
Review-time form: for every ingress in the `authority` table's "every
authoring ingress" cell, state whether reach feeds it to the seam and
whether C5 covers it. Two of the five ingresses feed reach and C5 covers
neither.

AT-2 (C-2) — One kind, one comparison, inside one resolution
```
Given a view with n = "07" via an owned reader
  And row A with [match] n = "7", row B with [guard.all] n eq 7
When intrastate flow next --model m.toml --all --as json
Then A and B report the SAME disposition for key n
```
Review-time form: list every kernel comparison of a declared-int value
(`TagSet.matches`, seam `eq`, seam `in`, seam `lt..gte`) and its
reading (bytes/parsed). Two readings for one kind is a contradiction
row, not a Key Discovery.

AT-3 (C-3) — No accepted silent flip in a record whose thesis is "no silent flips"
```
Given the record's own A6 If-wrong ("cannot introduce an unbounded set")
When any Failure Mode paragraph contains "accepted" + "silently"
Then the record either refuses that input loudly or downgrades its thesis
```
Concrete: `flow resolve --tag iter=07` against `iter eq 7` must either
refuse or match with the SAME verdict as before the change. It does
neither.

AT-4 (C-4) — The gate covers the assumption that bounds the blast radius
```
Given A6 is Pending with Evidence "none yet"
When Prerequisites are read
Then A6 is in the list or the record states why the silent held-side leg may ship unmeasured
```

AT-5 (C-5) — The beneficiary can reach the ingress
```
Given the Problem Statement names the rdr navigator
When its model and invocation are inspected
Then it declares an int or bool tag of provenance owned read through an accessor
  And it calls resolve.Resolve or a CLI path that reaches OwnedSnapshot
```
Both fail today; the sentence and the MVV acceptance paragraph go.

AT-6 (C-6) — "LOUD" is checked per consumer, not per seam
```
For each C4 site:
  Construct guard.Evaluator{} (nil mapping)
  Evaluate an owned eq atom whose held value equals its literal
  Observe the SITE's output (refusal / AssignmentSet{} / reachability)
Then every site produces an observable failure
```
`atomAdmitsValue` produces "reachable." Fails.

AT-7 (C-7) — Threading is decided at every site, with the same rule
```
Given C4 forbids per-row rebuilds at probeRow
When atomAdmitsValue is inspected
Then the record names how the constructed evaluator reaches it (parameter / field / cache) and why that does not violate one-model-one-evaluator
```

AT-8 (C-8) — Canonicalizing a literal does not break a byte comparison
```
Given [match] n = "007" and every caller passes --tag n=007 (works on main)
When Phase 4 lands and the author applies the diagnostic's rewrite (write 7)
Then flow resolve with --tag n=007 still selects the row
```
Fails: byte compare of "007" vs "7".

AT-9 (C-9) — Overflow verdict is build-independent
```
Given the C3 overflow case
When run under GOARCH=386 and GOARCH=amd64
Then expected verdicts are identical (contract fixes width, e.g. ParseInt 64)
```

AT-10 (C-10) — The narrowed reflection proxy is implementable
```
Given the evaluator gains a field of type map[string]string holding HELD values
When the narrowed REQ-34 test runs
Then it fails
```
It cannot; the record must name a different structural proxy (e.g.
constructor is the only writer, field is unexported and of a named
non-string-map type) or admit the proxy is gone.

AT-11 (C-13) — Scenario 4 names its runner
```
Given S4 asserts "no zero-value construction survives in non-test code"
When make check runs
Then some target fails on a var / new / embedded-field zero value of guard.Evaluator
```
No target is named.
