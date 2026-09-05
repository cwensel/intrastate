Model: claude-opus-5[1m]

# cove pre-lock lens — 0019, SECOND (delta) pass

Repo HEAD swept: `679a536` (working tree carries the 0019 rewrite,
uncommitted). Projector-only reading of the record; source read directly
on `main`. Pass-1 findings at `../findings.md` read once for context; none
re-raised.

## Delta scope checked

Still-open anchors: `A5`, `BR3`, `C1`, `D-naming`, `D-selection-predicate`,
`MVV`, `S7`, `§existing-infrastructure-audit`.

New/materially-changed surface: `A6`, `§technical-design` (second-new-data-
flow paragraph), `§consequences` (carrier-scope negative), `§mini-checks`
(five tables), `RT2` (scoped to file-backed carrier).

Out of scope by instruction: `§finalization-gate`.

Elements pulled for cross-consistency (not re-critiqued as such):
`A2`, `RT1`, `RT3`, `RT4`, `F3`, `F5`, `S1`–`S9`.

## The nine delta questions

**Q1 — C1 class arm: three named load refusals, and are they exhaustive?**
**CONFIRMED on all three strings; the "three arms are exact" framing is
imprecise but not wrong (see Finding 3).**

- `written tag %s is served by %d writers; want exactly one` —
  `internal/table/load.go:1570-1572`, reached via the `[initial]` fold at
  `internal/table/load.go:1538-1540`. Reproduced: `[initial]` over the
  recognized tag of `models/examples/routing-decision-table.toml` →
  `malformed_accessor_binding: written tag recognized is served by 0
  writers; want exactly one`.
- `writer %s names the non-owned tag %s` —
  `internal/table/load.go:928-929` (`CatWriteToNonOwnedTag`,
  `internal/table/category.go:24`), in `(*loader).loadAccessors`.
  Reproduced with an observed tag (`severity`) plus a matching writer →
  `write_to_non_owned_tag: writer w names the non-owned tag severity`.
- `[model] class %s declares owned=%d; a %s model declares zero tags of
  provenance %s` — `internal/table/load.go:833-838`,
  `(*loader).checkClassAgreement` at `internal/table/load.go:819`.
  Reproduced by declaring an owned `scalar` tag →
  `malformed_model_declaration: [model] class "decision-table" declares
  owned=1; a decision-table model declares zero tags of provenance owned`.

Exhaustiveness: the CLAIM (no decision-table model carrying `[initial]`
can load) holds — I found no fourth construction that loads. But I did
find a **fourth distinct refusal spelling**: a writer over a *recognized*
tag refuses `malformed_accessor_binding: write w names the reserved key
recognized` BEFORE the non-owned check, so C1's "adding a writer for it
is `write_to_non_owned_tag`" is true only for observed tags, not for the
recognized tag every decision-table model declares. Class (c), minor —
Finding 3.

**Q2 — C1 conform passage: loader already conforms; byte equality from
`::canonicalSet` alone.** **CONFIRMED against source; but it CONTRADICTS
two sibling elements the rewrite did not update (Finding 1).**
`internal/table/load.go:1701` `conform(decl, "eq", members)`; `conform`'s
non-operator fallthrough at `internal/table/load.go:1839-1848` runs
`conformKind` then `conformDomain` per member — the exact body of
`table.ConformValue` (`internal/table/load.go:1813-1818`). Encoder:
`internal/cli/flow_input.go::canonicalSet` at `internal/cli/flow_input.go:770`,
doc comment `:764-769` naming it "THE encoder … read-back equality byte
equality (REQ-71)". C1's new sentence is correct.

**Q3 — C1 ALL quantifier over every bound artifact; multi-writer
iteration.** **CONFIRMED in source; consistent in every place I checked.**
`internal/table/model.go:521` `Writers map[string]Accessor`.
`internal/cli/flow_state.go:351-362` — `runFlowSetState` loops
`for _, name := range slices.Sorted(maps(byWriter))` and checks
`req.artifacts[def.Accessor.Role]` per writer; a second identical loop at
`:367-378` executes each writer separately. C1's description is exact.
The ALL rule is stated consistently in C1 ("EVERY bound artifact carries
NO key … The quantifier is ALL, not ANY and not per-artifact"),
`D-selection-predicate` ("iff EVERY bound artifact carries no key — the
ALL quantifier, not ANY and not per-artifact"), `RT2` ("the bytes of
EVERY bound artifact are unchanged"), the `disposition` table ("Every
bound artifact empty" / "Any bound artifact non-empty"), and the `trace`
table ("ALL bound artifacts, STORE keys"). `RT1`/`RT4`/`MVV`/`S1`–`S9`
are single-artifact fixtures and say nothing that conflicts. **No
inconsistency found — pass-1 Finding 5 is genuinely closed.**

**Q4 — C1 carrier scope: non-file-backed carriers refused.** **Selection
logic CONFIRMED; the refusal is consistently stated in three of four
places and MISSING from the testing scenarios (Finding 2).**
`internal/cli/flowbind/registry.go:60-82` — default
`binding = &Writer{Path: acc.Path}`, `case len(acc.Edit) != 0: binding =
NewEditWriter(acc, name)` (`:78-79`), `case commandBacked(acc): binding =
&cmdbind.Writer{…}` (`:80-81`); `commandBacked` at
`internal/cli/flowbind/registry.go:121-123`. `NewEditWriter` at
`internal/cli/flowbind/edit.go:74`. Refusal is reflected in C1's CARRIER
SCOPE paragraph, `D-selection-predicate` ("with the other two
(`::NewEditWriter`, `cmdbind.Writer`) refused, per C1"), the
`disposition` table ("Non-file-backed write carrier | 2 | zero | carrier
refusal"), `§consequences` (the new negative), and the `trace` table
("carrier check … edit/command-backed writer → exit 2"). **It is absent
from `§testing-strategy` (`S1`–`S9`) and from `MVV` — a silence.**

**Q5 — emptiness-is-a-new-data-flow, and A6's three candidate carriers.**
**Emptiness gap CONFIRMED; A6's carrier enumeration is REFUTED (Finding 4).**
`internal/accessor/binding.go:34-42` — `ReadBinding.Read(ctx, art,
requested) (values []KeyValue, unreadable []string, err error)`, no
cardinality. `internal/cli/flowbind/flowbind.go:108` `type store
map[string]string` (unexported); the package's exported surface is
`Reader`/`Writer`/`Gate` and `Registry`/`OwnedTags`/`NewEditWriter`, none
returning a count. Three candidate CARRIERS in A6's proposition line
(new `ReadBinding` capability / exported `flowbind` probe /
`read-state`-family surface) are all real surfaces. But A6's *Evidence*
names the `ReadBinding` IMPLEMENTERS as "(`flowbind.Reader`, `cmdbind`,
edit)" — **there are only two**: `internal/cli/flowbind/flowbind.go:189`
`func (r Reader) Read(…)` and `internal/cli/cmdbind/cmdbind.go:1123`
`func (r Reader) Read(…)`. `EditWriter` implements no `Read` at all
(`internal/cli/flowbind/edit.go` — `Capability() → accessor.CapWrite` at
`:81`, no `Read` method), and the registry builds an edit binding only in
the WRITE loop (`registry.go:78`, whose own comment says "`edit` is a
WRITE carrier only"). Class (a) — Finding 4.

**Q6 — D-naming's new PREDICATE ground for rejecting `--from-initial`.**
**SOUND, and no surviving contradiction.** `D-naming` now states the
refuted ground explicitly as refuted and cites `--plan`
(`internal/cli/flow_input.go:273-296` `planFlagName`;
`internal/cli/flow_state.go:434-470` `parseWrites`, duplicate refusal at
`:455-457` naming "`--plan`, `--write` and `--clear`") — verified. `BR3`
was rewritten to match ("NOT rejected for carrying a second write-plan
source — `--plan` already is one"). I found no other passage still
asserting the two-sources ground: `ALT3`, `§decision-rationale`, and
`§approach` were checked and are silent on it. The new ground —
`--from-initial` would make a verb whose grammar is "commit the writes I
named" conditional on store emptiness, with no-op success
indistinguishable from a commit — is internally consistent with C1's
no-op-success arm. **Pass-1 Finding 4 is genuinely closed.**

**Q7 — the five mini-check tables.**
- `authority`: every row checks out. "Store emptiness | … | **no shipped
  reader** — carrier undecided (A6)" agrees with C1 and
  `§technical-design`. "Model class | loader `[model] class` (0010:C1) |
  `::IsDecisionTable` callers | `len(owned)==0` re-derivation" agrees with
  `internal/table/model.go:60-70` (doc comment says "never re-derives the
  class from `len(owned) == 0`"). "Value encoding | `::canonicalSet` …
  Sibling arms: loader `conform` (kind/domain only)" is precisely
  consistent with C1's rewritten conform passage. "Writer for a key |
  loader `::checkAccessorBindings` (exactly-one)" —
  `internal/table/load.go:1566-1573`. CLEAN.
- `oracle`: all five negative controls are constructible. Row 4's
  ("move the class check after routing → artifact created") is
  constructible only as a *code mutation*, not as a fixture — but it is
  labelled a negative control, i.e. a mutation-style check, and rows 2,
  3 and 5 are the same shape ("remove the predicate", "per-key variant",
  "count owned keys only"). Consistent within the table. CLEAN.
- `fidelity`: rows agree with `RT1`–`RT4`. Row 2's "Lossy exemption:
  file-backed carrier only" matches `RT2`'s new scope note. Row 3's "the
  3 kinds only the loader admits — no argv spelling exists" matches
  `RT3`. CLEAN as a table; but see Finding 1, which lands on what "the 3
  kinds" are.
- `disposition`: exit codes agree with C1's prose on all seven rows
  (empty→0/seed, non-empty→0/no-op, sealed→0/no-op, decision-table→2,
  routing/binding→2, non-file-backed carrier→2, read-back mismatch→
  non-zero/partial). C1's read-back arm says "a distinct terminal
  refusal", and the table says "non-zero" rather than pinning a code —
  consistent, deliberately unpinned. CLEAN.
- `trace`: I checked every row's "assertions in force" against the cited
  element and found no pair in conflict, so the declared "No
  CONTRADICTION row" is **correct** — the one unresolved cell (emptiness
  carrier, A6) is genuinely under-specification, not conflict. CLEAN.
  Note the `carrier check` row is the ONLY place a carrier-refusal
  witness is written down, and it is not a testing scenario (Finding 2).

**Q8 — S7 / MVV step 8: "no accessor process ran and no artifact was
created or modified".** **Discriminating AND constructible. CONFIRMED.**
A loadable decision-table model declares zero owned tags
(`checkClassAgreement`, `internal/table/load.go:819-839`) so it can bind
no owned-tag write accessor — S7 says exactly this and correctly
abandons the "writer that fails if invoked" instrument. The replacement
assertion is discriminating because every write path increments an
observable counter and touches the file: `internal/cli/flowbind/flowbind.go:251`
`w.invocations++` then `save(art.Path, s)` at `:283`
(`(*Writer).Invocations()` is exported at `:235`, and `EditWriter` has the
same pair at `edit.go:86`). So "artifact path still absent" is a real
discriminator: a refusal raised after routing would have called `Apply`,
which unconditionally `save`s. Constructible on a decision-table fixture:
the fixture needs no writer at all, and the assertion is on the artifact
path, not on a writer. **Pass-1 Findings 1 and 9 are genuinely closed.**

**Q9 — cross-element consistency after the four-site C1 edit.** C1's four
edited regions (class arm, conform passage, ALL quantifier, carrier scope
+ emptiness-is-a-new-data-flow) agree with each other, with `A5`, with
`RT1`–`RT4`, with `D-selection-predicate`, and with all five mini-check
tables. `§technical-design`'s new paragraph now says "The FIRST new data
flow is the *source*…" and "There is a SECOND new data flow", correctly
retracting the sentence pass-1 refuted; C1 still quotes the old sentence
("The Technical Design's claim that 'the only new data flow is the SOURCE
of the planned writes' is therefore wrong") — that is now a quotation of
a sentence that no longer exists in `§technical-design`, which reads as
narration of the rewrite rather than as contract. Minor; noted under
Finding 3 rather than raised separately, since the sentence is
argumentatively correct and no reader is misled about the design.

**The rewrite did NOT propagate C1's conform retraction into `A2` and
`F5`** — the one genuine internal contradiction. Finding 1.

## NEW findings

**Finding 1 — `0019:A2` and `0019:F5` — class (c), internal
contradiction with the rewritten C1.**
C1's rewrite now states, load-bearingly: "Conformance to the declaration
is ALREADY HELD by the loader and is **not re-established here** … A
re-conform pass on a loader-normalized value cannot fail and buys
nothing; the byte-equality property comes from the encoder alone." Two
siblings still assert the opposite, both unedited:
- `A2` Evidence: "§D13's canonical form … is preserved by **conforming
  each member with `table.ConformValue`** and rendering through the same
  `::canonicalSet`." A2 credits half the byte-equality property to a step
  C1 has just deleted.
- `F5`: "**C1's conform step** accepts it with no clause change" — F5's
  entire forward-compatibility disclosure for the empty-scalar asymmetry
  is anchored to a C1 step that no longer exists. If the empty scalar is
  later admitted, the record now gives no clause that governs it, so F5's
  "no clause change" promise is unbacked.
Source agrees with C1, not with A2/F5 (`internal/table/load.go:1701`,
`:1839-1848`, `:1813-1818`). This is the pass-1 Finding 6 fix pinned two
ways across three elements — the exact "pins one field two ways" defect
the delta pass is for. F5 is the sharper half: it makes a normative
promise about a clause that was removed.

Secondary, same class, same cause: `RT3` says "the three kinds only the
loader admits (bare scalar for a set tag, array for a scalar tag, **empty
scalar**)", while `S4` says "the **3 kinds** only the loader admits (bare
scalar for a set tag, array for a scalar tag)" — names 3, lists 2, and the
omitted one is the empty scalar that `F5` places out of scope and that
`A2` records as unwritable by either route. S4's parenthetical is
therefore both an undercount and, if read as authored, a demand to test a
kind the record says is unwritable.

**Finding 2 — `0019:§testing-strategy` (`S1`–`S9`) and `0019:MVV` —
class (d), silence on the newly-added refusal.**
The rewrite added a whole new terminal refusal — "Against a model any of
whose bound write accessors is non-file-backed, init-state MUST refuse —
a distinct terminal refusal in the `flow-*` family naming the accessor
and its carrier" (C1, CARRIER SCOPE) — and reflected it in
`D-selection-predicate`, the `disposition` table, `§consequences`, and
the `trace` table. **No scenario in `S1`–`S9` and no `MVV` step covers
it.** Every other terminal arm C1 names has a scenario: class refusal S7,
writer routing S5, artifact binding S6, read-back mismatch S8, sealed
store S9. The carrier refusal is the only one with no oracle, and it is
the one guarding an emptiness answer C1 says is UNDEFINED for the other
two carriers — i.e. the arm whose failure mode is silently seeding on a
wrong answer, exactly what the clause exists to prevent. It is also
cheaply constructible: an `[write.x.edit]` accessor or a
`[write.x] command = [...]` accessor on the S1 fixture
(`internal/cli/flowbind/registry.go:78-81`, `:121-123`).

**Finding 3 — `0019:C1` — class (a), the class arm's refusal enumeration
is incomplete for the recognized tag.**
C1 asserts "The three load arms are **exact**" and that "adding a writer
for it is `write_to_non_owned_tag`". Reproduced at HEAD against
`models/examples/routing-decision-table.toml`: adding
`[write.w] keys = ["recognized"]` refuses
`malformed_accessor_binding: write w names the reserved key recognized`,
NOT `write_to_non_owned_tag` — a reserved-key guard fires first for the
`recognized` tag, which every decision-table model in this repo declares
(`models/examples/routing-decision-table.toml:38-42`). Only for an
*observed* tag does the `write_to_non_owned_tag` arm at
`internal/table/load.go:928-929` fire. The CONCLUSION C1 draws
(unconstructible; `len(Model.Initial) == 0` for every admitted
decision-table model) is unaffected and remains CONFIRMED — I found no
loading construction — so this is an exactness defect in a passage that
says "exact", not an approach problem. Fix is one clause: either scope
the second arm to observed tags or drop "exact" for "on every path we
found".

Also folded here, sub-finding, class (c) cosmetic: C1's
emptiness-is-a-new-data-flow paragraph quotes and refutes
"The Technical Design's claim that 'the only new data flow is the SOURCE
of the planned writes'", but `§technical-design` no longer makes that
claim — it was rewritten to "The FIRST new data flow is the *source*…
There is a SECOND new data flow". C1 now argues against a sentence that
does not exist in the record, which reads as change-history narration in
a normative block.

**Finding 4 — `0019:A6` — class (a), the implementer enumeration in the
verification plan is codebase-refuted.**
A6's Evidence sets the verification method as "enumerate the three
candidate carriers against `::ReadBinding`'s implementers
(`flowbind.Reader`, `cmdbind`, **edit**) and confirm one can answer
emptiness for the file-backed carrier with no change to the other
implementers' behavior." There are exactly **two** `ReadBinding`
implementers in the repo: `internal/cli/flowbind/flowbind.go:189`
`func (r Reader) Read(_ context.Context, art accessor.Artifact, requested
[]string)` and `internal/cli/cmdbind/cmdbind.go:1123`
`func (r Reader) Read(ctx context.Context, art accessor.Artifact,
requested []string)`. The edit carrier implements no `Read`:
`internal/cli/flowbind/edit.go:81` `func (w *EditWriter) Capability()
accessor.Capability { return accessor.CapWrite }`, no `Read` method on
the type, and the registry constructs it only in the write loop
(`internal/cli/flowbind/registry.go:78`, whose comment states "`edit` is
a WRITE carrier only, which is why this arm appears in this loop alone").
The consequence is not cosmetic: A6 is Pending and its Method is the
instruction the resolver will follow. As written it sends the resolver to
check a non-existent implementer for a contract disturbance that cannot
occur, and — more costly — it makes the implementer set look like it
covers all three write carriers, when in fact an edit-carried write
accessor has NO read binding at all. That is independent evidence for
C1's carrier refusal (an edit-carried model cannot answer emptiness
because it has no reader, not merely because it has no JSON store), and
A6 should say so rather than mis-enumerate.

## What was checked and found clean

Recorded so a third pass does not re-walk them: the ALL quantifier
(Q3, consistent in six places), `D-naming`/`BR3`'s new predicate ground
(Q6, sound, no surviving contradiction), S7/MVV-8's refusal oracle (Q8,
discriminating and constructible), the `authority` / `oracle` /
`fidelity` / `disposition` / `trace` tables (Q7, all rows agree with C1
and with source; `trace`'s "No CONTRADICTION row" is correct), and
`§existing-infrastructure-audit`'s Class-discrimination row, which the
rewrite corrected to "Reuse (new caller layer)" with the
`internal/cli`-is-the-first-CLI-reader limit stated — pass-1 Finding 7
closed, verified against `internal/table/model.go:66-67` and a non-test
reference sweep.

Pass-1 Findings 1, 2, 3, 4, 5, 6, 7 and 9 are all closed by the rewrite.
Finding 8 (`§finalization-gate`) is out of delta scope by instruction.
