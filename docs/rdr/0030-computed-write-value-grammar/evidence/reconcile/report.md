Model: claude-opus-5[1m]

# Stage 6 Reconcile — cli/0030

Iteration 1. Preflight: `--outcome lens` = `/rdr-reconcile`, row complete
(`lens-foundational-row-complete`); `--outcome critique` = `none`
(two base models, diffed); `--outcome repeatability` = `none` (variant
full, runs 1/2/3 + diff). `Determinacy: fired` is written in Normative
Contracts, so the chain resolved without a route to run 1.

## Open set — four sources

1. **Pre-Lock needs-verification list** — the nine Pending Critical
   Assumptions the lenses left in the record: A5, A6, A10, A11, A12, A13,
   A14, A15, A16. No separate list file; the record carries them.
2. **Status Pending/Unverified** — `ca_pending_ids` is exactly those nine;
   `ca_unverified=0`, `ca_off_vocabulary=0`. Source 2 = source 1.
3. **`spikes_unrun`** — `[]`. The one named spike (A7 row growth) has its
   artifact under `evidence/spikes/a7-row-growth.md`. Plus the absorption
   audit over cove / 3amigo / critique / repeatability / author-round:
   **zero unabsorbed findings**, so no findings file names a spike the
   record does not. Nothing added.
4. **Exactness-word delta** — `recs lint 0030` emits no `prose:exactness`
   finding. Lint header: `blocking=0 resolution=0 placeholder=5
   advisory=6`; all five placeholder hits sit at 1997-2116, the
   Finalization Gate region, which is Stage 7's to move to `gate.md`, not
   a Stage 6 obligation. Nothing added.

Completeness: no `_Draft placeholder._` and no seed-skeleton header
survives in any body section; `## References` is fully authored (peer
clause ids, source anchors, authoring-guide sections, prior-art record,
the related issue) with no bracketed template text.

## Dispositions

| item | src | disposition | evidence / plan |
|---|---|---|---|
| A12 evaluator relocation | 1,2 | VERIFIED | Source Search vs `main`: `grammar.go` imports only stdlib + `internal/resolve`, names no `table` type; `::intWidth` is `int`-only with three `guard`-side callers, hence exported at its new home; `::IntDomain`/`::domainSize` carry `table.TagDecl` and stay; `resolve.GuardEvaluator` + `guardcontract.go` already there; `go list -deps ./internal/table` = `resolve`, not `guard`/`graphlint`. `evidence/reconcile/source-verify.md` §A12 |
| A14 step/clear collision | 1,2 | VERIFIED | Source Search: write loop then clear loop in one `::renderWrites` body, no branch between; clear ends with unconditional `assignments[key] = []string{ClearSentinel}` and mints no collision refusal; only refusals are `CatUnknownTag` / `CatWriteToNonOwnedTag`; `::conform` runs inside the write loop, before the clear loop. §A14 |
| A16 `expand` totality | 1,2 | VERIFIED | Source Search: `::expand` package-level, `[]Row`, no error; `::renderWrites` is the `*loader` method reaching the decl via `l.model.Tags[key]` and already calling `::conform`; return triple `([]TagValue, []string, error)` as stated; `TagValue` is `{Key; Value []string}`, literal-only; `::normalizeRule` hands `predicates` to `expand` alone. §A16 |
| A15 lint attribution | 1,2 | VERIFIED | Source Search: eight per-row producer sites (five in `groups.go`, two in `analysis.go`, one in `coverage.go`), every one holding the `Row` as a loop variable; four group-level sites disjoint from them; `::ruleIDsOf` bare and unchanged; no production `Row.Identity()` caller; recovery exact because `#` is banned in both rule ids (`normalize.go::normalizeRules`) and candidate members (`load.go`). Read-back census found a SECOND reader beyond `::groupHasOverlap` — `engine.go::identityKey` — verified SAFE and now enumerated in A15, C3 and the authority table. `evidence/reconcile/source-verify-a15.md` |
| A5 no new overlap | 1,2 | DOWNGRADED | Unrunnable before Phase 2 by construction. Plan: MVV items 2 and 5, scenarios S1 (both tiers), S5b, S5c. |
| A6 plan equality | 1,2 | DOWNGRADED | Same construction, with A5. Plan: MVV item 3, scenario S2. |
| A10 subsumption | 1,2 | DOWNGRADED | Static half already discharged by A1 (`::compareAtoms` order, `BlockAll` before `BlockMatch`). Plan: scenario S8, authored + context-inherited arms. |
| A11 both carriers | 1,2 | DOWNGRADED | Census half is a recorded source fact; the population is what needs a build. Plan: scenario S10, the discriminating oracle. |
| A13 bound-before-render | 1,2 | DOWNGRADED | Premise verified — `::conform` calls `::conformKind` before `::conformDomain`, and `conform` is reached only from the write loop this record widens (§A14 d). Refusal TEXT needs the build. Plan: scenario S5's A13 overflow fixture. |

### The A15 read-back finding — judged, not waved through

The A15 search returned one unenumerated consumer: `engine.go::identityKey`,
the key `::sortFindings` orders on, reads both `Rule` and `Element` and was
named in neither A15 nor C3. The sub-agent flagged it BLOCK per its brief's
instruction to report any such find as blocking; this stage judges it
**not a refutation**. The hard rule's test is "the design says X, the source
does not-X". Here the source says X: `identityKey` branches on the slot being
NON-EMPTY to pick a sort namespace and then uses the raw string only as a
within-namespace tie-break, so a suffixed value stays in the same bucket, the
rule-before-element ordering holds, and no test pins the literal string. It
also never resolves the field against authored `[[rules]]` ids by equality,
which is exactly what C3's fence forbids — so C3's normative rule was already
TOTAL over both readers and no clause changes. What was short was the
ENUMERATION, in an assumption whose whole point is that a producer census
misses consumers; its own consumer census then stopped one short of that
standard. Closed by naming the second reader as verified-safe in A15's
evidence, in C3's supporting paragraph, and in the authority table's consumer
cell (§amendment-sweep: the three sites a grep for `groupHasOverlap` /
read-back reaches; no other site claims a single reader). No route-back, no
punt-ledger row — nothing was refuted and no stage is owed rework.

### On the MVV-defer hard rule

Five items (A5, A6, A10, A11, A13) are declared `MVV Test` against a
form that does not exist: the fixture pair cannot be authored while the
loader refuses the shape, so these are unrunnable before Phase 2 **by
construction**, not by choice. The gate's rule — nothing the MVV depends
on may defer past lock — is satisfied not by running them now (impossible)
but by pinning each to a NAMED authored scenario that must pass before the
implementation is accepted. Every one of the five names its scenario in the
Testing Strategy, and those scenarios are written with their own
discriminating oracles (S10 fails a `Writes`-only build that passes
everything else; S11(c) fails a naive suffixing that passes everything
else). That is the disposition this stage makes explicitly: DOWNGRADED
with a pinned plan, not deferred silently and not falsely marked Verified.

Where a Pending item had a half that IS decidable on `main`, that half was
run rather than carried: A13's ordering premise, A11's consumer census and
A10's sort-order fact are each grounded above, so what remains Pending is
strictly the behaviour of code not yet written. Two items declared
`MVV Test`/`Source Search` against an unimplemented form turned out to be
fully decidable on `main` after all and were VERIFIED rather than
downgraded: A15 (a payload change whose producer and consumer censuses are
both facts about today's `internal/graphlint`) and A12/A14/A16 (all three
statements about today's call order and signatures). That is the split this
stage owed: unrunnable-by-construction is a property of the SUBJECT, not of
the declared Method, and four of the nine survived the check.

## Verdict

**RECONCILED.** Nine open items, nine terminal dispositions: four VERIFIED
by source search against `main` (A12, A14, A15, A16), five DOWNGRADED with
named, authored MVV scenarios (A5, A6, A10, A11, A13). No BLOCKER. No
refutation. Final `recs lint 0030`: `blocking=0 resolution=0`, the residual
placeholder findings all in the Stage-7 gate region.

Carried to Finalize, both advisory and both Gate-time judgments rather
than Stage 6 debts:
- The Finalization Gate's five responses are still inlined at 1997-2116
  and owe the move to `artifacts/gate.md`, which is Stage 7's own `lock`
  step.
- Two `evidence:over-budget` advisories, on A12 and A15, are NEW from this
  pass: both Evidence fields now carry their Stage 6 verification inline.
  Not truncated, per the finding's own remedy — the mass is verification
  content, every load-bearing anchor is a `path::Symbol` that stays
  findable, and the balance would move to `{ARTIFACT_DIR}` only at the
  Gate's discretion.
