Model: claude-opus-5

# Critique Diff — RDR 0003 Guard Predicate Exhaustiveness (iteration 2)

Reconciles two independent hostile critiques of the same draft:

- **Pass A** — `critique.md`, model `claude-opus-5`, 14 rows (A/C-1..C-14)
- **Pass B** — `critique-modelB.md`, model `claude-sonnet-5`, 10 rows (B/C-1..C-10)

`C-N` IDs do not correspond across the two files. Everything below is
reconciled **by RDR passage anchor**. Verification was performed
independently against the source documents by this pass; neither critique's
prose was taken on trust.

Headline: the two passes agree on **four** passages, and those four are the
hotspots. They disagree in substance on **two**. Pass A's six harsh claims —
the ones that attack `Verified` stamps Pass B accepted — verify **five
VERIFIED, one PARTIAL, zero REFUTED**, with one important scoping correction
inside R-01 and one inside R-14.

---

## Reconciled ledger

Severity key: **BLOCK** = must close before lock; **DEFER** = real defect,
closable after lock or at cluster reconcile; **NOISE** = correct observation,
no action warranted.

| RID | Passage anchor | A-id | B-id | Bucket | Verification verdict | Severity |
|-----|----------------|------|------|--------|----------------------|----------|
| R-01 | A2 Evidence (`0003:135`) "finite domains … enum/boolean values as declared sets … bounded integers as `{min..max}`", `Status: Verified` | C-1 | — | SOLO-A (B never reached) | **VERIFIED** — RDR 0002 declares no domain field of any kind | **BLOCK** |
| R-02 | §Technical Design `0003:374-377` "RDR 0006 supplies the graph-lint grouping context" | C-2 | — | SOLO-A (B never reached) | **VERIFIED** — `row group` appears 0× in RDR 0006; 0006 A2 delegates back to 0003 A2 | **BLOCK** |
| R-03 | §Technical Design `0003:378-383` default-on / "opts out by leaving a dimension undeclared" | C-3 | C-1, C-4 | **AGREED (hotspot)** | **VERIFIED** — direct contradiction with `disposition` row `0003:554` | **BLOCK** |
| R-04 | §Normative Contracts `0003:495-501` "can refuse" / `present-for-every-reachable-predecessor` | C-4 | — | SOLO-A (B never reached) | **VERIFIED** — 3 uses in 0003, 1 in 0006, 0 definitions anywhere; **but** see V-note: A's "RDR 0006 has no A6" is right, and the line-512 cite is a mis-cite in the 3amigo record, not in the draft | **BLOCK** |
| R-05 | A6 (`0003:162-167`) `Verified` on provenance-declaration evidence | C-5 | — | SOLO-A (B never reached) | **VERIFIED** — evidence proves provenance is declared, not that write-reachability is decidable | **BLOCK** |
| R-06 | A1 Evidence (`0003:128`) `Verified`, Method: Spike, `check.sh` | C-6 | — | SOLO-A (B *partially* considered — B/C-8 cites the same spike's limits for `contains` but accepts A1) | **VERIFIED** — `check.sh` is two `awk` passes over operator spellings; evaluates no predicate | **BLOCK** |
| R-07 | §Technical Design `0003:369-373` + Identity `0003:584-592`; `block` survival through RDR 0002 normalization | C-7 | — | SOLO-A (B never reached) | **PARTIAL** — `0002:305-309` does say "combine both into one candidate-row predicate set", but 0003 `:371-373` explicitly concedes normalization "may combine them into one internal constraint object"; the erasure is *not contracted*, it is *unguaranteed* | **BLOCK** (as a producer-request, not as a contradiction) |
| R-08 | Entire RDR — escape rows absent from coverage identity and `disposition` table | C-8 | — | SOLO-A (B never reached) | **VERIFIED** — `grep escape` over 0003 = 4 hits, none about rows; `0006:277-279` counts "a declared escape row" as coverage | **BLOCK** |
| R-09 | §Normative Contracts `0003:475-482` withheld claim "MUST name … the atom"; §authority `0003:538` | C-9 | — | SOLO-A (B raised the *adjacent* A8 binding at B/C-3, not this payload binding) | **VERIFIED** — `0006:352-356` lists no atom-level field | **DEFER** |
| R-10 | §Prerequisites `0003:932` `- [x] A1-A6 verified`; §Assumption Verification `0003:1104` | C-10 | — | SOLO-A (B never reached) | **VERIFIED (meta)** — §Assumption Verification re-asserts label hygiene, does not re-test Evidence | **DEFER** (mechanism, not a defect in itself) |
| R-11 | §MVV `0003:960-969` + Scenario 6; §Scope Verification `0003:1148` | C-11 | (C-8 partial) | SOLO-A, **overlaps** B/C-8 on the `contains` arm | **VERIFIED** — MVV needs optionality (A7), element universe (A9), and domains (R-01); gate says "in scope, not deferred" | **BLOCK** |
| R-12 | §Metadata Status `0003:9-12` + §Prerequisites `0003:946-953` — four external blockers, no fallback | C-12 | C-10 | **AGREED (hotspot)** | **VERIFIED** — 3 of 5 prerequisites unchecked; Phase list carries no gate | **DEFER** |
| R-13 | §Technical Design `0003:354-360` atom's "four conceptual fields … source identity" | C-13 | — | SOLO-A (B never reached) | **VERIFIED** — `0007:1260-1265` fixes `Key`/`Operator`/`Literal`/`Block`; 0003 substitutes `source identity` for `Block` | **BLOCK** |
| R-14 | §Approach operator matrix `0003:334-338` `in`/`contains` set literals; A9 | C-14 | C-8 | **CONTESTED** (see §CON-2) | **VERIFIED with correction** — `0007:1585-1595` does file the canonical-spelling duty on 0003 and does name `in`, but the exact phrase A quotes ("does not yet [have a byte-comparable spelling]") is a paraphrase; the substance holds verbatim | **BLOCK** |
| R-15 | A7 (`0003:171-200`) + Implementation Plan Phase 1 — projection lives in a side evidence file, no phase gate | — | C-2 | SOLO-B (A booked A7 as "correctly open", declined to attack it) | **VERIFIED** — `evidence/research/iter-2-projection-derivation.md` is cited only from A7; Phase 1 (`0003:978-981`) does not mention A7 or `exists` | **DEFER** |
| R-16 | A8 (`0003:1128-1136`) — lock blocked on `Final` RDR 0006 via an unowned venue | — | C-3, C-7 | SOLO-B (A treats A8 as "correctly booked"; explicitly declines at A `§closing`) | **VERIFIED** — `0006:344-350` is scoped to the non-finite case only; nothing in Phases 1-4 gates on A8 | **DEFER** |
| R-17 | Identity `0003:584-592` — atom identity vs. semantic equality `(tag, operator, literal)`; no clause says which operation uses which | — | C-5 | SOLO-B (A read the same passage for R-07/R-13 and never asked this question) | **VERIFIED** — `0003:592` states both equalities; no normative clause or MVV scenario selects between them | **BLOCK** |
| R-18 | §Technical Design `0003:398-402` "too large for the implementation to prove deterministically" — no threshold | — | C-6 | SOLO-B (A never reached) | **VERIFIED** — no bound, complexity class, or heuristic anywhere; Performance Expectations concedes "the spike uses four rows" | **DEFER** |
| R-19 | §Failure Modes — `guard_unevaluable` vs. predicate-parse-kind rendering to the author | — | C-9 | SOLO-B (A never reached) | **VERIFIED (as stated)** — the distinction is reconstructible only across three RDRs; but RDR 0005 owns the envelope by A4, so this is correctly out of scope | **NOISE** |
| R-20 | §Approach / A2 — is A2's central defect the *missing producer* or the *unbounded product size*? | C-1, §3 | C-6, §3 | **CONTESTED** (see §CON-1) | Both diagnoses are independently **VERIFIED**; they are not alternatives | A: **BLOCK** / B: **DEFER** |

---

## Contested items

Two genuine substantive disagreements. Both are worth surfacing because in each
case one pass's evidence bears on the other's conclusion.

### CON-1 — "The assumption that will not survive first contact" is A2 in both passes, for opposite reasons

Both passes nominate **A2** in their §3. They then diagnose it incompatibly:

- **Pass A** says A2's math is *correct* and its inputs are *nonexistent*: "A2
  is a valid derivation over inputs that do not exist." The failure is that
  `declared` domains cannot be declared (R-01). Under this reading a user
  never hits a product-size problem, because no product is ever built.
- **Pass B** says A2's math is *correct* and its inputs are *unbounded*: the
  first real author writes `6×8×12×50 = 28,800` and hits an undocumented
  "too large" threshold (R-18). Under this reading the product is built
  routinely and the cliff is performance.

These are not two readings of one defect; they are two defects, and **Pass A's
finding strictly precedes Pass B's in time**. If R-01 is true, no author can
declare a domain at all, so no author reaches B's 28,800-element product — B's
scenario silently assumes the very `domain = [...]` / `min` / `max` schema that
R-01 proves does not exist in RDR 0002. Pass B's own §3 quotes tag kinds
("enum/bool/set/bounded-int") as though they carry domains; verification shows
they do not.

**Resolution.** Pass A is right about ordering; Pass B is right that the
threshold gap is real and survives R-01's fix. R-01 blocks lock. R-18 becomes
live the moment R-01 closes and should be booked now so it is not discovered
after the producer request lands. Neither dismisses the other.

Note also that Pass B's §3 evidence is *derived from the spike fixture*
(`guard-fixture.toml`, which carries `domain = ["Draft","Final","Implemented"]`
and `min = 0 / max = 3`) rather than from RDR 0002. This is exactly the
circularity Pass A names at `critique.md:72-75`: the fixture invented a schema,
and both the draft and a downstream critic read it back as if it were the
producer contract. Independently confirmed — the fixture's `domain`/`min`/`max`
keys appear nowhere in RDR 0002.

### CON-2 — Is the `contains`/`in` gap one gap (A9) or two?

- **Pass B (C-8)** treats it as **one** gap: `contains` has no worked example
  and no element-universe declaration, so the "closed five-operator
  vocabulary" claim should be scoped to "four verified, one
  specified-pending-fixture." B's remedy is documentary — add an illustrative
  TOML fragment and re-scope the closed-vocabulary claim.
- **Pass A (C-14)** asserts B's framing **understates it**: the element
  universe (A9, `contains`-only) and the *canonical literal byte-encoding*
  (`in` **and** `contains`) are distinct duties from distinct producers, and
  the second is filed on RDR 0003 by a `Final` peer.

Verification favors Pass A. `0007:1585-1595` reads verbatim:

> for the SET-valued literals RDR 0003 gives `in` ("non-empty typed scalar
> set") and `contains` ("non-empty typed element set") it does not yet —
> `resolve.Tag.Value` is a bare `string` and no element encoding is declared
> (Phase 3). Until 0003 declares it, two `in` atoms on one key differing only
> by literal-set ORDER are not distinguishable by this tuple, so the kernel
> MUST compare the literal as the exact bytes the normalizer emitted and
> **0003's declaration MUST make that spelling canonical.**

That is a MUST assigned to this RDR by a `Final` document, naming `in`
explicitly. RDR 0003's Capability Dependencies (`0003:645`) book only "Set-valued
element encoding (declared element universe for `contains`)" — `contains` only.
So the `in` half is genuinely unbooked, and B's proposed remedy (a `contains`
example) would not discharge it.

**One correction to Pass A**, recorded for fairness: A's ledger row asserts
RDR 0007 "states the consequence against its own contract" and quotes with a
bracketed interpolation. The interpolation is A's own gloss, not 0007's words.
The substance — the duty, the naming of `in`, the order-sensitivity — is
verbatim. **Not refuted**; the quotation is loose, the claim is sound.

**Resolution.** Book two rows: A9 (element universe, `contains`, producer =
RDR 0002) and a new A10 (canonical set-literal byte spelling, `in` +
`contains`, owner = **this RDR**, filed by RDR 0007). The second is closable
inside this document today and is a lock blocker because a `Final` peer's
determinism claim rests on it.

### Near-contest, resolved: R-03's failure direction

Both passes hit `0003:378-383`, and their symptom predictions are opposite:

- **Pass A (C-3)** predicts **silent green** — every dimension is undeclared,
  so every group opts out, so lint proves nothing and returns exit 0.
- **Pass B (C-1)** predicts **surprise blocking** — a previously-green group
  starts failing `graph-unprovable-coverage` on a dimension the author never
  touched.

This is not a contest; it is the **same defect seen from both sides**, and its
existence confirms the defect. Verified at source: `0003:378-383` makes an
undeclared dimension the opt-out (no finding), while the `disposition` table at
`0003:554` makes the identical input "**Loud** — Inability-to-prove finding".
Two clauses, same input, opposite outcomes. Whichever an implementer picks, the
other pass's user is the one who files the bug. This is the strongest
AGREED hotspot in the run.

---

## AGREED hotspots

Ranked. Convergence of two isolated models on one passage is signal.

### H-1 — `0003:378-383` default-on / opt-out-by-omission (A/C-3 ↔ B/C-1, B/C-4)

Both passes independently selected this passage, and Pass B selected it
**twice** (as its C-1 ledger row and again as its §2 "section rewritten within
six weeks"). Pass A also names it in its §1.3 and §2. Four of the two documents'
six "most likely failure" and "will be rewritten" slots point here.

The defect is internal, not cross-document: `0003:378-383` and `0003:554` state
opposite outcomes for the same input. That is closable inside this document
today and does not need cluster reconcile. The cluster-reconcile question (does
RDR 0006 mean opt-in?) is a *second*, separable issue that both passes also
flag — B/C-4 wants lock blocked on it, A treats it as a scheduled rewrite.
Recommend: fix the internal contradiction now, book the 0006 reading question
with A8.

### H-2 — External-blocker sequencing with no phase gate (A/C-12 ↔ B/C-10, B/C-7)

Both passes read §Prerequisites and reached the same structural conclusion: the
Implementation Plan's Phase 1-4 sequence is written as linear while three of
five prerequisites are unchecked and blocked outside this document, and **no
phase references any of A7/A8/A9 or the RDR 0007 reshape as a gate**. Verified:
`0003:978-1000` (Phases 1-4) mentions A9 once, obliquely, in Phase 3 prose; A7
and A8 appear in no phase.

Pass B adds the sharper mechanism (`critique-modelB.md:66`): implementers build
against the *shipped* `Row.Guard string` because it compiles, then must redo the
work. That is a concrete, checkable prediction Pass A did not make. Pass A adds
the sharper framing: there is no stated fallback if RDR 0002 declines the
producer request. Both are worth carrying.

### H-3 — Set-literal / `contains` under-specification (A/C-14 ↔ B/C-8)

See CON-2. Agreement on the passage, disagreement on scope, resolved in
Pass A's favor with a booking correction.

### H-4 — A8 / narrowing not binding on the document that emits the finding

Pass B raises this in three places (C-3, C-7, §1 Failure 1, AT-1) and treats it
as its top finding. Pass A does **not** contest it — A explicitly categorizes A8
among "three genuinely open questions, all correctly booked"
(`critique.md:612`) — but A also declines to rank it, arguing the binding
constraints lie elsewhere. This is agreement-by-non-contest rather than
convergence, so it is ranked below H-1..H-3, but B's AT-1 (make the gate fail
rather than note A8 Pending) is the most actionable single remedy either pass
produced for it.

---

## SOLO findings — considered-and-declined vs. never-reached

Recorded because the task asks whether the silent pass *looked*.

**Pass B never reached (no trace in B's prose):** R-01, R-02, R-04, R-05,
R-08, R-10, R-13, R-18-adjacent grouping. B's document contains no mention of
RDR 0002's domain fields, of row-group definitions, of `reachable predecessor`,
of escape rows, or of RDR 0007's atom field names. These are not declined —
they are unvisited. Notably B *did* open `0002:217-218` (B cites it verbatim at
`critique-modelB.md:32`) and, exactly as Pass A predicted at `critique.md:598-603`,
used it to establish the missing **optionality** subfield for A7 while not
registering that the same line shows no domain field at all. Pass A's
prediction about the shape of the miss is confirmed by Pass B's own text.

**Pass B partially considered R-06:** B cites the spike's limits
(`critique-modelB.md:40-42`, "the Resolve spike fixture declares no set-valued
tag") but treats the spike as valid for the four operators it did cover. B
accepted A1 as `Verified`. So B considered the spike's *coverage* and never
questioned its *method*. Not a decline.

**Pass A considered-and-declined R-15, R-16:** A explicitly books A7/A8/A9 as
"three genuinely open questions, all correctly booked" and argues they are not
the binding constraints. A's decline is reasoned, and B's counter (they are
booked but not *gated*) is a different point A did not address. Both stand.

**Pass A never reached R-17, R-18, R-19.**

**R-17 deserves promotion.** Pass B's C-5 is the single most underrated finding
in either document: `0003:592` states two equalities — atom identity
`(RuleID, SourceLocator, key, block, operator, literal)` and semantic equality
`(tag, operator, literal)` — and no clause anywhere says which one overlap
detection uses. Verified: no normative contract in `0003:460-530` references
either, and MVV Scenario 5 ("same atom under reordering") does not disambiguate.
Since overlap detection is a core deliverable, this is a **BLOCK**, and it is
closable inside this document today. Pass A read this exact passage twice (for
R-07 and R-13) and did not ask the question.

---

## Verification results

Each claim checked directly against the source file named. Verdicts are this
pass's own; neither critique's characterization was accepted on trust.

### V-1 — Does RDR 0002's tag declaration carry no domain/enum/min-max field? (A/C-1)

**VERIFIED.**

`docs/rdr/0002-transition-table-as-reviewable-data.md:217-218`, read in full:

> 2. Tag declarations: tag name, provenance (`owned`, `observed`,
>    `recognized`), value kind, and optional accessor reference for observed or
>    owned read-back.

That is the complete declaration. A grep over all of RDR 0002 for
`domain|min|max|finite|value kind|enum|optionality|single-valued` returns 20
lines; **not one** is a domain, enum-value-list, or min/max field. The hits are
`enumerating`/`enumeration` (`:65`, `:103`) in unrelated prose about Cartesian
expansion, `at minimum` (`:346`) in an error-list clause, `deterministic`
(many), and `:218` itself. The Normative Contracts block (`0002:265-360`) was
read in full: it requires provenance declaration (`0002:334-336`), source rule
id and locator retention (`0002:312-315`), escape-row normalization
(`0002:294-298`), and deterministic dump ordering. **No clause obliges an
author to declare a domain.**

Cross-check confirming the orphan: `0006:251-252` lists as a lint input
"declared tags with provenance, value kind, **finite-domain metadata when
exhaustiveness is claimed**, and single-valued grouping when applicable." Three
documents (0003 A2, 0006 input contract, 0003's fixture) consume finite-domain
metadata; none produces it. `single-valued` likewise appears 0× in RDR 0002 —
the identical orphan shape, confirming Pass A's "systematic pattern" claim.

Confirming the circularity: `evidence/spikes/guard-fixture.toml` carries
`domain = ["Draft", "Final", "Implemented"]`, `min = 0`, `max = 3`, and
`domain = [true, false]`. None of those keys exists in RDR 0002's schema. The
fixture invented them.

**One scoping correction to Pass A.** A's ledger says A7 and A9 "book two
specific missing subfields … while the base domain declaration they are
subfields of is itself absent and unbooked." Verified as to the base
declaration. But note `0003:645`'s Capability Dependencies row for A7 already
states the full field list — "RDR 0002's tag declaration carries name,
provenance, value kind, and accessor reference only" — so the draft *has* the
correct fact written down; it fails to draw the consequence for A2. The defect
is a missed inference, not a missing fact. This slightly weakens A's rhetorical
frame ("nobody has looked") without weakening the finding.

### V-2 — Does RDR 0002 flatten `all`/`unless` into one predicate set, erasing `block`? (A/C-7)

**PARTIAL.**

`0002:305-309`, verbatim:

> Guard predicates MUST be represented as positive `all` predicates and
> negative `unless` predicates. Normalization MUST combine both into one
> candidate-row predicate set before ambiguity checks.

Pass A reads this as contracting the block away. That over-reads it in one
direction and under-reads it in another:

- **Over-read:** "one candidate-row predicate set" constrains the *container*,
  not the per-atom fields. It does not say atoms lose a discriminator. RDR 0003
  itself (`0003:371-373`) says "Normalization **may** combine them into one
  internal constraint object, but authoring keeps the separation" — the two
  documents are compatible on their face. And RDR 0007's SEAM clause
  (`0007:1249-1253`) independently *requires* `block ∈ {all, unless}` on every
  parsed atom, which a `Final` document imposes on the normalizer regardless of
  what 0002 says.
- **Under-read:** RDR 0002 nonetheless **does not require** the discriminator
  to survive, and RDR 0002's own normalized-row contract (`0002:312-315`)
  enumerates only "source rule id and source locator" as required retentions.
  So the block's survival is guaranteed by RDR 0007, not by the producer RDR
  0003 names.

**Verdict: PARTIAL — not a contradiction, a producer-contract gap.** Pass A's
predicted symptom (`unless` becoming per-atom negation) is not licensed by
`0002:305-309` and is contradicted by `0007:1249-1266`. But A's underlying point
stands in weaker form: RDR 0003's identity tuple names a field its named
producer does not promise. The remedy is a one-line producer request to RDR
0002 (or a citation to RDR 0007's SEAM clause as the guarantor), not a rewrite.

Downgrading A's framing here matters: the fix half should not chase a
contradiction that does not exist. It should book the retention request.

### V-3 — Does RDR 0007's atom shape have Key/Operator/Literal/Block with no `source identity`? (A/C-13)

**VERIFIED.**

`docs/rdr/0007-guard-predicate-totality.md:1249-1253`:

> SEAM. A candidate row carries its guard as a slice of parsed atoms — key,
> operator token, literal, block ∈ {all, unless} — the shape JDR 0001 §D1
> fixes; this RDR cites it and does not restate the grammar.

And `0007:1260-1265`:

> (1) The atom's four fields are named `Key`, `Operator`, `Literal`, `Block` —
> the payload's `UndecidedAtom` mirrors them by name, so divergent spellings
> would make the mirror a mapping. (2) `Block` is an exported named STRING type
> carrying exactly two constants, `BlockAll = "all"` and
> `BlockUnless = "unless"`.

RDR 0003 `:354-356`, verbatim:

> A predicate atom has four conceptual fields: tag name, operator, expected
> value, and source identity.

Both say "four fields." The sets differ: `{Key, Operator, Literal, Block}` vs.
`{tag name, operator, expected value, source identity}`. **`Block` is absent
from 0003's list and `source identity` is absent from 0007's.** And `0003:359-360`
claims "this RDR cites it rather than restating it" — immediately after
restating it incorrectly.

The internal contradiction A alleges is also confirmed: `0003:584-586` names
`block` as a member of the identity tuple, 230 lines after the field list omits
it. Two clauses in one document disagree about whether the atom carries a block.

This is the highest-consequence VERIFIED finding in either critique that is
**closable inside this document today** — it is a field-list correction, not a
peer negotiation. Rank it with R-01.

### V-4 — Is `check.sh` an awk regex that evaluates no predicate? (A/C-6)

**VERIFIED.**

`evidence/spikes/check.sh` read in full. It is a POSIX `sh` script containing
exactly three `awk` invocations:

1. Pass one collects operator tokens from `[rule.match|guard.all|guard.unless.*]`
   sections and tests each against the literal regex
   `/^(eq|in|lt|lte|gt|gte|exists|contains)$/`, printing any that fail.
2. Pass two re-scans and prints the *set* of seen operators, ordered by a
   hardcoded `split("eq in lt gte exists", order, " ")`.
3. Pass three counts `[[rule]]` headers.

It then `printf`s a hardcoded coverage string:
`coverage=status/profile routing, cap-3 handling, prelock lens sets, cluster
eligibility, rewind legality`. **That line is a literal in the script, not a
computed result.**

`output.txt` in full is four lines:

```
fixture=guard-fixture.toml
operators=eq,in,lt,gte,exists
rules=4
coverage=status/profile routing, cap-3 handling, prelock lens sets, cluster eligibility, rewind legality
```

The script never reads a tag value, never evaluates an atom, never builds a
product, never computes coverage or overlap, and never type-checks a literal
against a declared `kind`. A1's claim is that target flows "fit a closed typed
predicate vocabulary" — the script can observe neither *typing* nor *fit*.

**Two aggravations Pass A did not name**, found in verification:

- The script's step-2 ordering array is hardcoded to `eq in lt gte exists` —
  the exact five operators A1's Evidence reports. `lte`, `gt`, and `contains`
  are accepted by the validation regex but **cannot appear in the output**,
  because the printing loop iterates only the hardcoded five. The transcript
  could not have reported a sixth operator even if the fixture used one.
- The `coverage=` line — the part of `output.txt` that appears to substantiate
  A1's "four representative guard rows covering status/profile routing, cap-3
  handling, prelock lens sets, cluster eligibility, and rewind legality" — is a
  constant `printf` in the script. A1's Evidence quotes back to the RDR a string
  the RDR's own spike author typed into the script.

Pass A's characterization ("runs a spelling checker against a text file") is
accurate and, if anything, generous. **VERIFIED.**

### V-5 — Does `graph-unprovable-coverage` have no atom-level field? (A/C-9)

**VERIFIED.**

`0006:302` lists the code in the mandatory blocking table:

> | Required finite-domain proof unavailable | `graph-unprovable-coverage` |

`0006:352-356`, the finding payload contract, verbatim:

> Every blocking finding MUST carry a stable code, model identity, severity,
> human-readable message, and the source rule/context id or source span when
> the normalized model can provide one.

No atom-level field. Cross-checked `0006:288-289` (the prose restatement):
"Each finding carries a stable code, severity, model id, rule/context id when
available, source span when available, and a concise human message." Same set,
no atom.

RDR 0003 `:475-482` requires the finding to "name the participating row **and
the atom that can refuse**." That is a payload extension imposed on a `Final`
peer's shipped contract. And `0003:538` (§authority) states this RDR "reuses
0006's existing blocking code rather than minting a category" — which is
precisely the move that makes the extension invisible: it looks like reuse and
is actually an extension.

A's further claim that A8 does not cover this is confirmed: A8's text
(`0003:1128-1136`) is entirely about *which document records the narrowing*,
never about payload fields. **The payload extension is unregistered.**

**Severity note.** Pass A ranks this alongside its BLOCK findings. I rank it
**DEFER**: it is real, it is unbooked, and it is a one-line addition to A8's
scope or a new A-record. It does not prevent the design from being coherent, and
it closes at the same cluster-reconcile venue A8 already names.

### V-6 — Does RDR 0006 invariant 4 count "a declared escape row" as coverage, and does RDR 0003 never mention escape rows? (A/C-8)

**VERIFIED, both halves.**

`0006:277-279`, verbatim:

> 4. **Guard exhaustiveness / gap** — for each state/outcome pair that claims
>    closed coverage, finite-domain input assignments must either match exactly
>    one modeled row **or a declared escape row**.

RDR 0003's coverage identity (A2 Evidence, `0003:135`) is
`union(row_i accepted assignments) == scoped product` — **no escape-row term.**

`grep escape` over RDR 0003 returns exactly four lines: `:80`, `:112`, `:114`,
`:689`. Read in context, all four are about *human* escape hatches in prose
("an explicit escape remains necessary where judgment is policy", "authors will
request escape hatches"). **None is about escape rows as a normalized row
class.** The ten-row `disposition` table (`0003:545-557`) enumerates every input
class — qualifying rows, pruned rows, disabled rows, absent keys, zero matches,
multiple matches, undeclared domains, oversized products, parse failures,
withheld claims — and has no escape row.

Confirmed at the producer end: `0002:294-298` requires normalization to "render
an escape rule as a candidate row with row kind `escape`, **its normal predicate
set**, source rule id, source locator, and modeled failure class list." So
escape rows arrive at lint carrying guards, indistinguishable from ordinary rows
except by `row kind` — a discriminator RDR 0003 never reads.

Also confirmed: `0002:324-331` makes the resolver consult escape rows **only**
after ordinary matching fails ("If zero or multiple non-escape rows match, the
resolver MAY return a modeled escape disposition"). So Pass A's premortem
consequence is exactly right — an escape row can supply lint coverage for a
combination the resolver will still refuse first. **Two documents define
coverage over different row sets, and the difference is load-bearing.**

This is the finding most likely to produce a shipped-and-wrong system, and it
is a genuine SOLO-A: Pass B's document contains no mention of escape rows.

### V-7 — Supplementary checks on A's remaining specific claims

Checked because they carry weight in A's ranking.

**"`reachable predecessor` is never defined" (A/C-4).** **VERIFIED with a count
correction.** Grep across 0002/0003/0006 returns: RDR 0003 `:192`, `:497`,
`:517` (three, not A's claimed four — A's `:417` is `reachable predecessor
write` at `0003:417`, so the count is four if that phrasing counts, three for
the exact hyphenated form); RDR 0006 `:130` (one, not A's claimed five). RDR
0002: **zero**. No document defines the predecessor relation, its carrier, or
its computation. The count is off; **the finding is sound and the direction of
the error makes A's case *stronger*, not weaker** — RDR 0006 uses the term even
less than A claimed, and its single use (`0006:130`) is inside A3's Evidence,
citing RDR 0003's Normative Contracts. Meanwhile RDR 0003 `:517` states the MUST
using the undefined term. The citation loop is real: `0006:130` → 0003's
contracts; `0003:165-167` (A6) → RDR 0006's Technical Design.

**"RDR 0006 has no A6" (A/C-4).** **VERIFIED.** RDR 0006 carries exactly A1-A5
at `:98`, `:110`, `:122`, `:137`, `:150`. There is no A6. The 3amigo
disposition's warrant citing "RDR 0006 via A6" points at a record that does not
exist. Note this is a defect in the *3amigo evidence record*, not in the RDR
draft itself — the fix half should not go looking for a bad citation in
`0003.md`.

**"`row group` appears 0× in RDR 0006" (A/C-2).** **VERIFIED.** `grep -c "row
group"` over RDR 0006 returns **0**. `scoped` appears twice: `:115` (inside A2's
Evidence, quoting RDR 0003's own derivation) and `:501` ("properly scoped to
design-time graph acceptance" — unrelated). `selection context`: zero. So RDR
0003 `:374-377` delegates grouping to a document that contains no grouping
definition, and RDR 0006 A2 (`0006:110-119`) delegates the coverage derivation
back to RDR 0003 A2. **The cycle is real and closed.**

Additional confirmation Pass A asserted and I checked: RDR 0006's own Normative
Contracts (`0006:340-370`) never mention state/outcome pairs, row groups, or
partitioning. Invariant 4 (`0006:277-279`) uses "state/outcome pair" but is a
numbered invariant, not a definition, and RDR 0006's input contract
(`0006:248-256`) never says how a row induces a state.

**"RDR 0006's exhaustiveness clause covers only the non-finite case" (B/C-3,
B §1).** **VERIFIED.** `0006:344-350`:

> Graph lint MAY claim exhaustiveness only over finite declared domains
> supplied by the predicate/tag model. If a required dimension is not finite,
> lint MUST emit a blocking inability-to-prove finding for any contract that
> depends on closed coverage.

Scoped to "if a required dimension is not finite." Says nothing about a
fully-finite product whose participating row can still refuse. Pass B's Failure 1
is grounded. (Note the same clause is the source of Pass A's R-01 orphan: it says
domains are "supplied by the predicate/tag model" — a supplier that, per V-1,
declares none.)

**A/C-1's claim that RDR 0006 A2 is `Verified` on nonexistent domains.**
**VERIFIED.** `0006:110-119` is `Status: Verified`, `Method: Peer RDR`, Evidence
citing "declared set-universe" domains among others. A9 of RDR 0003 concedes
those do not exist. So a `Final` RDR carries a `Verified` assumption resting on
a `Draft` RDR's undeclared capability. Locking 0003 does not fix this; it
propagates it. This belongs in the cluster-reconcile packet regardless of what
happens to 0003.

### Verdict summary for the six claims the task named

| Pass A claim | Verdict |
|---|---|
| C-1 — RDR 0002 declares no domain field | **VERIFIED** (scoping note: draft states the fact at `0003:645`, fails to draw the consequence) |
| C-5 — A6 `Verified` on an adjacent fact | **VERIFIED** |
| C-6 — `check.sh` evaluates no predicate | **VERIFIED** (plus two aggravations A missed) |
| C-7 — RDR 0002 erases `block` | **PARTIAL** — gap, not contradiction; RDR 0007 SEAM independently guarantees `Block` |
| C-13 — 0003's four atom fields ≠ 0007's four | **VERIFIED** |
| C-14 — RDR 0007 files an unbooked set-literal duty on 0003, naming `in` | **VERIFIED** (quotation loose, substance verbatim) |

**Zero REFUTED.** One PARTIAL (C-7). Pass A's harsh line held up under
independent check.

---

## Ranking for the fix half

### BLOCK — must close before lock

Ordered by (a) whether closable inside this document, then (b) consequence.

**Closable here, today — no peer negotiation needed:**

1. **R-13** (V-3) — atom field list says `source identity` where the kernel seam
   says `Block`, and contradicts this RDR's own identity tuple 230 lines later.
   One-line fix. Ships the inverse of every `unless` guard if missed.
2. **R-03 / H-1** — `0003:378-383` and `0003:554` state opposite outcomes for an
   undeclared dimension. Pick one. Both models found this.
3. **R-02** (V-7) — define the row group here, or state that no document defines
   it and book it. The cycle with RDR 0006 is closed and neither end will break
   on its own.
4. **R-17** (B/C-5) — say which equality overlap detection uses. One clause.
5. **R-08** (V-6) — state whether escape rows join their group's product, form
   their own group, or are excluded. RDR 0006 counts them; RDR 0003 is silent.
6. **R-04** (V-7) — either define `reachable predecessor` or downgrade the two
   MUSTs (`0003:497-501`, `0003:517-519`) that decide on it.

**Requires a peer or a status change:**

7. **R-01** (V-1) — the finite-domain producer request. This is the single
   largest item. Merge it with A7 and A9 into one RDR 0002 tag-declaration
   request; A2 reverts to `Pending` until RDR 0002 accepts.
8. **R-05** (V-7), **R-06** (V-4) — A6 and A1 revert to `Pending`. A6's evidence
   proves an adjacent fact; A1's spike is a token scanner. Both are `Verified`
   stamps the draft's own Method vocabulary (`0003:286-291`) forbids.
9. **R-14 / CON-2** — book the canonical set-literal byte spelling as a new
   assumption owned by *this* RDR. A `Final` peer's determinism claim depends
   on it.
10. **R-11** — the MVV cannot be authored today (needs R-01, A7, A9). Either
    descope §Scope Verification's "in scope, not deferred" claim or block on the
    producer request.

### DEFER — real, closable after lock or at cluster reconcile

- **R-09** (V-5) — register the `graph-unprovable-coverage` atom-field extension
  as a peer obligation. Travels with A8.
- **R-16** (B/C-3, B/C-7) — A8's narrowing has no phase gate. Adopt B's AT-1:
  make the Finalization Gate fail rather than note A8 `Pending`.
- **R-15** (B/C-2) — promote the A7 projection derivation out of the side
  evidence file into a normative clause, and name `exists` scope in Phase 1.
- **R-12 / H-2** — sequence the Phase list against the four external blockers,
  including B's specific `Row.Guard string` interim-grammar risk.
- **R-18** (B/C-6) — state a scoped-product size bound or complexity class.
  Becomes live the moment R-01 closes (see CON-1).
- **R-10** — the §Assumption Verification section audits label hygiene and
  reports it as verification. Meta-defect; the remedy is procedural, and it is
  already discharged in practice by this iteration's findings.
- **R-07** (V-2) — book the block-retention request to RDR 0002, or cite RDR
  0007's SEAM clause as the guarantor. Do **not** treat as a contradiction.

### NOISE — correct observation, no action

- **R-19** (B/C-9) — cross-RDR failure-class rendering. Correct that the
  distinction spans three documents, but A4 (`0003:143-155`) explicitly assigns
  the CLI envelope mapping to RDR 0005. Out of scope by design, not by omission.

---

## What the dual-model run bought

Neither pass alone would have produced this list.

Pass A found every cross-document factual defect — the missing producer, the
delegation cycle, the atom field-list mismatch, the escape-row divergence, the
unbooked `in` duty. All five verify. Pass A's method was to open every peer
document and read the cited line. Its cost: it accepted A7/A8/A9 as "correctly
booked" and never asked whether they were *gated*.

Pass B found every process and sequencing defect — the ungated lock blockers,
the unmeasured threshold, the identity/semantic-equality ambiguity, the
interim-grammar risk. Four of those five are absent from Pass A. Its cost: it
read `0002:217-218` and, like the two lenses before it, extracted only the
subfield it was looking for.

The one place they overlap hardest (H-1, `0003:378-383`) is the one place they
predicted **opposite user symptoms** from the same passage — which is not a
contest but the strongest possible confirmation the passage is broken.
