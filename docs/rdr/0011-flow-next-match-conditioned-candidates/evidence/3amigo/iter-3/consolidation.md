Model: claude-opus-5[1m] (dispatcher; persona passes stamped in their own files)

# 3amigo Consolidation — RDR cli/0011, iter-3

Three isolated persona passes (`persona-1-pm.md`, `persona-2-implementer.md`,
`persona-3-qa.md`), each a separate sub-agent seeing only the RDR and its own
persona block. No file references another's output — isolation held.

**Why iter-3 is a FULL pass, not a delta.** The iter-1/iter-2 3amigo evidence
reviewed the pre-route-back draft, whose approach was a KERNEL three-valued match
seam. The critique route-back and re-propose (`fc93bfc`, "CLI presence-split,
kernel untouched") replaced that approach wholesale — it is now `0011:ALT4`,
rejected — so the earlier ledger's open entries do not carry, the same reason the
grounding lens re-swept at iter-2. The findings below are this cycle's origin
ledger.

Consolidated **mechanically**: hotspots are the set intersection over the element
ids each persona anchored to, not a re-judgement by a model that read all three.
Overlap marks a hotspot passage, not a validated finding; a single-persona
finding is not thereby weaker.

## Hotspots (≥2 personas anchored the same id)

| id | personas | what they independently landed on |
| --- | --- | --- |
| `0011:C1` | PM, Impl, QA | the contract every finding orbits: no diagnostic surface for an excluded row, the reader-demand gap that voids the presence test, dedup/sort unpinned |
| `0011:C2` | PM, Impl, QA | `--all` strips match facts — which breaks F1's remedy (PM), collides with C1's owned-key walk (Impl), and leaves the structural negative homeless (QA) |
| `0011:C3` | PM, Impl, QA | the documentation census is short (README, contract doc), `stringsAt`'s fate unstated, registration vs. structural negative conflated |
| `0011:MVV` | PM, Impl, QA | step 5 names a row that is not in the fixture it names; the success criterion and the oracles quantify differently |
| `0011:S1` | PM, QA | the narrowing scenario vs. the universally-quantified outcome claim |
| `0011:S7` | PM, QA | `uncomparable` reporting and what is actually pinned |
| `0011:§mini-checks` | PM, QA | the `oracle` table's negative control for MVV 5 dies with the fixture error |
| `0011:D-undecided-reporting-shape` | Impl, QA | text mode's `key (reason)` rendering has no carrier and no scenario |
| `0011:A11` | PM, Impl | CLI/kernel view agreement — and the key set that never reaches either |

**The strongest signal**: Impl-S1 and QA-1 are independent, and both are
*vacuity* defects — a contract that silently degrades to the mode it was written
to replace, and an oracle that passes in every possible build. Neither is
reachable from the draft's own text; both needed the shipped code. Separately,
PM-H1 and PM-H2 are the same defect reached from two directions: the payload has
no match block in either mode, so the Failure Mode's remedy and the Briefly
Rejected ground are both written against information the design does not emit.

## Origin ledger

| L | origin | finding (one line) | anchors | sev |
| --- | --- | --- | --- | --- |
| L1 | Impl-S1 | an owned key that is only MATCHED on invokes no reader, so it is absent from the view and C1 silently degrades to `--all` for that model | `C1`, `A11`, `JC1`, §existing-infrastructure-audit | HIGH |
| L2 | QA-1 | S2/MVV-5 assert `gated-excluded` absent from `flowMVVModel`, but that row lives in `flowGatedNextModel` — vacuous in every build | `S2`, `MVV`, `§mini-checks` | HIGH |
| L3 | QA-2 | S4's "row with an empty match pattern" is unbuildable — `normalize.go:378` refuses an authored-empty match block | `S4`, `C1` | HIGH |
| L4 | PM-H1 | F1's remedy routes to a `dump` verb that does not exist, and to `--all`, which C2 strips of match facts — an excluded row has no diagnostic surface | `F1`, `C2`, `§failure-modes` | HIGH |
| L5 | PM-H2 | BR3 rejects `excluded_by` because "`--all` already yields" the information — C2 says the opposite | `BR3`, `C2` | HIGH |
| L6 | Impl-S2 | C2's "filter BlockMatch out of `unknown`, an oracle MUST assert their absence" collides with C1's owned-key walk on any key both matched and written (every `models/rdr.toml` rule) | `C2`, `C1` | HIGH |
| L7 | Impl-S3 / QA-3 | `D-undecided-reporting-shape` mandates `key (reason)` text rendering; `respond`'s generic `flatten` cannot produce it, and no Phase/Audit/Overrides row scopes the gateway edit | `D-undecided-reporting-shape` | MAJOR |
| L8 | QA-4 | `TestReq3` is presence-only, so C3's registration clause and C2's structural negative are two oracles; only one has a named home | `C2`, `C3` | MAJOR |
| L9 | QA-5 | S3/MVV-6's fixture is specified by row property, not model; A5's shape forces a dummy owned key into `unknown`, so "present-equal ⇒ nothing in `unknown`" is false as written | `S3`, `MVV` | MAJOR |
| L10 | Impl-S4 | C1 states the restriction over ATOMS, but `KernelRow()` converts match atoms to `[]resolve.Tag`; C1 names no seam, and rebuilding from `row.Atoms` would fork `seamValue`'s canonicalization | `C1`, `A12` | MAJOR |
| L11 | Impl-S5 | C3 doesn't say what `stringsAt` becomes — a shared harness helper used well beyond the five reads; a new helper is a sixth touched site with no census slot | `C3`, `S6` | MAJOR |
| L12 | PM-M1 | §consequences says "without a contract-doc edit"; Phase 3 plans that exact edit; `cli-output-contract.md:125` still says "Enumerate" | `§consequences`, `§phase-3-fixtures-and-docs` | MEDIUM |
| L13 | PM-M2 | `README.md:43` ("list legal next outcomes") is a shipped user-facing description absent from C3's census | `C3`, `§consequences` | MEDIUM |
| L14 | PM-M3 | the success criterion is universally quantified and explicitly "not a target cardinality", but every oracle is existential and MVV 3 is a cardinality assertion; `G-scope` is still template placeholder text | `§problem-statement`, `MVV`, `G-scope` | MEDIUM |
| L15 | QA-6 | no fixture produces a "fully-resolved candidate" for MVV 7; C1's "every owned key no invoked reader established" may widen `summarize`'s current `RequiresOwned` scope | `MVV`, `S2`, `C1` | MINOR |
| L16 | QA-7 | sort stability is asserted nowhere; nothing in the record discriminates pair-dedup from key-dedup | `C1`, `S7` | MINOR |
| L17 | QA-8 | S6's "none re-homed under `--all`" is a diff claim with no automatable criterion | `S6` | MINOR |
| L18 | QA-9 | S5 states the kernel-untouched criterion three ways; only MVV 8's `git diff --stat` form is mechanically checkable | `S5`, `MVV` | OBS |
| L19 | PM-L1 | `{key, absent}` does not distinguish "bind a reader" from "supply a tag" — the fork the caller actually faces; the Consequences bullet overstates the remedy | `§consequences`, `A5` | LOW |
| L20 | PM-L2 | every `--all` precedent cited widens to equally-informative rows; here it widens to strictly less informative ones | `D-naming` | LOW |

## Grounding already executed by the dispatcher

Confirmed on `main` before any disposition (the resolve half's gate 1):

- **L1** — `flow_exec.go::guardOwnedKeys` skips any atom that is not `BlockAll`/
  `BlockUnless` (`:133-135`); `RequiresOwned` is returned by `renderWrites`
  (`normalize.go:438`), i.e. writes + clears only. A match-only owned key is in
  neither set. CONFIRMED.
- **L2** — `gated-excluded` is authored at `flow_fixtures_0005_test.go:317`,
  inside `flowGatedNextModel` (opens `:245`); `flowMVVModel` spans `:77-170`.
  CONFIRMED.
- **L3** — `normalize.go:378-380`: "both an absent `[rule.match]` and a
  present-but-empty one decode to nil, and both are equally malformed" → refused
  `CatMalformedRuleShape`. CONFIRMED.
- **L4** — shipped verbs are `version`, `lint`, `flow next|resolve|read-state|
  set-state`; no `dump` (`flow_surface_0005_test.go:105` asserts `dump` is
  FOREIGN to the flow group). `candidate` (`flow_next.go:39-58`) carries no match
  block. `owned`/`observed` ARE on the payload (`:32-33`), so half of F1's
  mechanism exists; the row's match block is only in the model file. CONFIRMED.
- **L7** — `respond.OK` → `writeTextPayload` → `flatten` (`respond/text.go:42`,
  `:65`), a generic reflective flattener over the marshalled JSON; it would emit
  `candidates[0].unknown[0].key:` / `.reason:`. No verb-specific renderer exists.
  CONFIRMED.
- **L12** — `docs/cli-output-contract.md:125` reads "Enumerate the legal outcomes
  and their candidate rules". CONFIRMED.
- **L13** — `README.md:43` reads "list legal next outcomes". CONFIRMED.
- **L14** — `0011:G-scope` projects as verbatim template bracket text. CONFIRMED
  (also a §mechanical-gate hit).

## Dispositions (resolve half, same cycle)

| L | disposition | section touched |
| --- | --- | --- |
| L1 | **fixed** | C1 — new demand-set paragraph requiring `invokedReaders` to add match-block owned keys, with the degradation it prevents and why `checkAccessorBindings` does not catch it; Audit `Reader narrowing` row Reuse→**Extend**; `JC1` shared-anchor clause; `authority` mini-check gains a demand-set row; Approach placement ⇒-clause; Consequences blast-radius bullet; Phase 1. New **A15** (Pending) + **S8** + MVV 8 + `oracle`/`trace` rows |
| L2 | **fixed** | MVV 5, S2, `trace` step 5 — the `gated-excluded` assertion re-homed to `flowGatedNextModel` (where `:317` authors it), `gated-reported` added as negative control, and the vacuity named so it cannot be re-introduced |
| L3 | **fixed** | S4 — the unbuildable authored-empty-match row replaced by the reachable all-atoms-omitted form; `disposition` table row rewritten to distinguish the probe-minted empty PATTERN from an authored empty block (refused at load) |
| L4 | **fixed** (remedy corrected) + **charted** (item 1) | F1 — the `dump` reference and the false `--all` remedy replaced with what the payload actually supports, and the absence of a model-inspection verb stated |
| L5 | **fixed** (ground corrected) + **charted** (item 1) | BR3 — the false "`--all` already yields" ground replaced with the real rejection reason, and the reopening path named |
| L6 | **fixed** | C2 — the filter scoped to the ATOM WALK's contribution only, with the matched-and-written collision named; C3's absent-key fixture must write a DIFFERENT key |
| L7 | **fixed** (as descriptive) + **charted** (item 3) | `D-undecided-reporting-shape` — `key (reason)` reframed as "the pairs are carried", with the gateway invariant and its 0005 clauses cited |
| L8 | **fixed** | C3 — `TestReq3` named presence-only, so the positive registration and C2's structural negative are stated as two oracles with distinct homes |
| L9 | **fixed** | S3 + MVV 6 — the assertion moved to the MATCH key's absence from `unknown`, never an empty list |
| L10 | **fixed** | C1 — the restriction pinned as a post-`KernelRow()` filter over `probe.Match` by `Tag.Key`, with the `seamValue` fork named as the reason |
| L11 | **fixed** | C3 — `stringsAt` explicitly NOT repointed; a sibling harness helper is added, and it is harness rather than a sixth re-homing |
| L12 | **fixed** | Consequences + Phase 3 — the contradiction resolved: no contract CLAUSE is invalidated, but the doc's prose states the old default and Phase 3 corrects it |
| L13 | **fixed** | C3 census + Consequences + Phase 3 — `README.md:43` and `cli-output-contract.md:125` added; C3's "the ONE shipped user-facing string" corrected to "in this file" |
| L14 | **fixed** | Problem Statement — the quantifier gap stated (universal criterion, existential oracles, closed by the `disposition` enumeration). NOT written into `G-scope`: the Finalization Gate's sub-sections spec `gate.md`'s content and are authored at Stage 7, never inlined here |
| L15 | **fixed** | MVV 7 — the fully-resolved-candidate fixture must be supplied, with the conditions that make one (incl. `--evaluate-gates`, else gate ids fill `unknown`) |
| L16 | **fixed** | S7 — sort-by-`(key, reason)` given an oracle; the pair-vs-key dedup distinction recorded as NOT discriminable by any reachable input, rather than asserted |
| L17 | **fixed** | S6 — "none re-homed under `--all`" reclassified as a diff-review claim with no runtime oracle, and tied to refuting A4 |
| L18 | **fixed** | S5 — the kernel-untouched criterion collapsed to MVV 9's `git diff --stat` as the one mechanical form, the suite demoted to corroboration |
| L19 | **fixed** (limit stated) + **charted** (item 2) | Consequences — `absent` names the key, not which remedy; provenance not carried |
| L20 | **fixed** | `D-naming` — the precedent asymmetry recorded (cited `--all`s widen to equally-informative rows; this one widens to strictly less informative) |

No entry was dismissed-with-cite: every finding grounded. Three were charted in
part (L4, L5, L19 → `Charted.md` items 1–2; L7 → item 3) with the in-scope half
fixed in the draft.

### Tiebreakers escalated: none

Two forks were collapsed on evidence rather than escalated:

- **L1's fork** ("extend `invokedReaders`" vs "state the limit") collapsed once
  `runFlowNext` was read: `invokedReaders(req.model, "")` is called ONCE before
  the row loop and `--all` is never consulted in reader selection, so extending
  the demand set preserves C1's own "fixed once per invocation, same under
  `--all`" invariant instead of breaking it. The `JC1` risk collapsed too: 0010's
  `decision-table` class declares zero owned tags, so the added term is empty
  there and `0010:C4` holds verbatim.
- **L7's fork** (normative `key (reason)` vs descriptive) collapsed on
  `respond/text.go`'s own doc pinning the two-modes invariant to 0005 clauses this
  RDR does not override — making the normative reading an out-of-scope edit to a
  predecessor's surface.

### Amendment sweep (rdr-common §amendment-sweep)

The demand-set change (L1) was swept to every site stating the blast radius or
the reader-narrowing disposition: Approach ⇒-clause, Audit row, `JC1`,
`authority` mini-check, Consequences, Phase 1, Prerequisites. The `gated-excluded`
fixture correction (L2) was swept to all three sites (MVV 5, S2, `trace` step 5).
The S8/MVV-8 insertion renumbered MVV 8→9, swept to S5's citation and the `trace`
table's final row.

### Mini-checks

Fired cues are unchanged — `authority`, `oracle`, `disposition`, `trace` — and
their tables are live in the draft. This pass's fixes add no new cue: the
`authority` table gained a demand-set row and `oracle`/`trace` gained an S8 row,
both within tables the cue already owns. `round-trip`/`fidelity` stays absent.
No re-read owed.

### Needs (re)verification → Stage 6

- **A15 (NEW, Pending)** — the demand-set extension: invokes strictly more readers
  and never fewer; the added term is empty over the zero-owned-tag class
  (`0010:C4` unaffected); a previously-skipped refusing reader now aborting the
  verb is correct-but-changed behaviour needing a fixture; and the widening is
  uniform across callers rather than `next`-special-cased. Method: Source Search +
  MVV Test. This is the only assumption this pass moved.
- **A11 re-read, NOT flipped** — its claim is that the CLI's and the kernel's views
  AGREE on key set over the inputs `flow next` passes, which the demand-set change
  does not disturb (both still consume the one `owned` slice `runReaders` returns).
  Its Evidence needs no edit; recorded here so Stage 6 does not have to re-derive
  that it was considered.
- **A4 unchanged** — its BREAKS-0 is predicate-scoped; nothing this pass fixed
  touches the predicate's effect on the 0005 oracles. L17's reclassification makes
  A4 the thing "none re-homed under `--all`" would refute, which is a reporting
  clarification, not a status change.
- No spike is called for by any finding.
