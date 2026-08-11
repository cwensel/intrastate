Model: claude-opus-5[1m]
Delta-scoped re-run — 10 revised passages, iteration 2

Scope: only the ten rewritten passages named in the task. Prior-round findings
already adjudicated by the RDR are not re-litigated. All judgements grounded on
`main` against `internal/resolve/resolve.go`, its `_test.go` files, RDR 0002 /
0007 / 0009, and the three committed TOML fixtures.

---

## 1. A4 — reserved-key rule, both clauses swept + migration inventory (~line 424)

**Verdict: NEW-GAP**

The core fix lands. The prior defect — `Verified` stamped against only the
owned/observed clause — is genuinely closed. I confirmed both halves:

- *Second clause*: `Key: "recognized"` appears in Go table data exactly once,
  `internal/resolve/fixtures_test.go:237`, inside `recognizedTagSensitiveTable`'s
  `Row.Match`. That row's `RequiresOwned` is `["status"]` and its
  `NextTags`/`Writes` are `{Key: "status"}` (fixtures_test.go:239-241). Match
  position only. A4's claim is exact.
- *First clause*: the three declarations are the complete set. A `find`-driven
  sweep of every `*.toml` under `docs/` returns exactly three files carrying
  `provenance`, and each carries exactly one recognized-provenance declaration:
  `rdr-fixture.toml:25` `[tags.outcome]`, `kata-fixture.toml:22` `[tags.outcome]`,
  `guard-fixture.toml:36` `[tags.rewind_target]`. Line refs verified individually.
  "The complete set of committed recognized-provenance declarations" is true.

**The new gap is in the migration inventory the fix added.** The inventory reads:

> the three declarations above plus their reference sites — `[rule.match.outcome]`
> at `rdr-fixture.toml:59,74,86` and `kata-fixture.toml:47,57`.
> `[rule.match.<tag>]` is a *tag* predicate (`0002…md:272`), so a renamed
> declaration carries its references with it.

All five cited `[rule.match.outcome]` sites verify. But the inventory enumerates
reference sites **only for the two `[tags.outcome]` declarations**. The third
declaration — `[tags.rewind_target]`, which A4's own first clause counts as a
violation — has a reference site the inventory never lists:

    docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml:72
    [rule.guard.all.rewind_target]

This is not a stray. RDR 0002's field-layout block (`0002…md:266-274`, read on
`main`) puts `[rule.guard.all.<tag>]` in the same tag-predicate family as
`[rule.match.<tag>]`: "rule predicates live under `[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, and `[rule.guard.unless.<tag>]`". So the inventory's own
justifying sentence — a renamed declaration carries its tag-predicate references
with it — applies to `:72` with equal force, and the site is simply missing.

Concrete failure: Phase 2 says "rename the three committed recognized-provenance
declarations and their five `[rule.match.outcome]` reference sites (A4's
inventory)." An implementer executing that literally renames
`[tags.rewind_target]` → `[tags.recognized]` in `guard-fixture.toml` and leaves
`[rule.guard.all.rewind_target]` at `:72` pointing at a now-undeclared tag. That
fixture then fails RDR 0002's pre-existing `unknown tag` category
(`0002…md:345`) — a red test on RDR 0003's canonical guard evidence, introduced
*by* the migration, in a category this RDR does not own. The count "five
reference sites" is propagated into two more places that inherit the error: the
Decision Rationale blast-radius row ("a rename of three committed fixtures plus
five reference sites (A4)") and Trade-offs / Consequences ("five reference
sites"), plus the Risks bullet ("three declarations, five `[rule.match.outcome]`
reference sites").

Anchor: A4 **Migration inventory** bullet; propagated to Decision Rationale
blast-radius row, Trade-offs / Consequences bullet 2, Risks-and-Mitigations final
bullet, and Phase 2.

Repair: inventory is six reference sites across three fixtures, and the
guard-predicate positions belong in it. Note this also widens A9's reach — the
`Row.Match` vs `Row.Outcome` question A9 books is about `[rule.match.outcome]`
specifically, and `[rule.guard.all.rewind_target]` normalizes into `Row.Guard`,
a third field A9 does not currently contemplate.

---

## 2. A9 — `[rule.match.outcome]` → `Row.Match` or `Row.Outcome` (Pending, ~line 691)

**Verdict: CLOSED**

Booking this as `Pending` is the right disposition and it is booked honestly.
The underdetermination is real, not manufactured: RDR 0002's field-layout block
(`0002…md:266-274`) states only that "rule predicates live under
`[rule.match.<tag>]`" and nowhere maps a predicate to a normalized kernel field.
The kernel has both fields on `Row` — `Outcome string` and `Match []Tag`
(resolve.go `Row` struct) — and `Resolve` gates on them separately
(`row.Outcome != in.Recognized` at `:332`, then `view.matches(row.Match)`), so
both readings are structurally live. The cited fixture coordinates verify:
`rdr-fixture.toml:1` is the root `outcomes` alphabet, `:25` is `[tags.outcome]`
with `provenance = "recognized"`, `:59` is `[rule.match.outcome]`.

Status hygiene holds, which is where a `Pending` assumption usually breaks. I
checked every downstream consumer of the first reading and each one carries the
conditional rather than asserting the fact: Phase 2 says "**Do not run the
rename before A9 resolves**"; the Risks bullet says "mechanical **only under
A9's first reading**"; A4's own sub-question bullet says "this RDR must not
assert which." No settled-fact prose depends on it. That satisfies the
Finalization Gate's status-consistency rule.

The only defect touching A9 is inherited from A4's inventory (finding 1) —
A9's frame is `[rule.match.outcome]` alone, and the missed
`[rule.guard.all.rewind_target]` site raises a third normalization target
(`Row.Guard`). That is scored against A4, not against A9's own text.

---

## 3. A10 — accessor-produced owned key data-derived as `recognized` (Pending)

**Verdict: CLOSED**

This closes a real hole that A8's sweep could not reach, and it says so in the
right place. A8's scope-limit bullet now concedes "A literal sweep is
structurally incapable of refuting that case; it is governed by A10, not by this
assumption" — which is accurate. `Input.Owned` is documented in `resolve.go` as
"accessor-produced owned tag snapshot" (verified in the `Input` struct doc
comment), so a key spelled `recognized` can arrive with no Go literal anywhere,
and no `grep` over `*_test.go` can see it.

The load-bearing framing is the strongest part: "block 4 sends this breach down
the Go error path RDR 0001 reserves for *programmer* mistakes. If a key can be
data-derived, user data can be reclassified as a programmer mistake." That is
the correct hazard, correctly attributed to block 4, with a named repair
("scoping the precondition to producers the kernel can hold responsible, and
routing the accessor-origin case to a typed refusal or to RDR 0004's own
validation — a design change this RDR would have to make before lock").

Downstream consistency holds: Failure Modes carries the conditional explicitly
("Whether an *accessor-produced* owned key can spell `recognized` from artifact
data is open (A10); if it can, that origin must not travel the
programmer-mistake path"). No settled-fact prose depends on A10 resolving one
way. Verification plan names concrete artifacts (`rdr-fixture.toml:34-40`, RDR
0004's accessor binding) rather than gesturing.

---

## 4. A11 — composition/ordering vs RDR 0009's `Resolve`-entry precondition (Pending)

**Verdict: CLOSED**

The independence half is correct, and I verified it against `main` rather than
taking the RDR's word. The two predicates read genuinely disjoint fields of the
same `Row` values:

- This RDR (after block 5's widening): `Input.Owned`, `Input.Observed`, and each
  row's `RequiresOwned []string`.
- RDR 0009: each row's `Escape []RefusalKind` and `Writes []Tag`, predicate
  "exactly `len(row.Escape) != 0 && len(row.Writes) != 0`" (`0009…md`
  Load-Bearing Decisions / Selection-predicate, read on `main`).

These are four distinct fields on the `Row` struct (confirmed against the struct
definition in resolve.go). Neither predicate writes anything, so neither can
change the other's verdict — the RDR's claim survives contact with the source.
A11 correctly reduces the open question to *order only*, which is the accurate
residue. Its "If wrong" is properly consequential: "a doubly-breaching `Input`
reports different errors on different builds, and RDR 0005's exit-code mapping is
nondeterministic for it."

Also correct not to bind 0009's symbol name — I confirmed 0009's predicate form
is explicitly provisional ("form sharpened at Pre-Lock; e.g. a
ValidateTable-style function", `0009…md:510-517`), so A11's caution is grounded
rather than decorative. And 0009 does write the structurally identical clause
over the same entry point, so the collision A11 books is real.

The self-referential honesty in block 5 — "this widens the `Input` predicate to
read `in.Table.Rows`, so it is no longer decidable without reading the table, and
the independence half of A11 must be re-checked" — is exactly the cross-reference
that keeps A11 from going stale. It is present and points the right way.

---

## 5. A12 — 0002 implementer reading 0002 alone (Pending)

**Verdict: CLOSED**

This books a genuine structural risk that the RDR previously left implicit, and
the grounding is exact. I confirmed RDR 0002's `[tags.<tag>]` grammar block on
`main` (`0002…md:266-274`): it fixes the field layout and admits any
`[tags.<name>]`, carries no name constraint, and points at nothing. A12's
premise — "0002's own text does not carry" the constraint — is literally true.

The linkage to A1 is the part that makes this more than a filing note. A1's
scope-limit bullet correctly separates the two halves (the "including at
minimum" fence at `0002…md:342-348` licenses the new *category*; it does not
license the new *constraint on tag names*, which lands in the separate grammar
block at `0002…md:266-274`), and then hands the second half to the Overrides
field and to A12. I verified both cited blocks are in fact distinct normative
blocks in 0002 with different extensibility postures — the category block says
"including at minimum", the grammar block says "MUST use the Resolve spike field
layout" with no fence. The distinction A1 draws is real.

"If wrong" is properly severe and names the repair ("a pointer landing in 0002
(a cross-RDR edit) or an explicit implementation-ordering prerequisite here").
No settled-fact prose asserts the mechanism exists.

---

## 6. Normative block 2 — reference-position split, non-coverage, lower bound

**Verdict: CLOSED**

Three sub-fixes, all sound.

*Predicate vs write split.* The reasoning is correct and I checked the
categories exist. Predicate positions (`[rule.match.<tag>]`,
`[rule.guard.all.<tag>]`, `[rule.guard.unless.<tag>]`) need no separate
reserved-key check because a conforming model declares `recognized` and those
references resolve to it, while a non-declaring model already trips `unknown
tag`. `[rule.write]` is correctly carved out with the right reason: once
`[tags.recognized]` is legally declared the name is *known*, so `unknown tag`
stops firing and what rejects the write is `write to non-owned tag`. Both
categories verify in 0002's list at `0002…md:345`. The three predicate positions
match 0002's grammar block verbatim. This is the fall-through argument done
properly — it names which pre-existing rule catches each position rather than
asserting coverage.

*Explicit non-coverage.* The added paragraph is the strongest passage in the
delta. It constructs the surviving-symptom case concretely (`[tags.recognized]`
declared correctly, `[tags.outcome]` owned, `[rule.match.outcome]` written where
the recognized one was meant), walks each rule to show all pass, and states the
row still never fires. Then names it precisely: "not a naming defect: no name is
wrong, a reference points at the wrong declared tag." And closes with the
anti-misreading instruction — "a reader must not read block 2 as a guarantee
that a conforming model's outcome row fires." A normative block that fences its
own guarantee is doing what a normative block should.

*Lower bound.* Correct and well-grounded. "This RDR imposes **no lower bound**"
with the reason — matching the recognized outcome as a tag is an affordance, not
an obligation — is confirmed by the kernel: outcome gating is
`row.Outcome != in.Recognized` (resolve.go:332) and never reads the tag key, so
a model gating only on `Row.Outcome` is fully functional with no
recognized-provenance declaration. Requiring the declaration would indeed break
those models. The consequence is stated rather than left implicit ("lints clean
and refuses `no_match` at resolve time") and routed to the same intent-channel
gap the Failure Modes section charts.

The upper-bound argument also holds: because `[tags.<tag>]` is name-keyed, one
model admits at most one declaration named `recognized`, and a second
recognized-provenance declaration necessarily carries a different name and fails
the naming rule. Cardinality as a *consequence* of naming rather than a separate
check is the correct derivation, and the duplicate-key case is properly routed to
`malformed TOML`.

---

## 7. Normative block 4 — RDR 0009 collision paragraph + unconditional justification

**Verdict: OPEN**

The unconditional-on-`Input.Recognized` justification is fully closed. It states
the cost plainly rather than assuming it away ("when `Input.Recognized` is empty
`assemble` injects nothing (`internal/resolve/resolve.go:151`), so there is no
collision to prevent and the error is a *reservation* being enforced, not a
shadowing being averted"). I verified `:151` — `if in.Recognized != ""` is
exactly the guard, so the concession is accurate to the line. The reason for
keeping it unconditional is principled and not merely aesthetic: "a key whose
reservation lapses per-call is not reserved: a producer would have to know
whether an outcome is in flight to know whether its own tag key is legal, which
is exactly the cross-side coupling this RDR exists to remove." And it correctly
disclaims symmetry-with-block-1 as the reason, pinning scenario 6 instead.
Nothing survives here.

**What survives is in the 0009 collision paragraph.** It issues a MUST that
nothing in the RDR can satisfy at implementation time:

> This RDR binds only what it owns: its breach MUST NOT mask or be masked by a
> table-shape breach — a caller breaching both is entitled to a deterministic,
> reproducible error rather than whichever check the implementer sequenced first.

then, ten lines later:

> The **total order** between them is a cross-RDR decision neither RDR may take
> unilaterally; it is reconciled at Stage 7.1 (A11). Absent that reconciliation
> an implementer MUST NOT invent an order silently.

For two checks that both return a single non-nil `error` from one `Resolve` call,
"neither masks the other" *is* a statement about order — whichever runs first and
returns is the one the caller sees, and the other is masked by construction.
`Resolve`'s signature is `(Result, error)` (verified on `main`), a single error
return with no join, and block 4 itself requires "returning a non-nil error and
no `Result` disposition on breach". So the only ways to satisfy non-masking are to
fix an order, or to report both (an aggregate error) — and block 4 forbids the
first to the implementer while never licensing the second.

The result is that an implementer who reaches block 4 with A11 still `Pending`
is handed two MUSTs that cannot both be discharged: implement the predicate
(block 4 requires the `Resolve` entry call), do not let either breach mask the
other, and do not invent an order. The only conforming action is to not implement
— which contradicts Testing Strategy's "Done for this RDR's own implementation"
listing scenario 6 as green-at-HEAD, and contradicts MVV's kernel half, which is
declared "executed during this RDR's implementation."

This is narrower than the pre-revision finding — the independence half is now
correct and verified (see finding 4), and the deferral to A11 is legitimate in
itself. What is unresolved is that block 4 states the non-masking requirement as
a binding MUST *on the implementation* while deferring the only mechanism that
could satisfy it to a stage that has not run. Either the non-masking sentence
should be demoted to a constraint on the eventual Stage 7.1 reconciliation
(where it belongs), or block 4 should license the interim behavior explicitly —
e.g. that until A11 resolves, `Resolve` MUST report both breaches, or that this
RDR's kernel half MUST NOT land before A11 closes (which the Testing Strategy's
own cluster-read caveat already gestures at for a different reason).

Anchor: normative block 4, the sentence "its breach MUST NOT mask or be masked
by a table-shape breach", against the same block's "an implementer MUST NOT
invent an order silently" and Testing Strategy's placement of scenario 6 in the
runnable-at-HEAD half.

---

## 8. Normative block 5 — enforced by the same exported `Input` predicate (widened to `RequiresOwned`)

**Verdict: CLOSED**

The finding this was written for — block 5 was an unfalsifiable MUST discharged
by a derivation argument conditional on a `Pending` assumption — is closed, and
closed the right way. The added paragraph states the principle explicitly ("An
obligation discharged by construction still needs an artifact that fails when the
construction changes, or it is text nothing executes") and then supplies the
artifact: the block-4 predicate "MUST also reject a `Row.RequiresOwned` entry
naming `recognized`" on the rows of the supplied table. That converts a
derivation into a check, at genuinely no new surface — it is the predicate block
4 already mandates.

Conditionality is now handled correctly throughout. Block 5 says "**if A7
holds**, the obligation is therefore discharged by construction on today's
derivation path", and the Approach item 5 mirrors it ("**Conditional on A7
(`Pending`)**... that discharge is a consequence of A7, not an independently
established fact"). The Capability Dependencies row matches ("Discharge on the
current derivation path is **conditional on A7**, not established"). So the
`Pending`-assumption-feeding-settled-prose defect is gone — the predicate carries
the obligation unconditionally and A7 only governs whether a *second*
source-lint check is additionally needed.

**On the specific question of whether the widening contradicts block 4:** it does
not. I checked this against source rather than against the RDR's assertion. The
widened predicate reads `in.Table.Rows[].RequiresOwned []string`; RDR 0009's
reads `in.Table.Rows[].Escape []RefusalKind` and `Writes []Tag`. Those are three
distinct fields on the `Row` struct (verified against the definition in
resolve.go). Both predicates are pure reads. So block 4's claim — "neither can
change the other's verdict and both are decidable in one pass" — remains true
after the widening, and block 4's parenthetical already states the widened field
list correctly ("this one: reserved-key occupancy in
`Input.Owned`/`Input.Observed` and in each row's `RequiresOwned`"). No
contradiction. Traversing the same `in.Table.Rows` slice is shared *iteration*,
not shared state, and does not couple the verdicts.

Block 5 also flags the consequence itself rather than leaving it for a reviewer
to catch ("this widens the `Input` predicate to read `in.Table.Rows`, so it is no
longer decidable without reading the table, and the independence half of A11 must
be re-checked against RDR 0009's row-shape precondition rather than assumed"),
and A11's verification plan picks up exactly that thread ("so neither can change
the other's verdict even though both now traverse `in.Table.Rows`"). The
cross-reference is bidirectional and correct.

The kernel-side residual is accurate: `missingOwned` iterates
`row.RequiresOwned` and tests `view.has(key, ProvenanceOwned)` (resolve.go:444-453),
and `TagSet.has` is provenance-specific, so a key present under
`ProvenanceRecognized` fails the owned-only test and is reported missing,
yielding `owned_state_unavailable` naming the reserved key. Verified on `main`.
Scenario 9's two halves pin both sides.

(The ordering defect that touches this predicate is scored against block 4 in
finding 7, not here — block 5's own text correctly routes it to A11.)

---

## 9. Normative Contracts preamble — "closure argument for three channels"

**Verdict: CLOSED**

This is a real closure argument, not a more confident assertion. I tested it the
way it invites: it claims the channels are "derived from the code, not enumerated
by inspection of this document," so I checked the code.

The argument's structure is a disjunction over how a key can become visible to
rule evaluation: either it enters the assembled view, or it is named in a row
field the kernel reads against that view. Both legs verify.

*View-entry leg.* `assemble` is the sole constructor of `TagSet` in the resolve
path, and it writes exactly three sources — `in.Recognized` (`:151-156`),
`in.Observed` (`:157-159`), `in.Owned` (`:160-162`). The cited range `:151-162`
is exact. `Resolve` calls `assemble` once (`view := assemble(in)`) and threads
that single value to the matcher, the guard gate, and `escapeOrRefuse`; `TagSet`
wraps a map so pass-by-value shares one backing map. Three write-sources, no
fourth.

*Row-field leg.* This is the part that could have been hand-waved and is not.
The claim is that `RequiresOwned` "is the one row field read against the view *by
provenance* rather than by value, which is what makes it a separate channel." I
swept every provenance-sensitive read in the package: `view.has(key,
ProvenanceOwned)` at `resolve.go:449` inside `missingOwned` is the only one.
Every other view read is value-based (`matches`) or provenance-blind. And
`missingOwned` is driven exclusively by `row.RequiresOwned`. So the claim that
`RequiresOwned` is the unique provenance-sensitive row field is literally true at
HEAD, not merely plausible.

I also checked the fields the argument implicitly excludes, since an incomplete
disjunction is how this kind of argument usually fails. `Row.NextTags` and
`Row.Writes` are consumed only by `planOf`, which copies them into the `Plan`
output — they are never read against the view, so they are correctly outside the
enumeration. `Row.Outcome` is compared to `in.Recognized` as a raw string
(`:332`, `:479`) and never routes through `TagSet` — correctly outside.
`Row.Match` and `Row.Guard` are view reads but value-based, and the argument
folds them into the declaration channel with a valid reason ("a reference must
resolve to a declared tag"), consistent with block 2's fall-through. `Row.Escape`
is 0009's field, not a view read.

The falsifiability condition is stated and is the right one: "A fourth channel
would require either a second view constructor or a new row field the kernel
reads against the view — both of which are kernel changes this RDR would have to
re-open anyway." That is a genuine closure condition — it names what would have
to change for the enumeration to break, and both are changes this RDR would see.

The preamble also earns the meta-honesty it opens with ("not an assertion — the
enumeration was extended twice under review, so it owes one"), and correctly
keeps the split signal untripped: `foundational` on the cross-RDR axis, not on
contract count. That matches the Metadata Profile field.

---

## 10. Testing Strategy — binding consequences of the Done split

**Verdict: CLOSED**

The prior defect was that the Done split let this RDR reach "implemented" with
its user-facing outcome unshipped and left that as an accounting fact a reader
had to notice. The revision makes it binding and self-disclosing.

The disclosure is blunt and correct: "every scenario that delivers the Problem
Statement's outcome — the typed load/lint failure replacing the silent no-match —
is in the carried half. This RDR's own half ratifies the key, adds the
`Input`/`RequiresOwned` predicate, and pins view plumbing; none of it is
perceptible to a table author. So this RDR can reach 'implemented' with its
user-facing outcome still unshipped, and that is a real property of the split,
not an accounting artifact." That is accurate — I checked the scenario
allocation. Scenarios 2, 4, 5, 7 and the MVV lint half all require a normalizer,
which does not exist at HEAD; scenarios 1-kernel-half, 3, 6, 8, 9 are all
kernel-side and run against `internal/resolve` as it stands.

The two consequences are genuinely binding rather than advisory, which is the fix:

1. "the Prerequisite 'RDR 0002 implementation underway' is a **gate on declaring
   this RDR's problem solved**, not merely on running the carried tests: this
   RDR's Close MUST NOT claim the Problem Statement outcome until the carried
   half is green." This converts a checklist item into a Close-gate obligation.
2. The implementation-ordering fact is routed to Stage 7.1 "because it determines
   whether locking 0008 before 0002 is implemented buys anything," with a named
   honest disposition if the answer is no ("hold this RDR at Final-unimplemented
   until 0002's work starts rather than implement a kernel predicate in
   isolation").

The split is also correctly distinguished from deferral in the Finalization
Gate's sense ("they are *scoped to the peer that owns the code*"), which matters
because the Gate's Scope Verification item would otherwise catch the carried half
as deferred MVV. That distinction is legitimate here: the code genuinely lives in
a peer, and this RDR's Normative Contracts are the authority 0002's
implementation prompt extracts.

Scenario-level grounding spot-checked and holds: scenario 6's reference to
`resolve_test.go:748` is correct — `TestReq20_EmptyInputTupleStillYieldsAValueDisposition`
calls `resolve.Resolve(resolve.Input{})` and fails on a non-nil error, so the new
precondition must keep it green, and the RDR says so. Scenario 3's claim that the
existing guard seam discards the view is correct
(`fixtureGuards.Evaluate(guard string, _ resolve.TagSet)`), making that coverage
genuinely net-new, and the RDR correctly requires a *new* seam alongside rather
than a change to it.

---

## Summary

| # | Passage | Verdict |
| --- | --- | --- |
| 1 | A4 — both-clause sweep + migration inventory | NEW-GAP |
| 2 | A9 — `[rule.match.outcome]` normalization target | CLOSED |
| 3 | A10 — accessor-derived `recognized` key | CLOSED |
| 4 | A11 — composition/ordering vs 0009 | CLOSED |
| 5 | A12 — 0002-alone implementer reach | CLOSED |
| 6 | Block 2 — reference split, non-coverage, lower bound | CLOSED |
| 7 | Block 4 — 0009 collision + unconditional justification | OPEN |
| 8 | Block 5 — enforced by widened `Input` predicate | CLOSED |
| 9 | Preamble — three-channel closure argument | CLOSED |
| 10 | Testing Strategy — binding Done-split consequences | CLOSED |

Two items to repair before lock:

- **A4's migration inventory** is short one reference site
  (`guard-fixture.toml:72`, `[rule.guard.all.rewind_target]`) and the "five
  reference sites" count is propagated to four other passages. Executing Phase 2
  as written leaves an undeclared-tag reference in RDR 0003's canonical guard
  fixture.
- **Block 4's non-masking MUST** is not dischargeable while A11 is `Pending`,
  because block 4 simultaneously forbids the implementer from fixing an order and
  never licenses reporting both breaches. Demote it to a constraint on the Stage
  7.1 reconciliation, or license the interim behavior.

On the two questions posed specifically: the block 5 widening does **not**
contradict block 4 — the predicates read disjoint `Row` fields and shared slice
iteration does not couple their verdicts. The preamble's closure argument
**does** close over the channels — its row-field leg is verified unique at HEAD
(`view.has(key, ProvenanceOwned)` at `resolve.go:449` is the only
provenance-sensitive read in the package) and its falsifiability condition is
real.
