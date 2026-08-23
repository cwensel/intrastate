Model: claude-opus-5[1m]

# Repeatability-Lite DIFF — RDR 0003 (Guard Predicate Exhaustiveness), iteration 3

One alternate-model reconstruction (`run-1.md`, Claude Sonnet 5,
`variant: lite (profile: large)`) diffed against the current draft. This is a
lite pass: no same-model consistency sample, no three-run disagreement count.
Findings below are contract silences only — places where a specific passage in
the *current* draft (or a specific silence in it) left an algorithmic contract
under-determined and the run had to invent a resolution. Naming, wording,
package path, Go identifiers, helper decomposition, and representation choices
the RDR deliberately leaves open are not findings; they are listed in §3 under
Inadmissible.

Iter-2's D-1..D-7 were all resolved into the live text (see §4). Nothing below
is a restatement of them; each finding is grounded against a passage that exists
in the draft as it stands today, verified by line number.

## 1. Findings

| ID | RDR anchor | The silence | How run-1 resolved it | Why it is a determinacy gap |
| --- | --- | --- | --- | --- |
| D-1 | `Normative Contracts`, the bound clause's per-kind assignment-count table (`:1212-1218`) vs. the optionality-default clause (`:846-850`) and the single-valued clause (`:893-901`); `trace` step 3 (`:1313`) | **The assignment-count table's inputs are never joined to the declaration defaults.** The table gives `\|domain\|` for a *single-valued* `enum`/`int`, `2^\|domain\|` for a finite kind **without** the marker, and "×2 for any **optional** key". The optionality clause fixes *absent marker ⇒ optional*; the single-valued clause states no default, so *absent marker ⇒ not single-valued* by construction of the table's own two rows. The RDR never states that these defaults feed the table, and its own worked example does not apply them: step 3 gives `profile` = 4, `prelock_iterations` = 4, `cluster_eligible` = 2 with a presence dimension for `cluster_eligible` alone — the single-valued counts, with `profile` and `prelock_iterations` getting no ×2 — over a fixture (`evidence/spikes/guard-fixture.toml`) whose tags declare neither `optional` nor `single_valued`. | Applies the table uniformly against the declared defaults: helper 2 (`:219-226`) computes "`\|domain\|` for single-valued enum/int, `2` for single-valued bool, `2^\|domain\|` for unmarked finite kinds, `2^\|universe\|` for `set`, and ×2 per optional key's presence dimension"; pseudo-code step 2 (`:308-309`) adds a presence dimension for *every* optional participating key with no `exists`-atom condition. | Changes the product cardinality by orders of magnitude on the RDR's own fixture — 64 under the trace, 8,192 under the run's uniform reading (`2^4·2 × 2^4·2 × 2^2·2`) — therefore changes the coverage verdict, the gap witness, and whether the too-large bound fires. Two conforming implementations return different verdicts on the same model, which is exactly what the published-bound clause (`:1198-1200`) promises cannot happen. |
| D-2 | `Normative Contracts`, the `unless`-Kleene clause (`:1091-1105`); the accepted-assignments rule (`:769-772`); the report-every-defect clause (`:1137-1147`) | **What a can-refuse row contributes to the coverage union and the overlap intersections is unstated.** `:769-772` defines a row's accepted assignments as `all`-intersection minus the `unless` block, unconditionally. `:1092` then says the subtraction is valid "**only when its atoms are decided**" — so for a row whose `unless` atom is over an optional key the subtraction does not apply, and the RDR stops there. It fixes that the group's *claim* is withheld; it never says whether that row still contributes `allSet` unsubtracted, contributes nothing, or contributes something else. The report-every-defect clause makes this load-bearing rather than moot: it requires "any coverage gap over a provable product is its own finding" **alongside** the withholding, so the union is still computed for a withheld group. | Contributes **nothing**. Pseudo-code step 5 (`:333-341`) assigns `accepted[row]` only in the all-decided branch; the else-branch marks the row "can refuse" and leaves the entry unset. Step 6 then unions over rows including that one, so its assignments silently drop out of `covered`. | Determines whether a withheld group *also* emits a `graph-coverage-gap`. Under the run's reading every withheld group with a refusing ordinary row emits a spurious gap finding naming a witness assignment that a row does in fact accept; under the unsubtracted reading it does not. Both conform to the text. The gap finding is mandatory-blocking and carries a concrete witness assignment (`:1237-1242`), so this is a difference in emitted blocking diagnostics on the same model, not a presentation choice. |
| D-3 | `Normative Contracts`, the finite-domain clauses (`:976-980`, `:971-974`) and the bound clause's cardinality table (`:1206-1218`); `disposition` rows `:1294` and `:1295`; report-every-defect (`:1137-1147`) | **The cardinality of a product containing a non-finite dimension is undefined, and the RDR requires it to be reported anyway.** The per-kind table has rows for `enum`/`bool`/`int`/`set` only — no row for `scalar`, and no row for a finite kind that declares no domain. Yet `:1294` and `:1295` are separate input classes that the report-every-defect clause forbids short-circuiting between: a group can carry both an undeclared dimension and an over-large product, and `:1202-1204` requires the over-large refusal to "report the product size it computed". The RDR never states whether the too-large check runs at all once a dimension is unprovable, nor what assignment count an unprovable dimension contributes to the arithmetic. | Runs both unconditionally and in sequence: pseudo-code step 3 (`:311-316`) collects unprovable dimensions, then step 4 (`:318-321`) computes `cardinality = product(assignmentCount(decl) for key, decl in dims)` over **all** dims — including the ones just found to have no finite domain, for which `assignmentCount` has no defined value — and compares it to the bound. | Fixes whether an author with one `scalar` dimension sees one finding or two, and what number the second one reports. `:1203-1204` makes the reported figure normative ("so an author can tell an over-large product from an undeclared dimension") — the run's resolution reports a figure computed from an undefined term, which cannot be model-independent. Cross-implementation: `:1229-1230` requires "each MUST return the same verdict for the same model". |
| D-4 | `Normative Contracts`, the escape-population clause (`:1052-1063`): "overlap **among escape rows for the same failure class**"; `disposition` row `:1299`; A17 (`:541-565`) | **The RDR never says how an escape row's failure class is read, nor how rows declaring more than one class are partitioned.** RDR 0002 does own the field — an escape row's `escape` list contains the failure classes it models (`0002:328-334`, `0002:738`) — but that list is plural, and this RDR's clause quantifies over "the same failure class" (singular) without citing the field or stating the partition. Two escape rows sharing only *one* of two declared classes, or one row declaring both, have no stated population membership. | Left as an undefined predicate and then flattened: pseudo-code step 7 (`:355-357`) iterates `pairs(escapeRowsForSameFailureClass)` as a single flat pairing rather than one population per class, and the run's data model (§1, §3) carries no failure-class field on any type — the concept never reaches the reconstructed contract. | Decides whether two escape rows overlapping on one shared class of two draw a finding. `:1058-1061` grounds the check on RDR 0002 admitting a rescue "only when *exactly one* escape row matches" — that is per class, so a flat pairing over-reports and a naive per-row pairing under-reports. Cross-RDR: RDR 0002 owns the field, RDR 0006 emits the finding, and A17 is the open venue where the two-population reading is being confirmed — it will be confirmed without this term. |
| D-5 | `Normative Contracts`, the narrowing clause (`:1084-1089`) — "when a **participating row** can refuse" — vs. the escape clause (`:1045-1047`) — "a group whose **ordinary** row can refuse" — and the can-refuse clause (`:1162-1169`), which says "a row" | **Whether an escape row can trigger the withholding is unstated.** The three clauses use three different quantifiers over the same predicate. An escape row carrying a value atom over an optional key satisfies the can-refuse test at `:1163-1165` verbatim; the escape clause's phrasing implies the trigger is the ordinary population; the narrowing clause says "participating", which the participation clause (`:987-988`) defines as "any row in that group". | Loops over all rows: pseudo-code step 8 (`:362-364`) is `for row in group.rows` with no population filter, so a guarded escape row withholds the whole group's claim. | The escape row is consulted only after ordinary resolution fails, so a guarded escape row that could refuse is a different runtime situation from an ordinary one that could — which is the exact distinction the two-population overlap clause was written to preserve. Under the run's reading a group with a fully-provable ordinary population is blocked by its escape row; under the ordinary-only reading it is green. Load-bearing on the flagship false-green defence (MVV Scenario 8, `:2049-2057`, which adds a bare escape row but not a guarded one). |

## 2. Detail — load-bearing and cross-RDR

**D-1 is the highest-consequence finding and it is self-checkable inside the
document.** The critique iteration-3 pass added the per-kind assignment-count
table precisely because "`set` cardinality [was] understated by an exponent"
(Status line, `:26-27`), and the table now carries three rows whose trigger is a
*declaration marker*: single-valued present, single-valued absent, optional. Two
other clauses fix what an absent marker means — `:846-848` makes an absent
optionality marker declare the key optional ("the conservative default"), and the
table's own two-row split makes an absent single-valued marker mean the powerset
count. Nothing joins them. The desk trace, which is the document's only worked
arithmetic, does not apply either default: it computes the fixture's Draft group
at 64 assignments using single-valued counts for `profile` and
`prelock_iterations` and a presence dimension for `cluster_eligible` alone,
against a fixture (`evidence/spikes/guard-fixture.toml:5-40`) in which no tag
declares `optional` or `single_valued` at all. Applying the table as written
gives 8,192. Run-1 applied the table as written.

The cross-RDR reach is direct. The product cardinality is the quantity the
published bound is measured against (`:1206-1208`), and the bound clause promises
that "the same model MUST receive the same verdict on every conforming
implementation" (`:1198-1200`) — a promise that cannot hold while two readings of
the same declaration produce cardinalities two orders of magnitude apart. It also
lands on RDR 0006: the coverage verdict and its witness assignment
(`:1237-1242`) are RDR 0006's blocking finding payload, and on A15, whose MVV
Scenario 3 constructs equal-cardinality product pairs relative to the published
bound — a construction that presupposes cardinality is computable from the
declarations without a further convention. And it lands on A16: the
single-valued marker's authoring location is still open at RDR 0002, so every
fixture authored before A16 lands declares no marker, which is exactly the input
class the table's `2^|domain|` row governs and the trace ignores.

**D-2 is the second load-bearing one and it is new to iteration 3.** The
`unless`-Kleene clause is a critique iteration-3 addition (Status line,
`:27-28`); it correctly closes the false-green where an undecided `unless` block
was subtracted as if decided, and it correctly routes the row to the withholding.
What it does not do is finish the arithmetic it interrupted. `:769-772` states a
row's accepted assignments as one formula with no undecided branch; `:1092`
disables that formula's second term for a can-refuse row and supplies no
replacement. The report-every-defect clause is what makes the missing term
observable rather than academic: it requires a coverage gap "over a provable
product" to be emitted as its own finding alongside the withheld claim, so the
union is genuinely computed for the very group whose row-contribution is
undefined. Run-1's "contributes nothing" reading makes every such group emit a
spurious `graph-coverage-gap` naming a witness assignment the withheld row
actually accepts — a blocking finding on a model that has no gap. The
alternative reading (contribute `allSet` unsubtracted) emits the withholding
alone. Both are conforming; the RDR selects neither.

**D-3 exposes an ordering the report-every-defect clause created and did not
close.** Before that clause, an unprovable dimension could reasonably be read as
terminating the proof (that was iter-2's D-2, now fixed). Now it explicitly must
not: `:1140-1142` requires "each unprovable dimension, each refusing row, each
overlapping pair, and any coverage gap over a provable product" as separate
findings, and `:1146-1147` allows a withheld claim and a coverage finding
together. That forces the too-large check and the unprovable-dimension check to
coexist on one group — and the cardinality the too-large check reports is then
computed over a dimension set containing a term the per-kind table does not
define. The table's own justification paragraph (`:1220-1225`) argues the
arithmetic is what makes the "a bound discovered by exhausting memory is not a
conforming bound" sentence enforceable; that argument only reaches dimensions the
table covers. A `scalar` dimension, or an `enum` with no declared domain, is
outside it. Run-1 multiplied through anyway.

**D-4 and D-5 are both escape-population silences and both sit on A17's open
agenda**, which is why they are worth surfacing now rather than after lock. A17
carries the cross-document confirmation that overlap is checked in two
populations, and its done-condition (`:1821-1827`) is that "RDR 0006's invariant
3 repair matches this RDR's two-population clause". If the clause reaching RDR
0006's refine does not state how failure classes partition the escape population
(D-4) and does not say whether an escape row participates in the withholding
(D-5), the confirmation will ratify an underspecified clause and the divergence
will surface at implementation instead. D-4 in particular is a one-line fix that
costs nothing: RDR 0002 already carries the `escape` failure-class list
(`0002:330-332`, `0002:738`), so this RDR need only cite it and state that the
escape population is partitioned per declared class, with a row declaring two
classes appearing in both.

## 3. GUESS clusters

The run carries 30 GUESS markers and is explicit about where they concentrate:
"Almost entirely in Section 1 (concrete Go API)" (`run-1.md:377`). Most land on
implementation detail the RDR deliberately leaves open. Grouped by the contract
they touch.

### Admissible — GUESSes landing on an algorithmic or normative contract

- **`ScopedProduct.Cardinality uint64` with the comment "must handle up to `2^N`
  for set kinds"** (`run-1.md:155`) → **D-1** and **D-3**. The run recognised
  that the table's exponential rows make the cardinality unbounded in practice
  and had to pick a width, without a rule for what an undeclared or `scalar`
  dimension contributes to the same arithmetic.
- **Pseudo-code step 2's unconditional presence dimension** (`:308-309`) → **D-1**.
  The RDR's presence-dimension clause (`:1016-1027`) is headed "An `exists` atom
  projects…", while the bound table's row is "any optional key ×2" with no atom
  condition; the run applied the declaration-triggered reading. The two readings
  differ on any group with an optional guard key carrying only value atoms.
- **Pseudo-code step 5's else-branch, which marks a row can-refuse and leaves
  `accepted[row]` unassigned** (`:338-341`) → **D-2**. Not marked GUESS by the
  run, which is itself the signal: it read the silence as settled.
- **`assignmentCount(decl)` applied to every dim after the unprovable set was
  collected** (`:319`) → **D-3**. The run kept the two checks in sequence with no
  guard between them.
- **`pairs(escapeRowsForSameFailureClass)` as a bare undefined predicate**
  (`:355`) → **D-4**. The run named the concept and then carried no failure-class
  field anywhere in its data model (§1, §3), so the term is unreconstructable
  from the RDR.
- **Step 8's unfiltered `for row in group.rows`** (`:362`) → **D-5**.

### Inadmissible — impl detail the RDR deliberately leaves open

Listed only to show they were considered and dismissed; none is a finding.

- Package path (`package guard`, "could equally be `internal/guard` or
  `internal/predicate`", `run-1.md:23-25`). The RDR says it "owns the evaluator
  semantics, not command I/O" (`:1448`) and delegates lint authority to RDR 0006;
  layout is free.
- All Go identifier choices — `Kind`/`KindEnum`, `Operator`/`OpEq`, `Block`/
  `BlockAll`, `TagDecl`, `Atom`, `Verdict`, `RowGroup`. Naming is out of scope by
  rule, and the run reproduced the five kind *tokens* and the eight operator
  *tokens* exactly, which is the part the RDR does fix (`:819-822`, `:947-949`).
- `TagDecl` field shapes — `Domain []Literal` vs. `Min, Max *int64` vs.
  `Elements []Literal` (`run-1.md:45-52`). The RDR fixes that `{min..max}` is
  "notation, not wire spelling: the authored form is two fields, `min` and `max`,
  and both endpoints are **inclusive**" (`:836-838`) and the run rendered exactly
  that; the Go representation is free, and the authoring location is RDR 0002's.
- `Literal{Kind, Value any, Set []any}` (`run-1.md:97-101`). The canonical
  spelling *is* fixed now (`:1179-1188`) and the run cited it correctly
  ("canonicalize to an unordered, duplicate-free typed set before entering the
  identity tuple"); the in-memory shape is free. This closes iter-2's D-5.
- `TagDecl.Validate() error`, `ParseAtom(...)`, `Evaluator.EvaluateAtom(...)`,
  `ProveExhaustiveness(...)` signatures and their existence as entry points. The
  RDR names well-formedness rules and semantic kinds, not an API. Helper
  decomposition is free — and the run said so itself for the prover
  (`run-1.md:243-244`).
- The three named internal helpers as *functions* rather than one routine.
  Decomposition taste; the semantics they implement are pinned and were rendered
  correctly (see §4).
- `AtomIdentity.RuleID string` / `SourceLocator string` types (`run-1.md:107-108`).
  Explicitly RDR 0002's fields, cited not owned (`:1265`).
- The run's note that `ProveExhaustiveness` "may in fact live inside RDR 0006's
  lint package rather than here" (`run-1.md:165`). That is the run reading the
  authority census correctly (`:1263`): this RDR owns the *verdict*, RDR 0006
  owns the *finding*. Not a gap.
- The run's note that the atom-naming half of the withheld-claim finding "has no
  producer on either surface yet" (`run-1.md:174-175`, `:364-367`). That is A20
  read correctly, with a booked plan (`:601-617`) — not a silence.
- Restating well-formedness at lint time "since lint depends on it"
  (`run-1.md:292-298`), flagged `# GUESS ordering` at `:298`. The RDR does fix
  the phase boundary that matters — declaration errors before normalization,
  literal-outside-domain at guard parse after normalization and before resolution
  (`:936-943`) — so the residual is intra-phase diagnostic ordering for an atom
  with two simultaneous defects. Real but not self-evidently load-bearing at
  lite, and the run did not diverge from the RDR on it.

## 4. Agreement

Contracts the run rendered exactly as the draft states them. This confirms
determinacy at each of these passages.

Closed by the current text — iter-2 findings that no longer reproduce:

- **iter-2 D-1 (participating dimension).** The participation clause
  (`:987-991`) now pins the *referenced* reading, and the match/guard split is
  fixed per atom per group (`:1004-1013`). Run-1 reproduced both without
  guessing: helper 2 (`:220-223`) states participation as "any row in *that*
  group carries a guard atom over the key, in `all` or `unless`, regardless of
  whether it discriminates" and excludes match keys as "the grouping context,
  never a dimension".
- **iter-2 D-2 / D-6 (finding-set scope and completeness).** The
  report-every-defect clause (`:1137-1147`) closed the short-circuit. Run-1's
  step 9 (`:369-372`) states "All findings from steps 3-8 are emitted together,
  not short-circuited at the first one", and its `Verdict` carries plural
  `Gaps`, `Overlaps`, `Withheld`. The one residual short-circuit in the run's
  step 4 (`:325`) is annotated against itself ("but continue to steps 5-7 too —
  withholding MUST NOT suppress other findings"), i.e. the clause won.
- **iter-2 D-3 (too-large bound), unit and shape half.** Now pinned as
  cardinality of the scoped product, one published integer constant reported
  beside the computed cardinality (`:1206-1230`). Run-1 threaded `bound uint64`
  through `ProveExhaustiveness` and emitted
  `InabilityToProve(product=cardinality, bound=bound)` (`:321`) — both payload
  fields present, which is exactly what iter-2 recorded as missing. The
  *adequacy* half remains open as A15 by design.
- **iter-2 D-4 (coverage-gap attribution).** The gap-attribution clause
  (`:1237-1242`) now requires the selection context, every rule id in the group,
  and one concrete uncovered assignment. Run-1 emitted
  `gap finding(context, all rule ids in group, one witness assignment)` (`:349`)
  verbatim.
- **iter-2 D-5 (canonical set-literal spelling).** The set-literal clause
  (`:1179-1188`) closed A13. Run-1 cited it rather than inventing an ordered
  slice (`run-1.md:94-96`, `:267-268`).
- **iter-2 D-7 (interim atom-naming representation).** Dismissed-with-cite at
  iter-2 and correspondingly not re-raised: run-1 recorded the missing producer
  as A20 rather than minting an interim contract.

Determinate in the current text — rendered without invention:

- **Seam scope.** The evaluator decides value semantics over a present value
  only, never reads the tag view, never sees `exists`; presence, existence
  atoms, absent-key unevaluability and per-atom combination are the kernel's
  (`:958-963`). Run-1 reproduced all five clauses including the two-valued /
  three-valued split ("three-valued combination (Kleene, `¬U = U` for unless) is
  the kernel's job, not this evaluator's", `run-1.md:139-141`).
- **The five value kinds and their exact tokens** (`:819-822`), including
  `scalar` as the opaque kind that "can never bear an exhaustiveness claim"
  (`:822-826`).
- **Domain/kind agreement.** The four-column table (`:919-925`) and its three
  named errors — `{min..max}` on an `enum`, an element universe on a non-`set`,
  a single-valued marker on a `set` or `scalar` — reproduced exactly
  (`run-1.md:54-56`, `:190-192`), including the phase boundary (declaration
  loader, before normalization completes).
- **Optionality default.** Absent marker ⇒ optional, the conservative default
  (`:846-850`). Run-1: "`Optional bool // default true (no marker = optional,
  conservative default)`" (`run-1.md:50`). (Its *consequence* for the product is
  D-1; the default itself is determinate.)
- **The single-valued marker's operational content.** One dimension of
  `|domain|` with the marker, one boolean dimension per value without it, never
  on a `set` (`:893-901`). Run-1 stated both halves in helper 2.
- **Atom shape and source identity.** Four fields (`Key`, `Operator`, `Literal`,
  `Block`) cited from JDR 0001 §D1, with source identity explicitly *not* an atom
  field and carried by the enclosing row (`:727-731`). Run-1 flagged this as
  stated rather than guessed (`run-1.md:84-86`, `:270-273`), and correctly noted
  `Row.Guard string` is pending RDR 0007's reshape and "not this RDR's field to
  change".
- **Identity vs. semantic equality.** Total tuple
  `(RuleID, SourceLocator, key, block, operator, literal)` for diagnostics and
  dedup; `(tag, operator, literal)` for domain computation only (`:1343-1360`).
  Reproduced verbatim (`run-1.md:103-113`, `:275-277`).
- **Row-group definition owned here.** Same source state, same recognized
  outcome, the set RDR 0001 resolves exact-one over, not read back from RDR 0006
  (`:745-754`). Run-1: "defined normatively in this RDR, not RDR 0006"
  (`run-1.md:146-147`).
- **Refuse and downgrade are one blocking outcome**, over all three input
  classes, with no non-blocking tier (`:1150-1160`). Run-1: "**all three are the
  same single blocking outcome**, never a soft/advisory tier" (`run-1.md:198-199`).
- **The can-refuse test is syntactic over declarations, in either block**
  (`:1162-1169`). Run-1: "decided syntactically from the optional-declared field
  alone, in either block" (`run-1.md:239-240`). (Whether escape rows are in scope
  is D-5; the *procedure* is determinate.)
- **`unless` is one conjunctive exclusion block, never per-atom negation, never
  source-order priority** (`:740-744`, `:965-969`). Rendered correctly in helper
  3 and pseudo-code step 5.
- **Escape rows participate in the coverage union like any other row; a bare
  escape row denotes the whole product and closes coverage by itself; "an escape
  row exists" is never a coverage-satisfying fact outside the union; and the
  closure must be named in the verdict rather than certified as a bare green**
  (`:1031-1035`, `:1069-1081`). Run-1 reproduced all four, including the A19
  observability requirement (`run-1.md:168-175`, `:344-345`).
- **Escape/guarded overlap is not a runtime ambiguity and MUST NOT be reported**
  (`:1057-1058`). Run-1's step 7 closes with "escape-vs-guarded overlap is NOT
  checked — not a runtime ambiguity" (`run-1.md:358`).
- **No wire format and no round-trip owned here** (`:1388-1391`). Run-1 cited
  "Round-Trip / Inverse Invariants" explicitly and introduced no encode/decode
  pair (`run-1.md:279-281`).
- **The predicate semantic kinds** — unknown operator, unknown tag,
  operator/kind mismatch, literal parse failure, literal-outside-declared-domain,
  declaration/kind disagreement — and their routing through
  `internal/cli/clierr::CLIError` / `internal/cli/respond::Fail` per A4
  (`:214`, `:1296-1298`). Enumerated exactly (`run-1.md:184-205`), with
  `guard_unevaluable` correctly excluded as RDR 0007's.
- **Assumption ledger.** A16, A17, A19, A20 all cited at the right places and
  none fabricated as closed.

The determinacy of the seam split, the declaration model's five fields and their
kind-agreement table, the identity tuple, the row-group definition, the
one-blocking-outcome rule, and escape-row coverage participation is strong
evidence those passages are doing their job.

## 5. Escalation verdict

**Escalate.**

Trigger: *the lite diff surfaces a determinacy silence on a load-bearing or
cross-RDR contract*. The other trigger — suspiciously clean agreement with no
GUESS markers — does **not** fire: run-1 carries 30 GUESS markers and five
substantive silences, so there is no overfit signature.

The specific contract is **the scoped product's cardinality arithmetic**: the
per-kind assignment-count table at `:1212-1218` and its join to the two
declaration defaults it depends on (absent optionality marker ⇒ optional,
`:846-848`; absent single-valued marker ⇒ powerset count, the table's own second
row). It is load-bearing because cardinality is the quantity the published bound
gates on and the object the coverage identity is stated over, and it is cross-RDR
because that verdict is RDR 0006's blocking-finding input, because the missing
single-valued marker is A16's open request against RDR 0002, and because A15's
MVV Scenario 3 constructs its test pair from a cardinality it assumes computable.
**D-2** compounds it on the same output surface: the union whose cardinality D-1
governs has an undefined contribution from every can-refuse row.

Note the contrast with iter-2's escalation reasoning, which declined on the
grounds that its findings were textual contradictions readable directly from the
draft and that further runs would only poll model opinion. That reasoning does
not transfer here. D-1 is not a contradiction between two prose phrasings that an
author can simply pin: the table and the defaults are each individually correct,
and the document's only worked arithmetic — the desk trace at `:1313`, revised by
critique iteration 3 specifically to *compute* rather than assert — silently uses
a third reading. Which reading is intended is a live design question with a
different product, a different bound interaction, and a different A16 dependency
under each, and a second and third reconstruction would show whether an ordinary
reader lands on the table or on the trace. That is the question the full ×3 lens
answers and one run cannot.

Escalation would require rewriting `run-1.md`'s `variant:` line to
`full (escalated: <reason>)`. Not done here — this pass records the verdict and
its reasoning only.

---

## 6. Escalation disposition (resolve pass, same iteration)

**Declined — lite stands. `run-1.md`'s `variant:` line remains `lite`.**

The escalate trigger above was the load-bearing/cross-RDR silence at D-1, on the
reasoning that one reconstruction shows a silence is *reachable* but not which
reading a second reader would pick. That question is moot as of this iteration's
resolve: D-1 is pinned (the assignment-count table now reads the declaration
defaults normatively), and the deeper silence the pin exposed — value-atom
projection over a non-single-valued dimension — is pinned too, collapsed on the
RDR's own operator matrix rather than on a vote (`eq`/`in` are restricted to the
single-valued-shaped kinds; `contains` is the co-occurring-value operator with
already-fixed semantics). Runs 2 and 3 would reconstruct against a draft that no
longer contains the silence they were called to sample.

The residual risk is real but is not a repeatability question: the projection
clause is a **new** load-bearing claim authored this pass and exercised by no
fixture. That is booked as **A21** (Method: MVV Test) with a plan that extends
MVV Scenario 4's negative control and declares the marker across the Phase 3
fixture. Stage 6 and the MVV are the instruments for an unexercised clause; a
cross-model re-draw is not.

Decision taken by the driver on the resolve pass, 2026-08-22.
